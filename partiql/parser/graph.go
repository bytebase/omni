package parser

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/partiql/ast"
)

// ---------------------------------------------------------------------------
// PartiQL Graph Pattern Matching Language (GPML).
//
// MATCH is context-sensitive: it only appears INSIDE parentheses, in the form
//
//	( lhsExpr MATCH gpmlPatternList )
//
// The opening paren and the lhs expression are parsed by parseParenExpr in
// exprprimary.go; when it sees tokMATCH after the inner expression it calls
// parseGraphMatch(first), which owns everything from MATCH through the closing
// paren.
//
// Grammar (PartiQLParser.g4):
//
//	exprGraphMatchMany   : PAREN_LEFT exprPrimary MATCH gpmlPatternList PAREN_RIGHT ;  (625-626)
//	exprGraphMatchOne    : exprPrimary MATCH gpmlPattern ;                              (628-629)
//	gpmlPattern          : selector=matchSelector? matchPattern ;                       (315-316)
//	gpmlPatternList      : selector=matchSelector? matchPattern ( COMMA matchPattern )*; (318-319)
//	matchPattern         : restrictor=patternRestrictor? variable=patternPathVariable? graphPart*; (321-322)
//	graphPart            : node | edge | pattern ;                                       (324-328)
//	matchSelector        : mod=(ANY|ALL) SHORTEST          # SelectorBasic
//	                     | ANY k=LITERAL_INTEGER?          # SelectorAny
//	                     | SHORTEST k=LITERAL_INTEGER GROUP? # SelectorShortest ;        (330-334)
//	patternPathVariable  : symbolPrimitive EQ ;                                          (336-337)
//	patternRestrictor    : restrictor=IDENTIFIER ;  // TRAIL / ACYCLIC / SIMPLE          (339-340)
//	node                 : PAREN_LEFT symbolPrimitive? patternPartLabel? whereClause? PAREN_RIGHT; (342-343)
//	edge                 : edgeWSpec quantifier=patternQuantifier?   # EdgeWithSpec
//	                     | edgeAbbrev quantifier=patternQuantifier?  # EdgeAbbreviated ; (345-348)
//	pattern              : PAREN_LEFT  restrictor? variable? graphPart+ where=whereClause? PAREN_RIGHT  quantifier?
//	                     | BRACKET_LEFT restrictor? variable? graphPart+ where=whereClause? BRACKET_RIGHT quantifier? ; (350-353)
//	patternQuantifier    : quant=( PLUS | ASTERISK )
//	                     | BRACE_LEFT lower=LITERAL_INTEGER COMMA upper=LITERAL_INTEGER? BRACE_RIGHT ; (355-358)
//	edgeWSpec            : MINUS edgeSpec MINUS ANGLE_RIGHT            # EdgeSpecRight
//	                     | TILDE edgeSpec TILDE                       # EdgeSpecUndirected
//	                     | ANGLE_LEFT MINUS edgeSpec MINUS            # EdgeSpecLeft
//	                     | TILDE edgeSpec TILDE ANGLE_RIGHT           # EdgeSpecUndirectedRight
//	                     | ANGLE_LEFT TILDE edgeSpec TILDE            # EdgeSpecUndirectedLeft
//	                     | ANGLE_LEFT MINUS edgeSpec MINUS ANGLE_RIGHT # EdgeSpecBidirectional
//	                     | MINUS edgeSpec MINUS                       # EdgeSpecUndirectedBidirectional ; (360-368)
//	edgeSpec            : BRACKET_LEFT symbolPrimitive? patternPartLabel? whereClause? BRACKET_RIGHT; (370-371)
//	patternPartLabel   : COLON symbolPrimitive ;                                         (373-374)
//	edgeAbbrev         : TILDE | TILDE ANGLE_RIGHT | ANGLE_LEFT TILDE | ANGLE_LEFT? MINUS ANGLE_RIGHT? ; (376-381)
//
// NOTE on `<`/`>`: the lexer scans `->`, `<-`, `~>`, `<~`, `<->` as separate
// single-character tokens (tokMINUS/tokLT/tokGT/tokTILDE), exactly matching the
// ANTLR grammar (MINUS, ANGLE_LEFT, ANGLE_RIGHT, TILDE). `<<`/`>>` are greedily
// lexed as bag-literal tokens, but those never appear in edge syntax.
// ---------------------------------------------------------------------------

