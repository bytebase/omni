package review

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/bytebase/omni/pg"
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/internal/metacmd"
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

// statementRanges lists every non-empty statement's byte range, from the
// lexical splitter, the same pg.Split Bytebase splits with. The ranges
// are the units parse works on.
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

// parse splits the text with pg.Split and parses each non-empty segment
// in place with parser.ParseRange, so every location and error position
// is an offset into the text. The first segment that fails ends the
// review with its one Syntax finding: nothing after a syntax error can
// be trusted, and one finding per sheet is the contract. A non-empty
// segment that produced no statement is a failure too, a guard against
// text the parser drops.
func parse(sql string, ranges []review.Range) ([]statement, *review.Finding) {
	var stmts []statement
	for index, r := range ranges {
		list, err := parser.ParseRange(sql, r.Start, r.End)
		if err != nil {
			return nil, syntaxFinding(err, ranges)
		}
		found := false
		if list != nil {
			for _, item := range list.Items {
				raw, ok := item.(*ast.RawStmt)
				if !ok || raw.Stmt == nil {
					continue
				}
				// A node the parser gave no location has NoLoc, or the
				// zero Loc when its kind carries one the parser did not
				// fill in; either falls outside the segment's bytes and is
				// replaced by them.
				loc := raw.Loc
				if loc.Start < r.Start || loc.End > r.End || loc.End <= loc.Start {
					loc = contentLoc(sql, r.Start, r.End)
				}
				stmts = append(stmts, statement{index: index, node: raw.Stmt, loc: loc})
				found = true
			}
		}
		if !found {
			return nil, unparsedFinding(sql, index, r)
		}
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

// isSpace is the lexer's whitespace set.
func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' || c == '\v'
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

// unparsedFinding is the Syntax finding for a non-empty range the parser
// produced no statement for. The parser gives no message for it, so the
// finding names the first token of the text, past any comment or psql
// metacommand before it, the way a syntax error does.
func unparsedFinding(sql string, index int, r review.Range) *review.Finding {
	start := skipTrivia(sql, r.Start, r.End)
	text := sql[start:r.End]
	if i := strings.IndexAny(text, " \t\r\n\f\v"); i >= 0 {
		text = text[:i]
	}
	if len(text) > 40 {
		text = text[:40]
	}
	return &review.Finding{
		Rule:      review.Syntax,
		Statement: index,
		Range:     review.Range{Start: start, End: start},
		Message:   escapeLines(fmt.Sprintf("syntax error at or near %q", text)),
	}
}

// skipTrivia returns the offset of the first byte in [start, end) that is
// not whitespace, a complete comment, or a psql metacommand line. An
// unterminated block comment is not trivia: it is the text at fault.
func skipTrivia(sql string, start, end int) int {
	i := start
	for i < end {
		switch {
		case isSpace(sql[i]):
			i++
		case metacmd.IsMetaCommand(sql, i):
			i = metacmd.SkipLine(sql, i)
		case strings.HasPrefix(sql[i:end], "--"):
			for i < end && sql[i] != '\n' && sql[i] != '\r' {
				i++
			}
		case strings.HasPrefix(sql[i:end], "/*"):
			depth, j := 0, i
			for j < end {
				switch {
				case strings.HasPrefix(sql[j:end], "/*"):
					depth++
					j += 2
				case strings.HasPrefix(sql[j:end], "*/"):
					depth--
					j += 2
				default:
					j++
				}
				if depth == 0 {
					break
				}
			}
			if depth > 0 {
				return i
			}
			i = j
		default:
			return i
		}
	}
	return i
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
