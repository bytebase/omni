package parser

import "github.com/bytebase/omni/trino/ast"

// This file is the `parser-match-recognize` DAG node (with row_pattern.go): it
// implements Trino's row-pattern-recognition subsystem — the `patternRecognition`
// relation clause `MATCH_RECOGNIZE ( … )` and the window-frame row-pattern
// additions (the windowFrame rule's MEASURES / PATTERN / SUBSET / DEFINE / AFTER
// MATCH / INITIAL|SEEK extensions). row_pattern.go owns the rowPattern grammar
// itself; this file owns the clause that surrounds it and the two integration
// hooks into the parser-select / expressions nodes.
//
// Legacy ANTLR grammar (TrinoParser.g4):
//
//	sampledRelation : patternRecognition (TABLESAMPLE …)? ;
//	patternRecognition
//	    : aliasedRelation (
//	        MATCH_RECOGNIZE_ LPAREN_
//	          (PARTITION_ BY_ expression (, expression)*)?
//	          (ORDER_ BY_ sortItem (, sortItem)*)?
//	          (MEASURES_ measureDefinition (, measureDefinition)*)?
//	          rowsPerMatch?
//	          (AFTER_ MATCH_ skipTo)?
//	          (INITIAL_ | SEEK_)?
//	          PATTERN_ LPAREN_ rowPattern RPAREN_
//	          (SUBSET_ subsetDefinition (, subsetDefinition)*)?
//	          DEFINE_ variableDefinition (, variableDefinition)*
//	        RPAREN_ (AS_? identifier columnAliases?)?
//	      )? ;
//	measureDefinition : expression AS_ identifier ;
//	rowsPerMatch      : ONE_ ROW_ PER_ MATCH_
//	                  | ALL_ ROWS_ PER_ MATCH_ emptyMatchHandling? ;
//	emptyMatchHandling: SHOW_ EMPTY_ MATCHES_ | OMIT_ EMPTY_ MATCHES_
//	                  | WITH_ UNMATCHED_ ROWS_ ;
//	skipTo            : SKIP_ (TO_ (NEXT_ ROW_ | (FIRST_|LAST_)? identifier)
//	                          | PAST_ LAST_ ROW_) ;
//	subsetDefinition  : identifier EQ_ LPAREN_ identifier (, identifier)* RPAREN_ ;
//	variableDefinition: identifier AS_ expression ;
//	windowFrame
//	    : (MEASURES_ measureDefinition (, measureDefinition)*)? frameExtent
//	      (AFTER_ MATCH_ skipTo)? (INITIAL_ | SEEK_)?
//	      (PATTERN_ LPAREN_ rowPattern RPAREN_)?
//	      (SUBSET_ subsetDefinition (, subsetDefinition)*)?
//	      (DEFINE_ variableDefinition (, variableDefinition)*)? ;
//
// Adjudicated against the live Trino 481 oracle. Oracle-confirmed structure
// (probed; see oracle_match_recognize_test.go):
//
//	M1 (PATTERN and DEFINE are both mandatory in a relation MATCH_RECOGNIZE).
//	   `MATCH_RECOGNIZE (PATTERN (A))`, `MATCH_RECOGNIZE (DEFINE A AS true)`, and
//	   `MATCH_RECOGNIZE ()` are all SYNTAX_ERRORs. PATTERN's body parens are also
//	   required (`PATTERN A …` is rejected).
//	M2 (fixed clause order). The optional clauses appear in exactly the grammar
//	   order — PARTITION BY, ORDER BY, MEASURES, rowsPerMatch, AFTER MATCH,
//	   INITIAL|SEEK, then PATTERN. Reordering (e.g. ORDER BY before PARTITION BY)
//	   is rejected, so each is parsed by a single positional check (no loop).
//	M3 (double alias). The inner aliasedRelation may carry its own alias AND the
//	   MATCH_RECOGNIZE may carry a trailing `AS? identifier columnAliases?`:
//	   `t AS r MATCH_RECOGNIZE (…) AS m` is accepted. The trailing alias is on the
//	   PatternRecognition node; the inner alias stays on the wrapped relation.
//	M4 (MATCH_RECOGNIZE binds before TABLESAMPLE). Because sampledRelation is
//	   `patternRecognition (TABLESAMPLE …)?`, MATCH_RECOGNIZE attaches to the
//	   aliasedRelation and TABLESAMPLE wraps the result: `t MATCH_RECOGNIZE (…)
//	   TABLESAMPLE …` is accepted, `t TABLESAMPLE … MATCH_RECOGNIZE (…)` rejected.
//	   The relation-level hook therefore runs inside parseSampledRelation, between
//	   aliasedRelation and TABLESAMPLE (the D4 boundary noted in relation.go).
//	M5 (INITIAL|SEEK and AFTER MATCH SKIP forms parse). INITIAL / SEEK are
//	   accepted at the grammar level (the planner may answer NOT_SUPPORTED, which
//	   is a semantic outcome, not a syntax rejection); every skipTo form
//	   (TO NEXT ROW, TO [FIRST|LAST] var, TO var, PAST LAST ROW) is accepted.
//	M6 (window-frame pattern). The OVER-clause windowFrame may carry the same
//	   MEASURES / PATTERN / SUBSET / DEFINE / AFTER MATCH / INITIAL|SEEK around a
//	   frameExtent. There, every part is optional EXCEPT frameExtent — a window
//	   frame may have MEASURES with no PATTERN, or PATTERN with no DEFINE
//	   (unlike the mandatory relation form, M1). The expressions node's
//	   parseWindowFrame delegates these additions here (the B2 boundary).

