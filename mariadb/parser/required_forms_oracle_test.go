package parser

import (
	"maps"
	"slices"
	"testing"
)

// A column has a type, and a keyword is no function name. Both are forms the parser read
// and the engine does not, which a corpus of accepting rows cannot catch. `a srid 0` also
// used to panic the parser, on the type it had let the column go without. Every verdict
// was read from mariadb:11.8.8 on 2026-11-05.
//
// Left out because they are no part of this: `a serial default value`, which reads `value`
// as an expression, and the type list of CAST, which the parser reads as a column's.

type requiredFormCase struct {
	sql    string
	accept bool
}

var requiredFormCases = []requiredFormCase{
	// no type.
	{"CREATE TABLE tl (a check (a > 0))", false},
	{"CREATE TABLE tl (a)", false},
	{"CREATE TABLE tl (a, b int)", false},
	{"CREATE TABLE tl (a) SELECT 1 AS a", false},
	{"CREATE TABLE tl (a) AS SELECT 1 AS a", false},
	{"CREATE TABLE tl (a, b) AS SELECT 1 AS a, 2 AS b", false},
	// still reads.
	{"CREATE TABLE tl (a int key)", true},
	{"CREATE TABLE tl (a int generated always as (1))", true},
	{"CREATE TABLE tl (a int as (1) stored)", true},
	{"CREATE TABLE tl (a serial)", true},
	{"ALTER TABLE t ADD COLUMN a int comment 'x'", true},
	{"ALTER TABLE t MODIFY a int not null", true},
	{"ALTER TABLE t CHANGE a b int key", true},
	// function names still read.
	{"SELECT d.comment(1)", true},
	{"SELECT `comment`(1)", true},
	{"SELECT d.vector(1)", true},
	{"SELECT comment FROM t", true},
	{"SELECT 1 FROM comment", true},
	{"INSERT INTO t VALUES (1) ON DUPLICATE KEY UPDATE a = VALUES(a)", true},
	{"SELECT 1 = ANY (SELECT 1)", true},
	{"SELECT 1 FROM t WHERE a = SOME (SELECT 1)", true},
	{"SELECT rank() OVER ()", true},
	{"SELECT lag(a) OVER () FROM t", true},
	{"SELECT ntile(1) OVER ()", true},
	{"SET @a = comment(1)", false},
	{"DO comment(1)", false},
	{"SELECT a FROM t ORDER BY comment(1)", false},
	{"SELECT a FROM t WHERE comment(1) = 1", false},
	{"CREATE TABLE tf (a int DEFAULT (comment(1)))", false},
	{"CREATE TABLE tf (a int, b int AS (comment(1)))", false},
	{"CREATE VIEW vf AS SELECT comment(1)", false},
	{"CREATE FUNCTION comment() RETURNS INT RETURN 1", true},
	{"CALL comment(1)", true},
	{"DROP FUNCTION comment", true},
}

// TestRequiredForms holds the parser to the recorded verdicts without the server.
func TestRequiredForms(t *testing.T) {
	for _, c := range requiredFormCases {
		t.Run(c.sql, func(t *testing.T) {
			_, err := Parse(c.sql)
			if omniAccepts := err == nil; omniAccepts != c.accept {
				t.Errorf("omni accepts=%v, the record says %v for %q: %v", omniAccepts, c.accept, c.sql, err)
			}
		})
	}
}

// TestKeywordFunctionNames: a listed word is refused as a call, and only as a call of that
// spelling.
func TestKeywordFunctionNames(t *testing.T) {
	if len(keywordsNotFunctionNames) == 0 {
		t.Fatal("the list is empty")
	}
	for word := range keywordsNotFunctionNames {
		for _, refused := range []string{"SELECT " + word + "(1)", "SELECT " + word + "()", "SELECT 1 + " + word + " (1)"} {
			if _, err := Parse(refused); err == nil {
				t.Errorf("omni reads %q", refused)
			}
		}
		for _, read := range []string{"SELECT `" + word + "`(1)", "SELECT d." + word + "(1)"} {
			if _, err := Parse(read); err != nil {
				t.Errorf("omni refuses %q: %v", read, err)
			}
		}
	}
}

// functionNameShapes are the argument lists the sweep tries a word with: a word is taken
// as a function name if the engine parses it in any of them, which leaves out every word
// with a grammar of its own, such as the OVER of a window function.
var functionNameShapes = []string{
	"()", "(1)", "(1,2)", "(1,2,3)", "(1,2,3,4)", "(*)", "('a')", "(a USING utf8mb4)", "(1 AS SIGNED)",
	"(1, INTERVAL 1 DAY)", "(1 FROM 2)", "(1 IN 2)", "(DISTINCT 1)", "() OVER ()", "(1) OVER ()", "(1,2) OVER ()",
	"(1,2,3) OVER ()", "(1) FROM t", "(SELECT 1)", "(YEAR FROM 1)", "(DAY, 1, 2)", "(a)", "(a, b)", "(a) OVER ()",
}

