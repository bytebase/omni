package catalog

import (
	"strings"
	"testing"
)

// planCase is a change and the statements its plan must be, in order.
type planCase struct {
	name     string
	from, to string
	want     []string
}

func checkPlans(t *testing.T, cases []planCase) {
	t.Helper()
	for _, tc := range cases {
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

// TestMigration_StatementOrder pins the statements OceanBase is given and their order:
// one clause each, but for a primary key's replacement, with the statement that frees a
// column ahead of the one that needs it free. TestOracle_OceanBaseAlterTableClauses asks
// the engine about each rule, and the apply tests run plans like these on it.
func TestMigration_StatementOrder(t *testing.T) {
	checkPlans(t, []planCase{
		{"an AUTO_INCREMENT column and its key are two statements",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT);",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT, seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_seq (seq));",
			[]string{
				"ALTER TABLE `app`.`t` ADD COLUMN `seq` bigint unsigned NOT NULL AUTO_INCREMENT",
				"ALTER TABLE `app`.`t` ADD UNIQUE KEY `uk_seq` (`seq`)",
			}},
		{"a new table gets that key after it is created",
			"",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, seq BIGINT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_seq (seq));",
			[]string{
				"CREATE TABLE `app`.`t` (\n  `id` int NOT NULL,\n  `seq` bigint NOT NULL AUTO_INCREMENT,\n  PRIMARY KEY (`id`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4",
				"ALTER TABLE `app`.`t` ADD UNIQUE KEY `uk_seq` (`seq`)",
			}},
		{"AUTO_INCREMENT leaves a column before another column takes it",
			"CREATE TABLE t (b INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, KEY kb (b), KEY ka (a));",
			"CREATE TABLE t (b INT NOT NULL, a INT NOT NULL AUTO_INCREMENT, KEY kb (b), KEY ka (a));",
			[]string{
				"ALTER TABLE `app`.`t` MODIFY COLUMN `b` int NOT NULL",
				"ALTER TABLE `app`.`t` MODIFY COLUMN `a` int NOT NULL AUTO_INCREMENT",
			}},
		{"and before a new column takes it",
			"CREATE TABLE t (z INT NOT NULL AUTO_INCREMENT, KEY kz (z));",
			"CREATE TABLE t (z INT NOT NULL, a INT NOT NULL AUTO_INCREMENT, KEY kz (z), KEY ka (a));",
			[]string{
				"ALTER TABLE `app`.`t` MODIFY COLUMN `z` int NOT NULL",
				"ALTER TABLE `app`.`t` ADD COLUMN `a` int NOT NULL AUTO_INCREMENT",
				"ALTER TABLE `app`.`t` ADD KEY `ka` (`a`)",
			}},
		{"a primary key is replaced in one statement",
			"CREATE TABLE t (id INT NOT NULL, a INT NOT NULL, PRIMARY KEY (id));",
			"CREATE TABLE t (id INT NOT NULL, a INT NOT NULL, PRIMARY KEY (a));",
			[]string{"ALTER TABLE `app`.`t` DROP PRIMARY KEY, ADD PRIMARY KEY (`a`)"}},
		{"a column that leaves the key is dropped after the key is replaced",
			"CREATE TABLE t (a INT NOT NULL, c INT NOT NULL, PRIMARY KEY (c, a));",
			"CREATE TABLE t (c INT NOT NULL, PRIMARY KEY (c));",
			[]string{
				"ALTER TABLE `app`.`t` DROP PRIMARY KEY, ADD PRIMARY KEY (`c`)",
				"ALTER TABLE `app`.`t` DROP COLUMN `a`",
			}},
		{"and made nullable after it",
			"CREATE TABLE t (id INT NOT NULL, a INT NOT NULL, PRIMARY KEY (id, a));",
			"CREATE TABLE t (id INT NOT NULL, a INT, PRIMARY KEY (id));",
			[]string{
				"ALTER TABLE `app`.`t` DROP PRIMARY KEY, ADD PRIMARY KEY (`id`)",
				"ALTER TABLE `app`.`t` MODIFY COLUMN `a` int",
			}},
		{"of a column that leaves the key, changes type and becomes nullable, only NULL waits",
			"CREATE TABLE t (id INT NOT NULL, a VARCHAR(5) NOT NULL, PRIMARY KEY (id, a));",
			"CREATE TABLE t (id INT NOT NULL, a VARCHAR(20), PRIMARY KEY (id), KEY ka (a(10)));",
			[]string{
				"ALTER TABLE `app`.`t` MODIFY COLUMN `a` varchar(20) NOT NULL",
				"ALTER TABLE `app`.`t` ADD KEY `ka` (`a`(10))",
				"ALTER TABLE `app`.`t` DROP PRIMARY KEY, ADD PRIMARY KEY (`id`)",
				"ALTER TABLE `app`.`t` MODIFY COLUMN `a` varchar(20)",
			}},
		{"a key column that hands AUTO_INCREMENT on and is dropped does so in two statements",
			"CREATE TABLE t (c INT NOT NULL AUTO_INCREMENT, b INT NOT NULL, PRIMARY KEY (c));",
			"CREATE TABLE t (b INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (b));",
			[]string{
				"ALTER TABLE `app`.`t` MODIFY COLUMN `c` int NOT NULL",
				"ALTER TABLE `app`.`t` MODIFY COLUMN `b` int NOT NULL AUTO_INCREMENT",
				"ALTER TABLE `app`.`t` DROP PRIMARY KEY, ADD PRIMARY KEY (`b`)",
				"ALTER TABLE `app`.`t` DROP COLUMN `c`",
			}},
		{"and one that becomes nullable likewise",
			"CREATE TABLE t (c INT NOT NULL AUTO_INCREMENT, b INT NOT NULL, PRIMARY KEY (c));",
			"CREATE TABLE t (c INT, b INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (b));",
			[]string{
				"ALTER TABLE `app`.`t` MODIFY COLUMN `c` int NOT NULL",
				"ALTER TABLE `app`.`t` MODIFY COLUMN `b` int NOT NULL AUTO_INCREMENT",
				"ALTER TABLE `app`.`t` DROP PRIMARY KEY, ADD PRIMARY KEY (`b`)",
				"ALTER TABLE `app`.`t` MODIFY COLUMN `c` int",
			}},
		{"a key that is dropped goes before its columns",
			"CREATE TABLE t (c INT NOT NULL AUTO_INCREMENT PRIMARY KEY, a INT);",
			"CREATE TABLE t (a INT);",
			[]string{
				"ALTER TABLE `app`.`t` DROP PRIMARY KEY",
				"ALTER TABLE `app`.`t` DROP COLUMN `c`",
			}},
		{"an index is dropped in one statement and added in another",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL, a INT, UNIQUE KEY uk_c (c));",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL, a INT, UNIQUE KEY uk_c (c, a));",
			[]string{
				"ALTER TABLE `app`.`t` DROP INDEX `uk_c`",
				"ALTER TABLE `app`.`t` ADD UNIQUE KEY `uk_c` (`c`,`a`)",
			}},
		{"a NOT NULL spatial column is added nullable first",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY);",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, g GEOMETRY NOT NULL SRID 0, p POINT);",
			[]string{
				"ALTER TABLE `app`.`t` ADD COLUMN `g` geometry /*!80003 SRID 0 */",
				"ALTER TABLE `app`.`t` MODIFY COLUMN `g` geometry NOT NULL /*!80003 SRID 0 */",
				"ALTER TABLE `app`.`t` ADD COLUMN `p` point",
			}},
		{"a NOT NULL spatial column is modified through nullable",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, g GEOMETRY NOT NULL SRID 0, j JSON NOT NULL);",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, g POINT NOT NULL SRID 0 COMMENT 'c', j JSON NOT NULL COMMENT 'c');",
			[]string{
				"ALTER TABLE `app`.`t` MODIFY COLUMN `g` point /*!80003 SRID 0 */ COMMENT 'c'",
				"ALTER TABLE `app`.`t` MODIFY COLUMN `g` point NOT NULL /*!80003 SRID 0 */ COMMENT 'c'",
				"ALTER TABLE `app`.`t` MODIFY COLUMN `j` json NOT NULL COMMENT 'c'",
			}},
		{"a repartition that lets a column go runs before the key is replaced",
			"CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2010), PARTITION pmax VALUES LESS THAN MAXVALUE);",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY) PARTITION BY HASH (id) PARTITIONS 2;",
			[]string{
				"ALTER TABLE `app`.`t` /*!50100 PARTITION BY HASH (`id`)\n(PARTITION p0 ENGINE = innodb,\n PARTITION p1 ENGINE = innodb) */",
				"ALTER TABLE `app`.`t` DROP PRIMARY KEY, ADD PRIMARY KEY (`id`)",
				"ALTER TABLE `app`.`t` DROP COLUMN `dt`",
			}},
	})
}

