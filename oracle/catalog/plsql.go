package catalog

import (
	"fmt"
	"strings"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/oracle/ast"
	oracleparser "github.com/bytebase/omni/oracle/parser"
)

// This file compiles a PL/SQL unit's names as Oracle does when it creates
// the unit, to tell whether Oracle stores it VALID or INVALID and with
// which first error. testdata/programs/simulate.tsv holds it to the pinned
// Oracle. It reports only the errors those rows pin, and takes a name it
// cannot place, such as a public synonym's or another schema's, as one
// that resolves.

// plChecker compiles one PL/SQL unit's names.
type plChecker struct {
	e          *oracleMetadataExtractor
	kind, name string
	// reads holds the schema objects the unit names.
	reads map[string]bool
	err   *compileError
	scope *plScope
	// function reports a function body, whose RETURN takes a value, and
	// pipelined a pipelined function's.
	function, pipelined bool
	// row holds the trigger's row names (NEW, OLD) with the table's
	// columns, and rowTable the table's name as the trigger names it.
	row      map[string]map[string]bool
	rowTable string
	// labels holds the labels of the unit's blocks and statements.
	labels map[string]bool
	// reading holds the packages whose specs the checker is reading, which
	// a spec that names itself, or two that name each other, would read
	// again without end.
	reading map[string]bool
}

// plScope is one declarative region.
type plScope struct {
	parent *plScope
	names  map[string]*plSym
}

type plKind int

const (
	plVariable plKind = iota
	plConstant
	plInParam
	plCursor
	plException
	plType
	plSubprogram
)

// plSym is one declared name.
type plSym struct {
	kind plKind
	// fields holds a record's fields, or nil when the checker does not
	// know them.
	fields map[string]bool
	// subprograms holds a subprogram's overloads, by their parameters.
	subprograms []plSignature
	// duplicate marks a name declared twice in one scope (PLS-00371 where
	// it is used).
	duplicate bool
}

// plSignature is one subprogram's parameters and kind, and its
// signatureText.
type plSignature struct {
	function bool
	params   []*ast.Parameter
	text     string
}

func newPLChecker(e *oracleMetadataExtractor, kind, name string) *plChecker {
	return &plChecker{e: e, kind: kind, name: name, reads: map[string]bool{}, labels: map[string]bool{}, reading: map[string]bool{}}
}

// fail records the unit's first error.
func (c *plChecker) fail(code, format string, args ...any) {
	if c.err == nil {
		c.err = &compileError{Kind: c.kind, Name: c.name, Code: code, Message: fmt.Sprintf(format, args...)}
	}
}

func (c *plChecker) push() { c.scope = &plScope{parent: c.scope, names: map[string]*plSym{}} }
func (c *plChecker) pop()  { c.scope = c.scope.parent }

func (c *plChecker) declare(name string, sym *plSym) {
	if c.scope == nil {
		c.push()
	}
	if old := c.scope.names[name]; old != nil && old.kind == plSubprogram && sym.kind == plSubprogram {
		old.subprograms = append(old.subprograms, sym.subprograms...)
		return
	}
	c.scope.names[name] = sym
}

// lookup finds a declared name, innermost first.
func (c *plChecker) lookup(name string) *plSym {
	for s := c.scope; s != nil; s = s.parent {
		if sym := s.names[name]; sym != nil {
			return sym
		}
	}
	return nil
}

// read records a schema object the unit names, other than itself.
func (c *plChecker) read(name string) {
	if name != c.name {
		c.reads[name] = true
	}
}

// readVia records a name the unit reads and the object it resolves to,
// through a synonym, so that a change to either reaches the unit.
func (c *plChecker) readVia(name, target string) {
	c.read(name)
	c.read(target)
}

// unit compiles a CREATE statement's unit.
func (c *plChecker) unit(n ast.Node) {
	// A GOTO may name a label anywhere in the unit; the checker does not
	// hold it to the label's scope.
	ast.Inspect(n, func(x ast.Node) bool {
		switch x := x.(type) {
		case *ast.PLSQLLabeledStatement:
			for _, l := range x.Labels {
				c.labels[l] = true
			}
		case *ast.PLSQLBlock:
			c.labels[x.Label] = true
		case *ast.PLSQLLoop:
			c.labels[x.Label] = true
		}
		return true
	})
	c.push()
	switch n := n.(type) {
	case *ast.CreateProcedureStmt:
		c.declare(c.name, &plSym{kind: plSubprogram, subprograms: []plSignature{{params: paramList(n.Parameters)}}})
		c.subprogram(n.Parameters, nil, n.Body, false, false)
	case *ast.CreateFunctionStmt:
		c.declare(c.name, &plSym{kind: plSubprogram, subprograms: []plSignature{{function: true, params: paramList(n.Parameters)}}})
		c.subprogram(n.Parameters, n.ReturnType, n.Body, true, n.Pipelined)
	case *ast.CreatePackageStmt:
		if n.IsBody {
			c.packageBody(n)
		} else {
			c.declarations(n.Body)
		}
	case *ast.CreateTriggerStmt:
		c.trigger(n)
	case *ast.CreateTypeStmt:
		c.objectType(n)
	}
}

