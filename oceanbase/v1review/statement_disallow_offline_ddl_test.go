package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/oceanbase/parser"
	"github.com/bytebase/omni/v1review"
)

func TestStatementDisallowOfflineDDL(t *testing.T) {
	sql := "ALTER TABLE t ADD COLUMN a INT AUTO_INCREMENT, DROP COLUMN g, DROP COLUMN v, ADD PRIMARY KEY (a); DROP TABLE t; TRUNCATE TABLE u; ALTER TABLE db2.t DROP COLUMN g; ALTER TABLE t COMMENT = 'x'"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	g := v1review.ColumnRef{Database: "db", Table: "t", Column: "g"}
	v := v1review.ColumnRef{Database: "db", Table: "t", Column: "v"}
	g2 := v1review.ColumnRef{Database: "db2", Table: "t", Column: "g"}
	if diff := cmp.Diff([]v1review.ColumnRef{g, v, g2}, StatementDisallowOfflineDDLLookups(list.Items, "db")); diff != "" {
		t.Errorf("lookups (-want +got):\n%s", diff)
	}
	if got := StatementDisallowOfflineDDLLookups(list.Items, ""); len(got) != 1 || got[0] != g2 {
		t.Errorf("no current database: got %v, want only the qualified table's column", got)
	}
	type finding struct{ Name, Text string }
	var got []finding
	for _, f := range StatementDisallowOfflineDDL(list.Items, "db", map[v1review.ColumnRef]string{g: "STORED GENERATED", v: "VIRTUAL GENERATED"}) {
		got = append(got, finding{f.Name, sql[f.Start:f.End]})
	}
	want := []finding{
		{"Adding auto increment columns", "ADD COLUMN a INT AUTO_INCREMENT"},
		{"Dropping stored generated columns", "DROP COLUMN g"},
		{"Adding primary keys", "ADD PRIMARY KEY (a)"},
		{"Dropping tables", "DROP TABLE t"},
		{"Truncating tables", "TRUNCATE TABLE u"},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("findings (-want +got):\n%s", diff)
	}
}
