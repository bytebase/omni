package review

import (
	"strings"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/review/walk"
)

// drop dispatches a DROP by the kind of object it names.
func (b *builder) drop(v *ast.DropStmt) bool {
	switch v.ObjectType {
	case ast.DropTable:
		return b.dropTables(v)
	case ast.DropIndex:
		return b.dropIndexes(v)
	case ast.DropSchema:
		return b.dropSchema(v)
	case ast.DropType:
		return b.dropTypes(v)
	case ast.DropView:
		return b.dropKind(v, kindView)
	case ast.DropProcedure:
		return b.dropKind(v, kindProcedure)
	case ast.DropFunction:
		return b.dropKind(v, kindFunction)
	case ast.DropSequence:
		return b.dropKind(v, kindSequence)
	case ast.DropSynonym:
		return b.dropKind(v, kindSynonym)
	case ast.DropTrigger:
		if v.OnDatabase || v.OnAllServer {
			return false
		}
		return b.dropKind(v, kindTrigger)
	}
	return false
}

// kindSpace returns the space where the sync lists every object of a
// kind, which then answers for it: the relation space for views and
// sequences, the object space, where absence is the change's alone, for
// the others.
func (r objRef) kindSpace(kind string) string {
	if isRelation(kind) {
		return spaceRelation
	}
	return spaceObject
}

// otherSpace is the space that answers a kind's wrong kind besides
// kindSpace.
func otherSpace(kind string) string {
	if isRelation(kind) {
		return spaceObject
	}
	return spaceRelation
}

// dropKind follows DROP VIEW, PROCEDURE, FUNCTION, SEQUENCE, SYNONYM, or
// TRIGGER [IF EXISTS] x, ...: each x must exist (3701) and be of the
// statement's kind, under IF EXISTS too (3705). The items run in order,
// each judged against the state before the statement, as DROP TABLE's.
// A view takes its triggers along.
func (b *builder) dropKind(v *ast.DropStmt, kind string) bool {
	return b.dropEach(v, spaceObject, func(r objRef) bool {
		if r.temp {
			return false
		}
		if r.where != placed {
			b.leaveObject(r, kind)
			return true
		}
		space := r.kindSpace(kind)
		absent, present := msgMissing(kind, r.written), msgNotKind(r.written, kind)
		if !v.IfExists {
			b.needs(space, r.parent, r.n, []string{kind}, note(r.rng, absent, present))
		} else {
			b.ofKind(space, r.parent, r.n, []string{kind}, note(r.rng, "", present))
			if b.see(space, r.parent, r.n).fact == walk.Unknown {
				b.maybe()
			}
		}
		b.ofKind(otherSpace(kind), r.parent, r.n, []string{kind}, note(r.rng, "", present))
		b.leaveObject(r, kind)
		if kind == kindView {
			b.MayRemove(walk.Scope{Space: spaceObject, Parent: r.parent})
		}
		return true
	})
}

// dropTypes follows DROP TYPE [IF EXISTS]: types are not synced, so only
// the change's own facts say one is missing (218).
func (b *builder) dropTypes(v *ast.DropStmt) bool {
	return b.dropEach(v, spaceType, func(r objRef) bool {
		if r.temp {
			return false
		}
		switch r.where {
		case unplaced:
			b.forgetName(spaceType, r.n)
			return true
		case elsewhere:
			return true
		}
		if !v.IfExists {
			b.needs(spaceType, r.parent, r.n, nil, note(r.rng, msgMissing("type", r.written), ""))
		} else if b.see(spaceType, r.parent, r.n).fact == walk.Unknown {
			b.maybe()
		}
		b.leave(spaceType, r.parent, r.n)
		return true
	})
}

// createType follows CREATE TYPE, alias or table: a type of the name must
// not exist (219), which only the change's own facts tell.
func (b *builder) createType(v *ast.CreateTypeStmt) bool {
	r, ok := b.ref(v.Name, true)
	if !ok || r.temp {
		return false
	}
	switch r.where {
	case elsewhere:
		return true
	case unplaced:
		b.forgetName(spaceType, r.n)
		return true
	}
	if b.refTooLong(v.Name, r) {
		return true
	}
	b.schemaNeeds(r)
	b.lacks(spaceType, r.parent, r.n, note(r.rng, "", "type "+r.written+" already exists"))
	b.arrive(spaceType, r.parent, r.n, "")
	return true
}

