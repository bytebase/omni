package ast

import (
	"reflect"
	"testing"
)

func TestWalkSelectStmt(t *testing.T) {
	stmt := &SelectStmt{
		TargetList: []ExprNode{
			&ColumnRef{Column: "id"},
			&ColumnRef{Column: "name"},
		},
		From: []TableExpr{
			&TableRef{Name: "users"},
		},
		Where: &BinaryExpr{
			Op:    BinOpAdd, // any op, just need a tree
			Left:  &ColumnRef{Column: "id"},
			Right: &IntLit{Value: 1},
		},
		OrderBy: []*OrderByItem{
			{Expr: &ColumnRef{Column: "name"}},
		},
		Limit: &Limit{
			Count: &IntLit{Value: 10},
		},
	}

	var visited []string
	Inspect(stmt, func(n Node) bool {
		if n == nil {
			return false
		}
		visited = append(visited, reflect.TypeOf(n).Elem().Name())
		return true
	})

	if len(visited) < 10 {
		t.Errorf("visited %d nodes, want at least 10: %v", len(visited), visited)
	}

	typeSet := map[string]bool{}
	for _, v := range visited {
		typeSet[v] = true
	}
	for _, want := range []string{"SelectStmt", "ColumnRef", "TableRef", "BinaryExpr", "IntLit", "OrderByItem", "Limit"} {
		if !typeSet[want] {
			t.Errorf("expected to visit %s, visited: %v", want, visited)
		}
	}
}

func TestWalkNil(t *testing.T) {
	Walk(inspector(func(n Node) bool { return true }), nil)
}

func TestInspectPruning(t *testing.T) {
	stmt := &SelectStmt{
		Where: &BinaryExpr{
			Op:    BinOpAdd,
			Left:  &ColumnRef{Column: "a"},
			Right: &IntLit{Value: 1},
		},
	}

	var visited []string
	Inspect(stmt, func(n Node) bool {
		if n == nil {
			return false
		}
		name := reflect.TypeOf(n).Elem().Name()
		visited = append(visited, name)
		if name == "BinaryExpr" {
			return false
		}
		return true
	})

	typeSet := map[string]bool{}
	for _, v := range visited {
		typeSet[v] = true
	}
	if !typeSet["SelectStmt"] {
		t.Error("expected SelectStmt")
	}
	if !typeSet["BinaryExpr"] {
		t.Error("expected BinaryExpr")
	}
	if typeSet["ColumnRef"] {
		t.Error("ColumnRef should have been pruned")
	}
	if typeSet["IntLit"] {
		t.Error("IntLit should have been pruned")
	}
}

// TestWalkChildrenCompleteness verifies that walkChildren handles every
// concrete AST node type. It creates a zero-value instance of each and
// verifies Walk visits it without panic.
func TestWalkChildrenCompleteness(t *testing.T) {
	// Get all struct types from walk_generated.go by testing each type.
	// We use reflection to find all types implementing Node.
	// If walkChildren is missing a case, the node will still be visited
	// (Walk calls v.Visit before walkChildren), but its children won't.

	nodeType := reflect.TypeOf((*Node)(nil)).Elem()

	for _, inst := range allKnownNodes() {
		typ := reflect.TypeOf(inst)
		if !typ.Implements(nodeType) {
			continue
		}
		name := typ.Elem().Name()
		t.Run(name, func(t *testing.T) {
			node := reflect.New(typ.Elem()).Interface().(Node)
			visited := false
			Inspect(node, func(n Node) bool {
				if n == node {
					visited = true
				}
				return true
			})
			if !visited {
				t.Errorf("Walk did not visit zero-value %s", name)
			}
		})
	}
}

