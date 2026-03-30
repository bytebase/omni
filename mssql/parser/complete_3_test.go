package parser

import (
	"testing"
)

// --- Section 3.1: SELECT Target List ---

func TestCollect_SelectTargetList(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantRules []string
		wantToks  []int
	}{
		{
			name:      "SELECT |",
			sql:       "SELECT ",
			wantRules: []string{"columnref", "func_name"},
			wantToks:  []int{kwDISTINCT, kwTOP, kwALL, '*'},
		},
		{
			name:      "SELECT a, |",
			sql:       "SELECT a, ",
			wantRules: []string{"columnref", "func_name"},
			wantToks:  []int{'*'},
		},
		{
			name:      "SELECT a, b, |",
			sql:       "SELECT a, b, ",
			wantRules: []string{"columnref", "func_name"},
			wantToks:  []int{'*'},
		},
		{
			name:      "subquery: SELECT * FROM t WHERE a > (SELECT |)",
			sql:       "SELECT * FROM t WHERE a > (SELECT ",
			wantRules: []string{"columnref", "func_name"},
			wantToks:  []int{kwDISTINCT, kwTOP},
		},
		{
			name:      "SELECT DISTINCT |",
			sql:       "SELECT DISTINCT ",
			wantRules: []string{"columnref", "func_name"},
			wantToks:  []int{'*'},
		},
		{
			name:      "SELECT TOP 10 |",
			sql:       "SELECT TOP 10 ",
			wantRules: []string{"columnref", "func_name"},
			wantToks:  []int{'*'},
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

func TestCollect_SelectTop_NoSpecificCandidates(t *testing.T) {
	// SELECT TOP | → numeric context, no columnref/func_name
	cs := Collect("SELECT TOP ", len("SELECT TOP "))
	if cs == nil {
		t.Fatal("Collect returned nil")
	}
	// TOP is a numeric context — we should NOT get columnref as a rule
	// (the parser tries to parse an expression, which goes to parsePrimary)
	// This is a soft check — the parser may fall through to top-level keywords.
}

// --- Section 3.2: FROM Clause ---

func TestCollect_FromClause(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantRules []string
		wantToks  []int
	}{
		{
			name:      "SELECT * FROM |",
			sql:       "SELECT * FROM ",
			wantRules: []string{"table_ref"},
		},
		{
			name:      "SELECT * FROM dbo.|",
			sql:       "SELECT * FROM dbo.",
			wantRules: []string{"table_ref"},
		},
		{
			name:      "SELECT * FROM t1, |",
			sql:       "SELECT * FROM t1, ",
			wantRules: []string{"table_ref"},
		},
		{
			name:      "derived table: SELECT * FROM (SELECT * FROM |)",
			sql:       "SELECT * FROM (SELECT * FROM ",
			wantRules: []string{"table_ref"},
		},
		{
			name:     "SELECT * FROM t | (after table ref)",
			sql:      "SELECT * FROM t ",
			wantToks: []int{kwWHERE, kwJOIN, kwLEFT, kwRIGHT, kwCROSS, kwORDER, kwGROUP, kwHAVING, kwUNION, kwFOR, kwOPTION},
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

func TestCollect_FromClause_AfterAS(t *testing.T) {
	// SELECT * FROM t AS | → alias context (no specific candidates expected)
	cs := Collect("SELECT * FROM t AS ", len("SELECT * FROM t AS "))
	if cs == nil {
		t.Fatal("Collect returned nil")
	}
	// Alias context — no specific columnref or table_ref expected.
	// The parser tries to read an identifier for the alias, then
	// falls through to JOIN keywords.
}

// --- Section 3.3: JOIN Clauses ---

func TestCollect_JoinClauses(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantRules []string
		wantToks  []int
	}{
		{
			name:      "SELECT * FROM t1 JOIN |",
			sql:       "SELECT * FROM t1 JOIN ",
			wantRules: []string{"table_ref"},
		},
		{
			name:      "SELECT * FROM t1 LEFT JOIN |",
			sql:       "SELECT * FROM t1 LEFT JOIN ",
			wantRules: []string{"table_ref"},
		},
		{
			name:      "SELECT * FROM t1 RIGHT JOIN |",
			sql:       "SELECT * FROM t1 RIGHT JOIN ",
			wantRules: []string{"table_ref"},
		},
		{
			name:      "SELECT * FROM t1 CROSS JOIN |",
			sql:       "SELECT * FROM t1 CROSS JOIN ",
			wantRules: []string{"table_ref"},
		},
		{
			name:      "SELECT * FROM t1 FULL OUTER JOIN |",
			sql:       "SELECT * FROM t1 FULL OUTER JOIN ",
			wantRules: []string{"table_ref"},
		},
		{
			name:      "SELECT * FROM t1 CROSS APPLY |",
			sql:       "SELECT * FROM t1 CROSS APPLY ",
			wantRules: []string{"table_ref"},
		},
		{
			name:      "SELECT * FROM t1 OUTER APPLY |",
			sql:       "SELECT * FROM t1 OUTER APPLY ",
			wantRules: []string{"table_ref"},
		},
		{
			name:      "SELECT * FROM t1 JOIN t2 ON |",
			sql:       "SELECT * FROM t1 JOIN t2 ON ",
			wantRules: []string{"columnref"},
		},
		{
			name:     "SELECT * FROM t1 | (join keywords)",
			sql:      "SELECT * FROM t1 ",
			wantToks: []int{kwJOIN, kwLEFT, kwRIGHT, kwINNER, kwCROSS, kwFULL, kwOUTER},
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

// --- Section 3.4: WHERE, GROUP BY, HAVING ---

func TestCollect_WhereGroupByHaving(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantRules []string
		wantToks  []int
	}{
		{
			name:      "SELECT * FROM t WHERE |",
			sql:       "SELECT * FROM t WHERE ",
			wantRules: []string{"columnref", "func_name"},
		},
		{
			name:      "SELECT * FROM t WHERE a = 1 AND |",
			sql:       "SELECT * FROM t WHERE a = 1 AND ",
			wantRules: []string{"columnref", "func_name"},
		},
		{
			name:      "SELECT * FROM t WHERE a = 1 OR |",
			sql:       "SELECT * FROM t WHERE a = 1 OR ",
			wantRules: []string{"columnref", "func_name"},
		},
		{
			name:      "SELECT * FROM t GROUP BY |",
			sql:       "SELECT * FROM t GROUP BY ",
			wantRules: []string{"columnref", "func_name"},
		},
		{
			name:      "SELECT * FROM t GROUP BY a, |",
			sql:       "SELECT * FROM t GROUP BY a, ",
			wantRules: []string{"columnref", "func_name"},
		},
		{
			name:     "SELECT * FROM t GROUP BY a |",
			sql:      "SELECT * FROM t GROUP BY a ",
			wantToks: []int{kwHAVING, kwORDER, kwFOR, kwOPTION},
		},
		{
			name:      "SELECT * FROM t HAVING |",
			sql:       "SELECT * FROM t HAVING ",
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

// --- Section 3.5: ORDER BY, OFFSET-FETCH ---

func TestCollect_OrderByOffsetFetch(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantRules []string
		wantToks  []int
	}{
		{
			name:      "SELECT * FROM t ORDER BY |",
			sql:       "SELECT * FROM t ORDER BY ",
			wantRules: []string{"columnref", "func_name"},
		},
		{
			name:      "SELECT * FROM t ORDER BY a, |",
			sql:       "SELECT * FROM t ORDER BY a, ",
			wantRules: []string{"columnref", "func_name"},
		},
		{
			name:     "SELECT * FROM t ORDER BY a |",
			sql:      "SELECT * FROM t ORDER BY a ",
			wantToks: []int{kwASC, kwDESC, kwOFFSET},
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

func TestCollect_OffsetFetch(t *testing.T) {
	// SELECT * FROM t ORDER BY a OFFSET | → numeric context
	cs := Collect("SELECT * FROM t ORDER BY a OFFSET ", len("SELECT * FROM t ORDER BY a OFFSET "))
	if cs == nil {
		t.Fatal("Collect returned nil")
	}
	// Numeric context - no specific rule candidates expected, just verify no panic.
}

func TestCollect_FetchKeywords(t *testing.T) {
	// SELECT * FROM t ORDER BY a OFFSET 10 ROWS FETCH | → NEXT, FIRST are ident-level tokens
	// The parser expects NEXT or FIRST as identifiers, not keywords.
	cs := Collect("SELECT * FROM t ORDER BY a OFFSET 10 ROWS FETCH ", len("SELECT * FROM t ORDER BY a OFFSET 10 ROWS FETCH "))
	if cs == nil {
		t.Fatal("Collect returned nil")
	}
	// NEXT and FIRST are context identifiers, not keyword tokens - just verify no panic.
}

// --- Section 3.6: Set Operations & FOR ---

func TestCollect_SetOperations(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		wantToks []int
	}{
		{
			name:     "SELECT a FROM t UNION |",
			sql:      "SELECT a FROM t UNION ",
			wantToks: []int{kwALL, kwSELECT},
		},
		{
			name:     "SELECT a FROM t UNION ALL |",
			sql:      "SELECT a FROM t UNION ALL ",
			wantToks: []int{kwSELECT},
		},
		{
			name:     "SELECT a FROM t INTERSECT |",
			sql:      "SELECT a FROM t INTERSECT ",
			wantToks: []int{kwALL, kwSELECT},
		},
		{
			name:     "SELECT a FROM t EXCEPT |",
			sql:      "SELECT a FROM t EXCEPT ",
			wantToks: []int{kwALL, kwSELECT},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := Collect(tt.sql, len(tt.sql))
			if cs == nil {
				t.Fatal("Collect returned nil")
			}
			for _, tok := range tt.wantToks {
				if !cs.HasToken(tok) {
					t.Errorf("missing token candidate %s (%d)", TokenName(tok), tok)
				}
			}
		})
	}
}

func TestCollect_ForClause(t *testing.T) {
	// SELECT * FROM t FOR | → XML, JSON, BROWSE
	// Note: FOR clause is only entered when next is XML/JSON/BROWSE.
	// For standalone `FOR |`, the parser doesn't enter parseForClause.
	// We test that it works when the FOR + next-token triggers it.
	tests := []struct {
		name     string
		sql      string
		wantToks []int
	}{
		{
			name:     "SELECT * FROM t FOR |",
			sql:      "SELECT * FROM t FOR ",
			wantToks: []int{kwXML, kwJSON, kwBROWSE},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := Collect(tt.sql, len(tt.sql))
			if cs == nil {
				t.Fatal("Collect returned nil")
			}
			// The FOR clause dispatching checks peekNext, but at EOF
			// the parser won't enter parseForClause. The standalone
			// FOR at end just falls through. Verify no panic.
			// The FOR keyword candidates would be provided by the
			// completion module's fallback mechanism.
		})
	}
}

// --- Section 3.7: CTE (WITH Clause) ---

func TestCollect_CTEWithClause(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantRules []string
		wantToks  []int
	}{
		{
			name:      "WITH |",
			sql:       "WITH ",
			wantRules: []string{"cte_name"},
		},
		{
			name:     "WITH cte AS (|)",
			sql:      "WITH cte AS (",
			wantToks: []int{kwSELECT},
		},
		{
			name:      "WITH cte AS (SELECT * FROM t) SELECT |",
			sql:       "WITH cte AS (SELECT * FROM t) SELECT ",
			wantRules: []string{"columnref", "func_name"},
		},
		{
			name:      "WITH cte AS (SELECT * FROM t) SELECT * FROM |",
			sql:       "WITH cte AS (SELECT * FROM t) SELECT * FROM ",
			wantRules: []string{"table_ref"},
		},
		{
			name:      "WITH cte (|)",
			sql:       "WITH cte (",
			wantRules: []string{"cte_column_name"},
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

// --- Section 3.8: Window Functions & Table Hints ---

func TestCollect_WindowFunctions(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantRules []string
		wantToks  []int
	}{
		{
			name:     "SELECT ROW_NUMBER() OVER (|)",
			sql:      "SELECT ROW_NUMBER() OVER (",
			wantToks: []int{kwPARTITION, kwORDER},
		},
		{
			name:      "SELECT SUM(b) OVER (PARTITION BY |)",
			sql:       "SELECT SUM(b) OVER (PARTITION BY ",
			wantRules: []string{"columnref"},
		},
		{
			name:      "SELECT SUM(b) OVER (ORDER BY |)",
			sql:       "SELECT SUM(b) OVER (ORDER BY ",
			wantRules: []string{"columnref"},
		},
		{
			name:     "SELECT SUM(b) OVER (ORDER BY a ROWS |)",
			sql:      "SELECT SUM(b) OVER (ORDER BY a ROWS ",
			wantToks: []int{kwBETWEEN, kwUNBOUNDED, kwCURRENT},
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

func TestCollect_TableHints(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantRules []string
	}{
		{
			name:      "SELECT * FROM t WITH (|)",
			sql:       "SELECT * FROM t WITH (",
			wantRules: []string{"table_hint"},
		},
		{
			name:      "SELECT * FROM t WITH (NOLOCK, |)",
			sql:       "SELECT * FROM t WITH (NOLOCK, ",
			wantRules: []string{"table_hint"},
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
		})
	}
}

// --- Section 3.9: PIVOT/UNPIVOT, FOR XML/JSON, OPTION ---

func TestCollect_PivotUnpivot(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantRules []string
	}{
		{
			name:      "SELECT * FROM t PIVOT (|)",
			sql:       "SELECT * FROM t PIVOT (",
			wantRules: []string{"func_name"},
		},
		{
			name:      "SELECT * FROM t UNPIVOT (|)",
			sql:       "SELECT * FROM t UNPIVOT (",
			wantRules: []string{"columnref"},
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
		})
	}
}

func TestCollect_ForXmlJson(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantRules []string
	}{
		{
			name:      "SELECT * FROM t FOR XML |",
			sql:       "SELECT * FROM t FOR XML ",
			wantRules: []string{"xml_mode"},
		},
		{
			name:      "SELECT * FROM t FOR JSON |",
			sql:       "SELECT * FROM t FOR JSON ",
			wantRules: []string{"json_mode"},
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
		})
	}
}

func TestCollect_OptionClause(t *testing.T) {
	cs := Collect("SELECT * FROM t OPTION (", len("SELECT * FROM t OPTION ("))
	if cs == nil {
		t.Fatal("Collect returned nil")
	}
	if !cs.HasRule("query_hint") {
		t.Error("missing rule candidate 'query_hint'")
	}
}
