package review

import (
	"strings"

	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/tidb/ast"
)

// tableOf resolves the table an index statement names, which must exist
// (1146). It returns false when the walk does not follow the statement on
// it: a temporary table (8200, verified), one it cannot tell, or a view or
// a sequence, on which CREATE INDEX of an expression succeeds (verified).
func (b *builder) tableOf(tr *ast.TableRef) (tableRef, bool, bool) {
	r, st := b.ref(tr)
	switch st {
	case refFailed:
		return r, false, true
	case refUnknown:
		return r, false, false
	}
	f, temp, kind := b.resolve(r)
	if temp || f == walk.Unknown || kind == relView || kind == relSequence {
		return r, false, false
	}
	b.need(relKey(r.db, r.name), walk.Present, nil, r.rng, "Table "+r.written+" does not exist", "")
	return r, f == walk.Present, true
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

// dropView follows DROP VIEW, which looks every view up first. TiDB drops
// the views before a name it then fails on (verified); the statement
// fails either way. IF EXISTS does not cover a table or a sequence of the
// name (1347, verified).
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
		if !v.IfExists || f == walk.Present {
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

// createSequence follows CREATE SEQUENCE [IF NOT EXISTS]. A sequence is a
// relation of its own kind: a table, view, or sequence of the name fails
// it (1050), and a temporary table does not (verified).
func (b *builder) createSequence(v *ast.CreateSequenceStmt) bool {
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
	prior, _ := b.l.Lookup(k)
	if !v.IfNotExists {
		b.need(k, walk.Absent, nil, r.rng, "", "Table "+r.written+" already exists") // 1050
		prior = walk.Absent
	}
	switch prior {
	case walk.Present:
		// A note: nothing is created (verified).
	case walk.Unknown:
		b.Next()
		b.Optional()
		b.Set(k, walk.Present, relSequence)
	default:
		b.Set(k, walk.Present, relSequence)
	}
	return true
}

// dropSequence follows DROP SEQUENCE [IF EXISTS], which looks every
// sequence up first. TiDB drops the sequences before a name it then fails
// on (verified); the statement fails either way. Each must exist (4139)
// as a sequence (1347); IF EXISTS makes either a note. A temporary table
// is no sequence and shadows none (verified). TiDB checks no length here
// (4139 for a name of 65 characters).
func (b *builder) dropSequence(v *ast.DropSequenceStmt) bool {
	b.Order = walk.Lookups
	for _, t := range v.Sequences {
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
			b.need(k, walk.Present, []string{relSequence}, r.rng, "Sequence "+r.written+" does not exist", r.written+" is not a sequence") // 4139, 1347
		}
		switch {
		case kind == relSequence:
			b.Set(k, walk.Absent, "")
		case f == walk.Present && kind == "":
			b.Set(k, walk.Unknown, "")
		}
	}
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
// must exist (1146), as a table, a temporary table, a view, or a
// sequence. TiDB has no triggers or stored functions (1064), so DML runs
// no code that could create or drop a temporary table.
func (b *builder) dml(targets []*ast.TableRef) bool {
	for _, t := range targets {
		r, st := b.ref(t)
		switch st {
		case refFailed:
			return true
		case refUnknown:
			continue
		}
		if f, temp, _ := b.resolve(r); !temp && f == walk.Absent {
			b.need(relKey(r.db, r.name), walk.Present, nil, r.rng, "Table "+r.written+" does not exist", "") // 1146
		}
	}
	return true
}

// set follows SET: foreign_key_checks of the session changes whether a
// foreign key's target is checked. A global setting reaches only later
// sessions.
func (b *builder) set(v *ast.SetStmt) bool {
	if strings.EqualFold(v.Scope, "GLOBAL") || strings.EqualFold(v.Scope, "PERSIST") || strings.EqualFold(v.Scope, "PERSIST_ONLY") {
		return true
	}
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
		default:
			name = strings.TrimPrefix(name, "@@")
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
