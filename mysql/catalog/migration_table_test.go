package catalog

import (
	"context"
	"fmt"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// Oracle proof for the AUTO_INCREMENT-supporting-key inlining in formatCreateTable (the errno-1075
// fix), against the LIVE MySQL engines (5.7 :13307, 8.0 :13306). MySQL rejects CREATE TABLE for an
// AUTO_INCREMENT column that is not the first column of a key ("there can be only one auto column
// and it must be defined as a key", errno 1075). When the only backing key is a non-PK index,
// generate-core previously rendered the CREATE without that key, so the statement failed on apply.
//
// Two properties per case, on every supported version:
//   - APPLY-CORRECTNESS: the full generated plan (CREATE TABLE + any follow-up index ADDs) applies
//     to a real database and the result reads back equal to `to` — proving the CREATE is accepted
//     (the supporting key is inline) AND that the inlined key is NOT also re-added by the index
//     node (which would fail errno 1061, duplicate key name).
//   - IDEMPOTENCE: the stored form self-diffs empty and the user form vs the stored form diffs
//     empty (no phantom diff from the inlined key).
//
// The apply is SELF-CONTAINED (its own per-probe database via loadOnePartTable), like the partition
// oracle test, so it does not contend on the shared `diffdb`. It reuses connectOracle / NormalizerFor
// / serverCharsetFor / both from the shared harness and skips cleanly when the engines are down.

// autoIncProbe is a single CREATE-TABLE form whose AUTO_INCREMENT column is backed only by a non-PK
// key; the generated plan must apply cleanly and round-trip empty.
type autoIncProbe struct {
	id        string
	table     string
	createSQL string
	versions  []Version
}

func autoIncSupportProbes() []autoIncProbe {
	return []autoIncProbe{
		// The headline case: AUTO_INCREMENT on a single-column NON-PK UNIQUE key (errno 1075 without
		// the inline key).
		{"autoinc-nonpk-unique", "t",
			"CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, name VARCHAR(20), UNIQUE KEY uk (id)) ENGINE=InnoDB", both()},
		// AUTO_INCREMENT first in a COMPOSITE unique key (valid: it is the leading column).
		{"autoinc-nonpk-unique-composite", "t",
			"CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, UNIQUE KEY uk (id, a)) ENGINE=InnoDB", both()},
		// AUTO_INCREMENT backed only by a PLAIN (non-unique) secondary key.
		{"autoinc-plain-key", "t",
			"CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, a INT, KEY k (id)) ENGINE=InnoDB", both()},
		// AUTO_INCREMENT backed by a non-PK unique while the table ALSO has a separate PRIMARY KEY on
		// another column (the PK does not cover the auto column, so the unique key must be inline).
		{"autoinc-unique-with-other-pk", "t",
			"CREATE TABLE t (pk INT NOT NULL, id INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (pk), UNIQUE KEY uk (id)) ENGINE=InnoDB", both()},
		// REGRESSION GUARD: AUTO_INCREMENT AS the PRIMARY KEY — the common case must keep rendering
		// the PK inline (and NOT additionally inline a supporting key).
		{"autoinc-primary-key-guard", "t",
			"CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY, name VARCHAR(20)) ENGINE=InnoDB", both()},
	}
}

// autoIncFKProbes covers the AUTO_INCREMENT column whose ONLY first-column key is the index MySQL
// auto-creates for a FOREIGN KEY on it. The CREATE TABLE must inline that backing key (else errno
// 1075), and the deferred FK must reuse it (no duplicate-index / errno-1061). The referenced parent
// table is created in the same plan, so these are full multi-table apply cases.
func autoIncFKProbes() []struct {
	id, parentDDL, childDDL, table string
	versions                       []Version
} {
	return []struct {
		id, parentDDL, childDDL, table string
		versions                       []Version
	}{
		{"autoinc-fk-implicit-only",
			"CREATE TABLE p (id INT NOT NULL PRIMARY KEY)",
			"CREATE TABLE c (id INT NOT NULL AUTO_INCREMENT, x INT, FOREIGN KEY (id) REFERENCES p(id)) ENGINE=InnoDB",
			"c", both()},
	}
}

// TestOracle_AutoIncFKImplicitApply proves the FK-implicit-backed AUTO_INCREMENT case: the
// generated multi-table plan (parent + child + deferred FK) applies cleanly and the child reads
// back equal to its stored form.
func TestOracle_AutoIncFKImplicitApply(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		n := NormalizerFor(version)
		for _, p := range autoIncFKProbes() {
			if !containsVersion(p.versions, version) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, p.id), func(t *testing.T) {
				assertAutoIncFKApply(t, o, n, p.id, p.parentDDL, p.childDDL, p.table)
			})
		}
	}
}

