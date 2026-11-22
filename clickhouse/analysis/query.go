package analysis

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/clickhouse/parser"
)

// QueryType is the permission class of a statement.
type QueryType = analysis.QueryType

// ValidateQuery reports whether statement is a single read-only ClickHouse
// query that returns rows.
//
// Classification uses clickhouse/parser's lexer and AST: SELECT (including
// parenthesized, FROM-first, WITH/CTE, UNION/INTERSECT/EXCEPT, FORMAT, and
// SETTINGS queries, excluding INTO OUTFILE), SHOW, DESCRIBE/DESC, EXISTS,
// CHECK, and non-executing EXPLAIN statements (including EXPLAIN TABLE OVERRIDE
// and EXPLAIN CURRENT TRANSACTION) are allowed. Empty/comment-only inputs and
// multi-statement batches return (false, false, nil). EXPLAIN ANALYZE (which
// executes what it plans) and all DDL/DML/admin statements return
// (false, false, nil). Lexical and syntax errors from clickhouse/parser are
// propagated with statement-relative byte offsets.
func ValidateQuery(statement string) (readOnly bool, returnsData bool, err error) {
	seg, ok := singleNonEmptySegment(statement)
	if !ok {
		return false, false, nil
	}
	stmts, err := parser.Parse(seg.Text)
	if err != nil {
		return false, false, wrapLexError(offsetParseError(err, seg.ByteStart))
	}
	if len(stmts) != 1 || !isReadOnlyAST(stmts[0]) {
		return false, false, nil
	}
	return true, true, nil
}

// Classify parses a single ClickHouse SQL statement and returns its QueryType.
// Empty or comment-only input returns QueryTypeUnknown with a nil error;
// multi-statement input or invalid syntax returns an error.
func Classify(statement string) (QueryType, error) {
	var nonEmpty []parser.Segment
	for _, seg := range parser.Split(statement) {
		if !seg.Empty() {
			nonEmpty = append(nonEmpty, seg)
		}
	}
	if len(nonEmpty) == 0 {
		return analysis.QueryTypeUnknown, nil
	}
	if len(nonEmpty) > 1 {
		return analysis.QueryTypeUnknown, errors.New("expected a single ClickHouse statement")
	}
	stmts, err := parser.Parse(nonEmpty[0].Text)
	if err != nil {
		return analysis.QueryTypeUnknown, wrapLexError(offsetParseError(err, nonEmpty[0].ByteStart))
	}
	if len(stmts) != 1 {
		return analysis.QueryTypeUnknown, errors.New("expected a single ClickHouse statement")
	}
	return classifyParsedStatement(stmts[0]), nil
}

// StatementReadsRows reports whether the server answers statement with rows.
// EXPLAIN ANALYZE is classified as non-read-only by ValidateQuery because it
// executes what it plans, yet the server still returns the plan as rows.
// Lexical and syntax errors from clickhouse/parser are propagated with
// statement-relative byte offsets rather than downgraded to a non-row command.
func StatementReadsRows(statement string) (bool, error) {
	seg, ok := singleNonEmptySegment(statement)
	if !ok {
		return false, nil
	}
	stmts, err := parser.Parse(seg.Text)
	if err != nil {
		return false, wrapLexError(offsetParseError(err, seg.ByteStart))
	}
	if len(stmts) != 1 {
		return false, nil
	}
	return astReadsRows(stmts[0]), nil
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
	if stmt == nil {
		return false
	}
	switch s := stmt.(type) {
	case *parser.SelectUnionStmt:
		return !hasOutfile(s.Output)
	case *parser.SelectStmt:
		return true
	case *parser.ShowStmt:
		return !hasOutfile(s.Output)
	case *parser.DescribeStmt:
		return !hasOutfile(s.Output)
	case *parser.ExistsStmt:
		return !hasOutfile(s.Output)
	case *parser.CheckStmt:
		return !hasOutfile(s.Output)
	case *parser.ExplainStmt:
		if hasOutfile(s.Output) || strings.EqualFold(s.Kind, "ANALYZE") {
			return false
		}
		return true
	default:
		return false
	}
}

func astReadsRows(stmt parser.Statement) bool {
	switch s := stmt.(type) {
	case *parser.ExplainStmt:
		return !hasOutfile(s.Output) && (s.Inner == nil || !statementHasOutfile(s.Inner))
	default:
		return isReadOnlyAST(stmt)
	}
}

func classifyParsedStatement(stmt parser.Statement) QueryType {
	if stmt == nil {
		return analysis.QueryTypeUnknown
	}
	switch stmt.(type) {
	case *parser.SelectUnionStmt, *parser.SelectStmt:
		if statementHasOutfile(stmt) {
			return analysis.DML
		}
		return analysis.Select
	case *parser.ExplainStmt:
		return analysis.Explain
	case *parser.ShowStmt, *parser.DescribeStmt, *parser.ExistsStmt, *parser.CheckStmt:
		if statementHasOutfile(stmt) {
			return analysis.DML
		}
		return analysis.SelectInfoSchema
	case *parser.InsertStmt, *parser.UpdateStmt, *parser.DeleteStmt:
		return analysis.DML
	case *parser.CreateDatabaseStmt, *parser.CreateTableStmt, *parser.CreateViewStmt,
		*parser.CreateDictionaryStmt, *parser.CreateFunctionStmt, *parser.AlterTableStmt,
		*parser.DropStmt, *parser.AttachStmt, *parser.RenameStmt, *parser.ExchangeStmt,
		*parser.TruncateStmt, *parser.UndropStmt, *parser.OptimizeStmt:
		return analysis.DDL
	default:
		return analysis.QueryTypeUnknown
	}
}

func statementHasOutfile(stmt parser.Statement) bool {
	switch s := stmt.(type) {
	case *parser.SelectUnionStmt:
		return hasOutfile(s.Output)
	case *parser.ShowStmt:
		return hasOutfile(s.Output)
	case *parser.DescribeStmt:
		return hasOutfile(s.Output)
	case *parser.ExplainStmt:
		return hasOutfile(s.Output)
	case *parser.ExistsStmt:
		return hasOutfile(s.Output)
	case *parser.CheckStmt:
		return hasOutfile(s.Output)
	default:
		return false
	}
}

func hasOutfile(out *parser.OutputClause) bool {
	return out != nil && out.Outfile != ""
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
