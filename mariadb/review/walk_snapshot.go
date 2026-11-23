package review

import (
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review/walk"
)

// The walk's spaces (docs/specs/2026-10-22-walk-through-mysql.md §2, as
// docs/specs/2026-10-26-walk-through-mariadb-tidb.md §2 changes it). A
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

// Relation kinds: tables, views, and sequences share one space (1050,
// verified). MariaDB has no invisible primary key, so the walk keeps no
// column of one.
const (
	relTable    = "table"
	relView     = "view"
	relSequence = "sequence"
)

// primaryIndex is the key of the primary key in a table's index space.
const primaryIndex = "primary"

// nonASCII marks the key of a name with a byte outside ASCII. MariaDB
// compares such names by rules the walk does not follow (namesOutsideASCII),
// and may take one for an ASCII name: ſ is s as a column, and é is e as
// a routine. Such a name is never known to be absent, and a space the
// sync lists one in has no name known to be absent.
const nonASCII = "\x01"

// nameKey returns the key of a name in a space, false when the walk cannot
// read it. Database, relation, and temporary table names fold to lower
// case when foldTables (lower_case_table_names 1 or 2); every other name
// folds (verified on MySQL 5.7, 8.0, and 8.4), triggers too, whose
// spelling the walk keeps beside the key (MariaDB compares them with case,
// as MySQL 5.7 does). Only ASCII letters fold.
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

// mariadbSnapshot is a target's synced schema as WalkThrough trusts it:
// one database, the target. It answers Present for what Bytebase's sync
// lists and Absent where that sync is complete (MariaDB spec §2.1): the
// target's relations, which it lists with its sequences and
// system-versioned tables since Bytebase #21665, and its triggers; the
// columns, the indexes, the primary key among them, and the foreign keys
// of the tables it lists; the target's foreign key names; and its
// routines and events. The columns of a system-versioned table, which the
// sync marks with WITH SYSTEM VERSIONING in its create options, are never
// Absent:
// information_schema does not list its hidden row_start and row_end
// (verified; see walker.versioned). Anything of another database is
// Unknown, and the session starts with no temporary table.
type mariadbSnapshot struct {
	db        string
	relations map[string]string
	// versioned lists the relation keys of the system-versioned tables.
	versioned map[string]bool
	// partitioned lists the relation keys of the partitioned tables, and
	// rangeOrList those partitioned by RANGE or LIST.
	partitioned map[string]bool
	rangeOrList map[string]bool
	columns     map[string]map[string]bool
	// colTypes holds each listed column's base type (baseType), the kind
	// Lookup gives its key.
	colTypes map[string]map[string]string
	// columnNames lists the keys of each listed table's columns.
	columnNames map[string][]string
	// sequences holds the values the sync lists of each sequence it
	// read them for.
	sequences map[string]sequenceValues
	// fkTables holds the foreign keys and indexes of the InnoDB tables
	// (syncedFKTable).
	fkTables map[string]*fkTable
	// tables holds what the sync says of each base table beyond its
	// names (syncedTable).
	tables      map[string]*syncedTable
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
	// fkGeneratedIdx lists synced non-unique secondary indexes that may be
	// InnoDB auto-generated foreign-key backing indexes (HA_GENERATED_KEY),
	// mapped by table key -> index key -> lowercased index columns (nil when
	// column metadata was not provided). MariaDB silently drops such an index
	// when a newly added explicit index covers its columns, even under the same
	// index name.
	fkGeneratedIdx map[string]map[string][]string
}

