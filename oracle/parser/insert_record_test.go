package parser

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// TestInsertValuesRecord pins the PL/SQL record form of INSERT:
// INSERT INTO t VALUES record — engine-verified on Oracle 23ai Free,
// including FORALL collection elements (VALUES arr(i)) and RETURNING.
// The tuple form stays unchanged.
func TestInsertValuesRecord(t *testing.T) {
	accepts := []struct {
		name string
		sql  string
	}{
		{"bare record", "BEGIN INSERT INTO sms_card VALUES v_row_smscard; END;"},
		{"record with returning", "BEGIN INSERT INTO t VALUES v_rec RETURNING id INTO v_id; END;"},
		{"package-qualified record", "BEGIN INSERT INTO t VALUES pkg.g_rec; END;"},
		{"forall collection element", "BEGIN FORALL i IN 1 .. 2 INSERT INTO t9 VALUES arr(i); END;"},
		{"tuple form control", "BEGIN INSERT INTO t VALUES (1, 2); END;"},
	}
	for _, tc := range accepts {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := Parse(tc.sql); err != nil {
				t.Fatalf("parse failed: %v", err)
			}
		})
	}

	t.Run("record lands on ValuesRecord not Values", func(t *testing.T) {
		list, err := Parse("BEGIN INSERT INTO t VALUES v_rec; END;")
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		block := list.Items[0].(*nodes.RawStmt).Stmt.(*nodes.PLSQLBlock)
		ins, ok := block.Statements.Items[0].(*nodes.InsertStmt)
		if !ok {
			t.Fatalf("expected InsertStmt, got %T", block.Statements.Items[0])
		}
		if ins.ValuesRecord == nil {
			t.Fatal("ValuesRecord not set")
		}
		if ins.Values != nil {
			t.Fatal("Values must stay nil for the record form")
		}
		out := nodes.NodeToString(list)
		if !strings.Contains(out, ":valuesRecord") {
			t.Fatalf("expected :valuesRecord in output, got %s", out)
		}
	})

	t.Run("bare VALUES still rejected", func(t *testing.T) {
		if _, err := Parse("INSERT INTO t VALUES"); err == nil {
			t.Fatal("expected parse error for bare VALUES")
		}
	})
}
