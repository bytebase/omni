package ddl

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// --- View AST helpers (reverse conversion for testing) ---

// reverseConvertSelectStmt converts a legacy test SelectStmt back to a nodes.SelectStmt.
// Used by test helpers to build pgparser AST from legacy test types.
func reverseConvertSelectStmt(s *SelectStmt) *nodes.SelectStmt {
	if s == nil {
		return nil
	}

	// Set operation.
	if s.Op != SetOpNone {
		ns := &nodes.SelectStmt{
			Larg: reverseConvertSelectStmt(s.Left),
			Rarg: reverseConvertSelectStmt(s.Right),
		}
		switch s.Op {
		case SetOpUnion:
			ns.Op = nodes.SETOP_UNION
		case SetOpUnionAll:
			ns.Op = nodes.SETOP_UNION
			ns.All = true
		case SetOpIntersect:
			ns.Op = nodes.SETOP_INTERSECT
		case SetOpIntersectAll:
			ns.Op = nodes.SETOP_INTERSECT
			ns.All = true
		case SetOpExcept:
			ns.Op = nodes.SETOP_EXCEPT
		case SetOpExceptAll:
			ns.Op = nodes.SETOP_EXCEPT
			ns.All = true
		}
		return ns
	}

	// Simple SELECT.
	ns := &nodes.SelectStmt{}

	// Target list.
	if len(s.TargetList) > 0 {
		items := make([]nodes.Node, len(s.TargetList))
		for i, rt := range s.TargetList {
			items[i] = &nodes.ResTarget{
				Name: rt.Name,
				Val:  reverseConvertExpr(rt.Val),
			}
		}
		ns.TargetList = &nodes.List{Items: items}
	}

	// FROM clause.
	if len(s.From) > 0 {
		items := make([]nodes.Node, len(s.From))
		for i, fi := range s.From {
			items[i] = reverseConvertFromItem(fi)
		}
		ns.FromClause = &nodes.List{Items: items}
	}

	return ns
}

// reverseConvertExpr converts a catalog Expr back to a nodes.Node.
func reverseConvertExpr(e *Expr) nodes.Node {
	if e == nil {
		return nil
	}
	switch e.Kind {
	case ExprColumnRef:
		var items []nodes.Node
		if e.TableName != "" {
			items = append(items, &nodes.String{Str: e.TableName})
		}
		items = append(items, &nodes.String{Str: e.ColumnName})
		return &nodes.ColumnRef{Fields: &nodes.List{Items: items}}

	case ExprStar:
		var items []nodes.Node
		if e.StarTable != "" {
			items = append(items, &nodes.String{Str: e.StarTable})
		}
		items = append(items, &nodes.A_Star{})
		return &nodes.ColumnRef{Fields: &nodes.List{Items: items}}

	case ExprLiteral:
		ac := &nodes.A_Const{}
		switch e.LiteralType {
		case INT4OID:
			ac.Val = &nodes.Integer{Ival: 0}
		case FLOAT8OID:
			ac.Val = &nodes.Float{Fval: "0"}
		case BOOLOID:
			ac.Val = &nodes.Boolean{Boolval: false}
		default:
			ac.Val = &nodes.String{Str: ""}
		}
		return ac

	case ExprFuncCall:
		fc := &nodes.FuncCall{
			Funcname: &nodes.List{Items: []nodes.Node{&nodes.String{Str: e.FuncName}}},
			AggStar:  e.FuncStar,
		}
		if len(e.FuncArgs) > 0 {
			items := make([]nodes.Node, len(e.FuncArgs))
			for i, arg := range e.FuncArgs {
				items[i] = reverseConvertExpr(arg)
			}
			fc.Args = &nodes.List{Items: items}
		}
		return fc

	case ExprOpExpr:
		ae := &nodes.A_Expr{
			Lexpr: reverseConvertExpr(e.Left),
			Rexpr: reverseConvertExpr(e.Right),
		}
		if e.OpName != "" {
			ae.Name = &nodes.List{Items: []nodes.Node{&nodes.String{Str: e.OpName}}}
		}
		return ae

	case ExprTypeCast:
		tn := makeTypeNameNode(e.CastType)
		return &nodes.TypeCast{
			Arg:      reverseConvertExpr(e.CastArg),
			TypeName: tn,
		}

	case ExprNullConst:
		return &nodes.A_Const{Isnull: true}

	case ExprCaseExpr:
		ce := &nodes.CaseExpr{}
		if len(e.CaseResults) > 0 {
			items := make([]nodes.Node, len(e.CaseResults))
			for i, r := range e.CaseResults {
				items[i] = &nodes.CaseWhen{Result: reverseConvertExpr(r)}
			}
			ce.Args = &nodes.List{Items: items}
		}
		if e.CaseElse != nil {
			ce.Defresult = reverseConvertExpr(e.CaseElse)
		}
		return ce

	case ExprCoalesceExpr:
		ce := &nodes.CoalesceExpr{}
		if len(e.CoalesceArgs) > 0 {
			items := make([]nodes.Node, len(e.CoalesceArgs))
			for i, arg := range e.CoalesceArgs {
				items[i] = reverseConvertExpr(arg)
			}
			ce.Args = &nodes.List{Items: items}
		}
		return ce

	case ExprSubquery:
		return &nodes.SubLink{
			Subselect: reverseConvertSelectStmt(e.Subquery),
		}

	case ExprBoolExpr:
		be := &nodes.BoolExpr{}
		switch e.BoolOp {
		case BoolAnd:
			be.Boolop = nodes.AND_EXPR
		case BoolOr:
			be.Boolop = nodes.OR_EXPR
		case BoolNot:
			be.Boolop = nodes.NOT_EXPR
		}
		if len(e.BoolArgs) > 0 {
			items := make([]nodes.Node, len(e.BoolArgs))
			for i, arg := range e.BoolArgs {
				items[i] = reverseConvertExpr(arg)
			}
			be.Args = &nodes.List{Items: items}
		}
		return be

	default:
		return &nodes.A_Const{Isnull: true}
	}
}

