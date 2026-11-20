package completion

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/mssql/catalog"
	"github.com/bytebase/omni/mssql/parser"
)

// tsqlTypes are the hardcoded T-SQL type keywords returned for the "type_name" rule.
var tsqlTypes = []string{
	"INT", "BIGINT", "SMALLINT", "TINYINT",
	"VARCHAR", "NVARCHAR", "CHAR", "NCHAR",
	"TEXT", "NTEXT",
	"DATETIME", "DATETIME2", "DATE", "TIME",
	"DECIMAL", "NUMERIC", "FLOAT", "REAL",
	"BIT",
	"MONEY", "SMALLMONEY",
	"UNIQUEIDENTIFIER",
	"XML",
	"VARBINARY", "IMAGE",
	"SQL_VARIANT",
}

var tsqlTableHints = []string{
	"NOLOCK", "READUNCOMMITTED", "READCOMMITTED", "REPEATABLEREAD",
	"SERIALIZABLE", "HOLDLOCK", "UPDLOCK", "TABLOCK", "TABLOCKX",
	"ROWLOCK", "PAGLOCK", "INDEX", "FORCESEEK",
}

var tsqlQueryHints = []string{
	"RECOMPILE", "OPTIMIZE", "MAXDOP", "HASH JOIN", "MERGE JOIN",
	"LOOP JOIN", "FORCE ORDER",
}

// builtinFunctions lists T-SQL built-in scalar, aggregate, and metadata functions.
var builtinFunctions = []string{
	"ABS", "ACOS", "APP_NAME", "APPLOCK_MODE", "APPLOCK_TEST", "ASCII", "ASIN",
	"ASSEMBLYPROPERTY", "ATAN", "ATN2", "BINARY_CHECKSUM", "CAST", "CEILING",
	"CERT_ID", "CERTENCODED", "CERTPRIVATEKEY", "CHAR", "CHARINDEX", "CHECKSUM",
	"COALESCE", "COL_LENGTH", "COL_NAME", "COLUMNPROPERTY", "COMPRESS", "CONCAT",
	"CONCAT_WS", "CONNECTIONPROPERTY", "CONTEXT_INFO", "CONVERT", "COS", "COT",
	"CURRENT_DATE", "CURRENT_REQUEST_ID", "CURRENT_TIMESTAMP", "CURRENT_TIMEZONE",
	"CURRENT_TIMEZONE_ID", "CURRENT_TRANSACTION_ID", "CURRENT_USER", "CURSOR_ROWS",
	"CURSOR_STATUS", "DATABASE_PRINCIPAL_ID", "DATABASEPROPERTYEX", "DATALENGTH",
	"DATE_BUCKET", "DATEADD", "DATEDIFF", "DATEDIFF_BIG", "DATEFROMPARTS",
	"DATENAME", "DATEPART", "DATETIME2FROMPARTS", "DATETIMEFROMPARTS",
	"DATETIMEOFFSETFROMPARTS", "DATETRUNC", "DAY", "DB_ID", "DB_NAME",
	"DECOMPRESS", "DEGREES", "DIFFERENCE", "EOMONTH", "ERROR_LINE",
	"ERROR_MESSAGE", "ERROR_NUMBER", "ERROR_PROCEDURE", "ERROR_SEVERITY",
	"ERROR_STATE", "EXP", "FETCH_STATUS", "FILE_ID", "FILE_IDEX", "FILE_NAME",
	"FILEGROUP_ID", "FILEGROUP_NAME", "FILEGROUPPROPERTY", "FILEPROPERTY",
	"FILEPROPERTYEX", "FLOOR", "FORMAT", "FORMATMESSAGE",
	"FULLTEXTCATALOGPROPERTY", "FULLTEXTSERVICEPROPERTY",
	"GET_FILESTREAM_TRANSACTION_CONTEXT", "GETANSINULL", "GETDATE", "GETUTCDATE",
	"GREATEST", "HAS_DBACCESS", "HAS_PERMS_BY_NAME", "HOST_ID", "HOST_NAME",
	"IDENT_CURRENT", "IDENT_INCR", "IDENT_SEED", "IDENTITY", "IIF", "INDEX_COL",
	"INDEXKEY_PROPERTY", "INDEXPROPERTY", "IS_MEMBER", "IS_ROLEMEMBER",
	"IS_SRVROLEMEMBER", "ISDATE", "ISJSON", "ISNULL", "ISNUMERIC", "JSON_ARRAY",
	"JSON_MODIFY", "JSON_OBJECT", "JSON_PATH_EXISTS", "JSON_QUERY", "JSON_VALUE",
	"LEAST", "LEFT", "LEN", "LOG", "LOG10", "LOGINPROPERTY", "LOWER", "LTRIM",
	"MATH_SIGN", "MIN_ACTIVE_ROWVERSION", "MONTH", "NCHAR", "NEWID",
	"NEWSEQUENTIALID", "NEXT_VALUE_FOR", "NULLIF", "OBJECT_DEFINITION",
	"OBJECT_ID", "OBJECT_NAME", "OBJECT_SCHEMA_NAME", "OBJECTPROPERTY",
	"OBJECTPROPERTYEX", "ORIGINAL_DB_NAME", "ORIGINAL_LOGIN", "PARSE",
	"PARSENAME", "PATINDEX", "PERMISSIONS", "PI", "POWER", "PWDCOMPARE",
	"PWDENCRYPT", "QUOTENAME", "RADIANS", "RAND", "REPLACE", "REPLICATE",
	"REVERSE", "RIGHT", "ROUND", "ROWCOUNT_BIG", "RTRIM", "SCHEMA_ID",
	"SCHEMA_NAME", "SCOPE_IDENTITY", "SERVERPROPERTY", "SESSION_CONTEXT",
	"SESSIONPROPERTY", "SIN", "SMALLDATETIMEFROMPARTS", "SOUNDEX", "SPACE",
	"SQL_VARIANT_PROPERTY", "SQRT", "SQUARE", "STATS_DATE", "STR",
	"STRING_ESCAPE", "STRINGAGG", "STUFF", "SUBSTRING", "SUSER_ID", "SUSER_NAME",
	"SUSER_SID", "SUSER_SNAME", "SWITCHOFFSET", "SYSDATETIME",
	"SYSDATETIMEOFFSET", "SYSTEM_USER", "SYSUTCDATETIME", "TAN", "TIMEFROMPARTS",
	"TODATETIMEOFFSET", "TRANSLATE", "TRIM", "TRY_CAST", "TYPE_ID", "TYPE_NAME",
	"TYPEPROPERTY", "UNICODE", "UPPER", "USER", "USER_ID", "USER_NAME",
	"XACT_STATE", "XML_DATA_TYPE_FUNC", "YEAR",
}

const asciiWhitespace = " \t\r\n"

// resolve converts parser CandidateSet into completion Candidates.
// cat may be nil; in that case only keyword and hardcoded type candidates are returned.
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

	result = append(result, contextualKeywordCandidates(sql, cursorOffset)...)

	if cat == nil {
		for _, rc := range cs.Rules {
			result = append(result, resolveRule(rc.Rule, nil, sql, cursorOffset)...)
		}
		if ctx := parser.CollectCompletion(sql, cursorOffset); ctx != nil && ctx.Intent != nil {
			for _, kind := range ctx.Intent.ObjectKinds {
				if kind == parser.ObjectKindType {
					result = append(result, resolveTypeNames()...)
				}
			}
		}
		return dedup(result)
	}

	r := newCatalogResolver(cat, sql, cursorOffset)
	result = append(result, r.resolveCandidates(cs)...)
	return dedup(result)
}

