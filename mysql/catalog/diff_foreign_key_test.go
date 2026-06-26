package catalog

import (
	"context"
	"fmt"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// The oracle proof for the foreign-key differ (correctness-protocol.md gates 1 & 2), against the
// LIVE MySQL engines (5.7 :13307, 8.0 :13306). A foreign key always lives on a CHILD table that
// references a PARENT, so every FK form is proven through a multi-table round-trip: parent `p` +
// child `c` are applied to the real engine, BOTH are read back via SHOW CREATE, reloaded as one
// schema, and the self-diff must be empty. The headline properties:
//
//  1. Idempotence (the spine): the engine's stored form of an FK schema self-diffs empty.
//  2. Action-default normalization: the user form of an FK (with no explicit ON DELETE, or an
//     explicit RESTRICT / NO ACTION) and the engine's stored form must collapse to empty — even
//     though SHOW CREATE echoes the default action OPPOSITELY per version (8.0 omits NO ACTION,
//     echoes RESTRICT; 5.7 omits RESTRICT, echoes NO ACTION). canonicalFKAction folds all three
//     onto one key, so the diff is empty on both versions regardless of which form the readback
//     took.
//  3. Auto-name normalization: an unnamed FK (auto-named `<table>_ibfk_N`) and an unnamed-backing
//     index round-trip empty.
//
// The harness reuses connectOracle / serverCharsetFor / both / only / containsVersion /
// describeDiff from the diff + normalize oracle tests; it skips cleanly when the engines are
// unreachable.

// fkDiffCase is a parent+child schema with one or more foreign keys, used to prove the FK diff
// round-trips empty. childDDL defines table `c` referencing table `p`; parentDDL defines `p`.
type fkDiffCase struct {
	id        string
	parentDDL string
	childDDL  string
	versions  []Version
}

// fkParentSingle is the common single-column-PK parent (most FK cases reference p(id)).
const fkParentSingle = "CREATE TABLE p (id INT NOT NULL PRIMARY KEY, a INT NOT NULL, b INT NOT NULL, UNIQUE KEY ua (a), UNIQUE KEY ub (b))"

// fkParentComposite is the composite-PK parent for composite-FK cases.
const fkParentComposite = "CREATE TABLE p (a INT NOT NULL, b INT NOT NULL, c INT, PRIMARY KEY (a,b))"

// fkDiffCases enumerate the FK FORMS that must round-trip empty. Each builds parent + child,
// reads BOTH back, reloads as one schema, and asserts the self-diff is empty (and the user form
// vs the stored form is empty).
func fkDiffCases() []fkDiffCase {
	return []fkDiffCase{
		// ---- naming ----
		// Named FK, no backing index → MySQL auto-creates KEY `fk`, constraint kept verbatim.
		{"named", fkParentSingle,
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))", both()},
		// Unnamed FK → constraint auto-named `c_ibfk_1`; backing index named after the column.
		{"unnamed-autoname", fkParentSingle,
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT, FOREIGN KEY (pid) REFERENCES p(id))", both()},
		// User-chosen name that matches the `<table>_ibfk_N` auto shape.
		{"user-named-ibfk", fkParentSingle,
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT c_ibfk_1 FOREIGN KEY (pid) REFERENCES p(id))", both()},

		// ---- ON DELETE / ON UPDATE actions (the action-default normalization core) ----
		// Bare FK: no action clause. Stored form has no clause on both versions.
		{"action-bare", fkParentSingle,
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))", both()},
		// Explicit RESTRICT on both: 8.0 ECHOES it, 5.7 OMITS it — must round-trip empty on both.
		{"action-restrict", fkParentSingle,
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id) ON DELETE RESTRICT ON UPDATE RESTRICT)", both()},
		// Explicit NO ACTION on both: 8.0 OMITS it, 5.7 ECHOES it — must round-trip empty on both.
		{"action-no-action", fkParentSingle,
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id) ON DELETE NO ACTION ON UPDATE NO ACTION)", both()},
		// CASCADE / SET NULL: echoed verbatim on both. (pid must be nullable for SET NULL.)
		{"action-cascade-setnull", fkParentSingle,
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id) ON DELETE CASCADE ON UPDATE SET NULL)", both()},
		// ON DELETE only (ON UPDATE defaults/omitted).
		{"action-delete-only", fkParentSingle,
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id) ON DELETE CASCADE)", both()},
		// ON UPDATE only.
		{"action-update-only", fkParentSingle,
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id) ON UPDATE CASCADE)", both()},
		// SET NULL on both actions.
		{"action-setnull-both", fkParentSingle,
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id) ON DELETE SET NULL ON UPDATE SET NULL)", both()},

		// ---- columns / shape ----
		// Composite FK with its own auto-created composite backing index.
		{"composite", fkParentComposite,
			"CREATE TABLE c (id INT PRIMARY KEY, x INT, y INT, CONSTRAINT fk FOREIGN KEY (x,y) REFERENCES p(a,b))", both()},
		// FK reusing an explicit user index (distinct name) → no separate backing index.
		{"explicit-covering-index", fkParentSingle,
			"CREATE TABLE c (id INT PRIMARY KEY, pid INT, KEY my_idx (pid), CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))", both()},
		// Unnamed FK whose first-column index name collides with a user index → backing `pid_2`.
		{"backing-name-collision", fkParentSingle,
			"CREATE TABLE c (id INT PRIMARY KEY, a INT, pid INT, KEY pid (a), FOREIGN KEY (pid) REFERENCES p(id))", both()},
		// Multiple FKs on one child (two distinct parents' columns).
		{"multiple-fks", fkParentSingle,
			"CREATE TABLE c (id INT PRIMARY KEY, pa INT, pb INT, CONSTRAINT fka FOREIGN KEY (pa) REFERENCES p(a), CONSTRAINT fkb FOREIGN KEY (pb) REFERENCES p(b))", both()},
	}
}

