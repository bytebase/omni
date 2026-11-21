package deparse

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/trino/analysis"
	"github.com/bytebase/omni/trino/ast"
	trinoparser "github.com/bytebase/omni/trino/parser"
)

// DatabaseDefinitionFromProto renders a Trino database schema proto as SDL.
func DatabaseDefinitionFromProto(metadata *metadatapb.DatabaseSchemaMetadata) (string, error) {
	converted, err := fromProtoMetadata(metadata)
	if err != nil {
		return "", err
	}
	return GetDatabaseDefinition(converted)
}

func fromProtoMetadata(metadata *metadatapb.DatabaseSchemaMetadata) (*DatabaseSchemaMetadata, error) {
	if metadata == nil {
		return nil, nil
	}
	if metadata.CharacterSet != "" || metadata.Collation != "" || metadata.Owner != "" || metadata.SearchPath != "" || metadata.ServiceName != "" || metadata.DefaultSchema != "" || metadata.Datashare || len(metadata.Extensions) > 0 || slices.ContainsFunc(metadata.EventTriggers, func(et *metadatapb.EventTriggerMetadata) bool { return et != nil && !et.SkipDump }) || len(metadata.PublicSynonyms) > 0 || len(metadata.LinkedDatabases) > 0 || len(metadata.ForeignServers) > 0 || len(metadata.ForeignDataWrappers) > 0 {
		return nil, fmt.Errorf("database %s contains database attributes or objects unsupported by Trino SDL renderer", quoteIdentifier(metadata.Name))
	}
	out := &DatabaseSchemaMetadata{
		Name:    metadata.Name,
		Schemas: make([]*SchemaMetadata, 0, len(metadata.Schemas)),
	}
	seenSchemas := make(map[string]struct{}, len(metadata.Schemas))
	for _, s := range metadata.Schemas {
		if s == nil {
			continue
		}
		if s.SkipDump {
			if len(s.Functions) > 0 {
				skippedFns := make([]*FunctionMetadata, 0, len(s.Functions))
				for _, fn := range s.Functions {
					if fn != nil {
						skippedFns = append(skippedFns, convertSkippedTrinoFunction(fn))
					}
				}
				if len(skippedFns) > 0 {
					out.Schemas = append(out.Schemas, &SchemaMetadata{
						Name:      s.Name,
						Functions: skippedFns,
						SkipDump:  true,
					})
				}
			}
			continue
		}
		schemaKey := strings.ToLower(s.Name)
		if _, exists := seenSchemas[schemaKey]; exists {
			return nil, fmt.Errorf("database %s has duplicate schema %s", quoteIdentifier(metadata.Name), quoteIdentifier(s.Name))
		}
		seenSchemas[schemaKey] = struct{}{}
		if s.Comment != "" {
			return nil, fmt.Errorf("schema %s comment cannot be restored because Trino does not support COMMENT ON SCHEMA", s.Name)
		}
		if s.Owner != "" {
			return nil, fmt.Errorf("schema %s owner %q cannot be restored unambiguously because Trino SET AUTHORIZATION requires user vs role distinction", s.Name, s.Owner)
		}
		schemaMeta := &SchemaMetadata{
			Name:              s.Name,
			Functions:         make([]*FunctionMetadata, 0, len(s.Functions)),
			Tables:            make([]*TableMetadata, 0, len(s.Tables)),
			Views:             make([]*ViewMetadata, 0, len(s.Views)),
			MaterializedViews: make([]*ViewMetadata, 0, len(s.MaterializedViews)),
		}
		seenFuncSigs := make(map[string]struct{}, len(s.Functions))
		for _, fn := range s.Functions {
			if fn == nil {
				continue
			}
			if fn.SkipDump {
				schemaMeta.Functions = append(schemaMeta.Functions, convertSkippedTrinoFunction(fn))
				continue
			}
			def, sigKey, err := prepareTrinoFunctionDefinition(metadata.Name, s.Name, fn)
			if err != nil {
				return nil, err
			}
			if _, exists := seenFuncSigs[sigKey]; exists {
				return nil, fmt.Errorf("schema %s has duplicate function signature %s", s.Name, sigKey)
			}
			seenFuncSigs[sigKey] = struct{}{}
			schemaMeta.Functions = append(schemaMeta.Functions, &FunctionMetadata{
				Name:       fn.Name,
				Definition: def,
				Comment:    fn.Comment,
			})
		}
		for _, p := range s.Procedures {
			if p == nil || p.SkipDump {
				continue
			}
			return nil, fmt.Errorf("procedure %s.%s cannot be restored because Trino does not support restoring procedures", s.Name, p.Name)
		}
		for _, seq := range s.Sequences {
			if seq != nil && !seq.SkipDump {
				return nil, fmt.Errorf("sequence %s.%s cannot be restored because Trino does not support sequences", s.Name, seq.Name)
			}
		}
		for _, et := range s.ExternalTables {
			if et != nil && !et.SkipDump {
				return nil, fmt.Errorf("external table %s.%s cannot be restored by Trino SDL renderer", s.Name, et.Name)
			}
		}
		if len(s.Streams) > 0 || len(s.Tasks) > 0 || len(s.Packages) > 0 || len(s.Events) > 0 || len(s.Synonyms) > 0 ||
			slices.ContainsFunc(s.EnumTypes, func(e *metadatapb.EnumTypeMetadata) bool { return e != nil && !e.SkipDump }) ||
			slices.ContainsFunc(s.CompositeTypes, func(c *metadatapb.CompositeTypeMetadata) bool { return c != nil && !c.SkipDump }) ||
			slices.ContainsFunc(s.Collations, func(c *metadatapb.CollationMetadata) bool { return c != nil && !c.SkipDump }) ||
			slices.ContainsFunc(s.RangeTypes, func(r *metadatapb.RangeTypeMetadata) bool { return r != nil && !r.SkipDump }) {
			return nil, fmt.Errorf("schema %s contains schema objects unsupported by Trino SDL renderer", s.Name)
		}
		seenRelations := make(map[string]struct{}, len(s.Tables)+len(s.Views)+len(s.MaterializedViews))
		for _, t := range s.Tables {
			if t == nil || t.SkipDump {
				continue
			}
			relKey := strings.ToLower(t.Name)
			if _, exists := seenRelations[relKey]; exists {
				return nil, fmt.Errorf("schema %s has duplicate table %s", s.Name, quoteIdentifier(t.Name))
			}
			seenRelations[relKey] = struct{}{}
			if len(t.Indexes) > 0 || len(t.ForeignKeys) > 0 || len(t.CheckConstraints) > 0 || len(t.ExcludeConstraints) > 0 || len(t.Partitions) > 0 || len(t.Triggers) > 0 || len(t.Rules) > 0 || t.CreateOptions != "" || t.Engine != "" || t.Collation != "" || t.Charset != "" || len(t.SortingKeys) > 0 || t.ShardingInfo != "" || t.PrimaryKeyType != "" {
				return nil, fmt.Errorf("table %s.%s constraints or table options cannot be restored by Trino SDL renderer", s.Name, t.Name)
			}
			if t.Owner != "" {
				return nil, fmt.Errorf("table %s.%s owner %q cannot be restored unambiguously because Trino SET AUTHORIZATION requires user vs role distinction", s.Name, t.Name, t.Owner)
			}
			cols, err := fromProtoColumns(s.Name, t.Name, t.Columns, true)
			if err != nil {
				return nil, err
			}
			// GetDatabaseDefinition keeps emitting empty CREATE TABLE () for legacy
			// direct-struct oracle/golden parity, while DatabaseDefinitionFromProto
			// rejects columnless tables because Trino cannot restore them.
			if len(cols) == 0 {
				return nil, fmt.Errorf("table %s.%s has no columns", s.Name, t.Name)
			}
			schemaMeta.Tables = append(schemaMeta.Tables, &TableMetadata{
				Name:    t.Name,
				Columns: cols,
				Comment: t.Comment,
			})
		}
		for _, v := range s.Views {
			if v == nil || v.SkipDump {
				continue
			}
			if len(v.Triggers) > 0 || len(v.Rules) > 0 {
				return nil, fmt.Errorf("view %s.%s triggers or rules cannot be restored by Trino SDL renderer", s.Name, v.Name)
			}
			relKey := strings.ToLower(v.Name)
			if _, exists := seenRelations[relKey]; exists {
				return nil, fmt.Errorf("schema %s has duplicate view %s", s.Name, quoteIdentifier(v.Name))
			}
			seenRelations[relKey] = struct{}{}
			viewCols, err := fromProtoColumns(s.Name, v.Name, v.Columns, false)
			if err != nil {
				return nil, err
			}
			schemaMeta.Views = append(schemaMeta.Views, &ViewMetadata{
				Name:         v.Name,
				Definition:   v.Definition,
				Comment:      v.Comment,
				Columns:      viewCols,
				Dependencies: fromProtoDependencies(v.DependencyColumns),
			})
		}
		for _, mv := range s.MaterializedViews {
			if mv == nil || mv.SkipDump {
				continue
			}
			if mv.Comment != "" {
				return nil, fmt.Errorf("materialized view %s.%s comment cannot be restored on pinned Trino 482", s.Name, mv.Name)
			}
			if len(mv.Indexes) > 0 || len(mv.Triggers) > 0 {
				return nil, fmt.Errorf("materialized view %s.%s indexes or triggers cannot be restored by Trino SDL renderer", s.Name, mv.Name)
			}
			relKey := strings.ToLower(mv.Name)
			if _, exists := seenRelations[relKey]; exists {
				return nil, fmt.Errorf("schema %s has duplicate materialized view %s", s.Name, quoteIdentifier(mv.Name))
			}
			seenRelations[relKey] = struct{}{}
			schemaMeta.MaterializedViews = append(schemaMeta.MaterializedViews, &ViewMetadata{
				Name:         mv.Name,
				Definition:   mv.Definition,
				Comment:      mv.Comment,
				Dependencies: fromProtoDependencies(mv.DependencyColumns),
			})
		}
		out.Schemas = append(out.Schemas, schemaMeta)
	}
	return out, nil
}

func qualifyTrinoFunctionDeclaration(dbName, schemaName, fnName, definition string) (*trinoparser.CreateFunctionStmt, string, bool, error) {
	cfs, hasSQLBody, _, ok := parseTrinoCreateFunction(definition)
	if !ok || cfs.Spec.Declaration.Name == nil || len(cfs.Spec.Declaration.Name.Parts) == 0 {
		def, err := qualifyTrinoFunctionHeaderLexically(dbName, schemaName, fnName, definition)
		return nil, def, false, err
	}
	if hasSQLBody && cfs.Spec.Body != nil {
		w := &trinoASTWalker{callerSchema: schemaName}
		w.walkStmt(cfs.Spec.Body, nil)
		if w.hasQuery {
			return nil, "", false, fmt.Errorf("function %s.%s contains a query or subquery, which Trino SQL UDFs cannot process", schemaName, fnName)
		}
	}
	declName := cfs.Spec.Declaration.Name
	parts := declName.Parts
	mismatch := func() error {
		return fmt.Errorf("function %s.%s definition name %q does not match metadata identity", schemaName, fnName, declName.String())
	}
	def := definition
	switch len(parts) {
	case 1:
		if !strings.EqualFold(parts[0].Value, fnName) {
			return nil, "", false, mismatch()
		}
		if declName.Loc.Start >= 0 && declName.Loc.End <= len(definition) && declName.Loc.Start < declName.Loc.End {
			if prefix := qualifyTrinoObjectPrefix(dbName, schemaName); prefix != "" {
				def = definition[:declName.Loc.Start] + prefix + definition[declName.Loc.Start:]
			}
		}
	case 2:
		if !strings.EqualFold(parts[0].Value, schemaName) || !strings.EqualFold(parts[1].Value, fnName) {
			return nil, "", false, mismatch()
		}
		if dbName != "" && declName.Loc.Start >= 0 && declName.Loc.End <= len(definition) && declName.Loc.Start < declName.Loc.End {
			def = definition[:declName.Loc.Start] + quoteIdentifier(dbName) + "." + definition[declName.Loc.Start:]
		}
	case 3:
		if (dbName != "" && !strings.EqualFold(parts[0].Value, dbName)) ||
			!strings.EqualFold(parts[1].Value, schemaName) ||
			!strings.EqualFold(parts[2].Value, fnName) {
			return nil, "", false, mismatch()
		}
	default:
		return nil, "", false, mismatch()
	}
	return cfs, def, true, nil
}

func scanTrinoDeclIdentParts(definition string) ([]struct {
	text  string
	start int
	end   int
}, int, bool) {
	type identSpan struct {
		text  string
		start int
		end   int
	}
	readIdent := func(pos int) (identSpan, int, bool) {
		pos = skipSQLWhitespaceAndComments(definition, pos)
		if pos >= len(definition) {
			return identSpan{}, pos, false
		}
		if definition[pos] == '"' {
			start := pos
			pos++
			var sb strings.Builder
			for pos < len(definition) {
				if definition[pos] == '"' {
					if pos+1 < len(definition) && definition[pos+1] == '"' {
						sb.WriteByte('"')
						pos += 2
						continue
					}
					pos++
					return identSpan{text: sb.String(), start: start, end: pos}, pos, true
				}
				sb.WriteByte(definition[pos])
				pos++
			}
			return identSpan{}, pos, false
		}
		if !isIdentStartChar(definition[pos]) {
			return identSpan{}, pos, false
		}
		start := pos
		pos++
		for pos < len(definition) && isIdentContChar(definition[pos]) {
			pos++
		}
		return identSpan{text: definition[start:pos], start: start, end: pos}, pos, true
	}

	pos := 0
	kw1, nextPos, ok := readIdent(pos)
	if !ok || !strings.EqualFold(kw1.text, "CREATE") {
		return nil, 0, false
	}
	pos = nextPos
	kw2, nextPos, ok := readIdent(pos)
	if !ok {
		return nil, 0, false
	}
	pos = nextPos
	if strings.EqualFold(kw2.text, "OR") {
		kwReplace, p2, ok := readIdent(pos)
		if !ok || !strings.EqualFold(kwReplace.text, "REPLACE") {
			return nil, 0, false
		}
		kwFunc, p3, ok := readIdent(p2)
		if !ok || !strings.EqualFold(kwFunc.text, "FUNCTION") {
			return nil, 0, false
		}
		pos = p3
	} else if !strings.EqualFold(kw2.text, "FUNCTION") {
		return nil, 0, false
	}

	var parts []struct {
		text  string
		start int
		end   int
	}
	for len(parts) < 3 {
		part, pNext, ok := readIdent(pos)
		if !ok {
			return nil, 0, false
		}
		parts = append(parts, struct {
			text  string
			start int
			end   int
		}(part))
		pos = skipSQLWhitespaceAndComments(definition, pNext)
		if pos < len(definition) && definition[pos] == '.' {
			pos++
			continue
		}
		break
	}
	pos = skipSQLWhitespaceAndComments(definition, pos)
	if pos >= len(definition) || definition[pos] != '(' {
		return nil, 0, false
	}
	return parts, pos, true
}