// createObject follows CREATE SEQUENCE and CREATE SYNONYM: the schema must
// exist (2760) and the name be free (2714).
func (b *builder) createObject(tr *ast.TableRef, kind string) bool {
	r, ok := b.ref(tr, true)
	if !ok || r.temp {
		return false
	}
	if r.where == elsewhere {
		return true
	}
	if b.refTooLong(tr, r) {
		return true
	}
	b.schemaNeeds(r)
	b.lacksObject(r)
	b.arriveObject(r, kind)
	return true
}

// alterSequence follows ALTER SEQUENCE, of a sequence that must exist
// (15151, for another kind too).
func (b *builder) alterSequence(v *ast.AlterSequenceStmt) bool {
	r, ok := b.ref(v.Name, false)
	if !ok || r.temp {
		return false
	}
	if r.where == placed {
		b.needs(spaceRelation, r.parent, r.n, []string{kindSequence}, note(r.rng, msgMissing("sequence", r.written), msgNotKind(r.written, "sequence")))
	}
	return true
}

// createModule follows CREATE, CREATE OR ALTER, and ALTER of a view,
// procedure, or function, which must begin its batch (111). CREATE needs
// the name free (2714); CREATE OR ALTER an object of the name of the same
// kind if any (2010); ALTER one that exists (208) and is of the kind
// (2010). A procedure's body runs to the end of the batch, which the
// walker skips (spec §6).
func (b *builder) createModule(loc ast.Loc, tr *ast.TableRef, orAlter bool, kind string) bool {
	verb := b.moduleVerb(loc, orAlter)
	if !b.first(loc, verb+" "+strings.ToUpper(kind)) {
		return true
	}
	r, ok := b.ref(tr, verb != "ALTER")
	if !ok || r.temp {
		return false
	}
	if r.where == elsewhere {
		return true
	}
	if b.refTooLong(tr, r) {
		return true
	}
	if r.where == placed {
		b.moduleChecks(r, verb, kind)
	}
	b.arriveObject(r, kind)
	return true
}

// moduleChecks adds a module statement's checks on its name.
func (b *builder) moduleChecks(r objRef, verb, kind string) {
	space := r.kindSpace(kind)
	present := msgNotKind(r.written, kind)
	switch verb {
	case "CREATE":
		b.schemaNeeds(r)
		b.lacksObject(r)
		return
	case "CREATE OR ALTER":
		b.schemaNeeds(r)
		b.ofKind(space, r.parent, r.n, []string{kind}, note(r.rng, "", present))
	default:
		b.needs(space, r.parent, r.n, []string{kind}, note(r.rng, msgMissing(kind, r.written), present))
	}
	b.ofKind(otherSpace(kind), r.parent, r.n, []string{kind}, note(r.rng, "", present))
}

// createTrigger follows CREATE, CREATE OR ALTER, and ALTER of a DML
// trigger: its table must be a table or view (8197), and the trigger in
// the table's schema (2103). A trigger may fire on later DML.
func (b *builder) createTrigger(v *ast.CreateTriggerStmt) bool {
	if v.OnDatabase || v.OnAllServer || v.Table == nil || v.Name == nil {
		return false
	}
	verb := b.moduleVerb(v.Loc, v.OrAlter)
	if !b.first(v.Loc, verb+" TRIGGER") {
		return true
	}
	t, ok := b.ref(v.Table, false)
	if !ok || t.temp {
		return false
	}
	if t.where == elsewhere {
		return true
	}
	if b.refTooLong(v.Table, t) || b.refTooLong(v.Name, objRef{}) {
		return true
	}
	if v.Name.Schema != "" && t.where == placed && differ(v.Name.Schema, t.schema.spelling) {
		b.fail(rangeOf(v.Name.Loc), "trigger "+object(v.Name)+" must be in the schema of its table "+t.written)
		return true
	}
	b.tableNeeds(t, []string{kindTable, kindView}, msgMissing("table", t.written), msgNotKind(t.written, "table or view"))
	r := objRef{schema: t.schema, parent: t.parent, where: t.where, written: object(v.Name), rng: rangeOf(v.Name.Loc)}
	if v.Name.Schema != "" {
		r, ok = b.ref(v.Name, true)
		if !ok {
			return false
		}
	} else if r.n, ok = keyed(v.Name.Object); !ok {
		return false
	}
	if r.where == placed {
		b.moduleChecks(r, verb, kindTrigger)
	}
	b.arriveObject(r, kindTrigger)
	b.Session(func() { b.w.triggers = true })
	return true
}

