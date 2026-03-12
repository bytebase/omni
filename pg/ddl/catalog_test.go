package ddl

import "testing"

func TestNewCatalogIndexes(t *testing.T) {
	c := New()

	// Verify index sizes match generated data.
	if got := len(c.typeByOID); got != len(BuiltinTypes) {
		t.Errorf("typeByOID: got %d, want %d", got, len(BuiltinTypes))
	}
	if got := len(c.castIndex); got != len(BuiltinCasts) {
		t.Errorf("castIndex: got %d, want %d", got, len(BuiltinCasts))
	}
	if got := len(c.operByOID); got != len(BuiltinOperators) {
		t.Errorf("operByOID: got %d, want %d", got, len(BuiltinOperators))
	}
	if got := len(c.procByOID); got != len(BuiltinProcs) {
		t.Errorf("procByOID: got %d, want %d", got, len(BuiltinProcs))
	}
}

func TestTypeLookup(t *testing.T) {
	c := New()

	// Spot-check BOOLOID.
	bt := c.TypeByOID(BOOLOID)
	if bt == nil || bt.TypeName != "bool" {
		t.Fatal("expected bool type for BOOLOID")
	}

	// Spot-check INT4OID.
	it := c.TypeByOID(INT4OID)
	if it == nil || it.TypeName != "int4" {
		t.Fatal("expected int4 type for INT4OID")
	}
}

func TestCastLookup(t *testing.T) {
	c := New()

	// int4 → int8 cast should exist and be implicit.
	cast := c.LookupCast(INT4OID, INT8OID)
	if cast == nil {
		t.Fatal("expected int4 → int8 cast")
	}
	if cast.Context != 'i' {
		t.Errorf("int4 → int8 cast context: got %c, want 'i'", cast.Context)
	}
}

func TestOperatorLookup(t *testing.T) {
	c := New()

	// int4 = int4 should exist.
	ops := c.LookupOperatorExact("=", INT4OID, INT4OID)
	if len(ops) == 0 {
		t.Fatal("expected int4 = int4 operator")
	}
	if ops[0].Result != BOOLOID {
		t.Errorf("int4 = int4 result: got %d, want %d", ops[0].Result, BOOLOID)
	}
}

func TestProcLookup(t *testing.T) {
	c := New()

	// "int4in" should exist.
	procs := c.LookupProcByName("int4in")
	if len(procs) == 0 {
		t.Fatal("expected int4in proc")
	}

	// Lookup by OID.
	p := c.LookupProcByOID(procs[0].OID)
	if p == nil || p.Name != "int4in" {
		t.Fatal("expected int4in proc by OID")
	}
}

func TestBuiltinSchemas(t *testing.T) {
	c := New()

	for _, name := range []string{"pg_catalog", "pg_toast", "public"} {
		if s := c.GetSchema(name); s == nil {
			t.Errorf("expected built-in schema %q", name)
		}
	}
}
