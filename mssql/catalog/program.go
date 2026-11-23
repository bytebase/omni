package catalog

import (
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/mssql/ast"
	mssqlparser "github.com/bytebase/omni/mssql/parser"
)

// This file binds the body of a procedure, function or trigger as SQL Server
// does when it creates one, as measured by testdata/programs/simulate.tsv.
//
// SQL Server declares a module's variables as it parses the body, so a
// variable used before its DECLARE is refused. It then binds each statement
// whose tables all exist, and defers a statement that names a table it cannot
// find, or a temporary table, until the module runs: deferred name
// resolution. A user-defined function a body calls is resolved when it runs
// too. An inline table-valued function is bound as a view is.

// programKind is the kind of module whose body is bound.
type programKind int

const (
	programProcedure programKind = iota
	programScalarFunction
	programTableFunction // multi-statement table-valued function
	programTrigger
)

// programBinder binds one module's body.
type programBinder struct {
	*binder
	kind programKind
}

// moduleParams declares a module's parameters.
func moduleParams(b *binder, params *ast.List) {
	if params == nil {
		return
	}
	for _, item := range params.Items {
		if p, ok := item.(*ast.ParamDef); ok {
			b.vars[strings.ToUpper(p.Name)] = true
		}
	}
}

// bindProcedure binds CREATE or ALTER PROCEDURE.
func (c *Catalog) bindProcedure(n *ast.CreateProcedureStmt, source, name string) error {
	if n.Body == nil || n.ExternalName != "" {
		return nil
	}
	pb := c.newProgramBinder(programProcedure, source, name)
	moduleParams(pb.binder, n.Params)
	return pb.body(n.Body)
}

// bindFunction binds CREATE or ALTER FUNCTION.
func (c *Catalog) bindFunction(n *ast.CreateFunctionStmt, source, name string) error {
	if n.ExternalName != "" || n.Body == nil {
		return nil
	}
	kind := programScalarFunction
	if n.ReturnsTable != nil {
		if n.ReturnsTable.Variable == "" {
			return c.bindInlineFunction(n, source, name)
		}
		kind = programTableFunction
	}
	pb := c.newProgramBinder(kind, source, name)
	pb.schemaBound = hasOptionWord(n.Options, "SCHEMABINDING")
	moduleParams(pb.binder, n.Params)
	// A recursive call binds against the signature being created, not the
	// one an ALTER replaces.
	schemaName, _ := tableRefSchemaObject(n.Name, c.defaultSchema)
	pb.self = &selfFunction{schema: schemaName, name: name, info: fnInfo{category: functionCategory(n)}}
	if n.Params != nil {
		pb.self.info.params = len(n.Params.Items)
	}
	if kind == programTableFunction {
		if err := returnTableColumns(n.ReturnsTable); err != nil {
			return err
		}
		v := strings.ToUpper(n.ReturnsTable.Variable)
		pb.vars[v] = true
		pb.tableVars[v] = columnDefs(n.ReturnsTable.Columns)
	}
	if err := pb.body(n.Body); err != nil {
		return err
	}
	return pb.lastIsReturn(n.Body)
}

// returnTableColumns checks the table a multi-statement function returns.
func returnTableColumns(rt *ast.ReturnsTableDef) error {
	if rt.Columns == nil {
		return nil
	}
	var names []string
	for i, item := range rt.Columns.Items {
		cd, ok := item.(*ast.ColumnDef)
		if !ok {
			continue
		}
		if cd.Computed == nil && cd.DataType != nil && cd.DataType.Schema == "" && !isValidDataType(cd.DataType, false) {
			return mssqlErr(ErrUnknownColumnType, "Column, parameter, or variable #%d: Cannot find data type %s.", i+1, cd.DataType.Name)
		}
		for _, prev := range names {
			if strings.EqualFold(prev, cd.Name) {
				return mssqlErr(ErrDupColumn, "Column names in each table must be unique. Column name '%s' in table '%s' is specified more than once.", cd.Name, rt.Variable)
			}
		}
		names = append(names, cd.Name)
	}
	return nil
}

func columnDefs(cols *ast.List) []qcol {
	if cols == nil {
		return nil
	}
	out := []qcol{}
	for _, item := range cols.Items {
		if cd, ok := item.(*ast.ColumnDef); ok {
			out = append(out, qcol{name: cd.Name, typ: dataTypeClass(cd.DataType), identity: cd.Identity != nil})
		}
	}
	return out
}

