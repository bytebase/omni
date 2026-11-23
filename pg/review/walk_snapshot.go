package review

import (
	"strings"

	"github.com/bytebase/omni/metadata"
	pganalysis "github.com/bytebase/omni/pg/analysis"
	"github.com/bytebase/omni/pg/ast"
	pgcatalog "github.com/bytebase/omni/pg/catalog"
	"github.com/bytebase/omni/pg/parser"
	"github.com/bytebase/omni/review/walk"
)

// The spaces WalkThrough keeps facts in. A relation, type, function, and
// procedure are held by a schema, a column, constraint, trigger, policy, and
// the primary key by a table; the primary key space has one key per table,
// named "".
const (
	spaceSchema     = "schema"
	spaceRelation   = "relation"
	spaceType       = "type"
	spaceColumn     = "column"
	spaceConstraint = "constraint"
	spacePrimaryKey = "primary key"
	spaceTrigger    = "trigger"
	spacePolicy     = "policy"
	spaceFunction   = "function"
	spaceProcedure  = "procedure"
)

// The kinds of a relation: every name in pg_class's namespace of a schema.
const (
	relTable     = "table" // ordinary or partitioned
	relView      = "view"
	relMatview   = "materialized view"
	relSequence  = "sequence"
	relIndex     = "index"
	relForeign   = "foreign table"
	relComposite = "composite type"
)

// The kinds of a column and a constraint.
const (
	kindColumn = "column"

	// Kinds in the type space: a domain, or a type known to be no domain.
	// A type of unknown kind has none.
	typDomain = "domain"
	typOther  = "type"

	conPrimary   = "primary key"
	conUnique    = "unique"
	conForeign   = "foreign key"
	conCheck     = "check"
	conExclusion = "exclusion"
)

// tempSchema is the name that stands for the session's temporary schema.
const tempSchema = "pg_temp"

// systemColumns are the system columns every table has, which a user
// column cannot be named after (PostgreSQL 12 and later; oid stopped
// being one when WITH OIDS went).
//
// pg: src/backend/catalog/heap.c — SysAtt
var systemColumns = map[string]bool{
	"tableoid": true, "cmax": true, "xmax": true, "cmin": true, "xmin": true, "ctid": true,
}

// pgSnapshot is a target's synced schema as WalkThrough trusts it. It
// answers Present for what Bytebase's sync lists and Absent only where
// that sync is complete; everywhere else it answers Unknown.
//
// What the sync leaves out, from bytebase/backend/plugin/db/pg/sync.go and
// backend/plugin/parser/pg/system_objects.go:
//   - the system schemas (pg_catalog, information_schema, pg_toast*,
//     pg_temp*, and the TimescaleDB, Citus, pg_cron, pg_partman, and
//     squeeze schemas);
//   - tables and views by name (systemTables and systemViews, every
//     name of which begins with pg_, and the g_columnar_,
//     google_db_advisor_, g_agg_stat_, and hypopg view and materialized
//     view prefixes of IsSystemView);
//   - the schema of a partition, which the metadata keeps under its
//     parent, and a partition's columns;
//   - the partitions of a partitioned partition in another schema than
//     the table it descends from: warpTablePartitions looks a partition's
//     own partitions up under that table's schema;
//   - domains, range types, and base types;
//   - NOT NULL constraints, which PostgreSQL 18 names (t_a_not_null).
//
// So a schema's absence is trusted unless its name is a system one, a
// relation's in a synced schema unless its name is one the sync drops or
// a partition's or a partition index's, which may sit in any schema, or
// the sync may have lost partitions, and a column's in a synced table or
// foreign table. A type's, a constraint's, and a primary key's absence is
// never trusted.
//
// Assumed: a relation the sync drops by its name or its schema holds no
// index, partition, trigger, or rule, which the sync drops with it. Such
// a relation is a table named like a system catalog, a view or
// materialized view named like an extension's, or a partition placed in a
// system schema.
type pgSnapshot struct {
	schemas   map[string]bool
	relations map[[2]string]string
	types     map[[2]string]bool
	// columns lists the columns of the tables and foreign tables whose
	// columns the sync lists in full.
	columns     map[[2]string]map[string]bool
	constraints map[[2]string]map[string]string
	primaryKeys map[[2]string]bool
	triggers    map[[2]string]map[string]bool
	functions   map[[2]string]bool
	procedures  map[[2]string]bool
	// ambigRoutines marks entries from SchemaMetadata.Functions with an
	// empty Definition whose kind (function vs procedure) is uncertain, so
	// both spaceFunction and spaceProcedure stay Unknown rather than assuming
	// one kind or claiming the other absent.
	ambigRoutines map[[2]string]bool
	// rangeConstructors marks pre-existing range and multirange constructor
	// names omitted from SchemaMetadata.Functions by pg/driver/sync.go.
	rangeConstructors map[[2]string]bool
	viewCols          map[[2]string]int
	// anywhere lists the partition and partition index names, Unknown in
	// every schema.
	anywhere map[string]bool
	// lost is true when the sync may have lost partitions, of any name in
	// any schema: no relation's absence is trusted then.
	lost bool
	// inSchema lists the relations, types, and routines each schema holds.
	inSchema map[string][]walk.Key

	// Synced detail maps populated from pg/driver/sync.go metadata (cleared
	// under OtherSpecs and invalidated per-object when altered in the change).
	colMeta     map[[2]string]map[string]syncedCol
	pkCols      map[[2]string][]string
	partKeyCols map[[2]string]map[string]bool
	fkRefs      map[[2]string][]syncedFKRef
	compAttrs   map[[2]string]map[string]bool
	indexDefs   map[[2]string]*syncedIndex
	funcSigs    map[[2]string]*syncedFuncInfo
	trigFuncs   map[[2]string][]syncedTrigRef
}