// resolveRule resolves a single grammar rule into completion candidates.
func resolveRule(rule string, cat *catalog.Catalog, sql string, cursorOffset int) []Candidate {
	// type_name always returns hardcoded T-SQL types regardless of catalog.
	if rule == "type_name" {
		return resolveTypeNames()
	}

	// pseudocolumn_action suggests the $action pseudo-column (OUTPUT clause);
	// static, no catalog needed.
	if rule == "pseudocolumn_action" {
		return []Candidate{{Text: "$action", Type: CandidateColumn}}
	}

	// All other rules require a catalog.
	if cat == nil {
		return nil
	}

	switch rule {
	case "table_ref", "columnref", "schema_ref", "func_name",
		"proc_ref", "proc_name", "view_ref", "view_name",
		"index_ref", "trigger_ref", "database_ref", "sequence_ref":
		r := newCatalogResolver(cat, sql, cursorOffset)
		return dedup(r.resolveCandidates(&parser.CandidateSet{
			Rules: []parser.RuleCandidate{{Rule: rule}},
		}))
	case "login_ref", "user_ref", "role_ref":
		return nil
	}
	return nil
}

// resolveTypeNames returns hardcoded T-SQL type keyword candidates.
func resolveTypeNames() []Candidate {
	result := make([]Candidate, 0, len(tsqlTypes))
	for _, t := range tsqlTypes {
		result = append(result, Candidate{Text: t, Type: CandidateType_})
	}
	return result
}

func resolveBuiltinFunctions() []Candidate {
	result := make([]Candidate, 0, len(builtinFunctions))
	for _, fn := range builtinFunctions {
		result = append(result, Candidate{Text: fn, Type: CandidateFunction})
	}
	return result
}

func contextualKeywordCandidates(sql string, cursorOffset int) []Candidate {
	if cursorOffset < 0 || len(sql) == 0 {
		return nil
	}
	before := strings.ToUpper(strings.TrimSpace(sql[:min(cursorOffset, len(sql))]))
	if before == "" {
		return nil
	}
	var result []Candidate
	appendKeywords := func(keywords ...string) {
		for _, kw := range keywords {
			result = append(result, Candidate{Text: kw, Type: CandidateKeyword})
		}
	}
	if hasUnclosedKeywordParen(before, "WITH") {
		appendKeywords(tsqlTableHints...)
	}
	if hasUnclosedKeywordParen(before, "OPTION") {
		appendKeywords(tsqlQueryHints...)
	}
	if strings.HasSuffix(before, "FOR XML") {
		appendKeywords("PATH", "RAW", "AUTO", "EXPLICIT")
	}
	if strings.HasSuffix(before, "FOR JSON") {
		appendKeywords("PATH", "AUTO")
	}
	if strings.HasSuffix(before, " WHEN") || before == "WHEN" {
		appendKeywords("MATCHED", "NOT")
	}
	return result
}

func hasUnclosedKeywordParen(before, keyword string) bool {
	start := strings.LastIndex(before, keyword+" (")
	if start < 0 {
		return false
	}
	return strings.LastIndex(before[start:], ")") < 0
}

type physicalTableRef struct {
	Server   string
	Database string
	Schema   string
	Table    string
	Alias    string
}

type virtualTableRef struct {
	Table   string
	Columns []string
}

type resolvedTableRef struct {
	physical *physicalTableRef
	virtual  *virtualTableRef
}

type catalogResolver struct {
	cat              *catalog.Catalog
	sql              string
	cursorByteOffset int
	tokens           []parser.Token
	caretTokenIndex  int
	completionCtx    *parser.CompletionContext
	completionIntent *parser.CompletionIntent
	references       []resolvedTableRef
	cteTables        []virtualTableRef
}

func newCatalogResolver(cat *catalog.Catalog, sql string, cursorOffset int) *catalogResolver {
	if cursorOffset < 0 {
		cursorOffset = 0
	}
	if cursorOffset > len(sql) {
		cursorOffset = len(sql)
	}
	tokens := parser.Tokenize(sql)
	r := &catalogResolver{
		cat:              cat,
		sql:              sql,
		cursorByteOffset: cursorOffset,
		tokens:           tokens,
		caretTokenIndex:  findCaretTokenIndex(tokens, cursorOffset),
	}
	r.completionCtx = parser.CollectCompletion(sql, cursorOffset)
	if r.completionCtx != nil {
		r.completionIntent = r.completionCtx.Intent
		r.collectCompletionCTEs(r.completionCtx)
		r.collectCompletionScopeReferences(r.completionCtx)
	}
	if len(r.references) == 0 {
		for _, ref := range extractTableRefs(sql, cursorOffset) {
			r.appendPhysicalRef(physicalTableRef{
				Server:   unquoteIdent(ref.Server),
				Database: unquoteIdent(ref.Database),
				Schema:   unquoteIdent(ref.Schema),
				Table:    unquoteIdent(ref.Table),
				Alias:    unquoteIdent(ref.Alias),
			})
		}
	}
	return r
}

func findCaretTokenIndex(tokens []parser.Token, byteOffset int) int {
	for i, tok := range tokens {
		if tok.Loc >= byteOffset {
			return i
		}
	}
	return len(tokens)
}

func (r *catalogResolver) defaultDatabase() string {
	if r.cat == nil {
		return ""
	}
	return r.cat.DatabaseName()
}

func (r *catalogResolver) defaultSchema() string {
	if r.cat == nil || r.cat.DefaultSchema() == "" {
		return "dbo"
	}
	return r.cat.DefaultSchema()
}

func (r *catalogResolver) resolveCandidates(cs *parser.CandidateSet) []Candidate {
	var (
		typeEntries     []Candidate
		functionEntries []Candidate
		databaseEntries []Candidate
		schemaEntries   []Candidate
		tableEntries    []Candidate
		columnEntries   []Candidate
		viewEntries     []Candidate
		sequenceEntries []Candidate
		routineEntries  []Candidate
		indexEntries    []Candidate
		triggerEntries  []Candidate
	)

	handledDotted := r.insertDottedObjectContextCandidates(&databaseEntries, &schemaEntries, &tableEntries, &viewEntries, &sequenceEntries, &routineEntries)
	if !handledDotted {
		r.insertCompletionIntentCandidates(&typeEntries, &functionEntries, &databaseEntries, &schemaEntries, &tableEntries, &viewEntries, &sequenceEntries, &routineEntries)
	}

	if cs != nil {
		for _, rc := range cs.Rules {
			if handledDotted && isObjectCompletionRule(rc.Rule) {
				continue
			}
			switch rc.Rule {
			case "type_name":
				typeEntries = append(typeEntries, resolveTypeNames()...)
			case "func_name":
				functionEntries = append(functionEntries, resolveBuiltinFunctions()...)
				functionEntries = append(functionEntries, r.metadataFunctions("", "", "")...)
			case "database_ref":
				databaseEntries = append(databaseEntries, r.metadataDatabases("")...)
			case "schema_ref":
				schemaEntries = append(schemaEntries, r.metadataSchemas("", "")...)
			case "proc_ref", "proc_name":
				for _, ctx := range r.determineObjectNameContext() {
					if ctx.flags&objectFlagShowObject != 0 {
						routineEntries = append(routineEntries, r.metadataProcedures(ctx.linkedServer, ctx.database, ctx.schema)...)
					}
				}
			case "sequence_ref":
				for _, ctx := range r.determineObjectNameContext() {
					if ctx.flags&objectFlagShowObject != 0 {
						sequenceEntries = append(sequenceEntries, r.metadataSequences(ctx.linkedServer, ctx.database, ctx.schema)...)
					}
				}
			case "table_ref", "asterisk":
				for _, ctx := range r.determineObjectNameContext() {
					r.insertObjectCandidates(ctx, &databaseEntries, &schemaEntries, &tableEntries, &viewEntries, &sequenceEntries)
				}
			case "view_name", "view_ref":
				for _, ctx := range r.determineObjectNameContext() {
					if ctx.flags&objectFlagShowObject != 0 {
						viewEntries = append(viewEntries, r.metadataViews(ctx.linkedServer, ctx.database, ctx.schema)...)
					}
				}
			case "index_ref":
				for _, ctx := range r.determineObjectNameContext() {
					if ctx.flags&objectFlagShowObject != 0 {
						indexEntries = append(indexEntries, r.metadataIndexes(ctx.linkedServer, ctx.database, ctx.schema)...)
					}
				}
			case "trigger_ref":
				for _, ctx := range r.determineObjectNameContext() {
					if ctx.flags&objectFlagShowObject != 0 {
						triggerEntries = append(triggerEntries, r.metadataTriggers(ctx.linkedServer, ctx.database, ctx.schema)...)
					}
				}
			case "pseudocolumn_action":
				columnEntries = append(columnEntries, Candidate{Text: "$action", Type: CandidateColumn})
			case "columnref":
				r.resolveColumnRefCandidates(cs, &databaseEntries, &schemaEntries, &tableEntries, &viewEntries, &sequenceEntries, &columnEntries)
			}
		}
	}

	r.applyContextualColumnCandidates(&columnEntries)

	var result []Candidate
	result = append(result, typeEntries...)
	result = append(result, functionEntries...)
	result = append(result, databaseEntries...)
	result = append(result, schemaEntries...)
	result = append(result, tableEntries...)
	result = append(result, columnEntries...)
	result = append(result, viewEntries...)
	result = append(result, sequenceEntries...)
	result = append(result, routineEntries...)
	result = append(result, indexEntries...)
	result = append(result, triggerEntries...)
	return result
}

