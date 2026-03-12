package ddl

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// =============================================================================
// Round 8: Type system hardening tests
// =============================================================================

// -----------------------------------------------------------------------------
// moveArrayTypeName: Array name collision for DefineEnum
// -----------------------------------------------------------------------------

func TestRound8EnumArrayNameCollisionAutoRename(t *testing.T) {
	c := New()
	// First, create an enum called "_myenum" so its array type is "__myenum".
	// Then create an enum called "myenum" — its array type "_myenum" would collide
	// with the first enum. The old "_myenum" should be renamed.
	stmt1 := &nodes.CreateEnumStmt{
		TypeName: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "_myenum"}}},
		Vals:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: "x"}}},
	}
	if err := c.DefineEnum(stmt1); err != nil {
		t.Fatal(err)
	}

	// Verify "_myenum" type and "__myenum" array type exist.
	oid1, _, err := c.ResolveType(TypeName{Name: "_myenum", TypeMod: -1})
	if err != nil {
		t.Fatal(err)
	}
	bt1 := c.typeByOID[oid1]
	if bt1.Type != 'e' {
		t.Errorf("expected enum type, got %c", bt1.Type)
	}

	// Now create "myenum" — array type "_myenum" collides with enum "_myenum".
	// Since "_myenum" is NOT an auto-generated array type (it's an enum), we should
	// NOT rename it, and the collision should fail.
	// Actually, let me re-read the requirement: the array type of "_myenum" enum is "__myenum".
	// When we create "myenum" enum, its array type would be "_myenum".
	// The existing "_myenum" is the enum type itself (Type='e', Category='E'), NOT an array type.
	// So moveArrayTypeName should return an error because it's not an auto-generated array.
	stmt2 := &nodes.CreateEnumStmt{
		TypeName: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "myenum"}}},
		Vals:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: "a"}}},
	}
	err = c.DefineEnum(stmt2)
	assertCode(t, err, CodeDuplicateObject)
}

func TestRound8EnumArrayNameCollisionAutoArrayRename(t *testing.T) {
	c := New()
	// Create an enum "foo" — this creates array type "_foo".
	stmt1 := &nodes.CreateEnumStmt{
		TypeName: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "foo"}}},
		Vals:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: "a"}}},
	}
	if err := c.DefineEnum(stmt1); err != nil {
		t.Fatal(err)
	}

	// Drop enum "foo" but keep array type "_foo" around by registering a new type "foo"
	// that has "_foo" as its array (simulating the scenario).
	// Actually, let's test the real scenario:
	// Create domain "foo" to test that when a domain creates array type "_foo",
	// if "_foo" already exists as an auto-generated array (from enum), it gets renamed.

	// The enum created "_foo" array type. Now create domain "foo" — this will
	// try to create "_foo" array type. But "_foo" already exists as an auto-generated
	// array type. moveArrayTypeName should rename it.
	// Wait — we can't create a type named "foo" because the enum already occupies it.

	// Better test: Create an enum "bar", then create a type named "_bar" which is
	// also an auto-generated array. This is a bit contrived, but let's use a
	// more realistic scenario:
	// 1. Create type "x" which gets array "_x"
	// 2. Drop type "x" (but array "_x" stays as orphan in name map)
	// Actually, let's just test that the rename works on the enum scenario properly.

	// The simplest real-world scenario: Create composite type first, then enum.
	// composite type "ct" gets array "_ct", then enum "ct" would fail because
	// "ct" already exists. But we want to test array collision specifically.

	// Let's just directly verify moveArrayTypeName works:
	schema := c.schemaByName["public"]

	// Manually register an array type "_testtype" that looks auto-generated.
	fakeArrayOID := c.oidGen.Next()
	fakeArray := &BuiltinType{
		OID:       fakeArrayOID,
		TypeName:  "_testtype",
		Namespace: schema.OID,
		Len:       -1,
		ByVal:     false,
		Type:      'b',
		Category:  'A',
		IsDefined: true,
		Delim:     ',',
		Elem:      999, // points to something
		Align:     'i',
		Storage:   'x',
		TypeMod:   -1,
	}
	c.typeByOID[fakeArrayOID] = fakeArray
	c.typeByName[typeKey{ns: schema.OID, name: "_testtype"}] = fakeArray

	// Now create enum "testtype" — its array "_testtype" collides with the existing array.
	// moveArrayTypeName should rename the old one to "_testtype_1".
	stmt := &nodes.CreateEnumStmt{
		TypeName: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "testtype"}}},
		Vals:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: "val1"}}},
	}
	if err := c.DefineEnum(stmt); err != nil {
		t.Fatalf("expected success, got: %v", err)
	}

	// The old array should have been renamed to "_testtype_1".
	renamed := c.typeByOID[fakeArrayOID]
	if renamed.TypeName != "_testtype_1" {
		t.Errorf("expected renamed to _testtype_1, got %q", renamed.TypeName)
	}
	// And the new enum's array should be "_testtype".
	enumOID, _, _ := c.ResolveType(TypeName{Name: "testtype", TypeMod: -1})
	enumBt := c.typeByOID[enumOID]
	newArray := c.typeByOID[enumBt.Array]
	if newArray.TypeName != "_testtype" {
		t.Errorf("expected new array to be _testtype, got %q", newArray.TypeName)
	}
}

