package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/partiql/ast"
)

// TestParser_FromLedDML covers the FROM-led DML form (the 2nd alternative
// of dml#DmlBaseWrapper in PartiQLParser.g4):
//
//	FROM fromClause whereClause? dmlBaseCommand+ returningClause?
//
// Every accept-case below was confirmed ACCEPTED by the generated ANTLR
// PartiQL parser (Script() entry rule) via a differential probe; every
// reject-case was confirmed REJECTED (syntax error) by the same parser.
// The ANTLR grammar is the executable oracle (antlr_fallback) for this node.
func TestParser_FromLedDML(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantStr string // substring NodeToString must contain
	}{
		{
			name:    "from_set",
			input:   "FROM t SET a = 1",
			wantStr: "UpdateStmt{From:true Source:VarRef{Name:t",
		},
		{
			name:    "from_where_set",
			input:   "FROM t WHERE x = 1 SET a = 1",
			wantStr: "Where:",
		},
		{
			name:    "from_set_multi",
			input:   "FROM t SET a = 1, b = 2",
			wantStr: "Sets:[SetAssignment{Target:PathExpr{Root:VarRef{Name:a",
		},
		{
			name:    "from_set_remove",
			input:   "FROM t SET a = 1 REMOVE b",
			wantStr: "RemoveStmt{Path:",
		},
		{
			name:    "from_remove",
			input:   "FROM t REMOVE a.b",
			wantStr: "UpdateStmt{From:true Source:VarRef{Name:t",
		},
		{
			name:    "from_insert",
			input:   "FROM t INSERT INTO a VALUE 1",
			wantStr: "InsertStmt{",
		},
		{
			name:    "from_replace_bare_expr",
			input:   "FROM t REPLACE INTO a 1",
			wantStr: "ReplaceStmt{",
		},
		{
			name:    "from_upsert_bare_expr",
			input:   "FROM t UPSERT INTO a 1",
			wantStr: "UpsertStmt{",
		},
		{
			name:    "from_where_set_returning",
			input:   "FROM t WHERE x = 1 SET a = 1 RETURNING MODIFIED NEW *",
			wantStr: "Returning:ReturningClause{",
		},
		{
			name:    "from_set_returning",
			input:   "FROM t SET a = 1 RETURNING ALL OLD *",
			wantStr: "Returning:ReturningClause{",
		},
		{
			name:    "from_join_set",
			input:   "FROM a, b SET x = 1",
			wantStr: "JoinExpr{",
		},
		{
			name:    "from_path_set",
			input:   "FROM t.x.y SET a = 1",
			wantStr: "UpdateStmt{From:true",
		},
		{
			name:    "from_alias_set",
			input:   "FROM t AS x SET a = 1",
			wantStr: "AliasedSource{",
		},
		// FROM source is a full tableReference, so joins / paren / unpivot
		// all flow through (all ACCEPTed by the ANTLR oracle).
		{
			name:    "from_join_on_set",
			input:   "FROM a JOIN b ON a.id = b.id SET x = 1",
			wantStr: "JoinExpr{Kind:INNER",
		},
		{
			name:    "from_paren_set",
			input:   "FROM (t) SET a = 1",
			wantStr: "UpdateStmt{From:true Source:VarRef{Name:t",
		},
		{
			name:    "from_unpivot_set",
			input:   "FROM UNPIVOT x SET a = 1",
			wantStr: "Source:UnpivotExpr{",
		},
		// Multiple non-SET commands of the same kind, captured in order.
		{
			name:    "from_multi_insert",
			input:   "FROM t INSERT INTO a VALUE 1 INSERT INTO b VALUE 2",
			wantStr: "Commands:[InsertStmt{Target:PathExpr{Root:VarRef{Name:a} Steps:[]} Value:NumberLit{Val:1}} InsertStmt{Target:PathExpr{Root:VarRef{Name:b}",
		},
		{
			name:    "from_set_insert_remove",
			input:   "FROM t SET a = 1 INSERT INTO b VALUE 2 REMOVE c",
			wantStr: "Commands:[InsertStmt{Target:PathExpr{Root:VarRef{Name:b} Steps:[]} Value:NumberLit{Val:2}} RemoveStmt{Path:",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("ParseStatement(%q) unexpected error: %v", tc.input, err)
			}
			us, ok := stmt.(*ast.UpdateStmt)
			if !ok {
				t.Fatalf("ParseStatement(%q) = %T, want *ast.UpdateStmt", tc.input, stmt)
			}
			if !us.From {
				t.Errorf("FROM-led form should have From=true, got false")
			}
			got := ast.NodeToString(stmt)
			if !strings.Contains(got, tc.wantStr) {
				t.Errorf("NodeToString = %q, want to contain %q", got, tc.wantStr)
			}
		})
	}
}

