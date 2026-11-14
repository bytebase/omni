package analysis

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/parser"
)

const maxCommentLength = 1000

// RestoreOptions configures RestoreSQL.
type RestoreOptions struct {
	// Statement is the SQL that changed the table, quoted in a comment that
	// opens the restore SQL.
	Statement string
	// Database and Table are the changed table.
	Database string
	Table    string
	// BackupDatabase and BackupTable hold the backed-up rows.
	BackupDatabase string
	BackupTable    string
	// Columns are the table's columns other than generated ones, which the
	// restore copies back.
	Columns []string
	// Catalog is the instance's metadata. An UPDATE's restore reads the
	// table's unique keys from it.
	Catalog Catalog
}

// NormalColumnSet returns the lower-cased names of database.table's columns
// other than generated ones.
func NormalColumnSet(catalog Catalog, database, table string) (map[string]bool, error) {
	tableMetadata, err := getTable(catalog, database, table)
	if err != nil {
		return nil, err
	}
	result := make(map[string]bool)
	for _, col := range tableMetadata.Columns {
		if col == nil || col.Generation != nil {
			continue
		}
		result[strings.ToLower(col.Name)] = true
	}
	return result, nil
}

// getTable returns database.table from the catalog.
func getTable(catalog Catalog, database, table string) (*metadata.TableMetadata, error) {
	if catalog.GetDatabase == nil {
		return nil, errors.New("the catalog has no GetDatabase")
	}
	db, err := catalog.GetDatabase(database)
	if err != nil {
		return nil, fmt.Errorf("failed to get database metadata for %s: %w", database, err)
	}
	if db == nil {
		return nil, fmt.Errorf("database metadata is nil for %s", database)
	}
	if !db.HasSchema() {
		return nil, fmt.Errorf("schema is nil for %s", database)
	}
	tableMetadata := db.GetTable(table)
	if tableMetadata == nil {
		return nil, fmt.Errorf("table metadata is nil for %s.%s", database, table)
	}
	return tableMetadata, nil
}

// MutatesTable reports whether an UPDATE or DELETE of statement changes
// database.table. columns is NormalColumnSet's answer for the table, which
// decides an unqualified SET column. It returns the parser's error as is.
func MutatesTable(statement, database, table string, columns map[string]bool) (bool, error) {
	stmtList, err := parser.Parse(statement)
	if err != nil {
		return false, err
	}
	for _, node := range stmtList.Items {
		if containsTable(node, database, table, columns) {
			return true, nil
		}
	}
	return false, nil
}

// MatchingDMLs returns every UPDATE and DELETE of statement that changes
// database.table, as MutatesTable decides. It returns the parser's error as
// is.
func MatchingDMLs(statement, database, table string, columns map[string]bool) ([]ast.Node, error) {
	stmtList, err := parser.Parse(statement)
	if err != nil {
		return nil, err
	}
	var result []ast.Node
	for _, item := range stmtList.Items {
		switch item.(type) {
		case *ast.UpdateStmt, *ast.DeleteStmt:
			if containsTable(item, database, table, columns) {
				result = append(result, item)
			}
		default:
		}
	}
	return result, nil
}

