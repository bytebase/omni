package catalog

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
)

// The cases a code review on this engine's MariaDB migration turned up, each
// one a place the catalog or the view deparser had taken MySQL's answer. Every
// expectation here was read from a live mariadb:11.8.8, and the oracle side
// re-reads it, so these stay honest if the server changes.

// TestReview_RegexpPredicatePrecedence covers REGEXP's place in the precedence
// table. It sits with IN and LIKE, above comparison, so a comparison operand
// keeps its parentheses — and the parentheses change the value:
// `(2 = 1) REGEXP '0'` is 1 while `2 = 1 REGEXP '0'` is 0.
func TestReview_RegexpPredicatePrecedence(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)

	runOnBoth(t, mc, c, `CREATE TABLE t (a INT, b INT);`)
	for _, tc := range []struct{ name, ddl, view string }{
		{"comparison_operand_keeps_parens", `CREATE VIEW v_rx1 AS SELECT (a = 1) REGEXP '0' AS x FROM t;`, "v_rx1"},
		{"regexp_under_comparison_is_bare", `CREATE VIEW v_rx2 AS SELECT a REGEXP 'x' = 1 AS x FROM t;`, "v_rx2"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			runOnBoth(t, mc, c, tc.ddl)
			assertViewBodyMatches(t, mc, c, tc.view)
		})
	}
}

// TestReview_ExistsLimitIsReducedNotReplaced covers the LIMIT an EXISTS
// subquery is stored with. The server caps it at one row only where that
// cannot change the answer: a LIMIT 0 and any offset are left alone, since
// replacing either turns "no rows" into "a row".
func TestReview_ExistsLimitIsReducedNotReplaced(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)

	runOnBoth(t, mc, c, `CREATE TABLE t (a INT);`)
	for _, tc := range []struct{ name, ddl, view string }{
		{"no_limit_gains_one", `CREATE VIEW v_ex1 AS SELECT a FROM t WHERE EXISTS(SELECT 1 FROM t);`, "v_ex1"},
		{"limit_above_one_is_reduced", `CREATE VIEW v_ex2 AS SELECT a FROM t WHERE EXISTS(SELECT 1 FROM t LIMIT 5);`, "v_ex2"},
		{"limit_zero_is_kept", `CREATE VIEW v_ex3 AS SELECT a FROM t WHERE EXISTS(SELECT 1 FROM t LIMIT 0);`, "v_ex3"},
		{"offset_is_kept", `CREATE VIEW v_ex4 AS SELECT a FROM t WHERE EXISTS(SELECT 1 FROM t LIMIT 1 OFFSET 1);`, "v_ex4"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			runOnBoth(t, mc, c, tc.ddl)
			assertViewBodyMatches(t, mc, c, tc.view)
		})
	}
}

// TestReview_NestedNamedWindowIsInlined covers a window reference that is not
// at the top of its target. The WINDOW clause is dropped either way, so a
// reference left behind would name a window the body does not define.
func TestReview_NestedNamedWindowIsInlined(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)

	runOnBoth(t, mc, c, `CREATE TABLE t (a INT, b INT);`)
	runOnBoth(t, mc, c,
		`CREATE VIEW v_win AS SELECT 1 + ROW_NUMBER() OVER w AS x FROM t WINDOW w AS (ORDER BY a);`)
	assertViewBodyMatches(t, mc, c, "v_win")

	body := c.ShowCreateView("testdb", "v_win")
	if strings.Contains(strings.ToLower(body), "window `w`") {
		t.Errorf("the WINDOW clause should be gone from the stored body: %s", body)
	}
	if strings.Contains(body, "over `w`") {
		t.Errorf("the window reference should have been inlined: %s", body)
	}
}

// TestReview_NondeterministicGeneratedColumn covers where a non-deterministic
// expression is allowed: a VIRTUAL column no index covers, and nowhere else.
func TestReview_NondeterministicGeneratedColumn(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct {
		name    string
		ddl     string
		accepts bool
	}{
		{"virtual_unindexed_accepted", "CREATE TABLE t (a DOUBLE AS (RAND()) VIRTUAL)", true},
		{"stored_refused", "CREATE TABLE t (a DOUBLE AS (RAND()) STORED)", false},
		{"virtual_indexed_refused", "CREATE TABLE t (a DOUBLE AS (RAND()) VIRTUAL, KEY(a))", false},
		{"virtual_unique_refused", "CREATE TABLE t (a DOUBLE AS (RAND()) VIRTUAL UNIQUE)", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)

			_, oracleErr := mc.db.ExecContext(mc.ctx, tc.ddl)
			if tc.accepts && oracleErr != nil {
				t.Errorf("oracle: expected %q to be accepted, got %v", tc.ddl, oracleErr)
			}
			if !tc.accepts && oracleErr == nil {
				t.Errorf("oracle: expected %q to be refused", tc.ddl)
			}

			results, err := c.Exec(tc.ddl+";", nil)
			omniErrored := err != nil
			for _, r := range results {
				if r.Error != nil {
					omniErrored = true
				}
			}
			if omniErrored == tc.accepts {
				t.Errorf("omni: %q accepted=%v, want accepted=%v", tc.ddl, !omniErrored, tc.accepts)
			}
		})
	}
}

// TestReview_CheckNameScope covers a CHECK constraint's name scope: per table,
// so two tables may each hold a `ck`, but one table may not hold two.
func TestReview_CheckNameScope(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)

	runOnBoth(t, mc, c, `CREATE TABLE t1 (a INT, CONSTRAINT ck CHECK (a > 0));`)

	t.Run("reused_in_another_table", func(t *testing.T) {
		runOnBoth(t, mc, c, `CREATE TABLE t2 (a INT, CONSTRAINT ck CHECK (a < 9));`)
	})

	t.Run("refused_twice_in_one_table", func(t *testing.T) {
		const ddl = "ALTER TABLE t1 ADD CONSTRAINT ck CHECK (a < 10)"
		if _, oracleErr := mc.db.ExecContext(mc.ctx, ddl); oracleErr == nil {
			t.Errorf("oracle: expected a duplicate CHECK name in one table to be refused")
		}
		results, err := c.Exec(ddl+";", nil)
		omniErrored := err != nil
		for _, r := range results {
			if r.Error != nil {
				omniErrored = true
			}
		}
		if !omniErrored {
			t.Errorf("omni: expected a duplicate CHECK name in one table to be refused")
		}
	})
}

// TestReview_ExpressionDefaultParens covers which expression defaults keep
// their parentheses. A single function call is written bare; an operator
// expression keeps them, which the DEFAULT grammar requires.
func TestReview_ExpressionDefaultParens(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, ddl, want string }{
		{"function_call_is_bare", "CREATE TABLE t (a VARCHAR(36) DEFAULT (uuid()))", "DEFAULT uuid()"},
		{"nested_call_is_bare", "CREATE TABLE t (a INT DEFAULT (floor(rand() * 100)))", "DEFAULT floor(rand() * 100)"},
		{"arithmetic_keeps_parens", "CREATE TABLE t (a INT DEFAULT (1+1))", "DEFAULT (1 + 1)"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, tc.ddl+";")

			if create := oracleShow(t, mc, "SHOW CREATE TABLE t"); !strings.Contains(create, tc.want) {
				t.Errorf("oracle: expected %q in %s", tc.want, create)
			}
			if create := c.ShowCreateTable("testdb", "t"); !strings.Contains(create, tc.want) {
				t.Errorf("omni: expected %q in %s", tc.want, create)
			}
		})
	}
}

// TestReview_CheckNamedAfterGeneratedColumn covers a table-level CHECK whose
// name happens to match a generated column. A generated column's definition
// carries no check, so the constraint stays table-level; writing it on the
// column would drop it from the output altogether.
func TestReview_CheckNamedAfterGeneratedColumn(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)

	runOnBoth(t, mc, c, `CREATE TABLE t (a INT, b INT AS (a+1), CONSTRAINT b CHECK (a > 0));`)

	oracleCreate := oracleShow(t, mc, "SHOW CREATE TABLE t")
	omniCreate := c.ShowCreateTable("testdb", "t")
	for _, k := range []struct{ label, create string }{
		{"oracle", oracleCreate}, {"omni", omniCreate},
	} {
		if !strings.Contains(k.create, "CONSTRAINT `b` CHECK (`a` > 0)") {
			t.Errorf("%s: expected the check to stay table-level: %s", k.label, k.create)
		}
	}
}

// TestReview_JSONCheckFollowsItsColumn covers the implicit json_valid check,
// which the catalog names after its column: it goes when the column is
// dropped, follows a rename, and comes and goes with the type.
func TestReview_JSONCheckFollowsItsColumn(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, setup, alter, want, absent string }{
		{
			name:  "dropped_with_the_column",
			setup: "CREATE TABLE t (a INT, j JSON)", alter: "ALTER TABLE t DROP COLUMN j",
			absent: "json_valid",
		},
		{
			name:  "renamed_with_the_column",
			setup: "CREATE TABLE t (a INT, j JSON)", alter: "ALTER TABLE t CHANGE COLUMN j k JSON",
			want: "CHECK (json_valid(`k`))", absent: "json_valid(`j`)",
		},
		{
			name:  "added_by_alter",
			setup: "CREATE TABLE t (a INT)", alter: "ALTER TABLE t ADD COLUMN j JSON",
			want: "CHECK (json_valid(`j`))",
		},
		{
			name:  "dropped_when_retyped_away",
			setup: "CREATE TABLE t (a INT, j JSON)", alter: "ALTER TABLE t MODIFY COLUMN j INT",
			absent: "json_valid",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, tc.setup+";")
			runOnBoth(t, mc, c, tc.alter+";")

			for _, k := range []struct{ label, create string }{
				{"oracle", oracleShow(t, mc, "SHOW CREATE TABLE t")},
				{"omni", c.ShowCreateTable("testdb", "t")},
			} {
				if tc.want != "" && !strings.Contains(k.create, tc.want) {
					t.Errorf("%s: expected %q in %s", k.label, tc.want, k.create)
				}
				if tc.absent != "" && strings.Contains(k.create, tc.absent) {
					t.Errorf("%s: expected no %q in %s", k.label, tc.absent, k.create)
				}
			}
		})
	}
}

// TestReview_DropCompositePKColumn covers dropping a column a composite
// primary key is built on. It is refused while the key survives the statement,
// and accepted when the same ALTER drops the key too, whichever order the
// clauses come in.
func TestReview_DropCompositePKColumn(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct {
		name    string
		alter   string
		accepts bool
	}{
		{"column_alone_refused", "ALTER TABLE t DROP COLUMN a", false},
		{"with_the_key_accepted", "ALTER TABLE t DROP COLUMN a, DROP PRIMARY KEY", true},
		{"key_first_accepted", "ALTER TABLE t DROP PRIMARY KEY, DROP COLUMN a", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, "CREATE TABLE t (a INT, b INT, PRIMARY KEY (a, b));")

			_, oracleErr := mc.db.ExecContext(mc.ctx, tc.alter)
			if tc.accepts && oracleErr != nil {
				t.Errorf("oracle: expected %q to be accepted, got %v", tc.alter, oracleErr)
			}
			if !tc.accepts && oracleErr == nil {
				t.Errorf("oracle: expected %q to be refused", tc.alter)
			}

			results, err := c.Exec(tc.alter+";", nil)
			omniErrored := err != nil
			for _, r := range results {
				if r.Error != nil {
					omniErrored = true
				}
			}
			if omniErrored == tc.accepts {
				t.Errorf("omni: %q accepted=%v, want accepted=%v", tc.alter, !omniErrored, tc.accepts)
			}
		})
	}
}

// The two review findings that were not bugs. Both asked omni to improve on
// what MariaDB itself writes, and both are pinned here so the next reader can
// see that the odd-looking output is the engine's.

// TestReview_CTENameIsWrittenUnquoted covers a CTE name MariaDB cannot write
// back: 11.8.8 stores `WITH `select` AS (...)` with the name unquoted, so its
// own SHOW CREATE VIEW fails to parse the result. Quoting it in omni would be
// an improvement on the engine, which is not what this catalog is for.
func TestReview_CTENameIsWrittenUnquoted(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)

	runOnBoth(t, mc, c, "CREATE TABLE t (a INT);")
	runOnBoth(t, mc, c, "CREATE VIEW v_cte AS WITH `select` AS (SELECT a FROM t) SELECT a FROM `select`;")

	stored := oracleViewDefinition(t, mc, "v_cte")
	if !strings.Contains(stored, "with select as (") {
		t.Errorf("oracle: expected the CTE name stored unquoted, got %s", stored)
	}
	// The server cannot read its own definition back, which is the whole
	// reason the quoting is not omni's to add.
	if _, err := mc.db.ExecContext(mc.ctx, "SELECT 1 FROM v_cte"); err == nil {
		t.Errorf("oracle: expected the stored definition to fail to parse")
	}

	omni := c.ShowCreateView("testdb", "v_cte")
	if !strings.Contains(omni, "with select as (") {
		t.Errorf("omni: expected the CTE name written unquoted, got %s", omni)
	}
}

// TestReview_UnionOrderByKeepsWrittenAlias covers ORDER BY over a UNION: the
// alias is echoed as the statement wrote it, and the first branch's target
// keeps that alias, so the reference resolves without rewriting.
func TestReview_UnionOrderByKeepsWrittenAlias(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)

	runOnBoth(t, mc, c, "CREATE TABLE t (a INT, b INT);")
	runOnBoth(t, mc, c, "CREATE VIEW v_un AS SELECT a AS x FROM t UNION SELECT b FROM t ORDER BY x;")

	for _, k := range []struct{ label, body string }{
		{"oracle", oracleViewDefinition(t, mc, "v_un")},
		{"omni", c.ShowCreateView("testdb", "v_un")},
	} {
		if !strings.Contains(k.body, "AS `x`") || !strings.Contains(k.body, "order by `x`") {
			t.Errorf("%s: expected the written alias on both the target and the ORDER BY, got %s", k.label, k.body)
		}
	}
}

// oracleViewDefinition reads a view's stored body from information_schema,
// which works even where SHOW CREATE VIEW cannot parse what it stored.
func oracleViewDefinition(t *testing.T, mc *mariadbContainer, view string) string {
	t.Helper()
	var def string
	oracleScan(t, mc,
		"SELECT VIEW_DEFINITION FROM information_schema.VIEWS "+
			"WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = '"+view+"'", &def)
	return def
}

// assertViewBodyMatches compares the SELECT body the server stored for the
// view with the one the catalog deparsed.
func assertViewBodyMatches(t *testing.T, mc *mariadbContainer, c *Catalog, view string) {
	t.Helper()
	oracleBody := extractSelectBody(oracleShow(t, mc, "SHOW CREATE VIEW "+view))
	omniBody := extractSelectBody(c.ShowCreateView("testdb", view))
	if omniBody == "" {
		t.Fatalf("%s: omni stored no view body", view)
	}
	if oracleBody != omniBody {
		t.Errorf("%s body mismatch:\n--- oracle ---\n%s\n--- omni ---\n%s", view, oracleBody, omniBody)
	}
}

// The third round of review findings, measured the same way.

// TestReview_WindowReferenceExtendsDefinition covers an OVER clause that names
// a window and adds to it. The name supplies what it declares and the
// reference adds what it leaves out, so neither side is lost.
func TestReview_WindowReferenceExtendsDefinition(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)

	runOnBoth(t, mc, c, "CREATE TABLE t (a INT, b INT);")
	for _, tc := range []struct{ name, ddl string }{
		{"w_order_added", "CREATE VIEW w_order_added AS SELECT ROW_NUMBER() OVER (w ORDER BY b) AS x FROM t WINDOW w AS (PARTITION BY a)"},
		{"w_frame_added", "CREATE VIEW w_frame_added AS SELECT SUM(b) OVER (w ROWS UNBOUNDED PRECEDING) AS x FROM t WINDOW w AS (ORDER BY a)"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			runOnBoth(t, mc, c, tc.ddl+";")
			assertViewBodyMatches(t, mc, c, tc.name)
		})
	}
}

// TestReview_WindowNameIsScopedToItsQuery covers where a window name resolves.
// It belongs to the query that declares it, so an inner query keeps its own
// definition even where an outer one declares the same name, and a query
// nested anywhere has its own clause inlined and dropped.
func TestReview_WindowNameIsScopedToItsQuery(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)

	runOnBoth(t, mc, c, "CREATE TABLE t (a INT, b INT);")
	for _, tc := range []struct{ name, ddl string }{
		{"w_shadowed", "CREATE VIEW w_shadowed AS SELECT (SELECT SUM(b) OVER w FROM t WINDOW w AS (ORDER BY b) LIMIT 1) AS x FROM t WINDOW w AS (ORDER BY a)"},
		{"w_subquery", "CREATE VIEW w_subquery AS SELECT (SELECT SUM(b) OVER w FROM t WINDOW w AS (ORDER BY b) LIMIT 1) AS x FROM t"},
		{"w_derived", "CREATE VIEW w_derived AS SELECT x FROM (SELECT ROW_NUMBER() OVER w AS x FROM t WINDOW w AS (ORDER BY a)) d"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			runOnBoth(t, mc, c, tc.ddl+";")
			assertViewBodyMatches(t, mc, c, tc.name)

			body := c.ShowCreateView("testdb", tc.name)
			if strings.Contains(strings.ToLower(body), "window `w`") {
				t.Errorf("no WINDOW clause should be left in the body: %s", body)
			}
		})
	}
}

// TestReview_TriggerDDLEchoesItsStatement covers SHOW CREATE TRIGGER, which
// echoes the statement the trigger was created by rather than rebuilding it.
// The writer's quoting, spacing and case all survive, so a name that needs
// backticks keeps them and one that does not never gains them.
func TestReview_TriggerDDLEchoesItsStatement(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)

	runOnBoth(t, mc, c, "CREATE TABLE `order` (a INT);")
	runOnBoth(t, mc, c, "CREATE TABLE tg (a INT);")
	for _, tc := range []struct{ name, ddl string }{
		{"my trigger", "CREATE TRIGGER `my trigger` BEFORE INSERT ON `order` FOR EACH ROW SET NEW.a = 1"},
		{"tr_bare", "CREATE TRIGGER tr_bare BEFORE INSERT ON tg FOR EACH ROW SET NEW.a = 1"},
		{"tr_spaced", "create   trigger   tr_spaced   before  insert on   tg\n  for each row\n   set NEW.a = 1"},
		{"tr_ifne", "CREATE TRIGGER IF NOT EXISTS tr_ifne BEFORE INSERT ON tg FOR EACH ROW SET NEW.a = 1"},
		{"tr_replace", "CREATE OR REPLACE TRIGGER tr_replace BEFORE INSERT ON tg FOR EACH ROW SET NEW.a = 1"},
		{"tr_follows", "CREATE TRIGGER tr_follows BEFORE INSERT ON tg FOR EACH ROW FOLLOWS tr_bare SET NEW.a = 1"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			runOnBoth(t, mc, c, tc.ddl+";")
			want := oracleColumn(t, mc, "SHOW CREATE TRIGGER "+backtickIDQ(tc.name), 2)
			if got := c.ShowCreateTrigger("testdb", tc.name); got != want {
				t.Errorf("trigger DDL mismatch\n  oracle: %q\n  omni:   %q", want, got)
			}
		})
	}
}

// TestReview_TriggerDDLWithoutAStatement covers the synced trigger, which has
// no statement to echo. Its DDL is rebuilt with every identifier backticked,
// since a name that needs them would otherwise come out unparseable.
func TestReview_TriggerDDLWithoutAStatement(t *testing.T) {
	c := New()
	results, err := c.Exec("CREATE DATABASE d; USE d; CREATE TABLE `order` (a INT);", nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatal(r.Error)
		}
	}
	stmt := mustParseTrigger(t, "CREATE TRIGGER `my trigger` BEFORE INSERT ON `order` FOR EACH ROW SET NEW.a = 1")
	if err := c.DefineTrigger(stmt); err != nil {
		t.Fatalf("DefineTrigger: %v", err)
	}
	got := c.ShowCreateTrigger("d", "my trigger")
	const want = "CREATE DEFINER=`root`@`%` TRIGGER `my trigger` BEFORE INSERT ON `order` FOR EACH ROW SET NEW.a = 1"
	if got != want {
		t.Errorf("rebuilt trigger DDL:\n  got:  %q\n  want: %q", got, want)
	}
}

// TestReview_FunctionReturnsCollation covers which collation a RETURNS clause
// is shown with: the one it names, else the charset's default where it names a
// charset, else the database's own.
func TestReview_FunctionReturnsCollation(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)

	fns := []string{
		"CREATE FUNCTION f_plain() RETURNS VARCHAR(20) RETURN 'a'",
		"CREATE FUNCTION f_cs_other() RETURNS VARCHAR(20) CHARACTER SET utf8mb3 RETURN 'a'",
		"CREATE FUNCTION f_cs_same() RETURNS VARCHAR(20) CHARACTER SET utf8mb4 RETURN 'a'",
		"CREATE FUNCTION f_text() RETURNS TEXT RETURN 'a'",
		"CREATE FUNCTION f_collate() RETURNS VARCHAR(20) COLLATE utf8mb4_general_ci RETURN 'a'",
		"CREATE FUNCTION f_int() RETURNS INT RETURN 1",
	}
	const createDB = "CREATE DATABASE fdb CHARACTER SET utf8mb4 COLLATE utf8mb4_bin"
	for _, stmt := range append([]string{"DROP DATABASE IF EXISTS fdb", createDB, "USE fdb"}, fns...) {
		if _, err := mc.db.ExecContext(mc.ctx, stmt); err != nil {
			t.Fatalf("oracle %q: %v", stmt, err)
		}
	}

	c := New()
	ddl := createDB + "; USE fdb;"
	for _, fn := range fns {
		ddl += fn + ";"
	}
	results, err := c.Exec(ddl, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatalf("omni stmt %d: %v", r.Index, r.Error)
		}
	}

	for _, name := range []string{
		"f_plain", "f_cs_other", "f_cs_same", "f_text", "f_collate", "f_int",
	} {
		t.Run(name, func(t *testing.T) {
			want := oracleColumn(t, mc, "SHOW CREATE FUNCTION "+name, 2)
			if got := c.ShowCreateFunction("fdb", name); got != want {
				t.Errorf("function DDL mismatch\n  oracle: %q\n  omni:   %q", want, got)
			}
		})
	}
}

// TestReview_QuantifiedComparisonDropsParens is the finding that was not a
// bug. MariaDB writes a quantified comparison's operands with no parentheses
// of its own, and its own stored view then means something else: the
// expression is 0 where the view it was stored in reads 1. Adding the
// parentheses would make omni disagree with the engine it mirrors.
func TestReview_QuantifiedComparisonDropsParens(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)

	runOnBoth(t, mc, c, "CREATE TABLE t (a INT);")
	runOnBoth(t, mc, c, "INSERT INTO t VALUES (1);")
	runOnBoth(t, mc, c, "CREATE VIEW q AS SELECT (1 OR 0) > ALL (SELECT 2) AS x FROM t;")
	assertViewBodyMatches(t, mc, c, "q")

	var direct, throughView int
	oracleScan(t, mc, "SELECT (1 OR 0) > ALL (SELECT 2)", &direct)
	oracleScan(t, mc, "SELECT x FROM q", &throughView)
	if direct != 0 || throughView != 1 {
		t.Errorf("oracle: expected the expression to be 0 and the stored view to read 1, got %d and %d",
			direct, throughView)
	}
}

// oracleColumn runs a SHOW statement and returns one column of its first row.
// SHOW CREATE TRIGGER and SHOW CREATE FUNCTION both hold their DDL in column
// 2, after the name and the sql_mode.
func oracleColumn(t *testing.T, mc *mariadbContainer, stmt string, col int) string {
	t.Helper()
	rows := oracleRows(t, mc, stmt)
	if len(rows) == 0 || col >= len(rows[0]) {
		t.Errorf("oracleColumn: %q gave no column %d", stmt, col)
		return ""
	}
	s, _ := rows[0][col].(string)
	return s
}

// TestReview_WindowReferenceRules covers what a window reference may do. A
// reference extends the window it names, but only with what that window
// leaves out, and the names have to resolve — so every case here is compared
// with the server by error code and message, the accepted ones included.
//
// The one that reads oddly is real: `OVER (w)` is refused where w has a frame
// while the bare `OVER w` inherits it, because the first is a specification
// referencing a window and the second is only a reference.
func TestReview_WindowReferenceRules(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	if _, err := mc.db.ExecContext(mc.ctx, "CREATE TABLE t (a INT, b INT)"); err != nil {
		t.Fatal(err)
	}

	for i, tc := range []struct{ name, query string }{
		{"adds_partition", "SELECT ROW_NUMBER() OVER (w PARTITION BY b) FROM t WINDOW w AS (ORDER BY a)"},
		{"adds_partition_to_empty", "SELECT ROW_NUMBER() OVER (w PARTITION BY b) FROM t WINDOW w AS ()"},
		{"adds_order_to_ordered", "SELECT ROW_NUMBER() OVER (w ORDER BY b) FROM t WINDOW w AS (ORDER BY a)"},
		{"adds_order_to_partitioned", "SELECT ROW_NUMBER() OVER (w ORDER BY b) FROM t WINDOW w AS (PARTITION BY a)"},
		{"spec_refs_framed", "SELECT SUM(b) OVER (w ROWS UNBOUNDED PRECEDING) FROM t WINDOW w AS (ORDER BY a ROWS 2 PRECEDING)"},
		{"empty_spec_refs_framed", "SELECT SUM(b) OVER (w) FROM t WINDOW w AS (ORDER BY a ROWS 2 PRECEDING)"},
		{"bare_ref_inherits_frame", "SELECT SUM(b) OVER w FROM t WINDOW w AS (ORDER BY a ROWS 2 PRECEDING)"},
		{"adds_order_to_framed", "SELECT SUM(b) OVER (w ORDER BY b) FROM t WINDOW w AS (ROWS 2 PRECEDING)"},
		{"adds_frame", "SELECT SUM(b) OVER (w ROWS UNBOUNDED PRECEDING) FROM t WINDOW w AS (ORDER BY a)"},
		{"clause_entry_refs", "SELECT ROW_NUMBER() OVER w2 FROM t WINDOW w1 AS (ORDER BY a), w2 AS (w1)"},
		{"clause_entry_adds_partition", "SELECT ROW_NUMBER() OVER w2 FROM t WINDOW w1 AS (ORDER BY a), w2 AS (w1 PARTITION BY b)"},
		{"clause_entry_adds_order", "SELECT ROW_NUMBER() OVER w2 FROM t WINDOW w1 AS (PARTITION BY a), w2 AS (w1 ORDER BY b)"},
		{"clause_entry_refs_framed", "SELECT ROW_NUMBER() OVER w2 FROM t WINDOW w1 AS (ORDER BY a ROWS 2 PRECEDING), w2 AS (w1)"},
		{"clause_entry_order_on_framed", "SELECT ROW_NUMBER() OVER w2 FROM t WINDOW w1 AS (ORDER BY a ROWS 2 PRECEDING), w2 AS (w1 ORDER BY b)"},
		{"forward_reference", "SELECT ROW_NUMBER() OVER w1 FROM t WINDOW w1 AS (w2), w2 AS (ORDER BY a)"},
		{"duplicate_name", "SELECT ROW_NUMBER() OVER w FROM t WINDOW w AS (ORDER BY a), w AS (ORDER BY b)"},
		{"unknown_in_spec", "SELECT ROW_NUMBER() OVER (nope ORDER BY b) FROM t"},
		{"unknown_bare", "SELECT ROW_NUMBER() OVER nope FROM t"},
		{"nested_scope_is_its_own", "SELECT (SELECT SUM(b) OVER w FROM t WINDOW w AS (ORDER BY b) LIMIT 1) FROM t WINDOW w AS (ORDER BY a)"},
		{"nested_unknown_is_caught", "SELECT (SELECT SUM(b) OVER nope FROM t LIMIT 1) FROM t WINDOW w AS (ORDER BY a)"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var oracleMsg string
			view := fmt.Sprintf("wv%d", i)
			if _, err := mc.db.ExecContext(mc.ctx, "CREATE VIEW "+view+" AS "+tc.query); err != nil {
				oracleMsg = serverErrorMessage(err)
			}

			c := New()
			results, err := c.Exec(
				"CREATE DATABASE d; USE d; CREATE TABLE t (a INT, b INT); CREATE VIEW v AS "+tc.query+";", nil)
			if err != nil {
				t.Fatalf("omni parse error: %v", err)
			}
			var omniMsg string
			for _, r := range results {
				if r.Error == nil {
					continue
				}
				e, ok := r.Error.(*Error)
				if !ok {
					t.Fatalf("omni returned %T, want *Error: %v", r.Error, r.Error)
				}
				omniMsg = fmt.Sprintf("%d %s", e.Code, e.Message)
			}

			if omniMsg != oracleMsg {
				t.Errorf("answers differ\n  oracle: %q\n  omni:   %q", oracleMsg, omniMsg)
			}
		})
	}
}

// serverErrorMessage renders a driver error as the code and message, dropping
// the SQLSTATE so it can be compared with a catalog Error.
func serverErrorMessage(err error) string {
	s := err.Error()
	var code int
	if n, _ := fmt.Sscanf(s, "Error %d", &code); n != 1 {
		return s
	}
	if i := strings.Index(s, "): "); i >= 0 {
		return fmt.Sprintf("%d %s", code, s[i+len("): "):])
	}
	return s
}

// The fourth round of review findings.

// TestReview_ExpressionDefaultQuotedParens covers an expression default whose
// string literals hold parentheses. Those are text, not structure, so they do
// not decide whether the default is a single call: `length('(') + length(')')`
// is two calls around an operator and keeps its parentheses, while
// `length('()')` is one call and loses them.
func TestReview_ExpressionDefaultQuotedParens(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, ddl string }{
		{"quoted_parens_around_operator", "CREATE TABLE t (a INT DEFAULT (length('(') + length(')')))"},
		{"quoted_parens_in_one_call", "CREATE TABLE t (a INT DEFAULT (length('()')))"},
		{"quoted_close_paren_arg", "CREATE TABLE t (a VARCHAR(9) DEFAULT (concat('a', ')')))"},
		{"quoted_open_paren_arg", "CREATE TABLE t (a VARCHAR(9) DEFAULT (concat('a', '(')))"},
		{"operator_keeps_parens", "CREATE TABLE t (a INT DEFAULT (1 + 1))"},
		{"single_call_is_bare", "CREATE TABLE t (a VARCHAR(36) DEFAULT (uuid()))"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, tc.ddl+";")
			assertDefaultMatches(t, mc, c)
		})
	}
}