// bindInlineFunction binds an inline table-valued function's query as a
// view's is, with the function's parameters declared.
func (c *Catalog) bindInlineFunction(n *ast.CreateFunctionStmt, source, name string) error {
	var sel *ast.SelectStmt
	ast.Inspect(n.Body, func(node ast.Node) bool {
		if sel != nil {
			return false
		}
		if s, ok := node.(*ast.SelectStmt); ok {
			sel = s
			return false
		}
		return true
	})
	if sel == nil {
		return nil
	}
	b := newBinder(c, source, name, 0)
	b.inView = true
	b.schemaBound = hasOptionWord(n.Options, "SCHEMABINDING")
	moduleParams(b, n.Params)
	if err := b.viewClauses(sel); err != nil {
		return err
	}
	if err := b.variables(sel); err != nil {
		return err
	}
	cols, opaque, err := b.query(sel, nil)
	if err != nil || opaque {
		return err
	}
	for i, col := range cols {
		if col.name == "" {
			return mssqlErr(4514, "CREATE FUNCTION failed because a column name is not specified for column %d.", i+1)
		}
		for _, prev := range cols[:i] {
			if strings.EqualFold(prev.name, col.name) {
				return mssqlErr(4506, "Column names in each view or function must be unique. Column name '%s' in view or function '%s' is specified more than once.", col.name, name)
			}
		}
	}
	return nil
}

// bindTrigger binds a trigger's body: a DML trigger's reads inserted and
// deleted, which hold the columns of the table or view it is on.
func (c *Catalog) bindTrigger(n *ast.CreateTriggerStmt, source, name string, parent *qrel) error {
	if n.Body == nil || n.ExternalName != "" {
		return nil
	}
	pb := c.newProgramBinder(programTrigger, source, name)
	if parent != nil {
		pb.pseudo = map[string]*qrel{"INSERTED": parent, "DELETED": parent}
	}
	return pb.body(n.Body)
}

func (c *Catalog) newProgramBinder(kind programKind, source, name string) *programBinder {
	b := newBinder(c, source, name, 0)
	b.deferNames = true
	b.tableVars = map[string][]qcol{}
	return &programBinder{binder: b, kind: kind}
}

// body binds a module's body: its variables first, then its statements.
func (pb *programBinder) body(body ast.Node) error {
	if err := pb.declared(body); err != nil {
		return err
	}
	return pb.stmt(body)
}

// declared refuses a variable the body uses before it declares it, in the
// order SQL Server's parser reads the body.
func (pb *programBinder) declared(body ast.Node) error {
	loc := ast.NodeLoc(body)
	if loc.Start < 0 || loc.End > len(pb.source) || loc.Start >= loc.End {
		return nil
	}
	declared := map[int]bool{}
	named := map[int]bool{}
	tableUse := map[int]bool{}
	ast.Inspect(body, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.VariableDecl:
			declared[v.Loc.Start] = true
		case *ast.ExecArg:
			if v.Name != "" {
				named[v.Loc.Start] = true
			}
		case *ast.TableVarRef:
			tableUse[v.Loc.Start] = true
		}
		return true
	})
	vars := map[string]bool{}
	for v := range pb.vars {
		vars[v] = true
	}
	for _, tok := range mssqlparser.Tokenize(pb.source[loc.Start:loc.End]) {
		if tok.Type != mssqlparser.VARIABLE {
			continue
		}
		at := loc.Start + tok.Loc
		name := strings.ToUpper(tok.Str)
		switch {
		case declared[at]:
			vars[name] = true
		case named[at], vars[name]:
		case tableUse[at]:
			return mssqlErr(1087, "Must declare the table variable \"%s\".", tok.Str)
		default:
			return mssqlErr(137, "Must declare the scalar variable \"%s\".", tok.Str)
		}
	}
	for v := range vars {
		pb.vars[v] = true
	}
	return nil
}