func qualifyTrinoFunctionHeaderLexically(dbName, schemaName, fnName, definition string) (string, error) {
	parts, _, ok := scanTrinoDeclIdentParts(definition)
	if !ok || len(parts) == 0 {
		return "", fmt.Errorf("function %s.%s definition must be a CREATE FUNCTION statement", schemaName, fnName)
	}
	mismatch := func() error {
		return fmt.Errorf("function %s.%s definition name does not match metadata identity", schemaName, fnName)
	}
	switch len(parts) {
	case 1:
		if !strings.EqualFold(parts[0].text, fnName) {
			return "", mismatch()
		}
		if prefix := qualifyTrinoObjectPrefix(dbName, schemaName); prefix != "" {
			return definition[:parts[0].start] + prefix + definition[parts[0].start:], nil
		}
		return definition, nil
	case 2:
		if !strings.EqualFold(parts[0].text, schemaName) || !strings.EqualFold(parts[1].text, fnName) {
			return "", mismatch()
		}
		if dbName != "" {
			return definition[:parts[0].start] + quoteIdentifier(dbName) + "." + definition[parts[0].start:], nil
		}
		return definition, nil
	case 3:
		if (dbName != "" && !strings.EqualFold(parts[0].text, dbName)) ||
			!strings.EqualFold(parts[1].text, schemaName) ||
			!strings.EqualFold(parts[2].text, fnName) {
			return "", mismatch()
		}
		return definition, nil
	default:
		return "", mismatch()
	}
}

func convertSkippedTrinoFunction(fn *metadatapb.FunctionMetadata) *FunctionMetadata {
	def := fn.Definition
	if strings.TrimSpace(def) == "" && strings.TrimSpace(fn.Signature) != "" {
		sig := strings.TrimSpace(fn.Signature)
		if !strings.HasPrefix(sig, "(") {
			sig = "(" + sig + ")"
		}
		def = "CREATE FUNCTION " + quoteIdentifier(fn.Name) + sig + " RETURNS varchar RETURN ''"
	}
	return &FunctionMetadata{
		Name:       fn.Name,
		Definition: def,
		Comment:    fn.Comment,
		SkipDump:   true,
	}
}

func prepareTrinoFunctionDefinition(dbName, schemaName string, fn *metadatapb.FunctionMetadata) (string, string, error) {
	if strings.TrimSpace(fn.Definition) == "" {
		return "", "", fmt.Errorf("function %s.%s has no definition to restore", schemaName, fn.Name)
	}
	if !isTrinoCreateFunctionHeader(fn.Definition) {
		return "", "", fmt.Errorf("function %s.%s definition must be a CREATE FUNCTION statement", schemaName, fn.Name)
	}
	fallbackSigKey := strings.ToLower(fn.Name) + "(" + strings.ToLower(strings.TrimSpace(fn.Signature)) + ")"
	cfs, def, ok, err := qualifyTrinoFunctionDeclaration(dbName, schemaName, fn.Name, fn.Definition)
	if err != nil {
		return "", "", err
	}
	if !ok {
		return fn.Definition, fallbackSigKey, nil
	}
	if err := validateTrinoFunctionSignature(schemaName, fn.Name, fn.Signature, cfs.Spec.Declaration.Parameters); err != nil {
		return "", "", err
	}
	paramTypes := make([]string, len(cfs.Spec.Declaration.Parameters))
	for i, p := range cfs.Spec.Declaration.Parameters {
		paramTypes[i] = normalizeTrinoDataType(p.Type)
	}
	sigKey := strings.ToLower(fn.Name) + "(" + strings.Join(paramTypes, ",") + ")"
	return def, sigKey, nil
}

func isTrinoCreateFunctionHeader(definition string) bool {
	var words []string
	i := 0
	for len(words) < 4 {
		i = skipSQLWhitespaceAndComments(definition, i)
		if i >= len(definition) || !isIdentStartChar(definition[i]) {
			break
		}
		j := i + 1
		for j < len(definition) && isIdentContChar(definition[j]) {
			j++
		}
		words = append(words, strings.ToUpper(definition[i:j]))
		i = j
	}
	if len(words) >= 2 && words[0] == "CREATE" && words[1] == "FUNCTION" {
		return true
	}
	return len(words) >= 4 && words[0] == "CREATE" && words[1] == "OR" && words[2] == "REPLACE" && words[3] == "FUNCTION"
}

func qualifyTrinoObjectPrefix(dbName, schemaName string) string {
	switch {
	case dbName != "" && schemaName != "":
		return quoteIdentifier(dbName) + "." + quoteIdentifier(schemaName) + "."
	case schemaName != "":
		return quoteIdentifier(schemaName) + "."
	default:
		return ""
	}
}

func qualifyTrinoViewBody(dbName, schemaName, definition string, deps []ViewDependency, cat *trinoFuncCatalog) string {
	if strings.TrimSpace(definition) == "" {
		return definition
	}
	file, err := trinoparser.Parse(definition)
	if err != nil || file == nil || len(file.Stmts) != 1 {
		return definition
	}
	var query *trinoparser.Query
	var inlineFns []trinoparser.FunctionSpecification
	switch s := file.Stmts[0].(type) {
	case *trinoparser.QueryStmt:
		query = s.Query
	case *trinoparser.WithFunctionStmt:
		query = s.Query
		inlineFns = s.Functions
	}
	if query == nil {
		return definition
	}
	calls, unqualTables := collectTrinoQueryASTRefs(query, inlineFns, schemaName, deps, cat)
	return applyTrinoCallPrefixEdits(definition, calls, unqualTables, dbName, schemaName, deps, cat)
}

func resolveUnqualifiedRelationSchema(relName, defaultSchema string, deps []ViewDependency) string {
	if len(deps) == 0 {
		return defaultSchema
	}
	var matched string
	for _, dep := range deps {
		if dep.Schema == "" || !strings.EqualFold(dep.Table, relName) {
			continue
		}
		if matched == "" {
			matched = dep.Schema
		} else if !strings.EqualFold(matched, dep.Schema) {
			// Ambiguous across schemas: keep the caller's own schema if it matches.
			if strings.EqualFold(dep.Schema, defaultSchema) || strings.EqualFold(matched, defaultSchema) {
				return defaultSchema
			}
			return ""
		}
	}
	if matched != "" {
		return matched
	}
	return defaultSchema
}

func applyTrinoCallPrefixEdits(sql string, calls []functionCallRef, unqualTables []*ast.QualifiedName, dbName, schemaName string, deps []ViewDependency, cat *trinoFuncCatalog) string {
	type prefixEdit struct {
		offset int
		prefix string
	}
	var edits []prefixEdit
	seen := make(map[int]struct{})
	addEdit := func(qnLoc ast.Loc, prefix string) {
		if prefix == "" || qnLoc.Start < 0 || qnLoc.End <= qnLoc.Start || qnLoc.End > len(sql) {
			return
		}
		if _, dup := seen[qnLoc.Start]; !dup {
			seen[qnLoc.Start] = struct{}{}
			edits = append(edits, prefixEdit{offset: qnLoc.Start, prefix: prefix})
		}
	}
	relPrefix := qualifyTrinoObjectPrefix(dbName, schemaName)
	for _, qn := range unqualTables {
		targetSchema := schemaName
		if len(qn.Parts) == 1 {
			targetSchema = resolveUnqualifiedRelationSchema(qn.Parts[0].Value, schemaName, deps)
		}
		addEdit(qn.Loc, qualifyTrinoObjectPrefix(dbName, targetSchema))
	}
	for _, call := range calls {
		if len(cat.resolveCall(schemaName, call)) == 0 || call.catalog != "" {
			continue
		}
		var prefix string
		if call.schema == "" {
			prefix = relPrefix
		} else if dbName != "" {
			prefix = quoteIdentifier(dbName) + "."
		}
		addEdit(ast.Loc{Start: call.nameStart, End: call.nameEnd}, prefix)
	}
	if len(edits) == 0 {
		return sql
	}
	slices.SortFunc(edits, func(a, b prefixEdit) int { return b.offset - a.offset })
	out := sql
	for _, ed := range edits {
		out = out[:ed.offset] + ed.prefix + out[ed.offset:]
	}
	return out
}

func validateTrinoFunctionSignature(schemaName, fnName, signature string, declParams []trinoparser.ParameterDeclaration) error {
	sig := strings.TrimSpace(signature)
	if sig == "" {
		return nil
	}
	if open := strings.IndexByte(sig, '('); open >= 0 {
		close := strings.LastIndexByte(sig, ')')
		if close <= open {
			return fmt.Errorf("function %s.%s metadata signature %q is malformed", schemaName, fnName, signature)
		}
		sig = strings.TrimSpace(sig[open+1 : close])
	}
	sigFile, err := trinoparser.Parse("CREATE FUNCTION _sig(" + sig + ") RETURNS int RETURN 1")
	if err != nil || sigFile == nil || len(sigFile.Stmts) != 1 {
		return nil
	}
	sigStmt, ok := sigFile.Stmts[0].(*trinoparser.CreateFunctionStmt)
	if !ok {
		return nil
	}
	sigParams := sigStmt.Spec.Declaration.Parameters
	if len(sigParams) != len(declParams) {
		return fmt.Errorf("function %s.%s definition parameter count (%d) does not match signature %q", schemaName, fnName, len(declParams), signature)
	}
	for i := range sigParams {
		st := normalizeTrinoDataType(sigParams[i].Type)
		dt := normalizeTrinoDataType(declParams[i].Type)
		if st != "" && dt != "" && st != dt {
			return fmt.Errorf("function %s.%s parameter %d type %q does not match signature %q", schemaName, fnName, i+1, dt, signature)
		}
	}
	return nil
}

func fromProtoColumns(schemaName, tableName string, cols []*metadatapb.ColumnMetadata, requireType bool) ([]*ColumnMetadata, error) {
	if len(cols) == 0 {
		return nil, nil
	}
	out := make([]*ColumnMetadata, 0, len(cols))
	seenCols := make(map[string]struct{}, len(cols))
	for _, c := range cols {
		if c == nil {
			continue
		}
		colKey := strings.ToLower(c.Name)
		if _, exists := seenCols[colKey]; exists {
			if requireType {
				return nil, fmt.Errorf("table %s.%s has duplicate column %s", schemaName, tableName, quoteIdentifier(c.Name))
			}
			return nil, fmt.Errorf("view %s.%s has duplicate column %s", schemaName, tableName, quoteIdentifier(c.Name))
		}
		seenCols[colKey] = struct{}{}
		if requireType && strings.TrimSpace(c.Type) == "" {
			return nil, fmt.Errorf("column %s.%s.%s has no data type", schemaName, tableName, c.Name)
		}
		if !requireType && c.Default != "" {
			return nil, fmt.Errorf("column %s.%s.%s default cannot be restored on a view by Trino SDL renderer", schemaName, tableName, c.Name)
		}
		if c.Generation != nil {
			return nil, fmt.Errorf("column %s.%s.%s generated expression cannot be restored by Trino SDL renderer", schemaName, tableName, c.Name)
		}
		if c.IsIdentity || c.IdentityGeneration != 0 || c.IdentitySeed != 0 || c.IdentityIncrement != 0 || c.DefaultOnNull || c.OnUpdate != "" || c.CharacterSet != "" || c.Collation != "" || c.CollationSchema != "" || len(c.Options) > 0 || c.Srid != nil || c.IsInvisible || c.DefaultConstraintName != "" {
			return nil, fmt.Errorf("column %s.%s.%s attributes cannot be restored by Trino SDL renderer", schemaName, tableName, c.Name)
		}
		out = append(out, &ColumnMetadata{
			Name:     c.Name,
			Type:     c.Type,
			Nullable: c.Nullable,
			Default:  c.Default,
			Comment:  c.Comment,
		})
	}
	return out, nil
}

func fromProtoDependencies(deps []*metadatapb.DependencyColumn) []ViewDependency {
	if len(deps) == 0 {
		return nil
	}
	out := make([]ViewDependency, 0, len(deps))
	for _, dep := range deps {
		if dep == nil || dep.Table == "" {
			continue
		}
		out = append(out, ViewDependency{
			Schema: dep.Schema,
			Table:  dep.Table,
		})
	}
	return out
}

