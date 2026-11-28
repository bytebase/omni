package analysis

import (
	"slices"
	"strings"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/pg/ast"
)

// ChangedResource is a table a statement writes.
type ChangedResource = analysis.ChangedResource

// DMLStatement is a DML statement whose affected rows Bytebase estimates.
type DMLStatement = analysis.DMLStatement

// Changes is what a script changes.
type Changes = analysis.Changes

// ChangeStatement is one statement of a script, parsed on its own: Node's
// locations are byte offsets into Text, as when each split segment is parsed
// separately.
type ChangeStatement struct {
	Node ast.Node
	Text string
}

// ChangeMetadata is the synced metadata of the database a script changes.
// Its lookups follow the metadata's own name matching.
type ChangeMetadata interface {
	// HasSchema reports whether the database has the schema.
	HasSchema(name string) bool
	// RelationSchema returns the first schema in searchPath that has a table
	// (or a partition), view, materialized view or external table named name,
	// or "".
	RelationSchema(searchPath []string, name string) string
	// IndexTable returns the schema of the first index named name in
	// searchPath's schemas and the table it indexes; ok is false when there
	// is no such index or it names no table.
	IndexTable(searchPath []string, name string) (schema, table string, ok bool)
}

// ExtractChangedResources returns the tables a script writes and the DML
// statements to estimate, replaying the search path changes it makes: SET,
// RESET, DISCARD ALL and the transactions that scope SET LOCAL. database is
// the current database and searchPath the initial search path; an empty one
// means public. EXPLAIN ANALYZE counts as the statement it runs.
func ExtractChangedResources(stmts []ChangeStatement, database string, searchPath []string, metadata ChangeMetadata) Changes {
	var changes Changes
	if len(searchPath) == 0 {
		searchPath = []string{"public"}
	}

	paths := newSearchPathTracker(searchPath)
	sampleDML := func(text string) {
		changes.DMLCount++
		text = strings.TrimSpace(text)
		if !strings.HasSuffix(text, ";") {
			text += ";"
		}
		statement := DMLStatement{Text: text}
		// The EXPLAIN connection starts with the initial search path, so a changed one is replayed.
		if !slices.Equal(searchPath, paths.initial) {
			statement.SearchPath = searchPath
		}
		changes.DMLStatements = append(changes.DMLStatements, statement)
	}
	addTable := func(db, schema, table string, affectsData bool) {
		changes.Resources = append(changes.Resources, ChangedResource{Database: db, Schema: schema, Table: table, AffectsData: affectsData})
	}
	addTarget := func(rv *ast.RangeVar, affectedTable bool) {
		db, schema, table := extractExistingRangeVarNames(rv, database, searchPath, metadata)
		addTable(db, schema, table, affectedTable)
	}
	// addDataModifyingCTETargets reports whether the WITH clause modifies data.
	addDataModifyingCTETargets := func(with *ast.WithClause) bool {
		targets := getDataModifyingCTETargets(with)
		for _, rv := range targets {
			addTarget(rv, false)
		}
		return len(targets) > 0
	}

	for _, stmt := range stmts {
		if stmt.Node == nil {
			continue
		}
		node, text := UnwrapExplainAnalyze(stmt.Node, stmt.Text)

		paths.apply(node)
		searchPath = paths.current

		switch n := node.(type) {
		case *ast.CreateStmt:
			if n.Relation != nil {
				db, schema, table := extractNewRangeVarNames(n.Relation, database, searchPath, metadata)
				addTable(db, schema, table, false)
			}

		case *ast.DropStmt:
			objType := ast.ObjectType(n.RemoveType)
			switch objType {
			case ast.OBJECT_INDEX:
				handleDropIndex(n, database, searchPath, metadata, addTable)
			case ast.OBJECT_TABLE, ast.OBJECT_MATVIEW:
				handleDropTable(n, database, searchPath, metadata, addTable)
			default:
			}

		case *ast.AlterTableStmt:
			if ast.ObjectType(n.ObjType) == ast.OBJECT_VIEW {
				continue
			}
			if n.Relation != nil {
				db, schema, table := extractExistingRangeVarNames(n.Relation, database, searchPath, metadata)
				addTable(db, schema, table, true)
			}

		case *ast.RenameStmt:
			if n.Relation == nil {
				continue
			}
			// Newname names a table only when the table itself is renamed.
			switch n.RenameType {
			case ast.OBJECT_TABLE, ast.OBJECT_MATVIEW:
				db, schema, oldTableName := extractExistingRangeVarNames(n.Relation, database, searchPath, metadata)
				addTable(db, schema, oldTableName, true)
				if n.Newname != "" {
					addTable(db, schema, n.Newname, false)
				}
			case ast.OBJECT_COLUMN, ast.OBJECT_TABCONSTRAINT:
				if n.RelationType == ast.OBJECT_TABLE || n.RelationType == ast.OBJECT_MATVIEW {
					addTarget(n.Relation, true)
				}
			case ast.OBJECT_TRIGGER, ast.OBJECT_RULE, ast.OBJECT_POLICY:
				addTarget(n.Relation, true)
			default:
			}

		case *ast.IndexStmt:
			if n.Relation != nil {
				db, schema, table := extractExistingRangeVarNames(n.Relation, database, searchPath, metadata)
				addTable(db, schema, table, false)
			}

		case *ast.InsertStmt:
			hasDataModifyingCTE := addDataModifyingCTETargets(n.WithClause)
			if n.Relation != nil {
				addTarget(n.Relation, false)
			}
			// A data-modifying CTE makes the whole statement a sample, whose estimate already
			// includes the VALUES rows.
			if rows, ok := getInsertValuesRowCount(n); ok && !hasDataModifyingCTE {
				changes.InsertCount += rows
			} else {
				sampleDML(text)
			}

		case *ast.UpdateStmt:
			addDataModifyingCTETargets(n.WithClause)
			if n.Relation != nil {
				addTarget(n.Relation, false)
			}
			sampleDML(text)

		case *ast.DeleteStmt:
			addDataModifyingCTETargets(n.WithClause)
			if n.Relation != nil {
				addTarget(n.Relation, false)
			}
			sampleDML(text)

		case *ast.TruncateStmt:
			if n.Relations != nil {
				for _, item := range n.Relations.Items {
					rv, ok := item.(*ast.RangeVar)
					if !ok {
						continue
					}
					db, schema, table := extractExistingRangeVarNames(rv, database, searchPath, metadata)
					addTable(db, schema, table, true)
				}
			}

		case *ast.MergeStmt:
			addDataModifyingCTETargets(n.WithClause)
			if n.Relation != nil {
				addTarget(n.Relation, false)
			}
			sampleDML(text)

		case *ast.CreateTableAsStmt:
			if n.Into != nil && n.Into.Rel != nil {
				db, schema, table := extractNewRangeVarNames(n.Into.Rel, database, searchPath, metadata)
				addTable(db, schema, table, false)
			}
			if query, ok := n.Query.(*ast.SelectStmt); ok && addDataModifyingCTETargets(query.WithClause) {
				sampleDML(text)
			}

		case *ast.SelectStmt:
			// SELECT ... INTO target_table is a write (classified DDL via the INTO clause,
			// which may sit on the first arm of a set operation).
			if into := intoClause(n); into != nil && into.Rel != nil {
				db, schema, table := extractNewRangeVarNames(into.Rel, database, searchPath, metadata)
				addTable(db, schema, table, false)
			}
			if addDataModifyingCTETargets(n.WithClause) {
				sampleDML(text)
			}

		default:
		}
	}

	return changes
}