// stmt binds one statement of the body, and the statements it holds.
func (pb *programBinder) stmt(n ast.Node) error {
	if n == nil || isNilNode(n) {
		return nil
	}
	switch s := n.(type) {
	case *ast.BeginEndStmt:
		return pb.list(s.Stmts)
	case *ast.List:
		return pb.list(s)
	case *ast.TryCatchStmt:
		if err := pb.list(s.TryBlock); err != nil {
			return err
		}
		return pb.list(s.CatchBlock)
	case *ast.IfStmt:
		if err := pb.expr(s.Condition); err != nil {
			return err
		}
		if err := pb.stmt(s.Then); err != nil {
			return err
		}
		return pb.stmt(s.Else)
	case *ast.WhileStmt:
		if err := pb.expr(s.Condition); err != nil {
			return err
		}
		return pb.stmt(s.Body)
	case *ast.ReturnStmt:
		if s.Value != nil && pb.kind == programTableFunction {
			return mssqlErr(178, "A RETURN statement with a return value cannot be used in this context.")
		}
		return pb.expr(s.Value)
	case *ast.SetStmt:
		if s.Cursor != nil {
			return pb.stmt(s.Cursor)
		}
		return pb.expr(s.Value)
	case *ast.DeclareStmt:
		if s.Variables == nil {
			return nil
		}
		for _, item := range s.Variables.Items {
			v, ok := item.(*ast.VariableDecl)
			if !ok {
				continue
			}
			if v.IsTable {
				pb.tableVars[strings.ToUpper(v.Name)] = columnDefs(v.TableDef)
			}
			if err := pb.expr(v.Default); err != nil {
				return err
			}
		}
		return nil
	case *ast.DeclareCursorStmt:
		return pb.query(s.Query)
	case *ast.SelectStmt:
		if pb.function() && s.IntoTable == nil && returnsData(s) {
			return mssqlErr(444, "Select statements included within a function cannot return data to a client.")
		}
		return pb.query(s)
	case *ast.InsertStmt:
		if err := pb.sideEffect(s.Relation, "INSERT", nil); err != nil {
			return err
		}
		if pb.deferred(s) {
			return nil
		}
		return pb.insert(s)
	case *ast.UpdateStmt:
		if err := pb.sideEffect(s.Relation, "UPDATE", s.FromClause); err != nil {
			return err
		}
		if pb.deferred(s) {
			return nil
		}
		return pb.update(s)
	case *ast.DeleteStmt:
		if err := pb.sideEffect(s.Relation, "DELETE", s.FromClause); err != nil {
			return err
		}
		if pb.deferred(s) {
			return nil
		}
		return pb.delete(s)
	case *ast.MergeStmt:
		if err := pb.sideEffect(s.Target, "MERGE", nil); err != nil {
			return err
		}
		if pb.deferred(s) {
			return nil
		}
		return pb.merge(s)
	case *ast.TruncateStmt:
		if pb.function() {
			return mssqlErr(443, "Invalid use of a side-effecting operator 'TRUNCATE TABLE' within a function.")
		}
	}
	return nil
}

func (pb *programBinder) list(l *ast.List) error {
	if l == nil {
		return nil
	}
	for _, item := range l.Items {
		if err := pb.stmt(item); err != nil {
			return err
		}
	}
	return nil
}

func (pb *programBinder) function() bool {
	return pb.kind == programScalarFunction || pb.kind == programTableFunction
}

// sideEffect refuses, in a function, a statement that writes anything but a
// table variable.
func (pb *programBinder) sideEffect(target ast.Node, verb string, from *ast.List) error {
	if !pb.function() {
		return nil
	}
	if _, ok := target.(*ast.TableVarRef); ok {
		return nil
	}
	// UPDATE r ... FROM @r AS r writes the table variable its alias names.
	if ref, ok := target.(*ast.TableRef); ok && ref.Schema == "" && ref.Database == "" && from != nil {
		aliased := false
		ast.Inspect(from, func(n ast.Node) bool {
			if v, ok := n.(*ast.TableVarRef); ok && strings.EqualFold(v.Alias, ref.Object) {
				aliased = true
			}
			return !aliased
		})
		if aliased {
			return nil
		}
	}
	return mssqlErr(443, "Invalid use of a side-effecting operator '%s' within a function.", verb)
}

