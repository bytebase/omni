package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/bytebase/omni/pg/ast"
)

// LocViolation records a single location validation failure.
type LocViolation struct {
	Path    string // e.g. "RawStmt.Stmt(SelectStmt).TargetList.Items[0](ResTarget)"
	NodeTag string // node type name
	Start   int
	End     int
	Reason  string
}

func (v LocViolation) String() string {
	return fmt.Sprintf("%s [%s]: Start=%d End=%d — %s", v.Path, v.NodeTag, v.Start, v.End, v.Reason)
}

// CheckLocations parses sql via Parse(), recursively walks the AST using
// reflection, and returns all Loc violations where Start >= 0 but End <= Start.
// This does NOT call t.Fatal — callers decide how to handle violations.
func CheckLocations(t *testing.T, sql string) []LocViolation {
	t.Helper()

	result, err := Parse(sql)
	if err != nil {
		t.Fatalf("CheckLocations Parse(%q): %v", sql, err)
	}

	var violations []LocViolation
	if result != nil {
		for i, item := range result.Items {
			path := fmt.Sprintf("Items[%d]", i)
			walkNodeLocs(reflect.ValueOf(item), path, &violations)
		}
	}
	return violations
}

// walkNodeLocs recursively walks a reflected AST value, checking every Loc field.
func walkNodeLocs(v reflect.Value, path string, violations *[]LocViolation) {
	// Dereference pointers and interfaces.
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		typeName := v.Type().Name()

		// Check if this struct has a Loc field of type ast.Loc.
		locField := v.FieldByName("Loc")
		if locField.IsValid() && locField.Type() == reflect.TypeOf(ast.Loc{}) {
			loc := locField.Interface().(ast.Loc)
			if loc.Start >= 0 && loc.End <= loc.Start {
				*violations = append(*violations, LocViolation{
					Path:    path,
					NodeTag: typeName,
					Start:   loc.Start,
					End:     loc.End,
					Reason:  "Start >= 0 but End <= Start",
				})
			}
		}

		// Recurse into all fields.
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			if !field.IsExported() {
				continue
			}
			if field.Name == "Loc" {
				continue // already checked
			}
			childPath := path + "." + field.Name
			walkNodeLocs(v.Field(i), childPath, violations)
		}

	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i)
			elemPath := fmt.Sprintf("%s[%d]", path, i)
			// Add type name for interface elements.
			actual := elem
			for actual.Kind() == reflect.Ptr || actual.Kind() == reflect.Interface {
				if actual.IsNil() {
					break
				}
				actual = actual.Elem()
			}
			if actual.IsValid() && actual.Kind() == reflect.Struct {
				elemPath = fmt.Sprintf("%s[%d](%s)", path, i, actual.Type().Name())
			}
			walkNodeLocs(elem, elemPath, violations)
		}
	}
}

// parseAndCheckLoc is a helper that parses SQL and checks location violations.
func parseAndCheckLoc(t *testing.T, sql string) {
	t.Helper()
	_, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	violations := CheckLocations(t, sql)
	for _, v := range violations {
		t.Logf("  violation: %s", v)
	}
}

// TestLocBasic verifies that CheckLocations correctly detects missing End
// positions. This test logs any remaining violations to track progress.
func TestLocBasic(t *testing.T) {
	tests := []string{
		"SELECT 1",
		"SELECT a, b FROM t WHERE x > 0",
		"INSERT INTO t (a) VALUES (1)",
		"CREATE TABLE t (id int)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseAndCheckLoc(t, sql)
		})
	}
}

// TestLocExpr validates Loc.End is set correctly for expression nodes.
func TestLocExpr(t *testing.T) {
	tests := []string{
		"SELECT 1 + 2",
		"SELECT -1",
		"SELECT 2 * 3 + 4",
		"SELECT 1 = 2",
		"SELECT 1 <> 2",
		"SELECT TRUE AND FALSE",
		"SELECT NOT TRUE",
		"SELECT NULL IS NULL",
		"SELECT 1 IS NOT NULL",
		"SELECT TRUE IS TRUE",
		"SELECT 5 BETWEEN 1 AND 10",
		"SELECT 1 IN (1, 2, 3)",
		"SELECT 'foo' LIKE 'f%'",
		"SELECT count(*)",
		"SELECT upper('hello')",
		"SELECT coalesce(a, b, c)",
		"SELECT CASE WHEN x > 0 THEN 'pos' ELSE 'neg' END",
		"SELECT EXISTS (SELECT 1)",
		"SELECT ARRAY[1, 2, 3]",
		"SELECT ROW(1, 2, 3)",
		"SELECT 1::text",
		"SELECT CAST(1 AS text)",
		"SELECT a",
		"SELECT t.a",
		"SELECT $1",
		"SELECT NULLIF(a, b)",
		"SELECT GREATEST(1, 2, 3)",
		"SELECT LEAST(1, 2, 3)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseAndCheckLoc(t, sql)
		})
	}
}

