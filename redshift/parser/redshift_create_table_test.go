package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
)

func TestRedshiftCreateTableOptionsParse(t *testing.T) {
	tests := []string{
		"CREATE TABLE users (id INT) DISTSTYLE EVEN;",
		"CREATE TABLE users (id INT) DISTSTYLE KEY DISTKEY(id);",
		"CREATE TABLE users (id INT) SORTKEY(id);",
		"CREATE TABLE users (id INT) SORTKEY AUTO;",
		"CREATE TABLE users (id INT) INTERLEAVED SORTKEY(id);",
		"CREATE TABLE users (id INT) ENCODE AUTO;",
		"CREATE TABLE users (id INT IDENTITY(1,1));",
		"CREATE TABLE users (id INT ENCODE lzo);",
		"CREATE TABLE users (id INT DISTKEY);",
		"CREATE TABLE users (id INT SORTKEY);",
	}

	for _, sql := range tests {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Fatalf("Parse returned error: %v", err)
			}
		})
	}
}

func TestRedshiftCreateTableOptionsAST(t *testing.T) {
	tree, err := Parse("CREATE TABLE users (id INT ENCODE lzo IDENTITY(1,1)) DISTSTYLE KEY DISTKEY(id) INTERLEAVED SORTKEY(id) BACKUP NO;")
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	stmt := firstCreateStmt(t, tree)
	assertDefElemString(t, stmt.Options, "diststyle", "key")
	assertDefElemList(t, stmt.Options, "distkey")
	assertDefElemString(t, stmt.Options, "sortstyle", "interleaved")
	assertDefElemList(t, stmt.Options, "sortkey")
	assertDefElemString(t, stmt.Options, "backup", "no")

	if stmt.TableElts == nil || len(stmt.TableElts.Items) != 1 {
		t.Fatalf("expected one table element, got %#v", stmt.TableElts)
	}
	col, ok := stmt.TableElts.Items[0].(*nodes.ColumnDef)
	if !ok {
		t.Fatalf("expected ColumnDef, got %T", stmt.TableElts.Items[0])
	}
	if col.Compression != "lzo" {
		t.Fatalf("expected column compression lzo, got %q", col.Compression)
	}
	if col.Constraints == nil || len(col.Constraints.Items) != 1 {
		t.Fatalf("expected one column constraint, got %#v", col.Constraints)
	}
	constraint, ok := col.Constraints.Items[0].(*nodes.Constraint)
	if !ok {
		t.Fatalf("expected identity Constraint, got %T", col.Constraints.Items[0])
	}
	if constraint.Contype != nodes.CONSTR_IDENTITY {
		t.Fatalf("expected identity constraint, got %v", constraint.Contype)
	}
	if constraint.Options == nil || len(constraint.Options.Items) != 2 {
		t.Fatalf("expected identity start and increment options, got %#v", constraint.Options)
	}
}

func TestRedshiftCreateTableAsOptionsAST(t *testing.T) {
	tree, err := Parse("CREATE TABLE users DISTSTYLE EVEN SORTKEY(id) AS SELECT 1 AS id;")
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(tree.Items) != 1 {
		t.Fatalf("expected one statement, got %d", len(tree.Items))
	}
	raw, ok := tree.Items[0].(*nodes.RawStmt)
	if !ok {
		t.Fatalf("expected RawStmt, got %T", tree.Items[0])
	}
	stmt, ok := raw.Stmt.(*nodes.CreateTableAsStmt)
	if !ok {
		t.Fatalf("expected CreateTableAsStmt, got %T", raw.Stmt)
	}
	assertDefElemString(t, stmt.Into.Options, "diststyle", "even")
	assertDefElemList(t, stmt.Into.Options, "sortkey")
}

func TestRedshiftCreateTableAsHashTempTable(t *testing.T) {
	tree, err := Parse("CREATE TABLE #session_data AS SELECT 1 AS session_id;")
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(tree.Items) != 1 {
		t.Fatalf("expected one statement, got %d", len(tree.Items))
	}
	raw, ok := tree.Items[0].(*nodes.RawStmt)
	if !ok {
		t.Fatalf("expected RawStmt, got %T", tree.Items[0])
	}
	stmt, ok := raw.Stmt.(*nodes.CreateTableAsStmt)
	if !ok {
		t.Fatalf("expected CreateTableAsStmt, got %T", raw.Stmt)
	}
	if stmt.Into == nil || stmt.Into.Rel == nil {
		t.Fatalf("expected target relation")
	}
	if stmt.Into.Rel.Relname != "session_data" {
		t.Fatalf("Relname = %q, want session_data", stmt.Into.Rel.Relname)
	}
	if stmt.Into.Rel.Relpersistence != nodes.RELPERSISTENCE_TEMP {
		t.Fatalf("Relpersistence = %v, want temporary", stmt.Into.Rel.Relpersistence)
	}
}

func TestRedshiftSortkeyAutoAST(t *testing.T) {
	tree, err := Parse("CREATE TABLE users (id INT) SORTKEY AUTO;")
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	stmt := firstCreateStmt(t, tree)
	assertDefElemString(t, stmt.Options, "sortkey", "auto")
}

func firstCreateStmt(t *testing.T, tree *nodes.List) *nodes.CreateStmt {
	t.Helper()
	if len(tree.Items) != 1 {
		t.Fatalf("expected one statement, got %d", len(tree.Items))
	}
	raw, ok := tree.Items[0].(*nodes.RawStmt)
	if !ok {
		t.Fatalf("expected RawStmt, got %T", tree.Items[0])
	}
	stmt, ok := raw.Stmt.(*nodes.CreateStmt)
	if !ok {
		t.Fatalf("expected CreateStmt, got %T", raw.Stmt)
	}
	return stmt
}

func findDefElem(list *nodes.List, name string) *nodes.DefElem {
	if list == nil {
		return nil
	}
	for _, item := range list.Items {
		elem, ok := item.(*nodes.DefElem)
		if ok && elem.Defname == name {
			return elem
		}
	}
	return nil
}

func assertDefElemString(t *testing.T, list *nodes.List, name string, want string) {
	t.Helper()
	elem := findDefElem(list, name)
	if elem == nil {
		t.Fatalf("expected option %q in %#v", name, list)
	}
	arg, ok := elem.Arg.(*nodes.String)
	if !ok {
		t.Fatalf("expected option %q string arg, got %T", name, elem.Arg)
	}
	if arg.Str != want {
		t.Fatalf("expected option %q=%q, got %q", name, want, arg.Str)
	}
}

func assertDefElemList(t *testing.T, list *nodes.List, name string) {
	t.Helper()
	elem := findDefElem(list, name)
	if elem == nil {
		t.Fatalf("expected option %q in %#v", name, list)
	}
	if _, ok := elem.Arg.(*nodes.List); !ok {
		t.Fatalf("expected option %q list arg, got %T", name, elem.Arg)
	}
}
