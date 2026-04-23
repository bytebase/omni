package parser

import (
	"strings"
	"testing"
)

// This file is the regression guard for the "callable CoreKeyword" design.
//
// Background: prior to the 2026-04 parser-optimize refactor, omni's
// parsePrimary() default branch fell through to parseIdentExpr() for ANY
// keyword whose string was non-empty. That silently routed CORE keyword
// tokens like CONTAINS / FREETEXT into a FuncCallExpr shape and hid the
// fact that the explicit case list in parsePrimary was incomplete. When
// the refactor tightened the default branch to ContextKeyword only,
// CONTAINS / FREETEXT started failing to parse — an invisible regression
// for multiple call-sites at once.
//
// The fix landed in search_functions.go adds dedicated AST nodes and
// parse functions aligned with SqlScriptDOM's grammar. This file
// institutionalises the alignment: every T-SQL keyword that looks like a
// function call has an explicit role in the parser's dispatch tables, and
// this manifest-driven test asserts that each keyword parses in its
// declared context AND fails in all the others.
//
// If you add a new CoreKeyword-callable (e.g. a future TVF that becomes a
// reserved keyword), add it to the manifest below. Tests will then guide
// you to the correct dispatch site.

// keywordRole classifies where a Core keyword is allowed to appear in a
// function-call-like position.
type keywordRole int

const (
	// roleExprScalar — callable in scalar expression positions, typically
	// produces a dedicated scalar AST (CAST/CONVERT/CASE/...) or, absent
	// a custom parser, a FuncCallExpr.
	roleExprScalar keywordRole = iota
	// roleFullTextPredicate — valid only inside <search_condition> (WHERE /
	// HAVING / JOIN ON / CASE WHEN / ...), produces ast.FullTextPredicate.
	// CONTAINS, FREETEXT.
	roleFullTextPredicate
	// roleRowsetTable — rowset function, valid only as a FROM table source,
	// produces a FuncCallExpr-based AliasedTableRef. OPENROWSET / OPENQUERY
	// / OPENJSON / OPENDATASOURCE / OPENXML.
	roleRowsetTable
	// roleFullTextTable — valid only as a FROM table source, produces
	// ast.FullTextTableRef. CONTAINSTABLE, FREETEXTTABLE.
	roleFullTextTable
	// roleSemanticTable — valid only as a FROM table source, produces
	// ast.SemanticTableRef. SEMANTICKEYPHRASETABLE /
	// SEMANTICSIMILARITYTABLE / SEMANTICSIMILARITYDETAILSTABLE.
	roleSemanticTable
)

type callableKeyword struct {
	keyword string
	role    keywordRole
	// acceptSample is a SQL snippet that exercises the keyword in its
	// declared context. It must parse without error.
	acceptSample string
	// alsoScalar is true for FROM-side callables that SqlScriptDOM also
	// happens to accept as generic scalar FunctionCall in expression
	// position. OPENJSON is the canonical example: it is a ContextKeyword
	// (not reserved) so `SELECT OPENJSON(…)` parses as a scalar
	// FunctionCall. Keywords with alsoScalar=true are exempt from the
	// scalar-rejection half of TestCallableKeywordRoleRejection.
	alsoScalar bool
}

