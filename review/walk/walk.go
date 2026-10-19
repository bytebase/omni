// Package walk is the engine-neutral core of the WalkThrough review rule:
// a ledger of what is certainly there, certainly not there, or unknown,
// followed through a change statement by statement.
//
// The walk asks one question of every statement: right before it runs,
// does each object it names certainly exist, or certainly not? It models
// nothing else: no types, expressions, dependencies, or privileges. Each
// fact is Present, Absent, or Unknown, and Unknown never yields a finding,
// so a finding can be wrong only when a fact is. Facts come from three
// places, and each is kept sound by its own rule:
//
//   - The synced schema, through the engine's Snapshot: Present for what
//     it lists, and Absent only where the engine knows the sync to be
//     complete.
//   - The change's own statements, through the engine's lowering of each
//     one into a Step: what it needs (Checks) and what it does when it
//     succeeds (Apply). A statement the engine does not lower ends the
//     walk, so coverage can only cost recall.
//   - Bounded uncertainty: a statement whose effect the engine cannot
//     state exactly names the scope it may change, with MayAdd, MayRemove,
//     or Forget, and every fact in that scope that the statement could
//     have falsified becomes Unknown.
//
// A violation is a claim about one statement given that the statements
// before it succeeded: the walk stops at the first one, because a
// statement that fails does nothing and what follows it no longer runs.
//
// The package knows no engine. Spaces, kinds, and the form of names are
// the engine's; the ledger compares keys only for equality.
package walk

import (
	"context"
	"slices"
	"sort"
	"strings"
)

// Fact is what the walk knows about one key.
type Fact uint8

const (
	// Unknown: the object may or may not exist. It never yields a finding.
	Unknown Fact = iota
	// Present: the object certainly exists.
	Present
	// Absent: the object certainly does not exist.
	Absent
)

func (f Fact) String() string {
	switch f {
	case Present:
		return "present"
	case Absent:
		return "absent"
	default:
		return "unknown"
	}
}

// sep joins the names of a parent path. No engine allows it in a name the
// walk is asked about, and an engine that cannot rule it out makes the
// name Unknown before it reaches the ledger.
const sep = "\x00"

// Path is the parent of a key held by the named objects, outermost first:
// Path("public") for a relation in schema public, Path("public", "t") for
// a column of public.t. Path() is the root, the parent of top-level keys
// such as schemas.
func Path(names ...string) string {
	return strings.Join(names, sep)
}

// Names splits a path into the names Path joined.
func Names(path string) []string {
	if path == "" {
		return nil
	}
	return strings.Split(path, sep)
}

// Key names one object. Space is the engine's namespace ("relation",
// "column", ...), Parent the path of the objects that hold it, and Name
// the name as the engine stores it.
type Key struct {
	Space  string
	Parent string
	Name   string
}

// Snapshot is the synced schema as the engine reads it. Lookup returns
// Present for an object the schema lists, with its kind; Absent only for a
// key in a space the engine trusts the sync to list completely; and
// Unknown otherwise. Kind is "" for a fact other than Present.
type Snapshot interface {
	Lookup(Key) (Fact, string)
}

// Scope selects keys: those of Space ("" for every space) whose parent is
// Parent, or, with Deep, Parent or any path under it, and whose name
// begins with Prefix. Scope{Deep: true} selects every key. Kinds narrows a
// MayRemove to the Present facts of those kinds, and of unknown kind; the
// other operations ignore it.
type Scope struct {
	Space  string
	Parent string
	Deep   bool
	Prefix string
	Kinds  []string
}

// Contains reports whether the scope selects the key.
func (s Scope) Contains(k Key) bool {
	if s.Space != "" && s.Space != k.Space {
		return false
	}
	if s.Deep {
		if !Under(k.Parent, s.Parent) {
			return false
		}
	} else if k.Parent != s.Parent {
		return false
	}
	return strings.HasPrefix(k.Name, s.Prefix)
}

// Under reports whether path is root or a path below it. Every path is
// under the root.
func Under(path, root string) bool {
	return root == "" || path == root || strings.HasPrefix(path, root+sep)
}

type opKind uint8

const (
	opSet opKind = iota
	opClear
	opMayAdd
	opMayRemove
	opMove
)

// op is one entry of the ledger's log.
type op struct {
	kind  opKind
	key   Key    // opSet
	fact  Fact   // opSet
	what  string // opSet: the kind of a Present fact
	scope Scope  // opClear, opMayAdd, opMayRemove
	from  string // opMove
	to    string // opMove
}

