package parser

import (
	"errors"
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// Two statements need a ';' between them. Parse read a statement and then
// started on whatever token came next, so "SELECT 1 SELECT 2" parsed as two
// statements where the server answers 1064, and a tail the grammar refuses
// read as a second statement: "SELECT 1 = ANY SELECT 1" passed that way.
//
// Requiring the separator showed where this parser had leaned on the gap.
// CREATE TABLE read only SELECT as its query after the column list, the table
// options or the partitioning, so "CREATE TABLE t (a INT) VALUES ROW(1)"
// parsed as a CREATE TABLE followed by VALUES ROW(1).
// createTableQueryFormCases records every position CREATE TABLE takes its
// query in, and the forms around them.
//
// MariaDB ignores an executable comment whose version reads as MySQL's
// (/*!80000 ... */), so the statement inside it is no second statement there,
// and it has no TABLE query or VALUES ROW constructor.
//
// Every case was measured against MariaDB 11.8.8 (mariadb:11.8.8) through
// PREPARE on 2026-11-05. TestStatementSeparatorFormsOracle keeps it a
// recording: an image upgrade that changes an answer fails there rather than
// leaving a stale expectation behind.

// separatorFormCase is one statement text and whether the engine's parser
// accepts it.
type separatorFormCase struct {
	sql    string
	accept bool
}

var separatorFormCases = []separatorFormCase{
	// Two statements with no separator between them, across statement kinds.
	{"select 1 select 2", false},
	{"select 1 update t set a=1", false},
	{"create table t (a int) create table u (a int)", false},
	{"begin commit", false},
	{"commit begin", false},
	{"rollback select 1", false},
	{"set @a = 1 select 1", false},
	{"set @a = 1 set @b = 2", false},
	{"select 1 from dual where 1 = 1 select 2", false},
	{"select 1 union select 2 select 3", false},
	{"select 1 for update select 2", false},
	{"with c as (select 1) select * from c select 2", false},
	{"delete from t where a = 1 delete from t", false},
	{"insert into t values (1) insert into t values (2)", false},
	{"update t set a = 1 update t set a = 2", false},
	{"show tables show databases", false},
	{"drop table t drop table u", false},
	{"alter table t add column b int alter table t drop column b", false},
	{"create view v as select 1 select 2", false},
	{"explain select 1 select 2", false},
	{"do 1 do 2", false},
	{"handler t open handler t close", false},
	{"prepare s from 'select 1' execute s", false},
	{"lock tables t read unlock tables", false},
	{"xa start 'x' xa end 'x'", false},
	{"signal sqlstate '45000' select 1", false},
	{"table t table u", false},
	{"values row(1) values row(2)", false},
	// A tail the grammar refuses, which a missing separator check read as a
	// second statement.
	{"select 1 = any select 1", false},
	// A comment between the two statements is no separator.
	{"select 1 -- c\nselect 2", false},
	{"select 1 # c\nselect 2", false},
	{"select 1 /* c */ select 2", false},
	// An executable comment's body is read in place, so a statement inside one
	// follows the first with no separator.
	{"select 1 /*! select 2 */", false},
	{"select 1 /*!80000 select 2 */", true},
	{"select 1 /*!100000 select 2 */", false},
	// A statement after a stored program body, with no separator after its END.
	{"create procedure p() begin select 1; end select 2", false},
	{"create procedure p() select 1 select 2", false},
	{"create trigger tr before insert on t for each row begin set new.a = 1; end insert into t values (1)", false},
	{"create trigger tr before insert on t for each row set new.a = 1 insert into t values (1)", false},
	// Inside a stored program body, every declaration and statement ends with
	// ';'.
	{"create procedure p() begin select 1 select 2; end", false},
	{"create procedure p() begin select 1; select 2 end", false},
	{"create procedure p() begin if 1 then select 1 end if; end", false},
	{"create procedure p() begin declare a int declare b int; end", false},
	{"create procedure p() begin begin select 1; end end", false},
	{"create procedure p() begin while 1 do select 1 end while; end", false},
	{"create procedure p() begin select 1; select 2; end", true},
	{"create procedure p() begin end", true},
	{"create procedure p() begin if 1 then select 1; end if; end", true},
	{"create procedure p() begin declare a int; declare b int; end", true},
	{"create procedure p() begin begin select 1; end; end", true},
	// One statement, with or without its own ';'.
	{"select 1;", true},
	{"select 1 ;", true},
	{"select 1; -- c", true},
	{"begin", true},
	{"begin work", true},
}

var createTableQueryFormCases = []separatorFormCase{
	// The query after the column list, with no AS: any query expression.
	{"create table t (a int) select 1", true},
	{"create table t (a int) (select 1)", true},
	{"create table t (a int) (select 1) union (select 2)", true},
	{"create table t (a int) (select 1) order by 1", true},
	{"create table t (a int) ((select 1))", true},
	{"create table t (a int) with c as (select 1) select * from c", true},
	{"create table t (a int) with c as (select 1) table c", false},
	{"create table t (a int) table u", false},
	{"create table t (a int) (table u)", false},
	{"create table t (a int) table u order by a limit 1", false},
	{"create table t (a int) table u union select 1", false},
	{"create table t (a int) values row(1)", false},
	{"create table t (a int) (values row(1))", false},
	{"create table t (a int) values row(1) order by 1", false},
	{"create table t (a int) values row(1) union table u", false},
	{"create table t (a int) values (1)", true},
	// The query with no column list.
	{"create table t select 1", true},
	{"create table t table u", false},
	{"create table t values row(1)", false},
	{"create table t values row(1), row(2)", false},
	{"create table t values (1)", true},
	{"create table t with c as (select 1) select * from c", true},
	// A parenthesized query right after the table name.
	{"create table t (select 1)", true},
	{"create table t (select 1) union (select 2)", true},
	{"create table t (select 1) order by 1", true},
	{"create table t ((select 1))", true},
	{"create table t ((select 1)) union (select 2)", true},
	{"create table t (with c as (select 1) select * from c)", true},
	{"create table t (table u)", false},
	{"create table t (values row(1))", false},
	// After the table options, the partitioning, or both.
	{"create table t (a int) engine=innodb (select 1)", true},
	{"create table t (a int) engine=innodb table u", false},
	{"create table t (a int) engine=innodb values row(1)", false},
	{"create table t (a int) engine=innodb with c as (select 1) select * from c", true},
	{"create table t engine=innodb table u", false},
	{"create table t (a int) partition by hash(a) partitions 2 table u", false},
	{"create table t (a int) partition by hash(a) partitions 2 values row(1)", false},
	{"create table t (a int) partition by hash(a) partitions 2 with c as (select 1) select * from c", true},
	{"create table t (a int) partition by hash(a) (partition p0) (select 1)", true},
	{"create table t (a int) partition by hash(a) partitions 2 (select 1)", false},
	// A comma stands between two table options and never ends the list.
	{"create table t (a int) engine=innodb comment 'x'", true},
	{"create table t (a int) engine=innodb, comment 'x'", true},
	{"create table t (a int) engine=innodb,", false},
	{"create table t (a int) engine=innodb, select 1", false},
	{"create table t (a int) engine=innodb, (select 1)", false},
	{"create table t (a int) engine=innodb, partition by hash(a)", false},
	{"create table t (a int) engine=innodb,, comment 'x'", false},
	{"create table t engine=innodb, select 1", false},
	// IGNORE or REPLACE, and the query that has to follow it.
	{"create table t (a int) ignore (select 1)", true},
	{"create table t (a int) ignore table u", false},
	{"create table t (a int) ignore with c as (select 1) select * from c", true},
	{"create table t (a int) replace table u", false},
	{"create table t (a int) replace values row(1)", false},
	{"create table t (a int) partition by hash(a) partitions 2 ignore (select 1)", true},
	{"create table t ignore select 1", true},
	{"create table t (a int) ignore as select 1", true},
	{"create table t (a int) ignore as table u", false},
	{"create table t (a int) as ignore select 1", false},
	{"create table t (a int) ignore", false},
	{"create table t (a int) replace", false},
	{"create table t ignore", false},
	// AS, then the query.
	{"create table t as table u", false},
	{"create table t as values row(1)", false},
	{"create table t as (select 1)", true},
	{"create table t (a int) as (select 1)", true},
	{"create table t (a int) as table u", false},
	{"create table t (a int) as with c as (select 1) select * from c", true},
	{"create table t (a int) as", false},
	{"create table t (a int) ignore as", false},
	{"create table t as", false},
	// A locking clause on the query.
	{"create table t (a int) select 1 for update", true},
	{"create table t (a int) (select 1) for update", true},
	{"create table t (a int) table u for update", false},
	{"create table t (a int) (select 1) union (select 2) for update", true},
	{"create table t (a int) select 1 union select 2 for update", true},
	{"create table t (a int) (select 1) order by 1 for update", true},
	{"create table t as (select 1) for update", true},
	{"create table t (a int) as (select 1) for update", true},
	// The query still ends the statement.
	{"create table t (a int) as select 1 select 2", false},
	{"create table t (a int) table u table v", false},
	{"create table t (a int) values row(1) select 1", false},
	{"create table t (a int) (select 1) (select 2)", false},
}

func TestStatementSeparatorForms(t *testing.T) {
	for _, cases := range [][]separatorFormCase{separatorFormCases, createTableQueryFormCases} {
		for _, c := range cases {
			t.Run(c.sql, func(t *testing.T) {
				_, err := Parse(c.sql)
				if accepts := err == nil; accepts != c.accept {
					t.Errorf("omni accepts=%v, want %v for %q: %v", accepts, c.accept, c.sql, err)
				}
			})
		}
	}
}

// TestCreateTableQueryIsOneStatement holds the tree: an accepted CREATE TABLE
// with a query is one CreateTableStmt that carries the query, not a CREATE
// TABLE followed by a second statement.
func TestCreateTableQueryIsOneStatement(t *testing.T) {
	for _, c := range createTableQueryFormCases {
		if !c.accept {
			continue
		}
		t.Run(c.sql, func(t *testing.T) {
			list, err := Parse(c.sql)
			if err != nil {
				t.Fatalf("Parse(%q): %v", c.sql, err)
			}
			if len(list.Items) != 1 {
				t.Fatalf("Parse(%q) = %d statements, want 1", c.sql, len(list.Items))
			}
			stmt, ok := list.Items[0].(*nodes.CreateTableStmt)
			if !ok {
				t.Fatalf("Parse(%q) = %T, want *ast.CreateTableStmt", c.sql, list.Items[0])
			}
			hasQuery := strings.Contains(c.sql, "select") || strings.Contains(c.sql, "table u") ||
				strings.Contains(c.sql, "values")
			if got := stmt.Select != nil; got != hasQuery {
				t.Errorf("Parse(%q).Select set=%v, want %v", c.sql, got, hasQuery)
			}
		})
	}
}

// TestStatementSeparatorScripts holds what the separator check leaves alone and
// where its error lands. Split cuts a script at each ';' outside a stored
// program body, and at the DELIMITER in force, so a segment can still hold
// statements a ';' separates (under "DELIMITER //"); the check is that a
// statement ends at a ';' or the end of its segment, not that a segment holds
// one statement. A rejection points at the first token of the statement that
// had no separator in front of it, as an offset into the whole script.
func TestStatementSeparatorScripts(t *testing.T) {
	for _, c := range []struct {
		sql   string
		stmts int
	}{
		{"", 0},
		{";", 0},
		{"-- c", 0},
		{"select 1;", 1},
		{"select 1; select 2", 2},
		{"select 1;;select 2;", 2},
		{"select 1 ; -- c\n", 1},
		{"/* c */ select 1 /* c */;", 1},
		{"begin; select 1; commit", 3},
		{"create procedure p() begin select 1; select 2; end; select 3;", 2},
		{"create table t (a int) (select 1); select 1", 2},
		{"DELIMITER //\ncreate procedure p() begin select 1; end//\nDELIMITER ;\nselect 1;", 2},
		{"DELIMITER //\nselect 1; select 2//\nDELIMITER ;", 2},
		{"DELIMITER //\ncreate procedure p() begin select 1; end; select 2//\nDELIMITER ;", 2},
	} {
		list, err := Parse(c.sql)
		if err != nil {
			t.Errorf("Parse(%q): %v", c.sql, err)
			continue
		}
		if got := len(list.Items); got != c.stmts {
			t.Errorf("Parse(%q) = %d statements, want %d", c.sql, got, c.stmts)
		}
	}

	for _, c := range []struct {
		sql string
		at  string // the error is at the last occurrence of this text
		msg string // the message starts with this
	}{
		{"select 1 select 2", "select 2", "syntax error"},
		{"select 1;\nselect 2 select 3", "select 3", "syntax error"},
		{"select 1; select 2\nselect 3", "select 3", "syntax error"},
		{"DELIMITER //\nselect 1 select 2//\nDELIMITER ;", "select 2//", "syntax error"},
		{"create procedure p() begin select 1; end select 2", "select 2", "syntax error"},
		{"create table t (a int) engine=innodb, (select 1)", "(select 1)", "syntax error"},
		// A token the lexer could not finish is the more specific cause.
		{"select 1 'abc", "'abc", "unterminated"},
	} {
		_, err := Parse(c.sql)
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Errorf("Parse(%q) = %v, want a *ParseError", c.sql, err)
			continue
		}
		pos := strings.LastIndex(c.sql, c.at)
		line := strings.Count(c.sql[:pos], "\n") + 1
		col := pos - strings.LastIndex(c.sql[:pos], "\n")
		if pe.Position != pos || pe.Line != line || pe.Column != col || !strings.HasPrefix(pe.Message, c.msg) {
			t.Errorf("Parse(%q) error = %q at %d (line %d, column %d), want %q... at %d (line %d, column %d)",
				c.sql, pe.Message, pe.Position, pe.Line, pe.Column, c.msg, pos, line, col)
		}
	}
}