// reverseConvertFromItem converts a catalog FromItem back to a nodes.Node.
func reverseConvertFromItem(fi *FromItem) nodes.Node {
	if fi == nil {
		return nil
	}
	switch fi.Kind {
	case FromTable:
		rv := &nodes.RangeVar{
			Schemaname: fi.Schema,
			Relname:    fi.Table,
		}
		if fi.Alias != "" {
			rv.Alias = &nodes.Alias{Aliasname: fi.Alias}
		}
		return rv

	case FromJoin:
		je := &nodes.JoinExpr{
			Larg: reverseConvertFromItem(fi.JoinLeft),
			Rarg: reverseConvertFromItem(fi.JoinRight),
		}
		switch fi.JoinType {
		case JoinInner:
			je.Jointype = nodes.JOIN_INNER
		case JoinLeft:
			je.Jointype = nodes.JOIN_LEFT
		case JoinRight:
			je.Jointype = nodes.JOIN_RIGHT
		case JoinFull:
			je.Jointype = nodes.JOIN_FULL
		case JoinCross:
			je.Jointype = nodes.JOIN_INNER
			je.IsNatural = true
		}
		return je

	case FromSubquery:
		rs := &nodes.RangeSubselect{
			Subquery: reverseConvertSelectStmt(fi.Subquery),
		}
		if fi.SubAlias != "" {
			rs.Alias = &nodes.Alias{Aliasname: fi.SubAlias}
		}
		return rs

	default:
		return nil
	}
}

// makeViewStmt builds a nodes.ViewStmt from test parameters.
func makeViewStmt(schema, name string, query *SelectStmt, columnNames []string) *nodes.ViewStmt {
	stmt := &nodes.ViewStmt{
		View: &nodes.RangeVar{Schemaname: schema, Relname: name},
	}
	if query != nil {
		stmt.Query = reverseConvertSelectStmt(query)
	}
	if len(columnNames) > 0 {
		stmt.Aliases = makeStringListNode(columnNames)
	}
	return stmt
}

// makeDropViewStmt builds a nodes.DropStmt for DROP VIEW from test parameters.
func makeDropViewStmt(schema, name string, ifExists, cascade bool) *nodes.DropStmt {
	behavior := int(nodes.DROP_RESTRICT)
	if cascade {
		behavior = int(nodes.DROP_CASCADE)
	}
	var items []nodes.Node
	if schema != "" {
		items = []nodes.Node{&nodes.String{Str: schema}, &nodes.String{Str: name}}
	} else {
		items = []nodes.Node{&nodes.String{Str: name}}
	}
	return &nodes.DropStmt{
		Objects:    &nodes.List{Items: []nodes.Node{&nodes.List{Items: items}}},
		RemoveType: int(nodes.OBJECT_VIEW),
		Behavior:   behavior,
		Missing_ok: ifExists,
	}
}

