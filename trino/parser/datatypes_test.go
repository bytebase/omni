package parser

import (
	"testing"
)

// This file is the `types` node's correctness gate. It has three layers, in
// order of authority (correctness-protocol.md):
//
//  1. Structural unit tests — assert ParseDataType builds the right DataType
//     shape (Kind, Name, params, nesting, interval fields, time-zone flags) for
//     representative inputs. These pin the AST contract downstream nodes embed.
//
//  2. Differential oracle gate — the authoritative accept/reject check. Each
//     type is wrapped as `CAST(NULL AS <type>)` (the canonical type position)
//     and omni's accept/reject of the standalone type must equal Trino 481's
//     accept/reject of the wrapped statement. The corpus carries BOTH accepted
//     types (every documented form + legacy forms) and rejected ones (negative
//     coverage: empty parens, reversed/cross-family interval ranges, malformed
//     angle brackets, trailing commas). Classification keys on SYNTAX_ERROR
//     only — TYPE_MISMATCH / NOT_SUPPORTED / GENERIC_INTERNAL_ERROR all mean
//     Trino's grammar ACCEPTED and the failure is semantic.
//
//  3. Deparse round-trip — parse → String() → re-parse must yield an equal
//     accept verdict and a stable structure, and the re-rendered type must
//     still be accepted by the oracle. This is the structural gate (no
//     reference parser exists for Trino types).
//
// All oracle-backed subtests skip cleanly when no Trino is reachable.

// ---------------------------------------------------------------------------
// Layer 1 — structural unit tests
// ---------------------------------------------------------------------------

