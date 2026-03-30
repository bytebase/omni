package parser

import (
	"testing"
)

// --- Section 6.1: DECLARE & SET ---

func TestCollect_DeclareSet(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantRules []string
		wantToks  []int
	}{
		{
			name:      "DECLARE |",
			sql:       "DECLARE ",
			wantRules: []string{"@variable"},
		},
		{
			name:      "DECLARE @v |",
			sql:       "DECLARE @v ",
			wantRules: []string{"type_name"},
		},
		{
			name:      "SET |",
			sql:       "SET ",
			wantRules: []string{"@variable"},
			wantToks:  []int{kwNOCOUNT, kwXACT_ABORT},
		},
		{
			name:      "SET @v = |",
			sql:       "SET @v = ",
			wantRules: []string{"columnref", "func_name"},
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

// --- Section 6.2: Control Flow ---

func TestCollect_ControlFlow(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantRules []string
		wantToks  []int
	}{
		{
			name:      "IF |",
			sql:       "IF ",
			wantRules: []string{"columnref", "func_name"},
		},
		{
			name:      "WHILE |",
			sql:       "WHILE ",
			wantRules: []string{"columnref", "func_name"},
		},
		{
			name:     "BEGIN |",
			sql:      "BEGIN ",
			wantToks: []int{kwTRANSACTION, kwTRY, kwSELECT, kwINSERT},
		},
		{
			name:     "BEGIN TRY |",
			sql:      "BEGIN TRY ",
			wantToks: []int{kwSELECT, kwINSERT, kwUPDATE, kwDELETE},
		},
		{
			name:     "BEGIN CATCH |",
			sql:      "BEGIN TRY SELECT 1 END TRY BEGIN CATCH ",
			wantToks: []int{kwSELECT, kwINSERT, kwUPDATE, kwDELETE},
		},
		{
			name:      "RETURN |",
			sql:       "RETURN ",
			wantRules: []string{"columnref", "func_name"},
		},
		{
			name:      "EXEC |",
			sql:       "EXEC ",
			wantRules: []string{"proc_ref"},
		},
		{
			name:      "EXECUTE |",
			sql:       "EXECUTE ",
			wantRules: []string{"proc_ref"},
		},
		{
			name:      "EXEC p @param = |",
			sql:       "EXEC p @param = ",
			wantRules: []string{"columnref", "func_name"},
		},
		{
			name:     "WAITFOR |",
			sql:      "WAITFOR ",
			wantToks: []int{kwDELAY, kwTIME},
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

// --- Section 6.3: Transaction ---

func TestCollect_Transaction(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantRules []string
		wantToks  []int
	}{
		{
			name:      "BEGIN TRANSACTION |",
			sql:       "BEGIN TRANSACTION ",
			wantRules: []string{"identifier"},
		},
		{
			name:     "COMMIT |",
			sql:      "COMMIT ",
			wantToks: []int{kwTRANSACTION},
		},
		{
			name:     "ROLLBACK |",
			sql:      "ROLLBACK ",
			wantToks: []int{kwTRANSACTION},
		},
		{
			name:      "SAVE TRANSACTION |",
			sql:       "SAVE TRANSACTION ",
			wantRules: []string{"identifier"},
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

// --- Section 6.4: Cursor Operations ---

func TestCollect_CursorOperations(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantRules []string
		wantToks  []int
	}{
		{
			name:     "DECLARE c CURSOR FOR |",
			sql:      "DECLARE c CURSOR FOR ",
			wantToks: []int{kwSELECT},
		},
		{
			name:      "OPEN |",
			sql:       "OPEN ",
			wantRules: []string{"cursor_name"},
		},
		{
			name:      "FETCH NEXT FROM |",
			sql:       "FETCH NEXT FROM ",
			wantRules: []string{"cursor_name"},
		},
		{
			name:      "CLOSE |",
			sql:       "CLOSE ",
			wantRules: []string{"cursor_name"},
		},
		{
			name:      "DEALLOCATE |",
			sql:       "DEALLOCATE ",
			wantRules: []string{"cursor_name"},
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

// --- Section 7.1: GRANT, REVOKE, DENY ---

func TestCollect_GrantRevokeDeny(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantRules []string
		wantToks  []int
	}{
		{
			name:     "GRANT |",
			sql:      "GRANT ",
			wantToks: []int{kwSELECT, kwINSERT, kwUPDATE, kwDELETE, kwEXECUTE, kwALTER, kwREFERENCES, kwALL},
		},
		{
			name:      "GRANT SELECT ON |",
			sql:       "GRANT SELECT ON ",
			wantRules: []string{"table_ref"},
		},
		{
			name:      "GRANT SELECT ON t TO |",
			sql:       "GRANT SELECT ON t TO ",
			wantRules: []string{"user_ref", "role_ref"},
		},
		{
			name:     "REVOKE |",
			sql:      "REVOKE ",
			wantToks: []int{kwSELECT, kwINSERT, kwUPDATE, kwDELETE, kwEXECUTE, kwALTER, kwREFERENCES, kwALL},
		},
		{
			name:     "DENY |",
			sql:      "DENY ",
			wantToks: []int{kwSELECT, kwINSERT, kwUPDATE, kwDELETE, kwEXECUTE, kwALTER, kwREFERENCES, kwALL},
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

// --- Section 7.2: USE & Utility ---

func TestCollect_UseUtility(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantRules []string
		wantToks  []int
	}{
		{
			name:      "USE |",
			sql:       "USE ",
			wantRules: []string{"database_ref"},
		},
		{
			name:      "PRINT |",
			sql:       "PRINT ",
			wantRules: []string{"columnref", "func_name"},
		},
		{
			name:      "RAISERROR(|)",
			sql:       "RAISERROR(",
			wantRules: []string{"columnref", "func_name"},
		},
		{
			name:      "THROW |",
			sql:       "THROW ",
			wantRules: []string{"columnref", "func_name"},
		},
		{
			name:     "BACKUP DATABASE |",
			sql:      "BACKUP DATABASE ",
			wantRules: []string{"database_ref"},
		},
		{
			name:     "RESTORE DATABASE |",
			sql:      "RESTORE DATABASE ",
			wantRules: []string{"database_ref"},
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
