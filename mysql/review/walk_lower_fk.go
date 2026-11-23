package review

import (
	"sort"
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// referencingFKs returns the keys of the foreign keys, as their own tables
// hold them, that reference the table rel from another table not in skip:
// each one the walk knows of whose table and foreign key are still there.
// The walk knows those the sync lists in the target database.
func (b *builder) referencingFKs(rel walk.Key, skip map[walk.Key]bool) []walk.Key {
	var out []walk.Key
	for ref := range b.w.referencedByFK[rel] {
		if ref.child == rel || skip[ref.child] {
			continue
		}
		if f, kind := b.l.Lookup(ref.child); f != walk.Present || kind != relTable {
			continue
		}
		names := walk.Names(ref.child.Parent)
		if len(names) != 1 {
			continue
		}
		k := walk.Key{Space: spaceTableForeignKey, Parent: tablePath(names[0], ref.child.Name), Name: ref.fk}
		if f, _ := b.l.Lookup(k); f != walk.Present {
			continue
		}
		out = append(out, k)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Parent != out[j].Parent {
			return out[i].Parent < out[j].Parent
		}
		return out[i].Name < out[j].Name
	})
	return out
}

// ownFKs returns the keys of a table's own foreign keys the walk knows are
// still there.
func (b *builder) ownFKs(r tableRef) []walk.Key {
	var out []walk.Key
	for name := range b.w.tableFKs[relKey(r.db, r.name)] {
		k := walk.Key{Space: spaceTableForeignKey, Parent: tablePath(r.db, r.name), Name: name}
		if f, _ := b.l.Lookup(k); f == walk.Present {
			out = append(out, k)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// failReferenced records that the statement fails because of foreign keys
// that are there, which the violation claims present.
func (b *builder) failReferenced(rng review.Range, keys []walk.Key, message string) {
	b.Fail(keys, walk.Present, &walkNote{rng: rng, absent: message, present: message, absentCode: review.UnsupportedForm, presentCode: review.UnsupportedForm})
}

// referencedBy is the text for a table another table's foreign key
// references.
func referencedBy(table string, fk walk.Key) string {
	names := walk.Names(fk.Parent)
	return "Table " + review.Code(table) + " is referenced by foreign key " + review.Code(fk.Name) + " of " + review.Code(names[len(names)-1]) + "."
}

// dropTablesReferencedFail fails a DROP TABLE, with foreign keys checked,
// that drops a table another table's foreign key references while it
// leaves that table (3730, verified on mysql:8.0.46, IF EXISTS or not).
// MySQL looks every table up first, so the walk judges only a statement
// whose every table it knows is a table.
func (b *builder) dropTablesReferencedFail(dropped []tableRef) bool {
	if b.w.fkChecks != review.SwitchOn || len(dropped) == 0 {
		return false
	}
	skip := make(map[walk.Key]bool, len(dropped))
	for _, r := range dropped {
		skip[relKey(r.db, r.name)] = true
	}
	for _, r := range dropped {
		if refs := b.referencingFKs(relKey(r.db, r.name), skip); len(refs) > 0 {
			b.failReferenced(r.rng, refs[:1], referencedBy(r.written, refs[0]))
			return true
		}
	}
	return false
}

// truncateReferencedFails fails a TRUNCATE, with foreign keys checked, of a
// table another table's foreign key references (1701; a foreign key of the
// table to itself doesn't count, verified on mysql:8.0.46).
func (b *builder) truncateReferencedFails(r tableRef) bool {
	if b.w.fkChecks != review.SwitchOn {
		return false
	}
	if refs := b.referencingFKs(relKey(r.db, r.name), nil); len(refs) > 0 {
		b.failReferenced(r.rng, refs[:1], referencedBy(r.written, refs[0]))
		return true
	}
	return false
}

// engineChangeFails fails an ALTER TABLE that moves a table in a foreign
// key, its own or another's referencing it, to an engine that keeps none
// (3776, foreign_key_checks on or off, verified on mysql:8.0.46).
func (b *builder) engineChangeFails(r tableRef, engine string) bool {
	if keepsFKs(engine) || !mysqlEngines[strings.ToLower(engine)] {
		return false
	}
	keys := b.ownFKs(r)
	if len(keys) == 0 {
		rel := relKey(r.db, r.name)
		keys = b.referencingFKs(rel, nil)
		if len(keys) == 0 {
			for ref := range b.w.referencedByFK[rel] {
				if ref.child != rel {
					continue
				}
				k := walk.Key{Space: spaceTableForeignKey, Parent: tablePath(r.db, r.name), Name: ref.fk}
				if f, _ := b.l.Lookup(k); f == walk.Present {
					keys = append(keys, k)
					break
				}
			}
		}
	}
	if len(keys) == 0 {
		return false
	}
	owner := walk.Names(keys[0].Parent)
	in := " is in foreign key " + review.Code(keys[0].Name) + " of " + review.Code(owner[len(owner)-1])
	if relKey(owner[0], owner[len(owner)-1]) == relKey(r.db, r.name) {
		in = " has foreign key " + review.Code(keys[0].Name)
	}
	b.failReferenced(b.anchor(), keys[:1], "Table "+review.Code(r.written)+in+", and engine "+review.Code(engine)+" keeps no foreign keys.")
	return true
}

// fkRef is a foreign key the sync lists in the target database: the table
// it references and the columns it pairs, by key.
type fkRef struct {
	parent   walk.Key
	cols     []string
	refCols  []string
	refTable string
}

// isStringType reports a character or binary string type.
func isStringType(t string) bool {
	switch t {
	case "char", "varchar", "binary", "varbinary", "tinytext", "text", "mediumtext", "longtext",
		"tinyblob", "blob", "mediumblob", "longblob":
		return true
	}
	return false
}

// fkTypesClash reports two column types a foreign key can't pair (3780):
// integers of two sizes, or an integer and a string. Others are left
// alone.
func fkTypesClash(a, b string) bool {
	if a == "" || b == "" || a == b {
		return false
	}
	ai, bi := isIntegerBaseType(a), isIntegerBaseType(b)
	return ai && bi || ai && isStringType(b) || bi && isStringType(a)
}

// fkClash is the text for a foreign key pairing columns of clashing types.
func fkClash(fk, col, colType, parent, refCol, refType string) string {
	what := "A foreign key"
	if fk != "" {
		what = "Foreign key " + review.Code(fk)
	}
	return what + " pairs " + review.Code(ident(col)) + " (" + review.Code(colType) + ") with " +
		review.Code(parent) + "." + review.Code(ident(refCol)) + " (" + review.Code(refType) + "), types that don't match."
}

// fkAddFails fails a foreign key a CREATE TABLE or ALTER TABLE adds, with
// foreign keys checked, whose referenced table the walk knows: a referenced
// column the table doesn't have (3734), or a pair of columns of clashing
// types (3780). childType gives the type each column of the child has
// once the statement runs, "" when not known.
func (b *builder) fkAddFails(r tableRef, c *ast.Constraint, childType func(string) string) bool {
	if b.w.fkChecks != review.SwitchOn || c == nil || c.Type != ast.ConstrForeignKey || c.RefTable == nil || b.renames {
		return false
	}
	ref := *c.RefTable
	if ref.Schema == "" {
		ref.Schema = r.db
	}
	t, st := b.ref(&ref)
	if st != refOK || t.db == r.db && t.name == r.name {
		return false
	}
	if f, temp, kind := b.resolve(t); temp || f != walk.Present || kind != relTable {
		return false
	}
	parent := table(c.RefTable)
	for i, rc := range c.RefColumns {
		k, ok := b.key(spaceColumn, rc)
		if !ok {
			return false
		}
		ck := colKey(t.db, t.name, k)
		f, _ := b.l.Lookup(ck)
		if f == walk.Absent {
			b.fail(rangeOf(c.Loc), []walk.Key{ck}, review.ObjectNotFound, "Referenced column "+review.Code(ident(rc))+" doesn't exist in "+review.Code(parent)+".")
			return true
		}
		if f != walk.Present || i >= len(c.Columns) {
			continue
		}
		if ct, pt := childType(c.Columns[i]), b.w.colTypes[ck]; fkTypesClash(ct, pt) {
			name := c.Name
			if name == "" {
				name = c.IndexName
			}
			b.fail(rangeOf(c.Loc), nil, review.UnsupportedForm, fkClash(name, c.Columns[i], ct, parent, rc, pt))
			return true
		}
	}
	return false
}

// fkTypeChangeFails fails an ALTER TABLE, with foreign keys checked, that
// changes the type of a column a foreign key the sync lists pairs, on
// either side, to one that clashes with the other side's (3780). A foreign
// key of the table the statement drops (dropped) is left alone. newType
// gives a column's type once the statement runs, "" when the statement
// doesn't change it or the walk doesn't know.
func (b *builder) fkTypeChangeFails(r tableRef, dropped map[string]bool, newType func(string) string) bool {
	if b.w.fkChecks != review.SwitchOn {
		return false
	}
	rel := relKey(r.db, r.name)
	present := func(child walk.Key, fk string) bool {
		names := walk.Names(child.Parent)
		if len(names) != 1 {
			return false
		}
		f, _ := b.l.Lookup(walk.Key{Space: spaceTableForeignKey, Parent: tablePath(names[0], child.Name), Name: fk})
		return f == walk.Present
	}
	// The table's own foreign keys.
	fks := b.w.fkRefs[rel]
	names := make([]string, 0, len(fks))
	for n := range fks {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		fr := fks[n]
		if fr.parent == rel || dropped[n] || !present(rel, n) || !b.w.referencedByFK[fr.parent][fkRefKey{child: rel, fk: n}] {
			continue
		}
		pnames := walk.Names(fr.parent.Parent)
		for i, col := range fr.cols {
			nt := newType(col)
			if nt == "" || i >= len(fr.refCols) {
				continue
			}
			pt := b.w.colTypes[colKey(pnames[0], fr.parent.Name, fr.refCols[i])]
			if fkTypesClash(nt, pt) {
				b.fail(b.anchor(), nil, review.UnsupportedForm, fkClash(n, col, nt, fr.refTable, fr.refCols[i], pt))
				return true
			}
		}
	}
	// Foreign keys of other tables referencing it.
	refs := make([]fkRefKey, 0, len(b.w.referencedByFK[rel]))
	for ref := range b.w.referencedByFK[rel] {
		refs = append(refs, ref)
	}
	sort.Slice(refs, func(i, j int) bool {
		if refs[i].child != refs[j].child {
			return refs[i].child.Name < refs[j].child.Name
		}
		return refs[i].fk < refs[j].fk
	})
	for _, ref := range refs {
		if ref.child == rel || !present(ref.child, ref.fk) {
			continue
		}
		fr, ok := b.w.fkRefs[ref.child][ref.fk]
		if !ok || fr.parent != rel {
			continue
		}
		cnames := walk.Names(ref.child.Parent)
		for i, refCol := range fr.refCols {
			nt := newType(refCol)
			if nt == "" || i >= len(fr.cols) {
				continue
			}
			ct := b.w.colTypes[colKey(cnames[0], ref.child.Name, fr.cols[i])]
			if fkTypesClash(ct, nt) {
				b.fail(b.anchor(), nil, review.UnsupportedForm, fkClash(ref.fk, fr.cols[i], ct, r.written, refCol, nt))
				return true
			}
		}
	}
	return false
}

// droppedReferencedColumnFails fails an ALTER TABLE, with foreign keys
// checked, that drops a column another table's foreign key the sync lists
// references (1829, verified on mysql:8.0.46).
func (b *builder) droppedReferencedColumnFails(r tableRef, v *ast.AlterTableStmt) bool {
	if b.w.fkChecks != review.SwitchOn {
		return false
	}
	rel := relKey(r.db, r.name)
	for _, c := range v.Commands {
		if c == nil || c.Type != ast.ATDropColumn {
			continue
		}
		k, ok := b.key(spaceColumn, c.Name)
		if !ok {
			continue
		}
		for _, fk := range b.referencingFKs(rel, nil) {
			names := walk.Names(fk.Parent)
			child := relKey(names[0], names[1])
			fr, ok := b.w.fkRefs[child][fk.Name]
			if !ok || fr.parent != rel || !containsString(fr.refCols, k) {
				continue
			}
			b.failReferenced(rangeOf(c.Loc), []walk.Key{fk}, "Column "+review.Code(ident(c.Name))+" is referenced by foreign key "+review.Code(fk.Name)+" of "+review.Code(names[1])+".")
			return true
		}
	}
	return false
}
