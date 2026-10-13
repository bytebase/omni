package review

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/snowflake/parser"
)

// statement is one parsed statement with its index in Result.Statements.
type statement struct {
	index int
	node  ast.Node
	// loc is the statement's absolute byte range in the input, from its
	// first token to its last, without the trailing semicolon.
	loc ast.Loc
}

// parseReview splits the text with parser.Split, the splitter Bytebase
// splits Snowflake with, and parses it with parser.Parse, which parses each
// of those segments in place: the lexer starts at the segment's offset, so
// every location and error position is an offset into the text. The
// statements pair with the segments in order. A text that does not parse
// gets its one Syntax finding, at the first error: nothing after a syntax
// error can be trusted, and one finding per sheet is the contract.
//
// bytebase: backend/plugin/parser/snowflake/split.go — SplitSQL
func parseReview(sql string) ([]review.Range, []statement, *review.Finding, error) {
	var ranges []review.Range
	for _, seg := range parser.Split(sql) {
		if loc, ok := contentLoc(seg); ok {
			ranges = append(ranges, review.Range{Start: loc.Start, End: loc.End})
		}
	}
	file, err := parser.Parse(sql)
	if err != nil {
		return ranges, nil, syntaxFinding(err, ranges), nil
	}
	if file == nil || len(file.Stmts) != len(ranges) {
		n := 0
		if file != nil {
			n = len(file.Stmts)
		}
		return nil, nil, nil, fmt.Errorf("snowflake/review: the parser returned %d statements for %d segments", n, len(ranges))
	}
	stmts := make([]statement, len(ranges))
	for i, node := range file.Stmts {
		stmts[i] = statement{index: i, node: node, loc: ast.Loc{Start: ranges[i].Start, End: ranges[i].End}}
	}
	return ranges, stmts, nil, nil
}

// contentLoc is a segment's range from its first token to its last, without
// the whitespace and comments around them; ok is false for a segment with
// no token. The parser's own lexer decides what a token is, so a comment is
// trivia exactly when the parser skips it.
func contentLoc(seg parser.Segment) (loc ast.Loc, ok bool) {
	lx := parser.NewLexerWithOffset(seg.Text, seg.ByteStart)
	for {
		tok := lx.NextToken()
		if parser.TokenName(tok.Type) == "EOF" {
			return loc, ok
		}
		if !ok {
			loc.Start, ok = tok.Loc.Start, true
		}
		loc.End = tok.Loc.End
	}
}

// syntaxFinding turns a parse failure into the Syntax finding. A parser
// error carries the byte offset of the offending token and its message;
// the finding's range starts and ends at that offset, which the caller
// anchors to its line. Any other error addresses the whole change.
func syntaxFinding(err error, ranges []review.Range) *review.Finding {
	f := &review.Finding{Rule: review.Syntax, Statement: -1, Message: escapeLines(err.Error())}
	var perr *parser.ParseError
	if errors.As(err, &perr) {
		// The message quotes the offending token, which may hold a line
		// break inside a quoted identifier.
		f.Message = escapeLines(perr.Message)
		if perr.Position >= 0 {
			if len(ranges) > 0 {
				f.Statement = max(0, statementAt(ranges, perr.Position))
			}
			f.Range = review.Range{Start: perr.Position, End: perr.Position}
		}
	}
	return f
}

// statementAt returns the index of the range containing the byte offset,
// or the last range starting before it, or -1 when none does.
func statementAt(ranges []review.Range, offset int) int {
	i := sort.Search(len(ranges), func(i int) bool { return ranges[i].Start > offset })
	return i - 1
}

// rangeOf converts a node's location to a finding range, the zero range
// for a location the parser did not set.
func rangeOf(loc ast.Loc) review.Range {
	if loc.Start < 0 || loc.End <= loc.Start {
		return review.Range{}
	}
	return review.Range{Start: loc.Start, End: loc.End}
}

// anchor is the range a finding on a statement node anchors: the whole
// statement when the node is the statement, else the node's location. The
// parser starts an ALTER statement's location, and DROP DATABASE's and
// DROP SCHEMA's, at the object keyword after ALTER or DROP; the statement's
// own range starts at its first token.
func anchor(s *statement, n ast.Node, loc ast.Loc) review.Range {
	if n == s.node {
		return rangeOf(s.loc)
	}
	return rangeOf(loc)
}

// oneLine writes source text on one line: whitespace between tokens
// becomes one space, and text inside quotes is kept as written, with a
// line break escaped.
func oneLine(text string) string {
	var b strings.Builder
	var quote byte
	space := false
	for i := 0; i < len(text); i++ {
		c := text[i]
		switch {
		case quote != 0:
			if c == quote {
				quote = 0
			}
		case c == '\'' || c == '"':
			quote = c
		case c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f':
			space = true
			continue
		}
		if space && b.Len() > 0 {
			b.WriteByte(' ')
		}
		space = false
		b.WriteByte(c)
	}
	return escapeLines(b.String())
}

// escapeLines keeps a message on one line: a line break inside quoted text
// would otherwise split the comment that carries it.
func escapeLines(s string) string {
	return strings.NewReplacer("\r", `\r`, "\n", `\n`).Replace(s)
}
