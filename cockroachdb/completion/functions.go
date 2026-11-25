package completion

import (
	"slices"
	"strings"
	"sync"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/lexbase"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/types"
)

// specialForms are the keywords CockroachDB's grammar reads as a function
// call of its own syntax, such as COALESCE(a, b) or NULLIF(a, b): the call
// parses to a node other than a plain call of a named function. They are not
// in pg_proc, so they are offered as functions beside builtinFunctions.
var specialForms = sync.OnceValue(func() []string {
	var out []string
	for _, kw := range lexbase.KeywordNames {
		if lexbase.KeywordsCategories[kw] == "U" || lexerRewrites[kw] {
			continue
		}
		for _, args := range []string{"()", "(1)", "(1, 1)", "(1, 1, 1)"} {
			if isSpecialForm(kw + args) {
				out = append(out, kw)
				break
			}
		}
	}
	return out
})

func isSpecialForm(call string) bool {
	stmt, err := parseOne("SELECT " + call)
	if err != nil {
		return false
	}
	sel, ok := stmt.(*tree.Select)
	if !ok {
		return false
	}
	clause, ok := sel.Select.(*tree.SelectClause)
	if !ok || len(clause.Exprs) != 1 {
		return false
	}
	switch e := clause.Exprs[0].Expr.(type) {
	case *tree.FuncExpr:
		_, named := e.Func.FunctionReference.(*tree.UnresolvedName)
		return !named
	case *tree.CoalesceExpr, *tree.NullIfExpr, *tree.IfExpr, *tree.IfErrExpr:
		return true
	}
	return false
}

// typeKeywords are the keywords CockroachDB's grammar reads as a built-in
// type, such as INT, BOOLEAN or REAL: a cast to them parses to a resolved
// type rather than a type name to look up. They are offered as types beside
// builtinTypes, which holds the catalog's spellings.
var typeKeywords = sync.OnceValue(func() []string {
	var out []string
	for _, kw := range lexbase.KeywordNames {
		if lexerRewrites[kw] {
			continue
		}
		stmt, err := parseOne("SELECT NULL::" + kw)
		if err != nil {
			continue
		}
		sel, ok := stmt.(*tree.Select)
		if !ok {
			continue
		}
		clause, ok := sel.Select.(*tree.SelectClause)
		if !ok || len(clause.Exprs) != 1 {
			continue
		}
		cast, ok := clause.Exprs[0].Expr.(*tree.CastExpr)
		if !ok {
			continue
		}
		// Pseudo-types, which no value or column has, are left out as
		// builtinTypes leaves them out.
		if typ, ok := cast.Type.(*types.T); ok && typ.Family() != types.UnknownFamily && typ.Family() != types.TriggerFamily {
			out = append(out, strings.ToUpper(kw))
		}
	}
	return out
})

// formatNames writes the text of each function and type candidate. A
// built-in type keeps CockroachDB's spelling (INT8, STRING); a user-defined
// type is quoted like any name. A function is quoted when CockroachDB would
// not read the name back bare, with "()" when the caret starts a call. A
// name that is also a keyword is probed at the caret: it stays bare when
// "name(" parses there, as substring( and session_user( do in CockroachDB,
// and is quoted otherwise. A function or type with a schema is qualified by
// it. Keyword candidates a bare function covers are dropped.
func formatNames(prefix string, candidates []pendingName, quoted, call bool) []Candidate {
	base := ""
	out := make([]Candidate, 0, len(candidates))
	covered := map[string]bool{}
	for _, p := range candidates {
		c := p.Candidate
		qualifier := ""
		if p.schema != "" {
			qualifier = QuoteIdentifierIfNeeded(p.schema, false) + "."
		}
		switch c.Type {
		case CandidateType_:
			if qualifier != "" || !isBuiltinType(c.Text) {
				c.Text = qualifier + QuoteIdentifierIfNeeded(c.Text, quoted)
			}
		case CandidateFunction:
			name := QuoteIdentifierIfNeeded(c.Text, quoted)
			if _, keyword := lexbase.KeywordsCategories[c.Text]; keyword && !quoted && qualifier == "" {
				if base == "" {
					base = probePrefix(prefix) + " "
				}
				if reachesPast(base+c.Text+" (", len(base)+len(c.Text)+1) {
					name = c.Text
					covered[strings.ToUpper(c.Text)] = true
				}
			}
			c.Text = qualifier + name
			if call && !quoted {
				c.Text += "()"
			}
		}
		out = append(out, c)
	}
	return slices.DeleteFunc(out, func(c Candidate) bool {
		return c.Type == CandidateKeyword && covered[c.Text]
	})
}

func isBuiltinType(name string) bool {
	return slices.Contains(builtinTypes, name) || slices.Contains(typeKeywords(), name)
}
