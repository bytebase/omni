package completion

import (
	"context"
	"sort"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/yacc"
)

// CandidateType represents the type of completion candidate.
type CandidateType string

const (
	CandidateColumn   CandidateType = "COLUMN"
	CandidateTable    CandidateType = "TABLE"
	CandidateView     CandidateType = "VIEW"
	CandidateSchema   CandidateType = "SCHEMA"
	CandidateSequence CandidateType = "SEQUENCE"
	CandidateKeyword  CandidateType = "KEYWORD"
	CandidateFunction CandidateType = "FUNCTION"
)

// Candidate represents a single completion suggestion.
type Candidate struct {
	Text       string        `yaml:"text"`
	Type       CandidateType `yaml:"type"`
	Definition string        `yaml:"definition"`
	Comment    string        `yaml:"comment"`
}

// MetadataProvider provides database schema information for completion.
type MetadataProvider interface {
	GetSchemaNames(ctx context.Context) []string
	GetTables(ctx context.Context, schema string) []TableInfo
	GetViews(ctx context.Context, schema string) []string
	GetSequences(ctx context.Context, schema string) []string
	GetColumns(ctx context.Context, schema, table string) []ColumnInfo
}

// TableInfo describes a table.
type TableInfo struct {
	Name string
}

// ColumnInfo describes a column.
type ColumnInfo struct {
	Name       string
	Type       string
	NotNull    bool
	Definition string
}

// completionHint is the result of parse-table analysis.
type completionHint struct {
	validTokens map[int]bool
	contexts    []GrammarContext
}

// tableRef represents a table reference found in the SQL.
type tableRef struct {
	schema      string
	table       string
	alias       string
	isSubq      bool
	subqCols    []string
	subqFromRef []tableRef // inner FROM refs for subquery (for SELECT *)
}

// cteInfo represents a CTE found in the SQL.
type cteInfo struct {
	name    string
	columns []string
	selStmt *nodes.SelectStmt
}

// aliasInfo maps alias names to their resolved table/schema.
type aliasInfo struct {
	alias  string
	schema string
	table  string
}

// expand converts a completionHint into user-facing candidates.
func expand(ctx context.Context, hint completionHint, fullSQL string, cursorOffset int, meta MetadataProvider) []Candidate {
	var candidates []Candidate
	candidateMap := make(map[string]bool)

	addCandidate := func(c Candidate) {
		key := string(c.Type) + ":" + c.Text
		if !candidateMap[key] {
			candidateMap[key] = true
			candidates = append(candidates, c)
		}
	}

	refs, ctes, _ := collectAllRefs(fullSQL, cursorOffset)

	qualifier := extractQualifier(fullSQL, cursorOffset)
	if qualifier.hasDot {
		expandQualified(ctx, qualifier, refs, ctes, meta, addCandidate)
		return sortCandidates(candidates)
	}

	hasContext := func(c GrammarContext) bool {
		for _, ctx := range hint.contexts {
			if ctx == c {
				return true
			}
		}
		return false
	}

	isRelationCtx := hasContext(CtxRelationRef)
	isColumnCtx := hasContext(CtxColumnRef)
	isFuncCtx := hasContext(CtxFuncName)
	isTypeCtx := hasContext(CtxTypeName)

	if isRelationCtx {
		activeSchema := inferActiveSchema(refs)
		expandRelations(ctx, meta, refs, ctes, activeSchema, addCandidate)
	}
	if isColumnCtx {
		expandColumns(ctx, refs, ctes, meta, addCandidate)
		expandSelectAliases(fullSQL, cursorOffset, addCandidate)
		if !isRelationCtx {
			activeSchema := inferActiveSchema(refs)
			expandRelationsMinimal(ctx, meta, refs, ctes, activeSchema, addCandidate)
		}
	}
	if isFuncCtx {
		expandFunctions(addCandidate)
	}
	if isTypeCtx {
		expandTypes(addCandidate)
	}

	// Keywords: map valid tokens to keyword candidates
	expandKeywords(hint.validTokens, addCandidate)

	return sortCandidates(candidates)
}

// qualifierInfo holds info about a dotted qualifier before the cursor.
type qualifierInfo struct {
	hasDot bool
	parts  []string
}