func paramList(list *ast.List) []*ast.Parameter {
	var out []*ast.Parameter
	for _, item := range listItems(list) {
		if p, ok := item.(*ast.Parameter); ok {
			out = append(out, p)
		}
	}
	return out
}

// subprogram compiles a procedure's or function's parameters, return
// type, and body. A pipelined function returns its rows by PIPE ROW and
// ends with a bare RETURN.
func (c *plChecker) subprogram(params *ast.List, ret *ast.TypeName, body ast.StmtNode, function, pipelined bool) {
	seen := map[string]bool{}
	for _, p := range paramList(params) {
		c.typeName(p.TypeName)
		if p.Default != nil {
			c.expr(p.Default)
		}
		if seen[p.Name] {
			c.fail("PLS-00410", "duplicate fields in RECORD,TABLE or argument list are not permitted")
		}
		seen[p.Name] = true
	}
	if ret != nil {
		c.typeName(ret)
	}
	if body == nil {
		return
	}
	c.push()
	for _, p := range paramList(params) {
		kind := plVariable
		if mode := strings.ToUpper(p.Mode); mode == "" || mode == "IN" {
			kind = plInParam
		}
		c.declare(p.Name, &plSym{kind: kind, fields: c.typeFields(p.TypeName)})
	}
	savedFunction, savedPipelined := c.function, c.pipelined
	c.function, c.pipelined = function, pipelined
	c.block(body)
	c.function, c.pipelined = savedFunction, savedPipelined
	c.pop()
}

// declarations compiles a declarative part: a block's, or a package
// spec's or body's.
func (c *plChecker) declarations(list *ast.List) {
	for _, item := range listItems(list) {
		c.declaration(item)
	}
}

func (c *plChecker) declaration(item ast.Node) {
	switch d := item.(type) {
	case *ast.PLSQLVarDecl:
		if d.TypeName != nil && typeWord(d.TypeName) == "EXCEPTION" {
			c.declare(d.Name, &plSym{kind: plException})
			return
		}
		c.typeName(d.TypeName)
		if d.Default != nil {
			c.expr(d.Default)
		}
		kind := plVariable
		if d.Constant {
			kind = plConstant
		}
		if old := c.scope.names[d.Name]; old != nil && (old.kind == plVariable || old.kind == plConstant) {
			// A name declared twice is an error only where it is used.
			old.duplicate = true
			return
		}
		c.declare(d.Name, &plSym{kind: kind, fields: c.typeFields(d.TypeName)})
	case *ast.PLSQLCursorDecl:
		c.push()
		for _, item := range listItems(d.Parameters) {
			// The parser keeps a cursor's parameters as declarations.
			if p, ok := item.(*ast.PLSQLVarDecl); ok {
				c.typeName(p.TypeName)
				c.declare(p.Name, &plSym{kind: plInParam})
			}
		}
		var fields map[string]bool
		if d.ReturnType != nil {
			c.typeName(d.ReturnType)
			fields = c.typeFields(d.ReturnType)
		}
		if d.Query != nil {
			fields = c.sql(d.Query)
		}
		c.pop()
		c.declare(d.Name, &plSym{kind: plCursor, fields: fields})
	case *ast.PLSQLTypeDecl:
		c.typeName(d.ElementType)
		c.typeName(d.ReturnType)
		var fields map[string]bool
		if strings.EqualFold(d.Kind, "RECORD") {
			fields = map[string]bool{}
			for _, f := range listItems(d.Fields) {
				switch f := f.(type) {
				case *ast.ColumnDef:
					fields[f.Name] = true
					c.typeName(f.TypeName)
				case *ast.PLSQLVarDecl:
					fields[f.Name] = true
					c.typeName(f.TypeName)
				default:
					fields = nil
				}
				if fields == nil {
					break
				}
			}
		}
		c.declare(d.Name, &plSym{kind: plType, fields: fields})
	case *ast.PLSQLSubtypeDecl:
		c.typeName(d.BaseType)
		c.declare(d.Name, &plSym{kind: plType, fields: c.typeFields(d.BaseType)})
	case *ast.CreateProcedureStmt:
		name := objectName(d.Name)
		c.declare(name, &plSym{kind: plSubprogram, subprograms: []plSignature{{params: paramList(d.Parameters), text: signatureText(d.Parameters, nil)}}})
		c.subprogram(d.Parameters, nil, d.Body, false, false)
	case *ast.CreateFunctionStmt:
		name := objectName(d.Name)
		c.declare(name, &plSym{kind: plSubprogram, subprograms: []plSignature{{function: true, params: paramList(d.Parameters), text: signatureText(d.Parameters, d.ReturnType)}}})
		c.subprogram(d.Parameters, d.ReturnType, d.Body, true, d.Pipelined)
	}
}

// typeWord is a type name's first word.
func typeWord(t *ast.TypeName) string {
	parts := typeParts(t)
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}