// TestParser_FromLedDML_Script verifies FROM-led DML also dispatches through
// the script-level entry point (Parse), which uses parseRootStatement rather
// than ParseStatement.
func TestParser_FromLedDML_Script(t *testing.T) {
	list, err := Parse("FROM t SET a = 1; FROM u WHERE x = 1 REMOVE p.q")
	if err != nil {
		t.Fatalf("Parse unexpected error: %v", err)
	}
	if list.Len() != 2 {
		t.Fatalf("Parse len = %d, want 2", list.Len())
	}
	for i, item := range list.Items {
		if _, ok := item.(*ast.UpdateStmt); !ok {
			t.Errorf("item[%d] = %T, want *ast.UpdateStmt", i, item)
		}
	}
}

// TestParser_FromLedDML_Explain verifies the EXPLAIN prefix composes with
// FROM-led DML through the script entry (the EXPLAIN prefix is handled by
// parseRoot, then the inner statement dispatches to FROM-led DML). ANTLR
// accepts `EXPLAIN FROM t SET a = 1`.
func TestParser_FromLedDML_Explain(t *testing.T) {
	list, err := Parse("EXPLAIN FROM t SET a = 1")
	if err != nil {
		t.Fatalf("Parse unexpected error: %v", err)
	}
	if list.Len() != 1 {
		t.Fatalf("Parse len = %d, want 1", list.Len())
	}
	ex, ok := list.Items[0].(*ast.ExplainStmt)
	if !ok {
		t.Fatalf("item[0] = %T, want *ast.ExplainStmt", list.Items[0])
	}
	us, ok := ex.Inner.(*ast.UpdateStmt)
	if !ok {
		t.Fatalf("EXPLAIN inner = %T, want *ast.UpdateStmt", ex.Inner)
	}
	if !us.From {
		t.Errorf("inner UpdateStmt From = false, want true")
	}
}

// TestParser_UpdateParenSource_KnownGap documents a PRE-EXISTING limitation
// unrelated to this fix: the UPDATE source parser (parseUpdateSource) only
// accepts a bare/quoted identifier, not a full tableBaseReference, so a
// parenthesised source `UPDATE (t) SET ...` is rejected even though the
// ANTLR oracle accepts it. This is out of scope for the multi-command fix
// (B7/B8 concern the dmlBaseCommand+ loop, not the source production) and is
// flagged in the migration divergence ledger. The test pins the current
// behavior so a future source-parser node can flip it deliberately.
func TestParser_UpdateParenSource_KnownGap(t *testing.T) {
	p := NewParser("UPDATE (t) SET a = 1")
	if _, err := p.ParseStatement(); err == nil {
		t.Fatal("expected error for parenthesised UPDATE source (known gap), got nil")
	}
}

// TestParser_UpdateMultiCommand covers the UPDATE multi-command form (the
// 1st alternative of dml#DmlBaseWrapper):
//
//	UPDATE updateClause dmlBaseCommand+ whereClause? returningClause?
//
// All accept/reject verdicts confirmed against the generated ANTLR parser.
func TestParser_UpdateMultiCommand(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantStr string
	}{
		{
			name:    "update_set",
			input:   "UPDATE t SET a = 1",
			wantStr: "UpdateStmt{Source:VarRef{Name:t",
		},
		{
			name:    "update_remove",
			input:   "UPDATE t REMOVE a.b",
			wantStr: "Commands:[RemoveStmt{Path:",
		},
		{
			name:    "update_set_remove",
			input:   "UPDATE t SET a = 1 REMOVE b",
			wantStr: "RemoveStmt{Path:",
		},
		{
			name:    "update_insert",
			input:   "UPDATE t INSERT INTO a VALUE 1",
			wantStr: "Commands:[InsertStmt{",
		},
		{
			name:    "update_set_insert",
			input:   "UPDATE t SET a = 1 INSERT INTO a VALUE 2",
			wantStr: "InsertStmt{",
		},
		{
			name:    "update_replace_bare_expr",
			input:   "UPDATE t REPLACE INTO a 1",
			wantStr: "ReplaceStmt{",
		},
		{
			name:    "update_upsert_bare_expr",
			input:   "UPDATE t UPSERT INTO a 1",
			wantStr: "UpsertStmt{",
		},
		{
			name:    "update_remove_set",
			input:   "UPDATE t REMOVE a.b SET c = 1",
			wantStr: "Commands:[RemoveStmt{Path:",
		},
		{
			name:    "update_set_where",
			input:   "UPDATE t SET a = 1 WHERE x = 1",
			wantStr: "Where:",
		},
		{
			name:    "update_remove_where",
			input:   "UPDATE t REMOVE a.b WHERE x = 1",
			wantStr: "Where:",
		},
		{
			name:    "update_set_remove_where_returning",
			input:   "UPDATE t SET a = 1 REMOVE b WHERE x = 1 RETURNING MODIFIED NEW *",
			wantStr: "Returning:ReturningClause{",
		},
		{
			name:    "update_multi_set",
			input:   "UPDATE t SET a = 1 SET b = 2",
			wantStr: "Sets:[SetAssignment{Target:PathExpr{Root:VarRef{Name:a} Steps:[]} Value:NumberLit{Val:1}} SetAssignment{Target:PathExpr{Root:VarRef{Name:b}",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("ParseStatement(%q) unexpected error: %v", tc.input, err)
			}
			us, ok := stmt.(*ast.UpdateStmt)
			if !ok {
				t.Fatalf("ParseStatement(%q) = %T, want *ast.UpdateStmt", tc.input, stmt)
			}
			if us.From {
				t.Errorf("UPDATE form should have From=false, got true")
			}
			got := ast.NodeToString(stmt)
			if !strings.Contains(got, tc.wantStr) {
				t.Errorf("NodeToString = %q, want to contain %q", got, tc.wantStr)
			}
		})
	}
}

