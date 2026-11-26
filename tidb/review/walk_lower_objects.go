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
	return b.tableOfRef(r, st)
}

func (b *builder) tableOfRef(r tableRef, st refState) (tableRef, bool, bool) {
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
	b.need(relKey(r.db, r.name), walk.Present, nil, r.rng, notFound("table", r.written, ""), "")
	return r, f == walk.Present, true
}

func (b *builder) createIndex(v *ast.CreateIndexStmt) bool {
	r, st := b.ref(v.Table)
	switch st {
	case refFailed:
		return true
	case refUnknown:
		return false
	}
	f, temp, kind := b.resolve(r)
	if temp || f == walk.Unknown {
		return false
	}
	if f == walk.Absent {
		b.need(relKey(r.db, r.name), walk.Present, nil, r.rng, notFound("table", r.written, ""), "")
		return true
	}
	rng := b.anchor()
	// TiDB v8.5.5 refuses CREATE FULLTEXT INDEX and CREATE SPATIAL INDEX (8200)
	// once the relation exists, including on views and sequences.
	if v.Fulltext {
		b.fail(rng, nil, review.UnsupportedForm, engineName+" doesn't support FULLTEXT indexes.")
		return true
	}
	if v.Spatial {
		b.fail(rng, nil, review.UnsupportedForm, engineName+" doesn't support SPATIAL indexes.")
		return true
	}
	if kind == relView {
		// Functional expression indexes succeed on views, while a column key
		// part fails with 1105 ("Unknown charset ").
		for _, ic := range v.Columns {
			if ic != nil && !ic.Functional {
				if _, ok := ic.Expr.(*ast.ColumnRef); ok {
					b.fail(rng, []walk.Key{relKey(r.db, r.name)}, review.UnsupportedForm, engineName+" doesn't support column indexes on view "+review.Code(r.written)+".")
					return true
				}
			}
		}
		return false
	}
	if kind == relSequence {
		return false
	}
	if b.tooLong(v.IndexName, rng) || b.indexFuncFails(v.Columns, rng) {
		return true
	}
	k, ok := b.key(spaceIndex, v.IndexName)
	if !ok {
		return false
	}
	if k == primaryIndex {
		b.fail(rng, nil, review.ReservedName, reservedPrimary) // 1280
		return true
	}
	ik := idxKey(r.db, r.name, k)
	if v.IfNotExists {
		if existing, _ := b.l.Lookup(ik); existing == walk.Absent {
			if !v.Vector {
				b.checkIndexColumnsExist(r, v.Columns, nil, rng)
			}
		}
		b.Next()
		b.Optional()
	} else {
		b.need(ik, walk.Absent, nil, rng, "", alreadyExists("index", ident(v.IndexName), r.written)) // 1061
		if !v.Vector {
			b.checkIndexColumnsExist(r, v.Columns, nil, rng)
		}
	}
	b.Set(ik, walk.Present, "")
	b.Set(indexSetKey(r.db, r.name), walk.Unknown, "")
	b.Set(columnSetKey(r.db, r.name), walk.Unknown, "")
	return true
}

