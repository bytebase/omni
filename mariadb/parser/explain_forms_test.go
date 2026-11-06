package parser

import (
	"errors"
	"regexp"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// EXPLAIN's modifier list and its three slots -- the format name, the
// connection id and the explained statement -- as MariaDB 11.8 reads them.
//
// The list is not MySQL 8.0's, and this parser was forked from a parser that
// read MySQL's. EXTENDED and PARTITIONS are still here, where MySQL 8.0.46 and
// TiDB each answer 1064 for both, and they are one choice rather than a list.
// There is no EXPLAIN ANALYZE at all: MariaDB spells that ANALYZE [FORMAT =
// JSON] stmt, a statement of its own that this parser does not read yet, so
// every EXPLAIN ANALYZE below is a rejection rather than a form to port from
// MySQL.
//
// The rejections are the point of this file. A parser that accepts what the
// engine will not lets a review pass a statement that then fails on the
// server, and nothing else in this package covers a FORMAT with no value, a
// reserved word as a format name, two modifiers that contradict each other, or
// a modifier in front of the table form. Every case was measured against
// mariadb:11.8.8 and TestExplainFormsMatchMariaDB holds this file to the
// engine, so an image upgrade that changes a verdict fails there instead of
// passing silently.
//
// Ref: https://mariadb.com/kb/en/explain/; describe_command and
// opt_describe_column in MariaDB 11.8 sql/sql_yacc.yy.

// explainFormCase is a statement and whether the engine's parser accepts it.
type explainFormCase struct {
	sql    string
	accept bool
}

var explainFormCases = []explainFormCase{
	// The format name is ident_or_text: an identifier, a non-reserved
	// keyword, a backquoted identifier or a string literal in either quote.
	// MariaDB checks the name itself after parsing (1791), so a name it does
	// not have still parses.
	{"EXPLAIN FORMAT = JSON SELECT 1", true},
	{"EXPLAIN FORMAT = TRADITIONAL SELECT 1", true},
	{"EXPLAIN FORMAT = TREE SELECT 1", false},
	{"EXPLAIN FORMAT = 'json' SELECT 1", true},
	{"EXPLAIN FORMAT = 'JSON' SELECT 1", true},
	{`EXPLAIN FORMAT = "json" SELECT 1`, true},
	{"EXPLAIN FORMAT = `json` SELECT 1", true},
	{"EXPLAIN FORMAT = 'traditional' SELECT 1", true},
	{"EXPLAIN FORMAT = 'nosuchformat' SELECT 1", false},
	{"EXPLAIN FORMAT = nosuchformat SELECT 1", false},
	{"EXPLAIN FORMAT = '' SELECT 1", false},
	// ...and nothing else. A reserved word is no name, and neither is an
	// integer.
	{"EXPLAIN FORMAT = binary SELECT 1", false},
	{"EXPLAIN FORMAT = select SELECT 1", false},
	{"EXPLAIN FORMAT = 1 SELECT 1", false},
	// The value is required, the '=' is required, and one FORMAT is the most
	// a statement takes.
	{"EXPLAIN FORMAT = SELECT 1", false},
	{"EXPLAIN FORMAT SELECT 1", false},
	{"EXPLAIN FORMAT 'json' SELECT 1", false},
	{"EXPLAIN FORMAT = 'json' FORMAT = 'traditional' SELECT 1", false},
	{"EXPLAIN FORMAT = 'traditional' 'traditional' SELECT 1", false},
	{"EXPLAIN FORMAT = JSON FORMAT = TRADITIONAL SELECT 1", false},

	// EXTENDED and PARTITIONS, which MySQL 8.0 and TiDB both dropped. They
	// are one choice: neither stands with the other, neither stands with a
	// FORMAT, and neither reaches FOR CONNECTION.
	{"EXPLAIN EXTENDED SELECT 1", true},
	{"EXPLAIN PARTITIONS SELECT 1", true},
	{"EXPLAIN EXTENDED SELECT * FROM t1", true},
	{"EXPLAIN PARTITIONS SELECT * FROM t1", true},
	{"EXPLAIN EXTENDED FORMAT = JSON SELECT 1", false},
	{"EXPLAIN FORMAT = JSON EXTENDED SELECT 1", false},
	{"EXPLAIN EXTENDED FOR CONNECTION 1", false},
	{"EXPLAIN PARTITIONS FOR CONNECTION 1", false},
	// EXPLAIN EXTENDED PARTITIONS SELECT 1 belongs here too, and is not:
	// EXTENDED and PARTITIONS read as a table name and a column name, and
	// this parser then takes the SELECT that follows as a second statement
	// with no separator in front of it, which every engine in this family
	// does and no server does. That is a defect of its own -- SELECT 1
	// SELECT 2 parses here too -- and not one of EXPLAIN.

	// There is no EXPLAIN ANALYZE: MariaDB's own spelling is a statement of
	// its own, ANALYZE [FORMAT = JSON] stmt, which MySQL and TiDB both lack
	// and this parser does not read yet.
	{"EXPLAIN ANALYZE SELECT 1", false},
	{"EXPLAIN ANALYZE FORMAT = JSON SELECT 1", false},
	{"EXPLAIN ANALYZE FORMAT = 'json' SELECT 1", false},
	{"EXPLAIN ANALYZE FOR CONNECTION 1", false},
	{"EXPLAIN ANALYZE t1", false},
	{"EXPLAIN ANALYZE UPDATE t1 SET a = 1", false},
	{"DESC ANALYZE SELECT 1", false},

	// FOR CONNECTION, with and without a format. CONNECTION and the id are
	// both required.
	{"EXPLAIN FOR CONNECTION 1", true},
	{"EXPLAIN FOR connection 42", true},
	{"EXPLAIN FOR CONNECTION 0", true},
	{"EXPLAIN FOR CONNECTION 18446744073709551615", true},
	{"EXPLAIN FORMAT = 'traditional' FOR CONNECTION 1", true},
	{"EXPLAIN FORMAT = JSON FOR CONNECTION 1", true},
	{"EXPLAIN FOR CONNECTION", false},
	{"EXPLAIN FOR 1", false},
	{"EXPLAIN FORMAT = JSON FOR CONNECTION", false},

	// The explainable statement.
	{"EXPLAIN SELECT 1", true},
	{"EXPLAIN INSERT INTO t1 VALUES (1)", true},
	{"EXPLAIN REPLACE INTO t1 VALUES (1)", true},
	{"EXPLAIN UPDATE t1 SET a = 1", true},
	{"EXPLAIN DELETE FROM t1", true},
	{"EXPLAIN VALUES (1)", true},
	{"EXPLAIN FORMAT = JSON UPDATE t1 SET a = 1", true},
	{"EXPLAIN FORMAT = JSON DELETE FROM t1", true},
	{"EXPLAIN FORMAT = JSON INSERT INTO t1 VALUES (1)", true},
	{"EXPLAIN FORMAT = JSON REPLACE INTO t1 VALUES (1)", true},
	{"EXPLAIN FORMAT = JSON VALUES (1)", true},
	{"EXPLAIN EXTENDED UPDATE t1 SET a = 1", true},
	{"EXPLAIN (SELECT 1)", true},
	{"EXPLAIN WITH c AS (SELECT 1) SELECT * FROM c", true},
	{"EXPLAIN FORMAT = JSON (SELECT 1)", true},
	{"EXPLAIN FORMAT = JSON WITH c AS (SELECT 1) SELECT * FROM c", true},
	// TABLE is no statement on MariaDB, and neither is MySQL's VALUES ROW().
	{"EXPLAIN TABLE t1", false},
	{"EXPLAIN VALUES ROW(1)", false},
	{"EXPLAIN FORMAT = JSON TABLE t1", false},
	// ...and no DDL at all, which is where TiDB parts: it explains ALTER
	// TABLE and MariaDB explains nothing.
	{"EXPLAIN ALTER TABLE t1 ADD INDEX (a)", false},
	{"EXPLAIN ALTER TABLE t1 ADD a varchar(255)", false},
	{"EXPLAIN FORMAT = JSON ALTER TABLE t1 ADD INDEX (a)", false},
	{"EXPLAIN CREATE TABLE t2 (a int)", false},
	{"EXPLAIN CREATE INDEX i ON t1 (a)", false},
	{"EXPLAIN DROP TABLE t1", false},
	{"EXPLAIN TRUNCATE TABLE t1", false},
	{"EXPLAIN RENAME TABLE t1 TO t2", false},
	{"EXPLAIN SET @a = 1", false},
	{"EXPLAIN FORMAT = JSON SET @a = 1", false},
	{"EXPLAIN FORMAT = JSON SHOW TABLES", false},

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
	{"EXPLAIN t1 a.b", false},
	{"EXPLAIN t1 123", false},
	{"EXPLAIN t1 select", false},
	{"EXPLAIN t1 a b", false},
	// A modifier rules the table form out, so a table name after one is a
	// 1064 and so is a column name after that.
	{"EXPLAIN FORMAT = JSON t1", false},
	{"EXPLAIN FORMAT = 'json' t1", false},
	{"EXPLAIN FORMAT = JSON t1 a", false},
	// EXTENDED and PARTITIONS are not reserved, so each is a table name as
	// well as a modifier: EXPLAIN EXTENDED t describes a table called
	// `extended`, and EXPLAIN EXTENDED on its own describes that table. What
	// follows the word decides which it was.
	{"EXPLAIN EXTENDED t1", true},
	{"EXPLAIN PARTITIONS t1", true},
	{"EXPLAIN EXTENDED", true},
	{"EXPLAIN PARTITIONS `a`", true},
	{"EXPLAIN FORMAT = JSON t1", false},
	// A statement in a string is TiDB grammar, and not MariaDB's at all.
	{"EXPLAIN 'sqldigest'", false},
}

// explainFormConnectionIDDivergences are the FOR CONNECTION ids MariaDB 11.8.8
// parses and this parser turns away. MariaDB reads a whole expression there --
// a signed number, a string literal, a user variable, a sum, a parenthesized
// value and a column reference all reach the server, where MySQL 8.0.46 and
// TiDB v8.5.5 answer 1064 for every one of them -- while ExplainStmt holds a
// connection id, which is the documented form. Reading the rest would mean an
// expression field beside that id, for input no layer above the parser could
// use, so the id stays an unsigned integer here.
//
// The divergence is pinned in both directions by
// TestExplainFormConnectionIDDivergences: MariaDB must still accept each one
// and this parser must still reject it, so the day either side changes, this
// list is what fails.
var explainFormConnectionIDDivergences = []string{
	"EXPLAIN FOR CONNECTION 'x'",
	"EXPLAIN FOR CONNECTION 1.5",
	"EXPLAIN FOR CONNECTION @a",
	"EXPLAIN FOR CONNECTION -1",
	"EXPLAIN FOR CONNECTION 1+1",
	"EXPLAIN FOR CONNECTION (1)",
	"EXPLAIN FOR CONNECTION a",
	"EXPLAIN FOR CONNECTION t.a",
	"EXPLAIN FORMAT = JSON FOR CONNECTION 'x'",
}

// explainFormEmptyFormatTableForm are statements where an empty format name
// stands in front of the table form. A modifier rules that form out, so each
// is a syntax error, and the guard that says so has to track whether FORMAT
// was read rather than whether the name it read is empty.
//
// They are held to the parser alone and not to the server: MariaDB 11.8.8 raises
// 1791 for the empty name, which it does not have, while it is still
// parsing, before it reads the shape at all, so it gives no verdict on
// these. The non-empty
// FORMAT = JSON t, in explainFormCases, is the one the server settles, and
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
				t.Fatalf("Parse(%q) was accepted, but MariaDB 11.8 rejects it", c.sql)
			}
		})
	}
}

