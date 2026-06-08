package parser

import (
	"testing"

	"github.com/bytebase/omni/redshift/plpgsql/ast"
)

// --------------------------------------------------------------------------
// Section 4.1: OPEN Cursor
// --------------------------------------------------------------------------

func TestOpenBoundCursor(t *testing.T) {
	block := parseOK(t, `BEGIN OPEN cur; END`)
	if len(block.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(block.Body))
	}
	node, ok := block.Body[0].(*ast.PLOpen)
	if !ok {
		t.Fatalf("expected PLOpen, got %T", block.Body[0])
	}
	if node.CursorVar != "cur" {
		t.Errorf("CursorVar = %q, want %q", node.CursorVar, "cur")
	}
	if node.Query != "" {
		t.Errorf("Query = %q, want empty", node.Query)
	}
	if node.DynQuery != "" {
		t.Errorf("DynQuery = %q, want empty", node.DynQuery)
	}
	if node.Args != "" {
		t.Errorf("Args = %q, want empty", node.Args)
	}
}

func TestOpenBoundCursorWithArgs(t *testing.T) {
	block := parseOK(t, `BEGIN OPEN cur(1, 'a'); END`)
	if len(block.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(block.Body))
	}
	node, ok := block.Body[0].(*ast.PLOpen)
	if !ok {
		t.Fatalf("expected PLOpen, got %T", block.Body[0])
	}
	if node.CursorVar != "cur" {
		t.Errorf("CursorVar = %q, want %q", node.CursorVar, "cur")
	}
	if node.Args != "1, 'a'" {
		t.Errorf("Args = %q, want %q", node.Args, "1, 'a'")
	}
}

func TestOpenUnboundCursorForQuery(t *testing.T) {
	block := parseOK(t, `BEGIN OPEN cur FOR SELECT * FROM t; END`)
	if len(block.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(block.Body))
	}
	node, ok := block.Body[0].(*ast.PLOpen)
	if !ok {
		t.Fatalf("expected PLOpen, got %T", block.Body[0])
	}
	if node.CursorVar != "cur" {
		t.Errorf("CursorVar = %q, want %q", node.CursorVar, "cur")
	}
	if node.Query != "SELECT * FROM t" {
		t.Errorf("Query = %q, want %q", node.Query, "SELECT * FROM t")
	}
	if node.Scroll != ast.ScrollNone {
		t.Errorf("Scroll = %d, want ScrollNone", node.Scroll)
	}
}

func TestOpenUnboundCursorWithScroll(t *testing.T) {
	block := parseOK(t, `BEGIN OPEN cur SCROLL FOR SELECT * FROM t; END`)
	if len(block.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(block.Body))
	}
	node, ok := block.Body[0].(*ast.PLOpen)
	if !ok {
		t.Fatalf("expected PLOpen, got %T", block.Body[0])
	}
	if node.CursorVar != "cur" {
		t.Errorf("CursorVar = %q, want %q", node.CursorVar, "cur")
	}
	if node.Query != "SELECT * FROM t" {
		t.Errorf("Query = %q, want %q", node.Query, "SELECT * FROM t")
	}
	if node.Scroll != ast.ScrollYes {
		t.Errorf("Scroll = %d, want ScrollYes", node.Scroll)
	}
}

func TestOpenUnboundCursorWithNoScroll(t *testing.T) {
	block := parseOK(t, `BEGIN OPEN cur NO SCROLL FOR SELECT * FROM t; END`)
	if len(block.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(block.Body))
	}
	node, ok := block.Body[0].(*ast.PLOpen)
	if !ok {
		t.Fatalf("expected PLOpen, got %T", block.Body[0])
	}
	if node.Scroll != ast.ScrollNo {
		t.Errorf("Scroll = %d, want ScrollNo", node.Scroll)
	}
	if node.Query != "SELECT * FROM t" {
		t.Errorf("Query = %q, want %q", node.Query, "SELECT * FROM t")
	}
}

func TestOpenForExecute(t *testing.T) {
	block := parseOK(t, `BEGIN OPEN cur FOR EXECUTE 'SELECT * FROM t'; END`)
	if len(block.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(block.Body))
	}
	node, ok := block.Body[0].(*ast.PLOpen)
	if !ok {
		t.Fatalf("expected PLOpen, got %T", block.Body[0])
	}
	if node.CursorVar != "cur" {
		t.Errorf("CursorVar = %q, want %q", node.CursorVar, "cur")
	}
	if node.DynQuery != "'SELECT * FROM t'" {
		t.Errorf("DynQuery = %q, want %q", node.DynQuery, "'SELECT * FROM t'")
	}
}

