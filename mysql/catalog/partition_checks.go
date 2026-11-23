package catalog

import (
	"fmt"
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// The checks MySQL makes on a partitioned table's definitions after it parsed
// them, where a definition meets the table's method and the partitions it
// already has. Each code and the order between them was read on mysql:8.0.46.

// errPartitionWrongValues is 1480: a definition's VALUES clause belongs to
// another method, VALUES IN outside LIST and VALUES LESS THAN outside RANGE.
func errPartitionWrongValues(in bool) error {
	method, values := "RANGE", "LESS THAN"
	if in {
		method, values = "LIST", "IN"
	}
	return &Error{Code: 1480, SQLState: "HY000",
		Message: fmt.Sprintf("Only %s PARTITIONING can use VALUES %s in partition definition", method, values)}
}

func errRangeNotIncreasing() error {
	return &Error{Code: 1493, SQLState: "HY000", Message: "VALUES LESS THAN value must be strictly increasing for each partition"}
}

func errListConstantTwice() error {
	return &Error{Code: 1495, SQLState: "HY000", Message: "Multiple definition of same constant in list partitioning"}
}

func errPartitionList(action string) error {
	return &Error{Code: 1507, SQLState: "HY000", Message: "Error in list of partitions to " + action}
}

func errDropOnlyRangeList() error {
	return &Error{Code: 1512, SQLState: "HY000", Message: "DROP PARTITION can only be used on RANGE/LIST partitions"}
}

func errCoalesceNone() error {
	return &Error{Code: 1515, SQLState: "HY000", Message: "At least one partition must be coalesced"}
}

func errReorganizeChangesCount() error {
	return &Error{Code: 1510, SQLState: "HY000", Message: "REORGANIZE PARTITION can only be used to reorganize partitions not to change their numbers"}
}

func errReorganizeNotConsecutive() error {
	return &Error{Code: 1519, SQLState: "HY000", Message: "When reorganizing a set of partitions they must be in consecutive order"}
}

func errReorganizeRange() error {
	return &Error{Code: 1520, SQLState: "HY000", Message: "Reorganize of range partitions cannot change total ranges except for last partition where it can extend the range"}
}

func errUnknownPartition(name, table string) error {
	return &Error{Code: 1735, SQLState: "HY000", Message: fmt.Sprintf("Unknown partition '%s' in table '%s'", name, table)}
}

func errExchangePartitioned(table string) error {
	return &Error{Code: 1732, SQLState: "HY000", Message: fmt.Sprintf("Table to exchange with partition is partitioned: '%s'", table)}
}

func errPartitionFieldType(col string) error {
	return &Error{Code: 1659, SQLState: "HY000", Message: fmt.Sprintf("Field '%s' is of a not allowed type for this type of partitioning", col)}
}

func errBlobInPartitionFunc() error {
	return &Error{Code: 1502, SQLState: "HY000", Message: "A BLOB field is not allowed in partition function"}
}

// checkPartitionValuesKind refuses a definition whose VALUES clause is not its
// method's (1480). The parser keeps VALUES IN as a list and VALUES LESS THAN
// as its one bound, or as a list of one bound per column for RANGE COLUMNS, so
// a list is read as VALUES IN where the method takes no such list: under RANGE
// on one column, and a list of rows under RANGE COLUMNS.
func checkPartitionValuesKind(ptype nodes.PartitionType, columns int, v nodes.Node) error {
	if v == nil {
		return nil
	}
	list, isList := v.(*nodes.List)
	in := isList
	if isList && ptype == nodes.PartitionRange && columns > 1 {
		in = len(list.Items) > 0
		for _, item := range list.Items {
			if _, row := item.(*nodes.List); !row {
				in = false
			}
		}
	}
	switch ptype {
	case nodes.PartitionRange:
		if in {
			return errPartitionWrongValues(true)
		}
	case nodes.PartitionList:
		if !in {
			return errPartitionWrongValues(false)
		}
	default:
		return errPartitionWrongValues(in)
	}
	return nil
}

// checkPartitionFieldTypes refuses a partitioning function on a column of a
// type it cannot partition by: a BLOB, TEXT or JSON column under any method
// (1502), and under a method other than KEY a column outside the integers (an
// expression's bare column) or outside the integers, strings and dates (a
// COLUMNS list) with 1659. MySQL answers both before the unique keys the
// function must be covered by (1503).
func checkPartitionFieldTypes(tbl *Table, typ nodes.PartitionType, columns []string, expr nodes.ExprNode) error {
	var names []string
	bare := false
	if len(columns) > 0 {
		names = columns
	} else if ref, ok := expr.(*nodes.ColumnRef); ok && ref.Table == "" {
		names, bare = []string{ref.Column}, true
	}
	for _, name := range names {
		col := tbl.GetColumn(name)
		if col == nil {
			continue
		}
		dt := strings.ToLower(col.DataType)
		switch dt {
		case "tinyblob", "blob", "mediumblob", "longblob", "tinytext", "text", "mediumtext", "longtext", "json":
			return errBlobInPartitionFunc()
		}
		if typ == nodes.PartitionKey {
			continue
		}
		integer := false
		switch dt {
		case "tinyint", "smallint", "mediumint", "int", "integer", "bigint", "bit", "year":
			integer = true
		}
		switch {
		case bare && !integer:
			return errPartitionFieldType(col.Name)
		case !bare:
			switch dt {
			case "tinyint", "smallint", "mediumint", "int", "integer", "bigint",
				"char", "varchar", "binary", "varbinary", "date", "datetime", "time":
			default:
				return errPartitionFieldType(col.Name)
			}
		}
	}
	return nil
}

// checkPartitionNamesUnique refuses a list that names one partition twice
// (1517), which MySQL answers before it compares the bounds.
func checkPartitionNamesUnique(defs []*nodes.PartitionDef) error {
	seen := map[string]bool{}
	for _, pd := range defs {
		if pd == nil || pd.Name == "" {
			continue
		}
		if seen[toLower(pd.Name)] {
			return errSameNamePartition(pd.Name)
		}
		seen[toLower(pd.Name)] = true
	}
	return nil
}

// checkPartitionBounds holds a RANGE table's partitions to bounds that
// strictly increase (1493), and a LIST table's to values no two partitions
// share (1495). Two bounds are compared only where their order is certain,
// as rangeColumnsNotIncreasing compares them.
func checkPartitionBounds(tbl *Table, pi *PartitionInfo) error {
	if pi == nil {
		return nil
	}
	switch pi.Type {
	case "RANGE", "RANGE COLUMNS":
		// RANGE on an expression takes MAXVALUE in its last partition alone
		// (1481), which ADD PARTITION after one answers; RANGE COLUMNS takes it
		// anywhere and answers such a partition with 1493.
		if pi.Type == "RANGE" {
			for i, pd := range pi.Partitions {
				if strings.EqualFold(pd.ValueExpr, "MAXVALUE") && i != len(pi.Partitions)-1 {
					return &Error{Code: 1481, SQLState: "HY000", Message: "MAXVALUE can only be used in last partition definition"}
				}
			}
		}
		kinds := []rangeBoundKind{rangeBoundInteger}
		if pi.Type == "RANGE COLUMNS" {
			kinds = rangeColumnsBoundKinds(tbl, pi.Columns)
		}
		for i := 1; i < len(pi.Partitions); i++ {
			prev, cur := storedBound(pi.Partitions[i-1].ValueExpr), storedBound(pi.Partitions[i].ValueExpr)
			if prev != nil && cur != nil && rangeColumnsNotIncreasing(prev, cur, kinds) {
				return errRangeNotIncreasing()
			}
		}
	case "LIST", "LIST COLUMNS":
		// A value of a column whose collation ignores case is compared
		// without it, each column of a tuple by its own collation.
		var fold []bool
		if pi.Type == "LIST COLUMNS" {
			for _, name := range pi.Columns {
				col := tbl.GetColumn(name)
				fold = append(fold, col != nil && strings.HasSuffix(columnCollation(tbl, col), "_ci"))
			}
		}
		key := func(v string) string {
			parts := []string{v}
			if len(fold) > 1 && strings.HasPrefix(v, "(") && strings.HasSuffix(v, ")") {
				parts = splitPartitionValues(v[1 : len(v)-1])
			}
			for i := range parts {
				if i < len(fold) && fold[i] {
					parts[i] = strings.ToLower(parts[i])
				}
			}
			return strings.Join(parts, ",")
		}
		seen := map[string]bool{}
		for _, pd := range pi.Partitions {
			for _, v := range splitPartitionValues(pd.ValueExpr) {
				v = key(v)
				if seen[v] {
					return errListConstantTwice()
				}
				seen[v] = true
			}
		}
	}
	return nil
}

// rangeBoundOrder orders two stored RANGE bounds of one table, a against b,
// when the order is certain.
func rangeBoundOrder(tbl *Table, pi *PartitionInfo, a, b string) (int, bool) {
	x, y := storedBound(a), storedBound(b)
	if x == nil || y == nil {
		return 0, false
	}
	kinds := []rangeBoundKind{rangeBoundInteger}
	if pi.Type == "RANGE COLUMNS" {
		kinds = rangeColumnsBoundKinds(tbl, pi.Columns)
	}
	xs, ys := partitionBoundItems(x), partitionBoundItems(y)
	if len(xs) != len(ys) || len(xs) != len(kinds) {
		return 0, false
	}
	for i := range xs {
		xm, ym := isPartitionMaxvalue(xs[i]), isPartitionMaxvalue(ys[i])
		switch {
		case xm && ym:
			continue
		case xm:
			return 1, true
		case ym:
			return -1, true
		}
		order, known := comparePartitionBounds(xs[i], ys[i], kinds[i])
		if !known {
			return 0, false
		}
		if order != 0 {
			return order, true
		}
	}
	return 0, true
}

// storedBound reads back a RANGE bound the catalog stored as text
// (partitionValueToString): MAXVALUE, an integer or a quoted string, or a
// list of them for RANGE COLUMNS. It returns nil for anything else.
func storedBound(text string) nodes.Node {
	items := splitPartitionValues(text)
	if len(items) == 0 {
		return nil
	}
	nodesOf := make([]nodes.Node, 0, len(items))
	for _, item := range items {
		switch {
		case strings.EqualFold(item, "MAXVALUE"):
			nodesOf = append(nodesOf, &nodes.String{Str: "MAXVALUE"})
		case len(item) >= 2 && item[0] == '\'' && item[len(item)-1] == '\'':
			body := item[1 : len(item)-1]
			body = strings.ReplaceAll(strings.ReplaceAll(body, `\'`, `'`), `''`, `'`)
			nodesOf = append(nodesOf, &nodes.StringLit{Value: body})
		default:
			v, err := strconv.ParseInt(item, 10, 64)
			if err != nil {
				return nil
			}
			nodesOf = append(nodesOf, &nodes.IntLit{Value: v, Text: item})
		}
	}
	if len(nodesOf) == 1 {
		return nodesOf[0]
	}
	return &nodes.List{Items: nodesOf}
}

// splitPartitionValues splits a stored value list at its top-level commas,
// outside quotes and parentheses.
func splitPartitionValues(text string) []string {
	var out []string
	depth, start := 0, 0
	var quote byte
	for i := 0; i < len(text); i++ {
		ch := text[i]
		switch {
		case quote != 0:
			if ch == '\\' {
				i++
			} else if ch == quote {
				quote = 0
			}
		case ch == '\'' || ch == '"':
			quote = ch
		case ch == '(':
			depth++
		case ch == ')':
			depth--
		case ch == ',' && depth == 0:
			out = append(out, strings.TrimSpace(text[start:i]))
			start = i + 1
		}
	}
	if rest := strings.TrimSpace(text[start:]); rest != "" {
		out = append(out, rest)
	}
	return out
}

// checkPartitionNamesKnown refuses a list of a table's partitions that names
// one it does not have, or one twice (1507, for DROP, REORGANIZE and REBUILD).
func checkPartitionNamesKnown(pi *PartitionInfo, names []string, action string) error {
	have := map[string]bool{}
	for _, pd := range pi.Partitions {
		have[toLower(pd.Name)] = true
	}
	seen := map[string]bool{}
	for _, name := range names {
		key := toLower(name)
		if !have[key] || seen[key] {
			return errPartitionList(action)
		}
		seen[key] = true
	}
	return nil
}

// firstUnknownPartition returns the first name a list holds that the table
// has no partition of, or "": TRUNCATE and EXCHANGE answer it with 1735.
func firstUnknownPartition(pi *PartitionInfo, names []string) string {
	for _, name := range names {
		if strings.EqualFold(name, "ALL") && len(names) == 1 {
			return ""
		}
		found := false
		for _, pd := range pi.Partitions {
			if strings.EqualFold(pd.Name, name) {
				found = true
				break
			}
		}
		if !found {
			return name
		}
	}
	return ""
}
