package review

import (
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// tableOf resolves the table an index or trigger statement names, which
// must be a base table. It returns false when the walk does not follow
// the statement on it: a temporary table, or one it cannot tell.
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
		b.fail(rng, nil, "Index name PRIMARY is reserved for the primary key") // 1280
		return true
	}
	ik := idxKey(r.db, r.name, k)
	if v.IfNotExists {
		b.Next()
		b.Optional()
	} else {
		b.need(ik, walk.Absent, nil, rng, "", "Index "+ident(v.IndexName)+" already exists in table "+r.written) // 1061
	}
	b.Set(ik, walk.Present, "")
	return true
}

func (b *builder) dropIndex(v *ast.DropIndexStmt) bool {
	r, ok, lowered := b.tableOf(v.Table)
	if !ok {
		return lowered
	}
	return b.dropIndexName(r, v.Name, b.anchor())
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
	if v.OrReplace {
		if f, _ := b.l.Lookup(k); f == walk.Present {
			b.need(k, walk.Present, []string{relView}, r.rng, "", r.written+" is not a view") // 1347
		}
	} else {
		b.need(k, walk.Absent, nil, r.rng, "", "Table "+r.written+" already exists") // 1050
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
		b.need(relKey(r.db, r.name), walk.Present, []string{relView}, r.rng, "Table "+r.written+" does not exist", r.written+" is not a view") // 1146, 1347
	}
	return true
}

// dropView follows DROP VIEW, which looks every view up first.
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
			b.need(k, walk.Present, []string{relView}, r.rng, "View "+r.written+" does not exist", r.written+" is not a view") // 1051, 1347
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

// createTrigger follows CREATE TRIGGER [IF NOT EXISTS]. Trigger names are
// unique in the database and compare without case on 8.0 but with case on
// 5.7, so the name is claimed taken only by the spelling the walk has.
func (b *builder) createTrigger(v *ast.CreateTriggerStmt) bool {
	r, ok, lowered := b.tableOf(v.Table)
	if !ok {
		return lowered
	}
	rng := b.anchor()
	if b.tooLong(v.Name, rng) {
		return true
	}
	k, ok := b.key(spaceTrigger, v.Name)
	if !ok {
		return false
	}
	tk := objKey(spaceTrigger, r.db, k)
	f, spelling := b.l.Lookup(tk)
	switch {
	case f == walk.Present && spelling == v.Name:
		if !v.IfNotExists {
			b.need(tk, walk.Absent, nil, rng, "", "Trigger "+ident(v.Name)+" already exists") // 1359
		}
	case f == walk.Present:
		// A case variant, or a spelling the walk does not have: 8.0 and
		// 8.4 refuse it (1359), and 5.7 keeps both under the one key
		// (verified).
		b.Set(tk, walk.Unknown, "")
	case f == walk.Unknown:
		b.Set(tk, walk.Present, "")
	default:
		b.Set(tk, walk.Present, v.Name)
	}
	rel := relKey(r.db, r.name)
	b.Session(func() { b.w.triggered[rel] = true })
	return true
}

// dropTrigger follows DROP TRIGGER, which drops the trigger of the
// spelling on 5.7 and of any case variant on 8.0 and 8.4: the key is
// Absent after it only when it held that one spelling.
func (b *builder) dropTrigger(v *ast.DropTriggerStmt) bool {
	k, db, st := b.dbObject(spaceTrigger, v.Name)
	if st != refOK {
		return st == refFailed
	}
	tk := objKey(spaceTrigger, db, k)
	if !v.IfExists {
		b.need(tk, walk.Present, nil, b.anchor(), "Trigger "+ident(v.Name.Name)+" does not exist", "") // 1360
	}
	if f, spelling := b.l.Lookup(tk); f == walk.Present && spelling == v.Name.Name || f == walk.Absent {
		b.Set(tk, walk.Absent, "")
	} else {
		b.Set(tk, walk.Unknown, "")
	}
	return true
}

// dbObject resolves the name of a trigger, routine, or event, which a
// database may qualify. With no database selected, an unqualified name
// fails, IF EXISTS or not (1046, verified).
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
		b.fail(b.anchor(), []walk.Key{dbKey(b.w.current)}, noDatabaseMessage)
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
		return spaceProcedure, "PROCEDURE "
	}
	return spaceFunction, "FUNCTION "
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
		b.need(rk, walk.Absent, nil, b.anchor(), "", word+ident(v.Name.Name)+" already exists") // 1304
	}
	b.Set(rk, walk.Present, "")
	if !v.IsProcedure {
		b.Session(func() { b.w.functions[rk] = true })
	}
	return true
}

