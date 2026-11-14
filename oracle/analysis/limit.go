package analysis

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/oracle/parser"
)

const dbVersion12 = 12

// RewriteResultLimit limits the rows a query returns to limit. engineVersion
// picks the clause: ROWNUM before version 12 or when it does not parse, FETCH
// NEXT from 12c on. A statement that is not a query, or a simple SELECT ...
// FROM DUAL, is returned unchanged. An error means the FETCH NEXT rewrite
// failed; RowNumResultLimit is the fallback.
func RewriteResultLimit(statement string, limit int, engineVersion string) (string, error) {
	if shouldSkipLimit(statement) {
		return statement, nil
	}
	if isOracle11gOrEarlier(engineVersion) {
		return addLimitFor11g(statement, limit), nil
	}
	if !isSelectOrWithStatement(statement) {
		return statement, nil
	}
	return addFetchNextClause(statement, limit)
}

// RowNumResultLimit wraps a query in a SELECT that keeps its first limit rows
// by ROWNUM. A statement that is not a query is returned unchanged.
func RowNumResultLimit(statement string, limit int) string {
	return addLimitFor11g(statement, limit)
}

func trimStatement(statement string) string {
	return strings.TrimLeftFunc(strings.TrimRightFunc(statement, isSpaceOrSemicolon), unicode.IsSpace)
}

func isSpaceOrSemicolon(r rune) bool {
	return unicode.IsSpace(r) || r == ';'
}

// shouldSkipLimit checks if the statement needs a limit clause
func shouldSkipLimit(stmt string) bool {
	ok, err := skipAddLimit(stmt)
	return err == nil && ok
}

// isOracle11gOrEarlier checks if the Oracle version is 11g or earlier
func isOracle11gOrEarlier(engineVersion string) bool {
	versionIdx := strings.Index(engineVersion, ".")
	if versionIdx < 0 {
		return true // Default to 11g behavior for invalid version
	}
	versionNumber, err := strconv.Atoi(engineVersion[:versionIdx])
	if err != nil {
		return true // Default to 11g behavior for parsing errors
	}
	return versionNumber < dbVersion12
}

// addLimitFor11g adds a ROWNUM-based limit for Oracle 11g and earlier versions.
// Uses the legacy approach with subquery and ROWNUM.
func addLimitFor11g(statement string, limitCount int) string {
	if !isSelectOrWithStatement(statement) {
		return statement
	}
	return fmt.Sprintf("SELECT * FROM (%s) WHERE ROWNUM <= %d", trimStatement(statement), limitCount)
}

// isSelectOrWithStatement checks if the statement is a query: SELECT, WITH, or
// a parenthesized query.
func isSelectOrWithStatement(statement string) bool {
	trimmedStatement := strings.ToLower(strings.TrimLeftFunc(statement, unicode.IsSpace))
	return strings.HasPrefix(trimmedStatement, "select") || strings.HasPrefix(trimmedStatement, "with") || strings.HasPrefix(trimmedStatement, "(")
}

// addFetchNextClause adds a FETCH NEXT clause to a SELECT statement using AST parsing.
// This provides more precise placement of the limit clause compared to simple string wrapping.
func addFetchNextClause(statement string, limitCount int) (string, error) {
	list, err := Parse(statement)
	if err != nil {
		return "", err
	}
	if list == nil || len(list.Items) == 0 {
		return "", errors.New("no parse results")
	}
	if len(list.Items) > 1 {
		return "", fmt.Errorf("expected single statement, got %d statements", len(list.Items))
	}
	raw, ok := list.Items[0].(*ast.RawStmt)
	if !ok {
		return "", fmt.Errorf("expected raw statement, got %T", list.Items[0])
	}
	selectStmt, ok := raw.Stmt.(*ast.SelectStmt)
	if !ok {
		return statement, nil
	}

	res, err := rewriteOracleSelectFetch(statement, selectStmt, limitCount)
	if err != nil {
		return "", err
	}
	// https://stackoverflow.com/questions/27987882/how-can-i-solve-ora-00911-invalid-character-error
	res = strings.TrimRightFunc(res, isSpaceOrSemicolon)

	return res, nil
}

