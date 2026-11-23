package review

import (
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
	"github.com/bytebase/omni/tidb/ast"
)

// A table's partitioning is one key of its own, under the table's path:
// Present with the partitioning's kind, partNone for a table that is not
// partitioned. Only Present is read; an Absent key, as a cleared table path
// leaves it, says nothing, so a table the walk cannot read the partitioning
// of is never taken for one it can. A partition is a key of its table's
// path too, read only when Present: TiDB folds partition names (verified:
// ADD PARTITION (PARTITION P1 ...) beside p1 fails with 1517).
const (
	spacePartitioning = "partitioning"
	spacePartition    = "partition"
	partitioningName  = "by"
)

// Partitioning kinds. LINEAR HASH and LINEAR KEY are HASH and KEY to every
// check the walk makes.
const (
	partNone      = "none"
	partRange     = "range"
	partRangeCols = "range columns"
	partList      = "list"
	partListCols  = "list columns"
	partHash      = "hash"
	partKey       = "key"
)

func partitioningKey(db, table string) walk.Key {
	return walk.Key{Space: spacePartitioning, Parent: tablePath(db, table), Name: partitioningName}
}

func partitionKey(db, table, name string) walk.Key {
	return walk.Key{Space: spacePartition, Parent: tablePath(db, table), Name: name}
}

func partitionScope(db, table string) walk.Scope {
	return walk.Scope{Space: spacePartition, Parent: tablePath(db, table)}
}

// syncedPartitioning is the kind of a synced table's partitioning, from the
// partitions Bytebase's sync reads from information_schema.PARTITIONS: none
// for a table it lists none of, false for a method the walk does not know.
func syncedPartitioning(parts []*metadata.TablePartitionMetadata) (string, bool) {
	if len(parts) == 0 {
		return partNone, true
	}
	switch parts[0].GetType() {
	case metadata.TablePartitionMetadata_RANGE:
		return partRange, true
	case metadata.TablePartitionMetadata_RANGE_COLUMNS:
		return partRangeCols, true
	case metadata.TablePartitionMetadata_LIST:
		return partList, true
	case metadata.TablePartitionMetadata_LIST_COLUMNS:
		return partListCols, true
	case metadata.TablePartitionMetadata_HASH, metadata.TablePartitionMetadata_LINEAR_HASH:
		return partHash, true
	case metadata.TablePartitionMetadata_KEY, metadata.TablePartitionMetadata_LINEAR_KEY:
		return partKey, true
	}
	return "", false
}

// clausePartitioning is the kind a PARTITION BY clause gives a table, false
// when the walk does not follow it: SYSTEM_TIME, which TiDB reads and does
// not apply, and subpartitioning, which it does not support.
func clausePartitioning(pc *ast.PartitionClause) (string, bool) {
	if pc == nil {
		return partNone, true
	}
	if pc.SubPartExpr != nil || len(pc.SubPartColumns) > 0 || pc.NumSubParts > 0 {
		return "", false
	}
	for _, pd := range pc.Partitions {
		if pd == nil || len(pd.SubPartitions) > 0 {
			return "", false
		}
	}
	switch pc.Type {
	case ast.PartitionRange:
		if len(pc.Columns) > 0 {
			return partRangeCols, true
		}
		return partRange, true
	case ast.PartitionList:
		if len(pc.Columns) > 0 {
			return partListCols, true
		}
		return partList, true
	case ast.PartitionHash:
		return partHash, true
	case ast.PartitionKey:
		return partKey, true
	}
	return "", false
}

// setPartitioning records the partitioning a CREATE TABLE or a PARTITION BY
// gives a table, and the partitions it names. Partitions TiDB names itself
// (PARTITIONS n, INTERVAL) are a May-effect.
func (b *builder) setPartitioning(r tableRef, pc *ast.PartitionClause) {
	k := partitioningKey(r.db, r.name)
	scope := partitionScope(r.db, r.name)
	kind, ok := clausePartitioning(pc)
	if !ok {
		b.Set(k, walk.Unknown, "")
		b.Forget(scope)
		return
	}
	b.Set(k, walk.Present, kind)
	b.Clear(scope)
	if pc == nil {
		return
	}
	if len(pc.Partitions) == 0 || pc.Interval != nil {
		b.MayAdd(scope)
	}
	for _, pd := range pc.Partitions {
		if pk, ok := b.key(spacePartition, pd.Name); ok {
			b.Set(partitionKey(r.db, r.name, pk), walk.Present, "")
		}
	}
}

