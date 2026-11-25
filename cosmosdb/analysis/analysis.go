package analysis

import (
	"strconv"

	"github.com/bytebase/omni/cosmosdb/ast"
)

// Projection represents one item in the SELECT list.
type Projection struct {
	// Name is the property the item is returned under, as Cosmos DB names
	// it: the alias; else the identifier or last property name of a property
	// path (c.name, c["name"], t); else $1, $2, ... counted over the unnamed
	// items in order. A SELECT VALUE item is returned bare and has no name.
	Name        string
	SourcePaths []FieldPath // source field paths this projection references
}

// QueryAnalysis is the result of analyzing a parsed CosmosDB SELECT statement.
type QueryAnalysis struct {
	// Projections is the SELECT list. It is empty when SelectStar is true. A
	// SELECT * over a FROM subquery that names its items is returned as that
	// subquery's items, which is what Cosmos DB returns, so it is not a star.
	Projections []Projection
	SelectStar  bool // true if SELECT *
	// StarPaths is, when SelectStar is true, where each returned document
	// comes from: the container's root for SELECT * FROM c, the property for
	// SELECT * FROM c.address, the array iterated for SELECT * FROM t IN
	// c.tags.
	StarPaths []FieldPath
	// Predicates is the field paths the query filters on: its WHERE clause's,
	// and those of every subquery's WHERE clause.
	Predicates []FieldPath
}

// Analyze extracts field-level information from a parsed CosmosDB SELECT statement.
func Analyze(stmt *ast.SelectStmt) *QueryAnalysis {
	a := &analyzer{}
	out := a.selectStmt(stmt, nil)
	qa := &QueryAnalysis{
		SelectStar: out.star,
		StarPaths:  copyPaths(out.starPaths),
		Predicates: copyPaths(a.predicates),
	}
	for i, name := range out.names {
		qa.Projections = append(qa.Projections, Projection{Name: name, SourcePaths: copyPaths(out.shapes[i].paths(nil))})
	}
	return qa
}

// analyzer resolves one statement, collecting the predicates of every SELECT
// it holds. The paths it builds are shared between the shapes that hold
// them and never changed once built; Analyze copies what it returns.
type analyzer struct {
	predicates []FieldPath
	elements   map[*shape]*shape
	selections map[selection]*shape
	suffixes   map[suffix]int
	leaves     map[string]*shape
}

// selectOutput is what one SELECT returns: its items' names and shapes, or,
// for a SELECT * that stays a star, where its documents come from.
type selectOutput struct {
	value     bool // SELECT VALUE: one bare item
	star      bool
	starPaths []FieldPath
	names     []string
	shapes    []*shape
}

// item is the shape of one value the SELECT returns: the bare item of a
// SELECT VALUE, the document of a SELECT *, or an object of the named items.
func (o selectOutput) item() *shape {
	switch {
	case o.star:
		return leaf(o.starPaths)
	case o.value && len(o.shapes) == 1:
		return o.shapes[0]
	}
	return &shape{kind: objectShape, keys: o.names, members: o.shapes}
}

// scope is the names one SELECT's FROM clause binds, each to the shape of
// the value it names, inside the SELECT that holds it: a subquery may name
// its outer query's sources.
type scope struct {
	parent *scope
	names  map[string]*shape
}

func (s *scope) lookup(name string) *shape {
	for ; s != nil; s = s.parent {
		if b, ok := s.names[name]; ok {
			return b
		}
	}
	return nil
}