// graphRestrictorNames maps the (uppercased) identifier spellings recognised as
// pattern restrictors to their AST enum. ANTLR's patternRestrictor rule admits
// ANY identifier, but only TRAIL/ACYCLIC/SIMPLE are documented and the AST's
// PatternRestrictor enum can only represent those three; see the divergence note
// in the PR body.
var graphRestrictorNames = map[string]ast.PatternRestrictor{
	"TRAIL":   ast.PatternRestrictorTrail,
	"ACYCLIC": ast.PatternRestrictorAcyclic,
	"SIMPLE":  ast.PatternRestrictorSimple,
}

// parseGraphMatch parses the MATCH keyword through the (already-open) closing
// paren of `( lhs MATCH gpmlPatternList )`. On entry, p.cur is tokMATCH and
// `lhs` is the already-parsed left-hand graph expression. On success p.cur is
// the closing PAREN_RIGHT (left for parseParenExpr to consume), matching the
// other branches of parseParenExpr which all stop on the closing paren.
//
// Grammar: exprGraphMatchMany (PartiQLParser.g4:625-626) + gpmlPatternList.
func (p *Parser) parseGraphMatch(lhs ast.ExprNode) (ast.ExprNode, error) {
	start := lhs.GetLoc().Start
	if _, err := p.expect(tokMATCH); err != nil {
		return nil, err
	}

	// gpmlPatternList: selector=matchSelector? matchPattern ( COMMA matchPattern )*
	// The optional selector binds before the first matchPattern, so it is
	// attached to the first GraphPattern's Selector field.
	selector, err := p.parseMatchSelector()
	if err != nil {
		return nil, err
	}

	patterns := make([]*ast.GraphPattern, 0, 1)
	first, err := p.parseMatchPattern()
	if err != nil {
		return nil, err
	}
	first.Selector = selector
	patterns = append(patterns, first)

	for p.cur.Type == tokCOMMA {
		p.advance() // consume COMMA
		gp, err := p.parseMatchPattern()
		if err != nil {
			return nil, err
		}
		patterns = append(patterns, gp)
	}

	// The closing paren is NOT consumed here: parseParenExpr consumes it after
	// we return, exactly like its valueList / plain-expr branches. We do peek
	// to give a graph-specific error if the pattern list ended early.
	if p.cur.Type != tokPAREN_RIGHT {
		return nil, &ParseError{
			Message: fmt.Sprintf("expected ')' to close graph MATCH expression, got %q", p.cur.Str),
			Loc:     p.cur.Loc,
		}
	}

	return &ast.MatchExpr{
		Expr:     lhs,
		Patterns: patterns,
		Loc:      ast.Loc{Start: start, End: p.cur.Loc.End},
	}, nil
}

// parseMatchSelector parses the optional matchSelector that may precede a
// pattern list. Returns nil when no selector is present.
//
// Grammar: matchSelector (PartiQLParser.g4:330-334):
//
//	mod=(ANY|ALL) SHORTEST              # SelectorBasic   -> ALL_SHORTEST
//	ANY k=LITERAL_INTEGER?              # SelectorAny     -> ANY (k ignored when present)
//	SHORTEST k=LITERAL_INTEGER GROUP?   # SelectorShortest-> SHORTEST_K
func (p *Parser) parseMatchSelector() (*ast.PatternSelector, error) {
	startLoc := p.cur.Loc
	switch p.cur.Type {
	case tokANY:
		p.advance()
		// ANY SHORTEST (SelectorBasic) vs ANY [k] (SelectorAny).
		if p.cur.Type == tokSHORTEST {
			p.advance()
			return &ast.PatternSelector{
				Kind: ast.SelectorKindAllShortest,
				Loc:  ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End},
			}, nil
		}
		// ANY k=LITERAL_INTEGER? — the integer is part of the selector but the
		// AST PatternSelector carries K only for SHORTEST_K; for ANY it is a
		// count that does not change the kind, so we consume and discard it.
		if p.cur.Type == tokICONST {
			p.advance()
		}
		return &ast.PatternSelector{
			Kind: ast.SelectorKindAny,
			Loc:  ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End},
		}, nil

	case tokALL:
		// ALL is only valid as `ALL SHORTEST` (SelectorBasic). A bare ALL is a
		// syntax error in GPML selector position.
		p.advance()
		if _, err := p.expect(tokSHORTEST); err != nil {
			return nil, &ParseError{
				Message: fmt.Sprintf("expected SHORTEST after ALL in graph selector, got %q", p.cur.Str),
				Loc:     p.cur.Loc,
			}
		}
		return &ast.PatternSelector{
			Kind: ast.SelectorKindAllShortest,
			Loc:  ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End},
		}, nil

	case tokSHORTEST:
		// SHORTEST k=LITERAL_INTEGER GROUP? — k is REQUIRED here (the ANY-less
		// form). GROUP is parsed-and-discarded (no AST slot, no semantics).
		p.advance()
		if p.cur.Type != tokICONST {
			return nil, &ParseError{
				Message: fmt.Sprintf("expected an integer count after SHORTEST in graph selector, got %q", p.cur.Str),
				Loc:     p.cur.Loc,
			}
		}
		k, err := parseIntLiteral(p.cur.Str)
		if err != nil {
			return nil, &ParseError{Message: err.Error(), Loc: p.cur.Loc}
		}
		p.advance()
		p.match(tokGROUP) // optional GROUP, no AST representation
		return &ast.PatternSelector{
			Kind: ast.SelectorKindShortestK,
			K:    k,
			Loc:  ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End},
		}, nil

	default:
		return nil, nil // no selector
	}
}

