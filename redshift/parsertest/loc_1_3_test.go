package parsertest

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
	"github.com/bytebase/omni/redshift/parser"
)

// ---------------------------------------------------------------------------
// Section 1.3: JSON nodes
// ---------------------------------------------------------------------------

func TestLocJsonKeyValue(t *testing.T) {
	sql := "SELECT JSON_OBJECT('key' VALUE 'val')"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)
	rt := sel.TargetList.Items[0].(*nodes.ResTarget)
	obj := rt.Val.(*nodes.JsonObjectConstructor)
	if obj.Exprs == nil || len(obj.Exprs.Items) == 0 {
		t.Fatal("expected key-value pairs")
	}
	kv := obj.Exprs.Items[0].(*nodes.JsonKeyValue)
	got := sql[kv.Loc.Start:kv.Loc.End]
	want := "'key' VALUE 'val'"
	if got != want {
		t.Errorf("JsonKeyValue text = %q, want %q", got, want)
	}
}

func TestLocJsonKeyValueColon(t *testing.T) {
	sql := "SELECT JSON_OBJECT('key': 'val')"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)
	rt := sel.TargetList.Items[0].(*nodes.ResTarget)
	obj := rt.Val.(*nodes.JsonObjectConstructor)
	if obj.Exprs == nil || len(obj.Exprs.Items) == 0 {
		t.Fatal("expected key-value pairs")
	}
	kv := obj.Exprs.Items[0].(*nodes.JsonKeyValue)
	got := sql[kv.Loc.Start:kv.Loc.End]
	want := "'key': 'val'"
	if got != want {
		t.Errorf("JsonKeyValue text = %q, want %q", got, want)
	}
}

func TestLocJsonObjectAgg(t *testing.T) {
	sql := "SELECT JSON_OBJECTAGG(k VALUE v) FROM t"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)
	rt := sel.TargetList.Items[0].(*nodes.ResTarget)
	agg := rt.Val.(*nodes.JsonObjectAgg)
	got := sql[agg.Loc.Start:agg.Loc.End]
	want := "JSON_OBJECTAGG(k VALUE v)"
	if got != want {
		t.Errorf("JsonObjectAgg text = %q, want %q", got, want)
	}
}

func TestLocJsonArrayAgg(t *testing.T) {
	sql := "SELECT JSON_ARRAYAGG(v) FROM t"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)
	rt := sel.TargetList.Items[0].(*nodes.ResTarget)
	agg := rt.Val.(*nodes.JsonArrayAgg)
	got := sql[agg.Loc.Start:agg.Loc.End]
	want := "JSON_ARRAYAGG(v)"
	if got != want {
		t.Errorf("JsonArrayAgg text = %q, want %q", got, want)
	}
}

func TestLocJsonOutput(t *testing.T) {
	sql := "SELECT JSON_OBJECT('k': 'v' RETURNING text)"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)
	rt := sel.TargetList.Items[0].(*nodes.ResTarget)
	obj := rt.Val.(*nodes.JsonObjectConstructor)
	if obj.Output == nil {
		t.Fatal("expected JsonOutput to be set")
	}
	got := sql[obj.Output.Loc.Start:obj.Output.Loc.End]
	want := "RETURNING text"
	if got != want {
		t.Errorf("JsonOutput text = %q, want %q", got, want)
	}
}

func TestLocJsonValueExpr(t *testing.T) {
	// In JSON_OBJECT('key' VALUE expr), the value expression is a JsonValueExpr.
	sql := "SELECT JSON_OBJECT('key' VALUE 'val')"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)
	rt := sel.TargetList.Items[0].(*nodes.ResTarget)
	obj := rt.Val.(*nodes.JsonObjectConstructor)
	if obj.Exprs == nil || len(obj.Exprs.Items) == 0 {
		t.Fatal("expected key-value pairs")
	}
	kv := obj.Exprs.Items[0].(*nodes.JsonKeyValue)
	if kv.Value == nil {
		t.Fatal("expected JsonValueExpr to be set")
	}
	got := sql[kv.Value.Loc.Start:kv.Value.Loc.End]
	want := "'val'"
	if got != want {
		t.Errorf("JsonValueExpr text = %q, want %q", got, want)
	}
}

func TestLocJsonArgument(t *testing.T) {
	// JsonArgument appears in PASSING clause of JSON_VALUE, JSON_QUERY, etc.
	sql := "SELECT JSON_VALUE('{\"x\":1}', '$.x' PASSING 42 AS val)"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	sel := raw.Stmt.(*nodes.SelectStmt)
	rt := sel.TargetList.Items[0].(*nodes.ResTarget)
	fn := rt.Val.(*nodes.JsonFuncExpr)
	if fn.Passing == nil || len(fn.Passing.Items) == 0 {
		t.Fatal("expected PASSING arguments")
	}
	arg := fn.Passing.Items[0].(*nodes.JsonArgument)
	got := sql[arg.Loc.Start:arg.Loc.End]
	want := "42 AS val"
	if got != want {
		t.Errorf("JsonArgument text = %q, want %q", got, want)
	}
}
