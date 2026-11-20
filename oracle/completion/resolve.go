package completion

import (
	"fmt"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oracle"
	"github.com/bytebase/omni/oracle/analysis"
	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/oracle/parser"
)

var oracleDataTypes = []string{
	"NUMBER",
	"VARCHAR2",
	"NVARCHAR2",
	"CHAR",
	"NCHAR",
	"DATE",
	"TIMESTAMP",
	"CLOB",
	"BLOB",
	"RAW",
	"INTEGER",
	"FLOAT",
	"DECIMAL",
	"SMALLINT",
	"LONG",
	"NCLOB",
}

var oracleBuiltinFunctionKeywords = map[string]bool{
	"CAST":                true,
	"COLLECT":             true,
	"DECODE":              true,
	"DENSE_RANK":          true,
	"JSON_ARRAY":          true,
	"JSON_EXISTS":         true,
	"JSON_MERGEPATCH":     true,
	"JSON_OBJECT":         true,
	"JSON_QUERY":          true,
	"JSON_VALUE":          true,
	"SYS_CONNECT_BY_PATH": true,
	"TREAT":               true,
	"XMLAGG":              true,
	"XMLELEMENT":          true,
	"XMLFOREST":           true,
	"XMLPARSE":            true,
	"XMLROOT":             true,
	"XMLSERIALIZE":        true,
}

type schemaStore struct {
	defaultSchema   string
	schemas         map[string]*metadatapb.SchemaMetadata
	schemaOrder     []string
	publicSynonyms  []*metadatapb.SynonymMetadata
	linkedDatabases []*metadatapb.LinkedDatabaseMetadata
	loadedMisses    map[string]bool
	schemaLoader    func(schema string) *metadatapb.DatabaseSchemaMetadata
}

func newSchemaStore(meta *metadatapb.DatabaseSchemaMetadata, cfg options) *schemaStore {
	s := &schemaStore{
		defaultSchema: cfg.defaultSchema,
		schemas:       make(map[string]*metadatapb.SchemaMetadata),
		loadedMisses:  make(map[string]bool),
		schemaLoader:  cfg.schemaLoader,
	}
	if s.defaultSchema == "" && meta != nil && meta.GetName() != "" {
		s.defaultSchema = meta.GetName()
	}
	if cfg.catalog != nil {
		s.ingestDatabaseMeta("", cfg.catalog.Metadata())
	}
	if meta != nil {
		s.ingestDatabaseMeta("", meta)
	}
	for name, dbMeta := range cfg.schemas {
		s.ingestDatabaseMeta(name, dbMeta)
	}
	for _, name := range cfg.schemaNames {
		s.registerSchemaName(name)
	}
	return s
}

func (s *schemaStore) registerSchemaName(name string) {
	if name == "" {
		return
	}
	for _, existing := range s.schemaOrder {
		if existing == name {
			return
		}
	}
	s.schemaOrder = append(s.schemaOrder, name)
}

func (s *schemaStore) ingestDatabaseMeta(explicitName string, dbMeta *metadatapb.DatabaseSchemaMetadata) {
	if dbMeta == nil {
		if explicitName != "" {
			s.registerSchemaName(explicitName)
		}
		return
	}
	for _, syn := range dbMeta.GetPublicSynonyms() {
		if syn != nil && syn.GetName() != "" {
			s.publicSynonyms = append(s.publicSynonyms, syn)
		}
	}
	for _, link := range dbMeta.GetLinkedDatabases() {
		if link != nil && link.GetName() != "" {
			s.linkedDatabases = append(s.linkedDatabases, link)
		}
	}
	schemas := dbMeta.GetSchemas()
	if len(schemas) == 0 {
		name := explicitName
		if name == "" {
			name = dbMeta.GetName()
		}
		if name != "" {
			s.registerSchemaName(name)
		}
		return
	}
	for _, sm := range schemas {
		if sm == nil {
			continue
		}
		name := sm.GetName()
		if name == "" {
			name = explicitName
		}
		if name == "" {
			name = dbMeta.GetName()
		}
		if name == "" {
			name = s.defaultSchema
		}
		s.schemas[name] = sm
		if name != "" {
			s.registerSchemaName(name)
		}
	}
}

func (s *schemaStore) listSchemas() []string {
	var result []string
	seen := make(map[string]bool)
	if s.defaultSchema != "" {
		result = append(result, s.defaultSchema)
		seen[s.defaultSchema] = true
	}
	for _, name := range s.schemaOrder {
		if name != "" && !seen[name] {
			seen[name] = true
			result = append(result, name)
		}
	}
	return result
}

func (s *schemaStore) getSchema(name string) (string, *metadatapb.SchemaMetadata) {
	return s.getSchemaWithCase(name, false)
}

func (s *schemaStore) getSchemaWithCase(name string, quoted bool) (string, *metadatapb.SchemaMetadata) {
	if sm, ok := s.schemas[name]; ok {
		return name, sm
	}
	if !quoted {
		for k, sm := range s.schemas {
			if strings.EqualFold(k, name) {
				return k, sm
			}
		}
	}
	if name == "" {
		if s.defaultSchema != "" {
			return s.getSchema(s.defaultSchema)
		}
		return "", nil
	}

	targetName := name
	if !quoted {
		for _, known := range s.schemaOrder {
			if strings.EqualFold(known, name) {
				targetName = known
				break
			}
		}
	}
	if s.loadedMisses[targetName] || s.schemaLoader == nil {
		return "", nil
	}
	loaded := s.schemaLoader(targetName)
	if loaded == nil {
		s.loadedMisses[targetName] = true
		return "", nil
	}
	s.ingestDatabaseMeta(targetName, loaded)
	if sm, ok := s.schemas[targetName]; ok {
		return targetName, sm
	}
	if !quoted {
		for k, sm := range s.schemas {
			if strings.EqualFold(k, targetName) {
				return k, sm
			}
		}
	}
	s.loadedMisses[targetName] = true
	return "", nil
}

type tableRef interface {
	isTableRef()
}

type physicalTableRef struct {
	Schema       string
	SchemaQuoted bool
	Table        string
	TableQuoted  bool
	Alias        string
	AliasQuoted  bool
}

func (*physicalTableRef) isTableRef() {}

type virtualTableRef struct {
	Table   string
	Quoted  bool
	Columns []string
}

func (*virtualTableRef) isTableRef() {}

type tableTarget struct {
	schema         string
	schemaQuoted   bool
	schemaExplicit bool
	table          string
	tableQuoted    bool
}

type schemaQuery struct {
	name   string
	quoted bool
}

type qualifierPart struct {
	name   string
	quoted bool
}

type scopeRefLevel struct {
	refs      []tableRef
	scopeRefs []parser.RangeReference
}

type resolver struct {
	sql              string
	collectOffset    int
	tokens           []parser.Token
	caretTokenIndex  int
	completionCtx    *parser.CompletionContext
	completionIntent *parser.CompletionIntent
	store            *schemaStore
	references       []tableRef
	scopeLevels      []scopeRefLevel
	cteTables        []*virtualTableRef
	scopeReferences  []parser.RangeReference
}

