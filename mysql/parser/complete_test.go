package parser

import (
	"testing"
)

func TestCollect_1_2_EmptyInput(t *testing.T) {
	cs := Collect("", 0)
	if cs == nil {
		t.Fatal("Collect returned nil")
	}
	// Should have top-level statement keywords.
	want := []int{kwSELECT, kwINSERT, kwUPDATE, kwDELETE, kwCREATE, kwALTER, kwDROP}
	for _, tok := range want {
		if !cs.HasToken(tok) {
			t.Errorf("missing expected token %d", tok)
		}
	}
}

func TestCollect_1_2_AfterSemicolon(t *testing.T) {
	sql := "SELECT 1; "
	cs := Collect(sql, len(sql))
	if cs == nil {
		t.Fatal("Collect returned nil")
	}
	// After semicolon, should have top-level statement keywords.
	want := []int{kwSELECT, kwINSERT, kwUPDATE, kwDELETE, kwCREATE, kwALTER, kwDROP}
	for _, tok := range want {
		if !cs.HasToken(tok) {
			t.Errorf("missing expected token %d for new statement after semicolon", tok)
		}
	}
}

func TestCollect_1_2_SelectCursor(t *testing.T) {
	sql := "SELECT "
	cs := Collect(sql, len(sql))
	if cs == nil {
		t.Fatal("Collect returned nil")
	}
	// SELECT | should offer DISTINCT, ALL keywords and columnref, func_name rules.
	if !cs.HasToken(kwDISTINCT) {
		t.Error("missing DISTINCT keyword candidate")
	}
	if !cs.HasToken(kwALL) {
		t.Error("missing ALL keyword candidate")
	}
	if !cs.HasRule("columnref") {
		t.Error("missing columnref rule candidate")
	}
	if !cs.HasRule("func_name") {
		t.Error("missing func_name rule candidate")
	}
}

func TestCollect_1_2_CreateCursor(t *testing.T) {
	sql := "CREATE "
	cs := Collect(sql, len(sql))
	if cs == nil {
		t.Fatal("Collect returned nil")
	}
	want := []int{kwTABLE, kwINDEX, kwVIEW, kwDATABASE, kwFUNCTION, kwPROCEDURE, kwTRIGGER, kwEVENT}
	for _, tok := range want {
		if !cs.HasToken(tok) {
			t.Errorf("missing expected token %d after CREATE", tok)
		}
	}
}

func TestCollect_1_2_AlterCursor(t *testing.T) {
	sql := "ALTER "
	cs := Collect(sql, len(sql))
	if cs == nil {
		t.Fatal("Collect returned nil")
	}
	want := []int{kwTABLE, kwDATABASE, kwVIEW, kwFUNCTION, kwPROCEDURE, kwEVENT}
	for _, tok := range want {
		if !cs.HasToken(tok) {
			t.Errorf("missing expected token %d after ALTER", tok)
		}
	}
}

func TestCollect_1_2_DropCursor(t *testing.T) {
	sql := "DROP "
	cs := Collect(sql, len(sql))
	if cs == nil {
		t.Fatal("Collect returned nil")
	}
	want := []int{kwTABLE, kwINDEX, kwVIEW, kwDATABASE, kwFUNCTION, kwPROCEDURE, kwTRIGGER, kwEVENT, kwIF}
	for _, tok := range want {
		if !cs.HasToken(tok) {
			t.Errorf("missing expected token %d after DROP", tok)
		}
	}
}

// --- SELECT clause positions ---

func TestCollectSelectClauses(t *testing.T) {
	// "SELECT 1 " — after target list, should offer clause keywords
	cs := Collect("SELECT 1 ", 9)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	for _, tok := range []int{kwFROM, kwWHERE, kwGROUP, kwORDER, kwLIMIT} {
		if !cs.HasToken(tok) {
			t.Errorf("missing clause keyword after target list: %s", TokenName(tok))
		}
	}
}

func TestCollectAfterFrom(t *testing.T) {
	cs := Collect("SELECT 1 FROM ", 14)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !cs.HasRule("table_ref") {
		t.Error("expected table_ref rule candidate after FROM")
	}
	if !cs.HasRule("database_ref") {
		t.Error("expected database_ref rule candidate after FROM")
	}
}

func TestCollectAfterFromTable(t *testing.T) {
	// "SELECT 1 FROM t " — after FROM clause table
	cs := Collect("SELECT 1 FROM t ", 16)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	for _, tok := range []int{kwWHERE, kwGROUP, kwORDER, kwLIMIT, kwJOIN, kwLEFT, kwRIGHT, kwCROSS, kwINNER, kwNATURAL, kwUNION} {
		if !cs.HasToken(tok) {
			t.Errorf("missing clause keyword after FROM table: %s", TokenName(tok))
		}
	}
}

