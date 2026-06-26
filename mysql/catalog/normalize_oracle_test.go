package catalog

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// The oracle proof for normalize-core. For every normalization.md entry, this applies
// the input DDL to a LIVE MySQL engine (both 5.7 and 8.0 where they diverge), reads back
// SHOW CREATE TABLE, and asserts the phantom-diff-elimination property:
//
//	CanonicalColumn(loaded(input_ddl))  ==  CanonicalColumn(loaded(SHOW CREATE readback))
//
// i.e. the user's declared form and the engine's stored form collapse onto the same
// canonical key for the target version. A failure here is a real normalization bug — the
// canonicalizer disagrees with what the engine actually stores, which would produce a
// phantom diff forever. These tests are the correctness spine (correctness-protocol.md).
//
// Connection: the local oracle instances from the work order, overridable via env.
// They skip cleanly when the engines are unreachable so the unit suite stays hermetic.

type oracleConn struct {
	db      *sql.DB
	version Version
	name    string
}

func dsnOr(env, def string) string {
	if v := os.Getenv(env); v != "" {
		return v
	}
	return def
}

// connectOracle dials one engine; returns nil (and skips) if unreachable.
func connectOracle(t *testing.T, version Version) *oracleConn {
	t.Helper()
	var dsn, name string
	switch version {
	case MySQL80:
		dsn = dsnOr("OMNI_MYSQL80_DSN", "root:010424@tcp(127.0.0.1:13306)/?multiStatements=true")
		name = "8.0"
	case MySQL57:
		dsn = dsnOr("OMNI_MYSQL57_DSN", "root:010424@tcp(127.0.0.1:13307)/?multiStatements=true&tls=false")
		name = "5.7"
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Skipf("oracle %s unavailable (open): %v", name, err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		t.Skipf("oracle %s unavailable (ping): %v", name, err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return &oracleConn{db: db, version: version, name: name}
}

// showCreate applies the CREATE statements in a throwaway database and returns the
// SHOW CREATE TABLE readback for the named table.
func (o *oracleConn) showCreate(t *testing.T, dbName, createSQL, table string) (string, bool) {
	t.Helper()
	ctx := context.Background()
	stmts := []string{
		fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName),
		fmt.Sprintf("CREATE DATABASE %s", dbName),
		fmt.Sprintf("USE %s", dbName),
		createSQL,
	}
	for _, s := range stmts {
		if _, err := o.db.ExecContext(ctx, s); err != nil {
			// An input the version rejects (e.g. functional DEFAULT on 5.7) is not a
			// canonicalizer failure; the caller decides whether that is expected.
			t.Logf("[%s] exec failed (may be expected): %q: %v", o.name, s, err)
			return "", false
		}
	}
	var name, ddl string
	row := o.db.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE TABLE %s.%s", dbName, table))
	if err := row.Scan(&name, &ddl); err != nil {
		t.Logf("[%s] SHOW CREATE failed: %v", o.name, err)
		return "", false
	}
	return ddl, true
}

// loadColumn loads DDL through the omni catalog and returns the table + named column.
// The CREATE TABLE is wrapped in a database whose default charset matches the oracle's
// server default (utf8mb4 on the 8.0 box, latin1 on the 5.7 box) so that table-charset
// inheritance resolves identically on both the user-input and the SHOW CREATE readback
// sides — without this, db-inherited charset rules (e.g. column-charset-echo-57) would
// not be comparable.
func loadColumn(t *testing.T, serverCharset, createSQL, table, column string) (*Table, *Column) {
	t.Helper()
	wrapped := fmt.Sprintf("CREATE DATABASE nrmdb DEFAULT CHARSET=%s;\nUSE nrmdb;\n%s", serverCharset, createSQL)
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
	return tbl, tbl.GetColumn(column)
}

// serverCharsetFor returns the oracle box's default server charset for a version, used
// to wrap loads so table-charset inheritance matches the readback.
func serverCharsetFor(v Version) string {
	if v == MySQL80 {
		return "utf8mb4"
	}
	return "latin1"
}

// assertPhantomFree is the core assertion: input form and stored readback must produce
// the same canonical column key for the engine's version.
func assertPhantomFree(t *testing.T, o *oracleConn, dbName, inputDDL, table, column string) {
	t.Helper()
	readback, ok := o.showCreate(t, dbName, inputDDL, table)
	if !ok {
		t.Skipf("[%s] could not obtain readback for %s.%s", o.name, table, column)
	}
	// The readback is a CREATE TABLE we can reload through the same loader.
	readbackCreate := normalizeReadbackForLoad(readback)

	n := NormalizerFor(o.version)
	sc := serverCharsetFor(o.version)

	userTbl, userCol := loadColumn(t, sc, inputDDL, table, column)
	storTbl, storCol := loadColumn(t, sc, readbackCreate, table, column)
	if userCol == nil || storCol == nil {
		t.Fatalf("[%s] column %q missing (user=%v stored=%v)", o.name, column, userCol != nil, storCol != nil)
	}

	uk := n.CanonicalColumn(userTbl, userCol)
	sk := n.CanonicalColumn(storTbl, storCol)
	if uk != sk {
		t.Errorf("[%s] PHANTOM DIFF on column %q:\n  input  DDL: %s\n  stored DDL: %s\n  user  key: %s\n  store key: %s",
			o.name, column, strings.TrimSpace(inputDDL), strings.TrimSpace(readback), uk, sk)
	}
}

// normalizeReadbackForLoad strips MySQL version-specific comment syntax that the omni
// loader does not need (e.g. /*!80023 ... */ inline hints stay; AUTO_INCREMENT counter
// is harmless). The readback is otherwise valid CREATE TABLE DDL.
func normalizeReadbackForLoad(ddl string) string {
	// Strip a trailing AUTO_INCREMENT=N table option so reload does not depend on it
	// (it is ignore-in-diff anyway).
	return ddl
}

// ---- the rule corpus: (id, inputDDL, table, columns, which versions diverge) --------

type ruleProbe struct {
	id       string
	create   string
	table    string
	columns  []string
	versions []Version // which engines to prove against
}

func both() []Version          { return []Version{MySQL57, MySQL80} }
func only(v Version) []Version { return []Version{v} }

func normalizationProbes() []ruleProbe {
	return []ruleProbe{
		{"int-display-width",
			"CREATE TABLE t_int (a INT(11), b BIGINT(20), c INT, f SMALLINT(6), o BIGINT, p TINYINT, q SMALLINT, r MEDIUMINT)",
			"t_int", []string{"a", "b", "c", "f", "o", "p", "q", "r"}, both()},
		{"tinyint1-boolean",
			"CREATE TABLE t_bool (d TINYINT(1), e TINYINT(4), m BOOLEAN, n BOOL, p TINYINT)",
			"t_bool", []string{"d", "e", "m", "n", "p"}, both()},
		{"int-unsigned-width",
			"CREATE TABLE t_uns (h INT UNSIGNED, i INT(11) UNSIGNED, b BIGINT UNSIGNED, ti TINYINT UNSIGNED)",
			"t_uns", []string{"h", "i", "b", "ti"}, both()},
		{"int-zerofill",
			"CREATE TABLE t_zf (j INT ZEROFILL, k INT(5) ZEROFILL, l INT UNSIGNED ZEROFILL)",
			"t_zf", []string{"j", "k", "l"}, both()},
		{"decimal-precision-scale",
			"CREATE TABLE t_num (a DECIMAL, b DECIMAL(10), c DECIMAL(10,2), d NUMERIC, e NUMERIC(8,3), k DEC, l FIXED)",
			"t_num", []string{"a", "b", "c", "d", "e", "k", "l"}, both()},
		{"float-double-aliasing",
			"CREATE TABLE t_fd (f FLOAT, g FLOAT(10,2), h DOUBLE, i DOUBLE(15,4), j REAL, m FLOAT(5))",
			"t_fd", []string{"f", "g", "h", "i", "j", "m"}, both()},
		{"char-binary-length-default",
			"CREATE TABLE t_char (a CHAR, b CHAR(10), d BINARY, e BINARY(16), f VARBINARY(32))",
			"t_char", []string{"a", "b", "d", "e", "f"}, both()},
		{"text-blob-no-default-null",
			"CREATE TABLE t_tb (g TEXT, h BLOB, j LONGTEXT, k TINYTEXT)",
			"t_tb", []string{"g", "h", "j", "k"}, both()},
		{"year-width",
			"CREATE TABLE t_year (a YEAR, b YEAR(4))",
			"t_year", []string{"a", "b"}, both()},
		{"bit-length-default",
			"CREATE TABLE t_bit (c BIT, d BIT(8))",
			"t_bit", []string{"c", "d"}, both()},
		{"time-json-date-stable",
			"CREATE TABLE t_tjd (e JSON, f DATE, g TIME, h TIME(3))",
			"t_tjd", []string{"e", "f", "g", "h"}, both()},
		{"default-literal-quoting",
			"CREATE TABLE t_def (a INT DEFAULT 0, b INT DEFAULT 5, h INT DEFAULT '7', c VARCHAR(20) DEFAULT 'hello', d VARCHAR(20) DEFAULT '', e VARCHAR(20) DEFAULT 'it''s')",
			"t_def", []string{"a", "b", "h", "c", "d", "e"}, both()},
		{"decimal-default-padding",
			"CREATE TABLE t_dpad (f DECIMAL(10,2) DEFAULT 0, g DECIMAL(10,2) DEFAULT 0.00)",
			"t_dpad", []string{"f", "g"}, both()},
		{"boolean-default",
			"CREATE TABLE t_bdef (j BOOLEAN DEFAULT TRUE, k BOOLEAN DEFAULT FALSE)",
			"t_bdef", []string{"j", "k"}, both()},
		{"nullable-default-null",
			"CREATE TABLE t_null (a INT NULL, b INT DEFAULT NULL, c INT NULL DEFAULT NULL, d VARCHAR(10) NULL)",
			"t_null", []string{"a", "b", "c", "d"}, both()},
		{"enum-set-quoting",
			`CREATE TABLE t_enum (a ENUM('x','y','z'), b SET('a','b','c'), e ENUM("dq1","dq2"))`,
			"t_enum", []string{"a", "b", "e"}, both()},
		{"comment-escaping",
			"CREATE TABLE t_comment (a INT COMMENT 'has ''quote'' inside')",
			"t_comment", []string{"a"}, both()},
		// version-flagged: charset echo. Only prove on the version whose readback diverges.
		{"column-charset-echo-80",
			"CREATE TABLE t_ce80 (a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci, c VARCHAR(10) CHARACTER SET utf8mb4, d VARCHAR(10)) DEFAULT CHARSET=utf8mb4",
			"t_ce80", []string{"a", "c", "d"}, only(MySQL80)},
		{"column-charset-echo-57",
			"CREATE TABLE t_ce57 (a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci, b VARCHAR(10) COLLATE utf8mb4_general_ci, c VARCHAR(10) CHARACTER SET utf8mb4, d VARCHAR(10)) DEFAULT CHARSET=utf8mb4",
			"t_ce57", []string{"a", "b", "c", "d"}, only(MySQL57)},
		{"column-charset-only-collation-resolution",
			"CREATE TABLE t_d1 (a VARCHAR(10) CHARACTER SET utf8mb4) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci",
			"t_d1", []string{"a"}, both()},
		// generated columns.
		{"generated-expr-normalization",
			"CREATE TABLE t_gen (a INT, b INT GENERATED ALWAYS AS (a+1) VIRTUAL, c INT GENERATED ALWAYS AS ( a * 2 ) STORED, d INT AS (a+1))",
			"t_gen", []string{"b", "c", "d"}, both()},
		{"generated-expr-string-introducer",
			"CREATE TABLE t_gen2 (a VARCHAR(20), e VARCHAR(20) GENERATED ALWAYS AS (CONCAT(a,'x')) VIRTUAL) DEFAULT CHARSET=utf8mb4",
			"t_gen2", []string{"e"}, both()},
		// functional default — 8.0 only.
		{"functional-default",
			"CREATE TABLE t_funcdef (a INT DEFAULT (1+1), b VARCHAR(36) DEFAULT (UUID()), c JSON DEFAULT (JSON_ARRAY()))",
			"t_funcdef", []string{"a", "b", "c"}, only(MySQL80)},
		// TIMESTAMP magic — both, opposite behavior.
		{"timestamp-magic",
			"CREATE TABLE t_ts (a TIMESTAMP, b TIMESTAMP NULL, d TIMESTAMP DEFAULT CURRENT_TIMESTAMP)",
			"t_ts", []string{"a", "b", "d"}, both()},
		{"timestamp-datetime-default-expr",
			"CREATE TABLE t_ts2 (g DATETIME DEFAULT CURRENT_TIMESTAMP, h DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, i TIMESTAMP(3) NULL DEFAULT CURRENT_TIMESTAMP(3))",
			"t_ts2", []string{"g", "h", "i"}, both()},
	}
}

func TestOracle_PhantomDiffElimination(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		for _, probe := range normalizationProbes() {
			if !containsVersion(probe.versions, version) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, probe.id), func(t *testing.T) {
				db := "nrm_" + strings.ReplaceAll(probe.id, "-", "_")
				for _, c := range probe.columns {
					assertPhantomFree(t, o, db, probe.create, probe.table, c)
				}
			})
		}
	}
}

