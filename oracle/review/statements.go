package review

import (
	"errors"
	"strings"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/oracle/parser"
	"github.com/bytebase/omni/review"
)

// statement is one parsed statement with its index in Result.Statements.
type statement struct {
	index int
	node  ast.StmtNode
	// loc is the statement's absolute byte range in the input, from its
	// first token to its last: a SQL statement without its semicolon, a
	// PL/SQL unit through the semicolon after its END.
	loc ast.Loc
}

// parseReview splits the text with parser.Split, the splitter Bytebase
// splits Oracle with, and parses each segment in place, so every location
// and error position is an offset into the text. A SQL*Plus command, such
// as SET DEFINE OFF or the / that ends a PL/SQL unit, is not a statement:
// Bytebase drops it. The first failure ends the review with its one Syntax
// finding: nothing after a syntax error can be trusted, and one finding
// per sheet is the contract. The segments after it are still listed.
//
// bytebase: backend/plugin/parser/plsql/split.go — SplitSQL
// bytebase: backend/plugin/parser/plsql/plsql.go — parsePLSQLStatements
func parseReview(sql string) ([]review.Range, []statement, *review.Finding) {
	var ranges []review.Range
	var stmts []statement
	var syntax *review.Finding
	for _, seg := range parser.Split(sql) {
		if seg.Empty() || seg.Kind == parser.SegmentSQLPlusCommand {
			continue
		}
		if syntax != nil {
			ranges = append(ranges, contentRange(sql, seg.ByteStart, seg.ByteEnd))
			continue
		}
		list, err := parser.ParseRange(sql, seg.ByteStart, seg.ByteEnd)
		if err != nil {
			syntax = syntaxFinding(err, len(ranges))
			ranges = append(ranges, contentRange(sql, seg.ByteStart, seg.ByteEnd))
			continue
		}
		parsed := 0
		if list != nil {
			for _, item := range list.Items {
				raw, ok := item.(*ast.RawStmt)
				if !ok || raw.Stmt == nil {
					continue
				}
				loc := raw.Loc
				if loc.Start < 0 || loc.End <= loc.Start {
					content := contentRange(sql, seg.ByteStart, seg.ByteEnd)
					loc = ast.Loc{Start: content.Start, End: content.End}
				}
				stmts = append(stmts, statement{index: len(ranges), node: raw.Stmt, loc: loc})
				ranges = append(ranges, review.Range{Start: loc.Start, End: loc.End})
				parsed++
			}
		}
		// A guard against text the parser drops: a segment with content is
		// a statement or an error.
		if parsed == 0 {
			content := contentRange(sql, seg.ByteStart, seg.ByteEnd)
			syntax = &review.Finding{
				Rule:      review.Syntax,
				Statement: len(ranges),
				Range:     review.Range{Start: content.Start, End: content.Start},
				Message:   "syntax error at or near " + quoteToken(sql, content),
			}
			ranges = append(ranges, content)
		}
	}
	if syntax != nil {
		return ranges, nil, syntax
	}
	return ranges, stmts, nil
}

// contentRange is the range of a segment's text from its first token to
// its last, without the comments around it and a final semicolon. Where
// the lexer cannot read the text, only whitespace is trimmed.
func contentRange(sql string, start, end int) review.Range {
	lex := parser.NewLexerRange(sql, start, end)
	first, last, beforeLast, lastType := -1, -1, -1, 0
	for {
		tok := lex.NextToken()
		if tok.Type == 0 || lex.Err != nil {
			break
		}
		if first < 0 {
			first = tok.Loc
		}
		beforeLast, last, lastType = last, tok.End, tok.Type
	}
	if lastType == ';' && beforeLast > first {
		last = beforeLast
	}
	if first < 0 || lex.Err != nil {
		s, e := start, end
		for s < e && isSpace(sql[s]) {
			s++
		}
		for e > s && (isSpace(sql[e-1]) || sql[e-1] == ';') {
			e--
		}
		return review.Range{Start: s, End: e}
	}
	return review.Range{Start: first, End: last}
}

// quoteToken quotes the first word of a range, for a message about text
// the parser produced nothing for.
func quoteToken(sql string, r review.Range) string {
	text := sql[r.Start:r.End]
	if i := strings.IndexAny(text, " \t\r\n\f"); i >= 0 {
		text = text[:i]
	}
	if len(text) > 40 {
		text = text[:40]
	}
	return `"` + escapeLines(text) + `"`
}

// isSpace is the lexer's whitespace set.
func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f'
}

// syntaxFinding turns a parse failure of the statement at the given index
// into the Syntax finding. A parser error carries the byte offset of the
// offending token and its message; the finding's range starts and ends at
// that offset, which the caller anchors to its line. An error without an
// offset anchors the whole statement.
func syntaxFinding(err error, index int) *review.Finding {
	f := &review.Finding{Rule: review.Syntax, Statement: index, Message: escapeLines(err.Error())}
	var perr *parser.ParseError
	if errors.As(err, &perr) {
		f.Message = escapeLines(perr.Message)
		if perr.Position >= 0 {
			f.Range = review.Range{Start: perr.Position, End: perr.Position}
		}
	}
	return f
}

// rangeOf converts a node's location to a finding range, the zero range
// for a location the parser did not set.
func rangeOf(loc ast.Loc) review.Range {
	if loc.Start < 0 || loc.End <= loc.Start {
		return review.Range{}
	}
	return review.Range{Start: loc.Start, End: loc.End}
}

// textOf returns the source text of a location, "" for one not set.
func textOf(sql string, loc ast.Loc) string {
	if loc.Start < 0 || loc.End > len(sql) || loc.End <= loc.Start {
		return ""
	}
	return sql[loc.Start:loc.End]
}

// escapeLines keeps a message on one line: a line break inside quoted text
// would otherwise split the comment that carries it.
func escapeLines(s string) string {
	return strings.NewReplacer("\r", `\r`, "\n", `\n`).Replace(s)
}

// tokens returns the tokens of a location of the input.
func tokens(sql string, loc ast.Loc) []parser.Token {
	if loc.Start < 0 || loc.End > len(sql) || loc.End <= loc.Start {
		return nil
	}
	lex := parser.NewLexerRange(sql, loc.Start, loc.End)
	var out []parser.Token
	for {
		tok := lex.NextToken()
		if tok.Type == 0 || lex.Err != nil {
			return out
		}
		out = append(out, tok)
	}
}

// keywordCascade is the token type of the CASCADE keyword. The parser
// reads CASCADE in a DROP clause without recording it, so the rules find
// it among the clause's tokens; a quoted "CASCADE" is an identifier, of
// another type.
var keywordCascade = parser.Tokenize("CASCADE")[0].Type

// cascades reports whether a clause's text holds the CASCADE keyword.
func cascades(sql string, loc ast.Loc) bool {
	for _, tok := range tokens(sql, loc) {
		if tok.Type == keywordCascade {
			return true
		}
	}
	return false
}
