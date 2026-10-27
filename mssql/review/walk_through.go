package review

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// checkWalkThrough reports the first statement of the change that
// certainly fails on the target, given that the statements before it
// succeeded. The walk follows the change on review/walk's ledger, seeded
// with the synced schema as mssqlSnapshot trusts it; the scope is
// docs/specs/2026-10-26-walk-through-mssql.md.
func checkWalkThrough(ctx context.Context, sql string, stmts []statement, t review.Target, r *reporter) error {
	v, err := walkThrough(ctx, sql, stmts, t)
	if err != nil || v == nil {
		return err
	}
	note, ok := v.Check.Note.(*walkNote)
	if !ok {
		return nil
	}
	r.report(review.WalkThrough, stmts[v.Statement].index, note.rng, escapeLines(note.message(v)))
	return nil
}

// walkThrough walks the change on one target and returns the first
// violation, or nil.
func walkThrough(ctx context.Context, sql string, stmts []statement, t review.Target) (*walk.Violation, error) {
	if t.Schema == nil {
		return nil, nil
	}
	w := newWalker(sql, stmts, t)
	var snap walk.Snapshot = w.snap
	if t.OtherSpecs {
		// Another spec of the plan may have changed anything first, and
		// added triggers the sync does not show.
		snap = nil
		w.triggers = true
	}
	return walk.Run(ctx, len(stmts), snap, w.lower)
}

// walkNote is a check's finding: where it anchors and what it says when
// the object is certainly absent or certainly present (or of another
// kind).
type walkNote struct {
	rng     review.Range
	absent  string
	present string
}

func (n *walkNote) message(v *walk.Violation) string {
	if v.Got == walk.Absent {
		return n.absent
	}
	return n.present
}

// dbState is what the walk knows of the database a name means.
type dbState uint8

const (
	dbTarget dbState = iota
	// dbOther: certainly another database, which the walk does not keep.
	dbOther
	// dbUnknown: maybe the target, maybe not.
	dbUnknown
)

// walker is the session the change runs in, as far as the walk follows
// it.
type walker struct {
	sql    string
	stmts  []statement
	snap   *mssqlSnapshot
	target review.Target
	// db is the current database.
	db dbState
	// triggers is set when a DML trigger may fire on DML: the sync lists
	// one, or the change created or enabled one. A trigger can run DDL,
	// and a cascaded action fires the triggers of other tables (verified).
	triggers bool
	batch    []batchPlace
	// made lists, by schema path, the exact keys of the objects the change
	// made there, which DROP SCHEMA asks the ledger about.
	made map[string][]walk.Key
}

// batchPlace is where a statement sits in its batch, as Bytebase splits
// the change at GO lines (spec §1).
type batchPlace struct {
	// first: the batch's first statement.
	first bool
	// body: part of the body of a procedure or trigger before it in the
	// batch, which runs to the end of the batch, so the statement does
	// not run.
	body bool
	// afterOnly: after a view, function, or CREATE SCHEMA in its batch,
	// which then fails (156).
	afterOnly bool
}

func newWalker(sql string, stmts []statement, t review.Target) *walker {
	w := &walker{
		sql:    sql,
		stmts:  stmts,
		snap:   newMSSQLSnapshot(t.Schema),
		target: t,
		batch:  make([]batchPlace, len(stmts)),
		made:   make(map[string][]walk.Key),
	}
	w.triggers = w.snap.triggers
	first, body, only := true, false, false
	for i := range stmts {
		if _, isGo := stmts[i].node.(*ast.GoStmt); isGo {
			first, body, only = true, false, false
			continue
		}
		w.batch[i] = batchPlace{first: first, body: body, afterOnly: only && !body}
		first = false
		if body || only {
			continue
		}
		switch moduleOf(stmts[i].node) {
		case moduleBody:
			body = true
		case moduleOnly:
			only = true
		}
	}
	return w
}

// module classes a statement by what the rest of its batch is to it.
type module uint8

const (
	moduleNone module = iota
	// moduleBody: a procedure or trigger, whose body runs to the end of
	// the batch.
	moduleBody
	// moduleOnly: a view, function, or schema, after which the batch may
	// hold nothing (156).
	moduleOnly
)

