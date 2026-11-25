package completion

import (
	"slices"
	"strings"
)

// patchTails complete a statement the caret leaves unfinished, such as
// UPDATE t SET | or INSERT INTO |, so the placeholder can be parsed in place.
var patchTails = []string{
	"",
	" TRUE",
	" = 1",
	" VALUES (1)",
	" SET omni_completion_ident = 1",
	" ON TRUE",
	" (omni_completion_ident)",
	" AS omni_completion_ident",
}

// nameContextAt returns what the word at the caret names, or nil when no
// name fits there.
func nameContextAt(stmt string, word caretWord, cat *catalog) *nameContext {
	prefix := stmt[:word.start] + placeholder
	suffix := stmt[word.end:]
	for _, text := range patchedTexts(prefix, suffix) {
		parsed, err := parseOne(text)
		if err != nil {
			continue
		}
		nc := analyzeStatement(parsed, cat)
		if nc == nil {
			// The placeholder parsed as a name completion has no objects
			// for, such as a new alias or column definition.
			return nil
		}
		if (nc.kind == kindExpr || nc.kind == kindFunction) && len(nc.scope) == 0 {
			nc.scope = lexicalSources(stmt)
		}
		return nc
	}
	return lexicalContext(stmt, word)
}

// patchedTexts returns the statement with the placeholder at the caret, then
// the same text completed by each tail and the closing parentheses it lacks,
// then the text before the caret completed the same way.
func patchedTexts(prefix, suffix string) []string {
	var out []string
	add := func(text string) {
		if !slices.Contains(out, text) {
			out = append(out, text)
		}
	}
	for _, head := range []string{prefix + suffix, prefix} {
		head = strings.TrimRight(head, " \t\r\n;")
		for _, tail := range patchTails {
			add(head + tail + closers(head+tail))
			add(head + closers(head) + tail)
		}
	}
	return out
}

// closers returns the ")" and "]" that close what text leaves open.
func closers(text string) string {
	var open []byte
	for _, t := range tokenize(text) {
		switch t.id {
		case '(':
			open = append(open, ')')
		case '[':
			open = append(open, ']')
		case ')', ']':
			if len(open) > 0 {
				open = open[:len(open)-1]
			}
		}
	}
	slices.Reverse(open)
	return string(open)
}

// nameCandidates lists the objects nc names. Hidden columns, such as the
// rowid CockroachDB adds to a table without a primary key, are offered only
// when the user has started typing a name.
func nameCandidates(nc *nameContext, cat *catalog, typing bool) []Candidate {
	if nc == nil {
		return nil
	}
	var out []Candidate
	switch nc.kind {
	case kindObject:
		if len(nc.qualifier) > 0 {
			return relationCandidates(cat.qualifiedSchema(nc.qualifier[len(nc.qualifier)-1]))
		}
		if cat != nil {
			for _, s := range cat.schemas {
				out = append(out, Candidate{Text: s.name, Type: CandidateSchema})
				out = append(out, relationCandidates(s)...)
			}
		}
		for _, cte := range nc.ctes {
			if cat.relation("", cte) == nil {
				out = append(out, Candidate{Text: cte, Type: CandidateTable})
			}
		}
	case kindSchema:
		if cat != nil {
			for _, s := range cat.schemas {
				out = append(out, Candidate{Text: s.name, Type: CandidateSchema})
			}
		}
	case kindExpr, kindFunction:
		a := &analyzer{cat: cat}
		for _, src := range nc.scope {
			if matchesQualifier(src, nc.qualifier) {
				out = append(out, columnCandidates(a.sourceColumns(src, typing))...)
			}
		}
		if len(nc.qualifier) == 0 && cat != nil {
			for _, s := range cat.schemas {
				for _, fn := range s.functions {
					out = append(out, Candidate{Text: fn, Type: CandidateFunction})
				}
			}
		}
	case kindTargetColumn:
		a := &analyzer{cat: cat}
		out = columnCandidates(a.sourceColumns(*nc.target, typing))
	case kindIndex:
		if nc.target != nil && !nc.target.derived {
			if rel := cat.relation(nc.target.schema, nc.target.name); rel != nil {
				for _, idx := range rel.indexes {
					out = append(out, Candidate{Text: idx, Type: CandidateIndex})
				}
			}
		}
	}
	return out
}

// qualifiedSchema resolves the qualifier of an object name: a schema, or the
// current database, which CockroachDB reads as db.public.
func (c *catalog) qualifiedSchema(name string) *schema {
	if c == nil {
		return nil
	}
	if s := c.schema(name); s != nil {
		return s
	}
	if name == c.database {
		return c.schema("public")
	}
	return nil
}

func relationCandidates(s *schema) []Candidate {
	if s == nil {
		return nil
	}
	out := make([]Candidate, 0, len(s.relations))
	for _, r := range s.relations {
		out = append(out, Candidate{Text: r.name, Type: r.kind})
	}
	return out
}

func columnCandidates(cols []string) []Candidate {
	out := make([]Candidate, 0, len(cols))
	for _, c := range cols {
		out = append(out, Candidate{Text: c, Type: CandidateColumn})
	}
	return out
}

// matchesQualifier reports whether a column qualified by qual can come from
// src: t names the alias or table, s.t the schema-qualified table.
func matchesQualifier(src source, qual []string) bool {
	switch len(qual) {
	case 0:
		return true
	case 1:
		return src.ref == qual[0] || strings.EqualFold(src.ref, qual[0])
	default:
		t, s := qual[len(qual)-1], qual[len(qual)-2]
		return !src.derived && strings.EqualFold(src.name, t) && (src.schema == "" || strings.EqualFold(src.schema, s))
	}
}