// ---------------------------------------------------------------------------
// patternRecognition node types (parser-package types; not ast.Node — matching
// the Relation / Expr convention, see relation.go / expr.go headers)
// ---------------------------------------------------------------------------

// PatternRecognition is the relation produced by an aliasedRelation followed by
// a `MATCH_RECOGNIZE ( … )` clause (the patternRecognition rule's
// MATCH_RECOGNIZE alternative). It is a Relation, wrapping the underlying
// aliasedRelation in Input. The body fields mirror the grammar; Pattern and
// Definitions are always non-nil/non-empty (M1). Alias / ColumnAliases carry the
// clause's own trailing `AS? identifier columnAliases?` (M3), nil when absent.
//
// When an aliasedRelation is NOT followed by MATCH_RECOGNIZE, no
// PatternRecognition node is produced — the aliasedRelation is returned as-is
// (the patternRecognition rule's empty alternative).
type PatternRecognition struct {
	Input          Relation             // the wrapped aliasedRelation
	PartitionBy    []Expr               // PARTITION BY …, nil when absent
	OrderBy        []SortItem           // ORDER BY …, nil when absent
	Measures       []MeasureDefinition  // MEASURES …, nil when absent
	RowsPerMatch   *RowsPerMatch        // ONE ROW / ALL ROWS PER MATCH, nil when absent
	AfterMatchSkip *SkipTo              // AFTER MATCH SKIP …, nil when absent
	SearchMode     string               // "", "INITIAL", or "SEEK"
	Pattern        RowPattern           // PATTERN ( rowPattern ) — required (M1)
	Subsets        []SubsetDefinition   // SUBSET …, nil when absent
	Definitions    []VariableDefinition // DEFINE … — required, non-empty (M1)
	Alias          *ast.Identifier      // trailing AS? identifier, nil when absent
	ColumnAliases  []*ast.Identifier    // trailing columnAliases, nil when absent
	Loc            ast.Loc
}

func (n *PatternRecognition) Span() ast.Loc { return n.Loc }
func (*PatternRecognition) relationNode()   {}

// MeasureDefinition is one `expression AS identifier` of a MEASURES clause (the
// measureDefinition rule): a named row-pattern measure. The RUNNING/FINAL
// semantics live inside Expr (the expression grammar's processingMode).
type MeasureDefinition struct {
	Expr Expr
	Name *ast.Identifier
	Loc  ast.Loc
}

// RowsPerMatchKind enumerates the rowsPerMatch shapes.
type RowsPerMatchKind int

const (
	// OneRowPerMatch is `ONE ROW PER MATCH`.
	OneRowPerMatch RowsPerMatchKind = iota
	// AllRowsPerMatch is `ALL ROWS PER MATCH [emptyMatchHandling]`.
	AllRowsPerMatch
)

// RowsPerMatch is the `ONE ROW PER MATCH` / `ALL ROWS PER MATCH
// [emptyMatchHandling]` clause (the rowsPerMatch rule). EmptyHandling is "" for
// ONE ROW PER MATCH and for ALL ROWS PER MATCH with no handling; otherwise it is
// "SHOW EMPTY MATCHES", "OMIT EMPTY MATCHES", or "WITH UNMATCHED ROWS".
type RowsPerMatch struct {
	Kind          RowsPerMatchKind
	EmptyHandling string
	Loc           ast.Loc
}

// SkipToKind enumerates the AFTER MATCH SKIP target shapes (the skipTo rule).
type SkipToKind int

