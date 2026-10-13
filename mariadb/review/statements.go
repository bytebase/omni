package review

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/mariadb/parser"
	"github.com/bytebase/omni/review"
)

// statement is one parsed statement with its index in Result.Statements.
type statement struct {
	index int
	node  ast.Node
	// loc is the statement's absolute byte range in the input, from its
	// first token to its last, without the trailing delimiter.
	loc ast.Loc
}

// statementRanges lists every non-empty statement's byte range, from the
// lexical splitter, the same parser.Split Bytebase splits with. The ranges
// are the units parse works on.
func statementRanges(sql string) []review.Range {
	var ranges []review.Range
	for _, seg := range parser.Split(sql) {
		if seg.Empty() {
			continue
		}
		ranges = append(ranges, review.Range{Start: seg.ByteStart, End: seg.ByteEnd})
	}
	return ranges
}

// parse parses the text with parser.Parse, which splits it the same way
// and parses each segment in place, so every location and error position
// is an offset into the text. The first failure ends the review with its
// one Syntax finding: nothing after a syntax error can be trusted, and one
// finding per sheet is the contract. A non-empty range that produced no
// statement is a failure too, a guard against text the parser drops.
func parse(sql string, ranges []review.Range) ([]statement, *review.Finding) {
	list, err := parser.Parse(sql)
	if err != nil {
		return nil, syntaxFinding(err, ranges)
	}
	var stmts []statement
	found := make([]bool, len(ranges))
	if list != nil {
		for _, node := range list.Items {
			loc := nodeLoc(node)
			index := statementAt(ranges, loc.Start)
			if index < 0 {
				index = 0
			}
			if index >= len(ranges) {
				continue
			}
			// The statement's range is its segment's text: the parser starts
			// the location of a CREATE, ALTER, or DROP statement after its
			// first keyword.
			r := ranges[index]
			stmts = append(stmts, statement{index: index, node: node, loc: contentLoc(sql, r.Start, r.End)})
			found[index] = true
		}
	}
	for index, r := range ranges {
		if !found[index] {
			return nil, unparsedFinding(sql, index, r)
		}
	}
	sort.SliceStable(stmts, func(i, j int) bool { return stmts[i].loc.Start < stmts[j].loc.Start })
	return stmts, nil
}

// nodeLoc reads a statement node's Loc field. The MariaDB AST gives its
// nodes no common accessor; every statement struct carries one.
func nodeLoc(n ast.Node) ast.Loc {
	v := reflect.ValueOf(n)
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return ast.Loc{}
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return ast.Loc{}
	}
	f := v.FieldByName("Loc")
	if !f.IsValid() {
		return ast.Loc{}
	}
	loc, _ := f.Interface().(ast.Loc)
	return loc
}

// contentLoc is the range of a statement's text without the whitespace
// and comments before it, and the whitespace and semicolon after it. A
// /*! ... */ comment is not skipped: MariaDB runs the text inside it.
func contentLoc(sql string, start, end int) ast.Loc {
	start = skipTrivia(sql, start, end)
	for end > start && (isSpace(sql[end-1]) || sql[end-1] == ';') {
		end--
	}
	return ast.Loc{Start: start, End: end}
}

// skipTrivia returns the offset of the first byte in [start, end) that is
// not whitespace or a complete comment.
func skipTrivia(sql string, start, end int) int {
	i := start
	for i < end {
		switch {
		case isSpace(sql[i]):
			i++
		case sql[i] == '#', strings.HasPrefix(sql[i:end], "-- "), strings.HasPrefix(sql[i:end], "--\t"), strings.HasPrefix(sql[i:end], "--\n"), strings.HasPrefix(sql[i:end], "--\r"):
			for i < end && sql[i] != '\n' {
				i++
			}
		case strings.HasPrefix(sql[i:end], "/*") && !strings.HasPrefix(sql[i:end], "/*!"):
			j := strings.Index(sql[i+2:end], "*/")
			if j < 0 {
				return i
			}
			i += 2 + j + 2
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
// error carries the byte offset of the offending token and the message
// MariaDB would print; the finding's range starts and ends at that offset,
// which the caller anchors to its line. Any other error, which the parser
// does not produce today, addresses the whole change.
func syntaxFinding(err error, ranges []review.Range) *review.Finding {
	f := &review.Finding{Rule: review.Syntax, Statement: -1, Message: escapeLines(err.Error())}
	var perr *parser.ParseError
	if errors.As(err, &perr) {
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
// finding names the first token of the text the way a syntax error does.
func unparsedFinding(sql string, index int, r review.Range) *review.Finding {
	start := contentLoc(sql, r.Start, r.End).Start
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

// statementAt returns the index of the range containing the byte offset,
// or the last range starting before it, or -1 when none does.
func statementAt(ranges []review.Range, offset int) int {
	i := sort.Search(len(ranges), func(i int) bool { return ranges[i].Start > offset })
	return i - 1
}

// anchor returns the range a finding on a statement node anchors: the
// statement's own range for the statement itself, and otherwise the
// node's location, which for a statement of a routine body may start
// after its leading CREATE, ALTER, or DROP, extended back over it.
func anchor(sql string, s *statement, n ast.Node, loc ast.Loc) review.Range {
	if n == s.node {
		return rangeOf(s.loc)
	}
	i := loc.Start
	for i > 0 && i <= len(sql) && isSpace(sql[i-1]) {
		i--
	}
	j := i
	for j > 0 && isWordByte(sql[j-1]) {
		j--
	}
	switch strings.ToUpper(sql[j:i]) {
	case "CREATE", "ALTER", "DROP":
		loc.Start = j
	}
	return rangeOf(loc)
}

func isWordByte(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_' || c == '$'
}

// spanOf returns the range a comparison covers. The parser starts a binary
// expression's location at its operator and may end it past the space
// after it, and gives a NULL literal no location: the left operand's
// location, or a NULL just before the operator, widens it to the whole
// comparison, and trailing space is cut.
func spanOf(sql string, v *ast.BinaryExpr) review.Range {
	loc := v.Loc
	if l := nodeLoc(v.Left); l.End > l.Start && l.Start < loc.Start {
		loc.Start = l.Start
	} else if isNull(v.Left) {
		i := loc.Start
		for i > 0 && isSpace(sql[i-1]) {
			i--
		}
		if i >= 4 && strings.EqualFold(sql[i-4:i], "NULL") {
			loc.Start = i - 4
		}
	}
	if r := nodeLoc(v.Right); r.End > r.Start && r.End > loc.End {
		loc.End = r.End
	}
	for loc.End > loc.Start && loc.End <= len(sql) && isSpace(sql[loc.End-1]) {
		loc.End--
	}
	return rangeOf(loc)
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
// becomes one space, and text inside quotes or backticks is kept as
// written, with a line break escaped.
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
		case c == '\'' || c == '"' || c == '`':
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
