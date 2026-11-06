package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// Every unique key of a partitioned table, the primary key included, must include
// every column the partitioning function reads, and include it whole, unless it is a
// GLOBAL index, which this parser does not read. TiDB v8.5.5 checks a CREATE TABLE's
// keys in its own order (pkg/ddl, checkPartitioningKeysConstraints) and each key an
// ALTER TABLE or CREATE INDEX adds as that command runs, and names the key it
// refuses:
//
//   - a clustered primary key: 1503 "A CLUSTERED INDEX must include all columns in
//     the table's partitioning function";
//   - a NONCLUSTERED primary key or a unique key: 8264 "Global Index is needed for
//     index '<name>', since the unique index is not including all partitioning
//     columns, and GLOBAL is not given as IndexOption".
//
// TiDB ignores SUBPARTITION BY, so only the partitioning function counts.

// partitionKeyColumns returns the columns pc's partitioning function reads: the
// column list of a COLUMNS or KEY partitioning, the column references of a RANGE,
// LIST or HASH expression, and for KEY () the primary key. A KEY () on a table with no
// primary key contributes nothing, so the check never asks for a column the engine
// does not.
func partitionKeyColumns(tbl *Table, pc *nodes.PartitionClause) []string {
	var names []string
	switch {
	case len(pc.Columns) > 0:
		names = pc.Columns
	case pc.Type == nodes.PartitionKey:
		names = effectivePrimaryKeyColumns(tbl)
	default:
		names = exprColumnNames(pc.Expr)
	}
	var cols []string
	for _, name := range names {
		if !containsFold(cols, name) {
			cols = append(cols, name)
		}
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

// effectivePrimaryKeyColumns returns the columns KEY () partitions by: the primary
// key's, or else those of the first unique key whose parts are all whole NOT NULL
// columns.
func effectivePrimaryKeyColumns(tbl *Table) []string {
	var key *Index
	for _, idx := range tbl.Indexes {
		if idx.Primary {
			key = idx
			break
		}
	}
	if key == nil {
		for _, idx := range tbl.Indexes {
			if !idx.Unique || idx.Fulltext || idx.Spatial || len(idx.Columns) == 0 {
				continue
			}
			whole := true
			for _, ic := range idx.Columns {
				col := tbl.GetColumn(ic.Name)
				if ic.Name == "" || col == nil || col.Nullable || !keyPartIsWholeColumn(tbl, ic) {
					whole = false
					break
				}
			}
			if whole {
				key = idx
				break
			}
		}
	}
	if key == nil {
		return nil
	}
	cols := make([]string, 0, len(key.Columns))
	for _, ic := range key.Columns {
		cols = append(cols, ic.Name)
	}
	return cols
}

// keyPartIsWholeColumn reports whether a key part holds its whole column. A prefix
// does not, unless it is as long as the column: UNIQUE KEY (a(20)) on a VARCHAR(20)
// column partitioned by KEY(a) is accepted and (a(19)) is an 8264. A prefix on a
// column whose length the catalog cannot read counts as whole, so no error is raised
// that the engine might not raise.
func keyPartIsWholeColumn(tbl *Table, ic *IndexColumn) bool {
	if ic.Length == 0 {
		return true
	}
	col := tbl.GetColumn(ic.Name)
	if col == nil {
		return true
	}
	switch strings.ToLower(col.DataType) {
	case "char", "varchar", "binary", "varbinary":
	default:
		return true
	}
	open := strings.IndexByte(col.ColumnType, '(')
	end := strings.IndexByte(col.ColumnType, ')')
	if open < 0 || end <= open+1 {
		return true
	}
	n := 0
	for _, ch := range col.ColumnType[open+1 : end] {
		if ch < '0' || ch > '9' {
			return true
		}
		n = n*10 + int(ch-'0')
	}
	return ic.Length >= n
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

// primaryKeyClustered reports whether TiDB clusters the table by its primary key:
// as CLUSTERED or NONCLUSTERED says, and otherwise by default in CREATE TABLE
// (tidb_enable_clustered_index is ON) but not in ALTER TABLE ... ADD PRIMARY KEY.
func primaryKeyClustered(tbl *Table, pk *Index) bool {
	for _, con := range tbl.Constraints {
		if con.Type == ConPrimaryKey && con.Clustered != nil {
			return *con.Clustered
		}
	}
	return !pk.addedByAlter
}

// checkUniqueKeysCoverPartitionColumns checks idxs, in the order TiDB does. A
// clustered primary key over one integer column is the row handle, which TiDB checks
// after the indexes.
func checkUniqueKeysCoverPartitionColumns(tbl *Table, idxs []*Index, partCols []string) error {
	if len(partCols) == 0 {
		return nil
	}
	handle := false
	for _, idx := range idxs {
		if !idx.Unique && !idx.Primary {
			continue
		}
		// A global index is exempt: covering the partitioning columns is what
		// GLOBAL buys an exemption from, and v8.5.5 runs
		// `ADD UNIQUE INDEX i (b) GLOBAL` on a table partitioned by a where
		// it answers 8264 without the option.
		if idx.Global {
			continue
		}
		// TiDB names an unnamed key over an expression expression_index, which the
		// catalog leaves unnamed, so it cannot name the key TiDB refuses.
		if idx.Name == "" || uniqueKeyCovers(tbl, idx, partCols) {
			continue
		}
		if idx.Primary && primaryKeyClustered(tbl, idx) {
			if isRowHandle(tbl, idx) {
				handle = true
				continue
			}
			return errClusteredIndexNeedAllFieldsInPF()
		}
		return errGlobalIndexNeeded(idx.Name)
	}
	if handle {
		return errClusteredIndexNeedAllFieldsInPF()
	}
	return nil
}

func isRowHandle(tbl *Table, pk *Index) bool {
	if len(pk.Columns) != 1 || pk.Columns[0].Length != 0 {
		return false
	}
	col := tbl.GetColumn(pk.Columns[0].Name)
	return col != nil && isIntegerType(strings.ToLower(col.DataType))
}

func errClusteredIndexNeedAllFieldsInPF() error {
	return &Error{
		Code:     1503,
		SQLState: "HY000",
		Message:  "A CLUSTERED INDEX must include all columns in the table's partitioning function",
	}
}

func errGlobalIndexNeeded(name string) error {
	return &Error{
		Code:     8264,
		SQLState: "HY000",
		Message:  fmt.Sprintf("Global Index is needed for index '%s', since the unique index is not including all partitioning columns, and GLOBAL is not given as IndexOption", name),
	}
}

// validateCreateTableUniqueKeys checks the keys of a CREATE TABLE in TiDB's order:
// the table's constraints as written, then the columns' (TiDB appends a column's
// PRIMARY KEY or UNIQUE after the table constraints), which is the order the
// catalog keeps them in.
func validateCreateTableUniqueKeys(tbl *Table, pc *nodes.PartitionClause) error {
	return checkUniqueKeysCoverPartitionColumns(tbl, tbl.Indexes, partitionKeyColumns(tbl, pc))
}

// validateNewUniqueKey checks a key an ALTER TABLE or CREATE INDEX is about to add.
func validateNewUniqueKey(tbl *Table, idx *Index) error {
	if tbl.Partitioning == nil {
		return nil
	}
	return checkUniqueKeysCoverPartitionColumns(tbl, []*Index{idx}, tbl.Partitioning.keyColumns)
}

// validateRepartition checks the unique keys of a table ALTER TABLE ... PARTITION BY
// partitions anew. The catalog does not apply a repartition, so the partitioning it
// keeps no longer says which columns a unique key needs.
func validateRepartition(tbl *Table, pc *nodes.PartitionClause) error {
	if pc == nil {
		return nil
	}
	// UPDATE INDEXES says what each named index becomes once the table is
	// repartitioned, so it has to apply before the coverage check rather
	// than after: v8.5.5 runs
	// `ALTER TABLE t PARTITION BY HASH(a) PARTITIONS 2 UPDATE INDEXES
	// (i GLOBAL)` on a table whose unique key covers none of the new
	// partitioning columns, and answers 8264 for the same statement without
	// the clause (read on 2026-11-04). Checking the index's old setting left
	// the clause parsed and unread, so the statement it exists to allow
	// failed.
	undo, err := applyUpdateIndexes(tbl, pc.UpdateIndexes)
	if err != nil {
		return err
	}
	if err := checkUniqueKeysCoverPartitionColumns(tbl, tbl.Indexes, partitionKeyColumns(tbl, pc)); err != nil {
		// The statement is refused whole, so the settings go back.
		undo()
		return err
	}
	if tbl.Partitioning != nil {
		tbl.Partitioning.keyColumns = nil
	}
	return nil
}

// applyUpdateIndexes sets Global on each index an UPDATE INDEXES clause names
// and returns the function that puts the old settings back, for the caller
// that refuses the statement afterwards.
//
// Two names are errors, both 1280 ("Incorrect index name '<name>'") on
// v8.5.5: one the table does not have, and one the clause names twice --
// `UPDATE INDEXES (i GLOBAL, i LOCAL)` and even `(i GLOBAL, i GLOBAL)` are
// refused, so the list names each index at most once (read on 2026-11-04).
// Each index is therefore saved once, which also means undo does not depend
// on the order it runs in.
func applyUpdateIndexes(tbl *Table, specs []*nodes.UpdateIndexSpec) (func(), error) {
	type saved struct {
		idx    *Index
		global bool
	}
	var was []saved
	undo := func() {
		for _, s := range was {
			s.idx.Global = s.global
		}
	}
	seen := make(map[*Index]bool, len(specs))
	for _, spec := range specs {
		if spec == nil {
			continue
		}
		idx := findIndexByName(tbl, spec.Name)
		if idx == nil || seen[idx] {
			undo()
			return nil, errWrongNameForIndex(spec.Name)
		}
		seen[idx] = true
		was = append(was, saved{idx: idx, global: idx.Global})
		idx.Global = spec.Global
	}
	return undo, nil
}

// findIndexByName returns the index of the name, matched the way the server
// matches an index name: case-insensitively.
func findIndexByName(tbl *Table, name string) *Index {
	for _, idx := range tbl.Indexes {
		if toLower(idx.Name) == toLower(name) {
			return idx
		}
	}
	return nil
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
