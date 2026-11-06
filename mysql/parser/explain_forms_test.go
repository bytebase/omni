package parser

import (
	"regexp"
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// EXPLAIN's modifier list and its three slots -- the format name, the
// connection id and the explained statement -- as MySQL 8.0 reads them.
//
// The list is not the one this parser used to read, and it is not TiDB's
// either. MySQL 8.0 dropped EXTENDED and PARTITIONS, so each is a 1064 here;
// the format name is ident_or_text, so a quoted name and a name MySQL does not
// have both parse and only a reserved word and an integer do not; ANALYZE
// reaches FOR CONNECTION, which TiDB refuses; and no DDL is explainable, where
// TiDB explains ALTER TABLE.
//
// The rejections are the point of this file. A parser that accepts what the
// engine will not lets a review pass a statement that then fails on the
// server, and nothing else in this package covers a FORMAT with no value, a
// reserved word as a format name, two modifiers that contradict each other,
// ANALYZE on a statement MySQL will not analyze, or a modifier in front of the
// table form. Every case was measured against MySQL 8.0.46 and
// TestExplainFormsMatchMySQL holds this file to the engine, so an image
// upgrade that changes a verdict fails there instead of passing silently.
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/explain.html; describe_stmt,
// opt_explain_format, opt_describe_column and explainable_stmt in MySQL 8.0.46
// sql/sql_yacc.yy.

// explainFormCase is a statement and whether the engine's parser accepts it.
type explainFormCase struct {
	sql    string
	accept bool
}

var explainFormCases = []explainFormCase{
	// The format name is ident_or_text: an identifier, a non-reserved
	// keyword, a backquoted identifier or a string literal in either quote.
	// MySQL checks the name itself after parsing (1791), so a name it does
	// not have still parses.
	{"EXPLAIN FORMAT = TRADITIONAL SELECT 1", true},
	{"EXPLAIN FORMAT = JSON SELECT 1", true},
	{"EXPLAIN FORMAT = TREE SELECT 1", true},
	{"EXPLAIN FORMAT = tree SELECT 1", true},
	{"EXPLAIN FORMAT = 'traditional' SELECT 1", true},
	{"EXPLAIN FORMAT = 'TRADITIONAL' SELECT 1", true},
	{`EXPLAIN FORMAT = "traditional" SELECT 1`, true},
	{"EXPLAIN FORMAT = `traditional` SELECT 1", true},
	{"EXPLAIN FORMAT = 'json' SELECT 1", true},
	{"EXPLAIN FORMAT = 'row' SELECT 1", false},
	{"EXPLAIN FORMAT = 'nosuchformat' SELECT 1", false},
	{"EXPLAIN FORMAT = nosuchformat SELECT 1", false},
	{"EXPLAIN FORMAT = `row` SELECT 1", false},
	{"EXPLAIN FORMAT = '' SELECT 1", false},
	{"EXPLAIN ANALYZE FORMAT = TREE SELECT 1", true},
	{"EXPLAIN ANALYZE FORMAT = 'tree' SELECT 1", true},
	{"EXPLAIN ANALYZE FORMAT = 'traditional' SELECT 1", true},
	// ...and nothing else. A reserved word is no name, and neither is an
	// integer: this is where MySQL and TiDB part, since TiDB reads a closed
	// list of bare names and any string at all.
	{"EXPLAIN FORMAT = binary SELECT 1", false},
	{"EXPLAIN FORMAT = select SELECT 1", false},
	{"EXPLAIN FORMAT = 1 SELECT 1", false},
	// The value is required, the '=' is required, and one FORMAT is the most
	// a statement takes.
	{"EXPLAIN FORMAT = SELECT 1", false},
	{"EXPLAIN FORMAT SELECT 1", false},
	{"EXPLAIN FORMAT 'row' SELECT 1", false},
	{"EXPLAIN FORMAT = 'traditional' FORMAT = 'json' SELECT 1", false},
	{"EXPLAIN FORMAT = 'traditional' 'traditional' SELECT 1", false},
	{"EXPLAIN FORMAT = TRADITIONAL FORMAT = JSON SELECT 1", false},
	{"EXPLAIN FORMAT = 'traditional' SELECT 1 FORMAT = 'json'", false},

	// FOR CONNECTION, with a format and with ANALYZE. CONNECTION and the id
	// are both required and the id is an unsigned integer.
	{"EXPLAIN FOR CONNECTION 1", true},
	{"EXPLAIN FOR connection 42", true},
	{"EXPLAIN FOR CONNECTION 0", true},
	{"EXPLAIN FOR CONNECTION 18446744073709551615", true},
	{"EXPLAIN FORMAT = 'traditional' FOR CONNECTION 1", true},
	{"EXPLAIN FORMAT = TREE FOR CONNECTION 1", true},
	// MySQL parses EXPLAIN ANALYZE FOR CONNECTION and answers 1235, where
	// TiDB rejects it outright.
	{"EXPLAIN ANALYZE FOR CONNECTION 1", true},
	{"EXPLAIN ANALYZE FORMAT = TREE FOR CONNECTION 1", true},
	{"EXPLAIN FOR CONNECTION", false},
	{"EXPLAIN FOR 1", false},
	{"EXPLAIN FOR CONNECTION 'x'", false},
	{"EXPLAIN FOR CONNECTION 1.5", false},
	{"EXPLAIN FOR CONNECTION @a", false},
	{"EXPLAIN FOR CONNECTION -1", false},
	{"EXPLAIN FOR CONNECTION 1+1", false},
	{"EXPLAIN FOR CONNECTION (1)", false},
	{"EXPLAIN FOR CONNECTION a", false},
	{"EXPLAIN FOR CONNECTION t.a", false},
	{"EXPLAIN FORMAT = TRADITIONAL FOR CONNECTION", false},
	// Neither FOR DATABASE nor an INTO is MySQL 8.0 grammar.
	{"EXPLAIN FOR DATABASE d SELECT 1", false},
	{"EXPLAIN FOR SCHEMA d SELECT 1", false},
	{"EXPLAIN FORMAT = JSON INTO @v SELECT 1", false},
	{"EXPLAIN INTO @v SELECT 1", false},

	// EXTENDED and PARTITIONS went with MySQL 5.7: 8.0 folded both into the
	// default output and answers 1064 for each, as TiDB does. MariaDB and
	// OceanBase are the engines in this family that still take them.
	{"EXPLAIN EXTENDED SELECT 1", false},
	{"EXPLAIN PARTITIONS SELECT 1", false},
	{"EXPLAIN ANALYZE EXTENDED SELECT 1", false},
	{"EXPLAIN EXTENDED FORMAT = JSON SELECT 1", false},
	{"EXPLAIN FORMAT = JSON EXTENDED SELECT 1", false},
	{"DESC EXTENDED SELECT 1", false},

	// The explainable statement.
	{"EXPLAIN SELECT 1", true},
	{"EXPLAIN ANALYZE SELECT 1", true},
	{"EXPLAIN INSERT INTO t1 VALUES (1)", true},
	{"EXPLAIN REPLACE INTO t1 VALUES (1)", true},
	{"EXPLAIN UPDATE t1 SET a = 1", true},
	{"EXPLAIN DELETE FROM t1", true},
	{"EXPLAIN TABLE t1", true},
	{"EXPLAIN VALUES ROW(1)", true},
	{"EXPLAIN FORMAT = 'traditional' UPDATE t1 SET a = 1", true},
	{"EXPLAIN FORMAT = 'traditional' DELETE FROM t1", true},
	{"EXPLAIN FORMAT = 'traditional' INSERT INTO t1 VALUES (1)", true},
	{"EXPLAIN FORMAT = 'traditional' REPLACE INTO t1 VALUES (1)", true},
	{"EXPLAIN FORMAT = 'traditional' TABLE t1", true},
	{"EXPLAIN FORMAT = 'traditional' VALUES ROW(1)", true},
	{"EXPLAIN ANALYZE TABLE t1", true},
	{"EXPLAIN ANALYZE VALUES ROW(1)", true},
	{"EXPLAIN ANALYZE UPDATE t1 SET a = 1", true},
	// A parenthesized query and a WITH clause in front of one.
	{"EXPLAIN (SELECT 1)", true},
	{"EXPLAIN WITH c AS (SELECT 1) SELECT * FROM c", true},
	{"EXPLAIN FORMAT = 'traditional' (SELECT 1)", true},
	{"EXPLAIN FORMAT = 'traditional' WITH c AS (SELECT 1) SELECT * FROM c", true},
	{"EXPLAIN ANALYZE (SELECT 1)", true},
	{"EXPLAIN ANALYZE WITH c AS (SELECT 1) SELECT * FROM c", true},
	// ...and no DDL at all, which is where TiDB parts from MySQL: it
	// explains ALTER TABLE and MySQL explains nothing.
	{"EXPLAIN ALTER TABLE t1 ADD INDEX (a)", false},
	{"EXPLAIN ALTER TABLE t1 ADD a varchar(255)", false},
	{"EXPLAIN ANALYZE ALTER TABLE t1 ADD a int", false},
	{"EXPLAIN FORMAT = 'traditional' ALTER TABLE t1 ADD INDEX (a)", false},
	{"EXPLAIN CREATE TABLE t2 (a int)", false},
	{"EXPLAIN CREATE INDEX i ON t1 (a)", false},
	{"EXPLAIN DROP TABLE t1", false},
	{"EXPLAIN TRUNCATE TABLE t1", false},
	{"EXPLAIN RENAME TABLE t1 TO t2", false},
	{"EXPLAIN LOAD DATA INFILE 'x' INTO TABLE t1", false},
	{"EXPLAIN SET @a = 1", false},
	{"EXPLAIN ANALYZE SET @a = 1", false},
	{"EXPLAIN FORMAT = 'traditional' SET @a = 1", false},
	{"EXPLAIN FORMAT = 'traditional' SHOW TABLES", false},

	// The table form: a table name, then a column name or a wild pattern.
	// The string is a pattern here, which TiDB refuses, and the qualified
	// name TiDB takes is a 1064 here.
	{"EXPLAIN t1", true},
	{"EXPLAIN t1 a", true},
	{"EXPLAIN t1 `a`", true},
	{"EXPLAIN t1 'abc'", true},
	{"EXPLAIN t1 '%a%'", true},
	{"EXPLAIN d1.t1 a", true},
	{"DESC t1 a", true},
	{"DESC t1 'abc'", true},
	{"DESCRIBE t1 'abc'", true},
	{"EXPLAIN BEGIN", true},
	// EXTENDED and PARTITIONS are not reserved, so each is a table name
	// here, and that is all they are: the modifier they used to be is gone.
	{"EXPLAIN EXTENDED t1", true},
	{"EXPLAIN PARTITIONS t1", true},
	{"EXPLAIN EXTENDED", true},
	{"EXPLAIN t1 a.b", false},
	{"EXPLAIN t1 123", false},
	{"EXPLAIN t1 select", false},
	{"EXPLAIN t1 a b", false},
	// A modifier rules the table form out, so a table name after one is a
	// 1064 and so is a column name after that.
	{"EXPLAIN FORMAT = 'traditional' t1", false},
	{"EXPLAIN FORMAT = 'traditional' t1 a", false},
	{"EXPLAIN FORMAT = TRADITIONAL t1", false},
	{"EXPLAIN ANALYZE t1", false},
	{"EXPLAIN ANALYZE t1 a", false},
	// A statement in a string is TiDB grammar, and not MySQL's at all.
	{"EXPLAIN 'sqldigest'", false},
	{"EXPLAIN ANALYZE 'sqldigest'", false},
}

// explainFormEmptyFormatTableForm are statements where an empty format name
// stands in front of the table form. A modifier rules that form out, so each
// is a syntax error, and the guard that says so has to track whether FORMAT
// was read rather than whether the name it read is empty.
//
// They are held to the parser alone and not to the server: MySQL 8.0.46 raises
// 1791 for the empty name, which it does not have, while it is still
// parsing, before it reads the shape at all, so it gives no verdict on
// these. The non-empty
// FORMAT = 'tree' t, in explainFormCases, is the one the server settles, and
// TiDB settles the empty one too (tidb/parser/explain_forms_test.go).
var explainFormEmptyFormatTableForm = []string{
	"EXPLAIN FORMAT = '' t1",
	"EXPLAIN FORMAT = '' t1 a",
	`EXPLAIN FORMAT = "" t1`,
}

func TestExplainFormEmptyFormatRulesTheTableFormOut(t *testing.T) {
	for _, sql := range explainFormEmptyFormatTableForm {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err == nil {
				t.Fatalf("Parse(%q) was accepted; a FORMAT, empty name and all, rules the table form out", sql)
			}
		})
	}
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
				t.Fatalf("Parse(%q) was accepted, but MySQL 8.0 rejects it", c.sql)
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
			{"EXPLAIN FORMAT = 'TREE' SELECT 1", "TREE"},
			{"EXPLAIN FORMAT = 'tree' SELECT 1", "tree"},
			{"EXPLAIN FORMAT = TREE SELECT 1", "TREE"},
			{`EXPLAIN FORMAT = "tree" SELECT 1`, "tree"},
			{"EXPLAIN FORMAT = `tree` SELECT 1", "tree"},
		} {
			stmt := parseOneExplain(t, tc.sql)
			if stmt.Format != tc.want {
				t.Errorf("%s: Format = %q, want %q", tc.sql, stmt.Format, tc.want)
			}
			if stmt.Analyze {
				t.Errorf("%s: Analyze = true, want false", tc.sql)
			}
			if _, ok := stmt.Stmt.(*nodes.SelectStmt); !ok {
				t.Errorf("%s: Stmt = %T, want *ast.SelectStmt", tc.sql, stmt.Stmt)
			}
		}
	})

	t.Run("analyze with a format", func(t *testing.T) {
		stmt := parseOneExplain(t, "EXPLAIN ANALYZE FORMAT = 'tree' SELECT 1")
		if !stmt.Analyze {
			t.Error("Analyze = false, want true")
		}
		if stmt.Format != "tree" {
			t.Errorf("Format = %q, want %q", stmt.Format, "tree")
		}
		// MySQL 8.0 has neither modifier, so no statement can set either.
		if stmt.Extended || stmt.Partitions {
			t.Errorf("Extended = %v, Partitions = %v, want both false", stmt.Extended, stmt.Partitions)
		}
	})

	t.Run("for connection", func(t *testing.T) {
		stmt := parseOneExplain(t, "EXPLAIN FORMAT = TREE FOR CONNECTION 42")
		if stmt.Format != "TREE" {
			t.Errorf("Format = %q, want TREE", stmt.Format)
		}
		if stmt.ForConnection != 42 {
			t.Errorf("ForConnection = %d, want 42", stmt.ForConnection)
		}
		// The connection form explains no statement of its own, which is how
		// a consumer tells it from a connection id that happens to be 0.
		if stmt.Stmt != nil {
			t.Errorf("Stmt = %T, want nil for the FOR CONNECTION form", stmt.Stmt)
		}
		analyzed := parseOneExplain(t, "EXPLAIN ANALYZE FOR CONNECTION 0")
		if !analyzed.Analyze || analyzed.ForConnection != 0 || analyzed.Stmt != nil {
			t.Errorf("ANALYZE FOR CONNECTION 0: Analyze = %v, ForConnection = %d, Stmt = %T",
				analyzed.Analyze, analyzed.ForConnection, analyzed.Stmt)
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
			"TABLE t1",
			"VALUES ROW(1)",
			"(SELECT 1)",
			"WITH c AS (SELECT 1) SELECT * FROM c",
		} {
			alone, err := Parse(inner)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", inner, err)
			}
			for _, prefix := range []string{"EXPLAIN ", "EXPLAIN ANALYZE ", "EXPLAIN FORMAT = 'tree' "} {
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

// TestExplainFormsMatchMySQL holds the table above to the engine itself (container test).
func TestExplainFormsMatchMySQL(t *testing.T) {
	o := startParserOracle(t)
	for _, c := range explainFormCases {
		t.Run(c.sql, func(t *testing.T) {
			if got := o.prepareParses(t, c.sql); got != c.accept {
				t.Errorf("MySQL parses %q: %v, this file records accept=%v", c.sql, got, c.accept)
			}
		})
	}
}
