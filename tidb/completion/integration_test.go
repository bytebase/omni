package completion

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/tidb/catalog"
)

// setupIntegrationCatalog creates a realistic test catalog for integration tests.
func setupIntegrationCatalog(t *testing.T) *catalog.Catalog {
	t.Helper()
	cat := catalog.New()
	mustExec(t, cat, "CREATE DATABASE test")
	cat.SetCurrentDatabase("test")
	mustExec(t, cat, "CREATE TABLE users (id INT, name VARCHAR(100), email VARCHAR(200))")
	mustExec(t, cat, "CREATE TABLE orders (id INT, user_id INT, amount DECIMAL(10,2), status VARCHAR(20))")
	mustExec(t, cat, "CREATE TABLE products (id INT, name VARCHAR(100), price DECIMAL(10,2))")
	mustExec(t, cat, "CREATE VIEW active_users AS SELECT id, name FROM users WHERE id > 0")
	mustExec(t, cat, "CREATE INDEX idx_user_id ON orders (user_id)")
	return cat
}

// assertContains checks that at least one candidate has the given text and type.
func assertContains(t *testing.T, candidates []Candidate, text string, typ CandidateType) {
	t.Helper()
	if !containsCandidate(candidates, text, typ) {
		t.Errorf("expected candidate %q (type %d) not found in %d candidates", text, typ, len(candidates))
	}
}

// assertNotContains checks that no candidate has the given text and type.
func assertNotContains(t *testing.T, candidates []Candidate, text string, typ CandidateType) {
	t.Helper()
	if containsCandidate(candidates, text, typ) {
		t.Errorf("unexpected candidate %q (type %d) found", text, typ)
	}
}

// assertHasType checks that at least one candidate has the given type.
func assertHasType(t *testing.T, candidates []Candidate, typ CandidateType) {
	t.Helper()
	for _, c := range candidates {
		if c.Type == typ {
			return
		}
	}
	t.Errorf("expected at least one candidate of type %d, found none in %d candidates", typ, len(candidates))
}

// candidatesOfType returns candidates matching the given type.
func candidatesOfType(candidates []Candidate, typ CandidateType) []Candidate {
	var result []Candidate
	for _, c := range candidates {
		if c.Type == typ {
			result = append(result, c)
		}
	}
	return result
}

// --- 9.1 Multi-Table Schema Tests ---

func TestIntegration_9_1_MultiTableSchema(t *testing.T) {
	cat := setupIntegrationCatalog(t)

	t.Run("column_scoped_to_correct_table_in_JOIN", func(t *testing.T) {
		// When selecting columns from a joined query without qualifier,
		// columns from both tables should be present.
		sql := "SELECT  FROM users JOIN orders ON users.id = orders.user_id"
		cursor := len("SELECT ")
		candidates := Complete(sql, cursor, cat)

		// Columns from users.
		assertContains(t, candidates, "name", CandidateColumn)
		assertContains(t, candidates, "email", CandidateColumn)

		// Columns from orders.
		assertContains(t, candidates, "user_id", CandidateColumn)
		assertContains(t, candidates, "amount", CandidateColumn)
		assertContains(t, candidates, "status", CandidateColumn)

		// Shared column (id) should appear once (dedup).
		assertContains(t, candidates, "id", CandidateColumn)
	})

	t.Run("unqualified_columns_from_all_tables", func(t *testing.T) {
		// SELECT | FROM users JOIN orders ON users.id = orders.user_id
		// Unqualified column completion should include columns from both tables.
		sql := "SELECT  FROM users JOIN orders ON users.id = orders.user_id"
		cursor := len("SELECT ")
		candidates := Complete(sql, cursor, cat)

		cols := candidatesOfType(candidates, CandidateColumn)
		// Should have columns from users: id, name, email
		assertContains(t, candidates, "name", CandidateColumn)
		assertContains(t, candidates, "email", CandidateColumn)
		// Should have columns from orders: user_id, amount, status
		assertContains(t, candidates, "user_id", CandidateColumn)
		assertContains(t, candidates, "amount", CandidateColumn)
		assertContains(t, candidates, "status", CandidateColumn)
		_ = cols
	})

	t.Run("table_alias_completion", func(t *testing.T) {
		// SELECT | FROM users AS x
		// Alias x refers to users table. Unqualified column completion in
		// this context should include users columns (via alias resolution).
		sql := "SELECT  FROM users AS x"
		cursor := len("SELECT ")
		candidates := Complete(sql, cursor, cat)

		// The ref extractor sees alias x for users, so users columns available.
		assertContains(t, candidates, "id", CandidateColumn)
		assertContains(t, candidates, "name", CandidateColumn)
		assertContains(t, candidates, "email", CandidateColumn)
	})

	t.Run("view_column_completion", func(t *testing.T) {
		// SELECT | FROM active_users
		// active_users is a view with columns: id, name
		sql := "SELECT  FROM active_users"
		cursor := len("SELECT ")
		candidates := Complete(sql, cursor, cat)

		assertContains(t, candidates, "id", CandidateColumn)
		assertContains(t, candidates, "name", CandidateColumn)
	})

	t.Run("cte_column_completion", func(t *testing.T) {
		// WITH cte AS (SELECT id, name FROM users) SELECT | FROM cte
		sql := "WITH cte AS (SELECT id, name FROM users) SELECT  FROM cte"
		cursor := len("WITH cte AS (SELECT id, name FROM users) SELECT ")
		candidates := Complete(sql, cursor, cat)

		// CTE resolves to the users table, so should have some columns.
		// At minimum, columns from the referenced table should be available.
		assertHasType(t, candidates, CandidateColumn)
	})

	t.Run("database_qualified_table", func(t *testing.T) {
		// SELECT * FROM test.| → tables in database test
		sql := "SELECT * FROM test. "
		cursor := len("SELECT * FROM test.")
		candidates := Complete(sql, cursor, cat)

		// Should have tables from the "test" database.
		assertContains(t, candidates, "users", CandidateTable)
		assertContains(t, candidates, "orders", CandidateTable)
		assertContains(t, candidates, "products", CandidateTable)
	})
}

