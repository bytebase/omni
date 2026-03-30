package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo/ast"
)

func TestPlanCacheGetPlanCache(t *testing.T) {
	node := mustParse(t, `db.users.getPlanCache()`)
	stmt := node.(*ast.PlanCacheStatement)
	if stmt.Collection != "users" {
		t.Errorf("expected collection 'users', got %q", stmt.Collection)
	}
	if stmt.AccessMethod != "dot" {
		t.Errorf("expected access method 'dot', got %q", stmt.AccessMethod)
	}
	if len(stmt.ChainedMethods) != 0 {
		t.Errorf("expected 0 chained methods, got %d", len(stmt.ChainedMethods))
	}
}

func TestPlanCacheClear(t *testing.T) {
	node := mustParse(t, `db.users.getPlanCache().clear()`)
	stmt := node.(*ast.PlanCacheStatement)
	if len(stmt.ChainedMethods) != 1 {
		t.Fatalf("expected 1 chained method, got %d", len(stmt.ChainedMethods))
	}
	if stmt.ChainedMethods[0].Name != "clear" {
		t.Errorf("expected 'clear', got %q", stmt.ChainedMethods[0].Name)
	}
}

func TestPlanCacheList(t *testing.T) {
	node := mustParse(t, `db.users.getPlanCache().list()`)
	stmt := node.(*ast.PlanCacheStatement)
	if len(stmt.ChainedMethods) != 1 {
		t.Fatalf("expected 1 chained method, got %d", len(stmt.ChainedMethods))
	}
	if stmt.ChainedMethods[0].Name != "list" {
		t.Errorf("expected 'list', got %q", stmt.ChainedMethods[0].Name)
	}
}

func TestPlanCacheListWithArgs(t *testing.T) {
	node := mustParse(t, `db.users.getPlanCache().list([{ $match: { isActive: true } }])`)
	stmt := node.(*ast.PlanCacheStatement)
	if len(stmt.ChainedMethods) != 1 {
		t.Fatalf("expected 1 chained method, got %d", len(stmt.ChainedMethods))
	}
	if len(stmt.ChainedMethods[0].Args) != 1 {
		t.Errorf("expected 1 arg, got %d", len(stmt.ChainedMethods[0].Args))
	}
}

func TestPlanCacheGetCollectionAccess(t *testing.T) {
	node := mustParse(t, `db.getCollection("users").getPlanCache().clear()`)
	stmt := node.(*ast.PlanCacheStatement)
	if stmt.Collection != "users" {
		t.Errorf("expected collection 'users', got %q", stmt.Collection)
	}
	if stmt.AccessMethod != "getCollection" {
		t.Errorf("expected access method 'getCollection', got %q", stmt.AccessMethod)
	}
}

func TestPlanCacheBracketAccess(t *testing.T) {
	node := mustParse(t, `db["users"].getPlanCache().clear()`)
	stmt := node.(*ast.PlanCacheStatement)
	if stmt.Collection != "users" {
		t.Errorf("expected collection 'users', got %q", stmt.Collection)
	}
	if stmt.AccessMethod != "bracket" {
		t.Errorf("expected access method 'bracket', got %q", stmt.AccessMethod)
	}
}
