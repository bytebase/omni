package review

import (
	"slices"
	"strconv"
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
		return b.createEnum(v)
	case *ast.CreateDomainStmt:
		return b.createDomain(v)
	case *ast.CreateRangeStmt:
		return b.createRange(v)
	case *ast.DefineStmt:
		return b.defineObject(v)
	case *ast.AlterDomainStmt:
		return b.alterDomain(v)
	case *ast.AlterEnumStmt:
		return b.alterEnum(v)
	case *ast.AlterTypeStmt:
		return b.alterType(v)
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
		return b.alterSequence(v)
	case *ast.RefreshMatViewStmt:
		b.existsOfKind(v.Relation, relMatview, "materialized view")
		// It runs the view's stored query, which may call a function
		// that changes the schema (verified on PostgreSQL 17 and 18).
		b.Last = true
		return true
	case *ast.CreatePolicyStmt:
		return b.createPolicy(v)
	case *ast.AlterPolicyStmt:
		return b.alterPolicy(v)
	case *ast.AlterObjectDependsStmt:
		return b.alterObjectDepends(v)
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
			// leaves the path unknown. Unqualified, set_config may be a
			// user function a schema before pg_catalog on the path holds
			// (verified on PostgreSQL 17 and 18).
			if path, ok := dumpSearchPath(v, calls); ok && qualifiedCatalog(calls[0]) {
				known := append([]string{}, path...)
				b.session(func() { b.w.path = known })
				return true
			}
			b.session(func() { b.w.path = nil })
		}
		b.Last = b.w.callsUserCode(v)
		return true
	case *ast.VariableSetStmt:
		return b.setVariable(v)
	case *ast.TransactionStmt:
		switch v.Kind {
		case ast.TRANS_STMT_BEGIN, ast.TRANS_STMT_START:
			b.session(func() { b.w.inBlock = true })
			return true
		case ast.TRANS_STMT_COMMIT:
			// The block ends, AND CHAIN starting another, and with it
			// what SET LOCAL set, as Bytebase's SET LOCAL ROLE in tenant
			// mode.
			b.session(func() {
				if b.w.tenant && b.w.inBlock {
					b.w.user = ""
				}
				b.w.inBlock = v.Chain
			})
			return true
		case ast.TRANS_STMT_SAVEPOINT, ast.TRANS_STMT_RELEASE:
			return true
		}
		// A rollback undoes statements the walk took as done.
		return false
	case *ast.CreateFunctionStmt:
		return b.createFunction(v)
	case *ast.AlterFunctionStmt:
		return b.alterFunction(v)
	case *ast.CreateTrigStmt:
		return b.createTrigger(v)
	case *ast.AlterOwnerStmt:
		return b.alterOwner(v)
	case *ast.VariableShowStmt, *ast.NotifyStmt, *ast.ListenStmt, *ast.UnlistenStmt,
		*ast.CheckPointStmt, *ast.VacuumStmt, *ast.ReindexStmt, *ast.ClusterStmt,
		*ast.AlterDefaultPrivilegesStmt,
		*ast.CreateRoleStmt, *ast.AlterRoleStmt, *ast.AlterRoleSetStmt, *ast.DropRoleStmt,
		*ast.GrantRoleStmt:
		// None changes a name the walk keeps.
		return true
	}
	return false
}

// qualifiedCatalog reports whether a call names pg_catalog's function.
func qualifiedCatalog(fc *ast.FuncCall) bool {
	parts := nameParts(fc.Funcname)
	return len(parts) == 2 && parts[0] == "pg_catalog"
}

// missing is the message for a relation that resolves nowhere, of the
// kind the statement names, "" for any.
func missing(r resolved, kind, written string) string {
	if r.missingSchema != "" {
		return schemaNotFound(r.missingSchema)
	}
	return notFound(kind, written)
}

// exists checks that a relation the statement names exists (42P01; a
// schema that does not exist is 3F000).
func (b *builder) exists(rv *ast.RangeVar, missingOK bool) resolved {
	r := b.rangeVar(rv)
	switch {
	case r.fail != "":
		b.failAt(rangeOf(rv.Loc), nil, review.CrossDatabaseReference, r.fail)
	case r.fact == walk.Absent && !missingOK:
		b.failAt(rangeOf(rv.Loc), r.keys, review.ObjectNotFound, missing(r, "", relation(rv)))
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
	switch {
	case r.fail != "":
		b.failAt(rangeOf(rv.Loc), nil, review.CrossDatabaseReference, r.fail)
	case r.fact == walk.Absent:
		b.failAt(rangeOf(rv.Loc), r.keys, review.ObjectNotFound, notFound("", relation(rv)))
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
		b.need(r.key, walk.Present, []string{kind}, ofKind(rangeOf(rv.Loc), relation(rv), word))
	}
}

// crossDatabase fails a statement that creates in another database
// (0A000).
func (b *builder) crossDatabase(rv *ast.RangeVar) bool {
	if !b.w.crossDatabase(rv.Catalogname) {
		return false
	}
	b.failAt(rangeOf(rv.Loc), nil, review.CrossDatabaseReference, crossDatabaseMessage(relation(rv)))
	return true
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
		b.fail(nil, walk.Present, review.ReservedName, reservedSchema(written)) // 42939
		return true
	}
	k := schemaKey(name)
	prior, _ := b.l.Lookup(k)
	if !v.IfNotExists {
		b.need(k, walk.Absent, nil, &walkNote{rng: b.anchor(), present: alreadyExists("schema", ident(written))}) // 42P06
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
// without CASCADE finds. A new object is no partition the walker listed
// under that name before.
func (b *builder) made(k walk.Key) {
	b.session(func() { b.w.made = append(b.w.made, k) })
	b.unlisted(k)
}

// unlisted forgets k as a partition, and its own partitions, once the
// relation that had them is gone or replaced: a later drop of the parent
// must not reach whatever takes the name next.
func (b *builder) unlisted(k walk.Key) {
	b.session(func() { b.w.unlist(k) })
}

// unlist forgets k as a partitioned table, its partitions, and k as a
// partition.
func (w *walker) unlist(k walk.Key) {
	delete(w.partitions, k)
	delete(w.partitioned, k)
	for parent, children := range w.partitions {
		w.partitions[parent] = slices.DeleteFunc(children, func(c walk.Key) bool { return c == k })
	}
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

	if b.crossDatabase(rv) {
		return true
	}
	c := b.creationSchema(rv.Schemaname, rv.Relname, temp)
	if c.fail != "" {
		b.failAt(rng, c.keys, c.code, c.fail) // 3F000
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
		b.need(k, walk.Absent, nil, taken(rng, relTable, relation(rv)))                                 // 42P07
		b.need(tk, walk.Absent, nil, &walkNote{rng: rng, present: alreadyExists("type", relation(rv))}) // 42710
		prior = walk.Absent
	case prior == walk.Present:
		return true // a notice: nothing is created
	case prior == walk.Absent:
		b.need(tk, walk.Absent, nil, &walkNote{rng: rng, present: alreadyExists("type", relation(rv))})
	}

	if prior == walk.Absent {
		if v.AccessMethod == "" || v.AccessMethod == "heap" {
			b.checkHeapRelOptions(v.Options, rng)
		}
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
			if con.Contype == ast.CONSTR_FOREIGN && con.Pktable != nil && !b.mayBeSelf(con.Pktable, name, c.schema) {
				b.exists(con.Pktable, false)
			}
		}
		if v.OfTypename != nil {
			if !b.checkCompositeType(v.OfTypename, rangeOf(v.OfTypename.Loc)) {
				if impossible, keys := b.selfTypeImpossible(c.schema, name, v.OfTypename); impossible {
					b.failAt(rangeOf(v.OfTypename.Loc), keys, review.UnsupportedForm, "Table "+review.Code(relation(rv))+" can't be defined `OF` itself.")
					return true
				}
			}
		}
		if !temp && c.schema != tempSchema && v.OnCommit != ast.ONCOMMIT_NOOP && v.OnCommit != 0 {
			b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("`ON COMMIT` on a non-temporary table")) // 42P16
			return true
		}
		// Within the statement, judged from the text alone. A repeated
		// column fails before a system column's name does (verified on
		// PostgreSQL 17 and 18).
		if !partitionOf && v.OfTypename == nil {
			seen := make(map[string]bool, len(cols))
			for _, col := range cols {
				n, ok := storedKey(col.Colname)
				if !ok {
					continue
				}
				if seen[n] {
					b.failAt(rangeOf(col.Loc), nil, review.DuplicateName, twice("declares", "column", ident(col.Colname))) // 42701
					return true
				}
				seen[n] = true
			}
			for _, col := range cols {
				if n, ok := walkName(col.Colname); ok && systemColumns[n] {
					b.failAt(rangeOf(col.Loc), nil, review.ReservedName, systemColumnName(col.Colname)) // 42701
					return true
				}
			}
			if plain {
				for _, con := range cons {
					if con.Contype != ast.CONSTR_CHECK {
						continue
					}
					if col, ok := bareColumnRef(con.RawExpr); ok && !seen[col] && !systemColumns[col] {
						if col == name {
							b.failAt(rangeOf(con.Loc), nil, review.UnsupportedForm, unsupported("non-boolean expression in `CHECK` constraint")) // 42804
						} else {
							b.failAt(rangeOf(con.Loc), nil, review.ObjectNotFound, notFoundIn("column", ident(col), relation(rv))) // 42703
						}
						return true
					}
				}
			}
			for _, col := range cols {
				if col.TypeName == nil || isSerial(col.TypeName) {
					continue
				}
				if b.checkColumnType(col.TypeName, rangeOf(col.Loc)) {
					return true
				}
				if impossible, keys := b.selfTypeImpossible(c.schema, name, col.TypeName); impossible {
					b.failAt(rangeOf(col.Loc), keys, review.UnsupportedForm, "Table "+review.Code(relation(rv))+" can't use its own type in `CREATE TABLE`.")
					return true
				}
			}
		}
		primaries := 0
		for _, con := range cons {
			if con.Contype == ast.CONSTR_PRIMARY {
				primaries++
				if primaries > 1 {
					b.failAt(rangeOf(con.Loc), nil, review.NameTaken, hasPrimaryKey(relation(rv))) // 42P16
					return true
				}
			}
		}
		if msg, rng, ok := duplicateConstraint(cons); ok {
			b.failAt(rng, nil, review.DuplicateName, msg)
			return true
		}
		// A key's index takes the name in the relation namespace the
		// table already holds (42P07, verified on PostgreSQL 17).
		for _, con := range cons {
			if n, ok := walkName(con.Conname); ok && n == name && indexBacked(constraintKind(con.Contype)) {
				b.failAt(rangeOf(con.Loc), nil, review.DuplicateName, twice("declares", "", ident(con.Conname)))
				return true
			}
		}
	}

	parent := tableOf(k)
	switch prior {
	case walk.Absent:
		b.set(k, walk.Present, relTable)
		b.set(tk, walk.Present, typOther)
		b.mayAddArrayType(c.schema)
		b.clear(walk.Scope{Parent: parent, Deep: true})
		if !plain {
			b.forget(walk.Scope{Space: spaceColumn, Parent: parent})
			b.forget(walk.Scope{Space: spaceConstraint, Parent: parent})
			b.forget(walk.Scope{Space: spaceTrigger, Parent: parent})
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
		// may have no row type: a sequence or an index. The indexes and
		// sequences of the keys may or may not be there.
		b.set(k, walk.Present, "")
		b.set(tk, walk.Unknown, "")
		b.forget(walk.Scope{Parent: parent, Deep: true})
		b.mayIndex(k, cons)
		if v.Partspec != nil || partitionOf {
			// Partitioned, or a partition a drop of its parent takes,
			// only if it was created.
			b.set(k, walk.Unknown, "")
		}
	}
	if generated {
		b.mayAdd(walk.Scope{Space: spaceRelation, Parent: walk.Path(c.schema), Prefix: generatedPrefix(name)})
	}
	if temp || c.schema == tempSchema {
		b.session(func() { b.w.temporary = true })
	}
	if v.Partspec != nil {
		b.session(func() {
			if prior == walk.Absent {
				b.w.partitioned[k] = true
			}
			b.w.anyPartitioned = true
		})
	}
	if partitionOf && prior == walk.Absent {
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
			// Only a partitioned table takes a partition, and the sync
			// shows one without partitions as an ordinary table.
			b.session(func() {
				b.w.partitions[parentKey] = append(b.w.partitions[parentKey], k)
				b.w.partitioned[parentKey] = true
				b.w.anyPartitioned = true
			})
		}
	}
	if v.OnCommit == ast.ONCOMMIT_DROP {
		// Dropped at the commit, which the change may run.
		b.set(k, walk.Unknown, "")
		b.set(tk, walk.Unknown, "")
	}
	return true
}

// storedKey is a key two names share exactly when PostgreSQL stores them
// alike: the stored name, or, for a name walkName cannot read, the
// spelling itself, which stores alike only when it is identical (a
// non-ASCII name; verified on PostgreSQL 17).
func storedKey(name string) (string, bool) {
	if name == "" {
		return "", false
	}
	if n, ok := walkName(name); ok {
		return n, true
	}
	return "\x00" + name, true
}

// duplicateConstraint finds two constraints a CREATE TABLE names alike,
// which PostgreSQL refuses (verified on PostgreSQL 17): two keys collide
// on their index's name (42P07), two checks on the check (42710), and any
// other pair on the constraint (42710). A NOT NULL constraint is left
// out: PostgreSQL 18 names it and refuses a second constraint of its name,
// while 17 ignores the name.
//
// pg: src/backend/catalog/heap.c — AddRelationNewConstraints;
// src/backend/catalog/index.c — index_create
func duplicateConstraint(cons []*ast.Constraint) (string, review.Range, bool) {
	seen := make(map[string]bool)
	for _, con := range cons {
		key, ok := storedKey(con.Conname)
		if constraintKind(con.Contype) == "" || !ok {
			continue
		}
		if !seen[key] {
			seen[key] = true
			continue
		}
		return twice("declares", "constraint", ident(con.Conname)), rangeOf(con.Loc), true
	}
	return "", review.Range{}, false
}

// mayBeSelf reports whether a foreign key's reference may name the table
// a CREATE TABLE creates, name in schema. PostgreSQL resolves the
// reference after creating the table, so one of the same name, written
// with that schema or without one where the search path reaches it, may
// be the table itself; one that names another schema, or another
// database, or that the search path does not reach, is not (verified on
// PostgreSQL 17 and 18).
//
// pg: src/backend/parser/parse_utilcmd.c — transformFKConstraints (the
// constraints are added once the table exists)
func (b *builder) mayBeSelf(ref *ast.RangeVar, name, schema string) bool {
	if b.w.crossDatabase(ref.Catalogname) {
		return false
	}
	if n, ok := walkName(ref.Relname); !ok || n != name {
		return !ok
	}
	if ref.Schemaname == "" {
		// Unqualified, it finds the new table only where the search path
		// reaches its schema; pg_temp is searched first.
		path, ok := b.w.searchPath()
		return !ok || schema == tempSchema || containsString(path, schema)
	}
	s, ok := walkName(ref.Schemaname)
	return !ok || s == schema
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
		// PostgreSQL 18 records a NOT NULL constraint, under the name the
		// SQL gives it or one it generates; 17 records none.
		if name, ok := walkName(con.Conname); ok {
			b.set(constraintKey(t, name), walk.Unknown, "")
		} else {
			b.mayAdd(walk.Scope{Space: spaceConstraint, Parent: parent, Prefix: generatedPrefix(t.Name)})
		}
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
	var fkTarget resolved
	if con.Contype == ast.CONSTR_FOREIGN && con.Pktable != nil && alter {
		fkTarget = b.exists(con.Pktable, false)
	}
	if alter && con.Indexname == "" {
		b.checkConstraintColumns(t, rv, con, rng)
	}
	if alter && con.Contype == ast.CONSTR_CHECK {
		if col, ok := bareColumnRef(con.RawExpr); ok && b.colTouch[col] == 0 {
			if sc, ok := b.w.syncedCol(t, col); ok && isPinnedNonBoolType(sc.typ) {
				b.failAt(rng, []walk.Key{t}, review.UnsupportedForm, unsupported("non-boolean expression in `CHECK` constraint")) // 42804
				return
			}
		}
	}
	if alter && con.Contype == ast.CONSTR_FOREIGN && con.Pktable != nil && fkTarget.fact == walk.Present &&
		(con.PkAttrs == nil || len(con.PkAttrs.Items) == 0) && con.FkAttrs != nil &&
		!b.w.isTableDirty(fkTarget.key) && !(fkTarget.key == t && b.pkChanged) {
		if pkCols := b.w.pg.pkCols[[2]string{schemaOfKey(fkTarget.key), fkTarget.key.Name}]; len(pkCols) > 0 {
			if len(con.FkAttrs.Items) != len(pkCols) {
				b.failAt(rng, []walk.Key{fkTarget.key}, review.UnsupportedForm, unsupported("foreign key column count mismatch with primary key of "+review.Code(relation(con.Pktable)))) // 42830
				return
			}
		}
	}
	if con.Contype == ast.CONSTR_PRIMARY {
		pk := primaryKeyOf(t)
		if alter {
			b.need(pk, walk.Absent, nil, &walkNote{rng: rng, present: hasPrimaryKey(relation(rv))}) // 42P16
		}
		b.set(pk, walk.Present, "")
	}
	if con.Indexname != "" {
		b.usingIndex(t, rv, con, kind, alter)
		return
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
		b.need(ik, walk.Absent, nil, taken(rng, relIndex, ident(con.Conname))) // 42P07
		b.set(ik, walk.Present, relIndex)
		b.forget(walk.Scope{Parent: tableOf(ik), Deep: true})
	}
	ck := constraintKey(t, name)
	if alter {
		b.need(ck, walk.Absent, nil, &walkNote{rng: rng, present: alreadyExistsIn("constraint", ident(con.Conname), relation(rv))}) // 42710
	}
	if inh && ((con.Contype == ast.CONSTR_CHECK && !con.IsNoInherit) || partitionIndexes) {
		b.mayAdd(walk.Scope{Space: spaceConstraint, Deep: true, Prefix: name})
	}
	b.set(ck, walk.Present, kind)
}

