package completion

import (
	"testing"

	"github.com/bytebase/omni/mysql/catalog"
)

// containsCandidate returns true if candidates contains one with the given text and type.
func containsCandidate(candidates []Candidate, text string, typ CandidateType) bool {
	for _, c := range candidates {
		if c.Text == text && c.Type == typ {
			return true
		}
	}
	return false
}

// containsText returns true if any candidate has the given text.
func containsText(candidates []Candidate, text string) bool {
	for _, c := range candidates {
		if c.Text == text {
			return true
		}
	}
	return false
}

// hasDuplicates returns true if there are duplicate (text, type) pairs (case-insensitive).
func hasDuplicates(candidates []Candidate) bool {
	type key struct {
		text string
		typ  CandidateType
	}
	seen := make(map[key]bool)
	for _, c := range candidates {
		k := key{text: c.Text, typ: c.Type}
		if seen[k] {
			return true
		}
		seen[k] = true
	}
	return false
}

func TestComplete_2_1_CompleteReturnsSlice(t *testing.T) {
	// Scenario: Complete(sql, cursorOffset, catalog) returns []Candidate
	cat := catalog.New()
	candidates := Complete("SELECT ", 7, cat)
	if candidates == nil {
		// nil is acceptable (no candidates), but the function should not panic
		candidates = []Candidate{}
	}
	// Just verify it returns a slice (type is enforced by compiler).
	_ = candidates
}

func TestComplete_2_1_CandidateFields(t *testing.T) {
	// Scenario: Candidate struct has Text, Type, Definition, Comment fields
	c := Candidate{
		Text:       "SELECT",
		Type:       CandidateKeyword,
		Definition: "SQL SELECT statement",
		Comment:    "Retrieves data",
	}
	if c.Text != "SELECT" {
		t.Errorf("Text = %q, want SELECT", c.Text)
	}
	if c.Type != CandidateKeyword {
		t.Errorf("Type = %d, want CandidateKeyword", c.Type)
	}
	if c.Definition != "SQL SELECT statement" {
		t.Errorf("Definition = %q", c.Definition)
	}
	if c.Comment != "Retrieves data" {
		t.Errorf("Comment = %q", c.Comment)
	}
}

func TestComplete_2_1_CandidateTypeEnum(t *testing.T) {
	// Scenario: CandidateType enum with all types
	types := []CandidateType{
		CandidateKeyword,
		CandidateDatabase,
		CandidateTable,
		CandidateView,
		CandidateColumn,
		CandidateFunction,
		CandidateProcedure,
		CandidateIndex,
		CandidateTrigger,
		CandidateEvent,
		CandidateVariable,
		CandidateCharset,
		CandidateEngine,
		CandidateType_,
	}
	// All types should be distinct.
	seen := make(map[CandidateType]bool)
	for _, ct := range types {
		if seen[ct] {
			t.Errorf("duplicate CandidateType value %d", ct)
		}
		seen[ct] = true
	}
	if len(types) != 14 {
		t.Errorf("expected 14 CandidateType values, got %d", len(types))
	}
}

func TestComplete_2_1_NilCatalog(t *testing.T) {
	// Scenario: Complete with nil catalog returns keyword-only candidates
	// (plus built-in function names, which are always available regardless of catalog).
	candidates := Complete("SELECT ", 7, nil)
	for _, c := range candidates {
		if c.Type != CandidateKeyword && c.Type != CandidateFunction {
			t.Errorf("with nil catalog, got unexpected candidate type: %+v", c)
		}
	}
	// Should still return some keywords (e.g., DISTINCT, ALL from SELECT context).
	if len(candidates) == 0 {
		t.Error("expected some keyword candidates with nil catalog")
	}
	// No catalog-dependent types should appear.
	for _, c := range candidates {
		switch c.Type {
		case CandidateTable, CandidateView, CandidateColumn, CandidateDatabase,
			CandidateProcedure, CandidateIndex, CandidateTrigger, CandidateEvent:
			t.Errorf("with nil catalog, got catalog-dependent candidate: %+v", c)
		}
	}
}

