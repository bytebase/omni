package v1review

import (
	"testing"

	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

func TestColumnTypeDisallowList(t *testing.T) {
	bad := func(table, name, tp, text string) placed {
		return placed{Code: v1review.ColumnTypeDisallowed, Name: name, Table: table, Value: tp, Text: text}
	}
	disallowed := []string{"json", "INT(11)", "varchar(255)", "ENUM('A','B''C')", "TINYINT(1)", "DECIMAL(8,2)", "FLOAT"}
	for _, tt := range []struct {
		name string
		sql  string
		want []placed
	}{
		{name: "allowed", sql: "CREATE TABLE t (a INT(10), b VARCHAR(10), c TEXT, d FLOAT(10,2), e BIGINT);"},
		{
			name: "rendered forms",
			sql:  "CREATE TABLE t (a JSON, b INT UNSIGNED, c VARCHAR(255) CHARSET utf8mb4, d ENUM('a','b''c'), e BOOLEAN, f NUMERIC(8,2), g FLOAT(10), h INTEGER);",
			want: []placed{
				bad("t", "a", "JSON", "a JSON"),
				bad("t", "b", "INT(11)", "b INT UNSIGNED"),
				bad("t", "c", "VARCHAR(255)", "c VARCHAR(255) CHARSET utf8mb4"),
				bad("t", "d", "ENUM('A','B''C')", "d ENUM('a','b''c')"),
				bad("t", "e", "TINYINT(1)", "e BOOLEAN"),
				bad("t", "f", "DECIMAL(8,2)", "f NUMERIC(8,2)"),
				bad("t", "g", "FLOAT", "g FLOAT(10)"),
				bad("t", "h", "INT(11)", "h INTEGER"),
			},
		},
		{
			name: "alter anchors the statement",
			sql:  "ALTER TABLE t ADD a JSON, MODIFY b INT;",
			want: []placed{
				bad("t", "a", "JSON", "ALTER TABLE t ADD a JSON, MODIFY b INT"),
				bad("t", "b", "INT(11)", "ALTER TABLE t ADD a JSON, MODIFY b INT"),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			checkPlaced(t, tt.sql, func(stmts []ast.Node) []v1review.Finding {
				return ColumnTypeDisallowList(stmts, disallowed)
			}, tt.want)
		})
	}
}

func TestCompactType(t *testing.T) {
	for _, tt := range []struct {
		dt   ast.DataType
		want string
	}{
		{ast.DataType{Name: "DECIMAL"}, "decimal(10,0)"},
		{ast.DataType{Name: "FIXED", Length: 5}, "decimal(5,0)"},
		{ast.DataType{Name: "DOUBLE", Length: 10, Scale: 2}, "double(10,2)"},
		{ast.DataType{Name: "REAL", Length: 10}, "real"},
		{ast.DataType{Name: "TINYINT"}, "tinyint(4)"},
		{ast.DataType{Name: "SMALLINT"}, "smallint(6)"},
		{ast.DataType{Name: "MEDIUMINT"}, "mediumint(9)"},
		{ast.DataType{Name: "BIGINT"}, "bigint(20)"},
		{ast.DataType{Name: "BIT"}, "bit(1)"},
		{ast.DataType{Name: "BINARY"}, "binary(1)"},
		{ast.DataType{Name: "YEAR"}, "year(4)"},
		{ast.DataType{Name: "BOOL"}, "tinyint(1)"},
		{ast.DataType{Name: "SET"}, "set()"},
		{ast.DataType{Name: "DATETIME", Length: 3}, "datetime(3)"},
	} {
		if got := compactType(&tt.dt); got != tt.want {
			t.Errorf("compactType(%+v) = %q, want %q", tt.dt, got, tt.want)
		}
	}
}

func TestCompactTypeTable(t *testing.T) {
	for _, tc := range []struct {
		name string
		dt   *ast.DataType
		want string
	}{
		{"JSON", &ast.DataType{Name: "JSON"}, "json"},
		{"BLOB", &ast.DataType{Name: "BLOB"}, "blob"},
		{"DATE", &ast.DataType{Name: "DATE"}, "date"},
		{"VARCHAR(255)", &ast.DataType{Name: "VARCHAR", Length: 255}, "varchar(255)"},
		{"CHAR(10)", &ast.DataType{Name: "CHAR", Length: 10}, "char(10)"},
		{"BINARY(16)", &ast.DataType{Name: "BINARY", Length: 16}, "binary(16)"},
		{"INT bare → int(11)", &ast.DataType{Name: "INT"}, "int(11)"},
		{"TINYINT bare → tinyint(4)", &ast.DataType{Name: "TINYINT"}, "tinyint(4)"},
		{"SMALLINT bare → smallint(6)", &ast.DataType{Name: "SMALLINT"}, "smallint(6)"},
		{"MEDIUMINT bare → mediumint(9)", &ast.DataType{Name: "MEDIUMINT"}, "mediumint(9)"},
		{"BIGINT bare → bigint(20)", &ast.DataType{Name: "BIGINT"}, "bigint(20)"},
		{"BIT bare → bit(1)", &ast.DataType{Name: "BIT"}, "bit(1)"},
		{"BINARY bare → binary(1)", &ast.DataType{Name: "BINARY"}, "binary(1)"},
		{"YEAR bare → year(4)", &ast.DataType{Name: "YEAR"}, "year(4)"},
		{"TINYINT(1)", &ast.DataType{Name: "TINYINT", Length: 1}, "tinyint(1)"},
		{"INT(11)", &ast.DataType{Name: "INT", Length: 11}, "int(11)"},
		{"DECIMAL bare → decimal(10,0)", &ast.DataType{Name: "DECIMAL"}, "decimal(10,0)"},
		{"DECIMAL(10)", &ast.DataType{Name: "DECIMAL", Length: 10}, "decimal(10,0)"},
		{"DECIMAL(10,2)", &ast.DataType{Name: "DECIMAL", Length: 10, Scale: 2}, "decimal(10,2)"},
		{"NUMERIC(10,2) → decimal", &ast.DataType{Name: "NUMERIC", Length: 10, Scale: 2}, "decimal(10,2)"},
		{"FLOAT bare", &ast.DataType{Name: "FLOAT"}, "float"},
		{"FLOAT(10) → float", &ast.DataType{Name: "FLOAT", Length: 10}, "float"},
		{"FLOAT(10,2)", &ast.DataType{Name: "FLOAT", Length: 10, Scale: 2}, "float(10,2)"},
		{"ENUM('x','y')", &ast.DataType{Name: "ENUM", EnumValues: []string{"x", "y"}}, "enum('x','y')"},
		{"SET('a','b')", &ast.DataType{Name: "SET", EnumValues: []string{"a", "b"}}, "set('a','b')"},
		{"BOOLEAN → tinyint(1)", &ast.DataType{Name: "BOOLEAN"}, "tinyint(1)"},
		{"BIGINT UNSIGNED → compact has no UNSIGNED", &ast.DataType{Name: "BIGINT", Unsigned: true}, "bigint(20)"},
		{"INT(11) UNSIGNED → compact has no UNSIGNED", &ast.DataType{Name: "INT", Length: 11, Unsigned: true}, "int(11)"},
		{"INT ZEROFILL → compact has no ZEROFILL", &ast.DataType{Name: "INT", Zerofill: true}, "int(11)"},
		{"nil", nil, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := compactType(tc.dt); got != tc.want {
				t.Errorf("compactType = %q, want %q", got, tc.want)
			}
		})
	}
}