// makeViewStmtReplace builds a nodes.ViewStmt with Replace=true.
func makeViewStmtReplace(schema, name string, query *SelectStmt, columnNames []string) *nodes.ViewStmt {
	stmt := makeViewStmt(schema, name, query, columnNames)
	stmt.Replace = true
	return stmt
}

// --- View tests ---

func TestCreateViewSimple(t *testing.T) {
	c := newTestCatalogWithTable()

	err := c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			TargetList: []*ResTarget{
				{Val: &Expr{Kind: ExprColumnRef, ColumnName: "id"}},
				{Val: &Expr{Kind: ExprColumnRef, ColumnName: "name"}},
			},
			From: []*FromItem{{Kind: FromTable, Table: "t"}},
		}, nil))
	if err != nil {
		t.Fatal(err)
	}

	r := c.GetRelation("", "v")
	if r == nil {
		t.Fatal("view not found")
	}
	if r.RelKind != 'v' {
		t.Errorf("RelKind: got %c, want 'v'", r.RelKind)
	}
	if len(r.Columns) != 2 {
		t.Fatalf("columns: got %d, want 2", len(r.Columns))
	}
	if r.Columns[0].TypeOID != INT4OID {
		t.Errorf("col 0 type: got %d, want %d", r.Columns[0].TypeOID, INT4OID)
	}
	if r.Columns[1].TypeOID != TEXTOID {
		t.Errorf("col 1 type: got %d, want %d", r.Columns[1].TypeOID, TEXTOID)
	}
	if r.Columns[0].Name != "id" {
		t.Errorf("col 0 name: got %q, want %q", r.Columns[0].Name, "id")
	}
}

func TestCreateViewWithStar(t *testing.T) {
	c := newTestCatalogWithTable()

	err := c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprStar}}},
			From:       []*FromItem{{Kind: FromTable, Table: "t"}},
		}, nil))
	if err != nil {
		t.Fatal(err)
	}

	r := c.GetRelation("", "v")
	if len(r.Columns) != 3 {
		t.Fatalf("columns: got %d, want 3", len(r.Columns))
	}
}

func TestCreateViewWithExplicitColumnNames(t *testing.T) {
	c := newTestCatalogWithTable()

	err := c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			TargetList: []*ResTarget{
				{Val: &Expr{Kind: ExprColumnRef, ColumnName: "id"}},
				{Val: &Expr{Kind: ExprColumnRef, ColumnName: "name"}},
			},
			From: []*FromItem{{Kind: FromTable, Table: "t"}},
		}, []string{"a", "b"}))
	if err != nil {
		t.Fatal(err)
	}

	r := c.GetRelation("", "v")
	if r.Columns[0].Name != "a" {
		t.Errorf("col 0 name: got %q, want %q", r.Columns[0].Name, "a")
	}
	if r.Columns[1].Name != "b" {
		t.Errorf("col 1 name: got %q, want %q", r.Columns[1].Name, "b")
	}
}

func TestCreateViewColumnCountMismatch(t *testing.T) {
	c := newTestCatalogWithTable()

	err := c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			TargetList: []*ResTarget{
				{Val: &Expr{Kind: ExprColumnRef, ColumnName: "id"}},
				{Val: &Expr{Kind: ExprColumnRef, ColumnName: "name"}},
			},
			From: []*FromItem{{Kind: FromTable, Table: "t"}},
		}, []string{"a"})) // only 1, but query produces 2
	assertErrorCode(t, err, CodeDatatypeMismatch)
}

func TestCreateViewWithExpression(t *testing.T) {
	c := newTestCatalogWithTable()

	err := c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			TargetList: []*ResTarget{{
				Name: "total",
				Val: &Expr{
					Kind:     ExprFuncCall,
					FuncName: "count",
					FuncStar: true,
				},
			}},
			From: []*FromItem{{Kind: FromTable, Table: "t"}},
		}, nil))
	if err != nil {
		t.Fatal(err)
	}

	r := c.GetRelation("", "v")
	if r.Columns[0].TypeOID != INT8OID {
		t.Errorf("col type: got %d, want %d", r.Columns[0].TypeOID, INT8OID)
	}
	if r.Columns[0].Name != "total" {
		t.Errorf("col name: got %q, want %q", r.Columns[0].Name, "total")
	}
}

