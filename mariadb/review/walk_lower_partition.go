package review

import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// partitionInt reads a partition bound written as an integer literal, with
// a minus or not.
func partitionInt(n ast.Node) (int64, bool) {
	switch x := n.(type) {
	case *ast.IntLit:
		if x.Text != "" {
			return 0, false
		}
		return x.Value, true
	case *ast.UnaryExpr:
		if x.Op != ast.UnaryMinus {
			return 0, false
		}
		v, ok := partitionInt(x.Operand)
		return -v, ok
	case *ast.ParenExpr:
		return partitionInt(x.Expr)
	}
	return 0, false
}

// isMaxValue reports a VALUES LESS THAN MAXVALUE.
func isMaxValue(n ast.Node) bool {
	s, ok := n.(*ast.String)
	return ok && strings.EqualFold(s.Str, "MAXVALUE")
}

// listInts reads the integers a VALUES IN list names, false when one is
// anything else.
func listInts(n ast.Node) ([]int64, bool) {
	l, ok := n.(*ast.List)
	if !ok {
		return nil, false
	}
	out := make([]int64, 0, len(l.Items))
	for _, item := range l.Items {
		v, ok := partitionInt(item)
		if !ok {
			return nil, false
		}
		out = append(out, v)
	}
	return out, true
}

// partitionDefsFail checks partition definitions a statement adds to a
// partitioning of a kind, after the partitions it already has, as
// MariaDB checks them (each verified on mariadb:11.8.8): a name twice,
// without case (1517); a RANGE bound no greater than the one before it
// (1493), or any after MAXVALUE (1481); and a LIST value twice (1495).
// bounds and values are nil when the walk does not know the partitions'
// own; names holds their keys.
func (b *builder) partitionDefsFail(kind metadata.TablePartitionMetadata_Type, defs []*ast.PartitionDef, names map[string]bool, bounds []int64, maxvalue bool, values map[int64]bool, keys []walk.Key) bool {
	seen := make(map[string]bool, len(names)+len(defs))
	for n := range names {
		seen[n] = true
	}
	for _, d := range defs {
		if d == nil {
			continue
		}
		k, ok := b.key(spaceIndex, d.Name)
		if !ok || strings.HasPrefix(k, nonASCII) {
			return false
		}
		if seen[k] {
			b.fail(rangeOf(d.Loc), keys, review.DuplicateName, "Partition "+review.Code(ident(d.Name))+" already exists.")
			return true
		}
		seen[k] = true
	}
	switch {
	case rangeLike(kind) && bounds != nil:
		last := bounds
		for _, d := range defs {
			if d == nil || d.ValuesIn || d.Values == nil {
				return false
			}
			if maxvalue {
				b.fail(rangeOf(d.Loc), keys, review.UnsupportedForm, "Partition "+review.Code(ident(d.Name))+" comes after the partition of "+review.Code("MAXVALUE")+", which must be the last.")
				return true
			}
			if isMaxValue(d.Values) {
				maxvalue = true
				continue
			}
			v, ok := partitionInt(d.Values)
			if !ok {
				return false
			}
			if len(last) > 0 && v <= last[len(last)-1] {
				b.fail(rangeOf(d.Loc), keys, review.UnsupportedForm, "Partition "+review.Code(ident(d.Name))+"'s bound must be greater than the bound of the partition before it.")
				return true
			}
			last = append(last[:len(last):len(last)], v)
		}
	case listLike(kind) && values != nil:
		all := make(map[int64]bool, len(values))
		for v := range values {
			all[v] = true
		}
		for _, d := range defs {
			if d == nil || !d.ValuesIn {
				return false
			}
			vs, ok := listInts(d.Values)
			if !ok {
				return false
			}
			for _, v := range vs {
				if all[v] {
					b.fail(rangeOf(d.Loc), keys, review.UnsupportedForm, "Partition "+review.Code(ident(d.Name))+" lists "+review.Code(strconv.FormatInt(v, 10))+", which another partition lists.")
					return true
				}
				all[v] = true
			}
		}
	}
	return false
}

