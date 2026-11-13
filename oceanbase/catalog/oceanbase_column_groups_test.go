package catalog

import (
	"context"
	"strings"
	"testing"
)

// TestOracle_OceanBaseColumnGroups holds column_groups.go to the engine: which column
// groups OceanBase takes, in which order it checks them, and what a table holds after
// ADD and DROP COLUMN GROUP, read through the next command's answer.
func TestOracle_OceanBaseColumnGroups(t *testing.T) {
	runExecSteps(t, "oceanbase_catalog_column_groups", []struct {
		name  string
		steps []engineStep
	}{
		{"a group that names its columns", []engineStep{
			{"CREATE TABLE t (a INT PRIMARY KEY, b INT) WITH COLUMN GROUP(g(a))", 1235},
			{"CREATE TABLE t (a INT PRIMARY KEY, b INT) WITH COLUMN GROUP(all columns, g(a, b))", 1235},
			{"CREATE TABLE t (a INT PRIMARY KEY, b INT)", taken},
			{"ALTER TABLE t ADD COLUMN GROUP(g(a, b))", 1235},
			{"ALTER TABLE t ADD COLUMN GROUP(EACH COLUMN, g(a))", 1235},
			{"ALTER TABLE t ADD COLUMN GROUP(`each`(a))", 1235},
			{"ALTER TABLE t DROP COLUMN GROUP(g(a))", 1235},
			{"ALTER TABLE t DROP COLUMN GROUP(ALL COLUMNS, g(a))", 1235},
			{"ALTER TABLE nosuch ADD COLUMN GROUP(g(a))", 1146},
		}},
		// A named group and a repeated one are refused in written order.
		{"a group written twice", []engineStep{
			{"CREATE TABLE t (a INT PRIMARY KEY) WITH COLUMN GROUP(all columns, all columns)", 1060},
			{"CREATE TABLE t (a INT PRIMARY KEY) WITH COLUMN GROUP(each column, each column)", 1060},
			{"CREATE TABLE t (a INT PRIMARY KEY)", taken},
			{"ALTER TABLE t ADD COLUMN GROUP(g(a), ALL COLUMNS, ALL COLUMNS)", 1235},
			{"ALTER TABLE t ADD COLUMN GROUP(ALL COLUMNS, ALL COLUMNS, g(a))", 1060},
			{"ALTER TABLE t ADD COLUMN GROUP(EACH COLUMN, EACH COLUMN)", 1060},
			{"ALTER TABLE t DROP COLUMN GROUP(ALL COLUMNS, ALL COLUMNS)", 1060},
			{"ALTER TABLE t ADD COLUMN GROUP(EACH COLUMN)", taken},
			{"ALTER TABLE t DROP COLUMN GROUP(EACH COLUMN, EACH COLUMN)", 1060},
		}},
		{"before the table's name", []engineStep{
			{"CREATE TABLE t (a INT PRIMARY KEY)", taken},
			{"CREATE TABLE t (a INT) WITH COLUMN GROUP(g(a))", 1235},
			{"CREATE TABLE IF NOT EXISTS t (a INT) WITH COLUMN GROUP(g(a))", 1235},
			{"CREATE TABLE IF NOT EXISTS t (a INT) WITH COLUMN GROUP(all columns, all columns)", 1060},
			{"CREATE TABLE IF NOT EXISTS t (a INT) WITH COLUMN GROUP(all columns)", taken},
			{"CREATE TABLE t (a INT) WITH COLUMN GROUP(all columns)", 1050},
		}},
		// A group already on the table is 1060 and one it lacks 4185, ALL COLUMNS
		// checked before EACH COLUMN in either order.
		{"add and drop", []engineStep{
			{"CREATE TABLE t (a INT PRIMARY KEY, b INT)", taken},
			{"ALTER TABLE t DROP COLUMN GROUP(ALL COLUMNS)", 4185},
			{"ALTER TABLE t DROP COLUMN GROUP(EACH COLUMN)", 4185},
			{"ALTER TABLE t ADD COLUMN GROUP(EACH COLUMN)", taken},
			{"ALTER TABLE t ADD COLUMN GROUP(EACH COLUMN)", 1060},
			{"ALTER TABLE t ADD COLUMN GROUP(ALL COLUMNS, EACH COLUMN)", 1060},
			{"ALTER TABLE t DROP COLUMN GROUP(EACH COLUMN, ALL COLUMNS)", 4185},
			{"ALTER TABLE t ADD COLUMN GROUP(ALL COLUMNS)", taken},
			{"ALTER TABLE t ADD COLUMN GROUP(EACH COLUMN, ALL COLUMNS)", 1060},
			{"ALTER TABLE t DROP COLUMN GROUP(ALL COLUMNS, EACH COLUMN)", taken},
			{"ALTER TABLE t DROP COLUMN GROUP(EACH COLUMN)", 4185},
			{"ALTER TABLE t ADD COLUMN GROUP(EACH COLUMN, ALL COLUMNS)", taken},
			{"ALTER TABLE t DROP COLUMN GROUP(ALL COLUMNS)", taken},
			{"ALTER TABLE t DROP COLUMN GROUP(ALL COLUMNS)", 4185},
			{"ALTER TABLE t DROP COLUMN GROUP(EACH COLUMN)", taken},
		}},
		// ADD takes DELAYED only beside EACH COLUMN, and checks it after the list and
		// before the table's groups; DROP takes it never, before anything else.
		{"delayed", []engineStep{
			{"CREATE TABLE t (a INT PRIMARY KEY, b INT)", taken},
			{"ALTER TABLE t ADD COLUMN GROUP(ALL COLUMNS) DELAYED", 1235},
			{"ALTER TABLE t ADD COLUMN GROUP(ALL COLUMNS, ALL COLUMNS) DELAYED", 1060},
			{"ALTER TABLE t ADD COLUMN GROUP(EACH COLUMN, g(a)) DELAYED", 1235},
			{"ALTER TABLE t DROP COLUMN GROUP(ALL COLUMNS) DELAYED", 1235},
			{"ALTER TABLE t ADD COLUMN GROUP(EACH COLUMN) DELAYED", taken},
			{"ALTER TABLE t DROP COLUMN GROUP(EACH COLUMN) DELAYED", 1235},
			{"ALTER TABLE t DROP COLUMN GROUP(EACH COLUMN, EACH COLUMN) DELAYED", 1235},
			{"ALTER TABLE t ADD COLUMN GROUP(EACH COLUMN) DELAYED", 1060},
			{"ALTER TABLE t DROP COLUMN GROUP(EACH COLUMN)", taken},
			{"ALTER TABLE t ADD COLUMN GROUP(ALL COLUMNS, EACH COLUMN) DELAYED", taken},
			{"ALTER TABLE t ADD COLUMN GROUP(ALL COLUMNS) DELAYED", 1235},
		}},
		{"CREATE TABLE LIKE copies them", []engineStep{
			{"CREATE TABLE s (a INT PRIMARY KEY) WITH COLUMN GROUP(each column)", taken},
			{"CREATE TABLE t LIKE s", taken},
			{"ALTER TABLE t ADD COLUMN GROUP(EACH COLUMN)", 1060},
			{"ALTER TABLE t DROP COLUMN GROUP(EACH COLUMN)", taken},
		}},
	})
}

