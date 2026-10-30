package catalog

import (
	"strings"
	"testing"
)

// In-process, hermetic tests for generate-core. They assert MigrationPlan shape, SQL()
// joining, Filter, the idempotence terminal (empty diff → empty plan), deterministic op
// ordering, and the structural content of the table/column DDL — loading both sides through
// the omni catalog (LoadSDL) so the model state is authentic. The oracle-backed
// apply-correctness + idempotence proofs (correctness-protocol.md gates 1 & 2) live in
// migration_oracle_test.go and run against the live engines.

// planFor diffs two SDL schemas and generates a migration plan with the given normalizer.
func planFor(t *testing.T, fromSDL, toSDL string, n *Normalizer) *MigrationPlan {
	t.Helper()
	from := loadCat(t, fromSDL)
	to := loadCat(t, toSDL)
	diff := DiffWithNormalizer(from, to, n)
	return GenerateMigrationWithNormalizer(from, to, diff, n)
}

// onlyOp returns the single op in a plan, failing if the count differs.
func onlyOp(t *testing.T, p *MigrationPlan) MigrationOp {
	t.Helper()
	if len(p.Ops) != 1 {
		t.Fatalf("want exactly 1 op, got %d: %s", len(p.Ops), p.SQL())
	}
	return p.Ops[0]
}

// ---- idempotence terminal: empty diff → empty plan -------------------------

func TestMigration_EmptyDiffEmptyPlan(t *testing.T) {
	schemas := []string{
		dbDDL + "CREATE TABLE t (id INT NOT NULL, name VARCHAR(50), PRIMARY KEY (id));",
		dbDDL + "CREATE TABLE g (a INT, b INT GENERATED ALWAYS AS (a+1) STORED, c VARCHAR(10) DEFAULT 'x');",
		dbDDL + "CREATE TABLE ts (id INT PRIMARY KEY, created TIMESTAMP DEFAULT CURRENT_TIMESTAMP, n INT DEFAULT 0);",
		dbDDL + "CREATE TABLE a (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY, v DECIMAL(10,2) DEFAULT 0) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;",
	}
	for _, s := range schemas {
		c := loadCat(t, s)
		for _, n := range []*Normalizer{NormalizerFor(MySQL80), NormalizerFor(MySQL57)} {
			diff := DiffWithNormalizer(c, c, n)
			plan := GenerateMigrationWithNormalizer(c, c, diff, n)
			if len(plan.Ops) != 0 {
				t.Errorf("self-diff plan not empty (version=%d) for %q:\n%s", n.Version, s, plan.SQL())
			}
			if plan.SQL() != "" {
				t.Errorf("self-diff SQL not empty (version=%d) for %q: %q", n.Version, s, plan.SQL())
			}
		}
		// Parameterless GenerateMigration must also be empty.
		if plan := GenerateMigration(c, c, Diff(c, c)); plan.SQL() != "" {
			t.Errorf("parameterless self-diff SQL not empty for %q: %q", s, plan.SQL())
		}
	}
}

func TestMigration_NilDiffEmptyPlan(t *testing.T) {
	if p := GenerateMigration(New(), New(), nil); p == nil || len(p.Ops) != 0 || p.SQL() != "" {
		t.Errorf("nil diff must produce empty plan, got %+v", p)
	}
}

// ---- MigrationPlan helpers -------------------------------------------------

func TestMigrationPlan_SQLJoin(t *testing.T) {
	p := &MigrationPlan{Ops: []MigrationOp{
		{SQL: "DROP TABLE `a`"},
		{SQL: "CREATE TABLE `b` (\n  `id` int\n)"},
	}}
	got := p.SQL()
	if !strings.Contains(got, ";\n") {
		t.Errorf("SQL() must join ops with \";\\n\", got %q", got)
	}
	if strings.Count(got, ";\n") != 1 {
		t.Errorf("two ops must have exactly one separator, got %q", got)
	}
	if (&MigrationPlan{}).SQL() != "" {
		t.Errorf("empty plan SQL() must be \"\"")
	}
}

