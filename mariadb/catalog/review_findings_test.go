package catalog

import (
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
