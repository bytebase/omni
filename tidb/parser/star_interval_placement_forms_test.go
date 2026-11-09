package parser

import "testing"

// Placement of `*` and of a bare `INTERVAL expr unit` in the core expression
// parser. Neither is an Expression in TiDB's grammar:
//
//   - `*` (and `t.*`, `s.t.*`) is a Field wildcard (parser.y:6464-6473), valid
//     only in a select list and in COUNT's own `'(' '*' ')'` production
//     (parser.y:8546). It is not a SimpleExpr, so it is rejected as a
//     parenthesized expression, an IN value, a row element, or a call argument.
//
//   - `INTERVAL expr unit` is only an operand of the binary +/- date-arithmetic
//     productions — `BitExpr '+' "INTERVAL" Expression TimeUnit`,
//     `BitExpr '-' "INTERVAL" ...` and `"INTERVAL" Expression TimeUnit '+' BitExpr`
//     (parser.y:7793, 7804, 7815) — and the second argument of the date-arith
//     builtins. It is not a SimpleExpr, so it is rejected as the operand of a
//     unary prefix operator, as a row element, or as a bare/parenthesized
//     expression.
//
// The five corpus rows these fixes flip (DO(*), DO ident IN(*),
// DO +/- INTERVAL ..., DO(ident, INTERVAL ...)) carry the engine verdict measured
// against pingcap/tidb:v8.5.5 in parser_y_negatives_corpus.tsv. The cases below
// are grammar-derived and held to omni's parser; accept records what TiDB's
// grammar allows.

type starIntervalPlacementCase struct {
	sql    string
	accept bool
}

var starIntervalPlacementCases = []starIntervalPlacementCase{
	// Mechanism A: `*` is not a general expression. Rejected off the select list.
	{`do (*)`, false},
	{`do (t.*)`, false},
	{`do (s.t.*)`, false},
	{`do a in (*)`, false},
	{`do a in (t.*)`, false},
	{`do a in (1, *)`, false},
	{`select (*)`, false},
	{`select 1 in (*) from t`, false},
	{`do (a, *)`, false},

	// Mechanism A: the valid homes for `*` still parse.
	{`select * from t`, true},
	{`select t.* from t`, true},
	{`select s.t.* from t`, true},
	{`select count(*) from t`, true},
	{`select count(*)`, true},
	{`do a in (1, 2)`, true},
	{`select c in (1, 2) from t`, true},
	{`do (a)`, true},
	{`do (a, b)`, true},

	// Mechanism B: a bare INTERVAL is not the operand of a unary prefix operator.
	{`do - interval 1 day`, false},
	{`do + interval 1 day`, false},
	{`do - interval a second_microsecond`, false},
	{`do + interval a second_microsecond`, false},
	{`do ~ interval 1 day`, false},
	{`do ! interval 1 day`, false},
	{`do binary interval 1 day`, false},

	// Mechanism B: a bare INTERVAL is not a row element.
	{`do (a, interval 1 day)`, false},
	{`do (interval 1 day, a)`, false},
	{`select (a, interval 1 day) from t`, false},
	{`do row(a, interval 1 day)`, false},

	// Deferred, not fixed here: a bare or parenthesized INTERVAL standing as a
	// whole expression (`do interval 1 day`, `do (interval 1 day)`) is also a 1064
	// on v8.5.5 but reaches the parser through parseDoStmt and the single-expression
	// parenthesized form, outside the unary/row/IN/paren-star sites this PR narrows.
	// Those over-acceptances are left for a follow-up.

	// Mechanism B: the valid homes for INTERVAL still parse. It is an operand of
	// the binary +/- date-arithmetic forms (either side) and of the date-arith
	// builtins, and those forms nest inside parens, rows and unary operators.
	{`do d + interval 1 day`, true},
	{`do d - interval 1 day`, true},
	{`do interval 1 day + d`, true},
	{`do 1 + interval 1 day`, true},
	{`do (d + interval 1 day)`, true},
	{`do (a, d + interval 1 day)`, true},
	{`do -(d + interval 1 day)`, true},
	{`select date_add(d, interval 1 day) from t`, true},
	{`select date_sub(d, interval 1 day) from t`, true},
	{`do d + interval 1 day + interval 2 hour`, true},
}

func TestTiDBStarIntervalPlacementForms(t *testing.T) {
	for _, c := range starIntervalPlacementCases {
		t.Run(c.sql, func(t *testing.T) {
			list, err := Parse(c.sql)
			if got := err == nil; got != c.accept {
				t.Fatalf("Parse(%q) accepts = %v, want %v (err: %v)", c.sql, got, c.accept, err)
			}
			if err == nil && len(list.Items) != 1 {
				t.Errorf("Parse(%q) read %d statements, want 1", c.sql, len(list.Items))
			}
		})
	}
}