// RestoreSQL returns the SQL that copies the backed-up rows of nodes, the
// UPDATE and DELETE statements one backup covers, back into the table. When
// any is an UPDATE, the rows overwrite every column the UPDATEs change, which
// needs a unique key they leave alone; otherwise they are inserted again.
func RestoreSQL(nodes []ast.Node, opts RestoreOptions) (string, error) {
	// Partition nodes by type. If any UPDATE exists in the bundle, ODKU
	// is required — pure INSERT SELECT would FAIL on the surviving UPDATE-
	// modified rows due to duplicate key. ODKU restores the listed columns
	// for those rows AND inserts the deleted rows (no conflict). If the
	// bundle is pure-DELETE, simple INSERT SELECT suffices (no surviving
	// rows means no conflicts).
	var updateStmts []*ast.UpdateStmt
	for _, n := range nodes {
		switch v := n.(type) {
		case *ast.UpdateStmt:
			updateStmts = append(updateStmts, v)
		case *ast.DeleteStmt:
			// DELETE doesn't contribute columns; tracked only by partition.
		default:
			return "", fmt.Errorf("unexpected statement type: %T", n)
		}
	}

	var result string
	if len(updateStmts) > 0 {
		r, err := generateUpdateRestore(opts.Catalog, updateStmts, opts.Database, opts.Table, opts.BackupDatabase, opts.BackupTable, opts.Columns)
		if err != nil {
			return "", err
		}
		result = r
	} else {
		result = generateDeleteRestore(opts.Database, opts.Table, opts.BackupDatabase, opts.BackupTable, opts.Columns)
	}

	sqlForComment, truncated := truncateString(opts.Statement, maxCommentLength)
	if truncated {
		sqlForComment += "..."
	}
	// A "*/" in the statement would end the comment early.
	sqlForComment = strings.ReplaceAll(sqlForComment, "*/", "* /")
	return fmt.Sprintf("/*\nOriginal SQL:\n%s\n*/\n%s", sqlForComment, result), nil
}

func generateDeleteRestore(originalDatabase, originalTable, backupDatabase, backupTable string, normalColumns []string) string {
	quotedColumnList := quoteColumnList(normalColumns, "", ", ")
	return fmt.Sprintf("INSERT INTO %s.%s (%s) SELECT %s FROM %s.%s;", quoteIdentifier(originalDatabase), quoteIdentifier(originalTable), quotedColumnList, quotedColumnList, quoteIdentifier(backupDatabase), quoteIdentifier(backupTable))
}

func generateUpdateRestore(catalog Catalog, stmts []*ast.UpdateStmt, originalDatabase, originalTable, backupDatabase, backupTable string, normalColumns []string) (string, error) {
	// Union update columns across every statement: the single-table backup
	// path bundles more than maxMixedDMLCount same-table DMLs into one backup,
	// and the restore must roll back every column any of them touches.
	//
	// Iteration order is preserved (first-seen-first) for deterministic
	// ODKU clause output; case-insensitive dedup matches the
	// hasDisjointUniqueKey lower-case map convention.
	seen := make(map[string]bool)
	var updateColumns []string
	var assignments [][]*ast.Assignment
	for _, stmt := range stmts {
		singleTables := extractSingleTablesFromTableExprs(originalDatabase, stmt.Tables)
		assignments = append(assignments, targetAssignments(stmt.SetList, originalDatabase, originalTable, singleTables, normalColumns))
		for _, col := range extractUpdateColumns(stmt.SetList, originalDatabase, originalTable, singleTables, normalColumns) {
			key := strings.ToLower(col)
			if !seen[key] {
				seen[key] = true
				updateColumns = append(updateColumns, col)
			}
		}
	}

	tableProto, err := getTable(catalog, originalDatabase, originalTable)
	if err != nil {
		return "", err
	}
	// A column the server sets ON UPDATE changed too, and is restored with
	// the rest.
	normal := make(map[string]bool)
	for _, column := range normalColumns {
		normal[strings.ToLower(column)] = true
	}
	for _, column := range tableProto.Columns {
		key := strings.ToLower(column.GetName())
		if column.GetOnUpdate() != "" && column.GetGeneration() == nil && normal[key] && !seen[key] {
			seen[key] = true
			updateColumns = append(updateColumns, column.GetName())
		}
	}

	if !hasDisjointUniqueKey(tableProto, updateColumns) {
		return "", fmt.Errorf("no disjoint unique key found for %s.%s", originalDatabase, originalTable)
	}
	// The upsert matches a backed-up row on any unique key, so a key the
	// UPDATEs change can match another row.
	if key := changedUniqueKey(tableProto, updateColumns, assignments); key != "" {
		return "", fmt.Errorf("prior backup cannot restore %s.%s: the UPDATE changes unique key %s, so the restore could match another row", originalDatabase, originalTable, key)
	}

	var buf strings.Builder
	quotedColumnList := quoteColumnList(normalColumns, "", ", ")
	fmt.Fprintf(&buf, "INSERT INTO %s.%s (%s) SELECT %s FROM %s.%s ON DUPLICATE KEY UPDATE ", quoteIdentifier(originalDatabase), quoteIdentifier(originalTable), quotedColumnList, quotedColumnList, quoteIdentifier(backupDatabase), quoteIdentifier(backupTable))

	for i, field := range updateColumns {
		if i > 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(&buf, "%s = VALUES(%s)", quoteIdentifier(field), quoteIdentifier(field))
	}
	buf.WriteString(";")
	return buf.String(), nil
}

