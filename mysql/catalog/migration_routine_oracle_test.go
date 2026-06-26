package catalog

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// The oracle apply-correctness proof for the routine generator (correctness-protocol.md gate 2),
// against the LIVE MySQL engines (5.7 :13307, 8.0 :13306). For representative (from, to) pairs
// covering every routine variant — function vs procedure, CREATE / DROP / characteristic-ALTER /
// body-change(=DROP+CREATE) / param-change / RETURNS-change / DETERMINISTIC-change — the generated
// DDL applied to a real `from` database must yield a routine whose canonical form equals `to`.
//
// Both from/to catalogs are loaded from the ENGINE'S OWN SHOW CREATE readbacks under the SAME
// database name (so routine identity matches), then the plan is applied to a from-state database
// and the result read back and compared. Idempotence (gate 1) is covered by
// TestOracle_RoutineDiffIdempotence; this file proves the EMITTED DDL is correct, and additionally
// asserts the chosen PATH (single ALTER vs DROP+CREATE) for the characteristic-only cases.
//
// The harness reuses connectOracle / serverCharsetFor / both() and skips cleanly when the engines
// are unreachable.

const routineApplyDB = "rtn_apply"

// routineMigrationProbe is one apply-correctness case: transform a routine from `fromDDL` to
// `toDDL`. An empty fromDDL means the routine does not exist yet (a pure CREATE); an empty toDDL
// means it is dropped. wantAlter asserts the plan path: true → a single ALTER, false → CREATE-only
// or DROP-only or DROP+CREATE (the default; only checked when both sides are present).
type routineMigrationProbe struct {
	id        string
	kind      string // "FUNCTION" | "PROCEDURE"
	name      string
	fromDDL   string
	toDDL     string
	wantAlter bool // only meaningful for a modify (both DDLs present)
	versions  []Version
}

