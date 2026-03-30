package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo/ast"
)

func TestShStatus(t *testing.T) {
	node := mustParse(t, `sh.status()`)
	stmt := node.(*ast.ShStatement)
	if stmt.MethodName != "status" {
		t.Errorf("expected method 'status', got %q", stmt.MethodName)
	}
	if len(stmt.Args) != 0 {
		t.Errorf("expected 0 args, got %d", len(stmt.Args))
	}
}

func TestShAddShard(t *testing.T) {
	node := mustParse(t, `sh.addShard("rs1/mongo1:27017")`)
	stmt := node.(*ast.ShStatement)
	if stmt.MethodName != "addShard" {
		t.Errorf("expected method 'addShard', got %q", stmt.MethodName)
	}
	if len(stmt.Args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(stmt.Args))
	}
}

func TestShEnableSharding(t *testing.T) {
	node := mustParse(t, `sh.enableSharding("mydb")`)
	stmt := node.(*ast.ShStatement)
	if stmt.MethodName != "enableSharding" {
		t.Errorf("expected method 'enableSharding', got %q", stmt.MethodName)
	}
}

func TestShShardCollection(t *testing.T) {
	node := mustParse(t, `sh.shardCollection("mydb.orders", { _id: 1 })`)
	stmt := node.(*ast.ShStatement)
	if stmt.MethodName != "shardCollection" {
		t.Errorf("expected method 'shardCollection', got %q", stmt.MethodName)
	}
	if len(stmt.Args) != 2 {
		t.Fatalf("expected 2 args, got %d", len(stmt.Args))
	}
}

func TestShGetBalancerState(t *testing.T) {
	node := mustParse(t, `sh.getBalancerState()`)
	stmt := node.(*ast.ShStatement)
	if stmt.MethodName != "getBalancerState" {
		t.Errorf("expected method 'getBalancerState', got %q", stmt.MethodName)
	}
}

func TestShStartBalancer(t *testing.T) {
	node := mustParse(t, `sh.startBalancer()`)
	stmt := node.(*ast.ShStatement)
	if stmt.MethodName != "startBalancer" {
		t.Errorf("expected method 'startBalancer', got %q", stmt.MethodName)
	}
}

func TestShStopBalancer(t *testing.T) {
	node := mustParse(t, `sh.stopBalancer()`)
	stmt := node.(*ast.ShStatement)
	if stmt.MethodName != "stopBalancer" {
		t.Errorf("expected method 'stopBalancer', got %q", stmt.MethodName)
	}
}

func TestShMoveChunk(t *testing.T) {
	node := mustParse(t, `sh.moveChunk("mydb.orders", { _id: 1 }, "shard2")`)
	stmt := node.(*ast.ShStatement)
	if stmt.MethodName != "moveChunk" {
		t.Errorf("expected method 'moveChunk', got %q", stmt.MethodName)
	}
	if len(stmt.Args) != 3 {
		t.Fatalf("expected 3 args, got %d", len(stmt.Args))
	}
}

func TestShLoc(t *testing.T) {
	node := mustParse(t, `sh.status()`)
	assertLoc(t, node, 0, 11)
}
