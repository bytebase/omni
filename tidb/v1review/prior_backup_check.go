package v1review

import (
	"sort"
	"strconv"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/v1review"
)

// maxMixedDMLCount is the number of DML statements TiDB's backup transformer
// handles one by one: above it the statements must all change one table, since
// they are grouped into a single backup. It mirrors tidb/analysis's unexported
// limit, which the transformer enforces.
const maxMixedDMLCount = 5

// The two kinds of change a top-level DML makes. V1 kept the kind as the
// statement's name, which it also matched its DDL kinds against.
const (
	priorBackupUpdate = "UPDATE"
	priorBackupDelete = "DELETE"
)

// priorBackupTable is a table an UPDATE or DELETE changes.
type priorBackupTable struct {
	database string
	table    string
}

// priorBackupDML is one change the reviewed statements make, and its kind.
type priorBackupDML struct {
	table priorBackupTable
	kind  string
}

// PriorBackupCheck reports what keeps a prior backup from taking the
// statements, in the V1 advisor's order: the reviewed SQL is larger than facts
// allows; each statement is DDL; the database the backup writes to is missing,
// which does not end the check; one table is changed by more than one kind of
// DML; and the batch changes rows in more than maxMixedDMLCount statements
// across more than one table. V1 reported every DDL statement, whether or not
// the batch changes rows as well, and anchored it at the statement's first
// token, so a DDL finding carries the index of its statement and no range.
//
// database is the database V1 resolved an unqualified table reference to: the
// reviewed database. Tables are grouped as "database.table", both parts
// lower-cased, which is also how V1 printed the table. The mixed DML tables
// come in sorted key order, as V1 printed them; a batch that changes rows
// through an INSERT is not grouped, since V1 collected only UPDATE and DELETE.
//
// Two judgments here are not V1's. A reference to a WITH clause's CTE is not a
// table and does not count as a change, where V1 grouped the reference under a
// table of the CTE's name and the backup transformer drops it. A multi-table
// DELETE's target is resolved through its USING clause to the table that
// clause binds the name to, where V1 grouped the name as written, so an alias
// grouped under a table of the alias's name. Both made this rule disagree with
// the transform it predicts: this one approved a change the transform rejects,
// and reported tables that do not exist.
//
// Ported from Bytebase's V1 rule BUILTIN_PRIOR_BACKUP_CHECK; Bytebase calls it
// only when the review enables prior backup.
func PriorBackupCheck(stmts []ast.Node, database string, schema v1review.Schema, facts v1review.PriorBackupFacts) []v1review.Finding {
	var findings []v1review.Finding

	if facts.TotalSize > facts.MaxSize {
		findings = append(findings, v1review.Finding{Code: v1review.PriorBackupScriptTooLarge, Value: strconv.Itoa(facts.MaxSize)})
	}

	for i, stmt := range stmts {
		if ast.IsDDL(stmt) {
			findings = append(findings, v1review.Finding{Code: v1review.PriorBackupMixedDDLAndDML, Statement: i})
		}
	}

	if !facts.DatabaseExists {
		findings = append(findings, v1review.Finding{Code: v1review.PriorBackupDatabaseMissing, Name: facts.Database})
	}

	var changes []priorBackupDML
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.UpdateStmt:
			// A SET column names the table the UPDATE changes, not the tables
			// its joins read.
			aliases := newUpdateTableAliases(n.Tables, database, cteNames(n.CTEs))
			for _, table := range updateTargets(n.SetList, aliases, schema.Metadata) {
				changes = append(changes, priorBackupDML{table: table, kind: priorBackupUpdate})
			}
		case *ast.DeleteStmt:
			// A DELETE changes the tables it deletes from; its USING clause
			// only filters.
			for _, table := range deleteTargets(n, database) {
				changes = append(changes, priorBackupDML{table: table, kind: priorBackupDelete})
			}
		default:
		}
	}

	kinds := make(map[string]map[string]bool, len(changes))
	for _, change := range changes {
		key := priorBackupTableKey(change.table)
		if kinds[key] == nil {
			kinds[key] = make(map[string]bool, 2)
		}
		kinds[key][change.kind] = true
	}
	keys := make([]string, 0, len(kinds))
	for key := range kinds {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if len(kinds[key]) > 1 {
			findings = append(findings, v1review.Finding{Code: v1review.PriorBackupMixedDMLOnTable, Table: key})
		}
	}

	if len(changes) > maxMixedDMLCount && len(kinds) > 1 {
		findings = append(findings, v1review.Finding{Code: v1review.PriorBackupTooManyDMLStatements, Value: strconv.Itoa(maxMixedDMLCount)})
	}

	return findings
}

