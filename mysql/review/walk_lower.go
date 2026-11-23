package review

import (
	"strconv"
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
	case *ast.CreateSpatialRefSysStmt:
		return b.createSpatialRefSys(v)
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
	case *ast.ExplainStmt:
		if b.callsStoredFunction(n) {
			b.forgetTemporary()
		}
		return b.explain(v)
	case *ast.SelectStmt, *ast.DoStmt:
		if b.callsStoredFunction(n) {
			b.forgetTemporary()
		}
		return true
	case *ast.BeginStmt, *ast.CommitStmt, *ast.SavepointStmt, *ast.SetTransactionStmt, *ast.UnlockTablesStmt,
		*ast.ShowStmt, *ast.HelpStmt, *ast.AnalyzeTableStmt, *ast.OptimizeTableStmt, *ast.CheckTableStmt,
		*ast.ChecksumTableStmt, *ast.RepairTableStmt, *ast.FlushStmt,
		*ast.GrantStmt, *ast.RevokeStmt, *ast.GrantRoleStmt, *ast.RevokeRoleStmt, *ast.SetRoleStmt, *ast.SetDefaultRoleStmt,
		*ast.CreateUserStmt, *ast.AlterUserStmt, *ast.DropUserStmt, *ast.RenameUserStmt,
		*ast.CreateRoleStmt, *ast.DropRoleStmt, *ast.SetPasswordStmt:
		// None changes a name the walk keeps. SAVEPOINT and RELEASE keep
		// what the change did; a ROLLBACK, which undoes it, is not here.
		// RENAME USER and SET PASSWORD join the other user statements:
		// both parsed already, and both reached the default arm, so
		// either one in the middle of a change script ended the walk and
		// hid every finding below it.
		return true
	}
	return false
}

// explain lowers EXPLAIN. Reporting a plan changes nothing, so a plain
// EXPLAIN is followed without recording anything. ANALYZE runs the statement
// it carries, so that form is lowered as the statement it performs. No DDL is
// explainable on MySQL 8.0 (mysql/parser/explain_forms_test.go), so what
// ANALYZE reaches here is a query or a DML, not a schema change.
func (b *builder) explain(v *ast.ExplainStmt) bool {
	if !v.Analyze || v.Stmt == nil {
		return true
	}
	return b.statement(v.Stmt)
}

// tableNeeds adds the check that a statement's table exists as a base
// table (1146 missing, 1347 a view).
func (b *builder) tableNeeds(r tableRef) {
	b.needKind(relKey(r.db, r.name), []string{relTable}, r.rng, notFound("table", r.written, ""), r.written, "table")
}

// databaseNeeds adds the check that a database exists (1049).
func (b *builder) databaseNeeds(db, written string, rng review.Range) {
	b.need(dbKey(db), walk.Present, nil, rng, notFound("database", ident(written), ""), "")
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
		b.need(k, walk.Absent, nil, rng, "", alreadyExists("database", ident(v.Name), "")) // 1007
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
		b.need(k, walk.Present, nil, b.anchor(), notFound("database", ident(v.Name), ""), "") // 1008
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
	for _, space := range []string{spaceRelation, spaceColumn, spaceIndex, spaceForeignKey, spaceTableForeignKey, spaceCheckConstraint, spaceTrigger, spaceProcedure, spaceFunction, spaceEvent} {
		b.Clear(walk.Scope{Space: space, Parent: walk.Path(d), Deep: true})
	}
	b.Session(func() {
		for gk := range b.w.fkGeneratedIndexes {
			if names := walk.Names(gk.Parent); len(names) > 0 && names[0] == d {
				delete(b.w.fkGeneratedIndexes, gk)
			}
		}
		for rk := range b.w.engines {
			if names := walk.Names(rk.Parent); len(names) > 0 && names[0] == d {
				delete(b.w.engines, rk)
			}
		}
		for rk := range b.w.partitionState {
			if names := walk.Names(rk.Parent); len(names) > 0 && names[0] == d {
				delete(b.w.partitionState, rk)
			}
		}
		for rk := range b.w.tableIndexes {
			if names := walk.Names(rk.Parent); len(names) > 0 && names[0] == d {
				delete(b.w.tableIndexes, rk)
				delete(b.w.tableIndexesComplete, rk)
			}
		}
		for rk := range b.w.tableFKs {
			if names := walk.Names(rk.Parent); len(names) > 0 && names[0] == d {
				delete(b.w.tableFKs, rk)
			}
		}
		for ck := range b.w.colTypes {
			if names := walk.Names(ck.Parent); len(names) > 0 && names[0] == d {
				delete(b.w.colTypes, ck)
			}
		}
		for pRel, refs := range b.w.referencedByFK {
			if names := walk.Names(pRel.Parent); len(names) > 0 && names[0] == d {
				delete(b.w.referencedByFK, pRel)
				continue
			}
			for ref := range refs {
				if names := walk.Names(ref.child.Parent); len(names) > 0 && names[0] == d {
					delete(refs, ref)
				}
			}
			if len(refs) == 0 {
				delete(b.w.referencedByFK, pRel)
			}
		}
	})
}