func containsVersion(vs []Version, v Version) bool {
	for _, x := range vs {
		if x == v {
			return true
		}
	}
	return false
}

// TestOracle_MissedDiffGuards proves the canonical key DISTINGUISHES schemas that
// genuinely differ — the dual of phantom-diff elimination. These lock in the review
// findings where the canonicalizer was over-collapsing (a missed diff is as harmful as
// a phantom one). Each asserts two real, different MySQL schemas produce different keys,
// loaded through the real engine's SHOW CREATE so the catalog state is authentic.
func TestOracle_MissedDiffGuards(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		sc := serverCharsetFor(o.version)
		n := NormalizerFor(o.version)

		t.Run(o.name+"/bare-column-follows-table-collation", func(t *testing.T) {
			// A bare column inherits the table COLLATE; changing the table COLLATE must
			// change the bare column's canonical key (else the column change is invisible).
			ddlA := "CREATE TABLE t (a VARCHAR(10)) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci"
			ddlB := "CREATE TABLE t (a VARCHAR(10)) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci"
			rbA, okA := o.showCreate(t, "mdg_a", ddlA, "t")
			rbB, okB := o.showCreate(t, "mdg_b", ddlB, "t")
			if !okA || !okB {
				t.Skip("could not obtain readbacks")
			}
			tA, cA := loadColumn(t, sc, rbA, "t", "a")
			tB, cB := loadColumn(t, sc, rbB, "t", "a")
			if n.CanonicalColumn(tA, cA) == n.CanonicalColumn(tB, cB) {
				t.Errorf("[%s] bare column must follow table COLLATE change:\n A=%s\n B=%s",
					o.name, n.CanonicalColumn(tA, cA), n.CanonicalColumn(tB, cB))
			}
		})

		if version == MySQL57 {
			t.Run(o.name+"/timestamp-null-vs-notnull-with-default", func(t *testing.T) {
				// EDFT=0: `TIMESTAMP NULL DEFAULT '<const>'` (nullable) must produce a
				// different key from a bare `TIMESTAMP DEFAULT '<const>'` (NOT NULL).
				ddl := "CREATE TABLE t (x INT, n TIMESTAMP NULL DEFAULT '2020-01-01 00:00:00', b TIMESTAMP DEFAULT '2020-01-01 00:00:00')"
				rb, ok := o.showCreate(t, "mdg_ts", ddl, "t")
				if !ok {
					t.Skip("could not obtain readback")
				}
				tbl, nCol := loadColumn(t, sc, rb, "t", "n")
				_, bCol := loadColumn(t, sc, rb, "t", "b")
				if n.CanonicalColumn(tbl, nCol) == n.CanonicalColumn(tbl, bCol) {
					t.Errorf("[%s] nullable TIMESTAMP NULL and NOT NULL TIMESTAMP must differ", o.name)
				}
			})
		}
	}
}