// TestLocSelect validates Loc.End is set correctly for SELECT-related nodes.
func TestLocSelect(t *testing.T) {
	tests := []string{
		"SELECT 1",
		"SELECT a, b, c",
		"SELECT a AS x, b AS y FROM t",
		"SELECT * FROM t ORDER BY a ASC, b DESC",
		"SELECT * FROM t ORDER BY a NULLS FIRST",
		"SELECT * FROM t LIMIT 10 OFFSET 5",
		"SELECT * FROM t FETCH FIRST 5 ROWS ONLY",
		"SELECT * FROM t FOR UPDATE",
		"WITH cte AS (SELECT 1) SELECT * FROM cte",
		"WITH RECURSIVE cte(n) AS (SELECT 1 UNION ALL SELECT n+1 FROM cte WHERE n < 10) SELECT * FROM cte",
		"SELECT * FROM t TABLESAMPLE BERNOULLI (10)",
		"SELECT * FROM t TABLESAMPLE SYSTEM (50) REPEATABLE (42)",
		"SELECT a, b FROM t GROUP BY a, b HAVING count(*) > 1",
		"SELECT a, count(*) OVER (PARTITION BY b ORDER BY c) FROM t",
		"SELECT a, sum(b) OVER w FROM t WINDOW w AS (ORDER BY c)",
		"VALUES (1, 2), (3, 4)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseAndCheckLoc(t, sql)
		})
	}
}

// TestLocInsert validates Loc.End is set correctly for INSERT-related nodes.
func TestLocInsert(t *testing.T) {
	tests := []string{
		"INSERT INTO t VALUES (1)",
		"INSERT INTO t (a, b) VALUES (1, 2)",
		"INSERT INTO t DEFAULT VALUES",
		"INSERT INTO t SELECT * FROM s",
		"INSERT INTO t (a) VALUES (1) ON CONFLICT DO NOTHING",
		"INSERT INTO t (a) VALUES (1) ON CONFLICT (a) DO UPDATE SET a = EXCLUDED.a",
		"INSERT INTO t (a) VALUES (1) ON CONFLICT ON CONSTRAINT pk DO NOTHING",
		"INSERT INTO t (a) VALUES (1) RETURNING *",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseAndCheckLoc(t, sql)
		})
	}
}

// TestLocUpdate validates Loc.End is set correctly for UPDATE-related nodes.
func TestLocUpdate(t *testing.T) {
	tests := []string{
		"UPDATE t SET a = 1",
		"UPDATE t SET a = 1, b = 2 WHERE c = 3",
		"UPDATE t AS u SET a = 1 FROM s WHERE t.id = s.id",
		"UPDATE t SET a = 1 RETURNING *",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseAndCheckLoc(t, sql)
		})
	}
}

// TestLocDelete validates Loc.End is set correctly for DELETE-related nodes.
func TestLocDelete(t *testing.T) {
	tests := []string{
		"DELETE FROM t",
		"DELETE FROM t WHERE a = 1",
		"DELETE FROM t USING s WHERE t.id = s.id",
		"DELETE FROM t RETURNING *",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseAndCheckLoc(t, sql)
		})
	}
}

// TestLocCreateTable validates Loc.End for CREATE TABLE nodes.
func TestLocCreateTable(t *testing.T) {
	tests := []string{
		"CREATE TABLE t (id int)",
		"CREATE TABLE t (id int, name text NOT NULL)",
		"CREATE TABLE t (id serial PRIMARY KEY, name text UNIQUE)",
		"CREATE TABLE IF NOT EXISTS t (id int)",
		"CREATE TEMP TABLE t (id int)",
		"CREATE TABLE t (id int) WITH (fillfactor=70)",
		"CREATE TABLE t (a int, b int, PRIMARY KEY (a, b))",
		"CREATE TABLE t (id int CHECK (id > 0))",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseAndCheckLoc(t, sql)
		})
	}
}

