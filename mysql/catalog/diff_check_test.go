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
		{"unnamed-single", "CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (a > 0))", "t", only(MySQL80)},
		// Named single CHECK — the user name is preserved verbatim on both sides.
		{"named-single", "CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT chk_a CHECK (a > 0))", "t", only(MySQL80)},
		// Multiple unnamed → t_chk_1, t_chk_2, t_chk_3 in source order; the differ keys on the
		// (auto-)name, so all three collapse onto their readback.
		{"multiple-unnamed", "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CHECK (a > 0), CHECK (b < 100), CHECK (a < b))", "t", only(MySQL80)},
		// Mixed named + unnamed in one table.
		{"mixed-named-unnamed", "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ck1 CHECK (a > 0), CHECK (b < 100))", "t", only(MySQL80)},
		// NOT ENFORCED — the /*!80016 NOT ENFORCED */ marker; the enforced-state boolean must
		// round-trip.
		{"not-enforced", "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CONSTRAINT ne CHECK (a < b) NOT ENFORCED)", "t", only(MySQL80)},
		// ENFORCED explicitly stated (default) — the readback omits ENFORCED; must still be equal.
		{"enforced-explicit", "CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT en CHECK (a > 0) ENFORCED)", "t", only(MySQL80)},
		// Compound logical AND/OR — MySQL wraps each operand in parens; the loader's nodeToSQL
		// reproduces them on the user side too, so the canonical keys match.
		{"compound-and", "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CHECK (a >= 0 AND b <= 100))", "t", only(MySQL80)},
		{"compound-or", "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CHECK (a = 1 OR b = 2))", "t", only(MySQL80)},
		// BETWEEN / IN — operator keywords lowercased in the stored form.
		{"between", "CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (a BETWEEN 1 AND 10))", "t", only(MySQL80)},
		{"in-list", "CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (a IN (1, 2, 3)))", "t", only(MySQL80)},
		// String literal → 8.0 injects the `_utf8mb4` charset introducer on storage; CanonicalCheckExpr
		// must strip it so the user form (no introducer) and the readback compare equal.
		{"string-literal", "CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(20), CHECK (s = 'hello')) DEFAULT CHARSET=utf8mb4", "t", only(MySQL80)},
		// Function call — function name lowercased in the stored form (UPPER -> upper).
		{"function-call", "CREATE TABLE t (id INT PRIMARY KEY, s VARCHAR(20), CHECK (CHAR_LENGTH(s) > 0))", "t", only(MySQL80)},
		// Backtick-quoted identifiers in the user form vs bare-then-rebackticked stored form.
		{"backtick-idents", "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CHECK (`a` < `b`))", "t", only(MySQL80)},
		// Extra whitespace / odd spacing in the user form collapses to the stored single-spacing.
		{"odd-spacing", "CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (   a    >0   ))", "t", only(MySQL80)},
		// Rich compound: function + string + IN + AND on one table (the probe from the live readback
		// investigation) plus a second NOT ENFORCED check — the realistic multi-check table.
		{"rich-compound", "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, s VARCHAR(20), CONSTRAINT c1 CHECK (UPPER(s) = 'HI' AND a IN (1,2,3)), CONSTRAINT c2 CHECK (a < b) NOT ENFORCED) DEFAULT CHARSET=utf8mb4", "t", only(MySQL80)},
		// Column-level CHECK syntax (attached to a column definition) — MySQL treats it identically
		// to a table-level CHECK and auto-names it; must round-trip.
		{"column-level", "CREATE TABLE t (id INT PRIMARY KEY, a INT CHECK (a > 0))", "t", only(MySQL80)},
	}
}