// TestReview_BinaryLiteralDefault covers a literal default on a binary
// column. The two halves of the family answer differently: BINARY and
// VARBINARY store the coerced value and write it back as lowercase hex, zero
// padded to the width where the type is fixed, while the BLOB family stores
// the expression and keeps the literal as written, only normalized.
func TestReview_BinaryLiteralDefault(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, ddl string }{
		{"binary_hex", "CREATE TABLE t (a BINARY(1) DEFAULT X'FF')"},
		{"binary_hex_padded", "CREATE TABLE t (a BINARY(4) DEFAULT X'FF')"},
		{"binary_0x", "CREATE TABLE t (a BINARY(1) DEFAULT 0xFF)"},
		{"binary_0x_padded", "CREATE TABLE t (a BINARY(3) DEFAULT 0x4142)"},
		{"binary_bits_padded", "CREATE TABLE t (a BINARY(2) DEFAULT b'101')"},
		{"binary_string_padded", "CREATE TABLE t (a BINARY(4) DEFAULT 'a')"},
		{"binary_empty_hex", "CREATE TABLE t (a BINARY(1) DEFAULT X'')"},
		{"binary_parenthesized", "CREATE TABLE t (a BINARY(1) DEFAULT (X'FF'))"},
		{"varbinary_hex", "CREATE TABLE t (a VARBINARY(4) DEFAULT X'FF')"},
		{"varbinary_string", "CREATE TABLE t (a VARBINARY(4) DEFAULT 'a')"},
		{"varbinary_bits", "CREATE TABLE t (a VARBINARY(4) DEFAULT b'1000001')"},
		{"varbinary_expression", "CREATE TABLE t (a VARBINARY(4) DEFAULT (concat(X'FF', X'00')))"},
		{"blob_hex_keeps_spelling", "CREATE TABLE t (a BLOB DEFAULT X'FF')"},
		{"blob_hex_two_bytes", "CREATE TABLE t (a BLOB DEFAULT X'4142')"},
		{"blob_0x_keeps_spelling", "CREATE TABLE t (a BLOB DEFAULT 0xFF)"},
		{"blob_bits_become_0x", "CREATE TABLE t (a BLOB DEFAULT b'1000001')"},
		{"blob_string_stays_string", "CREATE TABLE t (a BLOB DEFAULT 'a')"},
		{"blob_parenthesized", "CREATE TABLE t (a BLOB DEFAULT (X'FF'))"},
		{"tinyblob_hex", "CREATE TABLE t (a TINYBLOB DEFAULT X'FF')"},
		{"longblob_hex", "CREATE TABLE t (a LONGBLOB DEFAULT X'FF')"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, tc.ddl+";")
			assertDefaultMatches(t, mc, c)
		})
	}
}

// TestReview_QualifiedFunctionCall covers a schema-qualified call, which names
// a stored routine rather than a built-in. It keeps the name it was declared
// with and is written quoted, where an unqualified call is folded onto the
// built-in's canonical name even if a routine of that name exists.
func TestReview_QualifiedFunctionCall(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)

	runOnBoth(t, mc, c, "CREATE TABLE t (a INT);")
	for _, ddl := range []string{
		"CREATE FUNCTION `upper`(x INT) RETURNS INT DETERMINISTIC RETURN x+1;",
		"CREATE FUNCTION plainfn(x INT) RETURNS INT DETERMINISTIC RETURN x+2;",
		"CREATE FUNCTION `odd name`(x INT) RETURNS INT DETERMINISTIC RETURN x+3;",
		"CREATE FUNCTION MixedFn(x INT) RETURNS INT DETERMINISTIC RETURN x+4;",
	} {
		runOnBoth(t, mc, c, ddl)
	}

	for _, tc := range []struct{ name, view, query string }{
		{"qualified_shadows_builtin", "qv_upper", "SELECT testdb.`upper`(1) AS x"},
		{"qualified_plain", "qv_plain", "SELECT testdb.plainfn(1) AS x"},
		{"qualified_needs_quoting", "qv_odd", "SELECT testdb.`odd name`(1) AS x"},
		{"qualified_keeps_case", "qv_mixed", "SELECT testdb.MixedFn(1) AS x"},
		{"unqualified_folds_to_builtin", "qv_bare", "SELECT `upper`(1) AS x"},
		{"builtin_folds", "qv_builtin", "SELECT UPPER('a') AS x"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			runOnBoth(t, mc, c, "CREATE VIEW "+tc.view+" AS "+tc.query+";")
			want := oracleViewDefinition(t, mc, tc.view)
			if got := extractSelectBody(c.ShowCreateView("testdb", tc.view)); got != want {
				t.Errorf("view body mismatch\n  oracle: %s\n  omni:   %s", want, got)
			}
		})
	}
}

// TestReview_ByteLengthAndMidNames covers two names a stored expression is
// written back under, in a view body and a DEFAULT alike: LENGTH counts bytes
// and is stored as octet_length, and MID is SUBSTRING's third spelling.
func TestReview_ByteLengthAndMidNames(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)

	for i, tc := range []struct{ name, expr string }{
		{"length", "length('ab')"},
		{"mid", "mid('abc',1,1)"},
		{"octet_length", "octet_length('ab')"},
		{"char_length", "char_length('ab')"},
		{"substring", "substring('abc',1)"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			view := fmt.Sprintf("nv%d", i)
			runOnBoth(t, mc, c, "CREATE VIEW "+view+" AS SELECT "+tc.expr+" AS x;")
			want := oracleViewDefinition(t, mc, view)
			if got := extractSelectBody(c.ShowCreateView("testdb", view)); got != want {
				t.Errorf("view body mismatch\n  oracle: %s\n  omni:   %s", want, got)
			}
		})
	}
}

// TestReview_ValuesViewDropsColumnNames is the finding that was not a bug.
// MariaDB discards the column list a VALUES view declares: the columns of
// `CREATE VIEW v(x, y) AS VALUES (1, 2)` are named 1 and 2 on the server
// itself, and naming x fails, so there is no interface for omni to preserve.
func TestReview_ValuesViewDropsColumnNames(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)

	runOnBoth(t, mc, c, "CREATE VIEW v_values(x, y) AS VALUES (1, 2);")
	assertViewBodyMatches(t, mc, c, "v_values")

	var names []string
	for _, row := range oracleRows(t, mc,
		"SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE()"+
			" AND TABLE_NAME = 'v_values' ORDER BY ORDINAL_POSITION") {
		if s, ok := row[0].(string); ok {
			names = append(names, s)
		}
	}
	if len(names) != 2 || names[0] != "1" || names[1] != "2" {
		t.Errorf("oracle: expected the declared names to be dropped for 1 and 2, got %v", names)
	}
	if _, err := mc.db.ExecContext(mc.ctx, "SELECT x FROM v_values"); err == nil {
		t.Errorf("oracle: expected the declared column name to be unusable")
	}

	// A SELECT body keeps its list, which is why that one is dropped from the
	// stored DDL rather than kept.
	runOnBoth(t, mc, c, "CREATE VIEW v_select(x, y) AS SELECT 1, 2;")
	assertViewBodyMatches(t, mc, c, "v_select")
}

// assertDefaultMatches compares the DEFAULT clause the server stored for table
// t with the one the catalog wrote.
func assertDefaultMatches(t *testing.T, mc *mariadbContainer, c *Catalog) {
	t.Helper()
	want := defaultClauseOf(oracleShow(t, mc, "SHOW CREATE TABLE t"))
	if got := defaultClauseOf(c.ShowCreateTable("testdb", "t")); got != want {
		t.Errorf("DEFAULT clause mismatch\n  oracle: %s\n  omni:   %s", want, got)
	}
}

// defaultClauseOf pulls the DEFAULT clause out of a SHOW CREATE TABLE.
func defaultClauseOf(create string) string {
	for _, line := range strings.Split(create, "\n") {
		if i := strings.Index(line, "DEFAULT "); i >= 0 {
			return strings.TrimRight(strings.TrimSpace(line[i:]), ",")
		}
	}
	return "(no DEFAULT)"
}

// TestReview_JSONColumnCopiedByLike covers CREATE TABLE ... LIKE over a JSON
// column. The copy is a JSON column too, so the marker travels with it and a
// later ALTER reads the column for what it is rather than as newly typed.
func TestReview_JSONColumnCopiedByLike(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct {
		name  string
		stmts []string
	}{
		{"plain_copy", []string{
			"CREATE TABLE t (j JSON)",
			"CREATE TABLE t2 LIKE t",
		}},
		{"copy_then_retype_to_json", []string{
			"CREATE TABLE t (j JSON)",
			"CREATE TABLE t2 LIKE t",
			"ALTER TABLE t2 MODIFY COLUMN j JSON",
		}},
		{"copy_then_retype_away", []string{
			"CREATE TABLE t (j JSON)",
			"CREATE TABLE t2 LIKE t",
			"ALTER TABLE t2 MODIFY COLUMN j INT",
		}},
		{"copy_then_drop_column", []string{
			"CREATE TABLE t (a INT, j JSON)",
			"CREATE TABLE t2 LIKE t",
			"ALTER TABLE t2 DROP COLUMN j",
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			for _, stmt := range tc.stmts {
				runOnBoth(t, mc, c, stmt+";")
			}
			want := oracleShow(t, mc, "SHOW CREATE TABLE t2")
			if got := c.ShowCreateTable("testdb", "t2"); got != want {
				t.Errorf("table mismatch\n--- oracle ---\n%s\n--- omni ---\n%s", want, got)
			}
		})
	}
}

// TestReview_ImplicitJSONCheckIsNotDroppableByName covers the implicit check's
// name, which is the column's: the server does not accept it as a constraint
// to drop, answering 1091.
func TestReview_ImplicitJSONCheckIsNotDroppableByName(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)

	runOnBoth(t, mc, c, "CREATE TABLE t (j JSON);")

	const drop = "ALTER TABLE t DROP CONSTRAINT j"
	_, oracleErr := mc.db.ExecContext(mc.ctx, drop)
	if oracleErr == nil {
		t.Fatalf("oracle: expected the implicit check to be undroppable by name")
	}

	results, err := c.Exec(drop+";", nil)
	if err != nil {
		t.Fatalf("omni: %v", err)
	}
	var omniErr error
	for _, r := range results {
		if r.Error != nil {
			omniErr = r.Error
		}
	}
	if omniErr == nil {
		t.Errorf("omni: accepted %q, which the server refuses with %v", drop, oracleErr)
	}
}

// TestReview_LiteralDefaultByColumnFamily covers a hex, bit or string literal
// default on every family of column type. A column that holds its default as
// a value coerces the literal and writes it in its own notation; one that
// holds it as an expression keeps the literal as written.
//
// All of it measured on mariadb:11.8.8 through SHOW CREATE TABLE.
func TestReview_LiteralDefaultByColumnFamily(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, ddl string }{
		// char, varchar, enum and set hold a value, written as a string.
		{"varchar_0x", "CREATE TABLE t (a VARCHAR(4) DEFAULT 0x41)"},
		{"varchar_hex", "CREATE TABLE t (a VARCHAR(4) DEFAULT X'41')"},
		{"varchar_bits", "CREATE TABLE t (a VARCHAR(8) DEFAULT b'1000001')"},
		{"varchar_parenthesized", "CREATE TABLE t (a VARCHAR(4) DEFAULT (0x41))"},
		{"char_0x", "CREATE TABLE t (a CHAR(4) DEFAULT 0x41)"},
		{"char_two_bytes", "CREATE TABLE t (a CHAR(4) DEFAULT X'4142')"},
		{"enum_0x", "CREATE TABLE t (a ENUM('A','B') DEFAULT 0x41)"},
		{"set_0x", "CREATE TABLE t (a SET('A','B') DEFAULT 0x41)"},
		{"varchar_plain_string", "CREATE TABLE t (a VARCHAR(8) DEFAULT 'plain')"},
		// The five bytes MariaDB escapes, and some it does not.
		{"escape_quote", "CREATE TABLE t (a VARCHAR(4) DEFAULT 0x2741)"},
		{"escape_backslash", "CREATE TABLE t (a VARCHAR(4) DEFAULT 0x5C41)"},
		{"escape_newline", "CREATE TABLE t (a VARCHAR(4) DEFAULT 0x0A)"},
		{"escape_return", "CREATE TABLE t (a VARCHAR(4) DEFAULT 0x0D)"},
		{"escape_nul", "CREATE TABLE t (a VARCHAR(4) DEFAULT 0x00)"},
		{"raw_tab", "CREATE TABLE t (a VARCHAR(4) DEFAULT 0x09)"},
		{"raw_double_quote", "CREATE TABLE t (a VARCHAR(4) DEFAULT 0x22)"},
		{"raw_utf8", "CREATE TABLE t (a VARCHAR(4) DEFAULT 0xC3A4)"},
		// The text family holds an expression, so the literal survives.
		{"text_0x", "CREATE TABLE t (a TEXT DEFAULT 0x41)"},
		{"text_hex", "CREATE TABLE t (a TEXT DEFAULT X'41')"},
		{"text_bits", "CREATE TABLE t (a TEXT DEFAULT b'1000001')"},
		{"tinytext_0x", "CREATE TABLE t (a TINYTEXT DEFAULT 0x41)"},
		{"longtext_hex", "CREATE TABLE t (a LONGTEXT DEFAULT X'41')"},
		{"text_plain_string", "CREATE TABLE t (a TEXT DEFAULT 'plain')"},
		// A number reads the bytes as one.
		{"int_0x", "CREATE TABLE t (a INT DEFAULT 0xFF)"},
		{"int_0x_two_bytes", "CREATE TABLE t (a INT DEFAULT 0xFFFF)"},
		{"bigint_0x", "CREATE TABLE t (a BIGINT DEFAULT 0xFFFFFFFF)"},
		{"tinyint_0x", "CREATE TABLE t (a TINYINT DEFAULT 0x7F)"},
		{"decimal_keeps_scale", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT 0x41)"},
		{"double_0x", "CREATE TABLE t (a DOUBLE DEFAULT 0x41)"},
		{"year_0x", "CREATE TABLE t (a YEAR DEFAULT 0x07D0)"},
		{"int_bits", "CREATE TABLE t (a INT DEFAULT b'1111111111')"},
		{"int_plain", "CREATE TABLE t (a INT DEFAULT 7)"},
		{"int_quoted_number", "CREATE TABLE t (a INT DEFAULT '7')"},
		// And the expression defaults stay expressions.
		{"expression_operator", "CREATE TABLE t (a INT DEFAULT (1 + 1))"},
		{"expression_call", "CREATE TABLE t (a VARCHAR(36) DEFAULT (uuid()))"},
		{"expression_quoted_parens", "CREATE TABLE t (a INT DEFAULT (length('(') + length(')')))"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, tc.ddl+";")
			assertDefaultMatches(t, mc, c)
		})
	}
}

// TestReview_WindowClauseEntryChain covers a WINDOW clause entry that
// references an earlier entry. Each entry is resolved against the ones before
// it, so a window reached through a chain carries everything it inherits.
func TestReview_WindowClauseEntryChain(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)

	runOnBoth(t, mc, c, "CREATE TABLE t (a INT, b INT);")
	for _, tc := range []struct{ name, ddl string }{
		{"entry_extends_entry", "CREATE VIEW wc1 AS SELECT SUM(b) OVER w2 AS x FROM t" +
			" WINDOW w1 AS (PARTITION BY a), w2 AS (w1 ORDER BY b)"},
		{"chain_of_three", "CREATE VIEW wc2 AS SELECT SUM(b) OVER w3 AS x FROM t" +
			" WINDOW w1 AS (PARTITION BY a), w2 AS (w1), w3 AS (w2 ORDER BY b)"},
		{"spec_over_a_chain", "CREATE VIEW wc3 AS SELECT SUM(b) OVER (w2) AS x FROM t" +
			" WINDOW w1 AS (PARTITION BY a), w2 AS (w1 ORDER BY b)"},
		{"frame_over_a_chain", "CREATE VIEW wc4 AS SELECT SUM(b) OVER (w2 ROWS UNBOUNDED PRECEDING) AS x FROM t" +
			" WINDOW w1 AS (PARTITION BY a), w2 AS (w1 ORDER BY b)"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			runOnBoth(t, mc, c, tc.ddl+";")
			assertViewBodyMatches(t, mc, c, strings.Fields(tc.ddl)[2])
		})
	}
}

// TestReview_PositionIsStoredAsLocate covers POSITION, which MariaDB stores as
// the locate() call with its arguments separated by a comma. That is also what
// keeps the output parseable: POSITION's first argument is a bit_expr, so a
// lower-precedence operand there would need parentheses the engine does not
// write, while locate's arguments carry no such constraint.
func TestReview_PositionIsStoredAsLocate(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)

	runOnBoth(t, mc, c, "CREATE TABLE t (a INT);")
	for _, tc := range []struct{ name, view, query string }{
		{"plain_strings", "pv1", "SELECT POSITION('1' IN '1') AS x FROM t"},
		{"comparison_operand", "pv2", "SELECT POSITION((1 = 1) IN '1') AS x FROM t"},
		{"arithmetic_operand", "pv3", "SELECT POSITION((a + 1) IN '1') AS x FROM t"},
		{"locate_itself", "pv4", "SELECT LOCATE((1 = 1), '1') AS x FROM t"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			runOnBoth(t, mc, c, "CREATE VIEW "+tc.view+" AS "+tc.query+";")
			assertViewBodyMatches(t, mc, c, tc.view)

			// The stored body has to be something the server can read back.
			body := extractSelectBody(c.ShowCreateView("testdb", tc.view))
			if _, err := mc.db.ExecContext(mc.ctx, "CREATE VIEW re_"+tc.view+" AS "+body); err != nil {
				t.Errorf("the stored body does not reparse: %v\n  %s", err, body)
			}
		})
	}
}

// TestReview_IntroducedAndScaledDefaults covers a default written with a
// character-set introducer, which is another spelling of a string literal and
// so follows the same family rule, and the scale a DECIMAL carries whatever
// the default was written as.
func TestReview_IntroducedAndScaledDefaults(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, ddl string }{
		// An introducer on a column that holds a value: coerced, introducer gone.
		{"int_introduced", "CREATE TABLE t (a INT DEFAULT _utf8mb4'7')"},
		{"int_introduced_latin1", "CREATE TABLE t (a INT DEFAULT _latin1'7')"},
		{"varchar_introduced", "CREATE TABLE t (a VARCHAR(4) DEFAULT _utf8mb4'x')"},
		{"varchar_introduced_latin1", "CREATE TABLE t (a VARCHAR(4) DEFAULT _latin1'x')"},
		{"binary_introduced", "CREATE TABLE t (a BINARY(2) DEFAULT _binary'ab')"},
		{"decimal_introduced", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT _utf8mb4'7')"},
		// On one that holds an expression, the literal survives as written.
		{"text_introduced", "CREATE TABLE t (a TEXT DEFAULT _utf8mb4'x')"},
		// N'' is the national-charset spelling, which the parser folds away.
		{"int_national", "CREATE TABLE t (a INT DEFAULT N'7')"},
		{"varchar_national", "CREATE TABLE t (a VARCHAR(4) DEFAULT N'x')"},
		// A DECIMAL shows its declared scale however the default was written.
		{"decimal_plain", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT 7)"},
		{"decimal_quoted", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT '7')"},
		{"decimal_fraction", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT 7.5)"},
		{"decimal_quoted_fraction", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT '7.5')"},
		{"decimal_hex", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT 0x41)"},
		{"decimal_no_scale", "CREATE TABLE t (a DECIMAL(5,0) DEFAULT 7)"},
		{"double_keeps_fraction", "CREATE TABLE t (a DOUBLE DEFAULT '7.5')"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, tc.ddl+";")
			assertDefaultMatches(t, mc, c)
		})
	}
}

// TestReview_ExponentNumberDefault covers a number default in exponent
// notation, which is worked out rather than carried: padding it to a DECIMAL's
// scale would have produced `1e2.00`, which does not parse.
func TestReview_ExponentNumberDefault(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, ddl string }{
		{"decimal_exponent", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT 1e2)"},
		{"decimal_exponent_quoted", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT '1e2')"},
		{"decimal_exponent_upper", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT 1E2)"},
		{"decimal_exponent_fraction", "CREATE TABLE t (a DECIMAL(6,2) DEFAULT 1.5e2)"},
		{"decimal_exponent_negative", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT 1e-1)"},
		{"double_exponent", "CREATE TABLE t (a DOUBLE DEFAULT 1e2)"},
		{"int_exponent", "CREATE TABLE t (a INT DEFAULT 1e2)"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, tc.ddl+";")
			assertDefaultMatches(t, mc, c)

			// The written default has to be something the server can read.
			create := c.ShowCreateTable("testdb", "t")
			if _, err := mc.db.ExecContext(mc.ctx,
				strings.Replace(create, "CREATE TABLE `t`", "CREATE TABLE `t_reparse`", 1)); err != nil {
				t.Errorf("the written table does not reparse: %v\n%s", err, create)
			}
		})
	}
}

// TestReview_AnalyzedChecksSkipTheImplicitOne covers the analyzed CHECK
// expressions a table exposes. They are matched to the statement's CHECK
// clauses in order, and a JSON column's implicit json_valid check is not one
// of those clauses, so it is passed over rather than counted.
func TestReview_AnalyzedChecksSkipTheImplicitOne(t *testing.T) {
	for _, tc := range []struct {
		name     string
		ddl      string
		analyzed map[string]bool
	}{
		{
			name:     "json_column_before_a_check",
			ddl:      "CREATE TABLE t (j JSON, a INT, CONSTRAINT positive CHECK (a > 0))",
			analyzed: map[string]bool{"j": false, "positive": true},
		},
		{
			name: "json_column_between_two_checks",
			ddl: "CREATE TABLE t (a INT, CONSTRAINT lo CHECK (a > 0), j JSON," +
				" CONSTRAINT hi CHECK (a < 9))",
			analyzed: map[string]bool{"lo": true, "j": false, "hi": true},
		},
		{
			name:     "column_check_after_a_json_column",
			ddl:      "CREATE TABLE t (j JSON, a INT CHECK (a > 0))",
			analyzed: map[string]bool{"j": false, "a": true},
		},
		{
			name:     "no_json_column",
			ddl:      "CREATE TABLE t (a INT, CONSTRAINT positive CHECK (a > 0))",
			analyzed: map[string]bool{"positive": true},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := New()
			results, err := c.Exec("CREATE DATABASE d; USE d; "+tc.ddl+";", nil)
			if err != nil {
				t.Fatal(err)
			}
			for _, r := range results {
				if r.Error != nil {
					t.Fatal(r.Error)
				}
			}
			tbl := c.GetDatabase("d").GetTable("t")
			seen := map[string]bool{}
			for _, con := range tbl.Constraints {
				if con.Type != ConCheck {
					continue
				}
				seen[con.Name] = con.CheckAnalyzed != nil
			}
			for name, want := range tc.analyzed {
				got, ok := seen[name]
				if !ok {
					t.Errorf("no CHECK constraint named %q; got %v", name, seen)
					continue
				}
				if got != want {
					t.Errorf("constraint %q analyzed=%v, want %v", name, got, want)
				}
			}
			if len(seen) != len(tc.analyzed) {
				t.Errorf("CHECK constraints: got %v, want %v", seen, tc.analyzed)
			}
		})
	}
}

// TestReview_ExpressionOfSeveralLiterals covers an expression default that
// opens and closes with a string literal without being one. Reading its
// source text as a value encoded the text itself on a binary column, and on a
// text one dropped the parentheses the DEFAULT grammar needs.
func TestReview_ExpressionOfSeveralLiterals(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, ddl string }{
		{"binary", "CREATE TABLE t (a BINARY(8) DEFAULT ('1' + '2'))"},
		{"text", "CREATE TABLE t (a TEXT DEFAULT ('1' + '2'))"},
		{"varchar", "CREATE TABLE t (a VARCHAR(8) DEFAULT ('1' + '2'))"},
		{"blob", "CREATE TABLE t (a BLOB DEFAULT ('1' + '2'))"},
		{"three_literals", "CREATE TABLE t (a VARCHAR(8) DEFAULT ('a' + 'b' + 'c'))"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, tc.ddl+";")
			assertDefaultMatches(t, mc, c)
		})
	}
}

// TestReview_ParenthesizedLiteralDefault covers a literal default written in
// parentheses. It is still a literal, however many layers it carries, and is
// written in the column's own notation without them, where a real expression
// keeps its parentheses.
func TestReview_ParenthesizedLiteralDefault(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, ddl string }{
		{"string", "CREATE TABLE t (a VARCHAR(8) DEFAULT ('x'))"},
		{"adjacent_strings", "CREATE TABLE t (a VARCHAR(8) DEFAULT ('a' 'b'))"},
		{"integer", "CREATE TABLE t (a INT DEFAULT (7))"},
		{"negative", "CREATE TABLE t (a INT DEFAULT (-1))"},
		{"two_layers", "CREATE TABLE t (a INT DEFAULT ((7)))"},
		{"fraction", "CREATE TABLE t (a DOUBLE DEFAULT (1.5))"},
		{"exponent", "CREATE TABLE t (a DOUBLE DEFAULT (1e2))"},
		{"hex", "CREATE TABLE t (a BINARY(1) DEFAULT (X'FF'))"},
		{"null", "CREATE TABLE t (a INT DEFAULT (NULL))"},
		{"operator_keeps_parens", "CREATE TABLE t (a INT DEFAULT (1 + 1))"},
		{"call_is_bare", "CREATE TABLE t (a VARCHAR(36) DEFAULT (uuid()))"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, tc.ddl+";")
			assertDefaultMatches(t, mc, c)
		})
	}
}

// TestReview_ExplicitJSONValidCheckIsOrdinary covers a CHECK a statement
// declares that happens to look exactly like the one a JSON column carries.
// It is an ordinary constraint: it can be dropped by name, and its analyzed
// expression belongs to it.
//
// Measured on mariadb:11.8.8, which tells the two apart: the explicit one is
// droppable and the JSON column's answers 1091. So the catalog records how a
// constraint was made rather than reading it off the SQL.
func TestReview_ExplicitJSONValidCheckIsOrdinary(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct {
		name    string
		setup   string
		drop    string
		accepts bool
	}{
		{
			name:  "explicit_named_after_its_column",
			setup: "CREATE TABLE t (j LONGTEXT, CONSTRAINT j CHECK (json_valid(j)))",
			drop:  "ALTER TABLE t DROP CONSTRAINT j", accepts: true,
		},
		{
			name:  "explicit_under_another_name",
			setup: "CREATE TABLE t (j LONGTEXT, CONSTRAINT ck CHECK (json_valid(j)))",
			drop:  "ALTER TABLE t DROP CONSTRAINT ck", accepts: true,
		},
		{
			name:  "explicit_beside_a_json_column",
			setup: "CREATE TABLE t (j JSON, CONSTRAINT ck CHECK (json_valid(j)))",
			drop:  "ALTER TABLE t DROP CONSTRAINT ck", accepts: true,
		},
		{
			name:  "the_json_column_s_own",
			setup: "CREATE TABLE t (j JSON)",
			drop:  "ALTER TABLE t DROP CONSTRAINT j", accepts: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, tc.setup+";")

			_, oracleErr := mc.db.ExecContext(mc.ctx, tc.drop)
			if tc.accepts != (oracleErr == nil) {
				t.Fatalf("oracle: %q gave %v, want accepted=%v", tc.drop, oracleErr, tc.accepts)
			}

			results, err := c.Exec(tc.drop+";", nil)
			omniErr := err
			for _, r := range results {
				if r.Error != nil {
					omniErr = r.Error
				}
			}
			if tc.accepts != (omniErr == nil) {
				t.Errorf("omni: %q gave %v, want accepted=%v", tc.drop, omniErr, tc.accepts)
			}
		})
	}
}

// TestReview_ExplicitJSONValidCheckIsAnalyzed covers the other half of the
// same classification: an explicitly declared json_valid check is one of the
// statement's CHECK clauses, so it keeps its own analyzed expression and does
// not shift the ones after it.
func TestReview_ExplicitJSONValidCheckIsAnalyzed(t *testing.T) {
	for _, tc := range []struct {
		name     string
		ddl      string
		analyzed map[string]bool
	}{
		{
			name:     "explicit_named_after_its_column",
			ddl:      "CREATE TABLE t (j LONGTEXT, CONSTRAINT j CHECK (json_valid(j)))",
			analyzed: map[string]bool{"j": true},
		},
		{
			name: "explicit_then_another_check",
			ddl: "CREATE TABLE t (j LONGTEXT, a INT, CONSTRAINT j CHECK (json_valid(j))," +
				" CONSTRAINT positive CHECK (a > 0))",
			analyzed: map[string]bool{"j": true, "positive": true},
		},
		{
			name: "a_json_column_and_an_explicit_check",
			ddl: "CREATE TABLE t (j JSON, a INT, CONSTRAINT ck CHECK (json_valid(j))," +
				" CONSTRAINT positive CHECK (a > 0))",
			analyzed: map[string]bool{"j": false, "ck": true, "positive": true},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := New()
			results, err := c.Exec("CREATE DATABASE d; USE d; "+tc.ddl+";", nil)
			if err != nil {
				t.Fatal(err)
			}
			for _, r := range results {
				if r.Error != nil {
					t.Fatal(r.Error)
				}
			}
			seen := map[string]bool{}
			for _, con := range c.GetDatabase("d").GetTable("t").Constraints {
				if con.Type == ConCheck {
					seen[con.Name] = con.CheckAnalyzed != nil
				}
			}
			for name, want := range tc.analyzed {
				if got, ok := seen[name]; !ok || got != want {
					t.Errorf("constraint %q analyzed=%v (present=%v), want %v", name, got, ok, want)
				}
			}
			if len(seen) != len(tc.analyzed) {
				t.Errorf("CHECK constraints: got %v, want %v", seen, tc.analyzed)
			}
		})
	}
}