// extractQualifier checks if the cursor is right after a dot.
func extractQualifier(sql string, cursorOffset int) qualifierInfo {
	if cursorOffset <= 0 || cursorOffset > len(sql) {
		return qualifierInfo{}
	}

	pos := cursorOffset - 1
	for pos >= 0 && sql[pos] == ' ' {
		pos--
	}
	if pos < 0 || sql[pos] != '.' {
		return qualifierInfo{}
	}

	var parts []string
	end := pos
	for {
		end--
		if end < 0 {
			break
		}
		identEnd := end + 1
		for end >= 0 && isIdentChar(rune(sql[end])) {
			end--
		}
		if end+1 < identEnd {
			parts = append([]string{strings.ToLower(sql[end+1 : identEnd])}, parts...)
		}
		if end >= 0 && sql[end] == '.' {
			continue
		}
		break
	}
	if len(parts) == 0 {
		return qualifierInfo{}
	}
	return qualifierInfo{hasDot: true, parts: parts}
}

// expandQualified handles completion after a dot (schema.| or alias.| or schema.table.|).
func expandQualified(ctx context.Context, q qualifierInfo, refs []tableRef, ctes []cteInfo, meta MetadataProvider, add func(Candidate)) {
	if meta == nil || len(q.parts) == 0 {
		return
	}

	// Two-level qualifier: schema.table.| → expand columns
	if len(q.parts) >= 2 {
		schema := q.parts[0]
		table := q.parts[1]
		cols := meta.GetColumns(ctx, schema, table)
		for _, col := range cols {
			add(Candidate{Text: QuoteIdentifier(col.Name), Type: CandidateColumn, Definition: col.Definition})
		}
		return
	}

	name := q.parts[0]

	// 1. Try as schema — but only if NOT also a table/alias reference.
	// Check alias/table match first so that alias.column takes priority over schema.table.
	// 2. Try as CTE alias
	for _, cte := range ctes {
		if strings.EqualFold(cte.name, name) {
			if len(cte.columns) > 0 {
				for _, col := range cte.columns {
					add(Candidate{Text: QuoteIdentifier(col), Type: CandidateColumn})
				}
				return
			}
			if cte.selStmt != nil {
				for _, col := range resolveSelectColumns(ctx, cte.selStmt, meta) {
					add(Candidate{Text: QuoteIdentifier(col), Type: CandidateColumn})
				}
				return
			}
		}
	}

	// 3. Try as table/subquery alias
	for _, ref := range refs {
		refName := ref.alias
		if refName == "" {
			refName = ref.table
		}
		if strings.EqualFold(refName, name) {
			if ref.isSubq {
				if len(ref.subqCols) > 0 {
					for _, col := range ref.subqCols {
						add(Candidate{Text: QuoteIdentifier(col), Type: CandidateColumn})
					}
					return
				}
				// SELECT * subquery — resolve columns from the subquery's source tables
				resolveSubqueryStarColumns(ctx, ref, meta, add)
				return
			}
			schema := ref.schema
			if schema == "" {
				schema = "public"
			}
			cols := meta.GetColumns(ctx, schema, ref.table)
			for _, col := range cols {
				add(Candidate{Text: QuoteIdentifier(col.Name), Type: CandidateColumn, Definition: col.Definition})
			}
			return
		}
	}

	// 4. Try as schema name
	tables := meta.GetTables(ctx, name)
	views := meta.GetViews(ctx, name)
	seqs := meta.GetSequences(ctx, name)

	if len(tables) > 0 || len(views) > 0 || len(seqs) > 0 {
		for _, t := range tables {
			add(Candidate{Text: QuoteIdentifier(t.Name), Type: CandidateTable})
		}
		for _, v := range views {
			add(Candidate{Text: QuoteIdentifier(v), Type: CandidateView})
		}
		for _, s := range seqs {
			add(Candidate{Text: QuoteIdentifier(s), Type: CandidateSequence})
		}
		return
	}
}