func assertAutoIncFKApply(t *testing.T, o *oracleConn, n *Normalizer, id, parentDDL, childDDL, childTable string) {
	t.Helper()
	slug := strings.ReplaceAll(id, "-", "_")
	sc := serverCharsetFor(o.version)
	ctx := context.Background()

	conn, err := o.db.Conn(ctx)
	if err != nil {
		t.Fatalf("[%s] grab conn: %v", o.name, err)
	}
	defer func() { _ = conn.Close() }()

	applyDB := "aifk_" + slug
	// Build the `to` catalog by applying both tables in the apply db, then reading them back.
	for _, s := range []string{
		"DROP DATABASE IF EXISTS " + applyDB,
		fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s", applyDB, sc),
		"USE " + applyDB,
		parentDDL,
		childDDL,
	} {
		if _, err := conn.ExecContext(ctx, s); err != nil {
			t.Skipf("[%s] %s `to` setup failed: %q: %v", o.name, id, s, err)
		}
	}
	readback := func(tbl string) string {
		var name, ddl string
		if err := conn.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE TABLE %s.%s", applyDB, tbl)).Scan(&name, &ddl); err != nil {
			t.Skipf("[%s] %s SHOW CREATE %s failed: %v", o.name, id, tbl, err)
		}
		return ddl
	}
	parentRB, childRB := readback("p"), readback(childTable)
	// Compose both readbacks into one apply-db catalog so FK references resolve.
	toSQL := fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s;\nUSE %s;\n%s;\n%s;", applyDB, sc, applyDB, parentRB, childRB)
	toCat, err := LoadSQL(toSQL)
	if err != nil {
		t.Fatalf("[%s] %s load `to`: %v", o.name, id, err)
	}

	from := New()
	diff := DiffWithNormalizer(from, toCat, n)
	plan := GenerateMigrationWithNormalizer(from, toCat, diff, n)

	// Apply to a fresh database.
	for _, s := range []string{
		"DROP DATABASE IF EXISTS " + applyDB,
		fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s", applyDB, sc),
		"USE " + applyDB,
	} {
		if _, err := conn.ExecContext(ctx, s); err != nil {
			t.Skipf("[%s] %s apply db reset failed: %v", o.name, id, err)
		}
	}
	for _, op := range plan.Ops {
		if _, err := conn.ExecContext(ctx, op.SQL); err != nil {
			t.Fatalf("[%s] APPLY FAILED for %s:\n  stmt: %s\n  err: %v\n  full plan:\n%s",
				o.name, id, op.SQL, err, plan.SQL())
		}
	}
	// Reload the applied child and diff against `to`.
	resultSQL := fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s;\nUSE %s;\n%s;\n%s;", applyDB, sc, applyDB, readback("p"), readback(childTable))
	resultCat, err := LoadSQL(resultSQL)
	if err != nil {
		t.Fatalf("[%s] %s load result: %v", o.name, id, err)
	}
	if d := DiffWithNormalizer(resultCat, toCat, n); !d.IsEmpty() {
		t.Errorf("[%s] APPLY-CORRECTNESS FAILED for %s: result != to\n  plan:\n%s\n  diff: %s",
			o.name, id, plan.SQL(), describeDiff(d))
	}
}

// TestOracle_AutoIncSupportingKeyApply proves the generated plan for an AUTO_INCREMENT column backed
// only by a non-PK key applies cleanly on a real engine and reads back equal to `to`.
func TestOracle_AutoIncSupportingKeyApply(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		n := NormalizerFor(version)
		for _, p := range autoIncSupportProbes() {
			if !containsVersion(p.versions, version) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, p.id), func(t *testing.T) {
				assertAutoIncApply(t, o, n, p)
			})
		}
	}
}

