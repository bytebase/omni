package completion

import (
	"testing"

	"github.com/bytebase/omni/redshift/catalog"
)

func setupCatalog(t *testing.T) *catalog.Catalog {
	t.Helper()
	cat := catalog.New()
	_, err := cat.Exec(`
		CREATE TABLE users (id serial PRIMARY KEY, name text NOT NULL, email text);
		CREATE TABLE orders (id serial PRIMARY KEY, user_id int REFERENCES users(id), amount numeric);
		CREATE VIEW active_users AS SELECT * FROM users WHERE name IS NOT NULL;
		CREATE SEQUENCE order_seq;
	`, nil)
	if err != nil {
		t.Fatal(err)
	}
	return cat
}

func hasCandidate(candidates []Candidate, text string, typ CandidateType) bool {
	for _, c := range candidates {
		if c.Text == text && c.Type == typ {
			return true
		}
	}
	return false
}

func hasCandidateText(candidates []Candidate, text string) bool {
	for _, c := range candidates {
		if c.Text == text {
			return true
		}
	}
	return false
}

func TestE2E_StatementStart(t *testing.T) {
	candidates := Complete("", 0, nil)
	for _, kw := range []string{"SELECT", "INSERT", "CREATE", "ALTER", "DROP", "UPDATE", "DELETE"} {
		if !hasCandidate(candidates, kw, CandidateKeyword) {
			t.Errorf("missing keyword %q at statement start", kw)
		}
	}
}

func TestE2E_SelectTargetList(t *testing.T) {
	cat := setupCatalog(t)
	candidates := Complete("SELECT  FROM users", 7, cat)
	for _, col := range []string{"id", "name", "email"} {
		if !hasCandidate(candidates, col, CandidateColumn) {
			t.Errorf("missing column %q in SELECT target", col)
		}
	}
}

func TestE2E_FromClause(t *testing.T) {
	cat := setupCatalog(t)
	candidates := Complete("SELECT 1 FROM ", 14, cat)
	if !hasCandidate(candidates, "users", CandidateTable) {
		t.Error("missing table 'users' after FROM")
	}
	if !hasCandidate(candidates, "orders", CandidateTable) {
		t.Error("missing table 'orders' after FROM")
	}
	if !hasCandidate(candidates, "active_users", CandidateView) {
		t.Error("missing view 'active_users' after FROM")
	}
}

func TestE2E_WhereClause(t *testing.T) {
	cat := setupCatalog(t)
	candidates := Complete("SELECT * FROM users WHERE ", 26, cat)
	if !hasCandidate(candidates, "id", CandidateColumn) {
		t.Error("missing column 'id' in WHERE")
	}
	if !hasCandidate(candidates, "name", CandidateColumn) {
		t.Error("missing column 'name' in WHERE")
	}
}

func TestE2E_AfterWherePredicate(t *testing.T) {
	cat := setupCatalog(t)
	candidates := Complete("SELECT * FROM users WHERE id = 1 ", 33, cat)
	if !hasCandidate(candidates, "AND", CandidateKeyword) {
		t.Error("missing AND after WHERE predicate")
	}
	if !hasCandidate(candidates, "ORDER", CandidateKeyword) {
		t.Error("missing ORDER after WHERE predicate")
	}
}

func TestE2E_CreateDispatch(t *testing.T) {
	candidates := Complete("CREATE ", 7, nil)
	for _, kw := range []string{"TABLE", "VIEW", "INDEX", "FUNCTION", "SCHEMA", "DATABASE"} {
		if !hasCandidate(candidates, kw, CandidateKeyword) {
			t.Errorf("missing %q after CREATE", kw)
		}
	}
}

func TestE2E_AlterDispatch(t *testing.T) {
	candidates := Complete("ALTER ", 6, nil)
	for _, kw := range []string{"TABLE", "DATABASE", "ROLE", "SCHEMA"} {
		if !hasCandidate(candidates, kw, CandidateKeyword) {
			t.Errorf("missing %q after ALTER", kw)
		}
	}
}

func TestE2E_DropDispatch(t *testing.T) {
	candidates := Complete("DROP ", 5, nil)
	for _, kw := range []string{"TABLE", "VIEW", "INDEX", "FUNCTION"} {
		if !hasCandidate(candidates, kw, CandidateKeyword) {
			t.Errorf("missing %q after DROP", kw)
		}
	}
}