// parseMatchPattern parses one matchPattern:
//
//	restrictor=patternRestrictor? variable=patternPathVariable? graphPart*
//
// At the top level (inside the MATCH list) graphPart is `*` (zero or more); the
// grouped sub-pattern form (parsePatternGroup) requires `+` and is handled there.
//
// Grammar: matchPattern (PartiQLParser.g4:321-322).
func (p *Parser) parseMatchPattern() (*ast.GraphPattern, error) {
	startLoc := p.cur.Loc
	gp := &ast.GraphPattern{}

	// restrictor? — a bare IDENT (TRAIL/ACYCLIC/SIMPLE) NOT followed by EQ.
	if r, ok := p.tryRestrictor(); ok {
		gp.Restrictor = r
	}

	// variable? — symbolPrimitive EQ (path variable binding).
	if v, ok, err := p.tryPathVariable(); err != nil {
		return nil, err
	} else if ok {
		gp.Variable = v
	}

	// graphPart* — node | edge | pattern.
	parts, err := p.parseGraphParts(false /* requireOne */)
	if err != nil {
		return nil, err
	}
	gp.Parts = parts

	end := p.prev.Loc.End
	if len(parts) == 0 {
		end = startLoc.End
	}
	gp.Loc = ast.Loc{Start: startLoc.Start, End: end}
	return gp, nil
}

// tryRestrictor consumes a leading restrictor keyword (TRAIL/ACYCLIC/SIMPLE,
// case-insensitive) when present and not immediately followed by EQ (which
// would make it a path-variable name). Returns (restrictor, true) on match.
//
// Grammar: patternRestrictor (PartiQLParser.g4:339-340). ANTLR admits any
// IDENTIFIER as a restrictor; we restrict to the three documented spellings the
// AST can represent (divergence flagged in the PR body).
func (p *Parser) tryRestrictor() (ast.PatternRestrictor, bool) {
	if p.cur.Type != tokIDENT {
		return ast.PatternRestrictorNone, false
	}
	if p.peekNext().Type == tokEQ {
		return ast.PatternRestrictorNone, false // it's a path variable name
	}
	r, ok := graphRestrictorNames[strings.ToUpper(p.cur.Str)]
	if !ok {
		return ast.PatternRestrictorNone, false
	}
	p.advance()
	return r, true
}

// tryPathVariable consumes a path-variable binding `symbolPrimitive EQ` when the
// current tokens form one. Returns (varRef, true, nil) on match, (nil, false,
// nil) when absent.
//
// Grammar: patternPathVariable (PartiQLParser.g4:336-337).
func (p *Parser) tryPathVariable() (*ast.VarRef, bool, error) {
	if p.cur.Type != tokIDENT && p.cur.Type != tokIDENT_QUOTED {
		return nil, false, nil
	}
	if p.peekNext().Type != tokEQ {
		return nil, false, nil
	}
	name, caseSensitive, nameLoc, err := p.parseSymbolPrimitive()
	if err != nil {
		return nil, false, err
	}
	if _, err := p.expect(tokEQ); err != nil {
		return nil, false, err
	}
	return &ast.VarRef{
		Name:          name,
		CaseSensitive: caseSensitive,
		Loc:           nameLoc,
	}, true, nil
}