func newResolver(sql string, collectOffset int, ctx *parser.CompletionContext, store *schemaStore) *resolver {
	tokens := parser.Tokenize(sql)
	r := &resolver{
		sql:             sql,
		collectOffset:   collectOffset,
		tokens:          tokens,
		caretTokenIndex: findCaretTokenIndex(tokens, collectOffset),
		completionCtx:   ctx,
		store:           store,
	}
	if ctx != nil {
		r.completionIntent = ctx.Intent
		r.collectCompletionCTEs(ctx)
		r.collectCompletionScopeReferences(ctx)
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

func (r *resolver) resolve() []Candidate {
	if r.completionCtx == nil {
		return nil
	}
	cs := r.completionCtx.Candidates

	var (
		keywords    []Candidate
		functions   []Candidate
		schemas     []Candidate
		tables      []Candidate
		columns     []Candidate
		views       []Candidate
		mviews      []Candidate
		sequences   []Candidate
		synonyms    []Candidate
		dbLinks     []Candidate
		procs       []Candidate
		packages    []Candidate
		indexes     []Candidate
		triggers    []Candidate
		constraints []Candidate
		variables   []Candidate
		types       []Candidate
	)

	addSchemaCandidates := func() {
		for _, name := range r.store.listSchemas() {
			schemas = append(schemas, Candidate{Type: CandidateSchema, Text: name})
		}
	}

	addTypeCandidates := func(objectOnly bool) {
		schema, quoted := r.completionSchemaQualifier()
		if objectOnly && schema == "" {
			addSchemaCandidates()
		}
		r.insertSchemaCompositeTypes(&types, schema, quoted)
		if !objectOnly {
			for _, typ := range oracleDataTypes {
				types = append(types, Candidate{Type: CandidateType_, Text: typ})
			}
		}
	}

	addRelationCandidates := func(schemaSet map[schemaQuery]bool, includeTables, includeViews, includeMViews, includeSeqs, includeSyns bool) {
		includePublicSyns := false
		for sq := range schemaSet {
			if sq.name == "" && len(r.cteTables) > 0 && includeTables {
				for _, cte := range r.cteTables {
					tables = append(tables, Candidate{Type: CandidateTable, Text: cte.Table})
				}
			}
			if sq.name == "" {
				includePublicSyns = true
			}
			if sq.name == "" && r.store.defaultSchema != "" {
				continue
			}
			_, sm := r.store.getSchemaWithCase(sq.name, sq.quoted)
			if sm == nil {
				continue
			}
			if includeTables {
				for _, t := range sm.GetTables() {
					tables = append(tables, Candidate{Type: CandidateTable, Text: t.GetName()})
				}
				for _, ext := range sm.GetExternalTables() {
					tables = append(tables, Candidate{Type: CandidateTable, Text: ext.GetName()})
				}
			}
			if includeViews {
				for _, v := range sm.GetViews() {
					views = append(views, Candidate{Type: CandidateView, Text: v.GetName()})
				}
			}
			if includeMViews {
				for _, mv := range sm.GetMaterializedViews() {
					mviews = append(mviews, Candidate{Type: CandidateMaterializedView, Text: mv.GetName()})
				}
			}
			if includeSeqs {
				for _, seq := range sm.GetSequences() {
					sequences = append(sequences, Candidate{Type: CandidateSequence, Text: seq.GetName()})
				}
			}
			if includeSyns {
				for _, syn := range sm.GetSynonyms() {
					synonyms = append(synonyms, Candidate{Type: CandidateSynonym, Text: syn.GetName()})
				}
			}
		}
		if includeSyns && includePublicSyns {
			for _, syn := range r.store.publicSynonyms {
				synonyms = append(synonyms, Candidate{Type: CandidateSynonym, Text: syn.GetName()})
			}
		}
	}

	if cs != nil {
		for _, tok := range cs.Tokens {
			if tok > 0 && tok < 256 {
				continue
			}
			name := parser.TokenName(tok)
			if name == "" {
				continue
			}
			keywords = append(keywords, Candidate{Type: CandidateKeyword, Text: name})
		}

		for _, rule := range cs.Rules {
			switch rule.Rule {
			case "func_name":
				for _, tok := range cs.Tokens {
					name := parser.TokenName(tok)
					if oracleBuiltinFunctionKeywords[name] {
						functions = append(functions, Candidate{Type: CandidateFunction, Text: name})
					}
				}
				schema, quoted := r.completionSchemaQualifier()
				r.insertSchemaFunctions(&functions, schema, quoted)
			case "schema_ref":
				addSchemaCandidates()
			case "table_ref":
				schema, quoted := r.completionSchemaQualifier()
				switch {
				case r.completionIntentOnly(parser.ObjectKindTable):
					if schema == "" {
						addSchemaCandidates()
					}
					addRelationCandidates(r.completionSchemas(schema, quoted), true, false, false, false, false)
				case r.completionIntentOnly(parser.ObjectKindView):
					if schema == "" {
						addSchemaCandidates()
					}
					if r.isMaterializedViewContext() {
						addRelationCandidates(r.completionSchemas(schema, quoted), false, false, true, false, false)
					} else {
						addRelationCandidates(r.completionSchemas(schema, quoted), false, true, false, false, false)
					}
				default:
					if schema == "" {
						addSchemaCandidates()
					}
					addRelationCandidates(r.completionSchemas(schema, quoted), true, true, true, true, true)
				}
			case "sequence_ref":
				addRelationCandidates(r.sequenceCompletionSchemas(), false, false, false, true, false)
			case "sequence_member_ref":
				r.insertSequenceMemberCandidates(&keywords)
			case "synonym_ref":
				schema, quoted := r.completionSchemaQualifier()
				addRelationCandidates(r.completionSchemas(schema, quoted), false, false, false, false, true)
			case "database_link_ref":
				r.insertDatabaseLinkCandidates(&dbLinks)
			case "proc_ref":
				schema, quoted := r.completionSchemaQualifier()
				r.insertSchemaProcedures(&procs, schema, quoted)
			case "package_ref":
				schema, quoted := r.completionSchemaQualifier()
				r.insertSchemaPackages(&packages, schema, quoted)
			case "package_member_ref":
				r.insertPackageMemberCandidates(&procs, &functions)
			case "index_ref":
				schema, quoted := r.completionSchemaQualifier()
				r.insertSchemaIndexes(&indexes, schema, quoted)
			case "trigger_ref":
				schema, quoted := r.completionSchemaQualifier()
				r.insertSchemaTriggers(&triggers, schema, quoted)
			case "constraint_ref":
				r.insertConstraintCandidates(&constraints)
			case "variable_ref":
				r.insertLocalVariableCandidates(&variables)
			case "type_name":
				addTypeCandidates(r.isObjectTypeDDLContext())
			case "columnref":
				if !r.isMaterializedViewContext() && !r.isObjectTypeDDLContext() {
					r.insertColumnContextCandidates(&tables, &columns, addSchemaCandidates, addRelationCandidates)
				}
			}
		}
	}

	if r.completionIntent != nil {
		for _, kind := range r.completionIntent.ObjectKinds {
			switch kind {
			case parser.ObjectKindSchema:
				addSchemaCandidates()
			case parser.ObjectKindTable:
				schema, quoted := r.completionSchemaQualifier()
				if schema == "" {
					addSchemaCandidates()
				}
				addRelationCandidates(r.completionSchemas(schema, quoted), true, false, false, false, false)
			case parser.ObjectKindView:
				schema, quoted := r.completionSchemaQualifier()
				if schema == "" {
					addSchemaCandidates()
				}
				if r.isMaterializedViewContext() {
					addRelationCandidates(r.completionSchemas(schema, quoted), false, false, true, false, false)
				} else {
					addRelationCandidates(r.completionSchemas(schema, quoted), false, true, false, false, false)
				}
			case parser.ObjectKindSequence:
				addRelationCandidates(r.sequenceCompletionSchemas(), false, false, false, true, false)
			case parser.ObjectKindSequenceMember:
				if !r.isMaterializedViewContext() && !r.isObjectTypeDDLContext() {
					r.insertSequenceMemberCandidates(&keywords)
				}
			case parser.ObjectKindSynonym:
				schema, quoted := r.completionSchemaQualifier()
				addRelationCandidates(r.completionSchemas(schema, quoted), false, false, false, false, true)
			case parser.ObjectKindDatabaseLink:
				r.insertDatabaseLinkCandidates(&dbLinks)
			case parser.ObjectKindColumn:
				if !r.isMaterializedViewContext() && !r.isObjectTypeDDLContext() {
					r.insertColumnContextCandidates(&tables, &columns, addSchemaCandidates, addRelationCandidates)
				}
			case parser.ObjectKindFunction:
				schema, quoted := r.completionSchemaQualifier()
				r.insertSchemaFunctions(&functions, schema, quoted)
			case parser.ObjectKindProcedure:
				schema, quoted := r.completionSchemaQualifier()
				r.insertSchemaProcedures(&procs, schema, quoted)
			case parser.ObjectKindPackage:
				schema, quoted := r.completionSchemaQualifier()
				r.insertSchemaPackages(&packages, schema, quoted)
			case parser.ObjectKindPackageMember:
				if !r.isMaterializedViewContext() && !r.isObjectTypeDDLContext() {
					r.insertPackageMemberCandidates(&procs, &functions)
				}
			case parser.ObjectKindIndex:
				schema, quoted := r.completionSchemaQualifier()
				r.insertSchemaIndexes(&indexes, schema, quoted)
			case parser.ObjectKindTrigger:
				schema, quoted := r.completionSchemaQualifier()
				r.insertSchemaTriggers(&triggers, schema, quoted)
			case parser.ObjectKindConstraint:
				r.insertConstraintCandidates(&constraints)
			case parser.ObjectKindVariable:
				r.insertLocalVariableCandidates(&variables)
			case parser.ObjectKindType:
				addTypeCandidates(r.isObjectTypeDDLContext())
			}
		}
	}

	switch {
	case r.isMaterializedViewContext():
		if parts := r.qualifierParts(); len(parts) == 1 {
			addRelationCandidates(r.completionSchemas(parts[0].name, parts[0].quoted), false, false, true, false, false)
		}
	case r.isDropOrAlterTypeContext():
		addTypeCandidates(true)
	case r.isGrantExecuteContext():
		if parts := r.qualifierParts(); len(parts) == 1 {
			r.insertSchemaProcedures(&procs, parts[0].name, parts[0].quoted)
			r.insertSchemaFunctions(&functions, parts[0].name, parts[0].quoted)
			r.insertSchemaPackages(&packages, parts[0].name, parts[0].quoted)
			r.insertSchemaCompositeTypes(&types, parts[0].name, parts[0].quoted)
		}
	case len(r.qualifierParts()) > 0 && len(columns) == 0:
		r.insertColumnContextCandidates(&tables, &columns, addSchemaCandidates, addRelationCandidates)
	}

	var result []Candidate
	result = append(result, keywords...)
	result = append(result, types...)
	result = append(result, functions...)
	result = append(result, procs...)
	result = append(result, packages...)
	result = append(result, schemas...)
	result = append(result, tables...)
	result = append(result, columns...)
	result = append(result, views...)
	result = append(result, mviews...)
	result = append(result, sequences...)
	result = append(result, synonyms...)
	result = append(result, dbLinks...)
	result = append(result, indexes...)
	result = append(result, triggers...)
	result = append(result, constraints...)
	result = append(result, variables...)
	return result
}

func (r *resolver) completionIntentOnly(kind parser.ObjectKind) bool {
	return r.completionIntent != nil && len(r.completionIntent.ObjectKinds) == 1 && r.completionIntent.ObjectKinds[0] == kind
}

func (r *resolver) completionSchemaQualifier() (string, bool) {
	parts := r.qualifierParts()
	if r.completionIntent != nil && r.completionIntent.Qualifier.Schema != "" {
		if len(parts) == 1 && strings.EqualFold(parts[0].name, r.completionIntent.Qualifier.Schema) {
			return parts[0].name, parts[0].quoted
		}
		return r.completionIntent.Qualifier.Schema, false
	}
	if len(parts) == 1 {
		return parts[0].name, parts[0].quoted
	}
	return "", false
}

func (r *resolver) sequenceCompletionSchemas() map[schemaQuery]bool {
	schema, quoted := r.completionSchemaQualifier()
	return r.completionSchemas(schema, quoted)
}

func (r *resolver) completionSchemas(schema string, quoted bool) map[schemaQuery]bool {
	schemas := make(map[schemaQuery]bool)
	if schema != "" {
		schemas[schemaQuery{name: schema, quoted: quoted}] = true
		return schemas
	}
	if r.store.defaultSchema != "" {
		schemas[schemaQuery{name: r.store.defaultSchema}] = true
	}
	schemas[schemaQuery{}] = true
	return schemas
}

func (r *resolver) insertDatabaseLinkCandidates(dst *[]Candidate) {
	for _, link := range r.store.linkedDatabases {
		*dst = append(*dst, Candidate{Type: CandidateDatabaseLink, Text: link.GetName()})
	}
}

func (r *resolver) insertSchemaFunctions(dst *[]Candidate, schema string, quoted bool) {
	for sq := range r.completionSchemas(schema, quoted) {
		if sq.name == "" && r.store.defaultSchema != "" {
			continue
		}
		_, sm := r.store.getSchemaWithCase(sq.name, sq.quoted)
		if sm == nil {
			continue
		}
		for _, fn := range sm.GetFunctions() {
			*dst = append(*dst, Candidate{Type: CandidateFunction, Text: fn.GetName(), Comment: fn.GetComment()})
		}
	}
}

func (r *resolver) insertSchemaProcedures(dst *[]Candidate, schema string, quoted bool) {
	for sq := range r.completionSchemas(schema, quoted) {
		if sq.name == "" && r.store.defaultSchema != "" {
			continue
		}
		_, sm := r.store.getSchemaWithCase(sq.name, sq.quoted)
		if sm == nil {
			continue
		}
		for _, p := range sm.GetProcedures() {
			*dst = append(*dst, Candidate{Type: CandidateProcedure, Text: p.GetName(), Comment: p.GetComment()})
		}
	}
}

func (r *resolver) insertSchemaCompositeTypes(dst *[]Candidate, schema string, quoted bool) {
	for sq := range r.completionSchemas(schema, quoted) {
		if sq.name == "" && r.store.defaultSchema != "" {
			continue
		}
		_, sm := r.store.getSchemaWithCase(sq.name, sq.quoted)
		if sm == nil {
			continue
		}
		for _, ct := range sm.GetCompositeTypes() {
			if ct.GetName() != "" {
				*dst = append(*dst, Candidate{Type: CandidateType_, Text: ct.GetName()})
			}
		}
	}
}

func (r *resolver) isObjectTypeDDLContext() bool {
	return r.isDropOrAlterTypeContext() || r.isGrantExecuteContext()
}

func (r *resolver) isDropOrAlterTypeContext() bool {
	stmtStart, _ := r.statementTokenBounds()
	idx := r.caretTokenIndex - 1
	if idx >= stmtStart && r.tokens[idx].Type == '.' {
		idx--
		if idx >= stmtStart && r.isIdentifierToken(r.tokens[idx]) {
			idx--
		} else {
			return false
		}
	}
	if idx-1 < stmtStart {
		return false
	}
	kw := strings.ToUpper(parser.TokenName(r.tokens[idx].Type))
	prevKw := strings.ToUpper(parser.TokenName(r.tokens[idx-1].Type))
	return kw == "TYPE" && (prevKw == "DROP" || prevKw == "ALTER")
}

func (r *resolver) isGrantExecuteContext() bool {
	stmtStart, _ := r.statementTokenBounds()
	if stmtStart >= r.caretTokenIndex {
		return false
	}
	stmtKw := strings.ToUpper(parser.TokenName(r.tokens[stmtStart].Type))
	if stmtKw != "GRANT" && stmtKw != "REVOKE" {
		return false
	}
	seenExecute := false
	onIdx := -1
	for i := stmtStart + 1; i < r.caretTokenIndex; i++ {
		switch strings.ToUpper(parser.TokenName(r.tokens[i].Type)) {
		case "EXECUTE":
			if onIdx < 0 {
				seenExecute = true
			}
		case "ON":
			onIdx = i
		}
	}
	if !seenExecute || onIdx < 0 {
		return false
	}
	rest := r.caretTokenIndex - 1 - onIdx
	if rest == 0 {
		return true
	}
	return rest == 2 && r.isIdentifierToken(r.tokens[onIdx+1]) && r.tokens[onIdx+2].Type == '.'
}

func (r *resolver) insertSchemaPackages(dst *[]Candidate, schema string, quoted bool) {
	for sq := range r.completionSchemas(schema, quoted) {
		if sq.name == "" && r.store.defaultSchema != "" {
			continue
		}
		_, sm := r.store.getSchemaWithCase(sq.name, sq.quoted)
		if sm == nil {
			continue
		}
		for _, pkg := range sm.GetPackages() {
			*dst = append(*dst, Candidate{Type: CandidatePackage, Text: pkg.GetName()})
		}
	}
}

func (r *resolver) insertSchemaIndexes(dst *[]Candidate, schema string, quoted bool) {
	for sq := range r.completionSchemas(schema, quoted) {
		if sq.name == "" && r.store.defaultSchema != "" {
			continue
		}
		_, sm := r.store.getSchemaWithCase(sq.name, sq.quoted)
		if sm == nil {
			continue
		}
		for _, t := range sm.GetTables() {
			for _, idx := range t.GetIndexes() {
				*dst = append(*dst, Candidate{Type: CandidateIndex, Text: idx.GetName()})
			}
		}
		for _, mv := range sm.GetMaterializedViews() {
			for _, idx := range mv.GetIndexes() {
				*dst = append(*dst, Candidate{Type: CandidateIndex, Text: idx.GetName()})
			}
		}
	}
}

func (r *resolver) insertPackageMemberCandidates(procsOut, funcsOut *[]Candidate) {
	schema, pkgName, schemaQuoted, pkgQuoted := r.columnQualifier()
	if pkgName == "" && r.completionIntent != nil {
		pkgName = r.completionIntent.Qualifier.Object
	}
	if pkgName == "" {
		return
	}
	for sq := range r.completionSchemas(schema, schemaQuoted) {
		if sq.name == "" && r.store.defaultSchema != "" {
			continue
		}
		_, sm := r.store.getSchemaWithCase(sq.name, sq.quoted)
		if sm == nil {
			continue
		}
		for _, pkg := range sm.GetPackages() {
			if nameMatches(pkg.GetName(), pkgName, pkgQuoted) {
				procs, funcs := parsePackageMembers(pkg.GetDefinition())
				for _, p := range procs {
					*procsOut = append(*procsOut, Candidate{Type: CandidateProcedure, Text: p})
				}
				for _, f := range funcs {
					*funcsOut = append(*funcsOut, Candidate{Type: CandidateFunction, Text: f})
				}
				return
			}
		}
	}
}

func (r *resolver) insertConstraintCandidates(dst *[]Candidate) {
	for _, ref := range r.references {
		phys, ok := ref.(*physicalTableRef)
		if !ok {
			continue
		}
		schema := phys.Schema
		schemaQuoted := phys.SchemaQuoted
		if schema == "" {
			schema = r.store.defaultSchema
			schemaQuoted = false
		}
		_, sm := r.store.getSchemaWithCase(schema, schemaQuoted)
		if sm == nil {
			continue
		}
		t := findTableMeta(sm, phys.Table, phys.TableQuoted)
		if t == nil {
			continue
		}
		for _, cc := range t.GetCheckConstraints() {
			if cc.GetName() != "" {
				*dst = append(*dst, Candidate{Type: CandidateConstraint, Text: cc.GetName()})
			}
		}
		for _, fk := range t.GetForeignKeys() {
			if fk.GetName() != "" {
				*dst = append(*dst, Candidate{Type: CandidateConstraint, Text: fk.GetName()})
			}
		}
		for _, idx := range t.GetIndexes() {
			if len(idx.GetConstraints()) > 0 {
				for _, c := range idx.GetConstraints() {
					if c.GetName() != "" {
						*dst = append(*dst, Candidate{Type: CandidateConstraint, Text: c.GetName()})
					}
				}
			} else if idx.GetIsConstraint() && idx.GetName() != "" {
				*dst = append(*dst, Candidate{Type: CandidateConstraint, Text: idx.GetName()})
			}
		}
	}
}

func (r *resolver) insertSchemaTriggers(dst *[]Candidate, schema string, quoted bool) {
	for sq := range r.completionSchemas(schema, quoted) {
		if sq.name == "" && r.store.defaultSchema != "" {
			continue
		}
		_, sm := r.store.getSchemaWithCase(sq.name, sq.quoted)
		if sm == nil {
			continue
		}
		for _, t := range sm.GetTables() {
			for _, trig := range t.GetTriggers() {
				*dst = append(*dst, Candidate{Type: CandidateTrigger, Text: trig.GetName()})
			}
		}
		for _, ext := range sm.GetExternalTables() {
			for _, trig := range ext.GetTriggers() {
				*dst = append(*dst, Candidate{Type: CandidateTrigger, Text: trig.GetName()})
			}
		}
		for _, v := range sm.GetViews() {
			for _, trig := range v.GetTriggers() {
				*dst = append(*dst, Candidate{Type: CandidateTrigger, Text: trig.GetName()})
			}
		}
		for _, mv := range sm.GetMaterializedViews() {
			for _, trig := range mv.GetTriggers() {
				*dst = append(*dst, Candidate{Type: CandidateTrigger, Text: trig.GetName()})
			}
		}
	}
}

func (r *resolver) insertSequenceMemberCandidates(keywords *[]Candidate) {
	schema, seqName, schemaQuoted, seqQuoted := r.columnQualifier()
	if seqName == "" && r.completionIntent != nil {
		seqName = r.completionIntent.Qualifier.Object
	}
	if seqName == "" {
		return
	}
	for sq := range r.completionSchemas(schema, schemaQuoted) {
		if sq.name == "" && r.store.defaultSchema != "" {
			continue
		}
		_, sm := r.store.getSchemaWithCase(sq.name, sq.quoted)
		if sm == nil {
			continue
		}
		for _, seq := range sm.GetSequences() {
			if nameMatches(seq.GetName(), seqName, seqQuoted) {
				*keywords = append(*keywords,
					Candidate{Type: CandidateKeyword, Text: "NEXTVAL"},
					Candidate{Type: CandidateKeyword, Text: "CURRVAL"},
				)
				return
			}
		}
	}
}

func (r *resolver) insertColumnContextCandidates(
	tablesOut, columnsOut *[]Candidate,
	addSchemaCandidates func(),
	addRelationCandidates func(schemaSet map[schemaQuery]bool, includeTables, includeViews, includeMViews, includeSeqs, includeSyns bool),
) {
	schema, table, schemaQuoted, tableQuoted := r.columnQualifier()
	if table == "" && r.completionIntent != nil {
		table = r.completionIntent.Qualifier.Object
	}

	schemas := make(map[schemaQuery]bool)
	if schema != "" {
		schemas[schemaQuery{name: schema, quoted: schemaQuoted}] = true
	} else if table == "" {
		for _, ref := range r.references {
			if physical, ok := ref.(*physicalTableRef); ok && physical.Schema != "" {
				schemas[schemaQuery{name: physical.Schema, quoted: physical.SchemaQuoted}] = true
			}
		}
	}

	var targets []tableTarget
	addTarget := func(s, t string, sQuoted, tQuoted bool) {
		if t == "" {
			return
		}
		explicit := s != ""
		if s == "" {
			s = r.store.defaultSchema
			sQuoted = false
		}
		for _, existing := range targets {
			if existing.schema == s && existing.table == t && existing.schemaQuoted == sQuoted && existing.tableQuoted == tQuoted && existing.schemaExplicit == explicit {
				return
			}
		}
		targets = append(targets, tableTarget{
			schema:         s,
			schemaQuoted:   sQuoted,
			schemaExplicit: explicit,
			table:          t,
			tableQuoted:    tQuoted,
		})
	}

	if table != "" {
		matchedRef := false
		for _, level := range r.scopeLevels {
			for _, ref := range level.refs {
				switch ref := ref.(type) {
				case *physicalTableRef:
					if schema != "" && ref.Schema != "" && !nameMatches(ref.Schema, schema, ref.SchemaQuoted || schemaQuoted) {
						continue
					}
					if matchPhysicalQualifier(ref, table, tableQuoted) {
						targetSchema := ref.Schema
						targetSchemaQuoted := ref.SchemaQuoted
						if targetSchema == "" {
							targetSchema = schema
							targetSchemaQuoted = schemaQuoted
						}
						addTarget(targetSchema, ref.Table, targetSchemaQuoted, ref.TableQuoted)
						matchedRef = true
					}
				case *virtualTableRef:
					if schema == "" && nameMatches(ref.Table, table, tableQuoted) {
						for _, col := range ref.Columns {
							*columnsOut = append(*columnsOut, Candidate{Type: CandidateColumn, Text: col})
						}
						matchedRef = true
					}
				}
			}
			if matchedRef {
				break
			}
		}
		if !matchedRef && (len(r.references) == 0 || r.isCommentOnColumnContext()) {
			// Fallback for DDL/standalone qualified names (e.g. COMMENT ON COLUMN schema1.t2.|).
			addTarget(schema, table, schemaQuoted, tableQuoted)
		}
	} else if len(r.references) > 0 {
		for _, ref := range r.references {
			switch ref := ref.(type) {
			case *physicalTableRef:
				addTarget(ref.Schema, ref.Table, ref.SchemaQuoted, ref.TableQuoted)
			case *virtualTableRef:
				for _, col := range ref.Columns {
					*columnsOut = append(*columnsOut, Candidate{Type: CandidateColumn, Text: col})
				}
			}
		}
	} else {
		r.insertAllColumns(columnsOut)
	}

	if len(schemas) == 0 {
		schemas[schemaQuery{name: r.store.defaultSchema}] = true
	}
	schemas[schemaQuery{}] = true

	if table == "" {
		addSchemaCandidates()
		addRelationCandidates(schemas, true, true, true, true, true)
		for _, alias := range r.selectItemAliases() {
			*columnsOut = append(*columnsOut, Candidate{Type: CandidateColumn, Text: alias})
		}
		r.insertReferenceAliasCandidates(tablesOut)
	}

	if len(targets) > 0 {
		r.insertColumnsForTargets(columnsOut, targets)
	}
	if table != "" && len(*columnsOut) == 0 {
		r.insertAliasColumns(columnsOut, table, tableQuoted)
	}
}

func matchPhysicalQualifier(ref *physicalTableRef, qualifier string, qualifierQuoted bool) bool {
	if ref.Alias != "" {
		return nameMatches(ref.Alias, qualifier, ref.AliasQuoted || qualifierQuoted)
	}
	return nameMatches(ref.Table, qualifier, ref.TableQuoted || qualifierQuoted)
}

func (r *resolver) isMaterializedViewContext() bool {
	stmtStart, _ := r.statementTokenBounds()
	idx := r.caretTokenIndex - 1
	if idx >= stmtStart && r.tokens[idx].Type == '.' {
		idx--
		if idx >= stmtStart && r.isIdentifierToken(r.tokens[idx]) {
			idx--
		} else {
			return false
		}
	}
	return idx-1 >= stmtStart &&
		r.tokens[idx].Type == parser.VIEW &&
		strings.EqualFold(r.tokens[idx-1].Str, "MATERIALIZED")
}

func (r *resolver) isCommentOnColumnContext() bool {
	stmtStart, _ := r.statementTokenBounds()
	if stmtStart+2 >= r.caretTokenIndex {
		return false
	}
	return strings.EqualFold(parser.TokenName(r.tokens[stmtStart].Type), "COMMENT") &&
		strings.EqualFold(parser.TokenName(r.tokens[stmtStart+1].Type), "ON") &&
		strings.EqualFold(parser.TokenName(r.tokens[stmtStart+2].Type), "COLUMN")
}

func (r *resolver) insertAllColumns(columnsOut *[]Candidate) {
	for _, schema := range r.store.listSchemas() {
		resolvedSchema, sm := r.store.getSchema(schema)
		if sm == nil {
			continue
		}
		for _, t := range sm.GetTables() {
			for _, col := range t.GetColumns() {
				*columnsOut = append(*columnsOut, columnCandidate(resolvedSchema, t.GetName(), col))
			}
		}
	}
	if r.store.defaultSchema == "" {
		if _, sm := r.store.getSchema(""); sm != nil {
			for _, t := range sm.GetTables() {
				for _, col := range t.GetColumns() {
					*columnsOut = append(*columnsOut, columnCandidate("", t.GetName(), col))
				}
			}
		}
	}
}

func (r *resolver) insertColumnsForTargets(columnsOut *[]Candidate, targets []tableTarget) {
	for _, target := range targets {
		if !target.schemaExplicit {
			for _, cte := range r.cteTables {
				if nameMatches(cte.Table, target.table, cte.Quoted || target.tableQuoted) {
					for _, col := range cte.Columns {
						*columnsOut = append(*columnsOut, Candidate{Type: CandidateColumn, Text: col})
					}
				}
			}
		}
		r.insertColumnsForObject(columnsOut, target.schema, target.table, target.schemaQuoted, target.tableQuoted, !target.schemaExplicit, make(map[string]bool))
	}
}

func (r *resolver) insertColumnsForObject(columnsOut *[]Candidate, schema, name string, schemaQuoted, nameQuoted, allowPublicSyn bool, visited map[string]bool) {
	if name == "" {
		return
	}
	resolvedSchema, sm := r.store.getSchemaWithCase(schema, schemaQuoted)
	visitKey := resolvedSchema + "\x00" + name
	if visited[visitKey] {
		return
	}
	visited[visitKey] = true
	if sm != nil {
		if t := findTableMeta(sm, name, nameQuoted); t != nil {
			for _, col := range t.GetColumns() {
				*columnsOut = append(*columnsOut, columnCandidate(resolvedSchema, t.GetName(), col))
			}
			return
		}
		if ext := findExternalTableMeta(sm, name, nameQuoted); ext != nil {
			for _, col := range ext.GetColumns() {
				*columnsOut = append(*columnsOut, columnCandidate(resolvedSchema, ext.GetName(), col))
			}
			return
		}
		if v := findViewMeta(sm, name, nameQuoted); v != nil {
			if len(v.GetColumns()) > 0 {
				for _, col := range v.GetColumns() {
					*columnsOut = append(*columnsOut, columnCandidate(resolvedSchema, v.GetName(), col))
				}
				return
			}
			if v.GetDefinition() != "" {
				r.insertResolvedRelationColumns(columnsOut, resolvedSchema, v.GetName(), v.GetDefinition())
				return
			}
		}
		if mv := findMaterializedViewMeta(sm, name, nameQuoted); mv != nil {
			if mv.GetDefinition() != "" {
				r.insertResolvedRelationColumns(columnsOut, resolvedSchema, mv.GetName(), mv.GetDefinition())
			}
			return
		}
		if syn := findSynonymMeta(sm.GetSynonyms(), name, nameQuoted); syn != nil && syn.GetTargetName() != "" {
			targetSchema := syn.GetTargetSchema()
			targetExplicit := targetSchema != ""
			if targetSchema == "" {
				targetSchema = resolvedSchema
			}
			r.insertColumnsForObject(columnsOut, targetSchema, syn.GetTargetName(), false, false, !targetExplicit, visited)
			return
		}
	}
	if allowPublicSyn && (schema == "" || nameMatches(schema, r.store.defaultSchema, schemaQuoted)) {
		if syn := findSynonymMeta(r.store.publicSynonyms, name, nameQuoted); syn != nil && syn.GetTargetName() != "" {
			targetSchema := syn.GetTargetSchema()
			targetExplicit := targetSchema != ""
			if targetSchema == "" {
				targetSchema = r.store.defaultSchema
			}
			r.insertColumnsForObject(columnsOut, targetSchema, syn.GetTargetName(), false, false, !targetExplicit, visited)
		}
	}
}

func nameMatches(candidate, query string, quoted bool) bool {
	if quoted {
		return candidate == query
	}
	return candidate == query || strings.EqualFold(candidate, query)
}

func findExternalTableMeta(sm *metadatapb.SchemaMetadata, name string, quoted bool) *metadatapb.ExternalTableMetadata {
	for _, ext := range sm.GetExternalTables() {
		if ext.GetName() == name {
			return ext
		}
	}
	if !quoted {
		for _, ext := range sm.GetExternalTables() {
			if strings.EqualFold(ext.GetName(), name) {
				return ext
			}
		}
	}
	return nil
}

func findMaterializedViewMeta(sm *metadatapb.SchemaMetadata, name string, quoted bool) *metadatapb.MaterializedViewMetadata {
	for _, mv := range sm.GetMaterializedViews() {
		if mv.GetName() == name {
			return mv
		}
	}
	if !quoted {
		for _, mv := range sm.GetMaterializedViews() {
			if strings.EqualFold(mv.GetName(), name) {
				return mv
			}
		}
	}
	return nil
}

func (r *resolver) insertResolvedRelationColumns(columnsOut *[]Candidate, schema, objectName, definition string) {
	target := quoteIdentifierForSQL(objectName)
	defaultSchema := r.store.defaultSchema
	if schema != "" {
		target = quoteIdentifierForSQL(schema) + "." + target
		defaultSchema = schema
	}
	span, err := analysis.ResolveQuerySpan(fmt.Sprintf("SELECT * FROM %s", target), analysis.ResolveOptions{
		DefaultDatabase: defaultSchema,
		Catalog:         r.analysisCatalog(),
	})
	if (err != nil || span == nil || span.NotFoundError != nil || len(span.Results) == 0) && definition != "" {
		if queryText := extractSelectFromViewDefinition(definition); queryText != "" {
			span, err = analysis.ResolveQuerySpan(queryText, analysis.ResolveOptions{
				DefaultDatabase: defaultSchema,
				Catalog:         r.analysisCatalog(),
			})
		}
	}
	if err != nil || span == nil || span.NotFoundError != nil {
		return
	}
	for _, col := range span.Results {
		if col.Name != "" {
			*columnsOut = append(*columnsOut, Candidate{Type: CandidateColumn, Text: col.Name})
		}
	}
}

func extractSelectFromViewDefinition(definition string) string {
	// Definition comes from an engine-accepted synced view or materialized view; downgrade parse errors to no inferred columns.
	stmts, err := oracle.Parse(definition)
	if err != nil {
		return ""
	}
	for _, stmt := range stmts {
		cv, ok := stmt.AST.(*ast.CreateViewStmt)
		if !ok || cv == nil || cv.Query == nil {
			continue
		}
		if sel, ok := cv.Query.(*ast.SelectStmt); ok && sel.Loc.Start >= 0 && sel.Loc.End > sel.Loc.Start && sel.Loc.End <= len(definition) {
			return definition[sel.Loc.Start:sel.Loc.End]
		}
	}
	return ""
}

func findSynonymMeta(syns []*metadatapb.SynonymMetadata, name string, quoted bool) *metadatapb.SynonymMetadata {
	for _, syn := range syns {
		if syn.GetName() == name {
			return syn
		}
	}
	if !quoted {
		for _, syn := range syns {
			if strings.EqualFold(syn.GetName(), name) {
				return syn
			}
		}
	}
	return nil
}

func findTableMeta(sm *metadatapb.SchemaMetadata, name string, quoted bool) *metadatapb.TableMetadata {
	for _, t := range sm.GetTables() {
		if t.GetName() == name {
			return t
		}
	}
	if !quoted {
		for _, t := range sm.GetTables() {
			if strings.EqualFold(t.GetName(), name) {
				return t
			}
		}
	}
	return nil
}

func findViewMeta(sm *metadatapb.SchemaMetadata, name string, quoted bool) *metadatapb.ViewMetadata {
	for _, v := range sm.GetViews() {
		if v.GetName() == name {
			return v
		}
	}
	if !quoted {
		for _, v := range sm.GetViews() {
			if strings.EqualFold(v.GetName(), name) {
				return v
			}
		}
	}
	return nil
}

func columnCandidate(schema, table string, col *metadatapb.ColumnMetadata) Candidate {
	target := table
	if schema != "" {
		target = schema + "." + table
	}
	definition := fmt.Sprintf("%s | %s", target, col.GetType())
	if !col.GetNullable() {
		definition += ", NOT NULL"
	}
	return Candidate{
		Type:       CandidateColumn,
		Text:       col.GetName(),
		Definition: definition,
		Comment:    col.GetComment(),
	}
}

func (r *resolver) collectCompletionCTEs(ctx *parser.CompletionContext) {
	for _, ref := range ctx.CTEs {
		// Hide later CTEs, and hide the current CTE inside its own anchor member
		// while exposing it in any recursive UNION member.
		inRecursiveMember := false
		unionLoc := -1
		if ref.Loc.End >= 0 && r.collectOffset < ref.Loc.End {
			if ref.Loc.Start >= 0 && r.collectOffset < ref.Loc.Start {
				continue
			}
			unionLoc = r.firstTopLevelUnionBefore(ref.BodyLoc.Start, r.collectOffset)
			if unionLoc < 0 {
				continue
			}
			inRecursiveMember = true
		}
		_, cteQuoted, _, _ := r.refQuotedFlags(ref)
		virtual := convertCompletionVirtualReference(ref, cteQuoted)
		if virtual == nil {
			continue
		}
		if len(virtual.Columns) == 0 {
			if inRecursiveMember {
				virtual.Columns = r.inferCTEAnchorColumns(ref, unionLoc)
			} else {
				virtual.Columns = r.inferCompletionVirtualColumns(ref)
			}
		}
		r.appendVirtualReferenceIfMissing(virtual)
	}
}

func (r *resolver) firstTopLevelUnionBefore(bodyStart, limit int) int {
	if bodyStart < 0 || bodyStart >= limit || limit > len(r.sql) {
		return -1
	}
	depth := 0
	for _, tok := range r.tokens {
		if tok.Loc < bodyStart {
			continue
		}
		if tok.Loc >= limit {
			break
		}
		switch tok.Type {
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		case parser.UNION:
			if depth == 0 {
				return tok.Loc
			}
		}
	}
	return -1
}

func (r *resolver) collectCompletionScopeReferences(ctx *parser.CompletionContext) {
	if ctx.Scope == nil {
		return
	}
	collectLevel := func(raw []parser.RangeReference) {
		var level scopeRefLevel
		for _, ref := range raw {
			if converted := r.convertCompletionScopeReference(ref); converted != nil {
				r.appendReferenceIfMissing(converted)
				r.appendScopeReferenceIfMissing(ref)
				level.refs = append(level.refs, converted)
				level.scopeRefs = append(level.scopeRefs, ref)
			}
		}
		if len(level.refs) > 0 {
			r.scopeLevels = append(r.scopeLevels, level)
		}
	}
	collectLevel(ctx.Scope.LocalReferences)
	for _, outer := range ctx.Scope.OuterReferences {
		collectLevel(outer)
	}
}

func (r *resolver) convertCompletionScopeReference(ref parser.RangeReference) tableRef {
	if ref.Unsupported || r.hasDatabaseLink(ref) {
		table := ref.Alias
		if table == "" {
			table = ref.Name
		}
		if table == "" {
			return nil
		}
		_, refNameQuoted, refAliasQuoted, hasExplicitAlias := r.refQuotedFlags(ref)
		quoted := refNameQuoted
		if hasExplicitAlias {
			quoted = refAliasQuoted
		}
		return &virtualTableRef{Table: table, Quoted: quoted}
	}
	switch ref.Kind {
	case parser.RangeReferenceRelation, parser.RangeReferenceDMLTarget, parser.RangeReferenceMergeSource:
		if r.hasDatabaseLink(ref) {
			table := ref.Alias
			if table == "" {
				table = ref.Name
			}
			if table == "" {
				return nil
			}
			_, refNameQuoted, refAliasQuoted, hasExplicitAlias := r.refQuotedFlags(ref)
			quoted := refNameQuoted
			if hasExplicitAlias {
				quoted = refAliasQuoted
			}
			return &virtualTableRef{Table: table, Quoted: quoted}
		}
		schemaQuoted, tableQuoted, aliasQuoted, hasExplicitAlias := r.refQuotedFlags(ref)
		alias := ""
		if hasExplicitAlias {
			alias = ref.Alias
		}
		return &physicalTableRef{
			Schema:       ref.Schema,
			SchemaQuoted: schemaQuoted,
			Table:        ref.Name,
			TableQuoted:  tableQuoted,
			Alias:        alias,
			AliasQuoted:  aliasQuoted,
		}
	case parser.RangeReferenceCTE, parser.RangeReferenceSubquery, parser.RangeReferenceJoinAlias:
		_, refNameQuoted, refAliasQuoted, hasExplicitAlias := r.refQuotedFlags(ref)
		quoted := refNameQuoted
		if hasExplicitAlias {
			quoted = refAliasQuoted
		}
		virtual := convertCompletionVirtualReference(ref, quoted)
		if virtual != nil && (ref.Kind == parser.RangeReferenceSubquery || len(virtual.Columns) == 0) {
			if columns := r.inferCompletionVirtualColumns(ref); len(columns) > 0 {
				virtual.Columns = columns
			}
		}
		return virtual
	default:
		return nil
	}
}

func (r *resolver) hasDatabaseLink(ref parser.RangeReference) bool {
	if ref.Loc.Start < 0 {
		return false
	}
	for _, tok := range r.tokens {
		if tok.Loc < ref.Loc.Start {
			continue
		}
		if ref.Loc.End >= 0 && tok.Loc > ref.Loc.End {
			break
		}
		if tok.Type == '@' {
			return true
		}
	}
	return false
}

func (r *resolver) refQuotedFlags(ref parser.RangeReference) (schemaQuoted, tableQuoted, aliasQuoted, hasExplicitAlias bool) {
	if ref.Loc.Start < 0 || ref.Loc.End <= ref.Loc.Start || ref.Loc.End > len(r.sql) {
		return false, false, false, false
	}
	var toks []parser.Token
	for _, tok := range r.tokens {
		if tok.Loc < ref.Loc.Start {
			continue
		}
		if tok.Loc >= ref.Loc.End {
			break
		}
		toks = append(toks, tok)
	}
	if ref.Kind == parser.RangeReferenceSubquery && ref.BodyLoc.End >= 0 {
		for k := len(toks) - 1; k >= 0; k-- {
			if toks[k].Loc > ref.BodyLoc.End && toks[k].Str != "" && strings.EqualFold(toks[k].Str, ref.Alias) {
				aliasQuoted = r.isQuotedToken(toks[k])
				hasExplicitAlias = true
				break
			}
		}
		return false, aliasQuoted, aliasQuoted, hasExplicitAlias
	}
	tableIdx := -1
	if ref.Schema != "" {
		for k := 0; k+2 < len(toks); k++ {
			if toks[k+1].Type == '.' && strings.EqualFold(toks[k].Str, ref.Schema) {
				schemaQuoted = r.isQuotedToken(toks[k])
				tableQuoted = r.isQuotedToken(toks[k+2])
				tableIdx = k + 2
				break
			}
		}
	} else {
		for k, tok := range toks {
			if tok.Str != "" && strings.EqualFold(tok.Str, ref.Name) {
				tableQuoted = r.isQuotedToken(tok)
				tableIdx = k
				break
			}
		}
	}
	if ref.Alias != "" {
		for k := len(toks) - 1; k >= 0; k-- {
			if k == tableIdx {
				continue
			}
			if toks[k].Str != "" && strings.EqualFold(toks[k].Str, ref.Alias) {
				aliasQuoted = r.isQuotedToken(toks[k])
				hasExplicitAlias = true
				break
			}
		}
	}
	return schemaQuoted, tableQuoted, aliasQuoted, hasExplicitAlias
}

func (r *resolver) isQuotedToken(tok parser.Token) bool {
	return tok.Loc >= 0 && tok.Loc < len(r.sql) && r.sql[tok.Loc] == '"'
}

func (r *resolver) insertAliasColumns(columnsOut *[]Candidate, alias string, aliasQuoted bool) {
	for _, level := range r.scopeLevels {
		matched := false
		for _, ref := range level.scopeRefs {
			if (ref.Unsupported && len(ref.Columns) == 0) || r.hasDatabaseLink(ref) {
				continue
			}
			_, refNameQuoted, refAliasQuoted, hasExplicitAlias := r.refQuotedFlags(ref)
			if hasExplicitAlias {
				if !nameMatches(ref.Alias, alias, refAliasQuoted || aliasQuoted) {
					continue
				}
			} else if !nameMatches(ref.Name, alias, refNameQuoted || aliasQuoted) {
				continue
			}
			matched = true
			columns := append([]string{}, ref.Columns...)
			if len(columns) == 0 {
				columns = r.inferCompletionVirtualColumns(ref)
			}
			for _, col := range columns {
				*columnsOut = append(*columnsOut, Candidate{Type: CandidateColumn, Text: col})
			}
		}
		if matched {
			return
		}
	}
}

func (r *resolver) insertLocalVariableCandidates(dst *[]Candidate) {
	if r.caretTokenIndex <= 0 || r.caretTokenIndex > len(r.tokens) {
		return
	}
	toks := r.tokens[:r.caretTokenIndex]

	type varScope struct {
		inDecl bool
		vars   []string
	}
	type blockMarker struct {
		kind     string
		hasScope bool
	}
	var scopes []varScope
	var blocks []blockMarker

	pushScope := func(params []string) {
		scopes = append(scopes, varScope{inDecl: true, vars: params})
	}
	addVar := func(name string) {
		if name == "" || len(scopes) == 0 {
			return
		}
		scopes[len(scopes)-1].vars = append(scopes[len(scopes)-1].vars, name)
	}

	i := 0
	for i < len(toks) {
		tok := toks[i]
		kw := strings.ToUpper(parser.TokenName(tok.Type))
		switch kw {
		case "PROCEDURE", "FUNCTION":
			// Only start a subprogram scope when declaring/defining a subprogram header with IS/AS.
			if i+1 < len(toks) && r.isIdentifierToken(toks[i+1]) {
				j := i + 2
				var params []string
				if j < len(toks) && toks[j].Type == '(' {
					closeIdx := findMatchingTokenIdx(toks, j)
					params = r.extractParamNames(toks[j+1 : closeIdx])
					j = closeIdx
					if j < len(toks) {
						j++
					}
				}
				foundHeaderEnd := false
				for j < len(toks) {
					jkw := strings.ToUpper(parser.TokenName(toks[j].Type))
					if toks[j].Type == ';' {
						break
					}
					if jkw == "IS" || jkw == "AS" {
						pushScope(params)
						i = j
						foundHeaderEnd = true
						break
					}
					j++
				}
				if foundHeaderEnd {
					i++
					continue
				}
			}
		case "DECLARE":
			pushScope(nil)
		case "BEGIN":
			hasDeclScope := false
			if len(scopes) > 0 && scopes[len(scopes)-1].inDecl {
				scopes[len(scopes)-1].inDecl = false
				hasDeclScope = true
			}
			blocks = append(blocks, blockMarker{kind: "BEGIN", hasScope: hasDeclScope})
		case "CASE":
			if len(scopes) > 0 && !scopes[len(scopes)-1].inDecl {
				blocks = append(blocks, blockMarker{kind: "CASE"})
			}
		case "END":
			nextKw := ""
			if i+1 < len(toks) {
				nextKw = strings.ToUpper(parser.TokenName(toks[i+1].Type))
			}
			switch nextKw {
			case "IF", "LOOP":
				i++
			case "CASE":
				if len(blocks) > 0 && blocks[len(blocks)-1].kind == "CASE" {
					blocks = blocks[:len(blocks)-1]
				}
				i++
			default:
				if len(blocks) > 0 {
					top := blocks[len(blocks)-1]
					blocks = blocks[:len(blocks)-1]
					if top.kind == "BEGIN" && top.hasScope && len(scopes) > 0 {
						scopes = scopes[:len(scopes)-1]
					}
				}
			}
		default:
			if len(scopes) > 0 && scopes[len(scopes)-1].inDecl && r.isIdentifierToken(tok) &&
				(i == 0 || toks[i-1].Type == ';' || strings.ToUpper(parser.TokenName(toks[i-1].Type)) == "DECLARE" ||
					strings.ToUpper(parser.TokenName(toks[i-1].Type)) == "IS" || strings.ToUpper(parser.TokenName(toks[i-1].Type)) == "AS") {
				if i+1 < len(toks) && toks[i+1].Type != ';' &&
					strings.ToUpper(parser.TokenName(toks[i+1].Type)) != "EXCEPTION" &&
					toks[i+1].Type != '(' {
					addVar(tok.Str)
				}
				depth := 0
				for i < len(toks) {
					if toks[i].Type == '(' {
						depth++
					} else if toks[i].Type == ')' && depth > 0 {
						depth--
					} else if toks[i].Type == ';' && depth == 0 {
						break
					}
					i++
				}
			}
		}
		i++
	}

	for _, sc := range scopes {
		for _, v := range sc.vars {
			*dst = append(*dst, Candidate{Type: CandidateVariable, Text: v})
		}
	}
}

func (r *resolver) extractParamNames(paramToks []parser.Token) []string {
	var names []string
	depth := 0
	atItemStart := true
	for _, tok := range paramToks {
		switch tok.Type {
		case '(':
			depth++
			atItemStart = false
		case ')':
			if depth > 0 {
				depth--
			}
			atItemStart = false
		case ',':
			if depth == 0 {
				atItemStart = true
			}
		default:
			if depth == 0 && atItemStart {
				if r.isIdentifierToken(tok) {
					names = append(names, tok.Str)
				}
				atItemStart = false
			}
		}
	}
	return names
}

func findMatchingTokenIdx(toks []parser.Token, openIdx int) int {
	depth := 1
	for k := openIdx + 1; k < len(toks); k++ {
		switch toks[k].Type {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return k
			}
		}
	}
	return len(toks)
}

func (r *resolver) insertReferenceAliasCandidates(tablesOut *[]Candidate) {
	for _, ref := range r.references {
		switch ref := ref.(type) {
		case *physicalTableRef:
			if ref.Alias == "" || strings.EqualFold(ref.Alias, ref.Table) {
				continue
			}
			*tablesOut = append(*tablesOut, Candidate{Type: CandidateTable, Text: ref.Alias})
		case *virtualTableRef:
			*tablesOut = append(*tablesOut, Candidate{Type: CandidateTable, Text: ref.Table})
		}
	}
}

func convertCompletionVirtualReference(ref parser.RangeReference, quoted bool) *virtualTableRef {
	table := ref.Alias
	if table == "" {
		table = ref.Name
	}
	if table == "" {
		return nil
	}
	return &virtualTableRef{
		Table:   table,
		Quoted:  quoted,
		Columns: append([]string{}, ref.Columns...),
	}
}

func (r *resolver) inferCTEAnchorColumns(ref parser.RangeReference, unionLoc int) []string {
	if ref.BodyLoc.Start < 0 || unionLoc <= ref.BodyLoc.Start || unionLoc > len(r.sql) {
		return nil
	}
	anchor := r.sql[ref.BodyLoc.Start:unionLoc]
	return r.resolveVirtualColumnsFromSQL(anchor)
}

func (r *resolver) inferCompletionVirtualColumns(ref parser.RangeReference) []string {
	if ref.BodyLoc.Start < 0 || ref.BodyLoc.End <= ref.BodyLoc.Start || ref.BodyLoc.End > len(r.sql) {
		return nil
	}
	body := r.sql[ref.BodyLoc.Start:ref.BodyLoc.End]
	statement := body
	if ref.Kind == parser.RangeReferenceSubquery {
		alias := ref.Alias
		if alias == "" {
			alias = "x"
		}
		statement = fmt.Sprintf("SELECT * FROM (%s) %s", body, quoteIdentifierForSQL(alias))
		if withPrefix := r.withClausePrefix(); withPrefix != "" {
			statement = fmt.Sprintf("%s SELECT * FROM (%s) %s", withPrefix, body, quoteIdentifierForSQL(alias))
		}
	}
	if cols := r.resolveVirtualColumnsFromSQL(statement); len(cols) > 0 {
		return cols
	}
	if ref.Kind == parser.RangeReferenceCTE && ref.Name != "" {
		if withPrefix := r.withClausePrefix(); withPrefix != "" {
			return r.resolveVirtualColumnsFromSQL(fmt.Sprintf("%s SELECT * FROM %s", withPrefix, quoteIdentifierForSQL(ref.Name)))
		}
	}
	return nil
}

func (r *resolver) resolveVirtualColumnsFromSQL(statement string) []string {
	span, err := analysis.ResolveQuerySpan(statement, analysis.ResolveOptions{
		DefaultDatabase: r.store.defaultSchema,
		Catalog:         r.analysisCatalog(),
	})
	if err != nil || span == nil || span.NotFoundError != nil {
		return nil
	}
	columns := make([]string, 0, len(span.Results))
	for _, col := range span.Results {
		columns = append(columns, col.Name)
	}
	return columns
}

func (r *resolver) analysisCatalog() analysis.Catalog {
	return analysis.Catalog{
		GetDatabase: func(_, database string) (analysis.Database, error) {
			_, sm := r.store.getSchema(database)
			if sm == nil {
				return nil, nil
			}
			return &analysisDatabase{schema: sm}, nil
		},
	}
}

type analysisDatabase struct {
	schema *metadatapb.SchemaMetadata
}

func (d *analysisDatabase) HasSchema() bool { return d.schema != nil }

func (d *analysisDatabase) GetTable(name string) *metadatapb.TableMetadata {
	return findTableMeta(d.schema, name, false)
}

func (d *analysisDatabase) GetView(name string) *metadatapb.ViewMetadata {
	return findViewMeta(d.schema, name, false)
}

func (d *analysisDatabase) GetMaterializedView(name string) *metadatapb.MaterializedViewMetadata {
	return findMaterializedViewMeta(d.schema, name, false)
}

func (d *analysisDatabase) GetExternalTable(name string) *metadatapb.ExternalTableMetadata {
	return findExternalTableMeta(d.schema, name, false)
}

func (r *resolver) withClausePrefix() string {
	first := -1
	for i, tok := range r.tokens {
		if tok.Type != ';' {
			first = i
			break
		}
	}
	if first < 0 || r.tokens[first].Type != parser.WITH {
		return ""
	}
	depth := 0
	for i := first + 1; i < len(r.tokens); i++ {
		switch r.tokens[i].Type {
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		case parser.SELECT:
			if depth == 0 {
				return strings.TrimSpace(r.sql[:r.tokens[i].Loc])
			}
		}
	}
	return ""
}

func (r *resolver) appendVirtualReferenceIfMissing(ref *virtualTableRef) {
	for _, existing := range r.cteTables {
		if strings.EqualFold(existing.Table, ref.Table) {
			if len(existing.Columns) == 0 && len(ref.Columns) > 0 {
				existing.Columns = ref.Columns
			}
			return
		}
	}
	r.cteTables = append(r.cteTables, ref)
}

func (r *resolver) appendReferenceIfMissing(ref tableRef) {
	key := tableRefKey(ref)
	for _, existing := range r.references {
		if tableRefKey(existing) == key {
			return
		}
	}
	r.references = append(r.references, ref)
}

func (r *resolver) appendScopeReferenceIfMissing(ref parser.RangeReference) {
	key := scopeRefKey(ref)
	for _, existing := range r.scopeReferences {
		if scopeRefKey(existing) == key {
			return
		}
	}
	r.scopeReferences = append(r.scopeReferences, ref)
}

func tableRefKey(ref tableRef) string {
	switch ref := ref.(type) {
	case *physicalTableRef:
		return strings.Join([]string{"physical", ref.Schema, ref.Table, ref.Alias}, "\x00")
	case *virtualTableRef:
		return strings.Join([]string{"virtual", ref.Table}, "\x00")
	default:
		return fmt.Sprintf("%T", ref)
	}
}

func scopeRefKey(ref parser.RangeReference) string {
	return strings.Join([]string{
		fmt.Sprint(ref.Kind),
		ref.Schema,
		ref.Name,
		ref.Alias,
		fmt.Sprint(ref.Loc.Start),
		fmt.Sprint(ref.Loc.End),
	}, "\x00")
}

func (r *resolver) qualifierParts() []qualifierPart {
	offset := min(max(r.collectOffset, 0), len(r.sql))
	before := strings.TrimRight(r.sql[:offset], " \t\r\n")
	if !strings.HasSuffix(before, ".") {
		return nil
	}
	before = strings.TrimSuffix(before, ".")
	var parts []qualifierPart
	for {
		before = strings.TrimRight(before, " \t\r\n")
		part, start := lastIdentifierSpan(before)
		if part == "" {
			break
		}
		quoted := len(part) >= 2 && part[0] == '"' && part[len(part)-1] == '"'
		parts = append([]qualifierPart{{name: unquoteIdent(part), quoted: quoted}}, parts...)
		before = strings.TrimRight(before[:start], " \t\r\n")
		if !strings.HasSuffix(before, ".") {
			break
		}
		before = strings.TrimSuffix(before, ".")
	}
	return parts
}

func lastIdentifierSpan(s string) (string, int) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", 0
	}
	if s[len(s)-1] == '"' {
		for i := len(s) - 2; i >= 0; i-- {
			if s[i] == '"' {
				if i > 0 && s[i-1] == '"' {
					i--
					continue
				}
				return s[i:], i
			}
		}
		return "", 0
	}
	end := len(s)
	start := end
	for start > 0 {
		r, size := utf8.DecodeLastRuneInString(s[:start])
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '$' && r != '#' {
			break
		}
		start -= size
	}
	return s[start:end], start
}

