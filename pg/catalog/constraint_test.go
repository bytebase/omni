package catalog

import "testing"

func TestCreateTableWithPK(t *testing.T) {
	c := New()

	err := c.DefineRelation(makeCreateTableStmt("", "users", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "integer", TypeMod: -1}},
		{Name: "name", Type: TypeName{Name: "text", TypeMod: -1}},
	}, []ConstraintDef{
		{Type: ConstraintPK, Columns: []string{"id"}},
	}, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "users")
	if rel == nil {
		t.Fatal("table not found")
	}

	// PK columns should be NOT NULL.
	if !rel.Columns[0].NotNull {
		t.Error("PK column should be NOT NULL")
	}

	// Constraint should exist.
	cons := c.ConstraintsOf(rel.OID)
	if len(cons) != 1 {
		t.Fatalf("expected 1 constraint, got %d", len(cons))
	}
	if cons[0].Type != ConstraintPK {
		t.Errorf("expected PK constraint, got %c", cons[0].Type)
	}
	if cons[0].Name != "users_pkey" {
		t.Errorf("constraint name: got %q, want %q", cons[0].Name, "users_pkey")
	}

	// Backing index should exist.
	idxs := c.IndexesOf(rel.OID)
	if len(idxs) != 1 {
		t.Fatalf("expected 1 index, got %d", len(idxs))
	}
	if !idxs[0].IsPrimary {
		t.Error("index should be primary")
	}
	if !idxs[0].IsUnique {
		t.Error("index should be unique")
	}
}

func TestCreateTableDuplicatePK(t *testing.T) {
	c := New()

	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "b", Type: TypeName{Name: "int4", TypeMod: -1}},
	}, []ConstraintDef{
		{Type: ConstraintPK, Columns: []string{"a"}},
		{Type: ConstraintPK, Columns: []string{"b"}},
	}, false), 'r')
	assertErrorCode(t, err, CodeDuplicatePKey)

	// Table should have been rolled back.
	if r := c.GetRelation("", "t"); r != nil {
		t.Error("table should have been rolled back after constraint error")
	}
}

func TestCreateTableWithUnique(t *testing.T) {
	c := New()

	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "email", Type: TypeName{Name: "text", TypeMod: -1}},
	}, []ConstraintDef{
		{Type: ConstraintUnique, Columns: []string{"email"}},
	}, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	cons := c.ConstraintsOf(rel.OID)
	if len(cons) != 1 {
		t.Fatalf("expected 1 constraint, got %d", len(cons))
	}
	if cons[0].Type != ConstraintUnique {
		t.Errorf("expected UNIQUE constraint, got %c", cons[0].Type)
	}
	if cons[0].Name != "t_email_key" {
		t.Errorf("constraint name: got %q, want %q", cons[0].Name, "t_email_key")
	}

	// Backing index.
	idxs := c.IndexesOf(rel.OID)
	if len(idxs) != 1 {
		t.Fatalf("expected 1 index, got %d", len(idxs))
	}
	if !idxs[0].IsUnique {
		t.Error("index should be unique")
	}
}

func TestCreateTableWithFK(t *testing.T) {
	c := New()

	// Create referenced table with PK.
	c.DefineRelation(makeCreateTableStmt("", "parents", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "integer", TypeMod: -1}},
	}, []ConstraintDef{
		{Type: ConstraintPK, Columns: []string{"id"}},
	}, false), 'r')

	// Create table with FK.
	err := c.DefineRelation(makeCreateTableStmt("", "children", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "integer", TypeMod: -1}},
		{Name: "parent_id", Type: TypeName{Name: "integer", TypeMod: -1}},
	}, []ConstraintDef{
		{Type: ConstraintFK, Columns: []string{"parent_id"}, RefTable: "parents"},
	}, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "children")
	cons := c.ConstraintsOf(rel.OID)
	var fk *Constraint
	for _, con := range cons {
		if con.Type == ConstraintFK {
			fk = con
		}
	}
	if fk == nil {
		t.Fatal("FK constraint not found")
	}
	if fk.Name != "children_parent_id_fkey" {
		t.Errorf("FK name: got %q, want %q", fk.Name, "children_parent_id_fkey")
	}
}

func TestFKIncompatibleTypes(t *testing.T) {
	c := New()

	c.DefineRelation(makeCreateTableStmt("", "parents", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "integer", TypeMod: -1}},
	}, []ConstraintDef{
		{Type: ConstraintPK, Columns: []string{"id"}},
	}, false), 'r')

	err := c.DefineRelation(makeCreateTableStmt("", "children", []ColumnDef{
		{Name: "parent_id", Type: TypeName{Name: "boolean", TypeMod: -1}},
	}, []ConstraintDef{
		{Type: ConstraintFK, Columns: []string{"parent_id"}, RefTable: "parents"},
	}, false), 'r')
	assertErrorCode(t, err, CodeDatatypeMismatch)
}

