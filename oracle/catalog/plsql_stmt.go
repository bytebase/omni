package catalog

import (
	"strings"

	"github.com/bytebase/omni/oracle/ast"
)

// block compiles a block: its declarations, statements and handlers.
func (c *plChecker) block(n ast.StmtNode) {
	b, ok := n.(*ast.PLSQLBlock)
	if !ok {
		c.stmt(n)
		return
	}
	c.push()
	defer c.pop()
	c.declarations(b.Declarations)
	c.stmts(b.Statements)
	for _, item := range listItems(b.Exceptions) {
		h, ok := item.(*ast.ExceptionHandler)
		if !ok {
			continue
		}
		for _, x := range listItems(h.Exceptions) {
			if s, ok := x.(*ast.String); ok {
				c.exception(s.Str)
			}
		}
		c.stmts(h.Statements)
	}
}

func (c *plChecker) stmts(list *ast.List) {
	for _, item := range listItems(list) {
		if s, ok := item.(ast.StmtNode); ok {
			c.stmt(s)
		}
	}
}

func (c *plChecker) stmtSlice(list []ast.StmtNode) {
	for _, s := range list {
		c.stmt(s)
	}
}

// exception checks a name a handler or RAISE names.
func (c *plChecker) exception(name string) {
	if plsqlExceptions[name] {
		return
	}
	if sym := c.lookup(name); sym != nil {
		return
	}
	if strings.Contains(name, ".") {
		c.qualifiedName(strings.Split(name, "."))
		return
	}
	c.fail("PLS-00201", "identifier '%s' must be declared", name)
}

func (c *plChecker) stmt(n ast.StmtNode) {
	switch s := n.(type) {
	case nil:
	case *ast.PLSQLBlock:
		c.block(s)
	case *ast.PLSQLLabeledStatement:
		c.stmt(s.Statement)
	case *ast.PLSQLNull, *ast.PLSQLPragma, *ast.CommitStmt, *ast.RollbackStmt, *ast.SavepointStmt:
	case *ast.PLSQLAssign:
		c.expr(s.Value)
		c.target(s.Target)
	case *ast.PLSQLIf:
		c.expr(s.Condition)
		c.stmts(s.Then)
		for _, item := range listItems(s.ElsIfs) {
			if e, ok := item.(*ast.PLSQLElsIf); ok {
				c.expr(e.Condition)
				c.stmts(e.Then)
			}
		}
		c.stmts(s.Else)
	case *ast.PLSQLCase:
		c.expr(s.Expr)
		for _, w := range s.Whens {
			c.expr(w.Expr)
			c.stmtSlice(w.Stmts)
		}
		c.stmtSlice(s.Else)
	case *ast.PLSQLLoop:
		c.loop(s)
	case *ast.PLSQLPipeRow:
		c.expr(s.Row)
	case *ast.PLSQLExit:
		c.expr(s.Condition)
	case *ast.PLSQLContinue:
		c.expr(s.Condition)
	case *ast.PLSQLReturn:
		switch {
		case s.Expr != nil && !c.function:
			c.fail("PLS-00372", "In a procedure, RETURN statement cannot contain an expression")
		case s.Expr == nil && c.function && !c.pipelined:
			c.fail("PLS-00503", "RETURN <value> statement required for this return from function")
		default:
			c.expr(s.Expr)
		}
	case *ast.PLSQLGoto:
		if !c.labels[s.Label] {
			c.fail("PLS-00201", "identifier '%s' must be declared", s.Label)
		}
	case *ast.PLSQLRaise:
		if s.Exception != "" {
			c.exception(s.Exception)
		}
	case *ast.PLSQLCall:
		c.call(s.Name)
	case *ast.PLSQLExecImmediate:
		c.expr(s.SQL)
		c.exprList(s.Using)
		c.targets(s.Into)
	case *ast.PLSQLOpen:
		c.cursor(s.Cursor)
		c.exprList(s.Args)
		if s.ForQuery != nil {
			c.sql(s.ForQuery)
		}
		c.expr(s.ForExpr)
	case *ast.PLSQLFetch:
		c.cursor(s.Cursor)
		c.targets(s.Into)
	case *ast.PLSQLClose:
		c.cursor(s.Cursor)
	case *ast.PLSQLForall:
		c.expr(s.Lower)
		c.expr(s.Upper)
		c.push()
		c.declare(s.Index, &plSym{kind: plInParam})
		c.stmt(s.Body)
		c.pop()
	case *ast.SelectStmt:
		if s.IntoVars == nil {
			c.fail("PLS-00428", "an INTO clause is expected in this SELECT statement")
			return
		}
		c.sql(s)
	case *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.MergeStmt:
		c.sql(s)
	}
}

// cursor checks a cursor or cursor variable OPEN, FETCH or CLOSE names. A
// qualified one is a package's, which the checker takes as declared.
func (c *plChecker) cursor(name string) {
	if strings.ContainsAny(name, ".:") || c.lookup(name) != nil {
		return
	}
	c.fail("PLS-00201", "identifier '%s' must be declared", name)
}

