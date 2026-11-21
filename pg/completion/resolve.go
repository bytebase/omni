package completion

import (
	"strings"

	"github.com/bytebase/omni/pg/ast"
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

	// Words are clause text no keyword names, such as PARTITIONED BY.
	for _, word := range cs.Words {
		result = append(result, Candidate{Text: word, Type: CandidateKeyword})
	}

	// Rule candidates -> catalog objects and statement-local scopes (CTEs/aliases).
	for _, rc := range cs.Rules {
		result = append(result, resolveRule(rc.Rule, cat, sql, cursorOffset)...)
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
		if cat == nil {
			return nil
		}
		return resolveRelationsInSchema(cat, schema)
	}

	var result []Candidate
	if cat != nil {
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
	}
	if ctx := parser.CollectCompletion(sql, offset); ctx != nil {
		for _, cte := range ctx.CTEs {
			if cte.Name != "" && (cat == nil || cat.GetRelation("", cte.Name) == nil) {
				result = append(result, Candidate{Text: cte.Name, Type: CandidateTable})
			}
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
		var ct CandidateType
		switch rel.RelKind {
		case 'r', 'p', 'f':
			ct = CandidateTable
		case 'v':
			ct = CandidateView
		case 'm':
			ct = CandidateMaterializedView
		default:
			continue
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
	if i > 0 && sql[i-1] == '"' {
		end := i - 1
		i -= 2
		for i >= 0 {
			if sql[i] == '"' {
				if i > 0 && sql[i-1] == '"' {
					i -= 2
					continue
				}
				return strings.ReplaceAll(sql[i+1:end], `""`, `"`)
			}
			i--
		}
		return ""
	}
	for i > 0 && isPGIdentCont(sql[i-1]) {
		i--
	}
	if i == offset-1 || !isPGIdentStart(sql[i]) {
		return ""
	}
	return sql[i : offset-1]
}

func isPGIdentStart(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' || c >= 128
}

func isPGIdentCont(c byte) bool {
	return isPGIdentStart(c) || (c >= '0' && c <= '9') || c == '$'
}

func tableRefMatchesQualifier(ref TableRef, qual string) bool {
	if qual == "" {
		return true
	}
	if ref.Alias != "" {
		return strings.EqualFold(ref.Alias, qual)
	}
	return strings.EqualFold(ref.Table, qual)
}

func resolveColumns(cat *catalog.Catalog, sql string, offset int) []Candidate {
	qual := qualifierBeforeCursor(sql, offset)
	var result []Candidate
	if cat != nil {
		refs := extractTableRefs(sql, offset)
		for _, ref := range refs {
			if !tableRefMatchesQualifier(ref, qual) {
				continue
			}
			rel := cat.GetRelation(ref.Schema, ref.Table)
			if rel == nil {
				continue
			}
			for _, col := range rel.Columns {
				result = append(result, Candidate{Text: col.Name, Type: CandidateColumn})
			}
		}
	}
	result = append(result, resolveVirtualColumns(sql, offset, qual)...)
	return result
}

func resolveVirtualColumns(sql string, offset int, qual string) []Candidate {
	ctx := parser.CollectCompletion(sql, offset)
	if ctx == nil || ctx.Scope == nil {
		return nil
	}

	var result []Candidate
	for _, ref := range ctx.Scope.References {
		if qual != "" {
			refName := ref.Alias
			if refName == "" {
				refName = ref.Name
			}
			if !strings.EqualFold(refName, qual) {
				continue
			}
		}
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
	if !ok {
		return nil
	}
	for stmt != nil && stmt.TargetList == nil && stmt.Larg != nil {
		stmt = stmt.Larg
	}
	if stmt == nil || stmt.TargetList == nil {
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
	if cat == nil {
		return nil
	}
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
