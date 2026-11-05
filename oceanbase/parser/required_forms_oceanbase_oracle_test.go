package parser

import (
	"context"
	"testing"
)

// A column has a type, and a keyword is no function name. Both are forms the parser read
// and the engine does not, which a corpus of accepting rows cannot catch. `a srid 0` also
// used to panic the parser, on the type it had let the column go without. Every verdict
// was read from oceanbase/oceanbase-ce:4.3.5-lts on 2026-11-05.
//
// Left out because they are no part of this: `a serial default value`, which reads `value`
// as an expression, and the type list of CAST, which the parser reads as a column's.

type requiredFormCase struct {
	sql    string
	accept bool
}

var requiredFormCases = []requiredFormCase{
	{"SELECT _latin1 X'61'", true},
	{"SELECT _latin1 0x61", true},
	{"SELECT _latin1 b'01100001'", true},
	{"SELECT JSON_VALUE('{}', '$.a' DEFAULT _latin1 X'61' ON EMPTY)", true},
	{"SELECT JSON_VALUE('{}', '$.a' DEFAULT b'' ON EMPTY)", true},
	{"SELECT JSON_VALUE('{}', '$.a' DEFAULT +~1 ON EMPTY)", false},
	{"SELECT JSON_VALUE('{}', '$.a' DEFAULT +!1 ON EMPTY)", false},
	{"SELECT JSON_VALUE('{}', '$.a' DEFAULT -~1 ON EMPTY)", false},
	{"SELECT JSON_VALUE('{}', '$.a' DEFAULT ~1 ON EMPTY)", false},
	{"SELECT JSON_VALUE('{}', '$.a' DEFAULT - 1 ON EMPTY)", true},
	{"SELECT JSON_VALUE('{}', '$.a' DEFAULT + 1 ON EMPTY)", true},
	{"SELECT JSON_VALUE('{}', '$.a' DEFAULT --1 ON EMPTY)", false},
	{"SELECT JSON_VALUE('{}', '$.a' DEFAULT -+1 ON EMPTY)", false},
	{"SELECT JSON_VALUE('{}', '$.a' DEFAULT +-1 ON EMPTY)", false},
	{"SELECT JSON_VALUE('{}', '$.a' DEFAULT -0x61 ON EMPTY)", false},
	{"SELECT JSON_VALUE('{}', '$.a' DEFAULT -b'1' ON EMPTY)", false},
	{"SELECT JSON_VALUE('{}', '$.a' DEFAULT -1.5 ON EMPTY)", true},
	{"SELECT JSON_VALUE('{}', '$.a' DEFAULT -1e2 ON EMPTY)", true},
	{"SELECT JSON_VALUE('{}', '$.a' DEFAULT +1.5 ON EMPTY)", true},
	{"SELECT JSON_VALUE('{}', '$.a' DEFAULT -.5 ON EMPTY)", true},
	{"SELECT JSON_VALUE('{}', '$.a' DEFAULT - -1 ON EMPTY)", false},
	{"SELECT JSON_VALUE('{}', '$.a' DEFAULT -(1) ON EMPTY)", false},
	{"SELECT JSON_VALUE('{}', '$.a' DEFAULT -'1' ON EMPTY)", false},
	{"SELECT JSON_VALUE('{}', '$.a' DEFAULT -TRUE ON EMPTY)", false},
	{"SELECT JSON_VALUE('{}', _latin1 X'24')", true},
	{"SELECT JSON_VALUE('{}', _latin1 b'00100100')", true},
	{"SELECT JSON_VALUE('{}', _latin1 0x24)", true},
	{"SELECT JSON_VALUE('{}', _latin1'$')", true},
	{"SELECT JSON_VALUE('{}', '$' DEFAULT OFF ON EMPTY)", false},
	{"SELECT JSON_VALUE('{}', '$' DEFAULT FALSE ON EMPTY)", true},
	{"SELECT JSON_VALUE((1 IN (SELECT 1)), '$' NULL ON EMPTY)", true},
	{"SELECT JSON_VALUE(1 IN (SELECT 1), '$' NULL ON EMPTY)", false},
	{"SELECT JSON_VALUE((1 = 1), '$')", true},
	{"SELECT JSON_VALUE(1 = 1, '$')", false},
	{"SELECT JSON_VALUE('{}', '$' DEFAULT b'2' ON EMPTY)", false},
	{"SELECT JSON_VALUE('{}', '$' DEFAULT X'6' ON EMPTY)", false},
	{"SELECT JSON_VALUE('{}', '$' DEFAULT X'GG' ON EMPTY)", false},
	{"SELECT JSON_VALUE('{}', '$' DEFAULT 0x6 ON EMPTY)", true},
	{"SELECT JSON_VALUE('{}', '$' DEFAULT b'1' ON EMPTY)", true},
	{"SELECT `_latin1` X'61'", false},
	{"SELECT JSON_VALUE('{}', '$' DEFAULT `_latin1` X'61' ON EMPTY)", false},
	{"SELECT `_latin1`'x'", true},
	{"SELECT _latin1 1", false},
	{"SELECT _latin1 X'GG'", false},
	{"SELECT _latin1 b'2'", false},
	{"SELECT JSON_VALUE(a,'$.a' DEFAULT N'x' ON ERROR)", true},
	{"SELECT JSON_VALUE(a,'$.a' DEFAULT n'x' ON EMPTY)", true},
	{"SELECT JSON_VALUE(a,'$.a' DEFAULT N'x' 'y' ON EMPTY)", true},
	{"SELECT JSON_VALUE(a,'$.a' DEFAULT N 'x' ON ERROR)", false},
	{"SELECT JSON_VALUE(a,'$.a' DEFAULT -N'x' ON ERROR)", false},
	{"SELECT JSON_VALUE(a,'$.a' DEFAULT N ON ERROR)", false},
	{"SELECT JSON_VALUE(a,'$.a' DEFAULT -'x' ON ERROR)", false},
	{"SELECT JSON_VALUE(a,'$.a' DEFAULT -TRUE ON ERROR)", false},
	{"SELECT JSON_VALUE(a,'$.a' DEFAULT +'x' ON ERROR)", false},
	{"SELECT JSON_VALUE(a,'$.a' DEFAULT -X'41' ON ERROR)", false},
	{"SELECT JSON_VALUE(a,'$.a' DEFAULT -1.5 ON ERROR)", true},
	{"SELECT JSON_VALUE(a,'$.a' DEFAULT - 1 ON ERROR)", true},
	{"SELECT JSON_VALUE(a,'$.a' DEFAULT --1 ON ERROR)", false},
	// JSON_VALUE has its own arm: a simple_expr, a string literal, then its clauses. RETURNING INT,
	// VARCHAR and TEXT wait for the CAST list in bytebase/omni#669, which this reads through.
	{"SELECT JSON_VALUE('{}','$.a')", true},
	{"SELECT json_value(a,'$.a')", true},
	{"SELECT JSON_VALUE(a,'$.a' RETURNING CHAR(3))", true},
	{"SELECT JSON_VALUE(a,'$.a' RETURNING SIGNED)", true},
	{"SELECT JSON_VALUE(a,'$.a' RETURNING INT)", true},
	{"SELECT JSON_VALUE(a,'$.a' RETURNING DECIMAL(5,2))", true},
	{"SELECT JSON_VALUE(a,'$.a' RETURNING JSON)", true},
	{"SELECT JSON_VALUE(a,'$.a' RETURNING CHAR(3) ARRAY)", false},
	{"SELECT JSON_VALUE(a,'$.a' RETURNING)", false},
	{"SELECT JSON_VALUE(a,'$.a' NULL ON ERROR)", true},
	{"SELECT JSON_VALUE(a,'$.a' ERROR ON ERROR)", true},
	{"SELECT JSON_VALUE(a,'$.a' DEFAULT 'x' ON ERROR)", true},
	{"SELECT JSON_VALUE(a,'$.a' DEFAULT 1 ON ERROR)", true},
	{"SELECT JSON_VALUE(a,'$.a' DEFAULT b ON ERROR)", false},
	{"SELECT JSON_VALUE(a,'$.a' NULL ON EMPTY)", true},
	{"SELECT JSON_VALUE(a,'$.a' ERROR ON EMPTY)", true},
	{"SELECT JSON_VALUE(a,'$.a' DEFAULT 'x' ON EMPTY)", true},
	{"SELECT JSON_VALUE(a,'$.a' NULL ON EMPTY NULL ON ERROR)", true},
	{"SELECT JSON_VALUE(a,'$.a' NULL ON ERROR NULL ON EMPTY)", false},
	{"SELECT JSON_VALUE(a,'$.a' NULL ON EMPTY NULL ON EMPTY)", false},
	{"SELECT JSON_VALUE(a,'$.a' NULL ON ERROR NULL ON ERROR)", false},
	{"SELECT JSON_VALUE(a,'$.a' RETURNING CHAR(3) NULL ON EMPTY DEFAULT 'x' ON ERROR)", true},
	{"SELECT JSON_VALUE(a,'$.a' NULL ON EMPTY RETURNING CHAR(3))", false},
	{"SELECT JSON_VALUE(a,'$.a' RETURNING CHAR(3) ASCII)", true},
	{"SELECT JSON_VALUE(a,'$.a' RETURNING CHAR(3) TRUNCATE)", true},
	{"SELECT JSON_VALUE(a,'$.a' TRUNCATE)", true},
	{"SELECT JSON_VALUE(a,'$.a' ASCII)", true},
	{"SELECT JSON_VALUE(a,'$.a',1)", false},
	{"SELECT JSON_VALUE(a)", false},
	{"SELECT JSON_VALUE()", false},
	{"SELECT JSON_VALUE(1,2)", false},
	{"SELECT JSON_VALUE(a,b)", false},
	{"SELECT JSON_VALUE(a,concat('$','a'))", false},
	{"SELECT JSON_VALUE(a,_utf8mb4'$.a')", true},
	{"SELECT JSON_VALUE(a,N'$.a')", true},
	{"SELECT JSON_VALUE(a,X'24')", false},
	{"SELECT JSON_VALUE(a,\"$.a\")", true},
	{"SELECT JSON_VALUE(a,'$' '.a')", true},
	{"SELECT JSON_VALUE(1+1,'$.a')", false},
	{"SELECT JSON_VALUE(a+1,'$.a')", false},
	{"SELECT JSON_VALUE((1+1),'$.a')", true},
	{"SELECT JSON_VALUE(a->'$.b','$.a')", true},
	{"SELECT JSON_VALUE(NOT a,'$.a')", false},
	{"SELECT JSON_VALUE(a IS NULL,'$.a')", false},
	{"SELECT JSON_VALUE(DISTINCT a,'$.a')", false},
	{"SELECT JSON_VALUE(*,'$.a')", false},
	{"SELECT `json_value`(a,'$.a')", true},
	{"SELECT d.json_value(a,'$.a')", true},
	{"SELECT JSON_VALUE (a,'$.a')", true},
	{"SELECT json_value FROM t", true},
	{"CREATE TABLE json_value (json_value INT)", true},
	{"CREATE TABLE t (a JSON, b INT AS (JSON_VALUE(a,'$.b' RETURNING SIGNED)))", true},
	{"CREATE TABLE t (a JSON, KEY k ((JSON_VALUE(a,'$.b' RETURNING CHAR(3)))))", true},
	{"SELECT JSON_VALUE(a,'$.a' DEFAULT -1 ON ERROR)", true},
	{"SELECT JSON_VALUE(a,'$.a' DEFAULT +1 ON ERROR)", true},
	{"SELECT JSON_VALUE(a,'$.a' DEFAULT NULL ON ERROR)", true},
	{"SELECT JSON_VALUE(a,'$.a' DEFAULT 1.5 ON ERROR)", true},
	{"SELECT JSON_VALUE(a,'$.a' DEFAULT TRUE ON ERROR)", true},
	{"SELECT JSON_VALUE(a,'$.a' DEFAULT DATE'2000-01-01' ON ERROR)", true},
	{"SELECT JSON_VALUE(a,'$.a' DEFAULT X'41' ON ERROR)", true},
	{"SELECT JSON_VALUE(a,'$.a' DEFAULT _utf8mb4'x' ON ERROR)", true},
	{"SELECT JSON_VALUE(a,'$.a' DEFAULT 1+1 ON ERROR)", false},
	{"SELECT JSON_VALUE(a,'$.a' DEFAULT (1) ON ERROR)", false},
	{"SELECT JSON_VALUE(a,'$.a' DEFAULT @x ON ERROR)", false},
	{"SELECT JSON_VALUE(a,'$.a' DEFAULT ? ON ERROR)", false},
	{"SELECT JSON_VALUE(a,'$.a' RETURNING SIGNED INT)", true},
	{"SELECT JSON_VALUE(a,'$.a' RETURNING UNSIGNED)", true},
	{"SELECT JSON_VALUE(a,'$.a' RETURNING DATE)", true},
	{"SELECT JSON_VALUE(a,'$.a' RETURNING DATETIME(3))", true},
	{"SELECT JSON_VALUE(a,'$.a' RETURNING FLOAT)", true},
	{"SELECT JSON_VALUE(a,'$.a' RETURNING DOUBLE)", true},
	{"SELECT JSON_VALUE(a,'$.a' RETURNING BINARY(3))", true},
	{"SELECT JSON_VALUE(a,'$.a' RETURNING YEAR)", true},
	{"SELECT JSON_VALUE(a,'$.a' RETURNING POINT)", true},
	{"SELECT JSON_VALUE(a,'$.a' RETURNING SIGNED RETURNING SIGNED)", false},
	{"SELECT JSON_VALUE(a,'$.a' RETURNING SIGNED TRUNCATE ASCII)", true},
	{"SELECT JSON_VALUE(a,'$.a' RETURNING CHAR(3) ASCII TRUNCATE)", false},
	{"SELECT JSON_VALUE(a,'$.a' RETURNING CHAR(3) TRUNCATE ASCII)", true},
	{"SELECT JSON_VALUE(a,'$.a' ASCII TRUNCATE)", false},
	{"SELECT JSON_VALUE(a,'$.a' TRUNCATE NULL ON EMPTY)", true},
	{"SELECT JSON_VALUE(a,'$.a' NULL ON EMPTY TRUNCATE)", false},
	{"SELECT JSON_VALUE(a,'$.a' RETURNING SIGNED ASCII)", true},
	{"SELECT JSON_VALUE(-a,'$.a')", true},
	{"SELECT JSON_VALUE(!a,'$.a')", true},
	{"SELECT JSON_VALUE(a COLLATE utf8mb4_bin,'$.a')", true},
	{"SELECT JSON_VALUE(f(a),'$.a')", true},
	{"SELECT JSON_VALUE(@x,'$.a')", true},
	{"SELECT JSON_VALUE(t.a,'$.a')", true},
	{"SELECT JSON_VALUE(CASE WHEN 1 THEN a END,'$.a')", true},
	{"SELECT JSON_VALUE(a->>'$.b','$.a')", true},
	{"SELECT JSON_VALUE(a||b,'$.a')", false},
	{"SELECT JSON_VALUE(a,'')", true},
	{"SELECT JSON_VALUE(a,NULL)", false},
	{"SELECT JSON_VALUE(a,@p)", false},
	{"SELECT JSON_VALUE(a,?)", false},
	// no type.
	{"CREATE TABLE tl (a key)", false},
	{"CREATE TABLE tl (a primary key)", false},
	{"CREATE TABLE tl (a unique)", false},
	{"CREATE TABLE tl (a unique key)", false},
	{"CREATE TABLE tl (a comment 'x')", false},
	{"CREATE TABLE tl (a not null)", false},
	{"CREATE TABLE tl (a null)", false},
	{"CREATE TABLE tl (a default 1)", false},
	{"CREATE TABLE tl (a auto_increment)", false},
	{"CREATE TABLE tl (a collate utf8mb4_bin)", false},
	{"CREATE TABLE tl (a references r (id))", false},
	{"CREATE TABLE tl (a check (a > 0))", false},
	{"CREATE TABLE tl (a generated always as (1))", false},
	{"CREATE TABLE tl (a generated always as (1) stored)", false},
	{"CREATE TABLE tl (a generated always as (1) virtual)", false},
	{"CREATE TABLE tl (a as (1))", false},
	{"CREATE TABLE tl (a as (1) stored)", false},
	{"CREATE TABLE tl (a on update now())", false},
	{"CREATE TABLE tl (a column_format fixed)", false},
	{"CREATE TABLE tl (a storage disk)", false},
	{"CREATE TABLE tl (a srid 0)", false},
	{"CREATE TABLE tl (a invisible)", false},
	{"CREATE TABLE tl (a visible)", false},
	{"CREATE TABLE tl (a auto_random)", false},
	{"CREATE TABLE tl (a)", false},
	{"CREATE TABLE tl (a, b int)", false},
	{"CREATE TABLE tl (a) SELECT 1 AS a", false},
	{"CREATE TABLE tl (a) AS SELECT 1 AS a", false},
	{"CREATE TABLE tl (a, b) AS SELECT 1 AS a, 2 AS b", false},
	{"CREATE TABLE tl (a key) SELECT 1 AS a", false},
	{"CREATE TABLE tl (a default 1) SELECT 1 AS a", false},
	{"CREATE TABLE tl (a not null) AS SELECT 1 AS a", false},
	{"CREATE TABLE tl (a comment 'x') AS SELECT 1 AS a", false},
	{"CREATE TEMPORARY TABLE tl (a key)", false},
	{"ALTER TABLE t ADD a key", false},
	{"ALTER TABLE t ADD COLUMN a comment 'x'", false},
	{"ALTER TABLE t ADD COLUMN a generated always as (1)", false},
	{"ALTER TABLE t ADD (a key)", false},
	{"ALTER TABLE t ADD COLUMN (a not null, b int)", false},
	{"ALTER TABLE t MODIFY a comment 'x'", false},
	{"ALTER TABLE t MODIFY COLUMN a not null", false},
	{"ALTER TABLE t MODIFY a generated always as (1)", false},
	{"ALTER TABLE t CHANGE a b key", false},
	{"ALTER TABLE t CHANGE COLUMN a b comment 'x'", false},
	{"ALTER TABLE t CHANGE a b generated always as (1)", false},
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
	"(1,2,3) OVER ()", "(1) FROM t", "(SELECT 1)", "(YEAR FROM 1)", "(DAY, 1, 2)", "(1, 'a')", "(a)", "(a, b)", "(a) OVER ()",
}