// Ledger is the walk's state: the snapshot, overlaid by what the change's
// statements did, in order. It answers a lookup by reading its log back
// from the latest entry that bears on the key, so every operation is
// recorded in constant time and judged only when asked.
type Ledger struct {
	snap Snapshot
	ops  []op
	// sets indexes the positions of the opSet entries by key, exact the
	// positions of the scoped entries that are not Deep by parent, space,
	// and prefix, and deep the positions of the Deep scoped entries and
	// moves. Each list ascends. A lookup reads only the lists that can
	// bear on its key, so a long change of many tables, each leaving
	// scopes of its own, costs each lookup little.
	sets  map[Key][]int
	exact map[scopeKey][]int
	deep  []int
}

// scopeKey indexes the scoped entries that are not Deep.
type scopeKey struct{ parent, space, prefix string }

// NewLedger starts a ledger on a snapshot. A nil snapshot knows nothing:
// every key it is asked about is Unknown.
func NewLedger(snap Snapshot) *Ledger {
	return &Ledger{
		snap:  snap,
		sets:  make(map[Key][]int),
		exact: make(map[scopeKey][]int),
	}
}

// Set records the fact a statement established for one key: Present of a
// kind, Absent, or Unknown when the statement may have changed it in a way
// the engine does not follow.
func (l *Ledger) Set(k Key, f Fact, kind string) {
	if f != Present {
		kind = ""
	}
	l.sets[k] = append(l.sets[k], len(l.ops))
	l.ops = append(l.ops, op{kind: opSet, key: k, fact: f, what: kind})
}

// Clear records that nothing exists in the scope: a schema or table the
// statement created holds nothing yet, and one it dropped holds nothing
// any more.
func (l *Ledger) Clear(s Scope) { l.scoped(opClear, s) }

// MayAdd records that the statement may have created objects in the scope
// under names the engine does not derive. An Absent fact in it becomes
// Unknown; a Present one stays.
func (l *Ledger) MayAdd(s Scope) { l.scoped(opMayAdd, s) }

// MayRemove records that the statement may have dropped objects in the
// scope the engine does not list. A Present fact in it becomes Unknown;
// an Absent one stays.
func (l *Ledger) MayRemove(s Scope) { l.scoped(opMayRemove, s) }

// Forget records that the statement may have done anything to the scope:
// every fact in it becomes Unknown.
func (l *Ledger) Forget(s Scope) {
	l.MayAdd(s)
	l.MayRemove(s)
}

func (l *Ledger) scoped(kind opKind, s Scope) {
	if s.Deep {
		l.deep = append(l.deep, len(l.ops))
	} else {
		k := scopeKey{s.Parent, s.Space, s.Prefix}
		l.exact[k] = append(l.exact[k], len(l.ops))
	}
	l.ops = append(l.ops, op{kind: kind, scope: s})
}

// Move records that the objects under the path from are now under the
// path to, as when a table holding columns or a schema holding relations
// is renamed. Nothing is under from afterwards. The move covers
// the keys below the path, not the key that names the moved object
// itself: the engine sets that one. Both paths name an object; a move
// from or to the root, which no statement makes, forgets everything.
func (l *Ledger) Move(from, to string) {
	if from == to {
		return
	}
	if from == "" || to == "" {
		l.Forget(Scope{Deep: true})
		return
	}
	l.deep = append(l.deep, len(l.ops))
	l.ops = append(l.ops, op{kind: opMove, from: from, to: to})
}

// Lookup returns what the ledger knows about the key: the latest fact a
// statement set for it, adjusted by the scopes later statements may have
// changed, or else what the snapshot says, adjusted the same way.
func (l *Ledger) Lookup(k Key) (Fact, string) {
	var added bool
	var removed []Scope
	hi := len(l.ops)
	lists := l.candidates(k)
	for {
		p := -1
		for _, list := range lists {
			p = max(p, lastBefore(list, hi))
		}
		if p < 0 {
			break
		}
		hi = p
		o := &l.ops[p]
		switch o.kind {
		case opSet:
			return adjust(o.fact, o.what, added, removed)
		case opClear:
			if o.scope.Contains(k) {
				return adjust(Absent, "", added, removed)
			}
		case opMayAdd:
			if o.scope.Contains(k) {
				added = true
			}
		case opMayRemove:
			if o.scope.Contains(k) {
				removed = append(removed, o.scope)
			}
		case opMove:
			switch {
			case Under(k.Parent, o.to):
				k.Parent = o.from + k.Parent[len(o.to):]
				lists = l.candidates(k)
			case Under(k.Parent, o.from):
				return adjust(Absent, "", added, removed)
			}
		}
	}
	f, kind := Unknown, ""
	if l.snap != nil {
		f, kind = l.snap.Lookup(k)
	}
	return adjust(f, kind, added, removed)
}