// TestReview_NumberDefaultIsCoerced covers the number a numeric column's
// default spells, against the text it was written as. The server holds the
// value, so whitespace, a sign, leading zeros, a bare point, exponent
// notation and a fraction the type cannot hold are all notation it drops.
// Writing the text back with only its quotes removed had left whitespace
// inside the number, so DECIMAL(5,2) DEFAULT '7 ' came out 7 .00, which does
// not parse.
func TestReview_NumberDefaultIsCoerced(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, ddl string }{
		// Whitespace around a quoted number.
		{"decimal_trailing_space", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT '7 ')"},
		{"decimal_leading_space", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT ' 7')"},
		{"decimal_both_sides", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT ' 7 ')"},
		{"decimal_tab_and_newline", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT '\t7\n')"},
		{"decimal_fraction_space", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT '7.5 ')"},
		{"int_trailing_space", "CREATE TABLE t (a INT DEFAULT '7 ')"},
		{"int_both_sides", "CREATE TABLE t (a INT DEFAULT ' 7 ')"},
		{"bigint_spaces_and_sign", "CREATE TABLE t (a BIGINT DEFAULT '  -12  ')"},
		{"float_spaces", "CREATE TABLE t (a FLOAT DEFAULT ' 3.14 ')"},
		{"double_exponent_spaces", "CREATE TABLE t (a DOUBLE DEFAULT ' 1e2 ')"},
		{"parenthesized", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT ('7 '))"},
		{"introduced", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT _utf8mb4'7 ')"},

		// A sign, leading zeros and a bare point.
		{"decimal_plus", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT '+7')"},
		{"int_plus", "CREATE TABLE t (a INT DEFAULT ' +7')"},
		{"decimal_leading_zeros", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT '0007.00')"},
		{"int_leading_zeros", "CREATE TABLE t (a INT DEFAULT '007')"},
		{"int_minus_zero", "CREATE TABLE t (a INT DEFAULT '-0')"},
		{"decimal_negative", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT '-7')"},
		{"decimal_trailing_point", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT '7.')"},
		{"decimal_leading_point", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT '  .7  ')"},
		{"double_leading_zeros", "CREATE TABLE t (a DOUBLE DEFAULT '007.5')"},
		{"float_plus", "CREATE TABLE t (a FLOAT DEFAULT '+3.5')"},

		// A fraction the type cannot hold, rounded away from zero.
		{"decimal_rounds_up", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT '7.125')"},
		{"decimal_rounds_unquoted", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT 7.125)"},
		{"decimal_rounds_negative", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT -7.125)"},
		{"decimal_rounds_long", "CREATE TABLE t (a DECIMAL(5,4) DEFAULT '7.12345')"},
		{"decimal_carries", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT '9.999')"},
		{"decimal_carries_whole", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT '7.999')"},
		{"decimal_rounds_to_zero", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT '-0.001')"},
		{"decimal_tiny_rounds_up", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT 0.005)"},
		{"decimal_no_scale_rounds", "CREATE TABLE t (a DECIMAL DEFAULT 7.9)"},
		{"decimal_zero_scale_rounds", "CREATE TABLE t (a DECIMAL(5,0) DEFAULT '7.9')"},
		{"int_rounds_up", "CREATE TABLE t (a INT DEFAULT '7.9')"},
		{"int_rounds_tie", "CREATE TABLE t (a INT DEFAULT 7.5)"},
		{"int_rounds_tie_negative", "CREATE TABLE t (a INT DEFAULT '-7.5')"},
		{"int_rounds_down", "CREATE TABLE t (a INT DEFAULT 7.4)"},
		{"int_rounds_to_zero", "CREATE TABLE t (a INT DEFAULT '-0.4')"},
		{"tinyint_rounds", "CREATE TABLE t (a TINYINT DEFAULT '1.7')"},
		{"boolean_rounds", "CREATE TABLE t (a BOOLEAN DEFAULT '1.7')"},

		// A float's fraction: a scaleless one prints as the double it holds,
		// where a declared scale rounds as a binary double rather than as
		// decimal digits.
		{"double_trailing_zeros", "CREATE TABLE t (a DOUBLE DEFAULT 7.000)"},
		{"float_trailing_zeros", "CREATE TABLE t (a FLOAT DEFAULT '7.000')"},
		{"double_fraction_zero", "CREATE TABLE t (a DOUBLE DEFAULT 7.10)"},
		{"double_scaled_pads", "CREATE TABLE t (a DOUBLE(10,2) DEFAULT 7)"},
		{"double_scaled_tie_down", "CREATE TABLE t (a DOUBLE(10,2) DEFAULT 0.125)"},
		{"double_scaled_tie_up", "CREATE TABLE t (a DOUBLE(10,2) DEFAULT 0.375)"},
		{"float_scaled_space", "CREATE TABLE t (a FLOAT(10,2) DEFAULT '7 ')"},

		// The aliases, which are the same type.
		{"numeric_alias", "CREATE TABLE t (a NUMERIC(5,2) DEFAULT '7 ')"},
		{"dec_alias", "CREATE TABLE t (a DEC(5,2) DEFAULT '7 ')"},
		{"fixed_alias", "CREATE TABLE t (a FIXED(5,2) DEFAULT '7 ')"},
		{"decimal_unsigned", "CREATE TABLE t (a DECIMAL(5,2) UNSIGNED DEFAULT '7 ')"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, tc.ddl+";")
			assertDefaultMatches(t, mc, c)
		})
	}
}

// TestReview_NonNumberDefaultStaysParseable covers a numeric column given a
// default that spells no number. The server refuses it outright with 1067 and
// omni does not check that yet, so what omni writes has to be a statement
// that still parses rather than bare text where a number belongs.
func TestReview_NonNumberDefaultStaysParseable(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, ddl string }{
		{"word", "CREATE TABLE t (a INT DEFAULT 'abc')"},
		{"trailing_word", "CREATE TABLE t (a INT DEFAULT '7abc')"},
		{"empty", "CREATE TABLE t (a INT DEFAULT '')"},
		{"only_space", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT ' ')"},
		{"bare_point", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT '.')"},
		{"bare_sign", "CREATE TABLE t (a INT DEFAULT '-')"},
		{"bare_exponent", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT '1e')"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			if _, err := mc.db.ExecContext(mc.ctx, tc.ddl); err == nil {
				t.Fatalf("oracle: expected %q to be refused", tc.ddl)
			}

			c := scenarioNewCatalog(t)
			results, err := c.Exec(tc.ddl+";", nil)
			if err != nil {
				t.Fatalf("omni parse error: %v", err)
			}
			for _, r := range results {
				if r.Error != nil {
					return // omni refused it too, which is the server's answer
				}
			}
			// Whatever omni wrote has to go back in.
			create := c.ShowCreateTable("testdb", "t")
			again := scenarioNewCatalog(t)
			results, err = again.Exec(create+";", nil)
			if err != nil {
				t.Fatalf("omni wrote a default it cannot read back: %v\n%s", err, create)
			}
			for _, r := range results {
				if r.Error != nil {
					t.Fatalf("omni wrote a default it cannot read back: %v\n%s", r.Error, create)
				}
			}
		})
	}
}

// TestReview_CheckBelongsToItsColumn covers where a CHECK is written. MariaDB
// writes a check declared inside a column definition, and the json_valid
// check a JSON column carries, on that column; one declared at table level it
// writes at table level, whatever either is named. Reading the placement off
// the name instead had a table-level check named after a column written on the
// column, and dropped one of two checks of the same name from the output
// altogether.
//
// All of it measured on mariadb:11.8.8 through SHOW CREATE TABLE.
//
// What a renamed column does to a check's own text is a separate matter, and
// not covered here: the server rewrites the references inside every stored
// expression, so after `CREATE TABLE t (j JSON, CONSTRAINT k CHECK
// (json_valid(j))); ALTER TABLE t RENAME COLUMN j TO k` the table-level check
// reads json_valid(`k`), where the catalog rewrites only the json_valid check
// it wrote itself.
func TestReview_CheckBelongsToItsColumn(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct {
		name  string
		stmts []string
	}{
		// A table-level check keeps its place even when it carries a
		// column's name.
		{"table_level_named_after_a_column", []string{
			"CREATE TABLE t (k INT, a INT, CONSTRAINT k CHECK (a > 0))",
		}},
		{"table_level_named_after_its_own_column", []string{
			"CREATE TABLE t (a INT, CONSTRAINT a CHECK (a > 0))",
		}},
		{"table_level_survives_a_column_rename", []string{
			"CREATE TABLE t (k INT, a INT, CONSTRAINT k CHECK (a > 0))",
			"ALTER TABLE t RENAME COLUMN k TO z",
		}},

		// A check declared on a column takes no name of its own, so it
		// collides with none: the two live side by side in either order,
		// where two checks the statement named answer 1826.
		{"column_check_then_a_table_check_of_that_name", []string{
			"CREATE TABLE t (a INT, b INT CHECK (b > 0), CONSTRAINT b CHECK (a > 0))",
		}},
		{"table_check_then_a_column_check_of_that_name", []string{
			"CREATE TABLE t (a INT, CONSTRAINT b CHECK (a > 0), b INT CHECK (b > 0))",
		}},
		{"table_check_named_after_another_columns_check", []string{
			"CREATE TABLE t (b INT CHECK (b > 0), c INT, CONSTRAINT b CHECK (c > 0))",
		}},

		// A JSON column's own check beside a table-level check of the same
		// name: the server carries both and writes both.
		{"change_a_json_column_onto_a_check_name", []string{
			"CREATE TABLE t (j JSON, a INT, CONSTRAINT k CHECK (a > 0))",
			"ALTER TABLE t CHANGE j k JSON",
		}},
		{"rename_a_json_column_onto_a_check_name", []string{
			"CREATE TABLE t (j JSON, a INT, CONSTRAINT k CHECK (a > 0))",
			"ALTER TABLE t RENAME COLUMN j TO k",
		}},
		{"declare_a_check_named_after_a_json_column", []string{
			"CREATE TABLE t (j JSON, a INT, CONSTRAINT j CHECK (a > 0))",
		}},
		{"add_a_check_named_after_a_json_column", []string{
			"CREATE TABLE t (j JSON, a INT)",
			"ALTER TABLE t ADD CONSTRAINT j CHECK (a > 0)",
		}},
		{"retype_to_json_where_the_name_is_taken", []string{
			"CREATE TABLE t (j LONGTEXT, a INT, CONSTRAINT j CHECK (a > 0))",
			"ALTER TABLE t MODIFY COLUMN j JSON",
		}},
		{"add_a_json_column_whose_name_is_taken", []string{
			"CREATE TABLE t (a INT, CONSTRAINT k CHECK (a > 0))",
			"ALTER TABLE t ADD COLUMN k JSON",
		}},

		// A JSON column that declares its own check carries that one
		// instead: the json_valid check is what a JSON column has in place
		// of one, not on top of it.
		{"json_column_with_its_own_check", []string{
			"CREATE TABLE t (j JSON CHECK (j IS NOT NULL))",
		}},
		{"json_column_with_its_own_check_and_a_table_one", []string{
			"CREATE TABLE t (j JSON CHECK (j IS NOT NULL), a INT, CONSTRAINT ck CHECK (a > 0))",
		}},

		// CHANGE recreates the column, so the check declared on it goes;
		// a JSON column gets a fresh json_valid check named after the new
		// column.
		{"change_drops_the_columns_check", []string{
			"CREATE TABLE t (a INT CHECK (a > 0))",
			"ALTER TABLE t CHANGE a z INT",
		}},
		{"change_resynthesizes_the_json_check", []string{
			"CREATE TABLE t (j JSON)",
			"ALTER TABLE t CHANGE j z JSON",
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			for _, stmt := range tc.stmts {
				runOnBoth(t, mc, c, stmt+";")
			}
			want := oracleShow(t, mc, "SHOW CREATE TABLE t")
			if got := c.ShowCreateTable("testdb", "t"); got != want {
				t.Errorf("table mismatch\n--- oracle ---\n%s\n--- omni ---\n%s", want, got)
			}
		})
	}
}

// TestReview_ColumnCheckIsNotDroppableByName covers the other side of where a
// check belongs: a check that belongs to a column is not a constraint a
// statement can name, so DROP CONSTRAINT does not find it and the server
// answers 1091, while a table-level check of the very same name is dropped.
//
// Measured on mariadb:11.8.8.
func TestReview_ColumnCheckIsNotDroppableByName(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct {
		name    string
		setup   string
		drop    string
		accepts bool
	}{
		{
			name:  "a_json_columns_own",
			setup: "CREATE TABLE t (j JSON)",
			drop:  "ALTER TABLE t DROP CONSTRAINT j", accepts: false,
		},
		{
			name:  "declared_on_its_column",
			setup: "CREATE TABLE t (a INT CHECK (a > 0))",
			drop:  "ALTER TABLE t DROP CONSTRAINT a", accepts: false,
		},
		{
			name:  "declared_at_table_level",
			setup: "CREATE TABLE t (a INT, CONSTRAINT ck CHECK (a > 0))",
			drop:  "ALTER TABLE t DROP CONSTRAINT ck", accepts: true,
		},
		{
			name:  "declared_at_table_level_under_a_columns_name",
			setup: "CREATE TABLE t (k INT, a INT, CONSTRAINT k CHECK (a > 0))",
			drop:  "ALTER TABLE t DROP CONSTRAINT k", accepts: true,
		},
		{
			name:  "declared_at_table_level_beside_a_json_columns_own",
			setup: "CREATE TABLE t (j JSON, a INT, CONSTRAINT j CHECK (a > 0))",
			drop:  "ALTER TABLE t DROP CONSTRAINT j", accepts: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, tc.setup+";")

			_, oracleErr := mc.db.ExecContext(mc.ctx, tc.drop)
			if tc.accepts != (oracleErr == nil) {
				t.Fatalf("oracle: %q gave %v, want accepted=%v", tc.drop, oracleErr, tc.accepts)
			}

			results, err := c.Exec(tc.drop+";", nil)
			omniErr := err
			for _, r := range results {
				if r.Error != nil {
					omniErr = r.Error
				}
			}
			if tc.accepts != (omniErr == nil) {
				t.Errorf("omni: %q gave %v, want accepted=%v", tc.drop, omniErr, tc.accepts)
			}
			if !tc.accepts {
				return
			}
			want := oracleShow(t, mc, "SHOW CREATE TABLE t")
			if got := c.ShowCreateTable("testdb", "t"); got != want {
				t.Errorf("table mismatch\n--- oracle ---\n%s\n--- omni ---\n%s", want, got)
			}
		})
	}
}

// TestReview_BitColumnDefault covers a BIT column's default. BIT holds a
// number and writes it as a bit literal of exactly its own digits, so every
// notation a default can be written in — hex, bit, a quoted string, an
// introduced one, a plain number — is read as that number and comes back as
// b'...'. The type was in no family at all, so a hex default came back wrapped
// in string quotes, with its own quotes still inside, which does not parse.
//
// All of it measured on mariadb:11.8.8 through SHOW CREATE TABLE.
func TestReview_BitColumnDefault(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, ddl string }{
		// The hex, bit and string notations, all of them the same value.
		{"hex_upper", "CREATE TABLE t (a BIT(8) DEFAULT X'41')"},
		{"hex_lower", "CREATE TABLE t (a BIT(8) DEFAULT x'41')"},
		{"hex_0x", "CREATE TABLE t (a BIT(8) DEFAULT 0x41)"},
		{"hex_leading_zeros", "CREATE TABLE t (a BIT(8) DEFAULT X'0041')"},
		{"bits", "CREATE TABLE t (a BIT(8) DEFAULT b'1000001')"},
		{"bits_upper", "CREATE TABLE t (a BIT(8) DEFAULT B'1000001')"},
		{"string_is_bytes", "CREATE TABLE t (a BIT(8) DEFAULT 'A')"},
		{"string_two_bytes", "CREATE TABLE t (a BIT(16) DEFAULT 'AB')"},
		{"introduced", "CREATE TABLE t (a BIT(8) DEFAULT _utf8mb4'A')"},
		{"parenthesized", "CREATE TABLE t (a BIT(8) DEFAULT (X'41'))"},

		// A plain number is the value itself.
		{"number", "CREATE TABLE t (a BIT(8) DEFAULT 65)"},
		{"number_parenthesized", "CREATE TABLE t (a BIT(8) DEFAULT (65))"},
		{"number_signed", "CREATE TABLE t (a BIT(8) DEFAULT +65)"},
		{"number_leading_zeros", "CREATE TABLE t (a BIT(8) DEFAULT 007)"},
		{"number_rounds_up", "CREATE TABLE t (a BIT(8) DEFAULT 1.7)"},
		{"number_rounds_down", "CREATE TABLE t (a BIT(8) DEFAULT 1.2)"},
		{"number_exponent", "CREATE TABLE t (a BIT(8) DEFAULT 1e2)"},
		{"number_max", "CREATE TABLE t (a BIT(8) DEFAULT 255)"},

		// Zero and the widths around it.
		{"zero", "CREATE TABLE t (a BIT(8) DEFAULT 0)"},
		{"empty_bits", "CREATE TABLE t (a BIT(8) DEFAULT b'')"},
		{"empty_hex", "CREATE TABLE t (a BIT(8) DEFAULT X'')"},
		{"one_bit", "CREATE TABLE t (a BIT(1) DEFAULT 1)"},
		{"one_bit_literal", "CREATE TABLE t (a BIT(1) DEFAULT b'1')"},
		{"no_width", "CREATE TABLE t (a BIT DEFAULT b'1')"},
		{"wider_column_does_not_pad", "CREATE TABLE t (a BIT(16) DEFAULT X'41')"},
		{"widest", "CREATE TABLE t (a BIT(64) DEFAULT X'FFFFFFFFFFFFFFFF')"},
		{"null", "CREATE TABLE t (a BIT(8) DEFAULT NULL)"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, tc.ddl+";")
			assertDefaultMatches(t, mc, c)
		})
	}
}

// TestReview_BitDefaultStaysParseable covers the BIT defaults the server
// refuses outright with 1067 — a value wider than the column, a negative one,
// a string whose bytes are too wide. omni does not check that yet, so what it
// writes has to be a statement that still parses.
func TestReview_BitDefaultStaysParseable(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, ddl string }{
		{"too_wide_hex", "CREATE TABLE t (a BIT(4) DEFAULT X'41')"},
		{"too_wide_number", "CREATE TABLE t (a BIT(8) DEFAULT 256)"},
		{"too_wide_bits", "CREATE TABLE t (a BIT(8) DEFAULT b'100000011')"},
		{"too_wide_string", "CREATE TABLE t (a BIT(8) DEFAULT 'AB')"},
		{"digits_in_a_string", "CREATE TABLE t (a BIT(8) DEFAULT '65')"},
		{"negative", "CREATE TABLE t (a BIT(8) DEFAULT -1)"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			if _, err := mc.db.ExecContext(mc.ctx, tc.ddl); err == nil {
				t.Fatalf("oracle: expected %q to be refused", tc.ddl)
			}

			c := scenarioNewCatalog(t)
			results, err := c.Exec(tc.ddl+";", nil)
			if err != nil {
				t.Fatalf("omni parse error: %v", err)
			}
			for _, r := range results {
				if r.Error != nil {
					return // omni refused it too, which is the server's answer
				}
			}
			create := c.ShowCreateTable("testdb", "t")
			again := scenarioNewCatalog(t)
			results, err = again.Exec(create+";", nil)
			if err != nil {
				t.Fatalf("omni wrote a default it cannot read back: %v\n%s", err, create)
			}
			for _, r := range results {
				if r.Error != nil {
					t.Fatalf("omni wrote a default it cannot read back: %v\n%s", r.Error, create)
				}
			}
		})
	}
}

// TestReview_IndexOnNondeterministicGeneratedColumn covers a key added over a
// VIRTUAL generated column after the fact. The key materializes the column's
// value, and the server then refuses a non-deterministic expression for it
// with 1901, exactly as it refuses the key declared with the table. Only the
// column's own definition was checked, so CREATE INDEX and ALTER ... ADD INDEX
// let it through.
//
// Measured on mariadb:11.8.8, including the message, which names the function
// as the server folded it: NOW() is reported as current_timestamp().
func TestReview_IndexOnNondeterministicGeneratedColumn(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct {
		name    string
		setup   string
		add     string
		accepts bool
	}{
		{
			name:  "create_index",
			setup: "CREATE TABLE t (a DOUBLE AS (RAND()) VIRTUAL)",
			add:   "CREATE INDEX idx ON t(a)",
		},
		{
			name:  "alter_add_index",
			setup: "CREATE TABLE t (a DOUBLE AS (RAND()) VIRTUAL)",
			add:   "ALTER TABLE t ADD INDEX idx(a)",
		},
		{
			name:  "alter_add_unique",
			setup: "CREATE TABLE t (a DOUBLE AS (RAND()) VIRTUAL)",
			add:   "ALTER TABLE t ADD UNIQUE idx(a)",
		},
		{
			name:  "covered_second",
			setup: "CREATE TABLE t (b INT, a DOUBLE AS (RAND()) VIRTUAL)",
			add:   "CREATE INDEX idx ON t(b,a)",
		},
		{
			name:  "now_is_reported_as_current_timestamp",
			setup: "CREATE TABLE t (a DATETIME AS (NOW()) VIRTUAL)",
			add:   "CREATE INDEX idx ON t(a)",
		},
		{
			name:  "uuid",
			setup: "CREATE TABLE t (a CHAR(36) AS (UUID()) VIRTUAL)",
			add:   "CREATE INDEX idx ON t(a)",
		},
		{
			name:  "text_column_before_the_key_length",
			setup: "CREATE TABLE t (a TEXT AS (UUID()) VIRTUAL)",
			add:   "CREATE INDEX idx ON t(a)",
		},
		{
			name:    "deterministic_is_fine",
			setup:   "CREATE TABLE t (b INT, a INT AS (b + 1) VIRTUAL)",
			add:     "CREATE INDEX idx ON t(a)",
			accepts: true,
		},
		{
			name:    "an_ordinary_column_is_fine",
			setup:   "CREATE TABLE t (a INT)",
			add:     "CREATE INDEX idx ON t(a)",
			accepts: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, tc.setup+";")

			_, oracleErr := mc.db.ExecContext(mc.ctx, tc.add)
			if tc.accepts != (oracleErr == nil) {
				t.Fatalf("oracle: %q gave %v, want accepted=%v", tc.add, oracleErr, tc.accepts)
			}

			results, err := c.Exec(tc.add+";", nil)
			omniErr := err
			for _, r := range results {
				if r.Error != nil {
					omniErr = r.Error
				}
			}
			if tc.accepts != (omniErr == nil) {
				t.Fatalf("omni: %q gave %v, want accepted=%v", tc.add, omniErr, tc.accepts)
			}
			if tc.accepts {
				want := oracleShow(t, mc, "SHOW CREATE TABLE t")
				if got := c.ShowCreateTable("testdb", "t"); got != want {
					t.Errorf("table mismatch\n--- oracle ---\n%s\n--- omni ---\n%s", want, got)
				}
				return
			}
			if want := messageTailOf(oracleErr); !strings.Contains(omniErr.Error(), want) {
				t.Errorf("message mismatch\n  oracle: %s\n  omni:   %s", want, omniErr)
			}
		})
	}
}

// messageTailOf strips the leading code and SQL state from a server error, so
// that the message itself can be compared with the catalog's.
func messageTailOf(err error) string {
	msg := err.Error()
	if i := strings.Index(msg, "): "); i >= 0 {
		return msg[i+3:]
	}
	return msg
}

// TestReview_YearColumnDefault covers a YEAR column's default, whose value is
// more than its digits: the number is expanded to the four-digit year it
// denotes, and zero means one year or another depending on how it was written.
// Coercing the literal into plain digits had turned DEFAULT '0', which is the
// year 2000, into DEFAULT 0, which is the zero year.
//
// All of it measured on mariadb:11.8.8 through SHOW CREATE TABLE, with the
// stored value read back from a row to confirm the year each one means.
func TestReview_YearColumnDefault(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, ddl string }{
		// Zero: the year zero where it was written as a number, 2000 where it
		// was written as a string — except as a four-character one.
		{"number_zero", "CREATE TABLE t (a YEAR DEFAULT 0)"},
		{"number_zero_padded", "CREATE TABLE t (a YEAR DEFAULT 0000)"},
		{"number_zero_long", "CREATE TABLE t (a YEAR DEFAULT 00000)"},
		{"number_zero_signed", "CREATE TABLE t (a YEAR DEFAULT +0)"},
		{"number_zero_fraction", "CREATE TABLE t (a YEAR DEFAULT 0.0)"},
		{"number_zero_parenthesized", "CREATE TABLE t (a YEAR DEFAULT (0))"},
		{"string_zero", "CREATE TABLE t (a YEAR DEFAULT '0')"},
		{"string_zero_two", "CREATE TABLE t (a YEAR DEFAULT '00')"},
		{"string_zero_three", "CREATE TABLE t (a YEAR DEFAULT '000')"},
		{"string_zero_four", "CREATE TABLE t (a YEAR DEFAULT '0000')"},
		{"string_zero_five", "CREATE TABLE t (a YEAR DEFAULT '00000')"},
		{"string_zero_four_with_sign", "CREATE TABLE t (a YEAR DEFAULT '+000')"},
		{"string_zero_spaced", "CREATE TABLE t (a YEAR DEFAULT ' 0000 ')"},
		{"string_zero_negative", "CREATE TABLE t (a YEAR DEFAULT '-0')"},
		{"string_zero_fraction", "CREATE TABLE t (a YEAR DEFAULT '0.0')"},
		{"string_zero_parenthesized", "CREATE TABLE t (a YEAR DEFAULT ('0'))"},
		{"introduced_zero", "CREATE TABLE t (a YEAR DEFAULT _utf8mb4'0')"},

		// The two-digit years, quoted or not.
		{"one", "CREATE TABLE t (a YEAR DEFAULT 1)"},
		{"string_one", "CREATE TABLE t (a YEAR DEFAULT '1')"},
		{"string_one_padded", "CREATE TABLE t (a YEAR DEFAULT '01')"},
		{"one_fraction", "CREATE TABLE t (a YEAR DEFAULT 1.0)"},
		{"sixty_nine", "CREATE TABLE t (a YEAR DEFAULT 69)"},
		{"seventy", "CREATE TABLE t (a YEAR DEFAULT 70)"},
		{"string_seventy", "CREATE TABLE t (a YEAR DEFAULT '70')"},
		{"ninety_nine", "CREATE TABLE t (a YEAR DEFAULT 99)"},
		{"string_ninety_nine", "CREATE TABLE t (a YEAR DEFAULT '99')"},
		{"string_ninety_nine_spaced", "CREATE TABLE t (a YEAR DEFAULT ' 99 ')"},
		{"introduced_ninety_nine", "CREATE TABLE t (a YEAR DEFAULT _utf8mb4'99')"},

		// The four-digit years.
		{"first_year", "CREATE TABLE t (a YEAR DEFAULT 1901)"},
		{"string_first_year", "CREATE TABLE t (a YEAR DEFAULT '1901')"},
		{"string_first_year_padded", "CREATE TABLE t (a YEAR DEFAULT '01901')"},
		{"last_year", "CREATE TABLE t (a YEAR DEFAULT 2155)"},
		{"two_thousand", "CREATE TABLE t (a YEAR DEFAULT 2000)"},
		{"string_two_thousand", "CREATE TABLE t (a YEAR DEFAULT '2000')"},

		// Hex and bit notations are the number their bytes spell.
		{"hex_0x", "CREATE TABLE t (a YEAR DEFAULT 0x41)"},
		{"hex_0x_zero", "CREATE TABLE t (a YEAR DEFAULT 0x0)"},
		{"hex_0x_four_digit", "CREATE TABLE t (a YEAR DEFAULT 0x7E3)"},
		{"bits", "CREATE TABLE t (a YEAR DEFAULT b'1100011')"},
		{"bits_zero", "CREATE TABLE t (a YEAR DEFAULT b'0')"},

		{"null", "CREATE TABLE t (a YEAR DEFAULT NULL)"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, tc.ddl+";")
			assertDefaultMatches(t, mc, c)

			// What omni wrote has to mean the same year to the server, which
			// is what the quoting carried and plain digits would have lost.
			if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
				t.Fatalf("oracle insert: %v", err)
			}
			var stored any
			oracleScan(t, mc, "SELECT a FROM t", &stored)
			if _, err := mc.db.ExecContext(mc.ctx, "DROP TABLE t"); err != nil {
				t.Fatalf("oracle drop: %v", err)
			}
			replay := c.ShowCreateTable("testdb", "t")
			if _, err := mc.db.ExecContext(mc.ctx, replay); err != nil {
				t.Fatalf("oracle cannot read back what omni wrote: %v\n%s", err, replay)
			}
			if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
				t.Fatalf("oracle insert after replay: %v", err)
			}
			var replayed any
			oracleScan(t, mc, "SELECT a FROM t", &replayed)
			if fmt.Sprint(stored) != fmt.Sprint(replayed) {
				t.Errorf("replaying what omni wrote changed the year: %v became %v\n%s",
					stored, replayed, replay)
			}
		})
	}
}

// TestReview_BooleanLiteralDefault covers a hex or bit literal default on a
// BOOLEAN column, which is a TINYINT and reads the literal as a number.
// BOOLEAN was in no family, so the literal came back as written.
func TestReview_BooleanLiteralDefault(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, ddl string }{
		{"hex_0x", "CREATE TABLE t (a BOOLEAN DEFAULT 0x41)"},
		{"bits", "CREATE TABLE t (a BOOLEAN DEFAULT b'1000001')"},
		{"bool_alias", "CREATE TABLE t (a BOOL DEFAULT 0x41)"},
		{"number", "CREATE TABLE t (a BOOLEAN DEFAULT 1)"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, tc.ddl+";")
			assertDefaultMatches(t, mc, c)
		})
	}
}

// TestReview_DroppedIndexKeepsUnrelatedChecks covers the constraint cleanup a
// dropped column sets off. An index goes when the column it covered goes, and
// the constraints that index held up go with it — but a CHECK is held up by no
// index, so one that merely carries the index's name stays, enforcement and
// all.
//
// Measured on mariadb:11.8.8.
func TestReview_DroppedIndexKeepsUnrelatedChecks(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct {
		name  string
		stmts []string
	}{
		{"json_check_named_like_the_index", []string{
			"CREATE TABLE t (a INT, j JSON, KEY j(a))",
			"ALTER TABLE t DROP COLUMN a",
		}},
		{"table_check_named_like_the_index", []string{
			"CREATE TABLE t (a INT, b INT, CONSTRAINT ck CHECK (b > 0), KEY ck(a))",
			"ALTER TABLE t DROP COLUMN a",
		}},
		{"column_check_named_like_the_index", []string{
			"CREATE TABLE t (a INT, b INT CHECK (b > 0), KEY b(a))",
			"ALTER TABLE t DROP COLUMN a",
		}},
		{"unique_constraint_is_kept", []string{
			"CREATE TABLE t (a INT, b INT, UNIQUE u(b), KEY u2(a))",
			"ALTER TABLE t DROP COLUMN a",
		}},
		// The case the name match was there for: an index-backed constraint
		// over the dropped column goes with its index.
		{"unique_over_the_dropped_column_goes", []string{
			"CREATE TABLE t (a INT, b INT, UNIQUE u(a))",
			"ALTER TABLE t DROP COLUMN a",
		}},
		{"primary_key_over_the_dropped_column_goes", []string{
			"CREATE TABLE t (a INT, b INT, PRIMARY KEY (a))",
			"ALTER TABLE t DROP COLUMN a",
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			for _, stmt := range tc.stmts {
				runOnBoth(t, mc, c, stmt+";")
			}
			want := oracleShow(t, mc, "SHOW CREATE TABLE t")
			if got := c.ShowCreateTable("testdb", "t"); got != want {
				t.Errorf("table mismatch\n--- oracle ---\n%s\n--- omni ---\n%s", want, got)
			}
		})
	}
}

