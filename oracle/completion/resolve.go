package completion

import (
	"fmt"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oracle/analysis"
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
	if sm, ok := s.schemas[name]; ok {
		return name, sm
	}
	for k, sm := range s.schemas {
		if strings.EqualFold(k, name) {
			return k, sm
		}
	}
	if name == "" {
		if s.defaultSchema != "" {
			return s.getSchema(s.defaultSchema)
		}
		return "", nil
	}

	targetName := name
	for _, known := range s.schemaOrder {
		if known == name {
			targetName = known
			break
		}
		if strings.EqualFold(known, name) {
			targetName = known
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
	for k, sm := range s.schemas {
		if strings.EqualFold(k, targetName) {
			return k, sm
		}
	}
	s.loadedMisses[targetName] = true
	return "", nil
}

type tableRef interface {
	isTableRef()
}

type physicalTableRef struct {
	Schema string
	Table  string
	Alias  string
}

func (*physicalTableRef) isTableRef() {}

type virtualTableRef struct {
	Table   string
	Columns []string
}

func (*virtualTableRef) isTableRef() {}

type tableTarget struct {
	schema string
	table  string
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
		types       []Candidate
	)

	addSchemaCandidates := func() {
		for _, name := range r.store.listSchemas() {
			schemas = append(schemas, Candidate{Type: CandidateSchema, Text: name})
		}
	}

	addRelationCandidates := func(schemaSet map[string]bool, includeTables, includeViews, includeMViews, includeSeqs, includeSyns bool) {
		includePublicSyns := false
		for schema := range schemaSet {
			if schema == "" && len(r.cteTables) > 0 && includeTables {
				for _, cte := range r.cteTables {
					tables = append(tables, Candidate{Type: CandidateTable, Text: cte.Table})
				}
			}
			if schema == "" || strings.EqualFold(schema, r.store.defaultSchema) {
				includePublicSyns = true
			}
			if schema == "" && r.store.defaultSchema != "" {
				continue
			}
			_, sm := r.store.getSchema(schema)
			if sm == nil {
				continue
			}
			if includeTables {
				for _, t := range sm.GetTables() {
					tables = append(tables, Candidate{Type: CandidateTable, Text: t.GetName()})
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
				r.insertSchemaFunctions(&functions, r.completionSchemaQualifier())
			case "schema_ref":
				addSchemaCandidates()
			case "table_ref":
				switch {
				case r.completionIntentOnly(parser.ObjectKindTable):
					schema := r.completionSchemaQualifier()
					if schema == "" {
						addSchemaCandidates()
					}
					addRelationCandidates(r.completionSchemas(schema), true, false, false, false, false)
				case r.completionIntentOnly(parser.ObjectKindView):
					schema := r.completionSchemaQualifier()
					if schema == "" {
						addSchemaCandidates()
					}
					addRelationCandidates(r.completionSchemas(schema), false, true, false, false, false)
				default:
					schema := r.completionSchemaQualifier()
					if schema == "" {
						addSchemaCandidates()
					}
					addRelationCandidates(r.completionSchemas(schema), true, true, true, true, true)
				}
			case "sequence_ref":
				addRelationCandidates(r.sequenceCompletionSchemas(), false, false, false, true, false)
			case "sequence_member_ref":
				r.insertSequenceMemberCandidates(&keywords)
			case "synonym_ref":
				schema := r.completionSchemaQualifier()
				addRelationCandidates(r.completionSchemas(schema), false, false, false, false, true)
			case "database_link_ref":
				r.insertDatabaseLinkCandidates(&dbLinks)
			case "proc_ref":
				r.insertSchemaProcedures(&procs, r.completionSchemaQualifier())
			case "package_ref":
				r.insertSchemaPackages(&packages, r.completionSchemaQualifier())
			case "package_member_ref":
				r.insertPackageMemberCandidates(&procs, &functions)
			case "index_ref":
				r.insertSchemaIndexes(&indexes, r.completionSchemaQualifier())
			case "trigger_ref":
				r.insertSchemaTriggers(&triggers, r.completionSchemaQualifier())
			case "constraint_ref":
				r.insertConstraintCandidates(&constraints)
			case "type_name":
				for _, typ := range oracleDataTypes {
					types = append(types, Candidate{Type: CandidateType_, Text: typ})
				}
			case "columnref":
				r.insertColumnContextCandidates(&tables, &columns, addSchemaCandidates, addRelationCandidates)
			}
		}
	}

	if r.completionIntent != nil {
		for _, kind := range r.completionIntent.ObjectKinds {
			switch kind {
			case parser.ObjectKindSchema:
				addSchemaCandidates()
			case parser.ObjectKindTable:
				schema := r.completionSchemaQualifier()
				if schema == "" {
					addSchemaCandidates()
				}
				addRelationCandidates(r.completionSchemas(schema), true, false, false, false, false)
			case parser.ObjectKindView:
				schema := r.completionSchemaQualifier()
				if schema == "" {
					addSchemaCandidates()
				}
				addRelationCandidates(r.completionSchemas(schema), false, true, false, false, false)
			case parser.ObjectKindSequence:
				addRelationCandidates(r.sequenceCompletionSchemas(), false, false, false, true, false)
			case parser.ObjectKindSequenceMember:
				r.insertSequenceMemberCandidates(&keywords)
			case parser.ObjectKindSynonym:
				schema := r.completionSchemaQualifier()
				addRelationCandidates(r.completionSchemas(schema), false, false, false, false, true)
			case parser.ObjectKindDatabaseLink:
				r.insertDatabaseLinkCandidates(&dbLinks)
			case parser.ObjectKindColumn:
				r.insertColumnContextCandidates(&tables, &columns, addSchemaCandidates, addRelationCandidates)
			case parser.ObjectKindFunction:
				r.insertSchemaFunctions(&functions, r.completionSchemaQualifier())
			case parser.ObjectKindProcedure:
				r.insertSchemaProcedures(&procs, r.completionSchemaQualifier())
			case parser.ObjectKindPackage:
				r.insertSchemaPackages(&packages, r.completionSchemaQualifier())
			case parser.ObjectKindPackageMember:
				r.insertPackageMemberCandidates(&procs, &functions)
			case parser.ObjectKindIndex:
				r.insertSchemaIndexes(&indexes, r.completionSchemaQualifier())
			case parser.ObjectKindTrigger:
				r.insertSchemaTriggers(&triggers, r.completionSchemaQualifier())
			case parser.ObjectKindConstraint:
				r.insertConstraintCandidates(&constraints)
			case parser.ObjectKindType:
				for _, typ := range oracleDataTypes {
					types = append(types, Candidate{Type: CandidateType_, Text: typ})
				}
			}
		}
	}

	if len(r.qualifierParts()) > 0 && len(columns) == 0 {
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
	return result
}

func (r *resolver) completionIntentOnly(kind parser.ObjectKind) bool {
	return r.completionIntent != nil && len(r.completionIntent.ObjectKinds) == 1 && r.completionIntent.ObjectKinds[0] == kind
}

func (r *resolver) completionSchemaQualifier() string {
	if r.completionIntent != nil && r.completionIntent.Qualifier.Schema != "" {
		return r.completionIntent.Qualifier.Schema
	}
	parts := r.qualifierParts()
	if len(parts) == 1 {
		return parts[0]
	}
	return ""
}

func (r *resolver) sequenceCompletionSchemas() map[string]bool {
	return r.completionSchemas(r.completionSchemaQualifier())
}

func (r *resolver) completionSchemas(schema string) map[string]bool {
	schemas := make(map[string]bool)
	if schema != "" {
		schemas[schema] = true
		return schemas
	}
	if r.store.defaultSchema != "" {
		schemas[r.store.defaultSchema] = true
	}
	schemas[""] = true
	return schemas
}

func (r *resolver) insertDatabaseLinkCandidates(dst *[]Candidate) {
	for _, link := range r.store.linkedDatabases {
		*dst = append(*dst, Candidate{Type: CandidateDatabaseLink, Text: link.GetName()})
	}
}

func (r *resolver) insertSchemaFunctions(dst *[]Candidate, schema string) {
	for s := range r.completionSchemas(schema) {
		if s == "" && r.store.defaultSchema != "" {
			continue
		}
		_, sm := r.store.getSchema(s)
		if sm == nil {
			continue
		}
		for _, fn := range sm.GetFunctions() {
			*dst = append(*dst, Candidate{Type: CandidateFunction, Text: fn.GetName(), Comment: fn.GetComment()})
		}
	}
}

func (r *resolver) insertSchemaProcedures(dst *[]Candidate, schema string) {
	for s := range r.completionSchemas(schema) {
		if s == "" && r.store.defaultSchema != "" {
			continue
		}
		_, sm := r.store.getSchema(s)
		if sm == nil {
			continue
		}
		for _, p := range sm.GetProcedures() {
			*dst = append(*dst, Candidate{Type: CandidateProcedure, Text: p.GetName(), Comment: p.GetComment()})
		}
	}
}

func (r *resolver) insertSchemaPackages(dst *[]Candidate, schema string) {
	for s := range r.completionSchemas(schema) {
		if s == "" && r.store.defaultSchema != "" {
			continue
		}
		_, sm := r.store.getSchema(s)
		if sm == nil {
			continue
		}
		for _, pkg := range sm.GetPackages() {
			*dst = append(*dst, Candidate{Type: CandidatePackage, Text: pkg.GetName()})
		}
	}
}

func (r *resolver) insertSchemaIndexes(dst *[]Candidate, schema string) {
	for s := range r.completionSchemas(schema) {
		if s == "" && r.store.defaultSchema != "" {
			continue
		}
		_, sm := r.store.getSchema(s)
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
	schema, pkgName := r.columnQualifier()
	if pkgName == "" && r.completionIntent != nil {
		pkgName = r.completionIntent.Qualifier.Object
	}
	if pkgName == "" {
		return
	}
	for s := range r.completionSchemas(schema) {
		if s == "" && r.store.defaultSchema != "" {
			continue
		}
		_, sm := r.store.getSchema(s)
		if sm == nil {
			continue
		}
		for _, pkg := range sm.GetPackages() {
			if strings.EqualFold(pkg.GetName(), pkgName) {
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
		if schema == "" {
			schema = r.store.defaultSchema
		}
		_, sm := r.store.getSchema(schema)
		if sm == nil {
			continue
		}
		t := findTableMeta(sm, phys.Table)
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
			if idx.GetIsConstraint() && idx.GetName() != "" {
				*dst = append(*dst, Candidate{Type: CandidateConstraint, Text: idx.GetName()})
			}
		}
	}
}

func (r *resolver) insertSchemaTriggers(dst *[]Candidate, schema string) {
	for s := range r.completionSchemas(schema) {
		if s == "" && r.store.defaultSchema != "" {
			continue
		}
		_, sm := r.store.getSchema(s)
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
	schema, seqName := r.columnQualifier()
	if seqName == "" && r.completionIntent != nil {
		seqName = r.completionIntent.Qualifier.Object
	}
	if seqName == "" {
		return
	}
	for s := range r.completionSchemas(schema) {
		if s == "" && r.store.defaultSchema != "" {
			continue
		}
		_, sm := r.store.getSchema(s)
		if sm == nil {
			continue
		}
		for _, seq := range sm.GetSequences() {
			if strings.EqualFold(seq.GetName(), seqName) {
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
	addRelationCandidates func(schemaSet map[string]bool, includeTables, includeViews, includeMViews, includeSeqs, includeSyns bool),
) {
	schema, table := r.columnQualifier()
	if table == "" && r.completionIntent != nil {
		table = r.completionIntent.Qualifier.Object
	}

	schemas := make(map[string]bool)
	if schema != "" {
		schemas[schema] = true
	} else if table == "" {
		for _, ref := range r.references {
			if physical, ok := ref.(*physicalTableRef); ok && physical.Schema != "" {
				schemas[physical.Schema] = true
			}
		}
	}

	var targets []tableTarget
	addTarget := func(s, t string) {
		if t == "" {
			return
		}
		if s == "" {
			s = r.store.defaultSchema
		}
		for _, existing := range targets {
			if strings.EqualFold(existing.schema, s) && strings.EqualFold(existing.table, t) {
				return
			}
		}
		targets = append(targets, tableTarget{schema: s, table: t})
	}

	if table != "" {
		matchedRef := false
		for _, ref := range r.references {
			switch ref := ref.(type) {
			case *physicalTableRef:
				if schema != "" && ref.Schema != "" && !strings.EqualFold(ref.Schema, schema) {
					continue
				}
				if strings.EqualFold(ref.Alias, table) || strings.EqualFold(ref.Table, table) {
					targetSchema := ref.Schema
					if targetSchema == "" {
						targetSchema = schema
					}
					addTarget(targetSchema, ref.Table)
					matchedRef = true
				}
			case *virtualTableRef:
				if schema == "" && strings.EqualFold(ref.Table, table) {
					for _, col := range ref.Columns {
						*columnsOut = append(*columnsOut, Candidate{Type: CandidateColumn, Text: col})
					}
					matchedRef = true
				}
			}
		}
		if !matchedRef {
			// Fallback for DDL/standalone qualified names (e.g. COMMENT ON COLUMN schema1.t2.|).
			addTarget(schema, table)
		}
	} else if len(r.references) > 0 {
		for _, ref := range r.references {
			switch ref := ref.(type) {
			case *physicalTableRef:
				addTarget(ref.Schema, ref.Table)
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
		schemas[r.store.defaultSchema] = true
	}
	schemas[""] = true

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
		r.insertAliasColumns(columnsOut, table)
	}
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
		if target.schema == "" {
			for _, cte := range r.cteTables {
				if strings.EqualFold(cte.Table, target.table) {
					for _, col := range cte.Columns {
						*columnsOut = append(*columnsOut, Candidate{Type: CandidateColumn, Text: col})
					}
				}
			}
		}
		r.insertColumnsForObject(columnsOut, target.schema, target.table, 0)
	}
}

func (r *resolver) insertColumnsForObject(columnsOut *[]Candidate, schema, name string, depth int) {
	if depth > 4 || name == "" {
		return
	}
	resolvedSchema, sm := r.store.getSchema(schema)
	if sm != nil {
		if t := findTableMeta(sm, name); t != nil {
			for _, col := range t.GetColumns() {
				*columnsOut = append(*columnsOut, columnCandidate(resolvedSchema, t.GetName(), col))
			}
			return
		}
		if ext := findExternalTableMeta(sm, name); ext != nil {
			for _, col := range ext.GetColumns() {
				*columnsOut = append(*columnsOut, columnCandidate(resolvedSchema, ext.GetName(), col))
			}
			return
		}
		if v := findViewMeta(sm, name); v != nil {
			for _, col := range v.GetColumns() {
				*columnsOut = append(*columnsOut, columnCandidate(resolvedSchema, v.GetName(), col))
			}
			return
		}
		if syn := findSynonymMeta(sm.GetSynonyms(), name); syn != nil && syn.GetTargetName() != "" {
			targetSchema := syn.GetTargetSchema()
			if targetSchema == "" {
				targetSchema = resolvedSchema
			}
			r.insertColumnsForObject(columnsOut, targetSchema, syn.GetTargetName(), depth+1)
			return
		}
	}
	if schema == "" || strings.EqualFold(schema, r.store.defaultSchema) {
		if syn := findSynonymMeta(r.store.publicSynonyms, name); syn != nil && syn.GetTargetName() != "" {
			targetSchema := syn.GetTargetSchema()
			if targetSchema == "" {
				targetSchema = r.store.defaultSchema
			}
			r.insertColumnsForObject(columnsOut, targetSchema, syn.GetTargetName(), depth+1)
		}
	}
}

func findExternalTableMeta(sm *metadatapb.SchemaMetadata, name string) *metadatapb.ExternalTableMetadata {
	for _, ext := range sm.GetExternalTables() {
		if ext.GetName() == name {
			return ext
		}
	}
	for _, ext := range sm.GetExternalTables() {
		if strings.EqualFold(ext.GetName(), name) {
			return ext
		}
	}
	return nil
}

func findSynonymMeta(syns []*metadatapb.SynonymMetadata, name string) *metadatapb.SynonymMetadata {
	for _, syn := range syns {
		if syn.GetName() == name {
			return syn
		}
	}
	for _, syn := range syns {
		if strings.EqualFold(syn.GetName(), name) {
			return syn
		}
	}
	return nil
}

func findTableMeta(sm *metadatapb.SchemaMetadata, name string) *metadatapb.TableMetadata {
	for _, t := range sm.GetTables() {
		if t.GetName() == name {
			return t
		}
	}
	for _, t := range sm.GetTables() {
		if strings.EqualFold(t.GetName(), name) {
			return t
		}
	}
	return nil
}

func findViewMeta(sm *metadatapb.SchemaMetadata, name string) *metadatapb.ViewMetadata {
	for _, v := range sm.GetViews() {
		if v.GetName() == name {
			return v
		}
	}
	for _, v := range sm.GetViews() {
		if strings.EqualFold(v.GetName(), name) {
			return v
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
		// Only expose CTEs completely defined before the cursor so completing
		// inside an earlier CTE body does not suggest later CTEs.
		if ref.Loc.End >= 0 && r.collectOffset < ref.Loc.End {
			continue
		}
		virtual := convertCompletionVirtualReference(ref)
		if virtual == nil {
			continue
		}
		if len(virtual.Columns) == 0 {
			virtual.Columns = r.inferCompletionVirtualColumns(ref)
		}
		r.appendVirtualReferenceIfMissing(virtual)
	}
}

func (r *resolver) collectCompletionScopeReferences(ctx *parser.CompletionContext) {
	if ctx.Scope == nil {
		return
	}
	appendRef := func(ref parser.RangeReference) {
		if converted := r.convertCompletionScopeReference(ref); converted != nil {
			r.appendReferenceIfMissing(converted)
			r.appendScopeReferenceIfMissing(ref)
		}
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

func (r *resolver) convertCompletionScopeReference(ref parser.RangeReference) tableRef {
	switch ref.Kind {
	case parser.RangeReferenceRelation, parser.RangeReferenceDMLTarget, parser.RangeReferenceMergeSource:
		return &physicalTableRef{
			Schema: ref.Schema,
			Table:  ref.Name,
			Alias:  ref.Alias,
		}
	case parser.RangeReferenceCTE, parser.RangeReferenceSubquery, parser.RangeReferenceJoinAlias:
		virtual := convertCompletionVirtualReference(ref)
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

func (r *resolver) insertAliasColumns(columnsOut *[]Candidate, alias string) {
	for _, ref := range r.scopeReferences {
		if !strings.EqualFold(ref.Alias, alias) && !strings.EqualFold(ref.Name, alias) {
			continue
		}
		columns := append([]string{}, ref.Columns...)
		if len(columns) == 0 {
			columns = r.inferCompletionVirtualColumns(ref)
		}
		for _, col := range columns {
			*columnsOut = append(*columnsOut, Candidate{Type: CandidateColumn, Text: col})
		}
		return
	}
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

func convertCompletionVirtualReference(ref parser.RangeReference) *virtualTableRef {
	table := ref.Alias
	if table == "" {
		table = ref.Name
	}
	if table == "" {
		return nil
	}
	return &virtualTableRef{
		Table:   table,
		Columns: append([]string{}, ref.Columns...),
	}
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
	span, err := analysis.ResolveQuerySpan(statement, analysis.ResolveOptions{
		DefaultDatabase: r.store.defaultSchema,
		Catalog:         r.analysisCatalog(),
	})
	if (err != nil || span == nil || span.NotFoundError != nil) && ref.Kind == parser.RangeReferenceCTE && ref.Name != "" {
		if withPrefix := r.withClausePrefix(); withPrefix != "" {
			statement = fmt.Sprintf("%s SELECT * FROM %s", withPrefix, quoteIdentifierForSQL(ref.Name))
			span, err = analysis.ResolveQuerySpan(statement, analysis.ResolveOptions{
				DefaultDatabase: r.store.defaultSchema,
				Catalog:         r.analysisCatalog(),
			})
		}
	}
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
	return findTableMeta(d.schema, name)
}

func (d *analysisDatabase) GetView(name string) *metadatapb.ViewMetadata {
	return findViewMeta(d.schema, name)
}

func (d *analysisDatabase) GetMaterializedView(name string) *metadatapb.MaterializedViewMetadata {
	for _, mv := range d.schema.GetMaterializedViews() {
		if mv.GetName() == name {
			return mv
		}
	}
	for _, mv := range d.schema.GetMaterializedViews() {
		if strings.EqualFold(mv.GetName(), name) {
			return mv
		}
	}
	return nil
}

func (d *analysisDatabase) GetExternalTable(name string) *metadatapb.ExternalTableMetadata {
	for _, ext := range d.schema.GetExternalTables() {
		if ext.GetName() == name {
			return ext
		}
	}
	for _, ext := range d.schema.GetExternalTables() {
		if strings.EqualFold(ext.GetName(), name) {
			return ext
		}
	}
	return nil
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

func (r *resolver) qualifierParts() []string {
	offset := min(max(r.collectOffset, 0), len(r.sql))
	before := strings.TrimRight(r.sql[:offset], " \t\r\n")
	if !strings.HasSuffix(before, ".") {
		return nil
	}
	before = strings.TrimSuffix(before, ".")
	var parts []string
	for {
		before = strings.TrimRight(before, " \t\r\n")
		part, start := lastIdentifierSpan(before)
		if part == "" {
			break
		}
		parts = append([]string{unquoteIdent(part)}, parts...)
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

func (r *resolver) columnQualifier() (schema, table string) {
	parts := r.qualifierParts()
	switch len(parts) {
	case 0:
		return "", ""
	case 1:
		return "", parts[0]
	default:
		return parts[len(parts)-2], parts[len(parts)-1]
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
	depth := 0
	for i := stmtStart; i < stmtEnd && i < r.caretTokenIndex; i++ {
		if r.tokens[i].Type == ')' && depth > 0 {
			depth--
		}
		if depth == targetDepth {
			switch r.tokens[i].Type {
			case parser.SELECT:
				selectIdx = i
				fromIdx = -1
			case parser.FROM:
				if selectIdx >= 0 && fromIdx < 0 {
					fromIdx = i
				}
			}
		}
		if r.tokens[i].Type == '(' {
			depth++
		}
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