func TestDataType_Structure(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(t *testing.T, dt *DataType)
	}{
		{
			name:  "generic_no_params",
			input: "bigint",
			check: func(t *testing.T, dt *DataType) {
				if dt.Kind != TypeGeneric || dt.Name != "bigint" || dt.Params != nil {
					t.Errorf("got Kind=%v Name=%q Params=%v", dt.Kind, dt.Name, dt.Params)
				}
			},
		},
		{
			name:  "generic_one_int_param",
			input: "varchar(100)",
			check: func(t *testing.T, dt *DataType) {
				if dt.Kind != TypeGeneric || dt.Name != "varchar" || len(dt.Params) != 1 {
					t.Fatalf("got Kind=%v Name=%q Params=%v", dt.Kind, dt.Name, dt.Params)
				}
				if !dt.Params[0].IsInt || dt.Params[0].IntVal != 100 {
					t.Errorf("param[0]=%+v, want int 100", dt.Params[0])
				}
			},
		},
		{
			name:  "generic_two_int_params",
			input: "decimal(38, 10)",
			check: func(t *testing.T, dt *DataType) {
				if len(dt.Params) != 2 || !dt.Params[0].IsInt || dt.Params[0].IntVal != 38 ||
					!dt.Params[1].IsInt || dt.Params[1].IntVal != 10 {
					t.Errorf("params=%+v, want [38 10]", dt.Params)
				}
			},
		},
		{
			// An int parameter beyond int64 range keeps its exact source text so
			// String() round-trips it faithfully (regression: the cross-review
			// gate flagged that rendering from the parsed value alone prints 0).
			name:  "generic_overflow_int_param",
			input: "decimal(9223372036854775808)",
			check: func(t *testing.T, dt *DataType) {
				if len(dt.Params) != 1 || !dt.Params[0].IsInt {
					t.Fatalf("params=%+v", dt.Params)
				}
				if dt.Params[0].IntText != "9223372036854775808" {
					t.Errorf("IntText=%q, want exact source", dt.Params[0].IntText)
				}
				if got := dt.String(); got != "decimal(9223372036854775808)" {
					t.Errorf("String()=%q, want decimal(9223372036854775808)", got)
				}
			},
		},
		{
			name:  "generic_nested_type_param",
			input: "qdigest(bigint)",
			check: func(t *testing.T, dt *DataType) {
				if len(dt.Params) != 1 || dt.Params[0].IsInt || dt.Params[0].Type == nil {
					t.Fatalf("params=%+v, want one nested-type param", dt.Params)
				}
				if dt.Params[0].Type.Name != "bigint" {
					t.Errorf("nested param type=%q, want bigint", dt.Params[0].Type.Name)
				}
			},
		},
		{
			name:  "double_precision",
			input: "double precision",
			check: func(t *testing.T, dt *DataType) {
				if dt.Kind != TypeGeneric || dt.Name != "DOUBLE PRECISION" {
					t.Errorf("got Kind=%v Name=%q, want Generic DOUBLE PRECISION", dt.Kind, dt.Name)
				}
			},
		},
		{
			name:  "double_alone_is_generic",
			input: "double",
			check: func(t *testing.T, dt *DataType) {
				if dt.Kind != TypeGeneric || dt.Name != "double" {
					t.Errorf("got Kind=%v Name=%q, want Generic double", dt.Kind, dt.Name)
				}
			},
		},
		{
			name:  "row_named_fields",
			input: "ROW(x bigint, y double)",
			check: func(t *testing.T, dt *DataType) {
				if dt.Kind != TypeRow || len(dt.Fields) != 2 {
					t.Fatalf("got Kind=%v Fields=%+v", dt.Kind, dt.Fields)
				}
				if dt.Fields[0].Name == nil || dt.Fields[0].Name.Value != "x" || dt.Fields[0].Type.Name != "bigint" {
					t.Errorf("field0=%+v, want name x type bigint", dt.Fields[0])
				}
				if dt.Fields[1].Name == nil || dt.Fields[1].Name.Value != "y" {
					t.Errorf("field1=%+v, want name y", dt.Fields[1])
				}
			},
		},
		{
			name:  "row_unnamed_fields",
			input: "ROW(bigint, varchar)",
			check: func(t *testing.T, dt *DataType) {
				if dt.Kind != TypeRow || len(dt.Fields) != 2 {
					t.Fatalf("got Kind=%v Fields=%+v", dt.Kind, dt.Fields)
				}
				if dt.Fields[0].Name != nil || dt.Fields[0].Type.Name != "bigint" {
					t.Errorf("field0=%+v, want unnamed bigint", dt.Fields[0])
				}
			},
		},
		{
			name:  "row_quoted_field_name",
			input: `ROW("x y" bigint)`,
			check: func(t *testing.T, dt *DataType) {
				if len(dt.Fields) != 1 || dt.Fields[0].Name == nil || !dt.Fields[0].Name.Quoted ||
					dt.Fields[0].Name.Value != "x y" {
					t.Errorf("field0=%+v, want quoted name 'x y'", dt.Fields[0])
				}
			},
		},
		{
			name:  "row_nested_type_field",
			input: "ROW(a ARRAY(int))",
			check: func(t *testing.T, dt *DataType) {
				if len(dt.Fields) != 1 || dt.Fields[0].Name == nil || dt.Fields[0].Name.Value != "a" {
					t.Fatalf("fields=%+v", dt.Fields)
				}
				ft := dt.Fields[0].Type
				if ft == nil || ft.Kind != TypeGeneric || ft.Name != "ARRAY" || len(ft.Params) != 1 {
					t.Errorf("field type=%+v, want generic ARRAY(int)", ft)
				}
			},
		},
		{
			// Unnamed field whose type is a multi-token dateTimeType: the leading
			// TIME must be parsed as the type, NOT consumed as a field name.
			name:  "row_unnamed_datetime_field",
			input: "ROW(TIME WITHOUT TIME ZONE)",
			check: func(t *testing.T, dt *DataType) {
				if len(dt.Fields) != 1 {
					t.Fatalf("fields=%+v", dt.Fields)
				}
				f := dt.Fields[0]
				if f.Name != nil {
					t.Errorf("field should be unnamed, got name=%v", f.Name)
				}
				if f.Type == nil || f.Type.Kind != TypeDateTime || f.Type.Name != "TIME" || f.Type.WithTimeZone {
					t.Errorf("field type=%+v, want unnamed TIME WITHOUT TIME ZONE", f.Type)
				}
			},
		},
		{
			// Unnamed field whose type is a multi-token intervalType.
			name:  "row_unnamed_interval_field",
			input: "ROW(INTERVAL DAY TO SECOND)",
			check: func(t *testing.T, dt *DataType) {
				if len(dt.Fields) != 1 || dt.Fields[0].Name != nil {
					t.Fatalf("fields=%+v, want one unnamed field", dt.Fields)
				}
				ft := dt.Fields[0].Type
				if ft == nil || ft.Kind != TypeInterval || ft.IntervalFrom != IntervalDay ||
					ft.IntervalTo == nil || *ft.IntervalTo != IntervalSecond {
					t.Errorf("field type=%+v, want INTERVAL DAY TO SECOND", ft)
				}
			},
		},
		{
			// Named field whose type is a legacy angle-bracket array: `a` is the
			// name, ARRAY<int> the type (needs >2-token disambiguation).
			name:  "row_named_legacy_array_field",
			input: "ROW(a ARRAY<int>)",
			check: func(t *testing.T, dt *DataType) {
				if len(dt.Fields) != 1 || dt.Fields[0].Name == nil || dt.Fields[0].Name.Value != "a" {
					t.Fatalf("fields=%+v", dt.Fields)
				}
				ft := dt.Fields[0].Type
				if ft == nil || ft.Kind != TypeArray || ft.ElementType == nil || ft.ElementType.Name != "int" {
					t.Errorf("field type=%+v, want ARRAY<int>", ft)
				}
			},
		},
		{
			// Unnamed DOUBLE PRECISION field (ANTLR's unnamed alt is tried first,
			// so this is one type, not name 'double' + type 'precision').
			name:  "row_unnamed_double_precision_field",
			input: "ROW(double precision)",
			check: func(t *testing.T, dt *DataType) {
				if len(dt.Fields) != 1 || dt.Fields[0].Name != nil {
					t.Fatalf("fields=%+v, want one unnamed field", dt.Fields)
				}
				if dt.Fields[0].Type.Name != "DOUBLE PRECISION" {
					t.Errorf("field type name=%q, want DOUBLE PRECISION", dt.Fields[0].Type.Name)
				}
			},
		},
		{
			// An all-integer ROW body has no valid rowField reading (rowField has
			// no integer alternative), so it falls back to the genericType "ROW"
			// reading with INTEGER typeParameters — NOT a TypeRow.
			name:  "row_all_integers_is_generic",
			input: "ROW(1, 2)",
			check: func(t *testing.T, dt *DataType) {
				if dt.Kind != TypeGeneric || dt.Name != "ROW" || len(dt.Params) != 2 {
					t.Fatalf("got Kind=%v Name=%q Params=%+v, want generic ROW(1,2)", dt.Kind, dt.Name, dt.Params)
				}
				if !dt.Params[0].IsInt || dt.Params[0].IntVal != 1 || !dt.Params[1].IsInt || dt.Params[1].IntVal != 2 {
					t.Errorf("params=%+v, want int [1 2]", dt.Params)
				}
			},
		},
		{
			name:  "interval_single_field",
			input: "INTERVAL DAY",
			check: func(t *testing.T, dt *DataType) {
				if dt.Kind != TypeInterval || dt.IntervalFrom != IntervalDay || dt.IntervalTo != nil {
					t.Errorf("got Kind=%v from=%v to=%v", dt.Kind, dt.IntervalFrom, dt.IntervalTo)
				}
			},
		},
		{
			name:  "interval_range",
			input: "INTERVAL DAY TO SECOND",
			check: func(t *testing.T, dt *DataType) {
				if dt.Kind != TypeInterval || dt.IntervalFrom != IntervalDay ||
					dt.IntervalTo == nil || *dt.IntervalTo != IntervalSecond {
					t.Errorf("got Kind=%v from=%v to=%v", dt.Kind, dt.IntervalFrom, dt.IntervalTo)
				}
			},
		},
		{
			name:  "interval_year_to_month",
			input: "INTERVAL YEAR TO MONTH",
			check: func(t *testing.T, dt *DataType) {
				if dt.IntervalFrom != IntervalYear || dt.IntervalTo == nil || *dt.IntervalTo != IntervalMonth {
					t.Errorf("got from=%v to=%v", dt.IntervalFrom, dt.IntervalTo)
				}
			},
		},
		{
			name:  "interval_bare_is_generic",
			input: "INTERVAL",
			check: func(t *testing.T, dt *DataType) {
				// D2: bare INTERVAL (no field) falls through to genericType.
				if dt.Kind != TypeGeneric || dt.Name != "INTERVAL" {
					t.Errorf("got Kind=%v Name=%q, want Generic INTERVAL", dt.Kind, dt.Name)
				}
			},
		},
		{
			name:  "interval_paren_is_generic",
			input: "INTERVAL(5)",
			check: func(t *testing.T, dt *DataType) {
				// D2: INTERVAL '(' is genericType, not intervalType.
				if dt.Kind != TypeGeneric || dt.Name != "INTERVAL" || len(dt.Params) != 1 {
					t.Errorf("got Kind=%v Name=%q params=%+v", dt.Kind, dt.Name, dt.Params)
				}
			},
		},
		{
			name:  "timestamp_bare",
			input: "TIMESTAMP",
			check: func(t *testing.T, dt *DataType) {
				if dt.Kind != TypeDateTime || dt.Name != "TIMESTAMP" || dt.Precision != nil || dt.HasTimeZoneClause {
					t.Errorf("got %+v", dt)
				}
			},
		},
		{
			name:  "timestamp_precision_with_tz",
			input: "TIMESTAMP(6) WITH TIME ZONE",
			check: func(t *testing.T, dt *DataType) {
				if dt.Kind != TypeDateTime || dt.Precision == nil || !dt.Precision.IsInt || dt.Precision.IntVal != 6 {
					t.Fatalf("precision=%+v", dt.Precision)
				}
				if !dt.HasTimeZoneClause || !dt.WithTimeZone {
					t.Errorf("tz: has=%v with=%v, want true/true", dt.HasTimeZoneClause, dt.WithTimeZone)
				}
			},
		},
		{
			name:  "time_without_tz",
			input: "TIME WITHOUT TIME ZONE",
			check: func(t *testing.T, dt *DataType) {
				if dt.Name != "TIME" || !dt.HasTimeZoneClause || dt.WithTimeZone {
					t.Errorf("got Name=%q has=%v with=%v", dt.Name, dt.HasTimeZoneClause, dt.WithTimeZone)
				}
			},
		},
		{
			name:  "legacy_array",
			input: "ARRAY<bigint>",
			check: func(t *testing.T, dt *DataType) {
				if dt.Kind != TypeArray || dt.ElementType == nil || dt.ElementType.Name != "bigint" {
					t.Errorf("got %+v", dt)
				}
			},
		},
		{
			name:  "legacy_map",
			input: "MAP<varchar, integer>",
			check: func(t *testing.T, dt *DataType) {
				if dt.Kind != TypeMap || dt.KeyType == nil || dt.KeyType.Name != "varchar" ||
					dt.ValueType == nil || dt.ValueType.Name != "integer" {
					t.Errorf("got %+v", dt)
				}
			},
		},
		{
			name:  "array_paren_is_generic",
			input: "ARRAY(bigint)",
			check: func(t *testing.T, dt *DataType) {
				// ARRAY '(' → genericType ARRAY with a nested-type param.
				if dt.Kind != TypeGeneric || dt.Name != "ARRAY" || len(dt.Params) != 1 || dt.Params[0].Type == nil {
					t.Errorf("got %+v", dt)
				}
			},
		},
		{
			name:  "postfix_array",
			input: "bigint ARRAY",
			check: func(t *testing.T, dt *DataType) {
				if dt.Kind != TypeArray || dt.ElementType == nil || dt.ElementType.Name != "bigint" {
					t.Errorf("got %+v", dt)
				}
			},
		},
		{
			name:  "postfix_array_with_dim",
			input: "bigint ARRAY[5]",
			check: func(t *testing.T, dt *DataType) {
				if dt.Kind != TypeArray || dt.ArrayDim != 5 {
					t.Errorf("got Kind=%v dim=%d", dt.Kind, dt.ArrayDim)
				}
			},
		},
		{
			name:  "postfix_array_nested",
			input: "bigint ARRAY ARRAY",
			check: func(t *testing.T, dt *DataType) {
				if dt.Kind != TypeArray || dt.ElementType == nil || dt.ElementType.Kind != TypeArray {
					t.Fatalf("got %+v", dt)
				}
				if dt.ElementType.ElementType == nil || dt.ElementType.ElementType.Name != "bigint" {
					t.Errorf("inner element=%+v", dt.ElementType.ElementType)
				}
			},
		},
		{
			name:  "nested_array_of_array_paren",
			input: "ARRAY(ARRAY(bigint))",
			check: func(t *testing.T, dt *DataType) {
				if dt.Kind != TypeGeneric || len(dt.Params) != 1 || dt.Params[0].Type == nil {
					t.Fatalf("got %+v", dt)
				}
				inner := dt.Params[0].Type
				if inner.Kind != TypeGeneric || inner.Name != "ARRAY" {
					t.Errorf("inner=%+v, want generic ARRAY", inner)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dt, errs := ParseDataType(tc.input)
			if len(errs) != 0 {
				t.Fatalf("ParseDataType(%q) errors: %v", tc.input, errs)
			}
			if dt == nil {
				t.Fatalf("ParseDataType(%q) returned nil", tc.input)
			}
			tc.check(t, dt)
		})
	}
}

