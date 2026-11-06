package catalog

import (
	"errors"
	"fmt"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"

	nodes "github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/parser"
)

// viewCollationTable is the table every viewCollationCases body reads.
const viewCollationTable = `CREATE TABLE t (
	a VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin,
	g VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
	u VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
	z VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
	l VARCHAR(10) CHARACTER SET latin1 COLLATE latin1_bin,
	s VARCHAR(10) CHARACTER SET ascii,
	vb VARBINARY(10),
	i INT,
	d DATE
)`

// viewCollationCases records TiDB v8.5.5's answer for view bodies over
// viewCollationTable, with collation_connection utf8mb4_bin (the server
// default): SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) over the view, or
// the error CREATE VIEW answers. view_collation.go is the model; these are the
// cases its rules were read against.
var viewCollationCases = []struct {
	body      string
	charset   string
	collation string
	coer      int
	code      int
	msg       string
}{
	// Same charset: a _bin collation wins, two others give NONE.
	{body: "SELECT CONCAT(a, g) AS c FROM t", charset: "utf8mb4", collation: "utf8mb4_bin", coer: 2},
	{body: "SELECT CONCAT(g, u) AS c FROM t", charset: "utf8mb4", collation: "utf8mb4_bin", coer: 1},
	{body: "SELECT CONCAT(z, g) AS c FROM t", charset: "utf8mb4", collation: "utf8mb4_bin", coer: 1},
	{body: "SELECT CONCAT(g, u, a) AS c FROM t", charset: "utf8mb4", collation: "utf8mb4_bin", coer: 1},
	{body: "SELECT CONCAT(CONCAT(g, u), l) AS c FROM t", charset: "utf8mb4", collation: "utf8mb4_bin", coer: 1},
	{body: "SELECT REPEAT(CONCAT(g, u), 2) AS c FROM t", charset: "utf8mb4", collation: "utf8mb4_bin", coer: 1},
	{body: "SELECT IF(1, g, u) AS c FROM t", charset: "utf8mb4", collation: "utf8mb4_bin", coer: 1},
	// Different charsets: Unicode wins at one coercibility.
	{body: "SELECT CONCAT(a, l) AS c FROM t", charset: "utf8mb4", collation: "utf8mb4_bin", coer: 2},
	{body: "SELECT CONCAT(l, a) AS c FROM t", charset: "utf8mb4", collation: "utf8mb4_bin", coer: 2},
	{body: "SELECT IFNULL(g, l) AS c FROM t", charset: "utf8mb4", collation: "utf8mb4_general_ci", coer: 2},
	{body: "SELECT CONCAT(s, g) AS c FROM t", charset: "utf8mb4", collation: "utf8mb4_general_ci", coer: 2},
	{body: "SELECT CONCAT(s, l) AS c FROM t", charset: "latin1", collation: "latin1_bin", coer: 2},
	// Literals are COERCIBLE, with or without an introducer.
	{body: "SELECT CONCAT(l, 'é') AS c FROM t", charset: "latin1", collation: "latin1_bin", coer: 2},
	{body: "SELECT CONCAT('é', l) AS c FROM t", charset: "latin1", collation: "latin1_bin", coer: 2},
	{body: "SELECT CONCAT(l, _utf8mb4'x') AS c FROM t", charset: "latin1", collation: "latin1_bin", coer: 2},
	{body: "SELECT CONCAT(_latin1'x', g) AS c FROM t", charset: "utf8mb4", collation: "utf8mb4_general_ci", coer: 2},
	{body: "SELECT CONCAT(l, CONCAT('é', 'x')) AS c FROM t", charset: "latin1", collation: "latin1_bin", coer: 2},
	{body: "SELECT LPAD(g, 5, 'x') AS c FROM t", charset: "utf8mb4", collation: "utf8mb4_general_ci", coer: 2},
	{body: "SELECT CONCAT(s, 'é') AS c FROM t", code: 1267,
		msg: "Illegal mix of collations (ascii_bin,IMPLICIT) and (utf8mb4_bin,COERCIBLE) for operation 'concat'"},
	// Numbers, NULL and binary strings.
	{body: "SELECT CONCAT(1, 2) AS c FROM t", charset: "utf8mb4", collation: "utf8mb4_bin", coer: 4},
	{body: "SELECT CONCAT(g, 1) AS c FROM t", charset: "utf8mb4", collation: "utf8mb4_general_ci", coer: 2},
	{body: "SELECT CONCAT(i, g) AS c FROM t", charset: "utf8mb4", collation: "utf8mb4_general_ci", coer: 2},
	{body: "SELECT CONCAT(d, g) AS c FROM t", charset: "utf8mb4", collation: "utf8mb4_general_ci", coer: 2},
	{body: "SELECT CONCAT(NULL, g) AS c FROM t", charset: "utf8mb4", collation: "utf8mb4_general_ci", coer: 2},
	{body: "SELECT CONCAT(vb, g) AS c FROM t", charset: "binary", collation: "binary", coer: 2},
	{body: "SELECT CONCAT(g, vb) AS c FROM t", charset: "binary", collation: "binary", coer: 2},
	{body: "SELECT CONCAT(g = l, s) AS c FROM t", charset: "ascii", collation: "ascii_bin", coer: 2},
	// Functions that take the collation of some arguments only.
	{body: "SELECT LOWER(g) AS c FROM t", charset: "utf8mb4", collation: "utf8mb4_general_ci", coer: 2},
	{body: "SELECT TRIM(BOTH 'x' FROM g) AS c FROM t", charset: "utf8mb4", collation: "utf8mb4_general_ci", coer: 2},
	{body: "SELECT LPAD(l, 20, g) AS c FROM t", charset: "utf8mb4", collation: "utf8mb4_general_ci", coer: 2},
	{body: "SELECT LPAD(g, 20, u) AS c FROM t", charset: "utf8mb4", collation: "utf8mb4_bin", coer: 1},
	// COLLATE and CONVERT.
	{body: "SELECT CONCAT(g, u COLLATE utf8mb4_general_ci) AS c FROM t", charset: "utf8mb4", collation: "utf8mb4_general_ci", coer: 0},
	{body: "SELECT CONCAT(l COLLATE latin1_bin, 'é') AS c FROM t", charset: "latin1", collation: "latin1_bin", coer: 0},
	{body: "SELECT CONVERT(g USING latin1) AS c FROM t", charset: "latin1", collation: "latin1_bin", coer: 2},
	{body: "SELECT CONCAT(l, CONVERT(g USING latin1)) AS c FROM t", charset: "latin1", collation: "latin1_bin", coer: 2},
	{body: "SELECT CONCAT(u, CONVERT(g USING utf8mb4)) AS c FROM t", charset: "utf8mb4", collation: "utf8mb4_bin", coer: 2},
	{body: "SELECT CONCAT(g COLLATE utf8mb4_general_ci, u COLLATE utf8mb4_bin) AS c FROM t", code: 1267,
		msg: "Illegal mix of collations (utf8mb4_general_ci,EXPLICIT) and (utf8mb4_bin,EXPLICIT) for operation 'concat'"},
	{body: "SELECT LENGTH(CONCAT(g COLLATE utf8mb4_general_ci, u COLLATE utf8mb4_bin)) AS c FROM t", code: 1267,
		msg: "Illegal mix of collations (utf8mb4_general_ci,EXPLICIT) and (utf8mb4_bin,EXPLICIT) for operation 'concat'"},
	{body: "SELECT CONCAT(l COLLATE latin1_bin, g) AS c FROM t", code: 1267,
		msg: "Illegal mix of collations (latin1_bin,EXPLICIT) and (utf8mb4_general_ci,IMPLICIT) for operation 'concat'"},
	// Comparisons refuse NONE.
	{body: "SELECT a AS c FROM t WHERE a = g", charset: "utf8mb4", collation: "utf8mb4_bin", coer: 2},
	{body: "SELECT a AS c FROM t WHERE l = g", charset: "utf8mb4", collation: "utf8mb4_bin", coer: 2},
	{body: "SELECT a AS c FROM t WHERE a = 'x' OR i = g OR d = g OR vb = g", charset: "utf8mb4", collation: "utf8mb4_bin", coer: 2},
	{body: "SELECT a AS c FROM t WHERE g = u", code: 1267,
		msg: "Illegal mix of collations (utf8mb4_general_ci,IMPLICIT) and (utf8mb4_unicode_ci,IMPLICIT) for operation '='"},
	{body: "SELECT a AS c FROM t WHERE z = g", code: 1267,
		msg: "Illegal mix of collations (utf8mb4_0900_ai_ci,IMPLICIT) and (utf8mb4_general_ci,IMPLICIT) for operation '='"},
	{body: "SELECT a AS c FROM t WHERE g <> u", code: 1267,
		msg: "Illegal mix of collations (utf8mb4_general_ci,IMPLICIT) and (utf8mb4_unicode_ci,IMPLICIT) for operation '!='"},
	{body: "SELECT a AS c FROM t WHERE g <=> u", code: 1267,
		msg: "Illegal mix of collations (utf8mb4_general_ci,IMPLICIT) and (utf8mb4_unicode_ci,IMPLICIT) for operation '<=>'"},
	{body: "SELECT a AS c FROM t WHERE g < u", code: 1267,
		msg: "Illegal mix of collations (utf8mb4_general_ci,IMPLICIT) and (utf8mb4_unicode_ci,IMPLICIT) for operation '<'"},
	{body: "SELECT a AS c FROM t WHERE g NOT LIKE u", code: 1267,
		msg: "Illegal mix of collations (utf8mb4_general_ci,IMPLICIT) and (utf8mb4_unicode_ci,IMPLICIT) for operation 'like'"},
	{body: "SELECT a AS c FROM t WHERE STRCMP(g, u) = 0", code: 1267,
		msg: "Illegal mix of collations (utf8mb4_general_ci,IMPLICIT) and (utf8mb4_unicode_ci,IMPLICIT) for operation 'strcmp'"},
	{body: "SELECT a AS c FROM t WHERE CONCAT(g, u) = CONCAT(g, u)", code: 1267,
		msg: "Illegal mix of collations (utf8mb4_bin,NONE) and (utf8mb4_bin,NONE) for operation '='"},
	{body: "SELECT a AS c FROM t WHERE COALESCE(g, u, z) = 'x'", code: 1267,
		msg: "Illegal mix of collations (utf8mb4_bin,NONE) and (utf8mb4_bin,COERCIBLE) for operation '='"},
	{body: "SELECT t.a AS c FROM t JOIN t AS t2 ON t.g = t2.u", code: 1267,
		msg: "Illegal mix of collations (utf8mb4_general_ci,IMPLICIT) and (utf8mb4_unicode_ci,IMPLICIT) for operation '='"},
}

