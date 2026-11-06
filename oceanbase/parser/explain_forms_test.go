package parser

import (
	"regexp"
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// EXPLAIN's modifier list and its three slots -- the format name, the
// connection id and the explained statement -- as OceanBase reads them.
//
// The list is neither MySQL 8.0's nor the one this parser inherited from the
// MySQL tree it was forked from, and it is the narrowest in this family.
// EXTENDED and PARTITIONS are still here, as they are in MySQL 5.7, whose
// protocol version this server reports. The format name is a closed list of
// two bare words and is never quoted, where MySQL reads any ident_or_text and
// TiDB reads a closed list plus any string. There is no EXPLAIN ANALYZE and no
// FOR CONNECTION at all, and TABLE is no statement, so four forms this parser
// used to take are 1064 here.
//
// The rejections are the point of this file: a parser that accepts what the
// engine will not lets a review pass a statement that then fails on the
// server. Every case was measured against 5.7.25-OceanBase_CE-v4.3.5.6 and
// TestExplainFormsMatchOceanBase holds this file to the engine; the tables
// below are also checked against the parser on their own.
//
// Ref: https://en.oceanbase.com/docs (EXPLAIN); measured on
// 5.7.25-OceanBase_CE-v4.3.5.6.

// explainFormCase is a statement and whether the engine's parser accepts it.
type explainFormCase struct {
	sql    string
	accept bool
}

var explainFormCases = []explainFormCase{
	// The format name is TRADITIONAL or JSON, bare. A quoted name is a 1064
	// here, where MySQL and TiDB both read one, and so is every other word.
	{"EXPLAIN FORMAT = TRADITIONAL SELECT 1", true},
	{"EXPLAIN FORMAT = JSON SELECT 1", true},
	{"EXPLAIN FORMAT = traditional SELECT 1", true},
	{"EXPLAIN FORMAT = json SELECT 1", true},
	{"EXPLAIN FORMAT = 'traditional' SELECT 1", false},
	{"EXPLAIN FORMAT = 'TRADITIONAL' SELECT 1", false},
	{`EXPLAIN FORMAT = "traditional" SELECT 1`, false},
	{"EXPLAIN FORMAT = `TRADITIONAL` SELECT 1", false},
	{"EXPLAIN FORMAT = 'json' SELECT 1", false},
	{"EXPLAIN FORMAT = '' SELECT 1", false},
	{"EXPLAIN FORMAT = TREE SELECT 1", false},
	{"EXPLAIN FORMAT = ROW SELECT 1", false},
	{"EXPLAIN FORMAT = DOT SELECT 1", false},
	{"EXPLAIN FORMAT = BASIC SELECT 1", false},
	{"EXPLAIN FORMAT = EXTENDED SELECT 1", false},
	{"EXPLAIN FORMAT = EXTENDED_NOADDR SELECT 1", false},
	{"EXPLAIN FORMAT = PLANREGRESS SELECT 1", false},
	{"EXPLAIN FORMAT = nosuchformat SELECT 1", false},
	{"EXPLAIN FORMAT = binary SELECT 1", false},
	{"EXPLAIN FORMAT = select SELECT 1", false},
	{"EXPLAIN FORMAT = 1 SELECT 1", false},
	// The value is required, the '=' is required, and one FORMAT is the most
	// a statement takes.
	{"EXPLAIN FORMAT = SELECT 1", false},
	{"EXPLAIN FORMAT SELECT 1", false},
	{"EXPLAIN FORMAT 'json' SELECT 1", false},
	{"EXPLAIN FORMAT = JSON FORMAT = TRADITIONAL SELECT 1", false},

	// EXTENDED and PARTITIONS, which MySQL 8.0 and TiDB both dropped. They
	// are one choice: neither stands with the other and neither stands with
	// a FORMAT.
	{"EXPLAIN EXTENDED SELECT 1", true},
	{"EXPLAIN PARTITIONS SELECT 1", true},
	{"EXPLAIN EXTENDED SELECT * FROM t1", true},
	{"EXPLAIN PARTITIONS UPDATE t1 SET a = 1", true},
	{"EXPLAIN EXTENDED FORMAT = TRADITIONAL SELECT 1", false},
	{"EXPLAIN FORMAT = TRADITIONAL EXTENDED SELECT 1", false},
	// EXPLAIN EXTENDED PARTITIONS SELECT 1 and EXPLAIN PARTITIONS EXTENDED
	// SELECT 1 belong here too, and are not: the two words read as a table
	// name and a column name, and this parser then takes the SELECT that
	// follows as a second statement with no separator in front of it, which
	// every engine in this family does and no server does. That is a defect
	// of its own -- SELECT 1 SELECT 2 parses here too -- and not one of
	// EXPLAIN.

	// There is no EXPLAIN ANALYZE, which MySQL and TiDB both take.
	{"EXPLAIN ANALYZE SELECT 1", false},
	{"EXPLAIN ANALYZE FORMAT = TRADITIONAL SELECT 1", false},
	{"EXPLAIN ANALYZE FORMAT = 'traditional' SELECT 1", false},
	{"EXPLAIN ANALYZE UPDATE t1 SET a = 1", false},
	{"EXPLAIN ANALYZE t1", false},
	{"DESC ANALYZE SELECT 1", false},

	// ...and no FOR CONNECTION either, in any spelling.
	{"EXPLAIN FOR CONNECTION 1", false},
	{"EXPLAIN FOR connection 42", false},
	{"EXPLAIN FORMAT = TRADITIONAL FOR CONNECTION 1", false},
	{"EXPLAIN FORMAT = JSON FOR CONNECTION 1", false},
	{"EXPLAIN EXTENDED FOR CONNECTION 1", false},
	{"EXPLAIN FOR CONNECTION", false},
	{"EXPLAIN FOR 1", false},
	{"EXPLAIN FOR CONNECTION 'x'", false},
	{"EXPLAIN ANALYZE FOR CONNECTION 1", false},

	// The explainable statement.
	{"EXPLAIN SELECT 1", true},
	{"EXPLAIN INSERT INTO t1 VALUES (1)", true},
	{"EXPLAIN REPLACE INTO t1 VALUES (1)", true},
	{"EXPLAIN UPDATE t1 SET a = 1", true},
	{"EXPLAIN DELETE FROM t1", true},
	{"EXPLAIN VALUES ROW(1)", true},
	{"EXPLAIN VALUES ROW(1), ROW(2) ORDER BY 1", true},
	{"EXPLAIN FORMAT = TRADITIONAL UPDATE t1 SET a = 1", true},
	{"EXPLAIN FORMAT = TRADITIONAL DELETE FROM t1", true},
	{"EXPLAIN FORMAT = TRADITIONAL INSERT INTO t1 VALUES (1)", true},
	{"EXPLAIN FORMAT = TRADITIONAL REPLACE INTO t1 VALUES (1)", true},
	{"EXPLAIN (SELECT 1)", true},
	{"EXPLAIN WITH c AS (SELECT 1) SELECT * FROM c", true},
	{"EXPLAIN FORMAT = TRADITIONAL (SELECT 1)", true},
	{"EXPLAIN FORMAT = TRADITIONAL WITH c AS (SELECT 1) SELECT * FROM c", true},
	// TABLE is no statement here, and neither is MySQL 8.0's VALUES (1)
	// without ROW.
	{"EXPLAIN TABLE t1", false},
	{"EXPLAIN FORMAT = TRADITIONAL TABLE t1", false},
	{"EXPLAIN VALUES (1)", false},
	// ...and no DDL at all, which is where TiDB parts: it explains ALTER
	// TABLE and OceanBase explains nothing.
	{"EXPLAIN ALTER TABLE t1 ADD INDEX (a)", false},
	{"EXPLAIN ALTER TABLE t1 ADD a varchar(255)", false},
	{"EXPLAIN FORMAT = TRADITIONAL ALTER TABLE t1 ADD INDEX (a)", false},
	{"EXPLAIN CREATE TABLE t2 (a int)", false},
	{"EXPLAIN DROP TABLE t1", false},
	{"EXPLAIN SET @a = 1", false},
	{"EXPLAIN FORMAT = TRADITIONAL SET @a = 1", false},
	{"EXPLAIN FORMAT = TRADITIONAL SHOW TABLES", false},

	// The table form: a table name, then a column name or a wild pattern.
	// The string is a pattern here, as on MySQL, which TiDB refuses, and the
	// qualified name TiDB takes is a 1064 here.
	{"EXPLAIN t1", true},
	{"EXPLAIN t1 a", true},
	{"EXPLAIN t1 `a`", true},
	{"EXPLAIN t1 'abc'", true},
	{`EXPLAIN t1 "abc"`, true},
	{"EXPLAIN t1 '%a%'", true},
	{"EXPLAIN d1.t1 a", true},
	{"DESC t1 a", true},
	{"DESC t1 'abc'", true},
	{"EXPLAIN BEGIN", true},
	{"EXPLAIN t1 a.b", false},
	{"EXPLAIN t1 123", false},
	{"EXPLAIN t1 select", false},
	{"EXPLAIN t1 a b", false},
	// A modifier rules the table form out, so a table name after one is a
	// 1064 and so is a column name after that. EXTENDED and PARTITIONS are
	// not reserved, so each is also a table name of its own, and what
	// follows the word decides which it was.
	{"EXPLAIN FORMAT = TRADITIONAL t1", false},
	{"EXPLAIN FORMAT = TRADITIONAL t1 a", false},
	{"EXPLAIN EXTENDED t1", true},
	{"EXPLAIN PARTITIONS t1", true},
	{"EXPLAIN EXTENDED", true},
	// A statement in a string is TiDB grammar, and not OceanBase's at all.
	{"EXPLAIN 'sqldigest'", false},
}

// explainFormPendingGrammar is the additive side: EXPLAIN modifiers OceanBase
// parses that this parser does not read yet. BASIC and OUTLINE are output
// modes of their own, not format names -- EXPLAIN BASIC SELECT 1 parses and
// FORMAT = BASIC does not -- and reading them means an AST field each, so they
// are recorded here rather than guessed at.
//
// TestExplainFormPendingGrammar pins both halves: the engine must still accept
// each one and this parser must still reject it, so the day either side
// changes, this list is what fails.
var explainFormPendingGrammar = []string{
	"EXPLAIN BASIC SELECT 1",
	"EXPLAIN OUTLINE SELECT 1",
}

func TestExplainFormsAccepted(t *testing.T) {
	for _, c := range explainFormCases {
		if !c.accept {
			continue
		}
		t.Run(c.sql, func(t *testing.T) {
			if _, err := Parse(c.sql); err != nil {
				t.Fatalf("Parse(%q) error: %v", c.sql, err)
			}
		})
	}
}

func TestExplainFormsRejected(t *testing.T) {
	for _, c := range explainFormCases {
		if c.accept {
			continue
		}
		t.Run(c.sql, func(t *testing.T) {
			if _, err := Parse(c.sql); err == nil {
				t.Fatalf("Parse(%q) was accepted, but OceanBase 4.3.5 rejects it", c.sql)
			}
		})
	}
}

func TestExplainFormPendingGrammarRejected(t *testing.T) {
	for _, sql := range explainFormPendingGrammar {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err == nil {
				t.Fatalf("Parse(%q) was accepted: clear it from explainFormPendingGrammar in the same change", sql)
			}
		})
	}
}