// parseGraphParts parses a run of graphParts (node | edge | pattern). When
// requireOne is true at least one part must be present (the grouped sub-pattern
// form `graphPart+`); otherwise zero parts are allowed (top-level matchPattern
// `graphPart*`).
//
// Grammar: graphPart (PartiQLParser.g4:324-328).
func (p *Parser) parseGraphParts(requireOne bool) ([]ast.PatternNode, error) {
	var parts []ast.PatternNode
	for {
		switch {
		case p.cur.Type == tokPAREN_LEFT:
			// `(` begins either a node or a parenthesised sub-pattern group.
			part, err := p.parseParenGraphPart()
			if err != nil {
				return nil, err
			}
			parts = append(parts, part)
		case p.cur.Type == tokBRACKET_LEFT:
			// `[` begins a bracketed sub-pattern group.
			part, err := p.parsePatternGroup(tokBRACKET_LEFT, tokBRACKET_RIGHT)
			if err != nil {
				return nil, err
			}
			parts = append(parts, part)
		case p.isEdgeStart():
			// An edge must be preceded by a node/group; reject a leading edge.
			if len(parts) == 0 {
				return nil, &ParseError{
					Message: "graph pattern cannot begin with an edge",
					Loc:     p.cur.Loc,
				}
			}
			edge, err := p.parseEdge()
			if err != nil {
				return nil, err
			}
			parts = append(parts, edge)
		default:
			// No more graphParts.
			if requireOne && len(parts) == 0 {
				return nil, &ParseError{
					Message: fmt.Sprintf("expected a node, edge, or sub-pattern in graph pattern, got %q", p.cur.Str),
					Loc:     p.cur.Loc,
				}
			}
			return parts, nil
		}
	}
}

// parseParenGraphPart disambiguates a `(` between a node pattern and a
// parenthesised sub-pattern group, then parses the chosen form.
//
// A node body is `symbolPrimitive? patternPartLabel? whereClause?` followed by
// `)`. A group body is `restrictor? variable? graphPart+ where?` — its first
// significant token is `(`/`[` (a nested node/group), or an IDENT that is a
// restrictor / path-variable name preceding such a nested part. The two are
// distinguished by bounded lookahead at the token after `(`.
func (p *Parser) parseParenGraphPart() (ast.PatternNode, error) {
	if p.parenStartsNode() {
		return p.parseNode()
	}
	return p.parsePatternGroup(tokPAREN_LEFT, tokPAREN_RIGHT)
}

// parenStartsNode reports whether the `(` at p.cur opens a node (as opposed to
// a sub-pattern group). It does not consume input. The decision uses the two
// tokens after `(`.
//
// node     : ( symbolPrimitive? patternPartLabel? whereClause? )
// pattern  : ( restrictor? variable? graphPart+ ... )   // graphPart starts with ( or [
func (p *Parser) parenStartsNode() bool {
	inner := p.peekNext() // token after '('
	switch inner.Type {
	case tokPAREN_RIGHT:
		return true // () empty node
	case tokCOLON:
		return true // (:Label) — label-only node
	case tokWHERE:
		return true // (WHERE ...) — degenerate but node-shaped
	case tokPAREN_LEFT, tokBRACKET_LEFT:
		return false // ((...)) / ([...]) — group whose first part is nested
	case tokIDENT, tokIDENT_QUOTED:
		// Need the token after the identifier to decide.
		after := p.peekAfterParenIdent()
		switch after {
		case tokPAREN_RIGHT, tokCOLON, tokWHERE:
			return true // (a) / (a:Label) / (a WHERE ...) — node with a variable
		default:
			// (a = ...) path variable, or (TRAIL (x)...) restrictor, or
			// (a (x)...) — all are sub-pattern groups (a leading IDENT that is
			// not a node body).
			return false
		}
	default:
		// Anything else after '(' is not a valid node body; treat as a group so
		// the group parser produces a precise graphPart error.
		return false
	}
}

// peekAfterParenIdent returns the token type two positions after p.cur (which is
// `(`): i.e. the token following the identifier that follows `(`. It uses
// save/restore so it consumes no input.
func (p *Parser) peekAfterParenIdent() int {
	st := p.save()
	p.advance() // consume '('
	p.advance() // consume the identifier
	t := p.cur.Type
	p.restore(st)
	return t
}

