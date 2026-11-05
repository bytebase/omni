package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// TiDB's RECOMMEND INDEX, SET CONFIG, ALTER RANGE and CANCEL statements.
//
// None of this is MySQL's: mysql:8.0, mariadb:11.8.8 and OceanBase 4.3.5 answer a
// syntax error for every form (tidb_admin_refused_test.go in each of those
// parsers holds them to it). The grammar below is v8.5.5's, read from parser.y and
// confirmed form by form against a live server in
// tidb_admin_forms_oracle_test.go.
//
// None of the words here is registered in the lexer for this: RECOMMEND, APPLY,
// CANCEL and BR are read with isBareWord, which keeps each an ordinary name
// everywhere else, and a backquoted spelling is not the word the grammar means.
// RUN, CONFIG, DISTRIBUTION and JOB are older unreserved keywords, read as tokens.

// parseRecommendIndexStmt parses RECOMMEND INDEX, TiDB's index advisor.
//
//	RECOMMEND INDEX RUN [FOR 'sql'] [WITH option, ...]
//	RECOMMEND INDEX SHOW
//	RECOMMEND INDEX {APPLY | IGNORE} id
//	RECOMMEND INDEX SET option, ...
//
// p.cur is the word RECOMMEND.
//
// Ref: TiDB v8.5.5 parser.y:14346-14426.
func (p *Parser) parseRecommendIndexStmt() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume RECOMMEND
	if _, err := p.expect(kwINDEX); err != nil {
		return nil, err
	}

	stmt := &nodes.RecommendIndexStmt{Loc: nodes.Loc{Start: start}}
	var err error
	switch {
	case p.cur.Type == kwRUN:
		p.advance()
		stmt.Action = nodes.RecommendIndexRun
		if _, ok := p.match(kwFOR); ok {
			// One stringLit: `FOR 'a' 'b'` is a 1064, where the same two
			// literals concatenate as an option's value.
			stmt.SQL, _, err = p.parseStringLit()
			if err != nil {
				return nil, err
			}
			stmt.HasSQL = true
		}
		if _, ok := p.match(kwWITH); ok {
			stmt.Options, err = p.parseRecommendIndexOptions()
			if err != nil {
				return nil, err
			}
		}
	case p.cur.Type == kwSHOW:
		p.advance()
		stmt.Action = nodes.RecommendIndexShow
	case p.isBareWord("APPLY"), p.isBareWord("IGNORE"):
		stmt.Action = nodes.RecommendIndexApply
		if p.isBareWord("IGNORE") {
			stmt.Action = nodes.RecommendIndexIgnore
		}
		p.advance()
		// The grammar says NUM and the action asserts an int64, so a number
		// past that range is refused while parsing with 1105, and one past
		// uint64 is a 1064. Int64Num's rule says the same.
		stmt.ID, err = p.parseAdminInt64()
		if err != nil {
			return nil, err
		}
	case p.cur.Type == kwSET:
		p.advance()
		stmt.Action = nodes.RecommendIndexSet
		stmt.Options, err = p.parseRecommendIndexOptions()
		if err != nil {
			return nil, err
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseRecommendIndexOptions parses RecommendIndexOptionList, one or more of
// `name = literal`. The `=` is the eq token alone: `a := 1` is a 1064, where SET's
// own assignments take both.
//
// Ref: TiDB v8.5.5 parser.y:14411-14426.
func (p *Parser) parseRecommendIndexOptions() ([]*nodes.RecommendIndexOption, error) {
	var opts []*nodes.RecommendIndexOption
	for {
		start := p.pos()
		name, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		val, err := p.parseGrammarLiteral()
		if err != nil {
			return nil, err
		}
		opts = append(opts, &nodes.RecommendIndexOption{
			Loc:   nodes.Loc{Start: start, End: p.prev.End},
			Name:  name,
			Value: val,
		})
		if _, ok := p.match(','); !ok {
			return opts, nil
		}
	}
}

// parseGrammarLiteral parses the grammar's Literal: a constant and nothing built
// from one, so `-1`, `(1)`, `1 + 1` and `now()` are 1064s.
//
// Adjacent string literals concatenate, and only plain ones do: `_utf8mb4 'x' 'y'`,
// `N'x' 'y'`, `'x' N'y'` and `x'10' 'y'` are 1064s. An introducer names a character
// set the lexer knows, and `_foo 'x'` is a 1064 for reading _foo as a name. One it
// knows and has no support for, `_utf16 'x'`, answers 1115 from the action. That is
// a name looked up while parsing, which this parser does not refuse.
//
// Ref: TiDB v8.5.5 parser.y:7573 (Literal), pkg/parser/lexer.go (startWithNn).
func (p *Parser) parseGrammarLiteral() (nodes.ExprNode, error) {
	switch p.cur.Type {
	case kwTRUE, kwFALSE, kwNULL, tokICONST, tokFCONST, tokXCONST, tokBCONST:
		return p.parsePrimaryExpr()
	case tokSCONST:
		tok := p.advance()
		lit := &nodes.StringLit{Loc: nodes.Loc{Start: tok.Loc}, Value: tok.Str}
		for p.cur.Type == tokSCONST {
			lit.Value += p.advance().Str
		}
		lit.Loc.End = p.prev.End
		return lit, nil
	}

	if p.cur.Type != tokIDENT || p.backtickQuotedAt(p.cur.Loc) {
		return nil, p.syntaxErrorAtCur()
	}
	word := p.cur.Str
	next := p.peekNext()
	switch {
	case strings.EqualFold(word, "n") && next.Type == tokSCONST && next.Loc == p.cur.End:
		start := p.pos()
		p.advance() // consume N
		tok := p.advance()
		return &nodes.StringLit{Loc: nodes.Loc{Start: start, End: tok.End}, Value: tok.Str, Charset: "utf8"}, nil
	case isCharsetIntroducerWord(word) && next.Type == tokSCONST:
		start := p.pos()
		p.advance() // consume the introducer
		tok := p.advance()
		return &nodes.StringLit{Loc: nodes.Loc{Start: start, End: tok.End}, Value: tok.Str, Charset: word}, nil
	case isCharsetIntroducerWord(word) && (next.Type == tokXCONST || next.Type == tokBCONST):
		p.advance() // consume the introducer
		return p.parsePrimaryExpr()
	}
	return nil, p.syntaxErrorAtCur()
}

// parseSetConfigStmt parses SET CONFIG, after SET.
//
//	SET CONFIG {type | 'instance'} item {= | :=} value
//
// p.cur is kwCONFIG. CONFIG is unreserved and `SET config = 1`, `SET config := 1`
// and `SET config.a = 1` assign a variable called config, which the caller has
// sent elsewhere. There is one item and no list: `SET CONFIG tikv a = 1, b = 2` is a
// 1064.
//
// Ref: TiDB v8.5.5 parser.y:10781-10788.
func (p *Parser) parseSetConfigStmt(start int) (nodes.Node, error) {
	p.advance() // consume CONFIG

	stmt := &nodes.SetConfigStmt{Loc: nodes.Loc{Start: start}}
	if p.cur.Type == tokSCONST {
		stmt.Instance = p.advance().Str
		stmt.HasInstance = true
	} else {
		name, _, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		stmt.Type = strings.ToLower(name)
	}

	name, err := p.parseConfigItemName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	asg, err := p.parseSetAssignmentValue(start, nil)
	if err != nil {
		return nil, err
	}
	stmt.Value = asg.Value
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseConfigItemName parses ConfigItemName, names joined by `.` and `-`.
//
// The same name reads differently on each side of the two. After a dot the lexer
// makes any word a name, so `a.select` and `a.1` parse, while `a-select` and `a-1`
// are 1064s: there the word is an ordinary token, and `a-1b` parses only because 1b
// lexes as a name. The lexer's identDot and dotNamesWord are what decide all of that.
//
// Ref: TiDB v8.5.5 parser.y:10931-10940.
func (p *Parser) parseConfigItemName() (string, error) {
	var name strings.Builder
	for {
		part, _, err := p.parseIdent()
		if err != nil {
			return "", err
		}
		name.WriteString(part)
		if p.cur.Type != '.' && p.cur.Type != '-' {
			return name.String(), nil
		}
		name.WriteByte(byte(p.cur.Type))
		p.advance()
	}
}

// parseAlterRangeStmt parses ALTER RANGE, after ALTER.
//
//	ALTER RANGE name PLACEMENT POLICY [=] {name | 'name' | DEFAULT}
//	ALTER RANGE name PLACEMENT POLICY SET DEFAULT
//
// p.cur is kwRANGE. There is one option and it is this one: the direct placement
// options ALTER TABLE takes, and a second PLACEMENT POLICY, are 1064s.
//
// Ref: TiDB v8.5.5 parser.y:13764-13769, 2068-2084.
func (p *Parser) parseAlterRangeStmt(start int) (nodes.Node, error) {
	p.advance() // consume RANGE

	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	stmt := &nodes.AlterRangeStmt{Loc: nodes.Loc{Start: start}, Name: name}

	if _, err := p.expect(kwPLACEMENT); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwPOLICY); err != nil {
		return nil, err
	}

	if _, ok := p.match(kwSET); ok {
		if _, err := p.expect(kwDEFAULT); err != nil {
			return nil, err
		}
		stmt.Policy = "DEFAULT"
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}
	p.eqOpt()
	switch p.cur.Type {
	case kwDEFAULT:
		p.advance()
		stmt.Policy = "DEFAULT"
	case tokSCONST:
		stmt.Policy = p.advance().Str
	default:
		stmt.Policy, _, err = p.parseIdent()
		if err != nil {
			return nil, err
		}
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCancelStmt parses CANCEL, which cancels a job by its id.
//
//	CANCEL BR JOB id
//	CANCEL DISTRIBUTION JOB id
//
// p.cur is the word CANCEL. The id is Int64Num: an integer literal that fits an
// int64, so -1, 0x10 and 9223372036854775808 are all 1064s.
//
// Ref: TiDB v8.5.5 parser.y:5694-5700 (BR), 3198-3204 (DISTRIBUTION).
func (p *Parser) parseCancelStmt() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume CANCEL

	switch {
	case p.isBareWord("BR"):
		p.advance()
		id, err := p.parseJobID()
		if err != nil {
			return nil, err
		}
		return &nodes.BRIEStmt{
			Loc:   nodes.Loc{Start: start, End: p.prev.End},
			Kind:  nodes.BRIECancelJob,
			JobID: id,
		}, nil
	case p.cur.Type == kwDISTRIBUTION:
		p.advance()
		id, err := p.parseJobID()
		if err != nil {
			return nil, err
		}
		return &nodes.CancelDistributionJobStmt{
			Loc:   nodes.Loc{Start: start, End: p.prev.End},
			JobID: id,
		}, nil
	}
	return nil, p.syntaxErrorAtCur()
}

// parseJobID parses JOB id: the word and the Int64Num after it.
func (p *Parser) parseJobID() (int64, error) {
	if _, err := p.expect(kwJOB); err != nil {
		return 0, err
	}
	return p.parseAdminInt64()
}