// locs matches the position every dumped node carries, so two dumps of the
// same statement can be compared when one of them was parsed after a prefix.
var locs = regexp.MustCompile(` :loc \d+`)

// parseOneExplain parses sql and returns the ExplainStmt it is.
func parseOneExplain(t *testing.T, sql string) *nodes.ExplainStmt {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q) error: %v", sql, err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("Parse(%q) returned %d statements, want 1", sql, len(list.Items))
	}
	stmt, ok := list.Items[0].(*nodes.ExplainStmt)
	if !ok {
		t.Fatalf("Parse(%q) returned %T, want *ast.ExplainStmt", sql, list.Items[0])
	}
	return stmt
}

// TestExplainFormsAST asserts the node each form parses into, not only that it
// parses. A modifier that lands in the wrong field, or an explained statement
// that is no longer the node it is at statement level, parses clean and
// misleads every layer above the parser.
func TestExplainFormsAST(t *testing.T) {
	t.Run("format name reaches the AST as written", func(t *testing.T) {
		for _, tc := range []struct{ sql, want string }{
			{"EXPLAIN FORMAT = JSON SELECT 1", "JSON"},
			{"EXPLAIN FORMAT = json SELECT 1", "json"},
			{"EXPLAIN FORMAT = TRADITIONAL SELECT 1", "TRADITIONAL"},
		} {
			stmt := parseOneExplain(t, tc.sql)
			if stmt.Format != tc.want {
				t.Errorf("%s: Format = %q, want %q", tc.sql, stmt.Format, tc.want)
			}
			if stmt.Extended || stmt.Partitions {
				t.Errorf("%s: Extended = %v, Partitions = %v, want both false",
					tc.sql, stmt.Extended, stmt.Partitions)
			}
			if _, ok := stmt.Stmt.(*nodes.SelectStmt); !ok {
				t.Errorf("%s: Stmt = %T, want *ast.SelectStmt", tc.sql, stmt.Stmt)
			}
		}
	})

	t.Run("extended and partitions", func(t *testing.T) {
		ext := parseOneExplain(t, "EXPLAIN EXTENDED SELECT 1")
		if !ext.Extended || ext.Partitions || ext.Format != "" {
			t.Errorf("EXTENDED: Extended = %v, Partitions = %v, Format = %q",
				ext.Extended, ext.Partitions, ext.Format)
		}
		part := parseOneExplain(t, "EXPLAIN PARTITIONS SELECT 1")
		if !part.Partitions || part.Extended || part.Format != "" {
			t.Errorf("PARTITIONS: Partitions = %v, Extended = %v, Format = %q",
				part.Partitions, part.Extended, part.Format)
		}
		// The same word as a table name sets no flag.
		table := parseOneExplain(t, "EXPLAIN EXTENDED t1")
		if table.Extended {
			t.Error("EXPLAIN EXTENDED t1 set Extended, but EXTENDED is the table name there")
		}
		show, ok := table.Stmt.(*nodes.ShowStmt)
		if !ok {
			t.Fatalf("EXPLAIN EXTENDED t1: Stmt = %T, want *ast.ShowStmt", table.Stmt)
		}
		if show.From == nil || !strings.EqualFold(show.From.Name, "extended") {
			t.Errorf("EXPLAIN EXTENDED t1: table = %+v, want extended", show.From)
		}
		// OceanBase has no EXPLAIN ANALYZE and no FOR CONNECTION, so no
		// statement can set either.
		for _, c := range explainFormCases {
			if !c.accept {
				continue
			}
			stmt := parseOneExplain(t, c.sql)
			if stmt.Analyze {
				t.Errorf("%q set Analyze, which OceanBase's EXPLAIN has no syntax for", c.sql)
			}
			if stmt.ForConnection != 0 {
				t.Errorf("%q set ForConnection, which OceanBase's EXPLAIN has no syntax for", c.sql)
			}
		}
	})

	t.Run("table form holds a column name or a pattern", func(t *testing.T) {
		stmt := parseOneExplain(t, "EXPLAIN d1.t1 a")
		show, ok := stmt.Stmt.(*nodes.ShowStmt)
		if !ok {
			t.Fatalf("Stmt = %T, want *ast.ShowStmt", stmt.Stmt)
		}
		if show.Type != "COLUMNS" {
			t.Errorf("Type = %q, want COLUMNS", show.Type)
		}
		if show.From == nil || show.From.Schema != "d1" || show.From.Name != "t1" {
			t.Errorf("From = %+v, want d1.t1", show.From)
		}
		col, ok := show.Like.(*nodes.ColumnRef)
		if !ok {
			t.Fatalf("Like = %T, want *ast.ColumnRef", show.Like)
		}
		if col.Column != "a" {
			t.Errorf("column = %q, want a", col.Column)
		}
		// The wild pattern is a string literal, not a column name, so a
		// consumer can tell EXPLAIN t 'a%' from EXPLAIN t a.
		wild := parseOneExplain(t, "EXPLAIN t1 'a%'")
		show = wild.Stmt.(*nodes.ShowStmt)
		lit, ok := show.Like.(*nodes.StringLit)
		if !ok {
			t.Fatalf("Like = %T, want *ast.StringLit", show.Like)
		}
		if lit.Value != "a%" {
			t.Errorf("pattern = %q, want a%%", lit.Value)
		}
	})

	t.Run("a backquoted name is read at its own offset", func(t *testing.T) {
		// The bare format name is matched on its text, so the backquoted
		// form is turned away by reading the source at the token's offset.
		// In a script that offset is not the statement's, and the lexer's
		// base offset is what makes the two agree.
		for _, tc := range []struct {
			sql    string
			accept bool
		}{
			{"SELECT 1; EXPLAIN FORMAT = `json` SELECT 1", false},
			{"SELECT 1; EXPLAIN FORMAT = JSON SELECT 1", true},
			{"SELECT 1;\n/* c */ EXPLAIN FORMAT = `traditional` SELECT 1", false},
			{"SELECT 1;\n/* c */ EXPLAIN FORMAT = TRADITIONAL SELECT 1", true},
		} {
			_, err := Parse(tc.sql)
			if accepts := err == nil; accepts != tc.accept {
				t.Errorf("Parse(%q) accepted=%v, want %v: %v", tc.sql, accepts, tc.accept, err)
			}
		}
	})

	t.Run("the explained statement is the statement", func(t *testing.T) {
		// Each explainable form parses into the same node it is on its own,
		// positions aside. A statement that landed in the wrong field, or
		// read a clause differently because EXPLAIN stood in front of it,
		// shows up here.
		for _, inner := range []string{
			"SELECT 1",
			"SELECT a FROM t1 WHERE a > 1 ORDER BY a LIMIT 3",
			"select c1 from t1 union (select c2 from t2) limit 1, 1",
			"INSERT INTO t1 VALUES (1)",
			"REPLACE INTO t1 VALUES (1)",
			"UPDATE t1 SET a = 1 WHERE b = 2",
			"DELETE FROM t1 WHERE a = 1",
			"VALUES ROW(1)",
			"(SELECT 1)",
			"WITH c AS (SELECT 1) SELECT * FROM c",
		} {
			alone, err := Parse(inner)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", inner, err)
			}
			for _, prefix := range []string{"EXPLAIN ", "EXPLAIN EXTENDED ", "EXPLAIN FORMAT = JSON "} {
				sql := prefix + inner
				stmt := parseOneExplain(t, sql)
				want := locs.ReplaceAllString(nodes.NodeToString(alone.Items[0]), "")
				got := locs.ReplaceAllString(nodes.NodeToString(stmt.Stmt), "")
				if got != want {
					t.Errorf("%q explained a different node than %q:\n got %s\nwant %s", sql, inner, got, want)
				}
				// The explained statement keeps its place in the text it was
				// parsed in, which is what every position a consumer reports
				// is relative to.
				if start := innerStart(stmt.Stmt); start != len(prefix) {
					t.Errorf("%q: inner Loc.Start = %d, want %d", sql, start, len(prefix))
				}
				if stmt.Loc.Start != 0 || stmt.Loc.End != len(sql) {
					t.Errorf("%q: Loc = %d..%d, want 0..%d", sql, stmt.Loc.Start, stmt.Loc.End, len(sql))
				}
			}
		}
	})
}