func newMariaDBSnapshot(db *metadata.DatabaseSchemaMetadata, foldTables bool) *mariadbSnapshot {
	s := &mariadbSnapshot{
		relations:        make(map[string]string),
		versioned:        make(map[string]bool),
		partitioned:      make(map[string]bool),
		rangeOrList:      make(map[string]bool),
		columns:          make(map[string]map[string]bool),
		colTypes:         make(map[string]map[string]string),
		columnNames:      make(map[string][]string),
		sequences:        make(map[string]sequenceValues),
		fkTables:         make(map[string]*fkTable),
		tables:           make(map[string]*syncedTable),
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
		fkGeneratedIdx:   make(map[string]map[string][]string),
	}
	s.db, _ = nameKey(spaceDatabase, db.GetName(), foldTables)
	// mark notes a key outside ASCII in a space of a table, or of the
	// database for table "".
	mark := func(space, table, k string) {
		if strings.HasPrefix(k, nonASCII) {
			s.mixed[[2]string{space, table}] = true
		}
	}
	// set lists a name in a space of a table, or of the database.
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
			if strings.Contains(strings.ToUpper(t.GetCreateOptions()), "WITH SYSTEM VERSIONING") {
				s.versioned[tk] = true
			}
			if parts := t.GetPartitions(); len(parts) > 0 {
				s.partitioned[tk] = true
				switch parts[0].GetType() {
				case metadata.TablePartitionMetadata_RANGE, metadata.TablePartitionMetadata_RANGE_COLUMNS,
					metadata.TablePartitionMetadata_LIST, metadata.TablePartitionMetadata_LIST_COLUMNS:
					s.rangeOrList[tk] = true
				}
			}
			if !keepsFKs(t.GetEngine()) {
				s.ignoresFKs[tk] = true
			}
			if ft, ok := syncedFKTable(t, foldTables); ok {
				s.fkTables[tk] = ft
			}
			s.tables[tk] = newSyncedTable(s.db, t, foldTables)
			cols := make(map[string]bool)
			types := make(map[string]string)
			for _, c := range t.GetColumns() {
				set(cols, spaceColumn, tk, c.GetName())
				if k, ok := nameKey(spaceColumn, c.GetName(), foldTables); ok {
					types[k] = baseType(c.GetType())
					s.columnNames[tk] = append(s.columnNames[tk], k)
				}
			}
			s.columns[tk] = cols
			s.colTypes[tk] = types
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
			if len(t.GetForeignKeys()) > 0 {
				for _, i := range t.GetIndexes() {
					if cols, ok := mayBeFKGeneratedIndex(i, t.GetForeignKeys()); ok {
						if ik, valid := nameKey(spaceIndex, i.GetName(), foldTables); valid {
							if s.fkGeneratedIdx[tk] == nil {
								s.fkGeneratedIdx[tk] = make(map[string][]string)
							}
							s.fkGeneratedIdx[tk][ik] = cols
						}
					}
				}
			}
			for _, tr := range t.GetTriggers() {
				if k, ok := nameKey(spaceTrigger, tr.GetName(), foldTables); ok {
					// MariaDB keeps case variants apart: a key two
					// triggers share has no one spelling.
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
		for _, q := range schema.GetSequences() {
			if k, ok := nameKey(spaceRelation, q.GetName(), foldTables); ok {
				s.relations[k] = relSequence
				mark(spaceRelation, "", k)
				if values, ok := syncedSequenceValues(q); ok {
					s.sequences[k] = values
				}
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
func (s *mariadbSnapshot) Lookup(k walk.Key) (walk.Fact, string) {
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
		if len(path) != 2 || k.Space == spaceColumn && s.versioned[path[1]] {
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
		if k.Space == spaceColumn && m[k.Name] {
			return walk.Present, s.colTypes[path[1]][k.Name]
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

// mayBeFKGeneratedIndex reports whether a synced index could be an InnoDB
// auto-generated foreign-key backing index (HA_GENERATED_KEY) for one of fks,
// and returns its lowercased column list (nil when column metadata is absent).
func mayBeFKGeneratedIndex(idx *metadata.IndexMetadata, fks []*metadata.ForeignKeyMetadata) ([]string, bool) {
	if idx == nil || idx.GetPrimary() || idx.GetUnique() {
		return nil, false
	}
	switch strings.ToUpper(idx.GetType()) {
	case "FULLTEXT", "SPATIAL", "VECTOR", "HASH", "RTREE":
		return nil, false
	}
	for _, l := range idx.GetKeyLength() {
		if l > 0 {
			return nil, false
		}
	}
	idxName := strings.ToLower(idx.GetName())
	exprs := idx.GetExpressions()
	for _, fk := range fks {
		if fk == nil {
			continue
		}
		fkCols := fk.GetColumns()
		if len(exprs) > 0 && len(fkCols) > 0 {
			if len(exprs) != len(fkCols) {
				continue
			}
			match := true
			lowerCols := make([]string, len(fkCols))
			for j := range fkCols {
				if !strings.EqualFold(exprs[j], fkCols[j]) {
					match = false
					break
				}
				lowerCols[j] = strings.ToLower(fkCols[j])
			}
			if !match {
				continue
			}
			if isFKGeneratedIndexName(idxName, fk.GetName(), lowerCols[0]) {
				return lowerCols, true
			}
		} else if fk.GetName() != "" && idxName == strings.ToLower(fk.GetName()) {
			return nil, true
		}
	}
	return nil, false
}

func isFKGeneratedIndexName(idxLower, fkName, firstColLower string) bool {
	if fkName != "" && idxLower == strings.ToLower(fkName) {
		return true
	}
	if firstColLower == "" {
		return false
	}
	if idxLower == firstColLower {
		return true
	}
	if rest, ok := strings.CutPrefix(idxLower, firstColLower+"_"); ok && rest != "" {
		for i := 0; i < len(rest); i++ {
			if rest[i] < '0' || rest[i] > '9' {
				return false
			}
		}
		return true
	}
	return false
}
