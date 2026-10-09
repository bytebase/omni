package review

import (
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
	// existed marks a table of the synced schema, and inherits one made
	// with INHERITS: a cascading drop of a parent may drop either.
	existed  bool
	inherits bool

	statement int
	rng       review.Range
	message   string
	settled   bool
}

// matching returns the pending tables a name may refer to: any with that
// name, unless both the name and the table are in known, different
// schemas.
func (s *scan) matching(rv *ast.RangeVar) []*pendingTable {
	var out []*pendingTable
	for _, p := range s.pending {
		if p.names[rv.Relname] && (p.schema == "" || rv.Schemaname == "" || rv.Schemaname == p.schema) {
			out = append(out, p)
		}
	}
	return out
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
	if !s.on[review.RequirePrimaryKey] || v.Relation == nil || v.Partbound != nil || isTemp(v.Relation) {
		return
	}
	if v.IfNotExists && !s.createsNew(v.Relation) {
		return
	}
	keyed, known := s.createsKey(v)
	if keyed || !known {
		return
	}
	s.pending = append(s.pending, &pendingTable{
		names:     map[string]bool{v.Relation.Relname: true},
		schema:    v.Relation.Schemaname,
		inherits:  v.InhRelations != nil && len(v.InhRelations.Items) > 0,
		statement: st.index,
		rng:       rangeOf(v.Loc),
		message:   "creates table " + relation(v.Relation) + " without a primary key",
	})
}

// createsNew reports whether CREATE TABLE IF NOT EXISTS creates a table:
// the schema it creates in is known, the synced one, and holds no relation
// of that name. An unqualified name is created in the first schema of the
// search path that exists.
func (s *scan) createsNew(rv *ast.RangeVar) bool {
	if s.index == nil {
		return false
	}
	schema := rv.Schemaname
	if schema == "" {
		path, ok := s.searchPath()
		if !ok {
			return false
		}
		for _, name := range path {
			// A system schema is not in the snapshot, and a schema the
			// change created or dropped may or may not exist.
			if s.schemas[name] || strings.HasPrefix(name, "pg_") || name == "information_schema" {
				return false
			}
			if s.index.schemas[name] != nil {
				schema = name
				break
			}
		}
	}
	ns := s.index.schemas[schema]
	if ns == nil || s.schemas[schema] {
		return false
	}
	if s.freed[[2]string{schema, rv.Relname}] {
		return true
	}
	if s.isTouched(&ast.RangeVar{Schemaname: schema, Relname: rv.Relname}) {
		return false
	}
	_, exists := ns.relations[rv.Relname]
	return !exists
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
			case !ok || kind == kindPartition || kind == kindTable && s.isUnsettled(tableRef{schema, e.Relation.Relname}):
				known = false
			case kind == kindTable:
				pk := primaryKey(s.index.schemas[schema].tables[e.Relation.Relname])
				if pk != nil && !s.constraintKnown(tableRef{schema, e.Relation.Relname}, pk.GetName()) {
					known = false
				} else if pk != nil {
					keyed = true
				}
			}
			// A view, materialized view, foreign table, or composite type
			// has no primary key to copy.
		}
	}
	return keyed, known
}