// priorBackupTableKey returns the text V1 grouped and printed a table as.
func priorBackupTableKey(table priorBackupTable) string {
	return strings.ToLower(table.database) + "." + strings.ToLower(table.table)
}

// cteNames returns the names a WITH clause declares, lower-cased: TiDB
// resolves a reference to a CTE case-insensitively.
func cteNames(ctes []*ast.CommonTableExpr) map[string]bool {
	if len(ctes) == 0 {
		return nil
	}
	names := make(map[string]bool, len(ctes))
	for _, cte := range ctes {
		if cte != nil {
			names[strings.ToLower(cte.Name)] = true
		}
	}
	return names
}

// isCTEReference reports whether a table reference names a CTE of the
// statement holding it rather than a table. A CTE reference is never
// schema-qualified, so a qualified name is a table even when it matches one.
func isCTEReference(ref *ast.TableRef, ctes map[string]bool) bool {
	return ref.Schema == "" && ctes[strings.ToLower(ref.Name)]
}

// deleteTables returns the tables a table list names, an unqualified name
// resolved to database. A reference to a CTE and a derived table name none.
func deleteTables(tables []ast.TableExpr, database string, ctes map[string]bool) []priorBackupTable {
	var result []priorBackupTable
	for _, table := range tables {
		result = append(result, tableExprRefs(table, database, ctes)...)
	}
	return result
}

// tableExprRefs returns the base tables one table expression names: the table
// a reference names, both sides of a join, and nothing for a derived table.
func tableExprRefs(expr ast.TableExpr, database string, ctes map[string]bool) []priorBackupTable {
	switch n := expr.(type) {
	case *ast.TableRef:
		if isCTEReference(n, ctes) {
			return nil
		}
		if n.Schema != "" {
			database = n.Schema
		}
		return []priorBackupTable{{database: database, table: n.Name}}
	case *ast.JoinClause:
		left := tableExprRefs(n.Left, database, ctes)
		right := tableExprRefs(n.Right, database, ctes)
		return append(left, right...)
	case *ast.SubqueryExpr:
		return nil
	default:
		return nil
	}
}

// deleteTargets returns the tables a DELETE deletes from. A multi-table DELETE
// names its targets by the aliases or names its USING clause declares, so each
// target is resolved there to the table that clause binds it to; a name the
// clause does not declare is kept as written.
func deleteTargets(stmt *ast.DeleteStmt, database string) []priorBackupTable {
	ctes := cteNames(stmt.CTEs)
	targets := deleteTables(stmt.Tables, database, ctes)
	if len(stmt.Using) == 0 {
		return targets
	}
	for i, target := range targets {
		if table, ok := deleteTargetTable(target.table, stmt.Using, database, ctes); ok {
			targets[i] = table
		}
	}
	return targets
}

// deleteTargetTable returns the table a multi-table DELETE's USING clause
// binds name to, by an alias or by a table's own name. ok is false when the
// clause declares no such name.
func deleteTargetTable(name string, using []ast.TableExpr, database string, ctes map[string]bool) (priorBackupTable, bool) {
	for _, expr := range using {
		if table, ok := tableExprByName(name, expr, database, ctes); ok {
			return table, true
		}
	}
	return priorBackupTable{}, false
}

func tableExprByName(name string, expr ast.TableExpr, database string, ctes map[string]bool) (priorBackupTable, bool) {
	switch n := expr.(type) {
	case *ast.TableRef:
		if isCTEReference(n, ctes) {
			return priorBackupTable{}, false
		}
		if !strings.EqualFold(n.Name, name) && (n.Alias == "" || !strings.EqualFold(n.Alias, name)) {
			return priorBackupTable{}, false
		}
		if n.Schema != "" {
			database = n.Schema
		}
		return priorBackupTable{database: database, table: n.Name}, true
	case *ast.JoinClause:
		if table, ok := tableExprByName(name, n.Left, database, ctes); ok {
			return table, true
		}
		return tableExprByName(name, n.Right, database, ctes)
	case *ast.SubqueryExpr:
		return priorBackupTable{}, false
	default:
		return priorBackupTable{}, false
	}
}

// updateTableAliases is what an UPDATE's table list says about the SET
// columns that name their table instead of a column list: the tables keyed by
// their schema-qualified name, the tables keyed by their alias, the tables
// keyed by their bare name, and the distinct tables in the list.
type updateTableAliases struct {
	bySchemaName  map[string]priorBackupTable
	byAlias       map[string]priorBackupTable
	byBareName    map[string][]priorBackupTable
	distinctBases []priorBackupTable
}

