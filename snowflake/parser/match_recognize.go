package parser

// MATCH_RECOGNIZE table clause (T5.3):
//
//	MATCH_RECOGNIZE (
//	  [ PARTITION BY <expr>, … ]
//	  [ ORDER BY <expr> [ASC|DESC] [NULLS …], … ]
//	  [ MEASURES [FINAL|RUNNING] <expr> [AS] <alias>, … ]
//	  [ { ONE ROW | ALL ROWS } PER MATCH
//	      [ SHOW EMPTY MATCHES | OMIT EMPTY MATCHES | WITH UNMATCHED ROWS ] ]
//	  [ AFTER MATCH SKIP { PAST LAST ROW | TO NEXT ROW | TO [FIRST|LAST] <var> } ]
//	  PATTERN ( <row_pattern> )
//	  [ DEFINE <var> AS <cond>, … ]
//	) [ [AS] alias ]
//
// MEASURES and DEFINE expressions are parsed by the shared expression parser.
// The PATTERN body is captured as raw text (RowPattern.Raw): the full
// row-pattern mini-grammar (quantifiers +, *, {n,m}, alternation |, grouping,
// PERMUTE, anchors ^ $) has no executable oracle to validate a structured
// tree against, and consumers re-lex the raw text when they need it. The
// legacy SnowflakeParser.g4 likewise leaves PATTERN/DEFINE as placeholders
// (its `symbol: DUMMY` rule never matches real input), so the docs + corpus
// are authoritative here.

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// parseMatchRecognizeClause parses a MATCH_RECOGNIZE (...) clause. The caller
// has verified that p.cur is kwMATCH_RECOGNIZE.
func (p *Parser) parseMatchRecognizeClause() (*ast.MatchRecognizeClause, error) {
	mrTok := p.advance() // consume MATCH_RECOGNIZE
	clause := &ast.MatchRecognizeClause{
		Loc: ast.Loc{Start: mrTok.Loc.Start},
	}

	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	// PARTITION BY <expr_list>
	if p.cur.Type == kwPARTITION {
		p.advance() // consume PARTITION
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		exprs, err := p.parseExprList()
		if err != nil {
			return nil, err
		}
		clause.PartitionBy = exprs
	}

	// ORDER BY <order_items>
	if p.cur.Type == kwORDER {
		p.advance() // consume ORDER
		if _, err := p.expect(kwBY); err != nil {
			return nil, err
		}
		items, err := p.parseOrderByList()
		if err != nil {
			return nil, err
		}
		clause.OrderBy = items
	}

	// MEASURES <measure_list>
	if p.cur.Type == kwMEASURES {
		p.advance() // consume MEASURES
		measures, err := p.parseMeasureList()
		if err != nil {
			return nil, err
		}
		clause.Measures = measures
	}

	// { ONE ROW | ALL ROWS } PER MATCH [opt]
	if p.cur.Type == kwONE || p.cur.Type == kwALL {
		rpm, err := p.parseRowsPerMatch()
		if err != nil {
			return nil, err
		}
		clause.RowsPerMatch = rpm
	}

	// AFTER MATCH SKIP …
	if p.cur.Type == kwAFTER {
		am, err := p.parseAfterMatchSkip()
		if err != nil {
			return nil, err
		}
		clause.AfterMatch = am
	}

	// PATTERN ( <row_pattern> )
	if p.cur.Type == kwPATTERN {
		pat, err := p.parseRowPattern()
		if err != nil {
			return nil, err
		}
		clause.Pattern = pat
	}

	// DEFINE <define_list>
	if p.cur.Type == kwDEFINE {
		p.advance() // consume DEFINE
		defs, err := p.parseDefineList()
		if err != nil {
			return nil, err
		}
		clause.Define = defs
	}

	closeTok, err := p.expect(')')
	if err != nil {
		return nil, err
	}
	clause.Loc.End = closeTok.Loc.End

	// Optional trailing [AS] alias.
	if alias, has := p.parseOptionalAlias(); has {
		clause.Alias = alias
		clause.Loc.End = p.prev.Loc.End
	}

	return clause, nil
}

