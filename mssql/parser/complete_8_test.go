package parser

import (
	"testing"
)

// --- Section 8.1: Expression Positions ---

func TestCollect_ExpressionPositions(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantRules []string
		wantToks  []int
	}{
		{
			name:      "after = operator",
			sql:       "SELECT * FROM t WHERE a = ",
			wantRules: []string{"columnref", "func_name"},
		},
		{
			name:      "after AND",
			sql:       "SELECT * FROM t WHERE a = 1 AND ",
			wantRules: []string{"columnref", "func_name"},
		},
		{
			name:      "after OR",
			sql:       "SELECT * FROM t WHERE a = 1 OR ",
			wantRules: []string{"columnref", "func_name"},
		},
		{
			name:      "after ( in expression",
			sql:       "SELECT * FROM t WHERE (",
			wantRules: []string{"columnref", "func_name"},
			wantToks:  []int{kwSELECT},
		},
		{
			name:      "CASE WHEN |",
			sql:       "SELECT CASE WHEN ",
			wantRules: []string{"columnref", "func_name"},
		},
		{
			name:      "CASE WHEN x THEN |",
			sql:       "SELECT CASE WHEN 1=1 THEN ",
			wantRules: []string{"columnref", "func_name"},
		},
		{
			name:      "CASE WHEN x THEN y ELSE |",
			sql:       "SELECT CASE WHEN 1=1 THEN 1 ELSE ",
			wantRules: []string{"columnref", "func_name"},
		},
		{
			name:      "IN (|)",
			sql:       "SELECT * FROM t WHERE a IN (",
			wantRules: []string{"columnref", "func_name"},
			wantToks:  []int{kwSELECT},
		},
		{
			name:      "BETWEEN |",
			sql:       "SELECT * FROM t WHERE a BETWEEN ",
			wantRules: []string{"columnref", "func_name"},
		},
		{
			name:      "LIKE |",
			sql:       "SELECT * FROM t WHERE a LIKE ",
			wantRules: []string{"columnref", "func_name"},
		},
		{
			name:      "CAST(|",
			sql:       "SELECT CAST(",
			wantRules: []string{"columnref", "func_name"},
		},
		{
			name:      "CAST(x AS |)",
			sql:       "SELECT CAST(x AS ",
			wantRules: []string{"type_name"},
		},
		{
			name:      "CONVERT(|",
			sql:       "SELECT CONVERT(",
			wantRules: []string{"type_name"},
		},
		{
			name:      "COALESCE(|",
			sql:       "SELECT COALESCE(",
			wantRules: []string{"columnref", "func_name"},
		},
		{
			name:      "NULLIF(|",
			sql:       "SELECT NULLIF(",
			wantRules: []string{"columnref", "func_name"},
		},
		{
			name:      "IIF(|",
			sql:       "SELECT IIF(",
			wantRules: []string{"columnref", "func_name"},
		},
		{
			name:     "WHERE a IS |",
			sql:      "SELECT * FROM t WHERE a IS ",
			wantToks: []int{kwNULL, kwNOT},
		},
		{
			name:     "WHERE EXISTS (|)",
			sql:      "SELECT * FROM t WHERE EXISTS (",
			wantToks: []int{kwSELECT},
		},
		{
			name:      "BETWEEN 1 AND |",
			sql:       "SELECT * FROM t WHERE a BETWEEN 1 AND ",
			wantRules: []string{"columnref", "func_name"},
		},
		{
			name:      "CASE a WHEN |",
			sql:       "SELECT CASE a WHEN ",
			wantRules: []string{"columnref", "func_name"},
		},
		{
			name:      "TRY_CAST(x AS |)",
			sql:       "SELECT TRY_CAST(x AS ",
			wantRules: []string{"type_name"},
		},
		{
			name:      "TRY_CONVERT(|",
			sql:       "SELECT TRY_CONVERT(",
			wantRules: []string{"type_name"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := Collect(tt.sql, len(tt.sql))
			if cs == nil {
				t.Fatal("Collect returned nil")
			}
			for _, r := range tt.wantRules {
				if !cs.HasRule(r) {
					t.Errorf("missing rule candidate %q", r)
				}
			}
			for _, tok := range tt.wantToks {
				if !cs.HasToken(tok) {
					t.Errorf("missing token candidate %s (%d)", TokenName(tok), tok)
				}
			}
		})
	}
}