// TestMigration_ForeignKeyHoldsItsLastIndex pins when a foreign key the diff leaves alone
// is dropped and added back: when the plan drops the last index over exactly its columns,
// which OceanBase refuses while the foreign key lives.
// TestOracle_OceanBaseForeignKeyIndex asks the engine.
func TestMigration_ForeignKeyHoldsItsLastIndex(t *testing.T) {
	const parent = "CREATE TABLE p (id INT NOT NULL PRIMARY KEY); "
	const fk = "CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id)"
	const addFK = "ALTER TABLE `app`.`c` ADD CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `p` (`id`)"
	checkPlans(t, []planCase{
		{"the last index over its column is dropped",
			parent + "CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY my_idx (pid), " + fk + ");",
			parent + "CREATE TABLE c (id INT PRIMARY KEY, pid INT, " + fk + ");",
			[]string{
				"ALTER TABLE `app`.`c` DROP FOREIGN KEY `fk`",
				"ALTER TABLE `app`.`c` DROP INDEX `my_idx`",
				addFK,
			}},
		{"it is replaced by another, which is not there yet when it goes",
			parent + "CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY a1 (pid), " + fk + ");",
			parent + "CREATE TABLE c (id INT PRIMARY KEY, pid INT, UNIQUE KEY a2 (pid), " + fk + ");",
			[]string{
				"ALTER TABLE `app`.`c` DROP FOREIGN KEY `fk`",
				"ALTER TABLE `app`.`c` DROP INDEX `a1`",
				"ALTER TABLE `app`.`c` ADD UNIQUE KEY `a2` (`pid`)",
				addFK,
			}},
		{"another index over its column stays",
			parent + "CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY a1 (pid), KEY a2 (pid), " + fk + ");",
			parent + "CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY a1 (pid), " + fk + ");",
			[]string{"ALTER TABLE `app`.`c` DROP INDEX `a2`"}},
		{"a wider index is not held",
			parent + "CREATE TABLE c (id INT PRIMARY KEY, pid INT, x INT, KEY w (pid, x), " + fk + ");",
			parent + "CREATE TABLE c (id INT PRIMARY KEY, pid INT, x INT, " + fk + ");",
			[]string{"ALTER TABLE `app`.`c` DROP INDEX `w`"}},
		{"its index only changes visibility",
			parent + "CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY a1 (pid), " + fk + ");",
			parent + "CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY a1 (pid) INVISIBLE, " + fk + ");",
			[]string{"ALTER TABLE `app`.`c` ALTER INDEX `a1` INVISIBLE"}},
		{"the primary key it is over is dropped",
			parent + "CREATE TABLE c (pid INT NOT NULL, a INT NOT NULL, PRIMARY KEY (pid), " + fk + ");",
			parent + "CREATE TABLE c (pid INT NOT NULL, a INT NOT NULL, " + fk + ");",
			[]string{
				"ALTER TABLE `app`.`c` DROP FOREIGN KEY `fk`",
				"ALTER TABLE `app`.`c` DROP PRIMARY KEY",
				addFK,
			}},
		{"the primary key it is over is replaced, which OceanBase takes",
			parent + "CREATE TABLE c (pid INT NOT NULL, a INT NOT NULL, PRIMARY KEY (pid), " + fk + ");",
			parent + "CREATE TABLE c (pid INT NOT NULL, a INT NOT NULL, PRIMARY KEY (a), " + fk + ");",
			[]string{"ALTER TABLE `app`.`c` DROP PRIMARY KEY, ADD PRIMARY KEY (`a`)"}},
		{"the foreign key is dropped with its index",
			parent + "CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY my_idx (pid), " + fk + ");",
			parent + "CREATE TABLE c (id INT PRIMARY KEY, pid INT);",
			[]string{
				"ALTER TABLE `app`.`c` DROP FOREIGN KEY `fk`",
				"ALTER TABLE `app`.`c` DROP INDEX `my_idx`",
			}},
		{"a foreign key that had no index of its own drops alone",
			parent + "CREATE TABLE c (id INT PRIMARY KEY, pid INT, " + fk + ");",
			parent + "CREATE TABLE c (id INT PRIMARY KEY, pid INT);",
			[]string{"ALTER TABLE `app`.`c` DROP FOREIGN KEY `fk`"}},
	})
}

