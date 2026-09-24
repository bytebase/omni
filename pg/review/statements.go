package review

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/bytebase/omni/pg"
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
	"github.com/bytebase/omni/review"
)

// statement is one parsed statement with its index in Result.Statements.
type statement struct {
	index int
	node  ast.Node
	// loc is the statement's absolute byte range in the input, from its
	// first token to its last, without the trailing semicolon.
	loc ast.Loc
}

// statementRanges lists every non-empty statement's byte range. It comes
// from the lexical splitter, the same pg.Split Bytebase splits with, so
// it does not depend on the text parsing.
func statementRanges(sql string) []review.Range {
	var ranges []review.Range
	for _, seg := range pg.Split(sql) {
		if seg.Empty() {
			continue
		}
		ranges = append(ranges, review.Range{Start: seg.ByteStart, End: seg.ByteEnd})
	}
	return ranges
}

// parse parses the whole text through pg.Parse, the entry Bytebase
// parses with. On success it returns the statements mapped onto ranges.
// On failure it returns the one Syntax finding of the review: the parser
// stops at the first error, so there is never more than one, and the
// review ends there. A non-empty range that produced no statement is a
// failure too: the parser drops an unterminated string, identifier, or
// comment that begins a statement instead of reporting it.
func parse(sql string, ranges []review.Range) ([]statement, *review.Finding) {
	parsed, err := pg.Parse(sql)
	if err != nil {
		return nil, syntaxFinding(err, ranges)
	}
	stmts := make([]statement, 0, len(parsed))
	covered := make([]bool, len(ranges))
	for _, p := range parsed {
		if p.Empty() {
			continue
		}
		// A node the parser gave no location has NoLoc, or the zero Loc
		// when its kind carries one the parser did not fill in; either
		// falls outside the statement's own bytes and is replaced by
		// them.
		loc := ast.NodeLoc(p.AST)
		if loc.Start < p.ByteStart || loc.End > p.ByteEnd || loc.End <= loc.Start {
			loc = contentLoc(sql, p.ByteStart, p.ByteEnd)
		}
		index := statementAt(ranges, loc.Start)
		if index >= 0 {
			covered[index] = true
		}
		stmts = append(stmts, statement{index: index, node: p.AST, loc: loc})
	}
	for i, ok := range covered {
		if !ok {
			return nil, unparsedFinding(sql, i, ranges[i])
		}
	}
	return stmts, nil
}

// unparsedFinding is the Syntax finding for a non-empty range the parser
// produced no statement for. The parser gives no message for it, so the
// finding names the first token of the text the way a syntax error does.
func unparsedFinding(sql string, index int, r review.Range) *review.Finding {
	loc := contentLoc(sql, r.Start, r.End)
	text := sql[loc.Start:loc.End]
	if i := strings.IndexAny(text, " \t\r\n"); i >= 0 {
		text = text[:i]
	}
	if len(text) > 40 {
		text = text[:40]
	}
	return &review.Finding{
		Rule:      review.Syntax,
		Statement: index,
		Range:     review.Range{Start: loc.Start, End: loc.Start},
		Message:   escapeLines(fmt.Sprintf("syntax error at or near %q", text)),
	}
}

// contentLoc is the range of a statement's text without its surrounding
// whitespace and trailing semicolon, for a node that carries no location.
func contentLoc(sql string, start, end int) ast.Loc {
	for start < end && isSpace(sql[start]) {
		start++
	}
	for end > start && (isSpace(sql[end-1]) || sql[end-1] == ';') {
		end--
	}
	return ast.Loc{Start: start, End: end}
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

// syntaxFinding turns a parse failure into the Syntax finding. A parser
// error carries the byte offset of the offending token and the message
// PostgreSQL would print; the finding's range starts and ends at that
// offset, which the caller anchors to its line. Any other error, which
// the parser does not produce today, addresses the whole change.
func syntaxFinding(err error, ranges []review.Range) *review.Finding {
	f := &review.Finding{Rule: review.Syntax, Statement: -1, Message: escapeLines(err.Error())}
	var perr *parser.ParseError
	if errors.As(err, &perr) {
		// The message quotes the offending token, which may hold a line
		// break inside a quoted identifier.
		f.Message = escapeLines(perr.Message)
		if perr.Position >= 0 {
			f.Statement = statementAt(ranges, perr.Position)
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

func rangeOf(loc ast.Loc) review.Range {
	if loc.Start < 0 || loc.End < loc.Start {
		return review.Range{}
	}
	return review.Range{Start: loc.Start, End: loc.End}
}