// usingIndex lowers ADD CONSTRAINT ... USING INDEX: an existing index of
// the table's schema becomes the constraint's, renamed to the
// constraint's name when it has another (verified on PostgreSQL 17: the
// index must exist, 42704, and be an index, 42809, and a new name must be
// free, 42P07).
//
// pg: src/backend/parser/parse_utilcmd.c — transformIndexConstraint;
// src/backend/commands/tablecmds.c — ATExecAddIndexConstraint
func (b *builder) usingIndex(t walk.Key, rv *ast.RangeVar, con *ast.Constraint, kind string, alter bool) {
	s := schemaOfKey(t)
	rng := rangeOf(con.Loc)
	idx, ok := walkName(con.Indexname)
	if !ok {
		b.forget(walk.Scope{Space: spaceRelation, Parent: walk.Path(s)})
		b.forget(walk.Scope{Space: spaceConstraint, Parent: tableOf(t)})
		return
	}
	ik := relKey(s, idx)
	note := ofKind(rng, ident(con.Indexname), relIndex)
	note.absent = notFound(relIndex, ident(con.Indexname))
	b.need(ik, walk.Present, []string{relIndex}, note)
	name, written := idx, con.Indexname
	if con.Conname != "" {
		n, ok := walkName(con.Conname)
		if !ok {
			b.forget(walk.Scope{Space: spaceRelation, Parent: walk.Path(s)})
			b.forget(walk.Scope{Space: spaceConstraint, Parent: tableOf(t)})
			return
		}
		name, written = n, con.Conname
	}
	if name != idx {
		nk := relKey(s, name)
		b.need(nk, walk.Absent, nil, taken(rng, relIndex, ident(written))) // 42P07
		b.set(ik, walk.Absent, "")
		b.set(nk, walk.Present, relIndex)
		b.forget(walk.Scope{Parent: tableOf(nk), Deep: true})
	}
	ck := constraintKey(t, name)
	if alter {
		b.need(ck, walk.Absent, nil, &walkNote{rng: rng, present: alreadyExistsIn("constraint", ident(written), relation(rv))}) // 42710
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

// mayIndex records the indexes and sequences that constraints of table t
// may or may not create, where a statement may skip them: a name they
// take may now be in use.
func (b *builder) mayIndex(t walk.Key, cons []*ast.Constraint) {
	s := walk.Path(schemaOfKey(t))
	for _, con := range cons {
		switch {
		case con.Contype == ast.CONSTR_IDENTITY && con.Options != nil && len(con.Options.Items) > 0:
			// The sequence may be named anywhere.
			b.mayAdd(walk.Scope{Space: spaceRelation, Deep: true})
		case con.Contype == ast.CONSTR_IDENTITY:
			b.mayAdd(walk.Scope{Space: spaceRelation, Parent: s, Prefix: generatedPrefix(t.Name)})
		case indexBacked(constraintKind(con.Contype)):
			if b.w.partitioned[t] {
				b.mayAdd(walk.Scope{Space: spaceRelation, Deep: true})
			}
			name, ok := walkName(con.Conname)
			switch {
			case ok:
				b.set(relKey(schemaOfKey(t), name), walk.Unknown, "")
			case con.Conname == "":
				b.mayAdd(walk.Scope{Space: spaceRelation, Parent: s, Prefix: generatedPrefix(t.Name)})
			default:
				b.mayAdd(walk.Scope{Space: spaceRelation, Parent: s})
			}
		}
	}
}

// createIndex follows CREATE INDEX. The index goes to its table's schema.
//
// pg: src/backend/commands/indexcmds.c — DefineIndex (an index only on a
// table, a partitioned table, or a materialized view, 42809)
func (b *builder) createIndex(v *ast.IndexStmt) bool {
	b.Last = b.w.callsUserCode(v)
	r := b.exists(v.Relation, false)
	if r.fact == walk.Present {
		b.need(r.key, walk.Present, []string{relTable, relMatview}, ofKind(rangeOf(v.Relation.Loc), relation(v.Relation), "table or materialized view"))
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
		// A relation of unknown kind may have replaced the partitioned
		// table: only a known table is judged.
		if v.Concurrent && r.kind == relTable {
			// Unsupported (0A000, verified on PostgreSQL 17 and 18 outside
			// a transaction, as Bytebase runs CONCURRENTLY).
			b.failAt(rangeOf(v.Relation.Loc), nil, review.UnsupportedForm, unsupported("`CREATE INDEX CONCURRENTLY` on partitioned table "+review.Code(relation(v.Relation))))
			return true
		}
		// Each partition gets an index of its own, in its own schema.
		b.mayAdd(walk.Scope{Space: spaceRelation, Deep: true})
	}
	if v.Idxname == "" {
		b.checkIndexElems(r, v)
		b.mayAdd(walk.Scope{Space: spaceRelation, Parent: s, Prefix: generatedPrefix(r.key.Name)})
		return true
	}
	name, ok := walkName(v.Idxname)
	if !ok {
		if !v.IfNotExists {
			b.checkIndexElems(r, v)
		}
		b.mayAdd(walk.Scope{Space: spaceRelation, Parent: s})
		return true
	}
	ik := relKey(schemaOfKey(r.key), name)
	prior, _ := b.l.Lookup(ik)
	if !v.IfNotExists {
		b.need(ik, walk.Absent, nil, taken(b.anchor(), relIndex, ident(v.Idxname))) // 42P07
		prior = walk.Absent
	}
	if prior == walk.Absent {
		b.checkIndexElems(r, v)
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
	if b.crossDatabase(rv) {
		return true
	}
	c := b.creationSchema(rv.Schemaname, rv.Relname, temp)
	if c.fail != "" {
		b.failAt(rng, c.keys, c.code, c.fail)
		return true
	}
	name, ok := walkName(rv.Relname)
	if c.fact != walk.Present || !ok {
		b.placeUnknown(rv.Relname, false)
		return true
	}
	k := relKey(c.schema, name)
	tk := typeKey(c.schema, name)
	prior, _ := b.l.Lookup(k)
	switch {
	case !v.Replace:
		b.need(k, walk.Absent, nil, taken(rng, relView, relation(rv)))                                  // 42P07
		b.need(tk, walk.Absent, nil, &walkNote{rng: rng, present: alreadyExists("type", relation(rv))}) // 42710
	case prior == walk.Present:
		// OR REPLACE replaces a view and nothing else (42809, verified on
		// PostgreSQL 17 and 18).
		b.need(k, walk.Present, []string{relView}, ofKind(rng, relation(rv), relView))
	case prior == walk.Absent:
		// It creates the view, whose row type takes the name (42710).
		b.need(tk, walk.Absent, nil, &walkNote{rng: rng, present: alreadyExists("type", relation(rv))})
	}
	b.checkViewQuery(v, k, prior)
	// OR REPLACE may create the view as well as replace it; holds asks
	// the ledger whether it is still there.
	b.made(k)
	b.set(k, walk.Present, relView)
	b.set(tk, walk.Present, typOther)
	b.mayAddArrayType(c.schema)
	b.forget(walk.Scope{Parent: tableOf(k), Deep: true})
	b.session(func() {
		b.w.modifiedViews[k] = true
		b.w.dirtyTables[k] = true
	})
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
	_, isExecute := query.(*ast.ExecuteStmt)
	if isExecute {
		// EXECUTE runs a prepared statement whose effects are not followed,
		// so end the walk after checking the target relation's existence.
		b.Last = true
	} else if !into.SkipData && b.w.callsUserCode(query) {
		b.Last = true
	}
	rv := into.Rel
	rng := rangeOf(rv.Loc)
	temp := rv.Relpersistence == ast.RELPERSISTENCE_TEMP
	if b.crossDatabase(rv) {
		return true
	}
	c := b.creationSchema(rv.Schemaname, rv.Relname, temp)
	if c.fail != "" {
		b.failAt(rng, c.keys, c.code, c.fail)
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
	existing := prior == walk.Present
	// PostgreSQL analyzes the query of CREATE TABLE/MATERIALIZED VIEW AS
	// before checking IF NOT EXISTS, so an unbound bare column in a
	// FROM-less SELECT fails even when the relation already exists (42703).
	if !isExecute && selectHasUnboundColumn(query) {
		b.fail(nil, walk.Absent, review.ObjectNotFound, "Query references a column with no `FROM` clause.") // 42703
		return true
	}
	if !isExecute && b.queryRefused(query) {
		return true
	}
	switch {
	case !ifNotExists:
		b.need(k, walk.Absent, nil, taken(rng, kind, relation(rv)))                                     // 42P07
		b.need(tk, walk.Absent, nil, &walkNote{rng: rng, present: alreadyExists("type", relation(rv))}) // 42710
		prior = walk.Absent
	case prior == walk.Present:
		return true
	case prior == walk.Absent:
		// IF NOT EXISTS looks at the relation alone (verified on
		// PostgreSQL 17).
		b.need(tk, walk.Absent, nil, &walkNote{rng: rng, present: alreadyExists("type", relation(rv))}) // 42710
	}
	if prior == walk.Absent {
		if into.AccessMethod == "" || into.AccessMethod == "heap" {
			b.checkHeapRelOptions(into.Options, rng)
		}
		if !temp && c.schema != tempSchema && into.OnCommit != ast.ONCOMMIT_NOOP && into.OnCommit != 0 {
			b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("`ON COMMIT` on a non-temporary table")) // 42P16
			return true
		}
		if into.ColNames != nil && len(into.ColNames.Items) > 0 {
			if sel, ok := query.(*ast.SelectStmt); ok && sel.Op == ast.SETOP_NONE && sel.ValuesLists == nil {
				if n, exact := selectTargetCount(sel); exact && len(into.ColNames.Items) > n {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("more column names than columns in `CREATE TABLE AS`")) // 42601
					return true
				}
			}
		}
		if isExecute && !existing {
			if ex, ok := query.(*ast.ExecuteStmt); ok && ex.Name != "" {
				b.fail(nil, walk.Absent, review.ObjectNotFound, notFound("prepared statement", ident(ex.Name))) // 26000
			}
		}
	}
	if isExecute {
		return true
	}
	parent := tableOf(k)
	if prior == walk.Absent {
		b.set(k, walk.Present, kind)
		b.set(tk, walk.Present, "")
		b.mayAddArrayType(c.schema)
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
	if b.crossDatabase(rv) {
		return true
	}
	c := b.creationSchema(rv.Schemaname, rv.Relname, temp)
	if c.fail != "" {
		b.failAt(rng, c.keys, c.code, c.fail)
		return true
	}
	name, ok := walkName(rv.Relname)
	if c.fact != walk.Present || !ok {
		b.placeUnknown(rv.Relname, false)
		return true
	}
	k := relKey(c.schema, name)
	tk := typeKey(c.schema, name)
	prior, _ := b.l.Lookup(k)
	switch {
	case !v.IfNotExists:
		b.need(k, walk.Absent, nil, taken(rng, relSequence, relation(rv)))                              // 42P07
		b.need(tk, walk.Absent, nil, &walkNote{rng: rng, present: alreadyExists("type", relation(rv))}) // 42710
		prior = walk.Absent
	case prior == walk.Absent:
		// IF NOT EXISTS skips a relation conflict, not a standalone type
		// conflict in heap_create_with_catalog (42710, verified on PG 17).
		b.need(tk, walk.Absent, nil, &walkNote{rng: rng, present: alreadyExists("type", relation(rv))})
	}
	if prior == walk.Absent {
		b.checkSequenceOptions(v.Options, false)
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
	if b.crossDatabase(rv) {
		return true
	}
	c := b.creationSchema(rv.Schemaname, rv.Relname, false)
	if c.fail != "" {
		b.failAt(rng, c.keys, c.code, c.fail)
		return true
	}
	name, ok := walkName(rv.Relname)
	if c.fact != walk.Present || !ok {
		b.placeUnknown(rv.Relname, false)
		return true
	}
	k := relKey(c.schema, name)
	tk := typeKey(c.schema, name)
	b.need(tk, walk.Absent, nil, &walkNote{rng: rng, present: alreadyExists("type", relation(rv))}) // 42710
	b.need(k, walk.Absent, nil, taken(rng, relComposite, relation(rv)))                             // 42P07
	if v.Coldeflist != nil {
		seen := make(map[string]bool, len(v.Coldeflist.Items))
		for _, item := range v.Coldeflist.Items {
			col, ok := item.(*ast.ColumnDef)
			if !ok {
				continue
			}
			cn, ok := storedKey(col.Colname)
			if !ok {
				continue
			}
			if seen[cn] {
				b.failAt(rangeOf(col.Loc), nil, review.DuplicateName, twice("declares", "column", ident(col.Colname))) // 42701
				return true
			}
			seen[cn] = true
			if col.TypeName != nil {
				if b.checkColumnType(col.TypeName, rangeOf(col.Loc)) {
					return true
				}
				if impossible, keys := b.selfTypeImpossible(c.schema, name, col.TypeName); impossible {
					b.failAt(rangeOf(col.Loc), keys, review.UnsupportedForm, "Composite type "+review.Code(relation(rv))+" can't reference itself.")
					return true
				}
			}
		}
	}
	b.set(k, walk.Present, relComposite)
	b.set(tk, walk.Present, typOther)
	b.mayAddArrayType(c.schema)
	b.forget(walk.Scope{Parent: tableOf(k), Deep: true})
	b.made(k)
	b.session(func() { b.w.dirtyTypes[k] = true })
	return true
}

// mayAddArrayType weakens Absent facts for PostgreSQL's auto-generated
// companion array types (_<name>) in schema s.
func (b *builder) mayAddArrayType(schema string) {
	b.mayAdd(walk.Scope{Space: spaceType, Parent: walk.Path(schema), Prefix: "_"})
}

// mayRemoveArrayType weakens Present facts for PostgreSQL's auto-generated
// companion array types (_<name>) in schema s.
func (b *builder) mayRemoveArrayType(schema string) {
	b.mayRemove(walk.Scope{Space: spaceType, Parent: walk.Path(schema), Prefix: "_"})
}

// createEnum follows CREATE TYPE ... AS ENUM (...).
func (b *builder) createEnum(v *ast.CreateEnumStmt) bool {
	if !b.createType(nameParts(v.TypeName), typOther) {
		return false
	}
	if v.Vals != nil {
		seen := make(map[string]bool, len(v.Vals.Items))
		for _, item := range v.Vals.Items {
			s, ok := item.(*ast.String)
			if !ok {
				continue
			}
			if seen[s.Str] {
				b.fail(nil, walk.Present, review.DuplicateName, "This statement declares enum label "+review.Code("'"+s.Str+"'")+" twice.") // 23505
				return true
			}
			seen[s.Str] = true
		}
	}
	return true
}

// createRange follows CREATE TYPE ... AS RANGE (...). A range definition may
// complete a shell type created earlier in the change (pg: DefineRange).
func (b *builder) createRange(v *ast.CreateRangeStmt) bool {
	parts := nameParts(v.TypeName)
	if len(parts) == 3 {
		if b.w.crossDatabase(parts[0]) {
			b.fail(nil, walk.Absent, review.CrossDatabaseReference, crossDatabaseMessage(qualified(parts)))
			return true
		}
		if !b.w.currentDatabase(parts[0]) {
			return false
		}
		parts = parts[1:]
	}
	schema, written, ok := splitName(parts)
	if !ok {
		return false
	}
	c := b.creationSchema(schema, written, false)
	if c.fail != "" {
		b.fail(c.keys, walk.Absent, c.code, c.fail)
		return true
	}
	name, ok := walkName(written)
	if c.fact != walk.Present || !ok {
		b.placeUnknown(written, false)
		return true
	}
	tk := typeKey(c.schema, name)
	if f, _ := b.l.Lookup(tk); f != walk.Present || !b.w.shellTypes[tk] {
		b.need(tk, walk.Absent, nil, &walkNote{rng: b.anchor(), present: alreadyExists("type", qualified(parts))}) // 42710
	}
	hasSubtype := false
	var multirangeArg ast.Node
	if v.Params != nil {
		for _, item := range v.Params.Items {
			de, ok := item.(*ast.DefElem)
			if !ok {
				continue
			}
			switch strings.ToLower(de.Defname) {
			case "subtype":
				hasSubtype = true
			case "multirange_type_name":
				multirangeArg = de.Arg
			}
		}
	}
	if !hasSubtype {
		b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("`CREATE TYPE ... AS RANGE` without `SUBTYPE`")) // 42601
		return true
	}
	b.set(tk, walk.Present, typOther)
	b.mayAddArrayType(c.schema)
	b.session(func() { delete(b.w.shellTypes, tk) })
	b.made(tk)
	// PostgreSQL 14+ also creates a companion multirange type alongside the
	// range type (pg: src/backend/commands/typecmds.c — DefineRange), and
	// range / multirange constructor functions in pg_proc. Weaken their target
	// scopes so later references are not falsely reported absent.
	b.mayAdd(walk.Scope{Space: spaceFunction, Parent: walk.Path(c.schema)})
	b.session(func() { b.w.allDirty = true })
	if multirangeArg != nil {
		var mrParts []string
		switch arg := multirangeArg.(type) {
		case *ast.TypeName:
			mrParts = nameParts(arg.Names)
		case *ast.List:
			mrParts = nameParts(arg)
		case *ast.String:
			if arg.Str != "" {
				mrParts = []string{arg.Str}
			}
		}
		switch len(mrParts) {
		case 1:
			if mn, ok := walkName(mrParts[0]); ok {
				b.mayAdd(walk.Scope{Space: spaceType, Parent: walk.Path(c.schema), Prefix: mn})
			} else {
				b.mayAdd(walk.Scope{Space: spaceType, Parent: walk.Path(c.schema)})
			}
		case 2:
			ms, ok1 := walkName(mrParts[0])
			mn, ok2 := walkName(mrParts[1])
			if ok1 && ok2 {
				b.mayAdd(walk.Scope{Space: spaceType, Parent: walk.Path(ms), Prefix: mn})
				b.mayAdd(walk.Scope{Space: spaceFunction, Parent: walk.Path(ms)})
			} else {
				b.mayAdd(walk.Scope{Space: spaceType, Deep: true})
				b.mayAdd(walk.Scope{Space: spaceFunction, Deep: true})
			}
		default:
			b.mayAdd(walk.Scope{Space: spaceType, Deep: true})
			b.mayAdd(walk.Scope{Space: spaceFunction, Deep: true})
		}
	} else {
		b.mayAdd(walk.Scope{Space: spaceType, Parent: walk.Path(c.schema)})
	}
	return true
}

// defineObject follows CREATE TYPE name (shell type) and CREATE TYPE name (...)
// (base type). Other DefineStmt kinds end the walk.
func (b *builder) defineObject(v *ast.DefineStmt) bool {
	if v.Kind != ast.OBJECT_TYPE {
		return false
	}
	parts := nameParts(v.Defnames)
	if len(parts) > 2 {
		if len(parts) == 3 && b.w.crossDatabase(parts[0]) {
			b.fail(nil, walk.Absent, review.CrossDatabaseReference, crossDatabaseMessage(qualified(parts))) // 0A000
			return true
		}
		if len(parts) > 3 {
			b.fail(nil, walk.Absent, review.CrossDatabaseReference, crossDatabaseMessage(qualified(parts))) // 0A000
			return true
		}
		if !b.w.currentDatabase(parts[0]) {
			return false
		}
		parts = parts[1:]
	}
	schema, written, ok := splitName(parts)
	if !ok {
		return false
	}
	c := b.creationSchema(schema, written, false)
	if c.fail != "" {
		b.fail(c.keys, walk.Absent, c.code, c.fail) // 3F000
		return true
	}
	name, ok := walkName(written)
	if c.fact != walk.Present || !ok {
		b.placeUnknown(written, false)
		return true
	}
	tk := typeKey(c.schema, name)
	isShell := v.Definition == nil || len(v.Definition.Items) == 0
	if isShell {
		b.need(tk, walk.Absent, nil, &walkNote{rng: b.anchor(), present: alreadyExists("type", qualified(parts))}) // 42710
		b.set(tk, walk.Present, typOther)
		b.session(func() { b.w.shellTypes[tk] = true })
		b.made(tk)
		return true
	}
	if f, _ := b.l.Lookup(tk); f == walk.Present && !b.w.shellTypes[tk] {
		b.need(tk, walk.Absent, nil, &walkNote{rng: b.anchor(), present: alreadyExists("type", qualified(parts))}) // 42710
		return true
	}
	hasInput, hasOutput, hasTypmodIn := false, false, false
	for _, item := range v.Definition.Items {
		de, ok := item.(*ast.DefElem)
		if !ok {
			continue
		}
		switch strings.ToLower(de.Defname) {
		case "input":
			hasInput = true
		case "output":
			hasOutput = true
		case "typmod_in":
			hasTypmodIn = true
		}
	}
	if !hasInput || !hasOutput {
		b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("`CREATE TYPE` without `INPUT` and `OUTPUT` functions")) // 42P17 / 42710
		return true
	}
	b.set(tk, walk.Present, typOther)
	b.mayAddArrayType(c.schema)
	b.session(func() {
		delete(b.w.shellTypes, tk)
		if hasTypmodIn {
			b.w.typmodTypes[tk] = true
		} else {
			delete(b.w.typmodTypes, tk)
		}
	})
	b.made(tk)
	return true
}

// checkColumnType checks that a referenced column/attribute/domain type is not
// certainly absent or cross-database (0A000 / 3F000 / 42704), returning true
// when it emitted a failure.
func (b *builder) checkColumnType(tn *ast.TypeName, rng review.Range) bool {
	if tn == nil {
		return false
	}
	parts := nameParts(tn.Names)
	if len(parts) == 0 {
		return false
	}
	r := b.resolveType(parts)
	switch {
	case r.fail != "":
		b.failAt(rng, nil, review.CrossDatabaseReference, r.fail)
		return true
	case r.fact == walk.Absent:
		if r.missingSchema != "" {
			b.failAt(rng, r.keys, review.ObjectNotFound, schemaNotFound(r.missingSchema))
		} else {
			b.failAt(rng, r.keys, review.ObjectNotFound, notFound("type", typeName(tn)))
		}
		return true
	}
	return false
}

// checkCompositeType checks that tn is not certainly a non-composite type
// (42809), used by CREATE TABLE ... OF and ALTER TABLE ... OF. Returns true
// when it emitted a failure.
func (b *builder) checkCompositeType(tn *ast.TypeName, rng review.Range) bool {
	if tn == nil {
		return false
	}
	parts := nameParts(tn.Names)
	if len(parts) == 0 {
		return false
	}
	r := b.resolveType(parts)
	if r.fail != "" {
		b.failAt(rng, nil, review.CrossDatabaseReference, r.fail)
		return true
	}
	if r.fact == walk.Absent {
		if r.missingSchema != "" {
			b.failAt(rng, r.keys, review.ObjectNotFound, schemaNotFound(r.missingSchema))
		} else {
			b.failAt(rng, r.keys, review.ObjectNotFound, notFound("type", typeName(tn)))
		}
		return true
	}
	if r.fact != walk.Present {
		return false
	}
	if tn.ArrayBounds != nil || r.builtin || r.kind == typDomain {
		b.failAt(rng, r.keys, review.WrongObjectKind, review.Code(typeName(tn))+" isn't a composite type.") // 42809
		return true
	}
	rk := relKey(schemaOfKey(r.key), r.key.Name)
	rf, rkind := b.l.Lookup(rk)
	switch {
	case rf == walk.Present && rkind != "" && rkind != relComposite:
		b.failAt(rng, []walk.Key{rk}, review.WrongObjectKind, review.Code(typeName(tn))+" is "+article(rkind)+", not a composite type.") // 42809
		return true
	case rf == walk.Absent:
		// Present in spaceType with no backing pg_class entry (e.g. an enum).
		b.failAt(rng, []walk.Key{r.key, rk}, review.WrongObjectKind, review.Code(typeName(tn))+" isn't a composite type.") // 42809
		return true
	}
	return false
}

// checkDomain checks that a type named by parts is not certainly absent or a
// non-domain type (42809).
func (b *builder) checkDomain(parts []string, rng review.Range) resolved {
	r := b.resolveType(parts)
	written := qualified(parts)
	switch {
	case r.fail != "":
		b.failAt(rng, nil, review.CrossDatabaseReference, r.fail)
	case r.fact == walk.Absent:
		if r.missingSchema != "" {
			b.failAt(rng, r.keys, review.ObjectNotFound, schemaNotFound(r.missingSchema))
		} else {
			b.failAt(rng, r.keys, review.ObjectNotFound, notFound("type", written))
		}
	case r.fact == walk.Present:
		if r.builtin {
			b.failAt(rng, r.keys, review.WrongObjectKind, review.Code(written)+" isn't a domain.") // 42809
		} else {
			b.need(r.key, walk.Present, []string{typDomain}, &walkNote{rng: rng, present: review.Code(written) + " isn't a domain."}) // 42809
		}
	}
	return r
}

// alterDomain follows ALTER DOMAIN. ADD CONSTRAINT ... CHECK and VALIDATE
// CONSTRAINT validate existing domain values and may execute user functions
// that change the schema, so the walk ends after them.
func (b *builder) alterDomain(v *ast.AlterDomainStmt) bool {
	if v.Subtype == 'V' || b.w.callsUserCode(v) {
		b.Last = true
	}
	parts := nameParts(v.Typname)
	if len(parts) == 0 {
		return true
	}
	r := b.checkDomain(parts, b.anchor())
	if r.fail != "" || r.fact == walk.Absent || r.builtin {
		return true
	}
	// Law of excluded middle for domain DEFAULT and CHECK expressions:
	// - SET DEFAULT <bare_col>: column references are never allowed in DEFAULT
	//   (0A000 if the domain exists, 42704 if it does not).
	// - ADD CHECK (<bare_col>): only VALUE is valid as a bare identifier in a
	//   domain CHECK expression (42703 if the domain exists, 42704 if not).
	switch v.Subtype {
	case 'T':
		if _, ok := bareColumnRef(v.Def); ok {
			b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("column reference in domain `DEFAULT`")) // 0A000 / 42704
		}
	case 'C':
		if con, ok := v.Def.(*ast.Constraint); ok && con.Contype == ast.CONSTR_CHECK {
			if col, ok := bareColumnRef(con.RawExpr); ok && !strings.EqualFold(col, "value") {
				b.fail(nil, walk.Absent, review.ObjectNotFound, unsupported("column reference in domain `CHECK`")) // 42703 / 42704
			}
		}
	}
	return true
}

// alterEnum follows ALTER TYPE ... ADD VALUE / RENAME VALUE.
func (b *builder) alterEnum(v *ast.AlterEnumStmt) bool {
	parts := nameParts(v.Typname)
	if len(parts) == 0 {
		return true
	}
	written := qualified(parts)
	rng := b.anchor()
	r := b.resolveType(parts)
	switch {
	case r.fail != "":
		b.failAt(rng, nil, review.CrossDatabaseReference, r.fail)
		return true
	case r.fact == walk.Absent:
		if r.missingSchema != "" {
			b.failAt(rng, r.keys, review.ObjectNotFound, schemaNotFound(r.missingSchema))
		} else {
			b.failAt(rng, r.keys, review.ObjectNotFound, notFound("type", written))
		}
		return true
	case r.fact == walk.Present:
		if r.builtin || r.kind == typDomain {
			b.failAt(rng, r.keys, review.WrongObjectKind, review.Code(written)+" isn't an enum.") // 42809
			return true
		}
		rk := relKey(schemaOfKey(r.key), r.key.Name)
		if rf, rkind := b.l.Lookup(rk); rf == walk.Present && hasRowType(rkind) {
			b.failAt(rng, []walk.Key{rk}, review.WrongObjectKind, review.Code(written)+" is "+article(rkind)+", not an enum.") // 42809
			return true
		}
	}
	// ADD VALUE 'x' BEFORE/AFTER 'x' (without IF NOT EXISTS) and RENAME VALUE
	// 'x' TO 'x' cannot succeed: if 'x' exists it fails with 42710, and if it
	// does not exist it fails with 22023 (or 42704/42809 on the type).
	if !v.SkipIfNewvalExists && v.Oldval == "" && v.NewvalNeighbor != "" && v.Newval == v.NewvalNeighbor {
		b.failAt(rng, nil, review.UnsupportedForm, unsupported("adding enum label "+review.Code("'"+v.Newval+"'")+" relative to itself"))
		return true
	}
	if v.Oldval != "" && v.Oldval == v.Newval {
		b.failAt(rng, nil, review.UnsupportedForm, unsupported("renaming enum label "+review.Code("'"+v.Oldval+"'")+" to itself"))
		return true
	}
	return true
}

// validAlterTypeProps lists the base-type properties ALTER TYPE ... SET (...)
// accepts in PostgreSQL.
//
// pg: src/backend/commands/typecmds.c — AlterType
var validAlterTypeProps = map[string]bool{
	"receive": true, "send": true, "typmod_in": true, "typmod_out": true,
	"analyze": true, "analyse": true, "subscript": true, "storage": true,
}

// alterType follows ALTER TYPE ... SET (...).
func (b *builder) alterType(v *ast.AlterTypeStmt) bool {
	parts := nameParts(v.TypeName)
	var r resolved
	if len(parts) > 0 {
		r = b.resolveType(parts)
		switch {
		case r.fail != "":
			b.fail(nil, walk.Absent, review.CrossDatabaseReference, r.fail)
			return true
		case r.fact == walk.Absent:
			if r.missingSchema != "" {
				b.fail(r.keys, walk.Absent, review.ObjectNotFound, schemaNotFound(r.missingSchema))
			} else {
				b.fail(r.keys, walk.Absent, review.ObjectNotFound, notFound("type", qualified(parts)))
			}
			return true
		}
	}
	setTypmodIn := false
	if v.Options != nil {
		for _, item := range v.Options.Items {
			de, ok := item.(*ast.DefElem)
			if !ok {
				continue
			}
			prop := strings.ToLower(de.Defname)
			if !validAlterTypeProps[prop] {
				b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("type attribute "+review.Code(ident(de.Defname))+" in `ALTER TYPE ... SET`")) // 42601
				return true
			}
			if prop == "typmod_in" {
				setTypmodIn = true
			}
		}
	}
	if setTypmodIn {
		switch {
		case r.fact == walk.Present && !r.builtin:
			b.session(func() { b.w.typmodTypes[r.key] = true })
		case r.fact == walk.Unknown:
			b.session(func() { b.w.typmodTypes[walk.Key{}] = true })
		}
	}
	return true
}

// createDomain follows CREATE DOMAIN, checking the creation schema, target
// type collision, underlying type existence, and impossible self-references.
func (b *builder) createDomain(v *ast.CreateDomainStmt) bool {
	parts := nameParts(v.Domainname)
	if len(parts) == 3 {
		if b.w.crossDatabase(parts[0]) {
			b.fail(nil, walk.Absent, review.CrossDatabaseReference, crossDatabaseMessage(qualified(parts))) // 0A000
			return true
		}
		if !b.w.currentDatabase(parts[0]) {
			return false
		}
		parts = parts[1:]
	}
	schema, written, ok := splitName(parts)
	if !ok {
		return false
	}
	c := b.creationSchema(schema, written, false)
	if c.fail != "" {
		b.fail(c.keys, walk.Absent, c.code, c.fail)
		return true
	}
	name, ok := walkName(written)
	if c.fact != walk.Present || !ok {
		b.placeUnknown(written, false)
		return true
	}
	tk := typeKey(c.schema, name)
	b.need(tk, walk.Absent, nil, &walkNote{rng: b.anchor(), present: alreadyExists("type", qualified(parts))}) // 42710
	if v.Typname != nil {
		if b.checkColumnType(v.Typname, rangeOf(v.Typname.Loc)) {
			return true
		}
		if f, _ := b.l.Lookup(tk); f == walk.Unknown {
			if impossible, keys := b.selfTypeImpossible(c.schema, name, v.Typname); impossible {
				b.fail(keys, walk.Absent, review.UnsupportedForm, "Domain "+review.Code(qualified(parts))+" can't use itself as its underlying type.")
				return true
			}
		}
	}
	b.set(tk, walk.Present, typDomain)
	b.mayAddArrayType(c.schema)
	b.made(tk)
	return true
}