// returnsData reports whether a SELECT sends rows to the client, rather than
// assigning variables.
func returnsData(s *ast.SelectStmt) bool {
	if s.Op != ast.SetOpNone || s.TargetList == nil {
		return true
	}
	for _, item := range s.TargetList.Items {
		if rt, ok := item.(*ast.ResTarget); ok {
			if _, assign := rt.Val.(*ast.SelectAssign); assign {
				continue
			}
			return true
		}
		if _, assign := item.(*ast.SelectAssign); !assign {
			return true
		}
	}
	return false
}

// lastIsReturn holds a scalar or multi-statement function to ending its body
// with RETURN.
func (pb *programBinder) lastIsReturn(body ast.Node) error {
	be, ok := body.(*ast.BeginEndStmt)
	if !ok || be.Stmts == nil || len(be.Stmts.Items) == 0 {
		return nil
	}
	if _, ok := be.Stmts.Items[len(be.Stmts.Items)-1].(*ast.ReturnStmt); ok {
		return nil
	}
	return mssqlErr(455, "The last statement included within a function must be a return statement.")
}

func emptyScope() *qscope {
	return &qscope{block: &qblock{windows: map[string]*ast.WindowDef{}}}
}

// expr binds an expression of the body unless it names a table SQL Server
// defers.
func (pb *programBinder) expr(e ast.Node) error {
	if e == nil || isNilNode(e) {
		return nil
	}
	if pb.function() {
		if err := sideEffectCall(e); err != nil {
			return err
		}
	}
	if pb.deferred(e) {
		return nil
	}
	_, err := pb.binder.expr(e, emptyScope(), ectx{clause: clOther})
	return err
}

// query binds a query of the body unless SQL Server defers it.
func (pb *programBinder) query(q ast.Node) error {
	sel, ok := q.(*ast.SelectStmt)
	if !ok || sel == nil {
		return nil
	}
	if pb.function() {
		if err := sideEffectCall(sel); err != nil {
			return err
		}
	}
	if pb.deferred(sel) {
		return nil
	}
	_, _, err := pb.binder.query(sel, nil)
	return err
}

// sideEffectCall refuses a side-effecting built-in in a function.
func sideEffectCall(n ast.Node) error {
	var err error
	ast.Inspect(n, func(node ast.Node) bool {
		if err != nil {
			return false
		}
		if f, ok := node.(*ast.FuncCallExpr); ok && f.Name != nil && f.Name.Schema == "" {
			switch strings.ToUpper(f.Name.Object) {
			case "NEWID", "NEWSEQUENTIALID":
				err = mssqlErr(443, "Invalid use of a side-effecting operator '%s' within a function.", strings.ToLower(f.Name.Object))
			}
		}
		return err == nil
	})
	return err
}

// deferred reports whether SQL Server defers binding n: it names a table
// that does not exist, or a temporary table.
func (pb *programBinder) deferred(n ast.Node) bool {
	ctes := map[string]bool{}
	ast.Inspect(n, func(node ast.Node) bool {
		if cte, ok := node.(*ast.CommonTableExpr); ok {
			ctes[strings.ToUpper(cte.Name)] = true
		}
		return true
	})
	missing := false
	var source func(te ast.Node)
	source = func(te ast.Node) {
		switch t := te.(type) {
		case *ast.TableRef:
			if pb.missingObject(t, ctes) {
				missing = true
			}
		case *ast.AliasedTableRef:
			source(t.Table)
		case *ast.JoinClause:
			source(t.Left)
			source(t.Right)
		case *ast.PivotExpr:
			source(t.Source)
		case *ast.UnpivotExpr:
			source(t.Source)
		case *ast.FuncCallExpr:
			if t.Name != nil && (t.Name.Schema != "" || !builtinTableFunctions[strings.ToUpper(t.Name.Object)]) && pb.missingObject(t.Name, ctes) {
				missing = true
			}
		}
	}
	from := func(l *ast.List) {
		if l == nil {
			return
		}
		for _, item := range l.Items {
			source(item)
		}
	}
	ast.Inspect(n, func(node ast.Node) bool {
		if missing {
			return false
		}
		switch v := node.(type) {
		case *ast.SelectStmt:
			from(v.FromClause)
		case *ast.InsertStmt:
			source(v.Relation)
		case *ast.UpdateStmt:
			source(v.Relation)
			from(v.FromClause)
		case *ast.DeleteStmt:
			source(v.Relation)
			from(v.FromClause)
		case *ast.MergeStmt:
			source(v.Target)
			source(v.Source)
		}
		return !missing
	})
	return missing
}

