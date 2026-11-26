package review

import (
	"cmp"
	"slices"
	"strconv"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review/walk"
)

// syncedTable is what the sync says of a base table beyond the names the
// ledger keeps: its engine, partitions, CHECK constraints, generated
// columns, and foreign keys. The walk trusts it only while no statement of
// the change may have changed the table (walker.synced), as it does
// fkTables.
type syncedTable struct {
	// engine is the table's storage engine, lower case.
	engine string
	// part is the table's partitioning, nil when it has none or the walk
	// cannot read it.
	part *partInfo
	// checks maps each CHECK constraint's key to the keys of the columns
	// its expression reads, nil when the walk cannot read them.
	checks map[string][]string
	// generated maps each generated column's key to whether it is
	// STORED; a column the map lacks is no generated column.
	generated map[string]bool
	// fks holds the table's own foreign keys by their keys.
	fks map[string]fkRef
	// triggers holds the table's triggers by their keys, each with its
	// name and its time and event ("BEFORE INSERT"); a key two spellings
	// share holds neither.
	triggers map[string]syncedTrigger
}

// syncedTrigger is a synced trigger's name, time and event.
type syncedTrigger struct {
	name, when string
}

// fkRef is a synced foreign key: its columns, the table it references,
// and the columns there, each as keys.
type fkRef struct {
	cols    []string
	parent  walk.Key
	refCols []string
}

// partInfo is a synced table's partitioning.
type partInfo struct {
	kind metadata.TablePartitionMetadata_Type
	// names holds the keys of the partitions, count how many there are.
	names map[string]bool
	count int
	// cols holds the keys of the columns the partitioning function
	// reads, nil when the walk cannot read them.
	cols []string
	// subparts is set when the partitions have subpartitions.
	subparts bool
	// bounds holds a RANGE partitioning's bounds in order, each an
	// integer, and maxvalue is set when the last is MAXVALUE; bounds is
	// nil when the walk cannot read one.
	bounds   []int64
	maxvalue bool
	// values holds the integers a LIST partitioning's partitions name,
	// nil when the walk cannot read one.
	values map[int64]bool
}

// rangeLike and listLike tell RANGE and LIST partitioning, with COLUMNS
// or without; hashLike is HASH or KEY, LINEAR or not.
func rangeLike(t metadata.TablePartitionMetadata_Type) bool {
	return t == metadata.TablePartitionMetadata_RANGE || t == metadata.TablePartitionMetadata_RANGE_COLUMNS
}

func listLike(t metadata.TablePartitionMetadata_Type) bool {
	return t == metadata.TablePartitionMetadata_LIST || t == metadata.TablePartitionMetadata_LIST_COLUMNS
}

func hashLike(t metadata.TablePartitionMetadata_Type) bool {
	switch t {
	case metadata.TablePartitionMetadata_HASH, metadata.TablePartitionMetadata_LINEAR_HASH,
		metadata.TablePartitionMetadata_KEY, metadata.TablePartitionMetadata_LINEAR_KEY:
		return true
	}
	return false
}

// quotedNames returns the backquoted names an expression the sync lists
// reads, as keys of space, false when it reads none or one the walk
// cannot key. MariaDB writes every column of a stored expression in
// backquotes (information_schema's CHECK_CLAUSE, GENERATION_EXPRESSION,
// and PARTITION_EXPRESSION); a string literal is skipped.
func quotedNames(expr, space string, foldTables bool) ([]string, bool) {
	var out []string
	for i := 0; i < len(expr); i++ {
		switch expr[i] {
		case '\'', '"':
			q := expr[i]
			for i++; i < len(expr) && expr[i] != q; i++ {
				if expr[i] == '\\' {
					i++
				}
			}
		case '`':
			var name strings.Builder
			for i++; i < len(expr); i++ {
				if expr[i] == '`' {
					if i+1 < len(expr) && expr[i+1] == '`' {
						name.WriteByte('`')
						i++
						continue
					}
					break
				}
				name.WriteByte(expr[i])
			}
			k, ok := nameKey(space, name.String(), foldTables)
			if !ok || strings.HasPrefix(k, nonASCII) {
				return nil, false
			}
			out = append(out, k)
		}
	}
	return out, len(out) > 0
}

