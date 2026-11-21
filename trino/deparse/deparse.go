package deparse

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/trino/analysis"
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
		schemaMeta := &SchemaMetadata{
			Name:              s.Name,
			Owner:             s.Owner,
			Functions:         make([]*FunctionMetadata, 0, len(s.Functions)),
			Tables:            make([]*TableMetadata, 0, len(s.Tables)),
			Views:             make([]*ViewMetadata, 0, len(s.Views)),
			MaterializedViews: make([]*ViewMetadata, 0, len(s.MaterializedViews)),
		}
		for _, fn := range s.Functions {
			if fn == nil || fn.SkipDump {
				continue
			}
			def, err := prepareTrinoFunctionDefinition(metadata.Name, s.Name, fn)
			if err != nil {
				return nil, err
			}
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
			cols, err := fromProtoColumns(s.Name, t.Name, t.Columns)
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
				Owner:   t.Owner,
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
			viewCols, err := fromProtoColumns(s.Name, v.Name, v.Columns)
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

func prepareTrinoFunctionDefinition(dbName, schemaName string, fn *metadatapb.FunctionMetadata) (string, error) {
	cfs, _, _, ok := parseTrinoCreateFunction(fn.Definition)
	if !ok || cfs.Spec.Declaration.Name == nil || len(cfs.Spec.Declaration.Name.Parts) == 0 {
		return fn.Definition, nil
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
			return "", mismatch()
		}
		if declName.Loc.Start >= 0 && declName.Loc.End <= len(fn.Definition) && declName.Loc.Start < declName.Loc.End {
			if prefix := qualifyTrinoObjectPrefix(dbName, schemaName); prefix != "" {
				def = fn.Definition[:declName.Loc.Start] + prefix + fn.Definition[declName.Loc.Start:]
			}
		}
	case 2:
		if !strings.EqualFold(parts[0].Value, schemaName) || !strings.EqualFold(parts[1].Value, fn.Name) {
			return "", mismatch()
		}
		if dbName != "" && declName.Loc.Start >= 0 && declName.Loc.End <= len(fn.Definition) && declName.Loc.Start < declName.Loc.End {
			def = fn.Definition[:declName.Loc.Start] + quoteIdentifier(dbName) + "." + fn.Definition[declName.Loc.Start:]
		}
	case 3:
		if (dbName != "" && !strings.EqualFold(parts[0].Value, dbName)) ||
			!strings.EqualFold(parts[1].Value, schemaName) ||
			!strings.EqualFold(parts[2].Value, fn.Name) {
			return "", mismatch()
		}
	default:
		return "", mismatch()
	}
	if err := validateTrinoFunctionSignature(schemaName, fn.Name, fn.Signature, cfs.Spec.Declaration.Parameters); err != nil {
		return "", err
	}
	return def, nil
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

func fromProtoColumns(schemaName, tableName string, cols []*metadatapb.ColumnMetadata) ([]*ColumnMetadata, error) {
	if len(cols) == 0 {
		return nil, nil
	}
	out := make([]*ColumnMetadata, 0, len(cols))
	for _, c := range cols {
		if c == nil {
			continue
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
	for _, fn := range orderedFunctions(metadata) {
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
	for _, view := range orderedViews(metadata) {
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

func orderedFunctions(metadata *DatabaseSchemaMetadata) []orderedFunction {
	var fns []orderedFunction
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
	if len(fns) == 0 {
		return fns
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

	cat := &trinoFuncCatalog{
		dbName: metadata.Name,
		fns:    fns,
		parsed: parsed,
		byKey:  byKey,
	}

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
		return fns
	}
	indices := topoOrder(len(fns), dependencies)
	ordered := make([]orderedFunction, len(indices))
	for i, idx := range indices {
		ordered[i] = fns[idx]
	}
	return ordered
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
		return "interval"
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

func collectTrinoASTFunctionCalls(root trinoparser.RoutineStatement, callerSchema string, params []trinoFuncParam, cat *trinoFuncCatalog) []functionCallRef {
	var calls []functionCallRef
	env := make(map[string]string, len(params))
	for _, p := range params {
		if p.name != "" {
			env[strings.ToLower(p.name)] = p.typeName
		}
	}

	var walkExpr func(expr trinoparser.Expr, curEnv map[string]string) string
	var walkQuery func(q *trinoparser.Query, curEnv map[string]string) string
	var walkQueryNode func(qn trinoparser.QueryNode, curEnv map[string]string) string
	var walkRelation func(rel trinoparser.Relation, curEnv map[string]string)

	walkRelation = func(rel trinoparser.Relation, curEnv map[string]string) {
		if rel == nil {
			return
		}
		switch r := rel.(type) {
		case *trinoparser.Join:
			walkRelation(r.Left, curEnv)
			walkRelation(r.Right, curEnv)
			walkExpr(r.On, curEnv)
		case *trinoparser.AliasedRelation:
			walkRelation(r.Inner, curEnv)
			if r.Sample != nil {
				walkExpr(r.Sample.Percentage, curEnv)
			}
		case *trinoparser.TableRelation:
			if r.Period != nil {
				walkExpr(r.Period.End, curEnv)
			}
		case *trinoparser.NearestRelation:
			walkRelation(r.From, curEnv)
			walkExpr(r.Where, curEnv)
			walkExpr(r.Match, curEnv)
		case *trinoparser.SubqueryRelation:
			walkQuery(r.Query, curEnv)
		case *trinoparser.LateralRelation:
			walkQuery(r.Query, curEnv)
		case *trinoparser.UnnestRelation:
			for _, ex := range r.Exprs {
				walkExpr(ex, curEnv)
			}
		case *trinoparser.ParenRelation:
			walkRelation(r.Inner, curEnv)
		case *trinoparser.TableFunctionRelation:
			if r.Call != nil {
				for _, arg := range r.Call.Args {
					walkExpr(arg.Expr, curEnv)
					if arg.Table != nil {
						walkQuery(arg.Table.Query, curEnv)
						for _, pb := range arg.Table.PartitionBy {
							walkExpr(pb, curEnv)
						}
						for _, ob := range arg.Table.OrderBy {
							walkExpr(ob.Expr, curEnv)
						}
					}
				}
			}
		}
	}

	walkQueryNode = func(qn trinoparser.QueryNode, curEnv map[string]string) string {
		if qn == nil {
			return ""
		}
		switch n := qn.(type) {
		case *trinoparser.QuerySpec:
			for _, rel := range n.From {
				walkRelation(rel, curEnv)
			}
			walkExpr(n.Where, curEnv)
			if n.GroupBy != nil {
				for _, el := range n.GroupBy.Elements {
					for _, ex := range el.Exprs {
						walkExpr(ex, curEnv)
					}
					for _, set := range el.Sets {
						for _, ex := range set {
							walkExpr(ex, curEnv)
						}
					}
				}
			}
			walkExpr(n.Having, curEnv)
			for _, w := range n.Windows {
				for _, pb := range w.Spec.PartitionBy {
					walkExpr(pb, curEnv)
				}
				for _, ob := range w.Spec.OrderBy {
					walkExpr(ob.Expr, curEnv)
				}
				if w.Spec.Frame != nil {
					walkExpr(w.Spec.Frame.Start.Value, curEnv)
					if w.Spec.Frame.End != nil {
						walkExpr(w.Spec.Frame.End.Value, curEnv)
					}
				}
			}
			var singleType string
			for idx, item := range n.Items {
				it := walkExpr(item.Expr, curEnv)
				if idx == 0 && len(n.Items) == 1 && item.Kind == trinoparser.SelectSingle {
					singleType = it
				}
			}
			return singleType
		case *trinoparser.SetOperation:
			lt := walkQueryNode(n.Left, curEnv)
			rt := walkQueryNode(n.Right, curEnv)
			return combineTrinoTypes(lt, rt)
		case *trinoparser.ParenQuery:
			return walkQuery(n.Inner, curEnv)
		case *trinoparser.ValuesQuery:
			for _, row := range n.Rows {
				walkExpr(row, curEnv)
			}
		}
		return ""
	}

	walkQuery = func(q *trinoparser.Query, curEnv map[string]string) string {
		if q == nil {
			return ""
		}
		if q.With != nil {
			for _, cte := range q.With.CTEs {
				walkQuery(cte.Query, curEnv)
			}
		}
		ret := walkQueryNode(q.Body, curEnv)
		for _, ob := range q.OrderBy {
			walkExpr(ob.Expr, curEnv)
		}
		if q.Offset != nil {
			walkExpr(q.Offset.Count, curEnv)
		}
		if q.Limit != nil {
			walkExpr(q.Limit.Count, curEnv)
		}
		if q.Fetch != nil {
			walkExpr(q.Fetch.Count, curEnv)
		}
		return ret
	}
	walkJSONPath := func(path *trinoparser.JSONPathInvocation, curEnv map[string]string) {
		if path == nil {
			return
		}
		if path.Input != nil {
			walkExpr(path.Input.Expr, curEnv)
		}
		for _, p := range path.Passing {
			if p.Value != nil {
				walkExpr(p.Value.Expr, curEnv)
			}
		}
	}
	walkExpr = func(expr trinoparser.Expr, curEnv map[string]string) string {
		if expr == nil {
			return ""
		}
		switch e := expr.(type) {
		case *trinoparser.NamedArgument:
			return walkExpr(e.Value, curEnv)
		case *trinoparser.ColumnRef:
			if e.Name != nil {
				return curEnv[strings.ToLower(e.Name.Value)]
			}
		case *trinoparser.Dereference:
			baseType := walkExpr(e.Base, curEnv)
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
			baseType := walkExpr(e.Value, curEnv)
			walkExpr(e.Index, curEnv)
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
			return "interval"
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
			walkExpr(e.Expr, curEnv)
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
			return walkExpr(e.Operand, curEnv)
		case *trinoparser.ParenExpr:
			return walkExpr(e.Expr, curEnv)
		case *trinoparser.NotExpr:
			walkExpr(e.Operand, curEnv)
			return "boolean"
		case *trinoparser.LogicalExpr:
			walkExpr(e.Left, curEnv)
			walkExpr(e.Right, curEnv)
			return "boolean"
		case *trinoparser.BinaryExpr:
			lt := walkExpr(e.Left, curEnv)
			rt := walkExpr(e.Right, curEnv)
			switch e.Op {
			case "||":
				return "varchar"
			case "+", "-", "*", "/", "%":
				return commonTrinoNumericType(lt, rt)
			}
		case *trinoparser.CaseExpr:
			walkExpr(e.Operand, curEnv)
			res := ""
			for _, w := range e.Whens {
				walkExpr(w.Cond, curEnv)
				res = combineTrinoTypes(res, walkExpr(w.Result, curEnv))
			}
			res = combineTrinoTypes(res, walkExpr(e.Else, curEnv))
			return res
		case *trinoparser.AtTimeZoneExpr:
			t := walkExpr(e.Value, curEnv)
			walkExpr(e.Zone, curEnv)
			return zonedTrinoDateTimeType(t)
		case *trinoparser.AtLocalExpr:
			return zonedTrinoDateTimeType(walkExpr(e.Value, curEnv))
		case *trinoparser.ArrayConstructor:
			elemType := ""
			for _, el := range e.Elements {
				elemType = combineTrinoTypes(elemType, walkExpr(el, curEnv))
			}
			if elemType != "" {
				return "array(" + elemType + ")"
			}
			return "array"
		case *trinoparser.RowConstructor:
			fieldTypes := make([]string, len(e.Elements))
			allKnown := len(e.Elements) > 0
			for j, el := range e.Elements {
				ft := walkExpr(el, curEnv)
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
			walkExpr(e.Source, curEnv)
			walkExpr(e.Replacing, curEnv)
			walkExpr(e.From, curEnv)
			walkExpr(e.For, curEnv)
			return "varchar"
		case *trinoparser.ExtractExpr:
			walkExpr(e.Source, curEnv)
			return "bigint"
		case *trinoparser.SubstringExpr:
			walkExpr(e.Source, curEnv)
			walkExpr(e.From, curEnv)
			walkExpr(e.For, curEnv)
			return "varchar"
		case *trinoparser.ComparisonExpr:
			walkExpr(e.Left, curEnv)
			walkExpr(e.Right, curEnv)
			return "boolean"
		case *trinoparser.QuantifiedComparisonExpr:
			walkExpr(e.Left, curEnv)
			walkExpr(e.Subquery, curEnv)
			return "boolean"
		case *trinoparser.BetweenExpr:
			walkExpr(e.Value, curEnv)
			walkExpr(e.Lower, curEnv)
			walkExpr(e.Upper, curEnv)
			return "boolean"
		case *trinoparser.InListExpr:
			walkExpr(e.Value, curEnv)
			for _, item := range e.List {
				walkExpr(item, curEnv)
			}
			return "boolean"
		case *trinoparser.InSubqueryExpr:
			walkExpr(e.Value, curEnv)
			walkExpr(e.Subquery, curEnv)
			return "boolean"
		case *trinoparser.LikeExpr:
			walkExpr(e.Value, curEnv)
			walkExpr(e.Pattern, curEnv)
			walkExpr(e.Escape, curEnv)
			return "boolean"
		case *trinoparser.IsNullExpr:
			walkExpr(e.Value, curEnv)
			return "boolean"
		case *trinoparser.IsDistinctFromExpr:
			walkExpr(e.Value, curEnv)
			walkExpr(e.Right, curEnv)
			return "boolean"
		case *trinoparser.IsBooleanExpr:
			walkExpr(e.Value, curEnv)
			return "boolean"
		case *trinoparser.MatchExpr:
			walkExpr(e.Value, curEnv)
			walkExpr(e.Subquery, curEnv)
			return "boolean"
		case *trinoparser.UniqueExpr:
			walkExpr(e.Subquery, curEnv)
			return "boolean"
		case *trinoparser.SubqueryExpr:
			var subType string
			if e.Query != nil {
				subType = walkQuery(e.Query, curEnv)
			} else if e.RawText != "" {
				calls = append(calls, scanFunctionCallRefs(e.RawText)...)
			}
			if e.Kind == trinoparser.SubqueryExists || e.Kind == trinoparser.SubqueryUnique {
				return "boolean"
			}
			if e.Kind == trinoparser.SubqueryScalar {
				return subType
			}
		case *trinoparser.TrimExpr:
			walkExpr(e.Char, curEnv)
			walkExpr(e.Source, curEnv)
			return "varchar"
		case *trinoparser.NormalizeExpr:
			walkExpr(e.Source, curEnv)
			return "varchar"
		case *trinoparser.PositionExpr:
			walkExpr(e.Needle, curEnv)
			walkExpr(e.Haystack, curEnv)
			return "bigint"
		case *trinoparser.ListaggExpr:
			walkExpr(e.Arg, curEnv)
			walkExpr(e.Separator, curEnv)
			for _, ob := range e.WithinGroupBy {
				walkExpr(ob.Expr, curEnv)
			}
			walkExpr(e.Filter, curEnv)
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
			walkExpr(e.Body, lambdaEnv)
		case *trinoparser.JSONExistsExpr:
			walkJSONPath(e.Path, curEnv)
			return "boolean"
		case *trinoparser.JSONValueFunc:
			walkJSONPath(e.Path, curEnv)
			if e.OnEmpty != nil {
				walkExpr(e.OnEmpty.Default, curEnv)
			}
			if e.OnError != nil {
				walkExpr(e.OnError.Default, curEnv)
			}
			return normalizeTrinoDataType(e.Returning)
		case *trinoparser.JSONQueryFunc:
			walkJSONPath(e.Path, curEnv)
			return normalizeTrinoDataType(e.Returning)
		case *trinoparser.JSONObjectExpr:
			for _, m := range e.Members {
				walkExpr(m.Key, curEnv)
				if m.Value != nil {
					walkExpr(m.Value.Expr, curEnv)
				}
			}
			return normalizeTrinoDataType(e.Returning)
		case *trinoparser.JSONArrayExpr:
			for _, el := range e.Elements {
				if el != nil {
					walkExpr(el.Expr, curEnv)
				}
			}
			return normalizeTrinoDataType(e.Returning)
		case *trinoparser.MethodCall:
			walkExpr(e.Value, curEnv)
			for _, a := range e.Args {
				walkExpr(a, curEnv)
			}
		case *trinoparser.StaticMethodCall:
			for _, a := range e.Args {
				walkExpr(a, curEnv)
			}
		case *trinoparser.FuncCall:
			argTypes := make([]string, len(e.Args))
			var argNames []string
			for j, a := range e.Args {
				if na, ok := a.(*trinoparser.NamedArgument); ok && na != nil {
					argTypes[j] = walkExpr(na.Value, curEnv)
					if na.Name != nil && na.Name.Value != "" {
						if argNames == nil {
							argNames = make([]string, len(e.Args))
						}
						argNames[j] = na.Name.Value
					}
				} else {
					argTypes[j] = walkExpr(a, curEnv)
				}
			}
			for _, ob := range e.OrderBy {
				walkExpr(ob.Expr, curEnv)
			}
			walkExpr(e.Filter, curEnv)
			if e.Over != nil {
				for _, pb := range e.Over.PartitionBy {
					walkExpr(pb, curEnv)
				}
				for _, ob := range e.Over.OrderBy {
					walkExpr(ob.Expr, curEnv)
				}
				if e.Over.Frame != nil {
					walkExpr(e.Over.Frame.Start.Value, curEnv)
					if e.Over.Frame.End != nil {
						walkExpr(e.Over.Frame.End.Value, curEnv)
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
					// 4+ parts cannot reference a local catalog function.
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
				calls = append(calls, ref)
				if matches := cat.resolveCall(callerSchema, ref); len(matches) > 0 {
					ret := cat.parsed[matches[0]].returnType
					same := ret != ""
					for _, m := range matches[1:] {
						if cat.parsed[m].returnType != ret {
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

	var walkStmt func(stmt trinoparser.RoutineStatement, curEnv map[string]string)
	walkStmt = func(stmt trinoparser.RoutineStatement, curEnv map[string]string) {
		if stmt == nil {
			return
		}
		switch s := stmt.(type) {
		case *trinoparser.ReturnStatement:
			walkExpr(s.Value, curEnv)
		case *trinoparser.AssignmentStatement:
			walkExpr(s.Value, curEnv)
		case *trinoparser.IfStatement:
			walkExpr(s.Condition, curEnv)
			for _, st := range s.Then {
				walkStmt(st, curEnv)
			}
			for _, elif := range s.ElseIfs {
				walkExpr(elif.Condition, curEnv)
				for _, st := range elif.Body {
					walkStmt(st, curEnv)
				}
			}
			for _, st := range s.Else {
				walkStmt(st, curEnv)
			}
		case *trinoparser.CaseStatement:
			walkExpr(s.Operand, curEnv)
			for _, w := range s.Whens {
				walkExpr(w.Condition, curEnv)
				for _, st := range w.Body {
					walkStmt(st, curEnv)
				}
			}
			for _, st := range s.Else {
				walkStmt(st, curEnv)
			}
		case *trinoparser.LoopStatement:
			for _, st := range s.Body {
				walkStmt(st, curEnv)
			}
		case *trinoparser.WhileStatement:
			walkExpr(s.Condition, curEnv)
			for _, st := range s.Body {
				walkStmt(st, curEnv)
			}
		case *trinoparser.RepeatStatement:
			for _, st := range s.Body {
				walkStmt(st, curEnv)
			}
			walkExpr(s.Condition, curEnv)
		case *trinoparser.CompoundStatement:
			childEnv := make(map[string]string, len(curEnv)+len(s.Declarations))
			for k, v := range curEnv {
				childEnv[k] = v
			}
			for _, decl := range s.Declarations {
				walkExpr(decl.Default, childEnv)
				dt := normalizeTrinoDataType(decl.Type)
				for _, n := range decl.Names {
					if n != nil && dt != "" {
						childEnv[strings.ToLower(n.Value)] = dt
					}
				}
			}
			for _, st := range s.Body {
				walkStmt(st, childEnv)
			}
		}
	}

	walkStmt(root, env)
	return calls
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
func orderedViews(metadata *DatabaseSchemaMetadata) []orderedView {
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
		return views
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
	indices := topoOrder(len(views), dependencies)
	ordered := make([]orderedView, len(indices))
	for i, idx := range indices {
		ordered[i] = views[idx]
	}
	return ordered
}

// objectKey identifies a schema-scoped object by its owning schema and its name.
func objectKey(schema, name string) string {
	return schema + "\x00" + name
}

// topoOrder returns item indices ordered so each item follows its dependencies,
// keeping input order among independent items and breaking cycles in input order.
func topoOrder(n int, dependencies []map[int]struct{}) []int {
	remaining := make([]int, n)
	for i, deps := range dependencies {
		remaining[i] = len(deps)
	}
	ordered := make([]int, 0, n)
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
					done[i] = true
					ordered = append(ordered, i)
				}
			}
		}
	}
	return ordered
}