// getInsertValuesRowCount returns the rows inserted by INSERT ... VALUES or INSERT ... DEFAULT
// VALUES; ok is false when a query supplies the rows.
func getInsertValuesRowCount(n *ast.InsertStmt) (rows int, ok bool) {
	if n.SelectStmt == nil {
		return 1, true
	}
	if sel, isSelect := n.SelectStmt.(*ast.SelectStmt); isSelect && sel.ValuesLists != nil {
		return len(sel.ValuesLists.Items), true
	}
	return 0, false
}

// getDataModifyingCTETargets returns the tables that INSERT, UPDATE, DELETE, and MERGE statements in
// a WITH clause modify. PostgreSQL only allows data-modifying statements in a top-level WITH.
func getDataModifyingCTETargets(with *ast.WithClause) []*ast.RangeVar {
	if with == nil || with.Ctes == nil {
		return nil
	}
	var targets []*ast.RangeVar
	for _, item := range with.Ctes.Items {
		cte, ok := item.(*ast.CommonTableExpr)
		if !ok {
			continue
		}
		var target *ast.RangeVar
		switch query := cte.Ctequery.(type) {
		case *ast.InsertStmt:
			target = query.Relation
		case *ast.UpdateStmt:
			target = query.Relation
		case *ast.DeleteStmt:
			target = query.Relation
		case *ast.MergeStmt:
			target = query.Relation
		default:
			continue
		}
		if target != nil {
			targets = append(targets, target)
		}
	}
	return targets
}