func TestCreateViewWithJoin(t *testing.T) {
	c := New()
	c.DefineRelation(makeCreateTableStmt("", "t1", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "a", Type: TypeName{Name: "text", TypeMod: -1}},
	}, nil, false), 'r')
	c.DefineRelation(makeCreateTableStmt("", "t2", []ColumnDef{
		{Name: "id", Type: TypeName{Name: "int4", TypeMod: -1}},
		{Name: "b", Type: TypeName{Name: "int8", TypeMod: -1}},
	}, nil, false), 'r')

	err := c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			TargetList: []*ResTarget{
				{Val: &Expr{Kind: ExprColumnRef, TableName: "t1", ColumnName: "a"}},
				{Val: &Expr{Kind: ExprColumnRef, TableName: "t2", ColumnName: "b"}},
			},
			From: []*FromItem{{
				Kind:      FromJoin,
				JoinType:  JoinInner,
				JoinLeft:  &FromItem{Kind: FromTable, Table: "t1"},
				JoinRight: &FromItem{Kind: FromTable, Table: "t2"},
			}},
		}, nil))
	if err != nil {
		t.Fatal(err)
	}

	r := c.GetRelation("", "v")
	if len(r.Columns) != 2 {
		t.Fatalf("columns: got %d, want 2", len(r.Columns))
	}
	if r.Columns[0].TypeOID != TEXTOID {
		t.Errorf("col 0 type: got %d, want %d", r.Columns[0].TypeOID, TEXTOID)
	}
	if r.Columns[1].TypeOID != INT8OID {
		t.Errorf("col 1 type: got %d, want %d", r.Columns[1].TypeOID, INT8OID)
	}
}

func TestCreateViewWithUnion(t *testing.T) {
	c := newTestCatalogWithTable()

	err := c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			Op: SetOpUnion,
			Left: &SelectStmt{
				TargetList: []*ResTarget{{Val: &Expr{Kind: ExprColumnRef, ColumnName: "id"}}},
				From:       []*FromItem{{Kind: FromTable, Table: "t"}},
			},
			Right: &SelectStmt{
				TargetList: []*ResTarget{{Val: &Expr{Kind: ExprColumnRef, ColumnName: "val"}}},
				From:       []*FromItem{{Kind: FromTable, Table: "t"}},
			},
		}, nil))
	if err != nil {
		t.Fatal(err)
	}

	r := c.GetRelation("", "v")
	if len(r.Columns) != 1 {
		t.Fatalf("columns: got %d, want 1", len(r.Columns))
	}
	// int4 UNION int8 -> int8
	if r.Columns[0].TypeOID != INT8OID {
		t.Errorf("col type: got %d, want %d", r.Columns[0].TypeOID, INT8OID)
	}
}

func TestCreateViewDuplicate(t *testing.T) {
	c := newTestCatalogWithTable()

	query := &SelectStmt{
		TargetList: []*ResTarget{{Val: &Expr{Kind: ExprColumnRef, ColumnName: "id"}}},
		From:       []*FromItem{{Kind: FromTable, Table: "t"}},
	}
	c.DefineView(makeViewStmt("", "v", query, nil))

	err := c.DefineView(makeViewStmt("", "v", query, nil))
	assertErrorCode(t, err, CodeDuplicateTable)
}

func TestDropView(t *testing.T) {
	c := newTestCatalogWithTable()

	c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprStar}}},
			From:       []*FromItem{{Kind: FromTable, Table: "t"}},
		}, nil))

	err := c.RemoveRelations(makeDropViewStmt("", "v", false, false))
	if err != nil {
		t.Fatal(err)
	}
	if r := c.GetRelation("", "v"); r != nil {
		t.Fatal("view still exists after DROP")
	}
}

func TestDropViewIfExists(t *testing.T) {
	c := New()
	err := c.RemoveRelations(makeDropViewStmt("", "nosuch", true, false))
	if err != nil {
		t.Fatalf("IF EXISTS should not error, got: %v", err)
	}
}

func TestDropViewNonexistent(t *testing.T) {
	c := New()
	err := c.RemoveRelations(makeDropViewStmt("", "nosuch", false, false))
	assertErrorCode(t, err, CodeUndefinedTable)
}

func TestDropTableOnView(t *testing.T) {
	c := newTestCatalogWithTable()

	c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprStar}}},
			From:       []*FromItem{{Kind: FromTable, Table: "t"}},
		}, nil))

	err := c.RemoveRelations(makeDropTableStmt("", "v", false, false))
	assertErrorCode(t, err, CodeWrongObjectType)
}

