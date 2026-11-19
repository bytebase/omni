package parser

import (
	"strings"

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
// accept/reject is recorded by the differential (the suite:graph-query rows of testdata/compat.tsv) and
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

	blocks, end, err := p.parseGraphOperationBlock()
	if err != nil {
		return nil, err
	}
	stmt.Blocks = blocks
	stmt.Loc.End = end
	return stmt, nil
}

// parseGraphOperationBlock parses a graph_operation_block:
//
//	graph_composite_query_block ("NEXT" graph_composite_query_block)*
//
// It returns the blocks and the end offset of the last one. Both the GRAPH
// statement and a braced graph subquery hold this list.
func (p *Parser) parseGraphOperationBlock() ([]ast.Node, int, error) {
	var blocks []ast.Node
	end := p.cur.Loc.Start
	for {
		block, err := p.parseGraphCompositeQueryBlock()
		if err != nil {
			return nil, 0, err
		}
		blocks = append(blocks, block)
		end = nodeLoc(block).End
		if p.cur.Type != kwNEXT {
			break
		}
		p.advance() // NEXT
	}
	return blocks, end, nil
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
	return p.parseGraphSetOpChain(first)
}

// parseGraphSetOpChain parses the set-operation tail of a
// graph_composite_query_prefix after its first linear query operation: the
// block is `linear (meta linear)*`. first is already parsed; when no set
// operation follows, the block is just first.
func (p *Parser) parseGraphSetOpChain(first *ast.GraphLinearQuery) (ast.Node, error) {
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
	// graph_set_operation_metadata is
	// opt_corresponding_outer_mode query_set_operation_type all_or_distinct, so a
	// leading FULL/OUTER/LEFT/INNER belongs to it when a set operator follows. The
	// words are also join keywords, so the same lookahead the query set operation
	// uses decides it here.
	if atSetOpKeyword(p.cur.Type) {
		return true
	}
	switch p.cur.Type {
	case kwFULL, kwLEFT:
		return p.setOpFollowsOuterMode()
	case kwOUTER, kwINNER:
		return atSetOpKeyword(p.peekNext().Type)
	}
	return false
}

// atSetOpKeyword reports whether tokenType is a query_set_operation_type.
func atSetOpKeyword(tokenType int) bool {
	switch tokenType {
	case kwUNION, kwEXCEPT, kwINTERSECT:
		return true
	}
	return false
}

// parseGraphSetOpMeta parses a graph_set_operation_metadata:
// query_set_operation_type all_or_distinct. Both the operation type
// (UNION|EXCEPT|INTERSECT) and the quantifier (ALL|DISTINCT) are REQUIRED.
func (p *Parser) parseGraphSetOpMeta() (*ast.GraphSetOpMeta, error) {
	meta := &ast.GraphSetOpMeta{Loc: p.cur.Loc}
	// opt_corresponding_outer_mode: FULL [OUTER] | OUTER | LEFT [OUTER] | INNER.
	switch p.cur.Type {
	case kwFULL, kwLEFT:
		modeTok := p.advance()
		meta.OuterMode = upperTokenName(modeTok.Type)
		if p.cur.Type == kwOUTER {
			p.advance()
			meta.Loc.End = p.prev.Loc.End
		}
	case kwOUTER, kwINNER:
		modeTok := p.advance()
		meta.OuterMode = upperTokenName(modeTok.Type)
	}
	opTok := p.advance() // UNION | EXCEPT | INTERSECT (atGraphSetOpMeta gated)
	meta.Op = TokenName(opTok.Type)
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
	var ops []ast.Node
	for p.atGraphLinearOperator() {
		op, err := p.parseGraphLinearOperator()
		if err != nil {
			return nil, err
		}
		ops = append(ops, op)
	}
	return p.finishGraphLinearQuery(ops)
}