// TestReview_ReplacedGeneratedColumnKeepsItsIndex covers MODIFY and CHANGE
// over a generated column a key already covers. The key materializes the
// column's value, so the new definition is held to the same rule as one
// declared with a key — and under CHANGE the question is about the column
// being replaced, not about the name the replacement carries.
//
// Measured on mariadb:11.8.8.
func TestReview_ReplacedGeneratedColumnKeepsItsIndex(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct {
		name    string
		setup   string
		alter   string
		accepts bool
	}{
		{
			name:  "change_to_nondeterministic",
			setup: "CREATE TABLE t (a DOUBLE AS (1) VIRTUAL, KEY (a))",
			alter: "ALTER TABLE t CHANGE a b DOUBLE AS (RAND()) VIRTUAL",
		},
		{
			name:  "modify_to_nondeterministic",
			setup: "CREATE TABLE t (a DOUBLE AS (1) VIRTUAL, KEY (a))",
			alter: "ALTER TABLE t MODIFY a DOUBLE AS (RAND()) VIRTUAL",
		},
		{
			name:  "change_to_nondeterministic_under_a_second_key_part",
			setup: "CREATE TABLE t (x INT, a DOUBLE AS (1) VIRTUAL, KEY (x, a))",
			alter: "ALTER TABLE t CHANGE a b DOUBLE AS (RAND()) VIRTUAL",
		},
		{
			name:    "change_to_nondeterministic_without_a_key",
			setup:   "CREATE TABLE t (a DOUBLE AS (1) VIRTUAL)",
			alter:   "ALTER TABLE t CHANGE a b DOUBLE AS (RAND()) VIRTUAL",
			accepts: true,
		},
		{
			name:    "change_staying_deterministic",
			setup:   "CREATE TABLE t (a DOUBLE AS (1) VIRTUAL, KEY (a))",
			alter:   "ALTER TABLE t CHANGE a b DOUBLE AS (2) VIRTUAL",
			accepts: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, tc.setup+";")

			_, oracleErr := mc.db.ExecContext(mc.ctx, tc.alter)
			if tc.accepts != (oracleErr == nil) {
				t.Fatalf("oracle: %q gave %v, want accepted=%v", tc.alter, oracleErr, tc.accepts)
			}

			results, err := c.Exec(tc.alter+";", nil)
			omniErr := err
			for _, r := range results {
				if r.Error != nil {
					omniErr = r.Error
				}
			}
			if tc.accepts != (omniErr == nil) {
				t.Fatalf("omni: %q gave %v, want accepted=%v", tc.alter, omniErr, tc.accepts)
			}
			if tc.accepts {
				want := oracleShow(t, mc, "SHOW CREATE TABLE t")
				if got := c.ShowCreateTable("testdb", "t"); got != want {
					t.Errorf("table mismatch\n--- oracle ---\n%s\n--- omni ---\n%s", want, got)
				}
				return
			}
			if want := messageTailOf(oracleErr); !strings.Contains(omniErr.Error(), want) {
				t.Errorf("message mismatch\n  oracle: %s\n  omni:   %s", want, omniErr)
			}
		})
	}
}

// TestReview_TemporalLiteralDefault covers a hex or introduced literal default
// on a date or time column. A temporal value is held as the text of a quoted
// string, so the literal's bytes are read as that text: the types were in no
// family, so a hex default came back wrapped in string quotes with its own
// quotes still inside, which does not parse.
//
// Measured on mariadb:11.8.8. The server coerces further than the catalog
// does, and those normalizations are left out here: a date is written in its
// canonical spelling, a DATETIME given a date is padded with a time, a
// fraction is scaled to the column's precision, and a number is read as a
// date. See familyTemporal.
func TestReview_TemporalLiteralDefault(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, ddl string }{
		{"date_hex", "CREATE TABLE t (a DATE DEFAULT X'323032302d30312d3031')"},
		{"date_hex_0x", "CREATE TABLE t (a DATE DEFAULT 0x323032302d30312d3031)"},
		{"date_string", "CREATE TABLE t (a DATE DEFAULT '2020-01-01')"},
		{"date_introduced", "CREATE TABLE t (a DATE DEFAULT _utf8mb4'2020-01-01')"},
		{"date_parenthesized", "CREATE TABLE t (a DATE DEFAULT ('2020-01-01'))"},
		{"datetime_hex", "CREATE TABLE t (a DATETIME DEFAULT X'323032302d30312d30312030303a30303a3030')"},
		{"datetime_string", "CREATE TABLE t (a DATETIME DEFAULT '2020-01-01 00:00:00')"},
		{"time_hex", "CREATE TABLE t (a TIME DEFAULT X'31323a33343a3536')"},
		{"time_string", "CREATE TABLE t (a TIME DEFAULT '12:34:56')"},
		{"timestamp_hex", "CREATE TABLE t (a TIMESTAMP DEFAULT X'323032302d30312d30312030303a30303a3030')"},
		{"timestamp_string", "CREATE TABLE t (a TIMESTAMP DEFAULT '2020-01-01 00:00:00')"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, tc.ddl+";")
			assertDefaultMatches(t, mc, c)

			// What omni wrote has to go back in, which is what the quoted hex
			// literal did not.
			if _, err := mc.db.ExecContext(mc.ctx, "DROP TABLE t"); err != nil {
				t.Fatalf("oracle drop: %v", err)
			}
			replay := c.ShowCreateTable("testdb", "t")
			if _, err := mc.db.ExecContext(mc.ctx, replay); err != nil {
				t.Fatalf("oracle cannot read back what omni wrote: %v\n%s", err, replay)
			}
		})
	}
}

// TestReview_AddColumnCheck covers a CHECK declared on a column an ALTER adds.
// It belongs to that column and takes its name, as one declared with the table
// does — and a JSON column that declares one carries it in place of
// json_valid. The ADD COLUMN path handled no column-level CHECK at all, so the
// check was dropped and a JSON column got json_valid instead of it.
//
// Measured on mariadb:11.8.8.
func TestReview_AddColumnCheck(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct {
		name  string
		stmts []string
	}{
		{"json_column_with_its_own_check", []string{
			"CREATE TABLE t (a INT)",
			"ALTER TABLE t ADD COLUMN j JSON CHECK (j IS NOT NULL)",
		}},
		{"json_column_without_one", []string{
			"CREATE TABLE t (a INT)",
			"ALTER TABLE t ADD COLUMN j JSON",
		}},
		{"plain_column", []string{
			"CREATE TABLE t (a INT)",
			"ALTER TABLE t ADD COLUMN b INT CHECK (b > 0)",
		}},
		{"check_naming_another_column", []string{
			"CREATE TABLE t (a INT)",
			"ALTER TABLE t ADD COLUMN b INT CHECK (a > 0)",
		}},
		{"the_name_is_already_a_table_checks", []string{
			"CREATE TABLE t (a INT, CONSTRAINT b CHECK (a > 0))",
			"ALTER TABLE t ADD COLUMN b INT CHECK (b > 0)",
		}},
		{"two_columns_at_once", []string{
			"CREATE TABLE t (a INT)",
			"ALTER TABLE t ADD COLUMN b INT CHECK (b > 0), ADD COLUMN c INT CHECK (c > 0)",
		}},
		{"dropping_the_column_takes_its_check", []string{
			"CREATE TABLE t (a INT)",
			"ALTER TABLE t ADD COLUMN b INT CHECK (b > 0)",
			"ALTER TABLE t DROP COLUMN b",
		}},
		{"the_check_is_not_droppable_by_name", []string{
			"CREATE TABLE t (a INT)",
			"ALTER TABLE t ADD COLUMN b INT CHECK (b > 0)",
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			for _, stmt := range tc.stmts {
				runOnBoth(t, mc, c, stmt+";")
			}
			want := oracleShow(t, mc, "SHOW CREATE TABLE t")
			if got := c.ShowCreateTable("testdb", "t"); got != want {
				t.Errorf("table mismatch\n--- oracle ---\n%s\n--- omni ---\n%s", want, got)
			}
			if tc.name != "the_check_is_not_droppable_by_name" {
				return
			}
			// The added check belongs to its column, so it is no more
			// droppable by name than one declared with the table.
			const drop = "ALTER TABLE t DROP CONSTRAINT b"
			if _, err := mc.db.ExecContext(mc.ctx, drop); err == nil {
				t.Fatalf("oracle: expected %q to be refused", drop)
			}
			results, err := c.Exec(drop+";", nil)
			omniErr := err
			for _, r := range results {
				if r.Error != nil {
					omniErr = r.Error
				}
			}
			if omniErr == nil {
				t.Errorf("omni: accepted %q, which the server refuses", drop)
			}
		})
	}
}

// TestReview_CTENameIsWrittenBare is the finding that was not a bug, with the
// measurement that settles it. MariaDB writes a CTE's declared name bare in
// the stored view body whatever it is called — a name with a space, a reserved
// word, a backtick inside — while the references to it keep their backticks.
// So the server's own stored text does not parse, and the view it just created
// is unusable: measured on mariadb:11.8.8, a view selecting from a CTE named
// "my cte" is created without complaint, and then both SELECT * FROM it and
// SHOW CREATE VIEW of it answer 1064, while information_schema.VIEWS reports
// the body as "with my cte as (...)".
//
// Quoting the name in the catalog would make omni's body differ from the one
// the server holds, and would make omni's view work where the server's is
// dead. The engine's output is reproduced instead, and this is what pins it.
//
// The queries here keep to one FROM item and no subquery, because two
// differences of their own sit there and are not this test's to cover: the
// server writes a comma join without the parentheses it puts around an
// explicit JOIN — "from t1 join t2" against "from (t1 join t2)", both with
// their names backticked — where the catalog parenthesizes both, and inside a
// scalar subquery it qualifies a column the catalog writes bare.
func TestReview_CTENameIsWrittenBare(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)

	for i, tc := range []struct {
		name     string
		query    string
		readable bool
	}{
		{"plain", "WITH cte AS (SELECT 1 AS x) SELECT x FROM cte", true},
		{"mixed_case", "WITH MyCte AS (SELECT 1 AS x) SELECT x FROM MyCte", true},
		{"quoted_but_plain", "WITH `MyCte` AS (SELECT 1 AS x) SELECT x FROM `MyCte`", true},
		{"unreserved_keyword", "WITH `comment` AS (SELECT 1 AS x) SELECT x FROM `comment`", true},
		{"with_space", "WITH `my cte` AS (SELECT 1 AS x) SELECT x FROM `my cte`", false},
		{"reserved_word", "WITH `select` AS (SELECT 1 AS x) SELECT x FROM `select`", false},
		// A name the lexer still takes bare: `1cte` needs its backticks in
		// the statement but reads back without them.
		{"leading_digit", "WITH `1cte` AS (SELECT 1 AS x) SELECT x FROM `1cte`", true},
		{"with_dot", "WITH `a.b` AS (SELECT 1 AS x) SELECT x FROM `a.b`", false},
		{"with_dash", "WITH `a-b` AS (SELECT 1 AS x) SELECT x FROM `a-b`", false},
		{"with_backtick", "WITH `a``b` AS (SELECT 1 AS x) SELECT x FROM `a``b`", false},
		{"column_list", "WITH `my cte`(`x`) AS (SELECT 1) SELECT x FROM `my cte`", false},
		{"two_ctes", "WITH `my cte` AS (SELECT 1 AS x), plain AS (SELECT 2 AS y)" +
			" SELECT x FROM `my cte`", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			view := fmt.Sprintf("cv%d", i)
			runOnBoth(t, mc, c, "CREATE VIEW "+view+" AS "+tc.query+";")

			// The body the server holds, read from information_schema, since
			// SHOW CREATE VIEW cannot print the broken ones at all.
			want := oracleViewDefinition(t, mc, view)
			if got := extractSelectBody(c.ShowCreateView("testdb", view)); got != want {
				t.Errorf("view body mismatch\n  oracle: %s\n  omni:   %s", want, got)
			}

			_, selectErr := mc.db.ExecContext(mc.ctx, "SELECT * FROM "+view)
			_, showErr := mc.db.ExecContext(mc.ctx, "SHOW CREATE VIEW "+view)
			if tc.readable {
				if selectErr != nil || showErr != nil {
					t.Errorf("oracle: expected the view to be usable, got %v / %v",
						selectErr, showErr)
				}
				return
			}
			// The name needs quoting, so what the server stored does not
			// parse and the view it created cannot be read or printed.
			if selectErr == nil {
				t.Errorf("oracle: expected SELECT * FROM %s to fail on the stored body", view)
			}
			if showErr == nil {
				t.Errorf("oracle: expected SHOW CREATE VIEW %s to fail on the stored body", view)
			}
		})
	}
}

// TestReview_RenamedColumnKeepsItsCheck covers what a renamed column does to
// the check it carries. The check belongs to the column, so the ownership
// moves with it — and ownership left on the old name would make the check the
// next column of that name's, then go with that one when it is dropped, while
// the column it really belongs to still stands.
//
// Measured on mariadb:11.8.8.
//
// Renaming the column a check *references* is left out: the server rewrites
// the references inside every stored expression, turning CHECK (`b` > 0) into
// CHECK (`y` > 0), which is the general rewrite the catalog does not do. The
// ownership is right there either way, since the check belongs to neither
// renamed column.
func TestReview_RenamedColumnKeepsItsCheck(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct {
		name  string
		stmts []string
	}{
		// The whole sequence: rename the column, let another take its old
		// name, then drop that one. The check stays where it belongs.
		{"the_old_name_is_taken_then_dropped", []string{
			"CREATE TABLE t (a INT CHECK (b > 0), b INT)",
			"ALTER TABLE t RENAME COLUMN a TO z",
			"ALTER TABLE t ADD COLUMN a INT",
			"ALTER TABLE t DROP COLUMN a",
		}},
		{"the_old_name_is_taken", []string{
			"CREATE TABLE t (a INT CHECK (b > 0), b INT)",
			"ALTER TABLE t RENAME COLUMN a TO z",
			"ALTER TABLE t ADD COLUMN a INT",
		}},
		{"a_check_naming_another_column", []string{
			"CREATE TABLE t (a INT CHECK (b > 0), b INT)",
			"ALTER TABLE t RENAME COLUMN a TO z",
		}},
		{"dropping_the_renamed_column_takes_it", []string{
			"CREATE TABLE t (a INT CHECK (b > 0), b INT)",
			"ALTER TABLE t RENAME COLUMN a TO z",
			"ALTER TABLE t DROP COLUMN z",
		}},

		// The same for the check a JSON column carries.
		{"a_json_columns_old_name_is_taken_then_dropped", []string{
			"CREATE TABLE t (j JSON)",
			"ALTER TABLE t RENAME COLUMN j TO z",
			"ALTER TABLE t ADD COLUMN j INT",
			"ALTER TABLE t DROP COLUMN j",
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			for _, stmt := range tc.stmts {
				runOnBoth(t, mc, c, stmt+";")
			}
			want := oracleShow(t, mc, "SHOW CREATE TABLE t")
			if got := c.ShowCreateTable("testdb", "t"); got != want {
				t.Errorf("table mismatch\n--- oracle ---\n%s\n--- omni ---\n%s", want, got)
			}
		})
	}
}

// TestReview_RenamedColumnKeepsItsCheckName covers the name of a check whose
// column is renamed, which SHOW CREATE TABLE does not show: the check is
// written on its column, without one. The server leaves the name as it was,
// measured on mariadb:11.8.8 through information_schema.CHECK_CONSTRAINTS,
// where a check declared on `a` is still called a after the column becomes z
// and again after it becomes w.
func TestReview_RenamedColumnKeepsItsCheckName(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct {
		name  string
		stmts []string
	}{
		{"declared_on_its_column", []string{
			"CREATE TABLE t (a INT CHECK (a > 0))",
			"ALTER TABLE t RENAME COLUMN a TO z",
		}},
		{"renamed_twice", []string{
			"CREATE TABLE t (a INT CHECK (a > 0))",
			"ALTER TABLE t RENAME COLUMN a TO z",
			"ALTER TABLE t RENAME COLUMN z TO w",
		}},
		{"a_json_columns_own", []string{
			"CREATE TABLE t (j JSON)",
			"ALTER TABLE t RENAME COLUMN j TO z",
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			for _, stmt := range tc.stmts {
				runOnBoth(t, mc, c, stmt+";")
			}

			var want string
			oracleScan(t, mc,
				"SELECT CONSTRAINT_NAME FROM information_schema.CHECK_CONSTRAINTS"+
					" WHERE CONSTRAINT_SCHEMA = DATABASE() AND TABLE_NAME = 't'", &want)

			var got []string
			for _, con := range c.GetDatabase("testdb").GetTable("t").Constraints {
				if con.Type == ConCheck {
					got = append(got, con.Name)
				}
			}
			if len(got) != 1 || got[0] != want {
				t.Errorf("check names: omni %v, oracle [%s]", got, want)
			}
		})
	}
}

// TestReview_ReturnTypeCollationCharset covers a RETURNS clause that names a
// collation and no charset. The charset comes from that collation, not from
// the database: pairing the named collation with the database's charset wrote
// `CHARSET utf8mb4 COLLATE latin1_bin`, a combination the server refuses with
// 1253, so the function could not be replayed.
//
// Measured on mariadb:11.8.8 through SHOW CREATE FUNCTION, in a utf8mb4
// database.
func TestReview_ReturnTypeCollationCharset(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for i, tc := range []struct{ name, returns string }{
		{"collation_of_another_charset", "RETURNS VARCHAR(20) COLLATE latin1_bin"},
		{"collation_of_the_database_charset", "RETURNS VARCHAR(20) COLLATE utf8mb4_bin"},
		{"collation_of_utf8mb3", "RETURNS VARCHAR(20) COLLATE utf8mb3_bin"},
		{"the_charsets_default_collation", "RETURNS VARCHAR(20) COLLATE latin1_swedish_ci"},
		{"charset_only", "RETURNS VARCHAR(20) CHARACTER SET latin1"},
		{"charset_and_its_collation", "RETURNS VARCHAR(20) CHARACTER SET latin1 COLLATE latin1_bin"},
		{"neither", "RETURNS VARCHAR(20)"},
		{"text", "RETURNS TEXT COLLATE latin1_bin"},
		{"char", "RETURNS CHAR(5) COLLATE latin1_bin"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			fn := fmt.Sprintf("rf%d", i)
			ddl := fmt.Sprintf("CREATE FUNCTION %s() %s DETERMINISTIC RETURN 'a'", fn, tc.returns)
			runOnBoth(t, mc, c, ddl+";")

			want, err := mc.showCreateFunction(fn)
			if err != nil {
				t.Fatalf("oracle: %v", err)
			}
			got := c.ShowCreateFunction("testdb", fn)
			if got != want {
				t.Errorf("function DDL mismatch\n--- oracle ---\n%s\n--- omni ---\n%s", want, got)
			}

			// What omni wrote has to go back in: the charset and collation it
			// pairs have to be a pair the server accepts.
			if _, err := mc.db.ExecContext(mc.ctx, "DROP FUNCTION "+fn); err != nil {
				t.Fatalf("oracle drop: %v", err)
			}
			if _, err := mc.db.ExecContext(mc.ctx, got); err != nil {
				t.Fatalf("oracle cannot read back what omni wrote: %v\n%s", err, got)
			}
		})
	}
}

// TestReview_UnionOrderByKeepsTheOldAlias is the finding that was not a bug.
// A view that renames a UNION's output through its own column list keeps the
// global ORDER BY as it was written, naming an alias the body no longer has.
// The server's own stored text is what does that, and the view it just created
// is unusable: measured on mariadb:11.8.8, after
// `CREATE VIEW v (x) AS SELECT a AS z FROM t UNION SELECT b FROM t ORDER BY z`
// the stored body reads `select ... AS x ... union ... order by z`, selecting
// from the view answers 1356 and replaying that body answers 1054, Unknown
// column 'z' in 'ORDER BY'.
//
// Rewriting the reference would make omni's body differ from the one the
// server holds and would make omni's view work where the server's is dead.
// The engine's output is reproduced instead. Writing the ORDER BY against the
// new name is not an option either: the server refuses `ORDER BY x` outright
// at creation, also with 1054.
func TestReview_UnionOrderByKeepsTheOldAlias(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)
	runOnBoth(t, mc, c, "CREATE TABLE t (a INT, b INT);")

	for i, tc := range []struct {
		name     string
		cols     string
		body     string
		readable bool
	}{
		{"union_order_by_alias", "(x)",
			"SELECT a AS z FROM t UNION SELECT b FROM t ORDER BY z", false},
		{"union_all", "(x)",
			"SELECT a AS z FROM t UNION ALL SELECT b FROM t ORDER BY z", false},
		{"order_desc", "(x)",
			"SELECT a AS z FROM t UNION SELECT b FROM t ORDER BY z DESC", false},
		{"order_expression", "(x)",
			"SELECT a AS z FROM t UNION SELECT b FROM t ORDER BY z + 1", false},
		{"two_columns", "(x, y)",
			"SELECT a AS z, b AS w FROM t UNION SELECT b, a FROM t ORDER BY w, z", false},
		// An ordinal survives the rename, and so does a view with no column
		// list of its own.
		{"order_by_ordinal", "(x)",
			"SELECT a AS z FROM t UNION SELECT b FROM t ORDER BY 1", true},
		{"no_column_list", "",
			"SELECT a AS z FROM t UNION SELECT b FROM t ORDER BY z", true},
		// Only a set operation keeps the alias: a plain SELECT has its ORDER
		// BY rewritten to the column underneath.
		{"plain_select", "(x)", "SELECT a AS z FROM t ORDER BY z", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			view := fmt.Sprintf("uv%d", i)
			runOnBoth(t, mc, c, "CREATE VIEW "+view+" "+tc.cols+" AS "+tc.body+";")

			// SHOW CREATE VIEW still prints these: the stored text parses,
			// it just names a column the body no longer has.
			assertViewBodyMatches(t, mc, c, view)

			_, selectErr := mc.db.ExecContext(mc.ctx, "SELECT * FROM "+view)
			if tc.readable {
				if selectErr != nil {
					t.Errorf("oracle: expected the view to be usable, got %v", selectErr)
				}
				return
			}
			if selectErr == nil {
				t.Errorf("oracle: expected SELECT * FROM %s to fail on the stored body", view)
			}
		})
	}

	// The ORDER BY cannot be written against the new name instead: the server
	// refuses that at creation.
	const refused = "CREATE VIEW uv_new (x) AS SELECT a AS z FROM t UNION SELECT b FROM t ORDER BY x"
	if _, err := mc.db.ExecContext(mc.ctx, refused); err == nil {
		t.Errorf("oracle: expected ORDER BY on the declared name to be refused")
	}
}

// TestReview_QuantifiedComparisonOperandParens covers the grouping on the left
// of a quantified comparison, which the two forms the server rewrites keep and
// the rest lose. Writing the left operand with no precedence constraint at all
// dropped it everywhere, and `(1 OR 0) = ANY (SELECT 0)` then read as
// `1 or 0 in (select 0)`, which is 1 where the expression it came from is 0.
//
// Measured on mariadb:11.8.8, with the value of each expression read from the
// server beside the view's own, so a difference in grouping shows up as a
// difference in the answer and not only in the text.
func TestReview_QuantifiedComparisonOperandParens(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)

	for i, expr := range []string{
		// The `in` form keeps the parentheses an operator looser than IN
		// needs.
		"(1 OR 0) = ANY (SELECT 0)",
		"(1 AND 1) = ANY (SELECT 1)",
		"(1 XOR 0) = ANY (SELECT 1)",
		"(1 = 1) = ANY (SELECT 1)",
		"(1 BETWEEN 0 AND 2) = ANY (SELECT 1)",
		"(1 OR 0) = SOME (SELECT 0)",
		"(1 OR 0) <> ALL (SELECT 1)",
		// An operand that binds at least as tightly as IN needs none.
		"1 + 1 = ANY (SELECT 2)",
		"(NOT 1) = ANY (SELECT 0)",
		"(1 LIKE '1') = ANY (SELECT 1)",
		"1 = ANY (SELECT 1)",
		// The forms the server leaves as written lose the grouping, and the
		// answer with it. That is the server's own output, so it is what the
		// catalog writes.
		"(1 OR 0) > ALL (SELECT 1)",
		"(1 OR 0) < ANY (SELECT 1)",
	} {
		t.Run(fmt.Sprintf("case%d", i), func(t *testing.T) {
			view := fmt.Sprintf("qc%d", i)
			runOnBoth(t, mc, c, "CREATE VIEW "+view+" AS SELECT "+expr+" AS x;")
			assertViewBodyMatches(t, mc, c, view)

			// The view's own answer, which is what the stored grouping decides.
			var want string
			oracleScan(t, mc, "SELECT x FROM "+view, &want)
			var replay string
			body := extractSelectBody(c.ShowCreateView("testdb", view))
			oracleScan(t, mc, "SELECT * FROM ("+body+") AS q", &replay)
			if replay != want {
				t.Errorf("replaying what omni wrote changed the answer: %s became %s\n  %s",
					want, replay, body)
			}
		})
	}
}

// TestReview_IntroducedLiteralDefaultIsConverted covers an introduced literal
// as the default of a column that stores a value. The literal's bytes are read
// in the charset it names and written in the column's own, which is what the
// server holds: quoting the bytes as they stood had changed the value, two
// latin1 characters becoming one.
//
// Measured on mariadb:11.8.8 in a utf8mb4 database, with the stored value read
// back from a row so that a changed encoding shows up as a changed value.
//
// A column whose own charset is not the one the DDL is read in is covered by
// TestReview_DefaultIsWrittenInTheOutputCharset, which is the other half of
// the same conversion: this one reads the literal into the column's charset,
// that one writes the column's bytes into the output's.
func TestReview_IntroducedLiteralDefaultIsConverted(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, ddl string }{
		{"latin1_high_byte", "CREATE TABLE t (a VARCHAR(20) DEFAULT (_latin1'ä'))"},
		{"latin1_unparenthesized", "CREATE TABLE t (a VARCHAR(20) DEFAULT _latin1'ä')"},
		{"latin1_mixed", "CREATE TABLE t (a VARCHAR(20) DEFAULT (_latin1'aäb'))"},
		{"latin1_ascii_only", "CREATE TABLE t (a VARCHAR(20) DEFAULT (_latin1'abc'))"},
		{"latin1_empty", "CREATE TABLE t (a VARCHAR(20) DEFAULT (_latin1''))"},
		{"char_column", "CREATE TABLE t (a CHAR(5) DEFAULT (_latin1'ä'))"},
		{"same_charset", "CREATE TABLE t (a VARCHAR(20) DEFAULT (_utf8mb4'ä'))"},
		{"utf8mb3", "CREATE TABLE t (a VARCHAR(20) DEFAULT (_utf8mb3'ä'))"},
		{"ascii_introducer", "CREATE TABLE t (a VARCHAR(20) DEFAULT (_ascii'abc'))"},
		{"binary_introducer", "CREATE TABLE t (a VARCHAR(20) DEFAULT (_binary'ä'))"},
		{"no_introducer", "CREATE TABLE t (a VARCHAR(20) DEFAULT 'ä')"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, tc.ddl+";")
			assertDefaultMatches(t, mc, c)

			// The value the server stores for its own DDL, against the value
			// it stores for omni's.
			if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
				t.Fatalf("oracle insert: %v", err)
			}
			var want string
			oracleScan(t, mc, "SELECT HEX(a) FROM t", &want)
			if _, err := mc.db.ExecContext(mc.ctx, "DROP TABLE t"); err != nil {
				t.Fatalf("oracle drop: %v", err)
			}
			replay := c.ShowCreateTable("testdb", "t")
			if _, err := mc.db.ExecContext(mc.ctx, replay); err != nil {
				t.Fatalf("oracle cannot read back what omni wrote: %v\n%s", err, replay)
			}
			if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
				t.Fatalf("oracle insert after replay: %v", err)
			}
			var got string
			oracleScan(t, mc, "SELECT HEX(a) FROM t", &got)
			if got != want {
				t.Errorf("replaying what omni wrote changed the default: %s became %s\n%s",
					want, got, replay)
			}
		})
	}
}

// TestReview_GeneratedColumnDependencyNondeterminism covers a materialized
// generated column that reads another generated column. MariaDB carries
// nondeterminism through the columns an expression reads, so a STORED column —
// or a VIRTUAL one a key covers — is refused for reading a column that calls
// RAND(), however many steps away the call is. Only the column's own function
// calls were checked, so the catalog accepted a schema the server refuses.
//
// Measured on mariadb:11.8.8, message included: the server names the column
// the value came through, where a direct call it names with its parentheses.
func TestReview_GeneratedColumnDependencyNondeterminism(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct {
		name    string
		stmts   []string
		accepts bool
	}{
		{"stored_reads_a_rand_column", []string{
			"CREATE TABLE t (a DOUBLE AS (RAND()) VIRTUAL, b DOUBLE AS (a) STORED)",
		}, false},
		{"stored_reads_it_in_an_expression", []string{
			"CREATE TABLE t (a DOUBLE AS (RAND()) VIRTUAL, b DOUBLE AS (a + 1) STORED)",
		}, false},
		{"stored_reads_it_beside_a_plain_column", []string{
			"CREATE TABLE t (x INT, a DOUBLE AS (RAND()) VIRTUAL, b DOUBLE AS (x + a) STORED)",
		}, false},
		{"two_hops", []string{
			"CREATE TABLE t (a DOUBLE AS (RAND()) VIRTUAL, b DOUBLE AS (a) VIRTUAL," +
				" c DOUBLE AS (b) STORED)",
		}, false},
		{"keyed_virtual_reads_it", []string{
			"CREATE TABLE t (a DOUBLE AS (RAND()) VIRTUAL, b DOUBLE AS (a) VIRTUAL, KEY (b))",
		}, false},
		{"now_through_a_dependency", []string{
			"CREATE TABLE t (a DATETIME AS (NOW()) VIRTUAL, b DATETIME AS (a) STORED)",
		}, false},
		{"the_stored_column_is_added_later", []string{
			"CREATE TABLE t (a DOUBLE AS (RAND()) VIRTUAL)",
			"ALTER TABLE t ADD COLUMN b DOUBLE AS (a) STORED",
		}, false},
		{"the_key_arrives_later", []string{
			"CREATE TABLE t (a DOUBLE AS (RAND()) VIRTUAL, b DOUBLE AS (a) VIRTUAL)",
			"CREATE INDEX i ON t(b)",
		}, false},
		// An unmaterialized VIRTUAL column may read it, and a deterministic
		// chain is fine however long.
		{"virtual_reads_it", []string{
			"CREATE TABLE t (a DOUBLE AS (RAND()) VIRTUAL, b DOUBLE AS (a) VIRTUAL)",
		}, true},
		{"deterministic_chain", []string{
			"CREATE TABLE t (x INT, a DOUBLE AS (x) VIRTUAL, b DOUBLE AS (a) STORED)",
		}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)

			var oracleErr, omniErr error
			for _, stmt := range tc.stmts {
				if _, err := mc.db.ExecContext(mc.ctx, stmt); err != nil {
					oracleErr = err
				}
				results, err := c.Exec(stmt+";", nil)
				if err != nil {
					omniErr = err
				}
				for _, r := range results {
					if r.Error != nil {
						omniErr = r.Error
					}
				}
			}
			if tc.accepts != (oracleErr == nil) {
				t.Fatalf("oracle: gave %v, want accepted=%v", oracleErr, tc.accepts)
			}
			if tc.accepts != (omniErr == nil) {
				t.Fatalf("omni: gave %v, want accepted=%v", omniErr, tc.accepts)
			}
			if tc.accepts {
				want := oracleShow(t, mc, "SHOW CREATE TABLE t")
				if got := c.ShowCreateTable("testdb", "t"); got != want {
					t.Errorf("table mismatch\n--- oracle ---\n%s\n--- omni ---\n%s", want, got)
				}
				return
			}
			if want := messageTailOf(oracleErr); !strings.Contains(omniErr.Error(), want) {
				t.Errorf("message mismatch\n  oracle: %s\n  omni:   %s", want, omniErr)
			}
		})
	}
}

