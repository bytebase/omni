package catalog

import (
	"context"
	"fmt"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// The oracle proof for diff-core (correctness-protocol.md gates 1 & 2), against the
// LIVE MySQL engines (5.7 and 8.0). Two properties are proven mechanically:
//
//  1. Idempotence (the spine): for a schema in its real stored form (the engine's own
//     SHOW CREATE readback, loaded into a catalog), Diff(c, c).IsEmpty() — and, more
//     strongly, Diff(userForm, storedForm).IsEmpty() where storedForm is the engine's
//     readback of userForm. A non-empty self/round-trip diff is a normalization bug.
//
//  2. Canonicalization-driven empties: a schema written in a USER form (INT(11),
//     BOOLEAN, implicit charset, DEFAULT 0, generated-col spacing, ...) diffed against
//     the engine's STORED form of the same schema must be EMPTY. This is the proof that
//     diff-core routes column equality through normalize-core's CanonicalColumn — if it
//     compared surface tokens, these would phantom-diff forever.
//
// The harness reuses the connection + load helpers from normalize_oracle_test.go
// (connectOracle, loadColumn's wrapping convention, serverCharsetFor). It skips cleanly
// when the engines are unreachable, so the unit suite stays hermetic.

// applyAndReadback applies a single CREATE TABLE in a throwaway database and returns the
// SHOW CREATE TABLE readback (the engine's canonical stored form). It reuses oracleConn
// from normalize_oracle_test.go.
func (o *oracleConn) applyAndReadback(t *testing.T, dbName, createSQL, table string) (string, bool) {
	t.Helper()
	return o.showCreate(t, dbName, createSQL, table)
}

// loadOneTable loads a CREATE TABLE (wrapped in a database whose default charset matches
// the oracle box, so table-charset inheritance resolves identically on both sides) and
// returns the named table.
func loadOneTable(t *testing.T, serverCharset, createSQL, table string) (*Catalog, *Table) {
	t.Helper()
	wrapped := fmt.Sprintf("CREATE DATABASE diffdb DEFAULT CHARSET=%s;\nUSE diffdb;\n%s", serverCharset, createSQL)
	cat, err := LoadSQL(wrapped)
	if err != nil {
		t.Fatalf("LoadSQL failed for %q: %v", createSQL, err)
	}
	var tbl *Table
	for _, db := range cat.Databases() {
		if tt := db.GetTable(table); tt != nil {
			tbl = tt
			break
		}
	}
	if tbl == nil {
		t.Fatalf("table %q not found after load of %q", table, createSQL)
	}
	return cat, tbl
}

// assertDiffEmptyAgainstReadback applies userDDL to the engine, reads back the stored
// form, loads BOTH the user form and the stored form, and asserts:
//   - Diff(userCat, storedCat) is empty  (canonicalization-driven, gate 2)
//   - Diff(storedCat, storedCat) is empty (idempotence spine, gate 1)
//   - Diff(userCat, userCat) is empty     (user form self-consistency)
//
// all under a version-correct Normalizer.
func assertDiffEmptyAgainstReadback(t *testing.T, o *oracleConn, dbName, userDDL, table string) {
	t.Helper()
	readback, ok := o.applyAndReadback(t, dbName, userDDL, table)
	if !ok {
		t.Skipf("[%s] could not obtain readback for %s", o.name, table)
	}
	sc := serverCharsetFor(o.version)
	n := NormalizerFor(o.version)

	userCat, _ := loadOneTable(t, sc, userDDL, table)
	storedCat, _ := loadOneTable(t, sc, readback, table)

	if d := DiffWithNormalizer(storedCat, storedCat, n); !d.IsEmpty() {
		t.Errorf("[%s] IDEMPOTENCE: self-diff of stored form not empty for %s:\n  stored: %s\n  diff: %s",
			o.name, table, strings.TrimSpace(readback), describeDiff(d))
	}
	if d := DiffWithNormalizer(userCat, userCat, n); !d.IsEmpty() {
		t.Errorf("[%s] self-diff of user form not empty for %s:\n  user: %s\n  diff: %s",
			o.name, table, strings.TrimSpace(userDDL), describeDiff(d))
	}
	// The headline property: user form vs engine-stored form must collapse to empty.
	if d := DiffWithNormalizer(userCat, storedCat, n); !d.IsEmpty() {
		t.Errorf("[%s] CANONICALIZATION: user form vs stored form not empty for %s:\n  user:   %s\n  stored: %s\n  diff: %s",
			o.name, table, strings.TrimSpace(userDDL), strings.TrimSpace(readback), describeDiff(d))
	}
	// Symmetry: stored vs user must also be empty (Drop/Add are direction-sensitive).
	if d := DiffWithNormalizer(storedCat, userCat, n); !d.IsEmpty() {
		t.Errorf("[%s] CANONICALIZATION (reverse): stored vs user not empty for %s:\n  diff: %s",
			o.name, table, describeDiff(d))
	}
}

// describeDiff renders a compact human description of a SchemaDiff for failure output.
func describeDiff(d *SchemaDiff) string {
	var b strings.Builder
	for _, te := range d.Tables {
		fmt.Fprintf(&b, "[table %s %s", te.Name, te.Action)
		for _, ce := range te.Columns {
			fmt.Fprintf(&b, " col(%s %s)", ce.Name, ce.Action)
		}
		b.WriteString("]")
	}
	return b.String()
}

// diffProbe is one (id, userDDL, table, versions) idempotence/canonicalization case.
type diffProbe struct {
	id       string
	create   string
	table    string
	versions []Version
}

// diffIdempotenceProbes enumerates representative table/column FORMS that must round-trip
// empty: user-form DDL whose engine-stored form differs but must canonicalize equal.
// These exercise the high-risk normalization surfaces through the full diff path.
func diffIdempotenceProbes() []diffProbe {
	return []diffProbe{
		// Integer display width: 8.0 drops, 5.7 injects default — must round-trip on both.
		{"int-widths", "CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT(11), b BIGINT(20), c INT, p TINYINT, h INT UNSIGNED)", "t", both()},
		// BOOLEAN/BOOL -> tinyint(1).
		{"boolean", "CREATE TABLE t (id INT PRIMARY KEY, m BOOLEAN, n BOOL, d TINYINT(1))", "t", both()},
		// decimal/float aliasing + default precision.
		{"numeric", "CREATE TABLE t (id INT PRIMARY KEY, a DECIMAL, b DECIMAL(10,2), d NUMERIC, j REAL, f FLOAT)", "t", both()},
		// char/binary/bit length-1, year width.
		{"char-bit-year", "CREATE TABLE t (id INT PRIMARY KEY, a CHAR, d BINARY, c BIT, y YEAR)", "t", both()},
		// literal default quoting (numeric + string), decimal scale padding, boolean default.
		{"defaults", "CREATE TABLE t (id INT PRIMARY KEY, a INT DEFAULT 0, b INT DEFAULT '7', f DECIMAL(10,2) DEFAULT 0, c VARCHAR(20) DEFAULT 'x', j BOOLEAN DEFAULT TRUE)", "t", both()},
		// nullability collapse to DEFAULT NULL; PK forces NOT NULL.
		{"nullability", "CREATE TABLE t (id INT, a INT NULL, b INT DEFAULT NULL, d VARCHAR(10) NULL, PRIMARY KEY (id))", "t", both()},
		// enum/set quoting + order.
		{"enum-set", `CREATE TABLE t (id INT PRIMARY KEY, a ENUM('x','y','z'), b SET('a','b','c'), e ENUM("dq1","dq2"))`, "t", both()},
		// generated column expression normalization.
		{"generated", "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT GENERATED ALWAYS AS (a+1) VIRTUAL, c INT GENERATED ALWAYS AS ( a * 2 ) STORED, d INT AS (a+1))", "t", both()},
		// implicit table charset/engine (no clause): must not phantom an ADD CHARSET/ENGINE.
		{"implicit-table-opts", "CREATE TABLE t (id INT PRIMARY KEY, a VARCHAR(10))", "t", both()},
		// explicit table charset + column echo (version-flagged echo rules).
		{"charset-echo", "CREATE TABLE t (id INT PRIMARY KEY, a VARCHAR(10) CHARACTER SET utf8mb4, b VARCHAR(10)) DEFAULT CHARSET=utf8mb4", "t", both()},
		// comment escaping (table COMMENT carries an embedded single-quoted token `'c'` via
		// doubled quotes; the closing quote balances the literal — the previous fixture dropped it,
		// so MySQL rejected the DDL and the probe silently skipped on both engines).
		{"comment", "CREATE TABLE t (id INT PRIMARY KEY, a INT COMMENT 'has ''quote'' inside') COMMENT='tbl ''c'''", "t", both()},
		// AUTO_INCREMENT column attribute (real schema) + table counter (ignored).
		{"autoinc", "CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY, a INT) AUTO_INCREMENT=100", "t", both()},
		// TIMESTAMP magic — opposite behavior across versions (EDFT box default).
		{"timestamp", "CREATE TABLE t (id INT PRIMARY KEY, a TIMESTAMP, b TIMESTAMP NULL, d TIMESTAMP DEFAULT CURRENT_TIMESTAMP)", "t", both()},
		{"datetime-default", "CREATE TABLE t (id INT PRIMARY KEY, g DATETIME DEFAULT CURRENT_TIMESTAMP, h DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)", "t", both()},
		// functional / expression default — 8.0 only.
		{"functional-default", "CREATE TABLE t (id INT PRIMARY KEY, a INT DEFAULT (1+1), b VARCHAR(36) DEFAULT (UUID()))", "t", only(MySQL80)},
		// 8.0 expression string introducer (_latin1'x') inside a generated expr.
		{"gen-introducer", "CREATE TABLE t (id INT PRIMARY KEY, a VARCHAR(20), e VARCHAR(20) GENERATED ALWAYS AS (CONCAT(a,'x')) VIRTUAL) DEFAULT CHARSET=utf8mb4", "t", only(MySQL80)},
		// multi-column ordering preserved.
		{"multi-column", "CREATE TABLE t (id INT NOT NULL, c1 VARCHAR(10), c2 INT DEFAULT 0, c3 DATE, c4 DECIMAL(8,3), PRIMARY KEY (id))", "t", both()},
		// explicit ROW_FORMAT round-trips (SHOW CREATE echoes the declared option, so the
		// user form and the readback both carry it → empty diff). Idempotence guard for the
		// rowFormatChanged gate.
		{"row-format-explicit", "CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8", "t", both()},
	}
}

// TestOracle_DiffIdempotence proves gate 1 & 2 for every form: user DDL vs its engine
// readback diffs EMPTY, and the stored form self-diffs empty, on every supported version.
func TestOracle_DiffIdempotence(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		for _, probe := range diffIdempotenceProbes() {
			if !containsVersion(probe.versions, version) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, probe.id), func(t *testing.T) {
				db := "diff_" + strings.ReplaceAll(probe.id, "-", "_")
				assertDiffEmptyAgainstReadback(t, o, db, probe.create, probe.table)
			})
		}
	}
}

