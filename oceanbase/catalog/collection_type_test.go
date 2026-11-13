package catalog

import (
	"strings"
	"testing"
)

// TestCollectionColumnTypes holds ARRAY and MAP columns to the text OceanBase prints
// for them in SHOW CREATE TABLE (simulate.tsv's tape), elements included, so Diff
// sees a change of element type and the text reads back.
func TestCollectionColumnTypes(t *testing.T) {
	build := func(column string) *Catalog {
		t.Helper()
		c := New()
		mustExec(t, c, "CREATE DATABASE test")
		c.SetCurrentDatabase("test")
		mustExec(t, c, "CREATE TABLE t (id INT PRIMARY KEY, a "+column+")")
		return c
	}
	for column, want := range map[string]string{
		"ARRAY(TINYINT)":        "ARRAY(TINYINT)",
		"ARRAY<ARRAY<TINYINT>>": "ARRAY(ARRAY(TINYINT))",
		"TINYINT[]":             "ARRAY(TINYINT)",
		"MAP(TINYINT, TINYINT)": "MAP(TINYINT,TINYINT)",
		"SPARSEVECTOR":          "SPARSEVECTOR",
	} {
		c := build(column)
		show := c.ShowCreateTable("test", "t")
		if !strings.Contains(show, "`a` "+want) {
			t.Errorf("%s: SHOW CREATE TABLE =\n%s\nwant `a` %s", column, show, want)
		}
		back := New()
		mustExec(t, back, "CREATE DATABASE test")
		back.SetCurrentDatabase("test")
		if _, err := back.Exec(show, nil); err != nil {
			t.Errorf("%s: SHOW CREATE TABLE does not read back: %v", column, err)
		}
	}
	if Diff(build("ARRAY(INT)"), build("ARRAY(VARCHAR(3))")).IsEmpty() {
		t.Error("Diff saw no change between ARRAY(INT) and ARRAY(VARCHAR(3))")
	}
}