// TestReview_ReplacingAGeneratedDependency covers the other side of the
// dependency rule: replacing a column's expression can make a materialized
// column that reads it invalid. Only the replacement's own materialization was
// checked, so a MODIFY could leave behind a table the catalog's own CREATE
// TABLE validation refuses.
//
// Measured on mariadb:11.8.8, message included. A refusal leaves the table as
// it was, which each refused case checks.
func TestReview_ReplacingAGeneratedDependency(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct {
		name    string
		setup   string
		alter   string
		accepts bool
	}{
		{
			name:  "a_stored_column_reads_it",
			setup: "CREATE TABLE t (a DOUBLE AS (1) VIRTUAL, b DOUBLE AS (a) STORED)",
			alter: "ALTER TABLE t MODIFY a DOUBLE AS (RAND()) VIRTUAL",
		},
		{
			name:  "a_keyed_virtual_column_reads_it",
			setup: "CREATE TABLE t (a DOUBLE AS (1) VIRTUAL, b DOUBLE AS (a) VIRTUAL, KEY (b))",
			alter: "ALTER TABLE t MODIFY a DOUBLE AS (RAND()) VIRTUAL",
		},
		{
			name: "two_hops_from_the_column_replaced",
			setup: "CREATE TABLE t (a DOUBLE AS (1) VIRTUAL, b DOUBLE AS (a) VIRTUAL," +
				" c DOUBLE AS (b) STORED)",
			alter: "ALTER TABLE t MODIFY a DOUBLE AS (RAND()) VIRTUAL",
		},
		{
			name:  "the_replacement_is_itself_materialized",
			setup: "CREATE TABLE t (a DOUBLE AS (1) STORED)",
			alter: "ALTER TABLE t MODIFY a DOUBLE AS (RAND()) STORED",
		},
		// An unmaterialized dependent may read it, a deterministic
		// replacement keeps everything valid, and a column nothing reads is
		// its own business.
		{
			name:    "only_an_unmaterialized_column_reads_it",
			setup:   "CREATE TABLE t (a DOUBLE AS (1) VIRTUAL, b DOUBLE AS (a) VIRTUAL)",
			alter:   "ALTER TABLE t MODIFY a DOUBLE AS (RAND()) VIRTUAL",
			accepts: true,
		},
		{
			name:    "the_replacement_stays_deterministic",
			setup:   "CREATE TABLE t (a DOUBLE AS (1) VIRTUAL, b DOUBLE AS (a) STORED)",
			alter:   "ALTER TABLE t MODIFY a DOUBLE AS (2) VIRTUAL",
			accepts: true,
		},
		{
			name:    "nothing_reads_it",
			setup:   "CREATE TABLE t (a DOUBLE AS (1) VIRTUAL)",
			alter:   "ALTER TABLE t MODIFY a DOUBLE AS (RAND()) VIRTUAL",
			accepts: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, tc.setup+";")
			before := c.ShowCreateTable("testdb", "t")

			_, oracleErr := mc.db.ExecContext(mc.ctx, tc.alter)
			if tc.accepts != (oracleErr == nil) {
				t.Fatalf("oracle: %q gave %v, want accepted=%v", tc.alter, oracleErr, tc.accepts)
			}

			results, err := c.Exec(tc.alter+";", nil)
			omniErr := err
			for _, r := range results {
				if r.Error != nil {
					omniErr = r.Error
				}
			}
			if tc.accepts != (omniErr == nil) {
				t.Fatalf("omni: %q gave %v, want accepted=%v", tc.alter, omniErr, tc.accepts)
			}

			want := oracleShow(t, mc, "SHOW CREATE TABLE t")
			if got := c.ShowCreateTable("testdb", "t"); got != want {
				t.Errorf("table mismatch\n--- oracle ---\n%s\n--- omni ---\n%s", want, got)
			}
			if tc.accepts {
				return
			}
			// A refused ALTER leaves the table as it was.
			if after := c.ShowCreateTable("testdb", "t"); after != before {
				t.Errorf("the refused ALTER changed the table\n--- before ---\n%s\n--- after ---\n%s",
					before, after)
			}
			if msg := messageTailOf(oracleErr); !strings.Contains(omniErr.Error(), msg) {
				t.Errorf("message mismatch\n  oracle: %s\n  omni:   %s", msg, omniErr)
			}
		})
	}
}

// TestReview_ChangingAGeneratedDependencyName covers the same replacement
// under CHANGE, which renames the column its dependents read. The server
// rewrites those references and reports the name as it now stands: measured on
// mariadb:11.8.8, `ALTER TABLE t CHANGE a z DOUBLE AS (RAND()) VIRTUAL` with
// `b DOUBLE AS (a) STORED` in place answers 1901 naming z.
func TestReview_ChangingAGeneratedDependencyName(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)

	runOnBoth(t, mc, c, "CREATE TABLE t (a DOUBLE AS (1) VIRTUAL, b DOUBLE AS (a) STORED);")

	const alter = "ALTER TABLE t CHANGE a z DOUBLE AS (RAND()) VIRTUAL"
	_, oracleErr := mc.db.ExecContext(mc.ctx, alter)
	if oracleErr == nil {
		t.Fatalf("oracle: expected %q to be refused", alter)
	}

	results, err := c.Exec(alter+";", nil)
	omniErr := err
	for _, r := range results {
		if r.Error != nil {
			omniErr = r.Error
		}
	}
	if omniErr == nil {
		t.Fatalf("omni: accepted %q, which the server refuses with %v", alter, oracleErr)
	}
	if msg := messageTailOf(oracleErr); !strings.Contains(omniErr.Error(), msg) {
		t.Errorf("message mismatch\n  oracle: %s\n  omni:   %s", msg, omniErr)
	}
}

// TestReview_DefaultIsWrittenInTheOutputCharset covers a default the writer
// reconstructs from bytes on a column whose charset is not the one the DDL is
// read in. The value is the column's bytes, and the text has to carry them in
// the output's charset, which is where the server writes them too: writing the
// column's bytes into the text changed the value on replay, a latin1 column's
// two bytes coming back as one character.
//
// Measured on mariadb:11.8.8 with a utf8mb4 client, the stored bytes read back
// from a row on each side so a changed encoding fails as a changed value.
func TestReview_DefaultIsWrittenInTheOutputCharset(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, ddl string }{
		{"latin1_column_introduced", "CREATE TABLE t (a VARCHAR(20) CHARACTER SET latin1 DEFAULT (_latin1'ä'))"},
		{"latin1_column_hex", "CREATE TABLE t (a VARCHAR(20) CHARACTER SET latin1 DEFAULT (X'E4'))"},
		{"latin1_column_hex_0x", "CREATE TABLE t (a VARCHAR(20) CHARACTER SET latin1 DEFAULT (0xE4))"},
		{"latin1_column_ascii_hex", "CREATE TABLE t (a VARCHAR(20) CHARACTER SET latin1 DEFAULT (X'41'))"},
		{"latin1_column_plain_literal", "CREATE TABLE t (a VARCHAR(20) CHARACTER SET latin1 DEFAULT 'ä')"},
		{"latin1_char_column", "CREATE TABLE t (a CHAR(5) CHARACTER SET latin1 DEFAULT (_latin1'ä'))"},
		{"utf8mb4_column_hex", "CREATE TABLE t (a VARCHAR(20) DEFAULT (X'41'))"},
		{"utf8mb4_column_introduced", "CREATE TABLE t (a VARCHAR(20) DEFAULT (_latin1'ä'))"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, tc.ddl+";")
			assertDefaultMatches(t, mc, c)

			if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
				t.Fatalf("oracle insert: %v", err)
			}
			var want string
			oracleScan(t, mc, "SELECT HEX(a) FROM t", &want)
			if _, err := mc.db.ExecContext(mc.ctx, "DROP TABLE t"); err != nil {
				t.Fatalf("oracle drop: %v", err)
			}
			replay := c.ShowCreateTable("testdb", "t")
			if _, err := mc.db.ExecContext(mc.ctx, replay); err != nil {
				t.Fatalf("oracle cannot read back what omni wrote: %v\n%s", err, replay)
			}
			if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
				t.Fatalf("oracle insert after replay: %v", err)
			}
			var got string
			oracleScan(t, mc, "SELECT HEX(a) FROM t", &got)
			if got != want {
				t.Errorf("replaying what omni wrote changed the default: %s became %s\n%s",
					want, got, replay)
			}
		})
	}
}

// TestReview_DefaultOutsideTheConvertibleBand covers the one band the writer
// will not convert: MariaDB's latin1 is cp1252, so its bytes 0x80 to 0x9F are
// not the code points of the same number. A default reaching into it keeps the
// literal as written — the server writes the character instead — and the
// value, which is what the replay here checks.
//
// Measured on mariadb:11.8.8: a latin1 column defaulting to X'80' is written
// back as '€', and the byte stored is 0x80.
func TestReview_DefaultOutsideTheConvertibleBand(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)

	const ddl = "CREATE TABLE t (a VARCHAR(20) CHARACTER SET latin1 DEFAULT (X'80'))"
	runOnBoth(t, mc, c, ddl+";")

	if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
		t.Fatalf("oracle insert: %v", err)
	}
	var want string
	oracleScan(t, mc, "SELECT HEX(a) FROM t", &want)
	if want != "80" {
		t.Fatalf("oracle stored %s, want the byte 80", want)
	}
	if _, err := mc.db.ExecContext(mc.ctx, "DROP TABLE t"); err != nil {
		t.Fatalf("oracle drop: %v", err)
	}
	replay := c.ShowCreateTable("testdb", "t")
	if _, err := mc.db.ExecContext(mc.ctx, replay); err != nil {
		t.Fatalf("oracle cannot read back what omni wrote: %v\n%s", err, replay)
	}
	if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
		t.Fatalf("oracle insert after replay: %v", err)
	}
	var got string
	oracleScan(t, mc, "SELECT HEX(a) FROM t", &got)
	if got != want {
		t.Errorf("replaying what omni wrote changed the default: %s became %s\n%s",
			want, got, replay)
	}
}

// TestReview_RenamedColumnKeepsItsDependents covers what a renamed column does
// to the record of what reads it. The names a generated expression reads are
// recorded when the column is defined, and a rename has to carry them, or the
// dependency goes unseen and a later statement is accepted that the server
// refuses.
//
// Measured on mariadb:11.8.8, message included: the server names the column as
// it now stands, since it rewrites the expressions that read it.
func TestReview_RenamedColumnKeepsItsDependents(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct {
		name    string
		stmts   []string
		accepts bool
	}{
		{"rename_then_modify_to_nondeterministic", []string{
			"CREATE TABLE t (a DOUBLE AS (1) VIRTUAL, b DOUBLE AS (a) STORED)",
			"ALTER TABLE t RENAME COLUMN a TO z",
			"ALTER TABLE t MODIFY z DOUBLE AS (RAND()) VIRTUAL",
		}, false},
		{"rename_then_index_the_dependent", []string{
			"CREATE TABLE t (a DOUBLE AS (RAND()) VIRTUAL, b DOUBLE AS (a) VIRTUAL)",
			"ALTER TABLE t RENAME COLUMN a TO z",
			"CREATE INDEX i ON t(b)",
		}, false},
		{"rename_then_add_a_stored_dependent", []string{
			"CREATE TABLE t (a DOUBLE AS (RAND()) VIRTUAL)",
			"ALTER TABLE t RENAME COLUMN a TO z",
			"ALTER TABLE t ADD COLUMN b DOUBLE AS (z) STORED",
		}, false},
		{"rename_then_modify_to_deterministic", []string{
			"CREATE TABLE t (a DOUBLE AS (1) VIRTUAL, b DOUBLE AS (a) STORED)",
			"ALTER TABLE t RENAME COLUMN a TO z",
			"ALTER TABLE t MODIFY z DOUBLE AS (2) VIRTUAL",
		}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)

			var oracleErr, omniErr error
			for _, stmt := range tc.stmts {
				if _, err := mc.db.ExecContext(mc.ctx, stmt); err != nil && oracleErr == nil {
					oracleErr = err
				}
				results, err := c.Exec(stmt+";", nil)
				if err != nil && omniErr == nil {
					omniErr = err
				}
				for _, r := range results {
					if r.Error != nil && omniErr == nil {
						omniErr = r.Error
					}
				}
			}
			if tc.accepts != (oracleErr == nil) {
				t.Fatalf("oracle: gave %v, want accepted=%v", oracleErr, tc.accepts)
			}
			if tc.accepts != (omniErr == nil) {
				t.Fatalf("omni: gave %v, want accepted=%v", omniErr, tc.accepts)
			}
			if tc.accepts {
				return
			}
			if msg := messageTailOf(oracleErr); !strings.Contains(omniErr.Error(), msg) {
				t.Errorf("message mismatch\n  oracle: %s\n  omni:   %s", msg, omniErr)
			}
		})
	}
}

// TestReview_RejectedAlterLeavesDependenciesAlone covers what a rejected
// multi-clause ALTER leaves behind. One clause renames a column and a later
// one fails, so the statement is refused and the table is restored — and the
// record of what each generated expression reads has to be restored with it,
// or a later statement is accepted that the server refuses. The snapshot was
// sharing that record with the table it was taken from, so a rename written
// before the failure survived the rollback.
//
// Measured on mariadb:11.8.8, message included. The server leaves the table as
// it was, which each case checks before going on.
func TestReview_RejectedAlterLeavesDependenciesAlone(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct {
		name   string
		setup  string
		failed string
		then   string
	}{
		{
			name:   "rename_first",
			setup:  "CREATE TABLE t (a DOUBLE AS (1) VIRTUAL, b DOUBLE AS (a) STORED)",
			failed: "ALTER TABLE t RENAME COLUMN a TO z, ADD COLUMN b INT",
			then:   "ALTER TABLE t MODIFY a DOUBLE AS (RAND()) VIRTUAL",
		},
		{
			name:   "rename_second",
			setup:  "CREATE TABLE t (a DOUBLE AS (1) VIRTUAL, b DOUBLE AS (a) STORED)",
			failed: "ALTER TABLE t ADD COLUMN b INT, RENAME COLUMN a TO z",
			then:   "ALTER TABLE t MODIFY a DOUBLE AS (RAND()) VIRTUAL",
		},
		{
			name:   "then_an_index",
			setup:  "CREATE TABLE t (a DOUBLE AS (RAND()) VIRTUAL, b DOUBLE AS (a) VIRTUAL)",
			failed: "ALTER TABLE t RENAME COLUMN a TO z, ADD COLUMN b INT",
			then:   "CREATE INDEX i ON t(b)",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, tc.setup+";")
			before := c.ShowCreateTable("testdb", "t")

			// The ALTER both sides refuse, which must leave the table alone.
			if _, err := mc.db.ExecContext(mc.ctx, tc.failed); err == nil {
				t.Fatalf("oracle: expected %q to be refused", tc.failed)
			}
			results, err := c.Exec(tc.failed+";", nil)
			omniErr := err
			for _, r := range results {
				if r.Error != nil {
					omniErr = r.Error
				}
			}
			if omniErr == nil {
				t.Fatalf("omni: accepted %q, which the server refuses", tc.failed)
			}
			want := oracleShow(t, mc, "SHOW CREATE TABLE t")
			if got := c.ShowCreateTable("testdb", "t"); got != want {
				t.Fatalf("after the refused ALTER\n--- oracle ---\n%s\n--- omni ---\n%s", want, got)
			}
			if after := c.ShowCreateTable("testdb", "t"); after != before {
				t.Errorf("the refused ALTER changed the table\n--- before ---\n%s\n--- after ---\n%s",
					before, after)
			}

			// And the dependency is still seen, which is what the restored
			// record decides.
			_, oracleErr := mc.db.ExecContext(mc.ctx, tc.then)
			if oracleErr == nil {
				t.Fatalf("oracle: expected %q to be refused", tc.then)
			}
			results, err = c.Exec(tc.then+";", nil)
			omniErr = err
			for _, r := range results {
				if r.Error != nil {
					omniErr = r.Error
				}
			}
			if omniErr == nil {
				t.Fatalf("omni: accepted %q, which the server refuses with %v", tc.then, oracleErr)
			}
			if msg := messageTailOf(oracleErr); !strings.Contains(omniErr.Error(), msg) {
				t.Errorf("message mismatch\n  oracle: %s\n  omni:   %s", msg, omniErr)
			}
		})
	}
}

// TestReview_HexNotationsPartCompany covers the two hex notations, which are
// not the same literal: X'..' is a string, where 0x.. and a bit literal are
// numbers. Reading all three as bytes made `INT DEFAULT (X'37')` 55 instead of
// 7, changing the value on replay.
//
// Measured on mariadb:11.8.8 through SHOW CREATE TABLE. Where a column reads
// bytes rather than text — BIT, BINARY — all three are the bytes they hold,
// and the notations agree again.
//
// https://mariadb.com/docs/server/reference/sql-structure/sql-language-structure/hexadecimal-literals
func TestReview_HexNotationsPartCompany(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, ddl string }{
		// A numeric column reads X'..' as the characters it spells and the
		// other two as the number their bytes spell.
		{"int_x_quoted", "CREATE TABLE t (a INT DEFAULT (X'37'))"},
		{"int_x_quoted_bare", "CREATE TABLE t (a INT DEFAULT X'37')"},
		{"int_0x", "CREATE TABLE t (a INT DEFAULT 0x37)"},
		{"int_bits", "CREATE TABLE t (a INT DEFAULT b'110111')"},
		{"int_x_quoted_two_digits", "CREATE TABLE t (a INT DEFAULT (X'3137'))"},
		{"bigint_x_quoted", "CREATE TABLE t (a BIGINT DEFAULT (X'3132333435'))"},
		{"decimal_x_quoted", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT (X'37'))"},
		{"decimal_0x", "CREATE TABLE t (a DECIMAL(5,2) DEFAULT 0x37)"},
		{"double_x_quoted_fraction", "CREATE TABLE t (a DOUBLE DEFAULT (X'372E35'))"},
		{"boolean_x_quoted", "CREATE TABLE t (a BOOLEAN DEFAULT (X'37'))"},
		{"boolean_0x", "CREATE TABLE t (a BOOLEAN DEFAULT 0x37)"},

		// A YEAR column reads them the same two ways.
		{"year_x_quoted", "CREATE TABLE t (a YEAR DEFAULT (X'37'))"},
		{"year_x_quoted_two_digits", "CREATE TABLE t (a YEAR DEFAULT (X'3939'))"},
		{"year_0x", "CREATE TABLE t (a YEAR DEFAULT 0x37)"},

		// A column that reads bytes reads all three alike.
		{"bit_x_quoted", "CREATE TABLE t (a BIT(8) DEFAULT X'37')"},
		{"bit_x_quoted_letter", "CREATE TABLE t (a BIT(8) DEFAULT X'41')"},
		{"bit_0x", "CREATE TABLE t (a BIT(8) DEFAULT 0x37)"},
		{"binary_x_quoted", "CREATE TABLE t (a BINARY(2) DEFAULT (X'41'))"},
		{"varchar_x_quoted", "CREATE TABLE t (a VARCHAR(4) DEFAULT (X'41'))"},
		{"varchar_0x", "CREATE TABLE t (a VARCHAR(4) DEFAULT 0x41)"},
		{"text_keeps_the_literal", "CREATE TABLE t (a TEXT DEFAULT (X'41'))"},
		{"date_x_quoted", "CREATE TABLE t (a DATE DEFAULT (X'323032302D30312D3031'))"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, tc.ddl+";")
			assertDefaultMatches(t, mc, c)

			// And what omni wrote means the same to the server.
			if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
				t.Fatalf("oracle insert: %v", err)
			}
			var want string
			oracleScan(t, mc, "SELECT HEX(a) FROM t", &want)
			if _, err := mc.db.ExecContext(mc.ctx, "DROP TABLE t"); err != nil {
				t.Fatalf("oracle drop: %v", err)
			}
			replay := c.ShowCreateTable("testdb", "t")
			if _, err := mc.db.ExecContext(mc.ctx, replay); err != nil {
				t.Fatalf("oracle cannot read back what omni wrote: %v\n%s", err, replay)
			}
			if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
				t.Fatalf("oracle insert after replay: %v", err)
			}
			var got string
			oracleScan(t, mc, "SELECT HEX(a) FROM t", &got)
			if got != want {
				t.Errorf("replaying what omni wrote changed the default: %s became %s\n%s",
					want, got, replay)
			}
		})
	}
}

// TestReview_XQuotedNonNumberIsRefused covers the other half of X'..' being a
// string: on a numeric column the characters have to spell a number, and the
// server refuses the default outright where they do not — which is how the
// two notations were told apart. Measured on mariadb:11.8.8:
// `INT DEFAULT (X'41')` answers 1067, the character A spelling no number,
// where `INT DEFAULT 0x41` is 65.
func TestReview_XQuotedNonNumberIsRefused(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)

	for _, ddl := range []string{
		"CREATE TABLE t (a INT DEFAULT (X'41'))",
		"CREATE TABLE t (a INT DEFAULT (X''))",
		"CREATE TABLE t (a YEAR DEFAULT (X'41'))",
	} {
		_, oracleErr := mc.db.ExecContext(mc.ctx, ddl)
		if oracleErr == nil {
			t.Errorf("oracle: expected %q to be refused", ddl)
			if _, err := mc.db.ExecContext(mc.ctx, "DROP TABLE t"); err != nil {
				t.Fatalf("oracle drop: %v", err)
			}
			continue
		}
		// omni does not check 1067 yet, so what it writes has to parse and
		// mean what the statement said.
		results, err := c.Exec(ddl+";", nil)
		if err != nil {
			continue // a parse error counts as a refusal
		}
		refused := false
		for _, r := range results {
			if r.Error != nil {
				refused = true
			}
		}
		if refused {
			continue
		}
		create := c.ShowCreateTable("testdb", "t")
		again := scenarioNewCatalog(t)
		results, err = again.Exec(create+";", nil)
		if err != nil {
			t.Errorf("omni wrote a default it cannot read back: %v\n%s", err, create)
			continue
		}
		for _, r := range results {
			if r.Error != nil {
				t.Errorf("omni wrote a default it cannot read back: %v\n%s", r.Error, create)
			}
		}
		if _, err := c.Exec("DROP TABLE t;", nil); err != nil {
			t.Fatalf("omni drop: %v", err)
		}
	}
}

// TestReview_ForeignKeyIndexOverAGeneratedColumn covers the index a foreign
// key brings with it, which materializes a VIRTUAL generated column as any
// other key does. Every explicit index path checked that; the implicit one did
// not, so a foreign key could leave an index over RAND() behind.
//
// Measured on mariadb:11.8.8, message included.
func TestReview_ForeignKeyIndexOverAGeneratedColumn(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct {
		name    string
		stmts   []string
		accepts bool
	}{
		{"added_by_alter", []string{
			"CREATE TABLE p (a DOUBLE PRIMARY KEY)",
			"CREATE TABLE t (a DOUBLE AS (RAND()) VIRTUAL)",
			"ALTER TABLE t ADD CONSTRAINT fk FOREIGN KEY (a) REFERENCES p(a)",
		}, false},
		{"declared_with_the_table", []string{
			"CREATE TABLE p (a DOUBLE PRIMARY KEY)",
			"CREATE TABLE t (a DOUBLE AS (RAND()) VIRTUAL," +
				" CONSTRAINT fk FOREIGN KEY (a) REFERENCES p(a))",
		}, false},
		{"over_a_column_reading_a_rand_one", []string{
			"CREATE TABLE p (a DOUBLE PRIMARY KEY)",
			"CREATE TABLE t (r DOUBLE AS (RAND()) VIRTUAL, a DOUBLE AS (r) VIRTUAL)",
			"ALTER TABLE t ADD CONSTRAINT fk FOREIGN KEY (a) REFERENCES p(a)",
		}, false},
		{"deterministic_is_fine", []string{
			"CREATE TABLE p (a DOUBLE PRIMARY KEY)",
			"CREATE TABLE t (x DOUBLE, a DOUBLE AS (x) VIRTUAL)",
			"ALTER TABLE t ADD CONSTRAINT fk FOREIGN KEY (a) REFERENCES p(a)",
		}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)

			var oracleErr, omniErr error
			for _, stmt := range tc.stmts {
				if _, err := mc.db.ExecContext(mc.ctx, stmt); err != nil && oracleErr == nil {
					oracleErr = err
				}
				results, err := c.Exec(stmt+";", nil)
				if err != nil && omniErr == nil {
					omniErr = err
				}
				for _, r := range results {
					if r.Error != nil && omniErr == nil {
						omniErr = r.Error
					}
				}
			}
			if tc.accepts != (oracleErr == nil) {
				t.Fatalf("oracle: gave %v, want accepted=%v", oracleErr, tc.accepts)
			}
			if tc.accepts != (omniErr == nil) {
				t.Fatalf("omni: gave %v, want accepted=%v", omniErr, tc.accepts)
			}
			if tc.accepts {
				want := oracleShow(t, mc, "SHOW CREATE TABLE t")
				if got := c.ShowCreateTable("testdb", "t"); got != want {
					t.Errorf("table mismatch\n--- oracle ---\n%s\n--- omni ---\n%s", want, got)
				}
				return
			}
			if msg := messageTailOf(oracleErr); !strings.Contains(omniErr.Error(), msg) {
				t.Errorf("message mismatch\n  oracle: %s\n  omni:   %s", msg, omniErr)
			}
		})
	}
}

// TestReview_ResolverRewritesOutsideViews covers the rewrites the server
// applies to a stored expression, which are not a view's own: a CHECK, a
// DEFAULT and a generated column come back with the same ones a view body
// gets. Only the view path applied them, so a CHECK or DEFAULT kept an
// operator the server does not store — and SOUNDS LIKE, kept at comparison
// precedence, lost the parentheses its left operand needed and changed the
// expression's value.
//
// Measured on mariadb:11.8.8 through SHOW CREATE TABLE.
//
// A generated column's expression is written by a second writer of its own,
// which none of these rewrites reach, and it is left out here: measured, the
// server stores `b INT AS ((2=1) SOUNDS LIKE '0')` as
// soundex(2 = 1) = soundex('0'), `AS (a AND c)` as `a` <> 0 and `c` <> 0 and
// `AS (NOT (a = 1))` as `a` <> 1, where that writer keeps the operator and
// doubles the parentheses around it.
func TestReview_ResolverRewritesOutsideViews(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, ddl string }{
		{"check_sounds_like", "CREATE TABLE t (a INT, CHECK ((2=1) SOUNDS LIKE '0'))"},
		{"default_sounds_like", "CREATE TABLE t (a INT DEFAULT ((2=1) SOUNDS LIKE '0'))"},
		{"sounds_like_plain", "CREATE TABLE t (a INT, CHECK (a SOUNDS LIKE '1'))"},
		{"check_not_comparison", "CREATE TABLE t (a INT, CHECK (NOT (a = 1)))"},
		{"check_not_between", "CREATE TABLE t (a INT, CHECK (NOT (a BETWEEN 1 AND 2)))"},
		{"check_and_of_columns", "CREATE TABLE t (a INT, b INT, CHECK (a AND b))"},
		{"check_or_of_columns", "CREATE TABLE t (a INT, b INT, CHECK (a OR b))"},
		{"check_bare_column", "CREATE TABLE t (a INT, CHECK (a))"},
		{"default_not_comparison", "CREATE TABLE t (a INT DEFAULT (NOT (2 = 1)))"},
		{"default_and_of_literals", "CREATE TABLE t (a INT DEFAULT (1 AND 2))"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, tc.ddl+";")

			want := oracleShow(t, mc, "SHOW CREATE TABLE t")
			if got := c.ShowCreateTable("testdb", "t"); got != want {
				t.Errorf("table mismatch\n--- oracle ---\n%s\n--- omni ---\n%s", want, got)
			}

			// And what omni wrote goes back in, which the dropped
			// parentheses did not.
			if _, err := mc.db.ExecContext(mc.ctx, "DROP TABLE t"); err != nil {
				t.Fatalf("oracle drop: %v", err)
			}
			replay := c.ShowCreateTable("testdb", "t")
			if _, err := mc.db.ExecContext(mc.ctx, replay); err != nil {
				t.Fatalf("oracle cannot read back what omni wrote: %v\n%s", err, replay)
			}
			if got := oracleShow(t, mc, "SHOW CREATE TABLE t"); got != want {
				t.Errorf("replaying what omni wrote changed the table\n--- was ---\n%s\n--- now ---\n%s",
					want, got)
			}
		})
	}
}

// TestReview_TriggerOrderingLeavesTheBodyReadable covers the text left behind
// when the ordering clause is cut out of a trigger's source. The clause goes
// with the whitespace after it and the text before it is kept as written: a
// line comment ending right before the clause needs its newline, or the body
// ends up inside the comment.
//
// Measured on mariadb:11.8.8. The server drops a comment altogether — a
// difference of its own, so the cases carrying one compare what omni wrote by
// reading it back rather than against the server's text.
func TestReview_TriggerOrderingLeavesTheBodyReadable(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct {
		name     string
		ddl      string
		comments bool
	}{
		{"line_comment_before_follows",
			"CREATE TRIGGER tr2 BEFORE INSERT ON t FOR EACH ROW -- keep order\nFOLLOWS tr1 SET NEW.a = 2",
			true},
		{"hash_comment_before_follows",
			"CREATE TRIGGER tr2 BEFORE INSERT ON t FOR EACH ROW # keep order\nFOLLOWS tr1 SET NEW.a = 2",
			true},
		{"line_comment_before_precedes",
			"CREATE TRIGGER tr2 BEFORE INSERT ON t FOR EACH ROW -- keep order\nPRECEDES tr1 SET NEW.a = 2",
			true},
		{"newline_before_follows",
			"CREATE TRIGGER tr2 BEFORE INSERT ON t FOR EACH ROW\nFOLLOWS tr1\nSET NEW.a = 2",
			false},
		{"spaces_before_follows",
			"CREATE TRIGGER tr2 BEFORE INSERT ON t FOR EACH ROW FOLLOWS tr1 SET NEW.a = 2",
			false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, "CREATE TABLE t (a INT);")
			runOnBoth(t, mc, c, "CREATE TRIGGER tr1 BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1;")
			runOnBoth(t, mc, c, tc.ddl+";")

			got := c.ShowCreateTrigger("testdb", "tr2")
			if got == "" {
				t.Fatalf("omni wrote no trigger DDL")
			}
			if strings.Contains(strings.ToUpper(got), "FOLLOWS") ||
				strings.Contains(strings.ToUpper(got), "PRECEDES") {
				t.Errorf("the ordering clause survived: %s", got)
			}
			if !strings.Contains(got, "SET NEW.a = 2") {
				t.Errorf("the body is gone: %s", got)
			}

			// What omni wrote has to go back in: a comment that lost its
			// newline swallowed the body, and the statement then did nothing
			// a trigger body can do.
			again := scenarioNewCatalog(t)
			if _, err := again.Exec("CREATE TABLE t (a INT);", nil); err != nil {
				t.Fatalf("omni setup: %v", err)
			}
			results, err := again.Exec(got+";", nil)
			if err != nil {
				t.Fatalf("omni wrote a trigger it cannot read back: %v\n%s", err, got)
			}
			for _, r := range results {
				if r.Error != nil {
					t.Fatalf("omni wrote a trigger it cannot read back: %v\n%s", r.Error, got)
				}
			}
			if tc.comments {
				return
			}
			// With no comment in the way, the text is the server's own.
			want := oracleShowTrigger(t, mc, "tr2")
			if got != want {
				t.Errorf("trigger mismatch\n--- oracle ---\n%s\n--- omni ---\n%s", want, got)
			}
		})
	}
}

