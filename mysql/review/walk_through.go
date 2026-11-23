package review

import (
	"context"
	"slices"
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
		w.fkGeneratedIndexes = make(map[walk.Key]map[string]bool)
		w.engines = make(map[walk.Key]string)
		w.colTypes = make(map[walk.Key]string)
		w.partitionState = make(map[walk.Key]*tablePartitionInfo)
		w.tableIndexes = make(map[walk.Key]map[string]*indexMeta)
		w.tableIndexesComplete = make(map[walk.Key]bool)
		w.tableFKs = make(map[walk.Key]map[string][]string)
		w.referencedByFK = make(map[walk.Key]map[fkRefKey]bool)
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
const engineName = "MySQL"

// A WALK_THROUGH finding is one sentence, shaped for the finding's code as
// the finding message standard gives it
// (docs/specs/2026-11-02-review-finding-messages.md, Appendix E). A name
// is passed as MySQL reads it back (ident, table) and written as a code
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

// fkRefKey identifies one foreign key on a child table referencing a parent.
type fkRefKey struct {
	child walk.Key
	fk    string
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
	// defaultEngine and defaultTmpEngine track session default_storage_engine
	// and default_tmp_storage_engine ("" means default/InnoDB); the *Unknown
	// flags are set when a SET assigns a value the walk cannot read.
	defaultEngine           string
	defaultEngineUnknown    bool
	defaultTmpEngine        string
	defaultTmpEngineUnknown bool
	// innodbStrictOff is set when the session turns innodb_strict_mode off or
	// to an unknown value, which downgrades InnoDB ROW_FORMAT=FIXED to a warning.
	innodbStrictOff bool
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
	// fkGeneratedIndexes maps synced index keys that may be MySQL-generated
	// foreign-key backing indexes (HA_GENERATED_KEY) to the set of synced FK
	// keys whose columns match them. Dropping all matching FKs retires the
	// candidate; adding an index consumes it (and marks other candidates on the
	// table Unknown).
	fkGeneratedIndexes map[walk.Key]map[string]bool
	// engines records known table storage engines ("" means default/InnoDB);
	// absence from the map means the engine is unknown.
	engines map[walk.Key]string
	// colTypes records known base data types keyed by colKey(db, tbl, col).
	colTypes map[walk.Key]string
	// partitionState records known partitioning state keyed by relKey(db, tbl).
	partitionState map[walk.Key]*tablePartitionInfo
	// tableIndexes records known index column lists keyed by relKey(db, tbl).
	tableIndexes map[walk.Key]map[string]*indexMeta
	// tableIndexesComplete reports whether tableIndexes[rel] lists every index
	// on the table (false once an unnamed/unknown index may have been added).
	tableIndexesComplete map[walk.Key]bool
	// tableFKs records known foreign key column lists keyed by relKey(db, tbl).
	tableFKs map[walk.Key]map[string][]string
	// referencedByFK maps a parent relKey to the child FKs referencing it.
	referencedByFK map[walk.Key]map[fkRefKey]bool
}

func newWalker(sql string, stmts []statement, t review.Target) *walker {
	fold := t.LowerCaseTableNames != 0
	snap := newMySQLSnapshot(t.Schema, fold)
	w := &walker{
		sql:                  sql,
		stmts:                stmts,
		snap:                 snap,
		foldTables:           fold,
		current:              snap.db,
		fkChecks:             t.ForeignKeyChecks,
		triggered:            make(map[walk.Key]bool),
		ignoresFKs:           make(map[walk.Key]bool),
		fkGeneratedIndexes:   make(map[walk.Key]map[string]bool),
		engines:              make(map[walk.Key]string),
		colTypes:             make(map[walk.Key]string),
		partitionState:       make(map[walk.Key]*tablePartitionInfo),
		tableIndexes:         make(map[walk.Key]map[string]*indexMeta),
		tableIndexesComplete: make(map[walk.Key]bool),
		tableFKs:             make(map[walk.Key]map[string][]string),
		referencedByFK:       make(map[walk.Key]map[fkRefKey]bool),
	}
	for name := range snap.triggered {
		w.triggered[relKey(snap.db, name)] = true
	}
	for name := range snap.ignoresFKs {
		w.ignoresFKs[relKey(snap.db, name)] = true
	}
	for name, eng := range snap.engines {
		w.engines[relKey(snap.db, name)] = eng
	}
	for tbl, cols := range snap.colTypes {
		for col, bt := range cols {
			w.colTypes[colKey(snap.db, tbl, col)] = bt
		}
	}
	for tbl, pi := range snap.partitionState {
		w.partitionState[relKey(snap.db, tbl)] = pi.clone()
	}
	for tbl, idxs := range snap.tableIndexes {
		rel := relKey(snap.db, tbl)
		cp := make(map[string]*indexMeta, len(idxs))
		for k, im := range idxs {
			cp[k] = im.clone()
		}
		w.tableIndexes[rel] = cp
		w.tableIndexesComplete[rel] = true
	}
	for tbl, fks := range snap.tableFKs {
		rel := relKey(snap.db, tbl)
		cp := make(map[string][]string, len(fks))
		for k, cols := range fks {
			cp[k] = append([]string(nil), cols...)
		}
		w.tableFKs[rel] = cp
	}
	for parent, refs := range snap.referencedByFK {
		pRel := relKey(snap.db, parent)
		cp := make(map[fkRefKey]bool, len(refs))
		for pair := range refs {
			cp[fkRefKey{child: relKey(snap.db, pair[0]), fk: pair[1]}] = true
		}
		w.referencedByFK[pRel] = cp
	}
	for tbl, idxs := range snap.fkGeneratedIndexes {
		for idx, owners := range idxs {
			set := make(map[string]bool, len(owners))
			for _, fk := range owners {
				set[fk] = true
			}
			w.fkGeneratedIndexes[idxKey(snap.db, tbl, idx)] = set
		}
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
	// renames is set while an ALTER TABLE also renames its table, which
	// leaves the foreign keys it adds no parent the walk can resolve
	// (foreignKeyTarget).
	renames bool
	// early is set while ALTER TABLE lowers its drops and renames, which
	// run together against the table as it was; arrived holds the column
	// and index names they bring in (leave, arrive).
	early   bool
	arrived map[walk.Key]bool
	// consumedFKGen tracks synced FK-generated index candidates already
	// consumed earlier in the current statement so later subcommands in the
	// same ALTER TABLE do not still treat them as generated candidates or
	// overwrite their explicit Present/Absent facts with Unknown.
	consumedFKGen map[walk.Key]bool
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

func chkKey(db, table, name string) walk.Key {
	return walk.Key{Space: spaceCheckConstraint, Parent: tablePath(db, table), Name: name}
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
		if o != nil && (strings.EqualFold(o.Name, "ENGINE") || strings.EqualFold(o.Name, "STORAGE ENGINE")) {
			engine = strings.Trim(o.Value, "`'\"")
		}
	}
	return engine
}

// generatedFKs returns the prefix of the foreign key names the server
// generates for a table: MySQL's <table>_ibfk_<n>, which RENAME TABLE
// renames along with the table (verified on 8.0 and 8.4). For a name
// outside ASCII it returns "", which takes in every name: such a foreign
// key's key is its exact spelling, which the table's key need not be.
func (*builder) generatedFKs(table string) string {
	if strings.HasPrefix(table, nonASCII) {
		return ""
	}
	return strings.ToLower(table) + "_ibfk_"
}

// generatedChecks returns the prefix of the check constraint names MySQL
// generates and rewrites on RENAME TABLE / ALTER TABLE RENAME TO:
// <table>_chk_.
func (*builder) generatedChecks(table string) string {
	if strings.HasPrefix(table, nonASCII) {
		return ""
	}
	return strings.ToLower(table) + "_chk_"
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

// callsStoredFunction reports whether the statement may call a stored
// function, which may create or drop a temporary table (verified on
// 8.0.46): a qualified call, an unqualified one to a function the current
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
// but a temporary table or a base table the walk knows.
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
	return !temp && (f != walk.Present || kind != relTable)
}

// forgetTemporary records that a trigger or stored function may have
// created or dropped temporary tables of any name.
func (b *builder) forgetTemporary() {
	b.Forget(walk.Scope{Space: spaceTemporary, Deep: true})
}

// isFKGenerated reports whether k is still a tracked synced FK-generated index
// candidate (not already consumed by an earlier subcommand in this statement).
func (b *builder) isFKGenerated(k walk.Key) bool {
	return len(b.w.fkGeneratedIndexes[k]) > 0 && !b.consumedFKGen[k]
}

// consumeFKGenerated marks k as no longer a generated-FK candidate for the
// remainder of this statement and removes it from w.fkGeneratedIndexes once
// the statement applies.
func (b *builder) consumeFKGenerated(k walk.Key) {
	if !b.isFKGenerated(k) {
		return
	}
	if b.consumedFKGen == nil {
		b.consumedFKGen = make(map[walk.Key]bool)
	}
	b.consumedFKGen[k] = true
	b.Session(func() {
		delete(b.w.fkGeneratedIndexes, k)
	})
}

// onForeignKeyDropped retires any tracked FK-generated index candidates on r
// whose matching foreign keys have all been dropped. Because the update runs
// in b.Session, an ADD INDEX in the same ALTER TABLE still sees the candidate
// as generated (matching mysql:8.0.46), while subsequent statements treat the
// surviving index as an ordinary index.
func (b *builder) onForeignKeyDropped(r tableRef, fkKey string) {
	parent := tablePath(r.db, r.name)
	b.Session(func() {
		for gk, owners := range b.w.fkGeneratedIndexes {
			if gk.Parent != parent {
				continue
			}
			delete(owners, fkKey)
			if len(owners) == 0 {
				delete(b.w.fkGeneratedIndexes, gk)
			}
		}
	})
}

// onIndexAdded updates tracked FK-generated indexes on r after a statement adds
// an index (with key ik when knownName is true): MySQL drops any generated FK
// backing index that the new index makes redundant, so other synced candidates
// on r become Unknown in the ledger, and ik itself becomes an explicit index.
func (b *builder) onIndexAdded(r tableRef, ik walk.Key, knownName bool) {
	parent := tablePath(r.db, r.name)
	var others []walk.Key
	for gk := range b.w.fkGeneratedIndexes {
		if gk.Parent != parent || b.consumedFKGen[gk] {
			continue
		}
		if knownName && gk == ik {
			continue
		}
		others = append(others, gk)
	}
	slices.SortFunc(others, func(a, b walk.Key) int { return strings.Compare(a.Name, b.Name) })
	for _, gk := range others {
		b.Set(gk, walk.Unknown, "")
		b.consumeFKGenerated(gk)
	}
	if knownName {
		b.consumeFKGenerated(ik)
	}
}

// clearTableFKGenerated clears any tracked FK-generated index candidates for
// the table at (db, name) when the table is dropped or recreated.
func (b *builder) clearTableFKGenerated(db, name string) {
	parent := tablePath(db, name)
	rel := relKey(db, name)
	b.Session(func() {
		for gk := range b.w.fkGeneratedIndexes {
			if gk.Parent == parent {
				delete(b.w.fkGeneratedIndexes, gk)
			}
		}
		delete(b.w.engines, rel)
		delete(b.w.partitionState, rel)
		delete(b.w.tableIndexes, rel)
		delete(b.w.tableIndexesComplete, rel)
		delete(b.w.tableFKs, rel)
		for ck := range b.w.colTypes {
			if ck.Parent == parent {
				delete(b.w.colTypes, ck)
			}
		}
		for pRel, refs := range b.w.referencedByFK {
			for ref := range refs {
				if ref.child == rel {
					delete(refs, ref)
				}
			}
			if len(refs) == 0 {
				delete(b.w.referencedByFK, pRel)
			}
		}
	})
}
