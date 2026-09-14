package parser

import (
	"testing"

	"github.com/bytebase/omni/mariadb/ast"
)

// TestBlockCommentDoesNotNest_Oracle proves against a live MariaDB 11.8
// container that omni's parser now agrees with the server: a
// "/* /* */ ... */"-shaped comment closes at the FIRST */, so text after it
// is a live, separate statement — not comment content the way the old
// depth-nesting scanner read it. The oracle runs the raw SQL on a
// multiStatements connection so the server's own comment handling decides
// what actually executes; omni's Split and Parse must describe the same two
// statements.
func TestBlockCommentDoesNotNest_Oracle(t *testing.T) {
	o := startMariaDB(t)

	if _, err := o.db.ExecContext(o.ctx, "CREATE TABLE t1 (id INT)"); err != nil {
		t.Fatalf("create t1: %v", err)
	}
	if _, err := o.db.ExecContext(o.ctx, "CREATE TABLE t2 (id INT)"); err != nil {
		t.Fatalf("create t2: %v", err)
	}

	statement := "INSERT INTO t1 VALUES (1) /* /* */; DROP TABLE t2; -- */"

	// Ground truth: what does the real server do with this text?
	if _, err := o.db.ExecContext(o.ctx, statement); err != nil {
		t.Fatalf("server rejected the oracle statement outright: %v", err)
	}
	var t2Count int
	err := o.db.QueryRowContext(o.ctx, `
		SELECT COUNT(*) FROM information_schema.tables
		WHERE table_schema = 'test' AND table_name = 't2'
	`).Scan(&t2Count)
	if err != nil {
		t.Fatalf("check t2 existence: %v", err)
	}
	if t2Count != 0 {
		t.Fatalf("server did not treat %q as two statements (t2 still exists) — the oracle assumption is wrong", statement)
	}

	// omni must describe the same two statements the server actually ran.
	segments := Split(statement)
	if len(segments) != 2 {
		t.Fatalf("Split(%q) returned %d segments, want 2 (matching the server, which ran an INSERT then a DROP TABLE)", statement, len(segments))
	}

	list, err := Parse(statement)
	if err != nil {
		t.Fatalf("Parse(%q): %v", statement, err)
	}
	if len(list.Items) != 2 {
		t.Fatalf("Parse(%q) returned %d statements, want 2", statement, len(list.Items))
	}
	if _, ok := list.Items[0].(*ast.InsertStmt); !ok {
		t.Errorf("statement 1: got %T, want *ast.InsertStmt", list.Items[0])
	}
	if _, ok := list.Items[1].(*ast.DropTableStmt); !ok {
		t.Errorf("statement 2: got %T, want *ast.DropTableStmt", list.Items[1])
	}
}

// TestBlockCommentDoesNotHidePredicate_Oracle proves against a live MariaDB
// 11.8 container that a predicate following a "/* /* */ ... */"-shaped
// comment is live and takes effect, matching omni's parsed WHERE clause. The
// old depth-nesting scanner read the OR clause as comment content, so omni's
// AST would have shown only the first half of the predicate.
func TestBlockCommentDoesNotHidePredicate_Oracle(t *testing.T) {
	o := startMariaDB(t)

	if _, err := o.db.ExecContext(o.ctx, "CREATE TABLE t3 (id INT)"); err != nil {
		t.Fatalf("create t3: %v", err)
	}
	if _, err := o.db.ExecContext(o.ctx, "INSERT INTO t3 VALUES (1), (2), (3)"); err != nil {
		t.Fatalf("seed t3: %v", err)
	}

	// id=100 never matches; the row for id=1 is only deleted if the OR clause
	// is live rather than swallowed by the comment.
	statement := "DELETE FROM t3 WHERE id = 100 /* /* */ OR id = 1 -- */"

	result, err := o.db.ExecContext(o.ctx, statement)
	if err != nil {
		t.Fatalf("server rejected the oracle statement outright: %v", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		t.Fatalf("rows affected: %v", err)
	}
	if affected != 1 {
		t.Fatalf("server deleted %d rows, want 1 (id=1) — the oracle assumption that the OR clause is live is wrong", affected)
	}

	list, err := Parse(statement)
	if err != nil {
		t.Fatalf("Parse(%q): %v", statement, err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("Parse(%q) returned %d statements, want 1", statement, len(list.Items))
	}
	del, ok := list.Items[0].(*ast.DeleteStmt)
	if !ok {
		t.Fatalf("expected *ast.DeleteStmt, got %T", list.Items[0])
	}
	where, ok := del.Where.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected the WHERE clause to be a binary OR expression, got %T (comment may have swallowed the OR clause)", del.Where)
	}
	if where.Op != ast.BinOpOr {
		t.Errorf("expected top-level WHERE operator OR, got %v", where.Op)
	}
}
