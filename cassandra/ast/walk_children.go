package ast

func walkChildren(v Visitor, node Node) {
	switch n := node.(type) {
	case *List:
		for _, item := range n.Items {
			Walk(v, item)
		}
	case *RawStmt:
		Walk(v, n.Stmt)

	// Identifiers / Names
	case *Identifier:
	case *QualifiedName:
		for _, p := range n.Parts {
			Walk(v, p)
		}

	// Literals
	case *StringLit:
	case *IntegerLit:
	case *FloatLit:
	case *BoolLit:
	case *NullLit:
	case *UUIDLit:
	case *HexLit:
	case *CodeBlock:
	case *StarExpr:

	// Collections
	case *MapLit:
		for _, k := range n.Keys {
			Walk(v, k)
		}
		for _, val := range n.Values {
			Walk(v, val)
		}
	case *SetLit:
		for _, e := range n.Elements {
			Walk(v, e)
		}
	case *ListLit:
		for _, e := range n.Elements {
			Walk(v, e)
		}
	case *TupleLit:
		for _, e := range n.Elements {
			Walk(v, e)
		}
	case *VectorLit:
		for _, e := range n.Elements {
			Walk(v, e)
		}

	// Expressions
	case *FunctionCall:
		Walk(v, n.Name)
		for _, a := range n.Args {
			Walk(v, a)
		}
	case *BinaryExpr:
		Walk(v, n.Left)
		Walk(v, n.Right)
	case *InExpr:
		Walk(v, n.Column)
		for _, val := range n.Values {
			Walk(v, val)
		}
	case *ContainsExpr:
		Walk(v, n.Column)
		Walk(v, n.Value)
	case *TupleCompareExpr:
		for _, c := range n.Columns {
			Walk(v, c)
		}
		for _, val := range n.Values {
			Walk(v, val)
		}
	case *TupleInExpr:
		for _, c := range n.Columns {
			Walk(v, c)
		}
		for _, t := range n.Tuples {
			Walk(v, t)
		}
	case *IndexAccess:
		Walk(v, n.Collection)
		Walk(v, n.Index)
	case *DotAccess:
		Walk(v, n.Object)
		Walk(v, n.Field)

	// Types
	case *DataType:
		Walk(v, n.Name)
		for _, p := range n.TypeParams {
			Walk(v, p)
		}
		if n.Dimension != nil {
			Walk(v, n.Dimension)
		}

	// Clauses / helpers
	case *ColumnDef:
		Walk(v, n.Name)
		Walk(v, n.Type)
	case *PrimaryKeyDef:
		for _, k := range n.PartitionKeys {
			Walk(v, k)
		}
		for _, k := range n.ClusteringKeys {
			Walk(v, k)
		}
	case *ClusteringOrder:
		Walk(v, n.Column)
	case *TableOption:
		Walk(v, n.Name)
		Walk(v, n.Value)
	case *OptionHash:
		for _, item := range n.Items {
			Walk(v, item)
		}
	case *OptionHashItem:
		Walk(v, n.Key)
		Walk(v, n.Value)
	case *SelectElement:
		Walk(v, n.Expr)
		if n.Alias != nil {
			Walk(v, n.Alias)
		}
	case *AssignmentElement:
		Walk(v, n.Target)
		Walk(v, n.Value)
	case *IfCondition:
		Walk(v, n.Column)
		Walk(v, n.Value)
	case *UsingClause:
		if n.TTL != nil {
			Walk(v, n.TTL)
		}
		if n.Timestamp != nil {
			Walk(v, n.Timestamp)
		}
	case *OrderByElement:
		Walk(v, n.Column)
		if n.AnnVector != nil {
			Walk(v, n.AnnVector)
		}
		if n.AnnLimit != nil {
			Walk(v, n.AnnLimit)
		}

	// DML
	case *SelectStmt:
		for _, e := range n.Elements {
			Walk(v, e)
		}
		if n.From != nil {
			Walk(v, n.From)
		}
		for _, w := range n.Where {
			Walk(v, w)
		}
		for _, o := range n.OrderBy {
			Walk(v, o)
		}
		if n.Limit != nil {
			Walk(v, n.Limit)
		}
	case *InsertStmt:
		Walk(v, n.Table)
		for _, c := range n.Columns {
			Walk(v, c)
		}
		for _, val := range n.Values {
			Walk(v, val)
		}
		if n.JSONValue != nil {
			Walk(v, n.JSONValue)
		}
		if n.Using != nil {
			Walk(v, n.Using)
		}
	case *UpdateStmt:
		Walk(v, n.Table)
		if n.Using != nil {
			Walk(v, n.Using)
		}
		for _, a := range n.Assignments {
			Walk(v, a)
		}
		for _, w := range n.Where {
			Walk(v, w)
		}
		for _, c := range n.IfConditions {
			Walk(v, c)
		}
	case *DeleteStmt:
		for _, c := range n.Columns {
			Walk(v, c)
		}
		Walk(v, n.From)
		if n.Using != nil {
			Walk(v, n.Using)
		}
		for _, w := range n.Where {
			Walk(v, w)
		}
		for _, c := range n.IfConditions {
			Walk(v, c)
		}
	case *BatchStmt:
		if n.Using != nil {
			Walk(v, n.Using)
		}
		for _, s := range n.Statements {
			Walk(v, s)
		}
	case *TruncateStmt:
		Walk(v, n.Table)
	case *UseStmt:
		Walk(v, n.Keyspace)

	// DDL
	case *CreateKeyspaceStmt:
		Walk(v, n.Name)
		if n.Replication != nil {
			Walk(v, n.Replication)
		}
		if n.DurableWrites != nil {
			Walk(v, n.DurableWrites)
		}
	case *AlterKeyspaceStmt:
		Walk(v, n.Name)
		if n.Replication != nil {
			Walk(v, n.Replication)
		}
		if n.DurableWrites != nil {
			Walk(v, n.DurableWrites)
		}
	case *DropKeyspaceStmt:
		Walk(v, n.Name)
	case *CreateTableStmt:
		Walk(v, n.Name)
		for _, c := range n.Columns {
			Walk(v, c)
		}
		if n.PrimaryKey != nil {
			Walk(v, n.PrimaryKey)
		}
		for _, o := range n.Options {
			Walk(v, o)
		}
		for _, co := range n.ClusteringOrders {
			Walk(v, co)
		}
	case *AlterTableStmt:
		Walk(v, n.Name)
		for _, c := range n.AddColumns {
			Walk(v, c)
		}
		for _, c := range n.DropColumns {
			Walk(v, c)
		}
		if n.RenameFrom != nil {
			Walk(v, n.RenameFrom)
		}
		if n.RenameTo != nil {
			Walk(v, n.RenameTo)
		}
		for _, o := range n.Options {
			Walk(v, o)
		}
	case *DropTableStmt:
		Walk(v, n.Name)
	case *CreateIndexStmt:
		if n.IndexName != nil {
			Walk(v, n.IndexName)
		}
		Walk(v, n.Table)
		Walk(v, n.Column)
		if n.UsingClass != nil {
			Walk(v, n.UsingClass)
		}
		if n.Options != nil {
			Walk(v, n.Options)
		}
	case *DropIndexStmt:
		Walk(v, n.Name)
	case *CreateTypeStmt:
		Walk(v, n.Name)
		for _, f := range n.Fields {
			Walk(v, f)
		}
	case *AlterTypeStmt:
		Walk(v, n.Name)
		if n.AlterColumn != nil {
			Walk(v, n.AlterColumn)
		}
		if n.AlterType != nil {
			Walk(v, n.AlterType)
		}
		for _, f := range n.AddFields {
			Walk(v, f)
		}
		for _, r := range n.Renames {
			Walk(v, r)
		}
	case *AlterTypeRenameItem:
		Walk(v, n.From)
		Walk(v, n.To)
	case *DropTypeStmt:
		Walk(v, n.Name)
	case *CreateMVStmt:
		Walk(v, n.Name)
		for _, c := range n.SelectColumns {
			Walk(v, c)
		}
		Walk(v, n.FromTable)
		for _, c := range n.WhereNotNull {
			Walk(v, c)
		}
		for _, r := range n.WhereRelations {
			Walk(v, r)
		}
		if n.PrimaryKey != nil {
			Walk(v, n.PrimaryKey)
		}
		for _, o := range n.Options {
			Walk(v, o)
		}
		for _, co := range n.ClusteringOrders {
			Walk(v, co)
		}
	case *AlterMVStmt:
		Walk(v, n.Name)
		for _, o := range n.Options {
			Walk(v, o)
		}
	case *DropMVStmt:
		Walk(v, n.Name)
	case *CreateFunctionStmt:
		Walk(v, n.Name)
		for _, p := range n.Params {
			Walk(v, p)
		}
		if n.ReturnType != nil {
			Walk(v, n.ReturnType)
		}
		if n.Language != nil {
			Walk(v, n.Language)
		}
		if n.Body != nil {
			Walk(v, n.Body)
		}
	case *FunctionParam:
		Walk(v, n.Name)
		Walk(v, n.Type)
	case *DropFunctionStmt:
		Walk(v, n.Name)
	case *CreateAggregateStmt:
		Walk(v, n.Name)
		if n.ParamType != nil {
			Walk(v, n.ParamType)
		}
		if n.SFunc != nil {
			Walk(v, n.SFunc)
		}
		if n.SType != nil {
			Walk(v, n.SType)
		}
		if n.FinalFunc != nil {
			Walk(v, n.FinalFunc)
		}
		if n.InitCond != nil {
			Walk(v, n.InitCond)
		}
	case *DropAggregateStmt:
		Walk(v, n.Name)
	case *CreateTriggerStmt:
		Walk(v, n.Name)
		if n.UsingClass != nil {
			Walk(v, n.UsingClass)
		}
	case *DropTriggerStmt:
		if n.Name != nil {
			Walk(v, n.Name)
		}
		if n.Table != nil {
			Walk(v, n.Table)
		}

	// Auth
	case *CreateRoleStmt:
		Walk(v, n.Name)
		for _, o := range n.Options {
			Walk(v, o)
		}
	case *RoleOption:
		if n.Value != nil {
			Walk(v, n.Value)
		}
	case *AlterRoleStmt:
		Walk(v, n.Name)
		for _, o := range n.Options {
			Walk(v, o)
		}
	case *DropRoleStmt:
		Walk(v, n.Name)
	case *CreateUserStmt:
		Walk(v, n.Name)
		if n.Password != nil {
			Walk(v, n.Password)
		}
	case *AlterUserStmt:
		Walk(v, n.Name)
		if n.Password != nil {
			Walk(v, n.Password)
		}
	case *DropUserStmt:
		Walk(v, n.Name)
	case *GrantStmt:
		if n.Resource != nil {
			Walk(v, n.Resource)
		}
		Walk(v, n.Role)
	case *RevokeStmt:
		if n.Resource != nil {
			Walk(v, n.Resource)
		}
		Walk(v, n.Role)
	case *Resource:
		if n.Name != nil {
			Walk(v, n.Name)
		}
	case *ListPermissionsStmt:
		if n.Resource != nil {
			Walk(v, n.Resource)
		}
		if n.Role != nil {
			Walk(v, n.Role)
		}
	case *ListRolesStmt:
		if n.Of != nil {
			Walk(v, n.Of)
		}
	}
}
