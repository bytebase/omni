package parser

import (
	"math"
	"runtime"
	"strings"
	"testing"
	"time"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// execFallbackCase is a statement with executable comments, whether OceanBase
// accepts it, and, when it does, the statement it runs: the text without the
// executable comments it reads as plain ones, "" for an empty query.
type execFallbackCase struct {
	sql    string
	accept bool
	reads  string
}

// execFallbackCases are OceanBase CE 4.3.5.6's verdicts on executable comments
// whose bodies do not parse in place (exec_fallback.go), measured on 2026-11-05
// through PREPARE. Run as a query, each got the same verdict, and a SELECT
// returned the rows its reads returns. TestExecFallback_OceanBase keeps them a
// recording.
var execFallbackCases = []execFallbackCase{
	// A body that does not parse drops, and the statement runs without it.
	{"select 1 /*! select 2 */", true, "select 1"},
	{"select 1 /*!50000 select 2 */", true, "select 1"},
	{"select 1 /*! from */", true, "select 1"},
	{"select 1 /*!50000 from */", true, "select 1"},
	{"select 1 /*! from dual */ , 2", true, "select 1 , 2"},
	{"select 1 /*! as */ , 2", true, "select 1 , 2"},
	// The retry reads the whole statement again, so every executable comment
	// in it drops, the ones that parsed too.
	{"select 1 /*! + 1 */ /*! from */", true, "select 1"},
	{"select 1 /*! from */ /*! + 1 */", true, "select 1"},
	{"select 1 /*! + 1 */ + 1 /*! from */", true, "select 1 + 1"},
	{"select 1 /*! + 10 */ + 100 /*! , 2 */ /*! from */", true, "select 1 + 100"},
	{"select 1 /*! + 1 /*! from */", true, "select 1"},
	{"select 1 /*! from /*! + 1 */", true, "select 1"},
	// A body that parses is read in place.
	{"select 1 /*! + 1 */", true, "select 1 + 1"},
	{"select 1 /*! , 2 */", true, "select 1 , 2"},
	{"select /*! 1 */", true, "select 1"},
	{"select 1 /*! from dual */ where 1 = 1", true, "select 1 from dual where 1 = 1"},
	{"/*!50000 select 1 */ /*!50000 + 1 */", true, "select 1 + 1"},
	// At the start of the statement. A statement left with no token is an
	// empty query, which runs and returns nothing.
	{"/*! select 1 */", true, "select 1"},
	{"/*! select */ 1", true, "select 1"},
	{"/*! x */ select 1", true, "select 1"},
	{"/*! select 1 from */", true, ""},
	{"/*! select from */", true, ""},
	{"/*! select 1 */ /*! from */", true, ""},
	{"/*!50000 SELECT FROM t */", true, ""},
	{"/*! SELECT 'abc */", true, ""},
	{"/*!50100 SELECT 'x */", true, ""},
	{"/*!80000 SELECT a FROM WHERE */", true, ""},
	{"/*! select 1 */ from", false, ""},
	{"/*! select 1 */ + 1 /*! from */", false, ""},
	// A comment that never closes. The statement is refused when it parses
	// with the comment open, and read again when it does not.
	{"select 1 /*! + 1", false, ""},
	{"select 1 /*! from", false, ""},
	{"select 1 /*! + 'a*/ -- '", false, ""},
	{"select 1 /*! from 'a*/ , 2 -- '", true, "select 1 , 2"},
	{"select 1 /*! from -- c */\n", true, "select 1"},
	{"select 1 /*! , 'abc */", true, "select 1"},
	// The statement read again does not parse either.
	{"select 1 /*! , 2 */ from", false, ""},
	{"select 1 /*! from */ */", false, ""},
	{"select 1 /*! + 1 /*! + 2 */ from */", false, ""},
	// The second reading follows OceanBase's select grammar as the first does:
	// it has no TABLE query, and no locking clause right after a lone
	// parenthesized query.
	{"select 1 /*! table t */", true, "select 1"},
	{"/*! x */ table t", false, ""},
	{"(select 1) /*! for update */", true, "(select 1)"},
	{"(select 1) /*! from */ for update", false, ""},
	{"(select 1) order by 1 /*! from */ for update", true, "(select 1) order by 1 for update"},
	// Other statements the SQL parser reads.
	{"create table t (a int) /*! engine */", true, "create table t (a int)"},
	{"create table t (a int) /*!50000 select 2 */", true, "create table t (a int) select 2"},
	{"create database d /*!40100 default character set utf8mb4 */ /*!80016 default encryption='N' */", true, "create database d"},
	{"insert into t values (1) /*! from */", true, "insert into t values (1)"},
	{"update t set a = 2 /*! from */", true, "update t set a = 2"},
	{"delete from t /*! from */", true, "delete from t"},
	{"set @v = 1 /*! from */", true, "set @v = 1"},
	{"create view v as select 1 /*! from */", true, "create view v as select 1"},
	{"create or replace view v as select 1 /*! from */", true, "create or replace view v as select 1"},
	{"create /*! definer=root@'%' */ view v as select 1 /*! from */", true, "create view v as select 1"},
	{"drop table if exists t /*! from */", true, "drop table if exists t"},
	{"begin /*! from */", true, "begin"},
	{"begin work /*! from */", true, "begin work"},
	{"drop/**/function if exists f /*! from */", true, "drop function if exists f"},
	{"create aggregate function f returns int /*! from */ soname 'f.so'", true, "create aggregate function f returns int soname 'f.so'"},
	{"/*! x */ get diagnostics @n = number", true, "get diagnostics @n = number"},
	// The PL parser reads a stored-program statement and has no retry.
	{"call p() /*! from */", false, ""},
	{"call p() /*!50000 from t */", false, ""},
	{"/* c */ call p() /*! from */", false, ""},
	{"-- c\ncall p() /*! from */", false, ""},
	{"call p() /*! , 2 */", false, ""},
	{"signal sqlstate '45000' /*! from */", false, ""},
	{"drop procedure if exists p /*! from */", false, ""},
	{"alter procedure p /*! from */ comment 'x'", false, ""},
	{"create procedure p() /*! from */ begin select 1; end", false, ""},
	{"create or replace procedure p() /*! from */ select 1", false, ""},
	{"create procedure p() begin declare x int /*! from */; select 1; end", false, ""},
	{"create definer=root@'%' view v as select 1 /*! from */", false, ""},
	{"create event e /*! from */ on schedule every 1 hour do select 1", false, ""},
	{"create event e /*! if not exists */ on schedule every 1 hour do select 1", false, ""},
	{"alter event e /*! from */ disable", false, ""},
	{"drop event if exists e /*! from */", false, ""},
	// CREATE FUNCTION goes to the SQL parser once the PL parser fails, and the
	// SQL parser reads only a loadable function.
	{"create function f returns int /*! from */ soname 'f.so'", true, "create function f returns int soname 'f.so'"},
	{"create function f() returns int /*! from */ return 1", false, ""},
	{"create function f /*! from */ () returns int return 1", false, ""},
	// A retry the SQL parser has no statement for.
	{"do(1) /*! from */", false, ""},
	{"call/**/p() /*! from */", false, ""},
	{"drop/**/procedure if exists p /*! from */", false, ""},
	{"drop/**/trigger if exists tr /*! from */", false, ""},
	{"create/**/procedure p() select 1 /*! from */", false, ""},
	{"/*! x */ call p()", false, ""},
	{"/*! x */ do 1", false, ""},
	{"/*! x */ signal sqlstate '45000'", false, ""},
	{"/*! x */ create procedure p() select 1", false, ""},
	{"/*! x */ create trigger tr before insert on t for each row set new.a = 1", false, ""},
	{"/*! x */ drop trigger if exists tr", false, ""},
	{"/*! x */ alter procedure p comment 'x'", false, ""},
	{"/*! x */ create event e on schedule every 1 hour do select 1", false, ""},
}

// execVersionCases are how OceanBase CE 4.3.5.6 reads the version of an
// executable comment (execOpener): five digits after "/*!", whatever their
// value, and none when fewer than five follow. Measured on 2026-11-05 as
// execFallbackCases are.
var execVersionCases = []execFallbackCase{
	{"select 1 /*!50000 + 1 */", true, "select 1 + 1"},
	{"select 1 /*!99999 + 1 */", true, "select 1 + 1"},
	{"select /*!99999 3 */", true, "select 3"},
	{"select /*!00000 1 */", true, "select 1"},
	{"select /*!12345 */ 7", true, "select 7"},
	{"select 1 /*!50000+1 */", true, "select 1 + 1"},
	{"select 1 /*!50000*/", true, "select 1"},
	{"select /*!1 */", true, "select 1"},
	{"select /*!1234 */", true, "select 1234"},
	{"select /*!5000*/", true, "select 5000"},
	{"select 1 /*!5000*/", true, "select 1"},
	{"select /*!123456 */", true, "select 6"},
	{"select /*!500001 */", true, "select 1"},
	{"select 2 * /*!100005 */", true, "select 2 * 5"},
	{"select 1 + /*!1234567 */", true, "select 1 + 67"},
	{"select 1 /*!100000 + 1 */", true, "select 1"},
	{"select 1 /*!80000 select 2 */", true, "select 1"},
	{"select 1 /*!100000 select 2 */", true, "select 1"},
	{"select 1 /*!M + 1 */", true, "select 1"},
	{"select /*!*/ 1", true, "select 1"},
	{"select 1 /*!*/", true, "select 1"},
}

// treeText is the parse of sql without locations, which differ between a
// statement and the text it is read as.
func treeText(t *testing.T, sql string) string {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	var parts []string
	for _, item := range list.Items {
		parts = append(parts, locs.ReplaceAllString(nodes.NodeToString(item), ""))
	}
	return strings.Join(parts, "\n")
}

// TestExecFallback pins omni to execFallbackCases and execVersionCases: the
// verdict, and for an accepted statement the tree of the text OceanBase runs.
func TestExecFallback(t *testing.T) {
	for _, cases := range [][]execFallbackCase{execFallbackCases, execVersionCases} {
		for _, c := range cases {
			t.Run(c.sql, func(t *testing.T) {
				list, err := Parse(c.sql)
				if (err == nil) != c.accept {
					t.Fatalf("Parse(%q) error = %v, want accept %v", c.sql, err, c.accept)
				}
				if !c.accept {
					return
				}
				if c.reads == "" {
					if len(list.Items) != 0 {
						t.Errorf("Parse(%q) = %d statements, want none: OceanBase runs an empty query", c.sql, len(list.Items))
					}
					return
				}
				if got, want := treeText(t, c.sql), treeText(t, c.reads); got != want {
					t.Errorf("Parse(%q) is not the tree of %q:\n got %s\nwant %s", c.sql, c.reads, got, want)
				}
			})
		}
	}
}

// TestExecFallbackPerStatement: the retry reads one statement again, and the
// statements around it keep their executable comments. OceanBase 4.3.5.6 runs
// "select 1 /*! + 1 */; select 1 /*! from */" as a query and returns 2 and then
// 1, measured on 2026-11-05 (TestExecFallbackScript_OceanBase).
func TestExecFallbackPerStatement(t *testing.T) {
	for _, c := range []struct{ sql, reads string }{
		{"select 1 /*! + 1 */; select 1 /*! from */", "select 1 + 1; select 1"},
		{"select 1 /*! from */; select 1 /*! + 1 */", "select 1; select 1 + 1"},
		{"select 1; /*! select from */; select 2", "select 1; select 2"},
		{"select 1 /*! from ; select 2 */", "select 1"},
		// Under a client delimiter, a segment holds statements a ';'
		// separates, and each is read on its own.
		{"DELIMITER //\nselect 1 /*! + 1 */; select 1 /*! from */; select 2 /*! + 2 */ //\nDELIMITER ;", "select 1 + 1; select 1; select 2 + 2"},
	} {
		if got, want := treeText(t, c.sql), treeText(t, c.reads); got != want {
			t.Errorf("Parse(%q) is not the tree of %q:\n got %s\nwant %s", c.sql, c.reads, got, want)
		}
	}
}

// TestExecFallbackStaysLinear: under a client delimiter, one segment holds
// every statement, and each one read again must cost no pass over the
// segment or over its line. Each reading started a new parser, which found
// the newlines of the whole segment again for the syntax error that set off
// the retry, and the error's line was found by scanning out from the error:
// 32,000 such statements took 23 seconds one to a line and 11 seconds all on
// one line.
func TestExecFallbackStaysLinear(t *testing.T) {
	for _, stmt := range []string{
		"select 1 /*! from */;\n",
		"select 1 /*! from */; ",
		"select 1 /*! select 2 */;\n",
		"select 1 /*! select 2 */; ",
	} {
		assertLinearGrowth(t, stmt, func(n int) string {
			return "DELIMITER //\n" + strings.Repeat(stmt, n) + "//\nDELIMITER ;"
		}, 2000)
	}
}

// assertLinearGrowth parses sql(n) and sql(8n), each the fastest of three
// runs measured in CPU time, and fails when the larger took more than 24
// times as long as the smaller: linear work grows eightfold and quadratic
// work sixty-fourfold, so the bound holds on a slow machine and still catches
// quadratic work.
func assertLinearGrowth(t *testing.T, name string, sql func(n int) string, n int) {
	t.Helper()
	fastest := func(text string) time.Duration {
		best := time.Duration(math.MaxInt64)
		for range 3 {
			runtime.GC()
			start := cpuTime(t)
			if _, err := Parse(text); err != nil {
				t.Fatalf("%s: Parse: %v", name, err)
			}
			best = min(best, cpuSince(t, start))
		}
		return best
	}
	small, large := fastest(sql(n)), fastest(sql(8*n))
	if large > 24*small {
		t.Errorf("%q: %d took %v and %d took %v, %.0f times as long in CPU time where linear work takes 8",
			name, n, small, 8*n, large, float64(large)/float64(small))
	}
}

// TestExecFallbackKeepsNestingLimit: a statement that reaches the nesting
// limit or the tree height limit inside an executable comment is refused as
// it is outside one. Read again without the comment, it would lose the
// expression that went too deep and pass as something else.
func TestExecFallbackKeepsNestingLimit(t *testing.T) {
	deep := strings.Repeat("(", maxNestingDepth+10) + "1" + strings.Repeat(")", maxNestingDepth+10)
	long := strings.Repeat(" OR 1", maxTreeHeight+10)
	for _, sql := range []string{
		"select 0 /*!" + long + " */",
		"select 0 /*! + " + deep + " */",
		// The nesting limit stops the first reading before the FROM it would
		// fail on.
		"select 0 /*! + " + deep + " */ /*! from */",
	} {
		_, err := Parse(sql)
		if !isNestingError(err) {
			t.Errorf("Parse(%.40q...) error = %v, want the nesting error", sql, err)
		}
	}
	// A statement whose first reading fails before it reaches a limit is read
	// again as any other. A chain is read by a loop, so the reading of this
	// one fails at FROM with no tree built, as OceanBase's does
	// (TestExecFallbackNestingLimit_OceanBase).
	for _, sql := range []string{
		"select 0 /*!" + long + " */ /*! from */",
		"select 0 /*!" + strings.Repeat(" OR 1", 1000) + " */ /*! from */",
		"select 0 /*! + " + strings.Repeat("(", 100) + "1" + strings.Repeat(")", 100) + " */ /*! from */",
	} {
		if got, want := treeText(t, sql), treeText(t, "select 0"); got != want {
			t.Errorf("Parse(%.40q...) = %s, want %s", sql, got, want)
		}
	}
}

// TestExecFallbackLocs: a statement read again keeps every location an offset
// into the text given to the parser, from its first token to its last, and
// an error in it is at the token the second reading stops on.
func TestExecFallbackLocs(t *testing.T) {
	for _, c := range []struct {
		sql  string
		want []string // the text each statement's Loc slices
	}{
		{"select 1 /*! from */", []string{"select 1"}},
		{"select 1 /*! from dual */ , 2", []string{"select 1 /*! from dual */ , 2"}},
		{"/*! x */ select 1 /*! from */", []string{"select 1"}},
		{"select 0;\nselect 1 /*! + 1 */ /*! from */ ;\nselect 2", []string{"select 0", "select 1", "select 2"}},
	} {
		list, err := Parse(c.sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", c.sql, err)
		}
		var got []string
		for _, item := range list.Items {
			loc := locSpans(item)[0].loc
			got = append(got, c.sql[loc.Start:loc.End])
		}
		if strings.Join(got, "|") != strings.Join(c.want, "|") {
			t.Errorf("Parse(%q) statements slice %q, want %q", c.sql, got, c.want)
		}
		if v := CheckLocations(t, c.sql); len(v) > 0 {
			t.Errorf("Parse(%q) locations: %v", c.sql, v)
		}
	}

	for _, c := range []struct {
		sql  string
		want int
	}{
		// SELECT 1 * /, read again: the error is at the "/".
		{"select 1 /*! from */ */", 22},
		// SELECT 1 FROM, read again: the error is at the end of the input.
		{"select 1 /*! , 2 */ from", 24},
		{"select 0; select 1 /*! from */ from", 35},
		// A retry to a statement the SQL parser does not have fails at its start.
		{"/*! x */ call p()", 9},
		{"do(1) /*! from */", 0},
		// The PL parser has no retry, and its error is in the comment, where it
		// read the body.
		{"call p() /*!50000 from t */", 18},
	} {
		_, err := Parse(c.sql)
		pe, ok := err.(*ParseError)
		if !ok {
			t.Fatalf("Parse(%q) error = %v, want a *ParseError", c.sql, err)
		}
		if pe.Position != c.want {
			t.Errorf("Parse(%q) error at %d (%q), want %d", c.sql, pe.Position, pe.Message, c.want)
		}
	}
}

// TestRetriesAsSQL pins the routing ObParser::parse does before it parses, on
// the statements execFallbackCases measures it through and the words it reads.
func TestRetriesAsSQL(t *testing.T) {
	for _, c := range []struct {
		stmt string
		sql  bool
	}{
		{"select 1", true},
		{"  /*! call p() */", true},
		{"/*!50003 create procedure p() begin end */", true},
		{"call p()", false},
		{"CALL p()", false},
		{"/* c */ call p()", false},
		{"-- c\ncall p()", false},
		{"--c\ncall p()", false},
		{"# c\ncall p()", false},
		{"call/**/p()", true},
		{"do 1", false},
		{"do(1)", true},
		{"declare x int", false},
		{"<<label>> begin end", false},
		{"signal sqlstate '45000'", false},
		{"resignal", false},
		{"create procedure p() select 1", false},
		{"create or replace procedure p() select 1", false},
		{"create definer=root procedure p() select 1", false},
		{"create definer = root view v as select 1", false},
		{"create /*! definer=root */ view v as select 1", true},
		{"create trigger tr before insert on t for each row set @a = 1", false},
		{"create function f() returns int return 1", true},
		{"create or replace function f() returns int return 1", true},
		{"create aggregate function f returns int soname 'f.so'", true},
		{"create view v as select 1", true},
		{"create table t (a int)", true},
		{"create\tprocedure p() select 1", false},
		{"create/**/procedure p() select 1", true},
		{"drop procedure p", false},
		{"drop function f", false},
		{"drop trigger tr", false},
		{"create event e on schedule every 1 hour do select 1", false},
		{"alter event e disable", false},
		{"drop event e", false},
		{"event", false},
		{"drop definer", false},
		{"drop table t", true},
		{"alter procedure p comment 'x'", false},
		{"alter table t add b int", true},
		{"begin", true},
		{"begin work", true},
		{"begin select 1; end", false},
		{"update t set a = 1", true},
		{"update of", false},
		{"submit job", false},
		{"cancel job", false},
		{"dd begin", false},
		{"/* never closes", true},
		{"", true},
	} {
		if got := retriesAsSQL(c.stmt); got != c.sql {
			t.Errorf("retriesAsSQL(%q) = %v, want %v", c.stmt, got, c.sql)
		}
	}
}