// -----------------------------------------------------------------------------
// moveArrayTypeName: Array name collision for DefineDomain
// -----------------------------------------------------------------------------

func TestRound8DomainArrayNameCollisionAutoArrayRename(t *testing.T) {
	c := New()
	schema := c.schemaByName["public"]

	// Manually register an auto-generated array type "_mydom".
	fakeArrayOID := c.oidGen.Next()
	fakeArray := &BuiltinType{
		OID:       fakeArrayOID,
		TypeName:  "_mydom",
		Namespace: schema.OID,
		Len:       -1,
		ByVal:     false,
		Type:      'b',
		Category:  'A',
		IsDefined: true,
		Delim:     ',',
		Elem:      999,
		Align:     'i',
		Storage:   'x',
		TypeMod:   -1,
	}
	c.typeByOID[fakeArrayOID] = fakeArray
	c.typeByName[typeKey{ns: schema.OID, name: "_mydom"}] = fakeArray

	// Create domain "mydom" — its array type "_mydom" collides.
	stmt := makeCreateDomainStmt("", "mydom", TypeName{Name: "int4", TypeMod: -1}, false, "", "", "")
	if err := c.DefineDomain(stmt); err != nil {
		t.Fatalf("expected success, got: %v", err)
	}

	// The old array should have been renamed.
	renamed := c.typeByOID[fakeArrayOID]
	if renamed.TypeName != "_mydom_1" {
		t.Errorf("expected renamed to _mydom_1, got %q", renamed.TypeName)
	}
}

// -----------------------------------------------------------------------------
// moveArrayTypeName: Array name collision for DefineRange
// -----------------------------------------------------------------------------

func TestRound8RangeArrayNameCollisionAutoArrayRename(t *testing.T) {
	c := New()
	schema := c.schemaByName["public"]

	// Manually register an auto-generated array type "_myrange".
	fakeArrayOID := c.oidGen.Next()
	fakeArray := &BuiltinType{
		OID:       fakeArrayOID,
		TypeName:  "_myrange",
		Namespace: schema.OID,
		Len:       -1,
		ByVal:     false,
		Type:      'b',
		Category:  'A',
		IsDefined: true,
		Delim:     ',',
		Elem:      999,
		Align:     'i',
		Storage:   'x',
		TypeMod:   -1,
	}
	c.typeByOID[fakeArrayOID] = fakeArray
	c.typeByName[typeKey{ns: schema.OID, name: "_myrange"}] = fakeArray

	// Create range "myrange" — its array type "_myrange" collides.
	stmt := &nodes.CreateRangeStmt{
		TypeName: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "myrange"}}},
		Params: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{
				Defname: "subtype",
				Arg: &nodes.TypeName{
					Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "int4"}}},
				},
			},
		}},
	}
	if err := c.DefineRange(stmt); err != nil {
		t.Fatalf("expected success, got: %v", err)
	}

	// The old array should have been renamed.
	renamed := c.typeByOID[fakeArrayOID]
	if renamed.TypeName != "_myrange_1" {
		t.Errorf("expected renamed to _myrange_1, got %q", renamed.TypeName)
	}
}