// functionArmsNotCovered is what the list cannot reach: `values` is read by an arm of its
// own for ON DUPLICATE KEY UPDATE, and the engine refuses it everywhere else.
var functionArmsNotCovered = map[string]bool{"values": true}

// judgeKeywordFunctionNames sweeps every keyword the lexer has against the engine.
func judgeKeywordFunctionNames(t *testing.T, accepts func(sql string) bool) {
	t.Helper()
	// In name order: a test that asks the server asks in the same order every run.
	for _, word := range slices.Sorted(maps.Keys(keywords)) {
		var engine, omni bool
		for _, shape := range functionNameShapes {
			sql := "SELECT " + word + shape
			engine = engine || accepts(sql)
			if _, err := Parse(sql); err == nil {
				omni = true
			}
		}
		if keywordsNotFunctionNames[word] && engine {
			t.Errorf("%q is listed and the engine parses it as a function name", word)
		}
		if !engine && omni && !functionArmsNotCovered[word] {
			t.Errorf("the engine refuses %q as a function name and omni reads it: add it to the list", word)
		}
	}
}

func TestRequiredFormsOracle(t *testing.T) {
	o := startMariaDB(t)
	for _, c := range requiredFormCases {
		t.Run(c.sql, func(t *testing.T) {
			if got := o.prepareAccepts(t, c.sql); got != c.accept {
				t.Errorf("MariaDB 11.8.8 accepts=%v, the record says %v for %q", got, c.accept, c.sql)
			}
		})
	}
}

func TestKeywordFunctionNamesOracle(t *testing.T) {
	o := startMariaDB(t)
	judgeKeywordFunctionNames(t, func(sql string) bool { return o.prepareAccepts(t, sql) })
}

// typelessColumnCases read the keyword after the name as a type, and the grammar's action
// answers 4161 for it. They can never run, so the parser refuses them, as it does any
// other name that is no type.
var typelessColumnCases = []parseCodeCase{
	{"CREATE TABLE tl (a key)", 4161},
	{"CREATE TABLE tl (a primary key)", 4161},
	{"CREATE TABLE tl (a unique)", 4161},
	{"CREATE TABLE tl (a unique key)", 4161},
	{"CREATE TABLE tl (a comment 'x')", 4161},
	{"CREATE TABLE tl (a not null)", 4161},
	{"CREATE TABLE tl (a null)", 4161},
	{"CREATE TABLE tl (a default 1)", 4161},
	{"CREATE TABLE tl (a auto_increment)", 4161},
	{"CREATE TABLE tl (a collate utf8mb4_bin)", 4161},
	{"CREATE TABLE tl (a references r (id))", 4161},
	{"CREATE TABLE tl (a generated always as (1))", 4161},
	{"CREATE TABLE tl (a generated always as (1) stored)", 4161},
	{"CREATE TABLE tl (a generated always as (1) virtual)", 4161},
	{"CREATE TABLE tl (a as (1))", 4161},
	{"CREATE TABLE tl (a as (1) stored)", 4161},
	{"CREATE TABLE tl (a on update now())", 4161},
	{"CREATE TABLE tl (a column_format fixed)", 4161},
	{"CREATE TABLE tl (a storage disk)", 4161},
	{"CREATE TABLE tl (a srid 0)", 4161},
	{"CREATE TABLE tl (a invisible)", 4161},
	{"CREATE TABLE tl (a visible)", 4161},
	{"CREATE TABLE tl (a auto_random)", 4161},
	{"CREATE TABLE tl (a key) SELECT 1 AS a", 4161},
	{"CREATE TABLE tl (a default 1) SELECT 1 AS a", 4161},
	{"CREATE TABLE tl (a not null) AS SELECT 1 AS a", 4161},
	{"CREATE TABLE tl (a comment 'x') AS SELECT 1 AS a", 4161},
	{"CREATE TEMPORARY TABLE tl (a key)", 4161},
	{"ALTER TABLE t ADD a key", 4161},
	{"ALTER TABLE t ADD COLUMN a comment 'x'", 4161},
	{"ALTER TABLE t ADD COLUMN a generated always as (1)", 4161},
	{"ALTER TABLE t ADD (a key)", 4161},
	{"ALTER TABLE t ADD COLUMN (a not null, b int)", 4161},
	{"ALTER TABLE t MODIFY a comment 'x'", 4161},
	{"ALTER TABLE t MODIFY COLUMN a not null", 4161},
	{"ALTER TABLE t MODIFY a generated always as (1)", 4161},
	{"ALTER TABLE t CHANGE a b key", 4161},
	{"ALTER TABLE t CHANGE COLUMN a b comment 'x'", 4161},
	{"ALTER TABLE t CHANGE a b generated always as (1)", 4161},
}

// TestTypelessColumns holds the parser to them without the server.
func TestTypelessColumns(t *testing.T) {
	for _, c := range typelessColumnCases {
		if _, err := Parse(c.sql); err == nil {
			t.Errorf("omni reads %q", c.sql)
		}
	}
}

func TestTypelessColumnsOracle(t *testing.T) {
	assertParseCodes(t, typelessColumnCases)
}