// TestMigration_WarnsWhereOceanBaseHasNoStatement pins the changes the plan cannot make on
// OceanBase: the statement MySQL would run is kept, and its op says that OceanBase refuses
// it. The apply tests check each warned statement is refused (applyPlan).
func TestMigration_WarnsWhereOceanBaseHasNoStatement(t *testing.T) {
	for _, tc := range []struct {
		name     string
		from, to string
		stmt     string
		warning  string
	}{
		{"partitioning removed",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT) PARTITION BY HASH (id) PARTITIONS 4;",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT);",
			"ALTER TABLE `app`.`t` REMOVE PARTITIONING", "cannot remove the partitioning"},
		{"a spatial column's SRID changed",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, g GEOMETRY NOT NULL SRID 0);",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, g GEOMETRY NOT NULL SRID 4326);",
			"ALTER TABLE `app`.`t` MODIFY COLUMN `g` geometry NOT NULL /*!80003 SRID 4326 */", "cannot change the SRID"},
		{"a spatial column's SRID taken away",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, g GEOMETRY NOT NULL SRID 0);",
			"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, g GEOMETRY NOT NULL);",
			"ALTER TABLE `app`.`t` MODIFY COLUMN `g` geometry NOT NULL", "cannot change the SRID"},
		{"a generated column's expression, under the primary key",
			"CREATE TABLE t (a INT, g INT GENERATED ALWAYS AS (a+1) STORED NOT NULL, PRIMARY KEY (g));",
			"CREATE TABLE t (a INT, g INT GENERATED ALWAYS AS (a+2) STORED NOT NULL, PRIMARY KEY (g));",
			"ALTER TABLE `app`.`t` MODIFY COLUMN `g` int GENERATED ALWAYS AS ((`a` + 2)) STORED NOT NULL", "cannot change the expression"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			plan := planFor(t, dbDDL+tc.from, dbDDL+tc.to, NormalizerFor(MySQL80))
			warned := plan.Warnings()
			if len(plan.Ops) != 1 || len(warned) != 1 || warned[0].SQL != tc.stmt || !strings.Contains(warned[0].Warning, tc.warning) {
				t.Errorf("want one statement, warned %q:\n  %s\ngot:\n%s\n  warnings: %v", tc.warning, tc.stmt, plan.SQL(), warned)
			}
		})
	}
	t.Run("a change OceanBase takes carries no warning", func(t *testing.T) {
		plan := planFor(t, dbDDL+"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, g GEOMETRY NOT NULL SRID 0);",
			dbDDL+"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, g GEOMETRY NOT NULL SRID 0 COMMENT 'c');", NormalizerFor(MySQL80))
		if len(plan.Ops) == 0 || plan.HasWarnings() {
			t.Errorf("want a plan with no warning, got:\n%s\n  warnings: %v", plan.SQL(), plan.Warnings())
		}
	})
}
