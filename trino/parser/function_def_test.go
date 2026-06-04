package parser

import (
	"testing"
)

// This file is the structural (Layer 2) gate for the parser-routines node's
// statement shells — CREATE FUNCTION, DROP FUNCTION, and the WITH FUNCTION
// inline-routine prefix (function_def.go). Accept/reject alone does not catch a
// form that "accepts" but parses into the wrong node shape (OR REPLACE flag, the
// parameter name/type split, the characteristic list, IF EXISTS, the inline
// function count). These tests pin the parse-node fields of one representative
// per alternative. The authoritative accept/reject differential against the live
// Trino 481 oracle lives in oracle_routines_test.go.

func TestRoutines_StructureCreateFunction(t *testing.T) {
	t.Run("minimal_return", func(t *testing.T) {
		s, ok := parseOneStmt(t, "CREATE FUNCTION meaning_of_life() RETURNS bigint RETURN 42").(*CreateFunctionStmt)
		if !ok {
			t.Fatalf("not a *CreateFunctionStmt")
		}
		if s.OrReplace {
			t.Errorf("OrReplace = true, want false")
		}
		if got := s.Spec.Declaration.Name.Normalize(); got != "meaning_of_life" {
			t.Errorf("name = %q, want meaning_of_life", got)
		}
		if n := len(s.Spec.Declaration.Parameters); n != 0 {
			t.Errorf("params = %d, want 0", n)
		}
		if s.Spec.ReturnType == nil || s.Spec.ReturnType.Name != "bigint" {
			t.Errorf("return type = %v, want bigint", s.Spec.ReturnType)
		}
		if _, ok := s.Spec.Body.(*ReturnStatement); !ok {
			t.Errorf("body = %T, want *ReturnStatement", s.Spec.Body)
		}
	})

	t.Run("or_replace_qualified_name", func(t *testing.T) {
		s := parseOneStmt(t, "CREATE OR REPLACE FUNCTION example.default.f(x bigint) RETURNS bigint RETURN x * x").(*CreateFunctionStmt)
		if !s.OrReplace {
			t.Errorf("OrReplace = false, want true")
		}
		if got := s.Spec.Declaration.Name.Normalize(); got != "example.default.f" {
			t.Errorf("name = %q, want example.default.f", got)
		}
		if n := len(s.Spec.Declaration.Parameters); n != 1 {
			t.Fatalf("params = %d, want 1", n)
		}
		p := s.Spec.Declaration.Parameters[0]
		if p.Name == nil || p.Name.Normalize() != "x" {
			t.Errorf("param[0] name = %v, want x", p.Name)
		}
		if p.Type == nil || p.Type.Name != "bigint" {
			t.Errorf("param[0] type = %v, want bigint", p.Type)
		}
	})

	t.Run("unnamed_and_mixed_parameters", func(t *testing.T) {
		// F2: parameterDeclaration is `identifier? type`. A type-only parameter
		// has a nil Name; a named one carries it.
		s := parseOneStmt(t, "CREATE FUNCTION f(x bigint, varchar, z double) RETURNS bigint RETURN 1").(*CreateFunctionStmt)
		params := s.Spec.Declaration.Parameters
		if len(params) != 3 {
			t.Fatalf("params = %d, want 3", len(params))
		}
		if params[0].Name == nil || params[0].Name.Normalize() != "x" {
			t.Errorf("param[0] name = %v, want x", params[0].Name)
		}
		if params[1].Name != nil {
			t.Errorf("param[1] name = %v, want nil (type-only)", params[1].Name)
		}
		if params[1].Type == nil || params[1].Type.Name != "varchar" {
			t.Errorf("param[1] type = %v, want varchar", params[1].Type)
		}
		if params[2].Name == nil || params[2].Name.Normalize() != "z" {
			t.Errorf("param[2] name = %v, want z", params[2].Name)
		}
	})

	t.Run("multitoken_typeonly_parameters", func(t *testing.T) {
		// Two tokens of lookahead cannot distinguish a type-only multi-token type
		// (DOUBLE PRECISION, INTERVAL DAY TO SECOND) from a named parameter; the
		// speculative resolution must read them as a single unnamed type (nil
		// Name), not as name + truncated type.
		for _, tc := range []struct {
			sql      string
			wantType string
		}{
			{"CREATE FUNCTION f(double precision) RETURNS int RETURN 1", "DOUBLE PRECISION"},
			{"CREATE FUNCTION f(interval day to second) RETURNS int RETURN 1", "INTERVAL"},
		} {
			s := parseOneStmt(t, tc.sql).(*CreateFunctionStmt)
			ps := s.Spec.Declaration.Parameters
			if len(ps) != 1 {
				t.Fatalf("%q: params = %d, want 1", tc.sql, len(ps))
			}
			if ps[0].Name != nil {
				t.Errorf("%q: param name = %v, want nil (type-only)", tc.sql, ps[0].Name)
			}
			if ps[0].Type == nil || ps[0].Type.Name != tc.wantType {
				t.Errorf("%q: param type = %v, want %s", tc.sql, ps[0].Type, tc.wantType)
			}
		}
	})

	t.Run("nonreserved_keyword_parameter_name", func(t *testing.T) {
		// A non-reserved keyword (COMMENT) is a valid parameter NAME when a type
		// follows it.
		s := parseOneStmt(t, "CREATE FUNCTION f(comment bigint) RETURNS int RETURN 1").(*CreateFunctionStmt)
		ps := s.Spec.Declaration.Parameters
		if len(ps) != 1 || ps[0].Name == nil || ps[0].Name.Normalize() != "comment" {
			t.Errorf("param name = %v, want comment", ps[0].Name)
		}
		if ps[0].Type == nil || ps[0].Type.Name != "bigint" {
			t.Errorf("param type = %v, want bigint", ps[0].Type)
		}
	})

	t.Run("all_keyword_characteristics", func(t *testing.T) {
		// F3: a representative full characteristic list (each keyword-only form
		// plus LANGUAGE and COMMENT) parses into the kinds in source order.
		s := parseOneStmt(t, "CREATE FUNCTION f(x bigint) RETURNS bigint LANGUAGE SQL DETERMINISTIC RETURNS NULL ON NULL INPUT SECURITY DEFINER COMMENT 'doc' RETURN x").(*CreateFunctionStmt)
		chs := s.Spec.Characteristics
		wantKinds := []RoutineCharacteristicKind{
			CharacteristicLanguage,
			CharacteristicDeterministic,
			CharacteristicReturnsNullOnNullInput,
			CharacteristicSecurityDefiner,
			CharacteristicComment,
		}
		if len(chs) != len(wantKinds) {
			t.Fatalf("characteristics = %d, want %d", len(chs), len(wantKinds))
		}
		for i, want := range wantKinds {
			if chs[i].Kind != want {
				t.Errorf("characteristic[%d] kind = %v, want %v", i, chs[i].Kind, want)
			}
		}
		if chs[0].Language == nil || chs[0].Language.Normalize() != "sql" {
			t.Errorf("LANGUAGE id = %v, want sql", chs[0].Language)
		}
	})

	t.Run("not_deterministic_called_invoker", func(t *testing.T) {
		s := parseOneStmt(t, "CREATE FUNCTION f(x bigint) RETURNS bigint NOT DETERMINISTIC CALLED ON NULL INPUT SECURITY INVOKER RETURN x").(*CreateFunctionStmt)
		wantKinds := []RoutineCharacteristicKind{
			CharacteristicNotDeterministic,
			CharacteristicCalledOnNullInput,
			CharacteristicSecurityInvoker,
		}
		chs := s.Spec.Characteristics
		if len(chs) != len(wantKinds) {
			t.Fatalf("characteristics = %d, want %d", len(chs), len(wantKinds))
		}
		for i, want := range wantKinds {
			if chs[i].Kind != want {
				t.Errorf("characteristic[%d] kind = %v, want %v", i, chs[i].Kind, want)
			}
		}
	})

	t.Run("properties_characteristic_docs_ahead", func(t *testing.T) {
		// F4: WITH ( name = expr, ... ) — a docs-ahead-of-legacy characteristic
		// (flagged divergence). Confirm it parses into a properties characteristic
		// with the named entries.
		s := parseOneStmt(t, "CREATE FUNCTION f() RETURNS int LANGUAGE PYTHON WITH (handler = 'x', other_prop = 1) RETURN 1").(*CreateFunctionStmt)
		chs := s.Spec.Characteristics
		if len(chs) != 2 {
			t.Fatalf("characteristics = %d, want 2", len(chs))
		}
		if chs[1].Kind != CharacteristicProperties {
			t.Fatalf("characteristic[1] kind = %v, want Properties", chs[1].Kind)
		}
		props := chs[1].Properties
		if len(props) != 2 {
			t.Fatalf("properties = %d, want 2", len(props))
		}
		if props[0].Name.Normalize() != "handler" {
			t.Errorf("property[0] name = %q, want handler", props[0].Name.Normalize())
		}
		if props[1].Name.Normalize() != "other_prop" {
			t.Errorf("property[1] name = %q, want other_prop", props[1].Name.Normalize())
		}
	})

	t.Run("compound_body", func(t *testing.T) {
		s := parseOneStmt(t, "CREATE FUNCTION f(x bigint) RETURNS bigint BEGIN DECLARE y bigint DEFAULT 0; SET y = x + 1; RETURN y; END").(*CreateFunctionStmt)
		body, ok := s.Spec.Body.(*CompoundStatement)
		if !ok {
			t.Fatalf("body = %T, want *CompoundStatement", s.Spec.Body)
		}
		if len(body.Declarations) != 1 {
			t.Errorf("declarations = %d, want 1", len(body.Declarations))
		}
		if len(body.Body) != 2 {
			t.Errorf("body statements = %d, want 2", len(body.Body))
		}
	})
}