func unquoteIdent(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return strings.ReplaceAll(s[1:len(s)-1], `""`, `"`)
	}
	return strings.ToUpper(s)
}

func (r *resolver) columnQualifier() (schema, table string, schemaQuoted, tableQuoted bool) {
	parts := r.qualifierParts()
	switch len(parts) {
	case 0:
		return "", "", false, false
	case 1:
		return "", parts[0].name, false, parts[0].quoted
	default:
		s := parts[len(parts)-2]
		t := parts[len(parts)-1]
		return s.name, t.name, s.quoted, t.quoted
	}
}

func (r *resolver) selectItemAliases() []string {
	if !r.isInAliasAllowedContext() {
		return nil
	}
	selectIdx, fromIdx := r.currentQueryBlockSelectFrom()
	if selectIdx < 0 || fromIdx < 0 {
		return nil
	}
	aliasMap := make(map[string]bool)
	start := selectIdx + 1
	depth := 0
	for i := start; i <= fromIdx; i++ {
		if i == fromIdx || (depth == 0 && r.tokens[i].Type == ',') {
			if alias := r.selectAliasFromTokens(r.tokens[start:i]); alias != "" {
				aliasMap[alias] = true
			}
			start = i + 1
			continue
		}
		switch r.tokens[i].Type {
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		}
	}
	aliases := make([]string, 0, len(aliasMap))
	for alias := range aliasMap {
		aliases = append(aliases, alias)
	}
	slices.Sort(aliases)
	return aliases
}

