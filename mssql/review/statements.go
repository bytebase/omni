package review

import (
	"errors"
	"sort"
	"strings"

	"github.com/bytebase/omni/mssql"
	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/mssql/parser"
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

// parseReview parses the text with mssql.Parse, the entry Bytebase parses
// SQL Server with, and returns the statements and their ranges. T-SQL has
// no lexical splitter: statements need not end with a semicolon, so only
// the parser tells where one ends. A text that does not parse is split at
// its semicolons, and gets its Syntax finding: nothing after a syntax
// error can be trusted, and one finding per sheet is the contract. A GO
// batch separator is a statement of its own, which no rule reads.
func parseReview(sql string) ([]review.Range, []statement, *review.Finding) {
	parsed, err := mssql.Parse(sql)
	if err != nil {
		ranges := semicolonRanges(sql)
		return ranges, nil, syntaxFinding(err, ranges)
	}
	var ranges []review.Range
	var stmts []statement
	for _, p := range parsed {
		if p.Empty() {
			continue
		}
		loc := contentLoc(sql, p.ByteStart, p.ByteEnd)
		if loc.End <= loc.Start {
			continue
		}
		index := len(ranges)
		ranges = append(ranges, review.Range{Start: loc.Start, End: loc.End})
		stmts = append(stmts, statement{index: index, node: p.AST, loc: loc})
	}
	return ranges, stmts, nil
}

// semicolonRanges splits a text that does not parse at its semicolons, as
// Bytebase does then, so the review lists the statements Bytebase lists.
// Without a tree, nothing better tells where a T-SQL statement ends.
//
// bytebase: backend/plugin/parser/tsql/split.go — SplitSQL
func semicolonRanges(sql string) []review.Range {
	var ranges []review.Range
	add := func(start, end int) {
		if loc := contentLoc(sql, start, end); loc.End > loc.Start {
			ranges = append(ranges, review.Range{Start: loc.Start, End: loc.End})
		}
	}
	start := 0
	for _, tok := range parser.Tokenize(sql) {
		if tok.Type == ';' {
			add(start, tok.Loc+1)
			start = tok.Loc + 1
		}
	}
	add(start, len(sql))
	return ranges
}

// contentLoc is the range of a statement's text without the whitespace
// and comments before it, and the whitespace and semicolon after it.
func contentLoc(sql string, start, end int) ast.Loc {
	start = skipTrivia(sql, start, end)
	for end > start && (isSpace(sql[end-1]) || sql[end-1] == ';') {
		end--
	}
	return ast.Loc{Start: start, End: end}
}

// skipTrivia returns the offset of the first byte in [start, end) that is
// not whitespace, a complete comment, or the semicolon of an empty
// statement before it. T-SQL block comments nest.
func skipTrivia(sql string, start, end int) int {
	i := start
	for i < end {
		switch {
		case isSpace(sql[i]), sql[i] == ';':
			i++
		case strings.HasPrefix(sql[i:end], "--"):
			for i < end && sql[i] != '\n' {
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

// isSpace is the lexer's whitespace set.
func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' || c == '\v'
}

// syntaxFinding turns a parse failure into the Syntax finding. A parser
// error carries the byte offset of the offending token and its message;
// the finding's range starts and ends at that offset, which the caller
// anchors to its line. Any other error addresses the whole change.
func syntaxFinding(err error, ranges []review.Range) *review.Finding {
	f := &review.Finding{Rule: review.Syntax, Statement: -1, Message: escapeLines(err.Error())}
	var perr *parser.ParseError
	if errors.As(err, &perr) {
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

// oneLine writes source text on one line: whitespace between tokens
// becomes one space, and text inside quotes or brackets is kept as
// written, with a line break escaped.
func oneLine(text string) string {
	var b strings.Builder
	var quote byte
	space := false
	for i := 0; i < len(text); i++ {
		c := text[i]
		switch {
		case quote == ']' && c == ']' && i+1 < len(text) && text[i+1] == ']':
			// ]] is a ] inside a bracketed name: both are kept, and the name
			// stays open.
			b.WriteByte(c)
			i++
		case quote != 0:
			if c == quote {
				quote = 0
			}
		case c == '\'' || c == '"' || c == '[':
			quote = c
			if c == '[' {
				quote = ']'
			}
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
