package completion

import (
	"testing"

	"github.com/bytebase/omni/trino/catalog"
)

// These tests cover the context-detection branches and scope edge cases that
// the headline scenario tests in completion_test.go do not exercise directly,
// keeping every completion context accounted for (completeness gate).

func TestContext_SelectListComma(t *testing.T) {
	cat := buildCatalog()
	// A comma in the SELECT list extends the column list.
	sql := "SELECT orderkey,  FROM orders"
	got := Complete(sql, len("SELECT orderkey, "), cat)
	if !has(got, CandidateColumn, "custkey") {
		t.Errorf("comma in SELECT list should offer columns; columns=%v", texts(got, CandidateColumn))
	}
}

func TestContext_FromListComma(t *testing.T) {
	cat := buildCatalog()
	// A comma in the FROM list extends the relation list (implicit cross join).
	sql := "SELECT * FROM orders, "
	got := Complete(sql, len("SELECT * FROM orders, "), cat)
	if !has(got, CandidateTable, "customer") {
		t.Errorf("comma in FROM list should offer relations; tables=%v", texts(got, CandidateTable))
	}
}

func TestContext_GroupByColumns(t *testing.T) {
	cat := buildCatalog()
	sql := "SELECT custkey, count(*) FROM orders GROUP BY "
	got := Complete(sql, len(sql), cat)
	if !has(got, CandidateColumn, "custkey") {
		t.Errorf("GROUP BY should offer columns; columns=%v", texts(got, CandidateColumn))
	}
}

func TestContext_OrderByColumns(t *testing.T) {
	cat := buildCatalog()
	sql := "SELECT * FROM orders ORDER BY "
	got := Complete(sql, len(sql), cat)
	if !has(got, CandidateColumn, "totalprice") {
		t.Errorf("ORDER BY should offer columns; columns=%v", texts(got, CandidateColumn))
	}
}

func TestContext_AfterComparisonOperator(t *testing.T) {
	cat := buildCatalog()
	// After a comparison operator an expression operand is expected: columns +
	// expression keywords.
	sql := "SELECT * FROM orders WHERE totalprice > "
	got := Complete(sql, len(sql), cat)
	if !has(got, CandidateColumn, "custkey") {
		t.Errorf("after '>' should offer columns; columns=%v", texts(got, CandidateColumn))
	}
	if !has(got, CandidateKeyword, "CAST") {
		t.Errorf("after '>' should offer expression keywords; keywords=%v", texts(got, CandidateKeyword))
	}
}

func TestContext_AfterAndOr(t *testing.T) {
	cat := buildCatalog()
	sql := "SELECT * FROM customer WHERE custkey = 1 AND "
	got := Complete(sql, len(sql), cat)
	if !has(got, CandidateColumn, "name") {
		t.Errorf("after AND should offer columns; columns=%v", texts(got, CandidateColumn))
	}
}

func TestContext_IntoRelation(t *testing.T) {
	cat := buildCatalog()
	// INSERT INTO <relation>.
	sql := "INSERT INTO "
	got := Complete(sql, len(sql), cat)
	if !has(got, CandidateTable, "orders") {
		t.Errorf("INSERT INTO should offer relations; tables=%v", texts(got, CandidateTable))
	}
}

func TestContext_UpdateRelation(t *testing.T) {
	cat := buildCatalog()
	sql := "UPDATE "
	got := Complete(sql, len(sql), cat)
	if !has(got, CandidateTable, "customer") {
		t.Errorf("UPDATE should offer relations; tables=%v", texts(got, CandidateTable))
	}
}

func TestContext_DescribeRelation(t *testing.T) {
	cat := buildCatalog()
	sql := "DESCRIBE "
	got := Complete(sql, len(sql), cat)
	if !has(got, CandidateTable, "orders") {
		t.Errorf("DESCRIBE should offer relations; tables=%v", texts(got, CandidateTable))
	}
}