// GetDatabaseDefinition renders the schema metadata as the DDL that recreates
// it: one CREATE SCHEMA per schema, then one CREATE TABLE per table, then one
// CREATE VIEW per view. It is the omni-side implementation of the bytebase
// schema.GetDatabaseDefinition contract for the Trino engine.
//
// The output is in snapshot order within each kind, each statement separated by
// a blank line and terminated with the same trailing "\n\n" the bytebase source
// emits. Names are written always double-quoted, with an embedded double quote
// escaped as two. Column types are emitted as-is. Non-nullable columns get a
// trailing " NOT NULL".
//
// Schemas come first because a restore starts in a catalog that may not have
// them; tables come before views because a view may read a table in another
// schema. Views are ordered by their definition's references to other views, so
// a view that reads another is written second. A nil metadata yields the empty
// string. A view with no definition is an error: writing it would silently lose
// the view.
func GetDatabaseDefinition(metadata *DatabaseSchemaMetadata) (string, error) {
	if metadata == nil {
		return "", nil
	}
	seenSchemas := make(map[string]struct{}, len(metadata.Schemas))
	for _, schema := range metadata.Schemas {
		if schema == nil || schema.SkipDump {
			continue
		}
		schemaKey := strings.ToLower(schema.Name)
		if _, exists := seenSchemas[schemaKey]; exists {
			return "", fmt.Errorf("database %s has duplicate schema %s", metadata.Name, quoteIdentifier(schema.Name))
		}
		seenSchemas[schemaKey] = struct{}{}
		if schema.Owner != "" {
			return "", fmt.Errorf("schema %s owner %q cannot be restored unambiguously because Trino SET AUTHORIZATION requires user vs role distinction", schema.Name, schema.Owner)
		}
		seenRelations := make(map[string]struct{}, len(schema.Tables)+len(schema.Views)+len(schema.MaterializedViews))
		checkRelation := func(kind, name string) error {
			relKey := strings.ToLower(name)
			if _, exists := seenRelations[relKey]; exists {
				return fmt.Errorf("schema %s has duplicate %s %s", schema.Name, kind, quoteIdentifier(name))
			}
			seenRelations[relKey] = struct{}{}
			return nil
		}
		for _, table := range schema.Tables {
			if table != nil {
				if table.Owner != "" {
					return "", fmt.Errorf("table %s.%s owner %q cannot be restored unambiguously because Trino SET AUTHORIZATION requires user vs role distinction", schema.Name, table.Name, table.Owner)
				}
				if err := checkRelation("table", table.Name); err != nil {
					return "", err
				}
			}
		}
		for _, view := range schema.Views {
			if view != nil {
				if err := checkRelation("view", view.Name); err != nil {
					return "", err
				}
			}
		}
		for _, mv := range schema.MaterializedViews {
			if mv != nil {
				if err := checkRelation("materialized view", mv.Name); err != nil {
					return "", err
				}
			}
		}
	}

	var buf strings.Builder
	for _, schema := range metadata.Schemas {
		if schema == nil || schema.SkipDump {
			continue
		}
		writeCreateSchema(&buf, schema.Name)
		buf.WriteString("\n\n")
	}
	fns, err := orderedFunctions(metadata)
	if err != nil {
		return "", err
	}
	views, err := orderedViews(metadata)
	if err != nil {
		return "", err
	}
	for _, fn := range fns {
		if err := writeCreateFunction(&buf, metadata.Name, fn.schema, fn.fn); err != nil {
			return "", err
		}
		buf.WriteString("\n\n")
	}
	funcCat := buildTrinoFuncCatalog(metadata)
	for _, schema := range metadata.Schemas {
		if schema == nil || schema.SkipDump {
			continue
		}
		for _, table := range schema.Tables {
			if table == nil {
				continue
			}
			if err := writeCreateTableWithCatalog(&buf, metadata.Name, schema.Name, table, funcCat); err != nil {
				return "", err
			}
			buf.WriteString("\n\n")
		}
	}
	for _, view := range views {
		viewCopy := *view.view
		viewCopy.Definition = qualifyTrinoViewBody(metadata.Name, view.schema, view.view.Definition, view.view.Dependencies, funcCat)
		if err := writeCreateView(&buf, view.schema, &viewCopy, view.materialized); err != nil {
			return "", err
		}
		buf.WriteString("\n\n")
	}

	return buf.String(), nil
}

// quoteIdentifier writes name as a double-quoted Trino identifier. A double
// quote inside name is doubled, as Trino's own lexer reads it back.
func quoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

// writeCreateSchema appends the CREATE SCHEMA statement for schemaName. It does
// not append the inter-statement blank line; the caller does.
func writeCreateSchema(buf *strings.Builder, schemaName string) {
	fmt.Fprintf(buf, "CREATE SCHEMA IF NOT EXISTS %s;", quoteIdentifier(schemaName))
}

func qualifyTrinoDefault(dbName, schemaName, raw string, cat *trinoFuncCatalog) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}
	expr, errs := trinoparser.ParseExpression(trimmed)
	if len(errs) == 0 && expr != nil && isTrinoDefaultAST(expr) {
		if cat == nil {
			return trimmed
		}
		w := &trinoASTWalker{callerSchema: schemaName, cat: cat}
		w.walkExpr(expr, nil, nil)
		return applyTrinoCallPrefixEdits(trimmed, w.calls, w.unqualTables, dbName, schemaName, nil, cat)
	}
	return fmt.Sprintf("'%s'", escapeStringLiteral(raw))
}

func isTrinoDefaultAST(expr trinoparser.Expr) bool {
	switch e := expr.(type) {
	case *trinoparser.Literal, *trinoparser.IntervalLiteral:
		return true
	case *trinoparser.UnaryExpr:
		return isTrinoDefaultAST(e.Operand)
	case *trinoparser.BinaryExpr:
		return isTrinoDefaultAST(e.Left) && isTrinoDefaultAST(e.Right)
	case *trinoparser.ParenExpr:
		return isTrinoDefaultAST(e.Expr)
	case *trinoparser.RowConstructor:
		for _, el := range e.Elements {
			if !isTrinoDefaultAST(el) {
				return false
			}
		}
		return true
	case *trinoparser.ColumnRef, *trinoparser.Dereference, *trinoparser.Parameter:
		return false
	default:
		return true
	}
}

func writeCreateFunction(buf *strings.Builder, dbName, schemaName string, fn *FunctionMetadata) error {
	definition := strings.TrimSpace(fn.Definition)
	if definition == "" {
		return fmt.Errorf("function %s.%s has no definition to restore", schemaName, fn.Name)
	}
	if fn.Comment != "" {
		var err error
		definition, err = injectFunctionComment(schemaName, fn.Name, definition, fn.Comment)
		if err != nil {
			return err
		}
	}
	_, definition, _, err := qualifyTrinoFunctionDeclaration(dbName, schemaName, fn.Name, definition)
	if err != nil {
		return err
	}
	buf.WriteString(appendTerminator(definition))
	return nil
}

func injectFunctionComment(schemaName, fnName, definition, comment string) (string, error) {
	cfs, _, bodyStart, ok := parseTrinoCreateFunction(definition)
	if !ok {
		return "", fmt.Errorf("function %s.%s comment cannot be placed in unparseable definition", schemaName, fnName)
	}
	clause := fmt.Sprintf("COMMENT '%s'", escapeStringLiteral(comment))
	for _, ch := range cfs.Spec.Characteristics {
		if ch.Kind == trinoparser.CharacteristicComment && ch.Loc.Start >= 0 && ch.Loc.End <= len(definition) && ch.Loc.Start < ch.Loc.End {
			return definition[:ch.Loc.Start] + clause + definition[ch.Loc.End:], nil
		}
	}
	if bodyStart < 0 || bodyStart > len(definition) {
		return "", fmt.Errorf("function %s.%s has invalid body span", schemaName, fnName)
	}
	prefix := strings.TrimRight(definition[:bodyStart], " \t\r\n")
	sep := " "
	if endsInLineComment(prefix) {
		sep = "\n"
	}
	return prefix + sep + clause + " " + definition[bodyStart:], nil
}

// parseTrinoCreateFunction parses a CREATE FUNCTION statement, including non-SQL
// routines (such as LANGUAGE PYTHON ... AS $$...$$) whose header and routine
// characteristics share Trino's CREATE FUNCTION grammar up to the AS body.
func parseTrinoCreateFunction(definition string) (*trinoparser.CreateFunctionStmt, bool, int, bool) {
	if file, err := trinoparser.Parse(definition); err == nil && file != nil && len(file.Stmts) == 1 {
		if cfs, ok := file.Stmts[0].(*trinoparser.CreateFunctionStmt); ok && cfs.Spec.Body != nil {
			return cfs, true, cfs.Spec.Body.Span().Start, true
		}
	}
	asOffset, ok := splitTrinoDollarFunctionBody(definition)
	if !ok {
		return nil, false, -1, false
	}
	file, err := trinoparser.Parse(definition[:asOffset] + "RETURN 1")
	if err != nil || file == nil || len(file.Stmts) != 1 {
		return nil, false, -1, false
	}
	cfs, ok := file.Stmts[0].(*trinoparser.CreateFunctionStmt)
	if !ok || cfs.Spec.Body == nil {
		return nil, false, -1, false
	}
	return cfs, false, asOffset, true
}

func splitTrinoDollarFunctionBody(definition string) (int, bool) {
	var (
		inSingleQuote  bool
		inDoubleQuote  bool
		inLineComment  bool
		inBlockComment bool
		parenDepth     int
	)
	for i := 0; i < len(definition); {
		ch := definition[i]
		if inLineComment {
			if ch == '\n' {
				inLineComment = false
			}
			i++
			continue
		}
		if inBlockComment {
			if ch == '*' && i+1 < len(definition) && definition[i+1] == '/' {
				inBlockComment = false
				i += 2
				continue
			}
			i++
			continue
		}
		if inSingleQuote {
			if ch == '\'' {
				if i+1 < len(definition) && definition[i+1] == '\'' {
					i += 2
					continue
				}
				inSingleQuote = false
			}
			i++
			continue
		}
		if inDoubleQuote {
			if ch == '"' {
				if i+1 < len(definition) && definition[i+1] == '"' {
					i += 2
					continue
				}
				inDoubleQuote = false
			}
			i++
			continue
		}
		switch {
		case ch == '-' && i+1 < len(definition) && definition[i+1] == '-':
			inLineComment = true
			i += 2
		case ch == '/' && i+1 < len(definition) && definition[i+1] == '*':
			inBlockComment = true
			i += 2
		case ch == '\'':
			inSingleQuote = true
			i++
		case ch == '"':
			inDoubleQuote = true
			i++
		case ch == '(':
			parenDepth++
			i++
		case ch == ')':
			if parenDepth > 0 {
				parenDepth--
			}
			i++
		case parenDepth == 0 && isIdentStartChar(ch) && (i == 0 || (!isIdentContChar(definition[i-1]) && definition[i-1] != '$')):
			j := i + 1
			for j < len(definition) && isIdentContChar(definition[j]) {
				j++
			}
			if strings.EqualFold(definition[i:j], "AS") {
				k := skipSQLWhitespaceAndComments(definition, j)
				if bodyEnd, ok := matchTrinoDollarLiteral(definition, k); ok && onlyTrailingStatementEnd(definition[bodyEnd:]) {
					return i, true
				}
			}
			i = j
		default:
			i++
		}
	}
	return 0, false
}

func skipSQLWhitespaceAndComments(s string, i int) int {
	for i < len(s) {
		switch {
		case s[i] == ' ' || s[i] == '\t' || s[i] == '\r' || s[i] == '\n':
			i++
		case s[i] == '-' && i+1 < len(s) && s[i+1] == '-':
			i += 2
			for i < len(s) && s[i] != '\n' {
				i++
			}
		case s[i] == '/' && i+1 < len(s) && s[i+1] == '*':
			i += 2
			for i < len(s) {
				if s[i] == '*' && i+1 < len(s) && s[i+1] == '/' {
					i += 2
					break
				}
				i++
			}
		default:
			return i
		}
	}
	return i
}

func matchTrinoDollarLiteral(s string, start int) (int, bool) {
	if start+1 >= len(s) || s[start] != '$' {
		return 0, false
	}
	j := start + 1
	if j < len(s) && s[j] != '$' {
		if !isIdentStartChar(s[j]) {
			return 0, false
		}
		j++
		for j < len(s) && isIdentContChar(s[j]) {
			j++
		}
	}
	if j >= len(s) || s[j] != '$' {
		return 0, false
	}
	tag := s[start : j+1]
	closeRel := strings.Index(s[j+1:], tag)
	if closeRel < 0 {
		return 0, false
	}
	return j + 1 + closeRel + len(tag), true
}

func onlyTrailingStatementEnd(rest string) bool {
	i := skipSQLWhitespaceAndComments(rest, 0)
	if i < len(rest) && rest[i] == ';' {
		i = skipSQLWhitespaceAndComments(rest, i+1)
	}
	return i == len(rest)
}

func endsInLineComment(s string) bool {
	var inSingle, inDouble, inLine, inBlock bool
	for i := 0; i < len(s); {
		ch := s[i]
		switch {
		case inLine:
			if ch == '\n' {
				inLine = false
			}
			i++
		case inBlock:
			if ch == '*' && i+1 < len(s) && s[i+1] == '/' {
				inBlock = false
				i += 2
				continue
			}
			i++
		case inSingle:
			if ch == '\'' {
				if i+1 < len(s) && s[i+1] == '\'' {
					i += 2
					continue
				}
				inSingle = false
			}
			i++
		case inDouble:
			if ch == '"' {
				if i+1 < len(s) && s[i+1] == '"' {
					i += 2
					continue
				}
				inDouble = false
			}
			i++
		case ch == '-' && i+1 < len(s) && s[i+1] == '-':
			inLine = true
			i += 2
		case ch == '/' && i+1 < len(s) && s[i+1] == '*':
			inBlock = true
			i += 2
		case ch == '\'':
			inSingle = true
			i++
		case ch == '"':
			inDouble = true
			i++
		default:
			i++
		}
	}
	return inLine
}

type orderedFunction struct {
	schema  string
	fn      *FunctionMetadata
	skipped bool
}

type trinoFuncParam struct {
	name     string
	typeName string
}

type parsedTrinoFn struct {
	params     []trinoFuncParam
	returnType string
	bodyStmt   trinoparser.RoutineStatement
	parsed     bool
}

type trinoFuncCatalog struct {
	dbName    string
	fns       []orderedFunction
	parsed    []parsedTrinoFn
	byKey     map[string][]int
	tableCols map[string]map[string]string
}

func parseTrinoColumnTypeString(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}
	if file, err := trinoparser.Parse("CREATE FUNCTION _c(x " + trimmed + ") RETURNS int RETURN 1"); err == nil && file != nil && len(file.Stmts) == 1 {
		if cfs, ok := file.Stmts[0].(*trinoparser.CreateFunctionStmt); ok && len(cfs.Spec.Declaration.Parameters) == 1 {
			if dt := normalizeTrinoDataType(cfs.Spec.Declaration.Parameters[0].Type); dt != "" {
				return dt
			}
		}
	}
	return normalizeTrinoTypeName(trimmed)
}

func buildTrinoFuncCatalog(metadata *DatabaseSchemaMetadata) *trinoFuncCatalog {
	var fns []orderedFunction
	tableCols := make(map[string]map[string]string)
	if metadata != nil {
		for _, schema := range metadata.Schemas {
			if schema == nil {
				continue
			}
			for _, fn := range schema.Functions {
				if fn == nil {
					continue
				}
				fns = append(fns, orderedFunction{
					schema:  schema.Name,
					fn:      fn,
					skipped: schema.SkipDump || fn.SkipDump,
				})
			}
			if schema.SkipDump {
				continue
			}
			for _, tbl := range schema.Tables {
				if tbl == nil || len(tbl.Columns) == 0 {
					continue
				}
				cols := make(map[string]string, len(tbl.Columns))
				for _, col := range tbl.Columns {
					if col == nil || col.Name == "" {
						continue
					}
					if dt := parseTrinoColumnTypeString(col.Type); dt != "" {
						cols[strings.ToLower(col.Name)] = dt
					}
				}
				if len(cols) > 0 {
					tableCols[objectKey(strings.ToLower(schema.Name), strings.ToLower(tbl.Name))] = cols
				}
			}
		}
	}
	parsed := make([]parsedTrinoFn, len(fns))
	byKey := make(map[string][]int, len(fns))
	for i, f := range fns {
		key := objectKey(strings.ToLower(f.schema), strings.ToLower(f.fn.Name))
		byKey[key] = append(byKey[key], i)

		if cfs, hasSQLBody, _, ok := parseTrinoCreateFunction(f.fn.Definition); ok {
			params := make([]trinoFuncParam, len(cfs.Spec.Declaration.Parameters))
			for j, p := range cfs.Spec.Declaration.Parameters {
				name := ""
				if p.Name != nil {
					name = p.Name.Value
				}
				params[j] = trinoFuncParam{
					name:     name,
					typeName: normalizeTrinoDataType(p.Type),
				}
			}
			var bodyStmt trinoparser.RoutineStatement
			if hasSQLBody {
				bodyStmt = cfs.Spec.Body
			}
			parsed[i] = parsedTrinoFn{
				params:     params,
				returnType: normalizeTrinoDataType(cfs.Spec.ReturnType),
				bodyStmt:   bodyStmt,
				parsed:     true,
			}
		}
	}
	dbName := ""
	if metadata != nil {
		dbName = metadata.Name
	}
	return &trinoFuncCatalog{
		dbName:    dbName,
		fns:       fns,
		parsed:    parsed,
		byKey:     byKey,
		tableCols: tableCols,
	}
}