// TestOracle_DiffMultiTableRoundTrip proves a multi-table schema (with cross-table FKs
// and an index) loaded from its real engine readbacks self-diffs empty — the realistic
// release-path idempotence check, not just single-table forms.
func TestOracle_DiffMultiTableRoundTrip(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		sc := serverCharsetFor(o.version)
		n := NormalizerFor(o.version)
		ctx := context.Background()

		t.Run(o.name+"/two-tables-fk", func(t *testing.T) {
			dbName := "diff_multi"
			stmts := []string{
				fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName),
				fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s", dbName, sc),
				fmt.Sprintf("USE %s", dbName),
				"CREATE TABLE parent (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY, name VARCHAR(50) NOT NULL)",
				"CREATE TABLE child (id INT NOT NULL PRIMARY KEY, pid INT, amount DECIMAL(10,2) DEFAULT 0, KEY pid_idx (pid), CONSTRAINT fk_pid FOREIGN KEY (pid) REFERENCES parent (id))",
			}
			for _, s := range stmts {
				if _, err := o.db.ExecContext(ctx, s); err != nil {
					t.Skipf("[%s] setup failed (may be expected): %q: %v", o.name, s, err)
				}
			}
			// Read back both tables and reload as a single schema.
			var b strings.Builder
			fmt.Fprintf(&b, "CREATE DATABASE %s DEFAULT CHARSET=%s;\nUSE %s;\n", dbName, sc, dbName)
			for _, tbl := range []string{"parent", "child"} {
				var name, ddl string
				row := o.db.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE TABLE %s.%s", dbName, tbl))
				if err := row.Scan(&name, &ddl); err != nil {
					t.Fatalf("[%s] SHOW CREATE %s: %v", o.name, tbl, err)
				}
				b.WriteString(ddl)
				b.WriteString(";\n")
			}
			cat, err := LoadSQL(b.String())
			if err != nil {
				t.Fatalf("[%s] reload of multi-table readback failed: %v\n%s", o.name, err, b.String())
			}
			if d := DiffWithNormalizer(cat, cat, n); !d.IsEmpty() {
				t.Errorf("[%s] multi-table self-diff not empty: %s", o.name, describeDiff(d))
			}
		})
	}
}

