package review

import (
	"context"
	"slices"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/review"
)

// scan follows the change statement by statement for the rules that ask
// what a statement does to an object of the target's schema:
// DisallowDropConstraint and RequirePrimaryKey. It simulates nothing. It
// checks each statement against the synced schema directly, and keeps
// track of what earlier statements may have changed: a name whose meaning
// a statement may have changed is not looked up again, and a statement
// whose effect the scan cannot bound (procedural code, a rollback, a
// cascading drop of anything but a relation) ends it. Every way the scan
// can be unsure leads to no finding.
//
// A finding is a claim about its own statement, given that the statements
// before it ran: the statement it anchors never fails for a reason the
// SQL and the synced schema show. Whether the change as a whole can run
// is WalkThrough's question, not this scan's. Where the scan notices on
// the way that the server refuses a statement (a missing object, a name of
// the wrong kind, a dependent object without CASCADE) it reports nothing
// for it and ends, since it cannot say what follows; it does not look for
// every way a statement can fail. A finding after such a statement still
// describes the change as written. What user code does when it runs, a
// trigger fired by DML or a function a query calls, is outside what a
// static check can bound, and the scan does not predict it.
type scan struct {
	on    map[review.Rule]bool
	r     *reporter
	index *schemaIndex // nil without a schema

	// path is the search path as last set, "$user" unexpanded, or nil
	// when unknown. user is the current user and session the session
	// user, "" when unknown.
	path    []string
	user    string
	session string

	// touched lists the relation names a statement created, dropped,
	// renamed, or moved, by (schema, name), with "" for a schema the
	// statement left to the search path: the synced schema no longer says
	// what they name. touchedName lists the same names alone.
	touched     map[[2]string]bool
	touchedName map[string]bool
	// freed lists the (schema, name) a drop the scan resolved removed for
	// certain, until a statement uses the name again.
	freed map[[2]string]bool
	// unsettled lists the (schema, table) whose columns, constraints, or
	// inheritance a statement changed in a way the scan does not follow,
	// with "" for a schema the scan could not tell.
	unsettled map[[2]string]bool
	// constraints lists the (schema, table, constraint) names a statement
	// added, dropped, or renamed, with "" for a schema the scan could not
	// tell, and renamedKeys the names an index rename may have given to or
	// taken from a constraint, on whichever table.
	constraints map[[3]string]bool
	renamedKeys map[string]bool
	// dropped lists the (schema, table, constraint) a DROP CONSTRAINT of a
	// resolved table removed.
	dropped map[[3]string]bool
	// columns lists the (schema, table, column) names a statement added,
	// and renamedFrom and renamedTo the old and new names of renamed
	// columns, with "" for a schema the scan could not tell.
	columns     map[[3]string]bool
	renamedFrom map[[3]string]bool
	renamedTo   map[[3]string]bool
	// schemas lists the schemas a statement created or dropped.
	schemas map[string]bool
	// newReferences lists the keys the foreign keys the change created
	// reference, newReads the relations its views read, and newReturns the
	// relations whose row type its functions return, none of which the
	// synced schema records.
	newReferences []newReference
	newReads      []newDependency
	newReturns    []newDependency
	// pathVersion counts search path changes, so two unqualified names
	// are known to resolve alike only under the same path.
	pathVersion int
	// droppedFunctions lists the synced functions, by (schema, name), a
	// DROP FUNCTION certainly removed, and newFunctions counts the functions
	// the change created, by name as written.
	droppedFunctions map[tableRef]bool
	newFunctions     map[tableRef]int
	// replacedViews lists the synced views CREATE OR REPLACE VIEW
	// redefined, whose synced reads no longer hold, and reowned the
	// sequences, by name, whose owner ALTER SEQUENCE changed.
	replacedViews map[tableRef]bool
	reowned       map[string]bool
	// generated marks that a statement made relations under names the
	// server generates (an unnamed index, a key's index, a serial or
	// identity column's sequence), which the scan does not derive.
	generated bool

	// pending are the tables RequirePrimaryKey reports at the end unless a
	// later statement settles them.
	pending []*pendingTable
	stopped bool
}

// scanTarget scans the change on one target. RequirePrimaryKey is about
// the end of the change, so it reports only when the scan reaches it.
func scanTarget(ctx context.Context, stmts []statement, on map[review.Rule]bool, target review.Target, r *reporter) error {
	s := &scan{
		on:          on,
		r:           r,
		user:        target.SessionUser,
		session:     target.SessionUser,
		touched:     make(map[[2]string]bool),
		touchedName: make(map[string]bool),
		unsettled:   make(map[[2]string]bool),
		freed:       make(map[[2]string]bool),
		constraints: make(map[[3]string]bool),
		renamedKeys: make(map[string]bool),
		dropped:     make(map[[3]string]bool),
		columns:     make(map[[3]string]bool),
		renamedFrom: make(map[[3]string]bool),
		renamedTo:   make(map[[3]string]bool),
		schemas:     make(map[string]bool),

		droppedFunctions: make(map[tableRef]bool),
		newFunctions:     make(map[tableRef]int),
		replacedViews:    make(map[tableRef]bool),
		reowned:          make(map[string]bool),
	}
	if target.Schema != nil {
		s.index = newSchemaIndex(target.Schema)
		// A schema synced without its search path leaves it unknown.
		if path := parseSearchPath(target.Schema.GetSearchPath()); len(path) > 0 {
			s.path = path
		}
	}
	if target.Schema != nil && eventTriggerFires(target.Schema, stmts) {
		return nil
	}
	for i := range stmts {
		if err := ctx.Err(); err != nil {
			return err
		}
		s.statement(&stmts[i])
		if s.stopped {
			return nil
		}
	}
	for _, p := range s.pending {
		if !p.settled {
			r.report(review.RequirePrimaryKey, p.statement, p.rng, p.message)
		}
	}
	return nil
}

