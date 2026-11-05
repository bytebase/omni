package review

import (
	"context"
	"strings"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/mariadb/parser"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// checkWalkThrough reports the first statement of the change that
// certainly fails on the target, given that the statements before it
// succeeded. The walk follows the change on review/walk's ledger, seeded
// with the synced schema as mariadbSnapshot trusts it; the scope is
// docs/specs/2026-10-22-walk-through-mysql.md, as
// docs/specs/2026-10-26-walk-through-mariadb-tidb.md §2 changes it for
// MariaDB.
func checkWalkThrough(ctx context.Context, sql string, stmts []statement, t review.Target, r *reporter) error {
	v, err := walkThrough(ctx, sql, stmts, t)
	if err != nil || v == nil {
		return err
	}
	note, ok := v.Check.Note.(*walkNote)
	if !ok {
		return nil
	}
	r.report(note.code(v), stmts[v.Statement].index, note.rng, escapeLines(note.message(v)))
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
		// added triggers and functions the sync does not show.
		snap = nil
		w.anyTrigger = true
		w.anyFunction = true
	}
	return walk.Run(ctx, len(stmts), snap, w.lower)
}

// namesOutsideASCII reports whether a statement's text holds a name with a
// character outside ASCII, which the walk does not follow. MariaDB compares
// such names by rules of its own for each space (verified on 11.8). A
// column, an index, a foreign key, a routine, and an event fold their
// case (CAFé is café), while a table under lower_case_table_names 0 and a
// trigger do not. ſ is s as a column and an index, and ı is i as a
// column, though the Kelvin sign is not k. é is e as a routine and an
// event.
// A token that is no string literal and holds such a byte counts, which
// takes in a little more than names.
func namesOutsideASCII(sql string, loc ast.Loc) bool {
	if loc.Start < 0 || loc.End > len(sql) || loc.Start >= loc.End {
		return false
	}
	text := sql[loc.Start:loc.End]
	if !outsideASCII(text) {
		return false
	}
	for _, tok := range parser.Tokenize(text) {
		t := tokenText(text, tok)
		if t != "" && t[0] != '\'' && t[0] != '"' && outsideASCII(t) {
			return true
		}
	}
	return false
}

func outsideASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 0x80 {
			return true
		}
	}
	return false
}

// walkNote is a check's finding: where it anchors, and what it says and
// the condition it reports when the object is certainly absent or
// certainly present. kinded, when set, writes the text for a relation
// present from the kind the ledger gives it.
type walkNote struct {
	rng         review.Range
	absent      string
	present     string
	kinded      func(kind string) string
	absentCode  review.FindingCode
	presentCode review.FindingCode
}

func (n *walkNote) message(v *walk.Violation) string {
	switch {
	case v.Got == walk.Absent:
		return n.absent
	case n.kinded != nil:
		return n.kinded(v.Kind)
	}
	return n.present
}

func (n *walkNote) code(v *walk.Violation) review.FindingCode {
	if v.Got == walk.Absent {
		return n.absentCode
	}
	return n.presentCode
}

// engineName is the engine a WALK_THROUGH text names.
const engineName = "MariaDB"

// A WALK_THROUGH finding is one sentence, shaped for the finding's code as
// the finding message standard gives it
// (docs/specs/2026-11-02-review-finding-messages.md, Appendix E). A name
// is passed as MariaDB reads it back (ident, table) and written as a code
// span.

// notFound reads "Table `t` doesn't exist.", with an owner "Column `c`
// doesn't exist in `t`.", and with no kind "`t` doesn't exist.".
func notFound(kind, name, owner string) string {
	return named(kind, name) + " doesn't exist" + ownedBy(owner) + "."
}

// alreadyExists reads "Table `t` already exists.", with an owner "Column
// `c` already exists in `t`.".
func alreadyExists(kind, name, owner string) string {
	return named(kind, name) + " already exists" + ownedBy(owner) + "."
}

// nameTaken is the text for a relation name another kind of relation
// holds, or one of a kind the walk does not know.
func nameTaken(name string) string {
	return "The name " + review.Code(name) + " is already taken."
}

// wrongKind reads "`v` is a view, not a table.", and "`v` isn't a table."
// when the kind it is is not known.
func wrongKind(name, kind, need string) string {
	if kind == "" {
		return review.Code(name) + " isn't a " + need + "."
	}
	return review.Code(name) + " is a " + kind + ", not a " + need + "."
}

