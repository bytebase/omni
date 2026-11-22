package analysis

import (
	"errors"
	"strings"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/hive/parser"
)

// QueryType is the permission class of a statement.
type QueryType = analysis.QueryType

// ValidateQuery reports whether every statement in script is a read-only HiveQL
// query that returns rows.
//
// Classification uses hive/parser's Hive 4.0.1 grammar rather than keyword
// prefixes: CTE-prefixed or FROM-first inserts (`WITH ... INSERT ...`,
// `FROM ... INSERT ...`), external script transforms (`MAP`/`REDUCE`/
// `TRANSFORM`), and `EXPLAIN ANALYZE` (which executes what it plans) are all
// rejected. Read-only metadata statements (`SHOW`, `DESCRIBE`, `DESC`) and
// non-executing `EXPLAIN` statements are allowed. Lexical or syntax errors are
// returned with script-relative byte offsets.
func ValidateQuery(script string) (readOnly bool, returnsData bool, err error) {
	if err := firstLexError(script); err != nil {
		return false, false, err
	}
	for _, seg := range parser.Split(script) {
		if seg.Empty() {
			continue
		}
		trimmed := strings.TrimRightFunc(seg.Text, isSpaceOrSemicolon)
		res, err := parser.Parse(trimmed)
		if err != nil {
			return false, false, offsetParseError(err, seg.ByteStart)
		}
		if !isReadOnlyStatement(res.Statement) {
			return false, false, nil
		}
	}
	return true, true, nil
}

// Classify parses a single HiveQL statement and returns its QueryType.
// Empty or comment-only input returns QueryTypeUnknown with a nil error;
// multi-statement input or invalid syntax returns an error.
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
		return analysis.QueryTypeUnknown, errors.New("expected a single HiveQL statement")
	}
	trimmed := strings.TrimRightFunc(nonEmpty[0].Text, isSpaceOrSemicolon)
	res, err := parser.Parse(trimmed)
	if err != nil {
		return analysis.QueryTypeUnknown, offsetParseError(err, nonEmpty[0].ByteStart)
	}
	return classifyParsedStatement(res.Statement), nil
}

func isReadOnlyStatement(stmt *parser.Statement) bool {
	if stmt == nil {
		return false
	}
	switch stmt.Kind {
	case parser.StmtQuery:
		return !stmt.HasTransform
	case parser.StmtExplain:
		return stmt.Explain != nil && !stmt.Explain.Analyze
	case parser.StmtDDL:
		return stmt.FirstKeyword == parser.KW_SHOW ||
			stmt.FirstKeyword == parser.KW_DESCRIBE ||
			stmt.FirstKeyword == parser.KW_DESC
	default:
		return false
	}
}

func classifyParsedStatement(stmt *parser.Statement) QueryType {
	if stmt == nil {
		return analysis.QueryTypeUnknown
	}
	switch stmt.Kind {
	case parser.StmtQuery:
		return analysis.Select
	case parser.StmtExplain:
		return analysis.Explain
	case parser.StmtDDL:
		if stmt.FirstKeyword == parser.KW_SHOW ||
			stmt.FirstKeyword == parser.KW_DESCRIBE ||
			stmt.FirstKeyword == parser.KW_DESC {
			return analysis.SelectInfoSchema
		}
		return analysis.DDL
	case parser.StmtInsert, parser.StmtUpdate, parser.StmtDelete,
		parser.StmtMerge, parser.StmtLoad, parser.StmtExport, parser.StmtImport:
		return analysis.DML
	default:
		return analysis.QueryTypeUnknown
	}
}

func firstLexError(text string) error {
	l := parser.NewLexer(text)
	for l.NextToken().Type != parser.TokEOF {
		if err := l.Err(); err != nil {
			return err
		}
	}
	return l.Err()
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
