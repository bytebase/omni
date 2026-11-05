package parser

import "testing"

// The clauses a plain index takes on MariaDB and MySQL's grammar does not have: [WAIT n |
// NOWAIT] on CREATE INDEX and DROP INDEX, and IF NOT EXISTS after ALTER TABLE ADD. And the
// ones the shared index_option list read for every index where MariaDB reads them for
// some: WITH PARSER on FULLTEXT alone, and no index type on FULLTEXT or SPATIAL. Every
// verdict was read from mariadb:11.8.8 on 2026-11-05.

type indexClauseGapCase struct {
	sql    string
	accept bool
}

var indexClauseGapCases = []indexClauseGapCase{
	// WAIT and NOWAIT on CREATE INDEX.
	{"CREATE INDEX w1 ON t (a) WAIT 3", true},
	{"CREATE INDEX w2 ON t (a) NOWAIT", true},
	{"CREATE INDEX w3 ON t (a) WAIT 0", true},
	{"CREATE INDEX w4 ON t (a) WAIT 1.5", true},
	{"CREATE INDEX w5 ON t (a) WAIT 1e3", true},
	{"CREATE INDEX w6 ON t (a) WAIT 0x10", true},
	{"CREATE INDEX w7 ON t (a) WAIT '3'", false},
	{"CREATE INDEX w8 ON t (a) WAIT a", false},
	{"CREATE INDEX w9 ON t (a) WAIT", false},
	{"CREATE INDEX w10 ON t (a) WAIT -1", false},
	{"CREATE INDEX w11 ON t (a) WAIT 99999999999999999999", true},
	{"CREATE INDEX w12 ON t (a) WAIT 3 WAIT 4", false},
	{"CREATE INDEX w13 ON t (a) WAIT 3 NOWAIT", false},
	{"CREATE INDEX w14 ON t (a) NOWAIT WAIT 3", false},
	{"CREATE INDEX w15 ON t (a) NOWAIT 3", false},
	{"CREATE INDEX w16 ON t (a) WAIT 3 COMMENT 'c'", true},
	{"CREATE INDEX w17 ON t (a) COMMENT 'c' WAIT 3", false},
	{"CREATE INDEX w18 ON t (a) NOWAIT USING BTREE", true},
	{"CREATE INDEX w19 ON t (a) USING BTREE NOWAIT", false},
	{"CREATE INDEX w20 ON t (a) WAIT 3 ALGORITHM=DEFAULT", true},
	{"CREATE INDEX w21 ON t (a) WAIT 3 LOCK=NONE", true},
	{"CREATE INDEX w22 ON t (a) ALGORITHM=DEFAULT WAIT 3", false},
	{"CREATE INDEX w23 WAIT 3 ON t (a)", false},
	{"CREATE INDEX w24 NOWAIT ON t (a)", false},
	{"CREATE INDEX w25 ON WAIT 3 t (a)", false},
	{"CREATE INDEX w26 ON t WAIT 3 (a)", false},
	{"CREATE INDEX w27 USING BTREE ON t (a) WAIT 3", true},
	{"CREATE UNIQUE INDEX w28 ON t (a) WAIT 3", true},
	{"CREATE UNIQUE INDEX w29 ON t (a) NOWAIT", true},
	{"CREATE FULLTEXT INDEX w30 ON t (c) WAIT 3", true},
	{"CREATE FULLTEXT INDEX w31 ON t (c) NOWAIT", true},
	{"CREATE SPATIAL INDEX w32 ON t (g) WAIT 3", true},
	{"CREATE SPATIAL INDEX w33 ON t (g) NOWAIT", true},
	{"CREATE OR REPLACE INDEX w34 ON t (a) WAIT 3", true},
	{"CREATE INDEX IF NOT EXISTS w35 ON t (a) NOWAIT", true},
	{"CREATE VECTOR INDEX w36 ON t (a) WAIT 1.5", true},
	{"CREATE VECTOR INDEX w37 ON t (a) WAIT 1e3", true},
	{"CREATE VECTOR INDEX w38 ON t (a) WAIT 3 NOWAIT", false},
	{"CREATE INDEX wait ON t (a)", true},
	{"CREATE INDEX nowait ON t (a)", true},
	// IF NOT EXISTS after ADD.
	{"ALTER TABLE t ADD INDEX IF NOT EXISTS n1 (a)", true},
	{"ALTER TABLE t ADD KEY IF NOT EXISTS n2 (a)", true},
	{"ALTER TABLE t ADD INDEX IF NOT EXISTS (a)", true},
	{"ALTER TABLE t ADD INDEX IF NOT EXISTS USING BTREE (a)", true},
	{"ALTER TABLE t ADD INDEX IF NOT EXISTS n3 USING BTREE (a)", true},
	{"ALTER TABLE t ADD INDEX IF NOT n4 (a)", false},
	{"ALTER TABLE t ADD INDEX IF EXISTS n5 (a)", false},
	{"ALTER TABLE t ADD INDEX n6 IF NOT EXISTS (a)", false},
	{"ALTER TABLE t ADD FULLTEXT IF NOT EXISTS n7 (c)", true},
	{"ALTER TABLE t ADD FULLTEXT INDEX IF NOT EXISTS n8 (c)", true},
	{"ALTER TABLE t ADD FULLTEXT KEY IF NOT EXISTS n9 (c)", true},
	{"ALTER TABLE t ADD SPATIAL IF NOT EXISTS n10 (g)", true},
	{"ALTER TABLE t ADD SPATIAL INDEX IF NOT EXISTS n11 (g)", true},
	{"ALTER TABLE t ADD UNIQUE IF NOT EXISTS n12 (a)", true},
	{"ALTER TABLE t ADD UNIQUE INDEX IF NOT EXISTS n13 (a)", true},
	{"ALTER TABLE t ADD UNIQUE KEY IF NOT EXISTS n14 (a)", true},
	{"ALTER TABLE t ADD CONSTRAINT c1 UNIQUE IF NOT EXISTS n15 (a)", true},
	{"ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS c2 UNIQUE (a)", false},
	{"ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS UNIQUE (a)", false},
	{"ALTER TABLE t ADD PRIMARY KEY IF NOT EXISTS (a)", true},
	{"ALTER TABLE t ADD PRIMARY KEY IF NOT EXISTS n16 (a)", true},
	{"ALTER TABLE t ADD CONSTRAINT c3 PRIMARY KEY IF NOT EXISTS (a)", true},
	{"ALTER TABLE t ADD FOREIGN KEY IF NOT EXISTS n17 (a) REFERENCES t (id)", true},
	{"ALTER TABLE t ADD FOREIGN KEY IF NOT EXISTS (a) REFERENCES t (id)", true},
	{"ALTER TABLE t ADD CONSTRAINT c4 FOREIGN KEY IF NOT EXISTS n18 (a) REFERENCES t (id)", true},
	{"ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS c5 FOREIGN KEY (a) REFERENCES t (id)", false},
	{"ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS c6 CHECK (a > 0)", true},
	{"ALTER TABLE t ADD CHECK IF NOT EXISTS (a > 0)", false},
	{"ALTER TABLE t ADD (INDEX IF NOT EXISTS n19 (a))", true},
	{"ALTER TABLE t ADD (KEY IF NOT EXISTS n20 (a), UNIQUE IF NOT EXISTS n21 (a))", true},
	{"ALTER TABLE t ADD (d9 INT, INDEX IF NOT EXISTS n22 (a))", true},
	{"ALTER TABLE t ADD COLUMN (d8 INT, INDEX IF NOT EXISTS n23 (a))", true},
	{"ALTER TABLE t ADD (FULLTEXT IF NOT EXISTS n24 (c))", true},
	{"ALTER TABLE t ADD (PRIMARY KEY IF NOT EXISTS (a))", true},
	{"ALTER TABLE t ADD (CONSTRAINT IF NOT EXISTS c7 UNIQUE (a))", false},
	{"ALTER TABLE t ADD INDEX IF NOT EXISTS n25 (a), ADD INDEX IF NOT EXISTS n26 (b)", true},
	{"ALTER TABLE t ADD INDEX IF NOT EXISTS n27 (a) COMMENT 'c' USING BTREE", true},
	{"CREATE TABLE u1 (a INT, INDEX IF NOT EXISTS i (a))", false},
	{"CREATE TABLE u2 (a INT, KEY IF NOT EXISTS i (a))", false},
	{"CREATE TABLE u3 (a INT, UNIQUE IF NOT EXISTS i (a))", false},
	{"CREATE TABLE u4 (a INT, CONSTRAINT IF NOT EXISTS c UNIQUE (a))", false},
	{"CREATE TABLE u5 (a INT PRIMARY KEY, PRIMARY KEY IF NOT EXISTS (a))", false},
	{"CREATE TABLE u6 (a INT, b INT, FOREIGN KEY IF NOT EXISTS i (a) REFERENCES t (id))", false},
	{"CREATE TABLE u7 (c TEXT, FULLTEXT IF NOT EXISTS i (c))", false},
	// WITH PARSER.
	{"CREATE INDEX x1 ON t (a) WITH PARSER p", false},
	{"CREATE UNIQUE INDEX x2 ON t (a) WITH PARSER p", false},
	{"CREATE FULLTEXT INDEX x3 ON t (c) WITH PARSER p", true},
	{"CREATE SPATIAL INDEX x4 ON t (g) WITH PARSER p", false},
	{"CREATE TABLE y1 (a INT, KEY k (a) WITH PARSER p)", false},
	{"CREATE TABLE y2 (a INT, UNIQUE k (a) WITH PARSER p)", false},
	{"CREATE TABLE y3 (a INT PRIMARY KEY, c TEXT, FULLTEXT k (c) WITH PARSER p)", true},
	{"CREATE TABLE y4 (g GEOMETRY NOT NULL, SPATIAL k (g) WITH PARSER p)", false},
	{"CREATE TABLE y5 (a INT, PRIMARY KEY (a) WITH PARSER p)", false},
	{"ALTER TABLE t ADD INDEX x5 (a) WITH PARSER p", false},
	{"ALTER TABLE t ADD UNIQUE x6 (a) WITH PARSER p", false},
	{"ALTER TABLE t ADD FULLTEXT x7 (c) WITH PARSER p", true},
	{"ALTER TABLE t ADD SPATIAL x8 (g) WITH PARSER p", false},
	{"ALTER TABLE t ADD PRIMARY KEY (a) WITH PARSER p", false},
	{"CREATE FULLTEXT INDEX x9 ON t (c) WITH PARSER", false},
	{"CREATE FULLTEXT INDEX x10 ON t (c) WITH p", false},
	{"CREATE FULLTEXT INDEX x11 ON t (c) COMMENT 'c' WITH PARSER p", true},
	{"CREATE FULLTEXT INDEX x12 ON t (c) WITH PARSER p WITH PARSER q", true},
	{"CREATE FULLTEXT INDEX x13 ON t (c) WITH PARSER 'p'", false},
	{"CREATE FULLTEXT INDEX x14 ON t (c) WITH PARSER `p`", true},
	{"CREATE INDEX x15 ON t (a) COMMENT 'c' WITH PARSER p", false},
	{"CREATE FULLTEXT INDEX x16 ON t (c) USING BTREE", false},
	{"CREATE TABLE y6 (c TEXT, FULLTEXT k (c) USING BTREE)", false},
	{"CREATE SPATIAL INDEX x17 ON t (g) USING BTREE", false},
	{"CREATE TABLE y7 (g GEOMETRY NOT NULL, SPATIAL k (g) USING BTREE)", false},
	// type names.
	{"CREATE TABLE z4 (a foo())", false},
	{"CREATE TABLE z17 (a 'foo')", false},
	{"CREATE TABLE z18 (a foo(a))", false},
	{"CREATE TABLE z19 (a foo('x'))", false},
	{"CREATE TABLE z20 (a `INET6`)", true},
	{"CREATE TABLE z21 (a `UUID`)", true},
	{"CREATE TABLE z22 (a `INET6`(3))", true},
	{"CREATE TABLE z23 (a inet6(3))", true},
	{"CREATE TABLE z24 (a uuid(3))", true},
	{"SELECT * FROM JSON_TABLE('[]', '$[*]' COLUMNS (a foo PATH '$')) AS j", false},
	{"SELECT CAST('x' AS INET6)", true},
	{"SELECT CAST('x' AS inet6(3))", false},
	{"SELECT CAST('x' AS uuid(3))", false},
	{"SELECT CAST('x' AS `INET6`(3))", false},
	{"SELECT CONVERT('x', uuid(3))", false},
	{"CREATE FUNCTION f3(a inet6(3)) RETURNS uuid(3) RETURN a", true},
	{"SELECT * FROM JSON_TABLE('[]', '$[*]' COLUMNS (a inet6 PATH '$')) AS j", false},
	{"SELECT * FROM JSON_TABLE('[]', '$[*]' COLUMNS (a inet6(3) PATH '$')) AS j", false},
	{"SELECT * FROM JSON_TABLE('[]', '$[*]' COLUMNS (a uuid PATH '$')) AS j", false},
	{"ALTER TABLE t ADD PRIMARY IF NOT EXISTS (a)", false},
	{"ALTER TABLE t ADD PRIMARY (a)", false},
	{"CREATE TABLE p1 (a INT, PRIMARY (a))", false},
	{"ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS `c9` CHECK (a > 0)", true},
	{"ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS c10 CHECK (a > 0) NOT ENFORCED", false},
	{"ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS c14 (a > 0)", false},
	{"ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS c11", false},
	{"ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS c12 PRIMARY KEY (a)", false},
	{"ALTER TABLE t ADD CONSTRAINT IF c13 CHECK (a > 0)", false},
	{"ALTER TABLE t ADD INDEX IF NOT EXISTS select (a)", false},
	{"ALTER TABLE t ADD INDEX IF NOT EXISTS `select` (a)", true},
	{"ALTER TABLE t ADD INDEX IF NOT EXISTS (a) COMMENT 'c'", true},
	{"ALTER TABLE t ADD UNIQUE IF NOT EXISTS USING BTREE (a)", true},
	{"ALTER TABLE t ADD UNIQUE IF NOT EXISTS n40 USING BTREE (a)", true},
	{"ALTER TABLE t ADD PRIMARY KEY IF NOT EXISTS USING BTREE (a)", true},
	{"CREATE INDEX n41 ON t (a) WAIT 0X10", false},
	{"CREATE INDEX n42 ON t (a) WAIT X'10'", false},
	{"CREATE INDEX n43 ON t (a) WAIT b'1'", false},
	{"CREATE INDEX n44 ON t (a) WAIT 0b10", false},
	{"CREATE INDEX n45 ON t (a) WAIT (5)", false},
	{"CREATE INDEX n46 ON t (a) WAIT @v", false},
	{"DROP INDEX ki ON t WAIT 3 ALGORITHM=DEFAULT", false},
	{"DROP INDEX ki ON t WAIT 3 NOWAIT", false},
	{"DROP INDEX ki ON t COMMENT 'c'", false},
	// more syntax.
	{"DROP INDEX ki ON t WAIT 3", true},
	{"DROP INDEX ki ON t NOWAIT", true},
	{"DROP INDEX IF EXISTS ki ON t WAIT 1.5", true},
	{"DROP INDEX ki ON t WAIT", false},
	{"DROP INDEX ki WAIT 3 ON t", false},
	{"ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS CHECK (a > 0)", false},
	{"ALTER TABLE t ADD (CONSTRAINT IF NOT EXISTS c8 CHECK (a > 0))", false},
	{"ALTER TABLE t ADD (CONSTRAINT IF NOT EXISTS CHECK (a > 0))", false},
	{"ALTER TABLE t ADD (FOREIGN KEY IF NOT EXISTS n30 (a) REFERENCES t (id))", true},
	{"ALTER TABLE t ADD (SPATIAL IF NOT EXISTS n31 (g))", true},
	{"CREATE TABLE u8 (a INT, CONSTRAINT IF NOT EXISTS c CHECK (a > 0))", false},
	{"CREATE TABLE u9 (a INT, CONSTRAINT IF NOT EXISTS CHECK (a > 0))", false},
	{"CREATE FULLTEXT INDEX r1 USING BTREE ON t (c)", false},
	{"CREATE FULLTEXT INDEX r2 TYPE BTREE ON t (c)", false},
	{"CREATE SPATIAL INDEX r3 USING RTREE ON t (g)", false},
	{"CREATE SPATIAL INDEX r4 USING BTREE ON t (g)", false},
	{"CREATE FULLTEXT INDEX r5 ON t (c) TYPE BTREE", false},
	{"CREATE SPATIAL INDEX r6 ON t (g) TYPE RTREE", false},
	{"CREATE TABLE r7 (c TEXT, FULLTEXT k USING BTREE (c))", false},
	{"CREATE TABLE r8 (g GEOMETRY NOT NULL, SPATIAL k USING RTREE (g))", false},
	{"ALTER TABLE t ADD FULLTEXT r9 (c) USING BTREE", false},
	{"ALTER TABLE t ADD SPATIAL r10 (g) USING RTREE", false},
	{"ALTER TABLE t ADD FULLTEXT r11 USING BTREE (c)", false},
	{"ALTER TABLE t ADD SPATIAL r12 USING RTREE (g)", false},
	{"CREATE TABLE z30 (a uuid(3,2))", true},
	{"CREATE TABLE z31 (a inet6(1.5))", true},
	{"CREATE TABLE z32 (a inet4(3))", true},
	{"CREATE TABLE z33 (a inet6())", false},
	{"CREATE TABLE z34 (a inet6(a))", false},
	{"CREATE TABLE z35 (a uuid(1e5))", false},
}