// functionArmsNotCovered is what the list cannot reach: `values` is read by an arm of its
// own for ON DUPLICATE KEY UPDATE, and the engine refuses it everywhere else.
var functionArmsNotCovered = map[string]bool{"values": true}

// judgeKeywordFunctionNames sweeps every keyword the lexer has against the engine.
func judgeKeywordFunctionNames(t *testing.T, accepts func(sql string) bool) {
	t.Helper()
	for word := range keywords {
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

func TestRequiredFormsOceanBaseOracle(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, stmt := range []string{
		"DROP DATABASE IF EXISTS obreq",
		"CREATE DATABASE obreq",
		"USE obreq",
		"CREATE TABLE t (a INT, b INT)",
	} {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("seed %q: %v", stmt, err)
		}
	}
	t.Cleanup(func() { _, _ = db.ExecContext(ctx, "DROP DATABASE IF EXISTS obreq") })
	for _, c := range requiredFormCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := engineAccepts(ctx, db, c.sql)
			if err != nil {
				t.Fatalf("the oracle gave no verdict on %q: %v", c.sql, err)
			}
			if got != c.accept {
				t.Errorf("OceanBase accepts %q = %v, the record says %v", c.sql, got, c.accept)
			}
		})
	}
}

func TestKeywordFunctionNamesOceanBaseOracle(t *testing.T) {
	db := oceanBaseOracle(t)
	ctx := context.Background()
	for _, stmt := range []string{"CREATE DATABASE IF NOT EXISTS obfn", "USE obfn"} {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("seed %q: %v", stmt, err)
		}
	}
	t.Cleanup(func() { _, _ = db.ExecContext(ctx, "DROP DATABASE IF EXISTS obfn") })
	judgeKeywordFunctionNames(t, func(sql string) bool {
		return prepareAcceptsOceanBase(t, ctx, db, sql)
	})
}
