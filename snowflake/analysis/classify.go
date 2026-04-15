// Package analysis provides utilities for classifying and analyzing parsed
// Snowflake SQL AST nodes.
//
// The primary entry points are Classify, which maps any ast.Node to a
// Category, and ClassifySQL, which parses a raw SQL string and classifies
// its first statement.
//
// See docs/superpowers/specs/2026-04-15-snowflake-stmt-classify-design.md
// for the design rationale.
package analysis

import (
	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/snowflake/parser"
)

// Category is the broad statement category assigned by Classify.
type Category int

const (
	// CategoryUnknown is returned for nil nodes or node types that are not
	// recognized as a known statement kind. It is also returned when parsing
	// fails completely.
	CategoryUnknown Category = iota

	// CategorySelect covers SELECT statements, WITH+SELECT (CTEs), and set
	// operations (UNION / INTERSECT / EXCEPT).
	CategorySelect

	// CategoryDML covers data-manipulation statements: INSERT (single-table and
	// multi-table), UPDATE, DELETE, and MERGE.
	CategoryDML

	// CategoryDDL covers data-definition statements: CREATE / ALTER / DROP /
	// UNDROP for tables, views, materialized views, databases, and schemas.
	CategoryDDL

	// CategoryShow covers SHOW <objects> statements. This is a forward-compatible
	// placeholder; the parser does not yet emit a dedicated SHOW node.
	CategoryShow

	// CategoryDescribe covers DESCRIBE / DESC statements. This is a
	// forward-compatible placeholder; the parser does not yet emit a dedicated
	// DESCRIBE node.
	CategoryDescribe

	// CategoryOther covers recognized SQL statements that are none of the above:
	// USE, SET, CALL, EXPLAIN, GRANT, REVOKE, etc. This is a forward-compatible
	// placeholder; the parser does not yet emit dedicated nodes for these.
	CategoryOther
)

// String returns a human-readable label for the category.
func (c Category) String() string {
	switch c {
	case CategorySelect:
		return "SELECT"
	case CategoryDML:
		return "DML"
	case CategoryDDL:
		return "DDL"
	case CategoryShow:
		return "SHOW"
	case CategoryDescribe:
		return "DESCRIBE"
	case CategoryOther:
		return "Other"
	default:
		return "Unknown"
	}
}

// Classify returns the Category for a single statement node.
//
// It uses a type switch over the concrete ast.Node types defined in
// snowflake/ast/parsenodes.go. The mapping is:
//
//   - *ast.SelectStmt, *ast.SetOperationStmt          → CategorySelect
//   - *ast.InsertStmt, *ast.InsertMultiStmt,
//     *ast.UpdateStmt, *ast.DeleteStmt, *ast.MergeStmt → CategoryDML
//   - *ast.CreateTableStmt, *ast.AlterTableStmt,
//     *ast.CreateDatabaseStmt, *ast.AlterDatabaseStmt,
//     *ast.DropDatabaseStmt, *ast.UndropDatabaseStmt,
//     *ast.CreateSchemaStmt, *ast.AlterSchemaStmt,
//     *ast.DropSchemaStmt, *ast.UndropSchemaStmt,
//     *ast.CreateViewStmt, *ast.CreateMaterializedViewStmt,
//     *ast.AlterViewStmt, *ast.AlterMaterializedViewStmt,
//     *ast.DropStmt, *ast.UndropStmt                  → CategoryDDL
//   - nil or unrecognized node                         → CategoryUnknown
func Classify(node ast.Node) Category {
	if node == nil {
		return CategoryUnknown
	}
	switch node.(type) {
	// SELECT / set operations
	case *ast.SelectStmt, *ast.SetOperationStmt:
		return CategorySelect

	// DML
	case *ast.InsertStmt, *ast.InsertMultiStmt,
		*ast.UpdateStmt, *ast.DeleteStmt,
		*ast.MergeStmt:
		return CategoryDML

	// DDL — tables
	case *ast.CreateTableStmt, *ast.AlterTableStmt:
		return CategoryDDL

	// DDL — databases
	case *ast.CreateDatabaseStmt, *ast.AlterDatabaseStmt,
		*ast.DropDatabaseStmt, *ast.UndropDatabaseStmt:
		return CategoryDDL

	// DDL — schemas
	case *ast.CreateSchemaStmt, *ast.AlterSchemaStmt,
		*ast.DropSchemaStmt, *ast.UndropSchemaStmt:
		return CategoryDDL

	// DDL — views
	case *ast.CreateViewStmt, *ast.CreateMaterializedViewStmt,
		*ast.AlterViewStmt, *ast.AlterMaterializedViewStmt:
		return CategoryDDL

	// DDL — generic DROP / UNDROP (TABLE, TAG, DYNAMIC TABLE, etc.)
	case *ast.DropStmt, *ast.UndropStmt:
		return CategoryDDL

	default:
		return CategoryUnknown
	}
}

// ClassifySQL parses the first statement in sql and returns its Category.
//
// Parse errors are silently ignored; only the first successfully-parsed
// statement node is classified. Returns CategoryUnknown if the input is
// empty, consists only of whitespace or comments, or if parsing yields no
// nodes.
//
// Callers that need strict parse-error reporting should use the
// snowflake/diagnostics package instead.
func ClassifySQL(sql string) Category {
	result := parser.ParseBestEffort(sql)
	if result.File == nil || len(result.File.Stmts) == 0 {
		return CategoryUnknown
	}
	return Classify(result.File.Stmts[0])
}
