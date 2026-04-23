package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

// TestFullTextPredicate verifies CONTAINS / FREETEXT parse as FullTextPredicate
// inside a WHERE search condition, matching SqlScriptDOM's FullTextPredicate
// shape.
func TestFullTextPredicate(t *testing.T) {
	tests := []struct {
		name           string
		sql            string
		wantFunc       ast.FullTextFunc
		wantColCount   int
		wantStarIdx    []int // indices in Columns that should be StarExpr
		wantProperty   string
		wantValue      string // expected Literal.Str or Variable.Name when value is a variable
		wantValueVar   bool
		wantLanguage   string // empty => expect no LanguageTerm
		wantLanguageN  bool   // language is N'...'
	}{
		{
			name:         "contains-single-col",
			sql:          "SELECT * FROM t WHERE CONTAINS(b, 'foo')",
			wantFunc:     ast.FullTextContains,
			wantColCount: 1,
			wantValue:    "foo",
		},
		{
			name:         "freetext-single-col",
			sql:          "SELECT * FROM t WHERE FREETEXT(b, 'foo')",
			wantFunc:     ast.FullTextFreeText,
			wantColCount: 1,
			wantValue:    "foo",
		},
		{
			name:         "contains-star",
			sql:          "SELECT * FROM t WHERE CONTAINS(*, 'foo')",
			wantFunc:     ast.FullTextContains,
			wantColCount: 1,
			wantStarIdx:  []int{0},
			wantValue:    "foo",
		},
		{
			name:         "contains-paren-list",
			sql:          "SELECT * FROM t WHERE CONTAINS((a, b), 'foo')",
			wantFunc:     ast.FullTextContains,
			wantColCount: 2,
			wantValue:    "foo",
		},
		{
			name:         "contains-paren-star",
			sql:          "SELECT * FROM t WHERE CONTAINS((*), 'foo')",
			wantFunc:     ast.FullTextContains,
			wantColCount: 1,
			wantStarIdx:  []int{0},
			wantValue:    "foo",
		},
		{
			name:         "contains-property",
			sql:          "SELECT * FROM t WHERE CONTAINS(PROPERTY(b, 'Title'), 'foo')",
			wantFunc:     ast.FullTextContains,
			wantColCount: 1,
			wantProperty: "Title",
			wantValue:    "foo",
		},
		{
			name:         "contains-language",
			sql:          "SELECT * FROM t WHERE CONTAINS((a, b), 'foo', LANGUAGE 'English')",
			wantFunc:     ast.FullTextContains,
			wantColCount: 2,
			wantValue:    "foo",
			wantLanguage: "English",
		},
		{
			name:         "contains-variable-value",
			sql:          "DECLARE @v NVARCHAR(100); SELECT * FROM t WHERE CONTAINS(b, @v)",
			wantFunc:     ast.FullTextContains,
			wantColCount: 1,
			wantValue:    "@v",
			wantValueVar: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			pred := findFullTextPredicate(t, list)
			if pred == nil {
				t.Fatalf("no FullTextPredicate in parse tree: %s", ast.NodeToString(list))
			}
			if pred.Func != tt.wantFunc {
				t.Errorf("func = %v, want %v", pred.Func, tt.wantFunc)
			}
			if pred.Columns == nil {
				t.Fatalf("Columns is nil")
			}
			if got := len(pred.Columns.Items); got != tt.wantColCount {
				t.Errorf("Columns len = %d, want %d", got, tt.wantColCount)
			}
			for _, idx := range tt.wantStarIdx {
				if idx >= len(pred.Columns.Items) {
					t.Errorf("expected star at idx %d but list too short", idx)
					continue
				}
				if _, ok := pred.Columns.Items[idx].(*ast.StarExpr); !ok {
					t.Errorf("Columns[%d] = %T, want *StarExpr", idx, pred.Columns.Items[idx])
				}
			}
			if tt.wantProperty != "" {
				if pred.PropertyName == nil {
					t.Errorf("PropertyName is nil, want literal %q", tt.wantProperty)
				} else if pred.PropertyName.Str != tt.wantProperty {
					t.Errorf("PropertyName = %q, want %q", pred.PropertyName.Str, tt.wantProperty)
				}
			} else if pred.PropertyName != nil {
				t.Errorf("PropertyName = %q, want nil", pred.PropertyName.Str)
			}
			if tt.wantValueVar {
				v, ok := pred.Value.(*ast.VariableRef)
				if !ok {
					t.Fatalf("Value = %T, want *VariableRef", pred.Value)
				}
				if v.Name != tt.wantValue {
					t.Errorf("Value.Name = %q, want %q", v.Name, tt.wantValue)
				}
			} else {
				lit, ok := pred.Value.(*ast.Literal)
				if !ok {
					t.Fatalf("Value = %T, want *Literal", pred.Value)
				}
				if lit.Str != tt.wantValue {
					t.Errorf("Value.Str = %q, want %q", lit.Str, tt.wantValue)
				}
			}
			if tt.wantLanguage != "" {
				lit, ok := pred.LanguageTerm.(*ast.Literal)
				if !ok {
					t.Fatalf("LanguageTerm = %T, want *Literal", pred.LanguageTerm)
				}
				if lit.Str != tt.wantLanguage {
					t.Errorf("LanguageTerm = %q, want %q", lit.Str, tt.wantLanguage)
				}
			} else if pred.LanguageTerm != nil {
				t.Errorf("LanguageTerm = %v, want nil", pred.LanguageTerm)
			}
			if pred.Loc.Start < 0 || pred.Loc.End < 0 {
				t.Errorf("Loc not populated: %+v", pred.Loc)
			}
		})
	}
}

