package analysis

import (
	"errors"
	"fmt"

	"github.com/bytebase/omni/clickhouse/parser"
)

// ValidateQuery reports whether statement is a single read-only ClickHouse
// query that returns rows.
//
// Classification uses clickhouse/parser's lexer and AST: SELECT (including
// parenthesized, FROM-first, WITH/CTE, UNION/INTERSECT/EXCEPT, FORMAT, and
// SETTINGS queries, excluding INTO OUTFILE), SHOW, DESCRIBE/DESC, EXISTS,
// CHECK, and non-executing EXPLAIN statements are allowed. Multi-statement
// batches, EXPLAIN ANALYZE (which executes what it plans), WITH ... DML, and
// all DDL/DML/admin statements are rejected. Lexical errors (unterminated
// strings, dollar-quoted strings, quoted identifiers, or block comments) are
// returned as errors.
func ValidateQuery(statement string) (readOnly bool, returnsData bool, err error) {
	if err := firstLexError(statement); err != nil {
		return false, false, err
	}
	seg, ok := singleNonEmptySegment(statement)
	if !ok {
		return false, false, nil
	}
	stmts, err := parser.Parse(seg.Text)
	if err != nil || len(stmts) != 1 {
		return false, false, nil
	}
	if !isReadOnlyAST(stmts[0]) {
		return false, false, nil
	}
	return true, true, nil
}

// StatementReadsRows reports whether the server answers statement with rows.
// EXPLAIN ANALYZE is classified as non-read-only by ValidateQuery because it
// executes what it plans, yet the server still returns the plan as rows.
func StatementReadsRows(statement string) (bool, error) {
	if err := firstLexError(statement); err != nil {
		return false, err
	}
	seg, ok := singleNonEmptySegment(statement)
	if !ok {
		return false, nil
	}
	stmts, err := parser.Parse(seg.Text)
	if err == nil && len(stmts) == 1 {
		return astReadsRows(stmts[0]), nil
	}
	// If a single EXPLAIN statement wraps syntax outside the parser grammar,
	// ClickHouse still answers an EXPLAIN statement with rows.
	return inspectExplain(seg.Text).isExplain, nil
}

// CheckStatementWithoutQuotedTextAndComment classifies a single statement as a
// read-only query. Retained for API compatibility; classification now uses
// clickhouse/parser directly.
func CheckStatementWithoutQuotedTextAndComment(statement string) bool {
	readOnly, _, err := ValidateQuery(statement)
	return err == nil && readOnly
}

func singleNonEmptySegment(script string) (parser.Segment, bool) {
	var found parser.Segment
	count := 0
	for _, seg := range parser.Split(script) {
		if seg.Empty() {
			continue
		}
		found = seg
		count++
		if count > 1 {
			return parser.Segment{}, false
		}
	}
	return found, count == 1
}

func isReadOnlyAST(stmt parser.Statement) bool {
	switch s := stmt.(type) {
	case *parser.SelectUnionStmt:
		return !hasOutfile(s.Output)
	case *parser.ShowStmt:
		return !hasOutfile(s.Output)
	case *parser.DescribeStmt:
		return !hasOutfile(s.Output)
	case *parser.ExistsStmt:
		return !hasOutfile(s.Output)
	case *parser.CheckStmt:
		return !hasOutfile(s.Output)
	case *parser.ExplainStmt:
		if hasOutfile(s.Output) {
			return false
		}
		return s.Kind != "ANALYZE"
	default:
		return false
	}
}

func astReadsRows(stmt parser.Statement) bool {
	switch s := stmt.(type) {
	case *parser.ExplainStmt:
		return !hasOutfile(s.Output)
	default:
		return isReadOnlyAST(stmt)
	}
}

func hasOutfile(out *parser.OutputClause) bool {
	return out != nil && out.Outfile != ""
}

func firstLexError(text string) error {
	l := parser.NewLexer(text)
	for l.NextToken().Type != parser.TokEOF {
		if err := l.Err(); err != nil {
			return wrapLexError(err)
		}
	}
	if err := l.Err(); err != nil {
		return wrapLexError(err)
	}
	return nil
}

// wrapLexError preserves *parser.ParseError unwrapping while keeping the
// "invalid string" phrase expected by legacy callers for unterminated string
// literals.
func wrapLexError(err error) error {
	var pe *parser.ParseError
	if errors.As(err, &pe) && pe.Message == "unterminated string literal" {
		return fmt.Errorf("invalid string: %w", err)
	}
	return err
}