// TestViewCollationRecord holds the catalog to viewCollationCases without a
// server: the error CREATE VIEW answers, or the charset and collation it keeps
// in View.ColumnMetadata and the coercibility it derives.
func TestViewCollationRecord(t *testing.T) {
	for _, tc := range viewCollationCases {
		t.Run(tc.body, func(t *testing.T) {
			c := setupWithDB(t)
			mustExec(t, c, viewCollationTable)
			results, err := c.Exec("CREATE VIEW v AS "+tc.body, nil)
			if err == nil && len(results) > 0 {
				err = results[0].Error
			}
			if tc.code != 0 {
				var catErr *Error
				if !errors.As(err, &catErr) {
					t.Fatalf("got %v, want error %d %q", err, tc.code, tc.msg)
				}
				if catErr.Code != tc.code || catErr.Message != tc.msg {
					t.Errorf("got %d %q, want %d %q", catErr.Code, catErr.Message, tc.code, tc.msg)
				}
				if c.GetDatabase("testdb").Views["v"] != nil {
					t.Error("the refused view was created")
				}
				return
			}
			if err != nil {
				t.Fatalf("CREATE VIEW: %v", err)
			}
			v := c.GetDatabase("testdb").Views["v"]
			if v == nil || len(v.ColumnMetadata) != 1 {
				t.Fatalf("view v: %+v", v)
			}
			if got := v.ColumnMetadata[0]; got.Charset != tc.charset || got.Collation != tc.collation {
				t.Errorf("got %s %s, want %s %s", got.Charset, got.Collation, tc.charset, tc.collation)
			}

			list, err := parser.Parse(tc.body)
			if err != nil {
				t.Fatal(err)
			}
			sel := list.Items[0].(*nodes.SelectStmt)
			ecs := c.viewTargetCollations(sel, c.GetDatabase("testdb"))
			if len(ecs) != 1 || !ecs[0].known {
				t.Fatalf("no derivation: %+v", ecs)
			}
			if int(ecs[0].coer) != tc.coer {
				t.Errorf("coercibility %d, want %d", ecs[0].coer, tc.coer)
			}
		})
	}
}

