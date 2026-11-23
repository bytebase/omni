package review

import (
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
		for ik, icols := range ft.indexes {
			if !serves(icols, cols) {
				continue
			}
			if droppedIdx[ik] {
				dropped = ik
				continue
			}
			left := slices.DeleteFunc(slices.Clone(icols), func(c string) bool { return droppedCol[c] })
			if serves(left, cols) {
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
			b.Session(func() { clear(b.w.fkTables) })
			return
		}
		keys = append(keys, relKey(db, name))
	}
	b.Session(func() {
		for _, k := range keys {
			delete(b.w.fkTables, k)
		}
	})
}