func TestDropViewOnTable(t *testing.T) {
	c := newTestCatalogWithTable()

	err := c.RemoveRelations(makeDropViewStmt("", "t", false, false))
	assertErrorCode(t, err, CodeWrongObjectType)
}

func TestDropTableCascadeDropsViews(t *testing.T) {
	c := newTestCatalogWithTable()

	c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprStar}}},
			From:       []*FromItem{{Kind: FromTable, Table: "t"}},
		}, nil))

	// DROP TABLE t without CASCADE should fail (view depends on it).
	err := c.RemoveRelations(makeDropTableStmt("", "t", false, false))
	assertErrorCode(t, err, CodeDependentObjects)

	// DROP TABLE t CASCADE should succeed and also drop the view.
	err = c.RemoveRelations(makeDropTableStmt("", "t", false, true))
	if err != nil {
		t.Fatal(err)
	}
	if r := c.GetRelation("", "t"); r != nil {
		t.Error("table still exists after CASCADE DROP")
	}
	if r := c.GetRelation("", "v"); r != nil {
		t.Error("view still exists after CASCADE DROP of source table")
	}
}

func TestViewDependencySourceTableRequired(t *testing.T) {
	c := newTestCatalogWithTable()

	c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprColumnRef, ColumnName: "id"}}},
			From:       []*FromItem{{Kind: FromTable, Table: "t"}},
		}, nil))

	// Cannot drop source table without CASCADE.
	err := c.RemoveRelations(makeDropTableStmt("", "t", false, false))
	assertErrorCode(t, err, CodeDependentObjects)
}

func TestViewChaining(t *testing.T) {
	c := newTestCatalogWithTable()

	// v1 depends on t
	c.DefineView(makeViewStmt("", "v1",
		&SelectStmt{
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprStar}}},
			From:       []*FromItem{{Kind: FromTable, Table: "t"}},
		}, nil))

	// v2 depends on v1
	c.DefineView(makeViewStmt("", "v2",
		&SelectStmt{
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprStar}}},
			From:       []*FromItem{{Kind: FromTable, Table: "v1"}},
		}, nil))

	// Cannot drop v1 without CASCADE (v2 depends on it).
	err := c.RemoveRelations(makeDropViewStmt("", "v1", false, false))
	assertErrorCode(t, err, CodeDependentObjects)

	// DROP VIEW v1 CASCADE drops both v1 and v2.
	err = c.RemoveRelations(makeDropViewStmt("", "v1", false, true))
	if err != nil {
		t.Fatal(err)
	}
	if r := c.GetRelation("", "v1"); r != nil {
		t.Error("v1 still exists")
	}
	if r := c.GetRelation("", "v2"); r != nil {
		t.Error("v2 still exists after CASCADE")
	}
	// t should still exist.
	if r := c.GetRelation("", "t"); r == nil {
		t.Error("source table should still exist")
	}
}

func TestDropViewCascadeChain(t *testing.T) {
	c := newTestCatalogWithTable()

	// t -> v1 -> v2 -> v3
	c.DefineView(makeViewStmt("", "v1",
		&SelectStmt{
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprStar}}},
			From:       []*FromItem{{Kind: FromTable, Table: "t"}},
		}, nil))
	c.DefineView(makeViewStmt("", "v2",
		&SelectStmt{
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprStar}}},
			From:       []*FromItem{{Kind: FromTable, Table: "v1"}},
		}, nil))
	c.DefineView(makeViewStmt("", "v3",
		&SelectStmt{
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprStar}}},
			From:       []*FromItem{{Kind: FromTable, Table: "v2"}},
		}, nil))

	// DROP TABLE t CASCADE drops everything.
	err := c.RemoveRelations(makeDropTableStmt("", "t", false, true))
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"t", "v1", "v2", "v3"} {
		if r := c.GetRelation("", name); r != nil {
			t.Errorf("%q still exists after CASCADE", name)
		}
	}
}

