package parser

// simpleTypenameLeadTokens is the single source of truth for tokens that
// unambiguously start a SimpleTypename production (excluding GenericType
// which is covered by the type_function_name category check).
//
// Grammar reference: SimpleTypename in postgres/src/backend/parser/gram.y:14339.
//
// DO NOT extend this list without also extending parseSimpleTypename's
// dispatch switch in type.go, and vice versa. The oracle test
// TestSimpleTypenameLeadTokensMatchPG enforces this invariant against
// PG 17.
var simpleTypenameLeadTokens = []int{
	INT_P, INTEGER, SMALLINT, BIGINT, REAL, FLOAT_P, DOUBLE_P,
	DECIMAL_P, DEC, NUMERIC,
	BIT, CHARACTER, CHAR_P, VARCHAR, NATIONAL, NCHAR,
	BOOLEAN_P, JSON,
	TIMESTAMP, TIME, INTERVAL,
}

var simpleTypenameLeadSet = func() map[int]bool {
	m := make(map[int]bool, len(simpleTypenameLeadTokens))
	for _, t := range simpleTypenameLeadTokens {
		m[t] = true
	}
	return m
}()

// isSimpleTypenameStart reports whether the current token can start a
// SimpleTypename production. This is the single source of truth for all
// SimpleTypename FIRST-set probes. Callers must use this function rather
// than hand-writing token case clusters.
func (p *Parser) isSimpleTypenameStart() bool {
	if simpleTypenameLeadSet[p.cur.Type] {
		return true
	}
	// GenericType fallthrough: type_function_name opt_type_modifiers
	return p.isTypeFunctionName()
}

// isConstTypenameStart reports whether the current token starts a
// type-cast literal at PG's AexprConst grammar position. This is the
// FIRST set used by isAExprConstTypeCast / parseTypeCastedConst, which
// handles the union of two PG productions:
//
//	ConstTypename:  Numeric | ConstBit | ConstCharacter | ConstDatetime | JsonType  (gram.y:14370)
//	ConstInterval:  INTERVAL                                                          (gram.y:14688)
//
// Strictly speaking, INTERVAL is in ConstInterval, not ConstTypename —
// PG's AexprConst rule has both `ConstTypename Sconst` and
// `ConstInterval Sconst` as separate alternatives. omni's parser routes
// both alternatives through the same parseTypeCastedConst entry point
// (see expr.go), so the predicate's FIRST set is the union. Hence the
// "Const" in the name refers to the broader "type-cast literal lead"
// rather than the strict ConstTypename production.
//
// In omni this union coincides exactly with simpleTypenameLeadSet
// (SimpleTypename = the union plus GenericType fallback plus a couple
// of historical bits), so we reuse simpleTypenameLeadSet rather than
// maintaining a second slice. The only behavioral difference between
// isSimpleTypenameStart and isConstTypenameStart is the GenericType /
// type_function_name fallthrough — SimpleTypename has it, this predicate
// does not, because typed-string-literal contexts do not accept
// arbitrary identifier-typed names.
//
// NOTE: this predicate CANNOT be validated via a PG oracle at the
// AexprConst grammar position because that position is ambiguous with
// `func_name Sconst`, which accepts almost any keyword. See
// TestIsConstTypenameStartRejectsNonTypeStarters for the negative
// coverage that replaces the oracle test.
func (p *Parser) isConstTypenameStart() bool {
	return simpleTypenameLeadSet[p.cur.Type]
}

// isTypenameStart reports whether the current token starts a Typename
// production.
//
//	Typename:
//	    SimpleTypename opt_array_bounds
//	    | SETOF SimpleTypename opt_array_bounds
//	    | SimpleTypename ARRAY '[' Iconst ']'
//	    | SETOF SimpleTypename ARRAY '[' Iconst ']'
//	    | SimpleTypename ARRAY
//	    | SETOF SimpleTypename ARRAY
//
// Composed from isSimpleTypenameStart plus SETOF. There is no
// typenameLeadTokens slice — the FIRST set is fully expressible as
// {SETOF} ∪ FIRST(SimpleTypename), and we reuse the latter.
//
// At the FIRST-set level, FIRST(Typename) equals FIRST(func_type),
// because func_type's %TYPE alternatives all start with
// type_function_name, which is already in FIRST(Typename) via
// SimpleTypename's GenericType path. The oracle test below probes the
// RETURNS position of CREATE FUNCTION, which is grammatically func_type,
// and that's a sound oracle for Typename's FIRST set.
func (p *Parser) isTypenameStart() bool {
	if p.cur.Type == SETOF {
		return true
	}
	return p.isSimpleTypenameStart()
}
