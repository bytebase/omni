package catalog

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
)

// sqlFnParams are the input parameters of the SQL function whose body is
// bound. A name no column takes is one of them, bare or qualified with the
// function's name, and $n is the nth.
//
// pg: src/backend/executor/functions.c — sql_fn_post_column_ref, sql_fn_param_ref
type sqlFnParams struct {
	name    string
	procOID uint32 // set once the function has its OID; it names the function in a deparse
	names   []string
	types   []uint32
}

// sqlFnParamRef is a reference bound to a parameter.
type sqlFnParamRef struct {
	num int
	// qualify reports that a range table is in scope, where
	// print_function_sqlbody qualifies a parameter with the function's name.
	qualify bool
}

func (p *sqlFnParams) lookup(name string) int {
	if p == nil || name == "" {
		return 0
	}
	for i, n := range p.names {
		if n == name {
			return i + 1
		}
	}
	return 0
}

// deparse prints a parameter reference as get_parameter does.
//
// pg: src/backend/utils/adt/ruleutils.c — get_parameter
func (p *sqlFnParams) deparse(c *Catalog, ref sqlFnParamRef) string {
	if ref.num < 1 || ref.num > len(p.names) || p.names[ref.num-1] == "" {
		return fmt.Sprintf("$%d", ref.num)
	}
	name := quoteIdentifier(p.names[ref.num-1])
	if !ref.qualify {
		return name
	}
	fname := p.name
	if c != nil {
		if up := c.userProcs[p.procOID]; up != nil {
			fname = up.Name
		}
	}
	return quoteIdentifier(fname) + "." + name
}

func (b *ruleBindings) bindFnParam(n nodes.Node, num int, ctx *ruleEvalContext) {
	qualify := b.fnQualify
	for cur := ctx; cur != nil && !qualify; cur = cur.outer {
		qualify = len(cur.scope) > 0 || cur.hasOpenSource
	}
	b.paramByNode[n] = sqlFnParamRef{num: num, qualify: qualify}
	if t := b.fn.types[num-1]; t != 0 {
		b.exprTypeByNode[n] = t
	}
}

// sqlFnInputParams returns the parameters a SQL function's body can name:
// its IN, INOUT and VARIADIC ones.
func sqlFnInputParams(name string, params *nodes.List, argOIDs []uint32) *sqlFnParams {
	fn := &sqlFnParams{name: name, types: argOIDs}
	if params != nil {
		for _, item := range params.Items {
			fp, ok := item.(*nodes.FunctionParameter)
			if !ok || fp.Mode == nodes.FUNC_PARAM_OUT || fp.Mode == nodes.FUNC_PARAM_TABLE {
				continue
			}
			fn.names = append(fn.names, fp.Name)
		}
	}
	for len(fn.names) < len(fn.types) {
		fn.names = append(fn.names, "")
	}
	return fn
}

// sqlBodyStmts returns the statements of a SQL-standard body: the RETURN
// statement, or those between BEGIN ATOMIC and END.
func sqlBodyStmts(sqlBody nodes.Node) []nodes.Node {
	switch n := sqlBody.(type) {
	case nil:
		return nil
	case *nodes.List:
		stmts := n.Items
		if len(n.Items) == 1 {
			if inner, ok := n.Items[0].(*nodes.List); ok {
				stmts = inner.Items
			}
		}
		return stmts
	default:
		return []nodes.Node{sqlBody}
	}
}

