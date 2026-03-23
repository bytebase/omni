package parsertest

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
)

func TestLocExplainStmt(t *testing.T) {
	sql := "EXPLAIN SELECT 1"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.ExplainStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("ExplainStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("ExplainStmt text = %q, want %q", got, sql)
	}
}

func TestLocCallStmt(t *testing.T) {
	sql := "CALL myproc()"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CallStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("CallStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("CallStmt text = %q, want %q", got, sql)
	}
}

func TestLocDoStmt(t *testing.T) {
	sql := "DO $$ BEGIN RAISE NOTICE 'hi'; END $$"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.DoStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("DoStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("DoStmt text = %q, want %q", got, sql)
	}
}

func TestLocCheckPointStmt(t *testing.T) {
	sql := "CHECKPOINT"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.CheckPointStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("CheckPointStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("CheckPointStmt text = %q, want %q", got, sql)
	}
}

func TestLocDiscardStmt(t *testing.T) {
	sql := "DISCARD ALL"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.DiscardStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("DiscardStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("DiscardStmt text = %q, want %q", got, sql)
	}
}

func TestLocListenStmt(t *testing.T) {
	sql := "LISTEN mychannel"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.ListenStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("ListenStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("ListenStmt text = %q, want %q", got, sql)
	}
}

func TestLocNotifyStmt(t *testing.T) {
	sql := "NOTIFY mychannel, 'payload'"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.NotifyStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("NotifyStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("NotifyStmt text = %q, want %q", got, sql)
	}
}

func TestLocUnlistenStmt(t *testing.T) {
	sql := "UNLISTEN mychannel"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.UnlistenStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("UnlistenStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("UnlistenStmt text = %q, want %q", got, sql)
	}
}

func TestLocLoadStmt(t *testing.T) {
	sql := "LOAD 'mylib'"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.LoadStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("LoadStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("LoadStmt text = %q, want %q", got, sql)
	}
}

func TestLocReassignOwnedStmt(t *testing.T) {
	sql := "REASSIGN OWNED BY oldrole TO newrole"
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	raw := list.Items[0].(*nodes.RawStmt)
	stmt := raw.Stmt.(*nodes.ReassignOwnedStmt)
	if stmt.Loc.Start == -1 || stmt.Loc.End == -1 {
		t.Fatalf("ReassignOwnedStmt Loc not set: %+v", stmt.Loc)
	}
	got := sql[stmt.Loc.Start:stmt.Loc.End]
	if got != sql {
		t.Errorf("ReassignOwnedStmt text = %q, want %q", got, sql)
	}
}