type syncedCol struct {
	typ       string // normalized lowercase ColumnMetadata.Type
	nullable  bool
	identity  bool
	generated bool
}

type syncedFKRef struct {
	fromSchema string
	fromTable  string
	conName    string
}

type syncedIndex struct {
	colNames     []string
	isExpr       []bool
	nkey         int
	hasExpr      bool
	accessMethod string
}

type syncedFuncInfo struct {
	complete             bool
	arities              map[int]bool
	anySetof             bool
	zeroArg              bool
	hasDefaultOrVariadic bool
}

type syncedTrigRef struct {
	tableSchema string
	tableName   string
	trigName    string
}

// newPgSnapshot reads a target's synced schema. A nil schema gives a
// snapshot that knows nothing.
func newPgSnapshot(db *metadata.DatabaseSchemaMetadata) *pgSnapshot {
	s := &pgSnapshot{
		schemas:           make(map[string]bool),
		relations:         make(map[[2]string]string),
		types:             make(map[[2]string]bool),
		columns:           make(map[[2]string]map[string]bool),
		constraints:       make(map[[2]string]map[string]string),
		primaryKeys:       make(map[[2]string]bool),
		triggers:          make(map[[2]string]map[string]bool),
		functions:         make(map[[2]string]bool),
		procedures:        make(map[[2]string]bool),
		ambigRoutines:     make(map[[2]string]bool),
		rangeConstructors: make(map[[2]string]bool),
		viewCols:          make(map[[2]string]int),
		anywhere:          make(map[string]bool),
		inSchema:          make(map[string][]walk.Key),
		colMeta:           make(map[[2]string]map[string]syncedCol),
		pkCols:            make(map[[2]string][]string),
		partKeyCols:       make(map[[2]string]map[string]bool),
		fkRefs:            make(map[[2]string][]syncedFKRef),
		compAttrs:         make(map[[2]string]map[string]bool),
		indexDefs:         make(map[[2]string]*syncedIndex),
		funcSigs:          make(map[[2]string]*syncedFuncInfo),
		trigFuncs:         make(map[[2]string][]syncedTrigRef),
	}
	if db == nil {
		return s
	}
	for _, schema := range db.GetSchemas() {
		name := schema.GetName()
		if trustedSchemaName(name) {
			s.schemas[name] = true
		}
		rel := func(n, kind string) {
			s.relations[[2]string{name, n}] = kind
			s.inSchema[name] = append(s.inSchema[name], relKey(name, n))
		}
		rowType := func(n string) {
			s.types[[2]string{name, n}] = true
			s.inSchema[name] = append(s.inSchema[name], typeKey(name, n))
		}
		for _, t := range schema.GetTables() {
			key := [2]string{name, t.GetName()}
			rel(t.GetName(), relTable)
			rowType(t.GetName())
			cols := make(map[string]bool, len(t.GetColumns()))
			var cmeta map[string]syncedCol
			for _, c := range t.GetColumns() {
				cols[c.GetName()] = true
				if c.GetType() != "" && c.GetPosition() > 0 {
					if cmeta == nil {
						cmeta = make(map[string]syncedCol, len(t.GetColumns()))
					}
					cmeta[c.GetName()] = syncedCol{
						typ:       strings.ToLower(strings.TrimSpace(c.GetType())),
						nullable:  c.GetNullable(),
						identity:  c.GetIdentityGeneration() != metadata.ColumnMetadata_IDENTITY_GENERATION_UNSPECIFIED,
						generated: c.GetGeneration() != nil,
					}
				}
			}
			s.columns[key] = cols
			if len(cmeta) > 0 {
				s.colMeta[key] = cmeta
			}
			cons := make(map[string]string)
			for _, idx := range t.GetIndexes() {
				rel(idx.GetName(), relIndex)
				if idxInfo := parseSyncedIndex(idx.GetDefinition()); idxInfo != nil {
					s.indexDefs[[2]string{name, idx.GetName()}] = idxInfo
					if idx.GetPrimary() && !idxInfo.hasExpr && idxInfo.nkey > 0 {
						s.pkCols[key] = idxInfo.colNames[:idxInfo.nkey]
					}
				}
				switch {
				case idx.GetPrimary():
					cons[idx.GetName()] = conPrimary
					s.primaryKeys[key] = true
				case idx.GetIsConstraint() && idx.GetUnique():
					cons[idx.GetName()] = conUnique
				}
			}
			for _, fk := range t.GetForeignKeys() {
				cons[fk.GetName()] = conForeign
				if fk.GetReferencedTable() != "" {
					refSchema := fk.GetReferencedSchema()
					if refSchema == "" {
						refSchema = name
					}
					refKey := [2]string{refSchema, fk.GetReferencedTable()}
					if refKey != key {
						s.fkRefs[refKey] = append(s.fkRefs[refKey], syncedFKRef{
							fromSchema: name,
							fromTable:  t.GetName(),
							conName:    fk.GetName(),
						})
					}
				}
			}
			for _, c := range t.GetCheckConstraints() {
				cons[c.GetName()] = conCheck
			}
			for _, x := range t.GetExcludeConstraints() {
				cons[x.GetName()] = conExclusion
			}
			s.constraints[key] = cons
			trigs := make(map[string]bool, len(t.GetTriggers()))
			for _, tr := range t.GetTriggers() {
				trigs[tr.GetName()] = true
				if fSchema, fName, ok := parseTriggerFunc(tr.GetBody()); ok {
					s.trigFuncs[[2]string{fSchema, fName}] = append(s.trigFuncs[[2]string{fSchema, fName}], syncedTrigRef{
						tableSchema: name,
						tableName:   t.GetName(),
						trigName:    tr.GetName(),
					})
				}
			}
			s.triggers[key] = trigs
			for _, p := range t.GetPartitions() {
				if pcols := parsePartKeyCols(p.GetExpression()); len(pcols) > 0 {
					s.partKeyCols[key] = pcols
					break
				}
			}
			s.addPartitions(name, t.GetPartitions())
		}
		for _, t := range schema.GetExternalTables() {
			rel(t.GetName(), relForeign)
			rowType(t.GetName())
			cols := make(map[string]bool, len(t.GetColumns()))
			for _, c := range t.GetColumns() {
				cols[c.GetName()] = true
			}
			s.columns[[2]string{name, t.GetName()}] = cols
		}
		for _, v := range schema.GetViews() {
			rel(v.GetName(), relView)
			rowType(v.GetName())
			if len(v.GetColumns()) > 0 {
				s.viewCols[[2]string{name, v.GetName()}] = len(v.GetColumns())
			}
		}
		for _, v := range schema.GetMaterializedViews() {
			rel(v.GetName(), relMatview)
			rowType(v.GetName())
			for _, idx := range v.GetIndexes() {
				rel(idx.GetName(), relIndex)
				if idxInfo := parseSyncedIndex(idx.GetDefinition()); idxInfo != nil {
					s.indexDefs[[2]string{name, idx.GetName()}] = idxInfo
				}
			}
		}
		for _, q := range schema.GetSequences() {
			rel(q.GetName(), relSequence)
		}
		for _, c := range schema.GetCompositeTypes() {
			rel(c.GetName(), relComposite)
			rowType(c.GetName())
			if len(c.GetAttributes()) > 0 {
				attrs := make(map[string]bool, len(c.GetAttributes()))
				for _, a := range c.GetAttributes() {
					attrs[a.GetName()] = true
				}
				s.compAttrs[[2]string{name, c.GetName()}] = attrs
			}
		}
		for _, e := range schema.GetEnumTypes() {
			rowType(e.GetName())
		}
		for _, rt := range schema.GetRangeTypes() {
			rowType(rt.GetName())
			s.rangeConstructors[[2]string{name, rt.GetName()}] = true
			if mr := rt.GetMultirangeName(); mr != "" {
				mrParts := identifierParts(mr)
				switch len(mrParts) {
				case 1:
					rowType(mrParts[0])
					s.rangeConstructors[[2]string{name, mrParts[0]}] = true
				case 2:
					s.types[[2]string{mrParts[0], mrParts[1]}] = true
					s.rangeConstructors[[2]string{mrParts[0], mrParts[1]}] = true
				}
			} else {
				s.rangeConstructors[[2]string{name, rt.GetName() + "_multirange"}] = true
			}
		}
		for _, fn := range schema.GetFunctions() {
			// pg/driver/sync.go stores procedures in SchemaMetadata.Functions
			// with their pg_get_functiondef() text (CREATE OR REPLACE PROCEDURE).
			// Classify them via kindOf; an empty Definition may be either kind.
			k := kindOf(fn)
			fk := [2]string{name, fn.GetName()}
			if k == kindFunction|kindProcedure {
				s.ambigRoutines[fk] = true
				if fi := s.funcSigs[fk]; fi != nil {
					fi.complete = false
				} else {
					s.funcSigs[fk] = &syncedFuncInfo{complete: false}
				}
				continue
			}
			if k&kindFunction != 0 {
				s.functions[fk] = true
				s.inSchema[name] = append(s.inSchema[name], funcKey(name, fn.GetName()))
				fi := s.funcSigs[fk]
				if fi == nil {
					fi = &syncedFuncInfo{complete: true, arities: make(map[int]bool)}
					s.funcSigs[fk] = fi
				}
				if fi.complete {
					if arity, setof, defOrVar, ok := parseFuncDef(fn.GetDefinition()); ok {
						fi.arities[arity] = true
						if arity == 0 {
							fi.zeroArg = true
						}
						if setof {
							fi.anySetof = true
						}
						if defOrVar {
							fi.hasDefaultOrVariadic = true
						}
					} else {
						fi.complete = false
					}
				}
			}
			if k&kindProcedure != 0 {
				s.procedures[fk] = true
				s.inSchema[name] = append(s.inSchema[name], procKey(name, fn.GetName()))
			}
		}
		// Mark any range/multirange constructor names incomplete in funcSigs.
		for rk := range s.rangeConstructors {
			if fi := s.funcSigs[rk]; fi != nil {
				fi.complete = false
			}
		}
		for _, pr := range schema.GetProcedures() {
			s.procedures[[2]string{name, pr.GetName()}] = true
			s.inSchema[name] = append(s.inSchema[name], procKey(name, pr.GetName()))
		}
	}
	return s
}

