package review

import (
	"context"
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/parser"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// checkWalkThrough reports the first statement of the change that
// certainly fails on the target, given that the statements before it
// succeeded. The walk follows the change on review/walk's ledger, seeded
// with the synced schema as mysqlSnapshot trusts it; the scope is
// docs/specs/2026-10-22-walk-through-mysql.md.
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
		// added triggers and functions the sync does not show.
		snap = nil
		w.anyTrigger = true
		w.anyFunction = true
	}
	return walk.Run(ctx, len(stmts), snap, w.lower)
}

// namesOutsideASCII reports whether a statement's text holds a name with a
// character outside ASCII, which the walk does not follow. MySQL compares
// such names by rules of its own for each space (verified on 8.0 and
// 8.4). A column, an index, a foreign key, a trigger, a routine, and an
// event fold their case (CAFé is café), and so do a table and a database
// under lower_case_table_names 1 but not 0. The Kelvin sign is k and İ is
// i as a column and an index. é is e as a routine, an event, a trigger,
// and a foreign key's name in its database. ſ is not s.
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

// walkNote is a check's finding: where it anchors and what it says when
// the object is certainly absent or certainly present.
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

// walker is the session the change runs in, as far as the walk follows
// it.
type walker struct {
	sql        string
	stmts      []statement
	snap       *mysqlSnapshot
	foldTables bool
	// current is the current database's key. noDatabase is set once the
	// change drops it, and lost once a USE names one the walk cannot read.
	current    string
	noDatabase bool
	lost       bool
	fkChecks   review.Switch
	// triggered lists the relations with a trigger, synced or created;
	// anyTrigger is set when any table may have one the walk does not
	// know. functions lists the stored functions the target or the change
	// defines; anyFunction is set when a call may run one the walk does
	// not know. A trigger or stored function can create or drop a
	// temporary table, and nothing else (verified: a trigger's DDL fails
	// with 1422).
	triggered   map[walk.Key]bool
	anyTrigger  bool
	functions   map[walk.Key]bool
	anyFunction bool
	// ignoresFKs lists the tables whose engine may ignore foreign keys.
	ignoresFKs map[walk.Key]bool
}

func newWalker(sql string, stmts []statement, t review.Target) *walker {
	fold := t.LowerCaseTableNames != 0
	snap := newMySQLSnapshot(t.Schema, fold)
	w := &walker{
		sql:        sql,
		stmts:      stmts,
		snap:       snap,
		foldTables: fold,
		current:    snap.db,
		fkChecks:   t.ForeignKeyChecks,
		triggered:  make(map[walk.Key]bool),
		functions:  make(map[walk.Key]bool),
		ignoresFKs: make(map[walk.Key]bool),
	}
	for name := range snap.triggered {
		w.triggered[relKey(snap.db, name)] = true
	}
	for name := range snap.ignoresFKs {
		w.ignoresFKs[relKey(snap.db, name)] = true
	}
	for name := range snap.functions {
		w.functions[objKey(spaceFunction, snap.db, name)] = true
	}
	return w
}

// lower implements walk.Lower.
func (w *walker) lower(i int, l *walk.Ledger) *walk.Step {
	if namesOutsideASCII(w.sql, w.stmts[i].loc) {
		return nil
	}
	b := &builder{w: w, l: l, s: &w.stmts[i]}
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
	// fks is set when the statement's table keeps the foreign keys it
	// names (keepsFKs). Until a lowering sets it, the statement claims
	// nothing of foreign keys and leaves their names Unknown.
	fks bool
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
// certainly present (or of another kind).
func (b *builder) need(k walk.Key, want walk.Fact, kinds []string, rng review.Range, absent, present string) {
	b.Need(k, want, kinds, &walkNote{rng: rng, absent: absent, present: present})
}

// fail records that the statement certainly fails, established from the
// keys it read; with no keys, from its text.
func (b *builder) fail(rng review.Range, keys []walk.Key, message string) {
	b.Fail(keys, walk.Absent, &walkNote{rng: rng, absent: message, present: message})
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

// keepsFKs reports whether a table engine keeps foreign keys: InnoDB and
// NDB do. Another engine parses and ignores them: it accepts a foreign
// key to a missing table, a name another one has, and DROP FOREIGN KEY of
// any name (verified with MyISAM on 8.0 and 8.4). An empty engine is the
// server's default, assumed InnoDB (spec §3).
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

// generatedFKs returns the prefix of the foreign key names MySQL generates
// for a table, <table>_ibfk_<n>; RENAME TABLE renames every foreign key
// whose name begins with it along with the table (verified on 8.0 and
// 8.4). For a name outside ASCII it returns "", which takes in every name:
// such a foreign key's key is its exact spelling, which the table's key
// need not be.
func generatedFKs(table string) string {
	if strings.HasPrefix(table, nonASCII) {
		return ""
	}
	return strings.ToLower(table) + "_ibfk_"
}

func (b *builder) key(space, name string) (string, bool) {
	return nameKey(space, name, b.w.foldTables)
}

// tableRef is a table reference as the walk resolves it: the keys of its
// database and name, and how the SQL wrote it.
type tableRef struct {
	db, name string
	written  string
	rng      review.Range
}

const noDatabaseMessage = "No database selected"

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
		b.fail(r.rng, []walk.Key{dbKey(b.w.current)}, noDatabaseMessage)
		return r, refFailed
	}
	r.db = b.w.current
	return r, refOK
}

// resolve returns what a statement that uses an existing table meets
// under the reference: a temporary table of the name shadows a base
// relation (verified). kind is the relation's kind when Present.
func (b *builder) resolve(r tableRef) (f walk.Fact, temp bool, kind string) {
	switch tf, _ := b.l.Lookup(tempKey(r.db, r.name)); tf {
	case walk.Present:
		return walk.Present, true, relTable
	case walk.Unknown:
		return walk.Unknown, false, ""
	}
	f, kind = b.l.Lookup(relKey(r.db, r.name))
	return f, false, kind
}

// tooLong fails a statement that names an identifier longer than 64
// characters, before anything is looked up (1059, verified). Only an ASCII
// name is judged.
func (b *builder) tooLong(name string, rng review.Range) bool {
	if len(name) <= 64 {
		return false
	}
	for i := 0; i < len(name); i++ {
		if name[i] >= 0x80 {
			return false
		}
	}
	b.fail(rng, nil, "Identifier name "+ident(name)+" is too long")
	return true
}

// callsStoredFunction reports whether the statement may call a stored
// function, which may create or drop a temporary table.
func (b *builder) callsStoredFunction(n ast.Node) bool {
	calls := false
	ast.Inspect(n, func(m ast.Node) bool {
		fc, ok := m.(*ast.FuncCallExpr)
		if !ok || calls {
			return !calls
		}
		switch {
		case b.w.anyFunction, fc.Schema != "":
			calls = true
		default:
			if k, ok := b.key(spaceFunction, fc.Name); ok && b.w.functions[objKey(spaceFunction, b.w.current, k)] {
				calls = true
			}
		}
		return !calls
	})
	return calls
}

// forgetTemporary records that a trigger or stored function may have
// created or dropped temporary tables of any name.
func (b *builder) forgetTemporary() {
	b.Forget(walk.Scope{Space: spaceTemporary, Deep: true})
}