func moduleOf(n ast.Node) module {
	switch n.(type) {
	case *ast.CreateProcedureStmt, *ast.CreateTriggerStmt:
		return moduleBody
	case *ast.CreateViewStmt, *ast.CreateFunctionStmt, *ast.CreateSchemaStmt:
		return moduleOnly
	}
	return moduleNone
}

// lower implements walk.Lower.
func (w *walker) lower(i int, l *walk.Ledger) *walk.Step {
	place := w.batch[i]
	switch {
	case place.body:
		// The statement is part of a module's body and does not run.
		return &walk.Step{}
	case place.afterOnly:
		return nil
	}
	b := &builder{w: w, l: l, s: &w.stmts[i], staged: make(map[walk.Key]staged)}
	if !b.statement(w.stmts[i].node) {
		return nil
	}
	return b.Step()
}

// builder lowers one statement into a walk.Step.
type builder struct {
	walk.Builder
	w *walker
	l *walk.Ledger
	s *statement
	// optional is set while an IF's branches are lowered: each of their
	// parts may not run. weak is set for a part the statement may skip.
	optional bool
	weak     bool
	// staged holds the class facts the statement's earlier parts set, for
	// the effects of its later parts.
	staged map[walk.Key]staged
}

type staged struct {
	fact walk.Fact
	kind string
}

// Next starts the statement's next part, optional within an IF.
func (b *builder) Next() {
	b.Builder.Next()
	b.weak = b.optional
	if b.optional {
		b.Builder.Optional()
	}
}

// maybe marks the current part as one the statement may skip.
func (b *builder) maybe() {
	b.Builder.Optional()
	b.weak = true
}

func (b *builder) anchor() review.Range { return rangeOf(b.s.loc) }

// fail records that the statement certainly fails, from its text.
func (b *builder) fail(rng review.Range, message string) {
	b.Fail(nil, walk.Absent, &walkNote{rng: rng, absent: message, present: message})
}

// sight is what the ledger shows of one spelling of a name: its fact,
// the kind of a Present object, and the key that shows it.
type sight struct {
	fact walk.Fact
	kind string
	key  walk.Key
}

// see reads one spelling of a name: its exact key, else its class key,
// which is Absent when the class holds nothing and the space no marked
// name, and shows the kind of the class's one object when that is spelled
// so.
func (b *builder) see(space, parent string, n name) sight {
	ek := n.exactKey(space, parent)
	switch f, kind := b.l.Lookup(ek); f {
	case walk.Present, walk.Absent:
		return sight{f, kind, ek}
	}
	ck := n.classKey(space, parent)
	switch f, kind := b.l.Lookup(ck); f {
	case walk.Absent:
		// An empty class leaves the spelling missing unless a marked name
		// of the space may equal it.
		if mf, _ := b.l.Lookup(markedKey(space, parent)); mf == walk.Absent {
			return sight{walk.Absent, "", ck}
		}
	case walk.Present:
		if k, sp, ok := single(kind); ok && sp == n.spelling {
			return sight{walk.Present, k, ck}
		}
	}
	return sight{}
}

// needs adds the check that the name exists, and is of one of kinds when
// kinds is not nil. It adds only a check the ledger decides as see does.
func (b *builder) needs(space, parent string, n name, kinds []string, note *walkNote) {
	s := b.see(space, parent, n)
	switch s.fact {
	case walk.Absent:
		b.Need(s.key, walk.Present, nil, note)
	case walk.Present:
		if kinds == nil || s.kind == "" || contains(kinds, s.kind) {
			return
		}
		if isClassKey(s.key) {
			coded := make([]string, len(kinds))
			for i, k := range kinds {
				coded[i] = classKind(k, n.spelling)
			}
			kinds = coded
		}
		b.Need(s.key, walk.Present, kinds, note)
	}
}

// lacks adds the check that the name, spelled so, does not exist.
func (b *builder) lacks(space, parent string, n name, note *walkNote) {
	if s := b.see(space, parent, n); s.fact == walk.Present {
		b.Need(s.key, walk.Absent, nil, note)
	}
}

func contains(xs []string, x string) bool {
	for _, y := range xs {
		if y == x {
			return true
		}
	}
	return false
}

// classFact reads a class key as the statement's earlier parts left it.
func (b *builder) classFact(k walk.Key) (walk.Fact, string) {
	if s, ok := b.staged[k]; ok {
		return s.fact, s.kind
	}
	return b.l.Lookup(k)
}

