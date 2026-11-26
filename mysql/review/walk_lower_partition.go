package review

import (
	"strconv"
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/review/walk"
)

// partitionIntBound returns the integer a partition bound writes, false
// for anything else.
func partitionIntBound(n ast.Node) (int64, bool) {
	switch x := n.(type) {
	case *ast.IntLit:
		return x.Value, x.Text == "" || x.Text == strconv.FormatInt(x.Value, 10)
	case *ast.UnaryExpr:
		if x.Op == ast.UnaryMinus {
			if v, ok := partitionIntBound(x.Operand); ok && v > 0 {
				return -v, true
			}
		}
	case *ast.ParenExpr:
		return partitionIntBound(x.Expr)
	}
	return 0, false
}

func isMaxValue(n ast.Node) bool {
	s, ok := n.(*ast.String)
	return ok && s.Str == "MAXVALUE"
}

// listConstant returns a constant a LIST partition's VALUES IN holds as
// INFORMATION_SCHEMA.PARTITIONS writes it: an integer as digits and a
// string quoted. ok is false for anything else, or a string the walk
// doesn't compare.
func listConstant(n ast.Node) (string, bool) {
	if v, ok := partitionIntBound(n); ok {
		return strconv.FormatInt(v, 10), true
	}
	if s, ok := n.(*ast.StringLit); ok && s.Charset == "" && !strings.ContainsAny(s.Value, `'\`) {
		return "'" + s.Value + "'", true
	}
	return "", false
}

// splitPartitionValues splits a LIST partition's description into its
// constants, false when one holds a quote the walk doesn't read.
func splitPartitionValues(desc string) ([]string, bool) {
	var out []string
	for _, v := range strings.Split(desc, ",") {
		v = strings.TrimSpace(v)
		if v == "" || strings.Count(v, "'")%2 != 0 || strings.Contains(v, `\`) {
			return nil, false
		}
		out = append(out, v)
	}
	return out, true
}

// boundKindFault returns the fault of a partition definition whose bound
// is of another kind than the partitioning takes (1480): VALUES IN for
// RANGE, VALUES LESS THAN for LIST, any VALUES for HASH or KEY. A bound of
// one partitioning column is a single value for LESS THAN and a list for
// IN; with more columns both are lists, which the walk leaves alone.
func boundKindFault(typ ast.PartitionType, single bool, pd *ast.PartitionDef) string {
	if pd == nil || pd.Values == nil {
		return ""
	}
	_, list := pd.Values.(*ast.List)
	switch typ {
	case ast.PartitionHash, ast.PartitionKey:
		return "Partition " + review.Code(ident(pd.Name)) + " can't take a " + review.Code("VALUES") + " bound: only " + review.Code("RANGE") + " and " + review.Code("LIST") + " partitioning take one."
	case ast.PartitionRange:
		if single && list {
			return "Partition " + review.Code(ident(pd.Name)) + " can't take " + review.Code("VALUES IN") + ": only " + review.Code("LIST") + " partitioning takes it."
		}
	case ast.PartitionList:
		if single && !list {
			return "Partition " + review.Code(ident(pd.Name)) + " can't take " + review.Code("VALUES LESS THAN") + ": only " + review.Code("RANGE") + " partitioning takes it."
		}
	}
	return ""
}

// rangeBoundsFault returns the fault of RANGE bounds of one column, read in
// order after the table's last bound (last, "" for none): MAXVALUE anywhere
// but last among them (1481), or a bound no greater than the one before it,
// the table's MAXVALUE included (1493). Bounds the walk can't compare end
// the check.
func rangeBoundsFault(last string, defs []*ast.PartitionDef) string {
	for i, pd := range defs {
		if pd == nil {
			return ""
		}
		if last == "MAXVALUE" {
			if i > 0 {
				return review.Code("MAXVALUE") + " can only bound the last partition."
			}
			return "Partition " + review.Code(ident(pd.Name)) + "'s bound must be greater than the partition's before it."
		}
		if isMaxValue(pd.Values) {
			last = "MAXVALUE"
			continue
		}
		v, ok := partitionIntBound(pd.Values)
		if !ok {
			return ""
		}
		if last != "" {
			prev, err := strconv.ParseInt(last, 10, 64)
			if err != nil {
				return ""
			}
			if v <= prev {
				return "Partition " + review.Code(ident(pd.Name)) + "'s bound must be greater than the partition's before it."
			}
		}
		last = strconv.FormatInt(v, 10)
	}
	return ""
}

// listValuesFault returns the fault of LIST values of one column that
// repeat a constant, of the definitions or of the values the table holds
// (seen), 1495.
func listValuesFault(seen map[string]bool, defs []*ast.PartitionDef) string {
	for _, pd := range defs {
		list, ok := pd.Values.(*ast.List)
		if !ok {
			return ""
		}
		for _, item := range list.Items {
			c, ok := listConstant(item)
			if !ok {
				return ""
			}
			if seen[c] {
				return "The partitions' " + review.Code("VALUES IN") + " name " + review.Code(c) + " twice."
			}
			seen[c] = true
		}
	}
	return ""
}

// notIntegerPartitionType reports a column type that RANGE, LIST, or HASH
// partitioning can't read by itself (1659): a string or a binary string.
func notIntegerPartitionType(t string) bool {
	switch t {
	case "char", "varchar", "binary", "varbinary", "tinytext", "text", "mediumtext", "longtext",
		"tinyblob", "blob", "mediumblob", "longblob":
		return true
	}
	return false
}

// isBlobType reports a TEXT or BLOB type, which no partitioning reads
// (1502).
func isBlobType(t string) bool {
	switch t {
	case "tinytext", "text", "mediumtext", "longtext", "tinyblob", "blob", "mediumblob", "longblob":
		return true
	}
	return false
}

// partitionClauseTextFails fails a PARTITION BY clause, of CREATE TABLE or
// ALTER TABLE, that MySQL refuses for its own text and the types of the
// columns it reads (colType, "" when not known), each verified on
// mysql:8.0.46 in mysql/catalog/testdata/simulate.tsv:
//
//   - a bound of another kind than the partitioning takes (1480);
//   - two partitions of one name (1517);
//   - a column of a type the partitioning can't read: a string for RANGE,
//     LIST, or HASH of a bare column (1659), TEXT or BLOB for KEY (1502);
//   - RANGE bounds of one column out of order (1481, 1493), or LIST values
//     of one column named twice (1495).
func (b *builder) partitionClauseTextFails(pc *ast.PartitionClause, colType func(string) string, rng review.Range) bool {
	if pc == nil {
		return false
	}
	fail := func(msg string) bool {
		b.fail(rng, nil, review.UnsupportedForm, msg)
		return true
	}
	single := len(pc.Columns) <= 1
	for _, pd := range pc.Partitions {
		if msg := boundKindFault(pc.Type, single, pd); msg != "" {
			return fail(msg)
		}
	}
	names := make(map[string]bool, len(pc.Partitions))
	for _, pd := range pc.Partitions {
		if pd == nil || pd.Name == "" {
			continue
		}
		k := strings.ToLower(pd.Name)
		if names[k] {
			b.fail(rng, nil, review.DuplicateName, statesTwice("names", "partition", ident(pd.Name)))
			return true
		}
		names[k] = true
	}
	switch pc.Type {
	case ast.PartitionKey:
		for _, c := range pc.Columns {
			if t := colType(c); isBlobType(t) {
				return fail("Partition column " + review.Code(ident(c)) + " is " + review.Code(t) + ", which no partitioning can read.")
			}
		}
	case ast.PartitionRange, ast.PartitionList, ast.PartitionHash:
		if len(pc.Columns) == 0 {
			if c, ok := unparen(pc.Expr).(*ast.ColumnRef); ok && c.Table == "" && !c.Star {
				if t := colType(c.Column); notIntegerPartitionType(t) {
					return fail("Partition column " + review.Code(ident(c.Column)) + " is " + review.Code(t) + ", and the partitioning function must be an integer.")
				}
			}
		}
	}
	if !single {
		return false
	}
	switch pc.Type {
	case ast.PartitionRange:
		if msg := rangeBoundsFault("", pc.Partitions); msg != "" {
			return fail(msg)
		}
	case ast.PartitionList:
		if msg := listValuesFault(make(map[string]bool), pc.Partitions); msg != "" {
			return fail(msg)
		}
	}
	return false
}

// partitionOpFails fails an ALTER TABLE partition operation on a table the
// walk knows the partitioning of (pi), each verified on mysql:8.0.46:
//
//   - a partition definition of a bound of another kind (1480), or one out
//     of order after the table's own RANGE bounds (1481, 1493), or naming a
//     LIST value the table holds (1495);
//   - DROP PARTITION of HASH or KEY partitioning (1512);
//   - a partition name the table doesn't have: DROP (1507), and TRUNCATE,
//     EXCHANGE and the other operations on named partitions (1735);
//   - COALESCE PARTITION 0 (1515), or of every partition (1508);
//   - REORGANIZE PARTITION of HASH or KEY partitioning into another number
//     of partitions (1510);
//   - EXCHANGE PARTITION with a table that is partitioned itself (1732).
func (b *builder) partitionOpFails(r tableRef, pi *tablePartitionInfo, c *ast.AlterTableCmd, rng review.Range) bool {
	if pi == nil || !pi.partitioned || !pi.hasType {
		return false
	}
	fail := func(msg string) bool {
		b.fail(rng, nil, review.UnsupportedForm, msg)
		return true
	}
	hashOrKey := pi.partType == ast.PartitionHash || pi.partType == ast.PartitionKey
	unknown := func(names []string) string {
		if pi.names == nil {
			return ""
		}
		for _, n := range names {
			if !pi.names[strings.ToLower(n)] {
				return n
			}
		}
		return ""
	}
	switch c.Type {
	case ast.ATAddPartition:
		for _, pd := range c.PartitionDefs {
			if msg := boundKindFault(pi.partType, pi.single, pd); msg != "" {
				return fail(msg)
			}
		}
		if !pi.single || pi.bounds == nil {
			return false
		}
		switch pi.partType {
		case ast.PartitionRange:
			last := ""
			if n := len(pi.bounds); n > 0 {
				last = pi.bounds[n-1].value
				if _, err := strconv.ParseInt(last, 10, 64); err != nil && last != "MAXVALUE" {
					return false
				}
			}
			if msg := rangeBoundsFault(last, c.PartitionDefs); msg != "" {
				return fail(msg)
			}
		case ast.PartitionList:
			seen := make(map[string]bool)
			for _, bd := range pi.bounds {
				vals, ok := splitPartitionValues(bd.value)
				if !ok {
					return false
				}
				for _, v := range vals {
					seen[v] = true
				}
			}
			if msg := listValuesFault(seen, c.PartitionDefs); msg != "" {
				return fail(msg)
			}
		}
	case ast.ATDropPartition:
		if hashOrKey {
			return fail(review.Code("DROP PARTITION") + " works only on " + review.Code("RANGE") + " or " + review.Code("LIST") + " partitioning.")
		}
		if n := unknown(c.PartitionNames); n != "" {
			return fail(notFound("partition", ident(n), r.written))
		}
	case ast.ATTruncatePartition, ast.ATExchangePartition, ast.ATRebuildPartition,
		ast.ATDiscardPartitionTablespace, ast.ATImportPartitionTablespace:
		names := c.PartitionNames
		if c.Type == ast.ATExchangePartition {
			names = []string{c.Name}
		}
		// A subpartition's name, which pi.names doesn't hold, works here
		// too.
		if !c.AllPartitions && !pi.subparts {
			if n := unknown(names); n != "" {
				return fail(notFound("partition", ident(n), r.written))
			}
		}
		if c.Type == ast.ATExchangePartition && c.ExchangeTable != nil {
			if er, st := b.ref(c.ExchangeTable); st == refOK {
				if f, temp, kind := b.resolve(er); !temp && f == walk.Present && kind == relTable {
					if epi := b.w.partitionState[relKey(er.db, er.name)]; epi != nil && epi.partitioned {
						return fail("Table " + review.Code(er.written) + " is partitioned, and a partition can be exchanged only with a table that isn't.")
					}
				}
			}
		}
	case ast.ATCoalescePartition:
		if !hashOrKey {
			return false
		}
		if c.Number == 0 {
			return fail(review.Code("COALESCE PARTITION") + " must remove at least one partition.")
		}
		if pi.count > 0 && c.Number >= pi.count {
			return fail("Cannot remove all partitions from " + review.Code(r.written) + ".")
		}
	case ast.ATReorganizePartition:
		if hashOrKey && len(c.PartitionNames) > 0 && len(c.PartitionDefs) > 0 && len(c.PartitionNames) != len(c.PartitionDefs) && unknown(c.PartitionNames) == "" && pi.names != nil {
			return fail(review.Code("REORGANIZE PARTITION") + " of " + review.Code("HASH") + " or " + review.Code("KEY") + " partitioning must keep the number of partitions.")
		}
		if len(c.PartitionNames) == 0 || pi.names == nil || unknown(c.PartitionNames) != "" {
			return false
		}
		if msg := reorganizeFault(r, pi, c); msg != "" {
			return fail(msg)
		}
	}
	return false
}

// reorganizeFault returns the fault of REORGANIZE PARTITION of partitions
// the table has, each verified on mysql:8.0.46: a new partition named after
// one the statement keeps (1517); for RANGE of one value, reorganized
// partitions that aren't consecutive (1519), or new bounds that don't rise
// from the partition before them to the one after them (1493).
func reorganizeFault(r tableRef, pi *tablePartitionInfo, c *ast.AlterTableCmd) string {
	reorged := make(map[string]bool, len(c.PartitionNames))
	for _, n := range c.PartitionNames {
		reorged[strings.ToLower(n)] = true
	}
	for _, pd := range c.PartitionDefs {
		if pd == nil {
			return ""
		}
		if k := strings.ToLower(pd.Name); pi.names[k] && !reorged[k] {
			return alreadyExists("partition", ident(pd.Name), r.written)
		}
	}
	if pi.partType != ast.PartitionRange || !pi.single || len(pi.bounds) != len(pi.names) {
		return ""
	}
	first, last := -1, -1
	for i, bd := range pi.bounds {
		if reorged[bd.name] {
			if first < 0 {
				first = i
			}
			last = i
		}
	}
	if first < 0 || last-first+1 != len(reorged) {
		return review.Code("REORGANIZE PARTITION") + " must name consecutive partitions."
	}
	prev := ""
	if first > 0 {
		prev = pi.bounds[first-1].value
	}
	if msg := rangeBoundsFault(prev, c.PartitionDefs); msg != "" || last+1 == len(pi.bounds) || len(c.PartitionDefs) == 0 {
		return msg
	}
	end, ok := partitionIntBound(c.PartitionDefs[len(c.PartitionDefs)-1].Values)
	next := pi.bounds[last+1]
	nv, err := strconv.ParseInt(next.value, 10, 64)
	if ok && err == nil && nv <= end {
		return "Partition " + review.Code(ident(next.name)) + "'s bound must be greater than the partition's before it."
	}
	return ""
}

// uniqueKeyPartitionFails fails an ALTER TABLE that adds a primary or
// unique key, of whole columns, leaving out a column the table's
// partitioning function reads (1503, verified on mysql:8.0.46). A statement
// that also repartitions the table, or drops or renames a column, is left
// alone.
func (b *builder) uniqueKeyPartitionFails(r tableRef, v *ast.AlterTableStmt, d *alterDelta) bool {
	pi := b.w.partitionState[relKey(r.db, r.name)]
	if pi == nil || !pi.partitioned || len(pi.cols) == 0 || d.repartition != nil || d.removesPartitioning ||
		len(d.droppedCols) > 0 || len(d.renamedFrom) > 0 {
		return false
	}
	for _, c := range v.Commands {
		if c == nil || (c.Type != ast.ATAddConstraint && c.Type != ast.ATAddIndex) || c.Constraint == nil {
			continue
		}
		if c.Constraint.Type != ast.ConstrUnique && c.Constraint.Type != ast.ConstrPrimaryKey {
			continue
		}
		covered := make(map[string]bool)
		whole := true
		for _, p := range constraintKeyParts(c.Constraint) {
			if p.functional || p.length > 0 {
				whole = false
				break
			}
			if f, _, _ := b.alterFinalCol(r, d, p.name); f != walk.Present {
				whole = false
				break
			}
			covered[strings.ToLower(p.name)] = true
		}
		if !whole {
			continue
		}
		for col := range pi.cols {
			if !covered[col] {
				b.fail(rangeOf(c.Loc), nil, review.UnsupportedForm, "Every unique key of "+review.Code(r.written)+" must include all columns in its partitioning function.")
				return true
			}
		}
	}
	return false
}
