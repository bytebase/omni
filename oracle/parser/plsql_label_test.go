package parser

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// TestPLSQLStatementLabels pins <<label>> support on every statement kind:
// one or more labels may precede ANY PL/SQL statement, not only blocks
// (engine-verified on Oracle 23ai Free, including multiple labels and GOTO
// targets). Blocks and loops keep their first label on their own Label
// field for END-label matching; other statements carry labels on a
// PLSQLLabeledStatement wrapper.
func TestPLSQLStatementLabels(t *testing.T) {
	accepts := []struct {
		name string
		sql  string
	}{
		{"label + assignment", "BEGIN <<lbl>> x := 1; END;"},
		{"label + record field assignment", "BEGIN <<end_check>> v_row.input_note := a || '|' || b; END;"},
		{"label + null", "BEGIN <<lbl>> NULL; END;"},
		{"label + procedure call", "BEGIN <<lbl>> dbms_output.put_line('x'); END;"},
		{"label + if", "BEGIN <<lbl>> IF a THEN NULL; END IF; END;"},
		{"label + goto target flow", "BEGIN GOTO done; <<done>> NULL; END;"},
		{"two labels one statement", "BEGIN <<a>> <<b>> x := 1; END;"},
		{"label + block", "BEGIN <<outer>> BEGIN NULL; END; END;"},
		{"label + basic loop with end label", "BEGIN <<lp>> LOOP EXIT lp; END LOOP lp; END;"},
		{"label + for loop", "BEGIN <<lp>> FOR i IN 1 .. 3 LOOP NULL; END LOOP lp; END;"},
		{"label + while loop", "BEGIN <<lp>> WHILE a LOOP NULL; END LOOP lp; END;"},
		{"customer shape: label then labeled assignment in package body",
			"CREATE OR REPLACE PACKAGE BODY p IS PROCEDURE q IS BEGIN <<end_check>> v_row.input_note := a; END; END;"},
	}
	for _, tc := range accepts {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := Parse(tc.sql); err != nil {
				t.Fatalf("parse failed: %v", err)
			}
		})
	}

	// Documented over-accepts (parser-unsure -> accept direction, per the
	// failure-direction contract; the engine rejects these):
	//   - a label inside a DECLARE section
	//   - INSERT INTO t (a, b) VALUES rec (column list with record form)
	//   - the record form in plain SQL outside PL/SQL
	// Recorded by cross-validation; revisit only with engine-anchored
	// tightening.

	t.Run("unclosed label rejected", func(t *testing.T) {
		if _, err := Parse("BEGIN <<lbl x := 1; END;"); err == nil {
			t.Fatal("expected parse error for unclosed label")
		}
	})

	t.Run("empty label rejected", func(t *testing.T) {
		// Engine: PLS-00103 at <<>>.
		if _, err := Parse("BEGIN <<>> NULL; END;"); err == nil {
			t.Fatal("expected parse error for empty label")
		}
	})

	t.Run("quoted label with space accepted and quoted in output", func(t *testing.T) {
		// Engine-verified: <<"a b">> is legal PL/SQL.
		list, err := Parse(`BEGIN <<"a b">> NULL; END;`)
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		out := nodes.NodeToString(list)
		if !strings.Contains(out, `"a b"`) {
			t.Fatalf("expected quoted label in output, got %s", out)
		}
	})

	t.Run("mismatched end loop label stays accepted like the engine", func(t *testing.T) {
		// Engine-verified: Oracle ACCEPTS <<a>> LOOP ... END LOOP b; —
		// review-flagged as a candidate reject, refuted by the engine.
		if _, err := Parse("BEGIN <<a>> LOOP EXIT; END LOOP b; END;"); err != nil {
			t.Fatalf("parse failed: %v", err)
		}
	})

	t.Run("labeled block range covers its label", func(t *testing.T) {
		sql := "BEGIN <<x>> BEGIN NULL; END; END;"
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		outer := list.Items[0].(*nodes.RawStmt).Stmt.(*nodes.PLSQLBlock)
		inner := outer.Statements.Items[0].(*nodes.PLSQLBlock)
		if got := sql[inner.Loc.Start:]; !strings.HasPrefix(got, "<<x>>") {
			t.Fatalf("inner block starts at %q, want at the label", got[:10])
		}
		if nodes.NodeLoc(&nodes.PLSQLLabeledStatement{Loc: nodes.Loc{Start: 1, End: 2}}).Start != 1 {
			t.Fatal("NodeLoc missing PLSQLLabeledStatement case")
		}
	})

	t.Run("loop keeps first label on its Label field", func(t *testing.T) {
		list, err := Parse("BEGIN <<lp>> LOOP EXIT lp; END LOOP lp; END;")
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		block := list.Items[0].(*nodes.RawStmt).Stmt.(*nodes.PLSQLBlock)
		loop, ok := block.Statements.Items[0].(*nodes.PLSQLLoop)
		if !ok {
			t.Fatalf("expected PLSQLLoop, got %T", block.Statements.Items[0])
		}
		if loop.Label != "LP" && loop.Label != "lp" {
			t.Fatalf("loop label = %q, want lp", loop.Label)
		}
	})

	t.Run("plain statement labels ride the wrapper", func(t *testing.T) {
		list, err := Parse("BEGIN <<a>> <<b>> x := 1; END;")
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		block := list.Items[0].(*nodes.RawStmt).Stmt.(*nodes.PLSQLBlock)
		wrapped, ok := block.Statements.Items[0].(*nodes.PLSQLLabeledStatement)
		if !ok {
			t.Fatalf("expected PLSQLLabeledStatement, got %T", block.Statements.Items[0])
		}
		if len(wrapped.Labels) != 2 {
			t.Fatalf("labels = %v, want 2", wrapped.Labels)
		}
		out := nodes.NodeToString(list)
		if !strings.Contains(out, ":labels") {
			t.Fatalf("expected :labels in output, got %s", out)
		}
	})
}