// createType follows CREATE TYPE ... AS ENUM, which takes a name in the type
// namespace alone.
func (b *builder) createType(parts []string, kind string) bool {
	schema, written, ok := splitName(parts)
	if !ok {
		return false
	}
	c := b.creationSchema(schema, written, false)
	if c.fail != "" {
		b.fail(c.keys, walk.Absent, c.code, c.fail)
		return true
	}
	name, ok := walkName(written)
	if c.fact != walk.Present || !ok {
		b.placeUnknown(written, false)
		return true
	}
	tk := typeKey(c.schema, name)
	b.need(tk, walk.Absent, nil, &walkNote{rng: b.anchor(), present: alreadyExists("type", qualified(parts))}) // 42710
	b.set(tk, walk.Present, kind)
	b.mayAddArrayType(c.schema)
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

// alterForm is the form of ALTER of a relation wrongKind judges.
type alterForm int

const (
	formAlter alterForm = iota
	formRename
	formSetSchema
)

// anyRelation lists every relation kind the walk keeps.
var anyRelation = []string{relTable, relView, relMatview, relSequence, relIndex, relForeign, relComposite}

// wrongKind checks that an ALTER of a relation names one of a kind it
// takes, and reports whether the relation is certainly of another (42809,
// verified on PostgreSQL 17 and 18, IF EXISTS or not). ALTER VIEW,
// MATERIALIZED VIEW, SEQUENCE, and FOREIGN TABLE take their own kind
// only, and ALTER INDEX an index only but for RENAME. ALTER TYPE's
// attribute forms, which parse as an ALTER of a relation, take a composite
// type only. ALTER TABLE takes any relation but a composite type, for
// compatibility, and moves no index to another schema.
//
// pg: src/backend/commands/tablecmds.c — RangeVarCallbackForAlterRelation,
// AlterTableNamespace
func (b *builder) wrongKind(r resolved, rv *ast.RangeVar, obj ast.ObjectType, form alterForm) bool {
	if r.fact != walk.Present || r.kind == "" {
		return false
	}
	// The kind the statement names, which the text says the relation is
	// not.
	want, own := dropKinds[obj]
	needed := want.word
	if !own {
		needed = relTable
	}
	var kinds []string
	switch {
	case obj == ast.OBJECT_TYPE:
		kinds = []string{relComposite}
		needed = relComposite
	case own && obj != ast.OBJECT_TABLE && (obj != ast.OBJECT_INDEX || form != formRename):
		kinds = []string{want.kind}
	case r.kind == relComposite:
		kinds = slices.DeleteFunc(slices.Clone(anyRelation), func(k string) bool { return k == relComposite })
	case r.kind == relIndex && form == formSetSchema:
		kinds = slices.DeleteFunc(slices.Clone(anyRelation), func(k string) bool { return k == relComposite || k == relIndex })
	default:
		return false
	}
	if slices.Contains(kinds, r.kind) {
		return false
	}
	b.need(r.key, walk.Present, kinds, ofKind(rangeOf(rv.Loc), relation(rv), needed))
	return true
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
	b.Last = b.w.callsUserCode(v)
	r := b.exists(v.Relation, v.Missing_ok)
	objType := ast.ObjectType(v.ObjType)
	if b.wrongKind(r, v.Relation, objType, formAlter) {
		return true
	}
	if r.fact == walk.Present && r.kind != "" && objType != ast.OBJECT_TYPE {
		for _, cmd := range cmds {
			sub := ast.AlterTableType(cmd.Subtype)
			if allowed, word := atAllowedKinds(sub); allowed != nil && !slices.Contains(allowed, r.kind) {
				b.need(r.key, walk.Present, allowed, ofKind(rangeOf(v.Relation.Loc), relation(v.Relation), word)) // 42809
				return true
			}
		}
	}
	switch {
	case r.fail != "":
		return true
	case r.fact == walk.Absent:
		return true // a notice under IF EXISTS
	case r.fact == walk.Unknown, objType != ast.OBJECT_TABLE && objType != ast.OBJECT_FOREIGN_TABLE, !hasColumns(r.kind):
		if r.fact == walk.Present && objType == ast.OBJECT_TYPE {
			droppedAttrs := make(map[string]bool)
			for _, cmd := range cmds {
				if ast.AlterTableType(cmd.Subtype) == ast.AT_DropColumn {
					if n, ok := storedKey(cmd.Name); ok {
						if droppedAttrs[n] && !cmd.Missing_ok {
							b.failAt(rangeOf(cmd.Loc), nil, review.DuplicateName, twice("drops", "attribute", ident(cmd.Name))) // 42703
							return true
						}
						droppedAttrs[n] = true
					}
				}
			}
			for _, cmd := range cmds {
				switch ast.AlterTableType(cmd.Subtype) {
				case ast.AT_AddColumn:
					if def, ok := cmd.Def.(*ast.ColumnDef); ok && !b.w.isTypeDirty(r.key) {
						if an, ok := walkName(def.Colname); ok && !droppedAttrs[an] {
							if b.w.pg.compAttrs[[2]string{schemaOfKey(r.key), r.key.Name}][an] {
								b.failAt(rangeOf(cmd.Loc), []walk.Key{r.key}, review.DuplicateName, alreadyExistsIn("attribute", ident(def.Colname), relation(v.Relation))) // 42701
								return true
							}
						}
					}
				case ast.AT_AlterColumnType:
					if def, ok := cmd.Def.(*ast.ColumnDef); ok && def.TypeName != nil && def.TypeName.ArrayBounds == nil {
						if tr := b.resolveType(nameParts(def.TypeName.Names)); tr.fact == walk.Present && !tr.builtin && tr.key == typeKey(schemaOfKey(r.key), r.key.Name) {
							b.failAt(rangeOf(cmd.Loc), []walk.Key{r.key}, review.UnsupportedForm, "Composite type "+review.Code(relation(v.Relation))+" can't reference itself.") // 42P16
							return true
						}
					}
				}
			}
			b.session(func() { b.w.dirtyTypes[r.key] = true })
		}
		// Where the relation is, or what the subcommands do to a relation
		// of its kind, is not followed: check known target-index/type
		// subcommands on present relations, and pass only no-fact ones.
		if r.fact == walk.Present && r.kind != "" && objType != ast.OBJECT_TYPE {
			for _, cmd := range cmds {
				b.checkAlterTargetSubcmd(v.Relation, r, cmd)
			}
		}
		for _, cmd := range cmds {
			if !atNoFacts[ast.AlterTableType(cmd.Subtype)] {
				if r.fact == walk.Present && r.kind != "" && objType != ast.OBJECT_TYPE {
					b.Last = true
					return true
				}
				return false
			}
		}
		return true
	}
	b.repeated(v.Relation, cmds)
	b.drops = make(map[string]bool)
	b.added = make(map[string]bool)
	b.colTouch = make(map[string]int)
	b.conDropped = false
	b.pkChanged = false
	for _, cmd := range cmds {
		sub := ast.AlterTableType(cmd.Subtype)
		switch sub {
		case ast.AT_DropColumn:
			if k, ok := storedKey(cmd.Name); ok {
				b.drops[k] = true
				b.colTouch[k]++
			}
		case ast.AT_AddColumn:
			if def, ok := cmd.Def.(*ast.ColumnDef); ok {
				if k, ok := storedKey(def.Colname); ok {
					b.colTouch[k]++
				}
			}
		case ast.AT_AlterColumnType, ast.AT_SetNotNull, ast.AT_DropNotNull,
			ast.AT_AddIdentity, ast.AT_DropIdentity, ast.AT_DropExpression:
			if k, ok := storedKey(cmd.Name); ok {
				b.colTouch[k]++
			}
		case ast.AT_DropConstraint:
			b.conDropped = true
			b.pkChanged = true
		case ast.AT_AddConstraint:
			if con, ok := cmd.Def.(*ast.Constraint); ok && con.Contype == ast.CONSTR_PRIMARY {
				b.pkChanged = true
				if con.Keys != nil {
					for _, item := range con.Keys.Items {
						if s, ok := item.(*ast.String); ok {
							if k, ok := storedKey(s.Str); ok {
								b.colTouch[k]++
							}
						}
					}
				}
			}
		}
	}
	for _, cmd := range passOrder(cmds) {
		b.next()
		if !b.alterTableCmd(v.Relation, r.key, cmd) {
			return false
		}
	}
	b.session(func() {
		if v.Relation.Inh {
			b.w.allDirty = true
		} else {
			b.w.dirtyTables[r.key] = true
		}
	})
	return true
}

// passOrder returns the subcommands in PostgreSQL's pass order: drops first,
// then column additions, then everything else, preserving textual order within
// each pass. PostgreSQL runs every DROP COLUMN and DROP CONSTRAINT before
// anything else, and ADD COLUMN before ADD CONSTRAINT, so ADD CONSTRAINT
// uq UNIQUE (b), ADD COLUMN b int succeeds (verified on PostgreSQL 17 and 18).
//
// pg: src/backend/commands/tablecmds.c — AlterTablePass (AT_PASS_DROP, AT_PASS_ADD_COL)
func passOrder(cmds []*ast.AlterTableCmd) []*ast.AlterTableCmd {
	out := make([]*ast.AlterTableCmd, 0, len(cmds))
	pass := func(cmd *ast.AlterTableCmd) int {
		switch ast.AlterTableType(cmd.Subtype) {
		case ast.AT_DropColumn, ast.AT_DropConstraint:
			return 0
		case ast.AT_AddColumn:
			return 1
		default:
			return 2
		}
	}
	for p := range 3 {
		for _, cmd := range cmds {
			if pass(cmd) == p {
				out = append(out, cmd)
			}
		}
	}
	return out
}

// repeated fails an ALTER TABLE that adds, or drops, one column twice,
// drops one constraint twice, or adds two primary keys or two
// constraints of one name, whatever the table held (verified on PostgreSQL 17 and 18):
//
//   - A second primary key fails as the subcommands are prepared, before
//     any runs (42P16).
//   - Subcommands of one kind run in the order written, so the second
//     column, or the second drop of a constraint, fails unless it says IF
//     NOT EXISTS or IF EXISTS (42701, 42703, 42704). Names compare as
//     storedKey keeps them. A system column's name fails for being one,
//     which its subcommand reports. Columns are added before constraints.
//   - Two keys collide on their indexes' name (42P07), and any other pair
//     on the constraint's (42710), whichever pass adds each.
//
// A NOT NULL constraint is left out, as in duplicateConstraint, and so
// are an index USING INDEX attaches and the constraints of a column ADD
// COLUMN IF NOT EXISTS may skip.
//
// pg: src/backend/parser/parse_utilcmd.c — transformIndexConstraints
func (b *builder) repeated(rv *ast.RangeVar, cmds []*ast.AlterTableCmd) {
	var cons []*ast.Constraint
	for _, cmd := range cmds {
		switch ast.AlterTableType(cmd.Subtype) {
		case ast.AT_AddConstraint:
			if con, ok := cmd.Def.(*ast.Constraint); ok {
				cons = append(cons, con)
			}
		case ast.AT_AddColumn:
			if def, ok := cmd.Def.(*ast.ColumnDef); ok && !cmd.Missing_ok && def.Constraints != nil {
				for _, item := range def.Constraints.Items {
					if con, ok := item.(*ast.Constraint); ok {
						cons = append(cons, con)
					}
				}
			}
		}
	}
	primaries := 0
	for _, con := range cons {
		if con.Contype == ast.CONSTR_PRIMARY {
			primaries++
			if primaries > 1 {
				b.failAt(rangeOf(con.Loc), nil, review.NameTaken, hasPrimaryKey(relation(rv)))
				return
			}
		}
	}
	added := make(map[string]bool)
	dropped := make(map[string]bool)
	droppedCons := make(map[string]bool)
	for _, cmd := range cmds {
		switch ast.AlterTableType(cmd.Subtype) {
		case ast.AT_AddColumn:
			def, ok := cmd.Def.(*ast.ColumnDef)
			if !ok {
				continue
			}
			n, ok := storedKey(def.Colname)
			if !ok || systemColumns[n] {
				continue
			}
			if added[n] && !cmd.Missing_ok {
				b.failAt(rangeOf(def.Loc), nil, review.DuplicateName, twice("adds", "column", ident(def.Colname)))
				return
			}
			added[n] = true
		case ast.AT_DropColumn:
			n, ok := storedKey(cmd.Name)
			if !ok || systemColumns[n] {
				continue
			}
			if dropped[n] && !cmd.Missing_ok {
				b.failAt(rangeOf(cmd.Loc), nil, review.DuplicateName, twice("drops", "column", ident(cmd.Name)))
				return
			}
			dropped[n] = true
		case ast.AT_DropConstraint:
			n, ok := storedKey(cmd.Name)
			if !ok {
				continue
			}
			if droppedCons[n] && !cmd.Missing_ok {
				b.failAt(rangeOf(cmd.Loc), nil, review.DuplicateName, twice("drops", "constraint", ident(cmd.Name)))
				return
			}
			droppedCons[n] = true
		}
	}
	named := make(map[string]bool)
	for _, con := range cons {
		key, ok := storedKey(con.Conname)
		if constraintKind(con.Contype) == "" || !ok || con.Indexname != "" {
			continue
		}
		if !named[key] {
			named[key] = true
			continue
		}
		b.failAt(rangeOf(con.Loc), nil, review.DuplicateName, twice("adds", "constraint", ident(con.Conname)))
		return
	}
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
	case sub == ast.AT_ClusterOn, sub == ast.AT_DropCluster, sub == ast.AT_ReplicaIdentity, sub == ast.AT_AddOf, sub == ast.AT_DropOf,
		sub == ast.AT_AddInherit, sub == ast.AT_DropInherit,
		sub == ast.AT_AttachPartition, sub == ast.AT_DetachPartition, sub == ast.AT_DetachPartitionFinalize,
		sub == ast.AT_SetRelOptions, sub == ast.AT_ReplaceRelOptions,
		sub == ast.AT_EnableTrig, sub == ast.AT_EnableAlwaysTrig, sub == ast.AT_EnableReplicaTrig, sub == ast.AT_DisableTrig,
		sub == ast.AT_EnableRule, sub == ast.AT_EnableAlwaysRule, sub == ast.AT_EnableReplicaRule, sub == ast.AT_DisableRule:
		b.checkAlterTargetSubcmd(rv, resolved{fact: walk.Present, key: t, kind: relTable}, cmd)
		if !atNoFacts[sub] {
			b.Last = true
		}
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
		b.need(constraintKey(t, name), walk.Present, nil, &walkNote{rng: rng, absent: notFoundIn("constraint", ident(written), relation(rv))})
	}
}

// addColumn follows ADD COLUMN. A system column's name fails even with IF
// NOT EXISTS, and the column reaches the inheritance children.
//
// pg: src/backend/commands/tablecmds.c — check_for_column_name_collision
func (b *builder) addColumn(rv *ast.RangeVar, t walk.Key, def *ast.ColumnDef, ifNotExists, inh bool) {
	rng := rangeOf(def.Loc)
	// A column the statement drops is gone when this one is added, so IF
	// NOT EXISTS adds it (verified on PostgreSQL 17 and 18).
	dropped := false
	if k, ok := storedKey(def.Colname); ok {
		if ifNotExists && b.added[k] {
			// An earlier subcommand added it: this one is a notice
			// (verified on PostgreSQL 17 and 18), which still adds the
			// column's constraints on PostgreSQL 12.
			b.mayConstrain(t, def, inh)
			return
		}
		b.added[k] = true
		dropped = b.drops[k]
	}
	name, ok := walkName(def.Colname)
	if !ok {
		b.mayAdd(walk.Scope{Space: spaceColumn, Parent: tableOf(t)})
		if inh {
			b.mayAdd(walk.Scope{Space: spaceColumn, Deep: true})
		}
		if ifNotExists && !dropped {
			b.mayConstrain(t, def, inh)
			return
		}
		b.columnConstraints(t, rv, def, inh)
		return
	}
	if systemColumns[name] {
		b.failAt(rng, nil, review.ReservedName, systemColumnName(def.Colname)) // 42701
		return
	}
	ck := columnKey(t, name)
	prior, _ := b.l.Lookup(ck)
	switch {
	case !ifNotExists:
		b.need(ck, walk.Absent, nil, &walkNote{rng: rng, present: alreadyExistsIn("column", ident(def.Colname), relation(rv))}) // 42701
		prior = walk.Absent
	case dropped:
		prior = walk.Absent
	}
	if prior == walk.Present {
		// A notice: the subcommand adds no column, and on PostgreSQL 12
		// still adds the column's constraints (verified on 12, 17, and
		// 18).
		b.mayConstrain(t, def, inh)
		return
	}
	if prior == walk.Absent && def.TypeName != nil && def.TypeName.ArrayBounds == nil {
		if tr := b.resolveType(nameParts(def.TypeName.Names)); tr.fact == walk.Present && !tr.builtin && tr.key == typeKey(schemaOfKey(t), t.Name) {
			b.failAt(rng, []walk.Key{t}, review.UnsupportedForm, unsupported("column of table "+review.Code(relation(rv))+"'s own row type")) // 42P16
			return
		}
	}
	if inh {
		b.mayAdd(walk.Scope{Space: spaceColumn, Deep: true, Prefix: name})
	}
	b.set(ck, walk.Present, kindColumn)
	if prior == walk.Unknown {
		// The column may be there already and the subcommand a notice:
		// its constraints may or may not be added, and nothing they need
		// is checked.
		b.mayConstrain(t, def, inh)
		return
	}
	b.columnConstraints(t, rv, def, inh)
}

// mayConstrain records the constraints of a column ADD COLUMN IF NOT
// EXISTS may or may not add.
func (b *builder) mayConstrain(t walk.Key, def *ast.ColumnDef, inh bool) {
	if isSerial(def.TypeName) {
		b.mayAdd(walk.Scope{Space: spaceRelation, Parent: walk.Path(schemaOfKey(t)), Prefix: generatedPrefix(t.Name)})
	}
	if def.Constraints == nil || len(def.Constraints.Items) == 0 {
		return
	}
	var cons []*ast.Constraint
	for _, item := range def.Constraints.Items {
		if con, ok := item.(*ast.Constraint); ok {
			cons = append(cons, con)
			if con.Contype == ast.CONSTR_PRIMARY {
				b.set(primaryKeyOf(t), walk.Unknown, "")
			}
		}
	}
	b.mayIndex(t, cons)
	b.mayAdd(walk.Scope{Space: spaceConstraint, Parent: tableOf(t)})
	if inh {
		b.mayAdd(walk.Scope{Space: spaceConstraint, Deep: true})
	}
}