// addPartitions records the partitions of a table in schema root, at any
// depth. Bytebase's sync finds a partition's own partitions under root,
// whatever the partition's schema, so it loses them when that is another:
// unless a partition is known to be a leaf or in root, they may be lost.
//
// bytebase: backend/plugin/db/pg/sync.go — warpTablePartitions,
// listTablePartitionQuery (keyed by the parent's own schema)
func (s *pgSnapshot) addPartitions(root string, partitions []*metadata.TablePartitionMetadata) {
	for _, p := range partitions {
		s.anywhere[p.GetName()] = true
		for _, idx := range p.GetIndexes() {
			s.anywhere[idx.GetName()] = true
		}
		if schema, partitioned, ok := partitionPlace(p); !ok || partitioned && schema != root {
			s.lost = true
		}
		s.addPartitions(root, p.GetSubpartitions())
	}
}

// partitionPlace reads, from the definitions of a partition's indexes,
// whether the partition is partitioned itself and, if so, its schema.
// pg_get_indexdef writes ON ONLY for every index of a partitioned table
// and qualifies the table with its schema, as Bytebase's sync runs it on
// an empty search path (verified on PostgreSQL 12 and 17). A partition
// with no index tells neither.
func partitionPlace(p *metadata.TablePartitionMetadata) (schema string, partitioned, ok bool) {
	for _, idx := range p.GetIndexes() {
		def := idx.GetDefinition()
		if def == "" {
			continue
		}
		if !strings.Contains(def, " ON ONLY ") {
			return "", false, true
		}
		// Feature extraction from a definition the server wrote: one that
		// does not parse tells nothing.
		list, err := parser.Parse(def)
		if err != nil || list == nil || len(list.Items) != 1 {
			continue
		}
		raw, ok := list.Items[0].(*ast.RawStmt)
		if !ok {
			continue
		}
		stmt, ok := raw.Stmt.(*ast.IndexStmt)
		if !ok || stmt.Relation == nil || stmt.Relation.Relname != p.GetName() || stmt.Relation.Schemaname == "" {
			continue
		}
		return stmt.Relation.Schemaname, !stmt.Relation.Inh, true
	}
	return "", false, false
}

