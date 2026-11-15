package catalog

import "testing"

func TestGenerateMigrationOrder(t *testing.T) {
	from := execCatalogIn(t, "db", `
CREATE TABLE p (id int PRIMARY KEY NONCLUSTERED);
CREATE TABLE c (id int PRIMARY KEY, pid int, a int, b int, e varchar(10), KEY ab (a, b), KEY ex ((lower(e))),
  CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p (id));
CREATE VIEW v AS SELECT id FROM c;`)
	to := execCatalogIn(t, "db", `
CREATE TABLE p (id bigint PRIMARY KEY NONCLUSTERED);
CREATE TABLE c (id int PRIMARY KEY, n int, pid bigint, a int, e varchar(20), KEY ex ((lower(e))), KEY fk (pid),
  CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p (id), CONSTRAINT ck CHECK (a > 0));
CREATE VIEW v2 AS SELECT n FROM v;
CREATE VIEW v AS SELECT id, n FROM c;
CREATE SEQUENCE s;`)
	got := GenerateMigration(from, to, Diff(from, to)).SQL()
	want := "ALTER TABLE `db`.`c` DROP FOREIGN KEY `fk`;\n" +
		"ALTER TABLE `db`.`c` DROP INDEX `ab`;\n" +
		"ALTER TABLE `db`.`c` DROP INDEX `ex`;\n" +
		"ALTER TABLE `db`.`c` DROP COLUMN `b`;\n" +
		"CREATE SEQUENCE `db`.`s` start with 1 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB;\n" +
		"ALTER TABLE `db`.`c` ADD COLUMN `n` int DEFAULT NULL AFTER `id`;\n" +
		"ALTER TABLE `db`.`c` MODIFY COLUMN `e` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL;\n" +
		"ALTER TABLE `db`.`c` MODIFY COLUMN `pid` bigint DEFAULT NULL;\n" +
		"ALTER TABLE `db`.`c` ADD KEY `ex` ((lower(`e`)));\n" +
		"ALTER TABLE `db`.`c` ADD CONSTRAINT `ck` CHECK (`a` > 0);\n" +
		"ALTER TABLE `db`.`p` MODIFY COLUMN `id` bigint NOT NULL;\n" +
		"CREATE OR REPLACE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `db`.`v` (`id`, `n`) AS SELECT `id` AS `id`,`n` AS `n` FROM `db`.`c`;\n" +
		"CREATE ALGORITHM=UNDEFINED SQL SECURITY DEFINER VIEW `db`.`v2` (`n`) AS SELECT `n` AS `n` FROM `db`.`v`;\n" +
		"ALTER TABLE `db`.`c` ADD CONSTRAINT `fk` FOREIGN KEY (`pid`) REFERENCES `p` (`id`)"
	if got != want {
		t.Errorf("got\n%s\nwant\n%s", got, want)
	}
	if plan := GenerateMigration(to, to, Diff(to, to)); len(plan.Ops) != 0 {
		t.Errorf("an empty diff gave %s", plan.SQL())
	}
}

func TestGenerateMigrationSetsTiFlashReplicas(t *testing.T) {
	from := execCatalogIn(t, "db", "CREATE TABLE t (a int); ALTER TABLE t SET TIFLASH REPLICA 1")
	to := execCatalogIn(t, "db", "CREATE TABLE t (a int); ALTER TABLE t SET TIFLASH REPLICA 2 LOCATION LABELS 'zone'")
	if got, want := GenerateMigration(from, to, Diff(from, to)).SQL(), "ALTER TABLE `db`.`t` SET TIFLASH REPLICA 2 LOCATION LABELS 'zone'"; got != want {
		t.Errorf("got %s, want %s", got, want)
	}
	to = execCatalogIn(t, "db", "CREATE TABLE t (a int)")
	if got, want := GenerateMigration(from, to, Diff(from, to)).SQL(), "ALTER TABLE `db`.`t` SET TIFLASH REPLICA 0"; got != want {
		t.Errorf("got %s, want %s", got, want)
	}
	if got, want := GenerateMigration(to, from, Diff(to, from)).SQL(), "ALTER TABLE `db`.`t` SET TIFLASH REPLICA 1"; got != want {
		t.Errorf("got %s, want %s", got, want)
	}
	empty := execCatalogIn(t, "db", "")
	if got, want := GenerateMigration(empty, from, Diff(empty, from)).SQL(), "CREATE TABLE `db`.`t` (\n  `a` int DEFAULT NULL\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;\nALTER TABLE `db`.`t` SET TIFLASH REPLICA 1"; got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