// dropColumn follows DROP COLUMN, which drops the column's indexes and
// constraints, a sequence it owns, and the column of the inheritance
// children.
func (b *builder) dropColumn(rv *ast.RangeVar, t walk.Key, cmd *ast.AlterTableCmd, inh bool) {
	rng := rangeOf(cmd.Loc)
	name, ok := walkName(cmd.Name)
	if ok && systemColumns[name] {
		b.failAt(rng, nil, review.SystemColumn, systemColumnMessage(cmd.Name, "dropped")) // 0A000
		return
	}
	if ok {
		ck := columnKey(t, name)
		if !cmd.Missing_ok {
			b.need(ck, walk.Present, nil, &walkNote{rng: rng, absent: notFoundIn("column", ident(cmd.Name), relation(rv))}) // 42703
		}
		if !b.w.isTableDirty(t) && b.w.pg != nil && b.w.pg.partKeyCols[[2]string{schemaOfKey(t), t.Name}][name] {
			b.failAt(rng, []walk.Key{t}, review.UnsupportedForm, unsupported("dropping partition key column "+review.Code(ident(cmd.Name)))) // 42P16
			return
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
	b.mayRemove(walk.Scope{Space: spaceTrigger, Parent: tableOf(t)})
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
	if !ok && sub == ast.AT_SetStatistics && cmd.Num != 0 {
		// Column-number ALTER COLUMN <n> SET STATISTICS is only valid on indexes (0A000).
		b.failAt(rng, []walk.Key{t}, review.UnsupportedForm, unsupported("`ALTER TABLE ... ALTER COLUMN` by column number"))
		return
	}
	if ok && systemColumns[name] {
		switch sub {
		case ast.AT_AlterColumnType, ast.AT_ColumnDefault, ast.AT_SetNotNull, ast.AT_DropNotNull:
			b.failAt(rng, nil, review.SystemColumn, systemColumnMessage(cmd.Name, "altered")) // 0A000
		}
		return
	}
	if ok {
		b.need(columnKey(t, name), walk.Present, nil, &walkNote{rng: rng, absent: notFoundIn("column", ident(cmd.Name), relation(rv))}) // 42703
	}
	if sub == ast.AT_AlterColumnType {
		if def, ok := cmd.Def.(*ast.ColumnDef); ok && def.TypeName != nil && def.TypeName.ArrayBounds == nil {
			if tr := b.resolveType(nameParts(def.TypeName.Names)); tr.fact == walk.Present && !tr.builtin && tr.key == typeKey(schemaOfKey(t), t.Name) {
				b.failAt(rng, []walk.Key{t}, review.UnsupportedForm, unsupported("column of table "+review.Code(relation(rv))+"'s own row type")) // 42P16
				return
			}
		}
	}
	if sub == ast.AT_SetStorage {
		if s, ok := cmd.Def.(*ast.String); ok {
			b.checkStorageMode(s.Str, rng)
		}
	}
	if sub == ast.AT_SetCompression {
		if s, ok := cmd.Def.(*ast.String); ok {
			b.checkCompressionMethod(s.Str, rng)
		}
	}
	if sub == ast.AT_SetOptions {
		if opts, ok := cmd.Def.(*ast.List); ok {
			b.checkAttrOptions(opts, rng)
		}
	}
	if sub == ast.AT_SetIdentity {
		if opts, ok := cmd.Def.(*ast.List); ok {
			b.checkIdentityOptions(opts, rng)
		}
	}
	if sub == ast.AT_AddIdentity {
		if def, ok := cmd.Def.(*ast.Constraint); ok && def.Options != nil {
			b.checkIdentityOptions(def.Options, rng)
		}
	}
	if ok && !b.w.isTableDirty(t) {
		sc, hasCol := b.w.syncedCol(t, name)
		switch sub {
		case ast.AT_DropExpression:
			if !cmd.Missing_ok && hasCol && !sc.generated && b.colTouch[name] <= 1 {
				b.failAt(rng, []walk.Key{columnKey(t, name)}, review.UnsupportedForm, unsupported("`DROP EXPRESSION` on non-generated column "+review.Code(ident(cmd.Name)))) // 55000
				return
			}
		case ast.AT_SetExpression:
			if hasCol && !sc.generated && b.colTouch[name] == 0 {
				b.failAt(rng, []walk.Key{columnKey(t, name)}, review.UnsupportedForm, unsupported("`SET EXPRESSION` on non-generated column "+review.Code(ident(cmd.Name)))) // 55000
				return
			}
		case ast.AT_DropIdentity:
			if !cmd.Missing_ok && hasCol && !sc.identity && b.colTouch[name] <= 1 {
				b.failAt(rng, []walk.Key{columnKey(t, name)}, review.UnsupportedForm, unsupported("`DROP IDENTITY` on non-identity column "+review.Code(ident(cmd.Name)))) // 55000
				return
			}
		case ast.AT_SetIdentity:
			if hasCol && !sc.identity && b.colTouch[name] == 0 {
				b.failAt(rng, []walk.Key{columnKey(t, name)}, review.UnsupportedForm, unsupported("identity column alteration on non-identity column "+review.Code(ident(cmd.Name)))) // 55000
				return
			}
		case ast.AT_AddIdentity:
			if hasCol && b.colTouch[name] <= 1 {
				if isKnownNonIntegerBuiltinType(sc.typ) {
					b.failAt(rng, []walk.Key{columnKey(t, name)}, review.UnsupportedForm, unsupported("`ADD GENERATED AS IDENTITY` on non-integer column "+review.Code(ident(cmd.Name)))) // 22023
					return
				}
				if sc.nullable {
					b.failAt(rng, []walk.Key{columnKey(t, name)}, review.UnsupportedForm, unsupported("`ADD GENERATED AS IDENTITY` on nullable column "+review.Code(ident(cmd.Name)))) // 55000
					return
				}
			}
		case ast.AT_DropNotNull:
			if !b.conDropped && b.colTouch[name] <= 1 && b.w.pg != nil && slices.Contains(b.w.pg.pkCols[[2]string{schemaOfKey(t), t.Name}], name) {
				b.failAt(rng, []walk.Key{columnKey(t, name)}, review.UnsupportedForm, unsupported("`DROP NOT NULL` on primary key column "+review.Code(ident(cmd.Name)))) // 42P16
				return
			}
		}
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
	b.mayRemove(walk.Scope{Space: spaceTrigger, Parent: tableOf(t)})
	if b.w.partitioned[t] {
		b.mayRemove(walk.Scope{Space: spaceRelation, Deep: true, Kinds: []string{relIndex}})
		b.mayRemove(walk.Scope{Space: spaceTrigger, Deep: true})
	}
	if inh {
		b.mayRemove(walk.Scope{Space: spaceTrigger, Deep: true})
	}
	if ast.DropBehavior(cmd.Behavior) == ast.DROP_CASCADE {
		b.mayRemove(walk.Scope{Deep: true})
	}
}

// rename follows the RENAME forms of ALTER.
func (b *builder) rename(v *ast.RenameStmt) bool {
	b.session(func() { b.w.allDirty = true })
	switch v.RenameType {
	case ast.OBJECT_TABLE, ast.OBJECT_VIEW, ast.OBJECT_MATVIEW, ast.OBJECT_INDEX, ast.OBJECT_SEQUENCE, ast.OBJECT_FOREIGN_TABLE:
		return b.renameRelation(v)
	case ast.OBJECT_COLUMN:
		return b.renameColumn(v)
	case ast.OBJECT_ATTRIBUTE:
		return b.renameAttribute(v)
	case ast.OBJECT_TABCONSTRAINT:
		return b.renameConstraint(v)
	case ast.OBJECT_SCHEMA:
		return b.renameSchema(v)
	case ast.OBJECT_DOMAIN, ast.OBJECT_DOMCONSTRAINT:
		return b.renameDomain(v)
	case ast.OBJECT_TYPE:
		return b.renameType(v)
	case ast.OBJECT_FUNCTION, ast.OBJECT_PROCEDURE, ast.OBJECT_ROUTINE:
		return b.renameRoutine(v)
	case ast.OBJECT_AGGREGATE:
		name := v.Newname
		b.session(func() { b.w.routines[name] = true })
		return true
	case ast.OBJECT_TRIGGER:
		return b.renameTrigger(v)
	case ast.OBJECT_POLICY:
		return b.renamePolicy(v)
	case ast.OBJECT_RULE:
		// A renamed rule moves between names the walk may have taken
		// from the sync; it no longer knows which rules any table has.
		b.forget(walk.Scope{Space: spaceRule, Deep: true})
		return true
	case ast.OBJECT_EVENT_TRIGGER, ast.OBJECT_COLLATION,
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
	switch {
	case r.fail != "", r.fact == walk.Absent, b.wrongKind(r, v.Relation, v.RenameType, formRename):
		return true
	case r.fact == walk.Unknown:
		return false
	}
	newName, ok := walkName(v.Newname)
	if !ok {
		return false
	}
	s := schemaOfKey(r.key)
	nk := relKey(s, newName)
	b.need(nk, walk.Absent, nil, taken(b.anchor(), r.kind, ident(v.Newname))) // 42P07
	rowType := hasRowType(r.kind)
	if rowType || r.kind == "" {
		b.need(typeKey(s, newName), walk.Absent, nil, &walkNote{rng: b.anchor(), present: alreadyExists("type", ident(v.Newname))}) // 42710
	}
	b.set(r.key, walk.Absent, "")
	b.set(nk, walk.Present, r.kind)
	switch {
	case rowType:
		b.mayRemoveArrayType(s)
		b.mayAddArrayType(s)
		b.set(typeKey(s, r.key.Name), walk.Absent, "")
		b.set(typeKey(s, newName), walk.Present, "")
	case r.kind == "":
		b.mayRemoveArrayType(s)
		b.mayAddArrayType(s)
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

// renamed moves what the walker keeps by relation key. The relation
// holds its schema under the new name, listed or made, and takes nothing
// a relation that had the name before left there.
func (b *builder) renamed(from, to walk.Key) {
	b.session(func() {
		w := b.w
		if from == to {
			return
		}
		w.modifiedViews[from] = true
		w.modifiedViews[to] = true
		w.dirtyTables[from] = true
		w.dirtyTables[to] = true
		w.unlist(to)
		if !slices.Contains(w.made, from) {
			w.made = append(w.made, to)
		}
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
	switch {
	case r.fail != "", r.fact == walk.Absent:
		return true
	case r.fact == walk.Unknown:
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
		if r.kind == relSequence {
			b.need(r.key, walk.Present, []string{relTable, relView, relMatview, relForeign, relComposite}, ofKind(rangeOf(v.Relation.Loc), relation(v.Relation), "table or view")) // 42809
			return true
		}
		if oldName == newName && r.kind != "" {
			b.failAt(rng, nil, review.DuplicateName, alreadyExistsIn("column", ident(v.Newname), relation(v.Relation))) // 42701 / 42703
			return true
		}
		// Another kind's columns are not followed, and some kinds have
		// none to rename (42809).
		b.forget(walk.Scope{Space: spaceColumn, Parent: tableOf(r.key)})
		if inh {
			b.forget(walk.Scope{Space: spaceColumn, Deep: true})
		}
		return true
	}
	if systemColumns[oldName] {
		b.failAt(rng, nil, review.SystemColumn, systemColumnMessage(v.Subname, "renamed")) // 0A000
		return true
	}
	if f, _ := b.l.Lookup(oldKey); systemColumns[newName] && f != walk.Absent {
		b.failAt(rng, nil, review.ReservedName, systemColumnName(v.Newname)) // 42701
		return true
	}
	b.need(oldKey, walk.Present, nil, &walkNote{rng: rng, absent: notFoundIn("column", ident(v.Subname), relation(v.Relation))})      // 42703
	b.need(newKey, walk.Absent, nil, &walkNote{rng: rng, present: alreadyExistsIn("column", ident(v.Newname), relation(v.Relation))}) // 42701
	if inh {
		b.mayRemove(walk.Scope{Space: spaceColumn, Deep: true, Prefix: oldName})
		b.mayAdd(walk.Scope{Space: spaceColumn, Deep: true, Prefix: newName})
	}
	b.set(oldKey, walk.Absent, "")
	b.set(newKey, walk.Present, kindColumn)
	b.session(func() {
		if inh {
			b.w.allDirty = true
		} else {
			b.w.dirtyTables[r.key] = true
		}
	})
	return true
}

// renameConstraint follows RENAME CONSTRAINT: a key's index takes the
// new name too, and the inheritance children's constraint follows.
func (b *builder) renameConstraint(v *ast.RenameStmt) bool {
	r := b.exists(v.Relation, v.MissingOk)
	switch {
	case r.fail != "", r.fact == walk.Absent:
		return true
	case r.fact == walk.Unknown:
		return false
	}
	if r.kind != "" && r.kind != relTable && r.kind != relForeign {
		b.need(r.key, walk.Present, []string{relTable, relForeign}, ofKind(rangeOf(v.Relation.Loc), relation(v.Relation), relTable)) // 42809
		return true
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
	b.need(oldKey, walk.Present, nil, &walkNote{rng: rng, absent: notFoundIn("constraint", ident(v.Subname), relation(v.Relation))})      // 42704
	b.need(newKey, walk.Absent, nil, &walkNote{rng: rng, present: alreadyExistsIn("constraint", ident(v.Newname), relation(v.Relation))}) // 42710
	f, kind := b.l.Lookup(oldKey)
	if oldName == newName && f == walk.Unknown {
		// Renaming a constraint to itself always fails: 42710 if it exists, 42704 if not.
		b.failAt(rng, []walk.Key{r.key}, review.DuplicateName, "Constraint "+review.Code(ident(v.Subname))+" in "+review.Code(relation(v.Relation))+" can't be renamed to itself.")
		return true
	}
	switch {
	case f == walk.Present && indexBacked(kind):
		b.need(relKey(s, newName), walk.Absent, nil, taken(rng, relIndex, ident(v.Newname))) // 42P07
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
		b.fail(nil, walk.Present, review.ReservedName, reservedSchema(v.Newname)) // 42939
		return true
	}
	oldKey := schemaKey(oldName)
	newKey := schemaKey(newName)
	b.need(oldKey, walk.Present, nil, &walkNote{rng: b.anchor(), absent: schemaNotFound(v.Subname)})                 // 3F000
	b.need(newKey, walk.Absent, nil, &walkNote{rng: b.anchor(), present: alreadyExists("schema", ident(v.Newname))}) // 42P06
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
		shellTypes := make(map[walk.Key]bool, len(w.shellTypes))
		for k := range w.shellTypes {
			shellTypes[rekey(k)] = true
		}
		w.shellTypes = shellTypes
		typmodTypes := make(map[walk.Key]bool, len(w.typmodTypes))
		for k := range w.typmodTypes {
			typmodTypes[rekey(k)] = true
		}
		w.typmodTypes = typmodTypes
		rangeMultirange := make(map[walk.Key]walk.Key, len(w.rangeMultirange))
		for rk, mr := range w.rangeMultirange {
			rangeMultirange[rekey(rk)] = rekey(mr)
		}
		w.rangeMultirange = rangeMultirange
		// What the synced schema listed in it is there under the new name.
		for _, k := range w.pg.inSchema[oldName] {
			w.made = append(w.made, rekey(k))
			w.modifiedViews[k] = true
			w.modifiedViews[rekey(k)] = true
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
	b.session(func() { b.w.allDirty = true })
	switch v.ObjectType {
	case ast.OBJECT_TABLE, ast.OBJECT_VIEW, ast.OBJECT_MATVIEW, ast.OBJECT_SEQUENCE, ast.OBJECT_FOREIGN_TABLE:
	case ast.OBJECT_DOMAIN:
		return b.setSchemaDomain(v)
	case ast.OBJECT_TYPE:
		return b.setSchemaType(v)
	case ast.OBJECT_FUNCTION, ast.OBJECT_PROCEDURE, ast.OBJECT_ROUTINE:
		return b.setSchemaRoutine(v)
	case ast.OBJECT_AGGREGATE, ast.OBJECT_COLLATION,
		ast.OBJECT_CONVERSION, ast.OBJECT_OPERATOR, ast.OBJECT_OPCLASS, ast.OBJECT_OPFAMILY, ast.OBJECT_STATISTIC_EXT,
		ast.OBJECT_TSCONFIGURATION, ast.OBJECT_TSDICTIONARY, ast.OBJECT_TSPARSER, ast.OBJECT_TSTEMPLATE:
		return true
	default:
		return false
	}
	r := b.exists(v.Relation, v.MissingOk)
	switch {
	case r.fail != "", r.fact == walk.Absent, b.wrongKind(r, v.Relation, v.ObjectType, formSetSchema):
		return true
	case r.fact == walk.Unknown:
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
	b.need(schemaKey(ns), walk.Present, nil, &walkNote{rng: rng, absent: schemaNotFound(v.Newschema)}) // 3F000
	nk := relKey(ns, r.key.Name)
	moved := qualified([]string{v.Newschema, r.key.Name})
	b.need(nk, walk.Absent, nil, taken(rng, r.kind, moved)) // 42P07
	rowType := hasRowType(r.kind)
	if rowType || r.kind == "" {
		b.need(typeKey(ns, r.key.Name), walk.Absent, nil, &walkNote{rng: rng, present: alreadyExists("type", moved)}) // 42710
	}
	b.set(r.key, walk.Absent, "")
	b.set(nk, walk.Present, r.kind)
	switch {
	case rowType:
		b.mayRemoveArrayType(old)
		b.mayAddArrayType(ns)
		b.set(typeKey(old, r.key.Name), walk.Absent, "")
		b.set(typeKey(ns, r.key.Name), walk.Present, "")
	case r.kind == "":
		b.mayRemoveArrayType(old)
		b.mayAddArrayType(ns)
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
// the word a message uses for it.
var dropKinds = map[ast.ObjectType]struct{ kind, word string }{
	ast.OBJECT_TABLE:         {relTable, "table"},
	ast.OBJECT_VIEW:          {relView, "view"},
	ast.OBJECT_MATVIEW:       {relMatview, "materialized view"},
	ast.OBJECT_INDEX:         {relIndex, "index"},
	ast.OBJECT_SEQUENCE:      {relSequence, "sequence"},
	ast.OBJECT_FOREIGN_TABLE: {relForeign, "foreign table"},
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
	// Every object is looked up before any is dropped, so each check
	// holds against the state before the statement, whatever another
	// object's drop does (verified on PostgreSQL 17 and 18).
	//
	// pg: src/backend/commands/tablecmds.c — RemoveRelations;
	// src/backend/commands/dropcmds.c — RemoveObjects
	b.Order = walk.Lookups
	if v.Concurrent && kind == ast.OBJECT_INDEX {
		// Refused before any name is looked up, IF EXISTS or not (0A000,
		// verified on PostgreSQL 17 and 18 outside a transaction, as
		// Bytebase runs CONCURRENTLY).
		//
		// pg: src/backend/commands/tablecmds.c — RemoveRelations
		switch {
		case len(v.Objects.Items) != 1:
			b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("`DROP INDEX CONCURRENTLY` of several indexes"))
			return true
		case cascade:
			b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("`DROP INDEX CONCURRENTLY` with `CASCADE`"))
			return true
		}
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
			var list *ast.List
			var tn *ast.TypeName
			switch o := obj.(type) {
			case *ast.List:
				list = o
			case *ast.TypeName:
				tn = o
				list = o.Names
			}
			parts := nameParts(list)
			if list == nil || len(parts) != len(list.Items) {
				return false
			}
			b.dropType(parts, tn, v.Missing_ok, kind == ast.OBJECT_DOMAIN)
		}
	case dropKinds[kind].kind != "":
		var dropping map[walk.Key]bool
		if kind == ast.OBJECT_TABLE && !cascade {
			dropping = make(map[walk.Key]bool, len(v.Objects.Items))
			for _, obj := range v.Objects.Items {
				if r := b.nameList(nameParts(listOf(obj))); r.fact == walk.Present {
					dropping[r.key] = true
				}
			}
		}
		for _, obj := range v.Objects.Items {
			b.next()
			b.dropRelation(kind, nameParts(listOf(obj)), v.Missing_ok, cascade, dropping)
		}
	case kind == ast.OBJECT_FUNCTION || kind == ast.OBJECT_PROCEDURE || kind == ast.OBJECT_ROUTINE:
		for _, obj := range v.Objects.Items {
			b.next()
			b.dropRoutine(obj, kind, v.Missing_ok, cascade)
		}
	case kind == ast.OBJECT_TRIGGER:
		for _, obj := range v.Objects.Items {
			b.next()
			b.dropTrigger(nameParts(listOf(obj)), v.Missing_ok)
		}
	case kind == ast.OBJECT_POLICY:
		for _, obj := range v.Objects.Items {
			b.next()
			b.dropPolicy(nameParts(listOf(obj)), v.Missing_ok)
		}
	case !dropNoFacts[kind]:
		return false
	}
	if cascade {
		b.mayRemove(walk.Scope{Deep: true})
		b.session(func() { b.w.allDirty = true })
	}
	return true
}

func (b *builder) dropRelation(kind ast.ObjectType, parts []string, missingOK, cascade bool, dropping map[walk.Key]bool) {
	want := dropKinds[kind]
	written := qualified(parts)
	r := b.nameList(parts)
	if r.fail != "" {
		b.fail(nil, walk.Absent, review.CrossDatabaseReference, r.fail)
		return
	}
	switch r.fact {
	case walk.Absent:
		if missingOK {
			return
		}
		b.fail(r.keys, walk.Absent, review.ObjectNotFound, missing(r, want.word, written)) // 42P01; 42704 for an index
		return
	case walk.Unknown:
		b.dropUnknown(parts)
		b.session(func() { b.w.allDirty = true })
		return
	}
	b.need(r.key, walk.Present, []string{want.kind}, ofKind(b.anchor(), written, want.word)) // 42809
	if kind == ast.OBJECT_TABLE && !cascade && r.kind == relTable {
		if ref, ok := b.activeFKDependent(r.key, dropping, nil); ok {
			b.fail([]walk.Key{r.key, relKey(ref.fromSchema, ref.fromTable)}, walk.Present, review.UnsupportedForm, "Table "+review.Code(written)+" is referenced by foreign key on "+review.Code(ident(ref.fromTable))+" and can't be dropped without `CASCADE`.") // 2BP01
			return
		}
	}
	b.session(func() { b.w.dirtyTables[r.key] = true })
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
	b.mayRemove(walk.Scope{Space: spaceTrigger, Deep: true})
	b.mayRemove(walk.Scope{Space: spacePolicy, Deep: true})
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
		b.mayRemoveArrayType(s)
		b.set(typeKey(s, k.Name), walk.Absent, "")
	case kind == "":
		b.mayRemoveArrayType(s)
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
	b.session(func() {
		b.w.modifiedViews[k] = true
		b.w.dirtyTables[k] = true
	})
	b.unlisted(k)
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
		b.need(k, walk.Present, nil, &walkNote{rng: b.anchor(), absent: schemaNotFound(parts[0])}) // 3F000
	}
	if f, _ := b.l.Lookup(k); f == walk.Present && !cascade {
		if held, ok := b.holds(name); ok {
			b.fail([]walk.Key{k, held}, walk.Present, review.SchemaNotEmpty, "Schema "+review.Code(ident(parts[0]))+" still holds objects.")
		}
	}
	b.set(k, walk.Absent, "")
	b.clear(walk.Scope{Parent: walk.Path(name), Deep: true})
	for rk, mrKey := range b.w.rangeMultirange {
		if rk.Parent == walk.Path(name) && mrKey.Parent != walk.Path(name) {
			if mf, _ := b.l.Lookup(mrKey); mf == walk.Present {
				b.set(mrKey, walk.Unknown, "")
			}
			b.mayRemoveArrayType(schemaOfKey(mrKey))
			b.mayRemove(walk.Scope{Space: spaceFunction, Parent: mrKey.Parent})
		}
	}
	b.session(func() {
		b.w.allDirty = true
		// What it held is gone, and so is the partition state of it.
		w := b.w
		in := func(k walk.Key) bool { return k.Parent == walk.Path(name) }
		for k := range w.partitioned {
			if in(k) {
				delete(w.partitioned, k)
			}
		}
		for parent, children := range w.partitions {
			if in(parent) {
				delete(w.partitions, parent)
				continue
			}
			w.partitions[parent] = slices.DeleteFunc(children, in)
		}
		for k := range w.shellTypes {
			if in(k) {
				delete(w.shellTypes, k)
			}
		}
		for k := range w.typmodTypes {
			if in(k) {
				delete(w.typmodTypes, k)
			}
		}
		for rk, mrKey := range w.rangeMultirange {
			if in(rk) || in(mrKey) {
				delete(w.rangeMultirange, rk)
			}
		}
		for _, k := range w.pg.inSchema[name] {
			w.modifiedViews[k] = true
		}
	})
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

// dropType follows DROP TYPE and DROP DOMAIN.
func (b *builder) dropType(parts []string, tn *ast.TypeName, missingOK, domain bool) {
	written := qualified(parts)
	if tn != nil {
		written = typeName(tn)
	}
	rng := b.anchor()
	// Check built-in pg_catalog types explicitly qualified by the parser (e.g.
	// INT, BIT, CHARACTER(1), INTERVAL, TIMESTAMP WITH TIME ZONE).
	if len(parts) == 2 && parts[0] == "pg_catalog" && builtinCatalogTypes[parts[1]] {
		if tn != nil && tn.Typmods != nil && hasNonConstTypmod(tn.Typmods) {
			b.failAt(rng, nil, review.UnsupportedForm, unsupported("invalid type modifier for "+review.Code(written))) // 22P02
			return
		}
		if domain {
			b.failAt(rng, nil, review.WrongObjectKind, review.Code(written)+" isn't a domain.") // 42809
			return
		}
		b.failAt(rng, nil, review.UnsupportedForm, "Built-in type "+review.Code(written)+" can't be dropped.") // 2BP01
		return
	}
	r := b.resolveType(parts)
	if r.fail != "" {
		b.failAt(rng, nil, review.CrossDatabaseReference, r.fail)
		return
	}
	if r.fact == walk.Absent {
		if !missingOK {
			if r.missingSchema != "" {
				b.failAt(rng, r.keys, review.ObjectNotFound, schemaNotFound(r.missingSchema)) // 3F000
			} else {
				kind := typOther
				if domain {
					kind = typDomain
				}
				b.failAt(rng, r.keys, review.ObjectNotFound, notFound(kind, written)) // 42704
			}
		}
		return
	}
	if r.fact == walk.Unknown {
		if !domain && !missingOK && tn != nil && tn.ArrayBounds != nil {
			// DROP TYPE t[] without IF EXISTS always fails: 2BP01 if t exists, 42704 if not.
			b.failAt(rng, nil, review.UnsupportedForm, "Array type "+review.Code(written)+" can't be dropped directly.")
			return
		}
		if name, ok := walkName(parts[len(parts)-1]); ok {
			b.mayRemove(walk.Scope{Space: spaceType, Deep: true, Prefix: name})
			b.mayRemove(walk.Scope{Space: spaceRelation, Deep: true, Prefix: name, Kinds: []string{relComposite}})
		} else {
			b.mayRemove(walk.Scope{Space: spaceType, Deep: true})
			b.mayRemove(walk.Scope{Space: spaceRelation, Deep: true, Kinds: []string{relComposite}})
		}
		if !domain {
			b.mayRemove(walk.Scope{Space: spaceFunction, Deep: true})
			for _, mrKey := range b.w.rangeMultirange {
				if mf, _ := b.l.Lookup(mrKey); mf == walk.Present {
					b.set(mrKey, walk.Unknown, "")
				}
			}
		}
		b.session(func() { b.w.allDirty = true })
		return
	}
	// r.fact == walk.Present.
	if tn != nil && tn.Typmods != nil {
		if r.builtin && hasNonConstTypmod(tn.Typmods) {
			b.failAt(rng, nil, review.UnsupportedForm, unsupported("invalid type modifier for "+review.Code(written))) // 22P02
			return
		}
		if !r.builtin && !b.w.typmodTypes[r.key] && !b.w.typmodTypes[walk.Key{}] {
			b.failAt(rng, r.keys, review.UnsupportedForm, unsupported("type modifier for "+review.Code(written))) // 42601
			return
		}
	}
	if domain {
		if r.builtin || (tn != nil && tn.ArrayBounds != nil) {
			b.failAt(rng, r.keys, review.WrongObjectKind, review.Code(written)+" isn't a domain.") // 42809
			return
		}
		// DROP DOMAIN drops a domain and no other type, IF EXISTS or not
		// (42809, verified on PostgreSQL 17 and 18).
		b.need(r.key, walk.Present, []string{typDomain}, &walkNote{rng: rng, absent: notFound(typDomain, written), present: review.Code(written) + " isn't a domain."})
	} else {
		if tn != nil && tn.ArrayBounds != nil {
			b.failAt(rng, r.keys, review.UnsupportedForm, "Array type "+review.Code(written)+" can't be dropped directly.") // 2BP01
			return
		}
		if r.builtin {
			b.failAt(rng, r.keys, review.UnsupportedForm, "Built-in type "+review.Code(written)+" can't be dropped.") // 2BP01
			return
		}
		s := schemaOfKey(r.key)
		rk := relKey(s, r.key.Name)
		if rf, rkind := b.l.Lookup(rk); rf == walk.Present && hasRowType(rkind) && rkind != relComposite {
			b.failAt(rng, []walk.Key{rk}, review.UnsupportedForm, "Row type "+review.Code(written)+" belongs to "+article(rkind)+" and can't be dropped directly.") // 2BP01
			return
		}
	}
	s := schemaOfKey(r.key)
	b.mayRemoveArrayType(s)
	if !domain {
		// Dropping a range type also drops its companion multirange type (and
		// array type) and its range/multirange constructor functions in pg_proc.
		b.mayRemove(walk.Scope{Space: spaceFunction, Parent: walk.Path(s)})
		if mrKey, ok := b.w.rangeMultirange[r.key]; ok {
			if mf, _ := b.l.Lookup(mrKey); mf == walk.Present {
				b.set(mrKey, walk.Unknown, "")
			}
			if ms := schemaOfKey(mrKey); ms != s {
				b.mayRemoveArrayType(ms)
				b.mayRemove(walk.Scope{Space: spaceFunction, Parent: mrKey.Parent})
			}
		}
	}
	b.set(r.key, walk.Absent, "")
	rk := relKey(s, r.key.Name)
	b.session(func() {
		delete(b.w.shellTypes, r.key)
		delete(b.w.typmodTypes, r.key)
		delete(b.w.rangeMultirange, r.key)
		for k, mr := range b.w.rangeMultirange {
			if mr == r.key {
				delete(b.w.rangeMultirange, k)
			}
		}
		b.w.dirtyTypes[rk] = true
	})
	switch f, kind := b.l.Lookup(rk); {
	case f == walk.Present && kind == relComposite:
		b.set(rk, walk.Absent, "")
		b.clear(walk.Scope{Parent: tableOf(rk), Deep: true})
	case f == walk.Unknown, f == walk.Present && kind == "":
		b.set(rk, walk.Unknown, "")
	}
}

func hasNonConstTypmod(mods *ast.List) bool {
	if mods == nil {
		return false
	}
	for _, item := range mods.Items {
		switch item.(type) {
		case *ast.A_Const, *ast.Integer:
		default:
			return true
		}
	}
	return false
}

// commentKinds is the relation kind each COMMENT ON object kind names.
var commentKinds = map[ast.ObjectType]struct{ kind, word string }{
	ast.OBJECT_TABLE:         {relTable, "table"},
	ast.OBJECT_VIEW:          {relView, "view"},
	ast.OBJECT_MATVIEW:       {relMatview, "materialized view"},
	ast.OBJECT_INDEX:         {relIndex, "index"},
	ast.OBJECT_SEQUENCE:      {relSequence, "sequence"},
	ast.OBJECT_FOREIGN_TABLE: {relForeign, "foreign table"},
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
		switch {
		case r.fail != "":
			b.fail(nil, walk.Absent, review.CrossDatabaseReference, r.fail)
		case r.fact == walk.Absent:
			b.fail(r.keys, walk.Absent, review.ObjectNotFound, missing(r, want.word, qualified(parts))) // 42P01
		case r.fact == walk.Present:
			b.need(r.key, walk.Present, []string{want.kind}, ofKind(b.anchor(), qualified(parts), want.word)) // 42809
		}
	case kind == ast.OBJECT_COLUMN:
		if len(parts) < 2 {
			b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("unqualified column name in `COMMENT ON COLUMN`")) // 42601
			return true
		}
		table := parts[:len(parts)-1]
		r := b.nameList(table)
		switch {
		case r.fail != "":
			b.fail(nil, walk.Absent, review.CrossDatabaseReference, r.fail)
		case r.fact == walk.Absent:
			b.fail(r.keys, walk.Absent, review.ObjectNotFound, missing(r, "", qualified(table)))
		case r.fact == walk.Present:
			if col, ok := walkName(parts[len(parts)-1]); ok && !systemColumns[col] && hasColumns(r.kind) {
				b.need(columnKey(r.key, col), walk.Present, nil, &walkNote{rng: b.anchor(), absent: notFoundIn("column", ident(parts[len(parts)-1]), qualified(table))}) // 42703
			}
		}
	case kind == ast.OBJECT_TABCONSTRAINT:
		if len(parts) < 2 {
			return true
		}
		table := parts[:len(parts)-1]
		r := b.nameList(table)
		switch {
		case r.fail != "":
			b.fail(nil, walk.Absent, review.CrossDatabaseReference, r.fail)
		case r.fact == walk.Absent:
			b.fail(r.keys, walk.Absent, review.ObjectNotFound, missing(r, "", qualified(table))) // 42P01
		case r.fact == walk.Present:
			if r.kind == relIndex || r.kind == relComposite {
				b.need(r.key, walk.Present, []string{relTable, relForeign}, ofKind(b.anchor(), qualified(table), relTable)) // 42809
				return true
			}
			if r.kind == relView || r.kind == relMatview || r.kind == relSequence {
				b.fail([]walk.Key{r.key}, walk.Present, review.ObjectNotFound, notFoundIn("constraint", ident(parts[len(parts)-1]), qualified(table))) // 42704
				return true
			}
			if con, ok := walkName(parts[len(parts)-1]); ok {
				b.need(constraintKey(r.key, con), walk.Present, nil, &walkNote{rng: b.anchor(), absent: notFoundIn("constraint", ident(parts[len(parts)-1]), qualified(table))}) // 42704
			}
		}
	case kind == ast.OBJECT_TRIGGER:
		if len(parts) < 2 {
			return true
		}
		table := parts[:len(parts)-1]
		r := b.nameList(table)
		switch {
		case r.fail != "":
			b.fail(nil, walk.Absent, review.CrossDatabaseReference, r.fail)
		case r.fact == walk.Absent:
			b.fail(r.keys, walk.Absent, review.ObjectNotFound, missing(r, "", qualified(table))) // 42P01
		case r.fact == walk.Present:
			if r.kind == relIndex || r.kind == relComposite {
				b.need(r.key, walk.Present, []string{relTable, relView, relForeign}, ofKind(b.anchor(), qualified(table), "table or view")) // 42809
				return true
			}
			if r.kind == relMatview || r.kind == relSequence {
				b.fail([]walk.Key{r.key}, walk.Present, review.ObjectNotFound, notFoundIn("trigger", ident(parts[len(parts)-1]), qualified(table))) // 42704
				return true
			}
			if trig, ok := walkName(parts[len(parts)-1]); ok {
				b.need(trigKey(r.key, trig), walk.Present, nil, &walkNote{rng: b.anchor(), absent: notFoundIn("trigger", ident(parts[len(parts)-1]), qualified(table))}) // 42704
			}
		}
	case kind == ast.OBJECT_POLICY:
		if len(parts) < 2 {
			return true
		}
		table := parts[:len(parts)-1]
		r := b.nameList(table)
		switch {
		case r.fail != "":
			b.fail(nil, walk.Absent, review.CrossDatabaseReference, r.fail)
		case r.fact == walk.Absent:
			b.fail(r.keys, walk.Absent, review.ObjectNotFound, missing(r, "", qualified(table))) // 42P01
		case r.fact == walk.Present:
			if r.kind == relIndex || r.kind == relComposite {
				b.need(r.key, walk.Present, []string{relTable}, ofKind(b.anchor(), qualified(table), relTable)) // 42809
				return true
			}
			if r.kind != "" && r.kind != relTable {
				b.fail([]walk.Key{r.key}, walk.Present, review.ObjectNotFound, notFoundIn("policy", ident(parts[len(parts)-1]), qualified(table))) // 42704
				return true
			}
			if pol, ok := walkName(parts[len(parts)-1]); ok {
				b.need(polKey(r.key, pol), walk.Present, nil, &walkNote{rng: b.anchor(), absent: notFoundIn("policy", ident(parts[len(parts)-1]), qualified(table))}) // 42704
			}
		}
	case kind == ast.OBJECT_TYPE:
		var tParts []string
		if tn, ok := v.Object.(*ast.TypeName); ok {
			tParts = nameParts(tn.Names)
		} else {
			tParts = parts
		}
		if len(tParts) > 0 {
			r := b.resolveType(tParts)
			switch {
			case r.fail != "":
				b.fail(nil, walk.Absent, review.CrossDatabaseReference, r.fail)
			case r.fact == walk.Absent:
				if r.missingSchema != "" {
					b.fail(r.keys, walk.Absent, review.ObjectNotFound, schemaNotFound(r.missingSchema))
				} else {
					b.fail(r.keys, walk.Absent, review.ObjectNotFound, notFound("type", qualified(tParts)))
				}
			}
		}
	case kind == ast.OBJECT_SCHEMA:
		if s, ok := v.Object.(*ast.String); ok {
			if name, ok := walkName(s.Str); ok {
				b.need(schemaKey(name), walk.Present, nil, &walkNote{rng: b.anchor(), absent: schemaNotFound(s.Str)}) // 3F000
			}
		}
	case kind == ast.OBJECT_DOMAIN:
		var dParts []string
		if tn, ok := v.Object.(*ast.TypeName); ok {
			dParts = nameParts(tn.Names)
		} else {
			dParts = parts
		}
		if len(dParts) > 0 {
			b.checkDomain(dParts, b.anchor())
		}
	case kind == ast.OBJECT_FUNCTION, kind == ast.OBJECT_PROCEDURE, kind == ast.OBJECT_ROUTINE:
		if owa, ok := v.Object.(*ast.ObjectWithArgs); ok {
			// COMMENT ON FUNCTION / ROUTINE also accepts aggregates (prokind = 'a'),
			// which pg/driver/sync.go omits from SchemaMetadata.Functions.
			b.checkRoutine(owa, kind, false, kind != ast.OBJECT_PROCEDURE)
		}
	case kind == ast.OBJECT_DOMCONSTRAINT:
		if list := listOf(v.Object); list != nil && len(list.Items) >= 2 {
			var dParts []string
			if tn, ok := list.Items[0].(*ast.TypeName); ok {
				dParts = nameParts(tn.Names)
			} else if len(list.Items) == 2 {
				if s, ok := list.Items[0].(*ast.String); ok {
					dParts = []string{s.Str}
				}
			} else if len(list.Items) == 3 {
				if s0, ok0 := list.Items[0].(*ast.String); ok0 {
					if s1, ok1 := list.Items[1].(*ast.String); ok1 {
						dParts = []string{s0.Str, s1.Str}
					}
				}
			}
			if len(dParts) > 0 {
				b.checkDomain(dParts, b.anchor())
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
				b.need(schemaKey(name), walk.Present, nil, &walkNote{rng: b.anchor(), absent: schemaNotFound(parts[0])}) // 3F000
			}
		}
	}
	return true
}

// truncate follows TRUNCATE, which changes no name: each relation must be
// a table (42P01, 42809). A TRUNCATE trigger may run anything.
func (b *builder) truncate(v *ast.TruncateStmt) bool {
	if v.Relations != nil {
		cascade := ast.DropBehavior(v.Behavior) == ast.DROP_CASCADE
		var truncating map[walk.Key]bool
		var inhTargets []walk.Key
		hasUnknownTarget := false
		if !cascade {
			truncating = make(map[walk.Key]bool, len(v.Relations.Items))
			for _, item := range v.Relations.Items {
				if rv, ok := item.(*ast.RangeVar); ok {
					r := b.rangeVar(rv)
					if r.fact == walk.Present {
						truncating[r.key] = true
						if rv.Inh {
							inhTargets = append(inhTargets, r.key)
						}
					} else if r.fact == walk.Unknown {
						hasUnknownTarget = true
					}
				}
			}
		}
		for _, item := range v.Relations.Items {
			rv, ok := item.(*ast.RangeVar)
			if !ok {
				continue
			}
			r := b.exists(rv, false)
			switch {
			case r.fact != walk.Present:
			case !rv.Inh && r.kind == relTable && b.w.partitioned[r.key]:
				// ONLY leaves a partitioned table nothing to truncate
				// (42809, verified on PostgreSQL 17 and 18).
				//
				// pg: src/backend/commands/tablecmds.c — ExecuteTruncate
				b.failAt(rangeOf(rv.Loc), []walk.Key{r.key}, review.UnsupportedForm, unsupported("`TRUNCATE ONLY` of partitioned table "+review.Code(relation(rv))))
			default:
				b.need(r.key, walk.Present, []string{relTable, relForeign}, ofKind(rangeOf(rv.Loc), relation(rv), relTable))
				if !cascade && !hasUnknownTarget && r.kind == relTable {
					if ref, ok := b.activeFKDependent(r.key, truncating, inhTargets); ok {
						b.failAt(rangeOf(rv.Loc), []walk.Key{r.key, relKey(ref.fromSchema, ref.fromTable)}, review.UnsupportedForm, unsupported("`TRUNCATE` on table "+review.Code(relation(rv))+" referenced by a foreign key without `CASCADE`")) // 0A000
					}
				}
			}
		}
	}
	b.Last = b.w.triggered
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
	b.Last = b.w.triggered || b.w.callsUserCode(n)
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

// createTrigger follows CREATE [OR REPLACE] [CONSTRAINT] TRIGGER: checks
// target relation existence and relkind, PostgreSQL's unconditional timing /
// event / transition-table rules, and duplicate trigger name on the target.
//
// pg: src/backend/commands/trigger.c — CreateTrigger
func (b *builder) createTrigger(v *ast.CreateTrigStmt) bool {
	// CREATE OR REPLACE CONSTRAINT TRIGGER never reaches here: gram.y's
	// action refuses it, and so does the parser.
	if v.IsConstraint && v.Constrrel != nil {
		b.exists(v.Constrrel, false)
	}
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
	if r.fail != "" || r.fact == walk.Absent {
		return true
	}
	rng := rangeOf(v.Relation.Loc)
	isInstead := v.Timing&ast.TRIGGER_TYPE_INSTEAD != 0
	isBefore := v.Timing&ast.TRIGGER_TYPE_BEFORE != 0
	isAfter := !isInstead && !isBefore
	hasTruncate := v.Events&ast.TRIGGER_TYPE_TRUNCATE != 0
	hasInsert := v.Events&ast.TRIGGER_TYPE_INSERT != 0
	hasUpdate := v.Events&ast.TRIGGER_TYPE_UPDATE != 0
	hasDelete := v.Events&ast.TRIGGER_TYPE_DELETE != 0
	hasTransitions := v.TransitionRels != nil && len(v.TransitionRels.Items) > 0

	if r.fact == walk.Present && r.kind != "" {
		switch r.kind {
		case relTable:
			if isInstead {
				b.need(r.key, walk.Present, []string{relView}, ofKind(rng, relation(v.Relation), relView)) // 42809
				return true
			}
			if v.Row && hasTransitions && b.w.partitioned[r.key] {
				b.failAt(rng, []walk.Key{r.key}, review.UnsupportedForm, unsupported("`FOR EACH ROW` trigger with transition tables on partitioned table "+review.Code(relation(v.Relation)))) // 0A000
				return true
			}
		case relView:
			if !isInstead && (v.Row || v.IsConstraint) {
				b.need(r.key, walk.Present, []string{relTable}, ofKind(rng, relation(v.Relation), relTable)) // 42809
				return true
			}
			if hasTruncate {
				b.need(r.key, walk.Present, []string{relTable}, ofKind(rng, relation(v.Relation), relTable)) // 42809
				return true
			}
		case relForeign:
			if isInstead {
				b.need(r.key, walk.Present, []string{relView}, ofKind(rng, relation(v.Relation), relView)) // 42809
				return true
			}
			if hasTruncate {
				b.need(r.key, walk.Present, []string{relTable}, ofKind(rng, relation(v.Relation), relTable)) // 42809
				return true
			}
		default:
			b.need(r.key, walk.Present, []string{relTable, relView, relForeign}, ofKind(rng, relation(v.Relation), "table or view")) // 42809
			return true
		}
	}
	if v.Row && hasTruncate {
		b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("`TRUNCATE FOR EACH ROW` trigger")) // 42P17
		return true
	}
	if isInstead {
		if !v.Row {
			b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("`INSTEAD OF` trigger without `FOR EACH ROW`")) // 0A000
			return true
		}
		if v.WhenClause != nil {
			b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("`INSTEAD OF` trigger with `WHEN` condition")) // 42P17
			return true
		}
		if v.Columns != nil && len(v.Columns.Items) > 0 {
			b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("`INSTEAD OF` trigger with column list")) // 42P17
			return true
		}
	}
	if hasTransitions {
		for _, item := range v.TransitionRels.Items {
			if tr, ok := item.(*ast.TriggerTransition); ok && !tr.IsTable {
				b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("`REFERENCING ... ROW` in `CREATE TRIGGER`")) // 0A000
				return true
			}
		}
		if r.fact == walk.Present && (r.kind == relView || r.kind == relForeign) {
			b.need(r.key, walk.Present, []string{relTable}, ofKind(rng, relation(v.Relation), relTable)) // 42809
			return true
		}
		if !isAfter {
			b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("transition tables on a non-`AFTER` trigger")) // 42P17
			return true
		}
		if hasTruncate {
			b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("transition tables on a `TRUNCATE` trigger")) // 42P17
			return true
		}
		dmlEvents := 0
		if hasInsert {
			dmlEvents++
		}
		if hasUpdate {
			dmlEvents++
		}
		if hasDelete {
			dmlEvents++
		}
		if dmlEvents > 1 {
			b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("transition tables on a trigger with multiple events")) // 42P17
			return true
		}
		if v.Columns != nil && len(v.Columns.Items) > 0 {
			b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("transition tables on a trigger with a column list")) // 42P17
			return true
		}
		var seenNew, seenOld bool
		seenNames := make(map[string]bool)
		for _, item := range v.TransitionRels.Items {
			tr, ok := item.(*ast.TriggerTransition)
			if !ok {
				continue
			}
			if tr.IsNew {
				if seenNew {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("multiple `NEW TABLE` transition relations in `CREATE TRIGGER`")) // 42P17
					return true
				}
				seenNew = true
				if !hasInsert && !hasUpdate {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("`NEW TABLE` transition relation on a `DELETE` trigger")) // 42P17
					return true
				}
			} else {
				if seenOld {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("multiple `OLD TABLE` transition relations in `CREATE TRIGGER`")) // 42P17
					return true
				}
				seenOld = true
				if !hasDelete && !hasUpdate {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("`OLD TABLE` transition relation on an `INSERT` trigger")) // 42P17
					return true
				}
			}
			if seenNames[tr.Name] {
				b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("duplicate transition table name "+review.Code(ident(tr.Name))+" in `CREATE TRIGGER`")) // 42P17
				return true
			}
			seenNames[tr.Name] = true
		}
	}
	if v.WhenClause != nil && hasBareTriggerColumn(v.WhenClause) {
		b.fail(nil, walk.Absent, review.UnsupportedForm, unsupported("unqualified column reference in trigger `WHEN` condition")) // 42702 / 42703
		return true
	}
	dupTrig := false
	if r.fact == walk.Present {
		if tn, ok := walkName(v.Trigname); ok {
			tk := trigKey(r.key, tn)
			if !v.Replace {
				if f, _ := b.l.Lookup(tk); f == walk.Present {
					dupTrig = true
				}
				b.need(tk, walk.Absent, nil, &walkNote{rng: b.anchor(), present: alreadyExistsIn("trigger", ident(v.Trigname), relation(v.Relation))}) // 42710
			}
			b.set(tk, walk.Present, "")
			if b.w.partitioned[r.key] || r.kind == "" {
				// PostgreSQL clones row/constraint triggers on a partitioned
				// table onto its partitions (pg: cloneSingleTriggerToPartitions).
				b.mayAdd(walk.Scope{Space: spaceTrigger, Deep: true, Prefix: tn})
			}
		} else if b.w.partitioned[r.key] || r.kind == "" {
			b.mayAdd(walk.Scope{Space: spaceTrigger, Deep: true})
		}
	} else {
		if tn, ok := walkName(v.Trigname); ok {
			b.mayAdd(walk.Scope{Space: spaceTrigger, Deep: true, Prefix: tn})
		} else {
			b.mayAdd(walk.Scope{Space: spaceTrigger, Deep: true})
		}
	}
	if !dupTrig {
		if fnParts := nameParts(v.Funcname); len(fnParts) > 0 {
			b.checkRoutineName(fnParts, ast.OBJECT_FUNCTION, false, false)
		}
	}
	if r.fact == walk.Present {
		b.session(func() { b.w.dirtyTables[r.key] = true })
	} else {
		b.session(func() { b.w.allDirty = true })
	}
	return true
}