// parseNode parses a node pattern.
//
//	node : PAREN_LEFT symbolPrimitive? patternPartLabel? whereClause? PAREN_RIGHT
//
// Grammar: node (PartiQLParser.g4:342-343).
func (p *Parser) parseNode() (*ast.NodePattern, error) {
	start := p.cur.Loc.Start
	if _, err := p.expect(tokPAREN_LEFT); err != nil {
		return nil, err
	}
	n := &ast.NodePattern{}

	// symbolPrimitive? (the node variable).
	if p.cur.Type == tokIDENT || p.cur.Type == tokIDENT_QUOTED {
		name, caseSensitive, nameLoc, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		n.Variable = &ast.VarRef{Name: name, CaseSensitive: caseSensitive, Loc: nameLoc}
	}

	// patternPartLabel?
	if p.cur.Type == tokCOLON {
		label, err := p.parsePatternPartLabel()
		if err != nil {
			return nil, err
		}
		n.Labels = []string{label}
	}

	// whereClause?
	if p.cur.Type == tokWHERE {
		where, err := p.parseGraphWhere()
		if err != nil {
			return nil, err
		}
		n.Where = where
	}

	rp, err := p.expect(tokPAREN_RIGHT)
	if err != nil {
		return nil, err
	}
	n.Loc = ast.Loc{Start: start, End: rp.Loc.End}
	return n, nil
}

// parsePatternGroup parses a parenthesised or bracketed sub-pattern group:
//
//	( restrictor? variable? graphPart+ where=whereClause? ) quantifier?
//	[ restrictor? variable? graphPart+ where=whereClause? ] quantifier?
//
// It is represented as a nested *ast.GraphPattern (which satisfies PatternNode).
// open/close are the matching delimiter token types.
//
// Grammar: pattern (PartiQLParser.g4:350-353).
func (p *Parser) parsePatternGroup(open, closeTok int) (*ast.GraphPattern, error) {
	start := p.cur.Loc.Start
	if _, err := p.expect(open); err != nil {
		return nil, err
	}
	gp := &ast.GraphPattern{}

	if r, ok := p.tryRestrictor(); ok {
		gp.Restrictor = r
	}
	if v, ok, err := p.tryPathVariable(); err != nil {
		return nil, err
	} else if ok {
		gp.Variable = v
	}

	// graphPart+ (at least one).
	parts, err := p.parseGraphParts(true /* requireOne */)
	if err != nil {
		return nil, err
	}
	gp.Parts = parts

	// where=whereClause? — a filter over the whole sub-pattern group, stored on
	// GraphPattern.Where (a nested-group-only field).
	if p.cur.Type == tokWHERE {
		where, err := p.parseGraphWhere()
		if err != nil {
			return nil, err
		}
		gp.Where = where
	}

	end, err := p.expect(closeTok)
	if err != nil {
		return nil, err
	}
	gp.Loc = ast.Loc{Start: start, End: end.Loc.End}

	// quantifier? on the whole group, stored on GraphPattern.Quantifier.
	if p.peekQuantifierStart() {
		q, err := p.parseQuantifier()
		if err != nil {
			return nil, err
		}
		gp.Quantifier = q
		gp.Loc.End = q.Loc.End
	}
	return gp, nil
}

// parsePatternPartLabel parses `COLON symbolPrimitive` and returns the label
// name.
//
// Grammar: patternPartLabel (PartiQLParser.g4:373-374).
func (p *Parser) parsePatternPartLabel() (string, error) {
	if _, err := p.expect(tokCOLON); err != nil {
		return "", err
	}
	name, _, _, err := p.parseSymbolPrimitive()
	if err != nil {
		return "", err
	}
	return name, nil
}

// parseGraphWhere parses `WHERE expr` (the whereClause used inside node and edge
// patterns). Uses the plain expression entry point — graph whereClauses bind a
// full expr, not an exprSelect.
//
// Grammar: whereClause (PartiQLParser.g4:207-208).
func (p *Parser) parseGraphWhere() (ast.ExprNode, error) {
	if _, err := p.expect(tokWHERE); err != nil {
		return nil, err
	}
	return p.parseExprTop()
}

