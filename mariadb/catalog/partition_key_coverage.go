package catalog

import nodes "github.com/bytebase/omni/mariadb/ast"

// Every unique key of a partitioned table, the primary key included, must include
// every column its partitioning and subpartitioning functions read, and include it
// whole: MariaDB 11.8 answers 1503 for CREATE TABLE, for each ALTER TABLE form that
// adds a unique key, for CREATE UNIQUE INDEX, and for ALTER TABLE ... PARTITION BY
// over a table whose unique keys leave a column out, whatever the partitioning type
// (sql_partition.cc, check_primary_key and check_unique_keys). SYSTEM_TIME
// partitioning reads no column a unique key could leave out.

// partitionKeyColumns returns the columns pc's partitioning and subpartitioning
// functions read: the column list of a COLUMNS or KEY partitioning, the column
// references of a RANGE, LIST or HASH expression, and for KEY () the primary key,
// which MariaDB partitions by. A KEY () on a table with no primary key contributes
// nothing, so the check never asks for a column the engine does not.
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
		add(effectivePrimaryKeyColumns(tbl))
	default:
		add(exprColumnNames(pc.Expr))
	}
	switch {
	case len(pc.SubPartColumns) > 0:
		add(pc.SubPartColumns)
	case pc.SubPartType == nodes.PartitionKey:
		add(effectivePrimaryKeyColumns(tbl))
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

// effectivePrimaryKey returns the index MariaDB treats as the primary key: the
// PRIMARY KEY, or else the first unique key whose parts are all whole NOT NULL
// columns, which MariaDB promotes in its place (TABLE_SHARE::primary_key). So a SERIAL
// column on a table partitioned by another column is refused as "A PRIMARY KEY must
// include ...".
func effectivePrimaryKey(tbl *Table) *Index {
	for _, idx := range tbl.Indexes {
		if idx.Primary {
			return idx
		}
	}
	for _, idx := range tbl.Indexes {
		if !idx.Unique || idx.Fulltext || idx.Spatial || len(idx.Columns) == 0 {
			continue
		}
		promoted := true
		for _, ic := range idx.Columns {
			col := tbl.GetColumn(ic.Name)
			if ic.Name == "" || col == nil || col.Nullable || !keyPartIsWholeColumn(tbl, ic) {
				promoted = false
				break
			}
		}
		if promoted {
			return idx
		}
	}
	return nil
}

func effectivePrimaryKeyColumns(tbl *Table) []string {
	pk := effectivePrimaryKey(tbl)
	if pk == nil {
		return nil
	}
	cols := make([]string, 0, len(pk.Columns))
	for _, ic := range pk.Columns {
		cols = append(cols, ic.Name)
	}
	return cols
}

// keyPartIsWholeColumn reports whether a key part holds its whole column. A prefix
// does not, unless it is as long as the column: UNIQUE KEY (a(20)) on a VARCHAR(20)
// column partitioned by KEY(a) is accepted and (a(19)) is a 1503. A prefix on a column whose length the
// catalog cannot read counts as whole, so no error is raised that the engine might not
// raise.
func keyPartIsWholeColumn(tbl *Table, ic *IndexColumn) bool {
	if ic.Length == 0 {
		return true
	}
	col := tbl.GetColumn(ic.Name)
	if col == nil {
		return true
	}
	n, ok := fixedLengthStringColumnLimit(col)
	return !ok || ic.Length >= n
}

// uniqueKeyCovers reports whether idx holds every one of partCols whole.
func uniqueKeyCovers(tbl *Table, idx *Index, partCols []string) bool {
	for _, pc := range partCols {
		covered := false
		for _, ic := range idx.Columns {
			if ic.Name != "" && toLower(ic.Name) == toLower(pc) && keyPartIsWholeColumn(tbl, ic) {
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
	return checkUniqueKeysCoverPartitionColumns(tbl, partCols, func(*Index) bool { return true })
}

// checkUniqueKeysCoverPartitionColumns checks the unique keys of tbl that check
// selects. MariaDB tests the primary key first, so it names the primary key when both
// it and another unique key leave a column out.
func checkUniqueKeysCoverPartitionColumns(tbl *Table, partCols []string, check func(*Index) bool) error {
	if len(partCols) == 0 {
		return nil
	}
	if pk := effectivePrimaryKey(tbl); pk != nil && check(pk) && !uniqueKeyCovers(tbl, pk, partCols) {
		return errUniqueKeyNeedAllFieldsInPF("PRIMARY KEY")
	}
	for _, idx := range tbl.Indexes {
		if (idx.Unique || idx.Primary) && check(idx) && !uniqueKeyCovers(tbl, idx, partCols) {
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

// validateAlteredUniqueKeys checks the table an ALTER TABLE leaves. Against a new
// partitioning every unique key is checked; against the one the table keeps, only the
// unique keys the statement added (before holds the indexes it started with), so an
// earlier statement the catalog accepted and the engine would not cannot fail this one.
func validateAlteredUniqueKeys(tbl *Table, before map[*Index]bool, repartition *nodes.PartitionClause) error {
	if repartition != nil {
		return validateUniqueKeysCoverPartitionColumns(tbl, partitionKeyColumns(tbl, repartition))
	}
	if tbl.Partitioning == nil {
		return nil
	}
	return checkUniqueKeysCoverPartitionColumns(tbl, tbl.Partitioning.keyColumns, func(idx *Index) bool { return !before[idx] })
}

// validateNewIndexCoversPartitioning checks the index CREATE INDEX is about to add,
// with it in place so that it can be the key MariaDB promotes to primary.
func validateNewIndexCoversPartitioning(tbl *Table, idx *Index) error {
	if tbl.Partitioning == nil || !idx.Unique {
		return nil
	}
	n := len(tbl.Indexes)
	tbl.Indexes = append(tbl.Indexes, idx)
	err := checkUniqueKeysCoverPartitionColumns(tbl, tbl.Partitioning.keyColumns, func(i *Index) bool { return i == idx })
	tbl.Indexes = tbl.Indexes[:n]
	return err
}

// renamePartitionKeyColumn follows a column rename in the partitioning's key columns,
// as the catalog follows it in the indexes.
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