func typeParts(t *ast.TypeName) []string {
	if t == nil {
		return nil
	}
	var parts []string
	for _, item := range listItems(t.Names) {
		if s, ok := item.(*ast.String); ok {
			parts = append(parts, s.Str)
		}
	}
	return parts
}

// unseen reports a relation name the catalog cannot see into: a
// dictionary view, or a synonym for another schema's table.
func (c *plChecker) unseen(name string) bool {
	n := &ast.ObjectName{Name: name}
	return c.dictionary(n) || c.outside(n)
}

// relation finds the table, view or materialized view a name reads,
// through a synonym. found reports one that exists; cols is empty when the
// catalog does not know its columns (a materialized view's, a view it could
// not derive), and then no column is missing from it.
func (c *plChecker) relation(name string) (cols []*metadatapb.ColumnMetadata, found bool) {
	target := c.e.resolveSynonymTarget(name, "")
	switch {
	case c.e.views[target] != nil:
		cols, found = c.e.views[target].Columns, true
	case c.e.tables[target] != nil:
		cols, found = c.e.tables[target].Columns, true
	case c.e.materializedViews[target] != nil:
		found = true
	}
	if found {
		c.readVia(name, target)
	}
	return cols, found
}

func columnSet(cols []*metadatapb.ColumnMetadata) map[string]bool {
	m := map[string]bool{}
	for _, col := range cols {
		m[col.Name] = true
	}
	return m
}

// typeName checks a declared type: a %TYPE or %ROWTYPE names a declared
// item, a table or a view, and a type name a built-in, declared or schema
// type.
func (c *plChecker) typeName(t *ast.TypeName) {
	if t == nil {
		return
	}
	parts := typeParts(t)
	if len(parts) == 0 {
		return
	}
	switch {
	case t.IsPercRowtype:
		if len(parts) != 1 {
			return
		}
		if sym := c.lookup(parts[0]); sym != nil {
			return
		}
		if _, found := c.relation(parts[0]); !found && !c.unseen(parts[0]) {
			c.fail("PLS-00201", "identifier '%s' must be declared", parts[0])
		}
		c.read(c.e.resolveSynonymTarget(parts[0], "") + ".*")
	case t.IsPercType:
		if sym := c.lookup(parts[0]); sym != nil {
			if len(parts) == 2 && sym.fields != nil && !sym.fields[parts[1]] {
				c.fail("PLS-00302", "component '%s' must be declared", parts[1])
			}
			return
		}
		if len(parts) != 2 {
			return
		}
		cols, found := c.relation(parts[0])
		if !found {
			if c.knownObject(parts[0]) || c.unseen(parts[0]) {
				return
			}
			c.fail("PLS-00201", "identifier '%s.%s' must be declared", parts[0], parts[1])
			return
		}
		if len(cols) > 0 && !columnSet(cols)[parts[1]] {
			c.fail("PLS-00302", "component '%s' must be declared", parts[1])
		}
		c.read(c.e.resolveSynonymTarget(parts[0], "") + "." + parts[1])
	case len(parts) == 1:
		name := parts[0]
		if plsqlTypes[name] || c.lookup(name) != nil {
			return
		}
		if target := c.e.resolveSynonymTarget(name, ""); c.e.hasType(target) {
			c.readVia(name, target)
			return
		}
		if c.e.units.typesUnknown {
			// A schema loaded from metadata may hold the type.
			return
		}
		c.fail("PLS-00201", "identifier '%s' must be declared", name)
	default:
		// pkg.type or schema.type.
		if c.e.packages[parts[0]] != nil {
			c.read(parts[0])
		}
	}
}

// knownObject reports a schema object of any kind under name.
func (c *plChecker) knownObject(name string) bool {
	return c.e.objectExistsInNamespace(name) || c.e.publicSynonyms[name] != nil
}

// typeFields returns the fields of a record type: a %ROWTYPE's columns,
// or a declared record's, or nil.
func (c *plChecker) typeFields(t *ast.TypeName) map[string]bool {
	parts := typeParts(t)
	if t == nil || len(parts) != 1 {
		return nil
	}
	if t.IsPercRowtype {
		if sym := c.lookup(parts[0]); sym != nil {
			return sym.fields
		}
		target := c.e.resolveSynonymTarget(parts[0], "")
		if v := c.e.views[target]; v != nil {
			return columnSet(v.Columns)
		}
		if tb := c.e.tables[target]; tb != nil {
			return columnSet(tb.Columns)
		}
		return nil
	}
	if sym := c.lookup(parts[0]); sym != nil && sym.kind == plType {
		return sym.fields
	}
	return nil
}

// parseSaved parses a unit's stored text, which Oracle took when it was
// created, for its declarations.
func parseSaved(text string) ast.Node {
	if strings.TrimSpace(text) == "" {
		return nil
	}
	res, err := oracleparser.Parse(text)
	if err != nil || res == nil {
		// A stored unit's text the parser does not read: its
		// declarations are unknown, not missing.
		return nil
	}
	for _, item := range res.Items {
		if raw, ok := item.(*ast.RawStmt); ok {
			return raw.Stmt
		}
		return item
	}
	return nil
}
