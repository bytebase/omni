package ddl

import (
	"strings"
	"testing"
)

// =============================================================================
// Phase 5: Type System Deepening Tests
// =============================================================================

// defaultSchemaOID returns the OID of the default schema (first in search path).
func defaultSchemaOID(t *testing.T, c *Catalog) uint32 {
	t.Helper()
	s, err := c.resolveTargetSchema("")
	if err != nil {
		t.Fatalf("resolveTargetSchema: %v", err)
	}
	return s.OID
}

// --- 5c. Array coercion ---

func TestPhase5_ArrayCoercion(t *testing.T) {
	c := New()

	// int4[] → int8[] should be possible via element implicit coercion.
	int4ArrayOID := c.findArrayType(INT4OID)
	int8ArrayOID := c.findArrayType(INT8OID)
	if int4ArrayOID == 0 || int8ArrayOID == 0 {
		t.Skip("missing array types for INT4/INT8")
	}

	if !c.CanCoerce(int4ArrayOID, int8ArrayOID, 'i') {
		t.Errorf("expected int4[] → int8[] implicit coercion to be possible")
	}
}

func TestPhase5_ArrayCoercionExplicit(t *testing.T) {
	c := New()

	// text[] → int4[] should NOT be possible implicitly.
	textArrayOID := c.findArrayType(TEXTOID)
	int4ArrayOID := c.findArrayType(INT4OID)
	if textArrayOID == 0 || int4ArrayOID == 0 {
		t.Skip("missing array types for TEXT/INT4")
	}

	if c.CanCoerce(textArrayOID, int4ArrayOID, 'i') {
		t.Errorf("expected text[] → int4[] implicit coercion to NOT be possible")
	}
}

func TestPhase5_ArrayCoercionSameType(t *testing.T) {
	c := New()

	// int4[] → int4[] should be relabel (same type).
	int4ArrayOID := c.findArrayType(INT4OID)
	if int4ArrayOID == 0 {
		t.Skip("missing array type for INT4")
	}

	p, _ := c.FindCoercionPathway(int4ArrayOID, int4ArrayOID, 'i')
	if p != CoercionRelabel {
		t.Errorf("expected CoercionRelabel for same array type, got %d", p)
	}
}