func TestOpenForExecuteWithUsing(t *testing.T) {
	block := parseOK(t, `BEGIN OPEN cur FOR EXECUTE 'SELECT * FROM t WHERE id=$1' USING my_id; END`)
	if len(block.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(block.Body))
	}
	node, ok := block.Body[0].(*ast.PLOpen)
	if !ok {
		t.Fatalf("expected PLOpen, got %T", block.Body[0])
	}
	if node.DynQuery != "'SELECT * FROM t WHERE id=$1'" {
		t.Errorf("DynQuery = %q, want %q", node.DynQuery, "'SELECT * FROM t WHERE id=$1'")
	}
	if len(node.Params) != 1 || node.Params[0] != "my_id" {
		t.Errorf("Params = %v, want [my_id]", node.Params)
	}
}

// --------------------------------------------------------------------------
// Section 4.2: FETCH and MOVE
// --------------------------------------------------------------------------

func TestFetchDefault(t *testing.T) {
	block := parseOK(t, `BEGIN FETCH cur INTO x; END`)
	if len(block.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(block.Body))
	}
	node, ok := block.Body[0].(*ast.PLFetch)
	if !ok {
		t.Fatalf("expected PLFetch, got %T", block.Body[0])
	}
	if node.CursorVar != "cur" {
		t.Errorf("CursorVar = %q, want %q", node.CursorVar, "cur")
	}
	if node.Direction != ast.FetchNext {
		t.Errorf("Direction = %d, want FetchNext", node.Direction)
	}
	if len(node.Into) != 1 || node.Into[0] != "x" {
		t.Errorf("Into = %v, want [x]", node.Into)
	}
	if node.IsMove {
		t.Errorf("IsMove = true, want false")
	}
}

func TestFetchNext(t *testing.T) {
	block := parseOK(t, `BEGIN FETCH NEXT FROM cur INTO x; END`)
	node := block.Body[0].(*ast.PLFetch)
	if node.Direction != ast.FetchNext {
		t.Errorf("Direction = %d, want FetchNext", node.Direction)
	}
	if node.CursorVar != "cur" {
		t.Errorf("CursorVar = %q, want %q", node.CursorVar, "cur")
	}
	if len(node.Into) != 1 || node.Into[0] != "x" {
		t.Errorf("Into = %v, want [x]", node.Into)
	}
}

func TestFetchPrior(t *testing.T) {
	block := parseOK(t, `BEGIN FETCH PRIOR FROM cur INTO x; END`)
	node := block.Body[0].(*ast.PLFetch)
	if node.Direction != ast.FetchPrior {
		t.Errorf("Direction = %d, want FetchPrior", node.Direction)
	}
	if node.CursorVar != "cur" {
		t.Errorf("CursorVar = %q, want %q", node.CursorVar, "cur")
	}
}

func TestFetchFirst(t *testing.T) {
	block := parseOK(t, `BEGIN FETCH FIRST FROM cur INTO x; END`)
	node := block.Body[0].(*ast.PLFetch)
	if node.Direction != ast.FetchFirst {
		t.Errorf("Direction = %d, want FetchFirst", node.Direction)
	}
}

func TestFetchLast(t *testing.T) {
	block := parseOK(t, `BEGIN FETCH LAST FROM cur INTO x; END`)
	node := block.Body[0].(*ast.PLFetch)
	if node.Direction != ast.FetchLast {
		t.Errorf("Direction = %d, want FetchLast", node.Direction)
	}
}

func TestFetchAbsolute(t *testing.T) {
	block := parseOK(t, `BEGIN FETCH ABSOLUTE 5 FROM cur INTO x; END`)
	node := block.Body[0].(*ast.PLFetch)
	if node.Direction != ast.FetchAbsolute {
		t.Errorf("Direction = %d, want FetchAbsolute", node.Direction)
	}
	if node.Count != "5" {
		t.Errorf("Count = %q, want %q", node.Count, "5")
	}
}

func TestFetchRelative(t *testing.T) {
	block := parseOK(t, `BEGIN FETCH RELATIVE -1 FROM cur INTO x; END`)
	node := block.Body[0].(*ast.PLFetch)
	if node.Direction != ast.FetchRelative {
		t.Errorf("Direction = %d, want FetchRelative", node.Direction)
	}
	if node.Count != "-1" {
		t.Errorf("Count = %q, want %q", node.Count, "-1")
	}
}

