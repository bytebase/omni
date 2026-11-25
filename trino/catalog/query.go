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
	// names is what a reference may qualify the relation's columns by,
	// as a suffix: its alias, or catalog, schema and name.
	names []string
	cols  []*Column
	// using names the columns a JOIN ... USING gave to the relation
	// before this one, which * shows once; usingOrder is their order.
	using      map[string]bool
	usingOrder []string
	// usingFrom is the index, among the scope's relations, of the first
	// relation of the join's left side.
	usingFrom int
}

// queryColumns returns the output columns of a CREATE TABLE AS or view
// query: a name (empty when Trino would make one up) and a type (empty when
// omni does not know it). It refuses only a relation a query reads that is
// missing: its catalog, schema or table.
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
			if q.With.Recursive {
				// A recursive term reads its own query, whose columns
				// omni leaves unknown.
				inner[fold(cte.Name)] = nil
			}
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
	case *parser.ValuesQuery:
		return valuesColumns(n), nil
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
			out = append(out, joinedColumns(scope.rels, 0)...)
		case parser.SelectAllFrom:
			var cols []*Column
			if q, ok := nameChain(item.Expr); ok {
				for _, r := range scope.rels {
					if r.matches(q) {
						cols = append(cols, copyColumns(r.cols)...)
					}
				}
			}
			for i, a := range item.Aliases {
				if i < len(cols) {
					cols[i].Name = fold(a)
				}
			}
			out = append(out, cols...)
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
		if n.Alias != nil && len(scope.rels) > before+1 {
			// An alias on a parenthesized join names one relation of all
			// the join's columns.
			cols := joinedColumns(scope.rels[before:], before)
			scope.rels = append(scope.rels[:before], scopeRel{cols: cols})
		}
		if len(scope.rels) == before+1 {
			r := &scope.rels[before]
			if n.Alias != nil {
				r.names = []string{fold(n.Alias)}
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
		names := []string{fold(n.Name.Parts[0])}
		if _, cte := ctes[names[0]]; !cte || len(n.Name.Parts) > 1 {
			name, _ := c.qualify(n.Name)
			names = []string{name.catalog, name.schema, name.name}
		}
		scope.rels = append(scope.rels, scopeRel{names: names, cols: cols})
	case *parser.Join:
		left := len(scope.rels)
		if err := c.addRelation(scope, n.Left, ctes); err != nil {
			return err
		}
		before := len(scope.rels)
		if err := c.addRelation(scope, n.Right, ctes); err != nil {
			return err
		}
		if len(n.Using) > 0 && len(scope.rels) == before+1 {
			// A USING column is the join's once, from its left side.
			r := &scope.rels[before]
			r.using, r.usingFrom = map[string]bool{}, left
			for _, u := range n.Using {
				r.using[fold(u)] = true
				r.usingOrder = append(r.usingOrder, fold(u))
			}
		}
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
	case *parser.ParenExpr:
		return s.typeExpr(n.Expr)
	case *parser.ColumnRef:
		return s.column(nil, fold(n.Name))
	case *parser.Dereference:
		parts, ok := nameChain(n)
		if !ok {
			return &Column{Name: fold(n.FieldName), Nullable: true}, nil
		}
		// The longest prefix that names a relation qualifies the column;
		// the names after it are fields of a ROW value.
		k := len(parts) - 1
		for k > 0 && !s.hasRelation(parts[:k]) {
			k--
		}
		col, err := s.column(parts[:k], parts[k])
		if err != nil {
			return nil, err
		}
		for _, field := range parts[k+1:] {
			if col.Type == "" {
				col = &Column{Name: field, Nullable: true}
				continue
			}
			typ, _ := rowField(col.Type, field)
			col = &Column{Name: field, Type: typ, Nullable: true}
		}
		return col, nil
	case *parser.Literal:
		return &Column{Type: literalType(n), Nullable: true}, nil
	case *parser.TypeConstructor:
		if strings.EqualFold(n.Name, "DATE") {
			return &Column{Type: "date", Nullable: true}, nil
		}
	}
	return &Column{Nullable: true}, nil
}

// joinedColumns is what * shows of relations joined left to right, rels
// being the scope's relations from index offset on: a
// JOIN ... USING puts its columns first, in its order, once, then the
// other columns of each side.
func joinedColumns(rels []scopeRel, offset int) []*Column {
	var out []*Column
	starts := make([]int, len(rels))
	for i, r := range rels {
		starts[i] = len(out)
		if len(r.usingOrder) > 0 {
			// Only the join's left side is coalesced.
			from := 0
			if k := r.usingFrom - offset; k >= 0 && k < i {
				from = starts[k]
			}
			side := out[from:]
			var keys, rest []*Column
			for _, name := range r.usingOrder {
				for _, col := range side {
					if col.Name == name {
						keys = append(keys, col)
						break
					}
				}
			}
			for _, col := range side {
				if !r.using[col.Name] {
					rest = append(rest, col)
				}
			}
			out = append(append(out[:from:from], keys...), rest...)
		}
		for _, col := range copyColumns(r.cols) {
			if !r.using[col.Name] {
				out = append(out, col)
			}
		}
	}
	return out
}

