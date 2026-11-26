package review

import (
	"strings"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// tableOf resolves the table an index or trigger statement names, which
// must be a base table. It returns false when the walk does not follow
// the statement on it: a temporary table or one it cannot tell.
func (b *builder) tableOf(tr *ast.TableRef) (tableRef, bool, bool) {
	r, st := b.ref(tr)
	switch st {
	case refFailed:
		return r, false, true
	case refUnknown:
		return r, false, false
	}
	f, temp, _ := b.resolve(r)
	if temp || f == walk.Unknown {
		return r, false, false
	}
	b.tableNeeds(r)
	return r, true, true
}

func (b *builder) createIndex(v *ast.CreateIndexStmt) bool {
	r, ok, lowered := b.tableOf(v.Table)
	if !ok {
		return lowered
	}
	rng := b.anchor()
	if b.tooLong(v.IndexName, rng) {
		return true
	}
	k, ok := b.key(spaceIndex, v.IndexName)
	if !ok {
		return false
	}
	if k == primaryIndex {
		b.fail(rng, nil, review.ReservedName, reservedPrimary) // 1280
		return true
	}
	if b.indexAlgorithmFails(v, r) || b.createUniqueFails(v, r) {
		return true
	}
	ik := idxKey(r.db, r.name, k)
	idxFact, _ := b.l.Lookup(ik)
	replacedFKIdx := false
	if !v.Fulltext && !v.Spatial && !v.Vector && indexTypeSupportsFK(v.IndexType, v.Options) {
		replacedFKIdx = b.replaceFKGeneratedIndexes(r, ik, v.Columns)
	}
	tf, _, tkind := b.resolve(r)
	switch {
	case v.OrReplace:
		// It replaces an index of the name (verified).
	case v.IfNotExists:
		if idxFact == walk.Present {
			return true
		}
		if idxFact == walk.Unknown {
			b.Next()
			b.Optional()
		}
	case !replacedFKIdx:
		b.need(ik, walk.Absent, nil, rng, "", alreadyExists("index", ident(v.IndexName), r.written)) // 1061
	}
	if tf == walk.Present && tkind == relTable && (!v.IfNotExists || idxFact == walk.Absent) {
		b.checkIndexColumns(r, v.Columns)
		if !v.Spatial && !v.Fulltext && !v.Vector && usesRTree(v.IndexType, v.Options) {
			b.fail(rng, nil, review.UnsupportedForm, notRTree) // 1210
			return true
		}
		if !b.w.ignoresFKs[relKey(r.db, r.name)] && b.keyPartTypesFail(createIndexKeyKind(v), v.Columns, b.typeOf(r)) {
			return true
		}
	}
	b.Set(ik, walk.Present, "")
	return true
}

func (b *builder) dropIndex(v *ast.DropIndexStmt) bool {
	r, st := b.ref(v.Table)
	switch st {
	case refFailed:
		return true
	case refUnknown:
		return false
	}
	f, temp, kind := b.resolve(r)
	if temp || f == walk.Unknown {
		return false
	}
	if f == walk.Present && kind == relSequence {
		if v.IfExists {
			return true
		}
		b.fail(b.anchor(), []walk.Key{relKey(r.db, r.name)}, review.ObjectNotFound, notFound("index", ident(v.Name), r.written)) // 1091
		return true
	}
	b.tableNeeds(r)
	if f == walk.Present && kind == relTable && b.fkIndexDropFails(r, []string{v.Name}, nil, b.anchor()) {
		return true
	}
	return b.dropIndexName(r, v.Name, b.anchor(), v.IfExists)
}

func (b *builder) createView(v *ast.CreateViewStmt) bool {
	r, st := b.ref(v.Name)
	switch st {
	case refFailed:
		return true
	case refUnknown:
		return false
	}
	if b.tooLong(v.Name.Name, r.rng) {
		return true
	}
	if v.Name.Schema != "" {
		b.databaseNeeds(r.db, v.Name.Schema, r.rng) // 1049
	}
	k := relKey(r.db, r.name)
	if v.IfNotExists {
		// IF NOT EXISTS makes a name already taken a note, after the body is
		// judged: mariadb:11.8.8 runs `CREATE VIEW IF NOT EXISTS a AS SELECT 1`
		// over a view or a table a and changes nothing.
		switch f, _ := b.l.Lookup(k); f {
		case walk.Unknown:
			return false
		case walk.Present:
			if ref := selectBareColumnNoFrom(v.Select); ref != nil {
				b.fail(rangeOf(ref.Loc), nil, review.ObjectNotFound, notFound("column", ident(ref.Column), "")) // 1054
				return true
			}
			b.viewQueryFails(v.Select, len(v.Columns) > 0)
			return true
		}
	}
	if v.OrReplace {
		if f, _ := b.l.Lookup(k); f == walk.Present {
			b.needKind(k, []string{relView}, r.rng, "", r.written, "view") // 1347
		}
	} else {
		b.needFree(k, r.rng, r.written, relView) // 1050
	}
	if f, kind := b.l.Lookup(k); f == walk.Absent || (v.OrReplace && f == walk.Present && kind == relView) {
		if ref := selectBareColumnNoFrom(v.Select); ref != nil {
			b.fail(rangeOf(ref.Loc), nil, review.ObjectNotFound, notFound("column", ident(ref.Column), "")) // 1054
			return true
		}
		if b.viewQueryFails(v.Select, len(v.Columns) > 0) {
			return true
		}
	}
	b.Set(k, walk.Present, relView)
	return true
}

func (b *builder) alterView(v *ast.AlterViewStmt) bool {
	r, st := b.ref(v.Name)
	switch st {
	case refFailed:
		return true
	case refOK:
		k := relKey(r.db, r.name)
		b.needKind(k, []string{relView}, r.rng, notFound("view", r.written, ""), r.written, "view") // 1146, 1347
		if f, kind := b.l.Lookup(k); f == walk.Present && kind == relView {
			if ref := selectBareColumnNoFrom(v.Select); ref != nil {
				b.fail(rangeOf(ref.Loc), nil, review.ObjectNotFound, notFound("column", ident(ref.Column), "")) // 1054
				return true
			}
			if b.viewQueryFails(v.Select, len(v.Columns) > 0) {
				return true
			}
		}
	}
	return true
}

// dropView follows DROP VIEW, which looks every view up first, and drops
// those that exist before it fails on one that does not (verified).
func (b *builder) dropView(v *ast.DropViewStmt) bool {
	b.Order = walk.Lookups
	for _, t := range v.Views {
		b.Next()
		r, st := b.ref(t)
		switch st {
		case refFailed:
			return true
		case refUnknown:
			return false
		}
		k := relKey(r.db, r.name)
		f, kind := b.l.Lookup(k)
		if !v.IfExists {
			b.needKind(k, []string{relView}, r.rng, notFound("view", r.written, ""), r.written, "view") // 4092
		}
		switch {
		case kind == relView:
			b.Set(k, walk.Absent, "")
		case f == walk.Present && kind == "":
			b.Set(k, walk.Unknown, "")
		}
	}
	return true
}

// createTrigger follows CREATE [OR REPLACE] TRIGGER [IF NOT EXISTS].
// Trigger names are unique in the database and compare with case on
// MariaDB, as on MySQL 5.7, so the name is claimed taken only by the
// spelling the walk has. OR REPLACE replaces the trigger of the name; of
// another table's it fails (1360), which the walk does not claim.
// createTrigger follows CREATE TRIGGER. A trigger lives with its table: its database is the one
// its name is qualified with, or the current one, and an unqualified table is looked up there, so
// with shop current CREATE TRIGGER q.tr ... ON t reads q.t. A table in another database answers
// 1435 and the statement changes nothing (mysql:8.0.46, mariadb:11.8.8).
func (b *builder) createTrigger(v *ast.CreateTriggerStmt) bool {
	table := v.Table
	if v.Schema != "" && table != nil && table.Schema == "" {
		qualified := *table
		qualified.Schema = v.Schema
		table = &qualified
	}
	r, ok, lowered := b.tableOf(table)
	if !ok {
		return lowered
	}
	rng := b.anchor()
	if b.tooLong(v.Name, rng) {
		return true
	}
	k, triggerDB, st := b.dbObject(spaceTrigger, &ast.TableRef{Schema: v.Schema, Name: v.Name})
	if st != refOK {
		return st == refFailed
	}
	if triggerDB != r.db {
		return true
	}
	// MariaDB reads the body's NEW and OLD columns and the ordering clause
	// before it looks for a trigger of the name, IF NOT EXISTS or not
	// (verified on mariadb:11.8.8).
	b.triggerChecks(v, r)
	tk := objKey(spaceTrigger, r.db, k)
	f, spelling := b.l.Lookup(tk)
	if f == walk.Present && spelling == v.Name && !v.IfNotExists && !v.OrReplace {
		b.need(tk, walk.Absent, nil, rng, "", alreadyExists("trigger", ident(v.Name), "")) // 1359
	}
	switch {
	case f == walk.Present && spelling == v.Name:
		// OR REPLACE replaces it.
	case f == walk.Present:
		// A case variant, or a spelling the walk does not have: MariaDB
		// keeps both under the one key (verified).
		b.Set(tk, walk.Unknown, "")
	case f == walk.Unknown:
		b.Set(tk, walk.Present, "")
	default:
		b.Set(tk, walk.Present, v.Name)
	}
	rel := relKey(r.db, r.name)
	b.Session(func() { b.w.triggered[rel] = true; b.w.triggersMade[tk] = true })
	return true
}

// triggerChecks checks what CREATE TRIGGER reads first: its body's NEW
// and OLD columns, then the trigger FOLLOWS or PRECEDES names.
func (b *builder) triggerChecks(v *ast.CreateTriggerStmt, r tableRef) {
	b.triggerRowColumns(v.Body, r)
	o := v.Order
	if o == nil {
		return
	}
	// FOLLOWS and PRECEDES name a trigger of the same table, time and
	// event (4031, verified on mariadb:11.8.8); the walk claims only one
	// that doesn't exist.
	k, ok := b.key(spaceTrigger, o.TriggerName)
	if !ok {
		return
	}
	clause := "PRECEDES"
	if o.Follows {
		clause = "FOLLOWS"
	}
	ordered := objKey(spaceTrigger, r.db, k)
	b.need(ordered, walk.Present, nil, rangeOf(o.Loc),
		review.Code(clause)+" names trigger "+review.Code(ident(o.TriggerName))+", which doesn't exist.", "") // 4031
	// One the sync lists on the table, which the change has not dropped or
	// made since, answers 4031 too when its time or event is another.
	st := b.syncedOf(r)
	if st == nil || b.w.triggersMade[ordered] {
		return
	}
	when := strings.ToUpper(v.Timing + " " + v.Event)
	if tr := st.triggers[k]; tr.name == o.TriggerName && tr.when != "" && tr.when != when {
		if f, spelling := b.l.Lookup(ordered); f == walk.Present && spelling == o.TriggerName {
			b.fail(rangeOf(o.Loc), []walk.Key{ordered}, review.ObjectNotFound,
				review.Code(clause)+" names trigger "+review.Code(ident(o.TriggerName))+", which is a "+review.Code(tr.when)+" trigger, not "+review.Code(when)+".") // 4031
		}
	}
}

// triggerRowColumns checks that each NEW or OLD column a trigger's body
// names is a column of its table (1054, verified on mariadb:11.8.8). It
// leaves the ones in a query alone, where a table may be aliased NEW or
// OLD. The parser refuses a row the trigger's event has none of, and an
// update of OLD or of NEW after the event.
func (b *builder) triggerRowColumns(body ast.Node, r tableRef) {
	if body == nil {
		return
	}
	ast.Inspect(body, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.SelectStmt:
			return false
		case *ast.ColumnRef:
			row := strings.ToUpper(x.Table)
			if x.Schema != "" || x.Star || row != "NEW" && row != "OLD" {
				return false
			}
			if k, ok := b.key(spaceColumn, x.Column); ok {
				b.need(colKey(r.db, r.name, k), walk.Present, nil, rangeOf(x.Loc), notFound("column", ident(x.Column), row), "") // 1054
			}
			return false
		}
		return true
	})
}

