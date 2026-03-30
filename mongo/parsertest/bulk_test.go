package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo/ast"
)

func TestBulkOrderedInit(t *testing.T) {
	node := mustParse(t, `db.users.initializeOrderedBulkOp()`)
	stmt := node.(*ast.BulkStatement)
	if stmt.Collection != "users" {
		t.Errorf("expected collection 'users', got %q", stmt.Collection)
	}
	if !stmt.Ordered {
		t.Error("expected ordered to be true")
	}
	if stmt.AccessMethod != "dot" {
		t.Errorf("expected access method 'dot', got %q", stmt.AccessMethod)
	}
	if len(stmt.Operations) != 0 {
		t.Errorf("expected 0 operations, got %d", len(stmt.Operations))
	}
}

func TestBulkUnorderedInit(t *testing.T) {
	node := mustParse(t, `db.users.initializeUnorderedBulkOp()`)
	stmt := node.(*ast.BulkStatement)
	if stmt.Ordered {
		t.Error("expected ordered to be false")
	}
}

func TestBulkInsert(t *testing.T) {
	node := mustParse(t, `db.users.initializeUnorderedBulkOp().insert({ name: "alice" })`)
	stmt := node.(*ast.BulkStatement)
	if len(stmt.Operations) != 1 {
		t.Fatalf("expected 1 operation, got %d", len(stmt.Operations))
	}
	if stmt.Operations[0].Method != "insert" {
		t.Errorf("expected operation 'insert', got %q", stmt.Operations[0].Method)
	}
}

func TestBulkMultipleInserts(t *testing.T) {
	node := mustParse(t, `db.users.initializeOrderedBulkOp().insert({ name: "bob" }).insert({ name: "charlie" })`)
	stmt := node.(*ast.BulkStatement)
	if len(stmt.Operations) != 2 {
		t.Fatalf("expected 2 operations, got %d", len(stmt.Operations))
	}
	if stmt.Operations[0].Method != "insert" || stmt.Operations[1].Method != "insert" {
		t.Error("expected both operations to be 'insert'")
	}
}

func TestBulkInsertAndExecute(t *testing.T) {
	node := mustParse(t, `db.users.initializeOrderedBulkOp().insert({ name: "grace" }).execute()`)
	stmt := node.(*ast.BulkStatement)
	if len(stmt.Operations) != 2 {
		t.Fatalf("expected 2 operations, got %d", len(stmt.Operations))
	}
	if stmt.Operations[1].Method != "execute" {
		t.Errorf("expected 'execute', got %q", stmt.Operations[1].Method)
	}
}

func TestBulkExecuteWithWriteConcern(t *testing.T) {
	node := mustParse(t, `db.users.initializeOrderedBulkOp().insert({ name: "bob" }).execute({ w: "majority" })`)
	stmt := node.(*ast.BulkStatement)
	if len(stmt.Operations) != 2 {
		t.Fatalf("expected 2 operations, got %d", len(stmt.Operations))
	}
	if stmt.Operations[1].Method != "execute" {
		t.Errorf("expected 'execute', got %q", stmt.Operations[1].Method)
	}
	if len(stmt.Operations[1].Args) != 1 {
		t.Errorf("expected 1 arg for execute, got %d", len(stmt.Operations[1].Args))
	}
}

func TestBulkFindRemove(t *testing.T) {
	node := mustParse(t, `db.users.initializeUnorderedBulkOp().find({ status: "inactive" }).remove().execute()`)
	stmt := node.(*ast.BulkStatement)
	if len(stmt.Operations) != 3 {
		t.Fatalf("expected 3 operations, got %d", len(stmt.Operations))
	}
	if stmt.Operations[0].Method != "find" {
		t.Errorf("expected 'find', got %q", stmt.Operations[0].Method)
	}
	if stmt.Operations[1].Method != "remove" {
		t.Errorf("expected 'remove', got %q", stmt.Operations[1].Method)
	}
	if stmt.Operations[2].Method != "execute" {
		t.Errorf("expected 'execute', got %q", stmt.Operations[2].Method)
	}
}

func TestBulkGetCollectionAccess(t *testing.T) {
	node := mustParse(t, `db.getCollection("users").initializeOrderedBulkOp()`)
	stmt := node.(*ast.BulkStatement)
	if stmt.Collection != "users" {
		t.Errorf("expected collection 'users', got %q", stmt.Collection)
	}
	if stmt.AccessMethod != "getCollection" {
		t.Errorf("expected access method 'getCollection', got %q", stmt.AccessMethod)
	}
}

func TestBulkBracketAccess(t *testing.T) {
	node := mustParse(t, `db["users"].initializeOrderedBulkOp()`)
	stmt := node.(*ast.BulkStatement)
	if stmt.Collection != "users" {
		t.Errorf("expected collection 'users', got %q", stmt.Collection)
	}
	if stmt.AccessMethod != "bracket" {
		t.Errorf("expected access method 'bracket', got %q", stmt.AccessMethod)
	}
}