// parseMeasureList parses MEASURES items: [FINAL|RUNNING] <expr> [AS] <alias>.
func (p *Parser) parseMeasureList() ([]*ast.MatchMeasure, error) {
	var measures []*ast.MatchMeasure
	for {
		m, err := p.parseMeasure()
		if err != nil {
			return nil, err
		}
		measures = append(measures, m)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return measures, nil
}

// parseMeasure parses one MEASURES entry.
func (p *Parser) parseMeasure() (*ast.MatchMeasure, error) {
	startLoc := p.cur.Loc
	m := &ast.MatchMeasure{
		Loc: ast.Loc{Start: startLoc.Start},
	}

	// Optional FINAL / RUNNING semantics prefix. These are non-reserved
	// identifiers, so detect them by text. Treat them as the prefix only when
	// the following token actually begins a measure operand — not '(' (which
	// would make FINAL/RUNNING itself a function call) and not a terminator
	// such as AS / ',' / ')' (which would mean FINAL/RUNNING is the operand,
	// e.g. a column literally named "final").
	if p.cur.Type == tokIdent {
		if up := strings.ToUpper(p.cur.Str); (up == "FINAL" || up == "RUNNING") && measureOperandFollows(p.peekNext().Type) {
			p.advance()
			if up == "FINAL" {
				m.Semantics = ast.MatchSemanticsFinal
			} else {
				m.Semantics = ast.MatchSemanticsRunning
			}
		}
	}

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	m.Expr = expr
	m.Loc.End = ast.NodeLoc(expr).End

	// Alias (AS <alias> or bare <alias>). Required by Snowflake but parsed
	// optionally so a missing alias surfaces as a downstream error rather
	// than swallowing the next clause keyword.
	if alias, has := p.parseOptionalAlias(); has {
		m.Alias = alias
		m.Loc.End = p.prev.Loc.End
	}
	return m, nil
}

// measureOperandFollows reports whether a token type can begin the operand of
// a MEASURES item following a FINAL/RUNNING prefix. It returns false for '('
// (FINAL/RUNNING is then a function call), and for the terminators AS, ',',
// ')', and EOF (FINAL/RUNNING is then itself the operand, e.g. a column named
// "final").
func measureOperandFollows(tokType int) bool {
	switch tokType {
	case '(', ')', ',', kwAS, tokEOF:
		return false
	}
	return true
}

// parseRowsPerMatch parses { ONE ROW | ALL ROWS } PER MATCH [opt].
func (p *Parser) parseRowsPerMatch() (*ast.RowsPerMatch, error) {
	startTok := p.advance() // consume ONE or ALL
	rpm := &ast.RowsPerMatch{
		Loc: ast.Loc{Start: startTok.Loc.Start},
	}
	switch startTok.Type {
	case kwONE:
		rpm.Kind = ast.OneRowPerMatch
		if _, err := p.expect(kwROW); err != nil {
			return nil, err
		}
	case kwALL:
		rpm.Kind = ast.AllRowsPerMatch
		if _, err := p.expect(kwROWS); err != nil {
			return nil, err
		}
	}
	if _, err := p.expect(kwPER); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwMATCH); err != nil {
		return nil, err
	}
	rpm.Loc.End = p.prev.Loc.End

	// Optional modifier (only valid for ALL ROWS, but accepted generally).
	switch p.cur.Type {
	case kwSHOW:
		p.advance() // consume SHOW
		if _, err := p.expect(kwEMPTY); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwMATCHES); err != nil {
			return nil, err
		}
		rpm.Opt = ast.RowsPerMatchShowEmpty
		rpm.Loc.End = p.prev.Loc.End
	case kwOMIT:
		p.advance() // consume OMIT
		if _, err := p.expect(kwEMPTY); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwMATCHES); err != nil {
			return nil, err
		}
		rpm.Opt = ast.RowsPerMatchOmitEmpty
		rpm.Loc.End = p.prev.Loc.End
	case kwWITH:
		p.advance() // consume WITH
		if _, err := p.expect(kwUNMATCHED); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwROWS); err != nil {
			return nil, err
		}
		rpm.Opt = ast.RowsPerMatchWithUnmatched
		rpm.Loc.End = p.prev.Loc.End
	}
	return rpm, nil
}

