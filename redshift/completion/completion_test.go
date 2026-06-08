package completion

import (
	"testing"

	"github.com/bytebase/omni/redshift/catalog"
	"github.com/bytebase/omni/redshift/parser"
)

func TestCompleteCandidateTypes(t *testing.T) {
	// Verify all candidate types are distinct.
	types := []CandidateType{
		CandidateKeyword, CandidateSchema, CandidateTable, CandidateView,
		CandidateMaterializedView, CandidateColumn, CandidateFunction,
		CandidateSequence, CandidateIndex, CandidateType_, CandidateTrigger,
		CandidatePolicy, CandidateExtension,
	}
	seen := make(map[CandidateType]bool)
	for _, ct := range types {
		if seen[ct] {
			t.Errorf("duplicate CandidateType value: %d", ct)
		}
		seen[ct] = true
	}
	if len(seen) != 13 {
		t.Errorf("expected 13 distinct types, got %d", len(seen))
	}
}

func TestCompleteEmpty(t *testing.T) {
	// Complete with empty input and no catalog should not panic.
	result := Complete("", 0, nil)
	// Result may be nil or empty for now; just verify no panic.
	_ = result
}

func TestCompleteTableFromCatalog(t *testing.T) {
	cat := catalog.New()
	cat.Exec("CREATE TABLE users (id int, name text);", nil)

	candidates := Complete("SELECT * FROM ", 14, cat)
	found := false
	for _, c := range candidates {
		if c.Type == CandidateTable && c.Text == "users" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'users' table candidate")
	}
}

func TestCompleteKeywords(t *testing.T) {
	candidates := Complete("", 0, nil)
	found := false
	for _, c := range candidates {
		if c.Type == CandidateKeyword && c.Text == "SELECT" {
			found = true
		}
	}
	if !found {
		t.Error("expected SELECT keyword candidate")
	}
}

func TestCompletePrefixFilter(t *testing.T) {
	candidates := Complete("SEL", 3, nil)
	for _, c := range candidates {
		if c.Type == CandidateKeyword && c.Text == "INSERT" {
			t.Error("INSERT should be filtered out by prefix 'SEL'")
		}
	}
	found := false
	for _, c := range candidates {
		if c.Type == CandidateKeyword && c.Text == "SELECT" {
			found = true
		}
	}
	if !found {
		t.Error("expected SELECT to match prefix 'SEL'")
	}
}

func TestCompleteCTE(t *testing.T) {
	cat := catalog.New()
	cat.Exec("CREATE TABLE users (id int, name text);", nil)

	sql := "WITH active AS (SELECT * FROM users) SELECT * FROM "
	offset := len(sql)
	candidates := Complete(sql, offset, cat)

	// Should include CTE name "active" as a table ref
	found := false
	for _, c := range candidates {
		if c.Text == "active" {
			found = true
		}
	}
	if !found {
		t.Error("expected CTE 'active' in table candidates")
	}
}

func TestTrickyCompletePartialSQL(t *testing.T) {
	cat := catalog.New()
	cat.Exec("CREATE TABLE orders (id int, amount numeric);", nil)

	// "SELECT * FROM " — ends abruptly, standard may or may not work
	// but tricky should handle it
	candidates := Complete("SELECT * FROM ", 14, cat)
	found := false
	for _, c := range candidates {
		if c.Type == CandidateTable && c.Text == "orders" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'orders' table in tricky completion")
	}
}

// ---------------------------------------------------------------------------
// Batch 1: Parenthesized column-list completion
// ---------------------------------------------------------------------------

