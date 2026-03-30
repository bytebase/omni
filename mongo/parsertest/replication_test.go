package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo/ast"
)

func TestRsStatus(t *testing.T) {
	node := mustParse(t, `rs.status()`)
	stmt := node.(*ast.RsStatement)
	if stmt.MethodName != "status" {
		t.Errorf("expected method 'status', got %q", stmt.MethodName)
	}
	if len(stmt.Args) != 0 {
		t.Errorf("expected 0 args, got %d", len(stmt.Args))
	}
}

func TestRsStatusWithArgs(t *testing.T) {
	node := mustParse(t, `rs.status({ initialSync: 1 })`)
	stmt := node.(*ast.RsStatement)
	if stmt.MethodName != "status" {
		t.Errorf("expected method 'status', got %q", stmt.MethodName)
	}
	if len(stmt.Args) != 1 {
		t.Errorf("expected 1 arg, got %d", len(stmt.Args))
	}
}

func TestRsInitiate(t *testing.T) {
	node := mustParse(t, `rs.initiate()`)
	stmt := node.(*ast.RsStatement)
	if stmt.MethodName != "initiate" {
		t.Errorf("expected method 'initiate', got %q", stmt.MethodName)
	}
}

func TestRsInitiateWithConfig(t *testing.T) {
	node := mustParse(t, `rs.initiate({ _id: "myRS", members: [{ _id: 0, host: "mongo1:27017" }] })`)
	stmt := node.(*ast.RsStatement)
	if stmt.MethodName != "initiate" {
		t.Errorf("expected method 'initiate', got %q", stmt.MethodName)
	}
	if len(stmt.Args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(stmt.Args))
	}
}

func TestRsConf(t *testing.T) {
	node := mustParse(t, `rs.conf()`)
	stmt := node.(*ast.RsStatement)
	if stmt.MethodName != "conf" {
		t.Errorf("expected method 'conf', got %q", stmt.MethodName)
	}
}

func TestRsReconfig(t *testing.T) {
	node := mustParse(t, `rs.reconfig({ _id: "rs0", members: [] })`)
	stmt := node.(*ast.RsStatement)
	if stmt.MethodName != "reconfig" {
		t.Errorf("expected method 'reconfig', got %q", stmt.MethodName)
	}
}

func TestRsAdd(t *testing.T) {
	node := mustParse(t, `rs.add("mongo4:27017")`)
	stmt := node.(*ast.RsStatement)
	if stmt.MethodName != "add" {
		t.Errorf("expected method 'add', got %q", stmt.MethodName)
	}
}

func TestRsFreeze(t *testing.T) {
	node := mustParse(t, `rs.freeze(120)`)
	stmt := node.(*ast.RsStatement)
	if stmt.MethodName != "freeze" {
		t.Errorf("expected method 'freeze', got %q", stmt.MethodName)
	}
}

func TestRsStepDown(t *testing.T) {
	node := mustParse(t, `rs.stepDown(60)`)
	stmt := node.(*ast.RsStatement)
	if stmt.MethodName != "stepDown" {
		t.Errorf("expected method 'stepDown', got %q", stmt.MethodName)
	}
}

func TestRsLoc(t *testing.T) {
	node := mustParse(t, `rs.status()`)
	assertLoc(t, node, 0, 11)
}
