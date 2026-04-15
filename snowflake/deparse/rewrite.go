// Package deparse — rewrite.go implements AST-level rewrites applied before
// deparsing. InjectLimit caps result-set size, used by the Bytebase read-only
// SQL editor before sending a query to Snowflake.
package deparse

import (
	"fmt"

	"github.com/bytebase/omni/snowflake/ast"
	"github.com/bytebase/omni/snowflake/parser"
)

// InjectLimit rewrites the given SQL to cap the result-set size to maxRows.
//
// The semantics depend on what the SQL already has:
//   - No LIMIT, no OFFSET, no FETCH → add LIMIT maxRows
//   - Existing LIMIT n with n <= maxRows → leave unchanged
//   - Existing LIMIT n with n > maxRows → replace with LIMIT maxRows
//   - Existing LIMIT n with non-literal n (expr/variable) → wrap: SELECT * FROM (...) LIMIT maxRows
//   - Existing FETCH FIRST n ROWS ONLY → same literal vs. expr handling as LIMIT
//   - SetOperationStmt (UNION/INTERSECT/EXCEPT) → wrap the whole thing: SELECT * FROM (...) LIMIT maxRows
//   - Non-SELECT statements (INSERT/UPDATE/DELETE/MERGE/DDL) → return the input unchanged
//   - Multi-statement input → apply to each statement independently
//
// Returns the rewritten SQL and any error encountered during parsing/deparsing.
func InjectLimit(sql string, maxRows int) (string, error) {
	if maxRows <= 0 {
		return "", fmt.Errorf("maxRows must be positive")
	}
	file, err := parser.Parse(sql)
	if err != nil {
		return "", err
	}
	for i, stmt := range file.Stmts {
		file.Stmts[i] = rewriteStmt(stmt, maxRows)
	}
	return DeparseFile(file)
}

// rewriteStmt dispatches on the concrete node type. Non-query statements are
// returned unchanged.
func rewriteStmt(node ast.Node, maxRows int) ast.Node {
	switch s := node.(type) {
	case *ast.SelectStmt:
		return rewriteSelect(s, maxRows)
	case *ast.SetOperationStmt:
		// Set operations are always wrapped because the LIMIT belongs to the
		// outermost query, not to either branch.
		return wrapInLimit(s, maxRows)
	default:
		return node
	}
}

// rewriteSelect inspects an existing SelectStmt and either adds, lowers, or
// wraps the LIMIT/FETCH clause so that at most maxRows rows are returned.
func rewriteSelect(s *ast.SelectStmt, maxRows int) ast.Node {
	limitLit := asIntLiteral(s.Limit)
	fetchLit := asFetchIntLiteral(s.Fetch)

	switch {
	// -----------------------------------------------------------------------
	// FETCH FIRST n ROWS ONLY is present
	// -----------------------------------------------------------------------
	case s.Fetch != nil:
		if fetchLit == nil {
			// Non-literal FETCH expression — wrap to be safe.
			return wrapInLimit(s, maxRows)
		}
		if fetchLit.Ival > int64(maxRows) {
			// Lower the FETCH count in place.
			s.Fetch = &ast.FetchClause{
				Count: intLiteral(maxRows),
			}
		}
		// fetchLit.Ival <= maxRows: leave unchanged.
		return s

	// -----------------------------------------------------------------------
	// LIMIT is present
	// -----------------------------------------------------------------------
	case s.Limit != nil:
		if limitLit == nil {
			// Non-literal LIMIT expression (bind variable, subquery, etc.) — wrap.
			return wrapInLimit(s, maxRows)
		}
		if limitLit.Ival > int64(maxRows) {
			s.Limit = intLiteral(maxRows)
		}
		// limitLit.Ival <= maxRows: leave unchanged.
		return s

	// -----------------------------------------------------------------------
	// No LIMIT and no FETCH — add LIMIT maxRows
	// -----------------------------------------------------------------------
	default:
		s.Limit = intLiteral(maxRows)
		return s
	}
}

// wrapInLimit builds:
//
//	SELECT * FROM (<inner>) AS _q LIMIT maxRows
func wrapInLimit(inner ast.Node, maxRows int) *ast.SelectStmt {
	alias := ast.Ident{Name: "_q"}
	return &ast.SelectStmt{
		Targets: []*ast.SelectTarget{
			{Star: true, Expr: &ast.StarExpr{}},
		},
		From: []ast.Node{
			&ast.TableRef{
				Subquery: inner,
				Alias:    alias,
			},
		},
		Limit: intLiteral(maxRows),
	}
}

// intLiteral creates an integer Literal node with the given value.
func intLiteral(v int) *ast.Literal {
	return &ast.Literal{
		Kind:  ast.LitInt,
		Value: fmt.Sprintf("%d", v),
		Ival:  int64(v),
	}
}

// asIntLiteral returns the *ast.Literal if node is an integer literal,
// otherwise nil.
func asIntLiteral(node ast.Node) *ast.Literal {
	if node == nil {
		return nil
	}
	lit, ok := node.(*ast.Literal)
	if !ok || lit.Kind != ast.LitInt {
		return nil
	}
	return lit
}

// asFetchIntLiteral is like asIntLiteral but extracts from a *FetchClause.
func asFetchIntLiteral(f *ast.FetchClause) *ast.Literal {
	if f == nil {
		return nil
	}
	return asIntLiteral(f.Count)
}
