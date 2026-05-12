package catalog

import (
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
	pgparser "github.com/bytebase/omni/pg/parser"
)

// ExecOptions controls execution behavior.
type ExecOptions struct {
	// ContinueOnError, when true, continues executing subsequent statements
	// after a failure, collecting all errors. When false (default), execution
	// stops at the first error.
	ContinueOnError bool
}

// ExecResult is the execution result for a single statement.
type ExecResult struct {
	Index    int       // statement position in the batch (0-based)
	SQL      string    // original SQL text for this statement
	Line     int       // 1-based start line in the original SQL
	Skipped  bool      // true if the statement was not processed (DML, etc.)
	Error    error     // nil = success or skipped
	Warnings []Warning // non-fatal notices emitted during execution
}

// Exec parses and executes one or more SQL statements against the catalog.
// DDL statements modify catalog state. DML, transaction control, and other
// non-utility statements are skipped (Skipped=true in the result).
//
// pg: src/backend/tcop/postgres.c — exec_simple_query
func (c *Catalog) Exec(sql string, opts *ExecOptions) ([]ExecResult, error) {
	list, err := pgparser.Parse(sql)
	if err != nil {
		return nil, err
	}
	if list == nil || len(list.Items) == 0 {
		return nil, nil
	}

	// Compute statement boundaries (byte offsets) by scanning for semicolons.
	boundaries := scanStatementBoundaries(sql, len(list.Items))

	// Build line number index for offset→line conversion.
	lineIndex := buildLineIndex(sql)

	// Schema dumps are often not strictly executable order. Pre-declare
	// CREATE OR REPLACE functions when possible so earlier views can resolve
	// RETURNS TABLE/OUT columns during analysis. Errors are still surfaced
	// during normal statement execution below.
	predeclaredFunctions := make(map[int]bool)
	firstFunctionIndex := firstFunctionDefinitionIndexes(list)

	continueOnError := false
	if opts != nil {
		continueOnError = opts.ContinueOnError
	}

	results := make([]ExecResult, 0, len(list.Items))
	for i, item := range list.Items {
		// Drain any leftover warnings from previous operations.
		c.DrainWarnings()

		c.predeclareReplaceableFunctions(list, predeclaredFunctions, firstFunctionIndex)

		// Unwrap RawStmt if present.
		stmt := unwrapRawStmt(item)

		// Determine SQL text and line number for this statement.
		stmtSQL := ""
		stmtLine := 1
		if i < len(boundaries) {
			b := boundaries[i]
			stmtSQL = b.sql
			stmtLine = offsetToLine(lineIndex, b.start)
		}

		result := ExecResult{
			Index: i,
			SQL:   stmtSQL,
			Line:  stmtLine,
		}

		// Route: DML and non-utility statements are skipped.
		if isDML(stmt) {
			result.Skipped = true
			results = append(results, result)
			continue
		}

		// Execute DDL/utility via ProcessUtility.
		execErr := c.ProcessUtility(stmt)
		result.Error = execErr
		result.Warnings = c.DrainWarnings()
		results = append(results, result)

		if execErr != nil && !continueOnError {
			break
		}
	}
	return results, nil
}

func firstFunctionDefinitionIndexes(list *nodes.List) map[string]int {
	first := make(map[string]int)
	for i, item := range list.Items {
		fn, ok := unwrapRawStmt(item).(*nodes.CreateFunctionStmt)
		if !ok {
			continue
		}
		identity := functionIdentity(fn)
		if _, exists := first[identity]; !exists {
			first[identity] = i
		}
	}
	return first
}

func (c *Catalog) predeclareReplaceableFunctions(list *nodes.List, predeclared map[int]bool, firstFunctionIndex map[string]int) {
	for i, item := range list.Items {
		if predeclared[i] {
			continue
		}
		stmt := unwrapRawStmt(item)
		fn, ok := stmt.(*nodes.CreateFunctionStmt)
		if !ok || !fn.IsOrReplace {
			continue
		}
		if firstFunctionIndex[functionIdentity(fn)] != i {
			continue
		}
		if err := c.CreateFunctionStmt(fn); err == nil {
			predeclared[i] = true
		}
		c.DrainWarnings()
	}
}

func unwrapRawStmt(item nodes.Node) nodes.Node {
	if raw, ok := item.(*nodes.RawStmt); ok {
		return raw.Stmt
	}
	return item
}

// isDML returns true for statements that are not utility statements
// and should be skipped by Exec (they go through the planner/executor
// path in PostgreSQL, not ProcessUtility).
func isDML(stmt nodes.Node) bool {
	switch stmt.(type) {
	case *nodes.SelectStmt:
		return true
	case *nodes.InsertStmt:
		return true
	case *nodes.UpdateStmt:
		return true
	case *nodes.DeleteStmt:
		return true
	case *nodes.MergeStmt:
		return true
	default:
		return false
	}
}

// stmtBoundary records the byte range and SQL text for a single statement.
type stmtBoundary struct {
	start int    // byte offset of statement start
	sql   string // trimmed SQL text
}

// scanStatementBoundaries uses the pgparser Lexer to find statement boundaries
// (semicolons) in the SQL text, respecting string literals and comments.
func scanStatementBoundaries(sql string, expectedCount int) []stmtBoundary {
	lex := pgparser.NewLexer(sql)
	lex.StandardConformingStrings = true

	boundaries := make([]stmtBoundary, 0, expectedCount)
	stmtStart := 0
	firstTokenLoc := -1

	for {
		tok := lex.NextToken()
		if tok.Type == 0 { // EOF
			break
		}

		if firstTokenLoc < 0 {
			firstTokenLoc = tok.Loc
		}

		if tok.Type == ';' {
			// End of statement: extract SQL from stmtStart to semicolon.
			end := tok.Loc
			startOff := stmtStart
			if firstTokenLoc >= 0 {
				startOff = firstTokenLoc
			}
			stmtText := strings.TrimSpace(sql[startOff:end])
			if stmtText != "" {
				boundaries = append(boundaries, stmtBoundary{
					start: startOff,
					sql:   stmtText,
				})
			}
			stmtStart = end + 1
			firstTokenLoc = -1
		}
	}

	// Handle last statement without trailing semicolon.
	if stmtStart < len(sql) {
		startOff := stmtStart
		if firstTokenLoc >= 0 {
			startOff = firstTokenLoc
		}
		remaining := strings.TrimSpace(sql[startOff:])
		if remaining != "" {
			boundaries = append(boundaries, stmtBoundary{
				start: startOff,
				sql:   remaining,
			})
		}
	}

	return boundaries
}

// buildLineIndex returns the byte offset of each line start (0-indexed line numbers).
func buildLineIndex(sql string) []int {
	index := []int{0} // line 0 starts at offset 0
	for i, ch := range sql {
		if ch == '\n' {
			index = append(index, i+1)
		}
	}
	return index
}

// offsetToLine converts a byte offset to a 1-based line number.
func offsetToLine(lineIndex []int, offset int) int {
	// Binary search for the last line start <= offset.
	lo, hi := 0, len(lineIndex)-1
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if lineIndex[mid] <= offset {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return lo + 1 // 1-based
}