const (
	// SkipPastLastRow is `SKIP PAST LAST ROW`.
	SkipPastLastRow SkipToKind = iota
	// SkipToNextRow is `SKIP TO NEXT ROW`.
	SkipToNextRow
	// SkipToFirst is `SKIP TO FIRST identifier`.
	SkipToFirst
	// SkipToLast is `SKIP TO LAST identifier`.
	SkipToLast
	// SkipToVariable is `SKIP TO identifier` (no FIRST/LAST).
	SkipToVariable
)

// SkipTo is an `AFTER MATCH SKIP …` target (the skipTo rule). Kind is the target
// shape; Variable is the pattern-variable name for the SkipToFirst / SkipToLast /
// SkipToVariable forms (nil otherwise).
type SkipTo struct {
	Kind     SkipToKind
	Variable *ast.Identifier // for SKIP TO [FIRST|LAST] variable, nil otherwise
	Loc      ast.Loc
}

// SubsetDefinition is one `identifier = ( identifier (, identifier)* )` of a
// SUBSET clause (the subsetDefinition rule): a union variable Name defined as the
// union of the Of pattern variables.
type SubsetDefinition struct {
	Name *ast.Identifier
	Of   []*ast.Identifier
	Loc  ast.Loc
}

// VariableDefinition is one `identifier AS expression` of a DEFINE clause (the
// variableDefinition rule): a row-pattern variable defined by a boolean
// condition expression.
type VariableDefinition struct {
	Name      *ast.Identifier
	Condition Expr
	Loc       ast.Loc
}

// ---------------------------------------------------------------------------
// Relation-level integration hook (D4 in relation.go's parseSampledRelation)
// ---------------------------------------------------------------------------

// parsePatternRecognitionSuffix is the patternRecognition rule applied to an
// already-parsed aliasedRelation `input`. If the current token is MATCH_RECOGNIZE
// it parses the full `MATCH_RECOGNIZE ( … ) (AS? identifier columnAliases?)?`
// clause and returns a *PatternRecognition wrapping input; otherwise it returns
// input unchanged (the rule's empty alternative). It is called by
// parseSampledRelation between the aliasedRelation and the optional TABLESAMPLE
// (M4), so MATCH_RECOGNIZE binds tighter than TABLESAMPLE.
func (p *Parser) parsePatternRecognitionSuffix(input Relation) (Relation, error) {
	if p.cur.Kind != kwMATCH_RECOGNIZE {
		return input, nil
	}
	p.advance() // consume MATCH_RECOGNIZE
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}

	// The span starts at the wrapped relation (input), which already covers the
	// table/subquery and its inner alias; the MATCH_RECOGNIZE keyword always
	// follows it, so input.Span().Start is the node start.
	pr := &PatternRecognition{
		Input: input,
		Loc:   ast.Loc{Start: input.Span().Start},
	}
	if err := p.parseMatchRecognizeBody(pr); err != nil {
		return nil, err
	}

	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, err
	}
	pr.Loc.End = closeTok.Loc.End

	// Optional trailing `AS? identifier columnAliases?` (M3).
	if alias, cols, end, ok, err := p.parseTrailingMatchRecognizeAlias(); err != nil {
		return nil, err
	} else if ok {
		pr.Alias = alias
		pr.ColumnAliases = cols
		pr.Loc.End = end
	}
	return pr, nil
}

// parseTrailingMatchRecognizeAlias parses the optional `AS? identifier
// columnAliases?` alias that may follow a `MATCH_RECOGNIZE ( … )` clause. It
// mostly delegates to the shared relation alias helper (tryParseRelationAlias),
// but first handles the one case that helper deliberately rejects: a no-AS alias
// literally named MATCH_RECOGNIZE followed by '('. In the aliasedRelation
// position tryParseRelationAlias's atRelationClauseStart guard reads
// `MATCH_RECOGNIZE (` as the start of the MATCH_RECOGNIZE clause (so it is NOT
// swallowed as an alias); but here the clause has already been parsed and a
// second one cannot follow, so `… ) MATCH_RECOGNIZE (x)` is a valid bare alias
// "match_recognize" with a column list — oracle-confirmed (it parses exactly as
// `… ) AS MATCH_RECOGNIZE (x)`). This helper consumes that form directly; every
// other alias (including a bare MATCH_RECOGNIZE not followed by '(') goes through
// the shared helper unchanged.
func (p *Parser) parseTrailingMatchRecognizeAlias() (alias *ast.Identifier, cols []*ast.Identifier, end int, ok bool, err error) {
	if p.cur.Kind == kwMATCH_RECOGNIZE && p.peekNext().Kind == int('(') {
		alias = identFromToken(p.advance()) // consume MATCH_RECOGNIZE as the alias
		end = alias.Loc.End
		c, e, cerr := p.parseColumnAliases()
		if cerr != nil {
			return nil, nil, 0, false, cerr
		}
		return alias, c, e, true, nil
	}
	return p.tryParseRelationAlias()
}

