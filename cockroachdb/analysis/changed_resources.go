package analysis

import (
	"slices"
	"strings"

	"github.com/bytebase/omni/analysis"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/parser/statements"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
)

// ChangedResource is a table a statement writes.
type ChangedResource = analysis.ChangedResource

// DMLStatement is a DML statement whose affected rows Bytebase estimates.
type DMLStatement struct {
	// Text is the statement, as CockroachDB's parser reports it and without a
	// trailing semicolon.
	Text string
	// SearchPath is the search path the statement runs under when an earlier
	// statement changed it, and nil when it runs under the initial one.
	SearchPath []string
}

// Changes is what a script changes.
type Changes struct {
	// Resources holds the tables the statements write, in the order they name
	// them, duplicates included.
	Resources []ChangedResource
	// Databases holds the databases a statement changes without naming a
	// table, such as an index of another database.
	Databases []string
	// DMLStatements holds the DML statements to estimate, in order.
	DMLStatements []DMLStatement
	// DMLCount counts the DML statements DMLStatements holds.
	DMLCount int
	// InsertCount counts the rows of INSERT ... VALUES and INSERT ... DEFAULT VALUES.
	InsertCount int
}

// ChangeMetadata is the synced metadata of the database a script changes.
// Its lookups follow the metadata's own name matching.
type ChangeMetadata interface {
	// RelationSchema returns the first schema in searchPath that has a table,
	// view or materialized view named name, or "".
	RelationSchema(searchPath []string, name string) string
	// IndexTable returns the schema of the first index named name in
	// searchPath's schemas and the table it indexes; ok is false when there
	// is no such index or it names no table.
	IndexTable(searchPath []string, name string) (schema, table string, ok bool)
}

// ExtractChangedResourcesFromParsed returns the changed resources and DML
// statements of opaque ParsedStatement values, keeping cockroachdb-parser
// types out of callers.
func ExtractChangedResourcesFromParsed(stmts []ParsedStatement, database string, searchPath []string, metadata ChangeMetadata) Changes {
	return ExtractChangedResources(unwrapStatements(stmts), database, searchPath, metadata)
}