// routineMigrationProbes enumerates the routine ADD/DROP/MODIFY FORMS the generator covers.
// Functions declare a binlog-safe characteristic (DETERMINISTIC / NO SQL / READS SQL DATA) so they
// create on the 8.0 box (log_bin_trust_function_creators off → error 1418 otherwise).
func routineMigrationProbes() []routineMigrationProbe {
	return []routineMigrationProbe{
		// ---- CREATE (from empty) ----
		{"create-fn-simple", "FUNCTION", "f", "", "CREATE FUNCTION f(a INT, b INT) RETURNS INT DETERMINISTIC RETURN a + b", false, both()},
		{"create-fn-string", "FUNCTION", "f", "", "CREATE FUNCTION f(s VARCHAR(50)) RETURNS VARCHAR(100) DETERMINISTIC RETURN CONCAT(s, s)", false, both()},
		{"create-fn-all-chars", "FUNCTION", "f", "", "CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC READS SQL DATA SQL SECURITY INVOKER COMMENT 'c' RETURN a", false, both()},
		{"create-fn-begin", "FUNCTION", "f", "", "CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC\nBEGIN\n  DECLARE r INT;\n  SET r = a * 2;\n  RETURN r;\nEND", false, both()},
		{"create-proc-inout", "PROCEDURE", "p", "", "CREATE PROCEDURE p(IN x INT, OUT y INT, INOUT z INT) BEGIN SET y = x + z; SET z = z + 1; END", false, both()},
		{"create-proc-modifies", "PROCEDURE", "p", "", "CREATE PROCEDURE p() MODIFIES SQL DATA BEGIN SET @v = 1; END", false, both()},
		{"create-proc-comment", "PROCEDURE", "p", "", "CREATE PROCEDURE p(IN a INT) SQL SECURITY INVOKER COMMENT 'pc' BEGIN SET @v = a; END", false, both()},
		// Type-rendering surfaces: DECIMAL precision, ENUM/SET member lists, reserved-word param.
		{"create-fn-decimal", "FUNCTION", "f", "", "CREATE FUNCTION f(a DECIMAL(10,2)) RETURNS DECIMAL(12,4) DETERMINISTIC RETURN a", false, both()},
		{"create-fn-enum-return", "FUNCTION", "f", "", "CREATE FUNCTION f() RETURNS ENUM('a','b','c') DETERMINISTIC RETURN 'a'", false, both()},
		{"create-proc-enum-set-params", "PROCEDURE", "p", "", "CREATE PROCEDURE p(IN s SET('x','y'), IN e ENUM('p','q')) BEGIN SET @v = s; END", false, both()},
		{"create-proc-reserved-param", "PROCEDURE", "p", "", "CREATE PROCEDURE p(IN `select` INT, IN `from` INT) BEGIN SET @v = `select`; END", false, both()},

		// ---- DROP (to empty) ----
		{"drop-fn", "FUNCTION", "f", "CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a", "", false, both()},
		{"drop-proc", "PROCEDURE", "p", "CREATE PROCEDURE p() BEGIN SET @v = 1; END", "", false, both()},

		// ---- MODIFY via ALTER (only SQL SECURITY / COMMENT / DATA ACCESS change) ----
		{"alter-fn-comment", "FUNCTION", "f",
			"CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC COMMENT 'old' RETURN a",
			"CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC COMMENT 'new' RETURN a", true, both()},
		{"alter-fn-security", "FUNCTION", "f",
			"CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC SQL SECURITY DEFINER RETURN a",
			"CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC SQL SECURITY INVOKER RETURN a", true, both()},
		{"alter-fn-dataaccess", "FUNCTION", "f",
			"CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a",
			"CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC READS SQL DATA RETURN a", true, both()},
		{"alter-fn-add-comment", "FUNCTION", "f",
			"CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a",
			"CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC COMMENT 'added' RETURN a", true, both()},
		{"alter-proc-security", "PROCEDURE", "p",
			"CREATE PROCEDURE p(IN a INT) SQL SECURITY DEFINER BEGIN SET @v = a; END",
			"CREATE PROCEDURE p(IN a INT) SQL SECURITY INVOKER BEGIN SET @v = a; END", true, both()},
		{"alter-proc-dataaccess", "PROCEDURE", "p",
			"CREATE PROCEDURE p() CONTAINS SQL BEGIN SET @v = 1; END",
			"CREATE PROCEDURE p() MODIFIES SQL DATA BEGIN SET @v = 1; END", true, both()},

		// ---- MODIFY via DROP+CREATE (body / params / RETURNS / DETERMINISTIC change) ----
		{"body-change-fn", "FUNCTION", "f",
			"CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a + 1",
			"CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a + 2", false, both()},
		{"param-change-fn", "FUNCTION", "f",
			"CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a",
			"CREATE FUNCTION f(a INT, b INT) RETURNS INT DETERMINISTIC RETURN a + b", false, both()},
		{"returns-change-fn", "FUNCTION", "f",
			"CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a",
			"CREATE FUNCTION f(a INT) RETURNS BIGINT DETERMINISTIC RETURN a", false, both()},
		{"determ-change-fn", "FUNCTION", "f",
			"CREATE FUNCTION f(a INT) RETURNS INT NOT DETERMINISTIC READS SQL DATA RETURN a",
			"CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC READS SQL DATA RETURN a", false, both()},
		{"body-change-proc", "PROCEDURE", "p",
			"CREATE PROCEDURE p() BEGIN SET @v = 1; END",
			"CREATE PROCEDURE p() BEGIN SET @v = 2; END", false, both()},
		{"param-change-proc", "PROCEDURE", "p",
			"CREATE PROCEDURE p(IN x INT) BEGIN SET @v = x; END",
			"CREATE PROCEDURE p(IN x INT, OUT y INT) BEGIN SET y = x; END", false, both()},
		// ENUM member-case change in a parameter: a real domain change → DROP+CREATE, and the
		// re-created routine must carry the NEW member case (proves member case survives render).
		{"param-enum-case-change-proc", "PROCEDURE", "p",
			"CREATE PROCEDURE p(IN e ENUM('a','b')) BEGIN SET @v = e; END",
			"CREATE PROCEDURE p(IN e ENUM('A','b')) BEGIN SET @v = e; END", false, both()},
		// Combined: a body change AND a comment change at once → still DROP+CREATE (body forces it),
		// and the re-created routine carries the new comment.
		{"body-and-comment-fn", "FUNCTION", "f",
			"CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC COMMENT 'old' RETURN a + 1",
			"CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC COMMENT 'new' RETURN a + 2", false, both()},
	}
}

// TestOracle_RoutineMigrationApplyCorrectness proves gate 2 for every routine probe: the generated
// DDL transforms a real `from` database into a `to`-equal one (compared via canonical readback),
// and the chosen path (ALTER vs DROP+CREATE) matches the probe's expectation.
func TestOracle_RoutineMigrationApplyCorrectness(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		n := NormalizerFor(version)
		for _, probe := range routineMigrationProbes() {
			if !containsVersion(probe.versions, version) {
				continue
			}
			t.Run(fmt.Sprintf("%s/%s", o.name, probe.id), func(t *testing.T) {
				assertRoutineApplyCorrect(t, o, n, probe)
			})
		}
	}
}

