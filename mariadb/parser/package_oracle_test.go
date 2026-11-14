package parser

import (
	"strings"
	"testing"
	"time"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// TestPackageFormsOracle holds omni to MariaDB 11.8.8 on its packages: CREATE PACKAGE, DROP
// PACKAGE [BODY], SHOW CREATE PACKAGE [BODY], SHOW PACKAGE [BODY] STATUS and SHOW PACKAGE
// BODY CODE. omni reads no CREATE PACKAGE BODY, which the server refuses (1305) for a package
// that does not exist, as every one here does not.
func TestPackageFormsOracle(t *testing.T) {
	o := startMariaDB(t)
	for _, c := range []struct {
		sql    string
		accept bool
	}{
		{"CREATE PACKAGE pk END", true},
		{"CREATE OR REPLACE PACKAGE pk END", true},
		{"CREATE PACKAGE IF NOT EXISTS pk END", true},
		{"CREATE OR REPLACE PACKAGE IF NOT EXISTS pk END", false},
		{"CREATE DEFINER = root@localhost PACKAGE pk END", true},
		{"CREATE DEFINER = CURRENT_USER() PACKAGE pk END", true},
		{"CREATE DEFINER = root SQL SECURITY INVOKER PACKAGE pk END", false},
		{"CREATE PACKAGE d.pk END", true},
		{"CREATE PACKAGE pk END pk", false},
		{"CREATE PACKAGE `BODY` END", true},
		{"CREATE PACKAGE body END", false},
		{"CREATE TEMPORARY PACKAGE pk END", false},
		{"CREATE AGGREGATE PACKAGE pk END", false},
		{"CREATE PACKAGE pk SQL SECURITY INVOKER COMMENT 'x' SQL SECURITY DEFINER END", true},
		{"CREATE PACKAGE pk SQL SECURITY pk END", false},
		{"CREATE PACKAGE pk LANGUAGE SQL END", false},
		{"CREATE PACKAGE pk DETERMINISTIC END", false},
		{"CREATE PACKAGE pk PROCEDURE p(IN x INT, OUT y INT) COMMENT 'c' SQL SECURITY INVOKER; " +
			"FUNCTION f(x INT) RETURNS INT DETERMINISTIC NO SQL; END", true},
		{"CREATE PACKAGE pk PROCEDURE p() NOT DETERMINISTIC MODIFIES SQL DATA; END", true},
		{"CREATE PACKAGE pk PROCEDURE p(); FUNCTION p() RETURNS INT; END", true},
		{"CREATE PACKAGE pk PROCEDURE p(); PROCEDURE P(); END", false},
		{"CREATE PACKAGE pk FUNCTION f() RETURNS INT; FUNCTION f() RETURNS INT; END", false},
		{"CREATE PACKAGE pk PROCEDURE p() BEGIN END; END", false},
		{"CREATE PACKAGE pk PROCEDURE p; END", false},
		{"CREATE PACKAGE pk FUNCTION f RETURNS INT; END", false},
		{"CREATE PACKAGE pk FUNCTION f() RETURNS INT AGGREGATE; END", false},
		{"CREATE PACKAGE pk PROCEDURE d.p(); END", false},
		{"CREATE PACKAGE pk DECLARE x INT; END", false},
		{"CREATE PACKAGE `package` PROCEDURE `end`(); END", true},
		{"CREATE /* c */ PACKAGE pk PROCEDURE p(); END", true},
		{"CREATE DEFINER = 'a b'@'localhost' PACKAGE pk PROCEDURE p(); END", true},
		{"CREATE OR /* x */ REPLACE DEFINER = CURRENT_USER ( ) PACKAGE pk PROCEDURE p(); END", true},
		{"CREATE DEFINER = 用户 PACKAGE pk PROCEDURE p(); END", true},
		{"CREATE DEFINER = a$b@localhost PACKAGE pk PROCEDURE p(); END", true},
		{"CREATE /*! DEFINER = root */ PACKAGE pk PROCEDURE p(); END", true},
		{"CREATE /*M!100000 OR REPLACE */ PACKAGE pk PROCEDURE p(); END", true},
		{"CREATE PACKAGE pk PROCEDURE p() LANGUAGE SQL; END", true},
		{"CREATE /*! /*! DEFINER = root */ PACKAGE pk PROCEDURE p(); END", true},
		{"CREATE PACKAGE pk PROCEDURE p(x INT, x INT); END", false},
		{"CREATE PACKAGE pk PROCEDURE p(a ROW(x INT), b TYPE OF t.c, c ROW TYPE OF t); END", true},
		{"CREATE PACKAGE pk FUNCTION f(a TYPE OF a) RETURNS INT; END", false},
		{"CREATE DEFINER=root@127.0.0.1 PACKAGE pk PROCEDURE p(); END", true},
		{"CREATE DEFINER = CURRENT_USER( PACKAGE pk END", false},
		{"CREATE DEFINER PACKAGE pk END", false},
		{"CREATE DEFINER = PACKAGE pk END", false},
		{"CREATE DEFINER = root@ PACKAGE pk END", true},
		{"CREATE DEFINER=root@ PACKAGE pk PROCEDURE p(); END", true},
		{"CREATE DEFINER=CURRENT_ROLE() PACKAGE pk PROCEDURE p(); END", true},
		{"CREATE DEFINER=root@package VIEW v AS SELECT 1", true},
		{"CREATE DEFINER = root@ 'h' VIEW v AS SELECT 1", false},
		{"CREATE DEFINER = 'root'@ VIEW v AS SELECT 1", true},
		{"CREATE PACKAGE pk PROCEDURE `p.q`(); END", false},
		{"CREATE PACKAGE pk FUNCTION `f.g`() RETURNS INT; END", false},
		{"CREATE PACKAGE pk PROCEDURE ``(); END", false},
		{"CREATE PACKAGE pk PROCEDURE `p `(); END", false},
		{"CREATE PACKAGE pk PROCEDURE `" + strings.Repeat("a", 65) + "`(); END", false},
		{"CREATE PACKAGE pk FUNCTION f(x INT, X INT) RETURNS INT; END", false},
		{"CREATE PACKAGE pk PROCEDURE p() CONTAINS; END", false},
		{"CREATE PACKAGE pk PROCEDURE p() READS DATA; END", false},
		{"CREATE PACKAGE pk PROCEDURE p() SQL SECURITY garbage; END", false},

		{"DROP PACKAGE pk", true},
		{"DROP PACKAGE IF EXISTS d.pk", true},
		{"DROP PACKAGE BODY IF EXISTS pk", true},
		{"DROP PACKAGE BODY BODY", true},
		{"DROP PACKAGE BODY", false},
		{"DROP PACKAGE IF EXISTS", false},
		{"DROP TEMPORARY PACKAGE pk", false},

		{"SHOW CREATE PACKAGE pk", true},
		{"SHOW CREATE PACKAGE BODY pk", true},
		{"SHOW CREATE PACKAGE BODY", true},
		{"SHOW CREATE PACKAGE BODY.pk", true},
		{"SHOW CREATE PACKAGE BODY d.pk", true},
		{"SHOW PACKAGE STATUS", true},
		{"SHOW PACKAGE STATUS LIKE 'p%'", true},
		{"SHOW PACKAGE BODY STATUS WHERE Db = 'd'", true},
		{"SHOW PACKAGE BODY CODE pk", true},
		{"SHOW PACKAGE CODE pk", false},
		{"SHOW PACKAGE BODY CODE", false},
		{"SHOW PACKAGE", false},
	} {
		if got := o.prepareParses(t, c.sql); got != c.accept {
			t.Fatalf("MariaDB prepareParses(%q)=%v, want %v", c.sql, got, c.accept)
		}
		_, perr := Parse(c.sql)
		if omniAccepts := perr == nil; omniAccepts != c.accept {
			t.Errorf("DIVERGENCE omni accepts=%v, MariaDB accepts=%v for %q: %v",
				omniAccepts, c.accept, c.sql, perr)
		}
	}
}

// TestPackageTree pins what the package statements read into.
func TestPackageTree(t *testing.T) {
	list, err := Parse("CREATE OR REPLACE DEFINER = u PACKAGE d.pk COMMENT 'c' " +
		"PROCEDURE p(IN a INT); FUNCTION f() RETURNS INT DETERMINISTIC; END; " +
		"DROP PACKAGE BODY IF EXISTS pk; SHOW CREATE PACKAGE BODY pk; SHOW PACKAGE BODY STATUS")
	if err != nil {
		t.Fatal(err)
	}
	create, ok := list.Items[0].(*nodes.CreatePackageStmt)
	if !ok {
		t.Fatalf("got %T, want *ast.CreatePackageStmt", list.Items[0])
	}
	if !create.OrReplace || create.Definer == "" || create.Name.Schema != "d" || create.Name.Name != "pk" ||
		len(create.Characteristics) != 1 || len(create.Routines) != 2 {
		t.Errorf("CREATE PACKAGE = %+v", create)
	}
	if p := create.Routines[0]; !p.IsProcedure || p.Name != "p" || len(p.Params) != 1 || p.Returns != nil {
		t.Errorf("procedure = %+v", p)
	}
	if f := create.Routines[1]; f.IsProcedure || f.Name != "f" || f.Returns == nil || len(f.Characteristics) != 1 {
		t.Errorf("function = %+v", f)
	}
	if drop, ok := list.Items[1].(*nodes.DropPackageStmt); !ok || !drop.Body || !drop.IfExists || drop.Name.Name != "pk" {
		t.Errorf("DROP PACKAGE = %+v", list.Items[1])
	}
	if show, ok := list.Items[2].(*nodes.ShowStmt); !ok || show.Type != "CREATE PACKAGE BODY" || show.From.Name != "pk" {
		t.Errorf("SHOW CREATE PACKAGE BODY = %+v", list.Items[2])
	}
	if show, ok := list.Items[3].(*nodes.ShowStmt); !ok || show.Type != "PACKAGE BODY STATUS" {
		t.Errorf("SHOW PACKAGE BODY STATUS = %+v", list.Items[3])
	}
}

// TestSplitPackageWordsStayLinear holds Split to linear time on a statement that names
// package many times: each PACKAGE asks whether the statement's head is CREATE PACKAGE, and
// that question reads a few tokens, not the whole statement again.
func TestSplitPackageWordsStayLinear(t *testing.T) {
	sql := strings.Repeat("SELECT package UNION ALL ", 16000) + "SELECT package"
	start := time.Now()
	if segs := Split(sql); len(segs) != 1 {
		t.Fatalf("Split gave %d segments, want 1", len(segs))
	}
	if d := time.Since(start); d > 5*time.Second {
		t.Fatalf("Split of 16001 package words took %v, over 5s", d)
	}
}

// TestSplitPackageExecOpenersStayLinear holds the CREATE PACKAGE head check to linear time
// on executable-comment openers inside one, which the lexer drops rather than nests.
func TestSplitPackageExecOpenersStayLinear(t *testing.T) {
	sql := "CREATE " + strings.Repeat("/*! ", 16000) + " DEFINER = root */ PACKAGE pk PROCEDURE p(); END; SELECT 1"
	start := time.Now()
	if segs := Split(sql); len(segs) != 2 {
		t.Fatalf("Split gave %d segments, want 2", len(segs))
	}
	if d := time.Since(start); d > 5*time.Second {
		t.Fatalf("Split of 16000 nested openers took %v, over 5s", d)
	}
}
