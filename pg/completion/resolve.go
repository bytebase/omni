package completion

import (
	"strings"

	"github.com/bytebase/omni/pg/catalog"
	"github.com/bytebase/omni/pg/parser"
)

func resolve(cs *parser.CandidateSet, cat *catalog.Catalog, sql string, cursorOffset int) []Candidate {
	if cs == nil {
		return nil
	}
	var result []Candidate

	// Token candidates -> keywords
	for _, tok := range cs.Tokens {
		name := parser.TokenName(tok)
		if name == "" {
			continue
		}
		result = append(result, Candidate{Text: name, Type: CandidateKeyword})
	}

	// Rule candidates -> catalog objects
	if cat != nil {
		for _, rc := range cs.Rules {
			result = append(result, resolveRule(rc.Rule, cat, sql, cursorOffset)...)
		}
	}

	return dedup(result)
}

func resolveRule(rule string, cat *catalog.Catalog, sql string, offset int) []Candidate {
	switch rule {
	case "columnref":
		return resolveColumns(cat, sql, offset)
	case "relation_expr":
		return resolveRelations(cat)
	case "func_name":
		return resolveFunctions(cat)
	}
	return nil
}

func resolveRelations(cat *catalog.Catalog) []Candidate {
	var result []Candidate
	for _, s := range cat.UserSchemas() {
		result = append(result, Candidate{Text: s.Name, Type: CandidateSchema})
		for name, rel := range s.Relations {
			ct := CandidateTable
			switch rel.RelKind {
			case 'v':
				ct = CandidateView
			case 'm':
				ct = CandidateMaterializedView
			}
			result = append(result, Candidate{Text: name, Type: ct})
		}
		for name := range s.Sequences {
			result = append(result, Candidate{Text: name, Type: CandidateSequence})
		}
	}
	return result
}

func resolveColumns(cat *catalog.Catalog, sql string, offset int) []Candidate {
	refs := extractTableRefs(sql, offset)
	var result []Candidate
	for _, ref := range refs {
		rel := cat.GetRelation(ref.Schema, ref.Table)
		if rel == nil {
			continue
		}
		for _, col := range rel.Columns {
			result = append(result, Candidate{Text: col.Name, Type: CandidateColumn})
		}
	}
	return result
}

func resolveFunctions(cat *catalog.Catalog) []Candidate {
	names := cat.AllProcNames()
	result := make([]Candidate, 0, len(names))
	for _, name := range names {
		result = append(result, Candidate{Text: name, Type: CandidateFunction})
	}
	return result
}

func dedup(cs []Candidate) []Candidate {
	type key struct {
		text string
		typ  CandidateType
	}
	seen := make(map[key]bool, len(cs))
	result := make([]Candidate, 0, len(cs))
	for _, c := range cs {
		k := key{strings.ToLower(c.Text), c.Type}
		if seen[k] {
			continue
		}
		seen[k] = true
		result = append(result, c)
	}
	return result
}