// statesTwice reads "This statement adds column `c` twice.".
func statesTwice(verb, kind, name string) string {
	return "This statement " + verb + " " + kind + " " + review.Code(name) + " twice."
}

// hasPrimaryKey and noPrimaryKey are the texts for a table that has a
// primary key, or has none, where the statement needs the other.
func hasPrimaryKey(table string) string {
	return review.Code(table) + " already has a primary key."
}

func noPrimaryKey(table string) string {
	return review.Code(table) + " has no primary key."
}

// reservedPrimary is the text for an index named PRIMARY (1280).
var reservedPrimary = review.Code("PRIMARY") + " is reserved for the primary key's index."

// named writes a kind and a name, "Table `t`", or the name alone.
func named(kind, name string) string {
	if kind == "" {
		return review.Code(name)
	}
	return strings.ToUpper(kind[:1]) + kind[1:] + " " + review.Code(name)
}

// ownedBy writes " in `t`" for an owner, nothing for none.
func ownedBy(owner string) string {
	if owner == "" {
		return ""
	}
	return " in " + review.Code(owner)
}

// walker is the session the change runs in, as far as the walk follows
// it.
type walker struct {
	sql        string
	stmts      []statement
	snap       *mariadbSnapshot
	foldTables bool
	// current is the current database's key. noDatabase is set once the
	// change drops it, and lost once a USE names one the walk cannot read.
	current    string
	noDatabase bool
	lost       bool
	fkChecks   review.Switch
	// triggered lists the relations with a trigger, synced or created;
	// anyTrigger is set when any table may have one the walk does not
	// know. anyFunction is set when a call may run a stored function the
	// ledger does not show. A trigger or stored function can create or
	// drop a temporary table, and nothing else (verified: a trigger's DDL
	// fails with 1422).
	triggered   map[walk.Key]bool
	anyTrigger  bool
	anyFunction bool
	// ignoresFKs lists the tables whose engine may ignore foreign keys.
	ignoresFKs map[walk.Key]bool
	// versioned lists the system-versioned tables. Their hidden row_start
	// and row_end act as absent to every check the walk makes but ADD
	// COLUMN (1060), and by default an ALTER TABLE of one fails for its
	// versioning (4119) whatever else the walk would say (verified). So
	// their columns are Unknown, and an ALTER TABLE of one ends the walk
	// (MariaDB spec §2.3), as a precaution.
	versioned map[walk.Key]bool
}

func newWalker(sql string, stmts []statement, t review.Target) *walker {
	fold := t.LowerCaseTableNames != 0
	snap := newMariaDBSnapshot(t.Schema, fold)
	w := &walker{
		sql:        sql,
		stmts:      stmts,
		snap:       snap,
		foldTables: fold,
		current:    snap.db,
		fkChecks:   t.ForeignKeyChecks,
		triggered:  make(map[walk.Key]bool),
		ignoresFKs: make(map[walk.Key]bool),
		versioned:  make(map[walk.Key]bool),
	}
	for name := range snap.triggered {
		w.triggered[relKey(snap.db, name)] = true
	}
	for name := range snap.ignoresFKs {
		w.ignoresFKs[relKey(snap.db, name)] = true
	}
	for name := range snap.versioned {
		w.versioned[relKey(snap.db, name)] = true
	}
	return w
}

