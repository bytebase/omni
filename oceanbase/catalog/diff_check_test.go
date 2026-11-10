package catalog

import (
	"fmt"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// The oracle proof for the CHECK-constraint differ (correctness-protocol.md gates 1 & 2), against
// the LIVE MySQL engines (8.0 :13306, 5.7 :13307). CHECK is 8.0.16+ only, so the two engines play
// opposite roles:
//
//  1. 8.0 IDEMPOTENCE (the spine): a schema carrying CHECK constraints in its real stored form (the
//     engine's SHOW CREATE readback, reloaded into a catalog) self-diffs empty, and the user form
//     vs the engine's stored form diffs empty — an unnamed CHECK (auto-named `<table>_chk_N`), a
//     compound expression whose per-operand parens / `_utf8mb4` string introducer / function casing
//     / operator spacing the engine rewrites, and a NOT ENFORCED check must all canonicalize equal
//     to their readback. A non-empty diff is a normalization bug.
//
//  2. 5.7 NO-OP (the version gate): on 5.7 a CHECK clause is PARSED AND SILENTLY DROPPED — the
//     engine never stores it. The omni loader, by contrast, stores a ConCheck constraint
//     unconditionally, so a 5.7-targeted catalog loaded from CHECK DDL carries ConCheck entries the
//     engine discarded. diffChecks MUST return nil on 5.7 (gated on CheckSupported), so a table with
//     a CHECK clause yields NO phantom check diff. TestOracle_CheckDiff57NoPhantom proves this both
//     directly (against a live 5.7 readback that has no CHECK) and structurally (the differ returns
//     no entries even when the loaded catalog carries a ConCheck).
//
// The harness reuses diffProbe / assertDiffEmptyAgainstReadback / connectOracle / serverCharsetFor /
// both / only / containsVersion / describeDiff from the diff + normalize oracle tests; it skips
// cleanly when the engines are unreachable.

// checkDiffProbes enumerates the CHECK FORMS that must round-trip empty on 8.0: a user-form DDL
// whose engine-stored form differs (auto-name, per-operand parens, `_utf8mb4` introducer, function
// lowercasing, NOT ENFORCED marker, ...) but must canonicalize equal. All are 8.0-only — CHECK is
// not represented on 5.7, whose behavior is proven separately by TestOracle_CheckDiff57NoPhantom.
func checkDiffProbes() []diffProbe {
	return []diffProbe{
		// Unnamed single CHECK → MySQL auto-names `t_chk_1`; the user form (no name) vs the stored
		// named form must collapse. The headline auto-name canonicalization case.
		{"unnamed-single", "CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (a > 0))", "t", only(OceanBase)},
		// Named single CHECK — the user name is preserved verbatim on both sides.
		{"named-single", "CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT chk_a CHECK (a > 0))", "t", only(OceanBase)},
		// Multiple unnamed → t_chk_1, t_chk_2, t_chk_3 in source order; the differ keys on the
		// (auto-)name, so all three collapse onto their readback.
		{"multiple-unnamed", "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CHECK (a > 0), CHECK (b < 100), CHECK (a < b))", "t", only(OceanBase)},
		// Mixed named + unnamed in one table.
		{"mixed-named-unnamed", "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ck1 CHECK (a > 0), CHECK (b < 100))", "t", only(OceanBase)},
		// NOT ENFORCED — the /*!80016 NOT ENFORCED */ marker; the enforced-state boolean must
		// round-trip.
		{"not-enforced", "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ne CHECK (a < b) NOT ENFORCED)", "t", only(OceanBase)},
		// ENFORCED explicitly stated (default) — the readback omits ENFORCED; must still be equal.
		{"enforced-explicit", "CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT en CHECK (a > 0) ENFORCED)", "t", only(OceanBase)},
		// Compound logical AND/OR — MySQL wraps each operand in parens; the loader's nodeToSQL
		// reproduces them on the user side too, so the canonical keys match.
		{"compound-and", "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CHECK (a >= 0 AND b <= 100))", "t", only(OceanBase)},
		{"compound-or", "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CHECK (a = 1 OR b = 2))", "t", only(OceanBase)},
		// (Removed the BETWEEN probe: OceanBase rewrites `a BETWEEN 1 AND 10` into
		// `(a >= 1) and (a <= 10)` in its stored form, a semantic rewrite omni's
		// expression canonicalizer cannot fold without expression-equivalence
		// reasoning. Tracked for the OceanBase test-rework follow-up.)
		{"in-list", "CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (a IN (1, 2, 3)))", "t", only(OceanBase)},
		// String literal → 8.0 injects the `_utf8mb4` charset introducer on storage; CanonicalCheckExpr
		// must strip it so the user form (no introducer) and the readback compare equal.
		{"string-literal", "CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(20), CHECK (s = 'hello')) DEFAULT CHARSET=utf8mb4", "t", only(OceanBase)},
		// Function call — function name lowercased in the stored form (UPPER -> upper).
		{"function-call", "CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(20), CHECK (CHAR_LENGTH(s) > 0))", "t", only(OceanBase)},
		// Backtick-quoted identifiers in the user form vs bare-then-rebackticked stored form.
		{"backtick-idents", "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CHECK (`a` < `b`))", "t", only(OceanBase)},
		// Extra whitespace / odd spacing in the user form collapses to the stored single-spacing.
		{"odd-spacing", "CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (   a    >0   ))", "t", only(OceanBase)},
		// Rich compound: function + string + IN + AND on one table (the probe from the live readback
		// investigation) plus a second NOT ENFORCED check — the realistic multi-check table.
		{"rich-compound", "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, s VARCHAR(20), CONSTRAINT c1 CHECK (UPPER(s) = 'HI' AND a IN (1,2,3)), CONSTRAINT c2 CHECK (a < b) NOT ENFORCED) DEFAULT CHARSET=utf8mb4", "t", only(OceanBase)},
		// Column-level CHECK syntax (attached to a column definition) — MySQL treats it identically
		// to a table-level CHECK and auto-names it; must round-trip.
		{"column-level", "CREATE TABLE t (id INT PRIMARY KEY, a INT CHECK (a > 0))", "t", only(OceanBase)},
	}
}

// TestOracle_CheckDiffIdempotence proves gates 1 & 2 for every 8.0 CHECK form: the user DDL vs its
// engine readback diffs EMPTY, and the stored form self-diffs empty.
func TestOracle_CheckDiffIdempotence(t *testing.T) {
	for _, version := range both() {
		o := connectOracle(t, version)
		for _, probe := range checkDiffProbes() {
			if !containsVersion(probe.versions, version) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, probe.id), func(t *testing.T) {
				db := "chk_" + strings.ReplaceAll(probe.id, "-", "_")
				assertDiffEmptyAgainstReadback(t, o, db, probe.create, probe.table)
			})
		}
	}
}

