package review

import (
	"slices"
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
	// it, and readsColumn a column to the ones that read it.
	readers     map[tableRef][]tableRef
	readsColumn map[columnRef][]tableRef
	// returnedBy maps a relation to the functions, by (schema, name), that
	// return its row type, and functions counts the functions of each
	// (schema, name).
	returnedBy map[tableRef][]tableRef
	functions  map[tableRef]int
	// constraintIndexes marks the indexes a constraint owns.
	constraintIndexes map[tableRef]bool
	// extensions lists the installed extensions by name.
	extensions map[string]bool
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
	// types are the enum types, and routines counts the functions and
	// procedures, which keep a schema from being empty.
	types    []string
	routines int
}

type relationKind int

const (
	kindTable relationKind = iota + 1
	kindPartition
	kindView
	kindMatView
	kindForeignTable
	kindSequence
	kindIndex
	kindCompositeType
	// kindAmbiguous is a name the snapshot lists twice.
	kindAmbiguous
)

type tableRef struct{ schema, table string }

type columnRef struct{ schema, table, column string }

type foreignKeyRef struct {
	owner      tableRef // the referencing table
	local      []string // the referencing columns
	name       string
	referenced tableRef
	columns    []string // the referenced columns
}

func newSchemaIndex(db *metadata.DatabaseSchemaMetadata) *schemaIndex {
	idx := &schemaIndex{
		database:    db.GetName(),
		schemas:     make(map[string]*namespace),
		readers:     make(map[tableRef][]tableRef),
		readsColumn: make(map[columnRef][]tableRef),
		returnedBy:  make(map[tableRef][]tableRef),
		functions:   make(map[tableRef]int),

		constraintIndexes: make(map[tableRef]bool),
		extensions:        make(map[string]bool),
	}
	for _, e := range db.GetExtensions() {
		idx.extensions[e.GetName()] = true
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
				ns.add(i.GetName(), kindIndex)
				ns.dependents[t.GetName()] = append(ns.dependents[t.GetName()], i.GetName())
				if i.GetIsConstraint() {
					idx.constraintIndexes[tableRef{s.GetName(), i.GetName()}] = true
				}
			}
			ns.addPartitions(t.GetName(), t.GetPartitions())
			idx.partitionConstraintIndexes(s.GetName(), t.GetPartitions())
			for _, fk := range t.GetForeignKeys() {
				ref := tableRef{fk.GetReferencedSchema(), fk.GetReferencedTable()}
				if ref.schema == "" {
					ref.schema = s.GetName()
				}
				idx.foreignKeys = append(idx.foreignKeys, foreignKeyRef{owner: tableRef{s.GetName(), t.GetName()}, local: fk.GetColumns(), name: fk.GetName(), referenced: ref, columns: fk.GetReferencedColumns()})
			}
		}
		for _, v := range s.GetViews() {
			ns.add(v.GetName(), kindView)
			idx.addReader(tableRef{s.GetName(), v.GetName()}, v.GetDependencyColumns())
		}
		for _, v := range s.GetMaterializedViews() {
			ns.add(v.GetName(), kindMatView)
			idx.addReader(tableRef{s.GetName(), v.GetName()}, v.GetDependencyColumns())
			for _, i := range v.GetIndexes() {
				ns.add(i.GetName(), kindIndex)
				ns.dependents[v.GetName()] = append(ns.dependents[v.GetName()], i.GetName())
			}
		}
		for _, t := range s.GetExternalTables() {
			ns.add(t.GetName(), kindForeignTable)
		}
		for _, q := range s.GetSequences() {
			ns.add(q.GetName(), kindSequence)
			if q.GetOwnerTable() != "" {
				ns.dependents[q.GetOwnerTable()] = append(ns.dependents[q.GetOwnerTable()], q.GetName())
			}
		}
		for _, c := range s.GetCompositeTypes() {
			ns.add(c.GetName(), kindCompositeType)
		}
		for _, e := range s.GetEnumTypes() {
			ns.types = append(ns.types, e.GetName())
		}
		ns.routines = len(s.GetFunctions()) + len(s.GetProcedures())
		for _, f := range s.GetFunctions() {
			fn := tableRef{s.GetName(), f.GetName()}
			idx.functions[fn]++
			for _, d := range f.GetDependencyTables() {
				t := tableRef{d.GetSchema(), d.GetTable()}
				idx.returnedBy[t] = append(idx.returnedBy[t], fn)
			}
		}
	}
	return idx
}

// add records a name. A name the snapshot lists twice is not known to be
// a table.
func (ns *namespace) add(name string, kind relationKind) {
	if _, dup := ns.relations[name]; dup {
		kind = kindAmbiguous
	}
	ns.relations[name] = kind
}

// addPartitions records a table's partitions, at any depth, with their
// indexes, as what goes with the table and with each partition above
// them.
func (ns *namespace) addPartitions(table string, partitions []*metadata.TablePartitionMetadata) {
	ns.addPartitionsUnder([]string{table}, partitions)
}

func (ns *namespace) addPartitionsUnder(owners []string, partitions []*metadata.TablePartitionMetadata) {
	for _, p := range partitions {
		ns.add(p.GetName(), kindPartition)
		names := []string{p.GetName()}
		for _, i := range p.GetIndexes() {
			ns.add(i.GetName(), kindIndex)
			names = append(names, i.GetName())
		}
		for _, owner := range owners {
			ns.dependents[owner] = append(ns.dependents[owner], names...)
		}
		// The partition's own indexes go with it too.
		ns.dependents[p.GetName()] = append(ns.dependents[p.GetName()], names[1:]...)
		ns.addPartitionsUnder(append(slices.Clone(owners), p.GetName()), p.GetSubpartitions())
	}
}

// partitionConstraintIndexes marks the constraint indexes of partitions.
func (idx *schemaIndex) partitionConstraintIndexes(schema string, partitions []*metadata.TablePartitionMetadata) {
	for _, p := range partitions {
		for _, i := range p.GetIndexes() {
			if i.GetIsConstraint() {
				idx.constraintIndexes[tableRef{schema, i.GetName()}] = true
			}
		}
		idx.partitionConstraintIndexes(schema, p.GetSubpartitions())
	}
}

func (idx *schemaIndex) addReader(view tableRef, columns []*metadata.DependencyColumn) {
	seen := make(map[tableRef]bool)
	for _, c := range columns {
		t := tableRef{c.GetSchema(), c.GetTable()}
		col := columnRef{t.schema, t.table, c.GetColumn()}
		idx.readsColumn[col] = append(idx.readsColumn[col], view)
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
