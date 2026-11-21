package analysis

import (
	"strings"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/oracle/ast"
)

// ChangedResource is a resource a statement writes.
type ChangedResource = analysis.ChangedResource

// StatementChanges is what one statement changes.
type StatementChanges = analysis.StatementChanges

// IndexTableFunc returns the table that the index named index in schema
// belongs to, and false when the metadata has no such index.
type IndexTableFunc func(schema, index string) (table string, ok bool)

// ExtractChangedResources returns the resources one statement writes and its
// DML counts. currentSchema is the schema unqualified names are in, text is the
// statement's text, which a DML statement is sampled from, and indexTable
// resolves the table of a dropped index; with a nil indexTable, DROP INDEX
// changes nothing.
func ExtractChangedResources(node ast.Node, currentSchema, text string, indexTable IndexTableFunc) StatementChanges {
	extractor := &changedResourceExtractor{
		currentSchema: currentSchema,
		indexTable:    indexTable,
	}
	extractor.extract(node, text)
	return extractor.changes
}

type changedResourceExtractor struct {
	currentSchema string
	indexTable    IndexTableFunc
	changes       StatementChanges
}

func (e *changedResourceExtractor) extract(node ast.Node, text string) {
	if node == nil {
		return
	}

	switch n := node.(type) {
	case *ast.CreateTableStmt:
		e.addTable(n.Name, false)
	case *ast.DropStmt:
		if n.ObjectType == ast.OBJECT_TABLE {
			for _, name := range objectNameList(n.Names) {
				e.addTable(name, true)
			}
		}
		if n.ObjectType == ast.OBJECT_INDEX {
			for _, name := range objectNameList(n.Names) {
				e.addIndexTable(name)
			}
		}
		switch n.ObjectType {
		case ast.OBJECT_VIEW, ast.OBJECT_SEQUENCE, ast.OBJECT_SYNONYM,
			ast.OBJECT_PROCEDURE, ast.OBJECT_FUNCTION, ast.OBJECT_TRIGGER,
			ast.OBJECT_PACKAGE, ast.OBJECT_PACKAGE_BODY, ast.OBJECT_TYPE,
			ast.OBJECT_TYPE_BODY, ast.OBJECT_MATERIALIZED_VIEW, ast.OBJECT_CLUSTER:
			for _, name := range objectNameList(n.Names) {
				e.addObjectDatabase(name)
			}
		default:
		}
	case *ast.AlterTableStmt:
		affected := true
		for _, cmd := range alterTableCmdList(n.Actions) {
			if cmd.Action == ast.AT_RENAME {
				affected = false
				break
			}
		}
		e.addTable(n.Name, affected)
	case *ast.CreateIndexStmt:
		e.addTable(n.Table, false)
	case *ast.InsertStmt:
		e.extractInsert(n, text)
	case *ast.UpdateStmt:
		e.addTable(n.Table, false)
		e.trackDML(text)
	case *ast.DeleteStmt:
		e.addTable(n.Table, false)
		e.trackDML(text)
	case *ast.MergeStmt:
		e.addTable(n.Target, false)
		e.trackDML(text)
	case *ast.TruncateStmt:
		e.addTable(n.Table, true)
	case *ast.CreateViewStmt:
		e.addObjectDatabase(n.Name)
	case *ast.CreateProcedureStmt:
		e.addObjectDatabase(n.Name)
	case *ast.CreateFunctionStmt:
		e.addObjectDatabase(n.Name)
	case *ast.CreateTriggerStmt:
		e.addObjectDatabase(n.Name)
	case *ast.CreateSequenceStmt:
		e.addObjectDatabase(n.Name)
	case *ast.CreateSynonymStmt:
		e.addObjectDatabase(n.Name)
	case *ast.AlterViewStmt:
		e.addObjectDatabase(n.Name)
	case *ast.AlterSequenceStmt:
		e.addObjectDatabase(n.Name)
	case *ast.AlterProcedureStmt:
		e.addObjectDatabase(n.Name)
	case *ast.AlterFunctionStmt:
		e.addObjectDatabase(n.Name)
	case *ast.AlterTriggerStmt:
		e.addObjectDatabase(n.Name)
	case *ast.AlterSynonymStmt:
		e.addObjectDatabase(n.Name)
	case *ast.CreatePackageStmt:
		e.addObjectDatabase(n.Name)
	case *ast.AlterPackageStmt:
		e.addObjectDatabase(n.Name)
	case *ast.CreateTypeStmt:
		e.addObjectDatabase(n.Name)
	case *ast.AlterTypeStmt:
		e.addObjectDatabase(n.Name)
	case *ast.CreateClusterStmt:
		e.addObjectDatabase(n.Name)
	case *ast.AlterMaterializedViewStmt:
		e.addObjectDatabase(n.Name)
	default:
	}
}

