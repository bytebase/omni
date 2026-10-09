package review

import (
	"strings"

	"github.com/bytebase/omni/metadata"
)

// schemaIndex is a target's synced schema arranged for the questions the
// rules ask of it: what a name in a schema's relation namespace is, what
// constraints a table has, and which foreign keys and views depend on a
// table. It reads the metadata as synced; nothing is rebuilt from it.
type schemaIndex struct {
	database string
	schemas  map[string]*namespace
	// foreignKeys lists every foreign key of every table.
	foreignKeys []foreignKeyRef
	// readers maps a table to the views and materialized views that read
	// it, and readsColumn marks the columns they read.
	readers     map[tableRef][]tableRef
	readsColumn map[columnRef]bool
}

// namespace is one schema: the kind of every name in its relation
// namespace (tables, partitions, views, materialized views, foreign
// tables, sequences, indexes, and composite types share it) and its
// tables.
type namespace struct {
	relations map[string]relationKind
	tables    map[string]*metadata.TableMetadata
	// partitions, indexes, and owned sequences of each table, and the
	// indexes of each materialized view, which go with it when it is
	// dropped.
	dependents map[string][]string
}

type relationKind int

const (
	kindTable relationKind = iota + 1
	kindPartition
	// kindOther is a view, materialized view, foreign table, sequence,
	// index, or composite type, or a name the snapshot lists twice.
	kindOther
)

type tableRef struct{ schema, table string }

type columnRef struct{ schema, table, column string }

type foreignKeyRef struct {
	owner      tableRef // the referencing table
	name       string
	referenced tableRef
	columns    []string // the referenced columns
}

func newSchemaIndex(db *metadata.DatabaseSchemaMetadata) *schemaIndex {
	idx := &schemaIndex{
		database:    db.GetName(),
		schemas:     make(map[string]*namespace),
		readers:     make(map[tableRef][]tableRef),
		readsColumn: make(map[columnRef]bool),
	}
	for _, s := range db.GetSchemas() {
		ns := &namespace{
			relations:  make(map[string]relationKind),
			tables:     make(map[string]*metadata.TableMetadata),
			dependents: make(map[string][]string),
		}
		idx.schemas[s.GetName()] = ns
		for _, t := range s.GetTables() {
			ns.add(t.GetName(), kindTable)
			ns.tables[t.GetName()] = t
			for _, i := range t.GetIndexes() {
				ns.add(i.GetName(), kindOther)
				ns.dependents[t.GetName()] = append(ns.dependents[t.GetName()], i.GetName())
			}
			ns.addPartitions(t.GetName(), t.GetPartitions())
			for _, fk := range t.GetForeignKeys() {
				ref := tableRef{fk.GetReferencedSchema(), fk.GetReferencedTable()}
				if ref.schema == "" {
					ref.schema = s.GetName()
				}
				idx.foreignKeys = append(idx.foreignKeys, foreignKeyRef{owner: tableRef{s.GetName(), t.GetName()}, name: fk.GetName(), referenced: ref, columns: fk.GetReferencedColumns()})
			}
		}
		for _, v := range s.GetViews() {
			ns.add(v.GetName(), kindOther)
			idx.addReader(tableRef{s.GetName(), v.GetName()}, v.GetDependencyColumns())
		}
		for _, v := range s.GetMaterializedViews() {
			ns.add(v.GetName(), kindOther)
			idx.addReader(tableRef{s.GetName(), v.GetName()}, v.GetDependencyColumns())
			for _, i := range v.GetIndexes() {
				ns.add(i.GetName(), kindOther)
				ns.dependents[v.GetName()] = append(ns.dependents[v.GetName()], i.GetName())
			}
		}
		for _, t := range s.GetExternalTables() {
			ns.add(t.GetName(), kindOther)
		}
		for _, q := range s.GetSequences() {
			ns.add(q.GetName(), kindOther)
			if q.GetOwnerTable() != "" {
				ns.dependents[q.GetOwnerTable()] = append(ns.dependents[q.GetOwnerTable()], q.GetName())
			}
		}
		for _, c := range s.GetCompositeTypes() {
			ns.add(c.GetName(), kindOther)
		}
	}
	return idx
}

// add records a name. A name the snapshot lists twice is not known to be
// a table.
func (ns *namespace) add(name string, kind relationKind) {
	if _, dup := ns.relations[name]; dup {
		kind = kindOther
	}
	ns.relations[name] = kind
}

func (ns *namespace) addPartitions(table string, partitions []*metadata.TablePartitionMetadata) {
	for _, p := range partitions {
		ns.add(p.GetName(), kindPartition)
		ns.dependents[table] = append(ns.dependents[table], p.GetName())
		for _, i := range p.GetIndexes() {
			ns.add(i.GetName(), kindOther)
			ns.dependents[table] = append(ns.dependents[table], i.GetName())
		}
		ns.addPartitions(table, p.GetSubpartitions())
	}
}

