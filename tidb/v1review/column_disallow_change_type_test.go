package v1review

import (
	"testing"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// columnSchema is a snapshot of schema "" whose table "book" has the
// columns, each named and typed by a pair, nullable when its type is
// prefixed with "?", with the indexes, each listing its expressions; and a
// partition "book_p0" of it.
func columnSchema(columns [][2]string, indexes ...[]string) v1review.Schema {
	table := &metadata.TableMetadata{Name: "book", Partitions: []*metadata.TablePartitionMetadata{{Name: "book_p0"}}}
	for _, c := range columns {
		tp, nullable := c[1], false
		if len(tp) > 0 && tp[0] == '?' {
			tp, nullable = tp[1:], true
		}
		table.Columns = append(table.Columns, &metadata.ColumnMetadata{Name: c[0], Type: tp, Nullable: nullable})
	}
	for i, expressions := range indexes {
		table.Indexes = append(table.Indexes, &metadata.IndexMetadata{Name: string(rune('a' + i)), Expressions: expressions})
	}
	return v1review.Schema{Metadata: &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{Tables: []*metadata.TableMetadata{table}}}}}
}

func TestColumnDisallowChangeType(t *testing.T) {
	schema := columnSchema([][2]string{{"id", "int(11)"}, {"Name", "varchar(20)"}, {"flag", "tinyint(1)"}, {"amount", "decimal(10,0) unsigned"}, {"kind", "enum('a','b')"}})
	changed := func(text string) placed {
		return placed{Code: v1review.ColumnTypeChanged, Table: "book", Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []placed
	}{
		{name: "same types", sql: "ALTER TABLE book MODIFY id INT, MODIFY name VARCHAR(20), CHANGE flag f BOOLEAN, MODIFY amount DECIMAL UNSIGNED, MODIFY kind ENUM('a','b');"},
		{name: "unknown column or table", sql: "ALTER TABLE book MODIFY other BIGINT; ALTER TABLE shelf MODIFY id BIGINT;"},
		{
			name: "once per statement",
			sql:  "ALTER TABLE book MODIFY id BIGINT, CHANGE name n VARCHAR(30);\nALTER TABLE db.BOOK CHANGE id id INT UNSIGNED;",
			want: []placed{
				changed("ALTER TABLE book MODIFY id BIGINT, CHANGE name n VARCHAR(30)"),
				{Code: v1review.ColumnTypeChanged, Table: "BOOK", Text: "ALTER TABLE db.BOOK CHANGE id id INT UNSIGNED"},
			},
		},
		{name: "change looks up the old name", sql: "ALTER TABLE book CHANGE other id BIGINT;"},
		{name: "a partition finds its table", sql: "ALTER TABLE book_p0 MODIFY id BIGINT;", want: []placed{{Code: v1review.ColumnTypeChanged, Table: "book_p0", Text: "ALTER TABLE book_p0 MODIFY id BIGINT"}}},
		{name: "other actions are not judged", sql: "ALTER TABLE book ADD id BIGINT, ALTER COLUMN id SET DEFAULT 1, RENAME COLUMN id TO x; CREATE TABLE book (id BIGINT);"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			checkPlaced(t, tt.sql, func(stmts []ast.Node) []v1review.Finding {
				return ColumnDisallowChangeType(stmts, schema)
			}, tt.want)
		})
	}
	t.Run("object case sensitive", func(t *testing.T) {
		schema := schema
		schema.ObjectCaseSensitive = true
		checkPlaced(t, "ALTER TABLE BOOK MODIFY id BIGINT;", func(stmts []ast.Node) []v1review.Finding {
			return ColumnDisallowChangeType(stmts, schema)
		}, nil)
	})
	t.Run("no snapshot", func(t *testing.T) {
		checkPlaced(t, "ALTER TABLE book MODIFY id BIGINT;", func(stmts []ast.Node) []v1review.Finding {
			return ColumnDisallowChangeType(stmts, v1review.Schema{})
		}, nil)
	})
}

