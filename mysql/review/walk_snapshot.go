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
	spaceCheckConstraint = "check constraint"
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
	// checkConstraints holds each table's CHECK constraints.
	checkConstraints map[string]map[string]bool
	triggers         map[string]string // key → the spelling the sync lists
	procedures       map[string]bool
	functions        map[string]bool
	events           map[string]bool
	// triggered lists the relation keys of the tables with a trigger.
	triggered map[string]bool
	// ignoresFKs lists those of the tables whose engine ignores foreign
	// keys.
	ignoresFKs map[string]bool
	// engines records each synced table's engine ("" means default/InnoDB).
	engines map[string]string
	// colTypes records known base data types for synced columns (omitted when
	// sync metadata left ColumnMetadata.Type empty).
	colTypes map[string]map[string]string
	// partitionState records partitioning facts for each synced table.
	partitionState map[string]*tablePartitionInfo
	// tableIndexes records index metadata per synced table when every synced
	// index on the table has known expressions.
	tableIndexes map[string]map[string]*indexMeta
	// tableFKs records each synced table's foreign key column lists.
	tableFKs map[string]map[string][]string
	// referencedByFK records synced tables referenced by a foreign key in this
	// database, mapped to the referencing (childTableKey, fkKey) pairs.
	referencedByFK map[string]map[[2]string]bool
	// fkGeneratedIndexes lists synced indexes that could be implicitly
	// generated foreign-key backing indexes (MySQL's HA_GENERATED_KEY), keyed
	// by table key then index key, mapping to the keys of the synced foreign
	// keys whose column list matches the index. Because SHOW CREATE TABLE and
	// INFORMATION_SCHEMA.STATISTICS do not distinguish a generated FK backing
	// index from an explicit one, MySQL may accept CREATE INDEX or ALTER TABLE
	// ADD INDEX reusing its name (dropping the generated index first).
	fkGeneratedIndexes map[string]map[string][]string
	// partial lists the spaces with no name known to be absent, as the
	// space and its table's key ("" for a space of the database): where the
	// sync lists a name outside ASCII.
	partial map[[2]string]bool
}