// partitioning returns the kind of a table's partitioning before the
// statement, "" when the walk does not know it.
func (b *builder) partitioning(r tableRef) string {
	if f, kind := b.l.Lookup(partitioningKey(r.db, r.name)); f == walk.Present {
		return kind
	}
	return ""
}

func notPartitioned(table string) string {
	return "Table " + review.Code(table) + " is not partitioned."
}

// partitionFails fails an ALTER TABLE subcommand the table's partitioning
// refutes, kind as the table had it before the statement (verified on
// v8.5.5, each against a table of the kind):
//
//   - On a table that is not partitioned, every partition-management
//     action is 1505: ADD, DROP, TRUNCATE, ANALYZE, COALESCE, EXCHANGE and
//     REORGANIZE PARTITION, REMOVE PARTITIONING, FIRST and LAST PARTITION
//     LESS THAN, and PARTITION p ATTRIBUTES; IF [NOT] EXISTS does not make
//     it a note. COMPACT PARTITION is 1105. ALTER TABLE t PARTITION p, with
//     no option, runs; CHECK, OPTIMIZE, REBUILD, REPAIR, DISCARD and IMPORT
//     PARTITION fail with 8200 whatever the table, which walk_validate.go
//     reports.
//   - On a partitioned table, CACHE is 8242.
//   - On RANGE or LIST partitioning, ADD PARTITION with no definitions is
//     1492, COALESCE PARTITION 1509, and REORGANIZE PARTITION naming none
//     8200. A definition ADD or REORGANIZE PARTITION gives must carry the
//     partitioning's own VALUES clause: none is 1479, the other one 1480,
//     and VALUES LESS THAN of more than one value on RANGE (not RANGE
//     COLUMNS) is 1657, as TiDB's PartitionDefinitionClause.Validate reads
//     each. ADD PARTITION of a name the table has is 1517, unless written
//     IF NOT EXISTS.
func (b *builder) partitionFails(r tableRef, c *ast.AlterTableCmd, kind string, rng review.Range) bool {
	switch kind {
	case "":
		return false
	case partNone:
		switch c.Type {
		case ast.ATAddPartition, ast.ATDropPartition, ast.ATTruncatePartition, ast.ATAnalyzePartition,
			ast.ATCoalescePartition, ast.ATExchangePartition, ast.ATReorganizePartition, ast.ATRemovePartitioning,
			ast.ATAddLastPartition, ast.ATDropFirstPartition, ast.ATPartitionAttributes:
		case ast.ATCompact:
			if len(c.PartitionNames) == 0 {
				return false
			}
		default:
			return false
		}
		b.fail(rng, nil, review.UnsupportedForm, notPartitioned(r.written))
		return true
	}
	rangeOrList := kind == partRange || kind == partRangeCols || kind == partList || kind == partListCols
	// A HISTORY or CURRENT partition is refused first (4113), which
	// checkPartitionDefs reports.
	for _, pd := range c.PartitionDefs {
		if pd != nil && pd.BoundKind == ast.PartBoundHistory {
			return false
		}
	}
	switch c.Type {
	case ast.ATCache:
		b.fail(rng, nil, review.UnsupportedForm, engineName+" can't cache "+review.Code(r.written)+", which is partitioned.")
		return true
	case ast.ATCoalescePartition:
		if rangeOrList {
			b.fail(rng, nil, review.UnsupportedForm, review.Code("COALESCE PARTITION")+" works only on "+review.Code("HASH")+" or "+review.Code("KEY")+" partitioning, and "+review.Code(r.written)+" has "+review.Code(strings.ToUpper(kind))+".")
			return true
		}
	case ast.ATReorganizePartition:
		if rangeOrList && len(c.PartitionNames) == 0 {
			b.fail(rng, nil, review.UnsupportedForm, review.Code("REORGANIZE PARTITION")+" of "+review.Code(r.written)+" has to name the partitions it reorganizes.")
			return true
		}
		return rangeOrList && b.partitionDefsFail(c.PartitionDefs, kind, rng)
	case ast.ATAddPartition:
		if !rangeOrList {
			return false
		}
		if len(c.PartitionDefs) == 0 {
			b.fail(rng, nil, review.UnsupportedForm, "Each partition of "+review.Code(r.written)+" has to be defined for "+review.Code(strings.ToUpper(kind))+" partitioning.")
			return true
		}
		if b.partitionDefsFail(c.PartitionDefs, kind, rng) {
			return true
		}
		if c.IfNotExists {
			return false
		}
		for _, pd := range c.PartitionDefs {
			pk, ok := b.key(spacePartition, pd.Name)
			if !ok {
				continue
			}
			if f, _ := b.l.Lookup(partitionKey(r.db, r.name, pk)); f == walk.Present {
				b.fail(rng, nil, review.NameTaken, alreadyExists("partition", ident(pd.Name), r.written))
				return true
			}
		}
	}
	return false
}

