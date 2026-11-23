package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// The checks MariaDB makes on a partitioned table's definitions after it
// parsed them, where a definition meets the table's columns and the partitions
// it already has. Each code and the order between them was read on
// mariadb:11.8.8.

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

func errCoalesceNone() error {
	return &Error{Code: 1515, SQLState: "HY000", Message: "At least one partition must be coalesced"}
}

func errReorganizeChangesCount() error {
	return &Error{Code: 1510, SQLState: "HY000", Message: "REORGANIZE PARTITION can only be used to reorganize partitions not to change their numbers"}
}

func errReorganizeNotConsecutive() error {
	return &Error{Code: 1519, SQLState: "HY000", Message: "When reorganizing a set of partitions they must be in consecutive order"}
}

// checkPartitionFieldTypes refuses a partitioning function on a column of a
// type it cannot partition by: a BLOB or TEXT column under any method (1502),
// and under a method other than KEY a column outside the integers (an
// expression's bare column) or outside the integers, strings and dates (a
// COLUMNS list) with 1659. MariaDB answers both before the unique keys the
// function must be covered by (1503): `PARTITION BY HASH (amount)` over a
// DECIMAL column outside the primary key is 1659.
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

// checkPartitionClauseFields holds a PARTITION BY clause's function and its
// subpartitioning's to the types they can read.
func checkPartitionClauseFields(tbl *Table, pc *nodes.PartitionClause) error {
	if err := checkPartitionFieldTypes(tbl, pc.Type, pc.Columns, pc.Expr); err != nil {
		return err
	}
	if pc.SubPartType != 0 {
		return checkPartitionFieldTypes(tbl, pc.SubPartType, pc.SubPartColumns, pc.SubPartExpr)
	}
	return nil
}

// checkPartitionNamesUnique refuses partitions that share a name (1517): the
// names a statement writes, among themselves and beside the ones the table
// keeps. MariaDB answers it before it compares bounds.
func checkPartitionNamesUnique(keep []*PartitionDefInfo, defs []*nodes.PartitionDef) error {
	seen := map[string]bool{}
	for _, pd := range keep {
		seen[toLower(pd.Name)] = true
	}
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

// checkPartitionSequence holds a RANGE or LIST table's partitions, in their
// order, to bounds that strictly increase (1493, 1481) and values no two
// share (1495), as checkAddedRangeBounds holds added ones.
func checkPartitionSequence(tbl *Table, pi *PartitionInfo, parts []*PartitionDefInfo) error {
	empty := *pi
	empty.Partitions = nil
	defs := make([]*nodes.PartitionDef, 0, len(parts))
	for _, pd := range parts {
		if pd.bound == nil {
			continue
		}
		defs = append(defs, &nodes.PartitionDef{Name: pd.Name, Values: pd.bound, ValuesIn: strings.HasPrefix(pi.Type, "LIST")})
	}
	return checkRangeListBounds(tbl, &empty, defs)
}

// firstUnknownPartition returns the first name a list holds that the table
// has no partition or subpartition of, or "": TRUNCATE and EXCHANGE answer it
// with 1735. ALL alone names every partition.
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
			for _, sp := range pd.SubPartitions {
				if sp != nil && strings.EqualFold(sp.Name, name) {
					found = true
				}
			}
		}
		if !found {
			return name
		}
	}
	return ""
}