// trustedSchemaName reports whether the sync lists a schema of this name
// whenever it exists: every name but the system ones, which it skips, and
// the pg_ names PostgreSQL reserves for them.
func trustedSchemaName(name string) bool {
	if strings.HasPrefix(name, "pg_") {
		return false
	}
	switch name {
	case "information_schema", "rw_catalog",
		"timescaledb_information", "timescaledb_experimental",
		"_timescaledb_cache", "_timescaledb_catalog", "_timescaledb_internal", "_timescaledb_config",
		"squeeze", "citus", "citus_internal", "columnar", "columnar_internal", "cron", "partman":
		return false
	}
	return true
}

// trustedRelationName reports whether the sync lists a relation of this
// name in a synced schema whenever it exists.
//
// bytebase: backend/plugin/parser/pg/system_objects.go — IsSystemTable,
// IsSystemView
func trustedRelationName(name string) bool {
	for _, prefix := range []string{"pg_", "g_columnar_", "google_db_advisor_", "g_agg_stat_", "hypopg"} {
		if strings.HasPrefix(name, prefix) {
			return false
		}
	}
	return true
}

// trustedRoutineName reports whether Bytebase's sync lists a routine of this
// name in a synced schema whenever it exists: every name except system/built-in
// names and extension prefixes that pg/analysis.IsSystemFunction drops.
//
// bytebase: backend/plugin/parser/pg/system_objects.go — IsSystemFunction
func trustedRoutineName(name string) bool {
	if builtinCalls[strings.ToLower(name)] || pganalysis.IsSystemFunction(name, "") {
		return false
	}
	return !strings.HasPrefix(name, "pg_")
}