// isEdgeStart reports whether p.cur could begin an edge (edgeWSpec or
// edgeAbbrev). Edges start with MINUS, TILDE, or ANGLE_LEFT (`<`).
//
// Grammar: edgeWSpec / edgeAbbrev (PartiQLParser.g4:360-381).
func (p *Parser) isEdgeStart() bool {
	switch p.cur.Type {
	case tokMINUS, tokTILDE, tokLT:
		return true
	default:
		return false
	}
}

// parseEdge parses an edge (with or without a bracketed spec) plus an optional
// trailing quantifier.
//
//	edge : edgeWSpec quantifier?      # EdgeWithSpec
//	     | edgeAbbrev quantifier?     # EdgeAbbreviated
//
// Grammar: edge (PartiQLParser.g4:345-348).
func (p *Parser) parseEdge() (*ast.EdgePattern, error) {
	start := p.cur.Loc.Start
	edge, err := p.parseEdgeShape()
	if err != nil {
		return nil, err
	}
	edge.Loc.Start = start

	// quantifier? on the edge.
	if p.peekQuantifierStart() {
		q, err := p.parseQuantifier()
		if err != nil {
			return nil, err
		}
		edge.Quantifier = q
	}
	edge.Loc.End = p.prev.Loc.End
	return edge, nil
}

// parseEdgeShape parses the directional shape of an edge, branching between the
// bracketed (edgeWSpec) and abbreviated (edgeAbbrev) forms. The returned
// EdgePattern has Direction (and, for edgeWSpec, the body fields) populated; Loc
// is partially filled (End set, Start fixed up by the caller).
func (p *Parser) parseEdgeShape() (*ast.EdgePattern, error) {
	switch p.cur.Type {
	case tokMINUS:
		// MINUS ... — either `-[spec]-(>)` (edgeWSpec) or `-`/`->` (edgeAbbrev).
		p.advance() // consume '-'
		if p.cur.Type == tokBRACKET_LEFT {
			// edgeWSpec: MINUS edgeSpec MINUS ANGLE_RIGHT?  (Right / UndirectedBidirectional)
			return p.parseEdgeWSpecAfterMinus()
		}
		// edgeAbbrev: MINUS ANGLE_RIGHT?  -> `->` or `-`
		if p.cur.Type == tokGT {
			p.advance()
			return &ast.EdgePattern{Direction: ast.EdgeDirRight, Loc: ast.Loc{End: p.prev.Loc.End}}, nil
		}
		return &ast.EdgePattern{Direction: ast.EdgeDirUndirectedBidirectional, Loc: ast.Loc{End: p.prev.Loc.End}}, nil

	case tokTILDE:
		// TILDE ... — `~[spec]~(>)` (edgeWSpec) or `~`/`~>` (edgeAbbrev).
		p.advance() // consume '~'
		if p.cur.Type == tokBRACKET_LEFT {
			return p.parseEdgeWSpecAfterTilde()
		}
		// edgeAbbrev: TILDE ANGLE_RIGHT?  -> `~>` or `~`
		if p.cur.Type == tokGT {
			p.advance()
			return &ast.EdgePattern{Direction: ast.EdgeDirRightOrUndirected, Loc: ast.Loc{End: p.prev.Loc.End}}, nil
		}
		return &ast.EdgePattern{Direction: ast.EdgeDirUndirected, Loc: ast.Loc{End: p.prev.Loc.End}}, nil

	case tokLT:
		// ANGLE_LEFT ... — `<-[spec]-(>)`, `<~[spec]~` (edgeWSpec) or `<-`,
		// `<->`, `<~` (edgeAbbrev).
		p.advance() // consume '<'
		switch p.cur.Type {
		case tokMINUS:
			p.advance() // consume '-'
			if p.cur.Type == tokBRACKET_LEFT {
				// edgeWSpec: ANGLE_LEFT MINUS edgeSpec MINUS ANGLE_RIGHT?  (Left / Bidirectional)
				return p.parseEdgeWSpecAfterLeftMinus()
			}
			// edgeAbbrev: ANGLE_LEFT MINUS ANGLE_RIGHT?  -> `<-` or `<->`
			if p.cur.Type == tokGT {
				p.advance()
				return &ast.EdgePattern{Direction: ast.EdgeDirLeftOrRight, Loc: ast.Loc{End: p.prev.Loc.End}}, nil
			}
			return &ast.EdgePattern{Direction: ast.EdgeDirLeft, Loc: ast.Loc{End: p.prev.Loc.End}}, nil
		case tokTILDE:
			p.advance() // consume '~'
			if p.cur.Type == tokBRACKET_LEFT {
				// edgeWSpec: ANGLE_LEFT TILDE edgeSpec TILDE  (UndirectedLeft)
				return p.parseEdgeWSpecAfterLeftTilde()
			}
			// edgeAbbrev: ANGLE_LEFT TILDE  -> `<~`
			return &ast.EdgePattern{Direction: ast.EdgeDirLeftOrUndirected, Loc: ast.Loc{End: p.prev.Loc.End}}, nil
		default:
			return nil, &ParseError{
				Message: fmt.Sprintf("expected '-' or '~' after '<' in edge, got %q", p.cur.Str),
				Loc:     p.cur.Loc,
			}
		}

	default:
		return nil, &ParseError{
			Message: fmt.Sprintf("expected an edge ('-', '<', or '~'), got %q", p.cur.Str),
			Loc:     p.cur.Loc,
		}
	}
}