// ExtractChangedResources returns the tables a script writes and the DML
// statements to estimate, replaying the search path changes it makes: SET,
// RESET, DISCARD ALL and the transactions that scope SET LOCAL. database is
// the current database and searchPath the initial search path; an empty one
// means public. Like PostgreSQL's extractor, an existing table resolves to the
// first schema that has it in metadata, and a new table to the first schema.
// EXPLAIN ANALYZE counts as the statement it runs.
func ExtractChangedResources(stmts statements.Statements, database string, searchPath []string, metadata ChangeMetadata) Changes {
	if len(searchPath) == 0 {
		searchPath = []string{"public"}
	}
	initialSearchPath := searchPath
	// sessionSearchPath leaves out SET LOCAL, which lasts until the transaction ends.
	sessionSearchPath := searchPath
	// sessionGeneration counts non-local SETs, so restoring a savepoint can tell
	// one made since it even when the value it assigned equals the saved one.
	sessionGeneration := 0
	// savepoints holds the search path each savepoint saw, so ROLLBACK TO SAVEPOINT
	// restores it; a savepoint is scoped to its transaction.
	type searchPathSnapshot struct {
		searchPath        []string
		sessionGeneration int
	}
	savepoints := map[tree.Name]searchPathSnapshot{}
	var savepointOrder []tree.Name
	// restoreSavepoint makes a savepoint's search path current and drops the
	// savepoints nested after it.
	restoreSavepoint := func(name tree.Name) bool {
		index := slices.Index(savepointOrder, name)
		if index < 0 {
			return false
		}
		snapshot := savepoints[name]
		if sessionGeneration == snapshot.sessionGeneration {
			searchPath = snapshot.searchPath
		} else {
			// A non-local SET since the savepoint is not rolled back.
			searchPath = sessionSearchPath
		}
		for _, later := range savepointOrder[index+1:] {
			delete(savepoints, later)
		}
		savepointOrder = savepointOrder[:index+1]
		return true
	}
	changes := Changes{}
	schemaOf := func(table string) string {
		if metadata != nil {
			if schema := metadata.RelationSchema(searchPath, table); schema != "" {
				return schema
			}
		}
		return searchPath[0]
	}
	addTable := func(name *tree.TableName, affectedTable bool) {
		db, schema, table := resolveTableName(name, database, schemaOf(name.Table()))
		changes.Resources = append(changes.Resources, ChangedResource{Database: db, Schema: schema, Table: table, AffectsData: affectedTable})
	}
	addDatabase := func(name string) {
		if !slices.Contains(changes.Databases, name) {
			changes.Databases = append(changes.Databases, name)
		}
	}
	// addIndexTable records the table of an index, which metadata supplies when the name omits it.
	addIndexTable := func(index *tree.TableIndexName) {
		if index.Table.ObjectName != "" {
			addTable(&index.Table, false)
			return
		}
		// metadata describes only the current database, so an index of another is recorded by its database.
		if index.Table.ExplicitCatalog && index.Table.Catalog() != database {
			addDatabase(index.Table.Catalog())
			return
		}
		if metadata == nil {
			return
		}
		schemas := searchPath
		if index.Table.ExplicitSchema {
			schemas = []string{index.Table.Schema()}
		}
		schema, table, ok := metadata.IndexTable(schemas, string(index.Index))
		if !ok {
			return
		}
		db, _, _ := resolveTableName(&index.Table, database, schema)
		changes.Resources = append(changes.Resources, ChangedResource{Database: db, Schema: schema, Table: table})
	}
	// addMutations records the target of each mutation in the statement and returns their number.
	addMutations := func(stmt tree.Statement) int {
		targets := getMutationTargets(stmt)
		for _, target := range targets {
			if name, ok := getTableName(target); ok {
				addTable(name, false)
			}
		}
		return len(targets)
	}
	addDML := func(text string) {
		changes.DMLCount++
		statement := DMLStatement{Text: strings.TrimSpace(text)}
		// The EXPLAIN connection starts with the initial search path, so a changed one is replayed.
		if !slices.Equal(searchPath, initialSearchPath) {
			statement.SearchPath = searchPath
		}
		changes.DMLStatements = append(changes.DMLStatements, statement)
	}

	for _, parsed := range stmts {
		stmt, text := parsed.AST, parsed.SQL
		// EXPLAIN ANALYZE executes the statement it explains.
		if explain, ok := stmt.(*tree.ExplainAnalyze); ok {
			stmt, text = explain.Statement, tree.AsString(explain.Statement)
		}
		switch n := stmt.(type) {
		case *tree.Insert:
			mutationCount := addMutations(n)
			// Nested mutations make the whole statement a sample, whose estimate includes the VALUES rows.
			if rows, ok := getInsertValuesRowCount(n); ok && mutationCount == 1 {
				changes.InsertCount += rows
			} else {
				addDML(text)
			}
		case *tree.Update, *tree.Delete:
			addMutations(n)
			addDML(text)
		case *tree.Select:
			if addMutations(n) > 0 {
				addDML(text)
			}
		case *tree.CreateTable:
			db, schema, table := resolveTableName(&n.Table, database, searchPath[0])
			changes.Resources = append(changes.Resources, ChangedResource{Database: db, Schema: schema, Table: table})
			if n.AsSource != nil && addMutations(n.AsSource) > 0 {
				addDML(text)
			}
		case *tree.CreateView:
			if n.Materialized {
				db, schema, table := resolveTableName(&n.Name, database, searchPath[0])
				changes.Resources = append(changes.Resources, ChangedResource{Database: db, Schema: schema, Table: table})
			}
		case *tree.Import:
			// The plan cannot estimate the rows an import loads, so, as for LOAD DATA in the MySQL family,
			// only its target is recorded.
			if n.Into && n.Table != nil {
				addTable(n.Table, false)
			}
		case *tree.AlterTable:
			name := n.Table.ToTableName()
			addTable(&name, true)
		case *tree.AlterTableLocality:
			name := n.Name.ToTableName()
			addTable(&name, true)
		case *tree.AlterTableSetSchema:
			if !n.IsSequence && (!n.IsView || n.IsMaterialized) {
				name := n.Name.ToTableName()
				addTable(&name, true)
			}
		case *tree.AlterTableOwner:
			if !n.IsSequence && (!n.IsView || n.IsMaterialized) {
				name := n.Name.ToTableName()
				addTable(&name, true)
			}
		case *tree.DropTable:
			for i := range n.Names {
				addTable(&n.Names[i], true)
			}
		case *tree.DropView:
			// A materialized view holds data, so its drop is recorded like a table's.
			if n.IsMaterialized {
				for i := range n.Names {
					addTable(&n.Names[i], true)
				}
			}
		case *tree.Truncate:
			for i := range n.Tables {
				addTable(&n.Tables[i], true)
			}
		case *tree.RenameTable:
			if n.IsSequence || (n.IsView && !n.IsMaterialized) {
				continue
			}
			name := n.Name.ToTableName()
			db, schema, table := resolveTableName(&name, database, schemaOf(name.Table()))
			changes.Resources = append(changes.Resources, ChangedResource{Database: db, Schema: schema, Table: table, AffectsData: true})
			// An unqualified new name stays in the schema of the renamed table.
			newName := n.NewName.ToTableName()
			newDB, newSchema, newTable := resolveTableName(&newName, db, schema)
			changes.Resources = append(changes.Resources, ChangedResource{Database: newDB, Schema: newSchema, Table: newTable})
		case *tree.CreateIndex:
			addTable(&n.Table, false)
		case *tree.DropIndex:
			for _, index := range n.IndexList {
				addIndexTable(index)
			}
		case *tree.AlterIndex:
			addIndexTable(&n.Index)
		case *tree.AlterIndexVisible:
			addIndexTable(&n.Index)
		case *tree.RenameIndex:
			addIndexTable(n.Index)
		case *tree.SetZoneConfig:
			switch {
			case n.TableOrIndex.Index != "":
				addIndexTable(&n.TableOrIndex)
			case n.TableOrIndex.Table.ObjectName != "":
				addTable(&n.TableOrIndex.Table, true)
			default:
			}
		case *tree.BeginTransaction:
			savepoints = map[tree.Name]searchPathSnapshot{}
			savepointOrder = nil
		case *tree.Savepoint:
			if index := slices.Index(savepointOrder, n.Name); index >= 0 {
				savepointOrder = slices.Delete(savepointOrder, index, index+1)
			}
			savepoints[n.Name] = searchPathSnapshot{searchPath: searchPath, sessionGeneration: sessionGeneration}
			savepointOrder = append(savepointOrder, n.Name)
		case *tree.RollbackToSavepoint:
			// CockroachDB rolls a SET LOCAL back to the savepoint; a SET since then stays.
			restoreSavepoint(n.Savepoint)
		case *tree.ReleaseSavepoint:
			// CockroachDB also rolls a SET LOCAL back when the savepoint is released.
			if restoreSavepoint(n.Savepoint) {
				delete(savepoints, n.Savepoint)
				savepointOrder = savepointOrder[:len(savepointOrder)-1]
			}
		case *tree.Discard:
			if n.Mode == tree.DiscardModeAll {
				searchPath, sessionSearchPath = initialSearchPath, initialSearchPath
				sessionGeneration++
			}
		case *tree.CommitTransaction, *tree.RollbackTransaction:
			// Unlike PostgreSQL, CockroachDB keeps a SET made in a transaction that rolls back.
			searchPath = sessionSearchPath
			savepoints = map[tree.Name]searchPathSnapshot{}
			savepointOrder = nil
		case *tree.SetVar:
			if n.ResetAll || strings.EqualFold(n.Name, "search_path") {
				searchPath = getSearchPath(n, initialSearchPath)
				if !n.Local {
					sessionSearchPath = searchPath
					sessionGeneration++
				}
			}
		default:
		}
	}
	return changes
}

