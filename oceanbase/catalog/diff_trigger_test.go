package catalog

import (
	"strings"
	"testing"
)

// Hermetic unit tests for the trigger differ + generator (no oracle). They pin the diff
// semantics — identity, the compared-field set, the deliberate Definer/Order ignores — and the
// rendered DDL shape. The live-engine idempotence + apply-correctness proofs are in
// migration_trigger_oracle_test.go.

// loadTriggerCatalog loads an SDL snippet (wrapped in a utf8mb4 database) and returns the catalog.
func loadTriggerCatalog(t *testing.T, body string) *Catalog {
	t.Helper()
	cat, err := LoadSDL("CREATE DATABASE d DEFAULT CHARSET=utf8mb4;\nUSE d;\n" + body)
	if err != nil {
		t.Fatalf("LoadSDL failed for %q: %v", body, err)
	}
	return cat
}

// triggerEntry returns the single trigger diff entry for name, or fails.
func triggerEntry(t *testing.T, d *SchemaDiff, name string) *TriggerDiffEntry {
	t.Helper()
	for i := range d.Triggers {
		if strings.EqualFold(d.Triggers[i].Name, name) {
			return &d.Triggers[i]
		}
	}
	t.Fatalf("no trigger diff entry for %q (have %d entries)", name, len(d.Triggers))
	return nil
}

func TestTriggerDiff_SelfEmpty(t *testing.T) {
	n := NormalizerFor(MySQL80)
	cat := loadTriggerCatalog(t, `
CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT);
CREATE TRIGGER bi BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1;
CREATE TRIGGER au AFTER UPDATE ON t FOR EACH ROW SET @x = NEW.b;`)
	if d := DiffWithNormalizer(cat, cat, n); !d.IsEmpty() {
		t.Fatalf("self-diff not empty: %s", describeTriggerDiff(d))
	}
}

func TestTriggerDiff_Add(t *testing.T) {
	n := NormalizerFor(MySQL80)
	from := loadTriggerCatalog(t, `CREATE TABLE t (id INT PRIMARY KEY, a INT);`)
	to := loadTriggerCatalog(t, `
CREATE TABLE t (id INT PRIMARY KEY, a INT);
CREATE TRIGGER bi BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1;`)
	d := DiffWithNormalizer(from, to, n)
	e := triggerEntry(t, d, "bi")
	if e.Action != DiffAdd || e.To == nil || e.From != nil {
		t.Fatalf("want DiffAdd with To set, got %v", e.Action)
	}
}

func TestTriggerDiff_Drop(t *testing.T) {
	n := NormalizerFor(MySQL80)
	from := loadTriggerCatalog(t, `
CREATE TABLE t (id INT PRIMARY KEY, a INT);
CREATE TRIGGER bi BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1;`)
	to := loadTriggerCatalog(t, `CREATE TABLE t (id INT PRIMARY KEY, a INT);`)
	d := DiffWithNormalizer(from, to, n)
	e := triggerEntry(t, d, "bi")
	if e.Action != DiffDrop || e.From == nil || e.To != nil {
		t.Fatalf("want DiffDrop with From set, got %v", e.Action)
	}
}

// TestTriggerDiff_ModifyEachField proves every compared field independently triggers a MODIFY.
func TestTriggerDiff_ModifyEachField(t *testing.T) {
	n := NormalizerFor(MySQL80)
	base := "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT);\nCREATE TABLE u (id INT PRIMARY KEY, a INT);\n"
	cases := []struct {
		name string
		from string
		to   string
	}{
		{"timing", "CREATE TRIGGER x BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1;",
			"CREATE TRIGGER x AFTER INSERT ON t FOR EACH ROW SET @z = NEW.a;"},
		{"event", "CREATE TRIGGER x BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1;",
			"CREATE TRIGGER x BEFORE UPDATE ON t FOR EACH ROW SET NEW.a = 1;"},
		{"table", "CREATE TRIGGER x BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1;",
			"CREATE TRIGGER x BEFORE INSERT ON u FOR EACH ROW SET NEW.a = 1;"},
		{"body", "CREATE TRIGGER x BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1;",
			"CREATE TRIGGER x BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 2;"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			from := loadTriggerCatalog(t, base+tc.from)
			to := loadTriggerCatalog(t, base+tc.to)
			d := DiffWithNormalizer(from, to, n)
			e := triggerEntry(t, d, "x")
			if e.Action != DiffModify {
				t.Fatalf("want DiffModify for changed %s, got %v", tc.name, e.Action)
			}
		})
	}
}