// innerStart reports the start offset of an explained statement, whatever
// node it is.
func innerStart(n nodes.StmtNode) int {
	dump := nodes.NodeToString(n)
	m := locs.FindString(dump)
	if m == "" {
		return -1
	}
	start := 0
	for _, c := range strings.TrimPrefix(m, " :loc ") {
		start = start*10 + int(c-'0')
	}
	return start
}

// TestExplainFormsMatchOceanBase holds the tables above to the engine itself.
// It runs on the package's own OceanBase (see sharedserver.PrivateOceanBase):
//
//	go test ./oceanbase/parser/ -run Explain -count=1
//
// The verdict comes from the statement's error: 1064 is the parser's
// rejection, and any other server error means the text parsed and failed
// later.
func TestExplainFormsMatchOceanBase(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := t.Context()
	for _, c := range explainFormCases {
		t.Run(c.sql, func(t *testing.T) {
			accepts, err := engineAccepts(ctx, db, c.sql)
			if err != nil {
				t.Fatalf("the engine gave no verdict for %q: %v", c.sql, err)
			}
			if accepts != c.accept {
				t.Errorf("OceanBase accepts=%v for %q, this file records accept=%v", accepts, c.sql, c.accept)
			}
		})
	}
	for _, sql := range explainFormPendingGrammar {
		t.Run(sql, func(t *testing.T) {
			accepts, err := engineAccepts(ctx, db, sql)
			if err != nil {
				t.Fatalf("the engine gave no verdict for %q: %v", sql, err)
			}
			if !accepts {
				t.Errorf("OceanBase now rejects %q: it is no longer grammar this parser owes", sql)
			}
		})
	}
}