// isInternalConstraintTrigger reports whether name matches one of PostgreSQL's
// system-generated internal constraint trigger prefixes (tgisinternal = true),
// which pg/driver/sync.go omits from TableMetadata.Triggers.
//
// pg: src/backend/commands/tablecmds.c — createForeignKeyCheckTriggers;
// src/backend/catalog/index.c — index_constraint_create
func isInternalConstraintTrigger(name string) bool {
	return strings.HasPrefix(name, "RI_ConstraintTrigger_") ||
		strings.HasPrefix(name, "Unique_ConstraintTrigger_") ||
		strings.HasPrefix(name, "PK_ConstraintTrigger_")
}

// parseSyncedIndex parses a pg_get_indexdef() CREATE INDEX definition and
// extracts its column names and expression flags.
func parseSyncedIndex(def string) *syncedIndex {
	if def == "" {
		return nil
	}
	list, err := parser.Parse(def)
	if err != nil || list == nil || len(list.Items) != 1 {
		return nil
	}
	raw, ok := list.Items[0].(*ast.RawStmt)
	if !ok {
		return nil
	}
	stmt, ok := raw.Stmt.(*ast.IndexStmt)
	if !ok || stmt.IndexParams == nil || len(stmt.IndexParams.Items) == 0 {
		return nil
	}
	info := &syncedIndex{
		colNames:     make([]string, 0, len(stmt.IndexParams.Items)),
		isExpr:       make([]bool, 0, len(stmt.IndexParams.Items)),
		nkey:         len(stmt.IndexParams.Items),
		accessMethod: strings.ToLower(stmt.AccessMethod),
	}
	for _, item := range stmt.IndexParams.Items {
		elem, ok := item.(*ast.IndexElem)
		if !ok {
			return nil
		}
		if elem.Expr != nil {
			info.colNames = append(info.colNames, "")
			info.isExpr = append(info.isExpr, true)
			info.hasExpr = true
		} else {
			name, ok := walkName(elem.Name)
			if !ok {
				return nil
			}
			info.colNames = append(info.colNames, name)
			info.isExpr = append(info.isExpr, false)
		}
	}
	if stmt.IndexIncludingParams != nil {
		for _, item := range stmt.IndexIncludingParams.Items {
			elem, ok := item.(*ast.IndexElem)
			if !ok || elem.Name == "" {
				return nil
			}
			name, ok := walkName(elem.Name)
			if !ok {
				return nil
			}
			info.colNames = append(info.colNames, name)
			info.isExpr = append(info.isExpr, false)
		}
	}
	return info
}