func (b *builder) alterRoutine(v *ast.AlterRoutineStmt) bool {
	space, word := routineSpace(v.IsProcedure)
	if k, db, st := b.dbObject(space, v.Name); st == refOK {
		b.need(objKey(space, db, k), walk.Present, nil, b.anchor(), word+ident(v.Name.Name)+" does not exist", "") // 1305
	}
	return true
}

// dropRoutine follows DROP PROCEDURE and DROP FUNCTION. An unqualified
// DROP FUNCTION drops a loadable function of the name first, which the
// sync does not list: it succeeds with no stored function, and leaves the
// stored one when both exist (verified on 8.0 and 8.4).
func (b *builder) dropRoutine(v *ast.DropRoutineStmt) bool {
	space, word := routineSpace(v.IsProcedure)
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
		b.need(rk, walk.Present, nil, b.anchor(), word+ident(v.Name.Name)+" does not exist", "") // 1305
	}
	b.Set(rk, walk.Absent, "")
	return true
}

// event resolves an event name. The parser takes no qualified one, so
// the name is one identifier, a period and all: `shop.e` is an event of
// the current database (verified).
func (b *builder) event(name string) (walk.Key, refState) {
	k, db, st := b.dbObject(spaceEvent, &ast.TableRef{Name: name})
	return objKey(spaceEvent, db, k), st
}

func (b *builder) createEvent(v *ast.CreateEventStmt) bool {
	ek, st := b.event(v.Name)
	if st != refOK {
		return st == refFailed
	}
	if !v.IfNotExists {
		b.need(ek, walk.Absent, nil, b.anchor(), "", "Event "+ident(v.Name)+" already exists") // 1537
	}
	b.Set(ek, walk.Present, "")
	return true
}

func (b *builder) alterEvent(v *ast.AlterEventStmt) bool {
	ek, st := b.event(v.Name)
	if st != refOK {
		return st == refFailed
	}
	b.need(ek, walk.Present, nil, b.anchor(), "Event "+ident(v.Name)+" does not exist", "") // 1539
	if v.RenameTo != "" {
		nk, st := b.event(v.RenameTo)
		if st != refOK {
			return st == refFailed
		}
		if nk != ek {
			b.need(nk, walk.Absent, nil, b.anchor(), "", "Event "+ident(v.RenameTo)+" already exists") // 1537
			b.Set(ek, walk.Absent, "")
		}
		b.Set(nk, walk.Present, "")
	}
	return true
}

func (b *builder) dropEvent(v *ast.DropEventStmt) bool {
	ek, st := b.event(v.Name)
	if st != refOK {
		return st == refFailed
	}
	if !v.IfExists {
		b.need(ek, walk.Present, nil, b.anchor(), "Event "+ident(v.Name)+" does not exist", "") // 1539
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
			b.need(relKey(r.db, r.name), walk.Present, nil, r.rng, "Table "+r.written+" does not exist", "") // 1146
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
// sessions. Each assignment has a scope of its own (verified on 8.0 and
// 8.4): @@session., @@local., and a bare @@ name the session's,
// @@global., @@persist., and @@persist_only. not, and an unprefixed name
// the statement's GLOBAL, PERSIST, or PERSIST_ONLY, past any prefixed
// assignment between. So SET GLOBAL sql_mode = DEFAULT,
// @@session.foreign_key_checks = 0 turns the session's checks off, and
// SET GLOBAL sql_mode = DEFAULT, foreign_key_checks = 0 leaves them as
// they are. omni's parser takes a scope keyword only first.
func (b *builder) set(v *ast.SetStmt) bool {
	if b.callsStoredFunction(v) {
		b.forgetTemporary()
	}
	global := strings.EqualFold(v.Scope, "GLOBAL") || strings.EqualFold(v.Scope, "PERSIST") || strings.EqualFold(v.Scope, "PERSIST_ONLY")
	for _, a := range v.Assignments {
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
		if name != "foreign_key_checks" {
			continue
		}
		on := review.SwitchUnknown
		switch x := a.Value.(type) {
		case *ast.IntLit:
			on = review.SwitchOff
			if x.Value != 0 {
				on = review.SwitchOn
			}
		case *ast.BoolLit:
			on = review.SwitchOff
			if x.Value {
				on = review.SwitchOn
			}
		}
		b.Session(func() { b.w.fkChecks = on })
	}
	return true
}
