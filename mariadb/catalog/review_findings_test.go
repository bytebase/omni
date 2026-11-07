package catalog

import (
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