// missingObject reports whether a table name does not resolve now.
func (pb *programBinder) missingObject(ref *ast.TableRef, ctes map[string]bool) bool {
	if strings.HasPrefix(ref.Object, "#") {
		return true
	}
	if !pb.localObject(ref) || systemSchema(ref.Schema) {
		return false
	}
	if ref.Schema == "" {
		if ctes[strings.ToUpper(ref.Object)] {
			return false
		}
		if _, ok := pb.pseudo[strings.ToUpper(ref.Object)]; ok {
			return false
		}
	}
	schema := pb.objectSchema(ref.Schema, ref.Object)
	return schema == nil || pb.c.lookupSchemaObject(schema, ref.Object, false) == kindNone
}

// cteScope binds a DML statement's WITH clause into the scope its table
// sources resolve names in.
func (pb *programBinder) cteScope(wc *ast.WithClause) (*qscope, error) {
	scope := &qscope{}
	if wc != nil {
		if err := pb.ctes(wc, scope); err != nil {
			return nil, err
		}
	}
	return scope, nil
}

// dmlScope binds a DML statement's WITH clause, target and FROM clause and
// returns the scope its expressions see and the target's relation.
func (pb *programBinder) dmlScope(wc *ast.WithClause, target ast.Node, fromClause *ast.List) (*qscope, *qrel, error) {
	ctes, err := pb.cteScope(wc)
	if err != nil {
		return nil, nil, err
	}
	sc := emptyScope()
	sc.parent = ctes
	if fromClause != nil {
		for _, item := range fromClause.Items {
			rels, err := pb.tableExpr(item, sc, nil)
			if err != nil {
				return nil, nil, err
			}
			if err := checkExposedNames(sc.rels, rels); err != nil {
				return nil, nil, err
			}
			sc.rels = append(sc.rels, rels...)
		}
	}
	// The target may name a table source of the FROM clause by its alias.
	if ref, ok := target.(*ast.TableRef); ok && ref.Schema == "" && ref.Database == "" {
		for _, r := range sc.rels {
			if strings.EqualFold(r.name, ref.Object) {
				return sc, r, nil
			}
		}
	}
	rels, err := pb.tableExpr(target, sc, nil)
	if err != nil {
		return nil, nil, err
	}
	if len(rels) != 1 {
		return sc, &qrel{opaque: true}, nil
	}
	// A target the FROM clause does not name is one more table its
	// expressions see: UPDATE t SET ... FROM u WHERE t.id = u.n.
	sc.rels = append(sc.rels, rels[0])
	return sc, rels[0], nil
}

// insertColumns is how many columns an INSERT without a column list fills.
func insertColumns(r *qrel) int {
	n := 0
	for _, c := range r.cols {
		if !c.identity && !c.noInsert {
			n++
		}
	}
	return n
}

func targetColumn(r *qrel, name string) error {
	if r.opaque || r.col(name) != nil {
		return nil
	}
	return mssqlErr(ErrInvalidColumnName, "Invalid column name '%s'.", name)
}

func (pb *programBinder) insert(s *ast.InsertStmt) error {
	sc, target, err := pb.dmlScope(s.WithClause, s.Relation, nil)
	if err != nil {
		return err
	}
	width := -1
	if cols := stringNodeList(s.Cols); len(cols) > 0 {
		for _, col := range cols {
			if err := targetColumn(target, col); err != nil {
				return err
			}
		}
		width = len(cols)
	} else if !target.opaque {
		width = insertColumns(target)
	}
	switch src := s.Source.(type) {
	case *ast.ValuesClause:
		if src.Rows == nil {
			return nil
		}
		for _, row := range src.Rows.Items {
			list, ok := row.(*ast.List)
			if !ok {
				continue
			}
			if width >= 0 && len(list.Items) > width {
				return mssqlErr(110, "There are fewer columns in the INSERT statement than values specified in the VALUES clause. The number of values in the VALUES clause must match the number of columns specified in the INSERT statement.")
			}
			if width >= 0 && len(list.Items) < width {
				return mssqlErr(109, "There are more columns in the INSERT statement than values specified in the VALUES clause. The number of values in the VALUES clause must match the number of columns specified in the INSERT statement.")
			}
			if _, err := pb.binder.expr(list, emptyScope(), ectx{clause: clOther}); err != nil {
				return err
			}
		}
	case *ast.SelectStmt:
		cols, opaque, err := pb.binder.query(src, sc.parent)
		if err != nil {
			return err
		}
		if width >= 0 && !opaque {
			if len(cols) < width {
				return mssqlErr(120, "The select list for the INSERT statement contains fewer items than the insert list. The number of SELECT values must match the number of INSERT columns.")
			}
			if len(cols) > width {
				return mssqlErr(121, "The select list for the INSERT statement contains more items than the insert list. The number of SELECT values must match the number of INSERT columns.")
			}
		}
	}
	return nil
}

