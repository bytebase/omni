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
