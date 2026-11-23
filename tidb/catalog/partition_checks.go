package catalog

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// The checks TiDB makes on a partitioned table's partitions after it parsed
// them, where a definition meets the table's method, its columns and the
// partitions it already has. Each code, message and the order between them
// was read on v8.5.5 (the partitioned@ fixtures of simulate.tsv).

func errRangeNotIncreasing() error {
	return &Error{Code: 1493, SQLState: "HY000", Message: "VALUES LESS THAN value must be strictly increasing for each partition"}
}

func errListConstantTwice() error {
	return &Error{Code: 1495, SQLState: "HY000", Message: "Multiple definition of same constant in list partitioning"}
}

func errDropOnlyRangeList() error {
	return &Error{Code: 1512, SQLState: "HY000", Message: "DROP PARTITION can only be used on RANGE/LIST partitions"}
}

func errDropPartitionList() error {
	return &Error{Code: 1507, SQLState: "HY000", Message: "Error in list of partitions to DROP"}
}

func errDropLastPartition() error {
	return &Error{Code: 1508, SQLState: "HY000", Message: "Cannot remove all partitions, use DROP TABLE instead"}
}

func errCoalesceNone() error {
	return &Error{Code: 1515, SQLState: "HY000", Message: "At least one partition must be coalesced"}
}

func errUnknownPartition(name, table string) error {
	return &Error{Code: 1735, SQLState: "HY000", Message: fmt.Sprintf("Unknown partition '%s' in table '%s'", name, table)}
}

func errUnsupportedReorganize() error {
	return &Error{Code: ErrUnsupportedDDL, SQLState: sqlState(ErrUnsupportedDDL), Message: "Unsupported reorganize partition"}
}

func errPartitionColumnDependency(col string) error {
	return &Error{Code: 3855, SQLState: "HY000", Message: fmt.Sprintf("Column '%s' has a partitioning function dependency and cannot be dropped or renamed", col)}
}

func errForeignKeyOnPartitioned() error {
	return &Error{Code: 1506, SQLState: "HY000", Message: "Foreign key clause is not yet supported in conjunction with partitioning"}
}

// boundKind is how a RANGE bound of one column compares where the order is
// certain: integers by value, and dates or datetimes written in one fixed,
// zero-padded shape by their text.
type boundKind int

const (
	boundOpen boundKind = iota
	boundInteger
	boundTemporal
)

var temporalBound = regexp.MustCompile(`^'\d{4}-\d{2}-\d{2}( \d{2}:\d{2}:\d{2})?'$`)

// rangeBoundKinds reads the kind of each RANGE bound column: the integer the
// expression form yields, or each COLUMNS column's.
func rangeBoundKinds(tbl *Table, pi *PartitionInfo) []boundKind {
	if pi.Type != "RANGE COLUMNS" {
		return []boundKind{boundInteger}
	}
	kinds := make([]boundKind, len(pi.Columns))
	for i, name := range pi.Columns {
		col := tbl.GetColumn(name)
		if col == nil {
			continue
		}
		switch strings.ToLower(col.DataType) {
		case "tinyint", "smallint", "mediumint", "int", "integer", "bigint":
			kinds[i] = boundInteger
		case "date", "datetime":
			kinds[i] = boundTemporal
		}
	}
	return kinds
}

// compareBounds orders two stored RANGE bounds of one table, a against b,
// where the order is certain.
func compareBounds(a, b string, kinds []boundKind) (int, bool) {
	xs, ys := splitPartitionValues(a), splitPartitionValues(b)
	if len(xs) != len(ys) || len(xs) != len(kinds) {
		return 0, false
	}
	for i := range xs {
		xm, ym := strings.EqualFold(xs[i], "MAXVALUE"), strings.EqualFold(ys[i], "MAXVALUE")
		switch {
		case xm && ym:
			continue
		case xm:
			return 1, true
		case ym:
			return -1, true
		}
		order, ok := 0, false
		switch kinds[i] {
		case boundInteger:
			x, errX := strconv.ParseInt(xs[i], 10, 64)
			y, errY := strconv.ParseInt(ys[i], 10, 64)
			if errX == nil && errY == nil {
				order, ok = cmpInt(x, y), true
			}
		case boundTemporal:
			if temporalBound.MatchString(xs[i]) && temporalBound.MatchString(ys[i]) && len(xs[i]) == len(ys[i]) {
				order, ok = strings.Compare(xs[i], ys[i]), true
			}
		}
		if !ok {
			return 0, false
		}
		if order != 0 {
			return order, true
		}
	}
	return 0, true
}