// parseMatchRecognizeBody parses the body of a relation `MATCH_RECOGNIZE ( … )`
// between the opening '(' (already consumed) and the closing ')' (parsed by the
// caller). The clauses appear in a fixed grammar order (M2); PATTERN and DEFINE
// are mandatory (M1). It fills the body fields of pr.
func (p *Parser) parseMatchRecognizeBody(pr *PatternRecognition) error {
	// PARTITION BY expression (, expression)*
	if p.cur.Kind == kwPARTITION {
		exprs, err := p.parsePartitionByExprs()
		if err != nil {
			return err
		}
		pr.PartitionBy = exprs
	}

	// ORDER BY sortItem (, sortItem)*
	if p.cur.Kind == kwORDER {
		items, err := p.parseOrderByItems()
		if err != nil {
			return err
		}
		pr.OrderBy = items
	}

	// MEASURES measureDefinition (, measureDefinition)*
	if p.cur.Kind == kwMEASURES {
		measures, err := p.parseMeasures()
		if err != nil {
			return err
		}
		pr.Measures = measures
	}

	// rowsPerMatch (ONE ROW / ALL ROWS PER MATCH …)
	if p.cur.Kind == kwONE || p.cur.Kind == kwALL {
		rpm, err := p.parseRowsPerMatch()
		if err != nil {
			return err
		}
		pr.RowsPerMatch = rpm
	}

	// AFTER MATCH skipTo
	if p.cur.Kind == kwAFTER {
		skip, err := p.parseAfterMatchSkip()
		if err != nil {
			return err
		}
		pr.AfterMatchSkip = skip
	}

	// INITIAL | SEEK
	if mode, ok := p.match(kwINITIAL, kwSEEK); ok {
		pr.SearchMode = mode.Str
	}

	// PATTERN ( rowPattern ) — mandatory (M1)
	pattern, err := p.parsePatternClause()
	if err != nil {
		return err
	}
	pr.Pattern = pattern

	// SUBSET subsetDefinition (, subsetDefinition)*
	if p.cur.Kind == kwSUBSET {
		subsets, err := p.parseSubsets()
		if err != nil {
			return err
		}
		pr.Subsets = subsets
	}

	// DEFINE variableDefinition (, variableDefinition)* — mandatory (M1)
	defs, err := p.parseDefine()
	if err != nil {
		return err
	}
	pr.Definitions = defs
	return nil
}

// ---------------------------------------------------------------------------
// Window-frame integration hook (B2 in function.go's parseWindowFrame)
// ---------------------------------------------------------------------------

// MatchRecognizeFrame holds the row-pattern additions a windowFrame may carry
// AROUND its frameExtent (the windowFrame rule's MEASURES / AFTER MATCH /
// INITIAL|SEEK / PATTERN / SUBSET / DEFINE parts). Unlike the relation
// MATCH_RECOGNIZE (M1), every part here is optional: a frame may carry MEASURES
// with no PATTERN, or a PATTERN with no DEFINE (M6). Empty slices / nil pointers
// / "" mark absent parts. It hangs off WindowFrame (function.go) via an optional
// pointer set by parseWindowFramePattern.
type MatchRecognizeFrame struct {
	Measures       []MeasureDefinition
	AfterMatchSkip *SkipTo
	SearchMode     string     // "", "INITIAL", or "SEEK"
	Pattern        RowPattern // nil when no PATTERN clause
	Subsets        []SubsetDefinition
	Definitions    []VariableDefinition
}