// extractSingleTablesFromTableExprs walks omni TableExpr nodes and returns
// a map of alias-or-name to TableReference for simple table references.
func extractSingleTablesFromTableExprs(databaseName string, exprs []ast.TableExpr) map[string]*tableReference {
	result := make(map[string]*tableReference)
	for _, expr := range exprs {
		collectTableRefs(databaseName, expr, result)
	}
	return result
}

func collectTableRefs(databaseName string, expr ast.TableExpr, result map[string]*tableReference) {
	switch n := expr.(type) {
	case *ast.TableRef:
		database := n.Schema
		if database == "" {
			database = databaseName
		}
		table := &tableReference{
			Database: database,
			Table:    n.Name,
			Alias:    n.Alias,
		}
		// Map keys are normalized to lowercase to match TiDB's
		// case-insensitive identifier semantics — a SET clause that
		// references the alias with different case (e.g., `T1` declared,
		// `t1` referenced in SET) must still resolve.
		var key string
		if n.Alias != "" {
			key = strings.ToLower(n.Alias)
		} else {
			key = strings.ToLower(n.Name)
		}
		result[key] = table
	case *ast.JoinClause:
		collectTableRefs(databaseName, n.Left, result)
		collectTableRefs(databaseName, n.Right, result)
	default:
	}
}

// extractUpdateColumns extracts column names from SET assignments that
// belong to the original table.
//
// For qualified SET columns (`SET t1.a = ...`), the qualifier is looked
// up in singleTables (the alias-or-name → TableReference map for this
// stmt) — if the qualifier resolves to a TableReference whose .Table
// equals originalTable, the column is included. This is deterministic
// even for self-join UPDATEs where the same physical table appears
// under multiple aliases (e.g., `UPDATE test t1 JOIN test t2 SET t1.a
// = ...`); each SET clause's qualifier independently resolves to the
// correct alias entry.
func extractUpdateColumns(setList []*ast.Assignment, database, originalTable string, singleTables map[string]*tableReference, normalColumns []string) []string {
	var result []string
	for _, assignment := range setList {
		col := assignment.Column
		// A generated column only takes DEFAULT, so the restore leaves it out.
		if col == nil || !slices.ContainsFunc(normalColumns, func(c string) bool { return strings.EqualFold(c, col.Column) }) {
			continue
		}

		if col.Schema != "" && !strings.EqualFold(col.Schema, database) {
			continue
		}

		if col.Table == "" {
			result = append(result, col.Column)
			continue
		}

		// Qualified column: resolve the qualifier through the alias map
		// (handles self-joins deterministically). Lookup key is
		// lowercased to match the case-insensitive insert convention
		// in collectTableRefs. BOTH Database AND Table must match —
		// without the Database check, cross-database joins with
		// homonymous tables (e.g. `UPDATE db1.test t1 JOIN db2.test t2
		// SET t2.a = ...`) would incorrectly include the joined-DB
		// SET column in the target-DB's rollback.
		if entry, ok := singleTables[strings.ToLower(col.Table)]; ok && strings.EqualFold(entry.Database, database) && strings.EqualFold(entry.Table, originalTable) {
			result = append(result, col.Column)
			continue
		}
		// Fallback: qualifier IS the bare table name (no alias used).
		// col.Schema check above already filtered explicit-schema
		// mismatches.
		if strings.EqualFold(col.Table, originalTable) {
			result = append(result, col.Column)
		}
	}
	return result
}

