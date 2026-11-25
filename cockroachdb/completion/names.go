package completion

import (
	"context"
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
func nameContextAt(ctx context.Context, stmt string, word caretWord, cat *catalog) *nameContext {
	prefix := stmt[:word.start] + placeholder
	suffix := stmt[word.end:]
	for _, p := range patchedTexts(prefix, suffix) {
		if ctx.Err() != nil {
			break
		}
		parsed, err := parseOne(p.text)
		if err != nil {
			continue
		}
		nc := analyzeStatement(ctx, parsed, cat)
		if nc == nil {
			// The placeholder parsed as a name completion has no objects
			// for, such as a new alias or column definition.
			return nil
		}
		if (nc.kind == kindExpr || nc.kind == kindFunction) && (!nc.scoped || !p.withSuffix) {
			// The tree lacks the clauses the patch dropped, or the
			// analyzer did not walk to the name: read the tables the
			// statement names from its tokens as well.
			nc.scope = append(nc.scope, lexicalSources(stmt)...)
		}
		return nc
	}
	return lexicalContext(stmt, word)
}

// patchedText is one text nameContextAt tries to parse.
type patchedText struct {
	text       string
	withSuffix bool // the text keeps everything after the caret
}

// patchedTexts returns the statement with the placeholder at the caret, then
// the same text completed by each tail and the closing parentheses it lacks,
// then the text before the caret completed the same way.
func patchedTexts(prefix, suffix string) []patchedText {
	var out []patchedText
	add := func(text string, withSuffix bool) {
		for _, p := range out {
			if p.text == text {
				return
			}
		}
		out = append(out, patchedText{text, withSuffix})
	}
	for _, head := range []struct {
		text       string
		withSuffix bool
	}{{prefix + suffix, true}, {prefix, false}} {
		text := strings.TrimRight(head.text, " \t\r\n;")
		for _, tail := range patchTails {
			add(text+tail+closers(text+tail), head.withSuffix)
			add(text+closers(text)+tail, head.withSuffix)
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

// nameCandidates lists the objects nc names whose names start with typed,
// ignoring case. Schema, relation, column and index texts are final:
// quoted as CockroachDB reads them back, or only escaped when quoted (the
// caret is inside a double quote). Function and type texts are left to
// formatNames. Hidden columns, such as the rowid CockroachDB adds to a table
// without a primary key, are offered only when typed is not empty.
func nameCandidates(nc *nameContext, cat *catalog, typed string, quoted bool) []Candidate {
	if nc == nil {
		return nil
	}
	var out []Candidate
	add := func(name string, typ CandidateType) {
		if hasPrefixFold(name, typed) {
			out = append(out, Candidate{Text: QuoteIdentifierIfNeeded(name, quoted), Type: typ})
		}
	}
	addRaw := func(name string, typ CandidateType) {
		if hasPrefixFold(name, typed) {
			out = append(out, Candidate{Text: name, Type: typ})
		}
	}
	typing := typed != ""
	switch nc.kind {
	case kindObject:
		if len(nc.qualifier) > 0 {
			if s := cat.qualifiedSchema(nc.qualifier[len(nc.qualifier)-1]); s != nil {
				for _, r := range s.relations {
					add(r.name, r.kind)
				}
			}
			return out
		}
		if cat != nil {
			for _, s := range cat.schemas {
				add(s.name, CandidateSchema)
				for _, r := range s.relations {
					if cat.relation("", r.name) == r {
						add(r.name, r.kind)
					} else if !quoted && hasPrefixFold(r.name, typed) {
						// Outside the search path, or behind a table of the
						// same name earlier on it, a relation needs its schema.
						out = append(out, Candidate{Text: QuoteIdentifierIfNeeded(s.name, false) + "." + QuoteIdentifierIfNeeded(r.name, false), Type: r.kind})
					}
				}
			}
		}
		for _, cte := range nc.ctes {
			if cat.relation("", cte) == nil {
				add(cte, CandidateTable)
			}
		}
	case kindSchema:
		if cat != nil {
			for _, s := range cat.schemas {
				add(s.name, CandidateSchema)
			}
		}
	case kindExpr, kindFunction:
		a := &analyzer{ctx: context.Background(), cat: cat}
		for _, src := range resolveQualifier(nc.scope, nc.qualifier) {
			for _, col := range a.sourceColumns(src, typing) {
				add(col, CandidateColumn)
			}
		}
		builtins := slices.Concat(builtinFunctions, specialForms())
		for _, fn := range schemaObjects(cat, nc.qualifier, builtins, func(s *schema) []string { return s.functions }) {
			addRaw(fn, CandidateFunction)
		}
	case kindTargetColumn:
		a := &analyzer{ctx: context.Background(), cat: cat}
		for _, col := range a.sourceColumns(*nc.target, typing) {
			add(col, CandidateColumn)
		}
	case kindColumnList:
		for _, col := range nc.columns {
			add(col, CandidateColumn)
		}
	case kindType:
		builtins := slices.Concat(builtinTypes, typeKeywords())
		for _, typ := range schemaObjects(cat, nc.qualifier, builtins, func(s *schema) []string { return s.types }) {
			addRaw(typ, CandidateType_)
		}
	case kindIndex:
		if nc.target != nil && !nc.target.derived {
			if rel := cat.relation(nc.target.schema, nc.target.name); rel != nil {
				for _, idx := range rel.indexes {
					add(idx, CandidateIndex)
				}
			}
		}
	}
	return out
}

// resolveQualifier returns the sources a column qualified by qual can come
// from. A qualifier resolves in the innermost query that defines it, so an
// inner alias hides an outer one of the same name.
func resolveQualifier(scope []source, qual []string) []source {
	if len(qual) == 0 {
		return scope
	}
	var out []source
	depth := -1
	for _, src := range scope {
		if !matchesQualifier(src, qual) {
			continue
		}
		switch {
		case src.depth > depth:
			out, depth = []source{src}, src.depth
		case src.depth == depth:
			out = append(out, src)
		}
	}
	return out
}

func hasPrefixFold(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.EqualFold(s[:len(prefix)], prefix)
}

// schemaObjects lists functions or types a name qualified by qual can
// reach: unqualified, the built-ins and every schema's own; qualified by
// pg_catalog, the built-ins; qualified by a schema, that schema's own.
func schemaObjects(cat *catalog, qual []string, builtins []string, own func(*schema) []string) []string {
	switch len(qual) {
	case 0:
		out := slices.Clone(builtins)
		if cat != nil {
			for _, s := range cat.schemas {
				out = append(out, own(s)...)
			}
		}
		return out
	case 1:
		if strings.EqualFold(qual[0], "pg_catalog") {
			return builtins
		}
		if s := cat.qualifiedSchema(qual[0]); s != nil {
			return own(s)
		}
	}
	return nil
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