// assertRoutineApplyCorrect loads from/to catalogs from the engine's readbacks (under the SAME db
// name so identity matches), generates the plan, applies it to a from-state database, reads the
// result back, and asserts the result canonicalizes equal to `to`.
func assertRoutineApplyCorrect(t *testing.T, o *oracleConn, n *Normalizer, p routineMigrationProbe) {
	t.Helper()
	ctx := context.Background()
	sc := serverCharsetFor(o.version)
	db := routineApplyDB + "_" + strings.ReplaceAll(p.id, "-", "_")

	// Build the `from` and `to` catalogs from engine readbacks, BOTH under db name `db`.
	loadSide := func(ddl string) (*Catalog, bool) {
		if strings.TrimSpace(ddl) == "" {
			return New(), true
		}
		cat, _, ok := o.loadRoutineReadback(t, db, ddl, p.kind, p.name)
		return cat, ok
	}
	toCat, ok := loadSide(p.toDDL)
	if !ok {
		t.Skipf("[%s] could not obtain `to` readback for %s", o.name, p.id)
	}
	fromCat, ok := loadSide(p.fromDDL)
	if !ok {
		t.Skipf("[%s] could not obtain `from` readback for %s", o.name, p.id)
	}

	diff := DiffWithNormalizer(fromCat, toCat, n)
	plan := GenerateMigrationWithNormalizer(fromCat, toCat, diff, n)

	// Path assertion (modify cases only): ALTER vs DROP+CREATE.
	bothPresent := strings.TrimSpace(p.fromDDL) != "" && strings.TrimSpace(p.toDDL) != ""
	if bothPresent {
		alterOps, createOps, dropOps := classifyRoutineOps(plan.Ops)
		if p.wantAlter {
			if alterOps != 1 || createOps != 0 || dropOps != 0 {
				t.Errorf("[%s] %s: expected a single ALTER, got plan:\n%s", o.name, p.id, plan.SQL())
			}
		} else {
			if dropOps != 1 || createOps != 1 {
				t.Errorf("[%s] %s: expected DROP+CREATE, got plan:\n%s", o.name, p.id, plan.SQL())
			}
		}
	}

	// Apply the plan to a real database in `from` state, on ONE dedicated connection.
	conn, err := o.db.Conn(ctx)
	if err != nil {
		t.Fatalf("[%s] grab conn: %v", o.name, err)
	}
	defer func() { _ = conn.Close() }()

	setup := []string{
		"DROP DATABASE IF EXISTS " + db,
		fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s", db, sc),
		"USE " + db,
	}
	if strings.TrimSpace(p.fromDDL) != "" {
		setup = append(setup, p.fromDDL)
	}
	for _, s := range setup {
		if _, err := conn.ExecContext(ctx, s); err != nil {
			t.Skipf("[%s] could not set up `from` state for %s: %q: %v", o.name, p.id, s, err)
		}
	}

	for _, op := range plan.Ops {
		if _, err := conn.ExecContext(ctx, op.SQL); err != nil {
			t.Fatalf("[%s] APPLY FAILED for %s:\n  stmt: %s\n  err: %v\n  full plan:\n%s",
				o.name, p.id, op.SQL, err, plan.SQL())
		}
	}

	// Readback the result on the same connection.
	readback := func() (string, bool) {
		var c1, c2, c3, c4, c5, c6 sql.NullString
		row := conn.QueryRowContext(ctx, fmt.Sprintf("SHOW CREATE %s %s.%s", p.kind, db, p.name))
		if err := row.Scan(&c1, &c2, &c3, &c4, &c5, &c6); err != nil {
			return "", false
		}
		return c3.String, true
	}

	if strings.TrimSpace(p.toDDL) == "" {
		if _, ok := readback(); ok {
			t.Errorf("[%s] %s: routine %s still exists after DROP plan:\n%s", o.name, p.id, p.name, plan.SQL())
		}
		return
	}

	resultRB, ok := readback()
	if !ok {
		t.Fatalf("[%s] %s: result routine %s missing after apply:\n%s", o.name, p.id, p.name, plan.SQL())
	}
	resultWrapped := fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s;\nUSE %s;\n%s", db, sc, db, resultRB)
	resultCat, err := LoadSQL(resultWrapped)
	if err != nil {
		t.Fatalf("[%s] %s: result load failed: %v\n%s", o.name, p.id, err, resultWrapped)
	}

	if d := DiffWithNormalizer(resultCat, toCat, n); !d.IsEmpty() {
		t.Errorf("[%s] APPLY-CORRECTNESS FAILED for %s: result != to\n  plan:\n%s\n  result: %s\n  diff: %s",
			o.name, p.id, plan.SQL(), strings.TrimSpace(resultRB), describeRoutineDiff(d))
	}
	// Symmetry.
	if d := DiffWithNormalizer(toCat, resultCat, n); !d.IsEmpty() {
		t.Errorf("[%s] APPLY-CORRECTNESS (reverse) FAILED for %s: %s", o.name, p.id, describeRoutineDiff(d))
	}
}

