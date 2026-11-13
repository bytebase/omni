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
				{Name: "a", Type: "int(11)"}, {Name: "b", Type: "BLOB"}, {Name: "c", Type: "bigint"},
				{Name: "d", Type: "VARCHAR(20)"}, {Name: "e", Type: "blob(100)"},
			},
		}},
	}}}}
}

// statementFinding is a finding of a rule over a whole script.
type statementFinding struct {
	Code              v1review.Code
	Table, Name, Type string
	Statement         int
	Text              string
}

func scriptFindings(sqls []string, got []v1review.Finding) []statementFinding {
	var findings []statementFinding
	for _, f := range got {
		findings = append(findings, statementFinding{f.Code, f.Table, f.Name, f.Type, f.Statement, sqls[f.Statement][f.Start:f.End]})
	}
	return findings
}

func TestIndexPKTypeLimit(t *testing.T) {
	bad := func(table, column, columnType string, statement int, text string) statementFinding {
		return statementFinding{v1review.PrimaryKeyTypeNotAllowed, table, column, columnType, statement, text}
	}
	for _, tt := range []struct {
		name string
		sqls []string
		want []statementFinding
	}{
		{
			name: "column attributes, with length and scale",
			sqls: []string{"CREATE TABLE n (a INT PRIMARY KEY, b INTEGER PRIMARY KEY, c INT UNSIGNED PRIMARY KEY, d DECIMAL(10,2) PRIMARY KEY, e VARCHAR(20) PRIMARY KEY, f TEXT PRIMARY KEY)"},
			want: []statementFinding{
				bad("n", "d", "decimal(10,2)", 0, "d DECIMAL(10,2) PRIMARY KEY"),
				bad("n", "e", "varchar(20)", 0, "e VARCHAR(20) PRIMARY KEY"),
				bad("n", "f", "text", 0, "f TEXT PRIMARY KEY"),
			},
		},
		{
			name: "a listed column takes the script's type, else the snapshot's as spelled",
			sqls: []string{"ALTER TABLE t ADD PRIMARY KEY (a, b, c, missing)", "ALTER TABLE t MODIFY COLUMN a BIGINT", "ALTER TABLE t ADD PRIMARY KEY (a)"},
			want: []statementFinding{
				bad("t", "a", "int(11)", 0, "ALTER TABLE t ADD PRIMARY KEY (a, b, c, missing)"),
				bad("t", "b", "BLOB", 0, "ALTER TABLE t ADD PRIMARY KEY (a, b, c, missing)"),
			},
		},
		{
			name: "change column forgets the old name",
			sqls: []string{"CREATE TABLE n (a VARCHAR(5))", "ALTER TABLE n CHANGE COLUMN a b INT", "ALTER TABLE n ADD PRIMARY KEY (a, b)"},
		},
		{
			name: "the script's names are matched as written",
			sqls: []string{"CREATE TABLE n (a VARCHAR(5))", "ALTER TABLE N ADD PRIMARY KEY (a)", "ALTER TABLE n ADD PRIMARY KEY (A)"},
		},
		{
			name: "letter case does not matter in the snapshot",
			sqls: []string{"ALTER TABLE T ADD PRIMARY KEY (D)"},
			want: []statementFinding{bad("T", "D", "VARCHAR(20)", 0, "ALTER TABLE T ADD PRIMARY KEY (D)")},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IndexPKTypeLimit(parseEach(t, tt.sqls), keyColumnSchema())
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.want, scriptFindings(tt.sqls, got)); diff != "" {
				t.Errorf("findings (-want +got):\n%s", diff)
			}
		})
	}
}

func TestIndexPrimaryKeyTypeAllowlist(t *testing.T) {
	sqls := []string{
		"CREATE TABLE n (a VARCHAR(5) PRIMARY KEY, b INT PRIMARY KEY)",
		"ALTER TABLE t ADD PRIMARY KEY (a, c, d)",
	}
	got, err := IndexPrimaryKeyTypeAllowlist(parseEach(t, sqls), []string{"INT", "BIGINT", "Varchar"}, keyColumnSchema())
	if err != nil {
		t.Fatal(err)
	}
	want := []statementFinding{
		{v1review.PrimaryKeyTypeNotAllowed, "t", "a", "int(11)", 1, "ALTER TABLE t ADD PRIMARY KEY (a, c, d)"},
		{v1review.PrimaryKeyTypeNotAllowed, "t", "d", "varchar(20)", 1, "ALTER TABLE t ADD PRIMARY KEY (a, c, d)"},
	}
	if diff := cmp.Diff(want, scriptFindings(sqls, got)); diff != "" {
		t.Errorf("findings (-want +got):\n%s", diff)
	}
}

func TestIndexTypeNoBlob(t *testing.T) {
	sqls := []string{
		"CREATE TABLE n (a BLOB PRIMARY KEY, b TINYBLOB UNIQUE, c TEXT UNIQUE, d LONGBLOB, KEY (d), FULLTEXT KEY (d))",
		"CREATE INDEX i ON t (b, e, a)",
		"CREATE FULLTEXT INDEX f ON t (b)",
	}
	got, err := IndexTypeNoBlob(parseEach(t, sqls), keyColumnSchema())
	if err != nil {
		t.Fatal(err)
	}
	want := []statementFinding{
		{v1review.IndexColumnTypeNotAllowed, "n", "a", "blob", 0, "a BLOB PRIMARY KEY"},
		{v1review.IndexColumnTypeNotAllowed, "n", "b", "tinyblob", 0, "b TINYBLOB UNIQUE"},
		{v1review.IndexColumnTypeNotAllowed, "n", "d", "longblob", 0, "KEY (d)"},
		{v1review.IndexColumnTypeNotAllowed, "t", "b", "blob", 1, "CREATE INDEX i ON t (b, e, a)"},
	}
	if diff := cmp.Diff(want, scriptFindings(sqls, got)); diff != "" {
		t.Errorf("findings (-want +got):\n%s", diff)
	}
}

func TestKeyColumnsSchemaMissing(t *testing.T) {
	defined := parseEach(t, []string{"CREATE TABLE n (a BLOB)", "ALTER TABLE n ADD PRIMARY KEY (a)"})
	if _, err := IndexTypeNoBlob(defined, v1review.Schema{}); err != nil {
		t.Errorf("a script that defines every listed column read the schema: %v", err)
	}
	undefined := parseEach(t, []string{"CREATE TABLE n (a BLOB PRIMARY KEY)", "ALTER TABLE n ADD PRIMARY KEY (b)"})
	got, err := IndexPKTypeLimit(undefined, v1review.Schema{})
	if !errors.Is(err, v1review.ErrSchemaMissing) || got != nil {
		t.Errorf("got %v, %v; want no findings and ErrSchemaMissing", got, err)
	}
}
