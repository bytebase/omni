package catalog

import (
	"strings"
	"testing"
)

// Hermetic plan-shape tests for the AUTO_INCREMENT / backing-key ALTER grouping
// (migration_autoinc.go). They assert WHICH clauses share a statement and how the surrounding
// ops are (not) disturbed; the live-engine proof that the grouped statements are the forms
// MySQL accepts is migration_autoinc_oracle_test.go.

// aigPlan diffs two single-table schemas (already wrapped in dbDDL) and returns the plan under
// the 8.0 normalizer.
func aigPlan(t *testing.T, fromDDL, toDDL string) *MigrationPlan {
	t.Helper()
	return planFor(t, dbDDL+fromDDL, dbDDL+toDDL, NormalizerFor(MySQL80))
}

// opsContaining returns the plan ops whose SQL contains every fragment.
func opsContaining(p *MigrationPlan, fragments ...string) []MigrationOp {
	var out []MigrationOp
	for _, op := range p.Ops {
		all := true
		for _, f := range fragments {
			if !strings.Contains(op.SQL, f) {
				all = false
				break
			}
		}
		if all {
			out = append(out, op)
		}
	}
	return out
}

// requireOneStatementWith asserts exactly one op carries ALL fragments (i.e. they were grouped
// into one statement) and returns it.
func requireOneStatementWith(t *testing.T, p *MigrationPlan, fragments ...string) MigrationOp {
	t.Helper()
	got := opsContaining(p, fragments...)
	if len(got) != 1 {
		t.Fatalf("want exactly 1 statement containing %q, got %d:\n%s", fragments, len(got), p.SQL())
	}
	return got[0]
}

func TestAutoIncGrouping_AddColumnWithUniqueKey(t *testing.T) {
	p := aigPlan(t,
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT);",
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT, seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_seq (seq));")
	if len(p.Ops) != 1 {
		t.Fatalf("want 1 grouped op, got %d:\n%s", len(p.Ops), p.SQL())
	}
	op := p.Ops[0]
	for _, frag := range []string{"ADD COLUMN `seq`", "AUTO_INCREMENT", "ADD UNIQUE KEY `uk_seq` (`seq`)"} {
		if !strings.Contains(op.SQL, frag) {
			t.Errorf("grouped statement missing %q:\n%s", frag, op.SQL)
		}
	}
	if op.Type != OpAddColumn || op.Phase != PhaseMain || op.Priority != priorityColumn {
		t.Errorf("grouped op keeps the column op identity; got type=%s phase=%d prio=%d", op.Type, op.Phase, op.Priority)
	}
}

func TestAutoIncGrouping_AddColumnWithPlainKeyAndNewPK(t *testing.T) {
	// Plain KEY backing.
	p := aigPlan(t,
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT);",
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT, seq INT NOT NULL AUTO_INCREMENT, KEY k_seq (seq));")
	requireOneStatementWith(t, p, "ADD COLUMN `seq`", "ADD KEY `k_seq` (`seq`)")

	// New PRIMARY KEY backing on a PK-less table.
	p = aigPlan(t,
		"CREATE TABLE t (id INT NOT NULL, a INT);",
		"CREATE TABLE t (id INT NOT NULL, a INT, seq INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (seq));")
	requireOneStatementWith(t, p, "ADD COLUMN `seq`", "ADD PRIMARY KEY (`seq`)")
}

func TestAutoIncGrouping_PromotePKCombined(t *testing.T) {
	// The new AUTO_INCREMENT column replaces a natural PK: the column add joins the combined
	// DROP PRIMARY KEY / ADD PRIMARY KEY statement (three clauses, one statement).
	p := aigPlan(t,
		"CREATE TABLE t (a INT NOT NULL PRIMARY KEY, b INT);",
		"CREATE TABLE t (a INT NOT NULL, b INT, seq INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (seq));")
	requireOneStatementWith(t, p, "ADD COLUMN `seq`", "DROP PRIMARY KEY", "ADD PRIMARY KEY (`seq`)")
}

