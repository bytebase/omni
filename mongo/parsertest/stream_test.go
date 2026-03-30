package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo/ast"
)

func TestSpCreateStreamProcessor(t *testing.T) {
	node := mustParse(t, `sp.createStreamProcessor("myProc", { pipeline: [] })`)
	stmt := node.(*ast.SpStatement)
	if stmt.MethodName != "createStreamProcessor" {
		t.Errorf("expected method 'createStreamProcessor', got %q", stmt.MethodName)
	}
	if stmt.SubMethod != "" {
		t.Errorf("expected no sub-method, got %q", stmt.SubMethod)
	}
	if len(stmt.Args) != 2 {
		t.Fatalf("expected 2 args, got %d", len(stmt.Args))
	}
}

func TestSpListStreamProcessors(t *testing.T) {
	node := mustParse(t, `sp.listStreamProcessors()`)
	stmt := node.(*ast.SpStatement)
	if stmt.MethodName != "listStreamProcessors" {
		t.Errorf("expected method 'listStreamProcessors', got %q", stmt.MethodName)
	}
	if len(stmt.Args) != 0 {
		t.Errorf("expected 0 args, got %d", len(stmt.Args))
	}
}

func TestSpProcessorStats(t *testing.T) {
	node := mustParse(t, `sp.myProcessor.stats()`)
	stmt := node.(*ast.SpStatement)
	if stmt.MethodName != "myProcessor" {
		t.Errorf("expected method 'myProcessor', got %q", stmt.MethodName)
	}
	if stmt.SubMethod != "stats" {
		t.Errorf("expected sub-method 'stats', got %q", stmt.SubMethod)
	}
}

func TestSpProcessorStart(t *testing.T) {
	node := mustParse(t, `sp.analyticsProcessor.start()`)
	stmt := node.(*ast.SpStatement)
	if stmt.MethodName != "analyticsProcessor" {
		t.Errorf("expected method 'analyticsProcessor', got %q", stmt.MethodName)
	}
	if stmt.SubMethod != "start" {
		t.Errorf("expected sub-method 'start', got %q", stmt.SubMethod)
	}
}

func TestSpProcessorStop(t *testing.T) {
	node := mustParse(t, `sp.dataProcessor.stop()`)
	stmt := node.(*ast.SpStatement)
	if stmt.MethodName != "dataProcessor" {
		t.Errorf("expected method 'dataProcessor', got %q", stmt.MethodName)
	}
	if stmt.SubMethod != "stop" {
		t.Errorf("expected sub-method 'stop', got %q", stmt.SubMethod)
	}
}

func TestSpProcessorDrop(t *testing.T) {
	node := mustParse(t, `sp.oldProcessor.drop()`)
	stmt := node.(*ast.SpStatement)
	if stmt.MethodName != "oldProcessor" {
		t.Errorf("expected method 'oldProcessor', got %q", stmt.MethodName)
	}
	if stmt.SubMethod != "drop" {
		t.Errorf("expected sub-method 'drop', got %q", stmt.SubMethod)
	}
}

func TestSpProcessorSample(t *testing.T) {
	node := mustParse(t, `sp.sensorProcessor.sample()`)
	stmt := node.(*ast.SpStatement)
	if stmt.MethodName != "sensorProcessor" {
		t.Errorf("expected method 'sensorProcessor', got %q", stmt.MethodName)
	}
	if stmt.SubMethod != "sample" {
		t.Errorf("expected sub-method 'sample', got %q", stmt.SubMethod)
	}
}

func TestSpLoc(t *testing.T) {
	node := mustParse(t, `sp.process()`)
	assertLoc(t, node, 0, 12)
}
