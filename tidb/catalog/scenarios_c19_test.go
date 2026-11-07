package catalog

import (
	"errors"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// TestScenario_C19 checks TiDB's expression-index metadata, inferred key types,
// restrictions and lifecycle against the native engine. Internal storage column
// names are not part of its public catalog contract.
func TestScenario_C19(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	// -----------------------------------------------------------------
	// 19.1 Functional index creates a hidden VIRTUAL generated column
	// -----------------------------------------------------------------
	t.Run("19_1_hidden_virtual_column_created", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, "CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(64))")
		runOnBoth(t, mc, c, "CREATE INDEX idx_lower ON t ((LOWER(name)))")

		// Oracle: information_schema.STATISTICS exposes the functional
		// expression with COLUMN_NAME=NULL.
		rows := oracleRows(t, mc, `
			SELECT COLUMN_NAME, EXPRESSION
			  FROM information_schema.STATISTICS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND INDEX_NAME='idx_lower'`)
		if len(rows) != 1 {
			t.Errorf("oracle: expected 1 STATISTICS row for idx_lower, got %d", len(rows))
		} else {
			colName := rows[0][0]
			expr := asString(rows[0][1])
			if colName != nil && !strings.EqualFold(asString(colName), "NULL") {
				t.Errorf("oracle: functional index COLUMN_NAME should be NULL, got %v", colName)
			}
			if !strings.Contains(strings.ToLower(expr), "lower(`name`)") {
				t.Errorf("oracle: STATISTICS.EXPRESSION should contain lower(`name`), got %q", expr)
			}
		}

		// Oracle: SHOW CREATE renders functional key part as ((expr)).
		mysqlCreate := oracleShow(t, mc, "SHOW CREATE TABLE t")
		if !strings.Contains(strings.ToLower(mysqlCreate), "((lower(`name`)))") {
			t.Errorf("oracle: SHOW CREATE should contain ((lower(`name`))); got %q", mysqlCreate)
		}

		// omni: expose the hidden functional column. omni has no Hidden
		// flag today, so we assert on what's observable: the IndexColumn
		// should carry an expression and SHOW CREATE should match the
		// oracle's ((expr)) form byte-for-byte.
		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Fatal("omni: table t not found")
		}
		var idx *Index
		for _, i := range tbl.Indexes {
			if strings.EqualFold(i.Name, "idx_lower") {
				idx = i
				break
			}
		}
		if idx == nil {
			t.Fatal("omni: idx_lower not found")
		}
		if len(idx.Columns) != 1 {
			t.Errorf("omni: idx_lower expected 1 key part, got %d", len(idx.Columns))
		} else if idx.Columns[0].Expr == "" {
			t.Errorf("omni: idx_lower key part has no Expr; MySQL stores this as a hidden generated column")
		}
		if len(tbl.Columns) != 2 {
			t.Errorf("catalog exposed %d columns, want the two declared columns", len(tbl.Columns))
		}
		assertScenarioReadback(t, mc, c, "t")

		// Byte-exact SHOW CREATE comparison on the key part line.
		omniCreate := c.ShowCreateTable("testdb", "t")
		myKey := c19Line(mysqlCreate, "idx_lower")
		omniKey := c19Line(omniCreate, "idx_lower")
		if myKey != omniKey {
			t.Errorf("omni idx_lower key-part render differs from oracle:\noracle: %q\nomni:   %q", myKey, omniKey)
		}
	})

	// -----------------------------------------------------------------
	// 19.2 Hidden column type inferred from expression return type
	// -----------------------------------------------------------------
	t.Run("19_2_type_inferred_from_expression", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		ddl := `CREATE TABLE t (
			a INT, b INT,
			name VARCHAR(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
			payload JSON,
			INDEX k_sum ((a + b)),
			INDEX k_low ((LOWER(name))),
			INDEX k_cast ((CAST(payload->'$.age' AS UNSIGNED)))
		)`
		runOnBoth(t, mc, c, ddl)

		// Oracle: information_schema.STATISTICS should list all three
		// functional indexes, each with COLUMN_NAME=NULL and a non-empty
		// EXPRESSION. MySQL uses these to inform optimizer decisions.
		rows := oracleRows(t, mc, `
			SELECT INDEX_NAME, EXPRESSION
			  FROM information_schema.STATISTICS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
			 ORDER BY INDEX_NAME`)
		if len(rows) != 3 {
			t.Errorf("oracle: expected 3 functional index rows, got %d", len(rows))
		}

		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Fatal("omni: table t not found")
		}
		// Query result metadata independently exposes the types TiDB infers for
		// these same expressions; no MySQL-specific hidden columns are assumed.
		nativeRows, err := mc.db.QueryContext(mc.ctx, "SELECT a+b,LOWER(name),CAST(payload->'$.age' AS UNSIGNED) FROM t LIMIT 0")
		if err != nil {
			t.Fatal(err)
		}
		defer nativeRows.Close()
		types, err := nativeRows.ColumnTypes()
		if err != nil {
			t.Fatal(err)
		}
		if err := nativeRows.Close(); err != nil {
			t.Fatal(err)
		}
		for i, want := range []struct {
			name   string
			base   MySQLBaseType
			native string
		}{{"k_sum", BaseTypeBigInt, "BIGINT"}, {"k_low", BaseTypeVarchar, "VARCHAR"}, {"k_cast", BaseTypeBigInt, "UNSIGNED BIGINT"}} {
			if got := types[i].DatabaseTypeName(); got != want.native {
				t.Errorf("native %s type %q, want %q", want.name, got, want.native)
			}
			var part *IndexColumn
			for _, idx := range tbl.Indexes {
				if idx.Name == want.name && len(idx.Columns) == 1 {
					part = idx.Columns[0]
				}
			}
			if part == nil || part.ExprType == nil || part.ExprType.BaseType != want.base {
				t.Errorf("catalog %s inferred key type: %+v", want.name, part)
			} else if i == 2 && !part.ExprType.Unsigned {
				t.Error("unsigned CAST lost its unsigned result type")
			}
		}
		assertScenarioReadback(t, mc, c, "t")
	})

	// -----------------------------------------------------------------
	// 19.3 Hidden column suppressed in SELECT * and user I_S.COLUMNS
	// -----------------------------------------------------------------
	t.Run("19_3_hidden_suppressed_in_select_star", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, "CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(64))")
		runOnBoth(t, mc, c, "CREATE INDEX idx_lower ON t ((LOWER(name)))")

		// Oracle: user-scoped information_schema.COLUMNS does NOT list the
		// hidden column. Only (id, name) should appear.
		rows := oracleRows(t, mc, `
			SELECT COLUMN_NAME FROM information_schema.COLUMNS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
			 ORDER BY ORDINAL_POSITION`)
		if len(rows) != 2 {
			t.Errorf("oracle: expected 2 visible columns, got %d", len(rows))
		}
		for _, r := range rows {
			name := asString(r[0])
			if strings.HasPrefix(name, "!hidden!") {
				t.Errorf("oracle: user I_S.COLUMNS leaked hidden column %q", name)
			}
		}

		// Oracle: STATISTICS still records the expression (confirming the
		// hidden column exists at the storage layer).
		stats := oracleRows(t, mc, `
			SELECT EXPRESSION FROM information_schema.STATISTICS
			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND INDEX_NAME='idx_lower'`)
		if len(stats) != 1 || asString(stats[0][0]) == "" {
			t.Errorf("oracle: STATISTICS.EXPRESSION missing for idx_lower; got %v", stats)
		}

		// The model and SELECT * must expose the declared columns only.
		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Fatal("omni: table t not found")
		}
		var visible, hidden int
		for _, col := range tbl.Columns {
			if strings.HasPrefix(col.Name, "!hidden!") {
				hidden++
			} else {
				visible++
			}
		}
		if visible != 2 {
			t.Errorf("omni: visible column count = %d, want 2", visible)
		}
		if hidden != 0 {
			t.Errorf("catalog leaked %d internal columns", hidden)
		}
		starRows, err := mc.db.QueryContext(mc.ctx, "SELECT * FROM t LIMIT 0")
		if err != nil {
			t.Fatal(err)
		}
		defer starRows.Close()
		names, err := starRows.Columns()
		if err != nil || strings.Join(names, ",") != "id,name" {
			t.Errorf("native SELECT * columns: %v, error %v", names, err)
		}
	})

	// -----------------------------------------------------------------
	// 19.4 Functional expression must be deterministic / non-LOB
	// -----------------------------------------------------------------
	t.Run("19_4_disallowed_expression_rejected", func(t *testing.T) {
		scenarioReset(t, mc)

		// Oracle verification: run each bad DDL against MySQL directly
		// and confirm it's rejected. omni should reject the same.
		cases := []struct {
			label          string
			ddl            string
			mysqlErrSubstr string // expected substring of oracle error text
		}{
			{
				"rand_disallowed",
				"CREATE TABLE t4a (a INT, INDEX ((a + RAND())))",
				"disallowed function",
			},
			{
				"bare_column_rejected",
				"CREATE TABLE t4b (a INT, INDEX ((a)))",
				"Expression index on a column",
			},
			{
				// MySQL 8.0.45 actually reports ER 3753 "JSON or GEOMETRY"
				// for `->` (JSON return type). The broader substring
				// "functional index" covers both 3753 and 3754.
				"lob_json_rejected",
				"CREATE TABLE t4c (payload JSON, INDEX ((payload->'$.name')))",
				"expression index",
			},
		}
		for _, tc := range cases {
			t.Run(tc.label, func(t *testing.T) {
				// Oracle should reject.
				_, mysqlErr := mc.db.ExecContext(mc.ctx, tc.ddl)
				if mysqlErr == nil {
					t.Errorf("oracle: %s DDL unexpectedly succeeded", tc.label)
				} else if !strings.Contains(mysqlErr.Error(), tc.mysqlErrSubstr) {
					t.Errorf("oracle: %s error = %q, want substring %q", tc.label, mysqlErr.Error(), tc.mysqlErrSubstr)
				}

				// omni should reject as well. Use a fresh catalog so
				// earlier cases don't pollute.
				cc := scenarioNewCatalog(t)
				results, parseErr := cc.Exec(tc.ddl, nil)
				rejected := false
				if parseErr != nil {
					rejected = true
				} else {
					for _, r := range results {
						if r.Error != nil {
							rejected = true
							break
						}
					}
				}
				if !rejected {
					t.Errorf("omni: %s DDL was accepted; MySQL rejects it with %q",
						tc.label, tc.mysqlErrSubstr)
				}
			})
		}
	})

	// -----------------------------------------------------------------
	// 19.5 Functional index on JSON path via (col->>'$.path')
	// -----------------------------------------------------------------
	t.Run("19_5_json_path_functional_index", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		ddl := `CREATE TABLE t (
			id INT PRIMARY KEY,
			doc JSON,
			INDEX idx_name ((CAST(doc->>'$.name' AS CHAR(64))))
		)`
		runOnBoth(t, mc, c, ddl)

		mysqlCreate := oracleShow(t, mc, "SHOW CREATE TABLE t")
		omniCreate := c.ShowCreateTable("testdb", "t")

		// Oracle should render the full cast expression with ->> form.
		if !strings.Contains(strings.ToLower(mysqlCreate), "cast(") ||
			!strings.Contains(mysqlCreate, "idx_name") {
			t.Errorf("oracle: expected idx_name with CAST(...) clause; got %q", mysqlCreate)
		}

		// omni: key-part line must match the oracle byte-for-byte. This is
		// the round-trip test the scenario calls out as "the key test".
		myKey := c19Line(mysqlCreate, "idx_name")
		omniKey := c19Line(omniCreate, "idx_name")
		if myKey != omniKey {
			t.Errorf("omni idx_name render differs from oracle (round-trip test):\noracle: %q\nomni:   %q",
				myKey, omniKey)
		}

		// Plain `doc->>'$.name'` with no CAST must be rejected (LOB return).
		badDDL := "CREATE TABLE t_bad (doc JSON, INDEX ((doc->>'$.name')))"
		if _, err := mc.db.ExecContext(mc.ctx, badDDL); err == nil {
			t.Errorf("oracle: uncast ->> should be rejected as LOB")
		} else if !strings.Contains(err.Error(), "expression index") {
			t.Errorf("oracle: unexpected error for uncast ->>: %v", err)
		}
		cc := scenarioNewCatalog(t)
		results, parseErr := cc.Exec(badDDL, nil)
		rejected := parseErr != nil
		for _, r := range results {
			if r.Error != nil {
				rejected = true
			}
		}
		if !rejected {
			t.Errorf("omni: uncast ->> in functional index was accepted; MySQL rejects as LOB")
		}
	})

	// -----------------------------------------------------------------
	// 19.6 DROP INDEX cascades to hidden generated column
	// -----------------------------------------------------------------
	t.Run("19_6_drop_index_cascades_hidden", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, "CREATE TABLE t (id INT PRIMARY KEY, name VARCHAR(64))")
		runOnBoth(t, mc, c, "CREATE INDEX idx_lower ON t ((LOWER(name)))")
		runOnBoth(t, mc, c, "DROP INDEX idx_lower ON t")
		assertScenarioReadback(t, mc, c, "t")

		// Oracle: SHOW CREATE must show no trace of idx_lower or any
		// hidden column.
		mysqlCreate := oracleShow(t, mc, "SHOW CREATE TABLE t")
		if strings.Contains(mysqlCreate, "idx_lower") {
			t.Errorf("oracle: idx_lower still present after DROP INDEX: %q", mysqlCreate)
		}
		if strings.Contains(mysqlCreate, "!hidden!") {
			t.Errorf("oracle: hidden column leaked after DROP INDEX: %q", mysqlCreate)
		}

		// omni: same round-trip assertion.
		omniCreate := c.ShowCreateTable("testdb", "t")
		if strings.Contains(omniCreate, "idx_lower") {
			t.Errorf("omni: idx_lower still present after DROP INDEX: %q", omniCreate)
		}
		if strings.Contains(omniCreate, "!hidden!") {
			t.Errorf("omni: hidden column leaked after DROP INDEX: %q", omniCreate)
		}

		// The inherited MySQL internal name is not a TiDB column. Keep the
		// original DROP probe and compare its native unknown-column error.
		runOnBoth(t, mc, c, "CREATE INDEX idx_lower ON t ((LOWER(name)))")
		_, dropErr := mc.db.ExecContext(mc.ctx, "ALTER TABLE t DROP COLUMN `!hidden!idx_lower!0!0`")
		var nativeErr *gomysql.MySQLError
		if !errors.As(dropErr, &nativeErr) || nativeErr.Number != 1091 {
			t.Fatalf("native internal-name DROP error: %v", dropErr)
		}

		// omni: the same ALTER. omni has no hidden column, so it should
		// either reject the column name as "unknown" OR reject it with a
		// 3108-equivalent error. Accepting the DROP COLUMN silently is
		// the bug we want to catch.
		results, parseErr := c.Exec("ALTER TABLE t DROP COLUMN `!hidden!idx_lower!0!0`", nil)
		rejected := parseErr != nil
		for _, r := range results {
			if r.Error != nil {
				rejected = true
			}
		}
		if !rejected {
			t.Errorf("omni: DROP COLUMN `!hidden!idx_lower!0!0` silently accepted; MySQL returns ER_CANNOT_DROP_COLUMN_FUNCTIONAL_INDEX (3108)")
		}
		if parseErr != nil || len(results) != 1 {
			t.Fatalf("catalog DROP results: %v, %v", results, parseErr)
		}
		var catalogErr *Error
		if !errors.As(results[0].Error, &catalogErr) || catalogErr.Code != int(nativeErr.Number) {
			t.Errorf("catalog internal-name DROP error: %v, native: %v", results[0].Error, nativeErr)
		}
		assertScenarioReadback(t, mc, c, "t")
	})
}

// c19Line returns the line from s that contains needle, trimmed of the
// trailing comma some SHOW CREATE outputs use. Empty string if no match.
func c19Line(s, needle string) string {
	for _, line := range strings.Split(s, "\n") {
		if strings.Contains(line, needle) {
			return strings.TrimRight(strings.TrimSpace(line), ",")
		}
	}
	return ""
}