// parsePartKeyCols extracts plain column names from pg_get_partkeydef()
// (e.g. "RANGE (a)" or "HASH (a, b)") via the PostgreSQL parser.
func parsePartKeyCols(expr string) map[string]bool {
	if expr == "" {
		return nil
	}
	list, err := parser.Parse("CREATE TABLE _t (_ int) PARTITION BY " + expr)
	if err != nil || list == nil || len(list.Items) != 1 {
		return nil
	}
	raw, ok := list.Items[0].(*ast.RawStmt)
	if !ok {
		return nil
	}
	stmt, ok := raw.Stmt.(*ast.CreateStmt)
	if !ok || stmt.Partspec == nil || stmt.Partspec.PartParams == nil {
		return nil
	}
	var out map[string]bool
	for _, item := range stmt.Partspec.PartParams.Items {
		pe, ok := item.(*ast.PartitionElem)
		if !ok || pe.Name == "" || pe.Expr != nil {
			continue
		}
		if name, ok := walkName(pe.Name); ok {
			if out == nil {
				out = make(map[string]bool)
			}
			out[name] = true
		}
	}
	return out
}

// parseTriggerFunc parses a pg_get_triggerdef() CREATE TRIGGER definition and
// returns the schema and name of the executed trigger function when qualified.
func parseTriggerFunc(body string) (schema, name string, ok bool) {
	if body == "" {
		return "", "", false
	}
	list, err := parser.Parse(body)
	if err != nil || list == nil || len(list.Items) != 1 {
		return "", "", false
	}
	raw, ok := list.Items[0].(*ast.RawStmt)
	if !ok {
		return "", "", false
	}
	stmt, ok := raw.Stmt.(*ast.CreateTrigStmt)
	if !ok || stmt.Funcname == nil {
		return "", "", false
	}
	parts := nameParts(stmt.Funcname)
	if len(parts) != 2 {
		return "", "", false
	}
	s, ok1 := walkName(parts[0])
	fn, ok2 := walkName(parts[1])
	if !ok1 || !ok2 {
		return "", "", false
	}
	return s, fn, true
}