// targetAssignments returns the SET assignments extractUpdateColumns takes
// the column names of.
func targetAssignments(setList []*ast.Assignment, database, originalTable string, singleTables map[string]*tableReference, normalColumns []string) []*ast.Assignment {
	var result []*ast.Assignment
	for _, assignment := range setList {
		col := assignment.Column
		// A generated column only takes DEFAULT, so the restore leaves it out.
		if col == nil || !slices.ContainsFunc(normalColumns, func(c string) bool { return strings.EqualFold(c, col.Column) }) {
			continue
		}

		if col.Schema != "" && !strings.EqualFold(col.Schema, database) {
			continue
		}

		if col.Table == "" {
			result = append(result, assignment)
			continue
		}

		// Qualified column: resolve the qualifier through the alias map
		// (handles self-joins deterministically). Lookup key is
		// lowercased to match the case-insensitive insert convention
		// in collectTableRefs. BOTH Database AND Table must match —
		// without the Database check, cross-database joins with
		// homonymous tables (e.g. `UPDATE db1.test t1 JOIN db2.test t2
		// SET t2.a = ...`) would incorrectly include the joined-DB
		// SET column in the target-DB's rollback.
		if entry, ok := singleTables[strings.ToLower(col.Table)]; ok && strings.EqualFold(entry.Database, database) && strings.EqualFold(entry.Table, originalTable) {
			result = append(result, assignment)
			continue
		}
		// Fallback: qualifier IS the bare table name (no alias used).
		// col.Schema check above already filtered explicit-schema
		// mismatches.
		if strings.EqualFold(col.Table, originalTable) {
			result = append(result, assignment)
		}
	}
	return result
}

func hasDisjointUniqueKey(tableProto *metadata.TableMetadata, updateColumns []string) bool {
	columnMap := make(map[string]bool)
	for _, column := range updateColumns {
		columnMap[strings.ToLower(column)] = true
	}

	// Build a fast lookup of regular (non-generated) column names. Used
	// to filter out unique keys whose Expressions reference generated
	// columns or non-column expressions (functional indexes) — those are
	// unsafe for ODKU rollback because string-comparison disjoint can't
	// tell whether the SET clause indirectly affects them.
	regularColumns := make(map[string]bool)
	for _, col := range tableProto.Columns {
		// NULLs never collide, so a key on a nullable column cannot match
		// the backed-up row to the changed one either.
		if col == nil || col.Generation != nil || col.Nullable {
			continue
		}
		regularColumns[strings.ToLower(col.Name)] = true
	}

	for _, index := range tableProto.Indexes {
		if !index.Primary && !index.Unique {
			continue
		}
		// Skip UKs whose Expressions reference anything other than
		// regular (non-generated) columns. A UK on a generated column
		// (e.g., `c_generated = a + b`) appears disjoint from SET cols
		// {a, b} via string comparison — but updating a or b changes
		// c_generated's value, so the UK is NOT safe for ODKU matching.
		// Same problem for functional indexes (Expressions = `(LOWER(email))`).
		// Conservative: treat any UK with non-regular-column expressions
		// as overlapping (skip). Never returns false-positive disjoint.
		//
		// ALSO skip UKs with empty Expressions — disjoint([]) returns
		// vacuously true, which would falsely mark the UK as safe.
		// Bytebase's TiDB metadata sync leaves Expressions empty for
		// expression-based index parts.
		if len(index.Expressions) == 0 {
			continue
		}
		safe := true
		for _, expr := range index.Expressions {
			if !regularColumns[strings.ToLower(expr)] {
				safe = false
				break
			}
		}
		if !safe {
			continue
		}
		if disjoint(index.Expressions, columnMap) {
			return true
		}
	}

	return false
}

