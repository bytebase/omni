package review

import (
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review/walk"
)

// The walk's spaces (docs/specs/2026-10-26-walk-through-oracle.md §2). The
// object space is a schema's object namespace, which tables, views,
// materialized views, sequences, private synonyms, routines, packages, and
// types share; its key's kind is what the object is. The table, view, and
// sequence spaces repeat what the object space holds for those kinds, so
// that a statement of one kind can trust the sync where it is complete for
// that kind. A column, the primary key, and a constraint as DROP
// CONSTRAINT finds it on its table are keys of the table's path; a
// constraint name, an index, and a trigger are keys of the schema, each in
// a namespace of its own, of the kind of the table path they are on.
const (
	spaceObject          = "object"
	spaceTable           = "table"
	spaceView            = "view"
	spaceSequence        = "sequence"
	spaceColumn          = "column"
	spacePrimary         = "primary key"
	spaceTableConstraint = "table constraint"
	spaceConstraint      = "constraint"
	spaceIndex           = "index"
	spaceTrigger         = "trigger"
)

// Object kinds, as Oracle names them.
const (
	kindTable     = "TABLE"
	kindView      = "VIEW"
	kindMView     = "MATERIALIZED VIEW"
	kindSequence  = "SEQUENCE"
	kindSynonym   = "SYNONYM"
	kindProcedure = "PROCEDURE"
	kindFunction  = "FUNCTION"
	kindPackage   = "PACKAGE"
	kindType      = "TYPE"
)

// generated reports whether a name is one the walk never claims in a
// space: a name Oracle generates, SYS_ for constraints, indexes, LOB
// segments, hidden columns, and IOT tables, ISEQ$$_ for identity
// sequences, BIN$ for the recycle bin, and ORA$ for private temporary
// tables; and, among columns, the hidden columns Oracle adds under names
// of its own, ORA_ARCHIVE_STATE for row archival and ORABCTAB_..$ for
// blockchain and immutable tables (verified: DDL on them fails with
// 38398 and 05715), which the sync does not list.
func generated(space, name string) bool {
	for _, p := range []string{"SYS_", "ISEQ$$_", "BIN$", "ORA$"} {
		if strings.HasPrefix(name, p) {
			return true
		}
	}
	return space == spaceColumn && (strings.HasPrefix(name, "ORA_") || strings.HasPrefix(name, "ORABCTAB_"))
}

// tablePath is the parent of a table's columns, primary key, and
// constraints, and the kind of its indexes, constraint names, and
// triggers.
func tablePath(schema, name string) string { return walk.Path(schema, name) }

// oracleSnapshot is a target's synced schema as WalkThrough trusts it: one
// schema, the database the target names (spec §3). It answers Present for
// what Bytebase's sync lists, and Absent where that sync lists every object
// of a space: tables and materialized views, views, and sequences; the
// columns of a table whose columns it lists in full; and whether a table
// has a primary key. Constraint names, indexes, and triggers it lists only
// in part, and the object namespace holds synonyms and types the sync does
// not list, so their absence is Unknown. Anything of another schema is
// Unknown.
type oracleSnapshot struct {
	schema  string
	objects map[string]string
	tables  map[string]*syncedTable
	// constraints, indexes, and triggers map a name to the table or view
	// it is on. constraintIndexes holds the names of the primary and unique
	// constraints, whose indexes may have those names, another, or none.
	constraints       map[string]string
	indexes           map[string]string
	constraintIndexes map[string]bool
	triggers          map[string]string
	// domainTables holds the tables and materialized views with a domain
	// index, and domainIndexes those indexes: an indextype's code runs on
	// the table's DDL and DML, and can create and drop tables of its own,
	// as Oracle Text's DR$ tables. storage holds the tables the sync lists
	// no column of, a nested table's storage tables: DROP TABLE of the
	// parent, and DROP COLUMN of the nested table, drop them (verified).
	domainTables  map[string]bool
	domainIndexes map[string]bool
	storage       []string
}

// syncedTable is what the snapshot holds of a table or materialized view.
type syncedTable struct {
	mview   bool
	columns map[string]bool
	// complete is set when the sync lists every column of the table.
	complete    bool
	primary     bool
	constraints map[string]bool
}

func newOracleSnapshot(db *metadata.DatabaseSchemaMetadata) *oracleSnapshot {
	s := &oracleSnapshot{
		schema:            db.GetName(),
		objects:           make(map[string]string),
		tables:            make(map[string]*syncedTable),
		constraints:       make(map[string]string),
		indexes:           make(map[string]string),
		constraintIndexes: make(map[string]bool),
		triggers:          make(map[string]string),
		domainTables:      make(map[string]bool),
		domainIndexes:     make(map[string]bool),
	}
	for _, schema := range db.GetSchemas() {
		if name := schema.GetName(); name != "" && name != s.schema {
			continue
		}
		for _, t := range schema.GetTables() {
			s.objects[t.GetName()] = kindTable
			st := &syncedTable{
				columns:     make(map[string]bool),
				constraints: make(map[string]bool),
			}
			for _, c := range t.GetColumns() {
				st.columns[c.GetName()] = true
			}
			st.complete = columnsComplete(t, st.columns)
			if len(st.columns) == 0 {
				s.storage = append(s.storage, t.GetName())
			}
			s.indexesOf(t.GetName(), st, t.GetIndexes())
			for _, fk := range t.GetForeignKeys() {
				st.constraints[fk.GetName()] = true
				s.constraints[fk.GetName()] = t.GetName()
			}
			for _, ck := range t.GetCheckConstraints() {
				st.constraints[ck.GetName()] = true
				s.constraints[ck.GetName()] = t.GetName()
			}
			s.triggersOf(t.GetName(), t.GetTriggers())
			s.tables[t.GetName()] = st
		}
		for _, v := range schema.GetViews() {
			s.objects[v.GetName()] = kindView
			s.triggersOf(v.GetName(), v.GetTriggers())
		}
		for _, mv := range schema.GetMaterializedViews() {
			s.objects[mv.GetName()] = kindMView
			st := &syncedTable{mview: true, constraints: make(map[string]bool)}
			s.indexesOf(mv.GetName(), st, mv.GetIndexes())
			s.triggersOf(mv.GetName(), mv.GetTriggers())
			s.tables[mv.GetName()] = st
		}
		for _, seq := range schema.GetSequences() {
			s.objects[seq.GetName()] = kindSequence
		}
		for _, f := range schema.GetFunctions() {
			s.objects[f.GetName()] = kindFunction
		}
		for _, p := range schema.GetProcedures() {
			s.objects[p.GetName()] = kindProcedure
		}
		for _, p := range schema.GetPackages() {
			s.objects[p.GetName()] = kindPackage
		}
	}
	return s
}