func TestAutoIncGrouping_ModifyGainAI(t *testing.T) {
	// Backing key arrives in the same plan: grouped.
	p := aigPlan(t,
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL);",
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_c (c));")
	requireOneStatementWith(t, p, "MODIFY COLUMN `c`", "AUTO_INCREMENT", "ADD UNIQUE KEY `uk_c` (`c`)")

	// CONTROL: column already keyed — single ungrouped MODIFY.
	p = aigPlan(t,
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL, UNIQUE KEY uk_c (c));",
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_c (c));")
	op := onlyOp(t, p)
	if !strings.Contains(op.SQL, "MODIFY COLUMN `c`") || strings.Contains(op.SQL, "ADD") {
		t.Errorf("already-keyed gain must stay a bare MODIFY: %s", op.SQL)
	}
}

func TestAutoIncGrouping_CompositeBackingPullsSecondColumn(t *testing.T) {
	p := aigPlan(t,
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY);",
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, seq INT NOT NULL AUTO_INCREMENT, d INT NOT NULL, UNIQUE KEY uk_sd (seq, d));")
	op := requireOneStatementWith(t, p, "ADD COLUMN `seq`", "ADD COLUMN `d`", "ADD UNIQUE KEY `uk_sd` (`seq`,`d`)")
	// The key clause must come after both column clauses it references.
	if strings.Index(op.SQL, "ADD UNIQUE KEY") < strings.Index(op.SQL, "ADD COLUMN `d`") {
		t.Errorf("key clause must follow the pulled column clause: %s", op.SQL)
	}
	if len(p.Ops) != 1 {
		t.Errorf("want 1 grouped op, got %d:\n%s", len(p.Ops), p.SQL())
	}
}

func TestAutoIncGrouping_AIMigratesBetweenColumns(t *testing.T) {
	// Both columns stay keyed; the de-AUTO_INCREMENT of `a` must share the statement that makes
	// `b` AUTO_INCREMENT (two statements would momentarily hold two auto columns — errno 1075).
	p := aigPlan(t,
		"CREATE TABLE t (a INT NOT NULL AUTO_INCREMENT, b INT NOT NULL, UNIQUE KEY uk_a (a), UNIQUE KEY uk_b (b));",
		"CREATE TABLE t (a INT NOT NULL, b INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_a (a), UNIQUE KEY uk_b (b));")
	op := requireOneStatementWith(t, p, "MODIFY COLUMN `a`", "MODIFY COLUMN `b`", "AUTO_INCREMENT")
	// The de-AI clause precedes the gain clause.
	if strings.Index(op.SQL, "MODIFY COLUMN `a`") > strings.Index(op.SQL, "MODIFY COLUMN `b`") {
		t.Errorf("de-AUTO_INCREMENT clause must precede the gaining clause: %s", op.SQL)
	}
	if len(p.Ops) != 1 {
		t.Errorf("want 1 grouped op, got %d:\n%s", len(p.Ops), p.SQL())
	}
}

