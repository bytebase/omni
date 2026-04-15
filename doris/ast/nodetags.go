package ast

// NodeTag identifies the concrete type of an AST node.
//
// Every concrete node type defined under doris/ast must declare a unique
// NodeTag constant in this file and return it from Tag(). This enables fast
// switch dispatch and code-generated walker support.
//
// The numeric values are NOT stable -- do not persist them. Tags are assigned
// by source order; reorder freely as the package evolves.
type NodeTag int

const (
	// T_Invalid is the zero-value tag, returned only by uninitialized nodes
	// or test stubs that have no need for a real tag.
	T_Invalid NodeTag = iota

	// T_File is the tag for *File, the root statement-list container
	// returned by the parser entry point.
	T_File

	// T_ObjectName is the tag for *ObjectName, a qualified multi-part name
	// (e.g., catalog.db.table).
	T_ObjectName

	// T_TypeName is the tag for *TypeName, a Doris data type as it appears in
	// SQL source (e.g., INT, VARCHAR(255), ARRAY<INT>, MAP<STRING,INT>).
	T_TypeName

	// T_CreateIndexStmt is the tag for *CreateIndexStmt.
	T_CreateIndexStmt

	// T_DropIndexStmt is the tag for *DropIndexStmt.
	T_DropIndexStmt

	// T_BuildIndexStmt is the tag for *BuildIndexStmt.
	T_BuildIndexStmt

	// T_CreateDatabaseStmt is the tag for *CreateDatabaseStmt.
	T_CreateDatabaseStmt

	// T_AlterDatabaseStmt is the tag for *AlterDatabaseStmt.
	T_AlterDatabaseStmt

	// T_DropDatabaseStmt is the tag for *DropDatabaseStmt.
	T_DropDatabaseStmt

	// T_Property is the tag for *Property, a key=value pair in PROPERTIES clauses.
	T_Property

	// Expression nodes (T1.3).

	// T_BinaryExpr is the tag for *BinaryExpr (left op right).
	T_BinaryExpr

	// T_UnaryExpr is the tag for *UnaryExpr (op expr).
	T_UnaryExpr

	// T_IsExpr is the tag for *IsExpr (expr IS [NOT] NULL/TRUE/FALSE).
	T_IsExpr

	// T_BetweenExpr is the tag for *BetweenExpr (expr [NOT] BETWEEN low AND high).
	T_BetweenExpr

	// T_InExpr is the tag for *InExpr (expr [NOT] IN (...)).
	T_InExpr

	// T_LikeExpr is the tag for *LikeExpr (expr [NOT] LIKE pattern [ESCAPE esc]).
	T_LikeExpr

	// T_RegexpExpr is the tag for *RegexpExpr (expr [NOT] REGEXP/RLIKE pattern).
	T_RegexpExpr

	// T_FuncCallExpr is the tag for *FuncCallExpr (name(args...)).
	T_FuncCallExpr

	// T_CastExpr is the tag for *CastExpr (CAST/TRY_CAST).
	T_CastExpr

	// T_CaseExpr is the tag for *CaseExpr (CASE...END).
	T_CaseExpr

	// T_WhenClause is the tag for *WhenClause (WHEN cond THEN result).
	T_WhenClause

	// T_SubqueryExpr is the tag for *SubqueryExpr (placeholder for subqueries).
	T_SubqueryExpr

	// T_ColumnRef is the tag for *ColumnRef (qualified column reference).
	T_ColumnRef

	// T_Literal is the tag for *Literal (numeric, string, boolean, NULL).
	T_Literal

	// T_ParenExpr is the tag for *ParenExpr ((expr)).
	T_ParenExpr

	// T_ExistsExpr is the tag for *ExistsExpr (EXISTS (subquery)).
	T_ExistsExpr

	// T_IntervalExpr is the tag for *IntervalExpr (INTERVAL expr unit).
	T_IntervalExpr

	// T_OrderByItem is the tag for *OrderByItem (expr ASC/DESC NULLS FIRST/LAST).
	T_OrderByItem

	// CTE / WITH clause nodes (T1.6).

	// T_WithClause is the tag for *WithClause.
	T_WithClause

	// T_CTE is the tag for *CTE (one named CTE entry inside a WITH clause).
	T_CTE

	// SELECT statement nodes (T1.4).

	// T_SelectStmt is the tag for *SelectStmt.
	T_SelectStmt

	// T_SelectItem is the tag for *SelectItem (one item in the SELECT list).
	T_SelectItem

	// T_TableRef is the tag for *TableRef (table reference in FROM clause).
	T_TableRef

	// T_JoinClause is the tag for *JoinClause (JOIN expression in FROM).
	T_JoinClause

	// T_SetOpStmt is the tag for *SetOpStmt (UNION/INTERSECT/EXCEPT/MINUS).
	T_SetOpStmt

	// DDL — CREATE TABLE nodes (T2.1).

	// T_CreateTableStmt is the tag for *CreateTableStmt.
	T_CreateTableStmt

	// T_ColumnDef is the tag for *ColumnDef.
	T_ColumnDef

	// T_IndexDef is the tag for *IndexDef (inline index in CREATE TABLE).
	T_IndexDef

	// T_TableConstraint is the tag for *TableConstraint (PRIMARY KEY, UNIQUE).
	T_TableConstraint

	// T_KeyDesc is the tag for *KeyDesc (AGGREGATE/UNIQUE/DUPLICATE KEY).
	T_KeyDesc

	// T_PartitionDesc is the tag for *PartitionDesc (PARTITION BY RANGE/LIST).
	T_PartitionDesc

	// T_PartitionItem is the tag for *PartitionItem.
	T_PartitionItem

	// T_DistributionDesc is the tag for *DistributionDesc (DISTRIBUTED BY HASH/RANDOM).
	T_DistributionDesc

	// T_RollupDef is the tag for *RollupDef.
	T_RollupDef

	// T_RawQuery is the tag for *RawQuery (placeholder for unparsed SELECT).
	T_RawQuery

	// DDL — ALTER TABLE nodes (T2.2).

	// T_AlterTableStmt is the tag for *AlterTableStmt.
	T_AlterTableStmt

	// T_AlterTableAction is the tag for *AlterTableAction.
	T_AlterTableAction

	// DDL — VIEW nodes (T2.4).

	// T_ViewColumn is the tag for *ViewColumn (one column in a view column list).
	T_ViewColumn

	// T_CreateViewStmt is the tag for *CreateViewStmt.
	T_CreateViewStmt

	// T_AlterViewStmt is the tag for *AlterViewStmt.
	T_AlterViewStmt

	// T_DropViewStmt is the tag for *DropViewStmt.
	T_DropViewStmt

	// DML — INSERT statement (T4.1).

	// T_InsertStmt is the tag for *InsertStmt.
	T_InsertStmt

	// DML — UPDATE / DELETE nodes (T4.2).

	// T_Assignment is the tag for *Assignment (col = expr in UPDATE SET clause).
	T_Assignment

	// T_UpdateStmt is the tag for *UpdateStmt.
	T_UpdateStmt

	// T_DeleteStmt is the tag for *DeleteStmt.
	T_DeleteStmt

	// DML — MERGE INTO nodes (T4.3).

	// T_MergeStmt is the tag for *MergeStmt.
	T_MergeStmt

	// T_MergeClause is the tag for *MergeClause (one WHEN clause inside MERGE).
	T_MergeClause
)

