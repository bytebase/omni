package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/mariadb/ast"
)

// TestSystemTimeAccept covers the FOR SYSTEM_TIME temporal query clause on a
// base table factor (system-versioned time travel). All forms are
// container-verified against mariadb:11.8.8 (GAP: MariaDB parses, omni rejected).
//
// The clause shares its leading FOR with the SELECT-tail locking clause
// (FOR UPDATE / FOR SHARE), disambiguated by a single-token peek for SYSTEM_TIME.
func TestSystemTimeAccept(t *testing.T) {
	accept := []string{
		// AS OF expr — literal, typed literal, function, parenthesized expr
		"SELECT * FROM t FOR SYSTEM_TIME AS OF '2020-01-01'",
		"SELECT * FROM t FOR SYSTEM_TIME AS OF TIMESTAMP'2020-01-01 00:00:00'",
		"SELECT * FROM t FOR SYSTEM_TIME AS OF NOW()",
		"SELECT * FROM t FOR SYSTEM_TIME AS OF (NOW() - INTERVAL 1 YEAR)",
		// BETWEEN .. AND .. / FROM .. TO .. / ALL
		"SELECT * FROM t FOR SYSTEM_TIME BETWEEN '2020-01-01' AND '2020-06-01'",
		"SELECT * FROM t FOR SYSTEM_TIME FROM '2020-01-01' TO '2020-06-01'",
		"SELECT * FROM t FOR SYSTEM_TIME ALL",
		// Ordering: temporal binds BEFORE the alias and index hints
		"SELECT * FROM t FOR SYSTEM_TIME AS OF '2020-01-01' AS x",
		"SELECT * FROM t FOR SYSTEM_TIME AS OF '2020-01-01' x",
		"SELECT * FROM t FOR SYSTEM_TIME AS OF '2020-01-01' AS x USE INDEX (i)",
		// After PARTITION
		"SELECT * FROM t PARTITION (p0) FOR SYSTEM_TIME AS OF '2020-01-01'",
		// Coexists with the SELECT-tail locking clause (two distinct FOR clauses)
		"SELECT * FROM t FOR SYSTEM_TIME AS OF '2020-01-01' FOR UPDATE",
		"SELECT * FROM t FOR SYSTEM_TIME ALL FOR UPDATE",
		"SELECT * FROM t FOR SYSTEM_TIME ALL LOCK IN SHARE MODE",
		// Per-table in joins and comma lists
		"SELECT * FROM t FOR SYSTEM_TIME AS OF '2020-01-01' JOIN u FOR SYSTEM_TIME AS OF '2020-01-01' ON t.id = u.id",
		"SELECT * FROM t1 FOR SYSTEM_TIME AS OF '2020-01-01', t2 FOR SYSTEM_TIME ALL",
	}
	for _, sql := range accept {
		t.Run(sql, func(t *testing.T) { ParseAndCheck(t, sql) })
	}
}

// TestSystemTimeAsOfTransaction covers `FOR SYSTEM_TIME AS OF TRANSACTION n`
// (time travel by transaction id), container-verified vs mariadb:11.8.8.
func TestSystemTimeAsOfTransaction(t *testing.T) {
	ParseAndCheck(t, "SELECT * FROM t FOR SYSTEM_TIME AS OF TRANSACTION 12345")
}

// TestSystemTimeReject covers the 1064 edges (all AGREE_REJECT vs mariadb:11.8.8):
// the locking clause must keep its FOR, temporal binds before the alias only,
// it attaches to base tables (not derived), and each form needs its full spec.
func TestSystemTimeReject(t *testing.T) {
	reject := []string{
		// Locking clause is untouched: FOR UPDATE/SHARE still demand UPDATE/SHARE
		"SELECT * FROM t FOR SYSTEM TIME AS OF '2020-01-01'", // two-word form is not MariaDB
		// Alias-first is rejected (temporal must precede the alias)
		"SELECT * FROM t AS x FOR SYSTEM_TIME AS OF '2020-01-01'",
		"SELECT * FROM t x FOR SYSTEM_TIME AS OF '2020-01-01'",
		// Derived tables / subqueries take no temporal clause
		"SELECT * FROM (SELECT 1) AS d FOR SYSTEM_TIME AS OF '2020-01-01'",
		// Incomplete specifications
		"SELECT * FROM t FOR SYSTEM_TIME",
		"SELECT * FROM t FOR SYSTEM_TIME AS OF",
		"SELECT * FROM t FOR SYSTEM_TIME BETWEEN '2020-01-01'",
		"SELECT * FROM t FOR SYSTEM_TIME FROM '2020-01-01'",
	}
	for _, sql := range reject {
		t.Run(sql, func(t *testing.T) { ParseExpectError(t, sql) })
	}
}