// hasRelation reports whether a FROM relation goes by a qualifier.
func (s *queryScope) hasRelation(q []string) bool {
	for _, r := range s.rels {
		if r.matches(q) {
			return true
		}
	}
	return false
}

// matches reports whether a qualifier names the relation: it ends its
// names.
func (r scopeRel) matches(q []string) bool {
	if len(q) == 0 || len(q) > len(r.names) {
		return false
	}
	tail := r.names[len(r.names)-len(q):]
	for i := range q {
		if q[i] != tail[i] {
			return false
		}
	}
	return true
}

// nameChain reads a.b.c, a column reference and the fields after it, as
// its folded names.
func nameChain(e parser.Expr) ([]string, bool) {
	switch n := e.(type) {
	case *parser.ColumnRef:
		return []string{fold(n.Name)}, true
	case *parser.Dereference:
		base, ok := nameChain(n.Base)
		return append(base, fold(n.FieldName)), ok
	}
	return nil, false
}

// rowField returns the type of a field of a stored row type, row(x
// integer, y varchar).
func rowField(rowType, field string) (string, bool) {
	inner, ok := strings.CutPrefix(rowType, "row(")
	if !ok {
		return "", false
	}
	inner = strings.TrimSuffix(inner, ")")
	depth, start := 0, 0
	for i := 0; i <= len(inner); i++ {
		if i < len(inner) {
			switch inner[i] {
			case '(':
				depth++
				continue
			case ')':
				depth--
				continue
			case ',':
				if depth > 0 {
					continue
				}
			default:
				continue
			}
		}
		name, typ := splitRowField(strings.TrimSpace(inner[start:i]))
		if name == field {
			return typ, true
		}
		start = i + 1
	}
	return "", false
}

// splitRowField reads one field of a stored row type, a name, quoted when
// quoteName quoted it, and its type.
func splitRowField(f string) (name, typ string) {
	if !strings.HasPrefix(f, `"`) {
		name, typ, _ = strings.Cut(f, " ")
		return name, typ
	}
	for i := 1; i < len(f); i++ {
		if f[i] != '"' {
			continue
		}
		if i+1 < len(f) && f[i+1] == '"' {
			i++
			continue
		}
		return strings.ReplaceAll(f[1:i], `""`, `"`), strings.TrimSpace(f[i+1:])
	}
	return f, ""
}

// valuesColumns types the columns of a VALUES list by its rows' literals:
// unnamed, and typed when every row agrees.
func valuesColumns(v *parser.ValuesQuery) []*Column {
	var out []*Column
	for i, row := range v.Rows {
		var elems []parser.Expr
		switch r := row.(type) {
		case *parser.RowConstructor:
			elems = r.Elements
		case *parser.ParenExpr:
			elems = []parser.Expr{r.Expr}
		default:
			elems = []parser.Expr{row}
		}
		if i == 0 {
			for range elems {
				out = append(out, &Column{Nullable: true})
			}
		}
		for j, e := range elems {
			if j >= len(out) {
				break
			}
			typ := ""
			if lit, ok := e.(*parser.Literal); ok {
				typ = literalType(lit)
			}
			if i == 0 {
				out[j].Type = typ
			} else if out[j].Type != typ {
				out[j].Type = ""
			}
		}
	}
	return out
}

// column resolves a column reference, qualified by a relation's names or
// not, to the first relation that holds it; a USING column is the left
// side's. A name no relation holds, or one a relation of unknown columns
// may, has no type: Exec refuses no reference (see testdata/README.md).
func (s *queryScope) column(qualifier []string, name string) (*Column, error) {
	known := true
	var found *Column
	for _, r := range s.rels {
		if len(qualifier) > 0 && !r.matches(qualifier) || len(qualifier) == 0 && r.using[name] {
			continue
		}
		if r.cols == nil {
			known = false
		}
		for _, col := range r.cols {
			if col.Name != name {
				continue
			}
			if found == nil {
				found = &Column{Name: col.Name, Type: col.Type, Nullable: true}
			}
		}
	}
	if found != nil && known {
		return found, nil
	}
	return &Column{Name: name, Nullable: true}, nil
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