// String returns a human-readable representation of the tag.
func (t NodeTag) String() string {
	switch t {
	case T_Invalid:
		return "Invalid"
	case T_File:
		return "File"
	case T_ObjectName:
		return "ObjectName"
	case T_TypeName:
		return "TypeName"
	case T_CreateIndexStmt:
		return "CreateIndexStmt"
	case T_DropIndexStmt:
		return "DropIndexStmt"
	case T_BuildIndexStmt:
		return "BuildIndexStmt"
	case T_CreateDatabaseStmt:
		return "CreateDatabaseStmt"
	case T_AlterDatabaseStmt:
		return "AlterDatabaseStmt"
	case T_DropDatabaseStmt:
		return "DropDatabaseStmt"
	case T_Property:
		return "Property"
	case T_BinaryExpr:
		return "BinaryExpr"
	case T_UnaryExpr:
		return "UnaryExpr"
	case T_IsExpr:
		return "IsExpr"
	case T_BetweenExpr:
		return "BetweenExpr"
	case T_InExpr:
		return "InExpr"
	case T_LikeExpr:
		return "LikeExpr"
	case T_RegexpExpr:
		return "RegexpExpr"
	case T_FuncCallExpr:
		return "FuncCallExpr"
	case T_CastExpr:
		return "CastExpr"
	case T_CaseExpr:
		return "CaseExpr"
	case T_WhenClause:
		return "WhenClause"
	case T_SubqueryExpr:
		return "SubqueryExpr"
	case T_ColumnRef:
		return "ColumnRef"
	case T_Literal:
		return "Literal"
	case T_ParenExpr:
		return "ParenExpr"
	case T_ExistsExpr:
		return "ExistsExpr"
	case T_IntervalExpr:
		return "IntervalExpr"
	case T_OrderByItem:
		return "OrderByItem"
	case T_WithClause:
		return "WithClause"
	case T_CTE:
		return "CTE"
	case T_SelectStmt:
		return "SelectStmt"
	case T_SelectItem:
		return "SelectItem"
	case T_TableRef:
		return "TableRef"
	case T_JoinClause:
		return "JoinClause"
	case T_SetOpStmt:
		return "SetOpStmt"
	case T_CreateTableStmt:
		return "CreateTableStmt"
	case T_ColumnDef:
		return "ColumnDef"
	case T_IndexDef:
		return "IndexDef"
	case T_TableConstraint:
		return "TableConstraint"
	case T_KeyDesc:
		return "KeyDesc"
	case T_PartitionDesc:
		return "PartitionDesc"
	case T_PartitionItem:
		return "PartitionItem"
	case T_DistributionDesc:
		return "DistributionDesc"
	case T_RollupDef:
		return "RollupDef"
	case T_RawQuery:
		return "RawQuery"
	case T_AlterTableStmt:
		return "AlterTableStmt"
	case T_AlterTableAction:
		return "AlterTableAction"
	case T_ViewColumn:
		return "ViewColumn"
	case T_CreateViewStmt:
		return "CreateViewStmt"
	case T_AlterViewStmt:
		return "AlterViewStmt"
	case T_DropViewStmt:
		return "DropViewStmt"
	case T_InsertStmt:
		return "InsertStmt"
	case T_Assignment:
		return "Assignment"
	case T_UpdateStmt:
		return "UpdateStmt"
	case T_DeleteStmt:
		return "DeleteStmt"
	case T_MergeStmt:
		return "MergeStmt"
	case T_MergeClause:
		return "MergeClause"
	default:
		return "Unknown"
	}
}