func TestComplete_2_1_EmptySQL(t *testing.T) {
	// Scenario: Complete with empty sql returns top-level statement keywords
	candidates := Complete("", 0, nil)
	if len(candidates) == 0 {
		t.Fatal("expected top-level keywords for empty SQL")
	}
	// Should contain core statement keywords.
	for _, kw := range []string{"SELECT", "INSERT", "UPDATE", "DELETE", "CREATE", "ALTER", "DROP"} {
		if !containsCandidate(candidates, kw, CandidateKeyword) {
			t.Errorf("missing expected keyword %s", kw)
		}
	}
	// All should be keywords.
	for _, c := range candidates {
		if c.Type != CandidateKeyword {
			t.Errorf("non-keyword candidate in empty SQL: %+v", c)
		}
	}
}

func TestComplete_2_1_PrefixFiltering(t *testing.T) {
	// Scenario: Prefix filtering: `SEL|` matches SELECT keyword
	candidates := Complete("SEL", 3, nil)
	if !containsCandidate(candidates, "SELECT", CandidateKeyword) {
		t.Error("expected SELECT in candidates for prefix SEL")
	}
	// Should not contain non-matching keywords.
	if containsCandidate(candidates, "INSERT", CandidateKeyword) {
		t.Error("INSERT should not match prefix SEL")
	}
}

func TestComplete_2_1_PrefixCaseInsensitive(t *testing.T) {
	// Scenario: Prefix filtering is case-insensitive
	candidates := Complete("sel", 3, nil)
	if !containsCandidate(candidates, "SELECT", CandidateKeyword) {
		t.Error("expected SELECT in candidates for lowercase prefix sel")
	}
	// Mixed case
	candidates2 := Complete("Sel", 3, nil)
	if !containsCandidate(candidates2, "SELECT", CandidateKeyword) {
		t.Error("expected SELECT in candidates for mixed-case prefix Sel")
	}
}

func TestComplete_2_1_Deduplication(t *testing.T) {
	// Scenario: Deduplication: same candidate not returned twice
	// Use a context that might produce duplicate token candidates.
	candidates := Complete("", 0, nil)
	if hasDuplicates(candidates) {
		t.Error("found duplicate candidates in results")
	}

	// Also test with a prefix context.
	candidates2 := Complete("SELECT ", 7, nil)
	if hasDuplicates(candidates2) {
		t.Error("found duplicate candidates in SELECT context")
	}
}

// --- Section 2.2: Candidate Resolution ---

// setupCatalog creates a catalog with a test database for resolution tests.
func setupCatalog(t *testing.T) *catalog.Catalog {
	t.Helper()
	cat := catalog.New()
	mustExec(t, cat, "CREATE DATABASE testdb")
	cat.SetCurrentDatabase("testdb")
	mustExec(t, cat, "CREATE TABLE users (id INT, name VARCHAR(100), email VARCHAR(200))")
	mustExec(t, cat, "CREATE TABLE orders (id INT, user_id INT, total DECIMAL(10,2))")
	mustExec(t, cat, "CREATE INDEX idx_name ON users (name)")
	mustExec(t, cat, "CREATE INDEX idx_user_id ON orders (user_id)")
	mustExec(t, cat, "CREATE VIEW active_users AS SELECT * FROM users WHERE id > 0")
	mustExec(t, cat, "CREATE FUNCTION my_func() RETURNS INT DETERMINISTIC RETURN 1")
	mustExec(t, cat, "CREATE PROCEDURE my_proc() BEGIN SELECT 1; END")
	mustExec(t, cat, "CREATE TRIGGER my_trig BEFORE INSERT ON users FOR EACH ROW SET NEW.name = UPPER(NEW.name)")
	// Event creation requires schedule — use Exec directly.
	mustExec(t, cat, "CREATE EVENT my_event ON SCHEDULE EVERY 1 HOUR DO SELECT 1")
	return cat
}

