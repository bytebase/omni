// Package analysis provides query analysis for PartiQL statements.
// It mirrors bytebase's query validation requirements: rejecting
// DDL/DML/EXEC and extracting basic query structure from SELECT.
package analysis

import (
	"fmt"

	"github.com/bytebase/omni/partiql/ast"
	"github.com/bytebase/omni/partiql/parser"
)

// ValidateQuery parses the input and validates it is a DQL-only query
// (SELECT or set operation). Returns an error if the input contains DDL,
// DML, EXEC, or EXPLAIN wrapping non-DQL.
//
// This mirrors the legacy bytebase validateQuery function for PartiQL.
func ValidateQuery(input string) error {
	list, err := parser.Parse(input)
	if err != nil {
		return err
	}
	for _, item := range list.Items {
		if err := validateStatement(item); err != nil {
			return err
		}
	}
	return nil
}

func validateStatement(node ast.Node) error {
	switch node.(type) {
	case *ast.SelectStmt:
		return nil // DQL — allowed
	case *ast.SetOpStmt:
		return nil // DQL (UNION/INTERSECT/EXCEPT) — allowed
	case *ast.ExplainStmt:
		// EXPLAIN is always allowed regardless of the inner statement,
		// matching the legacy ANTLR validator behavior. EXPLAIN is a
		// read-only operation even when wrapping DML/DDL.
		return nil
	case *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt,
		*ast.UpsertStmt, *ast.ReplaceStmt, *ast.RemoveStmt:
		return fmt.Errorf("DML statements are not allowed in query context")
	case *ast.CreateTableStmt, *ast.CreateIndexStmt,
		*ast.DropTableStmt, *ast.DropIndexStmt:
		return fmt.Errorf("DDL statements are not allowed in query context")
	case *ast.ExecStmt:
		return fmt.Errorf("EXEC statements are not allowed in query context")
	default:
		return fmt.Errorf("unexpected statement type %T", node)
	}
}

// QueryAnalysis is the result of analyzing a parsed PartiQL SELECT statement.
type QueryAnalysis struct {
	Projections []Projection
	SelectStar  bool
	Tables      []string // table names from FROM clause
}

// Projection represents one item in the SELECT list.
type Projection struct {
	Name string // output alias; empty if none
	Expr string // string representation of the expression
}

// Analyze extracts basic structure from a PartiQL SELECT statement.
func Analyze(stmt *ast.SelectStmt) *QueryAnalysis {
	qa := &QueryAnalysis{}
	switch {
	case stmt.Star:
		qa.SelectStar = true
	case stmt.Value != nil:
		// SELECT VALUE expr
		qa.Projections = append(qa.Projections, Projection{
			Expr: ast.NodeToString(stmt.Value),
		})
	case stmt.Pivot != nil:
		// SELECT PIVOT v AT k — represent as a single unnamed projection
		qa.Projections = append(qa.Projections, Projection{
			Expr: ast.NodeToString(stmt.Pivot),
		})
	default:
		for _, te := range stmt.Targets {
			p := Projection{Expr: ast.NodeToString(te.Expr)}
			if te.Alias != nil {
				p.Name = *te.Alias
			}
			qa.Projections = append(qa.Projections, p)
		}
	}
	qa.Tables = extractTables(stmt.From)
	return qa
}

// extractTables walks a TableExpr and collects bare table/relation names.
func extractTables(te ast.TableExpr) []string {
	if te == nil {
		return nil
	}
	switch v := te.(type) {
	case *ast.TableRef:
		return []string{v.Name}
	case *ast.AliasedSource:
		return extractTables(v.Source)
	case *ast.VarRef:
		return []string{v.Name}
	case *ast.PathExpr:
		if vr, ok := v.Root.(*ast.VarRef); ok {
			return []string{vr.Name}
		}
	case *ast.JoinExpr:
		left := extractTables(v.Left)
		right := extractTables(v.Right)
		return append(left, right...)
	case *ast.SubLink:
		// subquery in FROM — no simple table name
		return nil
	case *ast.UnpivotExpr:
		return nil
	}
	return nil
}