func trinoFunctionSignatureKey(fnName, definition string, parsedFn parsedTrinoFn) string {
	if parsedFn.parsed {
		paramTypes := make([]string, len(parsedFn.params))
		for i, p := range parsedFn.params {
			paramTypes[i] = p.typeName
		}
		return strings.ToLower(fnName) + "(" + strings.Join(paramTypes, ",") + ")"
	}
	if _, openParen, ok := scanTrinoDeclIdentParts(definition); ok {
		depth := 0
		inSingle := false
		inDouble := false
		closeParen := -1
		for i := openParen; i < len(definition); i++ {
			ch := definition[i]
			if inSingle {
				if ch == '\'' {
					if i+1 < len(definition) && definition[i+1] == '\'' {
						i++
					} else {
						inSingle = false
					}
				}
				continue
			}
			if inDouble {
				if ch == '"' {
					if i+1 < len(definition) && definition[i+1] == '"' {
						i++
					} else {
						inDouble = false
					}
				}
				continue
			}
			switch ch {
			case '\'':
				inSingle = true
			case '"':
				inDouble = true
			case '(':
				depth++
			case ')':
				depth--
				if depth == 0 {
					closeParen = i
				}
			}
			if closeParen >= 0 {
				break
			}
		}
		if closeParen > openParen {
			inner := strings.TrimSpace(definition[openParen+1 : closeParen])
			if sigFile, err := trinoparser.Parse("CREATE FUNCTION _sig(" + inner + ") RETURNS int RETURN 1"); err == nil && sigFile != nil && len(sigFile.Stmts) == 1 {
				if cfs, ok := sigFile.Stmts[0].(*trinoparser.CreateFunctionStmt); ok {
					paramTypes := make([]string, len(cfs.Spec.Declaration.Parameters))
					for i, p := range cfs.Spec.Declaration.Parameters {
						paramTypes[i] = normalizeTrinoDataType(p.Type)
					}
					return strings.ToLower(fnName) + "(" + strings.Join(paramTypes, ",") + ")"
				}
			}
		}
	}
	return strings.ToLower(fnName)
}

func orderedFunctions(metadata *DatabaseSchemaMetadata) ([]orderedFunction, error) {
	cat := buildTrinoFuncCatalog(metadata)
	allFns := cat.fns
	if len(allFns) == 0 {
		return nil, nil
	}
	emittedIndices := make([]int, 0, len(allFns))
	emittedPos := make(map[int]int, len(allFns))
	for i, f := range allFns {
		if f.skipped {
			continue
		}
		emittedPos[i] = len(emittedIndices)
		emittedIndices = append(emittedIndices, i)
	}
	if len(emittedIndices) == 0 {
		return nil, nil
	}
	parsed := cat.parsed
	seenSigs := make(map[string]struct{}, len(emittedIndices))
	for _, origIdx := range emittedIndices {
		f := allFns[origIdx]
		sigKey := trinoFunctionSignatureKey(f.fn.Name, f.fn.Definition, parsed[origIdx])
		fullKey := strings.ToLower(f.schema) + "\x00" + sigKey
		if _, exists := seenSigs[fullKey]; exists {
			return nil, fmt.Errorf("schema %s has duplicate function signature %s", f.schema, sigKey)
		}
		seenSigs[fullKey] = struct{}{}
	}

	type prefixEdit struct {
		offset int
		prefix string
	}
	dependencies := make([]map[int]struct{}, len(emittedIndices))
	for ei, origIdx := range emittedIndices {
		f := allFns[origIdx]
		var calls []functionCallRef
		switch {
		case parsed[origIdx].bodyStmt != nil:
			calls = collectTrinoASTFunctionCalls(parsed[origIdx].bodyStmt, f.schema, parsed[origIdx].params, cat)
		case !parsed[origIdx].parsed:
			calls = scanFunctionBodyCallRefs(f.fn.Definition)
		}
		var edits []prefixEdit
		seenEdit := make(map[int]struct{})
		for _, call := range calls {
			matches := cat.resolveCall(f.schema, call)
			if len(matches) == 0 {
				continue
			}
			for _, depOrigIdx := range matches {
				if allFns[depOrigIdx].skipped {
					continue
				}
				depIdx := emittedPos[depOrigIdx]
				if depIdx == ei {
					return nil, fmt.Errorf("function %s.%s is self-recursive and cannot be restored because Trino disallows recursive SQL UDFs", f.schema, f.fn.Name)
				}
				if dependencies[ei] == nil {
					dependencies[ei] = map[int]struct{}{}
				}
				dependencies[ei][depIdx] = struct{}{}
			}
			if parsed[origIdx].parsed && call.catalog == "" && call.nameEnd > call.nameStart && call.nameEnd <= len(f.fn.Definition) {
				var prefix string
				if call.schema == "" {
					prefix = qualifyTrinoObjectPrefix(metadata.Name, f.schema)
				} else if metadata.Name != "" {
					prefix = quoteIdentifier(metadata.Name) + "."
				}
				if prefix != "" {
					if _, dup := seenEdit[call.nameStart]; !dup {
						seenEdit[call.nameStart] = struct{}{}
						edits = append(edits, prefixEdit{offset: call.nameStart, prefix: prefix})
					}
				}
			}
		}
		if len(edits) > 0 {
			slices.SortFunc(edits, func(a, b prefixEdit) int {
				return b.offset - a.offset
			})
			def := f.fn.Definition
			for _, ed := range edits {
				def = def[:ed.offset] + ed.prefix + def[ed.offset:]
			}
			allFns[origIdx].fn = &FunctionMetadata{
				Name:       f.fn.Name,
				Definition: def,
				Comment:    f.fn.Comment,
			}
		}
	}

	if len(emittedIndices) < 2 {
		return []orderedFunction{allFns[emittedIndices[0]]}, nil
	}
	indices, cyclic := topoOrder(len(emittedIndices), dependencies)
	if len(cyclic) > 0 {
		bad := allFns[emittedIndices[cyclic[0]]]
		return nil, fmt.Errorf("cyclic function dependency involving %s.%s cannot be restored because Trino requires UDFs to be declared before use", bad.schema, bad.fn.Name)
	}
	ordered := make([]orderedFunction, len(indices))
	for i, idx := range indices {
		ordered[i] = allFns[emittedIndices[idx]]
	}
	return ordered, nil
}

type functionCallRef struct {
	catalog   string
	schema    string
	name      string
	nameStart int
	nameEnd   int
	argNames  []string
	argTypes  []string
	argCount  int
}

func (c *trinoFuncCatalog) resolveCall(callerSchema string, call functionCallRef) []int {
	if c == nil || len(c.fns) == 0 {
		return nil
	}
	if call.catalog != "" && (c.dbName == "" || !strings.EqualFold(call.catalog, c.dbName)) {
		return nil
	}
	if call.catalog == "" && call.schema == "" && hasMatchingTrinoBuiltin(call) {
		return nil
	}
	targetSchema := call.schema
	if targetSchema == "" {
		targetSchema = callerSchema
	}
	candidates := c.byKey[objectKey(strings.ToLower(targetSchema), strings.ToLower(call.name))]
	if len(candidates) == 0 {
		return nil
	}
	var arityMatches []int
	for _, idx := range candidates {
		if !c.parsed[idx].parsed || len(c.parsed[idx].params) == call.argCount {
			arityMatches = append(arityMatches, idx)
		}
	}
	if len(arityMatches) == 0 {
		return nil
	}
	if len(call.argTypes) == 0 {
		return arityMatches
	}
	bestScore := -1
	var best []int
	for _, idx := range arityMatches {
		if !c.parsed[idx].parsed {
			if bestScore < 0 {
				bestScore = 0
				best = []int{idx}
			} else if bestScore == 0 {
				best = append(best, idx)
			}
			continue
		}
		argTypes, ok := reorderTrinoCallArgs(c.parsed[idx].params, call)
		if !ok {
			continue
		}
		score := 0
		compatible := true
		for i, argType := range argTypes {
			paramType := c.parsed[idx].params[i].typeName
			s, compat := trinoTypeCompatScore(argType, paramType)
			if !compat {
				compatible = false
				break
			}
			score += s
		}
		if !compatible {
			continue
		}
		if score > bestScore {
			bestScore = score
			best = []int{idx}
		} else if score == bestScore {
			best = append(best, idx)
		}
	}
	return best
}

func isTrinoSpecialFormCall(name string, call functionCallRef) bool {
	firstArg := ""
	if len(call.argTypes) > 0 {
		firstArg = call.argTypes[0]
	}
	switch name {
	case "coalesce", "grouping":
		return call.argCount >= 1
	case "format":
		return call.argCount >= 2 && trinoBuiltinTypeMatch(firstArg, "varchar")
	case "if":
		return (call.argCount == 2 || call.argCount == 3) && trinoBuiltinTypeMatch(firstArg, "boolean")
	case "nullif":
		return call.argCount == 2
	case "try":
		return call.argCount == 1
	default:
		return false
	}
}

func hasMatchingTrinoBuiltin(call functionCallRef) bool {
	for _, n := range call.argNames {
		if n != "" {
			return false
		}
	}
	name := strings.ToLower(call.name)
	if isTrinoSpecialFormCall(name, call) {
		return true
	}
	sigs, ok := trinoBuiltinFunctions[name]
	if !ok {
		return false
	}
	// SHOW FUNCTIONS lists variadic concat/greatest/least with a single element parameter.
	switch name {
	case "greatest", "least":
		if isCompatibleTrinoOrderableCall(call.argTypes, call.argCount) {
			return true
		}
	case "map_concat":
		if call.argCount >= 1 {
			allMaps := true
			for i := 0; i < call.argCount; i++ {
				argType := ""
				if i < len(call.argTypes) {
					argType = call.argTypes[i]
				}
				if !trinoBuiltinTypeMatch(argType, "map(*,*)") {
					allMaps = false
					break
				}
			}
			if allMaps {
				return true
			}
		}
	case "concat_ws":
		if call.argCount >= 2 {
			sepType := ""
			if len(call.argTypes) > 0 {
				sepType = call.argTypes[0]
			}
			if trinoBuiltinTypeMatch(sepType, "varchar") {
				allVarchar := true
				for i := 1; i < call.argCount; i++ {
					argType := ""
					if i < len(call.argTypes) {
						argType = call.argTypes[i]
					}
					if !trinoBuiltinTypeMatch(argType, "varchar") {
						allVarchar = false
						break
					}
				}
				if allVarchar {
					return true
				}
			}
		}
	case "concat":
		if call.argCount >= 2 {
			for _, elemType := range []string{"varchar", "varbinary", "array(*)"} {
				allMatch := true
				for i := 0; i < call.argCount; i++ {
					argType := ""
					if i < len(call.argTypes) {
						argType = call.argTypes[i]
					}
					if !trinoBuiltinTypeMatch(argType, elemType) {
						allMatch = false
						break
					}
				}
				if allMatch {
					return true
				}
			}
		}
	}
	for _, sig := range sigs {
		if len(sig) != call.argCount {
			continue
		}
		matched := true
		for i, paramType := range sig {
			argType := ""
			if i < len(call.argTypes) {
				argType = call.argTypes[i]
			}
			if !trinoBuiltinTypeMatch(argType, paramType) {
				matched = false
				break
			}
		}
		if matched {
			return true
		}
	}
	return false
}

func isOrderableTrinoType(t string) bool {
	if t == "" || t == "function" || t == "json" {
		return false
	}
	base, _, _ := splitParameterizedType(t)
	return base != "map" && base != "row"
}

func isCompatibleTrinoOrderableCall(argTypes []string, argCount int) bool {
	if argCount < 1 || len(argTypes) != argCount {
		return false
	}
	common := ""
	for _, t := range argTypes {
		if !isOrderableTrinoType(t) {
			return false
		}
		if common == "" {
			common = t
			continue
		}
		if _, ok := trinoTypeCompatScore(t, common); ok {
			continue
		}
		if _, ok := trinoTypeCompatScore(common, t); ok {
			common = t
			continue
		}
		return false
	}
	return common != ""
}

func trinoBuiltinTypeMatch(argType, paramType string) bool {
	if paramType == "*" {
		return true
	}
	if argType == "" {
		return false
	}
	if paramType == "function" {
		return argType == "function"
	}
	if innerParam, ok := strings.CutPrefix(paramType, "array("); ok && strings.HasSuffix(innerParam, ")") {
		innerParam = innerParam[:len(innerParam)-1]
		if argType == "array" {
			return true
		}
		innerArg, ok := strings.CutPrefix(argType, "array(")
		if !ok || !strings.HasSuffix(innerArg, ")") {
			return false
		}
		return trinoBuiltinTypeMatch(innerArg[:len(innerArg)-1], innerParam)
	}
	if innerParam, ok := strings.CutPrefix(paramType, "map("); ok && strings.HasSuffix(innerParam, ")") {
		if argType == "map" {
			return true
		}
		innerArg, ok := strings.CutPrefix(argType, "map(")
		if !ok || !strings.HasSuffix(innerArg, ")") {
			return false
		}
		pParts := splitTopLevelComma(innerParam[:len(innerParam)-1])
		aParts := splitTopLevelComma(innerArg[:len(innerArg)-1])
		if len(pParts) != 2 || len(aParts) != 2 {
			return false
		}
		return trinoBuiltinTypeMatch(aParts[0], pParts[0]) && trinoBuiltinTypeMatch(aParts[1], pParts[1])
	}
	if paramType == "row" && strings.HasPrefix(argType, "row") {
		return argType == "row" || strings.HasPrefix(argType, "row(")
	}
	_, compat := trinoTypeCompatScore(argType, paramType)
	return compat
}