// partitionCmdOnly returns an ALTER TABLE's one partition operation, nil
// when it does anything else but ALGORITHM and LOCK.
func partitionCmdOnly(v *ast.AlterTableStmt) *ast.AlterTableCmd {
	var cmd *ast.AlterTableCmd
	for _, c := range v.Commands {
		switch {
		case c == nil:
			return nil
		case c.Type == ast.ATAlgorithm || c.Type == ast.ATLock:
		case cmd == nil && isPartitionMgmtCmd(c.Type):
			cmd = c
		default:
			return nil
		}
	}
	return cmd
}

// unknownPartition is the text for a partition a table doesn't have
// (1507, 1735).
func unknownPartition(name, table string) string {
	return "Partition " + review.Code(ident(name)) + " doesn't exist in " + review.Code(table) + "."
}

// alterPartitionFails fails an ALTER TABLE of a synced partitioned table
// the change has not touched whose partition operation MariaDB refuses
// (each verified on mariadb:11.8.8):
//
//   - DROP PARTITION of HASH or KEY partitioning (1512), of as many
//     partitions as the table has (1508), or of one it lacks (1507);
//   - TRUNCATE, EXCHANGE, REBUILD, or REORGANIZE PARTITION of one it lacks
//     (1735 for the first two, 1507), and EXCHANGE with a partitioned table
//     (1732), on a table without subpartitions, whose names they take too;
//   - ADD PARTITION of a name it has (1517), of a bound of the wrong kind
//     (1480), and partitionDefsFail's;
//   - COALESCE PARTITION of none (1515) or of all (1508), and a REORGANIZE
//     PARTITION of HASH or KEY partitioning that changes how many there
//     are (1510);
//   - DISCARD or IMPORT TABLESPACE of the whole table (1031).
func (b *builder) alterPartitionFails(v *ast.AlterTableStmt, r tableRef) bool {
	st := b.syncedOf(r)
	if st == nil || st.part == nil {
		return false
	}
	p := st.part
	rel := relKey(r.db, r.name)
	keys := []walk.Key{rel}
	for _, c := range v.Commands {
		if c != nil && (c.Type == ast.ATDiscardTablespace || c.Type == ast.ATImportTablespace) {
			b.fail(rangeOf(c.Loc), keys, review.UnsupportedForm, "Table "+review.Code(r.written)+" is partitioned, and "+review.Code("DISCARD")+" and "+review.Code("IMPORT TABLESPACE")+" take only a partition of it.")
			return true
		}
	}
	c := partitionCmdOnly(v)
	if c == nil || c.IfExists || c.IfNotExists {
		return false
	}
	rng := rangeOf(c.Loc)
	// missing returns the first name of a list the table lacks.
	missing := func(names []string) (string, bool) {
		for _, n := range names {
			k, ok := b.key(spaceIndex, n)
			if !ok || strings.HasPrefix(k, nonASCII) {
				return "", false
			}
			if !p.names[k] {
				return n, true
			}
		}
		return "", false
	}
	switch c.Type {
	case ast.ATDropPartition:
		switch {
		case hashLike(p.kind):
			b.fail(rng, keys, review.UnsupportedForm, review.Code("DROP PARTITION")+" works only on "+review.Code("RANGE")+" or "+review.Code("LIST")+" partitioning.")
			return true
		case len(c.PartitionNames) >= p.count:
			b.fail(rng, keys, review.UnsupportedForm, "Cannot remove all partitions of "+review.Code(r.written)+"; drop the table instead.")
			return true
		}
		if n, ok := missing(c.PartitionNames); ok && !p.subparts {
			b.fail(rng, keys, review.ObjectNotFound, unknownPartition(n, r.written))
			return true
		}
	case ast.ATTruncatePartition, ast.ATRebuildPartition:
		if c.AllPartitions || p.subparts {
			return false
		}
		if n, ok := missing(c.PartitionNames); ok {
			b.fail(rng, keys, review.ObjectNotFound, unknownPartition(n, r.written))
			return true
		}
	case ast.ATReorganizePartition:
		if len(c.PartitionNames) == 0 {
			return false
		}
		if n, ok := missing(c.PartitionNames); ok && !p.subparts {
			b.fail(rng, keys, review.ObjectNotFound, unknownPartition(n, r.written))
			return true
		}
		if hashLike(p.kind) && len(c.PartitionDefs) > 0 && len(c.PartitionDefs) != len(c.PartitionNames) {
			b.fail(rng, keys, review.UnsupportedForm, review.Code("REORGANIZE PARTITION")+" of "+review.Code("HASH")+" or "+review.Code("KEY")+" partitioning must keep the number of partitions.")
			return true
		}
	case ast.ATExchangePartition:
		if p.subparts {
			return false
		}
		if n, ok := missing([]string{c.Name}); ok {
			b.fail(rng, keys, review.ObjectNotFound, unknownPartition(n, r.written))
			return true
		}
		if c.ExchangeTable != nil {
			if x, st := b.ref(c.ExchangeTable); st == refOK {
				if xs := b.syncedOf(x); xs != nil && xs.part != nil {
					b.fail(rng, []walk.Key{rel, relKey(x.db, x.name)}, review.UnsupportedForm, "Table "+review.Code(x.written)+" is partitioned, and a partition can be exchanged only with a table that isn't.")
					return true
				}
			}
		}
	case ast.ATCoalescePartition:
		if !hashLike(p.kind) {
			return false
		}
		switch {
		case c.Number == 0:
			b.fail(rng, keys, review.UnsupportedForm, review.Code("COALESCE PARTITION")+" must remove at least one partition.")
			return true
		case c.Number >= p.count:
			b.fail(rng, keys, review.UnsupportedForm, "Cannot remove all partitions of "+review.Code(r.written)+"; drop the table instead.")
			return true
		}
	case ast.ATAddPartition:
		for _, d := range c.PartitionDefs {
			if d == nil || d.Values == nil {
				continue
			}
			wrong := ""
			switch {
			case rangeLike(p.kind) && d.ValuesIn:
				wrong = "VALUES IN"
			case listLike(p.kind) && !d.ValuesIn:
				wrong = "VALUES LESS THAN"
			case hashLike(p.kind) && d.ValuesIn:
				wrong = "VALUES IN"
			case hashLike(p.kind):
				wrong = "VALUES LESS THAN"
			}
			if wrong != "" {
				b.fail(rangeOf(d.Loc), keys, review.UnsupportedForm, "Partition "+review.Code(ident(d.Name))+" can't take "+review.Code(wrong)+" in "+review.Code(r.written)+"'s partitioning.")
				return true
			}
		}
		return b.partitionDefsFail(p.kind, c.PartitionDefs, p.names, p.bounds, p.maxvalue, p.values, keys)
	}
	return false
}