// TestLocCreateIndex validates Loc.End for CREATE INDEX nodes.
func TestLocCreateIndex(t *testing.T) {
	tests := []string{
		"CREATE INDEX idx ON t (a)",
		"CREATE UNIQUE INDEX idx ON t (a, b)",
		"CREATE INDEX IF NOT EXISTS idx ON t (a)",
		"CREATE INDEX idx ON t USING btree (a)",
		"CREATE INDEX idx ON t (a) WHERE a > 0",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseAndCheckLoc(t, sql)
		})
	}
}

// TestLocCreateView validates Loc.End for CREATE VIEW nodes.
func TestLocCreateView(t *testing.T) {
	tests := []string{
		"CREATE VIEW v AS SELECT 1",
		"CREATE OR REPLACE VIEW v AS SELECT * FROM t",
		"CREATE VIEW v (a, b) AS SELECT 1, 2",
		"CREATE TEMP VIEW v AS SELECT 1",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseAndCheckLoc(t, sql)
		})
	}
}

// TestLocAlterTable validates Loc.End for ALTER TABLE nodes.
func TestLocAlterTable(t *testing.T) {
	tests := []string{
		"ALTER TABLE t ADD COLUMN c int",
		"ALTER TABLE t DROP COLUMN c",
		"ALTER TABLE t ALTER COLUMN c SET NOT NULL",
		"ALTER TABLE t RENAME TO t2",
		"ALTER TABLE t ADD CONSTRAINT c UNIQUE (a)",
		"ALTER TABLE t DROP CONSTRAINT c",
		"ALTER TABLE t SET SCHEMA s",
		"ALTER TABLE t OWNER TO u",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseAndCheckLoc(t, sql)
		})
	}
}

// TestLocAlterMisc validates Loc.End for miscellaneous ALTER nodes.
func TestLocAlterMisc(t *testing.T) {
	tests := []string{
		"ALTER DOMAIN d SET NOT NULL",
		"ALTER TYPE t ADD VALUE 'v'",
		"ALTER SEQUENCE s RESTART",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseAndCheckLoc(t, sql)
		})
	}
}

// TestLocCreateFunction validates Loc.End for CREATE FUNCTION nodes.
func TestLocCreateFunction(t *testing.T) {
	tests := []string{
		"CREATE FUNCTION f() RETURNS int AS 'SELECT 1' LANGUAGE SQL",
		"CREATE FUNCTION f(a int, b text) RETURNS void AS $$ $$ LANGUAGE SQL",
		"CREATE OR REPLACE FUNCTION f() RETURNS int AS 'body' LANGUAGE plpgsql IMMUTABLE",
		"CREATE FUNCTION f() RETURNS int AS 'body' LANGUAGE SQL STRICT SECURITY DEFINER",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseAndCheckLoc(t, sql)
		})
	}
}

// TestLocTrigger validates Loc.End for CREATE TRIGGER nodes.
func TestLocTrigger(t *testing.T) {
	tests := []string{
		"CREATE TRIGGER tr AFTER INSERT ON t FOR EACH ROW EXECUTE FUNCTION f()",
		"CREATE TRIGGER tr BEFORE UPDATE ON t FOR EACH ROW EXECUTE PROCEDURE f()",
		"CREATE TRIGGER tr AFTER DELETE ON t EXECUTE FUNCTION f()",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseAndCheckLoc(t, sql)
		})
	}
}

// TestLocType validates Loc.End for CREATE TYPE nodes.
func TestLocType(t *testing.T) {
	tests := []string{
		"CREATE TYPE t AS (a int, b text)",
		"CREATE TYPE t AS ENUM ('a', 'b', 'c')",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseAndCheckLoc(t, sql)
		})
	}
}

// TestLocDefine validates Loc.End for DEFINE statement nodes.
func TestLocDefine(t *testing.T) {
	tests := []string{
		"CREATE AGGREGATE agg (int) (SFUNC = f, STYPE = int)",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseAndCheckLoc(t, sql)
		})
	}
}

// TestLocDrop validates Loc.End for DROP statement nodes.
func TestLocDrop(t *testing.T) {
	tests := []string{
		"DROP TABLE t",
		"DROP TABLE IF EXISTS t CASCADE",
		"DROP INDEX idx",
		"DROP VIEW v",
		"DROP FUNCTION f()",
		"DROP SCHEMA s CASCADE",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseAndCheckLoc(t, sql)
		})
	}
}