// classifyRoutineOps counts the routine op kinds in a plan: ALTER (a CREATE-type op whose SQL
// begins with ALTER), CREATE (a CREATE-type op whose SQL begins with CREATE), and DROP.
func classifyRoutineOps(ops []MigrationOp) (alter, create, drop int) {
	for _, op := range ops {
		switch op.Type {
		case OpDropFunction, OpDropProcedure:
			drop++
		case OpCreateFunction, OpCreateProcedure:
			if strings.HasPrefix(strings.TrimSpace(op.SQL), "ALTER ") {
				alter++
			} else {
				create++
			}
		}
	}
	return alter, create, drop
}

// TestOracle_RoutineMultiObjectRoundTrip proves a database carrying several routines of both kinds
// (loaded from real engine readbacks) self-diffs empty and yields an empty plan — the realistic
// release-path idempotence check across the whole routine set, not just single forms.
func TestOracle_RoutineMultiObjectRoundTrip(t *testing.T) {
	if testing.Short() {
		t.Skip("oracle test skipped in short mode")
	}
	for _, version := range both() {
		o := connectOracle(t, version)
		sc := serverCharsetFor(o.version)
		n := NormalizerFor(o.version)
		ctx := context.Background()

		t.Run(o.name+"/mixed-routines", func(t *testing.T) {
			dbName := "rtn_multi"
			stmts := []string{
				"DROP DATABASE IF EXISTS " + dbName,
				fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARSET=%s", dbName, sc),
				"USE " + dbName,
				"CREATE FUNCTION add2(a INT, b INT) RETURNS INT DETERMINISTIC RETURN a + b",
				"CREATE FUNCTION greet(s VARCHAR(20)) RETURNS VARCHAR(40) DETERMINISTIC READS SQL DATA COMMENT 'g' RETURN CONCAT('hi ', s)",
				"CREATE PROCEDURE bump(INOUT n INT) BEGIN SET n = n + 1; END",
				"CREATE PROCEDURE noop() MODIFIES SQL DATA BEGIN SET @v = 1; END",
			}
			for _, s := range stmts {
				if _, err := o.db.ExecContext(ctx, s); err != nil {
					t.Skipf("[%s] setup failed (may be expected): %q: %v", o.name, s, err)
				}
			}

			// Read back each routine and reload as a single schema.
			var b strings.Builder
			fmt.Fprintf(&b, "CREATE DATABASE %s DEFAULT CHARSET=%s;\nUSE %s;\n", dbName, sc, dbName)
			routines := []struct{ kind, name string }{
				{"FUNCTION", "add2"}, {"FUNCTION", "greet"},
				{"PROCEDURE", "bump"}, {"PROCEDURE", "noop"},
			}
			for _, r := range routines {
				rb, ok := o.showCreateRoutine(t, dbName, r.kind, r.name)
				if !ok {
					t.Fatalf("[%s] readback %s %s failed", o.name, r.kind, r.name)
				}
				b.WriteString(rb)
				b.WriteString(";\n")
			}
			cat, err := LoadSQL(b.String())
			if err != nil {
				t.Fatalf("[%s] reload of multi-routine readback failed: %v\n%s", o.name, err, b.String())
			}
			diff := DiffWithNormalizer(cat, cat, n)
			if !diff.IsEmpty() {
				t.Errorf("[%s] multi-routine self-diff not empty: %s", o.name, describeRoutineDiff(diff))
			}
			plan := GenerateMigrationWithNormalizer(cat, cat, diff, n)
			if plan.SQL() != "" {
				t.Errorf("[%s] multi-routine no-op plan not empty:\n%s", o.name, plan.SQL())
			}
		})
	}
}