func (r *resolver) currentQueryBlockSelectFrom() (int, int) {
	stmtStart, stmtEnd := r.statementTokenBounds()
	targetDepth := r.depthAtToken(r.caretTokenIndex)
	selectIdx := -1
	fromIdx := -1
	firstArmSelectIdx := -1
	firstArmFromIdx := -1
	inSetOp := false
	inOrderByAfterSetOp := false
	depth := 0
	for i := stmtStart; i < stmtEnd && i < r.caretTokenIndex; i++ {
		if r.tokens[i].Type == ')' && depth > 0 {
			depth--
		}
		if depth < targetDepth {
			selectIdx, fromIdx = -1, -1
			firstArmSelectIdx, firstArmFromIdx = -1, -1
			inSetOp, inOrderByAfterSetOp = false, false
		}
		if depth == targetDepth {
			switch strings.ToUpper(parser.TokenName(r.tokens[i].Type)) {
			case "SELECT":
				selectIdx = i
				fromIdx = -1
				if firstArmSelectIdx < 0 {
					firstArmSelectIdx = i
				}
			case "FROM":
				if selectIdx >= 0 && fromIdx < 0 {
					fromIdx = i
					if firstArmFromIdx < 0 {
						firstArmFromIdx = i
					}
				}
			case "UNION", "INTERSECT", "MINUS", "EXCEPT":
				inSetOp = true
				inOrderByAfterSetOp = false
			case "ORDER":
				if inSetOp {
					inOrderByAfterSetOp = true
				}
			}
		}
		if r.tokens[i].Type == '(' {
			depth++
		}
	}
	if inOrderByAfterSetOp && firstArmSelectIdx >= 0 && firstArmFromIdx >= 0 {
		return firstArmSelectIdx, firstArmFromIdx
	}
	return selectIdx, fromIdx
}

