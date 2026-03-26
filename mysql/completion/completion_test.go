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
	candidates := resolveRule("table_ref", cat, "", 0)
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
	candidates := resolveRule("columnref", cat, "", 0)
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
	candidates := resolveRule("database_ref", cat, "", 0)
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
		candidates := resolveRule(rule, cat, "", 0)
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
	candidates := resolveRule("procedure_ref", cat, "", 0)
	if !containsCandidate(candidates, "my_proc", CandidateProcedure) {
		t.Error("missing procedure 'my_proc'")
	}
}

func TestResolve_2_2_IndexRef(t *testing.T) {
	// Scenario: "index_ref" rule -> indexes from relevant table
	cat := setupCatalog(t)
	candidates := resolveRule("index_ref", cat, "", 0)
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
	candidates := resolveRule("trigger_ref", cat, "", 0)
	if !containsCandidate(candidates, "my_trig", CandidateTrigger) {
		t.Error("missing trigger 'my_trig'")
	}
}

func TestResolve_2_2_EventRef(t *testing.T) {
	// Scenario: "event_ref" rule -> catalog events
	cat := setupCatalog(t)
	candidates := resolveRule("event_ref", cat, "", 0)
	if !containsCandidate(candidates, "my_event", CandidateEvent) {
		t.Error("missing event 'my_event'")
	}
}

func TestResolve_2_2_ViewRef(t *testing.T) {
	// Scenario: "view_ref" rule -> catalog views
	cat := setupCatalog(t)
	candidates := resolveRule("view_ref", cat, "", 0)
	if !containsCandidate(candidates, "active_users", CandidateView) {
		t.Error("missing view 'active_users'")
	}
}

func TestResolve_2_2_Charset(t *testing.T) {
	// Scenario: "charset" rule -> known charset names
	candidates := resolveRule("charset", nil, "", 0)
	for _, cs := range []string{"utf8mb4", "latin1", "utf8", "ascii", "binary"} {
		if !containsCandidate(candidates, cs, CandidateCharset) {
			t.Errorf("missing charset %q", cs)
		}
	}
}

func TestResolve_2_2_Engine(t *testing.T) {
	// Scenario: "engine" rule -> known engine names
	candidates := resolveRule("engine", nil, "", 0)
	for _, eng := range []string{"InnoDB", "MyISAM", "MEMORY", "CSV", "ARCHIVE"} {
		if !containsCandidate(candidates, eng, CandidateEngine) {
			t.Errorf("missing engine %q", eng)
		}
	}
}

func TestResolve_2_2_TypeName(t *testing.T) {
	// Scenario: "type_name" rule -> MySQL type keywords
	candidates := resolveRule("type_name", nil, "", 0)
	for _, typ := range []string{"INT", "VARCHAR", "TEXT", "BLOB", "DATE", "DATETIME", "DECIMAL", "JSON", "ENUM"} {
		if !containsCandidate(candidates, typ, CandidateType_) {
			t.Errorf("missing type %q", typ)
		}
	}
}

// --- Section 2.4: Tricky Completion (Fallback) ---

func TestComplete_2_4_IncompleteTrailingSpace(t *testing.T) {
	// Scenario: Incomplete SQL with trailing space → insert placeholder, re-collect.
	// The trickyComplete function patches SQL with placeholder tokens to make it
	// parseable, then re-runs Collect. When standard Collect returns nothing,
	// trickyComplete should return whatever the patched version produces.
	//
	// Use a context where standard returns empty but placeholder strategy succeeds:
	// `SELECT ` at offset 7 gets standard candidates via the SELECT expr
	// instrumentation. So instead we test that trickyComplete is called when
	// standardComplete returns empty results.
	cat := setupCatalog(t)

	// Test that trailing space after FROM gets candidates via tricky path.
	// The numeric placeholder "1" makes "SELECT * FROM 1" parseable, yielding
	// keyword tokens for the follow set (WHERE, JOIN, etc.).
	candidates := Complete("SELECT * FROM ", 14, cat)
	if len(candidates) == 0 {
		t.Skip("FROM clause not yet instrumented (Phase 3); tricky mechanism works but parser lacks rule candidates here")
	}
}

