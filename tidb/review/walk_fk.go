package review

import (
	"strings"

	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/tidb/ast"
)

// The foreign key checks read the foreign keys the sync lists, as the
// walk follows them (walk_table_shape.go): a table's own, under its path
// wherever it moves, and the ones that reference a table, found from the
// referencing table's synced path. A foreign key the change adds is one
// the walk claims nothing of. Each refusal was verified on v8.5.5.
//
// TiDB enforces the foreign keys it created from v6.6.0 on; one an older
// release created stays unenforced after an upgrade, and the sync can't
// tell the two apart, so the walk takes every synced foreign key for an
// enforced one.

// ownFK is a foreign key of a table, and its key.
type ownFK struct {
	name string
	fk   fkShape
}

// ownFKs returns the foreign keys of the table at r before the statement,
// in the order the sync lists them.
func (b *builder) ownFKs(r tableRef) []ownFK {
	var out []ownFK
	seen := make(map[string]bool)
	for _, ref := range b.w.snap.fkRefs {
		if seen[ref.name] {
			continue
		}
		seen[ref.name] = true
		f, kind := b.l.Lookup(fkDefKey(r.db, r.name, ref.name))
		if f != walk.Present {
			continue
		}
		if fk, ok := decodeFK(kind); ok {
			out = append(out, ownFK{name: ref.name, fk: fk})
		}
	}
	return out
}

// fkUse is a foreign key of another table that references one.
type fkUse struct {
	child, name string
	fk          fkShape
}

// referencingFKs returns the foreign keys that reference the table at r
// before the statement: each of a table still where the sync put it, and
// naming the table r holds by the name the sync knew it by. TiDB follows a
// rename of either table, so a foreign key goes on referencing the table
// it named (verified); the walk finds the ones whose table did not move.
// A table that references itself is left out.
func (b *builder) referencingFKs(r tableRef) []fkUse {
	db := b.w.snap.db
	if r.db != db {
		return nil
	}
	f, origin := b.l.Lookup(originKey(r.db, r.name))
	if f != walk.Present {
		return nil
	}
	var out []fkUse
	for _, ref := range b.w.snap.fkRefs {
		f, kind := b.l.Lookup(fkDefKey(db, ref.child, ref.name))
		if f != walk.Present {
			continue
		}
		fk, ok := decodeFK(kind)
		if !ok || fk.refDB != db || fk.refTable != origin || fk.refTable == ref.child {
			continue
		}
		out = append(out, fkUse{child: ref.child, name: ref.name, fk: fk})
	}
	return out
}

// referencedTable returns where the table a foreign key of the sync
// references is before the statement: at the path the sync knew it by,
// false once it moved, or a table of its name replaced it.
func (b *builder) referencedTable(fk fkShape) (tableRef, bool) {
	if fk.refDB != b.w.snap.db {
		return tableRef{}, false
	}
	if f, origin := b.l.Lookup(originKey(fk.refDB, fk.refTable)); f != walk.Present || origin != fk.refTable {
		return tableRef{}, false
	}
	return tableRef{db: fk.refDB, name: fk.refTable, written: fk.refTable}, true
}

// referencedFail fails DROP TABLE or TRUNCATE of a table a foreign key of
// another table references, while the session checks foreign keys (3730
// for DROP TABLE, 1701 for TRUNCATE; both pass with the checks off).
// DROP TABLE that drops the referencing table too passes. kept reports
// whether a table survives the statement.
func (b *builder) referencedFail(r tableRef, kept func(child string) bool) bool {
	if b.w.fkChecks != review.SwitchOn {
		return false
	}
	for _, u := range b.referencingFKs(r) {
		if !kept(u.child) {
			continue
		}
		b.fail(r.rng, nil, review.UnsupportedForm, "Table "+review.Code(r.written)+" is referenced by foreign key "+review.Code(ident(u.name))+" of "+review.Code(ident(u.child))+".")
		return true
	}
	return false
}

// dropColumnFKFails fails DROP COLUMN of a column a foreign key of the
// table names (1828), or one a foreign key of another table references
// (1829). TiDB refuses either with the checks off too, and when the same
// statement drops the foreign key (verified).
func (b *builder) dropColumnFKFails(r tableRef, col string, rng review.Range) bool {
	k, ok := b.key(spaceColumn, col)
	if !ok {
		return false
	}
	for _, o := range b.ownFKs(r) {
		if contains(o.fk.cols, k) {
			b.fail(rng, nil, review.UnsupportedForm, "Column "+review.Code(ident(col))+" is in foreign key "+review.Code(ident(o.name))+" of "+review.Code(r.written)+", and can't be dropped.")
			return true
		}
	}
	if b.w.renamedColumns {
		return false
	}
	for _, u := range b.referencingFKs(r) {
		if contains(u.fk.refCols, k) {
			b.fail(rng, nil, review.UnsupportedForm, "Column "+review.Code(ident(col))+" is referenced by foreign key "+review.Code(ident(u.name))+" of "+review.Code(ident(u.child))+", and can't be dropped.")
			return true
		}
	}
	return false
}

// prefixOf reports whether cols are an index's first key parts, in order,
// as TiDB matches an index to a foreign key (IsIndexPrefixCovered).
func prefixOf(i indexShape, cols []string) bool {
	if len(i.cols) < len(cols) {
		return false
	}
	for n, c := range cols {
		if i.cols[n] != c {
			return false
		}
	}
	return true
}