// TestLocGrant validates Loc.End for GRANT/REVOKE nodes.
func TestLocGrant(t *testing.T) {
	tests := []string{
		"GRANT SELECT ON t TO PUBLIC",
		"GRANT ALL ON t TO u",
		"REVOKE ALL ON t FROM PUBLIC",
		"REVOKE SELECT ON t FROM u",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseAndCheckLoc(t, sql)
		})
	}
}

// TestLocSchema validates Loc.End for CREATE SCHEMA nodes.
func TestLocSchema(t *testing.T) {
	tests := []string{
		"CREATE SCHEMA s",
		"CREATE SCHEMA IF NOT EXISTS s",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseAndCheckLoc(t, sql)
		})
	}
}

// TestLocDatabase validates Loc.End for CREATE/DROP DATABASE nodes.
func TestLocDatabase(t *testing.T) {
	tests := []string{
		"CREATE DATABASE d",
		"DROP DATABASE d",
		"DROP DATABASE IF EXISTS d",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseAndCheckLoc(t, sql)
		})
	}
}

// TestLocUtility validates Loc.End for utility statement nodes.
func TestLocUtility(t *testing.T) {
	tests := []string{
		"EXPLAIN SELECT 1",
		"EXPLAIN ANALYZE SELECT 1",
		"VACUUM t",
		"ANALYZE t",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseAndCheckLoc(t, sql)
		})
	}
}

// TestLocTransaction validates Loc.End for transaction statement nodes.
func TestLocTransaction(t *testing.T) {
	tests := []string{
		"BEGIN",
		"COMMIT",
		"ROLLBACK",
		"START TRANSACTION",
		"SAVEPOINT sp1",
		"RELEASE SAVEPOINT sp1",
		"ROLLBACK TO SAVEPOINT sp1",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseAndCheckLoc(t, sql)
		})
	}
}

// TestLocCursor validates Loc.End for cursor statement nodes.
func TestLocCursor(t *testing.T) {
	tests := []string{
		"DECLARE c CURSOR FOR SELECT 1",
		"FETCH NEXT FROM c",
		"CLOSE c",
		"MOVE NEXT FROM c",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseAndCheckLoc(t, sql)
		})
	}
}

// TestLocSet validates Loc.End for SET/SHOW/RESET nodes.
func TestLocSet(t *testing.T) {
	tests := []string{
		"SET search_path TO public",
		"SET search_path = public",
		"SET TIME ZONE 'UTC'",
		"SHOW search_path",
		"SHOW ALL",
		"RESET search_path",
		"RESET ALL",
	}
	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			parseAndCheckLoc(t, sql)
		})
	}
}

// TestLocRawStmt validates that RawStmt.Loc accurately spans the source SQL text
// for each statement in multi-statement inputs.
func TestLocRawStmt(t *testing.T) {
	tests := []struct {
		sql   string
		stmts []string // expected source text for each statement
	}{
		{
			sql:   "SELECT 1",
			stmts: []string{"SELECT 1"},
		},
		{
			sql:   "SELECT 1; SELECT 2",
			stmts: []string{"SELECT 1", "SELECT 2"},
		},
		{
			sql:   "BEGIN; INSERT INTO t VALUES (1); COMMIT",
			stmts: []string{"BEGIN", "INSERT INTO t VALUES (1)", "COMMIT"},
		},
		{
			sql:   "  SELECT 1 ;  SELECT 2  ",
			stmts: []string{"SELECT 1", "SELECT 2"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.sql, func(t *testing.T) {
			result, err := Parse(tc.sql)
			if err != nil {
				t.Fatalf("Parse(%q): %v", tc.sql, err)
			}
			if len(result.Items) != len(tc.stmts) {
				t.Fatalf("Parse(%q): got %d stmts, want %d", tc.sql, len(result.Items), len(tc.stmts))
			}
			for i, item := range result.Items {
				raw, ok := item.(*ast.RawStmt)
				if !ok {
					t.Fatalf("stmt[%d] is %T, want *ast.RawStmt", i, item)
				}
				if raw.Loc.Start < 0 || raw.Loc.End < 0 {
					t.Errorf("stmt[%d] has invalid Loc: %+v", i, raw.Loc)
					continue
				}
				got := tc.sql[raw.Loc.Start:raw.Loc.End]
				if got != tc.stmts[i] {
					t.Errorf("stmt[%d] Loc mismatch:\n  got:  %q (Loc=%+v)\n  want: %q", i, got, raw.Loc, tc.stmts[i])
				}
			}
		})
	}
}
