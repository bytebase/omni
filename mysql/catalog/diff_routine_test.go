package catalog

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// The oracle proof for the routine (function + procedure) differ (correctness-protocol.md
// gate 1: idempotence), against the LIVE MySQL engines (5.7 :13307, 8.0 :13306). Two
// properties are proven mechanically for every routine FORM:
//
//  1. Idempotence (the spine): a routine in its real stored form (the engine's own SHOW
//     CREATE FUNCTION/PROCEDURE readback, reloaded into a catalog) self-diffs empty —
//     Diff(c, c).IsEmpty().
//  2. Canonicalization-driven empties: a routine written in a USER form (explicit default
//     characteristics that SHOW CREATE drops, varied whitespace) diffed against the engine's
//     STORED form of the same routine must be EMPTY. This proves the differ folds defaults
//     (CONTAINS SQL, SQL SECURITY DEFINER, NOT DETERMINISTIC) and compares the body opaquely —
//     if it compared surface tokens these would phantom-diff forever.
//
// Both sides are loaded under the SAME database name (routineProbeDB) so identity matches the
// single-database release path; the harness reuses connectOracle / serverCharsetFor / both()
// from the diff + normalize oracle tests and skips cleanly when the engines are unreachable.

const routineProbeDB = "rtn_idem"

// routineProbe is one idempotence/canonicalization case: a routine DDL whose engine-stored form
// may differ from the user form but must canonicalize equal. kind is "FUNCTION"/"PROCEDURE".
type routineProbe struct {
	id       string
	kind     string
	name     string
	create   string // bare CREATE FUNCTION/PROCEDURE (no DEFINER, no leading CREATE DATABASE)
	versions []Version
}

// routineIdempotenceProbes enumerates representative routine FORMS that must round-trip empty.
// Functions on 8.0 must declare DETERMINISTIC / NO SQL / READS SQL DATA (the box has
// log_bin_trust_function_creators off, so a non-deterministic CONTAINS-SQL function is rejected
// at CREATE — error 1418); each function probe satisfies that. Procedures have no such rule.
func routineIdempotenceProbes() []routineProbe {
	return []routineProbe{
		// ---- functions ----
		{"fn-simple", "FUNCTION", "f", "CREATE FUNCTION f(a INT, b INT) RETURNS INT DETERMINISTIC RETURN a + b", both()},
		{"fn-string-return", "FUNCTION", "f", "CREATE FUNCTION f(s VARCHAR(50)) RETURNS VARCHAR(100) DETERMINISTIC RETURN CONCAT(s, s)", both()},
		{"fn-all-characteristics", "FUNCTION", "f", "CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC READS SQL DATA SQL SECURITY INVOKER COMMENT 'hi' RETURN a", both()},
		{"fn-no-sql", "FUNCTION", "f", "CREATE FUNCTION f() RETURNS INT NO SQL RETURN 42", both()},
		{"fn-modifies", "FUNCTION", "f", "CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC MODIFIES SQL DATA RETURN a", both()},
		// Explicit-default characteristics the readback DROPS — must still canonicalize empty.
		{"fn-explicit-defaults", "FUNCTION", "f", "CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC CONTAINS SQL SQL SECURITY DEFINER RETURN a", both()},
		// LANGUAGE SQL is never echoed distinctly — must not phantom-diff.
		{"fn-language-sql", "FUNCTION", "f", "CREATE FUNCTION f() RETURNS INT LANGUAGE SQL DETERMINISTIC RETURN 7", both()},
		// Multi-statement BEGIN...END body with internal whitespace (stored verbatim).
		{"fn-begin-body", "FUNCTION", "f", "CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC\nBEGIN\n  DECLARE r INT;\n  SET r = a * 2;\n  RETURN r;\nEND", both()},
		// Comment with an embedded single quote (escaping round-trip).
		{"fn-comment-quote", "FUNCTION", "f", "CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC COMMENT 'it''s here' RETURN a", both()},
		// Comment with an embedded backslash: the engine echoes it backslash-escaped and the loader
		// decodes it back consistently, so the no-op diff stays empty (guards against an
		// escape-asymmetry idempotence break).
		{"fn-comment-backslash", "FUNCTION", "f", "CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC COMMENT 'a\\\\b' RETURN a", both()},
		// DECIMAL RETURNS: precision/scale padding via the type canonicalizer.
		{"fn-decimal-return", "FUNCTION", "f", "CREATE FUNCTION f(a DECIMAL(10,2)) RETURNS DECIMAL(10,2) DETERMINISTIC RETURN a", both()},
		// ENUM RETURNS + ENUM/SET params: the member list must survive the load (was dropped before
		// the formatDataType fix) so the type round-trips.
		{"fn-enum-return", "FUNCTION", "f", "CREATE FUNCTION f() RETURNS ENUM('a','b','c') DETERMINISTIC RETURN 'a'", both()},
		{"proc-enum-set-params", "PROCEDURE", "p", "CREATE PROCEDURE p(IN s SET('x','y'), IN e ENUM('p','q')) BEGIN SET @v = s; END", both()},
		// Reserved-word parameter name (must be backtick-quoted in the rendered CREATE).
		{"proc-reserved-param", "PROCEDURE", "p", "CREATE PROCEDURE p(IN `select` INT, IN `from` INT) BEGIN SET @v = `select`; END", both()},

		// ---- procedures ----
		{"proc-in-out-inout", "PROCEDURE", "p", "CREATE PROCEDURE p(IN x INT, OUT y INT, INOUT z INT) BEGIN SET y = x + z; SET z = z + 1; END", both()},
		{"proc-no-params", "PROCEDURE", "p", "CREATE PROCEDURE p() BEGIN DECLARE n INT; SET n = 1; END", both()},
		{"proc-modifies", "PROCEDURE", "p", "CREATE PROCEDURE p() MODIFIES SQL DATA BEGIN SET @v = 1; END", both()},
		{"proc-default-in-direction", "PROCEDURE", "p", "CREATE PROCEDURE p(x INT, y INT) BEGIN SET @v = x + y; END", both()},
		{"proc-all-characteristics", "PROCEDURE", "p", "CREATE PROCEDURE p(IN a INT) DETERMINISTIC READS SQL DATA SQL SECURITY INVOKER COMMENT 'pc' BEGIN SET @v = a; END", both()},
		{"proc-explicit-defaults", "PROCEDURE", "p", "CREATE PROCEDURE p(IN a INT) NOT DETERMINISTIC CONTAINS SQL SQL SECURITY DEFINER BEGIN SET @v = a; END", both()},
		// Adjacent string-literal concatenation in the body ('...' '...' across a line
		// break — the stock-8.0 sys.ps_trace_thread shape). Bodies are stored VERBATIM,
		// so the engine readback still carries the adjacency and the reload re-parses
		// it; before adjacency support that failed with "unexpected token".
		{"proc-adjacent-literals", "PROCEDURE", "p", "CREATE PROCEDURE p()\nBEGIN\n    SELECT CONCAT('tmp disk tables: ', 3, '\\n'\n                  'select scan: ', 4, '\\n');\nEND", both()},
	}
}