// mustExec executes SQL on the catalog, failing the test on error.
func mustExec(t *testing.T, cat *catalog.Catalog, sql string) {
	t.Helper()
	if _, err := cat.Exec(sql, nil); err != nil {
		t.Fatalf("Exec(%q) failed: %v", sql, err)
	}
}

func TestResolve_2_2_TokenCandidatesKeywords(t *testing.T) {
	// Scenario: Token candidates -> keyword strings (from token type mapping)
	// Tested via Complete — empty SQL yields token-only candidates resolved as keywords.
	candidates := Complete("", 0, nil)
	if len(candidates) == 0 {
		t.Fatal("expected keyword candidates")
	}
	for _, c := range candidates {
		if c.Type != CandidateKeyword {
			t.Errorf("expected keyword type, got %d for %q", c.Type, c.Text)
		}
	}
}

func TestResolve_2_2_TableRef(t *testing.T) {
	// Scenario: "table_ref" rule -> catalog tables + views
	cat := setupCatalog(t)
	candidates := resolveRule("table_ref", cat)
	if !containsCandidate(candidates, "users", CandidateTable) {
		t.Error("missing table 'users'")
	}
	if !containsCandidate(candidates, "orders", CandidateTable) {
		t.Error("missing table 'orders'")
	}
	if !containsCandidate(candidates, "active_users", CandidateView) {
		t.Error("missing view 'active_users'")
	}
}

func TestResolve_2_2_ColumnRef(t *testing.T) {
	// Scenario: "columnref" rule -> columns from tables in scope
	// For now, returns all columns from all tables in current database.
	cat := setupCatalog(t)
	candidates := resolveRule("columnref", cat)
	// users: id, name, email
	if !containsCandidate(candidates, "id", CandidateColumn) {
		t.Error("missing column 'id'")
	}
	if !containsCandidate(candidates, "name", CandidateColumn) {
		t.Error("missing column 'name'")
	}
	if !containsCandidate(candidates, "email", CandidateColumn) {
		t.Error("missing column 'email'")
	}
	// orders: user_id, total (id is deduped)
	if !containsCandidate(candidates, "user_id", CandidateColumn) {
		t.Error("missing column 'user_id'")
	}
	if !containsCandidate(candidates, "total", CandidateColumn) {
		t.Error("missing column 'total'")
	}
}

func TestResolve_2_2_DatabaseRef(t *testing.T) {
	// Scenario: "database_ref" rule -> catalog databases
	cat := setupCatalog(t)
	// Add another database.
	mustExec(t, cat, "CREATE DATABASE otherdb")
	candidates := resolveRule("database_ref", cat)
	if !containsCandidate(candidates, "testdb", CandidateDatabase) {
		t.Error("missing database 'testdb'")
	}
	if !containsCandidate(candidates, "otherdb", CandidateDatabase) {
		t.Error("missing database 'otherdb'")
	}
}

func TestResolve_2_2_FunctionRef(t *testing.T) {
	// Scenario: "function_ref" / "func_name" rule -> catalog functions + built-in names
	cat := setupCatalog(t)
	for _, rule := range []string{"function_ref", "func_name"} {
		candidates := resolveRule(rule, cat)
		// Should include built-in functions.
		if !containsCandidate(candidates, "COUNT", CandidateFunction) {
			t.Errorf("[%s] missing built-in function COUNT", rule)
		}
		if !containsCandidate(candidates, "CONCAT", CandidateFunction) {
			t.Errorf("[%s] missing built-in function CONCAT", rule)
		}
		if !containsCandidate(candidates, "NOW", CandidateFunction) {
			t.Errorf("[%s] missing built-in function NOW", rule)
		}
		// Should include catalog function.
		if !containsCandidate(candidates, "my_func", CandidateFunction) {
			t.Errorf("[%s] missing catalog function 'my_func'", rule)
		}
	}
}