func TestColumnType(t *testing.T) {
	for _, tt := range []struct {
		dt   ast.DataType
		want string
	}{
		{ast.DataType{Name: "BLOB"}, "blob"},
		{ast.DataType{Name: "DECIMAL", Length: 10}, "decimal(10,0)"},
		{ast.DataType{Name: "NUMERIC", Length: 10, Scale: 2, Unsigned: true}, "decimal(10,2) unsigned"},
		{ast.DataType{Name: "INT", Length: 11, Zerofill: true}, "int(11) unsigned zerofill"},
		{ast.DataType{Name: "INT", Unsigned: true, Zerofill: true}, "int(11) unsigned zerofill"},
		{ast.DataType{Name: "DECIMAL", Zerofill: true}, "decimal(10,0) unsigned zerofill"},
		{ast.DataType{Name: "FLOAT", Length: 10}, "float"},
		{ast.DataType{Name: "FLOAT", Unsigned: true}, "float unsigned"},
		{ast.DataType{Name: "ENUM", EnumValues: []string{"with'quote"}}, "enum('with''quote')"},
		{ast.DataType{Name: "BOOLEAN"}, "tinyint(1)"},
		{ast.DataType{Name: "TIMESTAMP", Length: 3}, "timestamp(3)"},
	} {
		if got := columnType(&tt.dt); got != tt.want {
			t.Errorf("columnType(%+v) = %q, want %q", tt.dt, got, tt.want)
		}
	}
}