// oracleShowTrigger returns the CREATE TRIGGER statement the server stored,
// which SHOW CREATE TRIGGER reports in its third column.
func oracleShowTrigger(t *testing.T, mc *mariadbContainer, name string) string {
	t.Helper()
	ddl, err := mc.showCreateTrigger(name)
	if err != nil {
		t.Fatalf("oracle SHOW CREATE TRIGGER %s: %v", name, err)
	}
	return ddl
}

// TestReview_OrderByAliasOfAConstant covers an ORDER BY or GROUP BY item that
// names an alias whose target is a whole-number constant. Substituting the
// expression turned the item into an output position: `SELECT 2 AS z FROM t
// ORDER BY z` came out `order by 2`, which names a second output column the
// view does not have. MariaDB writes an empty string there instead.
//
// Measured on mariadb:11.8.8 through SHOW CREATE VIEW, with the view read to
// confirm the body it stored is one it can use.
func TestReview_OrderByAliasOfAConstant(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)
	runOnBoth(t, mc, c, "CREATE TABLE t (a INT, b INT);")

	for i, query := range []string{
		// A whole-number constant is written as '' — anything else would be
		// read as a position.
		"SELECT 2 AS z FROM t ORDER BY z",
		"SELECT 0 AS z FROM t ORDER BY z",
		"SELECT -2 AS z FROM t ORDER BY z",
		"SELECT +2 AS z FROM t ORDER BY z",
		"SELECT (2) AS z FROM t ORDER BY z",
		"SELECT TRUE AS z FROM t ORDER BY z",
		"SELECT 2 AS z, a FROM t ORDER BY z",
		"SELECT a, 2 AS z FROM t ORDER BY z",
		"SELECT 2 AS z FROM t ORDER BY z DESC",
		"SELECT 2 AS z FROM t GROUP BY z",
		"SELECT 2 AS z FROM t GROUP BY z ORDER BY z",
		// A constant that cannot be read as a position keeps its own text.
		"SELECT 2.5 AS z FROM t ORDER BY z",
		"SELECT 1e2 AS z FROM t ORDER BY z",
		"SELECT 'x' AS z FROM t ORDER BY z",
		"SELECT NULL AS z FROM t ORDER BY z",
		// And an alias of a column, or an item that is not the bare alias, is
		// unaffected.
		"SELECT a AS z FROM t ORDER BY z",
		"SELECT a AS z FROM t GROUP BY z",
		"SELECT 1 + 1 AS z FROM t ORDER BY z",
		"SELECT 2 AS z FROM t ORDER BY z + 0",
		"SELECT 2 AS z FROM t ORDER BY 1",
		"SELECT a AS z, b FROM t ORDER BY z, b",
	} {
		t.Run(fmt.Sprintf("case%d", i), func(t *testing.T) {
			view := fmt.Sprintf("oa%d", i)
			runOnBoth(t, mc, c, "CREATE VIEW "+view+" AS "+query+";")

			assertViewBodyMatches(t, mc, c, view)
			if _, err := mc.db.ExecContext(mc.ctx, "SELECT * FROM "+view); err != nil {
				t.Errorf("oracle: the stored body is unusable: %v", err)
			}
		})
	}
}

// TestReview_DroppingEveryPrimaryKeyColumn covers an ALTER that drops every
// column a composite primary key is built on. The key goes with them, so the
// statement is accepted; 1072 is for a drop that would leave part of the key
// standing. Asking the question one command at a time refused the first drop
// of a pair that together take the whole key.
//
// Measured on mariadb:11.8.8.
func TestReview_DroppingEveryPrimaryKeyColumn(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct {
		name    string
		setup   string
		alter   string
		accepts bool
	}{
		{
			name:    "both_parts_in_one_alter",
			setup:   "CREATE TABLE t (a INT, b INT, c INT, PRIMARY KEY (a,b))",
			alter:   "ALTER TABLE t DROP COLUMN a, DROP COLUMN b",
			accepts: true,
		},
		{
			name:    "both_parts_in_reverse_order",
			setup:   "CREATE TABLE t (a INT, b INT, c INT, PRIMARY KEY (a,b))",
			alter:   "ALTER TABLE t DROP COLUMN b, DROP COLUMN a",
			accepts: true,
		},
		{
			name:    "three_parts_all_dropped",
			setup:   "CREATE TABLE t (a INT, b INT, c INT, d INT, PRIMARY KEY (a,b,c))",
			alter:   "ALTER TABLE t DROP COLUMN a, DROP COLUMN b, DROP COLUMN c",
			accepts: true,
		},
		{
			name:    "with_an_explicit_drop_primary_key",
			setup:   "CREATE TABLE t (a INT, b INT, c INT, PRIMARY KEY (a,b))",
			alter:   "ALTER TABLE t DROP PRIMARY KEY, DROP COLUMN a",
			accepts: true,
		},
		{
			name:    "a_single_column_key_goes_with_its_column",
			setup:   "CREATE TABLE t (a INT, b INT, PRIMARY KEY (a))",
			alter:   "ALTER TABLE t DROP COLUMN a",
			accepts: true,
		},
		{
			name:    "a_plain_key_is_rebuilt_over_what_is_left",
			setup:   "CREATE TABLE t (a INT, b INT, c INT, KEY k (a,b))",
			alter:   "ALTER TABLE t DROP COLUMN a",
			accepts: true,
		},
		// One part dropped while another stands is the case 1072 is for.
		{
			name:  "one_part_of_two",
			setup: "CREATE TABLE t (a INT, b INT, c INT, PRIMARY KEY (a,b))",
			alter: "ALTER TABLE t DROP COLUMN a",
		},
		{
			name:  "two_parts_of_three",
			setup: "CREATE TABLE t (a INT, b INT, c INT, d INT, PRIMARY KEY (a,b,c))",
			alter: "ALTER TABLE t DROP COLUMN a, DROP COLUMN b",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, tc.setup+";")

			_, oracleErr := mc.db.ExecContext(mc.ctx, tc.alter)
			if tc.accepts != (oracleErr == nil) {
				t.Fatalf("oracle: %q gave %v, want accepted=%v", tc.alter, oracleErr, tc.accepts)
			}
			results, err := c.Exec(tc.alter+";", nil)
			omniErr := err
			for _, r := range results {
				if r.Error != nil {
					omniErr = r.Error
				}
			}
			if tc.accepts != (omniErr == nil) {
				t.Fatalf("omni: %q gave %v, want accepted=%v", tc.alter, omniErr, tc.accepts)
			}
			if !tc.accepts {
				if msg := messageTailOf(oracleErr); !strings.Contains(omniErr.Error(), msg) {
					t.Errorf("message mismatch\n  oracle: %s\n  omni:   %s", msg, omniErr)
				}
				return
			}
			want := oracleShow(t, mc, "SHOW CREATE TABLE t")
			if got := c.ShowCreateTable("testdb", "t"); got != want {
				t.Errorf("table mismatch\n--- oracle ---\n%s\n--- omni ---\n%s", want, got)
			}
		})
	}
}

// TestReview_NumericDefaultOnAStringColumn covers a numeric literal default on
// a column that stores a value as text or bytes. The column takes the number's
// own text, so the literal is converted rather than quoted as written:
// `VARCHAR(20) DEFAULT (1e2)` is '100', where quoting the spelling made it the
// three characters 1e2.
//
// Measured on mariadb:11.8.8, with the stored value read back on each side.
//
// An integer literal too wide for an int64 is left out, and is a difference
// before the formatter sees it: the parser holds such a literal clamped, so
// `DEFAULT (12345678901234567890)` reaches here as 9223372036854775807 where
// the server stores the digits as written.
func TestReview_NumericDefaultOnAStringColumn(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, ddl string }{
		{"exponent", "CREATE TABLE t (a VARCHAR(20) DEFAULT (1e2))"},
		{"exponent_unparenthesized", "CREATE TABLE t (a VARCHAR(20) DEFAULT 1e2)"},
		{"negative_exponent", "CREATE TABLE t (a VARCHAR(20) DEFAULT (1e-3))"},
		{"exponent_with_a_mantissa", "CREATE TABLE t (a VARCHAR(40) DEFAULT (1.5e3))"},
		{"integer", "CREATE TABLE t (a VARCHAR(20) DEFAULT (100))"},
		{"fraction", "CREATE TABLE t (a VARCHAR(20) DEFAULT (1.5))"},
		{"trailing_zero_is_kept", "CREATE TABLE t (a VARCHAR(20) DEFAULT (0.10))"},
		{"leading_zeros_go", "CREATE TABLE t (a VARCHAR(20) DEFAULT (007))"},
		{"negative", "CREATE TABLE t (a VARCHAR(20) DEFAULT (-1))"},
		{"char_column", "CREATE TABLE t (a CHAR(8) DEFAULT (1e2))"},
		{"varbinary_column", "CREATE TABLE t (a VARBINARY(20) DEFAULT (1e2))"},
		{"varbinary_fraction", "CREATE TABLE t (a VARBINARY(20) DEFAULT (1.5))"},
		{"binary_column_pads", "CREATE TABLE t (a BINARY(4) DEFAULT (1e2))"},
		{"text_keeps_the_literal", "CREATE TABLE t (a TEXT DEFAULT (1e2))"},
		{"string_literal_is_untouched", "CREATE TABLE t (a VARCHAR(20) DEFAULT ('1e2'))"},
		// Past the significant positions a double prints, the server keeps the
		// literal's own spelling and so does the catalog.
		{"large_exponent", "CREATE TABLE t (a VARCHAR(20) DEFAULT (1e20))"},
		{"the_largest_expanded", "CREATE TABLE t (a VARCHAR(40) DEFAULT (1e14))"},
		{"the_smallest_expanded", "CREATE TABLE t (a VARCHAR(40) DEFAULT (1e-15))"},
		{"too_small_to_expand", "CREATE TABLE t (a VARCHAR(40) DEFAULT (1e-16))"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, tc.ddl+";")
			assertDefaultMatches(t, mc, c)

			if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
				t.Fatalf("oracle insert: %v", err)
			}
			var want string
			oracleScan(t, mc, "SELECT HEX(a) FROM t", &want)
			if _, err := mc.db.ExecContext(mc.ctx, "DROP TABLE t"); err != nil {
				t.Fatalf("oracle drop: %v", err)
			}
			replay := c.ShowCreateTable("testdb", "t")
			if _, err := mc.db.ExecContext(mc.ctx, replay); err != nil {
				t.Fatalf("oracle cannot read back what omni wrote: %v\n%s", err, replay)
			}
			if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
				t.Fatalf("oracle insert after replay: %v", err)
			}
			var got string
			oracleScan(t, mc, "SELECT HEX(a) FROM t", &got)
			if got != want {
				t.Errorf("replaying what omni wrote changed the default: %s became %s\n%s",
					want, got, replay)
			}
		})
	}
}

// TestReview_GroupByPrefersTheSourceColumn covers an alias that collides with
// a column of the FROM clause. GROUP BY resolves against the FROM clause
// first and takes the column, where ORDER BY takes the alias — so the
// substitution added for ORDER BY had to be held back here, having changed
// which rows a view groups.
//
// Measured on mariadb:11.8.8 over a table holding both a and b, with the rows
// read back so a changed grouping fails as a changed answer.
func TestReview_GroupByPrefersTheSourceColumn(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)
	runOnBoth(t, mc, c, "CREATE TABLE t (a INT, b INT);")
	if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES (1,10),(1,20)"); err != nil {
		t.Fatalf("oracle insert: %v", err)
	}

	for i, query := range []string{
		// The alias collides with a column, so GROUP BY takes the column.
		"SELECT a AS b, COUNT(*) AS n FROM t GROUP BY b",
		"SELECT a AS b, COUNT(*) AS n FROM t GROUP BY b ORDER BY b",
		"SELECT 2 AS b, COUNT(*) AS n FROM t GROUP BY b",
		"SELECT a + 1 AS b, COUNT(*) AS n FROM t GROUP BY b",
		"SELECT a AS b, COUNT(*) AS n FROM t GROUP BY t.b",
		// No collision, so the alias is taken as before.
		"SELECT a AS z, COUNT(*) AS n FROM t GROUP BY z",
		"SELECT a + 1 AS z, COUNT(*) AS n FROM t GROUP BY z",
		"SELECT 2 AS z, COUNT(*) AS n FROM t GROUP BY z",
		// ORDER BY takes the alias either way.
		"SELECT a AS b FROM t ORDER BY b",
		"SELECT 2 AS b FROM t ORDER BY b",
	} {
		t.Run(fmt.Sprintf("case%d", i), func(t *testing.T) {
			view := fmt.Sprintf("gc%d", i)
			runOnBoth(t, mc, c, "CREATE VIEW "+view+" AS "+query+";")
			assertViewBodyMatches(t, mc, c, view)

			// The rows the server gives for its own body, against the rows it
			// gives for the catalog's.
			want := oracleRows(t, mc, "SELECT * FROM "+view+" ORDER BY 1")
			body := extractSelectBody(c.ShowCreateView("testdb", view))
			got := oracleRows(t, mc, "SELECT * FROM ("+body+") AS q ORDER BY 1")
			if fmt.Sprint(got) != fmt.Sprint(want) {
				t.Errorf("replaying what omni wrote changed the rows: %v became %v\n  %s",
					want, got, body)
			}
		})
	}
}

// TestReview_ApproximateIntegerDefault covers a default written with an
// exponent, which is read as a double: a tie then rounds to its even
// neighbour, where an exact decimal rounds away from zero. Rounding the
// normalized text one way for both made `INT DEFAULT (2.5e0)` 3 where the
// server stores 2.
//
// Measured on mariadb:11.8.8 through SHOW CREATE TABLE.
func TestReview_ApproximateIntegerDefault(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, ddl string }{
		// An approximate literal rounds to even.
		{"int_two_and_a_half", "CREATE TABLE t (a INT DEFAULT (2.5e0))"},
		{"int_minus_two_and_a_half", "CREATE TABLE t (a INT DEFAULT (-2.5e0))"},
		{"int_three_and_a_half", "CREATE TABLE t (a INT DEFAULT (3.5e0))"},
		{"int_one_and_a_half", "CREATE TABLE t (a INT DEFAULT (1.5e0))"},
		{"int_half", "CREATE TABLE t (a INT DEFAULT (0.5e0))"},
		{"int_unparenthesized", "CREATE TABLE t (a INT DEFAULT 2.5e0)"},
		{"int_below_the_tie", "CREATE TABLE t (a INT DEFAULT (2.4e0))"},
		{"int_above_the_tie", "CREATE TABLE t (a INT DEFAULT (2.6e0))"},
		{"bigint", "CREATE TABLE t (a BIGINT DEFAULT (4.5e0))"},
		{"tinyint", "CREATE TABLE t (a TINYINT DEFAULT (2.5e0))"},
		{"bit", "CREATE TABLE t (a BIT(8) DEFAULT (2.5e0))"},
		{"year", "CREATE TABLE t (a YEAR DEFAULT (2.5e0))"},
		// An exact one rounds away from zero.
		{"int_exact_two_and_a_half", "CREATE TABLE t (a INT DEFAULT (2.5))"},
		{"int_exact_minus", "CREATE TABLE t (a INT DEFAULT (-2.5))"},
		{"int_exact_three_and_a_half", "CREATE TABLE t (a INT DEFAULT (3.5))"},
		// A DECIMAL column rounds away from zero either way.
		{"decimal_approximate", "CREATE TABLE t (a DECIMAL(5,0) DEFAULT (2.5e0))"},
		{"decimal_exact", "CREATE TABLE t (a DECIMAL(5,0) DEFAULT (2.5))"},
		{"decimal_scaled_approximate", "CREATE TABLE t (a DECIMAL(5,1) DEFAULT (2.25e0))"},
		{"decimal_scaled_exact", "CREATE TABLE t (a DECIMAL(5,1) DEFAULT (2.25))"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, tc.ddl+";")
			assertDefaultMatches(t, mc, c)

			if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
				t.Fatalf("oracle insert: %v", err)
			}
			var want string
			oracleScan(t, mc, "SELECT a + 0 FROM t", &want)
			if _, err := mc.db.ExecContext(mc.ctx, "DROP TABLE t"); err != nil {
				t.Fatalf("oracle drop: %v", err)
			}
			replay := c.ShowCreateTable("testdb", "t")
			if _, err := mc.db.ExecContext(mc.ctx, replay); err != nil {
				t.Fatalf("oracle cannot read back what omni wrote: %v\n%s", err, replay)
			}
			if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
				t.Fatalf("oracle insert after replay: %v", err)
			}
			var got string
			oracleScan(t, mc, "SELECT a + 0 FROM t", &got)
			if got != want {
				t.Errorf("replaying what omni wrote changed the default: %s became %s\n%s",
					want, got, replay)
			}
		})
	}
}

// TestReview_ExponentDefaultOutsideThePlainRange covers a numeric default on a
// string column whose exponent reaches past what a plain decimal holds. The
// value is still the double's, written as the double's own shortest form —
// mantissa and exponent with nothing spare — where echoing the literal's
// spelling kept whatever it was written as: (0001e15) came out '0001e15'
// instead of '1e15'.
//
// Measured on mariadb:11.8.8, with the stored characters read back from a row
// so a difference in spelling fails as a difference in value.
func TestReview_ExponentDefaultOutsideThePlainRange(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, literal string }{
		// Every one of these is the same double, and comes back '1e15'.
		{"leading_zeros", "0001e15"},
		{"canonical", "1e15"},
		{"upper_case_e", "1E15"},
		{"signed_exponent", "1e+15"},
		{"zero_padded_exponent", "1e015"},
		{"trailing_fraction_zero", "1.0e15"},
		{"mantissa_above_ten", "10e14"},
		{"mantissa_below_one", "0.1e16"},
		// A mantissa of its own is kept, normalized.
		{"with_a_mantissa", "1.5e15"},
		{"mantissa_padded", "1.50e15"},
		// Further out, in both directions.
		{"just_past", "1e16"},
		{"far_out", "1e20"},
		{"very_far_out", "1e100"},
		{"the_largest", "1e308"},
		{"small", "1e-16"},
		{"small_with_leading_zeros", "0001e-16"},
		{"small_with_a_mantissa", "1.5e-16"},
		{"very_small", "1e-100"},
		// A literal written without an exponent keeps its digits however many
		// there are.
		{"eighteen_digits", "123456789012345678"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			ddl := "CREATE TABLE t (a VARCHAR(60) DEFAULT (" + tc.literal + "))"
			runOnBoth(t, mc, c, ddl+";")
			assertDefaultMatches(t, mc, c)

			if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
				t.Fatalf("oracle insert: %v", err)
			}
			var want string
			oracleScan(t, mc, "SELECT a FROM t", &want)
			if _, err := mc.db.ExecContext(mc.ctx, "DROP TABLE t"); err != nil {
				t.Fatalf("oracle drop: %v", err)
			}
			replay := c.ShowCreateTable("testdb", "t")
			if _, err := mc.db.ExecContext(mc.ctx, replay); err != nil {
				t.Fatalf("oracle cannot read back what omni wrote: %v\n%s", err, replay)
			}
			if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
				t.Fatalf("oracle insert after replay: %v", err)
			}
			var got string
			oracleScan(t, mc, "SELECT a FROM t", &got)
			if got != want {
				t.Errorf("replaying what omni wrote changed the default: %q became %q\n%s",
					want, got, replay)
			}
		})
	}
}

// TestReview_LongExactNumericDefault covers an exact numeric default longer
// than a double holds. It is not a double, so every digit it was written with
// is kept — only the leading zeros and a leading plus go. Taking the 15
// positions that bound the exponent expansion as a limit on conversion
// altogether had left the literal's own spelling quoted, leading zeros and
// all: (000123456789012345678.50) came out '000123456789012345678.50'.
//
// Measured on mariadb:11.8.8, with the stored characters read back on both
// sides.
//
// Every literal here carries a fraction or fits an int64, because a whole
// number wider than that is clamped by the parser before the writer sees it:
// `DEFAULT (123456789012345678901)` arrives as 9223372036854775807, which no
// formatting can undo.
func TestReview_LongExactNumericDefault(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, literal string }{
		{"leading_zeros_and_a_fraction", "000123456789012345678.50"},
		{"no_leading_zeros", "123456789012345678.50"},
		{"zeros_down_to_one", "0000000000000000000001"},
		{"a_leading_plus", "+000123456789012345678.50"},
		{"negative_with_leading_zeros", "-000123456789012345678.50"},
		{"thirty_digits_with_a_fraction", "00123456789012345678901234567890.5"},
		{"twenty_one_decimals", "0.000000000000000000001"},
		{"a_long_fraction", "0.123456789012345678901"},
		{"zeros_before_a_fraction", "00.100"},
	} {
		for _, typ := range []string{"VARCHAR(60)", "VARBINARY(60)"} {
			name := tc.name
			if strings.HasPrefix(typ, "VARBINARY") {
				name += "_as_bytes"
			}
			t.Run(name, func(t *testing.T) {
				scenarioReset(t, mc)
				c := scenarioNewCatalog(t)
				ddl := "CREATE TABLE t (a " + typ + " DEFAULT (" + tc.literal + "))"
				runOnBoth(t, mc, c, ddl+";")
				assertDefaultMatches(t, mc, c)

				if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
					t.Fatalf("oracle insert: %v", err)
				}
				var want string
				oracleScan(t, mc, "SELECT HEX(a) FROM t", &want)
				if _, err := mc.db.ExecContext(mc.ctx, "DROP TABLE t"); err != nil {
					t.Fatalf("oracle drop: %v", err)
				}
				replay := c.ShowCreateTable("testdb", "t")
				if _, err := mc.db.ExecContext(mc.ctx, replay); err != nil {
					t.Fatalf("oracle cannot read back what omni wrote: %v\n%s", err, replay)
				}
				if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
					t.Fatalf("oracle insert after replay: %v", err)
				}
				var got string
				oracleScan(t, mc, "SELECT HEX(a) FROM t", &got)
				if got != want {
					t.Errorf("replaying what omni wrote changed the default: %s became %s\n%s",
						want, got, replay)
				}
			})
		}
	}
}

// TestReview_ApproximateStringDefaultIsTheDouble covers a default written with
// an exponent on a column that stores a value. The literal is a double, so the
// value is the double's and not the digits it was written with — and where the
// point sits in that value decides whether it is written plainly or with an
// exponent. Taking the expanded digits whenever they fit a plain decimal kept
// digits no double holds: (123456789012345.123456789e0) came out with all
// twenty-four.
//
// Measured on mariadb:11.8.8, with the stored characters read back on both
// sides.
func TestReview_ApproximateStringDefaultIsTheDouble(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, literal string }{
		// Digits past what a double holds are the double's, not the
		// literal's.
		{"more_digits_than_a_double_holds", "123456789012345.123456789e0"},
		{"mantissa_past_the_point", "1.23456789012345678e5"},
		{"seventeen_significant", "1234567890123456.7e0"},
		{"sixteen_significant", "123456789012345.6e0"},
		{"rounds_to_a_whole", "9.9999999999999999e0"},
		{"the_classic_tie", "2.675e0"},
		{"a_double_that_is_not_the_decimal", "0.30000000000000004e0"},
		{"just_past_one", "1.0000000000000002e0"},
		{"pi", "3.14159265358979e0"},
		// Where the point sits decides the notation.
		{"plain_at_fifteen", "1e14"},
		{"plain_at_fifteen_with_a_mantissa", "1.23e14"},
		{"exponent_at_sixteen", "1e15"},
		{"exponent_at_sixteen_with_a_mantissa", "1.23e15"},
		{"sixteen_digits_no_fraction", "1234567890123456e0"},
		{"seventeen_digits_no_fraction", "12345678901234567e0"},
		{"small_plain", "1e-15"},
		{"small_plain_with_a_mantissa", "1.5e-15"},
		{"small_exponent", "1e-16"},
		{"small_digits", "0.00012345678901234567e0"},
		{"whole", "1e0"},
		{"scaled_whole", "1.5e1"},
		// The exact twins keep every digit.
		{"exact_twin", "123456789012345.123456789"},
		{"exact_tie", "2.675"},
	} {
		for _, typ := range []string{"VARCHAR(60)", "VARBINARY(60)"} {
			name := tc.name
			if strings.HasPrefix(typ, "VARBINARY") {
				name += "_as_bytes"
			}
			t.Run(name, func(t *testing.T) {
				scenarioReset(t, mc)
				c := scenarioNewCatalog(t)
				ddl := "CREATE TABLE t (a " + typ + " DEFAULT (" + tc.literal + "))"
				runOnBoth(t, mc, c, ddl+";")
				assertDefaultMatches(t, mc, c)

				if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
					t.Fatalf("oracle insert: %v", err)
				}
				var want string
				oracleScan(t, mc, "SELECT HEX(a) FROM t", &want)
				if _, err := mc.db.ExecContext(mc.ctx, "DROP TABLE t"); err != nil {
					t.Fatalf("oracle drop: %v", err)
				}
				replay := c.ShowCreateTable("testdb", "t")
				if _, err := mc.db.ExecContext(mc.ctx, replay); err != nil {
					t.Fatalf("oracle cannot read back what omni wrote: %v\n%s", err, replay)
				}
				if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
					t.Fatalf("oracle insert after replay: %v", err)
				}
				var got string
				oracleScan(t, mc, "SELECT HEX(a) FROM t", &got)
				if got != want {
					t.Errorf("replaying what omni wrote changed the default: %s became %s\n%s",
						want, got, replay)
				}
			})
		}
	}
}

// TestReview_EveryPrimaryKeyDropSpelling covers the spellings that take the
// primary key away, which the statement-wide check has to count all of: the
// parser files DROP PRIMARY KEY and the DROP CONSTRAINT forms as a dropped
// constraint and DROP INDEX `PRIMARY` and DROP KEY `PRIMARY` as a dropped
// index, and only the first was counted — so dropping a key's column beside
// one of the index spellings was refused.
//
// Measured on mariadb:11.8.8: each of these is accepted whichever order the
// clauses come in, where the column drop alone answers 1072.
func TestReview_EveryPrimaryKeyDropSpelling(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct {
		name    string
		alter   string
		accepts bool
	}{
		{"drop_index_primary", "ALTER TABLE t DROP COLUMN a, DROP INDEX `PRIMARY`", true},
		{"drop_index_primary_first", "ALTER TABLE t DROP INDEX `PRIMARY`, DROP COLUMN a", true},
		{"drop_key_primary", "ALTER TABLE t DROP COLUMN a, DROP KEY `PRIMARY`", true},
		{"drop_primary_key", "ALTER TABLE t DROP COLUMN a, DROP PRIMARY KEY", true},
		{"drop_constraint_primary", "ALTER TABLE t DROP COLUMN a, DROP CONSTRAINT `PRIMARY`", true},
		{"the_column_alone", "ALTER TABLE t DROP COLUMN a", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, "CREATE TABLE t (a INT, b INT, c INT, PRIMARY KEY (a,b));")

			_, oracleErr := mc.db.ExecContext(mc.ctx, tc.alter)
			if tc.accepts != (oracleErr == nil) {
				t.Fatalf("oracle: %q gave %v, want accepted=%v", tc.alter, oracleErr, tc.accepts)
			}
			results, err := c.Exec(tc.alter+";", nil)
			omniErr := err
			for _, r := range results {
				if r.Error != nil {
					omniErr = r.Error
				}
			}
			if tc.accepts != (omniErr == nil) {
				t.Fatalf("omni: %q gave %v, want accepted=%v", tc.alter, omniErr, tc.accepts)
			}
			if !tc.accepts {
				if msg := messageTailOf(oracleErr); !strings.Contains(omniErr.Error(), msg) {
					t.Errorf("message mismatch\n  oracle: %s\n  omni:   %s", msg, omniErr)
				}
				return
			}
			want := oracleShow(t, mc, "SHOW CREATE TABLE t")
			if got := c.ShowCreateTable("testdb", "t"); got != want {
				t.Errorf("table mismatch\n--- oracle ---\n%s\n--- omni ---\n%s", want, got)
			}
		})
	}
}