// expandRelations adds schemas, tables, views, sequences as candidates.
func expandRelations(ctx context.Context, meta MetadataProvider, refs []tableRef, ctes []cteInfo, activeSchema string, add func(Candidate)) {
	if meta == nil {
		return
	}

	schemas := meta.GetSchemaNames(ctx)
	defaultSchema := "public"

	for _, s := range schemas {
		add(Candidate{Text: QuoteIdentifier(s), Type: CandidateSchema})
	}

	if activeSchema != "" && activeSchema != defaultSchema {
		// Non-default schema context: show tables from active + default schemas
		for _, t := range meta.GetTables(ctx, activeSchema) {
			add(Candidate{Text: QuoteIdentifier(t.Name), Type: CandidateTable})
		}
		for _, t := range meta.GetTables(ctx, defaultSchema) {
			add(Candidate{Text: QuoteIdentifier(t.Name), Type: CandidateTable})
		}
		for _, v := range meta.GetViews(ctx, defaultSchema) {
			add(Candidate{Text: QuoteIdentifier(v), Type: CandidateView})
		}
	} else {
		for _, t := range meta.GetTables(ctx, defaultSchema) {
			add(Candidate{Text: QuoteIdentifier(t.Name), Type: CandidateTable})
		}
		for _, v := range meta.GetViews(ctx, defaultSchema) {
			add(Candidate{Text: QuoteIdentifier(v), Type: CandidateView})
		}
		for _, s := range meta.GetSequences(ctx, defaultSchema) {
			add(Candidate{Text: QuoteIdentifier(s), Type: CandidateSequence})
		}
		for _, schema := range schemas {
			if schema == defaultSchema {
				continue
			}
			for _, t := range meta.GetTables(ctx, schema) {
				add(Candidate{Text: QuoteIdentifier(schema) + "." + QuoteIdentifier(t.Name), Type: CandidateTable})
			}
			for _, s := range meta.GetSequences(ctx, schema) {
				add(Candidate{Text: QuoteIdentifier(schema) + "." + QuoteIdentifier(s), Type: CandidateSequence})
			}
		}
	}

	// Add CTE names as tables
	for _, cte := range ctes {
		add(Candidate{Text: QuoteIdentifier(cte.name), Type: CandidateTable})
	}

	// Add subquery/alias names as tables
	for _, ref := range refs {
		if ref.isSubq && ref.alias != "" {
			add(Candidate{Text: QuoteIdentifier(ref.alias), Type: CandidateTable})
		} else if ref.alias != "" && ref.alias != ref.table {
			add(Candidate{Text: QuoteIdentifier(ref.alias), Type: CandidateTable})
		}
	}
}

// expandColumns adds column candidates from table references.
func expandColumns(ctx context.Context, refs []tableRef, ctes []cteInfo, meta MetadataProvider, add func(Candidate)) {
	if meta == nil {
		return
	}

	cteNames := make(map[string]bool)
	for _, cte := range ctes {
		cteNames[strings.ToLower(cte.name)] = true
	}

	for _, ref := range refs {
		if ref.isSubq {
			if len(ref.subqCols) > 0 {
				for _, col := range ref.subqCols {
					add(Candidate{Text: QuoteIdentifier(col), Type: CandidateColumn})
				}
			} else {
				// SELECT * — resolve from inner FROM
				resolveSubqueryStarColumns(ctx, ref, meta, add)
			}
			continue
		}
		if ref.table == "" || cteNames[strings.ToLower(ref.table)] {
			continue
		}
		schema := ref.schema
		if schema == "" {
			schema = "public"
		}
		cols := meta.GetColumns(ctx, schema, ref.table)
		for _, col := range cols {
			add(Candidate{Text: QuoteIdentifier(col.Name), Type: CandidateColumn, Definition: col.Definition})
		}
	}

	for _, cte := range ctes {
		if len(cte.columns) > 0 {
			for _, col := range cte.columns {
				add(Candidate{Text: QuoteIdentifier(col), Type: CandidateColumn})
			}
		}
	}
}

// expandSelectAliases adds SELECT-list aliases for ORDER BY completion.
func expandSelectAliases(fullSQL string, cursorOffset int, add func(Candidate)) {
	cleanSQL := removeCursorMarker(fullSQL)
	upperSQL := strings.ToUpper(cleanSQL)

	beforeCursor := upperSQL
	if cursorOffset > 0 && cursorOffset <= len(upperSQL) {
		beforeCursor = upperSQL[:cursorOffset]
	}
	if !strings.Contains(beforeCursor, "ORDER BY") {
		return
	}

	// Parse the full SQL with a repair to get the AST
	fixed := repairSQL(cleanSQL, cursorOffset)
	if fixed == "" {
		return
	}
	result, err := yacc.Parse(fixed)
	if err != nil || result == nil {
		return
	}

	for _, item := range result.Items {
		var stmt nodes.Node
		if raw, ok := item.(*nodes.RawStmt); ok {
			stmt = raw.Stmt
		} else {
			stmt = item
		}
		sel, ok := stmt.(*nodes.SelectStmt)
		if !ok || sel == nil || sel.TargetList == nil {
			continue
		}
		for _, tgt := range sel.TargetList.Items {
			rt, ok := tgt.(*nodes.ResTarget)
			if !ok || rt.Name == "" {
				continue
			}
			add(Candidate{Text: QuoteIdentifier(rt.Name), Type: CandidateColumn})
		}
	}
}