// -----------------------------------------------------------------------------
// moveArrayTypeName: Non-array type collision should error
// -----------------------------------------------------------------------------

func TestRound8ArrayNameCollisionNonArrayError(t *testing.T) {
	c := New()
	// Create a non-array type named "_myenum2" (an enum, not an auto-generated array).
	stmt1 := &nodes.CreateEnumStmt{
		TypeName: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "_myenum2"}}},
		Vals:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: "x"}}},
	}
	if err := c.DefineEnum(stmt1); err != nil {
		t.Fatal(err)
	}

	// Try to create "myenum2" — its array "_myenum2" collides with the enum
	// (which is not an auto-generated array), so it should error.
	stmt2 := &nodes.CreateEnumStmt{
		TypeName: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "myenum2"}}},
		Vals:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: "a"}}},
	}
	err := c.DefineEnum(stmt2)
	assertCode(t, err, CodeDuplicateObject)
}

// -----------------------------------------------------------------------------
// Domain COLLATE validation: non-collatable type error
// -----------------------------------------------------------------------------

func TestRound8DomainCollateNonCollatableError(t *testing.T) {
	c := New()
	// CREATE DOMAIN d AS integer COLLATE "C" — int4 is not collatable.
	stmt := &nodes.CreateDomainStmt{
		Domainname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "d"}}},
		Typname:    makeTypeNameNode(TypeName{Name: "int4", TypeMod: -1}),
		CollClause: &nodes.CollateClause{
			Collname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "C"}}},
		},
	}
	err := c.DefineDomain(stmt)
	assertCode(t, err, CodeDatatypeMismatch)
	if !strings.Contains(err.Error(), "collations are not supported by type") {
		t.Errorf("unexpected message: %s", err)
	}
}

func TestRound8DomainCollateOnTextOK(t *testing.T) {
	c := New()
	// CREATE DOMAIN d AS text COLLATE "C" — text is collatable.
	stmt := &nodes.CreateDomainStmt{
		Domainname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "d"}}},
		Typname:    makeTypeNameNode(TypeName{Name: "text", TypeMod: -1}),
		CollClause: &nodes.CollateClause{
			Collname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "C"}}},
		},
	}
	err := c.DefineDomain(stmt)
	if err != nil {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------------------
// Range collation validation: non-collatable subtype error
// -----------------------------------------------------------------------------

func TestRound8RangeCollationNonCollatableError(t *testing.T) {
	c := New()
	// CREATE TYPE myrange AS RANGE (subtype = integer, collation = "C")
	// integer is not collatable — should error.
	stmt := &nodes.CreateRangeStmt{
		TypeName: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "myrange"}}},
		Params: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{
				Defname: "subtype",
				Arg: &nodes.TypeName{
					Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "int4"}}},
				},
			},
			&nodes.DefElem{
				Defname: "collation",
				Arg:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: "C"}}},
			},
		}},
	}
	err := c.DefineRange(stmt)
	assertCode(t, err, CodeDatatypeMismatch)
	if !strings.Contains(err.Error(), "collations are not supported by type") {
		t.Errorf("unexpected message: %s", err)
	}
}

func TestRound8RangeCollationOnTextOK(t *testing.T) {
	c := New()
	// CREATE TYPE myrange AS RANGE (subtype = text, collation = "C") — text is collatable.
	stmt := &nodes.CreateRangeStmt{
		TypeName: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "myrange"}}},
		Params: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{
				Defname: "subtype",
				Arg: &nodes.TypeName{
					Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "text"}}},
				},
			},
			&nodes.DefElem{
				Defname: "collation",
				Arg:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: "C"}}},
			},
		}},
	}
	err := c.DefineRange(stmt)
	if err != nil {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------------------
// Range duplicate parameter detection
// -----------------------------------------------------------------------------

func TestRound8RangeDuplicateSubtypeError(t *testing.T) {
	c := New()
	// CREATE TYPE myrange AS RANGE (subtype = integer, subtype = text) — duplicate
	stmt := &nodes.CreateRangeStmt{
		TypeName: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "myrange"}}},
		Params: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{
				Defname: "subtype",
				Arg: &nodes.TypeName{
					Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "int4"}}},
				},
			},
			&nodes.DefElem{
				Defname: "subtype",
				Arg: &nodes.TypeName{
					Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "text"}}},
				},
			},
		}},
	}
	err := c.DefineRange(stmt)
	assertCode(t, err, CodeInvalidParameterValue)
	if !strings.Contains(err.Error(), "conflicting or redundant") {
		t.Errorf("unexpected message: %s", err)
	}
}