// loadFKSchema applies parent+child in a throwaway db, then reloads BOTH tables wrapped in a FIXED
// logical db name so the from/to catalogs share a database identity. Returns the reloaded catalog.
// Skips the subtest if the engine rejects the setup DDL.
func loadFKSchema(t *testing.T, o *oracleConn, logicalDB, execDB, parentDDL, childDDL string) *Catalog {
	t.Helper()
	ctx := context.Background()
	sc := serverCharsetFor(o.version)
	stmts := []string{
		"DROP DATABASE IF EXISTS " + execDB,
		fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s", execDB, sc),
		"USE " + execDB,
		parentDDL,
		childDDL,
	}
	for _, s := range stmts {
		if _, err := o.db.ExecContext(ctx, s); err != nil {
			t.Skipf("[%s] setup failed (may be expected): %q: %v", o.name, s, err)
		}
	}
	var b strings.Builder
	fmt.Fprintf(&b, "CREATE DATABASE %s DEFAULT CHARSET=%s;\nUSE %s;\n", logicalDB, sc, logicalDB)
	for _, tbl := range []string{"p", "c"} {
		var nm, ddl string
		row := o.db.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE TABLE %s.%s", execDB, tbl))
		if err := row.Scan(&nm, &ddl); err != nil {
			t.Fatalf("[%s] SHOW CREATE %s: %v", o.name, tbl, err)
		}
		b.WriteString(ddl + ";\n")
	}
	cat, err := LoadSQL(b.String())
	if err != nil {
		t.Fatalf("[%s] reload of FK readback failed: %v\n%s", o.name, err, b.String())
	}
	return cat
}

// TestOracle_ForeignKeyDiffIdempotence proves gates 1 & 2 for every FK form: the engine's stored
// form self-diffs empty, AND the user form vs the engine's stored form diffs empty (action-default
// + auto-name normalization), on every supported version.
func TestOracle_ForeignKeyDiffIdempotence(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		n := NormalizerFor(version)
		for _, fc := range fkDiffCases() {
			if !containsVersion(fc.versions, version) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, fc.id), func(t *testing.T) {
				slug := strings.ReplaceAll(fc.id, "-", "_")
				// Stored form (engine readback) — the idempotence spine.
				storedCat := loadFKSchema(t, o, "fkdiff", "fkdiff_s_"+slug, fc.parentDDL, fc.childDDL)
				if d := DiffWithNormalizer(storedCat, storedCat, n); !d.IsEmpty() {
					t.Errorf("[%s] IDEMPOTENCE: FK stored-form self-diff not empty for %s: %s",
						o.name, fc.id, describeForeignKeyDiff(d))
				}

				// User form: load the raw user DDL (wrapped in the same logical db) and diff it
				// against the engine's stored form. This is the action-default / auto-name
				// canonicalization proof — they must collapse to empty.
				userCat := loadUserFKSchema(t, serverCharsetFor(o.version), fc.parentDDL, fc.childDDL)
				if d := DiffWithNormalizer(userCat, userCat, n); !d.IsEmpty() {
					t.Errorf("[%s] user-form FK self-diff not empty for %s: %s",
						o.name, fc.id, describeForeignKeyDiff(d))
				}
				if d := DiffWithNormalizer(userCat, storedCat, n); !d.IsEmpty() {
					t.Errorf("[%s] CANONICALIZATION: user form vs stored form not empty for %s: %s",
						o.name, fc.id, describeForeignKeyDiff(d))
				}
				if d := DiffWithNormalizer(storedCat, userCat, n); !d.IsEmpty() {
					t.Errorf("[%s] CANONICALIZATION (reverse): stored vs user not empty for %s: %s",
						o.name, fc.id, describeForeignKeyDiff(d))
				}
			})
		}
	}
}

