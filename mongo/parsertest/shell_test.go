package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo/ast"
)

func TestShowDbs(t *testing.T) {
	node := mustParse(t, "show dbs")
	cmd := node.(*ast.ShowCommand)
	if cmd.Target != "databases" {
		t.Errorf("expected 'databases', got %q", cmd.Target)
	}
}

func TestShowDatabases(t *testing.T) {
	node := mustParse(t, "show databases")
	cmd := node.(*ast.ShowCommand)
	if cmd.Target != "databases" {
		t.Errorf("expected 'databases', got %q", cmd.Target)
	}
}

func TestShowCollections(t *testing.T) {
	node := mustParse(t, "show collections")
	cmd := node.(*ast.ShowCommand)
	if cmd.Target != "collections" {
		t.Errorf("expected 'collections', got %q", cmd.Target)
	}
}

func TestShowTables(t *testing.T) {
	node := mustParse(t, "show tables")
	cmd := node.(*ast.ShowCommand)
	if cmd.Target != "tables" {
		t.Errorf("expected 'tables', got %q", cmd.Target)
	}
}

func TestShowProfile(t *testing.T) {
	node := mustParse(t, "show profile")
	cmd := node.(*ast.ShowCommand)
	if cmd.Target != "profile" {
		t.Errorf("expected 'profile', got %q", cmd.Target)
	}
}

func TestShowUsers(t *testing.T) {
	node := mustParse(t, "show users")
	cmd := node.(*ast.ShowCommand)
	if cmd.Target != "users" {
		t.Errorf("expected 'users', got %q", cmd.Target)
	}
}

func TestShowRoles(t *testing.T) {
	node := mustParse(t, "show roles")
	cmd := node.(*ast.ShowCommand)
	if cmd.Target != "roles" {
		t.Errorf("expected 'roles', got %q", cmd.Target)
	}
}

func TestShowLog(t *testing.T) {
	node := mustParse(t, "show log")
	cmd := node.(*ast.ShowCommand)
	if cmd.Target != "log" {
		t.Errorf("expected 'log', got %q", cmd.Target)
	}
}

func TestShowInvalidTarget(t *testing.T) {
	mustFail(t, "show something")
}

func TestMultipleShowCommands(t *testing.T) {
	nodes := mustParseN(t, "show dbs; show collections", 2)
	cmd1 := nodes[0].(*ast.ShowCommand)
	cmd2 := nodes[1].(*ast.ShowCommand)
	if cmd1.Target != "databases" {
		t.Errorf("expected 'databases', got %q", cmd1.Target)
	}
	if cmd2.Target != "collections" {
		t.Errorf("expected 'collections', got %q", cmd2.Target)
	}
}

func TestShowCommandLoc(t *testing.T) {
	node := mustParse(t, "show dbs")
	assertLoc(t, node, 0, 8)
}