func TestViewRowType(t *testing.T) {
	c := newTestCatalogWithTable()

	c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			TargetList: []*ResTarget{{Val: &Expr{Kind: ExprColumnRef, ColumnName: "id"}}},
			From:       []*FromItem{{Kind: FromTable, Table: "t"}},
		}, nil))

	// View row type should be resolvable.
	oid, _, err := c.ResolveType(TypeName{Name: "v", TypeMod: -1})
	if err != nil {
		t.Fatalf("view row type not resolvable: %v", err)
	}
	typ := c.TypeByOID(oid)
	if typ == nil || typ.Type != 'c' {
		t.Fatal("expected composite type")
	}

	// Drop the view -- type should be gone.
	c.RemoveRelations(makeDropViewStmt("", "v", false, false))
	_, _, err = c.ResolveType(TypeName{Name: "v", TypeMod: -1})
	if err == nil {
		t.Fatal("view row type should not be resolvable after DROP")
	}
}

// --- CREATE OR REPLACE VIEW tests ---

func TestCreateOrReplaceViewSameColumns(t *testing.T) {
	c := newTestCatalogWithTable()

	query := &SelectStmt{
		TargetList: []*ResTarget{
			{Val: &Expr{Kind: ExprColumnRef, ColumnName: "id"}},
			{Val: &Expr{Kind: ExprColumnRef, ColumnName: "name"}},
		},
		From: []*FromItem{{Kind: FromTable, Table: "t"}},
	}
	if err := c.DefineView(makeViewStmt("", "v", query, nil)); err != nil {
		t.Fatal(err)
	}

	// Replace with same columns should succeed.
	if err := c.DefineView(makeViewStmtReplace("", "v", query, nil)); err != nil {
		t.Fatal(err)
	}

	r := c.GetRelation("", "v")
	if len(r.Columns) != 2 {
		t.Fatalf("columns: got %d, want 2", len(r.Columns))
	}
}

func TestCreateOrReplaceViewAddColumn(t *testing.T) {
	c := newTestCatalogWithTable()

	// Create view with 1 column.
	c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			TargetList: []*ResTarget{
				{Val: &Expr{Kind: ExprColumnRef, ColumnName: "id"}},
			},
			From: []*FromItem{{Kind: FromTable, Table: "t"}},
		}, nil))

	// Replace with 2 columns (superset) should succeed.
	err := c.DefineView(makeViewStmtReplace("", "v",
		&SelectStmt{
			TargetList: []*ResTarget{
				{Val: &Expr{Kind: ExprColumnRef, ColumnName: "id"}},
				{Val: &Expr{Kind: ExprColumnRef, ColumnName: "name"}},
			},
			From: []*FromItem{{Kind: FromTable, Table: "t"}},
		}, nil))
	if err != nil {
		t.Fatal(err)
	}

	r := c.GetRelation("", "v")
	if len(r.Columns) != 2 {
		t.Fatalf("columns: got %d, want 2", len(r.Columns))
	}
	if r.Columns[1].Name != "name" {
		t.Errorf("col 1 name: got %q, want %q", r.Columns[1].Name, "name")
	}
}

func TestCreateOrReplaceViewDropColumnFails(t *testing.T) {
	c := newTestCatalogWithTable()

	// Create view with 2 columns.
	c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			TargetList: []*ResTarget{
				{Val: &Expr{Kind: ExprColumnRef, ColumnName: "id"}},
				{Val: &Expr{Kind: ExprColumnRef, ColumnName: "name"}},
			},
			From: []*FromItem{{Kind: FromTable, Table: "t"}},
		}, nil))

	// Replace with 1 column (fewer) should fail.
	err := c.DefineView(makeViewStmtReplace("", "v",
		&SelectStmt{
			TargetList: []*ResTarget{
				{Val: &Expr{Kind: ExprColumnRef, ColumnName: "id"}},
			},
			From: []*FromItem{{Kind: FromTable, Table: "t"}},
		}, nil))
	assertErrorCode(t, err, CodeInvalidObjectDefinition)
}

func TestCreateOrReplaceViewChangeTypeFails(t *testing.T) {
	c := newTestCatalogWithTable()

	// Create view with int4 column.
	c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			TargetList: []*ResTarget{
				{Val: &Expr{Kind: ExprColumnRef, ColumnName: "id"}},
			},
			From: []*FromItem{{Kind: FromTable, Table: "t"}},
		}, nil))

	// Replace with int8 column (type mismatch) should fail.
	err := c.DefineView(makeViewStmtReplace("", "v",
		&SelectStmt{
			TargetList: []*ResTarget{
				{Val: &Expr{Kind: ExprColumnRef, ColumnName: "val"}},
			},
			From: []*FromItem{{Kind: FromTable, Table: "t"}},
		}, []string{"id"}))
	assertErrorCode(t, err, CodeInvalidObjectDefinition)
}

