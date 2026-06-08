package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
)

// ---------------------------------------------------------------------------
// GQL — graph query language (parser-gql node)
// ---------------------------------------------------------------------------
//
// This file ports the legacy ANTLR §2.12 GQL sub-language
// (GoogleSQLParser.g4 gql_statement and its graph_* / label_expression rules) —
// a hand-port of ZetaSQL's GoogleSQL graph extension (BigQuery / Spanner Graph).
// The top-level GQL statement is `GRAPH <path> <graph_operation_block>`, a
// NEXT-separated list of composite query blocks; each composite block is either
// a single linear query operation or a set-operation-joined chain. A linear
// query operation is a sequence of linear operators (MATCH / OPTIONAL MATCH /
// LET / FILTER / ORDER BY / PAGE / WITH / FOR / TABLESAMPLE) terminated by a
// required RETURN operator. The graph pattern is a path-pattern list of
// node `( … )` and edge `-[ … ]->` patterns over a shared element-pattern
// filler (variable + label algebra + property spec + inline WHERE).
//
// ORACLE NOTE — GQL is BigQuery/Spanner-Graph syntax. The Spanner emulator's
// accept/reject is recorded by the differential (graph_query_oracle_test.go) and
// used where it does not disagree with the legacy .g4; the truth1 BigQuery corpus
// explicitly scopes the full GQL syntax OUT (INDEX.md), so the authoritative
// reference for the grammar shape is the pinned legacy GoogleSQLParser.g4 itself.
// One omni parser serves both dialects — it accepts the BigQuery+Spanner union.

// parseGQLStmt parses a gql_statement:
//
//	GRAPH <path_expression> <graph_operation_block>
//
// The GRAPH keyword has NOT been consumed (parseStmt peeks it). The graph
// operation block is a NEXT-separated list of composite query blocks.
func (p *Parser) parseGQLStmt() (ast.Node, error) {
	graph := p.advance() // consume GRAPH

	stmt := &ast.GQLStmt{}
	stmt.Loc.Start = graph.Loc.Start

	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// graph_operation_block: composite_query_block (NEXT composite_query_block)*.
	for {
		block, err := p.parseGraphCompositeQueryBlock()
		if err != nil {
			return nil, err
		}
		stmt.Blocks = append(stmt.Blocks, block)
		stmt.Loc.End = nodeLoc(block).End
		if p.cur.Type != kwNEXT {
			break
		}
		p.advance() // NEXT
	}
	return stmt, nil
}

// parseGraphCompositeQueryBlock parses a graph_composite_query_block: a single
// linear query operation, or a graph_composite_query_prefix (a set-operation-
// joined chain of >= 2 linear query operations). Because both alternatives begin
// with a linear query operation, we parse one and then check for a trailing
// set-operation metadata token to decide which alternative we are in.
func (p *Parser) parseGraphCompositeQueryBlock() (ast.Node, error) {
	first, err := p.parseGraphLinearQuery()
	if err != nil {
		return nil, err
	}
	if !p.atGraphSetOpMeta() {
		return first, nil
	}
	// graph_composite_query_prefix: linear (meta linear)+ — at least one more.
	setOp := &ast.GraphSetOp{Ops: []ast.Node{first}, Loc: first.Loc}
	for p.atGraphSetOpMeta() {
		meta, err := p.parseGraphSetOpMeta()
		if err != nil {
			return nil, err
		}
		next, err := p.parseGraphLinearQuery()
		if err != nil {
			return nil, err
		}
		setOp.Metas = append(setOp.Metas, meta)
		setOp.Ops = append(setOp.Ops, next)
		setOp.Loc.End = next.Loc.End
	}
	return setOp, nil
}

// atGraphSetOpMeta reports whether the cursor is at the start of a
// graph_set_operation_metadata (query_set_operation_type all_or_distinct): one
// of UNION / EXCEPT / INTERSECT.
func (p *Parser) atGraphSetOpMeta() bool {
	switch p.cur.Type {
	case kwUNION, kwEXCEPT, kwINTERSECT:
		return true
	}
	return false
}

// parseGraphSetOpMeta parses a graph_set_operation_metadata:
// query_set_operation_type all_or_distinct. Both the operation type
// (UNION|EXCEPT|INTERSECT) and the quantifier (ALL|DISTINCT) are REQUIRED.
func (p *Parser) parseGraphSetOpMeta() (*ast.GraphSetOpMeta, error) {
	opTok := p.advance() // UNION | EXCEPT | INTERSECT (atGraphSetOpMeta gated)
	meta := &ast.GraphSetOpMeta{Op: TokenName(opTok.Type), Loc: opTok.Loc}
	switch p.cur.Type {
	case kwALL:
		meta.Quantifier = "ALL"
	case kwDISTINCT:
		meta.Quantifier = "DISTINCT"
	default:
		return nil, p.syntaxErrorAtCur()
	}
	q := p.advance()
	meta.Loc.End = q.Loc.End
	return meta, nil
}

// parseGraphLinearQuery parses a graph_linear_query_operation:
//
//	graph_linear_operator_list? graph_return_operator
//
// The operator list is a sequence of linear operators; the RETURN operator is
// required at the end.
func (p *Parser) parseGraphLinearQuery() (*ast.GraphLinearQuery, error) {
	lq := &ast.GraphLinearQuery{Loc: p.cur.Loc}
	for p.atGraphLinearOperator() {
		op, err := p.parseGraphLinearOperator()
		if err != nil {
			return nil, err
		}
		lq.Operators = append(lq.Operators, op)
	}
	ret, err := p.parseGraphReturnOp()
	if err != nil {
		return nil, err
	}
	lq.Return = ret
	lq.Loc.End = ret.Loc.End
	if len(lq.Operators) > 0 {
		lq.Loc.Start = nodeLoc(lq.Operators[0]).Start
	} else {
		lq.Loc.Start = ret.Loc.Start
	}
	return lq, nil
}

// atGraphLinearOperator reports whether the cursor is at the start of a
// graph_linear_operator (one of MATCH / OPTIONAL MATCH / LET / FILTER /
// ORDER BY / PAGE(OFFSET|SKIP|LIMIT) / WITH / FOR / TABLESAMPLE). RETURN is NOT
// a linear operator (it terminates the linear query), so it is excluded.
func (p *Parser) atGraphLinearOperator() bool {
	switch p.cur.Type {
	case kwMATCH, kwOPTIONAL, kwLET, kwFILTER, kwORDER, kwWITH, kwFOR,
		kwTABLESAMPLE, kwOFFSET, kwSKIP, kwLIMIT:
		return true
	}
	return false
}

