package review

import (
	"maps"
	"slices"
	"strings"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// fkTable is a synced InnoDB table's foreign keys and its indexes, as the
// walk reads them to tell whether a drop leaves a foreign key without an
// index (1553). An index's columns are the keys of its key parts in
// order; any index whose leading parts are a foreign key's columns counts
// as one that serves it, whatever its type or prefixes, so the walk never
// takes an index InnoDB would use for one it would not.
type fkTable struct {
	indexes map[string][]string
	fks     map[string][]string
}

// syncedFKTable reads a synced table's foreign keys and indexes, false
// when it has no foreign key or the walk cannot read an index's parts.
func syncedFKTable(t *metadata.TableMetadata, foldTables bool) (*fkTable, bool) {
	if len(t.GetForeignKeys()) == 0 || !keepsFKs(t.GetEngine()) {
		return nil, false
	}
	plain := func(name string) (string, bool) {
		if name == "" || strings.ContainsAny(name, "`'\"() +-*/,.") {
			return "", false
		}
		return nameKey(spaceColumn, name, foldTables)
	}
	ft := &fkTable{indexes: make(map[string][]string), fks: make(map[string][]string)}
	for _, idx := range t.GetIndexes() {
		ik, ok := nameKey(spaceIndex, idx.GetName(), foldTables)
		if idx.GetPrimary() {
			ik, ok = primaryIndex, true
		}
		if !ok || len(idx.GetExpressions()) == 0 {
			return nil, false
		}
		var cols []string
		for _, e := range idx.GetExpressions() {
			k, ok := plain(e)
			if !ok {
				return nil, false
			}
			cols = append(cols, k)
		}
		ft.indexes[ik] = cols
	}
	for _, fk := range t.GetForeignKeys() {
		fkk, ok := nameKey(spaceTableForeignKey, fk.GetName(), foldTables)
		if !ok || len(fk.GetColumns()) == 0 {
			return nil, false
		}
		var cols []string
		for _, c := range fk.GetColumns() {
			k, ok := plain(c)
			if !ok {
				return nil, false
			}
			cols = append(cols, k)
		}
		ft.fks[fkk] = cols
	}
	return ft, true
}

// serves reports whether an index of cols serves a foreign key of fk.
func serves(cols, fk []string) bool {
	return len(cols) >= len(fk) && slices.Equal(cols[:len(fk)], fk)
}

// fkIndexDropFails fails a statement that drops only indexes and columns
// of an InnoDB table the change has not touched, while foreign keys are
// checked, when it drops a column of one of the table's foreign keys or
// leaves one with no index (1553, verified on mariadb:11.8.8). A drop
// takes a column out of each index, and an index with none left goes.
// Indexes of other tables that the table's foreign keys need are not
// judged.
func (b *builder) fkIndexDropFails(r tableRef, dropIdx, dropCols []string, rng review.Range) bool {
	rel := relKey(r.db, r.name)
	ft := b.w.fkTables[rel]
	if ft == nil || b.w.fkChecks != review.SwitchOn || b.w.ignoresFKs[rel] {
		return false
	}
	droppedCol := make(map[string]bool)
	for _, name := range dropCols {
		k, ok := b.key(spaceColumn, name)
		if !ok {
			return false
		}
		droppedCol[k] = true
	}
	droppedIdx := make(map[string]bool)
	for _, name := range dropIdx {
		k, ok := b.key(spaceIndex, name)
		if !ok {
			return false
		}
		droppedIdx[k] = true
	}
	fkNames := make([]string, 0, len(ft.fks))
	for name := range ft.fks {
		fkNames = append(fkNames, name)
	}
	slices.Sort(fkNames)
	for _, name := range fkNames {
		cols := ft.fks[name]
		own := walk.Key{Space: spaceTableForeignKey, Parent: tablePath(r.db, r.name), Name: name}
		for _, c := range cols {
			if droppedCol[c] {
				b.fail(rng, []walk.Key{colKey(r.db, r.name, c), own}, review.UnsupportedForm, "Column "+review.Code(ident(c))+" is needed by a foreign key constraint.") // 1553
				return true
			}
		}
		var dropped string
		survives := false
		// An index serves the foreign key as the drops leave it: one that
		// loses a column may come to serve it.
		for ik, icols := range ft.indexes {
			left := slices.DeleteFunc(slices.Clone(icols), func(c string) bool { return droppedCol[c] })
			switch {
			case droppedIdx[ik]:
				if serves(icols, cols) || serves(left, cols) {
					dropped = ik
				}
			case serves(left, cols):
				survives = true
			}
		}
		if dropped != "" && !survives {
			b.fail(rng, []walk.Key{idxKey(r.db, r.name, dropped), own}, review.UnsupportedForm, "Index "+review.Code(ident(dropped))+" is needed by a foreign key constraint.") // 1553
			return true
		}
	}
	return false
}

// alterDropsOnly returns the indexes and columns an ALTER TABLE drops,
// false when it does anything else. A drop IF EXISTS of a name the table
// lacks is a note, and drops nothing.
func alterDropsOnly(v *ast.AlterTableStmt) (dropIdx, dropCols []string, ok bool) {
	for _, c := range v.Commands {
		switch {
		case c == nil:
			return nil, nil, false
		case c.Type == ast.ATDropIndex && !sameName(c.Name, "PRIMARY"):
			dropIdx = append(dropIdx, c.Name)
		case c.Type == ast.ATDropColumn:
			dropCols = append(dropCols, c.Name)
		default:
			return nil, nil, false
		}
	}
	return dropIdx, dropCols, len(v.Commands) > 0
}

// touchTables records that a statement may change the indexes or foreign
// keys of the tables it names, or of any table when it names one the walk
// cannot resolve.
func (b *builder) touchTables(trs ...*ast.TableRef) {
	keys := make([]walk.Key, 0, len(trs))
	for _, tr := range trs {
		if tr == nil {
			continue
		}
		name, ok := nameKey(spaceRelation, tr.Name, b.w.foldTables)
		db := b.w.current
		if tr.Schema != "" {
			db, ok = nameKey(spaceDatabase, tr.Schema, b.w.foldTables)
		} else if b.w.lost || b.w.noDatabase {
			ok = false
		}
		if !ok {
			b.Session(func() { clear(b.w.fkTables); clear(b.w.synced) })
			return
		}
		keys = append(keys, relKey(db, name))
	}
	b.Session(func() {
		for _, k := range keys {
			delete(b.w.fkTables, k)
			delete(b.w.synced, k)
		}
	})
}

// integerTypes are the integer base types, each a type of its own to a
// foreign key.
var integerTypes = map[string]bool{"tinyint": true, "smallint": true, "mediumint": true, "int": true, "bigint": true}

// fkTypesClash reports whether InnoDB refuses to pair two base types in a
// foreign key, as far as the walk tells: two integer types of different
// sizes, or an integer and a character string (verified on
// mariadb:11.8.8). Strings of one kind and another, and every other pair,
// are left alone.
func fkTypesClash(a, b string) bool {
	switch {
	case integerTypes[a] && integerTypes[b]:
		return a != b
	case integerTypes[a]:
		return characterType(b)
	case integerTypes[b]:
		return characterType(a)
	}
	return false
}

// referencedBy is the text for a table another table's foreign key
// references, which a statement can't drop or empty (1451, 1701).
func referencedBy(table string, u fkUse) string {
	return "Table " + review.Code(table) + " is referenced by foreign key " + review.Code(ident(u.name)) + " of " + review.Code(ident(u.child.Name)) + "."
}

// fkParentFails fails a DROP TABLE or TRUNCATE of a table a foreign key of
// another synced table the change has not touched references, while the
// session checks foreign keys (1451 for DROP TABLE, 1701 for TRUNCATE;
// verified on mariadb:11.8.8). skip lists the tables the statement drops,
// whose foreign keys the walk leaves alone: MariaDB drops in the order
// written, and a child dropped first lets its parent go.
func (b *builder) fkParentFails(r tableRef, skip map[walk.Key]bool) bool {
	if b.w.fkChecks != review.SwitchOn {
		return false
	}
	rel := relKey(r.db, r.name)
	for _, u := range b.referencing(rel) {
		if skip[u.child] {
			continue
		}
		b.fail(r.rng, []walk.Key{rel, u.child}, review.UnsupportedForm, referencedBy(r.written, u))
		return true
	}
	return false
}

// alterOptionsOnly reports whether an ALTER TABLE does nothing but set
// table options, ALGORITHM, and LOCK.
func alterOptionsOnly(v *ast.AlterTableStmt) bool {
	for _, c := range v.Commands {
		if c == nil {
			return false
		}
		switch c.Type {
		case ast.ATTableOption, ast.ATAlgorithm, ast.ATLock:
		default:
			return false
		}
	}
	return len(v.Commands) > 0
}

// fkEngineFails fails an ALTER TABLE that moves a table with a foreign
// key, or one a foreign key of another table references, to an engine
// that keeps none (1217, whatever foreign_key_checks is; verified on
// mariadb:11.8.8). Only a synced table the change has not touched, and a
// statement that sets nothing but table options, are judged.
func (b *builder) fkEngineFails(v *ast.AlterTableStmt, r tableRef) bool {
	engine, sets := alterEngine(v)
	e := strings.ToLower(engine)
	if !sets || keepsFKs(e) || e == "innobase" || !alterOptionsOnly(v) {
		return false
	}
	if !builtinEngines[e] && (b.w.sqlModeSet || !mariadbEngines[e]) {
		return false
	}
	st := b.syncedOf(r)
	if st == nil || st.engine != "innodb" {
		return false
	}
	rel := relKey(r.db, r.name)
	if len(st.fks) > 0 {
		names := slices.Sorted(maps.Keys(st.fks))
		b.fail(b.anchor(), []walk.Key{rel}, review.UnsupportedForm, "Table "+review.Code(r.written)+" has foreign key "+review.Code(ident(names[0]))+", and engine "+review.Code(engine)+" keeps no foreign keys.")
		return true
	}
	if uses := b.referencing(rel); len(uses) > 0 {
		u := uses[0]
		b.fail(b.anchor(), []walk.Key{rel, u.child}, review.UnsupportedForm, "Table "+review.Code(r.written)+" is referenced by foreign key "+review.Code(ident(u.name))+" of "+review.Code(ident(u.child.Name))+", and engine "+review.Code(engine)+" keeps no foreign keys.")
		return true
	}
	return false
}

// droppedFKs returns the keys of the foreign keys or constraints an ALTER
// TABLE drops by name.
func (b *builder) droppedFKs(v *ast.AlterTableStmt) map[string]bool {
	out := make(map[string]bool)
	for _, c := range v.Commands {
		if c == nil || c.Type != ast.ATDropConstraint {
			continue
		}
		if k, ok := b.key(spaceTableForeignKey, c.Name); ok {
			out[k] = true
		}
	}
	return out
}

// fkColumnFails fails an ALTER TABLE that changes the type of a column a
// foreign key pairs, on its own table (1832) or on the table it references
// (1833), or drops a column a foreign key of another table references
// (1829), whatever foreign_key_checks is (verified on mariadb:11.8.8).
// Only foreign keys of synced tables the change has not touched count,
// and not one the statement drops.
func (b *builder) fkColumnFails(v *ast.AlterTableStmt, r tableRef) bool {
	rel := relKey(r.db, r.name)
	if f, temp, kind := b.resolve(r); temp || f != walk.Present || kind != relTable {
		return false
	}
	dropped := b.droppedFKs(v)
	var own map[string]fkRef
	if st := b.syncedOf(r); st != nil && st.engine == "innodb" {
		own = st.fks
	}
	uses := b.referencing(rel)
	for _, c := range v.Commands {
		if c == nil {
			continue
		}
		var from, to string
		switch c.Type {
		case ast.ATModifyColumn, ast.ATChangeColumn:
			if c.Column == nil || c.IfExists {
				continue
			}
			from, to = c.Column.Name, astBaseType(c.Column.TypeName)
			if c.Type == ast.ATChangeColumn {
				from = c.Name
			}
		case ast.ATDropColumn:
			if c.IfExists {
				continue
			}
			from = c.Name
		default:
			continue
		}
		k, ok := b.key(spaceColumn, from)
		if !ok || !b.namedOnce(v, k) {
			continue
		}
		old := b.ledgerType(r, k)
		for _, u := range uses {
			if !slices.Contains(u.ref.refCols, k) {
				continue
			}
			if c.Type == ast.ATDropColumn {
				b.fail(b.anchor(), []walk.Key{colKey(r.db, r.name, k), u.child}, review.UnsupportedForm, "Column "+review.Code(ident(from))+" is referenced by foreign key "+review.Code(ident(u.name))+" of "+review.Code(ident(u.child.Name))+".")
				return true
			}
			if fkTypesClash(old, to) {
				b.fail(b.anchor(), []walk.Key{colKey(r.db, r.name, k), u.child}, review.UnsupportedForm, "Column "+review.Code(ident(from))+" is referenced by foreign key "+review.Code(ident(u.name))+" of "+review.Code(ident(u.child.Name))+", which keeps its type.")
				return true
			}
		}
		if c.Type == ast.ATDropColumn {
			continue
		}
		for _, n := range slices.Sorted(maps.Keys(own)) {
			if dropped[n] || !slices.Contains(own[n].cols, k) {
				continue
			}
			if fkTypesClash(old, to) {
				b.fail(b.anchor(), []walk.Key{colKey(r.db, r.name, k), rel}, review.UnsupportedForm, "Column "+review.Code(ident(from))+" is in foreign key "+review.Code(ident(n))+", which keeps its type.")
				return true
			}
		}
	}
	return false
}