// finishGraphLinearQuery parses the RETURN operator that terminates a linear
// query operation and wraps the already-parsed operators with it. ops may be
// empty (`RETURN *` alone is a complete linear query operation).
func (p *Parser) finishGraphLinearQuery(ops []ast.Node) (*ast.GraphLinearQuery, error) {
	lq := &ast.GraphLinearQuery{Loc: p.cur.Loc, Operators: ops}
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
		kwTABLESAMPLE, kwOFFSET, kwSKIP, kwLIMIT, kwCALL:
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
		// OPTIONAL introduces graph_optional_match_operator or the OPTIONAL form
		// of graph_call_operator; no other linear operator takes it.
		if p.peekNext().Type == kwCALL {
			optional := p.advance() // OPTIONAL
			op, err := p.parseGraphCallOp(true)
			if err != nil {
				return nil, err
			}
			// The node spans OPTIONAL, its first token.
			call := op.(*ast.GraphCallOp)
			call.Loc.Start = optional.Loc.Start
			return call, nil
		}
		return p.parseGraphMatchOp(true)
	case kwCALL:
		return p.parseGraphCallOp(false)
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
	// partition_by_clause_prefix_no_hint? — PARTITION BY expr (, expr)* inside the
	// sample size. The expressions are kept: a scalar subquery in one is a read.
	if p.cur.Type == kwPARTITION {
		p.advance() // PARTITION
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		for {
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			op.PartitionBy = append(op.PartitionBy, expr)
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
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
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

	// opt_graph_search_prefix: ANY / ALL with an optional SHORTEST, CHEAPEST or
	// count. ANY and ALL are also usable as bare graph_identifiers, so they are
	// consumed as a prefix only when graphPrefixKeywordContinues confirms a prefix
	// follows — NOT when the keyword is the first node's own identifier (e.g. the
	// parenthesized path `( ANY )` is a single node named ANY). This is the same
	// guard the disambiguator (atGraphPathPrefix) uses, so the two never disagree
	// about whether a leading keyword is a prefix.
	search, err := p.parseGraphSearchPrefix()
	if err != nil {
		return nil, err
	}
	pp.Search = search

	// opt_graph_path_mode_prefix: opt_graph_path_mode (PATH | PATHS)?. The mode
	// keywords are likewise usable as bare identifiers, so they are consumed as a
	// mode prefix only when graphPrefixKeywordContinues confirms a prefix follows
	// (PATH/PATHS / a further prefix / a path-factor start). PATH and PATHS need no
	// mode in front of them (`PATH -` is a path in PATH mode).
	pp.Mode, err = p.parseGraphPathMode()
	if err != nil {
		return nil, err
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
	factor, err := p.parseGraphPathFactor()
	if err != nil {
		return nil, err
	}
	pp.Factors = append(pp.Factors, factor)
	pp.Loc.End = nodeLoc(factor).End
	for {
		// A '@' here opens a `hint? graph_path_factor` continuation. Consume the
		// hint, then REQUIRE a following factor (F5).
		hintLoc := ast.Loc{}
		hasHint := false
		if p.cur.Type == int('@') {
			hintLoc = p.cur.Loc
			if herr := p.skipHint(); herr != nil {
				return nil, herr
			}
			hasHint = true
			if !p.atGraphPathFactor() {
				return nil, p.syntaxErrorAtCur()
			}
		} else if !p.atGraphPathFactor() {
			break
		}
		next, err := p.parseGraphPathFactor()
		if err != nil {
			return nil, err
		}
		// A hint between two EDGE patterns is ambiguous and the grammar refuses it
		// ("Hint cannot be used in between two …"), so `- @1 -` rejects while a
		// quantified edge (a one-factor path pattern) or a node pattern does not.
		if hasHint {
			_, lastIsEdge := pp.Factors[len(pp.Factors)-1].(*ast.GraphEdgePattern)
			if _, nextIsEdge := next.(*ast.GraphEdgePattern); lastIsEdge && nextIsEdge {
				return nil, &ParseError{
					Position: hintLoc.Start,
					End:      hintLoc.End,
					Message:  "Syntax error: Hint cannot be used in between two edge patterns",
				}
			}
		}
		pp.Factors = append(pp.Factors, next)
		pp.Loc.End = nodeLoc(next).End
	}
	pp.Loc.Start = start
	return pp, nil
}

// atGraphEdgeStart reports whether tokenType begins a graph_edge_pattern (one of
// the six shapes: `-`, `<-`, `->`, and their bracketed forms).
func atGraphEdgeStart(tokenType int) bool {
	switch tokenType {
	case int('-'), int('<'), tokArrow:
		return true
	}
	return false
}

// atGraphPathFactor reports whether the cursor is at the start of a
// graph_path_factor: a node pattern `(`, an edge pattern (`-`, `<-`, `->`), or a
// parenthesized path pattern `(`. (Both node patterns and parenthesized paths
// start with '(', disambiguated inside parseGraphPathPrimary.)
func (p *Parser) atGraphPathFactor() bool {
	return atGraphPathFactorToken(p.cur)
}

// atGraphPathFactorToken reports whether tok begins a graph_path_factor.
func atGraphPathFactorToken(tok Token) bool {
	switch tok.Type {
	case int('('), int('-'), int('<'), tokArrow:
		return true
	}
	return false
}

// parseGraphPathFactor parses a graph_path_factor:
//
//	graph_path_primary [ graph_quantifier ]   (a quantified path primary)
func (p *Parser) parseGraphPathFactor() (ast.Node, error) {
	primary, err := p.parseGraphPathPrimary()
	if err != nil {
		return nil, err
	}
	// graph_quantified_path_primary: graph_path_primary graph_quantifier.
	if p.atGraphQuantifier() {
		return p.parseGraphQuantifier(primary)
	}
	return primary, nil
}

// atGraphQuantifier reports whether the cursor begins a graph_quantifier:
// `{ [n,] m }` (either bound optional), `{ n }`, or the `+` / `*` shorthand.
func (p *Parser) atGraphQuantifier() bool {
	switch p.cur.Type {
	case int('{'), int('+'), int('*'):
		return true
	}
	return false
}

// parseGraphQuantifier consumes a graph_quantifier following a path primary:
//
//	{ [n,] m } | { n }   (both bounds optional in the first form)
//	+ | *                (plus_or_star_quantifier)
//
// The bounds are parse-only (the consumer needs the inner pattern's tables, not
// the repetition count); the primary node is returned with its Loc extended over
// the quantifier. The grammar refuses a quantifier on a node pattern
// ("Quantifier cannot be used on a node pattern"), so that is an error here.
func (p *Parser) parseGraphQuantifier(primary ast.Node) (ast.Node, error) {
	if node, ok := primary.(*ast.GraphNodePattern); ok {
		return nil, &ParseError{
			Position: node.Loc.Start,
			End:      node.Loc.End,
			Message:  "Syntax error: Quantifier cannot be used on a node pattern",
		}
	}
	// The bounds are int_literal_or_parameter (integer / @param / ? / @@sysvar) —
	// NOT a full expression. A bare identifier or an arithmetic expression here is
	// a syntax error (oracle: the Spanner emulator rejects `{a,b}` / `{1+1,3}`).
	var closeTok Token
	if p.cur.Type == int('+') || p.cur.Type == int('*') {
		// plus_or_star_quantifier: a bare `+` or `*`.
		p.advance()
		closeTok = p.prev
	} else {
		p.advance() // {
		// int_literal_or_parameter? ',' int_literal_or_parameter? | int_literal_or_parameter
		// A brace with neither a bound nor the comma is not a quantifier (`{}`).
		if !atIntLiteralOrParameter(p.cur.Type) && p.cur.Type != int(',') {
			return nil, p.syntaxErrorAtCur()
		}
		if atIntLiteralOrParameter(p.cur.Type) {
			if _, err := p.parseGraphIntLiteralOrParameter(); err != nil {
				return nil, err
			}
		}
		if p.cur.Type == int(',') {
			p.advance() // ,
			if atIntLiteralOrParameter(p.cur.Type) {
				if _, err := p.parseGraphIntLiteralOrParameter(); err != nil {
					return nil, err
				}
			}
		}
		var err error
		closeTok, err = p.expect(int('}'))
		if err != nil {
			return nil, err
		}
	}
	// Extend the primary's Loc to include the quantifier so the factor spans it. A
	// quantified EDGE becomes a one-factor graph_path_pattern, as the grammar's
	// action wraps it (`quantifier_container = MakeNode<ASTGraphPathPattern>`), so a
	// quantified edge is a path and not a bare edge factor — which is what keeps
	// the hint rule below from reading `-+ @1 -` as a hint between two edges.
	switch n := primary.(type) {
	case *ast.GraphEdgePattern:
		// The edge keeps its own range; the quantifier belongs to the wrapper.
		return &ast.GraphPathPattern{
			Factors: []ast.Node{n},
			Loc:     ast.Loc{Start: n.Loc.Start, End: closeTok.Loc.End},
		}, nil
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
		hintConsumed := false
		hintLoc := p.cur.Loc
		if p.cur.Type == int('@') {
			if herr := p.skipHint(); herr != nil {
				return nil, herr
			}
			hintConsumed = true
		}
		if p.interiorStartsParenthesizedPath() {
			// A hint may open a node pattern's filler, but a parenthesized path
			// pattern takes none — the grammar refuses it outright ("Hint cannot be
			// used at beginning of path pattern"), so `(@1 -)` rejects.
			if hintConsumed {
				return nil, &ParseError{
					Position: hintLoc.Start,
					End:      hintLoc.End,
					Message:  "Syntax error: Hint cannot be used at beginning of path pattern",
				}
			}
			return p.parseGraphParenthesizedPathBody(open)
		}
		return p.parseGraphNodePatternBody(open, hintConsumed)
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
		return p.graphPrefixKeywordContinues(p.peekNext())
	case kwPATH, kwPATHS:
		// PATH/PATHS need no mode in front of them, but they are still only a
		// prefix when a path factor follows: `(PATH)` is a node named PATH.
		return atGraphPathFactorToken(p.peekNext())
	}
	if p.curIsGraphWord("SHORTEST") || p.curIsGraphWord("CHEAPEST") {
		return atIntLiteralOrParameter(p.peekNext().Type)
	}
	// opt_path_variable_assignment: graph_identifier '='.
	return isIdentifierStart(p.cur.Type) && p.peekNext().Type == int('=')
}

// graphPrefixKeywordContinues reports whether a token following a search/mode
// prefix keyword (ANY/ALL/WALK/TRAIL/SIMPLE/ACYCLIC) confirms that the keyword
// introduces a path-pattern prefix rather than standing in as a node-filler
// identifier. A prefix continues into another prefix word (SHORTEST or CHEAPEST
// after ANY/ALL; PATH/PATHS after a mode; or a further mode/search keyword), into
// a search count (int_literal_or_parameter, as in `ANY 2`), or directly into a
// path factor ('(' or an edge). A following filler-continuation (':', IS, '{',
// WHERE) or a closing ')' means the keyword was the node's own identifier, so it
// does NOT continue.
func (p *Parser) graphPrefixKeywordContinues(next Token) bool {
	// CHEAPEST (and the counts after ANY/SHORTEST/CHEAPEST) continue a search
	// prefix. CHEAPEST is a tm keyword the legacy lexer does not carry, so it
	// arrives as an identifier and is matched by its source text.
	if tokenIsUnquotedWord(next, "CHEAPEST") {
		return true
	}
	switch next.Type {
	case kwSHORTEST, kwPATH, kwPATHS,
		kwANY, kwALL, kwWALK, kwTRAIL, kwSIMPLE, kwACYCLIC,
		int('('), int('-'), int('<'), tokArrow,
		tokInteger, int('@'), int('?'), tokAtAt:
		return true
	}
	return false
}

// curIsGraphWord reports whether the current token is the unquoted word `word`,
// whether the lexer reads it as an identifier or as a keyword. The graph rules
// name words the legacy GoogleSQLLexer.g4 does not carry (CHEAPEST, COST,
// LABELED, PER, YIELD) — those arrive as identifiers — beside words it does
// (SHORTEST), so both spellings are matched. A backtick-quoted word is an
// identifier of that text and never the keyword, which the span check enforces.
func (p *Parser) curIsGraphWord(word string) bool {
	if p.cur.Type == tokIdentifier {
		return tokenIsUnquotedWord(p.cur, word)
	}
	return p.cur.Type >= keywordBase && strings.EqualFold(TokenName(p.cur.Type), word)
}

// tokenIsUnquotedWord reports whether tok is the unquoted word `word`. A
// backtick-quoted `cheapest` is the same text but names an identifier, so the
// check compares the token's byte span with its text, as curIsUnquotedWord does.
func tokenIsUnquotedWord(tok Token, word string) bool {
	return tok.Type == tokIdentifier && strings.EqualFold(tok.Str, word) &&
		tok.Loc.End-tok.Loc.Start == len(word)
}

// atIntLiteralOrParameter reports whether tokenType begins an
// int_literal_or_parameter: integer | @param | ? | @@sysvar.
func atIntLiteralOrParameter(tokenType int) bool {
	switch tokenType {
	case tokInteger, int('@'), int('?'), tokAtAt:
		return true
	}
	return false
}

// parseGraphSearchPrefix parses an opt_graph_search_prefix and returns the
// spelling the AST keeps:
//
//	ANY | ALL                                 (search with no bound)
//	ANY SHORTEST | ALL SHORTEST | ANY CHEAPEST | ALL CHEAPEST
//	ANY <n> | SHORTEST <n> | CHEAPEST <n>     (a count, as source text)
//
// ANY and ALL stand alone, so each is a prefix only when
// graphPrefixKeywordContinues confirms one follows; SHORTEST and CHEAPEST always
// take a count. The cursor is left after the prefix ("" when there is none).
func (p *Parser) parseGraphSearchPrefix() (string, error) {
	var lead string
	switch {
	case p.cur.Type == kwANY:
		lead = "ANY"
	case p.cur.Type == kwALL:
		lead = "ALL"
	case p.curIsGraphWord("SHORTEST"):
		lead = "SHORTEST"
	case p.curIsGraphWord("CHEAPEST"):
		lead = "CHEAPEST"
	default:
		return "", nil
	}
	if lead == "SHORTEST" || lead == "CHEAPEST" {
		// Unlike ANY and ALL, these never stand alone: the grammar takes them only
		// with a count (`SHORTEST 1`), so a bare SHORTEST is an identifier.
		if !atIntLiteralOrParameter(p.peekNext().Type) {
			return "", nil
		}
	} else if !p.graphPrefixKeywordContinues(p.peekNext()) {
		return "", nil
	}
	p.advance() // ANY | ALL | SHORTEST | CHEAPEST
	switch {
	case p.cur.Type == kwSHORTEST:
		p.advance()
		return lead + " SHORTEST", nil
	case p.curIsGraphWord("CHEAPEST"):
		p.advance()
		return lead + " CHEAPEST", nil
	case lead != "ALL" && atIntLiteralOrParameter(p.cur.Type):
		// Counts follow ANY, SHORTEST and CHEAPEST only: graph_search_prefix has no
		// `ALL <count>` alternative, so `ALL 2 -` is not a path.
		start := p.cur.Loc.Start
		if _, err := p.parseGraphIntLiteralOrParameter(); err != nil {
			return "", err
		}
		count := p.input[absIndex(p, start):absIndex(p, p.prev.Loc.End)]
		return lead + " " + count, nil
	}
	return lead, nil
}

// parseGraphPathMode parses opt_graph_path_mode_prefix: an optional
// graph_path_mode followed by an optional PATH | PATHS. It returns the spelling
// the AST keeps ("WALK", "TRAIL PATH", "PATH", …; "" when there is none).
//
// The mode words are also usable as bare graph identifiers, so a mode is
// consumed only when graphPrefixKeywordContinues confirms one follows. PATH and
// PATHS need no mode in front of them (`PATH -` is a path in PATH mode), and
// they are consumed only when a path factor follows, so `(PATH)` stays a node
// named PATH.
func (p *Parser) parseGraphPathMode() (string, error) {
	mode := ""
	switch p.cur.Type {
	case kwWALK, kwTRAIL, kwSIMPLE, kwACYCLIC:
		if p.graphPrefixKeywordContinues(p.peekNext()) {
			mode = TokenName(p.cur.Type)
			p.advance()
		}
	}
	switch p.cur.Type {
	case kwPATH:
		if mode != "" || atGraphPathFactorToken(p.peekNext()) {
			mode += " PATH"
			p.advance()
		}
	case kwPATHS:
		if mode != "" || atGraphPathFactorToken(p.peekNext()) {
			mode += " PATHS"
			p.advance()
		}
	}
	return strings.TrimSpace(mode), nil
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
func (p *Parser) parseGraphNodePatternBody(open Token, hintAlreadyConsumed bool) (ast.Node, error) {
	filler, err := p.parseGraphPatternFiller(hintAlreadyConsumed)
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
		// `<-` then optionally `[ filler ]-` (left-directed). The grammar's
		// ValidateNoWhitespace refuses `< -` and `< - [`, so the delimiters must
		// touch (`< -` and `< - [ ] -` are recorded rejections).
		p.advance() // <
		if _, err := p.expectAdjacent(int('-')); err != nil {
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
	// `-[` : the open bracket touches the dash (ValidateNoWhitespace).
	if _, err = p.expectAdjacent(int('[')); err != nil {
		return nil, false, 0, err
	}
	filler, err = p.parseGraphPatternFiller(false)
	if err != nil {
		return nil, false, 0, err
	}
	if _, err = p.expect(int(']')); err != nil {
		return nil, false, 0, err
	}
	if allowArrowClose && p.cur.Type == tokArrow {
		// `]->` : a right-directed full edge closes with `->` (tokArrow), adjacent
		// to the bracket (ValidateNoWhitespace).
		if end, err = p.expectAdjacent(tokArrow); err != nil {
			return nil, false, 0, err
		}
		return filler, true, end, nil
	}
	// `]-` : an undirected or left-directed full edge closes with a bare MINUS,
	// adjacent to the bracket (ValidateNoWhitespace).
	if end, err = p.expectAdjacent(int('-')); err != nil {
		return nil, false, 0, err
	}
	return filler, false, end, nil
}

// expectAdjacent consumes a token of the given type and reports the offset past
// it, requiring that it touches the previous token. graph_edge_pattern's
// delimiter pairs are written without whitespace (`-[`, `]-`, `]->`, `<-`), and
// the grammar's ValidateNoWhitespace actions refuse a space between them, so
// `- [ ] -` is a syntax error rather than a spaced edge pattern.
func (p *Parser) expectAdjacent(tokenType int) (int, error) {
	if p.cur.Type != tokenType {
		return 0, p.syntaxErrorAtCur()
	}
	if p.prev.Loc.End != p.cur.Loc.Start {
		return 0, &ParseError{
			Position: p.cur.Loc.Start,
			End:      p.cur.Loc.End,
			Message:  "Syntax error: Unexpected whitespace inside an edge pattern",
		}
	}
	tok := p.advance()
	return tok.Loc.End, nil
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
func (p *Parser) parseGraphPatternFiller(hintAlreadyConsumed bool) (*ast.GraphPatternFiller, error) {
	filler := &ast.GraphPatternFiller{Loc: p.cur.Loc}
	start := p.cur.Loc.Start
	end := p.prev.Loc.End // empty filler: zero-width at the prior token's end

	// graph_element_pattern_filler admits at most one leading hint. When the caller
	// (a node pattern via parseGraphPathPrimary) already consumed the single hint at
	// the '(' level, do NOT consume another here — a second '@' is a syntax error
	// (oracle: `(@{h1} @{h2} v)` -> live emulator `Expected ")" but got "@"`).
	if !hintAlreadyConsumed && p.cur.Type == int('@') {
		if herr := p.skipHint(); herr != nil {
			return nil, herr
		}
	}

	// opt_graph_element_identifier: a bare graph_identifier (NOT followed by '=' —
	// that is a path-variable assignment handled at the path level, which never
	// reaches here). An identifier directly before ':'/'IS'/'{'/WHERE/')'/']' is
	// the element variable.
	//
	// opt_graph_cost — the grammar's third filler tail, `COST <expression>` — is
	// NOT read here: the pinned emulator syntax-rejects it, compat.tsv records
	// that (`opt_graph_cost/COST_expression`, `SELECT EXISTS {(COST NULL) }`), and
	// the word parses as an element named COST, which is what the grammar's own
	// shift preference and that recording agree on.
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
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
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

// ---------------------------------------------------------------------------
// Graph subqueries and CALL
// ---------------------------------------------------------------------------

// parseGraphBracedSubquery parses a braced_graph_subquery — a whole graph query
// inside braces, used as an expression operand; the cursor is at '{':
//
//	{ graph_operation_block }
//	{ GRAPH <path_expression> graph_operation_block }
//
// The operation block always ends in RETURN, so the subquery yields rows. A
// subquery nests, so it passes through enterNesting.
func (p *Parser) parseGraphBracedSubquery() (*ast.GraphSubqueryExpr, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	open := p.advance() // {
	sub := &ast.GraphSubqueryExpr{Loc: open.Loc}
	graph, err := p.parseGraphSubqueryGraph()
	if err != nil {
		return nil, err
	}
	sub.Graph = graph
	blocks, _, err := p.parseGraphOperationBlock()
	if err != nil {
		return nil, err
	}
	closeTok, err := p.expect(int('}'))
	if err != nil {
		return nil, err
	}
	sub.Blocks = blocks
	sub.Loc.End = closeTok.Loc.End
	return sub, nil
}

// parseGraphSubqueryGraph parses the optional `GRAPH <path_expression>` prefix
// every braced graph subquery admits, after the '{' (or the EXISTS keyword's
// brace) has been consumed. It returns nil when the graph is left implicit.
func (p *Parser) parseGraphSubqueryGraph() (*ast.PathExpr, error) {
	if p.cur.Type != kwGRAPH {
		return nil, nil
	}
	p.advance() // GRAPH
	// braced_graph_subquery takes a plain path_expression — dotted identifier
	// components only — so a dashed name (`GRAPH a-b`) is not folded here the way
	// BigQuery's table paths are.
	return p.parsePathExpr()
}

// parseGraphExistsBraced parses the braced alternatives of exists_graph_subquery
// (exists_graph_pattern_subquery, exists_linear_ops_subquery, and
// `EXISTS [hint] braced_graph_subquery`). The EXISTS keyword and any hint have
// been consumed; start is EXISTS's start offset and the cursor is at '{'.
//
// What follows the brace decides the form:
//
//	RETURN …                → a complete graph query (GraphSubqueryExpr)
//	linear operators, RETURN→ a complete graph query (GraphSubqueryExpr)
//	linear operators        → a linear-operator run (GraphExistsExpr)
//	a graph pattern         → a bare pattern (GraphExistsExpr)
func (p *Parser) parseGraphExistsBraced(start int) (ast.Node, error) {
	if err := p.enterNesting(); err != nil {
		return nil, err
	}
	defer p.leaveNesting()
	p.advance() // {
	graph, err := p.parseGraphSubqueryGraph()
	if err != nil {
		return nil, err
	}
	if p.atGraphLinearOperator() || p.cur.Type == kwRETURN {
		var ops []ast.Node
		for p.atGraphLinearOperator() {
			op, err := p.parseGraphLinearOperator()
			if err != nil {
				return nil, err
			}
			ops = append(ops, op)
		}
		if p.cur.Type == kwRETURN {
			// A complete linear query operation: this is the braced-subquery
			// alternative, so the rest is an operation block.
			first, err := p.finishGraphLinearQuery(ops)
			if err != nil {
				return nil, err
			}
			block, err := p.parseGraphSetOpChain(first)
			if err != nil {
				return nil, err
			}
			blocks := []ast.Node{block}
			for p.cur.Type == kwNEXT {
				p.advance() // NEXT
				b, err := p.parseGraphCompositeQueryBlock()
				if err != nil {
					return nil, err
				}
				blocks = append(blocks, b)
			}
			closeTok, err := p.expect(int('}'))
			if err != nil {
				return nil, err
			}
			return &ast.GraphSubqueryExpr{
				Modifier: ast.GraphSubqueryExists,
				Graph:    graph,
				Blocks:   blocks,
				Loc:      ast.Loc{Start: start, End: closeTok.Loc.End},
			}, nil
		}
		closeTok, err := p.expect(int('}'))
		if err != nil {
			return nil, err
		}
		return &ast.GraphExistsExpr{
			Kind:  ast.GraphExistsLinearOps,
			Graph: graph,
			Ops:   ops,
			Loc:   ast.Loc{Start: start, End: closeTok.Loc.End},
		}, nil
	}
	pattern, err := p.parseGraphPattern()
	if err != nil {
		return nil, err
	}
	closeTok, err := p.expect(int('}'))
	if err != nil {
		return nil, err
	}
	return &ast.GraphExistsExpr{
		Kind:    ast.GraphExistsPattern,
		Graph:   graph,
		Pattern: pattern,
		Loc:     ast.Loc{Start: start, End: closeTok.Loc.End},
	}, nil
}

// parseGraphCallOp parses a graph_call_operator:
//
//	[OPTIONAL] CALL opt_per_clause tvf opt_yield_clause
//	[OPTIONAL] CALL opt_per_clause braced_graph_subquery
//	           CALL parenthesized_identifier_list braced_graph_subquery
//
// optional says the OPTIONAL keyword has already been consumed; the cursor is at
// CALL. PER and YIELD are tm keywords the legacy lexer does not carry, so they
// arrive as identifiers and are matched by source text.
func (p *Parser) parseGraphCallOp(optional bool) (ast.Node, error) {
	callTok := p.advance() // CALL
	op := &ast.GraphCallOp{Optional: optional, Loc: callTok.Loc}
	// The identifier-list spelling has no PER keyword and is not a partitioning
	// call: `CALL ( id, … ) { … }`.
	if p.cur.Type == int('(') {
		ids, end, err := p.parseGraphIdentifierList()
		if err != nil {
			return nil, err
		}
		op.Per = ids
		op.Loc.End = end
		if p.cur.Type != int('{') {
			return nil, p.syntaxErrorAtCur()
		}
		sub, err := p.parseGraphBracedSubquery()
		if err != nil {
			return nil, err
		}
		op.Subquery = sub
		op.Loc.End = sub.Loc.End
		return op, nil
	}
	// opt_per_clause: PER ( id, … ). A `PER (` is the clause; a table-valued
	// function merely named `per` is not (the grammar's LALR preference reads it
	// the same way, and `CALL PER() a()` needs the empty list).
	if p.curIsGraphWord("PER") && p.peekNext().Type == int('(') {
		p.advance() // PER
		ids, end, err := p.parseGraphIdentifierList()
		if err != nil {
			return nil, err
		}
		op.HasPer = true
		op.Per = ids
		op.Loc.End = end
	}
	if p.cur.Type == int('{') {
		sub, err := p.parseGraphBracedSubquery()
		if err != nil {
			return nil, err
		}
		op.Subquery = sub
		op.Loc.End = sub.Loc.End
		return op, nil
	}
	// A named table-valued function, then its optional YIELD list. graph_call_
	// operator_core takes the grammar's `tvf`, whose arguments admit TABLE /
	// MODEL / CONNECTION / INPUT TABLE, so this is the shared TVF parser rather
	// than an ordinary scalar call.
	fc, err := p.parseTVF()
	if err != nil {
		return nil, err
	}
	op.Call = fc
	op.Loc.End = fc.Loc.End
	if p.curIsGraphWord("YIELD") {
		p.advance() // YIELD
		items, end, err := p.parseGraphYieldList()
		if err != nil {
			return nil, err
		}
		op.Yield = items
		op.Loc.End = end
	}
	return op, nil
}

// parseGraphIdentifierList parses a parenthesized_identifier_list — `( id, … )`,
// possibly empty — and returns the identifiers and the end offset of the closing
// ')'. The cursor is at '('.
func (p *Parser) parseGraphIdentifierList() ([]string, int, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, 0, err
	}
	var ids []string
	for p.cur.Type != int(')') {
		tok, err := p.expectIdentifier()
		if err != nil {
			return nil, 0, err
		}
		id, err := p.identifierText(tok)
		if err != nil {
			return nil, 0, err
		}
		ids = append(ids, id)
		if p.cur.Type != int(',') {
			break
		}
		p.advance() // ,
		// parenthesized_identifier_list is an identifier_list: no trailing comma.
		if p.cur.Type == int(')') {
			return nil, 0, p.syntaxErrorAtCur()
		}
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, 0, err
	}
	return ids, closeTok.Loc.End, nil
}

// parseGraphYieldList parses a yield_item_list — `identifier [AS alias]` items
// separated by commas — and returns them with the end offset of the last item.
// The YIELD keyword has been consumed.
func (p *Parser) parseGraphYieldList() ([]*ast.GraphYieldItem, int, error) {
	var items []*ast.GraphYieldItem
	end := p.cur.Loc.End
	for {
		tok, err := p.expectIdentifier()
		if err != nil {
			return nil, 0, err
		}
		name, err := p.identifierText(tok)
		if err != nil {
			return nil, 0, err
		}
		item := &ast.GraphYieldItem{Name: name, Loc: tok.Loc}
		end = tok.Loc.End
		if p.cur.Type == kwAS {
			p.advance() // AS
			aliasTok, err := p.expectIdentifier()
			if err != nil {
				return nil, 0, err
			}
			alias, err := p.identifierText(aliasTok)
			if err != nil {
				return nil, 0, err
			}
			item.Alias = alias
			item.Loc.End = aliasTok.Loc.End
			end = aliasTok.Loc.End
		}
		items = append(items, item)
		if p.cur.Type != int(',') {
			break
		}
		p.advance() // ,
	}
	return items, end, nil
}

// ---------------------------------------------------------------------------
// Graph predicates (graph_expression)
// ---------------------------------------------------------------------------

// parseGraphIsPredicate parses the graph_expression predicates that begin with
// IS [NOT] (the cursor is just past IS, with `not` recording a NOT):
//
//	IS [NOT] SOURCE [OF] <expr>          (edge_source_endpoint_operator)
//	IS [NOT] DESTINATION [OF] <expr>     (edge_dest_endpoint_operator)
//	IS [NOT] LABELED <label_expression>  (is_labeled_operator)
//
// handled is false when the token is none of these, so the caller reads
// `IS NULL` / `IS DISTINCT FROM` as before. Left is the expression the IS
// applies to. The endpoint RHS is at the EHP level (the grammar's
// expression_higher_prec_than_and), and the whole predicate is non-associative
// with the other comparisons, as EDGE_ENDPOINT_PRECEDENCE sits among them.
func (p *Parser) parseGraphIsPredicate(left ast.Node, not bool) (ast.Node, bool, error) {
	kind := ast.GraphIsLabeled
	switch {
	case p.cur.Type == kwSOURCE:
		kind = ast.GraphIsSource
	case p.cur.Type == kwDESTINATION:
		kind = ast.GraphIsDestination
	case p.curIsGraphWord("LABELED"):
		// LABELED is a tm keyword the legacy lexer does not carry, so it arrives
		// as an identifier and is matched by its source text.
		p.advance() // LABELED
		label, err := p.parseLabelExpression()
		if err != nil {
			return nil, true, err
		}
		// The predicate ends at the label expression's last consumed token, which is
		// the closing ')' of a parenthesized label: the label node itself stays
		// transparent (its Loc excludes the parentheses, as the grammar says).
		return &ast.GraphPredicate{
			Kind:  ast.GraphIsLabeled,
			Not:   not,
			Left:  left,
			Right: label,
			Loc:   ast.Loc{Start: nodeLoc(left).Start, End: p.prev.Loc.End},
		}, true, nil
	default:
		return nil, false, nil
	}
	p.advance() // SOURCE | DESTINATION
	if p.cur.Type == kwOF {
		p.advance() // OF
	}
	right, err := p.parseBinaryExpr(bpBitOr)
	if err != nil {
		return nil, true, err
	}
	return &ast.GraphPredicate{
		Kind:  kind,
		Not:   not,
		Left:  left,
		Right: right,
		Loc:   spanNodes(left, right),
	}, true, nil
}

// parseValueGraphSubquery parses `VALUE [hint] { <graph subquery> }`
// (expression_subquery_with_keyword's VALUE braced_graph_subquery). VALUE is
// non-reserved, so a VALUE that is not followed by a brace or a hint stays an
// identifier (the caller handles that).
func (p *Parser) parseValueGraphSubquery() (ast.Node, error) {
	valTok := p.advance() // VALUE
	if p.cur.Type == int('@') {
		if herr := p.skipHint(); herr != nil {
			return nil, herr
		}
	}
	if p.cur.Type != int('{') {
		return nil, p.syntaxErrorAtCur()
	}
	sub, err := p.parseGraphBracedSubquery()
	if err != nil {
		return nil, err
	}
	sub.Modifier = ast.GraphSubqueryValue
	sub.Loc.Start = valTok.Loc.Start
	return sub, nil
}