// eventTriggerFires reports whether an enabled event trigger of the target
// fires on a statement of the change the scan follows: one with no tag
// filter, or one whose tags name the statement's command. Its function may
// run more DDL or reject the statement, so what the statement does is no
// longer what it says.
func eventTriggerFires(schema *metadata.DatabaseSchemaMetadata, stmts []statement) bool {
	var triggers []*metadata.EventTriggerMetadata
	for _, t := range schema.GetEventTriggers() {
		if t.GetEnabled() {
			triggers = append(triggers, t)
		}
	}
	if len(triggers) == 0 {
		return false
	}
	tags := make(map[string]bool)
	for i := range stmts {
		if tag := commandTag(stmts[i].node); tag != "" {
			tags[tag] = true
			// CREATE FUNCTION and CREATE PROCEDURE share a node.
			if tag == "CREATE FUNCTION" {
				tags["CREATE PROCEDURE"] = true
			}
		}
	}
	if len(tags) == 0 {
		return false
	}
	for _, t := range triggers {
		if len(t.GetTags()) == 0 {
			return true
		}
		for _, tag := range t.GetTags() {
			if tags[strings.ToUpper(tag)] {
				return true
			}
		}
	}
	return false
}

// commandTag is the command tag of a statement the scan follows, as an
// event trigger's tag filter names it, or "" for one it does not follow.
//
// pg: src/include/tcop/cmdtaglist.h
func commandTag(n ast.Node) string {
	objectTag := func(t ast.ObjectType) string {
		switch t {
		case ast.OBJECT_TABLE:
			return "TABLE"
		case ast.OBJECT_VIEW:
			return "VIEW"
		case ast.OBJECT_MATVIEW:
			return "MATERIALIZED VIEW"
		case ast.OBJECT_FOREIGN_TABLE:
			return "FOREIGN TABLE"
		case ast.OBJECT_SEQUENCE:
			return "SEQUENCE"
		case ast.OBJECT_INDEX:
			return "INDEX"
		case ast.OBJECT_TYPE:
			return "TYPE"
		case ast.OBJECT_SCHEMA:
			return "SCHEMA"
		case ast.OBJECT_FUNCTION:
			return "FUNCTION"
		case ast.OBJECT_ROUTINE:
			return "ROUTINE"
		}
		return ""
	}
	prefixed := func(verb string, t ast.ObjectType) string {
		if o := objectTag(t); o != "" {
			return verb + " " + o
		}
		return ""
	}
	switch v := n.(type) {
	case *ast.CreateStmt:
		return "CREATE TABLE"
	case *ast.CreateForeignTableStmt:
		return "CREATE FOREIGN TABLE"
	case *ast.CreateTableAsStmt:
		if v.Objtype == ast.OBJECT_MATVIEW {
			return "CREATE MATERIALIZED VIEW"
		}
		return "CREATE TABLE AS"
	case *ast.SelectStmt:
		if v.IntoClause != nil {
			return "SELECT INTO"
		}
	case *ast.ViewStmt:
		return "CREATE VIEW"
	case *ast.CreateSeqStmt:
		return "CREATE SEQUENCE"
	case *ast.CompositeTypeStmt:
		return "CREATE TYPE"
	case *ast.IndexStmt:
		return "CREATE INDEX"
	case *ast.CreateFunctionStmt:
		return "CREATE FUNCTION"
	case *ast.CreateSchemaStmt:
		return "CREATE SCHEMA"
	case *ast.AlterTableStmt:
		return prefixed("ALTER", ast.ObjectType(v.ObjType))
	case *ast.RenameStmt:
		if v.RenameType == ast.OBJECT_COLUMN || v.RenameType == ast.OBJECT_TABCONSTRAINT {
			if tag := prefixed("ALTER", v.RelationType); tag != "" {
				return tag
			}
			return "ALTER TABLE"
		}
		return prefixed("ALTER", v.RenameType)
	case *ast.AlterObjectSchemaStmt:
		return prefixed("ALTER", v.ObjectType)
	case *ast.DropStmt:
		return prefixed("DROP", ast.ObjectType(v.RemoveType))
	}
	return ""
}

// stop ends the scan: nothing after this statement is known, and no table
// is known to end the change without a key.
func (s *scan) stop() {
	s.stopped = true
}

func (s *scan) statement(st *statement) {
	switch v := st.node.(type) {
	case *ast.AlterTableStmt:
		s.alterTable(st, v)
	case *ast.CreateStmt:
		s.create(st, v)
	case *ast.CreateForeignTableStmt:
		s.touch(v.Base.Relation)
	case *ast.CreateTableAsStmt:
		if v.Into != nil {
			s.touch(v.Into.Rel)
		}
		if v.Objtype == ast.OBJECT_MATVIEW && v.Into != nil {
			s.reads(v.Into.Rel, v.Query)
		}
	case *ast.ViewStmt:
		if v.Replace {
			if schema, kind, ok := s.lookup(v.View); ok && kind == kindView {
				s.replacedViews[tableRef{schema, v.View.Relname}] = true
			}
		}
		s.touch(v.View)
		s.reads(v.View, v.Query)
	case *ast.CreateFunctionStmt:
		s.returns(v)
	case *ast.AlterSeqStmt:
		if v.Sequence != nil && v.Options != nil && slices.ContainsFunc(v.Options.Items, func(n ast.Node) bool {
			d, ok := n.(*ast.DefElem)
			return ok && d.Defname == "owned_by"
		}) {
			s.reowned[v.Sequence.Relname] = true
		}
	case *ast.CreateSeqStmt:
		s.touch(v.Sequence)
	case *ast.CompositeTypeStmt:
		s.touch(v.Typevar)
	case *ast.IndexStmt:
		if v.Idxname == "" {
			s.generated = true
		} else if v.Relation != nil {
			s.touchName(v.Relation.Schemaname, v.Idxname)
		}
	case *ast.SelectStmt:
		switch {
		case v.IntoClause != nil:
			s.touch(v.IntoClause.Rel)
		case selectsInto(v):
			s.stop()
		default:
			s.data(v)
		}
	case *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.MergeStmt:
		s.data(v)
	case *ast.ExplainStmt:
		if !explainAnalyzes(v) {
			return
		}
		switch q := v.Query.(type) {
		case *ast.CreateTableAsStmt:
			if q.Into != nil {
				s.touch(q.Into.Rel)
			}
		case *ast.SelectStmt:
			if q.IntoClause != nil || selectsInto(q) {
				s.stop()
				return
			}
			s.data(q)
		default:
			s.data(q)
		}
	case *ast.RenameStmt:
		s.rename(v)
	case *ast.AlterObjectSchemaStmt:
		s.setSchema(v)
	case *ast.DropStmt:
		s.drop(v)
	case *ast.CreateSchemaStmt:
		s.createSchema(st, v)
	case *ast.DoStmt, *ast.CallStmt:
		// Procedural code may run any DDL.
		s.stop()
	case *ast.DropOwnedStmt, *ast.ImportForeignSchemaStmt, *ast.AlterExtensionStmt, *ast.DiscardStmt, *ast.LoadStmt:
		// What these drop, create, or reset is not in the statement.
		s.stop()
	case *ast.CreateEventTrigStmt, *ast.AlterEventTrigStmt:
		// An event trigger runs code on every later DDL statement.
		s.stop()
	case *ast.TransactionStmt:
		switch v.Kind {
		case ast.TRANS_STMT_ROLLBACK, ast.TRANS_STMT_ROLLBACK_TO, ast.TRANS_STMT_PREPARE,
			ast.TRANS_STMT_COMMIT_PREPARED, ast.TRANS_STMT_ROLLBACK_PREPARED:
			s.stop()
		}
	case *ast.VariableSetStmt:
		s.setting(v)
	}
}

