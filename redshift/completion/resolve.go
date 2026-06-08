package completion

import (
	"strings"

	"github.com/bytebase/omni/redshift/ast"
	"github.com/bytebase/omni/redshift/catalog"
	"github.com/bytebase/omni/redshift/parser"
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
	case "relation_expr", "qualified_name":
		return resolveRelations(cat, sql, offset)
	case "func_name":
		return resolveFunctions(cat)
	case "any_name":
		// any_name appears in contexts like COMMENT ON COLUMN schema.table.col
		// where the name could be a column, table, or schema. Resolve all.
		var result []Candidate
		result = append(result, resolveColumns(cat, sql, offset)...)
		result = append(result, resolveRelations(cat, sql, offset)...)
		return result
	}
	return nil
}

func resolveRelations(cat *catalog.Catalog, sql string, offset int) []Candidate {
	if schema := qualifierBeforeCursor(sql, offset); schema != "" {
		return resolveRelationsInSchema(cat, schema)
	}

	var result []Candidate
	for _, s := range cat.UserSchemas() {
		result = append(result, Candidate{Text: s.Name, Type: CandidateSchema})
		result = append(result, relationCandidatesForSchema(s)...)
	}
	// Include CTE names from the query as table candidates.
	refs := extractTableRefs(sql, offset)
	for _, ref := range refs {
		if ref.Schema == "" && cat.GetRelation("", ref.Table) == nil {
			// This ref is not a real table — likely a CTE.
			result = append(result, Candidate{Text: ref.Table, Type: CandidateTable})
		}
	}
	return result
}

func resolveRelationsInSchema(cat *catalog.Catalog, schemaName string) []Candidate {
	for _, s := range cat.UserSchemas() {
		if strings.EqualFold(s.Name, schemaName) {
			return relationCandidatesForSchema(s)
		}
	}
	return nil
}

func relationCandidatesForSchema(s *catalog.Schema) []Candidate {
	var result []Candidate
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
	return result
}

func qualifierBeforeCursor(sql string, offset int) string {
	if offset < 0 {
		offset = 0
	}
	if offset > len(sql) {
		offset = len(sql)
	}
	if offset == 0 || sql[offset-1] != '.' {
		return ""
	}
	i := offset - 1
	for i > 0 {
		c := sql[i-1]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') || c == '_' {
			i--
			continue
		}
		break
	}
	return sql[i : offset-1]
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
	result = append(result, resolveVirtualColumns(sql, offset)...)
	return result
}

func resolveVirtualColumns(sql string, offset int) []Candidate {
	ctx := parser.CollectCompletion(sql, offset)
	if ctx == nil || ctx.Scope == nil {
		return nil
	}

	var result []Candidate
	for _, ref := range ctx.Scope.References {
		result = append(result, columnCandidates(ref.AliasColumns)...)
		result = append(result, columnCandidates(syntaxOutputColumns(sql, ref))...)
	}
	return result
}

func columnCandidates(cols []string) []Candidate {
	result := make([]Candidate, 0, len(cols))
	for _, col := range cols {
		result = append(result, Candidate{Text: col, Type: CandidateColumn})
	}
	return result
}

func syntaxOutputColumns(sql string, ref parser.RangeReference) []string {
	if len(ref.AliasColumns) > 0 {
		return nil
	}
	if ref.Kind != parser.RangeReferenceSubquery && ref.Kind != parser.RangeReferenceCTE {
		return nil
	}
	if ref.BodyLoc.Start < 0 || ref.BodyLoc.End <= ref.BodyLoc.Start || ref.BodyLoc.End > len(sql) {
		return nil
	}

	list, err := parser.Parse(sql[ref.BodyLoc.Start:ref.BodyLoc.End])
	if err != nil || list == nil || len(list.Items) == 0 {
		return nil
	}
	raw, ok := list.Items[0].(*ast.RawStmt)
	if !ok {
		return nil
	}
	stmt, ok := raw.Stmt.(*ast.SelectStmt)
	if !ok || stmt.TargetList == nil {
		return nil
	}

	var result []string
	for _, item := range stmt.TargetList.Items {
		target, ok := item.(*ast.ResTarget)
		if !ok {
			continue
		}
		if target.Name != "" {
			result = append(result, target.Name)
			continue
		}
		if col := columnRefOutputName(target.Val); col != "" {
			result = append(result, col)
		}
	}
	return result
}

func columnRefOutputName(n ast.Node) string {
	ref, ok := n.(*ast.ColumnRef)
	if !ok || ref.Fields == nil || len(ref.Fields.Items) == 0 {
		return ""
	}
	last := ref.Fields.Items[len(ref.Fields.Items)-1]
	name, ok := last.(*ast.String)
	if !ok {
		return ""
	}
	return name.Str
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