func TestAutoIncGrouping_DropDirection(t *testing.T) {
	// Column + backing key dropped together: one statement, in PhasePre at the column drop slot.
	p := aigPlan(t,
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_c (c));",
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY);")
	op := requireOneStatementWith(t, p, "DROP INDEX `uk_c`", "DROP COLUMN `c`")
	if op.Phase != PhasePre {
		t.Errorf("grouped drop must stay in PhasePre, got %d", op.Phase)
	}

	// De-AUTO_INCREMENT + backing key drop: one statement.
	p = aigPlan(t,
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_c (c));",
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL);")
	requireOneStatementWith(t, p, "MODIFY COLUMN `c`", "DROP INDEX `uk_c`")

	// Backing key reshaped (same name), column stays AUTO_INCREMENT: DROP+ADD in one statement.
	p = aigPlan(t,
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, UNIQUE KEY uk_c (c));",
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, UNIQUE KEY uk_c (c, a));")
	requireOneStatementWith(t, p, "DROP INDEX `uk_c`", "ADD UNIQUE KEY `uk_c` (`c`,`a`)")

	// PRIMARY KEY backing dropped with the column.
	p = aigPlan(t,
		"CREATE TABLE t (c INT NOT NULL AUTO_INCREMENT PRIMARY KEY, a INT NOT NULL);",
		"CREATE TABLE t (a INT NOT NULL);")
	requireOneStatementWith(t, p, "DROP PRIMARY KEY", "DROP COLUMN `c`")

	// PRIMARY KEY dropped, column de-AUTO_INCREMENTed.
	p = aigPlan(t,
		"CREATE TABLE t (c INT NOT NULL AUTO_INCREMENT PRIMARY KEY, a INT);",
		"CREATE TABLE t (c INT NOT NULL, a INT);")
	requireOneStatementWith(t, p, "MODIFY COLUMN `c`", "DROP PRIMARY KEY")
}

func TestAutoIncGrouping_DropControls(t *testing.T) {
	// CONTROL: another key still covers the column — the drop stays a lone statement.
	p := aigPlan(t,
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk1 (c), KEY k2 (c));",
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, KEY k2 (c));")
	op := onlyOp(t, p)
	if op.SQL != "ALTER TABLE `app`.`t` DROP INDEX `uk1`" {
		t.Errorf("covered-elsewhere drop must stay ungrouped: %s", op.SQL)
	}

	// CONTROL: dropping a key on a non-AUTO_INCREMENT column is untouched.
	p = aigPlan(t,
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL, UNIQUE KEY uk_c (c));",
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL);")
	op = onlyOp(t, p)
	if op.SQL != "ALTER TABLE `app`.`t` DROP INDEX `uk_c`" {
		t.Errorf("non-AUTO_INCREMENT key drop must stay ungrouped: %s", op.SQL)
	}
}

func TestAutoIncGrouping_MyISAMNonLeftmostBacking(t *testing.T) {
	// MyISAM accepts a non-first key position as backing (oracle-probed), so the grouping fires
	// with the (a, seq) key.
	p := aigPlan(t,
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT NOT NULL) ENGINE=MyISAM;",
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, seq INT NOT NULL AUTO_INCREMENT, KEY k_as (a, seq)) ENGINE=MyISAM;")
	requireOneStatementWith(t, p, "ADD COLUMN `seq`", "ADD KEY `k_as` (`a`,`seq`)")

	// InnoDB requires the FIRST position (oracle-probed): a non-leftmost key is NOT a backing
	// candidate, so no grouping happens and the plan is left to fail naturally on an invalid
	// target (matching the CREATE path's stance on unkeyable AUTO_INCREMENT columns).
	p = aigPlan(t,
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT NOT NULL);",
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, seq INT NOT NULL AUTO_INCREMENT, KEY k_as (a, seq));")
	if got := opsContaining(p, "ADD COLUMN `seq`", "ADD KEY"); len(got) != 0 {
		t.Errorf("InnoDB non-leftmost key must not be grouped as backing:\n%s", p.SQL())
	}
}

func TestAutoIncGrouping_UnrelatedOpsUndisturbed(t *testing.T) {
	p := aigPlan(t,
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT, e VARCHAR(10)); CREATE TABLE u (x INT NOT NULL PRIMARY KEY, y INT);",
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT, e VARCHAR(20), f INT, seq INT NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_seq (seq)); CREATE TABLE u (x INT NOT NULL PRIMARY KEY, y BIGINT);")
	requireOneStatementWith(t, p, "ADD COLUMN `seq`", "ADD UNIQUE KEY `uk_seq`")
	// The unrelated t changes and the u change remain their own statements.
	requireOneStatementWith(t, p, "MODIFY COLUMN `e` varchar(20)")
	fOp := requireOneStatementWith(t, p, "ADD COLUMN `f`")
	if strings.Contains(fOp.SQL, "seq") {
		t.Errorf("unrelated column add must not be pulled into the group: %s", fOp.SQL)
	}
	requireOneStatementWith(t, p, "MODIFY COLUMN `y` bigint")
}