// TestTriggerDiff_IgnoreDefiner proves a DEFINER difference does NOT produce a diff: the user form
// (no DEFINER → loader default `root`@`%`) and a stored form carrying an explicit DEFINER must
// collapse to empty.
func TestTriggerDiff_IgnoreDefiner(t *testing.T) {
	n := NormalizerFor(MySQL80)
	userForm := loadTriggerCatalog(t, `
CREATE TABLE t (id INT PRIMARY KEY, a INT);
CREATE TRIGGER bi BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1;`)
	storedForm := loadTriggerCatalog(t, "CREATE TABLE t (id INT PRIMARY KEY, a INT);\n"+
		"CREATE DEFINER=`someuser`@`localhost` TRIGGER bi BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1;")
	if d := DiffWithNormalizer(userForm, storedForm, n); !d.IsEmpty() {
		t.Fatalf("DEFINER difference phantom-diffed (should be ignored): %s", describeTriggerDiff(d))
	}
	if d := DiffWithNormalizer(storedForm, userForm, n); !d.IsEmpty() {
		t.Fatalf("DEFINER difference phantom-diffed (reverse): %s", describeTriggerDiff(d))
	}
}

// TestTriggerDiff_IgnoreOrder proves FOLLOWS/PRECEDES is NOT diffed: a trigger declared with
// FOLLOWS (non-nil Order) and the same trigger declared without it (the readback form, nil Order)
// must collapse to empty. This is the deliberate trigger-order-not-modelled flag.
func TestTriggerDiff_IgnoreOrder(t *testing.T) {
	n := NormalizerFor(MySQL80)
	withFollows := loadTriggerCatalog(t, `
CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT);
CREATE TRIGGER bi1 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1;
CREATE TRIGGER bi2 BEFORE INSERT ON t FOR EACH ROW FOLLOWS bi1 SET NEW.b = 2;`)
	// Sanity: the loader actually captured the FOLLOWS so the test is meaningful.
	if got := withFollows.GetDatabase("d").Triggers["bi2"]; got == nil || got.Order == nil {
		t.Fatalf("precondition: expected bi2.Order non-nil (FOLLOWS captured)")
	}
	noFollows := loadTriggerCatalog(t, `
CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT);
CREATE TRIGGER bi1 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1;
CREATE TRIGGER bi2 BEFORE INSERT ON t FOR EACH ROW SET NEW.b = 2;`)
	if d := DiffWithNormalizer(withFollows, noFollows, n); !d.IsEmpty() {
		t.Fatalf("FOLLOWS/PRECEDES phantom-diffed (should be ignored): %s", describeTriggerDiff(d))
	}
}

// TestTriggerDiff_BodyCaseWhitespaceTrim proves leading/trailing whitespace on the body is folded
// (the loader TrimSpaces). Interior whitespace/case is significant (MySQL stores it verbatim) and
// is NOT folded — a genuine body change must still diff.
func TestTriggerDiff_BodyWhitespaceTrim(t *testing.T) {
	n := NormalizerFor(MySQL80)
	a := &Trigger{Name: "x", Table: "t", Timing: "BEFORE", Event: "INSERT", Body: "SET NEW.a = 1"}
	b := &Trigger{Name: "x", Table: "t", Timing: "BEFORE", Event: "INSERT", Body: "  SET NEW.a = 1  "}
	if triggersChanged(a, b, n) {
		t.Fatalf("surrounding whitespace should not count as a change")
	}
	c := &Trigger{Name: "x", Table: "t", Timing: "BEFORE", Event: "INSERT", Body: "SET NEW.a = 2"}
	if !triggersChanged(a, c, n) {
		t.Fatalf("a genuine body change must diff")
	}
}

// ---- generator unit tests ----

func TestTriggerGenerate_CreateDDL(t *testing.T) {
	n := NormalizerFor(MySQL80)
	from := loadTriggerCatalog(t, `CREATE TABLE t (id INT PRIMARY KEY, a INT);`)
	to := loadTriggerCatalog(t, `
CREATE TABLE t (id INT PRIMARY KEY, a INT);
CREATE TRIGGER bi BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1;`)
	plan := GenerateMigrationWithNormalizer(from, to, DiffWithNormalizer(from, to, n), n)
	op := singleTriggerOp(t, plan, OpCreateTrigger)
	want := "CREATE TRIGGER `d`.`bi` BEFORE INSERT ON `d`.`t` FOR EACH ROW SET NEW.a = 1"
	if op.SQL != want {
		t.Fatalf("CREATE DDL mismatch:\n  got:  %s\n  want: %s", op.SQL, want)
	}
	if op.Phase != PhaseMain || op.Priority != priorityTrigger {
		t.Fatalf("CREATE op phase/priority wrong: phase=%v pri=%d", op.Phase, op.Priority)
	}
}

