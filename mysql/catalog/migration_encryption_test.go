package catalog

import (
	"strings"
	"testing"
)

// TestMigrationCarriesTableEncryption generates a migration to encrypted
// tables (one altered, one created file-per-table, one created in an encrypted
// general tablespace) and requires ENCRYPTION='Y' in each statement, and no
// difference left once the plan is applied. The pinned test server has no
// keyring, so ENCRYPTION='Y' itself cannot run there (3185); an explicit
// ENCRYPTION='N' reading back without the option is pinned in
// TestSimulationCausesOracle.
func TestMigrationCarriesTableEncryption(t *testing.T) {
	base := "CREATE DATABASE d; USE d; CREATE TABLESPACE ets ADD DATAFILE 'ets.ibd' ENCRYPTION='Y';"
	from, to := New(), New()
	execAlterClauses(t, from, base+" CREATE TABLE e (id INT);")
	execAlterClauses(t, to, base+" CREATE TABLE e (id INT) ENCRYPTION='Y'; CREATE TABLE f (id INT) ENCRYPTION='Y'; CREATE TABLE g (id INT) TABLESPACE ets ENCRYPTION='Y';")
	plan := GenerateMigration(from, to, Diff(from, to))
	for _, table := range []string{"`e`", "`f`", "`g`"} {
		found := false
		for _, op := range plan.Ops {
			if strings.Contains(op.SQL, table) && strings.Contains(op.SQL, "ENCRYPTION='Y'") {
				found = true
			}
		}
		if !found {
			t.Fatalf("no ENCRYPTION='Y' for %s in:\n%s", table, plan.SQL())
		}
	}
	for _, op := range plan.Ops {
		execAlterClauses(t, from, op.SQL+";")
	}
	if rest := GenerateMigration(from, to, Diff(from, to)); len(rest.Ops) != 0 {
		t.Fatalf("migration left a difference:\n%s", rest.SQL())
	}

	// An explicit ENCRYPTION='N' is no difference from none.
	nFrom, nTo := New(), New()
	execAlterClauses(t, nFrom, "CREATE DATABASE d; USE d; CREATE TABLE n (id INT) ENCRYPTION='N';")
	execAlterClauses(t, nTo, "CREATE DATABASE d; USE d; CREATE TABLE n (id INT);")
	if p := GenerateMigration(nFrom, nTo, Diff(nFrom, nTo)); len(p.Ops) != 0 {
		t.Fatalf("ENCRYPTION='N' against none produced:\n%s", p.SQL())
	}
}

// TestMigrationUppercaseFilePerTableTablespace moves a table into a general
// tablespace named INNODB_FILE_PER_TABLE in capitals, an ordinary name on
// mysql:8.0.46: only the lowercase innodb_file_per_table is the built-in.
func TestMigrationUppercaseFilePerTableTablespace(t *testing.T) {
	base := "CREATE DATABASE d; USE d; CREATE TABLESPACE INNODB_FILE_PER_TABLE ADD DATAFILE 'u.ibd';"
	from, to := New(), New()
	execAlterClauses(t, from, base+" CREATE TABLE m (id INT);")
	execAlterClauses(t, to, base+" CREATE TABLE m (id INT) TABLESPACE INNODB_FILE_PER_TABLE; CREATE TABLE n (id INT) TABLESPACE INNODB_FILE_PER_TABLE;")
	plan := GenerateMigration(from, to, Diff(from, to))
	for _, table := range []string{"`m`", "`n`"} {
		found := false
		for _, op := range plan.Ops {
			if strings.Contains(op.SQL, table) && strings.Contains(op.SQL, "`INNODB_FILE_PER_TABLE`") {
				found = true
			}
		}
		if !found {
			t.Fatalf("no TABLESPACE INNODB_FILE_PER_TABLE for %s in:\n%s", table, plan.SQL())
		}
	}
}