func (idx *schemaIndex) addReader(view tableRef, columns []*metadata.DependencyColumn) {
	seen := make(map[tableRef]bool)
	for _, c := range columns {
		t := tableRef{c.GetSchema(), c.GetTable()}
		idx.readsColumn[columnRef{t.schema, t.table, c.GetColumn()}] = true
		if !seen[t] {
			seen[t] = true
			idx.readers[t] = append(idx.readers[t], view)
		}
	}
}

// unlisted reports whether some schema of the database is not named in
// the search path, so that "$user" could name it.
func (idx *schemaIndex) unlisted(path []string) bool {
	listed := make(map[string]bool, len(path))
	for _, name := range path {
		listed[name] = true
	}
	for name := range idx.schemas {
		if !listed[name] && !strings.HasPrefix(name, "pg_") && name != "information_schema" {
			return true
		}
	}
	return false
}

// referenced reports whether a foreign key not in gone references the
// table's key on exactly these columns, which makes the server refuse to
// drop the key without CASCADE.
func (idx *schemaIndex) referenced(t tableRef, columns []string, gone map[[3]string]bool) bool {
	for _, fk := range idx.foreignKeys {
		if fk.referenced == t && !fk.in(gone) && sameColumns(fk.columns, columns) {
			return true
		}
	}
	return false
}

// in reports whether the foreign key is in a set of dropped constraints.
func (fk foreignKeyRef) in(gone map[[3]string]bool) bool {
	return gone[[3]string{fk.owner.schema, fk.owner.table, fk.name}]
}

// referencesColumn reports whether a foreign key not in gone references a
// key of the table that has the column.
func (idx *schemaIndex) referencesColumn(t tableRef, column string, gone map[[3]string]bool) bool {
	for _, fk := range idx.foreignKeys {
		if fk.referenced != t || fk.in(gone) {
			continue
		}
		for _, c := range fk.columns {
			if name := columnName(c); name == column || name == "" {
				return true
			}
		}
	}
	return false
}

// referencedAny reports whether a foreign key not in gone references the
// table.
func (idx *schemaIndex) referencedAny(t tableRef, gone map[[3]string]bool) bool {
	for _, fk := range idx.foreignKeys {
		if fk.referenced == t && !fk.in(gone) {
			return true
		}
	}
	return false
}

func sameColumns(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	set := make(map[string]bool, len(a))
	for _, c := range a {
		set[columnName(c)] = true
	}
	for _, c := range b {
		if !set[columnName(c)] {
			return false
		}
	}
	return true
}

// tableConstraint is a constraint of a table as the snapshot lists it.
type tableConstraint struct {
	// kind is the rule's name for it, or "" for an exclusion constraint.
	kind string
	// columns are a key's columns; nil when an index expression is not a
	// plain column.
	columns []string
}

// constraint finds a table's constraint by name.
func constraint(t *metadata.TableMetadata, name string) (tableConstraint, bool) {
	for _, fk := range t.GetForeignKeys() {
		if fk.GetName() == name {
			return tableConstraint{kind: "foreign key"}, true
		}
	}
	for _, c := range t.GetCheckConstraints() {
		if c.GetName() == name {
			return tableConstraint{kind: "check constraint"}, true
		}
	}
	for _, x := range t.GetExcludeConstraints() {
		if x.GetName() == name {
			return tableConstraint{}, true
		}
	}
	for _, i := range t.GetIndexes() {
		if i.GetName() != name || !i.GetIsConstraint() {
			continue
		}
		switch {
		case i.GetPrimary():
			return tableConstraint{kind: "primary key", columns: keyColumns(i)}, true
		case i.GetUnique():
			return tableConstraint{kind: "unique constraint", columns: keyColumns(i)}, true
		default:
			return tableConstraint{}, true
		}
	}
	return tableConstraint{}, false
}

// primaryKey returns a table's primary key index, or nil.
func primaryKey(t *metadata.TableMetadata) *metadata.IndexMetadata {
	for _, i := range t.GetIndexes() {
		if i.GetPrimary() {
			return i
		}
	}
	return nil
}

// keyColumns returns the columns of a key index, or nil when one of its
// expressions is not a plain column name.
func keyColumns(i *metadata.IndexMetadata) []string {
	var columns []string
	for _, e := range i.GetExpressions() {
		c := columnName(e)
		if c == "" {
			return nil
		}
		columns = append(columns, c)
	}
	return columns
}

// columnName reads a column as the snapshot writes it in an index or a
// foreign key: a plain lower-case name, or a double-quoted identifier.
// Anything else is an expression, and the result is "".
func columnName(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return strings.ReplaceAll(s[1:len(s)-1], `""`, `"`)
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !(c >= 'a' && c <= 'z' || c == '_' || c >= '0' && c <= '9' && i > 0 || c == '$' && i > 0) {
			return ""
		}
	}
	return s
}