// atMeasuresClause reports whether a leading MEASURES in a windowSpecification
// begins a pattern-recognition MEASURES clause rather than naming an existing
// window. MEASURES is a non-reserved keyword, so `OVER (measures ORDER BY …)`
// (a window named "measures") and `OVER (MEASURES x AS m …)` (the clause) both
// start with MEASURES — and the two cannot be told apart by the next token
// alone: `OVER (MEASURES ROWS CURRENT ROW)` (window named MEASURES + a ROWS
// frame) and `OVER (MEASURES rows AS m ROWS CURRENT ROW)` (the clause whose
// first measure is the column `rows`) both have a frame-type keyword after
// MEASURES. Oracle-confirmed: Trino accepts both.
//
// The disambiguation therefore mirrors Trino — try to parse a measure
// definition. MEASURES leads the CLAUSE iff, after consuming MEASURES, a single
// `expression` parses and is immediately followed by AS (the `expression AS
// identifier` shape of a measureDefinition). Otherwise MEASURES is an existing
// window name. The probe is fully speculative (a checkpoint is restored before
// returning), so the real parse happens later unaffected; a parse error during
// the probe simply means "not a measure clause". Used by the expressions node's
// parseWindowSpecification to keep MEASURES out of the existing-name slot (B2).
func (p *Parser) atMeasuresClause() bool {
	if p.cur.Kind != kwMEASURES {
		return false
	}
	cp := p.checkpoint()
	defer p.restore(cp)
	p.advance() // consume MEASURES
	if _, err := p.parseExpr(); err != nil {
		return false
	}
	return p.cur.Kind == kwAS
}

// parseLeadingFrameMeasures parses the windowFrame's leading
// `MEASURES measureDefinition (, …)*`, which precedes the frameExtent. It
// returns (nil, false, nil) when the current token is not MEASURES. The
// expressions node's parseWindowSpecification calls this before the frameExtent
// so the frame parser is entered even when MEASURES leads (B2 / M6).
func (p *Parser) parseLeadingFrameMeasures() ([]MeasureDefinition, bool, error) {
	if p.cur.Kind != kwMEASURES {
		return nil, false, nil
	}
	measures, err := p.parseMeasures()
	if err != nil {
		return nil, false, err
	}
	return measures, true, nil
}

// parseTrailingFramePattern parses the windowFrame's row-pattern additions that
// FOLLOW the frameExtent: `(AFTER MATCH skipTo)? (INITIAL|SEEK)?
// (PATTERN ( rowPattern ))? (SUBSET …)? (DEFINE …)?`. leadingMeasures carries the
// MEASURES already parsed before the frameExtent (M6). It returns (nil, nil) when
// none of these clauses are present and there were no leading measures — i.e. the
// frame is an ordinary window frame with no pattern recognition. The expressions
// node's parseWindowFrame calls this after the frameExtent (B2).
func (p *Parser) parseTrailingFramePattern(leadingMeasures []MeasureDefinition) (*MatchRecognizeFrame, error) {
	mrf := &MatchRecognizeFrame{Measures: leadingMeasures}
	saw := len(leadingMeasures) > 0

	if p.cur.Kind == kwAFTER {
		skip, err := p.parseAfterMatchSkip()
		if err != nil {
			return nil, err
		}
		mrf.AfterMatchSkip = skip
		saw = true
	}
	if mode, ok := p.match(kwINITIAL, kwSEEK); ok {
		mrf.SearchMode = mode.Str
		saw = true
	}
	if p.cur.Kind == kwPATTERN {
		pattern, err := p.parsePatternClause()
		if err != nil {
			return nil, err
		}
		mrf.Pattern = pattern
		saw = true
	}
	if p.cur.Kind == kwSUBSET {
		subsets, err := p.parseSubsets()
		if err != nil {
			return nil, err
		}
		mrf.Subsets = subsets
		saw = true
	}
	if p.cur.Kind == kwDEFINE {
		defs, err := p.parseDefine()
		if err != nil {
			return nil, err
		}
		mrf.Definitions = defs
		saw = true
	}
	if !saw {
		return nil, nil
	}
	return mrf, nil
}

// ---------------------------------------------------------------------------
// Shared clause parsers (used by both the relation and window-frame forms)
// ---------------------------------------------------------------------------

// parsePartitionByExprs parses `PARTITION BY expression (, expression)*`
// (PARTITION is current) and returns the non-empty expression list.
func (p *Parser) parsePartitionByExprs() ([]Expr, error) {
	p.advance() // consume PARTITION
	if _, err := p.expect(kwBY); err != nil {
		return nil, err
	}
	first, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	exprs := []Expr{first}
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		next, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		exprs = append(exprs, next)
	}
	return exprs, nil
}

// parseMeasures parses `MEASURES measureDefinition (, measureDefinition)*`
// (MEASURES is current); each measureDefinition is `expression AS identifier`.
func (p *Parser) parseMeasures() ([]MeasureDefinition, error) {
	p.advance() // consume MEASURES
	first, err := p.parseMeasureDefinition()
	if err != nil {
		return nil, err
	}
	measures := []MeasureDefinition{first}
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		next, err := p.parseMeasureDefinition()
		if err != nil {
			return nil, err
		}
		measures = append(measures, next)
	}
	return measures, nil
}