func reorderTrinoCallArgs(params []trinoFuncParam, call functionCallRef) ([]string, bool) {
	if len(call.argNames) == 0 {
		return call.argTypes, true
	}
	hasNamed := false
	for _, n := range call.argNames {
		if n != "" {
			hasNamed = true
			break
		}
	}
	if !hasNamed {
		return call.argTypes, true
	}
	ordered := make([]string, len(params))
	assigned := make([]bool, len(params))
	seenNamed := false
	for i, argName := range call.argNames {
		if argName == "" {
			if seenNamed || i >= len(params) || assigned[i] {
				return nil, false
			}
			assigned[i] = true
			ordered[i] = call.argTypes[i]
			continue
		}
		seenNamed = true
		matched := false
		for p, param := range params {
			if param.name != "" && strings.EqualFold(param.name, argName) {
				if assigned[p] {
					return nil, false
				}
				assigned[p] = true
				ordered[p] = call.argTypes[i]
				matched = true
				break
			}
		}
		if !matched {
			return nil, false
		}
	}
	for _, ok := range assigned {
		if !ok {
			return nil, false
		}
	}
	return ordered, true
}

func trinoTypeCompatScore(argType, paramType string) (int, bool) {
	if argType == "" || paramType == "" {
		return 0, true
	}
	if paramType == argType {
		return 2, true
	}
	if isTrinoNumericType(argType) && isTrinoNumericType(paramType) && trinoNumericRank(paramType) >= trinoNumericRank(argType) {
		return 1, true
	}
	if argType == "char" && paramType == "varchar" {
		return 1, true
	}
	argBase, argInner, argHasInner := splitParameterizedType(argType)
	paramBase, paramInner, paramHasInner := splitParameterizedType(paramType)
	if argBase != paramBase {
		return 0, false
	}
	if !argHasInner || !paramHasInner {
		return 0, true
	}
	argParts := splitTopLevelComma(argInner)
	paramParts := splitTopLevelComma(paramInner)
	if len(argParts) != len(paramParts) {
		return 0, false
	}
	allExact := true
	anyUnknown := false
	for i := range argParts {
		aPart, pPart := argParts[i], paramParts[i]
		if argBase == "row" {
			_, aPart = splitRowField(aPart)
			_, pPart = splitRowField(pPart)
		}
		s, ok := trinoTypeCompatScore(aPart, pPart)
		if !ok {
			return 0, false
		}
		if s < 2 {
			allExact = false
		}
		if s == 0 {
			anyUnknown = true
		}
	}
	if allExact {
		return 2, true
	}
	if anyUnknown {
		return 0, true
	}
	return 1, true
}

func splitRowField(part string) (name, typ string) {
	colon := strings.IndexByte(part, ':')
	paren := strings.IndexByte(part, '(')
	if colon > 0 && (paren < 0 || colon < paren) {
		return part[:colon], part[colon+1:]
	}
	return "", part
}

func splitParameterizedType(t string) (base, inner string, hasInner bool) {
	idx := strings.IndexByte(t, '(')
	if idx < 0 || !strings.HasSuffix(t, ")") {
		return t, "", false
	}
	return t[:idx], t[idx+1 : len(t)-1], true
}

func splitTopLevelComma(s string) []string {
	var parts []string
	depth := 0
	start := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		case ',':
			if depth == 0 {
				parts = append(parts, s[start:i])
				start = i + 1
			}
		}
	}
	parts = append(parts, s[start:])
	return parts
}

func normalizeTrinoDataType(dt *trinoparser.DataType) string {
	if dt == nil {
		return ""
	}
	switch dt.Kind {
	case trinoparser.TypeDateTime:
		base := strings.ToLower(dt.Name)
		if dt.WithTimeZone {
			return base + " with time zone"
		}
		return base
	case trinoparser.TypeInterval:
		return normalizeTrinoIntervalType(dt.IntervalFrom, dt.IntervalTo)
	case trinoparser.TypeRow:
		if len(dt.Fields) == 0 {
			return "row"
		}
		fields := make([]string, len(dt.Fields))
		for i, f := range dt.Fields {
			ft := normalizeTrinoDataType(f.Type)
			if ft == "" {
				return "row"
			}
			if f.Name != nil && f.Name.Value != "" {
				ft = strings.ToLower(f.Name.Value) + ":" + ft
			}
			fields[i] = ft
		}
		return "row(" + strings.Join(fields, ",") + ")"
	case trinoparser.TypeArray:
		if elem := normalizeTrinoDataType(dt.ElementType); elem != "" {
			return "array(" + elem + ")"
		}
		return "array"
	case trinoparser.TypeMap:
		kt := normalizeTrinoDataType(dt.KeyType)
		vt := normalizeTrinoDataType(dt.ValueType)
		if kt != "" && vt != "" {
			return "map(" + kt + "," + vt + ")"
		}
		return "map"
	default:
		base := normalizeTrinoTypeName(dt.Name)
		switch base {
		case "array":
			if len(dt.Params) == 1 && dt.Params[0].Type != nil {
				if elem := normalizeTrinoDataType(dt.Params[0].Type); elem != "" {
					return "array(" + elem + ")"
				}
			}
		case "map":
			if len(dt.Params) == 2 && dt.Params[0].Type != nil && dt.Params[1].Type != nil {
				kt := normalizeTrinoDataType(dt.Params[0].Type)
				vt := normalizeTrinoDataType(dt.Params[1].Type)
				if kt != "" && vt != "" {
					return "map(" + kt + "," + vt + ")"
				}
			}
		}
		return base
	}
}

func normalizeTrinoIntervalType(from trinoparser.IntervalField, to *trinoparser.IntervalField) string {
	// Trino distinguishes interval families and unit qualifiers in overload resolution.
	if to != nil {
		return "interval " + strings.ToLower(from.String()) + " to " + strings.ToLower(to.String())
	}
	return "interval " + strings.ToLower(from.String())
}

func normalizeTrinoTypeName(name string) string {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "int1", "tinyint":
		return "tinyint"
	case "int2", "smallint":
		return "smallint"
	case "int", "integer", "int4":
		return "integer"
	case "int8", "bigint":
		return "bigint"
	case "dec", "decimal", "numeric":
		return "decimal"
	case "float", "real":
		return "real"
	case "double precision", "double":
		return "double"
	case "bool", "boolean":
		return "boolean"
	case "char", "character":
		return "char"
	case "varchar", "character varying", "text":
		return "varchar"
	default:
		return strings.ToLower(strings.TrimSpace(name))
	}
}

func trinoNumericRank(t string) int {
	switch t {
	case "tinyint":
		return 1
	case "smallint":
		return 2
	case "integer":
		return 3
	case "bigint":
		return 4
	case "decimal":
		return 5
	case "real":
		return 6
	case "double":
		return 7
	default:
		return 0
	}
}

func isTrinoNumericType(t string) bool {
	return trinoNumericRank(t) > 0
}

func zonedTrinoDateTimeType(base string) string {
	if base == "time" || base == "time with time zone" {
		return "time with time zone"
	}
	return "timestamp with time zone"
}

func commonTrinoNumericType(a, b string) string {
	ra, rb := trinoNumericRank(a), trinoNumericRank(b)
	if ra == 0 || rb == 0 {
		return ""
	}
	if ra >= rb {
		return a
	}
	return b
}

func combineTrinoTypes(a, b string) string {
	if a == "" {
		return b
	}
	if b == "" || a == b {
		return a
	}
	if common := commonTrinoNumericType(a, b); common != "" {
		return common
	}
	if (a == "char" && b == "varchar") || (a == "varchar" && b == "char") {
		return "varchar"
	}
	aBase, aInner, aHas := splitParameterizedType(a)
	bBase, bInner, bHas := splitParameterizedType(b)
	if aBase == bBase && aHas && bHas {
		aParts := splitTopLevelComma(aInner)
		bParts := splitTopLevelComma(bInner)
		if len(aParts) == len(bParts) {
			combined := make([]string, len(aParts))
			for i := range aParts {
				ct := combineTrinoTypes(aParts[i], bParts[i])
				if ct == "" {
					return aBase
				}
				combined[i] = ct
			}
			return aBase + "(" + strings.Join(combined, ",") + ")"
		}
	}
	return ""
}

type trinoInlineFunc struct {
	params     []trinoFuncParam
	returnType string
}

type trinoProjCol struct {
	name string
	typ  string
}

type trinoASTWalker struct {
	callerSchema string
	deps         []ViewDependency
	cat          *trinoFuncCatalog
	calls        []functionCallRef
	unqualTables []*ast.QualifiedName
	inlineFuncs  map[string][]trinoInlineFunc
	aliasCols    map[string]map[string]string
	cteCols      map[string]map[string]string
	fromSeen     map[string]struct{}
	hasQuery     bool
}

func collectTrinoASTFunctionCalls(root trinoparser.RoutineStatement, callerSchema string, params []trinoFuncParam, cat *trinoFuncCatalog) []functionCallRef {
	env := make(map[string]string, len(params))
	for _, p := range params {
		if p.name != "" {
			env[strings.ToLower(p.name)] = p.typeName
		}
	}
	w := &trinoASTWalker{callerSchema: callerSchema, cat: cat}
	w.walkStmt(root, env)
	return w.calls
}

func collectTrinoQueryASTRefs(q *trinoparser.Query, extraInline []trinoparser.FunctionSpecification, callerSchema string, deps []ViewDependency, cat *trinoFuncCatalog) ([]functionCallRef, []*ast.QualifiedName) {
	w := &trinoASTWalker{callerSchema: callerSchema, deps: deps, cat: cat}
	w.walkQueryWithInline(q, extraInline, nil, nil)
	return w.calls, w.unqualTables
}

func (w *trinoASTWalker) lookupTableRelationCols(qn *ast.QualifiedName, ctes map[string]struct{}) (string, map[string]string) {
	if qn == nil || len(qn.Parts) == 0 {
		return "", nil
	}
	if len(qn.Parts) == 1 {
		tblLower := strings.ToLower(qn.Parts[0].Value)
		if _, isCTE := ctes[tblLower]; isCTE {
			if w.cteCols != nil {
				return tblLower, w.cteCols[tblLower]
			}
			return tblLower, nil
		}
	}
	if w.cat == nil || len(w.cat.tableCols) == 0 {
		return "", nil
	}
	var schema, tbl string
	switch len(qn.Parts) {
	case 1:
		tbl = qn.Parts[0].Value
		schema = resolveUnqualifiedRelationSchema(tbl, w.callerSchema, w.deps)
	case 2:
		schema, tbl = qn.Parts[0].Value, qn.Parts[1].Value
	case 3:
		if w.cat.dbName != "" && !strings.EqualFold(qn.Parts[0].Value, w.cat.dbName) {
			return "", nil
		}
		schema, tbl = qn.Parts[1].Value, qn.Parts[2].Value
	default:
		return "", nil
	}
	if schema == "" {
		schema = w.callerSchema
	}
	cols := w.cat.tableCols[objectKey(strings.ToLower(schema), strings.ToLower(tbl))]
	return strings.ToLower(tbl), cols
}

func (w *trinoASTWalker) bindRelationCols(alias string, cols map[string]string, curEnv map[string]string) {
	if len(cols) == 0 {
		return
	}
	if alias != "" && w.aliasCols != nil {
		w.aliasCols[strings.ToLower(alias)] = cols
	}
	if curEnv != nil && w.fromSeen != nil {
		for col, typ := range cols {
			if _, seen := w.fromSeen[col]; seen {
				delete(curEnv, col)
			} else {
				w.fromSeen[col] = struct{}{}
				curEnv[col] = typ
			}
		}
	}
}

func (w *trinoASTWalker) recordUnqualTable(qn *ast.QualifiedName, ctes map[string]struct{}) {
	if qn == nil || len(qn.Parts) != 1 {
		return
	}
	if _, isCTE := ctes[strings.ToLower(qn.Parts[0].Value)]; isCTE {
		return
	}
	w.unqualTables = append(w.unqualTables, qn)
}

func projColsToMap(proj []trinoProjCol, colAliases []*ast.Identifier) map[string]string {
	if len(proj) == 0 {
		return nil
	}
	cols := make(map[string]string, len(proj))
	for idx, pc := range proj {
		colName := pc.name
		if idx < len(colAliases) && colAliases[idx] != nil && colAliases[idx].Value != "" {
			colName = strings.ToLower(colAliases[idx].Value)
		}
		if colName != "" && pc.typ != "" {
			cols[colName] = pc.typ
		}
	}
	return cols
}

func (w *trinoASTWalker) walkRelation(rel trinoparser.Relation, curEnv map[string]string, ctes map[string]struct{}) {
	w.walkRelationBound(rel, curEnv, ctes, true)
}