func (r *catalogResolver) resolveColumnRefCandidates(
	cs *parser.CandidateSet,
	databaseEntries, schemaEntries, tableEntries, viewEntries, sequenceEntries, columnEntries *[]Candidate,
) {
	for _, ctx := range r.determineColumnNameContext() {
		r.insertObjectCandidates(ctx, databaseEntries, schemaEntries, tableEntries, viewEntries, sequenceEntries)
		if ctx.flags&objectFlagShowObject != 0 {
			for _, ref := range r.references {
				if p := ref.physical; p != nil {
					tableName := p.Table
					if p.Alias != "" {
						tableName = p.Alias
					}
					if ctx.linkedServer == "" && strings.EqualFold(p.Database, ctx.database) &&
						strings.EqualFold(p.Schema, ctx.schema) {
						*tableEntries = append(*tableEntries, Candidate{Text: tableName, Type: CandidateTable})
					}
				} else if v := ref.virtual; v != nil {
					if ctx.linkedServer == "" && ctx.database == "" && ctx.schema == "" {
						*tableEntries = append(*tableEntries, Candidate{Text: v.Table, Type: CandidateTable})
					}
				}
			}
		}
		if ctx.flags&objectFlagShowColumn != 0 {
			if aliasStartOffset, ok := r.selectItemAliasStartOffset(); ok && cs != nil {
				for _, alias := range r.fetchSelectItemAliases(cs.SelectAliasPositions, aliasStartOffset) {
					*columnEntries = append(*columnEntries, Candidate{Text: alias, Type: CandidateColumn})
				}
			}
			if !ctx.empty() || len(r.references) == 0 {
				*columnEntries = append(*columnEntries, r.metadataColumns(ctx.linkedServer, ctx.database, ctx.schema, ctx.object)...)
			}
			for _, ref := range r.references {
				matchedQualified := false
				if p := ref.physical; p != nil {
					if ctx.empty() {
						*columnEntries = append(*columnEntries, r.metadataColumns("", p.Database, p.Schema, p.Table)...)
						continue
					}
					inputLinkedServer := ctx.linkedServer
					inputDB := ctx.database
					if inputDB == "" {
						inputDB = r.defaultDatabase()
					}
					inputSchema := ctx.schema
					if inputSchema == "" {
						inputSchema = r.defaultSchema()
					}
					inputTable := ctx.object

					refDB := p.Database
					if refDB == "" {
						refDB = r.defaultDatabase()
					}
					refSchema := p.Schema
					if refSchema == "" {
						refSchema = r.defaultSchema()
					}
					refTable := p.Table
					refIsAlias := false
					if p.Alias != "" {
						refTable = p.Alias
						refIsAlias = true
					}

					if inputLinkedServer == "" && strings.EqualFold(refTable, inputTable) &&
						(refIsAlias || (strings.EqualFold(refDB, inputDB) && strings.EqualFold(refSchema, inputSchema))) {
						*columnEntries = append(*columnEntries, r.metadataColumns("", p.Database, p.Schema, p.Table)...)
						matchedQualified = refIsAlias && ctx.database == "" && ctx.schema == ""
					}
				} else if v := ref.virtual; v != nil {
					if ctx.linkedServer == "" && ctx.database == "" && ctx.schema == "" {
						if ctx.object != "" && !strings.EqualFold(v.Table, ctx.object) {
							continue
						}
						for _, col := range v.Columns {
							*columnEntries = append(*columnEntries, Candidate{Text: col, Type: CandidateColumn})
						}
						matchedQualified = ctx.object != ""
					}
				}
				if matchedQualified {
					break
				}
			}
			if ctx.linkedServer == "" && ctx.database == "" && ctx.schema == "" {
				for _, cte := range r.cteTables {
					if strings.EqualFold(cte.Table, ctx.object) {
						for _, col := range cte.Columns {
							*columnEntries = append(*columnEntries, Candidate{Text: col, Type: CandidateColumn})
						}
					}
				}
			}
			if ctx.empty() && len(r.references) == 0 {
				*columnEntries = append(*columnEntries, r.allColumns()...)
			}
		}
	}
}

func (r *catalogResolver) insertCompletionIntentCandidates(
	typeEntries, functionEntries, databaseEntries, schemaEntries, tableEntries, viewEntries, sequenceEntries, routineEntries *[]Candidate,
) {
	if r.completionIntent == nil {
		return
	}
	contexts := r.determineObjectNameContext()
	for _, kind := range r.completionIntent.ObjectKinds {
		for _, ctx := range contexts {
			switch kind {
			case parser.ObjectKindDatabase:
				*databaseEntries = append(*databaseEntries, r.metadataDatabases(ctx.linkedServer)...)
			case parser.ObjectKindSchema:
				*schemaEntries = append(*schemaEntries, r.metadataSchemas(ctx.linkedServer, ctx.database)...)
			case parser.ObjectKindTable:
				if ctx.flags&objectFlagShowObject != 0 {
					*tableEntries = append(*tableEntries, r.metadataTables(ctx.linkedServer, ctx.database, ctx.schema)...)
				}
			case parser.ObjectKindView:
				if ctx.flags&objectFlagShowObject != 0 {
					*viewEntries = append(*viewEntries, r.metadataViews(ctx.linkedServer, ctx.database, ctx.schema)...)
				}
			case parser.ObjectKindSequence:
				if ctx.flags&objectFlagShowObject != 0 {
					*sequenceEntries = append(*sequenceEntries, r.metadataSequences(ctx.linkedServer, ctx.database, ctx.schema)...)
				}
			case parser.ObjectKindProcedure:
				if ctx.flags&objectFlagShowObject != 0 {
					*routineEntries = append(*routineEntries, r.metadataProcedures(ctx.linkedServer, ctx.database, ctx.schema)...)
				}
			case parser.ObjectKindFunction:
				*functionEntries = append(*functionEntries, resolveBuiltinFunctions()...)
				*functionEntries = append(*functionEntries, r.metadataFunctions(ctx.linkedServer, ctx.database, ctx.schema)...)
			case parser.ObjectKindType:
				*typeEntries = append(*typeEntries, resolveTypeNames()...)
			}
		}
	}
}