// parseMeasureDefinition parses one `expression AS identifier` (the
// measureDefinition rule). The AS and the name are both mandatory.
func (p *Parser) parseMeasureDefinition() (MeasureDefinition, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return MeasureDefinition{}, err
	}
	if _, err := p.expect(kwAS); err != nil {
		return MeasureDefinition{}, err
	}
	name, err := p.parseIdentifier()
	if err != nil {
		return MeasureDefinition{}, err
	}
	return MeasureDefinition{
		Expr: expr,
		Name: name,
		Loc:  ast.Loc{Start: expr.Span().Start, End: name.Loc.End},
	}, nil
}

// parseRowsPerMatch parses `ONE ROW PER MATCH` or `ALL ROWS PER MATCH
// emptyMatchHandling?` (the leading ONE or ALL is current).
func (p *Parser) parseRowsPerMatch() (*RowsPerMatch, error) {
	lead := p.advance() // consume ONE or ALL
	rpm := &RowsPerMatch{Loc: ast.Loc{Start: lead.Loc.Start}}

	if lead.Kind == kwONE {
		// ONE ROW PER MATCH
		if _, err := p.expect(kwROW); err != nil {
			return nil, err
		}
		rpm.Kind = OneRowPerMatch
	} else {
		// ALL ROWS PER MATCH
		if _, err := p.expect(kwROWS); err != nil {
			return nil, err
		}
		rpm.Kind = AllRowsPerMatch
	}
	if _, err := p.expect(kwPER); err != nil {
		return nil, err
	}
	matchTok, err := p.expect(kwMATCH)
	if err != nil {
		return nil, err
	}
	rpm.Loc.End = matchTok.Loc.End

	// ALL ROWS PER MATCH takes an optional emptyMatchHandling.
	if rpm.Kind == AllRowsPerMatch {
		handling, end, ok, err := p.tryParseEmptyMatchHandling()
		if err != nil {
			return nil, err
		}
		if ok {
			rpm.EmptyHandling = handling
			rpm.Loc.End = end
		}
	}
	return rpm, nil
}

// tryParseEmptyMatchHandling parses the optional `emptyMatchHandling` of an
// ALL ROWS PER MATCH clause: `SHOW EMPTY MATCHES` | `OMIT EMPTY MATCHES` |
// `WITH UNMATCHED ROWS`. It returns ok=false when none of these leads. SHOW /
// OMIT / WITH are the only valid lead-ins here; SHOW and WITH are non-reserved
// but in this position can only begin the handling clause (a following PATTERN /
// MEASURES would have ended the clause already).
func (p *Parser) tryParseEmptyMatchHandling() (handling string, end int, ok bool, err error) {
	switch p.cur.Kind {
	case kwSHOW:
		p.advance() // consume SHOW
		if _, err := p.expect(kwEMPTY); err != nil {
			return "", 0, false, err
		}
		t, err := p.expect(kwMATCHES)
		if err != nil {
			return "", 0, false, err
		}
		return "SHOW EMPTY MATCHES", t.Loc.End, true, nil
	case kwOMIT:
		p.advance() // consume OMIT
		if _, err := p.expect(kwEMPTY); err != nil {
			return "", 0, false, err
		}
		t, err := p.expect(kwMATCHES)
		if err != nil {
			return "", 0, false, err
		}
		return "OMIT EMPTY MATCHES", t.Loc.End, true, nil
	case kwWITH:
		p.advance() // consume WITH
		if _, err := p.expect(kwUNMATCHED); err != nil {
			return "", 0, false, err
		}
		t, err := p.expect(kwROWS)
		if err != nil {
			return "", 0, false, err
		}
		return "WITH UNMATCHED ROWS", t.Loc.End, true, nil
	default:
		return "", 0, false, nil
	}
}

// parseAfterMatchSkip parses `AFTER MATCH skipTo` (AFTER is current). skipTo is
// `SKIP (TO (NEXT ROW | (FIRST|LAST)? identifier) | PAST LAST ROW)`.
func (p *Parser) parseAfterMatchSkip() (*SkipTo, error) {
	afterTok := p.advance() // consume AFTER
	if _, err := p.expect(kwMATCH); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwSKIP); err != nil {
		return nil, err
	}
	skip := &SkipTo{Loc: ast.Loc{Start: afterTok.Loc.Start}}

	switch p.cur.Kind {
	case kwPAST:
		// PAST LAST ROW
		p.advance() // consume PAST
		if _, err := p.expect(kwLAST); err != nil {
			return nil, err
		}
		rowTok, err := p.expect(kwROW)
		if err != nil {
			return nil, err
		}
		skip.Kind = SkipPastLastRow
		skip.Loc.End = rowTok.Loc.End
		return skip, nil
	case kwTO:
		p.advance() // consume TO
		return p.parseSkipToTarget(skip)
	default:
		return nil, p.exprErrorAt("expected TO or PAST after AFTER MATCH SKIP")
	}
}