// TestFullTextPredicateRejectsAsScalar verifies that CONTAINS/FREETEXT
// CANNOT appear as a scalar expression — they are only valid inside
// search_condition positions. This matches SqlScriptDOM's behavior.
func TestFullTextPredicateRejectsAsScalar(t *testing.T) {
	rejects := []string{
		// bare scalar
		"SELECT CONTAINS(b, 'foo')",
		"SELECT FREETEXT(b, 'foo')",
		// nested inside subquery select list
		"SELECT c FROM t WHERE b > (SELECT CONTAINS(col, 'x') FROM u)",
		// arithmetic context
		"SELECT 1 + CONTAINS(b, 'foo')",
		// function arg
		"SELECT COALESCE(CONTAINS(b, 'foo'), 0)",
	}
	for _, sql := range rejects {
		t.Run(sql, func(t *testing.T) {
			_, err := Parse(sql)
			if err == nil {
				t.Errorf("Parse(%q): expected error but got nil", sql)
			}
		})
	}
}

// TestFullTextPredicateWithNot verifies that NOT CONTAINS(...) produces a
// UnaryNot wrapping a FullTextPredicate, matching SqlScriptDOM's
// BooleanNotExpression{Expression: FullTextPredicate} layout.
func TestFullTextPredicateWithNot(t *testing.T) {
	sql := "SELECT * FROM t WHERE NOT CONTAINS(b, 'foo')"
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	sel, ok := list.Items[0].(*ast.SelectStmt)
	if !ok {
		t.Fatalf("stmt = %T, want SelectStmt", list.Items[0])
	}
	unary, ok := sel.WhereClause.(*ast.UnaryExpr)
	if !ok {
		t.Fatalf("WhereClause = %T, want *UnaryExpr (NOT)", sel.WhereClause)
	}
	if unary.Op != ast.UnaryNot {
		t.Errorf("UnaryExpr.Op = %v, want UnaryNot", unary.Op)
	}
	if _, ok := unary.Operand.(*ast.FullTextPredicate); !ok {
		t.Errorf("operand = %T, want *FullTextPredicate", unary.Operand)
	}
}

