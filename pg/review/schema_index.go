package review

import (
	"slices"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg/ast"
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
	// mayRead maps a relation to the views whose definition may read it
	// under a name several schemas have, and opaqueViews marks a view
	// whose definition the scan cannot read, which may read any.
	mayRead     map[tableRef][]tableRef
	opaqueViews bool
	// returnedBy maps a relation to the functions, by (schema, name), that
	// return its row type, and functions counts the functions and
	// procedures of each (schema, name).
	returnedBy map[tableRef][]tableRef
	functions  map[tableRef]int
	// constraintIndexes marks the indexes a constraint owns.
	constraintIndexes map[tableRef]bool
	// extensions lists the installed extensions by name, and routineKinds
	// the kinds of routine, function or procedure, each (schema, name) has.
	extensions   map[string]bool
	routineKinds map[tableRef]routineKind
	// noArgs marks the functions, by (schema, name), whose signature takes
	// no arguments, and signatures lists the argument signatures the scan
	// can read, written as argSignature writes them.
	noArgs     map[tableRef]bool
	signatures map[tableRef][]string
	// ownedSequences lists the sequences a column owns, which go with it.
	ownedSequences map[columnRef][]string
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
	// views are the views, for their columns.
	views map[string]*metadata.ViewMetadata
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
	// kindType is a type the change created, which takes a name in the
	// relation namespace's row types.
	kindType
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
		mayRead:     make(map[tableRef][]tableRef),
		returnedBy:  make(map[tableRef][]tableRef),
		functions:   make(map[tableRef]int),

		constraintIndexes: make(map[tableRef]bool),
		extensions:        make(map[string]bool),
		routineKinds:      make(map[tableRef]routineKind),
		noArgs:            make(map[tableRef]bool),
		signatures:        make(map[tableRef][]string),
		ownedSequences:    make(map[columnRef][]string),
	}
	for _, e := range db.GetExtensions() {
		idx.extensions[e.GetName()] = true
	}
	for _, s := range db.GetSchemas() {
		ns := &namespace{
			relations:  make(map[string]relationKind),
			tables:     make(map[string]*metadata.TableMetadata),
			dependents: make(map[string][]string),
			views:      make(map[string]*metadata.ViewMetadata),
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
			ns.views[v.GetName()] = v
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
				if q.GetOwnerColumn() != "" {
					c := columnRef{s.GetName(), q.GetOwnerTable(), q.GetOwnerColumn()}
					idx.ownedSequences[c] = append(idx.ownedSequences[c], q.GetName())
				}
			}
		}
		for _, c := range s.GetCompositeTypes() {
			ns.add(c.GetName(), kindCompositeType)
		}
		for _, e := range s.GetEnumTypes() {
			ns.types = append(ns.types, e.GetName())
		}
		ns.routines = len(s.GetFunctions()) + len(s.GetProcedures())
		for _, p := range s.GetProcedures() {
			fn := tableRef{s.GetName(), p.GetName()}
			idx.routineKinds[fn] |= kindProcedure
			idx.functions[fn]++
			if p.GetSignature() == p.GetName()+"()" {
				idx.noArgs[fn] = true
			}
			if sig, ok := readSignature(p.GetName(), p.GetSignature()); ok {
				idx.signatures[fn] = append(idx.signatures[fn], sig)
			}
		}
		for _, f := range s.GetFunctions() {
			fn := tableRef{s.GetName(), f.GetName()}
			idx.routineKinds[fn] |= kindOf(f)
			idx.functions[fn]++
			if f.GetSignature() == f.GetName()+"()" {
				idx.noArgs[fn] = true
			}
			if sig, ok := readSignature(f.GetName(), f.GetSignature()); ok {
				idx.signatures[fn] = append(idx.signatures[fn], sig)
			}
			for _, d := range f.GetDependencyTables() {
				t := tableRef{d.GetSchema(), d.GetTable()}
				idx.returnedBy[t] = append(idx.returnedBy[t], fn)
			}
		}
	}
	// A view may read a relation without naming a column of it, which the
	// snapshot's column dependencies leave out; its definition names it.
	path := parseSearchPath(db.GetSearchPath())
	for _, s := range db.GetSchemas() {
		for _, v := range s.GetViews() {
			idx.readersIn(tableRef{s.GetName(), v.GetName()}, v.GetDefinition(), path)
		}
		for _, v := range s.GetMaterializedViews() {
			idx.readersIn(tableRef{s.GetName(), v.GetName()}, v.GetDefinition(), path)
		}
	}
	return idx
}