// alterDatabase follows ALTER DATABASE, of the current database when it
// names none, which fails once that is dropped (1046, verified).
func (b *builder) alterDatabase(v *ast.AlterDatabaseStmt) bool {
	if v.Name == "" {
		if b.w.noDatabase && !b.w.lost {
			b.fail(b.anchor(), []walk.Key{dbKey(b.w.current)}, review.NoCurrentSchema, noDatabase("ALTER DATABASE"))
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
			b.needFree(tk, r.rng, r.written, relTable) // 1050
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
		b.needFree(k, r.rng, r.written, relTable) // 1050
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
				b.need(objKey(spaceForeignKey, r.db, n), walk.Absent, nil, rangeOf(c.Loc), "", alreadyExists("foreign key", ident(c.Name), "")) // 1826; 1022 on 5.7
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
// primary keys (1068), an index named PRIMARY (1280), empty table (4028),
// column/option/check/index/partition/select errors.
func (b *builder) tableTextFails(v *ast.CreateTableStmt, r tableRef) bool {
	if len(v.Columns) == 0 && v.Like == nil && v.Select == nil {
		b.fail(r.rng, nil, review.UnsupportedForm, "Table "+review.Code(r.written)+" must have at least one column.")
		return true
	}
	if b.tableOptionsFails(v.Options, engineOf(v.Options), v.Like == nil, r.rng) {
		return true
	}
	columns := make(map[string]bool)
	colBaseTypes := make(map[string]string)
	for _, c := range v.Columns {
		if b.columnDefTextFails(c) {
			return true
		}
		k, ok := b.key(spaceColumn, c.Name)
		if !ok {
			continue
		}
		if columns[k] {
			b.fail(rangeOf(c.Loc), nil, review.DuplicateName, statesTwice("declares", "column", ident(c.Name)))
			return true
		}
		columns[k] = true
		if bt := astBaseType(c.TypeName); bt != "" {
			colBaseTypes[k] = bt
		}
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
	hasFK := false
	for _, c := range v.Constraints {
		switch c.Type {
		case ast.ConstrPrimaryKey:
			primaries++
			continue
		case ast.ConstrForeignKey:
			hasFK = true
			if k, ok := b.key(spaceForeignKey, c.Name); ok && b.fks {
				if foreignKeys[k] {
					b.fail(rangeOf(c.Loc), nil, review.DuplicateName, statesTwice("declares", "foreign key", ident(c.Name)))
					return true
				}
				foreignKeys[k] = true
			}
			continue
		case ast.ConstrCheck:
			if c.Expr != nil && !checkExprIsBoolean(c.Expr) {
				b.fail(rangeOf(c.Loc), nil, review.UnsupportedForm, "Check constraint in "+review.Code(r.written)+" must be a boolean expression.")
				return true
			}
			continue
		}
		k, ok := b.key(spaceIndex, c.Name)
		switch {
		case !ok:
			continue
		case k == primaryIndex:
			b.fail(rangeOf(c.Loc), nil, review.ReservedName, reservedPrimary)
			return true
		case indexes[k]:
			b.fail(rangeOf(c.Loc), nil, review.DuplicateName, statesTwice("declares", "index", ident(c.Name)))
			return true
		}
		indexes[k] = true
	}
	if primaries > 1 {
		b.fail(r.rng, nil, review.NameTaken, hasPrimaryKey(r.written))
		return true
	}
	if b.partitionClausePreFKFails(v.Partitions, r.rng) {
		return true
	}
	if b.fks && v.Partitions != nil && hasFK {
		b.fail(r.rng, nil, review.UnsupportedForm, "Partitioned table "+review.Code(r.written)+" cannot define foreign keys.")
		return true
	}
	completeCols := v.Like == nil && v.Select == nil
	lookupCol := func(colName string) (walk.Fact, string) {
		k, ok := b.key(spaceColumn, colName)
		if !ok {
			return walk.Unknown, ""
		}
		if columns[k] {
			return walk.Present, colBaseTypes[k]
		}
		if completeCols {
			return walk.Absent, ""
		}
		return walk.Unknown, ""
	}
	partitioned := v.Partitions != nil
	for _, c := range v.Constraints {
		if c.Type == ast.ConstrCheck {
			continue
		}
		parts := constraintKeyParts(c)
		isFK := c.Type == ast.ConstrForeignKey
		fulltext := c.Type == ast.ConstrFulltextIndex
		spatial := c.Type == ast.ConstrSpatialIndex || isSpatialIndexType(c.IndexType)
		crng := rangeOf(c.Loc)
		if crng == (review.Range{}) {
			crng = r.rng
		}
		// Check missing column directly as a text failure for CREATE TABLE
		// (the table does not exist in the ledger yet).
		for _, p := range parts {
			if p.functional && isBareColumnRef(p.expr) {
				b.fail(crng, nil, review.UnsupportedForm, "Functional index key part cannot be a bare column reference.")
				return true
			}
		}
		seenPart := make(map[string]bool, len(parts))
		for _, p := range parts {
			if p.functional || p.name == "" {
				continue
			}
			k, ok := b.key(spaceColumn, p.name)
			if !ok {
				continue
			}
			if seenPart[k] {
				b.fail(crng, nil, review.DuplicateName, statesTwice("names", "column", ident(p.name)))
				return true
			}
			seenPart[k] = true
		}
		for _, p := range parts {
			if p.functional || p.name == "" {
				continue
			}
			if st, _ := lookupCol(p.name); st == walk.Absent {
				b.fail(crng, nil, review.ObjectNotFound, notFound("key column", ident(p.name), r.written))
				return true
			}
		}
		if fulltext && partitioned {
			b.fail(crng, nil, review.UnsupportedForm, "Partitioned table "+review.Code(r.written)+" doesn't support "+review.Code("FULLTEXT")+" indexes.")
			return true
		}
		for _, p := range parts {
			if p.functional || p.name == "" {
				continue
			}
			st, bt := lookupCol(p.name)
			if st != walk.Present || bt == "" {
				continue
			}
			if fulltext && !isFulltextBaseType(bt) {
				b.fail(crng, nil, review.UnsupportedForm, "Column "+review.Code(ident(p.name))+" cannot be part of a "+review.Code("FULLTEXT")+" index.")
				return true
			}
			if spatial && !isGeometryBaseType(bt) {
				b.fail(crng, nil, review.UnsupportedForm, "Column "+review.Code(ident(p.name))+" cannot be part of a "+review.Code("SPATIAL")+" index.")
				return true
			}
			if !fulltext && !spatial && !isFK && p.length > 0 && !isPrefixableBaseType(bt) {
				b.fail(crng, nil, review.UnsupportedForm, "Column "+review.Code(ident(p.name))+" cannot use an index prefix length.")
				return true
			}
		}
		if isFK && b.fks {
			if len(c.RefColumns) == 0 || (len(parts) > 0 && len(parts) != len(c.RefColumns)) {
				b.fail(crng, nil, review.UnsupportedForm, "Foreign key in "+review.Code(r.written)+" must name matching referenced columns.")
				return true
			}
		}
	}
	if completeCols {
		autoCount := 0
		for _, c := range v.Columns {
			if !colHasAutoIncrement(c) {
				continue
			}
			autoCount++
			if autoCount > 1 {
				b.fail(rangeOf(c.Loc), nil, review.UnsupportedForm, "Table "+review.Code(r.written)+" can have only one "+review.Code("AUTO_INCREMENT")+" column and it must be indexed.")
				return true
			}
			if !colHasInlineKey(c) {
				keyed := false
				for _, con := range v.Constraints {
					if con.Type == ast.ConstrCheck {
						continue
					}
					if keyPartsCoverColumn(constraintKeyParts(con), c.Name) {
						keyed = true
						break
					}
				}
				if !keyed {
					b.fail(rangeOf(c.Loc), nil, review.UnsupportedForm, "Column "+review.Code(ident(c.Name))+" is "+review.Code("AUTO_INCREMENT")+" and must be defined as a key.")
					return true
				}
			}
		}
	}
	if v.Partitions != nil {
		var uqIdxs map[string]*indexMeta
		if completeCols {
			uqIdxs = make(map[string]*indexMeta)
			for _, c := range v.Columns {
				for _, cc := range c.Constraints {
					if cc.Type == ast.ColConstrPrimaryKey || cc.Type == ast.ColConstrUnique {
						uqIdxs["col:"+strings.ToLower(c.Name)] = &indexMeta{
							unique:  true,
							primary: cc.Type == ast.ColConstrPrimaryKey,
							cols:    []string{strings.ToLower(c.Name)},
						}
					}
				}
			}
			for i, con := range v.Constraints {
				if con.Type != ast.ConstrPrimaryKey && con.Type != ast.ConstrUnique {
					continue
				}
				parts := constraintKeyParts(con)
				im := &indexMeta{unique: true, primary: con.Type == ast.ConstrPrimaryKey}
				for _, p := range parts {
					if p.functional {
						break
					}
					if p.length > 0 {
						im.hasPrefix = true
					}
					im.cols = append(im.cols, strings.ToLower(p.name))
				}
				uqIdxs["con:"+strconv.Itoa(i)] = im
			}
		}
		if b.partitionClauseColumnsAndUniqueFails(r, v.Partitions, r.rng, func(colName string) walk.Fact {
			st, _ := lookupCol(colName)
			return st
		}, uqIdxs, completeCols) {
			return true
		}
	}
	if v.Select != nil {
		if col, ok := selectFromlessColumn(v.Select); ok {
			b.fail(r.rng, nil, review.ObjectNotFound, "Column "+review.Code(ident(col))+" doesn't exist in the query.")
			return true
		}
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
	b.clearTableFKGenerated(r.db, r.name)
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
	rel := relKey(r.db, r.name)
	eng := engineOf(v.Options)
	hasLikeOrSelect := v.Like != nil || v.Select != nil
	var likeKey walk.Key
	hasLike := false
	if v.Like != nil {
		if src, st := b.ref(v.Like); st == refOK {
			likeKey = relKey(src.db, src.name)
			hasLike = true
		}
	}
	colTypesMap := make(map[walk.Key]string)
	for _, c := range v.Columns {
		if k, ok := b.key(spaceColumn, c.Name); ok {
			if bt := astBaseType(c.TypeName); bt != "" {
				colTypesMap[colKey(r.db, r.name, k)] = bt
			}
		}
	}
	var pinfo *tablePartitionInfo
	if !hasLike {
		pinfo = partitionInfoFromClause(v.Partitions)
	}
	idxsMap := make(map[string]*indexMeta)
	idxComplete := !hasLikeOrSelect
	for _, c := range v.Columns {
		for _, cc := range c.Constraints {
			switch cc.Type {
			case ast.ColConstrPrimaryKey:
				idxsMap[primaryIndex] = &indexMeta{primary: true, unique: true, cols: []string{strings.ToLower(c.Name)}}
			case ast.ColConstrUnique:
				idxComplete = false
			}
		}
	}
	fksMap := make(map[string][]string)
	for _, con := range v.Constraints {
		switch con.Type {
		case ast.ConstrCheck:
		case ast.ConstrForeignKey:
			idxComplete = false
			if b.fks {
				if k, ok := b.key(spaceForeignKey, con.Name); ok {
					var cols []string
					for _, p := range constraintKeyParts(con) {
						if !p.functional && p.name != "" {
							cols = append(cols, strings.ToLower(p.name))
						}
					}
					fksMap[k] = cols
				}
			}
		case ast.ConstrPrimaryKey:
			im := &indexMeta{primary: true, unique: true}
			for _, p := range constraintKeyParts(con) {
				if p.functional {
					break
				}
				if p.length > 0 {
					im.hasPrefix = true
				}
				im.cols = append(im.cols, strings.ToLower(p.name))
			}
			idxsMap[primaryIndex] = im
		default:
			if k, ok := b.key(spaceIndex, con.Name); ok {
				im := &indexMeta{
					unique:   con.Type == ast.ConstrUnique,
					fulltext: con.Type == ast.ConstrFulltextIndex,
					spatial:  con.Type == ast.ConstrSpatialIndex || isSpatialIndexType(con.IndexType),
				}
				for _, p := range constraintKeyParts(con) {
					if p.functional {
						break
					}
					if p.length > 0 {
						im.hasPrefix = true
					}
					im.cols = append(im.cols, strings.ToLower(p.name))
				}
				idxsMap[k] = im
			} else {
				idxComplete = false
			}
		}
	}
	b.Session(func() {
		if hasLike {
			if srcEng, ok := b.w.engines[likeKey]; ok {
				b.w.engines[rel] = srcEng
			}
			if srcP, ok := b.w.partitionState[likeKey]; ok {
				b.w.partitionState[rel] = srcP.clone()
			}
		} else {
			b.w.engines[rel] = eng
			b.w.partitionState[rel] = pinfo
		}
		for ck, bt := range colTypesMap {
			b.w.colTypes[ck] = bt
		}
		if idxComplete {
			b.w.tableIndexes[rel] = idxsMap
			b.w.tableIndexesComplete[rel] = true
		}
		if len(fksMap) > 0 {
			b.w.tableFKs[rel] = fksMap
		}
	})
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
			ik := idxKey(r.db, r.name, primaryIndex)
			b.Set(ik, walk.Present, "")
			b.onIndexAdded(r, ik, true)
			primary = true
		case ast.ColConstrUnique:
			b.MayAdd(walk.Scope{Space: spaceIndex, Parent: tablePath(r.db, r.name)})
			b.onIndexAdded(r, walk.Key{}, false)
		case ast.ColConstrReferences:
			b.mayAddForeignKeys(r)
		case ast.ColConstrCheck:
			if k, ok := b.key(spaceCheckConstraint, cc.Name); ok {
				b.Set(chkKey(r.db, r.name, k), walk.Present, "")
			} else {
				b.MayAdd(walk.Scope{Space: spaceCheckConstraint, Parent: tablePath(r.db, r.name)})
			}
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
		ik := idxKey(r.db, r.name, primaryIndex)
		b.Set(ik, walk.Present, "")
		b.onIndexAdded(r, ik, true)
		return true
	case ast.ConstrCheck:
		if k, ok := b.key(spaceCheckConstraint, c.Name); ok {
			b.Set(chkKey(r.db, r.name, k), walk.Present, "")
		} else {
			b.MayAdd(walk.Scope{Space: spaceCheckConstraint, Parent: tablePath(r.db, r.name)})
		}
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
			ik := idxKey(r.db, r.name, k)
			b.Set(ik, walk.Present, "")
			b.onIndexAdded(r, ik, true)
		} else {
			b.MayAdd(indexes)
			b.onIndexAdded(r, walk.Key{}, false)
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
	b.needKind(relKey(t.db, t.name), []string{relTable}, t.rng, notFound("referenced table", t.written, ""), t.written, "table a foreign key can reference")
}
