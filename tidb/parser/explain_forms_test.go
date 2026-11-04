package parser

import (
	"math"
	"regexp"
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// EXPLAIN's modifier list and its three slots -- the format name, the
// connection id and the explained statement -- as TiDB v8.5.5 reads them.
//
// The corpus in testdata pins the forms pingcap's own parser test writes. It
// writes almost none of the rejections below, and a parser that accepts what
// the engine will not lets a review pass a statement that then fails on the
// server, so this file carries them: a FORMAT with no value, a format name
// TiDB does not have, a backquoted one, two modifiers that contradict each
// other, ANALYZE on a statement TiDB will not analyze, a modifier in front of
// the table form, and the EXPLAIN of every DDL but ALTER TABLE.
//
// Every case was measured against TiDB v8.5.5 and the expectation is what that
// server answered; TestExplainFormsMatchTiDB holds this file to it, so an image
// upgrade that changes a verdict fails there instead of passing silently.
//
// Ref: ExplainStmt, ExplainFormatType and ExplainableStmt in TiDB v8.5.5
// parser/parser.y.

// explainFormCase is a statement and whether the engine's parser accepts it.
type explainFormCase struct {
	sql    string
	accept bool
}

var explainFormCases = []explainFormCase{
	// The eight bare format names, which are keyword tokens in TiDB's own
	// grammar, and the same eight after ANALYZE.
	{"EXPLAIN FORMAT = TRADITIONAL SELECT 1", true},
	{"EXPLAIN FORMAT = JSON SELECT 1", true},
	{"EXPLAIN FORMAT = ROW SELECT 1", true},
	{"EXPLAIN FORMAT = DOT SELECT 1", true},
	{"EXPLAIN FORMAT = BRIEF SELECT 1", true},
	{"EXPLAIN FORMAT = VERBOSE SELECT 1", true},
	{"EXPLAIN FORMAT = TRUE_CARD_COST SELECT 1", true},
	{"EXPLAIN FORMAT = TIDB_JSON SELECT 1", true},
	{"EXPLAIN FORMAT = tidb_json SELECT 1", true},
	{"EXPLAIN ANALYZE FORMAT = VERBOSE SELECT 1", true},
	{"EXPLAIN ANALYZE FORMAT = TIDB_JSON SELECT 1", true},
	// A quoted name is an arm of its own and carries any text at all: TiDB
	// checks it as a name when the statement runs (1791), not when it parses.
	{"EXPLAIN FORMAT = 'row' SELECT 1", true},
	{"EXPLAIN FORMAT = 'ROW' SELECT 1", true},
	{`EXPLAIN FORMAT = "row" SELECT 1`, true},
	{"EXPLAIN FORMAT = 'dot' SELECT 1", true},
	{"EXPLAIN FORMAT = 'brief' SELECT 1", true},
	{"EXPLAIN FORMAT = 'verbose' SELECT 1", true},
	{"EXPLAIN FORMAT = 'hint' SELECT 1", true},
	{"EXPLAIN FORMAT = 'nosuchformat' SELECT 1", true},
	{"EXPLAIN FORMAT = '' SELECT 1", true},
	{"EXPLAIN ANALYZE FORMAT = 'verbose' SELECT 1", true},
	{"EXPLAIN ANALYZE FORMAT = 'binary' SELECT 1", true},
	{"EXPLAIN ANALYZE format='TRUE_CARD_COST' SELECT 1", true},
	// ...and the bare arm is the closed list and nothing else. A name outside
	// it, a backquoted name, a reserved word and an integer are each a 1064,
	// which is where TiDB and MySQL part: MySQL reads any ident_or_text here.
	{"EXPLAIN FORMAT = tree SELECT 1", false},
	{"EXPLAIN FORMAT = hint SELECT 1", false},
	{"EXPLAIN FORMAT = binary SELECT 1", false},
	{"EXPLAIN FORMAT = cost_trace SELECT 1", false},
	{"EXPLAIN FORMAT = plan_cache SELECT 1", false},
	{"EXPLAIN FORMAT = nosuchformat SELECT 1", false},
	{"EXPLAIN FORMAT = `row` SELECT 1", false},
	{"EXPLAIN FORMAT = select SELECT 1", false},
	{"EXPLAIN FORMAT = 1 SELECT 1", false},
	{"EXPLAIN ANALYZE FORMAT = tree SELECT 1", false},
	{"EXPLAIN ANALYZE FORMAT = nosuchformat SELECT 1", false},
	// The value is required, the '=' is required, and one FORMAT is the most
	// a statement takes.
	{"EXPLAIN FORMAT = SELECT 1", false},
	{"EXPLAIN FORMAT SELECT 1", false},
	{"EXPLAIN FORMAT 'row' SELECT 1", false},
	{"EXPLAIN FORMAT = 'row' FORMAT = 'json' SELECT 1", false},
	{"EXPLAIN FORMAT = 'row' 'row' SELECT 1", false},
	{"EXPLAIN FORMAT = ROW ROW SELECT 1", false},

	// FOR CONNECTION, with and without a format. CONNECTION and the id are
	// both required and the id is an unsigned integer.
	{"EXPLAIN FOR CONNECTION 1", true},
	{"EXPLAIN FOR connection 42", true},
	{"EXPLAIN FOR CONNECTION 0", true},
	{"EXPLAIN FOR CONNECTION 18446744073709551615", true},
	{"EXPLAIN FORMAT = 'dot' FOR CONNECTION 1", true},
	{"EXPLAIN FORMAT = DOT FOR CONNECTION 1", true},
	{"EXPLAIN FORMAT = 'row' FOR connection 1", true},
	{"EXPLAIN FOR CONNECTION", false},
	{"EXPLAIN FOR 1", false},
	{"EXPLAIN FOR CONNECTION 'x'", false},
	{"EXPLAIN FOR CONNECTION 1.5", false},
	{"EXPLAIN FOR CONNECTION @a", false},
	{"EXPLAIN FOR CONNECTION -1", false},
	{"EXPLAIN FOR CONNECTION 1+1", false},
	{"EXPLAIN FOR CONNECTION (1)", false},
	{"EXPLAIN FOR CONNECTION a", false},
	{"EXPLAIN FORMAT = 'row' FOR CONNECTION", false},
	// ANALYZE does not reach FOR CONNECTION, which is where TiDB and MySQL
	// part again: MySQL parses it and answers 1235.
	{"EXPLAIN ANALYZE FOR CONNECTION 1", false},
	{"EXPLAIN ANALYZE FORMAT = 'row' FOR CONNECTION 1", false},

	// EXTENDED and PARTITIONS are MySQL 5.7 modifiers TiDB never had, and
	// MySQL 8.0.46 answers 1064 for both as well. MariaDB is the engine in
	// this family that still takes them.
	{"EXPLAIN EXTENDED SELECT 1", false},
	{"EXPLAIN PARTITIONS SELECT 1", false},
	{"EXPLAIN ANALYZE EXTENDED SELECT 1", false},
	{"EXPLAIN EXTENDED FORMAT = ROW SELECT 1", false},
	// EXPLAIN EXTENDED PARTITIONS SELECT 1 belongs here too, and is not:
	// EXTENDED and PARTITIONS read as a table name and a column name, and
	// this parser then takes the SELECT that follows as a second statement
	// with no separator in front of it, which every engine in this family
	// does and no server does. That is a defect of its own -- SELECT 1
	// SELECT 2 parses here too -- and not one of EXPLAIN.

	// The explainable statement. ALTER TABLE is the one DDL TiDB explains,
	// with and without ANALYZE and a format.
	{"EXPLAIN SELECT 1", true},
	{"EXPLAIN ANALYZE SELECT 1", true},
	{"EXPLAIN INSERT INTO t1 VALUES (1)", true},
	{"EXPLAIN REPLACE INTO t1 VALUES (1)", true},
	{"EXPLAIN UPDATE t1 SET a = 1", true},
	{"EXPLAIN DELETE FROM t1", true},
	{"EXPLAIN TABLE t1", true},
	{"EXPLAIN VALUES ROW(1)", true},
	{"EXPLAIN FORMAT = 'row' UPDATE t1 SET a = 1", true},
	{"EXPLAIN FORMAT = 'row' DELETE FROM t1", true},
	{"EXPLAIN FORMAT = 'row' INSERT INTO t1 VALUES (1)", true},
	{"EXPLAIN FORMAT = 'row' REPLACE INTO t1 VALUES (1)", true},
	{"EXPLAIN FORMAT = 'row' TABLE t1", true},
	{"EXPLAIN FORMAT = 'row' VALUES ROW(1)", true},
	{"EXPLAIN ANALYZE TABLE t1", true},
	{"EXPLAIN ANALYZE UPDATE t1 SET a = 1", true},
	{"EXPLAIN ALTER TABLE t1 ADD INDEX (a)", true},
	{"EXPLAIN ALTER TABLE t1 ADD a varchar(255)", true},
	{"EXPLAIN ALTER TABLE t1 MODIFY a bigint", true},
	{"EXPLAIN ALTER TABLE t1 DROP INDEX a", true},
	{"EXPLAIN ALTER TABLE t1 ADD INDEX (a), ADD INDEX (b)", true},
	{"EXPLAIN FORMAT = 'brief' ALTER TABLE t1 ADD INDEX (a)", true},
	{"EXPLAIN FORMAT = DOT ALTER TABLE t1 ADD INDEX (a)", true},
	{"EXPLAIN ANALYZE ALTER TABLE t1 ADD INDEX (a)", true},
	{"EXPLAIN ANALYZE FORMAT = 'brief' ALTER TABLE t1 ADD INDEX (a)", true},
	// A parenthesized query and a WITH clause in front of one, which are
	// SubSelect and SelectStmtWithClause in ExplainableStmt.
	{"EXPLAIN (SELECT 1)", true},
	{"EXPLAIN WITH c AS (SELECT 1) SELECT * FROM c", true},
	{"EXPLAIN FORMAT = 'row' (SELECT 1)", true},
	{"EXPLAIN FORMAT = 'row' WITH c AS (SELECT 1) SELECT * FROM c", true},
	{"explain format = 'row' select c1 from t1 union (select c2 from t2) limit 1, 1", true},
	// ...and no other DDL, no SET, no SHOW.
	{"EXPLAIN ALTER DATABASE d CHARSET = utf8", false},
	{"EXPLAIN CREATE TABLE t2 (a int)", false},
	{"EXPLAIN CREATE INDEX i ON t1 (a)", false},
	{"EXPLAIN DROP TABLE t1", false},
	{"EXPLAIN TRUNCATE TABLE t1", false},
	{"EXPLAIN RENAME TABLE t1 TO t2", false},
	{"EXPLAIN LOAD DATA INFILE 'x' INTO TABLE t1", false},
	{"EXPLAIN SET @a = 1", false},
	{"EXPLAIN ANALYZE SET @a = 1", false},
	{"EXPLAIN FORMAT = 'row' SET @a = 1", false},
	{"EXPLAIN FORMAT = 'row' SHOW TABLES", false},
	{"EXPLAIN FORMAT = ROW SHOW TABLES", false},

	// The table form: a table name and a column name, which is one to three
	// dotted identifiers. A string is no column name here, which is the
	// other half of the MySQL difference -- MySQL reads it as a wild pattern
	// -- and it is what made `explain explore 'digestxxx'` parse before this
	// slot was a column name. EXPLORE itself is grammar TiDB added after
	// v8.5.5 and v8.5.5 answers 1064 for all three corpus rows that write it.
	{"EXPLAIN t1", true},
	{"EXPLAIN t1 a", true},
	{"EXPLAIN t1 `a`", true},
	{"EXPLAIN t1 a.b", true},
	{"EXPLAIN d1.t1 a", true},
	{"DESC t1 a", true},
	{"DESC SCHE.TABL COLUM", true},
	{"EXPLAIN BEGIN", true},
	// EXTENDED and PARTITIONS are not reserved, so each is a table name
	// here, and that is all they are: the modifier they used to be is gone.
	{"EXPLAIN EXTENDED t1", true},
	{"EXPLAIN PARTITIONS t1", true},
	{"EXPLAIN EXTENDED", true},
	{"EXPLAIN t1 'abc'", false},
	{"EXPLAIN t1 '%a%'", false},
	{"EXPLAIN t1 123", false},
	{"EXPLAIN t1 select", false},
	{"EXPLAIN t1 a b", false},
	{"DESC t1 'abc'", false},
	{"DESCRIBE t1 'abc'", false},
	{"explain explore 'digestxxx'", false},
	{"explain explore '23adc8e6f62'", false},
	{"explain explore 'select a from t'", false},
	// A modifier rules the table form out, so a table name after one is a
	// 1064 and so is a column name after that.
	{"EXPLAIN FORMAT = 'row' t1", false},
	{"EXPLAIN FORMAT = 'row' t1 a", false},
	{"EXPLAIN FORMAT = ROW t1", false},
	{"EXPLAIN ANALYZE t1", false},
	{"EXPLAIN ANALYZE t1 a", false},
	// A statement digest in place of the statement is grammar TiDB added
	// after v8.5.5 as well.
	{"EXPLAIN 'sqldigest'", false},
	{"EXPLAIN ANALYZE 'sqldigest'", false},
	{"EXPLAIN format='json' 'sqldigest'", false},
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
				t.Fatalf("Parse(%q) was accepted, but TiDB %s rejects it", c.sql, tidbCorpusVersion)
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
	t.Run("quoted format name", func(t *testing.T) {
		stmt := parseOneExplain(t, "EXPLAIN FORMAT = 'row' SELECT 1")
		if stmt.Format != "row" {
			t.Errorf("Format = %q, want %q", stmt.Format, "row")
		}
		if stmt.Analyze {
			t.Error("Analyze = true, want false")
		}
		if _, ok := stmt.Stmt.(*nodes.SelectStmt); !ok {
			t.Errorf("Stmt = %T, want *ast.SelectStmt", stmt.Stmt)
		}
	})

	t.Run("quoted format name keeps its own spelling", func(t *testing.T) {
		// The name reaches the AST as written, not folded: TiDB reads
		// 'BRIEF' and BRIEF as the same format and this parser does not
		// rewrite either into the other.
		for _, tc := range []struct{ sql, want string }{
			{"EXPLAIN FORMAT = 'BRIEF' SELECT 1", "BRIEF"},
			{"EXPLAIN FORMAT = 'brief' SELECT 1", "brief"},
			{"EXPLAIN FORMAT = BRIEF SELECT 1", "BRIEF"},
			{`EXPLAIN FORMAT = "row" SELECT 1`, "row"},
			{"EXPLAIN FORMAT = tidb_json SELECT 1", "tidb_json"},
			{"EXPLAIN FORMAT = '' SELECT 1", ""},
		} {
			stmt := parseOneExplain(t, tc.sql)
			if stmt.Format != tc.want {
				t.Errorf("%s: Format = %q, want %q", tc.sql, stmt.Format, tc.want)
			}
		}
	})

	t.Run("analyze with a format", func(t *testing.T) {
		stmt := parseOneExplain(t, "EXPLAIN ANALYZE FORMAT = 'binary' SELECT 1")
		if !stmt.Analyze {
			t.Error("Analyze = false, want true")
		}
		if stmt.Format != "binary" {
			t.Errorf("Format = %q, want %q", stmt.Format, "binary")
		}
		if stmt.Extended || stmt.Partitions {
			t.Errorf("Extended = %v, Partitions = %v, want both false: TiDB has neither modifier",
				stmt.Extended, stmt.Partitions)
		}
	})

	t.Run("for connection", func(t *testing.T) {
		stmt := parseOneExplain(t, "EXPLAIN FORMAT = 'dot' FOR CONNECTION 42")
		if stmt.Format != "dot" {
			t.Errorf("Format = %q, want %q", stmt.Format, "dot")
		}
		if stmt.ForConnection != 42 {
			t.Errorf("ForConnection = %d, want 42", stmt.ForConnection)
		}
		// The connection form explains no statement of its own, which is how
		// a consumer tells it from a connection id that happens to be 0.
		if stmt.Stmt != nil {
			t.Errorf("Stmt = %T, want nil for the FOR CONNECTION form", stmt.Stmt)
		}
		zero := parseOneExplain(t, "EXPLAIN FOR CONNECTION 0")
		if zero.ForConnection != 0 || zero.Stmt != nil {
			t.Errorf("FOR CONNECTION 0: ForConnection = %d, Stmt = %T", zero.ForConnection, zero.Stmt)
		}
		// TiDB's connection id is unsigned and ForConnection is an int64, so
		// an id above the signed range comes out saturated. That is the
		// lexer's integer literal, which every statement in this engine
		// shares, rather than anything EXPLAIN does, and it is recorded here
		// rather than fixed in this slot.
		big := parseOneExplain(t, "EXPLAIN FOR CONNECTION 18446744073709551615")
		if big.ForConnection != math.MaxInt64 {
			t.Errorf("FOR CONNECTION 18446744073709551615: ForConnection = %d, want the saturated %d",
				big.ForConnection, int64(math.MaxInt64))
		}
	})

	t.Run("alter table", func(t *testing.T) {
		stmt := parseOneExplain(t, "EXPLAIN ALTER TABLE t1 ADD INDEX (a)")
		alt, ok := stmt.Stmt.(*nodes.AlterTableStmt)
		if !ok {
			t.Fatalf("Stmt = %T, want *ast.AlterTableStmt", stmt.Stmt)
		}
		if alt.Table == nil || alt.Table.Name != "t1" {
			t.Errorf("Table = %+v, want t1", alt.Table)
		}
		if len(alt.Commands) != 1 {
			t.Fatalf("Commands = %d, want 1", len(alt.Commands))
		}
		// The explained change is still the whole ALTER TABLE, with its own
		// positions into the text the caller holds.
		const prefix = len("EXPLAIN ")
		if alt.Loc.Start != prefix {
			t.Errorf("inner Loc.Start = %d, want %d", alt.Loc.Start, prefix)
		}
		if alt.Loc.End != len("EXPLAIN ALTER TABLE t1 ADD INDEX (a)") {
			t.Errorf("inner Loc.End = %d, want %d", alt.Loc.End, len("EXPLAIN ALTER TABLE t1 ADD INDEX (a)"))
		}
	})

	t.Run("table form holds a column name", func(t *testing.T) {
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
		qualified := parseOneExplain(t, "EXPLAIN t1 x.a")
		show = qualified.Stmt.(*nodes.ShowStmt)
		col, ok = show.Like.(*nodes.ColumnRef)
		if !ok {
			t.Fatalf("Like = %T, want *ast.ColumnRef", show.Like)
		}
		if col.Table != "x" || col.Column != "a" {
			t.Errorf("column = %+v, want x.a", col)
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
			{"SELECT 1; EXPLAIN FORMAT = `row` SELECT 1", false},
			{"SELECT 1; EXPLAIN FORMAT = ROW SELECT 1", true},
			{"SELECT 1;\n/* c */ EXPLAIN FORMAT = `dot` SELECT 1", false},
			{"SELECT 1;\n/* c */ EXPLAIN FORMAT = DOT SELECT 1", true},
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
			"TABLE t1",
			"VALUES ROW(1)",
			"(SELECT 1)",
			"WITH c AS (SELECT 1) SELECT * FROM c",
			"ALTER TABLE t1 ADD INDEX (a)",
			"ALTER TABLE t1 ADD a varchar(255)",
		} {
			alone, err := Parse(inner)
			if err != nil {
				t.Fatalf("Parse(%q) error: %v", inner, err)
			}
			for _, prefix := range []string{"EXPLAIN ", "EXPLAIN ANALYZE ", "EXPLAIN FORMAT = 'brief' "} {
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

// TestExplainFormsMatchTiDB holds the table above to the engine itself
// (container test). The verdict comes from PREPARE, which parses the statement
// without running it: 1064 is the parser's rejection, and any other server
// error means the text parsed and failed later.
func TestExplainFormsMatchTiDB(t *testing.T) {
	tc := startTiDB(t)
	assertTiDBVersion(t, tc)
	for _, c := range explainFormCases {
		t.Run(c.sql, func(t *testing.T) {
			accepts, err := tc.prepareSyntax(c.sql)
			if err != nil {
				t.Fatalf("the engine gave no verdict for %q: %v", c.sql, err)
			}
			if accepts != c.accept {
				t.Errorf("TiDB %ss %q, this file records %s",
					verdictName(accepts), c.sql, verdictName(c.accept))
			}
		})
	}
}

// TestExplainAnalyzeOfAnAlterTableRuns pins the fact review's walk depends on:
// EXPLAIN ALTER TABLE reports the plan for the change and leaves the table
// alone, while EXPLAIN ANALYZE ALTER TABLE makes the change. A walk that
// treated every EXPLAIN as a statement that changes nothing would lose the
// second one's column (container test).
func TestExplainAnalyzeOfAnAlterTableRuns(t *testing.T) {
	tc := startTiDB(t)
	// Every statement names its database: the shared pool selects none, and
	// one of this package's other container tests restarts the server.
	for _, stmt := range []string{
		"CREATE DATABASE IF NOT EXISTS explain_forms",
		"DROP TABLE IF EXISTS explain_forms.ddl",
		"CREATE TABLE explain_forms.ddl (id int)",
	} {
		if _, err := tc.db.ExecContext(tc.ctx, stmt); err != nil {
			t.Fatalf("seed %q: %v", stmt, err)
		}
	}
	columns := func() []string {
		rows, err := tc.db.QueryContext(tc.ctx,
			"SELECT column_name FROM information_schema.columns WHERE table_schema = 'explain_forms' AND table_name = 'ddl' ORDER BY column_name")
		if err != nil {
			t.Fatalf("read columns: %v", err)
		}
		defer rows.Close()
		var out []string
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				t.Fatalf("scan: %v", err)
			}
			out = append(out, name)
		}
		if err := rows.Err(); err != nil {
			t.Fatalf("rows: %v", err)
		}
		return out
	}
	if _, err := tc.db.ExecContext(tc.ctx, "EXPLAIN ALTER TABLE explain_forms.ddl ADD COLUMN plain int"); err != nil {
		t.Fatalf("EXPLAIN ALTER TABLE: %v", err)
	}
	if got := strings.Join(columns(), ","); got != "id" {
		t.Errorf("after EXPLAIN ALTER TABLE columns are %q, want %q: the plan was reported, not run", got, "id")
	}
	if _, err := tc.db.ExecContext(tc.ctx, "EXPLAIN ANALYZE ALTER TABLE explain_forms.ddl ADD COLUMN analyzed int"); err != nil {
		t.Fatalf("EXPLAIN ANALYZE ALTER TABLE: %v", err)
	}
	if got := strings.Join(columns(), ","); got != "analyzed,id" {
		t.Errorf("after EXPLAIN ANALYZE ALTER TABLE columns are %q, want %q: ANALYZE runs the change", got, "analyzed,id")
	}
}
