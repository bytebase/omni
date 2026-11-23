package catalog

import (
	"strings"

	"github.com/bytebase/omni/oracle/ast"
)

// expr checks the names an expression reads.
func (c *plChecker) expr(x ast.ExprNode) {
	if x == nil || c.err != nil {
		return
	}
	switch n := x.(type) {
	case *ast.ColumnRef:
		c.columnRef(n)
		return
	case *ast.BindVariable:
		c.bind(n)
		return
	case *ast.FuncCallExpr:
		c.funcCall(n)
		return
	case *ast.SubqueryExpr:
		c.sql(n.Subquery)
		return
	case *ast.CursorAttributeExpr:
		// SQL%ROWCOUNT reads the implicit cursor, which no one declares.
		if ref, ok := n.Cursor.(*ast.ColumnRef); !ok || ref.Table != "" || ref.Column != "SQL" {
			c.expr(n.Cursor)
		}
		c.expr(n.Index)
		return
	}
	ast.Inspect(x, func(child ast.Node) bool {
		if child == ast.Node(x) {
			return true
		}
		switch n := child.(type) {
		case ast.ExprNode:
			c.expr(n)
			return false
		case ast.StmtNode:
			c.sql(n)
			return false
		}
		return true
	})
}

// columnRef checks a name an expression reads: a declared item, a record's
// field, a package's item, a sequence's NEXTVAL or CURRVAL, a function
// called without arguments, or one of STANDARD's values.
func (c *plChecker) columnRef(n *ast.ColumnRef) {
	var parts []string
	for _, p := range []string{n.Schema, n.Table, n.Column} {
		if p != "" {
			parts = append(parts, p)
		}
	}
	if len(parts) == 1 {
		name := parts[0]
		if sym := c.lookup(name); sym != nil {
			if sym.duplicate {
				c.fail("PLS-00371", "at most one declaration for '%s' is permitted", name)
			}
			return
		}
		if plsqlValues[name] || c.row != nil && triggerPredicates[name] {
			return
		}
		if c.e.functions[name] != nil {
			c.read(name)
			return
		}
		if c.knownObject(name) {
			c.read(name)
			return
		}
		c.fail("PLS-00201", "identifier '%s' must be declared", name)
		return
	}
	c.qualifiedName(parts)
}

// qualifiedName checks a dotted name: a record's field, a package's item,
// a sequence's pseudocolumn, or an object's attribute. A first part the
// checker cannot place is taken as another schema or a public synonym.
func (c *plChecker) qualifiedName(parts []string) {
	head, rest := parts[0], parts[1:]
	if sym := c.lookup(head); sym != nil {
		if sym.fields != nil && !sym.fields[rest[0]] {
			c.fail("PLS-00302", "component '%s' must be declared", rest[0])
		}
		return
	}
	target := c.e.resolveSynonymTarget(head, "")
	switch {
	case c.e.sequences[target] != nil:
		c.read(head)
	case len(rest) == 1 && (rest[0] == "NEXTVAL" || rest[0] == "CURRVAL") && !c.knownObject(head):
		// Only a sequence has NEXTVAL; another schema's would be named
		// schema.sequence.NEXTVAL.
		c.fail("PLS-00201", "identifier '%s.%s' must be declared", head, rest[0])
	case c.e.packages[target] != nil:
		c.read(head)
		if names := c.packageNames(target); names != nil && !names[rest[0]] {
			c.fail("PLS-00302", "component '%s' must be declared", rest[0])
		}
	}
}

// packageNames lists the items a package spec declares, or nil when the
// checker cannot read the spec.
func (c *plChecker) packageNames(pkg string) map[string]bool {
	sigs := c.packageItems(pkg)
	if sigs == nil {
		return nil
	}
	names := map[string]bool{}
	for name := range sigs {
		names[name] = true
	}
	return names
}

