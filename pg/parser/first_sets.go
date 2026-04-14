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
// ConstTypename production.
//
// Grammar reference: gram.y:14370
//
//	ConstTypename: Numeric | ConstBit | ConstCharacter | ConstDatetime | JsonType
//
// In omni's parser, ConstTypename and SimpleTypename share the SAME
// hard-lead token set — both are validated by parseSimpleTypename's
// switch at type.go:105. The difference is that SimpleTypename additionally
// falls through to `isTypeFunctionName()` (GenericType path) while
// ConstTypename does not. We therefore reuse simpleTypenameLeadSet and
// simply skip the fallback.
//
// NOTE: this predicate CANNOT be validated via a PG oracle at the
// AexprConst grammar position because that position is ambiguous with
// `func_name Sconst`, which accepts almost any keyword. See
// TestIsConstTypenameStartRejectsNonTypeStarters for the hand-curated
// negative coverage that replaces the oracle test.
func (p *Parser) isConstTypenameStart() bool {
	return simpleTypenameLeadSet[p.cur.Type]
}