func TestPhase5_ArrayCoercionInView(t *testing.T) {
	c := New()
	stmts := parseStmts(t, `
		CREATE TABLE arr_test (vals int[]);
		CREATE VIEW v_arr_cast AS SELECT vals::bigint[] AS bvals FROM arr_test;
	`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	rel := c.GetRelation("", "v_arr_cast")
	if rel == nil {
		t.Fatal("view not found")
	}
	if len(rel.Columns) < 1 {
		t.Fatal("expected at least 1 column")
	}
	int8ArrayOID := c.findArrayType(INT8OID)
	if int8ArrayOID == 0 {
		t.Skip("missing array type for INT8")
	}
	if rel.Columns[0].TypeOID != int8ArrayOID {
		t.Errorf("expected column type %d (int8[]), got %d", int8ArrayOID, rel.Columns[0].TypeOID)
	}
}

// --- 5d. Domain unwrapping ---

func TestPhase5_DomainCoercion(t *testing.T) {
	c := New()
	stmts := parseStmts(t, `
		CREATE DOMAIN myint AS int;
		CREATE TABLE dt (val myint);
	`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	// myint (domain over int) should be coercible to bigint via unwrapping.
	rel := c.GetRelation("", "dt")
	if rel == nil {
		t.Fatal("table not found")
	}
	myintOID := rel.Columns[0].TypeOID
	if myintOID == 0 {
		t.Fatal("myint OID not resolved")
	}

	if !c.CanCoerce(myintOID, INT8OID, 'i') {
		t.Errorf("expected myint → bigint implicit coercion via domain unwrapping")
	}
}

func TestPhase5_DomainUnwrapChain(t *testing.T) {
	c := New()
	stmts := parseStmts(t, `
		CREATE DOMAIN myint AS int;
		CREATE DOMAIN myint2 AS myint;
	`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	// Find myint2 OID.
	nsOID := defaultSchemaOID(t, c)
	bt := c.typeByName[typeKey{ns: nsOID, name: "myint2"}]
	if bt == nil {
		t.Fatal("myint2 type not found")
	}

	// getBaseType should unwrap myint2 → myint → int4.
	base := c.getBaseType(bt.OID)
	if base != INT4OID {
		t.Errorf("expected base type INT4OID (%d), got %d", INT4OID, base)
	}
}

func TestPhase5_DomainCoercionInView(t *testing.T) {
	c := New()
	stmts := parseStmts(t, `
		CREATE DOMAIN myint AS int;
		CREATE TABLE dt (val myint);
		CREATE VIEW v_domain AS SELECT val + 1.5 AS result FROM dt;
	`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	rel := c.GetRelation("", "v_domain")
	if rel == nil {
		t.Fatal("view not found")
	}
	if len(rel.Columns) < 1 {
		t.Fatal("expected at least 1 column")
	}
	// val (myint→int) + 1.5 (numeric) should produce a numeric type.
	if rel.Columns[0].TypeOID != NUMERICOID && rel.Columns[0].TypeOID != FLOAT8OID {
		t.Errorf("expected numeric or float8 result, got type OID %d", rel.Columns[0].TypeOID)
	}
}

// --- 5b. Multirange auto-creation ---

func TestPhase5_MultirangeAutoCreation(t *testing.T) {
	c := New()
	stmts := parseStmts(t, `CREATE TYPE intrange AS RANGE (subtype = int4);`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("create range: %v", err)
		}
	}

	nsOID := defaultSchemaOID(t, c)

	// The range type should exist.
	rangeType := c.typeByName[typeKey{ns: nsOID, name: "intrange"}]
	if rangeType == nil {
		t.Fatal("range type 'intrange' not found")
	}

	// The multirange type should be auto-created.
	mrType := c.typeByName[typeKey{ns: nsOID, name: "intrange_multirange"}]
	if mrType == nil {
		t.Fatal("multirange type 'intrange_multirange' not found")
	}
	if mrType.Type != 'm' {
		t.Errorf("expected multirange type kind 'm', got %c", mrType.Type)
	}

	// The multirange array type should exist.
	mrArrayType := c.typeByName[typeKey{ns: nsOID, name: "_intrange_multirange"}]
	if mrArrayType == nil {
		t.Fatal("multirange array type '_intrange_multirange' not found")
	}
	if mrArrayType.Category != 'A' {
		t.Errorf("expected array category 'A', got %c", mrArrayType.Category)
	}

	// RangeType metadata should link to multirange.
	ri := c.RangeInfo(rangeType.OID)
	if ri == nil {
		t.Fatal("RangeInfo not found")
	}
	if ri.MultirangeOID != mrType.OID {
		t.Errorf("expected MultirangeOID=%d, got %d", mrType.OID, ri.MultirangeOID)
	}
}

func TestPhase5_MultirangeArrayLink(t *testing.T) {
	c := New()
	stmts := parseStmts(t, `CREATE TYPE floatrange AS RANGE (subtype = float8);`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("create range: %v", err)
		}
	}

	nsOID := defaultSchemaOID(t, c)
	mrType := c.typeByName[typeKey{ns: nsOID, name: "floatrange_multirange"}]
	if mrType == nil {
		t.Fatal("multirange type not found")
	}

	// multirange.Array should point to its array type.
	if mrType.Array == 0 {
		t.Error("multirange type has no array OID")
	}
	mrArray := c.typeByOID[mrType.Array]
	if mrArray == nil {
		t.Error("multirange array type not found by OID")
	} else if mrArray.Elem != mrType.OID {
		t.Errorf("multirange array Elem=%d, expected %d", mrArray.Elem, mrType.OID)
	}
}

func TestPhase5_MultirangeDependency(t *testing.T) {
	c := New()
	stmts := parseStmts(t, `CREATE TYPE textrange AS RANGE (subtype = text);`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("create range: %v", err)
		}
	}

	nsOID := defaultSchemaOID(t, c)
	rangeType := c.typeByName[typeKey{ns: nsOID, name: "textrange"}]
	mrType := c.typeByName[typeKey{ns: nsOID, name: "textrange_multirange"}]
	if rangeType == nil || mrType == nil {
		t.Fatal("types not found")
	}

	// Verify dependency: multirange → range.
	found := false
	for _, dep := range c.deps {
		if dep.ObjOID == mrType.OID && dep.RefOID == rangeType.OID {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected dependency from multirange to range type")
	}
}

// --- getBaseType edge cases ---

func TestPhase5_GetBaseTypeNonDomain(t *testing.T) {
	c := New()
	// Non-domain type should return itself.
	base := c.getBaseType(INT4OID)
	if base != INT4OID {
		t.Errorf("expected INT4OID, got %d", base)
	}
}

func TestPhase5_FindCoercionPathwayDomainToTarget(t *testing.T) {
	c := New()
	stmts := parseStmts(t, `CREATE DOMAIN mytext AS text;`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	nsOID := defaultSchemaOID(t, c)
	mytextType := c.typeByName[typeKey{ns: nsOID, name: "mytext"}]
	if mytextType == nil {
		t.Fatal("mytext type not found")
	}

	// mytext (domain over text) → varchar should work via domain unwrapping.
	p, _ := c.FindCoercionPathway(mytextType.OID, VARCHAROID, 'i')
	if p == CoercionNone {
		t.Error("expected coercion from mytext to varchar")
	}
}

// --- Comprehensive test: verify existing range tests still work ---

func TestPhase5_ExistingRangeStillWorks(t *testing.T) {
	c := New()
	stmts := parseStmts(t, `
		CREATE TYPE numrange2 AS RANGE (subtype = numeric);
		CREATE TABLE t_range (r numrange2);
	`)
	for _, s := range stmts {
		if err := c.ProcessUtility(s); err != nil {
			t.Fatalf("setup: %v", err)
		}
	}

	rel := c.GetRelation("", "t_range")
	if rel == nil {
		t.Fatal("table not found")
	}
	if len(rel.Columns) != 1 {
		t.Fatalf("expected 1 column, got %d", len(rel.Columns))
	}

	// Verify the column type is the range type.
	nsOID := defaultSchemaOID(t, c)
	rangeType := c.typeByName[typeKey{ns: nsOID, name: "numrange2"}]
	if rangeType == nil {
		t.Fatal("range type not found")
	}
	if rel.Columns[0].TypeOID != rangeType.OID {
		t.Errorf("expected column type %d, got %d", rangeType.OID, rel.Columns[0].TypeOID)
	}

	// Check that multirange was also created.
	mrType := c.typeByName[typeKey{ns: nsOID, name: "numrange2_multirange"}]
	if mrType == nil {
		t.Error("multirange type not auto-created for range")
	}
}

// --- Helper to check view def contains expected text ---

func assertViewDefContains(t *testing.T, c *Catalog, viewName, expected string) {
	t.Helper()
	def, err := c.GetViewDefinition("", viewName)
	if err != nil {
		t.Fatalf("get view def: %v", err)
	}
	if !strings.Contains(def, expected) {
		t.Errorf("expected %q in view def, got: %s", expected, def)
	}
}
