package review

import (
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
	"github.com/bytebase/omni/review/walk"
)

// The spaces WalkThrough keeps facts in. A relation and a type are held by
// a schema, a column, a constraint, and the primary key by a table; the
// primary key space has one key per table, named "".
const (
	spaceSchema     = "schema"
	spaceRelation   = "relation"
	spaceType       = "type"
	spaceColumn     = "column"
	spaceConstraint = "constraint"
	spacePrimaryKey = "primary key"
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
	// anywhere lists the partition and partition index names, Unknown in
	// every schema.
	anywhere map[string]bool
	// lost is true when the sync may have lost partitions, of any name in
	// any schema: no relation's absence is trusted then.
	lost bool
	// inSchema lists the relations and types each schema holds.
	inSchema map[string][]walk.Key
}

// newPgSnapshot reads a target's synced schema. A nil schema gives a
// snapshot that knows nothing.
func newPgSnapshot(db *metadata.DatabaseSchemaMetadata) *pgSnapshot {
	s := &pgSnapshot{
		schemas:     make(map[string]bool),
		relations:   make(map[[2]string]string),
		types:       make(map[[2]string]bool),
		columns:     make(map[[2]string]map[string]bool),
		constraints: make(map[[2]string]map[string]string),
		primaryKeys: make(map[[2]string]bool),
		anywhere:    make(map[string]bool),
		inSchema:    make(map[string][]walk.Key),
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
			for _, c := range t.GetColumns() {
				cols[c.GetName()] = true
			}
			s.columns[key] = cols
			cons := make(map[string]string)
			for _, idx := range t.GetIndexes() {
				rel(idx.GetName(), relIndex)
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
			}
			for _, c := range t.GetCheckConstraints() {
				cons[c.GetName()] = conCheck
			}
			for _, x := range t.GetExcludeConstraints() {
				cons[x.GetName()] = conExclusion
			}
			s.constraints[key] = cons
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
		}
		for _, v := range schema.GetMaterializedViews() {
			rel(v.GetName(), relMatview)
			rowType(v.GetName())
			for _, idx := range v.GetIndexes() {
				rel(idx.GetName(), relIndex)
			}
		}
		for _, q := range schema.GetSequences() {
			rel(q.GetName(), relSequence)
		}
		for _, c := range schema.GetCompositeTypes() {
			rel(c.GetName(), relComposite)
			rowType(c.GetName())
		}
		for _, e := range schema.GetEnumTypes() {
			rowType(e.GetName())
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
		if s.types[[2]string{k.Parent, k.Name}] {
			// The sync lists no domain.
			return walk.Present, typOther
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
	}
	return walk.Unknown, ""
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