// TestContainerViewCollationRecord asks TiDB that viewCollationCases is its
// answer.
func TestContainerViewCollationRecord(t *testing.T) {
	tc := startTiDBForCatalog(t)
	const db = "omni_view_collation"
	var prevCollation string
	if err := tc.db.QueryRowContext(tc.ctx, "SELECT @@collation_connection").Scan(&prevCollation); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = tc.db.ExecContext(tc.ctx, "DROP DATABASE IF EXISTS "+db)
		_, _ = tc.db.ExecContext(tc.ctx, "SET NAMES utf8mb4 COLLATE "+prevCollation)
	})
	// The pool holds one connection (startTiDBForCatalog), so the session
	// settings below hold for every statement.
	mustExecTiDB(t, tc, "SET NAMES utf8mb4 COLLATE utf8mb4_bin")
	mustExecTiDB(t, tc, "DROP DATABASE IF EXISTS "+db)
	mustExecTiDB(t, tc, "CREATE DATABASE "+db)
	mustExecTiDB(t, tc, "USE "+db)
	mustExecTiDB(t, tc, viewCollationTable)
	mustExecTiDB(t, tc, "INSERT INTO t VALUES ('x', 'x', 'x', 'x', 'x', 'x', 'x', 1, '2020-01-01')")

	for _, c := range viewCollationCases {
		t.Run(c.body, func(t *testing.T) {
			mustExecTiDB(t, tc, "DROP VIEW IF EXISTS v")
			_, err := tc.db.ExecContext(tc.ctx, "CREATE VIEW v AS "+c.body)
			if c.code != 0 {
				var myErr *gomysql.MySQLError
				if !errors.As(err, &myErr) {
					t.Fatalf("got %v, want error %d", err, c.code)
				}
				if int(myErr.Number) != c.code || myErr.Message != c.msg {
					t.Errorf("got %d %q, want %d %q", myErr.Number, myErr.Message, c.code, c.msg)
				}
				return
			}
			if err != nil {
				t.Fatalf("CREATE VIEW: %v", err)
			}
			var cs, coll string
			var coer int
			if err := tc.db.QueryRowContext(tc.ctx, "SELECT CHARSET(c), COLLATION(c), COERCIBILITY(c) FROM v").Scan(&cs, &coll, &coer); err != nil {
				t.Fatal(err)
			}
			if got, want := fmt.Sprint(cs, " ", coll, " ", coer), fmt.Sprint(c.charset, " ", c.collation, " ", c.coer); got != want {
				t.Errorf("got %s, want %s", got, want)
			}
		})
	}
}