// partitionDefsFail fails a partition definition whose VALUES clause RANGE
// or LIST partitioning refutes (1479, 1480, 1657).
func (b *builder) partitionDefsFail(defs []*ast.PartitionDef, kind string, rng review.Range) bool {
	isRange := kind == partRange || kind == partRangeCols
	want, other := ast.PartBoundLessThan, ast.PartBoundIn
	clause, otherClause := "VALUES LESS THAN", "VALUES IN"
	if !isRange {
		want, other = ast.PartBoundIn, ast.PartBoundLessThan
		clause, otherClause = "VALUES IN", "VALUES LESS THAN"
	}
	for _, pd := range defs {
		if pd == nil {
			continue
		}
		switch pd.BoundKind {
		case ast.PartBoundNone:
			b.fail(rng, nil, review.UnsupportedForm, "Partition "+review.Code(ident(pd.Name))+" needs "+review.Code(clause)+" for "+review.Code(strings.ToUpper(kind))+" partitioning.")
			return true
		case other:
			b.fail(rng, nil, review.UnsupportedForm, "Partition "+review.Code(ident(pd.Name))+" can't use "+review.Code(otherClause)+" with "+review.Code(strings.ToUpper(kind))+" partitioning.")
			return true
		case want:
			if kind == partRange {
				if l, ok := pd.Values.(*ast.List); ok && len(l.Items) > 1 {
					b.fail(rng, nil, review.UnsupportedForm, "Partition "+review.Code(ident(pd.Name))+" can't have more than one "+review.Code("VALUES LESS THAN")+" value for "+review.Code("RANGE")+" partitioning.")
					return true
				}
			}
		}
	}
	return false
}

// partitionEffects records what an ALTER TABLE subcommand does to the
// table's partitioning and partitions.
func (b *builder) partitionEffects(r tableRef, c *ast.AlterTableCmd) {
	scope := partitionScope(r.db, r.name)
	named := func(names []string, f walk.Fact) {
		for _, n := range names {
			if pk, ok := b.key(spacePartition, n); ok {
				b.Set(partitionKey(r.db, r.name, pk), f, "")
			}
		}
	}
	defs := func(defs []*ast.PartitionDef) {
		for _, pd := range defs {
			if pd != nil {
				named([]string{pd.Name}, walk.Present)
			}
		}
	}
	switch c.Type {
	case ast.ATPartitionBy:
		b.setPartitioning(r, c.PartitionBy)
	case ast.ATRemovePartitioning:
		b.Set(partitioningKey(r.db, r.name), walk.Present, partNone)
		b.Clear(scope)
	case ast.ATAddPartition:
		if len(c.PartitionDefs) == 0 || c.IfNotExists {
			b.MayAdd(scope)
		}
		if !c.IfNotExists {
			defs(c.PartitionDefs)
		}
	case ast.ATDropPartition:
		if c.IfExists {
			b.MayRemove(scope)
			return
		}
		named(c.PartitionNames, walk.Absent)
	case ast.ATReorganizePartition:
		if len(c.PartitionNames) == 0 {
			b.Forget(scope)
			return
		}
		named(c.PartitionNames, walk.Absent)
		defs(c.PartitionDefs)
	case ast.ATCoalescePartition, ast.ATAddLastPartition, ast.ATDropFirstPartition,
		ast.ATReorganizeLastPartition, ast.ATReorganizeFirstPartition:
		b.Forget(scope)
	}
}

// partitionChange reports whether a subcommand may change a table's
// partitioning or partitions.
func partitionChange(t ast.AlterTableCmdType) bool {
	switch t {
	case ast.ATPartitionBy, ast.ATRemovePartitioning, ast.ATAddPartition, ast.ATDropPartition,
		ast.ATReorganizePartition, ast.ATCoalescePartition, ast.ATAddLastPartition, ast.ATDropFirstPartition,
		ast.ATReorganizeLastPartition, ast.ATReorganizeFirstPartition:
		return true
	}
	return false
}
