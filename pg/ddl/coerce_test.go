package ddl

import "testing"

func TestCanCoerceSameType(t *testing.T) {
	c := New()
	if !c.CanCoerce(INT4OID, INT4OID, 'i') {
		t.Error("same type should always coerce")
	}
	if !c.CanCoerce(INT4OID, INT4OID, 'a') {
		t.Error("same type should always coerce")
	}
	if !c.CanCoerce(INT4OID, INT4OID, 'e') {
		t.Error("same type should always coerce")
	}
}

func TestCanCoerceImplicit(t *testing.T) {
	c := New()
	// int4 → int8 is implicit.
	if !c.CanCoerce(INT4OID, INT8OID, 'i') {
		t.Error("int4 → int8 should coerce implicitly")
	}
	if !c.CanCoerce(INT4OID, INT8OID, 'a') {
		t.Error("int4 → int8 should coerce in assignment context")
	}
	if !c.CanCoerce(INT4OID, INT8OID, 'e') {
		t.Error("int4 → int8 should coerce in explicit context")
	}
}

func TestCanCoerceAssignment(t *testing.T) {
	c := New()
	// int8 → int4 is assignment cast.
	if c.CanCoerce(INT8OID, INT4OID, 'i') {
		t.Error("int8 → int4 should NOT coerce implicitly")
	}
	if !c.CanCoerce(INT8OID, INT4OID, 'a') {
		t.Error("int8 → int4 should coerce in assignment context")
	}
	if !c.CanCoerce(INT8OID, INT4OID, 'e') {
		t.Error("int8 → int4 should coerce in explicit context")
	}
}

func TestCanCoerceExplicit(t *testing.T) {
	c := New()
	// bool → int4 is explicit only.
	if c.CanCoerce(BOOLOID, INT4OID, 'i') {
		t.Error("bool → int4 should NOT coerce implicitly")
	}
	if c.CanCoerce(BOOLOID, INT4OID, 'a') {
		t.Error("bool → int4 should NOT coerce in assignment context")
	}
	if !c.CanCoerce(BOOLOID, INT4OID, 'e') {
		t.Error("bool → int4 should coerce in explicit context")
	}
}

func TestCanCoerceNone(t *testing.T) {
	c := New()
	// bool → float8 has no direct cast.
	if c.CanCoerce(BOOLOID, FLOAT8OID, 'e') {
		t.Error("bool → float8 should not coerce")
	}
}

func TestIsBinaryCoercible(t *testing.T) {
	c := New()
	// varchar → text is binary coercible.
	if !c.IsBinaryCoercible(VARCHAROID, TEXTOID) {
		t.Error("varchar → text should be binary coercible")
	}
	// int4 → int8 is NOT binary coercible (requires function).
	if c.IsBinaryCoercible(INT4OID, INT8OID) {
		t.Error("int4 → int8 should NOT be binary coercible")
	}
	// Same type is always binary coercible.
	if !c.IsBinaryCoercible(INT4OID, INT4OID) {
		t.Error("same type should be binary coercible")
	}
}

func TestFindCoercionPathway(t *testing.T) {
	c := New()

	// Same type → relabel.
	p, foid := c.FindCoercionPathway(INT4OID, INT4OID, 'i')
	if p != CoercionRelabel || foid != 0 {
		t.Errorf("same type: got (%d, %d), want (Relabel, 0)", p, foid)
	}

	// int4 → int8 implicit → function cast.
	p, foid = c.FindCoercionPathway(INT4OID, INT8OID, 'i')
	if p != CoercionFunc {
		t.Errorf("int4→int8: got pathway %d, want CoercionFunc", p)
	}
	if foid == 0 {
		t.Error("int4→int8: expected non-zero function OID")
	}

	// varchar → text → binary coercible.
	p, _ = c.FindCoercionPathway(VARCHAROID, TEXTOID, 'i')
	if p != CoercionRelabel {
		t.Errorf("varchar→text: got pathway %d, want CoercionRelabel", p)
	}

	// bool → int4 in implicit context → none.
	p, _ = c.FindCoercionPathway(BOOLOID, INT4OID, 'i')
	if p != CoercionNone {
		t.Errorf("bool→int4 implicit: got pathway %d, want CoercionNone", p)
	}
}

func TestSelectCommonTypeIdentical(t *testing.T) {
	c := New()
	oid, err := c.selectCommonTypeFromOIDs([]uint32{INT4OID, INT4OID}, false)
	if err != nil {
		t.Fatal(err)
	}
	if oid != INT4OID {
		t.Errorf("got %d, want %d", oid, INT4OID)
	}
}

func TestSelectCommonTypeAllUnknown(t *testing.T) {
	c := New()
	oid, err := c.selectCommonTypeFromOIDs([]uint32{UNKNOWNOID, UNKNOWNOID}, false)
	if err != nil {
		t.Fatal(err)
	}
	if oid != TEXTOID {
		t.Errorf("got %d, want TEXTOID (%d)", oid, TEXTOID)
	}
}

func TestSelectCommonTypeMixedCategories(t *testing.T) {
	c := New()
	// bool (category B) + int4 (category N) → error.
	_, err := c.selectCommonTypeFromOIDs([]uint32{BOOLOID, INT4OID}, false)
	if err == nil {
		t.Error("expected error for mixed categories")
	}
}

func TestSelectCommonTypeNumericPromotion(t *testing.T) {
	c := New()
	// int4 + int8 → both numeric category. int8 can accept int4 implicitly.
	oid, err := c.selectCommonTypeFromOIDs([]uint32{INT4OID, INT8OID}, false)
	if err != nil {
		t.Fatal(err)
	}
	if oid != INT8OID {
		t.Errorf("got %d, want INT8OID (%d)", oid, INT8OID)
	}
}
