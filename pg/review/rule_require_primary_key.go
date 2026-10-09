package review

import (
	"slices"
	"strings"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/review"
)

// RequirePrimaryKey reports two things, each only when the change ends
// that way: a table a CREATE TABLE made without a primary key, and a table
// that had a primary key in the synced schema and lost it. A table that had
// none before the change is not reported, and neither is a table made by
// CREATE TABLE AS or SELECT INTO, a partition, whose keys come from its
// parent, or a temporary table.
//
// Without simulating the change, the rule follows the table by name. It is
// settled, and not reported, by any later statement that may give it a
// key, drop it, or make it a partition, under any name it may have by
// then; every other statement leaves a table's lack of a key as it is.

// pendingTable is a table without a primary key that the rule reports at
// the end of the change unless a later statement settles it.
type pendingTable struct {
	// names are the names the table may have, schema the schema it is in,
	// or "" when it may be in any.
	names  map[string]bool
	schema string
	// existed marks a table of the synced schema, and parents names the
	// tables one made with INHERITS inherits from, in the schema each
	// resolved to or was written in, under every name a rename or a move
	// may have given it: a cascading drop of a parent may drop either.
	existed bool
	parents []tableRef

	// made numbers the touch that created a table the change made.
	made int

	statement int
	rng       review.Range
	message   string
	settled   bool
}

// inheritsFrom names the tables CREATE TABLE ... INHERITS lists.
func (s *scan) inheritsFrom(v *ast.CreateStmt) []tableRef {
	if v.InhRelations == nil {
		return nil
	}
	var parents []tableRef
	for _, item := range v.InhRelations.Items {
		if rv, ok := item.(*ast.RangeVar); ok {
			parents = append(parents, tableRef{s.schemaOf(rv), rv.Relname})
		}
	}
	return parents
}

// mayBe reports whether two relations, each in a schema or "" when that is
// not known, may be one.
func mayBe(a, b tableRef) bool {
	return a.table == b.table && (a.schema == "" || b.schema == "" || a.schema == b.schema)
}

// followParents records that a parent a statement renamed or moved keeps
// its children under its new schema and name.
func (s *scan) followParents(rv *ast.RangeVar, to tableRef) {
	from := tableRef{s.schemaOf(rv), rv.Relname}
	for _, p := range s.pending {
		if slices.ContainsFunc(p.parents, func(t tableRef) bool { return mayBe(t, from) }) {
			p.parents = append(p.parents, to)
		}
	}
}

// matching returns the pending tables a name may refer to: any with that
// name, unless both the name and the table are in known, different
// schemas.
func (s *scan) matching(rv *ast.RangeVar) []*pendingTable {
	var out []*pendingTable
	for _, p := range s.pending {
		if !p.names[rv.Relname] || p.schema != "" && rv.Schemaname != "" && rv.Schemaname != p.schema {
			continue
		}
		if rv.Schemaname == "" && p.schema != "" && s.resolvesBefore(rv.Relname, p.schema) {
			continue
		}
		out = append(out, p)
	}
	return out
}

// resolvesBefore reports whether an unqualified name certainly means a
// synced relation of a schema the search path puts before the given one.
func (s *scan) resolvesBefore(name, schema string) bool {
	if s.index == nil || strings.HasPrefix(name, "pg_") {
		return false
	}
	path, ok := s.searchPath()
	if !ok {
		return false
	}
	for _, sch := range path {
		if sch == schema || sch == "information_schema" || s.schemas[sch] || s.isTouched(&ast.RangeVar{Schemaname: sch, Relname: name}) {
			return false
		}
		if ns := s.index.schemas[sch]; ns != nil {
			if _, exists := ns.relations[name]; exists {
				return true
			}
		}
	}
	return false
}

// keyed settles the pending tables a statement may give a primary key.
func (s *scan) keyed(rv *ast.RangeVar) {
	for _, p := range s.matching(rv) {
		p.settled = true
	}
}

// lostKey records that a statement removed the primary key of a table of
// the synced schema.
func (s *scan) lostKey(st *statement, v *ast.AlterTableStmt, cmd *ast.AlterTableCmd, t tableRef) {
	if !s.on[review.RequirePrimaryKey] {
		return
	}
	for _, p := range s.pending {
		if p.existed && !p.settled && p.schema == t.schema && p.names[t.table] {
			return
		}
	}
	s.pending = append(s.pending, &pendingTable{
		names:     map[string]bool{t.table: true},
		schema:    t.schema,
		existed:   true,
		statement: st.index,
		rng:       rangeOf(cmd.Loc),
		message:   "removes the primary key of " + relation(v.Relation) + ", and the change adds none back",
	})
}

