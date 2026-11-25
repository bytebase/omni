package completion

import (
	"strings"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/lexbase"
	crparser "github.com/cockroachdb/cockroachdb-parser/pkg/sql/parser"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/pgwire/pgerror"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/scanner"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"

	// Registers parser.ParseDoBlockFn; without it the SQL parser panics on a DO block.
	_ "github.com/cockroachdb/cockroachdb-parser/pkg/sql/plpgsql/parser"
)

// parseOne parses text as a single statement. A panic in the vendor parser is
// reported as a failure, since completion runs on whatever the user typed.
func parseOne(text string) (stmt tree.Statement, err error) {
	defer func() {
		if r := recover(); r != nil {
			stmt, err = nil, errParsePanic
		}
	}()
	stmts, err := crparser.Parse(text)
	if err != nil {
		return nil, err
	}
	if len(stmts) != 1 {
		return nil, errNotOneStatement
	}
	return stmts[0].AST, nil
}

type parseFailure string

func (e parseFailure) Error() string { return string(e) }

const (
	errParsePanic      = parseFailure("parser panicked")
	errNotOneStatement = parseFailure("not exactly one statement")
)

// reachesPast reports whether CockroachDB's parser, reading text, shifts the
// token starting at byte offset tokStart: the parse succeeds or fails at a
// later token. goyacc's LALR(1) tables never shift a token that cannot follow
// the prefix, so this is exactly "the token is valid here".
func reachesPast(text string, tokStart int) bool {
	_, err := parseOne(text)
	if err == nil {
		return true
	}
	return errorOffset(text, err) > tokStart
}

// errorOffset returns the byte offset in text of the token CockroachDB's
// parser failed at, read from the "source SQL" caret the parser attaches to
// a syntax error, or -1 when the error carries none.
func errorOffset(text string, err error) int {
	if _, ok := err.(parseFailure); ok {
		return -1
	}
	detail := pgerror.Flatten(err).Detail
	const header = "source SQL:\n"
	for _, block := range strings.Split(detail, "\n--\n") {
		if !strings.HasPrefix(block, header) {
			continue
		}
		block = block[len(header):]
		nl := strings.LastIndexByte(block, '\n')
		if nl < 0 {
			continue
		}
		snippet, caretLine := block[:nl], block[nl+1:]
		if !strings.HasSuffix(caretLine, "^") || strings.Trim(caretLine[:len(caretLine)-1], " ") != "" {
			continue
		}
		lineStart := strings.LastIndexByte(snippet, '\n') + 1
		return firstTokenStart(text) + lineStart + len(caretLine) - 1
	}
	return -1
}

// firstTokenStart returns where the parser's statement text begins: the
// vendor parser drops leading whitespace, comments and semicolons and reports
// positions relative to the first token.
func firstTokenStart(text string) int {
	var s scanner.SQLScanner
	s.Init(text)
	defer s.Cleanup()
	var sym scanSym
	for {
		s.Scan(&sym)
		if sym.id == 0 {
			return len(text)
		}
		if sym.id != ';' {
			return int(sym.pos)
		}
	}
}

// token is one lexical token of CockroachDB SQL. end is where the scanner
// stopped, which for a quoted token is one byte past its closing quote.
type token struct {
	id          int32
	start, end  int
	str         string
	quotedIdent bool // an unterminated double-quoted identifier, as an ERROR token
}

// tokenize scans text with CockroachDB's scanner. A lexical error ends the
// list with an ERROR token.
func tokenize(text string) (out []token) {
	defer func() {
		if recover() != nil {
			out = []token{{id: lexbase.ERROR, start: 0, end: len(text)}}
		}
	}()
	for _, t := range scanner.Inspect(text) {
		if t.ID == 0 {
			break
		}
		out = append(out, token{
			id:          t.ID,
			start:       int(t.Start),
			end:         int(t.End),
			str:         t.Str,
			quotedIdent: t.ID == lexbase.ERROR && t.MaybeID == lexbase.IDENT && t.Quoted,
		})
		if t.ID == lexbase.ERROR {
			break
		}
	}
	return out
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