// TestParser_DMLWrapperRejects covers inputs the ANTLR oracle REJECTS for
// both the FROM-led and UPDATE multi-command forms. omni must reject them
// too (a syntax error). Without these negative cases the dispatch could be
// over-permissive.
func TestParser_DMLWrapperRejects(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantErrIn string
	}{
		// FROM-led: at least one dmlBaseCommand is required.
		{
			name:      "from_only",
			input:     "FROM t",
			wantErrIn: "expected SET or another DML command",
		},
		{
			name:      "from_where_only",
			input:     "FROM t WHERE x = 1",
			wantErrIn: "expected SET or another DML command",
		},
		// FROM-led: DELETE is not a dmlBaseCommand (analysis.md:375's
		// "FROM..DELETE" claim is inaccurate — ANTLR rejects it).
		{
			name:      "from_delete",
			input:     "FROM t DELETE",
			wantErrIn: "expected SET or another DML command",
		},
		// FROM-led: SELECT is not a dmlBaseCommand.
		{
			name:      "from_select",
			input:     "FROM t SELECT a",
			wantErrIn: "expected SET or another DML command",
		},
		// FROM-led: WHERE must precede the commands; a trailing WHERE is
		// rejected (the alt has no whereClause after the commands).
		{
			name:      "from_set_then_where",
			input:     "FROM t SET a = 1 WHERE x = 1",
			wantErrIn: "after statement",
		},
		// UPDATE: at least one dmlBaseCommand is required.
		{
			name:      "update_only",
			input:     "UPDATE t",
			wantErrIn: "expected SET or another DML command",
		},
		// UPDATE: WHERE only follows the commands; a leading WHERE is rejected.
		{
			name:      "update_where_before_set",
			input:     "UPDATE t WHERE x = 1 SET a = 1",
			wantErrIn: "expected SET or another DML command",
		},
		// UPDATE: bare WHERE with no command is rejected.
		{
			name:      "update_where_only",
			input:     "UPDATE t WHERE x = 1",
			wantErrIn: "expected SET or another DML command",
		},
		// UPDATE: DELETE / SELECT are not dmlBaseCommands.
		{
			name:      "update_delete",
			input:     "UPDATE t DELETE",
			wantErrIn: "expected SET or another DML command",
		},
		{
			name:      "update_select",
			input:     "UPDATE t SELECT a",
			wantErrIn: "expected SET or another DML command",
		},
		// REPLACE/UPSERT as a command does NOT take a VALUE keyword (the
		// grammar is `REPLACE INTO sym asIdent? expr`); VALUE there is a
		// syntax error. omni rejects it while parsing the replace value
		// expression (the ANTLR oracle reports "extraneous input 'VALUE'");
		// both are syntax rejections — the verdict, not the wording, is the
		// contract.
		{
			name:      "from_replace_value_kw",
			input:     "FROM t REPLACE INTO a VALUE 1",
			wantErrIn: "unexpected token \"VALUE\" in expression",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			_, err := p.ParseStatement()
			if err == nil {
				t.Fatalf("ParseStatement(%q): expected error, got nil", tc.input)
			}
			if !strings.Contains(err.Error(), tc.wantErrIn) {
				t.Errorf("error = %q, want to contain %q", err.Error(), tc.wantErrIn)
			}
		})
	}
}