func TestMigrationPlan_Filter(t *testing.T) {
	p := &MigrationPlan{Ops: []MigrationOp{
		{Type: OpCreateTable, ObjectName: "keep"},
		{Type: OpDropTable, ObjectName: "drop"},
	}}
	kept := p.Filter(func(op MigrationOp) bool { return op.Type != OpDropTable })
	if len(kept.Ops) != 1 || kept.Ops[0].ObjectName != "keep" {
		t.Errorf("Filter did not keep the right op: %+v", kept.Ops)
	}
}

// ---- table-level DDL --------------------------------------------------------

func TestMigration_CreateTable(t *testing.T) {
	from := dbDDL
	to := dbDDL + "CREATE TABLE t (id INT NOT NULL PRIMARY KEY, name VARCHAR(20));"
	op := onlyOp(t, planFor(t, from, to, NormalizerFor(MySQL80)))
	if op.Type != OpCreateTable {
		t.Fatalf("want CreateTable, got %s", op.Type)
	}
	sql := op.SQL
	for _, want := range []string{"CREATE TABLE `app`.`t`", "`id` int", "`name` varchar(20)", "PRIMARY KEY (`id`)", "ENGINE=InnoDB"} {
		if !strings.Contains(sql, want) {
			t.Errorf("CREATE TABLE missing %q:\n%s", want, sql)
		}
	}
	// AUTO_INCREMENT=N must never appear (ignore-in-diff).
	if strings.Contains(sql, "AUTO_INCREMENT=") {
		t.Errorf("CREATE TABLE must not emit AUTO_INCREMENT=N counter:\n%s", sql)
	}
}

func TestMigration_DropTable(t *testing.T) {
	from := dbDDL + "CREATE TABLE a (id INT PRIMARY KEY); CREATE TABLE b (id INT PRIMARY KEY);"
	to := dbDDL + "CREATE TABLE a (id INT PRIMARY KEY);"
	op := onlyOp(t, planFor(t, from, to, NormalizerFor(MySQL80)))
	if op.Type != OpDropTable {
		t.Fatalf("want DropTable, got %s", op.Type)
	}
	if !strings.Contains(op.SQL, "DROP TABLE `app`.`b`") {
		t.Errorf("want DROP TABLE `b`, got %q", op.SQL)
	}
	if op.Phase != PhasePre {
		t.Errorf("DROP must be PhasePre, got %d", op.Phase)
	}
}

func TestMigration_TableEngineChange(t *testing.T) {
	from := dbDDL + "CREATE TABLE t (id INT PRIMARY KEY) ENGINE=InnoDB;"
	to := dbDDL + "CREATE TABLE t (id INT PRIMARY KEY) ENGINE=MyISAM;"
	op := onlyOp(t, planFor(t, from, to, NormalizerFor(MySQL80)))
	if op.Type != OpAlterTable {
		t.Fatalf("want AlterTable, got %s", op.Type)
	}
	if !strings.Contains(strings.ToUpper(op.SQL), "ENGINE=MYISAM") {
		t.Errorf("want ENGINE=MyISAM, got %q", op.SQL)
	}
}

func TestMigration_TableCommentChange(t *testing.T) {
	from := dbDDL + "CREATE TABLE t (id INT PRIMARY KEY) COMMENT='old';"
	to := dbDDL + "CREATE TABLE t (id INT PRIMARY KEY) COMMENT='new';"
	op := onlyOp(t, planFor(t, from, to, NormalizerFor(MySQL80)))
	if op.Type != OpAlterTable || !strings.Contains(op.SQL, "COMMENT='new'") {
		t.Errorf("want ALTER TABLE COMMENT='new', got %q", op.SQL)
	}
}

func TestMigration_TableCharsetChange(t *testing.T) {
	from := dbDDL + "CREATE TABLE t (id INT PRIMARY KEY) DEFAULT CHARSET=utf8mb4;"
	to := dbDDL + "CREATE TABLE t (id INT PRIMARY KEY) DEFAULT CHARSET=latin1;"
	op := onlyOp(t, planFor(t, from, to, NormalizerFor(MySQL80)))
	if op.Type != OpAlterTable || !strings.Contains(strings.ToUpper(op.SQL), "CHARSET=LATIN1") {
		t.Errorf("want ALTER TABLE DEFAULT CHARSET=latin1, got %q", op.SQL)
	}
}

