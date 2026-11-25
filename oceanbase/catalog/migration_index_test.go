package catalog

import (
	"strings"
	"testing"
)

// TestMigration_InvisibleIndex pins how an index is made invisible on OceanBase, whose
// ALTER TABLE ... ADD KEY does not take INVISIBLE: the index is added, then hidden; and an
// index that only changes visibility is altered in place, not built again.
// TestOracle_OceanBaseInvisibleIndex asks the engine.
func TestMigration_InvisibleIndex(t *testing.T) {
	const cols = "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT"
	tests := []struct {
		name     string
		from, to string
		want     []string
	}{
		{"added invisible",
			cols + ");", cols + ", KEY k (a) INVISIBLE);",
			[]string{"ALTER TABLE `app`.`t` ADD KEY `k` (`a`) /*!80000 INVISIBLE */", "ALTER TABLE `app`.`t` ALTER INDEX `k` INVISIBLE"}},
		{"added visible",
			cols + ");", cols + ", KEY k (a));",
			[]string{"ALTER TABLE `app`.`t` ADD KEY `k` (`a`)"}},
		{"hidden",
			cols + ", KEY k (a));", cols + ", KEY k (a) INVISIBLE);",
			[]string{"ALTER TABLE `app`.`t` ALTER INDEX `k` INVISIBLE"}},
		{"shown",
			cols + ", KEY k (a) INVISIBLE);", cols + ", KEY k (a));",
			[]string{"ALTER TABLE `app`.`t` ALTER INDEX `k` VISIBLE"}},
		{"hidden and rebuilt over another column",
			cols + ", KEY k (a));", cols + ", KEY k (b) INVISIBLE);",
			[]string{"ALTER TABLE `app`.`t` DROP INDEX `k`", "ALTER TABLE `app`.`t` ADD KEY `k` (`b`) /*!80000 INVISIBLE */", "ALTER TABLE `app`.`t` ALTER INDEX `k` INVISIBLE"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			plan := planFor(t, dbDDL+tc.from, dbDDL+tc.to, NormalizerFor(MySQL80))
			var got []string
			for _, op := range plan.Ops {
				got = append(got, op.SQL)
			}
			if strings.Join(got, "\n") != strings.Join(tc.want, "\n") {
				t.Errorf("plan:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(tc.want, "\n"))
			}
		})
	}
}

// TestAlterAddKeyInvisibleIgnored pins that ALTER TABLE ... ADD KEY leaves the key
// visible whatever it says, and that ALTER INDEX then hides it, as
// TestOracle_OceanBaseInvisibleIndex measures on the engine.
func TestAlterAddKeyInvisibleIgnored(t *testing.T) {
	c := New()
	mustExec(t, c, "CREATE DATABASE d; USE d; CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT);")
	mustExec(t, c, "ALTER TABLE t ADD KEY k_plain (a) INVISIBLE; ALTER TABLE t ADD KEY k_comment (a, b) /*!80000 INVISIBLE */;")
	tbl := c.GetDatabase("d").GetTable("t")
	for _, name := range []string{"k_plain", "k_comment"} {
		for _, idx := range tbl.Indexes {
			if idx.Name == name && !idx.Visible {
				t.Errorf("%s is invisible after ALTER TABLE ... ADD KEY; OceanBase leaves it visible", name)
			}
		}
	}
	mustExec(t, c, "ALTER TABLE t ALTER INDEX k_plain INVISIBLE;")
	for _, idx := range tbl.Indexes {
		if idx.Name == "k_plain" && idx.Visible {
			t.Error("k_plain is visible after ALTER INDEX ... INVISIBLE")
		}
	}
}