// lower implements walk.Lower.
func (w *walker) lower(i int, l *walk.Ledger) *walk.Step {
	if namesOutsideASCII(w.sql, w.stmts[i].loc) {
		return nil
	}
	b := &builder{w: w, l: l, s: &w.stmts[i]}
	// SET STATEMENT refuses an option that calls a stored function or reads a table (1970,
	// measured on 11.8.8), and then neither runs, so what it wraps does not happen. The walk
	// may not tell a stored function from a built-in one, so it follows no such statement.
	for _, opt := range w.stmts[i].options {
		if b.optionMayFail(opt.Value) {
			return nil
		}
	}
	// SET STATEMENT sets foreign_key_checks for this statement alone.
	checks := w.fkChecks
	if on, named := fkOption(w.stmts[i].options); named {
		w.fkChecks, b.fkOption = on, true
	}
	lowered := b.statement(w.stmts[i].node)
	w.fkChecks = checks
	if !lowered {
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
	// fkOption is set while the statement runs under a SET STATEMENT that
	// names foreign_key_checks.
	fkOption bool
	// fks is set when the statement's table keeps the foreign keys it
	// names (keepsFKs). Until a lowering sets it, the statement claims
	// nothing of foreign keys and leaves their names Unknown.
	fks bool
	// renames is set while an ALTER TABLE also renames its table, which
	// leaves the foreign keys it adds no parent the walk can resolve
	// (foreignKeyTarget).
	renames bool
	// early is set while ALTER TABLE lowers its drops and renames, which
	// run together against the table as it was; arrived holds the column
	// and index names they bring in (leave, arrive).
	early   bool
	arrived map[walk.Key]bool
}

// leave records that an ALTER TABLE subcommand takes a column or index
// name away. Among the drops and renames, a name another of them brings
// in stays, whatever their order: a swap of two columns or indexes keeps
// both, and RENAME INDEX i TO j, DROP INDEX j leaves j (verified). It is
// still an effect on the name, so the touched rule leaves out the other
// subcommand's check of it.
func (b *builder) leave(k walk.Key) {
	if b.early && b.arrived[k] {
		b.Set(k, walk.Present, "")
		return
	}
	b.Set(k, walk.Absent, "")
}

// arrive records that an ALTER TABLE subcommand renames a column or index
// to a name.
func (b *builder) arrive(k walk.Key) {
	b.Set(k, walk.Present, "")
	if b.early {
		if b.arrived == nil {
			b.arrived = make(map[walk.Key]bool)
		}
		b.arrived[k] = true
	}
}

func (b *builder) anchor() review.Range { return rangeOf(b.s.loc) }

// need adds a check, with the messages for a key certainly absent and
// certainly present (or of another kind). The check's want gives the
// condition: a key wanted present and absent is an object not found, and
// present of another kind a wrong kind; one wanted absent and present is
// a name taken.
func (b *builder) need(k walk.Key, want walk.Fact, kinds []string, rng review.Range, absent, present string) {
	note := &walkNote{rng: rng, absent: absent, present: present, absentCode: review.ObjectNotFound, presentCode: review.WrongObjectKind}
	if want == walk.Absent {
		note.presentCode = review.NameTaken
	}
	b.Need(k, want, kinds, note)
}

// needKind adds the check that a relation exists as one of kinds: absent
// is the text for a missing one, and one of another kind reads as
// wrongKind of the name and need, the kind wanted.
func (b *builder) needKind(k walk.Key, kinds []string, rng review.Range, absent, name, need string) {
	b.Need(k, walk.Present, kinds, &walkNote{
		rng:         rng,
		absent:      absent,
		kinded:      func(kind string) string { return wrongKind(name, relWord(k, kind), need) },
		absentCode:  review.ObjectNotFound,
		presentCode: review.WrongObjectKind,
	})
}

// needFree adds the check that a relation name is free for a relation of
// kind, "" when that is not known: one of that kind already exists, and
// any other takes the name.
func (b *builder) needFree(k walk.Key, rng review.Range, name, kind string) {
	b.Need(k, walk.Absent, nil, &walkNote{
		rng: rng,
		kinded: func(held string) string {
			if held != "" && held == kind {
				return alreadyExists(relWord(k, kind), name, "")
			}
			return nameTaken(name)
		},
		absentCode:  review.ObjectNotFound,
		presentCode: review.NameTaken,
	})
}

// relWord is the word for a relation of a kind under a key: a temporary
// one's is "temporary table".
func relWord(k walk.Key, kind string) string {
	if k.Space == spaceTemporary && kind != "" {
		return "temporary " + kind
	}
	return kind
}

// fail records that the statement certainly fails, established from the
// keys it read; with no keys, from its text. code is the condition.
func (b *builder) fail(rng review.Range, keys []walk.Key, code review.FindingCode, message string) {
	b.Fail(keys, walk.Absent, &walkNote{rng: rng, absent: message, present: message, absentCode: code, presentCode: code})
}

func dbKey(db string) walk.Key { return walk.Key{Space: spaceDatabase, Name: db} }

func relKey(db, name string) walk.Key {
	return walk.Key{Space: spaceRelation, Parent: walk.Path(db), Name: name}
}

func tempKey(db, name string) walk.Key {
	return walk.Key{Space: spaceTemporary, Parent: walk.Path(db), Name: name}
}

// tablePath is the parent path of a table's columns and indexes.
func tablePath(db, name string) string { return walk.Path(db, name) }

func colKey(db, table, name string) walk.Key {
	return walk.Key{Space: spaceColumn, Parent: tablePath(db, table), Name: name}
}

func idxKey(db, table, name string) walk.Key {
	return walk.Key{Space: spaceIndex, Parent: tablePath(db, table), Name: name}
}

// objKey keys a foreign key, trigger, routine, or event of a database.
func objKey(space, db, name string) walk.Key {
	return walk.Key{Space: space, Parent: walk.Path(db), Name: name}
}

// fkKeys returns a table's foreign key's two keys: its name in the
// database, and on the table.
func fkKeys(r tableRef, name string) (walk.Key, walk.Key) {
	return objKey(spaceForeignKey, r.db, name), walk.Key{Space: spaceTableForeignKey, Parent: tablePath(r.db, r.name), Name: name}
}

// keepsFKs reports whether a table engine keeps foreign keys: InnoDB
// does, and NDB, which MariaDB substitutes or refuses. Another engine
// parses and ignores them: MyISAM and Aria accept a foreign key to a
// missing table and keep none (verified on MariaDB 11.8). An empty engine
// is the server's default, assumed InnoDB (MySQL spec §3).
func keepsFKs(engine string) bool {
	switch strings.ToLower(engine) {
	case "", "innodb", "ndb", "ndbcluster":
		return true
	}
	return false
}

// engineOf returns the ENGINE a statement's table options name, "" for
// none.
func engineOf(opts []*ast.TableOption) string {
	engine := ""
	for _, o := range opts {
		if o != nil && strings.EqualFold(o.Name, "ENGINE") {
			engine = strings.Trim(o.Value, "`'\"")
		}
	}
	return engine
}

// generatedFKs returns the prefix of the foreign key names InnoDB
// generates for a table, <table>_ibfk_<n>; RENAME TABLE renames every
// foreign key whose name begins with it along with the table (verified on
// MySQL 8.0 and 8.4, and MariaDB 11.8). For a name outside ASCII it
// returns "", which takes in every name: such a foreign key's key is its
// exact spelling, which the table's key need not be.
func generatedFKs(table string) string {
	if strings.HasPrefix(table, nonASCII) {
		return ""
	}
	return strings.ToLower(table) + "_ibfk_"
}

// key returns the key of a name a statement writes, false when the walk
// cannot read it. A relation name the SEQUENCE engine answers for is one
// it cannot: such a table is in every database, though the sync lists
// none, and comes back when a table of its name is dropped.
func (b *builder) key(space, name string) (string, bool) {
	k, ok := nameKey(space, name, b.w.foldTables)
	if ok && space == spaceRelation && sequenceEngineTable(k) {
		return "", false
	}
	return k, ok
}

// sequenceEngineTable reports whether a relation key is a name MariaDB's
// SEQUENCE engine, built in and enabled by default, answers for in any
// database: seq_<from>_to_<to>, or that with _step_<step> after
// (storage/sequence/sequence.cc). On 11.8.8 SELECT reads seq_1_to_3 and
// mysql.seq_1_to_3, DROP TABLE and CREATE TRIGGER take it, CREATE TABLE
// finds it (1050), and DML fails, but not for want of a table (1031). The
// engine reads the step with sscanf, which takes a sign or a space before
// it; the walk takes anything after _step_.
func sequenceEngineTable(k string) bool {
	digits := func(s string) (string, bool) {
		i := 0
		for i < len(s) && '0' <= s[i] && s[i] <= '9' {
			i++
		}
		return s[i:], i > 0
	}
	rest, ok := strings.CutPrefix(k, "seq_")
	if !ok {
		return false
	}
	if rest, ok = digits(rest); !ok {
		return false
	}
	if rest, ok = strings.CutPrefix(rest, "_to_"); !ok {
		return false
	}
	if rest, ok = digits(rest); !ok {
		return false
	}
	return rest == "" || strings.HasPrefix(rest, "_step_")
}

// tableRef is a table reference as the walk resolves it: the keys of its
// database and name, and how the SQL wrote it.
type tableRef struct {
	db, name string
	written  string
	rng      review.Range
}

// noDatabase is the text for a name, or a statement that names none, with
// no database selected (1046).
func noDatabase(name string) string {
	return "No database is selected for " + review.Code(name) + "."
}

// refState is how ref resolved a reference.
type refState uint8

const (
	refOK refState = iota
	// refFailed: the statement certainly fails for want of a selected
	// database (1046), which ref recorded.
	refFailed
	// refUnknown: the walk cannot tell what the reference names, an
	// unreadable name or an unqualified one after a USE it cannot follow.
	refUnknown
)

// ref resolves a table reference's database and name.
func (b *builder) ref(tr *ast.TableRef) (tableRef, refState) {
	if tr == nil {
		return tableRef{}, refUnknown
	}
	r := tableRef{written: table(tr), rng: rangeOf(tr.Loc)}
	name, ok := b.key(spaceRelation, tr.Name)
	if !ok {
		return r, refUnknown
	}
	r.name = name
	if tr.Schema != "" {
		db, ok := b.key(spaceDatabase, tr.Schema)
		if !ok {
			return r, refUnknown
		}
		r.db = db
		return r, refOK
	}
	switch {
	case b.w.lost:
		return r, refUnknown
	case b.w.noDatabase:
		b.fail(r.rng, []walk.Key{dbKey(b.w.current)}, review.NoCurrentSchema, noDatabase(r.written))
		return r, refFailed
	}
	r.db = b.w.current
	return r, refOK
}

// resolve returns what a statement that uses an existing table meets
// under the reference: a temporary table or sequence of the name shadows a
// base relation (verified). kind is the relation's kind when Present, and
// a temporary one's when known.
func (b *builder) resolve(r tableRef) (f walk.Fact, temp bool, kind string) {
	switch tf, tkind := b.l.Lookup(tempKey(r.db, r.name)); tf {
	case walk.Present:
		return walk.Present, true, tkind
	case walk.Unknown:
		return walk.Unknown, false, ""
	}
	f, kind = b.l.Lookup(relKey(r.db, r.name))
	return f, false, kind
}

// tooLong fails a statement that names an identifier longer than 64
// characters, before anything is looked up (1059; 1102 for a database and
// 1103 for a table, view, or sequence; verified). Only an ASCII name is
// judged.
func (b *builder) tooLong(name string, rng review.Range) bool {
	if len(name) <= 64 {
		return false
	}
	for i := 0; i < len(name); i++ {
		if name[i] >= 0x80 {
			return false
		}
	}
	b.fail(rng, nil, review.NameTooLong, "The name "+review.Code(ident(name))+" is longer than "+engineName+"'s 64-character limit.")
	return true
}

// optionMayFail reports whether SET STATEMENT may refuse this option's value.
func (b *builder) optionMayFail(value ast.Node) bool {
	readsTable := false
	ast.Inspect(value, func(m ast.Node) bool {
		if _, ok := m.(*ast.TableRef); ok {
			readsTable = true
		}
		return !readsTable
	})
	return readsTable || b.callsStoredFunction(value)
}

// callsStoredFunction reports whether the statement may call a stored
// function, which may create or drop a temporary table (verified on
// 11.8.9): a qualified call, an unqualified one to a function the current
// database may hold, which after a USE of a database the change did not
// create the walk does not know, and a read of a relation that may be a
// view, whose query may call one.
func (b *builder) callsStoredFunction(n ast.Node) bool {
	calls := false
	ast.Inspect(n, func(m ast.Node) bool {
		// Inspect visits the siblings of a node whose children it skips,
		// so a call once found stays found.
		if calls {
			return false
		}
		switch v := m.(type) {
		case *ast.FuncCallExpr:
			calls = b.w.anyFunction || v.Schema != "" || b.w.lost
			if !calls {
				k, ok := b.key(spaceFunction, v.Name)
				f, _ := b.l.Lookup(objKey(spaceFunction, b.w.current, k))
				calls = !ok || f != walk.Absent
			}
		case *ast.TableRef:
			calls = b.mayBeView(v)
		}
		return !calls
	})
	return calls
}

// mayBeView reports whether a table reference may name a view: anything
// but a temporary table, or a base table or sequence the walk knows.
func (b *builder) mayBeView(tr *ast.TableRef) bool {
	name, ok := b.key(spaceRelation, tr.Name)
	if !ok {
		return true
	}
	r := tableRef{name: name, db: b.w.current}
	switch {
	case tr.Schema != "":
		if r.db, ok = b.key(spaceDatabase, tr.Schema); !ok {
			return true
		}
	case b.w.lost || b.w.noDatabase:
		return true
	}
	f, temp, kind := b.resolve(r)
	return !temp && (f != walk.Present || kind != relTable && kind != relSequence)
}

// forgetTemporary records that a trigger or stored function may have
// created or dropped temporary tables of any name.
func (b *builder) forgetTemporary() {
	b.Forget(walk.Scope{Space: spaceTemporary, Deep: true})
}