func cmpInt(a, b int64) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

// checkPartitionSequence holds a RANGE table's partitions, in their order, to
// bounds that strictly increase (1493), and a LIST table's to values no two
// partitions share (1495). Two bounds are compared only where their order is
// certain, and a LIST value the catalog cannot read as a constant is passed
// over.
func checkPartitionSequence(tbl *Table, pi *PartitionInfo, parts []*PartitionDefInfo) error {
	if pi == nil {
		return nil
	}
	switch pi.Type {
	case "RANGE", "RANGE COLUMNS":
		kinds := rangeBoundKinds(tbl, pi)
		for i := 1; i < len(parts); i++ {
			if order, ok := compareBounds(parts[i-1].ValueExpr, parts[i].ValueExpr, kinds); ok && order >= 0 {
				return errRangeNotIncreasing()
			}
		}
	case "LIST", "LIST COLUMNS":
		// A value of a column whose collation ignores case is compared
		// without it.
		fold := false
		if pi.Type == "LIST COLUMNS" && len(pi.Columns) == 1 {
			if col := tbl.GetColumn(pi.Columns[0]); col != nil {
				collation := col.Collation
				if collation == "" {
					collation = tbl.Collation
				}
				fold = strings.HasSuffix(strings.ToLower(collation), "_ci")
			}
		}
		seen := map[string]bool{}
		for _, pd := range parts {
			for _, v := range splitPartitionValues(pd.ValueExpr) {
				if fold {
					v = strings.ToLower(v)
				}
				if seen[v] {
					return errListConstantTwice()
				}
				seen[v] = true
			}
		}
	}
	return nil
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

// checkPartitionNamesUnique refuses a CREATE TABLE's partitions that share a
// name (1517).
func checkPartitionNamesUnique(defs []*nodes.PartitionDef) error {
	seen := map[string]bool{}
	for _, pd := range defs {
		if pd == nil || pd.Name == "" {
			continue
		}
		if seen[toLower(pd.Name)] {
			return errDuplicatePartitionName(pd.Name)
		}
		seen[toLower(pd.Name)] = true
	}
	return nil
}

// firstUnknownPartition returns the first name a list holds that the table
// has no partition of, or "": TRUNCATE and EXCHANGE answer it with 1735. ALL
// alone names every partition.
func firstUnknownPartition(pi *PartitionInfo, names []string) string {
	for _, name := range names {
		if strings.EqualFold(name, "ALL") && len(names) == 1 {
			return ""
		}
		found := false
		for _, pd := range pi.Partitions {
			if strings.EqualFold(pd.Name, name) {
				found = true
			}
		}
		if !found {
			return name
		}
	}
	return ""
}

// partitionReadsColumn reports whether a table's partitioning reads a column:
// a COLUMNS or KEY list that names it, or an expression that does.
func partitionReadsColumn(tbl *Table, name string) bool {
	pi := tbl.Partitioning
	if pi == nil {
		return false
	}
	for _, c := range pi.Columns {
		if strings.EqualFold(c, name) {
			return true
		}
	}
	for _, c := range pi.keyColumns {
		if strings.EqualFold(c, name) {
			return true
		}
	}
	return false
}

// checkPartitionDefinitions holds a PARTITION BY clause's partitions to
// distinct names (1517), then to rising RANGE bounds and distinct LIST values
// (1493, 1495), as v8.5.5 does for CREATE TABLE and ALTER TABLE ... PARTITION
// BY.
func checkPartitionDefinitions(tbl *Table, pc *nodes.PartitionClause) error {
	if err := checkPartitionNamesUnique(pc.Partitions); err != nil {
		return err
	}
	pi := buildPartitionInfo(pc)
	return checkPartitionSequence(tbl, pi, pi.Partitions)
}

// isErrorCode reports whether err is a catalog error of the code.
func isErrorCode(err error, code int) bool {
	e, ok := err.(*Error)
	return ok && e.Code == code
}