func (a *analyzer) selectStmt(stmt *ast.SelectStmt, parent *scope) selectOutput {
	sc := &scope{parent: parent, names: map[string]*shape{}}
	var first *shape
	if stmt.From != nil {
		first = a.register(stmt.From, sc)
	}
	for _, join := range stmt.Joins {
		a.register(join.Source, sc)
	}

	var out selectOutput
	switch {
	case stmt.Star && first != nil && first.kind == objectShape:
		// SELECT * over a subquery that returns objects of known members
		// returns those members.
		out.names, out.shapes = first.keys, first.members
	case stmt.Star:
		out.star = true
		if first != nil {
			out.starPaths = first.paths(nil)
		}
	default:
		out.value = stmt.Value
		unnamed := 0
		for _, t := range stmt.Targets {
			name := ""
			switch {
			case stmt.Value:
			case t.Alias != nil:
				name = *t.Alias
			default:
				if n, ok := propertyName(t.Expr); ok {
					name = n
				} else {
					unnamed++
					name = "$" + strconv.Itoa(unnamed)
				}
			}
			out.names = append(out.names, name)
			out.shapes = append(out.shapes, a.shapeOf(t.Expr, sc))
		}
	}

	if stmt.Where != nil {
		// A subquery in the clause adds its own predicates while it is
		// resolved, so the clause's paths are appended after it.
		where := a.exprPaths(nil, stmt.Where, sc)
		a.predicates = append(a.predicates, where...)
	}
	return out
}

// propertyName is the name Cosmos DB returns an unaliased item under, when
// the item is a property path: its identifier, its last property name, or
// the string its last bracket holds. Any other item is numbered.
func propertyName(expr ast.ExprNode) (string, bool) {
	switch e := expr.(type) {
	case *ast.ColumnRef:
		return e.Name, true
	case *ast.DotAccessExpr:
		return e.Property, true
	case *ast.BracketAccessExpr:
		if s, ok := e.Index.(*ast.StringLit); ok {
			return s.Val, true
		}
	}
	return "", false
}

// register binds a FROM or JOIN source in sc and returns its shape.
func (a *analyzer) register(te ast.TableExpr, sc *scope) *shape {
	var (
		name string
		b    *shape
	)
	switch src := te.(type) {
	case *ast.AliasedTableExpr:
		name, b = src.Alias, a.source(src.Source, sc)
	case *ast.ArrayIterationExpr:
		// The alias names each element of what IN reads.
		name, b = src.Alias, a.element(a.source(src.Source, sc))
	default:
		// An unaliased source is named by its last identifier: FROM c binds
		// c, FROM c.address binds address.
		name, b = sourceName(te), a.source(te, sc)
	}
	if name != "" {
		sc.names[name] = b
	}
	return b
}

// sourceName is the name an unaliased FROM source binds.
func sourceName(te ast.TableExpr) string {
	switch src := te.(type) {
	case *ast.ContainerRef:
		if src.Root {
			return "root"
		}
		return src.Name
	case *ast.DotAccessExpr:
		return src.Property
	}
	return ""
}

// source resolves what a FROM or JOIN source reads.
func (a *analyzer) source(te ast.TableExpr, sc *scope) *shape {
	switch src := te.(type) {
	case *ast.ContainerRef:
		name := src.Name
		if src.Root {
			name = "root"
		}
		if b := sc.lookup(name); b != nil {
			return b
		}
		return leaf([]FieldPath{{ItemSelector(name)}})
	case *ast.DotAccessExpr:
		return a.chainShape(src, sc)
	case *ast.BracketAccessExpr:
		return a.chainShape(src, sc)
	case *ast.SubqueryExpr:
		return a.selectStmt(src.Select, sc).item()
	case *ast.AliasedTableExpr:
		return a.source(src.Source, sc)
	case *ast.ArrayIterationExpr:
		return a.element(a.source(src.Source, sc))
	}
	return leaf(nil)
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

// copyPath returns a copy of a FieldPath.
func copyPath(fp FieldPath) FieldPath {
	cp := make(FieldPath, len(fp))
	copy(cp, fp)
	return cp
}

func copyPaths(fps []FieldPath) []FieldPath {
	if fps == nil {
		return nil
	}
	out := make([]FieldPath, len(fps))
	for i, fp := range fps {
		out[i] = copyPath(fp)
	}
	return out
}