func (e *changedResourceExtractor) extractInsert(stmt *ast.InsertStmt, text string) {
	if stmt.Table != nil {
		e.addTable(stmt.Table, false)
	}
	for _, item := range insertIntoList(stmt.MultiTable) {
		e.addTable(item.Table, false)
	}

	if stmt.Table != nil && stmt.Values != nil {
		e.changes.InsertCount++
		return
	}
	e.trackDML(text)
}

// addObjectDatabase records a database-only write target for a non-table object DDL (view,
// procedure, function, trigger, sequence, synonym) only when its name carries an explicit schema
// qualifier (Oracle's schema maps to a database). An unqualified name records nothing, keeping
// the request-database fallback.
func (e *changedResourceExtractor) addObjectDatabase(name *ast.ObjectName) {
	if name != nil && name.Schema != "" {
		e.changes.Resources = append(e.changes.Resources, ChangedResource{Database: name.Schema, DatabaseOnly: true})
	}
}

func (e *changedResourceExtractor) addTable(name *ast.ObjectName, affected bool) {
	if name == nil || name.Name == "" {
		return
	}
	schema := name.Schema
	if schema == "" {
		schema = e.currentSchema
	}
	e.changes.Resources = append(e.changes.Resources, ChangedResource{Database: schema, Table: name.Name, AffectsData: affected})
}

func (e *changedResourceExtractor) addIndexTable(name *ast.ObjectName) {
	if e.indexTable == nil || name == nil || name.Name == "" {
		return
	}
	schema := name.Schema
	if schema == "" {
		schema = e.currentSchema
	}
	table, ok := e.indexTable(schema, name.Name)
	if !ok {
		return
	}
	e.changes.Resources = append(e.changes.Resources, ChangedResource{Database: schema, Table: table})
}

func (e *changedResourceExtractor) trackDML(text string) {
	e.changes.DMLCount++
	e.changes.DMLStatements = append(e.changes.DMLStatements, statementText(text))
}

func statementText(text string) string {
	text = strings.TrimSpace(text)
	if strings.HasSuffix(text, ";") {
		return text
	}
	return text + ";"
}

func objectNameList(list *ast.List) []*ast.ObjectName {
	if list == nil {
		return nil
	}
	names := make([]*ast.ObjectName, 0, len(list.Items))
	for _, item := range list.Items {
		if name, ok := item.(*ast.ObjectName); ok {
			names = append(names, name)
		}
	}
	return names
}

func alterTableCmdList(list *ast.List) []*ast.AlterTableCmd {
	if list == nil {
		return nil
	}
	cmds := make([]*ast.AlterTableCmd, 0, len(list.Items))
	for _, item := range list.Items {
		if cmd, ok := item.(*ast.AlterTableCmd); ok {
			cmds = append(cmds, cmd)
		}
	}
	return cmds
}

func insertIntoList(list *ast.List) []*ast.InsertIntoClause {
	if list == nil {
		return nil
	}
	clauses := make([]*ast.InsertIntoClause, 0, len(list.Items))
	for _, item := range list.Items {
		if clause, ok := item.(*ast.InsertIntoClause); ok {
			clauses = append(clauses, clause)
		}
	}
	return clauses
}