func (b *builder) setClass(k walk.Key, f walk.Fact, kind string) {
	b.Set(k, f, kind)
	if b.weak {
		f, kind = walk.Unknown, ""
	}
	b.staged[k] = staged{f, kind}
}

// arrive records that an object of a kind arrives under a name: its
// spelling exists, the class holds it, and every other spelling of the
// class may now name it.
func (b *builder) arrive(space, parent string, n name, kind string) {
	ck := n.classKey(space, parent)
	next := ""
	switch f, k := b.classFact(ck); f {
	case walk.Absent:
		next = classKind(kind, n.spelling)
	case walk.Present:
		if _, sp, ok := single(k); ok && sp == n.spelling {
			next = classKind(kind, n.spelling)
		}
	}
	b.MayAdd(walk.Scope{Space: space, Parent: parent, Prefix: n.class})
	b.setClass(ck, walk.Present, next)
	b.Set(n.exactKey(space, parent), walk.Present, kind)
	b.across(space, parent, n, true)
}

// leave records that the object a name means is gone: nothing the
// collation equates to the spelling exists, the class is empty when it
// held only that object, and the class's other spellings may be gone too.
func (b *builder) leave(space, parent string, n name) {
	ck := n.classKey(space, parent)
	next := walk.Unknown
	switch f, k := b.classFact(ck); f {
	case walk.Absent:
		next = walk.Absent
	case walk.Present:
		// The class's one object is gone, unless the name meant a marked
		// one instead, which may still leave it.
		if _, sp, ok := single(k); ok && (sp == n.spelling || b.unmarked(space, parent)) {
			next = walk.Absent
		}
	}
	b.MayRemove(walk.Scope{Space: space, Parent: parent, Prefix: n.class})
	b.setClass(ck, next, "")
	b.Set(n.exactKey(space, parent), walk.Absent, "")
	b.across(space, parent, n, false)
}

// across records what a name's arrival or departure may do to the names
// of other classes: a plain name may equal a marked one, and a marked one
// any (spec §2.1).
func (b *builder) across(space, parent string, n name, add bool) {
	s := walk.Scope{Space: space, Parent: parent}
	if n.plain {
		s.Prefix = markedPrefix
	}
	if add {
		b.MayAdd(s)
	} else {
		b.MayRemove(s)
	}
	if add && !n.plain {
		b.setClass(markedKey(space, parent), walk.Present, "")
	}
}

// unmarked reports whether a space under a parent certainly holds no
// marked name.
func (b *builder) unmarked(space, parent string) bool {
	f, _ := b.classFact(markedKey(space, parent))
	return f == walk.Absent
}

// forgetName records that the statement may have done anything to a name
// of a space in a schema the walk does not know.
func (b *builder) forgetName(space string, n name) {
	s := walk.Scope{Space: space, Deep: true}
	if n.plain {
		s.Prefix = n.class
	}
	b.Forget(s)
}

// place is where a reference resolves.
type place uint8

const (
	// placed: in a known schema of the target, or as a temporary table.
	placed place = iota
	// unplaced: in the target, but the walk does not know where.
	unplaced
	// elsewhere: certainly in another database.
	elsewhere
)

// objRef is a reference to a schema object as the walk resolves it.
type objRef struct {
	n      name
	schema name
	// parent is the object space's parent: Path(schema class), or the
	// root for a temporary table.
	parent  string
	temp    bool
	where   place
	written string
	rng     review.Range
}

var dboName, _ = keyed("dbo")

// ref resolves a reference. An unqualified name is created in the
// default schema and looked up there and then in dbo (verified); with the
// default schema unknown (Target.DefaultSchema empty), it is unplaced. A
// temporary table is a session's, wherever the reference puts it. It
// returns false for a name the walk cannot key, which ends the walk.
func (b *builder) ref(tr *ast.TableRef, creating bool) (objRef, bool) {
	if tr == nil {
		return objRef{}, false
	}
	r := objRef{written: object(tr), rng: rangeOf(tr.Loc)}
	if strings.HasPrefix(tr.Object, "#") {
		n, ok := keyed(strings.TrimPrefix(tr.Object, "#"))
		if !ok {
			return r, false
		}
		r.n, r.temp = n, true
		if strings.HasPrefix(tr.Object, "##") || tr.Server != "" || tr.Database != "" {
			// A global temporary table is every session's.
			r.where = unplaced
		}
		return r, true
	}
	n, ok := keyed(tr.Object)
	if !ok {
		return r, false
	}
	r.n = n
	switch b.database(tr) {
	case dbOther:
		r.where = elsewhere
		return r, true
	case dbUnknown:
		r.where = unplaced
		return r, true
	}
	schema := tr.Schema
	if schema == "" {
		schema = b.defaultSchema(n, creating)
		if schema == "" {
			r.where = unplaced
			return r, true
		}
	}
	sn, ok := keyed(schema)
	if !ok {
		return r, false
	}
	r.schema, r.parent = sn, walk.Path(sn.class)
	if !b.schemaKnown(sn) || tr.Schema == "" && b.see(spaceSchema, walk.Path(), sn).fact != walk.Present {
		// A default schema that may not exist places nothing.
		r.where = unplaced
	}
	return r, true
}