func TestAutoIncGrouping_Determinism(t *testing.T) {
	fromDDL := "CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, UNIQUE KEY uk1 (c), KEY k2 (c));"
	toDDL := "CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, UNIQUE KEY uk3 (c));"
	first := aigPlan(t, fromDDL, toDDL).SQL()
	for i := 0; i < 20; i++ {
		if got := aigPlan(t, fromDDL, toDDL).SQL(); got != first {
			t.Fatalf("plan not deterministic:\n%s\nvs\n%s", first, got)
		}
	}
	// All covering drops share the covering add's statement.
	p := aigPlan(t, fromDDL, toDDL)
	requireOneStatementWith(t, p, "DROP INDEX `uk1`", "DROP INDEX `k2`", "ADD UNIQUE KEY `uk3` (`c`)")
}

func TestAutoIncGrouping_PKRelocateToSecondary(t *testing.T) {
	// The PK moves off the auto column onto b while a NEW secondary key (name sorting AFTER
	// "primary") takes over the coverage: split DROP PRIMARY KEY, the covering ADD KEY, and the
	// replacement ADD PRIMARY KEY must share ONE statement — a standalone ADD PRIMARY KEY would
	// sort before it and run while the old PK still exists (errno 1068).
	p := aigPlan(t,
		"CREATE TABLE t (c INT NOT NULL AUTO_INCREMENT, b INT NOT NULL, x INT, PRIMARY KEY (c));",
		"CREATE TABLE t (c INT NOT NULL AUTO_INCREMENT, b INT NOT NULL, x INT, PRIMARY KEY (b), KEY zz_c (c));")
	op := requireOneStatementWith(t, p, "DROP PRIMARY KEY", "ADD KEY `zz_c` (`c`)", "ADD PRIMARY KEY (`b`)")
	if len(p.Ops) != 1 {
		t.Errorf("want 1 grouped op, got %d:\n%s", len(p.Ops), p.SQL())
	}
	if strings.Count(op.SQL, "PRIMARY KEY") != 2 {
		t.Errorf("statement must carry exactly the drop and one re-add: %s", op.SQL)
	}
}

func TestAutoIncGrouping_DemotedMemberJoinsPKStatement(t *testing.T) {
	// The auto column keeps the shrinking PK while another old-PK member is demoted to
	// nullable: the demoting MODIFY joins the grouped DROP+ADD PRIMARY KEY statement (errno
	// 1171 if it ran earlier, errno 1075 if the PK drop ran alone).
	p := aigPlan(t,
		"CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, PRIMARY KEY (id, a));",
		"CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, a INT NULL, PRIMARY KEY (id));")
	requireOneStatementWith(t, p, "DROP PRIMARY KEY", "ADD PRIMARY KEY (`id`)", "MODIFY COLUMN `a`")
	if len(p.Ops) != 1 {
		t.Errorf("want 1 grouped op, got %d:\n%s", len(p.Ops), p.SQL())
	}
}

func TestAutoIncGrouping_GeneratedKeyPartPullsPlainAdds(t *testing.T) {
	// A generated key part g rides in the backing key and references a new plain column d that
	// is NOT itself a key part: d's ADD must be pulled into the grouped statement before g's.
	p := aigPlan(t,
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY);",
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, seq INT NOT NULL AUTO_INCREMENT, d INT NOT NULL, g INT GENERATED ALWAYS AS (d+1) VIRTUAL, UNIQUE KEY uk_sg (seq, g));")
	op := requireOneStatementWith(t, p, "ADD COLUMN `seq`", "ADD COLUMN `d`", "ADD COLUMN `g`", "ADD UNIQUE KEY `uk_sg`")
	if strings.Index(op.SQL, "ADD COLUMN `d`") > strings.Index(op.SQL, "ADD COLUMN `g`") {
		t.Errorf("plain add must precede the generated add that references it: %s", op.SQL)
	}
	if len(p.Ops) != 1 {
		t.Errorf("want 1 grouped op, got %d:\n%s", len(p.Ops), p.SQL())
	}
}