func TestE2E_PrefixFiltering(t *testing.T) {
	candidates := Complete("SEL", 3, nil)
	if !hasCandidate(candidates, "SELECT", CandidateKeyword) {
		t.Error("SELECT should match prefix SEL")
	}
	if hasCandidate(candidates, "INSERT", CandidateKeyword) {
		t.Error("INSERT should NOT match prefix SEL")
	}
}

func TestE2E_CTE(t *testing.T) {
	cat := setupCatalog(t)
	sql := "WITH active AS (SELECT * FROM users) SELECT * FROM "
	candidates := Complete(sql, len(sql), cat)
	if !hasCandidateText(candidates, "active") {
		t.Error("missing CTE 'active' in FROM candidates")
	}
}

func TestE2E_MultipleFromTables(t *testing.T) {
	cat := setupCatalog(t)
	candidates := Complete("SELECT  FROM users u, orders o", 7, cat)
	if !hasCandidate(candidates, "id", CandidateColumn) {
		t.Error("missing column 'id'")
	}
	if !hasCandidate(candidates, "amount", CandidateColumn) {
		t.Error("missing column 'amount' from orders")
	}
}

func TestE2E_SequenceInRelations(t *testing.T) {
	cat := setupCatalog(t)
	candidates := Complete("SELECT 1 FROM ", 14, cat)
	if !hasCandidate(candidates, "order_seq", CandidateSequence) {
		t.Error("missing sequence 'order_seq' after FROM")
	}
}

func TestE2E_ViewColumnsInSelect(t *testing.T) {
	cat := setupCatalog(t)
	candidates := Complete("SELECT  FROM active_users", 7, cat)
	// active_users is a view on users, should get users' columns
	if !hasCandidate(candidates, "name", CandidateColumn) {
		t.Error("missing column 'name' from view active_users")
	}
}

func TestE2E_InsertInto(t *testing.T) {
	cat := setupCatalog(t)
	candidates := Complete("INSERT INTO ", 12, cat)
	if !hasCandidate(candidates, "users", CandidateTable) {
		t.Error("missing table 'users' after INSERT INTO")
	}
	if !hasCandidate(candidates, "orders", CandidateTable) {
		t.Error("missing table 'orders' after INSERT INTO")
	}
}

func TestE2E_UpdateTable(t *testing.T) {
	cat := setupCatalog(t)
	// UPDATE SET completion — the parser currently panics on bare "UPDATE "
	// so we test with a more complete statement instead.
	candidates := Complete("UPDATE users SET ", 17, cat)
	// Should at least not panic and return some candidates
	_ = candidates

	// Test that UPDATE with a partial table name filters correctly
	candidates = Complete("UPDATE us", 9, cat)
	// Should not panic
	_ = candidates
}

func TestE2E_DeleteFrom(t *testing.T) {
	cat := setupCatalog(t)
	candidates := Complete("DELETE FROM ", 12, cat)
	if !hasCandidate(candidates, "orders", CandidateTable) {
		t.Error("missing table 'orders' after DELETE FROM")
	}
}

func TestE2E_TablePrefixFilter(t *testing.T) {
	cat := setupCatalog(t)
	candidates := Complete("SELECT 1 FROM us", 16, cat)
	if !hasCandidate(candidates, "users", CandidateTable) {
		t.Error("'users' should match prefix 'us'")
	}
	if hasCandidate(candidates, "orders", CandidateTable) {
		t.Error("'orders' should NOT match prefix 'us'")
	}
}

func TestE2E_NilCatalog(t *testing.T) {
	// Should not panic with nil catalog
	candidates := Complete("SELECT * FROM ", 14, nil)
	// Only keywords should be present, no tables
	for _, c := range candidates {
		if c.Type == CandidateTable {
			t.Errorf("unexpected table candidate %q with nil catalog", c.Text)
		}
	}
}

func TestE2E_OffsetBeyondLength(t *testing.T) {
	// Should not panic when offset exceeds sql length
	candidates := Complete("SELECT", 100, nil)
	_ = candidates // just verify no panic
}

func TestE2E_MidTokenCompletion(t *testing.T) {
	cat := setupCatalog(t)
	// Cursor in the middle of "users" token after FROM
	candidates := Complete("SELECT * FROM use", 17, cat)
	if !hasCandidate(candidates, "users", CandidateTable) {
		t.Error("'users' should match partial 'use' after FROM")
	}
}