// candidates lists the index lists whose entries can bear on the key: its
// sets, the scopes of its parent in its space or every space under each
// prefix of its name, and the deep entries.
func (l *Ledger) candidates(k Key) [][]int {
	lists := make([][]int, 0, 4)
	add := func(list []int) {
		if len(list) > 0 {
			lists = append(lists, list)
		}
	}
	add(l.sets[k])
	for n := 0; n <= len(k.Name); n++ {
		add(l.exact[scopeKey{k.Parent, k.Space, k.Name[:n]}])
		add(l.exact[scopeKey{k.Parent, "", k.Name[:n]}])
	}
	add(l.deep)
	return lists
}

// lastBefore returns the largest position in the ascending list below hi,
// or -1.
func lastBefore(list []int, hi int) int {
	i := sort.SearchInts(list, hi)
	if i == 0 {
		return -1
	}
	return list[i-1]
}

// adjust applies what later statements may have done to a fact: a scope
// they may have added to makes an Absent fact Unknown, and one they may
// have removed from, of the fact's kind, a Present one.
func adjust(f Fact, kind string, added bool, removed []Scope) (Fact, string) {
	switch f {
	case Absent:
		if added {
			return Unknown, ""
		}
	case Present:
		for _, s := range removed {
			if s.Kinds == nil || kind == "" || slices.Contains(s.Kinds, kind) {
				return Unknown, ""
			}
		}
		return f, kind
	}
	return f, ""
}

// Check is one condition a statement needs to succeed: that Key is Want,
// and, for Present with Kinds set, of one of those kinds. Note is the
// engine's own data for the finding, returned in the Violation.
type Check struct {
	Key   Key
	Want  Fact
	Kinds []string
	Note  any
}

// Violation is a check the ledger certainly contradicts: Got is the fact
// found, Absent or Present, and Kind the kind found when Present.
type Violation struct {
	Statement int
	Check     Check
	Got       Fact
	Kind      string
}

// Violates judges one check against the ledger. A check is violated only
// when the ledger shows the opposite for certain: Absent for a Present
// check, Present for an Absent check, or Present of a known kind outside
// Kinds. A Present fact of unknown kind satisfies any Kinds.
func (l *Ledger) Violates(c Check) (Violation, bool) {
	f, kind := l.Lookup(c.Key)
	v := Violation{Statement: -1, Check: c, Got: f, Kind: kind}
	switch c.Want {
	case Present:
		if f == Absent {
			return v, true
		}
		if f == Present && c.Kinds != nil && kind != "" && !slices.Contains(c.Kinds, kind) {
			return v, true
		}
	case Absent:
		if f == Present {
			return v, true
		}
	}
	return Violation{}, false
}

// Step is one statement as the engine lowered it: the checks it needs,
// judged against the state before it runs, and what it does when it
// succeeds. Fail is a violation the engine established itself while
// lowering, from the same state, for a condition no single check states,
// such as a name a search path resolves nowhere; it is reported before
// the checks are judged. Last ends the walk after the checks, for a
// statement whose effect is not bounded. A check on a key the same
// statement also changes would be judged in the wrong order, so the engine
// leaves it out.
type Step struct {
	Fail   *Violation
	Checks []Check
	Apply  func(*Ledger)
	Last   bool
}

// Lower lowers statement i against the ledger as the statements before it
// left it. It must not change the ledger: changes belong in Step.Apply. A
// nil Step ends the walk before statement i.
type Lower func(i int, l *Ledger) *Step

// Run walks n statements on the snapshot and returns the first violation,
// or nil when the walk reaches the end or a statement that ends it
// without one.
func Run(ctx context.Context, n int, snap Snapshot, lower Lower) (*Violation, error) {
	l := NewLedger(snap)
	for i := range n {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		step := lower(i, l)
		if step == nil {
			return nil, nil
		}
		if step.Fail != nil {
			v := *step.Fail
			v.Statement = i
			return &v, nil
		}
		for _, c := range step.Checks {
			if v, ok := l.Violates(c); ok {
				v.Statement = i
				return &v, nil
			}
		}
		if step.Last {
			return nil, nil
		}
		if step.Apply != nil {
			step.Apply(l)
		}
	}
	return nil, nil
}