// TestDataType_LocSpan asserts the parsed type's Loc covers exactly the type
// text (no leading/trailing slop), which downstream completion/diagnostics
// depend on for accurate ranges.
func TestDataType_LocSpan(t *testing.T) {
	for _, input := range []string{
		"bigint",
		"varchar(100)",
		"decimal(38, 10)",
		"ROW(x bigint, y double)",
		"INTERVAL DAY TO SECOND",
		"TIMESTAMP(6) WITH TIME ZONE",
		"ARRAY<bigint>",
		"MAP<varchar, integer>",
		"bigint ARRAY[5]",
		"double precision",
	} {
		t.Run(truncateName(input), func(t *testing.T) {
			dt, errs := ParseDataType(input)
			if len(errs) != 0 {
				t.Fatalf("errors: %v", errs)
			}
			if dt.Loc.Start != 0 || dt.Loc.End != len(input) {
				t.Errorf("Loc=%+v, want [0,%d) covering %q", dt.Loc, len(input), input)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Layer 2 — differential oracle corpus
// ---------------------------------------------------------------------------

// typeAcceptCorpus lists types Trino 481 accepts syntactically (it may still
// fail later semantically — that is fine, classification keys on SYNTAX_ERROR).
// It covers every documented type form (truth1) plus the legacy ARRAY<>/MAP<>
// spellings and the postfix-array forms.
var typeAcceptCorpus = []string{
	// --- generic scalar types (type names are identifiers, not keywords) ---
	"boolean",
	"tinyint", "smallint", "integer", "int", "bigint",
	"real", "double", "double precision",
	"decimal", "decimal(38)", "decimal(38, 10)",
	"varchar", "varchar(100)", "char", "char(10)",
	"varbinary",
	"json", "variant",
	"date",
	"uuid", "ipaddress",
	"hyperloglog", "p4hyperloglog", "tdigest", "setdigest",
	"qdigest(bigint)", "qdigest(double)",
	// quoted type name
	`"varchar"(10)`,
	// --- time / timestamp with precision and time zone ---
	"time", "time(3)", "time with time zone", "time(6) with time zone",
	"time without time zone",
	"timestamp", "timestamp(6)", "timestamp with time zone",
	"timestamp(9) with time zone", "timestamp without time zone",
	// --- interval (single field + valid ranges; D1) ---
	"interval year", "interval month", "interval day",
	"interval hour", "interval minute", "interval second",
	"interval year to month",
	"interval day to hour", "interval day to minute", "interval day to second",
	"interval hour to minute", "interval hour to second",
	"interval minute to second",
	// D2: INTERVAL not followed by a field is a generic type name
	"interval", "interval(5)", "interval array",
	// --- ROW (named, unnamed, nested, quoted, integer-spelling backtrack) ---
	"ROW(x bigint, y double)",
	"ROW(bigint, double)",
	"ROW(bigint)",
	`ROW("x y" bigint)`,
	"ROW(a ARRAY(int))",
	"ROW(a int ARRAY)",
	"ROW(x map(varchar, int))",
	"ROW(a double precision)",
	"ROW(a row(b int))",
	"ROW(1)",    // backtracks into genericType typeParameter
	"ROW(a b)",  // both identifiers (a=name, b=type)
	"ROW(a, b)", // two unnamed generic fields
	"ROW(a bigint, b)",
	"row(x double precision)",
	// ALL-typeParameter ROW bodies (every entry is INTEGER | bare type) fall
	// back to the genericType "ROW" reading, which is the only one that admits
	// integers. Mixing an integer with a NAMED field fits neither reading and is
	// a reject (see typeRejectCorpus).
	"ROW(1, 2)",
	"ROW(bigint, 1)",
	"ROW(varchar(10), 5)",
	"ROW(5, 10, 15)",
	// fields whose (unnamed) type is a structural/keyword type, and the
	// converse where a type keyword serves as a field NAME — accept/reject must
	// match regardless of how omni assigns names internally.
	"ROW(timestamp)",
	"ROW(interval day)",
	"ROW(double precision)",
	"ROW(array<int>)",
	"ROW(map<int, int>)",
	"ROW(time bigint)",      // 'time' keyword as a field NAME
	"ROW(timestamp bigint)", // 'timestamp' keyword as a field NAME
	"ROW(array bigint)",     // 'array' keyword as a field NAME
	"ROW(a bigint array)",   // named field of a postfix-array type
	"ROW(bigint array)",     // unnamed postfix-array field
	"ROW(a timestamp(3) with time zone)",
	"ROW(a interval day to second)",
	// UNNAMED multi-token field types — the leading keyword must be taken as a
	// type, not over-consumed as a field name (regression: caught by the
	// cross-review gate; a 2-token lookahead heuristic mis-rejected these).
	"ROW(TIME WITHOUT TIME ZONE)",
	"ROW(TIMESTAMP WITH TIME ZONE)",
	"ROW(INTERVAL DAY TO SECOND)",
	"ROW(BIGINT ARRAY[5])",
	"ROW(DOUBLE PRECISION)",
	"ROW(a ARRAY<int>)", // named field whose type is a legacy angle-bracket array
	"ROW(a ARRAY(int))", // named field whose type is a paren array
	"ROW(x TIME WITHOUT TIME ZONE, y INTERVAL DAY TO SECOND)",
	// many-field ROW mixing named/unnamed, parameterized, and multi-token types
	// (exercises the speculative name/type backtracking across several fields).
	"ROW(a int, b varchar(10), c TIMESTAMP WITH TIME ZONE, d INTERVAL DAY TO SECOND, e ARRAY<int>)",
	"ROW(a ARRAY<int>, b ARRAY(int), c bigint ARRAY[5])",
	// --- ARRAY / MAP, both legacy <> and modern () spellings ---
	"ARRAY<bigint>",
	"ARRAY(bigint)",
	"ARRAY(ARRAY(bigint))",
	"ARRAY<ARRAY<bigint>>",
	"MAP<varchar, integer>",
	"MAP(varchar, integer)",
	"MAP(varchar, ARRAY(bigint))",
	"MAP<varchar, ARRAY<bigint>>",
	"ARRAY(ROW(a bigint, b varchar))",
	// --- postfix array ---
	"bigint ARRAY",
	"bigint ARRAY ARRAY",
	"bigint ARRAY[5]",
	"decimal(10, 2) ARRAY",
	"INTERVAL DAY TO SECOND ARRAY",
	"TIMESTAMP(3) WITH TIME ZONE ARRAY",
	// --- nested combinations ---
	"MAP(varchar, ROW(a bigint, b ARRAY(double)))",
	"ARRAY(MAP(varchar, bigint))",
}

// typeRejectCorpus lists types Trino 481 rejects with SYNTAX_ERROR. This is the
// required negative coverage: without it the parser could be over-permissive
// and still pass every accept case.
var typeRejectCorpus = []string{
	// empty parameter lists
	"varchar()",
	"foo()",
	"bigint()",
	// empty / malformed ROW
	"ROW()",
	"ROW(a bigint,)",           // trailing comma
	"ROW(a b c)",               // three bare tokens: a name + type leaves a dangling token
	"ROW(a varchar(10) extra)", // parameterized field type followed by a stray token
	"ROW(a bigint, 1)",         // named field mixed with an integer fits neither ROW reading
	"ROW(1, a bigint)",         // integer mixed with a named field
	"ROW(x bigint, y 5)",       // 'y 5' is neither a rowField nor a typeParameter
	"ROW(varchar(10) x, 5)",    // a named field forces rowType, which rejects the integer
	// D1: reversed / cross-family interval ranges
	"interval month to year",
	"interval year to day",
	"interval day to year",
	"interval second to year",
	"interval second to minute",
	"interval to second", // TO not preceded by a field (INTERVAL '(' path? actually INTERVAL TO)
	// malformed angle brackets
	"MAP<varchar>",  // map needs two components
	"MAP<varchar,>", // missing value type
	"ARRAY<>",       // empty element
	"ROW<bigint>",   // ROW uses parens, not angle brackets
	// stray time-zone tokens
	"timestamp with zone",    // missing TIME
	"timestamp without zone", // missing TIME
	"double precision precision",
	// malformed postfix-array brackets: `[` requires exactly `[ INTEGER ]`
	"bigint ARRAY[]",    // empty brackets
	"bigint ARRAY[abc]", // non-integer dimension
	"bigint ARRAY[5",    // missing close bracket
	"bigint ARRAY[5 6]", // extra token in brackets
}

// wrapType embeds a type in the canonical type position for the oracle: a
// CAST target inside a SELECT (a bare `CAST(...)` is not a statement and Trino
// rejects it before reaching the type grammar).
func wrapType(typ string) string { return "SELECT CAST(NULL AS " + typ + ")" }

// TestDataType_AcceptCorpusParses verifies (oracle-free) that omni accepts
// every type in typeAcceptCorpus when wrapped in a SELECT CAST. This is the
// completeness smoke test; the oracle differential below proves the verdicts
// actually agree with Trino.
//
// Note: the wrapped statement reaches omni's CAST/SELECT parsing, which is NOT
// implemented in this node — so we exercise ParseDataType directly on the bare
// type instead, which is exactly the function this node ships.
func TestDataType_AcceptCorpusParses(t *testing.T) {
	for _, typ := range typeAcceptCorpus {
		t.Run(truncateName(typ), func(t *testing.T) {
			dt, errs := ParseDataType(typ)
			if len(errs) != 0 {
				t.Errorf("ParseDataType(%q) should accept, got errors: %v", typ, errs)
			}
			if dt == nil {
				t.Errorf("ParseDataType(%q) returned nil", typ)
			}
		})
	}
}

// TestDataType_RejectCorpusRejected verifies omni rejects every type in
// typeRejectCorpus (negative coverage). ParseDataType must surface at least one
// error for each.
func TestDataType_RejectCorpusRejected(t *testing.T) {
	for _, typ := range typeRejectCorpus {
		t.Run(truncateName(typ), func(t *testing.T) {
			_, errs := ParseDataType(typ)
			if len(errs) == 0 {
				t.Errorf("ParseDataType(%q) should reject, but accepted", typ)
			}
		})
	}
}

// TestDataType_OracleDifferential is the authoritative gate: for every type in
// both corpora, omni's ParseDataType accept/reject must equal Trino 481's
// accept/reject of `CAST(NULL AS <type>)`. Skipped when no oracle is reachable.
func TestDataType_OracleDifferential(t *testing.T) {
	if testing.Short() {
		t.Skip("trino oracle: skipped in -short mode")
	}
	o := connectOracle(t)

	check := func(t *testing.T, typ string) {
		_, errs := ParseDataType(typ)
		omniAccepts := len(errs) == 0

		trinoAccepts, ok := oracleAccepts(t, o, wrapType(typ))
		if !ok {
			t.Skip("oracle unreachable for this case")
		}
		if omniAccepts != trinoAccepts {
			t.Errorf("MISMATCH type=%q: omni accepts=%v (errs=%v), Trino accepts=%v",
				typ, omniAccepts, errs, trinoAccepts)
		}
	}

	t.Run("accept", func(t *testing.T) {
		for _, typ := range typeAcceptCorpus {
			typ := typ
			t.Run(truncateName(typ), func(t *testing.T) { check(t, typ) })
		}
	})
	t.Run("reject", func(t *testing.T) {
		for _, typ := range typeRejectCorpus {
			typ := typ
			t.Run(truncateName(typ), func(t *testing.T) { check(t, typ) })
		}
	})
}

// ---------------------------------------------------------------------------
// Layer 3 — deparse round-trip (structural gate)
// ---------------------------------------------------------------------------

// TestDataType_RoundTrip parses each accepted type, renders it via String(),
// re-parses the rendering, and asserts the re-parse also succeeds and yields
// the same top-level Kind. This is the structural correctness gate in the
// absence of a reference type parser.
func TestDataType_RoundTrip(t *testing.T) {
	for _, typ := range typeAcceptCorpus {
		t.Run(truncateName(typ), func(t *testing.T) {
			dt, errs := ParseDataType(typ)
			if len(errs) != 0 {
				t.Fatalf("first parse of %q failed: %v", typ, errs)
			}
			rendered := dt.String()
			dt2, errs2 := ParseDataType(rendered)
			if len(errs2) != 0 {
				t.Fatalf("re-parse of rendered %q (from %q) failed: %v", rendered, typ, errs2)
			}
			if dt2.Kind != dt.Kind {
				t.Errorf("round-trip changed Kind: %q -> %q (%v -> %v)", typ, rendered, dt.Kind, dt2.Kind)
			}
		})
	}
}

// TestDataType_RoundTripOracle verifies the String() rendering of every
// accepted type is itself accepted by Trino — proving deparse emits valid Trino
// type syntax. Skipped without an oracle.
func TestDataType_RoundTripOracle(t *testing.T) {
	if testing.Short() {
		t.Skip("trino oracle: skipped in -short mode")
	}
	o := connectOracle(t)
	for _, typ := range typeAcceptCorpus {
		typ := typ
		t.Run(truncateName(typ), func(t *testing.T) {
			dt, errs := ParseDataType(typ)
			if len(errs) != 0 {
				t.Fatalf("parse of %q failed: %v", typ, errs)
			}
			rendered := dt.String()
			accepted, ok := oracleAccepts(t, o, wrapType(rendered))
			if !ok {
				t.Skip("oracle unreachable for this case")
			}
			if !accepted {
				t.Errorf("deparse produced a type Trino rejects: %q -> %q", typ, rendered)
			}
		})
	}
}

// TestValidIntervalRange unit-checks the D1 qualifier table directly (the
// oracle differential validates it end-to-end, but a focused unit test pins the
// exact rule so a regression is obvious).
func TestValidIntervalRange(t *testing.T) {
	valid := [][2]IntervalField{
		{IntervalYear, IntervalMonth},
		{IntervalDay, IntervalHour},
		{IntervalDay, IntervalMinute},
		{IntervalDay, IntervalSecond},
		{IntervalHour, IntervalMinute},
		{IntervalHour, IntervalSecond},
		{IntervalMinute, IntervalSecond},
	}
	invalid := [][2]IntervalField{
		{IntervalMonth, IntervalYear}, // reversed
		{IntervalYear, IntervalDay},   // cross-family
		{IntervalDay, IntervalYear},   // cross-family reversed
		{IntervalSecond, IntervalMinute},
		{IntervalYear, IntervalYear}, // not strictly coarser
		{IntervalDay, IntervalDay},
	}
	for _, p := range valid {
		if !ValidIntervalRange(p[0], p[1]) {
			t.Errorf("ValidIntervalRange(%v, %v) = false, want true", p[0], p[1])
		}
	}
	for _, p := range invalid {
		if ValidIntervalRange(p[0], p[1]) {
			t.Errorf("ValidIntervalRange(%v, %v) = true, want false", p[0], p[1])
		}
	}
}
