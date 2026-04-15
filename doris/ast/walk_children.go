package ast

// walkChildren walks the child nodes of node, calling Walk(v, child)
// for each child. Maintained manually until the node count warrants a
// code generator.
func walkChildren(v Visitor, node Node) {
	switch n := node.(type) {
	case *File:
		walkNodes(v, n.Stmts)
	case *ObjectName:
		// leaf node, no children
	case *TypeName:
		if n.ElementType != nil {
			Walk(v, n.ElementType)
		}
		if n.ValueType != nil {
			Walk(v, n.ValueType)
		}
		for _, f := range n.Fields {
			if f.Type != nil {
				Walk(v, f.Type)
			}
		}
	case *CreateIndexStmt:
		Walk(v, n.Name)
		Walk(v, n.Table)
	case *DropIndexStmt:
		Walk(v, n.Name)
		Walk(v, n.Table)
	case *BuildIndexStmt:
		Walk(v, n.Name)
		Walk(v, n.Table)
	case *CreateDatabaseStmt:
		Walk(v, n.Name)
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *AlterDatabaseStmt:
		Walk(v, n.Name)
		if n.NewName != nil {
			Walk(v, n.NewName)
		}
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *DropDatabaseStmt:
		Walk(v, n.Name)
	case *Property:
		// leaf node, no children

	// Expression nodes (T1.3).
	case *BinaryExpr:
		Walk(v, n.Left)
		Walk(v, n.Right)
	case *UnaryExpr:
		Walk(v, n.Expr)
	case *IsExpr:
		Walk(v, n.Expr)
	case *BetweenExpr:
		Walk(v, n.Expr)
		Walk(v, n.Low)
		Walk(v, n.High)
	case *InExpr:
		Walk(v, n.Expr)
		walkNodes(v, n.Values)
	case *LikeExpr:
		Walk(v, n.Expr)
		Walk(v, n.Pattern)
		if n.Escape != nil {
			Walk(v, n.Escape)
		}
	case *RegexpExpr:
		Walk(v, n.Expr)
		Walk(v, n.Pattern)
	case *FuncCallExpr:
		Walk(v, n.Name)
		walkNodes(v, n.Args)
		for _, item := range n.OrderBy {
			Walk(v, item)
		}
	case *CastExpr:
		Walk(v, n.Expr)
		Walk(v, n.TypeName)
	case *CaseExpr:
		if n.Operand != nil {
			Walk(v, n.Operand)
		}
		for _, w := range n.Whens {
			Walk(v, w)
		}
		if n.Else != nil {
			Walk(v, n.Else)
		}
	case *WhenClause:
		Walk(v, n.Cond)
		Walk(v, n.Result)
	case *SubqueryExpr:
		// placeholder leaf — no parsed children yet
	case *ColumnRef:
		Walk(v, n.Name)
	case *Literal:
		// leaf node, no children
	case *ParenExpr:
		Walk(v, n.Expr)
	case *ExistsExpr:
		Walk(v, n.Subquery)
	case *IntervalExpr:
		Walk(v, n.Value)
	case *OrderByItem:
		Walk(v, n.Expr)

	// CTE / WITH clause nodes (T1.6).
	case *WithClause:
		for _, cte := range n.CTEs {
			Walk(v, cte)
		}
	case *CTE:
		if n.Query != nil {
			Walk(v, n.Query)
		}

	// SELECT statement nodes (T1.4).
	case *SelectStmt:
		if n.With != nil {
			Walk(v, n.With)
		}
		for _, item := range n.Items {
			Walk(v, item)
		}
		walkNodes(v, n.From)
		if n.Where != nil {
			Walk(v, n.Where)
		}
		walkNodes(v, n.GroupBy)
		if n.Having != nil {
			Walk(v, n.Having)
		}
		if n.Qualify != nil {
			Walk(v, n.Qualify)
		}
		for _, item := range n.OrderBy {
			Walk(v, item)
		}
		if n.Limit != nil {
			Walk(v, n.Limit)
		}
		if n.Offset != nil {
			Walk(v, n.Offset)
		}
	case *SelectItem:
		if n.Expr != nil {
			Walk(v, n.Expr)
		}
		if n.TableName != nil {
			Walk(v, n.TableName)
		}
	case *TableRef:
		if n.Name != nil {
			Walk(v, n.Name)
		}
	case *JoinClause:
		Walk(v, n.Left)
		Walk(v, n.Right)
		if n.On != nil {
			Walk(v, n.On)
		}
	case *SetOpStmt:
		Walk(v, n.Left)
		Walk(v, n.Right)

	// DDL — CREATE TABLE nodes (T2.1).
	case *CreateTableStmt:
		Walk(v, n.Name)
		for _, col := range n.Columns {
			Walk(v, col)
		}
		for _, idx := range n.Indexes {
			Walk(v, idx)
		}
		for _, c := range n.Constraints {
			Walk(v, c)
		}
		if n.KeyDesc != nil {
			Walk(v, n.KeyDesc)
		}
		if n.PartitionBy != nil {
			Walk(v, n.PartitionBy)
		}
		if n.DistributedBy != nil {
			Walk(v, n.DistributedBy)
		}
		for _, r := range n.Rollup {
			Walk(v, r)
		}
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
		if n.Like != nil {
			Walk(v, n.Like)
		}
		if n.AsSelect != nil {
			Walk(v, n.AsSelect)
		}
	case *ColumnDef:
		if n.Type != nil {
			Walk(v, n.Type)
		}
		if n.Default != nil {
			Walk(v, n.Default)
		}
		if n.Generated != nil {
			Walk(v, n.Generated)
		}
	case *IndexDef:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *TableConstraint:
		// leaf-ish node, no Node children
	case *KeyDesc:
		// leaf-ish node, no Node children
	case *PartitionDesc:
		for _, p := range n.Partitions {
			Walk(v, p)
		}
	case *PartitionItem:
		// leaf-ish node, values stored as strings
	case *DistributionDesc:
		// leaf-ish node, no Node children
	case *RollupDef:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *RawQuery:
		// leaf node, no parsed children

	// DDL — ALTER TABLE nodes (T2.2).
	case *AlterTableStmt:
		Walk(v, n.Name)
		for _, action := range n.Actions {
			Walk(v, action)
		}
	case *AlterTableAction:
		if n.Column != nil {
			Walk(v, n.Column)
		}
		if n.NewTableName != nil {
			Walk(v, n.NewTableName)
		}
		if n.Partition != nil {
			Walk(v, n.Partition)
		}
		if n.PartitionDist != nil {
			Walk(v, n.PartitionDist)
		}
		for _, prop := range n.PartitionProps {
			Walk(v, prop)
		}
		if n.Rollup != nil {
			Walk(v, n.Rollup)
		}
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
		if n.Distribution != nil {
			Walk(v, n.Distribution)
		}

	// DDL — VIEW nodes (T2.4).
	case *ViewColumn:
		// leaf-ish node, no Node children
	case *CreateViewStmt:
		Walk(v, n.Name)
		for _, col := range n.Columns {
			Walk(v, col)
		}
		if n.Query != nil {
			Walk(v, n.Query)
		}
	case *AlterViewStmt:
		Walk(v, n.Name)
		for _, col := range n.Columns {
			Walk(v, col)
		}
		if n.Query != nil {
			Walk(v, n.Query)
		}
	case *DropViewStmt:
		Walk(v, n.Name)

	// DML — INSERT statement (T4.1).
	case *InsertStmt:
		if n.Target != nil {
			Walk(v, n.Target)
		}
		for _, row := range n.Values {
			walkNodes(v, row)
		}
		if n.Query != nil {
			Walk(v, n.Query)
		}

	// DML — UPDATE / DELETE nodes (T4.2).
	case *Assignment:
		Walk(v, n.Column)
		Walk(v, n.Value)
	case *UpdateStmt:
		Walk(v, n.Target)
		for _, a := range n.Assignments {
			Walk(v, a)
		}
		walkNodes(v, n.From)
		if n.Where != nil {
			Walk(v, n.Where)
		}
	case *DeleteStmt:
		Walk(v, n.Target)
		walkNodes(v, n.Using)
		if n.Where != nil {
			Walk(v, n.Where)
		}

	// DML — MERGE INTO nodes (T4.3).
	case *MergeStmt:
		if n.Target != nil {
			Walk(v, n.Target)
		}
		if n.Source != nil {
			Walk(v, n.Source)
		}
		if n.On != nil {
			Walk(v, n.On)
		}
		for _, clause := range n.Clauses {
			Walk(v, clause)
		}
	case *MergeClause:
		if n.And != nil {
			Walk(v, n.And)
		}
		for _, a := range n.Assignments {
			if a.Value != nil {
				Walk(v, a.Value)
			}
		}
		for _, val := range n.Values {
			if val != nil {
				Walk(v, val)
			}
		}

	// DDL — STORAGE VAULT nodes (T5.3).
	case *CreateStorageVaultStmt:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *AlterStorageVaultStmt:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *DropStorageVaultStmt:
		// leaf-ish node, no Node children
	case *SetDefaultStorageVaultStmt:
		// leaf-ish node, no Node children
	case *UnsetDefaultStorageVaultStmt:
		// leaf-ish node, no Node children

	// DDL — STORAGE POLICY nodes (T5.3).
	case *CreateStoragePolicyStmt:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *AlterStoragePolicyStmt:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *DropStoragePolicyStmt:
		// leaf-ish node, no Node children

	// DDL — REPOSITORY nodes (T5.3).
	case *CreateRepositoryStmt:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *AlterRepositoryStmt:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *DropRepositoryStmt:
		// leaf-ish node, no Node children

	// DDL — STAGE nodes (T5.3).
	case *CreateStageStmt:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *DropStageStmt:
		// leaf-ish node, no Node children

	// DDL — FILE nodes (T5.3).
	case *CreateFileStmt:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *DropFileStmt:
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	}
}
