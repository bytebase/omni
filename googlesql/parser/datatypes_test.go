package parser

import (
	"strings"
	"testing"
)

// testParseType constructs a Parser over input and calls parseType, asserting
// the whole input is consumed (cur at EOF). It mirrors snowflake's
// testParseDataType helper. Returns the parsed *DataType and any error.
func testParseType(input string) (*DataType, error) {
	p := &Parser{lexer: NewLexer(input), input: input}
	p.advance()
	dt, err := p.parseType()
	if err != nil {
		return nil, err
	}
	if p.cur.Type != tokEOF {
		return dt, &ParseError{Loc: p.cur.Loc, Msg: "trailing tokens after type: " + TokenName(p.cur.Type)}
	}
	return dt, nil
}

// ---------------------------------------------------------------------------
// type_name (path_expression | INTERVAL) — scalars are plain identifiers.
// ---------------------------------------------------------------------------

func TestParseType_ScalarIdentifierNames(t *testing.T) {
	// Oracle (Spanner emulator): every one of these parses (accept) as a
	// type_name path_expression — INT64/STRING/BOOL/FLOAT64/FLOAT32/BYTES/
	// GEOGRAPHY/TOKENLIST are NOT keyword tokens, they are identifiers; the
	// "Type not found" cases (INT, BIGNUMERIC, GEOGRAPHY) are semantic, the
	// grammar still parsed them.
	names := []string{
		"INT64", "INT", "SMALLINT", "INTEGER", "BIGINT", "TINYINT", "BYTEINT",
		"FLOAT64", "FLOAT32", "BOOL", "BOOLEAN", "STRING", "BYTES",
		"GEOGRAPHY", "TOKENLIST", "BIGNUMERIC", "BIGDECIMAL",
	}
	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			dt, err := testParseType(name)
			if err != nil {
				t.Fatalf("parseType(%q): unexpected error: %v", name, err)
			}
			if dt.Kind != TypeName {
				t.Errorf("Kind = %v, want TypeName", dt.Kind)
			}
			if got := strings.ToUpper(dt.String()); got != name {
				t.Errorf("String() = %q, want %q", got, name)
			}
			if len(dt.NameParts) != 1 || strings.ToUpper(dt.NameParts[0]) != name {
				t.Errorf("NameParts = %v, want [%q]", dt.NameParts, name)
			}
		})
	}
}

func TestParseType_KeywordScalarNames(t *testing.T) {
	// These type names ARE keyword tokens but are non-reserved, so they parse
	// as the first component of a type_name path_expression (oracle: all accept,
	// some "Type not found" semantically on Spanner — still parsed).
	names := []string{"NUMERIC", "DECIMAL", "JSON", "DATE", "DATETIME", "TIMESTAMP", "TIME", "MAP", "FUNCTION"}
	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			dt, err := testParseType(name)
			if err != nil {
				t.Fatalf("parseType(%q): unexpected error: %v", name, err)
			}
			if dt.Kind != TypeName {
				t.Errorf("Kind = %v, want TypeName", dt.Kind)
			}
			if got := strings.ToUpper(dt.String()); got != name {
				t.Errorf("String() = %q, want %q", got, name)
			}
		})
	}
}

func TestParseType_Interval(t *testing.T) {
	// type_name's INTERVAL alternative. Oracle: `CAST(NULL AS INTERVAL)` accepts.
	dt, err := testParseType("INTERVAL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != TypeName {
		t.Fatalf("Kind = %v, want TypeName", dt.Kind)
	}
	if !dt.IsInterval {
		t.Errorf("IsInterval = false, want true")
	}
	if dt.String() != "INTERVAL" {
		t.Errorf("String() = %q, want INTERVAL", dt.String())
	}
}

func TestParseType_IntervalNoPathContinuation(t *testing.T) {
	// Oracle: `INTERVAL.foo` REJECTS ("Expected ... but got ."). INTERVAL is
	// its own type_name alt with no path continuation. testParseType requires
	// EOF, so the trailing `.foo` must surface as an error.
	if _, err := testParseType("INTERVAL.foo"); err == nil {
		t.Fatal("parseType(\"INTERVAL.foo\"): want error, got nil")
	}
}

func TestParseType_DottedPathName(t *testing.T) {
	// Oracle: `CAST(NULL AS foo.bar.Baz)` and `NUMERIC.foo` both accept — a
	// type_name path_expression may be multi-component (proto/enum type names).
	cases := []struct {
		input string
		parts []string
	}{
		{"foo.bar.Baz", []string{"foo", "bar", "Baz"}},
		{"NUMERIC.foo", []string{"NUMERIC", "foo"}},
		{"a.b", []string{"a", "b"}},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			dt, err := testParseType(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if dt.Kind != TypeName {
				t.Fatalf("Kind = %v, want TypeName", dt.Kind)
			}
			if len(dt.NameParts) != len(tc.parts) {
				t.Fatalf("NameParts = %v, want %v", dt.NameParts, tc.parts)
			}
			for i := range tc.parts {
				if dt.NameParts[i] != tc.parts[i] {
					t.Errorf("NameParts[%d] = %q, want %q", i, dt.NameParts[i], tc.parts[i])
				}
			}
		})
	}
}

