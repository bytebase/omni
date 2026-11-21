package catalog

import (
	"container/heap"
	"fmt"
	"io"
	"math"
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
	if to.Collation != "" {
		return "", fmt.Errorf("database %q collation %q cannot be restored by Redshift DDL renderer", to.Name, to.Collation)
	}
	if to.CharacterSet != "" && !strings.EqualFold(to.CharacterSet, "UNICODE") && !strings.EqualFold(to.CharacterSet, "UTF8") && !strings.EqualFold(to.CharacterSet, "UTF-8") {
		return "", fmt.Errorf("database %q character set %q cannot be restored by Redshift DDL renderer", to.Name, to.CharacterSet)
	}
	if to.SearchPath != "" || to.ServiceName != "" || to.DefaultSchema != "" ||
		len(to.Extensions) > 0 ||
		slices.ContainsFunc(to.EventTriggers, func(et *metadatapb.EventTriggerMetadata) bool { return et != nil && !et.SkipDump }) ||
		len(to.ForeignServers) > 0 || len(to.ForeignDataWrappers) > 0 || len(to.PublicSynonyms) > 0 || len(to.LinkedDatabases) > 0 {
		return "", fmt.Errorf("database %q contains objects unsupported by Redshift DDL renderer", to.Name)
	}
	caseSensitive := needsCaseSensitiveIdentifiers(to)
	seenSchemas := make(map[string]struct{}, len(to.Schemas))
	for _, schema := range to.Schemas {
		if schema == nil || schema.SkipDump {
			continue
		}
		if schema.Name == "" {
			return "", fmt.Errorf("database %s has schema with empty name", quoteDefinitionIdentifier(to.Name))
		}
		key := identKey(schema.Name, caseSensitive)
		if _, exists := seenSchemas[key]; exists {
			return "", fmt.Errorf("database %s has duplicate schema %s", quoteDefinitionIdentifier(to.Name), quoteDefinitionIdentifier(schema.Name))
		}
		seenSchemas[key] = struct{}{}
	}
	toState := convertToDatabaseState(to, caseSensitive)

	var sb strings.Builder

	if caseSensitive {
		if _, err := sb.WriteString("SET enable_case_sensitive_identifier TO true;\n"); err != nil {
			return "", err
		}
	}
	if err := writeSchemas(&sb, to); err != nil {
		return "", err
	}
	seenRoutineSigs := make(map[string]struct{})
	if err := writeOrderedObjects(&sb, to, toState, caseSensitive, seenRoutineSigs); err != nil {
		return "", err
	}
	if err := writeProcedures(&sb, to, caseSensitive, seenRoutineSigs); err != nil {
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
	caseSensitive := hasUpperASCII(schemaName) || tableHasUpperIdentifier(table)
	ts := convertToTableState(0, schemaName, table)
	if err := validateTableDefaultExprs(ts, nil); err != nil {
		return "", err
	}
	var sb strings.Builder
	if caseSensitive {
		if _, err := sb.WriteString("SET enable_case_sensitive_identifier TO true;\n\n"); err != nil {
			return "", err
		}
	}
	if err := writeTableWithOptions(&sb, ts, nil, true, true, caseSensitive); err != nil {
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
	caseSensitive := hasUpperASCII(schemaName) || viewHasUpperIdentifier(view)
	var sb strings.Builder
	if caseSensitive {
		if _, err := sb.WriteString("SET enable_case_sensitive_identifier TO true;\n\n"); err != nil {
			return "", err
		}
	}
	if err := writeView(&sb, convertToViewState(0, schemaName, view, caseSensitive)); err != nil {
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

func depsHaveUpperIdentifier(deps []*metadatapb.DependencyColumn) bool {
	for _, dep := range deps {
		if dep != nil && (hasUpperASCII(dep.Schema) || hasUpperASCII(dep.Table)) {
			return true
		}
	}
	return false
}

func depTablesHaveUpperIdentifier(deps []*metadatapb.DependencyTable) bool {
	for _, dep := range deps {
		if dep != nil && (hasUpperASCII(dep.Schema) || hasUpperASCII(dep.Table)) {
			return true
		}
	}
	return false
}

func viewHasUpperIdentifier(view *metadatapb.ViewMetadata) bool {
	if hasUpperASCII(view.Name) || hasUpperQuotedIdent(view.Definition) || depsHaveUpperIdentifier(view.DependencyColumns) {
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
			if mv != nil && !mv.SkipDump && (hasUpperASCII(mv.Name) || hasUpperQuotedIdent(mv.Definition) || depsHaveUpperIdentifier(mv.DependencyColumns)) {
				return true
			}
		}
		for _, fn := range schema.Functions {
			if fn != nil && !fn.SkipDump && (hasUpperASCII(fn.Name) || hasUpperQuotedIdent(fn.Definition) || depTablesHaveUpperIdentifier(fn.DependencyTables)) {
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
	parsed     bool
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
	if call.schema == "" && !slices.Contains(call.argTypes, "") && hasMatchingBuiltinProc(call.name, call.argTypes) {
		// Redshift searches pg_catalog before search_path for unqualified calls,
		// but only suppress catalog lookup when all argument types are known.
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

func redshiftTypeCompatScore(argType, paramType string) (int, bool) {
	if argType == "" || paramType == "" {
		return 0, true
	}
	if argType == "unknown_str" {
		switch paramType {
		case "varchar":
			return 2, true
		case "bpchar":
			return 1, true
		default:
			return 0, false
		}
	}
	if paramType == argType {
		return 2, true
	}
	if argType == "bpchar" && paramType == "varchar" {
		return 1, true
	}
	if argRank, paramRank := numericTypeRank(argType), numericTypeRank(paramType); argRank > 0 && paramRank >= argRank {
		return 1, true
	}
	return 0, false
}

func resolveFunctionOverload(fns []functionEntry, candidates []int, argTypes []string) []int {
	if len(candidates) == 0 {
		return nil
	}
	var arityMatches []int
	for _, idx := range candidates {
		if !fns[idx].parsed || len(fns[idx].params) == len(argTypes) {
			arityMatches = append(arityMatches, idx)
		}
	}
	if len(arityMatches) == 0 {
		return nil
	}
	if len(argTypes) == 0 {
		return arityMatches
	}
	bestScore := -1
	var best []int
	for _, idx := range arityMatches {
		if !fns[idx].parsed {
			if bestScore < 0 {
				bestScore = 0
				best = []int{idx}
			} else if bestScore == 0 {
				best = append(best, idx)
			}
			continue
		}
		score := 0
		compatible := true
		for i, argType := range argTypes {
			s, ok := redshiftTypeCompatScore(argType, fns[idx].params[i].typeName)
			if !ok {
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
	case "timestamp with time zone", "timestamptz":
		return "timestamptz"
	case "timestamp without time zone", "timestamp":
		return "timestamp"
	case "time with time zone", "timetz":
		return "timetz"
	case "time without time zone", "time":
		return "time"
	case "boolean", "bool":
		return "bool"
	case "char", "character", "bpchar", "nchar":
		return "bpchar"
	case "varchar", "character varying", "text", "nvarchar":
		return "varchar"
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
	var base string
	if len(parts) == 2 && parts[0] == "pg_catalog" {
		base = normalizeTypeAlias(parts[1])
	} else if len(parts) > 1 {
		base = strings.Join(parts, ".")
	} else {
		base = normalizeTypeAlias(parts[0])
	}
	if tn.ArrayBounds != nil && len(tn.ArrayBounds.Items) > 0 {
		return base + "[]"
	}
	return base
}

func parseFunctionSignature(definition string) (declName []string, params []functionParam, returnType string, parsed bool) {
	stmts, err := pgparser.Parse(definition)
	if err != nil || stmts == nil || len(stmts.Items) != 1 {
		return nil, nil, "", false
	}
	stmt := stmts.Items[0]
	if raw, ok := stmt.(*nodes.RawStmt); ok {
		stmt = raw.Stmt
	}
	cf, ok := stmt.(*nodes.CreateFunctionStmt)
	if !ok {
		return nil, nil, "", false
	}
	declName = extractNodeStringParts(cf.Funcname)
	returnType = normalizeSQLTypeName(cf.ReturnType)
	if cf.Parameters == nil {
		return declName, nil, returnType, true
	}
	for _, item := range cf.Parameters.Items {
		fp, ok := item.(*nodes.FunctionParameter)
		if !ok || fp.Mode == nodes.FUNC_PARAM_OUT || fp.Mode == nodes.FUNC_PARAM_TABLE {
			continue
		}
		params = append(params, functionParam{
			name:     fp.Name,
			typeName: normalizeSQLTypeName(fp.ArgType),
		})
	}
	return declName, params, returnType, true
}

func detectCreateRoutineKind(definition string) string {
	lex := pgparser.NewLexer(definition)
	sawCreate := false
	for i := 0; i < 6; i++ {
		tok := lex.NextToken()
		if tok.Type == 0 {
			return ""
		}
		word := strings.ToUpper(tok.Str)
		if !sawCreate {
			if word == "CREATE" {
				sawCreate = true
				continue
			}
			return ""
		}
		switch word {
		case "OR", "REPLACE":
			continue
		case "PROCEDURE":
			return "procedure"
		case "FUNCTION":
			return "function"
		default:
			return ""
		}
	}
	return ""
}

func parseSignatureParamTypes(signature string) ([]string, bool) {
	sig := strings.TrimSpace(signature)
	if sig == "" {
		return nil, false
	}
	if open := strings.IndexByte(sig, '('); open >= 0 {
		close := strings.LastIndexByte(sig, ')')
		if close <= open {
			return nil, false
		}
		sig = strings.TrimSpace(sig[open+1 : close])
	}
	if sig == "" {
		return []string{}, true
	}
	stmts, err := pgparser.Parse("CREATE FUNCTION _sig(" + sig + ") RETURNS int IMMUTABLE AS $$ SELECT 1 $$ LANGUAGE SQL")
	if err != nil || stmts == nil || len(stmts.Items) != 1 {
		return nil, false
	}
	stmt := stmts.Items[0]
	if raw, ok := stmt.(*nodes.RawStmt); ok {
		stmt = raw.Stmt
	}
	cf, ok := stmt.(*nodes.CreateFunctionStmt)
	if !ok || cf.Parameters == nil {
		return nil, false
	}
	out := make([]string, 0, len(cf.Parameters.Items))
	for _, item := range cf.Parameters.Items {
		fp, ok := item.(*nodes.FunctionParameter)
		if !ok || fp.Mode == nodes.FUNC_PARAM_OUT || fp.Mode == nodes.FUNC_PARAM_TABLE {
			continue
		}
		out = append(out, normalizeSQLTypeName(fp.ArgType))
	}
	return out, true
}

func validateRoutineIdentity(kind, schemaName, objName, signature, def string, expectProc, caseSensitive bool) ([]functionParam, string, bool, error) {
	detectedKind := detectCreateRoutineKind(def)
	declName, params, retType, parsed := parseFunctionSignature(def)
	if expectProc && (!parsed || detectedKind != "procedure") {
		return nil, "", false, fmt.Errorf("%s %s definition must be a CREATE PROCEDURE statement", kind, qualifyName(schemaName, objName))
	}
	if !expectProc && detectedKind != "function" {
		return nil, "", false, fmt.Errorf("%s %s definition must be a CREATE FUNCTION statement", kind, qualifyName(schemaName, objName))
	}
	if len(declName) > 0 {
		mismatch := func() error {
			return fmt.Errorf("%s %s definition name %q does not match metadata identity", kind, qualifyName(schemaName, objName), strings.Join(declName, "."))
		}
		switch len(declName) {
		case 1:
			if !identEqual(declName[0], objName, caseSensitive) {
				return nil, "", false, mismatch()
			}
		case 2:
			if !identEqual(declName[0], schemaName, caseSensitive) || !identEqual(declName[1], objName, caseSensitive) {
				return nil, "", false, mismatch()
			}
		default:
			return nil, "", false, mismatch()
		}
	}
	if parsed && signature != "" {
		if sigTypes, ok := parseSignatureParamTypes(signature); ok {
			if len(sigTypes) != len(params) {
				return nil, "", false, fmt.Errorf("%s %s definition parameter count (%d) does not match signature %q", kind, qualifyName(schemaName, objName), len(params), signature)
			}
			for i, st := range sigTypes {
				if st != "" && params[i].typeName != "" && st != params[i].typeName {
					return nil, "", false, fmt.Errorf("%s %s parameter %d type %q does not match signature type %q", kind, qualifyName(schemaName, objName), i+1, params[i].typeName, st)
				}
			}
		}
	}
	return params, retType, parsed, nil
}

func identKey(name string, caseSensitive bool) string {
	if caseSensitive {
		return name
	}
	return strings.ToLower(name)
}

func identEqual(a, b string, caseSensitive bool) bool {
	if caseSensitive {
		return a == b
	}
	return strings.EqualFold(a, b)
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
			if p.Kind != 'f' && p.Kind != 'a' && p.Kind != 'w' {
				continue
			}
			key := strings.ToLower(p.Name)
			builtinProcsByName[key] = append(builtinProcsByName[key], p)
		}
	})
}

func matchBuiltinProcArgs(p BuiltinProc, argTypes []string) (int, bool) {
	isVariadic := p.Variadic != 0
	minArgs := int(p.NArgs) - int(p.NArgDefaults)
	if minArgs < 0 {
		minArgs = 0
	}
	if !isVariadic {
		if len(argTypes) < minArgs || len(argTypes) > int(p.NArgs) || len(p.ArgTypes) != int(p.NArgs) {
			return 0, false
		}
	} else if len(argTypes) < minArgs || len(p.ArgTypes) == 0 {
		return 0, false
	}
	score := 0
	for i, argType := range argTypes {
		if argType == "" {
			continue
		}
		var paramOID uint32
		if isVariadic && i >= len(p.ArgTypes)-1 {
			paramOID = p.Variadic
		} else {
			paramOID = p.ArgTypes[i]
		}
		paramType := builtinTypeByOID[paramOID]
		if paramType == "" || strings.HasPrefix(paramType, "any") {
			continue
		}
		s, ok := redshiftTypeCompatScore(argType, paramType)
		if !ok {
			return 0, false
		}
		score += s
	}
	return score, true
}

func hasMatchingBuiltinProc(name string, argTypes []string) bool {
	initBuiltinCatalog()
	for _, p := range builtinProcsByName[strings.ToLower(name)] {
		if _, ok := matchBuiltinProcArgs(p, argTypes); ok {
			return true
		}
	}
	return false
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
		score, ok := matchBuiltinProcArgs(p, argTypes)
		if !ok {
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
		switch v := n.Val.(type) {
		case *nodes.Integer:
			if v.Ival >= math.MinInt32 && v.Ival <= math.MaxInt32 {
				return "int4"
			}
			return "int8"
		case *nodes.Float:
			return "numeric"
		case *nodes.String:
			return "unknown_str"
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
			return "varchar"
		case "+", "-", "*", "/", "%", "^":
			if n.Lexpr == nil {
				if op == "-" {
					if ac, ok := n.Rexpr.(*nodes.A_Const); ok && !ac.Isnull {
						switch v := ac.Val.(type) {
						case *nodes.Integer:
							if -v.Ival >= math.MinInt32 && -v.Ival <= math.MaxInt32 {
								return "int4"
							}
							return "int8"
						case *nodes.Float:
							if v.Fval == "9223372036854775808" {
								return "int8"
							}
						}
					}
				}
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
	return scanCallTokens(definition, true)
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
	return scanCallTokens(body, false)
}

func scanCallTokens(body string, skipHeader bool) []functionCall {
	lex := pgparser.NewLexer(body)
	var toks []pgparser.Token
	for {
		tok := lex.NextToken()
		if tok.Type == 0 {
			break
		}
		toks = append(toks, tok)
	}
	startIdx := 0
	if skipHeader {
		// Unparsed DDL often wraps its SQL body in `AS $$...$$`, which the lexer
		// returns as a single string-literal token; extract and scan that body directly.
		if bodyTok, ok := findFunctionASBodyToken(body, toks); ok {
			if parsedCalls := collectSQLFunctionCalls(bodyTok.Str, "", nil); len(parsedCalls) > 0 {
				return parsedCalls
			}
			return scanCallTokens(bodyTok.Str, false)
		}
		depth := 0
		for k := 0; k < len(toks); k++ {
			if toks[k].Type == '(' {
				depth++
			} else if toks[k].Type == ')' && depth > 0 {
				depth--
				if depth == 0 {
					startIdx = k + 1
					break
				}
			}
		}
		toks = toks[startIdx:]
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

func writeWithSchemaSearchPath(w io.StringWriter, schemaName, body string) error {
	if schemaName != "" {
		path := quoteDefinitionIdentifier(schemaName)
		if schemaName != "public" {
			path += ", public"
		}
		if _, err := w.WriteString(fmt.Sprintf("SET search_path TO %s;\n", path)); err != nil {
			return err
		}
	}
	if _, err := w.WriteString(body); err != nil {
		return err
	}
	if schemaName != "" {
		if _, err := w.WriteString("RESET search_path;\n"); err != nil {
			return err
		}
	}
	return nil
}

func writeProcedures(w io.StringWriter, to *metadatapb.DatabaseSchemaMetadata, caseSensitive bool, seenRoutineSigs map[string]struct{}) error {
	type validatedProc struct {
		schema string
		name   string
		def    string
	}
	var procs []validatedProc
	if seenRoutineSigs == nil {
		seenRoutineSigs = make(map[string]struct{})
	}
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
			params, _, parsed, err := validateRoutineIdentity("procedure", schema.Name, proc.Name, proc.Signature, def, true, caseSensitive)
			if err != nil {
				return err
			}
			sigKey := redshiftFunctionIdentityKey(schema.Name, proc.Name, proc.Signature, params, parsed, caseSensitive)
			if _, exists := seenRoutineSigs[sigKey]; exists {
				return fmt.Errorf("duplicate procedure or routine %s in schema %s cannot be restored", proc.Name, quoteDefinitionIdentifier(schema.Name))
			}
			seenRoutineSigs[sigKey] = struct{}{}
			procs = append(procs, validatedProc{schema: schema.Name, name: proc.Name, def: def})
		}
	}
	for _, proc := range procs {
		if _, err := w.WriteString(getProcedureAnnouncement(proc.name)); err != nil {
			return err
		}
		if err := writeWithSchemaSearchPath(w, proc.schema, proc.def+"\n"); err != nil {
			return err
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

func validationOnlyFunctionEntry(schemaName string, fn *metadatapb.FunctionMetadata) functionEntry {
	var params []functionParam
	var retType string
	var parsed bool
	if def := strings.TrimSpace(fn.Definition); def != "" {
		_, params, retType, parsed = parseFunctionSignature(def)
	}
	if !parsed && fn.Signature != "" {
		if sigTypes, ok := parseSignatureParamTypes(fn.Signature); ok {
			params = make([]functionParam, len(sigTypes))
			for i, t := range sigTypes {
				params[i] = functionParam{typeName: t}
			}
			parsed = true
		}
	}
	return functionEntry{
		schema:     schemaName,
		name:       fn.Name,
		params:     params,
		returnType: retType,
		parsed:     parsed,
	}
}

// writeOrderedObjects emits functions, tables, views, and materialized views
// in one topological order so cross-kind dependencies (such as view-dependent
// UDFs) always follow their prerequisites.
func writeOrderedObjects(w io.StringWriter, to *metadatapb.DatabaseSchemaMetadata, state *databaseState, caseSensitive bool, seenFnSigs map[string]struct{}) error {
	var items []orderedObject
	var fnEntries []functionEntry
	var fnNodeIndices []int
	if seenFnSigs == nil {
		seenFnSigs = make(map[string]struct{})
	}

	for _, schema := range to.Schemas {
		if schema == nil {
			continue
		}
		if schema.SkipDump {
			for _, fn := range schema.Functions {
				if fn != nil {
					fnEntries = append(fnEntries, validationOnlyFunctionEntry(schema.Name, fn))
					fnNodeIndices = append(fnNodeIndices, -1)
				}
			}
			continue
		}
		seenRels := make(map[string]string, len(schema.Tables)+len(schema.Views)+len(schema.MaterializedViews))
		checkRel := func(kind, name string) error {
			key := identKey(name, caseSensitive)
			if prev, exists := seenRels[key]; exists {
				return fmt.Errorf("schema %s has duplicate relation %s (%s and %s)", quoteDefinitionIdentifier(schema.Name), quoteDefinitionIdentifier(name), prev, kind)
			}
			seenRels[key] = kind
			return nil
		}
		for _, table := range schema.Tables {
			if table == nil || table.SkipDump {
				continue
			}
			if err := checkRel("table", table.Name); err != nil {
				return err
			}
		}
		for _, view := range schema.Views {
			if view == nil || view.SkipDump {
				continue
			}
			if err := checkRel("view", view.Name); err != nil {
				return err
			}
		}
		for _, mv := range schema.MaterializedViews {
			if mv == nil || mv.SkipDump {
				continue
			}
			if err := checkRel("materialized view", mv.Name); err != nil {
				return err
			}
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
		if len(schema.Streams) > 0 || len(schema.Tasks) > 0 || len(schema.Packages) > 0 || len(schema.Events) > 0 || len(schema.Synonyms) > 0 ||
			slices.ContainsFunc(schema.EnumTypes, func(e *metadatapb.EnumTypeMetadata) bool { return e != nil && !e.SkipDump }) ||
			slices.ContainsFunc(schema.CompositeTypes, func(c *metadatapb.CompositeTypeMetadata) bool { return c != nil && !c.SkipDump }) ||
			slices.ContainsFunc(schema.Collations, func(c *metadatapb.CollationMetadata) bool { return c != nil && !c.SkipDump }) ||
			slices.ContainsFunc(schema.RangeTypes, func(r *metadatapb.RangeTypeMetadata) bool { return r != nil && !r.SkipDump }) {
			return fmt.Errorf("schema %s contains schema objects unsupported by Redshift DDL renderer", quoteDefinitionIdentifier(schema.Name))
		}
		for _, fn := range schema.Functions {
			if fn == nil {
				continue
			}
			if fn.SkipDump {
				fnEntries = append(fnEntries, validationOnlyFunctionEntry(schema.Name, fn))
				fnNodeIndices = append(fnNodeIndices, -1)
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
			params, retType, parsed, err := validateRoutineIdentity("function", schema.Name, fn.Name, fn.Signature, def, false, caseSensitive)
			if err != nil {
				return err
			}
			sigKey := redshiftFunctionIdentityKey(schema.Name, fn.Name, fn.Signature, params, parsed, caseSensitive)
			if _, exists := seenFnSigs[sigKey]; exists {
				return fmt.Errorf("duplicate function %s in schema %s cannot be restored", fn.Name, quoteDefinitionIdentifier(schema.Name))
			}
			seenFnSigs[sigKey] = struct{}{}
			entry := functionEntry{
				schema:     schema.Name,
				name:       fn.Name,
				def:        def,
				params:     params,
				returnType: retType,
				parsed:     parsed,
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
		if dep == i && items[i].kind != objectKindView && items[i].kind != objectKindMaterializedView {
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
	addCallDeps := func(i int, callerSchema string, calls []functionCall) error {
		for _, call := range calls {
			for _, fnEntryIdx := range fnCat.resolveCall(callerSchema, call) {
				depNode := fnCat.nodeIdx[fnEntryIdx]
				if depNode < 0 {
					continue
				}
				if items[i].kind == objectKindFunc && depNode == i {
					return fmt.Errorf("function %s is self-recursive and cannot be restored because Redshift requires callees to exist first", qualifyName(items[i].schema, items[i].name))
				}
				addDep(i, depNode)
			}
		}
		return nil
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
			if len(item.fn.depTables) > 0 {
				item.fn.def = qualifyRedshiftFunctionDependencies(item.fn.def, item.schema, item.fn.depTables, caseSensitive)
				items[i].fn.def = item.fn.def
			}
			if err := addCallDeps(i, item.schema, collectFunctionCalls(item.fn.def, item.schema, item.fn.params, fnCat)); err != nil {
				return err
			}
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
				for _, ref := range collectRelationRefs(item.view.definition) {
					addRelDep(i, item.schema, ref.catalog, ref.schema, ref.name)
				}
			}
			if err := addCallDeps(i, item.schema, collectSQLFunctionCalls(item.view.definition, item.schema, fnCat)); err != nil {
				return err
			}
		}
	}

	for _, item := range items {
		if item.kind != objectKindTable {
			continue
		}
		if err := validateTableDefaultExprs(item.table, fnCat); err != nil {
			return err
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
			if item.fk.comment != "" && item.fk.name == "" {
				return fmt.Errorf("foreign key comment on %s cannot be restored because the constraint is unnamed", item.table.qualifiedName())
			}
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

	for _, item := range items {
		if item.kind != objectKindTable {
			continue
		}
		for _, fk := range item.table.foreignKeys {
			if refTable := findTableInItems(items, fk.referencedSchema, fk.referencedTable, caseSensitive); refTable != nil {
				for _, refCol := range fk.referencedColumns {
					if !hasTableColumn(refTable.columns, refCol, caseSensitive) {
						return fmt.Errorf("foreign key %s on %s references nonexistent column %q on %s", fk.name, item.table.qualifiedName(), refCol, refTable.qualifiedName())
					}
				}
				if !hasMatchingUniqueIndex(refTable, fk.referencedColumns, caseSensitive) {
					return fmt.Errorf("foreign key %s on %s references %s(%s), which is not covered by a primary key or unique constraint", fk.name, item.table.qualifiedName(), refTable.qualifiedName(), strings.Join(fk.referencedColumns, ", "))
				}
			}
		}
	}

	ordered, cyclic := topoOrder(len(items), dependencies)
	for _, idx := range cyclic {
		switch items[idx].kind {
		case objectKindFunc:
			return fmt.Errorf("cyclic function dependency involving %s cannot be restored because Redshift requires callees to exist first", qualifyName(items[idx].schema, items[idx].fn.name))
		case objectKindView, objectKindMaterializedView:
			return fmt.Errorf("cyclic view dependency involving %s cannot be restored because standard Redshift views require referenced relations to exist first", qualifyName(items[idx].schema, items[idx].name))
		}
	}
	for _, idx := range ordered {
		item := items[idx]
		switch item.kind {
		case objectKindFunc:
			if _, err := w.WriteString(getFunctionAnnouncement(item.fn.name)); err != nil {
				return err
			}
			if err := writeWithSchemaSearchPath(w, item.schema, item.fn.def+"\n"); err != nil {
				return err
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
			if err := writeTableWithOptions(w, item.table, inlineFKs, false, false, caseSensitive); err != nil {
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
	loc     int
}

func validateTableDefaultExprs(table *tableState, fnCat *funcCatalog) error {
	if table == nil {
		return nil
	}
	for _, expr := range table.defaultExprs {
		for _, call := range collectSQLFunctionCalls("SELECT "+expr, table.schema, fnCat) {
			if matches := fnCat.resolveCall(table.schema, call); len(matches) > 0 {
				fnEntry := fnCat.entries[matches[0]]
				return fmt.Errorf("table %s column default cannot call user-defined function %s because Redshift disallows UDFs in default_expr", table.qualifiedName(), qualifyName(fnEntry.schema, fnEntry.name))
			}
			if call.schema != "" && !strings.EqualFold(call.schema, "pg_catalog") && !strings.EqualFold(call.schema, "information_schema") {
				return fmt.Errorf("table %s column default cannot call user-defined function %s because Redshift disallows UDFs in default_expr", table.qualifiedName(), qualifyName(call.schema, call.name))
			}
		}
	}
	return nil
}

func collectRelationRefs(sql string) []relationRef {
	if strings.TrimSpace(sql) == "" {
		return nil
	}
	stmts, err := pgparser.Parse(sql)
	if err != nil || stmts == nil || len(stmts.Items) == 0 {
		return nil
	}
	var refs []relationRef
	for _, item := range stmts.Items {
		root := item
		if raw, ok := root.(*nodes.RawStmt); ok {
			root = raw.Stmt
		}
		var walkRoot nodes.Node
		switch s := root.(type) {
		case *nodes.ViewStmt:
			walkRoot = s.Query
		case *nodes.CreateTableAsStmt:
			walkRoot = s.Query
		case *nodes.CreateFunctionStmt:
			walkRoot = s.SqlBody
		default:
			walkRoot = root
		}
		if walkRoot != nil {
			nodes.Walk(&scopedRelVisitor{refs: &refs}, walkRoot)
		}
	}
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
			*v.refs = append(*v.refs, relationRef{catalog: n.Catalogname, schema: n.Schemaname, name: n.Relname, loc: n.Loc.Start})
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
func topoOrder(n int, dependencies []map[int]struct{}) ([]int, []int) {
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
	var cyclic []int
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
			cyclic = append(cyclic, u)
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
	return ordered, cyclic
}

func findTableInItems(items []orderedObject, schemaName, tableName string, caseSensitive bool) *tableState {
	if schemaName == "" || tableName == "" {
		return nil
	}
	for _, it := range items {
		if it.kind == objectKindTable && identEqual(it.schema, schemaName, caseSensitive) && identEqual(it.name, tableName, caseSensitive) {
			return it.table
		}
	}
	return nil
}

// writeTableWithOptions writes one table, its comments (table, column, and
// named key/FK constraints), and optionally its ownership transfer. Column
// comments follow the column order the CREATE TABLE body uses so output is
// deterministic across runs.
func writeTableWithOptions(w io.StringWriter, table *tableState, foreignKeys []*foreignKeyState, allFKs, emitOwner, caseSensitive bool) error {
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
	seenConstraints := make(map[string]struct{}, len(table.indexes)+len(table.foreignKeys))
	checkConstraintName := func(name string) error {
		if name == "" {
			return nil
		}
		key := identKey(name, caseSensitive)
		if _, exists := seenConstraints[key]; exists {
			return fmt.Errorf("table %s has duplicate constraint name %s", table.qualifiedName(), quoteDefinitionIdentifier(name))
		}
		seenConstraints[key] = struct{}{}
		return nil
	}
	primaryCount := 0
	for _, idx := range table.indexes {
		if !idx.primary && !idx.unique {
			return fmt.Errorf("index %s on %s cannot be restored because Redshift does not support secondary indexes", idx.name, table.qualifiedName())
		}
		if err := checkConstraintName(idx.name); err != nil {
			return err
		}
		if idx.primary {
			primaryCount++
			if primaryCount > 1 {
				return fmt.Errorf("table %s declares multiple primary keys", table.qualifiedName())
			}
		}
		if len(idx.keys) == 0 {
			return fmt.Errorf("index %s on %s has no columns", idx.name, table.qualifiedName())
		}
		if slices.ContainsFunc(idx.lengths, func(n int64) bool { return n != -1 && n != 0 }) {
			return fmt.Errorf("index %s on %s has key lengths, which Redshift constraints cannot restore", idx.name, table.qualifiedName())
		}
		if slices.Contains(idx.descending, true) {
			return fmt.Errorf("index %s on %s is descending, which Redshift constraints cannot restore", idx.name, table.qualifiedName())
		}
		if !isSupportedRedshiftIndexType(idx.tp) {
			return fmt.Errorf("index %s on %s has unsupported type %q", idx.name, table.qualifiedName(), idx.tp)
		}
		if idx.comment != "" && idx.name == "" {
			return fmt.Errorf("index comment on %s cannot be restored because the constraint is unnamed", table.qualifiedName())
		}
		for _, key := range idx.keys {
			colName, ok := parseConstraintColumnName(key)
			if !ok {
				return fmt.Errorf("index %s on %s key %q is not a table column", idx.name, table.qualifiedName(), key)
			}
			col, exists := lookupTableColumn(table.columns, colName, caseSensitive)
			if !exists {
				return fmt.Errorf("index %s on %s key %q is not a table column", idx.name, table.qualifiedName(), key)
			}
			if idx.primary && col.nullable {
				return fmt.Errorf("primary key %s on %s column %q cannot be nullable", idx.name, table.qualifiedName(), col.name)
			}
		}
	}
	for _, fk := range sortedForeignKeys(table.foreignKeys) {
		if err := checkConstraintName(fk.name); err != nil {
			return err
		}
		if len(fk.columns) == 0 || len(fk.columns) != len(fk.referencedColumns) {
			return fmt.Errorf("foreign key %s on %s has %d columns and %d referenced columns", fk.name, table.qualifiedName(), len(fk.columns), len(fk.referencedColumns))
		}
		for _, col := range fk.columns {
			if !hasTableColumn(table.columns, col, caseSensitive) {
				return fmt.Errorf("foreign key %s on %s column %q is not a table column", fk.name, table.qualifiedName(), col)
			}
		}
		if fk.comment != "" && fk.name == "" {
			return fmt.Errorf("foreign key comment on %s cannot be restored because the constraint is unnamed", table.qualifiedName())
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
	if view.duplicateColumn != "" {
		return fmt.Errorf("view %s has duplicate column %s", view.qualifiedName(), quoteDefinitionIdentifier(view.duplicateColumn))
	}
	if view.unsupportedReason != "" {
		return fmt.Errorf("view %s %s", view.qualifiedName(), view.unsupportedReason)
	}
	buf := &strings.Builder{}
	if err := view.toString(buf); err != nil {
		return err
	}
	if view.comment != "" {
		if _, err := fmt.Fprintf(buf, "COMMENT ON VIEW %s IS '%s';\n", view.qualifiedName(), escapeStringLiteral(view.comment)); err != nil {
			return err
		}
	}
	if isLateBindingView(view.definition) {
		for _, column := range sortedColumns(view.columns) {
			if column.comment != "" {
				return fmt.Errorf("column %s comment on late-binding view %s cannot be restored because Redshift does not support COMMENT ON COLUMN for late-binding views", quoteDefinitionIdentifier(column.name), view.qualifiedName())
			}
		}
	} else {
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
	if mv.unsupportedReason != "" {
		return fmt.Errorf("materialized view %s %s", mv.qualifiedName(), mv.unsupportedReason)
	}
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

func convertToDatabaseState(database *metadatapb.DatabaseSchemaMetadata, caseSensitive bool) *databaseState {
	state := newDatabaseState()
	state.name = database.Name
	for _, schema := range database.Schemas {
		if schema == nil || schema.SkipDump {
			continue
		}
		state.schemas[schema.Name] = convertToSchemaState(schema, caseSensitive)
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

func convertToSchemaState(schema *metadatapb.SchemaMetadata, caseSensitive bool) *schemaState {
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
		state.views[view.Name] = convertToViewState(i, schema.Name, view, caseSensitive)
	}
	for i, mv := range schema.MaterializedViews {
		if mv == nil || mv.SkipDump {
			continue
		}
		state.materializedViews[mv.Name] = convertToMaterializedViewState(i, schema.Name, mv, caseSensitive)
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
	case table.Engine != "":
		state.unsupportedReason = fmt.Sprintf("engine %q cannot be restored by Redshift DDL renderer", table.Engine)
	case table.Collation != "":
		state.unsupportedReason = fmt.Sprintf("table collation %q cannot be restored by Redshift DDL renderer", table.Collation)
	case table.Charset != "":
		state.unsupportedReason = fmt.Sprintf("table charset %q cannot be restored by Redshift DDL renderer", table.Charset)
	case table.PrimaryKeyType != "":
		state.unsupportedReason = fmt.Sprintf("primary key type %q cannot be restored by Redshift DDL renderer", table.PrimaryKeyType)
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
		if msg := unsupportedRedshiftColumnReason(column, false); msg != "" && state.unsupportedReason == "" {
			state.unsupportedReason = msg
		}
		if strings.TrimSpace(column.Type) == "" && state.unsupportedReason == "" {
			state.unsupportedReason = fmt.Sprintf("column %s has no data type", quoteDefinitionIdentifier(column.Name))
		}
		if hasDefaultSubquery(column.Default) && state.unsupportedReason == "" {
			state.unsupportedReason = fmt.Sprintf("column %s default %q contains a subquery, which Redshift does not support", quoteDefinitionIdentifier(column.Name), column.Default)
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
	matchType         string
	comment           string
}

func isSupportedRedshiftFKAction(action string) bool {
	trimmed := strings.TrimSpace(strings.ToUpper(action))
	return trimmed == "" || trimmed == "NO ACTION"
}

func isSupportedRedshiftFKMatchType(matchType string) bool {
	switch strings.TrimSpace(strings.ToUpper(matchType)) {
	case "", "SIMPLE", "S", "U":
		return true
	default:
		return false
	}
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
		matchType:         foreignKey.MatchType,
		comment:           foreignKey.Comment,
	}
}

func (f *foreignKeyState) toString(buf *strings.Builder) error {
	if f.referencedTable == "" {
		return fmt.Errorf("foreign key %s has empty referenced table", f.name)
	}
	if len(f.columns) == 0 || len(f.columns) != len(f.referencedColumns) {
		return fmt.Errorf("foreign key %s has %d columns and %d referenced columns", f.name, len(f.columns), len(f.referencedColumns))
	}
	if !isSupportedRedshiftFKAction(f.onDelete) {
		return fmt.Errorf("foreign key %s has unsupported ON DELETE action %q", f.name, f.onDelete)
	}
	if !isSupportedRedshiftFKAction(f.onUpdate) {
		return fmt.Errorf("foreign key %s has unsupported ON UPDATE action %q", f.name, f.onUpdate)
	}
	if !isSupportedRedshiftFKMatchType(f.matchType) {
		return fmt.Errorf("foreign key %s has unsupported MATCH type %q", f.name, f.matchType)
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
	id         int
	name       string
	keys       []string
	lengths    []int64
	descending []bool
	primary    bool
	unique     bool
	tp         string
	comment    string
}

func isSupportedRedshiftIndexType(tp string) bool {
	trimmed := strings.TrimSpace(strings.ToUpper(tp))
	return trimmed == "" || trimmed == "BTREE"
}

func convertToIndexState(id int, index *metadatapb.IndexMetadata) *indexState {
	return &indexState{
		id:         id,
		name:       index.Name,
		keys:       index.Expressions,
		lengths:    index.KeyLength,
		descending: index.Descending,
		primary:    index.Primary,
		unique:     index.Unique,
		tp:         index.Type,
		comment:    index.Comment,
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
	if c.collationSchema != "" || c.collation != "" {
		colKeyword, ok := supportedRedshiftCollation(c.collationSchema, c.collation)
		if !ok {
			return fmt.Errorf("column %s collation %s is not supported by Redshift CREATE TABLE", quoteDefinitionIdentifier(c.name), qualifyName(c.collationSchema, c.collation))
		}
		if _, err := fmt.Fprintf(buf, " COLLATE %s", colKeyword); err != nil {
			return err
		}
	}
	if c.isIdentity {
		if !isSupportedRedshiftIdentityType(c.tp) {
			return fmt.Errorf("column %s identity is not supported on non-integer type %q", quoteDefinitionIdentifier(c.name), c.tp)
		}
		seed, step := c.identitySeed, c.identityIncrement
		if step == 0 {
			if seed == 0 {
				seed = 1
			}
			step = 1
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
	if c.isIdentity && c.defaultValue != nil {
		return fmt.Errorf("column %s identity cannot also specify a default value", quoteDefinitionIdentifier(c.name))
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

func supportedRedshiftCollation(schema, name string) (string, bool) {
	if schema != "" {
		return "", false
	}
	switch strings.ToUpper(name) {
	case "CASE_SENSITIVE", "CS", "CASE_INSENSITIVE", "CI":
		return strings.ToUpper(name), true
	default:
		return "", false
	}
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
	hasDisallowed := false
	nodes.Inspect(rt.Val, func(n nodes.Node) bool {
		switch n.(type) {
		case *nodes.ColumnRef, *nodes.SubLink:
			hasDisallowed = true
			return false
		}
		return true
	})
	return !hasDisallowed
}

func hasDefaultSubquery(value string) bool {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || isQuotedSQLLiteral(trimmed) {
		return false
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
	if !ok || rt.Val == nil {
		return false
	}
	found := false
	nodes.Inspect(rt.Val, func(n nodes.Node) bool {
		if _, ok := n.(*nodes.SubLink); ok {
			found = true
			return false
		}
		return true
	})
	return found
}

func parseConstraintColumnName(key string) (string, bool) {
	stmts, err := pgparser.Parse("CREATE TABLE _t (PRIMARY KEY (" + key + "))")
	if err != nil || stmts == nil || len(stmts.Items) != 1 {
		return "", false
	}
	stmt := stmts.Items[0]
	if raw, ok := stmt.(*nodes.RawStmt); ok {
		stmt = raw.Stmt
	}
	cs, ok := stmt.(*nodes.CreateStmt)
	if !ok || cs.TableElts == nil || len(cs.TableElts.Items) != 1 {
		return "", false
	}
	con, ok := cs.TableElts.Items[0].(*nodes.Constraint)
	if !ok || con.Keys == nil || len(con.Keys.Items) != 1 {
		return "", false
	}
	s, ok := con.Keys.Items[0].(*nodes.String)
	if !ok || s.Str == "" {
		return "", false
	}
	return s.Str, true
}

func redshiftFunctionIdentityKey(schemaName, fnName, signature string, params []functionParam, parsed, caseSensitive bool) string {
	var sigPart string
	if parsed {
		parts := make([]string, len(params))
		for i, p := range params {
			parts[i] = p.typeName
		}
		sigPart = strings.Join(parts, ",")
	} else {
		sigPart = strings.ToLower(strings.TrimSpace(signature))
	}
	return identKey(schemaName, caseSensitive) + "\x00" + identKey(fnName, caseSensitive) + "\x00" + sigPart
}

func hasMatchingUniqueIndex(table *tableState, cols []string, caseSensitive bool) bool {
	if len(cols) == 0 {
		return false
	}
	for _, idx := range table.indexes {
		if (!idx.primary && !idx.unique) || len(idx.keys) != len(cols) {
			continue
		}
		if sameConstraintColumns(idx.keys, cols, caseSensitive) {
			return true
		}
	}
	return false
}

func sameConstraintColumns(idxKeys, fkCols []string, caseSensitive bool) bool {
	if len(idxKeys) != len(fkCols) {
		return false
	}
	remaining := make(map[string]int, len(idxKeys))
	for _, k := range idxKeys {
		col, ok := parseConstraintColumnName(k)
		if !ok {
			return false
		}
		remaining[identKey(col, caseSensitive)]++
	}
	for _, col := range fkCols {
		key := identKey(col, caseSensitive)
		if remaining[key] == 0 {
			return false
		}
		remaining[key]--
	}
	return true
}

func lookupTableColumn(cols map[string]*columnState, name string, caseSensitive bool) (*columnState, bool) {
	if col, ok := cols[name]; ok {
		return col, true
	}
	if caseSensitive {
		return nil, false
	}
	for k, col := range cols {
		if strings.EqualFold(k, name) {
			return col, true
		}
	}
	return nil, false
}

func hasTableColumn(cols map[string]*columnState, name string, caseSensitive bool) bool {
	_, ok := lookupTableColumn(cols, name, caseSensitive)
	return ok
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
	id                int
	schema            string
	name              string
	definition        string
	comment           string
	duplicateColumn   string
	unsupportedReason string
	columns           map[string]*columnState
	deps              []*metadatapb.DependencyColumn
	caseSensitive     bool
}

func unsupportedRedshiftColumnReason(column *metadatapb.ColumnMetadata, isView bool) string {
	switch {
	case column.Generation != nil:
		return fmt.Sprintf("column %s generated expression cannot be restored because Redshift does not support generated columns", quoteDefinitionIdentifier(column.Name))
	case column.DefaultOnNull || column.OnUpdate != "" || column.CharacterSet != "" || column.Srid != nil || column.IsInvisible || column.DefaultConstraintName != "":
		return fmt.Sprintf("column %s attributes cannot be restored by Redshift DDL renderer", quoteDefinitionIdentifier(column.Name))
	case isView && (len(column.Options) > 0 || column.Collation != "" || column.CollationSchema != "" || column.Default != "" || column.IsIdentity):
		return fmt.Sprintf("column %s attributes cannot be restored by Redshift DDL renderer", quoteDefinitionIdentifier(column.Name))
	case !isView && column.IsIdentity && !isSupportedRedshiftIdentityType(column.Type):
		return fmt.Sprintf("column %s identity is not supported on non-integer type %q", quoteDefinitionIdentifier(column.Name), column.Type)
	case !isView && column.IsIdentity && column.Default != "":
		return fmt.Sprintf("column %s identity cannot also specify a default value %q", quoteDefinitionIdentifier(column.Name), column.Default)
	case !isView && !column.IsIdentity && (column.IdentityGeneration != 0 || column.IdentitySeed != 0 || column.IdentityIncrement != 0):
		return fmt.Sprintf("column %s specifies identity attributes without IsIdentity", quoteDefinitionIdentifier(column.Name))
	case !isView && (column.CollationSchema != "" || column.Collation != ""):
		if _, ok := supportedRedshiftCollation(column.CollationSchema, column.Collation); !ok {
			return fmt.Sprintf("column %s collation %s is not supported by Redshift CREATE TABLE", quoteDefinitionIdentifier(column.Name), qualifyName(column.CollationSchema, column.Collation))
		}
		return ""
	default:
		return ""
	}
}

func isSupportedRedshiftIdentityType(tp string) bool {
	switch strings.ToLower(strings.TrimSpace(tp)) {
	case "int", "integer", "int4", "bigint", "int8":
		return true
	default:
		return false
	}
}

func convertToViewState(id int, schemaName string, view *metadatapb.ViewMetadata, caseSensitive bool) *viewState {
	var unsupportedReason string
	switch {
	case len(view.Triggers) > 0:
		unsupportedReason = "triggers cannot be restored because Redshift does not support triggers"
	case len(view.Rules) > 0:
		unsupportedReason = "rules cannot be restored because Redshift does not support rules"
	}
	cols := make(map[string]*columnState, len(view.Columns))
	seenCols := make(map[string]struct{}, len(view.Columns))
	var duplicateColumn string
	for i, col := range view.Columns {
		if col == nil {
			continue
		}
		key := identKey(col.Name, caseSensitive)
		if _, exists := seenCols[key]; exists && duplicateColumn == "" {
			duplicateColumn = col.Name
		}
		seenCols[key] = struct{}{}
		if msg := unsupportedRedshiftColumnReason(col, true); msg != "" && unsupportedReason == "" {
			unsupportedReason = msg
		}
		cols[col.Name] = convertToColumnState(i, col)
	}
	return &viewState{
		id:                id,
		schema:            schemaName,
		name:              view.Name,
		definition:        view.Definition,
		comment:           view.Comment,
		duplicateColumn:   duplicateColumn,
		unsupportedReason: unsupportedReason,
		columns:           cols,
		deps:              view.DependencyColumns,
		caseSensitive:     caseSensitive,
	}
}

func convertToMaterializedViewState(id int, schemaName string, mv *metadatapb.MaterializedViewMetadata, caseSensitive bool) *viewState {
	var unsupportedReason string
	switch {
	case len(mv.Indexes) > 0:
		unsupportedReason = "indexes cannot be restored by Redshift DDL renderer"
	case len(mv.Triggers) > 0:
		unsupportedReason = "triggers cannot be restored because Redshift does not support triggers"
	}
	return &viewState{
		id:                id,
		schema:            schemaName,
		name:              mv.Name,
		definition:        mv.Definition,
		comment:           mv.Comment,
		unsupportedReason: unsupportedReason,
		deps:              mv.DependencyColumns,
		caseSensitive:     caseSensitive,
	}
}

func (v *viewState) validateTarget(kind string, rel *nodes.RangeVar) error {
	if rel == nil {
		return nil
	}
	target := qualifyName(rel.Schemaname, rel.Relname)
	if rel.Catalogname != "" {
		target = quoteDefinitionIdentifier(rel.Catalogname) + "." + target
		return fmt.Errorf("%s %s definition target %s does not match metadata identity", kind, v.qualifiedName(), target)
	}
	if (rel.Schemaname != "" && !identEqual(rel.Schemaname, v.schema, v.caseSensitive)) ||
		(rel.Relname != "" && !identEqual(rel.Relname, v.name, v.caseSensitive)) {
		return fmt.Errorf("%s %s definition target %s does not match metadata identity", kind, v.qualifiedName(), target)
	}
	return nil
}

type redshiftDepTarget struct {
	schema string
	table  string
}

func resolveRedshiftDepTargetSchema(relName, defaultSchema string, deps []redshiftDepTarget, caseSensitive bool) string {
	var matched string
	for _, dep := range deps {
		if dep.schema == "" || !identEqual(dep.table, relName, caseSensitive) {
			continue
		}
		if matched == "" {
			matched = dep.schema
		} else if !identEqual(matched, dep.schema, caseSensitive) {
			if identEqual(dep.schema, defaultSchema, caseSensitive) || identEqual(matched, defaultSchema, caseSensitive) {
				return defaultSchema
			}
			return ""
		}
	}
	return matched
}

func qualifyRedshiftRelationDeps(stmt, schemaName string, deps []redshiftDepTarget, caseSensitive bool) string {
	if len(deps) == 0 {
		return stmt
	}
	refs := collectRelationRefs(stmt)
	if len(refs) == 0 {
		return stmt
	}
	type prefixEdit struct {
		offset int
		prefix string
	}
	var edits []prefixEdit
	seen := make(map[int]struct{})
	for _, ref := range refs {
		if ref.catalog != "" || ref.schema != "" || ref.loc < 0 || ref.loc >= len(stmt) {
			continue
		}
		depSchema := resolveRedshiftDepTargetSchema(ref.name, schemaName, deps, caseSensitive)
		if depSchema == "" || identEqual(depSchema, schemaName, caseSensitive) {
			continue
		}
		if _, dup := seen[ref.loc]; !dup {
			seen[ref.loc] = struct{}{}
			edits = append(edits, prefixEdit{offset: ref.loc, prefix: quoteDefinitionIdentifier(depSchema) + "."})
		}
	}
	if len(edits) == 0 {
		return stmt
	}
	slices.SortFunc(edits, func(a, b prefixEdit) int { return b.offset - a.offset })
	out := stmt
	for _, ed := range edits {
		out = out[:ed.offset] + ed.prefix + out[ed.offset:]
	}
	return out
}

func qualifyRedshiftViewDependencies(stmt, schemaName string, deps []*metadatapb.DependencyColumn, caseSensitive bool) string {
	if len(deps) == 0 {
		return stmt
	}
	targets := make([]redshiftDepTarget, 0, len(deps))
	for _, dep := range deps {
		if dep != nil {
			targets = append(targets, redshiftDepTarget{schema: dep.Schema, table: dep.Table})
		}
	}
	return qualifyRedshiftRelationDeps(stmt, schemaName, targets, caseSensitive)
}

var redshiftLexSCONST = pgparser.NewLexer("''").NextToken().Type

func findFunctionASBodyToken(src string, toks []pgparser.Token) (pgparser.Token, bool) {
	depth := 0
	afterParams := 0
	for i, tok := range toks {
		if tok.Type == '(' {
			depth++
		} else if tok.Type == ')' && depth > 0 {
			depth--
			if depth == 0 {
				afterParams = i + 1
				break
			}
		}
	}
	for i := afterParams; i+1 < len(toks); i++ {
		if strings.EqualFold(toks[i].Str, "as") && toks[i+1].Type == redshiftLexSCONST {
			bodyTok := toks[i+1]
			if bodyTok.Loc >= 0 && bodyTok.End <= len(src) && bodyTok.Loc < bodyTok.End {
				return bodyTok, true
			}
			break
		}
	}
	return pgparser.Token{}, false
}

func qualifyRedshiftFunctionDependencies(def, schemaName string, deps []*metadatapb.DependencyTable, caseSensitive bool) string {
	if len(deps) == 0 {
		return def
	}
	targets := make([]redshiftDepTarget, 0, len(deps))
	for _, dep := range deps {
		if dep != nil {
			targets = append(targets, redshiftDepTarget{schema: dep.Schema, table: dep.Table})
		}
	}
	lex := pgparser.NewLexer(def)
	var toks []pgparser.Token
	for {
		tok := lex.NextToken()
		if tok.Type == 0 {
			break
		}
		toks = append(toks, tok)
	}
	if bodyTok, ok := findFunctionASBodyToken(def, toks); ok {
		newBody := qualifyRedshiftRelationDeps(bodyTok.Str, schemaName, targets, caseSensitive)
		if newBody == bodyTok.Str {
			return def
		}
		raw := def[bodyTok.Loc:bodyTok.End]
		var quoted string
		switch {
		case strings.HasPrefix(raw, "$"):
			if tagEnd := strings.IndexByte(raw[1:], '$'); tagEnd >= 0 {
				tag := raw[:tagEnd+2]
				if strings.HasSuffix(raw, tag) {
					quoted = tag + newBody + tag
				}
			}
		case strings.HasPrefix(raw, "'"):
			quoted = "'" + strings.ReplaceAll(newBody, "'", "''") + "'"
		}
		if quoted == "" {
			return def
		}
		return def[:bodyTok.Loc] + quoted + def[bodyTok.End:]
	}
	return qualifyRedshiftRelationDeps(def, schemaName, targets, caseSensitive)
}

func (v *viewState) toString(buf io.StringWriter) error {
	trimmed := strings.TrimSpace(v.definition)
	if trimmed == "" {
		return fmt.Errorf("view %s has no definition to restore", v.qualifiedName())
	}
	var stmt string
	if vs := parseViewStmt(trimmed); vs != nil {
		if err := v.validateTarget("view", vs.View); err != nil {
			return err
		}
		stmt = trimmed
	} else {
		stmt = fmt.Sprintf("CREATE OR REPLACE VIEW %s AS %s", v.qualifiedName(), v.definition)
	}
	stmt = qualifyRedshiftViewDependencies(stmt, v.schema, v.deps, v.caseSensitive)
	stmt = appendTerminator(stmt) + "\n"
	if v.schema != "" {
		return writeWithSchemaSearchPath(buf, v.schema, stmt)
	}
	_, err := buf.WriteString(stmt)
	return err
}

func (v *viewState) toMaterializedString(buf io.StringWriter) error {
	trimmed := strings.TrimSpace(v.definition)
	if trimmed == "" {
		return fmt.Errorf("materialized view %s has no definition to restore", v.qualifiedName())
	}
	var stmt string
	if ctas := parseMaterializedViewStmt(trimmed); ctas != nil {
		if ctas.Into != nil {
			if err := v.validateTarget("materialized view", ctas.Into.Rel); err != nil {
				return err
			}
		}
		stmt = trimmed
	} else {
		stmt = fmt.Sprintf("CREATE MATERIALIZED VIEW %s AS %s", v.qualifiedName(), v.definition)
	}
	stmt = qualifyRedshiftViewDependencies(stmt, v.schema, v.deps, v.caseSensitive)
	stmt = appendTerminator(stmt) + "\n"
	if v.schema != "" {
		return writeWithSchemaSearchPath(buf, v.schema, stmt)
	}
	_, err := buf.WriteString(stmt)
	return err
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

func parseMaterializedViewStmt(definition string) *nodes.CreateTableAsStmt {
	stmts, err := pgparser.Parse(definition)
	if err != nil || stmts == nil || len(stmts.Items) != 1 {
		return nil
	}
	stmt := stmts.Items[0]
	if raw, ok := stmt.(*nodes.RawStmt); ok {
		stmt = raw.Stmt
	}
	ctas, ok := stmt.(*nodes.CreateTableAsStmt)
	if !ok || ctas.Objtype != nodes.OBJECT_MATVIEW {
		return nil
	}
	return ctas
}