// resolveSubqueryStarColumns resolves columns for a SELECT * subquery.
func resolveSubqueryStarColumns(ctx context.Context, ref tableRef, meta MetadataProvider, add func(Candidate)) {
	for _, innerRef := range ref.subqFromRef {
		schema := innerRef.schema
		if schema == "" {
			schema = "public"
		}
		cols := meta.GetColumns(ctx, schema, innerRef.table)
		for _, col := range cols {
			add(Candidate{Text: QuoteIdentifier(col.Name), Type: CandidateColumn})
		}
	}
}

// expandKeywords converts valid token IDs into keyword candidates.
func expandKeywords(validTokens map[int]bool, add func(Candidate)) {
	// Build a reverse map: internal goyacc token ID → keyword name
	for _, kw := range yacc.Keywords {
		internalID := mapKeywordToInternalID(kw.Token)
		if internalID > 0 && validTokens[internalID] {
			add(Candidate{Text: strings.ToUpper(kw.Name), Type: CandidateKeyword})
		}
	}
}

// mapKeywordToInternalID maps a parser keyword token constant to its goyacc-internal ID.
func mapKeywordToInternalID(tok int) int {
	if tok <= 0 {
		return 0
	}
	tok1 := yacc.Tok1()
	if tok < len(tok1) {
		return int(tok1[tok])
	}
	priv := yacc.Private()
	tok2 := yacc.Tok2()
	if tok >= priv && tok < priv+len(tok2) {
		return int(tok2[tok-priv])
	}
	tok3 := yacc.Tok3()
	for i := 0; i < len(tok3); i += 2 {
		if int(tok3[i]) == tok {
			return int(tok3[i+1])
		}
	}
	return 0
}

// expandFunctions adds common PostgreSQL built-in functions as candidates.
func expandFunctions(add func(Candidate)) {
	for _, fn := range builtinFunctions {
		add(Candidate{Text: fn, Type: CandidateFunction})
	}
}

// expandTypes adds common PostgreSQL data types as candidates.
func expandTypes(add func(Candidate)) {
	for _, t := range builtinTypes {
		add(Candidate{Text: t, Type: CandidateKeyword})
	}
}

// builtinFunctions is a list of common PostgreSQL built-in functions.
var builtinFunctions = []string{
	// Aggregate
	"avg", "count", "max", "min", "sum",
	"array_agg", "string_agg", "bool_and", "bool_or",
	// String
	"concat", "length", "lower", "upper", "trim", "substring",
	"replace", "position", "left", "right", "lpad", "rpad",
	"split_part", "regexp_replace", "regexp_match",
	// Math
	"abs", "ceil", "floor", "round", "trunc", "mod", "power", "sqrt",
	"random", "greatest", "least",
	// Date/Time
	"now", "current_date", "current_timestamp",
	"date_trunc", "date_part", "extract", "age",
	"to_char", "to_date", "to_timestamp", "to_number",
	// Type conversion
	"cast", "coalesce", "nullif",
	// JSON
	"json_build_object", "json_build_array", "json_agg",
	"jsonb_build_object", "jsonb_build_array", "jsonb_agg",
	"json_extract_path", "jsonb_extract_path",
	"json_extract_path_text", "jsonb_extract_path_text",
	// Array
	"array_length", "array_append", "array_cat", "unnest",
	// Window
	"row_number", "rank", "dense_rank", "ntile",
	"lag", "lead", "first_value", "last_value",
	// Conditional
	"case", "exists",
	// System
	"pg_typeof", "current_user", "current_schema",
	"gen_random_uuid",
}

// builtinTypes is a list of common PostgreSQL data types.
var builtinTypes = []string{
	"bigint", "bigserial", "bit", "boolean", "box",
	"bytea", "char", "character", "cidr", "circle",
	"date", "decimal", "double precision", "float4", "float8",
	"inet", "int", "int2", "int4", "int8",
	"integer", "interval", "json", "jsonb", "line",
	"lseg", "macaddr", "money", "numeric", "oid",
	"path", "point", "polygon", "real", "serial",
	"smallint", "smallserial", "text", "time", "timestamp",
	"timestamptz", "timetz", "tsquery", "tsvector", "uuid",
	"varchar", "xml",
}