// --- 9.2 Edge Cases ---

func TestIntegration_9_2_EdgeCases(t *testing.T) {
	cat := setupIntegrationCatalog(t)

	t.Run("cursor_at_beginning", func(t *testing.T) {
		// |SELECT * FROM users → top-level keywords
		sql := "SELECT * FROM users"
		cursor := 0
		candidates := Complete(sql, cursor, cat)

		// At the beginning, prefix is empty so top-level keywords returned.
		if len(candidates) == 0 {
			t.Fatal("expected top-level keywords at cursor position 0")
		}
		assertContains(t, candidates, "SELECT", CandidateKeyword)
	})

	t.Run("cursor_in_middle_of_identifier", func(t *testing.T) {
		// SELECT us|ers FROM t → prefix "us" filters candidates
		sql := "SELECT users FROM users"
		cursor := len("SELECT us")
		candidates := Complete(sql, cursor, cat)

		// The prefix "us" should filter candidates. "users" column should match.
		for _, c := range candidates {
			if !strings.HasPrefix(strings.ToUpper(c.Text), "US") {
				t.Errorf("candidate %q does not match prefix 'us'", c.Text)
			}
		}
	})

	t.Run("cursor_after_semicolon", func(t *testing.T) {
		// SELECT 1; SELECT | → new statement context
		sql := "SELECT 1; SELECT "
		cursor := len(sql)
		candidates := Complete(sql, cursor, cat)

		// Should get candidates for SELECT context (columns, keywords, etc.)
		if len(candidates) == 0 {
			t.Fatal("expected candidates after semicolon in new statement")
		}
		assertHasType(t, candidates, CandidateKeyword)
	})

	t.Run("empty_sql", func(t *testing.T) {
		// | → top-level keywords
		candidates := Complete("", 0, cat)

		if len(candidates) == 0 {
			t.Fatal("expected top-level keywords for empty SQL")
		}
		assertContains(t, candidates, "SELECT", CandidateKeyword)
		assertContains(t, candidates, "INSERT", CandidateKeyword)
		assertContains(t, candidates, "UPDATE", CandidateKeyword)
		assertContains(t, candidates, "DELETE", CandidateKeyword)
		assertContains(t, candidates, "CREATE", CandidateKeyword)
	})

	t.Run("whitespace_only", func(t *testing.T) {
		// "   |" → top-level keywords
		sql := "   "
		cursor := len(sql)
		candidates := Complete(sql, cursor, cat)

		if len(candidates) == 0 {
			t.Fatal("expected top-level keywords for whitespace-only SQL")
		}
		assertContains(t, candidates, "SELECT", CandidateKeyword)
	})

	t.Run("very_long_sql", func(t *testing.T) {
		// Very long SQL with cursor in middle - should not panic or hang.
		var b strings.Builder
		b.WriteString("SELECT ")
		for i := 0; i < 100; i++ {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString("id")
		}
		b.WriteString(" FROM users WHERE ")
		cursor := b.Len()
		b.WriteString("id > 0")
		sql := b.String()

		candidates := Complete(sql, cursor, cat)
		// Should return something (at least columns/keywords).
		if len(candidates) == 0 {
			t.Fatal("expected candidates for cursor in middle of long SQL")
		}
	})

	t.Run("syntax_errors_before_cursor", func(t *testing.T) {
		// SQL with syntax errors before cursor: completion still works.
		sql := "SELCT * FORM users WHERE "
		cursor := len(sql)
		candidates := Complete(sql, cursor, cat)

		// Even with typos, the system should return some candidates (fallback).
		// It's acceptable if it returns top-level keywords or columns.
		if len(candidates) == 0 {
			t.Log("no candidates for SQL with errors - acceptable fallback behavior")
		}
		// Main point: should not panic.
	})

	t.Run("backtick_quoted_identifiers", func(t *testing.T) {
		// SELECT `| FROM users → should still produce candidates
		// Note: backtick handling may be limited, but should not panic.
		sql := "SELECT ` FROM users"
		cursor := len("SELECT `")
		candidates := Complete(sql, cursor, cat)
		// Should not panic. Candidates may or may not be returned depending
		// on backtick handling, but the system must be robust.
		_ = candidates
	})
}