func TestRound8RangeDuplicateCollationError(t *testing.T) {
	c := New()
	stmt := &nodes.CreateRangeStmt{
		TypeName: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "myrange"}}},
		Params: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{
				Defname: "subtype",
				Arg: &nodes.TypeName{
					Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "text"}}},
				},
			},
			&nodes.DefElem{
				Defname: "collation",
				Arg:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: "C"}}},
			},
			&nodes.DefElem{
				Defname: "collation",
				Arg:     &nodes.List{Items: []nodes.Node{&nodes.String{Str: "POSIX"}}},
			},
		}},
	}
	err := c.DefineRange(stmt)
	assertCode(t, err, CodeInvalidParameterValue)
	if !strings.Contains(err.Error(), "conflicting or redundant") {
		t.Errorf("unexpected message: %s", err)
	}
}

// -----------------------------------------------------------------------------
// ALTER DOMAIN VALIDATE CONSTRAINT
// -----------------------------------------------------------------------------

func TestRound8AlterDomainValidateConstraint(t *testing.T) {
	c := New()
	// Create domain with CHECK constraint.
	stmt := makeCreateDomainStmt("", "d", TypeName{Name: "int4", TypeMod: -1}, false, "", "d_check", "VALUE > 0")
	if err := c.DefineDomain(stmt); err != nil {
		t.Fatal(err)
	}

	oid, _, _ := c.ResolveType(TypeName{Name: "d", TypeMod: -1})
	dt := c.DomainInfo(oid)
	// pg: CREATE DOMAIN always creates constraints with convalidated = true
	if !dt.Constraints[0].ConValidated {
		t.Error("constraint should be validated at creation time (PG behavior)")
	}

	// ALTER DOMAIN d VALIDATE CONSTRAINT d_check (no-op when already validated)
	alter := &nodes.AlterDomainStmt{
		Subtype: 'V',
		Typname: makeDomainTypname("", "d"),
		Name:    "d_check",
	}
	if err := c.AlterDomainStmt(alter); err != nil {
		t.Fatal(err)
	}

	if !dt.Constraints[0].ConValidated {
		t.Error("constraint should still be validated after ALTER DOMAIN VALIDATE CONSTRAINT")
	}
}

func TestRound8AlterDomainValidateConstraintNotFound(t *testing.T) {
	c := New()
	stmt := makeCreateDomainStmt("", "d", TypeName{Name: "int4", TypeMod: -1}, false, "", "", "")
	if err := c.DefineDomain(stmt); err != nil {
		t.Fatal(err)
	}

	alter := &nodes.AlterDomainStmt{
		Subtype: 'V',
		Typname: makeDomainTypname("", "d"),
		Name:    "nonexistent",
	}
	err := c.AlterDomainStmt(alter)
	assertCode(t, err, CodeUndefinedObject)
	if !strings.Contains(err.Error(), "nonexistent") {
		t.Errorf("unexpected message: %s", err)
	}
}

// -----------------------------------------------------------------------------
// Range without collation on non-collatable type should succeed (no error)
// -----------------------------------------------------------------------------

func TestRound8RangeNoCollationOnIntOK(t *testing.T) {
	c := New()
	stmt := &nodes.CreateRangeStmt{
		TypeName: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "intrange"}}},
		Params: &nodes.List{Items: []nodes.Node{
			&nodes.DefElem{
				Defname: "subtype",
				Arg: &nodes.TypeName{
					Names: &nodes.List{Items: []nodes.Node{&nodes.String{Str: "int4"}}},
				},
			},
		}},
	}
	err := c.DefineRange(stmt)
	if err != nil {
		t.Fatal(err)
	}
}
