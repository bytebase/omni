package catalog

import (
	"errors"
	"strings"
	"testing"
)

// TestColumnGroupRefusals pins each refusal's SQLSTATE and text, which
// TestOracle_OceanBaseColumnGroups does not compare: OceanBase answers 1235 0A000,
// 1060 42S21 and 4185 HY000 (its tape).
func TestColumnGroupRefusals(t *testing.T) {
	for _, tc := range []struct {
		sql   string
		code  int
		state string
		text  string
	}{
		{"CREATE TABLE t2 (a INT) WITH COLUMN GROUP(g(a))", 1235, "0A000", "Not supported feature or function"},
		{"ALTER TABLE t ADD COLUMN GROUP(ALL COLUMNS) DELAYED", 1235, "0A000", "Not supported feature or function"},
		{"ALTER TABLE t ADD COLUMN GROUP(ALL COLUMNS, ALL COLUMNS)", 1060, "42S21", "Duplicate column group name 'all column group'"},
		{"ALTER TABLE t ADD COLUMN GROUP(EACH COLUMN, EACH COLUMN)", 1060, "42S21", "Duplicate column group name 'single column group'"},
		{"ALTER TABLE t ADD COLUMN GROUP(EACH COLUMN)", 1060, "42S21", "Duplicate column group name 'each column group'"},
		{"ALTER TABLE t DROP COLUMN GROUP(ALL COLUMNS)", 4185, "HY000", "Column group 'all column group' not found"},
	} {
		c := New()
		if _, err := c.Exec("CREATE DATABASE d; USE d; CREATE TABLE t (a INT) WITH COLUMN GROUP(each column)", nil); err != nil {
			t.Fatal(err)
		}
		results, err := c.Exec(tc.sql, nil)
		if err != nil {
			t.Fatal(err)
		}
		var got *Error
		if len(results) != 1 || !errors.As(results[0].Error, &got) {
			t.Errorf("%s: no refusal: %+v", tc.sql, results)
			continue
		}
		if got.Code != tc.code || got.SQLState != tc.state || got.Message != tc.text {
			t.Errorf("%s: %v, want ERROR %d (%s): %s", tc.sql, got, tc.code, tc.state, tc.text)
		}
	}
}

// TestColumnGroupsReadAsWritten pins that the list is read with the parser's lexer: a
// comment or a quoted name inside it is what the parser saw.
func TestColumnGroupsReadAsWritten(t *testing.T) {
	for _, tc := range []struct {
		list string
		want ColumnGroups
		code int
	}{
		{"all columns /* , g(a) */, each column", ColumnGroups{AllColumns: true, EachColumn: true}, 0},
		{"EACH COLUMN -- , each column\n", ColumnGroups{EachColumn: true}, 0},
		{"`all`(a)", ColumnGroups{}, 1235},
		{"each column, `each`(a, b)", ColumnGroups{}, 1235},
	} {
		got, err := readColumnGroups(tc.list)
		var refusal *Error
		switch {
		case tc.code != 0 && (!errors.As(err, &refusal) || refusal.Code != tc.code):
			t.Errorf("%q: %v, want %d", tc.list, err, tc.code)
		case tc.code == 0 && (err != nil || got != tc.want):
			t.Errorf("%q: %+v, %v, want %+v", tc.list, got, err, tc.want)
		}
	}
}

// TestColumnGroupPlan pins the text of the plan between two tables' groups, which
// TestOracle_OceanBaseColumnGroupsStored applies: a DROP and an ADD of their own.
func TestColumnGroupPlan(t *testing.T) {
	load := func(create string) *Catalog {
		c, err := LoadSQL("CREATE DATABASE d; USE d; " + create)
		if err != nil {
			t.Fatal(err)
		}
		return c
	}
	from := load("CREATE TABLE t (a INT) WITH COLUMN GROUP(all columns)")
	to := load("CREATE TABLE t (a INT) WITH COLUMN GROUP(each column)")
	n := NormalizerFor(MySQL80)
	diff := DiffWithNormalizer(from, to, n)
	got := GenerateMigrationWithNormalizer(from, to, diff, n).SQL()
	want := "ALTER TABLE `d`.`t` DROP COLUMN GROUP(all columns);\nALTER TABLE `d`.`t` ADD COLUMN GROUP(each column)"
	if strings.TrimSpace(got) != want {
		t.Errorf("plan:\n%s\nwant:\n%s", got, want)
	}
	if d := DiffWithNormalizer(to, to, n); !d.IsEmpty() {
		t.Errorf("a table differs from itself: %+v", d)
	}
	if got := to.ShowCreateTable("d", "t"); !strings.HasSuffix(got, " WITH COLUMN GROUP(each column)") {
		t.Errorf("SHOW CREATE TABLE:\n%s", got)
	}
}