// readersIn records the relations a synced view's definition names, CTE
// names aside. One the definition qualifies, the only one of its name, or
// the one the synced search path finds first the view certainly reads:
// the server writes a name unqualified only when the search path finds
// it. Otherwise it may read each relation of the name.
//
// bytebase: backend/plugin/db/pg/sync.go — getViewDependencies keeps only
// column dependencies (attnum > 0).
func (idx *schemaIndex) readersIn(view tableRef, definition string, path []string) {
	if strings.TrimSpace(definition) == "" {
		return
	}
	stmts, failed := parse(definition, statementRanges(definition))
	if failed != nil || len(stmts) != 1 {
		idx.opaqueViews = true
		return
	}
	query := stmts[0].node
	if v, ok := query.(*ast.ViewStmt); ok {
		query = v.Query
	}
	relationsIn(query, nil, func(rv *ast.RangeVar) {
		var candidates []tableRef
		for name, ns := range idx.schemas {
			if _, ok := ns.relations[rv.Relname]; ok && (rv.Schemaname == "" || rv.Schemaname == name) {
				candidates = append(candidates, tableRef{name, rv.Relname})
			}
		}
		if len(candidates) > 1 && rv.Schemaname == "" {
			if first, ok := idx.firstOnPath(path, rv.Relname); ok {
				candidates = []tableRef{first}
			}
		}
		if len(candidates) == 1 {
			if !slices.Contains(idx.readers[candidates[0]], view) {
				idx.readers[candidates[0]] = append(idx.readers[candidates[0]], view)
			}
			return
		}
		for _, c := range candidates {
			idx.mayRead[c] = append(idx.mayRead[c], view)
		}
	})
}

// firstOnPath returns the relation of a name the search path finds first,
// when no "$user" entry before it may name another schema that has one.
func (idx *schemaIndex) firstOnPath(path []string, name string) (tableRef, bool) {
	for i, schema := range path {
		if schema == "$user" {
			continue
		}
		ns := idx.schemas[schema]
		if ns == nil {
			continue
		}
		if _, ok := ns.relations[name]; !ok {
			continue
		}
		if slices.Contains(path[:i], "$user") {
			for other, ons := range idx.schemas {
				if _, ok := ons.relations[name]; ok && other != schema && !slices.Contains(path, other) {
					return tableRef{}, false
				}
			}
		}
		return tableRef{schema, name}, true
	}
	return tableRef{}, false
}

// mayInheritCheck reports whether a synced table's check constraint may be
// inherited: another table that may be its parent has a check constraint
// of that name and expression. The snapshot records no inheritance but
// partitioning, and a child cannot drop what it inherits.
func (idx *schemaIndex) mayInheritCheck(t tableRef, name string) bool {
	expr := ""
	for _, c := range idx.schemas[t.schema].tables[t.table].GetCheckConstraints() {
		if c.GetName() == name {
			expr = c.GetExpression()
		}
	}
	return slices.ContainsFunc(idx.parentsOf(t), func(p *metadata.TableMetadata) bool {
		return slices.ContainsFunc(p.GetCheckConstraints(), func(c *metadata.CheckConstraintMetadata) bool {
			return c.GetName() == name && c.GetExpression() == expr
		})
	})
}