// loadUserFKSchema loads the raw user parent+child DDL (no engine round-trip) wrapped in the same
// FIXED logical db as loadFKSchema, so the user and stored catalogs share a database identity and
// only the FK canonicalization is under test.
func loadUserFKSchema(t *testing.T, serverCharset, parentDDL, childDDL string) *Catalog {
	t.Helper()
	wrapped := fmt.Sprintf("CREATE DATABASE fkdiff DEFAULT CHARSET=%s;\nUSE fkdiff;\n%s;\n%s;",
		serverCharset, parentDDL, childDDL)
	cat, err := LoadSQL(wrapped)
	if err != nil {
		t.Fatalf("LoadSQL failed for user FK schema: %v\n%s", err, wrapped)
	}
	return cat
}

// describeForeignKeyDiff renders a compact description of a SchemaDiff focused on FK changes for
// failure output.
func describeForeignKeyDiff(d *SchemaDiff) string {
	var b strings.Builder
	for _, te := range d.Tables {
		fmt.Fprintf(&b, "[table %s %s", te.Name, te.Action)
		for _, ce := range te.Columns {
			fmt.Fprintf(&b, " col(%s %s)", ce.Name, ce.Action)
		}
		for _, ie := range te.Indexes {
			fmt.Fprintf(&b, " idx(%s %s)", ie.Name, ie.Action)
		}
		for _, fe := range te.ForeignKeys {
			fmt.Fprintf(&b, " fk(%s %s)", fe.Name, fe.Action)
		}
		b.WriteString("]")
	}
	return b.String()
}

// TestOracle_ForeignKeyDiffNonEmptyCorrect proves the differ DISTINGUISHES genuinely different FK
// schemas (the dual of idempotence — a missed FK diff is as harmful as a phantom one). Each case
// loads two real engine readbacks and asserts the expected FK DiffAction.
func TestOracle_ForeignKeyDiffNonEmptyCorrect(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		n := NormalizerFor(version)

		cases := []struct {
			name      string
			fromChild string
			toChild   string
			wantFK    string
			wantFKA   DiffAction
		}{
			{"add-fk",
				"CREATE TABLE c (id INT PRIMARY KEY, pid INT)",
				"CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))",
				"fk", DiffAdd},
			{"drop-fk",
				"CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id))",
				"CREATE TABLE c (id INT PRIMARY KEY, pid INT)",
				"fk", DiffDrop},
			// Action change CASCADE→SET NULL is a real change (both echoed verbatim on both versions).
			{"modify-action",
				"CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id) ON DELETE CASCADE)",
				"CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(id) ON DELETE SET NULL)",
				"fk", DiffModify},
			// Referenced-column change a→b is a real change.
			{"modify-refcol",
				"CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(a))",
				"CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p(b))",
				"fk", DiffModify},
		}
		for _, tc := range cases {
			t.Run(o.name+"/"+tc.name, func(t *testing.T) {
				slug := strings.ReplaceAll(tc.name, "-", "_")
				from := loadFKSchema(t, o, "fknec", "fknec_f_"+slug, fkParentSingle, tc.fromChild)
				to := loadFKSchema(t, o, "fknec", "fknec_t_"+slug, fkParentSingle, tc.toChild)
				d := DiffWithNormalizer(from, to, n)
				te := findTable(d, "c")
				if te == nil {
					t.Fatalf("[%s] want table c modified, got %s", o.name, describeForeignKeyDiff(d))
				}
				var got *ForeignKeyDiffEntry
				for i := range te.ForeignKeys {
					if strings.EqualFold(te.ForeignKeys[i].Name, tc.wantFK) {
						got = &te.ForeignKeys[i]
						break
					}
				}
				if got == nil || got.Action != tc.wantFKA {
					t.Fatalf("[%s] want fk %s %s, got %s", o.name, tc.wantFK, tc.wantFKA, describeForeignKeyDiff(d))
				}
			})
		}
	}
}