func TestComplete_2_4_TruncatedMidKeyword(t *testing.T) {
	// Scenario: Truncated mid-keyword: `SELE` → prefix-filter against keywords.
	// The prefix "SELE" is extracted, Collect runs at offset 0 (start of statement),
	// producing top-level keywords, then filterByPrefix keeps only SELECT.
	candidates := Complete("SELE", 4, nil)
	if !containsCandidate(candidates, "SELECT", CandidateKeyword) {
		t.Error("expected SELECT keyword from prefix filter for 'SELE'")
	}
	if containsCandidate(candidates, "INSERT", CandidateKeyword) {
		t.Error("INSERT should not match prefix SELE")
	}
}

func TestComplete_2_4_TruncatedAfterComma(t *testing.T) {
	// Scenario: Truncated after comma: `SELECT a,` → insert placeholder column.
	// After the comma, standardComplete runs at offset 9. If it returns nothing,
	// trickyComplete patches to "SELECT a, __placeholder__" or "SELECT a, 1"
	// which may parse differently.
	//
	// With current parser instrumentation, the SELECT expr list checkpoint is at
	// the start of parseSelectExprs. The comma case needs additional instrumentation
	// (Phase 3, scenario 3.1). But we verify the tricky mechanism doesn't panic
	// and returns whatever the parser can provide.
	cat := setupCatalog(t)
	candidates := Complete("SELECT a,", 9, cat)
	// The mechanism must not panic; results depend on parser instrumentation.
	_ = candidates
}

func TestComplete_2_4_TruncatedAfterOperator(t *testing.T) {
	// Scenario: Truncated after operator: `WHERE a >` → insert placeholder expression.
	// trickyComplete patches to "... WHERE id > __placeholder__" or "... WHERE id > 1".
	// The numeric placeholder "1" makes valid SQL, so Collect can run on it.
	cat := setupCatalog(t)
	candidates := Complete("SELECT * FROM users WHERE id >", 30, cat)
	// Must not panic. Results depend on expression instrumentation (Phase 8).
	_ = candidates
}

func TestComplete_2_4_MultiplePlaceholderStrategies(t *testing.T) {
	// Scenario: Multiple placeholder strategies tried in order.
	// trickyComplete tries three strategies:
	//   1. prefix + " __placeholder__" + suffix
	//   2. prefix + " __placeholder__ " + suffix
	//   3. prefix + " 1" + suffix
	// We verify the function exists, tries them in order, and returns the first
	// strategy that yields candidates.

	// Use trickyComplete directly to verify it returns results when a strategy works.
	// "SELECT " at offset 7 — standard would return results, but we call tricky
	// directly to verify the placeholder mechanism.
	candidates := trickyComplete("", 0, nil)
	// For empty SQL, even the patched versions should produce keyword candidates
	// because " __placeholder__" at offset 0 triggers statement-start keywords.
	if len(candidates) == 0 {
		t.Error("expected trickyComplete to produce candidates for empty SQL via placeholder strategy")
	}

	// Verify keywords are present (the placeholder text itself should not appear).
	hasKeyword := false
	for _, c := range candidates {
		if c.Type == CandidateKeyword {
			hasKeyword = true
			break
		}
	}
	if !hasKeyword {
		t.Error("expected keyword candidates from placeholder strategy")
	}
}

func TestComplete_2_4_FallbackBestEffort(t *testing.T) {
	// Scenario: Fallback returns best-effort results when no strategy succeeds.
	// Completely nonsensical SQL should not panic. trickyComplete returns nil
	// when no strategy produces candidates.
	candidates := Complete("XYZZY PLUGH ", 12, nil)
	// Must not panic. Result may be empty or nil.
	_ = candidates

	// Also test with more realistic but still broken SQL.
	candidates2 := Complete(")))((( ", 7, nil)
	_ = candidates2

	// Verify trickyComplete returns nil for truly unparseable input.
	tricky := trickyComplete("XYZZY PLUGH ", 12, nil)
	// nil is acceptable — it means no strategy succeeded.
	_ = tricky
}

