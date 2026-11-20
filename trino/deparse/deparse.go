package deparse

import (
	"fmt"
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
	out := &DatabaseSchemaMetadata{
		Name:    metadata.Name,
		Schemas: make([]*SchemaMetadata, 0, len(metadata.Schemas)),
	}
	for _, s := range metadata.Schemas {
		if s == nil || s.SkipDump {
			continue
		}
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
			schemaMeta.Functions = append(schemaMeta.Functions, &FunctionMetadata{
				Name:       fn.Name,
				Definition: fn.Definition,
				Comment:    fn.Comment,
			})
		}
		for _, t := range s.Tables {
			if t == nil || t.SkipDump {
				continue
			}
			cols := fromProtoColumns(t.Columns)
			// GetDatabaseDefinition keeps emitting empty CREATE TABLE () for legacy
			// direct-struct oracle/golden parity, while DatabaseDefinitionFromProto
			// rejects columnless tables because Trino cannot restore them.
			if len(cols) == 0 {
				return nil, fmt.Errorf("table %s.%s has no columns", s.Name, t.Name)
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
			schemaMeta.Views = append(schemaMeta.Views, &ViewMetadata{
				Name:         v.Name,
				Definition:   v.Definition,
				Comment:      v.Comment,
				Columns:      fromProtoColumns(v.Columns),
				Dependencies: fromProtoDependencies(v.DependencyColumns),
			})
		}
		for _, mv := range s.MaterializedViews {
			if mv == nil || mv.SkipDump {
				continue
			}
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

func fromProtoColumns(cols []*metadatapb.ColumnMetadata) []*ColumnMetadata {
	if len(cols) == 0 {
		return nil
	}
	out := make([]*ColumnMetadata, 0, len(cols))
	for _, c := range cols {
		if c == nil {
			continue
		}
		out = append(out, &ColumnMetadata{
			Name:     c.Name,
			Type:     c.Type,
			Nullable: c.Nullable,
			Default:  c.Default,
			Comment:  c.Comment,
		})
	}
	return out
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
	file, err := trinoparser.Parse(definition)
	if err != nil || file == nil || len(file.Stmts) != 1 {
		return "", fmt.Errorf("function %s.%s comment cannot be placed in unparseable definition", schemaName, fnName)
	}
	cfs, ok := file.Stmts[0].(*trinoparser.CreateFunctionStmt)
	if !ok || cfs.Spec.Body == nil {
		return "", fmt.Errorf("function %s.%s comment cannot be placed in definition", schemaName, fnName)
	}
	clause := fmt.Sprintf("COMMENT '%s'", escapeStringLiteral(comment))
	for _, ch := range cfs.Spec.Characteristics {
		if ch.Kind == trinoparser.CharacteristicComment && ch.Loc.Start >= 0 && ch.Loc.End <= len(definition) && ch.Loc.Start < ch.Loc.End {
			return definition[:ch.Loc.Start] + clause + definition[ch.Loc.End:], nil
		}
	}
	insertAt := cfs.Spec.Body.Span().Start
	if insertAt < 0 || insertAt > len(definition) {
		return "", fmt.Errorf("function %s.%s has invalid body span", schemaName, fnName)
	}
	prefix := strings.TrimRight(definition[:insertAt], " \t\r\n")
	return prefix + " " + clause + " " + definition[insertAt:], nil
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
	if len(fns) < 2 {
		return fns
	}

	parsed := make([]parsedTrinoFn, len(fns))
	byKey := make(map[string][]int, len(fns))
	for i, f := range fns {
		key := objectKey(strings.ToLower(f.schema), strings.ToLower(f.fn.Name))
		byKey[key] = append(byKey[key], i)

		file, err := trinoparser.Parse(f.fn.Definition)
		if err == nil && file != nil && len(file.Stmts) == 1 {
			if cfs, ok := file.Stmts[0].(*trinoparser.CreateFunctionStmt); ok && cfs.Spec.Body != nil {
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
				parsed[i] = parsedTrinoFn{
					params:     params,
					returnType: normalizeTrinoDataType(cfs.Spec.ReturnType),
					bodyStmt:   cfs.Spec.Body,
					parsed:     true,
				}
			}
		}
	}

	cat := &trinoFuncCatalog{
		fns:    fns,
		parsed: parsed,
		byKey:  byKey,
	}

	dependencies := make([]map[int]struct{}, len(fns))
	for i, f := range fns {
		var calls []functionCallRef
		if parsed[i].parsed {
			calls = collectTrinoASTFunctionCalls(parsed[i].bodyStmt, f.schema, parsed[i].params, cat)
		} else {
			calls = scanFunctionCallRefs(f.fn.Definition)
		}
		for _, call := range calls {
			for _, depIdx := range cat.resolveCall(f.schema, call) {
				if depIdx == i {
					continue
				}
				if dependencies[i] == nil {
					dependencies[i] = map[int]struct{}{}
				}
				dependencies[i][depIdx] = struct{}{}
			}
		}
	}

	indices := topoOrder(len(fns), dependencies)
	ordered := make([]orderedFunction, len(indices))
	for i, idx := range indices {
		ordered[i] = fns[idx]
	}
	return ordered
}

type functionCallRef struct {
	schema   string
	name     string
	argTypes []string
	argCount int
}

func (c *trinoFuncCatalog) resolveCall(callerSchema string, call functionCallRef) []int {
	if c == nil || len(c.fns) == 0 {
		return nil
	}
	targetSchema := call.schema
	if targetSchema == "" {
		targetSchema = callerSchema
	}
	candidates := c.byKey[objectKey(strings.ToLower(targetSchema), strings.ToLower(call.name))]
	if len(candidates) <= 1 {
		if len(candidates) == 1 && len(call.argTypes) == 0 && call.argCount > 0 {
			// Single candidate matches name; check arity only if parsed.
			idx := candidates[0]
			if c.parsed[idx].parsed && len(c.parsed[idx].params) != call.argCount {
				return nil
			}
		}
		return candidates
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
	if len(arityMatches) == 1 || len(call.argTypes) == 0 {
		return arityMatches
	}
	bestScore := -1
	var best []int
	for _, idx := range arityMatches {
		if !c.parsed[idx].parsed {
			continue
		}
		score := 0
		compatible := true
		for i, argType := range call.argTypes {
			if argType == "" {
				continue
			}
			paramType := c.parsed[idx].params[i].typeName
			if paramType == "" {
				continue
			}
			if paramType == argType {
				score += 2
			} else if isTrinoNumericType(argType) && isTrinoNumericType(paramType) && trinoNumericRank(paramType) >= trinoNumericRank(argType) {
				score++
			} else {
				compatible = false
				break
			}
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
	if len(best) > 0 {
		return best
	}
	return arityMatches
}

func normalizeTrinoDataType(dt *trinoparser.DataType) string {
	if dt == nil {
		return ""
	}
	switch dt.Kind {
	case trinoparser.TypeDateTime:
		return strings.ToLower(dt.Name)
	case trinoparser.TypeInterval:
		return "interval"
	case trinoparser.TypeRow:
		return "row"
	case trinoparser.TypeArray:
		return "array"
	case trinoparser.TypeMap:
		return "map"
	default:
		return normalizeTrinoTypeName(dt.Name)
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
	case "varchar", "char", "character", "character varying", "text":
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

func collectTrinoASTFunctionCalls(root trinoparser.RoutineStatement, callerSchema string, params []trinoFuncParam, cat *trinoFuncCatalog) []functionCallRef {
	var calls []functionCallRef
	env := make(map[string]string, len(params))
	for _, p := range params {
		if p.name != "" {
			env[strings.ToLower(p.name)] = p.typeName
		}
	}

	var walkExpr func(expr trinoparser.Expr, curEnv map[string]string) string
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
			walkExpr(e.Base, curEnv)
		case *trinoparser.Subscript:
			walkExpr(e.Value, curEnv)
			walkExpr(e.Index, curEnv)
		case *trinoparser.Literal:
			switch e.Kind {
			case trinoparser.LiteralBool:
				return "boolean"
			case trinoparser.LiteralInteger:
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
		case *trinoparser.CastExpr:
			walkExpr(e.Expr, curEnv)
			return normalizeTrinoDataType(e.Type)
		case *trinoparser.UnaryExpr:
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
				if wt := walkExpr(w.Result, curEnv); wt != "" && res == "" {
					res = wt
				}
			}
			if et := walkExpr(e.Else, curEnv); et != "" && res == "" {
				res = et
			}
			return res
		case *trinoparser.AtTimeZoneExpr:
			t := walkExpr(e.Value, curEnv)
			walkExpr(e.Zone, curEnv)
			return t
		case *trinoparser.AtLocalExpr:
			return walkExpr(e.Value, curEnv)
		case *trinoparser.ArrayConstructor:
			for _, el := range e.Elements {
				walkExpr(el, curEnv)
			}
			return "array"
		case *trinoparser.RowConstructor:
			for _, el := range e.Elements {
				walkExpr(el, curEnv)
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
			for j, a := range e.Args {
				argTypes[j] = walkExpr(a, curEnv)
			}
			walkExpr(e.Filter, curEnv)
			var schema, name string
			if e.Name != nil && len(e.Name.Parts) > 0 {
				parts := e.Name.Parts
				switch len(parts) {
				case 1:
					name = parts[0].Value
				case 2:
					schema, name = parts[0].Value, parts[1].Value
				default:
					schema, name = parts[len(parts)-2].Value, parts[len(parts)-1].Value
				}
			}
			if name != "" {
				ref := functionCallRef{
					schema:   schema,
					name:     name,
					argTypes: argTypes,
					argCount: len(argTypes),
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
		if i+3 < len(toks) && toks[i].isIdent && toks[i+1].punct == '.' && toks[i+2].isIdent && toks[i+3].punct == '(' {
			calls = append(calls, functionCallRef{schema: toks[i].text, name: toks[i+2].text, argCount: countArgs(i + 3)})
			continue
		}
		if toks[i].isIdent && toks[i+1].punct == '(' && (i == 0 || toks[i-1].punct != '.') {
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
		index[objectKey(view.schema, view.view.Name)] = i
	}
	dependencies := make([]map[int]struct{}, len(views))
	addDep := func(i int, depSchema, depTable string) {
		if depSchema == "" {
			depSchema = views[i].schema
		}
		dependency, ok := index[objectKey(depSchema, depTable)]
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
			if access.Catalog != "" && metadata.Name != "" && access.Catalog != metadata.Name {
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