func TestCollectAfterJoin(t *testing.T) {
	// "SELECT 1 FROM t1 JOIN " — should offer table_ref
	cs := Collect("SELECT 1 FROM t1 JOIN ", 22)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !cs.HasRule("table_ref") {
		t.Error("expected table_ref rule after JOIN")
	}
}

func TestCollectAfterJoinOn(t *testing.T) {
	// "SELECT 1 FROM t1 JOIN t2 ON " — should offer columnref
	cs := Collect("SELECT 1 FROM t1 JOIN t2 ON ", 28)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !cs.HasRule("columnref") {
		t.Error("expected columnref rule after ON")
	}
}

func TestCollectAfterWhere(t *testing.T) {
	// "SELECT 1 FROM t WHERE " — should offer expression context
	cs := Collect("SELECT 1 FROM t WHERE ", 22)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !cs.HasRule("columnref") {
		t.Error("expected columnref rule in WHERE clause")
	}
	if !cs.HasRule("func_name") {
		t.Error("expected func_name rule in WHERE clause")
	}
}

func TestCollectAfterGroupBy(t *testing.T) {
	// "SELECT 1 FROM t GROUP BY x " — after GROUP BY, should offer HAVING, ORDER, LIMIT
	cs := Collect("SELECT 1 FROM t GROUP BY x ", 27)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	for _, tok := range []int{kwHAVING, kwORDER, kwLIMIT} {
		if !cs.HasToken(tok) {
			t.Errorf("missing keyword after GROUP BY: %s", TokenName(tok))
		}
	}
}

func TestCollectAfterOrderByExpr(t *testing.T) {
	// "SELECT 1 FROM t ORDER BY x " — should offer ASC, DESC, LIMIT
	cs := Collect("SELECT 1 FROM t ORDER BY x ", 27)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	for _, tok := range []int{kwASC, kwDESC, kwLIMIT} {
		if !cs.HasToken(tok) {
			t.Errorf("missing keyword after ORDER BY expr: %s", TokenName(tok))
		}
	}
}

// --- INSERT positions ---

func TestCollectInsertInto(t *testing.T) {
	// "INSERT INTO " — should offer table_ref
	cs := Collect("INSERT INTO ", 12)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !cs.HasRule("table_ref") {
		t.Error("expected table_ref rule after INSERT INTO")
	}
}