// getSearchPath returns the schemas a SET search_path names, skipping "$user", or defaultSearchPath
// when it names none or resets the search path.
func getSearchPath(set *tree.SetVar, defaultSearchPath []string) []string {
	if set.ResetAll {
		return defaultSearchPath
	}
	var searchPath []string
	for _, value := range set.Values {
		var schema string
		switch v := value.(type) {
		case *tree.UnresolvedName:
			schema = v.Parts[0]
		case *tree.StrVal:
			schema = v.RawString()
		default:
			return defaultSearchPath
		}
		if schema != "$user" {
			searchPath = append(searchPath, schema)
		}
	}
	if len(searchPath) == 0 {
		return defaultSearchPath
	}
	return searchPath
}

func resolveTableName(name *tree.TableName, defaultDatabase string, defaultSchema string) (string, string, string) {
	database, schema := defaultDatabase, defaultSchema
	if name.ExplicitCatalog {
		database = name.Catalog()
	}
	if name.ExplicitSchema {
		schema = name.Schema()
	}
	return database, schema, name.Table()
}

// getTableName returns the table name of a mutation target, which is false for a numeric table
// reference.
func getTableName(target tree.TableExpr) (*tree.TableName, bool) {
	if aliased, ok := target.(*tree.AliasedTableExpr); ok {
		target = aliased.Expr
	}
	name, ok := target.(*tree.TableName)
	return name, ok
}

