package catalog

import (
	"strings"

	"github.com/bytebase/omni/oracle/ast"
)

// packageBody compiles a package body against its spec: the spec's items
// are in scope, and each subprogram the spec declares needs a body with
// the same parameters and return type (PLS-00323).
func (c *plChecker) packageBody(n *ast.CreatePackageStmt) {
	spec := c.packageItems(c.name)
	if spec == nil {
		p := c.e.packages[c.name]
		if p == nil || func() bool { s, _ := splitPackageUnits(p.Definition); return s == "" }() {
			c.fail("PLS-00201", "identifier '%s' must be declared", c.name)
		}
		// A spec the checker cannot read: the body's names are unknown.
		return
	}
	for name, sym := range spec {
		c.declare(name, sym)
	}
	c.push()
	defer c.pop()
	c.declarations(n.Body)
	defined := map[string][]string{}
	for _, item := range listItems(n.Body) {
		switch d := item.(type) {
		case *ast.CreateProcedureStmt:
			if d.Body != nil || d.CallSpec != nil {
				defined[objectName(d.Name)] = append(defined[objectName(d.Name)], signatureText(d.Parameters, nil))
			}
		case *ast.CreateFunctionStmt:
			if d.Body != nil || d.CallSpec != nil {
				defined[objectName(d.Name)] = append(defined[objectName(d.Name)], signatureText(d.Parameters, d.ReturnType))
			}
		}
	}
	for name, sym := range spec {
		if sym.kind != plSubprogram {
			continue
		}
		for _, sig := range sym.subprograms {
			if sig.text == "" {
				continue
			}
			found := false
			for _, d := range defined[name] {
				if d == sig.text {
					found = true
				}
			}
			if !found {
				c.fail("PLS-00323", "subprogram or cursor '%s' is declared in a package specification and must be defined in the package body", name)
			}
		}
	}
	if n.InitBlock != nil {
		c.block(n.InitBlock)
	}
}

// signatureText renders a subprogram's parameters and return type as
// PLS-00323 compares a spec's declaration with the body's definition.
func signatureText(params *ast.List, ret *ast.TypeName) string {
	var b strings.Builder
	for _, p := range paramList(params) {
		mode := strings.ToUpper(strings.Join(strings.Fields(p.Mode), " "))
		if mode == "" {
			mode = "IN"
		}
		b.WriteString(p.Name + " " + mode + " " + typeText(p.TypeName) + ";")
	}
	if ret != nil {
		b.WriteString("RETURN " + typeText(ret))
	}
	return b.String()
}

func typeText(t *ast.TypeName) string {
	if t == nil {
		return ""
	}
	s := strings.Join(typeParts(t), ".")
	switch {
	case t.IsPercType:
		s += "%TYPE"
	case t.IsPercRowtype:
		s += "%ROWTYPE"
	}
	return s
}

// trigger compiles a trigger's body, whose :NEW and :OLD (or what
// REFERENCING names them) read the table's columns.
func (c *plChecker) trigger(n *ast.CreateTriggerStmt) {
	table := c.e.resolveSynonymTarget(n.Table.Name, n.Table.Schema)
	var cols map[string]bool
	if v := c.e.views[table]; v != nil {
		cols = columnSet(v.Columns)
	} else if t := c.e.tables[table]; t != nil {
		cols = columnSet(t.Columns)
	}
	c.read(n.Table.Name)
	if n.Compound || cols == nil {
		return
	}
	newName, oldName := "NEW", "OLD"
	if r := n.Referencing; r != nil {
		if r.NewAlias != "" {
			newName = strings.ToUpper(r.NewAlias)
		}
		if r.OldAlias != "" {
			oldName = strings.ToUpper(r.OldAlias)
		}
	}
	c.row = map[string]map[string]bool{newName: cols, oldName: cols}
	c.rowTable = table
	// The WHEN condition reads the row without colons: NEW.b > 0.
	ast.Inspect(n.When, func(x ast.Node) bool {
		if ref, ok := x.(*ast.ColumnRef); ok && c.row[ref.Table] != nil && cols[ref.Column] {
			c.read(table + "." + ref.Column)
		}
		return true
	})
	c.block(n.Body)
}

// objectType compiles an object or collection type's spec, or its body.
func (c *plChecker) objectType(n *ast.CreateTypeStmt) {
	if n.IsBody {
		c.typeBody(n)
		return
	}
	if n.Under != nil {
		super := c.e.units.types[n.Under.Name]
		switch {
		case super == nil || !super.spec:
			c.fail("PLS-00201", "identifier '%s' must be declared", n.Under.Name)
		case super.final:
			c.fail("PLS-00590", "attempt to create a subtype UNDER a FINAL type")
		}
		c.read(n.Under.Name)
	}
	for _, item := range listItems(n.Attributes) {
		if col, ok := item.(*ast.ColumnDef); ok {
			c.typeName(col.TypeName)
		}
	}
	c.typeName(n.AsTable)
	c.typeName(n.AsVarray)
	for _, item := range listItems(n.Methods) {
		m, ok := item.(*ast.TypeBodyMember)
		if !ok {
			continue
		}
		switch s := m.Subprog.(type) {
		case *ast.CreateProcedureStmt:
			c.subprogram(s.Parameters, nil, nil, false, false)
		case *ast.CreateFunctionStmt:
			c.subprogram(s.Parameters, memberReturn(m, s), nil, true, false)
		}
	}
}

// typeBody compiles a type body against its spec: each member is one the
// spec declares (PLS-00539), and reads SELF and the type's attributes.
func (c *plChecker) typeBody(n *ast.CreateTypeStmt) {
	t := c.e.units.types[c.name]
	if t == nil || !t.spec {
		c.fail("PLS-00201", "identifier '%s' must be declared", c.name)
		return
	}
	attrs := map[string]bool{}
	for _, a := range t.attrs {
		attrs[a] = true
		c.declare(a, &plSym{kind: plVariable})
	}
	methods := map[string]bool{}
	for _, m := range t.methods {
		methods[m] = true
		attrs[m] = true
		// A method may call another, or itself, without SELF.
		c.declare(m, &plSym{kind: plSubprogram})
	}
	c.declare("SELF", &plSym{kind: plVariable, fields: attrs})
	for _, item := range listItems(n.Body) {
		m, ok := item.(*ast.TypeBodyMember)
		if !ok {
			continue
		}
		name := subprogramName(m.Subprog)
		if !methods[name] && m.Kind != ast.TYPE_BODY_CONSTRUCTOR {
			c.fail("PLS-00539", "subprogram '%s' is declared in an object type body and must be defined in the object type specification", name)
		}
		switch s := m.Subprog.(type) {
		case *ast.CreateProcedureStmt:
			c.subprogram(s.Parameters, nil, s.Body, false, false)
		case *ast.CreateFunctionStmt:
			// A constructor ends with a bare RETURN, as a pipelined
			// function does.
			c.subprogram(s.Parameters, memberReturn(m, s), s.Body, true, m.Kind == ast.TYPE_BODY_CONSTRUCTOR)
		}
	}
}

// memberReturn is the return type a type's method declares, or nil for a
// constructor's RETURN SELF AS RESULT, which names no type.
func memberReturn(m *ast.TypeBodyMember, f *ast.CreateFunctionStmt) *ast.TypeName {
	if m.Kind == ast.TYPE_BODY_CONSTRUCTOR {
		return nil
	}
	return f.ReturnType
}
