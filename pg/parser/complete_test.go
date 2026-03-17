package parser

import "testing"

func TestCollectStmtKeywords(t *testing.T) {
	cs := Collect("", 0)
	if cs == nil {
		t.Fatal("Collect returned nil")
	}

	// At the start of a statement, the parser should offer statement-starting keywords.
	want := []int{SELECT, INSERT, CREATE, ALTER, DROP, UPDATE, DELETE_P, WITH,
		SET, SHOW, GRANT, REVOKE, TRUNCATE, BEGIN_P, COMMIT, ROLLBACK}

	for _, tok := range want {
		if !cs.HasToken(tok) {
			t.Errorf("expected token %d in candidates, but not found", tok)
		}
	}
}

func TestCollectBasic(t *testing.T) {
	cs := newCandidateSet()

	cs.addToken(SELECT)
	cs.addToken(INSERT)
	cs.addToken(SELECT) // duplicate

	if len(cs.Tokens) != 2 {
		t.Errorf("expected 2 tokens, got %d", len(cs.Tokens))
	}
	if !cs.HasToken(SELECT) {
		t.Error("expected HasToken(SELECT) = true")
	}
	if cs.HasToken(UPDATE) {
		t.Error("expected HasToken(UPDATE) = false")
	}

	cs.addRule("table_ref")
	cs.addRule("table_ref") // duplicate
	cs.addRule("column_ref")

	if len(cs.Rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(cs.Rules))
	}
	if !cs.HasRule("table_ref") {
		t.Error("expected HasRule(table_ref) = true")
	}
	if cs.HasRule("expr") {
		t.Error("expected HasRule(expr) = false")
	}
}

func TestCollectAfterSelect(t *testing.T) {
	// "SELECT " — cursor after SELECT keyword
	candidates := Collect("SELECT ", 7)
	if candidates == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !candidates.HasRule("columnref") {
		t.Error("expected columnref rule candidate in SELECT target list")
	}
	if !candidates.HasRule("func_name") {
		t.Error("expected func_name rule candidate in SELECT target list")
	}
	// Should also have expression-starting keywords
	if !candidates.HasToken(CASE) {
		t.Error("expected CASE token in expression context")
	}
	if !candidates.HasToken(EXISTS) {
		t.Error("expected EXISTS token in expression context")
	}
}

func TestCollectAfterFrom(t *testing.T) {
	candidates := Collect("SELECT 1 FROM ", 14)
	if candidates == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !candidates.HasRule("relation_expr") {
		t.Error("expected relation_expr rule candidate after FROM")
	}
}

func TestCollectInfixOperators(t *testing.T) {
	// "SELECT 1 FROM t WHERE x " — after expression, should suggest operators
	candidates := Collect("SELECT 1 FROM t WHERE x ", 24)
	if candidates == nil {
		t.Fatal("expected non-nil candidates")
	}
	for _, tok := range []int{AND, OR, IS, '<', '>', '=', BETWEEN, IN_P, LIKE} {
		if !candidates.HasToken(tok) {
			t.Errorf("missing infix operator token: %d", tok)
		}
	}
}

func TestCollectSelectClauses(t *testing.T) {
	// "SELECT 1 " — after target list
	candidates := Collect("SELECT 1 ", 9)
	if candidates == nil {
		t.Fatal("expected non-nil candidates")
	}
	for _, tok := range []int{FROM, WHERE, GROUP_P, ORDER, LIMIT} {
		if !candidates.HasToken(tok) {
			t.Errorf("missing clause keyword after target list: %d", tok)
		}
	}
}

func TestCollectAfterFromClause(t *testing.T) {
	// "SELECT 1 FROM t " — after FROM clause table
	candidates := Collect("SELECT 1 FROM t ", 16)
	if candidates == nil {
		t.Fatal("expected non-nil candidates")
	}
	for _, tok := range []int{WHERE, GROUP_P, ORDER, LIMIT} {
		if !candidates.HasToken(tok) {
			t.Errorf("missing clause keyword after FROM: %d", tok)
		}
	}
}

func TestCollectAfterWhere(t *testing.T) {
	// "SELECT 1 FROM t WHERE x = 1 " — after WHERE predicate
	candidates := Collect("SELECT 1 FROM t WHERE x = 1 ", 28)
	if candidates == nil {
		t.Fatal("expected non-nil candidates")
	}
	// Should have AND/OR (expression operators) and clause keywords
	if !candidates.HasToken(AND) {
		t.Error("expected AND after WHERE predicate")
	}
	if !candidates.HasToken(ORDER) {
		t.Error("expected ORDER after WHERE predicate")
	}
}

func TestCollectAfterCreate(t *testing.T) {
	candidates := Collect("CREATE ", 7)
	if candidates == nil {
		t.Fatal("expected non-nil candidates")
	}
	for _, tok := range []int{TABLE, VIEW, INDEX, FUNCTION, SCHEMA, DATABASE} {
		if !candidates.HasToken(tok) {
			t.Errorf("missing CREATE sub-keyword: %d", tok)
		}
	}
}

func TestCollectAfterAlter(t *testing.T) {
	candidates := Collect("ALTER ", 6)
	if candidates == nil {
		t.Fatal("expected non-nil candidates")
	}
	for _, tok := range []int{TABLE, DATABASE, ROLE, SCHEMA, FUNCTION} {
		if !candidates.HasToken(tok) {
			t.Errorf("missing ALTER sub-keyword: %d", tok)
		}
	}
}

func TestCollectAfterDrop(t *testing.T) {
	candidates := Collect("DROP ", 5)
	if candidates == nil {
		t.Fatal("expected non-nil candidates")
	}
	for _, tok := range []int{TABLE, VIEW, INDEX, FUNCTION, SCHEMA} {
		if !candidates.HasToken(tok) {
			t.Errorf("missing DROP sub-keyword: %d", tok)
		}
	}
}

func TestCollectKeywordCategories(t *testing.T) {
	candidates := Collect("SELECT ", 7)
	// Unreserved keywords like NAME should be candidates (valid as identifiers)
	if !candidates.HasToken(NAME_P) {
		t.Error("expected unreserved keyword NAME_P as candidate")
	}
}