// callableKeywordManifest is the single source of truth for which CoreKeyword
// tokens are legal as function-call-like callables and what role they fill.
var callableKeywordManifest = []callableKeyword{
	// --- roleExprScalar (scalar-expression callables) ---
	{keyword: "CAST", role: roleExprScalar, acceptSample: "SELECT CAST(1 AS INT)"},
	{keyword: "CONVERT", role: roleExprScalar, acceptSample: "SELECT CONVERT(VARCHAR, 1)"},
	{keyword: "TRY_CAST", role: roleExprScalar, acceptSample: "SELECT TRY_CAST('1' AS INT)"},
	{keyword: "TRY_CONVERT", role: roleExprScalar, acceptSample: "SELECT TRY_CONVERT(INT, '1')"},
	{keyword: "COALESCE", role: roleExprScalar, acceptSample: "SELECT COALESCE(1, 2)"},
	{keyword: "NULLIF", role: roleExprScalar, acceptSample: "SELECT NULLIF(1, 2)"},
	{keyword: "IIF", role: roleExprScalar, acceptSample: "SELECT IIF(1 = 1, 1, 0)"},
	// --- roleFullTextPredicate ---
	{keyword: "CONTAINS", role: roleFullTextPredicate, acceptSample: "SELECT 1 FROM t WHERE CONTAINS(b, 'foo')"},
	{keyword: "FREETEXT", role: roleFullTextPredicate, acceptSample: "SELECT 1 FROM t WHERE FREETEXT(b, 'foo')"},
	// --- roleRowsetTable ---
	//
	// The accept samples below are the minimal legal shape that exercises
	// each keyword's dispatch in parsePrimaryTableSource. Richer forms
	// (OPENROWSET BULK, OPENDATASOURCE.db.schema.tbl chains) are separately
	// tested elsewhere and are out of scope for this regression guard.
	{keyword: "OPENROWSET", role: roleRowsetTable, acceptSample: "SELECT * FROM OPENROWSET('SQLNCLI', 'srv', 'SELECT 1')"},
	{keyword: "OPENQUERY", role: roleRowsetTable, acceptSample: "SELECT * FROM OPENQUERY(srv, 'SELECT 1')"},
	// OPENJSON is a ContextKeyword in SqlScriptDOM: scalar `SELECT
	// OPENJSON(...)` parses as a generic FunctionCall there, so omni
	// mirrors that behavior. Special-path dispatch in FROM.
	{keyword: "OPENJSON", role: roleRowsetTable, acceptSample: "SELECT * FROM OPENJSON(@j)", alsoScalar: true},
	{keyword: "OPENDATASOURCE", role: roleRowsetTable, acceptSample: "SELECT * FROM OPENDATASOURCE('SQLNCLI', 'conn')"},
	{keyword: "OPENXML", role: roleRowsetTable, acceptSample: "SELECT * FROM OPENXML(@h, 'x')"},
	// --- roleFullTextTable ---
	{keyword: "CONTAINSTABLE", role: roleFullTextTable, acceptSample: "SELECT * FROM CONTAINSTABLE(t, b, 'foo') ct"},
	{keyword: "FREETEXTTABLE", role: roleFullTextTable, acceptSample: "SELECT * FROM FREETEXTTABLE(t, b, 'foo') ft"},
	// --- roleSemanticTable ---
	{keyword: "SEMANTICKEYPHRASETABLE", role: roleSemanticTable, acceptSample: "SELECT * FROM SEMANTICKEYPHRASETABLE(t, b) s"},
	{keyword: "SEMANTICSIMILARITYTABLE", role: roleSemanticTable, acceptSample: "SELECT * FROM SEMANTICSIMILARITYTABLE(t, b, 1) s"},
	{keyword: "SEMANTICSIMILARITYDETAILSTABLE", role: roleSemanticTable, acceptSample: "SELECT * FROM SEMANTICSIMILARITYDETAILSTABLE(t, a, 1, b, 2) s"},
}

// rejectContexts are the "wrong-context" shapes we spray each keyword into
// to prove the parser rejects it outside its declared role.
type rejectContext struct {
	name   string
	render func(keyword string) string
	// allowRoles is the set of roles that are allowed to use this context.
	// A keyword whose role is in this set is skipped for this context
	// (because the context matches its declared role).
	allowRoles []keywordRole
}