func TestParseType_KeywordPartCasePreserved(t *testing.T) {
	// Regression for identifierText folding keyword-spelled path parts to their
	// canonical UPPER-case keyword name. `Type`, `value`, `Key`, `Value` are all
	// NON-reserved keywords, so the lexer emits a kw* token for them — but it
	// also records the verbatim source spelling in Token.Str. identifierText must
	// return that spelling, not TokenName(tok.Type): otherwise `pkg.Type` parses
	// as `pkg.TYPE`, `x.value` as `x.VALUE`, silently corrupting proto/enum type
	// names (and, since identifierText backs every name part in the grammar,
	// every keyword-spelled identifier elsewhere too).
	cases := []struct {
		input string
		parts []string
	}{
		{"pkg.Type", []string{"pkg", "Type"}},
		{"x.value", []string{"x", "value"}},
		{"p.Key.Value", []string{"p", "Key", "Value"}},
		// Mixed case on the part itself must survive byte-for-byte.
		{"a.TyPe", []string{"a", "TyPe"}},
		// A keyword as the FIRST part (parseTypeName accepts a non-reserved
		// keyword there via isIdentifierStart) must also keep its casing.
		{"Value.x", []string{"Value", "x"}},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			dt, err := testParseType(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if dt.Kind != TypeName {
				t.Fatalf("Kind = %v, want TypeName", dt.Kind)
			}
			if len(dt.NameParts) != len(tc.parts) {
				t.Fatalf("NameParts = %v, want %v", dt.NameParts, tc.parts)
			}
			for i := range tc.parts {
				if dt.NameParts[i] != tc.parts[i] {
					t.Errorf("NameParts[%d] = %q, want %q (source casing must be preserved for keyword-spelled parts)", i, dt.NameParts[i], tc.parts[i])
				}
			}
		})
	}
}

func TestParseType_EmptyBacktickRejected(t *testing.T) {
	// Oracle: `CAST(NULL AS ``)` REJECTS ("Syntax error: Invalid empty
	// identifier"). The lexer admits the empty backtick pair `` as a
	// tokIdentifier with empty Str and NO lex error (a lexer gap, out of this
	// node's scope), so the type parser must reject the empty identifier rather
	// than fabricate a name from the token kind.
	if _, errs := ParseDataType("``"); len(errs) == 0 {
		t.Error("ParseDataType(\"``\"): want an error, got none")
	} else if !strings.Contains(errs[0].Msg, "empty identifier") {
		t.Errorf("error = %q, want it to mention 'empty identifier'", errs[0].Msg)
	}
	// Also as a dotted continuation and a struct field name.
	if _, errs := ParseDataType("foo.``"); len(errs) == 0 {
		t.Error("ParseDataType(\"foo.``\"): want an error, got none")
	}
}

func TestParseType_BacktickQuotedName(t *testing.T) {
	// Oracle: `CAST(NULL AS `my.proto.Type`)` accepts. The lexer strips the
	// backticks and emits a single tokIdentifier whose Str contains the dots;
	// it is ONE path component, not three.
	dt, err := testParseType("`my.proto.Type`")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != TypeName {
		t.Fatalf("Kind = %v, want TypeName", dt.Kind)
	}
	if len(dt.NameParts) != 1 || dt.NameParts[0] != "my.proto.Type" {
		t.Errorf("NameParts = %v, want [my.proto.Type]", dt.NameParts)
	}
}

// ---------------------------------------------------------------------------
// ARRAY<type>
// ---------------------------------------------------------------------------

func TestParseType_Array(t *testing.T) {
	dt, err := testParseType("ARRAY<INT64>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != TypeArray {
		t.Fatalf("Kind = %v, want TypeArray", dt.Kind)
	}
	if dt.ElementType == nil || dt.ElementType.Kind != TypeName {
		t.Fatalf("ElementType = %+v, want a TypeName INT64", dt.ElementType)
	}
	if got := strings.ToUpper(dt.ElementType.String()); got != "INT64" {
		t.Errorf("element = %q, want INT64", got)
	}
}

func TestParseType_ArrayWhitespaceInsensitive(t *testing.T) {
	// Oracle: `ARRAY< INT64 >`, `ARRAY<INT64 >`, `ARRAY< INT64>` all accept.
	for _, in := range []string{"ARRAY< INT64 >", "ARRAY<INT64 >", "ARRAY< INT64>"} {
		t.Run(in, func(t *testing.T) {
			if _, err := testParseType(in); err != nil {
				t.Errorf("parseType(%q): unexpected error: %v", in, err)
			}
		})
	}
}

func TestParseType_ArrayRequiresElement(t *testing.T) {
	// Oracle: `ARRAY<>` REJECTS ("Expected < but got <>") — the adjacent `<>`
	// lexes as one token (tokNotEqual2); an array needs an element type.
	if _, err := testParseType("ARRAY<>"); err == nil {
		t.Error("parseType(\"ARRAY<>\"): want error, got nil")
	}
	// `ARRAY< >` (spaced, two tokens) is also empty and rejects.
	if _, err := testParseType("ARRAY< >"); err == nil {
		t.Error("parseType(\"ARRAY< >\"): want error, got nil")
	}
}

func TestParseType_ArrayRejectsMultipleElements(t *testing.T) {
	// Oracle: `ARRAY<INT64, STRING>` REJECTS ("Expected > but got ,").
	if _, err := testParseType("ARRAY<INT64, STRING>"); err == nil {
		t.Error("parseType(\"ARRAY<INT64, STRING>\"): want error, got nil")
	}
}