// indexesOf records a table's indexes. The sync lists a primary key or
// unique constraint as an index named after the constraint, and the
// indexes that back no constraint of their name apart (getConstraints,
// getIndexesAndConstraints).
func (s *oracleSnapshot) indexesOf(table string, st *syncedTable, indexes []*metadata.IndexMetadata) {
	for _, idx := range indexes {
		name := idx.GetName()
		if idx.GetType() == "DOMAIN" {
			s.domainTables[table] = true
			s.domainIndexes[name] = true
		}
		if !idx.GetIsConstraint() {
			s.indexes[name] = table
			continue
		}
		s.constraintIndexes[name] = true
		st.constraints[name] = true
		s.constraints[name] = table
		if idx.GetPrimary() {
			st.primary = true
		}
	}
}

func (s *oracleSnapshot) triggersOf(table string, triggers []*metadata.TriggerMetadata) {
	for _, tr := range triggers {
		s.triggers[tr.GetName()] = table
	}
}

// columnsComplete reports whether the sync lists every column of a table
// whose columns are listed: the table is no object or XMLType table, which
// keep hidden columns of user names (XMLDATA, verified); it has a column,
// which the storage table of a nested table does not show (verified); and
// no check constraint names a column the sync does not list. A valid-time
// period adds hidden columns named after it, which the sync does not list,
// and a check constraint on them, which it does (verified: PERIOD FOR vt
// adds VT, VT_START, and VT_END, and a check "VT_START" < "VT_END").
func columnsComplete(t *metadata.TableMetadata, columns map[string]bool) bool {
	if t.GetSkipDump() || len(columns) == 0 {
		return false
	}
	for _, ck := range t.GetCheckConstraints() {
		for _, name := range quotedNames(ck.GetExpression()) {
			if !columns[name] {
				return false
			}
		}
	}
	return true
}

// quotedNames returns the double-quoted identifiers of a condition as
// Oracle stores it, outside its string literals.
func quotedNames(cond string) []string {
	var out []string
	for i := 0; i < len(cond); i++ {
		switch cond[i] {
		case '\'':
			for i++; i < len(cond); i++ {
				if cond[i] == '\'' {
					if i+1 < len(cond) && cond[i+1] == '\'' {
						i++
						continue
					}
					break
				}
			}
		case '"':
			end := strings.IndexByte(cond[i+1:], '"')
			if end < 0 {
				return out
			}
			out = append(out, cond[i+1:i+1+end])
			i += end + 1
		}
	}
	return out
}

// Lookup implements walk.Snapshot.
func (s *oracleSnapshot) Lookup(k walk.Key) (walk.Fact, string) {
	path := walk.Names(k.Parent)
	if len(path) == 0 || path[0] != s.schema {
		return walk.Unknown, ""
	}
	switch k.Space {
	case spaceObject:
		if kind, ok := s.objects[k.Name]; ok {
			return walk.Present, kind
		}
		return walk.Unknown, ""
	case spaceTable, spaceView, spaceSequence:
		kind := s.objects[k.Name]
		switch {
		case k.Space == spaceTable && (kind == kindTable || kind == kindMView),
			k.Space == spaceView && kind == kindView,
			k.Space == spaceSequence && kind == kindSequence:
			return walk.Present, kind
		case generated(k.Space, k.Name):
			return walk.Unknown, ""
		}
		return walk.Absent, ""
	case spaceColumn, spacePrimary, spaceTableConstraint:
		if len(path) != 2 {
			return walk.Unknown, ""
		}
		t, ok := s.tables[path[1]]
		if !ok || t.mview {
			return walk.Unknown, ""
		}
		switch k.Space {
		case spaceColumn:
			switch {
			case t.columns[k.Name]:
				return walk.Present, ""
			case t.complete && !generated(spaceColumn, k.Name):
				return walk.Absent, ""
			}
		case spacePrimary:
			if t.primary {
				return walk.Present, ""
			}
			return walk.Absent, ""
		case spaceTableConstraint:
			if t.constraints[k.Name] {
				return walk.Present, ""
			}
		}
		return walk.Unknown, ""
	case spaceConstraint:
		if table, ok := s.constraints[k.Name]; ok {
			return walk.Present, tablePath(s.schema, table)
		}
	case spaceIndex:
		if table, ok := s.indexes[k.Name]; ok && !s.constraintIndexes[k.Name] {
			return walk.Present, tablePath(s.schema, table)
		}
	case spaceTrigger:
		if table, ok := s.triggers[k.Name]; ok {
			return walk.Present, tablePath(s.schema, table)
		}
	}
	return walk.Unknown, ""
}
