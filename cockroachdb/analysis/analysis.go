// Package analysis is CockroachDB's statement judgment: the statement types a
// script changes, the tables it writes, its system catalogs, and the editor's
// splitting, diagnostics, read-only check, result limit and EXPLAIN. The wire
// protocol is PostgreSQL's, but CockroachDB's grammar is not: UPSERT,
// statement sources, AS OF SYSTEM TIME, index hints, regional DDL, zone
// configuration and IMPORT INTO have no PostgreSQL node, so the judgment walks
// CockroachDB's own AST.
package analysis

import (
	"strings"
	"unicode"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/lexbase"
	crrawparser "github.com/cockroachdb/cockroachdb-parser/pkg/sql/parser"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/parser/statements"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/pgwire/pgerror"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/scanner"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"

	// Registers parser.ParseDoBlockFn; without it the SQL parser panics on a DO block.
	_ "github.com/cockroachdb/cockroachdb-parser/pkg/sql/plpgsql/parser"
)

// ParsedStatement is an opaque parsed CockroachDB statement that keeps
// cockroachdb-parser AST types encapsulated inside this package.
type ParsedStatement struct {
	stmt statements.Statement[tree.Statement]
}

// Text returns the statement SQL as CockroachDB's parser read it.
func (s ParsedStatement) Text() string {
	return s.stmt.SQL
}

// ParseErrorInfo holds normalized error fields extracted from a CockroachDB
// parser error so callers outside this package do not depend on vendor error
// types. Position is the byte offset into the parsed SQL, or -1 when the parser
// reported no source location.
type ParseErrorInfo struct {
	Message  string
	Position int
	Code     string
}

// ParseRaw parses sql unmodified into opaque ParsedStatement values. On a
// parse failure it returns normalized error details with the byte offset into
// sql.
func ParseRaw(sql string) ([]ParsedStatement, *ParseErrorInfo, error) {
	stmts, err := crrawparser.Parse(sql)
	if err != nil {
		return nil, extractParseErrorInfo(sql, err), err
	}
	out := make([]ParsedStatement, len(stmts))
	for i, stmt := range stmts {
		out[i] = ParsedStatement{stmt: stmt}
	}
	return out, nil, nil
}

// ParseScript parses script, the whole script or one statement, into opaque
// ParsedStatement values.
func ParseScript(script string) ([]ParsedStatement, error) {
	stmts, err := Parse(script)
	if err != nil {
		return nil, err
	}
	out := make([]ParsedStatement, len(stmts))
	for i, stmt := range stmts {
		out[i] = ParsedStatement{stmt: stmt}
	}
	return out, nil
}

// SplitScript splits script lexically into statement segments, preserving
// empty segments as "" so callers can account for every semicolon boundary.
func SplitScript(script string) []string {
	return Split(script)
}

func unwrapStatements(stmts []ParsedStatement) statements.Statements {
	out := make(statements.Statements, len(stmts))
	for i, stmt := range stmts {
		out[i] = stmt.stmt
	}
	return out
}

// Parse parses script, the whole script or one statement, into statements
// carrying the AST and the text CockroachDB's parser read. A trailing
// semicolon is written so a script without one parses as CockroachDB's
// parser expects.
func Parse(script string) (statements.Statements, error) {
	script = strings.TrimRightFunc(script, isSpaceOrSemicolon) + ";"
	return crrawparser.Parse(script)
}

// Split returns each statement of script with its terminating semicolon,
// lexically, so a script the parser rejects still splits. A segment that holds
// no statement, such as a lone ";", is "" so a caller can account for every
// boundary; trailing text holding no statement is dropped.
func Split(script string) []string {
	segments := SplitSegments(script)
	result := make([]string, 0, len(segments))
	for i, segment := range segments {
		switch {
		case !segment.Empty():
			result = append(result, strings.TrimLeftFunc(segment.Text, unicode.IsSpace))
		case i < len(segments)-1 || strings.HasSuffix(segment.Text, ";"):
			result = append(result, "")
		}
	}
	return result
}

// Segment is one statement's span of a script: the text from the end of the
// previous segment through its terminating semicolon.
type Segment struct {
	Text      string
	ByteStart int
	ByteEnd   int
	empty     bool
}

// Empty reports whether the segment holds no statement: only whitespace,
// comments and semicolons.
func (s Segment) Empty() bool {
	return s.empty
}

// SplitSegments splits script lexically into segments that tile it, as
// CockroachDB's parser bounds statements: a semicolon inside a BEGIN ATOMIC
// body does not end its CREATE FUNCTION. After a lexical error, such as an
// unterminated string, the rest of the script is one segment.
func SplitSegments(script string) []Segment {
	var sc scanner.SQLScanner
	sc.Init(script)
	defer sc.Cleanup()

	var segments []Segment
	var lval scanSym
	start, holds, bodyDepth := 0, false, 0
	var prevID int32
	for {
		sc.Scan(&lval)
		if lval.id == 0 || lval.id == lexbase.ERROR {
			holds = holds || lval.id == lexbase.ERROR
			break
		}
		if lval.id == ';' && bodyDepth == 0 {
			end := sc.Pos()
			segments = append(segments, Segment{Text: script[start:end], ByteStart: start, ByteEnd: end, empty: !holds})
			start, holds = end, false
			prevID = lval.id
			continue
		}
		if lval.id != ';' {
			holds = true
		}
		// The same body tracking as cockroachdb-parser's scanOneStmt.
		if prevID == lexbase.BEGIN && lval.id == lexbase.ATOMIC {
			bodyDepth++
		}
		if bodyDepth > 0 && lval.id == lexbase.END {
			bodyDepth--
		}
		prevID = lval.id
	}
	if start < len(script) {
		segments = append(segments, Segment{Text: script[start:], ByteStart: start, ByteEnd: len(script), empty: !holds})
	}
	return segments
}