func newMySQLSnapshot(db *metadata.DatabaseSchemaMetadata, foldTables bool) *mysqlSnapshot {
	s := &mysqlSnapshot{
		relations:          make(map[string]string),
		columns:            make(map[string]map[string]bool),
		indexes:            make(map[string]map[string]bool),
		foreignKeys:        make(map[string]bool),
		triggers:           make(map[string]string),
		tableForeignKeys:   make(map[string]map[string]bool),
		checkConstraints:   make(map[string]map[string]bool),
		procedures:         make(map[string]bool),
		functions:          make(map[string]bool),
		events:             make(map[string]bool),
		triggered:          make(map[string]bool),
		ignoresFKs:         make(map[string]bool),
		engines:            make(map[string]string),
		colTypes:           make(map[string]map[string]string),
		partitionState:     make(map[string]*tablePartitionInfo),
		tableIndexes:       make(map[string]map[string]*indexMeta),
		tableFKs:           make(map[string]map[string][]string),
		referencedByFK:     make(map[string]map[[2]string]bool),
		fkGeneratedIndexes: make(map[string]map[string][]string),
		partial:            make(map[[2]string]bool),
	}
	s.db, _ = nameKey(spaceDatabase, db.GetName(), foldTables)
	// mark notes a key outside ASCII in a space of the database (table
	// "") or of a table.
	mark := func(space, table, k string) {
		if strings.HasPrefix(k, nonASCII) {
			s.partial[[2]string{space, table}] = true
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
			s.engines[tk] = t.GetEngine()
			if !keepsFKs(t.GetEngine()) {
				s.ignoresFKs[tk] = true
			}
			cols := make(map[string]bool)
			cTypes := make(map[string]string)
			for _, c := range t.GetColumns() {
				set(cols, spaceColumn, tk, c.GetName())
				if ck, ok := nameKey(spaceColumn, c.GetName(), foldTables); ok {
					if bt := normalizeBaseType(c.GetType()); bt != "" {
						cTypes[ck] = bt
					}
				}
			}
			s.columns[tk] = cols
			if len(cTypes) > 0 {
				s.colTypes[tk] = cTypes
			}
			idx := make(map[string]bool)
			idxMeta := make(map[string]*indexMeta)
			allIdxHaveExprs := true
			for _, i := range t.GetIndexes() {
				ik := primaryIndex
				if !i.GetPrimary() {
					var ok bool
					if ik, ok = nameKey(spaceIndex, i.GetName(), foldTables); ok {
						idx[ik] = true
						mark(spaceIndex, tk, ik)
					} else {
						allIdxHaveExprs = false
						continue
					}
				} else {
					idx[primaryIndex] = true
				}
				exprs := i.GetExpressions()
				if len(exprs) == 0 {
					allIdxHaveExprs = false
					continue
				}
				im := &indexMeta{
					unique:   i.GetUnique() || i.GetPrimary(),
					primary:  i.GetPrimary(),
					fulltext: strings.EqualFold(i.GetType(), "FULLTEXT"),
					spatial:  isSpatialIndexType(i.GetType()),
				}
				for _, kl := range i.GetKeyLength() {
					if kl > 0 {
						im.hasPrefix = true
					}
				}
				for _, ex := range exprs {
					trimmed := strings.TrimSpace(ex)
					if strings.HasPrefix(trimmed, "(") {
						break
					}
					im.cols = append(im.cols, strings.ToLower(strings.Trim(trimmed, "`")))
				}
				idxMeta[ik] = im
			}
			s.indexes[tk] = idx
			if allIdxHaveExprs {
				s.tableIndexes[tk] = idxMeta
			}
			fks := make(map[string]bool)
			fkCols := make(map[string][]string)
			for _, fk := range t.GetForeignKeys() {
				set(s.foreignKeys, spaceForeignKey, "", fk.GetName())
				set(fks, spaceTableForeignKey, tk, fk.GetName())
				if fkKey, ok := nameKey(spaceForeignKey, fk.GetName(), foldTables); ok {
					var colsLower []string
					for _, fc := range fk.GetColumns() {
						colsLower = append(colsLower, strings.ToLower(fc))
					}
					fkCols[fkKey] = colsLower
					refSchema := fk.GetReferencedSchema()
					sameDB := refSchema == ""
					if !sameDB {
						if refDB, ok := nameKey(spaceDatabase, refSchema, foldTables); ok && refDB == s.db {
							sameDB = true
						}
					}
					if sameDB {
						if refTbl, ok := nameKey(spaceRelation, fk.GetReferencedTable(), foldTables); ok {
							if s.referencedByFK[refTbl] == nil {
								s.referencedByFK[refTbl] = make(map[[2]string]bool)
							}
							s.referencedByFK[refTbl][[2]string{tk, fkKey}] = true
						}
					}
				}
			}
			s.tableForeignKeys[tk] = fks
			if len(fkCols) > 0 {
				s.tableFKs[tk] = fkCols
			}
			chks := make(map[string]bool)
			for _, cc := range t.GetCheckConstraints() {
				set(chks, spaceCheckConstraint, tk, cc.GetName())
			}
			s.checkConstraints[tk] = chks
			parts := t.GetPartitions()
			if len(parts) == 0 {
				s.partitionState[tk] = &tablePartitionInfo{partitioned: false}
			} else {
				pinfo := &tablePartitionInfo{
					partitioned: true,
					names:       make(map[string]bool, len(parts)),
					cols:        make(map[string]bool),
				}
				if pt, ok := partitionProtoType(parts[0].GetType()); ok {
					pinfo.hasType = true
					pinfo.partType = pt
				}
				for _, p := range parts {
					if p.GetName() != "" {
						pinfo.names[strings.ToLower(p.GetName())] = true
					}
					for _, col := range extractPartitionExprCols(p.GetExpression()) {
						pinfo.cols[col] = true
					}
					for _, sp := range p.GetSubpartitions() {
						for _, col := range extractPartitionExprCols(sp.GetExpression()) {
							pinfo.cols[col] = true
						}
					}
				}
				pinfo.count = len(pinfo.names)
				if len(pinfo.cols) == 0 {
					pinfo.cols = nil
				}
				s.partitionState[tk] = pinfo
			}
			if keepsFKs(t.GetEngine()) && len(t.GetForeignKeys()) > 0 {
				var fkGen map[string][]string
				for _, i := range t.GetIndexes() {
					owners := fkGeneratedIndexOwners(i, t.GetForeignKeys(), foldTables)
					if len(owners) == 0 {
						continue
					}
					if ik, ok := nameKey(spaceIndex, i.GetName(), foldTables); ok {
						if fkGen == nil {
							fkGen = make(map[string][]string)
						}
						fkGen[ik] = owners
					}
				}
				if len(fkGen) > 0 {
					s.fkGeneratedIndexes[tk] = fkGen
				}
			}
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
		if strings.HasPrefix(k.Name, nonASCII) || s.partial[[2]string{k.Space, table}] {
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
	case spaceColumn, spaceIndex, spaceTableForeignKey, spaceCheckConstraint:
		if len(path) != 2 {
			return walk.Unknown, ""
		}
		switch {
		case (k.Name == primaryIndex && k.Space == spaceIndex || k.Name == gipkColumn && k.Space == spaceColumn) && !s.indexes[path[1]][primaryIndex]:
			// An invisible primary key and its column, which the sync
			// does not list under show_gipk_in_create_table_and_information_schema=OFF
			// (verified on 8.0 and 8.4).
			return walk.Unknown, ""
		}
		m, ok := map[string]map[string]map[string]bool{
			spaceColumn:          s.columns,
			spaceIndex:           s.indexes,
			spaceTableForeignKey: s.tableForeignKeys,
			spaceCheckConstraint: s.checkConstraints,
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

// fkGeneratedIndexOwners returns the keys of synced foreign keys that could
// have implicitly generated idx as a backing index (HA_GENERATED_KEY): a plain,
// non-unique, non-primary, comment-less BTREE index with no prefix lengths or
// descending parts whose columns match the foreign key's columns in order.
func fkGeneratedIndexOwners(idx *metadata.IndexMetadata, fks []*metadata.ForeignKeyMetadata, foldTables bool) []string {
	if idx == nil || idx.GetPrimary() || idx.GetUnique() || idx.GetComment() != "" {
		return nil
	}
	if typ := idx.GetType(); typ != "" && !strings.EqualFold(typ, "BTREE") {
		return nil
	}
	for _, kl := range idx.GetKeyLength() {
		if kl > 0 {
			return nil
		}
	}
	for _, desc := range idx.GetDescending() {
		if desc {
			return nil
		}
	}
	exprs := idx.GetExpressions()
	if len(exprs) == 0 {
		return nil
	}
	var owners []string
	for _, fk := range fks {
		cols := fk.GetColumns()
		if len(cols) != len(exprs) {
			continue
		}
		match := true
		for i, c := range cols {
			if !strings.EqualFold(exprs[i], c) {
				match = false
				break
			}
		}
		if !match {
			continue
		}
		fkKey, ok := nameKey(spaceForeignKey, fk.GetName(), foldTables)
		if !ok {
			fkKey = ""
		}
		owners = append(owners, fkKey)
	}
	return owners
}
