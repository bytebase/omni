package catalog

import "testing"

func TestResolveBasicType(t *testing.T) {
	c := New()

	oid, typmod, err := c.ResolveType(TypeName{Name: "int4", TypeMod: -1})
	if err != nil {
		t.Fatal(err)
	}
	if oid != INT4OID {
		t.Errorf("got OID %d, want %d", oid, INT4OID)
	}
	if typmod != -1 {
		t.Errorf("got typmod %d, want -1", typmod)
	}
}

// TestResolveTypeAsWritten holds ResolveType to the catalog name: the
// grammar maps an unquoted SQL-standard name (integer) to int4, and a quoted
// "integer" names no type.
func TestResolveTypeAsWritten(t *testing.T) {
	c := New()

	tests := []struct {
		name    string
		wantOID uint32
	}{
		{"int4", INT4OID},
		{"bool", BOOLOID},
		{"int8", INT8OID},
		{"text", TEXTOID},
		{"integer", 0},
		{"boolean", 0},
		{"double precision", 0},
	}
	for _, tt := range tests {
		oid, _, err := c.ResolveType(TypeName{Name: tt.name, TypeMod: -1})
		if tt.wantOID == 0 {
			assertErrorCode(t, err, CodeUndefinedObject)
			continue
		}
		if err != nil {
			t.Errorf("ResolveType(%q): %v", tt.name, err)
			continue
		}
		if oid != tt.wantOID {
			t.Errorf("ResolveType(%q): got OID %d, want %d", tt.name, oid, tt.wantOID)
		}
	}
}

func TestResolveSchemaQualified(t *testing.T) {
	c := New()

	oid, _, err := c.ResolveType(TypeName{Schema: "pg_catalog", Name: "int4", TypeMod: -1})
	if err != nil {
		t.Fatal(err)
	}
	if oid != INT4OID {
		t.Errorf("got OID %d, want %d", oid, INT4OID)
	}
}

func TestResolveArrayType(t *testing.T) {
	c := New()

	oid, _, err := c.ResolveType(TypeName{Name: "int4", TypeMod: -1, IsArray: true})
	if err != nil {
		t.Fatal(err)
	}
	if oid != INT4ARRAYOID {
		t.Errorf("got OID %d, want %d (int4[])", oid, INT4ARRAYOID)
	}
}

func TestResolveVarcharWithTypmod(t *testing.T) {
	c := New()

	oid, typmod, err := c.ResolveType(TypeName{Name: "varchar", TypeMod: 100})
	if err != nil {
		t.Fatal(err)
	}
	if oid != VARCHAROID {
		t.Errorf("got OID %d, want %d", oid, VARCHAROID)
	}
	if typmod != 100 {
		t.Errorf("got typmod %d, want 100", typmod)
	}
}

func TestResolveTypmodRejected(t *testing.T) {
	c := New()

	// int4 does not accept type modifiers.
	_, _, err := c.ResolveType(TypeName{Name: "integer", TypeMod: 42})
	if err == nil {
		t.Fatal("expected error for typmod on integer")
	}
}

func TestResolveUndefinedType(t *testing.T) {
	c := New()

	_, _, err := c.ResolveType(TypeName{Name: "nosuchtype", TypeMod: -1})
	assertErrorCode(t, err, CodeUndefinedObject)
}