func TestFKNonExistentTable(t *testing.T) {
	c := New()

	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "ref_id", Type: TypeName{Name: "integer", TypeMod: -1}},
	}, []ConstraintDef{
		{Type: ConstraintFK, Columns: []string{"ref_id"}, RefTable: "nosuch"},
	}, false), 'r')
	assertErrorCode(t, err, CodeUndefinedTable)
}

func TestFKWithoutUniqueConstraint(t *testing.T) {
	c := New()

	// Create referenced table WITHOUT PK/UNIQUE.
	c.DefineRelation(makeCreateTableStmt("", "parents", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "integer", TypeMod: -1}},
	}, nil, false), 'r')

	err := c.DefineRelation(makeCreateTableStmt("", "children", []ColumnDef{
		{Name: "parent_id", Type: TypeName{Name: "integer", TypeMod: -1}},
	}, []ConstraintDef{
		{Type: ConstraintFK, Columns: []string{"parent_id"}, RefTable: "parents", RefColumns: []string{"id"}},
	}, false), 'r')
	assertErrorCode(t, err, CodeInvalidFK)
}

func TestFKEmptyRefColumnsUsesPK(t *testing.T) {
	c := New()

	c.DefineRelation(makeCreateTableStmt("", "parents", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "integer", TypeMod: -1}},
	}, []ConstraintDef{
		{Type: ConstraintPK, Columns: []string{"id"}},
	}, false), 'r')

	// FK without RefColumns -> should use PK of parents.
	err := c.DefineRelation(makeCreateTableStmt("", "children", []ColumnDef{
		{Name: "parent_id", Type: TypeName{Name: "integer", TypeMod: -1}},
	}, []ConstraintDef{
		{Type: ConstraintFK, Columns: []string{"parent_id"}, RefTable: "parents"},
	}, false), 'r')
	if err != nil {
		t.Fatal(err)
	}
}

// TestFKReferencesStandaloneUniqueIndex covers the pg_dump pattern: a unique
// index defined as a separate CREATE UNIQUE INDEX statement (not as an inline
// UNIQUE constraint) must be accepted as a FK referential target.
//
// pg: src/backend/commands/tablecmds.c — transformFkeyCheckAttrs
// PostgreSQL accepts a non-partial unique index over the referenced columns
// regardless of whether an explicit UNIQUE constraint exists. pg_dump relies
// on this when emitting CREATE UNIQUE INDEX as a separate statement.
func TestFKReferencesStandaloneUniqueIndex(t *testing.T) {
	c := New()

	// Referenced table has only a PK on a different column.
	c.DefineRelation(makeCreateTableStmt("", "principal", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "integer", TypeMod: -1}},
		{Name: "email", Type: TypeName{Name: "text", TypeMod: -1}},
	}, []ConstraintDef{
		{Type: ConstraintPK, Columns: []string{"id"}},
	}, false), 'r')

	// Standalone CREATE UNIQUE INDEX on email — no UNIQUE constraint object.
	if err := c.DefineIndex(makeIndexStmt("", "principal", "idx_principal_email", []string{"email"}, true, false)); err != nil {
		t.Fatal(err)
	}

	// FK references principal(email) — should be accepted because the unique
	// index satisfies the referential target requirement.
	err := c.DefineRelation(makeCreateTableStmt("", "session", []ColumnDef{
		{Name: "user_email", Type: TypeName{Name: "text", TypeMod: -1}},
	}, []ConstraintDef{
		{Type: ConstraintFK, Columns: []string{"user_email"}, RefTable: "principal", RefColumns: []string{"email"}},
	}, false), 'r')
	if err != nil {
		t.Fatalf("FK against standalone unique index should be accepted, got: %v", err)
	}
}