// createTable records a CREATE TABLE that makes a table without a primary
// key. CREATE TABLE IF NOT EXISTS is left out when a relation of that name
// may already exist, since then it makes nothing.
func (s *scan) createTable(st *statement, v *ast.CreateStmt) {
	if !s.on[review.RequirePrimaryKey] || v.Relation == nil || v.Partbound != nil || isTemp(v.Relation) || s.createsTemp(v.Relation) {
		return
	}
	if v.IfNotExists && !s.createsNew(v.Relation) {
		return
	}
	keyed, known := s.createsKey(v)
	if keyed || !known || s.mayBeGenerated(v.Relation) {
		return
	}
	s.pending = append(s.pending, &pendingTable{
		names:     map[string]bool{v.Relation.Relname: true},
		schema:    s.pendingSchema(v.Relation),
		parents:   s.inheritsFrom(v),
		statement: st.index,
		rng:       rangeOf(v.Loc),
		message:   "creates table " + relation(v.Relation) + " without a primary key",
	})
}

// keylessNew reports whether a qualified name certainly means a table the
// change made without a primary key, which no statement since may have
// given one: a pending table of that one name in that schema, whose name
// no statement used since it was made.
func (s *scan) keylessNew(rv *ast.RangeVar) bool {
	if rv.Schemaname == "" {
		return false
	}
	for _, p := range s.pending {
		if !p.existed && !p.settled && p.schema == rv.Schemaname && len(p.names) == 1 && p.names[rv.Relname] && s.lastTouch[rv.Relname] == p.made {
			return true
		}
	}
	return false
}

// pendingSchema is the schema a CREATE TABLE puts its table in, as written
// or as the search path resolves it, or "" when neither tells.
func (s *scan) pendingSchema(rv *ast.RangeVar) string {
	if rv.Schemaname != "" {
		return rv.Schemaname
	}
	schema, _ := s.creationSchema(rv)
	return schema
}

// createsTemp reports whether CREATE puts an unqualified relation in the
// session's temporary schema, or a system one: the known search path
// names one before any schema that exists.
func (s *scan) createsTemp(rv *ast.RangeVar) bool {
	if rv == nil || rv.Schemaname != "" || s.index == nil {
		return false
	}
	path, ok := s.searchPath()
	if !ok {
		return false
	}
	for _, name := range path {
		switch {
		case strings.HasPrefix(name, "pg_"):
			return true
		case s.schemas[name]:
			if !s.schemaGone[name] {
				return false
			}
		case s.index.schemas[name] != nil:
			return false
		}
	}
	return false
}

// createsNew reports whether CREATE TABLE IF NOT EXISTS creates a table:
// the schema it creates in is known, the synced one, and holds no relation
// of that name, or one the change created, where no statement made the
// name. An unqualified name is created in the first schema of the search
// path that exists.
func (s *scan) createsNew(rv *ast.RangeVar) bool {
	if s.index != nil && rv.Schemaname != "" && s.schemas[rv.Schemaname] && !s.schemaGone[rv.Schemaname] && !isTemp(rv) {
		// A schema the change created starts empty.
		return !s.isTouched(rv)
	}
	schema, ok := s.creationSchema(rv)
	if !ok {
		return false
	}
	ns := s.index.schemas[schema]
	if s.freed[[2]string{schema, rv.Relname}] {
		return true
	}
	if s.isTouched(&ast.RangeVar{Schemaname: schema, Relname: rv.Relname}) {
		return false
	}
	_, exists := ns.relations[rv.Relname]
	return !exists
}

// nameUncertain reports whether an earlier statement may have made a
// relation of that name where CREATE puts it, without the scan knowing it
// gone again: CREATE TABLE may then fail, and creates nothing the rule can
// report.
func (s *scan) nameUncertain(rv *ast.RangeVar) bool {
	if schema, ok := s.creationSchema(rv); ok {
		return s.isTouched(&ast.RangeVar{Schemaname: schema, Relname: rv.Relname}) && !s.freed[[2]string{schema, rv.Relname}]
	}
	if rv.Schemaname != "" {
		return s.isTouched(rv) && !s.freed[[2]string{rv.Schemaname, rv.Relname}]
	}
	return s.touchedName[rv.Relname]
}