func TestCreateOrReplaceViewChangeNameFails(t *testing.T) {
	c := newTestCatalogWithTable()

	c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			TargetList: []*ResTarget{
				{Val: &Expr{Kind: ExprColumnRef, ColumnName: "id"}},
			},
			From: []*FromItem{{Kind: FromTable, Table: "t"}},
		}, nil))

	// Replace with different column name should fail.
	err := c.DefineView(makeViewStmtReplace("", "v",
		&SelectStmt{
			TargetList: []*ResTarget{
				{Val: &Expr{Kind: ExprColumnRef, ColumnName: "id"}},
			},
			From: []*FromItem{{Kind: FromTable, Table: "t"}},
		}, []string{"other_id"}))
	assertErrorCode(t, err, CodeInvalidObjectDefinition)
}

func TestCreateOrReplaceOnTableFails(t *testing.T) {
	c := newTestCatalogWithTable()

	// Try to OR REPLACE a table name as a view.
	err := c.DefineView(makeViewStmtReplace("", "t",
		&SelectStmt{
			TargetList: []*ResTarget{
				{Val: &Expr{Kind: ExprColumnRef, ColumnName: "id"}},
			},
			From: []*FromItem{{Kind: FromTable, Table: "t"}},
		}, nil))
	assertErrorCode(t, err, CodeWrongObjectType)
}

// --- Collation propagation tests ---

func TestViewColumnCollationFromTextColumn(t *testing.T) {
	c := newTestCatalogWithTable()

	err := c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			TargetList: []*ResTarget{
				{Val: &Expr{Kind: ExprColumnRef, ColumnName: "name"}},
			},
			From: []*FromItem{{Kind: FromTable, Table: "t"}},
		}, nil))
	if err != nil {
		t.Fatal(err)
	}

	r := c.GetRelation("", "v")
	// text has type default collation 100 (DEFAULT_COLLATION_OID).
	if r.Columns[0].Collation != 100 {
		t.Errorf("Collation: got %d, want 100", r.Columns[0].Collation)
	}
}

func TestViewColumnCollationFromIntColumn(t *testing.T) {
	c := newTestCatalogWithTable()

	err := c.DefineView(makeViewStmt("", "v",
		&SelectStmt{
			TargetList: []*ResTarget{
				{Val: &Expr{Kind: ExprColumnRef, ColumnName: "id"}},
			},
			From: []*FromItem{{Kind: FromTable, Table: "t"}},
		}, nil))
	if err != nil {
		t.Fatal(err)
	}

	r := c.GetRelation("", "v")
	// int4 is not collatable, collation should be 0.
	if r.Columns[0].Collation != 0 {
		t.Errorf("Collation: got %d, want 0", r.Columns[0].Collation)
	}
}

func TestCreateOrReplaceViewSameCollationOK(t *testing.T) {
	c := newTestCatalogWithTable()

	query := &SelectStmt{
		TargetList: []*ResTarget{
			{Val: &Expr{Kind: ExprColumnRef, ColumnName: "name"}},
		},
		From: []*FromItem{{Kind: FromTable, Table: "t"}},
	}

	// Create view.
	if err := c.DefineView(makeViewStmt("", "v", query, nil)); err != nil {
		t.Fatal(err)
	}

	// OR REPLACE with same column (same collation) should succeed.
	if err := c.DefineView(makeViewStmtReplace("", "v", query, nil)); err != nil {
		t.Fatal(err)
	}
}

func TestCreateOrReplaceViewChangeCollationFails(t *testing.T) {
	c := newTestCatalogWithTable()

	query := &SelectStmt{
		TargetList: []*ResTarget{
			{Val: &Expr{Kind: ExprColumnRef, ColumnName: "name"}},
		},
		From: []*FromItem{{Kind: FromTable, Table: "t"}},
	}

	// Create view with text column (collation=100).
	if err := c.DefineView(makeViewStmt("", "v", query, nil)); err != nil {
		t.Fatal(err)
	}

	// Manually change the view column's collation to simulate a different collation.
	r := c.GetRelation("", "v")
	r.Columns[0].Collation = 950 // C collation

	// OR REPLACE should fail because new column has collation 100 but old has 950.
	err := c.DefineView(makeViewStmtReplace("", "v", query, nil))
	assertErrorCode(t, err, CodeInvalidObjectDefinition)
}
