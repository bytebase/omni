package catalog

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/parser"
)

// TestScenario_C21 covers section C21 (Parser-level implicit defaults) from
// mysql/catalog/SCENARIOS-mysql-implicit-behavior.md. Each subtest asserts that
// real TiDB v8.5.5 and the omni AST / catalog agree on the default value that
// the grammar fills in when the user omits a clause.
//
// Some of these scenarios are inherently parser-AST-level (JOIN type, ORDER
// direction, LIMIT offset, INSERT column list) so they assert against the
// parsed AST directly rather than via the catalog. DDL-level defaults
// (DEFAULT NULL, FK actions, CREATE INDEX USING, CREATE VIEW ALGORITHM, ENGINE)
// assert against both the container oracle and omni's catalog.
//
// Forked from MySQL; the FK-action, ORDER-direction and ENGINE subtests are
// re-baselined to TiDB v8.5.5's measured behavior (see the per-subtest notes).
func TestScenario_C21(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()

	// --- 21.1 DEFAULT without value on nullable column -> DEFAULT NULL ----
	t.Run("21_1_default_null_on_nullable", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		ddl := `CREATE TABLE t (a INT DEFAULT NULL);
CREATE TABLE t2 (c INT);`
		runOnBoth(t, mc, c, ddl)

		// Oracle: column a has COLUMN_DEFAULT NULL, IS_NULLABLE YES.
		var aDefault *string
		var aNullable string
		oracleScan(t, mc,
			`SELECT COLUMN_DEFAULT, IS_NULLABLE FROM information_schema.COLUMNS
               WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'`,
			&aDefault, &aNullable)
		if aDefault != nil {
			t.Errorf("oracle: expected COLUMN_DEFAULT NULL for t.a, got %q", *aDefault)
		}
		assertStringEq(t, "oracle t.a IS_NULLABLE", aNullable, "YES")

		// Oracle: column c (no DEFAULT clause at all) also has DEFAULT NULL.
		var cDefault *string
		var cNullable string
		oracleScan(t, mc,
			`SELECT COLUMN_DEFAULT, IS_NULLABLE FROM information_schema.COLUMNS
               WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t2' AND COLUMN_NAME='c'`,
			&cDefault, &cNullable)
		if cDefault != nil {
			t.Errorf("oracle: expected COLUMN_DEFAULT NULL for t2.c, got %q", *cDefault)
		}
		assertStringEq(t, "oracle t2.c IS_NULLABLE", cNullable, "YES")

		// omni: t.a explicit DEFAULT NULL — omni should model this as Nullable
		// true and either Default == nil or Default pointing at "NULL".
		tbl := c.GetDatabase("testdb").GetTable("t")
		if tbl == nil {
			t.Errorf("omni: table t missing")
		} else {
			col := tbl.GetColumn("a")
			if col == nil {
				t.Errorf("omni: column t.a missing")
			} else {
				assertBoolEq(t, "omni t.a nullable", col.Nullable, true)
				if col.Default != nil && strings.ToUpper(*col.Default) != "NULL" {
					t.Errorf("omni: t.a Default = %q, want nil or NULL", *col.Default)
				}
			}
		}

		// omni: t2.c — no DEFAULT clause at all. Expect Nullable true and
		// Default == nil (omni does not synthesize a "NULL" literal).
		tbl2 := c.GetDatabase("testdb").GetTable("t2")
		if tbl2 == nil {
			t.Errorf("omni: table t2 missing")
		} else {
			col := tbl2.GetColumn("c")
			if col == nil {
				t.Errorf("omni: column t2.c missing")
			} else {
				assertBoolEq(t, "omni t2.c nullable", col.Nullable, true)
				if col.Default != nil && strings.ToUpper(*col.Default) != "NULL" {
					t.Errorf("omni: t2.c Default = %q, want nil (no clause)", *col.Default)
				}
			}
		}
	})

	// --- 21.2 Bare JOIN -> INNER JOIN (JoinInner) --------------------------
	t.Run("21_2_join_type_default_inner", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		// Create tables so the container also accepts the SELECT.
		runOnBoth(t, mc, c, `CREATE TABLE t1 (a INT); CREATE TABLE t2 (a INT);`)

		cases := []struct {
			name string
			sql  string
		}{
			{"bare_JOIN", "SELECT * FROM t1 JOIN t2 ON t1.a = t2.a"},
			{"INNER_JOIN", "SELECT * FROM t1 INNER JOIN t2 ON t1.a = t2.a"},
			{"CROSS_JOIN", "SELECT * FROM t1 CROSS JOIN t2"},
		}
		for _, tc := range cases {
			// Oracle: just ensure MySQL accepts it.
			if _, err := mc.db.ExecContext(mc.ctx, "USE testdb; "+tc.sql); err != nil {
				t.Errorf("oracle %s: %v", tc.name, err)
			}

			// omni AST: parse and inspect JoinClause.Type.
			jt, ok := c21FirstJoinType(t, tc.sql)
			if !ok {
				continue
			}
			// Grammar-level normalization: bare JOIN and INNER JOIN should both
			// map to JoinInner. CROSS JOIN also maps to JoinInner per yacc but
			// omni tracks JoinCross distinctly for deparse fidelity, which is
			// acceptable — assert that the bare forms collapse.
			if tc.name == "CROSS_JOIN" {
				if jt != nodes.JoinInner && jt != nodes.JoinCross {
					t.Errorf("omni %s: JoinType = %d, want JoinInner or JoinCross", tc.name, jt)
				}
			} else {
				if jt != nodes.JoinInner {
					t.Errorf("omni %s: JoinType = %d, want JoinInner(%d)", tc.name, jt, nodes.JoinInner)
				}
			}
		}
	})

	// --- 21.3 ORDER BY without direction -> tri-state ---------------------
	t.Run("21_3_order_by_no_direction", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		runOnBoth(t, mc, c, `CREATE TABLE t (a INT);`)

		// Oracle: MySQL accepts both; we cannot read ORDER_NOT_RELEVANT from
		// info_schema directly, so just verify both parse.
		for _, s := range []string{
			"SELECT * FROM t ORDER BY a",
			"SELECT * FROM t ORDER BY a ASC",
		} {
			if _, err := mc.db.ExecContext(mc.ctx, "USE testdb; "+s); err != nil {
				t.Errorf("oracle %q: %v", s, err)
			}
		}

		// omni AST: OrderByItem has a single Desc bool, so bare ORDER BY and
		// explicit ASC both parse to Desc=false and explicit DESC to Desc=true.
		// The missing tri-state (ORDER_NOT_RELEVANT vs ASC) is a representational
		// asymmetry, not a parity divergence: TiDB v8.5.5 exposes no observable
		// difference between `ORDER BY a` and `ORDER BY a ASC` (both accepted,
		// both ascending, no information_schema/SHOW surface distinguishes them —
		// confirmed against the container), and omni deparses both as `ORDER BY
		// a`, matching TiDB's restore. A tri-state field would touch the AST
		// node, parser, generated walker and every consumer with no
		// oracle-measurable benefit, so it is documented and deferred.
		bare := c21FirstOrderByItem(t, "SELECT * FROM t ORDER BY a")
		asc := c21FirstOrderByItem(t, "SELECT * FROM t ORDER BY a ASC")
		desc := c21FirstOrderByItem(t, "SELECT * FROM t ORDER BY a DESC")
		if bare == nil || asc == nil || desc == nil {
			return
		}
		assertBoolEq(t, "omni bare ORDER BY Desc", bare.Desc, false)
		assertBoolEq(t, "omni ORDER BY ASC Desc", asc.Desc, false)
		assertBoolEq(t, "omni ORDER BY DESC Desc", desc.Desc, true)
	})

	// --- 21.4 LIMIT N without OFFSET -> opt_offset NULL -------------------
	t.Run("21_4_limit_without_offset", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		runOnBoth(t, mc, c, `CREATE TABLE t (a INT);`)

		// Oracle: verify both accept.
		for _, s := range []string{
			"SELECT * FROM t LIMIT 10",
			"SELECT * FROM t LIMIT 10 OFFSET 0",
		} {
			if _, err := mc.db.ExecContext(mc.ctx, "USE testdb; "+s); err != nil {
				t.Errorf("oracle %q: %v", s, err)
			}
		}

		// omni AST: LIMIT 10 -> Offset should be nil; LIMIT 10 OFFSET 0 ->
		// Offset should be non-nil.
		limBare := c21FirstLimit(t, "SELECT * FROM t LIMIT 10")
		limOff := c21FirstLimit(t, "SELECT * FROM t LIMIT 10 OFFSET 0")
		if limBare == nil {
			t.Errorf("omni: LIMIT 10 produced no Limit node")
		} else if limBare.Offset != nil {
			t.Errorf("omni: LIMIT 10 Offset = %v, want nil", limBare.Offset)
		}
		if limOff == nil {
			t.Errorf("omni: LIMIT 10 OFFSET 0 produced no Limit node")
		} else if limOff.Offset == nil {
			t.Errorf("omni: LIMIT 10 OFFSET 0 Offset = nil, want non-nil (Item_uint(0))")
		}
	})

	// --- 21.5 FK ON DELETE omitted -> FK_OPTION_UNDEF ---------------------
	t.Run("21_5_fk_on_delete_omitted_undef", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		ddl := `CREATE TABLE parent (id INT PRIMARY KEY);
CREATE TABLE child (p INT, FOREIGN KEY (p) REFERENCES parent(id));`
		runOnBoth(t, mc, c, ddl)

		// Oracle: information_schema.REFERENTIAL_CONSTRAINTS reports
		// DELETE_RULE and UPDATE_RULE as "NO ACTION" (rendering of UNDEF).
		var deleteRule, updateRule string
		oracleScan(t, mc,
			`SELECT DELETE_RULE, UPDATE_RULE FROM information_schema.REFERENTIAL_CONSTRAINTS
              WHERE CONSTRAINT_SCHEMA='testdb' AND TABLE_NAME='child'`,
			&deleteRule, &updateRule)
		assertStringEq(t, "oracle DELETE_RULE", deleteRule, "NO ACTION")
		assertStringEq(t, "oracle UPDATE_RULE", updateRule, "NO ACTION")

		// omni catalog: the FK constraint's OnDelete/OnUpdate should render as
		// "NO ACTION" (via refActionToString mapping RefActNone -> "NO ACTION").
		tbl := c.GetDatabase("testdb").GetTable("child")
		if tbl == nil {
			t.Errorf("omni: table child missing")
			return
		}
		var fk *Constraint
		for _, con := range tbl.Constraints {
			if con.Type == ConForeignKey {
				fk = con
				break
			}
		}
		if fk == nil {
			t.Errorf("omni: FK constraint missing on child")
			return
		}
		// omni's internal model holds "" for an unspecified action, which is the
		// TiDB-correct representation. Measured against the container: TiDB's
		// SHOW CREATE emits ON DELETE/ON UPDATE only for an explicit
		// NO ACTION/RESTRICT/CASCADE, never for an unspecified action, so the ""
		// sentinel is what keeps omni's SHOW CREATE from emitting the clause (it
		// matches TiDB exactly). TiDB's REFERENTIAL_CONSTRAINTS.DELETE_RULE does
		// report "NO ACTION" (asserted on the oracle above), but omni has no
		// information_schema query engine to surface, so there is nothing to map.
		// Accept the default set (empty / NO ACTION / RESTRICT), as C5 does.
		if !c5IsFKDefault(fk.OnDelete) {
			t.Errorf("omni FK OnDelete: got %q, want default (empty/NO ACTION/RESTRICT)", fk.OnDelete)
		}
		if !c5IsFKDefault(fk.OnUpdate) {
			t.Errorf("omni FK OnUpdate: got %q, want default (empty/NO ACTION/RESTRICT)", fk.OnUpdate)
		}
	})

	// --- 21.6 FK ON DELETE present, ON UPDATE omitted ---------------------
	t.Run("21_6_fk_on_update_independent_undef", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		ddl := `CREATE TABLE parent (id INT PRIMARY KEY);
CREATE TABLE child (p INT, FOREIGN KEY (p) REFERENCES parent(id) ON DELETE CASCADE);`
		runOnBoth(t, mc, c, ddl)

		// Oracle: DELETE_RULE=CASCADE, UPDATE_RULE=NO ACTION (not inherited).
		var deleteRule, updateRule string
		oracleScan(t, mc,
			`SELECT DELETE_RULE, UPDATE_RULE FROM information_schema.REFERENTIAL_CONSTRAINTS
              WHERE CONSTRAINT_SCHEMA='testdb' AND TABLE_NAME='child'`,
			&deleteRule, &updateRule)
		assertStringEq(t, "oracle DELETE_RULE", deleteRule, "CASCADE")
		assertStringEq(t, "oracle UPDATE_RULE", updateRule, "NO ACTION")

		tbl := c.GetDatabase("testdb").GetTable("child")
		if tbl == nil {
			t.Errorf("omni: table child missing")
			return
		}
		var fk *Constraint
		for _, con := range tbl.Constraints {
			if con.Type == ConForeignKey {
				fk = con
				break
			}
		}
		if fk == nil {
			t.Errorf("omni: FK constraint missing")
			return
		}
		// The explicit ON DELETE CASCADE is retained; the omitted ON UPDATE stays
		// at omni's default representation (see 21.5) rather than being inherited
		// from ON DELETE — TiDB reports UPDATE_RULE "NO ACTION" independently
		// (asserted on the oracle above).
		assertStringEq(t, "omni FK OnDelete", fk.OnDelete, "CASCADE")
		if !c5IsFKDefault(fk.OnUpdate) {
			t.Errorf("omni FK OnUpdate: got %q, want default (empty/NO ACTION/RESTRICT)", fk.OnUpdate)
		}
	})

	// --- 21.7 CREATE INDEX without USING -> nullptr (engine picks) --------
	t.Run("21_7_index_no_using_clause", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		ddl := `CREATE TABLE t (a INT, KEY k (a));
CREATE TABLE t_b (a INT, KEY kb (a) USING BTREE);`
		runOnBoth(t, mc, c, ddl)

		// Oracle: information_schema.STATISTICS.INDEX_TYPE is "BTREE" for both
		// on InnoDB — engine fills in the default. We check only that both
		// resolve to "BTREE" (InnoDB default), confirming the engine-fill.
		var t1Type, t2Type string
		oracleScan(t, mc,
			`SELECT INDEX_TYPE FROM information_schema.STATISTICS
              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND INDEX_NAME='k'`,
			&t1Type)
		oracleScan(t, mc,
			`SELECT INDEX_TYPE FROM information_schema.STATISTICS
              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t_b' AND INDEX_NAME='kb'`,
			&t2Type)
		assertStringEq(t, "oracle index k type", t1Type, "BTREE")
		assertStringEq(t, "oracle index kb type", t2Type, "BTREE")

		// omni: the parser should preserve the distinction (no USING -> empty,
		// USING BTREE -> "BTREE"). The engine default resolution is the
		// catalog's job, not the parser's — but omni may collapse both.
		tbl := c.GetDatabase("testdb").GetTable("t")
		tbl2 := c.GetDatabase("testdb").GetTable("t_b")
		if tbl == nil || tbl2 == nil {
			t.Errorf("omni: missing tables t/t_b")
			return
		}
		var kType, kbType string
		for _, idx := range tbl.Indexes {
			if idx.Name == "k" {
				kType = idx.IndexType
			}
		}
		for _, idx := range tbl2.Indexes {
			if idx.Name == "kb" {
				kbType = idx.IndexType
			}
		}
		// omni grammar: no USING clause should leave IndexType empty. If omni
		// defaults to "BTREE" at parse time, that's a bug (loses the info
		// needed to deparse faithfully).
		if kType != "" {
			t.Errorf("omni: index k IndexType = %q, want \"\" (no USING clause)", kType)
		}
		assertStringEq(t, "omni index kb IndexType", kbType, "BTREE")
	})

	// --- 21.8 INSERT without column list -> empty Columns -----------------
	t.Run("21_8_insert_no_column_list", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		runOnBoth(t, mc, c, `CREATE TABLE t (a INT, b INT, c INT);`)

		// Oracle: both forms accepted.
		for _, s := range []string{
			"INSERT INTO t VALUES (1, 2, 3)",
			"INSERT INTO t (a, b, c) VALUES (4, 5, 6)",
		} {
			if _, err := mc.db.ExecContext(mc.ctx, "USE testdb; "+s); err != nil {
				t.Errorf("oracle %q: %v", s, err)
			}
		}

		// omni AST: INSERT without column list -> Columns is nil / empty.
		bare := c21FirstInsert(t, "INSERT INTO t VALUES (1, 2, 3)")
		full := c21FirstInsert(t, "INSERT INTO t (a, b, c) VALUES (4, 5, 6)")
		if bare == nil || full == nil {
			return
		}
		if len(bare.Columns) != 0 {
			t.Errorf("omni: INSERT without column list produced %d columns, want 0",
				len(bare.Columns))
		}
		if len(full.Columns) != 3 {
			t.Errorf("omni: INSERT with explicit column list produced %d columns, want 3",
				len(full.Columns))
		}
	})

	// --- 21.9 CREATE VIEW without ALGORITHM -> UNDEFINED ------------------
	t.Run("21_9_view_no_algorithm_undefined", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)
		runOnBoth(t, mc, c, `CREATE TABLE t (a INT);`)

		// Three forms that all resolve to ALGORITHM=UNDEFINED.
		runOnBoth(t, mc, c, `CREATE VIEW v1 AS SELECT 1 AS x;`)
		runOnBoth(t, mc, c, `CREATE ALGORITHM=UNDEFINED VIEW v2 AS SELECT 1 AS x;`)

		// Oracle: information_schema.VIEWS has no ALGORITHM column in MySQL
		// 8.0; use SHOW CREATE VIEW which renders the algorithm verbatim.
		// Both forms should contain ALGORITHM=UNDEFINED.
		for _, name := range []string{"v1", "v2"} {
			create := oracleShow(t, mc, "SHOW CREATE VIEW "+name)
			if !strings.Contains(strings.ToUpper(create), "ALGORITHM=UNDEFINED") {
				t.Errorf("oracle %s: SHOW CREATE VIEW missing ALGORITHM=UNDEFINED: %s",
					name, create)
			}
		}

		// omni: View.Algorithm should either be "" (default) or "UNDEFINED"
		// for v1 (no explicit clause), and "UNDEFINED" for v2.
		db := c.GetDatabase("testdb")
		v1 := db.Views["v1"]
		v2 := db.Views["v2"]
		if v1 == nil {
			t.Errorf("omni: view v1 missing")
		} else if v1.Algorithm != "" && !strings.EqualFold(v1.Algorithm, "UNDEFINED") {
			t.Errorf("omni: v1 Algorithm = %q, want \"\" or UNDEFINED", v1.Algorithm)
		}
		if v2 == nil {
			t.Errorf("omni: view v2 missing")
		} else if !strings.EqualFold(v2.Algorithm, "UNDEFINED") {
			t.Errorf("omni: v2 Algorithm = %q, want UNDEFINED", v2.Algorithm)
		}
	})

	// --- 21.10 CREATE TABLE without ENGINE -> post-parse fill ------------
	t.Run("21_10_create_table_no_engine", func(t *testing.T) {
		scenarioReset(t, mc)
		c := scenarioNewCatalog(t)

		runOnBoth(t, mc, c, `CREATE TABLE t_noeng (a INT);`)
		runOnBoth(t, mc, c, `CREATE TABLE t_eng (a INT) ENGINE=InnoDB;`)

		// Oracle: both resolve to "InnoDB" via session default.
		var e1, e2 string
		oracleScan(t, mc,
			`SELECT ENGINE FROM information_schema.TABLES
              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t_noeng'`,
			&e1)
		oracleScan(t, mc,
			`SELECT ENGINE FROM information_schema.TABLES
              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t_eng'`,
			&e2)
		assertStringEq(t, "oracle t_noeng engine", e1, "InnoDB")
		assertStringEq(t, "oracle t_eng engine", e2, "InnoDB")

		// omni's catalog resolves the session-default engine (InnoDB) for a
		// table with no ENGINE clause. That matches TiDB v8.5.5, where
		// information_schema.TABLES.ENGINE and SHOW CREATE both report InnoDB for
		// t_noeng and t_eng alike (measured above). The MySQL-parser notion of a
		// NULL engine at parse time is not observable in TiDB, so omni holding
		// "InnoDB" is parity-correct, not a premature fill.
		tNo := c.GetDatabase("testdb").GetTable("t_noeng")
		tEx := c.GetDatabase("testdb").GetTable("t_eng")
		if tNo == nil || tEx == nil {
			t.Errorf("omni: missing t_noeng or t_eng")
			return
		}
		if !strings.EqualFold(tNo.Engine, "InnoDB") {
			t.Errorf("omni: t_noeng Engine = %q, want InnoDB (session default, TiDB-matching)", tNo.Engine)
		}
		if !strings.EqualFold(tEx.Engine, "InnoDB") {
			t.Errorf("omni: t_eng Engine = %q, want InnoDB (explicit)", tEx.Engine)
		}
	})
}

