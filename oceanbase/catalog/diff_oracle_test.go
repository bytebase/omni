package catalog

import (
	"context"
	"fmt"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// The oracle proof for diff-core (correctness-protocol.md gates 1 & 2), against
// OceanBase. Two properties are proven mechanically:
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
// (oceanBaseOracle, loadColumn's wrapping convention, serverCharsetFor)., so the unit suite stays hermetic.

// applyAndReadback applies a single CREATE TABLE in a throwaway database and returns the
// SHOW CREATE TABLE readback (the engine's canonical stored form). It reuses oracleConn
// from normalize_oracle_test.go.
func (o *oracleConn) applyAndReadback(t testing.TB, dbName, createSQL, table string) (string, bool) {
	t.Helper()
	return o.showCreate(t, dbName, createSQL, table)
}

// loadOneTable loads a CREATE TABLE (wrapped in a database whose default charset matches
// the oracle box, so table-charset inheritance resolves identically on both sides) and
// returns the named table.
func loadOneTable(t testing.TB, serverCharset, createSQL, table string) (*Catalog, *Table) {
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
// all under a version-correct Normalizer. A form the engine will not take has no stored
// form: the probe ends in refused.
func assertDiffEmptyAgainstReadback(t testing.TB, o *oracleConn, dbName, userDDL, table string) {
	t.Helper()
	readback, ok := o.applyAndReadback(t, dbName, userDDL, table)
	if !ok {
		o.refused(t, nil, "could not obtain readback for %s", table)
		return
	}
	assertDiffEmptyAgainstStored(t, o, readback, userDDL, table)
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
	id     string
	create string
	table  string
}

// diffIdempotenceProbes enumerates representative table/column FORMS that must round-trip
// empty: user-form DDL whose engine-stored form differs but must canonicalize equal.
// These exercise the high-risk normalization surfaces through the full diff path.
func diffIdempotenceProbes() []diffProbe {
	return []diffProbe{
		// Integer display width: 8.0 drops, 5.7 injects default — must round-trip on both.
		{"int-widths", "CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT(11), b BIGINT(20), c INT, p TINYINT, h INT UNSIGNED)", "t"},
		// BOOLEAN/BOOL -> tinyint(1).
		{"boolean", "CREATE TABLE t (id INT PRIMARY KEY, m BOOLEAN, n BOOL, d TINYINT(1))", "t"},
		// decimal/float aliasing + default precision.
		{"numeric", "CREATE TABLE t (id INT PRIMARY KEY, a DECIMAL, b DECIMAL(10,2), d NUMERIC, j REAL, f FLOAT)", "t"},
		// char/binary/bit length-1, year width.
		{"char-bit-year", "CREATE TABLE t (id INT PRIMARY KEY, a CHAR, d BINARY, c BIT, y YEAR)", "t"},
		// literal default quoting (numeric + string), decimal scale padding, boolean default.
		{"defaults", "CREATE TABLE t (id INT PRIMARY KEY, a INT DEFAULT 0, b INT DEFAULT '7', f DECIMAL(10,2) DEFAULT 0, c VARCHAR(20) DEFAULT 'x', j BOOLEAN DEFAULT TRUE)", "t"},
		// bug: BIGINT UNSIGNED max default (18446744073709551615) — must not int64-truncate;
		// the unquoted user form and the quoted stored readback ('18446744073709551615') must
		// diff empty. mid (9223372036854775808 = int64-max+1) also overflows int64.
		{"bigint-unsigned-max-default", "CREATE TABLE t (id INT PRIMARY KEY, big BIGINT UNSIGNED NOT NULL DEFAULT 18446744073709551615, mid BIGINT UNSIGNED NOT NULL DEFAULT 9223372036854775808)", "t"},
		// bug: signed-bigint boundaries — int64-min default (-9223372036854775808) must not
		// render off-by-one; int64-max round-trips too.
		{"bigint-signed-boundary-default", "CREATE TABLE t (id INT PRIMARY KEY, lo BIGINT NOT NULL DEFAULT -9223372036854775808, hi BIGINT NOT NULL DEFAULT 9223372036854775807)", "t"},
		// bug: BIT column default — bit-literal / hex / bare-int / quoted-byte-string forms are
		// all stored as b'<binary>' and must canonicalize by value (b'0' is the stated case).
		{"bit-default-forms", "CREATE TABLE t (id INT PRIMARY KEY, a BIT(8) NOT NULL DEFAULT b'0', b BIT(8) NOT NULL DEFAULT 0x05, c BIT(8) NOT NULL DEFAULT 5, d BIT(16) NOT NULL DEFAULT 'AB')", "t"},
		// bug: YEAR default — stored single-quoted ('2000') but written numeric (2000);
		// two-digit YEAR is expanded at storage (99 -> '1999'), 0 stays '0000'.
		{"year-default", "CREATE TABLE t (id INT PRIMARY KEY, a YEAR NOT NULL DEFAULT 2000, b YEAR NOT NULL DEFAULT 1999, c YEAR NOT NULL DEFAULT 99, d YEAR NOT NULL DEFAULT 0)", "t"},
		// nullability collapse to DEFAULT NULL; PK forces NOT NULL.
		{"nullability", "CREATE TABLE t (id INT, a INT NULL, b INT DEFAULT NULL, d VARCHAR(10) NULL, PRIMARY KEY (id))", "t"},
		// enum/set quoting + order.
		{"enum-set", `CREATE TABLE t (id INT PRIMARY KEY, a ENUM('x','y','z'), b SET('a','b','c'), e ENUM("dq1","dq2"))`, "t"},
		// generated column expression normalization.
		{"generated", "CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT GENERATED ALWAYS AS (a+1) VIRTUAL, c INT GENERATED ALWAYS AS ( a * 2 ) STORED, d INT AS (a+1))", "t"},
		// implicit table charset/engine (no clause): must not phantom an ADD CHARSET/ENGINE.
		{"implicit-table-opts", "CREATE TABLE t (id INT PRIMARY KEY, a VARCHAR(10))", "t"},
		// explicit table charset + column echo (version-flagged echo rules).
		{"charset-echo", "CREATE TABLE t (id INT PRIMARY KEY, a VARCHAR(10) CHARACTER SET utf8mb4, b VARCHAR(10)) DEFAULT CHARSET=utf8mb4", "t"},
		// bug (Roundcube et al.): the legacy `BINARY` string-column modifier resolves to the
		// charset's `_bin` collation (stored as CHARACTER SET <cs> COLLATE <cs>_bin on both
		// versions). The user's `varchar(n) BINARY` form must diff EMPTY against that stored
		// pair — the regression was a phantom MODIFY COLUMN on every no-op. Covers varchar/char/
		// text/tinytext on a utf8mb4 table (-> utf8mb4_bin) plus an explicit ascii column
		// (-> ascii_bin). (entry char-binary-attribute)
		{"binary-modifier-utf8mb4", "CREATE TABLE t (id INT PRIMARY KEY, k VARCHAR(128) BINARY NOT NULL, c CHAR(5) BINARY, tx TEXT BINARY, ti TINYTEXT BINARY) DEFAULT CHARSET=utf8mb4", "t"},
		// OceanBase's grammar puts BINARY before the charset (`VARCHAR string_length_i
		// opt_binary opt_charset`): it takes the first of these and refuses the second
		// (obRefused), and the parser reads them the same way.
		{"binary-modifier-before-charset", "CREATE TABLE t (id INT PRIMARY KEY, ac VARCHAR(10) BINARY CHARACTER SET ascii) DEFAULT CHARSET=utf8mb4", "t"},
		{"binary-modifier-after-charset", "CREATE TABLE t (id INT PRIMARY KEY, ac VARCHAR(10) CHARACTER SET ascii BINARY) DEFAULT CHARSET=utf8mb4", "t"},
		// latin1 table -> latin1_bin.
		{"binary-modifier-latin1", "CREATE TABLE t (id INT PRIMARY KEY, a VARCHAR(50) BINARY) DEFAULT CHARSET=latin1", "t"},
		// comment escaping (table COMMENT carries an embedded single-quoted token `'c'` via
		// doubled quotes; the closing quote balances the literal — the previous fixture dropped it,
		// so MySQL rejected the DDL and the probe silently skipped on both engines).
		{"comment", "CREATE TABLE t (id INT PRIMARY KEY, a INT COMMENT 'has ''quote'' inside') COMMENT='tbl ''c'''", "t"},
		// AUTO_INCREMENT column attribute (real schema) + table counter (ignored).
		{"autoinc", "CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY, a INT) AUTO_INCREMENT=100", "t"},
		// TIMESTAMP magic — opposite behavior across versions (EDFT box default).
		{"timestamp", "CREATE TABLE t (id INT PRIMARY KEY, a TIMESTAMP, b TIMESTAMP NULL, d TIMESTAMP DEFAULT CURRENT_TIMESTAMP)", "t"},
		{"datetime-default", "CREATE TABLE t (id INT PRIMARY KEY, g DATETIME DEFAULT CURRENT_TIMESTAMP, h DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)", "t"},
		// functional / expression default — 8.0 only.
		{"functional-default", "CREATE TABLE t (id INT PRIMARY KEY, a INT DEFAULT (1+1), b VARCHAR(36) DEFAULT (UUID()))", "t"},
		// 8.0 expression string introducer (_latin1'x') inside a generated expr.
		{"gen-introducer", "CREATE TABLE t (id INT PRIMARY KEY, a VARCHAR(20), e VARCHAR(20) GENERATED ALWAYS AS (CONCAT(a,'x')) VIRTUAL) DEFAULT CHARSET=utf8mb4", "t"},
		// multi-column ordering preserved.
		{"multi-column", "CREATE TABLE t (id INT NOT NULL, c1 VARCHAR(10), c2 INT DEFAULT 0, c3 DATE, c4 DECIMAL(8,3), PRIMARY KEY (id))", "t"},
		// explicit ROW_FORMAT round-trips (SHOW CREATE echoes the declared option, so the
		// user form and the readback both carry it → empty diff). Idempotence guard for the
		// rowFormatChanged gate.
		{"row-format-explicit", "CREATE TABLE t (id INT PRIMARY KEY) ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8", "t"},
		// adjacent string-literal concatenation in a bare column DEFAULT: the user form
		// DEFAULT 'x' 'y' ('x' 'y' 'z', "x" 'y') is stored as the folded DEFAULT 'xy'
		// ('xyz', 'xy') — the user form must diff EMPTY against that readback.
		{"adjacent-literal-default", `CREATE TABLE t (id INT PRIMARY KEY, a VARCHAR(20) DEFAULT 'x' 'y', b VARCHAR(20) DEFAULT 'x' 'y' 'z')`, "t"},
		// OceanBase folds single-quoted neighbours, as above, and refuses a double-quoted
		// literal with a neighbour.
		{"adjacent-literal-default-double-quoted", `CREATE TABLE t (id INT PRIMARY KEY, c VARCHAR(20) DEFAULT "x" 'y')`, "t"},
		// adjacency inside LIST COLUMNS partition values: VALUES IN ('a' 'b') is stored
		// folded as VALUES IN ('ab').
		{"adjacent-literal-partition", "CREATE TABLE t (c VARCHAR(10) NOT NULL PRIMARY KEY) PARTITION BY LIST COLUMNS(c) (PARTITION p0 VALUES IN ('a' 'b'), PARTITION p1 VALUES IN ('z'))", "t"},
		// charset-introduced string defaults, single and folded: the introducer is not
		// part of the stored default (SHOW CREATE echoes DEFAULT 'xy' / DEFAULT 'z'), so
		// the user form must diff EMPTY against the readback and never render
		// DEFAULT '_utf8mb4'xy'' (invalid SQL, the codex-review P2 case).
		{"adjacent-literal-default-introducer", "CREATE TABLE t (id INT PRIMARY KEY, a VARCHAR(20) CHARACTER SET utf8mb4 DEFAULT _utf8mb4'x' 'y', b VARCHAR(20) CHARACTER SET utf8mb4 DEFAULT _utf8mb4'z') DEFAULT CHARSET=utf8mb4", "t"},
	}
}

// TestOracle_DiffIdempotence proves gate 1 & 2 for every form: user DDL vs its engine
// readback diffs EMPTY, and the stored form self-diffs empty, on every supported version.
func TestOracle_DiffIdempotence(t *testing.T) {
	for range oracleVersions() {
		o := oceanBaseOracle(t)
		for _, probe := range diffIdempotenceProbes() {
			t.Run(fmt.Sprintf("%s/%s", o.name, probe.id), gapped(func(t testing.TB) {
				db := "diff_" + strings.ReplaceAll(probe.id, "-", "_")
				assertDiffEmptyAgainstReadback(t, o, db, probe.create, probe.table)
			}))
		}
	}
}

// TestOracle_DiffMultiTableRoundTrip proves a multi-table schema (with cross-table FKs
// and an index) loaded from its real engine readbacks self-diffs empty — the realistic
// release-path idempotence check, not just single-table forms.
func TestOracle_DiffMultiTableRoundTrip(t *testing.T) {
	for range oracleVersions() {
		o := oceanBaseOracle(t)
		sc := serverCharsetFor(o.version)
		n := NormalizerFor(o.version)
		ctx := context.Background()

		t.Run(o.name+"/two-tables-fk", gapped(func(t testing.TB) {
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
					o.refused(t, err, "setup failed: %q", s)
					return
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
		}))
	}
}

// TestOracle_DiffNonEmptyCorrect proves gate 3 against the live engine: representative
// (from, to) pairs loaded from REAL engine readbacks produce the correct DiffAction.
// (Apply-correctness — that the generated DDL transforms from into to — is generate-
// core's gate; here we assert the SchemaDiff content is structurally right.)
func TestOracle_DiffNonEmptyCorrect(t *testing.T) {
	for range oracleVersions() {
		o := oceanBaseOracle(t)
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
			t.Run(o.name+"/"+tc.name, gapped(func(t testing.TB) {
				rbFrom, ok1 := o.showCreate(t, "diff_nf", tc.fromDDL, tc.wantTable)
				rbTo, ok2 := o.showCreate(t, "diff_nt", tc.toDDL, tc.wantTable)
				if !ok1 || !ok2 {
					o.refused(t, nil, "could not obtain readbacks")
					return
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
			}))
		}
	}
}