// TestFullTextTableRef verifies CONTAINSTABLE / FREETEXTTABLE parse to a
// FullTextTableRef in FROM position.
func TestFullTextTableRef(t *testing.T) {
	tests := []struct {
		name         string
		sql          string
		wantFunc     ast.FullTextFunc
		wantColCount int
		wantProperty string
		wantSC       string
		wantAlias    string
		wantLanguage string
		wantTopN     int64
		hasTopN      bool
	}{
		{
			name:         "containstable-single",
			sql:          "SELECT * FROM CONTAINSTABLE(t, b, 'foo') ct",
			wantFunc:     ast.FullTextContains,
			wantColCount: 1,
			wantSC:       "foo",
			wantAlias:    "ct",
		},
		{
			name:         "freetexttable-paren-list",
			sql:          "SELECT * FROM FREETEXTTABLE(t, (a, b), 'foo') ft",
			wantFunc:     ast.FullTextFreeText,
			wantColCount: 2,
			wantSC:       "foo",
			wantAlias:    "ft",
		},
		{
			name:         "containstable-language-topn",
			sql:          "SELECT * FROM CONTAINSTABLE(t, (a, b), 'foo', LANGUAGE 'English', 10) ct",
			wantFunc:     ast.FullTextContains,
			wantColCount: 2,
			wantSC:       "foo",
			wantAlias:    "ct",
			wantLanguage: "English",
			wantTopN:     10,
			hasTopN:      true,
		},
		{
			name:         "containstable-topn-language",
			sql:          "SELECT * FROM CONTAINSTABLE(t, b, 'foo', 5, LANGUAGE 'English') ct",
			wantFunc:     ast.FullTextContains,
			wantColCount: 1,
			wantSC:       "foo",
			wantAlias:    "ct",
			wantLanguage: "English",
			wantTopN:     5,
			hasTopN:      true,
		},
		{
			name:         "containstable-property",
			sql:          "SELECT * FROM CONTAINSTABLE(t, PROPERTY(b, 'Title'), 'foo') ct",
			wantFunc:     ast.FullTextContains,
			wantColCount: 1,
			wantProperty: "Title",
			wantSC:       "foo",
			wantAlias:    "ct",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			ref := findFullTextTableRef(t, list)
			if ref == nil {
				t.Fatalf("no FullTextTableRef in parse tree: %s", ast.NodeToString(list))
			}
			if ref.Func != tt.wantFunc {
				t.Errorf("Func = %v, want %v", ref.Func, tt.wantFunc)
			}
			if ref.Table == nil || ref.Table.Object == "" {
				t.Errorf("Table missing: %+v", ref.Table)
			}
			if ref.Columns == nil || len(ref.Columns.Items) != tt.wantColCount {
				got := 0
				if ref.Columns != nil {
					got = len(ref.Columns.Items)
				}
				t.Errorf("Columns len = %d, want %d", got, tt.wantColCount)
			}
			if tt.wantProperty != "" {
				if ref.PropertyName == nil {
					t.Errorf("PropertyName = nil, want %q", tt.wantProperty)
				} else if ref.PropertyName.Str != tt.wantProperty {
					t.Errorf("PropertyName = %q, want %q", ref.PropertyName.Str, tt.wantProperty)
				}
			} else if ref.PropertyName != nil {
				t.Errorf("PropertyName = %q, want nil", ref.PropertyName.Str)
			}
			if lit, ok := ref.SearchCondition.(*ast.Literal); ok {
				if lit.Str != tt.wantSC {
					t.Errorf("SearchCondition = %q, want %q", lit.Str, tt.wantSC)
				}
			} else {
				t.Errorf("SearchCondition = %T, want *Literal", ref.SearchCondition)
			}
			if tt.wantAlias != "" && ref.Alias != tt.wantAlias {
				t.Errorf("Alias = %q, want %q", ref.Alias, tt.wantAlias)
			}
			if tt.wantLanguage != "" {
				lit, ok := ref.Language.(*ast.Literal)
				if !ok || lit.Str != tt.wantLanguage {
					t.Errorf("Language = %v, want %q", ref.Language, tt.wantLanguage)
				}
			} else if ref.Language != nil {
				t.Errorf("Language = %v, want nil", ref.Language)
			}
			if tt.hasTopN {
				lit, ok := ref.TopN.(*ast.Literal)
				if !ok || lit.Ival != tt.wantTopN {
					t.Errorf("TopN = %v, want %d", ref.TopN, tt.wantTopN)
				}
			} else if ref.TopN != nil {
				t.Errorf("TopN = %v, want nil", ref.TopN)
			}
		})
	}
}

// TestSemanticTableRef verifies the three semantic TVFs parse correctly.
func TestSemanticTableRef(t *testing.T) {
	tests := []struct {
		name            string
		sql             string
		wantFunc        ast.SemanticFunc
		wantColCount    int
		wantSourceKey   int64
		wantHasKey      bool
		wantMatchedCol  string
		wantMatchedKey  int64
		wantAlias       string
	}{
		{
			name:         "keyphrase-no-key",
			sql:          "SELECT * FROM SEMANTICKEYPHRASETABLE(t, b) s",
			wantFunc:     ast.SemanticKeyPhraseTable,
			wantColCount: 1,
			wantAlias:    "s",
		},
		{
			name:          "keyphrase-with-key",
			sql:           "SELECT * FROM SEMANTICKEYPHRASETABLE(t, b, 1) s",
			wantFunc:      ast.SemanticKeyPhraseTable,
			wantColCount:  1,
			wantSourceKey: 1,
			wantHasKey:    true,
			wantAlias:     "s",
		},
		{
			name:          "similarity",
			sql:           "SELECT * FROM SEMANTICSIMILARITYTABLE(t, (a, b), 1) s",
			wantFunc:      ast.SemanticSimilarityTable,
			wantColCount:  2,
			wantSourceKey: 1,
			wantHasKey:    true,
			wantAlias:     "s",
		},
		{
			name:           "similarity-details",
			sql:            "SELECT * FROM SEMANTICSIMILARITYDETAILSTABLE(t, a, 1, b, 2) s",
			wantFunc:       ast.SemanticSimilarityDetailsTable,
			wantColCount:   1,
			wantSourceKey:  1,
			wantHasKey:     true,
			wantMatchedCol: "b",
			wantMatchedKey: 2,
			wantAlias:      "s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			ref := findSemanticTableRef(t, list)
			if ref == nil {
				t.Fatalf("no SemanticTableRef in parse tree: %s", ast.NodeToString(list))
			}
			if ref.Func != tt.wantFunc {
				t.Errorf("Func = %v, want %v", ref.Func, tt.wantFunc)
			}
			if ref.Table == nil || ref.Table.Object == "" {
				t.Errorf("Table missing: %+v", ref.Table)
			}
			if ref.Columns == nil || len(ref.Columns.Items) != tt.wantColCount {
				got := 0
				if ref.Columns != nil {
					got = len(ref.Columns.Items)
				}
				t.Errorf("Columns len = %d, want %d", got, tt.wantColCount)
			}
			if tt.wantHasKey {
				lit, ok := ref.SourceKey.(*ast.Literal)
				if !ok || lit.Ival != tt.wantSourceKey {
					t.Errorf("SourceKey = %v, want %d", ref.SourceKey, tt.wantSourceKey)
				}
			} else if ref.SourceKey != nil {
				t.Errorf("SourceKey = %v, want nil", ref.SourceKey)
			}
			if tt.wantMatchedCol != "" {
				if ref.MatchedColumn == nil {
					t.Errorf("MatchedColumn = nil, want %q", tt.wantMatchedCol)
				} else if ref.MatchedColumn.Column != tt.wantMatchedCol {
					t.Errorf("MatchedColumn = %q, want %q", ref.MatchedColumn.Column, tt.wantMatchedCol)
				}
				lit, ok := ref.MatchedKey.(*ast.Literal)
				if !ok || lit.Ival != tt.wantMatchedKey {
					t.Errorf("MatchedKey = %v, want %d", ref.MatchedKey, tt.wantMatchedKey)
				}
			}
			if tt.wantAlias != "" && ref.Alias != tt.wantAlias {
				t.Errorf("Alias = %q, want %q", ref.Alias, tt.wantAlias)
			}
		})
	}
}

