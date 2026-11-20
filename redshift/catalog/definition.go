package catalog

import (
	"container/heap"
	"fmt"
	"io"
	"slices"
	"strings"
	"sync"

	metadatapb "github.com/bytebase/omni/metadata"
	nodes "github.com/bytebase/omni/redshift/ast"
	pgparser "github.com/bytebase/omni/redshift/parser"
)

// DatabaseDefinition renders a Redshift database snapshot from
// metadatapb.DatabaseSchemaMetadata as logical DDL (schemas, tables, views,
// materialized views, functions, procedures, comments, and owners). Physical
// storage attributes (DISTSTYLE, DISTKEY, SORTKEY, column ENCODE) are not
// captured by Redshift metadata sync into metadatapb; explicit unrendered table
// or column option fields fail rather than being silently dropped.
func DatabaseDefinition(to *metadatapb.DatabaseSchemaMetadata) (string, error) {
	if to == nil {
		return "", nil
	}
	if to.Datashare {
		return "", fmt.Errorf("datashare database %q cannot be restored because DatabaseSchemaMetadata does not capture FROM DATASHARE configuration", to.Name)
	}
	toState := convertToDatabaseState(to)

	var sb strings.Builder

	if needsCaseSensitiveIdentifiers(to) {
		if _, err := sb.WriteString("SET enable_case_sensitive_identifier TO true;\n"); err != nil {
			return "", err
		}
	}
	if err := writeSchemas(&sb, to); err != nil {
		return "", err
	}
	if err := writeOrderedObjects(&sb, to, toState); err != nil {
		return "", err
	}
	if err := writeProcedures(&sb, to); err != nil {
		return "", err
	}
	if err := writeOwners(&sb, to); err != nil {
		return "", err
	}

	s := sb.String()
	// Trim the leading newline so the output starts directly with the first statement or banner.
	s = strings.TrimLeft(s, "\n")
	return s, nil
}

// GetTableDefinition renders one table as DDL, matching DatabaseDefinition
// without the section banner.
func GetTableDefinition(schemaName string, table *metadatapb.TableMetadata) (string, error) {
	if table == nil || table.SkipDump {
		return "", nil
	}
	var sb strings.Builder
	if hasUpperASCII(schemaName) || tableHasUpperIdentifier(table) {
		if _, err := sb.WriteString("SET enable_case_sensitive_identifier TO true;\n\n"); err != nil {
			return "", err
		}
	}
	if err := writeTableWithOptions(&sb, convertToTableState(0, schemaName, table), nil, true, true); err != nil {
		return "", err
	}
	return sb.String(), nil
}

// GetViewDefinition renders one view as DDL, matching DatabaseDefinition
// without the section banner.
func GetViewDefinition(schemaName string, view *metadatapb.ViewMetadata) (string, error) {
	if view == nil || view.SkipDump {
		return "", nil
	}
	var sb strings.Builder
	if hasUpperASCII(schemaName) || viewHasUpperIdentifier(view) {
		if _, err := sb.WriteString("SET enable_case_sensitive_identifier TO true;\n\n"); err != nil {
			return "", err
		}
	}
	if err := writeView(&sb, convertToViewState(0, schemaName, view)); err != nil {
		return "", err
	}
	return sb.String(), nil
}

func hasUpperASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			return true
		}
	}
	return false
}

// hasUpperQuotedIdent reports whether sql contains a double-quoted identifier
// with an uppercase ASCII letter, skipping single-quoted string literals and
// SQL comments.
func hasUpperQuotedIdent(sql string) bool {
	blockDepth := 0
	for i := 0; i < len(sql); {
		if blockDepth > 0 {
			if sql[i] == '/' && i+1 < len(sql) && sql[i+1] == '*' {
				blockDepth++
				i += 2
				continue
			}
			if sql[i] == '*' && i+1 < len(sql) && sql[i+1] == '/' {
				blockDepth--
				i += 2
				continue
			}
			i++
			continue
		}
		switch {
		case sql[i] == '-' && i+1 < len(sql) && sql[i+1] == '-':
			i += 2
			for i < len(sql) && sql[i] != '\n' {
				i++
			}
		case sql[i] == '/' && i+1 < len(sql) && sql[i+1] == '*':
			blockDepth++
			i += 2
		case sql[i] == '\'':
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
		case sql[i] == '"':
			i++
			for i < len(sql) {
				if sql[i] == '"' {
					if i+1 < len(sql) && sql[i+1] == '"' {
						i += 2
						continue
					}
					i++
					break
				}
				if sql[i] >= 'A' && sql[i] <= 'Z' {
					return true
				}
				i++
			}
		default:
			i++
		}
	}
	return false
}

func viewHasUpperIdentifier(view *metadatapb.ViewMetadata) bool {
	if hasUpperASCII(view.Name) || hasUpperQuotedIdent(view.Definition) {
		return true
	}
	for _, col := range view.Columns {
		if col != nil && hasUpperASCII(col.Name) {
			return true
		}
	}
	return false
}

func tableHasUpperIdentifier(table *metadatapb.TableMetadata) bool {
	if hasUpperASCII(table.Name) || hasUpperASCII(table.Owner) {
		return true
	}
	for _, col := range table.Columns {
		if col != nil && hasUpperASCII(col.Name) {
			return true
		}
	}
	for _, idx := range table.Indexes {
		if idx != nil && (idx.Primary || idx.Unique) && hasUpperASCII(idx.Name) {
			return true
		}
	}
	for _, fk := range table.ForeignKeys {
		if fk == nil {
			continue
		}
		if hasUpperASCII(fk.Name) || hasUpperASCII(fk.ReferencedSchema) || hasUpperASCII(fk.ReferencedTable) {
			return true
		}
		for _, c := range fk.Columns {
			if hasUpperASCII(c) {
				return true
			}
		}
		for _, c := range fk.ReferencedColumns {
			if hasUpperASCII(c) {
				return true
			}
		}
	}
	return false
}

func needsCaseSensitiveIdentifiers(to *metadatapb.DatabaseSchemaMetadata) bool {
	if to.Owner != "" && (hasUpperASCII(to.Owner) || hasUpperASCII(to.Name)) {
		return true
	}
	for _, schema := range to.Schemas {
		if schema == nil || schema.SkipDump {
			continue
		}
		if hasUpperASCII(schema.Name) || hasUpperASCII(schema.Owner) {
			return true
		}
		for _, table := range schema.Tables {
			if table != nil && !table.SkipDump && tableHasUpperIdentifier(table) {
				return true
			}
		}
		for _, view := range schema.Views {
			if view != nil && !view.SkipDump && viewHasUpperIdentifier(view) {
				return true
			}
		}
		for _, mv := range schema.MaterializedViews {
			if mv != nil && !mv.SkipDump && (hasUpperASCII(mv.Name) || hasUpperQuotedIdent(mv.Definition)) {
				return true
			}
		}
		for _, fn := range schema.Functions {
			if fn != nil && !fn.SkipDump && (hasUpperASCII(fn.Name) || hasUpperQuotedIdent(fn.Definition)) {
				return true
			}
		}
		for _, proc := range schema.Procedures {
			if proc != nil && !proc.SkipDump && (hasUpperASCII(proc.Name) || hasUpperQuotedIdent(proc.Definition)) {
				return true
			}
		}
	}
	return false
}

type functionParam struct {
	name     string
	typeName string
}

type functionEntry struct {
	schema     string
	name       string
	def        string
	params     []functionParam
	returnType string
	depTables  []*metadatapb.DependencyTable
}

type funcCatalog struct {
	entries []functionEntry
	byName  map[string][]int
	nodeIdx []int
}

type functionCall struct {
	schema   string
	name     string
	argTypes []string
}

func (c *funcCatalog) resolveCall(callerSchema string, call functionCall) []int {
	if c == nil || len(c.entries) == 0 {
		return nil
	}
	var candidates []int
	if call.schema != "" {
		candidates = c.byName[objectKey(call.schema, call.name)]
	} else if cands, ok := c.byName[objectKey(callerSchema, call.name)]; ok {
		candidates = cands
	} else {
		candidates = c.byName[objectKey("public", call.name)]
	}
	return resolveFunctionOverload(c.entries, candidates, call.argTypes)
}