// TestOracle_CheckDiffIdempotence proves gates 1 & 2 for every 8.0 CHECK form: the user DDL vs its
// engine readback diffs EMPTY, and the stored form self-diffs empty.
func TestOracle_CheckDiffIdempotence(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
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

// check57Cases enumerate CHECK DDLs that 5.7 PARSES BUT DROPS. Each is applied to live 5.7, read
// back (the readback has NO check), and the user-form catalog (which the omni loader DID populate
// with ConCheck entries) is diffed against the 5.7 readback under a 5.7 Normalizer — the result
// must be EMPTY. This is the no-phantom-on-5.7 guarantee made concrete: diffChecks returns nil on
// 5.7, so the loader's ConCheck entries never surface as a diff against the check-less readback.
func check57Cases() []diffProbe {
	return []diffProbe{
		{"unnamed", "CREATE TABLE t (id INT PRIMARY KEY, a INT, CHECK (a > 0))", "t", only(MySQL57)},
		{"named", "CREATE TABLE t (id INT PRIMARY KEY, a INT, CONSTRAINT chk_a CHECK (a > 0))", "t", only(MySQL57)},
		{"multiple", "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, CHECK (a > 0), CONSTRAINT c2 CHECK (b < 100))", "t", only(MySQL57)},
		{"column-level", "CREATE TABLE t (id INT PRIMARY KEY, a INT CHECK (a > 0))", "t", only(MySQL57)},
	}
}

// TestOracle_CheckDiff57NoPhantom proves the version gate: on 5.7, a table with a CHECK clause
// yields NO check diff. It asserts the property two ways:
//
//  1. STRUCTURAL — the loaded user-form catalog DOES carry a ConCheck constraint (the loader stores
//     it unconditionally), yet diffChecks(from, to, NormalizerFor(MySQL57)) on that table returns
//     no entries (the gate fires), so a self-diff under the 5.7 normalizer is empty.
//  2. AGAINST THE LIVE 5.7 READBACK — apply the CHECK DDL to real 5.7, read it back (the engine
//     dropped the CHECK), and diff the user form vs the readback under the 5.7 normalizer; the
//     result must be EMPTY despite the user side carrying a ConCheck the readback lacks.
func TestOracle_CheckDiff57NoPhantom(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	o := connectOracle(t, MySQL57)
	n := NormalizerFor(MySQL57)
	sc := serverCharsetFor(MySQL57)

	// Sanity: the version gate must report CHECK unsupported on 5.7 (so diffChecks short-circuits).
	if n.CheckSupported() {
		t.Fatal("precondition failed: CheckSupported() must be false on MySQL57")
	}

	for _, c := range check57Cases() {
		t.Run(c.id, func(t *testing.T) {
			// (1) STRUCTURAL: the loaded user catalog carries a ConCheck, but the 5.7 differ ignores it.
			userCat, userTbl := loadOneTable(t, sc, c.create, c.table)
			if countChecks(userTbl) == 0 {
				t.Fatalf("precondition failed: loader stored no ConCheck for %q (expected ≥1)", c.create)
			}
			if entries := diffChecks(userTbl, userTbl, n); len(entries) != 0 {
				t.Errorf("diffChecks on 5.7 must return nil, got %d entries", len(entries))
			}
			// A full self-diff under the 5.7 normalizer must also be empty (no phantom from the check).
			if d := DiffWithNormalizer(userCat, userCat, n); !d.IsEmpty() {
				t.Errorf("5.7 self-diff of a CHECK table not empty: %s%s", describeDiff(d), describeChecks(d))
			}

			// (2) AGAINST LIVE 5.7 READBACK: the engine dropped the CHECK; the diff must still be empty.
			readback, ok := o.applyAndReadback(t, "chk57_"+strings.ReplaceAll(c.id, "-", "_"), c.create, c.table)
			if !ok {
				t.Skipf("[5.7] could not obtain readback for %s", c.id)
			}
			if strings.Contains(strings.ToUpper(readback), "CHECK") {
				t.Logf("[5.7] note: readback unexpectedly mentions CHECK (engine may have stored it): %s", strings.TrimSpace(readback))
			}
			storedCat, _ := loadOneTable(t, sc, readback, c.table)
			if d := DiffWithNormalizer(userCat, storedCat, n); !d.IsEmpty() {
				t.Errorf("[5.7] user CHECK form vs check-less readback not empty (phantom check diff):\n  user:   %s\n  stored: %s\n  diff: %s%s",
					strings.TrimSpace(c.create), strings.TrimSpace(readback), describeDiff(d), describeChecks(d))
			}
			if d := DiffWithNormalizer(storedCat, userCat, n); !d.IsEmpty() {
				t.Errorf("[5.7] reverse (readback vs user CHECK form) not empty: %s%s", describeDiff(d), describeChecks(d))
			}
		})
	}
}

// TestOracle_CheckDiffNonEmptyCorrect proves the differ reports the RIGHT action for representative
// (from, to) check changes on 8.0 (loaded from real engine readbacks): an added, dropped, modified
// (expression change), and enforced-state-flipped check each produce the expected CheckDiffEntry.
func TestOracle_CheckDiffNonEmptyCorrect(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	o := connectOracle(t, MySQL80)
	n := NormalizerFor(MySQL80)
	sc := serverCharsetFor(MySQL80)

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