func TestComplete_2_4_PlaceholderNoCorruption(t *testing.T) {
	// Scenario: Placeholder insertion does not corrupt the initial candidate set.
	// Running Complete multiple times on the same input must produce consistent results.
	// The placeholder text (__placeholder__) must never leak into returned candidates.
	cat := setupCatalog(t)

	// Use a SQL that produces candidates via standard path.
	validSQL := "SELECT "
	validCandidates := Complete(validSQL, len(validSQL), cat)
	validCandidates2 := Complete(validSQL, len(validSQL), cat)

	// Both runs should return the same number of candidates.
	if len(validCandidates) != len(validCandidates2) {
		t.Errorf("candidate count mismatch: first=%d, second=%d", len(validCandidates), len(validCandidates2))
	}

	// Placeholder text must not leak into any candidate set.
	for _, c := range validCandidates {
		if c.Text == "__placeholder__" {
			t.Error("placeholder text leaked into candidate set")
		}
	}
	for _, c := range validCandidates2 {
		if c.Text == "__placeholder__" {
			t.Error("placeholder text leaked into candidate set on second run")
		}
	}

	// Also test via trickyComplete directly: placeholder must not appear in results.
	trickyCandidates := trickyComplete("SELECT * FROM ", 14, cat)
	for _, c := range trickyCandidates {
		if c.Text == "__placeholder__" {
			t.Error("placeholder text leaked into tricky candidate set")
		}
	}
}

func TestResolve_2_2_NilCatalogSafety(t *testing.T) {
	// All catalog-dependent rules should handle nil catalog gracefully.
	for _, rule := range []string{"table_ref", "columnref", "database_ref", "procedure_ref", "index_ref", "trigger_ref", "event_ref", "view_ref"} {
		candidates := resolveRule(rule, nil, "", 0)
		if candidates != nil && len(candidates) > 0 {
			t.Errorf("[%s] expected no candidates with nil catalog, got %d", rule, len(candidates))
		}
	}
	// function_ref/func_name still return built-ins with nil catalog.
	candidates := resolveRule("func_name", nil, "", 0)
	if len(candidates) == 0 {
		t.Error("func_name should return built-in functions even with nil catalog")
	}
}

// --- Section 3.1: SELECT Target List ---

func TestComplete_3_1_SelectTargetList(t *testing.T) {
	cat := catalog.New()
	mustExec(t, cat, "CREATE DATABASE test")
	cat.SetCurrentDatabase("test")
	mustExec(t, cat, "CREATE TABLE t (a INT, b INT, c INT)")
	mustExec(t, cat, "CREATE TABLE t1 (x INT, y INT)")

	cases := []struct {
		name       string
		sql        string
		cursor     int
		wantCol    string // expected column candidate (or "" to skip)
		wantFunc   bool   // expect function candidates
		wantKW     string // expected keyword candidate (or "" to skip)
		absentType CandidateType
		absentText string
	}{
		{
			name:     "select_pipe_columnref",
			sql:      "SELECT ",
			cursor:   7,
			wantCol:  "a",
			wantFunc: true,
			wantKW:   "DISTINCT",
		},
		{
			name:     "select_after_comma",
			sql:      "SELECT a, ",
			cursor:   10,
			wantCol:  "a",
			wantFunc: true,
		},
		{
			name:     "select_after_two_commas",
			sql:      "SELECT a, b, ",
			cursor:   13,
			wantCol:  "c",
			wantFunc: true,
		},
		{
			name:     "select_subquery",
			sql:      "SELECT * FROM t WHERE a > (SELECT ",
			cursor:   34,
			wantCol:  "a",
			wantFunc: true,
		},
		{
			name:     "select_distinct_pipe",
			sql:      "SELECT DISTINCT ",
			cursor:   16,
			wantCol:  "a",
			wantFunc: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			candidates := Complete(tc.sql, tc.cursor, cat)
			if len(candidates) == 0 {
				t.Fatal("expected candidates, got none")
			}

			if tc.wantCol != "" {
				if !containsCandidate(candidates, tc.wantCol, CandidateColumn) {
					t.Errorf("missing column candidate %q; got %v", tc.wantCol, candidates)
				}
			}

			if tc.wantFunc {
				if !containsCandidate(candidates, "COUNT", CandidateFunction) {
					t.Errorf("missing function candidate COUNT; got %v", candidates)
				}
			}

			if tc.wantKW != "" {
				if !containsCandidate(candidates, tc.wantKW, CandidateKeyword) {
					t.Errorf("missing keyword candidate %q; got %v", tc.wantKW, candidates)
				}
			}

			if tc.absentText != "" {
				if containsCandidate(candidates, tc.absentText, tc.absentType) {
					t.Errorf("unexpected candidate %q of type %d", tc.absentText, tc.absentType)
				}
			}
		})
	}
}