// generatedName reports whether a name has the form of one the server
// generates for an index or a sequence: a suffix of _idx, _key, _pkey,
// _excl, or _seq, then digits the server adds to avoid a name in use. A
// table of such a name may collide with a relation an earlier statement
// made, which the scan does not name.
//
// pg: src/backend/commands/indexcmds.c — ChooseIndexName
func generatedName(name string) bool {
	name = strings.TrimRight(name, "0123456789")
	for _, suffix := range []string{"_idx", "_key", "_pkey", "_excl", "_seq"} {
		if strings.HasSuffix(name, suffix) {
			return true
		}
	}
	return false
}

// existsWhereCreated reports whether the schema a CREATE puts the relation
// in certainly holds a relation of that name already: the schema is the
// synced one, and no statement of the change dropped, renamed, or made the
// name.
func (s *scan) existsWhereCreated(rv *ast.RangeVar) bool {
	schema, ok := s.creationSchema(rv)
	if !ok || s.freed[[2]string{schema, rv.Relname}] || s.isTouched(&ast.RangeVar{Schemaname: schema, Relname: rv.Relname}) {
		return false
	}
	_, exists := s.index.schemas[schema].relations[rv.Relname]
	return exists
}

// creationSchema returns the synced schema a CREATE puts a relation in:
// the one written, or the first schema of the search path that exists. ok
// is false when that is not a schema the scan knows as synced.
func (s *scan) creationSchema(rv *ast.RangeVar) (string, bool) {
	if s.index == nil || isTemp(rv) {
		return "", false
	}
	schema := rv.Schemaname
	if schema == "" {
		path, ok := s.searchPath()
		if !ok {
			return "", false
		}
		for _, name := range path {
			// A system schema is not in the snapshot, and a schema the
			// change created or dropped may or may not exist.
			if s.schemas[name] || strings.HasPrefix(name, "pg_") || name == "information_schema" {
				return "", false
			}
			if s.index.schemas[name] != nil {
				schema = name
				break
			}
		}
	}
	if s.index.schemas[schema] == nil || s.schemas[schema] {
		return "", false
	}
	return schema, true
}

// likeIndexes is the LIKE option that copies indexes, the primary key
// among them; INCLUDING ALL sets it too.
//
// pg: src/include/nodes/parsenodes.h — CREATE_TABLE_LIKE_INDEXES
const likeIndexes = 1 << 6

// createsKey reports whether CREATE TABLE gives the table a primary key: a
// PRIMARY KEY of a column or of the table, or LIKE ... INCLUDING INDEXES of
// a table that has one. known is false when that depends on a table whose
// key the scan cannot tell.
func (s *scan) createsKey(v *ast.CreateStmt) (keyed, known bool) {
	if v.TableElts == nil {
		return false, true
	}
	known = true
	for _, item := range v.TableElts.Items {
		switch e := item.(type) {
		case *ast.ColumnDef:
			for _, c := range constraintsOf(e.Constraints) {
				if c.Contype == ast.CONSTR_PRIMARY {
					keyed = true
				}
			}
		case *ast.Constraint:
			if e.Contype == ast.CONSTR_PRIMARY {
				keyed = true
			}
		case *ast.TableLikeClause:
			if e.Options&likeIndexes == 0 {
				continue
			}
			schema, kind, ok := s.lookup(e.Relation)
			switch {
			case !ok && s.keylessNew(e.Relation):
				// A table the change made without a key has none to copy.
			case !ok || kind == kindPartition || kind == kindTable && s.isUnsettled(tableRef{schema, e.Relation.Relname}):
				known = false
			case kind == kindTable:
				t := tableRef{schema, e.Relation.Relname}
				pk := primaryKey(s.index.schemas[schema].tables[t.table])
				switch {
				case s.newKeys[[2]string{t.schema, t.table}] || s.newKeys[[2]string{"", t.table}]:
					// A statement may have given it a key since.
					known = false
				case pk == nil:
				case s.dropped[[3]string{t.schema, t.table, pk.GetName()}]:
					// The change dropped its key for certain.
				case !s.constraintKnown(t, pk.GetName()):
					known = false
				default:
					keyed = true
				}
			}
			// A view, materialized view, foreign table, or composite type
			// has no primary key to copy.
		}
	}
	return keyed, known
}