func TestAutoIncGrouping_FKBackingDropJoinsColumnDrop(t *testing.T) {
	// The AUTO_INCREMENT column's only covering key is the FK-implicit backing index and both
	// the FK and the column leave: the FK node's leftover DROP INDEX (PhasePre, before column
	// drops) would strand the still-AUTO_INCREMENT column, so its clause must join the DROP
	// COLUMN statement — and it must NOT vanish from the plan when consumed.
	p := aigPlan(t,
		"CREATE TABLE p (id INT NOT NULL PRIMARY KEY); CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT NOT NULL AUTO_INCREMENT, CONSTRAINT fk_c FOREIGN KEY (c) REFERENCES p (id));",
		"CREATE TABLE p (id INT NOT NULL PRIMARY KEY); CREATE TABLE t (id INT NOT NULL PRIMARY KEY);")
	op := requireOneStatementWith(t, p, "DROP INDEX `fk_c`", "DROP COLUMN `c`")
	if op.Phase != PhasePre {
		t.Errorf("grouped drop must stay in PhasePre, got %d", op.Phase)
	}
	// The FK constraint drop stays its own earlier statement.
	requireOneStatementWith(t, p, "DROP FOREIGN KEY `fk_c`")
	// No DROP INDEX clause may be lost: exactly one occurrence across the whole plan.
	if got := strings.Count(p.SQL(), "DROP INDEX `fk_c`"); got != 1 {
		t.Errorf("want the backing DROP INDEX exactly once in the plan, got %d:\n%s", got, p.SQL())
	}
}

func TestAutoIncGrouping_MyISAMSharedKeyAIMigration(t *testing.T) {
	// MyISAM corner: the composite key covers BOTH the old and the new auto column (any-position
	// rule). The old column and the shared key drop in PhasePre, so the gaining MODIFY must
	// still be grouped with the NEW backing key — a consumed covering drop that lands in a
	// PhasePre statement does not keep covering PhaseMain.
	p := aigPlan(t,
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c1 INT NOT NULL AUTO_INCREMENT, c2 INT NOT NULL, KEY k (c1, c2)) ENGINE=MyISAM;",
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c2 INT NOT NULL AUTO_INCREMENT, KEY k2 (c2)) ENGINE=MyISAM;")
	requireOneStatementWith(t, p, "DROP INDEX `k`", "DROP COLUMN `c1`")
	requireOneStatementWith(t, p, "MODIFY COLUMN `c2`", "ADD KEY `k2` (`c2`)")
}

func TestAutoIncGrouping_IdempotenceHermetic(t *testing.T) {
	// Self-diff of every grouped target form stays an empty plan under both normalizers.
	schemas := []string{
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, UNIQUE KEY uk_seq (seq));",
		"CREATE TABLE t (a INT NOT NULL, seq INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (seq));",
		"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, seq INT NOT NULL AUTO_INCREMENT, KEY k_as (a, seq)) ENGINE=MyISAM;",
	}
	for _, s := range schemas {
		c := loadCat(t, dbDDL+s)
		for _, n := range []*Normalizer{NormalizerFor(MySQL80), NormalizerFor(MySQL57)} {
			diff := DiffWithNormalizer(c, c, n)
			if plan := GenerateMigrationWithNormalizer(c, c, diff, n); plan.SQL() != "" {
				t.Errorf("self-diff plan not empty (version=%d) for %q:\n%s", n.Version, s, plan.SQL())
			}
		}
	}
}
