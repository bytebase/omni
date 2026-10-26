package review

import (
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review/walk"
)

// The walk's spaces (docs/specs/2026-10-22-walk-through-mysql.md §2, with
// the TiDB deltas of docs/specs/2026-10-26-walk-through-mariadb-tidb.md
// §3). A database is a key of the root; a relation or temporary table a
// key of its database's path; a column, an index, or a table's foreign key
// a key of its table's path. TiDB has no triggers, functions, or events
// (1064), and stored procedures fail (8108), so the walk keeps no space
// for them. Foreign key names are unique per table, not per database
// (verified: a name another table's foreign key has is accepted, 1826 on
// the same table), so a foreign key has the one key on its table.
const (
	spaceDatabase        = "database"
	spaceRelation        = "relation"
	spaceTemporary       = "temporary"
	spaceColumn          = "column"
	spaceIndex           = "index"
	spaceTableForeignKey = "table foreign key"
)

// rowIDColumn is the hidden row id of a table without a clustered integer
// key, which information_schema does not list. ADD COLUMN of the name
// fails (1166) whatever the table, so the walk claims nothing of it.
const rowIDColumn = "_tidb_rowid"

// Relation kinds: tables, views, and sequences share one space (1050 for
// a table, view, or sequence of a sequence's name, verified).
const (
	relTable    = "table"
	relView     = "view"
	relSequence = "sequence"
)

// primaryIndex is the key of the primary key in a table's index space.
const primaryIndex = "primary"

// nonASCII marks the key of a name with a byte outside ASCII. TiDB's
// comparison of such names depends on collations the walk does not
// follow: two spellings are one name only when they are identical, and
// such a name is never known to be absent.
const nonASCII = "\x01"

// nameKey returns the key of a name in a space, false when the walk cannot
// read it. Database, relation, and temporary table names fold to lower
// case when foldTables; every other name always folds. TiDB's
// lower_case_table_names is 2 and fixed, so the walk always folds
// (verified: CREATE TABLE T beside t, CREATE DATABASE SHOP beside shop,
// and CREATE TEMPORARY TABLE TT beside tt fail). Only ASCII letters fold.
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

// tidbSnapshot is a target's synced schema as WalkThrough trusts it: one
// database, the target. It answers Present for what Bytebase's sync lists
// and Absent where that sync is complete: the target's relations (base
// tables, views, and sequences, which Bytebase's TiDB sync reads from
// information_schema.SEQUENCES), the columns and indexes of its tables,
// the primary key included, and their foreign keys. TiDB has no invisible
// primary key (sql_generate_invisible_primary_key is unknown, 1193), so a
// table the sync lists without a primary key has none. Anything of
// another database is Unknown, and the session starts with no temporary
// table.
type tidbSnapshot struct {
	db        string
	relations map[string]string
	columns   map[string]map[string]bool
	indexes   map[string]map[string]bool
	// tableForeignKeys holds each table's foreign keys.
	tableForeignKeys map[string]map[string]bool
}

func newTiDBSnapshot(db *metadata.DatabaseSchemaMetadata, foldTables bool) *tidbSnapshot {
	s := &tidbSnapshot{
		relations:        make(map[string]string),
		columns:          make(map[string]map[string]bool),
		indexes:          make(map[string]map[string]bool),
		tableForeignKeys: make(map[string]map[string]bool),
	}
	s.db, _ = nameKey(spaceDatabase, db.GetName(), foldTables)
	set := func(m map[string]bool, space, name string) {
		if k, ok := nameKey(space, name, foldTables); ok {
			m[k] = true
		}
	}
	for _, schema := range db.GetSchemas() {
		for _, t := range schema.GetTables() {
			tk, ok := nameKey(spaceRelation, t.GetName(), foldTables)
			if !ok {
				continue
			}
			s.relations[tk] = relTable
			cols := make(map[string]bool)
			for _, c := range t.GetColumns() {
				set(cols, spaceColumn, c.GetName())
			}
			s.columns[tk] = cols
			idx := make(map[string]bool)
			for _, i := range t.GetIndexes() {
				if i.GetPrimary() {
					idx[primaryIndex] = true
					continue
				}
				set(idx, spaceIndex, i.GetName())
			}
			s.indexes[tk] = idx
			fks := make(map[string]bool)
			for _, fk := range t.GetForeignKeys() {
				set(fks, spaceTableForeignKey, fk.GetName())
			}
			s.tableForeignKeys[tk] = fks
		}
		for _, v := range schema.GetViews() {
			if k, ok := nameKey(spaceRelation, v.GetName(), foldTables); ok {
				s.relations[k] = relView
			}
		}
		for _, q := range schema.GetSequences() {
			if k, ok := nameKey(spaceRelation, q.GetName(), foldTables); ok {
				s.relations[k] = relSequence
			}
		}
	}
	return s
}

// Lookup implements walk.Snapshot.
func (s *tidbSnapshot) Lookup(k walk.Key) (walk.Fact, string) {
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
	// unless the name is one whose comparison the walk does not follow.
	absent := func() (walk.Fact, string) {
		if strings.HasPrefix(k.Name, nonASCII) {
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
		if k.Space == spaceColumn && k.Name == rowIDColumn {
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
	}
	return walk.Unknown, ""
}