// newReference is a foreign key the change created: the key it references,
// by the table as written, "" for a schema left to the search path, and
// the columns, nil for the primary key; and the table that owns it, as
// written, with the search path it was written under, and its name, ""
// when generated.
type newReference struct {
	schema, table string
	columns       []string

	owner     tableRef
	ownerPath int
	name      string
	retired   bool
}

// references records the keys new foreign keys of the owner reference.
func (s *scan) references(owner *ast.RangeVar, constraints ...*ast.Constraint) {
	for _, c := range constraints {
		if c.Contype == ast.CONSTR_FOREIGN && c.Pktable != nil {
			s.newReferences = append(s.newReferences, newReference{
				schema: s.schemaOf(c.Pktable), table: c.Pktable.Relname, columns: nameParts(c.PkAttrs),
				owner: tableRef{owner.Schemaname, owner.Relname}, ownerPath: s.pathVersion, name: c.Conname,
			})
		}
	}
}

// retireReferences records that a statement dropped the new foreign keys of
// a table, or the one of that name when name is set. Only a table named
// as the foreign key's owner was, under the same search path when neither
// is qualified, is known to be that owner.
func (s *scan) retireReferences(owner *ast.RangeVar, name string) {
	for i := range s.newReferences {
		r := &s.newReferences[i]
		same := r.owner.table == owner.Relname && r.owner.schema == owner.Schemaname &&
			(owner.Schemaname != "" || r.ownerPath == s.pathVersion)
		if same && (name == "" || r.name == name) {
			r.retired = true
		}
	}
}

// newlyReferenced reports whether a foreign key the change created may
// reference the table's key with these columns, nil when they cannot be
// read; primary tells whether the key is the primary key.
func (s *scan) newlyReferenced(t tableRef, columns []string, primary bool) bool {
	for _, r := range s.newReferences {
		if r.retired || r.table != t.table || r.schema != "" && r.schema != t.schema {
			continue
		}
		if columns == nil || r.columns == nil && primary || r.columns != nil && sameColumns(r.columns, columns) {
			return true
		}
	}
	return false
}

// newlyReferencedColumn reports whether a foreign key the change created
// may reference a key of the table that has the column; keyColumns are the
// table's primary key columns, nil when it has none or they cannot be read.
func (s *scan) newlyReferencedColumn(t tableRef, column string, keyColumns []string, hasKey bool) bool {
	for _, r := range s.newReferences {
		if r.retired || r.table != t.table || r.schema != "" && r.schema != t.schema {
			continue
		}
		if r.columns != nil && slices.Contains(r.columns, column) || r.columns == nil && hasKey && (keyColumns == nil || slices.Contains(keyColumns, column)) {
			return true
		}
	}
	return false
}

// newDependency is an object the change created that depends on a
// relation: the relation, in the schema it resolves to, the one written,
// or "" when neither is known, and the object as written, with the search
// path it was written under.
type newDependency struct {
	relation tableRef
	object   tableRef
	path     int
	retired  bool
}

// reads records the relations a new view's query names. An unqualified
// name of a CTE visible where it is used is the CTE, not a relation: a
// CTE is visible to the query of its WITH clause and to the CTEs after it,
// and, under WITH RECURSIVE, to every CTE of the clause.
func (s *scan) reads(view *ast.RangeVar, query ast.Node) {
	s.readsIn(view, query, nil)
}

func (s *scan) readsIn(view *ast.RangeVar, n ast.Node, visible map[string]bool) {
	ast.Inspect(n, func(n ast.Node) bool {
		if w, rest := splitWith(n); w != nil {
			scope := make(map[string]bool, len(visible))
			for name := range visible {
				scope[name] = true
			}
			var ctes []*ast.CommonTableExpr
			if w.Ctes != nil {
				for _, item := range w.Ctes.Items {
					if c, ok := item.(*ast.CommonTableExpr); ok {
						ctes = append(ctes, c)
					}
				}
			}
			if w.Recursive {
				for _, c := range ctes {
					scope[c.Ctename] = true
				}
			}
			for _, c := range ctes {
				s.readsIn(view, c.Ctequery, scope)
				scope[c.Ctename] = true
			}
			s.readsIn(view, rest, scope)
			return false
		}
		if rv, ok := n.(*ast.RangeVar); ok && !(rv.Schemaname == "" && visible[rv.Relname]) {
			s.newReads = append(s.newReads, newDependency{
				relation: tableRef{s.schemaOf(rv), rv.Relname},
				object:   tableRef{view.Schemaname, view.Relname},
				path:     s.pathVersion,
			})
		}
		return true
	})
}

// splitWith returns a statement's WITH clause and the statement without
// it, or nil when it has none.
func splitWith(n ast.Node) (*ast.WithClause, ast.Node) {
	switch v := n.(type) {
	case *ast.SelectStmt:
		if v.WithClause != nil {
			c := *v
			c.WithClause = nil
			return v.WithClause, &c
		}
	case *ast.InsertStmt:
		if v.WithClause != nil {
			c := *v
			c.WithClause = nil
			return v.WithClause, &c
		}
	case *ast.UpdateStmt:
		if v.WithClause != nil {
			c := *v
			c.WithClause = nil
			return v.WithClause, &c
		}
	case *ast.DeleteStmt:
		if v.WithClause != nil {
			c := *v
			c.WithClause = nil
			return v.WithClause, &c
		}
	case *ast.MergeStmt:
		if v.WithClause != nil {
			c := *v
			c.WithClause = nil
			return v.WithClause, &c
		}
	}
	return nil, nil
}