// --- section-local helpers -------------------------------------------------

// c21FirstJoinType parses a SELECT query and returns the JoinType of the
// first JoinClause found in the FROM list. Uses t.Errorf (not Fatal) so the
// subtest keeps running.
func c21FirstJoinType(t *testing.T, sql string) (nodes.JoinType, bool) {
	t.Helper()
	stmts, err := parser.Parse(sql)
	if err != nil {
		t.Errorf("omni parse %q: %v", sql, err)
		return 0, false
	}
	if len(stmts.Items) == 0 {
		t.Errorf("omni parse %q: no statements", sql)
		return 0, false
	}
	sel, ok := stmts.Items[0].(*nodes.SelectStmt)
	if !ok {
		t.Errorf("omni parse %q: expected SelectStmt, got %T", sql, stmts.Items[0])
		return 0, false
	}
	for _, te := range sel.From {
		if jc, ok := te.(*nodes.JoinClause); ok {
			return jc.Type, true
		}
	}
	t.Errorf("omni parse %q: no JoinClause in FROM list", sql)
	return 0, false
}

// c21FirstOrderByItem parses a SELECT and returns the first ORDER BY item.
func c21FirstOrderByItem(t *testing.T, sql string) *nodes.OrderByItem {
	t.Helper()
	stmts, err := parser.Parse(sql)
	if err != nil {
		t.Errorf("omni parse %q: %v", sql, err)
		return nil
	}
	if len(stmts.Items) == 0 {
		t.Errorf("omni parse %q: no statements", sql)
		return nil
	}
	sel, ok := stmts.Items[0].(*nodes.SelectStmt)
	if !ok {
		t.Errorf("omni parse %q: expected SelectStmt, got %T", sql, stmts.Items[0])
		return nil
	}
	if len(sel.OrderBy) == 0 {
		t.Errorf("omni parse %q: no ORDER BY items", sql)
		return nil
	}
	return sel.OrderBy[0]
}

