package review

import (
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review/walk"
)

// The walk's spaces (docs/specs/2026-10-22-walk-through-mysql.md §2). A
// database is a key of the root; a relation, temporary table, foreign key,
// trigger, routine, or event a key of its database's path; a column, an
// index, or a table's foreign key a key of its table's path. A foreign key
// has two keys: its name is unique in the database (1826), and DROP
// FOREIGN KEY finds it only on its own table (1091).
const (
	spaceDatabase        = "database"
	spaceRelation        = "relation"
	spaceTemporary       = "temporary"
	spaceColumn          = "column"
	spaceIndex           = "index"
	spaceForeignKey      = "foreign key"
	spaceTableForeignKey = "table foreign key"
	spaceTrigger         = "trigger"
	spaceProcedure       = "procedure"
	spaceFunction        = "function"
	spaceEvent           = "event"
)

// gipkColumn is the column of an invisible primary key
// (sql_generate_invisible_primary_key), which the sync may not list.
const gipkColumn = "my_row_id"

// Relation kinds: tables and views share one space.
const (
	relTable = "table"
	relView  = "view"
)

// primaryIndex is the key of the primary key in a table's index space.
const primaryIndex = "primary"

// nonASCII marks the key of a name with a byte outside ASCII. MySQL
// compares such names by rules the walk does not follow (namesOutsideASCII),
// and may take one for an ASCII name: the Kelvin sign is k as a column,
// and é is e as a routine. Such a name is never known to be absent, and a
// space the sync lists one in has no name known to be absent.
const nonASCII = "\x01"

// nameKey returns the key of a name in a space, false when the walk cannot
// read it. Database, relation, and temporary table names fold to lower
// case when foldTables (lower_case_table_names 1 or 2); every other name
// always folds (verified on MySQL 5.7, 8.0, and 8.4). Only ASCII letters
// fold.
func nameKey(space, name string, foldTables bool) (string, bool) {
	if name == "" {
		return "", false
	}
	for i := 0; i < len(name); i++ {
		switch c := name[i]; {
		case c == 0:
			return "", false
		case c >= 0x80:
			return nonASCII + name, true
		}
	}
	switch space {
	case spaceDatabase, spaceRelation, spaceTemporary:
		if !foldTables {
			return name, true
		}
	}
	return strings.ToLower(name), true
}

// mysqlSnapshot is a target's synced schema as WalkThrough trusts it: one
// database, the target. It answers Present for what Bytebase's sync lists
// and Absent where that sync is complete (spec §3): the target's
// relations, the columns and indexes of its tables but the primary key,
// and its foreign keys, triggers, routines, and events. Anything of
// another database is Unknown, and the session starts with no temporary
// table.
type mysqlSnapshot struct {
	db          string
	relations   map[string]string
	columns     map[string]map[string]bool
	indexes     map[string]map[string]bool
	foreignKeys map[string]bool
	// tableForeignKeys holds each table's foreign keys.
	tableForeignKeys map[string]map[string]bool
	triggers         map[string]string // key → the spelling the sync lists
	procedures       map[string]bool
	functions        map[string]bool
	events           map[string]bool
	// triggered lists the relation keys of the tables with a trigger.
	triggered map[string]bool
	// ignoresFKs lists those of the tables whose engine ignores foreign
	// keys.
	ignoresFKs map[string]bool
	// mixed lists the spaces where the sync lists a name outside ASCII, as
	// the space and its table's key ("" for a space of the database).
	mixed map[[2]string]bool
}

