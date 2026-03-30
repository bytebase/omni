package parsertest

import (
	"testing"

	"github.com/bytebase/omni/mongo"
	"github.com/bytebase/omni/mongo/ast"
)

func TestPositionTracking(t *testing.T) {
	input := `db.users.find({name: "alice"})`
	stmts, err := mongo.Parse(input)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}

	s := stmts[0]
	if s.ByteStart != 0 {
		t.Errorf("expected ByteStart=0, got %d", s.ByteStart)
	}
	if s.ByteEnd != len(input) {
		t.Errorf("expected ByteEnd=%d, got %d", len(input), s.ByteEnd)
	}
	if s.Start.Line != 1 {
		t.Errorf("expected Start.Line=1, got %d", s.Start.Line)
	}
	if s.Start.Column != 1 {
		t.Errorf("expected Start.Column=1, got %d", s.Start.Column)
	}
	if s.End.Line != 1 {
		t.Errorf("expected End.Line=1, got %d", s.End.Line)
	}
}

func TestPositionMultipleStatements(t *testing.T) {
	input := "show dbs\nshow collections"
	stmts, err := mongo.Parse(input)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(stmts) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(stmts))
	}

	// First statement: "show dbs" on line 1
	s1 := stmts[0]
	if s1.Start.Line != 1 {
		t.Errorf("stmt1: expected Start.Line=1, got %d", s1.Start.Line)
	}
	if s1.Start.Column != 1 {
		t.Errorf("stmt1: expected Start.Column=1, got %d", s1.Start.Column)
	}
	if s1.ByteStart != 0 {
		t.Errorf("stmt1: expected ByteStart=0, got %d", s1.ByteStart)
	}
	if s1.ByteEnd != 8 { // "show dbs" = 8 bytes
		t.Errorf("stmt1: expected ByteEnd=8, got %d", s1.ByteEnd)
	}

	// Second statement: "show collections" on line 2
	s2 := stmts[1]
	if s2.Start.Line != 2 {
		t.Errorf("stmt2: expected Start.Line=2, got %d", s2.Start.Line)
	}
	if s2.Start.Column != 1 {
		t.Errorf("stmt2: expected Start.Column=1, got %d", s2.Start.Column)
	}
	if s2.ByteStart != 9 { // after "show dbs\n"
		t.Errorf("stmt2: expected ByteStart=9, got %d", s2.ByteStart)
	}
}

func TestPositionWithSemicolons(t *testing.T) {
	input := "show dbs; show collections"
	stmts, err := mongo.Parse(input)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(stmts) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(stmts))
	}

	// First statement should include the trailing semicolon in text
	s1 := stmts[0]
	if s1.Text != "show dbs;" {
		t.Errorf("stmt1: expected text %q, got %q", "show dbs;", s1.Text)
	}

	// Second statement
	s2 := stmts[1]
	if s2.Text != "show collections" {
		t.Errorf("stmt2: expected text %q, got %q", "show collections", s2.Text)
	}
}

func TestPositionDocumentNodes(t *testing.T) {
	input := `db.c.find({k: "v"})`
	stmts, err := mongo.Parse(input)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}

	s := stmts[0]
	cs, ok := s.AST.(*ast.CollectionStatement)
	if !ok {
		t.Fatalf("expected *ast.CollectionStatement, got %T", s.AST)
	}
	if len(cs.Args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(cs.Args))
	}

	doc, ok := cs.Args[0].(*ast.Document)
	if !ok {
		t.Fatalf("expected *ast.Document, got %T", cs.Args[0])
	}

	// Document starts at '{' which is at index 10: db.c.find(
	expectedStart := 10
	if doc.Loc.Start != expectedStart {
		t.Errorf("expected doc start=%d, got %d", expectedStart, doc.Loc.Start)
	}

	// {k: "v"} — '}' is at index 17, End is exclusive = 18
	expectedEnd := 18
	if doc.Loc.End != expectedEnd {
		t.Errorf("expected doc end=%d, got %d", expectedEnd, doc.Loc.End)
	}
}

func TestPositionArrayNodes(t *testing.T) {
	input := `db.c.insert([1, 2, 3])`
	stmts, err := mongo.Parse(input)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}

	cs := stmts[0].AST.(*ast.CollectionStatement)
	if len(cs.Args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(cs.Args))
	}

	arr, ok := cs.Args[0].(*ast.Array)
	if !ok {
		t.Fatalf("expected *ast.Array, got %T", cs.Args[0])
	}

	// Array starts at '[' at index 12
	if arr.Loc.Start != 12 {
		t.Errorf("expected array start=12, got %d", arr.Loc.Start)
	}
	// Array ends after ']' at index 21, so End=21
	if arr.Loc.End != 21 {
		t.Errorf("expected array end=21, got %d", arr.Loc.End)
	}
}

func TestPositionStatementLoc(t *testing.T) {
	// Verify the AST node's own Loc matches what we expect
	input := `db.users.find()`
	node := mustParse(t, input)
	loc := node.GetLoc()
	if loc.Start != 0 {
		t.Errorf("expected loc.Start=0, got %d", loc.Start)
	}
	if loc.End != len(input) {
		t.Errorf("expected loc.End=%d, got %d", len(input), loc.End)
	}
}