func (r *resolver) depthAtToken(index int) int {
	stmtStart, _ := r.statementTokenBounds()
	depth := 0
	for i := stmtStart; i < index && i < len(r.tokens); i++ {
		switch r.tokens[i].Type {
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		}
	}
	return depth
}

func (r *resolver) selectAliasFromTokens(tokens []parser.Token) string {
	for len(tokens) > 0 && tokens[0].Type == ',' {
		tokens = tokens[1:]
	}
	for len(tokens) > 0 && tokens[len(tokens)-1].Type == ',' {
		tokens = tokens[:len(tokens)-1]
	}
	if len(tokens) < 2 {
		return ""
	}
	last := tokens[len(tokens)-1]
	if !r.isIdentifierToken(last) {
		return ""
	}
	prev := tokens[len(tokens)-2]
	if prev.Type == parser.AS {
		return last.Str
	}
	if r.canPrecedeImplicitAlias(prev) {
		return last.Str
	}
	return ""
}

func (r *resolver) isIdentifierToken(tok parser.Token) bool {
	if tok.Type == parser.IDENT {
		return tok.Str != ""
	}
	if parser.TokenName(tok.Type) == "" && tok.Str != "" && tok.Loc >= 0 && tok.Loc < len(r.sql) {
		return r.sql[tok.Loc] == '"'
	}
	return false
}