// newUpdateTableAliases reads an UPDATE's table list, an unqualified name
// resolved to database.
func newUpdateTableAliases(tables []ast.TableExpr, database string, ctes map[string]bool) *updateTableAliases {
	m := &updateTableAliases{
		bySchemaName: make(map[string]priorBackupTable),
		byAlias:      make(map[string]priorBackupTable),
		byBareName:   make(map[string][]priorBackupTable),
	}
	for _, table := range tables {
		collectTableAliases(table, database, ctes, m)
	}
	return m
}

func collectTableAliases(expr ast.TableExpr, database string, ctes map[string]bool, m *updateTableAliases) {
	switch n := expr.(type) {
	case *ast.TableRef:
		if isCTEReference(n, ctes) {
			return
		}
		if n.Schema != "" {
			database = n.Schema
		}
		base := priorBackupTable{database: database, table: n.Name}
		name := strings.ToLower(n.Name)
		key := strings.ToLower(database) + "." + name
		m.bySchemaName[key] = base
		if n.Alias != "" {
			m.byAlias[strings.ToLower(n.Alias)] = base
		}
		m.byBareName[name] = append(m.byBareName[name], base)
		for _, tracked := range m.distinctBases {
			if priorBackupTableKey(tracked) == key {
				return
			}
		}
		m.distinctBases = append(m.distinctBases, base)
	case *ast.JoinClause:
		collectTableAliases(n.Left, database, ctes, m)
		collectTableAliases(n.Right, database, ctes, m)
	case *ast.SubqueryExpr:
	default:
	}
}

// updateTargets returns the distinct tables an UPDATE's SET columns change, in
// the order they first name one. A column qualified by a schema and a table is
// looked up by both; one naming only a table is looked up by alias first, then
// by bare name, and is skipped where that names several tables. A column with
// no qualifier belongs to the UPDATE's only table, or, when it has several, to
// the tables the snapshot says hold that column.
func updateTargets(setList []*ast.Assignment, m *updateTableAliases, snapshot *metadata.DatabaseSchemaMetadata) []priorBackupTable {
	var result []priorBackupTable
	seen := make(map[string]bool)
	add := func(table priorBackupTable) {
		key := priorBackupTableKey(table)
		if seen[key] {
			return
		}
		seen[key] = true
		result = append(result, table)
	}
	for _, assignment := range setList {
		if assignment == nil || assignment.Column == nil {
			continue
		}
		column := assignment.Column
		switch {
		case column.Schema != "" && column.Table != "":
			if base, ok := m.bySchemaName[strings.ToLower(column.Schema)+"."+strings.ToLower(column.Table)]; ok {
				add(base)
			}
		case column.Table != "":
			if base, ok := m.byAlias[strings.ToLower(column.Table)]; ok {
				add(base)
				continue
			}
			if bases := m.byBareName[strings.ToLower(column.Table)]; len(bases) == 1 {
				add(bases[0])
			}
		default:
			if len(m.distinctBases) == 1 {
				add(m.distinctBases[0])
				continue
			}
			for _, table := range resolveUnqualifiedColumn(column.Column, m.distinctBases, snapshot) {
				add(table)
			}
		}
	}
	return result
}

// resolveUnqualifiedColumn returns the tables a SET column with no qualifier
// belongs to: those of bases the snapshot says hold a column of that name,
// read lower-cased as MySQL compares identifiers. A cross-database base is
// looked up only when the snapshot is of that database. No snapshot, no
// matching column, or an empty column name falls back to every base, since V1
// treated that as a change to all of them.
func resolveUnqualifiedColumn(column string, bases []priorBackupTable, snapshot *metadata.DatabaseSchemaMetadata) []priorBackupTable {
	if column == "" || snapshot == nil {
		return bases
	}
	var matches []priorBackupTable
	for _, base := range bases {
		if base.database != "" && snapshot.GetName() != "" && !strings.EqualFold(base.database, snapshot.GetName()) {
			continue
		}
		for _, schema := range snapshot.GetSchemas() {
			for _, table := range schema.GetTables() {
				if !strings.EqualFold(table.GetName(), base.table) {
					continue
				}
				for _, tableColumn := range table.GetColumns() {
					if strings.EqualFold(tableColumn.GetName(), column) {
						matches = append(matches, base)
						break
					}
				}
			}
		}
	}
	if len(matches) > 0 {
		return matches
	}
	return bases
}