// dropTrigger follows DROP TRIGGER [IF EXISTS] name ON table.
func (b *builder) dropTrigger(parts []string, missingOK bool) {
	if len(parts) < 2 {
		return
	}
	table := parts[:len(parts)-1]
	trigWritten := parts[len(parts)-1]
	r := b.nameList(table)
	switch {
	case r.fail != "":
		b.fail(nil, walk.Absent, review.CrossDatabaseReference, r.fail)
		return
	case r.fact == walk.Absent:
		if !missingOK {
			b.fail(r.keys, walk.Absent, review.ObjectNotFound, missing(r, "", qualified(table))) // 42P01
		}
		return
	case r.fact == walk.Unknown:
		if tn, ok := walkName(trigWritten); ok {
			b.mayRemove(walk.Scope{Space: spaceTrigger, Deep: true, Prefix: tn})
		} else {
			b.mayRemove(walk.Scope{Space: spaceTrigger, Deep: true})
		}
		return
	}
	if r.kind == relIndex || r.kind == relComposite {
		b.need(r.key, walk.Present, []string{relTable, relView, relForeign}, ofKind(b.anchor(), qualified(table), "table or view")) // 42809
		return
	}
	if r.kind == relMatview || r.kind == relSequence {
		if !missingOK {
			b.fail([]walk.Key{r.key}, walk.Present, review.ObjectNotFound, notFoundIn("trigger", ident(trigWritten), qualified(table))) // 42704
		}
		return
	}
	if tn, ok := walkName(trigWritten); ok {
		tk := trigKey(r.key, tn)
		if !missingOK {
			b.need(tk, walk.Present, nil, &walkNote{rng: b.anchor(), absent: notFoundIn("trigger", ident(trigWritten), qualified(table))}) // 42704
		}
		b.session(func() { b.w.dirtyTables[r.key] = true })
		b.set(tk, walk.Absent, "")
		if b.w.partitioned[r.key] || r.kind == "" {
			b.mayRemove(walk.Scope{Space: spaceTrigger, Deep: true, Prefix: tn})
		}
	} else if b.w.partitioned[r.key] || r.kind == "" {
		b.mayRemove(walk.Scope{Space: spaceTrigger, Deep: true})
	}
}

// renameTrigger follows ALTER TRIGGER old ON table RENAME TO new.
func (b *builder) renameTrigger(v *ast.RenameStmt) bool {
	r := b.exists(v.Relation, false)
	if r.fail != "" || r.fact == walk.Absent {
		return true
	}
	if r.fact == walk.Present && r.kind != "" && r.kind != relTable && r.kind != relView && r.kind != relForeign {
		b.need(r.key, walk.Present, []string{relTable, relView, relForeign}, ofKind(rangeOf(v.Relation.Loc), relation(v.Relation), "table or view")) // 42809
		return true
	}
	if r.fact == walk.Present {
		oldName, ok1 := walkName(v.Subname)
		newName, ok2 := walkName(v.Newname)
		if ok1 && ok2 {
			okKey := trigKey(r.key, oldName)
			b.need(okKey, walk.Present, nil, &walkNote{rng: b.anchor(), absent: notFoundIn("trigger", ident(v.Subname), relation(v.Relation))})
			if oldName != newName {
				nkKey := trigKey(r.key, newName)
				b.need(nkKey, walk.Absent, nil, &walkNote{rng: b.anchor(), present: alreadyExistsIn("trigger", ident(v.Newname), relation(v.Relation))})
				b.set(okKey, walk.Absent, "")
				b.set(nkKey, walk.Present, "")
				if b.w.partitioned[r.key] || r.kind == "" {
					b.mayRemove(walk.Scope{Space: spaceTrigger, Deep: true, Prefix: oldName})
					b.mayAdd(walk.Scope{Space: spaceTrigger, Deep: true, Prefix: newName})
				}
			}
		} else {
			b.forget(walk.Scope{Space: spaceTrigger, Parent: tableOf(r.key)})
			if b.w.partitioned[r.key] || r.kind == "" {
				b.forget(walk.Scope{Space: spaceTrigger, Deep: true})
			}
		}
	} else {
		if oldName, ok := walkName(v.Subname); ok {
			b.mayRemove(walk.Scope{Space: spaceTrigger, Deep: true, Prefix: oldName})
		} else {
			b.mayRemove(walk.Scope{Space: spaceTrigger, Deep: true})
		}
		if newName, ok := walkName(v.Newname); ok {
			b.mayAdd(walk.Scope{Space: spaceTrigger, Deep: true, Prefix: newName})
		} else {
			b.mayAdd(walk.Scope{Space: spaceTrigger, Deep: true})
		}
	}
	return true
}

// createPolicy follows CREATE POLICY name ON table.
func (b *builder) createPolicy(v *ast.CreatePolicyStmt) bool {
	r := b.exists(v.Table, false)
	if r.fail != "" || r.fact == walk.Absent {
		return true
	}
	if r.fact == walk.Present {
		b.need(r.key, walk.Present, []string{relTable}, ofKind(rangeOf(v.Table.Loc), relation(v.Table), relTable)) // 42809
		if r.kind == relTable {
			b.checkPolicyBoolExpr(r.key, v.Qual)
			b.checkPolicyBoolExpr(r.key, v.WithCheck)
		}
	}
	cmd := strings.ToLower(v.CmdName)
	if (cmd == "select" || cmd == "delete") && v.WithCheck != nil {
		b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("`WITH CHECK` on a `"+strings.ToUpper(cmd)+"` policy")) // 42601
		return true
	}
	if cmd == "insert" && v.Qual != nil {
		b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("`USING` on an `INSERT` policy")) // 42601
		return true
	}
	if r.fact == walk.Present {
		if pn, ok := walkName(v.PolicyName); ok {
			pk := polKey(r.key, pn)
			b.need(pk, walk.Absent, nil, &walkNote{rng: b.anchor(), present: alreadyExistsIn("policy", ident(v.PolicyName), relation(v.Table))}) // 42710
			b.set(pk, walk.Present, "")
		}
	} else {
		if pn, ok := walkName(v.PolicyName); ok {
			b.mayAdd(walk.Scope{Space: spacePolicy, Deep: true, Prefix: pn})
		} else {
			b.mayAdd(walk.Scope{Space: spacePolicy, Deep: true})
		}
	}
	return true
}

// alterPolicy follows ALTER POLICY name ON table.
func (b *builder) alterPolicy(v *ast.AlterPolicyStmt) bool {
	r := b.exists(v.Table, false)
	if r.fail != "" || r.fact == walk.Absent {
		return true
	}
	if r.fact == walk.Present {
		b.need(r.key, walk.Present, []string{relTable}, ofKind(rangeOf(v.Table.Loc), relation(v.Table), relTable)) // 42809
		if pn, ok := walkName(v.PolicyName); ok {
			b.need(polKey(r.key, pn), walk.Present, nil, &walkNote{rng: b.anchor(), absent: notFoundIn("policy", ident(v.PolicyName), relation(v.Table))}) // 42704
		}
		if r.kind == relTable {
			b.checkPolicyBoolExpr(r.key, v.Qual)
			b.checkPolicyBoolExpr(r.key, v.WithCheck)
		}
	}
	return true
}

// dropPolicy follows DROP POLICY [IF EXISTS] name ON table.
func (b *builder) dropPolicy(parts []string, missingOK bool) {
	if len(parts) < 2 {
		return
	}
	table := parts[:len(parts)-1]
	polWritten := parts[len(parts)-1]
	r := b.nameList(table)
	switch {
	case r.fail != "":
		b.fail(nil, walk.Absent, review.CrossDatabaseReference, r.fail)
		return
	case r.fact == walk.Absent:
		if !missingOK {
			b.fail(r.keys, walk.Absent, review.ObjectNotFound, missing(r, "", qualified(table))) // 42P01
		}
		return
	case r.fact == walk.Unknown:
		if pn, ok := walkName(polWritten); ok {
			b.mayRemove(walk.Scope{Space: spacePolicy, Deep: true, Prefix: pn})
		} else {
			b.mayRemove(walk.Scope{Space: spacePolicy, Deep: true})
		}
		return
	}
	if r.kind == relIndex || r.kind == relComposite {
		b.need(r.key, walk.Present, []string{relTable}, ofKind(b.anchor(), qualified(table), relTable)) // 42809
		return
	}
	if r.kind != "" && r.kind != relTable {
		if !missingOK {
			b.fail([]walk.Key{r.key}, walk.Present, review.ObjectNotFound, notFoundIn("policy", ident(polWritten), qualified(table))) // 42704
		}
		return
	}
	if pn, ok := walkName(polWritten); ok {
		pk := polKey(r.key, pn)
		if !missingOK {
			b.need(pk, walk.Present, nil, &walkNote{rng: b.anchor(), absent: notFoundIn("policy", ident(polWritten), qualified(table))}) // 42704
		}
		b.set(pk, walk.Absent, "")
	}
}

// renamePolicy follows ALTER POLICY old ON table RENAME TO new.
// PostgreSQL's rename_policy ignores MissingOk when opening the table.
func (b *builder) renamePolicy(v *ast.RenameStmt) bool {
	r := b.exists(v.Relation, false)
	if r.fail != "" || r.fact == walk.Absent {
		return true
	}
	if r.fact == walk.Present {
		b.need(r.key, walk.Present, []string{relTable}, ofKind(rangeOf(v.Relation.Loc), relation(v.Relation), relTable)) // 42809
		oldName, ok1 := walkName(v.Subname)
		newName, ok2 := walkName(v.Newname)
		if ok1 && ok2 {
			okKey := polKey(r.key, oldName)
			nkKey := polKey(r.key, newName)
			b.need(okKey, walk.Present, nil, &walkNote{rng: b.anchor(), absent: notFoundIn("policy", ident(v.Subname), relation(v.Relation))})
			b.need(nkKey, walk.Absent, nil, &walkNote{rng: b.anchor(), present: alreadyExistsIn("policy", ident(v.Newname), relation(v.Relation))})
			if oldName == newName {
				if f, _ := b.l.Lookup(okKey); f == walk.Unknown {
					b.failAt(b.anchor(), []walk.Key{r.key}, review.DuplicateName, "Policy "+review.Code(ident(v.Subname))+" on "+review.Code(relation(v.Relation))+" can't be renamed to itself.")
					return true
				}
			}
			b.set(okKey, walk.Absent, "")
			b.set(nkKey, walk.Present, "")
		} else {
			b.forget(walk.Scope{Space: spacePolicy, Parent: tableOf(r.key)})
		}
	} else {
		if oldName, ok := walkName(v.Subname); ok {
			b.mayRemove(walk.Scope{Space: spacePolicy, Deep: true, Prefix: oldName})
		} else {
			b.mayRemove(walk.Scope{Space: spacePolicy, Deep: true})
		}
		if newName, ok := walkName(v.Newname); ok {
			b.mayAdd(walk.Scope{Space: spacePolicy, Deep: true, Prefix: newName})
		} else {
			b.mayAdd(walk.Scope{Space: spacePolicy, Deep: true})
		}
	}
	return true
}

// alterObjectDepends follows ALTER ... [NO] DEPENDS ON EXTENSION.
func (b *builder) alterObjectDepends(v *ast.AlterObjectDependsStmt) bool {
	switch v.ObjectType {
	case ast.OBJECT_MATVIEW:
		b.existsOfKind(v.Relation, relMatview, "materialized view")
	case ast.OBJECT_INDEX:
		b.existsOfKind(v.Relation, relIndex, relIndex)
	case ast.OBJECT_TRIGGER:
		r := b.exists(v.Relation, false)
		if r.fact == walk.Present && r.kind != "" && r.kind != relTable && r.kind != relView && r.kind != relForeign {
			b.need(r.key, walk.Present, []string{relTable, relView, relForeign}, ofKind(rangeOf(v.Relation.Loc), relation(v.Relation), "table or view"))
		} else if r.fact == walk.Present {
			var trigName string
			if s, ok := v.Object.(*ast.String); ok {
				trigName = s.Str
			} else if parts := nameParts(listOf(v.Object)); len(parts) == 1 {
				trigName = parts[0]
			}
			if tn, ok := walkName(trigName); ok {
				b.need(trigKey(r.key, tn), walk.Present, nil, &walkNote{rng: b.anchor(), absent: notFoundIn("trigger", ident(trigName), relation(v.Relation))}) // 42704
			}
		}
	case ast.OBJECT_FUNCTION, ast.OBJECT_PROCEDURE, ast.OBJECT_ROUTINE:
		if owa, ok := v.Object.(*ast.ObjectWithArgs); ok {
			// ALTER FUNCTION / ROUTINE ... DEPENDS ON EXTENSION uses
			// LookupFuncWithArgs and accepts aggregates (prokind = 'a').
			b.checkRoutine(owa, v.ObjectType, false, v.ObjectType != ast.OBJECT_PROCEDURE)
		}
	}
	// DEPENDS and NO DEPENDS both look the extension up after the object
	// (42704). CREATE and DROP EXTENSION end the walk, so the synced list
	// holds throughout.
	//
	// pg: src/backend/commands/alter.c — ExecAlterObjectDependsStmt
	if s, ok := v.Extname.(*ast.String); ok {
		if name, ok := walkName(s.Str); ok {
			b.need(extensionKey(name), walk.Present, nil, &walkNote{rng: b.anchor(), absent: notFound("extension", ident(s.Str))}) // 42704
		}
	}
	return true
}

// alterSequence follows ALTER SEQUENCE.
func (b *builder) alterSequence(v *ast.AlterSeqStmt) bool {
	// Another kind of relation fails, under IF EXISTS too (42809,
	// verified on PostgreSQL 17).
	if r := b.exists(v.Sequence, v.MissingOk); r.fact == walk.Present {
		b.need(r.key, walk.Present, []string{relSequence}, ofKind(rangeOf(v.Sequence.Loc), relation(v.Sequence), relSequence))
		if r.kind == relSequence || r.kind == "" {
			b.checkSequenceOptions(v.Options, false)
		}
	}
	return true
}

// checkSequenceOptions validates unconditional CREATE/ALTER SEQUENCE option
// rules: duplicate options (42601), non-integer float literals (22P02),
// invalid OWNED BY / SEQUENCE NAME / AS type / bare LOGGED/UNLOGGED.
func (b *builder) checkSequenceOptions(opts *ast.List, isIdentity bool) {
	if opts == nil {
		return
	}
	seen := make(map[string]bool)
	for _, item := range opts.Items {
		de, ok := item.(*ast.DefElem)
		if !ok {
			continue
		}
		name := strings.ToLower(de.Defname)
		if seen[name] {
			if isIdentity {
				b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("conflicting or redundant options in identity column statement")) // 42601 / 55000
			} else {
				b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("conflicting or redundant options in sequence statement")) // 42601
			}
			return
		}
		seen[name] = true
		switch name {
		case "increment", "minvalue", "maxvalue", "start", "restart", "cache":
			if f, ok := de.Arg.(*ast.Float); ok && strings.ContainsAny(f.Fval, ".eE") {
				if isIdentity {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("non-integer value "+review.Code(f.Fval)+" for identity option "+review.Code(name))) // 22P02 / 55000
				} else {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("non-integer value "+review.Code(f.Fval)+" for sequence option "+review.Code(name))) // 22P02
				}
				return
			}
		case "sequence_name":
			if !isIdentity {
				b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("`SEQUENCE NAME` in `CREATE`/`ALTER SEQUENCE`")) // 42601
				return
			}
		case "logged":
			if !isIdentity && de.Arg == nil {
				b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("bare `LOGGED` in `CREATE`/`ALTER SEQUENCE`")) // XX000
				return
			}
		case "unlogged":
			if !isIdentity {
				b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("bare `UNLOGGED` in `CREATE`/`ALTER SEQUENCE`")) // XX000
				return
			}
		case "owned_by":
			if !isIdentity {
				if list, ok := de.Arg.(*ast.List); ok {
					parts := nameParts(list)
					switch len(parts) {
					case 1:
						if !strings.EqualFold(parts[0], "none") {
							b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("unqualified `OWNED BY` in sequence statement")) // 42601
							return
						}
					case 2, 3:
					case 4:
						if b.w.crossDatabase(parts[0]) {
							b.fail(nil, walk.Absent, review.CrossDatabaseReference, crossDatabaseMessage(qualified(parts))) // 0A000
							return
						}
					default:
						if len(parts) > 4 {
							b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("`OWNED BY` target "+review.Code(qualified(parts))+" in sequence statement")) // 42601
							return
						}
					}
				}
			}
		case "as":
			if !isIdentity {
				if tn, ok := de.Arg.(*ast.TypeName); ok && !b.checkSequenceType(tn) {
					return
				}
			}
		}
	}
}

func (b *builder) checkSequenceType(tn *ast.TypeName) bool {
	if tn == nil || tn.ArrayBounds != nil {
		b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("sequence data type "+review.Code(typeName(tn)))) // 42704 / 0A000
		return false
	}
	parts := nameParts(tn.Names)
	var name string
	switch len(parts) {
	case 1:
		name = parts[0]
	case 2:
		if parts[0] != "pg_catalog" {
			b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("sequence data type "+review.Code(typeName(tn)))) // 42704 / 0A000
			return false
		}
		name = parts[1]
	case 3:
		if b.w.crossDatabase(parts[0]) {
			b.fail(nil, walk.Absent, review.CrossDatabaseReference, crossDatabaseMessage(typeName(tn))) // 0A000
			return false
		}
		if parts[1] != "pg_catalog" {
			b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("sequence data type "+review.Code(typeName(tn)))) // 42704 / 0A000
			return false
		}
		name = parts[2]
	default:
		b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("sequence data type "+review.Code(typeName(tn)))) // 42704 / 0A000
		return false
	}
	switch strings.ToLower(name) {
	case "smallint", "int2", "integer", "int", "int4", "bigint", "int8":
		return true
	}
	b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("sequence data type "+review.Code(typeName(tn)))) // 42704 / 0A000
	return false
}

// atAllowedKinds returns the relation kinds PostgreSQL's ATSimplePermissions
// accepts for sub, and the word a wrong-kind message uses, or nil when sub is
// unrestricted or not judged.
//
// pg: src/backend/commands/tablecmds.c — ATSimplePermissions
func atAllowedKinds(sub ast.AlterTableType) ([]string, string) {
	switch sub {
	case ast.AT_GenericOptions, ast.AT_AlterColumnGenericOptions:
		return []string{relForeign}, "foreign table"
	case ast.AT_AddColumn, ast.AT_DropColumn, ast.AT_AlterColumnType,
		ast.AT_SetNotNull, ast.AT_DropNotNull,
		ast.AT_AddConstraint, ast.AT_DropConstraint, ast.AT_ValidateConstraint,
		ast.AT_AddInherit, ast.AT_DropInherit, ast.AT_DropOids,
		ast.AT_EnableTrig, ast.AT_EnableAlwaysTrig, ast.AT_EnableReplicaTrig, ast.AT_DisableTrig,
		ast.AT_EnableTrigAll, ast.AT_DisableTrigAll, ast.AT_EnableTrigUser, ast.AT_DisableTrigUser:
		return []string{relTable, relForeign}, relTable
	case ast.AT_AlterConstraint, ast.AT_SetExpression, ast.AT_DropExpression,
		ast.AT_EnableRowSecurity, ast.AT_DisableRowSecurity, ast.AT_ForceRowSecurity, ast.AT_NoForceRowSecurity,
		ast.AT_DetachPartition, ast.AT_DetachPartitionFinalize,
		ast.AT_AddOf, ast.AT_DropOf,
		ast.AT_EnableRule, ast.AT_EnableAlwaysRule, ast.AT_EnableReplicaRule, ast.AT_DisableRule:
		return []string{relTable}, relTable
	case ast.AT_AttachPartition:
		return []string{relTable, relIndex}, "table or index"
	case ast.AT_SetLogged, ast.AT_SetUnLogged:
		return []string{relTable, relSequence}, "table or sequence"
	case ast.AT_ClusterOn, ast.AT_DropCluster, ast.AT_ReplicaIdentity, ast.AT_SetAccessMethod:
		return []string{relTable, relMatview}, "table or materialized view"
	case ast.AT_SetTableSpace:
		return []string{relTable, relMatview, relIndex}, "table, materialized view, or index"
	case ast.AT_ColumnDefault, ast.AT_AddIdentity, ast.AT_SetIdentity, ast.AT_DropIdentity:
		return []string{relTable, relView, relForeign}, "table or view"
	case ast.AT_SetStatistics:
		return []string{relTable, relMatview, relForeign, relIndex}, "table, materialized view, or index"
	case ast.AT_SetOptions, ast.AT_ResetOptions, ast.AT_SetStorage, ast.AT_SetCompression:
		return []string{relTable, relMatview, relForeign}, "table or materialized view"
	case ast.AT_SetRelOptions, ast.AT_ResetRelOptions, ast.AT_ReplaceRelOptions:
		return []string{relTable, relView, relMatview, relIndex}, "table, view, materialized view, or index"
	}
	return nil, ""
}

