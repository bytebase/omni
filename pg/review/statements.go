package review

import (
	"errors"
	"sort"

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
// review ends there.
func parse(sql string, ranges []review.Range) ([]statement, *review.Finding) {
	parsed, err := pg.Parse(sql)
	if err != nil {
		return nil, syntaxFinding(err, ranges)
	}
	stmts := make([]statement, 0, len(parsed))
	for _, p := range parsed {
		if p.Empty() {
			continue
		}
		loc := ast.NodeLoc(p.AST)
		if loc.Start < 0 {
			loc = contentLoc(sql, p.ByteStart, p.ByteEnd)
		}
		stmts = append(stmts, statement{
			index: statementAt(ranges, loc.Start),
			node:  p.AST,
			loc:   loc,
		})
	}
	return stmts, nil
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
	f := &review.Finding{Rule: review.Syntax, Statement: -1, Message: err.Error()}
	var perr *parser.ParseError
	if errors.As(err, &perr) {
		f.Message = perr.Message
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