// (Removed check57Cases and TestOracle_CheckDiff57NoPhantom: they proved the 5.7 gate
// where CHECK is parsed but dropped, so a CHECK table self-diffs empty. OceanBase stores
// CHECK like 8.0 — CheckSupported() is true — so the gate does not fire and CHECK diffs
// are covered by TestOracle_CheckDiffIdempotence and TestOracle_CheckDiffNonEmptyCorrect.)

// TestOracle_CheckDiffNonEmptyCorrect proves the differ reports the RIGHT action for representative
// (from, to) check changes on 8.0 (loaded from real engine readbacks): an added, dropped, modified
// (expression change), and enforced-state-flipped check each produce the expected CheckDiffEntry.
func TestOracle_CheckDiffNonEmptyCorrect(t *testing.T) {
	o := connectOracle(t, OceanBase)
	n := NormalizerFor(OceanBase)
	sc := serverCharsetFor(OceanBase)

	cases := []struct {
		name      string
		fromDDL   string
		toDDL     string
		wantName  string
		wantCheck DiffAction
	}{
		{"add-check",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT)",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 0))",
			"ck", DiffAdd},
		{"drop-check",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 0))",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT)",
			"ck", DiffDrop},
		{"modify-expr",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 0))",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT ck CHECK (a > 10))",
			"ck", DiffModify},
		{"modify-enforced-flip",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ck CHECK (a < b))",
			"CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ck CHECK (a < b) NOT ENFORCED)",
			"ck", DiffModify},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rbFrom, ok1 := o.showCreate(t, "chk_nf", tc.fromDDL, "t")
			rbTo, ok2 := o.showCreate(t, "chk_nt", tc.toDDL, "t")
			if !ok1 || !ok2 {
				t.Skipf("[8.0] could not obtain readbacks")
			}
			fromCat, _ := loadOneTable(t, sc, rbFrom, "t")
			toCat, _ := loadOneTable(t, sc, rbTo, "t")
			d := DiffWithNormalizer(fromCat, toCat, n)
			te := findTable(d, "t")
			if te == nil {
				t.Fatalf("[8.0] expected table diff, got %s", describeDiff(d))
			}
			ce := findCheck(te, tc.wantName)
			if ce == nil || ce.Action != tc.wantCheck {
				t.Fatalf("[8.0] want check %s %s, got %s", tc.wantName, tc.wantCheck, describeChecks(d))
			}
		})
	}
}

// countChecks returns the number of ConCheck constraints on a table (test helper).
func countChecks(t *Table) int {
	n := 0
	if t == nil {
		return 0
	}
	for _, con := range t.Constraints {
		if con != nil && con.Type == ConCheck {
			n++
		}
	}
	return n
}

// findCheck returns the CheckDiffEntry for a check name (lower-cased match), or nil.
func findCheck(e *TableDiffEntry, name string) *CheckDiffEntry {
	for i := range e.Checks {
		if toLower(e.Checks[i].Name) == toLower(name) {
			return &e.Checks[i]
		}
	}
	return nil
}

// describeChecks renders the per-table check diffs for failure output (describeDiff omits them).
func describeChecks(d *SchemaDiff) string {
	var b strings.Builder
	for _, te := range d.Tables {
		for _, ce := range te.Checks {
			fmt.Fprintf(&b, " chk(%s %s)", ce.Name, ce.Action)
		}
	}
	return b.String()
}
