package catalog

import (
	"errors"
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

	// Two spellings of one element type are one type.
	if !Diff(build("ARRAY(INTEGER)"), build("ARRAY(INT)")).IsEmpty() {
		t.Error("Diff saw a change between ARRAY(INTEGER) and ARRAY(INT)")
	}
	if !Diff(build("MAP(INT, NUMERIC(10,2))"), build("MAP(INT, DECIMAL(10,2))")).IsEmpty() {
		t.Error("Diff saw a change between MAP(INT, NUMERIC(10,2)) and MAP(INT, DECIMAL(10,2))")
	}

	// An element OceanBase does not keep is its 1235, and the table is left as it was.
	c := build("INT")
	for _, sql := range []string{"ALTER TABLE t ADD b TINYTEXT[]", "ALTER TABLE t ADD b ARRAY(ARRAY(DECIMAL(5,2)))", "CREATE TABLE u (a ARRAY(ENUM('x')))"} {
		results, _ := c.Exec(sql, &ExecOptions{ContinueOnError: true})
		var ce *Error
		if len(results) != 1 || !errors.As(results[0].Error, &ce) || ce.Code != 1235 {
			t.Errorf("%s: got %v, want 1235", sql, results[0].Error)
		}
	}
	if c.GetDatabase("test").GetTable("t").GetColumn("b") != nil {
		t.Error("a refused ADD left its column")
	}

	// A routine's ARRAY and MAP parameters keep their elements, so it reads back.
	mustExec(t, c, "CREATE PROCEDURE p(a ARRAY(INT), m MAP(INT, VARCHAR(3))) SELECT 1")
	show := c.ShowCreateProcedure("test", "p")
	if !strings.Contains(show, "ARRAY(INT)") || !strings.Contains(show, "MAP(INT,VARCHAR(3))") {
		t.Errorf("SHOW CREATE PROCEDURE lost a collection's elements:\n%s", show)
	}
	back := New()
	mustExec(t, back, "CREATE DATABASE test")
	back.SetCurrentDatabase("test")
	if _, err := back.Exec(show, nil); err != nil {
		t.Errorf("SHOW CREATE PROCEDURE does not read back: %v\n%s", err, show)
	}
}