// returns records the relation whose row type a new function returns,
// when its result type names one.
func (s *scan) returns(v *ast.CreateFunctionStmt) {
	fn := nameParts(v.Funcname)
	if v.ReturnType == nil {
		if len(fn) > 0 {
			object := tableRef{table: fn[len(fn)-1]}
			if len(fn) >= 2 {
				object.schema = fn[len(fn)-2]
			}
			s.newFunctions[object]++
		}
		return
	}
	parts := nameParts(v.ReturnType.Names)
	if len(parts) == 0 || len(parts) > 2 || len(fn) == 0 {
		return
	}
	rv := &ast.RangeVar{Relname: parts[len(parts)-1]}
	if len(parts) == 2 {
		rv.Schemaname = parts[0]
	}
	object := tableRef{table: fn[len(fn)-1]}
	if len(fn) >= 2 {
		object.schema = fn[len(fn)-2]
	}
	s.newFunctions[object]++
	s.newReturns = append(s.newReturns, newDependency{relation: tableRef{s.schemaOf(rv), rv.Relname}, object: object, path: s.pathVersion})
}

// dependsOnNew reports whether an object the change created and has not
// dropped may depend on the relation.
func dependsOnNew(deps []newDependency, t tableRef) bool {
	for _, d := range deps {
		if !d.retired && d.relation.table == t.table && (d.relation.schema == "" || d.relation.schema == t.schema) {
			return true
		}
	}
	return false
}

// retire records that a statement dropped an object the change created,
// named as it was written, under the same search path when unqualified.
func (s *scan) retire(deps []newDependency, object tableRef) {
	for i := range deps {
		d := &deps[i]
		if d.object == object && (object.schema != "" || d.path == s.pathVersion) {
			d.retired = true
		}
	}
}

// create follows CREATE TABLE. A CREATE TABLE of a name its schema
// already holds is refused.
func (s *scan) create(st *statement, v *ast.CreateStmt) {
	exists := s.existsWhereCreated(v.Relation)
	switch {
	case s.schemaMissing(v.Relation) || exists && !v.IfNotExists:
		s.stop()
		return
	case exists:
		// IF NOT EXISTS of an existing table does nothing.
		return
	}
	if !s.nameUncertain(v.Relation) {
		s.createTable(st, v)
	}
	s.touch(v.Relation)
	s.generated = true
	if v.TableElts == nil {
		return
	}
	for _, item := range v.TableElts.Items {
		var constraints []*ast.Constraint
		switch e := item.(type) {
		case *ast.ColumnDef:
			constraints = constraintsOf(e.Constraints)
		case *ast.Constraint:
			constraints = []*ast.Constraint{e}
		}
		s.references(v.Relation, constraints...)
		for _, c := range constraints {
			// A key's index takes the constraint's name.
			if makesIndex(c) && c.Conname != "" {
				s.touchName(v.Relation.Schemaname, c.Conname)
			}
		}
	}
}

// createSchema follows CREATE SCHEMA and the statements it holds, each in
// the new schema unless it names another. A schema named after its owner
// takes the owner's name.
func (s *scan) createSchema(st *statement, v *ast.CreateSchemaStmt) {
	name := v.Schemaname
	if name == "" && v.Authrole != nil {
		name = v.Authrole.Rolename
	}
	if name == "" {
		if v.SchemaElts != nil && len(v.SchemaElts.Items) > 0 {
			s.stop()
		}
		return
	}
	s.schemas[name] = true
	if v.SchemaElts == nil {
		return
	}
	in := func(rv *ast.RangeVar) *ast.RangeVar {
		if rv == nil || rv.Schemaname != "" {
			return rv
		}
		c := *rv
		c.Schemaname = name
		return &c
	}
	for _, elt := range v.SchemaElts.Items {
		switch e := elt.(type) {
		case *ast.CreateStmt:
			c := *e
			c.Relation = in(e.Relation)
			s.create(st, &c)
		case *ast.ViewStmt:
			s.touch(in(e.View))
			s.reads(in(e.View), e.Query)
		case *ast.CreateSeqStmt:
			s.touch(in(e.Sequence))
		case *ast.IndexStmt:
			if e.Idxname != "" && e.Relation != nil {
				s.touchName(in(e.Relation).Schemaname, e.Idxname)
			}
		case *ast.CreateTrigStmt, *ast.GrantStmt:
		default:
			s.stop()
			return
		}
	}
}

// touch records a relation name a statement creates.
func (s *scan) touch(rv *ast.RangeVar) {
	if rv != nil {
		s.touchName(rv.Schemaname, rv.Relname)
	}
}

// touchName records a relation name a statement changed, in schema, or
// in whichever schema the search path finds when schema is "".
func (s *scan) touchName(schema, name string) {
	s.touched[[2]string{schema, name}] = true
	s.touchedName[name] = true
	for key := range s.freed {
		if key[1] == name && (schema == "" || key[0] == schema) {
			delete(s.freed, key)
		}
	}
}

// free records a name a resolved drop removed from its schema.
func (s *scan) free(schema, name string) {
	s.touchName(schema, name)
	s.freed[[2]string{schema, name}] = true
}

// unsettle records a table whose columns or constraints changed in a way
// the scan does not follow.
func (s *scan) unsettle(rv *ast.RangeVar) {
	s.unsettled[[2]string{s.schemaOf(rv), rv.Relname}] = true
}

// isUnsettled reports whether a statement may have changed the table's
// columns or constraints.
func (s *scan) isUnsettled(t tableRef) bool {
	return s.unsettled[[2]string{t.schema, t.table}] || s.unsettled[[2]string{"", t.table}]
}

// hasColumnIn reports whether a column set holds the table's column, in
// the table's schema or in an unknown one.
func hasColumnIn(set map[[3]string]bool, t tableRef, column string) bool {
	return set[[3]string{t.schema, t.table, column}] || set[[3]string{"", t.table, column}]
}

// isTouched reports whether a statement may have changed what the name
// means: a qualified name when a statement changed it in that schema or
// in an unknown one, an unqualified name when a statement changed it in
// any schema, which may come first on the search path.
func (s *scan) isTouched(rv *ast.RangeVar) bool {
	if rv.Schemaname == "" {
		return s.touchedName[rv.Relname]
	}
	return s.touched[[2]string{rv.Schemaname, rv.Relname}] || s.touched[[2]string{"", rv.Relname}]
}