func TestResolve_2_2_ProcedureRef(t *testing.T) {
	// Scenario: "procedure_ref" rule -> catalog procedures
	cat := setupCatalog(t)
	candidates := resolveRule("procedure_ref", cat)
	if !containsCandidate(candidates, "my_proc", CandidateProcedure) {
		t.Error("missing procedure 'my_proc'")
	}
}

func TestResolve_2_2_IndexRef(t *testing.T) {
	// Scenario: "index_ref" rule -> indexes from relevant table
	cat := setupCatalog(t)
	candidates := resolveRule("index_ref", cat)
	if !containsCandidate(candidates, "idx_name", CandidateIndex) {
		t.Error("missing index 'idx_name'")
	}
	if !containsCandidate(candidates, "idx_user_id", CandidateIndex) {
		t.Error("missing index 'idx_user_id'")
	}
}

func TestResolve_2_2_TriggerRef(t *testing.T) {
	// Scenario: "trigger_ref" rule -> catalog triggers
	cat := setupCatalog(t)
	candidates := resolveRule("trigger_ref", cat)
	if !containsCandidate(candidates, "my_trig", CandidateTrigger) {
		t.Error("missing trigger 'my_trig'")
	}
}

func TestResolve_2_2_EventRef(t *testing.T) {
	// Scenario: "event_ref" rule -> catalog events
	cat := setupCatalog(t)
	candidates := resolveRule("event_ref", cat)
	if !containsCandidate(candidates, "my_event", CandidateEvent) {
		t.Error("missing event 'my_event'")
	}
}

func TestResolve_2_2_ViewRef(t *testing.T) {
	// Scenario: "view_ref" rule -> catalog views
	cat := setupCatalog(t)
	candidates := resolveRule("view_ref", cat)
	if !containsCandidate(candidates, "active_users", CandidateView) {
		t.Error("missing view 'active_users'")
	}
}

func TestResolve_2_2_Charset(t *testing.T) {
	// Scenario: "charset" rule -> known charset names
	candidates := resolveRule("charset", nil)
	for _, cs := range []string{"utf8mb4", "latin1", "utf8", "ascii", "binary"} {
		if !containsCandidate(candidates, cs, CandidateCharset) {
			t.Errorf("missing charset %q", cs)
		}
	}
}

func TestResolve_2_2_Engine(t *testing.T) {
	// Scenario: "engine" rule -> known engine names
	candidates := resolveRule("engine", nil)
	for _, eng := range []string{"InnoDB", "MyISAM", "MEMORY", "CSV", "ARCHIVE"} {
		if !containsCandidate(candidates, eng, CandidateEngine) {
			t.Errorf("missing engine %q", eng)
		}
	}
}

func TestResolve_2_2_TypeName(t *testing.T) {
	// Scenario: "type_name" rule -> MySQL type keywords
	candidates := resolveRule("type_name", nil)
	for _, typ := range []string{"INT", "VARCHAR", "TEXT", "BLOB", "DATE", "DATETIME", "DECIMAL", "JSON", "ENUM"} {
		if !containsCandidate(candidates, typ, CandidateType_) {
			t.Errorf("missing type %q", typ)
		}
	}
}

func TestResolve_2_2_NilCatalogSafety(t *testing.T) {
	// All catalog-dependent rules should handle nil catalog gracefully.
	for _, rule := range []string{"table_ref", "columnref", "database_ref", "procedure_ref", "index_ref", "trigger_ref", "event_ref", "view_ref"} {
		candidates := resolveRule(rule, nil)
		if candidates != nil && len(candidates) > 0 {
			t.Errorf("[%s] expected no candidates with nil catalog, got %d", rule, len(candidates))
		}
	}
	// function_ref/func_name still return built-ins with nil catalog.
	candidates := resolveRule("func_name", nil)
	if len(candidates) == 0 {
		t.Error("func_name should return built-in functions even with nil catalog")
	}
}
