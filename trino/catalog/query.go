package catalog

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/bytebase/omni/trino/ast"
	"github.com/bytebase/omni/trino/parser"
)

// queryScope is the relations a query's FROM clause holds, each with the
// name a column may be qualified by and its columns.
type queryScope struct {
	rels []scopeRel
}

type scopeRel struct {
	alias string
	cols  []*Column
}

// queryColumns returns the output columns of a CREATE TABLE AS or view
// query: a name (empty when Trino would make one up) and a type (empty when
// omni does not know it). It reports what Trino's analyzer refuses about
// the relations a query reads: a missing catalog, schema or table, and a
// column no relation has.
func (c *Catalog) queryColumns(q *parser.Query) ([]*Column, error) {
	return c.typeQuery(q, map[string][]*Column{})
}

func (c *Catalog) typeQuery(q *parser.Query, ctes map[string][]*Column) ([]*Column, error) {
	if q == nil {
		return nil, nil
	}
	if q.With != nil {
		inner := make(map[string][]*Column, len(ctes)+len(q.With.CTEs))
		for k, v := range ctes {
			inner[k] = v
		}
		for _, cte := range q.With.CTEs {
			cols, err := c.typeQuery(cte.Query, inner)
			if err != nil {
				return nil, err
			}
			cols = copyColumns(cols)
			if len(cte.ColumnAliases) > 0 {
				for i, a := range cte.ColumnAliases {
					if i < len(cols) {
						cols[i].Name = fold(a)
					}
				}
			}
			inner[fold(cte.Name)] = cols
		}
		ctes = inner
	}
	return c.typeBody(q.Body, ctes)
}

func (c *Catalog) typeBody(body parser.QueryNode, ctes map[string][]*Column) ([]*Column, error) {
	switch n := body.(type) {
	case *parser.ParenQuery:
		return c.typeQuery(n.Inner, ctes)
	case *parser.TableQuery:
		cols, err := c.relationColumns(n.Name, ctes)
		return copyColumns(cols), err
	case *parser.QuerySpec:
		return c.typeSpec(n, ctes)
	case *parser.SetOperation:
		// The output takes the left side's names; a column keeps its type
		// when both sides agree on it.
		left, err := c.typeBody(n.Left, ctes)
		if err != nil {
			return nil, err
		}
		right, err := c.typeBody(n.Right, ctes)
		if err != nil {
			return nil, err
		}
		out := copyColumns(left)
		for i, col := range out {
			if i >= len(right) || right[i].Type != col.Type {
				col.Type = ""
			}
		}
		return out, nil
	}
	return nil, nil
}

func (c *Catalog) typeSpec(spec *parser.QuerySpec, ctes map[string][]*Column) ([]*Column, error) {
	scope := &queryScope{}
	for _, rel := range spec.From {
		if err := c.addRelation(scope, rel, ctes); err != nil {
			return nil, err
		}
	}
	var out []*Column
	for _, item := range spec.Items {
		switch item.Kind {
		case parser.SelectAll:
			for _, r := range scope.rels {
				out = append(out, copyColumns(r.cols)...)
			}
		case parser.SelectAllFrom:
			if ref, ok := item.Expr.(*parser.ColumnRef); ok {
				for _, r := range scope.rels {
					if r.alias == fold(ref.Name) {
						out = append(out, copyColumns(r.cols)...)
					}
				}
			}
		default:
			col, err := scope.typeExpr(item.Expr)
			if err != nil {
				return nil, err
			}
			if item.Alias != nil {
				col.Name = fold(item.Alias)
			}
			out = append(out, col)
		}
	}
	return out, nil
}

// addRelation adds a FROM item's relations to scope.
func (c *Catalog) addRelation(scope *queryScope, rel parser.Relation, ctes map[string][]*Column) error {
	switch n := rel.(type) {
	case *parser.AliasedRelation:
		before := len(scope.rels)
		if err := c.addRelation(scope, n.Inner, ctes); err != nil {
			return err
		}
		if len(scope.rels) == before+1 {
			r := &scope.rels[before]
			if n.Alias != nil {
				r.alias = fold(n.Alias)
			}
			if len(n.ColumnAliases) > 0 {
				r.cols = copyColumns(r.cols)
				for i, a := range n.ColumnAliases {
					if i < len(r.cols) {
						r.cols[i].Name = fold(a)
					}
				}
			}
		}
	case *parser.TableRelation:
		cols, err := c.relationColumns(n.Name, ctes)
		if err != nil {
			return err
		}
		alias := ""
		if len(n.Name.Parts) > 0 {
			alias = fold(n.Name.Parts[len(n.Name.Parts)-1])
		}
		scope.rels = append(scope.rels, scopeRel{alias: alias, cols: cols})
	case *parser.Join:
		if err := c.addRelation(scope, n.Left, ctes); err != nil {
			return err
		}
		return c.addRelation(scope, n.Right, ctes)
	case *parser.SubqueryRelation:
		cols, err := c.typeQuery(n.Query, ctes)
		if err != nil {
			return err
		}
		scope.rels = append(scope.rels, scopeRel{cols: cols})
	case *parser.ParenRelation:
		return c.addRelation(scope, n.Inner, ctes)
	default:
		// UNNEST, LATERAL, table functions: columns omni does not know.
		scope.rels = append(scope.rels, scopeRel{})
	}
	return nil
}