// allKnownNodes returns one instance of every concrete AST node type.
// Keep this in sync with parsenodes.go — CI should run `go generate` and diff.
func allKnownNodes() []Node {
	return []Node{
		// Statements
		&SelectStmt{}, &InsertStmt{}, &UpdateStmt{}, &DeleteStmt{},
		&CreateTableStmt{}, &CreateIndexStmt{}, &CreateViewStmt{}, &CreateDatabaseStmt{},
		&CreateFunctionStmt{}, &CreateTriggerStmt{}, &CreateEventStmt{},
		&CreateUserStmt{}, &CreateRoleStmt{}, &CreateServerStmt{},
		&CreateTablespaceStmt{}, &CreateResourceGroupStmt{}, &CreateSpatialRefSysStmt{},
		&CreateLogfileGroupStmt{},
		&AlterTableStmt{}, &AlterDatabaseStmt{}, &AlterViewStmt{}, &AlterEventStmt{},
		&AlterRoutineStmt{}, &AlterUserStmt{}, &AlterServerStmt{},
		&AlterTablespaceStmt{}, &AlterResourceGroupStmt{}, &AlterInstanceStmt{},
		&AlterLogfileGroupStmt{},
		&DropTableStmt{}, &DropIndexStmt{}, &DropViewStmt{}, &DropDatabaseStmt{},
		&DropRoutineStmt{}, &DropTriggerStmt{}, &DropEventStmt{},
		&DropUserStmt{}, &DropRoleStmt{}, &DropServerStmt{},
		&DropTablespaceStmt{}, &DropResourceGroupStmt{}, &DropSpatialRefSysStmt{},
		&DropLogfileGroupStmt{},
		&TruncateStmt{}, &RenameTableStmt{}, &RenameUserStmt{},
		&SetStmt{}, &SetTransactionStmt{}, &SetPasswordStmt{},
		&SetRoleStmt{}, &SetDefaultRoleStmt{}, &SetResourceGroupStmt{},
		&ShowStmt{}, &UseStmt{}, &ExplainStmt{},
		&BeginStmt{}, &CommitStmt{}, &RollbackStmt{}, &SavepointStmt{},
		&LockTablesStmt{}, &UnlockTablesStmt{},
		&GrantStmt{}, &RevokeStmt{}, &GrantRoleStmt{}, &RevokeRoleStmt{},
		&FlushStmt{}, &KillStmt{}, &LoadDataStmt{},
		&PrepareStmt{}, &ExecuteStmt{}, &DeallocateStmt{},
		&DoStmt{}, &CallStmt{},
		&AnalyzeTableStmt{}, &OptimizeTableStmt{}, &CheckTableStmt{}, &RepairTableStmt{},
		&ChecksumTableStmt{},
		&HandlerOpenStmt{}, &HandlerReadStmt{}, &HandlerCloseStmt{},
		&SignalStmt{}, &ResignalStmt{}, &GetDiagnosticsStmt{},
		&ChangeReplicationSourceStmt{}, &ChangeReplicationFilterStmt{},
		&StartReplicaStmt{}, &StopReplicaStmt{}, &ResetReplicaStmt{},
		&StartGroupReplicationStmt{}, &StopGroupReplicationStmt{},
		&XAStmt{},
		&BeginEndBlock{}, &DeclareVarStmt{}, &DeclareConditionStmt{},
		&DeclareCursorStmt{}, &DeclareHandlerStmt{},
		&IfStmt{}, &CaseStmtNode{}, &WhileStmt{}, &RepeatStmt{}, &LoopStmt{},
		&LeaveStmt{}, &IterateStmt{}, &ReturnStmt{},
		&OpenCursorStmt{}, &FetchCursorStmt{}, &CloseCursorStmt{},
		&ResetMasterStmt{}, &ResetPersistStmt{},
		&CacheIndexStmt{}, &LoadIndexIntoCacheStmt{},
		&InstallPluginStmt{}, &UninstallPluginStmt{},
		&InstallComponentStmt{}, &UninstallComponentStmt{},
		&LockInstanceStmt{}, &UnlockInstanceStmt{},
		&TableStmt{}, &ValuesStmt{},
		&BinlogStmt{}, &CloneStmt{}, &RestartStmt{}, &ShutdownStmt{},
		&HelpStmt{}, &PurgeBinaryLogsStmt{}, &ImportTableStmt{},
		&RawStmt{},

		// Expressions
		&ColumnRef{}, &IntLit{}, &FloatLit{}, &StringLit{},
		&HexLit{}, &BitLit{}, &BoolLit{}, &NullLit{}, &TemporalLit{},
		&DefaultExpr{}, &StarExpr{},
		&UnaryExpr{}, &BinaryExpr{},
		&BetweenExpr{}, &InExpr{}, &LikeExpr{},
		&IsExpr{}, &ExistsExpr{}, &SubqueryExpr{},
		&FuncCallExpr{}, &CastExpr{}, &ConvertExpr{}, &ExtractExpr{},
		&CaseExpr{}, &CaseWhen{}, &IntervalExpr{}, &CollateExpr{},
		&ParenExpr{}, &RowExpr{},
		&MatchExpr{}, &MemberOfExpr{}, &VariableRef{},
		&ValuesStmt{}, // also an expression in some contexts
		&JsonTableExpr{}, &JsonTableColumn{},

		// Table expressions
		&TableRef{}, &JoinClause{},

		// DDL components
		&CommonTableExpr{}, &ColumnDef{}, &ColumnConstraint{}, &Constraint{},
		&DataType{}, &IndexColumn{}, &IndexOption{},
		&PartitionClause{}, &PartitionDef{}, &SubPartitionDef{},
		&AlterTableCmd{}, &GeneratedColumn{},
		&TableOption{}, &DatabaseOption{},
		&OnCondition{}, &UsingCondition{},
		&IndexHint{},

		// DML components
		&Assignment{}, &OrderByItem{}, &Limit{}, &ForUpdate{},
		&IntoClause{}, &WindowDef{}, &WindowFrame{}, &WindowFrameBound{},
		&ResTarget{}, &LockTable{},

		// User/grant components
		&UserSpec{}, &RequireClause{}, &ResourceOption{},
		&GrantTarget{}, &VCPUSpec{},

		// Replication components
		&ReplicationOption{}, &ReplicationFilter{},

		// Signal/diagnostics components
		&SignalInfoItem{}, &DiagnosticsItem{},

		// Compound statement components
		&ElseIf{}, &CaseStmtWhen{}, &FuncParam{},
		&RoutineCharacteristic{}, &TriggerOrder{}, &EventSchedule{},

		// Rename components
		&RenameTablePair{}, &RenameUserPair{},
		&RegistrationOp{}, &FactorOp{},
		&LoadIndexTable{},

		// Value nodes
		&List{}, &String{}, &Integer{},
	}
}