// dropTrigger follows DROP TRIGGER, which drops the trigger of the
// spelling: the key is Absent after it only when it held that one
// spelling.
func (b *builder) dropTrigger(v *ast.DropTriggerStmt) bool {
	k, db, st := b.dbObject(spaceTrigger, v.Name)
	if st != refOK {
		return st == refFailed
	}
	tk := objKey(spaceTrigger, db, k)
	if !v.IfExists {
		b.need(tk, walk.Present, nil, b.anchor(), notFound("trigger", table(v.Name), ""), "") // 1360
	}
	b.Session(func() { b.w.triggersMade[tk] = true })
	if f, spelling := b.l.Lookup(tk); f == walk.Present && spelling == v.Name.Name || f == walk.Absent {
		b.Set(tk, walk.Absent, "")
	} else {
		b.Set(tk, walk.Unknown, "")
	}
	return true
}

// dbObject resolves the name of a trigger, routine, or event, which a
// database may qualify. With no database selected, an unqualified name
// fails, IF EXISTS or not (1046, verified), but for DROP FUNCTION
// (dropRoutine).
func (b *builder) dbObject(space string, tr *ast.TableRef) (string, string, refState) {
	if tr == nil {
		return "", "", refUnknown
	}
	k, ok := b.key(space, tr.Name)
	if !ok {
		return "", "", refUnknown
	}
	if tr.Schema != "" {
		db, ok := b.key(spaceDatabase, tr.Schema)
		if !ok {
			return "", "", refUnknown
		}
		return k, db, refOK
	}
	switch {
	case b.w.lost:
		return "", "", refUnknown
	case b.w.noDatabase:
		b.fail(b.anchor(), []walk.Key{dbKey(b.w.current)}, review.NoCurrentSchema, noDatabase(table(tr)))
		return "", "", refFailed
	}
	return k, b.w.current, refOK
}

