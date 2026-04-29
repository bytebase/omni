package parser

import (
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

func TestParseCreateView(t *testing.T) {
	p := newTestParser("VIEW emp_view AS SELECT employee_id, last_name FROM employees")
	stmt, parseErr1 := p.parseCreateViewStmt(0, false)
	if parseErr1 != nil {
		t.Fatalf("parse: %v", parseErr1)
	}
	if stmt == nil {
		t.Fatal("expected CreateViewStmt, got nil")
	}
	if stmt.Name == nil || stmt.Name.Name != "EMP_VIEW" {
		t.Errorf("expected view name EMP_VIEW, got %v", stmt.Name)
	}
	if stmt.Query == nil {
		t.Fatal("expected non-nil Query")
	}
	sel, ok := stmt.Query.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("expected SelectStmt, got %T", stmt.Query)
	}
	if sel.TargetList.Len() != 2 {
		t.Errorf("expected 2 select targets, got %d", sel.TargetList.Len())
	}
}

func TestParseCreateOrReplaceView(t *testing.T) {
	p := newTestParser("VIEW emp_view AS SELECT 1 FROM dual")
	stmt, parseErr2 := p.parseCreateViewStmt(0, true)
	if parseErr2 != nil {
		t.Fatalf("parse: %v", parseErr2)
	}
	if !stmt.OrReplace {
		t.Error("expected OrReplace to be true")
	}
}

func TestParseCreateForceView(t *testing.T) {
	p := newTestParser("FORCE VIEW emp_view AS SELECT 1 FROM dual")
	stmt, parseErr3 := p.parseCreateViewStmt(0, false)
	if parseErr3 != nil {
		t.Fatalf("parse: %v", parseErr3)
	}
	if !stmt.Force {
		t.Error("expected Force to be true")
	}
}

func TestParseCreateViewWithColumns(t *testing.T) {
	p := newTestParser("VIEW emp_view (id, name) AS SELECT employee_id, last_name FROM employees")
	stmt, parseErr4 := p.parseCreateViewStmt(0, false)
	if parseErr4 != nil {
		t.Fatalf("parse: %v", parseErr4)
	}
	if stmt.Columns == nil || stmt.Columns.Len() != 2 {
		t.Fatalf("expected 2 column aliases, got %d", stmt.Columns.Len())
	}
	col0 := stmt.Columns.Items[0].(*ast.String)
	if col0.Str != "ID" {
		t.Errorf("expected column alias ID, got %q", col0.Str)
	}
	col1 := stmt.Columns.Items[1].(*ast.String)
	if col1.Str != "NAME" {
		t.Errorf("expected column alias NAME, got %q", col1.Str)
	}
}

func TestParseCreateViewWithCheckOption(t *testing.T) {
	p := newTestParser("VIEW emp_view AS SELECT 1 FROM dual WITH CHECK OPTION")
	stmt, parseErr5 := p.parseCreateViewStmt(0, false)
	if parseErr5 != nil {
		t.Fatalf("parse: %v", parseErr5)
	}
	if !stmt.WithCheckOpt {
		t.Error("expected WithCheckOpt to be true")
	}
}

func TestParseCreateViewWithReadOnly(t *testing.T) {
	p := newTestParser("VIEW emp_view AS SELECT 1 FROM dual WITH READ ONLY")
	stmt, parseErr6 := p.parseCreateViewStmt(0, false)
	if parseErr6 != nil {
		t.Fatalf("parse: %v", parseErr6)
	}
	if !stmt.WithReadOnly {
		t.Error("expected WithReadOnly to be true")
	}
}

func TestParseCreateMaterializedView(t *testing.T) {
	p := newTestParser("MATERIALIZED VIEW mv_emp AS SELECT employee_id FROM employees")
	stmt, parseErr7 := p.parseCreateViewStmt(0, false)
	if parseErr7 != nil {
		t.Fatalf("parse: %v", parseErr7)
	}
	if !stmt.Materialized {
		t.Error("expected Materialized to be true")
	}
	if stmt.Name == nil || stmt.Name.Name != "MV_EMP" {
		t.Errorf("expected view name MV_EMP, got %v", stmt.Name)
	}
	if stmt.Query == nil {
		t.Fatal("expected non-nil Query")
	}
}

func TestParseCreateMaterializedViewRefresh(t *testing.T) {
	p := newTestParser("MATERIALIZED VIEW mv_emp REFRESH FAST AS SELECT 1 FROM dual")
	stmt, parseErr8 := p.parseCreateViewStmt(0, false)
	if parseErr8 != nil {
		t.Fatalf("parse: %v", parseErr8)
	}
	if stmt.RefreshMethod != "FAST" {
		t.Errorf("expected RefreshMethod FAST, got %q", stmt.RefreshMethod)
	}
}

func TestParseCreateMaterializedViewBuildDeferred(t *testing.T) {
	p := newTestParser("MATERIALIZED VIEW mv_emp BUILD DEFERRED AS SELECT 1 FROM dual")
	stmt, parseErr9 := p.parseCreateViewStmt(0, false)
	if parseErr9 != nil {
		t.Fatalf("parse: %v", parseErr9)
	}
	if stmt.BuildMode != "DEFERRED" {
		t.Errorf("expected BuildMode DEFERRED, got %q", stmt.BuildMode)
	}
}

func TestParseCreateViewLoc(t *testing.T) {
	p := newTestParser("VIEW v AS SELECT 1 FROM dual")
	stmt, parseErr10 := p.parseCreateViewStmt(0, false)
	if parseErr10 != nil {
		t.Fatalf("parse: %v", parseErr10)
	}
	if stmt.Loc.Start != 0 {
		t.Errorf("expected Loc.Start=0, got %d", stmt.Loc.Start)
	}
	if stmt.Loc.End <= stmt.Loc.Start {
		t.Errorf("expected Loc.End > Loc.Start, got %d", stmt.Loc.End)
	}
}
