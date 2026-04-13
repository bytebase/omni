package parser

import (
	"testing"

	"github.com/bytebase/omni/doris/ast"
)

// parseType is a test helper that creates a parser for the given input string
// and calls parseDataType. It returns the parsed TypeName and any error.
func parseType(input string) (*ast.TypeName, error) {
	p := makeParser(input)
	return p.parseDataType()
}

// assertTypeName checks that tn.Name equals want.
func assertTypeName(t *testing.T, tn *ast.TypeName, want string) {
	t.Helper()
	if tn == nil {
		t.Fatalf("TypeName is nil, want Name=%q", want)
	}
	if tn.Name != want {
		t.Errorf("TypeName.Name = %q, want %q", tn.Name, want)
	}
}

// assertParams checks that tn.Params matches want.
func assertParams(t *testing.T, tn *ast.TypeName, want ...int) {
	t.Helper()
	if len(tn.Params) != len(want) {
		t.Errorf("Params = %v, want %v", tn.Params, want)
		return
	}
	for i, w := range want {
		if tn.Params[i] != w {
			t.Errorf("Params[%d] = %d, want %d", i, tn.Params[i], w)
		}
	}
}

// ---------------------------------------------------------------------------
// Simple primitive types (no parameters)
// ---------------------------------------------------------------------------

func TestDataType_INT(t *testing.T) {
	tn, err := parseType("INT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "INT")
	if len(tn.Params) != 0 {
		t.Errorf("expected no params, got %v", tn.Params)
	}
}

func TestDataType_INTEGER(t *testing.T) {
	tn, err := parseType("INTEGER")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "INTEGER")
}

func TestDataType_TINYINT(t *testing.T) {
	tn, err := parseType("TINYINT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "TINYINT")
}

func TestDataType_SMALLINT(t *testing.T) {
	tn, err := parseType("SMALLINT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "SMALLINT")
}

func TestDataType_BIGINT(t *testing.T) {
	tn, err := parseType("BIGINT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "BIGINT")
}

func TestDataType_FLOAT(t *testing.T) {
	tn, err := parseType("FLOAT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "FLOAT")
}

func TestDataType_DOUBLE(t *testing.T) {
	tn, err := parseType("DOUBLE")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "DOUBLE")
}

func TestDataType_BOOLEAN(t *testing.T) {
	tn, err := parseType("BOOLEAN")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "BOOLEAN")
}

func TestDataType_DATE(t *testing.T) {
	tn, err := parseType("DATE")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "DATE")
}

func TestDataType_DATETIME(t *testing.T) {
	tn, err := parseType("DATETIME")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "DATETIME")
}

func TestDataType_TIME(t *testing.T) {
	tn, err := parseType("TIME")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "TIME")
}

func TestDataType_STRING(t *testing.T) {
	tn, err := parseType("STRING")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "STRING")
}

