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
	if metadata.CharacterSet != "" || metadata.Collation != "" || metadata.Owner != "" || metadata.SearchPath != "" || metadata.ServiceName != "" || metadata.DefaultSchema != "" || metadata.Datashare || len(metadata.Extensions) > 0 || len(metadata.EventTriggers) > 0 || len(metadata.PublicSynonyms) > 0 || len(metadata.LinkedDatabases) > 0 || len(metadata.ForeignServers) > 0 || len(metadata.ForeignDataWrappers) > 0 {
		return nil, fmt.Errorf("database %s contains database attributes or objects unsupported by Trino SDL renderer", quoteIdentifier(metadata.Name))
	}
	out := &DatabaseSchemaMetadata{
		Name:    metadata.Name,
		Schemas: make([]*SchemaMetadata, 0, len(metadata.Schemas)),
	}
	seenSchemas := make(map[string]struct{}, len(metadata.Schemas))
	for _, s := range metadata.Schemas {
		if s == nil || s.SkipDump {
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
			if fn == nil || fn.SkipDump {
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
		if len(s.Streams) > 0 || len(s.Tasks) > 0 || len(s.Packages) > 0 || len(s.Events) > 0 || len(s.EnumTypes) > 0 || len(s.CompositeTypes) > 0 || len(s.Synonyms) > 0 || len(s.Collations) > 0 || len(s.RangeTypes) > 0 {
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
			seenCols := make(map[string]struct{}, len(cols))
			for _, c := range cols {
				colKey := strings.ToLower(c.Name)
				if _, exists := seenCols[colKey]; exists {
					return nil, fmt.Errorf("table %s.%s has duplicate column %s", s.Name, t.Name, quoteIdentifier(c.Name))
				}
				seenCols[colKey] = struct{}{}
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
	funcCat := buildTrinoFuncCatalog(out)
	for _, schemaMeta := range out.Schemas {
		for _, v := range schemaMeta.Views {
			v.Definition = qualifyTrinoViewBody(metadata.Name, schemaMeta.Name, v.Definition, funcCat)
		}
		for _, mv := range schemaMeta.MaterializedViews {
			mv.Definition = qualifyTrinoViewBody(metadata.Name, schemaMeta.Name, mv.Definition, funcCat)
		}
	}
	return out, nil
}

func prepareTrinoFunctionDefinition(dbName, schemaName string, fn *metadatapb.FunctionMetadata) (string, string, error) {
	if strings.TrimSpace(fn.Definition) == "" {
		return "", "", fmt.Errorf("function %s.%s has no definition to restore", schemaName, fn.Name)
	}
	if !isTrinoCreateFunctionHeader(fn.Definition) {
		return "", "", fmt.Errorf("function %s.%s definition must be a CREATE FUNCTION statement", schemaName, fn.Name)
	}
	fallbackSigKey := strings.ToLower(fn.Name) + "(" + strings.ToLower(strings.TrimSpace(fn.Signature)) + ")"
	cfs, _, _, ok := parseTrinoCreateFunction(fn.Definition)
	if !ok || cfs.Spec.Declaration.Name == nil || len(cfs.Spec.Declaration.Name.Parts) == 0 {
		return fn.Definition, fallbackSigKey, nil
	}
	declName := cfs.Spec.Declaration.Name
	parts := declName.Parts
	mismatch := func() error {
		return fmt.Errorf("function %s.%s definition name %q does not match metadata identity", schemaName, fn.Name, declName.String())
	}
	def := fn.Definition
	switch len(parts) {
	case 1:
		if !strings.EqualFold(parts[0].Value, fn.Name) {
			return "", "", mismatch()
		}
		if declName.Loc.Start >= 0 && declName.Loc.End <= len(fn.Definition) && declName.Loc.Start < declName.Loc.End {
			if prefix := qualifyTrinoObjectPrefix(dbName, schemaName); prefix != "" {
				def = fn.Definition[:declName.Loc.Start] + prefix + fn.Definition[declName.Loc.Start:]
			}
		}
	case 2:
		if !strings.EqualFold(parts[0].Value, schemaName) || !strings.EqualFold(parts[1].Value, fn.Name) {
			return "", "", mismatch()
		}
		if dbName != "" && declName.Loc.Start >= 0 && declName.Loc.End <= len(fn.Definition) && declName.Loc.Start < declName.Loc.End {
			def = fn.Definition[:declName.Loc.Start] + quoteIdentifier(dbName) + "." + fn.Definition[declName.Loc.Start:]
		}
	case 3:
		if (dbName != "" && !strings.EqualFold(parts[0].Value, dbName)) ||
			!strings.EqualFold(parts[1].Value, schemaName) ||
			!strings.EqualFold(parts[2].Value, fn.Name) {
			return "", "", mismatch()
		}
	default:
		return "", "", mismatch()
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

func qualifyTrinoViewBody(dbName, schemaName, definition string, cat *trinoFuncCatalog) string {
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
	calls, unqualTables := collectTrinoQueryASTRefs(query, inlineFns, schemaName, cat)
	type prefixEdit struct {
		offset int
		prefix string
	}
	var edits []prefixEdit
	seen := make(map[int]struct{})
	addEdit := func(qnLoc ast.Loc, prefix string) {
		if prefix == "" || qnLoc.Start < 0 || qnLoc.End <= qnLoc.Start || qnLoc.End > len(definition) {
			return
		}
		if _, dup := seen[qnLoc.Start]; !dup {
			seen[qnLoc.Start] = struct{}{}
			edits = append(edits, prefixEdit{offset: qnLoc.Start, prefix: prefix})
		}
	}
	relPrefix := qualifyTrinoObjectPrefix(dbName, schemaName)
	for _, qn := range unqualTables {
		addEdit(qn.Loc, relPrefix)
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
		return definition
	}
	slices.SortFunc(edits, func(a, b prefixEdit) int { return b.offset - a.offset })
	out := definition
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
	for _, c := range cols {
		if c == nil {
			continue
		}
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

	var buf strings.Builder
	for _, schema := range metadata.Schemas {
		if schema == nil {
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
		if err := writeCreateFunction(&buf, fn.schema, fn.fn); err != nil {
			return "", err
		}
		buf.WriteString("\n\n")
	}
	for _, schema := range metadata.Schemas {
		if schema == nil {
			continue
		}
		for _, table := range schema.Tables {
			if table == nil {
				continue
			}
			writeCreateTable(&buf, schema.Name, table)
			buf.WriteString("\n\n")
		}
	}
	for _, view := range views {
		if err := writeCreateView(&buf, view.schema, view.view, view.materialized); err != nil {
			return "", err
		}
		buf.WriteString("\n\n")
	}
	for _, schema := range metadata.Schemas {
		if schema == nil {
			continue
		}
		for _, table := range schema.Tables {
			if table == nil || table.Owner == "" {
				continue
			}
			fmt.Fprintf(&buf, "ALTER TABLE %s.%s SET AUTHORIZATION %s;\n\n", quoteIdentifier(schema.Name), quoteIdentifier(table.Name), quoteIdentifier(table.Owner))
		}
		if schema.Owner != "" {
			fmt.Fprintf(&buf, "ALTER SCHEMA %s SET AUTHORIZATION %s;\n\n", quoteIdentifier(schema.Name), quoteIdentifier(schema.Owner))
		}
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

func formatTrinoDefault(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}
	expr, errs := trinoparser.ParseExpression(trimmed)
	if len(errs) == 0 && expr != nil && isTrinoDefaultAST(expr) {
		return trimmed
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

func writeCreateFunction(buf *strings.Builder, schemaName string, fn *FunctionMetadata) error {
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
	schema string
	fn     *FunctionMetadata
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
	dbName string
	fns    []orderedFunction
	parsed []parsedTrinoFn
	byKey  map[string][]int
}

func buildTrinoFuncCatalog(metadata *DatabaseSchemaMetadata) *trinoFuncCatalog {
	var fns []orderedFunction
	if metadata != nil {
		for _, schema := range metadata.Schemas {
			if schema == nil {
				continue
			}
			for _, fn := range schema.Functions {
				if fn == nil {
					continue
				}
				fns = append(fns, orderedFunction{schema: schema.Name, fn: fn})
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
		dbName: dbName,
		fns:    fns,
		parsed: parsed,
		byKey:  byKey,
	}
}

func orderedFunctions(metadata *DatabaseSchemaMetadata) ([]orderedFunction, error) {
	cat := buildTrinoFuncCatalog(metadata)
	fns := cat.fns
	if len(fns) == 0 {
		return fns, nil
	}
	parsed := cat.parsed

	type prefixEdit struct {
		offset int
		prefix string
	}
	dependencies := make([]map[int]struct{}, len(fns))
	for i, f := range fns {
		var calls []functionCallRef
		switch {
		case parsed[i].bodyStmt != nil:
			calls = collectTrinoASTFunctionCalls(parsed[i].bodyStmt, f.schema, parsed[i].params, cat)
		case !parsed[i].parsed:
			calls = scanFunctionCallRefs(f.fn.Definition)
		}
		var edits []prefixEdit
		seenEdit := make(map[int]struct{})
		for _, call := range calls {
			matches := cat.resolveCall(f.schema, call)
			if len(matches) == 0 {
				continue
			}
			for _, depIdx := range matches {
				if depIdx == i {
					continue
				}
				if dependencies[i] == nil {
					dependencies[i] = map[int]struct{}{}
				}
				dependencies[i][depIdx] = struct{}{}
			}
			if parsed[i].parsed && call.catalog == "" && call.nameEnd > call.nameStart && call.nameEnd <= len(f.fn.Definition) {
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
			fns[i].fn = &FunctionMetadata{
				Name:       f.fn.Name,
				Definition: def,
				Comment:    f.fn.Comment,
			}
		}
	}

	if len(fns) < 2 {
		return fns, nil
	}
	indices, cyclic := topoOrder(len(fns), dependencies)
	if len(cyclic) > 0 {
		return nil, fmt.Errorf("cyclic function dependency involving %s.%s cannot be restored because Trino requires UDFs to be declared before use", fns[cyclic[0]].schema, fns[cyclic[0]].fn.Name)
	}
	ordered := make([]orderedFunction, len(indices))
	for i, idx := range indices {
		ordered[i] = fns[idx]
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

func hasMatchingTrinoBuiltin(call functionCallRef) bool {
	for _, n := range call.argNames {
		if n != "" {
			return false
		}
	}
	name := strings.ToLower(call.name)
	sigs, ok := trinoBuiltinFunctions[name]
	if !ok {
		return false
	}
	// SHOW FUNCTIONS lists variadic concat/greatest/least with a single element parameter.
	switch name {
	case "greatest", "least":
		if call.argCount >= 1 {
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

func trinoBuiltinTypeMatch(argType, paramType string) bool {
	if paramType == "*" {
		return true
	}
	if argType == "" {
		return true
	}
	if paramType == "function" {
		return false
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

type trinoASTWalker struct {
	callerSchema string
	cat          *trinoFuncCatalog
	calls        []functionCallRef
	unqualTables []*ast.QualifiedName
	inlineFuncs  map[string][]trinoInlineFunc
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

func collectTrinoQueryASTRefs(q *trinoparser.Query, extraInline []trinoparser.FunctionSpecification, callerSchema string, cat *trinoFuncCatalog) ([]functionCallRef, []*ast.QualifiedName) {
	w := &trinoASTWalker{callerSchema: callerSchema, cat: cat}
	w.walkQueryWithInline(q, extraInline, nil, nil)
	return w.calls, w.unqualTables
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

func (w *trinoASTWalker) walkRelation(rel trinoparser.Relation, curEnv map[string]string, ctes map[string]struct{}) {
	if rel == nil {
		return
	}
	switch r := rel.(type) {
	case *trinoparser.Join:
		w.walkRelation(r.Left, curEnv, ctes)
		w.walkRelation(r.Right, curEnv, ctes)
		w.walkExpr(r.On, curEnv, ctes)
	case *trinoparser.AliasedRelation:
		w.walkRelation(r.Inner, curEnv, ctes)
		if r.Sample != nil {
			w.walkExpr(r.Sample.Percentage, curEnv, ctes)
		}
	case *trinoparser.TableRelation:
		w.recordUnqualTable(r.Name, ctes)
		if r.Period != nil {
			w.walkExpr(r.Period.End, curEnv, ctes)
		}
	case *trinoparser.NearestRelation:
		w.walkRelation(r.From, curEnv, ctes)
		w.walkExpr(r.Where, curEnv, ctes)
		w.walkExpr(r.Match, curEnv, ctes)
	case *trinoparser.SubqueryRelation:
		w.walkQuery(r.Query, curEnv, ctes)
	case *trinoparser.LateralRelation:
		w.walkQuery(r.Query, curEnv, ctes)
	case *trinoparser.UnnestRelation:
		for _, ex := range r.Exprs {
			w.walkExpr(ex, curEnv, ctes)
		}
	case *trinoparser.ParenRelation:
		w.walkRelation(r.Inner, curEnv, ctes)
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

func (w *trinoASTWalker) walkQueryNode(qn trinoparser.QueryNode, curEnv map[string]string, ctes map[string]struct{}) string {
	if qn == nil {
		return ""
	}
	switch n := qn.(type) {
	case *trinoparser.QuerySpec:
		for _, rel := range n.From {
			w.walkRelation(rel, curEnv, ctes)
		}
		w.walkExpr(n.Where, curEnv, ctes)
		if n.GroupBy != nil {
			for _, el := range n.GroupBy.Elements {
				for _, ex := range el.Exprs {
					w.walkExpr(ex, curEnv, ctes)
				}
				for _, set := range el.Sets {
					for _, ex := range set {
						w.walkExpr(ex, curEnv, ctes)
					}
				}
			}
		}
		w.walkExpr(n.Having, curEnv, ctes)
		for _, win := range n.Windows {
			for _, pb := range win.Spec.PartitionBy {
				w.walkExpr(pb, curEnv, ctes)
			}
			for _, ob := range win.Spec.OrderBy {
				w.walkExpr(ob.Expr, curEnv, ctes)
			}
			if win.Spec.Frame != nil {
				w.walkExpr(win.Spec.Frame.Start.Value, curEnv, ctes)
				if win.Spec.Frame.End != nil {
					w.walkExpr(win.Spec.Frame.End.Value, curEnv, ctes)
				}
			}
		}
		var singleType string
		for idx, item := range n.Items {
			it := w.walkExpr(item.Expr, curEnv, ctes)
			if idx == 0 && len(n.Items) == 1 && item.Kind == trinoparser.SelectSingle {
				singleType = it
			}
		}
		return singleType
	case *trinoparser.TableQuery:
		w.recordUnqualTable(n.Name, ctes)
	case *trinoparser.SetOperation:
		lt := w.walkQueryNode(n.Left, curEnv, ctes)
		rt := w.walkQueryNode(n.Right, curEnv, ctes)
		return combineTrinoTypes(lt, rt)
	case *trinoparser.ParenQuery:
		return w.walkQuery(n.Inner, curEnv, ctes)
	case *trinoparser.ValuesQuery:
		for _, row := range n.Rows {
			w.walkExpr(row, curEnv, ctes)
		}
	}
	return ""
}

func (w *trinoASTWalker) walkQuery(q *trinoparser.Query, curEnv map[string]string, ctes map[string]struct{}) string {
	return w.walkQueryWithInline(q, nil, curEnv, ctes)
}

func (w *trinoASTWalker) walkQueryWithInline(q *trinoparser.Query, extraInline []trinoparser.FunctionSpecification, curEnv map[string]string, ctes map[string]struct{}) string {
	if q == nil {
		return ""
	}
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
	if q.With != nil {
		for _, cte := range q.With.CTEs {
			if q.With.Recursive && cte.Name != nil && cte.Name.Value != "" {
				scope[strings.ToLower(cte.Name.Value)] = struct{}{}
			}
			w.walkQuery(cte.Query, curEnv, scope)
			if cte.Name != nil && cte.Name.Value != "" {
				scope[strings.ToLower(cte.Name.Value)] = struct{}{}
			}
		}
	}
	ret := w.walkQueryNode(q.Body, curEnv, scope)
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
	return ret
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
			if base, inner, ok := splitParameterizedType(baseType); ok && base == "row" {
				target := strings.ToLower(e.FieldName.Value)
				for _, part := range splitTopLevelComma(inner) {
					if name, typ := splitRowField(part); name == target {
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
		w.walkExpr(e.Value, curEnv, ctes)
		w.walkExpr(e.Subquery, curEnv, ctes)
		return "boolean"
	case *trinoparser.UniqueExpr:
		w.walkExpr(e.Subquery, curEnv, ctes)
		return "boolean"
	case *trinoparser.SubqueryExpr:
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
			for _, pb := range e.Over.PartitionBy {
				w.walkExpr(pb, curEnv, ctes)
			}
			for _, ob := range e.Over.OrderBy {
				w.walkExpr(ob.Expr, curEnv, ctes)
			}
			if e.Over.Frame != nil {
				w.walkExpr(e.Over.Frame.Start.Value, curEnv, ctes)
				if e.Over.Frame.End != nil {
					w.walkExpr(e.Over.Frame.End.Value, curEnv, ctes)
				}
			}
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
	fmt.Fprintf(buf, "CREATE TABLE IF NOT EXISTS %s.%s (\n", quoteIdentifier(schemaName), quoteIdentifier(table.Name))

	columnDefs := make([]string, 0, len(table.Columns))
	for _, column := range table.Columns {
		if column == nil {
			continue
		}
		colDef := fmt.Sprintf("    %s %s", quoteIdentifier(column.Name), column.Type)
		if column.Default != "" {
			colDef += " DEFAULT " + formatTrinoDefault(column.Default)
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
		return
	}
	buf.WriteString("\n);")
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
		if schema == nil {
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
	if len(views) < 2 {
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
		if !ok || dependency == i {
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
