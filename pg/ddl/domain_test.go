package ddl

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// ---------------------------------------------------------------------------
// domain_test helpers
// ---------------------------------------------------------------------------

func makeCreateDomainStmt(schema, name string, baseType TypeName, notNull bool, defaultExpr, checkName, checkExpr string) *nodes.CreateDomainStmt {
	var nameItems []nodes.Node
	if schema != "" {
		nameItems = []nodes.Node{&nodes.String{Str: schema}, &nodes.String{Str: name}}
	} else {
		nameItems = []nodes.Node{&nodes.String{Str: name}}
	}
	stmt := &nodes.CreateDomainStmt{
		Domainname: &nodes.List{Items: nameItems},
		Typname:    makeTypeNameNode(baseType),
	}
	var cons []nodes.Node
	if notNull {
		cons = append(cons, &nodes.Constraint{Contype: nodes.CONSTR_NOTNULL})
	}
	if defaultExpr != "" {
		cons = append(cons, &nodes.Constraint{Contype: nodes.CONSTR_DEFAULT, CookedExpr: defaultExpr})
	}
	if checkExpr != "" {
		cons = append(cons, &nodes.Constraint{Contype: nodes.CONSTR_CHECK, Conname: checkName, CookedExpr: checkExpr})
	}
	if len(cons) > 0 {
		stmt.Constraints = &nodes.List{Items: cons}
	}
	return stmt
}

func makeDomainTypname(schema, name string) *nodes.List {
	var items []nodes.Node
	if schema != "" {
		items = []nodes.Node{&nodes.String{Str: schema}, &nodes.String{Str: name}}
	} else {
		items = []nodes.Node{&nodes.String{Str: name}}
	}
	return &nodes.List{Items: items}
}

func makeAlterDomainSetNotNull(schema, name string) *nodes.AlterDomainStmt {
	return &nodes.AlterDomainStmt{Subtype: 'N', Typname: makeDomainTypname(schema, name)}
}

func makeAlterDomainDropNotNull(schema, name string) *nodes.AlterDomainStmt {
	return &nodes.AlterDomainStmt{Subtype: 'O', Typname: makeDomainTypname(schema, name)}
}

func makeAlterDomainSetDefault(schema, name, defaultExpr string) *nodes.AlterDomainStmt {
	return &nodes.AlterDomainStmt{
		Subtype: 'T',
		Typname: makeDomainTypname(schema, name),
		Def:     &nodes.Constraint{Contype: nodes.CONSTR_DEFAULT, CookedExpr: defaultExpr},
	}
}

func makeAlterDomainDropDefault(schema, name string) *nodes.AlterDomainStmt {
	return &nodes.AlterDomainStmt{Subtype: 'T', Typname: makeDomainTypname(schema, name)}
}

func makeAlterDomainAddCheck(schema, name, checkName, checkExpr string) *nodes.AlterDomainStmt {
	return &nodes.AlterDomainStmt{
		Subtype: 'C',
		Typname: makeDomainTypname(schema, name),
		Def:     &nodes.Constraint{Contype: nodes.CONSTR_CHECK, Conname: checkName, CookedExpr: checkExpr},
	}
}

func makeAlterDomainDropConstraint(schema, name, constraintName string, cascade bool) *nodes.AlterDomainStmt {
	behavior := nodes.DROP_RESTRICT
	if cascade {
		behavior = nodes.DROP_CASCADE
	}
	return &nodes.AlterDomainStmt{
		Subtype:  'X',
		Typname:  makeDomainTypname(schema, name),
		Name:     constraintName,
		Behavior: behavior,
	}
}

// makeDropTypeStmt is defined in enum_test.go (shared by both).

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestCreateDomain(t *testing.T) {
	c := New()
	err := c.DefineDomain(makeCreateDomainStmt("", "posint", TypeName{Name: "integer", TypeMod: -1}, false, "", "", ""))
	if err != nil {
		t.Fatal(err)
	}

	oid, _, err := c.ResolveType(TypeName{Name: "posint", TypeMod: -1})
	if err != nil {
		t.Fatal(err)
	}
	bt := c.typeByOID[oid]
	if bt.Type != 'd' {
		t.Errorf("type kind: got %c, want 'd'", bt.Type)
	}
	if bt.BaseType != INT4OID {
		t.Errorf("base type: got %d, want %d", bt.BaseType, INT4OID)
	}
}