func (pb *programBinder) setClause(list *ast.List, target *qrel, sc *qscope) error {
	if list == nil {
		return nil
	}
	for _, item := range list.Items {
		set, ok := item.(*ast.SetExpr)
		if !ok {
			continue
		}
		if set.Column != nil {
			if err := targetColumn(target, set.Column.Column); err != nil {
				return err
			}
		}
		if _, err := pb.binder.expr(set.Value, sc, ectx{clause: clOther}); err != nil {
			return err
		}
	}
	return nil
}

func (pb *programBinder) update(s *ast.UpdateStmt) error {
	sc, target, err := pb.dmlScope(s.WithClause, s.Relation, s.FromClause)
	if err != nil {
		return err
	}
	if err := pb.setClause(s.SetClause, target, sc); err != nil {
		return err
	}
	_, err = pb.binder.expr(s.WhereClause, sc, ectx{clause: clWhere, block: sc.block})
	return err
}

func (pb *programBinder) delete(s *ast.DeleteStmt) error {
	sc, _, err := pb.dmlScope(s.WithClause, s.Relation, s.FromClause)
	if err != nil {
		return err
	}
	_, err = pb.binder.expr(s.WhereClause, sc, ectx{clause: clWhere, block: sc.block})
	return err
}

func (pb *programBinder) merge(s *ast.MergeStmt) error {
	ctes, err := pb.cteScope(s.WithClause)
	if err != nil {
		return err
	}
	sc := emptyScope()
	sc.parent = ctes
	targets, err := pb.tableExpr(s.Target, sc, nil)
	if err != nil {
		return err
	}
	sources, err := pb.tableExpr(s.Source, sc, nil)
	if err != nil {
		return err
	}
	if err := checkExposedNames(targets, sources); err != nil {
		return err
	}
	target := &qrel{opaque: true}
	if len(targets) == 1 {
		target = targets[0]
	}
	sc.rels = append(append([]*qrel(nil), targets...), sources...)
	x := ectx{clause: clOn, block: sc.block}
	if _, err := pb.binder.expr(s.OnCondition, sc, x); err != nil {
		return err
	}
	if s.WhenClauses == nil {
		return nil
	}
	for _, item := range s.WhenClauses.Items {
		w, ok := item.(*ast.MergeWhenClause)
		if !ok {
			continue
		}
		if _, err := pb.binder.expr(w.Condition, sc, x); err != nil {
			return err
		}
		switch a := w.Action.(type) {
		case *ast.MergeUpdateAction:
			if err := pb.setClause(a.SetClause, target, sc); err != nil {
				return err
			}
		case *ast.MergeInsertAction:
			for _, col := range stringNodeList(a.Cols) {
				if err := targetColumn(target, col); err != nil {
					return err
				}
			}
			in := &qscope{rels: sources, block: sc.block}
			if _, err := pb.binder.expr(a.Values, in, ectx{clause: clOther}); err != nil {
				return err
			}
		}
	}
	return nil
}

// triggerParent is the relation a DML trigger's inserted and deleted tables
// copy: the table or view the trigger is on.
func (c *Catalog) triggerParent(schema *metadata.SchemaMetadata, name string) *qrel {
	if t := c.resolveTableInSchema(schema, name); t != nil {
		return &qrel{cols: tableColumns(t)}
	}
	if idx, ok := resolveExactOrFoldIndex(schema.Views, func(v *metadata.ViewMetadata) string { return v.Name }, name); ok {
		return newBinder(c, "", name, 0).viewRel(schema.Views[idx], schema.Name)
	}
	return &qrel{opaque: true}
}