func TestParseType_BareArrayRejected(t *testing.T) {
	// Oracle: bare `ARRAY` (no `<`) REJECTS ("Expected < but got )"). ARRAY is a
	// reserved keyword and can ONLY introduce array_type; it is not a valid
	// bare type_name. Same for ARRAY.foo (oracle reject "Expected < but got .").
	for _, in := range []string{"ARRAY", "ARRAY.foo"} {
		t.Run(in, func(t *testing.T) {
			if _, err := testParseType(in); err == nil {
				t.Errorf("parseType(%q): want error, got nil", in)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Nested types — the >> / >>> split (D1 divergence).
// ---------------------------------------------------------------------------

func TestParseType_NestedArray(t *testing.T) {
	// Oracle: `ARRAY<ARRAY<INT64>>` ACCEPTS (semantic "Arrays of arrays are not
	// supported" — the grammar parsed it). The adjacent `>>` lexes as
	// tokShiftRight and MUST be split into two template closers.
	dt, err := testParseType("ARRAY<ARRAY<INT64>>")
	if err != nil {
		t.Fatalf("parseType(\"ARRAY<ARRAY<INT64>>\"): unexpected error: %v", err)
	}
	if dt.Kind != TypeArray {
		t.Fatalf("outer Kind = %v, want TypeArray", dt.Kind)
	}
	if dt.ElementType == nil || dt.ElementType.Kind != TypeArray {
		t.Fatalf("inner = %+v, want a TypeArray", dt.ElementType)
	}
	inner := dt.ElementType.ElementType
	if inner == nil || strings.ToUpper(inner.String()) != "INT64" {
		t.Errorf("innermost = %+v, want INT64", inner)
	}
}

func TestParseType_TripleNestedClose(t *testing.T) {
	// Oracle: `STRUCT<x ARRAY<STRUCT<y INT64>>>` ACCEPTS — the `>>>` (lexed as
	// tokShiftRight then '>') must split into THREE template closers.
	dt, err := testParseType("STRUCT<x ARRAY<STRUCT<y INT64>>>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != TypeStruct {
		t.Fatalf("Kind = %v, want TypeStruct", dt.Kind)
	}
	if len(dt.Fields) != 1 {
		t.Fatalf("got %d fields, want 1", len(dt.Fields))
	}
	if dt.Fields[0].Type == nil || dt.Fields[0].Type.Kind != TypeArray {
		t.Errorf("field type = %+v, want TypeArray", dt.Fields[0].Type)
	}
}

func TestParseType_StructWithArrayDoubleClose(t *testing.T) {
	// Oracle: `STRUCT<ARRAY<INT64>>` ACCEPTS — anonymous struct field of array
	// type, closing `>>`.
	dt, err := testParseType("STRUCT<ARRAY<INT64>>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != TypeStruct || len(dt.Fields) != 1 {
		t.Fatalf("got %+v, want a 1-field struct", dt)
	}
	if dt.Fields[0].Name != "" {
		t.Errorf("field name = %q, want anonymous", dt.Fields[0].Name)
	}
	if dt.Fields[0].Type.Kind != TypeArray {
		t.Errorf("field type kind = %v, want TypeArray", dt.Fields[0].Type.Kind)
	}
}

func TestParseType_ExtraCloseRejected(t *testing.T) {
	// Oracle: `ARRAY<INT64>>` REJECTS ("Expected ) ... but got >"). After the
	// type is complete, a leftover `>` is a syntax error — the >>-split must NOT
	// borrow a closer when no template is open. testParseType requires EOF.
	if _, err := testParseType("ARRAY<INT64>>"); err == nil {
		t.Error("parseType(\"ARRAY<INT64>>\"): want error, got nil")
	}
}

// ---------------------------------------------------------------------------
// STRUCT<...>
// ---------------------------------------------------------------------------

func TestParseType_StructNamedFields(t *testing.T) {
	dt, err := testParseType("STRUCT<x INT64, y STRING>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != TypeStruct {
		t.Fatalf("Kind = %v, want TypeStruct", dt.Kind)
	}
	if len(dt.Fields) != 2 {
		t.Fatalf("got %d fields, want 2", len(dt.Fields))
	}
	if dt.Fields[0].Name != "x" || strings.ToUpper(dt.Fields[0].Type.String()) != "INT64" {
		t.Errorf("field0 = %+v, want {x INT64}", dt.Fields[0])
	}
	if dt.Fields[1].Name != "y" || strings.ToUpper(dt.Fields[1].Type.String()) != "STRING" {
		t.Errorf("field1 = %+v, want {y STRING}", dt.Fields[1])
	}
}

func TestParseType_StructAnonymousFields(t *testing.T) {
	// Oracle: `STRUCT<INT64, STRING>` accepts (anonymous fields).
	dt, err := testParseType("STRUCT<INT64, STRING>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(dt.Fields) != 2 {
		t.Fatalf("got %d fields, want 2", len(dt.Fields))
	}
	for i, f := range dt.Fields {
		if f.Name != "" {
			t.Errorf("field%d name = %q, want anonymous", i, f.Name)
		}
	}
}

func TestParseType_StructEmpty(t *testing.T) {
	// Oracle: `STRUCT<>` and `STRUCT< >` both ACCEPT (empty struct). The
	// adjacent `<>` lexes as tokNotEqual2; STRUCT treats it as an empty template.
	for _, in := range []string{"STRUCT<>", "STRUCT< >"} {
		t.Run(in, func(t *testing.T) {
			dt, err := testParseType(in)
			if err != nil {
				t.Fatalf("parseType(%q): unexpected error: %v", in, err)
			}
			if dt.Kind != TypeStruct {
				t.Fatalf("Kind = %v, want TypeStruct", dt.Kind)
			}
			if len(dt.Fields) != 0 {
				t.Errorf("got %d fields, want 0", len(dt.Fields))
			}
		})
	}
}

func TestParseType_StructNested(t *testing.T) {
	dt, err := testParseType("STRUCT<x STRUCT<y INT64, z INT64>>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != TypeStruct || len(dt.Fields) != 1 {
		t.Fatalf("got %+v, want 1-field struct", dt)
	}
	inner := dt.Fields[0].Type
	if inner == nil || inner.Kind != TypeStruct || len(inner.Fields) != 2 {
		t.Errorf("inner = %+v, want a 2-field struct", inner)
	}
}

func TestParseType_NamedFieldOfNestedType(t *testing.T) {
	// Regression for the named-field disambiguation + >>-split interaction: a
	// named field whose type is itself a template type closing with `>>` or
	// `>>>`. The named-form lookahead buffers the field type's leading keyword
	// via peekNext, then the inner template-close split pushes a `>` back —
	// these must compose. All oracle-accept.
	cases := []string{
		"STRUCT<a STRUCT<b INT64>>",          // named struct field of struct type, >> close
		"STRUCT<STRUCT<x INT64>>",            // anonymous struct field of struct type, >> close
		"STRUCT<x ARRAY<ARRAY<INT64>>>",      // named field, array-of-array, >>> close
		"ARRAY<STRUCT<a INT64 COLLATE 'x'>>", // collate inside a nested struct field, >> close
		"ARRAY<NUMERIC(10, 2)>",              // parameterized element type, then `)>` close
	}
	for _, in := range cases {
		t.Run(in, func(t *testing.T) {
			if _, errs := ParseDataType(in); len(errs) != 0 {
				t.Errorf("ParseDataType(%q): unexpected errors: %v", in, errs)
			}
		})
	}
}

func TestParseType_HexAndZeroTypeParam(t *testing.T) {
	// Oracle: STRING(0) and STRING(0x10) both accept (an integer_literal
	// type_parameter may be 0 or hex; the lexer parses 0x-hex into tokInteger).
	for _, in := range []string{"STRING(0)", "STRING(0x10)"} {
		t.Run(in, func(t *testing.T) {
			dt, err := testParseType(in)
			if err != nil {
				t.Fatalf("parseType(%q): unexpected error: %v", in, err)
			}
			if len(dt.Params) != 1 || !dt.Params[0].IsInt {
				t.Errorf("Params = %+v, want one int param", dt.Params)
			}
		})
	}
}

func TestParseType_StructFieldArray(t *testing.T) {
	// Oracle: `STRUCT<a ARRAY<INT64>>` accepts.
	dt, err := testParseType("STRUCT<a ARRAY<INT64>>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != TypeStruct || len(dt.Fields) != 1 {
		t.Fatalf("got %+v, want 1-field struct", dt)
	}
	if dt.Fields[0].Name != "a" || dt.Fields[0].Type.Kind != TypeArray {
		t.Errorf("field0 = %+v, want {a ARRAY<...>}", dt.Fields[0])
	}
}

func TestParseType_StructFieldCollate(t *testing.T) {
	// Oracle: `STRUCT<x INT64 COLLATE 'und:ci'>` and
	// `STRUCT<INT64 COLLATE 'x', y STRING>` both accept (collate on a struct
	// field type, named and anonymous).
	for _, in := range []string{"STRUCT<x INT64 COLLATE 'und:ci'>", "STRUCT<INT64 COLLATE 'x', y STRING>"} {
		t.Run(in, func(t *testing.T) {
			if _, err := testParseType(in); err != nil {
				t.Errorf("parseType(%q): unexpected error: %v", in, err)
			}
		})
	}
}

func TestParseType_StructKeywordFieldName(t *testing.T) {
	// Oracle: STRUCT<key INT64> and STRUCT<row INT64> accept — a NON-reserved
	// word-keyword (KEY, ROW, VALUE, DATA) is a valid struct field name
	// (struct_field: identifier type; identifier admits keyword_as_identifier).
	cases := []struct {
		input string
		name  string
	}{
		{"STRUCT<key INT64>", "KEY"},
		{"STRUCT<row INT64>", "ROW"},
		{"STRUCT<value STRING, data INT64>", "VALUE"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			dt, err := testParseType(tc.input)
			if err != nil {
				t.Fatalf("parseType(%q): unexpected error: %v", tc.input, err)
			}
			if len(dt.Fields) == 0 || strings.ToUpper(dt.Fields[0].Name) != tc.name {
				t.Errorf("first field name = %q, want %q", dt.Fields[0].Name, tc.name)
			}
		})
	}
}

func TestParseType_FieldNameEqualsTypeKeyword(t *testing.T) {
	// The hardest named/anonymous disambiguation: a field whose NAME is a
	// type-name keyword AND whose TYPE is the same keyword. Oracle accepts
	// STRUCT<date DATE> (field `date` : DATE), STRUCT<json JSON>. The lookahead
	// `name type` must win because the token after the name still begins a type.
	cases := []struct {
		input    string
		wantName string
		wantType string
	}{
		{"STRUCT<date DATE>", "DATE", "DATE"},
		{"STRUCT<json JSON>", "JSON", "JSON"},
		{"STRUCT<timestamp TIMESTAMP>", "TIMESTAMP", "TIMESTAMP"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			dt, err := testParseType(tc.input)
			if err != nil {
				t.Fatalf("parseType(%q): unexpected error: %v", tc.input, err)
			}
			if len(dt.Fields) != 1 {
				t.Fatalf("got %d fields, want 1", len(dt.Fields))
			}
			if strings.ToUpper(dt.Fields[0].Name) != tc.wantName {
				t.Errorf("field name = %q, want %q", dt.Fields[0].Name, tc.wantName)
			}
			if strings.ToUpper(dt.Fields[0].Type.String()) != tc.wantType {
				t.Errorf("field type = %q, want %q", dt.Fields[0].Type.String(), tc.wantType)
			}
		})
	}
}

func TestParseType_GreaterEqualSplitLeavesEquals(t *testing.T) {
	// The D1 split must also handle `>=` (tokGreaterEqual) as a template close:
	// the `>` closes the template and the `=` is pushed back. ParseDataType over
	// "ARRAY<INT64>=" must yield a TypeArray plus a trailing-token error for the
	// leftover `=`.
	dt, errs := ParseDataType("ARRAY<INT64>=")
	if dt == nil || dt.Kind != TypeArray {
		t.Fatalf("got %+v, want a TypeArray", dt)
	}
	if len(errs) == 0 {
		t.Error("want a trailing-token error for the leftover '='")
	}
}

func TestParseType_StructTrailingCommaRejected(t *testing.T) {
	// Oracle: `STRUCT<x INT64,>` REJECTS ("Unexpected >") — no trailing comma.
	if _, err := testParseType("STRUCT<x INT64,>"); err == nil {
		t.Error("parseType(\"STRUCT<x INT64,>\"): want error, got nil")
	}
}

func TestParseType_BareStructRejected(t *testing.T) {
	// Oracle: bare `STRUCT` (no `<`) REJECTS ("Expected < but got )").
	if _, err := testParseType("STRUCT"); err == nil {
		t.Error("parseType(\"STRUCT\"): want error, got nil")
	}
}

// ---------------------------------------------------------------------------
// RANGE<type>
// ---------------------------------------------------------------------------

func TestParseType_Range(t *testing.T) {
	// Oracle: `RANGE<DATE>` accepts (semantic "Type not found: RANGE<DATE>" on
	// Spanner, which lacks RANGE, but the grammar parsed it). `RANGE<INT64>`
	// also accepts at parse time (semantic "Unsupported type") — the grammar
	// puts ANY `type` inside RANGE, not just DATE/DATETIME/TIMESTAMP.
	for _, in := range []string{"RANGE<DATE>", "RANGE<DATETIME>", "RANGE<TIMESTAMP>", "RANGE<INT64>"} {
		t.Run(in, func(t *testing.T) {
			dt, err := testParseType(in)
			if err != nil {
				t.Fatalf("parseType(%q): unexpected error: %v", in, err)
			}
			if dt.Kind != TypeRange {
				t.Fatalf("Kind = %v, want TypeRange", dt.Kind)
			}
			if dt.ElementType == nil {
				t.Fatal("ElementType is nil")
			}
		})
	}
}

func TestParseType_RangeRejections(t *testing.T) {
	// Oracle: `RANGE<>` REJECTS ("Expected < but got <>"); bare `RANGE` REJECTS.
	for _, in := range []string{"RANGE<>", "RANGE"} {
		t.Run(in, func(t *testing.T) {
			if _, err := testParseType(in); err == nil {
				t.Errorf("parseType(%q): want error, got nil", in)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// MAP<key, value> and FUNCTION<...> (legacy raw_type alts).
// ---------------------------------------------------------------------------

func TestParseType_Map(t *testing.T) {
	// Oracle: `MAP<INT64, STRING>` accepts (semantic "MAP datatype is not
	// supported" — grammar parsed it). Bare `MAP` is a type_name (accept,
	// "Type not found: MAP") since MAP is non-reserved.
	dt, err := testParseType("MAP<INT64, STRING>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != TypeMap {
		t.Fatalf("Kind = %v, want TypeMap", dt.Kind)
	}
	if dt.KeyType == nil || dt.ValueType == nil {
		t.Fatalf("KeyType/ValueType = %v/%v, want both set", dt.KeyType, dt.ValueType)
	}
	if strings.ToUpper(dt.KeyType.String()) != "INT64" || strings.ToUpper(dt.ValueType.String()) != "STRING" {
		t.Errorf("MAP<%s, %s>, want MAP<INT64, STRING>", dt.KeyType, dt.ValueType)
	}
}

func TestParseType_BareMapIsTypeName(t *testing.T) {
	// Oracle: bare `MAP` accepts as a type_name (MAP is non-reserved).
	dt, err := testParseType("MAP")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != TypeName {
		t.Errorf("Kind = %v, want TypeName (bare MAP is a type_name)", dt.Kind)
	}
}

func TestParseType_Function(t *testing.T) {
	// Oracle: all three accept (semantic "FUNCTION type not supported"):
	//   FUNCTION<INT64 -> STRING>          single arg
	//   FUNCTION<() -> INT64>              no args
	//   FUNCTION<(INT64, STRING) -> BOOL>  multi args
	cases := []struct {
		input   string
		nArgs   int
		wantRet string
	}{
		{"FUNCTION<INT64 -> STRING>", 1, "STRING"},
		{"FUNCTION<() -> INT64>", 0, "INT64"},
		{"FUNCTION<(INT64, STRING) -> BOOL>", 2, "BOOL"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			dt, err := testParseType(tc.input)
			if err != nil {
				t.Fatalf("parseType(%q): unexpected error: %v", tc.input, err)
			}
			if dt.Kind != TypeFunction {
				t.Fatalf("Kind = %v, want TypeFunction", dt.Kind)
			}
			if len(dt.ArgTypes) != tc.nArgs {
				t.Errorf("got %d arg types, want %d", len(dt.ArgTypes), tc.nArgs)
			}
			if dt.ReturnType == nil || strings.ToUpper(dt.ReturnType.String()) != tc.wantRet {
				t.Errorf("return = %v, want %s", dt.ReturnType, tc.wantRet)
			}
		})
	}
}

func TestParseType_BareFunctionIsTypeName(t *testing.T) {
	// Oracle: bare `FUNCTION` accepts as a type_name (FUNCTION is non-reserved).
	dt, err := testParseType("FUNCTION")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != TypeName {
		t.Errorf("Kind = %v, want TypeName (bare FUNCTION is a type_name)", dt.Kind)
	}
}

func TestParseType_MapFunctionMalformed(t *testing.T) {
	// Oracle rejects (the bare MAP/FUNCTION type_name is parsed, then the
	// template tokens are unexpected, OR the template body is incomplete):
	//   MAP<>          reject ("Expected ) ... but got <>")  -> MAP type_name + stray <>
	//   FUNCTION<>     reject (same shape)
	//   MAP<INT64>     reject ("Expected , but got >")        -> map needs 2 types
	//   FUNCTION<INT64> reject ("Expected -> but got >")      -> function needs -> ret
	// (`<>` adjacent lexes as one token, so MAP</FUNCTION< fall through to a
	// bare type_name and the `<>` becomes a trailing token — exactly the oracle
	// shape.)
	for _, in := range []string{"MAP<>", "FUNCTION<>", "MAP<INT64>", "FUNCTION<INT64>"} {
		t.Run(in, func(t *testing.T) {
			if _, errs := ParseDataType(in); len(errs) == 0 {
				t.Errorf("ParseDataType(%q): want error, got none", in)
			}
		})
	}
	// Spaced MAP< … > is a real template and accepts.
	if _, errs := ParseDataType("MAP< INT64, STRING >"); len(errs) != 0 {
		t.Errorf("ParseDataType(\"MAP< INT64, STRING >\"): unexpected errors: %v", errs)
	}
}

// ---------------------------------------------------------------------------
// opt_type_parameters — (n), (n, m), (MAX), bool/string/bytes/float params.
// ---------------------------------------------------------------------------

func TestParseType_TypeParameters(t *testing.T) {
	// Oracle: all accept (semantic "Parameterized types are not supported" on
	// Spanner — the grammar parsed the parameter list). type_parameter is
	// int | bool | string | bytes | float | MAX; multiple are allowed.
	cases := []struct {
		input  string
		nParam int
	}{
		{"STRING(100)", 1},
		{"BYTES(256)", 1},
		{"NUMERIC(10, 2)", 2},
		{"STRING(MAX)", 1},
		{"STRING(MAX, 2)", 2},
		{"NUMERIC(10, 2, 3)", 3},
		{"NUMERIC(true)", 1},
		{"NUMERIC('foo')", 1},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			dt, err := testParseType(tc.input)
			if err != nil {
				t.Fatalf("parseType(%q): unexpected error: %v", tc.input, err)
			}
			if len(dt.Params) != tc.nParam {
				t.Errorf("got %d params, want %d: %+v", len(dt.Params), tc.nParam, dt.Params)
			}
		})
	}
}

func TestParseType_TypeParametersOnNameKind(t *testing.T) {
	// STRING(100) keeps Kind=TypeName (the params hang off the type_name).
	dt, err := testParseType("STRING(100)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != TypeName {
		t.Errorf("Kind = %v, want TypeName", dt.Kind)
	}
	if len(dt.Params) != 1 || !dt.Params[0].IsInt || dt.Params[0].IntVal != 100 {
		t.Errorf("Params = %+v, want one int 100", dt.Params)
	}
}

func TestParseType_TypeParametersMax(t *testing.T) {
	dt, err := testParseType("STRING(MAX)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(dt.Params) != 1 || !dt.Params[0].IsMax {
		t.Errorf("Params = %+v, want one MAX param", dt.Params)
	}
}

func TestParseType_TypeParameterRejections(t *testing.T) {
	// Oracle rejects (syntax):
	//   STRING()       empty list             ("Unexpected )")
	//   STRING(,)      leading comma          ("Unexpected ,")
	//   NUMERIC(-5)    signed (not an int lit)("Unexpected -")
	for _, in := range []string{"STRING()", "STRING(,)", "NUMERIC(-5)"} {
		t.Run(in, func(t *testing.T) {
			if _, err := testParseType(in); err == nil {
				t.Errorf("parseType(%q): want error, got nil", in)
			}
		})
	}
}

func TestParseType_TrailingCommaInTypeParameters(t *testing.T) {
	// Legacy grammar emits a dedicated "Trailing comma in type parameters list
	// is not allowed" error for `(n,)`. The oracle rejects equivalently
	// ("Unexpected )"). Either way it must be an error.
	if _, err := testParseType("NUMERIC(10,)"); err == nil {
		t.Error("parseType(\"NUMERIC(10,)\"): want error, got nil")
	}
}

func TestParseType_ArrayTypeParametersRejectNamedArg(t *testing.T) {
	// Oracle: `ARRAY<INT64>(vector_length=>128)` REJECTS ("Unexpected
	// identifier vector_length") in a bare `type` position — vector_length is a
	// DDL column-schema attribute, NOT a type_parameter. The general `type`
	// rule's opt_type_parameters only accepts int/bool/string/bytes/float/MAX.
	if _, err := testParseType("ARRAY<INT64>(vector_length=>128)"); err == nil {
		t.Error("parseType(\"ARRAY<INT64>(vector_length=>128)\"): want error, got nil")
	}
}

// testParseColumnSchemaType parses input as a column type (a column_schema_inner
// position, where the ARRAY vector-length parameter is admitted), asserting the
// whole input is consumed. It is testParseType with the inArrayColumnSchema flag
// set, mirroring how parseColumnDefinition / parseSpannerAlterColumn call in.
func testParseColumnSchemaType(input string) (*DataType, error) {
	p := &Parser{lexer: NewLexer(input), input: input}
	p.advance()
	dt, err := p.parseColumnSchemaType()
	if err != nil {
		return nil, err
	}
	if p.cur.Type != tokEOF {
		return dt, &ParseError{Loc: p.cur.Loc, Msg: "trailing tokens after type: " + TokenName(p.cur.Type)}
	}
	return dt, nil
}

// TestParseType_ArrayVectorLength is the focused regression for the Spanner
// ARRAY vector-length parameter `ARRAY<scalar>(vector_length => N)`
// (googlesql/maint-vector-array-type; divergence #202). The empirically-verified
// grammar (Spanner emulator @ sha256:caf1bd24) is admitted ONLY in a
// column-schema position, on an ARRAY type, with the exact bare name
// `vector_length` (case-insensitive) and a single integer value.
func TestParseType_ArrayVectorLength(t *testing.T) {
	// --- accepts in column-schema position (oracle: CREATE TABLE … accepts) ---
	accepts := []struct {
		input    string
		wantText string // round-trip rendering via DataType.String()
	}{
		{"ARRAY<FLOAT32>(vector_length=>128)", "ARRAY<FLOAT32>(vector_length => 128)"},
		{"ARRAY<FLOAT64>(vector_length => 256)", "ARRAY<FLOAT64>(vector_length => 256)"},
		// Case-insensitive name (oracle: VECTOR_LENGTH accepts); spelling preserved.
		{"ARRAY<FLOAT32>(VECTOR_LENGTH=>8)", "ARRAY<FLOAT32>(VECTOR_LENGTH => 8)"},
		// Grammar admits any element type + nested ARRAY (FLOAT-only / array-of-array
		// is a SEMANTIC restriction, not syntax — oracle parses then semantic-rejects).
		{"ARRAY<INT64>(vector_length=>4)", "ARRAY<INT64>(vector_length => 4)"},
		{"ARRAY<ARRAY<FLOAT32>(vector_length=>8)>", "ARRAY<ARRAY<FLOAT32>(vector_length => 8)>"},
	}
	for _, tc := range accepts {
		dt, err := testParseColumnSchemaType(tc.input)
		if err != nil {
			t.Errorf("column-schema parseType(%q): unexpected error: %v", tc.input, err)
			continue
		}
		if got := dt.String(); got != tc.wantText {
			t.Errorf("column-schema parseType(%q) round-trip = %q, want %q", tc.input, got, tc.wantText)
		}
	}

	// --- rejects even in column-schema position (oracle rejects each) ---
	rejects := []struct {
		input  string
		reason string
	}{
		// Wrong name (oracle: "Expecting 'VECTOR_LENGTH' but found 'length'").
		{"ARRAY<FLOAT32>(length=>8)", "non-vector_length name"},
		// MAX value (oracle: "Expecting '<INTEGER_LITERAL>'").
		{"ARRAY<FLOAT32>(vector_length=>MAX)", "non-integer value"},
		// String value — not an integer literal.
		{"ARRAY<FLOAT32>(vector_length=>'8')", "string value"},
		// Trailing comma / second param (oracle: "Expecting ')'").
		{"ARRAY<FLOAT32>(vector_length=>8,)", "trailing comma"},
		{"ARRAY<FLOAT32>(vector_length=>8, foo=>9)", "second parameter"},
		// Backtick-quoted name (oracle: "Expecting 'VECTOR_LENGTH'").
		{"ARRAY<FLOAT32>(`vector_length`=>8)", "quoted name"},
		// Not an ARRAY type — the parameter only attaches to ARRAY.
		{"STRING(vector_length=>8)", "non-array type"},
		// A positional type_parameter on an ARRAY is a syntax error in column
		// schema (oracle rejects "Error parsing Spanner DDL statement"); only the
		// vector_length named form is admissible. (In a general `type` position the
		// same `ARRAY<FLOAT32>(128)` PARSES — see TestParseType below.)
		{"ARRAY<FLOAT32>(128)", "positional param on array column"},
		{"ARRAY<FLOAT32>(MAX)", "positional MAX on array column"},
		// Missing `=>` after the name (oracle rejects).
		{"ARRAY<FLOAT32>(vector_length)", "missing => and value"},
		// Negative value — a sign is not part of an integer literal (oracle rejects).
		{"ARRAY<FLOAT32>(vector_length=>-5)", "negative value"},
		// Empty parens (oracle rejects).
		{"ARRAY<FLOAT32>()", "empty parens"},
	}
	for _, tc := range rejects {
		if _, err := testParseColumnSchemaType(tc.input); err == nil {
			t.Errorf("column-schema parseType(%q): want error (%s), got nil", tc.input, tc.reason)
		}
	}

	// --- the vector-length parameter is NOT admitted in a general `type` position
	// (a CAST target / standalone ParseDataType) even on an ARRAY (oracle: a
	// CAST<ARRAY<…>(vector_length=>…)> rejects "Unexpected identifier"). ---
	for _, in := range []string{
		"ARRAY<FLOAT32>(vector_length=>8)",
		"ARRAY<INT64>(vector_length=>128)",
	} {
		if _, err := testParseType(in); err == nil {
			t.Errorf("bare parseType(%q): want error (vector_length is column-schema-only), got nil", in)
		}
	}

	// --- conversely, a POSITIONAL type_parameter list on an ARRAY is fine in a
	// general `type` position (oracle: a CAST<ARRAY<FLOAT32>(128)> PARSES,
	// "Parameterized types are not supported" being a semantic message). Only the
	// column-schema position restricts an array to the vector_length form. ---
	for _, in := range []string{
		"ARRAY<FLOAT32>(128)",
		"ARRAY<INT64>(10, 2)",
	} {
		if _, err := testParseType(in); err != nil {
			t.Errorf("bare parseType(%q): unexpected error (positional array type-params parse in a general type position): %v", in, err)
		}
	}

	// --- value spellings the column-schema vector_length form accepts (oracle:
	// 0 and hex both accept). ---
	for _, in := range []string{
		"ARRAY<FLOAT32>(vector_length=>0)",
		"ARRAY<FLOAT32>(vector_length=>0x80)",
	} {
		if _, err := testParseColumnSchemaType(in); err != nil {
			t.Errorf("column-schema parseType(%q): unexpected error: %v", in, err)
		}
	}
}

// TestParseType_ArrayVectorLengthStatements exercises the ARRAY vector-length
// parameter through the real statement entry point Parse — the CREATE
// TABLE / ALTER TABLE column-schema positions a consumer hits — plus the
// negatives the live oracle rejects (STRUCT-of-ARRAY, and the SELECT value
// constructor, which the oracle itself rejects so omni's reject is parity).
func TestParseType_ArrayVectorLengthStatements(t *testing.T) {
	accepts := []string{
		// CREATE TABLE column (oracle: accept).
		"CREATE TABLE Products (ProductId INT64 NOT NULL, Embedding ARRAY<FLOAT32>(vector_length=>128)) PRIMARY KEY (ProductId)",
		"CREATE TABLE Products (ProductId INT64 NOT NULL, Embedding ARRAY<FLOAT64>(vector_length => 256)) PRIMARY KEY (ProductId)",
		// With a trailing NOT NULL on the same column (oracle: accept).
		"CREATE TABLE Products (ProductId INT64 NOT NULL, Embedding ARRAY<FLOAT32>(vector_length=>128) NOT NULL) PRIMARY KEY (ProductId)",
		// Spanner ALTER COLUMN <name> <column_schema_inner> (oracle: parses,
		// semantic-rejects the change).
		"ALTER TABLE Products ALTER COLUMN Embedding ARRAY<FLOAT32>(vector_length=>128)",
		// BigQuery ALTER COLUMN … SET DATA TYPE field_schema (column-schema position).
		"ALTER TABLE Products ALTER COLUMN Embedding SET DATA TYPE ARRAY<FLOAT32>(vector_length=>128)",
		// ADD COLUMN with the parameter (oracle: accept).
		"ALTER TABLE Products ADD COLUMN E2 ARRAY<FLOAT32>(vector_length=>64)",
	}
	for _, sql := range accepts {
		if _, errs := Parse(sql); len(errs) > 0 {
			t.Errorf("Parse(%q): unexpected errors: %v", sql, errs)
		}
	}

	rejects := []struct {
		sql    string
		reason string
	}{
		// STRUCT-of-ARRAY column (oracle: "'vector_length' is not supported in
		// STRUCT of ARRAY" — a true syntax reject). The flag is cleared inside a
		// STRUCT template, so omni rejects too.
		{
			"CREATE TABLE T (Id INT64 NOT NULL, E STRUCT<v ARRAY<FLOAT32>(vector_length=>8)>) PRIMARY KEY (Id)",
			"STRUCT-of-ARRAY",
		},
		// The SELECT value-constructor form: oracle REJECTS ("Expected '[' but got
		// '('"); omni rejects too — parity (defended against the stale divergence
		// #202 claim that it ACCEPTS).
		{"SELECT ARRAY<FLOAT32>(vector_length => 128)", "value constructor (oracle rejects)"},
		{"SELECT ARRAY<FLOAT32>()", "empty typed-array constructor (oracle rejects)"},
		// CAST target: vector_length is not a type_parameter in a general type
		// position (oracle: "Unexpected identifier vector_length").
		{"SELECT CAST(NULL AS ARRAY<FLOAT32>(vector_length=>2))", "CAST target"},
	}
	for _, tc := range rejects {
		if _, errs := Parse(tc.sql); len(errs) == 0 {
			t.Errorf("Parse(%q): want parse error (%s), got none", tc.sql, tc.reason)
		}
	}
}

// ---------------------------------------------------------------------------
// collate_clause — type COLLATE 'string-or-parameter'
// ---------------------------------------------------------------------------

func TestParseType_Collate(t *testing.T) {
	// Oracle: `INT64 COLLATE 'und:ci'` and `DATE COLLATE 'und:ci'` accept
	// (semantic "Type with collation name is not supported in cast" — parsed).
	for _, in := range []string{"INT64 COLLATE 'und:ci'", "DATE COLLATE 'und:ci'", "STRING COLLATE 'und'"} {
		t.Run(in, func(t *testing.T) {
			dt, err := testParseType(in)
			if err != nil {
				t.Fatalf("parseType(%q): unexpected error: %v", in, err)
			}
			if dt.Collate == "" {
				t.Errorf("Collate = empty, want the collation string")
			}
		})
	}
}

func TestParseType_CollateOnArrayElement(t *testing.T) {
	// Oracle: `ARRAY<INT64 COLLATE 'und:ci'>` accepts — collate binds to the
	// element type inside the array, not the array.
	dt, err := testParseType("ARRAY<INT64 COLLATE 'und:ci'>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dt.Kind != TypeArray {
		t.Fatalf("Kind = %v, want TypeArray", dt.Kind)
	}
	if dt.ElementType == nil || dt.ElementType.Collate == "" {
		t.Errorf("element collate not captured: %+v", dt.ElementType)
	}
	if dt.Collate != "" {
		t.Errorf("array Collate = %q, want empty (collate is on the element)", dt.Collate)
	}
}

func TestParseType_CollateParameter(t *testing.T) {
	// collate_clause: COLLATE string_literal_or_parameter, where parameter is
	// @name | ? | @@sysvar. The grammar accepts a parameter; capture it.
	for _, in := range []string{"STRING COLLATE @p", "STRING COLLATE ?"} {
		t.Run(in, func(t *testing.T) {
			if _, err := testParseType(in); err != nil {
				t.Errorf("parseType(%q): unexpected error: %v", in, err)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Round-trip: String() re-renders to syntax the parser re-accepts.
// ---------------------------------------------------------------------------

func TestParseType_StringRoundTrip(t *testing.T) {
	inputs := []string{
		"INT64",
		"STRING(100)",
		"NUMERIC(10, 2)",
		"STRING(MAX)",
		"ARRAY<INT64>",
		"ARRAY<ARRAY<INT64>>",
		"STRUCT<x INT64, y STRING>",
		"STRUCT<INT64, STRING>",
		"STRUCT<>",
		"STRUCT<a ARRAY<INT64>>",
		"RANGE<DATE>",
		"MAP<INT64, STRING>",
		"FUNCTION<(INT64, STRING) -> BOOL>",
		"INTERVAL",
		"foo.bar.Baz",
	}
	for _, in := range inputs {
		t.Run(in, func(t *testing.T) {
			dt1, err := testParseType(in)
			if err != nil {
				t.Fatalf("parseType(%q): %v", in, err)
			}
			rendered := dt1.String()
			dt2, err := testParseType(rendered)
			if err != nil {
				t.Fatalf("re-parse of String()=%q failed: %v", rendered, err)
			}
			if dt1.String() != dt2.String() {
				t.Errorf("round-trip unstable: %q -> %q -> %q", in, rendered, dt2.String())
			}
		})
	}
}

// ---------------------------------------------------------------------------
// ParseDataType — standalone string entry point.
// ---------------------------------------------------------------------------

func TestParseDataType_Standalone(t *testing.T) {
	dt, errs := ParseDataType("ARRAY<STRUCT<x INT64>>")
	if len(errs) != 0 {
		t.Fatalf("ParseDataType: unexpected errors: %v", errs)
	}
	if dt == nil || dt.Kind != TypeArray {
		t.Fatalf("got %+v, want a TypeArray", dt)
	}
}

func TestParseDataType_TrailingTokens(t *testing.T) {
	// A complete type followed by junk reports an error but still returns the
	// type (mirrors snowflake/trino ParseDataType).
	dt, errs := ParseDataType("INT64 garbage")
	if dt == nil {
		t.Fatal("ParseDataType returned nil type")
	}
	if len(errs) == 0 {
		t.Error("want a trailing-token error, got none")
	}
}

func TestParseDataType_Invalid(t *testing.T) {
	_, errs := ParseDataType("SELECT")
	if len(errs) == 0 {
		t.Error("ParseDataType(\"SELECT\"): want an error, got none")
	}
}

// ---------------------------------------------------------------------------
// Loc — every node carries a valid byte span.
// ---------------------------------------------------------------------------

func TestParseType_Loc(t *testing.T) {
	in := "ARRAY<INT64>"
	dt, err := testParseType(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !dt.Loc.IsValid() {
		t.Fatalf("Loc invalid: %+v", dt.Loc)
	}
	if dt.Loc.Start != 0 || dt.Loc.End != len(in) {
		t.Errorf("Loc = %+v, want {0, %d}", dt.Loc, len(in))
	}
	if dt.ElementType == nil || !dt.ElementType.Loc.IsValid() {
		t.Errorf("element Loc invalid: %+v", dt.ElementType)
	}
}
