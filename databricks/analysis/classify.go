package analysis

import (
	"errors"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/databricks/parser"
)

// QueryType is the permission class of a statement.
type QueryType = analysis.QueryType

// Classify parses a single Databricks SQL statement and returns its QueryType.
// Empty, whitespace-only, comment-only, or semicolon-only input returns
// QueryTypeUnknown with a nil error; multi-statement input or lexical/syntax
// errors return an error with statement-relative byte offsets.
func Classify(statement string) (QueryType, error) {
	if err := firstLexError(statement); err != nil {
		return analysis.QueryTypeUnknown, err
	}
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
		return analysis.QueryTypeUnknown, errors.New("expected a single Databricks SQL statement")
	}
	seg := nonEmpty[0]
	if body, off, ok := explainAnalyzeBody(seg.Text); ok {
		res, err := parser.Parse(body)
		if err != nil {
			return analysis.QueryTypeUnknown, offsetParseError(err, seg.ByteStart+off)
		}
		if res == nil || len(res.Statements) != 1 {
			return analysis.QueryTypeUnknown, errors.New("expected a single Databricks SQL statement")
		}
		return analysis.Explain, nil
	}
	res, err := parser.Parse(seg.Text)
	if err != nil {
		return analysis.QueryTypeUnknown, offsetParseError(err, seg.ByteStart)
	}
	if res == nil || len(res.Statements) != 1 {
		return analysis.QueryTypeUnknown, errors.New("expected a single Databricks SQL statement")
	}
	return classifyParsedStatement(res.Statements[0]), nil
}

// StatementReadsRows reports whether the server answers statement with rows.
// SELECT, WITH (when its main statement is a query), FROM-first queries,
// SHOW, DESCRIBE/DESC, EXPLAIN (including EXPLAIN ANALYZE), TABLE, and VALUES
// return rows. Empty/comment-only inputs or multi-statement scripts return
// (false, nil). Lexical and syntax errors from databricks/parser are propagated
// with statement-relative byte offsets.
func StatementReadsRows(statement string) (bool, error) {
	if err := firstLexError(statement); err != nil {
		return false, err
	}
	seg, ok := singleNonEmptySegment(statement)
	if !ok {
		return false, nil
	}
	if body, off, isEA := explainAnalyzeBody(seg.Text); isEA {
		res, err := parser.Parse(body)
		if err != nil {
			return false, offsetParseError(err, seg.ByteStart+off)
		}
		return res != nil && len(res.Statements) == 1, nil
	}
	res, err := parser.Parse(seg.Text)
	if err != nil {
		return false, offsetParseError(err, seg.ByteStart)
	}
	if res == nil || len(res.Statements) != 1 {
		return false, nil
	}
	return astReadsRows(res.Statements[0]), nil
}

// ClassifyStatement reports whether statement is expected to return rows.
// SELECT, WITH (unless it writes), FROM-first queries, SHOW, DESCRIBE, DESC,
// EXPLAIN, TABLE, and VALUES are reads; everything else (or an unparseable
// statement) is a command.
func ClassifyStatement(statement string) bool {
	reads, err := StatementReadsRows(statement)
	return err == nil && reads
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

func astReadsRows(stmt parser.Statement) bool {
	ps, ok := stmt.(*parser.ParsedStmt)
	if !ok || ps == nil {
		return false
	}
	switch ps.Kind {
	case parser.StmtQuery, parser.StmtShow, parser.StmtDescribe, parser.StmtExplain:
		return true
	default:
		return false
	}
}

func classifyParsedStatement(stmt parser.Statement) QueryType {
	ps, ok := stmt.(*parser.ParsedStmt)
	if !ok || ps == nil {
		return analysis.QueryTypeUnknown
	}
	switch ps.Kind {
	case parser.StmtQuery:
		return analysis.Select
	case parser.StmtExplain:
		return analysis.Explain
	case parser.StmtShow, parser.StmtDescribe:
		return analysis.SelectInfoSchema
	case parser.StmtInsert, parser.StmtMultiInsert, parser.StmtUpdate,
		parser.StmtDelete, parser.StmtMerge, parser.StmtLoadData:
		return analysis.DML
	case parser.StmtCreateNamespace, parser.StmtAlterNamespace, parser.StmtDropNamespace,
		parser.StmtCreateTable, parser.StmtCreateTableLike, parser.StmtReplaceTable,
		parser.StmtAlterTable, parser.StmtDropTable,
		parser.StmtCreateView, parser.StmtAlterView, parser.StmtDropView,
		parser.StmtCreateFunction, parser.StmtDropFunction,
		parser.StmtCreateIndex, parser.StmtDropIndex,
		parser.StmtTruncate, parser.StmtRepair, parser.StmtComment,
		parser.StmtAnalyze, parser.StmtRefresh,
		parser.StmtCache, parser.StmtUncache, parser.StmtClearCache:
		return analysis.DDL
	default:
		return analysis.QueryTypeUnknown
	}
}