// loadRoutineReadback applies a CREATE FUNCTION/PROCEDURE in a throwaway database, reads back its
// SHOW CREATE form (the engine's canonical stored form), and reloads it into a catalog under
// dbName. It returns the catalog + the routine's authentic stored DDL.
func (o *oracleConn) loadRoutineReadback(t *testing.T, dbName, createSQL, kind, name string) (*Catalog, string, bool) {
	t.Helper()
	ctx := context.Background()
	sc := serverCharsetFor(o.version)
	stmts := []string{
		"DROP DATABASE IF EXISTS " + dbName,
		fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s", dbName, sc),
		"USE " + dbName,
		createSQL,
	}
	for _, s := range stmts {
		if _, err := o.db.ExecContext(ctx, s); err != nil {
			t.Logf("[%s] exec failed (may be expected): %q: %v", o.name, s, err)
			return nil, "", false
		}
	}
	rb, ok := o.showCreateRoutine(t, dbName, kind, name)
	if !ok {
		return nil, "", false
	}
	wrapped := fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s;\nUSE %s;\n%s", dbName, sc, dbName, rb)
	cat, err := LoadSQL(wrapped)
	if err != nil {
		t.Fatalf("[%s] reload of routine readback failed: %v\n%s", o.name, err, wrapped)
	}
	return cat, rb, true
}

// showCreateRoutine runs SHOW CREATE FUNCTION/PROCEDURE and returns the create DDL column. The
// statement returns six columns (Name, sql_mode, Create <kind>, character_set_client,
// collation_connection, Database Collation); only the third (the DDL) is needed.
func (o *oracleConn) showCreateRoutine(t *testing.T, dbName, kind, name string) (string, bool) {
	t.Helper()
	var c1, c2, c3, c4, c5, c6 sql.NullString
	row := o.db.QueryRowContext(context.Background(), fmt.Sprintf("SHOW CREATE %s %s.%s", kind, dbName, name))
	if err := row.Scan(&c1, &c2, &c3, &c4, &c5, &c6); err != nil {
		t.Logf("[%s] SHOW CREATE %s %s.%s failed: %v", o.name, kind, dbName, name, err)
		return "", false
	}
	return c3.String, true
}

