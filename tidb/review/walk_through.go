package review

import (
	"context"
	"strings"

	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/parser"
)

// checkWalkThrough reports the first statement of the change that
// certainly fails on the target, given that the statements before it
// succeeded. The walk follows the change on review/walk's ledger, seeded
// with the synced schema as tidbSnapshot trusts it; the scope is
// docs/specs/2026-10-22-walk-through-mysql.md, with the TiDB deltas of
// docs/specs/2026-10-26-walk-through-mariadb-tidb.md §3.
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
		// Another spec of the plan may have changed anything first.
		snap = nil
	}
	return walk.Run(ctx, len(stmts), snap, w.lower)
}

// namesOutsideASCII reports whether a statement's text holds a name with a
// character outside ASCII, which the walk does not follow. TiDB compares
// such names by rules of its own for each space (verified): CAFé is café
// as a column and as a table, ſ is s as a column but not as a table, and
// the Kelvin sign is k as both, while é is not e. A token that is no
// string literal and holds such a byte counts, which takes in a little
// more than names.
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
const engineName = "TiDB"

// A WALK_THROUGH finding is one sentence, shaped for the finding's code as
// the finding message standard gives it
// (docs/specs/2026-11-02-review-finding-messages.md, Appendix E). A name
// is passed as TiDB reads it back (ident, table) and written as a code
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
	snap       *tidbSnapshot
	foldTables bool
	// current is the current database's key. noDatabase is set once the
	// change drops it, and lost once a USE names one the walk cannot read.
	current    string
	noDatabase bool
	lost       bool
	fkChecks   review.Switch
	// TiDB has no triggers or stored functions (1064), so no statement
	// runs code that creates or drops a temporary table: the walk keeps no
	// record of either.
}

func newWalker(sql string, stmts []statement, t review.Target) *walker {
	// TiDB's lower_case_table_names is 2 and fixed, whatever the target
	// says (spec §3.1).
	fold := true
	snap := newTiDBSnapshot(t.Schema, fold)
	return &walker{
		sql:        sql,
		stmts:      stmts,
		snap:       snap,
		foldTables: fold,
		current:    snap.db,
		fkChecks:   t.ForeignKeyChecks,
	}
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
	// names: every base table, since TiDB ignores ENGINE (verified: a
	// foreign key of an ENGINE=MyISAM table to a missing table fails with
	// 1824). Until a lowering sets it, the statement claims nothing of
	// foreign keys and leaves their names Unknown.
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

// fkKey returns the key of a table's foreign key. TiDB keeps foreign key
// names per table (spec §3.4), so a foreign key has no key of the
// database's.
func fkKey(r tableRef, name string) walk.Key {
	return walk.Key{Space: spaceTableForeignKey, Parent: tablePath(r.db, r.name), Name: name}
}

// generatedFKs is the prefix of the foreign key names TiDB generates for a
// table, fk_<n>, which RENAME TABLE leaves as they are (verified).
const generatedFKs = "fk_"

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
	b.fail(rng, nil, review.NameTooLong, "The name "+review.Code(ident(name))+" is longer than "+engineName+"'s 64-character limit.")
	return true
}