// createdObject resolves the name of a routine or event a statement
// creates, or an event is renamed to: a database it names must exist
// (1049, verified).
func (b *builder) createdObject(space string, tr *ast.TableRef) (string, string, refState) {
	k, db, st := b.dbObject(space, tr)
	if st == refOK && tr.Schema != "" {
		b.databaseNeeds(db, tr.Schema, b.anchor())
	}
	return k, db, st
}

func routineSpace(procedure bool) (string, string) {
	if procedure {
		return spaceProcedure, "procedure"
	}
	return spaceFunction, "function"
}

func (b *builder) createRoutine(v *ast.CreateFunctionStmt) bool {
	if v.Soname != "" {
		// A loadable function lives in no database.
		return true
	}
	space, word := routineSpace(v.IsProcedure)
	k, db, st := b.createdObject(space, v.Name)
	if st != refOK {
		return st == refFailed
	}
	if b.tooLong(v.Name.Name, b.anchor()) {
		return true
	}
	rk := objKey(space, db, k)
	if !v.IfNotExists && !v.OrReplace {
		b.need(rk, walk.Absent, nil, b.anchor(), "", alreadyExists(word, table(v.Name), "")) // 1304
	}
	b.Set(rk, walk.Present, "")
	return true
}

func (b *builder) alterRoutine(v *ast.AlterRoutineStmt) bool {
	space, word := routineSpace(v.IsProcedure)
	if k, db, st := b.dbObject(space, v.Name); st == refOK {
		b.need(objKey(space, db, k), walk.Present, nil, b.anchor(), notFound(word, table(v.Name), ""), "") // 1305
	}
	return true
}