// lookup resolves a relation name against the synced schema. ok is false
// when the name may no longer mean what the schema says, or names nothing
// the schema holds. An unqualified name needs a known search path, and a
// name beginning with pg_ may be a system catalog, which pg_catalog finds
// first.
func (s *scan) lookup(rv *ast.RangeVar) (schema string, kind relationKind, ok bool) {
	if s.index == nil || rv == nil || s.isTouched(rv) || s.schemas[rv.Schemaname] {
		return "", 0, false
	}
	if rv.Catalogname != "" && rv.Catalogname != s.index.database {
		return "", 0, false
	}
	if rv.Schemaname != "" {
		ns := s.index.schemas[rv.Schemaname]
		if ns == nil {
			return "", 0, false
		}
		kind, ok := ns.relations[rv.Relname]
		return rv.Schemaname, kind, ok
	}
	if strings.HasPrefix(rv.Relname, "pg_") {
		return "", 0, false
	}
	path, ok := s.searchPath()
	if !ok {
		return "", 0, false
	}
	for _, name := range path {
		if ns := s.index.schemas[name]; ns != nil {
			if kind, ok := ns.relations[rv.Relname]; ok {
				return name, kind, true
			}
			continue
		}
		// A system schema the snapshot leaves out may hold the name.
		if name == "information_schema" {
			return "", 0, false
		}
	}
	return "", 0, false
}

// searchPath expands "$user" in the current search path. Without a known
// current user, "$user" names no schema only when every schema of the
// database is listed by name.
func (s *scan) searchPath() ([]string, bool) {
	if s.path == nil {
		return nil, false
	}
	path := make([]string, 0, len(s.path))
	for _, name := range s.path {
		switch {
		case name != "$user":
			path = append(path, name)
		case s.user != "":
			path = append(path, s.user)
		case s.index.unlisted(s.path):
			return nil, false
		}
	}
	return path, true
}

// schemaOf returns the schema a relation name means now: the one the
// synced schema resolves it to, the one the SQL wrote, or "" when neither
// is known.
func (s *scan) schemaOf(rv *ast.RangeVar) string {
	if schema, _, ok := s.lookup(rv); ok {
		return schema
	}
	return rv.Schemaname
}

// table resolves the table an ALTER TABLE acts on, with its constraints
// as the synced schema lists them, when the change has not altered it in
// a way the scan does not follow. A partition is not a table here: its
// keys come from its parent.
func (s *scan) table(rv *ast.RangeVar) (tableRef, bool) {
	schema, kind, ok := s.lookup(rv)
	t := tableRef{schema, rv.Relname}
	if !ok || kind != kindTable || s.isUnsettled(t) {
		return tableRef{}, false
	}
	return t, true
}

// setting follows SET and RESET. The scan follows SET search_path to a
// list of plain names and SET ROLE. A transaction-local setting reverts at
// commit, and RESET returns to the session's default, which the target
// does not record, so after those the path or the user is unknown.
func (s *scan) setting(v *ast.VariableSetStmt) {
	if v.Kind == ast.VAR_RESET_ALL {
		s.path = nil
		s.pathVersion++
		return
	}
	switch v.Name {
	case "search_path":
		s.pathVersion++
		if v.Kind != ast.VAR_SET_VALUE || v.IsLocal || !plainNames(v.Args) {
			s.path = nil
			return
		}
		s.path = nil
		for _, item := range v.Args.Items {
			name, _ := constString(item)
			s.path = append(s.path, name)
		}
	case "role":
		s.pathVersion++
		switch {
		case v.IsLocal:
			s.user = ""
		case v.Kind == ast.VAR_SET_DEFAULT || v.Kind == ast.VAR_RESET:
			s.user = s.session
		case v.Kind == ast.VAR_SET_VALUE && v.Args != nil && len(v.Args.Items) == 1:
			name, _ := constString(v.Args.Items[0])
			if name == "none" {
				name = s.session
			}
			s.user = name
		default:
			s.user = ""
		}
	case "session_authorization":
		s.pathVersion++
		s.user, s.session = "", ""
	}
}

// plainNames reports whether every SET search_path value is a name the
// setting reads back unchanged: lower-case letters, digits, underscores,
// and dollar signs. A value such as 'MySchema' or 'a, b' is read as a list
// of identifiers, which the scan does not redo.
func plainNames(args *ast.List) bool {
	if args == nil {
		return false
	}
	for _, item := range args.Items {
		name, ok := constString(item)
		if !ok || name == "" {
			return false
		}
		for i := 0; i < len(name); i++ {
			c := name[i]
			if !(c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '_' || c == '$') {
				return false
			}
		}
	}
	return true
}

// data follows a statement that changes data, not the schema. A call of
// set_config on search_path changes how later names resolve: the form
// pg_dump writes, a lone SELECT of it with constant arguments, sets the
// path, and any other call leaves it unknown.
func (s *scan) data(n ast.Node) {
	calls := searchPathCalls(n)
	if len(calls) == 0 {
		return
	}
	s.pathVersion++
	s.path = nil
	if path, ok := dumpSearchPath(n, calls); ok {
		// An empty path is a known path that finds nothing.
		s.path = append([]string{}, path...)
	}
}

// searchPathCalls lists the set_config calls that may set search_path:
// every call whose setting name is not a constant naming another setting.
func searchPathCalls(n ast.Node) []*ast.FuncCall {
	var calls []*ast.FuncCall
	ast.Inspect(n, func(n ast.Node) bool {
		fc, ok := n.(*ast.FuncCall)
		if !ok {
			return true
		}
		parts := nameParts(fc.Funcname)
		builtin := len(parts) == 1 || len(parts) == 2 && parts[0] == "pg_catalog"
		if !builtin || parts[len(parts)-1] != "set_config" {
			return true
		}
		if name, ok := constString(firstArg(fc)); ok && !strings.EqualFold(name, "search_path") {
			return true
		}
		calls = append(calls, fc)
		return true
	})
	return calls
}