// getInsertValuesRowCount returns the rows of INSERT ... VALUES and INSERT ... DEFAULT VALUES; ok is
// false when a query supplies the rows.
func getInsertValuesRowCount(insert *tree.Insert) (int, bool) {
	if insert.Rows == nil || insert.Rows.Select == nil {
		return 1, true
	}
	// LIMIT and OFFSET make the row count the plan's to estimate, not VALUES'.
	if insert.Rows.Limit != nil {
		return 0, false
	}
	return valuesRowCount(insert.Rows.Select)
}

func valuesRowCount(stmt tree.SelectStatement) (int, bool) {
	switch n := stmt.(type) {
	case *tree.ParenSelect:
		if n.Select == nil || n.Select.Limit != nil {
			return 0, false
		}
		return valuesRowCount(n.Select.Select)
	case *tree.ValuesClause:
		return len(n.Rows), true
	default:
		return 0, false
	}
}

// getMutationTargets returns the target of every INSERT, UPSERT, UPDATE, and DELETE in stmt: stmt
// itself, its CTEs, and statement sources such as [DELETE FROM t RETURNING id] at any depth.
func getMutationTargets(stmt tree.Statement) []tree.TableExpr {
	c := &mutationCollector{}
	c.statement(stmt)
	return c.targets
}

// getMutationTypes returns the statement type of each mutation getMutationTargets finds, in order.
func getMutationTypes(stmt tree.Statement) []StatementType {
	c := &mutationCollector{}
	c.statement(stmt)
	return c.types
}

type mutationCollector struct {
	targets []tree.TableExpr
	types   []StatementType
}

func (c *mutationCollector) statement(stmt tree.Statement) {
	switch n := stmt.(type) {
	case *tree.Insert:
		c.targets = append(c.targets, n.Table)
		c.types = append(c.types, Insert)
		c.with(n.With)
		c.selectQuery(n.Rows)
		if n.OnConflict != nil {
			c.updateExprs(n.OnConflict.Exprs)
			c.where(n.OnConflict.Where)
			c.expr(n.OnConflict.ArbiterPredicate)
		}
		c.returning(n.Returning)
	case *tree.Update:
		c.targets = append(c.targets, n.Table)
		c.types = append(c.types, Update)
		c.with(n.With)
		c.updateExprs(n.Exprs)
		c.tableExprs(n.From)
		c.where(n.Where)
		c.orderBy(n.OrderBy)
		c.limit(n.Limit)
		c.returning(n.Returning)
	case *tree.Delete:
		c.targets = append(c.targets, n.Table)
		c.types = append(c.types, Delete)
		c.with(n.With)
		c.tableExprs(n.Using)
		c.where(n.Where)
		c.orderBy(n.OrderBy)
		c.limit(n.Limit)
		c.returning(n.Returning)
	case *tree.Select:
		c.selectQuery(n)
	case tree.SelectStatement:
		c.selectStatement(n)
	case *tree.ExplainAnalyze:
		// A statement source may hold an EXPLAIN ANALYZE, which runs the
		// statement it explains.
		c.statement(n.Statement)
	default:
	}
}