// checkAlterTargetSubcmd checks the secondary object named by ALTER TABLE
// subcommands like CLUSTER ON, REPLICA IDENTITY USING INDEX, OF, INHERIT,
// NO INHERIT, ATTACH PARTITION, and DETACH PARTITION.
func (b *builder) checkAlterTargetSubcmd(rv *ast.RangeVar, r resolved, cmd *ast.AlterTableCmd) {
	sub := ast.AlterTableType(cmd.Subtype)
	rng := rangeOf(cmd.Loc)
	s := schemaOfKey(r.key)
	switch sub {
	case ast.AT_ClusterOn:
		if idx, ok := walkName(cmd.Name); ok {
			ik := relKey(s, idx)
			note := ofKind(rng, ident(cmd.Name), relIndex)
			note.absent = notFound(relIndex, ident(cmd.Name))
			b.need(ik, walk.Present, []string{relIndex}, note) // 42704 / 42809
		}
	case ast.AT_DropCluster:
		if r.kind == relTable && b.w.partitioned[r.key] {
			b.failAt(rng, []walk.Key{r.key}, review.UnsupportedForm, unsupported("`SET WITHOUT CLUSTER` on partitioned table "+review.Code(relation(rv)))) // 0A000
		}
	case ast.AT_EnableTrig, ast.AT_EnableAlwaysTrig, ast.AT_EnableReplicaTrig, ast.AT_DisableTrig:
		if tn, ok := walkName(cmd.Name); ok {
			b.need(trigKey(r.key, tn), walk.Present, nil, &walkNote{rng: rng, absent: notFoundIn("trigger", ident(cmd.Name), relation(rv))}) // 42704
		}
	case ast.AT_EnableRule, ast.AT_EnableAlwaysRule, ast.AT_EnableReplicaRule, ast.AT_DisableRule:
		// A rule belongs to its table alone: no inheritance child or
		// partition shares it (pg: rewriteDefine.c — EnableDisableRule).
		if rn, ok := walkName(cmd.Name); ok {
			b.need(ruleKey(r.key, rn), walk.Present, nil, &walkNote{rng: rng, absent: notFoundIn("rule", ident(cmd.Name), relation(rv))}) // 42704
		}
	case ast.AT_SetRelOptions, ast.AT_ReplaceRelOptions:
		if opts, ok := cmd.Def.(*ast.List); ok {
			if r.kind == relTable && b.w.partitioned[r.key] && hasNonToastRelOption(opts) {
				b.failAt(rng, []walk.Key{r.key}, review.WrongObjectKind, unsupported("storage parameters on partitioned table "+review.Code(relation(rv)))) // 42809
				return
			}
			b.checkRelOptions(r.kind, opts, rng)
			if sub == ast.AT_SetRelOptions && (r.kind == relTable || r.kind == relMatview) {
				b.checkHeapRelOptions(opts, rng)
			}
			if r.kind == relIndex && !b.w.isTableDirty(r.key) && b.w.pg != nil {
				if idx := b.w.pg.indexDefs[[2]string{schemaOfKey(r.key), r.key.Name}]; idx != nil && idx.accessMethod == "btree" {
					b.checkBtreeRelOptions(opts, rng)
				}
			}
		}
	case ast.AT_SetOptions:
		if opts, ok := cmd.Def.(*ast.List); ok {
			b.checkAttrOptions(opts, rng)
		}
	case ast.AT_SetStorage:
		if s, ok := cmd.Def.(*ast.String); ok {
			b.checkStorageMode(s.Str, rng)
		}
	case ast.AT_SetCompression:
		if s, ok := cmd.Def.(*ast.String); ok {
			b.checkCompressionMethod(s.Str, rng)
		}
	case ast.AT_SetStatistics:
		if cmd.Num != 0 && r.kind != relIndex {
			b.failAt(rng, []walk.Key{r.key}, review.UnsupportedForm, unsupported("`ALTER TABLE ... ALTER COLUMN` by column number")) // 0A000
		} else if r.kind == relIndex && !b.w.isTableDirty(r.key) && b.w.pg != nil {
			if idx := b.w.pg.indexDefs[[2]string{schemaOfKey(r.key), r.key.Name}]; idx != nil {
				if cmd.Num > 0 {
					if int(cmd.Num) > len(idx.isExpr) {
						b.failAt(rng, []walk.Key{r.key}, review.ObjectNotFound, notFoundIn("column", strconv.Itoa(int(cmd.Num)), relation(rv))) // 42703
					} else if !idx.isExpr[cmd.Num-1] {
						b.failAt(rng, []walk.Key{r.key}, review.UnsupportedForm, unsupported("`SET STATISTICS` on non-expression index column")) // 0A000
					}
				} else if cmd.Name != "" && !idx.hasExpr {
					if col, ok := walkName(cmd.Name); ok && !slices.Contains(idx.colNames, col) {
						b.failAt(rng, []walk.Key{r.key}, review.ObjectNotFound, notFoundIn("column", ident(cmd.Name), relation(rv))) // 42703
					}
				}
			}
		}
	case ast.AT_AddIdentity, ast.AT_SetIdentity, ast.AT_DropIdentity:
		if r.kind == relView && (sub != ast.AT_DropIdentity || !cmd.Missing_ok) {
			b.failAt(rng, []walk.Key{r.key}, review.UnsupportedForm, unsupported("identity column action on view "+review.Code(relation(rv)))) // 55000 / 42809
		}
	case ast.AT_DropOf:
		if r.kind == relTable && b.w.partitioned[r.key] {
			b.failAt(rng, []walk.Key{r.key}, review.WrongObjectKind, unsupported("`NOT OF` on partitioned table "+review.Code(relation(rv)))) // 42809
		}
	case ast.AT_ReplicaIdentity:
		if ri, ok := cmd.Def.(*ast.ReplicaIdentityStmt); ok && ri.IdentityType == 'i' && ri.Name != "" {
			if idx, ok := walkName(ri.Name); ok {
				ik := relKey(s, idx)
				note := ofKind(rng, ident(ri.Name), relIndex)
				note.absent = notFound(relIndex, ident(ri.Name))
				b.need(ik, walk.Present, []string{relIndex}, note) // 42704 / 42809
			}
		}
	case ast.AT_AddOf:
		if tn, ok := cmd.Def.(*ast.TypeName); ok {
			b.checkCompositeType(tn, rng) // 42809 / 42704
		}
	case ast.AT_AddInherit:
		if parentRV, ok := cmd.Def.(*ast.RangeVar); ok {
			pr := b.exists(parentRV, false) // 42P01
			if pr.fact == walk.Present {
				b.need(pr.key, walk.Present, []string{relTable}, ofKind(rangeOf(parentRV.Loc), relation(parentRV), relTable)) // 42809
				if pr.key == r.key {
					b.failAt(rng, []walk.Key{r.key}, review.UnsupportedForm, unsupported("circular inheritance on "+review.Code(relation(rv)))) // 42P07
				}
			}
		}
	case ast.AT_DropInherit:
		if parentRV, ok := cmd.Def.(*ast.RangeVar); ok {
			pr := b.exists(parentRV, false) // 42P01
			if pr.fact == walk.Present && pr.key == r.key {
				b.failAt(rng, []walk.Key{r.key}, review.UnsupportedForm, unsupported("`NO INHERIT` of "+review.Code(relation(rv))+" from itself")) // 42P01
			}
		}
	case ast.AT_AttachPartition:
		if pc, ok := cmd.Def.(*ast.PartitionCmd); ok && pc.Name != nil {
			pr := b.exists(pc.Name, false) // 42P01
			if pr.fact == walk.Present {
				if r.kind == relIndex {
					b.need(pr.key, walk.Present, []string{relIndex}, ofKind(rangeOf(pc.Name.Loc), relation(pc.Name), relIndex)) // 42809
				} else {
					b.need(pr.key, walk.Present, []string{relTable, relForeign}, ofKind(rangeOf(pc.Name.Loc), relation(pc.Name), relTable)) // 42809
				}
				if pr.key == r.key {
					b.failAt(rng, []walk.Key{r.key}, review.UnsupportedForm, unsupported("attaching "+review.Code(relation(rv))+" as a partition of itself")) // 42P17
				}
			}
		}
	case ast.AT_DetachPartition, ast.AT_DetachPartitionFinalize:
		if pc, ok := cmd.Def.(*ast.PartitionCmd); ok && pc.Name != nil {
			pr := b.exists(pc.Name, false) // 42P01
			if pr.fact == walk.Present && pr.key == r.key {
				b.failAt(rng, []walk.Key{r.key}, review.UnsupportedForm, unsupported("detaching "+review.Code(relation(rv))+" from itself")) // 42P17 / 42P01
			}
		}
	}
}

// checkConstraintColumns checks that columns named in ADD PRIMARY KEY / UNIQUE
// / FOREIGN KEY exist in the table (42703).
func (b *builder) checkConstraintColumns(t walk.Key, rv *ast.RangeVar, con *ast.Constraint, rng review.Range) {
	checkList := func(list *ast.List) {
		if list == nil {
			return
		}
		for _, item := range list.Items {
			s, ok := item.(*ast.String)
			if !ok {
				continue
			}
			if sk, ok := storedKey(s.Str); ok && b.added[sk] {
				continue
			}
			if n, ok := walkName(s.Str); ok && !systemColumns[n] {
				b.need(columnKey(t, n), walk.Present, nil, &walkNote{rng: rng, absent: notFoundIn("column", ident(s.Str), relation(rv))}) // 42703
			}
		}
	}
	switch con.Contype {
	case ast.CONSTR_PRIMARY, ast.CONSTR_UNIQUE:
		checkList(con.Keys)
		checkList(con.Including)
	case ast.CONSTR_FOREIGN:
		checkList(con.FkAttrs)
	case ast.CONSTR_EXCLUSION:
		if con.Exclusions != nil {
			for _, item := range con.Exclusions.Items {
				pair, ok := item.(*ast.List)
				if !ok || len(pair.Items) == 0 {
					continue
				}
				elem, ok := pair.Items[0].(*ast.IndexElem)
				if !ok || elem.Name == "" || elem.Expr != nil {
					continue
				}
				if sk, ok := storedKey(elem.Name); ok && b.added[sk] {
					continue
				}
				if n, ok := walkName(elem.Name); ok && !systemColumns[n] {
					b.need(columnKey(t, n), walk.Present, nil, &walkNote{rng: rangeOf(elem.Loc), absent: notFoundIn("column", ident(elem.Name), relation(rv))}) // 42703
				}
			}
		}
	}
}

// checkIndexElems checks column existence (42703), non-immutable SQL value
// functions (42P17), and invalid INCLUDE elements (42P17) on CREATE INDEX.
func (b *builder) checkIndexElems(r resolved, v *ast.IndexStmt) {
	if v.IndexParams != nil {
		for _, item := range v.IndexParams.Items {
			elem, ok := item.(*ast.IndexElem)
			if !ok {
				continue
			}
			if elem.Name != "" && elem.Expr == nil && hasColumns(r.kind) {
				if col, ok := walkName(elem.Name); ok && !systemColumns[col] {
					b.need(columnKey(r.key, col), walk.Present, nil, &walkNote{rng: rangeOf(elem.Loc), absent: notFoundIn("column", ident(elem.Name), relation(v.Relation))}) // 42703
				}
			}
			if elem.Expr != nil && isDirectSQLValueFunction(elem.Expr) {
				b.failAt(rangeOf(elem.Loc), nil, review.UnsupportedForm, unsupported("non-immutable expression in index")) // 42P17
				return
			}
		}
	}
	if v.IndexIncludingParams != nil {
		for _, item := range v.IndexIncludingParams.Items {
			elem, ok := item.(*ast.IndexElem)
			if !ok {
				continue
			}
			if elem.Expr != nil || elem.Opclass != nil || elem.Collation != nil ||
				elem.Ordering != ast.SORTBY_DEFAULT || elem.NullsOrdering != ast.SORTBY_NULLS_DEFAULT {
				b.failAt(rangeOf(elem.Loc), nil, review.UnsupportedForm, unsupported("expressions or options in index `INCLUDE` list")) // 42P17
				return
			}
			if elem.Name != "" && hasColumns(r.kind) {
				if col, ok := walkName(elem.Name); ok && !systemColumns[col] {
					b.need(columnKey(r.key, col), walk.Present, nil, &walkNote{rng: rangeOf(elem.Loc), absent: notFoundIn("column", ident(elem.Name), relation(v.Relation))}) // 42703
				}
			}
		}
	}
}

// isDirectSQLValueFunction reports whether n is a bare SQLValueFunction
// (CURRENT_DATE, CURRENT_ROLE, etc.) that cannot be eliminated by PostgreSQL's
// constant-folding simplification of unreachable branches or strict-NULL ops.
func isDirectSQLValueFunction(n ast.Node) bool {
	_, ok := n.(*ast.SQLValueFunction)
	return ok
}

// checkViewQuery checks unconditional CREATE VIEW query rules: WITH CHECK
// OPTION without a FROM clause (0A000), unbound columns in a FROM-less SELECT
// (42703), too many column aliases (42601), CREATE RECURSIVE VIEW with too
// many column aliases (42P10), and CREATE OR REPLACE VIEW dropping synced view
// columns (42P16).
func (b *builder) checkViewQuery(v *ast.ViewStmt, k walk.Key, prior walk.Fact) {
	sel, ok := v.Query.(*ast.SelectStmt)
	if !ok {
		return
	}
	if selectHasUnboundColumn(sel) {
		b.fail(nil, walk.Absent, review.ObjectNotFound, "Query references a column with no `FROM` clause.") // 42703
		return
	}
	if sel.IntoClause != nil {
		// pg: src/backend/parser/analyze.c — transformSelectStmt's caller,
		// "views must not contain SELECT INTO"
		b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("`SELECT INTO` in a view")) // 0A000
		return
	}
	if b.queryRefused(sel) {
		return
	}
	if v.WithCheckOption != 0 && sel.Op == ast.SETOP_NONE && (sel.FromClause == nil || len(sel.FromClause.Items) == 0) {
		b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("`WITH CHECK OPTION` on a non-updatable view")) // 0A000
		return
	}
	if sel.WithClause != nil && sel.WithClause.Recursive && sel.WithClause.Ctes != nil && len(sel.WithClause.Ctes.Items) == 1 {
		if cte, ok := sel.WithClause.Ctes.Items[0].(*ast.CommonTableExpr); ok {
			if selectHasUnboundColumn(cte.Ctequery) {
				b.fail(nil, walk.Absent, review.ObjectNotFound, "Query references a column with no `FROM` clause.") // 42703
				return
			}
			if cte.Aliascolnames != nil && len(cte.Aliascolnames.Items) > 0 {
				if inner, ok := cte.Ctequery.(*ast.SelectStmt); ok && inner.Op == ast.SETOP_NONE && inner.ValuesLists == nil {
					if n, exact := selectTargetCount(inner); exact && len(cte.Aliascolnames.Items) > n {
						b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("more column names than columns in `CREATE RECURSIVE VIEW`")) // 42P10
						return
					}
				}
			}
		}
	}
	if v.Aliases != nil && len(v.Aliases.Items) > 0 && sel.Op == ast.SETOP_NONE && sel.ValuesLists == nil {
		if n, exact := selectTargetCount(sel); exact && len(v.Aliases.Items) > n {
			b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("more column names than columns in `CREATE VIEW`")) // 42601
			return
		}
	}
	if v.Replace && prior == walk.Present && !b.w.modifiedViews[k] && b.w.pg != nil && sel.Op == ast.SETOP_NONE && sel.ValuesLists == nil {
		if oldCols := b.w.pg.viewCols[[2]string{schemaOfKey(k), k.Name}]; oldCols > 0 {
			if n, exact := selectTargetCount(sel); exact && n < oldCols {
				b.fail([]walk.Key{k}, walk.Present, review.UnsupportedForm, unsupported("dropping columns in `CREATE OR REPLACE VIEW`")) // 42P16
				return
			}
		}
	}
}

// selectTargetCount returns the exact output column count of a plain SELECT
// with no `*` star expansion, or (0, false) when it cannot be known statically.
func selectTargetCount(sel *ast.SelectStmt) (int, bool) {
	if sel.TargetList == nil {
		return 0, true
	}
	for _, item := range sel.TargetList.Items {
		rt, ok := item.(*ast.ResTarget)
		if !ok {
			return 0, false
		}
		switch v := rt.Val.(type) {
		case *ast.ColumnRef:
			if v.Fields != nil {
				for _, f := range v.Fields.Items {
					if _, star := f.(*ast.A_Star); star {
						return 0, false
					}
				}
			}
		case *ast.A_Indirection:
			if v.Indirection != nil {
				for _, f := range v.Indirection.Items {
					if _, star := f.(*ast.A_Star); star {
						return 0, false
					}
				}
			}
		}
	}
	return len(sel.TargetList.Items), true
}

// renameDomain follows ALTER DOMAIN ... RENAME TO and RENAME CONSTRAINT.
func (b *builder) renameDomain(v *ast.RenameStmt) bool {
	var parts []string
	if tn, ok := v.Object.(*ast.TypeName); ok {
		parts = nameParts(tn.Names)
	} else {
		parts = nameParts(listOf(v.Object))
	}
	if len(parts) == 0 {
		return false
	}
	r := b.checkDomain(parts, b.anchor())
	if v.RenameType == ast.OBJECT_DOMCONSTRAINT {
		if oldC, ok1 := storedKey(v.Subname); ok1 {
			if newC, ok2 := storedKey(v.Newname); ok2 && oldC == newC && r.fail == "" && r.fact != walk.Absent && !r.builtin {
				b.failAt(b.anchor(), nil, review.DuplicateName, "Constraint "+review.Code(ident(v.Subname))+" on domain "+review.Code(qualified(parts))+" can't be renamed to itself.")
			}
		}
		return true
	}
	if r.fail != "" || r.fact == walk.Absent || r.builtin {
		return true
	}
	if r.fact == walk.Unknown {
		if oldN, ok1 := storedKey(parts[len(parts)-1]); ok1 {
			if newN, ok2 := storedKey(v.Newname); ok2 && oldN == newN {
				b.failAt(b.anchor(), nil, review.DuplicateName, "Domain "+review.Code(qualified(parts))+" can't be renamed to itself.")
				return true
			}
		}
		b.mayRemove(walk.Scope{Space: spaceType, Deep: true})
		b.mayAdd(walk.Scope{Space: spaceType, Deep: true})
		return true
	}
	newName, ok := walkName(v.Newname)
	if !ok {
		b.forget(walk.Scope{Space: spaceType, Parent: walk.Path(schemaOfKey(r.key))})
		return true
	}
	s := schemaOfKey(r.key)
	nk := typeKey(s, newName)
	b.need(nk, walk.Absent, nil, &walkNote{rng: b.anchor(), present: alreadyExists("type", ident(v.Newname))}) // 42710
	b.mayRemoveArrayType(s)
	b.mayAddArrayType(s)
	b.set(r.key, walk.Absent, "")
	b.set(nk, walk.Present, typDomain)
	b.made(nk)
	return true
}

// renameType follows ALTER TYPE ... RENAME TO.
func (b *builder) renameType(v *ast.RenameStmt) bool {
	var parts []string
	if tn, ok := v.Object.(*ast.TypeName); ok {
		parts = nameParts(tn.Names)
	} else {
		parts = nameParts(listOf(v.Object))
	}
	if len(parts) == 0 {
		return false
	}
	written := qualified(parts)
	rng := b.anchor()
	r := b.resolveType(parts)
	switch {
	case r.fail != "":
		b.failAt(rng, nil, review.CrossDatabaseReference, r.fail)
		return true
	case r.fact == walk.Absent:
		if r.missingSchema != "" {
			b.failAt(rng, r.keys, review.ObjectNotFound, schemaNotFound(r.missingSchema))
		} else {
			b.failAt(rng, r.keys, review.ObjectNotFound, notFound("type", written))
		}
		return true
	case r.fact == walk.Unknown:
		if oldN, ok1 := storedKey(parts[len(parts)-1]); ok1 {
			if newN, ok2 := storedKey(v.Newname); ok2 && oldN == newN {
				b.failAt(rng, nil, review.DuplicateName, "Type "+review.Code(written)+" can't be renamed to itself.")
				return true
			}
		}
		b.mayRemove(walk.Scope{Space: spaceType, Deep: true})
		b.mayAdd(walk.Scope{Space: spaceType, Deep: true})
		b.mayRemove(walk.Scope{Space: spaceRelation, Deep: true, Kinds: []string{relComposite}})
		b.mayAdd(walk.Scope{Space: spaceRelation, Deep: true, Kinds: []string{relComposite}})
		b.mayRemove(walk.Scope{Space: spaceFunction, Deep: true})
		b.mayAdd(walk.Scope{Space: spaceFunction, Deep: true})
		return true
	case r.builtin:
		return true
	}
	s := schemaOfKey(r.key)
	rk := relKey(s, r.key.Name)
	rf, rkind := b.l.Lookup(rk)
	if rf == walk.Present && hasRowType(rkind) && rkind != relComposite {
		b.failAt(rng, []walk.Key{rk}, review.WrongObjectKind, review.Code(written)+" is "+article(rkind)+", not a type.") // 42809
		return true
	}
	newName, ok := walkName(v.Newname)
	if !ok {
		b.forget(walk.Scope{Space: spaceType, Parent: walk.Path(s)})
		if rf == walk.Present && rkind == relComposite {
			b.forget(walk.Scope{Space: spaceRelation, Parent: walk.Path(s), Kinds: []string{relComposite}})
		}
		if mrKey, hasMR := b.w.rangeMultirange[r.key]; hasMR {
			if mf, _ := b.l.Lookup(mrKey); mf == walk.Present {
				b.set(mrKey, walk.Unknown, "")
			}
		}
		b.session(func() {
			delete(b.w.shellTypes, r.key)
			delete(b.w.typmodTypes, r.key)
			delete(b.w.rangeMultirange, r.key)
			for k, mr := range b.w.rangeMultirange {
				if mr == r.key {
					delete(b.w.rangeMultirange, k)
				}
			}
		})
		return true
	}
	nk := typeKey(s, newName)
	b.need(nk, walk.Absent, nil, &walkNote{rng: rng, present: alreadyExists("type", ident(v.Newname))}) // 42710
	if rf == walk.Present && rkind == relComposite {
		nrk := relKey(s, newName)
		b.need(nrk, walk.Absent, nil, taken(rng, relComposite, ident(v.Newname))) // 42P07
		b.set(rk, walk.Absent, "")
		b.set(nrk, walk.Present, relComposite)
		b.renamed(rk, nrk)
	} else if rf == walk.Unknown {
		b.mayAdd(walk.Scope{Space: spaceRelation, Parent: walk.Path(s), Kinds: []string{relComposite}})
	}
	b.mayRemoveArrayType(s)
	b.mayAddArrayType(s)
	b.mayRemove(walk.Scope{Space: spaceFunction, Parent: walk.Path(s)})
	b.mayAdd(walk.Scope{Space: spaceFunction, Parent: walk.Path(s)})
	b.set(r.key, walk.Absent, "")
	b.set(nk, walk.Present, r.kind)
	b.session(func() {
		if b.w.shellTypes[r.key] {
			delete(b.w.shellTypes, r.key)
			b.w.shellTypes[nk] = true
		}
		if b.w.typmodTypes[r.key] {
			delete(b.w.typmodTypes, r.key)
			b.w.typmodTypes[nk] = true
		}
		if mrKey, hasMR := b.w.rangeMultirange[r.key]; hasMR {
			delete(b.w.rangeMultirange, r.key)
			b.w.rangeMultirange[nk] = mrKey
		}
		for k, mr := range b.w.rangeMultirange {
			if mr == r.key {
				b.w.rangeMultirange[k] = nk
			}
		}
	})
	b.made(nk)
	return true
}

// setSchemaDomain follows ALTER DOMAIN ... SET SCHEMA.
func (b *builder) setSchemaDomain(v *ast.AlterObjectSchemaStmt) bool {
	var parts []string
	if tn, ok := v.Object.(*ast.TypeName); ok {
		parts = nameParts(tn.Names)
	} else {
		parts = nameParts(listOf(v.Object))
	}
	if len(parts) == 0 {
		return false
	}
	r := b.checkDomain(parts, b.anchor())
	if r.fail != "" || r.fact == walk.Absent || r.builtin {
		return true
	}
	ns, ok := walkName(v.Newschema)
	if !ok || ns == tempSchema || ns == "pg_catalog" {
		return false
	}
	b.need(schemaKey(ns), walk.Present, nil, &walkNote{rng: b.anchor(), absent: schemaNotFound(v.Newschema)}) // 3F000
	if r.fact == walk.Unknown {
		b.mayRemove(walk.Scope{Space: spaceType, Deep: true})
		b.mayAdd(walk.Scope{Space: spaceType, Parent: walk.Path(ns)})
		return true
	}
	old := schemaOfKey(r.key)
	if ns == old {
		return true
	}
	nk := typeKey(ns, r.key.Name)
	moved := qualified([]string{v.Newschema, r.key.Name})
	b.need(nk, walk.Absent, nil, &walkNote{rng: b.anchor(), present: alreadyExists("type", moved)}) // 42710
	b.mayRemoveArrayType(old)
	b.mayAddArrayType(ns)
	b.set(r.key, walk.Absent, "")
	b.set(nk, walk.Present, typDomain)
	b.made(nk)
	return true
}

// setSchemaType follows ALTER TYPE ... SET SCHEMA.
func (b *builder) setSchemaType(v *ast.AlterObjectSchemaStmt) bool {
	var parts []string
	if tn, ok := v.Object.(*ast.TypeName); ok {
		parts = nameParts(tn.Names)
	} else {
		parts = nameParts(listOf(v.Object))
	}
	if len(parts) == 0 {
		return false
	}
	written := qualified(parts)
	rng := b.anchor()
	r := b.resolveType(parts)
	switch {
	case r.fail != "":
		b.failAt(rng, nil, review.CrossDatabaseReference, r.fail)
		return true
	case r.fact == walk.Absent:
		if r.missingSchema != "" {
			b.failAt(rng, r.keys, review.ObjectNotFound, schemaNotFound(r.missingSchema))
		} else {
			b.failAt(rng, r.keys, review.ObjectNotFound, notFound("type", written))
		}
		return true
	}
	ns, ok := walkName(v.Newschema)
	if !ok || ns == tempSchema || ns == "pg_catalog" {
		return false
	}
	if r.fact == walk.Present && !r.builtin {
		old := schemaOfKey(r.key)
		rk := relKey(old, r.key.Name)
		rf, rkind := b.l.Lookup(rk)
		if rf == walk.Present && hasRowType(rkind) && rkind != relComposite {
			b.failAt(rng, []walk.Key{rk}, review.WrongObjectKind, review.Code(written)+" is "+article(rkind)+", not a type.") // 42809
			return true
		}
		b.need(schemaKey(ns), walk.Present, nil, &walkNote{rng: rng, absent: schemaNotFound(v.Newschema)}) // 3F000
		if ns == old {
			return true
		}
		nk := typeKey(ns, r.key.Name)
		moved := qualified([]string{v.Newschema, r.key.Name})
		b.need(nk, walk.Absent, nil, &walkNote{rng: rng, present: alreadyExists("type", moved)}) // 42710
		if rf == walk.Present && rkind == relComposite {
			nrk := relKey(ns, r.key.Name)
			b.need(nrk, walk.Absent, nil, taken(rng, relComposite, moved)) // 42P07
			b.set(rk, walk.Absent, "")
			b.set(nrk, walk.Present, relComposite)
			b.renamed(rk, nrk)
		} else if rf == walk.Unknown {
			b.mayAdd(walk.Scope{Space: spaceRelation, Parent: walk.Path(ns), Kinds: []string{relComposite}})
		}
		// Moving a type also moves its companion array type (and companion
		// multirange type and constructor functions for range types) into
		// the destination schema.
		b.mayRemove(walk.Scope{Space: spaceType, Parent: walk.Path(old)})
		b.mayAdd(walk.Scope{Space: spaceType, Parent: walk.Path(ns)})
		b.mayRemove(walk.Scope{Space: spaceFunction, Parent: walk.Path(old)})
		b.mayAdd(walk.Scope{Space: spaceFunction, Parent: walk.Path(ns)})
		if mrKey, hasMR := b.w.rangeMultirange[r.key]; hasMR && schemaOfKey(mrKey) != old {
			if mf, _ := b.l.Lookup(mrKey); mf == walk.Present {
				b.set(mrKey, walk.Unknown, "")
			}
			b.mayRemove(walk.Scope{Space: spaceFunction, Parent: mrKey.Parent})
		}
		b.set(r.key, walk.Absent, "")
		b.set(nk, walk.Present, r.kind)
		b.session(func() {
			if b.w.shellTypes[r.key] {
				delete(b.w.shellTypes, r.key)
				b.w.shellTypes[nk] = true
			}
			if b.w.typmodTypes[r.key] {
				delete(b.w.typmodTypes, r.key)
				b.w.typmodTypes[nk] = true
			}
			if mrKey, hasMR := b.w.rangeMultirange[r.key]; hasMR {
				delete(b.w.rangeMultirange, r.key)
				b.w.rangeMultirange[nk] = mrKey
			}
			for k, mr := range b.w.rangeMultirange {
				if mr == r.key {
					b.w.rangeMultirange[k] = nk
				}
			}
		})
		b.made(nk)
		return true
	}
	b.need(schemaKey(ns), walk.Present, nil, &walkNote{rng: rng, absent: schemaNotFound(v.Newschema)}) // 3F000
	if r.builtin {
		return true
	}
	b.mayRemove(walk.Scope{Space: spaceType, Deep: true})
	b.mayAdd(walk.Scope{Space: spaceType, Parent: walk.Path(ns)})
	b.mayRemove(walk.Scope{Space: spaceRelation, Deep: true, Kinds: []string{relComposite}})
	b.mayAdd(walk.Scope{Space: spaceRelation, Parent: walk.Path(ns), Kinds: []string{relComposite}})
	b.mayRemove(walk.Scope{Space: spaceFunction, Deep: true})
	b.mayAdd(walk.Scope{Space: spaceFunction, Parent: walk.Path(ns)})
	return true
}

