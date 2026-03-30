package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo/ast"
)

func TestNativeSleep(t *testing.T) {
	node := mustParse(t, `sleep(1000)`)
	stmt := node.(*ast.NativeFunctionCall)
	if stmt.Name != "sleep" {
		t.Errorf("expected name 'sleep', got %q", stmt.Name)
	}
	if len(stmt.Args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(stmt.Args))
	}
}

func TestNativeLoad(t *testing.T) {
	node := mustParse(t, `load("scripts/init.js")`)
	stmt := node.(*ast.NativeFunctionCall)
	if stmt.Name != "load" {
		t.Errorf("expected name 'load', got %q", stmt.Name)
	}
	if len(stmt.Args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(stmt.Args))
	}
}

func TestNativeQuit(t *testing.T) {
	node := mustParse(t, `quit()`)
	stmt := node.(*ast.NativeFunctionCall)
	if stmt.Name != "quit" {
		t.Errorf("expected name 'quit', got %q", stmt.Name)
	}
	if len(stmt.Args) != 0 {
		t.Errorf("expected 0 args, got %d", len(stmt.Args))
	}
}

func TestNativeQuitWithCode(t *testing.T) {
	node := mustParse(t, `quit(1)`)
	stmt := node.(*ast.NativeFunctionCall)
	if stmt.Name != "quit" {
		t.Errorf("expected name 'quit', got %q", stmt.Name)
	}
	if len(stmt.Args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(stmt.Args))
	}
}

func TestNativeExit(t *testing.T) {
	node := mustParse(t, `exit()`)
	stmt := node.(*ast.NativeFunctionCall)
	if stmt.Name != "exit" {
		t.Errorf("expected name 'exit', got %q", stmt.Name)
	}
}

func TestNativeVersion(t *testing.T) {
	node := mustParse(t, `version()`)
	stmt := node.(*ast.NativeFunctionCall)
	if stmt.Name != "version" {
		t.Errorf("expected name 'version', got %q", stmt.Name)
	}
}

func TestNativePrint(t *testing.T) {
	node := mustParse(t, `print("hello")`)
	stmt := node.(*ast.NativeFunctionCall)
	if stmt.Name != "print" {
		t.Errorf("expected name 'print', got %q", stmt.Name)
	}
}

func TestNativePrintjson(t *testing.T) {
	node := mustParse(t, `printjson({ x: 1 })`)
	stmt := node.(*ast.NativeFunctionCall)
	if stmt.Name != "printjson" {
		t.Errorf("expected name 'printjson', got %q", stmt.Name)
	}
}

func TestNativeCls(t *testing.T) {
	node := mustParse(t, `cls()`)
	stmt := node.(*ast.NativeFunctionCall)
	if stmt.Name != "cls" {
		t.Errorf("expected name 'cls', got %q", stmt.Name)
	}
}

func TestNativeLoc(t *testing.T) {
	node := mustParse(t, `sleep(100)`)
	assertLoc(t, node, 0, 10)
}