// defaultSchema returns the schema an unqualified name means, "" when the
// walk cannot tell: the default schema for a creation; for a lookup, the
// default schema when it holds an object of the name, and dbo when it
// certainly holds none.
func (b *builder) defaultSchema(n name, creating bool) string {
	d := b.w.target.DefaultSchema
	if d == "" || creating || d == dboName.spelling {
		return d
	}
	dn, ok := keyed(d)
	if !ok {
		return ""
	}
	switch b.see(spaceObject, walk.Path(dn.class), n).fact {
	case walk.Present:
		return d
	case walk.Absent:
		return dboName.spelling
	}
	return ""
}

// schemaKnown reports whether the walk can place objects in a schema: its
// class holds one schema, or none, in which case a creation in it fails
// (2760) and another check says so.
func (b *builder) schemaKnown(sn name) bool {
	f, kind := b.classFact(sn.classKey(spaceSchema, walk.Path()))
	switch f {
	case walk.Absent:
		return true
	case walk.Present:
		_, _, ok := single(kind)
		return ok
	}
	return false
}

// database returns what database a reference means.
func (b *builder) database(tr *ast.TableRef) dbState {
	switch {
	case tr.Server != "":
		// A linked server may be this one.
		return dbUnknown
	case tr.Database == "":
		return b.w.db
	}
	return b.w.databaseNamed(tr.Database)
}

// databaseNamed returns what database a name means: the target only for
// the spelling Bytebase gives it, and certainly another for a name of
// another class.
func (w *walker) databaseNamed(d string) dbState {
	target := w.target.Schema.GetName()
	if d == target {
		return dbTarget
	}
	if differ(d, target) {
		return dbOther
	}
	return dbUnknown
}

// differ reports whether two names certainly differ under every
// collation: both plain, of different classes.
func differ(a, b string) bool {
	an, ok1 := keyed(a)
	bn, ok2 := keyed(b)
	return ok1 && ok2 && an.plain && bn.plain && an.class != bn.class
}

// tablePath is the parent of a table's columns, indexes, constraints, and
// primary key.
func (r objRef) tablePath() string {
	if r.temp {
		return walk.Path(tempSchema, r.n.class)
	}
	return walk.Path(r.schema.class, r.n.class)
}

// relSpace returns the space and parent a reference's relation is under.
func (r objRef) relSpace() (string, string) {
	if r.temp {
		return spaceTemp, walk.Path()
	}
	return spaceRelation, r.parent
}

// knownTable reports whether the walk knows the children of the table a
// reference means: it is placed, and its class holds one object, a
// table. A reference spelled otherwise than that table means it, or
// nothing, and then the statement fails anyway, unless the space holds a
// marked name the reference may mean instead.
func (b *builder) knownTable(r objRef) bool {
	if r.where != placed {
		return false
	}
	space, parent := r.relSpace()
	f, kind := b.classFact(r.n.classKey(space, parent))
	if f != walk.Present {
		return false
	}
	k, sp, ok := single(kind)
	return ok && k == kindTable && (sp == r.n.spelling || b.unmarked(space, parent))
}

// isRelation reports whether a kind is listed in the relation space.
func isRelation(kind string) bool {
	return kind == kindTable || kind == kindView || kind == kindSequence
}