func TestFetchForward(t *testing.T) {
	block := parseOK(t, `BEGIN FETCH FORWARD FROM cur INTO x; END`)
	node := block.Body[0].(*ast.PLFetch)
	if node.Direction != ast.FetchForward {
		t.Errorf("Direction = %d, want FetchForward", node.Direction)
	}
}

func TestFetchForwardN(t *testing.T) {
	block := parseOK(t, `BEGIN FETCH FORWARD 5 FROM cur INTO x; END`)
	node := block.Body[0].(*ast.PLFetch)
	if node.Direction != ast.FetchForwardN {
		t.Errorf("Direction = %d, want FetchForwardN", node.Direction)
	}
	if node.Count != "5" {
		t.Errorf("Count = %q, want %q", node.Count, "5")
	}
}

func TestFetchForwardAll(t *testing.T) {
	block := parseOK(t, `BEGIN FETCH FORWARD ALL FROM cur INTO x; END`)
	node := block.Body[0].(*ast.PLFetch)
	if node.Direction != ast.FetchForwardAll {
		t.Errorf("Direction = %d, want FetchForwardAll", node.Direction)
	}
}

func TestFetchBackward(t *testing.T) {
	block := parseOK(t, `BEGIN FETCH BACKWARD FROM cur INTO x; END`)
	node := block.Body[0].(*ast.PLFetch)
	if node.Direction != ast.FetchBackward {
		t.Errorf("Direction = %d, want FetchBackward", node.Direction)
	}
}

func TestFetchBackwardN(t *testing.T) {
	block := parseOK(t, `BEGIN FETCH BACKWARD 3 FROM cur INTO x; END`)
	node := block.Body[0].(*ast.PLFetch)
	if node.Direction != ast.FetchBackwardN {
		t.Errorf("Direction = %d, want FetchBackwardN", node.Direction)
	}
	if node.Count != "3" {
		t.Errorf("Count = %q, want %q", node.Count, "3")
	}
}

func TestFetchBackwardAll(t *testing.T) {
	block := parseOK(t, `BEGIN FETCH BACKWARD ALL FROM cur INTO x; END`)
	node := block.Body[0].(*ast.PLFetch)
	if node.Direction != ast.FetchBackwardAll {
		t.Errorf("Direction = %d, want FetchBackwardAll", node.Direction)
	}
}

func TestFetchWithIN(t *testing.T) {
	block := parseOK(t, `BEGIN FETCH NEXT IN cur INTO x; END`)
	node := block.Body[0].(*ast.PLFetch)
	if node.Direction != ast.FetchNext {
		t.Errorf("Direction = %d, want FetchNext", node.Direction)
	}
	if node.CursorVar != "cur" {
		t.Errorf("CursorVar = %q, want %q", node.CursorVar, "cur")
	}
	if len(node.Into) != 1 || node.Into[0] != "x" {
		t.Errorf("Into = %v, want [x]", node.Into)
	}
}

func TestMoveNext(t *testing.T) {
	block := parseOK(t, `BEGIN MOVE NEXT FROM cur; END`)
	node := block.Body[0].(*ast.PLFetch)
	if !node.IsMove {
		t.Errorf("IsMove = false, want true")
	}
	if node.Direction != ast.FetchNext {
		t.Errorf("Direction = %d, want FetchNext", node.Direction)
	}
	if node.CursorVar != "cur" {
		t.Errorf("CursorVar = %q, want %q", node.CursorVar, "cur")
	}
	if len(node.Into) != 0 {
		t.Errorf("Into = %v, want empty", node.Into)
	}
}

func TestMoveForwardAll(t *testing.T) {
	block := parseOK(t, `BEGIN MOVE FORWARD ALL FROM cur; END`)
	node := block.Body[0].(*ast.PLFetch)
	if !node.IsMove {
		t.Errorf("IsMove = false, want true")
	}
	if node.Direction != ast.FetchForwardAll {
		t.Errorf("Direction = %d, want FetchForwardAll", node.Direction)
	}
}

func TestMoveAbsolute(t *testing.T) {
	block := parseOK(t, `BEGIN MOVE ABSOLUTE 0 FROM cur; END`)
	node := block.Body[0].(*ast.PLFetch)
	if !node.IsMove {
		t.Errorf("IsMove = false, want true")
	}
	if node.Direction != ast.FetchAbsolute {
		t.Errorf("Direction = %d, want FetchAbsolute", node.Direction)
	}
	if node.Count != "0" {
		t.Errorf("Count = %q, want %q", node.Count, "0")
	}
}