// TestOracle_DiffNonEmptyCorrect proves gate 3 against the live engine: representative
// (from, to) pairs loaded from REAL engine readbacks produce the correct DiffAction.
// (Apply-correctness — that the generated DDL transforms from into to — is generate-
// core's gate; here we assert the SchemaDiff content is structurally right.)
func TestOracle_DiffNonEmptyCorrect(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		sc := serverCharsetFor(o.version)
		n := NormalizerFor(o.version)

		// Each case: two DDLs applied to the engine; we diff their readbacks and assert.
		cases := []struct {
			name       string
			fromDDL    string
			toDDL      string
			wantTable  string
			wantTableA DiffAction
			wantCol    string
			wantColA   DiffAction
		}{
			{"add-column",
				"CREATE TABLE t (id INT PRIMARY KEY)",
				"CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(20))",
				"t", DiffModify, "name", DiffAdd},
			{"drop-column",
				"CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(20))",
				"CREATE TABLE t (id INT PRIMARY KEY)",
				"t", DiffModify, "name", DiffDrop},
			{"modify-type",
				"CREATE TABLE t (id INT PRIMARY KEY, v INT)",
				"CREATE TABLE t (id INT PRIMARY KEY, v BIGINT)",
				"t", DiffModify, "v", DiffModify},
			{"modify-default",
				"CREATE TABLE t (id INT PRIMARY KEY, v INT DEFAULT 0)",
				"CREATE TABLE t (id INT PRIMARY KEY, v INT DEFAULT 1)",
				"t", DiffModify, "v", DiffModify},
			{"modify-nullability",
				"CREATE TABLE t (id INT PRIMARY KEY, v INT NULL)",
				"CREATE TABLE t (id INT PRIMARY KEY, v INT NOT NULL)",
				"t", DiffModify, "v", DiffModify},
		}
		for _, tc := range cases {
			t.Run(o.name+"/"+tc.name, func(t *testing.T) {
				rbFrom, ok1 := o.showCreate(t, "diff_nf", tc.fromDDL, tc.wantTable)
				rbTo, ok2 := o.showCreate(t, "diff_nt", tc.toDDL, tc.wantTable)
				if !ok1 || !ok2 {
					t.Skipf("[%s] could not obtain readbacks", o.name)
				}
				fromCat, _ := loadOneTable(t, sc, rbFrom, tc.wantTable)
				toCat, _ := loadOneTable(t, sc, rbTo, tc.wantTable)
				d := DiffWithNormalizer(fromCat, toCat, n)
				te := findTable(d, tc.wantTable)
				if te == nil || te.Action != tc.wantTableA {
					t.Fatalf("[%s] want table %s %s, got %s", o.name, tc.wantTable, tc.wantTableA, describeDiff(d))
				}
				ce := findColumn(te, tc.wantCol)
				if ce == nil || ce.Action != tc.wantColA {
					t.Fatalf("[%s] want col %s %s, got %s", o.name, tc.wantCol, tc.wantColA, describeDiff(d))
				}
			})
		}
	}
}