func TestCollectInsertAfterTable(t *testing.T) {
	// "INSERT INTO t1 " — should offer VALUES, SET, SELECT, PARTITION
	cs := Collect("INSERT INTO t1 ", 15)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	for _, tok := range []int{kwVALUES, kwSET, kwSELECT} {
		if !cs.HasToken(tok) {
			t.Errorf("missing token after INSERT INTO table: %s", TokenName(tok))
		}
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

func TestCollectReplaceInto(t *testing.T) {
	// "REPLACE INTO " — should behave like INSERT INTO
	cs := Collect("REPLACE INTO ", 13)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !cs.HasRule("table_ref") {
		t.Error("expected table_ref rule after REPLACE INTO")
	}
}

// --- UPDATE positions ---

func TestCollectUpdateTable(t *testing.T) {
	// "UPDATE " — should offer table_ref
	cs := Collect("UPDATE ", 7)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !cs.HasRule("table_ref") {
		t.Error("expected table_ref rule after UPDATE")
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
	// "UPDATE t1 SET a = 1 " — should offer WHERE
	cs := Collect("UPDATE t1 SET a = 1 ", 20)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !cs.HasToken(kwWHERE) {
		t.Error("expected WHERE after SET clause")
	}
}

// --- DELETE positions ---

func TestCollectDeleteFrom(t *testing.T) {
	// "DELETE FROM " — should offer table_ref
	cs := Collect("DELETE FROM ", 12)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !cs.HasRule("table_ref") {
		t.Error("expected table_ref rule after DELETE FROM")
	}
}

func TestCollectDeleteAfterTable(t *testing.T) {
	// "DELETE FROM t1 " — should offer WHERE
	cs := Collect("DELETE FROM t1 ", 15)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !cs.HasToken(kwWHERE) {
		t.Error("expected WHERE after DELETE FROM table")
	}
}

func TestCollectDeleteWhere(t *testing.T) {
	// "DELETE FROM t1 WHERE " — should offer expression context
	cs := Collect("DELETE FROM t1 WHERE ", 21)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !cs.HasRule("columnref") {
		t.Error("expected columnref rule in DELETE WHERE clause")
	}
}

// --- Expression contexts ---

func TestCollectExprStart(t *testing.T) {
	// Expression start should offer columnref and func_name
	cs := Collect("SELECT ", 7)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !cs.HasRule("columnref") {
		t.Error("expected columnref rule in expression start")
	}
	if !cs.HasRule("func_name") {
		t.Error("expected func_name rule in expression start")
	}
}

func TestCollectAfterIS(t *testing.T) {
	// "SELECT 1 FROM t WHERE x IS " — should offer NULL, NOT, TRUE, FALSE
	cs := Collect("SELECT 1 FROM t WHERE x IS ", 27)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	for _, tok := range []int{kwNULL, kwNOT, kwTRUE, kwFALSE} {
		if !cs.HasToken(tok) {
			t.Errorf("missing IS operand: %s", TokenName(tok))
		}
	}
}

func TestCollectExistsSubquery(t *testing.T) {
	// "SELECT * FROM t WHERE EXISTS (" — should offer SELECT
	cs := Collect("SELECT * FROM t WHERE EXISTS (", 30)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !cs.HasToken(kwSELECT) {
		t.Error("expected SELECT inside EXISTS (...)")
	}
}

func TestCollectCastType(t *testing.T) {
	// "SELECT CAST(x AS " — should offer type_name
	cs := Collect("SELECT CAST(x AS ", 17)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !cs.HasRule("type_name") {
		t.Error("expected type_name rule in CAST(... AS ...)")
	}
}

func TestCollectWindowClause(t *testing.T) {
	// "SELECT SUM(x) OVER (" — should offer PARTITION, ORDER, ROWS, RANGE
	cs := Collect("SELECT SUM(x) OVER (", 20)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	for _, tok := range []int{kwPARTITION, kwORDER} {
		if !cs.HasToken(tok) {
			t.Errorf("missing window clause keyword: %s", TokenName(tok))
		}
	}
}

// --- CTE and alias position tracking ---

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

func TestCollectCTEPositionsMultiple(t *testing.T) {
	// Multiple CTEs — should record one position (the WITH keyword)
	sql := "WITH a AS (SELECT 1), b AS (SELECT 2) SELECT "
	cs := Collect(sql, len(sql))
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if len(cs.CTEPositions) != 1 {
		t.Errorf("expected 1 CTE position, got %d", len(cs.CTEPositions))
	}
}

func TestCollectNoAliasPositionsWithoutAlias(t *testing.T) {
	// "SELECT a, b FROM " — no aliases, should have 0 alias positions
	sql := "SELECT a, b FROM "
	cs := Collect(sql, len(sql))
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if len(cs.SelectAliasPositions) != 0 {
		t.Errorf("expected 0 alias positions, got %d", len(cs.SelectAliasPositions))
	}
}

// --- ALTER TABLE ---

func TestCollectAlterTableOps(t *testing.T) {
	// "ALTER TABLE t1 " — should offer ALTER TABLE sub-commands
	cs := Collect("ALTER TABLE t1 ", 15)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	for _, tok := range []int{kwADD, kwDROP, kwMODIFY, kwCHANGE, kwRENAME, kwALTER} {
		if !cs.HasToken(tok) {
			t.Errorf("missing ALTER TABLE sub-command: %s", TokenName(tok))
		}
	}
}

func TestCollectAlterTableAdd(t *testing.T) {
	// "ALTER TABLE t1 ADD " — should offer COLUMN, INDEX, KEY, UNIQUE, PRIMARY, etc.
	cs := Collect("ALTER TABLE t1 ADD ", 19)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	for _, tok := range []int{kwCOLUMN, kwINDEX, kwKEY, kwUNIQUE, kwPRIMARY, kwFOREIGN, kwCONSTRAINT, kwCHECK} {
		if !cs.HasToken(tok) {
			t.Errorf("missing ALTER TABLE ADD sub-keyword: %s", TokenName(tok))
		}
	}
}

func TestCollectAlterTableDrop(t *testing.T) {
	// "ALTER TABLE t1 DROP " — should offer COLUMN, INDEX, KEY, etc.
	cs := Collect("ALTER TABLE t1 DROP ", 20)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	for _, tok := range []int{kwCOLUMN, kwINDEX, kwKEY, kwFOREIGN, kwPRIMARY} {
		if !cs.HasToken(tok) {
			t.Errorf("missing ALTER TABLE DROP sub-keyword: %s", TokenName(tok))
		}
	}
}

func TestCollectAlterTableDropColumn(t *testing.T) {
	// "ALTER TABLE t1 DROP COLUMN " — should offer columnref
	cs := Collect("ALTER TABLE t1 DROP COLUMN ", 27)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !cs.HasRule("columnref") {
		t.Error("expected columnref rule after DROP COLUMN")
	}
}

func TestCollectAlterTableDropIndex(t *testing.T) {
	// "ALTER TABLE t1 DROP INDEX " — should offer index_ref
	cs := Collect("ALTER TABLE t1 DROP INDEX ", 26)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !cs.HasRule("index_ref") {
		t.Error("expected index_ref rule after DROP INDEX")
	}
}

func TestCollectAlterTableModify(t *testing.T) {
	// "ALTER TABLE t1 MODIFY " — should offer COLUMN and columnref
	cs := Collect("ALTER TABLE t1 MODIFY ", 22)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !cs.HasRule("columnref") {
		t.Error("expected columnref rule after MODIFY")
	}
}

func TestCollectAlterTableRenameIndex(t *testing.T) {
	// "ALTER TABLE t1 RENAME INDEX " — should offer index_ref
	cs := Collect("ALTER TABLE t1 RENAME INDEX ", 28)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !cs.HasRule("index_ref") {
		t.Error("expected index_ref rule after RENAME INDEX")
	}
}

// --- Dot-qualified completion ---

func TestCollectAfterDot(t *testing.T) {
	// "SELECT t." — should offer columnref after dot
	cs := Collect("SELECT t.", 9)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !cs.HasRule("columnref") {
		t.Error("expected columnref rule after table dot")
	}
}

func TestCollectAfterSchemaDot(t *testing.T) {
	// "SELECT s.t." — should offer columnref after schema.table.
	cs := Collect("SELECT s.t.", 11)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !cs.HasRule("columnref") {
		t.Error("expected columnref rule after schema.table dot")
	}
}

// --- INSERT empty parens ---

func TestCollectInsertEmptyParens(t *testing.T) {
	// "INSERT INTO t1()" — cursor between ( and ) should offer columnref
	cs := Collect("INSERT INTO t1()", 15)
	if cs == nil {
		t.Fatal("expected non-nil candidates")
	}
	if !cs.HasRule("columnref") {
		t.Error("expected columnref inside empty parens")
	}
}

// --- Derived table column alias ---

func TestParseDerivedTableColumnAlias(t *testing.T) {
	sql := "SELECT * FROM (SELECT c1, c2 FROM t1) AS dt(a, b)"
	_, err := Parse(sql)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
}

func TestParseDerivedTableColumnAliasSingle(t *testing.T) {
	sql := "SELECT * FROM (SELECT 1) AS t(x)"
	_, err := Parse(sql)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
}

func TestParseLateralDerivedColumnAlias(t *testing.T) {
	sql := "SELECT * FROM t1, LATERAL (SELECT c1 FROM t2) AS dt(a)"
	_, err := Parse(sql)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
}

// --- Tokenize and Token.End ---

func TestTokenize(t *testing.T) {
	tokens := Tokenize("SELECT a FROM t")
	if len(tokens) != 4 {
		t.Fatalf("expected 4 tokens, got %d", len(tokens))
	}
	if tokens[0].Str != "SELECT" {
		t.Errorf("expected first token str SELECT, got %q", tokens[0].Str)
	}
	if tokens[0].Loc != 0 {
		t.Errorf("expected first token Loc=0, got %d", tokens[0].Loc)
	}
	if tokens[0].End != 6 {
		t.Errorf("expected first token End=6, got %d", tokens[0].End)
	}
}

func TestTokenizeEmpty(t *testing.T) {
	tokens := Tokenize("")
	if len(tokens) != 0 {
		t.Errorf("expected 0 tokens for empty input, got %d", len(tokens))
	}
}

func TestTokenEnd(t *testing.T) {
	tokens := Tokenize("SELECT *, a AS b FROM t WHERE id = 1")
	for _, tok := range tokens {
		if tok.End <= tok.Loc {
			t.Errorf("token %q: End (%d) should be > Loc (%d)", tok.Str, tok.End, tok.Loc)
		}
	}
}

func TestTokenEndMultiChar(t *testing.T) {
	tokens := Tokenize("SELECT 1 <=> 2")
	// Find the <=> token
	found := false
	for _, tok := range tokens {
		if tok.Str == "<=>" {
			found = true
			if tok.End-tok.Loc != 3 {
				t.Errorf("<=> token: expected length 3, got End-Loc=%d", tok.End-tok.Loc)
			}
		}
	}
	if !found {
		t.Error("did not find <=> token")
	}
}