func (w *trinoASTWalker) walkRelationBound(rel trinoparser.Relation, curEnv map[string]string, ctes map[string]struct{}, bind bool) []trinoProjCol {
	if rel == nil {
		return nil
	}
	switch r := rel.(type) {
	case *trinoparser.Join:
		w.walkRelation(r.Left, curEnv, ctes)
		w.walkRelation(r.Right, curEnv, ctes)
		w.walkExpr(r.On, curEnv, ctes)
	case *trinoparser.AliasedRelation:
		var proj []trinoProjCol
		if tr, ok := r.Inner.(*trinoparser.TableRelation); ok {
			w.recordUnqualTable(tr.Name, ctes)
			if tr.Period != nil {
				w.walkExpr(tr.Period.End, curEnv, ctes)
			}
			defAlias, cols := w.lookupTableRelationCols(tr.Name, ctes)
			alias := defAlias
			if r.Alias != nil && r.Alias.Value != "" {
				alias = r.Alias.Value
			}
			w.bindRelationCols(alias, cols, curEnv)
		} else {
			proj = w.walkRelationBound(r.Inner, curEnv, ctes, false)
			alias := ""
			if r.Alias != nil {
				alias = r.Alias.Value
			}
			if len(r.ColumnAliases) > 0 && len(proj) > 0 {
				realiased := make([]trinoProjCol, len(proj))
				copy(realiased, proj)
				for idx := range realiased {
					if idx < len(r.ColumnAliases) && r.ColumnAliases[idx] != nil && r.ColumnAliases[idx].Value != "" {
						realiased[idx].name = strings.ToLower(r.ColumnAliases[idx].Value)
					}
				}
				proj = realiased
			}
			if bind {
				w.bindRelationCols(alias, projColsToMap(proj, nil), curEnv)
			}
		}
		if r.Sample != nil {
			w.walkExpr(r.Sample.Percentage, curEnv, ctes)
		}
		return proj
	case *trinoparser.TableRelation:
		w.recordUnqualTable(r.Name, ctes)
		if r.Period != nil {
			w.walkExpr(r.Period.End, curEnv, ctes)
		}
		alias, cols := w.lookupTableRelationCols(r.Name, ctes)
		if bind {
			w.bindRelationCols(alias, cols, curEnv)
		}
	case *trinoparser.NearestRelation:
		w.walkRelation(r.From, curEnv, ctes)
		w.walkExpr(r.Where, curEnv, ctes)
		w.walkExpr(r.Match, curEnv, ctes)
	case *trinoparser.SubqueryRelation:
		_, proj := w.walkQueryWithCols(r.Query, nil, curEnv, ctes)
		if bind {
			w.bindRelationCols("", projColsToMap(proj, nil), curEnv)
		}
		return proj
	case *trinoparser.LateralRelation:
		_, proj := w.walkQueryWithCols(r.Query, nil, curEnv, ctes)
		if bind {
			w.bindRelationCols("", projColsToMap(proj, nil), curEnv)
		}
		return proj
	case *trinoparser.UnnestRelation:
		for _, ex := range r.Exprs {
			w.walkExpr(ex, curEnv, ctes)
		}
	case *trinoparser.ParenRelation:
		return w.walkRelationBound(r.Inner, curEnv, ctes, bind)
	case *trinoparser.TableFunctionRelation:
		if r.Call != nil {
			for _, arg := range r.Call.Args {
				w.walkExpr(arg.Expr, curEnv, ctes)
				if arg.Table != nil {
					w.recordUnqualTable(arg.Table.Name, ctes)
					w.walkQuery(arg.Table.Query, curEnv, ctes)
					for _, pb := range arg.Table.PartitionBy {
						w.walkExpr(pb, curEnv, ctes)
					}
					for _, ob := range arg.Table.OrderBy {
						w.walkExpr(ob.Expr, curEnv, ctes)
					}
				}
			}
		}
	case *trinoparser.JSONTableRelation:
		w.walkJSONTable(r, curEnv, ctes)
	case *trinoparser.PatternRecognition:
		w.walkRelation(r.Input, curEnv, ctes)
		for _, pb := range r.PartitionBy {
			w.walkExpr(pb, curEnv, ctes)
		}
		for _, ob := range r.OrderBy {
			w.walkExpr(ob.Expr, curEnv, ctes)
		}
		for _, m := range r.Measures {
			w.walkExpr(m.Expr, curEnv, ctes)
		}
		for _, d := range r.Definitions {
			w.walkExpr(d.Condition, curEnv, ctes)
		}
	}
	return nil
}

func (w *trinoASTWalker) walkJSONTable(r *trinoparser.JSONTableRelation, curEnv map[string]string, ctes map[string]struct{}) {
	if r == nil {
		return
	}
	w.walkJSONPath(r.Path, curEnv, ctes)
	w.walkJSONTableColumns(r.Columns, curEnv, ctes)
}

func (w *trinoASTWalker) walkJSONTableColumns(cols []trinoparser.JSONTableColumn, curEnv map[string]string, ctes map[string]struct{}) {
	for _, col := range cols {
		if col.OnEmpty != nil {
			w.walkExpr(col.OnEmpty.Default, curEnv, ctes)
		}
		if col.OnError != nil {
			w.walkExpr(col.OnError.Default, curEnv, ctes)
		}
		if len(col.Columns) > 0 {
			w.walkJSONTableColumns(col.Columns, curEnv, ctes)
		}
	}
}

func (w *trinoASTWalker) walkQueryNodeCols(qn trinoparser.QueryNode, curEnv map[string]string, ctes map[string]struct{}) (string, []trinoProjCol) {
	if qn == nil {
		return "", nil
	}
	switch n := qn.(type) {
	case *trinoparser.QuerySpec:
		specEnv := make(map[string]string, len(curEnv))
		for k, v := range curEnv {
			specEnv[k] = v
		}
		prevAliasCols := w.aliasCols
		specAliasCols := make(map[string]map[string]string, len(prevAliasCols))
		for k, v := range prevAliasCols {
			specAliasCols[k] = v
		}
		prevFromSeen := w.fromSeen
		w.aliasCols = specAliasCols
		w.fromSeen = make(map[string]struct{})
		defer func() {
			w.aliasCols = prevAliasCols
			w.fromSeen = prevFromSeen
		}()

		for _, rel := range n.From {
			w.walkRelation(rel, specEnv, ctes)
		}
		w.walkExpr(n.Where, specEnv, ctes)
		if n.GroupBy != nil {
			for _, el := range n.GroupBy.Elements {
				for _, ex := range el.Exprs {
					w.walkExpr(ex, specEnv, ctes)
				}
				for _, set := range el.Sets {
					for _, ex := range set {
						w.walkExpr(ex, specEnv, ctes)
					}
				}
			}
		}
		w.walkExpr(n.Having, specEnv, ctes)
		for _, win := range n.Windows {
			w.walkWindowParts(win.Spec.PartitionBy, win.Spec.OrderBy, win.Spec.Frame, specEnv, ctes)
		}
		var singleType string
		proj := make([]trinoProjCol, 0, len(n.Items))
		for idx, item := range n.Items {
			it := w.walkExpr(item.Expr, specEnv, ctes)
			if idx == 0 && len(n.Items) == 1 && item.Kind == trinoparser.SelectSingle {
				singleType = it
			}
			if item.Kind == trinoparser.SelectSingle {
				colName := ""
				switch {
				case item.Alias != nil && item.Alias.Value != "":
					colName = strings.ToLower(item.Alias.Value)
				case item.Expr != nil:
					switch e := item.Expr.(type) {
					case *trinoparser.ColumnRef:
						if e.Name != nil {
							colName = strings.ToLower(e.Name.Value)
						}
					case *trinoparser.Dereference:
						if e.FieldName != nil {
							colName = strings.ToLower(e.FieldName.Value)
						}
					}
				}
				proj = append(proj, trinoProjCol{name: colName, typ: it})
			}
		}
		return singleType, proj
	case *trinoparser.TableQuery:
		w.recordUnqualTable(n.Name, ctes)
	case *trinoparser.SetOperation:
		lt, lcols := w.walkQueryNodeCols(n.Left, curEnv, ctes)
		rt, rcols := w.walkQueryNodeCols(n.Right, curEnv, ctes)
		if len(lcols) > 0 && len(lcols) == len(rcols) {
			for i := range lcols {
				lcols[i].typ = combineTrinoTypes(lcols[i].typ, rcols[i].typ)
			}
		}
		return combineTrinoTypes(lt, rt), lcols
	case *trinoparser.ParenQuery:
		return w.walkQueryWithCols(n.Inner, nil, curEnv, ctes)
	case *trinoparser.ValuesQuery:
		for _, row := range n.Rows {
			w.walkExpr(row, curEnv, ctes)
		}
	}
	return "", nil
}

func (w *trinoASTWalker) walkQuery(q *trinoparser.Query, curEnv map[string]string, ctes map[string]struct{}) string {
	return w.walkQueryWithInline(q, nil, curEnv, ctes)
}

func (w *trinoASTWalker) walkQueryWithInline(q *trinoparser.Query, extraInline []trinoparser.FunctionSpecification, curEnv map[string]string, ctes map[string]struct{}) string {
	ret, _ := w.walkQueryWithCols(q, extraInline, curEnv, ctes)
	return ret
}

func (w *trinoASTWalker) walkQueryWithCols(q *trinoparser.Query, extraInline []trinoparser.FunctionSpecification, curEnv map[string]string, ctes map[string]struct{}) (string, []trinoProjCol) {
	if q == nil {
		return "", nil
	}
	w.hasQuery = true
	totalInline := len(extraInline) + len(q.Functions)
	if totalInline > 0 {
		prevInline := w.inlineFuncs
		nextInline := make(map[string][]trinoInlineFunc, len(prevInline)+totalInline)
		for k, v := range prevInline {
			nextInline[k] = slices.Clone(v)
		}
		registerInline := func(specs []trinoparser.FunctionSpecification) {
			for _, fnSpec := range specs {
				if fnSpec.Declaration.Name == nil || len(fnSpec.Declaration.Name.Parts) != 1 {
					continue
				}
				key := strings.ToLower(fnSpec.Declaration.Name.Parts[0].Value)
				params := make([]trinoFuncParam, len(fnSpec.Declaration.Parameters))
				for i, p := range fnSpec.Declaration.Parameters {
					name := ""
					if p.Name != nil {
						name = p.Name.Value
					}
					params[i] = trinoFuncParam{
						name:     name,
						typeName: normalizeTrinoDataType(p.Type),
					}
				}
				nextInline[key] = append(nextInline[key], trinoInlineFunc{
					params:     params,
					returnType: normalizeTrinoDataType(fnSpec.ReturnType),
				})
			}
		}
		registerInline(extraInline)
		registerInline(q.Functions)
		w.inlineFuncs = nextInline
		defer func() {
			w.inlineFuncs = prevInline
		}()
		walkInlineBodies := func(specs []trinoparser.FunctionSpecification) {
			for _, fnSpec := range specs {
				fnEnv := make(map[string]string, len(fnSpec.Declaration.Parameters))
				for _, p := range fnSpec.Declaration.Parameters {
					if p.Name != nil && p.Name.Value != "" {
						if pt := normalizeTrinoDataType(p.Type); pt != "" {
							fnEnv[strings.ToLower(p.Name.Value)] = pt
						}
					}
				}
				w.walkStmt(fnSpec.Body, fnEnv)
			}
		}
		walkInlineBodies(extraInline)
		walkInlineBodies(q.Functions)
	}
	scope := make(map[string]struct{}, len(ctes))
	for k := range ctes {
		scope[k] = struct{}{}
	}
	prevCTECols := w.cteCols
	nextCTECols := make(map[string]map[string]string, len(prevCTECols))
	for k, v := range prevCTECols {
		nextCTECols[k] = v
	}
	w.cteCols = nextCTECols
	defer func() {
		w.cteCols = prevCTECols
	}()
	if q.With != nil {
		for _, cte := range q.With.CTEs {
			if q.With.Recursive && cte.Name != nil && cte.Name.Value != "" {
				scope[strings.ToLower(cte.Name.Value)] = struct{}{}
			}
			_, proj := w.walkQueryWithCols(cte.Query, nil, curEnv, scope)
			if cte.Name != nil && cte.Name.Value != "" {
				cteKey := strings.ToLower(cte.Name.Value)
				scope[cteKey] = struct{}{}
				nextCTECols[cteKey] = projColsToMap(proj, cte.ColumnAliases)
			}
		}
	}
	ret, proj := w.walkQueryNodeCols(q.Body, curEnv, scope)
	for _, ob := range q.OrderBy {
		w.walkExpr(ob.Expr, curEnv, scope)
	}
	if q.Offset != nil {
		w.walkExpr(q.Offset.Count, curEnv, scope)
	}
	if q.Limit != nil {
		w.walkExpr(q.Limit.Count, curEnv, scope)
	}
	if q.Fetch != nil {
		w.walkExpr(q.Fetch.Count, curEnv, scope)
	}
	return ret, proj
}

func (w *trinoASTWalker) walkJSONPath(path *trinoparser.JSONPathInvocation, curEnv map[string]string, ctes map[string]struct{}) {
	if path == nil {
		return
	}
	if path.Input != nil {
		w.walkExpr(path.Input.Expr, curEnv, ctes)
	}
	for _, p := range path.Passing {
		if p.Value != nil {
			w.walkExpr(p.Value.Expr, curEnv, ctes)
		}
	}
}

func (w *trinoASTWalker) walkWindowParts(partitionBy []trinoparser.Expr, orderBy []trinoparser.SortItem, frame *trinoparser.WindowFrame, curEnv map[string]string, ctes map[string]struct{}) {
	for _, pb := range partitionBy {
		w.walkExpr(pb, curEnv, ctes)
	}
	for _, ob := range orderBy {
		w.walkExpr(ob.Expr, curEnv, ctes)
	}
	if frame == nil {
		return
	}
	w.walkExpr(frame.Start.Value, curEnv, ctes)
	if frame.End != nil {
		w.walkExpr(frame.End.Value, curEnv, ctes)
	}
	if frame.Pattern != nil {
		for _, m := range frame.Pattern.Measures {
			w.walkExpr(m.Expr, curEnv, ctes)
		}
		for _, d := range frame.Pattern.Definitions {
			w.walkExpr(d.Condition, curEnv, ctes)
		}
	}
}