func isSpaceOrSemicolon(r rune) bool {
	return unicode.IsSpace(r) || r == ';'
}

type stmtSpan struct {
	start int
	end   int
}

type scanSym struct {
	id  int32
	pos int32
	str string
	val any
}

func (s *scanSym) ID() int32           { return s.id }
func (s *scanSym) SetID(id int32)      { s.id = id }
func (s *scanSym) Pos() int32          { return s.pos }
func (s *scanSym) SetPos(pos int32)    { s.pos = pos }
func (s *scanSym) Str() string         { return s.str }
func (s *scanSym) SetStr(str string)   { s.str = str }
func (s *scanSym) UnionVal() any       { return s.val }
func (s *scanSym) SetUnionVal(val any) { s.val = val }

// scanStmtSpans mirrors cockroachdb-parser's scanOneStmt loop so each
// statement's exact byte range in sql can be recovered when reporting a parse
// error's byte offset.
func scanStmtSpans(sql string) []stmtSpan {
	var sc scanner.SQLScanner
	sc.Init(sql)
	defer sc.Cleanup()

	var spans []stmtSpan
	var lval scanSym
	for {
		for {
			sc.Scan(&lval)
			if lval.id == 0 {
				return spans
			}
			if lval.id != ';' {
				break
			}
		}
		start := int(lval.pos)
		curFuncBodyCnt := 0
		for {
			if lval.id == lexbase.ERROR {
				return append(spans, stmtSpan{start: start, end: len(sql)})
			}
			prevID := lval.id
			sc.Scan(&lval)
			if prevID == lexbase.BEGIN && lval.id == lexbase.ATOMIC {
				curFuncBodyCnt++
			}
			if curFuncBodyCnt > 0 && lval.id == lexbase.END {
				curFuncBodyCnt--
			}
			if lval.id == 0 || (curFuncBodyCnt == 0 && lval.id == ';') {
				end := sc.Pos()
				if lval.id == ';' {
					end--
				}
				spans = append(spans, stmtSpan{start: start, end: end})
				if lval.id == 0 {
					return spans
				}
				break
			}
		}
	}
}

func extractParseErrorInfo(sql string, err error) *ParseErrorInfo {
	flat := pgerror.Flatten(err)
	info := &ParseErrorInfo{
		Message:  flat.Message,
		Position: -1,
		Code:     flat.Code,
	}
	snippet, relPos, ok := extractSourceSQLCaret(flat.Detail)
	if !ok {
		return info
	}
	span, found := failingStmtSpan(sql)
	if !found {
		return info
	}
	stmtText := sql[span.start:span.end]
	if strings.HasPrefix(stmtText, snippet) {
		pos := span.start + relPos
		if pos >= 0 && pos <= len(sql) {
			info.Position = pos
		}
		return info
	}
	if snippet != "" {
		if idx := strings.Index(stmtText, snippet); idx >= 0 {
			pos := span.start + idx + relPos
			if pos >= 0 && pos <= len(sql) {
				info.Position = pos
			}
		}
	}
	return info
}

func failingStmtSpan(sql string) (stmtSpan, bool) {
	spans := scanStmtSpans(sql)
	if len(spans) == 0 {
		return stmtSpan{}, false
	}
	if len(spans) == 1 {
		return spans[0], true
	}
	for _, span := range spans {
		if _, err := crrawparser.Parse(sql[span.start:span.end]); err != nil {
			return span, true
		}
	}
	return spans[0], true
}

func extractSourceSQLCaret(detail string) (snippet string, relPos int, ok bool) {
	if detail == "" {
		return "", -1, false
	}
	for _, block := range strings.Split(detail, "\n--\n") {
		if s, p, matched := parseSourceSQLBlock(block); matched {
			return s, p, true
		}
	}
	return "", -1, false
}

func parseSourceSQLBlock(block string) (snippet string, relPos int, ok bool) {
	const prefix = "source SQL:\n"
	if !strings.HasPrefix(block, prefix) {
		return "", -1, false
	}
	rest := block[len(prefix):]
	nl := strings.LastIndexByte(rest, '\n')
	if nl < 0 {
		return "", -1, false
	}
	snippet = rest[:nl]
	caretLine := rest[nl+1:]
	if len(caretLine) == 0 || caretLine[len(caretLine)-1] != '^' {
		return "", -1, false
	}
	for i := 0; i < len(caretLine)-1; i++ {
		if caretLine[i] != ' ' {
			return "", -1, false
		}
	}
	lastLineStart := strings.LastIndexByte(snippet, '\n') + 1
	return snippet, lastLineStart + len(caretLine) - 1, true
}