// dropRoutine follows DROP PROCEDURE and DROP FUNCTION. An unqualified
// DROP FUNCTION drops a loadable function of the name first, which the
// sync does not list: it succeeds with no stored function, and leaves the
// stored one when both exist (verified on MySQL 8.0 and 8.4). Kept for
// MariaDB, whose image has no loadable function to verify it with: it
// only gives claims up.
//
// With no database selected, an unqualified DROP FUNCTION looks for a
// loadable function alone: IF EXISTS makes a missing one a note, and
// without it the statement fails with 1305, not 1046, only where none is
// loaded, which the walk cannot know (verified on 11.8.9).
func (b *builder) dropRoutine(v *ast.DropRoutineStmt) bool {
	space, word := routineSpace(v.IsProcedure)
	if !v.IsProcedure && v.Name != nil && v.Name.Schema == "" && b.w.noDatabase && !b.w.lost {
		return true
	}
	k, db, st := b.dbObject(space, v.Name)
	if st != refOK {
		return st == refFailed
	}
	rk := objKey(space, db, k)
	if !v.IsProcedure && v.Name.Schema == "" {
		if f, _ := b.l.Lookup(rk); f == walk.Present {
			b.Set(rk, walk.Unknown, "")
		}
		return true
	}
	if !v.IfExists {
		b.need(rk, walk.Present, nil, b.anchor(), notFound(word, table(v.Name), ""), "") // 1305
	}
	b.Set(rk, walk.Absent, "")
	return true
}

