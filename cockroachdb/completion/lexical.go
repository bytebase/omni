package completion

import (
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/lexbase"
)

// relationKeywords are the keywords after which a lexical scan expects a
// table name.
var relationKeywords = map[string]bool{
	"from": true, "join": true, "into": true, "update": true, "table": true, "truncate": true, "using": true,
}

// lexicalContext classifies the caret from the tokens before it, for a
// statement no patch makes parseable: after FROM, JOIN, INTO and the like it
// expects a table; elsewhere, if an identifier fits, a column of the tables
// the statement names.
func lexicalContext(stmt string, word caretWord) *nameContext {
	prefix := stmt[:word.start]
	toks := tokenize(prefix)
	if n := len(toks); n > 0 && toks[n-1].id == lexbase.ERROR {
		return nil
	}
	// Peel a qualifier chain such as "s.t." off the end.
	var qual []string
	i := len(toks)
	for i >= 2 && toks[i-1].id == '.' && isNameToken(toks[i-2]) {
		qual = append([]string{toks[i-2].str}, qual...)
		i -= 2
	}
	if i > 0 && isKeyword(toks[i-1], relationKeywords) {
		return &nameContext{kind: kindObject, qualifier: qual}
	}
	if !reachesPast(prefix+identProbe, len(prefix)) {
		return nil
	}
	return &nameContext{kind: kindExpr, qualifier: qual, scope: lexicalSources(stmt)}
}

// lexicalSources reads the tables a statement names after FROM, JOIN,
// UPDATE, INTO and USING, with their aliases, from its tokens alone.
func lexicalSources(stmt string) []source {
	toks := tokenize(stmt)
	var out []source
	for i := 0; i < len(toks); i++ {
		if !isKeyword(toks[i], relationKeywords) || toks[i].str == "table" {
			continue
		}
		list := toks[i].str == "from" || toks[i].str == "using"
		for j := i + 1; j < len(toks); {
			var parts []string
			for j < len(toks) && isNameToken(toks[j]) {
				parts = append(parts, toks[j].str)
				j++
				if j < len(toks) && toks[j].id == '.' {
					j++
					continue
				}
				break
			}
			if len(parts) == 0 {
				break
			}
			src := source{name: parts[len(parts)-1], ref: parts[len(parts)-1]}
			if len(parts) >= 2 {
				src.qual = parts[max(0, len(parts)-3) : len(parts)-1]
				src.home = parts[len(parts)-2]
			}
			if j < len(toks) && toks[j].str == "as" && toks[j].id != lexbase.IDENT {
				j++
			}
			if j < len(toks) && toks[j].id == lexbase.IDENT {
				src.ref = toks[j].str
				j++
			}
			out = append(out, src)
			if !list || j >= len(toks) || toks[j].id != ',' {
				break
			}
			j++
		}
	}
	return out
}

// isNameToken reports whether t can be part of an object name: an
// identifier or a keyword CockroachDB does not reserve.
func isNameToken(t token) bool {
	if t.id == lexbase.IDENT {
		return true
	}
	cat, ok := lexbase.KeywordsCategories[t.str]
	return ok && cat != "R" && lexbase.GetKeywordID(t.str) == t.id
}

func isKeyword(t token, set map[string]bool) bool {
	return t.id != lexbase.IDENT && set[t.str]
}