// bindSQLFunctionBody binds each statement of a SQL function's body as parse
// analysis does at CREATE, which refuses a body naming what does not exist.
// A utility statement is not analyzed, and nothing a statement would create
// is there for the statements after it.
//
// pg: src/backend/executor/functions.c — fmgr_sql_validator
// pg: src/backend/commands/functioncmds.c — interpret_AS_clause
func (c *Catalog) bindSQLFunctionBody(stmts []nodes.Node, fn *sqlFnParams) (*ruleBindings, error) {
	b := newRuleBindings(0, 0)
	b.fn = fn
	for _, s := range stmts {
		if raw, ok := s.(*nodes.RawStmt); ok {
			s = raw.Stmt
		}
		var err error
		switch v := s.(type) {
		case *nodes.SelectStmt:
			if v.IntoClause != nil {
				// SELECT INTO is CREATE TABLE AS, a utility statement.
				continue
			}
			b.fnQualify = false
			err = c.bindRuleStmt(b, v, &ruleEvalContext{})
		case *nodes.InsertStmt, *nodes.UpdateStmt, *nodes.DeleteStmt, *nodes.MergeStmt:
			b.fnQualify = true
			err = c.bindRuleStmt(b, v, &ruleEvalContext{})
		case *nodes.ReturnStmt:
			b.fnQualify = false
			err = c.bindRuleExpr(b, v.Returnval, &ruleEvalContext{})
		case *nodes.CallStmt:
			b.fnQualify = false
			err = c.checkCall(b, v.Funccall, &ruleEvalContext{})
		}
		if err != nil {
			return nil, err
		}
	}
	return b, nil
}

// checkSQLFunctionBody runs what fmgr_sql_validator checks of a string body
// when check_function_bodies is on: every statement parses and, unless an
// argument is polymorphic, is analyzed, and the last one returns the
// function's result type.
//
// pg: src/backend/executor/functions.c — fmgr_sql_validator
func (c *Catalog) checkSQLFunctionBody(body string, fn *sqlFnParams, kind byte, retType uint32, outTypes []uint32) error {
	list, err := parser.Parse(body)
	if err != nil {
		var syntax *parser.ParseError
		if errors.As(err, &syntax) {
			code := syntax.Code
			if code == "" {
				code = CodeSyntaxError
			}
			return &Error{Code: code, Message: syntax.Message}
		}
		return err
	}
	for _, t := range fn.types {
		if isPolymorphic(t) {
			return nil
		}
	}
	var stmts []nodes.Node
	if list != nil {
		stmts = list.Items
	}
	b, err := c.bindSQLFunctionBody(stmts, fn)
	if err != nil {
		return err
	}
	return c.checkSQLFunctionResult(b, stmts, kind, retType, outTypes)
}

// checkFunctionBodies reports the session's check_function_bodies, which
// pg_dump turns off so a function may name what is created after it. A
// catalog loaded from a snapshot checks no body unless the session sets it
// on: a name the snapshot does not hold may still exist on the server.
func (c *Catalog) checkFunctionBodies() bool {
	switch strings.Trim(strings.ToLower(c.sessionSettings["check_function_bodies"]), "'\" ") {
	case "off", "false", "0", "no", "f", "n":
		return false
	case "on", "true", "1", "yes", "t", "y":
		return !c.installingSnapshot
	}
	return !c.snapshot && !c.installingSnapshot
}

