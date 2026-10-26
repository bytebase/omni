package walk

import "slices"

// Order is how the parts of one statement see one another.
type Order uint8

const (
	// Passes: the server runs the parts in an order of its own, as
	// PostgreSQL runs ALTER TABLE in passes, so a check or failure on a key
	// another part also changes is left out. The engine adds the parts in
	// the order the server applies their effects.
	Passes Order = iota
	// Lookups: the statement looks every name up before it changes any, as
	// a DROP of several objects does, so every check holds against the
	// state before the statement whatever another part does.
	Lookups
)

// Builder assembles the Step of one statement from its parts: a part is
// one subcommand or one object, with checks, failures, and effects of its
// own. Effects apply in the order they are added.
type Builder struct {
	// Order is how the statement's parts see one another.
	Order Order
	// Last ends the walk after the statement's checks.
	Last bool

	part     int
	optional map[int]bool
	checks   []partCheck
	fails    []partFail
	effects  []effect
}

type partCheck struct {
	part  int
	check Check
}

type partFail struct {
	part int
	v    Violation
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
	key   Key
	fact  Fact
	what  string
	scope Scope
	from  string
	to    string
	fn    func()
}

// Next starts the statement's next part.
func (b *Builder) Next() { b.part++ }

// Optional marks the current part as one the statement may skip, such as
// an IF NOT EXISTS on a name that may exist. Its checks and failures are
// left out, and its exact effects weaken to what holds either way: a Set
// leaves its key Unknown, a Clear may remove what it scopes, and a Move
// forgets both paths. Its session effects run unchanged, so the engine
// adds to an optional part only session effects that hold either way.
func (b *Builder) Optional() {
	if b.optional == nil {
		b.optional = make(map[int]bool)
	}
	b.optional[b.part] = true
}

// Need adds a check: the key must be want, and of one of kinds when
// Present. A check of an Unknown fact never fails. Note is returned in
// the Violation.
func (b *Builder) Need(k Key, want Fact, kinds []string, note any) {
	b.checks = append(b.checks, partCheck{b.part, Check{Key: k, Want: want, Kinds: kinds, Note: note}})
}

// Fail records that the statement certainly fails, as the engine
// established from the keys it read, which the Violation reports with
// their facts; with no keys, from the statement's text alone.
func (b *Builder) Fail(keys []Key, got Fact, note any) {
	b.fails = append(b.fails, partFail{b.part, Violation{Check: Check{Note: note}, Got: got, Keys: keys}})
}

// Set records the fact the statement establishes for a key; see
// Ledger.Set.
func (b *Builder) Set(k Key, f Fact, kind string) {
	b.effects = append(b.effects, effect{part: b.part, kind: effSet, key: k, fact: f, what: kind})
}

// Clear records that nothing exists in the scope; see Ledger.Clear.
func (b *Builder) Clear(s Scope) {
	b.effects = append(b.effects, effect{part: b.part, kind: effClear, scope: s})
}

// MayAdd records that the statement may create objects in the scope; see
// Ledger.MayAdd.
func (b *Builder) MayAdd(s Scope) {
	b.effects = append(b.effects, effect{part: b.part, kind: effMayAdd, scope: s})
}

// MayRemove records that the statement may drop objects in the scope; see
// Ledger.MayRemove.
func (b *Builder) MayRemove(s Scope) {
	b.effects = append(b.effects, effect{part: b.part, kind: effMayRemove, scope: s})
}

// Forget records that the statement may do anything to the scope.
func (b *Builder) Forget(s Scope) {
	b.MayAdd(s)
	b.MayRemove(s)
}

// Move records that the objects under one path are now under another; see
// Ledger.Move.
func (b *Builder) Move(from, to string) {
	b.effects = append(b.effects, effect{part: b.part, kind: effMove, from: from, to: to})
}

// Session records a change to the engine's own session state, made when
// the statement succeeds, in order with the ledger's effects.
func (b *Builder) Session(fn func()) {
	b.effects = append(b.effects, effect{part: b.part, kind: effSession, fn: fn})
}

// touched reports whether a part other than part changes the key.
func (b *Builder) touched(k Key, part int) bool {
	if b.Order == Lookups {
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
			if Under(k.Parent, e.from) || Under(k.Parent, e.to) {
				return true
			}
		}
	}
	return false
}

// Step assembles the statement's Step.
func (b *Builder) Step() *Step {
	s := &Step{Last: b.Last}
	for _, f := range b.fails {
		if b.optional[f.part] {
			continue
		}
		if !slices.ContainsFunc(f.v.Keys, func(k Key) bool { return b.touched(k, f.part) }) {
			v := f.v
			s.Fail = &v
			break
		}
	}
	for _, c := range b.checks {
		if !b.optional[c.part] && !b.touched(c.check.Key, c.part) {
			s.Checks = append(s.Checks, c.check)
		}
	}
	effects := b.effects
	optional := b.optional
	s.Apply = func(l *Ledger) {
		for _, e := range effects {
			weak := optional[e.part]
			switch e.kind {
			case effSet:
				if weak {
					l.Set(e.key, Unknown, "")
				} else {
					l.Set(e.key, e.fact, e.what)
				}
			case effClear:
				if weak {
					l.MayRemove(e.scope)
				} else {
					l.Clear(e.scope)
				}
			case effMayAdd:
				l.MayAdd(e.scope)
			case effMayRemove:
				l.MayRemove(e.scope)
			case effMove:
				if weak {
					l.Forget(Scope{Parent: e.from, Deep: true})
					l.Forget(Scope{Parent: e.to, Deep: true})
				} else {
					l.Move(e.from, e.to)
				}
			case effSession:
				e.fn()
			}
		}
	}
	return s
}