func (w *trinoASTWalker) walkExpr(expr trinoparser.Expr, curEnv map[string]string, ctes map[string]struct{}) string {
	if expr == nil {
		return ""
	}
	switch e := expr.(type) {
	case *trinoparser.NamedArgument:
		return w.walkExpr(e.Value, curEnv, ctes)
	case *trinoparser.ColumnRef:
		if e.Name != nil {
			return curEnv[strings.ToLower(e.Name.Value)]
		}
	case *trinoparser.Dereference:
		baseType := w.walkExpr(e.Base, curEnv, ctes)
		if e.FieldName != nil && e.FieldName.Value != "" {
			target := strings.ToLower(e.FieldName.Value)
			if base, inner, ok := splitParameterizedType(baseType); ok && base == "row" {
				for _, part := range splitTopLevelComma(inner) {
					if name, typ := splitRowField(part); name == target {
						return typ
					}
				}
			}
			if cr, ok := e.Base.(*trinoparser.ColumnRef); ok && cr.Name != nil && w.aliasCols != nil {
				if cols := w.aliasCols[strings.ToLower(cr.Name.Value)]; cols != nil {
					if typ, ok := cols[target]; ok {
						return typ
					}
				}
			}
		}
	case *trinoparser.Subscript:
		baseType := w.walkExpr(e.Value, curEnv, ctes)
		w.walkExpr(e.Index, curEnv, ctes)
		if base, inner, ok := splitParameterizedType(baseType); ok {
			switch base {
			case "array":
				return inner
			case "map":
				if parts := splitTopLevelComma(inner); len(parts) == 2 {
					return parts[1]
				}
			}
		}
	case *trinoparser.Literal:
		switch e.Kind {
		case trinoparser.LiteralBool:
			return "boolean"
		case trinoparser.LiteralInteger:
			if v, err := strconv.ParseInt(e.Value, 10, 64); err == nil && v >= math.MinInt32 && v <= math.MaxInt32 {
				return "integer"
			}
			return "bigint"
		case trinoparser.LiteralDecimal:
			return "decimal"
		case trinoparser.LiteralDouble:
			return "double"
		case trinoparser.LiteralString:
			return "varchar"
		case trinoparser.LiteralBinary:
			return "varbinary"
		}
	case *trinoparser.TypeConstructor:
		return normalizeTrinoTypeName(e.Name)
	case *trinoparser.IntervalLiteral:
		return normalizeTrinoIntervalType(e.From, e.To)
	case *trinoparser.SpecialFuncExpr:
		switch strings.ToUpper(e.Name) {
		case "CURRENT_DATE":
			return "date"
		case "LOCALTIME":
			return "time"
		case "LOCALTIMESTAMP":
			return "timestamp"
		case "CURRENT_TIME":
			return "time with time zone"
		case "CURRENT_TIMESTAMP":
			return "timestamp with time zone"
		case "CURRENT_USER", "CURRENT_CATALOG", "CURRENT_SCHEMA", "CURRENT_PATH":
			return "varchar"
		}
	case *trinoparser.CastExpr:
		w.walkExpr(e.Expr, curEnv, ctes)
		return normalizeTrinoDataType(e.Type)
	case *trinoparser.UnaryExpr:
		if e.Op == "-" {
			if lit, ok := e.Operand.(*trinoparser.Literal); ok && lit.Kind == trinoparser.LiteralInteger {
				if v, err := strconv.ParseInt("-"+lit.Value, 10, 64); err == nil && v >= math.MinInt32 && v <= math.MaxInt32 {
					return "integer"
				}
				return "bigint"
			}
		}
		return w.walkExpr(e.Operand, curEnv, ctes)
	case *trinoparser.ParenExpr:
		return w.walkExpr(e.Expr, curEnv, ctes)
	case *trinoparser.NotExpr:
		w.walkExpr(e.Operand, curEnv, ctes)
		return "boolean"
	case *trinoparser.LogicalExpr:
		w.walkExpr(e.Left, curEnv, ctes)
		w.walkExpr(e.Right, curEnv, ctes)
		return "boolean"
	case *trinoparser.BinaryExpr:
		lt := w.walkExpr(e.Left, curEnv, ctes)
		rt := w.walkExpr(e.Right, curEnv, ctes)
		switch e.Op {
		case "||":
			return "varchar"
		case "+", "-", "*", "/", "%":
			return commonTrinoNumericType(lt, rt)
		}
	case *trinoparser.CaseExpr:
		w.walkExpr(e.Operand, curEnv, ctes)
		res := ""
		for _, wh := range e.Whens {
			w.walkExpr(wh.Cond, curEnv, ctes)
			res = combineTrinoTypes(res, w.walkExpr(wh.Result, curEnv, ctes))
		}
		res = combineTrinoTypes(res, w.walkExpr(e.Else, curEnv, ctes))
		return res
	case *trinoparser.AtTimeZoneExpr:
		t := w.walkExpr(e.Value, curEnv, ctes)
		w.walkExpr(e.Zone, curEnv, ctes)
		return zonedTrinoDateTimeType(t)
	case *trinoparser.AtLocalExpr:
		return zonedTrinoDateTimeType(w.walkExpr(e.Value, curEnv, ctes))
	case *trinoparser.ArrayConstructor:
		elemType := ""
		for _, el := range e.Elements {
			elemType = combineTrinoTypes(elemType, w.walkExpr(el, curEnv, ctes))
		}
		if elemType != "" {
			return "array(" + elemType + ")"
		}
		return "array"
	case *trinoparser.RowConstructor:
		fieldTypes := make([]string, len(e.Elements))
		allKnown := len(e.Elements) > 0
		for j, el := range e.Elements {
			ft := w.walkExpr(el, curEnv, ctes)
			if ft == "" {
				allKnown = false
			}
			fieldTypes[j] = ft
		}
		if allKnown {
			return "row(" + strings.Join(fieldTypes, ",") + ")"
		}
		return "row"
	case *trinoparser.OverlayExpr:
		w.walkExpr(e.Source, curEnv, ctes)
		w.walkExpr(e.Replacing, curEnv, ctes)
		w.walkExpr(e.From, curEnv, ctes)
		w.walkExpr(e.For, curEnv, ctes)
		return "varchar"
	case *trinoparser.ExtractExpr:
		w.walkExpr(e.Source, curEnv, ctes)
		return "bigint"
	case *trinoparser.SubstringExpr:
		w.walkExpr(e.Source, curEnv, ctes)
		w.walkExpr(e.From, curEnv, ctes)
		w.walkExpr(e.For, curEnv, ctes)
		return "varchar"
	case *trinoparser.ComparisonExpr:
		w.walkExpr(e.Left, curEnv, ctes)
		w.walkExpr(e.Right, curEnv, ctes)
		return "boolean"
	case *trinoparser.QuantifiedComparisonExpr:
		w.hasQuery = true
		w.walkExpr(e.Left, curEnv, ctes)
		w.walkExpr(e.Subquery, curEnv, ctes)
		return "boolean"
	case *trinoparser.BetweenExpr:
		w.walkExpr(e.Value, curEnv, ctes)
		w.walkExpr(e.Lower, curEnv, ctes)
		w.walkExpr(e.Upper, curEnv, ctes)
		return "boolean"
	case *trinoparser.InListExpr:
		w.walkExpr(e.Value, curEnv, ctes)
		for _, item := range e.List {
			w.walkExpr(item, curEnv, ctes)
		}
		return "boolean"
	case *trinoparser.InSubqueryExpr:
		w.hasQuery = true
		w.walkExpr(e.Value, curEnv, ctes)
		w.walkExpr(e.Subquery, curEnv, ctes)
		return "boolean"
	case *trinoparser.LikeExpr:
		w.walkExpr(e.Value, curEnv, ctes)
		w.walkExpr(e.Pattern, curEnv, ctes)
		w.walkExpr(e.Escape, curEnv, ctes)
		return "boolean"
	case *trinoparser.IsNullExpr:
		w.walkExpr(e.Value, curEnv, ctes)
		return "boolean"
	case *trinoparser.IsDistinctFromExpr:
		w.walkExpr(e.Value, curEnv, ctes)
		w.walkExpr(e.Right, curEnv, ctes)
		return "boolean"
	case *trinoparser.IsBooleanExpr:
		w.walkExpr(e.Value, curEnv, ctes)
		return "boolean"
	case *trinoparser.MatchExpr:
		w.hasQuery = true
		w.walkExpr(e.Value, curEnv, ctes)
		w.walkExpr(e.Subquery, curEnv, ctes)
		return "boolean"
	case *trinoparser.UniqueExpr:
		w.hasQuery = true
		w.walkExpr(e.Subquery, curEnv, ctes)
		return "boolean"
	case *trinoparser.SubqueryExpr:
		w.hasQuery = true
		var subType string
		if e.Query != nil {
			subType = w.walkQuery(e.Query, curEnv, ctes)
		} else if e.RawText != "" {
			w.calls = append(w.calls, scanFunctionCallRefs(e.RawText)...)
		}
		if e.Kind == trinoparser.SubqueryExists || e.Kind == trinoparser.SubqueryUnique {
			return "boolean"
		}
		if e.Kind == trinoparser.SubqueryScalar {
			return subType
		}
	case *trinoparser.TrimExpr:
		w.walkExpr(e.Char, curEnv, ctes)
		w.walkExpr(e.Source, curEnv, ctes)
		return "varchar"
	case *trinoparser.NormalizeExpr:
		w.walkExpr(e.Source, curEnv, ctes)
		return "varchar"
	case *trinoparser.PositionExpr:
		w.walkExpr(e.Needle, curEnv, ctes)
		w.walkExpr(e.Haystack, curEnv, ctes)
		return "bigint"
	case *trinoparser.ListaggExpr:
		w.walkExpr(e.Arg, curEnv, ctes)
		w.walkExpr(e.Separator, curEnv, ctes)
		for _, ob := range e.WithinGroupBy {
			w.walkExpr(ob.Expr, curEnv, ctes)
		}
		w.walkExpr(e.Filter, curEnv, ctes)
		return "varchar"
	case *trinoparser.LambdaExpr:
		lambdaEnv := curEnv
		if len(e.Params) > 0 {
			lambdaEnv = make(map[string]string, len(curEnv))
			for k, v := range curEnv {
				lambdaEnv[k] = v
			}
			for _, p := range e.Params {
				if p != nil {
					delete(lambdaEnv, strings.ToLower(p.Value))
				}
			}
		}
		w.walkExpr(e.Body, lambdaEnv, ctes)
		return "function"
	case *trinoparser.JSONExistsExpr:
		w.walkJSONPath(e.Path, curEnv, ctes)
		return "boolean"
	case *trinoparser.JSONValueFunc:
		w.walkJSONPath(e.Path, curEnv, ctes)
		if e.OnEmpty != nil {
			w.walkExpr(e.OnEmpty.Default, curEnv, ctes)
		}
		if e.OnError != nil {
			w.walkExpr(e.OnError.Default, curEnv, ctes)
		}
		return normalizeTrinoDataType(e.Returning)
	case *trinoparser.JSONQueryFunc:
		w.walkJSONPath(e.Path, curEnv, ctes)
		return normalizeTrinoDataType(e.Returning)
	case *trinoparser.JSONObjectExpr:
		for _, m := range e.Members {
			w.walkExpr(m.Key, curEnv, ctes)
			if m.Value != nil {
				w.walkExpr(m.Value.Expr, curEnv, ctes)
			}
		}
		return normalizeTrinoDataType(e.Returning)
	case *trinoparser.JSONArrayExpr:
		for _, el := range e.Elements {
			if el != nil {
				w.walkExpr(el.Expr, curEnv, ctes)
			}
		}
		return normalizeTrinoDataType(e.Returning)
	case *trinoparser.MethodCall:
		w.walkExpr(e.Value, curEnv, ctes)
		for _, a := range e.Args {
			w.walkExpr(a, curEnv, ctes)
		}
	case *trinoparser.StaticMethodCall:
		for _, a := range e.Args {
			w.walkExpr(a, curEnv, ctes)
		}
	case *trinoparser.FuncCall:
		argTypes := make([]string, len(e.Args))
		var argNames []string
		for j, a := range e.Args {
			if na, ok := a.(*trinoparser.NamedArgument); ok && na != nil {
				argTypes[j] = w.walkExpr(na.Value, curEnv, ctes)
				if na.Name != nil && na.Name.Value != "" {
					if argNames == nil {
						argNames = make([]string, len(e.Args))
					}
					argNames[j] = na.Name.Value
				}
			} else {
				argTypes[j] = w.walkExpr(a, curEnv, ctes)
			}
		}
		for _, ob := range e.OrderBy {
			w.walkExpr(ob.Expr, curEnv, ctes)
		}
		w.walkExpr(e.Filter, curEnv, ctes)
		if e.Over != nil {
			w.walkWindowParts(e.Over.PartitionBy, e.Over.OrderBy, e.Over.Frame, curEnv, ctes)
		}
		var nameStart, nameEnd int
		var catalog, schema, name string
		if e.Name != nil && len(e.Name.Parts) > 0 {
			nameStart, nameEnd = e.Name.Loc.Start, e.Name.Loc.End
			parts := e.Name.Parts
			switch len(parts) {
			case 1:
				name = parts[0].Value
			case 2:
				schema, name = parts[0].Value, parts[1].Value
			case 3:
				catalog, schema, name = parts[0].Value, parts[1].Value, parts[2].Value
			default:
				return ""
			}
		}
		if name != "" {
			ref := functionCallRef{
				catalog:   catalog,
				schema:    schema,
				name:      name,
				nameStart: nameStart,
				nameEnd:   nameEnd,
				argNames:  argNames,
				argTypes:  argTypes,
				argCount:  len(argTypes),
			}
			if retType, isInline := w.resolveInlineFunc(ref); isInline {
				return retType
			}
			w.calls = append(w.calls, ref)
			if matches := w.cat.resolveCall(w.callerSchema, ref); len(matches) > 0 {
				ret := w.cat.parsed[matches[0]].returnType
				same := ret != ""
				for _, m := range matches[1:] {
					if w.cat.parsed[m].returnType != ret {
						same = false
						break
					}
				}
				if same {
					return ret
				}
			}
			switch strings.ToLower(name) {
			case "length", "cardinality", "strpos", "count":
				return "bigint"
			case "lower", "upper", "trim", "ltrim", "rtrim", "concat", "substr", "substring", "replace", "reverse":
				return "varchar"
			case "coalesce", "nullif":
				for _, at := range argTypes {
					if at != "" {
						return at
					}
				}
			}
		}
	}
	return ""
}

func (w *trinoASTWalker) resolveInlineFunc(call functionCallRef) (string, bool) {
	if call.catalog != "" || call.schema != "" || len(w.inlineFuncs) == 0 {
		return "", false
	}
	cands := w.inlineFuncs[strings.ToLower(call.name)]
	if len(cands) == 0 {
		return "", false
	}
	var matches []trinoInlineFunc
	for _, cand := range cands {
		if len(cand.params) != call.argCount {
			continue
		}
		argTypes, ok := reorderTrinoCallArgs(cand.params, call)
		if !ok {
			continue
		}
		compat := true
		for i, argType := range argTypes {
			if _, ok := trinoTypeCompatScore(argType, cand.params[i].typeName); !ok {
				compat = false
				break
			}
		}
		if compat {
			matches = append(matches, cand)
		}
	}
	if len(matches) == 0 {
		return "", false
	}
	ret := matches[0].returnType
	for _, m := range matches[1:] {
		if m.returnType != ret {
			return "", true
		}
	}
	return ret, true
}

func (w *trinoASTWalker) walkStmt(stmt trinoparser.RoutineStatement, curEnv map[string]string) {
	if stmt == nil {
		return
	}
	switch s := stmt.(type) {
	case *trinoparser.ReturnStatement:
		w.walkExpr(s.Value, curEnv, nil)
	case *trinoparser.AssignmentStatement:
		w.walkExpr(s.Value, curEnv, nil)
	case *trinoparser.IfStatement:
		w.walkExpr(s.Condition, curEnv, nil)
		for _, st := range s.Then {
			w.walkStmt(st, curEnv)
		}
		for _, elif := range s.ElseIfs {
			w.walkExpr(elif.Condition, curEnv, nil)
			for _, st := range elif.Body {
				w.walkStmt(st, curEnv)
			}
		}
		for _, st := range s.Else {
			w.walkStmt(st, curEnv)
		}
	case *trinoparser.CaseStatement:
		w.walkExpr(s.Operand, curEnv, nil)
		for _, wh := range s.Whens {
			w.walkExpr(wh.Condition, curEnv, nil)
			for _, st := range wh.Body {
				w.walkStmt(st, curEnv)
			}
		}
		for _, st := range s.Else {
			w.walkStmt(st, curEnv)
		}
	case *trinoparser.LoopStatement:
		for _, st := range s.Body {
			w.walkStmt(st, curEnv)
		}
	case *trinoparser.WhileStatement:
		w.walkExpr(s.Condition, curEnv, nil)
		for _, st := range s.Body {
			w.walkStmt(st, curEnv)
		}
	case *trinoparser.RepeatStatement:
		for _, st := range s.Body {
			w.walkStmt(st, curEnv)
		}
		w.walkExpr(s.Condition, curEnv, nil)
	case *trinoparser.CompoundStatement:
		childEnv := make(map[string]string, len(curEnv)+len(s.Declarations))
		for k, v := range curEnv {
			childEnv[k] = v
		}
		for _, decl := range s.Declarations {
			w.walkExpr(decl.Default, childEnv, nil)
			dt := normalizeTrinoDataType(decl.Type)
			for _, n := range decl.Names {
				if n != nil && dt != "" {
					childEnv[strings.ToLower(n.Value)] = dt
				}
			}
		}
		for _, st := range s.Body {
			w.walkStmt(st, childEnv)
		}
	}
}

