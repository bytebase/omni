package catalog

import (
	"strings"
	"testing"
)

func planFor(t *testing.T, from, to string) *MigrationPlan {
	t.Helper()
	a, b := loadCat(t, dbDDL+from), loadCat(t, dbDDL+to)
	return GenerateMigration(a, b, Diff(a, b))
}

func TestMigration_EmptyDiffEmptyPlan(t *testing.T) {
	if sql := planFor(t, "CREATE TABLE t (a INT);", "CREATE TABLE t (a INT(11));").SQL(); sql != "" {
		t.Errorf("plan = %q, want empty", sql)
	}
	if plan := GenerateMigration(New(), New(), nil); len(plan.Ops) != 0 {
		t.Errorf("nil diff: %d ops", len(plan.Ops))
	}
}

func TestMigration_Deterministic(t *testing.T) {
	from := "CREATE TABLE a (id INT); CREATE TABLE b (id INT); CREATE SEQUENCE s;"
	to := "CREATE TABLE a (id INT, x INT, KEY k (x)); CREATE TABLE c (id INT); CREATE SEQUENCE s INCREMENT BY 2; CREATE SEQUENCE u;"
	first := planFor(t, from, to).SQL()
	for i := 0; i < 20; i++ {
		if got := planFor(t, from, to).SQL(); got != first {
			t.Fatalf("plan changed between runs:\n%s\n---\n%s", first, got)
		}
	}
}

func TestMigration_SequenceLimitsLeavingOutTheStartWarn(t *testing.T) {
	plan := planFor(t, "CREATE SEQUENCE s;", "CREATE SEQUENCE s START WITH 100 MINVALUE 10;")
	if strings.Contains(plan.SQL(), "RESTART") {
		t.Errorf("the plan restarts the sequence:\n%s", plan.SQL())
	}
	if !plan.HasWarnings() {
		t.Errorf("no warning that the next value may lie outside the new limits:\n%s", plan.SQL())
	}
}

func TestMigration_ViewBodyQualifiesItsRelations(t *testing.T) {
	got := planFor(t, "", "CREATE TABLE t (a INT); CREATE VIEW v AS WITH c AS (SELECT a FROM t) SELECT c.a FROM c JOIN t ON c.a = t.a;").SQL()
	if strings.Count(got, "`app`.`t`") < 2 || strings.Contains(got, "`app`.`c`") {
		t.Errorf("want t qualified with app and the CTE c bare:\n%s", got)
	}
}

func TestMigration_AddedObjectsKeepTheirSessionContext(t *testing.T) {
	a := loadCat(t, dbDDL+"CREATE TABLE t (a VARCHAR(10));")
	b := loadCat(t, dbDDL+"CREATE TABLE t (a VARCHAR(10)); CREATE FUNCTION f() RETURNS VARCHAR(10) DETERMINISTIC RETURN 'a' || 'b'; CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 'x' || 'y'; CREATE EVENT e ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' DO SELECT 'a' || 'b';")
	ctx := SessionContext{SQLMode: "PIPES_AS_CONCAT"}
	b.ApplySessionContext(SessionContextMap{
		Functions: map[string]SessionContext{"f": ctx},
		Triggers:  map[string]SessionContext{"tr": ctx},
		Events:    map[string]SessionContext{"e": ctx},
	})
	for _, op := range GenerateMigration(a, b, Diff(a, b)).Ops {
		switch op.Type {
		case OpCreateFunction, OpCreateTrigger, OpCreateEvent:
			if !strings.Contains(op.SQL, "SET sql_mode = 'PIPES_AS_CONCAT'") {
				t.Errorf("%s created without its sql_mode:\n%s", op.ObjectName, op.SQL)
			}
		}
	}
}

func TestMigration_MariaDBStatements(t *testing.T) {
	for _, tc := range []struct{ name, from, to, want string }{
		{"sequence create", "", "CREATE SEQUENCE s START WITH 5;", "CREATE SEQUENCE `app`.`s` start with 5 minvalue 1 maxvalue 9223372036854775806 increment by 1 cache 1000 nocycle ENGINE=InnoDB"},
		{"sequence drop", "CREATE SEQUENCE s;", "", "DROP SEQUENCE `app`.`s`"},
		{"sequence alter", "CREATE SEQUENCE s;", "CREATE SEQUENCE s INCREMENT BY 2 COMMENT 'c';", "ALTER SEQUENCE `app`.`s` INCREMENT BY 2;\nALTER TABLE `app`.`s` COMMENT='c'"},
		{"check drop", "CREATE TABLE t (a INT, CONSTRAINT c CHECK (a > 0));", "CREATE TABLE t (a INT);", "ALTER TABLE `app`.`t` DROP CONSTRAINT `c`"},
		{"column check", "CREATE TABLE t (a INT);", "CREATE TABLE t (a INT CHECK (a > 0));", "MODIFY COLUMN `a` int(11) DEFAULT NULL CHECK (`a` > 0)"},
		{"versioning", "CREATE TABLE t (a INT);", "CREATE TABLE t (a INT) WITH SYSTEM VERSIONING;", "ALTER TABLE `app`.`t` ADD SYSTEM VERSIONING"},
		{"keep history", "CREATE TABLE t (a INT) WITH SYSTEM VERSIONING;", "CREATE TABLE t (a INT, b INT) WITH SYSTEM VERSIONING;", "SET @@system_versioning_alter_history = KEEP"},
		{"period", "CREATE TABLE t (s DATE NOT NULL, e DATE NOT NULL);", "CREATE TABLE t (s DATE NOT NULL, e DATE NOT NULL, PERIOD FOR p (s, e));", "ALTER TABLE `app`.`t` ADD PERIOD FOR `p` (`s`, `e`)"},
		{"option reset", "CREATE TABLE t (a INT) STATS_PERSISTENT=1 COMMENT 'x';", "CREATE TABLE t (a INT);", "ALTER TABLE `app`.`t` STATS_PERSISTENT=DEFAULT COMMENT=''"},
		{"event qualified", "", "CREATE EVENT e ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' DO SELECT 1;", "CREATE EVENT `app`.`e`"},
		{"ignored key", "CREATE TABLE t (a INT, KEY k (a));", "CREATE TABLE t (a INT, KEY k (a) IGNORED);", "ADD KEY `k` (`a`) IGNORED"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := planFor(t, tc.from, tc.to).SQL(); !strings.Contains(got, tc.want) {
				t.Errorf("plan:\n%s\nwant it to contain:\n%s", got, tc.want)
			}
		})
	}
}
