package catalog

import (
	"strings"
	"testing"
)

// TestScenario_C19 covers Section C19 "Virtual / functional indexes" from
// mysql/catalog/SCENARIOS-mysql-implicit-behavior.md. MySQL 8.0.13+ implements
// functional index key parts by synthesizing a hidden VIRTUAL generated
// column over the expression and building an ordinary index over it. These
// scenarios lock the catalog behavior against a real MySQL oracle: hidden
// column synthesis, expression typing and validation, visibility suppression,
// JSON expression normalization, and hidden-column lifecycle.
func TestScenario_C19(t *testing.T) {
	scenariosSkipIfShort(t)
	scenariosSkipIfNoDocker(t)

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
			if colName != nil {
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
		// omni gap: no hidden column is created alongside the index.
		// MySQL's dd.columns.is_hidden=HT_HIDDEN_SQL row has no omni analog.
		// We flag this as a bug: the count of user-visible columns stays at
		// 2 in both engines, but omni's Column list should gain a
		// HiddenBySystem entry for round-trip fidelity. Currently it does
		// not — confirm with a direct probe.
		hiddenFound := false
		for _, col := range tbl.Columns {
			if strings.HasPrefix(col.Name, "!hidden!") {
				hiddenFound = true
				break
			}
		}
		if hiddenFound {
			t.Log("omni: unexpectedly found hidden column — partial support?")
		} else {
			t.Errorf("omni: no hidden column synthesized for functional index (MySQL: !hidden!idx_lower!0!0)")
		}

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

		// omni: the only way to verify the inferred type today is to look
		// for a synthesized hidden column. None exists, so this will fail.
		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Fatal("omni: table t not found")
		}
		for _, want := range []struct {
			idx      string
			wantType string // expected hidden-column DataType
		}{
			{"k_sum", "bigint"},
			{"k_low", "varchar"},
			{"k_cast", "bigint"},
		} {
			var hiddenCol *Column
			for _, col := range tbl.Columns {
				if strings.Contains(col.Name, want.idx) && strings.HasPrefix(col.Name, "!hidden!") {
					hiddenCol = col
					break
				}
			}
			if hiddenCol == nil {
				t.Errorf("omni %s: no hidden column synthesized; MySQL types this as %s", want.idx, want.wantType)
				continue
			}
			if !strings.EqualFold(hiddenCol.DataType, want.wantType) {
				t.Errorf("omni %s: hidden col type %q, want %q", want.idx, hiddenCol.DataType, want.wantType)
			}
		}
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

		// omni: the catalog should expose exactly 2 user-visible columns
		// AND the hidden column must not leak into the visible list. Since
		// omni synthesizes no hidden column at all, visible-count is
		// accidentally correct, but the storage-layer invariant (that a
		// hidden column exists and is queryable via an internal API) is
		// missing. Assert the hidden column is discoverable — this is the
		// bug to track.
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
		if hidden != 1 {
			t.Errorf("omni: hidden column count = %d, want 1 (HT_HIDDEN_SQL for idx_lower)", hidden)
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
				"Functional index on a column",
			},
			{
				// MySQL 8.0.45 actually reports ER 3753 "JSON or GEOMETRY"
				// for `->` (JSON return type). The broader substring
				// "functional index" covers both 3753 and 3754.
				"lob_json_rejected",
				"CREATE TABLE t4c (payload JSON, INDEX ((payload->'$.name')))",
				"functional index",
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
		} else if !strings.Contains(err.Error(), "functional index") {
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

		// Oracle: dropping the hidden column by name must fail with 3108.
		// First recreate the functional index so there's a hidden column
		// to try to drop.
		if _, err := mc.db.ExecContext(mc.ctx, "CREATE INDEX idx_lower ON t ((LOWER(name)))"); err != nil {
			t.Fatalf("oracle: recreating idx_lower: %v", err)
		}
		omniResults, omniErr := c.Exec("CREATE INDEX idx_lower ON t ((LOWER(name)))", nil)
		if omniErr != nil {
			t.Fatalf("omni: recreating idx_lower: %v", omniErr)
		}
		for _, r := range omniResults {
			if r.Error != nil {
				t.Fatalf("omni: recreating idx_lower: %v", r.Error)
			}
		}
		_, dropErr := mc.db.ExecContext(mc.ctx, "ALTER TABLE t DROP COLUMN `!hidden!idx_lower!0!0`")
		if dropErr == nil {
			t.Errorf("oracle: dropping hidden column should be rejected with ER 3108")
		} else if !strings.Contains(dropErr.Error(), "functional index") {
			t.Errorf("oracle: unexpected error for hidden-column drop: %v", dropErr)
		}

		// omni: the same ALTER must reject the system-hidden column with a
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

		runOnBoth(t, mc, c, "ALTER TABLE t RENAME INDEX idx_lower TO idx_lc")
		mysqlCreate = oracleShow(t, mc, "SHOW CREATE TABLE t")
		if strings.Contains(mysqlCreate, "idx_lower") || !strings.Contains(mysqlCreate, "idx_lc") {
			t.Errorf("oracle: RENAME INDEX not reflected in SHOW CREATE: %q", mysqlCreate)
		}
		omniCreate = c.ShowCreateTable("testdb", "t")
		if strings.Contains(omniCreate, "idx_lower") || !strings.Contains(omniCreate, "idx_lc") {
			t.Errorf("omni: RENAME INDEX not reflected in SHOW CREATE: %q", omniCreate)
		}
		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Fatalf("omni: table t missing after RENAME INDEX")
		}
		if col := tbl.GetColumn("!hidden!idx_lower!0!0"); col != nil {
			t.Errorf("omni: old hidden column still present after RENAME INDEX: %#v", col)
		}
		if col := tbl.GetColumn("!hidden!idx_lc!0!0"); col == nil || col.Hidden != ColumnHiddenSystem {
			t.Errorf("omni: renamed hidden column missing or not system-hidden: %#v", col)
		}
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