func TestDataType_TEXT(t *testing.T) {
	tn, err := parseType("TEXT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "TEXT")
}

// ---------------------------------------------------------------------------
// Legacy Doris versioned type names
// ---------------------------------------------------------------------------

func TestDataType_LARGEINT(t *testing.T) {
	tn, err := parseType("LARGEINT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "LARGEINT")
}

func TestDataType_DATEV1(t *testing.T) {
	tn, err := parseType("DATEV1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "DATEV1")
}

func TestDataType_DATEV2(t *testing.T) {
	tn, err := parseType("DATEV2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "DATEV2")
}

func TestDataType_DATETIMEV1(t *testing.T) {
	tn, err := parseType("DATETIMEV1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "DATETIMEV1")
}

func TestDataType_DATETIMEV2(t *testing.T) {
	tn, err := parseType("DATETIMEV2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "DATETIMEV2")
}

func TestDataType_DECIMALV2(t *testing.T) {
	tn, err := parseType("DECIMALV2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "DECIMALV2")
}

func TestDataType_DECIMALV3(t *testing.T) {
	tn, err := parseType("DECIMALV3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "DECIMALV3")
}

// ---------------------------------------------------------------------------
// Special Doris types
// ---------------------------------------------------------------------------

func TestDataType_HLL(t *testing.T) {
	tn, err := parseType("HLL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "HLL")
}

func TestDataType_BITMAP(t *testing.T) {
	tn, err := parseType("BITMAP")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "BITMAP")
}

func TestDataType_QUANTILE_STATE(t *testing.T) {
	tn, err := parseType("QUANTILE_STATE")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "QUANTILE_STATE")
}

func TestDataType_JSON(t *testing.T) {
	tn, err := parseType("JSON")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "JSON")
}

func TestDataType_JSONB(t *testing.T) {
	tn, err := parseType("JSONB")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "JSONB")
}

func TestDataType_AGG_STATE(t *testing.T) {
	tn, err := parseType("AGG_STATE")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "AGG_STATE")
}

func TestDataType_IPV4(t *testing.T) {
	tn, err := parseType("IPV4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "IPV4")
}

func TestDataType_IPV6(t *testing.T) {
	tn, err := parseType("IPV6")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "IPV6")
}

func TestDataType_VARIANT(t *testing.T) {
	tn, err := parseType("VARIANT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "VARIANT")
}

func TestDataType_VARBINARY(t *testing.T) {
	tn, err := parseType("VARBINARY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "VARBINARY")
}

// ---------------------------------------------------------------------------
// Parameterized types
// ---------------------------------------------------------------------------

func TestDataType_VARCHAR_Param(t *testing.T) {
	tn, err := parseType("VARCHAR(255)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "VARCHAR")
	assertParams(t, tn, 255)
}

func TestDataType_CHAR_Param(t *testing.T) {
	tn, err := parseType("CHAR(1)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "CHAR")
	assertParams(t, tn, 1)
}

func TestDataType_DECIMAL_TwoParams(t *testing.T) {
	tn, err := parseType("DECIMAL(10,2)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "DECIMAL")
	assertParams(t, tn, 10, 2)
}

func TestDataType_DECIMAL_OneParam(t *testing.T) {
	tn, err := parseType("DECIMAL(18)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "DECIMAL")
	assertParams(t, tn, 18)
}

func TestDataType_DECIMALV3_TwoParams(t *testing.T) {
	tn, err := parseType("DECIMALV3(27,9)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "DECIMALV3")
	assertParams(t, tn, 27, 9)
}

func TestDataType_DATETIMEV2_Param(t *testing.T) {
	tn, err := parseType("DATETIMEV2(3)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "DATETIMEV2")
	assertParams(t, tn, 3)
}

// ---------------------------------------------------------------------------
// ARRAY types
// ---------------------------------------------------------------------------

func TestDataType_ARRAY_INT(t *testing.T) {
	tn, err := parseType("ARRAY<INT>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "ARRAY")
	if tn.ElementType == nil {
		t.Fatal("ElementType is nil")
	}
	assertTypeName(t, tn.ElementType, "INT")
}

func TestDataType_ARRAY_VARCHAR(t *testing.T) {
	tn, err := parseType("ARRAY<VARCHAR(255)>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "ARRAY")
	if tn.ElementType == nil {
		t.Fatal("ElementType is nil")
	}
	assertTypeName(t, tn.ElementType, "VARCHAR")
	assertParams(t, tn.ElementType, 255)
}

func TestDataType_ARRAY_Nested(t *testing.T) {
	tn, err := parseType("ARRAY<ARRAY<INT>>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "ARRAY")
	if tn.ElementType == nil {
		t.Fatal("outer ElementType is nil")
	}
	assertTypeName(t, tn.ElementType, "ARRAY")
	if tn.ElementType.ElementType == nil {
		t.Fatal("inner ElementType is nil")
	}
	assertTypeName(t, tn.ElementType.ElementType, "INT")
}

func TestDataType_ARRAY_BIGINT(t *testing.T) {
	tn, err := parseType("ARRAY<BIGINT>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "ARRAY")
	assertTypeName(t, tn.ElementType, "BIGINT")
}

// ---------------------------------------------------------------------------
// MAP types
// ---------------------------------------------------------------------------

func TestDataType_MAP_StringInt(t *testing.T) {
	tn, err := parseType("MAP<STRING, INT>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "MAP")
	if tn.ElementType == nil {
		t.Fatal("key type (ElementType) is nil")
	}
	assertTypeName(t, tn.ElementType, "STRING")
	if tn.ValueType == nil {
		t.Fatal("value type (ValueType) is nil")
	}
	assertTypeName(t, tn.ValueType, "INT")
}

func TestDataType_MAP_VarcharArray(t *testing.T) {
	tn, err := parseType("MAP<VARCHAR(10), ARRAY<INT>>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "MAP")
	if tn.ElementType == nil {
		t.Fatal("key type is nil")
	}
	assertTypeName(t, tn.ElementType, "VARCHAR")
	assertParams(t, tn.ElementType, 10)
	if tn.ValueType == nil {
		t.Fatal("value type is nil")
	}
	assertTypeName(t, tn.ValueType, "ARRAY")
	if tn.ValueType.ElementType == nil {
		t.Fatal("value array element type is nil")
	}
	assertTypeName(t, tn.ValueType.ElementType, "INT")
}

func TestDataType_MAP_NoSpaces(t *testing.T) {
	tn, err := parseType("MAP<STRING,INT>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "MAP")
	assertTypeName(t, tn.ElementType, "STRING")
	assertTypeName(t, tn.ValueType, "INT")
}

// ---------------------------------------------------------------------------
// STRUCT types
// ---------------------------------------------------------------------------

func TestDataType_STRUCT_SingleField(t *testing.T) {
	tn, err := parseType("STRUCT<name VARCHAR(50)>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "STRUCT")
	if len(tn.Fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(tn.Fields))
	}
	if tn.Fields[0].Name != "name" {
		t.Errorf("Fields[0].Name = %q, want %q", tn.Fields[0].Name, "name")
	}
	assertTypeName(t, tn.Fields[0].Type, "VARCHAR")
	assertParams(t, tn.Fields[0].Type, 50)
}

func TestDataType_STRUCT_TwoFields(t *testing.T) {
	tn, err := parseType("STRUCT<name VARCHAR(50), age INT>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "STRUCT")
	if len(tn.Fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(tn.Fields))
	}
	if tn.Fields[0].Name != "name" {
		t.Errorf("Fields[0].Name = %q, want %q", tn.Fields[0].Name, "name")
	}
	assertTypeName(t, tn.Fields[0].Type, "VARCHAR")
	assertParams(t, tn.Fields[0].Type, 50)
	if tn.Fields[1].Name != "age" {
		t.Errorf("Fields[1].Name = %q, want %q", tn.Fields[1].Name, "age")
	}
	assertTypeName(t, tn.Fields[1].Type, "INT")
}

func TestDataType_STRUCT_ThreeFields(t *testing.T) {
	tn, err := parseType("STRUCT<id BIGINT, label VARCHAR(100), score DOUBLE>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "STRUCT")
	if len(tn.Fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(tn.Fields))
	}
	wantNames := []string{"id", "label", "score"}
	wantTypes := []string{"BIGINT", "VARCHAR", "DOUBLE"}
	for i := range wantNames {
		if tn.Fields[i].Name != wantNames[i] {
			t.Errorf("Fields[%d].Name = %q, want %q", i, tn.Fields[i].Name, wantNames[i])
		}
		assertTypeName(t, tn.Fields[i].Type, wantTypes[i])
	}
}

func TestDataType_STRUCT_NestedArray(t *testing.T) {
	tn, err := parseType("STRUCT<tags ARRAY<STRING>, metadata MAP<STRING, INT>>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "STRUCT")
	if len(tn.Fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(tn.Fields))
	}
	assertTypeName(t, tn.Fields[0].Type, "ARRAY")
	assertTypeName(t, tn.Fields[1].Type, "MAP")
}

// ---------------------------------------------------------------------------
// Source location tracking
// ---------------------------------------------------------------------------

func TestDataType_Loc_Simple(t *testing.T) {
	tn, err := parseType("INT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tn.Loc.Start != 0 {
		t.Errorf("Loc.Start = %d, want 0", tn.Loc.Start)
	}
	if tn.Loc.End != 3 {
		t.Errorf("Loc.End = %d, want 3", tn.Loc.End)
	}
}

func TestDataType_Loc_Parameterized(t *testing.T) {
	tn, err := parseType("VARCHAR(255)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tn.Loc.Start != 0 {
		t.Errorf("Loc.Start = %d, want 0", tn.Loc.Start)
	}
	// "VARCHAR(255)" is 12 chars → End = 12
	if tn.Loc.End != 12 {
		t.Errorf("Loc.End = %d, want 12", tn.Loc.End)
	}
}

// ---------------------------------------------------------------------------
// NodeTag
// ---------------------------------------------------------------------------

func TestDataType_NodeTag(t *testing.T) {
	tn, err := parseType("INT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tn.Tag().String() != "TypeName" {
		t.Errorf("Tag().String() = %q, want %q", tn.Tag().String(), "TypeName")
	}
}

// ---------------------------------------------------------------------------
// Error cases
// ---------------------------------------------------------------------------

func TestDataType_Error_EmptyInput(t *testing.T) {
	_, err := parseType("")
	if err == nil {
		t.Fatal("expected error for empty input, got nil")
	}
}

func TestDataType_Error_NotAType(t *testing.T) {
	// SELECT is a reserved keyword but not a type name
	_, err := parseType("SELECT")
	if err == nil {
		t.Fatal("expected error for SELECT as type, got nil")
	}
}

func TestDataType_Error_ARRAY_MissingAngle(t *testing.T) {
	_, err := parseType("ARRAY INT")
	if err == nil {
		t.Fatal("expected error for ARRAY without '<'")
	}
}

func TestDataType_Error_MAP_MissingComma(t *testing.T) {
	_, err := parseType("MAP<STRING INT>")
	if err == nil {
		t.Fatal("expected error for MAP without ','")
	}
}

// ---------------------------------------------------------------------------
// Lowercase / case-insensitive keyword input
// ---------------------------------------------------------------------------

func TestDataType_Lowercase_int(t *testing.T) {
	tn, err := parseType("int")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Canonical name should be uppercase regardless of input case
	assertTypeName(t, tn, "INT")
}

func TestDataType_Lowercase_varchar(t *testing.T) {
	tn, err := parseType("varchar(100)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "VARCHAR")
	assertParams(t, tn, 100)
}

func TestDataType_Lowercase_array(t *testing.T) {
	tn, err := parseType("array<int>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertTypeName(t, tn, "ARRAY")
	assertTypeName(t, tn.ElementType, "INT")
}