// assertAutoIncApply builds `to` from the engine's own readback, generates the plan from empty→to,
// applies it to a fresh per-probe database, and asserts the apply succeeds and the result diffs
// empty against `to` (both directions). A failure here is either errno 1075 (supporting key not
// inline) or errno 1061 (inlined key also re-added by the index node).
func assertAutoIncApply(t *testing.T, o *oracleConn, n *Normalizer, p autoIncProbe) {
	t.Helper()
	slug := strings.ReplaceAll(p.id, "-", "_")
	sc := serverCharsetFor(o.version)
	ctx := context.Background()

	conn, err := o.db.Conn(ctx)
	if err != nil {
		t.Fatalf("[%s] grab conn: %v", o.name, err)
	}
	defer func() { _ = conn.Close() }()

	// Readback `to` in a throwaway db, then load it into the apply db so identities match.
	for _, s := range []string{
		"DROP DATABASE IF EXISTS rto_" + slug,
		fmt.Sprintf("CREATE DATABASE rto_%s DEFAULT CHARSET=%s", slug, sc),
		"USE rto_" + slug,
		p.createSQL,
	} {
		if _, err := conn.ExecContext(ctx, s); err != nil {
			t.Skipf("[%s] %s `to` readback setup failed: %q: %v", o.name, p.id, s, err)
		}
	}
	var name, rb string
	if err := conn.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE TABLE rto_%s.%s", slug, p.table)).Scan(&name, &rb); err != nil {
		t.Skipf("[%s] %s SHOW CREATE failed: %v", o.name, p.id, err)
	}

	applyDB := "ai_" + slug
	toCat, _ := loadOnePartTable(t, applyDB, sc, rb, p.table)

	from := New()
	diff := DiffWithNormalizer(from, toCat, n)
	plan := GenerateMigrationWithNormalizer(from, toCat, diff, n)

	for _, s := range []string{
		"DROP DATABASE IF EXISTS " + applyDB,
		fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s", applyDB, sc),
		"USE " + applyDB,
	} {
		if _, err := conn.ExecContext(ctx, s); err != nil {
			t.Skipf("[%s] %s apply db setup failed: %v", o.name, p.id, err)
		}
	}
	for _, op := range plan.Ops {
		if _, err := conn.ExecContext(ctx, op.SQL); err != nil {
			t.Fatalf("[%s] APPLY FAILED for %s:\n  stmt: %s\n  err: %v\n  full plan:\n%s",
				o.name, p.id, op.SQL, err, plan.SQL())
		}
	}

	var rname, resultRB string
	if err := conn.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE TABLE %s.%s", applyDB, p.table)).Scan(&rname, &resultRB); err != nil {
		t.Fatalf("[%s] %s: result table missing after apply:\n%s", o.name, p.id, plan.SQL())
	}
	resultCat, _ := loadOnePartTable(t, applyDB, sc, resultRB, p.table)

	if d := DiffWithNormalizer(resultCat, toCat, n); !d.IsEmpty() {
		t.Errorf("[%s] APPLY-CORRECTNESS FAILED for %s: result != to\n  plan:\n%s\n  result: %s\n  diff: %s",
			o.name, p.id, plan.SQL(), strings.TrimSpace(resultRB), describeDiff(d))
	}
	if d := DiffWithNormalizer(toCat, resultCat, n); !d.IsEmpty() {
		t.Errorf("[%s] APPLY-CORRECTNESS (reverse) FAILED for %s: %s", o.name, p.id, describeDiff(d))
	}
}

// TestOracle_AutoIncSupportingKeyIdempotence proves the AUTO_INCREMENT-supporting-key forms
// round-trip empty: the stored form self-diffs empty and the user form vs the stored form diffs
// empty (no phantom diff introduced by the inlined key).
func TestOracle_AutoIncSupportingKeyIdempotence(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		for _, p := range autoIncSupportProbes() {
			if !containsVersion(p.versions, version) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, p.id), func(t *testing.T) {
				db := "aidem_" + strings.ReplaceAll(p.id, "-", "_")
				assertDiffEmptyAgainstReadback(t, o, db, p.createSQL, p.table)
			})
		}
	}
}