func TestRoutines_StructureDropFunction(t *testing.T) {
	t.Run("with_param_types", func(t *testing.T) {
		s, ok := parseOneStmt(t, "DROP FUNCTION multiply_by_two(bigint)").(*DropFunctionStmt)
		if !ok {
			t.Fatalf("not a *DropFunctionStmt")
		}
		if s.IfExists {
			t.Errorf("IfExists = true, want false")
		}
		if got := s.Declaration.Name.Normalize(); got != "multiply_by_two" {
			t.Errorf("name = %q, want multiply_by_two", got)
		}
		if len(s.Declaration.Parameters) != 1 {
			t.Fatalf("params = %d, want 1", len(s.Declaration.Parameters))
		}
		if s.Declaration.Parameters[0].Type.Name != "bigint" {
			t.Errorf("param[0] type = %v, want bigint", s.Declaration.Parameters[0].Type)
		}
	})

	t.Run("if_exists_qualified", func(t *testing.T) {
		s := parseOneStmt(t, "DROP FUNCTION IF EXISTS example.default.f(integer, varchar)").(*DropFunctionStmt)
		if !s.IfExists {
			t.Errorf("IfExists = false, want true")
		}
		if got := s.Declaration.Name.Normalize(); got != "example.default.f" {
			t.Errorf("name = %q, want example.default.f", got)
		}
		if len(s.Declaration.Parameters) != 2 {
			t.Errorf("params = %d, want 2", len(s.Declaration.Parameters))
		}
	})

	t.Run("no_params", func(t *testing.T) {
		s := parseOneStmt(t, "DROP FUNCTION meaning_of_life()").(*DropFunctionStmt)
		if len(s.Declaration.Parameters) != 0 {
			t.Errorf("params = %d, want 0", len(s.Declaration.Parameters))
		}
	})
}