// routineWord returns the English kind word for a routine ObjectType.
func routineWord(objType ast.ObjectType) string {
	switch objType {
	case ast.OBJECT_PROCEDURE:
		return "procedure"
	case ast.OBJECT_ROUTINE:
		return "routine"
	default:
		return "function"
	}
}

// checkRoutine checks that the routine named by owa exists in the matching
// space (spaceFunction, spaceProcedure, or either for OBJECT_ROUTINE).
// When acceptsAggregate is true (or objType is OBJECT_ROUTINE), absence of the
// routine name in an existing schema is kept conservative because
// pg/driver/sync.go omits aggregates (prokind = 'a') from Functions.
func (b *builder) checkRoutine(owa *ast.ObjectWithArgs, objType ast.ObjectType, missingOK, acceptsAggregate bool) resolved {
	if owa == nil {
		return resolved{}
	}
	return b.checkRoutineName(nameParts(owa.Objname), objType, missingOK, acceptsAggregate)
}

func (b *builder) checkRoutineName(parts []string, objType ast.ObjectType, missingOK, acceptsAggregate bool) resolved {
	if len(parts) == 0 {
		return resolved{}
	}
	r := b.resolveRoutine(parts, objType)
	switch {
	case r.fail != "":
		b.fail(nil, walk.Absent, review.CrossDatabaseReference, r.fail) // 0A000
	case r.fact == walk.Absent && !missingOK:
		if r.missingSchema != "" {
			b.fail(r.keys, walk.Absent, review.ObjectNotFound, schemaNotFound(r.missingSchema)) // 3F000
		} else if !acceptsAggregate {
			b.fail(r.keys, walk.Absent, review.ObjectNotFound, notFound(routineWord(objType), qualified(parts))) // 42883
		}
	}
	return r
}

// dropRoutine follows DROP FUNCTION / PROCEDURE / ROUTINE. Because routines
// may be overloaded, dropping one overload weakens the name rather than
// marking it Absent.
func (b *builder) dropRoutine(obj ast.Node, objType ast.ObjectType, missingOK, cascade bool) {
	owa, ok := obj.(*ast.ObjectWithArgs)
	if !ok {
		return
	}
	parts := nameParts(owa.Objname)
	r := b.checkRoutineName(parts, objType, missingOK, objType == ast.OBJECT_ROUTINE)
	if r.fail != "" || r.fact == walk.Absent {
		return
	}
	if objType == ast.OBJECT_FUNCTION && r.fact == walk.Present && !b.w.isRoutineDirty(r.key) && b.w.pg != nil {
		fnKey := [2]string{schemaOfKey(r.key), r.key.Name}
		if !owa.ArgsUnspecified {
			if fi, ok := b.uniqueSyncedFunc(parts, r); ok {
				if argc, ok := routineArgCount(owa); ok && !fi.arities[argc] {
					if !missingOK {
						b.fail(r.keys, walk.Absent, review.ObjectNotFound, notFound("function", qualified(parts))) // 42883
					}
					return
				}
			}
		}
		if !cascade {
			if fi, ok := b.uniqueSyncedFunc(parts, r); ok && len(fi.arities) == 1 {
				for _, ref := range b.w.pg.trigFuncs[fnKey] {
					tk := relKey(ref.tableSchema, ref.tableName)
					if b.w.isTableDirty(tk) {
						continue
					}
					if tf, _ := b.l.Lookup(tk); tf != walk.Present {
						continue
					}
					if trf, _ := b.l.Lookup(trigKey(tk, ref.trigName)); trf != walk.Present {
						continue
					}
					b.fail([]walk.Key{r.key, trigKey(tk, ref.trigName)}, walk.Present, review.UnsupportedForm, "Function "+review.Code(qualified(parts))+" is used by trigger "+review.Code(ident(ref.trigName))+" and can't be dropped without `CASCADE`.") // 2BP01
					return
				}
			}
		}
	}
	b.weakenRoutine(r, parts, objType)
	b.session(func() {
		if len(parts) >= 2 && r.fact == walk.Present {
			b.w.dirtyRoutines[r.key] = true
		} else {
			b.w.allDirty = true
		}
	})
}

func (b *builder) weakenRoutine(r resolved, parts []string, objType ast.ObjectType) {
	// An unqualified routine name may resolve to an overload in a later
	// search-path schema than the first schema that holds the name, so weaken
	// across all schemas rather than only r.key's schema.
	if len(parts) >= 2 && r.fact == walk.Present {
		b.set(r.key, walk.Unknown, "")
		if objType == ast.OBJECT_ROUTINE {
			s := schemaOfKey(r.key)
			b.set(funcKey(s, r.key.Name), walk.Unknown, "")
			b.set(procKey(s, r.key.Name), walk.Unknown, "")
		}
		return
	}
	prefix := ""
	if len(parts) > 0 {
		if n, ok := walkName(parts[len(parts)-1]); ok {
			prefix = n
		}
	}
	switch objType {
	case ast.OBJECT_FUNCTION:
		b.mayRemove(walk.Scope{Space: spaceFunction, Deep: true, Prefix: prefix})
	case ast.OBJECT_PROCEDURE:
		b.mayRemove(walk.Scope{Space: spaceProcedure, Deep: true, Prefix: prefix})
	default:
		b.mayRemove(walk.Scope{Space: spaceFunction, Deep: true, Prefix: prefix})
		b.mayRemove(walk.Scope{Space: spaceProcedure, Deep: true, Prefix: prefix})
	}
}

func (b *builder) renameRoutine(v *ast.RenameStmt) bool {
	name := v.Newname
	b.session(func() { b.w.routines[name] = true })
	owa, ok := v.Object.(*ast.ObjectWithArgs)
	if !ok {
		return true
	}
	parts := nameParts(owa.Objname)
	r := b.checkRoutineName(parts, v.RenameType, false, v.RenameType == ast.OBJECT_ROUTINE)
	if r.fail != "" || r.fact == walk.Absent {
		return true
	}
	newName, ok := walkName(v.Newname)
	if (v.RenameType == ast.OBJECT_FUNCTION || v.RenameType == ast.OBJECT_ROUTINE) && ok {
		if src, ok := b.uniqueSyncedFunc(parts, r); ok && len(src.arities) == 1 && src.zeroArg {
			s := schemaOfKey(r.key)
			dstKey := funcKey(s, newName)
			dst := b.w.pg.funcSigs[[2]string{s, newName}]
			if dst != nil && dst.complete && dst.zeroArg && !b.w.isRoutineDirty(dstKey) &&
				(owa.ArgsUnspecified || func() bool { n, ok := routineArgCount(owa); return ok && n == 0 }()) {
				b.fail([]walk.Key{r.key, dstKey}, walk.Present, review.NameTaken, alreadyExists("function", ident(v.Newname))) // 42723 / 42725
				return true
			}
		}
	}
	b.weakenRoutine(r, parts, v.RenameType)
	b.session(func() { b.w.allDirty = true })
	if len(parts) >= 2 && r.fact == walk.Present && ok && v.RenameType != ast.OBJECT_ROUTINE {
		s := schemaOfKey(r.key)
		var nk walk.Key
		if r.key.Space == spaceProcedure {
			nk = procKey(s, newName)
		} else {
			nk = funcKey(s, newName)
		}
		b.set(nk, walk.Present, "")
		b.made(nk)
	} else {
		prefix := ""
		if ok {
			prefix = newName
		}
		scope := walk.Scope{Deep: true, Prefix: prefix}
		if len(parts) >= 2 && r.fact == walk.Present {
			scope = walk.Scope{Parent: walk.Path(schemaOfKey(r.key)), Prefix: prefix}
		}
		switch v.RenameType {
		case ast.OBJECT_FUNCTION:
			scope.Space = spaceFunction
			b.mayAdd(scope)
		case ast.OBJECT_PROCEDURE:
			scope.Space = spaceProcedure
			b.mayAdd(scope)
		default:
			scope.Space = spaceFunction
			b.mayAdd(scope)
			scope.Space = spaceProcedure
			b.mayAdd(scope)
		}
	}
	return true
}

func (b *builder) setSchemaRoutine(v *ast.AlterObjectSchemaStmt) bool {
	owa, ok := v.Object.(*ast.ObjectWithArgs)
	if !ok {
		return true
	}
	parts := nameParts(owa.Objname)
	r := b.checkRoutineName(parts, v.ObjectType, false, v.ObjectType == ast.OBJECT_ROUTINE)
	if r.fail != "" || r.fact == walk.Absent {
		return true
	}
	b.session(func() { b.w.allDirty = true })
	ns, ok := walkName(v.Newschema)
	if !ok || ns == tempSchema || ns == "pg_catalog" {
		b.weakenRoutine(r, parts, v.ObjectType)
		return true
	}
	b.need(schemaKey(ns), walk.Present, nil, &walkNote{rng: b.anchor(), absent: schemaNotFound(v.Newschema)}) // 3F000
	b.weakenRoutine(r, parts, v.ObjectType)
	if r.fact == walk.Present && v.ObjectType != ast.OBJECT_ROUTINE {
		var nk walk.Key
		if r.key.Space == spaceProcedure {
			nk = procKey(ns, r.key.Name)
		} else {
			nk = funcKey(ns, r.key.Name)
		}
		b.set(nk, walk.Present, "")
		b.made(nk)
	} else {
		switch v.ObjectType {
		case ast.OBJECT_FUNCTION:
			b.mayAdd(walk.Scope{Space: spaceFunction, Parent: walk.Path(ns)})
		case ast.OBJECT_PROCEDURE:
			b.mayAdd(walk.Scope{Space: spaceProcedure, Parent: walk.Path(ns)})
		default:
			b.mayAdd(walk.Scope{Space: spaceFunction, Parent: walk.Path(ns)})
			b.mayAdd(walk.Scope{Space: spaceProcedure, Parent: walk.Path(ns)})
		}
	}
	return true
}

// alterOwner follows ALTER ... OWNER TO for objects whose existence or kind
// the ledger knows (routines, domains, types, schemas, relations).
func (b *builder) alterOwner(v *ast.AlterOwnerStmt) bool {
	switch v.ObjectType {
	case ast.OBJECT_FUNCTION, ast.OBJECT_PROCEDURE, ast.OBJECT_ROUTINE:
		if owa, ok := v.Object.(*ast.ObjectWithArgs); ok {
			b.checkRoutine(owa, v.ObjectType, false, v.ObjectType == ast.OBJECT_ROUTINE)
		}
	case ast.OBJECT_DOMAIN:
		var parts []string
		if tn, ok := v.Object.(*ast.TypeName); ok {
			parts = nameParts(tn.Names)
		} else {
			parts = nameParts(listOf(v.Object))
		}
		if len(parts) > 0 {
			b.checkDomain(parts, b.anchor())
		}
	case ast.OBJECT_TYPE:
		var parts []string
		if tn, ok := v.Object.(*ast.TypeName); ok {
			parts = nameParts(tn.Names)
		} else {
			parts = nameParts(listOf(v.Object))
		}
		if len(parts) > 0 {
			r := b.resolveType(parts)
			switch {
			case r.fail != "":
				b.fail(nil, walk.Absent, review.CrossDatabaseReference, r.fail)
			case r.fact == walk.Absent:
				if r.missingSchema != "" {
					b.fail(r.keys, walk.Absent, review.ObjectNotFound, schemaNotFound(r.missingSchema))
				} else {
					b.fail(r.keys, walk.Absent, review.ObjectNotFound, notFound("type", qualified(parts)))
				}
			}
		}
	case ast.OBJECT_SCHEMA:
		if s, ok := v.Object.(*ast.String); ok {
			if name, ok := walkName(s.Str); ok {
				b.need(schemaKey(name), walk.Present, nil, &walkNote{rng: b.anchor(), absent: schemaNotFound(s.Str)}) // 3F000
			}
		}
	case ast.OBJECT_TABLE, ast.OBJECT_VIEW, ast.OBJECT_MATVIEW, ast.OBJECT_SEQUENCE, ast.OBJECT_FOREIGN_TABLE, ast.OBJECT_INDEX:
		if v.Relation != nil {
			r := b.exists(v.Relation, false)
			b.wrongKind(r, v.Relation, v.ObjectType, formAlter)
		}
	}
	return true
}

// renameAttribute follows ALTER TYPE t RENAME ATTRIBUTE old TO new.
func (b *builder) renameAttribute(v *ast.RenameStmt) bool {
	r := b.exists(v.Relation, v.MissingOk)
	switch {
	case r.fail != "", r.fact == walk.Absent:
		return true
	case r.fact == walk.Unknown:
		return false
	}
	if r.kind == relSequence {
		b.need(r.key, walk.Present, []string{relComposite, relTable, relView, relMatview, relForeign}, ofKind(rangeOf(v.Relation.Loc), relation(v.Relation), relComposite)) // 42809
		return true
	}
	rng := b.anchor()
	oldName, ok1 := walkName(v.Subname)
	newName, ok2 := walkName(v.Newname)
	if !ok1 || !ok2 {
		b.forget(walk.Scope{Space: spaceColumn, Parent: tableOf(r.key)})
		return true
	}
	if hasColumns(r.kind) {
		oldKey := columnKey(r.key, oldName)
		newKey := columnKey(r.key, newName)
		b.need(oldKey, walk.Present, nil, &walkNote{rng: rng, absent: notFoundIn("column", ident(v.Subname), relation(v.Relation))})      // 42703
		b.need(newKey, walk.Absent, nil, &walkNote{rng: rng, present: alreadyExistsIn("column", ident(v.Newname), relation(v.Relation))}) // 42701
		b.set(oldKey, walk.Absent, "")
		b.set(newKey, walk.Present, kindColumn)
		return true
	}
	if oldName == newName && r.kind != "" {
		b.failAt(rng, nil, review.DuplicateName, alreadyExistsIn("column", ident(v.Newname), relation(v.Relation))) // 42701 / 42703
		return true
	}
	b.forget(walk.Scope{Space: spaceColumn, Parent: tableOf(r.key)})
	b.session(func() { b.w.dirtyTypes[r.key] = true })
	return true
}

// checkStorageMode validates ALTER COLUMN ... SET STORAGE mode (22023).
func (b *builder) checkStorageMode(mode string, rng review.Range) {
	switch strings.ToLower(mode) {
	case "plain", "external", "extended", "main", "default":
	default:
		b.failAt(rng, nil, review.UnsupportedForm, unsupported("storage mode "+review.Code(ident(mode)))) // 22023
	}
}

// checkCompressionMethod validates ALTER COLUMN ... SET COMPRESSION method
// against PostgreSQL's hardcoded toast compression methods (22023 / 0A000).
//
// pg: src/backend/access/common/toast_compression.c — GetAttributeCompression
func (b *builder) checkCompressionMethod(method string, rng review.Range) {
	switch strings.ToLower(method) {
	case "", "default", "pglz", "lz4":
	default:
		b.failAt(rng, nil, review.UnsupportedForm, unsupported("compression method "+review.Code(ident(method)))) // 22023 / 0A000
	}
}

// checkAttrOptions validates ALTER COLUMN ... SET (...) against PostgreSQL's
// fixed attribute_reloptions table (22023).
//
// pg: src/backend/access/common/reloptions.c — attribute_reloptions
func (b *builder) checkAttrOptions(opts *ast.List, rng review.Range) {
	if opts == nil {
		return
	}
	for _, item := range opts.Items {
		de, ok := item.(*ast.DefElem)
		if !ok {
			continue
		}
		if de.Defnamespace != "" {
			b.failAt(rng, nil, review.UnsupportedForm, unsupported("attribute option "+review.Code(ident(de.Defnamespace)+"."+ident(de.Defname)))) // 22023
			return
		}
		switch strings.ToLower(de.Defname) {
		case "n_distinct", "n_distinct_inherited":
		default:
			b.failAt(rng, nil, review.UnsupportedForm, unsupported("attribute option "+review.Code(ident(de.Defname)))) // 22023
			return
		}
	}
}

// checkIdentityOptions validates sequence options inside ALTER COLUMN ...
// ADD/SET GENERATED AS IDENTITY / RESTART (42601 / 22P02 / 55000).
func (b *builder) checkIdentityOptions(opts *ast.List, _ review.Range) {
	b.checkSequenceOptions(opts, true)
}

// hasNonToastRelOption reports whether opts contains at least one option that
// is not in the "toast" namespace. PostgreSQL filters toast.* options out of
// the parent relation's reloptions before checking partitioned tables, so
// toast-only SET (...) on a partitioned table succeeds as a no-op.
//
// pg: src/backend/access/common/reloptions.c — transformRelOptions;
// src/backend/commands/tablecmds.c — ATExecSetRelOptions
func hasNonToastRelOption(opts *ast.List) bool {
	if opts == nil {
		return false
	}
	for _, item := range opts.Items {
		if de, ok := item.(*ast.DefElem); ok && !strings.EqualFold(de.Defnamespace, "toast") {
			return true
		}
	}
	return false
}

// checkRelOptions validates duplicate reloptions (rejected for all relation
// kinds and access methods, 22023) and view reloptions (not extensible in
// PostgreSQL, 22023).
//
// pg: src/backend/access/common/reloptions.c — parseRelOptionsInternal,
// view_reloptions
func (b *builder) checkRelOptions(relKind string, opts *ast.List, rng review.Range) {
	if opts == nil {
		return
	}
	seen := make(map[[2]string]bool, len(opts.Items))
	for _, item := range opts.Items {
		de, ok := item.(*ast.DefElem)
		if !ok {
			continue
		}
		key := [2]string{strings.ToLower(de.Defnamespace), strings.ToLower(de.Defname)}
		if seen[key] {
			b.failAt(rng, nil, review.UnsupportedForm, unsupported("duplicate storage parameter "+review.Code(ident(de.Defname)))) // 22023
			return
		}
		seen[key] = true
	}
	if relKind == relView {
		for _, item := range opts.Items {
			de, ok := item.(*ast.DefElem)
			if !ok {
				continue
			}
			if de.Defnamespace != "" {
				b.failAt(rng, nil, review.UnsupportedForm, unsupported("view option "+review.Code(ident(de.Defnamespace)+"."+ident(de.Defname)))) // 22023
				return
			}
			switch strings.ToLower(de.Defname) {
			case "check_option", "security_barrier", "security_invoker":
			default:
				b.failAt(rng, nil, review.UnsupportedForm, unsupported("view option "+review.Code(ident(de.Defname)))) // 22023
				return
			}
		}
	}
}

// heapRelOptions lists the storage parameters of PostgreSQL's heap, which
// tables and materialized views take, in any release from 12 to 18.
//
// pg: src/backend/access/common/reloptions.c — boolRelOpts, intRelOpts,
// realRelOpts, enumRelOpts (RELOPT_KIND_HEAP); "oids" is accepted as false
// by CREATE TABLE
var heapRelOptions = map[string]bool{
	"fillfactor":                            true,
	"autovacuum_enabled":                    true,
	"user_catalog_table":                    true,
	"vacuum_truncate":                       true,
	"vacuum_index_cleanup":                  true,
	"parallel_workers":                      true,
	"toast_tuple_target":                    true,
	"log_autovacuum_min_duration":           true,
	"oids":                                  true,
	"autovacuum_vacuum_threshold":           true,
	"autovacuum_vacuum_max_threshold":       true,
	"autovacuum_vacuum_insert_threshold":    true,
	"autovacuum_analyze_threshold":          true,
	"autovacuum_vacuum_cost_limit":          true,
	"autovacuum_vacuum_cost_delay":          true,
	"autovacuum_freeze_min_age":             true,
	"autovacuum_freeze_max_age":             true,
	"autovacuum_freeze_table_age":           true,
	"autovacuum_multixact_freeze_min_age":   true,
	"autovacuum_multixact_freeze_max_age":   true,
	"autovacuum_multixact_freeze_table_age": true,
	"autovacuum_vacuum_scale_factor":        true,
	"autovacuum_vacuum_insert_scale_factor": true,
	"autovacuum_analyze_scale_factor":       true,
	"vacuum_max_eager_freeze_failure_rate":  true,
}

// checkHeapRelOptions refuses a storage parameter a table or materialized
// view does not have, or a namespace other than toast (22023). An
// extension may register parameters or take a namespace's (TimescaleDB's
// timescaledb., Citus's columnar.), and may bring a table access method of
// its own, so the check holds only while the synced schema lists none but
// plpgsql; CREATE EXTENSION ends the walk.
func (b *builder) checkHeapRelOptions(opts *ast.List, rng review.Range) {
	if opts == nil || b.w.pg == nil || b.w.pg.extensions {
		return
	}
	for _, item := range opts.Items {
		de, ok := item.(*ast.DefElem)
		if !ok {
			continue
		}
		switch {
		case de.Defnamespace != "":
			if !strings.EqualFold(de.Defnamespace, "toast") {
				b.failAt(rng, nil, review.UnsupportedForm, unsupported("storage parameter namespace "+review.Code(ident(de.Defnamespace)))) // 22023
				return
			}
		case !heapRelOptions[strings.ToLower(de.Defname)]:
			b.failAt(rng, nil, review.UnsupportedForm, unsupported("storage parameter "+review.Code(ident(de.Defname)))) // 22023
			return
		}
	}
}

// selectHasUnboundColumn reports whether q references a column that no
// FROM item can supply, outside any subquery (42703). A query with no
// outer query (a view's, or CREATE TABLE AS's) resolves a column only
// against its own FROM items:
//   - in a FROM-less SELECT or VALUES, every column reference of the
//     target list, WHERE, HAVING, LIMIT, OFFSET, and WINDOW clause, and of
//     GROUP BY, ORDER BY, and DISTINCT ON except a bare name that may be an
//     output column's (pg: parse_clause.c — findTargetlistEntrySQL92);
//   - in each arm of a set operation, as in a query of its own;
//   - in the first FROM item's function arguments or table function
//     expressions, which no earlier item precedes to be LATERAL to.
func selectHasUnboundColumn(q ast.Node) bool {
	sel, ok := q.(*ast.SelectStmt)
	if !ok || sel == nil {
		return false
	}
	if sel.Op != ast.SETOP_NONE {
		// ORDER BY of a set operation sees its output columns, so only
		// the arms are judged.
		return selectHasUnboundColumn(sel.Larg) || selectHasUnboundColumn(sel.Rarg)
	}
	if sel.FromClause != nil && len(sel.FromClause.Items) > 0 {
		return fromItemHasColumnRef(sel.FromClause.Items[0])
	}
	if sel.ValuesLists != nil {
		return exprHasColumnRef(sel.ValuesLists)
	}
	if exprHasColumnRef(sel.TargetList) || exprHasColumnRef(sel.WhereClause) || exprHasColumnRef(sel.HavingClause) ||
		exprHasColumnRef(sel.LimitCount) || exprHasColumnRef(sel.LimitOffset) {
		return true
	}
	if sel.WindowClause != nil {
		// PARTITION BY and ORDER BY of a window follow SQL99 rules and
		// see no output column (pg: parse_clause.c — transformWindowDefinitions).
		for _, item := range sel.WindowClause.Items {
			if wd, ok := item.(*ast.WindowDef); ok && windowHasColumnRef(wd) {
				return true
			}
		}
	}
	names, known := outputNames(sel)
	if sel.SortClause != nil {
		for _, item := range sel.SortClause.Items {
			if sb, ok := item.(*ast.SortBy); ok && groupItemHasColumnRef(sb.Node, names, known, false) {
				return true
			}
		}
	}
	if sel.DistinctClause != nil {
		for _, item := range sel.DistinctClause.Items {
			if groupItemHasColumnRef(item, names, known, false) {
				return true
			}
		}
	}
	if sel.GroupClause != nil {
		for _, item := range sel.GroupClause.Items {
			if groupItemHasColumnRef(item, names, known, true) {
				return true
			}
		}
	}
	return false
}

// outputNames lists the names a SELECT's output columns take, and whether
// the list is complete: a column with no alias is named by
// FigureColname, which is followed here only for a constant and a function
// call.
func outputNames(sel *ast.SelectStmt) (map[string]bool, bool) {
	names := make(map[string]bool)
	if sel.TargetList == nil {
		return names, true
	}
	known := true
	for _, item := range sel.TargetList.Items {
		rt, ok := item.(*ast.ResTarget)
		if !ok {
			known = false
			continue
		}
		if rt.Name != "" {
			names[rt.Name] = true
			continue
		}
		switch v := rt.Val.(type) {
		case *ast.A_Const:
			names["?column?"] = true
		case *ast.FuncCall:
			if parts := nameParts(v.Funcname); len(parts) > 0 {
				names[parts[len(parts)-1]] = true
			} else {
				known = false
			}
		default:
			known = false
		}
	}
	return names, known
}

// groupItemHasColumnRef judges one ORDER BY, DISTINCT ON, or GROUP BY item
// of a FROM-less SELECT. A bare name is first matched against the output
// columns; in GROUP BY, so is each member of a grouping set or an implicit
// row (pg: parse_agg.c — flatten_grouping_sets).
func groupItemHasColumnRef(n ast.Node, names map[string]bool, known, group bool) bool {
	switch v := n.(type) {
	case nil:
		return false
	case *ast.ColumnRef:
		if v.Fields != nil && len(v.Fields.Items) == 1 {
			if s, ok := v.Fields.Items[0].(*ast.String); ok {
				return known && !names[s.Str]
			}
		}
		return true
	case *ast.GroupingSet:
		if group && v.Content != nil {
			for _, item := range v.Content.Items {
				if groupItemHasColumnRef(item, names, known, group) {
					return true
				}
			}
			return false
		}
	case *ast.RowExpr:
		if group {
			if v.Args != nil {
				for _, item := range v.Args.Items {
					if groupItemHasColumnRef(item, names, known, group) {
						return true
					}
				}
			}
			return false
		}
	}
	return exprHasColumnRef(n)
}

func windowHasColumnRef(wd *ast.WindowDef) bool {
	if exprHasColumnRef(wd.PartitionClause) || exprHasColumnRef(wd.StartOffset) || exprHasColumnRef(wd.EndOffset) {
		return true
	}
	if wd.OrderClause != nil {
		for _, item := range wd.OrderClause.Items {
			if sb, ok := item.(*ast.SortBy); ok && exprHasColumnRef(sb.Node) {
				return true
			}
		}
	}
	return false
}