func TestCreateDomainNotNull(t *testing.T) {
	c := New()
	err := c.DefineDomain(makeCreateDomainStmt("", "nn_text", TypeName{Name: "text", TypeMod: -1}, true, "", "", ""))
	if err != nil {
		t.Fatal(err)
	}

	oid, _, _ := c.ResolveType(TypeName{Name: "nn_text", TypeMod: -1})
	dt := c.DomainInfo(oid)
	if !dt.NotNull {
		t.Error("domain should be NOT NULL")
	}
}

func TestCreateDomainWithCheck(t *testing.T) {
	c := New()
	err := c.DefineDomain(makeCreateDomainStmt("", "posint", TypeName{Name: "integer", TypeMod: -1}, false, "", "", "VALUE > 0"))
	if err != nil {
		t.Fatal(err)
	}

	oid, _, _ := c.ResolveType(TypeName{Name: "posint", TypeMod: -1})
	dt := c.DomainInfo(oid)
	if len(dt.Constraints) != 1 {
		t.Fatalf("constraints: got %d, want 1", len(dt.Constraints))
	}
	if dt.Constraints[0].Name != "posint_check" {
		t.Errorf("constraint name: got %q, want %q", dt.Constraints[0].Name, "posint_check")
	}
}

func TestCreateDomainWithDefault(t *testing.T) {
	c := New()
	err := c.DefineDomain(makeCreateDomainStmt("", "myint", TypeName{Name: "integer", TypeMod: -1}, false, "42", "", ""))
	if err != nil {
		t.Fatal(err)
	}

	oid, _, _ := c.ResolveType(TypeName{Name: "myint", TypeMod: -1})
	dt := c.DomainInfo(oid)
	if dt.Default != "42" {
		t.Errorf("default: got %q, want %q", dt.Default, "42")
	}
}

func TestCreateDomainDuplicate(t *testing.T) {
	c := New()
	c.DefineDomain(makeCreateDomainStmt("", "d", TypeName{Name: "int4", TypeMod: -1}, false, "", "", ""))
	err := c.DefineDomain(makeCreateDomainStmt("", "d", TypeName{Name: "int4", TypeMod: -1}, false, "", "", ""))
	assertCode(t, err, CodeDuplicateObject)
}

func TestDomainAsColumnType(t *testing.T) {
	c := New()
	c.DefineDomain(makeCreateDomainStmt("", "posint", TypeName{Name: "integer", TypeMod: -1}, false, "", "", ""))
	err := c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "val", Type: TypeName{Name: "posint", TypeMod: -1}},
	}, nil, false), 'r')
	if err != nil {
		t.Fatal(err)
	}
}

func TestAlterDomainSetNotNull(t *testing.T) {
	c := New()
	c.DefineDomain(makeCreateDomainStmt("", "d", TypeName{Name: "int4", TypeMod: -1}, false, "", "", ""))

	err := c.AlterDomainStmt(makeAlterDomainSetNotNull("", "d"))
	if err != nil {
		t.Fatal(err)
	}

	oid, _, _ := c.ResolveType(TypeName{Name: "d", TypeMod: -1})
	dt := c.DomainInfo(oid)
	if !dt.NotNull {
		t.Error("should be NOT NULL")
	}
}

func TestAlterDomainDropNotNull(t *testing.T) {
	c := New()
	c.DefineDomain(makeCreateDomainStmt("", "d", TypeName{Name: "int4", TypeMod: -1}, true, "", "", ""))

	err := c.AlterDomainStmt(makeAlterDomainDropNotNull("", "d"))
	if err != nil {
		t.Fatal(err)
	}

	oid, _, _ := c.ResolveType(TypeName{Name: "d", TypeMod: -1})
	dt := c.DomainInfo(oid)
	if dt.NotNull {
		t.Error("should not be NOT NULL")
	}
}