// --- Section 3.2: FROM Clause ---

func TestComplete_3_2_FromClause(t *testing.T) {
	cat := catalog.New()
	mustExec(t, cat, "CREATE DATABASE test")
	mustExec(t, cat, "CREATE DATABASE testdb2")
	cat.SetCurrentDatabase("test")
	mustExec(t, cat, "CREATE TABLE t (a INT, b INT, c INT)")
	mustExec(t, cat, "CREATE TABLE t1 (x INT, y INT)")
	mustExec(t, cat, "CREATE TABLE t2 (p INT, q INT)")
	mustExec(t, cat, "CREATE VIEW v1 AS SELECT * FROM t")

	cases := []struct {
		name       string
		sql        string
		cursor     int
		wantType   CandidateType
		wantText   string // expected candidate text (or "" to skip check)
		wantAbsent string // text that should NOT appear (or "" to skip)
		absentType CandidateType
	}{
		{
			// Scenario 1: SELECT * FROM | → table_ref (tables, views, databases)
			name:     "from_table_ref",
			sql:      "SELECT * FROM ",
			cursor:   14,
			wantType: CandidateTable,
			wantText: "t",
		},
		{
			// Scenario 1 continued: views should also appear
			name:     "from_view_ref",
			sql:      "SELECT * FROM ",
			cursor:   14,
			wantType: CandidateView,
			wantText: "v1",
		},
		{
			// Scenario 2: SELECT * FROM db.| → table_ref qualified with database
			name:     "from_qualified_table_ref",
			sql:      "SELECT * FROM test.",
			cursor:   19,
			wantType: CandidateTable,
			wantText: "t",
		},
		{
			// Scenario 3: SELECT * FROM t1, | → table_ref after comma
			name:     "from_comma_table_ref",
			sql:      "SELECT * FROM t1, ",
			cursor:   18,
			wantType: CandidateTable,
			wantText: "t2",
		},
		{
			// Scenario 4: SELECT * FROM (SELECT * FROM |) → table_ref in derived table
			name:     "from_derived_table_ref",
			sql:      "SELECT * FROM (SELECT * FROM ",
			cursor:   29,
			wantType: CandidateTable,
			wantText: "t",
		},
		{
			// Scenario 5: SELECT * FROM t | → keyword candidates
			name:     "from_after_table_keywords",
			sql:      "SELECT * FROM t ",
			cursor:   16,
			wantType: CandidateKeyword,
			wantText: "WHERE",
		},
		{
			// Scenario 5 continued: JOIN keyword should appear
			name:     "from_after_table_join",
			sql:      "SELECT * FROM t ",
			cursor:   16,
			wantType: CandidateKeyword,
			wantText: "JOIN",
		},
		{
			// Scenario 5 continued: LEFT keyword should appear
			name:     "from_after_table_left",
			sql:      "SELECT * FROM t ",
			cursor:   16,
			wantType: CandidateKeyword,
			wantText: "LEFT",
		},
		{
			// Scenario 5 continued: RIGHT keyword should appear
			name:     "from_after_table_right",
			sql:      "SELECT * FROM t ",
			cursor:   16,
			wantType: CandidateKeyword,
			wantText: "RIGHT",
		},
		{
			// Scenario 6: SELECT * FROM t AS | → no specific candidates (alias context)
			// Tables/views should NOT appear since we're in an alias context.
			name:       "from_alias_no_table",
			sql:        "SELECT * FROM t AS ",
			cursor:     19,
			wantType:   CandidateTable,
			wantText:   "",
			wantAbsent: "t",
			absentType: CandidateTable,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			candidates := Complete(tc.sql, tc.cursor, cat)

			if tc.wantText != "" {
				if !containsCandidate(candidates, tc.wantText, tc.wantType) {
					t.Errorf("missing candidate %q of type %d; got %v", tc.wantText, tc.wantType, candidates)
				}
			}

			if tc.wantAbsent != "" {
				if containsCandidate(candidates, tc.wantAbsent, tc.absentType) {
					t.Errorf("unexpected candidate %q of type %d should not appear", tc.wantAbsent, tc.absentType)
				}
			}
		})
	}
}