// A column-only modification must NOT produce a table-option ALTER (the two generators must
// not double-emit the table). Only the column MODIFY should appear.
func TestMigration_ColumnChangeNoTableOptionAlter(t *testing.T) {
	from := dbDDL + "CREATE TABLE t (id INT PRIMARY KEY, v INT);"
	to := dbDDL + "CREATE TABLE t (id INT PRIMARY KEY, v BIGINT);"
	p := planFor(t, from, to, NormalizerFor(MySQL80))
	for _, op := range p.Ops {
		if op.Type == OpAlterTable {
			t.Errorf("column-only change must not emit a table-option ALTER, got %q", op.SQL)
		}
	}
	op := onlyOp(t, p)
	if op.Type != OpModifyColumn {
		t.Fatalf("want a single MODIFY COLUMN, got %s: %s", op.Type, op.SQL)
	}
}

// ---- column-level DDL -------------------------------------------------------

func TestMigration_AddColumn(t *testing.T) {
	from := dbDDL + "CREATE TABLE t (id INT PRIMARY KEY);"
	to := dbDDL + "CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(20) NOT NULL DEFAULT 'x');"
	op := onlyOp(t, planFor(t, from, to, NormalizerFor(MySQL80)))
	if op.Type != OpAddColumn {
		t.Fatalf("want AddColumn, got %s", op.Type)
	}
	for _, want := range []string{"ALTER TABLE `app`.`t` ADD COLUMN", "`name` varchar(20)", "NOT NULL", "DEFAULT 'x'"} {
		if !strings.Contains(op.SQL, want) {
			t.Errorf("ADD COLUMN missing %q:\n%s", want, op.SQL)
		}
	}
	if op.ParentObject != "t" {
		t.Errorf("column op ParentObject must be the table, got %q", op.ParentObject)
	}
}

func TestMigration_DropColumn(t *testing.T) {
	from := dbDDL + "CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(20));"
	to := dbDDL + "CREATE TABLE t (id INT PRIMARY KEY);"
	op := onlyOp(t, planFor(t, from, to, NormalizerFor(MySQL80)))
	if op.Type != OpDropColumn || !strings.Contains(op.SQL, "DROP COLUMN `name`") {
		t.Errorf("want DROP COLUMN `name`, got %q", op.SQL)
	}
	if op.Phase != PhasePre {
		t.Errorf("DROP COLUMN must be PhasePre, got %d", op.Phase)
	}
}

func TestMigration_ModifyColumnUsesModify(t *testing.T) {
	from := dbDDL + "CREATE TABLE t (id INT PRIMARY KEY, v INT);"
	to := dbDDL + "CREATE TABLE t (id INT PRIMARY KEY, v BIGINT NOT NULL);"
	op := onlyOp(t, planFor(t, from, to, NormalizerFor(MySQL80)))
	if op.Type != OpModifyColumn {
		t.Fatalf("want ModifyColumn, got %s", op.Type)
	}
	// MySQL uses MODIFY COLUMN, never PG's ALTER COLUMN ... TYPE.
	if !strings.Contains(op.SQL, "MODIFY COLUMN `v` bigint") {
		t.Errorf("want MODIFY COLUMN `v` bigint, got %q", op.SQL)
	}
	if strings.Contains(op.SQL, "ALTER COLUMN") {
		t.Errorf("MySQL must use MODIFY, not ALTER COLUMN: %q", op.SQL)
	}
}