func TestExplainFormConnectionIDDivergencesRejected(t *testing.T) {
	for _, sql := range explainFormConnectionIDDivergences {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err == nil {
				t.Fatalf("Parse(%q) was accepted; this parser reads an unsigned integer connection id", sql)
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
			{"EXPLAIN FORMAT = 'JSON' SELECT 1", "JSON"},
			{"EXPLAIN FORMAT = 'json' SELECT 1", "json"},
			{"EXPLAIN FORMAT = JSON SELECT 1", "JSON"},
			{`EXPLAIN FORMAT = "json" SELECT 1`, "json"},
			{"EXPLAIN FORMAT = `json` SELECT 1", "json"},
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
		// MariaDB has no EXPLAIN ANALYZE, so no statement can set the flag.
		for _, c := range explainFormCases {
			if !c.accept {
				continue
			}
			if stmt := parseOneExplain(t, c.sql); stmt.Analyze {
				t.Errorf("%q set Analyze, which MariaDB's EXPLAIN has no syntax for", c.sql)
			}
		}
	})

	t.Run("for connection", func(t *testing.T) {
		stmt := parseOneExplain(t, "EXPLAIN FORMAT = JSON FOR CONNECTION 42")
		if stmt.Format != "JSON" {
			t.Errorf("Format = %q, want JSON", stmt.Format)
		}
		if stmt.ForConnection != 42 {
			t.Errorf("ForConnection = %d, want 42", stmt.ForConnection)
		}
		// The connection form explains no statement of its own, which is how
		// a consumer tells it from a connection id that happens to be 0.
		if stmt.Stmt != nil {
			t.Errorf("Stmt = %T, want nil for the FOR CONNECTION form", stmt.Stmt)
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
			"VALUES (1)",
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

// explainFormVerdict asks MariaDB whether it parses sql. The verdict comes
// from PREPARE, which parses the statement without running it: 1064 is the
// parser's rejection, and any other server error means the text parsed and
// failed later.
func explainFormVerdict(t *testing.T, o *parserOracle, sql string) (accepts bool, code uint16) {
	t.Helper()
	_, err := o.db.ExecContext(o.ctx, "PREPARE parse_probe FROM '"+strings.ReplaceAll(sql, "'", "''")+"'")
	if err != nil {
		var myErr *gomysql.MySQLError
		if !errors.As(err, &myErr) {
			t.Fatalf("PREPARE %q: the oracle gave no verdict: %v", sql, err)
		}
		code = myErr.Number
	}
	return code != 1064, code
}

// TestExplainFormsMatchMariaDB holds the table above to the engine itself (container test).
func TestExplainFormsMatchMariaDB(t *testing.T) {
	o := startMariaDB(t)
	for _, c := range explainFormCases {
		t.Run(c.sql, func(t *testing.T) {
			if got := o.prepareParses(t, c.sql); got != c.accept {
				t.Errorf("MariaDB parses %q: %v, this file records accept=%v", c.sql, got, c.accept)
			}
		})
	}
}