// dropIndexFKFails fails DROP INDEX of the one index that serves a
// foreign key of the table, or one of another table that references it
// (1553): an index serves a foreign key whose columns are its first key
// parts, and the walk claims only an index no other index of the table
// holds every column of the foreign key in. TiDB refuses it with the
// checks off too, and when the same statement drops the foreign key
// (verified).
func (b *builder) dropIndexFKFails(r tableRef, name string, rng review.Range) bool {
	k, ok := b.key(spaceIndex, name)
	if !ok || k == primaryIndex {
		return false
	}
	ix, ok := b.indexesOf(r)
	if !ok {
		return false
	}
	var dropped *indexShape
	for n := range ix {
		if ix[n].name == k {
			dropped = &ix[n]
		}
	}
	if dropped == nil {
		return false
	}
	needs := func(cols []string) bool {
		if !prefixOf(*dropped, cols) {
			return false
		}
		for _, i := range ix {
			// A key part the walk can't read may be a column of cols.
			if i.name != k && (i.hasAll(cols) || contains(i.cols, "")) {
				return false
			}
		}
		return true
	}
	fail := func(fk string) bool {
		b.fail(rng, nil, review.IndexBacksConstraint, "Index "+review.Code(ident(name))+" of "+review.Code(r.written)+" serves foreign key "+review.Code(ident(fk))+", and no other index does.")
		return true
	}
	for _, o := range b.ownFKs(r) {
		if needs(o.fk.cols) {
			return fail(o.name)
		}
	}
	if b.w.renamedColumns {
		return false
	}
	for _, u := range b.referencingFKs(r) {
		if needs(u.fk.refCols) {
			return fail(u.name)
		}
	}
	return false
}

// typeChangeFKFails fails MODIFY or CHANGE that gives a column of a
// foreign key a type its counterpart does not match (3780): one of the
// table's own foreign keys, whose referenced column's type the walk knows,
// or one of another table that references the column. TiDB refuses it
// when the same statement drops the foreign key too (verified).
func (b *builder) typeChangeFKFails(r tableRef, col string, def *ast.ColumnDef, rng review.Range) bool {
	k, ok := b.key(spaceColumn, col)
	if !ok || def == nil || b.w.renamedColumns {
		return false
	}
	to, ok := parseColumnType(astColumnType(def.TypeName))
	if !ok {
		return false
	}
	if from, ok := b.typeOf(r.db, r.name, col); !ok || from.base == to.base {
		return false
	}
	mismatch := func(fk string, other tableRef, otherCol string) bool {
		t, ok := b.typeOf(other.db, other.name, otherCol)
		if !ok || !fkIncompatible(to, t) {
			return false
		}
		b.fail(rng, nil, review.UnsupportedForm, "Foreign key "+review.Code(ident(fk))+" would pair "+review.Code(ident(def.Name))+" ("+review.Code(to.String())+") with "+review.Code(ident(other.written))+"."+review.Code(ident(otherCol))+" ("+review.Code(t.String())+"), types that don't match.")
		return true
	}
	for _, o := range b.ownFKs(r) {
		parent, ok := b.referencedTable(o.fk)
		if !ok || parent.db == r.db && parent.name == r.name {
			continue
		}
		for n, c := range o.fk.cols {
			if c == k && mismatch(o.name, parent, o.fk.refCols[n]) {
				return true
			}
		}
	}
	for _, u := range b.referencingFKs(r) {
		child := tableRef{db: b.w.snap.db, name: u.child, written: u.child}
		for n, c := range u.fk.refCols {
			if c == k && mismatch(u.name, child, u.fk.cols[n]) {
				return true
			}
		}
	}
	return false
}

// addFKTypesFail fails a foreign key a CREATE or ALTER TABLE adds whose
// columns' types do not match the referenced columns' (3780), with the
// checks off too (verified). childType gives a column's type in the
// statement's table, false when the walk does not know it.
func (b *builder) addFKTypesFail(r tableRef, c *ast.Constraint, childType func(col string) (columnType, bool), rng review.Range) bool {
	if c == nil || c.Type != ast.ConstrForeignKey || c.RefTable == nil || len(c.Columns) == 0 || len(c.Columns) != len(c.RefColumns) {
		return false
	}
	ref := *c.RefTable
	if ref.Schema == "" {
		if strings.HasPrefix(r.db, nonASCII) {
			return false
		}
		ref.Schema = r.db
	}
	p, st := b.ref(&ref)
	if st != refOK || p.db == r.db && p.name == r.name {
		return false
	}
	if f, kind := b.l.Lookup(relKey(p.db, p.name)); f != walk.Present || kind != relTable {
		return false
	}
	for n, col := range c.Columns {
		ct, ok := childType(col)
		if !ok {
			continue
		}
		pt, ok := b.typeOf(p.db, p.name, c.RefColumns[n])
		if !ok || !fkIncompatible(ct, pt) {
			continue
		}
		b.fail(rng, nil, review.UnsupportedForm, "A foreign key would pair "+review.Code(ident(col))+" ("+review.Code(ct.String())+") with "+review.Code(table(c.RefTable))+"."+review.Code(ident(c.RefColumns[n]))+" ("+review.Code(pt.String())+"), types that don't match.")
		return true
	}
	return false
}
