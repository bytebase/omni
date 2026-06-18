package parser

import (
	"testing"

	"github.com/bytebase/omni/mariadb/ast"
)

// TestReturningAccept covers the RETURNING accept surface (BYT-9135 P2),
// container-verified vs mariadb:11.8.8: INSERT (all forms), REPLACE (the
// donorless arm), and single-table DELETE — with *, expr+alias, subquery, and
// aggregates (which parse and are semantic-rejected, i.e. in the parser's
// contract).
func TestReturningAccept(t *testing.T) {
	accept := []string{
		// INSERT — VALUES / list+expr+alias / star / SELECT / paren-query / SET / ON DUP KEY
		"INSERT INTO rt (id, name) VALUES (1, 'a') RETURNING id",
		"INSERT INTO rt (id, name) VALUES (2, 'b') RETURNING id, name, id + 100 AS big",
		"INSERT INTO rt (id, name) VALUES (3, 'c') RETURNING *",
		"INSERT INTO rt SELECT 9, 'q' RETURNING id, name",
		"INSERT INTO rt (SELECT 1, 'x') RETURNING id",
		"INSERT INTO rt SET id = 1, name = 'x' RETURNING id",
		"INSERT INTO rt (id, name) VALUES (4, 'd') ON DUPLICATE KEY UPDATE name = 'x' RETURNING id",
		// RETURNING list with subquery / aggregate (parse-accept)
		"INSERT INTO rt (id) VALUES (1) RETURNING id, (SELECT COUNT(*) FROM rt2)",
		"INSERT INTO rt (id) VALUES (1) RETURNING COUNT(*)",
		// REPLACE (donorless — pg has no REPLACE)
		"REPLACE INTO rt (id, name) VALUES (1, 'z') RETURNING id, name",
		"REPLACE INTO rt (id, name) VALUES (1, 'z') RETURNING *",
		// single-table DELETE
		"DELETE FROM rt WHERE id = 3 RETURNING id, name",
		"DELETE FROM rt WHERE id = 2 RETURNING *",
		"DELETE FROM rt RETURNING COUNT(*)",
	}
	for _, sql := range accept {
		t.Run(sql, func(t *testing.T) { ParseAndCheck(t, sql) })
	}
}

// TestReturningReject covers the 1064 edges (all AGREE_REJECT vs the container):
// MariaDB has no UPDATE RETURNING, RETURNING is single-table-DELETE only, and a
// bare RETURNING needs a select-list.
func TestReturningReject(t *testing.T) {
	reject := []string{
		"UPDATE rt SET name = 'x' WHERE id = 1 RETURNING id",                 // no UPDATE RETURNING
		"DELETE rt FROM rt JOIN rt2 ON rt.id = rt2.id RETURNING rt.id",       // multi-table delete (syntax 1, JOIN)
		"DELETE FROM rt USING rt JOIN rt2 ON rt.id = rt2.id RETURNING rt.id", // multi-table delete (syntax 2, USING)
		"INSERT INTO rt (id) VALUES (1) RETURNING",                           // bare RETURNING (no list)
	}
	for _, sql := range reject {
		t.Run(sql, func(t *testing.T) { ParseExpectError(t, sql) })
	}
}

// TestReturningAST verifies the Returning list is populated on INSERT/REPLACE
// (shared InsertStmt) and DELETE.
func TestReturningAST(t *testing.T) {
	t.Run("insert returning list", func(t *testing.T) {
		stmt := parseSeqStmt[*ast.InsertStmt](t, "INSERT INTO rt (id, name) VALUES (1, 'a') RETURNING id, name")
		if len(stmt.Returning) != 2 {
			t.Errorf("len(Returning) = %d, want 2", len(stmt.Returning))
		}
	})
	t.Run("insert returning star", func(t *testing.T) {
		stmt := parseSeqStmt[*ast.InsertStmt](t, "INSERT INTO rt (id) VALUES (1) RETURNING *")
		if len(stmt.Returning) != 1 {
			t.Errorf("len(Returning) = %d, want 1", len(stmt.Returning))
		}
	})
	t.Run("replace returning", func(t *testing.T) {
		stmt := parseSeqStmt[*ast.InsertStmt](t, "REPLACE INTO rt (id) VALUES (1) RETURNING id")
		if !stmt.IsReplace || len(stmt.Returning) != 1 {
			t.Errorf("IsReplace=%v len(Returning)=%d, want true/1", stmt.IsReplace, len(stmt.Returning))
		}
	})
	t.Run("delete returning", func(t *testing.T) {
		stmt := parseSeqStmt[*ast.DeleteStmt](t, "DELETE FROM rt WHERE id = 1 RETURNING id, name")
		if len(stmt.Returning) != 2 {
			t.Errorf("len(Returning) = %d, want 2", len(stmt.Returning))
		}
	})
}