func TestTriggerGenerate_DropDDL(t *testing.T) {
	n := NormalizerFor(MySQL80)
	from := loadTriggerCatalog(t, `
CREATE TABLE t (id INT PRIMARY KEY, a INT);
CREATE TRIGGER bi BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1;`)
	to := loadTriggerCatalog(t, `CREATE TABLE t (id INT PRIMARY KEY, a INT);`)
	plan := GenerateMigrationWithNormalizer(from, to, DiffWithNormalizer(from, to, n), n)
	op := singleTriggerOp(t, plan, OpDropTrigger)
	want := "DROP TRIGGER `d`.`bi`"
	if op.SQL != want {
		t.Fatalf("DROP DDL mismatch:\n  got:  %s\n  want: %s", op.SQL, want)
	}
	if op.Phase != PhasePre {
		t.Fatalf("DROP op must be PhasePre, got %v", op.Phase)
	}
}

// TestTriggerGenerate_ModifyIsDropThenCreate proves a changed trigger emits exactly DROP then
// CREATE (no ALTER TRIGGER), and the DROP (PhasePre) sorts before the CREATE (PhaseMain).
func TestTriggerGenerate_ModifyIsDropThenCreate(t *testing.T) {
	n := NormalizerFor(MySQL80)
	from := loadTriggerCatalog(t, `
CREATE TABLE t (id INT PRIMARY KEY, a INT);
CREATE TRIGGER bi BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1;`)
	to := loadTriggerCatalog(t, `
CREATE TABLE t (id INT PRIMARY KEY, a INT);
CREATE TRIGGER bi BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 2;`)
	plan := GenerateMigrationWithNormalizer(from, to, DiffWithNormalizer(from, to, n), n)
	var trigOps []MigrationOp
	for _, op := range plan.Ops {
		if op.Type == OpCreateTrigger || op.Type == OpDropTrigger {
			trigOps = append(trigOps, op)
		}
	}
	if len(trigOps) != 2 {
		t.Fatalf("want 2 trigger ops (DROP+CREATE), got %d: %s", len(trigOps), plan.SQL())
	}
	if trigOps[0].Type != OpDropTrigger || trigOps[1].Type != OpCreateTrigger {
		t.Fatalf("want DROP then CREATE, got %s then %s", trigOps[0].Type, trigOps[1].Type)
	}
}

// TestTriggerGenerate_DropTableSuppressesTriggerDrop proves that when a trigger's owning table is
// dropped, no explicit DROP TRIGGER is emitted (the DROP TABLE cascades).
func TestTriggerGenerate_DropTableSuppressesTriggerDrop(t *testing.T) {
	n := NormalizerFor(MySQL80)
	from := loadTriggerCatalog(t, `
CREATE TABLE t (id INT PRIMARY KEY, a INT);
CREATE TRIGGER bi BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1;`)
	to := loadTriggerCatalog(t, `CREATE TABLE keep (id INT PRIMARY KEY);`)
	plan := GenerateMigrationWithNormalizer(from, to, DiffWithNormalizer(from, to, n), n)
	for _, op := range plan.Ops {
		if op.Type == OpDropTrigger {
			t.Fatalf("DROP TRIGGER should be suppressed when its table is dropped:\n%s", plan.SQL())
		}
	}
	// And a DROP TABLE for t must be present (sanity that the table really is being dropped).
	sawDropTable := false
	for _, op := range plan.Ops {
		if op.Type == OpDropTable && strings.Contains(op.SQL, "`t`") {
			sawDropTable = true
		}
	}
	if !sawDropTable {
		t.Fatalf("expected DROP TABLE `t` in plan:\n%s", plan.SQL())
	}
}