// extractRangeVarNames extracts database, schema, table from a RangeVar with defaults.
func extractRangeVarNames(rv *ast.RangeVar, defaultDB string, searchPath []string) (string, string, string) {
	db := rv.Catalogname
	schema := rv.Schemaname
	table := rv.Relname
	if db == "" {
		db = defaultDB
	}
	if schema == "" && len(searchPath) > 0 {
		schema = searchPath[0]
	}
	return db, schema, table
}

func extractExistingRangeVarNames(rv *ast.RangeVar, defaultDB string, searchPath []string, metadata ChangeMetadata) (string, string, string) {
	db := rv.Catalogname
	schema := rv.Schemaname
	table := rv.Relname
	if db == "" {
		db = defaultDB
	}
	if schema == "" {
		schemaName := metadata.RelationSchema(searchPath, table)
		if schemaName != "" {
			schema = schemaName
		} else if len(searchPath) > 0 {
			schema = searchPath[0]
		}
	}
	return db, schema, table
}

func extractNewRangeVarNames(rv *ast.RangeVar, defaultDB string, searchPath []string, metadata ChangeMetadata) (string, string, string) {
	db, schema, table := extractRangeVarNames(rv, defaultDB, searchPath)
	if rv.Schemaname == "" && len(searchPath) > 1 && !metadata.HasSchema(schema) {
		schema = ""
	}
	return db, schema, table
}

// handleDropTable handles DROP TABLE/MATERIALIZED VIEW.
func handleDropTable(n *ast.DropStmt, database string, searchPath []string, metadata ChangeMetadata, addTable func(db, schema, table string, affectsData bool)) {
	if n.Objects == nil {
		return
	}
	for _, item := range n.Objects.Items {
		nameList, ok := item.(*ast.List)
		if !ok {
			continue
		}
		db, schema, name := extractNameListParts(nameList, database)
		if schema == "" {
			schemaName := metadata.RelationSchema(searchPath, name)
			if schemaName == "" {
				if len(searchPath) > 0 {
					schema = searchPath[0]
				}
			} else {
				schema = schemaName
			}
		}
		addTable(db, schema, name, true)
	}
}

// handleDropIndex handles DROP INDEX.
func handleDropIndex(n *ast.DropStmt, database string, searchPath []string, metadata ChangeMetadata, addTable func(db, schema, table string, affectsData bool)) {
	if n.Objects == nil {
		return
	}
	for _, item := range n.Objects.Items {
		nameList, ok := item.(*ast.List)
		if !ok {
			continue
		}
		db, schema, indexName := extractNameListParts(nameList, database)

		lookupPath := searchPath
		if schema != "" {
			lookupPath = []string{schema}
		}
		if schemaName, table, ok := metadata.IndexTable(lookupPath, indexName); ok && schemaName != "" {
			addTable(db, schemaName, table, false)
		}
	}
}

// extractNameListParts extracts db, schema, name from a list of String nodes.
func extractNameListParts(nameList *ast.List, defaultDB string) (string, string, string) {
	var parts []string
	for _, item := range nameList.Items {
		if s, ok := item.(*ast.String); ok {
			parts = append(parts, s.Str)
		}
	}
	switch len(parts) {
	case 1:
		return defaultDB, "", parts[0]
	case 2:
		return defaultDB, parts[0], parts[1]
	case 3:
		return parts[0], parts[1], parts[2]
	default:
		return defaultDB, "", ""
	}
}
