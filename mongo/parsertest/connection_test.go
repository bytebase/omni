package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo/ast"
)

func TestMongoBasic(t *testing.T) {
	node := mustParse(t, `Mongo()`)
	stmt := node.(*ast.ConnectionStatement)
	if stmt.Constructor != "Mongo" {
		t.Errorf("expected constructor 'Mongo', got %q", stmt.Constructor)
	}
	if len(stmt.Args) != 0 {
		t.Errorf("expected 0 args, got %d", len(stmt.Args))
	}
	if len(stmt.ChainedMethods) != 0 {
		t.Errorf("expected 0 chained methods, got %d", len(stmt.ChainedMethods))
	}
}

func TestMongoWithArg(t *testing.T) {
	node := mustParse(t, `Mongo("localhost:27017")`)
	stmt := node.(*ast.ConnectionStatement)
	if stmt.Constructor != "Mongo" {
		t.Errorf("expected constructor 'Mongo', got %q", stmt.Constructor)
	}
	if len(stmt.Args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(stmt.Args))
	}
}

func TestMongoGetDB(t *testing.T) {
	node := mustParse(t, `Mongo().getDB("test")`)
	stmt := node.(*ast.ConnectionStatement)
	if stmt.Constructor != "Mongo" {
		t.Errorf("expected constructor 'Mongo', got %q", stmt.Constructor)
	}
	if len(stmt.ChainedMethods) != 1 {
		t.Fatalf("expected 1 chained method, got %d", len(stmt.ChainedMethods))
	}
	if stmt.ChainedMethods[0].Name != "getDB" {
		t.Errorf("expected chained method 'getDB', got %q", stmt.ChainedMethods[0].Name)
	}
}

func TestMongoChainMultiple(t *testing.T) {
	node := mustParse(t, `Mongo("localhost").getDB("mydb").getCollectionNames()`)
	stmt := node.(*ast.ConnectionStatement)
	if len(stmt.ChainedMethods) != 2 {
		t.Fatalf("expected 2 chained methods, got %d", len(stmt.ChainedMethods))
	}
	if stmt.ChainedMethods[0].Name != "getDB" {
		t.Errorf("expected 'getDB', got %q", stmt.ChainedMethods[0].Name)
	}
	if stmt.ChainedMethods[1].Name != "getCollectionNames" {
		t.Errorf("expected 'getCollectionNames', got %q", stmt.ChainedMethods[1].Name)
	}
}

func TestConnectBasic(t *testing.T) {
	node := mustParse(t, `connect()`)
	stmt := node.(*ast.ConnectionStatement)
	if stmt.Constructor != "connect" {
		t.Errorf("expected constructor 'connect', got %q", stmt.Constructor)
	}
}

func TestConnectWithURL(t *testing.T) {
	node := mustParse(t, `connect("mongodb://localhost:27017/test")`)
	stmt := node.(*ast.ConnectionStatement)
	if stmt.Constructor != "connect" {
		t.Errorf("expected constructor 'connect', got %q", stmt.Constructor)
	}
	if len(stmt.Args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(stmt.Args))
	}
}

func TestConnectWithChain(t *testing.T) {
	node := mustParse(t, `connect("localhost").getDB("test")`)
	stmt := node.(*ast.ConnectionStatement)
	if stmt.Constructor != "connect" {
		t.Errorf("expected constructor 'connect', got %q", stmt.Constructor)
	}
	if len(stmt.ChainedMethods) != 1 {
		t.Fatalf("expected 1 chained method, got %d", len(stmt.ChainedMethods))
	}
}

func TestMongoLoc(t *testing.T) {
	node := mustParse(t, `Mongo()`)
	assertLoc(t, node, 0, 7)
}