// parseSkipToTarget parses the part of skipTo after `SKIP TO`: `NEXT ROW |
// (FIRST|LAST)? identifier`. skip already has its start offset; TO is consumed.
//
// NEXT, FIRST, and LAST are all NON-RESERVED keywords, so each may itself be a
// pattern-variable name (oracle-confirmed: `SKIP TO NEXT`, `SKIP TO FIRST`,
// `SKIP TO LAST`, `SKIP TO FIRST PATTERN (A)` all parse with NEXT / FIRST / LAST
// as the variable). The disambiguation mirrors Trino's adaptive prediction:
//
//   - NEXT is the `NEXT ROW` form ONLY when immediately followed by ROW;
//     otherwise NEXT is an ordinary variable identifier.
//   - FIRST / LAST is the `(FIRST|LAST) identifier` keyword ONLY when an
//     identifier follows it AND that identifier is in turn followed by a token
//     that legitimately ends the skipTo — INITIAL, SEEK, or PATTERN (the only
//     clauses that may follow `AFTER MATCH skipTo`). This is the discriminator
//     because PATTERN is mandatory and always trails the skipTo. It is checked
//     with a speculative checkpoint (peek beyond two tokens): e.g.
//     `FIRST PATTERN (A)` has FIRST as the variable (the `(` after PATTERN is not
//     a skipTo successor), whereas `FIRST zone PATTERN (A)` has FIRST as the
//     keyword with variable zone. Otherwise FIRST / LAST is itself the variable.
//
// Anything else is a bare `TO identifier` variable target.
func (p *Parser) parseSkipToTarget(skip *SkipTo) (*SkipTo, error) {
	if p.cur.Kind == kwNEXT && p.peekNext().Kind == kwROW {
		// TO NEXT ROW
		p.advance()           // consume NEXT
		rowTok := p.advance() // consume ROW
		skip.Kind = SkipToNextRow
		skip.Loc.End = rowTok.Loc.End
		return skip, nil
	}

	if p.cur.Kind == kwFIRST || p.cur.Kind == kwLAST {
		if kind, name, ok := p.tryParseSkipFirstLastKeyword(); ok {
			skip.Kind = kind
			skip.Variable = name
			skip.Loc.End = name.Loc.End
			return skip, nil
		}
	}

	// TO identifier (no FIRST/LAST keyword, or NEXT/FIRST/LAST used as the
	// variable name itself).
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	skip.Kind = SkipToVariable
	skip.Variable = name
	skip.Loc.End = name.Loc.End
	return skip, nil
}

// tryParseSkipFirstLastKeyword speculatively parses the `(FIRST|LAST) identifier`
// keyword form of a skipTo target. It returns the SkipToFirst/SkipToLast kind,
// the variable identifier, and ok=true only when the leading FIRST/LAST is
// genuinely the keyword — i.e. an identifier follows it AND that identifier is
// itself followed by a skipTo successor (INITIAL / SEEK / PATTERN), which is how
// Trino disambiguates FIRST/LAST as a keyword from FIRST/LAST used as the variable
// name. On a non-match it restores the parser and returns ok=false (the caller
// then reads FIRST/LAST as the variable).
func (p *Parser) tryParseSkipFirstLastKeyword() (SkipToKind, *ast.Identifier, bool) {
	cp := p.checkpoint()
	kind := SkipToFirst
	if p.cur.Kind == kwLAST {
		kind = SkipToLast
	}
	p.advance() // consume FIRST / LAST
	if !isIdentifierStart(p.cur.Kind) {
		p.restore(cp)
		return 0, nil, false
	}
	name := identFromToken(p.advance())
	if !isSkipToSuccessor(p.cur.Kind) {
		p.restore(cp)
		return 0, nil, false
	}
	return kind, name, true
}