func (c *mutationCollector) selectQuery(sel *tree.Select) {
	if sel == nil {
		return
	}
	c.with(sel.With)
	c.selectStatement(sel.Select)
	c.orderBy(sel.OrderBy)
	c.limit(sel.Limit)
}

func (c *mutationCollector) selectStatement(stmt tree.SelectStatement) {
	switch n := stmt.(type) {
	case *tree.ParenSelect:
		c.selectQuery(n.Select)
	case *tree.SelectClause:
		for _, selectExpr := range n.Exprs {
			c.expr(selectExpr.Expr)
		}
		c.tableExprs(n.From.Tables)
		c.where(n.Where)
		c.exprs(n.GroupBy)
		c.where(n.Having)
		c.exprs(n.DistinctOn)
		for _, window := range n.Window {
			c.exprs(window.Partitions)
			c.orderBy(window.OrderBy)
		}
	case *tree.UnionClause:
		c.selectQuery(n.Left)
		c.selectQuery(n.Right)
	case *tree.ValuesClause:
		for _, row := range n.Rows {
			c.exprs(row)
		}
	default:
	}
}

func (c *mutationCollector) with(with *tree.With) {
	if with == nil {
		return
	}
	for _, cte := range with.CTEList {
		c.statement(cte.Stmt)
	}
}

func (c *mutationCollector) tableExprs(exprs tree.TableExprs) {
	for _, expr := range exprs {
		c.tableExpr(expr)
	}
}

func (c *mutationCollector) tableExpr(expr tree.TableExpr) {
	switch n := expr.(type) {
	case *tree.AliasedTableExpr:
		c.tableExpr(n.Expr)
	case *tree.ParenTableExpr:
		c.tableExpr(n.Expr)
	case *tree.JoinTableExpr:
		c.tableExpr(n.Left)
		c.tableExpr(n.Right)
		if cond, ok := n.Cond.(*tree.OnJoinCond); ok {
			c.expr(cond.Expr)
		}
	case *tree.RowsFromExpr:
		c.exprs(n.Items)
	case *tree.Subquery:
		c.selectStatement(n.Select)
	case *tree.StatementSource:
		c.statement(n.Statement)
	default:
	}
}

func (c *mutationCollector) expr(expr tree.Expr) {
	if expr == nil {
		return
	}
	// The visit function never fails, so neither does the visit.
	_, _ = tree.SimpleVisit(expr, func(e tree.Expr) (bool, tree.Expr, error) {
		if subquery, ok := e.(*tree.Subquery); ok {
			c.selectStatement(subquery.Select)
			return false, e, nil
		}
		return true, e, nil
	})
}

func (c *mutationCollector) exprs(exprs []tree.Expr) {
	for _, expr := range exprs {
		c.expr(expr)
	}
}

func (c *mutationCollector) where(where *tree.Where) {
	if where != nil {
		c.expr(where.Expr)
	}
}

func (c *mutationCollector) orderBy(orderBy tree.OrderBy) {
	for _, order := range orderBy {
		c.expr(order.Expr)
	}
}

func (c *mutationCollector) limit(limit *tree.Limit) {
	if limit != nil {
		c.expr(limit.Count)
		c.expr(limit.Offset)
	}
}

func (c *mutationCollector) updateExprs(exprs tree.UpdateExprs) {
	for _, updateExpr := range exprs {
		c.expr(updateExpr.Expr)
	}
}

func (c *mutationCollector) returning(returning tree.ReturningClause) {
	if exprs, ok := returning.(*tree.ReturningExprs); ok {
		for _, selectExpr := range *exprs {
			c.expr(selectExpr.Expr)
		}
	}
}
