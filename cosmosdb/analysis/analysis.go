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
	return &QueryAnalysis{
		Projections: out.projections,
		SelectStar:  out.star,
		StarPaths:   out.starPaths,
		Predicates:  a.predicates,
	}
}

// analyzer resolves one statement, collecting the predicates of every SELECT
// it holds.
type analyzer struct {
	predicates []FieldPath
}

// selectOutput is what one SELECT returns.
type selectOutput struct {
	value       bool // SELECT VALUE: one bare item, Projections[0]
	star        bool
	starPaths   []FieldPath
	projections []Projection
}

// paths is every source of what the SELECT returns.
func (o selectOutput) paths() []FieldPath {
	if o.star {
		return copyPaths(o.starPaths)
	}
	var out []FieldPath
	for _, p := range o.projections {
		out = append(out, copyPaths(p.SourcePaths)...)
	}
	return out
}

// A binding is what a FROM or JOIN source names: the paths its value comes
// from and, for a subquery that names its items, those items, which a
// property access selects by name.
type binding struct {
	paths  []FieldPath
	fields []Projection
}

// field returns the sources of the bound value's property name.
func (b *binding) field(name string) ([]FieldPath, bool) {
	if b.fields == nil {
		return nil, false
	}
	for _, f := range b.fields {
		if f.Name == name {
			return copyPaths(f.SourcePaths), true
		}
	}
	// A property the subquery does not return is undefined; its sources are
	// all of the subquery's, which is the safe answer for masking.
	return copyPaths(b.paths), true
}

// scope is the names one SELECT's FROM clause binds, inside the SELECT that
// holds it: a subquery may name its outer query's sources.
type scope struct {
	parent *scope
	names  map[string]*binding
}

func (s *scope) lookup(name string) *binding {
	for ; s != nil; s = s.parent {
		if b, ok := s.names[name]; ok {
			return b
		}
	}
	return nil
}

func (a *analyzer) selectStmt(stmt *ast.SelectStmt, parent *scope) selectOutput {
	sc := &scope{parent: parent, names: map[string]*binding{}}
	var first *binding
	if stmt.From != nil {
		first = a.register(stmt.From, sc)
	}
	for _, join := range stmt.Joins {
		a.register(join.Source, sc)
	}

	var out selectOutput
	switch {
	case stmt.Star && first != nil && first.fields != nil:
		// SELECT * over a subquery returns the subquery's items.
		out.projections = make([]Projection, len(first.fields))
		for i, f := range first.fields {
			out.projections[i] = Projection{Name: f.Name, SourcePaths: copyPaths(f.SourcePaths)}
		}
	case stmt.Star:
		out.star = true
		if first != nil {
			out.starPaths = copyPaths(first.paths)
		}
	default:
		out.value = stmt.Value
		unnamed := 0
		for _, t := range stmt.Targets {
			p := Projection{SourcePaths: a.exprPaths(nil, t.Expr, sc)}
			switch {
			case stmt.Value:
			case t.Alias != nil:
				p.Name = *t.Alias
			default:
				if name, ok := propertyName(t.Expr); ok {
					p.Name = name
				} else {
					unnamed++
					p.Name = "$" + strconv.Itoa(unnamed)
				}
			}
			out.projections = append(out.projections, p)
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

// register binds a FROM or JOIN source in sc and returns its binding.
func (a *analyzer) register(te ast.TableExpr, sc *scope) *binding {
	var (
		name string
		b    *binding
	)
	switch src := te.(type) {
	case *ast.AliasedTableExpr:
		name, b = src.Alias, a.source(src.Source, sc)
	case *ast.ArrayIterationExpr:
		name, b = src.Alias, a.source(src.Source, sc)
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
func (a *analyzer) source(te ast.TableExpr, sc *scope) *binding {
	switch src := te.(type) {
	case *ast.ContainerRef:
		name := src.Name
		if src.Root {
			name = "root"
		}
		if b := sc.lookup(name); b != nil {
			return b
		}
		return &binding{paths: []FieldPath{{ItemSelector(name)}}}
	case *ast.DotAccessExpr:
		return &binding{paths: a.exprPaths(nil, src, sc)}
	case *ast.BracketAccessExpr:
		return &binding{paths: a.exprPaths(nil, src, sc)}
	case *ast.SubqueryExpr:
		out := a.selectStmt(src.Select, sc)
		b := &binding{paths: out.paths()}
		if !out.value && !out.star {
			b.fields = out.projections
		}
		return b
	case *ast.AliasedTableExpr:
		return a.source(src.Source, sc)
	case *ast.ArrayIterationExpr:
		return a.source(src.Source, sc)
	}
	return &binding{}
}

// resolveChain resolves a property path rooted at an identifier: the
// identifier's binding, then each selector. It returns false for a chain
// rooted at any other expression.
func (a *analyzer) resolveChain(dst []FieldPath, expr ast.ExprNode, sc *scope) ([]FieldPath, bool) {
	var sels []Selector
	root := expr
	for {
		switch e := root.(type) {
		case *ast.DotAccessExpr:
			sels = append(sels, ItemSelector(e.Property))
			root = e.Expr
			continue
		case *ast.BracketAccessExpr:
			sels = append(sels, bracketSelector(e.Index))
			root = e.Expr
			continue
		}
		break
	}
	ref, ok := root.(*ast.ColumnRef)
	if !ok {
		return dst, false
	}
	// sels runs from the last selector to the first.
	for i, j := 0, len(sels)-1; i < j; i, j = i+1, j-1 {
		sels[i], sels[j] = sels[j], sels[i]
	}
	var base []FieldPath
	b := sc.lookup(ref.Name)
	switch {
	case b == nil:
		base = []FieldPath{{ItemSelector(ref.Name)}}
	case len(sels) > 0 && !sels[0].IsArray() && b.fields != nil:
		base, _ = b.field(sels[0].Name)
		sels = sels[1:]
	default:
		base = copyPaths(b.paths)
	}
	for _, p := range base {
		dst = append(dst, append(p, sels...))
	}
	return dst, true
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

// copyPath returns a shallow copy of a FieldPath to avoid alias mutation.
func copyPath(fp FieldPath) FieldPath {
	cp := make(FieldPath, len(fp))
	copy(cp, fp)
	return cp
}

func copyPaths(fps []FieldPath) []FieldPath {
	out := make([]FieldPath, len(fps))
	for i, fp := range fps {
		out[i] = copyPath(fp)
	}
	return out
}