func (r *catalogResolver) insertObjectCandidates(
	ctx *objectRefContext,
	databaseEntries, schemaEntries, tableEntries, viewEntries, sequenceEntries *[]Candidate,
) {
	if ctx.flags&objectFlagShowDatabase != 0 {
		*databaseEntries = append(*databaseEntries, r.metadataDatabases(ctx.linkedServer)...)
	}
	if ctx.flags&objectFlagShowSchema != 0 {
		*schemaEntries = append(*schemaEntries, r.metadataSchemas(ctx.linkedServer, ctx.database)...)
	}
	if ctx.flags&objectFlagShowObject == 0 {
		return
	}
	*tableEntries = append(*tableEntries, r.metadataTables(ctx.linkedServer, ctx.database, ctx.schema)...)
	*viewEntries = append(*viewEntries, r.metadataViews(ctx.linkedServer, ctx.database, ctx.schema)...)
	*sequenceEntries = append(*sequenceEntries, r.metadataSequences(ctx.linkedServer, ctx.database, ctx.schema)...)
	if ctx.linkedServer == "" && ctx.database == "" && ctx.schema == "" {
		for _, cte := range r.cteTables {
			*tableEntries = append(*tableEntries, Candidate{Text: cte.Table, Type: CandidateTable})
		}
	}
}

func (r *catalogResolver) insertDottedObjectContextCandidates(
	databaseEntries, schemaEntries, tableEntries, viewEntries, sequenceEntries, routineEntries *[]Candidate,
) bool {
	if len(*databaseEntries) > 0 || len(*schemaEntries) > 0 || len(*tableEntries) > 0 ||
		len(*viewEntries) > 0 || len(*sequenceEntries) > 0 || len(*routineEntries) > 0 {
		return false
	}
	parts, beforeIdx, ok := r.multipartPartsBeforeDottedObject()
	if !ok {
		return false
	}
	kind, ok := r.dottedObjectKindBefore(beforeIdx)
	if !ok {
		return false
	}
	r.insertDottedObjectCandidatesForParts(kind, parts, schemaEntries, tableEntries, viewEntries, sequenceEntries, routineEntries)
	return true
}

func isObjectCompletionRule(rule string) bool {
	switch rule {
	case "database_ref", "schema_ref", "table_ref", "view_name", "view_ref", "sequence_ref", "proc_ref", "proc_name":
		return true
	default:
		return false
	}
}

func (r *catalogResolver) multipartPartsBeforeDottedObject() ([]string, int, bool) {
	idx := r.previousTokenIndex()
	if idx < 1 {
		return nil, -1, false
	}
	if r.completionCtx != nil && r.completionCtx.Prefix != "" && parser.IsIdentTokenType(r.tokens[idx].Type) {
		idx--
	}
	if idx < 1 || !isMultipartSeparator(r.tokenText(idx)) {
		return nil, -1, false
	}

	var parts []string
	for {
		separator := r.tokenText(idx)
		partIdx := idx - 1
		if partIdx < 0 {
			return nil, -1, false
		}
		if r.tokenText(partIdx) == "." {
			if separator == ".." {
				return nil, -1, false
			}
			parts = append([]string{""}, parts...)
			idx = partIdx
			continue
		}
		tok := r.tokens[partIdx]
		if !parser.IsIdentTokenType(tok.Type) {
			return nil, -1, false
		}
		part := unquoteIdent(r.tokenText(partIdx))
		if part == "" {
			return nil, -1, false
		}
		if separator == ".." {
			parts = append([]string{part, ""}, parts...)
		} else {
			parts = append([]string{part}, parts...)
		}
		if partIdx-1 < 0 || r.tokenText(partIdx-1) != "." {
			return parts, partIdx - 1, true
		}
		idx = partIdx - 1
	}
}

func isMultipartSeparator(text string) bool {
	return text == "." || text == ".."
}

type dottedObjectKind int

const (
	dottedObjectKindRelation dottedObjectKind = iota
	dottedObjectKindView
	dottedObjectKindSequence
	dottedObjectKindRoutine
)

func (r *catalogResolver) dottedObjectKindBefore(idx int) (dottedObjectKind, bool) {
	if idx < 0 {
		return dottedObjectKindRelation, false
	}
	if strings.EqualFold(r.tokenText(idx), "EXISTS") && idx >= 1 && strings.EqualFold(r.tokenText(idx-1), "IF") {
		idx -= 2
	}
	if idx < 0 {
		return dottedObjectKindRelation, false
	}
	switch strings.ToUpper(r.tokenText(idx)) {
	case "VIEW":
		return dottedObjectKindView, true
	case "SEQUENCE":
		return dottedObjectKindSequence, true
	case "PROCEDURE", "PROC", "EXEC", "EXECUTE":
		return dottedObjectKindRoutine, true
	case "FOR":
		if idx >= 2 && strings.EqualFold(r.tokenText(idx-1), "VALUE") && strings.EqualFold(r.tokenText(idx-2), "NEXT") {
			return dottedObjectKindSequence, true
		}
		return dottedObjectKindRelation, false
	case "FROM", "JOIN", "APPLY", "INTO", "UPDATE", "USING", "REFERENCES", "TABLE", "TRUNCATE":
		return dottedObjectKindRelation, true
	default:
		return dottedObjectKindRelation, false
	}
}

func (r *catalogResolver) insertDottedObjectCandidatesForParts(
	kind dottedObjectKind,
	parts []string,
	schemaEntries, tableEntries, viewEntries, sequenceEntries, routineEntries *[]Candidate,
) {
	switch len(parts) {
	case 0:
		return
	case 1:
		*schemaEntries = append(*schemaEntries, r.metadataSchemas("", parts[0])...)
		r.insertDottedObjectCandidatesForSchema(kind, "", parts[0], tableEntries, viewEntries, sequenceEntries, routineEntries)
	case 2:
		if parts[1] == "" {
			r.insertDottedObjectCandidatesForSchema(kind, parts[0], r.defaultSchema(), tableEntries, viewEntries, sequenceEntries, routineEntries)
			return
		}
		r.insertDottedObjectCandidatesForSchema(kind, parts[len(parts)-2], parts[len(parts)-1], tableEntries, viewEntries, sequenceEntries, routineEntries)
	default:
		return
	}
}

func (r *catalogResolver) insertDottedObjectCandidatesForSchema(
	kind dottedObjectKind,
	database, schema string,
	tableEntries, viewEntries, sequenceEntries, routineEntries *[]Candidate,
) {
	switch kind {
	case dottedObjectKindView:
		*viewEntries = append(*viewEntries, r.metadataViews("", database, schema)...)
	case dottedObjectKindSequence:
		*sequenceEntries = append(*sequenceEntries, r.metadataSequences("", database, schema)...)
	case dottedObjectKindRoutine:
		*routineEntries = append(*routineEntries, r.metadataProcedures("", database, schema)...)
	default:
		*tableEntries = append(*tableEntries, r.metadataTables("", database, schema)...)
		*viewEntries = append(*viewEntries, r.metadataViews("", database, schema)...)
		*sequenceEntries = append(*sequenceEntries, r.metadataSequences("", database, schema)...)
	}
}

func (r *catalogResolver) applyContextualColumnCandidates(columnEntries *[]Candidate) {
	if database, schema, table, usedColumns, ok := r.insertColumnListContext(); ok {
		cols := r.metadataColumnsExcept("", database, schema, table, usedColumns)
		if len(cols) > 0 {
			*columnEntries = cols
		}
	}
	if database, schema, table, ok := r.createIndexColumnListContext(); ok {
		cols := r.metadataColumns("", database, schema, table)
		if len(cols) > 0 {
			*columnEntries = cols
		}
	}
	if cols := r.createTableForeignKeySourceColumns(); len(cols) > 0 {
		*columnEntries = make([]Candidate, 0, len(cols))
		for _, col := range cols {
			*columnEntries = append(*columnEntries, Candidate{Text: col, Type: CandidateColumn})
		}
	}
	if database, schema, table, ok := r.referencesColumnListContext(); ok {
		cols := r.metadataColumns("", database, schema, table)
		if len(cols) > 0 {
			*columnEntries = cols
		}
	}
	if alias, ok := r.qualifiedObjectBeforeCaret(); ok {
		*columnEntries = append(*columnEntries, r.derivedTableColumns(alias)...)
	}
}