func TestContext_FallbackClauseKeywords(t *testing.T) {
	cat := buildCatalog()
	// Right after a complete table name, the useful next tokens are clause
	// keywords (WHERE, JOIN, GROUP, ...). This is the kindKeyword fallback.
	sql := "SELECT * FROM orders "
	got := Complete(sql, len(sql), cat)
	if !has(got, CandidateKeyword, "WHERE") {
		t.Errorf("after a table name should offer clause keyword WHERE; keywords=%v", texts(got, CandidateKeyword))
	}
	if !has(got, CandidateKeyword, "JOIN") {
		t.Errorf("after a table name should offer clause keyword JOIN; keywords=%v", texts(got, CandidateKeyword))
	}
}

func TestContext_EmptySelectListCaretInWhere(t *testing.T) {
	cat := buildCatalog()
	// The robustness fallback: select list is empty AND the caret is in the
	// WHERE. fillEmptySelectLists must recover the orders scope.
	sql := "SELECT  FROM orders WHERE "
	got := Complete(sql, len(sql), cat)
	if !has(got, CandidateColumn, "orderkey") {
		t.Errorf("empty select list + WHERE caret should still recover columns; columns=%v", texts(got, CandidateColumn))
	}
}

func TestContext_QuotedQualifier(t *testing.T) {
	// A quoted, case-sensitive schema qualifier must resolve against the
	// case-preserved catalog key.
	cat := catalog.New()
	tpch := cat.EnsureCatalog("tpch")
	// Schema stored with preserved case (as if created from a quoted name).
	tpch.EnsureSchema("MySchema").AddTable("t1", catalog.NewColumn("c", "bigint", false))
	cat.SetCurrentCatalog("tpch")
	cat.SetCurrentSchema("MySchema")

	// "MySchema". (quoted) -> its relations.
	sql := `SELECT * FROM "MySchema".`
	got := Complete(sql, len(sql), cat)
	if !has(got, CandidateTable, "t1") {
		t.Errorf(`"MySchema". should offer t1; tables=%v`, texts(got, CandidateTable))
	}
}

func TestContext_CaretAfterLastSemicolon(t *testing.T) {
	cat := buildCatalog()
	// Caret in the whitespace after the final ';' (no further statement text):
	// statement-start keywords, no crash.
	sql := "SELECT 1;   "
	got := Complete(sql, len(sql), cat)
	if !has(got, CandidateKeyword, "SELECT") {
		t.Errorf("caret after final ';' should offer statement-start keywords; got %v", texts(got, CandidateKeyword))
	}
}

func TestContext_SetClauseColumns(t *testing.T) {
	cat := buildCatalog()
	// UPDATE ... SET <column>. The SET keyword puts us in a column-ONLY context
	// (the assignment target must be a bare column), and the DML target (orders)
	// — which query-span omits — is recovered by the completion package's own
	// DML-target scan.
	sql := "UPDATE orders SET "
	got := Complete(sql, len(sql), cat)
	if !has(got, CandidateColumn, "totalprice") {
		t.Errorf("UPDATE ... SET should offer columns; columns=%v", texts(got, CandidateColumn))
	}
	// An expression keyword in the SET target is a syntax error (oracle 481:
	// "UPDATE t SET CASE" => SYNTAX_ERROR), so SET must NOT offer them.
	if has(got, CandidateKeyword, "CASE") || has(got, CandidateKeyword, "CAST") {
		t.Errorf("UPDATE ... SET must not offer expression keywords; keywords=%v", texts(got, CandidateKeyword))
	}
}

func TestContext_UsingNoBareColumns(t *testing.T) {
	cat := buildCatalog()
	// After a JOIN's USING keyword Trino requires "(col, ...)" (oracle 481:
	// "... USING id" => SYNTAX_ERROR), so the completer must NOT offer a bare
	// column candidate that would yield invalid syntax.
	sql := "SELECT * FROM orders o JOIN customer c USING "
	got := Complete(sql, len(sql), cat)
	if cols := texts(got, CandidateColumn); len(cols) != 0 {
		t.Errorf("after USING (before the paren) a bare column is invalid; must not offer columns, got %v", cols)
	}
}

