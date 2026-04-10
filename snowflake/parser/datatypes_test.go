package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// testParseDataType constructs a Parser from input and calls parseDataType.
func testParseDataType(input string) (*ast.TypeName, error) {
	p := &Parser{lexer: NewLexer(input), input: input}
	p.advance()
	return p.parseDataType()
}

func TestParseDataType_SimpleInt(t *testing.T) {
	for _, kw := range []string{"INT", "INTEGER", "SMALLINT", "TINYINT", "BYTEINT", "BIGINT"} {
		t.Run(kw, func(t *testing.T) {
			dt, err := testParseDataType(kw)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if dt.Kind != ast.TypeInt {
				t.Errorf("Kind = %v, want TypeInt", dt.Kind)
			}
			if strings.ToUpper(dt.Name) != kw {
				t.Errorf("Name = %q, want %q", dt.Name, kw)
			}
			if dt.Params != nil {
				t.Errorf("Params = %v, want nil", dt.Params)
			}
		})
	}
}

func TestParseDataType_NumberWithParams(t *testing.T) {
	cases := []struct {
		input      string
		wantName   string
		wantParams []int
	}{
		{"NUMBER(38, 0)", "NUMBER", []int{38, 0}},
		{"NUMERIC(10)", "NUMERIC", []int{10}},
		{"DECIMAL(18, 2)", "DECIMAL", []int{18, 2}},
		{"NUMBER", "NUMBER", nil},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			dt, err := testParseDataType(c.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if dt.Kind != ast.TypeNumber {
				t.Errorf("Kind = %v, want TypeNumber", dt.Kind)
			}
			if len(dt.Params) != len(c.wantParams) {
				t.Fatalf("Params = %v, want %v", dt.Params, c.wantParams)
			}
			for i, p := range dt.Params {
				if p != c.wantParams[i] {
					t.Errorf("Params[%d] = %d, want %d", i, p, c.wantParams[i])
				}
			}
		})
	}
}

func TestParseDataType_FloatAliases(t *testing.T) {
	for _, kw := range []string{"FLOAT", "FLOAT4", "FLOAT8", "REAL"} {
		t.Run(kw, func(t *testing.T) {
			dt, err := testParseDataType(kw)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if dt.Kind != ast.TypeFloat {
				t.Errorf("Kind = %v, want TypeFloat", dt.Kind)
			}
		})
	}
}

func TestParseDataType_DoublePrecision(t *testing.T) {
	dt, err := testParseDataType("DOUBLE PRECISION")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != ast.TypeFloat {
		t.Errorf("Kind = %v, want TypeFloat", dt.Kind)
	}
	if dt.Name != "DOUBLE PRECISION" {
		t.Errorf("Name = %q, want 'DOUBLE PRECISION'", dt.Name)
	}
}

func TestParseDataType_DoubleAlone(t *testing.T) {
	dt, err := testParseDataType("DOUBLE")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != ast.TypeFloat || dt.Name != "DOUBLE" {
		t.Errorf("got Kind=%v Name=%q, want TypeFloat DOUBLE", dt.Kind, dt.Name)
	}
}