// event resolves an event name, with the database it was qualified with or the current one: with
// shop current, DROP EVENT q.e drops q's, and ALTER EVENT q.e RENAME TO f makes shop.f
// (mysql:8.0.46, mariadb:11.8.8).
func (b *builder) event(schema, name string) (walk.Key, refState) {
	k, db, st := b.dbObject(spaceEvent, &ast.TableRef{Schema: schema, Name: name})
	return objKey(spaceEvent, db, k), st
}

func (b *builder) createEvent(v *ast.CreateEventStmt) bool {
	ek, st := b.event(v.Schema, v.Name)
	if st != refOK {
		return st == refFailed
	}
	if b.eventScheduleFails(v.Schedule) {
		return true
	}
	if ref := eventScheduleColumnRef(v.Schedule); ref != nil {
		b.fail(rangeOf(ref.Loc), nil, review.ObjectNotFound, notFound("column", ident(ref.Column), "")) // 1054
		return true
	}
	if !v.IfNotExists && !v.OrReplace {
		b.need(ek, walk.Absent, nil, b.anchor(), "", alreadyExists("event", ident(v.Name), "")) // 1537
	}
	b.Set(ek, walk.Present, "")
	return true
}

func (b *builder) alterEvent(v *ast.AlterEventStmt) bool {
	ek, st := b.event(v.Schema, v.Name)
	if st != refOK {
		return st == refFailed
	}
	if b.eventScheduleFails(v.Schedule) {
		return true
	}
	if ref := eventScheduleColumnRef(v.Schedule); ref != nil {
		b.fail(rangeOf(ref.Loc), nil, review.ObjectNotFound, notFound("column", ident(ref.Column), "")) // 1054
		return true
	}
	if v.RenameTo != "" {
		nk, st := b.event(v.RenameSchema, v.RenameTo)
		if st != refOK {
			return st == refFailed
		}
		if nk == ek {
			b.fail(b.anchor(), nil, review.DuplicateName, "Event "+review.Code(ident(v.Name))+" cannot be renamed to itself.") // 1551
			return true
		}
		b.need(ek, walk.Present, nil, b.anchor(), notFound("event", ident(v.Name), ""), "")         // 1539
		b.need(nk, walk.Absent, nil, b.anchor(), "", alreadyExists("event", ident(v.RenameTo), "")) // 1537
		b.Set(ek, walk.Absent, "")
		b.Set(nk, walk.Present, "")
		return true
	}
	b.need(ek, walk.Present, nil, b.anchor(), notFound("event", ident(v.Name), ""), "") // 1539
	return true
}

func (b *builder) dropEvent(v *ast.DropEventStmt) bool {
	ek, st := b.event(v.Schema, v.Name)
	if st != refOK {
		return st == refFailed
	}
	if !v.IfExists {
		b.need(ek, walk.Present, nil, b.anchor(), notFound("event", ident(v.Name), ""), "") // 1539
	}
	b.Set(ek, walk.Absent, "")
	return true
}

// singleTarget returns an UPDATE's or DELETE's target when it names one
// table, nil for a join, whose targets the walk does not follow.
func singleTarget(exprs []ast.TableExpr) []*ast.TableRef {
	if len(exprs) == 1 {
		if t, ok := exprs[0].(*ast.TableRef); ok {
			return []*ast.TableRef{t}
		}
	}
	return nil
}

