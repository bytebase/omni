package catalog

import (
	"context"
	"fmt"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// Apply probes for a new table whose AUTO_INCREMENT column is keyed by an index that is not
// the primary key, on OceanBase. MySQL refuses to create such a table without the key (1075),
// which is what the cases were written for; OceanBase lets the column stand without one
// (TestOracle_OceanBaseAlterTableClauses), so the plan is the CREATE TABLE and then the key.
//
// Two properties per case:
//   - APPLY-CORRECTNESS: the full generated plan (CREATE TABLE + the follow-up index ADDs)
//     applies to a real database and the result reads back equal to `to`.
//   - IDEMPOTENCE: the stored form self-diffs empty and the user form vs the stored form diffs
//     empty.
//
// The apply is SELF-CONTAINED (its own per-probe database via loadOnePartTable), like the partition
// oracle test, so it does not contend on the shared `diffdb`. It reuses oceanBaseOracle / NormalizerFor
// / serverCharsetFor from the shared harness.

// autoIncProbe is a single CREATE-TABLE form whose AUTO_INCREMENT column is keyed only by a non-PK
// key; the generated plan must apply cleanly and round-trip empty.
type autoIncProbe struct {
	id        string
	table     string
	createSQL string
}

func autoIncSupportProbes() []autoIncProbe {
	return []autoIncProbe{
		// The headline case: AUTO_INCREMENT on a single-column NON-PK UNIQUE key.
		{"autoinc-nonpk-unique", "t",
			"CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, name VARCHAR(20), UNIQUE KEY uk (id)) ENGINE=InnoDB"},
		// AUTO_INCREMENT first in a COMPOSITE unique key (valid: it is the leading column).
		{"autoinc-nonpk-unique-composite", "t",
			"CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, a INT NOT NULL, UNIQUE KEY uk (id, a)) ENGINE=InnoDB"},
		// AUTO_INCREMENT backed only by a PLAIN (non-unique) secondary key.
		{"autoinc-plain-key", "t",
			"CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, a INT, KEY k (id)) ENGINE=InnoDB"},
		// AUTO_INCREMENT keyed by a non-PK unique while the table ALSO has a separate PRIMARY KEY on
		// another column.
		{"autoinc-unique-with-other-pk", "t",
			"CREATE TABLE t (pk INT NOT NULL, id INT NOT NULL AUTO_INCREMENT, PRIMARY KEY (pk), UNIQUE KEY uk (id)) ENGINE=InnoDB"},
		// REGRESSION GUARD: AUTO_INCREMENT AS the PRIMARY KEY — the common case must keep rendering
		// the PK inline.
		{"autoinc-primary-key-guard", "t",
			"CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY, name VARCHAR(20)) ENGINE=InnoDB"},
	}
}

// autoIncFKProbes covers the AUTO_INCREMENT column whose only key, on MySQL, is the index MySQL
// adds for a FOREIGN KEY on it. On OceanBase the column has no key and the foreign key no
// index: the plan is the table and then the foreign key. The referenced parent table is created
// in the same plan, so these are full multi-table apply cases.
func autoIncFKProbes() []struct {
	id, parentDDL, childDDL, table string
} {
	return []struct {
		id, parentDDL, childDDL, table string
	}{
		{"autoinc-fk-implicit-only",
			"CREATE TABLE p (id INT NOT NULL PRIMARY KEY)",
			"CREATE TABLE c (id INT NOT NULL AUTO_INCREMENT, x INT, FOREIGN KEY (id) REFERENCES p(id)) ENGINE=InnoDB",
			"c"},
	}
}

// TestOracle_AutoIncFKImplicitApply proves the FK-implicit-backed AUTO_INCREMENT case: the
// generated multi-table plan (parent + child + deferred FK) applies cleanly and the child reads
// back equal to its stored form.
func TestOracle_AutoIncFKImplicitApply(t *testing.T) {
	for _, version := range oracleVersions() {
		o := oceanBaseOracle(t)
		n := NormalizerFor(version)
		for _, p := range autoIncFKProbes() {
			t.Run(fmt.Sprintf("%s/%s", o.name, p.id), gapped(func(t testing.TB) {
				assertAutoIncFKApply(t, o, n, p.id, p.parentDDL, p.childDDL, p.table)
			}))
		}
	}
}

func assertAutoIncFKApply(t testing.TB, o *oracleConn, n *Normalizer, id, parentDDL, childDDL, childTable string) {
	t.Helper()
	slug := strings.ReplaceAll(id, "-", "_")
	sc := serverCharsetFor(o.version)
	ctx := context.Background()

	conn, err := o.session(ctx)
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
			o.refused(t, err, "%s `to` setup failed: %q", id, s)
			return
		}
	}
	readback := func(tbl string) string {
		var name, ddl string
		if err := conn.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE TABLE %s.%s", applyDB, tbl)).Scan(&name, &ddl); err != nil {
			o.refused(t, err, "%s SHOW CREATE %s failed", id, tbl)
			return ""
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
			o.refused(t, err, "%s apply db reset failed", id)
			return
		}
	}
	if !o.applyPlan(t, conn, plan) {
		return
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
	for _, version := range oracleVersions() {
		o := oceanBaseOracle(t)
		n := NormalizerFor(version)
		for _, p := range autoIncSupportProbes() {
			t.Run(fmt.Sprintf("%s/%s", o.name, p.id), gapped(func(t testing.TB) {
				assertAutoIncApply(t, o, n, p)
			}))
		}
	}
}

// assertAutoIncApply builds `to` from the engine's own readback, generates the plan from empty→to,
// applies it to a fresh per-probe database, and asserts the apply succeeds and the result diffs
// empty against `to` (both directions).
func assertAutoIncApply(t testing.TB, o *oracleConn, n *Normalizer, p autoIncProbe) {
	t.Helper()
	slug := strings.ReplaceAll(p.id, "-", "_")
	sc := serverCharsetFor(o.version)
	ctx := context.Background()

	conn, err := o.session(ctx)
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
			o.refused(t, err, "%s `to` readback setup failed: %q", p.id, s)
			return
		}
	}
	var name, rb string
	if err := conn.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE TABLE rto_%s.%s", slug, p.table)).Scan(&name, &rb); err != nil {
		o.refused(t, err, "%s SHOW CREATE failed", p.id)
		return
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
			o.refused(t, err, "%s apply db setup failed", p.id)
			return
		}
	}
	if !o.applyPlan(t, conn, plan) {
		return
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
	for range oracleVersions() {
		o := oceanBaseOracle(t)
		for _, p := range autoIncSupportProbes() {
			t.Run(fmt.Sprintf("%s/%s", o.name, p.id), gapped(func(t testing.TB) {
				db := "aidem_" + strings.ReplaceAll(p.id, "-", "_")
				assertDiffEmptyAgainstReadback(t, o, db, p.createSQL, p.table)
			}))
		}
	}
}