// dumpSearchPath reads SELECT pg_catalog.set_config('search_path', '...',
// false), the only statement in n, and returns the path it sets.
func dumpSearchPath(n ast.Node, calls []*ast.FuncCall) ([]string, bool) {
	sel, ok := n.(*ast.SelectStmt)
	if !ok || len(calls) != 1 || sel.FromClause != nil || sel.WhereClause != nil || sel.WithClause != nil ||
		sel.TargetList == nil || len(sel.TargetList.Items) != 1 {
		return nil, false
	}
	target, ok := sel.TargetList.Items[0].(*ast.ResTarget)
	if !ok || target.Val != calls[0] || calls[0].Args == nil || len(calls[0].Args.Items) != 3 {
		return nil, false
	}
	value, ok := constString(calls[0].Args.Items[1])
	if !ok {
		return nil, false
	}
	local, ok := calls[0].Args.Items[2].(*ast.A_Const)
	if !ok || local.Isnull {
		return nil, false
	}
	if b, ok := local.Val.(*ast.Boolean); !ok || b.Boolval {
		return nil, false
	}
	return parseSearchPath(value), true
}

func firstArg(fc *ast.FuncCall) ast.Node {
	if fc.Args == nil || len(fc.Args.Items) == 0 {
		return nil
	}
	return fc.Args.Items[0]
}

// constString returns the value of a string constant.
func constString(n ast.Node) (string, bool) {
	c, ok := n.(*ast.A_Const)
	if !ok || c.Isnull {
		return "", false
	}
	str, ok := c.Val.(*ast.String)
	if !ok {
		return "", false
	}
	return str.Str, true
}

// selectsInto reports whether a branch of a set operation carries INTO.
func selectsInto(s *ast.SelectStmt) bool {
	if s == nil {
		return false
	}
	return s.IntoClause != nil || selectsInto(s.Larg) || selectsInto(s.Rarg)
}

// isTemp reports whether a relation is created in the temporary schema.
func isTemp(rv *ast.RangeVar) bool {
	return rv != nil && (rv.Relpersistence == 't' || rv.Schemaname == "pg_temp" || strings.HasPrefix(rv.Schemaname, "pg_temp_"))
}

// isRelationKind reports whether an object type lives in the relation
// namespace.
func isRelationKind(t ast.ObjectType) bool {
	switch t {
	case ast.OBJECT_TABLE, ast.OBJECT_VIEW, ast.OBJECT_MATVIEW, ast.OBJECT_FOREIGN_TABLE, ast.OBJECT_SEQUENCE, ast.OBJECT_INDEX:
		return true
	}
	return false
}

// rename follows RENAME. A relation renamed is touched under both names,
// and a pending table it may be keeps the new name too. An index rename
// renames the constraint the index backs, if any.
func (s *scan) rename(v *ast.RenameStmt) {
	switch {
	case isRelationKind(v.RenameType) && v.Relation != nil:
		s.touchName(v.Relation.Schemaname, v.Relation.Relname)
		s.touchName(v.Relation.Schemaname, v.Newname)
		for _, p := range s.matching(v.Relation) {
			p.names[v.Newname] = true
		}
		if v.RenameType == ast.OBJECT_INDEX {
			s.renamedKeys[v.Relation.Relname] = true
			s.renamedKeys[v.Newname] = true
		}
	case v.RenameType == ast.OBJECT_TYPE:
		if parts := nameParts(listOf(v.Object)); len(parts) > 0 {
			schema := ""
			if len(parts) >= 2 {
				schema = parts[len(parts)-2]
			}
			s.touchName(schema, parts[len(parts)-1])
			s.touchName(schema, v.Newname)
		}
	case v.RenameType == ast.OBJECT_TABCONSTRAINT && v.Relation != nil:
		schema := s.schemaOf(v.Relation)
		s.constraints[[3]string{schema, v.Relation.Relname, v.Subname}] = true
		s.constraints[[3]string{schema, v.Relation.Relname, v.Newname}] = true
	case v.RenameType == ast.OBJECT_COLUMN && v.Relation != nil:
		schema := s.schemaOf(v.Relation)
		s.renamedFrom[[3]string{schema, v.Relation.Relname, v.Subname}] = true
		s.renamedTo[[3]string{schema, v.Relation.Relname, v.Newname}] = true
	case v.RenameType == ast.OBJECT_SCHEMA:
		s.stop()
	}
}

// setSchema follows SET SCHEMA. A pending table it may be can then be in
// any schema.
func (s *scan) setSchema(v *ast.AlterObjectSchemaStmt) {
	switch {
	case isRelationKind(v.ObjectType) && v.Relation != nil:
		s.touch(v.Relation)
		s.touchName(v.Newschema, v.Relation.Relname)
		for _, p := range s.matching(v.Relation) {
			// The move is known only when the statement names the table's
			// own schema; otherwise it may be another table of that name.
			if v.Relation.Schemaname != "" && v.Relation.Schemaname == p.schema {
				p.schema = v.Newschema
			} else {
				p.schema = ""
			}
		}
	case v.ObjectType == ast.OBJECT_TYPE:
		if parts := nameParts(listOf(v.Object)); len(parts) > 0 {
			s.touchName("", parts[len(parts)-1])
		}
	case v.ObjectType == ast.OBJECT_EXTENSION:
		s.stop()
	}
}

