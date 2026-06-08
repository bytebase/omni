package parser

import (
	"testing"

	"github.com/bytebase/omni/redshift/plpgsql/ast"
)

// --------------------------------------------------------------------------
// 4.3 — CLOSE cursor and GET DIAGNOSTICS
// --------------------------------------------------------------------------

func TestCloseCursor(t *testing.T) {
	block := parseOK(t, `BEGIN CLOSE cur; END`)
	if len(block.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(block.Body))
	}
	n, ok := block.Body[0].(*ast.PLClose)
	if !ok {
		t.Fatalf("expected *ast.PLClose, got %T", block.Body[0])
	}
	if n.CursorVar != "cur" {
		t.Errorf("CursorVar = %q, want %q", n.CursorVar, "cur")
	}
}

func TestGetDiagCurrent(t *testing.T) {
	block := parseOK(t, `BEGIN GET DIAGNOSTICS rowcount = ROW_COUNT; END`)
	if len(block.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(block.Body))
	}
	n, ok := block.Body[0].(*ast.PLGetDiag)
	if !ok {
		t.Fatalf("expected *ast.PLGetDiag, got %T", block.Body[0])
	}
	if n.IsStacked {
		t.Error("IsStacked = true, want false")
	}
	if len(n.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(n.Items))
	}
	item := n.Items[0].(*ast.PLGetDiagItem)
	if item.Target != "rowcount" {
		t.Errorf("Target = %q, want %q", item.Target, "rowcount")
	}
	if item.Kind != "ROW_COUNT" {
		t.Errorf("Kind = %q, want %q", item.Kind, "ROW_COUNT")
	}
}

func TestGetCurrentDiagExplicit(t *testing.T) {
	block := parseOK(t, `BEGIN GET CURRENT DIAGNOSTICS rowcount = ROW_COUNT; END`)
	if len(block.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(block.Body))
	}
	n := block.Body[0].(*ast.PLGetDiag)
	if n.IsStacked {
		t.Error("IsStacked = true, want false")
	}
	if len(n.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(n.Items))
	}
	item := n.Items[0].(*ast.PLGetDiagItem)
	if item.Target != "rowcount" {
		t.Errorf("Target = %q, want %q", item.Target, "rowcount")
	}
	if item.Kind != "ROW_COUNT" {
		t.Errorf("Kind = %q, want %q", item.Kind, "ROW_COUNT")
	}
}

func TestGetDiagMultipleItems(t *testing.T) {
	block := parseOK(t, `BEGIN GET DIAGNOSTICS rc = ROW_COUNT, ctx = PG_CONTEXT; END`)
	n := block.Body[0].(*ast.PLGetDiag)
	if n.IsStacked {
		t.Error("IsStacked = true, want false")
	}
	if len(n.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(n.Items))
	}
	item0 := n.Items[0].(*ast.PLGetDiagItem)
	if item0.Target != "rc" || item0.Kind != "ROW_COUNT" {
		t.Errorf("item[0] = (%q, %q), want (rc, ROW_COUNT)", item0.Target, item0.Kind)
	}
	item1 := n.Items[1].(*ast.PLGetDiagItem)
	if item1.Target != "ctx" || item1.Kind != "PG_CONTEXT" {
		t.Errorf("item[1] = (%q, %q), want (ctx, PG_CONTEXT)", item1.Target, item1.Kind)
	}
}

func TestGetStackedDiag(t *testing.T) {
	block := parseOK(t, `BEGIN GET STACKED DIAGNOSTICS msg = MESSAGE_TEXT; END`)
	n := block.Body[0].(*ast.PLGetDiag)
	if !n.IsStacked {
		t.Error("IsStacked = false, want true")
	}
	if len(n.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(n.Items))
	}
	item := n.Items[0].(*ast.PLGetDiagItem)
	if item.Target != "msg" || item.Kind != "MESSAGE_TEXT" {
		t.Errorf("item = (%q, %q), want (msg, MESSAGE_TEXT)", item.Target, item.Kind)
	}
}

func TestGetStackedDiagMultiple(t *testing.T) {
	block := parseOK(t, `BEGIN GET STACKED DIAGNOSTICS st = RETURNED_SQLSTATE, msg = MESSAGE_TEXT, det = PG_EXCEPTION_DETAIL; END`)
	n := block.Body[0].(*ast.PLGetDiag)
	if !n.IsStacked {
		t.Error("IsStacked = false, want true")
	}
	if len(n.Items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(n.Items))
	}
	want := []struct{ target, kind string }{
		{"st", "RETURNED_SQLSTATE"},
		{"msg", "MESSAGE_TEXT"},
		{"det", "PG_EXCEPTION_DETAIL"},
	}
	for i, w := range want {
		item := n.Items[i].(*ast.PLGetDiagItem)
		if item.Target != w.target || item.Kind != w.kind {
			t.Errorf("item[%d] = (%q, %q), want (%q, %q)", i, item.Target, item.Kind, w.target, w.kind)
		}
	}
}

func TestGetDiagColonEquals(t *testing.T) {
	block := parseOK(t, `BEGIN GET DIAGNOSTICS rc := ROW_COUNT; END`)
	n := block.Body[0].(*ast.PLGetDiag)
	if len(n.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(n.Items))
	}
	item := n.Items[0].(*ast.PLGetDiagItem)
	if item.Target != "rc" || item.Kind != "ROW_COUNT" {
		t.Errorf("item = (%q, %q), want (rc, ROW_COUNT)", item.Target, item.Kind)
	}
}

func TestGetDiagAllCurrentItems(t *testing.T) {
	block := parseOK(t, `BEGIN GET DIAGNOSTICS a = ROW_COUNT, b = PG_CONTEXT, c = PG_ROUTINE_OID; END`)
	n := block.Body[0].(*ast.PLGetDiag)
	if len(n.Items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(n.Items))
	}
	want := []string{"ROW_COUNT", "PG_CONTEXT", "PG_ROUTINE_OID"}
	for i, w := range want {
		item := n.Items[i].(*ast.PLGetDiagItem)
		if item.Kind != w {
			t.Errorf("item[%d].Kind = %q, want %q", i, item.Kind, w)
		}
	}
}

func TestGetDiagAllStackedItems(t *testing.T) {
	body := `BEGIN GET STACKED DIAGNOSTICS
		a = RETURNED_SQLSTATE,
		b = MESSAGE_TEXT,
		c = PG_EXCEPTION_DETAIL,
		d = PG_EXCEPTION_HINT,
		e = PG_EXCEPTION_CONTEXT,
		f = COLUMN_NAME,
		g = CONSTRAINT_NAME,
		h = PG_DATATYPE_NAME,
		i = TABLE_NAME,
		j = SCHEMA_NAME;
	END`
	block := parseOK(t, body)
	n := block.Body[0].(*ast.PLGetDiag)
	if !n.IsStacked {
		t.Error("IsStacked = false, want true")
	}
	want := []string{
		"RETURNED_SQLSTATE", "MESSAGE_TEXT",
		"PG_EXCEPTION_DETAIL", "PG_EXCEPTION_HINT", "PG_EXCEPTION_CONTEXT",
		"COLUMN_NAME", "CONSTRAINT_NAME", "PG_DATATYPE_NAME",
		"TABLE_NAME", "SCHEMA_NAME",
	}
	if len(n.Items) != len(want) {
		t.Fatalf("expected %d items, got %d", len(want), len(n.Items))
	}
	for i, w := range want {
		item := n.Items[i].(*ast.PLGetDiagItem)
		if item.Kind != w {
			t.Errorf("item[%d].Kind = %q, want %q", i, item.Kind, w)
		}
	}
}