// expandRelationsMinimal adds schemas and search-path tables (no sequences, no cross-schema).
func expandRelationsMinimal(ctx context.Context, meta MetadataProvider, refs []tableRef, ctes []cteInfo, activeSchema string, add func(Candidate)) {
	if meta == nil {
		return
	}

	schemas := meta.GetSchemaNames(ctx)
	defaultSchema := "public"

	for _, s := range schemas {
		add(Candidate{Text: QuoteIdentifier(s), Type: CandidateSchema})
	}

	if activeSchema != "" && activeSchema != defaultSchema {
		for _, t := range meta.GetTables(ctx, activeSchema) {
			add(Candidate{Text: QuoteIdentifier(t.Name), Type: CandidateTable})
		}
		for _, t := range meta.GetTables(ctx, defaultSchema) {
			add(Candidate{Text: QuoteIdentifier(t.Name), Type: CandidateTable})
		}
		for _, v := range meta.GetViews(ctx, defaultSchema) {
			add(Candidate{Text: QuoteIdentifier(v), Type: CandidateView})
		}
	} else {
		for _, t := range meta.GetTables(ctx, defaultSchema) {
			add(Candidate{Text: QuoteIdentifier(t.Name), Type: CandidateTable})
		}
		for _, v := range meta.GetViews(ctx, defaultSchema) {
			add(Candidate{Text: QuoteIdentifier(v), Type: CandidateView})
		}
	}

	// Add CTE names
	for _, cte := range ctes {
		add(Candidate{Text: QuoteIdentifier(cte.name), Type: CandidateTable})
	}

	// Add subquery aliases as tables
	for _, ref := range refs {
		if ref.isSubq && ref.alias != "" {
			add(Candidate{Text: QuoteIdentifier(ref.alias), Type: CandidateTable})
		} else if ref.alias != "" && ref.alias != ref.table {
			add(Candidate{Text: QuoteIdentifier(ref.alias), Type: CandidateTable})
		}
	}
}

// inferActiveSchema determines which non-default schema is being used.
func inferActiveSchema(refs []tableRef) string {
	for _, ref := range refs {
		if ref.schema != "" && ref.schema != "public" {
			return ref.schema
		}
	}
	return ""
}

// resolveSelectColumns resolves column names from a CTE's SelectStmt.
func resolveSelectColumns(ctx context.Context, sel *nodes.SelectStmt, meta MetadataProvider) []string {
	if sel == nil || sel.TargetList == nil {
		return nil
	}
	var cols []string
	for _, item := range sel.TargetList.Items {
		rt, ok := item.(*nodes.ResTarget)
		if !ok {
			continue
		}
		if rt.Name != "" {
			cols = append(cols, rt.Name)
			continue
		}
		if cr, ok := rt.Val.(*nodes.ColumnRef); ok {
			name := columnRefLastName(cr)
			if name != "" {
				cols = append(cols, name)
				continue
			}
		}
		// SELECT * — expand
		if sel.FromClause != nil {
			fromRefs := extractFromClauseRefs(sel.FromClause)
			for _, ref := range fromRefs {
				schema := ref.schema
				if schema == "" {
					schema = "public"
				}
				for _, col := range meta.GetColumns(ctx, schema, ref.table) {
					cols = append(cols, col.Name)
				}
			}
		}
	}
	return cols
}

func columnRefLastName(cr *nodes.ColumnRef) string {
	if cr.Fields == nil || len(cr.Fields.Items) == 0 {
		return ""
	}
	last := cr.Fields.Items[len(cr.Fields.Items)-1]
	if s, ok := last.(*nodes.String); ok {
		return s.Str
	}
	return ""
}

func removeCursorMarker(sql string) string {
	idx := findCursorMarker(sql)
	if idx >= 0 {
		return sql[:idx] + sql[idx+1:]
	}
	return sql
}

// sortCandidates sorts candidates by Type then Text.
func sortCandidates(candidates []Candidate) []Candidate {
	typeOrder := map[CandidateType]int{
		CandidateColumn:   0,
		CandidateSchema:   1,
		CandidateSequence: 2,
		CandidateTable:    3,
		CandidateView:     4,
		CandidateFunction: 5,
		CandidateKeyword:  6,
	}

	sort.Slice(candidates, func(i, j int) bool {
		oi := typeOrder[candidates[i].Type]
		oj := typeOrder[candidates[j].Type]
		if oi != oj {
			return oi < oj
		}
		return candidates[i].Text < candidates[j].Text
	})
	return candidates
}
