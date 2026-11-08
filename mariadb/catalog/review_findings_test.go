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