// fromItemHasColumnRef reports whether the first FROM item is a function
// or table function whose expressions reference a column: nothing precedes
// it to supply one.
func fromItemHasColumnRef(n ast.Node) bool {
	switch v := n.(type) {
	case *ast.RangeFunction:
		return exprHasColumnRef(v.Functions)
	case *ast.RangeTableFunc:
		if exprHasColumnRef(v.Docexpr) || exprHasColumnRef(v.Rowexpr) || exprHasColumnRef(v.Namespaces) {
			return true
		}
		if v.Columns != nil {
			for _, item := range v.Columns.Items {
				if col, ok := item.(*ast.RangeTableFuncCol); ok && (exprHasColumnRef(col.Colexpr) || exprHasColumnRef(col.Coldefexpr)) {
					return true
				}
			}
		}
	}
	return false
}

// exprHasColumnRef reports whether an expression references a column
// outside any subquery. A type name's modifiers are skipped: they may be
// identifiers (geometry(Point, 4326)) rather than columns.
func exprHasColumnRef(n ast.Node) bool {
	if n == nil {
		return false
	}
	found := false
	ast.Inspect(n, func(m ast.Node) bool {
		if found {
			return false
		}
		switch v := m.(type) {
		case *ast.ColumnRef:
			found = true
			return false
		case *ast.SubLink:
			found = exprHasColumnRef(v.Testexpr)
			return false
		case *ast.SelectStmt, *ast.RangeSubselect, *ast.TypeName, *ast.WithClause, *ast.ColumnDef:
			return false
		}
		return true
	})
	return found
}

// hasBareTriggerColumn reports whether a trigger WHEN expression directly
// references an unqualified column other than OLD or NEW (outside any
// subquery), which PostgreSQL always rejects with 42702 or 42703.
func hasBareTriggerColumn(n ast.Node) bool {
	if n == nil {
		return false
	}
	switch v := n.(type) {
	case *ast.ColumnRef:
		if v.Fields != nil && len(v.Fields.Items) == 1 {
			if s, ok := v.Fields.Items[0].(*ast.String); ok {
				lo := strings.ToLower(s.Str)
				return lo != "old" && lo != "new"
			}
		}
		return false
	case *ast.A_Expr:
		return hasBareTriggerColumn(v.Lexpr) || hasBareTriggerColumn(v.Rexpr)
	case *ast.BoolExpr:
		if v.Args != nil {
			for _, item := range v.Args.Items {
				if hasBareTriggerColumn(item) {
					return true
				}
			}
		}
	case *ast.NullTest:
		return hasBareTriggerColumn(v.Arg)
	case *ast.BooleanTest:
		return hasBareTriggerColumn(v.Arg)
	case *ast.TypeCast:
		return hasBareTriggerColumn(v.Arg)
	case *ast.CoalesceExpr:
		if v.Args != nil {
			for _, item := range v.Args.Items {
				if hasBareTriggerColumn(item) {
					return true
				}
			}
		}
	case *ast.CaseExpr:
		if hasBareTriggerColumn(v.Arg) || hasBareTriggerColumn(v.Defresult) {
			return true
		}
		if v.Args != nil {
			for _, item := range v.Args.Items {
				if cw, ok := item.(*ast.CaseWhen); ok {
					if hasBareTriggerColumn(cw.Expr) || hasBareTriggerColumn(cw.Result) {
						return true
					}
				}
			}
		}
	case *ast.FuncCall:
		if v.Args != nil {
			for _, item := range v.Args.Items {
				if hasBareTriggerColumn(item) {
					return true
				}
			}
		}
	}
	return false
}

// createFunction follows CREATE [OR REPLACE] FUNCTION / PROCEDURE: validates
// unconditional PostgreSQL definition rules (schema, options, body, language,
// return type, parameter modes) and records the routine name.
//
// pg: src/backend/commands/functioncmds.c — CreateFunction
func (b *builder) createFunction(v *ast.CreateFunctionStmt) bool {
	parts := nameParts(v.Funcname)
	if len(parts) > 0 {
		name := parts[len(parts)-1]
		b.session(func() { b.w.routines[name] = true })
	}
	if len(parts) == 3 {
		if b.w.crossDatabase(parts[0]) {
			b.fail(nil, walk.Absent, review.CrossDatabaseReference, crossDatabaseMessage(qualified(parts))) // 0A000
			return true
		}
		if b.w.currentDatabase(parts[0]) {
			parts = parts[1:]
		} else {
			parts = nil
		}
	}
	var createdSchema, createdName string
	var hasCreatedKey bool
	if len(parts) >= 1 && len(parts) <= 2 {
		schema, written, ok := splitName(parts)
		if ok {
			c := b.creationSchema(schema, written, false)
			if c.fail != "" {
				b.fail(c.keys, walk.Absent, c.code, c.fail) // 3F000
				return true
			}
			if n, ok := walkName(written); c.fact == walk.Present && ok {
				createdSchema, createdName, hasCreatedKey = c.schema, n, true
			}
		}
	}

	isProcedure := false
	if v.Options != nil {
		for _, item := range v.Options.Items {
			if de, ok := item.(*ast.DefElem); ok && strings.EqualFold(de.Defname, "isProcedure") {
				switch arg := de.Arg.(type) {
				case nil:
					isProcedure = true
				case *ast.Boolean:
					isProcedure = arg.Boolval
				case *ast.Integer:
					isProcedure = arg.Ival != 0
				}
			}
		}
	}

	var hasAs, hasLang bool
	var seenVolatility, seenStrict, seenSecurity, seenParallel, seenLeakproof, seenWindow, seenCost, seenRows, seenSupport, seenLang bool
	if v.Options != nil {
		for _, item := range v.Options.Items {
			de, ok := item.(*ast.DefElem)
			if !ok {
				continue
			}
			switch strings.ToLower(de.Defname) {
			case "language":
				if seenLang {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("conflicting or redundant options in `CREATE FUNCTION`")) // 42601
					return true
				}
				seenLang = true
				hasLang = true
			case "as":
				if hasAs {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("conflicting or redundant options in `CREATE FUNCTION`")) // 42601
					return true
				}
				hasAs = true
			case "volatility":
				if isProcedure {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("volatility attribute in `CREATE PROCEDURE`")) // 42P13
					return true
				}
				if seenVolatility {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("conflicting or redundant options in `CREATE FUNCTION`")) // 42601
					return true
				}
				seenVolatility = true
			case "parallel":
				if isProcedure {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("`PARALLEL` in `CREATE PROCEDURE`")) // 42P13
					return true
				}
				if seenParallel {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("conflicting or redundant options in `CREATE FUNCTION`")) // 42601
					return true
				}
				seenParallel = true
				if s, ok := de.Arg.(*ast.String); ok {
					switch strings.ToLower(s.Str) {
					case "safe", "restricted", "unsafe":
					default:
						b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("`PARALLEL "+s.Str+"` in `CREATE FUNCTION`")) // 42601
						return true
					}
				}
			case "strict", "called_on_null_input":
				if isProcedure {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("strictness attribute in `CREATE PROCEDURE`")) // 42P13
					return true
				}
				if seenStrict {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("conflicting or redundant options in `CREATE FUNCTION`")) // 42601
					return true
				}
				seenStrict = true
			case "security":
				if seenSecurity {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("conflicting or redundant options in `CREATE FUNCTION`")) // 42601
					return true
				}
				seenSecurity = true
			case "leakproof":
				if seenLeakproof {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("conflicting or redundant options in `CREATE FUNCTION`")) // 42601
					return true
				}
				seenLeakproof = true
			case "window":
				if isProcedure {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("`WINDOW` in `CREATE PROCEDURE`")) // 42P13
					return true
				}
				if seenWindow {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("conflicting or redundant options in `CREATE FUNCTION`")) // 42601
					return true
				}
				seenWindow = true
			case "cost":
				if seenCost {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("conflicting or redundant options in `CREATE FUNCTION`")) // 42601
					return true
				}
				seenCost = true
			case "rows":
				if seenRows {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("conflicting or redundant options in `CREATE FUNCTION`")) // 42601
					return true
				}
				seenRows = true
			case "support":
				if isProcedure {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("`SUPPORT` in `CREATE PROCEDURE`")) // 42P13
					return true
				}
				if seenSupport {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("conflicting or redundant options in `CREATE FUNCTION`")) // 42601
					return true
				}
				seenSupport = true
			}
		}
	}

	hasSqlBody := v.SqlBody != nil
	if !hasAs && !hasSqlBody {
		b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("`CREATE FUNCTION` without a function body")) // 42P13
		return true
	}
	if !hasLang && !hasSqlBody {
		b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("`CREATE FUNCTION` without a `LANGUAGE` clause")) // 42P13
		return true
	}

	hasOutParam := false
	seenVariadic := false
	if v.Parameters != nil {
		for _, item := range v.Parameters.Items {
			fp, ok := item.(*ast.FunctionParameter)
			if !ok {
				continue
			}
			if fp.ArgType != nil && fp.ArgType.Setof {
				b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("`SETOF` parameter in `CREATE FUNCTION`")) // 42P13
				return true
			}
			mode := fp.Mode
			if mode == ast.FUNC_PARAM_DEFAULT {
				mode = ast.FUNC_PARAM_IN
			}
			if mode != ast.FUNC_PARAM_OUT && mode != ast.FUNC_PARAM_TABLE {
				if seenVariadic {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("parameter after `VARIADIC` in `CREATE FUNCTION`")) // 42P13
					return true
				}
			}
			if mode == ast.FUNC_PARAM_VARIADIC {
				seenVariadic = true
			}
			if mode == ast.FUNC_PARAM_OUT || mode == ast.FUNC_PARAM_INOUT || mode == ast.FUNC_PARAM_TABLE {
				hasOutParam = true
			}
		}
	}
	if !isProcedure && v.ReturnType == nil && !hasOutParam {
		b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("`CREATE FUNCTION` without a `RETURNS` clause")) // 42P13
		return true
	}
	if hasCreatedKey {
		var k walk.Key
		if isProcedure {
			k = procKey(createdSchema, createdName)
		} else {
			k = funcKey(createdSchema, createdName)
		}
		b.set(k, walk.Present, "")
		b.made(k)
		b.session(func() { b.w.dirtyRoutines[k] = true })
	} else {
		if isProcedure {
			b.mayAdd(walk.Scope{Space: spaceProcedure, Deep: true})
		} else {
			b.mayAdd(walk.Scope{Space: spaceFunction, Deep: true})
		}
		b.session(func() { b.w.allDirty = true })
	}
	return true
}

// alterFunction follows ALTER FUNCTION / PROCEDURE / ROUTINE: validates
// unconditional option checks and routine existence.
func (b *builder) alterFunction(v *ast.AlterFunctionStmt) bool {
	var seenVolatility, seenStrict, seenSecurity, seenParallel, seenLeakproof, seenCost, seenRows, seenSupport bool
	hasNonAggAction := false
	if v.Actions != nil {
		for _, item := range v.Actions.Items {
			de, ok := item.(*ast.DefElem)
			if !ok {
				continue
			}
			defName := strings.ToLower(de.Defname)
			if defName != "parallel" {
				hasNonAggAction = true
			}
			switch defName {
			case "volatility":
				if seenVolatility {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("conflicting or redundant options in `ALTER FUNCTION`")) // 42601
					return true
				}
				seenVolatility = true
			case "strict", "called_on_null_input":
				if seenStrict {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("conflicting or redundant options in `ALTER FUNCTION`")) // 42601
					return true
				}
				seenStrict = true
			case "security":
				if seenSecurity {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("conflicting or redundant options in `ALTER FUNCTION`")) // 42601
					return true
				}
				seenSecurity = true
			case "parallel":
				if seenParallel {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("conflicting or redundant options in `ALTER FUNCTION`")) // 42601
					return true
				}
				seenParallel = true
				if s, ok := de.Arg.(*ast.String); ok {
					switch strings.ToLower(s.Str) {
					case "safe", "restricted", "unsafe":
					default:
						b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("`PARALLEL "+s.Str+"` in `ALTER FUNCTION`")) // 42601
						return true
					}
				}
			case "leakproof":
				if seenLeakproof {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("conflicting or redundant options in `ALTER FUNCTION`")) // 42601
					return true
				}
				seenLeakproof = true
			case "cost":
				if seenCost {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("conflicting or redundant options in `ALTER FUNCTION`")) // 42601
					return true
				}
				seenCost = true
			case "rows":
				if seenRows {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("conflicting or redundant options in `ALTER FUNCTION`")) // 42601
					return true
				}
				seenRows = true
			case "support":
				if seenSupport {
					b.fail(nil, walk.Present, review.UnsupportedForm, unsupported("conflicting or redundant options in `ALTER FUNCTION`")) // 42601
					return true
				}
				seenSupport = true
			}
		}
	}
	if v.Func != nil {
		// PostgreSQL's AlterFunction rejects aggregates (prokind = 'a') with
		// 42809 for function-attribute actions (STRICT, VOLATILE, COST, ROWS,
		// etc.), even under ALTER ROUTINE. Keep PARALLEL-only ALTER ROUTINE
		// conservative for omitted aggregates.
		acceptsAgg := v.Objtype == ast.OBJECT_ROUTINE && !hasNonAggAction
		r := b.checkRoutine(v.Func, v.Objtype, false, acceptsAgg)
		if seenRows {
			if fi, ok := b.uniqueSyncedFunc(nameParts(v.Func.Objname), r); ok && !fi.anySetof {
				b.fail([]walk.Key{r.key}, walk.Present, review.UnsupportedForm, unsupported("`ROWS` on a function that doesn't return a set")) // 22023
				return true
			}
		}
	}
	return true
}

// uniqueSyncedFunc returns the complete synced function signature info for
// parts/r only when r is definitely the sole schema that can satisfy parts
// (for unqualified names, no other search-path schema or pg_catalog holds the
// name) and the routine has not been modified earlier in the change.
func (b *builder) uniqueSyncedFunc(parts []string, r resolved) (*syncedFuncInfo, bool) {
	if r.fact != walk.Present || r.key.Space != spaceFunction || b.w.isRoutineDirty(r.key) || b.w.pg == nil || b.w.pg.lostRoutines {
		return nil, false
	}
	targetSchema := schemaOfKey(r.key)
	if len(parts) == 1 {
		if !trustedRoutineName(r.key.Name) || builtinCatalogFunctions[r.key.Name] {
			return nil, false
		}
		sp, ok := b.w.searchPath()
		if !ok {
			return nil, false
		}
		for _, s := range sp {
			if s == targetSchema || s == "pg_catalog" {
				continue
			}
			if sf, _ := b.l.Lookup(schemaKey(s)); sf == walk.Absent {
				continue
			}
			if f, _ := b.l.Lookup(funcKey(s, r.key.Name)); f != walk.Absent {
				return nil, false
			}
		}
	}
	fi := b.w.pg.funcSigs[[2]string{targetSchema, r.key.Name}]
	if fi == nil || !fi.complete {
		return nil, false
	}
	return fi, true
}

// bareColumnRef returns the folded column name when expr is an unqualified
// single-part ColumnRef, or ("", false) otherwise.
func bareColumnRef(expr ast.Node) (string, bool) {
	cr, ok := expr.(*ast.ColumnRef)
	if !ok || cr.Fields == nil || len(cr.Fields.Items) != 1 {
		return "", false
	}
	s, ok := cr.Fields.Items[0].(*ast.String)
	if !ok {
		return "", false
	}
	return walkName(s.Str)
}

// isPinnedNonBoolType reports whether typ is PostgreSQL's built-in integer
// (int4) type, whose pg_cast entry (int4 -> bool) is a pinned explicit-only
// system cast that cannot be dropped or replaced with an implicit cast, and
// domains over integer ignore custom casts.
func isPinnedNonBoolType(typ string) bool {
	return typ == "integer" || typ == "int4"
}

// isKnownNonIntegerBuiltinType reports whether typ is a standard PostgreSQL
// built-in data type that is definitely not smallint/integer/bigint (so it
// cannot be made an identity column, 22023).
func isKnownNonIntegerBuiltinType(typ string) bool {
	if typ == "" {
		return false
	}
	base := typ
	if i := strings.IndexByte(base, '('); i >= 0 {
		base = strings.TrimSpace(base[:i])
	}
	switch base {
	case "boolean", "bool", "text", "character varying", "varchar", "character", "char", "bpchar",
		"numeric", "decimal", "real", "float4", "double precision", "float8",
		"date", "uuid", "bytea", "json", "jsonb", "xml", "money", "interval",
		"inet", "cidr", "macaddr", "macaddr8", "bit", "bit varying", "varbit",
		"time", "time without time zone", "time with time zone", "timetz",
		"timestamp", "timestamp without time zone", "timestamp with time zone", "timestamptz":
		return true
	}
	return false
}

func (b *builder) checkPolicyBoolExpr(t walk.Key, expr ast.Node) {
	if col, ok := bareColumnRef(expr); ok {
		if sc, ok := b.w.syncedCol(t, col); ok && isPinnedNonBoolType(sc.typ) {
			b.fail([]walk.Key{t}, walk.Present, review.UnsupportedForm, unsupported("non-boolean expression in policy")) // 42804
		}
	}
}

func (b *builder) activeFKDependent(target walk.Key, alsoRemoved map[walk.Key]bool, inhTargets []walk.Key) (syncedFKRef, bool) {
	if b.w.allDirty || b.w.pg == nil {
		return syncedFKRef{}, false
	}
	for _, ref := range b.w.pg.fkRefs[[2]string{schemaOfKey(target), target.Name}] {
		fromKey := relKey(ref.fromSchema, ref.fromTable)
		if fromKey == target || alsoRemoved[fromKey] || b.w.isTableDirty(fromKey) {
			continue
		}
		if f, _ := b.l.Lookup(fromKey); f != walk.Present {
			continue
		}
		if cf, _ := b.l.Lookup(constraintKey(fromKey, ref.conName)); cf != walk.Present {
			continue
		}
		// TRUNCATE without ONLY also truncates inheritance children of every
		// non-ONLY target in the statement. A child that INHERITS (parent)
		// must have every column of parent; skip fromTable if it could inherit
		// from any non-ONLY target (including any parent whose column list is
		// empty/unknown or a subset of fromTable's columns).
		if len(inhTargets) > 0 {
			fromCols := b.w.pg.columns[[2]string{ref.fromSchema, ref.fromTable}]
			mayInherit := false
			for _, parent := range inhTargets {
				parentCols := b.w.pg.columns[[2]string{schemaOfKey(parent), parent.Name}]
				if len(parentCols) == 0 || fromCols == nil {
					mayInherit = true
					break
				}
				subset := true
				for col := range parentCols {
					if !fromCols[col] {
						subset = false
						break
					}
				}
				if subset {
					mayInherit = true
					break
				}
			}
			if mayInherit {
				continue
			}
		}
		return ref, true
	}
	return syncedFKRef{}, false
}

// checkBtreeRelOptions validates storage parameters on PostgreSQL's built-in
// btree index access method, whose reloptions are fixed (no namespaces, only
// fillfactor, deduplicate_items, and PG 11-13 vacuum_cleanup_index_scale_factor).
//
// pg: src/backend/access/nbtree/nbtree.c — btoptions
func (b *builder) checkBtreeRelOptions(opts *ast.List, rng review.Range) {
	for _, item := range opts.Items {
		de, ok := item.(*ast.DefElem)
		if !ok {
			continue
		}
		if de.Defnamespace != "" {
			b.failAt(rng, nil, review.UnsupportedForm, unsupported("btree index option "+review.Code(ident(de.Defnamespace)+"."+ident(de.Defname)))) // 22023
			return
		}
		switch strings.ToLower(de.Defname) {
		case "fillfactor", "deduplicate_items", "vacuum_cleanup_index_scale_factor":
		default:
			b.failAt(rng, nil, review.UnsupportedForm, unsupported("btree index option "+review.Code(ident(de.Defname)))) // 22023
			return
		}
	}
}

// routineArgCount returns the number of input arguments specified in owa when
// ArgsUnspecified is false.
func routineArgCount(owa *ast.ObjectWithArgs) (int, bool) {
	if owa == nil || owa.ArgsUnspecified {
		return 0, false
	}
	if owa.Objargs == nil {
		return 0, true
	}
	n := 0
	for _, item := range owa.Objargs.Items {
		if fp, ok := item.(*ast.FunctionParameter); ok {
			if fp.Mode == ast.FUNC_PARAM_OUT || fp.Mode == ast.FUNC_PARAM_TABLE {
				continue
			}
		}
		n++
	}
	return n, true
}

// queryRefused fails a query PostgreSQL's analysis refuses whatever the
// data: a name two FROM items of one level both take (42712), or a call
// with no arguments to a function no schema on the path holds (42883).
func (b *builder) queryRefused(q ast.Node) bool {
	if name, ok := duplicateFromName(q); ok {
		b.fail(nil, walk.Present, review.UnsupportedForm, "Table name "+review.Code(ident(name))+" is specified more than once in `FROM`.") // 42712
		return true
	}
	if msg := refusedForm(q); msg != "" {
		b.fail(nil, walk.Present, review.UnsupportedForm, msg)
		return true
	}
	var call []string
	var r resolved
	ast.Inspect(q, func(n ast.Node) bool {
		if call != nil {
			return false
		}
		switch v := n.(type) {
		case *ast.TypeName:
			return false
		case *ast.FuncCall:
			// A call with no arguments is neither a type cast nor field
			// selection, and reaches no aggregate: a parameterless one is
			// called as f(*) and an ordered-set one with WITHIN GROUP
			// (pg: parse_func.c — ParseFuncOrColumn). The sync lists every
			// other function, so one it proves absent is absent.
			if (v.Args == nil || len(v.Args.Items) == 0) && !v.AggStar && !v.AggWithinGroup && !v.FuncVariadic {
				parts := nameParts(v.Funcname)
				if rr := b.resolveRoutine(parts, ast.OBJECT_FUNCTION); rr.fail == "" && rr.fact == walk.Absent && rr.missingSchema == "" && len(parts) > 0 {
					call, r = parts, rr
					return false
				}
			}
		}
		return true
	})
	if call != nil {
		b.fail(r.keys, walk.Absent, review.ObjectNotFound, notFound("function", qualified(call))) // 42883
		return true
	}
	return false
}

// fromName is a name a FROM item takes in its query's namespace. An
// unaliased table takes its own name, which another table of that name in
// another schema may share (pg: parse_relation.c — checkNameSpaceConflicts).
type fromName struct {
	name   string
	schema string
	table  bool
}

// duplicateFromName finds a name two FROM items of one level of a query
// both take, in a plain query or an arm of a set operation, outside
// subqueries.
func duplicateFromName(q ast.Node) (string, bool) {
	sel, ok := q.(*ast.SelectStmt)
	if !ok || sel == nil {
		return "", false
	}
	if sel.Op != ast.SETOP_NONE {
		if name, ok := duplicateFromName(sel.Larg); ok {
			return name, true
		}
		return duplicateFromName(sel.Rarg)
	}
	if sel.FromClause == nil {
		return "", false
	}
	var seen []fromName
	for _, item := range sel.FromClause.Items {
		names, dup, ok := fromItemNames(item)
		if ok {
			return dup, true
		}
		if dup, ok := conflict(seen, names); ok {
			return dup, true
		}
		seen = append(seen, names...)
	}
	return "", false
}

// fromItemNames lists the names a FROM item takes and reports a conflict
// inside a join. An item whose name is not followed here takes none.
func fromItemNames(n ast.Node) ([]fromName, string, bool) {
	switch v := n.(type) {
	case *ast.RangeVar:
		if v.Alias != nil && v.Alias.Aliasname != "" {
			return []fromName{{name: v.Alias.Aliasname}}, "", false
		}
		if v.Catalogname != "" {
			return nil, "", false
		}
		return []fromName{{name: v.Relname, schema: v.Schemaname, table: true}}, "", false
	case *ast.RangeSubselect:
		if v.Alias != nil && v.Alias.Aliasname != "" {
			return []fromName{{name: v.Alias.Aliasname}}, "", false
		}
	case *ast.RangeFunction:
		if v.Alias != nil && v.Alias.Aliasname != "" {
			return []fromName{{name: v.Alias.Aliasname}}, "", false
		}
	case *ast.JoinExpr:
		left, dup, ok := fromItemNames(v.Larg)
		if ok {
			return nil, dup, true
		}
		right, dup, ok := fromItemNames(v.Rarg)
		if ok {
			return nil, dup, true
		}
		if dup, ok := conflict(left, right); ok {
			return nil, dup, true
		}
		if v.Alias != nil && v.Alias.Aliasname != "" {
			// The alias hides the joined items.
			return []fromName{{name: v.Alias.Aliasname}}, "", false
		}
		return append(left, right...), "", false
	}
	return nil, "", false
}

func conflict(seen, names []fromName) (string, bool) {
	for _, a := range names {
		for _, s := range seen {
			if a.name != s.name {
				continue
			}
			if a.table && s.table && a.schema != s.schema {
				continue
			}
			return a.name, true
		}
	}
	return "", false
}

// refusedForm finds a form PostgreSQL refuses in any query with no outer
// statement, outside a data-modifying WITH query, which brings its own
// rules:
//   - a parameter, which nothing supplies (42P02);
//   - DEFAULT, which only INSERT and UPDATE take (42601);
//   - MERGE_ACTION(), which only MERGE's RETURNING takes (42601);
//   - an empty ARRAY[] not cast, whose type nothing determines (42P18);
//   - ARRAY(SELECT) or a scalar subquery of no column (42601);
//   - a WINDOW clause that defines a name twice (42P20) or refers to one
//     not defined before it (42704).
//
// pg: src/backend/parser/parse_expr.c — transformParamRef,
// transformArrayExpr, transformSubLink; parse_clause.c —
// transformWindowDefinitions
func refusedForm(q ast.Node) string {
	msg := ""
	ast.Inspect(q, func(n ast.Node) bool {
		if msg != "" {
			return false
		}
		switch v := n.(type) {
		case *ast.TypeName, *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.MergeStmt:
			return false
		case *ast.TypeCast:
			if _, ok := v.Arg.(*ast.A_ArrayExpr); ok {
				// The cast gives the array its type.
				return false
			}
		case *ast.ParamRef:
			msg = "Query references parameter `$" + strconv.Itoa(v.Number) + "`, which nothing supplies here."
		case *ast.SetToDefault:
			msg = unsupported("`DEFAULT` in a query")
		case *ast.MergeSupportFunc:
			msg = unsupported("`MERGE_ACTION()` outside `MERGE`")
		case *ast.A_ArrayExpr:
			if emptyArray(v) {
				msg = unsupported("an empty `ARRAY` with no cast to give its type")
			}
			// An element array takes its type from this one.
			return false
		case *ast.SubLink:
			if ast.SubLinkType(v.SubLinkType) == ast.ARRAY_SUBLINK || ast.SubLinkType(v.SubLinkType) == ast.EXPR_SUBLINK {
				if sel, ok := v.Subselect.(*ast.SelectStmt); ok && sel.Op == ast.SETOP_NONE && sel.ValuesLists == nil &&
					(sel.TargetList == nil || len(sel.TargetList.Items) == 0) {
					msg = unsupported("a subquery of no column as a value")
				}
			}
		case *ast.SelectStmt:
			if v.WindowClause != nil {
				defined := make(map[string]bool, len(v.WindowClause.Items))
				for _, item := range v.WindowClause.Items {
					wd, ok := item.(*ast.WindowDef)
					if !ok {
						continue
					}
					if wd.Refname != "" && !defined[wd.Refname] {
						msg = notFound("window", ident(wd.Refname))
						return false
					}
					if defined[wd.Name] {
						msg = unsupported("window " + review.Code(ident(wd.Name)) + " defined twice")
						return false
					}
					defined[wd.Name] = true
				}
			}
		}
		return msg == ""
	})
	return msg
}

// emptyArray reports whether an ARRAY[] holds no element at any depth.
func emptyArray(a *ast.A_ArrayExpr) bool {
	if a.Elements == nil {
		return true
	}
	for _, item := range a.Elements.Items {
		sub, ok := item.(*ast.A_ArrayExpr)
		if !ok || !emptyArray(sub) {
			return false
		}
	}
	return true
}