// parseFuncDef parses a pg_get_functiondef() CREATE FUNCTION definition and
// returns its input argument count (IN/INOUT/VARIADIC), whether it returns
// a set (SETOF or TABLE(...)), and whether any parameter has a DEFAULT value
// or is VARIADIC.
func parseFuncDef(def string) (arity int, setof, defOrVar, ok bool) {
	if def == "" {
		return 0, false, false, false
	}
	list, err := parser.Parse(def)
	if err != nil || list == nil || len(list.Items) != 1 {
		return 0, false, false, false
	}
	raw, ok := list.Items[0].(*ast.RawStmt)
	if !ok {
		return 0, false, false, false
	}
	stmt, ok := raw.Stmt.(*ast.CreateFunctionStmt)
	if !ok {
		return 0, false, false, false
	}
	if stmt.ReturnType != nil && stmt.ReturnType.Setof {
		setof = true
	}
	if stmt.Parameters != nil {
		for _, item := range stmt.Parameters.Items {
			fp, ok := item.(*ast.FunctionParameter)
			if !ok {
				return 0, false, false, false
			}
			if fp.Defexpr != nil || fp.Mode == ast.FUNC_PARAM_VARIADIC {
				defOrVar = true
			}
			switch fp.Mode {
			case ast.FUNC_PARAM_OUT:
			case ast.FUNC_PARAM_TABLE:
				setof = true
			default:
				arity++
			}
		}
	}
	return arity, setof, defOrVar, true
}

// builtinCatalogFunctions holds every built-in function name in PostgreSQL's
// pg_catalog (from pg_proc.dat via pg/catalog.BuiltinProcs). None of these
// names is ever claimed absent on an unqualified search path.
var builtinCatalogFunctions = func() map[string]bool {
	m := make(map[string]bool, len(pgcatalog.BuiltinProcs))
	for i := range pgcatalog.BuiltinProcs {
		if pgcatalog.BuiltinProcs[i].Kind != 'p' {
			m[pgcatalog.BuiltinProcs[i].Name] = true
		}
	}
	return m
}()

// Lookup implements walk.Snapshot.
func (s *pgSnapshot) Lookup(k walk.Key) (walk.Fact, string) {
	switch k.Space {
	case spaceSchema:
		switch {
		case s.schemas[k.Name]:
			return walk.Present, ""
		case trustedSchemaName(k.Name) && k.Parent == "":
			return walk.Absent, ""
		}
	case spaceRelation:
		schema := k.Parent
		if schema == tempSchema {
			// The change runs on a session of its own, which starts with
			// no temporary relations.
			return walk.Absent, ""
		}
		if kind, ok := s.relations[[2]string{schema, k.Name}]; ok {
			return walk.Present, kind
		}
		if s.schemas[schema] && !s.lost && trustedRelationName(k.Name) && !s.anywhere[k.Name] {
			return walk.Absent, ""
		}
	case spaceType:
		if k.Parent == tempSchema {
			// The session starts with no temporary types.
			return walk.Absent, ""
		}
		if s.types[[2]string{k.Parent, k.Name}] {
			// The sync lists no domain.
			return walk.Present, typOther
		}
	case spaceFunction:
		schema := k.Parent
		if schema == tempSchema {
			return walk.Absent, ""
		}
		if s.functions[[2]string{schema, k.Name}] {
			return walk.Present, ""
		}
		if s.ambigRoutines[[2]string{schema, k.Name}] || s.rangeConstructors[[2]string{schema, k.Name}] {
			return walk.Unknown, ""
		}
		if s.schemas[schema] && trustedRoutineName(k.Name) {
			return walk.Absent, ""
		}
	case spaceProcedure:
		schema := k.Parent
		if schema == tempSchema {
			return walk.Absent, ""
		}
		if s.procedures[[2]string{schema, k.Name}] {
			return walk.Present, ""
		}
		if s.ambigRoutines[[2]string{schema, k.Name}] {
			return walk.Unknown, ""
		}
		if s.schemas[schema] && trustedRoutineName(k.Name) {
			return walk.Absent, ""
		}
	case spaceColumn:
		table, ok := splitTable(k.Parent)
		if !ok || systemColumns[k.Name] {
			// A system column is in no column list; the lowering judges
			// the statements that name one from the text.
			break
		}
		if cols, ok := s.columns[table]; ok {
			if cols[k.Name] {
				return walk.Present, kindColumn
			}
			return walk.Absent, ""
		}
	case spaceConstraint:
		if table, ok := splitTable(k.Parent); ok {
			if kind, ok := s.constraints[table][k.Name]; ok {
				return walk.Present, kind
			}
		}
	case spacePrimaryKey:
		if table, ok := splitTable(k.Parent); ok && s.primaryKeys[table] {
			return walk.Present, ""
		}
	case spaceTrigger:
		if table, ok := splitTable(k.Parent); ok {
			if trigs, synced := s.triggers[table]; synced {
				if trigs[k.Name] {
					return walk.Present, ""
				}
				if !isInternalConstraintTrigger(k.Name) {
					return walk.Absent, ""
				}
			}
		}
	}
	return walk.Unknown, ""
}