// relationColumns returns the columns of a table, view or WITH query a
// FROM clause names, refusing a name Trino cannot resolve.
func (c *Catalog) relationColumns(qn *ast.QualifiedName, ctes map[string][]*Column) ([]*Column, error) {
	if len(qn.Parts) == 1 {
		if cols, ok := ctes[fold(qn.Parts[0])]; ok {
			return cols, nil
		}
	}
	name, err := c.qualify(qn)
	if err != nil {
		return nil, err
	}
	db, err := c.database(name.catalog)
	if err != nil {
		return nil, err
	}
	s := db.schemas[name.schema]
	if s == nil {
		return nil, errorf("SCHEMA_NOT_FOUND", "Schema '%s' does not exist", name.schemaString())
	}
	if t := s.tables[name.name]; t != nil {
		return t.columns, nil
	}
	if v := s.views[name.name]; v != nil {
		return v.columns, nil
	}
	return nil, errorf("TABLE_NOT_FOUND", "Table '%s' does not exist", name)
}

// typeExpr names and types one select item; an expression omni does not
// type has an empty type.
func (s *queryScope) typeExpr(e parser.Expr) (*Column, error) {
	switch n := e.(type) {
	case *parser.ColumnRef:
		return s.column("", fold(n.Name))
	case *parser.Dereference:
		if base, ok := n.Base.(*parser.ColumnRef); ok {
			return s.column(fold(base.Name), fold(n.FieldName))
		}
		return &Column{Name: fold(n.FieldName), Nullable: true}, nil
	case *parser.Literal:
		return &Column{Type: literalType(n), Nullable: true}, nil
	case *parser.TypeConstructor:
		if strings.EqualFold(n.Name, "DATE") {
			return &Column{Type: "date", Nullable: true}, nil
		}
	}
	return &Column{Nullable: true}, nil
}

// column resolves a column reference, qualified by a relation's name or
// not.
func (s *queryScope) column(qualifier, name string) (*Column, error) {
	known := true
	for _, r := range s.rels {
		if qualifier != "" && r.alias != qualifier {
			continue
		}
		if r.cols == nil {
			known = false
		}
		for _, col := range r.cols {
			if col.Name == name {
				return &Column{Name: col.Name, Type: col.Type, Nullable: true}, nil
			}
		}
	}
	if !known {
		return &Column{Name: name, Nullable: true}, nil
	}
	if qualifier != "" {
		return nil, errorf("COLUMN_NOT_FOUND", "Column '%s.%s' cannot be resolved", qualifier, name)
	}
	return nil, errorf("COLUMN_NOT_FOUND", "Column '%s' cannot be resolved", name)
}

// literalType is the type Trino gives a literal.
func literalType(l *parser.Literal) string {
	switch l.Kind {
	case parser.LiteralBool:
		return "boolean"
	case parser.LiteralString:
		return "varchar(" + strconv.Itoa(utf8.RuneCountInString(l.Value)) + ")"
	case parser.LiteralInteger:
		if v, err := strconv.ParseInt(l.Value, 10, 64); err == nil && v <= 2147483647 {
			return "integer"
		}
		return "bigint"
	case parser.LiteralDecimal:
		intPart, frac, _ := strings.Cut(l.Value, ".")
		intPart = strings.TrimLeft(intPart, "0")
		p := len(intPart) + len(frac)
		if p == 0 {
			p = 1
		}
		return "decimal(" + strconv.Itoa(p) + "," + strconv.Itoa(len(frac)) + ")"
	case parser.LiteralDouble:
		return "double"
	case parser.LiteralBinary:
		return "varbinary"
	}
	return "unknown"
}

// copyColumns returns a query's copy of columns: their names and types,
// nullable, with no default or comment.
func copyColumns(cols []*Column) []*Column {
	out := make([]*Column, len(cols))
	for i, col := range cols {
		out[i] = &Column{Name: col.Name, Type: col.Type, Nullable: true}
	}
	return out
}