// alterPartitionColumnsFail fails an ALTER TABLE of a synced partitioned
// table the change has not touched, and that the statement does not
// partition again, that drops a column its partitioning function reads
// (1054 or 1072), adds a unique key or primary key without every such
// column (1503), or adds a foreign key, which no partitioned table takes
// (1506) (each verified on mariadb:11.8.8).
func (b *builder) alterPartitionColumnsFail(v *ast.AlterTableStmt, r tableRef) bool {
	st := b.syncedOf(r)
	if st == nil || st.part == nil || alterHasPartitionBy(v.Commands) {
		return false
	}
	for _, c := range v.Commands {
		if c != nil && c.Type == ast.ATRemovePartitioning {
			return false
		}
	}
	rel := relKey(r.db, r.name)
	keys := []walk.Key{rel}
	partCols := make(map[string]bool, len(st.part.cols))
	for _, k := range st.part.cols {
		partCols[k] = true
	}
	// covers reports whether a key's columns hold every partition
	// column, false for a key whose parts the walk does not read.
	covers := func(cols []*ast.IndexColumn) (bool, bool) {
		have := make(map[string]bool, len(cols))
		for _, ic := range cols {
			k, ok := b.plainKeyPart(ic)
			if !ok {
				return false, false
			}
			have[k] = true
		}
		for k := range partCols {
			if !have[k] {
				return false, true
			}
		}
		return true, true
	}
	uniqueFails := func(rng review.Range) bool {
		b.fail(rng, keys, review.UnsupportedForm, "Every unique key of "+review.Code(r.written)+" must hold every column of its partitioning function.")
		return true
	}
	for _, c := range v.Commands {
		if c == nil {
			continue
		}
		switch c.Type {
		case ast.ATDropColumn:
			k, ok := b.key(spaceColumn, c.Name)
			if !ok || !partCols[k] || !b.namedOnce(v, k) {
				continue
			}
			// IF EXISTS drops a column the table has.
			if f, _ := b.l.Lookup(colKey(r.db, r.name, k)); c.IfExists && f != walk.Present {
				continue
			}
			b.fail(rangeOf(c.Loc), []walk.Key{rel, colKey(r.db, r.name, k)}, review.UnsupportedForm, "Column "+review.Code(ident(c.Name))+" is read by the partitioning function of "+review.Code(r.written)+".")
			return true
		case ast.ATAddConstraint, ast.ATAddIndex:
			con := c.Constraint
			if con == nil || con.IfNotExists && (con.Type != ast.ConstrForeignKey || b.fkIfNotExistsSkips(r, con)) {
				continue
			}
			switch con.Type {
			case ast.ConstrForeignKey:
				b.fail(rangeOf(c.Loc), keys, review.UnsupportedForm, "Table "+review.Code(r.written)+" is partitioned, and a partitioned table takes no foreign key.")
				return true
			case ast.ConstrPrimaryKey, ast.ConstrUnique:
				if len(partCols) == 0 {
					continue
				}
				if ok, read := covers(con.IndexColumns); read && !ok {
					return uniqueFails(rangeOf(c.Loc))
				}
			}
		case ast.ATAddColumn, ast.ATModifyColumn:
			if len(partCols) == 0 || c.IfNotExists || c.IfExists {
				continue
			}
			for _, col := range columnsOf(c) {
				if col == nil || !colDeclares(col, ast.ColConstrUnique) && !colDeclares(col, ast.ColConstrPrimaryKey) {
					continue
				}
				k, ok := b.key(spaceColumn, col.Name)
				if !ok {
					continue
				}
				if len(partCols) > 1 || !partCols[k] {
					return uniqueFails(rangeOf(col.Loc))
				}
			}
		}
	}
	return false
}