// parseAfterMatchSkip parses AFTER MATCH SKIP { PAST LAST ROW | TO NEXT ROW |
// TO [FIRST|LAST] <var> }.
func (p *Parser) parseAfterMatchSkip() (*ast.AfterMatchSkip, error) {
	afterTok := p.advance() // consume AFTER
	am := &ast.AfterMatchSkip{
		Loc: ast.Loc{Start: afterTok.Loc.Start},
	}
	if _, err := p.expect(kwMATCH); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwSKIP); err != nil {
		return nil, err
	}

	switch p.cur.Type {
	case kwPAST:
		p.advance() // consume PAST
		if _, err := p.expect(kwLAST); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwROW); err != nil {
			return nil, err
		}
		am.Kind = ast.AfterMatchSkipPastLastRow

	case kwTO:
		p.advance() // consume TO
		switch p.cur.Type {
		case kwNEXT:
			p.advance() // consume NEXT
			if _, err := p.expect(kwROW); err != nil {
				return nil, err
			}
			am.Kind = ast.AfterMatchSkipToNextRow
		case kwFIRST:
			p.advance() // consume FIRST
			sym, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			am.Kind = ast.AfterMatchSkipToFirst
			am.Symbol = sym
		case kwLAST:
			p.advance() // consume LAST
			sym, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			am.Kind = ast.AfterMatchSkipToLast
			am.Symbol = sym
		default:
			// SKIP TO <var>
			sym, err := p.parseIdent()
			if err != nil {
				return nil, err
			}
			am.Kind = ast.AfterMatchSkipToVar
			am.Symbol = sym
		}

	default:
		return nil, &ParseError{
			Loc: p.cur.Loc,
			Msg: "expected PAST or TO after AFTER MATCH SKIP",
		}
	}

	am.Loc.End = p.prev.Loc.End
	return am, nil
}

// parseRowPattern parses PATTERN ( <row_pattern> ) and captures the inner
// pattern text verbatim.
//
// The pattern body is scanned at the *byte* level rather than the token level:
// the row-pattern mini-grammar contains characters (notably '?') that are not
// valid SQL tokens, and tokenizing through them would make the lexer emit
// spurious tokInvalid errors. Byte scanning tracks parenthesis depth (so a
// grouped sub-pattern "( … )" is captured whole) and skips single-quoted
// string literals (so a quoted literal containing ')' does not end the scan
// early). After the matching ')' is found, the lexer is repositioned past it.
func (p *Parser) parseRowPattern() (*ast.RowPattern, error) {
	patTok := p.advance() // consume PATTERN
	openTok, err := p.expect('(')
	if err != nil {
		return nil, err
	}

	// Absolute byte offset just after '('. p.input is indexed by local
	// (unshifted) offsets, so subtract p.base.
	innerStartAbs := openTok.Loc.End
	i := innerStartAbs - p.base // local index into p.input
	depth := 1
	for i < len(p.input) {
		c := p.input[i]
		switch c {
		case '\'':
			// Skip a single-quoted string literal, honoring '' escapes.
			i++
			for i < len(p.input) {
				if p.input[i] == '\'' {
					if i+1 < len(p.input) && p.input[i+1] == '\'' {
						i += 2
						continue
					}
					break
				}
				i++
			}
			// i now points at the closing quote (or end of input).
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				innerEndAbs := i + p.base                // absolute offset of the closing ')'
				raw := p.input[innerStartAbs-p.base : i] // inner text, local slice
				// Reposition the lexer just past the matching ')' and re-prime.
				p.repositionLexer(i + 1)
				return &ast.RowPattern{
					Raw:    strings.TrimSpace(raw),
					RawLoc: ast.Loc{Start: innerStartAbs, End: innerEndAbs},
					Loc:    ast.Loc{Start: patTok.Loc.Start, End: innerEndAbs + 1},
				}, nil
			}
		}
		i++
	}
	return nil, &ParseError{
		Loc: ast.Loc{Start: patTok.Loc.Start, End: patTok.Loc.End},
		Msg: "unterminated PATTERN (...) in MATCH_RECOGNIZE",
	}
}

// parseDefineList parses DEFINE items: <var> AS <cond>.
func (p *Parser) parseDefineList() ([]*ast.MatchDefine, error) {
	var defs []*ast.MatchDefine
	for {
		d, err := p.parseDefine()
		if err != nil {
			return nil, err
		}
		defs = append(defs, d)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return defs, nil
}

// parseDefine parses one DEFINE entry: <var> AS <condition>.
func (p *Parser) parseDefine() (*ast.MatchDefine, error) {
	sym, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	d := &ast.MatchDefine{
		Symbol: sym,
		Loc:    ast.Loc{Start: sym.Loc.Start},
	}
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}
	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	d.Cond = cond
	d.Loc.End = ast.NodeLoc(cond).End
	return d, nil
}