// parentsOf returns the synced tables a table may inherit from: each
// other table whose columns, of the same names and types, are all its own.
func (idx *schemaIndex) parentsOf(t tableRef) []*metadata.TableMetadata {
	table := idx.schemas[t.schema].tables[t.table]
	if table == nil {
		return nil
	}
	types := make(map[string]string)
	for _, c := range table.GetColumns() {
		types[c.GetName()] = c.GetType()
	}
	var parents []*metadata.TableMetadata
	for name, ns := range idx.schemas {
		for tn, p := range ns.tables {
			if name == t.schema && tn == t.table || len(p.GetColumns()) == 0 {
				continue
			}
			if !slices.ContainsFunc(p.GetColumns(), func(c *metadata.ColumnMetadata) bool {
				ty, ok := types[c.GetName()]
				return !ok || ty != c.GetType()
			}) {
				parents = append(parents, p)
			}
		}
	}
	return parents
}

// readSignature reads a synced routine's signature, the name and the
// identity arguments PostgreSQL writes (each an optional mode, an optional
// name, and a type), into the form argSignature gives a statement's. ok is
// false for an argument it cannot read.
//
// bytebase: backend/plugin/db/pg/sync.go — getFunctions (Signature)
func readSignature(name, sig string) (string, bool) {
	if !strings.HasPrefix(sig, name+"(") || !strings.HasSuffix(sig, ")") {
		return "", false
	}
	args := strings.TrimSuffix(strings.TrimPrefix(sig, name+"("), ")")
	if args == "" {
		return "", true
	}
	var out []string
	for _, arg := range strings.Split(args, ", ") {
		t, ok := argType(arg)
		if !ok {
			return "", false
		}
		out = append(out, t)
	}
	return strings.Join(out, ","), true
}

// sqlTypes maps the names PostgreSQL writes built-in types under to the
// ones the parser gives them.
var sqlTypes = map[string]string{
	"integer": "int4", "bigint": "int8", "smallint": "int2", "boolean": "bool", "real": "float4",
	"double precision": "float8", "character varying": "varchar", "character": "bpchar",
	"timestamp without time zone": "timestamp", "timestamp with time zone": "timestamptz",
	"time without time zone": "time", "time with time zone": "timetz", "bit varying": "varbit",
	"numeric": "numeric", "text": "text", "uuid": "uuid", "date": "date", "interval": "interval",
	"json": "json", "jsonb": "jsonb", "bytea": "bytea",
}

// argType reads one identity argument into a type as argSignature writes
// it: a built-in type it knows, with an array marker. ok is false for any
// other.
func argType(arg string) (string, bool) {
	words := strings.Fields(arg)
	if len(words) > 0 && (words[0] == "IN" || words[0] == "INOUT" || words[0] == "VARIADIC") {
		words = words[1:]
	}
	array := false
	if n := len(words); n > 0 && strings.HasSuffix(words[n-1], "[]") {
		words[n-1] = strings.TrimSuffix(words[n-1], "[]")
		array = true
	}
	// The type is the whole rest, or the rest after an argument name.
	for _, rest := range [][]string{words, words[min(1, len(words)):]} {
		if t, ok := sqlTypes[strings.Join(rest, " ")]; ok && len(rest) > 0 {
			if array {
				t += "[]"
			}
			return t, true
		}
	}
	return "", false
}

// routineKind is a set of kinds of routine.
type routineKind int

const (
	kindFunction routineKind = 1 << iota
	kindProcedure
)

// kindOf tells a synced routine's kind. Bytebase lists a PostgreSQL
// procedure among the functions; its definition says which it is, and a
// routine without one may be either.
//
// bytebase: backend/plugin/db/pg/sync.go — listFunctionQuery
func kindOf(f *metadata.FunctionMetadata) routineKind {
	def := strings.ToUpper(f.GetDefinition())
	switch {
	case def == "":
		return kindFunction | kindProcedure
	case strings.Contains(def, "CREATE OR REPLACE PROCEDURE"):
		return kindProcedure
	}
	return kindFunction
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