func rewriteOracleSelectFetch(sql string, selectStmt *ast.SelectStmt, limitCount int) (string, error) {
	target, parenthesized := rightmostOracleSetSelect(sql, selectStmt)
	if !parenthesized {
		if target.FetchFirst != nil {
			return rewriteOracleFetchClause(sql, target.FetchFirst, limitCount)
		}
		if target.ForUpdate != nil && target.ForUpdate.Loc.Start > 0 && target.ForUpdate.Loc.Start <= len(sql) {
			return sql[:target.ForUpdate.Loc.Start] + fmt.Sprintf("FETCH NEXT %d ROWS ONLY ", limitCount) + sql[target.ForUpdate.Loc.Start:], nil
		}
	}

	loc := ast.NodeLoc(target)
	if loc.End < 0 || loc.End > len(sql) {
		return "", fmt.Errorf("invalid SELECT end position %d", loc.End)
	}
	return sql[:loc.End] + fmt.Sprintf(" FETCH NEXT %d ROWS ONLY", limitCount) + sql[loc.End:], nil
}

// rightmostOracleSetSelect returns the operand that omni attaches the ORDER
// BY, row limit, and FOR UPDATE ending a set operation to: its last operand.
// omni chains operands through Rarg, and only the text tells a parenthesized
// operand from the rest of the chain. The walk stops at a parenthesized one
// and reports it, as the clauses inside it limit that operand alone.
func rightmostOracleSetSelect(sql string, selectStmt *ast.SelectStmt) (*ast.SelectStmt, bool) {
	for selectStmt.Op != ast.SETOP_NONE && selectStmt.Rarg != nil {
		selectStmt = selectStmt.Rarg
		if isParenthesizedOracleQuery(sql, ast.NodeLoc(selectStmt)) {
			return selectStmt, true
		}
	}
	return selectStmt, false
}

// isParenthesizedOracleQuery reports whether one pair of parentheses encloses
// all of sql[loc.Start:loc.End].
func isParenthesizedOracleQuery(sql string, loc ast.Loc) bool {
	if loc.Start < 0 || loc.End > len(sql) || loc.Start >= loc.End || sql[loc.Start] != '(' {
		return false
	}
	lexer := parser.NewLexerRange(sql, loc.Start, loc.End)
	depth := 0
	for {
		tok := lexer.NextToken()
		if tok.End <= tok.Loc {
			return false
		}
		switch tok.Type {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return tok.End == loc.End
			}
		default:
		}
	}
}

func rewriteOracleFetchClause(sql string, fetch *ast.FetchFirstClause, limitCount int) (string, error) {
	if fetch.Count == nil {
		if fetch.Loc.Start < 0 || fetch.Loc.End < fetch.Loc.Start || fetch.Loc.End > len(sql) {
			return "", fmt.Errorf("invalid FETCH position %d:%d", fetch.Loc.Start, fetch.Loc.End)
		}
		if hasOracleFetchKeyword(sql, fetch.Loc) {
			return sql, nil
		}
		return sql[:fetch.Loc.End] + fmt.Sprintf(" FETCH NEXT %d ROWS ONLY", limitCount) + sql[fetch.Loc.End:], nil
	}
	if fetch.Percent {
		return "", fmt.Errorf("cannot rewrite PERCENT FETCH expression")
	}

	existingLimit := extractOracleFetchCount(fetch.Count)
	if existingLimit > 0 && existingLimit <= limitCount {
		return sql, nil
	}

	loc := ast.NodeLoc(fetch.Count)
	loc = trimOracleLocSpace(sql, loc)
	if loc.Start >= 0 && loc.End > loc.Start && loc.End <= len(sql) {
		if existingLimit <= 0 {
			return "", fmt.Errorf("cannot rewrite non-constant FETCH expression")
		}
		return sql[:loc.Start] + fmt.Sprintf("%d", limitCount) + sql[loc.End:], nil
	}
	return "", fmt.Errorf("cannot rewrite FETCH expression")
}