func newMySQLSnapshot(db *metadata.DatabaseSchemaMetadata, foldTables bool) *mysqlSnapshot {
	s := &mysqlSnapshot{
		relations:        make(map[string]string),
		columns:          make(map[string]map[string]bool),
		indexes:          make(map[string]map[string]bool),
		foreignKeys:      make(map[string]bool),
		triggers:         make(map[string]string),
		tableForeignKeys: make(map[string]map[string]bool),
		procedures:       make(map[string]bool),
		functions:        make(map[string]bool),
		events:           make(map[string]bool),
		triggered:        make(map[string]bool),
		ignoresFKs:       make(map[string]bool),
		mixed:            make(map[[2]string]bool),
	}
	s.db, _ = nameKey(spaceDatabase, db.GetName(), foldTables)
	// mark notes a key outside ASCII in a space of the database (table
	// "") or of a table.
	mark := func(space, table, k string) {
		if strings.HasPrefix(k, nonASCII) {
			s.mixed[[2]string{space, table}] = true
		}
	}
	set := func(m map[string]bool, space, table, name string) {
		if k, ok := nameKey(space, name, foldTables); ok {
			m[k] = true
			mark(space, table, k)
		}
	}
	for _, schema := range db.GetSchemas() {
		for _, t := range schema.GetTables() {
			tk, ok := nameKey(spaceRelation, t.GetName(), foldTables)
			if !ok {
				continue
			}
			s.relations[tk] = relTable
			mark(spaceRelation, "", tk)
			if !keepsFKs(t.GetEngine()) {
				s.ignoresFKs[tk] = true
			}
			cols := make(map[string]bool)
			for _, c := range t.GetColumns() {
				set(cols, spaceColumn, tk, c.GetName())
			}
			s.columns[tk] = cols
			idx := make(map[string]bool)
			for _, i := range t.GetIndexes() {
				if i.GetPrimary() {
					idx[primaryIndex] = true
					continue
				}
				set(idx, spaceIndex, tk, i.GetName())
			}
			s.indexes[tk] = idx
			fks := make(map[string]bool)
			for _, fk := range t.GetForeignKeys() {
				set(s.foreignKeys, spaceForeignKey, "", fk.GetName())
				set(fks, spaceTableForeignKey, tk, fk.GetName())
			}
			s.tableForeignKeys[tk] = fks
			for _, tr := range t.GetTriggers() {
				if k, ok := nameKey(spaceTrigger, tr.GetName(), foldTables); ok {
					// 5.7 keeps case variants apart: a key two triggers
					// share has no one spelling.
					if spelling, seen := s.triggers[k]; seen && spelling != tr.GetName() {
						s.triggers[k] = ""
					} else {
						s.triggers[k] = tr.GetName()
					}
					s.triggered[tk] = true
					mark(spaceTrigger, "", k)
				}
			}
		}
		for _, v := range schema.GetViews() {
			if k, ok := nameKey(spaceRelation, v.GetName(), foldTables); ok {
				s.relations[k] = relView
				mark(spaceRelation, "", k)
			}
		}
		for _, f := range schema.GetFunctions() {
			set(s.functions, spaceFunction, "", f.GetName())
		}
		for _, p := range schema.GetProcedures() {
			set(s.procedures, spaceProcedure, "", p.GetName())
		}
		for _, e := range schema.GetEvents() {
			set(s.events, spaceEvent, "", e.GetName())
		}
	}
	return s
}

// Lookup implements walk.Snapshot.
func (s *mysqlSnapshot) Lookup(k walk.Key) (walk.Fact, string) {
	switch k.Space {
	case spaceDatabase:
		if k.Name == s.db {
			return walk.Present, ""
		}
		return walk.Unknown, ""
	case spaceTemporary:
		// The change runs on a connection of its own.
		return walk.Absent, ""
	}
	path := walk.Names(k.Parent)
	if len(path) == 0 || path[0] != s.db {
		return walk.Unknown, ""
	}
	// absent is the answer for a name the sync does not list: Absent,
	// unless the name, or one the sync lists in its space, is one whose
	// comparison the walk does not follow.
	table := ""
	if len(path) == 2 {
		table = path[1]
	}
	absent := func() (walk.Fact, string) {
		if strings.HasPrefix(k.Name, nonASCII) || s.mixed[[2]string{k.Space, table}] {
			return walk.Unknown, ""
		}
		return walk.Absent, ""
	}
	in := func(m map[string]bool) (walk.Fact, string) {
		if m[k.Name] {
			return walk.Present, ""
		}
		return absent()
	}
	switch k.Space {
	case spaceRelation:
		if kind, ok := s.relations[k.Name]; ok {
			return walk.Present, kind
		}
		return absent()
	case spaceColumn, spaceIndex, spaceTableForeignKey:
		if len(path) != 2 {
			return walk.Unknown, ""
		}
		if (k.Name == primaryIndex && k.Space == spaceIndex || k.Name == gipkColumn && k.Space == spaceColumn) && !s.indexes[path[1]][primaryIndex] {
			// An invisible primary key and its column, which the sync
			// does not list under show_gipk_in_create_table_and_information_schema=OFF
			// (verified on 8.0 and 8.4).
			return walk.Unknown, ""
		}
		m, ok := map[string]map[string]map[string]bool{
			spaceColumn:          s.columns,
			spaceIndex:           s.indexes,
			spaceTableForeignKey: s.tableForeignKeys,
		}[k.Space][path[1]]
		if !ok {
			return walk.Unknown, ""
		}
		return in(m)
	case spaceForeignKey:
		return in(s.foreignKeys)
	case spaceTrigger:
		if spelling, ok := s.triggers[k.Name]; ok {
			return walk.Present, spelling
		}
		return absent()
	case spaceProcedure:
		return in(s.procedures)
	case spaceFunction:
		return in(s.functions)
	case spaceEvent:
		return in(s.events)
	}
	return walk.Unknown, ""
}