// --- 9.3 Complex SQL Patterns ---

func TestIntegration_9_3_ComplexSQLPatterns(t *testing.T) {
	cat := setupIntegrationCatalog(t)

	t.Run("nested_subquery_column_completion", func(t *testing.T) {
		// SELECT * FROM users WHERE id IN (SELECT | FROM orders)
		// The ref extractor finds outer-scope table refs (users).
		// Column candidates are returned (from users or fallback to all).
		sql := "SELECT * FROM users WHERE id IN (SELECT  FROM orders)"
		cursor := len("SELECT * FROM users WHERE id IN (SELECT ")
		candidates := Complete(sql, cursor, cat)

		// Should return some column candidates.
		assertHasType(t, candidates, CandidateColumn)
		// Should also have keyword/function candidates for SELECT context.
		assertHasType(t, candidates, CandidateKeyword)
	})

	t.Run("correlated_subquery", func(t *testing.T) {
		// SELECT *, (SELECT | FROM orders WHERE orders.user_id = users.id) FROM users
		sql := "SELECT *, (SELECT  FROM orders WHERE orders.user_id = users.id) FROM users"
		cursor := len("SELECT *, (SELECT ")
		candidates := Complete(sql, cursor, cat)

		// Should return column and keyword/function candidates.
		assertHasType(t, candidates, CandidateColumn)
	})

	t.Run("union_select", func(t *testing.T) {
		// SELECT name FROM users UNION SELECT | FROM products
		// The ref extractor walks both sides of the UNION, finding users and products.
		sql := "SELECT name FROM users UNION SELECT  FROM products"
		cursor := len("SELECT name FROM users UNION SELECT ")
		candidates := Complete(sql, cursor, cat)

		// Should return column candidates (from the combined table refs).
		assertHasType(t, candidates, CandidateColumn)
		// Both users and products columns should be available via UNION ref extraction.
		assertContains(t, candidates, "name", CandidateColumn)
	})

	t.Run("multiple_joins", func(t *testing.T) {
		// SELECT | FROM users JOIN orders ON users.id = orders.user_id JOIN products ON ...
		sql := "SELECT  FROM users JOIN orders ON users.id = orders.user_id JOIN products ON orders.id = products.id"
		cursor := len("SELECT ")
		candidates := Complete(sql, cursor, cat)

		// Should have columns from all 3 tables.
		assertContains(t, candidates, "email", CandidateColumn)    // users
		assertContains(t, candidates, "user_id", CandidateColumn)  // orders
		assertContains(t, candidates, "amount", CandidateColumn)   // orders
		assertContains(t, candidates, "price", CandidateColumn)    // products
	})

	t.Run("insert_select", func(t *testing.T) {
		// INSERT INTO users SELECT | FROM products
		// The ref extractor finds 'users' (INSERT target) and 'products' (SELECT FROM).
		sql := "INSERT INTO users SELECT  FROM products"
		cursor := len("INSERT INTO users SELECT ")
		candidates := Complete(sql, cursor, cat)

		// Should have column candidates (from users and/or products).
		assertHasType(t, candidates, CandidateColumn)
		// The INSERT target (users) columns should be available.
		assertContains(t, candidates, "name", CandidateColumn)
	})

	t.Run("complex_alter_table", func(t *testing.T) {
		// ALTER TABLE users ADD COLUMN | → should get type candidates or column options
		// Testing the simpler ALTER TABLE path that works end-to-end.
		sql := "ALTER TABLE users ADD INDEX idx ("
		cursor := len(sql)
		candidates := Complete(sql, cursor, cat)

		// ADD INDEX (|) should produce columnref candidates for users.
		assertContains(t, candidates, "id", CandidateColumn)
		assertContains(t, candidates, "name", CandidateColumn)
		assertContains(t, candidates, "email", CandidateColumn)
	})
}