// packageItems returns a package spec's items by name, with the overloads
// of each subprogram, or nil when the checker cannot read the spec.
func (c *plChecker) packageItems(pkg string) map[string]*plSym {
	p := c.e.packages[pkg]
	if p == nil || c.reading[pkg] || c.kind == unitPackage && c.name == pkg {
		// A spec that names itself, or one already being read: its items
		// are taken as declared.
		return nil
	}
	spec, _ := splitPackageUnits(p.Definition)
	stmt, ok := parseSaved(spec).(*ast.CreatePackageStmt)
	if !ok || stmt.IsBody {
		return nil
	}
	sub := newPLChecker(c.e, unitPackage, pkg)
	sub.reading = c.reading
	c.reading[pkg] = true
	defer delete(c.reading, pkg)
	sub.push()
	sub.declarations(stmt.Body)
	return sub.scope.names
}

// resolveCall resolves a called name to its overloads. known reports a
// name the checker found, whose arguments it can check; a name it cannot
// place is taken as resolving.
func (c *plChecker) resolveCall(parts []string) (sigs []plSignature, known bool) {
	if len(parts) == 1 {
		name := parts[0]
		if sym := c.lookup(name); sym != nil {
			if sym.kind != plSubprogram {
				return nil, false
			}
			// A package body's subprograms overload the spec's, a scope
			// out: a call may take any of them.
			for s := c.scope; s != nil; s = s.parent {
				if o := s.names[name]; o != nil && o.kind == plSubprogram {
					sigs = append(sigs, o.subprograms...)
				}
			}
			return sigs, true
		}
		target := c.e.resolveSynonymTarget(name, "")
		if p := c.e.procedures[target]; p != nil {
			c.read(name)
			return savedSignature(p.Definition, false), true
		}
		if f := c.e.functions[target]; f != nil {
			c.read(name)
			return savedSignature(f.Definition, true), true
		}
		if c.e.hasType(target) || plsqlFunctions[name] || plsqlValues[name] || c.row != nil && triggerPredicates[name] {
			return nil, false
		}
		if c.knownObject(name) {
			c.read(name)
			return nil, false
		}
		c.fail("PLS-00201", "identifier '%s' must be declared", name)
		return nil, false
	}
	head := parts[0]
	if c.lookup(head) != nil {
		// A method of a declared object or collection.
		return nil, false
	}
	target := c.e.resolveSynonymTarget(head, "")
	if c.e.packages[target] == nil || len(parts) != 2 {
		return nil, false
	}
	c.read(head)
	items := c.packageItems(target)
	if items == nil {
		return nil, false
	}
	sym := items[parts[1]]
	if sym == nil {
		c.fail("PLS-00302", "component '%s' must be declared", parts[1])
		return nil, false
	}
	if sym.kind != plSubprogram {
		return nil, false
	}
	return sym.subprograms, true
}

// savedSignature reads a stored procedure's or function's parameters from
// its text.
func savedSignature(text string, function bool) []plSignature {
	if function && !strings.HasPrefix(strings.ToUpper(strings.TrimSpace(text)), "CREATE") {
		text = "CREATE " + text
	}
	switch n := parseSaved(text).(type) {
	case *ast.CreateProcedureStmt:
		return []plSignature{{params: paramList(n.Parameters)}}
	case *ast.CreateFunctionStmt:
		return []plSignature{{function: true, params: paramList(n.Parameters)}}
	}
	return nil
}

// funcCall checks a function call in an expression: a declared or schema
// function, a package's function, a built-in, a type's constructor, or an
// element of a declared collection.
func (c *plChecker) funcCall(n *ast.FuncCallExpr) {
	c.exprList(n.Args)
	if n.FuncName == nil || n.FuncName.Name == "" || n.Over != nil {
		return
	}
	var parts []string
	if n.FuncName.Schema != "" {
		parts = append(parts, n.FuncName.Schema)
	}
	parts = append(parts, n.FuncName.Name)
	if len(parts) == 1 {
		if sym := c.lookup(parts[0]); sym != nil && sym.kind != plSubprogram {
			// A collection's element, or a cursor's or type's use.
			return
		}
	}
	sigs, known := c.resolveCall(parts)
	if known {
		c.checkArgs(parts[len(parts)-1], sigs, n.Args, true)
	}
}

// triggerPredicates are the conditional predicates a DML trigger's body
// reads.
var triggerPredicates = map[string]bool{"INSERTING": true, "UPDATING": true, "DELETING": true}
