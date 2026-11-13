package v1review

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/v1review"
)

// keyColumnSchema holds a table t whose columns have the types the column
// type rules branch on.
func keyColumnSchema() v1review.Schema {
	return v1review.Schema{Metadata: &metadata.DatabaseSchemaMetadata{Schemas: []*metadata.SchemaMetadata{{
		Tables: []*metadata.TableMetadata{{
			Name: "t",
			Columns: []*metadata.ColumnMetadata{
				{Name: "a", Type: "int(11)"}, {Name: "b", Type: "BLOB"}, {Name: "c", Type: "BIGINT"},
				{Name: "d", Type: "VARCHAR(20)"}, {Name: "e", Type: "blob(100)"}, {Name: "f", Type: "int"},
				{Name: "g", Type: "MediumText"},
			},
		}},
	}}}}
}

func TestIndexPKTypeLimit(t *testing.T) {
	const alter = "ALTER TABLE t ADD COLUMN x CHAR(2) PRIMARY KEY, MODIFY y CHAR(3) PRIMARY KEY, CHANGE COLUMN z w CHAR(4) PRIMARY KEY"
	bad := func(table, column, columnType string, statement int, text string) indexFinding {
		return indexFinding{Code: v1review.PrimaryKeyTypeNotAllowed, Table: table, Name: column, Type: columnType, Statement: statement, Text: text}
	}
	for _, tt := range []struct {
		name string
		sql  string
		want []indexFinding
	}{
		{
			name: "column attributes, with length and scale",
			sql:  "CREATE TABLE n (a INT PRIMARY KEY, b bigint PRIMARY KEY, c INT UNSIGNED PRIMARY KEY, d DECIMAL(10,2) PRIMARY KEY, e VARCHAR(20) PRIMARY KEY, f TEXT PRIMARY KEY)",
			want: []indexFinding{
				bad("n", "d", "decimal(10,2)", 0, "d DECIMAL(10,2) PRIMARY KEY"),
				bad("n", "e", "varchar(20)", 0, "e VARCHAR(20) PRIMARY KEY"),
				bad("n", "f", "text", 0, "f TEXT PRIMARY KEY"),
			},
		},
		{
			name: "a table-level primary key lists the script's columns",
			sql:  "CREATE TABLE n (a VARCHAR(5), b INT, PRIMARY KEY (a, b), UNIQUE (a))",
			want: []indexFinding{bad("n", "a", "varchar(5)", 0, "PRIMARY KEY (a, b)")},
		},
		{
			name: "a listed column takes the script's type, else the snapshot's exactly as spelled",
			sql:  "ALTER TABLE t ADD PRIMARY KEY (a, b, c, f, missing);\nALTER TABLE t MODIFY COLUMN a BIGINT;\nALTER TABLE t ADD PRIMARY KEY (a)",
			want: []indexFinding{
				bad("t", "a", "int(11)", 0, "ALTER TABLE t ADD PRIMARY KEY (a, b, c, f, missing)"),
				bad("t", "b", "BLOB", 0, "ALTER TABLE t ADD PRIMARY KEY (a, b, c, f, missing)"),
				bad("t", "f", "int", 0, "ALTER TABLE t ADD PRIMARY KEY (a, b, c, f, missing)"),
			},
		},
		{
			name: "alter table findings anchor the statement",
			sql:  alter,
			want: []indexFinding{
				bad("t", "x", "char(2)", 0, alter),
				bad("t", "y", "char(3)", 0, alter),
				bad("t", "w", "char(4)", 0, alter),
			},
		},
		{
			name: "change column forgets the old name",
			sql:  "CREATE TABLE n (a VARCHAR(5));\nALTER TABLE n CHANGE COLUMN a b INT;\nALTER TABLE n ADD PRIMARY KEY (a, b)",
		},
		{
			name: "the script's names are matched as written, the snapshot's without case",
			sql:  "CREATE TABLE T (a VARCHAR(5));\nALTER TABLE t ADD PRIMARY KEY (D)",
			want: []indexFinding{bad("t", "D", "VARCHAR(20)", 1, "ALTER TABLE t ADD PRIMARY KEY (D)")},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IndexPKTypeLimit(parseScript(t, tt.sql), keyColumnSchema())
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.want, indexFindings(tt.sql, got)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}

func TestIndexPrimaryKeyTypeAllowlist(t *testing.T) {
	sql := "CREATE TABLE n (a VARCHAR(5) PRIMARY KEY, b INT PRIMARY KEY, c DECIMAL(4,1) PRIMARY KEY);\nALTER TABLE t ADD PRIMARY KEY (a, c, d)"
	got, err := IndexPrimaryKeyTypeAllowlist(parseScript(t, sql), []string{"INT", "BIGINT", "Varchar"}, keyColumnSchema())
	if err != nil {
		t.Fatal(err)
	}
	const alter = "ALTER TABLE t ADD PRIMARY KEY (a, c, d)"
	want := []indexFinding{
		{Code: v1review.PrimaryKeyTypeNotAllowed, Table: "n", Name: "c", Type: "decimal", Text: "c DECIMAL(4,1) PRIMARY KEY"},
		{Code: v1review.PrimaryKeyTypeNotAllowed, Table: "t", Name: "a", Type: "int(11)", Statement: 1, Text: alter},
		{Code: v1review.PrimaryKeyTypeNotAllowed, Table: "t", Name: "d", Type: "varchar(20)", Statement: 1, Text: alter},
	}
	if diff := cmp.Diff(want, indexFindings(sql, got)); diff != "" {
		t.Errorf("findings (-want +got):\n%s", diff)
	}
}

func TestIndexTypeNoBlob(t *testing.T) {
	sql := "CREATE TABLE n (a BLOB PRIMARY KEY, b TINYBLOB UNIQUE, c TEXT UNIQUE, d LONGBLOB, KEY (d), FULLTEXT KEY (c), FOREIGN KEY (d) REFERENCES p (d));\n" +
		"ALTER TABLE n MODIFY a BLOB PRIMARY KEY;\n" +
		"CREATE INDEX i ON t (b, e, g, a);\n" +
		"CREATE FULLTEXT INDEX f ON t (b)"
	got, err := IndexTypeNoBlob(parseScript(t, sql), keyColumnSchema())
	if err != nil {
		t.Fatal(err)
	}
	notAllowed := func(table, column, columnType string, statement int, text string) indexFinding {
		return indexFinding{Code: v1review.IndexColumnTypeNotAllowed, Table: table, Name: column, Type: columnType, Statement: statement, Text: text}
	}
	want := []indexFinding{
		notAllowed("n", "b", "tinyblob", 0, "b TINYBLOB UNIQUE"),
		notAllowed("n", "c", "text", 0, "c TEXT UNIQUE"),
		notAllowed("n", "d", "longblob", 0, "KEY (d)"),
		notAllowed("n", "a", "blob", 1, "ALTER TABLE n MODIFY a BLOB PRIMARY KEY"),
		notAllowed("t", "b", "blob", 2, "CREATE INDEX i ON t (b, e, g, a)"),
		notAllowed("t", "g", "mediumtext", 2, "CREATE INDEX i ON t (b, e, g, a)"),
		notAllowed("t", "b", "blob", 3, "CREATE FULLTEXT INDEX f ON t (b)"),
	}
	if diff := cmp.Diff(want, indexFindings(sql, got)); diff != "" {
		t.Errorf("findings (-want +got):\n%s", diff)
	}
}

func TestKeyColumnsSchemaMissing(t *testing.T) {
	defined := parseScript(t, "CREATE TABLE n (a BLOB);\nALTER TABLE n ADD PRIMARY KEY (a)")
	if _, err := IndexTypeNoBlob(defined, v1review.Schema{}); err != nil {
		t.Errorf("a script that defines every listed column read the schema: %v", err)
	}
	undefined := parseScript(t, "CREATE TABLE n (a BLOB PRIMARY KEY);\nALTER TABLE n ADD PRIMARY KEY (b)")
	got, err := IndexPKTypeLimit(undefined, v1review.Schema{})
	if !errors.Is(err, v1review.ErrSchemaMissing) || got != nil {
		t.Errorf("got %v, %v; want no findings and ErrSchemaMissing", got, err)
	}
}