// spRename follows EXEC sp_rename with literal arguments (spec §4): an
// object, a column, or an index. The new name is taken as written.
func (b *builder) spRename(v *ast.ExecStmt) bool {
	objname, newname, typ, ok := spRename(v)
	if !ok {
		return false
	}
	nn, ok := keyed(newname)
	// sp_rename cuts a new name longer than 128 characters to 128
	// (verified), which the walk does not follow.
	if !ok || len(newname) > 128 {
		return false
	}
	switch strings.ToUpper(typ) {
	case "", "OBJECT":
		return b.renameObject(objname, nn)
	case "COLUMN":
		return b.renameChild(objname, nn, spaceColumn)
	case "INDEX":
		return b.renameChild(objname, nn, spaceIndex)
	}
	return false
}

// renameObject follows sp_rename of a table, view, procedure, function,
// sequence, or synonym: it must exist (15225), and the new name be free
// (15335). A table's columns, indexes, and constraints keep their names
// and move with it. A constraint, a trigger, or an object of a kind the
// walk does not know ends the walk.
func (b *builder) renameObject(objname string, nn name) bool {
	parts := splitParts(objname)
	if len(parts) < 1 || len(parts) > 2 {
		return false
	}
	tr := parseObjectName(objname)
	if tr == nil {
		return false
	}
	r, ok := b.ref(tr, false)
	if !ok || r.temp {
		return false
	}
	switch r.where {
	case elsewhere:
		return true
	case unplaced:
		b.forgetObject(r)
		b.forgetName(spaceObject, nn)
		b.forgetName(spaceRelation, nn)
		return true
	}
	rng := b.anchor()
	r.rng = rng
	os := b.see(spaceObject, r.parent, r.n)
	b.needs(spaceObject, r.parent, r.n, nil, note(rng, msgMissing("object", r.written), ""))
	if os.fact == walk.Absent {
		return true
	}
	kind := os.kind
	if rs := b.see(spaceRelation, r.parent, r.n); rs.fact == walk.Present {
		kind = rs.kind
	}
	switch kind {
	case kindTable, kindView, kindProcedure, kindFunction, kindSequence, kindSynonym:
	default:
		return false
	}
	nr := objRef{n: nn, schema: r.schema, parent: r.parent, where: placed, written: ident(nn.spelling), rng: rng}
	same := nn.class == r.n.class
	if !same {
		b.lacksObject(nr)
	}
	known := kind == kindTable && b.knownTable(r)
	freshTo, _ := b.classFact(nn.classKey(spaceRelation, r.parent))
	from, to := r.tablePath(), nr.tablePath()
	b.leaveObject(r, kind)
	b.arriveObject(nr, kind)
	switch {
	case kind != kindTable, same && known:
	case known && freshTo == walk.Absent:
		b.Move(from, to)
	default:
		b.forgetChildren(from)
		b.forgetChildren(to)
	}
	return true
}

// renameChild follows sp_rename of a column or index of a table: the
// table and the child must exist (15248), and the new name be free among
// the table's (15335). An index that backs a constraint renames the
// constraint too, which the walk does not follow.
func (b *builder) renameChild(objname string, nn name, space string) bool {
	tr, child := splitColumn(objname)
	if tr == nil {
		return false
	}
	cn, ok := keyed(child)
	if !ok {
		return false
	}
	r, ok := b.ref(tr, false)
	if !ok || r.temp {
		return false
	}
	switch r.where {
	case elsewhere:
		return true
	case unplaced:
		b.forgetName(spaceRelation, r.n)
		if space == spaceIndex {
			return false
		}
		return true
	}
	rng := b.anchor()
	what := "column"
	if space == spaceIndex {
		what = "index"
	}
	b.tableNeeds(r, nil, msgMissing("table", r.written), "")
	if !b.knownTable(r) {
		b.forgetMeant(r, space)
		return space != spaceIndex
	}
	path := r.tablePath()
	s := b.see(space, path, cn)
	if space == spaceIndex && s.fact != walk.Absent && (s.fact != walk.Present || s.kind != kindIndex) {
		return false
	}
	in := " in "
	if space == spaceIndex {
		in = " on "
	}
	b.needs(space, path, cn, nil, note(rng, what+" "+ident(child)+" does not exist"+in+r.written, ""))
	if nn.class != cn.class {
		b.lacks(space, path, nn, note(rng, "", what+" "+ident(nn.spelling)+" already exists"+in+r.written))
	}
	kind := ""
	if space == spaceIndex {
		kind = kindIndex
	}
	b.leave(space, path, cn)
	b.arrive(space, path, nn, kind)
	return true
}