func (r *catalogResolver) insertColumnListContext() (string, string, string, map[string]bool, bool) {
	before := r.sql[:min(r.cursorByteOffset, len(r.sql))]
	re := regexp.MustCompile(`(?is)\bINSERT\s+INTO\s+([^\s(]+)\s*\(([^()]*)$`)
	matches := re.FindAllStringSubmatch(before, -1)
	if len(matches) == 0 {
		return "", "", "", nil, false
	}
	database, schema, table := parseMultipartIdentifier(matches[len(matches)-1][1])
	if table == "" {
		return "", "", "", nil, false
	}
	excluded := make(map[string]bool)
	for _, col := range splitTopLevelCSV(matches[len(matches)-1][2]) {
		col = strings.TrimSpace(col)
		if col == "" {
			continue
		}
		excluded[strings.ToLower(unquoteIdent(col))] = true
	}
	return database, schema, table, excluded, true
}

func (r *catalogResolver) createIndexColumnListContext() (string, string, string, bool) {
	before := r.sql[:min(r.cursorByteOffset, len(r.sql))]
	re := regexp.MustCompile(`(?is)\bCREATE\s+(?:UNIQUE\s+)?(?:(?:CLUSTERED|NONCLUSTERED)\s+)?INDEX\s+[^\s]+\s+ON\s+([^\s(]+)\s*(?:\([^)]*\)\s*INCLUDE\s*)?\([^)]*$`)
	matches := re.FindAllStringSubmatch(before, -1)
	if len(matches) == 0 {
		return "", "", "", false
	}
	database, schema, table := parseMultipartIdentifier(matches[len(matches)-1][1])
	return database, schema, table, table != ""
}

func (r *catalogResolver) createTableForeignKeySourceColumns() []string {
	before := r.sql[:min(r.cursorByteOffset, len(r.sql))]
	re := regexp.MustCompile(`(?is)\bCREATE\s+TABLE\s+[^\s(]+\s*\((.*)\bFOREIGN\s+KEY\s*\([^)]*$`)
	matches := re.FindAllStringSubmatch(before, -1)
	if len(matches) == 0 {
		return nil
	}
	var columns []string
	for _, item := range splitTopLevelCSV(matches[len(matches)-1][1]) {
		fields := strings.Fields(strings.TrimSpace(item))
		if len(fields) < 2 {
			continue
		}
		kw := strings.ToUpper(fields[0])
		if kw == "CONSTRAINT" || kw == "PRIMARY" || kw == "FOREIGN" || kw == "CHECK" || kw == "UNIQUE" {
			continue
		}
		columns = append(columns, unquoteIdent(fields[0]))
	}
	return columns
}

func (r *catalogResolver) referencesColumnListContext() (string, string, string, bool) {
	before := r.sql[:min(r.cursorByteOffset, len(r.sql))]
	re := regexp.MustCompile(`(?is)\bREFERENCES\s+([^\s(]+)\s*\([^)]*$`)
	matches := re.FindAllStringSubmatch(before, -1)
	if len(matches) == 0 {
		return "", "", "", false
	}
	database, schema, table := parseMultipartIdentifier(matches[len(matches)-1][1])
	return database, schema, table, table != ""
}

func (r *catalogResolver) qualifiedObjectBeforeCaret() (string, bool) {
	idx := r.previousTokenIndex()
	if idx < 1 || r.tokenText(idx) != "." {
		return "", false
	}
	name := unquoteIdent(r.tokenText(idx - 1))
	if name == "" {
		return "", false
	}
	return name, true
}

func (r *catalogResolver) derivedTableColumns(alias string) []Candidate {
	var result []Candidate
	for _, derived := range extractDerivedTables(r.sql) {
		if !strings.EqualFold(derived.alias, alias) {
			continue
		}
		for _, col := range r.inferQueryOutputColumns(derived.query) {
			result = append(result, Candidate{Text: col, Type: CandidateColumn})
		}
	}
	return result
}

// --- Catalog lookup helpers ---

func (r *catalogResolver) metadataDatabases(linkedServer string) []Candidate {
	if linkedServer != "" || r.cat == nil {
		return nil
	}
	dbs := r.cat.Databases()
	result := make([]Candidate, 0, len(dbs))
	for _, db := range dbs {
		if db != "" {
			result = append(result, Candidate{Text: db, Type: CandidateDatabase})
		}
	}
	return result
}

func (r *catalogResolver) metadataSchemas(linkedServer, database string) []Candidate {
	if linkedServer != "" || r.cat == nil {
		return nil
	}
	target := r.lookupDatabaseCatalog(database)
	if target == nil {
		return nil
	}
	names := target.ListSchemaNames()
	result := make([]Candidate, 0, len(names))
	for _, s := range names {
		result = append(result, Candidate{Text: s, Type: CandidateSchema})
	}
	return result
}

func (r *catalogResolver) lookupDatabaseCatalog(database string) *catalog.Catalog {
	if r.cat == nil {
		return nil
	}
	if database == "" {
		return r.cat
	}
	return r.cat.DatabaseCatalog(database)
}

func (r *catalogResolver) lookupSchema(linkedServer, database, schema string) *metadata.SchemaMetadata {
	if linkedServer != "" || r.cat == nil {
		return nil
	}
	target := r.lookupDatabaseCatalog(database)
	if target == nil {
		return nil
	}
	schemaName := schema
	if schemaName == "" {
		schemaName = target.DefaultSchema()
	}
	return target.GetSchema(schemaName)
}

func (r *catalogResolver) metadataTables(linkedServer, database, schema string) []Candidate {
	s := r.lookupSchema(linkedServer, database, schema)
	if s == nil {
		return nil
	}
	result := make([]Candidate, 0, len(s.Tables))
	for _, t := range s.Tables {
		result = append(result, Candidate{Text: t.Name, Type: CandidateTable})
	}
	return result
}

func (r *catalogResolver) metadataViews(linkedServer, database, schema string) []Candidate {
	s := r.lookupSchema(linkedServer, database, schema)
	if s == nil {
		return nil
	}
	var result []Candidate
	for _, v := range s.Views {
		result = append(result, Candidate{Text: v.Name, Type: CandidateView})
	}
	for _, mv := range s.MaterializedViews {
		result = append(result, Candidate{Text: mv.Name, Type: CandidateView})
	}
	for _, et := range s.ExternalTables {
		result = append(result, Candidate{Text: et.Name, Type: CandidateView})
	}
	return result
}

func (r *catalogResolver) metadataSequences(linkedServer, database, schema string) []Candidate {
	s := r.lookupSchema(linkedServer, database, schema)
	if s == nil {
		return nil
	}
	result := make([]Candidate, 0, len(s.Sequences))
	for _, seq := range s.Sequences {
		result = append(result, Candidate{Text: seq.Name, Type: CandidateSequence})
	}
	return result
}

func (r *catalogResolver) metadataProcedures(linkedServer, database, schema string) []Candidate {
	s := r.lookupSchema(linkedServer, database, schema)
	if s == nil {
		return nil
	}
	result := make([]Candidate, 0, len(s.Procedures))
	for _, p := range s.Procedures {
		result = append(result, Candidate{Text: p.Name, Type: CandidateProcedure})
	}
	return result
}

func (r *catalogResolver) metadataFunctions(linkedServer, database, schema string) []Candidate {
	s := r.lookupSchema(linkedServer, database, schema)
	if s == nil {
		return nil
	}
	result := make([]Candidate, 0, len(s.Functions))
	for _, fn := range s.Functions {
		result = append(result, Candidate{Text: fn.Name, Type: CandidateFunction})
	}
	return result
}

func (r *catalogResolver) metadataIndexes(linkedServer, database, schema string) []Candidate {
	s := r.lookupSchema(linkedServer, database, schema)
	if s == nil {
		return nil
	}
	var result []Candidate
	for _, t := range s.Tables {
		for _, idx := range t.Indexes {
			if idx.Name != "" {
				result = append(result, Candidate{Text: idx.Name, Type: CandidateIndex})
			}
		}
	}
	return result
}