// c21FirstLimit parses a SELECT and returns its Limit node (or nil).
func c21FirstLimit(t *testing.T, sql string) *nodes.Limit {
	t.Helper()
	stmts, err := parser.Parse(sql)
	if err != nil {
		t.Errorf("omni parse %q: %v", sql, err)
		return nil
	}
	if len(stmts.Items) == 0 {
		return nil
	}
	sel, ok := stmts.Items[0].(*nodes.SelectStmt)
	if !ok {
		t.Errorf("omni parse %q: expected SelectStmt, got %T", sql, stmts.Items[0])
		return nil
	}
	return sel.Limit
}

// c21FirstInsert parses an INSERT and returns the InsertStmt node.
func c21FirstInsert(t *testing.T, sql string) *nodes.InsertStmt {
	t.Helper()
	stmts, err := parser.Parse(sql)
	if err != nil {
		t.Errorf("omni parse %q: %v", sql, err)
		return nil
	}
	if len(stmts.Items) == 0 {
		t.Errorf("omni parse %q: no statements", sql)
		return nil
	}
	ins, ok := stmts.Items[0].(*nodes.InsertStmt)
	if !ok {
		t.Errorf("omni parse %q: expected InsertStmt, got %T", sql, stmts.Items[0])
		return nil
	}
	return ins
}