func TestAlterDomainSetDropDefault(t *testing.T) {
	c := New()
	c.DefineDomain(makeCreateDomainStmt("", "d", TypeName{Name: "int4", TypeMod: -1}, false, "", "", ""))

	c.AlterDomainStmt(makeAlterDomainSetDefault("", "d", "42"))
	oid, _, _ := c.ResolveType(TypeName{Name: "d", TypeMod: -1})
	dt := c.DomainInfo(oid)
	if dt.Default != "42" {
		t.Errorf("default: got %q, want %q", dt.Default, "42")
	}

	c.AlterDomainStmt(makeAlterDomainDropDefault("", "d"))
	if dt.Default != "" {
		t.Errorf("default: got %q, want empty", dt.Default)
	}
}

func TestAlterDomainAddDropConstraint(t *testing.T) {
	c := New()
	c.DefineDomain(makeCreateDomainStmt("", "d", TypeName{Name: "int4", TypeMod: -1}, false, "", "", ""))

	err := c.AlterDomainStmt(makeAlterDomainAddCheck("", "d", "d_pos", "VALUE > 0"))
	if err != nil {
		t.Fatal(err)
	}

	oid, _, _ := c.ResolveType(TypeName{Name: "d", TypeMod: -1})
	dt := c.DomainInfo(oid)
	if len(dt.Constraints) != 1 {
		t.Fatalf("constraints: got %d, want 1", len(dt.Constraints))
	}

	err = c.AlterDomainStmt(makeAlterDomainDropConstraint("", "d", "d_pos", false))
	if err != nil {
		t.Fatal(err)
	}
	if len(dt.Constraints) != 0 {
		t.Fatalf("constraints after drop: got %d, want 0", len(dt.Constraints))
	}
}

func TestAlterDomainDropConstraintNotFound(t *testing.T) {
	c := New()
	c.DefineDomain(makeCreateDomainStmt("", "d", TypeName{Name: "int4", TypeMod: -1}, false, "", "", ""))
	err := c.AlterDomainStmt(makeAlterDomainDropConstraint("", "d", "nope", false))
	assertCode(t, err, CodeUndefinedObject)
}

func TestDropTypeDomain(t *testing.T) {
	c := New()
	c.DefineDomain(makeCreateDomainStmt("", "d", TypeName{Name: "int4", TypeMod: -1}, false, "", "", ""))
	err := c.RemoveObjects(makeDropTypeStmt("", "d", false, false))
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = c.ResolveType(TypeName{Name: "d", TypeMod: -1})
	assertCode(t, err, CodeUndefinedObject)
}

func TestDropTypeDomainCascade(t *testing.T) {
	c := New()
	c.DefineDomain(makeCreateDomainStmt("", "d", TypeName{Name: "int4", TypeMod: -1}, false, "", "", ""))
	c.DefineRelation(makeCreateTableStmt("", "t", []ColumnDef{
		{Name: "val", Type: TypeName{Name: "d", TypeMod: -1}},
	}, nil, false), 'r')

	err := c.RemoveObjects(makeDropTypeStmt("", "d", false, false))
	assertCode(t, err, CodeDependentObjects)

	err = c.RemoveObjects(makeDropTypeStmt("", "d", false, true))
	if err != nil {
		t.Fatal(err)
	}
	if c.GetRelation("", "t") != nil {
		t.Error("dependent table should be dropped after CASCADE")
	}
}

func TestDomainArrayType(t *testing.T) {
	c := New()
	c.DefineDomain(makeCreateDomainStmt("", "d", TypeName{Name: "int4", TypeMod: -1}, false, "", "", ""))

	oid, _, err := c.ResolveType(TypeName{Name: "d", TypeMod: -1, IsArray: true})
	if err != nil {
		t.Fatal(err)
	}
	at := c.typeByOID[oid]
	if at.Category != 'A' {
		t.Errorf("array category: got %c, want 'A'", at.Category)
	}
}