// drop follows DROP. A cascading drop of a relation also takes the
// relations that go with it (its partitions, indexes, owned sequences,
// and the views that read it), the foreign keys that reference it, and
// any inheritance child, which the synced schema does not record; a
// cascading drop of anything else may take columns of any table, which
// ends the scan. DROP EXTENSION drops the extension's objects.
func (s *scan) drop(v *ast.DropStmt) {
	kind := ast.ObjectType(v.RemoveType)
	cascade := v.Behavior == int(ast.DROP_CASCADE)
	if kind == ast.OBJECT_EXTENSION || cascade && !isRelationKind(kind) {
		s.stop()
		return
	}
	if kind == ast.OBJECT_FUNCTION || kind == ast.OBJECT_ROUTINE {
		s.dropFunctions(v)
		return
	}
	if kind == ast.OBJECT_SCHEMA && v.Objects != nil {
		// Without CASCADE the schema must be empty; its name may then be
		// created again.
		for _, obj := range v.Objects.Items {
			if parts := nameParts(listOf(obj)); len(parts) == 1 {
				s.schemas[parts[0]] = true
			}
		}
		return
	}
	if !isRelationKind(kind) && kind != ast.OBJECT_TYPE || v.Objects == nil {
		return
	}
	if !cascade && s.dropRefused(v) || !v.Missing_ok && s.dropsMissing(v) || s.dropsWrongKind(v) {
		s.stop()
		return
	}
	for _, obj := range v.Objects.Items {
		parts := nameParts(listOf(obj))
		if len(parts) == 0 {
			continue
		}
		rv := &ast.RangeVar{Relname: parts[len(parts)-1]}
		if len(parts) >= 2 {
			rv.Schemaname = parts[len(parts)-2]
		}
		schema, relKind, resolved := s.lookup(rv)
		switch {
		case resolved:
			s.free(schema, rv.Relname)
		case v.Missing_ok && s.index != nil && s.missing(rv):
			// IF EXISTS of a name the target lacks drops nothing.
			continue
		default:
			s.touch(rv)
			schema = ""
		}
		if kind == ast.OBJECT_INDEX {
			s.renamedKeys[rv.Relname] = true
		}
		for _, p := range s.matching(rv) {
			p.settled = true
		}
		switch kind {
		case ast.OBJECT_TABLE:
			s.retireReferences(rv, "")
		case ast.OBJECT_VIEW, ast.OBJECT_MATVIEW:
			s.retire(s.newReads, tableRef{rv.Schemaname, rv.Relname})
		}
		s.dropOwned(rv, schema, resolved && relKind == kindTable)
		if cascade {
			s.dropDependents(rv, schema)
		}
	}
	if cascade {
		// An inheritance child of a dropped table goes with it.
		for _, p := range s.pending {
			if p.inherits || p.existed {
				p.settled = true
			}
		}
	}
}

// dropFunctions follows DROP FUNCTION: a synced function whose name the
// scan resolves to one function, and a function the change created, no
// longer hold the relation whose row type they return.
func (s *scan) dropFunctions(v *ast.DropStmt) {
	if v.Objects == nil {
		return
	}
	for _, obj := range v.Objects.Items {
		owa, ok := obj.(*ast.ObjectWithArgs)
		if !ok {
			continue
		}
		parts := nameParts(owa.Objname)
		if len(parts) == 0 || len(parts) > 2 {
			continue
		}
		fn := tableRef{table: parts[len(parts)-1]}
		if len(parts) == 2 {
			fn.schema = parts[0]
		}
		// An overload the change created is retired only when the name
		// leaves no choice of which function was dropped.
		synced := 0
		if s.index != nil {
			for name, n := range s.index.functions {
				if name.table == fn.table && (fn.schema == "" || name.schema == fn.schema) {
					synced += n
				}
			}
		}
		if s.newFunctions[fn] == 1 && synced == 0 {
			s.retire(s.newReturns, fn)
		}
		if s.index == nil || s.newFunctions[fn] > 0 {
			continue
		}
		if fn.schema != "" {
			if s.index.functions[fn] == 1 {
				s.droppedFunctions[fn] = true
			}
			continue
		}
		path, ok := s.searchPath()
		if !ok {
			continue
		}
		// The server finds the first schema on the path with a function of
		// that name; one function there is the one dropped.
		for _, name := range path {
			if n := s.index.functions[tableRef{name, fn.table}]; n > 0 {
				if n == 1 {
					s.droppedFunctions[tableRef{name, fn.table}] = true
				}
				break
			}
		}
	}
}

// dropsWrongKind reports whether a DROP names a relation of another kind,
// which the server refuses: DROP TABLE of a view, DROP VIEW of a table, or
// DROP INDEX of the index a constraint owns.
func (s *scan) dropsWrongKind(v *ast.DropStmt) bool {
	kind := ast.ObjectType(v.RemoveType)
	if s.index == nil || v.Objects == nil || !isRelationKind(kind) {
		return false
	}
	for _, obj := range v.Objects.Items {
		parts := nameParts(listOf(obj))
		if len(parts) == 0 || len(parts) > 2 {
			continue
		}
		rv := &ast.RangeVar{Relname: parts[len(parts)-1]}
		if len(parts) == 2 {
			rv.Schemaname = parts[0]
		}
		schema, relKind, ok := s.lookup(rv)
		if !ok || relKind == kindAmbiguous {
			continue
		}
		if !dropKindMatches(kind, relKind) {
			return true
		}
		if kind == ast.OBJECT_INDEX && s.index.constraintIndexes[tableRef{schema, rv.Relname}] {
			return true
		}
	}
	return false
}

// viewHolds reports whether a synced view may still read what the snapshot
// says it reads: the change neither dropped it nor redefined it.
func (s *scan) viewHolds(view tableRef) bool {
	return !s.freed[[2]string{view.schema, view.table}] && !s.replacedViews[view]
}

// dropKindMatches reports whether a DROP of an object type can drop a
// relation of that kind.
func dropKindMatches(t ast.ObjectType, kind relationKind) bool {
	switch t {
	case ast.OBJECT_TABLE:
		return kind == kindTable || kind == kindPartition
	case ast.OBJECT_VIEW:
		return kind == kindView
	case ast.OBJECT_MATVIEW:
		return kind == kindMatView
	case ast.OBJECT_FOREIGN_TABLE:
		return kind == kindForeignTable
	case ast.OBJECT_SEQUENCE:
		return kind == kindSequence
	case ast.OBJECT_INDEX:
		return kind == kindIndex
	}
	return true
}

// schemaMissing reports whether a qualified name's schema certainly does
// not exist: the target lacks it and the change did not create it.
func (s *scan) schemaMissing(rv *ast.RangeVar) bool {
	return s.index != nil && rv != nil && rv.Schemaname != "" && !isTemp(rv) &&
		s.index.schemas[rv.Schemaname] == nil && !s.schemas[rv.Schemaname]
}

// dropsMissing reports whether a DROP without IF EXISTS names a relation
// the target certainly lacks: no statement of the change made the name,
// and the schema the SQL wrote, or every schema of a known search path,
// holds no relation of that name.
func (s *scan) dropsMissing(v *ast.DropStmt) bool {
	if s.index == nil || v.Objects == nil || !isRelationKind(ast.ObjectType(v.RemoveType)) {
		return false
	}
	for _, obj := range v.Objects.Items {
		parts := nameParts(listOf(obj))
		if len(parts) == 0 || len(parts) > 2 {
			continue
		}
		rv := &ast.RangeVar{Relname: parts[len(parts)-1]}
		if len(parts) == 2 {
			rv.Schemaname = parts[0]
		}
		if s.missing(rv) {
			return true
		}
	}
	return false
}