func TestRoutines_StructureWithFunction(t *testing.T) {
	t.Run("single_inline_function", func(t *testing.T) {
		s, ok := parseOneStmt(t, "WITH FUNCTION f(x bigint) RETURNS bigint RETURN x + 1 SELECT f(10)").(*WithFunctionStmt)
		if !ok {
			t.Fatalf("not a *WithFunctionStmt")
		}
		if len(s.Functions) != 1 {
			t.Fatalf("functions = %d, want 1", len(s.Functions))
		}
		if got := s.Functions[0].Declaration.Name.Normalize(); got != "f" {
			t.Errorf("function name = %q, want f", got)
		}
		if s.Query == nil {
			t.Errorf("Query = nil, want a parsed query")
		}
	})

	t.Run("multiple_inline_functions", func(t *testing.T) {
		s := parseOneStmt(t, "WITH FUNCTION a(x int) RETURNS int RETURN x, FUNCTION b(y int) RETURNS int RETURN y SELECT a(1), b(2)").(*WithFunctionStmt)
		if len(s.Functions) != 2 {
			t.Fatalf("functions = %d, want 2", len(s.Functions))
		}
		if s.Functions[0].Declaration.Name.Normalize() != "a" {
			t.Errorf("function[0] = %q, want a", s.Functions[0].Declaration.Name.Normalize())
		}
		if s.Functions[1].Declaration.Name.Normalize() != "b" {
			t.Errorf("function[1] = %q, want b", s.Functions[1].Declaration.Name.Normalize())
		}
	})

	t.Run("inline_function_then_cte", func(t *testing.T) {
		// F6: a CTE WITH may still follow the inline routines.
		s := parseOneStmt(t, "WITH FUNCTION f(x int) RETURNS int RETURN x WITH t AS (SELECT 1 AS c) SELECT f(c) FROM t").(*WithFunctionStmt)
		if len(s.Functions) != 1 {
			t.Fatalf("functions = %d, want 1", len(s.Functions))
		}
		if s.Query == nil || s.Query.With == nil {
			t.Errorf("trailing query CTE = nil, want a WITH clause")
		}
	})
}