func disjoint(a []string, b map[string]bool) bool {
	for _, item := range a {
		if _, ok := b[strings.ToLower(item)]; ok {
			return false
		}
	}
	return true
}

// containsTable checks whether the given AST node MUTATES the specified
// table (not merely references it via JOIN/USING).
//
// For UPDATE: the table must appear as the qualifier of at least one SET
// assignment — being only joined for filtering is not enough.
func containsTable(node ast.Node, database, table string, targetCols map[string]bool) bool {
	switch n := node.(type) {
	case *ast.UpdateStmt:
		return updateMutatesTable(n, database, table, targetCols)
	case *ast.DeleteStmt:
		return deleteMutatesTable(n, database, table)
	default:
	}
	return false
}

// deleteMutatesTable reports whether the DELETE removes rows from the specified
// physical table.
//
// Single-table DELETE (no USING): n.Tables is the physical table.
//
// Multi-table DELETE (DELETE t1, t2 FROM ...): n.Tables holds the delete-target
// ALIASES (or bare names), which must be resolved through the FROM/USING ref
// map to their physical tables, since the backup keys each item by the
// physical name. n.Using is the JOIN-only filter set, used here only to
// resolve aliases; it is never matched as a mutation target, because doing
// so would re-insert rows that were never deleted.
func deleteMutatesTable(n *ast.DeleteStmt, database, table string) bool {
	if len(n.Using) == 0 {
		for _, expr := range n.Tables {
			if tableExprReferences(expr, database, table) {
				return true
			}
		}
		return false
	}
	refs := extractSingleTablesFromTableExprs(database, n.Using)
	for _, target := range n.Tables {
		ref, ok := target.(*ast.TableRef)
		if !ok {
			continue
		}
		resolved, ok := refs[strings.ToLower(ref.Name)]
		if !ok {
			continue
		}
		if strings.EqualFold(resolved.Database, database) && strings.EqualFold(resolved.Table, table) {
			return true
		}
	}
	return false
}

// updateMutatesTable reports whether the UPDATE's SET clauses actually
// mutate (a column of) the specified table. JOIN-only references that
// do not appear as a SET-clause qualifier do NOT count as mutation.
//
// Resolution strategy mirrors extractUpdateColumns: build the alias map
// (lowercased), then for each SET assignment resolve the
// qualifier through the map. Unqualified columns are conservatively
// counted as mutation if the target table is in scope (multi-table
// UPDATE with unqualified col is ambiguous; safer to over-include than
// miss).
func updateMutatesTable(stmt *ast.UpdateStmt, database, table string, targetCols map[string]bool) bool {
	// Precondition: target must be in stmt.Tables at all (else there's
	// nothing to discuss).
	targetInScope := false
	for _, expr := range stmt.Tables {
		if tableExprReferences(expr, database, table) {
			targetInScope = true
			break
		}
	}
	if !targetInScope {
		return false
	}

	singleTables := extractSingleTablesFromTableExprs(database, stmt.Tables)
	for _, assignment := range stmt.SetList {
		col := assignment.Column
		if col == nil {
			continue
		}
		if col.Schema != "" && !strings.EqualFold(col.Schema, database) {
			continue
		}
		if col.Table == "" {
			// Unqualified column. Resolve against the target table's
			// actual columns — only count as mutation if the column
			// exists on the target: for `UPDATE test JOIN t1 ON ... SET
			// name = 1` where `name` is on t1 but not test, counting it
			// would produce an empty `... ON DUPLICATE KEY UPDATE ;`.
			if targetCols[strings.ToLower(col.Column)] {
				return true
			}
			continue
		}
		// Qualified — resolve qualifier through the alias map. Both
		// Database AND Table must match: for cross-database joins with
		// homonymous tables (e.g. `UPDATE db1.test t1 JOIN db2.test t2
		// SET t2.a = ...`), the alias t2 resolves to db2.test; without
		// the Database check, a backup item targeting db1.test would
		// incorrectly match db2.test's SET assignments.
		if entry, ok := singleTables[strings.ToLower(col.Table)]; ok && strings.EqualFold(entry.Database, database) && strings.EqualFold(entry.Table, table) {
			return true
		}
		// Fallback: qualifier IS the bare table name (no alias used).
		// The col.Schema check at the top of the loop already filtered
		// out assignments whose explicit schema doesn't match `database`,
		// so reaching here implies col.Schema is empty or matches.
		if strings.EqualFold(col.Table, table) {
			return true
		}
	}
	return false
}

