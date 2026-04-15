package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// unwrapGrantRole extracts the GrantRoleStmt from a single-stmt parse result.
func unwrapGrantRole(t *testing.T, sql string) *nodes.GrantRoleStmt {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if list == nil || len(list.Items) != 1 {
		t.Fatalf("Parse(%q): expected 1 stmt, got %d", sql, len(list.Items))
	}
	raw, ok := list.Items[0].(*nodes.RawStmt)
	if !ok {
		t.Fatalf("Parse(%q): item[0] not RawStmt", sql)
	}
	stmt, ok := raw.Stmt.(*nodes.GrantRoleStmt)
	if !ok {
		t.Fatalf("Parse(%q): stmt not GrantRoleStmt, got %T", sql, raw.Stmt)
	}
	return stmt
}

func defElemName(n nodes.Node) string {
	if de, ok := n.(*nodes.DefElem); ok {
		return de.Defname
	}
	return ""
}

func defElemBool(n nodes.Node) (bool, bool) {
	de, ok := n.(*nodes.DefElem)
	if !ok || de.Arg == nil {
		return false, false
	}
	b, ok := de.Arg.(*nodes.Boolean)
	if !ok {
		return false, false
	}
	return b.Boolval, true
}

func TestGrantRoleOptionsAndGrantedBy(t *testing.T) {
	t.Run("WITH ADMIN OPTION (baseline)", func(t *testing.T) {
		g := unwrapGrantRole(t, "GRANT r1 TO r2 WITH ADMIN OPTION")
		if !g.IsGrant {
			t.Fatalf("expected IsGrant=true")
		}
		if g.Opt == nil || len(g.Opt.Items) != 1 {
			t.Fatalf("expected 1 opt, got %v", g.Opt)
		}
		if defElemName(g.Opt.Items[0]) != "admin" {
			t.Fatalf("expected admin opt, got %q", defElemName(g.Opt.Items[0]))
		}
		if b, ok := defElemBool(g.Opt.Items[0]); !ok || !b {
			t.Fatalf("expected admin=true, got %v/%v", b, ok)
		}
		if g.Grantor != nil {
			t.Fatalf("expected nil Grantor")
		}
	})

	t.Run("no WITH (baseline)", func(t *testing.T) {
		g := unwrapGrantRole(t, "GRANT r1 TO r2")
		if g.Opt != nil {
			t.Fatalf("expected nil Opt, got %v", g.Opt)
		}
		if g.Grantor != nil {
			t.Fatalf("expected nil Grantor")
		}
	})

	t.Run("WITH INHERIT FALSE, ADMIN TRUE", func(t *testing.T) {
		g := unwrapGrantRole(t, "GRANT r1 TO r2 WITH INHERIT FALSE, ADMIN TRUE")
		if g.Opt == nil || len(g.Opt.Items) != 2 {
			t.Fatalf("expected 2 opts, got %v", g.Opt)
		}
		if defElemName(g.Opt.Items[0]) != "inherit" {
			t.Fatalf("[0] expected inherit, got %q", defElemName(g.Opt.Items[0]))
		}
		if b, _ := defElemBool(g.Opt.Items[0]); b {
			t.Fatalf("[0] expected false, got true")
		}
		if defElemName(g.Opt.Items[1]) != "admin" {
			t.Fatalf("[1] expected admin, got %q", defElemName(g.Opt.Items[1]))
		}
		if b, _ := defElemBool(g.Opt.Items[1]); !b {
			t.Fatalf("[1] expected true, got false")
		}
	})

	t.Run("WITH INHERIT TRUE, SET FALSE", func(t *testing.T) {
		g := unwrapGrantRole(t, "GRANT r1 TO r2 WITH INHERIT TRUE, SET FALSE")
		if g.Opt == nil || len(g.Opt.Items) != 2 {
			t.Fatalf("expected 2 opts, got %v", g.Opt)
		}
		if defElemName(g.Opt.Items[0]) != "inherit" || defElemName(g.Opt.Items[1]) != "set" {
			t.Fatalf("unexpected names: %q, %q", defElemName(g.Opt.Items[0]), defElemName(g.Opt.Items[1]))
		}
	})

	t.Run("GRANTED BY r3", func(t *testing.T) {
		g := unwrapGrantRole(t, "GRANT r1 TO r2 GRANTED BY r3")
		if g.Opt != nil {
			t.Fatalf("expected nil Opt, got %v", g.Opt)
		}
		if g.Grantor == nil || g.Grantor.Rolename != "r3" {
			t.Fatalf("expected grantor r3, got %+v", g.Grantor)
		}
	})

	t.Run("WITH ADMIN OPTION GRANTED BY r3", func(t *testing.T) {
		g := unwrapGrantRole(t, "GRANT r1 TO r2 WITH ADMIN OPTION GRANTED BY r3")
		if g.Opt == nil || len(g.Opt.Items) != 1 {
			t.Fatalf("expected 1 opt, got %v", g.Opt)
		}
		if defElemName(g.Opt.Items[0]) != "admin" {
			t.Fatalf("expected admin opt")
		}
		if g.Grantor == nil || g.Grantor.Rolename != "r3" {
			t.Fatalf("expected grantor r3, got %+v", g.Grantor)
		}
	})

	t.Run("WITH ADMIN TRUE GRANTED BY r3", func(t *testing.T) {
		g := unwrapGrantRole(t, "GRANT r1 TO r2 WITH ADMIN TRUE GRANTED BY r3")
		if g.Opt == nil || len(g.Opt.Items) != 1 {
			t.Fatalf("expected 1 opt, got %v", g.Opt)
		}
		if g.Grantor == nil || g.Grantor.Rolename != "r3" {
			t.Fatalf("expected grantor r3, got %+v", g.Grantor)
		}
	})
}