// TestOracle_OceanBaseColumnGroupsStored holds to the engine what the catalog stores for a
// table's column groups: SHOW CREATE TABLE prints all columns before each column whatever
// order they were written in, and the plan between two tables' groups, applied, stores the
// second's.
func TestOracle_OceanBaseColumnGroupsStored(t *testing.T) {
	o := oceanBaseOracle(t)
	ctx := context.Background()
	for i, form := range []struct{ create, stored string }{
		{"CREATE TABLE t (id INT PRIMARY KEY, a INT) WITH COLUMN GROUP(all columns)", "all columns"},
		{"CREATE TABLE t (id INT PRIMARY KEY, a INT) WITH COLUMN GROUP(EACH COLUMN)", "each column"},
		{"CREATE TABLE t (id INT PRIMARY KEY, a INT) WITH COLUMN GROUP(each column, all columns)", "all columns, each column"},
		{"CREATE TABLE t (id INT PRIMARY KEY, a INT) PARTITION BY HASH(id) PARTITIONS 2 WITH COLUMN GROUP(each column)", "each column"},
	} {
		t.Run(form.create, func(t *testing.T) {
			readback, err := storeAndRead(ctx, o, "obcg_"+string(rune('a'+i)), form.create, "t")
			if err != nil {
				t.Fatal(err)
			}
			if want := " WITH COLUMN GROUP(" + form.stored + ")"; !strings.HasSuffix(readback, want) {
				t.Errorf("recorded%s, but the engine stored\n%s", want, readback)
			}
			c, _ := loadOneTable(t, serverCharsetFor(o.version), form.create, "t")
			if got := c.ShowCreateTable(c.CurrentDatabase(), "t"); !strings.HasSuffix(got, " WITH COLUMN GROUP("+form.stored+")") {
				t.Errorf("the catalog prints\n%s", got)
			}
			assertDiffEmptyAgainstStored(t, o, readback, form.create, "t")
		})
	}
	n := NormalizerFor(o.version)
	const table = "CREATE TABLE t (id INT PRIMARY KEY, a INT)"
	for _, p := range []migrationProbe{
		{"cg-create", "t", "", table + " WITH COLUMN GROUP(all columns, each column)"},
		{"cg-add-all", "t", table, table + " WITH COLUMN GROUP(all columns)"},
		{"cg-add-both", "t", table, table + " WITH COLUMN GROUP(all columns, each column)"},
		{"cg-all-to-each", "t", table + " WITH COLUMN GROUP(all columns)", table + " WITH COLUMN GROUP(each column)"},
		{"cg-drop-both", "t", table + " WITH COLUMN GROUP(all columns, each column)", table},
		{"cg-and-a-column", "t", table + " WITH COLUMN GROUP(each column)", "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT) WITH COLUMN GROUP(all columns)"},
	} {
		t.Run(p.id, func(t *testing.T) { assertApplyCorrect(t, o, n, p) })
	}
}
