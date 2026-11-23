package analysis

import (
	"errors"
	"fmt"

	"github.com/bytebase/omni/databricks/parser"
)

// maxNestingDepth mirrors databricks/parser's nesting bound so deeply nested
// input fails closed with errNestingDepthExceeded.
const maxNestingDepth = 10000

var errNestingDepthExceeded = errors.New("maximum nesting depth exceeded")

// ValidateQuery reports whether statement may run in the SQL editor's read-only
// mode, and whether every statement in it returns rows:
// (readOnly, returnsData, error).
//
// Decision is AST-backed via databricks/parser: every non-empty statement in a
// script must parse as a read-only Databricks SQL query (SELECT, WITH ...
// SELECT/VALUES/TABLE/FROM pipeline, FROM ... query/pipeline, VALUES, TABLE,
// SHOW, DESCRIBE/DESC, or non-executing EXPLAIN over a read-only statement).
// External script transforms (MAP/REDUCE/SELECT TRANSFORM ... USING),
// EXPLAIN ANALYZE, EXPLAIN over writes, nested EXPLAIN, CTEs feeding writes,
// FROM-first inserts, DML, DDL, DCL, SET, USE, CALL, EXECUTE IMMEDIATE,
// compound BEGIN ... END blocks, and multi-statement batches containing any
// non-read-only statement are rejected. Empty, whitespace-only, comment-only,
// or semicolon-only input carries no statements to refuse and returns
// (true, true, nil). Lexical errors (unterminated comments, string literals, or
// quoted identifiers), syntax errors, and excessive nesting depth are returned
// with script-relative byte offsets so callers fail closed without losing
// diagnostic positions.
func ValidateQuery(statement string) (readOnly bool, returnsData bool, err error) {
	if err := firstLexError(statement); err != nil {
		return false, false, err
	}
	for _, seg := range parser.Split(statement) {
		if seg.Empty() {
			continue
		}
		if body, off, ok := explainAnalyzeBody(seg.Text); ok {
			if _, err := parser.ParseDatabricks(body); err != nil {
				return false, false, offsetParseError(err, seg.ByteStart+off)
			}
			return false, false, nil
		}
		res, err := parser.ParseDatabricks(seg.Text)
		if err != nil {
			offsetErr := offsetParseError(err, seg.ByteStart)
			var pe *parser.ParseError
			if errors.As(err, &pe) && pe.Message == errNestingDepthExceeded.Error() {
				return false, false, fmt.Errorf("%w: %w", errNestingDepthExceeded, offsetErr)
			}
			return false, false, offsetErr
		}
		if res == nil || len(res.Statements) == 0 {
			return false, false, nil
		}
		for _, stmt := range res.Statements {
			if !isReadOnlyAST(stmt) {
				return false, false, nil
			}
		}
	}
	return true, true, nil
}

func isReadOnlyAST(stmt parser.Statement) bool {
	ps, ok := stmt.(*parser.ParsedStmt)
	if !ok || ps == nil {
		return false
	}
	switch ps.Kind {
	case parser.StmtQuery:
		return !ps.HasTransform && !ps.HasEmptyPipe
	case parser.StmtShow, parser.StmtDescribe:
		return true
	case parser.StmtExplain:
		inner, ok := ps.Inner.(*parser.ParsedStmt)
		if !ok || inner == nil || inner.Kind == parser.StmtExplain {
			return false
		}
		return isReadOnlyAST(inner)
	default:
		return false
	}
}

// firstLexError scans src with databricks/parser's lexer and returns the first
// lexical error (unclosed block comment or unterminated string/backtick quote).
func firstLexError(src string) error {
	l := parser.NewLexer(src)
	for {
		tok := l.NextToken()
		if err := l.Err(); err != nil {
			return err
		}
		if tok.Kind == parser.TokEOF {
			return nil
		}
		if tok.Kind == parser.TokInvalid {
			switch tok.Text {
			case "'", "\"":
				return &parser.ParseError{
					Message:  fmt.Sprintf("unterminated string literal: missing closing %s", tok.Text),
					Position: tok.Start,
				}
			case "`":
				return &parser.ParseError{
					Message:  "unterminated quoted identifier: missing closing `",
					Position: tok.Start,
				}
			}
		}
	}
}

func offsetParseError(err error, base int) error {
	if base == 0 {
		return err
	}
	var pe *parser.ParseError
	if errors.As(err, &pe) {
		return &parser.ParseError{
			Message:  pe.Message,
			Position: pe.Position + base,
		}
	}
	return err
}