func (c *plChecker) loop(s *ast.PLSQLLoop) {
	c.push()
	defer c.pop()
	if s.Iteration != nil {
		// FOR i, v IN PAIRS OF xs: every iterand is the loop's own.
		for _, item := range listItems(s.Iteration.Iterands) {
			if it, ok := item.(*ast.PLSQLIterand); ok {
				c.declare(it.Name, &plSym{kind: plInParam})
			}
		}
	}
	switch s.Type {
	case ast.LOOP_WHILE:
		c.expr(s.Condition)
	case ast.LOOP_FOR:
		c.expr(s.LowerBound)
		c.expr(s.UpperBound)
		if s.Iterator != "" {
			c.declare(s.Iterator, &plSym{kind: plInParam})
		}
	case ast.LOOP_CURSOR_FOR:
		var fields map[string]bool
		switch {
		case s.CursorQuery != nil:
			fields = c.sql(s.CursorQuery)
		case s.CursorName != "":
			sym := c.lookup(s.CursorName)
			switch {
			case sym != nil:
				fields = sym.fields
			case s.Iteration == nil:
				c.fail("PLS-00201", "identifier '%s' must be declared", s.CursorName)
			default:
				// pk.c: the parser keeps the qualified cursor in the
				// iterator, which the checker takes as declared.
			}
			c.exprList(s.CursorArgs)
		}
		if s.Iterator != "" {
			c.declare(s.Iterator, &plSym{kind: plInParam, fields: fields})
		}
	}
	c.stmts(s.Statements)
}

func (c *plChecker) exprList(list *ast.List) {
	for _, item := range listItems(list) {
		if x, ok := item.(ast.ExprNode); ok {
			c.expr(x)
		}
	}
}

func (c *plChecker) targets(list *ast.List) {
	for _, item := range listItems(list) {
		if x, ok := item.(ast.ExprNode); ok {
			c.target(x)
		}
	}
}

// target checks an assignment's target: a variable, an OUT parameter, a
// record field, or a trigger's :NEW field.
func (c *plChecker) target(x ast.ExprNode) {
	switch t := x.(type) {
	case *ast.ColumnRef:
		if t.Table == "" {
			sym := c.lookup(t.Column)
			switch {
			case sym == nil:
				c.fail("PLS-00201", "identifier '%s' must be declared", t.Column)
			case sym.kind == plInParam || sym.kind == plConstant:
				c.fail("PLS-00363", "expression '%s' cannot be used as an assignment target", t.Column)
			case sym.duplicate:
				c.fail("PLS-00371", "at most one declaration for '%s' is permitted", t.Column)
			}
			return
		}
		c.columnRef(t)
	case *ast.BindVariable:
		c.bind(t)
	default:
		c.expr(x)
	}
}

// bind checks a trigger's :NEW.field or :OLD.field: the name is one the
// trigger reads its row by, and the field a column of the table.
func (c *plChecker) bind(b *ast.BindVariable) {
	if c.row == nil || b.Member == "" {
		return
	}
	name := strings.ToUpper(b.Name)
	cols := c.row[name]
	if cols == nil || !cols[b.Member] {
		c.fail("PLS-00049", "bad bind variable '%s.%s'", name, b.Member)
		return
	}
	c.read(c.rowTable + "." + b.Member)
}

// call checks a procedure call statement.
func (c *plChecker) call(x ast.ExprNode) {
	var parts []string
	var args *ast.List
	switch n := x.(type) {
	case *ast.FuncCallExpr:
		if n.FuncName == nil {
			return
		}
		if n.FuncName.Schema != "" {
			parts = append(parts, n.FuncName.Schema)
		}
		parts = append(parts, n.FuncName.Name)
		args = n.Args
		c.exprList(n.Args)
	case *ast.ColumnRef:
		if n.Schema != "" {
			parts = append(parts, n.Schema)
		}
		if n.Table != "" {
			parts = append(parts, n.Table)
		}
		parts = append(parts, n.Column)
	default:
		c.expr(x)
		return
	}
	sigs, known := c.resolveCall(parts)
	if known {
		c.checkArgs(parts[len(parts)-1], sigs, args, false)
	}
}

// checkArgs checks a call's arguments against a subprogram's overloads:
// one must take them, by count and by name.
func (c *plChecker) checkArgs(name string, sigs []plSignature, args *ast.List, function bool) {
	if len(sigs) == 0 {
		return
	}
	any := false
	for _, sig := range sigs {
		if sig.function != function {
			continue
		}
		any = true
		if argsFit(sig.params, args) {
			return
		}
	}
	if !any {
		if function {
			c.fail("PLS-00222", "no function with name '%s' exists in this scope", name)
		} else {
			c.fail("PLS-00221", "'%s' is not a procedure or is undefined", name)
		}
		return
	}
	c.fail("PLS-00306", "wrong number or types of arguments in call to '%s'", name)
}

// argsFit reports whether args can be passed to params: each named
// argument names a parameter, no more positional arguments than
// parameters, and every parameter without a default is given.
func argsFit(params []*ast.Parameter, args *ast.List) bool {
	given := map[string]bool{}
	positional := 0
	for _, item := range listItems(args) {
		if na, ok := item.(*ast.NamedArgExpr); ok {
			found := false
			for _, p := range params {
				if p.Name == na.Name {
					found = true
				}
			}
			if !found {
				return false
			}
			given[na.Name] = true
			continue
		}
		positional++
	}
	if positional > len(params) {
		return false
	}
	for i, p := range params {
		if i < positional {
			continue
		}
		if !given[p.Name] && p.Default == nil {
			return false
		}
	}
	return true
}