// parseGraphLinearOperator parses one graph_linear_operator by dispatching on
// the leading keyword.
func (p *Parser) parseGraphLinearOperator() (ast.Node, error) {
	switch p.cur.Type {
	case kwMATCH:
		return p.parseGraphMatchOp(false)
	case kwOPTIONAL:
		return p.parseGraphMatchOp(true)
	case kwLET:
		return p.parseGraphLetOp()
	case kwFILTER:
		return p.parseGraphFilterOp()
	case kwORDER:
		return p.parseGraphOrderByOp()
	case kwOFFSET, kwSKIP, kwLIMIT:
		return p.parseGraphPageOp()
	case kwWITH:
		return p.parseGraphWithOp()
	case kwFOR:
		return p.parseGraphForOp()
	case kwTABLESAMPLE:
		return p.parseGraphSampleOp()
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseGraphMatchOp parses a graph_match_operator / graph_optional_match_operator:
//
//	[OPTIONAL] MATCH [hint] <graph_pattern>
//
// When optional is true the cursor is at OPTIONAL; otherwise it is at MATCH.
func (p *Parser) parseGraphMatchOp(optional bool) (ast.Node, error) {
	start := p.cur.Loc
	if optional {
		p.advance() // OPTIONAL
		if _, err := p.expect(kwMATCH); err != nil {
			return nil, err
		}
	} else {
		p.advance() // MATCH
	}
	if p.cur.Type == int('@') {
		if herr := p.skipHint(); herr != nil {
			return nil, herr
		}
	}
	pattern, err := p.parseGraphPattern()
	if err != nil {
		return nil, err
	}
	return &ast.GraphMatchOp{
		Optional: optional,
		Pattern:  pattern,
		Loc:      ast.Loc{Start: start.Start, End: pattern.Loc.End},
	}, nil
}

// parseGraphLetOp parses a graph_let_operator: `LET <var defs>`. Each definition
// is `identifier = expression`; >= 1.
func (p *Parser) parseGraphLetOp() (ast.Node, error) {
	let := p.advance() // LET
	op := &ast.GraphLetOp{Loc: let.Loc}
	for {
		v, err := p.parseGraphLetVar()
		if err != nil {
			return nil, err
		}
		op.Vars = append(op.Vars, v)
		op.Loc.End = v.Loc.End
		if p.cur.Type != int(',') {
			break
		}
		p.advance() // ,
	}
	return op, nil
}

// parseGraphLetVar parses one graph_let_variable_definition:
// identifier = expression.
func (p *Parser) parseGraphLetVar() (*ast.GraphLetVar, error) {
	tok, err := p.expectIdentifier()
	if err != nil {
		return nil, err
	}
	name, err := p.identifierText(tok)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(int('=')); err != nil {
		return nil, err
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return &ast.GraphLetVar{
		Name: name,
		Expr: expr,
		Loc:  ast.Loc{Start: tok.Loc.Start, End: nodeLoc(expr).End},
	}, nil
}

// parseGraphFilterOp parses a graph_filter_operator:
//
//	FILTER WHERE <expr>   |   FILTER <expr>
func (p *Parser) parseGraphFilterOp() (ast.Node, error) {
	filter := p.advance() // FILTER
	op := &ast.GraphFilterOp{Loc: filter.Loc}
	if p.cur.Type == kwWHERE {
		p.advance() // WHERE
		op.HasWhere = true
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	op.Expr = expr
	op.Loc.End = nodeLoc(expr).End
	return op, nil
}

// parseGraphOrderByOp parses a graph_order_by_operator (graph_order_by_clause):
//
//	ORDER [hint] BY <ordering_expression> (, …)
//
// Each ordering expression is `expr [COLLATE c] [ASC|DESC|ASCENDING|DESCENDING]
// [NULLS FIRST|LAST]`, reusing the shared OrderItem.
func (p *Parser) parseGraphOrderByOp() (ast.Node, error) {
	order := p.advance() // ORDER
	if p.cur.Type == int('@') {
		if herr := p.skipHint(); herr != nil {
			return nil, herr
		}
	}
	if _, err := p.expect(kwBY); err != nil {
		return nil, err
	}
	items, end, err := p.parseGraphOrderingList()
	if err != nil {
		return nil, err
	}
	return &ast.GraphOrderByOp{Items: items, Loc: ast.Loc{Start: order.Loc.Start, End: end}}, nil
}

// parseGraphOrderingList parses a comma-separated list of graph_ordering_expression
// into OrderItems and returns the end offset of the last item. Used by both the
// ORDER BY operator and the trailing ORDER BY of RETURN.
func (p *Parser) parseGraphOrderingList() ([]*ast.OrderItem, int, error) {
	var items []*ast.OrderItem
	end := p.cur.Loc.End
	for {
		item, err := p.parseGraphOrderingExpr()
		if err != nil {
			return nil, 0, err
		}
		items = append(items, item)
		end = item.Loc.End
		if p.cur.Type != int(',') {
			break
		}
		p.advance() // ,
	}
	return items, end, nil
}

// parseGraphOrderingExpr parses a graph_ordering_expression:
//
//	expression [COLLATE c] [ASC|DESC|ASCENDING|DESCENDING] [NULLS FIRST|LAST]
//
// It returns an OrderItem (the shared ordering element). GQL additionally allows
// the long ASCENDING/DESCENDING spellings (opt_graph_asc_or_desc) on top of the
// usual ASC/DESC.
func (p *Parser) parseGraphOrderingExpr() (*ast.OrderItem, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	item := &ast.OrderItem{Expr: expr, Loc: ast.Loc{Start: nodeLoc(expr).Start, End: nodeLoc(expr).End}}
	if p.cur.Type == kwCOLLATE {
		p.advance() // COLLATE
		// collate_clause: COLLATE string_literal_or_parameter. Parse it as a Node
		// (bit-or precedence, matching the shared parseOrderItem) so the collation
		// operand is retained on the OrderItem for AST consumers, not just dropped.
		coll, err := p.parseBinaryExpr(bpBitOr)
		if err != nil {
			return nil, err
		}
		item.Collate = coll
		item.Loc.End = nodeLoc(coll).End
	}
	if dir, ok := p.matchGraphAscDesc(); ok {
		item.HasDir = true
		item.Desc = dir
		item.Loc.End = p.prev.Loc.End
	}
	if no := p.matchNullOrder(); no != "" {
		nf := no == "FIRST"
		item.NullsFirst = &nf
		item.Loc.End = p.prev.Loc.End
	}
	return item, nil
}

// matchGraphAscDesc consumes an opt_graph_asc_or_desc (ASC | DESC | ASCENDING |
// DESCENDING) if present and returns (isDesc, true), or (false, false) when
// absent.
func (p *Parser) matchGraphAscDesc() (bool, bool) {
	switch p.cur.Type {
	case kwASC, kwASCENDING:
		p.advance()
		return false, true
	case kwDESC, kwDESCENDING:
		p.advance()
		return true, true
	}
	return false, false
}

// parseGraphPageOp parses a graph_page_operator (graph_page_clause) appearing as
// a standalone linear operator:
//
//	(OFFSET|SKIP) <n> [LIMIT <n>]   |   LIMIT <n>
func (p *Parser) parseGraphPageOp() (ast.Node, error) {
	return p.parseGraphPageClause()
}

// parseGraphPageClause parses a graph_page_clause and returns it as a
// *GraphPageOp. Shared by the standalone PAGE operator and the trailing PAGE of
// a RETURN operator.
func (p *Parser) parseGraphPageClause() (*ast.GraphPageOp, error) {
	// The OFFSET/SKIP/LIMIT counts are possibly_cast_int_literal_or_parameter
	// (integer / @param / ? / @@sysvar / CAST(… AS …)) — NOT a full expression. A
	// bare identifier or an arithmetic expression is a syntax error (oracle: the
	// Spanner emulator rejects `LIMIT x` / `LIMIT 1+1`).
	op := &ast.GraphPageOp{Loc: p.cur.Loc}
	switch p.cur.Type {
	case kwOFFSET, kwSKIP:
		op.SkipIsSkip = p.cur.Type == kwSKIP
		p.advance() // OFFSET | SKIP
		skip, err := p.parseIntLiteralOrParameterOrCast()
		if err != nil {
			return nil, err
		}
		op.Skip = skip
		op.Loc.End = nodeLoc(skip).End
		if p.cur.Type == kwLIMIT {
			p.advance() // LIMIT
			lim, err := p.parseIntLiteralOrParameterOrCast()
			if err != nil {
				return nil, err
			}
			op.Limit = lim
			op.Loc.End = nodeLoc(lim).End
		}
	case kwLIMIT:
		p.advance() // LIMIT
		lim, err := p.parseIntLiteralOrParameterOrCast()
		if err != nil {
			return nil, err
		}
		op.Limit = lim
		op.Loc.End = nodeLoc(lim).End
	default:
		return nil, p.syntaxErrorAtCur()
	}
	return op, nil
}

// parseGraphWithOp parses a graph_with_operator:
//
//	WITH [ALL|DISTINCT] [hint] <return_item_list> [GROUP BY …]
func (p *Parser) parseGraphWithOp() (ast.Node, error) {
	with := p.advance() // WITH
	op := &ast.GraphWithOp{Loc: with.Loc}
	op.Quantifier = p.matchAllOrDistinct()
	if p.cur.Type == int('@') {
		if herr := p.skipHint(); herr != nil {
			return nil, herr
		}
	}
	items, end, err := p.parseGraphReturnItemList()
	if err != nil {
		return nil, err
	}
	op.Items = items
	op.Loc.End = end
	if p.cur.Type == kwGROUP {
		gb, err := p.parseGroupByClause()
		if err != nil {
			return nil, err
		}
		op.GroupBy = gb
		op.Loc.End = gb.Loc.End
	}
	return op, nil
}

// parseGraphForOp parses a graph_for_operator:
//
//	FOR <id> IN <expr> [WITH OFFSET [AS alias]]
func (p *Parser) parseGraphForOp() (ast.Node, error) {
	forTok := p.advance() // FOR
	op := &ast.GraphForOp{Loc: forTok.Loc}
	tok, err := p.expectIdentifier()
	if err != nil {
		return nil, err
	}
	name, err := p.identifierText(tok)
	if err != nil {
		return nil, err
	}
	op.Name = name
	if _, err := p.expect(kwIN); err != nil {
		return nil, err
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	op.Expr = expr
	op.Loc.End = nodeLoc(expr).End
	// opt_with_offset_and_alias_with_required_as: WITH OFFSET [AS alias].
	if p.cur.Type == kwWITH && p.peekNext().Type == kwOFFSET {
		p.advance()           // WITH
		offTok := p.advance() // OFFSET
		op.WithOffset = true
		op.Loc.End = offTok.Loc.End
		if p.cur.Type == kwAS {
			p.advance() // AS
			aliasTok, err := p.expectIdentifier()
			if err != nil {
				return nil, err
			}
			alias, err := p.identifierText(aliasTok)
			if err != nil {
				return nil, err
			}
			op.OffsetAlias = alias
			op.Loc.End = aliasTok.Loc.End
		}
	}
	return op, nil
}

// parseGraphSampleOp parses a graph_sample_clause:
//
//	TABLESAMPLE <method> ( <sample_size> ) [opt_graph_sample_clause_suffix]
//
// where sample_size is `<value> (ROWS|PERCENT) [PARTITION BY …]` and the suffix
// is one of `REPEATABLE(n)` / `WITH WEIGHT [REPEATABLE(n)]` /
// `WITH WEIGHT AS id [REPEATABLE(n)]`.
func (p *Parser) parseGraphSampleOp() (ast.Node, error) {
	ts := p.advance() // TABLESAMPLE
	op := &ast.GraphSampleOp{Loc: ts.Loc}
	methodTok, err := p.expectIdentifier()
	if err != nil {
		return nil, err
	}
	method, err := p.identifierText(methodTok)
	if err != nil {
		return nil, err
	}
	op.Method = method
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	// sample_size: sample_size_value sample_size_unit partition_by?
	// sample_size_value is possibly_cast_int_literal_or_parameter OR a
	// floating_point_literal — NOT a full expression. A bare identifier or an
	// arithmetic expression is a syntax error (oracle: the Spanner emulator
	// rejects `(foo ROWS)` / `(1+1 ROWS)`).
	size, err := p.parseGraphSampleSizeValue()
	if err != nil {
		return nil, err
	}
	op.Size = size
	switch p.cur.Type {
	case kwROWS:
		op.Unit = "ROWS"
	case kwPERCENT:
		op.Unit = "PERCENT"
	default:
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // ROWS | PERCENT
	// partition_by_clause_prefix_no_hint?  — PARTITION BY expr (, expr)* inside
	// the sample size. The expressions are parse-only here (the sample-size
	// partition is a rarely-used tail); consume them for parity.
	if p.cur.Type == kwPARTITION {
		p.advance() // PARTITION
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		for {
			if _, err := p.parseExpr(); err != nil {
				return nil, err
			}
			if p.cur.Type != int(',') {
				break
			}
			p.advance() // ,
		}
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	op.Loc.End = closeTok.Loc.End
	// opt_graph_sample_clause_suffix.
	if err := p.parseGraphSampleSuffix(op); err != nil {
		return nil, err
	}
	return op, nil
}

// parseGraphSampleSuffix consumes an opt_graph_sample_clause_suffix into op:
//
//	repeatable_clause
//	| WITH WEIGHT repeatable_clause?
//	| WITH WEIGHT AS identifier repeatable_clause?
func (p *Parser) parseGraphSampleSuffix(op *ast.GraphSampleOp) error {
	switch p.cur.Type {
	case kwREPEATABLE:
		return p.parseRepeatableInto(op)
	case kwWITH:
		p.advance() // WITH
		if _, err := p.expect(kwWEIGHT); err != nil {
			return err
		}
		op.WithWeight = true
		op.Loc.End = p.prev.Loc.End
		if p.cur.Type == kwAS {
			p.advance() // AS
			aliasTok, err := p.expectIdentifier()
			if err != nil {
				return err
			}
			alias, err := p.identifierText(aliasTok)
			if err != nil {
				return err
			}
			op.WeightAlias = alias
			op.Loc.End = aliasTok.Loc.End
		}
		if p.cur.Type == kwREPEATABLE {
			return p.parseRepeatableInto(op)
		}
	}
	return nil
}

// parseRepeatableInto parses a repeatable_clause `REPEATABLE ( <n> )` into
// op.Repeatable. The cursor is at REPEATABLE.
func (p *Parser) parseRepeatableInto(op *ast.GraphSampleOp) error {
	p.advance() // REPEATABLE
	if _, err := p.expect(int('(')); err != nil {
		return err
	}
	n, err := p.parseExpr()
	if err != nil {
		return err
	}
	op.Repeatable = n
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return err
	}
	op.Loc.End = closeTok.Loc.End
	return nil
}

// parseGraphIntLiteralOrParameter parses an int_literal_or_parameter:
//
//	integer_literal | parameter_expression | system_variable_expression
//
// This is the STRICT form (no CAST, no float) used for graph quantifier bounds
// (`{ [n,] m }`). It deliberately rejects a bare identifier, a float, a CAST, or
// any arithmetic expression — the legacy grammar's int_literal_or_parameter does
// not admit them, and the live Spanner emulator syntax-rejects e.g. `{a,b}` /
// `{1+1,3}`. (Distinct from parseIntLiteralOrParameterOrCast, the
// possibly_cast_int_literal_or_parameter superset used for page counts.)
func (p *Parser) parseGraphIntLiteralOrParameter() (ast.Node, error) {
	switch p.cur.Type {
	case tokInteger:
		tok := p.advance()
		return &ast.Literal{Kind: ast.LitInt, Value: tok.Str, Loc: tok.Loc}, nil
	case int('@'):
		return p.parseNamedParameter()
	case tokAtAt:
		return p.parseSystemVariable()
	case int('?'):
		tok := p.advance()
		return &ast.Parameter{Positional: true, Loc: tok.Loc}, nil
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseGraphSampleSizeValue parses a sample_size_value:
//
//	possibly_cast_int_literal_or_parameter | floating_point_literal
//
// i.e. the page-count operand set (integer / @param / ? / @@sysvar / CAST)
// PLUS a floating-point literal (TABLESAMPLE BERNOULLI (10.5 PERCENT)). It
// rejects a bare identifier or an arithmetic expression as the size (oracle: the
// Spanner emulator syntax-rejects `(foo ROWS)` / `(1+1 ROWS)`).
func (p *Parser) parseGraphSampleSizeValue() (ast.Node, error) {
	if p.cur.Type == tokFloat {
		tok := p.advance()
		return &ast.Literal{Kind: ast.LitFloat, Value: tok.Str, Loc: tok.Loc}, nil
	}
	return p.parseIntLiteralOrParameterOrCast()
}

// parseGraphReturnOp parses a graph_return_operator:
//
//	RETURN [hint] [ALL|DISTINCT] <return_item_list> [GROUP BY] [ORDER BY] [PAGE]
func (p *Parser) parseGraphReturnOp() (*ast.GraphReturnOp, error) {
	ret, err := p.expect(kwRETURN)
	if err != nil {
		return nil, err
	}
	op := &ast.GraphReturnOp{Loc: ret.Loc}
	if p.cur.Type == int('@') {
		if herr := p.skipHint(); herr != nil {
			return nil, herr
		}
	}
	op.Quantifier = p.matchAllOrDistinct()
	items, end, err := p.parseGraphReturnItemList()
	if err != nil {
		return nil, err
	}
	op.Items = items
	op.Loc.End = end
	if p.cur.Type == kwGROUP {
		gb, err := p.parseGroupByClause()
		if err != nil {
			return nil, err
		}
		op.GroupBy = gb
		op.Loc.End = gb.Loc.End
	}
	if p.cur.Type == kwORDER {
		p.advance() // ORDER
		if p.cur.Type == int('@') {
			if herr := p.skipHint(); herr != nil {
				return nil, herr
			}
		}
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		obItems, obEnd, err := p.parseGraphOrderingList()
		if err != nil {
			return nil, err
		}
		op.OrderBy = obItems
		op.Loc.End = obEnd
	}
	if p.atGraphPageClause() {
		page, err := p.parseGraphPageClause()
		if err != nil {
			return nil, err
		}
		op.Page = page
		op.Loc.End = page.Loc.End
	}
	return op, nil
}

// atGraphPageClause reports whether the cursor is at the start of a
// graph_page_clause (OFFSET | SKIP | LIMIT) — the trailing PAGE of a RETURN.
func (p *Parser) atGraphPageClause() bool {
	switch p.cur.Type {
	case kwOFFSET, kwSKIP, kwLIMIT:
		return true
	}
	return false
}

// matchAllOrDistinct consumes an all_or_distinct (ALL | DISTINCT) if present and
// returns "ALL"/"DISTINCT", or "" when absent. (In GQL the quantifier is
// optional on WITH/RETURN, distinct from a set-operation metadata where it is
// required.)
func (p *Parser) matchAllOrDistinct() string {
	switch p.cur.Type {
	case kwALL:
		p.advance()
		return "ALL"
	case kwDISTINCT:
		p.advance()
		return "DISTINCT"
	}
	return ""
}

// parseGraphReturnItemList parses a graph_return_item_list and returns the items
// plus the end offset of the last item. Each item is `expr [AS id]` or `*`.
func (p *Parser) parseGraphReturnItemList() ([]*ast.GraphReturnItem, int, error) {
	var items []*ast.GraphReturnItem
	end := p.cur.Loc.End
	for {
		item, err := p.parseGraphReturnItem()
		if err != nil {
			return nil, 0, err
		}
		items = append(items, item)
		end = item.Loc.End
		if p.cur.Type != int(',') {
			break
		}
		p.advance() // ,
	}
	return items, end, nil
}

// parseGraphReturnItem parses one graph_return_item:
//
//	<expr> [AS identifier]   |   *
func (p *Parser) parseGraphReturnItem() (*ast.GraphReturnItem, error) {
	if p.cur.Type == int('*') {
		star := p.advance()
		return &ast.GraphReturnItem{Star: true, Loc: star.Loc}, nil
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	item := &ast.GraphReturnItem{Expr: expr, Loc: ast.Loc{Start: nodeLoc(expr).Start, End: nodeLoc(expr).End}}
	if p.cur.Type == kwAS {
		p.advance() // AS
		tok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		alias, err := p.identifierText(tok)
		if err != nil {
			return nil, err
		}
		item.Alias = alias
		item.Loc.End = tok.Loc.End
	}
	return item, nil
}

// ---------------------------------------------------------------------------
// Graph patterns
// ---------------------------------------------------------------------------

// parseGraphPattern parses a graph_pattern:
//
//	<path_pattern_list> [WHERE <expr>]
func (p *Parser) parseGraphPattern() (*ast.GraphPattern, error) {
	gp := &ast.GraphPattern{Loc: p.cur.Loc}
	// graph_path_pattern_list: path (COMMA hint? path)*.
	for {
		path, err := p.parseGraphPathPattern()
		if err != nil {
			return nil, err
		}
		gp.Paths = append(gp.Paths, path)
		gp.Loc.End = path.Loc.End
		if p.cur.Type != int(',') {
			break
		}
		p.advance() // ,
		if p.cur.Type == int('@') {
			if herr := p.skipHint(); herr != nil {
				return nil, herr
			}
		}
	}
	if p.cur.Type == kwWHERE {
		p.advance() // WHERE
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		gp.Where = where
		gp.Loc.End = nodeLoc(where).End
	}
	gp.Loc.Start = gp.Paths[0].Loc.Start
	return gp, nil
}

// parseGraphPathPattern parses a graph_path_pattern:
//
//	[<var> =] [(ANY|ALL) [SHORTEST]] [<mode> [PATH|PATHS]] <path_pattern_expr>
func (p *Parser) parseGraphPathPattern() (*ast.GraphPathPattern, error) {
	pp := &ast.GraphPathPattern{Loc: p.cur.Loc}
	start := p.cur.Loc.Start

	// opt_path_variable_assignment: graph_identifier '='. Only when the '=' is
	// the next token (an element-pattern variable is NOT followed by '=').
	if isIdentifierStart(p.cur.Type) && p.peekNext().Type == int('=') {
		tok := p.advance()
		name, err := p.identifierText(tok)
		if err != nil {
			return nil, err
		}
		pp.PathVar = name
		p.advance() // =
	}

	// opt_graph_search_prefix: (ANY|ALL) SHORTEST?. ANY/ALL are also usable as
	// bare graph_identifiers, so they are consumed as a search prefix only when
	// graphPrefixKeywordContinues confirms a prefix follows (SHORTEST / a further
	// prefix / a path-factor start) — NOT when the keyword is the first node's own
	// identifier (e.g. the parenthesized path `( ANY )` is a single node named
	// ANY). This is the same guard the disambiguator (atGraphPathPrefix) uses, so
	// the two never disagree about whether a leading keyword is a prefix.
	switch p.cur.Type {
	case kwANY:
		if p.graphPrefixKeywordContinues(p.peekNext().Type) {
			p.advance()
			pp.Search = "ANY"
			if p.cur.Type == kwSHORTEST {
				p.advance()
				pp.Search = "ANY SHORTEST"
			}
		}
	case kwALL:
		if p.graphPrefixKeywordContinues(p.peekNext().Type) {
			p.advance()
			pp.Search = "ALL"
			if p.cur.Type == kwSHORTEST {
				p.advance()
				pp.Search = "ALL SHORTEST"
			}
		}
	}

	// opt_graph_path_mode_prefix: opt_graph_path_mode path_or_paths?. The mode
	// keywords are likewise usable as bare identifiers, so they are consumed as a
	// mode prefix only when graphPrefixKeywordContinues confirms a prefix follows
	// (PATH/PATHS / a further prefix / a path-factor start).
	switch p.cur.Type {
	case kwWALK, kwTRAIL, kwSIMPLE, kwACYCLIC:
		if p.graphPrefixKeywordContinues(p.peekNext().Type) {
			mode := TokenName(p.cur.Type)
			p.advance()
			switch p.cur.Type {
			case kwPATH:
				p.advance()
				mode += " PATH"
			case kwPATHS:
				p.advance()
				mode += " PATHS"
			}
			pp.Mode = mode
		}
	}

	// graph_path_pattern_expr: graph_path_factor (hint? graph_path_factor)*.
	//
	// A hint between factors is part of a `(hint? graph_path_factor)` group, so a
	// hint MUST be followed by another factor — a trailing hint with no factor
	// after it is a syntax error (oracle F5: the live emulator rejects
	// `(a) @{h} RETURN *` with `Expected "(" or "-" or "<" or -> but got …`). We
	// therefore only consume a hint once we know a factor follows: if the next
	// token is '@' we require either a factor right after the hint, or we leave the
	// '@' for the caller (which would have its own hint position). Because the only
	// position a trailing '@' could be consumed is here, we detect it and demand a
	// factor.
	for {
		factor, err := p.parseGraphPathFactor()
		if err != nil {
			return nil, err
		}
		pp.Factors = append(pp.Factors, factor)
		pp.Loc.End = nodeLoc(factor).End
		if p.cur.Type != int('@') {
			if !p.atGraphPathFactor() {
				break
			}
			continue
		}
		// A '@' here opens a `hint? graph_path_factor` continuation. Consume the
		// hint, then REQUIRE a following factor (F5).
		if herr := p.skipHint(); herr != nil {
			return nil, herr
		}
		if !p.atGraphPathFactor() {
			return nil, p.syntaxErrorAtCur()
		}
	}
	pp.Loc.Start = start
	return pp, nil
}

// atGraphPathFactor reports whether the cursor is at the start of a
// graph_path_factor: a node pattern `(`, an edge pattern (`-`, `<-`, `->`), or a
// parenthesized path pattern `(`. (Both node patterns and parenthesized paths
// start with '(', disambiguated inside parseGraphPathPrimary.)
func (p *Parser) atGraphPathFactor() bool {
	switch p.cur.Type {
	case int('('), int('-'), int('<'), tokArrow:
		return true
	}
	return false
}

// parseGraphPathFactor parses a graph_path_factor:
//
//	graph_path_primary [ { [n,] m } ]   (a quantified path primary)
func (p *Parser) parseGraphPathFactor() (ast.Node, error) {
	primary, err := p.parseGraphPathPrimary()
	if err != nil {
		return nil, err
	}
	// graph_quantified_path_primary: primary '{' int? ',' int '}' | primary '{' int '}'.
	if p.cur.Type == int('{') {
		return p.parseGraphQuantifier(primary)
	}
	return primary, nil
}

// parseGraphQuantifier consumes a `{ [n,] m }` quantifier following a path
// primary. The quantifier bounds are parse-only (the consumer needs the inner
// pattern's tables, not the repetition count); the primary node is returned with
// its Loc extended over the quantifier. The cursor is at '{'.
func (p *Parser) parseGraphQuantifier(primary ast.Node) (ast.Node, error) {
	p.advance() // {
	// int_literal_or_parameter? ',' int_literal_or_parameter  |  int_literal_or_parameter
	//
	// The bounds are int_literal_or_parameter (integer / @param / ? / @@sysvar) —
	// NOT a full expression. A bare identifier or an arithmetic expression here is
	// a syntax error (oracle: the Spanner emulator rejects `{a,b}` / `{1+1,3}`).
	if p.cur.Type != int(',') {
		if _, err := p.parseGraphIntLiteralOrParameter(); err != nil {
			return nil, err
		}
	}
	if p.cur.Type == int(',') {
		p.advance() // ,
		if _, err := p.parseGraphIntLiteralOrParameter(); err != nil {
			return nil, err
		}
	}
	closeTok, err := p.expect(int('}'))
	if err != nil {
		return nil, err
	}
	// Extend the primary's Loc to include the quantifier so the factor spans it.
	switch n := primary.(type) {
	case *ast.GraphNodePattern:
		n.Loc.End = closeTok.Loc.End
	case *ast.GraphEdgePattern:
		n.Loc.End = closeTok.Loc.End
	case *ast.GraphPathPattern:
		n.Loc.End = closeTok.Loc.End
	}
	return primary, nil
}

// parseGraphPathPrimary parses a graph_path_primary:
//
//	graph_element_pattern   |   graph_parenthesized_path_pattern
//
// A '(' begins either a node pattern `( <filler> )` or a parenthesized path
// `( [hint] <graph_path_pattern> [WHERE] )`. Both can carry a leading hint
// (the node hint lives INSIDE the filler — graph_element_pattern_filler:
// `hint? …`; the parenthesized-path hint sits between '(' and the path — and the
// hint is discarded either way), so we consume '(' and any leading hint FIRST,
// then disambiguate on the post-hint interior token. The parser has only
// one-token lookahead, so peeking PAST a multi-token `@{…}` hint is impossible;
// consuming it up front is what lets a hinted node `(@{h} v:L)` be recognized
// (the disambiguation bug F2). The interior is a parenthesized PATH when it
// begins with a nested path factor (another '(' or an edge `-`/`<`/`->`) or a
// path-pattern prefix (path-variable assignment `id =`, search prefix ANY/ALL,
// or path-mode prefix WALK/TRAIL/SIMPLE/ACYCLIC — bug F1); otherwise it is a
// node-pattern filler (identifier / ':' / IS / '{' / WHERE / empty ')').
func (p *Parser) parseGraphPathPrimary() (ast.Node, error) {
	switch p.cur.Type {
	case int('-'), int('<'), tokArrow:
		return p.parseGraphEdgePattern()
	case int('('):
		open := p.advance() // (
		if p.cur.Type == int('@') {
			if herr := p.skipHint(); herr != nil {
				return nil, herr
			}
		}
		if p.interiorStartsParenthesizedPath() {
			return p.parseGraphParenthesizedPathBody(open)
		}
		return p.parseGraphNodePatternBody(open)
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// interiorStartsParenthesizedPath decides, with cur positioned at the FIRST
// interior token of a '(' group (after the '(' and any leading hint have already
// been consumed), whether the group is a graph_parenthesized_path_pattern (vs a
// graph_node_pattern). It is a parenthesized path when the interior begins with:
//   - a path factor: another '(' or an edge start ('-', '<', '->') — a node
//     filler never starts with one of these; or
//   - a path-pattern prefix: a path-variable assignment (`id =`), a search
//     prefix (ANY/ALL …), or a path-mode prefix (WALK/TRAIL/SIMPLE/ACYCLIC …) —
//     see atGraphPathPrefix.
//
// Everything else (a bare identifier not followed by '=', ':', IS, '{', WHERE,
// or an immediate ')') is a node-pattern filler.
func (p *Parser) interiorStartsParenthesizedPath() bool {
	switch p.cur.Type {
	case int('('), int('-'), int('<'), tokArrow:
		return true
	}
	return p.atGraphPathPrefix()
}

// atGraphPathPrefix reports whether cur is at the start of a graph_path_pattern
// prefix — opt_path_variable_assignment? opt_graph_search_prefix?
// opt_graph_path_mode_prefix? — using the SAME detection parseGraphPathPattern
// applies when it consumes those prefixes (kept in lockstep so the disambiguator
// and the consumer never disagree):
//   - path-variable assignment: an identifier IMMEDIATELY followed by '='.
//   - search prefix: ANY or ALL. These are also usable as bare graph_identifiers,
//     so they are a prefix only when followed by SHORTEST or by something that
//     starts a path factor / further prefix — NOT when followed by a filler
//     continuation (':', IS, '{', WHERE) or a closing ')'.
//   - path-mode prefix: WALK/TRAIL/SIMPLE/ACYCLIC, with the same
//     identifier-vs-prefix disambiguation (a prefix only when followed by
//     PATH/PATHS or a path-factor / further-prefix start).
func (p *Parser) atGraphPathPrefix() bool {
	switch p.cur.Type {
	case kwANY, kwALL, kwWALK, kwTRAIL, kwSIMPLE, kwACYCLIC:
		return p.graphPrefixKeywordContinues(p.peekNext().Type)
	}
	// opt_path_variable_assignment: graph_identifier '='.
	return isIdentifierStart(p.cur.Type) && p.peekNext().Type == int('=')
}

// graphPrefixKeywordContinues reports whether a token following a search/mode
// prefix keyword (ANY/ALL/WALK/TRAIL/SIMPLE/ACYCLIC) confirms that the keyword
// introduces a path-pattern prefix rather than standing in as a node-filler
// identifier. A prefix continues into another prefix word (SHORTEST after
// ANY/ALL; PATH/PATHS after a mode; or a further mode/search keyword) or
// directly into a path factor ('(' or an edge). A following filler-continuation
// (':', IS, '{', WHERE) or a closing ')' means the keyword was the node's own
// identifier, so it does NOT continue.
func (p *Parser) graphPrefixKeywordContinues(next int) bool {
	switch next {
	case kwSHORTEST, kwPATH, kwPATHS,
		kwANY, kwALL, kwWALK, kwTRAIL, kwSIMPLE, kwACYCLIC,
		int('('), int('-'), int('<'), tokArrow:
		return true
	}
	return false
}

// parseGraphParenthesizedPathBody parses the body of a
// graph_parenthesized_path_pattern after its '(' (and any leading hint) have
// already been consumed by parseGraphPathPrimary:
//
//	<graph_path_pattern> [WHERE <expr>] )
//
// open is the already-consumed '(' token (for the Loc start).
func (p *Parser) parseGraphParenthesizedPathBody(open Token) (ast.Node, error) {
	path, err := p.parseGraphPathPattern()
	if err != nil {
		return nil, err
	}
	if p.cur.Type == kwWHERE {
		p.advance() // WHERE
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		path.Where = where
		path.Loc.End = nodeLoc(where).End
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	path.Loc.Start = open.Loc.Start
	path.Loc.End = closeTok.Loc.End
	return path, nil
}

// parseGraphNodePatternBody parses the body of a graph_node_pattern after its
// '(' (and any leading hint) have already been consumed by
// parseGraphPathPrimary:
//
//	<filler> )
//
// open is the already-consumed '(' token (for the Loc start). The filler's own
// optional leading hint has already been skipped at the '(' level, so any '@'
// the filler parser sees here would be a SECOND hint (the grammar admits at most
// one — a stray one is a syntax error from parseGraphPatternFiller's tail).
func (p *Parser) parseGraphNodePatternBody(open Token) (ast.Node, error) {
	filler, err := p.parseGraphPatternFiller()
	if err != nil {
		return nil, err
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	return &ast.GraphNodePattern{
		Filler: filler,
		Loc:    ast.Loc{Start: open.Loc.Start, End: closeTok.Loc.End},
	}, nil
}

// parseGraphEdgePattern parses a graph_edge_pattern — one of the six shapes the
// legacy grammar admits (graph_edge_pattern):
//
//	<-[ filler ]-      (EdgeLeftFull)        LT? MINUS LS filler RS MINUS, < present
//	 -[ filler ]-      (EdgeUndirectedFull)  LT? MINUS LS filler RS MINUS, < absent
//	 -[ filler ]->     (EdgeRightFull)       MINUS LS filler RS SUB_GT
//	-                  (EdgeAny)             MINUS
//	<-                 (EdgeLeft)            LT MINUS
//	->                 (EdgeRight)           SUB_GT
//
// The leading `<` of the first alternative is OPTIONAL, so a full edge with a
// trailing `-` (MINUS) is left-directed when `<` is present and UNDIRECTED when
// it is absent — both are valid. A full edge with a trailing `->` (tokArrow) is
// right-directed and must NOT carry a leading `<` (the oracle rejects `<-[..]->`
// with "Expected '-' but got '->'"). The lexer emits `->` as a single token
// (tokArrow); `<-`, `-[`, `]-` are sequences of single-char tokens.
func (p *Parser) parseGraphEdgePattern() (ast.Node, error) {
	start := p.cur.Loc
	switch p.cur.Type {
	case tokArrow:
		// `->`  (EdgeRight, abbreviated).
		arrow := p.advance()
		return &ast.GraphEdgePattern{Direction: ast.EdgeRight, Loc: ast.Loc{Start: start.Start, End: arrow.Loc.End}}, nil
	case int('<'):
		// `<-` then optionally `[ filler ]-` (left-directed).
		p.advance() // <
		if _, err := p.expect(int('-')); err != nil {
			return nil, err
		}
		if p.cur.Type == int('[') {
			// `<-[ filler ]-`  (EdgeLeftFull): a leading `<` forces the trailing close
			// to be a bare MINUS (`->` is rejected — see the doc note above).
			filler, end, err := p.parseBracketedEdgeFiller(false /*allowArrowClose*/)
			if err != nil {
				return nil, err
			}
			return &ast.GraphEdgePattern{Direction: ast.EdgeLeftFull, Filler: filler, Loc: ast.Loc{Start: start.Start, End: end}}, nil
		}
		// `<-`  (EdgeLeft, abbreviated).
		return &ast.GraphEdgePattern{Direction: ast.EdgeLeft, Loc: ast.Loc{Start: start.Start, End: p.prev.Loc.End}}, nil
	case int('-'):
		p.advance() // -
		if p.cur.Type == int('[') {
			// `-[ filler ]-`  (EdgeUndirectedFull) or `-[ filler ]->` (EdgeRightFull):
			// no leading `<`, so the trailing close may be either MINUS or `->`.
			filler, closeArrow, end, err := p.parseBracketedEdge(true /*allowArrowClose*/)
			if err != nil {
				return nil, err
			}
			dir := ast.EdgeUndirectedFull
			if closeArrow {
				dir = ast.EdgeRightFull
			}
			return &ast.GraphEdgePattern{Direction: dir, Filler: filler, Loc: ast.Loc{Start: start.Start, End: end}}, nil
		}
		// `-`  (EdgeAny, abbreviated undirected).
		return &ast.GraphEdgePattern{Direction: ast.EdgeAny, Loc: ast.Loc{Start: start.Start, End: p.prev.Loc.End}}, nil
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseBracketedEdgeFiller parses the `[ <filler> ]` and trailing close of a full
// edge whose direction is already fixed by a leading `<` (the EdgeLeftFull case),
// so the close MUST be a bare MINUS. allowArrowClose is false here; it returns
// only the filler and end offset.
func (p *Parser) parseBracketedEdgeFiller(allowArrowClose bool) (*ast.GraphPatternFiller, int, error) {
	filler, _, end, err := p.parseBracketedEdge(allowArrowClose)
	return filler, end, err
}

// parseBracketedEdge parses the `[ <filler> ]` of a full edge pattern and its
// trailing close. The current token is '['. The close is always `]` followed by
// either a bare MINUS (`]-`, an undirected/left edge) or, when allowArrowClose is
// true, `->` (`]->`, a right-directed edge). It returns the filler, whether the
// close was the `->` arrow (closeArrow), and the end offset.
func (p *Parser) parseBracketedEdge(allowArrowClose bool) (filler *ast.GraphPatternFiller, closeArrow bool, end int, err error) {
	p.advance() // [
	filler, err = p.parseGraphPatternFiller()
	if err != nil {
		return nil, false, 0, err
	}
	if _, err = p.expect(int(']')); err != nil {
		return nil, false, 0, err
	}
	if allowArrowClose && p.cur.Type == tokArrow {
		// `]->` : a right-directed full edge closes with `->` (tokArrow).
		arrow := p.advance()
		return filler, true, arrow.Loc.End, nil
	}
	// `]-` : an undirected or left-directed full edge closes with a bare MINUS.
	minus, merr := p.expect(int('-'))
	if merr != nil {
		return nil, false, 0, merr
	}
	return filler, false, minus.Loc.End, nil
}

// parseGraphPatternFiller parses a graph_element_pattern_filler — the interior
// of a node `( … )` or full-edge `[ … ]` pattern:
//
//	[hint] [<var>] [(IS|:) <label_expr>] [{ prop : expr, … }] [WHERE <expr>]
//
// Every part is optional (an empty filler `()` is valid, per the grammar's empty
// production). The grammar's three alternatives differ only in how WHERE relates
// to the property spec; this parser accepts the superset (var?, label?, props?,
// where?) which covers all three.
func (p *Parser) parseGraphPatternFiller() (*ast.GraphPatternFiller, error) {
	filler := &ast.GraphPatternFiller{Loc: p.cur.Loc}
	start := p.cur.Loc.Start
	end := p.prev.Loc.End // empty filler: zero-width at the prior token's end

	if p.cur.Type == int('@') {
		if herr := p.skipHint(); herr != nil {
			return nil, herr
		}
	}

	// opt_graph_element_identifier: a bare graph_identifier (NOT followed by '=' —
	// that is a path-variable assignment handled at the path level, which never
	// reaches here). An identifier directly before ':'/'IS'/'{'/WHERE/')'/']' is
	// the element variable.
	if isIdentifierStart(p.cur.Type) {
		tok := p.advance()
		name, err := p.identifierText(tok)
		if err != nil {
			return nil, err
		}
		filler.Var = name
		end = tok.Loc.End
	}

	// opt_is_label_expression: IS <label_expr> | ':' <label_expr>.
	if p.cur.Type == kwIS || p.cur.Type == int(':') {
		filler.LabelColon = p.cur.Type == int(':')
		p.advance() // IS | ':'
		label, err := p.parseLabelExpression()
		if err != nil {
			return nil, err
		}
		filler.Label = label
		end = label.Loc.End
	}

	// graph_property_specification?: `{ id : expr , … }`.
	if p.cur.Type == int('{') {
		props, err := p.parseGraphPropertySpec()
		if err != nil {
			return nil, err
		}
		filler.Properties = props
		end = props.Loc.End
	}

	// where_clause?: WHERE <expr>.
	if p.cur.Type == kwWHERE {
		p.advance() // WHERE
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		filler.Where = where
		end = nodeLoc(where).End
	}

	filler.Loc = ast.Loc{Start: start, End: end}
	return filler, nil
}

// parseGraphPropertySpec parses a graph_property_specification:
//
//	{ <id> : <expr> (, <id> : <expr>)* }
func (p *Parser) parseGraphPropertySpec() (*ast.GraphPropertySpec, error) {
	open := p.advance() // {
	spec := &ast.GraphPropertySpec{Loc: open.Loc}
	for {
		nv, err := p.parseGraphPropertyNameValue()
		if err != nil {
			return nil, err
		}
		spec.Properties = append(spec.Properties, nv)
		if p.cur.Type != int(',') {
			break
		}
		p.advance() // ,
	}
	closeTok, err := p.expect(int('}'))
	if err != nil {
		return nil, err
	}
	spec.Loc.End = closeTok.Loc.End
	return spec, nil
}

// parseGraphPropertyNameValue parses one graph_property_name_and_value:
// identifier ':' expression.
func (p *Parser) parseGraphPropertyNameValue() (*ast.GraphPropertyNameValue, error) {
	tok, err := p.expectIdentifier()
	if err != nil {
		return nil, err
	}
	name, err := p.identifierText(tok)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(int(':')); err != nil {
		return nil, err
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return &ast.GraphPropertyNameValue{
		Name: name,
		Expr: expr,
		Loc:  ast.Loc{Start: tok.Loc.Start, End: nodeLoc(expr).End},
	}, nil
}

// ---------------------------------------------------------------------------
// Label algebra
// ---------------------------------------------------------------------------

// parseLabelExpression parses a label_expression with the precedence the
// legacy grammar's left-recursive rule encodes: `!` (unary, tightest) binds
// before `&`, which binds before `|`. The grammar is:
//
//	label_expression:
//	    label_primary
//	  | label_expression '&' label_expression
//	  | label_expression '|' label_expression
//	  | '!' label_expression
//
// ANTLR resolves the ambiguous binary alternatives by listed order (& before |)
// and left-associativity. We implement that explicitly: parse an OR-chain whose
// operands are AND-chains whose operands are (possibly negated) primaries.
func (p *Parser) parseLabelExpression() (*ast.GraphLabelExpr, error) {
	return p.parseLabelOr()
}

// parseLabelOr parses an `|`-separated chain of AND-expressions (left-assoc).
func (p *Parser) parseLabelOr() (*ast.GraphLabelExpr, error) {
	left, err := p.parseLabelAnd()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == int('|') {
		p.advance() // |
		right, err := p.parseLabelAnd()
		if err != nil {
			return nil, err
		}
		left = &ast.GraphLabelExpr{
			Kind:  ast.LabelOr,
			Left:  left,
			Right: right,
			Loc:   ast.Loc{Start: left.Loc.Start, End: right.Loc.End},
		}
	}
	return left, nil
}

// parseLabelAnd parses an `&`-separated chain of unary label expressions
// (left-assoc).
func (p *Parser) parseLabelAnd() (*ast.GraphLabelExpr, error) {
	left, err := p.parseLabelUnary()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == int('&') {
		p.advance() // &
		right, err := p.parseLabelUnary()
		if err != nil {
			return nil, err
		}
		left = &ast.GraphLabelExpr{
			Kind:  ast.LabelAnd,
			Left:  left,
			Right: right,
			Loc:   ast.Loc{Start: left.Loc.Start, End: right.Loc.End},
		}
	}
	return left, nil
}

// parseLabelUnary parses a `!`-prefixed label expression (right-recursive) or a
// label primary.
func (p *Parser) parseLabelUnary() (*ast.GraphLabelExpr, error) {
	if p.cur.Type == int('!') {
		bang := p.advance() // !
		operand, err := p.parseLabelUnary()
		if err != nil {
			return nil, err
		}
		return &ast.GraphLabelExpr{
			Kind:    ast.LabelNot,
			Operand: operand,
			Loc:     ast.Loc{Start: bang.Loc.Start, End: operand.Loc.End},
		}, nil
	}
	return p.parseLabelPrimary()
}

// parseLabelPrimary parses a label_primary:
//
//	identifier   |   '%'   |   '(' label_expression ')'
func (p *Parser) parseLabelPrimary() (*ast.GraphLabelExpr, error) {
	switch {
	case p.cur.Type == int('%'):
		mod := p.advance()
		return &ast.GraphLabelExpr{Kind: ast.LabelWildcard, Loc: mod.Loc}, nil
	case p.cur.Type == int('('):
		p.advance() // (
		inner, err := p.parseLabelExpression()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
		// parenthesized_label_expression is transparent: the tree shape already
		// captures the grouping, so we return the inner expression directly.
		return inner, nil
	case isIdentifierStart(p.cur.Type):
		tok := p.advance()
		name, err := p.identifierText(tok)
		if err != nil {
			return nil, err
		}
		return &ast.GraphLabelExpr{Kind: ast.LabelName, Name: name, Loc: tok.Loc}, nil
	default:
		return nil, p.syntaxErrorAtCur()
	}
}