func TestCompleteJoinUsing(t *testing.T) {
	cat := catalog.New()
	cat.Exec("CREATE TABLE t1 (id int, name text); CREATE TABLE t2 (id int, age int);", nil)

	sql := "SELECT * FROM t1 JOIN t2 USING ("
	candidates := Complete(sql, len(sql), cat)
	found := false
	for _, c := range candidates {
		if c.Type == CandidateColumn {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected column candidates inside JOIN USING ()")
	}
}

func TestCollectJoinUsingEmitsColumnref(t *testing.T) {
	sql := "SELECT * FROM t1 JOIN t2 USING ("
	cs := parser.Collect(sql, len(sql))
	if cs == nil || !cs.HasRule("columnref") {
		t.Error("expected columnref rule inside JOIN USING ()")
	}
}

func TestCollectForeignKeyColumnList(t *testing.T) {
	sql := "CREATE TABLE t (a int, FOREIGN KEY ("
	cs := parser.Collect(sql, len(sql))
	if cs == nil || !cs.HasRule("columnref") {
		t.Error("expected columnref rule inside FOREIGN KEY ()")
	}
}

func TestCollectUniqueColumnList(t *testing.T) {
	sql := "CREATE TABLE t (a int, UNIQUE ("
	cs := parser.Collect(sql, len(sql))
	if cs == nil || !cs.HasRule("columnref") {
		t.Error("expected columnref rule inside UNIQUE ()")
	}
}

func TestCollectPrimaryKeyColumnList(t *testing.T) {
	sql := "CREATE TABLE t (a int, PRIMARY KEY ("
	cs := parser.Collect(sql, len(sql))
	if cs == nil || !cs.HasRule("columnref") {
		t.Error("expected columnref rule inside PRIMARY KEY ()")
	}
}

func TestCollectReferencesColumnList(t *testing.T) {
	// Column-level constraint: col REFERENCES other_table (
	sql := "CREATE TABLE t (a int REFERENCES other ("
	cs := parser.Collect(sql, len(sql))
	if cs == nil || !cs.HasRule("columnref") {
		t.Error("expected columnref rule inside REFERENCES table ()")
	}
}

func TestCollectReferencesInTableConstraint(t *testing.T) {
	// Table-level constraint: FOREIGN KEY (a) REFERENCES other (
	sql := "CREATE TABLE t (a int, FOREIGN KEY (a) REFERENCES other ("
	cs := parser.Collect(sql, len(sql))
	if cs == nil || !cs.HasRule("columnref") {
		t.Error("expected columnref rule inside table constraint REFERENCES ()")
	}
}

func TestCollectCreateIndexColumnList(t *testing.T) {
	sql := "CREATE INDEX idx ON t ("
	cs := parser.Collect(sql, len(sql))
	if cs == nil || !cs.HasRule("columnref") {
		t.Error("expected columnref rule inside CREATE INDEX ()")
	}
}

func TestCompleteCreateIndexColumn(t *testing.T) {
	cat := catalog.New()
	cat.Exec("CREATE TABLE orders (id int, amount numeric);", nil)

	sql := "CREATE INDEX idx ON orders ("
	candidates := Complete(sql, len(sql), cat)
	found := false
	for _, c := range candidates {
		if c.Type == CandidateColumn {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected column candidates inside CREATE INDEX ()")
	}
}

// ---------------------------------------------------------------------------
// Batch 2: DDL/Utility statement completion
// ---------------------------------------------------------------------------

func TestCollectAlterTableRename(t *testing.T) {
	// After RENAME, should suggest COLUMN, CONSTRAINT, TO, and columnref
	sql := "ALTER TABLE t RENAME "
	cs := parser.Collect(sql, len(sql))
	if cs == nil {
		t.Fatal("expected non-nil candidate set")
	}
	if !cs.HasRule("columnref") {
		t.Error("expected columnref rule after ALTER TABLE RENAME")
	}
	if !cs.HasToken(parser.COLUMN) {
		t.Error("expected COLUMN token after ALTER TABLE RENAME")
	}
	if !cs.HasToken(parser.TO) {
		t.Error("expected TO token after ALTER TABLE RENAME")
	}
	if !cs.HasToken(parser.CONSTRAINT) {
		t.Error("expected CONSTRAINT token after ALTER TABLE RENAME")
	}
}

func TestCollectAlterTableRenameColumn(t *testing.T) {
	// After RENAME COLUMN, should suggest columnref
	sql := "ALTER TABLE t RENAME COLUMN "
	cs := parser.Collect(sql, len(sql))
	if cs == nil || !cs.HasRule("columnref") {
		t.Error("expected columnref rule after ALTER TABLE RENAME COLUMN")
	}
}

func TestCompleteAlterTableRenameColumn(t *testing.T) {
	cat := catalog.New()
	cat.Exec("CREATE TABLE users (id int, name text);", nil)

	sql := "ALTER TABLE users RENAME COLUMN "
	candidates := Complete(sql, len(sql), cat)
	found := false
	for _, c := range candidates {
		if c.Type == CandidateColumn {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected column candidates after ALTER TABLE RENAME COLUMN")
	}
}

func TestCollectCommentOnColumn(t *testing.T) {
	// After COMMENT ON COLUMN, should suggest columnref/qualified_name
	sql := "COMMENT ON COLUMN "
	cs := parser.Collect(sql, len(sql))
	if cs == nil {
		t.Fatal("expected non-nil candidate set")
	}
	if !cs.HasRule("columnref") && !cs.HasRule("qualified_name") {
		t.Error("expected columnref or qualified_name rule after COMMENT ON COLUMN")
	}
}

func TestCollectCommentOn(t *testing.T) {
	// After COMMENT ON, should suggest object type keywords
	sql := "COMMENT ON "
	cs := parser.Collect(sql, len(sql))
	if cs == nil {
		t.Fatal("expected non-nil candidate set")
	}
	if !cs.HasToken(parser.COLUMN) {
		t.Error("expected COLUMN token after COMMENT ON")
	}
	if !cs.HasToken(parser.TABLE) {
		t.Error("expected TABLE token after COMMENT ON")
	}
}

func TestCompleteCommentOnColumn(t *testing.T) {
	cat := catalog.New()
	cat.Exec("CREATE TABLE users (id int, name text);", nil)

	sql := "COMMENT ON COLUMN users."
	candidates := Complete(sql, len(sql), cat)
	found := false
	for _, c := range candidates {
		if c.Type == CandidateColumn {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected column candidates after COMMENT ON COLUMN users.")
	}
}

func TestCollectGrantInSchema(t *testing.T) {
	// After GRANT ALL TABLES IN SCHEMA, should suggest schema name
	sql := "GRANT SELECT ON ALL TABLES IN SCHEMA "
	cs := parser.Collect(sql, len(sql))
	if cs == nil || !cs.HasRule("qualified_name") {
		t.Error("expected qualified_name rule after GRANT ... IN SCHEMA")
	}
}

// ---------------------------------------------------------------------------
// Batch 3: Parenthesized subquery completion
// ---------------------------------------------------------------------------

func TestCollectParenthesizedSubquery(t *testing.T) {
	sql := "SELECT * FROM ("
	cs := parser.Collect(sql, len(sql))
	if cs == nil {
		t.Fatal("expected non-nil candidate set")
	}
	if !cs.HasToken(parser.SELECT) {
		t.Error("expected SELECT token inside parenthesized subquery")
	}
}

func TestCompleteParenthesizedSubquery(t *testing.T) {
	candidates := Complete("SELECT * FROM (", 15, nil)
	found := false
	for _, c := range candidates {
		if c.Type == CandidateKeyword && c.Text == "SELECT" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected SELECT keyword inside parenthesized subquery")
	}
}

// ---------------------------------------------------------------------------
// Batch 4: Bytebase-reported remaining issues + systematic audit fixes
// ---------------------------------------------------------------------------

func TestCollectDropColumn(t *testing.T) {
	sql := "ALTER TABLE t1 DROP COLUMN "
	cs := parser.Collect(sql, len(sql))
	if cs == nil || !cs.HasRule("columnref") {
		t.Error("expected columnref rule after ALTER TABLE DROP COLUMN")
	}
}

func TestCompleteDropColumn(t *testing.T) {
	cat := catalog.New()
	cat.Exec("CREATE TABLE t1 (c1 int, c2 text);", nil)

	sql := "ALTER TABLE t1 DROP COLUMN "
	candidates := Complete(sql, len(sql), cat)
	found := false
	for _, c := range candidates {
		if c.Type == CandidateColumn {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected column candidates after ALTER TABLE DROP COLUMN")
	}
}

func TestCollectAlterColumn(t *testing.T) {
	sql := "ALTER TABLE t1 ALTER COLUMN "
	cs := parser.Collect(sql, len(sql))
	if cs == nil || !cs.HasRule("columnref") {
		t.Error("expected columnref rule after ALTER TABLE ALTER COLUMN")
	}
}

func TestCompleteAlterColumn(t *testing.T) {
	cat := catalog.New()
	cat.Exec("CREATE TABLE t1 (c1 int, c2 text);", nil)

	sql := "ALTER TABLE t1 ALTER COLUMN "
	candidates := Complete(sql, len(sql), cat)
	found := false
	for _, c := range candidates {
		if c.Type == CandidateColumn {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected column candidates after ALTER TABLE ALTER COLUMN")
	}
}

func TestCollectParenthesizedUnion(t *testing.T) {
	sql := "(SELECT c1 FROM t1) UNION (SELECT  FROM t2)"
	cs := parser.Collect(sql, 34) // cursor after "(SELECT "
	if cs == nil || !cs.HasRule("columnref") {
		t.Error("expected columnref rule in parenthesized UNION branch")
	}
}

func TestCompleteParenthesizedUnion(t *testing.T) {
	cat := catalog.New()
	cat.Exec("CREATE TABLE t1 (c1 int); CREATE TABLE t2 (c1 int, c2 int);", nil)

	sql := "(SELECT c1 FROM t1) UNION (SELECT  FROM t2)"
	candidates := Complete(sql, 34, cat)
	found := false
	for _, c := range candidates {
		if c.Type == CandidateColumn {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected column candidates in parenthesized UNION branch")
	}
}

func TestCollectDropConstraint(t *testing.T) {
	sql := "ALTER TABLE t1 DROP CONSTRAINT "
	cs := parser.Collect(sql, len(sql))
	if cs == nil || !cs.HasRule("qualified_name") {
		t.Error("expected qualified_name rule after ALTER TABLE DROP CONSTRAINT")
	}
}

func TestCollectRenameConstraint(t *testing.T) {
	sql := "ALTER TABLE t1 RENAME CONSTRAINT "
	cs := parser.Collect(sql, len(sql))
	if cs == nil || !cs.HasRule("qualified_name") {
		t.Error("expected qualified_name rule after ALTER TABLE RENAME CONSTRAINT")
	}
}

func TestCollectAlterConstraint(t *testing.T) {
	sql := "ALTER TABLE t1 ALTER CONSTRAINT "
	cs := parser.Collect(sql, len(sql))
	if cs == nil || !cs.HasRule("qualified_name") {
		t.Error("expected qualified_name rule after ALTER TABLE ALTER CONSTRAINT")
	}
}

func TestCollectSecurityLabelOnColumn(t *testing.T) {
	sql := "SECURITY LABEL ON COLUMN "
	cs := parser.Collect(sql, len(sql))
	if cs == nil {
		t.Fatal("expected non-nil candidate set")
	}
	if !cs.HasRule("columnref") && !cs.HasRule("qualified_name") {
		t.Error("expected columnref or qualified_name rule after SECURITY LABEL ON COLUMN")
	}
}

func TestCollectUniqueUsingIndex(t *testing.T) {
	sql := "CREATE TABLE t (a int, UNIQUE USING INDEX "
	cs := parser.Collect(sql, len(sql))
	if cs == nil || !cs.HasRule("qualified_name") {
		t.Error("expected qualified_name rule after UNIQUE USING INDEX")
	}
}

func TestCollectPrimaryKeyUsingIndex(t *testing.T) {
	sql := "CREATE TABLE t (a int, PRIMARY KEY USING INDEX "
	cs := parser.Collect(sql, len(sql))
	if cs == nil || !cs.HasRule("qualified_name") {
		t.Error("expected qualified_name rule after PRIMARY KEY USING INDEX")
	}
}