// 5.7 renders integer display widths; 8.0 drops them. The MODIFY must reflect the target
// version's canonical type (proving rendering routes through the Normalizer, not a fixed
// 8.0 form).
func TestMigration_ModifyColumnVersionedWidth(t *testing.T) {
	from := dbDDL + "CREATE TABLE t (id INT PRIMARY KEY, v SMALLINT);"
	to := dbDDL + "CREATE TABLE t (id INT PRIMARY KEY, v BIGINT);"
	op80 := onlyOp(t, planFor(t, from, to, NormalizerFor(MySQL80)))
	if !strings.Contains(op80.SQL, "`v` bigint") || strings.Contains(op80.SQL, "bigint(") {
		t.Errorf("8.0 MODIFY must render width-less bigint, got %q", op80.SQL)
	}
	op57 := onlyOp(t, planFor(t, from, to, NormalizerFor(MySQL57)))
	if !strings.Contains(op57.SQL, "`v` bigint(20)") {
		t.Errorf("5.7 MODIFY must render bigint(20), got %q", op57.SQL)
	}
}

func TestMigration_AddGeneratedColumn(t *testing.T) {
	from := dbDDL + "CREATE TABLE t (id INT PRIMARY KEY, a INT);"
	to := dbDDL + "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT GENERATED ALWAYS AS (a+1) STORED);"
	op := onlyOp(t, planFor(t, from, to, NormalizerFor(MySQL80)))
	if op.Type != OpAddColumn {
		t.Fatalf("want AddColumn, got %s", op.Type)
	}
	if !strings.Contains(op.SQL, "GENERATED ALWAYS AS") || !strings.Contains(op.SQL, "STORED") {
		t.Errorf("generated column must render GENERATED ALWAYS AS (..) STORED, got %q", op.SQL)
	}
	// A generated column must not carry a DEFAULT.
	if strings.Contains(op.SQL, "DEFAULT") {
		t.Errorf("generated column must not have DEFAULT, got %q", op.SQL)
	}
}

func TestMigration_AddAutoIncrementColumn(t *testing.T) {
	// Adding an AUTO_INCREMENT column to a table that has it as PK already.
	from := dbDDL + "CREATE TABLE t (id INT NOT NULL PRIMARY KEY);"
	to := dbDDL + "CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY);"
	op := onlyOp(t, planFor(t, from, to, NormalizerFor(MySQL80)))
	if op.Type != OpModifyColumn {
		t.Fatalf("want ModifyColumn, got %s", op.Type)
	}
	if !strings.Contains(op.SQL, "AUTO_INCREMENT") {
		t.Errorf("want AUTO_INCREMENT attribute in MODIFY, got %q", op.SQL)
	}
}

// ---- deterministic ordering -------------------------------------------------

// Drops must precede creates/modifies, and the whole plan must be byte-stable across runs
// (idempotence depends on it). We construct a multi-change diff and assert ordering + that
// two generations of the same diff are identical.
func TestMigration_DeterministicOrdering(t *testing.T) {
	from := dbDDL + "CREATE TABLE keep (id INT PRIMARY KEY, old_col INT); CREATE TABLE goner (id INT PRIMARY KEY); CREATE TABLE zeta (id INT PRIMARY KEY);"
	to := dbDDL + "CREATE TABLE keep (id INT PRIMARY KEY, new_col INT); CREATE TABLE zeta (id INT PRIMARY KEY); CREATE TABLE alpha (id INT PRIMARY KEY);"

	n := NormalizerFor(MySQL80)
	fromCat := loadCat(t, from)
	toCat := loadCat(t, to)
	diff := DiffWithNormalizer(fromCat, toCat, n)

	p1 := GenerateMigrationWithNormalizer(fromCat, toCat, diff, n)
	p2 := GenerateMigrationWithNormalizer(fromCat, toCat, diff, n)
	if p1.SQL() != p2.SQL() {
		t.Fatalf("plan generation not deterministic:\n--- p1 ---\n%s\n--- p2 ---\n%s", p1.SQL(), p2.SQL())
	}

	// All PhasePre (drops) ops must come before any PhaseMain op.
	seenMain := false
	for _, op := range p1.Ops {
		switch op.Phase {
		case PhasePre:
			if seenMain {
				t.Errorf("a drop (%s %s) appears after a create/alter — phase ordering broken:\n%s",
					op.Type, op.ObjectName, p1.SQL())
			}
		case PhaseMain:
			seenMain = true
		}
	}

	// Sanity: the expected op kinds are present (drop table goner, drop column old_col, create
	// table alpha, add column new_col).
	kinds := map[MigrationOpType]int{}
	for _, op := range p1.Ops {
		kinds[op.Type]++
	}
	if kinds[OpDropTable] != 1 || kinds[OpCreateTable] != 1 || kinds[OpAddColumn] != 1 || kinds[OpDropColumn] != 1 {
		t.Errorf("unexpected op mix: %+v\n%s", kinds, p1.SQL())
	}
}