func tableExprReferences(expr ast.TableExpr, database, table string) bool {
	switch n := expr.(type) {
	case *ast.TableRef:
		db := n.Schema
		if db == "" {
			db = database
		}
		// Both comparisons are case-insensitive — TiDB/MySQL identifier
		// comparisons are typically case-insensitive in practice, and the
		// table-name comparison below already uses EqualFold; using == on
		// the database side would inconsistently miss schema-qualified
		// references like `DB.test` when backupItem stores `db`.
		return strings.EqualFold(db, database) && strings.EqualFold(n.Name, table)
	case *ast.JoinClause:
		return tableExprReferences(n.Left, database, table) || tableExprReferences(n.Right, database, table)
	}
	return false
}

// changedUniqueKey returns the name of a unique key whose old values the
// restore could find in another row, or "". A key's value can move to
// another row only when more than one statement changes the columns it
// reads, through its parts, generated columns or expressions, or when the
// one that does sets them to anything but a literal, which every row it
// changes then shares, or when the server changes them ON UPDATE.
func changedUniqueKey(tableMetadata *metadata.TableMetadata, updateColumns []string, assignments [][]*ast.Assignment) string {
	changed := make(map[string]bool)
	for _, column := range updateColumns {
		changed[strings.ToLower(column)] = true
	}
	columns := make(map[string]*metadata.ColumnMetadata)
	for _, column := range tableMetadata.Columns {
		columns[strings.ToLower(column.GetName())] = column
	}
	for _, index := range tableMetadata.Indexes {
		if !index.Primary && !index.Unique {
			continue
		}
		reads := make(map[string]bool)
		for _, part := range index.Expressions {
			key := strings.ToLower(part)
			expression := part
			if column, ok := columns[key]; ok {
				if changed[key] {
					reads[key] = true
				}
				if column.GetGeneration() == nil {
					continue
				}
				expression = column.GetGeneration().GetExpression()
			}
			for _, token := range parser.Tokenize(expression) {
				if name := strings.ToLower(token.Str); changed[name] {
					reads[name] = true
				}
			}
		}
		if len(reads) == 0 {
			continue
		}
		statements := 0
		assigned := make(map[string]bool)
		for _, list := range assignments {
			touches := false
			for _, assignment := range list {
				name := strings.ToLower(assignment.Column.Column)
				if !reads[name] {
					continue
				}
				touches = true
				assigned[name] = true
				if !isLiteral(assignment.Value) {
					return index.Name
				}
			}
			if touches {
				statements++
			}
		}
		if statements > 1 || len(assigned) < len(reads) {
			return index.Name
		}
	}
	return ""
}

// isLiteral reports whether expr is a literal.
func isLiteral(expr ast.ExprNode) bool {
	switch expr.(type) {
	case *ast.IntLit, *ast.FloatLit, *ast.StringLit, *ast.BoolLit, *ast.NullLit, *ast.HexLit, *ast.BitLit:
		return true
	}
	return false
}