// checkSQLFunctionResult checks that the last statement of a SQL function
// returns what the function is declared to: one column coercible by
// assignment to a scalar result, or the columns of a row result.
//
// pg: src/backend/executor/functions.c — check_sql_fn_retval
func (c *Catalog) checkSQLFunctionResult(b *ruleBindings, stmts []nodes.Node, kind byte, retType uint32, outTypes []uint32) error {
	if retType == VOIDOID || kind == 'p' || isPolymorphic(retType) {
		return nil
	}
	mismatch := func(detail string) error {
		return &Error{
			Code:    CodeInvalidFunctionDefinition,
			Message: fmt.Sprintf("return type mismatch in function declared to return %s: %s", c.FormatType(retType, -1), detail),
		}
	}
	var last nodes.Node
	for _, s := range stmts {
		if raw, ok := s.(*nodes.RawStmt); ok {
			s = raw.Stmt
		}
		last = s
	}
	var cols []uint32
	switch v := last.(type) {
	case *nodes.SelectStmt:
		if v.IntoClause != nil {
			return mismatch("Function's final statement must be SELECT or INSERT/UPDATE/DELETE/MERGE RETURNING.")
		}
		types, unknown := c.extractRuleQueryColTypes(b, v)
		if unknown {
			return nil
		}
		cols = types
	case *nodes.ReturnStmt:
		cols = []uint32{c.inferRuleExprType(b, v.Returnval, nil)}
	case *nodes.InsertStmt, *nodes.UpdateStmt, *nodes.DeleteStmt, *nodes.MergeStmt:
		returning := dmlReturningList(v)
		if returning == nil || len(returning.Items) == 0 {
			return mismatch("Function's final statement must be SELECT or INSERT/UPDATE/DELETE/MERGE RETURNING.")
		}
		for _, item := range returning.Items {
			rt, ok := item.(*nodes.ResTarget)
			if !ok {
				return nil
			}
			if cr, ok := rt.Val.(*nodes.ColumnRef); ok && isStarColumnRef(cr) {
				return nil
			}
			cols = append(cols, c.inferRuleExprType(b, rt.Val, nil))
		}
	default:
		return mismatch("Function's final statement must be SELECT or INSERT/UPDATE/DELETE/MERGE RETURNING.")
	}
	coercible := func(from, to uint32) bool {
		switch {
		case from == 0 || isPolymorphic(from) || isPolymorphic(to):
			return true // a type the binder does not resolve
		case from == UNKNOWNOID:
			from = TEXTOID // a result column's literal resolves to text
		}
		return from == to || c.CanCoerce(from, to, 'a')
	}
	typ := c.typeByOID[retType]
	if typ == nil {
		return nil
	}
	switch {
	case typ.Type == 'c' || retType == RECORDOID:
		if len(cols) == 1 && kind != 'p' && cols[0] != 0 && coercible(cols[0], retType) {
			return nil
		}
		var want []uint32
		if retType == RECORDOID {
			if len(outTypes) == 0 {
				return nil
			}
			want = outTypes
		} else {
			rel := c.relationByOID[typ.RelID]
			if rel == nil {
				return nil
			}
			for _, col := range rel.Columns {
				want = append(want, col.TypeOID)
			}
		}
		if len(cols) > len(want) {
			return mismatch("Final statement returns too many columns.")
		}
		for i, t := range cols {
			if !coercible(t, want[i]) {
				return mismatch(fmt.Sprintf("Final statement returns %s instead of %s at column %d.", c.FormatType(t, -1), c.FormatType(want[i], -1), i+1))
			}
		}
		if len(cols) < len(want) {
			return mismatch("Final statement returns too few columns.")
		}
	case typ.Type == 'p':
		return nil
	default:
		if len(cols) != 1 {
			return mismatch("Final statement must return exactly one column.")
		}
		if !coercible(cols[0], retType) {
			return mismatch(fmt.Sprintf("Actual return type is %s.", c.FormatType(cols[0], -1)))
		}
	}
	return nil
}

func dmlReturningList(n nodes.Node) *nodes.List {
	switch v := n.(type) {
	case *nodes.InsertStmt:
		return v.ReturningList
	case *nodes.UpdateStmt:
		return v.ReturningList
	case *nodes.DeleteStmt:
		return v.ReturningList
	case *nodes.MergeStmt:
		return v.ReturningList
	}
	return nil
}