func hasOracleFetchKeyword(sql string, loc ast.Loc) bool {
	segment := sql[loc.Start:loc.End]
	lexer := parser.NewLexer(segment)
	for {
		tok := lexer.NextToken()
		if tok.Loc == tok.End && tok.End >= len(segment) {
			return false
		}
		if tok.Loc < 0 || tok.End > len(segment) || tok.End <= tok.Loc {
			return false
		}
		if strings.EqualFold(segment[tok.Loc:tok.End], "FETCH") {
			return true
		}
	}
}

func extractOracleFetchCount(node ast.Node) int {
	limit, ok := node.(*ast.NumberLiteral)
	if !ok || limit.IsFloat || limit.Ival <= 0 {
		return 0
	}
	return int(limit.Ival)
}

func trimOracleLocSpace(sql string, loc ast.Loc) ast.Loc {
	for loc.Start < loc.End && loc.Start < len(sql) && unicode.IsSpace(rune(sql[loc.Start])) {
		loc.Start++
	}
	for loc.End > loc.Start && loc.End <= len(sql) && unicode.IsSpace(rune(sql[loc.End-1])) {
		loc.End--
	}
	return loc
}

// skipAddLimit checks if the statement needs a limit clause.
// For Oracle, we think the statement like "SELECT xxx FROM DUAL" does not need a limit clause.
// More details, xxx can not be a subquery.
func skipAddLimit(stmt string) (bool, error) {
	list, err := Parse(stmt)
	if err != nil {
		return false, err
	}
	if list == nil || len(list.Items) == 0 {
		return false, nil
	}
	// Multiple statements should not skip limit
	if len(list.Items) > 1 {
		return false, nil
	}
	raw, ok := list.Items[0].(*ast.RawStmt)
	if !ok {
		return false, nil
	}
	selectStmt, ok := raw.Stmt.(*ast.SelectStmt)
	if !ok {
		return false, nil
	}
	if !isSimpleOracleSelect(selectStmt) {
		return false, nil
	}
	if !isOracleSelectFromDual(selectStmt) {
		return false, nil
	}
	return !hasOracleSubqueriesInSelection(selectStmt), nil
}

func isSimpleOracleSelect(selectStmt *ast.SelectStmt) bool {
	return selectStmt.WithClause == nil &&
		!selectStmt.Distinct &&
		!selectStmt.UniqueKw &&
		!selectStmt.All &&
		selectStmt.Into == nil &&
		selectStmt.IntoVars == nil &&
		selectStmt.WhereClause == nil &&
		selectStmt.Hierarchical == nil &&
		selectStmt.GroupClause == nil &&
		selectStmt.HavingClause == nil &&
		selectStmt.ModelClause == nil &&
		len(selectStmt.WindowDefs) == 0 &&
		selectStmt.QualifyClause == nil &&
		selectStmt.OrderBy == nil &&
		selectStmt.ForUpdate == nil &&
		selectStmt.FetchFirst == nil &&
		selectStmt.Op == ast.SETOP_NONE
}

func isOracleSelectFromDual(selectStmt *ast.SelectStmt) bool {
	if selectStmt.FromClause == nil || len(selectStmt.FromClause.Items) != 1 {
		return false
	}
	tableRef, ok := selectStmt.FromClause.Items[0].(*ast.TableRef)
	if !ok || tableRef.Name == nil {
		return false
	}
	return tableRef.Alias == nil &&
		tableRef.Name.Schema == "" &&
		tableRef.Name.DBLink == "" &&
		strings.EqualFold(tableRef.Name.Name, "DUAL")
}

func hasOracleSubqueriesInSelection(selectStmt *ast.SelectStmt) bool {
	if selectStmt.TargetList == nil || len(selectStmt.TargetList.Items) == 0 {
		return true
	}
	for _, item := range selectStmt.TargetList.Items {
		target, ok := item.(*ast.ResTarget)
		if !ok || target.Expr == nil {
			return true
		}
		if _, ok := target.Expr.(*ast.Star); ok {
			return true
		}
		hasSubquery := false
		ast.Inspect(target.Expr, func(node ast.Node) bool {
			if _, ok := node.(*ast.SubqueryExpr); ok {
				hasSubquery = true
				return false
			}
			return true
		})
		if hasSubquery {
			return true
		}
	}
	return false
}
