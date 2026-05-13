package ast

import (
	"reflect"
	"testing"
)

func TestWalkSelectStmt(t *testing.T) {
	stmt := &SelectStmt{
		TargetList: &List{Items: []Node{
			&ColumnRef{Column: "id"},
			&ColumnRef{Column: "name"},
		}},
		FromClause: &List{Items: []Node{
			&TableRef{Object: "users"},
		}},
		WhereClause: &BinaryExpr{
			Op:    BinOpAdd, // any op, just need a tree
			Left:  &ColumnRef{Column: "id"},
			Right: &Literal{Str: "1"},
		},
		OrderByClause: &List{Items: []Node{
			&OrderByItem{Expr: &ColumnRef{Column: "name"}},
		}},
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
	for _, want := range []string{"SelectStmt", "ColumnRef", "TableRef", "BinaryExpr", "Literal", "OrderByItem", "List"} {
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
		WhereClause: &BinaryExpr{
			Op:    BinOpAdd,
			Left:  &ColumnRef{Column: "a"},
			Right: &Literal{Str: "1"},
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
	if typeSet["Literal"] {
		t.Error("Literal should have been pruned")
	}
}

// TestWalkChildrenCompleteness verifies that walkChildren handles every
// concrete AST node type. It creates a zero-value instance of each and
// verifies Walk visits it without panic.
func TestWalkChildrenCompleteness(t *testing.T) {
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
		// Statements — DML
		&SelectStmt{}, &InsertStmt{}, &UpdateStmt{}, &DeleteStmt{}, &MergeStmt{},

		// Statements — DDL: CREATE
		&CreateTableStmt{}, &CreateIndexStmt{}, &CreateViewStmt{}, &CreateDatabaseStmt{},
		&CreateSchemaStmt{}, &CreateFunctionStmt{}, &CreateProcedureStmt{},
		&CreateTriggerStmt{}, &CreateTypeStmt{}, &CreateSequenceStmt{},
		&CreateSynonymStmt{}, &CreateStatisticsStmt{}, &CreateAssemblyStmt{},
		&CreatePartitionFunctionStmt{}, &CreatePartitionSchemeStmt{},
		&CreateFulltextIndexStmt{}, &CreateFulltextCatalogStmt{},
		&CreateFulltextStoplistStmt{}, &CreateSearchPropertyListStmt{},
		&CreateXmlSchemaCollectionStmt{}, &CreateXmlIndexStmt{},
		&CreateSelectiveXmlIndexStmt{}, &CreateSpatialIndexStmt{},
		&CreateAggregateStmt{}, &CreateJsonIndexStmt{}, &CreateVectorIndexStmt{},
		&CreateMaterializedViewStmt{}, &CreateFederationStmt{},
		&CreateExternalTableAsSelectStmt{}, &CreateTableCloneStmt{},
		&CreateTableAsSelectStmt{}, &CreateRemoteTableAsSelectStmt{},

		// Statements — DDL: ALTER
		&AlterTableStmt{}, &AlterDatabaseStmt{}, &AlterIndexStmt{},
		&AlterSchemaStmt{}, &AlterSequenceStmt{},
		&AlterPartitionFunctionStmt{}, &AlterPartitionSchemeStmt{},
		&AlterFulltextIndexStmt{}, &AlterFulltextCatalogStmt{},
		&AlterFulltextStoplistStmt{}, &AlterSearchPropertyListStmt{},
		&AlterXmlSchemaCollectionStmt{}, &AlterAssemblyStmt{},
		&AlterServerConfigurationStmt{}, &AlterMaterializedViewStmt{},
		&AlterFederationStmt{},

		// Statements — DDL: DROP
		&DropStmt{}, &DropStatisticsStmt{}, &DropAggregateStmt{},
		&DropFulltextStoplistStmt{}, &DropSearchPropertyListStmt{},
		&DropFederationStmt{},

		// Statements — DML components
		&TruncateStmt{}, &BulkInsertStmt{}, &InsertBulkStmt{},

		// Statements — Control flow
		&DeclareStmt{}, &SetStmt{}, &IfStmt{}, &WhileStmt{},
		&BeginEndStmt{}, &TryCatchStmt{}, &ReturnStmt{},
		&BreakStmt{}, &ContinueStmt{}, &GotoStmt{}, &LabelStmt{},
		&WaitForStmt{},

		// Statements — Execution
		&ExecStmt{}, &PrintStmt{}, &RaiseErrorStmt{}, &ThrowStmt{},

		// Statements — Transaction
		&BeginTransStmt{}, &CommitTransStmt{}, &RollbackTransStmt{},
		&SaveTransStmt{}, &BeginDistributedTransStmt{},

		// Statements — Cursor
		&DeclareCursorStmt{}, &OpenCursorStmt{}, &FetchCursorStmt{},
		&CloseCursorStmt{}, &DeallocateCursorStmt{},

		// Statements — Database
		&UseStmt{}, &GoStmt{},

		// Statements — Security
		&SecurityStmt{}, &GrantStmt{}, &SecurityKeyStmt{},
		&SecurityPolicyStmt{}, &SensitivityClassificationStmt{}, &SignatureStmt{},

		// Statements — Server & Admin
		&DbccStmt{}, &BackupStmt{}, &RestoreStmt{},
		&CheckpointStmt{}, &ReconfigureStmt{}, &ShutdownStmt{},
		&KillStmt{}, &KillStatsJobStmt{}, &KillQueryNotificationStmt{},
		&SetOptionStmt{}, &EnableDisableTriggerStmt{},
		&UpdateStatisticsStmt{}, &RenameStmt{}, &LinenoStmt{},

		// Statements — Text operations
		&ReadtextStmt{}, &WritetextStmt{}, &UpdatetextStmt{},

		// Statements — Service Broker
		&ServiceBrokerStmt{}, &ReceiveStmt{},

		// Statements — PREDICT / COPY INTO
		&PredictStmt{}, &CopyIntoStmt{},

		// Statements — USE FEDERATION
		&UseFederationStmt{},

		// Expressions
		&BinaryExpr{}, &UnaryExpr{}, &FuncCallExpr{},
		&CaseExpr{}, &CaseWhen{}, &BetweenExpr{}, &InExpr{},
		&LikeExpr{}, &IsExpr{}, &ExistsExpr{},
		&CastExpr{}, &ConvertExpr{}, &TryCastExpr{}, &TryConvertExpr{},
		&CoalesceExpr{}, &NullifExpr{}, &IifExpr{},
		&DatePart{}, &NextValueForExpr{}, &ParseExpr{}, &JsonKeyValueExpr{},
		&ColumnRef{}, &VariableRef{}, &StarExpr{}, &Literal{},
		&SubqueryExpr{}, &SubqueryComparisonExpr{},
		&CollateExpr{}, &AtTimeZoneExpr{}, &ParenExpr{},
		&CurrentOfExpr{}, &MethodCallExpr{},
		&PivotExpr{}, &UnpivotExpr{},
		&GroupingSetsExpr{}, &RollupExpr{}, &CubeExpr{},
		&SelectAssign{},

		// Table expressions
		&TableRef{}, &JoinClause{}, &AliasedTableRef{},

		// DDL components
		&ColumnDef{}, &ConstraintDef{}, &InlineIndexDef{},
		&IndexColumn{}, &DataType{}, &TableOption{}, &DatabaseOption{},
		&ComputedColumnDef{}, &EncryptedWithSpec{}, &GeneratedAlwaysSpec{},
		&NullableSpec{}, &IdentitySpec{}, &EdgeConnectionDef{},
		&ParamDef{}, &ReturnsTableDef{}, &TableTypeIndex{},
		&TriggerEvent{}, &TriggerOption{},
		&DatabaseFileSpec{}, &DatabaseFilegroup{}, &SizeValue{},
		&CommonTableExpr{}, &WithClause{}, &XmlNamespaceDecl{},

		// DML components
		&TopClause{}, &FetchClause{}, &ForClause{},
		&OutputClause{}, &MergeWhenClause{},
		&MergeUpdateAction{}, &MergeDeleteAction{}, &MergeInsertAction{},
		&ResTarget{}, &OrderByItem{}, &SetExpr{}, &ValuesClause{},
		&WindowDef{}, &OverClause{}, &WindowFrame{}, &WindowBound{},
		&AlterTableAction{}, &AlterColumnOption{},
		&ExecArg{}, &VariableDecl{},
		&TableSampleClause{}, &TableHint{}, &QueryHint{}, &OptimizeForParam{},
		&InsertBulkColumnDef{}, &CopyIntoColumn{},

		// Security components
		&SecurityPrincipalOption{}, &AuditSpecAction{},
		&EndpointOption{}, &AvailabilityGroupOption{},
		&EventNotificationOption{}, &ResourceGovernorOption{},
		&ExternalOption{}, &SecurityPredicate{}, &CryptoItem{},
		&SensitivityOption{},

		// Backup/restore components
		&BackupRestoreOption{}, &DbccOption{},

		// Service broker components
		&ServiceBrokerOption{}, &ReceiveColumn{},

		// Server config components
		&ServerConfigOption{},

		// Assembly components
		&AssemblyAction{}, &AssemblyFile{},

		// Value nodes
		&List{}, &String{}, &Integer{}, &Float{}, &Boolean{},
	}
}