func TestIndexClauseGapsOracle(t *testing.T) {
	o := startMariaDB(t)
	for _, tc := range indexClauseGapCases {
		t.Run(tc.sql, func(t *testing.T) {
			if got := o.prepareParses(t, tc.sql); got != tc.accept {
				t.Fatalf("MariaDB 11.8.8 accepts=%v, want %v for %q", got, tc.accept, tc.sql)
			}
			omniAccepts, err := parsesAsOneStatement(tc.sql)
			if omniAccepts != tc.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, MariaDB 11.8.8 accepts=%v for %q: %v",
					omniAccepts, tc.accept, tc.sql, err)
			}
		})
	}
}

// TestIndexClauseGaps holds the parser to the recorded verdicts without the server.
func TestIndexClauseGaps(t *testing.T) {
	for _, tc := range indexClauseGapCases {
		t.Run(tc.sql, func(t *testing.T) {
			omniAccepts, err := parsesAsOneStatement(tc.sql)
			if omniAccepts != tc.accept {
				t.Errorf("omni accepts=%v, the record says %v for %q: %v", omniAccepts, tc.accept, tc.sql, err)
			}
		})
	}
}

// unknownDataTypeCases name no type MariaDB has. The grammar reads any identifier in type
// position, even a reserved word, and its action then answers 4161, so none of these can
// ever run. The parser refuses them itself, as it does the 1221 of OR REPLACE with IF NOT
// EXISTS, rather than model a column of a type that does not exist.
var unknownDataTypeCases = []parseCodeCase{
	{"CREATE TABLE z1 (a foo)", 4161},
	{"CREATE TABLE z2 (a foo(3))", 4161},
	{"CREATE TABLE z3 (a foo(3,2))", 4161},
	{"CREATE TABLE z5 (a `INT`)", 4161},
	{"CREATE TABLE z6 (a `foo`)", 4161},
	{"CREATE TABLE z7 (a `foo`(3))", 4161},
	{"CREATE TABLE z8 (a foo UNSIGNED)", 4161},
	{"CREATE TABLE z9 (a foo NOT NULL)", 4161},
	{"CREATE TABLE z10 (a foo DEFAULT 1)", 4161},
	{"CREATE TABLE z11 (a foo CHARACTER SET utf8mb4)", 4161},
	{"CREATE TABLE z12 (a foo COLLATE utf8mb4_bin)", 4161},
	{"CREATE TABLE z13 (a foo BINARY)", 4161},
	{"CREATE TABLE z14 (a foo.bar)", 4161},
	{"CREATE TABLE z15 (a status)", 4161},
	{"CREATE TABLE z16 (a select)", 4161},
	{"ALTER TABLE t ADD COLUMN z25 foo", 4161},
	{"ALTER TABLE t MODIFY a foo", 4161},
	{"SELECT CAST('x' AS foo)", 4161},
	{"SELECT CAST('x' AS foo(3))", 4161},
	{"SELECT CAST('x' AS `foo`)", 4161},
	{"SELECT CAST('x' AS `INT`)", 4161},
	{"SELECT CONVERT('x', foo)", 4161},
	{"CREATE FUNCTION f1(a foo) RETURNS INT RETURN 1", 4161},
	{"CREATE FUNCTION f2() RETURNS foo RETURN 1", 4161},
	{"CREATE PROCEDURE p1() BEGIN DECLARE b foo; END", 4161},
}

func TestUnknownDataTypesOracle(t *testing.T) {
	assertParseCodes(t, unknownDataTypeCases)
}
