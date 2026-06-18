package parser

import (
	"testing"

	"github.com/bytebase/omni/mariadb/ast"
)

// TestSequenceValueAccept covers NEXT/PREVIOUS VALUE FOR in every position the
// container parse-accepts (BYT-9135). Per the position-matrix recon, MariaDB's
// grammar admits it as a universal primary expression — the CHECK/generated-col/
// partition restrictions are semantic (1970/1901/1564), not 1064 — so the parser
// must accept all of these.
func TestSequenceValueAccept(t *testing.T) {
	accept := []string{
		"SELECT NEXT VALUE FOR sq",
		"SELECT PREVIOUS VALUE FOR sq",
		"SELECT NEXT VALUE FOR sq + 1",
		"SELECT NEXT VALUE FOR s.sq",
		"INSERT INTO pt (id, v) VALUES (NEXT VALUE FOR sq, 1)",
		"INSERT INTO pt (id) VALUES (PREVIOUS VALUE FOR sq)",
		"UPDATE pt SET id = NEXT VALUE FOR sq WHERE v = 1",
		"SELECT * FROM pt WHERE id = NEXT VALUE FOR sq",
		"SELECT id FROM pt HAVING id = NEXT VALUE FOR sq",
		"SELECT id FROM pt ORDER BY NEXT VALUE FOR sq",
		"SELECT * FROM pt WHERE id IN (NEXT VALUE FOR sq)",
		"CREATE TABLE d1 (id INT DEFAULT NEXT VALUE FOR sq, v INT)",
		"CREATE TABLE d2 (id INT DEFAULT (NEXT VALUE FOR sq), v INT)",
		"CREATE TABLE c1 (id INT CHECK (id < NEXT VALUE FOR sq))",
		"CREATE TABLE g1 (id INT, g INT AS (NEXT VALUE FOR sq))",
	}
	for _, sql := range accept {
		t.Run(sql, func(t *testing.T) { ParseAndCheck(t, sql) })
	}
}

// TestSequenceValueReject covers the only parse-level (1064) rejects: a trailing
// OVER clause (SQL-Server superset MariaDB lacks) and NEXT VALUE without FOR
// (which MariaDB also rejects, since it reserves NEXT VALUE as special syntax).
func TestSequenceValueReject(t *testing.T) {
	reject := []string{
		"SELECT NEXT VALUE FOR sq OVER (ORDER BY 1)",
		"SELECT NEXT VALUE FOR sq OVER ()",
		"SELECT next value",
		"SELECT previous value",
	}
	for _, sql := range reject {
		t.Run(sql, func(t *testing.T) { ParseExpectError(t, sql) })
	}
}

// TestSequenceValueIdentStillWorks guards that wiring the NEXT/PREVIOUS dispatch
// does not break their use as plain identifiers when not followed by VALUE.
func TestSequenceValueIdentStillWorks(t *testing.T) {
	accept := []string{
		"SELECT next FROM pt",
		"SELECT previous FROM pt",
		"SELECT id AS next FROM pt",
		"SELECT next, value FROM pt",
		"SELECT next.value FROM pt",
	}
	for _, sql := range accept {
		t.Run(sql, func(t *testing.T) { ParseAndCheck(t, sql) })
	}
}

// TestSequenceValueAST verifies the produced expression nodes and their sequence
// reference (robust via Inspect, no byte-offset coupling).
func TestSequenceValueAST(t *testing.T) {
	t.Run("next value for", func(t *testing.T) {
		nvf := findNextValueFor(t, "SELECT NEXT VALUE FOR sq")
		if nvf.Sequence == nil || nvf.Sequence.Name != "sq" {
			t.Errorf("Sequence = %+v, want sq", nvf.Sequence)
		}
	})
	t.Run("next value for schema-qualified", func(t *testing.T) {
		nvf := findNextValueFor(t, "SELECT NEXT VALUE FOR s.sq")
		if nvf.Sequence == nil || nvf.Sequence.Schema != "s" || nvf.Sequence.Name != "sq" {
			t.Errorf("Sequence = %+v, want s.sq", nvf.Sequence)
		}
	})
	t.Run("previous value for", func(t *testing.T) {
		var pvf *ast.PreviousValueForExpr
		ast.Inspect(parseOne(t, "SELECT PREVIOUS VALUE FOR sq"), func(n ast.Node) bool {
			if x, ok := n.(*ast.PreviousValueForExpr); ok {
				pvf = x
				return false
			}
			return true
		})
		if pvf == nil || pvf.Sequence == nil || pvf.Sequence.Name != "sq" {
			t.Errorf("PreviousValueForExpr = %+v", pvf)
		}
	})
	t.Run("default bare produces NextValueForExpr", func(t *testing.T) {
		if findNextValueFor(t, "CREATE TABLE d (id INT DEFAULT NEXT VALUE FOR sq)") == nil {
			t.Error("bare DEFAULT did not yield a NextValueForExpr")
		}
	})
	// Watch-item #4: PREVIOUS VALUE FOR (keyword arm) and LASTVAL (generic func)
	// are two valid spellings of last-value; both must parse.
	t.Run("lastval generic func parses", func(t *testing.T) {
		ParseAndCheck(t, "SELECT LASTVAL(sq)")
	})
}

func parseOne(t *testing.T, sql string) ast.Node {
	t.Helper()
	result := ParseAndCheck(t, sql)
	if result.Len() == 0 {
		t.Fatalf("Parse(%q) returned no statements", sql)
	}
	return result.Items[0]
}

func findNextValueFor(t *testing.T, sql string) *ast.NextValueForExpr {
	t.Helper()
	var found *ast.NextValueForExpr
	ast.Inspect(parseOne(t, sql), func(n ast.Node) bool {
		if x, ok := n.(*ast.NextValueForExpr); ok {
			found = x
			return false
		}
		return true
	})
	if found == nil {
		t.Fatalf("Parse(%q) produced no NextValueForExpr", sql)
	}
	return found
}