// checkCall resolves the procedure a CALL names among the routines on the
// search path. It accepts a call whose argument types it cannot infer.
//
// pg: src/backend/parser/parse_func.c — ParseFuncOrColumn (proc_call)
func (c *Catalog) checkCall(b *ruleBindings, fc *nodes.FuncCall, ctx *ruleEvalContext) error {
	if fc == nil {
		return nil
	}
	var argTypes []uint32
	var argNames []string
	if fc.Args != nil {
		for i, arg := range fc.Args.Items {
			name := ""
			if na, ok := arg.(*nodes.NamedArgExpr); ok && na != nil {
				name, arg = na.Name, na.Arg
			}
			if err := c.bindRuleExpr(b, arg, ctx); err != nil {
				return err
			}
			t := c.inferRuleExprType(b, arg, ctx)
			if t == 0 {
				return nil
			}
			argTypes = append(argTypes, t)
			if name != "" && argNames == nil {
				argNames = make([]string, i, len(fc.Args.Items))
			}
			if argNames != nil {
				argNames = append(argNames, name)
			}
		}
	}
	schemaName, name := qualifiedName(fc.Funcname)
	var searchNS []uint32
	if schemaName != "" {
		ns := c.schemaByName[schemaName]
		if ns == nil {
			return errUndefinedSchema(schemaName)
		}
		searchNS = []uint32{ns.OID}
	} else {
		searchNS = c.routineSearchPathWithCatalog()
	}
	// Every routine visible on the search path is a candidate; one in a
	// later schema with the signature of an earlier one is hidden.
	// pg: src/backend/catalog/namespace.c — FuncnameGetCandidates
	var candidates []*BuiltinProc
	for _, nsOID := range searchNS {
		for _, p := range c.procByName[name] {
			pNS := uint32(PGCatalogNamespace)
			if p.OID >= FirstNormalObjectId {
				up := c.userProcs[p.OID]
				if up == nil || up.Schema == nil {
					continue
				}
				pNS = up.Schema.OID
			}
			if pNS != nsOID || slices.ContainsFunc(candidates, func(q *BuiltinProc) bool { return slices.Equal(q.ArgTypes, p.ArgTypes) }) {
				continue
			}
			candidates = append(candidates, p)
		}
	}
	// CALL passes an argument for each OUT parameter too, which the call
	// signature here leaves out; accept a call to such a procedure.
	for _, p := range candidates {
		if up := c.userProcs[p.OID]; up != nil && p.Kind == 'p' && slices.Contains(up.ArgModes, 'o') {
			return nil
		}
	}
	undefined := &Error{
		Code:    CodeUndefinedFunction,
		Message: fmt.Sprintf("procedure %s(%s) does not exist", name, c.formatArgTypes(argTypes)),
	}
	if len(candidates) == 0 {
		return undefined
	}
	best, _, err := c.resolveFuncOverload(candidates, argTypes, argNames, false, fc.FuncVariadic)
	if err != nil {
		var ce *Error
		if errors.As(err, &ce) && ce.Code == CodeUndefinedFunction {
			return undefined
		}
		// Ambiguity and the rest are left to PostgreSQL.
		return nil
	}
	if best == nil {
		return undefined
	}
	if best.Kind != 'p' {
		return &Error{
			Code:    CodeWrongObjectType,
			Message: fmt.Sprintf("%s(%s) is not a procedure", name, c.formatArgTypes(argTypes)),
		}
	}
	b.addFuncDep(best.OID)
	return nil
}

func (c *Catalog) formatArgTypes(types []uint32) string {
	parts := make([]string, len(types))
	for i, t := range types {
		parts[i] = c.FormatType(t, -1)
	}
	return strings.Join(parts, ", ")
}

// procBody is the body PostgreSQL stores for a routine: a SQL-standard body
// as print_function_sqlbody prints it today, with the names of what it uses
// as they are now; otherwise the text it was given.
func (c *Catalog) procBody(up *UserProc) string {
	if up == nil {
		return ""
	}
	if up.HasSqlBody && up.sqlBody != nil {
		return c.deparseSqlBody(up.sqlBody, up.bodyBindings)
	}
	return up.Body
}

// systemRelationName reports whether rv may name a system catalog or an
// information_schema view, which the catalog does not hold.
func systemRelationName(rv *nodes.RangeVar) bool {
	switch rv.Schemaname {
	case "pg_catalog", "information_schema":
		return true
	case "":
		return strings.HasPrefix(rv.Relname, "pg_")
	}
	return false
}
