package parser

import (
	"reflect"
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// rangeVarTexts returns the source text of every RangeVar in sql, in walk
// order. A RangeVar without a location yields "<NoLoc>".
func rangeVarTexts(t *testing.T, sql string) []string {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	var texts []string
	var visit func(nodes.Node) bool
	visit = func(n nodes.Node) bool {
		switch n := n.(type) {
		case *nodes.RangeVar:
			if n.Loc.Start < 0 || n.Loc.End < n.Loc.Start || n.Loc.End > len(sql) {
				texts = append(texts, "<NoLoc>")
			} else {
				texts = append(texts, sql[n.Loc.Start:n.Loc.End])
			}
		case *nodes.CreateForeignTableStmt:
			// The walker does not descend into the embedded CreateStmt.
			nodes.Inspect(&n.Base, visit)
		}
		return true
	}
	nodes.Inspect(list, visit)
	return texts
}

// TestRangeVarLoc pins that every RangeVar the parser builds carries a
// location, and that the location ends at the name's last token rather than
// at the start of the next one (root AGENTS.md, "Engine contract").
func TestRangeVarLoc(t *testing.T) {
	tests := []struct {
		sql  string
		want []string
	}{
		// Creation statements.
		{`CREATE TABLE x (a int)`, []string{"x"}},
		{`CREATE TABLE s.t (a int)`, []string{"s.t"}},
		{`CREATE TABLE "Q"."T" (a int)`, []string{`"Q"."T"`}},
		{`CREATE TABLE db.s.t (a int)`, []string{"db.s.t"}},
		{`CREATE TABLE U&"\00C4b" (a int)`, []string{`U&"\00C4b"`}},
		{`CREATE TABLE IF NOT EXISTS s . t (a int)`, []string{"s . t"}},
		{`CREATE TEMP TABLE x (a int)`, []string{"x"}},
		{`CREATE UNLOGGED TABLE x (a int)`, []string{"x"}},
		{`CREATE TABLE c PARTITION OF s.p FOR VALUES IN (1)`, []string{"c", "s.p"}},
		{`CREATE TABLE c (a int) INHERITS (p1, s.p2)`, []string{"c", "p1", "s.p2"}},
		{`CREATE TABLE c (a int REFERENCES r(a), FOREIGN KEY (a) REFERENCES s.r2 (b))`, []string{"c", "r", "s.r2"}},
		{`CREATE TABLE x (LIKE s.y)`, []string{"x", "s.y"}},
		{`CREATE TABLE x OF typ`, []string{"x"}},
		{`CREATE TABLE x AS SELECT 1`, []string{"x"}},
		{`CREATE TABLE s.x (a) AS SELECT 1`, []string{"s.x"}},
		{`CREATE TABLE x AS EXECUTE p`, []string{"x"}},
		{`CREATE SCHEMA s CREATE TABLE t (a int)`, []string{"t"}},
		{`CREATE VIEW v AS SELECT 1`, []string{"v"}},
		{`CREATE OR REPLACE VIEW s.v (a) AS SELECT 1`, []string{"s.v"}},
		// The FROM item of the expanded recursive view is synthesized.
		{`CREATE RECURSIVE VIEW v (a) AS SELECT 1`, []string{"v", "<NoLoc>"}},
		{`CREATE MATERIALIZED VIEW s.mv AS SELECT 1`, []string{"s.mv"}},
		{`CREATE SEQUENCE IF NOT EXISTS s.seq START 1`, []string{"s.seq"}},
		{`CREATE INDEX ON s.t (a)`, []string{"s.t"}},
		{`CREATE INDEX i ON ONLY t (a)`, []string{"ONLY t"}},
		{`CREATE FOREIGN TABLE s.ft (a int) SERVER srv`, []string{"s.ft"}},
		{`CREATE FOREIGN TABLE ft PARTITION OF s.p FOR VALUES IN (1) SERVER srv`, []string{"ft", "s.p"}},
		{`CREATE TYPE ct AS (a int)`, []string{"ct"}},
		{`CREATE TRIGGER tr BEFORE INSERT ON s.t FOR EACH ROW EXECUTE FUNCTION f()`, []string{"s.t"}},
		{`CREATE CONSTRAINT TRIGGER tr AFTER INSERT ON s.t FROM s.u FOR EACH ROW EXECUTE FUNCTION f()`, []string{"s.t", "s.u"}},
		{`CREATE RULE r AS ON INSERT TO s.t DO NOTHING`, []string{"s.t"}},
		{`CREATE POLICY p ON s.t USING (true)`, []string{"s.t"}},
		// ALTER statements.
		{`ALTER TABLE public.nosuch ADD COLUMN b int`, []string{"public.nosuch"}},
		{`ALTER TABLE ONLY t ADD COLUMN b int`, []string{"ONLY t"}},
		{`ALTER TABLE t * ADD COLUMN b int`, []string{"t *"}},
		{`ALTER TABLE IF EXISTS t RENAME TO u`, []string{"t"}},
		{`ALTER TABLE t ATTACH PARTITION s.c FOR VALUES IN (1)`, []string{"t", "s.c"}},
		{`ALTER TABLE t DETACH PARTITION s.c`, []string{"t", "s.c"}},
		{`ALTER TABLE t INHERIT s.p`, []string{"t", "s.p"}},
		{`ALTER TABLE t NO INHERIT s.p`, []string{"t", "s.p"}},
		{`ALTER INDEX s.i RENAME TO j`, []string{"s.i"}},
		{`ALTER INDEX s.i ATTACH PARTITION s.j`, []string{"s.i", "s.j"}},
		{`ALTER SEQUENCE s.seq RESTART`, []string{"s.seq"}},
		{`ALTER SEQUENCE s.seq RENAME TO q`, []string{"s.seq"}},
		{`ALTER VIEW s.v RENAME TO w`, []string{"s.v"}},
		{`ALTER MATERIALIZED VIEW s.mv RENAME TO w`, []string{"s.mv"}},
		{`ALTER TYPE s.ct ADD ATTRIBUTE b int`, []string{"s.ct"}},
		{`ALTER TYPE s.ct RENAME ATTRIBUTE a TO b`, []string{"s.ct"}},
		{`ALTER TRIGGER tr ON s.t RENAME TO u`, []string{"s.t"}},
		{`ALTER RULE r ON s.t RENAME TO u`, []string{"s.t"}},
		{`ALTER POLICY p ON s.t RENAME TO q`, []string{"s.t"}},
		// DML and queries.
		{`INSERT INTO s.t VALUES (1)`, []string{"s.t"}},
		{`INSERT INTO s.t AS x VALUES (1)`, []string{"s.t AS x"}},
		{`UPDATE s.t SET a = 1`, []string{"s.t"}},
		{`DELETE FROM s.t AS x USING u`, []string{"s.t AS x", "u"}},
		{`SELECT * FROM s.t`, []string{"s.t"}},
		{`SELECT * FROM db.s.t`, []string{"db.s.t"}},
		{`SELECT * FROM t *`, []string{"t *"}},
		{`SELECT * FROM ONLY s.t`, []string{"ONLY s.t"}},
		{`SELECT * FROM ONLY (s.t)`, []string{"ONLY (s.t)"}},
		{`SELECT * FROM t FOR UPDATE OF s.t`, []string{"t", "s.t"}},
		{`SELECT 1 INTO x FROM t`, []string{"x", "t"}},
		{`SELECT 1 INTO TEMP s.x FROM t`, []string{"TEMP s.x", "t"}},
		{`MERGE INTO s.t USING u ON true WHEN MATCHED THEN DELETE`, []string{"s.t", "u"}},
		// Utility statements.
		{`TRUNCATE s.t, u`, []string{"s.t", "u"}},
		{`LOCK TABLE s.t`, []string{"s.t"}},
		{`VACUUM (ANALYZE) s.t (a)`, []string{"s.t"}},
		{`ANALYZE s.t`, []string{"s.t"}},
		{`CLUSTER s.t USING i`, []string{"s.t"}},
		{`CLUSTER VERBOSE s.t`, []string{"s.t"}},
		{`CLUSTER i ON s.t`, []string{"s.t"}},
		{`REINDEX TABLE s.t`, []string{"s.t"}},
		{`REFRESH MATERIALIZED VIEW s.mv`, []string{"s.mv"}},
		{`COPY s.t FROM stdin`, []string{"s.t"}},
		{`GRANT SELECT ON TABLE s.t, u TO r`, []string{"s.t", "u"}},
		{`SET CONSTRAINTS s.c, d DEFERRED`, []string{"s.c", "d"}},
	}
	for _, tc := range tests {
		t.Run(tc.sql, func(t *testing.T) {
			if got := rangeVarTexts(t, tc.sql); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("RangeVar texts = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestCreateTableAsIntoLoc pins that the IntoClause of CREATE TABLE AS spans
// create_as_target, from the name through OptTableSpace, as it does for
// CREATE MATERIALIZED VIEW.
func TestCreateTableAsIntoLoc(t *testing.T) {
	for _, tc := range []struct{ sql, want string }{
		{`CREATE TABLE s.x AS SELECT 1`, "s.x"},
		{`CREATE TABLE s.x (a, b) WITH (fillfactor = 10) AS SELECT 1, 2`, "s.x (a, b) WITH (fillfactor = 10)"},
		{`CREATE TEMP TABLE x (a) ON COMMIT DROP AS SELECT 1`, "x (a) ON COMMIT DROP"},
	} {
		list, err := Parse(tc.sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", tc.sql, err)
		}
		into := list.Items[0].(*nodes.RawStmt).Stmt.(*nodes.CreateTableAsStmt).Into
		if got := tc.sql[into.Loc.Start:into.Loc.End]; got != tc.want {
			t.Errorf("%s: IntoClause text = %q, want %q", tc.sql, got, tc.want)
		}
	}
}