// newSyncedTable reads what syncedTable holds of a table the sync lists.
func newSyncedTable(db string, t *metadata.TableMetadata, foldTables bool) *syncedTable {
	st := &syncedTable{
		engine:    strings.ToLower(t.GetEngine()),
		checks:    make(map[string][]string),
		generated: make(map[string]bool),
		fks:       make(map[string]fkRef),
		triggers:  make(map[string]syncedTrigger),
	}
	for _, tr := range t.GetTriggers() {
		k, ok := nameKey(spaceTrigger, tr.GetName(), foldTables)
		if !ok {
			continue
		}
		if seen, dup := st.triggers[k]; dup && seen.name != tr.GetName() {
			st.triggers[k] = syncedTrigger{}
			continue
		}
		st.triggers[k] = syncedTrigger{name: tr.GetName(), when: strings.ToUpper(tr.GetTiming() + " " + tr.GetEvent())}
	}
	for _, c := range t.GetCheckConstraints() {
		k, ok := nameKey(spaceColumn, c.GetName(), foldTables)
		if !ok {
			continue
		}
		cols, _ := quotedNames(c.GetExpression(), spaceColumn, foldTables)
		st.checks[k] = cols
	}
	for _, c := range t.GetColumns() {
		g := c.GetGeneration()
		if g == nil {
			continue
		}
		if k, ok := nameKey(spaceColumn, c.GetName(), foldTables); ok {
			st.generated[k] = g.GetType() == metadata.GenerationMetadata_TYPE_STORED
		}
	}
	for _, fk := range t.GetForeignKeys() {
		k, ok := nameKey(spaceTableForeignKey, fk.GetName(), foldTables)
		if !ok {
			continue
		}
		ref := fkRef{}
		pdb := db
		if s := fk.GetReferencedSchema(); s != "" {
			pdb, ok = nameKey(spaceDatabase, s, foldTables)
		}
		parent, ok2 := nameKey(spaceRelation, fk.GetReferencedTable(), foldTables)
		if ok && ok2 {
			ref.parent = relKey(pdb, parent)
		}
		ref.cols = columnKeys(fk.GetColumns(), foldTables)
		ref.refCols = columnKeys(fk.GetReferencedColumns(), foldTables)
		st.fks[k] = ref
	}
	st.part = syncedPartitions(t.GetPartitions(), foldTables)
	return st
}

// columnKeys returns the keys of column names, nil when one is not a
// name the walk can key.
func columnKeys(names []string, foldTables bool) []string {
	out := make([]string, 0, len(names))
	for _, n := range names {
		k, ok := nameKey(spaceColumn, n, foldTables)
		if !ok || strings.HasPrefix(k, nonASCII) {
			return nil
		}
		out = append(out, k)
	}
	return out
}

// syncedPartitions reads a table's partitioning, nil when it has none or
// one the walk does not read.
func syncedPartitions(parts []*metadata.TablePartitionMetadata, foldTables bool) *partInfo {
	if len(parts) == 0 {
		return nil
	}
	p := &partInfo{kind: parts[0].GetType(), names: make(map[string]bool), count: len(parts)}
	if p.kind == metadata.TablePartitionMetadata_TYPE_UNSPECIFIED {
		return nil
	}
	p.cols, _ = quotedNames(parts[0].GetExpression(), spaceColumn, foldTables)
	single := p.kind == metadata.TablePartitionMetadata_RANGE || p.kind == metadata.TablePartitionMetadata_LIST ||
		len(p.cols) == 1 && (p.kind == metadata.TablePartitionMetadata_RANGE_COLUMNS || p.kind == metadata.TablePartitionMetadata_LIST_COLUMNS)
	if single && rangeLike(p.kind) {
		p.bounds = []int64{}
	}
	if single && listLike(p.kind) {
		p.values = make(map[int64]bool)
	}
	for i, part := range parts {
		k, ok := nameKey(spaceIndex, part.GetName(), foldTables)
		if !ok || strings.HasPrefix(k, nonASCII) || part.GetType() != p.kind {
			return nil
		}
		p.names[k] = true
		if len(part.GetSubpartitions()) > 0 {
			p.subparts = true
		}
		v := strings.TrimSpace(part.GetValue())
		switch {
		case p.bounds != nil:
			if strings.EqualFold(v, "MAXVALUE") && i == len(parts)-1 {
				p.maxvalue = true
				continue
			}
			n, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				p.bounds = nil
				continue
			}
			p.bounds = append(p.bounds, n)
		case p.values != nil:
			for _, item := range strings.Split(v, ",") {
				n, err := strconv.ParseInt(strings.TrimSpace(item), 10, 64)
				if err != nil {
					p.values = nil
					break
				}
				p.values[n] = true
			}
		}
	}
	if len(p.names) != p.count {
		// Two partitions of one key: the walk cannot tell them apart.
		return nil
	}
	return p
}

// syncedOf returns what the sync says of a base table, nil when the
// change may have changed it, the table is not the target database's, or
// a temporary table shadows it.
func (b *builder) syncedOf(r tableRef) *syncedTable {
	st := b.w.synced[relKey(r.db, r.name)]
	if st == nil {
		return nil
	}
	if f, _ := b.l.Lookup(tempKey(r.db, r.name)); f != walk.Absent {
		return nil
	}
	if f, kind := b.l.Lookup(relKey(r.db, r.name)); f != walk.Present || kind != relTable {
		return nil
	}
	return st
}

// referencing returns, sorted by child and name, the foreign keys of other
// synced tables the change has not touched that reference a table.
func (b *builder) referencing(parent walk.Key) []fkUse {
	var out []fkUse
	for child, st := range b.w.synced {
		if child == parent {
			continue
		}
		for name, ref := range st.fks {
			if ref.parent == parent {
				out = append(out, fkUse{child: child, name: name, ref: ref})
			}
		}
	}
	sortFKUses(out)
	return out
}

// fkUse is a foreign key of a child table.
type fkUse struct {
	child walk.Key
	name  string
	ref   fkRef
}

func sortFKUses(us []fkUse) {
	slices.SortFunc(us, func(a, c fkUse) int {
		return cmp.Or(strings.Compare(a.child.Name, c.child.Name), strings.Compare(a.name, c.name))
	})
}
