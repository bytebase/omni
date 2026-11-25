package parsertest

import (
	"reflect"
	"testing"
)

func TestExtractSQLBlocks(t *testing.T) {
	md := "# Title\n" +
		"\n" +
		"## DDL-001: CREATE DATABASE\n" +
		"\n" +
		"```\n" +
		"CREATE DATABASE database_id  -- grammar block, NOT sql-fenced; must be ignored\n" +
		"```\n" +
		"\n" +
		"```sql\n" +
		"CREATE DATABASE my_database;\n" +
		"```\n" +
		"\n" +
		"## DDL-002: indented body\n" +
		"\n" +
		"```sql\n" +
		"SELECT\n" +
		"  ```not a fence — trimmed line is not exactly three backticks\n" +
		"FROM t;\n" +
		"```\n" +
		"\n" +
		"## Note: a prose heading with a colon but no ID token\n" +
		"\n" +
		"```sql\n" +
		"SELECT 1;\n" +
		"```\n"

	got := ExtractSQLBlocks("x/y.md", md)

	// Only the three ```sql blocks are returned; the bare ``` grammar block is
	// skipped. Indices are document-order over the sql-fenced blocks only.
	want := []SQLBlock{
		{File: "x/y.md", Index: 0, SectionID: "DDL-001", SectionTitle: "CREATE DATABASE", Text: "CREATE DATABASE my_database;"},
		{File: "x/y.md", Index: 1, SectionID: "DDL-002", SectionTitle: "indented body", Text: "SELECT\n  ```not a fence — trimmed line is not exactly three backticks\nFROM t;"},
		{File: "x/y.md", Index: 2, SectionID: "", SectionTitle: "Note: a prose heading with a colon but no ID token", Text: "SELECT 1;"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ExtractSQLBlocks mismatch:\n got=%#v\nwant=%#v", got, want)
	}
}

func TestSQLBlockKeyAndLabel(t *testing.T) {
	b := SQLBlock{File: "bigquery/ddl.md", Index: 13, SectionID: "DDL-013", SectionTitle: "CREATE EXTERNAL SCHEMA"}
	if got, want := b.Key(), "bigquery/ddl.md#13"; got != want {
		t.Errorf("Key() = %q, want %q", got, want)
	}
	if got, want := b.Label(), "bigquery/ddl.md#13 (DDL-013: CREATE EXTERNAL SCHEMA)"; got != want {
		t.Errorf("Label() = %q, want %q", got, want)
	}
	bare := SQLBlock{File: "a.md", Index: 0}
	if got, want := bare.Label(), "a.md#0"; got != want {
		t.Errorf("bare Label() = %q, want %q", got, want)
	}
}

func TestItoa(t *testing.T) {
	for _, c := range []struct {
		n    int
		want string
	}{{0, "0"}, {7, "7"}, {13, "13"}, {350, "350"}, {-4, "-4"}} {
		if got := itoa(c.n); got != c.want {
			t.Errorf("itoa(%d) = %q, want %q", c.n, got, c.want)
		}
	}
}