func TestNormalizeColumnType(t *testing.T) {
	for in, want := range map[string]string{
		"boolean": "tinyint(1)", "BOOL": "tinyint(1)", "decimal": "decimal(10,0)", "numeric": "decimal(10,0)",
		"fixed": "decimal(10,0)", "bit": "bit(1)", "binary": "binary(1)", "year": "year(4)",
		"tinyint": "tinyint(4)", "tinyint unsigned": "tinyint(4) unsigned", "smallint": "smallint(6)",
		"smallint unsigned": "smallint(6) unsigned", "mediumint": "mediumint(9)", "mediumint unsigned": "mediumint(9) unsigned",
		"int": "int(11)", "int unsigned": "int(11) unsigned", "bigint": "bigint(20)", "bigint unsigned": "bigint(20) unsigned",
		"VARCHAR(255)": "varchar(255)", "": "", "integer": "integer",
	} {
		if got := normalizeColumnType(in); got != want {
			t.Errorf("normalizeColumnType(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestColumnTypeTable(t *testing.T) {
	for _, tc := range []struct {
		name string
		dt   *ast.DataType
		want string
	}{
		{"BLOB", &ast.DataType{Name: "BLOB"}, "blob"},
		{"TINYBLOB", &ast.DataType{Name: "TINYBLOB"}, "tinyblob"},
		{"MEDIUMBLOB", &ast.DataType{Name: "MEDIUMBLOB"}, "mediumblob"},
		{"LONGBLOB", &ast.DataType{Name: "LONGBLOB"}, "longblob"},
		{"TEXT", &ast.DataType{Name: "TEXT"}, "text"},
		{"DECIMAL(10)", &ast.DataType{Name: "DECIMAL", Length: 10}, "decimal(10,0)"},
		{"DECIMAL(10,0)", &ast.DataType{Name: "DECIMAL", Length: 10, Scale: 0}, "decimal(10,0)"},
		{"DECIMAL(10,2)", &ast.DataType{Name: "DECIMAL", Length: 10, Scale: 2}, "decimal(10,2)"},
		{"NUMERIC(8)", &ast.DataType{Name: "NUMERIC", Length: 8}, "decimal(8,0)"},
		{"NUMERIC(8,2)", &ast.DataType{Name: "NUMERIC", Length: 8, Scale: 2}, "decimal(8,2)"},
		{"FIXED(8,2)-defensive", &ast.DataType{Name: "FIXED", Length: 8, Scale: 2}, "decimal(8,2)"},
		{"NUMERIC(10,2)-UNSIGNED", &ast.DataType{Name: "NUMERIC", Length: 10, Scale: 2, Unsigned: true}, "decimal(10,2) unsigned"},
		{"INT(11)-ZEROFILL", &ast.DataType{Name: "INT", Length: 11, Zerofill: true}, "int(11) unsigned zerofill"},
		{"FLOAT", &ast.DataType{Name: "FLOAT"}, "float"},
		{"FLOAT(10)", &ast.DataType{Name: "FLOAT", Length: 10}, "float"},
		{"FLOAT(10,0)", &ast.DataType{Name: "FLOAT", Length: 10, Scale: 0}, "float"},
		{"FLOAT(10,2)", &ast.DataType{Name: "FLOAT", Length: 10, Scale: 2}, "float(10,2)"},
		{"DOUBLE(10,4)", &ast.DataType{Name: "DOUBLE", Length: 10, Scale: 4}, "double(10,4)"},
		{"REAL", &ast.DataType{Name: "REAL"}, "real"},
		{"ENUM-2", &ast.DataType{Name: "ENUM", EnumValues: []string{"a", "b"}}, "enum('a','b')"},
		{"ENUM-3", &ast.DataType{Name: "ENUM", EnumValues: []string{"x", "y", "z"}}, "enum('x','y','z')"},
		{"SET-3", &ast.DataType{Name: "SET", EnumValues: []string{"a", "b", "c"}}, "set('a','b','c')"},
		{"ENUM-escape", &ast.DataType{Name: "ENUM", EnumValues: []string{"with'quote"}}, "enum('with''quote')"},
		{"BOOLEAN", &ast.DataType{Name: "BOOLEAN"}, "tinyint(1)"},
		{"DECIMAL-bare", &ast.DataType{Name: "DECIMAL"}, "decimal(10,0)"},
		{"NUMERIC-bare", &ast.DataType{Name: "NUMERIC"}, "decimal(10,0)"},
		{"BIT-bare", &ast.DataType{Name: "BIT"}, "bit(1)"},
		{"BINARY-bare", &ast.DataType{Name: "BINARY"}, "binary(1)"},
		{"YEAR-bare", &ast.DataType{Name: "YEAR"}, "year(4)"},
		{"INT-bare", &ast.DataType{Name: "INT"}, "int(11)"},
		{"TINYINT-bare", &ast.DataType{Name: "TINYINT"}, "tinyint(4)"},
		{"BIGINT-bare", &ast.DataType{Name: "BIGINT"}, "bigint(20)"},
		{"DECIMAL-bare-UNSIGNED", &ast.DataType{Name: "DECIMAL", Unsigned: true}, "decimal(10,0) unsigned"},
		{"DECIMAL-bare-ZEROFILL", &ast.DataType{Name: "DECIMAL", Zerofill: true}, "decimal(10,0) unsigned zerofill"},
		{"DECIMAL-bare-UNSIGNED-ZEROFILL", &ast.DataType{Name: "DECIMAL", Unsigned: true, Zerofill: true}, "decimal(10,0) unsigned zerofill"},
		{"NUMERIC-bare-UNSIGNED", &ast.DataType{Name: "NUMERIC", Unsigned: true}, "decimal(10,0) unsigned"},
		{"INT-bare-ZEROFILL", &ast.DataType{Name: "INT", Zerofill: true}, "int(11) unsigned zerofill"},
		{"INT-bare-UNSIGNED-ZEROFILL", &ast.DataType{Name: "INT", Unsigned: true, Zerofill: true}, "int(11) unsigned zerofill"},
		{"TINYINT-bare-ZEROFILL", &ast.DataType{Name: "TINYINT", Zerofill: true}, "tinyint(4) unsigned zerofill"},
		{"SMALLINT-bare-ZEROFILL", &ast.DataType{Name: "SMALLINT", Zerofill: true}, "smallint(6) unsigned zerofill"},
		{"MEDIUMINT-bare-ZEROFILL", &ast.DataType{Name: "MEDIUMINT", Zerofill: true}, "mediumint(9) unsigned zerofill"},
		{"BIGINT-bare-ZEROFILL", &ast.DataType{Name: "BIGINT", Zerofill: true}, "bigint(20) unsigned zerofill"},
		{"INT-bare-UNSIGNED", &ast.DataType{Name: "INT", Unsigned: true}, "int(11) unsigned"},
		{"BIGINT-bare-UNSIGNED", &ast.DataType{Name: "BIGINT", Unsigned: true}, "bigint(20) unsigned"},
		{"FLOAT-bare-UNSIGNED", &ast.DataType{Name: "FLOAT", Unsigned: true}, "float unsigned"},
		{"DOUBLE-bare-UNSIGNED", &ast.DataType{Name: "DOUBLE", Unsigned: true}, "double unsigned"},
		{"VARCHAR(255)", &ast.DataType{Name: "VARCHAR", Length: 255}, "varchar(255)"},
		{"CHAR(10)", &ast.DataType{Name: "CHAR", Length: 10}, "char(10)"},
		{"TIMESTAMP(3)", &ast.DataType{Name: "TIMESTAMP", Length: 3}, "timestamp(3)"},
		{"DATETIME", &ast.DataType{Name: "DATETIME"}, "datetime"},
		{"JSON", &ast.DataType{Name: "JSON"}, "json"},
		{"INT(11)-UNSIGNED", &ast.DataType{Name: "INT", Length: 11, Unsigned: true}, "int(11) unsigned"},
		{"nil", nil, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := columnType(tc.dt); got != tc.want {
				t.Errorf("columnType = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestNormalizeColumnTypeTable(t *testing.T) {
	for _, tc := range []struct {
		name, in, want string
	}{
		{"BOOLEAN", "boolean", "tinyint(1)"},
		{"BOOL", "bool", "tinyint(1)"},
		{"tinyint(1) identity", "tinyint(1)", "tinyint(1)"},
		{"bare decimal", "decimal", "decimal(10,0)"},
		{"bare numeric", "numeric", "decimal(10,0)"},
		{"bare fixed", "fixed", "decimal(10,0)"},
		{"bare bit", "bit", "bit(1)"},
		{"bare binary", "binary", "binary(1)"},
		{"bare year", "year", "year(4)"},
		{"bare tinyint", "tinyint", "tinyint(4)"},
		{"bare smallint", "smallint", "smallint(6)"},
		{"bare mediumint", "mediumint", "mediumint(9)"},
		{"bare int", "int", "int(11)"},
		{"bare bigint", "bigint", "bigint(20)"},
		{"int unsigned", "int unsigned", "int(11) unsigned"},
		{"bigint unsigned", "bigint unsigned", "bigint(20) unsigned"},
		{"varchar(255) identity", "varchar(255)", "varchar(255)"},
		{"already-canonical decimal", "decimal(10,2)", "decimal(10,2)"},
		{"empty string", "", ""},
		{"uppercase input", "VARCHAR(255)", "varchar(255)"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := normalizeColumnType(tc.in); got != tc.want {
				t.Errorf("normalizeColumnType(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

// TestColumnTypeMatchesStored holds normalizeColumnType(columnType(dt)) equal
// to the normalized type MySQL stores for dt exactly when the type is
// unchanged; a mismatch makes the rule report a no-op change.
func TestColumnTypeMatchesStored(t *testing.T) {
	for _, tc := range []struct {
		name      string
		dt        *ast.DataType
		stored    string
		unchanged bool
	}{
		{"BLOB", &ast.DataType{Name: "BLOB"}, "blob", true},
		{"TINYBLOB", &ast.DataType{Name: "TINYBLOB"}, "tinyblob", true},
		{"VARBINARY(100)", &ast.DataType{Name: "VARBINARY", Length: 100}, "varbinary(100)", true},
		{"DECIMAL(10)", &ast.DataType{Name: "DECIMAL", Length: 10}, "decimal(10,0)", true},
		{"DECIMAL(10,0)", &ast.DataType{Name: "DECIMAL", Length: 10, Scale: 0}, "decimal(10,0)", true},
		{"FLOAT", &ast.DataType{Name: "FLOAT"}, "float", true},
		{"FLOAT(10)", &ast.DataType{Name: "FLOAT", Length: 10}, "float", true},
		{"FLOAT(10,2)", &ast.DataType{Name: "FLOAT", Length: 10, Scale: 2}, "float(10,2)", true},
		{"DOUBLE(10,4)", &ast.DataType{Name: "DOUBLE", Length: 10, Scale: 4}, "double(10,4)", true},
		{"ENUM-2", &ast.DataType{Name: "ENUM", EnumValues: []string{"a", "b"}}, "enum('a','b')", true},
		{"SET-3", &ast.DataType{Name: "SET", EnumValues: []string{"x", "y", "z"}}, "set('x','y','z')", true},
		{"BOOLEAN", &ast.DataType{Name: "BOOLEAN"}, "tinyint(1)", true},
		{"BOOL", &ast.DataType{Name: "BOOLEAN"}, "tinyint(1)", true},
		{"bare DECIMAL", &ast.DataType{Name: "DECIMAL"}, "decimal(10,0)", true},
		{"bare NUMERIC", &ast.DataType{Name: "NUMERIC"}, "decimal(10,0)", true},
		{"bare BIT", &ast.DataType{Name: "BIT"}, "bit(1)", true},
		{"bare BINARY", &ast.DataType{Name: "BINARY"}, "binary(1)", true},
		{"bare YEAR", &ast.DataType{Name: "YEAR"}, "year(4)", true},
		{"bare INT", &ast.DataType{Name: "INT"}, "int(11)", true},
		{"INT(11)", &ast.DataType{Name: "INT", Length: 11}, "int(11)", true},
		{"bare TINYINT", &ast.DataType{Name: "TINYINT"}, "tinyint(4)", true},
		{"VARCHAR(255)", &ast.DataType{Name: "VARCHAR", Length: 255}, "varchar(255)", true},
		{"DECIMAL UNSIGNED no-op", &ast.DataType{Name: "DECIMAL", Unsigned: true}, "decimal(10,0) unsigned", true},
		{"NUMERIC(8,2) no-op (catalog decimal)", &ast.DataType{Name: "NUMERIC", Length: 8, Scale: 2}, "decimal(8,2)", true},
		{"NUMERIC bare no-op (catalog decimal)", &ast.DataType{Name: "NUMERIC"}, "decimal(10,0)", true},
		{"DECIMAL ZEROFILL no-op", &ast.DataType{Name: "DECIMAL", Zerofill: true}, "decimal(10,0) unsigned zerofill", true},
		{"INT ZEROFILL no-op", &ast.DataType{Name: "INT", Zerofill: true}, "int(11) unsigned zerofill", true},
		{"TINYINT ZEROFILL no-op", &ast.DataType{Name: "TINYINT", Zerofill: true}, "tinyint(4) unsigned zerofill", true},
		{"INT UNSIGNED no-op", &ast.DataType{Name: "INT", Unsigned: true}, "int(11) unsigned", true},
		{"BIGINT UNSIGNED no-op", &ast.DataType{Name: "BIGINT", Unsigned: true}, "bigint(20) unsigned", true},
		{"INT → BIGINT", &ast.DataType{Name: "BIGINT"}, "int(11)", false},
		{"VARCHAR(50) → VARCHAR(255)", &ast.DataType{Name: "VARCHAR", Length: 255}, "varchar(50)", false},
		{"DECIMAL(10,2) → DECIMAL(10,4)", &ast.DataType{Name: "DECIMAL", Length: 10, Scale: 4}, "decimal(10,2)", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			built := columnType(tc.dt)
			if matched := normalizeColumnType(built) == normalizeColumnType(tc.stored); matched != tc.unchanged {
				t.Errorf("columnType = %q against stored %q: match %v, want %v", built, tc.stored, matched, tc.unchanged)
			}
		})
	}
}