// TestSearchFunctionSerialization verifies NodeToString round-trips through
// the outfuncs tag table — this is the regression lever that tripped the
// original CONTAINS regression (a node that the parser produced but outfuncs
// didn't know about would either panic or dump {UNKNOWN TYPE …}).
func TestSearchFunctionSerialization(t *testing.T) {
	sqls := []string{
		"SELECT * FROM t WHERE CONTAINS(b, 'foo')",
		"SELECT * FROM t WHERE FREETEXT((a, b), 'foo', LANGUAGE 'English')",
		"SELECT * FROM t WHERE CONTAINS(PROPERTY(b, 'Title'), 'foo')",
		"SELECT * FROM CONTAINSTABLE(t, b, 'foo') ct",
		"SELECT * FROM FREETEXTTABLE(t, (a,b), 'foo', LANGUAGE 'English', 10) ft",
		"SELECT * FROM SEMANTICKEYPHRASETABLE(t, b, 1) s",
		"SELECT * FROM SEMANTICSIMILARITYTABLE(t, (a, b), 1) s",
		"SELECT * FROM SEMANTICSIMILARITYDETAILSTABLE(t, a, 1, b, 2) s",
	}
	for _, sql := range sqls {
		t.Run(sql, func(t *testing.T) {
			list, err := Parse(sql)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			for _, stmt := range list.Items {
				s := ast.NodeToString(stmt)
				if strings.Contains(s, "UNKNOWN TYPE") {
					t.Errorf("NodeToString produced UNKNOWN TYPE sentinel: %s", s)
				}
				if s == "" {
					t.Errorf("NodeToString returned empty string")
				}
			}
		})
	}
}

// --- find helpers ---

func findFullTextPredicate(t *testing.T, list *ast.List) *ast.FullTextPredicate {
	t.Helper()
	var found *ast.FullTextPredicate
	ast.Walk(visitorFunc(func(n ast.Node) bool {
		if p, ok := n.(*ast.FullTextPredicate); ok {
			found = p
			return false
		}
		return true
	}), list)
	return found
}

func findFullTextTableRef(t *testing.T, list *ast.List) *ast.FullTextTableRef {
	t.Helper()
	var found *ast.FullTextTableRef
	ast.Walk(visitorFunc(func(n ast.Node) bool {
		if r, ok := n.(*ast.FullTextTableRef); ok {
			found = r
			return false
		}
		return true
	}), list)
	return found
}

func findSemanticTableRef(t *testing.T, list *ast.List) *ast.SemanticTableRef {
	t.Helper()
	var found *ast.SemanticTableRef
	ast.Walk(visitorFunc(func(n ast.Node) bool {
		if r, ok := n.(*ast.SemanticTableRef); ok {
			found = r
			return false
		}
		return true
	}), list)
	return found
}

type visitorFunc func(ast.Node) bool

func (f visitorFunc) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	if f(n) {
		return f
	}
	return nil
}
