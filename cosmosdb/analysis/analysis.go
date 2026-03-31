package analysis

import (
	"strconv"

	"github.com/bytebase/omni/cosmosdb/ast"
)

// Projection represents one item in the SELECT list.
type Projection struct {
	Name        string      // output alias; empty if none and not inferrable
	SourcePaths []FieldPath // source field paths this projection references
}

// QueryAnalysis is the result of analyzing a parsed CosmosDB SELECT statement.
type QueryAnalysis struct {
	Projections []Projection // SELECT list items (empty when SelectStar is true)
	SelectStar  bool         // true if SELECT *
	Predicates  []FieldPath  // field paths referenced in WHERE
}

// Analyze extracts field-level information from a parsed CosmosDB SELECT statement.
func Analyze(stmt *ast.SelectStmt) *QueryAnalysis {
	aliases := buildAliasMap(stmt)
	qa := &QueryAnalysis{}

	if stmt.Star {
		qa.SelectStar = true
	} else {
		qa.Projections = extractProjections(stmt.Targets, aliases)
	}

	if stmt.Where != nil {
		qa.Predicates = extractFieldPaths(stmt.Where, aliases)
	}

	return qa
}

// aliasMap maps alias names to their resolved source field paths.
type aliasMap map[string]FieldPath

// buildAliasMap walks FROM and JOINs to build the alias resolution map.
func buildAliasMap(stmt *ast.SelectStmt) aliasMap {
	m := make(aliasMap)
	if stmt.From == nil {
		return m
	}
	registerFromSource(stmt.From, m)
	for _, join := range stmt.Joins {
		registerFromSource(join.Source, m)
	}
	return m
}

// registerFromSource registers aliases from a single FROM/JOIN source.
func registerFromSource(te ast.TableExpr, m aliasMap) {
	switch src := te.(type) {
	case *ast.AliasedTableExpr:
		containerName := containerRefName(src.Source)
		if containerName != "" && src.Alias != "" {
			m[src.Alias] = FieldPath{ItemSelector(containerName)}
		}
	case *ast.ContainerRef:
		if src.Name != "" {
			m[src.Name] = FieldPath{ItemSelector(src.Name)}
		}
	case *ast.ArrayIterationExpr:
		sourcePath := resolveTableExprPath(src.Source, m)
		if src.Alias != "" && len(sourcePath) > 0 {
			m[src.Alias] = sourcePath
		}
	}
}

// containerRefName extracts the container name from a TableExpr.
func containerRefName(te ast.TableExpr) string {
	switch src := te.(type) {
	case *ast.ContainerRef:
		return src.Name
	case *ast.DotAccessExpr:
		return exprRootName(src)
	default:
		return ""
	}
}

// resolveTableExprPath resolves a TableExpr to a FieldPath using the alias map.
func resolveTableExprPath(te ast.TableExpr, m aliasMap) FieldPath {
	switch src := te.(type) {
	case *ast.ContainerRef:
		if resolved, ok := m[src.Name]; ok {
			return copyPath(resolved)
		}
		return FieldPath{ItemSelector(src.Name)}
	case *ast.DotAccessExpr:
		base := resolveExprPath(src.Expr, m)
		return append(base, ItemSelector(src.Property))
	case *ast.BracketAccessExpr:
		base := resolveExprPath(src.Expr, m)
		return append(base, bracketSelector(src.Index))
	default:
		return nil
	}
}

// resolveExprPath resolves an ExprNode to a FieldPath using the alias map.
func resolveExprPath(expr ast.ExprNode, m aliasMap) FieldPath {
	switch e := expr.(type) {
	case *ast.ColumnRef:
		if resolved, ok := m[e.Name]; ok {
			return copyPath(resolved)
		}
		return FieldPath{ItemSelector(e.Name)}
	case *ast.DotAccessExpr:
		base := resolveExprPath(e.Expr, m)
		return append(base, ItemSelector(e.Property))
	case *ast.BracketAccessExpr:
		base := resolveExprPath(e.Expr, m)
		return append(base, bracketSelector(e.Index))
	default:
		return nil
	}
}

// exprRootName extracts the root identifier name from a dot/bracket chain.
func exprRootName(expr ast.ExprNode) string {
	switch e := expr.(type) {
	case *ast.ColumnRef:
		return e.Name
	case *ast.DotAccessExpr:
		return exprRootName(e.Expr)
	case *ast.BracketAccessExpr:
		return exprRootName(e.Expr)
	default:
		return ""
	}
}

// bracketSelector creates a Selector from a bracket index expression.
func bracketSelector(index ast.ExprNode) Selector {
	switch idx := index.(type) {
	case *ast.NumberLit:
		if n, err := strconv.Atoi(idx.Val); err == nil {
			return ArraySelector(idx.Val, n)
		}
		return ItemSelector(idx.Val)
	case *ast.StringLit:
		return ItemSelector(idx.Val)
	default:
		return ItemSelector("?")
	}
}

// extractProjections builds the projection list from SELECT targets.
func extractProjections(targets []*ast.TargetEntry, aliases aliasMap) []Projection {
	var projections []Projection
	for _, t := range targets {
		paths := extractFieldPaths(t.Expr, aliases)
		name := ""
		if t.Alias != nil {
			name = *t.Alias
		} else if len(paths) == 1 && len(paths[0]) > 0 {
			last := paths[0][len(paths[0])-1]
			if !last.IsArray() {
				name = last.Name
			}
		}
		projections = append(projections, Projection{
			Name:        name,
			SourcePaths: paths,
		})
	}
	return projections
}

// copyPath returns a shallow copy of a FieldPath to avoid alias mutation.
func copyPath(fp FieldPath) FieldPath {
	cp := make(FieldPath, len(fp))
	copy(cp, fp)
	return cp
}