// ---- review regressions ----------------------------------------------------

// Blocker #4: the rendered table identifier must preserve the DECLARED casing, not the
// lower-cased diff identity key (on a case-sensitive server the stored name is significant).
func TestMigration_TableNameCasePreserved(t *testing.T) {
	from := "CREATE DATABASE App DEFAULT CHARSET=utf8mb4; USE App; CREATE TABLE MyTable (Id INT PRIMARY KEY); CREATE TABLE Goner (Id INT PRIMARY KEY);"
	to := "CREATE DATABASE App DEFAULT CHARSET=utf8mb4; USE App; CREATE TABLE MyTable (Id INT PRIMARY KEY, NewCol INT);"
	p := planFor(t, from, to, NormalizerFor(MySQL80))
	sql := p.SQL()
	if !strings.Contains(sql, "`App`.`MyTable`") {
		t.Errorf("table identifier must preserve case `App`.`MyTable`, got:\n%s", sql)
	}
	if !strings.Contains(sql, "`App`.`Goner`") {
		t.Errorf("drop must preserve case `App`.`Goner`, got:\n%s", sql)
	}
	if strings.Contains(sql, "`mytable`") || strings.Contains(sql, "`app`.`goner`") {
		t.Errorf("lower-cased identity key leaked into DDL:\n%s", sql)
	}
}

// Blocker #1: when both a generated column and a plain column it references are dropped, the
// generated column must be dropped FIRST (MySQL error 3108 otherwise).
func TestMigration_GeneratedColumnDroppedBeforeBase(t *testing.T) {
	from := dbDDL + "CREATE TABLE t (id INT PRIMARY KEY, a INT, z INT GENERATED ALWAYS AS (a+1) STORED);"
	to := dbDDL + "CREATE TABLE t (id INT PRIMARY KEY);"
	p := planFor(t, from, to, NormalizerFor(MySQL80))
	posA, posZ := -1, -1
	for i, op := range p.Ops {
		if op.Type == OpDropColumn && strings.Contains(op.SQL, "`a`") {
			posA = i
		}
		if op.Type == OpDropColumn && strings.Contains(op.SQL, "`z`") {
			posZ = i
		}
	}
	if posA < 0 || posZ < 0 {
		t.Fatalf("expected both column drops, got:\n%s", p.SQL())
	}
	if posZ > posA {
		t.Errorf("generated column `z` must be dropped before base column `a`:\n%s", p.SQL())
	}
}

// Blocker #2: a generated-column storage-mode flip (VIRTUAL↔STORED) must be rendered as
// DROP + ADD, never an in-place MODIFY (MySQL error 3106).
func TestMigration_GeneratedStorageFlipIsDropAdd(t *testing.T) {
	from := dbDDL + "CREATE TABLE t (id INT PRIMARY KEY, a INT, g INT GENERATED ALWAYS AS (a+1) VIRTUAL);"
	to := dbDDL + "CREATE TABLE t (id INT PRIMARY KEY, a INT, g INT GENERATED ALWAYS AS (a+1) STORED);"
	p := planFor(t, from, to, NormalizerFor(MySQL80))
	var sawDrop, sawAdd, sawModify bool
	for _, op := range p.Ops {
		switch op.Type {
		case OpDropColumn:
			if strings.Contains(op.SQL, "`g`") {
				sawDrop = true
			}
		case OpAddColumn:
			if strings.Contains(op.SQL, "`g`") {
				sawAdd = true
			}
		case OpModifyColumn:
			if strings.Contains(op.SQL, "`g`") {
				sawModify = true
			}
		}
	}
	if sawModify {
		t.Errorf("storage-mode flip must not use MODIFY:\n%s", p.SQL())
	}
	if !sawDrop || !sawAdd {
		t.Errorf("storage-mode flip must emit DROP + ADD, got drop=%v add=%v:\n%s", sawDrop, sawAdd, p.SQL())
	}
	// The re-added column must carry the NEW storage mode (STORED) and its generation expr.
	for _, op := range p.Ops {
		if op.Type == OpAddColumn && strings.Contains(op.SQL, "`g`") {
			if !strings.Contains(op.SQL, "STORED") || !strings.Contains(op.SQL, "GENERATED ALWAYS AS") {
				t.Errorf("re-added generated column must render the new STORED generation clause:\n%s", op.SQL)
			}
		}
	}
}

