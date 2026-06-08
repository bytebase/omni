package redshift

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bytebase/omni/redshift/ast"
)

// ChangeKind describes the top-level kind of table change.
type ChangeKind string

const (
	ChangeKindCreate ChangeKind = "CREATE"
	ChangeKindAlter  ChangeKind = "ALTER"
	ChangeKindDrop   ChangeKind = "DROP"
	ChangeKindDML    ChangeKind = "DML"
)

// ChangeSummary is an Omni-native Redshift changed-resource summary.
type ChangeSummary struct {
	Tables      []ChangedTable
	DMLCount    int
	InsertCount int
	UpdateCount int
	DeleteCount int
}

// ChangedTable identifies a changed table-like resource.
type ChangedTable struct {
	Database string
	Schema   string
	Name     string
	Kind     ChangeKind
	Affected bool
}

// ExtractChangedResources extracts table-level resource changes from Redshift SQL.
func ExtractChangedResources(sql, currentDatabase, currentSchema string) (*ChangeSummary, error) {
	stmts, err := Parse(sql)
	if err != nil {
		return nil, err
	}
	if currentSchema == "" {
		currentSchema = "public"
	}

	summary := &ChangeSummary{}
	tableByKey := make(map[string]ChangedTable)
	searchPath := currentSchema

	for _, stmt := range stmts {
		if stmt.Empty() {
			continue
		}
		switch n := stmt.AST.(type) {
		case *ast.SelectStmt:
			if n.IntoClause != nil {
				addChangedRangeVar(tableByKey, n.IntoClause.Rel, currentDatabase, searchPath, ChangeKindCreate, false)
			}
		case *ast.CreateStmt:
			addChangedRangeVar(tableByKey, n.Relation, currentDatabase, searchPath, ChangeKindCreate, false)
		case *ast.AlterTableStmt:
			if n.ObjType == int(ast.OBJECT_TABLE) {
				addChangedRangeVar(tableByKey, n.Relation, currentDatabase, searchPath, ChangeKindAlter, true)
			}
		case *ast.DropStmt:
			if n.RemoveType == int(ast.OBJECT_TABLE) {
				for _, item := range n.Objects.Items {
					addChangedNameParts(tableByKey, nameParts(item), currentDatabase, searchPath, ChangeKindDrop, true)
				}
			}
		case *ast.InsertStmt:
			summary.InsertCount++
			summary.DMLCount++
			addChangedRangeVar(tableByKey, n.Relation, currentDatabase, searchPath, ChangeKindDML, true)
		case *ast.UpdateStmt:
			summary.UpdateCount++
			summary.DMLCount++
			addChangedRangeVar(tableByKey, n.Relation, currentDatabase, searchPath, ChangeKindDML, true)
		case *ast.DeleteStmt:
			summary.DeleteCount++
			summary.DMLCount++
			addChangedRangeVar(tableByKey, n.Relation, currentDatabase, searchPath, ChangeKindDML, true)
		case *ast.VariableSetStmt:
			if schema, ok := searchPathFromSet(n, currentSchema); ok {
				searchPath = schema
			}
		}
	}

	for _, table := range tableByKey {
		summary.Tables = append(summary.Tables, table)
	}
	sort.Slice(summary.Tables, func(i, j int) bool {
		a := summary.Tables[i]
		b := summary.Tables[j]
		return fmt.Sprintf("%s.%s.%s", a.Database, a.Schema, a.Name) < fmt.Sprintf("%s.%s.%s", b.Database, b.Schema, b.Name)
	})
	return summary, nil
}

func addChangedRangeVar(tables map[string]ChangedTable, rv *ast.RangeVar, currentDatabase, currentSchema string, kind ChangeKind, affected bool) {
	if rv == nil || rv.Relname == "" {
		return
	}
	database := rv.Catalogname
	if database == "" {
		database = currentDatabase
	}
	schema := rv.Schemaname
	if schema == "" {
		schema = currentSchema
	}
	addChangedTable(tables, ChangedTable{
		Database: database,
		Schema:   schema,
		Name:     rv.Relname,
		Kind:     kind,
		Affected: affected,
	})
}

func addChangedNameParts(tables map[string]ChangedTable, parts []string, currentDatabase, currentSchema string, kind ChangeKind, affected bool) {
	if len(parts) == 0 {
		return
	}
	table := ChangedTable{
		Database: currentDatabase,
		Schema:   currentSchema,
		Name:     parts[len(parts)-1],
		Kind:     kind,
		Affected: affected,
	}
	switch len(parts) {
	case 2:
		table.Schema = parts[0]
	case 3:
		table.Database = parts[0]
		table.Schema = parts[1]
	}
	addChangedTable(tables, table)
}

func addChangedTable(tables map[string]ChangedTable, table ChangedTable) {
	if table.Name == "" {
		return
	}
	key := strings.Join([]string{table.Database, table.Schema, table.Name}, "\x00")
	if existing, ok := tables[key]; ok {
		if existing.Kind != table.Kind {
			existing.Kind = table.Kind
		}
		existing.Affected = existing.Affected || table.Affected
		tables[key] = existing
		return
	}
	tables[key] = table
}

func nameParts(node ast.Node) []string {
	switch n := node.(type) {
	case *ast.List:
		parts := make([]string, 0, len(n.Items))
		for _, item := range n.Items {
			if s, ok := item.(*ast.String); ok {
				parts = append(parts, s.Str)
			}
		}
		return parts
	case *ast.String:
		return []string{n.Str}
	default:
		return nil
	}
}

func searchPathFromSet(stmt *ast.VariableSetStmt, defaultSchema string) (string, bool) {
	if !strings.EqualFold(stmt.Name, "search_path") {
		return "", false
	}
	if stmt.Kind == ast.VAR_SET_DEFAULT || stmt.Kind == ast.VAR_RESET {
		return defaultSchema, true
	}
	if stmt.Args == nil {
		return "", false
	}
	for _, arg := range stmt.Args.Items {
		c, ok := arg.(*ast.A_Const)
		if !ok {
			continue
		}
		s, ok := c.Val.(*ast.String)
		if !ok {
			continue
		}
		if s.Str == "" || s.Str == "$user" {
			continue
		}
		return s.Str, true
	}
	return "", false
}
