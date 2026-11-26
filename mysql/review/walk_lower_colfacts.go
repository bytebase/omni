package review

import (
	"sort"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// colFact is what the sync states of a column: whether it is generated
// (genVirtual, genStored, or ""), whether it takes NULL, its character set
// ("" for none or one the walk doesn't know), its SRID, and whether it is
// invisible.
type colFact struct {
	gen       string
	nullable  bool
	charset   string
	hasSRID   bool
	srid      uint32
	invisible bool
}

func syncedColFact(c *metadata.ColumnMetadata) colFact {
	f := colFact{
		nullable:  c.GetNullable(),
		charset:   normalCharset(c.GetCharacterSet()),
		hasSRID:   c.Srid != nil,
		srid:      c.GetSrid(),
		invisible: c.GetIsInvisible(),
	}
	switch g := c.GetGeneration(); {
	case g == nil || g.GetExpression() == "":
	case g.GetType() == metadata.GenerationMetadata_TYPE_STORED:
		f.gen = genStored
	case g.GetType() == metadata.GenerationMetadata_TYPE_VIRTUAL:
		f.gen = genVirtual
	}
	return f
}

// colShape is what the walk knows of a foreign key's column once the
// statement runs: its base type and character set, "" when not known, and
// whether it is certainly NOT NULL and certainly a stored generated column.
type colShape struct {
	typ     string
	charset string
	notNull bool
	stored  bool
}

// defShape reads a column definition. A column takes the character set it
// states; one it doesn't state is left unknown.
func defShape(col *ast.ColumnDef) colShape {
	s := colShape{typ: astBaseType(col.TypeName), stored: col.Generated != nil && col.Generated.Stored}
	if col.TypeName != nil {
		s.charset = normalCharset(col.TypeName.Charset)
	}
	s.notNull = defNotNull(col)
	return s
}

// defNotNull reports a column definition that states NOT NULL or PRIMARY
// KEY, the last of NULL and NOT NULL counting.
func defNotNull(col *ast.ColumnDef) bool {
	notNull, primary := false, false
	for _, cc := range col.Constraints {
		if cc == nil {
			continue
		}
		switch cc.Type {
		case ast.ColConstrNotNull:
			notNull = true
		case ast.ColConstrNull:
			notNull = false
		case ast.ColConstrPrimaryKey:
			primary = true
		}
	}
	return notNull || primary
}

// factShape reads a synced column the statement doesn't redefine.
func factShape(typ string, f colFact, converts bool) colShape {
	s := colShape{typ: typ, notNull: !f.nullable, stored: f.gen == genStored}
	if !converts {
		s.charset = f.charset
	}
	return s
}

// isTextBaseType reports a character string type, which has a character
// set.
func isTextBaseType(t string) bool {
	switch t {
	case "char", "varchar", "tinytext", "text", "mediumtext", "longtext":
		return true
	}
	return false
}

// syncedCol returns the synced facts of a column of the table that is
// still there, false when the walk doesn't know them.
func (b *builder) syncedCol(r tableRef, name string) (colFact, walk.Key, bool) {
	k, ok := b.key(spaceColumn, name)
	if !ok {
		return colFact{}, walk.Key{}, false
	}
	ck := colKey(r.db, r.name, k)
	f, known := b.w.colFacts[ck]
	if !known {
		return colFact{}, ck, false
	}
	if p, _ := b.l.Lookup(ck); p != walk.Present {
		return colFact{}, ck, false
	}
	return f, ck, true
}

// alterChildShape is colShape for a column of the table an ALTER TABLE
// changes, as it is once the statement runs.
func (b *builder) alterChildShape(r tableRef, d *alterDelta, col string) colShape {
	f, typ, _ := b.alterFinalCol(r, d, col)
	k, ok := b.key(spaceColumn, col)
	if f != walk.Present || !ok {
		return colShape{typ: typ}
	}
	if def := d.defs[k]; def != nil {
		return defShape(def)
	}
	from := col
	if old, renamed := d.renamedTo[k]; renamed {
		from = old
	}
	facts, _, known := b.syncedCol(r, from)
	if !known {
		return colShape{typ: typ}
	}
	return factShape(typ, facts, d.converts)
}

// fkActionFault is the text for a foreign key whose actions its column
// can't take (each verified on mysql:8.0.46 and 8.4.13): SET NULL on a NOT
// NULL column (1830), and CASCADE or SET NULL on update, or SET NULL on
// delete, on a stored generated column (3104). "" when neither holds.
func fkActionFault(c *ast.Constraint, col string, s colShape) string {
	setsNull := c.OnDelete == ast.RefActSetNull || c.OnUpdate == ast.RefActSetNull
	if s.notNull && setsNull {
		return setNullOnNotNull(col, fkName(c))
	}
	if s.stored && (setsNull || c.OnUpdate == ast.RefActCascade) {
		return "Stored generated column " + review.Code(ident(col)) + " can't take a foreign key that sets or cascades it."
	}
	return ""
}

func fkName(c *ast.Constraint) string {
	if c.Name != "" {
		return c.Name
	}
	return c.IndexName
}

// setNullOnNotNull is the text for a NOT NULL column a foreign key sets
// NULL (1830).
func setNullOnNotNull(col, fk string) string {
	what := "a foreign key"
	if fk != "" {
		what = "foreign key " + review.Code(fk)
	}
	return "Column " + review.Code(ident(col)) + " can't be " + review.Code("NOT NULL") + " because " + what + " sets it " + review.Code("NULL") + "."
}

// parentKeyMissing reports a referenced table the walk knows every index
// of, none of which a foreign key on refCols can use (1822; 6125 on 8.4):
// MySQL takes a unique key on exactly those columns, in any order, or an
// index whose first columns they are, in order, counting the primary key's
// columns InnoDB appends to a secondary index. A key part's prefix
// length is left alone.
func (b *builder) parentKeyMissing(t tableRef, refCols []string) bool {
	rel := relKey(t.db, t.name)
	idxs := b.w.tableIndexes[rel]
	if len(idxs) == 0 || !b.w.tableIndexesComplete[rel] || b.w.ignoresFKs[rel] {
		return false
	}
	want := make([]string, 0, len(refCols))
	for _, rc := range refCols {
		k, ok := b.key(spaceColumn, rc)
		if !ok {
			return false
		}
		want = append(want, k)
	}
	var pk []string
	if p := idxs[primaryIndex]; p != nil {
		pk = p.cols
	}
	for _, im := range idxs {
		if im == nil {
			return false
		}
		if im.unique && sameSet(im.cols, want) {
			return false
		}
		cols := im.cols
		if !im.primary {
			cols = append(append([]string(nil), cols...), pk...)
		}
		if len(cols) >= len(want) && slicesEqual(cols[:len(want)], want) {
			return false
		}
	}
	return true
}

func sameSet(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for _, x := range b {
		if !containsString(a, x) {
			return false
		}
	}
	return true
}

func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// parentKeyMissingText is the text for parentKeyMissing.
func parentKeyMissingText(parent string, refCols []string) string {
	cols := make([]string, len(refCols))
	for i, c := range refCols {
		cols[i] = review.Code(ident(c))
	}
	return "Referenced table " + review.Code(parent) + " has no index on (" + strings.Join(cols, ", ") + ") for the foreign key."
}

// fkCharsetClash is the text for a foreign key pairing character columns
// of two character sets (3780).
func fkCharsetClash(fk, col, charset, parent, refCol, refCharset string) string {
	what := "A foreign key"
	if fk != "" {
		what = "Foreign key " + review.Code(fk)
	}
	return what + " pairs " + review.Code(ident(col)) + " (" + review.Code(charset) + ") with " +
		review.Code(parent) + "." + review.Code(ident(refCol)) + " (" + review.Code(refCharset) + "), character sets that don't match."
}

// spatialNullable is the text for a nullable column of a SPATIAL index
// (1252).
func spatialNullable(col string) string {
	return "Column " + review.Code(ident(col)) + " of a " + review.Code("SPATIAL") + " index must be " + review.Code("NOT NULL") + "."
}

// syncedColumnsFail fails an ALTER TABLE, the table's first in the change,
// that changes a column against what the sync states of it (each verified
// on mysql:8.0.46 and 8.4.13):
//
//   - making NOT NULL a column a foreign key of the table sets NULL, with
//     foreign keys checked (1830);
//   - giving a column a FULLTEXT index holds a type it can't index (1283);
//   - making a column a SPATIAL index holds nullable (1252), or changing
//     its SRID (3644); where it does both MySQL's order isn't pinned and
//     the walk says nothing;
//   - adding a SPATIAL index on a nullable column (1252);
//   - leaving the table no visible column (4028).
func (b *builder) syncedColumnsFail(r tableRef, v *ast.AlterTableStmt, d *alterDelta) bool {
	rel := relKey(r.db, r.name)
	idxs := b.w.tableIndexes[rel]
	holds := func(k string, kind func(*indexMeta) bool) bool {
		for name, im := range idxs {
			if im != nil && !d.droppedIndexes[name] && kind(im) && containsString(im.cols, k) {
				return true
			}
		}
		return false
	}
	primary := func(k string) bool {
		if holds(k, func(im *indexMeta) bool { return im.primary }) {
			return true
		}
		for _, im := range d.addedIndexes {
			if im.primary && containsString(im.cols, k) {
				return true
			}
		}
		return false
	}
	fulltext := func(im *indexMeta) bool { return im.fulltext }
	spatial := func(im *indexMeta) bool { return im.spatial }
	visibility := true
	setVisible := make(map[string]bool)
	for _, c := range v.Commands {
		if c == nil {
			continue
		}
		rng := rangeOf(c.Loc)
		switch c.Type {
		case ast.ATAlterColumnVisible, ast.ATAlterColumnInvisible:
			k, ok := b.key(spaceColumn, c.Name)
			if !ok {
				return false
			}
			setVisible[k] = c.Type == ast.ATAlterColumnVisible
			continue
		case ast.ATModifyColumn, ast.ATChangeColumn:
			if c.Column == nil {
				break
			}
			from := c.Column.Name
			if c.Type == ast.ATChangeColumn {
				from = c.Name
			}
			old, ck, ok := b.syncedCol(r, from)
			if !ok {
				break
			}
			k := ck.Name
			now := defShape(c.Column)
			if now.notNull && b.w.fkChecks == review.SwitchOn {
				if fk, ok := b.setsNullFK(r, k, d); ok {
					b.failReferenced(rng, []walk.Key{fk}, setNullOnNotNull(from, fk.Name))
					return true
				}
			}
			if now.typ != "" && !isFulltextBaseType(now.typ) && holds(k, fulltext) {
				b.fail(rng, nil, review.UnsupportedForm, "Column "+review.Code(ident(from))+" cannot be part of a "+review.Code("FULLTEXT")+" index.")
				return true
			}
			if isGeometryBaseType(now.typ) && holds(k, spatial) {
				nullable := !now.notNull && !primary(k)
				sameSRID := old.hasSRID == c.Column.TypeName.HasSRID && (!old.hasSRID || int(old.srid) == c.Column.TypeName.SRID)
				switch {
				case nullable && sameSRID:
					b.fail(rng, nil, review.UnsupportedForm, spatialNullable(from))
					return true
				case !nullable && !sameSRID:
					b.fail(rng, nil, review.UnsupportedForm, "Column "+review.Code(ident(from))+" can't change its SRID while a "+review.Code("SPATIAL")+" index holds it.")
					return true
				}
			}
		case ast.ATAddIndex, ast.ATAddConstraint:
			if c.Constraint == nil || !isSpatialConstraint(c.Constraint) {
				break
			}
			for _, p := range constraintKeyParts(c.Constraint) {
				if p.functional || p.name == "" {
					continue
				}
				k, ok := b.key(spaceColumn, p.name)
				if !ok || primary(k) {
					continue
				}
				nullable := false
				if def := d.defs[k]; def != nil {
					nullable = !defNotNull(def) && isGeometryBaseType(astBaseType(def.TypeName))
				} else if _, renamed := d.renamedTo[k]; !renamed && !d.droppedCols[k] {
					f, ck, known := b.syncedCol(r, p.name)
					nullable = known && f.nullable && isGeometryBaseType(b.w.colTypes[ck])
				}
				if nullable {
					b.fail(rng, nil, review.UnsupportedForm, spatialNullable(p.name))
					return true
				}
			}
		}
		visibility = false
	}
	if !visibility || len(setVisible) == 0 {
		return false
	}
	parent := tablePath(r.db, r.name)
	known, visible := 0, false
	for ck, f := range b.w.colFacts {
		if ck.Parent != parent {
			continue
		}
		known++
		shown, set := setVisible[ck.Name]
		visible = visible || (set && shown) || (!set && !f.invisible)
	}
	for k := range setVisible {
		if _, _, ok := b.syncedCol(r, k); !ok {
			return false
		}
	}
	if known == 0 || visible {
		return false
	}
	b.fail(b.anchor(), nil, review.UnsupportedForm, "Table "+review.Code(r.written)+" must have at least one visible column.")
	return true
}

// setsNullFK returns a foreign key of the table the sync lists, still
// there and not dropped by the statement, that sets column k NULL.
func (b *builder) setsNullFK(r tableRef, k string, d *alterDelta) (walk.Key, bool) {
	fks := b.w.fkRefs[relKey(r.db, r.name)]
	names := make([]string, 0, len(fks))
	for n := range fks {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		fr := fks[n]
		if !fr.setsNull || d.droppedFKs[n] || !containsString(fr.cols, k) {
			continue
		}
		_, own := fkKeys(r, n)
		if f, _ := b.l.Lookup(own); f == walk.Present {
			return own, true
		}
	}
	return walk.Key{}, false
}

// createPrimaryHolds reports a CREATE TABLE whose PRIMARY KEY constraint
// holds column col, which makes it NOT NULL.
func createPrimaryHolds(v *ast.CreateTableStmt, col string) bool {
	for _, c := range v.Constraints {
		if c == nil || c.Type != ast.ConstrPrimaryKey {
			continue
		}
		for _, p := range constraintKeyParts(c) {
			if !p.functional && strings.EqualFold(p.name, col) {
				return true
			}
		}
	}
	return false
}