func (b *builder) dropIndex(v *ast.DropIndexStmt) bool {
	// A view has no index to drop: DROP INDEX i ON v is 1091, not 1347
	// (verified on v8.5.5); with IF EXISTS it is a note.
	r, st := b.ref(v.Table)
	if st == refOK && !v.IfExists {
		if f, temp, kind := b.resolve(r); !temp && f == walk.Present && kind == relView {
			b.fail(b.anchor(), []walk.Key{relKey(r.db, r.name)}, review.ObjectNotFound, notFound("index", ident(v.Name), r.written))
			return true
		}
	}
	r, ok, lowered := b.tableOfRef(r, st)
	if !ok {
		return lowered
	}
	if b.dropIndexFKFails(r, v.Name, b.anchor()) {
		return true
	}
	b.Set(indexSetKey(r.db, r.name), walk.Unknown, "")
	b.Set(columnSetKey(r.db, r.name), walk.Unknown, "")
	return b.dropIndexName(r, v.Name, b.anchor(), v.IfExists)
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
	// TiDB v8.5.5 resolves a view body's stars before checking whether the
	// view name collides (1051 ahead of 1050).
	currentDB := ""
	if !b.w.lost && !b.w.noDatabase && !strings.HasPrefix(b.w.current, nonASCII) {
		currentDB = b.w.current
	}
	if viewBodyHasNoTableStar(v.Select, b.w.sql, currentDB) {
		b.fail(b.anchor(), nil, review.UnsupportedForm, review.Code("CREATE VIEW")+" expands "+review.Code("*")+" with no table to read from.")
		return true
	}
	if b.viewBodyFails(v) {
		return true
	}
	k := relKey(r.db, r.name)
	if v.OrReplace {
		if f, _ := b.l.Lookup(k); f == walk.Present {
			b.needKind(k, []string{relView}, r.rng, "", r.written, "view") // 1347
		}
	} else {
		b.needFree(k, r.rng, r.written, relView) // 1050
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
			b.needKind(k, []string{relView}, r.rng, notFound("view", r.written, ""), r.written, "view") // 1051, 1347
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
		b.needFree(k, r.rng, r.written, relSequence) // 1050
		prior = walk.Absent
	}
	if prior == walk.Absent && b.checkSequenceOptionsValid(r.written, v.Options, b.anchor()) {
		return true
	}
	if prior != walk.Present {
		b.Set(sequenceBoundsKey(r.db, r.name), walk.Unknown, "")
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

// alterSequence follows ALTER SEQUENCE [IF EXISTS]. The target must exist
// (1146, even with IF EXISTS on TiDB v8.5.5) as a sequence (1347 for a base
// table or view). A temporary table is no sequence and shadows none.
func (b *builder) alterSequence(v *ast.AlterSequenceStmt) bool {
	r, st := b.ref(v.Name)
	switch st {
	case refFailed:
		return true
	case refUnknown:
		return false
	}
	if v.Name.Schema != "" {
		b.databaseNeeds(r.db, v.Name.Schema, r.rng) // 1049
	}
	k := relKey(r.db, r.name)
	b.needKind(k, []string{relSequence}, r.rng, notFound("sequence", r.written, ""), r.written, "sequence") // 1146, 1347
	if f, kind := b.l.Lookup(k); f == walk.Present && kind == relSequence {
		if !b.checkSequenceOptionsValid(r.written, v.Options, b.anchor()) {
			b.sequenceBoundsFail(r, v.Options)
		}
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
			b.needKind(k, []string{relSequence}, r.rng, notFound("sequence", r.written, ""), r.written, "sequence") // 4139, 1347
		}
		switch {
		case kind == relSequence:
			b.Set(k, walk.Absent, "")
		case f == walk.Present && kind == "":
			b.Set(k, walk.Unknown, "")
		}
		b.Set(sequenceBoundsKey(r.db, r.name), walk.Unknown, "")
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
	return b.relationsExist(targets)
}

// relationsExist adds the check that each target exists (1146) as a
// relation of any kind: a table, a temporary table, a view, or a sequence.
func (b *builder) relationsExist(targets []*ast.TableRef) bool {
	for _, t := range targets {
		r, st := b.ref(t)
		switch st {
		case refFailed:
			return true
		case refUnknown:
			continue
		}
		if f, temp, _ := b.resolve(r); !temp && f == walk.Absent {
			b.need(relKey(r.db, r.name), walk.Present, nil, r.rng, notFound("", r.written, ""), "") // 1146
		}
	}
	return true
}

// set follows SET: foreign_key_checks changes whether a foreign key's
// target is checked. On TiDB a GLOBAL assignment changes the session's
// checks too, though @@session.foreign_key_checks still reads as before,
// and the assignments apply in order (verified: after SET
// foreign_key_checks = 0, SET GLOBAL foreign_key_checks = 1 makes a
// foreign key to a missing table fail with 1824, and SET
// @@session.foreign_key_checks = 1, GLOBAL foreign_key_checks = 0 takes
// it). So every assignment of it counts, whatever its scope.
func (b *builder) set(v *ast.SetStmt) bool {
	for _, a := range v.Assignments {
		if a.Column == nil {
			continue
		}
		name := strings.ToLower(a.Column.Column)
		if rest, ok := strings.CutPrefix(name, "@@"); ok {
			name = rest
			if i := strings.IndexByte(name, '.'); i >= 0 {
				name = name[i+1:]
			}
		}
		if name != "foreign_key_checks" {
			continue
		}
		on := setSwitch(a.Value)
		b.Session(func() { b.w.fkChecks = on })
	}
	return true
}

// setSwitch reads the value SET gives a switch: a number, a boolean, or the word ON or
// OFF, bare or quoted. The parser hands ON over as text and OFF as a name, because ON is
// reserved and OFF is not (verified: each form turns a foreign key to a missing table's
// 1824 on and off).
//
// A unary plus is read through, as TiDB reads it: +0 turns the checks off.
func setSwitch(value ast.ExprNode) review.Switch {
	word := ""
	switch x := value.(type) {
	case *ast.UnaryExpr:
		if x.Op == ast.UnaryPlus {
			return setSwitch(x.Operand)
		}
		return review.SwitchUnknown
	case *ast.IntLit:
		if x.Value != 0 {
			return review.SwitchOn
		}
		return review.SwitchOff
	case *ast.BoolLit:
		if x.Value {
			return review.SwitchOn
		}
		return review.SwitchOff
	case *ast.StringLit:
		word = x.Value
	case *ast.ColumnRef:
		if x.Table == "" {
			word = x.Column
		}
	}
	switch strings.ToLower(word) {
	case "on":
		return review.SwitchOn
	case "off":
		return review.SwitchOff
	}
	return review.SwitchUnknown
}