// TestTriggerGenerate_ModifyTableMoveOldTableDropped is a regression guard: a trigger whose owning
// table is RELOCATED (identity is (database, name), so a table move is a MODIFY) while its OLD table
// is dropped in the same plan must NOT emit an explicit DROP TRIGGER on the old table — the DROP
// TABLE cascades the trigger away first, and a standalone DROP TRIGGER would fail (errno 1360). Only
// the CREATE half (on the surviving new table) is emitted. (Confirmed against live MySQL: without
// the suppression the apply fails with "Trigger does not exist".)
func TestTriggerGenerate_ModifyTableMoveOldTableDropped(t *testing.T) {
	n := NormalizerFor(MySQL80)
	from := loadTriggerCatalog(t, `
CREATE TABLE old_t (id INT PRIMARY KEY, a INT);
CREATE TABLE new_t (id INT PRIMARY KEY, a INT);
CREATE TRIGGER tr BEFORE INSERT ON old_t FOR EACH ROW SET NEW.a = 1;`)
	to := loadTriggerCatalog(t, `
CREATE TABLE new_t (id INT PRIMARY KEY, a INT);
CREATE TRIGGER tr BEFORE INSERT ON new_t FOR EACH ROW SET NEW.a = 1;`)
	plan := GenerateMigrationWithNormalizer(from, to, DiffWithNormalizer(from, to, n), n)
	for _, op := range plan.Ops {
		if op.Type == OpDropTrigger {
			t.Fatalf("DROP TRIGGER must be suppressed when the trigger's old table is dropped (cascade):\n%s", plan.SQL())
		}
	}
	// The CREATE on the surviving new table must still be emitted.
	if op := singleTriggerOp(t, plan, OpCreateTrigger); !strings.Contains(op.SQL, "`new_t`") {
		t.Fatalf("CREATE TRIGGER must target the new table:\n%s", op.SQL)
	}
}

// TestTriggerGenerate_ModifyInPlaceStillDrops is the counterpart: an in-place body change (same
// table, table NOT dropped) must still emit DROP+CREATE — the suppression must not over-fire.
func TestTriggerGenerate_ModifyInPlaceStillDrops(t *testing.T) {
	n := NormalizerFor(MySQL80)
	from := loadTriggerCatalog(t, `
CREATE TABLE t (id INT PRIMARY KEY, a INT);
CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1;`)
	to := loadTriggerCatalog(t, `
CREATE TABLE t (id INT PRIMARY KEY, a INT);
CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 2;`)
	plan := GenerateMigrationWithNormalizer(from, to, DiffWithNormalizer(from, to, n), n)
	sawDrop, sawCreate := false, false
	for _, op := range plan.Ops {
		switch op.Type {
		case OpDropTrigger:
			sawDrop = true
		case OpCreateTrigger:
			sawCreate = true
		}
	}
	if !sawDrop || !sawCreate {
		t.Fatalf("in-place modify must emit both DROP and CREATE (suppression must not over-fire):\n%s", plan.SQL())
	}
}

// TestTriggerGenerate_CreateOnNewTable proves a trigger on a freshly-created table is emitted (the
// CREATE TRIGGER sorts after the CREATE TABLE via priorityTrigger > priorityTable).
func TestTriggerGenerate_CreateOnNewTable(t *testing.T) {
	n := NormalizerFor(MySQL80)
	from := New()
	to := loadTriggerCatalog(t, `
CREATE TABLE t (id INT PRIMARY KEY, a INT);
CREATE TRIGGER bi BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1;`)
	plan := GenerateMigrationWithNormalizer(from, to, DiffWithNormalizer(from, to, n), n)
	var createTableIdx, createTrigIdx = -1, -1
	for i, op := range plan.Ops {
		switch op.Type {
		case OpCreateTable:
			createTableIdx = i
		case OpCreateTrigger:
			createTrigIdx = i
		}
	}
	if createTableIdx < 0 || createTrigIdx < 0 {
		t.Fatalf("expected both CREATE TABLE and CREATE TRIGGER:\n%s", plan.SQL())
	}
	if createTrigIdx < createTableIdx {
		t.Fatalf("CREATE TRIGGER must follow CREATE TABLE:\n%s", plan.SQL())
	}
}

func singleTriggerOp(t *testing.T, plan *MigrationPlan, typ MigrationOpType) MigrationOp {
	t.Helper()
	var found []MigrationOp
	for _, op := range plan.Ops {
		if op.Type == typ {
			found = append(found, op)
		}
	}
	if len(found) != 1 {
		t.Fatalf("want exactly one %s op, got %d:\n%s", typ, len(found), plan.SQL())
	}
	return found[0]
}

// describeTriggerDiff renders a compact description of a SchemaDiff's trigger entries for failure
// output (the shared describeDiff only renders tables).
func describeTriggerDiff(d *SchemaDiff) string {
	var b strings.Builder
	for _, te := range d.Triggers {
		b.WriteString("[trigger ")
		b.WriteString(te.Name)
		b.WriteString(" ")
		b.WriteString(te.Action.String())
		b.WriteString("]")
	}
	if len(d.Tables) > 0 {
		b.WriteString(describeDiff(d))
	}
	return b.String()
}