// TestOracle_RoutineDiffIdempotence proves gate 1 & 2 for every routine form: the user DDL vs its
// engine readback diffs EMPTY, and the stored form self-diffs empty, on every supported version.
func TestOracle_RoutineDiffIdempotence(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		n := NormalizerFor(version)
		sc := serverCharsetFor(version)
		for _, probe := range routineIdempotenceProbes() {
			if !containsVersion(probe.versions, version) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, probe.id), func(t *testing.T) {
				db := routineProbeDB + "_" + strings.ReplaceAll(probe.id, "-", "_")
				storedCat, rb, ok := o.loadRoutineReadback(t, db, probe.create, probe.kind, probe.name)
				if !ok {
					t.Skipf("[%s] could not obtain readback for %s", o.name, probe.id)
				}

				// Idempotence spine: the stored form self-diffs empty.
				if d := DiffWithNormalizer(storedCat, storedCat, n); !d.IsEmpty() {
					t.Errorf("[%s] IDEMPOTENCE: self-diff of stored form not empty for %s:\n  stored: %s\n  diff: %s",
						o.name, probe.id, strings.TrimSpace(rb), describeRoutineDiff(d))
				}

				// Canonicalization: the USER form vs the engine STORED form must collapse to empty.
				userWrapped := fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s;\nUSE %s;\n%s", db, sc, db, probe.create)
				userCat, err := LoadSQL(userWrapped)
				if err != nil {
					t.Fatalf("[%s] user-form load failed for %s: %v", o.name, probe.id, err)
				}
				if d := DiffWithNormalizer(userCat, userCat, n); !d.IsEmpty() {
					t.Errorf("[%s] user-form self-diff not empty for %s: %s", o.name, probe.id, describeRoutineDiff(d))
				}
				if d := DiffWithNormalizer(userCat, storedCat, n); !d.IsEmpty() {
					t.Errorf("[%s] CANONICALIZATION: user vs stored not empty for %s:\n  user:   %s\n  stored: %s\n  diff: %s",
						o.name, probe.id, strings.TrimSpace(probe.create), strings.TrimSpace(rb), describeRoutineDiff(d))
				}
				// Symmetry: stored vs user must also be empty (Drop/Add are direction-sensitive).
				if d := DiffWithNormalizer(storedCat, userCat, n); !d.IsEmpty() {
					t.Errorf("[%s] CANONICALIZATION (reverse): stored vs user not empty for %s: %s",
						o.name, probe.id, describeRoutineDiff(d))
				}
			})
		}
	}
}

// TestOracle_SysProcedureAdjacentLiterals proves the real-world motivation for
// adjacent string-literal support against the engine's OWN schema: the stock
// sys.ps_trace_thread body (shipped in every 5.7/8.0 server) contains an
// adjacent-literal run ('tmp disk tables: ...\n' followed by 'select scan: ...'
// with no comma), so its SHOW CREATE readback must parse, load through LoadSDL,
// and self-diff empty. Before adjacency support this failed with "unexpected
// token" — hard-blocking the declarative path for any dump containing sys.
func TestOracle_SysProcedureAdjacentLiterals(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		n := NormalizerFor(o.version)
		t.Run(o.name, func(t *testing.T) {
			rb, ok := o.showCreateRoutine(t, "sys", "PROCEDURE", "ps_trace_thread")
			if !ok {
				t.Skipf("[%s] sys.ps_trace_thread unavailable (stripped sys schema?)", o.name)
			}
			sdl := "CREATE DATABASE sysdump;\nUSE sysdump;\n" + rb + ";\n"
			cat, err := LoadSDL(sdl)
			if err != nil {
				t.Fatalf("[%s] LoadSDL of stock sys.ps_trace_thread failed: %v", o.name, err)
			}
			if d := DiffWithNormalizer(cat, cat, n); !d.IsEmpty() {
				t.Errorf("[%s] self-diff of stock sys.ps_trace_thread not empty: %s",
					o.name, describeRoutineDiff(d))
			}
		})
	}
}

