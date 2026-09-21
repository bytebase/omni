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
	// loc is the statement's absolute byte range in the input, without
	// the trailing semicolon.
	loc ast.Loc
}

// statementRanges lists every non-empty statement's byte range, from the
// lexical splitter so it does not depend on the text parsing.
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

// parse parses the whole text. On success it returns the statements
// mapped onto ranges. On failure it returns the one Syntax finding of the
// review: the parser stops at the first error, so there is never more
// than one, and the review ends there.
func parse(sql string, ranges []review.Range) ([]statement, *review.Finding) {
	list, err := parser.Parse(sql)
	if err != nil {
		return nil, syntaxFinding(err, ranges)
	}
	if list == nil {
		return nil, nil
	}
	stmts := make([]statement, 0, len(list.Items))
	for _, item := range list.Items {
		raw, ok := item.(*ast.RawStmt)
		if !ok || raw.Stmt == nil {
			continue
		}
		stmts = append(stmts, statement{
			index: statementAt(ranges, raw.Loc.Start),
			node:  raw.Stmt,
			loc:   raw.Loc,
		})
	}
	return stmts, nil
}

// syntaxFinding turns a parse failure into the Syntax finding. A parser
// error carries the byte offset of the offending token; the finding's
// range starts and ends there, which the caller anchors to that line. An
// error without a position addresses the whole change.
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