// builtinCatalogTypes lists PostgreSQL's pinned built-in types in pg_catalog
// (pg: src/include/catalog/pg_type.dat). None can be dropped, none is a domain,
// none is an enum, and none is a composite type.
var builtinCatalogTypes = map[string]bool{
	"bool": true, "bytea": true, "char": true, "name": true,
	"int8": true, "int2": true, "int2vector": true, "int4": true,
	"regproc": true, "text": true, "oid": true, "tid": true, "xid": true, "cid": true,
	"oidvector": true, "json": true, "xml": true, "pg_node_tree": true,
	"point": true, "lseg": true, "path": true, "box": true, "polygon": true, "line": true,
	"float4": true, "float8": true, "unknown": true, "circle": true, "money": true,
	"macaddr": true, "inet": true, "cidr": true, "macaddr8": true,
	"aclitem": true, "bpchar": true, "varchar": true, "date": true, "time": true,
	"timestamp": true, "timestamptz": true, "interval": true, "timetz": true,
	"bit": true, "varbit": true, "numeric": true, "refcursor": true,
	"regprocedure": true, "regoper": true, "regoperator": true, "regclass": true,
	"regcollation": true, "regtype": true, "regrole": true, "regnamespace": true,
	"uuid": true, "pg_lsn": true, "tsvector": true, "gtsvector": true, "tsquery": true,
	"regconfig": true, "regdictionary": true, "jsonb": true, "jsonpath": true,
	"txid_snapshot": true, "pg_snapshot": true, "xid8": true,
	"int4range": true, "numrange": true, "tsrange": true, "tstzrange": true, "daterange": true, "int8range": true,
	"int4multirange": true, "nummultirange": true, "tsmultirange": true, "tstzmultirange": true, "datemultirange": true, "int8multirange": true,
	"record": true, "cstring": true, "any": true, "anyarray": true, "void": true,
	"trigger": true, "event_trigger": true, "language_handler": true, "internal": true,
	"anyelement": true, "anynonarray": true, "anyenum": true, "fdw_handler": true,
	"index_am_handler": true, "tsm_handler": true, "table_am_handler": true, "anyrange": true,
	"anycompatible": true, "anycompatiblearray": true, "anycompatiblenonarray": true,
	"anycompatiblerange": true, "anymultirange": true, "anycompatiblemultirange": true,
	"pg_brin_bloom_summary": true, "pg_brin_minmax_multi_summary": true,
	"pg_ndistinct": true, "pg_dependencies": true, "pg_mcv_list": true, "pg_ddl_command": true,
}

// splitTable reads the (schema, table) a column or constraint parent
// names.
func splitTable(parent string) ([2]string, bool) {
	names := walk.Names(parent)
	if len(names) != 2 {
		return [2]string{}, false
	}
	return [2]string{names[0], names[1]}, true
}