// isSkipToSuccessor reports whether kind is a token that may immediately follow a
// MATCH_RECOGNIZE skipTo target. This is the complete set across both contexts:
//
//   - relation MATCH_RECOGNIZE: only `(INITIAL|SEEK)? PATTERN` may follow, so
//     INITIAL / SEEK / PATTERN.
//   - windowFrame: the skipTo is followed by the OPTIONAL `(INITIAL|SEEK)?
//     (PATTERN …)? (SUBSET …)? (DEFINE …)?`, any of which may be absent — so the
//     successor may also be SUBSET, DEFINE, or the frame-closing ')'.
//
// Using the union is safe: a FIRST/LAST mis-promoted to the keyword form in a
// context where the trailing clause is then illegal (e.g. `SKIP TO FIRST A
// SUBSET …` in a relation MATCH_RECOGNIZE, where PATTERN is mandatory before
// SUBSET) still produces the same overall reject as Trino, because the body
// parser rejects the missing PATTERN. It is used to decide whether a FIRST/LAST
// is the skipTo keyword (followed by `variable <successor>`) or the variable.
func isSkipToSuccessor(kind TokenKind) bool {
	switch kind {
	case kwINITIAL, kwSEEK, kwPATTERN, kwSUBSET, kwDEFINE, int(')'):
		return true
	default:
		return false
	}
}

// parsePatternClause parses `PATTERN ( rowPattern )` (PATTERN is current). The
// body parens are mandatory (M1); the inner rowPattern is parsed by
// row_pattern.go.
func (p *Parser) parsePatternClause() (RowPattern, error) {
	if _, err := p.expect(kwPATTERN); err != nil {
		return nil, err
	}
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	pattern, err := p.parseRowPattern()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	return pattern, nil
}

// parseSubsets parses `SUBSET subsetDefinition (, subsetDefinition)*` (SUBSET is
// current); each subsetDefinition is `identifier = ( identifier (, identifier)* )`.
func (p *Parser) parseSubsets() ([]SubsetDefinition, error) {
	p.advance() // consume SUBSET
	first, err := p.parseSubsetDefinition()
	if err != nil {
		return nil, err
	}
	subsets := []SubsetDefinition{first}
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		next, err := p.parseSubsetDefinition()
		if err != nil {
			return nil, err
		}
		subsets = append(subsets, next)
	}
	return subsets, nil
}

// parseSubsetDefinition parses one `identifier = ( identifier (, identifier)* )`
// (the subsetDefinition rule): a union variable and its non-empty member list.
func (p *Parser) parseSubsetDefinition() (SubsetDefinition, error) {
	name, err := p.parseIdentifier()
	if err != nil {
		return SubsetDefinition{}, err
	}
	if _, err := p.expect(int('=')); err != nil {
		return SubsetDefinition{}, err
	}
	if _, err := p.expect(int('(')); err != nil {
		return SubsetDefinition{}, err
	}
	firstOf, err := p.parseIdentifier()
	if err != nil {
		return SubsetDefinition{}, err
	}
	of := []*ast.Identifier{firstOf}
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		nextOf, err := p.parseIdentifier()
		if err != nil {
			return SubsetDefinition{}, err
		}
		of = append(of, nextOf)
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return SubsetDefinition{}, err
	}
	return SubsetDefinition{
		Name: name,
		Of:   of,
		Loc:  ast.Loc{Start: name.Loc.Start, End: closeTok.Loc.End},
	}, nil
}

// parseDefine parses `DEFINE variableDefinition (, variableDefinition)*` (DEFINE
// is current); each variableDefinition is `identifier AS expression`. At least
// one definition is required (the rule has no `?`).
func (p *Parser) parseDefine() ([]VariableDefinition, error) {
	if _, err := p.expect(kwDEFINE); err != nil {
		return nil, err
	}
	first, err := p.parseVariableDefinition()
	if err != nil {
		return nil, err
	}
	defs := []VariableDefinition{first}
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		next, err := p.parseVariableDefinition()
		if err != nil {
			return nil, err
		}
		defs = append(defs, next)
	}
	return defs, nil
}

// parseVariableDefinition parses one `identifier AS expression` (the
// variableDefinition rule): a pattern variable and its boolean condition.
func (p *Parser) parseVariableDefinition() (VariableDefinition, error) {
	name, err := p.parseIdentifier()
	if err != nil {
		return VariableDefinition{}, err
	}
	if _, err := p.expect(kwAS); err != nil {
		return VariableDefinition{}, err
	}
	cond, err := p.parseExpr()
	if err != nil {
		return VariableDefinition{}, err
	}
	return VariableDefinition{
		Name:      name,
		Condition: cond,
		Loc:       ast.Loc{Start: name.Loc.Start, End: cond.Span().End},
	}, nil
}