// TestSystemTimeScope pins where FOR SYSTEM_TIME may appear. MariaDB accepts it
// on SELECT and UPDATE table references but rejects it in LOCK TABLES and in a
// single-table DELETE (all container-verified vs mariadb:11.8.8).
func TestSystemTimeScope(t *testing.T) {
	t.Run("update accepts", func(t *testing.T) {
		ParseAndCheck(t, "UPDATE t FOR SYSTEM_TIME AS OF '2020-01-01' SET id = 1")
	})
	t.Run("lock tables rejects", func(t *testing.T) {
		ParseExpectError(t, "LOCK TABLES t FOR SYSTEM_TIME AS OF '2020-01-01' READ")
	})
	t.Run("single-table delete rejects", func(t *testing.T) {
		ParseExpectError(t, "DELETE FROM t FOR SYSTEM_TIME AS OF '2020-01-01'")
	})
}

// TestSystemTimeLockingIntact guards the overload's other side: the SELECT-tail
// locking clause must still parse after the temporal change.
func TestSystemTimeLockingIntact(t *testing.T) {
	for _, sql := range []string{
		"SELECT * FROM t FOR UPDATE",
		"SELECT * FROM t LOCK IN SHARE MODE",
	} {
		t.Run(sql, func(t *testing.T) { ParseAndCheck(t, sql) })
	}
}

// TestSystemTimeAST verifies the temporal clause attaches to the TableRef with
// the right Kind and bounds for each form.
func TestSystemTimeAST(t *testing.T) {
	sysTime := func(t *testing.T, sql string) *ast.SystemTime {
		t.Helper()
		stmt := parseSeqStmt[*ast.SelectStmt](t, sql)
		ref, ok := stmt.From[0].(*ast.TableRef)
		if !ok {
			t.Fatalf("From[0] = %T, want *ast.TableRef", stmt.From[0])
		}
		if ref.SystemTime == nil {
			t.Fatalf("TableRef.SystemTime = nil, want non-nil")
		}
		return ref.SystemTime
	}

	t.Run("as of", func(t *testing.T) {
		st := sysTime(t, "SELECT * FROM t FOR SYSTEM_TIME AS OF '2020-01-01'")
		if st.Kind != ast.SystemTimeAsOf {
			t.Errorf("Kind = %v, want SystemTimeAsOf", st.Kind)
		}
		if st.From == nil || st.To != nil {
			t.Errorf("AS OF should set From only; got From=%v To=%v", st.From, st.To)
		}
	})
	t.Run("between", func(t *testing.T) {
		st := sysTime(t, "SELECT * FROM t FOR SYSTEM_TIME BETWEEN '2020-01-01' AND '2020-06-01'")
		if st.Kind != ast.SystemTimeBetween {
			t.Errorf("Kind = %v, want SystemTimeBetween", st.Kind)
		}
		if st.From == nil || st.To == nil {
			t.Errorf("BETWEEN should set From and To; got From=%v To=%v", st.From, st.To)
		}
	})
	t.Run("from to", func(t *testing.T) {
		st := sysTime(t, "SELECT * FROM t FOR SYSTEM_TIME FROM '2020-01-01' TO '2020-06-01'")
		if st.Kind != ast.SystemTimeFromTo {
			t.Errorf("Kind = %v, want SystemTimeFromTo", st.Kind)
		}
		if st.From == nil || st.To == nil {
			t.Errorf("FROM..TO should set From and To; got From=%v To=%v", st.From, st.To)
		}
	})
	t.Run("all", func(t *testing.T) {
		st := sysTime(t, "SELECT * FROM t FOR SYSTEM_TIME ALL")
		if st.Kind != ast.SystemTimeAll {
			t.Errorf("Kind = %v, want SystemTimeAll", st.Kind)
		}
		if st.From != nil || st.To != nil {
			t.Errorf("ALL should set no bounds; got From=%v To=%v", st.From, st.To)
		}
	})
}

// TestSystemTimeOutfuncs locks the outfuncs serialization of the temporal clause
// (the catalog-field checklist: a new node field must reach NodeToString).
func TestSystemTimeOutfuncs(t *testing.T) {
	stmt := parseSeqStmt[*ast.SelectStmt](t, "SELECT * FROM t FOR SYSTEM_TIME AS OF '2020-01-01'")
	got := ast.NodeToString(stmt)
	if !strings.Contains(got, ":system_time {SYSTEM_TIME") {
		t.Errorf("NodeToString missing :system_time clause:\n%s", got)
	}
}