func resolveFunctionOverload(fns []functionEntry, candidates []int, argTypes []string) []int {
	if len(candidates) <= 1 {
		return candidates
	}
	var arityMatches []int
	for _, idx := range candidates {
		if len(fns[idx].params) == len(argTypes) {
			arityMatches = append(arityMatches, idx)
		}
	}
	if len(arityMatches) == 0 {
		return candidates
	}
	if len(arityMatches) == 1 {
		return arityMatches
	}
	bestScore := -1
	var best []int
	for _, idx := range arityMatches {
		score := 0
		compatible := true
		for i, argType := range argTypes {
			if argType == "" {
				continue
			}
			paramType := fns[idx].params[i].typeName
			if paramType == "" {
				continue
			}
			if paramType == argType {
				score += 2
			} else if argRank, paramRank := numericTypeRank(argType), numericTypeRank(paramType); argRank > 0 && paramRank >= argRank {
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

func isNumericType(t string) bool {
	return numericTypeRank(t) > 0
}

func numericTypeRank(t string) int {
	switch t {
	case "int2":
		return 1
	case "int4":
		return 2
	case "int8":
		return 3
	case "numeric":
		return 4
	case "float4":
		return 5
	case "float8":
		return 6
	default:
		return 0
	}
}

func normalizeTypeAlias(name string) string {
	switch strings.ToLower(name) {
	case "smallint", "int2":
		return "int2"
	case "int", "integer", "int4":
		return "int4"
	case "bigint", "int8":
		return "int8"
	case "decimal", "numeric":
		return "numeric"
	case "real", "float4":
		return "float4"
	case "double precision", "float8", "float":
		return "float8"
	case "boolean", "bool":
		return "bool"
	case "varchar", "character varying", "text", "char", "character", "bpchar", "nvarchar":
		return "text"
	default:
		return strings.ToLower(name)
	}
}

func normalizeSQLTypeName(tn *nodes.TypeName) string {
	if tn == nil || tn.Names == nil || len(tn.Names.Items) == 0 {
		return ""
	}
	var parts []string
	for _, item := range tn.Names.Items {
		if s, ok := item.(*nodes.String); ok && s.Str != "" {
			parts = append(parts, strings.ToLower(s.Str))
		}
	}
	if len(parts) == 0 {
		return ""
	}
	return normalizeTypeAlias(parts[len(parts)-1])
}

func parseFunctionSignature(definition string) (params []functionParam, returnType string) {
	stmts, err := pgparser.Parse(definition)
	if err != nil || stmts == nil || len(stmts.Items) != 1 {
		return nil, ""
	}
	stmt := stmts.Items[0]
	if raw, ok := stmt.(*nodes.RawStmt); ok {
		stmt = raw.Stmt
	}
	cf, ok := stmt.(*nodes.CreateFunctionStmt)
	if !ok {
		return nil, ""
	}
	returnType = normalizeSQLTypeName(cf.ReturnType)
	if cf.Parameters == nil {
		return nil, returnType
	}
	for _, item := range cf.Parameters.Items {
		fp, ok := item.(*nodes.FunctionParameter)
		if !ok {
			continue
		}
		params = append(params, functionParam{
			name:     fp.Name,
			typeName: normalizeSQLTypeName(fp.ArgType),
		})
	}
	return params, returnType
}

var (
	builtinCatalogOnce sync.Once
	builtinTypeByOID   map[uint32]string
	builtinProcsByName map[string][]BuiltinProc
)

func initBuiltinCatalog() {
	builtinCatalogOnce.Do(func() {
		builtinTypeByOID = make(map[uint32]string, len(BuiltinTypes))
		for _, t := range BuiltinTypes {
			builtinTypeByOID[t.OID] = normalizeTypeAlias(t.TypeName)
		}
		builtinProcsByName = make(map[string][]BuiltinProc, len(BuiltinProcs))
		for _, p := range BuiltinProcs {
			if p.Kind != 'f' {
				continue
			}
			key := strings.ToLower(p.Name)
			builtinProcsByName[key] = append(builtinProcsByName[key], p)
		}
	})
}

func resolveBuiltinProcReturnType(schema, name string, argTypes []string) string {
	if schema != "" && !strings.EqualFold(schema, "pg_catalog") {
		return ""
	}
	initBuiltinCatalog()
	procs := builtinProcsByName[strings.ToLower(name)]
	if len(procs) == 0 {
		return ""
	}
	bestScore := -1
	var matchedRet string
	ambiguous := false
	for _, p := range procs {
		if int(p.NArgs) != len(argTypes) || len(p.ArgTypes) != len(argTypes) {
			continue
		}
		score := 0
		compatible := true
		for i, argType := range argTypes {
			if argType == "" {
				continue
			}
			paramType := builtinTypeByOID[p.ArgTypes[i]]
			if paramType == "" || strings.HasPrefix(paramType, "any") {
				continue
			}
			if paramType == argType {
				score += 2
			} else if isNumericType(argType) && isNumericType(paramType) {
				score++
			} else {
				compatible = false
				break
			}
		}
		if !compatible {
			continue
		}
		retType := builtinTypeByOID[p.RetType]
		if retType == "" || strings.HasPrefix(retType, "any") {
			continue
		}
		switch {
		case score > bestScore:
			bestScore = score
			matchedRet = retType
			ambiguous = false
		case score == bestScore && matchedRet != retType:
			ambiguous = true
		}
	}
	if ambiguous {
		return ""
	}
	return matchedRet
}

func numericRank(t string) int {
	switch t {
	case "int2":
		return 1
	case "int4":
		return 2
	case "int8":
		return 3
	case "numeric":
		return 4
	case "float4":
		return 5
	case "float8":
		return 6
	default:
		return 0
	}
}

func commonNumericType(a, b string) string {
	ra, rb := numericRank(a), numericRank(b)
	if ra == 0 || rb == 0 {
		return ""
	}
	if ra >= rb {
		return a
	}
	return b
}

func extractNodeStringParts(list *nodes.List) []string {
	if list == nil {
		return nil
	}
	var parts []string
	for _, item := range list.Items {
		if s, ok := item.(*nodes.String); ok && s.Str != "" {
			parts = append(parts, s.Str)
		}
	}
	return parts
}

func inferArgTypes(args *nodes.List, callerSchema string, callerParams []functionParam, fnCat *funcCatalog) []string {
	if args == nil || len(args.Items) == 0 {
		return nil
	}
	out := make([]string, len(args.Items))
	for i, arg := range args.Items {
		out[i] = inferExprType(arg, callerSchema, callerParams, fnCat)
	}
	return out
}

func inferExprType(node nodes.Node, callerSchema string, callerParams []functionParam, fnCat *funcCatalog) string {
	switch n := node.(type) {
	case *nodes.ParamRef:
		idx := n.Number - 1
		if idx >= 0 && idx < len(callerParams) {
			return callerParams[idx].typeName
		}
	case *nodes.ColumnRef:
		if n.Fields != nil && len(n.Fields.Items) == 1 {
			if s, ok := n.Fields.Items[0].(*nodes.String); ok {
				for _, p := range callerParams {
					if strings.EqualFold(p.name, s.Str) {
						return p.typeName
					}
				}
			}
		}
	case *nodes.TypeCast:
		return normalizeSQLTypeName(n.TypeName)
	case *nodes.A_Const:
		if n.Isnull {
			return ""
		}
		switch n.Val.(type) {
		case *nodes.Integer:
			return "int4"
		case *nodes.Float:
			return "numeric"
		case *nodes.String:
			return "text"
		case *nodes.Boolean:
			return "bool"
		}
	case *nodes.BoolExpr, *nodes.NullTest, *nodes.BooleanTest:
		return "bool"
	case *nodes.CoalesceExpr:
		if n.Args != nil {
			return inferCommonExprTypes(n.Args.Items, callerSchema, callerParams, fnCat)
		}
	case *nodes.MinMaxExpr:
		if n.Args != nil {
			return inferCommonExprTypes(n.Args.Items, callerSchema, callerParams, fnCat)
		}
	case *nodes.CaseExpr:
		var branches []nodes.Node
		if n.Args != nil {
			for _, item := range n.Args.Items {
				if cw, ok := item.(*nodes.CaseWhen); ok && cw.Result != nil {
					branches = append(branches, cw.Result)
				}
			}
		}
		if n.Defresult != nil {
			branches = append(branches, n.Defresult)
		}
		return inferCommonExprTypes(branches, callerSchema, callerParams, fnCat)
	case *nodes.A_Expr:
		if n.Kind == nodes.AEXPR_NULLIF {
			return inferExprType(n.Lexpr, callerSchema, callerParams, fnCat)
		}
		if n.Kind != nodes.AEXPR_OP {
			return "bool"
		}
		parts := extractNodeStringParts(n.Name)
		if len(parts) == 0 {
			return ""
		}
		op := parts[len(parts)-1]
		switch op {
		case "=", "<>", "!=", "<", "<=", ">", ">=", "~", "~*", "!~", "!~*":
			return "bool"
		case "||":
			return "text"
		case "+", "-", "*", "/", "%", "^":
			if n.Lexpr == nil {
				return inferExprType(n.Rexpr, callerSchema, callerParams, fnCat)
			}
			return commonNumericType(
				inferExprType(n.Lexpr, callerSchema, callerParams, fnCat),
				inferExprType(n.Rexpr, callerSchema, callerParams, fnCat),
			)
		}
	case *nodes.FuncCall:
		parts := extractNodeStringParts(n.Funcname)
		var callSchema, callName string
		switch len(parts) {
		case 1:
			callName = parts[0]
		case 2:
			callSchema, callName = parts[0], parts[1]
		default:
			return ""
		}
		argTypes := inferArgTypes(n.Args, callerSchema, callerParams, fnCat)
		if matches := fnCat.resolveCall(callerSchema, functionCall{schema: callSchema, name: callName, argTypes: argTypes}); len(matches) > 0 {
			ret := fnCat.entries[matches[0]].returnType
			same := ret != ""
			for _, m := range matches[1:] {
				if fnCat.entries[m].returnType != ret {
					same = false
					break
				}
			}
			if same {
				return ret
			}
		}
		return resolveBuiltinProcReturnType(callSchema, callName, argTypes)
	}
	return ""
}

func inferCommonExprTypes(exprs []nodes.Node, callerSchema string, callerParams []functionParam, fnCat *funcCatalog) string {
	var result string
	for _, item := range exprs {
		t := inferExprType(item, callerSchema, callerParams, fnCat)
		if t == "" {
			continue
		}
		if result == "" || result == t {
			result = t
			continue
		}
		if common := commonNumericType(result, t); common != "" {
			result = common
		}
	}
	return result
}

func walkASTFunctionCalls(node nodes.Node, callerSchema string, callerParams []functionParam, fnCat *funcCatalog, calls *[]functionCall) {
	nodes.Inspect(node, func(n nodes.Node) bool {
		fc, ok := n.(*nodes.FuncCall)
		if !ok || fc.Funcname == nil || len(fc.Funcname.Items) == 0 {
			return true
		}
		parts := extractNodeStringParts(fc.Funcname)
		argTypes := inferArgTypes(fc.Args, callerSchema, callerParams, fnCat)
		switch len(parts) {
		case 1:
			*calls = append(*calls, functionCall{name: parts[0], argTypes: argTypes})
		case 2:
			*calls = append(*calls, functionCall{schema: parts[0], name: parts[1], argTypes: argTypes})
		}
		return true
	})
}

func collectFunctionCalls(definition, callerSchema string, callerParams []functionParam, fnCat *funcCatalog) []functionCall {
	var calls []functionCall
	stmts, err := pgparser.Parse(definition)
	if err == nil && stmts != nil {
		for _, item := range stmts.Items {
			stmt := item
			if raw, ok := stmt.(*nodes.RawStmt); ok {
				stmt = raw.Stmt
			}
			cf, ok := stmt.(*nodes.CreateFunctionStmt)
			if !ok {
				continue
			}
			if cf.SqlBody != nil {
				walkASTFunctionCalls(cf.SqlBody, callerSchema, callerParams, fnCat, &calls)
			}
			if cf.Options != nil {
				for _, optItem := range cf.Options.Items {
					de, ok := optItem.(*nodes.DefElem)
					if !ok || de.Defname != "as" {
						continue
					}
					for _, body := range extractDefElemStrings(de.Arg) {
						if bodyStmts, err := pgparser.Parse(body); err == nil && bodyStmts != nil {
							walkASTFunctionCalls(bodyStmts, callerSchema, callerParams, fnCat, &calls)
						} else {
							calls = append(calls, scanCallRefs(body)...)
						}
					}
				}
			}
		}
		return calls
	}
	return scanCallRefs(definition)
}

func collectSQLFunctionCalls(sql, callerSchema string, fnCat *funcCatalog) []functionCall {
	if strings.TrimSpace(sql) == "" {
		return nil
	}
	stmts, err := pgparser.Parse(sql)
	if err != nil || stmts == nil {
		return nil
	}
	var calls []functionCall
	walkASTFunctionCalls(stmts, callerSchema, nil, fnCat, &calls)
	return calls
}

func extractDefElemStrings(arg nodes.Node) []string {
	switch v := arg.(type) {
	case *nodes.String:
		return []string{v.Str}
	case *nodes.List:
		var out []string
		for _, item := range v.Items {
			if s, ok := item.(*nodes.String); ok {
				out = append(out, s.Str)
			}
		}
		return out
	}
	return nil
}

func scanCallRefs(body string) []functionCall {
	lex := pgparser.NewLexer(body)
	var toks []pgparser.Token
	for {
		tok := lex.NextToken()
		if tok.Type == 0 {
			break
		}
		toks = append(toks, tok)
	}
	countArgs := func(openIdx int) int {
		if openIdx+1 < len(toks) && toks[openIdx+1].Type == ')' {
			return 0
		}
		depth := 0
		args := 1
		for j := openIdx; j < len(toks); j++ {
			switch toks[j].Type {
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
	var calls []functionCall
	for i := 0; i+1 < len(toks); i++ {
		if i+3 < len(toks) && toks[i+1].Type == '.' && toks[i+3].Type == '(' && toks[i].Str != "" && toks[i+2].Str != "" {
			calls = append(calls, functionCall{schema: toks[i].Str, name: toks[i+2].Str, argTypes: make([]string, countArgs(i+3))})
			continue
		}
		if toks[i+1].Type == '(' && toks[i].Str != "" && (i == 0 || toks[i-1].Type != '.') {
			calls = append(calls, functionCall{name: toks[i].Str, argTypes: make([]string, countArgs(i+1))})
		}
	}
	return calls
}

func writeProcedures(w io.StringWriter, to *metadatapb.DatabaseSchemaMetadata) error {
	for _, schema := range to.Schemas {
		if schema == nil || schema.SkipDump {
			continue
		}
		for _, proc := range schema.Procedures {
			if proc == nil || proc.SkipDump {
				continue
			}
			def := strings.TrimSpace(proc.Definition)
			if def == "" {
				return fmt.Errorf("procedure %s has no definition to restore", qualifyName(schema.Name, proc.Name))
			}
			if proc.Comment != "" {
				return fmt.Errorf("procedure %s comment cannot be restored because Redshift does not support COMMENT ON PROCEDURE", qualifyName(schema.Name, proc.Name))
			}
			def = appendTerminator(def)
			if _, err := w.WriteString(getProcedureAnnouncement(proc.Name) + def + "\n"); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeSchemas(w io.StringWriter, to *metadatapb.DatabaseSchemaMetadata) error {
	for _, schema := range to.Schemas {
		if schema == nil || schema.SkipDump || schema.Name == "" {
			continue
		}
		if schema.Name != "public" {
			if _, err := w.WriteString(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;\n", quoteDefinitionIdentifier(schema.Name))); err != nil {
				return err
			}
		}
		if schema.Comment != "" {
			if _, err := w.WriteString(fmt.Sprintf("COMMENT ON SCHEMA %s IS '%s';\n", quoteDefinitionIdentifier(schema.Name), escapeStringLiteral(schema.Comment))); err != nil {
				return err
			}
		}
	}
	return nil
}

// writeOwners transfers table and schema ownership at the end of the dump so
// all creates, comments, views, and cyclic foreign keys run while the restore
// role still owns the targets.
func writeOwners(w io.StringWriter, to *metadatapb.DatabaseSchemaMetadata) error {
	for _, schema := range to.Schemas {
		if schema == nil || schema.SkipDump {
			continue
		}
		for _, table := range schema.Tables {
			if table == nil || table.SkipDump || table.Owner == "" {
				continue
			}
			if _, err := w.WriteString(fmt.Sprintf("ALTER TABLE %s OWNER TO %s;\n", qualifyName(schema.Name, table.Name), quoteDefinitionIdentifier(table.Owner))); err != nil {
				return err
			}
		}
		if schema.Name != "" && schema.Owner != "" {
			if _, err := w.WriteString(fmt.Sprintf("ALTER SCHEMA %s OWNER TO %s;\n", quoteDefinitionIdentifier(schema.Name), quoteDefinitionIdentifier(schema.Owner))); err != nil {
				return err
			}
		}
	}
	if to.Owner != "" {
		if to.Name == "" {
			return fmt.Errorf("database owner %q cannot be restored because database name is empty", to.Owner)
		}
		if _, err := w.WriteString(fmt.Sprintf("ALTER DATABASE %s OWNER TO %s;\n", quoteDefinitionIdentifier(to.Name), quoteDefinitionIdentifier(to.Owner))); err != nil {
			return err
		}
	}
	return nil
}

// appendTerminator appends ";" when stmt does not already end in a
// statement terminator, placing it on a new line if stmt ends inside a "--"
// line comment.
func appendTerminator(stmt string) string {
	var (
		inSingleQuote bool
		inDoubleQuote bool
		inLineComment bool
		blockDepth    int
		dollarTag     string
		lastCodeByte  byte
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
		if blockDepth > 0 {
			if ch == '/' && i+1 < len(stmt) && stmt[i+1] == '*' {
				blockDepth++
				i += 2
				continue
			}
			if ch == '*' && i+1 < len(stmt) && stmt[i+1] == '/' {
				blockDepth--
				i += 2
				continue
			}
			i++
			continue
		}
		if dollarTag != "" {
			if strings.HasPrefix(stmt[i:], dollarTag) {
				i += len(dollarTag)
				dollarTag = ""
				lastCodeByte = '$'
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
			blockDepth++
			i += 2
		case ch == '\'':
			inSingleQuote = true
			lastCodeByte = ch
			i++
		case ch == '"':
			inDoubleQuote = true
			lastCodeByte = ch
			i++
		case ch == '$':
			if tag, ok := matchDollarTag(stmt[i:]); ok {
				dollarTag = tag
				lastCodeByte = '$'
				i += len(tag)
			} else {
				lastCodeByte = ch
				i++
			}
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

func matchDollarTag(s string) (string, bool) {
	if len(s) < 2 || s[0] != '$' {
		return "", false
	}
	for i := 1; i < len(s); i++ {
		c := s[i]
		if c == '$' {
			return s[:i+1], true
		}
		if i == 1 && c >= '0' && c <= '9' {
			return "", false
		}
		if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') && (c < '0' || c > '9') && c != '_' {
			return "", false
		}
	}
	return "", false
}

type objectKind int

const (
	objectKindFunc objectKind = iota
	objectKindTable
	objectKindView
	objectKindMaterializedView
)

type orderedObject struct {
	kind         objectKind
	schema       string
	name         string
	fn           functionEntry
	table        *tableState
	view         *viewState
	viewMeta     *metadatapb.ViewMetadata
	mvMeta       *metadatapb.MaterializedViewMetadata
	materialized bool
}

func isPythonUDF(definition string) (string, bool) {
	stmts, err := pgparser.Parse(definition)
	if err == nil && stmts != nil {
		for _, item := range stmts.Items {
			stmt := item
			if raw, ok := stmt.(*nodes.RawStmt); ok {
				stmt = raw.Stmt
			}
			cf, ok := stmt.(*nodes.CreateFunctionStmt)
			if !ok || cf.Options == nil {
				continue
			}
			for _, optItem := range cf.Options.Items {
				de, ok := optItem.(*nodes.DefElem)
				if !ok || !strings.EqualFold(de.Defname, "language") {
					continue
				}
				if s, ok := de.Arg.(*nodes.String); ok && strings.HasPrefix(strings.ToLower(s.Str), "plpython") {
					return s.Str, true
				}
			}
		}
		return "", false
	}
	lex := pgparser.NewLexer(definition)
	prevLang := false
	for {
		tok := lex.NextToken()
		if tok.Type == 0 {
			return "", false
		}
		if prevLang && tok.Str != "" && strings.HasPrefix(strings.ToLower(tok.Str), "plpython") {
			return tok.Str, true
		}
		prevLang = strings.EqualFold(tok.Str, "language")
	}
}

// writeOrderedObjects emits functions, tables, views, and materialized views
// in one topological order so cross-kind dependencies (such as view-dependent
// UDFs) always follow their prerequisites.
func writeOrderedObjects(w io.StringWriter, to *metadatapb.DatabaseSchemaMetadata, state *databaseState) error {
	var items []orderedObject
	var fnEntries []functionEntry
	var fnNodeIndices []int

	for _, schema := range to.Schemas {
		if schema == nil || schema.SkipDump {
			continue
		}
		for _, et := range schema.ExternalTables {
			if et != nil && !et.SkipDump {
				return fmt.Errorf("external table %s cannot be restored because SchemaMetadata does not capture external schema definitions", qualifyName(schema.Name, et.Name))
			}
		}
		for _, seq := range schema.Sequences {
			if seq != nil && !seq.SkipDump {
				return fmt.Errorf("sequence %s cannot be restored because Redshift does not support PostgreSQL sequences", qualifyName(schema.Name, seq.Name))
			}
		}
		for _, fn := range schema.Functions {
			if fn == nil || fn.SkipDump {
				continue
			}
			if fn.Comment != "" {
				return fmt.Errorf("function %s comment cannot be restored because Redshift does not support COMMENT ON FUNCTION", qualifyName(schema.Name, fn.Name))
			}
			def := strings.TrimSpace(fn.Definition)
			if def == "" {
				return fmt.Errorf("function %s has no definition to restore", qualifyName(schema.Name, fn.Name))
			}
			if lang, isPy := isPythonUDF(def); isPy {
				return fmt.Errorf("function %s uses language %s, which Redshift no longer supports", qualifyName(schema.Name, fn.Name), lang)
			}
			def = appendTerminator(def)
			params, retType := parseFunctionSignature(def)
			entry := functionEntry{
				schema:     schema.Name,
				name:       fn.Name,
				def:        def,
				params:     params,
				returnType: retType,
				depTables:  fn.DependencyTables,
			}
			fnEntries = append(fnEntries, entry)
			fnNodeIndices = append(fnNodeIndices, len(items))
			items = append(items, orderedObject{
				kind:   objectKindFunc,
				schema: schema.Name,
				name:   fn.Name,
				fn:     entry,
			})
		}
	}

	knownTables := make(map[string]bool)
	tablesRemaining := 0
	for _, schema := range to.Schemas {
		if schema == nil || schema.SkipDump {
			continue
		}
		schemaState, ok := state.schemas[schema.Name]
		if !ok {
			continue
		}
		for _, table := range schema.Tables {
			if table == nil || table.SkipDump {
				continue
			}
			ts, ok := schemaState.tables[table.Name]
			if !ok {
				continue
			}
			delete(schemaState.tables, table.Name)
			knownTables[objectKey(ts.schema, ts.name)] = true
			tablesRemaining++
			items = append(items, orderedObject{
				kind:   objectKindTable,
				schema: ts.schema,
				name:   ts.name,
				table:  ts,
			})
		}
	}

	for _, schema := range to.Schemas {
		if schema == nil || schema.SkipDump {
			continue
		}
		schemaState, ok := state.schemas[schema.Name]
		if !ok {
			continue
		}
		for _, view := range schema.Views {
			if view == nil || view.SkipDump {
				continue
			}
			vs, ok := schemaState.views[view.Name]
			if !ok {
				continue
			}
			delete(schemaState.views, view.Name)
			items = append(items, orderedObject{
				kind:     objectKindView,
				schema:   vs.schema,
				name:     vs.name,
				view:     vs,
				viewMeta: view,
			})
		}
	}

	for _, schema := range to.Schemas {
		if schema == nil || schema.SkipDump {
			continue
		}
		schemaState, ok := state.schemas[schema.Name]
		if !ok {
			continue
		}
		for _, mv := range schema.MaterializedViews {
			if mv == nil || mv.SkipDump {
				continue
			}
			vs, ok := schemaState.materializedViews[mv.Name]
			if !ok {
				continue
			}
			delete(schemaState.materializedViews, mv.Name)
			items = append(items, orderedObject{
				kind:         objectKindMaterializedView,
				schema:       vs.schema,
				name:         vs.name,
				view:         vs,
				mvMeta:       mv,
				materialized: true,
			})
		}
	}

	if len(items) == 0 {
		return nil
	}

	fnByName := make(map[string][]int, len(fnEntries))
	for i, fn := range fnEntries {
		key := objectKey(fn.schema, fn.name)
		fnByName[key] = append(fnByName[key], i)
	}
	fnCat := &funcCatalog{
		entries: fnEntries,
		byName:  fnByName,
		nodeIdx: fnNodeIndices,
	}

	relIndex := make(map[string]int, len(items))
	for i, item := range items {
		if item.kind == objectKindFunc {
			continue
		}
		relIndex[objectKey(item.schema, item.name)] = i
	}

	dependencies := make([]map[int]struct{}, len(items))
	addDep := func(i, dep int) {
		if dep == i {
			return
		}
		if dependencies[i] == nil {
			dependencies[i] = map[int]struct{}{}
		}
		dependencies[i][dep] = struct{}{}
	}
	addRelDep := func(i int, defaultSchema, relCatalog, relSchema, relName string) {
		if relName == "" {
			return
		}
		if relCatalog != "" && (to.Name == "" || !strings.EqualFold(relCatalog, to.Name)) {
			return
		}
		if relSchema != "" {
			if depIdx, ok := relIndex[objectKey(relSchema, relName)]; ok {
				addDep(i, depIdx)
			}
			return
		}
		if depIdx, ok := relIndex[objectKey(defaultSchema, relName)]; ok {
			addDep(i, depIdx)
		} else if depIdx, ok := relIndex[objectKey("public", relName)]; ok {
			addDep(i, depIdx)
		}
	}
	addCallDeps := func(i int, callerSchema string, calls []functionCall) {
		for _, call := range calls {
			for _, fnEntryIdx := range fnCat.resolveCall(callerSchema, call) {
				addDep(i, fnCat.nodeIdx[fnEntryIdx])
			}
		}
	}

	for i, item := range items {
		switch item.kind {
		case objectKindFunc:
			for _, dep := range item.fn.depTables {
				if dep == nil {
					continue
				}
				addRelDep(i, item.schema, "", dep.Schema, dep.Table)
			}
			addCallDeps(i, item.schema, collectFunctionCalls(item.fn.def, item.schema, item.fn.params, fnCat))
		case objectKindTable:
			for _, fk := range item.table.foreignKeys {
				refSchema := fk.referencedSchema
				if refSchema == "" {
					refSchema = item.schema
				}
				if depIdx, ok := relIndex[objectKey(refSchema, fk.referencedTable)]; ok {
					addDep(i, depIdx)
				}
			}
		case objectKindView, objectKindMaterializedView:
			// Late-binding views (WITH NO SCHEMA BINDING) resolve relations at query
			// time, so referenced relations do not need to exist prior to creation.
			if item.kind != objectKindView || !isLateBindingView(item.view.definition) {
				var depCols []*metadatapb.DependencyColumn
				if item.viewMeta != nil {
					depCols = item.viewMeta.DependencyColumns
				} else if item.mvMeta != nil {
					depCols = item.mvMeta.DependencyColumns
				}
				for _, dep := range depCols {
					if dep == nil {
						continue
					}
					addRelDep(i, item.schema, "", dep.Schema, dep.Table)
				}
				for _, ref := range collectViewRelationRefs(item.view.definition) {
					addRelDep(i, item.schema, ref.catalog, ref.schema, ref.name)
				}
			}
			addCallDeps(i, item.schema, collectSQLFunctionCalls(item.view.definition, item.schema, fnCat))
		}
	}

	for _, item := range items {
		if item.kind != objectKindTable {
			continue
		}
		for _, expr := range item.table.defaultExprs {
			for _, call := range collectSQLFunctionCalls("SELECT "+expr, item.schema, fnCat) {
				if matches := fnCat.resolveCall(item.schema, call); len(matches) > 0 {
					fnNode := fnCat.nodeIdx[matches[0]]
					return fmt.Errorf("table %s column default cannot call user-defined function %s because Redshift disallows UDFs in default_expr", item.table.qualifiedName(), qualifyName(items[fnNode].schema, items[fnNode].name))
				}
			}
		}
	}

	type deferredFK struct {
		table *tableState
		fk    *foreignKeyState
	}
	var deferred []deferredFK
	emittedTables := make(map[string]bool, len(knownTables))
	flushDeferred := func() error {
		if len(deferred) == 0 {
			return nil
		}
		if _, err := w.WriteString("\n"); err != nil {
			return err
		}
		for _, item := range deferred {
			var buf strings.Builder
			if err := item.fk.toString(&buf); err != nil {
				return err
			}
			if _, err := w.WriteString(fmt.Sprintf("ALTER TABLE %s ADD %s;\n", item.table.qualifiedName(), buf.String())); err != nil {
				return err
			}
			if item.fk.name != "" && item.fk.comment != "" {
				if _, err := w.WriteString(fmt.Sprintf("COMMENT ON CONSTRAINT %s ON %s IS '%s';\n", quoteDefinitionIdentifier(item.fk.name), item.table.qualifiedName(), escapeStringLiteral(item.fk.comment))); err != nil {
					return err
				}
			}
		}
		deferred = nil
		return nil
	}

	for _, idx := range topoOrder(len(items), dependencies) {
		item := items[idx]
		switch item.kind {
		case objectKindFunc:
			if _, err := w.WriteString(getFunctionAnnouncement(item.fn.name)); err != nil {
				return err
			}
			needsPath := item.schema != "" && item.schema != "public"
			if needsPath {
				if _, err := w.WriteString(fmt.Sprintf("SET search_path TO %s, public;\n", quoteDefinitionIdentifier(item.schema))); err != nil {
					return err
				}
			}
			if _, err := w.WriteString(item.fn.def + "\n"); err != nil {
				return err
			}
			if needsPath {
				if _, err := w.WriteString("RESET search_path;\n"); err != nil {
					return err
				}
			}
		case objectKindTable:
			if _, err := w.WriteString(getTableAnnouncement(item.table.name)); err != nil {
				return err
			}
			tableKey := objectKey(item.table.schema, item.table.name)
			var inlineFKs []*foreignKeyState
			for _, fk := range sortedForeignKeys(item.table.foreignKeys) {
				refSchema := fk.referencedSchema
				if refSchema == "" {
					refSchema = item.table.schema
				}
				refKey := objectKey(refSchema, fk.referencedTable)
				// Defer any FK whose target table exists in this dump and has not been
				// created yet (cyclic references) so every CREATE TABLE stays valid.
				if refKey != tableKey && knownTables[refKey] && !emittedTables[refKey] {
					deferred = append(deferred, deferredFK{table: item.table, fk: fk})
					continue
				}
				inlineFKs = append(inlineFKs, fk)
			}
			if err := writeTableWithOptions(w, item.table, inlineFKs, false, false); err != nil {
				return err
			}
			emittedTables[tableKey] = true
			tablesRemaining--
			if tablesRemaining == 0 {
				if err := flushDeferred(); err != nil {
					return err
				}
			}
		case objectKindView:
			if _, err := w.WriteString(getViewAnnouncement(item.view.name)); err != nil {
				return err
			}
			if err := writeView(w, item.view); err != nil {
				return err
			}
		case objectKindMaterializedView:
			if _, err := w.WriteString(getMaterializedViewAnnouncement(item.view.name)); err != nil {
				return err
			}
			if err := writeMaterializedView(w, item.view); err != nil {
				return err
			}
		}
	}
	return flushDeferred()
}

func objectKey(schemaName, objectName string) string {
	return schemaName + "\x00" + objectName
}

type relationRef struct {
	catalog string
	schema  string
	name    string
}

func collectViewRelationRefs(definition string) []relationRef {
	if strings.TrimSpace(definition) == "" {
		return nil
	}
	stmts, err := pgparser.Parse(definition)
	if err != nil || stmts == nil {
		return nil
	}
	var refs []relationRef
	nodes.Walk(&scopedRelVisitor{refs: &refs}, stmts)
	return refs
}

type scopedRelVisitor struct {
	cteScope map[string]bool
	refs     *[]relationRef
}

func (v *scopedRelVisitor) Visit(node nodes.Node) nodes.Visitor {
	if node == nil {
		return nil
	}
	switch n := node.(type) {
	case *nodes.RangeVar:
		if n.Relname != "" && (n.Catalogname != "" || n.Schemaname != "" || !v.cteScope[n.Relname]) {
			*v.refs = append(*v.refs, relationRef{catalog: n.Catalogname, schema: n.Schemaname, name: n.Relname})
		}
		return v
	case *nodes.SelectStmt:
		if n.WithClause == nil {
			return v
		}
		childScope := make(map[string]bool, len(v.cteScope))
		for k, val := range v.cteScope {
			childScope[k] = val
		}
		if n.WithClause.Recursive && n.WithClause.Ctes != nil {
			for _, item := range n.WithClause.Ctes.Items {
				if cte, ok := item.(*nodes.CommonTableExpr); ok && cte.Ctename != "" {
					childScope[cte.Ctename] = true
				}
			}
		}
		if n.WithClause.Ctes != nil {
			for _, item := range n.WithClause.Ctes.Items {
				cte, ok := item.(*nodes.CommonTableExpr)
				if !ok {
					continue
				}
				if cte.Ctequery != nil {
					cteVisitor := v
					if n.WithClause.Recursive || len(childScope) > len(v.cteScope) {
						cteVisitor = &scopedRelVisitor{cteScope: childScope, refs: v.refs}
					}
					nodes.Walk(cteVisitor, cte.Ctequery)
				}
				if !n.WithClause.Recursive && cte.Ctename != "" {
					childScope[cte.Ctename] = true
				}
			}
		}
		child := &scopedRelVisitor{cteScope: childScope, refs: v.refs}
		walkNodeList(child, n.DistinctClause)
		if n.IntoClause != nil {
			nodes.Walk(child, n.IntoClause)
		}
		walkNodeList(child, n.TargetList)
		walkNodeList(child, n.FromClause)
		nodes.Walk(child, n.WhereClause)
		walkNodeList(child, n.GroupClause)
		nodes.Walk(child, n.HavingClause)
		nodes.Walk(child, n.StartWith)
		nodes.Walk(child, n.ConnectBy)
		nodes.Walk(child, n.QualifyClause)
		walkNodeList(child, n.WindowClause)
		walkNodeList(child, n.ExcludeList)
		walkNodeList(child, n.ValuesLists)
		walkNodeList(child, n.SortClause)
		nodes.Walk(child, n.LimitOffset)
		nodes.Walk(child, n.LimitCount)
		walkNodeList(child, n.LockingClause)
		if n.Larg != nil {
			nodes.Walk(child, n.Larg)
		}
		if n.Rarg != nil {
			nodes.Walk(child, n.Rarg)
		}
		return nil
	default:
		return v
	}
}

func walkNodeList(v nodes.Visitor, list *nodes.List) {
	if list == nil {
		return
	}
	for _, item := range list.Items {
		nodes.Walk(v, item)
	}
}

type intMinHeap []int

func (h intMinHeap) Len() int           { return len(h) }
func (h intMinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h intMinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *intMinHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *intMinHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// topoOrder returns indices 0..n-1 ordered so each item follows its
// dependencies, preserving input order among independent items and breaking
// cycles in input order in O((n + E) log n) time.
func topoOrder(n int, dependencies []map[int]struct{}) []int {
	remaining := make([]int, n)
	dependents := make([][]int, n)
	ready := make(intMinHeap, 0, n)
	for i, deps := range dependencies {
		remaining[i] = len(deps)
		if len(deps) == 0 {
			ready = append(ready, i)
			continue
		}
		for dep := range deps {
			dependents[dep] = append(dependents[dep], i)
		}
	}
	heap.Init(&ready)

	ordered := make([]int, 0, n)
	done := make([]bool, n)
	cursor := 0
	for len(ordered) < n {
		var u int
		if ready.Len() > 0 {
			u = heap.Pop(&ready).(int)
		} else {
			for cursor < n && done[cursor] {
				cursor++
			}
			if cursor >= n {
				break
			}
			u = cursor
			cursor++
		}
		if done[u] {
			continue
		}
		done[u] = true
		ordered = append(ordered, u)
		for _, v := range dependents[u] {
			if done[v] {
				continue
			}
			remaining[v]--
			if remaining[v] == 0 {
				heap.Push(&ready, v)
			}
		}
	}
	return ordered
}

// writeTableWithOptions writes one table, its comments (table, column, and
// named key/FK constraints), and optionally its ownership transfer. Column
// comments follow the column order the CREATE TABLE body uses so output is
// deterministic across runs.
func writeTableWithOptions(w io.StringWriter, table *tableState, foreignKeys []*foreignKeyState, allFKs, emitOwner bool) error {
	if allFKs {
		foreignKeys = sortedForeignKeys(table.foreignKeys)
	}
	if table.duplicateColumn != "" {
		return fmt.Errorf("table %s has duplicate column %s", table.qualifiedName(), quoteDefinitionIdentifier(table.duplicateColumn))
	}
	if table.hasCheckConstraints {
		return fmt.Errorf("table %s check constraints cannot be restored because Redshift does not support CHECK constraints", table.qualifiedName())
	}
	if table.unsupportedReason != "" {
		return fmt.Errorf("table %s %s", table.qualifiedName(), table.unsupportedReason)
	}
	if len(table.columns) == 0 {
		return fmt.Errorf("table %s has no columns", table.qualifiedName())
	}
	for _, idx := range table.indexes {
		if (idx.primary || idx.unique) && len(idx.keys) == 0 {
			return fmt.Errorf("index %s on %s has no columns", idx.name, table.qualifiedName())
		}
	}
	for _, fk := range sortedForeignKeys(table.foreignKeys) {
		if len(fk.columns) == 0 || len(fk.columns) != len(fk.referencedColumns) {
			return fmt.Errorf("foreign key %s on %s has %d columns and %d referenced columns", fk.name, table.qualifiedName(), len(fk.columns), len(fk.referencedColumns))
		}
		if !isSupportedRedshiftFKAction(fk.onDelete) {
			return fmt.Errorf("foreign key %s on %s has unsupported ON DELETE action %q", fk.name, table.qualifiedName(), fk.onDelete)
		}
		if !isSupportedRedshiftFKAction(fk.onUpdate) {
			return fmt.Errorf("foreign key %s on %s has unsupported ON UPDATE action %q", fk.name, table.qualifiedName(), fk.onUpdate)
		}
	}
	buf := &strings.Builder{}
	if err := table.toStringWithFKs(buf, foreignKeys); err != nil {
		return err
	}
	if table.comment != "" {
		if _, err := fmt.Fprintf(buf, "COMMENT ON TABLE %s IS '%s';\n", table.qualifiedName(), escapeStringLiteral(table.comment)); err != nil {
			return err
		}
	}
	for _, column := range sortedColumns(table.columns) {
		if column.comment == "" {
			continue
		}
		if _, err := fmt.Fprintf(buf, "COMMENT ON COLUMN %s.%s IS '%s';\n", table.qualifiedName(), quoteDefinitionIdentifier(column.name), escapeStringLiteral(column.comment)); err != nil {
			return err
		}
	}
	for _, idx := range sortedKeyIndexes(table.indexes) {
		if idx.name == "" || idx.comment == "" {
			continue
		}
		if _, err := fmt.Fprintf(buf, "COMMENT ON CONSTRAINT %s ON %s IS '%s';\n", quoteDefinitionIdentifier(idx.name), table.qualifiedName(), escapeStringLiteral(idx.comment)); err != nil {
			return err
		}
	}
	for _, fk := range foreignKeys {
		if fk.name == "" || fk.comment == "" {
			continue
		}
		if _, err := fmt.Fprintf(buf, "COMMENT ON CONSTRAINT %s ON %s IS '%s';\n", quoteDefinitionIdentifier(fk.name), table.qualifiedName(), escapeStringLiteral(fk.comment)); err != nil {
			return err
		}
	}
	if emitOwner && table.owner != "" {
		if _, err := fmt.Fprintf(buf, "ALTER TABLE %s OWNER TO %s;\n", table.qualifiedName(), quoteDefinitionIdentifier(table.owner)); err != nil {
			return err
		}
	}
	_, err := w.WriteString(buf.String())
	return err
}

// writeView writes one view and the comments that belong to it and its columns.
func writeView(w io.StringWriter, view *viewState) error {
	buf := &strings.Builder{}
	if err := view.toString(buf); err != nil {
		return err
	}
	if view.comment != "" {
		if _, err := fmt.Fprintf(buf, "COMMENT ON VIEW %s IS '%s';\n", view.qualifiedName(), escapeStringLiteral(view.comment)); err != nil {
			return err
		}
	}
	if !isLateBindingView(view.definition) {
		for _, column := range sortedColumns(view.columns) {
			if column.comment == "" {
				continue
			}
			if _, err := fmt.Fprintf(buf, "COMMENT ON COLUMN %s.%s IS '%s';\n", view.qualifiedName(), quoteDefinitionIdentifier(column.name), escapeStringLiteral(column.comment)); err != nil {
				return err
			}
		}
	}
	_, err := w.WriteString(buf.String())
	return err
}

func writeMaterializedView(w io.StringWriter, mv *viewState) error {
	if mv.comment != "" {
		return fmt.Errorf("materialized view %s comment cannot be restored because Redshift does not support COMMENT ON MATERIALIZED VIEW", mv.qualifiedName())
	}
	buf := &strings.Builder{}
	if err := mv.toMaterializedString(buf); err != nil {
		return err
	}
	_, err := w.WriteString(buf.String())
	return err
}

func sanitizeBannerName(name string) string {
	return strings.NewReplacer("\r\n", " ", "\r", " ", "\n", " ").Replace(name)
}

func getFunctionAnnouncement(name string) string {
	return fmt.Sprintf("\n--\n-- Function structure for `%s`\n--\n", sanitizeBannerName(name))
}

func getProcedureAnnouncement(name string) string {
	return fmt.Sprintf("\n--\n-- Procedure structure for `%s`\n--\n", sanitizeBannerName(name))
}

func getTableAnnouncement(name string) string {
	return fmt.Sprintf("\n--\n-- Table structure for `%s`\n--\n", sanitizeBannerName(name))
}

func getViewAnnouncement(name string) string {
	return fmt.Sprintf("\n--\n-- View structure for `%s`\n--\n", sanitizeBannerName(name))
}

func getMaterializedViewAnnouncement(name string) string {
	return fmt.Sprintf("\n--\n-- Materialized view structure for `%s`\n--\n", sanitizeBannerName(name))
}

type databaseState struct {
	name    string
	schemas map[string]*schemaState
}

func newDatabaseState() *databaseState {
	return &databaseState{
		schemas: make(map[string]*schemaState),
	}
}

func convertToDatabaseState(database *metadatapb.DatabaseSchemaMetadata) *databaseState {
	state := newDatabaseState()
	state.name = database.Name
	for _, schema := range database.Schemas {
		if schema == nil || schema.SkipDump {
			continue
		}
		state.schemas[schema.Name] = convertToSchemaState(schema)
	}
	return state
}

type schemaState struct {
	name              string
	tables            map[string]*tableState
	views             map[string]*viewState
	materializedViews map[string]*viewState
}

// qualifyName prefixes the object with its schema, so DDL for analytics.orders
// cannot land in public. The prefix is dropped when no schema is known, which is
// how the single-object entry points behave when the caller passes none.
func qualifyName(schemaName, objectName string) string {
	if schemaName == "" {
		return quoteDefinitionIdentifier(objectName)
	}
	return quoteDefinitionIdentifier(schemaName) + "." + quoteDefinitionIdentifier(objectName)
}

// quoteDefinitionIdentifier renders a name as a quoted SQL identifier. Redshift
// folds an unquoted name to lower case and rejects one containing a space or a
// reserved word outright, so a name the catalog reports verbatim has to be
// quoted to come back as itself. An embedded double quote is doubled, as in
// PostgreSQL.
func quoteDefinitionIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

func (t *tableState) qualifiedName() string {
	return qualifyName(t.schema, t.name)
}

func (v *viewState) qualifiedName() string {
	return qualifyName(v.schema, v.name)
}

// escapeStringLiteral encodes a comment or raw default string for a Redshift
// SQL string literal, doubling both backslashes and single quotes.
func escapeStringLiteral(s string) string {
	return strings.NewReplacer(`\`, `\\`, `'`, `''`).Replace(s)
}

func newSchemaState() *schemaState {
	return &schemaState{
		tables:            make(map[string]*tableState),
		views:             make(map[string]*viewState),
		materializedViews: make(map[string]*viewState),
	}
}

func convertToSchemaState(schema *metadatapb.SchemaMetadata) *schemaState {
	state := newSchemaState()
	state.name = schema.Name
	for i, table := range schema.Tables {
		if table == nil || table.SkipDump {
			continue
		}
		state.tables[table.Name] = convertToTableState(i, schema.Name, table)
	}
	for i, view := range schema.Views {
		if view == nil || view.SkipDump {
			continue
		}
		state.views[view.Name] = convertToViewState(i, schema.Name, view)
	}
	for i, mv := range schema.MaterializedViews {
		if mv == nil || mv.SkipDump {
			continue
		}
		state.materializedViews[mv.Name] = convertToMaterializedViewState(i, schema.Name, mv)
	}
	return state
}

type tableState struct {
	id                  int
	schema              string
	name                string
	owner               string
	duplicateColumn     string
	hasCheckConstraints bool
	unsupportedReason   string
	columns             map[string]*columnState
	defaultExprs        []string
	indexes             []*indexState
	foreignKeys         []*foreignKeyState
	comment             string
}

func sortedForeignKeys(fks []*foreignKeyState) []*foreignKeyState {
	foreignKeys := slices.Clone(fks)
	slices.SortFunc(foreignKeys, func(a, b *foreignKeyState) int {
		if cmp := strings.Compare(a.name, b.name); cmp != 0 {
			return cmp
		}
		if a.id < b.id {
			return -1
		}
		if a.id > b.id {
			return 1
		}
		return 0
	})
	return foreignKeys
}

func (t *tableState) toStringWithFKs(buf *strings.Builder, foreignKeys []*foreignKeyState) error {
	if _, err := fmt.Fprintf(buf, "CREATE TABLE %s (\n  ", t.qualifiedName()); err != nil {
		return err
	}
	columns := sortedColumns(t.columns)
	for i, column := range columns {
		if i > 0 {
			if _, err := buf.WriteString(",\n  "); err != nil {
				return err
			}
		}
		if err := column.toString(buf); err != nil {
			return err
		}
	}

	indexes := sortedKeyIndexes(t.indexes)

	for i, index := range indexes {
		if i+len(columns) > 0 {
			if _, err := buf.WriteString(",\n  "); err != nil {
				return err
			}
		}
		if err := index.toString(buf); err != nil {
			return err
		}
	}

	for i, fk := range foreignKeys {
		if i+len(columns)+len(indexes) > 0 {
			if _, err := buf.WriteString(",\n  "); err != nil {
				return err
			}
		}
		if err := fk.toString(buf); err != nil {
			return err
		}
	}

	if _, err := buf.WriteString("\n);\n"); err != nil {
		return err
	}
	return nil
}

func newTableState(id int, name string) *tableState {
	return &tableState{
		id:      id,
		name:    name,
		columns: make(map[string]*columnState),
	}
}

// sortedColumns returns a table's columns in declaration order. Both the
// CREATE TABLE body and the trailing COMMENT ON COLUMN statements go through
// it, so the two agree and neither depends on map iteration order.
func sortedColumns(columns map[string]*columnState) []*columnState {
	sorted := make([]*columnState, 0, len(columns))
	for _, column := range columns {
		sorted = append(sorted, column)
	}
	slices.SortFunc(sorted, func(a, b *columnState) int {
		if a.id < b.id {
			return -1
		}
		if a.id > b.id {
			return 1
		}
		return 0
	})
	return sorted
}

func convertToTableState(id int, schemaName string, table *metadatapb.TableMetadata) *tableState {
	state := newTableState(id, table.Name)
	state.schema = schemaName
	state.owner = table.Owner
	state.comment = table.Comment
	for _, cc := range table.CheckConstraints {
		if cc != nil {
			state.hasCheckConstraints = true
			break
		}
	}
	switch {
	case table.CreateOptions != "":
		state.unsupportedReason = fmt.Sprintf("create options %q cannot be restored by Redshift DDL renderer", table.CreateOptions)
	case len(table.SortingKeys) > 0:
		state.unsupportedReason = "sorting keys cannot be restored by Redshift DDL renderer"
	case table.ShardingInfo != "":
		state.unsupportedReason = "sharding info cannot be restored by Redshift DDL renderer"
	case len(table.Partitions) > 0:
		state.unsupportedReason = "partitions cannot be restored because Redshift does not support declarative partitions"
	case len(table.Triggers) > 0:
		state.unsupportedReason = "triggers cannot be restored because Redshift does not support triggers"
	case len(table.Rules) > 0:
		state.unsupportedReason = "rules cannot be restored because Redshift does not support rules"
	case len(table.ExcludeConstraints) > 0:
		state.unsupportedReason = "exclude constraints cannot be restored because Redshift does not support EXCLUDE constraints"
	}
	for i, column := range table.Columns {
		if column == nil {
			continue
		}
		if len(column.Options) > 0 && state.unsupportedReason == "" {
			state.unsupportedReason = fmt.Sprintf("column %s options cannot be restored by Redshift DDL renderer", quoteDefinitionIdentifier(column.Name))
		}
		if _, exists := state.columns[column.Name]; exists && state.duplicateColumn == "" {
			state.duplicateColumn = column.Name
		}
		colState := convertToColumnState(i, column)
		state.columns[column.Name] = colState
		if _, ok := colState.defaultValue.(*defaultValueExpression); ok {
			state.defaultExprs = append(state.defaultExprs, column.Default)
		}
	}
	for i, index := range table.Indexes {
		if index == nil {
			continue
		}
		state.indexes = append(state.indexes, convertToIndexState(i, index))
	}
	for i, fk := range table.ForeignKeys {
		if fk == nil {
			continue
		}
		state.foreignKeys = append(state.foreignKeys, convertToForeignKeyState(i, schemaName, fk))
	}
	return state
}

type foreignKeyState struct {
	id                int
	name              string
	columns           []string
	referencedSchema  string
	referencedTable   string
	referencedColumns []string
	onDelete          string
	onUpdate          string
	comment           string
}

func isSupportedRedshiftFKAction(action string) bool {
	trimmed := strings.TrimSpace(strings.ToUpper(action))
	return trimmed == "" || trimmed == "NO ACTION"
}

func convertToForeignKeyState(id int, tableSchema string, foreignKey *metadatapb.ForeignKeyMetadata) *foreignKeyState {
	refSchema := foreignKey.ReferencedSchema
	if refSchema == "" {
		refSchema = tableSchema
	}
	return &foreignKeyState{
		id:                id,
		name:              foreignKey.Name,
		columns:           foreignKey.Columns,
		referencedSchema:  refSchema,
		referencedTable:   foreignKey.ReferencedTable,
		referencedColumns: foreignKey.ReferencedColumns,
		onDelete:          foreignKey.OnDelete,
		onUpdate:          foreignKey.OnUpdate,
		comment:           foreignKey.Comment,
	}
}

func (f *foreignKeyState) toString(buf *strings.Builder) error {
	if len(f.columns) == 0 || len(f.columns) != len(f.referencedColumns) {
		return fmt.Errorf("foreign key %s has %d columns and %d referenced columns", f.name, len(f.columns), len(f.referencedColumns))
	}
	if !isSupportedRedshiftFKAction(f.onDelete) {
		return fmt.Errorf("foreign key %s has unsupported ON DELETE action %q", f.name, f.onDelete)
	}
	if !isSupportedRedshiftFKAction(f.onUpdate) {
		return fmt.Errorf("foreign key %s has unsupported ON UPDATE action %q", f.name, f.onUpdate)
	}
	n := len(f.columns)
	if f.name != "" {
		if _, err := fmt.Fprintf(buf, "CONSTRAINT %s ", quoteDefinitionIdentifier(f.name)); err != nil {
			return err
		}
	}
	if _, err := buf.WriteString("FOREIGN KEY ("); err != nil {
		return err
	}
	for i := 0; i < n; i++ {
		if i > 0 {
			if _, err := buf.WriteString(", "); err != nil {
				return err
			}
		}
		if _, err := buf.WriteString(quoteDefinitionIdentifier(f.columns[i])); err != nil {
			return err
		}
	}
	// The sync strips the qualifier off ReferencedTable and keeps it in
	// ReferencedSchema, so an unqualified target here would bind to whatever
	// the search path resolves rather than the table the constraint names.
	referenced := qualifyName(f.referencedSchema, f.referencedTable)
	if _, err := fmt.Fprintf(buf, ") REFERENCES %s(", referenced); err != nil {
		return err
	}
	for i := 0; i < n; i++ {
		if i > 0 {
			if _, err := buf.WriteString(", "); err != nil {
				return err
			}
		}
		if _, err := buf.WriteString(quoteDefinitionIdentifier(f.referencedColumns[i])); err != nil {
			return err
		}
	}
	_, err := buf.WriteString(")")
	return err
}

type indexState struct {
	id      int
	name    string
	keys    []string
	lengths []int64
	primary bool
	unique  bool
	tp      string
	comment string
}

func convertToIndexState(id int, index *metadatapb.IndexMetadata) *indexState {
	return &indexState{
		id:      id,
		name:    index.Name,
		keys:    index.Expressions,
		lengths: index.KeyLength,
		primary: index.Primary,
		unique:  index.Unique,
		tp:      index.Type,
		comment: index.Comment,
	}
}

func sortedKeyIndexes(all []*indexState) []*indexState {
	indexes := make([]*indexState, 0, len(all))
	for _, index := range all {
		if index.primary || index.unique {
			indexes = append(indexes, index)
		}
	}
	slices.SortFunc(indexes, func(a, b *indexState) int {
		if a.primary && !b.primary {
			return -1
		}
		if !a.primary && b.primary {
			return 1
		}
		if cmp := strings.Compare(a.name, b.name); cmp != 0 {
			return cmp
		}
		if a.id < b.id {
			return -1
		}
		if a.id > b.id {
			return 1
		}
		return 0
	})
	return indexes
}

func (i *indexState) toString(buf *strings.Builder) error {
	var keyword string
	switch {
	case i.primary:
		keyword = "PRIMARY KEY ("
	case i.unique:
		keyword = "UNIQUE ("
	default:
		return nil
	}
	if len(i.keys) == 0 {
		return fmt.Errorf("index %s has no columns", i.name)
	}
	if i.name != "" {
		if _, err := fmt.Fprintf(buf, "CONSTRAINT %s ", quoteDefinitionIdentifier(i.name)); err != nil {
			return err
		}
	}
	if _, err := buf.WriteString(keyword); err != nil {
		return err
	}
	for idx, key := range i.keys {
		if idx > 0 {
			if _, err := buf.WriteString(", "); err != nil {
				return err
			}
		}
		// Written verbatim: these come from pg_get_indexdef(.., true), which
		// already quotes what needs it and can be an expression rather than
		// a column, so quoting again would name something else.
		if _, err := buf.WriteString(key); err != nil {
			return err
		}
	}
	_, err := buf.WriteString(")")
	return err
}

type defaultValue interface {
	toString() string
}

type defaultValueNull struct{}

func (*defaultValueNull) toString() string {
	return "NULL"
}

type defaultValueString struct {
	value string
}

func (d *defaultValueString) toString() string {
	return fmt.Sprintf("'%s'", escapeStringLiteral(d.value))
}

type defaultValueExpression struct {
	value string
}

func (d *defaultValueExpression) toString() string {
	return d.value
}

type columnState struct {
	id                 int
	name               string
	tp                 string
	collationSchema    string
	collation          string
	defaultValue       defaultValue
	comment            string
	nullable           bool
	isIdentity         bool
	identityGeneration metadatapb.ColumnMetadata_IdentityGeneration
	identitySeed       int64
	identityIncrement  int64
}

func (c *columnState) toString(buf *strings.Builder) error {
	if _, err := fmt.Fprintf(buf, "%s ", quoteDefinitionIdentifier(c.name)); err != nil {
		return err
	}
	if _, err := buf.WriteString(c.tp); err != nil {
		return err
	}
	if c.collation != "" {
		if _, err := fmt.Fprintf(buf, " COLLATE %s", formatRedshiftCollation(c.collationSchema, c.collation)); err != nil {
			return err
		}
	}
	if c.isIdentity {
		seed, step := c.identitySeed, c.identityIncrement
		if step == 0 {
			seed, step = 1, 1
		}
		if c.identityGeneration == metadatapb.ColumnMetadata_BY_DEFAULT {
			if _, err := fmt.Fprintf(buf, " GENERATED BY DEFAULT AS IDENTITY(%d, %d)", seed, step); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintf(buf, " IDENTITY(%d, %d)", seed, step); err != nil {
				return err
			}
		}
	}
	if !c.nullable {
		if _, err := buf.WriteString(" NOT NULL"); err != nil {
			return err
		}
	}
	if !c.isIdentity && c.defaultValue != nil {
		if _, err := fmt.Fprintf(buf, " DEFAULT %s", c.defaultValue.toString()); err != nil {
			return err
		}
	}
	return nil
}

func convertToColumnState(id int, column *metadatapb.ColumnMetadata) *columnState {
	result := &columnState{
		id:                 id,
		name:               column.Name,
		tp:                 column.Type,
		collationSchema:    column.CollationSchema,
		collation:          column.Collation,
		nullable:           column.Nullable,
		comment:            column.Comment,
		isIdentity:         column.IsIdentity,
		identityGeneration: column.IdentityGeneration,
		identitySeed:       column.IdentitySeed,
		identityIncrement:  column.IdentityIncrement,
	}
	if strings.EqualFold(strings.TrimSpace(column.Default), "NULL") {
		result.defaultValue = &defaultValueNull{}
	} else if column.Default != "" {
		if isExpression(column.Default) || isQuotedSQLLiteral(column.Default) || isNumericOrBoolLiteral(column.Default) {
			result.defaultValue = &defaultValueExpression{value: column.Default}
		} else {
			result.defaultValue = &defaultValueString{value: column.Default}
		}
	}
	return result
}

// formatRedshiftCollation emits Redshift's unquoted collation keywords
// (CASE_SENSITIVE, CS, CASE_INSENSITIVE, CI) verbatim and quotes other
// collation names as identifiers.
func formatRedshiftCollation(schema, name string) string {
	if schema == "" {
		switch strings.ToUpper(name) {
		case "CASE_SENSITIVE", "CS", "CASE_INSENSITIVE", "CI":
			return strings.ToUpper(name)
		}
	}
	return qualifyName(schema, name)
}

func isExpression(value string) bool {
	trimmed := strings.TrimSpace(value)
	switch strings.ToUpper(trimmed) {
	case "SYSDATE", "CURRENT_DATE", "CURRENT_TIME", "CURRENT_TIMESTAMP",
		"CURRENT_USER", "CURRENT_SCHEMA", "LOCALTIME", "LOCALTIMESTAMP",
		"SESSION_USER", "USER":
		return true
	}
	stmts, err := pgparser.Parse("SELECT " + trimmed)
	if err != nil || stmts == nil || len(stmts.Items) != 1 {
		return false
	}
	stmt := stmts.Items[0]
	if raw, ok := stmt.(*nodes.RawStmt); ok {
		stmt = raw.Stmt
	}
	sel, ok := stmt.(*nodes.SelectStmt)
	if !ok || sel.TargetList == nil || len(sel.TargetList.Items) != 1 {
		return false
	}
	rt, ok := sel.TargetList.Items[0].(*nodes.ResTarget)
	if !ok || rt.Name != "" {
		return false
	}
	if rt.Val == nil {
		return false
	}
	switch rt.Val.(type) {
	case *nodes.ColumnRef:
		return false
	case *nodes.A_Const:
		return true
	}
	// Column defaults cannot reference columns; rejecting ColumnRef keeps
	// raw strings like "pending (old)" from matching as FuncCall(ColumnRef).
	hasColRef := false
	nodes.Inspect(rt.Val, func(n nodes.Node) bool {
		if _, ok := n.(*nodes.ColumnRef); ok {
			hasColRef = true
			return false
		}
		return true
	})
	return !hasColRef
}

// isQuotedSQLLiteral reports whether s is already a single-quoted SQL string
// literal (such as 'active' or 'it”s ok'), so column defaults copied from
// information_schema.columns.column_default are preserved verbatim.
func isQuotedSQLLiteral(s string) bool {
	if len(s) < 2 || s[0] != '\'' || s[len(s)-1] != '\'' {
		return false
	}
	for i := 1; i < len(s)-1; i++ {
		switch s[i] {
		case '\\':
			if i+1 >= len(s)-1 {
				return false
			}
			i++
		case '\'':
			if i+1 >= len(s)-1 || s[i+1] != '\'' {
				return false
			}
			i++
		}
	}
	return true
}

func isNumericOrBoolLiteral(s string) bool {
	if strings.EqualFold(s, "true") || strings.EqualFold(s, "false") {
		return true
	}
	if s == "" {
		return false
	}
	i := 0
	if s[i] == '+' || s[i] == '-' {
		i++
	}
	digits := false
	dot := false
	for ; i < len(s); i++ {
		switch {
		case s[i] >= '0' && s[i] <= '9':
			digits = true
		case s[i] == '.' && !dot:
			dot = true
		default:
			return false
		}
	}
	return digits
}

type viewState struct {
	id         int
	schema     string
	name       string
	definition string
	comment    string
	columns    map[string]*columnState
}

func convertToViewState(id int, schemaName string, view *metadatapb.ViewMetadata) *viewState {
	cols := make(map[string]*columnState, len(view.Columns))
	for i, col := range view.Columns {
		if col == nil {
			continue
		}
		cols[col.Name] = convertToColumnState(i, col)
	}
	return &viewState{
		id:         id,
		schema:     schemaName,
		name:       view.Name,
		definition: view.Definition,
		comment:    view.Comment,
		columns:    cols,
	}
}

func convertToMaterializedViewState(id int, schemaName string, mv *metadatapb.MaterializedViewMetadata) *viewState {
	return &viewState{
		id:         id,
		schema:     schemaName,
		name:       mv.Name,
		definition: mv.Definition,
		comment:    mv.Comment,
	}
}

func (v *viewState) toString(buf io.StringWriter) error {
	trimmed := strings.TrimSpace(v.definition)
	if trimmed == "" {
		return fmt.Errorf("view %s has no definition to restore", v.qualifiedName())
	}
	var stmt string
	if isCompleteCreateView(trimmed) {
		stmt = trimmed
	} else {
		stmt = fmt.Sprintf("CREATE OR REPLACE VIEW %s AS %s", v.qualifiedName(), v.definition)
	}
	stmt = appendTerminator(stmt) + "\n"
	_, err := buf.WriteString(stmt)
	return err
}

func (v *viewState) toMaterializedString(buf io.StringWriter) error {
	trimmed := strings.TrimSpace(v.definition)
	if trimmed == "" {
		return fmt.Errorf("materialized view %s has no definition to restore", v.qualifiedName())
	}
	var stmt string
	if isCompleteCreateMaterializedView(trimmed) {
		stmt = trimmed
	} else {
		stmt = fmt.Sprintf("CREATE MATERIALIZED VIEW %s AS %s", v.qualifiedName(), v.definition)
	}
	stmt = appendTerminator(stmt) + "\n"
	_, err := buf.WriteString(stmt)
	return err
}

// isCompleteCreateView reports whether definition already starts with a full
// CREATE [OR REPLACE] VIEW statement (as SHOW VIEW returns for late-binding
// views with WITH NO SCHEMA BINDING).
func isCompleteCreateView(definition string) bool {
	return parseViewStmt(definition) != nil
}

// isLateBindingView reports whether definition is a CREATE VIEW ... WITH NO
// SCHEMA BINDING statement, whose columns Redshift does not allow commenting.
func isLateBindingView(definition string) bool {
	vs := parseViewStmt(strings.TrimSpace(definition))
	if vs == nil || vs.Options == nil {
		return false
	}
	for _, item := range vs.Options.Items {
		if de, ok := item.(*nodes.DefElem); ok && de.Defname == "no_schema_binding" {
			return true
		}
	}
	return false
}

func parseViewStmt(definition string) *nodes.ViewStmt {
	stmts, err := pgparser.Parse(definition)
	if err != nil || stmts == nil || len(stmts.Items) != 1 {
		return nil
	}
	stmt := stmts.Items[0]
	if raw, ok := stmt.(*nodes.RawStmt); ok {
		stmt = raw.Stmt
	}
	vs, _ := stmt.(*nodes.ViewStmt)
	return vs
}

func isCompleteCreateMaterializedView(definition string) bool {
	stmts, err := pgparser.Parse(definition)
	if err != nil || stmts == nil || len(stmts.Items) != 1 {
		return false
	}
	stmt := stmts.Items[0]
	if raw, ok := stmt.(*nodes.RawStmt); ok {
		stmt = raw.Stmt
	}
	ctas, ok := stmt.(*nodes.CreateTableAsStmt)
	return ok && ctas.Objtype == nodes.OBJECT_MATVIEW
}
