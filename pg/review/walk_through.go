package review

import (
	"context"
	"slices"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// checkWalkThrough reports the first statement of the change that
// certainly fails on the target, given that the statements before it
// succeeded. It follows the change with the engine-neutral ledger of
// review/walk: each statement is lowered into the checks it needs and the
// facts it establishes, and a statement this file does not lower ends the
// walk. It reads the synced schema as pgSnapshot trusts it, nothing else.
//
// With another spec of the plan changing the same database, which
// Bytebase does not order against this one, the synced schema may not be
// what the change meets: the walk then trusts no synced fact and judges
// the change by its own statements alone. The other spec may also have
// added a DML trigger, so DML ends the walk then. That it added no event
// trigger is assumed, as the walk assumes user code reached only
// implicitly changes no schema: stopping at every DDL statement would
// leave such a change nothing to report.
func checkWalkThrough(ctx context.Context, stmts []statement, target review.Target, r *reporter) error {
	v, err := walkThrough(ctx, stmts, target)
	if err != nil || v == nil {
		return err
	}
	note, ok := v.Check.Note.(*walkNote)
	if !ok {
		return nil
	}
	r.add(review.Finding{
		Rule:      review.WalkThrough,
		Statement: stmts[v.Statement].index,
		Range:     note.rng,
		Message:   escapeLines(note.message(v)),
	})
	return nil
}

// walkThrough walks the change on one target and returns the first
// violation, or nil.
func walkThrough(ctx context.Context, stmts []statement, target review.Target) (*walk.Violation, error) {
	if target.Schema == nil {
		return nil, nil
	}
	w := newWalker(stmts, target)
	var snap walk.Snapshot = w.pg
	if target.OtherSpecs {
		// Another spec may have replaced a synced partitioned table with
		// an ordinary one; anyPartitioned stays, as it only stops more.
		snap = nil
		w.triggered = true
		w.partitioned = make(map[walk.Key]bool)
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

// walker is the session the change runs in, as far as the walk follows
// it, and what of the target it needs besides the snapshot.
type walker struct {
	stmts    []statement
	pg       *pgSnapshot
	database string
	triggers []*metadata.EventTriggerMetadata

	// path is the search path as last set, "$user" unexpanded, or nil
	// when unknown; user is the current user and session the session
	// user, "" when unknown.
	path    []string
	user    string
	session string

	// triggered is true when a DML trigger or rule may fire: the target
	// or the change has one somewhere, or may have one the sync does not
	// list. Its function may run DDL, so DML ends the walk then.
	triggered bool
	// routines lists the names of the functions and procedures the target
	// or the change defines, which a call by that name may run.
	routines map[string]bool
	// partitioned marks the partitioned tables, by relation key, and
	// partitions the partitions the change created of each, which a drop
	// of the parent drops too.
	partitioned    map[walk.Key]bool
	partitions     map[walk.Key][]walk.Key
	anyPartitioned bool
	// temporary is true once the change made a temporary relation: a view
	// of one becomes temporary itself, wherever the SQL put it.
	temporary bool
	// made lists the relations and types the change created, or moved
	// to a name the synced schema does not list, which a DROP SCHEMA
	// without CASCADE finds.
	made []walk.Key
}

func newWalker(stmts []statement, target review.Target) *walker {
	db := target.Schema
	w := &walker{
		stmts:       stmts,
		pg:          newPgSnapshot(db),
		database:    db.GetName(),
		triggers:    enabledTriggers(db),
		user:        target.SessionUser,
		session:     target.SessionUser,
		routines:    make(map[string]bool),
		partitioned: make(map[walk.Key]bool),
		partitions:  make(map[walk.Key][]walk.Key),
	}
	if path := parseSearchPath(db.GetSearchPath()); len(path) > 0 {
		w.path = path
	}
	for _, schema := range db.GetSchemas() {
		for _, t := range schema.GetTables() {
			// A rule's action may call a function, as a trigger does. The
			// sync lists no trigger of a partition, so one may fire unseen.
			if len(t.GetTriggers()) > 0 || len(t.GetRules()) > 0 {
				w.triggered = true
			}
			if len(t.GetPartitions()) > 0 {
				w.partitioned[relKey(schema.GetName(), t.GetName())] = true
				w.anyPartitioned = true
				w.triggered = true
			}
		}
		// The sync lists no trigger of a foreign table either.
		if len(schema.GetExternalTables()) > 0 {
			w.triggered = true
		}
		for _, v := range schema.GetViews() {
			if len(v.GetTriggers()) > 0 || len(v.GetRules()) > 0 {
				w.triggered = true
			}
		}
		for _, v := range schema.GetMaterializedViews() {
			if len(v.GetTriggers()) > 0 {
				w.triggered = true
			}
		}
		for _, f := range schema.GetFunctions() {
			w.routines[f.GetName()] = true
		}
		for _, p := range schema.GetProcedures() {
			w.routines[p.GetName()] = true
		}
	}
	return w
}

// lower implements walk.Lower.
func (w *walker) lower(i int, l *walk.Ledger) *walk.Step {
	n := w.stmts[i].node
	if w.eventTriggerMayFire(n) {
		return nil
	}
	b := &builder{w: w, l: l, stmt: &w.stmts[i]}
	if !b.statement(n) {
		return nil
	}
	return b.step()
}

// eventTriggerMayFire reports whether an enabled event trigger of the
// target may fire on the statement. Its function may run DDL or reject
// the statement before it runs, so the walk ends before it. Every DDL
// statement raises ddl_command_start and ddl_command_end, one that may
// drop objects sql_drop, and one that may rewrite a table table_rewrite;
// a trigger filtered by tags fires only on the commands it names.
func (w *walker) eventTriggerMayFire(n ast.Node) bool {
	if len(w.triggers) == 0 || !isDDL(n) {
		return false
	}
	tag := commandTag(n)
	for _, t := range w.triggers {
		switch strings.ToUpper(t.GetEvent()) {
		case "SQL_DROP":
			if !mayDrop(n) {
				continue
			}
		case "TABLE_REWRITE":
			if !mayRewrite(n) {
				continue
			}
		}
		if len(t.GetTags()) == 0 || tag == "" {
			return true
		}
		for _, want := range t.GetTags() {
			if strings.EqualFold(strings.Trim(want, `"`), tag) {
				return true
			}
		}
	}
	return false
}

// builtinCalls are functions of pg_catalog a call may run without the walk
// losing track of the change: none changes the schema. A call by any
// other name, by a schema-qualified name outside pg_catalog, or by a name
// the target or the change also defines may run user code, which may.
var builtinCalls = map[string]bool{
	"now": true, "clock_timestamp": true, "statement_timestamp": true, "transaction_timestamp": true,
	"timeofday": true, "date_trunc": true, "date_part": true, "extract": true, "age": true,
	"to_char": true, "to_date": true, "to_timestamp": true, "to_number": true,
	"make_date": true, "make_time": true, "make_timestamp": true, "make_timestamptz": true, "make_interval": true,
	"lower": true, "upper": true, "initcap": true, "length": true, "char_length": true, "octet_length": true,
	"substr": true, "substring": true, "btrim": true, "ltrim": true, "rtrim": true, "replace": true,
	"split_part": true, "concat": true, "concat_ws": true, "format": true, "left": true, "right": true,
	"lpad": true, "rpad": true, "repeat": true, "reverse": true, "strpos": true, "position": true,
	"regexp_replace": true, "regexp_match": true, "regexp_matches": true, "regexp_split_to_array": true,
	"md5": true, "encode": true, "decode": true, "abs": true, "ceil": true, "ceiling": true, "floor": true,
	"round": true, "trunc": true, "mod": true, "power": true, "sqrt": true, "random": true,
	"gen_random_uuid": true, "count": true, "sum": true, "avg": true, "min": true, "max": true,
	"array_agg": true, "string_agg": true, "json_agg": true, "jsonb_agg": true, "bool_and": true, "bool_or": true,
	"json_build_object": true, "jsonb_build_object": true, "json_build_array": true, "jsonb_build_array": true,
	"to_json": true, "to_jsonb": true, "jsonb_set": true, "jsonb_insert": true, "jsonb_strip_nulls": true,
	"row_to_json": true, "array_length": true, "array_append": true, "array_prepend": true, "array_cat": true,
	"array_remove": true, "array_position": true, "cardinality": true, "unnest": true, "generate_series": true,
	"nextval": true, "currval": true, "setval": true, "lastval": true,
}

// callsUserCode reports whether the statement calls a function that may be
// user code. A call that may run anything ends the walk after its
// statement.
func (w *walker) callsUserCode(n ast.Node) bool {
	catalogListed := containsString(w.path, "pg_catalog")
	user := false
	ast.Inspect(n, func(m ast.Node) bool {
		if user {
			return false
		}
		fc, ok := m.(*ast.FuncCall)
		if !ok {
			return true
		}
		parts := nameParts(fc.Funcname)
		switch {
		case len(parts) == 0:
			user = true
		case len(parts) == 2 && parts[0] == "pg_catalog":
			user = !builtinCalls[parts[1]]
		case len(parts) == 1:
			name := parts[0]
			user = !builtinCalls[name] || w.routines[name] || catalogListed
		default:
			user = true
		}
		return !user
	})
	return user
}

// builder lowers one statement. A statement is made of parts, one per
// subcommand or object, each with its own checks and effects. A check or
// failure on a key another part also changes is dropped: PostgreSQL does
// not apply the parts in the order written. ALTER TABLE runs its
// subcommands in passes, so ADD COLUMN a, DROP COLUMN a drops a first: it
// fails where a is missing and leaves a new a where a exists (verified on
// PostgreSQL 17 and 18). Effects apply in the order of the parts, which
// passOrder puts in PostgreSQL's.
type builder struct {
	w    *walker
	l    *walk.Ledger
	stmt *statement

	part    int
	checks  []partCheck
	fails   []partFail
	effects []effect
	last    bool
	// drops lists, by storedKey, the columns an ALTER TABLE drops, which
	// are gone before any column is added (passOrder); added lists the
	// columns its ADD COLUMN subcommands so far add, in the order written.
	drops map[string]bool
	added map[string]bool
	// lookups is true when the statement looks every name up before it
	// changes any, so no part's effect reaches another part's checks.
	lookups bool
}

type partCheck struct {
	part  int
	check walk.Check
}

type partFail struct {
	part int
	keys []walk.Key
	v    walk.Violation
}

type effectKind uint8

const (
	effSet effectKind = iota
	effClear
	effMayAdd
	effMayRemove
	effMove
	effSession
)

type effect struct {
	part  int
	kind  effectKind
	key   walk.Key
	fact  walk.Fact
	what  string
	scope walk.Scope
	from  string
	to    string
	fn    func()
}

// next starts the statement's next part.
func (b *builder) next() { b.part++ }

// need adds a check: the key must be want, and of one of kinds when
// Present. A check of an Unknown fact never fails.
func (b *builder) need(k walk.Key, want walk.Fact, kinds []string, note *walkNote) {
	b.checks = append(b.checks, partCheck{b.part, walk.Check{Key: k, Want: want, Kinds: kinds, Note: note}})
}

// fail records that the statement certainly fails, as the lowering
// established from the keys it read, with message.
func (b *builder) fail(keys []walk.Key, got walk.Fact, message string) {
	note := &walkNote{rng: b.anchor(), absent: message, present: message}
	b.fails = append(b.fails, partFail{b.part, keys, walk.Violation{Check: walk.Check{Note: note}, Got: got}})
}

// failAt is fail with its own anchor.
func (b *builder) failAt(rng review.Range, keys []walk.Key, message string) {
	note := &walkNote{rng: rng, absent: message, present: message}
	b.fails = append(b.fails, partFail{b.part, keys, walk.Violation{Check: walk.Check{Note: note}, Got: walk.Absent}})
}

func (b *builder) set(k walk.Key, f walk.Fact, kind string) {
	b.effects = append(b.effects, effect{part: b.part, kind: effSet, key: k, fact: f, what: kind})
}

func (b *builder) clear(s walk.Scope) {
	b.effects = append(b.effects, effect{part: b.part, kind: effClear, scope: s})
}

func (b *builder) mayAdd(s walk.Scope) {
	b.effects = append(b.effects, effect{part: b.part, kind: effMayAdd, scope: s})
}

func (b *builder) mayRemove(s walk.Scope) {
	b.effects = append(b.effects, effect{part: b.part, kind: effMayRemove, scope: s})
}

func (b *builder) forget(s walk.Scope) {
	b.mayAdd(s)
	b.mayRemove(s)
}

func (b *builder) move(from, to string) {
	b.effects = append(b.effects, effect{part: b.part, kind: effMove, from: from, to: to})
}

// session changes the walker's session state when the statement succeeds.
func (b *builder) session(fn func()) {
	b.effects = append(b.effects, effect{part: b.part, kind: effSession, fn: fn})
}

// anchor is the statement's own range, for a finding about all of it.
func (b *builder) anchor() review.Range {
	return rangeOf(b.stmt.loc)
}

// touched reports whether a part other than part changes the key.
func (b *builder) touched(k walk.Key, part int) bool {
	if b.lookups {
		return false
	}
	for _, e := range b.effects {
		if e.part == part {
			continue
		}
		switch e.kind {
		case effSet:
			if e.key == k {
				return true
			}
		case effClear, effMayAdd, effMayRemove:
			if e.scope.Contains(k) {
				return true
			}
		case effMove:
			if walk.Under(k.Parent, e.from) || walk.Under(k.Parent, e.to) {
				return true
			}
		}
	}
	return false
}

// step assembles the walk.Step.
func (b *builder) step() *walk.Step {
	s := &walk.Step{Last: b.last}
	for _, f := range b.fails {
		if !slices.ContainsFunc(f.keys, func(k walk.Key) bool { return b.touched(k, f.part) }) {
			v := f.v
			s.Fail = &v
			break
		}
	}
	for _, c := range b.checks {
		if !b.touched(c.check.Key, c.part) {
			s.Checks = append(s.Checks, c.check)
		}
	}
	effects := b.effects
	s.Apply = func(l *walk.Ledger) {
		for _, e := range effects {
			switch e.kind {
			case effSet:
				l.Set(e.key, e.fact, e.what)
			case effClear:
				l.Clear(e.scope)
			case effMayAdd:
				l.MayAdd(e.scope)
			case effMayRemove:
				l.MayRemove(e.scope)
			case effMove:
				l.Move(e.from, e.to)
			case effSession:
				e.fn()
			}
		}
	}
	return s
}
