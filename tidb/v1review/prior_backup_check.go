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
			aliases := newUpdateTableAliases(n.Tables, database)
			for _, table := range updateTargets(n.SetList, aliases, schema.Metadata) {
				changes = append(changes, priorBackupDML{table: table, kind: priorBackupUpdate})
			}
		case *ast.DeleteStmt:
			// A DELETE's tables are the tables it deletes from; its USING
			// clause only filters. V1 read the names as written, so a
			// multi-table DELETE that targets an alias ("DELETE b FROM t AS b")
			// groups under the alias rather than the table the USING clause
			// resolves it to.
			for _, table := range deleteTables(n.Tables, database) {
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

// deleteTables returns the tables a DELETE's table list names, an unqualified
// name resolved to database. A derived table names none.
func deleteTables(tables []ast.TableExpr, database string) []priorBackupTable {
	var result []priorBackupTable
	for _, table := range tables {
		result = append(result, tableExprRefs(table, database)...)
	}
	return result
}

// tableExprRefs returns the base tables one table expression names: the table
// a reference names, both sides of a join, and nothing for a derived table.
func tableExprRefs(expr ast.TableExpr, database string) []priorBackupTable {
	switch n := expr.(type) {
	case *ast.TableRef:
		if n.Schema != "" {
			database = n.Schema
		}
		return []priorBackupTable{{database: database, table: n.Name}}
	case *ast.JoinClause:
		left := tableExprRefs(n.Left, database)
		right := tableExprRefs(n.Right, database)
		return append(left, right...)
	case *ast.SubqueryExpr:
		return nil
	default:
		return nil
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
func newUpdateTableAliases(tables []ast.TableExpr, database string) *updateTableAliases {
	m := &updateTableAliases{
		bySchemaName: make(map[string]priorBackupTable),
		byAlias:      make(map[string]priorBackupTable),
		byBareName:   make(map[string][]priorBackupTable),
	}
	for _, table := range tables {
		collectTableAliases(table, database, m)
	}
	return m
}

func collectTableAliases(expr ast.TableExpr, database string, m *updateTableAliases) {
	switch n := expr.(type) {
	case *ast.TableRef:
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
		collectTableAliases(n.Left, database, m)
		collectTableAliases(n.Right, database, m)
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