// missing reports whether the target certainly has no relation of that
// name where the name resolves.
func (s *scan) missing(rv *ast.RangeVar) bool {
	if s.isTouched(rv) {
		return false
	}
	if rv.Schemaname != "" {
		ns := s.index.schemas[rv.Schemaname]
		if ns == nil || s.schemas[rv.Schemaname] {
			return false
		}
		_, ok := ns.relations[rv.Relname]
		return !ok
	}
	if strings.HasPrefix(rv.Relname, "pg_") {
		return false
	}
	path, ok := s.searchPath()
	if !ok {
		return false
	}
	for _, name := range path {
		if s.schemas[name] || name == "information_schema" {
			return false
		}
		if ns := s.index.schemas[name]; ns != nil {
			if _, ok := ns.relations[rv.Relname]; ok {
				return false
			}
		}
	}
	return true
}

// dropRefused reports whether the server refuses a DROP without CASCADE of
// relations the scan resolves: a view the statement does not drop reads
// one, a function returns its row type, or, for a table, a foreign key of
// a table the statement does not drop references it, as the synced schema
// or the change itself records.
func (s *scan) dropRefused(v *ast.DropStmt) bool {
	kind := ast.ObjectType(v.RemoveType)
	if s.index == nil || kind != ast.OBJECT_TABLE && kind != ast.OBJECT_VIEW && kind != ast.OBJECT_MATVIEW && kind != ast.OBJECT_FOREIGN_TABLE {
		return false
	}
	var dropped []tableRef
	for _, obj := range v.Objects.Items {
		parts := nameParts(listOf(obj))
		if len(parts) == 0 {
			continue
		}
		rv := &ast.RangeVar{Relname: parts[len(parts)-1]}
		if len(parts) >= 2 {
			rv.Schemaname = parts[len(parts)-2]
		}
		if schema, _, ok := s.lookup(rv); ok {
			dropped = append(dropped, tableRef{schema, rv.Relname})
		}
	}
	for _, t := range dropped {
		if dependsOnNew(s.newReads, t) || dependsOnNew(s.newReturns, t) {
			return true
		}
		for _, fn := range s.index.returnedBy[t] {
			if !s.droppedFunctions[fn] {
				return true
			}
		}
		for _, view := range s.index.readers[t] {
			if s.viewHolds(view) && !slices.Contains(dropped, view) {
				return true
			}
		}
		if kind != ast.OBJECT_TABLE {
			continue
		}
		if s.newlyReferenced(t, nil, true) {
			return true
		}
		for _, fk := range s.index.foreignKeys {
			if fk.referenced == t && !fk.in(s.dropped) && !slices.Contains(dropped, fk.owner) {
				return true
			}
		}
	}
	return false
}

// dropOwned records what any drop of a relation takes with it, CASCADE or
// not: its partitions, indexes, and owned sequences, freed in its schema
// when the scan resolved it and touched in every schema the name may be in
// otherwise, and, for a table the scan resolved, its own foreign keys,
// which then no longer hold a key of another table.
func (s *scan) dropOwned(rv *ast.RangeVar, schema string, table bool) {
	if s.index == nil {
		return
	}
	for name, ns := range s.index.schemas {
		switch {
		case schema != "" && name == schema:
			for _, d := range ns.dependents[rv.Relname] {
				if s.reowned[d] {
					s.touchName(name, d)
				} else {
					s.free(name, d)
				}
			}
		case schema == "" && (rv.Schemaname == "" || rv.Schemaname == name):
			for _, d := range ns.dependents[rv.Relname] {
				s.touchName(name, d)
			}
		}
	}
	if !table {
		return
	}
	for _, fk := range s.index.foreignKeys {
		if fk.owner == (tableRef{schema, rv.Relname}) {
			s.dropped[[3]string{schema, rv.Relname, fk.name}] = true
		}
	}
}

// dropDependents records what a cascading drop of a relation takes with
// it besides what it owns, as far as the synced schema records it: the
// views that read it, freed when the scan resolved the relation to schema
// and touched in every schema the name may be in otherwise, and the
// foreign keys that reference it.
func (s *scan) dropDependents(rv *ast.RangeVar, schema string) {
	s.dropReferences(rv, tableRef{schema, rv.Relname}, schema != "", nil)
	if s.index == nil {
		return
	}
	var start []tableRef
	if schema != "" {
		start = []tableRef{{schema, rv.Relname}}
	} else {
		for name := range s.index.schemas {
			if rv.Schemaname == "" || rv.Schemaname == name {
				start = append(start, tableRef{name, rv.Relname})
			}
		}
	}
	// The views reading a dropped view go too.
	seen := make(map[tableRef]bool)
	for len(start) > 0 {
		t := start[0]
		start = start[1:]
		for _, view := range s.index.readers[t] {
			if seen[view] {
				continue
			}
			seen[view] = true
			if schema != "" {
				s.free(view.schema, view.table)
			} else {
				s.touchName(view.schema, view.table)
			}
			start = append(start, view)
		}
	}
}

// dropReferences records that foreign keys a cascading drop takes may be
// gone: those referencing t when the scan resolved it (known), and those
// referencing a table of that name, in the schema the SQL wrote if any,
// otherwise. affected, when set, narrows them to the ones that depend on
// what was dropped.
func (s *scan) dropReferences(rv *ast.RangeVar, t tableRef, known bool, affected func(foreignKeyRef) bool) {
	if s.index == nil {
		return
	}
	for _, fk := range s.index.foreignKeys {
		switch {
		case known && fk.referenced != t:
			continue
		case !known && (fk.referenced.table != rv.Relname || rv.Schemaname != "" && fk.referenced.schema != rv.Schemaname):
			continue
		case affected != nil && !affected(fk):
			continue
		}
		s.constraints[[3]string{fk.owner.schema, fk.owner.table, fk.name}] = true
	}
}

// constraintKnown reports whether a table's constraint of that name is
// still the one the synced schema lists.
func (s *scan) constraintKnown(t tableRef, name string) bool {
	return !s.constraints[[3]string{t.schema, t.table, name}] && !s.constraints[[3]string{"", t.table, name}] && !s.renamedKeys[name]
}