func (r *catalogResolver) metadataTriggers(linkedServer, database, schema string) []Candidate {
	s := r.lookupSchema(linkedServer, database, schema)
	if s == nil {
		return nil
	}
	var result []Candidate
	for _, t := range s.Tables {
		for _, trig := range t.Triggers {
			if trig.Name != "" {
				result = append(result, Candidate{Text: trig.Name, Type: CandidateTrigger})
			}
		}
	}
	for _, v := range s.Views {
		for _, trig := range v.Triggers {
			if trig.Name != "" {
				result = append(result, Candidate{Text: trig.Name, Type: CandidateTrigger})
			}
		}
	}
	return result
}

func (r *catalogResolver) metadataColumns(linkedServer, database, schema, table string) []Candidate {
	if linkedServer != "" || r.cat == nil {
		return nil
	}
	target := r.lookupDatabaseCatalog(database)
	if target == nil {
		return nil
	}
	dbName := target.DatabaseName()
	schemaName := schema
	if schemaName == "" {
		schemaName = target.DefaultSchema()
	}
	s := target.GetSchema(schemaName)
	if s == nil {
		return nil
	}
	var result []Candidate
	for _, t := range s.Tables {
		if table != "" && !strings.EqualFold(t.Name, table) {
			continue
		}
		for _, col := range t.Columns {
			def := fmt.Sprintf("%s.%s.%s | %s", dbName, s.Name, t.Name, col.Type)
			if !col.Nullable {
				def += ", NOT NULL"
			}
			result = append(result, Candidate{
				Text:       col.Name,
				Type:       CandidateColumn,
				Definition: def,
				Comment:    col.Comment,
			})
		}
	}
	for _, v := range s.Views {
		if table != "" && !strings.EqualFold(v.Name, table) {
			continue
		}
		for _, col := range v.Columns {
			result = append(result, Candidate{
				Text:    col.Name,
				Type:    CandidateColumn,
				Comment: col.Comment,
			})
		}
	}
	return result
}

