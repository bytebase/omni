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

func TestCollectNestedFromSubqueryAfterSelect(t *testing.T) {
	prefix := "SELECT * FROM (SELECT * FROM (SELECT "
	sql := prefix + " FROM t1) a) b"
	candidates := Collect(sql, len(prefix))
	if candidates == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !candidates.HasRule("columnref") {
		t.Error("expected columnref rule candidate in nested SELECT target list")
	}
	if !candidates.HasRule("func_name") {
		t.Error("expected func_name rule candidate in nested SELECT target list")
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

func TestCollectUpdateRelation(t *testing.T) {
	// "UPDATE " — should offer relation_expr
	cs := Collect("UPDATE ", 7)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !cs.HasRule("relation_expr") {
		t.Error("expected relation_expr rule after UPDATE")
	}
}

func TestCollectUpdateSet(t *testing.T) {
	// "UPDATE t1 SET " — should offer columnref
	cs := Collect("UPDATE t1 SET ", 14)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !cs.HasRule("columnref") {
		t.Error("expected columnref rule in SET clause")
	}
}

func TestCollectUpdateAfterSet(t *testing.T) {
	// "UPDATE t1 SET a = 1 " — should offer FROM, WHERE, RETURNING
	cs := Collect("UPDATE t1 SET a = 1 ", 20)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	for _, tok := range []int{FROM, WHERE, RETURNING} {
		if !cs.HasToken(tok) {
			t.Errorf("missing token after SET clause: %d (%s)", tok, TokenName(tok))
		}
	}
}

func TestCollectDeleteFrom(t *testing.T) {
	// "DELETE FROM " — should offer relation_expr
	cs := Collect("DELETE FROM ", 12)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !cs.HasRule("relation_expr") {
		t.Error("expected relation_expr rule after DELETE FROM")
	}
}

func TestCollectDeleteAfterRelation(t *testing.T) {
	// "DELETE FROM t1 " — should offer USING, WHERE, RETURNING
	cs := Collect("DELETE FROM t1 ", 15)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	for _, tok := range []int{USING, WHERE, RETURNING} {
		if !cs.HasToken(tok) {
			t.Errorf("missing token after DELETE relation: %d (%s)", tok, TokenName(tok))
		}
	}
}

func TestCollectInsertInto(t *testing.T) {
	// "INSERT INTO " — should offer qualified_name (via parseQualifiedName)
	cs := Collect("INSERT INTO ", 12)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !cs.HasRule("qualified_name") {
		t.Error("expected qualified_name rule after INSERT INTO")
	}
}

func TestCollectInsertColumns(t *testing.T) {
	// "INSERT INTO t1 (" — should offer columnref
	cs := Collect("INSERT INTO t1 (", 16)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !cs.HasRule("columnref") {
		t.Error("expected columnref rule in INSERT column list")
	}
}

func TestCollectInsertAfterTable(t *testing.T) {
	// "INSERT INTO t1 " — should offer (, DEFAULT, VALUES, SELECT etc.
	cs := Collect("INSERT INTO t1 ", 15)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	for _, tok := range []int{'(', DEFAULT, VALUES, SELECT} {
		if !cs.HasToken(tok) {
			t.Errorf("missing token after INSERT INTO table: %d", tok)
		}
	}
}

func TestCollectAlterTableCmd(t *testing.T) {
	// "ALTER TABLE t1 " — should offer ADD, DROP, ALTER, SET, etc.
	cs := Collect("ALTER TABLE t1 ", 15)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	for _, tok := range []int{ADD_P, DROP, ALTER, SET, OWNER, ENABLE_P, DISABLE_P, RENAME} {
		if !cs.HasToken(tok) {
			t.Errorf("missing ALTER TABLE sub-command: %d (%s)", tok, TokenName(tok))
		}
	}
}

func TestCollectAlterTableAdd(t *testing.T) {
	// "ALTER TABLE t1 ADD " — should offer COLUMN, CONSTRAINT, etc.
	cs := Collect("ALTER TABLE t1 ADD ", 19)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	for _, tok := range []int{COLUMN, CONSTRAINT, CHECK, UNIQUE, PRIMARY, FOREIGN} {
		if !cs.HasToken(tok) {
			t.Errorf("missing ALTER TABLE ADD sub-keyword: %d (%s)", tok, TokenName(tok))
		}
	}
}

func TestCollectAlterTableDrop(t *testing.T) {
	// "ALTER TABLE t1 DROP " — should offer COLUMN, CONSTRAINT
	cs := Collect("ALTER TABLE t1 DROP ", 20)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	for _, tok := range []int{COLUMN, CONSTRAINT} {
		if !cs.HasToken(tok) {
			t.Errorf("missing ALTER TABLE DROP sub-keyword: %d (%s)", tok, TokenName(tok))
		}
	}
	if !cs.HasRule("columnref") {
		t.Error("expected columnref rule in ALTER TABLE DROP (direct column name)")
	}
}

func TestCollectAlterTableAlter(t *testing.T) {
	// "ALTER TABLE t1 ALTER " — should offer COLUMN, CONSTRAINT, columnref
	cs := Collect("ALTER TABLE t1 ALTER ", 21)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !cs.HasToken(COLUMN) {
		t.Error("expected COLUMN after ALTER TABLE ALTER")
	}
	if !cs.HasToken(CONSTRAINT) {
		t.Error("expected CONSTRAINT after ALTER TABLE ALTER")
	}
	if !cs.HasRule("columnref") {
		t.Error("expected columnref rule after ALTER TABLE ALTER")
	}
}

func TestCollectMergeInto(t *testing.T) {
	// "MERGE INTO " — should offer relation_expr and qualified_name
	cs := Collect("MERGE INTO ", 11)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !cs.HasRule("relation_expr") {
		t.Error("expected relation_expr rule after MERGE INTO")
	}
	if !cs.HasRule("qualified_name") {
		t.Error("expected qualified_name rule after MERGE INTO")
	}
}

func TestCollectQualifiedNameRules(t *testing.T) {
	// "SELECT 1 FROM " — should offer both relation_expr and qualified_name
	cs := Collect("SELECT 1 FROM ", 14)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !cs.HasRule("relation_expr") {
		t.Error("expected relation_expr rule after FROM")
	}
	if !cs.HasRule("qualified_name") {
		t.Error("expected qualified_name rule after FROM")
	}
}

func TestCollectCTEPositions(t *testing.T) {
	// "WITH cte AS (SELECT 1) SELECT " — CTE position should be recorded
	sql := "WITH cte AS (SELECT 1) SELECT "
	cs := Collect(sql, len(sql))
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if len(cs.CTEPositions) == 0 {
		t.Error("expected CTE position to be recorded")
	}
	if len(cs.CTEPositions) > 0 && cs.CTEPositions[0] != 0 {
		t.Errorf("expected CTE position 0, got %d", cs.CTEPositions[0])
	}
}

func TestCollectSelectAliasPositions(t *testing.T) {
	// "SELECT a AS alias1, b alias2 FROM " — alias positions should be recorded
	sql := "SELECT a AS alias1, b alias2 FROM "
	cs := Collect(sql, len(sql))
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if len(cs.SelectAliasPositions) != 2 {
		t.Errorf("expected 2 alias positions, got %d", len(cs.SelectAliasPositions))
	}
}

// TestCollectMaterializedViewIcebergClauses pins the completion a cursor after
// USING ICEBERG offers: the Iceberg form's own clauses, two of which are
// keywords and one of which (PARTITIONED BY) is text no token can name.
func TestCollectMaterializedViewIcebergClauses(t *testing.T) {
	sql := "CREATE MATERIALIZED VIEW m USING ICEBERG "
	cs := Collect(sql, len(sql))
	if cs == nil {
		t.Fatal("Collect returned nil")
	}
	for _, tok := range []int{LOCATION, TABLE} {
		if !cs.HasToken(tok) {
			t.Errorf("expected token %d in candidates", tok)
		}
	}
	if !cs.HasWord("PARTITIONED BY") {
		t.Errorf("expected PARTITIONED BY in candidate words, got %v", cs.Words)
	}
}

// TestCandidateSetWords pins the word candidates' dedup and lookup.
func TestCandidateSetWords(t *testing.T) {
	cs := newCandidateSet()
	cs.addWord("PARTITIONED BY")
	cs.addWord("PARTITIONED BY")
	if len(cs.Words) != 1 {
		t.Fatalf("words = %v, want one", cs.Words)
	}
	if !cs.HasWord("PARTITIONED BY") || cs.HasWord("LOCATION") {
		t.Fatalf("HasWord = %v, want only PARTITIONED BY", cs.Words)
	}
}

// TestCollectHierarchicalAndIcebergContinuations pins the second words the new
// clauses take, which the peek-only FIRST-set predicates can no longer reach
// once the cursor sits after the first one.
func TestCollectHierarchicalAndIcebergContinuations(t *testing.T) {
	for _, tc := range []struct {
		sql       string
		tokens    []int
		wantWords []string
	}{
		// After the select list, both entries of the clause are offered.
		{"SELECT * ", []int{START}, []string{"CONNECT BY"}},
		// After START only its WITH can follow.
		{"SELECT * START ", []int{WITH}, nil},
		// After USING ICEBERG TABLE only PROPERTIES can follow.
		{"CREATE MATERIALIZED VIEW m USING ICEBERG TABLE ", nil, []string{"PROPERTIES"}},
		// After a parsed LOCATION, the clauses after it are still offered.
		{"CREATE MATERIALIZED VIEW m USING ICEBERG LOCATION 's3://b/' ", []int{TABLE}, []string{"PARTITIONED BY"}},
		// AUTO REFRESH YES takes INTERVAL.
		{"ALTER MATERIALIZED VIEW m AUTO REFRESH YES ", []int{INTERVAL}, nil},
		// VACUUM SORT takes ONLY.
		{"VACUUM SORT ", []int{ONLY}, nil},
		// A CREATE GROUP clause repeats only through WITH.
		{"CREATE GROUP g WITH USER u, ", []int{WITH}, nil},
		// DROP DATABASE takes FORCE, or PostgreSQL's option list.
		{"DROP DATABASE db ", []int{FORCE, '('}, nil},
	} {
		cs := Collect(tc.sql, len(tc.sql))
		if cs == nil {
			t.Fatalf("Collect(%q) returned nil", tc.sql)
		}
		for _, tok := range tc.tokens {
			if !cs.HasToken(tok) {
				t.Errorf("Collect(%q): expected token %d in candidates %v", tc.sql, tok, cs.Tokens)
			}
		}
		for _, word := range tc.wantWords {
			if !cs.HasWord(word) {
				t.Errorf("Collect(%q): expected word %q in candidates %v", tc.sql, word, cs.Words)
			}
		}
	}
}

// TestCollectKeepsAliasBeforeAClauseCursor pins cursorAfterCur to the position
// right after its token: START followed by another token is the table's alias,
// and completion reaches what follows it as usual.
func TestCollectKeepsAliasBeforeAClauseCursor(t *testing.T) {
	sql := "SELECT * FROM t START WHERE "
	cs := Collect(sql, len(sql))
	if cs == nil {
		t.Fatal("Collect returned nil")
	}
	if cs.HasToken(WITH) {
		t.Errorf("WITH offered after START WHERE, which keeps START an alias: %v", cs.Tokens)
	}
	if !cs.HasToken(TRUE_P) && !cs.HasToken(IDENT) {
		t.Errorf("expected the WHERE expression's candidates, got %v", cs.Tokens)
	}
}

// TestCollectRoleListAfterComma pins the role names a comma still takes, in an
// ordinary role list and in CREATE GROUP's, where the WITH clause repeats too.
func TestCollectRoleListAfterComma(t *testing.T) {
	for _, sql := range []string{
		"CREATE ROLE r ROLE u, ",
		"ALTER GROUP g ADD USER u, ",
	} {
		cs := Collect(sql, len(sql))
		if cs == nil {
			t.Fatalf("Collect(%q) returned nil", sql)
		}
		if len(cs.Tokens) == 0 {
			t.Errorf("Collect(%q) offered no candidate after the comma", sql)
		}
	}

	cs := Collect("CREATE GROUP g WITH USER u, ", len("CREATE GROUP g WITH USER u, "))
	if cs == nil {
		t.Fatal("Collect returned nil")
	}
	if !cs.HasToken(WITH) {
		t.Errorf("WITH is not offered after a CREATE GROUP comma: %v", cs.Tokens)
	}
}

// TestCollectIdentityAndConnectByContinuations pins the two positions where the
// production goes on with something no single token leads: IDENTITY's '(' opens
// both PostgreSQL's sequence options and Redshift's seed and step, and a START
// WITH condition goes on with CONNECT BY.
func TestCollectIdentityAndConnectByContinuations(t *testing.T) {
	cs := Collect("CREATE TABLE a (id INT GENERATED BY DEFAULT AS IDENTITY (", len("CREATE TABLE a (id INT GENERATED BY DEFAULT AS IDENTITY ("))
	if cs == nil {
		t.Fatal("Collect returned nil")
	}
	for _, tok := range []int{START, INCREMENT, CACHE, MAXVALUE, MINVALUE} {
		if !cs.HasToken(tok) {
			t.Errorf("IDENTITY's '(' offers no %s: %v", TokenName(tok), cs.Tokens)
		}
	}

	sql := "SELECT * FROM t START WITH id "
	cs = Collect(sql, len(sql))
	if cs == nil {
		t.Fatal("Collect returned nil")
	}
	if !cs.HasWord("CONNECT BY") {
		t.Errorf("START WITH offers no CONNECT BY: tokens=%v words=%v", cs.Tokens, cs.Words)
	}
}

// TestCollectCreateGroupSecondWithList pins the CREATE GROUP context around
// every repeated WITH list: a second block's comma still offers WITH.
func TestCollectCreateGroupSecondWithList(t *testing.T) {
	sql := "CREATE GROUP g WITH USER u, WITH USER v, "
	cs := Collect(sql, len(sql))
	if cs == nil {
		t.Fatal("Collect returned nil")
	}
	if !cs.HasToken(WITH) {
		t.Errorf("WITH is not offered after a repeated CREATE GROUP comma: %v", cs.Tokens)
	}
}