// TestReview_NumericDefaultIsTheColumnsOwnFloat covers a default on a numeric
// column, where the value written is the column's own rather than the digits
// the literal was written with. An exponent makes the literal a double, so a
// DECIMAL rounds the double's decimal and not the literal's: DECIMAL(20,0)
// DEFAULT (9007199254740993e0) came out with all seventeen digits where the
// server keeps the double's sixteen. A scaleless float column holds a float of
// its own either way, so a DOUBLE writes its value in the double's notation
// whether or not the literal carried an exponent, and a FLOAT writes the
// float32 it holds to six significant digits.
//
// Measured on mariadb:11.8.8, with the stored value read back on both sides.
func TestReview_NumericDefaultIsTheColumnsOwnFloat(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, typ, literal string }{
		// A DECIMAL rounds the double an exponent makes of the literal.
		{"decimal_past_a_double", "DECIMAL(20,0)", "9007199254740993e0"},
		{"decimal_exact_twin", "DECIMAL(20,0)", "9007199254740993"},
		{"decimal_tie_as_a_double", "DECIMAL(5,2)", "2.675e0"},
		{"decimal_tie_exact", "DECIMAL(5,2)", "2.675"},
		{"decimal_scale_past_the_double", "DECIMAL(30,20)", "0.1e0"},
		{"decimal_long_as_a_double", "DECIMAL(30,10)", "123456789012345.123456789e0"},
		{"decimal_long_exact", "DECIMAL(30,10)", "123456789012345.123456789"},
		{"decimal_from_an_exponent", "DECIMAL(20,0)", "1e15"},
		{"decimal_from_a_big_exponent", "DECIMAL(40,0)", "1e30"},
		// A scaleless DOUBLE writes the double, in the double's notation.
		{"double_past_a_double", "DOUBLE", "9007199254740993e0"},
		{"double_exact_twin", "DOUBLE", "9007199254740993"},
		{"double_tenth", "DOUBLE", "0.1e0"},
		{"double_tie_exact", "DOUBLE", "2.675"},
		{"double_long_exact", "DOUBLE", "123456789012345.123456789"},
		{"double_at_the_exponent_edge", "DOUBLE", "1e15"},
		{"double_past_the_edge", "DOUBLE", "1e20"},
		{"double_trailing_zeros", "DOUBLE", "7.000"},
		{"double_small_exponent", "DOUBLE", "1e-16"},
		{"double_largest", "DOUBLE", "1e308"},
		{"double_with_a_scale", "DOUBLE(10,2)", "2.675e0"},
		// A scaleless FLOAT writes its float32, to six significant digits.
		{"float_more_digits_than_it_holds", "FLOAT", "0.123456789"},
		{"float_seven_whole_digits", "FLOAT", "1234567"},
		{"float_eight_whole_digits", "FLOAT", "12345678"},
		{"float_past_its_integers", "FLOAT", "16777217"},
		{"float_just_past_one", "FLOAT", "1.0000001"},
		{"float_a_third", "FLOAT", "0.3333333333"},
		{"float_pi", "FLOAT", "3.14159265358979"},
		{"float_small_plain", "FLOAT", "1e-7"},
		{"float_largest", "FLOAT", "1e38"},
		{"float_at_the_exponent_edge", "FLOAT", "1e15"},
		{"float_long_exact", "FLOAT", "123456789012345.123456789"},
		{"float_rounds_up_a_place", "FLOAT", "999999.5"},
		{"float_smallest_subnormal", "FLOAT", "1e-45"},
		{"float_short_enough_to_keep", "FLOAT", "2.675"},
		{"float_negative", "FLOAT", "-0.123456789"},
		{"float_with_a_scale", "FLOAT(10,2)", "2.675e0"},
		// An integer column takes the double too, where the literal is one.
		{"bigint_past_a_double", "BIGINT", "9007199254740993e0"},
		{"bigint_exact_twin", "BIGINT", "9007199254740993"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			ddl := "CREATE TABLE t (a " + tc.typ + " DEFAULT (" + tc.literal + "))"
			runOnBoth(t, mc, c, ddl+";")
			assertDefaultMatches(t, mc, c)

			if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
				t.Fatalf("oracle insert: %v", err)
			}
			var want string
			oracleScan(t, mc, "SELECT CAST(a AS CHAR) FROM t", &want)
			if _, err := mc.db.ExecContext(mc.ctx, "DROP TABLE t"); err != nil {
				t.Fatalf("oracle drop: %v", err)
			}
			replay := c.ShowCreateTable("testdb", "t")
			if _, err := mc.db.ExecContext(mc.ctx, replay); err != nil {
				t.Fatalf("oracle cannot read back what omni wrote: %v\n%s", err, replay)
			}
			if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
				t.Fatalf("oracle insert after replay: %v", err)
			}
			var got string
			oracleScan(t, mc, "SELECT CAST(a AS CHAR) FROM t", &got)
			if got != want {
				t.Errorf("replaying what omni wrote changed the default: %s became %s\n%s",
					want, got, replay)
			}
		})
	}
}

// TestReview_QuotedNumericDefaultIsExact covers a numeric default whose number
// comes from a string, where an exponent in the text is no sign of a double:
// the string is converted as a decimal, so DECIMAL(20,0) DEFAULT
// '9007199254740993e0' keeps all seventeen digits where the bare literal is
// the double's sixteen, and a tie rounds away from zero rather than to even.
// Reading the text for an exponent took the double's reading for every one of
// these — the quoted, X-quoted and introduced forms alike.
//
// Measured on mariadb:11.8.8, with the stored value read back on both sides.
func TestReview_QuotedNumericDefaultIsExact(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	hexOf := func(s string) string {
		return "X'" + hex.EncodeToString([]byte(s)) + "'"
	}
	for _, tc := range []struct{ name, typ, literal string }{
		// A string keeps the digits an exponent would have spent.
		{"decimal_quoted_past_a_double", "DECIMAL(20,0)", "'9007199254740993e0'"},
		{"decimal_bare_past_a_double", "DECIMAL(20,0)", "9007199254740993e0"},
		{"decimal_quoted_long", "DECIMAL(30,10)", "'123456789012345.123456789e0'"},
		{"decimal_quoted_exponent", "DECIMAL(20,0)", "'1e15'"},
		{"decimal_quoted_big_exponent", "DECIMAL(40,0)", "'1e30'"},
		{"decimal_quoted_rounds_to_a_whole", "DECIMAL(20,0)", "'9.9999999999999999e0'"},
		{"decimal_quoted_tie", "DECIMAL(5,2)", "'2.675e0'"},
		// X-quoted and introduced strings are strings too.
		{"decimal_x_quoted_past_a_double", "DECIMAL(20,0)", hexOf("9007199254740993e0")},
		{"decimal_x_quoted_tie", "DECIMAL(5,2)", hexOf("2.675e0")},
		{"decimal_introduced_past_a_double", "DECIMAL(20,0)", "_utf8mb4'9007199254740993e0'"},
		// An integer column takes the string's number, and its tie.
		{"bigint_quoted_past_a_double", "BIGINT", "'9007199254740993e0'"},
		{"bigint_quoted_exact", "BIGINT", "'9007199254740993'"},
		{"int_quoted_tie", "INT", "'2.5e0'"},
		{"int_bare_tie", "INT", "2.5e0"},
		{"int_quoted_fraction", "INT", "'7.5e0'"},
		{"int_quoted_plain_fraction", "INT", "'7.5'"},
		// A float column holds a float either way, so the string reads as one.
		{"double_quoted_past_a_double", "DOUBLE", "'9007199254740993e0'"},
		{"double_quoted_exact", "DOUBLE", "'9007199254740993'"},
		{"double_quoted_exponent", "DOUBLE", "'1e15'"},
		{"double_quoted_tenth", "DOUBLE", "'0.1'"},
		{"float_quoted_long", "FLOAT", "'0.123456789'"},
		{"float_quoted_whole", "FLOAT", "'1234567'"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			ddl := "CREATE TABLE t (a " + tc.typ + " DEFAULT (" + tc.literal + "))"
			runOnBoth(t, mc, c, ddl+";")
			assertDefaultMatches(t, mc, c)

			if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
				t.Fatalf("oracle insert: %v", err)
			}
			var want string
			oracleScan(t, mc, "SELECT CAST(a AS CHAR) FROM t", &want)
			if _, err := mc.db.ExecContext(mc.ctx, "DROP TABLE t"); err != nil {
				t.Fatalf("oracle drop: %v", err)
			}
			replay := c.ShowCreateTable("testdb", "t")
			if _, err := mc.db.ExecContext(mc.ctx, replay); err != nil {
				t.Fatalf("oracle cannot read back what omni wrote: %v\n%s", err, replay)
			}
			if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
				t.Fatalf("oracle insert after replay: %v", err)
			}
			var got string
			oracleScan(t, mc, "SELECT CAST(a AS CHAR) FROM t", &got)
			if got != want {
				t.Errorf("replaying what omni wrote changed the default: %s became %s\n%s",
					want, got, replay)
			}
		})
	}
}

// TestReview_RewritesReachWindowsAndSubqueries covers the places inside a view
// body the resolver rewrites have to reach. A window specification is one: the
// rewrite walked a call's arguments but not its OVER clause, so a SOUNDS LIKE
// partitioned on came out at comparison precedence without the parentheses its
// left operand needed, which merges partitions the original keeps apart.
//
// Measured on mariadb:11.8.8: the server rewrites a window's PARTITION BY and
// ORDER BY exactly as it rewrites a target, and so it does every query nested
// in the body — a CTE's, a derived table's, a scalar, IN, EXISTS or quantified
// subquery's — and a join's ON condition. The walk reached none of those, so
// each one kept what it was written as. A frame bound is the one position with
// nothing to rewrite: it takes no expression at all, ROWS BETWEEN (1 + 1)
// PRECEDING being 1064.
//
// Each case compares the stored body and then replays what omni wrote and
// compares the view's own rows, so a changed expression fails as changed rows.
func TestReview_RewritesReachWindowsAndSubqueries(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)

	runOnBoth(t, mc, c, `CREATE TABLE t (a INT, b INT);`)
	if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES (1,2),(2,1)"); err != nil {
		t.Fatalf("oracle insert: %v", err)
	}
	for i, tc := range []struct{ name, body string }{
		{"named_window_partition", "SELECT SUM(a) OVER w AS x FROM t WINDOW w AS (PARTITION BY (a = 1) SOUNDS LIKE '0')"},
		{"inline_window_partition", "SELECT SUM(a) OVER (PARTITION BY (a = 1) SOUNDS LIKE '0') AS x FROM t"},
		{"inline_window_order", "SELECT SUM(a) OVER (ORDER BY (a = 1) SOUNDS LIKE '0') AS x FROM t"},
		{"window_partition_not", "SELECT SUM(a) OVER (PARTITION BY NOT (a > 0)) AS x FROM t"},
		{"window_partition_and", "SELECT SUM(a) OVER (PARTITION BY a AND b) AS x FROM t"},
		{"window_order_not", "SELECT SUM(a) OVER (ORDER BY NOT (a > 0)) AS x FROM t"},
		{"window_partition_de_morgan", "SELECT SUM(a) OVER (PARTITION BY NOT (a > 0 AND b > 0)) AS x FROM t"},
		{"derived_table_sounds_like", "SELECT d.x FROM (SELECT (a = 1) SOUNDS LIKE '0' AS x FROM t) d"},
		{"derived_table_not", "SELECT d.x FROM (SELECT NOT (a > 0) AS x FROM t) d"},
		{"scalar_subquery_sounds_like", "SELECT (SELECT (a = 1) SOUNDS LIKE '0' FROM t LIMIT 1) AS x"},
		{"in_subquery_sounds_like", "SELECT a FROM t WHERE a IN (SELECT (b = 1) SOUNDS LIKE '0' FROM t)"},
		{"exists_subquery_sounds_like", "SELECT a FROM t WHERE EXISTS (SELECT 1 FROM t t2 WHERE (t2.a = 1) SOUNDS LIKE '0')"},
		{"quantified_subquery_sounds_like", "SELECT a FROM t WHERE a = ANY (SELECT (b = 1) SOUNDS LIKE '0' FROM t)"},
		{"cte_body_sounds_like", "WITH cte AS (SELECT (a = 1) SOUNDS LIKE '0' AS x FROM t) SELECT cte.x FROM cte"},
		{"join_on_sounds_like", "SELECT t1.a FROM t t1 JOIN t t2 ON (t1.a = 1) SOUNDS LIKE '0'"},
		{"join_on_not", "SELECT t1.a FROM t t1 JOIN t t2 ON NOT (t1.a > 0)"},
		{"derived_table_in_a_join", "SELECT d.x FROM t t1 JOIN (SELECT NOT (a > 0) AS x FROM t) d ON d.x = t1.a"},
		{"window_in_a_derived_table", "SELECT d.x FROM (SELECT SUM(a) OVER (PARTITION BY NOT (a > 0)) AS x FROM t) d"},
	} {
		view := fmt.Sprintf("v_rw%d", i)
		t.Run(tc.name, func(t *testing.T) {
			runOnBoth(t, mc, c, "CREATE VIEW "+view+" AS "+tc.body+";")
			assertViewBodyMatches(t, mc, c, view)

			want := fmt.Sprint(oracleRows(t, mc, "SELECT * FROM "+view+" ORDER BY 1"))
			if _, err := mc.db.ExecContext(mc.ctx, "DROP VIEW "+view); err != nil {
				t.Fatalf("oracle drop: %v", err)
			}
			replay := c.ShowCreateView("testdb", view)
			if _, err := mc.db.ExecContext(mc.ctx, replay); err != nil {
				t.Fatalf("oracle cannot read back what omni wrote: %v\n%s", err, replay)
			}
			if got := fmt.Sprint(oracleRows(t, mc, "SELECT * FROM "+view+" ORDER BY 1")); got != want {
				t.Errorf("replaying what omni wrote changed the view's rows: %s became %s\n%s",
					want, got, replay)
			}
		})
	}
}

// TestReview_RewritesReachEveryChildExpression covers the expression nodes the
// rewrite walked past, each of which holds a child the server rewrites. A row
// was the costly one: its items are separated by commas alone, so an item left
// at comparison precedence reads back as the comparison and nothing else —
// ROW((2=1) SOUNDS LIKE '0',1) = ROW(1,1) came out as (2 = 1 sounds like
// '0',1) = (1,1), which is 0 where the original is 1.
//
// Measured on mariadb:11.8.8: a row's items, an INTERVAL's value, EXTRACT's
// operand, MATCH's AGAINST argument and an aggregate's own ORDER BY are each
// rewritten as an expression of its own. A separator is left alone because it
// can only be a string literal, which no rewrite touches.
//
// MATCH is rewritten the same way — `AGAINST (((2=1) SOUNDS LIKE '0'))` is
// stored as against (soundex(2 = 1) = soundex('0')) — but no case here can
// pin it: the deparser writes no MATCH at all yet, and the server's own
// stored body for that view is one SHOW CREATE VIEW answers 1064 on, since
// AGAINST takes only a literal. Both are their own gaps, reported.
func TestReview_RewritesReachEveryChildExpression(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)

	runOnBoth(t, mc, c, `CREATE TABLE t (a INT, b INT);`)
	if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES (1,2),(2,1)"); err != nil {
		t.Fatalf("oracle insert: %v", err)
	}
	for i, tc := range []struct{ name, body string }{
		{"row_item_sounds_like", "SELECT ROW((2=1) SOUNDS LIKE '0',1) = ROW(1,1) AS x"},
		{"row_item_not", "SELECT ROW(NOT (a > 0),1) = ROW(1,1) AS x FROM t"},
		{"row_on_the_right", "SELECT (a,b) = ((a = 1) SOUNDS LIKE '0', 1) AS x FROM t"},
		{"interval_value_sounds_like", "SELECT '2020-01-01' + INTERVAL ((2=1) SOUNDS LIKE '0') DAY AS x"},
		{"interval_value_not", "SELECT '2020-01-01' + INTERVAL (NOT (a > 0)) DAY AS x FROM t"},
		{"extract_operand", "SELECT EXTRACT(DAY FROM ((2=1) SOUNDS LIKE '0')) AS x"},
		{"aggregate_order_by_sounds_like", "SELECT GROUP_CONCAT(a ORDER BY (a = 1) SOUNDS LIKE '0') AS x FROM t"},
		{"aggregate_order_by_not", "SELECT GROUP_CONCAT(a ORDER BY NOT (a > 0)) AS x FROM t"},
	} {
		view := fmt.Sprintf("v_ch%d", i)
		t.Run(tc.name, func(t *testing.T) {
			runOnBoth(t, mc, c, "CREATE VIEW "+view+" AS "+tc.body+";")
			assertViewBodyMatches(t, mc, c, view)

			want := fmt.Sprint(oracleRows(t, mc, "SELECT * FROM "+view+" ORDER BY 1"))
			if _, err := mc.db.ExecContext(mc.ctx, "DROP VIEW "+view); err != nil {
				t.Fatalf("oracle drop: %v", err)
			}
			replay := c.ShowCreateView("testdb", view)
			if _, err := mc.db.ExecContext(mc.ctx, replay); err != nil {
				t.Fatalf("oracle cannot read back what omni wrote: %v\n%s", err, replay)
			}
			if got := fmt.Sprint(oracleRows(t, mc, "SELECT * FROM "+view+" ORDER BY 1")); got != want {
				t.Errorf("replaying what omni wrote changed the view's rows: %s became %s\n%s",
					want, got, replay)
			}
		})
	}
}

// TestReview_NumericTemporalDefault covers a default written as a number on a
// temporal column. The column holds the value the number converts to, and
// SHOW CREATE writes that value's own literal; omni quoted the digits as they
// had been written, so `TIME DEFAULT (1e2)` came out DEFAULT '1e2', which the
// server refuses outright on replay because '1e2' is no TIME.
//
// Measured on mariadb:11.8.8. A number is read from the right for a TIME and
// through the server's own two-digit-year ranges for a date, which is why a
// number of an unexpected length lands where it does — (20200) is
// '2002-02-00', the two-digit-year 020200, and (2020010112345) is
// '0202-00-10 11:23:45', the thirteen digits as they stand. A fraction is
// truncated to the precision the column declares, not rounded.
//
// Each case compares the DEFAULT clause and replays it, reading the stored
// value back from a row, so a changed value fails as a changed row.
func TestReview_NumericTemporalDefault(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, typ, literal string }{
		// TIME reads the number from the right.
		{"time_from_an_exponent", "TIME", "1e2"},
		{"time_seconds", "TIME", "7"},
		{"time_one_minute", "TIME", "100"},
		{"time_one_hour", "TIME", "10000"},
		{"time_widest", "TIME", "8385959"},
		{"time_day_boundary", "TIME", "240000"},
		{"time_zero", "TIME", "0"},
		{"time_negative_zero", "TIME", "-0"},
		{"time_negative", "TIME", "-100"},
		{"time_negative_widest", "TIME", "-8385959"},
		{"time_fraction_dropped", "TIME", "100.999"},
		{"time_mantissa", "TIME", "1.5e2"},
		{"time_from_a_date_and_time", "TIME", "20200101123456"},
		{"time_from_the_widest_date", "TIME", "99991231235959"},
		{"time_from_a_date_with_a_fraction", "TIME", "20200101123456.5"},
		// A declared precision keeps the fraction, truncated to it.
		{"time_fsp3", "TIME(3)", "100.5"},
		{"time_fsp6", "TIME(6)", "100.5"},
		{"time_fsp3_truncates", "TIME(3)", "100.4567"},
		{"time_fsp6_pads", "TIME(6)", "100.45"},
		{"time_fsp3_negative", "TIME(3)", "-100.5"},
		{"time_fsp3_from_a_date", "TIME(3)", "20200101123456.5"},
		// DATE takes the date alone.
		{"date_from_an_exponent", "DATE", "2.0200101e7"},
		{"date_eight_digits", "DATE", "20200101"},
		{"date_two_digit_year", "DATE", "200101"},
		{"date_short", "DATE", "101"},
		{"date_year_window_2069", "DATE", "691231"},
		{"date_year_window_1970", "DATE", "700101"},
		{"date_year_window_1999", "DATE", "991231"},
		{"date_widest", "DATE", "99991231"},
		{"date_zero", "DATE", "0"},
		{"date_leap_day", "DATE", "20200229"},
		{"date_zero_day", "DATE", "20200100"},
		{"date_zero_month", "DATE", "20200001"},
		{"date_drops_the_time", "DATE", "20200101123456"},
		{"date_drops_the_fraction", "DATE", "20200101.5"},
		{"date_five_digits", "DATE", "20200"},
		{"date_four_digits", "DATE", "1231"},
		{"date_nine_digits", "DATE", "202001011"},
		{"date_eleven_digits", "DATE", "20200101123"},
		// DATETIME and TIMESTAMP take both.
		{"datetime_fourteen_digits", "DATETIME", "20200101123456"},
		{"datetime_eight_digits", "DATETIME", "20200101"},
		{"datetime_six_digits", "DATETIME", "200101"},
		{"datetime_twelve_digits", "DATETIME", "200101123456"},
		{"datetime_thirteen_digits", "DATETIME", "2020010112345"},
		{"datetime_nine_digits", "DATETIME", "202001011"},
		{"datetime_zero", "DATETIME", "0"},
		{"datetime_widest", "DATETIME", "99991231235959"},
		{"datetime_from_an_exponent", "DATETIME", "2.0200101e13"},
		{"datetime_fraction_dropped", "DATETIME", "20200101123456.456789"},
		{"datetime_fsp3_zero", "DATETIME(3)", "0"},
		{"datetime_fsp3", "DATETIME(3)", "20200101123456.5"},
		{"datetime_fsp3_truncates", "DATETIME(3)", "20200101123456.456789"},
		{"datetime_fsp6", "DATETIME(6)", "20200101123456.5"},
		{"timestamp_fourteen_digits", "TIMESTAMP", "20200101123456"},
		{"timestamp_fraction_dropped", "TIMESTAMP", "20200101123456.5"},
		{"timestamp_fsp3", "TIMESTAMP(3)", "20200101123456.5"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			ddl := "CREATE TABLE t (a " + tc.typ + " DEFAULT (" + tc.literal + "))"
			runOnBoth(t, mc, c, ddl+";")
			assertDefaultMatches(t, mc, c)

			if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
				t.Fatalf("oracle insert: %v", err)
			}
			var want string
			oracleScan(t, mc, "SELECT CAST(a AS CHAR) FROM t", &want)
			if _, err := mc.db.ExecContext(mc.ctx, "DROP TABLE t"); err != nil {
				t.Fatalf("oracle drop: %v", err)
			}
			replay := c.ShowCreateTable("testdb", "t")
			if _, err := mc.db.ExecContext(mc.ctx, replay); err != nil {
				t.Fatalf("oracle cannot read back what omni wrote: %v\n%s", err, replay)
			}
			if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
				t.Fatalf("oracle insert after replay: %v", err)
			}
			var got string
			oracleScan(t, mc, "SELECT CAST(a AS CHAR) FROM t", &got)
			if got != want {
				t.Errorf("replaying what omni wrote changed the default: %s became %s\n%s",
					want, got, replay)
			}
		})
	}
}

// TestReview_InlineReferenceBackingIndexName covers the index a column-level
// REFERENCES brings with it. MariaDB honours such a reference, and the index
// is named after the constraint only where the statement named one: an unnamed
// reference's index takes the column's name while the constraint takes the
// generated symbol. omni gave the index that symbol, so a statement with a key
// of its own by that name ended up with two indexes called `t_ibfk_1` and a
// SHOW CREATE TABLE the server could not replay.
//
// The place the index takes is the statement's own order, which these pin too:
// measured, `FOREIGN KEY (a) REFERENCES p(id), KEY z(b)` shows KEY `a` then
// KEY `z` and `KEY z(b), FOREIGN KEY (a)` the other way round, where omni had
// put every backing index last.
//
// Measured on mariadb:11.8.8.
//
// Two neighbours are left out. A name an inline reference's index would
// duplicate is 1061 on the server — `a INT REFERENCES p(id), b INT, KEY a(b)`
// is refused where the table-level `KEY a(b), FOREIGN KEY (a)` allocates
// `a_2` — and omni, which gives every backing index its place only once the
// statement is read, accepts both. And `a INT CONSTRAINT fk1 REFERENCES
// p(id)`, which the server takes and names the index fk1 for, omni's parser
// answers 1064 on; the name is carried here for when it learns the spelling.
// Both are their own gaps, reported.
func TestReview_InlineReferenceBackingIndexName(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, ddl string }{
		{"unnamed_beside_a_key_of_that_symbol",
			"CREATE TABLE t (a INT REFERENCES p(id), b INT, KEY t_ibfk_1(b))"},
		{"unnamed_alone", "CREATE TABLE t (a INT REFERENCES p(id))"},
		{"two_unnamed", "CREATE TABLE t (a INT REFERENCES p(id), b INT REFERENCES p(id))"},
		{"with_an_action", "CREATE TABLE t (a INT REFERENCES p(id) ON DELETE CASCADE)"},
		{"covered_by_a_key_of_its_own", "CREATE TABLE t (a INT REFERENCES p(id), KEY a(a))"},
		{"table_level_unnamed", "CREATE TABLE t (a INT, FOREIGN KEY (a) REFERENCES p(id))"},
		{"key_after_an_inline_reference",
			"CREATE TABLE t (a INT REFERENCES p(id), b INT, KEY z(b))"},
		{"key_after_a_table_level_key",
			"CREATE TABLE t (a INT, b INT, FOREIGN KEY (a) REFERENCES p(id), KEY z(b))"},
		{"key_before_a_table_level_key",
			"CREATE TABLE t (a INT, b INT, KEY z(b), FOREIGN KEY (a) REFERENCES p(id))"},
		{"named_table_level_after_a_key",
			"CREATE TABLE t (a INT, b INT, KEY z(b), CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES p(id))"},
		{"table_level_beside_that_symbol",
			"CREATE TABLE t (a INT, b INT, KEY t_ibfk_1(b), FOREIGN KEY (a) REFERENCES p(id))"},
		{"table_level_named", "CREATE TABLE t (a INT, CONSTRAINT fk1 FOREIGN KEY (a) REFERENCES p(id))"},
		{"table_level_column_name_taken",
			"CREATE TABLE t (a INT, b INT, KEY a(b), FOREIGN KEY (a) REFERENCES p(id))"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, "CREATE TABLE p (id INT PRIMARY KEY);")
			runOnBoth(t, mc, c, tc.ddl+";")

			want := oracleShow(t, mc, "SHOW CREATE TABLE t")
			if got := c.ShowCreateTable("testdb", "t"); got != want {
				t.Errorf("table mismatch\n--- oracle ---\n%s\n--- omni ---\n%s", want, got)
			}

			// And what omni wrote goes back in, which two indexes of one name
			// did not.
			if _, err := mc.db.ExecContext(mc.ctx, "DROP TABLE t"); err != nil {
				t.Fatalf("oracle drop: %v", err)
			}
			replay := c.ShowCreateTable("testdb", "t")
			if _, err := mc.db.ExecContext(mc.ctx, replay); err != nil {
				t.Fatalf("oracle cannot read back what omni wrote: %v\n%s", err, replay)
			}
			if got := oracleShow(t, mc, "SHOW CREATE TABLE t"); got != want {
				t.Errorf("replaying what omni wrote changed the table\n--- was ---\n%s\n--- now ---\n%s",
					want, got)
			}
		})
	}
}

// TestReview_ApproximateTemporalDefault covers a temporal default written with
// an exponent. The literal is a double, so the value the column converts is
// the double's: DATETIME(6) DEFAULT (20200101123456.123456e0) is
// '2020-01-01 12:34:56.125000', where omni expanded the digits as written and
// kept .123456, a value no double holds.
//
// Measured on mariadb:11.8.8. The value read is the double's exact one and not
// its shortest decimal — (20200101123456.1e0) is '...56.101562', the double
// 20200101123456.1015625 cut at six places, where the shortest decimal that
// reads back as that double is ...56.102 — and the places the column does not
// keep are cut rather than rounded, which TIME(6) DEFAULT (8385959.999999e0)
// coming back '838:59:59.999998' shows.
func TestReview_ApproximateTemporalDefault(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, typ, literal string }{
		{"datetime_fsp6_past_a_double", "DATETIME(6)", "20200101123456.123456e0"},
		{"datetime_fsp6_exact_twin", "DATETIME(6)", "20200101123456.123456"},
		{"datetime_fsp3_past_a_double", "DATETIME(3)", "20200101123456.123e0"},
		{"datetime_fsp6_one_place", "DATETIME(6)", "20200101123456.1e0"},
		{"datetime_from_a_mantissa", "DATETIME", "2.0200101123456e13"},
		{"datetime_whole_exponent", "DATETIME", "20200101123456e0"},
		{"date_from_a_mantissa", "DATE", "2.0200101e7"},
		{"date_whole_exponent", "DATE", "20200101e0"},
		{"time_fsp6_past_a_double", "TIME(6)", "123456.123456e0"},
		{"time_fsp6_exact_twin", "TIME(6)", "123456.123456"},
		{"time_from_a_mantissa", "TIME", "1.005e2"},
		{"time_fsp3_half", "TIME(3)", "100.5e0"},
		{"time_fsp3_negative", "TIME(3)", "-100.5e0"},
		{"time_fsp6_widest", "TIME(6)", "8385959.999999e0"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			ddl := "CREATE TABLE t (a " + tc.typ + " DEFAULT (" + tc.literal + "))"
			runOnBoth(t, mc, c, ddl+";")
			assertDefaultMatches(t, mc, c)

			if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
				t.Fatalf("oracle insert: %v", err)
			}
			var want string
			oracleScan(t, mc, "SELECT CAST(a AS CHAR) FROM t", &want)
			if _, err := mc.db.ExecContext(mc.ctx, "DROP TABLE t"); err != nil {
				t.Fatalf("oracle drop: %v", err)
			}
			replay := c.ShowCreateTable("testdb", "t")
			if _, err := mc.db.ExecContext(mc.ctx, replay); err != nil {
				t.Fatalf("oracle cannot read back what omni wrote: %v\n%s", err, replay)
			}
			if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
				t.Fatalf("oracle insert after replay: %v", err)
			}
			var got string
			oracleScan(t, mc, "SELECT CAST(a AS CHAR) FROM t", &got)
			if got != want {
				t.Errorf("replaying what omni wrote changed the default: %s became %s\n%s",
					want, got, replay)
			}
		})
	}
}

// TestReview_SetOperationOrderByIsRewritten covers the ORDER BY over a whole
// set operation, which belongs to neither side of it. The rewrite walked both
// operands and returned, so that list was left as written and a SOUNDS LIKE
// ordered by lost the parentheses its left operand needs.
//
// Measured on mariadb:11.8.8, where such a list is rewritten like any other.
//
// The same list over a parenthesized query is rewritten too, and is left out
// here for a neighbour it runs into: the server folds it into that query and
// qualifies its columns — `(SELECT a FROM t) ORDER BY ((a=2) SOUNDS LIKE '0')`
// is stored as `(select `t`.`a` AS `a` from `t` order by soundex(`t`.`a` = 2)
// = soundex('0'))` — where omni keeps the outer shape and the bare name. That
// is its own gap, reported.
func TestReview_SetOperationOrderByIsRewritten(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)

	runOnBoth(t, mc, c, `CREATE TABLE t (a INT, b INT);`)
	if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES (1,2),(2,1)"); err != nil {
		t.Fatalf("oracle insert: %v", err)
	}
	for i, tc := range []struct{ name, body string }{
		{"union_sounds_like", "SELECT 2 AS a UNION SELECT 3 ORDER BY ((a=2) SOUNDS LIKE '0'), a DESC"},
		{"union_not", "SELECT 2 AS a UNION SELECT 3 ORDER BY NOT (a > 2)"},
		{"union_over_tables", "SELECT a FROM t UNION SELECT b FROM t ORDER BY ((a=2) SOUNDS LIKE '0')"},
		{"union_all_with_a_limit", "SELECT a FROM t UNION ALL SELECT b FROM t ORDER BY NOT (a > 2) LIMIT 5"},
	} {
		view := fmt.Sprintf("v_so%d", i)
		t.Run(tc.name, func(t *testing.T) {
			runOnBoth(t, mc, c, "CREATE VIEW "+view+" AS "+tc.body+";")
			assertViewBodyMatches(t, mc, c, view)

			want := fmt.Sprint(oracleRows(t, mc, "SELECT * FROM "+view))
			if _, err := mc.db.ExecContext(mc.ctx, "DROP VIEW "+view); err != nil {
				t.Fatalf("oracle drop: %v", err)
			}
			replay := c.ShowCreateView("testdb", view)
			if _, err := mc.db.ExecContext(mc.ctx, replay); err != nil {
				t.Fatalf("oracle cannot read back what omni wrote: %v\n%s", err, replay)
			}
			if got := fmt.Sprint(oracleRows(t, mc, "SELECT * FROM "+view)); got != want {
				t.Errorf("replaying what omni wrote changed the view's rows: %s became %s\n%s",
					want, got, replay)
			}
		})
	}
}

