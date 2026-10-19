package review

import (
	"strings"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// statement lowers one statement and reports whether the walk follows it.
// A statement it does not follow ends the walk: what it does is not
// bounded by anything the walk knows. The SQLSTATE beside a failure is
// the one PostgreSQL raises for it.
func (b *builder) statement(n ast.Node) bool {
	switch v := n.(type) {
	case *ast.CreateSchemaStmt:
		return b.createSchema(v)
	case *ast.CreateStmt:
		return b.createTable(v)
	case *ast.IndexStmt:
		return b.createIndex(v)
	case *ast.ViewStmt:
		return b.createView(v)
	case *ast.CreateTableAsStmt:
		return b.createTableAs(v.Into, ast.ObjectType(v.Objtype) == ast.OBJECT_MATVIEW, v.IfNotExists, v.Query)
	case *ast.CreateSeqStmt:
		return b.createSequence(v)
	case *ast.CompositeTypeStmt:
		return b.createCompositeType(v)
	case *ast.CreateEnumStmt:
		return b.createType(nameParts(v.TypeName))
	case *ast.CreateDomainStmt:
		return b.createType(nameParts(v.Domainname))
	case *ast.AlterTableStmt:
		return b.alterTable(v)
	case *ast.RenameStmt:
		return b.rename(v)
	case *ast.AlterObjectSchemaStmt:
		return b.setSchema(v)
	case *ast.DropStmt:
		return b.drop(v)
	case *ast.CommentStmt:
		return b.comment(v)
	case *ast.GrantStmt:
		return b.grant(v)
	case *ast.TruncateStmt:
		return b.truncate(v)
	case *ast.LockStmt:
		return b.lock(v)
	case *ast.AlterSeqStmt:
		b.exists(v.Sequence, v.MissingOk)
		return true
	case *ast.RefreshMatViewStmt:
		b.existsOfKind(v.Relation, relMatview, "materialized view")
		return true
	case *ast.CreatePolicyStmt:
		b.exists(v.Table, false)
		return true
	case *ast.InsertStmt:
		return b.dml(v.Relation, v)
	case *ast.UpdateStmt:
		return b.dml(v.Relation, v)
	case *ast.DeleteStmt:
		return b.dml(v.Relation, v)
	case *ast.MergeStmt:
		return b.dml(v.Relation, v)
	case *ast.SelectStmt:
		if v.IntoClause != nil {
			return b.createTableAs(v.IntoClause, false, false, v)
		}
		if calls := searchPathCalls(v); len(calls) > 0 {
			// pg_dump's SELECT pg_catalog.set_config('search_path', '',
			// false) sets a path the walk reads; any other call of it
			// leaves the path unknown.
			if path, ok := dumpSearchPath(v, calls); ok {
				known := append([]string{}, path...)
				b.session(func() { b.w.path = known })
				return true
			}
			b.session(func() { b.w.path = nil })
		}
		b.last = b.w.callsUserCode(v)
		return true
	case *ast.VariableSetStmt:
		return b.setVariable(v)
	case *ast.TransactionStmt:
		switch v.Kind {
		case ast.TRANS_STMT_BEGIN, ast.TRANS_STMT_START, ast.TRANS_STMT_COMMIT, ast.TRANS_STMT_SAVEPOINT, ast.TRANS_STMT_RELEASE:
			return true
		}
		// A rollback undoes statements the walk took as done.
		return false
	case *ast.CreateFunctionStmt:
		parts := nameParts(v.Funcname)
		if len(parts) > 0 {
			name := parts[len(parts)-1]
			b.session(func() { b.w.routines[name] = true })
		}
		return true
	case *ast.CreateTrigStmt:
		return b.createTrigger(v)
	case *ast.VariableShowStmt, *ast.NotifyStmt, *ast.ListenStmt, *ast.UnlistenStmt,
		*ast.CheckPointStmt, *ast.VacuumStmt, *ast.ReindexStmt, *ast.ClusterStmt,
		*ast.AlterOwnerStmt, *ast.AlterFunctionStmt, *ast.AlterDefaultPrivilegesStmt,
		*ast.CreateRoleStmt, *ast.AlterRoleStmt, *ast.AlterRoleSetStmt, *ast.DropRoleStmt,
		*ast.GrantRoleStmt, *ast.AlterPolicyStmt:
		// None changes a name the walk keeps.
		return true
	}
	return false
}

// missing is the message for a relation that resolves nowhere.
func missing(r resolved, written string) string {
	if r.missingSchema != "" {
		return "schema " + ident(r.missingSchema) + " does not exist"
	}
	return "relation " + written + " does not exist"
}

// exists checks that a relation the statement names exists (42P01; a
// schema that does not exist is 3F000).
func (b *builder) exists(rv *ast.RangeVar, missingOK bool) resolved {
	r := b.rangeVar(rv)
	if r.fact == walk.Absent && !missingOK {
		b.failAt(rangeOf(rv.Loc), r.keys, missing(r, relation(rv)))
	}
	return r
}

// opened checks that the relation a query opens exists. The parser opens
// it, and names the relation, not its schema, when the schema is missing
// (42P01).
//
// pg: src/backend/parser/parse_relation.c — parserOpenTable
func (b *builder) opened(rv *ast.RangeVar) resolved {
	r := b.rangeVar(rv)
	if r.fact == walk.Absent {
		b.failAt(rangeOf(rv.Loc), r.keys, "relation "+relation(rv)+" does not exist")
	}
	return r
}

// hasColumns reports whether the walk keeps the columns of a relation of
// the kind: a table's and a foreign table's. Other relations' are not
// followed, and a statement that names a column of one is judged on
// nothing about the column.
func hasColumns(kind string) bool {
	return kind == relTable || kind == relForeign
}

// existsOfKind checks that a relation exists and is of the kind (42809).
func (b *builder) existsOfKind(rv *ast.RangeVar, kind, word string) {
	if r := b.exists(rv, false); r.fact == walk.Present {
		b.need(r.key, walk.Present, []string{kind}, &walkNote{rng: rangeOf(rv.Loc), present: relation(rv) + " is not a " + word})
	}
}

// roleName returns the role a RoleSpec names, or "" when unknown.
func (b *builder) roleName(rs *ast.RoleSpec) string {
	if rs == nil {
		return ""
	}
	switch ast.RoleSpecType(rs.Roletype) {
	case ast.ROLESPEC_CSTRING:
		return rs.Rolename
	case ast.ROLESPEC_CURRENT_ROLE, ast.ROLESPEC_CURRENT_USER:
		return b.w.user
	case ast.ROLESPEC_SESSION_USER:
		return b.w.session
	}
	return ""
}

// createSchema follows CREATE SCHEMA without schema elements; PostgreSQL
// may reorder the elements, which the walk does not follow.
//
// pg: src/backend/commands/schemacmds.c — CreateSchemaCommand (the
// reserved-name check comes before IF NOT EXISTS)
func (b *builder) createSchema(v *ast.CreateSchemaStmt) bool {
	if v.SchemaElts != nil && len(v.SchemaElts.Items) > 0 {
		return false
	}
	written := v.Schemaname
	if written == "" {
		written = b.roleName(v.Authrole)
	}
	name, ok := walkName(written)
	if !ok {
		b.mayAdd(walk.Scope{Space: spaceSchema})
		return true
	}
	if strings.HasPrefix(name, "pg_") {
		b.fail(nil, walk.Present, "unacceptable schema name "+ident(written)) // 42939
		return true
	}
	k := schemaKey(name)
	prior, _ := b.l.Lookup(k)
	if !v.IfNotExists {
		b.need(k, walk.Absent, nil, &walkNote{rng: b.anchor(), present: "schema " + ident(written) + " already exists"}) // 42P06
		prior = walk.Absent
	}
	switch prior {
	case walk.Absent:
		b.set(k, walk.Present, "")
		b.clear(walk.Scope{Parent: walk.Path(name), Deep: true})
		b.made(k)
	case walk.Unknown:
		b.set(k, walk.Present, "")
	}
	return true
}

// made records an object the change created, which a later DROP SCHEMA
// without CASCADE finds.
func (b *builder) made(k walk.Key) {
	b.session(func() { b.w.made = append(b.w.made, k) })
}

// placeUnknown records a new object whose schema the walk cannot tell:
// the name may now exist in any schema, with anything under it.
func (b *builder) placeUnknown(name string, generated bool) {
	if n, ok := walkName(name); ok {
		b.mayAdd(walk.Scope{Space: spaceRelation, Deep: true, Prefix: n})
		b.mayAdd(walk.Scope{Space: spaceType, Deep: true, Prefix: n})
		if generated {
			b.mayAdd(walk.Scope{Space: spaceRelation, Deep: true, Prefix: generatedPrefix(n)})
		}
		return
	}
	b.mayAdd(walk.Scope{Space: spaceRelation, Deep: true})
	b.mayAdd(walk.Scope{Space: spaceType, Deep: true})
}

// isSerial reports whether a column type is one of the serial pseudo
// types, which create a sequence under a generated name.
func isSerial(tn *ast.TypeName) bool {
	if tn == nil {
		return false
	}
	parts := nameParts(tn.Names)
	if len(parts) != 1 {
		return false
	}
	switch parts[0] {
	case "serial", "serial4", "bigserial", "serial8", "smallserial", "serial2":
		return true
	}
	return false
}

// createTable follows CREATE TABLE. A plain column list gives the new
// table exactly those columns; LIKE, INHERITS, OF, and PARTITION OF copy
// columns and constraints the walk does not derive.
//
// pg: src/backend/parser/parse_utilcmd.c — transformCreateStmt (IF NOT
// EXISTS is judged before the columns), transformColumnDefinition
func (b *builder) createTable(v *ast.CreateStmt) bool {
	rv := v.Relation
	if rv == nil {
		return false
	}
	rng := rangeOf(rv.Loc)
	temp := rv.Relpersistence == ast.RELPERSISTENCE_TEMP
	var cols []*ast.ColumnDef
	var cons []*ast.Constraint
	var likes []*ast.TableLikeClause
	if v.TableElts != nil {
		for _, item := range v.TableElts.Items {
			switch e := item.(type) {
			case *ast.ColumnDef:
				cols = append(cols, e)
				if e.Constraints != nil {
					for _, c := range e.Constraints.Items {
						if con, ok := c.(*ast.Constraint); ok {
							cons = append(cons, con)
						}
					}
				}
			case *ast.Constraint:
				cons = append(cons, e)
			case *ast.TableLikeClause:
				likes = append(likes, e)
			default:
				return false
			}
		}
	}
	if v.Constraints != nil {
		for _, c := range v.Constraints.Items {
			if con, ok := c.(*ast.Constraint); ok {
				cons = append(cons, con)
			}
		}
	}
	var parents []*ast.RangeVar
	if v.InhRelations != nil {
		for _, item := range v.InhRelations.Items {
			if p, ok := item.(*ast.RangeVar); ok {
				parents = append(parents, p)
			}
		}
	}
	partitionOf := v.Partbound != nil
	plain := len(likes) == 0 && len(parents) == 0 && v.OfTypename == nil && !partitionOf
	generated := !plain
	for _, c := range cols {
		if isSerial(c.TypeName) {
			generated = true
		}
	}
	unnamedIndexes := false
	for _, con := range cons {
		switch con.Contype {
		case ast.CONSTR_PRIMARY, ast.CONSTR_UNIQUE, ast.CONSTR_EXCLUSION:
			if con.Conname == "" {
				unnamedIndexes = true
			}
		case ast.CONSTR_IDENTITY:
			generated = true
		}
	}
	generated = generated || unnamedIndexes

	c := b.creationSchema(rv.Schemaname, temp)
	if c.fail != "" {
		b.failAt(rng, c.keys, c.fail) // 3F000
		return true
	}
	name, ok := walkName(rv.Relname)
	if c.fact != walk.Present || !ok {
		b.placeUnknown(rv.Relname, true)
		return true
	}
	k := relKey(c.schema, name)
	tk := typeKey(c.schema, name)
	prior, _ := b.l.Lookup(k)
	switch {
	case !v.IfNotExists:
		b.need(k, walk.Absent, nil, &walkNote{rng: rng, present: "relation " + relation(rv) + " already exists"}) // 42P07
		b.need(tk, walk.Absent, nil, &walkNote{rng: rng, present: "type " + relation(rv) + " already exists"})    // 42710
		prior = walk.Absent
	case prior == walk.Present:
		return true // a notice: nothing is created
	case prior == walk.Absent:
		b.need(tk, walk.Absent, nil, &walkNote{rng: rng, present: "type " + relation(rv) + " already exists"})
	}

	if prior == walk.Absent {
		// The relations the new table names must exist (42P01): the ones
		// it copies from and the tables its foreign keys reference, but
		// itself. IF NOT EXISTS on a table that exists looks at none.
		for _, l := range likes {
			b.exists(l.Relation, false)
		}
		for _, p := range parents {
			b.exists(p, false)
		}
		for _, con := range cons {
			if con.Contype == ast.CONSTR_FOREIGN && con.Pktable != nil && !sameRelation(con.Pktable, rv) {
				b.exists(con.Pktable, false)
			}
		}
		// Within the statement, judged from the text alone.
		if !partitionOf && v.OfTypename == nil {
			seen := make(map[string]bool, len(cols))
			for _, col := range cols {
				if col.Colname != "" && seen[col.Colname] {
					b.failAt(rangeOf(col.Loc), nil, "column "+ident(col.Colname)+" specified more than once") // 42701
					return true
				}
				seen[col.Colname] = true
			}
		}
		primaries := 0
		for _, con := range cons {
			if con.Contype == ast.CONSTR_PRIMARY {
				primaries++
				if primaries > 1 {
					b.failAt(rangeOf(con.Loc), nil, "multiple primary keys for table "+relation(rv)+" are not allowed") // 42P16
					return true
				}
			}
		}
	}

	parent := tableOf(k)
	switch prior {
	case walk.Absent:
		b.set(k, walk.Present, relTable)
		b.set(tk, walk.Present, "")
		b.clear(walk.Scope{Parent: parent, Deep: true})
		if !plain {
			b.forget(walk.Scope{Space: spaceColumn, Parent: parent})
			b.forget(walk.Scope{Space: spaceConstraint, Parent: parent})
			b.set(primaryKeyOf(k), walk.Unknown, "")
		}
		for _, col := range cols {
			if n, ok := walkName(col.Colname); ok {
				b.set(columnKey(k, n), walk.Present, kindColumn)
			} else {
				b.mayAdd(walk.Scope{Space: spaceColumn, Parent: parent})
			}
		}
		if plain {
			b.set(primaryKeyOf(k), walk.Absent, "")
		}
		for _, con := range cons {
			b.constraint(k, rv, con, false, false)
		}
		b.made(k)
		b.made(tk)
	case walk.Unknown:
		// It existed already or was created: either way a relation is
		// there, of a kind and with contents unknown. One that existed
		// may have no row type: a sequence or an index.
		b.set(k, walk.Present, "")
		b.set(tk, walk.Unknown, "")
		b.forget(walk.Scope{Parent: parent, Deep: true})
	}
	if generated {
		b.mayAdd(walk.Scope{Space: spaceRelation, Parent: walk.Path(c.schema), Prefix: generatedPrefix(name)})
	}
	if temp || c.schema == tempSchema {
		b.session(func() { b.w.temporary = true })
	}
	if v.Partspec != nil {
		b.session(func() {
			b.w.partitioned[k] = true
			b.w.anyPartitioned = true
		})
	}
	if partitionOf {
		var parentKey walk.Key
		if len(parents) == 1 {
			if r := b.rangeVar(parents[0]); r.fact == walk.Present {
				parentKey = r.key
			}
		}
		if parentKey == (walk.Key{}) {
			// A drop of a parent the walk cannot name would drop this
			// partition unseen.
			b.set(k, walk.Unknown, "")
		} else {
			b.session(func() { b.w.partitions[parentKey] = append(b.w.partitions[parentKey], k) })
		}
	}
	if v.OnCommit == ast.ONCOMMIT_DROP {
		// Dropped at the commit, which the change may run.
		b.set(k, walk.Unknown, "")
		b.set(tk, walk.Unknown, "")
	}
	return true
}

// sameRelation reports whether a reference may name the relation a
// CREATE TABLE creates. PostgreSQL resolves a foreign key's table after
// creating the new one, so a reference of the same name, however
// qualified, may be the table itself.
//
// pg: src/backend/parser/parse_utilcmd.c — transformFKConstraints (the
// constraints are added once the table exists)
func sameRelation(ref, rv *ast.RangeVar) bool {
	return ref.Relname == rv.Relname
}

// constraintKind is the constraint space's kind of a constraint type.
func constraintKind(t ast.ConstrType) string {
	switch t {
	case ast.CONSTR_PRIMARY:
		return conPrimary
	case ast.CONSTR_UNIQUE:
		return conUnique
	case ast.CONSTR_EXCLUSION:
		return conExclusion
	case ast.CONSTR_CHECK:
		return conCheck
	case ast.CONSTR_FOREIGN:
		return conForeign
	}
	return ""
}

func indexBacked(kind string) bool {
	return kind == conPrimary || kind == conUnique || kind == conExclusion
}

// constraint lowers one constraint a statement gives table t, the table
// as written in rv. alter is true for ALTER TABLE, where the constraint
// meets the table's existing ones; inh is true when it also reaches the
// table's inheritance children and partitions.
//
// pg: src/backend/commands/tablecmds.c — ATExecAddConstraint,
// ATAddCheckNNConstraint and ATAddForeignKeyConstraint (a name in use
// fails with 42710); index_create (a key's index takes the constraint's
// name, 42P07)
func (b *builder) constraint(t walk.Key, rv *ast.RangeVar, con *ast.Constraint, alter, inh bool) {
	s := schemaOfKey(t)
	parent := tableOf(t)
	rng := rangeOf(con.Loc)
	switch con.Contype {
	case ast.CONSTR_NOTNULL:
		// PostgreSQL 18 names a NOT NULL constraint.
		b.mayAdd(walk.Scope{Space: spaceConstraint, Parent: parent, Prefix: generatedPrefix(t.Name)})
		if inh {
			b.mayAdd(walk.Scope{Space: spaceConstraint, Deep: true})
		}
		return
	case ast.CONSTR_IDENTITY:
		if con.Options != nil && len(con.Options.Items) > 0 {
			// The sequence may be named anywhere.
			b.mayAdd(walk.Scope{Space: spaceRelation, Deep: true})
		} else {
			b.mayAdd(walk.Scope{Space: spaceRelation, Parent: walk.Path(s), Prefix: generatedPrefix(t.Name)})
		}
		return
	case ast.CONSTR_PRIMARY, ast.CONSTR_UNIQUE, ast.CONSTR_EXCLUSION, ast.CONSTR_CHECK, ast.CONSTR_FOREIGN:
	default:
		return
	}
	kind := constraintKind(con.Contype)
	if con.Contype == ast.CONSTR_FOREIGN && con.Pktable != nil && alter {
		b.exists(con.Pktable, false)
	}
	if con.Contype == ast.CONSTR_PRIMARY {
		pk := primaryKeyOf(t)
		if alter {
			b.need(pk, walk.Absent, nil, &walkNote{rng: rng, present: "multiple primary keys for table " + relation(rv) + " are not allowed"}) // 42P16
		}
		b.set(pk, walk.Present, "")
	}
	if con.Indexname != "" {
		// USING INDEX: the index becomes the constraint's, under its name.
		if n, ok := walkName(con.Indexname); ok {
			b.set(relKey(s, n), walk.Unknown, "")
		} else {
			b.mayRemove(walk.Scope{Space: spaceRelation, Parent: walk.Path(s)})
		}
	}
	partitionIndexes := indexBacked(kind) && b.w.partitioned[t]
	if partitionIndexes {
		b.mayAdd(walk.Scope{Space: spaceRelation, Deep: true})
	}
	name, ok := walkName(con.Conname)
	if !ok {
		// A generated name, or one the walk cannot read.
		b.mayAdd(walk.Scope{Space: spaceConstraint, Parent: parent})
		if indexBacked(kind) {
			b.mayAdd(walk.Scope{Space: spaceRelation, Parent: walk.Path(s), Prefix: generatedPrefix(t.Name)})
			if con.Conname != "" {
				b.mayAdd(walk.Scope{Space: spaceRelation, Parent: walk.Path(s)})
			}
		}
		if inh && (con.Contype == ast.CONSTR_CHECK || partitionIndexes) {
			b.mayAdd(walk.Scope{Space: spaceConstraint, Deep: true})
		}
		return
	}
	// A key's index is created before its constraint is recorded, so a
	// name both in use fails on the index (verified on PostgreSQL 17).
	if indexBacked(kind) {
		ik := relKey(s, name)
		b.need(ik, walk.Absent, nil, &walkNote{rng: rng, present: "relation " + ident(con.Conname) + " already exists"}) // 42P07
		b.set(ik, walk.Present, relIndex)
		b.forget(walk.Scope{Parent: tableOf(ik), Deep: true})
	}
	ck := constraintKey(t, name)
	if alter {
		b.need(ck, walk.Absent, nil, &walkNote{rng: rng, present: "constraint " + ident(con.Conname) + " for relation " + relation(rv) + " already exists"}) // 42710
	}
	if inh && ((con.Contype == ast.CONSTR_CHECK && !con.IsNoInherit) || partitionIndexes) {
		b.mayAdd(walk.Scope{Space: spaceConstraint, Deep: true, Prefix: name})
	}
	b.set(ck, walk.Present, kind)
}

// columnConstraints lowers the constraints of a column ADD COLUMN adds.
func (b *builder) columnConstraints(t walk.Key, rv *ast.RangeVar, def *ast.ColumnDef, inh bool) {
	if isSerial(def.TypeName) {
		b.mayAdd(walk.Scope{Space: spaceRelation, Parent: walk.Path(schemaOfKey(t)), Prefix: generatedPrefix(t.Name)})
	}
	if def.Constraints == nil {
		return
	}
	for _, item := range def.Constraints.Items {
		if con, ok := item.(*ast.Constraint); ok {
			b.constraint(t, rv, con, true, inh)
		}
	}
}

// createIndex follows CREATE INDEX. The index goes to its table's schema.
//
// pg: src/backend/commands/indexcmds.c — DefineIndex (an index only on a
// table, a partitioned table, or a materialized view, 42809)
func (b *builder) createIndex(v *ast.IndexStmt) bool {
	b.last = b.w.callsUserCode(v)
	r := b.exists(v.Relation, false)
	if r.fact == walk.Present {
		b.need(r.key, walk.Present, []string{relTable, relMatview}, &walkNote{rng: rangeOf(v.Relation.Loc), present: "cannot create index on relation " + relation(v.Relation)})
	}
	if r.fact != walk.Present {
		if v.Idxname != "" {
			b.placeUnknown(v.Idxname, false)
		} else {
			b.placeUnknown(v.Relation.Relname, true)
		}
		if b.w.anyPartitioned {
			b.mayAdd(walk.Scope{Space: spaceRelation, Deep: true})
		}
		return true
	}
	s := walk.Path(schemaOfKey(r.key))
	if b.w.partitioned[r.key] {
		// Each partition gets an index of its own, in its own schema.
		b.mayAdd(walk.Scope{Space: spaceRelation, Deep: true})
	}
	if v.Idxname == "" {
		b.mayAdd(walk.Scope{Space: spaceRelation, Parent: s, Prefix: generatedPrefix(r.key.Name)})
		return true
	}
	name, ok := walkName(v.Idxname)
	if !ok {
		b.mayAdd(walk.Scope{Space: spaceRelation, Parent: s})
		return true
	}
	ik := relKey(schemaOfKey(r.key), name)
	prior, _ := b.l.Lookup(ik)
	if !v.IfNotExists {
		b.need(ik, walk.Absent, nil, &walkNote{rng: b.anchor(), present: "relation " + ident(v.Idxname) + " already exists"}) // 42P07
		prior = walk.Absent
	}
	switch prior {
	case walk.Absent:
		b.set(ik, walk.Present, relIndex)
		b.made(ik)
	case walk.Unknown:
		b.set(ik, walk.Present, "")
	}
	b.forget(walk.Scope{Parent: tableOf(ik), Deep: true})
	return true
}

// createView follows CREATE [OR REPLACE] VIEW. A view's columns are not
// followed.
//
// pg: src/backend/commands/view.c — DefineView (a view of a temporary
// relation becomes temporary itself)
func (b *builder) createView(v *ast.ViewStmt) bool {
	rv := v.View
	if rv == nil {
		return false
	}
	rng := rangeOf(rv.Loc)
	temp := rv.Relpersistence == ast.RELPERSISTENCE_TEMP
	if b.w.temporary && rv.Schemaname == "" && !temp {
		// It may read a temporary relation and land in pg_temp instead,
		// needing no schema of the search path.
		b.placeUnknown(rv.Relname, false)
		return true
	}
	c := b.creationSchema(rv.Schemaname, temp)
	if c.fail != "" {
		b.failAt(rng, c.keys, c.fail)
		return true
	}
	name, ok := walkName(rv.Relname)
	if c.fact != walk.Present || !ok {
		b.placeUnknown(rv.Relname, false)
		return true
	}
	k := relKey(c.schema, name)
	tk := typeKey(c.schema, name)
	if !v.Replace {
		b.need(k, walk.Absent, nil, &walkNote{rng: rng, present: "relation " + relation(rv) + " already exists"}) // 42P07
		b.need(tk, walk.Absent, nil, &walkNote{rng: rng, present: "type " + relation(rv) + " already exists"})    // 42710
		b.made(k)
	}
	b.set(k, walk.Present, relView)
	b.set(tk, walk.Present, "")
	b.forget(walk.Scope{Parent: tableOf(k), Deep: true})
	if temp || c.schema == tempSchema {
		b.session(func() { b.w.temporary = true })
	}
	return true
}

// createTableAs follows CREATE TABLE AS, SELECT INTO, and CREATE
// MATERIALIZED VIEW. The new relation's columns come from the query and
// are not followed; it has no constraints.
func (b *builder) createTableAs(into *ast.IntoClause, matview, ifNotExists bool, query ast.Node) bool {
	if into == nil || into.Rel == nil {
		return false
	}
	if _, ok := query.(*ast.ExecuteStmt); ok {
		return false
	}
	if !into.SkipData && b.w.callsUserCode(query) {
		b.last = true
	}
	rv := into.Rel
	rng := rangeOf(rv.Loc)
	temp := rv.Relpersistence == ast.RELPERSISTENCE_TEMP
	c := b.creationSchema(rv.Schemaname, temp)
	if c.fail != "" {
		b.failAt(rng, c.keys, c.fail)
		return true
	}
	name, ok := walkName(rv.Relname)
	if c.fact != walk.Present || !ok {
		b.placeUnknown(rv.Relname, false)
		return true
	}
	kind := relTable
	if matview {
		kind = relMatview
	}
	k := relKey(c.schema, name)
	tk := typeKey(c.schema, name)
	prior, _ := b.l.Lookup(k)
	switch {
	case !ifNotExists:
		b.need(k, walk.Absent, nil, &walkNote{rng: rng, present: "relation " + relation(rv) + " already exists"}) // 42P07
		b.need(tk, walk.Absent, nil, &walkNote{rng: rng, present: "type " + relation(rv) + " already exists"})    // 42710
		prior = walk.Absent
	case prior == walk.Present:
		return true
	}
	parent := tableOf(k)
	if prior == walk.Absent {
		b.set(k, walk.Present, kind)
		b.set(tk, walk.Present, "")
		b.clear(walk.Scope{Parent: parent, Deep: true})
		b.forget(walk.Scope{Space: spaceColumn, Parent: parent})
		b.made(k)
	} else {
		b.set(k, walk.Present, "")
		b.set(tk, walk.Unknown, "")
		b.forget(walk.Scope{Parent: parent, Deep: true})
	}
	if temp || c.schema == tempSchema {
		b.session(func() { b.w.temporary = true })
	}
	if into.OnCommit == ast.ONCOMMIT_DROP {
		b.set(k, walk.Unknown, "")
		b.set(tk, walk.Unknown, "")
	}
	return true
}

// createSequence follows CREATE SEQUENCE. A sequence has no row type.
func (b *builder) createSequence(v *ast.CreateSeqStmt) bool {
	rv := v.Sequence
	if rv == nil {
		return false
	}
	rng := rangeOf(rv.Loc)
	temp := rv.Relpersistence == ast.RELPERSISTENCE_TEMP
	c := b.creationSchema(rv.Schemaname, temp)
	if c.fail != "" {
		b.failAt(rng, c.keys, c.fail)
		return true
	}
	name, ok := walkName(rv.Relname)
	if c.fact != walk.Present || !ok {
		b.placeUnknown(rv.Relname, false)
		return true
	}
	k := relKey(c.schema, name)
	prior, _ := b.l.Lookup(k)
	if !v.IfNotExists {
		b.need(k, walk.Absent, nil, &walkNote{rng: rng, present: "relation " + relation(rv) + " already exists"}) // 42P07
		prior = walk.Absent
	}
	switch prior {
	case walk.Absent:
		b.set(k, walk.Present, relSequence)
		b.made(k)
	case walk.Unknown:
		b.set(k, walk.Present, "")
	}
	b.forget(walk.Scope{Parent: tableOf(k), Deep: true})
	if temp || c.schema == tempSchema {
		b.session(func() { b.w.temporary = true })
	}
	return true
}

// createCompositeType follows CREATE TYPE ... AS (...): a composite type
// takes its name in both the type and the relation namespaces.
func (b *builder) createCompositeType(v *ast.CompositeTypeStmt) bool {
	rv := v.Typevar
	if rv == nil {
		return false
	}
	rng := rangeOf(rv.Loc)
	c := b.creationSchema(rv.Schemaname, false)
	if c.fail != "" {
		b.failAt(rng, c.keys, c.fail)
		return true
	}
	name, ok := walkName(rv.Relname)
	if c.fact != walk.Present || !ok {
		b.placeUnknown(rv.Relname, false)
		return true
	}
	k := relKey(c.schema, name)
	tk := typeKey(c.schema, name)
	b.need(tk, walk.Absent, nil, &walkNote{rng: rng, present: "type " + relation(rv) + " already exists"})    // 42710
	b.need(k, walk.Absent, nil, &walkNote{rng: rng, present: "relation " + relation(rv) + " already exists"}) // 42P07
	b.set(k, walk.Present, relComposite)
	b.set(tk, walk.Present, "")
	b.forget(walk.Scope{Parent: tableOf(k), Deep: true})
	b.made(k)
	return true
}

// createType follows CREATE TYPE ... AS ENUM and CREATE DOMAIN, which
// take a name in the type namespace alone.
func (b *builder) createType(parts []string) bool {
	schema, written, ok := splitName(parts)
	if !ok {
		return false
	}
	c := b.creationSchema(schema, false)
	if c.fail != "" {
		b.fail(c.keys, walk.Absent, c.fail)
		return true
	}
	name, ok := walkName(written)
	if c.fact != walk.Present || !ok {
		b.placeUnknown(written, false)
		return true
	}
	tk := typeKey(c.schema, name)
	b.need(tk, walk.Absent, nil, &walkNote{rng: b.anchor(), present: "type " + qualified(parts) + " already exists"}) // 42710
	b.set(tk, walk.Present, "")
	b.made(tk)
	return true
}

// atNoFacts lists the ALTER TABLE subcommands that change no name the
// walk keeps.
var atNoFacts = map[ast.AlterTableType]bool{
	ast.AT_ChangeOwner: true, ast.AT_SetRelOptions: true, ast.AT_ResetRelOptions: true, ast.AT_ReplaceRelOptions: true,
	ast.AT_EnableTrig: true, ast.AT_EnableAlwaysTrig: true, ast.AT_EnableReplicaTrig: true, ast.AT_DisableTrig: true,
	ast.AT_EnableTrigAll: true, ast.AT_DisableTrigAll: true, ast.AT_EnableTrigUser: true, ast.AT_DisableTrigUser: true,
	ast.AT_EnableRule: true, ast.AT_EnableAlwaysRule: true, ast.AT_EnableReplicaRule: true, ast.AT_DisableRule: true,
	ast.AT_SetLogged: true, ast.AT_SetUnLogged: true, ast.AT_ClusterOn: true, ast.AT_DropCluster: true,
	ast.AT_SetTableSpace: true, ast.AT_SetAccessMethod: true, ast.AT_ReplicaIdentity: true,
	ast.AT_EnableRowSecurity: true, ast.AT_DisableRowSecurity: true, ast.AT_ForceRowSecurity: true, ast.AT_NoForceRowSecurity: true,
	ast.AT_GenericOptions: true,
}

// atColumn lists the ALTER TABLE subcommands that alter one existing
// column.
var atColumn = map[ast.AlterTableType]bool{
	ast.AT_ColumnDefault: true, ast.AT_DropNotNull: true, ast.AT_SetNotNull: true, ast.AT_AlterColumnType: true,
	ast.AT_SetStatistics: true, ast.AT_SetOptions: true, ast.AT_ResetOptions: true, ast.AT_SetStorage: true,
	ast.AT_SetCompression: true, ast.AT_DropExpression: true, ast.AT_SetExpression: true,
	ast.AT_AddIdentity: true, ast.AT_SetIdentity: true, ast.AT_DropIdentity: true, ast.AT_AlterColumnGenericOptions: true,
}

// alterTable follows ALTER TABLE and the ALTER INDEX, VIEW, SEQUENCE,
// MATERIALIZED VIEW, and FOREIGN TABLE statements that share its form.
// Each subcommand is a part of its own.
func (b *builder) alterTable(v *ast.AlterTableStmt) bool {
	if v.Cmds == nil || v.Relation == nil {
		return false
	}
	var cmds []*ast.AlterTableCmd
	for _, item := range v.Cmds.Items {
		cmd, ok := item.(*ast.AlterTableCmd)
		if !ok {
			return false
		}
		cmds = append(cmds, cmd)
	}
	b.last = b.w.callsUserCode(v)
	r := b.exists(v.Relation, v.Missing_ok)
	objType := ast.ObjectType(v.ObjType)
	switch {
	case r.fact == walk.Absent:
		return true // a notice under IF EXISTS
	case r.fact == walk.Unknown, objType != ast.OBJECT_TABLE && objType != ast.OBJECT_FOREIGN_TABLE, !hasColumns(r.kind):
		// Where the relation is, or what the subcommands do to a relation
		// of its kind, is not followed: only subcommands that change
		// nothing may pass.
		for _, cmd := range cmds {
			if !atNoFacts[ast.AlterTableType(cmd.Subtype)] {
				return false
			}
		}
		return true
	}
	for _, cmd := range cmds {
		b.next()
		if !b.alterTableCmd(v.Relation, r.key, cmd) {
			return false
		}
	}
	return true
}

func (b *builder) alterTableCmd(rv *ast.RangeVar, t walk.Key, cmd *ast.AlterTableCmd) bool {
	inh := rv.Inh
	sub := ast.AlterTableType(cmd.Subtype)
	switch {
	case sub == ast.AT_AddColumn:
		def, ok := cmd.Def.(*ast.ColumnDef)
		if !ok {
			return false
		}
		b.addColumn(rv, t, def, cmd.Missing_ok, inh)
	case sub == ast.AT_DropColumn:
		b.dropColumn(rv, t, cmd, inh)
	case atColumn[sub]:
		b.alterColumn(rv, t, cmd, inh)
	case sub == ast.AT_AddConstraint:
		con, ok := cmd.Def.(*ast.Constraint)
		if !ok {
			return false
		}
		b.constraint(t, rv, con, true, inh)
	case sub == ast.AT_DropConstraint:
		b.dropConstraint(rv, t, cmd, inh)
	case sub == ast.AT_ValidateConstraint:
		b.constraintExists(rv, t, cmd.Name, rangeOf(cmd.Loc))
	case sub == ast.AT_AlterConstraint:
		con, ok := cmd.Def.(*ast.Constraint)
		if !ok {
			return false
		}
		b.constraintExists(rv, t, con.Conname, rangeOf(cmd.Loc))
	default:
		return atNoFacts[sub]
	}
	return true
}

// constraintExists checks that a constraint the statement names exists
// (42704). Only a constraint the change dropped is certainly absent: the
// synced schema does not list every constraint.
func (b *builder) constraintExists(rv *ast.RangeVar, t walk.Key, written string, rng review.Range) {
	if name, ok := walkName(written); ok {
		b.need(constraintKey(t, name), walk.Present, nil, &walkNote{rng: rng, absent: "constraint " + ident(written) + " of relation " + relation(rv) + " does not exist"})
	}
}

// addColumn follows ADD COLUMN. A system column's name fails even with IF
// NOT EXISTS, and the column reaches the inheritance children.
//
// pg: src/backend/commands/tablecmds.c — check_for_column_name_collision
func (b *builder) addColumn(rv *ast.RangeVar, t walk.Key, def *ast.ColumnDef, ifNotExists, inh bool) {
	rng := rangeOf(def.Loc)
	name, ok := walkName(def.Colname)
	if !ok {
		b.mayAdd(walk.Scope{Space: spaceColumn, Parent: tableOf(t)})
		if inh {
			b.mayAdd(walk.Scope{Space: spaceColumn, Deep: true})
		}
		b.columnConstraints(t, rv, def, inh)
		return
	}
	if systemColumns[name] {
		b.failAt(rng, nil, "column name "+ident(def.Colname)+" conflicts with a system column name") // 42701
		return
	}
	ck := columnKey(t, name)
	prior, _ := b.l.Lookup(ck)
	if !ifNotExists {
		b.need(ck, walk.Absent, nil, &walkNote{rng: rng, present: "column " + ident(def.Colname) + " of relation " + relation(rv) + " already exists"}) // 42701
		prior = walk.Absent
	}
	if prior == walk.Present {
		return // a notice: the subcommand does nothing
	}
	if inh {
		b.mayAdd(walk.Scope{Space: spaceColumn, Deep: true, Prefix: name})
	}
	b.set(ck, walk.Present, kindColumn)
	b.columnConstraints(t, rv, def, inh)
}

// dropColumn follows DROP COLUMN, which drops the column's indexes and
// constraints, a sequence it owns, and the column of the inheritance
// children.
func (b *builder) dropColumn(rv *ast.RangeVar, t walk.Key, cmd *ast.AlterTableCmd, inh bool) {
	rng := rangeOf(cmd.Loc)
	name, ok := walkName(cmd.Name)
	if ok && systemColumns[name] {
		b.failAt(rng, nil, "cannot drop system column "+ident(cmd.Name)) // 0A000
		return
	}
	if ok {
		ck := columnKey(t, name)
		if !cmd.Missing_ok {
			b.need(ck, walk.Present, nil, &walkNote{rng: rng, absent: "column " + ident(cmd.Name) + " of relation " + relation(rv) + " does not exist"}) // 42703
		}
		if inh {
			b.mayRemove(walk.Scope{Space: spaceColumn, Deep: true, Prefix: name})
		}
		b.set(ck, walk.Absent, "")
	} else {
		b.forget(walk.Scope{Space: spaceColumn, Parent: tableOf(t)})
		if inh {
			b.mayRemove(walk.Scope{Space: spaceColumn, Deep: true})
		}
	}
	b.mayRemove(walk.Scope{Space: spaceConstraint, Parent: tableOf(t)})
	b.mayRemove(walk.Scope{Space: spacePrimaryKey, Parent: tableOf(t)})
	b.mayRemove(walk.Scope{Space: spaceRelation, Parent: walk.Path(schemaOfKey(t)), Kinds: []string{relIndex}})
	b.mayRemove(walk.Scope{Space: spaceRelation, Deep: true, Kinds: []string{relSequence, relIndex}})
	if ast.DropBehavior(cmd.Behavior) == ast.DROP_CASCADE {
		b.mayRemove(walk.Scope{Deep: true})
	}
}

// alterColumn follows the subcommands that alter an existing column.
//
// pg: src/backend/commands/tablecmds.c — "cannot alter system column"
func (b *builder) alterColumn(rv *ast.RangeVar, t walk.Key, cmd *ast.AlterTableCmd, inh bool) {
	rng := rangeOf(cmd.Loc)
	sub := ast.AlterTableType(cmd.Subtype)
	name, ok := walkName(cmd.Name)
	if ok && systemColumns[name] {
		switch sub {
		case ast.AT_AlterColumnType, ast.AT_ColumnDefault, ast.AT_SetNotNull, ast.AT_DropNotNull:
			b.failAt(rng, nil, "cannot alter system column "+ident(cmd.Name)) // 0A000
		}
		return
	}
	if ok {
		b.need(columnKey(t, name), walk.Present, nil, &walkNote{rng: rng, absent: "column " + ident(cmd.Name) + " of relation " + relation(rv) + " does not exist"}) // 42703
	}
	parent := tableOf(t)
	switch sub {
	case ast.AT_SetNotNull, ast.AT_DropNotNull:
		// PostgreSQL 18 adds or drops a named NOT NULL constraint, here
		// and in the inheritance children.
		b.forget(walk.Scope{Space: spaceConstraint, Parent: parent, Prefix: generatedPrefix(t.Name)})
		if inh {
			b.forget(walk.Scope{Space: spaceConstraint, Deep: true})
		}
	case ast.AT_AddIdentity:
		b.mayAdd(walk.Scope{Space: spaceRelation, Deep: true, Prefix: generatedPrefix(t.Name)})
		if def, ok := cmd.Def.(*ast.Constraint); ok && def.Options != nil && len(def.Options.Items) > 0 {
			b.mayAdd(walk.Scope{Space: spaceRelation, Deep: true})
		}
	case ast.AT_DropIdentity:
		b.mayRemove(walk.Scope{Space: spaceRelation, Deep: true, Kinds: []string{relSequence}})
	}
}

// dropConstraint follows DROP CONSTRAINT, which also drops a key's index
// and the constraint of the inheritance children.
func (b *builder) dropConstraint(rv *ast.RangeVar, t walk.Key, cmd *ast.AlterTableCmd, inh bool) {
	s := schemaOfKey(t)
	name, ok := walkName(cmd.Name)
	if !ok {
		b.forget(walk.Scope{Space: spaceConstraint, Parent: tableOf(t)})
		b.mayRemove(walk.Scope{Space: spacePrimaryKey, Parent: tableOf(t)})
		b.mayRemove(walk.Scope{Space: spaceRelation, Parent: walk.Path(s), Kinds: []string{relIndex}})
	} else {
		ck := constraintKey(t, name)
		if !cmd.Missing_ok {
			b.constraintExists(rv, t, cmd.Name, rangeOf(cmd.Loc))
		}
		f, kind := b.l.Lookup(ck)
		b.set(ck, walk.Absent, "")
		switch {
		case f == walk.Present && indexBacked(kind):
			b.set(relKey(s, name), walk.Absent, "")
			if kind == conPrimary {
				b.set(primaryKeyOf(t), walk.Absent, "")
			}
		case f == walk.Present:
		default:
			b.set(relKey(s, name), walk.Unknown, "")
			b.mayRemove(walk.Scope{Space: spacePrimaryKey, Parent: tableOf(t)})
		}
		if inh {
			b.mayRemove(walk.Scope{Space: spaceConstraint, Deep: true, Prefix: name})
		}
	}
	if b.w.partitioned[t] {
		b.mayRemove(walk.Scope{Space: spaceRelation, Deep: true, Kinds: []string{relIndex}})
	}
	if ast.DropBehavior(cmd.Behavior) == ast.DROP_CASCADE {
		b.mayRemove(walk.Scope{Deep: true})
	}
}

// rename follows the RENAME forms of ALTER.
func (b *builder) rename(v *ast.RenameStmt) bool {
	switch v.RenameType {
	case ast.OBJECT_TABLE, ast.OBJECT_VIEW, ast.OBJECT_MATVIEW, ast.OBJECT_INDEX, ast.OBJECT_SEQUENCE, ast.OBJECT_FOREIGN_TABLE:
		return b.renameRelation(v)
	case ast.OBJECT_COLUMN:
		return b.renameColumn(v)
	case ast.OBJECT_TABCONSTRAINT:
		return b.renameConstraint(v)
	case ast.OBJECT_SCHEMA:
		return b.renameSchema(v)
	case ast.OBJECT_FUNCTION, ast.OBJECT_PROCEDURE, ast.OBJECT_ROUTINE, ast.OBJECT_AGGREGATE:
		name := v.Newname
		b.session(func() { b.w.routines[name] = true })
		return true
	case ast.OBJECT_TRIGGER, ast.OBJECT_POLICY, ast.OBJECT_RULE, ast.OBJECT_EVENT_TRIGGER, ast.OBJECT_COLLATION,
		ast.OBJECT_CONVERSION, ast.OBJECT_STATISTIC_EXT, ast.OBJECT_OPCLASS, ast.OBJECT_OPFAMILY,
		ast.OBJECT_PUBLICATION, ast.OBJECT_SUBSCRIPTION, ast.OBJECT_FDW, ast.OBJECT_FOREIGN_SERVER,
		ast.OBJECT_LANGUAGE, ast.OBJECT_TABLESPACE, ast.OBJECT_TSCONFIGURATION, ast.OBJECT_TSDICTIONARY,
		ast.OBJECT_TSPARSER, ast.OBJECT_TSTEMPLATE:
		return true
	}
	return false
}

// renameRelation follows a relation's RENAME TO: its row type goes with
// it, and so do its columns and constraints. Its indexes and sequences
// keep their names; an index's constraint takes the index's new name.
func (b *builder) renameRelation(v *ast.RenameStmt) bool {
	r := b.exists(v.Relation, v.MissingOk)
	switch r.fact {
	case walk.Absent:
		return true
	case walk.Unknown:
		return false
	}
	newName, ok := walkName(v.Newname)
	if !ok {
		return false
	}
	s := schemaOfKey(r.key)
	nk := relKey(s, newName)
	b.need(nk, walk.Absent, nil, &walkNote{rng: b.anchor(), present: "relation " + ident(v.Newname) + " already exists"}) // 42P07
	rowType := hasRowType(r.kind)
	if rowType || r.kind == "" {
		b.need(typeKey(s, newName), walk.Absent, nil, &walkNote{rng: b.anchor(), present: "type " + ident(v.Newname) + " already exists"}) // 42710
	}
	b.set(r.key, walk.Absent, "")
	b.set(nk, walk.Present, r.kind)
	switch {
	case rowType:
		b.set(typeKey(s, r.key.Name), walk.Absent, "")
		b.set(typeKey(s, newName), walk.Present, "")
	case r.kind == "":
		b.set(typeKey(s, r.key.Name), walk.Unknown, "")
		b.set(typeKey(s, newName), walk.Unknown, "")
	}
	b.move(tableOf(r.key), tableOf(nk))
	if r.kind == relIndex || r.kind == "" {
		b.mayRemove(walk.Scope{Space: spaceConstraint, Deep: true, Prefix: r.key.Name})
		b.mayAdd(walk.Scope{Space: spaceConstraint, Deep: true, Prefix: newName})
	}
	b.renamed(r.key, nk)
	return true
}

// renamed moves what the walker keeps by relation key.
func (b *builder) renamed(from, to walk.Key) {
	b.session(func() {
		w := b.w
		if w.partitioned[from] {
			delete(w.partitioned, from)
			w.partitioned[to] = true
		}
		if children, ok := w.partitions[from]; ok {
			delete(w.partitions, from)
			w.partitions[to] = children
		}
		for parent, children := range w.partitions {
			for i, c := range children {
				if c == from {
					w.partitions[parent][i] = to
				}
			}
		}
		for i, k := range w.made {
			if k == from {
				w.made[i] = to
			}
		}
	})
}

// renameColumn follows RENAME COLUMN, which reaches the inheritance
// children too.
//
// pg: src/backend/commands/tablecmds.c — renameatt_internal
func (b *builder) renameColumn(v *ast.RenameStmt) bool {
	r := b.exists(v.Relation, v.MissingOk)
	switch r.fact {
	case walk.Absent:
		return true
	case walk.Unknown:
		return false
	}
	rng := b.anchor()
	oldName, ok1 := walkName(v.Subname)
	newName, ok2 := walkName(v.Newname)
	inh := v.Relation.Inh
	if !ok1 || !ok2 {
		b.forget(walk.Scope{Space: spaceColumn, Parent: tableOf(r.key)})
		if inh {
			b.forget(walk.Scope{Space: spaceColumn, Deep: true})
		}
		return true
	}
	oldKey := columnKey(r.key, oldName)
	newKey := columnKey(r.key, newName)
	if !hasColumns(r.kind) {
		// Another kind's columns are not followed, and some kinds have
		// none to rename (42809).
		b.forget(walk.Scope{Space: spaceColumn, Parent: tableOf(r.key)})
		if inh {
			b.forget(walk.Scope{Space: spaceColumn, Deep: true})
		}
		return true
	}
	if systemColumns[oldName] {
		b.failAt(rng, nil, "cannot rename system column "+ident(v.Subname)) // 0A000
		return true
	}
	if f, _ := b.l.Lookup(oldKey); systemColumns[newName] && f != walk.Absent {
		b.failAt(rng, nil, "column name "+ident(v.Newname)+" conflicts with a system column name") // 42701
		return true
	}
	b.need(oldKey, walk.Present, nil, &walkNote{rng: rng, absent: "column " + ident(v.Subname) + " does not exist"})                                          // 42703
	b.need(newKey, walk.Absent, nil, &walkNote{rng: rng, present: "column " + ident(v.Newname) + " of relation " + relation(v.Relation) + " already exists"}) // 42701
	if inh {
		b.mayRemove(walk.Scope{Space: spaceColumn, Deep: true, Prefix: oldName})
		b.mayAdd(walk.Scope{Space: spaceColumn, Deep: true, Prefix: newName})
	}
	b.set(oldKey, walk.Absent, "")
	b.set(newKey, walk.Present, kindColumn)
	return true
}

// renameConstraint follows RENAME CONSTRAINT: a key's index takes the
// new name too, and the inheritance children's constraint follows.
func (b *builder) renameConstraint(v *ast.RenameStmt) bool {
	r := b.exists(v.Relation, v.MissingOk)
	switch r.fact {
	case walk.Absent:
		return true
	case walk.Unknown:
		return false
	}
	rng := b.anchor()
	oldName, ok1 := walkName(v.Subname)
	newName, ok2 := walkName(v.Newname)
	s := schemaOfKey(r.key)
	if !ok1 || !ok2 {
		b.forget(walk.Scope{Space: spaceConstraint, Deep: true})
		b.forget(walk.Scope{Space: spaceRelation, Parent: walk.Path(s)})
		return true
	}
	oldKey := constraintKey(r.key, oldName)
	newKey := constraintKey(r.key, newName)
	b.need(oldKey, walk.Present, nil, &walkNote{rng: rng, absent: "constraint " + ident(v.Subname) + " for table " + relation(v.Relation) + " does not exist"})    // 42704
	b.need(newKey, walk.Absent, nil, &walkNote{rng: rng, present: "constraint " + ident(v.Newname) + " for relation " + relation(v.Relation) + " already exists"}) // 42710
	f, kind := b.l.Lookup(oldKey)
	switch {
	case f == walk.Present && indexBacked(kind):
		b.need(relKey(s, newName), walk.Absent, nil, &walkNote{rng: rng, present: "relation " + ident(v.Newname) + " already exists"}) // 42P07
		b.set(relKey(s, oldName), walk.Absent, "")
		b.set(relKey(s, newName), walk.Present, relIndex)
	case f != walk.Present:
		b.set(relKey(s, oldName), walk.Unknown, "")
		b.set(relKey(s, newName), walk.Unknown, "")
	}
	if v.Relation.Inh {
		b.mayRemove(walk.Scope{Space: spaceConstraint, Deep: true, Prefix: oldName})
		b.mayAdd(walk.Scope{Space: spaceConstraint, Deep: true, Prefix: newName})
	}
	b.set(oldKey, walk.Absent, "")
	b.set(newKey, walk.Present, kind)
	return true
}

// renameSchema follows ALTER SCHEMA ... RENAME TO: everything in it goes
// with it.
func (b *builder) renameSchema(v *ast.RenameStmt) bool {
	oldName, ok1 := walkName(v.Subname)
	newName, ok2 := walkName(v.Newname)
	if !ok1 || !ok2 {
		return false
	}
	if strings.HasPrefix(newName, "pg_") {
		b.fail(nil, walk.Present, "unacceptable schema name "+ident(v.Newname)) // 42939
		return true
	}
	oldKey := schemaKey(oldName)
	newKey := schemaKey(newName)
	b.need(oldKey, walk.Present, nil, &walkNote{rng: b.anchor(), absent: "schema " + ident(v.Subname) + " does not exist"}) // 3F000
	b.need(newKey, walk.Absent, nil, &walkNote{rng: b.anchor(), present: "schema " + ident(v.Newname) + " already exists"}) // 42P06
	b.set(oldKey, walk.Absent, "")
	b.set(newKey, walk.Present, "")
	b.move(walk.Path(oldName), walk.Path(newName))
	b.session(func() {
		w := b.w
		rekey := func(k walk.Key) walk.Key {
			if k.Parent == oldName {
				k.Parent = newName
			}
			return k
		}
		partitioned := make(map[walk.Key]bool, len(w.partitioned))
		for k := range w.partitioned {
			partitioned[rekey(k)] = true
		}
		w.partitioned = partitioned
		partitions := make(map[walk.Key][]walk.Key, len(w.partitions))
		for parent, children := range w.partitions {
			moved := make([]walk.Key, len(children))
			for i, c := range children {
				moved[i] = rekey(c)
			}
			partitions[rekey(parent)] = moved
		}
		w.partitions = partitions
		for i, k := range w.made {
			w.made[i] = rekey(k)
		}
	})
	return true
}

// setSchema follows ALTER ... SET SCHEMA of a relation: its row type,
// columns, and constraints go with it, and so do its indexes and the
// sequences it owns, under names the walk does not list.
//
// pg: src/backend/commands/tablecmds.c — AlterTableNamespaceInternal
func (b *builder) setSchema(v *ast.AlterObjectSchemaStmt) bool {
	switch v.ObjectType {
	case ast.OBJECT_TABLE, ast.OBJECT_VIEW, ast.OBJECT_MATVIEW, ast.OBJECT_SEQUENCE, ast.OBJECT_FOREIGN_TABLE:
	case ast.OBJECT_FUNCTION, ast.OBJECT_PROCEDURE, ast.OBJECT_ROUTINE, ast.OBJECT_AGGREGATE, ast.OBJECT_COLLATION,
		ast.OBJECT_CONVERSION, ast.OBJECT_OPERATOR, ast.OBJECT_OPCLASS, ast.OBJECT_OPFAMILY, ast.OBJECT_STATISTIC_EXT,
		ast.OBJECT_TSCONFIGURATION, ast.OBJECT_TSDICTIONARY, ast.OBJECT_TSPARSER, ast.OBJECT_TSTEMPLATE:
		return true
	default:
		return false
	}
	r := b.exists(v.Relation, v.MissingOk)
	switch r.fact {
	case walk.Absent:
		return true
	case walk.Unknown:
		return false
	}
	ns, ok := walkName(v.Newschema)
	old := schemaOfKey(r.key)
	if !ok || ns == tempSchema || ns == "pg_catalog" {
		return false
	}
	if ns == old {
		return true
	}
	rng := b.anchor()
	b.need(schemaKey(ns), walk.Present, nil, &walkNote{rng: rng, absent: "schema " + ident(v.Newschema) + " does not exist"}) // 3F000
	nk := relKey(ns, r.key.Name)
	b.need(nk, walk.Absent, nil, &walkNote{rng: rng, present: "relation " + ident(r.key.Name) + " already exists in schema " + ident(v.Newschema)}) // 42P07
	rowType := hasRowType(r.kind)
	if rowType || r.kind == "" {
		b.need(typeKey(ns, r.key.Name), walk.Absent, nil, &walkNote{rng: rng, present: "type " + ident(r.key.Name) + " already exists in schema " + ident(v.Newschema)}) // 42710
	}
	b.set(r.key, walk.Absent, "")
	b.set(nk, walk.Present, r.kind)
	switch {
	case rowType:
		b.set(typeKey(old, r.key.Name), walk.Absent, "")
		b.set(typeKey(ns, r.key.Name), walk.Present, "")
	case r.kind == "":
		b.set(typeKey(old, r.key.Name), walk.Unknown, "")
		b.set(typeKey(ns, r.key.Name), walk.Unknown, "")
	}
	b.move(tableOf(r.key), tableOf(nk))
	b.mayRemove(walk.Scope{Space: spaceRelation, Parent: walk.Path(old), Kinds: []string{relIndex, relSequence}})
	b.mayAdd(walk.Scope{Space: spaceRelation, Parent: walk.Path(ns)})
	b.renamed(r.key, nk)
	return true
}

// dropKinds is, for each relation kind DROP names, the kind it drops and
// the words PostgreSQL uses for it.
var dropKinds = map[ast.ObjectType]struct{ kind, word, article string }{
	ast.OBJECT_TABLE:         {relTable, "table", "a"},
	ast.OBJECT_VIEW:          {relView, "view", "a"},
	ast.OBJECT_MATVIEW:       {relMatview, "materialized view", "a"},
	ast.OBJECT_INDEX:         {relIndex, "index", "an"},
	ast.OBJECT_SEQUENCE:      {relSequence, "sequence", "a"},
	ast.OBJECT_FOREIGN_TABLE: {relForeign, "foreign table", "a"},
}

// dropNoFacts lists the objects whose DROP changes no name the walk keeps
// unless it cascades.
var dropNoFacts = map[ast.ObjectType]bool{
	ast.OBJECT_FUNCTION: true, ast.OBJECT_PROCEDURE: true, ast.OBJECT_ROUTINE: true, ast.OBJECT_AGGREGATE: true,
	ast.OBJECT_TRIGGER: true, ast.OBJECT_POLICY: true, ast.OBJECT_RULE: true, ast.OBJECT_OPERATOR: true,
	ast.OBJECT_CAST: true, ast.OBJECT_COLLATION: true, ast.OBJECT_CONVERSION: true, ast.OBJECT_OPCLASS: true,
	ast.OBJECT_OPFAMILY: true, ast.OBJECT_STATISTIC_EXT: true, ast.OBJECT_PUBLICATION: true,
	ast.OBJECT_EVENT_TRIGGER: true, ast.OBJECT_TSCONFIGURATION: true, ast.OBJECT_TSDICTIONARY: true,
	ast.OBJECT_TSPARSER: true, ast.OBJECT_TSTEMPLATE: true, ast.OBJECT_TRANSFORM: true,
}

// drop follows DROP. Each object is a part of its own.
//
// pg: src/backend/commands/tablecmds.c — RemoveRelations,
// DropErrorMsgNonExistent and DropErrorMsgWrongType (a relation of
// another kind fails even under IF EXISTS, verified on PostgreSQL 17)
func (b *builder) drop(v *ast.DropStmt) bool {
	kind := ast.ObjectType(v.RemoveType)
	cascade := ast.DropBehavior(v.Behavior) == ast.DROP_CASCADE
	if v.Objects == nil {
		return false
	}
	switch {
	case kind == ast.OBJECT_SCHEMA:
		for _, obj := range v.Objects.Items {
			b.next()
			if !b.dropSchema(nameParts(listOf(obj)), v.Missing_ok, cascade) {
				return false
			}
		}
	case kind == ast.OBJECT_TYPE || kind == ast.OBJECT_DOMAIN:
		for _, obj := range v.Objects.Items {
			b.next()
			tn, ok := obj.(*ast.TypeName)
			if !ok {
				return false
			}
			b.dropType(nameParts(tn.Names), v.Missing_ok)
		}
	case dropKinds[kind].kind != "":
		for _, obj := range v.Objects.Items {
			b.next()
			b.dropRelation(kind, nameParts(listOf(obj)), v.Missing_ok)
		}
	case !dropNoFacts[kind]:
		return false
	}
	if cascade {
		b.mayRemove(walk.Scope{Deep: true})
	}
	return true
}

func (b *builder) dropRelation(kind ast.ObjectType, parts []string, missingOK bool) {
	want := dropKinds[kind]
	written := qualified(parts)
	r := b.nameList(parts)
	switch r.fact {
	case walk.Absent:
		if missingOK {
			return
		}
		msg := want.word + " " + written + " does not exist" // 42P01; 42704 for an index
		if r.missingSchema != "" {
			msg = "schema " + ident(r.missingSchema) + " does not exist"
		}
		b.fail(r.keys, walk.Absent, msg)
		return
	case walk.Unknown:
		b.dropUnknown(parts)
		return
	}
	b.need(r.key, walk.Present, []string{want.kind}, &walkNote{rng: b.anchor(), present: written + " is not " + want.article + " " + want.word}) // 42809
	b.dropped(r.key, r.kind)
}

// dropUnknown records a drop of a relation the walk cannot place: the
// name may be gone from any schema, with what it held.
func (b *builder) dropUnknown(parts []string) {
	name, ok := walkName(parts[len(parts)-1])
	if !ok {
		b.mayRemove(walk.Scope{Space: spaceRelation, Deep: true})
		b.mayRemove(walk.Scope{Space: spaceType, Deep: true})
	} else {
		b.mayRemove(walk.Scope{Space: spaceRelation, Deep: true, Prefix: name})
		b.mayRemove(walk.Scope{Space: spaceType, Deep: true, Prefix: name})
	}
	b.mayRemove(walk.Scope{Space: spaceRelation, Deep: true, Kinds: []string{relIndex, relSequence}})
	for _, children := range b.w.partitions {
		for _, c := range children {
			b.set(c, walk.Unknown, "")
		}
	}
}

// dropped records the effects of dropping a relation the walk placed: its
// row type, columns, constraints, indexes, owned sequences, and
// partitions go with it.
func (b *builder) dropped(k walk.Key, kind string) {
	s := schemaOfKey(k)
	b.set(k, walk.Absent, "")
	switch {
	case hasRowType(kind):
		b.set(typeKey(s, k.Name), walk.Absent, "")
	case kind == "":
		b.set(typeKey(s, k.Name), walk.Unknown, "")
	}
	b.clear(walk.Scope{Parent: tableOf(k), Deep: true})
	switch kind {
	case relTable, relMatview, relForeign, "":
		b.mayRemove(walk.Scope{Space: spaceRelation, Parent: walk.Path(s), Kinds: []string{relIndex}})
		b.mayRemove(walk.Scope{Space: spaceRelation, Deep: true, Kinds: []string{relSequence}})
	}
	if b.w.partitioned[k] || kind == "" || (kind == relIndex && b.w.anyPartitioned) {
		b.mayRemove(walk.Scope{Space: spaceRelation, Deep: true, Kinds: []string{relIndex}})
	}
	for _, child := range b.w.partitions[k] {
		b.dropped(child, relTable)
	}
}

// dropSchema follows DROP SCHEMA. Without CASCADE it fails while the
// schema certainly holds an object (2BP01, verified on PostgreSQL 17).
func (b *builder) dropSchema(parts []string, missingOK, cascade bool) bool {
	if len(parts) != 1 {
		return false
	}
	name, ok := walkName(parts[0])
	if !ok {
		return false
	}
	k := schemaKey(name)
	if !missingOK {
		b.need(k, walk.Present, nil, &walkNote{rng: b.anchor(), absent: "schema " + ident(parts[0]) + " does not exist"}) // 3F000
	}
	if f, _ := b.l.Lookup(k); f == walk.Present && !cascade {
		if held, ok := b.holds(name); ok {
			b.fail([]walk.Key{k, held}, walk.Present, "cannot drop schema "+ident(parts[0])+" because other objects depend on it")
		}
	}
	b.set(k, walk.Absent, "")
	b.clear(walk.Scope{Parent: walk.Path(name), Deep: true})
	return true
}

// holds returns a relation or type that certainly exists in the schema:
// one the synced schema lists or the change made, still there.
func (b *builder) holds(schema string) (walk.Key, bool) {
	for _, k := range b.w.pg.inSchema[schema] {
		if f, _ := b.l.Lookup(k); f == walk.Present {
			return k, true
		}
	}
	for _, k := range b.w.made {
		if k.Parent == walk.Path(schema) {
			if f, _ := b.l.Lookup(k); f == walk.Present {
				return k, true
			}
		}
	}
	return walk.Key{}, false
}

// dropType follows DROP TYPE and DROP DOMAIN. An unqualified type name may
// be one of pg_catalog's, which the walk does not know, so only a
// qualified one is placed.
func (b *builder) dropType(parts []string, missingOK bool) {
	schema, written, ok := splitName(parts)
	name, ok2 := walkName(written)
	s, ok3 := walkName(schema)
	if !ok || !ok2 || !ok3 || s == tempSchema || s == "pg_catalog" {
		if ok2 {
			b.mayRemove(walk.Scope{Space: spaceType, Deep: true, Prefix: name})
			b.mayRemove(walk.Scope{Space: spaceRelation, Deep: true, Prefix: name, Kinds: []string{relComposite}})
		} else {
			b.mayRemove(walk.Scope{Space: spaceType, Deep: true})
			b.mayRemove(walk.Scope{Space: spaceRelation, Deep: true, Kinds: []string{relComposite}})
		}
		return
	}
	sk := schemaKey(s)
	switch f, _ := b.l.Lookup(sk); f {
	case walk.Absent:
		if !missingOK {
			b.fail([]walk.Key{sk}, walk.Absent, "schema "+ident(schema)+" does not exist") // 3F000
		}
		return
	case walk.Unknown:
		b.mayRemove(walk.Scope{Space: spaceType, Deep: true, Prefix: name})
		b.mayRemove(walk.Scope{Space: spaceRelation, Deep: true, Prefix: name, Kinds: []string{relComposite}})
		return
	}
	tk := typeKey(s, name)
	if !missingOK {
		b.need(tk, walk.Present, nil, &walkNote{rng: b.anchor(), absent: "type " + qualified(parts) + " does not exist"}) // 42704
	}
	b.set(tk, walk.Absent, "")
	rk := relKey(s, name)
	switch f, kind := b.l.Lookup(rk); {
	case f == walk.Present && kind == relComposite:
		b.set(rk, walk.Absent, "")
		b.clear(walk.Scope{Parent: tableOf(rk), Deep: true})
	case f == walk.Unknown, f == walk.Present && kind == "":
		b.set(rk, walk.Unknown, "")
	}
}

// commentKinds is the relation kind each COMMENT ON object kind names.
var commentKinds = map[ast.ObjectType]struct{ kind, word, article string }{
	ast.OBJECT_TABLE:         {relTable, "table", "a"},
	ast.OBJECT_VIEW:          {relView, "view", "a"},
	ast.OBJECT_MATVIEW:       {relMatview, "materialized view", "a"},
	ast.OBJECT_INDEX:         {relIndex, "index", "an"},
	ast.OBJECT_SEQUENCE:      {relSequence, "sequence", "a"},
	ast.OBJECT_FOREIGN_TABLE: {relForeign, "foreign table", "a"},
}

// comment follows COMMENT ON, which changes no name: the object must
// exist.
//
// pg: src/backend/catalog/objectaddress.c — get_object_address
func (b *builder) comment(v *ast.CommentStmt) bool {
	parts := nameParts(listOf(v.Object))
	switch kind := ast.ObjectType(v.Objtype); {
	case commentKinds[kind].kind != "":
		want := commentKinds[kind]
		r := b.nameList(parts)
		switch r.fact {
		case walk.Absent:
			b.fail(r.keys, walk.Absent, missing(r, qualified(parts))) // 42P01
		case walk.Present:
			b.need(r.key, walk.Present, []string{want.kind}, &walkNote{rng: b.anchor(), present: qualified(parts) + " is not " + want.article + " " + want.word}) // 42809
		}
	case kind == ast.OBJECT_COLUMN:
		if len(parts) < 2 {
			return true
		}
		table := parts[:len(parts)-1]
		r := b.nameList(table)
		switch r.fact {
		case walk.Absent:
			b.fail(r.keys, walk.Absent, missing(r, qualified(table)))
		case walk.Present:
			if col, ok := walkName(parts[len(parts)-1]); ok && !systemColumns[col] && hasColumns(r.kind) {
				b.need(columnKey(r.key, col), walk.Present, nil, &walkNote{rng: b.anchor(), absent: "column " + ident(parts[len(parts)-1]) + " of relation " + qualified(table) + " does not exist"}) // 42703
			}
		}
	case kind == ast.OBJECT_SCHEMA:
		if s, ok := v.Object.(*ast.String); ok {
			if name, ok := walkName(s.Str); ok {
				b.need(schemaKey(name), walk.Present, nil, &walkNote{rng: b.anchor(), absent: "schema " + ident(s.Str) + " does not exist"}) // 3F000
			}
		}
	}
	return true
}

// grant follows GRANT and REVOKE, which change no name: a table or schema
// they name must exist.
func (b *builder) grant(v *ast.GrantStmt) bool {
	if v.Objects == nil {
		return true
	}
	switch {
	case v.Targtype == ast.ACL_TARGET_OBJECT && (v.Objtype == ast.OBJECT_TABLE || v.Objtype == ast.OBJECT_SEQUENCE):
		for _, item := range v.Objects.Items {
			if rv, ok := item.(*ast.RangeVar); ok {
				b.exists(rv, false) // 42P01
			}
		}
	case v.Targtype == ast.ACL_TARGET_OBJECT && v.Objtype == ast.OBJECT_SCHEMA,
		v.Targtype == ast.ACL_TARGET_ALL_IN_SCHEMA:
		for _, item := range v.Objects.Items {
			parts := nameParts(listOf(item))
			if s, ok := item.(*ast.String); ok {
				parts = []string{s.Str}
			}
			if len(parts) != 1 {
				continue
			}
			if name, ok := walkName(parts[0]); ok {
				b.need(schemaKey(name), walk.Present, nil, &walkNote{rng: b.anchor(), absent: "schema " + ident(parts[0]) + " does not exist"}) // 3F000
			}
		}
	}
	return true
}

// truncate follows TRUNCATE, which changes no name: each relation must be
// a table (42P01, 42809). A TRUNCATE trigger may run anything.
func (b *builder) truncate(v *ast.TruncateStmt) bool {
	if v.Relations != nil {
		for _, item := range v.Relations.Items {
			rv, ok := item.(*ast.RangeVar)
			if !ok {
				continue
			}
			if r := b.exists(rv, false); r.fact == walk.Present {
				b.need(r.key, walk.Present, []string{relTable, relForeign}, &walkNote{rng: rangeOf(rv.Loc), present: relation(rv) + " is not a table"})
			}
		}
	}
	b.last = b.w.triggered
	return true
}

// lock follows LOCK, which changes no name: each relation must exist.
func (b *builder) lock(v *ast.LockStmt) bool {
	if v.Relations != nil {
		for _, item := range v.Relations.Items {
			if rv, ok := item.(*ast.RangeVar); ok {
				b.exists(rv, false)
			}
		}
	}
	return true
}

// dml follows INSERT, UPDATE, DELETE, and MERGE, which change no name:
// the target must exist (a WITH query's name is never the target,
// verified on PostgreSQL 17). A trigger or a called function may run
// anything, so the walk ends after one that may.
func (b *builder) dml(rv *ast.RangeVar, n ast.Node) bool {
	if rv != nil {
		b.opened(rv)
	}
	b.last = b.w.triggered || b.w.callsUserCode(n)
	return true
}

// setVariable follows SET and RESET of the search path and the role. A
// transaction-local setting reverts at commit and RESET returns to a
// default the target does not record, so after either the setting is
// unknown.
func (b *builder) setVariable(v *ast.VariableSetStmt) bool {
	w := b.w
	if v.Kind == ast.VAR_RESET_ALL {
		b.session(func() { w.path = nil })
		return true
	}
	switch v.Name {
	case "search_path":
		if v.Kind != ast.VAR_SET_VALUE || v.IsLocal || !plainNames(v.Args) {
			b.session(func() { w.path = nil })
			return true
		}
		path := make([]string, 0, len(v.Args.Items))
		for _, item := range v.Args.Items {
			name, _ := constString(item)
			path = append(path, name)
		}
		b.session(func() { w.path = path })
	case "role":
		var user string
		switch {
		case v.IsLocal:
		case v.Kind == ast.VAR_SET_DEFAULT || v.Kind == ast.VAR_RESET:
			user = w.session
		case v.Kind == ast.VAR_SET_VALUE && v.Args != nil && len(v.Args.Items) == 1:
			user, _ = constString(v.Args.Items[0])
			if user == "none" {
				user = w.session
			}
		}
		b.session(func() { w.user = user })
	case "session_authorization":
		b.session(func() { w.user, w.session = "", "" })
	}
	return true
}

// createTrigger follows CREATE TRIGGER: its table must exist, and from now
// on DML may fire it.
func (b *builder) createTrigger(v *ast.CreateTrigStmt) bool {
	r := b.exists(v.Relation, false)
	b.session(func() { b.w.triggered = true })
	if v.IsConstraint {
		// A constraint trigger is a constraint under the trigger's name.
		if r.fact == walk.Present {
			b.mayAdd(walk.Scope{Space: spaceConstraint, Parent: tableOf(r.key)})
		} else {
			b.mayAdd(walk.Scope{Space: spaceConstraint, Deep: true})
		}
	}
	return true
}