func TestParseDataType_CharVarying(t *testing.T) {
	for _, input := range []string{"CHAR VARYING", "NCHAR VARYING"} {
		t.Run(input, func(t *testing.T) {
			dt, err := testParseDataType(input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if dt.Kind != ast.TypeVarchar {
				t.Errorf("Kind = %v, want TypeVarchar", dt.Kind)
			}
			if dt.Name != input {
				t.Errorf("Name = %q, want %q", dt.Name, input)
			}
		})
	}
}

func TestParseDataType_CharVaryingWithLength(t *testing.T) {
	dt, err := testParseDataType("CHAR VARYING(100)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != ast.TypeVarchar || dt.Name != "CHAR VARYING" {
		t.Errorf("got Kind=%v Name=%q", dt.Kind, dt.Name)
	}
	if len(dt.Params) != 1 || dt.Params[0] != 100 {
		t.Errorf("Params = %v, want [100]", dt.Params)
	}
}

func TestParseDataType_StringTypes(t *testing.T) {
	for _, kw := range []string{"VARCHAR", "NVARCHAR", "NVARCHAR2", "STRING", "TEXT"} {
		t.Run(kw, func(t *testing.T) {
			dt, err := testParseDataType(kw)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if dt.Kind != ast.TypeVarchar {
				t.Errorf("Kind = %v, want TypeVarchar", dt.Kind)
			}
		})
	}
}

func TestParseDataType_VarcharWithLength(t *testing.T) {
	dt, err := testParseDataType("VARCHAR(100)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != ast.TypeVarchar {
		t.Errorf("Kind = %v, want TypeVarchar", dt.Kind)
	}
	if len(dt.Params) != 1 || dt.Params[0] != 100 {
		t.Errorf("Params = %v, want [100]", dt.Params)
	}
}

func TestParseDataType_CharTypes(t *testing.T) {
	for _, kw := range []string{"CHAR", "NCHAR", "CHARACTER"} {
		t.Run(kw, func(t *testing.T) {
			dt, err := testParseDataType(kw)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if dt.Kind != ast.TypeChar {
				t.Errorf("Kind = %v, want TypeChar", dt.Kind)
			}
		})
	}
}

func TestParseDataType_CharWithLength(t *testing.T) {
	dt, err := testParseDataType("CHAR(10)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != ast.TypeChar || len(dt.Params) != 1 || dt.Params[0] != 10 {
		t.Errorf("got Kind=%v Params=%v, want TypeChar [10]", dt.Kind, dt.Params)
	}
}

func TestParseDataType_Binary(t *testing.T) {
	for _, tc := range []struct {
		input string
		kind  ast.TypeKind
	}{
		{"BINARY", ast.TypeBinary},
		{"BINARY(16)", ast.TypeBinary},
		{"VARBINARY", ast.TypeVarbinary},
		{"VARBINARY(256)", ast.TypeVarbinary},
	} {
		t.Run(tc.input, func(t *testing.T) {
			dt, err := testParseDataType(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if dt.Kind != tc.kind {
				t.Errorf("Kind = %v, want %v", dt.Kind, tc.kind)
			}
		})
	}
}

func TestParseDataType_TimestampVariants(t *testing.T) {
	cases := []struct {
		input    string
		wantKind ast.TypeKind
	}{
		{"TIMESTAMP", ast.TypeTimestamp},
		{"TIMESTAMP(9)", ast.TypeTimestamp},
		{"TIMESTAMP_LTZ", ast.TypeTimestampLTZ},
		{"TIMESTAMP_LTZ(9)", ast.TypeTimestampLTZ},
		{"TIMESTAMPLTZ", ast.TypeTimestampLTZ},
		{"TIMESTAMP_NTZ", ast.TypeTimestampNTZ},
		{"TIMESTAMP_NTZ(6)", ast.TypeTimestampNTZ},
		{"TIMESTAMPNTZ", ast.TypeTimestampNTZ},
		{"TIMESTAMP_TZ", ast.TypeTimestampTZ},
		{"TIMESTAMP_TZ(3)", ast.TypeTimestampTZ},
		{"TIMESTAMPTZ", ast.TypeTimestampTZ},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			dt, err := testParseDataType(c.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if dt.Kind != c.wantKind {
				t.Errorf("Kind = %v, want %v", dt.Kind, c.wantKind)
			}
		})
	}
}

func TestParseDataType_TimeWithPrecision(t *testing.T) {
	dt, err := testParseDataType("TIME(9)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != ast.TypeTime || len(dt.Params) != 1 || dt.Params[0] != 9 {
		t.Errorf("got Kind=%v Params=%v, want TypeTime [9]", dt.Kind, dt.Params)
	}
}

func TestParseDataType_DateTime(t *testing.T) {
	dt, err := testParseDataType("DATETIME")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != ast.TypeDateTime {
		t.Errorf("Kind = %v, want TypeDateTime", dt.Kind)
	}
}

func TestParseDataType_Boolean(t *testing.T) {
	dt, err := testParseDataType("BOOLEAN")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != ast.TypeBoolean {
		t.Errorf("Kind = %v, want TypeBoolean", dt.Kind)
	}
}

func TestParseDataType_SemiStructured(t *testing.T) {
	for _, tc := range []struct {
		input string
		kind  ast.TypeKind
	}{
		{"VARIANT", ast.TypeVariant},
		{"OBJECT", ast.TypeObject},
		{"GEOGRAPHY", ast.TypeGeography},
		{"GEOMETRY", ast.TypeGeometry},
	} {
		t.Run(tc.input, func(t *testing.T) {
			dt, err := testParseDataType(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if dt.Kind != tc.kind {
				t.Errorf("Kind = %v, want %v", dt.Kind, tc.kind)
			}
		})
	}
}

func TestParseDataType_ArrayUntyped(t *testing.T) {
	dt, err := testParseDataType("ARRAY")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != ast.TypeArray {
		t.Errorf("Kind = %v, want TypeArray", dt.Kind)
	}
	if dt.ElementType != nil {
		t.Errorf("ElementType = %+v, want nil", dt.ElementType)
	}
}

func TestParseDataType_ArrayTyped(t *testing.T) {
	dt, err := testParseDataType("ARRAY(VARCHAR)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != ast.TypeArray {
		t.Errorf("Kind = %v, want TypeArray", dt.Kind)
	}
	if dt.ElementType == nil {
		t.Fatal("ElementType is nil, want *TypeName")
	}
	if dt.ElementType.Kind != ast.TypeVarchar {
		t.Errorf("ElementType.Kind = %v, want TypeVarchar", dt.ElementType.Kind)
	}
}

func TestParseDataType_ArrayNested(t *testing.T) {
	dt, err := testParseDataType("ARRAY(NUMBER(10, 2))")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != ast.TypeArray || dt.ElementType == nil {
		t.Fatalf("got Kind=%v ElementType=%v", dt.Kind, dt.ElementType)
	}
	if dt.ElementType.Kind != ast.TypeNumber {
		t.Errorf("ElementType.Kind = %v, want TypeNumber", dt.ElementType.Kind)
	}
	if len(dt.ElementType.Params) != 2 || dt.ElementType.Params[0] != 10 || dt.ElementType.Params[1] != 2 {
		t.Errorf("ElementType.Params = %v, want [10, 2]", dt.ElementType.Params)
	}
}

func TestParseDataType_Vector(t *testing.T) {
	cases := []struct {
		input    string
		elemKind ast.TypeKind
		dim      int
	}{
		{"VECTOR(INT, 256)", ast.TypeInt, 256},
		{"VECTOR(FLOAT, 768)", ast.TypeFloat, 768},
		{"VECTOR(FLOAT4, 128)", ast.TypeFloat, 128},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			dt, err := testParseDataType(c.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if dt.Kind != ast.TypeVector {
				t.Errorf("Kind = %v, want TypeVector", dt.Kind)
			}
			if dt.ElementType == nil || dt.ElementType.Kind != c.elemKind {
				t.Errorf("ElementType = %+v, want Kind=%v", dt.ElementType, c.elemKind)
			}
			if dt.VectorDim != c.dim {
				t.Errorf("VectorDim = %d, want %d", dt.VectorDim, c.dim)
			}
		})
	}
}

func TestParseDataType_Errors(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"unknown keyword", "FOOBAR"},
		{"missing close paren", "NUMBER(38"},
		{"vector missing dim", "VECTOR(INT)"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := testParseDataType(c.input)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

func TestParseDataType_Freestanding(t *testing.T) {
	dt, errs := ParseDataType("VARCHAR(100)")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %+v", errs)
	}
	if dt.Kind != ast.TypeVarchar || len(dt.Params) != 1 || dt.Params[0] != 100 {
		t.Errorf("got Kind=%v Params=%v", dt.Kind, dt.Params)
	}
}

func TestParseDataType_FreestandingTrailing(t *testing.T) {
	dt, errs := ParseDataType("INT foo")
	if dt == nil {
		t.Fatal("expected non-nil TypeName")
	}
	if len(errs) == 0 || !strings.Contains(errs[0].Msg, "unexpected token") {
		t.Errorf("expected 'unexpected token' error, got %+v", errs)
	}
}

func TestParseDataType_LocSpanning(t *testing.T) {
	// "DOUBLE PRECISION" — Loc should span both words.
	dt, err := testParseDataType("DOUBLE PRECISION")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// "DOUBLE" = 6 chars at 0..5, " " at 6, "PRECISION" at 7..15.
	// Loc should be {0, 16}.
	if dt.Loc.Start != 0 || dt.Loc.End != 16 {
		t.Errorf("Loc = %+v, want {0, 16}", dt.Loc)
	}
}

func TestParseDataType_LocWithParams(t *testing.T) {
	// "VARCHAR(100)" — Loc should span from V to closing ).
	dt, err := testParseDataType("VARCHAR(100)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Loc.Start != 0 || dt.Loc.End != 12 {
		t.Errorf("Loc = %+v, want {0, 12}", dt.Loc)
	}
}