// A generation EXPRESSION change with the SAME storage mode is a normal MODIFY (not DROP+ADD).
func TestMigration_GeneratedExprChangeIsModify(t *testing.T) {
	from := dbDDL + "CREATE TABLE t (id INT PRIMARY KEY, a INT, g INT GENERATED ALWAYS AS (a+1) STORED);"
	to := dbDDL + "CREATE TABLE t (id INT PRIMARY KEY, a INT, g INT GENERATED ALWAYS AS (a+2) STORED);"
	op := onlyOp(t, planFor(t, from, to, NormalizerFor(MySQL80)))
	if op.Type != OpModifyColumn {
		t.Fatalf("generation-expr change (same storage) must be MODIFY, got %s:\n%s", op.Type, op.SQL)
	}
}

// Review re-round N2: when a generated column is re-added (storage flip) in the same plan as a
// MODIFY of a base column it references, the generated ADD must come AFTER the base MODIFY — even
// when the generated name sorts lexicographically before the base name (g < z). The generated
// column's higher add-priority enforces this.
func TestMigration_GeneratedAddAfterBaseModify(t *testing.T) {
	// `g` (generated, refs z) sorts before `z`; flip g VIRTUAL->STORED and widen z's type.
	from := dbDDL + "CREATE TABLE t (id INT PRIMARY KEY, g INT GENERATED ALWAYS AS (z+1) VIRTUAL, z SMALLINT);"
	to := dbDDL + "CREATE TABLE t (id INT PRIMARY KEY, g INT GENERATED ALWAYS AS (z+1) STORED, z BIGINT);"
	p := planFor(t, from, to, NormalizerFor(MySQL80))
	posGAdd, posZMod := -1, -1
	for i, op := range p.Ops {
		if op.Type == OpAddColumn && strings.Contains(op.SQL, "`g`") {
			posGAdd = i
		}
		if op.Type == OpModifyColumn && strings.Contains(op.SQL, "`z`") {
			posZMod = i
		}
	}
	if posGAdd < 0 || posZMod < 0 {
		t.Fatalf("expected generated `g` ADD and base `z` MODIFY, got:\n%s", p.SQL())
	}
	if posGAdd < posZMod {
		t.Errorf("generated `g` ADD must come after base `z` MODIFY:\n%s", p.SQL())
	}
}

// sortMigrationOps must be a pure, stable function of its input ordering for equal keys.
func TestSortMigrationOps_StableAndPhased(t *testing.T) {
	ops := []MigrationOp{
		{Type: OpAddColumn, Phase: PhaseMain, Priority: priorityColumn, sortName: "d.t.b", SQL: "B"},
		{Type: OpDropTable, Phase: PhasePre, Priority: priorityTable, sortName: "d.z", SQL: "DROPZ"},
		{Type: OpAddColumn, Phase: PhaseMain, Priority: priorityColumn, sortName: "d.t.a", SQL: "A"},
		{Type: OpCreateTable, Phase: PhaseMain, Priority: priorityTable, sortName: "d.t", SQL: "CREATE"},
	}
	got := sortMigrationOps(ops)
	gotOrder := make([]string, len(got))
	for i, op := range got {
		gotOrder[i] = op.SQL
	}
	want := []string{"DROPZ", "CREATE", "A", "B"} // pre, then table-priority, then column a<b
	if strings.Join(gotOrder, ",") != strings.Join(want, ",") {
		t.Errorf("sortMigrationOps order = %v, want %v", gotOrder, want)
	}
}