// arriveObject records that an object of a kind arrives under a
// reference's name, in the object space and, for a table, view, or
// sequence, the relation space; an object of another kind leaves no
// relation of its name.
func (b *builder) arriveObject(r objRef, kind string) {
	switch r.where {
	case elsewhere:
		return
	case unplaced:
		b.forgetObject(r)
		return
	}
	if r.temp {
		b.arrive(spaceTemp, walk.Path(), r.n, kind)
		return
	}
	b.arrive(spaceObject, r.parent, r.n, kind)
	b.w.made[r.parent] = append(b.w.made[r.parent], r.n.exactKey(spaceObject, r.parent))
	if isRelation(kind) {
		b.arrive(spaceRelation, r.parent, r.n, kind)
	} else {
		b.Set(r.n.exactKey(spaceRelation, r.parent), walk.Absent, "")
	}
}

// leaveObject records that the object a reference means is gone.
func (b *builder) leaveObject(r objRef, kind string) {
	switch r.where {
	case elsewhere:
		return
	case unplaced:
		b.forgetObject(r)
		return
	}
	if r.temp {
		b.leave(spaceTemp, walk.Path(), r.n)
		return
	}
	b.leave(spaceObject, r.parent, r.n)
	if isRelation(kind) {
		b.leave(spaceRelation, r.parent, r.n)
	}
}

// forgetObject records that an unplaced reference's object may have
// changed in any schema, or among the temporary tables.
func (b *builder) forgetObject(r objRef) {
	if r.temp {
		b.forgetName(spaceTemp, r.n)
		return
	}
	b.forgetName(spaceObject, r.n)
	b.forgetName(spaceRelation, r.n)
}

// childSpaces are the spaces under a table.
var childSpaces = []string{spaceColumn, spaceIndex, spaceConstraint, spacePrimary}

// clearChildren records that a table holds nothing yet, or any more.
func (b *builder) clearChildren(path string) {
	for _, space := range childSpaces {
		b.Clear(walk.Scope{Space: space, Parent: path})
	}
}

// forgetChildren records that a table's children may have changed in any
// way.
func (b *builder) forgetChildren(path string) {
	for _, space := range childSpaces {
		b.Forget(walk.Scope{Space: space, Parent: path})
	}
}

// tooLong fails a statement that writes an identifier longer than SQL
// Server allows, 128 characters or 116 for a temporary table, before
// anything is looked up (103, 193, verified). Only an ASCII name is
// judged.
func (b *builder) tooLong(spelling string, temp bool, rng review.Range) bool {
	limit := 128
	if temp {
		limit = 116
	}
	if len(spelling) <= limit {
		return false
	}
	for i := 0; i < len(spelling); i++ {
		if spelling[i] >= 0x80 {
			return false
		}
	}
	b.fail(rng, "the identifier "+ident(spelling[:min(len(spelling), 32)])+"... is longer than "+strconv.Itoa(limit)+" characters")
	return true
}

// refTooLong checks the parts of a reference.
func (b *builder) refTooLong(tr *ast.TableRef, r objRef) bool {
	if r.temp {
		return b.tooLong(tr.Object, true, r.rng)
	}
	return b.tooLong(tr.Object, false, r.rng) || tr.Schema != "" && b.tooLong(tr.Schema, false, r.rng)
}

// goLine matches a GO line as Bytebase's batcher reads it, as sqlcmd
// does: blanks, GO, and blanks and a count or nothing.
//
// bytebase: backend/plugin/parser/tsql/batch/batch.go — goCommandMatcher
var goLine = regexp.MustCompile(`(?i)^[\t ]*?GO(?:[ ]+(.*)|)$`)

// goCount returns how many times Bytebase sends the batch before a GO,
// or 0 when Bytebase does not take the GO statement's line for one:
// Bytebase splits only at a line holding GO alone, while omni parses GO
// anywhere a statement may start.
func goCount(sql string, loc ast.Loc) int {
	if loc.Start < 0 || loc.Start > len(sql) {
		return 0
	}
	start := strings.LastIndexByte(sql[:loc.Start], '\n') + 1
	end := start
	// The batcher reads a line up to its first control character but a
	// tab.
	for end < len(sql) && (sql[end] >= 0x20 || sql[end] == '\t') && sql[end] != 0x7f {
		end++
	}
	if end < loc.End {
		return 0
	}
	m := goLine.FindStringSubmatch(sql[start:end])
	if m == nil {
		return 0
	}
	if m[1] == "" {
		return 1
	}
	n, err := strconv.Atoi(m[1])
	if err != nil || n < 1 {
		return 0
	}
	return n
}
