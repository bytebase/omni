package catalog

import (
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// Every unique key of a partitioned table, the primary key included, must include
// every column its partitioning and subpartitioning functions read, unless it is a
// GLOBAL unique key: 5.7.25-OceanBase_CE-v4.3.5.6 answers 1503 for CREATE TABLE, for
// each ALTER TABLE form that adds a unique key, for CREATE UNIQUE INDEX, and for
// ALTER TABLE ... PARTITION BY over a table whose unique keys leave a column out,
// whatever the partitioning type. Unlike MySQL it counts no prefix as the column, not
// even one as long as the column, and promotes no NOT NULL unique key to primary.

// partitionKeyColumns returns the columns pc's partitioning and subpartitioning
// functions read: the column list of a COLUMNS or KEY partitioning, the column
// references of a RANGE, LIST or HASH expression, and for KEY () the primary key.
func partitionKeyColumns(tbl *Table, pc *nodes.PartitionClause) []string {
	var cols []string
	add := func(names []string) {
		for _, name := range names {
			if !containsFold(cols, name) {
				cols = append(cols, name)
			}
		}
	}
	switch {
	case len(pc.Columns) > 0:
		add(pc.Columns)
	case pc.Type == nodes.PartitionKey:
		add(primaryKeyColumns(tbl))
	default:
		add(exprColumnNames(pc.Expr))
	}
	switch {
	case len(pc.SubPartColumns) > 0:
		add(pc.SubPartColumns)
	case pc.SubPartType == nodes.PartitionKey:
		add(primaryKeyColumns(tbl))
	case pc.SubPartExpr != nil:
		add(exprColumnNames(pc.SubPartExpr))
	}
	return cols
}

func exprColumnNames(expr nodes.ExprNode) []string {
	if expr == nil {
		return nil
	}
	var names []string
	nodes.Inspect(expr, func(n nodes.Node) bool {
		if ref, ok := n.(*nodes.ColumnRef); ok && ref.Column != "" {
			names = append(names, ref.Column)
		}
		return true
	})
	return names
}

func containsFold(names []string, name string) bool {
	for _, n := range names {
		if toLower(n) == toLower(name) {
			return true
		}
	}
	return false
}

// uniqueKeyCovers reports whether idx holds every one of partCols, none of them by a
// prefix: 4.3.5.6 refuses UNIQUE KEY (a(20)) on a VARCHAR(20) column partitioned by
// KEY(a), where MySQL accepts it.
func uniqueKeyCovers(idx *Index, partCols []string) bool {
	for _, pc := range partCols {
		covered := false
		for _, ic := range idx.Columns {
			if ic.Name != "" && ic.Length == 0 && toLower(ic.Name) == toLower(pc) {
				covered = true
				break
			}
		}
		if !covered {
			return false
		}
	}
	return true
}

// validateUniqueKeysCoverPartitionColumns checks every unique key of tbl against
// partCols.
func validateUniqueKeysCoverPartitionColumns(tbl *Table, partCols []string) error {
	return checkUniqueKeysCoverPartitionColumns(tbl, partCols, false, func(*Index) bool { return true })
}

// checkUniqueKeysCoverPartitionColumns checks the unique keys of tbl that check
// selects, the primary key first, as OceanBase does. A GLOBAL unique index is not
// partitioned with the table, so it enforces uniqueness across every partition and
// needs none of the partitioning columns. 4.3.5 accepts exactly that and rejects the
// LOCAL and unstated forms with this same 1503, which is why only GLOBAL is exempt.
// The primary key never is: CREATE TABLE refuses a scope on it (1064), and ALTER
// TABLE ... ADD PRIMARY KEY (a) GLOBAL on a table partitioned by another column is
// still a 1503 on 4.3.5.6.
//
// rebuilt says the statement rebuilds the table (it adds a column or repartitions),
// where 4.3.5.6 words a unique key's refusal "A UNIQUE INDEX/PRIMARY KEY must ...".
func checkUniqueKeysCoverPartitionColumns(tbl *Table, partCols []string, rebuilt bool, check func(*Index) bool) error {
	if len(partCols) == 0 {
		return nil
	}
	for _, idx := range tbl.Indexes {
		if idx.Primary && check(idx) && !uniqueKeyCovers(idx, partCols) {
			return errUniqueKeyNeedAllFieldsInPF("PRIMARY KEY")
		}
	}
	for _, idx := range tbl.Indexes {
		if !idx.Unique || idx.Primary || strings.EqualFold(idx.Scope, "GLOBAL") {
			continue
		}
		if check(idx) && !uniqueKeyCovers(idx, partCols) {
			if rebuilt {
				return errUniqueKeyNeedAllFieldsInPF("UNIQUE INDEX/PRIMARY KEY")
			}
			return errUniqueKeyNeedAllFieldsInPF("UNIQUE INDEX")
		}
	}
	return nil
}

func errUniqueKeyNeedAllFieldsInPF(key string) error {
	return &Error{
		Code:     1503,
		SQLState: "HY000",
		Message:  "A " + key + " must include all columns in the table's partitioning function",
	}
}

// alterRepartition returns the partitioning an ALTER TABLE ... PARTITION BY gives the
// table, or nil.
func alterRepartition(cmds []*nodes.AlterTableCmd) *nodes.PartitionClause {
	var pc *nodes.PartitionClause
	for _, cmd := range cmds {
		if cmd.Type == nodes.ATPartitionBy && cmd.PartitionBy != nil {
			pc = cmd.PartitionBy
		}
	}
	return pc
}

// refusedBeforePartitioning reports the ALTER TABLE shapes 4.3.5.6 refuses with 1235
// whatever the partitioning, which the catalog does not model: a column added with
// PRIMARY KEY, and a primary key added beside another index. The check stays out of
// their way rather than answer a 1503 the engine never gives.
func refusedBeforePartitioning(cmds []*nodes.AlterTableCmd) bool {
	addsPrimaryKey, addsIndex := false, false
	for _, cmd := range cmds {
		switch cmd.Type {
		case nodes.ATAddColumn:
			defs := cmd.Columns
			if cmd.Column != nil {
				defs = append(defs, cmd.Column)
			}
			for _, def := range defs {
				for _, cc := range def.Constraints {
					if cc.Type == nodes.ColConstrPrimaryKey {
						return true
					}
				}
			}
		case nodes.ATAddIndex, nodes.ATAddConstraint:
			if cmd.Constraint == nil {
				continue
			}
			switch cmd.Constraint.Type {
			case nodes.ConstrPrimaryKey:
				if addsPrimaryKey || addsIndex {
					return true
				}
				addsPrimaryKey = true
			case nodes.ConstrUnique, nodes.ConstrIndex, nodes.ConstrFulltextIndex, nodes.ConstrSpatialIndex:
				if addsPrimaryKey {
					return true
				}
				addsIndex = true
			}
		}
	}
	return false
}

// validateAlteredUniqueKeys checks the table an ALTER TABLE leaves. Against a new
// partitioning every unique key is checked; against the one the table keeps, only the
// unique keys the statement added (before holds the indexes it started with), so an
// earlier statement the catalog accepted and the engine would not cannot fail this one.
func validateAlteredUniqueKeys(tbl *Table, before map[*Index]bool, cmds []*nodes.AlterTableCmd, repartition *nodes.PartitionClause) error {
	if refusedBeforePartitioning(cmds) {
		return nil
	}
	if repartition != nil {
		return checkUniqueKeysCoverPartitionColumns(tbl, partitionKeyColumns(tbl, repartition), true, func(*Index) bool { return true })
	}
	if tbl.Partitioning == nil {
		return nil
	}
	rebuilt := false
	for _, cmd := range cmds {
		if cmd.Type == nodes.ATAddColumn {
			rebuilt = true
		}
	}
	return checkUniqueKeysCoverPartitionColumns(tbl, tbl.Partitioning.keyColumns, rebuilt, func(idx *Index) bool { return !before[idx] })
}

// validateNewIndexCoversPartitioning checks the index CREATE INDEX is about to add.
func validateNewIndexCoversPartitioning(tbl *Table, idx *Index) error {
	if tbl.Partitioning == nil || !idx.Unique {
		return nil
	}
	n := len(tbl.Indexes)
	tbl.Indexes = append(tbl.Indexes, idx)
	err := checkUniqueKeysCoverPartitionColumns(tbl, tbl.Partitioning.keyColumns, false, func(i *Index) bool { return i == idx })
	tbl.Indexes = tbl.Indexes[:n]
	return err
}

// renamePartitionKeyColumn follows a column rename in the partitioning's key columns,
// as the catalog follows it in the indexes. 4.3.5.6 takes ALTER TABLE ... CHANGE of a
// partitioning column to a new name.
func renamePartitionKeyColumn(tbl *Table, oldName, newName string) {
	if tbl.Partitioning == nil || len(tbl.Partitioning.keyColumns) == 0 {
		return
	}
	cols := make([]string, len(tbl.Partitioning.keyColumns))
	for i, name := range tbl.Partitioning.keyColumns {
		if toLower(name) == toLower(oldName) {
			name = newName
		}
		cols[i] = name
	}
	tbl.Partitioning.keyColumns = cols
}