// TestFKRejectsNonUniqueIndex confirms that a plain (non-unique) index over
// the referenced columns does NOT satisfy the FK target requirement, matching
// PostgreSQL's behavior.
func TestFKRejectsNonUniqueIndex(t *testing.T) {
	c := New()

	c.DefineRelation(makeCreateTableStmt("", "parents", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "integer", TypeMod: -1}},
		{Name: "label", Type: TypeName{Name: "text", TypeMod: -1}},
	}, []ConstraintDef{
		{Type: ConstraintPK, Columns: []string{"id"}},
	}, false), 'r')

	// Non-unique index on label.
	if err := c.DefineIndex(makeIndexStmt("", "parents", "idx_parents_label", []string{"label"}, false, false)); err != nil {
		t.Fatal(err)
	}

	err := c.DefineRelation(makeCreateTableStmt("", "children", []ColumnDef{
		{Name: "parent_label", Type: TypeName{Name: "text", TypeMod: -1}},
	}, []ConstraintDef{
		{Type: ConstraintFK, Columns: []string{"parent_label"}, RefTable: "parents", RefColumns: []string{"label"}},
	}, false), 'r')
	assertErrorCode(t, err, CodeInvalidFK)
}

func TestCheckConstraint(t *testing.T) {
	c := New()

	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "age", Type: TypeName{Name: "integer", TypeMod: -1}},
	}, []ConstraintDef{
		{Type: ConstraintCheck, Columns: []string{"age"}, CheckExpr: "age > 0"},
	}, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	cons := c.ConstraintsOf(rel.OID)
	if len(cons) != 1 {
		t.Fatalf("expected 1 constraint, got %d", len(cons))
	}
	if cons[0].Type != ConstraintCheck {
		t.Errorf("expected CHECK constraint, got %c", cons[0].Type)
	}
	if cons[0].CheckExpr != "age > 0" {
		t.Errorf("check expression: got %q, want %q", cons[0].CheckExpr, "age > 0")
	}
	if cons[0].Name != "t_age_check" {
		t.Errorf("constraint name: got %q, want %q", cons[0].Name, "t_age_check")
	}
}

func TestDropTableCascadeDropsFK(t *testing.T) {
	c := New()

	c.DefineRelation(makeCreateTableStmt("", "parents", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "integer", TypeMod: -1}},
	}, []ConstraintDef{
		{Type: ConstraintPK, Columns: []string{"id"}},
	}, false), 'r')
	c.DefineRelation(makeCreateTableStmt("", "children", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "integer", TypeMod: -1}},
		{Name: "parent_id", Type: TypeName{Name: "integer", TypeMod: -1}},
	}, []ConstraintDef{
		{Type: ConstraintFK, Columns: []string{"parent_id"}, RefTable: "parents"},
	}, false), 'r')

	// CASCADE should drop the FK constraint from children.
	err := c.RemoveRelations(makeDropTableStmt("", "parents", false, true))
	if err != nil {
		t.Fatal(err)
	}

	// Children table should still exist but FK constraint should be gone.
	childRel := c.GetRelation("", "children")
	if childRel == nil {
		t.Fatal("children table should still exist")
	}
	cons := c.ConstraintsOf(childRel.OID)
	for _, con := range cons {
		if con.Type == ConstraintFK {
			t.Error("FK constraint should have been dropped by CASCADE")
		}
	}
}

func TestDropTableRestrictWithFK(t *testing.T) {
	c := New()

	c.DefineRelation(makeCreateTableStmt("", "parents", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "integer", TypeMod: -1}},
	}, []ConstraintDef{
		{Type: ConstraintPK, Columns: []string{"id"}},
	}, false), 'r')
	c.DefineRelation(makeCreateTableStmt("", "children", []ColumnDef{
		{Name: "parent_id", Type: TypeName{Name: "integer", TypeMod: -1}},
	}, []ConstraintDef{
		{Type: ConstraintFK, Columns: []string{"parent_id"}, RefTable: "parents"},
	}, false), 'r')

	err := c.RemoveRelations(makeDropTableStmt("", "parents", false, false))
	assertErrorCode(t, err, CodeDependentObjects)
}

func TestAutoGeneratedConstraintNames(t *testing.T) {
	c := New()

	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "a", Type: TypeName{Name: "integer", TypeMod: -1}},
		{Name: "b", Type: TypeName{Name: "integer", TypeMod: -1}},
	}, []ConstraintDef{
		{Type: ConstraintPK, Columns: []string{"a"}},
		{Type: ConstraintUnique, Columns: []string{"b"}},
	}, false), 'r')
	if err != nil {
		t.Fatal(err)
	}

	rel := c.GetRelation("", "t")
	cons := c.ConstraintsOf(rel.OID)
	names := make(map[string]bool)
	for _, con := range cons {
		names[con.Name] = true
	}
	if !names["t_pkey"] {
		t.Error("expected constraint name t_pkey")
	}
	if !names["t_b_key"] {
		t.Error("expected constraint name t_b_key")
	}
}