func TestContext_UpdateWhereColumns(t *testing.T) {
	cat := buildCatalog()
	// The DML target is in scope for the WHERE predicate too.
	sql := "UPDATE customer SET name = 'x' WHERE "
	got := Complete(sql, len(sql), cat)
	if !has(got, CandidateColumn, "custkey") {
		t.Errorf("UPDATE ... WHERE should offer target columns; columns=%v", texts(got, CandidateColumn))
	}
}

func TestContext_DeleteWhereColumns(t *testing.T) {
	cat := buildCatalog()
	sql := "DELETE FROM orders WHERE "
	got := Complete(sql, len(sql), cat)
	if !has(got, CandidateColumn, "orderkey") {
		t.Errorf("DELETE ... WHERE should offer target columns; columns=%v", texts(got, CandidateColumn))
	}
}

func TestCandidateTypeString(t *testing.T) {
	cases := map[CandidateType]string{
		CandidateKeyword:  "keyword",
		CandidateCatalog:  "catalog",
		CandidateSchema:   "schema",
		CandidateTable:    "table",
		CandidateView:     "view",
		CandidateColumn:   "column",
		CandidateType(99): "unknown",
	}
	for typ, want := range cases {
		if got := typ.String(); got != want {
			t.Errorf("CandidateType(%d).String() = %q, want %q", typ, got, want)
		}
	}
}

func TestContext_CommaInGroupByNotSelectList(t *testing.T) {
	cat := buildCatalog()
	// A comma inside GROUP BY must offer columns (the GROUP BY list), and the
	// nearestListKeyword scan must NOT misattribute it to the SELECT list as a
	// relation context. Either way the user wants columns here.
	sql := "SELECT custkey FROM orders GROUP BY custkey, "
	got := Complete(sql, len(sql), cat)
	if !has(got, CandidateColumn, "totalprice") {
		t.Errorf("comma in GROUP BY should offer columns; columns=%v", texts(got, CandidateColumn))
	}
	// It must NOT be a relation context (no tables offered as the primary set).
	if has(got, CandidateTable, "customer") {
		t.Errorf("comma in GROUP BY wrongly offered relations; tables=%v", texts(got, CandidateTable))
	}
}

func TestContext_UpdateTargetDottedByTableName(t *testing.T) {
	cat := buildCatalog()
	// UPDATE has no target alias (oracle 481: "UPDATE t a SET" => SYNTAX_ERROR),
	// so a dotted reference uses the table name itself.
	sql := "UPDATE orders SET orders."
	got := Complete(sql, len(sql), cat)
	if !has(got, CandidateColumn, "totalprice") {
		t.Errorf("UPDATE orders ... orders. should offer target columns; columns=%v", texts(got, CandidateColumn))
	}
}

func TestContext_UpdateTargetAliasRejected(t *testing.T) {
	cat := buildCatalog()
	// "UPDATE orders o ..." is invalid Trino, so the completer must NOT invent an
	// alias scope for `o`: a dotted `o.` reference resolves to nothing.
	sql := "UPDATE orders o SET o."
	got := Complete(sql, len(sql), cat)
	if cols := texts(got, CandidateColumn); len(cols) != 0 {
		t.Errorf("UPDATE orders o (invalid alias) must not create an alias scope; columns=%v", cols)
	}
}

func TestContext_MergeTargetAliasColumns(t *testing.T) {
	cat := buildCatalog()
	// MERGE DOES allow a target alias (oracle 481 accepts it), so a dotted
	// reference through the alias resolves to the target's columns.
	sql := "MERGE INTO orders o USING customer c ON o.custkey = c.custkey WHEN MATCHED THEN UPDATE SET totalprice = o."
	got := Complete(sql, len(sql), cat)
	if !has(got, CandidateColumn, "totalprice") {
		t.Errorf("MERGE alias o. should offer target columns; columns=%v", texts(got, CandidateColumn))
	}
}
