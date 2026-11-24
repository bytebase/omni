package review

import (
	"slices"
	"strconv"
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
//   - On HASH or KEY partitioning, ADD PARTITION of a definition with a
//     VALUES clause is 1480, and of a name the table has 1517; DROP
//     PARTITION is 1512, REORGANIZE PARTITION 8200, and COALESCE
//     PARTITION 0 1515.
//   - A partitioned table takes no foreign key (1506), which TiDB judges
//     after the key's own text (1239 and 1105 first).
//
// l is the table's partitions when the statement does nothing else and
// the walk knows them (listPartitionsFail), nil otherwise.
func (b *builder) partitionFails(r tableRef, c *ast.AlterTableCmd, kind string, l *partitionList, rng review.Range) bool {
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
	case ast.ATAddConstraint, ast.ATAddIndex:
		if fk := c.Constraint; fk != nil && fk.Type == ast.ConstrForeignKey && len(fk.RefColumns) > 0 {
			for _, kp := range fk.RefKeyParts {
				if kp != nil && kp.Functional {
					return false
				}
			}
			b.fail(rng, nil, review.UnsupportedForm, engineName+" doesn't support a foreign key on partitioned table "+review.Code(r.written)+".")
			return true
		}
		return l != nil && b.addKeyUncovers(r, c.Constraint, l, rng)
	case ast.ATRenameColumn, ast.ATChangeColumn:
		return l != nil && b.renamePartitionColumnFails(r, c, l, rng)
	case ast.ATCoalescePartition:
		if rangeOrList {
			b.fail(rng, nil, review.UnsupportedForm, review.Code("COALESCE PARTITION")+" works only on "+review.Code("HASH")+" or "+review.Code("KEY")+" partitioning, and "+review.Code(r.written)+" has "+review.Code(strings.ToUpper(kind))+".")
			return true
		}
		if c.Number == 0 {
			b.fail(rng, nil, review.UnsupportedForm, review.Code("COALESCE PARTITION")+" must remove at least one partition.")
			return true
		}
	case ast.ATDropPartition:
		if !rangeOrList {
			b.fail(rng, nil, review.UnsupportedForm, review.Code("DROP PARTITION")+" works only on "+review.Code("RANGE")+" or "+review.Code("LIST")+" partitioning, and "+review.Code(r.written)+" has "+review.Code(strings.ToUpper(kind))+".")
			return true
		}
	case ast.ATReorganizePartition:
		if !rangeOrList {
			b.fail(rng, nil, review.UnsupportedForm, engineName+" doesn't support "+review.Code("REORGANIZE PARTITION")+" of "+review.Code(strings.ToUpper(kind))+" partitioning.")
			return true
		}
		if len(c.PartitionNames) == 0 {
			b.fail(rng, nil, review.UnsupportedForm, review.Code("REORGANIZE PARTITION")+" of "+review.Code(r.written)+" has to name the partitions it reorganizes.")
			return true
		}
		if b.partitionDefsFail(c.PartitionDefs, kind, rng) {
			return true
		}
	case ast.ATAddPartition:
		if !rangeOrList {
			for _, pd := range c.PartitionDefs {
				if pd != nil && (pd.BoundKind == ast.PartBoundLessThan || pd.BoundKind == ast.PartBoundIn) {
					b.fail(rng, nil, review.UnsupportedForm, "Partition "+review.Code(ident(pd.Name))+" can't use "+review.Code("VALUES")+" with "+review.Code(strings.ToUpper(kind))+" partitioning.")
					return true
				}
			}
			return !c.IfNotExists && b.partitionTaken(r, c.PartitionDefs, rng)
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
		if b.partitionTaken(r, c.PartitionDefs, rng) {
			return true
		}
	}
	return l != nil && b.listPartitionsFail(r, c, kind, l, rng)
}

// addKeyUncovers fails ADD PRIMARY KEY or UNIQUE, of columns the table
// has, on a partitioned table whose partitioning reads a column the key
// leaves out (8264, verified on v8.5.5); a GLOBAL key is left alone.
func (b *builder) addKeyUncovers(r tableRef, c *ast.Constraint, l *partitionList, rng review.Range) bool {
	if c == nil || (c.Type != ast.ConstrPrimaryKey && c.Type != ast.ConstrUnique) || len(l.cols) == 0 || globalKey(c) {
		return false
	}
	cols, ok := plainKeyParts(c)
	if !ok {
		return false
	}
	for _, col := range cols {
		k, ok := b.key(spaceColumn, col)
		if !ok {
			return false
		}
		if f, _ := b.l.Lookup(colKey(r.db, r.name, k)); f != walk.Present {
			return false
		}
	}
	if !uncoveringPartition(cols, l.cols) {
		return false
	}
	b.partitionKeyFail(r, keyWords(c), rng)
	return true
}

// renamePartitionColumnFails fails RENAME COLUMN, or CHANGE to another
// name, of a column the table's partitioning reads, to a name the table
// doesn't have (3855, verified on v8.5.5).
func (b *builder) renamePartitionColumnFails(r tableRef, c *ast.AlterTableCmd, l *partitionList, rng review.Range) bool {
	to := c.NewName
	if c.Type == ast.ATChangeColumn {
		if c.Column == nil {
			return false
		}
		to = c.Column.Name
	}
	from, ok1 := b.key(spaceColumn, c.Name)
	tk, ok2 := b.key(spaceColumn, to)
	if !ok1 || !ok2 || from == tk || !slices.Contains(l.cols, from) {
		return false
	}
	if f, _ := b.l.Lookup(colKey(r.db, r.name, tk)); f != walk.Absent {
		return false
	}
	b.fail(rng, nil, review.UnsupportedForm, "Column "+review.Code(ident(c.Name))+" is read by the partitioning of "+review.Code(r.written)+", and can't be renamed.")
	return true
}

// partitionTaken fails ADD PARTITION of a name the table has (1517).
func (b *builder) partitionTaken(r tableRef, defs []*ast.PartitionDef, rng review.Range) bool {
	for _, pd := range defs {
		pk, ok := b.key(spacePartition, pd.Name)
		if !ok {
			continue
		}
		if f, _ := b.l.Lookup(partitionKey(r.db, r.name, pk)); f == walk.Present {
			b.fail(rng, nil, review.NameTaken, alreadyExists("partition", ident(pd.Name), r.written))
			return true
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

// A synced table's partitions are one key of its shape, read only when
// Present: their names in order with the bound each has, as
// information_schema.PARTITIONS writes it, and the columns the
// partitioning reads. The sync gives it, any ALTER TABLE of the table
// leaves it Unknown, and a cleared table path Absent, which says nothing.
const spacePartList = "partition list"

func partListKey(db, table string) walk.Key {
	return walk.Key{Space: spacePartList, Parent: tablePath(db, table), Name: shapeName}
}

// partitionList is a table's partitions as the walk keeps them.
type partitionList struct {
	// cols holds the keys of the columns the partitioning reads, nil when
	// the walk does not read them.
	cols  []string
	parts []partEntry
}

// partEntry is one partition: its name's key and its bound, "" for none.
type partEntry struct {
	name, bound string
}

// syncedPartList reads a synced table's partitions, false when the walk
// cannot read one of them.
func syncedPartList(parts []*metadata.TablePartitionMetadata) (string, bool) {
	if len(parts) == 0 {
		return "", false
	}
	var l partitionList
	if expr := parts[0].GetExpression(); expr != "" {
		cols, ok := checkColumns(expr)
		if !ok {
			return "", false
		}
		l.cols = cols
	}
	for _, p := range parts {
		name, ok := plainName(spacePartition, p.GetName())
		if !ok || strings.ContainsAny(p.GetValue(), "\x1d\x1e\x1f") || len(p.GetSubpartitions()) > 0 {
			return "", false
		}
		l.parts = append(l.parts, partEntry{name: name, bound: p.GetValue()})
	}
	return l.encode(), true
}

func (l partitionList) encode() string {
	parts := make([]string, 0, len(l.parts)+1)
	parts = append(parts, strings.Join(l.cols, "\x1f"))
	for _, p := range l.parts {
		parts = append(parts, p.name+"\x1e"+p.bound)
	}
	return strings.Join(parts, "\x1d")
}

func decodePartList(kind string) partitionList {
	fields := strings.Split(kind, "\x1d")
	var l partitionList
	if fields[0] != "" {
		l.cols = strings.Split(fields[0], "\x1f")
	}
	for _, f := range fields[1:] {
		name, bound, _ := strings.Cut(f, "\x1e")
		l.parts = append(l.parts, partEntry{name: name, bound: bound})
	}
	return l
}

// partListOf returns a table's partitions before the statement, nil when
// the walk does not know them.
func (b *builder) partListOf(r tableRef) *partitionList {
	f, kind := b.l.Lookup(partListKey(r.db, r.name))
	if f != walk.Present {
		return nil
	}
	l := decodePartList(kind)
	return &l
}

// index returns the position of a partition, -1 when the table has none of
// the name.
func (l *partitionList) index(name string) int {
	for i, p := range l.parts {
		if p.name == name {
			return i
		}
	}
	return -1
}

// unknownPartition returns the first name of names the table has no
// partition of, "" when it has them all.
func (b *builder) unknownPartition(l *partitionList, names []string) string {
	for _, n := range names {
		k, ok := b.key(spacePartition, n)
		if !ok {
			return ""
		}
		if l.index(k) < 0 {
			return n
		}
	}
	return ""
}

// partitionBound returns a definition's bound as information_schema.PARTITIONS
// writes it, false for one the walk does not read: an integer, a plain
// string, MAXVALUE, or VALUES IN of those.
func partitionBound(pd *ast.PartitionDef) (string, bool) {
	if pd == nil {
		return "", false
	}
	switch pd.BoundKind {
	case ast.PartBoundLessThan:
		return boundConstant(pd.Values)
	case ast.PartBoundIn:
		list, ok := pd.Values.(*ast.List)
		if !ok {
			return "", false
		}
		vals := make([]string, 0, len(list.Items))
		for _, item := range list.Items {
			v, ok := boundConstant(item)
			if !ok || v == "MAXVALUE" {
				return "", false
			}
			vals = append(vals, v)
		}
		return strings.Join(vals, ","), true
	}
	return "", false
}

func boundConstant(n ast.Node) (string, bool) {
	switch x := n.(type) {
	case *ast.IntLit:
		return strconv.FormatInt(x.Value, 10), true
	case *ast.UnaryExpr:
		if x.Op == ast.UnaryMinus {
			if lit, ok := x.Operand.(*ast.IntLit); ok && lit.Value > 0 {
				return "-" + strconv.FormatInt(lit.Value, 10), true
			}
		}
	case *ast.ParenExpr:
		return boundConstant(x.Expr)
	case *ast.String:
		if x.Str == "MAXVALUE" {
			return x.Str, true
		}
	case *ast.StringLit:
		if x.Charset == "" && !x.Concatenated && !strings.ContainsAny(x.Value, `'\`+"\x1d\x1e\x1f") {
			return "'" + x.Value + "'", true
		}
	}
	return "", false
}

// rangeFault returns where RANGE bounds of one column, in order, stop
// rising, -1 where they rise as far as the walk compares them (integers
// and MAXVALUE), and whether the bound before it is MAXVALUE.
func rangeFault(bounds []string) (int, bool) {
	for i := 1; i < len(bounds); i++ {
		prev, cur := bounds[i-1], bounds[i]
		if prev == "MAXVALUE" {
			return i, true
		}
		if cur == "MAXVALUE" {
			continue
		}
		a, err1 := strconv.ParseInt(prev, 10, 64)
		c, err2 := strconv.ParseInt(cur, 10, 64)
		if err1 == nil && err2 == nil && c <= a {
			return i, false
		}
	}
	return -1, false
}

// listFault returns the first constant LIST bounds of one column name
// twice, "" for none. A bound the walk cannot split is passed over.
func listFault(bounds []string) string {
	seen := make(map[string]bool)
	for _, b := range bounds {
		if strings.Contains(b, `\`) {
			continue
		}
		vals := strings.Split(b, ",")
		for _, v := range vals {
			if strings.Count(v, "'")%2 != 0 {
				vals = nil
				break
			}
		}
		for _, v := range vals {
			v = strings.TrimSpace(v)
			if seen[v] {
				return v
			}
			seen[v] = true
		}
	}
	return ""
}

// sequenceFails fails partitions that leave RANGE bounds of one column
// out of order (1493) or LIST values of one column named twice (1495),
// verified on v8.5.5 for CREATE TABLE, PARTITION BY, ADD and REORGANIZE
// PARTITION. Bounds before from are the table's own, which ADD PARTITION
// follows: MAXVALUE among them is 1493 too, and the walk claims a MAXVALUE
// anywhere else (1481) only in a whole clause (from 0).
func (b *builder) sequenceFails(kind string, single bool, bounds []string, from int, rng review.Range) bool {
	if !single {
		return false
	}
	switch kind {
	case partRange, partRangeCols:
		i, afterMax := rangeFault(bounds)
		switch {
		case i < 0 || i < from:
			return false
		case afterMax && from == 0:
			b.fail(rng, nil, review.UnsupportedForm, review.Code("MAXVALUE")+" can only bound the last partition.")
		case afterMax && i != from:
			return false
		default:
			b.fail(rng, nil, review.UnsupportedForm, "Each partition's "+review.Code("VALUES LESS THAN")+" bound must be greater than the bound of the partition before it.")
		}
		return true
	case partList, partListCols:
		if v := listFault(bounds); v != "" {
			b.fail(rng, nil, review.UnsupportedForm, "The partitions' "+review.Code("VALUES IN")+" name "+review.Code(v)+" twice.")
			return true
		}
	}
	return false
}

// listPartitionsFail fails an ALTER TABLE, alone in its statement, of a
// table whose partitions the walk knows (verified on v8.5.5, each in the
// partitioned@ fixtures of tidb/catalog's simulate.tsv):
//
//   - DROP PARTITION of a partition the table doesn't have (1507), or of
//     every one (1508); COALESCE PARTITION of every one (1508);
//   - TRUNCATE PARTITION of a partition the table doesn't have (1735);
//   - REORGANIZE PARTITION of RANGE partitions that are not adjacent
//     (8200), or into a name a partition it keeps holds (1517);
//   - ADD or REORGANIZE PARTITION leaving RANGE bounds out of order (1493)
//     or LIST values twice (1495).
func (b *builder) listPartitionsFail(r tableRef, c *ast.AlterTableCmd, kind string, l *partitionList, rng review.Range) bool {
	// Two definitions of one name are 1517 first (checkPartitionDefs).
	names := make(map[string]bool, len(c.PartitionDefs))
	for _, pd := range c.PartitionDefs {
		k, ok := b.key(spacePartition, pd.Name)
		if !ok || names[k] {
			return false
		}
		names[k] = true
	}
	single := kind == partRange || kind == partList || len(l.cols) == 1
	bounds := func(parts []partEntry) []string {
		out := make([]string, len(parts))
		for i, p := range parts {
			out[i] = p.bound
		}
		return out
	}
	switch c.Type {
	case ast.ATDropPartition:
		if c.IfExists {
			return false
		}
		if n := b.unknownPartition(l, c.PartitionNames); n != "" {
			b.fail(rng, nil, review.ObjectNotFound, notFound("partition", ident(n), r.written))
			return true
		}
		dropped := make(map[string]bool)
		for _, n := range c.PartitionNames {
			k, _ := b.key(spacePartition, n)
			dropped[k] = true
		}
		if len(dropped) >= len(l.parts) {
			b.fail(rng, nil, review.UnsupportedForm, "This statement removes every partition of "+review.Code(r.written)+"; drop the table instead.")
			return true
		}
	case ast.ATCoalescePartition:
		if c.Number >= len(l.parts) {
			b.fail(rng, nil, review.UnsupportedForm, "This statement removes every partition of "+review.Code(r.written)+"; drop the table instead.")
			return true
		}
	case ast.ATTruncatePartition:
		if c.AllPartitions {
			return false
		}
		if n := b.unknownPartition(l, c.PartitionNames); n != "" {
			b.fail(rng, nil, review.ObjectNotFound, notFound("partition", ident(n), r.written))
			return true
		}
	case ast.ATAddPartition:
		if c.IfNotExists {
			return false
		}
		all := bounds(l.parts)
		for _, pd := range c.PartitionDefs {
			v, ok := partitionBound(pd)
			if !ok {
				return false
			}
			all = append(all, v)
		}
		return b.sequenceFails(kind, single, all, len(l.parts), rng)
	case ast.ATReorganizePartition:
		if b.unknownPartition(l, c.PartitionNames) != "" {
			return false
		}
		first, last := len(l.parts), -1
		gone := make(map[string]bool)
		for _, n := range c.PartitionNames {
			k, _ := b.key(spacePartition, n)
			i := l.index(k)
			first, last = min(first, i), max(last, i)
			gone[k] = true
		}
		if (kind == partRange || kind == partRangeCols) && last-first+1 != len(gone) {
			b.fail(rng, nil, review.UnsupportedForm, engineName+" doesn't support "+review.Code("REORGANIZE PARTITION")+" of partitions that aren't adjacent.")
			return true
		}
		var kept []partEntry
		for _, p := range l.parts {
			if !gone[p.name] {
				kept = append(kept, p)
			}
		}
		for _, pd := range c.PartitionDefs {
			k, ok := b.key(spacePartition, pd.Name)
			if !ok {
				return false
			}
			for _, p := range kept {
				if p.name == k {
					b.fail(rng, nil, review.NameTaken, alreadyExists("partition", ident(pd.Name), r.written))
					return true
				}
			}
		}
		var after []string
		for i, p := range l.parts {
			if i == first {
				for _, pd := range c.PartitionDefs {
					v, ok := partitionBound(pd)
					if !ok {
						return false
					}
					after = append(after, v)
				}
			}
			if !gone[p.name] {
				after = append(after, p.bound)
			}
		}
		i, afterMax := rangeFault(after)
		if (kind == partRange || kind == partRangeCols) && (i < 0 || afterMax) {
			return false
		}
		return b.sequenceFails(kind, single, after, 1, rng)
	}
	return false
}

// partitionClauseFails fails a PARTITION BY clause, of CREATE TABLE or
// ALTER TABLE, that TiDB refuses for its own text and the types of the
// columns it reads (colType, false when not known), each verified on
// v8.5.5 in tidb/catalog's simulate.tsv: two partitions of one name
// (1517); a lone column of a string type for RANGE, LIST or HASH, or of a
// TEXT or BLOB type for KEY (1659); RANGE bounds of one column out of
// order (1481, 1493), or LIST values of one column named twice (1495).
// TiDB's parser refuses the rest of what the clause alone settles (1480,
// 1484, 1492, 1504), and so does omni's.
func (b *builder) partitionClauseFails(pc *ast.PartitionClause, colType func(string) (columnType, bool), rng review.Range) bool {
	kind, ok := clausePartitioning(pc)
	if !ok || pc == nil || pc.Interval != nil {
		return false
	}
	names := make(map[string]bool, len(pc.Partitions))
	for _, pd := range pc.Partitions {
		if pd == nil || pd.Name == "" {
			continue
		}
		k := strings.ToLower(pd.Name)
		if names[k] {
			b.fail(rng, nil, review.DuplicateName, statesTwice("declares", "partition", ident(pd.Name)))
			return true
		}
		names[k] = true
	}
	switch kind {
	case partKey:
		for _, c := range pc.Columns {
			if t, ok := colType(c); ok && isTextOrBlob(t.base) {
				b.fail(rng, nil, review.UnsupportedForm, "Partition column "+review.Code(ident(c))+" is "+review.Code(t.base)+", which "+review.Code("KEY")+" partitioning can't read.")
				return true
			}
		}
	case partRange, partList, partHash:
		if c, ok := unwrapParens(pc.Expr).(*ast.ColumnRef); ok && c.Table == "" && !c.Star {
			if t, ok := colType(c.Column); ok && isStringType(t.base) {
				b.fail(rng, nil, review.UnsupportedForm, "Partition column "+review.Code(ident(c.Column))+" is "+review.Code(t.base)+", and the partitioning function must return an integer.")
				return true
			}
		}
	}
	var bounds []string
	for _, pd := range pc.Partitions {
		v, ok := partitionBound(pd)
		if !ok {
			return false
		}
		bounds = append(bounds, v)
	}
	return b.sequenceFails(kind, len(pc.Columns) <= 1, bounds, 0, rng)
}

// isStringType reports a character or binary string type, which RANGE,
// LIST and HASH can't read as a lone column (1659).
func isStringType(t string) bool {
	switch t {
	case "char", "varchar", "binary", "varbinary":
		return true
	}
	return isTextOrBlob(t)
}

// isTextOrBlob reports a TEXT or BLOB type, which no partitioning reads
// (1659).
func isTextOrBlob(t string) bool {
	switch t {
	case "tinytext", "text", "mediumtext", "longtext", "tinyblob", "blob", "mediumblob", "longblob":
		return true
	}
	return false
}

// exchangeFails fails EXCHANGE PARTITION with a table that is partitioned
// itself (1732), and then of a partition the table doesn't have (1735),
// in the order v8.5.5 checks them. The walk claims a missing partition
// only of a table it knows is not partitioned, and when the statement
// does nothing else.
func (b *builder) exchangeFails(r tableRef, c *ast.AlterTableCmd, ex tableRef, rng review.Range) bool {
	switch b.partitioning(ex) {
	case "", partNone:
	default:
		b.fail(rng, nil, review.UnsupportedForm, engineName+" can't exchange a partition of "+review.Code(r.written)+" with "+review.Code(ex.written)+", which is partitioned.")
		return true
	}
	if b.partitioning(ex) != partNone {
		return true
	}
	if v, ok := b.s.node.(*ast.AlterTableStmt); !ok || len(v.Commands) != 1 {
		return true
	}
	if l := b.partListOf(r); l != nil && b.unknownPartition(l, []string{c.Name}) != "" {
		b.fail(rng, nil, review.ObjectNotFound, notFound("partition", ident(c.Name), r.written))
	}
	return true
}