// dml follows INSERT, REPLACE, UPDATE, DELETE, and LOAD DATA: each target
// must exist (1146), as a table, a temporary table, or a view. A trigger
// it fires, or a stored function it calls, may create or drop temporary
// tables; nothing else (spec §6).
//
// A table the walk does not know, a view's base tables, and the tables of
// a join may be in another database, whose triggers the walk does not
// know either.
func (b *builder) dml(n ast.Node, targets []*ast.TableRef) bool {
	fires := targets == nil
	for _, t := range targets {
		r, st := b.ref(t)
		switch st {
		case refFailed:
			return true
		case refUnknown:
			fires = true
			continue
		}
		f, temp, kind := b.resolve(r)
		switch {
		case temp:
		case f == walk.Absent:
			b.need(relKey(r.db, r.name), walk.Present, nil, r.rng, notFound("", r.written, ""), "") // 1146
		case f == walk.Unknown, kind != relTable:
			fires = true
		default:
			fires = fires || b.w.anyTrigger || b.w.triggered[relKey(r.db, r.name)]
		}
	}
	if fires || b.callsStoredFunction(n) {
		b.forgetTemporary()
	}
	return true
}

// set follows SET: foreign_key_checks of the session changes whether a
// foreign key's target is checked. A global setting reaches only later
// sessions. Each assignment has a scope of its own, as on MySQL
// (verified on 11.8): @@session., @@local., and a bare @@ name the
// session's, @@global. not, and an unprefixed name the statement's
// GLOBAL, past any prefixed assignment between. So SET GLOBAL sql_mode =
// DEFAULT, @@session.foreign_key_checks = 0 turns the session's checks
// off, and SET GLOBAL sql_mode = DEFAULT, foreign_key_checks = 0 leaves
// them as they are.
func (b *builder) set(v *ast.SetStmt) bool {
	if b.callsStoredFunction(v) {
		b.forgetTemporary()
	}
	global := false
	for _, a := range v.Assignments {
		if a.Scope != "" {
			global = strings.EqualFold(a.Scope, "GLOBAL")
		}
		if a.Column == nil {
			continue
		}
		name := strings.ToLower(a.Column.Column)
		switch {
		case strings.HasPrefix(name, "@@global."), strings.HasPrefix(name, "@@persist"):
			continue
		case strings.HasPrefix(name, "@@session."), strings.HasPrefix(name, "@@local."):
			name = name[strings.IndexByte(name, '.')+1:]
		case strings.HasPrefix(name, "@@"):
			name = name[2:]
		case global:
			continue
		}
		if name == "sql_mode" && !b.sqlModeOption {
			b.Session(func() { b.w.sqlModeSet = true })
			continue
		}
		if name == "innodb_strict_mode" && !b.strictOption {
			b.Session(func() { b.w.strictSet = true })
			continue
		}
		// SET STATEMENT puts a variable it names back after the statement, whatever the
		// statement set it to (measured on 11.8).
		if name != "foreign_key_checks" || b.fkOption {
			continue
		}
		on := fkSwitch(a.Value)
		b.Session(func() { b.w.fkChecks = on })
	}
	return true
}

// fkSwitch reads the value foreign_key_checks is set to.
func fkSwitch(value ast.ExprNode) review.Switch {
	switch x := value.(type) {
	case *ast.IntLit:
		if x.Value != 0 {
			return review.SwitchOn
		}
		return review.SwitchOff
	case *ast.BoolLit:
		if x.Value {
			return review.SwitchOn
		}
		return review.SwitchOff
	}
	return review.SwitchUnknown
}

// fkOption reads the foreign_key_checks a SET STATEMENT sets for its statement, the last it
// names, and whether it names one.
func fkOption(options []*ast.Assignment) (review.Switch, bool) {
	on, named := review.SwitchUnknown, false
	for _, a := range options {
		if a.Column != nil && a.Column.Table == "" && strings.EqualFold(a.Column.Column, "foreign_key_checks") {
			on, named = fkSwitch(a.Value), true
		}
	}
	return on, named
}