func (r *resolver) canPrecedeImplicitAlias(tok parser.Token) bool {
	if tok.Type == ')' {
		return true
	}
	if name := parser.TokenName(tok.Type); name != "" {
		switch name {
		case "NULL", "TRUE", "FALSE", "SYSDATE", "SYSTIMESTAMP", "USER", "UID", "ROWNUM", "ROWID", "LEVEL", "END":
			return true
		case "DISTINCT", "ALL", "UNIQUE", "PRIOR", "CONNECT_BY_ROOT", "NOT", "EXISTS", "CASE", "WHEN", "THEN", "ELSE", "AND", "OR", "LIKE", "IN", "BETWEEN", "IS", "AS":
			return false
		default:
			return tok.Type >= 2000 && !parser.IsReservedKeyword(name)
		}
	}
	if tok.Loc >= 0 && tok.Loc < len(r.sql) {
		c := r.sql[tok.Loc]
		return unicode.IsLetter(rune(c)) || unicode.IsDigit(rune(c)) || c == '\'' || c == '"' || c == ':'
	}
	return false
}

func (r *resolver) isInAliasAllowedContext() bool {
	stmtStart, _ := r.statementTokenBounds()
	for i := r.caretTokenIndex - 1; i >= stmtStart; i-- {
		switch r.tokens[i].Type {
		case parser.ORDER, parser.GROUP:
			return true
		case parser.WHERE, parser.FROM, parser.SELECT:
			return false
		}
	}
	return false
}

func (r *resolver) statementTokenBounds() (int, int) {
	start := 0
	for i, tok := range r.tokens {
		if tok.Loc >= r.collectOffset {
			break
		}
		if tok.Type == ';' {
			start = i + 1
		}
	}
	end := len(r.tokens)
	for i := start; i < len(r.tokens); i++ {
		if r.tokens[i].Loc < r.collectOffset {
			continue
		}
		if r.tokens[i].Type == ';' {
			end = i
			break
		}
	}
	return start, end
}
