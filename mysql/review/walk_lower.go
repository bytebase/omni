package review

import (
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// statement lowers one statement (spec §5), false when the walk does not
// follow it, which ends the walk before it.
func (b *builder) statement(n ast.Node) bool {
	switch v := n.(type) {
	case *ast.CreateDatabaseStmt:
		return b.createDatabase(v)
	case *ast.DropDatabaseStmt:
		return b.dropDatabase(v)
	case *ast.AlterDatabaseStmt:
		return b.alterDatabase(v)
	case *ast.UseStmt:
		return b.use(v)
	case *ast.CreateTableStmt:
		return b.createTable(v)
	case *ast.AlterTableStmt:
		return b.alterTable(v)
	case *ast.RenameTableStmt:
		return b.renameTable(v)
	case *ast.DropTableStmt:
		return b.dropTable(v)
	case *ast.TruncateStmt:
		return b.truncate(v)
	case *ast.CreateIndexStmt:
		return b.createIndex(v)
	case *ast.DropIndexStmt:
		return b.dropIndex(v)
	case *ast.CreateViewStmt:
		return b.createView(v)
	case *ast.AlterViewStmt:
		return b.alterView(v)
	case *ast.DropViewStmt:
		return b.dropView(v)
	case *ast.CreateTriggerStmt:
		return b.createTrigger(v)
	case *ast.DropTriggerStmt:
		return b.dropTrigger(v)
	case *ast.CreateFunctionStmt:
		return b.createRoutine(v)
	case *ast.AlterRoutineStmt:
		return b.alterRoutine(v)
	case *ast.DropRoutineStmt:
		return b.dropRoutine(v)
	case *ast.CreateEventStmt:
		return b.createEvent(v)
	case *ast.AlterEventStmt:
		return b.alterEvent(v)
	case *ast.DropEventStmt:
		return b.dropEvent(v)
	case *ast.InsertStmt:
		return b.dml(n, []*ast.TableRef{v.Table})
	case *ast.LoadDataStmt:
		return b.dml(n, []*ast.TableRef{v.Table})
	case *ast.UpdateStmt:
		return b.dml(n, singleTarget(v.Tables))
	case *ast.DeleteStmt:
		if v.Using != nil {
			return b.dml(n, nil)
		}
		return b.dml(n, singleTarget(v.Tables))
	case *ast.SetStmt:
		return b.set(v)
	case *ast.SelectStmt, *ast.DoStmt, *ast.ValuesStmt, *ast.TableStmt, *ast.ExplainStmt:
		if b.callsStoredFunction(n) {
			b.forgetTemporary()
		}
		return true
	case *ast.BeginStmt, *ast.CommitStmt, *ast.SavepointStmt, *ast.SetTransactionStmt, *ast.UnlockTablesStmt,
		*ast.ShowStmt, *ast.HelpStmt, *ast.AnalyzeTableStmt, *ast.OptimizeTableStmt, *ast.CheckTableStmt,
		*ast.ChecksumTableStmt, *ast.RepairTableStmt, *ast.FlushStmt,
		*ast.GrantStmt, *ast.RevokeStmt, *ast.GrantRoleStmt, *ast.RevokeRoleStmt, *ast.SetRoleStmt, *ast.SetDefaultRoleStmt,
		*ast.CreateUserStmt, *ast.AlterUserStmt, *ast.DropUserStmt, *ast.CreateRoleStmt, *ast.DropRoleStmt:
		// None changes a name the walk keeps. SAVEPOINT and RELEASE keep
		// what the change did; a ROLLBACK, which undoes it, is not here.
		return true
	}
	return false
}

// tableNeeds adds the check that a statement's table exists as a base
// table (1146 missing, 1347 a view).
func (b *builder) tableNeeds(r tableRef) {
	b.need(relKey(r.db, r.name), walk.Present, []string{relTable}, r.rng,
		"Table "+r.written+" does not exist", r.written+" is not a base table")
}

// databaseNeeds adds the check that a database exists (1049).
func (b *builder) databaseNeeds(db, written string, rng review.Range) {
	b.need(dbKey(db), walk.Present, nil, rng, "Database "+ident(written)+" does not exist", "")
}

func (b *builder) createDatabase(v *ast.CreateDatabaseStmt) bool {
	rng := b.anchor()
	if b.tooLong(v.Name, rng) {
		return true
	}
	d, ok := b.key(spaceDatabase, v.Name)
	if !ok {
		return false
	}
	k := dbKey(d)
	prior, _ := b.l.Lookup(k)
	if !v.IfNotExists {
		b.need(k, walk.Absent, nil, rng, "", "Database "+ident(v.Name)+" already exists") // 1007
		prior = walk.Absent
	}
	switch prior {
	case walk.Absent:
		b.Set(k, walk.Present, "")
		b.clearDatabase(d)
	case walk.Unknown:
		b.Set(k, walk.Present, "")
	}
	return true
}

func (b *builder) dropDatabase(v *ast.DropDatabaseStmt) bool {
	d, ok := b.key(spaceDatabase, v.Name)
	if !ok {
		return false
	}
	k := dbKey(d)
	if !v.IfExists {
		b.need(k, walk.Present, nil, b.anchor(), "Database "+ident(v.Name)+" does not exist", "") // 1008
	}
	b.Set(k, walk.Absent, "")
	b.clearDatabase(d)
	if !b.w.lost && d == b.w.current {
		b.Session(func() { b.w.noDatabase = true })
	}
	return true
}

// clearDatabase records that a database holds nothing, temporary tables
// aside: a session's temporary tables outlive DROP DATABASE of theirs
// (verified on 8.0 and 8.4).
func (b *builder) clearDatabase(d string) {
	for _, space := range []string{spaceRelation, spaceColumn, spaceIndex, spaceForeignKey, spaceTableForeignKey, spaceTrigger, spaceProcedure, spaceFunction, spaceEvent} {
		b.Clear(walk.Scope{Space: space, Parent: walk.Path(d), Deep: true})
	}
}

// alterDatabase follows ALTER DATABASE, of the current database when it
// names none, which fails once that is dropped (1046, verified).
func (b *builder) alterDatabase(v *ast.AlterDatabaseStmt) bool {
	if v.Name == "" {
		if b.w.noDatabase && !b.w.lost {
			b.fail(b.anchor(), []walk.Key{dbKey(b.w.current)}, noDatabaseMessage)
		}
		return true
	}
	d, ok := b.key(spaceDatabase, v.Name)
	if ok {
		b.databaseNeeds(d, v.Name, b.anchor())
	}
	return true
}

func (b *builder) use(v *ast.UseStmt) bool {
	d, ok := b.key(spaceDatabase, v.Database)
	if !ok {
		b.Session(func() { b.w.lost = true })
		return true
	}
	b.databaseNeeds(d, v.Database, b.anchor())
	b.Session(func() {
		b.w.current, b.w.noDatabase, b.w.lost = d, false, false
	})
	return true
}

// createTable follows CREATE [TEMPORARY] TABLE [IF NOT EXISTS], with a
// column list, LIKE, or AS SELECT.
func (b *builder) createTable(v *ast.CreateTableStmt) bool {
	if b.namesTooLong(v) {
		return true
	}
	r, st := b.ref(v.Table)
	switch st {
	case refFailed:
		return true
	case refUnknown:
		return false
	}
	// One object: every check holds against the state before the
	// statement, and the creation is a part of its own, one IF NOT EXISTS
	// may skip.
	b.Order = walk.Lookups
	if v.Table.Schema != "" {
		b.databaseNeeds(r.db, v.Table.Schema, r.rng)
	}
	// The new table keeps foreign keys as its ENGINE does, or for LIKE as
	// its source's does, which the walk must know.
	b.fks = !v.Temporary && keepsFKs(engineOf(v.Options))
	// The LIKE source is checked even when IF NOT EXISTS finds the table
	// (verified).
	if v.Like != nil {
		b.fks = false
		src, st := b.ref(v.Like)
		switch st {
		case refFailed:
			return true
		case refOK:
			f, temp, kind := b.resolve(src)
			if !temp && f != walk.Unknown {
				b.tableNeeds(src)
			}
			b.fks = !v.Temporary && !temp && f == walk.Present && kind == relTable && !b.w.ignoresFKs[relKey(src.db, src.name)]
		}
	}
	// A stored function the query calls may create or drop temporary
	// tables (verified).
	if v.Select != nil && b.callsStoredFunction(v.Select) {
		b.forgetTemporary()
	}
	if v.Temporary {
		tk := tempKey(r.db, r.name)
		prior, _ := b.l.Lookup(tk)
		if !v.IfNotExists {
			b.need(tk, walk.Absent, nil, r.rng, "", "Table "+r.written+" already exists") // 1050
		}
		if v.IfNotExists && prior != walk.Absent {
			// A note when the table exists, its text unread (verified).
			if prior == walk.Unknown {
				b.Set(tk, walk.Present, relTable)
			}
			return true
		}
		if b.tableTextFails(v, r) {
			return true
		}
		b.Set(tk, walk.Present, relTable)
		return true
	}
	k := relKey(r.db, r.name)
	prior, _ := b.l.Lookup(k)
	if !v.IfNotExists {
		b.need(k, walk.Absent, nil, r.rng, "", "Table "+r.written+" already exists") // 1050
		prior = walk.Absent
	}
	if prior == walk.Present {
		return true // a note: nothing is created
	}
	b.Next()
	if prior == walk.Unknown {
		// The table may exist, as a table or a view, and the statement be
		// a note: what it would create may or may not be there.
		b.Optional()
	}
	if b.tableTextFails(v, r) {
		return true
	}
	if b.fks {
		for _, c := range v.Constraints {
			if c.Type != ast.ConstrForeignKey {
				continue
			}
			// A foreign key's name is the database's (1826, verified).
			if n, ok := b.key(spaceForeignKey, c.Name); ok {
				b.need(objKey(spaceForeignKey, r.db, n), walk.Absent, nil, rangeOf(c.Loc), "", "Foreign key "+ident(c.Name)+" already exists in the database") // 1826; 1022 on 5.7
			}
			b.foreignKeyTarget(r, c)
		}
	}
	b.created(v, r)
	// A table that may have been there keeps an engine the walk does not
	// know.
	ignores := !b.fks || prior == walk.Unknown
	b.Session(func() {
		if ignores {
			b.w.ignoresFKs[k] = true
		} else {
			delete(b.w.ignoresFKs, k)
		}
	})
	if prior == walk.Unknown {
		b.Next()
		b.Set(k, walk.Present, "")
	}
	return true
}

// namesTooLong fails a CREATE TABLE naming a table, column, index, or
// foreign key longer than its space allows. The name after FOREIGN KEY
// counts only without a CONSTRAINT symbol, which names the index in its
// place (verified).
func (b *builder) namesTooLong(v *ast.CreateTableStmt) bool {
	if v.Table != nil && b.tooLong(v.Table.Name, rangeOf(v.Table.Loc)) {
		return true
	}
	for _, c := range v.Columns {
		if b.tooLong(c.Name, rangeOf(c.Loc)) {
			return true
		}
	}
	for _, c := range v.Constraints {
		if name := b.constraintName(c); b.tooLong(name, rangeOf(c.Loc)) {
			return true
		}
	}
	return false
}

// constraintName returns the name a table constraint gives, for its
// length: the name itself, which a CHECK's length follows, or a foreign
// key's keyName.
func (*builder) constraintName(c *ast.Constraint) string {
	if c.Type != ast.ConstrForeignKey {
		return c.Name
	}
	return keyName(c)
}

// tableTextFails fails a CREATE TABLE its text alone refutes (verified):
// a column named twice (1060), two indexes of one name (1061), two foreign
// keys of one name (1826, or 1061 for the index each would add), two
// primary keys (1068), or an index named PRIMARY (1280).
func (b *builder) tableTextFails(v *ast.CreateTableStmt, r tableRef) bool {
	columns := make(map[string]bool)
	for _, c := range v.Columns {
		k, ok := b.key(spaceColumn, c.Name)
		if !ok {
			continue
		}
		if columns[k] {
			b.fail(rangeOf(c.Loc), nil, "Column "+ident(c.Name)+" is specified twice in table "+r.written)
			return true
		}
		columns[k] = true
	}
	primaries := 0
	indexes := make(map[string]bool)
	foreignKeys := make(map[string]bool)
	for _, c := range v.Columns {
		for _, cc := range c.Constraints {
			if cc.Type == ast.ColConstrPrimaryKey {
				primaries++
			}
		}
	}
	for _, c := range v.Constraints {
		switch c.Type {
		case ast.ConstrPrimaryKey:
			primaries++
			continue
		case ast.ConstrForeignKey:
			if k, ok := b.key(spaceForeignKey, c.Name); ok && b.fks {
				if foreignKeys[k] {
					b.fail(rangeOf(c.Loc), nil, "Foreign key "+ident(c.Name)+" is specified twice in table "+r.written)
					return true
				}
				foreignKeys[k] = true
			}
			continue
		case ast.ConstrCheck:
			continue
		}
		k, ok := b.key(spaceIndex, c.Name)
		switch {
		case !ok:
			continue
		case k == primaryIndex:
			b.fail(rangeOf(c.Loc), nil, "Index name PRIMARY is reserved for the primary key")
			return true
		case indexes[k]:
			b.fail(rangeOf(c.Loc), nil, "Index "+ident(c.Name)+" is specified twice in table "+r.written)
			return true
		}
		indexes[k] = true
	}
	if primaries > 1 {
		b.fail(r.rng, nil, "Multiple primary keys for table "+r.written)
		return true
	}
	return false
}

// created records what a CREATE TABLE makes: the table, holding exactly
// its column list's columns and named keys, with generated names as
// May-effects. A table copied with LIKE or filled from a query holds what
// the walk does not follow. A table with no primary key may get an
// invisible one, with its column (sql_generate_invisible_primary_key).
func (b *builder) created(v *ast.CreateTableStmt, r tableRef) {
	path := tablePath(r.db, r.name)
	b.Set(relKey(r.db, r.name), walk.Present, relTable)
	b.Clear(walk.Scope{Parent: path, Deep: true})
	if v.Like != nil || v.Select != nil {
		b.Forget(walk.Scope{Parent: path, Deep: true})
	}
	primary := false
	for _, c := range v.Columns {
		if k, ok := b.key(spaceColumn, c.Name); ok {
			b.Set(colKey(r.db, r.name, k), walk.Present, "")
		}
		primary = b.columnKeys(r, c) || primary
	}
	for _, c := range v.Constraints {
		primary = b.constraint(r, c) || primary
	}
	if !primary {
		b.Set(idxKey(r.db, r.name, primaryIndex), walk.Unknown, "")
		b.Set(colKey(r.db, r.name, gipkColumn), walk.Unknown, "")
	}
}

// columnKeys records the keys a column definition declares inline, and
// reports whether one is the primary key. An inline REFERENCES makes no
// foreign key (verified on 8.0); its generated name is a May-effect all
// the same.
func (b *builder) columnKeys(r tableRef, c *ast.ColumnDef) bool {
	primary := false
	for _, cc := range c.Constraints {
		switch cc.Type {
		case ast.ColConstrPrimaryKey:
			b.Set(idxKey(r.db, r.name, primaryIndex), walk.Present, "")
			primary = true
		case ast.ColConstrUnique:
			b.MayAdd(walk.Scope{Space: spaceIndex, Parent: tablePath(r.db, r.name)})
		case ast.ColConstrReferences:
			b.mayAddForeignKeys(r)
		}
	}
	return primary
}

// constraint records a table constraint a CREATE or ALTER TABLE adds, and
// reports whether it is the primary key: named keys and foreign keys are
// Present, unnamed ones generated names. A foreign key adds an index of
// its own when none serves it, named after it or its first column.
func (b *builder) constraint(r tableRef, c *ast.Constraint) bool {
	indexes := walk.Scope{Space: spaceIndex, Parent: tablePath(r.db, r.name)}
	switch c.Type {
	case ast.ConstrPrimaryKey:
		b.Set(idxKey(r.db, r.name, primaryIndex), walk.Present, "")
		return true
	case ast.ConstrCheck:
	case ast.ConstrForeignKey:
		// An engine that ignores foreign keys still adds the index (verified).
		k, ok := b.key(spaceForeignKey, c.Name)
		switch {
		case ok && b.fks:
			fk, own := fkKeys(r, k)
			b.Set(fk, walk.Present, "")
			b.Set(own, walk.Present, "")
		case ok:
			fk, own := fkKeys(r, k)
			b.Set(fk, walk.Unknown, "")
			b.Set(own, walk.Unknown, "")
		default:
			b.mayAddForeignKeys(r)
		}
		b.MayAdd(indexes)
	default:
		if k, ok := b.key(spaceIndex, c.Name); ok {
			b.Set(idxKey(r.db, r.name, k), walk.Present, "")
		} else {
			b.MayAdd(indexes)
		}
	}
	return false
}

// mayAddForeignKeys records a foreign key of a generated name.
func (b *builder) mayAddForeignKeys(r tableRef) {
	p := b.generatedFKs(r.name)
	b.MayAdd(walk.Scope{Space: spaceForeignKey, Parent: walk.Path(r.db), Prefix: p})
	b.MayAdd(walk.Scope{Space: spaceTableForeignKey, Parent: tablePath(r.db, r.name), Prefix: p})
}

// foreignKeyTarget checks a foreign key's referenced table while the
// session checks foreign keys: it must be a base table (1824, a view or a
// temporary table too, which shadows none here; 1215 on 5.7; verified). A
// table that refers to itself needs nothing.
//
// An ALTER TABLE that also renames its table resolves the parent by rules
// of its own (verified on 8.0.46): an unqualified one in the database the
// table moves to, the table's old name as the table itself, and its new
// name as a missing table (1824). The walk checks no parent there.
func (b *builder) foreignKeyTarget(r tableRef, c *ast.Constraint) {
	if b.w.fkChecks != review.SwitchOn || c.RefTable == nil || b.renames {
		return
	}
	// An unqualified parent is in the child's database (verified).
	ref := *c.RefTable
	if ref.Schema == "" {
		db := r.db
		if strings.HasPrefix(db, nonASCII) {
			return
		}
		ref.Schema = db
	}
	t, st := b.ref(&ref)
	if st != refOK || t.db == r.db && t.name == r.name {
		return
	}
	t.written = table(c.RefTable)
	b.need(relKey(t.db, t.name), walk.Present, []string{relTable}, t.rng, "Foreign key target table "+t.written+" does not exist", "Foreign key target "+t.written+" is not a base table")
}