// createUniqueFails is alterPartitionColumnsFail's 1503 for CREATE UNIQUE
// INDEX.
func (b *builder) createUniqueFails(v *ast.CreateIndexStmt, r tableRef) bool {
	if !v.Unique || v.IfNotExists || v.OrReplace {
		return false
	}
	st := b.syncedOf(r)
	if st == nil || st.part == nil || len(st.part.cols) == 0 {
		return false
	}
	have := make(map[string]bool, len(v.Columns))
	for _, ic := range v.Columns {
		k, ok := b.plainKeyPart(ic)
		if !ok {
			return false
		}
		have[k] = true
	}
	for _, k := range st.part.cols {
		if !have[k] {
			b.fail(b.anchor(), []walk.Key{relKey(r.db, r.name)}, review.UnsupportedForm, "Every unique key of "+review.Code(r.written)+" must hold every column of its partitioning function.")
			return true
		}
	}
	return false
}

// nonIntegerTypes are the base types no partitioning function of HASH,
// RANGE, or LIST reads (1659), and blobTypes those KEY partitioning
// can't read either (1502), each verified on mariadb:11.8.8.
var (
	nonIntegerTypes = map[string]bool{
		"char": true, "varchar": true, "binary": true, "varbinary": true, "decimal": true, "float": true, "double": true,
		"date": true, "time": true, "datetime": true, "timestamp": true,
		"tinytext": true, "text": true, "mediumtext": true, "longtext": true,
		"tinyblob": true, "blob": true, "mediumblob": true, "longblob": true,
	}
	blobTypes = map[string]bool{
		"tinytext": true, "text": true, "mediumtext": true, "longtext": true,
		"tinyblob": true, "blob": true, "mediumblob": true, "longblob": true,
	}
)