func (r *catalogResolver) metadataColumnsExcept(linkedServer, database, schema, table string, excluded map[string]bool) []Candidate {
	cols := r.metadataColumns(linkedServer, database, schema, table)
	if len(cols) == 0 || len(excluded) == 0 {
		return cols
	}
	var filtered []Candidate
	for _, c := range cols {
		if !excluded[strings.ToLower(unquoteIdent(c.Text))] {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

func (r *catalogResolver) allColumns() []Candidate {
	if r.cat == nil {
		return nil
	}
	dbName := r.defaultDatabase()
	var result []Candidate
	for _, s := range r.cat.Schemas("") {
		for _, t := range s.Tables {
			for _, col := range t.Columns {
				def := fmt.Sprintf("%s.%s.%s | %s", dbName, s.Name, t.Name, col.Type)
				if !col.Nullable {
					def += ", NOT NULL"
				}
				result = append(result, Candidate{
					Text:       col.Name,
					Type:       CandidateColumn,
					Definition: def,
					Comment:    col.Comment,
				})
			}
		}
	}
	return result
}

// --- Scope & CTE extraction ---

func (r *catalogResolver) collectCompletionCTEs(ctx *parser.CompletionContext) {
	if ctx == nil {
		return
	}
	for _, ref := range ctx.CTEs {
		v := r.convertVirtualReference(ref)
		if v == nil {
			continue
		}
		if len(v.Columns) == 0 {
			v.Columns = r.inferCTEColumns(ref)
		}
		r.appendCTEIfMissing(*v)
	}
}

func (r *catalogResolver) inferCTEColumns(ref parser.RangeReference) []string {
	if ref.BodyLoc.Start < 0 || ref.BodyLoc.End <= ref.BodyLoc.Start || ref.BodyLoc.End > len(r.sql) {
		return nil
	}
	body := strings.TrimSpace(r.sql[ref.BodyLoc.Start:ref.BodyLoc.End])
	if len(body) >= 2 && body[0] == '(' && body[len(body)-1] == ')' {
		body = strings.TrimSpace(body[1 : len(body)-1])
	}
	return r.inferQueryOutputColumns(body)
}

func (r *catalogResolver) inferQueryOutputColumns(querySQL string) []string {
	list, err := parser.Parse(querySQL)
	if err != nil || list == nil || len(list.Items) == 0 {
		return nil
	}
	stmt, ok := list.Items[0].(*ast.SelectStmt)
	if !ok || stmt == nil {
		return nil
	}
	return r.selectStmtOutputColumns(stmt)
}

func (r *catalogResolver) selectStmtOutputColumns(stmt *ast.SelectStmt) []string {
	if stmt == nil {
		return nil
	}
	if stmt.TargetList == nil {
		if stmt.Larg != nil {
			return r.selectStmtOutputColumns(stmt.Larg)
		}
		return nil
	}
	var cols []string
	for _, item := range stmt.TargetList.Items {
		switch target := item.(type) {
		case *ast.ResTarget:
			if target == nil {
				continue
			}
			if target.Name != "" {
				cols = append(cols, unquoteIdent(target.Name))
				continue
			}
			switch val := target.Val.(type) {
			case *ast.ColumnRef:
				if val == nil {
					continue
				}
				if val.Column == "*" {
					cols = append(cols, r.expandFromClauseStar(stmt, val.Table)...)
				} else if val.Column != "" {
					cols = append(cols, unquoteIdent(val.Column))
				}
			case *ast.StarExpr:
				if val != nil {
					cols = append(cols, r.expandFromClauseStar(stmt, val.Qualifier)...)
				}
			}
		case *ast.StarExpr:
			if target != nil {
				cols = append(cols, r.expandFromClauseStar(stmt, target.Qualifier)...)
			}
		}
	}
	return cols
}

func (r *catalogResolver) expandFromClauseStar(stmt *ast.SelectStmt, qualifier string) []string {
	if stmt == nil || stmt.FromClause == nil {
		return nil
	}
	refs := collectFromList(stmt.FromClause)
	var cols []string
	for _, ref := range refs {
		table := unquoteIdent(ref.Table)
		alias := unquoteIdent(ref.Alias)
		if qualifier != "" && !strings.EqualFold(qualifier, alias) && !strings.EqualFold(qualifier, table) {
			continue
		}
		if ref.Schema == "" && ref.Database == "" {
			if cteCols := r.lookupCTEColumns(table); len(cteCols) > 0 {
				cols = append(cols, cteCols...)
				continue
			}
		}
		for _, c := range r.metadataColumns(ref.Server, unquoteIdent(ref.Database), unquoteIdent(ref.Schema), table) {
			cols = append(cols, c.Text)
		}
	}
	return cols
}

func (r *catalogResolver) appendCTEIfMissing(v virtualTableRef) {
	for i := range r.cteTables {
		if strings.EqualFold(r.cteTables[i].Table, v.Table) {
			if len(r.cteTables[i].Columns) == 0 && len(v.Columns) > 0 {
				r.cteTables[i].Columns = v.Columns
			}
			return
		}
	}
	r.cteTables = append(r.cteTables, v)
}

func (r *catalogResolver) lookupCTEColumns(table string) []string {
	for _, cte := range r.cteTables {
		if strings.EqualFold(cte.Table, table) {
			return cte.Columns
		}
	}
	return nil
}

func (r *catalogResolver) collectCompletionScopeReferences(ctx *parser.CompletionContext) {
	if ctx == nil || ctx.Scope == nil {
		return
	}
	appendRef := func(ref parser.RangeReference) {
		r.appendConvertedScopeRef(ref)
	}
	if ctx.Scope.DMLTarget != nil {
		appendRef(*ctx.Scope.DMLTarget)
	}
	if ctx.Scope.MergeTarget != nil {
		appendRef(*ctx.Scope.MergeTarget)
	}
	if ctx.Scope.MergeSource != nil {
		appendRef(*ctx.Scope.MergeSource)
	}
	for _, ref := range ctx.Scope.LocalReferences {
		appendRef(ref)
	}
	for _, level := range ctx.Scope.OuterReferences {
		for _, ref := range level {
			appendRef(ref)
		}
	}
}

func (r *catalogResolver) appendConvertedScopeRef(ref parser.RangeReference) {
	switch ref.Kind {
	case parser.RangeReferenceRelation, parser.RangeReferenceDMLTarget, parser.RangeReferenceMergeTarget, parser.RangeReferenceMergeSource:
		alias := unquoteIdent(ref.Alias)
		table := unquoteIdent(ref.Object)
		if ref.Kind == parser.RangeReferenceMergeSource && ref.Schema == "" && ref.Database == "" {
			if cols := r.lookupCTEColumns(table); len(cols) > 0 {
				if alias != "" {
					table = alias
				}
				r.appendVirtualRef(virtualTableRef{Table: table, Columns: cols})
				return
			}
		}
		r.appendPhysicalRef(physicalTableRef{
			Server:   unquoteIdent(ref.Server),
			Database: unquoteIdent(ref.Database),
			Schema:   unquoteIdent(ref.Schema),
			Table:    table,
			Alias:    alias,
		})
	case parser.RangeReferenceCTE, parser.RangeReferenceValues, parser.RangeReferenceSubquery,
		parser.RangeReferenceTableVariable, parser.RangeReferenceJoinAlias, parser.RangeReferenceFunction:
		if v := r.convertVirtualReference(ref); v != nil {
			r.appendVirtualRef(*v)
		}
	}
}

func (r *catalogResolver) convertVirtualReference(ref parser.RangeReference) *virtualTableRef {
	table := unquoteIdent(ref.Alias)
	if table == "" {
		table = unquoteIdent(ref.Object)
	}
	if table == "" {
		return nil
	}
	rawCols := ref.Columns
	if len(rawCols) == 0 {
		rawCols = ref.AliasColumns
	}
	cols := make([]string, 0, len(rawCols))
	for _, c := range rawCols {
		cols = append(cols, unquoteIdent(c))
	}
	if len(cols) == 0 && ref.Kind == parser.RangeReferenceCTE {
		cols = r.lookupCTEColumns(table)
	}
	if len(cols) == 0 && ref.Kind == parser.RangeReferenceSubquery {
		cols = r.inferCTEColumns(ref)
	}
	return &virtualTableRef{
		Table:   table,
		Columns: cols,
	}
}

func (r *catalogResolver) appendPhysicalRef(p physicalTableRef) {
	for _, existing := range r.references {
		if ep := existing.physical; ep != nil {
			if ep.Server == p.Server && ep.Database == p.Database && ep.Schema == p.Schema && ep.Table == p.Table && ep.Alias == p.Alias {
				return
			}
		}
	}
	r.references = append(r.references, resolvedTableRef{physical: &p})
}

func (r *catalogResolver) appendVirtualRef(v virtualTableRef) {
	for _, existing := range r.references {
		if ev := existing.virtual; ev != nil && ev.Table == v.Table {
			return
		}
	}
	r.references = append(r.references, resolvedTableRef{virtual: &v})
}

// --- Object reference context ---

type objectFlag int

const (
	objectFlagShowLinkedServer objectFlag = 1 << iota
	objectFlagShowDatabase
	objectFlagShowSchema
	objectFlagShowObject
	objectFlagShowColumn
)

type objectRefContextOption func(*objectRefContext)

func withColumn() objectRefContextOption {
	return func(c *objectRefContext) {
		c.column = ""
		c.flags |= objectFlagShowColumn
	}
}

func newObjectRefContext(options ...objectRefContextOption) *objectRefContext {
	o := &objectRefContext{
		flags: objectFlagShowDatabase | objectFlagShowSchema | objectFlagShowObject,
	}
	for _, opt := range options {
		opt(o)
	}
	return o
}

type objectRefContext struct {
	linkedServer string
	database     string
	schema       string
	object       string
	column       string
	flags        objectFlag
}

func (o *objectRefContext) empty() bool {
	return o.linkedServer == "" && o.database == "" && o.schema == "" && o.object == "" && o.column == ""
}

func (o *objectRefContext) setLinkedServer(linkedServer string) {
	o.linkedServer = linkedServer
	o.flags &^= objectFlagShowLinkedServer
}

func (o *objectRefContext) setDatabase(database string) {
	o.database = database
	o.flags &^= objectFlagShowDatabase
}

func (o *objectRefContext) setSchema(schema string) {
	o.schema = schema
	o.flags &^= objectFlagShowSchema
}

func (o *objectRefContext) setObject(object string) {
	o.object = object
	o.flags &^= objectFlagShowObject
}

func (r *catalogResolver) determineObjectNameContext() []*objectRefContext {
	if r.completionIntent == nil {
		return []*objectRefContext{newObjectRefContext()}
	}
	ctx := newObjectRefContext()
	q := r.completionIntent.Qualifier
	if q.Server != "" {
		ctx.setLinkedServer(unquoteIdent(q.Server))
	}
	if q.Database != "" {
		ctx.setDatabase(unquoteIdent(q.Database))
	}
	if q.Schema != "" {
		ctx.setSchema(unquoteIdent(q.Schema))
		if q.Database == "" {
			ctx.flags &^= objectFlagShowDatabase
		}
	}
	if q.Object != "" {
		ctx.setObject(unquoteIdent(q.Object))
	}
	return []*objectRefContext{ctx}
}

func (r *catalogResolver) determineColumnNameContext() []*objectRefContext {
	if r.completionIntent == nil || !completionIntentHasObjectKind(r.completionIntent, parser.ObjectKindColumn) {
		return []*objectRefContext{newObjectRefContext(withColumn())}
	}
	ctx := newObjectRefContext(withColumn())
	q := r.completionIntent.Qualifier
	if q.Database == "" && q.Schema != "" && q.Object == "" &&
		(strings.EqualFold(unquoteIdent(q.Schema), "INSERTED") || strings.EqualFold(unquoteIdent(q.Schema), "DELETED")) {
		return []*objectRefContext{ctx}
	}
	if q.Object != "" {
		if q.Server != "" {
			ctx.setLinkedServer(unquoteIdent(q.Server))
		}
		if q.Database != "" {
			ctx.setDatabase(unquoteIdent(q.Database))
		}
		if q.Schema != "" {
			ctx.setSchema(unquoteIdent(q.Schema))
		}
		ctx.setObject(unquoteIdent(q.Object))
		ctx.flags &^= objectFlagShowDatabase | objectFlagShowSchema
		return []*objectRefContext{ctx}
	}
	if q.Server != "" && q.Database != "" && q.Schema != "" {
		ctx.setDatabase(unquoteIdent(q.Server))
		ctx.setSchema(unquoteIdent(q.Database))
		ctx.setObject(unquoteIdent(q.Schema))
		ctx.flags &^= objectFlagShowDatabase | objectFlagShowSchema
		return []*objectRefContext{ctx}
	}
	if q.Schema != "" {
		if q.Database != "" {
			ctx.setSchema(unquoteIdent(q.Database))
		}
		ctx.setObject(unquoteIdent(q.Schema))
		ctx.flags &^= objectFlagShowDatabase | objectFlagShowSchema
		return []*objectRefContext{ctx}
	}
	if q.Database != "" {
		ctx.setObject(unquoteIdent(q.Database))
		ctx.flags &^= objectFlagShowDatabase | objectFlagShowSchema
	}
	return []*objectRefContext{ctx}
}

func completionIntentHasObjectKind(intent *parser.CompletionIntent, kind parser.ObjectKind) bool {
	if intent == nil {
		return false
	}
	for _, k := range intent.ObjectKinds {
		if k == kind {
			return true
		}
	}
	return false
}

// --- SELECT item alias helpers ---

type selectAliasScope struct {
	clause      string
	selectStart int
	blocked     bool
}

func (r *catalogResolver) selectItemAliasStartOffset() (int, bool) {
	scopes := []selectAliasScope{{selectStart: -1}}
	previous := ""
	for idx, tok := range r.tokens {
		if tok.Loc >= r.cursorByteOffset {
			break
		}
		text := r.tokenText(idx)
		upper := strings.ToUpper(text)
		switch text {
		case "(":
			scopes = append(scopes, selectAliasScope{
				selectStart: -1,
				blocked:     previous == "OVER" || previous == "GROUP",
			})
		case ")":
			if len(scopes) > 1 {
				scopes = scopes[:len(scopes)-1]
			}
		default:
			scope := &scopes[len(scopes)-1]
			if !scope.blocked {
				switch upper {
				case "SELECT":
					scope.clause = "SELECT"
					scope.selectStart = tok.Loc
				case "FROM", "WHERE", "GROUP", "HAVING", "ORDER", "ON":
					scope.clause = upper
				case "UNION", "EXCEPT", "INTERSECT":
					scope.clause = ""
					scope.selectStart = -1
				}
			}
		}
		if upper != "" {
			previous = upper
		}
	}

	scope := scopes[len(scopes)-1]
	if scope.blocked || scope.selectStart < 0 {
		return 0, false
	}
	switch scope.clause {
	case "GROUP", "HAVING", "ORDER":
		return scope.selectStart, true
	default:
		return 0, false
	}
}

func (r *catalogResolver) fetchSelectItemAliases(aliasPositions []int, startOffset int) []string {
	aliasMap := make(map[string]bool)
	for _, pos := range aliasPositions {
		if pos < startOffset || pos >= r.cursorByteOffset {
			continue
		}
		if aliasText := r.extractAliasText(pos); aliasText != "" {
			aliasMap[aliasText] = true
		}
	}
	var result []string
	for alias := range aliasMap {
		result = append(result, alias)
	}
	slices.Sort(result)
	return result
}

func (r *catalogResolver) extractAliasText(pos int) string {
	if pos < 0 || pos >= len(r.sql) {
		return ""
	}
	for _, tok := range r.tokens {
		if tok.Loc == pos && parser.IsIdentTokenType(tok.Type) {
			return unquoteIdent(tok.Str)
		}
		if tok.Loc > pos {
			break
		}
	}
	return ""
}

func (r *catalogResolver) previousTokenIndex() int {
	idx := r.caretTokenIndex - 1
	for idx >= 0 && r.tokens[idx].End > r.cursorByteOffset {
		idx--
	}
	return idx
}

func (r *catalogResolver) tokenText(idx int) string {
	if idx < 0 || idx >= len(r.tokens) {
		return ""
	}
	tok := r.tokens[idx]
	if tok.Loc < 0 || tok.End > len(r.sql) || tok.Loc > tok.End {
		return ""
	}
	return r.sql[tok.Loc:tok.End]
}

// --- Identifier and derived-table parsing helpers ---

func unquoteIdent(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '[' && s[len(s)-1] == ']' {
		return strings.ReplaceAll(s[1:len(s)-1], "]]", "]")
	}
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return strings.ReplaceAll(s[1:len(s)-1], `""`, `"`)
	}
	return s
}

func parseMultipartIdentifier(text string) (string, string, string) {
	parts := splitMultipartIdentifier(text)
	for i, part := range parts {
		parts[i] = unquoteIdent(part)
	}
	switch len(parts) {
	case 0:
		return "", "", ""
	case 1:
		return "", "", parts[0]
	case 2:
		return "", parts[0], parts[1]
	default:
		return parts[len(parts)-3], parts[len(parts)-2], parts[len(parts)-1]
	}
}

func splitMultipartIdentifier(text string) []string {
	var parts []string
	var buf strings.Builder
	var bracket, quote bool
	for _, ch := range text {
		switch ch {
		case '[':
			bracket = true
		case ']':
			bracket = false
		case '"':
			quote = !quote
		case '.':
			if !bracket && !quote {
				parts = append(parts, buf.String())
				buf.Reset()
				continue
			}
		}
		buf.WriteRune(ch)
	}
	if buf.Len() > 0 {
		parts = append(parts, buf.String())
	}
	return parts
}

func splitTopLevelCSV(text string) []string {
	var parts []string
	var buf strings.Builder
	var bracket, quote bool
	level := 0
	for _, ch := range text {
		switch ch {
		case '[':
			bracket = true
		case ']':
			bracket = false
		case '"':
			quote = !quote
		case '(':
			if !bracket && !quote {
				level++
			}
		case ')':
			if !bracket && !quote && level > 0 {
				level--
			}
		case ',':
			if !bracket && !quote && level == 0 {
				parts = append(parts, buf.String())
				buf.Reset()
				continue
			}
		}
		buf.WriteRune(ch)
	}
	parts = append(parts, buf.String())
	return parts
}

type derivedTable struct {
	query string
	alias string
}

func extractDerivedTables(sql string) []derivedTable {
	var result []derivedTable
	upperSQL := strings.ToUpper(sql)
	for idx := 0; idx < len(sql); idx++ {
		if sql[idx] != '(' {
			continue
		}
		if !isDerivedTableOpen(upperSQL, idx) {
			continue
		}
		end := matchingParenInSQL(sql, idx)
		if end < 0 {
			continue
		}
		if alias, ok := readAliasAfter(sql, end+1); ok {
			result = append(result, derivedTable{
				query: sql[idx+1 : end],
				alias: alias,
			})
		}
		idx = end
	}
	return result
}

func isDerivedTableOpen(upperSQL string, idx int) bool {
	after := strings.TrimLeft(upperSQL[idx+1:], asciiWhitespace)
	if !strings.HasPrefix(after, "SELECT") && !strings.HasPrefix(after, "WITH") {
		return false
	}
	before := strings.TrimRight(upperSQL[:idx], asciiWhitespace)
	for _, kw := range []string{"FROM", "JOIN", "APPLY", "USING"} {
		if strings.HasSuffix(before, kw) {
			return true
		}
	}
	return false
}

func matchingParenInSQL(sql string, open int) int {
	level := 0
	var bracket, doubleQuote, singleQuote bool
	for i := open; i < len(sql); i++ {
		switch sql[i] {
		case '\'':
			if !bracket && !doubleQuote {
				singleQuote = !singleQuote
			}
		case '"':
			if !bracket && !singleQuote {
				doubleQuote = !doubleQuote
			}
		case '[':
			if !singleQuote && !doubleQuote {
				bracket = true
			}
		case ']':
			if bracket {
				bracket = false
			}
		case '(':
			if !bracket && !singleQuote && !doubleQuote {
				level++
			}
		case ')':
			if !bracket && !singleQuote && !doubleQuote {
				level--
				if level == 0 {
					return i
				}
			}
		}
	}
	return -1
}

func readAliasAfter(sql string, start int) (string, bool) {
	rest := strings.TrimLeft(sql[start:], asciiWhitespace)
	upperRest := strings.ToUpper(rest)
	if strings.HasPrefix(upperRest, "AS ") {
		rest = strings.TrimLeft(rest[2:], asciiWhitespace)
	}
	if rest == "" {
		return "", false
	}
	var alias string
	switch rest[0] {
	case '[':
		end := strings.IndexByte(rest, ']')
		if end <= 0 {
			return "", false
		}
		alias = rest[1:end]
	case '"':
		end := strings.IndexByte(rest[1:], '"')
		if end < 0 {
			return "", false
		}
		alias = rest[1 : end+1]
	default:
		end := 0
		for end < len(rest) {
			r, size := utf8.DecodeRuneInString(rest[end:])
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '@' && r != '#' && r != '$' {
				break
			}
			end += size
		}
		if end == 0 {
			return "", false
		}
		alias = rest[:end]
	}
	return alias, true
}

// dedup removes duplicate candidates (same text+type, case-insensitive).
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