var rejectContexts = []rejectContext{
	{
		name:   "scalar-select-list",
		render: func(k string) string { return "SELECT " + k + "(1, 'x')" },
		// CAST-family produce their own statement-level scalar
		// expressions; full-text predicates and rowset/table functions
		// must not appear here.
		allowRoles: []keywordRole{roleExprScalar},
	},
	{
		name:   "scalar-func-arg",
		render: func(k string) string { return "SELECT COALESCE(" + k + "(a, 'x'), 0)" },
		// Wrapping in another scalar function is still scalar position —
		// only expr-scalar callables are legal here (and CAST/CONVERT
		// take special-shaped args not `a, 'x'`; the wrapping test is
		// intentionally permissive for scalar callables so we skip
		// failure assertion for roleExprScalar here by listing it as
		// "allowed").
		allowRoles: []keywordRole{roleExprScalar},
	},
	{
		name:   "from-table-source",
		render: func(k string) string { return "SELECT * FROM " + k + "(1, 2)" },
		allowRoles: []keywordRole{
			roleRowsetTable,
			roleFullTextTable,
			roleSemanticTable,
			// Note: expr-scalar callables have degenerate function-call
			// syntax and would be treated as TVFs by the FROM path. We
			// still expect them to fail here because the args don't
			// match their expected grammar — except COALESCE/NULLIF/IIF
			// which are variadic-ish. Skip CAST/CONVERT/TRY_CAST/
			// TRY_CONVERT via allow for roleExprScalar — COALESCE etc.
			// failing here is fine even if allowed, we only assert
			// *rejections* for keywords NOT in allowRoles.
			roleExprScalar,
		},
	},
}

// TestCallableKeywordManifest asserts each manifest entry parses in its
// declared accepting context.
func TestCallableKeywordManifest(t *testing.T) {
	for _, entry := range callableKeywordManifest {
		t.Run(entry.keyword+"/accept", func(t *testing.T) {
			_, err := Parse(entry.acceptSample)
			if err != nil {
				t.Fatalf("accept-context parse failed for %s: %v\nSQL: %s",
					entry.keyword, err, entry.acceptSample)
			}
		})
	}
}

// TestCallableKeywordRoleRejection asserts each keyword is rejected in
// every context NOT in its declared role. This catches future drift where
// a refactor either widens a dispatch silently (letting a predicate into
// scalar position) or narrows one and breaks an accepted keyword.
func TestCallableKeywordRoleRejection(t *testing.T) {
	for _, entry := range callableKeywordManifest {
		for _, ctx := range rejectContexts {
			if roleIn(entry.role, ctx.allowRoles) {
				continue
			}
			// Hybrid scalar/FROM keywords (OPENJSON today) are also
			// legal as generic scalar FunctionCalls.
			if entry.alsoScalar && strings.HasPrefix(ctx.name, "scalar-") {
				continue
			}
			sql := ctx.render(entry.keyword)
			name := entry.keyword + "/" + ctx.name
			t.Run(name, func(t *testing.T) {
				_, err := Parse(sql)
				if err == nil {
					t.Errorf("keyword %s must be rejected in context %s\nSQL: %s",
						entry.keyword, ctx.name, sql)
				}
			})
		}
	}
}

// TestCallableKeywordManifestIsExhaustive enforces that every CoreKeyword
// whose lowercase form ends in a "callable-like" suffix (`table`) or is in
// the hand-picked tripwire list is listed in callableKeywordManifest. This
// is the drift detector: if someone adds a new `*TABLE` CoreKeyword and
// forgets to wire it, this test tells them where to look.
func TestCallableKeywordManifestIsExhaustive(t *testing.T) {
	// Build a lookup from the manifest.
	have := map[string]bool{}
	for _, e := range callableKeywordManifest {
		have[strings.ToLower(e.keyword)] = true
	}
	// Hand-picked tripwire list — historically all of these were
	// CoreKeyword function-call-like tokens. If a new keyword fits any
	// of these patterns add it to the manifest.
	tripwire := []string{}
	for _, kw := range sqlServerCoreKeywords {
		lw := strings.ToLower(kw)
		// Tripwire 1: every "*table" keyword in SqlScriptDOM is a TVF.
		if strings.HasSuffix(lw, "table") && lw != "table" {
			tripwire = append(tripwire, lw)
		}
		// Tripwire 2: the OPEN* family.
		if strings.HasPrefix(lw, "open") && lw != "open" {
			tripwire = append(tripwire, lw)
		}
	}
	for _, kw := range tripwire {
		if !have[kw] {
			t.Errorf("CoreKeyword %q fits a known callable pattern but is not in callableKeywordManifest. "+
				"Either add a manifest entry with the correct role or document why it is not a callable.",
				kw)
		}
	}
}

func roleIn(r keywordRole, set []keywordRole) bool {
	for _, x := range set {
		if x == r {
			return true
		}
	}
	return false
}