// partitionTypeName returns the metadata type of a partition clause's
// method.
func partitionTypeName(pc *ast.PartitionClause) metadata.TablePartitionMetadata_Type {
	switch pc.Type {
	case ast.PartitionRange:
		if len(pc.Columns) > 0 {
			return metadata.TablePartitionMetadata_RANGE_COLUMNS
		}
		return metadata.TablePartitionMetadata_RANGE
	case ast.PartitionList:
		if len(pc.Columns) > 0 {
			return metadata.TablePartitionMetadata_LIST_COLUMNS
		}
		return metadata.TablePartitionMetadata_LIST
	case ast.PartitionHash:
		return metadata.TablePartitionMetadata_HASH
	case ast.PartitionKey:
		return metadata.TablePartitionMetadata_KEY
	}
	return metadata.TablePartitionMetadata_TYPE_UNSPECIFIED
}

// createPartitionFails fails a CREATE TABLE whose partition clause
// MariaDB refuses against its own column list (each verified on
// mariadb:11.8.8): partitionDefsFail's faults among its partitions; a
// HASH, RANGE, or LIST function that is a column of a type other than an
// integer (1659), and KEY partitioning on a BLOB or TEXT column (1502); a
// foreign key (1506).
func (b *builder) createPartitionFails(v *ast.CreateTableStmt, r tableRef) bool {
	pc := v.Partitions
	if pc == nil || v.Like != nil || v.Select != nil {
		return false
	}
	kind := partitionTypeName(pc)
	if kind == metadata.TablePartitionMetadata_TYPE_UNSPECIFIED {
		return false
	}
	types := make(map[string]string, len(v.Columns))
	for _, c := range v.Columns {
		if k, ok := b.key(spaceColumn, c.Name); ok {
			types[k] = astBaseType(c.TypeName)
		}
	}
	for _, c := range v.Constraints {
		if c != nil && c.Type == ast.ConstrForeignKey {
			b.fail(rangeOf(c.Loc), nil, review.UnsupportedForm, "Table "+review.Code(r.written)+" is partitioned, and a partitioned table takes no foreign key.")
			return true
		}
	}
	for _, c := range v.Columns {
		if c != nil && colDeclares(c, ast.ColConstrReferences) {
			b.fail(rangeOf(c.Loc), nil, review.UnsupportedForm, "Table "+review.Code(r.written)+" is partitioned, and a partitioned table takes no foreign key.")
			return true
		}
	}
	colType := func(name string) (string, string) {
		k, ok := b.key(spaceColumn, name)
		if !ok {
			return "", ""
		}
		return name, types[k]
	}
	switch kind {
	case metadata.TablePartitionMetadata_HASH, metadata.TablePartitionMetadata_RANGE, metadata.TablePartitionMetadata_LIST:
		if ref, ok := pc.Expr.(*ast.ColumnRef); ok && !ref.Star && ref.Table == "" {
			if name, t := colType(ref.Column); nonIntegerTypes[t] {
				b.fail(rangeOf(pc.Loc), nil, review.UnsupportedForm, "Partition column "+review.Code(ident(name))+" is "+review.Code(t)+", and the partitioning function must be an integer.")
				return true
			}
		}
	case metadata.TablePartitionMetadata_KEY:
		for _, col := range pc.Columns {
			if name, t := colType(col); blobTypes[t] {
				b.fail(rangeOf(pc.Loc), nil, review.UnsupportedForm, "Partition column "+review.Code(ident(name))+" is "+review.Code(t)+", which no partitioning can read.")
				return true
			}
		}
	}
	var bounds []int64
	var values map[int64]bool
	single := len(pc.Columns) <= 1
	if single && rangeLike(kind) {
		bounds = []int64{}
	}
	if single && listLike(kind) {
		values = make(map[int64]bool)
	}
	return b.partitionDefsFail(kind, pc.Partitions, nil, bounds, false, values, nil)
}

// plainKeyPart returns the key of the column a key part names whole,
// false for a functional part, a prefix, or a period.
func (b *builder) plainKeyPart(ic *ast.IndexColumn) (string, bool) {
	if ic == nil || ic.Functional || ic.HasPrefix || ic.WithoutOverlaps {
		return "", false
	}
	ref, ok := ic.Expr.(*ast.ColumnRef)
	if !ok || ref.Star || ref.Table != "" || ref.Schema != "" {
		return "", false
	}
	return b.key(spaceColumn, ref.Column)
}
