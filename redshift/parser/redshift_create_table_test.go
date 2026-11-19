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

// TestRedshiftCreateTableGeneratedIdentitySeedStepParse pins the seed and step
// Redshift writes in IDENTITY's parentheses, bare and after GENERATED BY
// DEFAULT AS IDENTITY, where PostgreSQL writes a SeqOptList.
func TestRedshiftCreateTableGeneratedIdentitySeedStepParse(t *testing.T) {
	for _, tc := range []struct {
		sql       string
		when      byte
		options   int
		wantStart int64
	}{
		{"CREATE TABLE a (id INT GENERATED BY DEFAULT AS IDENTITY (1, 1));", 'd', 2, 1},
		{"CREATE TABLE a (id INT GENERATED ALWAYS AS IDENTITY (10, 5));", 'a', 2, 10},
		{"CREATE TABLE a (id INT IDENTITY (3, 2));", 'd', 2, 3},
		// PostgreSQL's SeqOptList form is still read.
		{"CREATE TABLE a (id INT GENERATED BY DEFAULT AS IDENTITY (START WITH 7));", 'd', 1, 7},
	} {
		stmt := firstCreateStmt(t, mustParse(t, tc.sql))
		col, ok := stmt.TableElts.Items[0].(*nodes.ColumnDef)
		if !ok {
			t.Fatalf("Parse(%q): expected ColumnDef, got %T", tc.sql, stmt.TableElts.Items[0])
		}
		if col.Constraints == nil || len(col.Constraints.Items) != 1 {
			t.Fatalf("Parse(%q): constraints = %#v, want one", tc.sql, col.Constraints)
		}
		constraint, ok := col.Constraints.Items[0].(*nodes.Constraint)
		if !ok || constraint.Contype != nodes.CONSTR_IDENTITY {
			t.Fatalf("Parse(%q): constraint = %#v, want an identity constraint", tc.sql, col.Constraints.Items[0])
		}
		if constraint.GeneratedWhen != tc.when {
			t.Fatalf("Parse(%q): generated when = %q, want %q", tc.sql, constraint.GeneratedWhen, tc.when)
		}
		if constraint.Options == nil || len(constraint.Options.Items) != tc.options {
			t.Fatalf("Parse(%q): options = %#v, want %d", tc.sql, constraint.Options, tc.options)
		}
		start := findDefElem(constraint.Options, "start")
		if start == nil {
			t.Fatalf("Parse(%q): expected a start option in %#v", tc.sql, constraint.Options)
		}
		value, ok := start.Arg.(*nodes.Integer)
		if !ok || value.Ival != tc.wantStart {
			t.Fatalf("Parse(%q): start = %#v, want %d", tc.sql, start.Arg, tc.wantStart)
		}
	}

	for _, sql := range []string{
		"CREATE TABLE a (id INT IDENTITY (, 1))",
		"CREATE TABLE a (id INT IDENTITY (1,))",
		"CREATE TABLE a (id INT GENERATED BY DEFAULT AS IDENTITY (1,))",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) succeeded, want a syntax error", sql)
		}
	}
}

func mustParse(t *testing.T, sql string) *nodes.List {
	t.Helper()
	tree, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	return tree
}