// TestRoutineDiffActions is a hermetic (no-oracle) check that the differ reports the right
// DiffAction for in-memory routine catalogs: add, drop, modify (body), and no-op for an identical
// pair. It pins the structural behavior independent of a live engine.
func TestRoutineDiffActions(t *testing.T) {
	n := NormalizerFor(MySQL80)
	load := func(t *testing.T, ddl string) *Catalog {
		t.Helper()
		cat, err := LoadSQL("CREATE DATABASE d DEFAULT CHARSET=utf8mb4;\nUSE d;\n" + ddl)
		if err != nil {
			t.Fatalf("load %q: %v", ddl, err)
		}
		return cat
	}

	base := "CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a"

	t.Run("identical-empty", func(t *testing.T) {
		from := load(t, base)
		to := load(t, base)
		if d := DiffWithNormalizer(from, to, n); !d.IsEmpty() {
			t.Fatalf("identical routines diff not empty: %s", describeRoutineDiff(d))
		}
	})
	t.Run("add", func(t *testing.T) {
		from := load(t, "CREATE FUNCTION other() RETURNS INT DETERMINISTIC RETURN 1")
		to := load(t, "CREATE FUNCTION other() RETURNS INT DETERMINISTIC RETURN 1;\n"+base)
		d := DiffWithNormalizer(from, to, n)
		if len(d.Functions) != 1 || d.Functions[0].Action != DiffAdd || d.Functions[0].Name != "f" {
			t.Fatalf("want one ADD of f, got %s", describeRoutineDiff(d))
		}
	})
	t.Run("drop", func(t *testing.T) {
		from := load(t, base)
		to := load(t, "CREATE FUNCTION g() RETURNS INT DETERMINISTIC RETURN 1")
		d := DiffWithNormalizer(from, to, n)
		// f dropped, g added.
		var dropF bool
		for _, e := range d.Functions {
			if e.Action == DiffDrop && e.Name == "f" {
				dropF = true
			}
		}
		if !dropF {
			t.Fatalf("want DROP of f, got %s", describeRoutineDiff(d))
		}
	})
	t.Run("modify-body", func(t *testing.T) {
		from := load(t, "CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a + 1")
		to := load(t, "CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a + 2")
		d := DiffWithNormalizer(from, to, n)
		if len(d.Functions) != 1 || d.Functions[0].Action != DiffModify {
			t.Fatalf("want one MODIFY of f, got %s", describeRoutineDiff(d))
		}
	})
	t.Run("function-and-procedure-same-name-distinct", func(t *testing.T) {
		// A function `x` and a procedure `x` are different objects; adding the procedure must not
		// be seen as modifying the function.
		from := load(t, "CREATE FUNCTION x() RETURNS INT DETERMINISTIC RETURN 1")
		to := load(t, "CREATE FUNCTION x() RETURNS INT DETERMINISTIC RETURN 1;\nCREATE PROCEDURE x() BEGIN SET @v = 1; END")
		d := DiffWithNormalizer(from, to, n)
		if len(d.Functions) != 0 {
			t.Fatalf("function x should be unchanged, got %s", describeRoutineDiff(d))
		}
		if len(d.Procedures) != 1 || d.Procedures[0].Action != DiffAdd {
			t.Fatalf("want one ADD procedure x, got %s", describeRoutineDiff(d))
		}
	})
	t.Run("modify-comment-only-alter-suffices", func(t *testing.T) {
		from := load(t, "CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC COMMENT 'old' RETURN a")
		to := load(t, "CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC COMMENT 'new' RETURN a")
		d := DiffWithNormalizer(from, to, n)
		if len(d.Functions) != 1 || d.Functions[0].Action != DiffModify {
			t.Fatalf("want MODIFY, got %s", describeRoutineDiff(d))
		}
		if !routineAlterSuffices(d.Functions[0].From, d.Functions[0].To, n) {
			t.Fatalf("comment-only change should be ALTER-able")
		}
		// And the generated plan must be a single ALTER, not DROP+CREATE.
		plan := GenerateMigrationWithNormalizer(from, to, d, n)
		if len(plan.Ops) != 1 || !strings.HasPrefix(plan.Ops[0].SQL, "ALTER FUNCTION") {
			t.Fatalf("want single ALTER FUNCTION, got plan:\n%s", plan.SQL())
		}
	})
	t.Run("enum-return-renders-members-and-self-diffs-empty", func(t *testing.T) {
		// The ENUM member list must survive load and render in the generated CREATE (else invalid
		// DDL). Self-diff empty proves the loaded form is stable; the rendered CREATE must carry
		// the members.
		from := load(t, "CREATE FUNCTION f() RETURNS ENUM('a','b','c') DETERMINISTIC RETURN 'a'")
		if d := DiffWithNormalizer(from, from, n); !d.IsEmpty() {
			t.Fatalf("enum-return self-diff not empty: %s", describeRoutineDiff(d))
		}
		plan := GenerateMigrationWithNormalizer(New(), from, DiffWithNormalizer(New(), from, n), n)
		if !strings.Contains(plan.SQL(), "enum('a','b','c')") {
			t.Fatalf("rendered CREATE dropped enum members:\n%s", plan.SQL())
		}
	})
	t.Run("reserved-param-name-is-quoted", func(t *testing.T) {
		from := load(t, "CREATE PROCEDURE p(IN `select` INT) BEGIN SET @v = `select`; END")
		plan := GenerateMigrationWithNormalizer(New(), from, DiffWithNormalizer(New(), from, n), n)
		if !strings.Contains(plan.SQL(), "`select`") {
			t.Fatalf("reserved param name not backtick-quoted in rendered CREATE:\n%s", plan.SQL())
		}
	})
	t.Run("procedure-implicit-in-equals-explicit-in", func(t *testing.T) {
		// `p(x INT)` and `p(IN x INT)` are the SAME procedure (IN is the default direction); the
		// differ must treat them as equal (no spurious DROP+CREATE on a cosmetic direction keyword).
		from := load(t, "CREATE PROCEDURE p(x INT) BEGIN SET @v = x; END")
		to := load(t, "CREATE PROCEDURE p(IN x INT) BEGIN SET @v = x; END")
		if d := DiffWithNormalizer(from, to, n); !d.IsEmpty() {
			t.Fatalf("implicit-IN vs explicit-IN should be equal, got %s", describeRoutineDiff(d))
		}
	})
	t.Run("param-enum-member-case-change-detected", func(t *testing.T) {
		// A case-only ENUM member change in a parameter is a real domain change and MUST be
		// detected (the differ must not fold member case). Guards the comparison-side mirror of the
		// render-side member-case fix.
		from := load(t, "CREATE PROCEDURE p(IN e ENUM('a','b')) BEGIN SET @v = e; END")
		to := load(t, "CREATE PROCEDURE p(IN e ENUM('A','b')) BEGIN SET @v = e; END")
		d := DiffWithNormalizer(from, to, n)
		if len(d.Procedures) != 1 || d.Procedures[0].Action != DiffModify {
			t.Fatalf("enum-member-case param change must be a MODIFY, got %s", describeRoutineDiff(d))
		}
		// A param-type change is not ALTER-able → DROP+CREATE.
		if routineAlterSuffices(d.Procedures[0].From, d.Procedures[0].To, n) {
			t.Fatalf("param type change must NOT be ALTER-able")
		}
	})
	t.Run("modify-deterministic-requires-dropcreate", func(t *testing.T) {
		from := load(t, "CREATE FUNCTION f(a INT) RETURNS INT NOT DETERMINISTIC READS SQL DATA RETURN a")
		to := load(t, "CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC READS SQL DATA RETURN a")
		d := DiffWithNormalizer(from, to, n)
		if len(d.Functions) != 1 {
			t.Fatalf("want one function diff, got %s", describeRoutineDiff(d))
		}
		if routineAlterSuffices(d.Functions[0].From, d.Functions[0].To, n) {
			t.Fatalf("DETERMINISTIC change must NOT be ALTER-able")
		}
		plan := GenerateMigrationWithNormalizer(from, to, d, n)
		if len(plan.Ops) != 2 {
			t.Fatalf("want DROP+CREATE (2 ops), got:\n%s", plan.SQL())
		}
		if plan.Ops[0].Type != OpDropFunction || plan.Ops[1].Type != OpCreateFunction {
			t.Fatalf("want DROP then CREATE, got %v then %v", plan.Ops[0].Type, plan.Ops[1].Type)
		}
	})
}

// describeRoutineDiff renders a compact human description of the routine slices of a SchemaDiff.
func describeRoutineDiff(d *SchemaDiff) string {
	var b strings.Builder
	for _, e := range d.Functions {
		fmt.Fprintf(&b, "[fn %s.%s %s]", e.Database, e.Name, e.Action)
	}
	for _, e := range d.Procedures {
		fmt.Fprintf(&b, "[proc %s.%s %s]", e.Database, e.Name, e.Action)
	}
	if b.Len() == 0 {
		// Fall back to the table description so a stray table diff is visible too.
		return describeDiff(d)
	}
	return b.String()
}
