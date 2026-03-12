package ddl

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// --- Phase 3 validation improvements and error path tests ---

func TestSetDefaultOnIdentityColumn(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	// Manually mark column as identity (ALWAYS).
	rel := c.GetRelation("", "t")
	rel.Columns[0].Identity = 'a'

	err := c.AlterTableStmt(makeAlterTableStmt("", "t",
		makeATSetDefault("id", "42"),
	))
	assertErrorCode(t, err, CodeInvalidObjectDefinition)
}

func TestSetDefaultOnGeneratedColumn(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "val", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	// Manually mark column as generated (stored).
	rel := c.GetRelation("", "t")
	rel.Columns[0].Generated = 's'

	err := c.AlterTableStmt(makeAlterTableStmt("", "t",
		makeATSetDefault("val", "1"),
	))
	assertErrorCode(t, err, CodeInvalidObjectDefinition)
}

func TestDropDefaultOnIdentityColumn(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	// Manually mark column as identity (ALWAYS).
	rel := c.GetRelation("", "t")
	rel.Columns[0].Identity = 'a'

	err := c.AlterTableStmt(makeAlterTableStmt("", "t",
		makeATDropDefault("id"),
	))
	assertErrorCode(t, err, CodeInvalidObjectDefinition)
}

func TestDropDefaultOnGeneratedColumn(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "val", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	// Manually mark column as generated (stored).
	rel := c.GetRelation("", "t")
	rel.Columns[0].Generated = 's'

	err := c.AlterTableStmt(makeAlterTableStmt("", "t",
		makeATDropDefault("val"),
	))
	assertErrorCode(t, err, CodeInvalidObjectDefinition)
}

func TestAlterColumnTypeOnIdentityColumn(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	// Manually mark column as identity (ALWAYS).
	rel := c.GetRelation("", "t")
	rel.Columns[0].Identity = 'a'

	err := c.AlterTableStmt(makeAlterTableStmt("", "t",
		makeATAlterColumnType("id", TypeName{Name: "int8", TypeMod: -1}),
	))
	assertErrorCode(t, err, CodeInvalidObjectDefinition)
}

func TestAlterColumnTypeWithDependentView(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "name", Type: TypeName{Name: "text", TypeMod: -1}},
	}, nil, false), 'r')

	// Create a view depending on the table.
	c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprStar}}},
			From:       []*FromItem{{Kind: FromTable, Table: "t"}},
		}, nil))

	// Attempt to change column type on the base table.
	err := c.AlterTableStmt(makeAlterTableStmt("", "t",
		makeATAlterColumnType("id", TypeName{Name: "int8", TypeMod: -1}),
	))
	assertErrorCode(t, err, CodeFeatureNotSupported)
}

func TestDropInheritedColumn(t *testing.T) {
	c := New()

	// Create parent table.
	c.DefineRelation(makeCreateTableStmt("", "parent", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	// Create child table with column "a".
	c.DefineRelation(makeCreateTableStmt("", "child", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	// Set up inheritance: mark column as inherited (not locally defined).
	child := c.GetRelation("", "child")
	child.Columns[0].InhCount = 1
	child.Columns[0].IsLocal = false

	// Try to drop inherited column from child -- should error.
	err := c.AlterTableStmt(makeAlterTableStmt("", "child",
		makeATDropColumn("a", false, false),
	))
	assertErrorCode(t, err, CodeInvalidObjectDefinition)
}

func TestExecuteTruncateOnView(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	// Create a view.
	c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprStar}}},
			From:       []*FromItem{{Kind: FromTable, Table: "t"}},
		}, nil))

	// Try to TRUNCATE the view -- should error with wrong object type.
	truncStmt := &nodes.TruncateStmt{
		Relations: &nodes.List{Items: []nodes.Node{
			&nodes.RangeVar{Relname: "v"},
		}},
	}
	err := c.ExecuteTruncate(truncStmt)
	assertErrorCode(t, err, CodeWrongObjectType)
}

func TestPersistenceTracking(t *testing.T) {
	c := New()

	cols := []ColumnDef{{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}}}
	stmt := makeCreateTableStmt("", "t", cols, nil, false)
	stmt.Relation.Relpersistence = 'u' // UNLOGGED
	c.DefineRelation(stmt, 'r')

	rel := c.GetRelation("", "t")
	if rel.Persistence != 'u' {
		t.Errorf("persistence: got %c, want 'u'", rel.Persistence)
	}
}

func TestRenameTriggerDuplicate(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, nil, false), 'r')

	// Create a trigger function.
	c.CreateFunctionStmt(makeTriggerFuncStmt("", "trig_fn", "plpgsql", "BEGIN RETURN NEW; END;"))

	// Create two triggers.
	c.CreateTriggerStmt(makeCreateTrigStmt("", "t", "trig1", TriggerBefore, TriggerEventInsert, false, "", "trig_fn"))
	c.CreateTriggerStmt(makeCreateTrigStmt("", "t", "trig2", TriggerBefore, TriggerEventUpdate, false, "", "trig_fn"))

	// Try to rename trig1 to trig2 -- should error with duplicate.
	err := c.ExecRenameStmt(makeRenameTriggerStmt("", "t", "trig1", "trig2"))
	assertErrorCode(t, err, CodeDuplicateObject)
}

func TestRenameConstraintDuplicateValidation(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "name", Type: TypeName{Name: "text", TypeMod: -1}},
	}, []ConstraintDef{
		{Type: ConstraintPK, Columns: []string{"id"}},
		{Type: ConstraintUnique, Columns: []string{"name"}},
	}, false), 'r')

	rel := c.GetRelation("", "t")
	cons := c.ConstraintsOf(rel.OID)
	if len(cons) < 2 {
		t.Fatalf("expected >= 2 constraints, got %d", len(cons))
	}

	// Try to rename the first constraint to the second constraint's name.
	err := c.ExecRenameStmt(makeRenameConstraintStmt("", "t", cons[0].Name, cons[1].Name))
	assertErrorCode(t, err, CodeDuplicateObject)
}