// parseEdgeWSpecAfterMinus handles `-[spec]-` then optional `>`:
//
//	MINUS edgeSpec MINUS              # EdgeSpecUndirectedBidirectional  (`-[]-`)
//	MINUS edgeSpec MINUS ANGLE_RIGHT  # EdgeSpecRight                    (`-[]->`)
//
// On entry the leading MINUS is already consumed and p.cur is BRACKET_LEFT.
func (p *Parser) parseEdgeWSpecAfterMinus() (*ast.EdgePattern, error) {
	edge, err := p.parseEdgeSpec()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokMINUS); err != nil {
		return nil, err
	}
	if p.cur.Type == tokGT {
		p.advance()
		edge.Direction = ast.EdgeDirRight
	} else {
		edge.Direction = ast.EdgeDirUndirectedBidirectional
	}
	edge.Loc.End = p.prev.Loc.End
	return edge, nil
}

// parseEdgeWSpecAfterTilde handles `~[spec]~` then optional `>`:
//
//	TILDE edgeSpec TILDE              # EdgeSpecUndirected       (`~[]~`)
//	TILDE edgeSpec TILDE ANGLE_RIGHT  # EdgeSpecUndirectedRight  (`~[]~>`)
//
// On entry the leading TILDE is already consumed and p.cur is BRACKET_LEFT.
func (p *Parser) parseEdgeWSpecAfterTilde() (*ast.EdgePattern, error) {
	edge, err := p.parseEdgeSpec()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokTILDE); err != nil {
		return nil, err
	}
	if p.cur.Type == tokGT {
		p.advance()
		edge.Direction = ast.EdgeDirRightOrUndirected
	} else {
		edge.Direction = ast.EdgeDirUndirected
	}
	edge.Loc.End = p.prev.Loc.End
	return edge, nil
}

// parseEdgeWSpecAfterLeftMinus handles `<-[spec]-` then optional `>`:
//
//	ANGLE_LEFT MINUS edgeSpec MINUS              # EdgeSpecLeft          (`<-[]-`)
//	ANGLE_LEFT MINUS edgeSpec MINUS ANGLE_RIGHT  # EdgeSpecBidirectional (`<-[]->`)
//
// On entry `<` and `-` are already consumed and p.cur is BRACKET_LEFT.
func (p *Parser) parseEdgeWSpecAfterLeftMinus() (*ast.EdgePattern, error) {
	edge, err := p.parseEdgeSpec()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokMINUS); err != nil {
		return nil, err
	}
	if p.cur.Type == tokGT {
		p.advance()
		edge.Direction = ast.EdgeDirLeftOrRight
	} else {
		edge.Direction = ast.EdgeDirLeft
	}
	edge.Loc.End = p.prev.Loc.End
	return edge, nil
}

// parseEdgeWSpecAfterLeftTilde handles `<~[spec]~`:
//
//	ANGLE_LEFT TILDE edgeSpec TILDE  # EdgeSpecUndirectedLeft  (`<~[]~`)
//
// On entry `<` and `~` are already consumed and p.cur is BRACKET_LEFT.
func (p *Parser) parseEdgeWSpecAfterLeftTilde() (*ast.EdgePattern, error) {
	edge, err := p.parseEdgeSpec()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(tokTILDE); err != nil {
		return nil, err
	}
	edge.Direction = ast.EdgeDirLeftOrUndirected
	edge.Loc.End = p.prev.Loc.End
	return edge, nil
}