// TestReview_ConstantOrderingItems covers a constant written straight into an
// ordering clause. GROUP BY and ORDER BY read a whole number as an output
// position, and the server reads a boolean there as one too: measured on
// mariadb:11.8.8, `GROUP BY TRUE` is stored as group by 1 and groups by the
// first column, exactly as `GROUP BY 1` does, while `GROUP BY FALSE`,
// `ORDER BY FALSE`, `ORDER BY -1` and `ORDER BY NOT TRUE` are all refused with
// 1054, the position 0 or -1 being no column. A constant that cannot be read
// as a position — 2.5, NULL, 1e2, X'31', b'1', 0x31 — is kept as it stands.
//
// Neither of a window specification's lists reads a number as a position, so
// there the same constant is written as the empty string an item taking a
// whole-number constant from an alias is written as — ORDER BY and PARTITION
// BY alike, and through a named WINDOW clause too. omni wrote the number,
// which is the one gap here. An aggregate's own ORDER BY keeps its constant as
// written.
//
// `GROUP_CONCAT(a ORDER BY -1)` is left out: both sides store it as
// order by -1, and selecting from the view the server itself stored answers
// 1054, so there is nothing a row comparison can say.
func TestReview_ConstantOrderingItems(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)

	runOnBoth(t, mc, c, `CREATE TABLE t (a INT, b INT);`)
	if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES (1,2),(2,1)"); err != nil {
		t.Fatalf("oracle insert: %v", err)
	}
	for i, tc := range []struct{ name, body string }{
		// A boolean is a position, as a plain number is.
		{"group_by_true", "SELECT a, COUNT(*) AS n FROM t GROUP BY TRUE"},
		{"group_by_one", "SELECT a, COUNT(*) AS n FROM t GROUP BY 1"},
		{"order_by_true", "SELECT a, b FROM t ORDER BY TRUE"},
		{"order_by_parenthesized_true", "SELECT a, b FROM t ORDER BY (TRUE)"},
		{"order_by_signed_position", "SELECT a, b FROM t ORDER BY +2"},
		{"order_by_true_then_a_column", "SELECT a, b FROM t ORDER BY TRUE, a DESC"},
		{"group_by_a_column_and_true", "SELECT a, b FROM t GROUP BY a, TRUE"},
		// A constant that is no position stays as written.
		{"order_by_a_decimal", "SELECT a, b FROM t ORDER BY 2.5"},
		{"order_by_null", "SELECT a, b FROM t ORDER BY NULL"},
		{"order_by_an_exponent", "SELECT a, b FROM t ORDER BY 1e2"},
		{"order_by_x_quoted", "SELECT a, b FROM t ORDER BY X'31'"},
		{"order_by_a_bit_literal", "SELECT a, b FROM t ORDER BY b'1'"},
		{"order_by_hex", "SELECT a, b FROM t ORDER BY 0x31"},
		// A window's ordering is not positional; an aggregate's is.
		{"window_order_by_true", "SELECT SUM(a) OVER (ORDER BY TRUE) AS x FROM t"},
		{"window_order_by_one", "SELECT SUM(a) OVER (ORDER BY 1) AS x FROM t"},
		{"window_order_by_negative_one", "SELECT SUM(a) OVER (ORDER BY -1) AS x FROM t"},
		{"window_order_by_one_desc", "SELECT SUM(a) OVER (ORDER BY 1 DESC) AS x FROM t"},
		{"window_order_by_one_and_a_column", "SELECT SUM(a) OVER (ORDER BY 1, a) AS x FROM t"},
		{"window_partition_by_true", "SELECT SUM(a) OVER (PARTITION BY TRUE) AS x FROM t"},
		{"window_partition_by_one", "SELECT SUM(a) OVER (PARTITION BY 1) AS x FROM t"},
		{"window_order_by_a_decimal", "SELECT SUM(a) OVER (ORDER BY 2.5) AS x FROM t"},
		{"window_order_by_null", "SELECT SUM(a) OVER (ORDER BY NULL) AS x FROM t"},
		{"window_order_by_an_exponent", "SELECT SUM(a) OVER (ORDER BY 1e2) AS x FROM t"},
		{"window_order_by_x_quoted", "SELECT SUM(a) OVER (ORDER BY X'31') AS x FROM t"},
		{"named_window_order_by_true", "SELECT SUM(a) OVER w AS x FROM t WINDOW w AS (ORDER BY TRUE)"},
		{"named_window_partition_by_one", "SELECT SUM(a) OVER w AS x FROM t WINDOW w AS (PARTITION BY 1)"},
		{"aggregate_order_by_true", "SELECT GROUP_CONCAT(a ORDER BY TRUE) AS x FROM t"},
		{"aggregate_order_by_a_decimal", "SELECT GROUP_CONCAT(a ORDER BY 2.5) AS x FROM t"},
		// The alias an item takes a whole-number constant from, already held.
		{"order_by_a_constant_alias", "SELECT 2 AS z FROM t ORDER BY z"},
	} {
		view := fmt.Sprintf("v_co%d", i)
		t.Run(tc.name, func(t *testing.T) {
			runOnBoth(t, mc, c, "CREATE VIEW "+view+" AS "+tc.body+";")
			assertViewBodyMatches(t, mc, c, view)

			want := fmt.Sprint(oracleRows(t, mc, "SELECT * FROM "+view))
			if _, err := mc.db.ExecContext(mc.ctx, "DROP VIEW "+view); err != nil {
				t.Fatalf("oracle drop: %v", err)
			}
			replay := c.ShowCreateView("testdb", view)
			if _, err := mc.db.ExecContext(mc.ctx, replay); err != nil {
				t.Fatalf("oracle cannot read back what omni wrote: %v\n%s", err, replay)
			}
			if got := fmt.Sprint(oracleRows(t, mc, "SELECT * FROM "+view)); got != want {
				t.Errorf("replaying what omni wrote changed the view's rows: %s became %s\n%s",
					want, got, replay)
			}
		})
	}
}

// TestReview_AliasReferenceInOrderingClauses covers an ORDER BY or GROUP BY
// item that mentions a target-list alias inside a larger expression, and what
// a view's own column list does to it. omni kept the alias everywhere, which a
// view naming its columns cannot even be replayed from: `CREATE VIEW v (x) AS
// SELECT a AS z FROM t ORDER BY z + 0` stored an order by `z` that names
// nothing.
//
// Measured on mariadb:11.8.8, and the three clauses each have their own rule.
// A GROUP BY item takes the target's expression whatever it is — over
// `a + 1 AS z` it is group by `t`.`a` + 1 + 0 and over `2.5 AS z`
// group by 2.5 + 0. An ORDER BY item takes it only where the target is a
// column, `a AS z` giving order by `t`.`a` + 0, and keeps the alias for
// anything else. HAVING keeps the output name and so takes the one the column
// list declares, coming back having `x` > 0.
func TestReview_AliasReferenceInOrderingClauses(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)

	runOnBoth(t, mc, c, `CREATE TABLE t (a INT, b INT);`)
	if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES (1,2),(2,1)"); err != nil {
		t.Fatalf("oracle insert: %v", err)
	}
	for i, tc := range []struct{ name, cols, body string }{
		{"renamed_order_by_an_expression", " (x)", "SELECT a AS z FROM t ORDER BY z + 0"},
		{"order_by_over_an_expression_target", "", "SELECT a + 1 AS z FROM t ORDER BY z + 0"},
		{"order_by_over_a_constant_target", "", "SELECT 2 AS z FROM t ORDER BY z + 0"},
		{"order_by_over_a_call_target", "", "SELECT ABS(a) AS z FROM t ORDER BY z + 0"},
		{"order_by_a_qualified_column_target", "", "SELECT t.a AS z FROM t ORDER BY z + 0"},
		{"order_by_negating_the_alias", "", "SELECT a AS z FROM t ORDER BY -z"},
		{"group_by_over_an_expression_target", "",
			"SELECT a + 1 AS z, COUNT(*) AS n FROM t GROUP BY z + 0"},
		{"group_by_over_a_constant_target", "",
			"SELECT 2.5 AS z, COUNT(*) AS n FROM t GROUP BY z + 0"},
		{"group_by_a_call_over_the_alias", "",
			"SELECT a AS z, COUNT(*) AS n FROM t GROUP BY ABS(z)"},
		{"group_by_a_call_over_a_constant_target", "",
			"SELECT 2 AS z, COUNT(*) AS n FROM t GROUP BY ABS(z)"},
		{"order_by_an_expression", "", "SELECT a AS z FROM t ORDER BY z + 0"},
		{"renamed_to_the_same_name", " (z)", "SELECT a AS z FROM t ORDER BY z + 0"},
		{"renamed_order_by_a_call", " (x)", "SELECT a AS z FROM t ORDER BY ABS(z)"},
		{"order_by_a_call", "", "SELECT a AS z FROM t ORDER BY ABS(z)"},
		{"renamed_order_by_the_alias", " (x)", "SELECT a AS z FROM t ORDER BY z"},
		{"renamed_group_by_an_expression", " (x, y)",
			"SELECT a AS z, COUNT(*) AS n FROM t GROUP BY z + 0"},
		{"group_by_an_expression", "",
			"SELECT a AS z, COUNT(*) AS n FROM t GROUP BY z + 0"},
		{"renamed_having", " (x)", "SELECT a AS z FROM t HAVING z > 0"},
		{"having", "", "SELECT a AS z FROM t HAVING z > 0"},
		{"renamed_order_by_a_constant_alias", " (x)", "SELECT 2 AS z FROM t ORDER BY z"},
		{"a_nested_query_keeps_its_own_alias", " (x)",
			"SELECT a AS z FROM t HAVING EXISTS(SELECT 1 FROM (SELECT 1 AS z HAVING z=1) AS d)"},
		{"renamed_beside_a_nested_alias", " (x)",
			"SELECT a AS z FROM t HAVING z > 0 AND EXISTS(SELECT 1 FROM (SELECT 1 AS z HAVING z=1) AS d)"},
		{"a_nested_query_keeps_its_own_having", " (x)",
			"SELECT a AS z FROM t HAVING (SELECT MAX(b) FROM t t2 HAVING MAX(b) > 0) > 0"},
		{"cte_columns_rename_the_having", "",
			"WITH c(x) AS (SELECT a AS z FROM t HAVING z > 0) SELECT c.x FROM c"},
		{"cte_columns_and_an_ordering_item", "",
			"WITH c(x) AS (SELECT a AS z FROM t ORDER BY z + 0) SELECT c.x FROM c"},
		{"cte_columns_and_a_group_by_item", "",
			"WITH c(x, y) AS (SELECT a AS z, COUNT(*) AS n FROM t GROUP BY z + 0) SELECT c.x FROM c"},
	} {
		view := fmt.Sprintf("v_al%d", i)
		t.Run(tc.name, func(t *testing.T) {
			runOnBoth(t, mc, c, "CREATE VIEW "+view+tc.cols+" AS "+tc.body+";")
			assertViewBodyMatches(t, mc, c, view)

			want := fmt.Sprint(oracleRows(t, mc, "SELECT * FROM "+view))
			if _, err := mc.db.ExecContext(mc.ctx, "DROP VIEW "+view); err != nil {
				t.Fatalf("oracle drop: %v", err)
			}
			replay := c.ShowCreateView("testdb", view)
			if _, err := mc.db.ExecContext(mc.ctx, replay); err != nil {
				t.Fatalf("oracle cannot read back what omni wrote: %v\n%s", err, replay)
			}
			if got := fmt.Sprint(oracleRows(t, mc, "SELECT * FROM "+view)); got != want {
				t.Errorf("replaying what omni wrote changed the view's rows: %s became %s\n%s",
					want, got, replay)
			}
		})
	}
}

// TestReview_EnumNumberDefault covers a default written as a number on an ENUM
// or SET column, where a number is a choice among the members rather than
// text. An ENUM takes it as a member's position and a SET as a bitmask of
// them; omni wrote the digits as a string, which names the member of that
// value instead — `ENUM('2','1') DEFAULT (1)` came out '1', the second
// member, where the server stores '2', the first.
//
// Measured on mariadb:11.8.8, with the stored value read back on both sides.
// The characters of a number are the other reading, and both are pinned here:
// '1', X'31' and 0x31 all name the member '1'.
//
// The cases the server refuses are left out, since it answers 1067 and omni
// does not check such a default yet: a position of zero or past the last
// member, a SET bitmask past them, a negative number, and any fraction,
// approximate or not. Its acceptance is stranger than the value it stores —
// the default is taken only where the number's digits are themselves a member,
// so `ENUM('a','b') DEFAULT (1)` is refused where `ENUM('a','1') DEFAULT (1)`
// is 'a' — but nothing omni writes turns on that, a statement it is given
// having been accepted already.
func TestReview_EnumNumberDefault(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, typ, literal string }{
		{"enum_first_member", "ENUM('2','1')", "1"},
		{"enum_second_member", "ENUM('2','1')", "2"},
		{"enum_from_an_exponent", "ENUM('2','1')", "1e0"},
		{"enum_third_member", "ENUM('3','b','c')", "3"},
		{"enum_member_of_its_own_digits", "ENUM('1','2','3')", "3"},
		{"enum_first_of_three", "ENUM('x','2','1')", "1"},
		{"enum_digits_that_are_a_member", "ENUM('a','1')", "1"},
		{"set_first_member", "SET('2','1')", "1"},
		{"set_second_member", "SET('2','1')", "2"},
		{"set_two_members", "SET('1','2','3')", "3"},
		{"set_three_members", "SET('7','b','c')", "7"},
		{"set_from_an_exponent", "SET('2','1')", "1e0"},
		{"set_digits_that_are_a_member", "SET('a','1')", "1"},
		// The characters of a number name the member of that value.
		{"enum_quoted_digits", "ENUM('2','1')", "'1'"},
		{"enum_x_quoted_digits", "ENUM('2','1')", "X'31'"},
		{"enum_hex_digits", "ENUM('2','1')", "0x31"},
		{"enum_quoted_member", "ENUM('2','1')", "'2'"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			ddl := "CREATE TABLE t (a " + tc.typ + " DEFAULT (" + tc.literal + "))"
			runOnBoth(t, mc, c, ddl+";")
			assertDefaultMatches(t, mc, c)

			if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
				t.Fatalf("oracle insert: %v", err)
			}
			var want string
			oracleScan(t, mc, "SELECT CAST(a AS CHAR) FROM t", &want)
			if _, err := mc.db.ExecContext(mc.ctx, "DROP TABLE t"); err != nil {
				t.Fatalf("oracle drop: %v", err)
			}
			replay := c.ShowCreateTable("testdb", "t")
			if _, err := mc.db.ExecContext(mc.ctx, replay); err != nil {
				t.Fatalf("oracle cannot read back what omni wrote: %v\n%s", err, replay)
			}
			if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
				t.Fatalf("oracle insert after replay: %v", err)
			}
			var got string
			oracleScan(t, mc, "SELECT CAST(a AS CHAR) FROM t", &got)
			if got != want {
				t.Errorf("replaying what omni wrote changed the default: %s became %s\n%s",
					want, got, replay)
			}
		})
	}
}

// TestReview_ApproximateBitAndYearDefault covers a BIT or YEAR default written
// with an exponent. Both cut the double's fraction where the integer types
// round it, which omni had them doing: BIT(8) DEFAULT (1.9e0) came out b'10'
// and YEAR DEFAULT (1.9e0) 2002, where the server stores 1 and 2001. The
// earlier 2.5e0 case could not tell the two apart, both readings giving 2.
//
// Measured on mariadb:11.8.8, with the integer and decimal families pinned
// beside them for the contrast: an exact fraction still rounds away from zero
// on BIT, where YEAR refuses one outright.
func TestReview_ApproximateBitAndYearDefault(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, typ, literal string }{
		// BIT and YEAR cut the fraction.
		{"bit_cuts_nine_tenths", "BIT(8)", "1.9e0"},
		{"bit_cuts_a_half", "BIT(8)", "1.5e0"},
		{"bit_cuts_two_and_a_half", "BIT(8)", "2.5e0"},
		{"bit_cuts_three_point_seven", "BIT(8)", "3.7e0"},
		{"bit_cuts_to_zero", "BIT(8)", "0.9e0"},
		{"bit_cuts_a_negative_to_zero", "BIT(8)", "-0.9e0"},
		{"year_cuts_nine_tenths", "YEAR", "1.9e0"},
		{"year_cuts_a_half", "YEAR", "1.5e0"},
		{"year_cuts_two_and_a_half", "YEAR", "2.5e0"},
		{"year_cuts_a_four_digit_year", "YEAR", "2001.9e0"},
		{"year_cuts_to_zero", "YEAR", "0.9e0"},
		// The value cut is the double's, not the digits written.
		{"bit_double_is_a_whole", "BIT(8)", "1.99999999999999999e0"},
		{"bit_double_is_one", "BIT(8)", "0.99999999999999999e0"},
		{"bit_double_just_past_two", "BIT(8)", "2.0000000000000001e0"},
		{"bit_double_short_of_two", "BIT(8)", "1.9999999999e0"},
		{"bit_wide_double", "BIT(64)", "9007199254740993e0"},
		{"bit_wide_exact_twin", "BIT(64)", "9007199254740993"},
		{"year_double_is_a_whole", "YEAR", "1.99999999999999999e0"},
		{"year_double_is_one", "YEAR", "0.99999999999999999e0"},
		{"year_double_rounds_a_year_up", "YEAR", "2000.9999999999999999e0"},
		// An exact fraction rounds away from zero on BIT.
		{"bit_exact_rounds_up", "BIT(8)", "1.9"},
		{"bit_exact_tie_rounds_up", "BIT(8)", "1.5"},
		{"bit_exact_seven_tenths", "BIT(8)", "1.7"},
		// The integer and decimal families round instead, and to even.
		{"int_rounds_nine_tenths", "INT", "1.9e0"},
		{"int_rounds_a_tie_to_even", "INT", "2.5e0"},
		{"int_rounds_the_other_tie_to_even", "INT", "3.5e0"},
		{"int_rounds_a_negative_tie_to_even", "INT", "-2.5e0"},
		{"int_exact_rounds_away_from_zero", "INT", "2.5"},
		{"decimal_rounds_away_from_zero", "DECIMAL(5,0)", "2.5e0"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			ddl := "CREATE TABLE t (a " + tc.typ + " DEFAULT (" + tc.literal + "))"
			runOnBoth(t, mc, c, ddl+";")
			assertDefaultMatches(t, mc, c)

			if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
				t.Fatalf("oracle insert: %v", err)
			}
			var want string
			oracleScan(t, mc, "SELECT CAST(a AS CHAR) FROM t", &want)
			if _, err := mc.db.ExecContext(mc.ctx, "DROP TABLE t"); err != nil {
				t.Fatalf("oracle drop: %v", err)
			}
			replay := c.ShowCreateTable("testdb", "t")
			if _, err := mc.db.ExecContext(mc.ctx, replay); err != nil {
				t.Fatalf("oracle cannot read back what omni wrote: %v\n%s", err, replay)
			}
			if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES ()"); err != nil {
				t.Fatalf("oracle insert after replay: %v", err)
			}
			var got string
			oracleScan(t, mc, "SELECT CAST(a AS CHAR) FROM t", &got)
			if got != want {
				t.Errorf("replaying what omni wrote changed the default: %s became %s\n%s",
					want, got, replay)
			}
		})
	}
}

// TestReview_UnaryBinaryIsACast covers the BINARY operator written before an
// expression. The server keeps no such operator: it stores the cast it stands
// for, so `2 * BINARY (1 + 2)` is 2 * cast(1 + 2 as char charset binary). omni
// wrote the operand bare — the operator and the grouping both gone — which
// turned `2 * BINARY (1 + 2)` into 2 * 1 + 2, 4 where the original is 6, and
// `BINARY 'a' = 'A'` into a comparison that holds where the original does not.
//
// Measured on mariadb:11.8.8, with the view's own rows read back.
func TestReview_UnaryBinaryIsACast(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)

	runOnBoth(t, mc, c, `CREATE TABLE t (a INT, b VARCHAR(10));`)
	if _, err := mc.db.ExecContext(mc.ctx, "INSERT INTO t VALUES (1,'x'),(2,'y')"); err != nil {
		t.Fatalf("oracle insert: %v", err)
	}
	for i, tc := range []struct{ name, body string }{
		{"under_a_product", "SELECT 2 * BINARY (1 + 2) AS x"},
		{"alone", "SELECT BINARY (1 + 2) AS x"},
		{"over_a_literal", "SELECT 2 * BINARY 3 AS x"},
		{"in_a_comparison", "SELECT BINARY 'a' = 'A' AS x"},
		{"under_a_minus", "SELECT -BINARY (1 + 2) AS x"},
		{"as_an_argument", "SELECT CONCAT(BINARY (1+2), 'x') AS x"},
		{"over_a_subquery", "SELECT 2 * BINARY (SELECT 3) AS x"},
		{"over_a_column", "SELECT BINARY b AS x FROM t"},
		// The cast the operator stands for, written out.
		{"the_cast_itself", "SELECT 2 * CAST((1 + 2) AS BINARY) AS x"},
		{"a_cast_over_a_column", "SELECT CAST(b AS BINARY) AS x FROM t"},
	} {
		view := fmt.Sprintf("v_bn%d", i)
		t.Run(tc.name, func(t *testing.T) {
			runOnBoth(t, mc, c, "CREATE VIEW "+view+" AS "+tc.body+";")
			assertViewBodyMatches(t, mc, c, view)

			want := fmt.Sprint(oracleRows(t, mc, "SELECT * FROM "+view))
			if _, err := mc.db.ExecContext(mc.ctx, "DROP VIEW "+view); err != nil {
				t.Fatalf("oracle drop: %v", err)
			}
			replay := c.ShowCreateView("testdb", view)
			if _, err := mc.db.ExecContext(mc.ctx, replay); err != nil {
				t.Fatalf("oracle cannot read back what omni wrote: %v\n%s", err, replay)
			}
			if got := fmt.Sprint(oracleRows(t, mc, "SELECT * FROM "+view)); got != want {
				t.Errorf("replaying what omni wrote changed the view's rows: %s became %s\n%s",
					want, got, replay)
			}
		})
	}
}

// TestReview_UnaryBinaryInStoredExpressions covers the same operator where a
// table stores the expression: a CHECK, a DEFAULT and a generated column each
// come back with the cast, the last through a writer of its own that no
// rewrite reaches.
//
// Measured on mariadb:11.8.8.
func TestReview_UnaryBinaryInStoredExpressions(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct{ name, ddl string }{
		{"a_check", "CREATE TABLE t (a INT, CHECK (BINARY 'a' = 'A'))"},
		{"a_default", "CREATE TABLE t (a INT DEFAULT (2 * BINARY (1 + 2)))"},
		{"a_generated_column",
			"CREATE TABLE t (a VARCHAR(5), b INT AS (BINARY (1 + 2)) VIRTUAL)"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			runOnBoth(t, mc, c, tc.ddl+";")

			want := oracleShow(t, mc, "SHOW CREATE TABLE t")
			if got := c.ShowCreateTable("testdb", "t"); got != want {
				t.Errorf("table mismatch\n--- oracle ---\n%s\n--- omni ---\n%s", want, got)
			}
			if _, err := mc.db.ExecContext(mc.ctx, "DROP TABLE t"); err != nil {
				t.Fatalf("oracle drop: %v", err)
			}
			replay := c.ShowCreateTable("testdb", "t")
			if _, err := mc.db.ExecContext(mc.ctx, replay); err != nil {
				t.Fatalf("oracle cannot read back what omni wrote: %v\n%s", err, replay)
			}
			if got := oracleShow(t, mc, "SHOW CREATE TABLE t"); got != want {
				t.Errorf("replaying what omni wrote changed the table\n--- was ---\n%s\n--- now ---\n%s",
					want, got)
			}
		})
	}
}

// TestReview_BacktickInAJSONColumnName covers an identifier holding a
// backtick, which has to be written twice. The implicit json_valid check a
// JSON column carries names that column and wrapped the raw name, so a column
// named a`b gave json_valid(`a`b`) — not that column, and not even parseable —
// and SHOW CREATE TABLE wrote the column, key, constraint and period names the
// same way.
//
// The FOREIGN KEY line is the one place the server does not double a backtick
// either, and omni follows it there rather than writing something the server
// does not: a constraint named fk`1 over a column c`1 comes back
// CONSTRAINT `fk`1` FOREIGN KEY (`c`1`) REFERENCES `p“q` (`i`d`), the
// referenced table doubled and nothing else, which is why such a table's own
// SHOW CREATE TABLE will not replay. These cases compare that text and stop
// there.
//
// Measured on mariadb:11.8.8 through SHOW CREATE TABLE, with what omni wrote
// replayed where the server can read it back.
func TestReview_BacktickInAJSONColumnName(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	for _, tc := range []struct {
		name, ddl string
		noReplay  bool // the text the server itself stores will not parse
	}{
		{"a_backtick", "CREATE TABLE t (`a``b` JSON)", false},
		{"a_backtick_with_a_check_of_its_own",
			"CREATE TABLE t (`a``b` JSON, CHECK (json_valid(`a``b`)))", false},
		{"a_backtick_on_a_longtext",
			"CREATE TABLE t (`a``b` LONGTEXT CHECK (json_valid(`a``b`)))", false},
		{"a_quote", "CREATE TABLE t (`q'u` JSON)", false},
		{"a_double_quote", "CREATE TABLE t (`a\"b` JSON)", false},
		{"renamed_to_a_backtick",
			"CREATE TABLE t (j JSON); ALTER TABLE t CHANGE COLUMN j `a``b` JSON", false},
		{"every_identifier_at_once",
			"CREATE TABLE `p``q` (`i``d` INT PRIMARY KEY); " +
				"CREATE TABLE t (`c``1` INT, `c``2` INT, KEY `k``1` (`c``2`), " +
				"CONSTRAINT `ck``1` CHECK (`c``1` > 0), " +
				"CONSTRAINT `fk``1` FOREIGN KEY (`c``1`) REFERENCES `p``q` (`i``d`))", true},
		{"a_generated_column_and_a_period",
			"CREATE TABLE t (`a``1` INT, `b``2` INT AS (`a``1` + 1) VIRTUAL, " +
				"`s``t` TIMESTAMP(6) AS ROW START, `e``n` TIMESTAMP(6) AS ROW END, " +
				"PERIOD FOR SYSTEM_TIME (`s``t`, `e``n`)) WITH SYSTEM VERSIONING", false},
		{"a_partition_on_a_backticked_column",
			"CREATE TABLE t (`a``1` INT) PARTITION BY HASH(`a``1`) PARTITIONS 2", false},
		{"a_range_partition_on_a_backticked_column",
			"CREATE TABLE t (`a``1` INT) PARTITION BY RANGE COLUMNS(`a``1`) " +
				"(PARTITION `p``0` VALUES LESS THAN (10))", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			scenarioReset(t, mc)
			c := scenarioNewCatalog(t)
			for _, stmt := range strings.Split(tc.ddl, "; ") {
				runOnBoth(t, mc, c, stmt+";")
			}

			want := oracleShow(t, mc, "SHOW CREATE TABLE t")
			if got := c.ShowCreateTable("testdb", "t"); got != want {
				t.Errorf("table mismatch\n--- oracle ---\n%s\n--- omni ---\n%s", want, got)
			}
			if tc.noReplay {
				// The server's own output holds the undoubled backticks of
				// its FOREIGN KEY line, which it answers 1064 on itself, so
				// matching it is all there is to check.
				return
			}

			if _, err := mc.db.ExecContext(mc.ctx, "DROP TABLE t"); err != nil {
				t.Fatalf("oracle drop: %v", err)
			}
			replay := c.ShowCreateTable("testdb", "t")
			if _, err := mc.db.ExecContext(mc.ctx, replay); err != nil {
				t.Fatalf("oracle cannot read back what omni wrote: %v\n%s", err, replay)
			}
			if got := oracleShow(t, mc, "SHOW CREATE TABLE t"); got != want {
				t.Errorf("replaying what omni wrote changed the table\n--- was ---\n%s\n--- now ---\n%s",
					want, got)
			}
		})
	}
}

// TestReview_ValuesRowsTakeTheRewrites covers a view backed by a table value
// constructor. Its rows were the last expression position the rewrite walk did
// not reach, so a SOUNDS LIKE in one came out at comparison precedence without
// the parentheses its left operand needs: `VALUES (((2=1) SOUNDS LIKE '0'))`
// stored values (2 = 1 sounds like '0'), which is 0 where the original is 1.
//
// Measured on mariadb:11.8.8. A VALUES node's own ORDER BY is positional, so
// `VALUES (1),(2) ORDER BY 1` keeps order by 1, and is pinned here beside the
// rows.
func TestReview_ValuesRowsTakeTheRewrites(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)
	c := scenarioNewCatalog(t)

	for i, tc := range []struct{ name, body string }{
		{"a_row_with_sounds_like", "VALUES (((2=1) SOUNDS LIKE '0'))"},
		{"a_row_with_a_not", "VALUES (NOT (1 > 0))"},
		{"two_rows", "VALUES (((2=1) SOUNDS LIKE '0')), (1)"},
		{"literals_under_an_and", "VALUES (1 AND 2)"},
		{"under_a_union", "SELECT 1 UNION VALUES (((2=1) SOUNDS LIKE '0'))"},
		{"ordered_rows", "VALUES (1),(2) ORDER BY 1"},
		{"ordered_rows_descending", "VALUES (2),(1) ORDER BY 1 DESC"},
		{"a_rewritten_row_ordered", "VALUES (((2=1) SOUNDS LIKE '0')) ORDER BY 1"},
		{"rows_with_a_limit", "VALUES (1),(2) LIMIT 1"},
	} {
		view := fmt.Sprintf("v_vs%d", i)
		t.Run(tc.name, func(t *testing.T) {
			runOnBoth(t, mc, c, "CREATE VIEW "+view+" AS "+tc.body+";")
			assertViewBodyMatches(t, mc, c, view)

			want := fmt.Sprint(oracleRows(t, mc, "SELECT * FROM "+view))
			if _, err := mc.db.ExecContext(mc.ctx, "DROP VIEW "+view); err != nil {
				t.Fatalf("oracle drop: %v", err)
			}
			replay := c.ShowCreateView("testdb", view)
			if _, err := mc.db.ExecContext(mc.ctx, replay); err != nil {
				t.Fatalf("oracle cannot read back what omni wrote: %v\n%s", err, replay)
			}
			if got := fmt.Sprint(oracleRows(t, mc, "SELECT * FROM "+view)); got != want {
				t.Errorf("replaying what omni wrote changed the view's rows: %s became %s\n%s",
					want, got, replay)
			}
		})
	}
}