// scanFunctionCallRefs extracts unqualified and schema-qualified call sites
// inside sql, ignoring comments and string literals.
func scanFunctionCallRefs(sql string) []functionCallRef {
	return scanFunctionCallTokens(sql, false)
}

func scanFunctionBodyCallRefs(definition string) []functionCallRef {
	return scanFunctionCallTokens(definition, true)
}

func scanFunctionCallTokens(sql string, skipHeader bool) []functionCallRef {
	type token struct {
		isIdent bool
		text    string
		punct   byte
	}
	var toks []token
	for i := 0; i < len(sql); {
		ch := sql[i]
		switch {
		case ch == '-' && i+1 < len(sql) && sql[i+1] == '-':
			i += 2
			for i < len(sql) && sql[i] != '\n' {
				i++
			}
		case ch == '/' && i+1 < len(sql) && sql[i+1] == '*':
			i += 2
			for i+1 < len(sql) && (sql[i] != '*' || sql[i+1] != '/') {
				i++
			}
			if i+1 < len(sql) {
				i += 2
			} else {
				i = len(sql)
			}
		case ch == '\'':
			i++
			for i < len(sql) {
				if sql[i] == '\'' {
					if i+1 < len(sql) && sql[i+1] == '\'' {
						i += 2
						continue
					}
					i++
					break
				}
				i++
			}
		case ch == '"':
			i++
			var sb strings.Builder
			for i < len(sql) {
				if sql[i] == '"' {
					if i+1 < len(sql) && sql[i+1] == '"' {
						sb.WriteByte('"')
						i += 2
						continue
					}
					i++
					break
				}
				sb.WriteByte(sql[i])
				i++
			}
			toks = append(toks, token{isIdent: true, text: sb.String()})
		case isIdentStartChar(ch):
			start := i
			i++
			for i < len(sql) && isIdentContChar(sql[i]) {
				i++
			}
			toks = append(toks, token{isIdent: true, text: sql[start:i]})
		case ch == '.' || ch == '(' || ch == ')' || ch == ',':
			toks = append(toks, token{punct: ch})
			i++
		default:
			i++
		}
	}
	if skipHeader {
		// Skip `CREATE [OR REPLACE] FUNCTION name(params)` so the declared
		// signature is not treated as a self-recursive call in unparsed DDL.
		depth := 0
		for k := 0; k < len(toks); k++ {
			if toks[k].punct == '(' {
				depth++
			} else if toks[k].punct == ')' && depth > 0 {
				depth--
				if depth == 0 {
					toks = toks[k+1:]
					break
				}
			}
		}
	}

	countArgs := func(openIdx int) int {
		if openIdx+1 < len(toks) && toks[openIdx+1].punct == ')' {
			return 0
		}
		depth := 0
		args := 1
		for j := openIdx; j < len(toks); j++ {
			switch toks[j].punct {
			case '(':
				depth++
			case ')':
				depth--
				if depth == 0 {
					return args
				}
			case ',':
				if depth == 1 {
					args++
				}
			}
		}
		return args
	}

	var calls []functionCallRef
	for i := 0; i+1 < len(toks); i++ {
		if i > 0 && toks[i-1].punct == '.' {
			continue
		}
		if i+5 < len(toks) && toks[i].isIdent && toks[i+1].punct == '.' && toks[i+2].isIdent && toks[i+3].punct == '.' && toks[i+4].isIdent && toks[i+5].punct == '(' {
			calls = append(calls, functionCallRef{catalog: toks[i].text, schema: toks[i+2].text, name: toks[i+4].text, argCount: countArgs(i + 5)})
			continue
		}
		if i+3 < len(toks) && toks[i].isIdent && toks[i+1].punct == '.' && toks[i+2].isIdent && toks[i+3].punct == '(' {
			calls = append(calls, functionCallRef{schema: toks[i].text, name: toks[i+2].text, argCount: countArgs(i + 3)})
			continue
		}
		if toks[i].isIdent && toks[i+1].punct == '(' {
			calls = append(calls, functionCallRef{name: toks[i].text, argCount: countArgs(i + 1)})
		}
	}
	return calls
}

func isIdentStartChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isIdentContChar(c byte) bool {
	return isIdentStartChar(c) || (c >= '0' && c <= '9')
}

// writeCreateTable appends a single CREATE TABLE statement for table (qualified
// by schemaName) to buf. It does not append the inter-statement blank line; the
// caller does.
func writeCreateTable(buf *strings.Builder, schemaName string, table *TableMetadata) {
	_ = writeCreateTableWithCatalog(buf, "", schemaName, table, nil)
}

func writeCreateTableWithCatalog(buf *strings.Builder, dbName, schemaName string, table *TableMetadata, funcCat *trinoFuncCatalog) error {
	seenCols := make(map[string]struct{}, len(table.Columns))
	for _, column := range table.Columns {
		if column == nil {
			continue
		}
		colKey := strings.ToLower(column.Name)
		if _, exists := seenCols[colKey]; exists {
			return fmt.Errorf("table %s.%s has duplicate column %s", schemaName, table.Name, quoteIdentifier(column.Name))
		}
		seenCols[colKey] = struct{}{}
		if strings.TrimSpace(column.Type) == "" {
			return fmt.Errorf("table %s.%s column %s has no data type", schemaName, table.Name, column.Name)
		}
	}
	fmt.Fprintf(buf, "CREATE TABLE IF NOT EXISTS %s.%s (\n", quoteIdentifier(schemaName), quoteIdentifier(table.Name))

	columnDefs := make([]string, 0, len(table.Columns))
	for _, column := range table.Columns {
		if column == nil {
			continue
		}
		colDef := fmt.Sprintf("    %s %s", quoteIdentifier(column.Name), column.Type)
		if column.Default != "" {
			colDef += " DEFAULT " + qualifyTrinoDefault(dbName, schemaName, column.Default, funcCat)
		}
		if !column.Nullable {
			colDef += " NOT NULL"
		}
		if column.Comment != "" {
			colDef += fmt.Sprintf(" COMMENT '%s'", escapeStringLiteral(column.Comment))
		}
		columnDefs = append(columnDefs, colDef)
	}

	buf.WriteString(strings.Join(columnDefs, ",\n"))
	if table.Comment != "" {
		fmt.Fprintf(buf, "\n) COMMENT '%s';", escapeStringLiteral(table.Comment))
		return nil
	}
	buf.WriteString("\n);")
	return nil
}

func escapeStringLiteral(s string) string {
	return strings.ReplaceAll(s, `'`, `''`)
}

// appendTerminator appends ";" when stmt does not already end in a
// statement terminator, placing it on a new line if stmt ends inside a "--"
// line comment.
func appendTerminator(stmt string) string {
	var (
		inSingleQuote  bool
		inDoubleQuote  bool
		inLineComment  bool
		inBlockComment bool
		lastCodeByte   byte
	)
	for i := 0; i < len(stmt); {
		ch := stmt[i]
		if inLineComment {
			if ch == '\n' {
				inLineComment = false
			}
			i++
			continue
		}
		if inBlockComment {
			if ch == '*' && i+1 < len(stmt) && stmt[i+1] == '/' {
				inBlockComment = false
				i += 2
				continue
			}
			i++
			continue
		}
		if inSingleQuote {
			if ch == '\'' {
				if i+1 < len(stmt) && stmt[i+1] == '\'' {
					i += 2
					continue
				}
				inSingleQuote = false
				lastCodeByte = '\''
			}
			i++
			continue
		}
		if inDoubleQuote {
			if ch == '"' {
				if i+1 < len(stmt) && stmt[i+1] == '"' {
					i += 2
					continue
				}
				inDoubleQuote = false
				lastCodeByte = '"'
			}
			i++
			continue
		}
		switch {
		case ch == '-' && i+1 < len(stmt) && stmt[i+1] == '-':
			inLineComment = true
			i += 2
		case ch == '/' && i+1 < len(stmt) && stmt[i+1] == '*':
			inBlockComment = true
			i += 2
		case ch == '\'':
			inSingleQuote = true
			lastCodeByte = ch
			i++
		case ch == '"':
			inDoubleQuote = true
			lastCodeByte = ch
			i++
		case ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n':
			i++
		default:
			lastCodeByte = ch
			i++
		}
	}
	if lastCodeByte == ';' {
		return stmt
	}
	if inLineComment {
		return stmt + "\n;"
	}
	return stmt + ";"
}

// writeCreateView appends a single CREATE OR REPLACE VIEW (or CREATE
// MATERIALIZED VIEW IF NOT EXISTS) statement for view (qualified by schemaName)
// to buf. It does not append the inter-statement blank line; the caller does.
func writeCreateView(buf *strings.Builder, schemaName string, view *ViewMetadata, materialized bool) error {
	if materialized && view.Comment != "" {
		return fmt.Errorf("materialized view %s.%s comment cannot be restored on pinned Trino 482", schemaName, view.Name)
	}
	definition := strings.TrimSpace(view.Definition)
	if definition == "" {
		if materialized {
			return fmt.Errorf("materialized view %s.%s has no definition to restore", schemaName, view.Name)
		}
		return fmt.Errorf("view %s.%s has no definition to restore", schemaName, view.Name)
	}
	commentClause := ""
	if view.Comment != "" {
		commentClause = fmt.Sprintf(" COMMENT '%s'", escapeStringLiteral(view.Comment))
	}
	var stmt string
	if materialized {
		stmt = fmt.Sprintf("CREATE MATERIALIZED VIEW IF NOT EXISTS %s.%s%s AS\n%s", quoteIdentifier(schemaName), quoteIdentifier(view.Name), commentClause, definition)
	} else {
		stmt = fmt.Sprintf("CREATE OR REPLACE VIEW %s.%s%s AS\n%s", quoteIdentifier(schemaName), quoteIdentifier(view.Name), commentClause, definition)
	}
	buf.WriteString(appendTerminator(stmt))
	if !materialized {
		for _, col := range view.Columns {
			if col == nil || col.Comment == "" {
				continue
			}
			fmt.Fprintf(buf, "\nCOMMENT ON COLUMN %s.%s.%s IS '%s';", quoteIdentifier(schemaName), quoteIdentifier(view.Name), quoteIdentifier(col.Name), escapeStringLiteral(col.Comment))
		}
	}
	return nil
}

// orderedView is one ordinary or materialized view and the schema that owns
// it, in snapshot order.
type orderedView struct {
	schema       string
	view         *ViewMetadata
	materialized bool
}

// orderedViews returns metadata's views and materialized views in dependency
// order: every entry follows the views/materialized views its explicit
// Dependencies or parsed query span reference. Independent entries keep
// snapshot order. A definition the local parser cannot read contributes no
// edges rather than failing the dump; Trino accepted it, and a dependency the
// local parser cannot see is no reason to lose the statement.
func orderedViews(metadata *DatabaseSchemaMetadata) ([]orderedView, error) {
	var views []orderedView
	for _, schema := range metadata.Schemas {
		if schema == nil || schema.SkipDump {
			continue
		}
		for _, view := range schema.Views {
			if view == nil {
				continue
			}
			views = append(views, orderedView{schema: schema.Name, view: view})
		}
		for _, mv := range schema.MaterializedViews {
			if mv == nil {
				continue
			}
			views = append(views, orderedView{schema: schema.Name, view: mv, materialized: true})
		}
	}
	if len(views) == 0 {
		return views, nil
	}

	index := make(map[string]int, len(views))
	for i, view := range views {
		index[objectKey(strings.ToLower(view.schema), strings.ToLower(view.view.Name))] = i
	}
	dependencies := make([]map[int]struct{}, len(views))
	addDep := func(i int, depSchema, depTable string) {
		if depSchema == "" {
			depSchema = views[i].schema
		}
		dependency, ok := index[objectKey(strings.ToLower(depSchema), strings.ToLower(depTable))]
		if !ok {
			return
		}
		if dependencies[i] == nil {
			dependencies[i] = map[int]struct{}{}
		}
		dependencies[i][dependency] = struct{}{}
	}
	for i, view := range views {
		for _, dep := range view.view.Dependencies {
			addDep(i, dep.Schema, dep.Table)
		}
		span, err := analysis.GetQuerySpan(view.view.Definition)
		if err != nil {
			continue
		}
		for _, access := range span.AccessTables {
			if access.Catalog != "" && metadata.Name != "" && !strings.EqualFold(access.Catalog, metadata.Name) {
				continue
			}
			addDep(i, access.Schema, access.Table)
		}
	}
	indices, cyclic := topoOrder(len(views), dependencies)
	if len(cyclic) > 0 {
		return nil, fmt.Errorf("cyclic view dependency involving %s.%s cannot be restored", views[cyclic[0]].schema, views[cyclic[0]].view.Name)
	}
	ordered := make([]orderedView, len(indices))
	for i, idx := range indices {
		ordered[i] = views[idx]
	}
	return ordered, nil
}

// objectKey identifies a schema-scoped object by its owning schema and its name.
func objectKey(schema, name string) string {
	return schema + "\x00" + name
}

// topoOrder returns item indices ordered so each item follows its dependencies,
// keeping input order among independent items and returning indices where a cycle was broken.
func topoOrder(n int, dependencies []map[int]struct{}) ([]int, []int) {
	remaining := make([]int, n)
	for i, deps := range dependencies {
		remaining[i] = len(deps)
	}
	ordered := make([]int, 0, n)
	var cyclic []int
	done := make([]bool, n)
	for len(ordered) < n {
		progressed := false
		for i := 0; i < n; i++ {
			if done[i] || remaining[i] != 0 {
				continue
			}
			done[i] = true
			ordered = append(ordered, i)
			progressed = true
			for j, deps := range dependencies {
				if _, ok := deps[i]; ok && !done[j] {
					remaining[j]--
				}
			}
		}
		if !progressed {
			for i := 0; i < n; i++ {
				if !done[i] {
					cyclic = append(cyclic, i)
					done[i] = true
					ordered = append(ordered, i)
				}
			}
		}
	}
	return ordered, cyclic
}