// parseEdgeSpec parses the bracketed body of an edge:
//
//	edgeSpec : BRACKET_LEFT symbolPrimitive? patternPartLabel? whereClause? BRACKET_RIGHT
//
// It returns an EdgePattern with Variable/Labels/Where populated (Direction and
// Loc are filled by the caller). On entry p.cur is BRACKET_LEFT.
//
// Grammar: edgeSpec (PartiQLParser.g4:370-371).
func (p *Parser) parseEdgeSpec() (*ast.EdgePattern, error) {
	if _, err := p.expect(tokBRACKET_LEFT); err != nil {
		return nil, err
	}
	edge := &ast.EdgePattern{}

	if p.cur.Type == tokIDENT || p.cur.Type == tokIDENT_QUOTED {
		name, caseSensitive, nameLoc, err := p.parseSymbolPrimitive()
		if err != nil {
			return nil, err
		}
		edge.Variable = &ast.VarRef{Name: name, CaseSensitive: caseSensitive, Loc: nameLoc}
	}
	if p.cur.Type == tokCOLON {
		label, err := p.parsePatternPartLabel()
		if err != nil {
			return nil, err
		}
		edge.Labels = []string{label}
	}
	if p.cur.Type == tokWHERE {
		where, err := p.parseGraphWhere()
		if err != nil {
			return nil, err
		}
		edge.Where = where
	}
	if _, err := p.expect(tokBRACKET_RIGHT); err != nil {
		return nil, err
	}
	return edge, nil
}

// peekQuantifierStart reports whether p.cur begins a patternQuantifier
// (`+`, `*`, or `{`).
//
// Grammar: patternQuantifier (PartiQLParser.g4:355-358).
func (p *Parser) peekQuantifierStart() bool {
	switch p.cur.Type {
	case tokPLUS, tokASTERISK, tokBRACE_LEFT:
		return true
	default:
		return false
	}
}

// parseQuantifier parses a patternQuantifier:
//
//	quant=( PLUS | ASTERISK )                                       -> +:{1,-1}, *:{0,-1}
//	BRACE_LEFT lower=LITERAL_INTEGER COMMA upper=LITERAL_INTEGER? BRACE_RIGHT  -> {m,n} / {m,}
//
// Grammar: patternQuantifier (PartiQLParser.g4:355-358).
func (p *Parser) parseQuantifier() (*ast.PatternQuantifier, error) {
	start := p.cur.Loc.Start
	switch p.cur.Type {
	case tokPLUS:
		p.advance()
		return &ast.PatternQuantifier{Min: 1, Max: -1, Loc: ast.Loc{Start: start, End: p.prev.Loc.End}}, nil
	case tokASTERISK:
		p.advance()
		return &ast.PatternQuantifier{Min: 0, Max: -1, Loc: ast.Loc{Start: start, End: p.prev.Loc.End}}, nil
	case tokBRACE_LEFT:
		p.advance() // consume '{'
		if p.cur.Type != tokICONST {
			return nil, &ParseError{
				Message: fmt.Sprintf("expected a lower bound integer in quantifier, got %q", p.cur.Str),
				Loc:     p.cur.Loc,
			}
		}
		lower, err := parseIntLiteral(p.cur.Str)
		if err != nil {
			return nil, &ParseError{Message: err.Error(), Loc: p.cur.Loc}
		}
		p.advance()
		if _, err := p.expect(tokCOMMA); err != nil {
			return nil, err
		}
		upper := -1 // {m,} == unbounded
		if p.cur.Type == tokICONST {
			u, err := parseIntLiteral(p.cur.Str)
			if err != nil {
				return nil, &ParseError{Message: err.Error(), Loc: p.cur.Loc}
			}
			upper = u
			p.advance()
		}
		end, err := p.expect(tokBRACE_RIGHT)
		if err != nil {
			return nil, err
		}
		return &ast.PatternQuantifier{Min: lower, Max: upper, Loc: ast.Loc{Start: start, End: end.Loc.End}}, nil
	default:
		return nil, &ParseError{
			Message: fmt.Sprintf("expected a quantifier ('+', '*', or '{m,n}'), got %q", p.cur.Str),
			Loc:     p.cur.Loc,
		}
	}
}
