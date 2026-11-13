-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestIndexClauseGapsOracle default#1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| CREATE INDEX w1 ON t (a) WAIT 3
< ok 0 0
> exec ^
: string
| CREATE INDEX w2 ON t (a) NOWAIT
< ok 0 0
> exec ^
: string
| CREATE INDEX w3 ON t (a) WAIT 0
< ok 0 0
> exec ^
: string
| CREATE INDEX w4 ON t (a) WAIT 1.5
< ok 0 0
> exec ^
: string
| CREATE INDEX w5 ON t (a) WAIT 1e3
< ok 0 0
> exec ^
: string
| CREATE INDEX w6 ON t (a) WAIT 0x10
< ok 0 0
> exec ^
: string
| CREATE INDEX w7 ON t (a) WAIT '3'
< error
| 1064 42000 …near ''3'' at line 1
> exec ^
: string
| CREATE INDEX w8 ON t (a) WAIT a
< error
| 1064 42000 …near 'a' at line 1
> exec ^
: string
| CREATE INDEX w9 ON t (a) WAIT
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE INDEX w10 ON t (a) WAIT -1
< error
| 1064 42000 …near '-1' at line 1
> exec ^
: string
| CREATE INDEX w11 ON t (a) WAIT 99999999999999999999
< ok 0 0
> exec ^
: string
| CREATE INDEX w12 ON t (a) WAIT 3 WAIT 4
< error
| 1064 42000 …near 'WAIT 4' at line 1
> exec ^
: string
| CREATE INDEX w13 ON t (a) WAIT 3 NOWAIT
< error
| 1064 42000 …near 'NOWAIT' at line 1
> exec ^
: string
| CREATE INDEX w14 ON t (a) NOWAIT WAIT 3
< error
| 1064 42000 …near 'WAIT 3' at line 1
> exec ^
: string
| CREATE INDEX w15 ON t (a) NOWAIT 3
< error
| 1064 42000 …near '3' at line 1
> exec ^
: string
| CREATE INDEX w16 ON t (a) WAIT 3 COMMENT 'c'
< ok 0 0
> exec ^
: string
| CREATE INDEX w17 ON t (a) COMMENT 'c' WAIT 3
< error
| 1064 42000 …near 'WAIT 3' at line 1
> exec ^
: string
| CREATE INDEX w18 ON t (a) NOWAIT USING BTREE
< ok 0 0
> exec ^
: string
| CREATE INDEX w19 ON t (a) USING BTREE NOWAIT
< error
| 1064 42000 …near 'NOWAIT' at line 1
> exec ^
: string
| CREATE INDEX w20 ON t (a) WAIT 3 ALGORITHM=DEFAULT
< ok 0 0
> exec ^
: string
| CREATE INDEX w21 ON t (a) WAIT 3 LOCK=NONE
< ok 0 0
> exec ^
: string
| CREATE INDEX w22 ON t (a) ALGORITHM=DEFAULT WAIT 3
< error
| 1064 42000 …near 'WAIT 3' at line 1
> exec ^
: string
| CREATE INDEX w23 WAIT 3 ON t (a)
< error
| 1064 42000 …near 'WAIT 3 ON t (a)' at line 1
> exec ^
: string
| CREATE INDEX w24 NOWAIT ON t (a)
< error
| 1064 42000 …near 'NOWAIT ON t (a)' at line 1
> exec ^
: string
| CREATE INDEX w25 ON WAIT 3 t (a)
< error
| 1064 42000 …near '3 t (a)' at line 1
> exec ^
: string
| CREATE INDEX w26 ON t WAIT 3 (a)
< error
| 1064 42000 …near 'WAIT 3 (a)' at line 1
> exec ^
: string
| CREATE INDEX w27 USING BTREE ON t (a) WAIT 3
< ok 0 0
> exec ^
: string
| CREATE UNIQUE INDEX w28 ON t (a) WAIT 3
< ok 0 0
> exec ^
: string
| CREATE UNIQUE INDEX w29 ON t (a) NOWAIT
< ok 0 0
> exec ^
: string
| CREATE FULLTEXT INDEX w30 ON t (c) WAIT 3
< ok 0 0
> exec ^
: string
| CREATE FULLTEXT INDEX w31 ON t (c) NOWAIT
< ok 0 0
> exec ^
: string
| CREATE SPATIAL INDEX w32 ON t (g) WAIT 3
< ok 0 0
> exec ^
: string
| CREATE SPATIAL INDEX w33 ON t (g) NOWAIT
< ok 0 0
> exec ^
: string
| CREATE OR REPLACE INDEX w34 ON t (a) WAIT 3
< ok 0 0
> exec ^
: string
| CREATE INDEX IF NOT EXISTS w35 ON t (a) NOWAIT
< ok 0 0
> exec ^
: string
| CREATE VECTOR INDEX w36 ON t (a) WAIT 1.5
< ok 0 0
> exec ^
: string
| CREATE VECTOR INDEX w37 ON t (a) WAIT 1e3
< ok 0 0
> exec ^
: string
| CREATE VECTOR INDEX w38 ON t (a) WAIT 3 NOWAIT
< error
| 1064 42000 …near 'NOWAIT' at line 1
> exec ^
: string
| CREATE INDEX wait ON t (a)
< ok 0 0
> exec ^
: string
| CREATE INDEX nowait ON t (a)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD INDEX IF NOT EXISTS n1 (a)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD KEY IF NOT EXISTS n2 (a)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD INDEX IF NOT EXISTS (a)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD INDEX IF NOT EXISTS USING BTREE (a)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD INDEX IF NOT EXISTS n3 USING BTREE (a)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD INDEX IF NOT n4 (a)
< error
| 1064 42000 …near 'n4 (a)' at line 1
> exec ^
: string
| ALTER TABLE t ADD INDEX IF EXISTS n5 (a)
< error
| 1064 42000 …near 'EXISTS n5 (a)' at line 1
> exec ^
: string
| ALTER TABLE t ADD INDEX n6 IF NOT EXISTS (a)
< error
| 1064 42000 …near 'IF NOT EXISTS (a)' at line 1
> exec ^
: string
| ALTER TABLE t ADD FULLTEXT IF NOT EXISTS n7 (c)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD FULLTEXT INDEX IF NOT EXISTS n8 (c)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD FULLTEXT KEY IF NOT EXISTS n9 (c)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD SPATIAL IF NOT EXISTS n10 (g)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD SPATIAL INDEX IF NOT EXISTS n11 (g)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD UNIQUE IF NOT EXISTS n12 (a)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD UNIQUE INDEX IF NOT EXISTS n13 (a)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD UNIQUE KEY IF NOT EXISTS n14 (a)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD CONSTRAINT c1 UNIQUE IF NOT EXISTS n15 (a)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS c2 UNIQUE (a)
< error
| 1064 42000 …near 'UNIQUE (a)' at line 1
> exec ^
: string
| ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS UNIQUE (a)
< error
| 1064 42000 …near 'UNIQUE (a)' at line 1
> exec ^
: string
| ALTER TABLE t ADD PRIMARY KEY IF NOT EXISTS (a)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD PRIMARY KEY IF NOT EXISTS n16 (a)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD CONSTRAINT c3 PRIMARY KEY IF NOT EXISTS (a)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD FOREIGN KEY IF NOT EXISTS n17 (a) REFERENCES t (id)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD FOREIGN KEY IF NOT EXISTS (a) REFERENCES t (id)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD CONSTRAINT c4 FOREIGN KEY IF NOT EXISTS n18 (a) REFERENCES t (id)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS c5 FOREIGN KEY (a) REFERENCES t (id)
< error
| 1064 42000 …near 'FOREIGN KEY (a) REFERENCES t (id)' at line 1
> exec ^
: string
| ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS c6 CHECK (a > 0)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD CHECK IF NOT EXISTS (a > 0)
< error
| 1064 42000 …near 'IF NOT EXISTS (a > 0)' at line 1
> exec ^
: string
| ALTER TABLE t ADD (INDEX IF NOT EXISTS n19 (a))
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD (KEY IF NOT EXISTS n20 (a), UNIQUE IF NOT EXISTS n21 (a))
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD (d9 INT, INDEX IF NOT EXISTS n22 (a))
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD COLUMN (d8 INT, INDEX IF NOT EXISTS n23 (a))
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD (FULLTEXT IF NOT EXISTS n24 (c))
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD (PRIMARY KEY IF NOT EXISTS (a))
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD (CONSTRAINT IF NOT EXISTS c7 UNIQUE (a))
< error
| 1064 42000 …near 'IF NOT EXISTS c7 UNIQUE (a))' at line 1
> exec ^
: string
| ALTER TABLE t ADD INDEX IF NOT EXISTS n25 (a), ADD INDEX IF NOT EXISTS n26 (b)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD INDEX IF NOT EXISTS n27 (a) COMMENT 'c' USING BTREE
< ok 0 0
> exec ^
: string
| CREATE TABLE u1 (a INT, INDEX IF NOT EXISTS i (a))
< error
| 1064 42000 …near '(a))' at line 1
> exec ^
: string
| CREATE TABLE u2 (a INT, KEY IF NOT EXISTS i (a))
< error
| 1064 42000 …near '(a))' at line 1
> exec ^
: string
| CREATE TABLE u3 (a INT, UNIQUE IF NOT EXISTS i (a))
< error
| 1064 42000 …near '(a))' at line 1
> exec ^
: string
| CREATE TABLE u4 (a INT, CONSTRAINT IF NOT EXISTS c UNIQUE (a))
< error
| 1064 42000 …near 'IF NOT EXISTS c UNIQUE (a))' at line 1
> exec ^
: string
| CREATE TABLE u5 (a INT PRIMARY KEY, PRIMARY KEY IF NOT EXISTS (a))
< error
| 1064 42000 …near '(a))' at line 1
> exec ^
: string
| CREATE TABLE u6 (a INT, b INT, FOREIGN KEY IF NOT EXISTS i (a) REFERENCES t (id))
< error
| 1064 42000 …near '(a) REFERENCES t (id))' at line 1
> exec ^
: string
| CREATE TABLE u7 (c TEXT, FULLTEXT IF NOT EXISTS i (c))
< error
| 1064 42000 …near '(c))' at line 1
> exec ^
: string
| CREATE INDEX x1 ON t (a) WITH PARSER p
< error
| 1064 42000 …near 'WITH PARSER p' at line 1
> exec ^
: string
| CREATE UNIQUE INDEX x2 ON t (a) WITH PARSER p
< error
| 1064 42000 …near 'WITH PARSER p' at line 1
> exec ^
: string
| CREATE FULLTEXT INDEX x3 ON t (c) WITH PARSER p
< error
| 1128 HY000 Function 'p' is not defined
> exec ^
: string
| CREATE SPATIAL INDEX x4 ON t (g) WITH PARSER p
< error
| 1064 42000 …near 'WITH PARSER p' at line 1
> exec ^
: string
| CREATE TABLE y1 (a INT, KEY k (a) WITH PARSER p)
< error
| 1064 42000 …near 'WITH PARSER p)' at line 1
> exec ^
: string
| CREATE TABLE y2 (a INT, UNIQUE k (a) WITH PARSER p)
< error
| 1064 42000 …near 'WITH PARSER p)' at line 1
> exec ^
: string
| CREATE TABLE y3 (a INT PRIMARY KEY, c TEXT, FULLTEXT k (c) WITH PARSER p)
< error
| 1128 HY000 Function 'p' is not defined
> exec ^
: string
| CREATE TABLE y4 (g GEOMETRY NOT NULL, SPATIAL k (g) WITH PARSER p)
< error
| 1064 42000 …near 'WITH PARSER p)' at line 1
> exec ^
: string
| CREATE TABLE y5 (a INT, PRIMARY KEY (a) WITH PARSER p)
< error
| 1064 42000 …near 'WITH PARSER p)' at line 1
> exec ^
: string
| ALTER TABLE t ADD INDEX x5 (a) WITH PARSER p
< error
| 1064 42000 …near 'WITH PARSER p' at line 1
> exec ^
: string
| ALTER TABLE t ADD UNIQUE x6 (a) WITH PARSER p
< error
| 1064 42000 …near 'WITH PARSER p' at line 1
> exec ^
: string
| ALTER TABLE t ADD FULLTEXT x7 (c) WITH PARSER p
< error
| 1128 HY000 Function 'p' is not defined
> exec ^
: string
| ALTER TABLE t ADD SPATIAL x8 (g) WITH PARSER p
< error
| 1064 42000 …near 'WITH PARSER p' at line 1
> exec ^
: string
| ALTER TABLE t ADD PRIMARY KEY (a) WITH PARSER p
< error
| 1064 42000 …near 'WITH PARSER p' at line 1
> exec ^
: string
| CREATE FULLTEXT INDEX x9 ON t (c) WITH PARSER
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| CREATE FULLTEXT INDEX x10 ON t (c) WITH p
< error
| 1064 42000 …near 'p' at line 1
> exec ^
: string
| CREATE FULLTEXT INDEX x11 ON t (c) COMMENT 'c' WITH PARSER p
< error
| 1128 HY000 Function 'p' is not defined
> exec ^
: string
| CREATE FULLTEXT INDEX x12 ON t (c) WITH PARSER p WITH PARSER q
< error
| 1128 HY000 Function 'p' is not defined
> exec ^
: string
| CREATE FULLTEXT INDEX x13 ON t (c) WITH PARSER 'p'
< error
| 1064 42000 …near ''p'' at line 1
> exec ^
: string
| CREATE FULLTEXT INDEX x14 ON t (c) WITH PARSER `p`
< error
| 1128 HY000 Function 'p' is not defined
> exec ^
: string
| CREATE INDEX x15 ON t (a) COMMENT 'c' WITH PARSER p
< error
| 1064 42000 …near 'WITH PARSER p' at line 1
> exec ^
: string
| CREATE FULLTEXT INDEX x16 ON t (c) USING BTREE
< error
| 1064 42000 …near 'USING BTREE' at line 1
> exec ^
: string
| CREATE TABLE y6 (c TEXT, FULLTEXT k (c) USING BTREE)
< error
| 1064 42000 …near 'USING BTREE)' at line 1
> exec ^
: string
| CREATE SPATIAL INDEX x17 ON t (g) USING BTREE
< error
| 1064 42000 …near 'USING BTREE' at line 1
> exec ^
: string
| CREATE TABLE y7 (g GEOMETRY NOT NULL, SPATIAL k (g) USING BTREE)
< error
| 1064 42000 …near 'USING BTREE)' at line 1
> exec ^
: string
| CREATE TABLE z4 (a foo())
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| CREATE TABLE z17 (a 'foo')
< error
| 1064 42000 …near ''foo')' at line 1
> exec ^
: string
| CREATE TABLE z18 (a foo(a))
< error
| 1064 42000 …near 'a))' at line 1
> exec ^
: string
| CREATE TABLE z19 (a foo('x'))
< error
| 1064 42000 …near ''x'))' at line 1
> exec ^
: string
| CREATE TABLE z20 (a `INET6`)
< ok 0 0
> exec ^
: string
| CREATE TABLE z21 (a `UUID`)
< ok 0 0
> exec ^
: string
| CREATE TABLE z22 (a `INET6`(3))
< ok 0 0
> exec ^
: string
| CREATE TABLE z23 (a inet6(3))
< ok 0 0
> exec ^
: string
| CREATE TABLE z24 (a uuid(3))
< ok 0 0
> exec ^
: string
| SELECT * FROM JSON_TABLE('[]', '$[*]' COLUMNS (a foo PATH '$')) AS j
< error
| 1064 42000 …near 'foo PATH '$')) AS j' at line 1
> exec ^
: string
| SELECT CAST('x' AS INET6)
< ok 0 0
> exec ^
: string
| SELECT CAST('x' AS inet6(3))
< error
| 1064 42000 …near '(3))' at line 1
> exec ^
: string
| SELECT CAST('x' AS uuid(3))
< error
| 1064 42000 …near '(3))' at line 1
> exec ^
: string
| SELECT CAST('x' AS `INET6`(3))
< error
| 1064 42000 …near '(3))' at line 1
> exec ^
: string
| SELECT CONVERT('x', uuid(3))
< error
| 1064 42000 …near '(3))' at line 1
> exec ^
: string
| CREATE FUNCTION f3(a inet6(3)) RETURNS uuid(3) RETURN a
< ok 0 0
> exec ^
: string
| SELECT * FROM JSON_TABLE('[]', '$[*]' COLUMNS (a inet6 PATH '$')) AS j
< error
| 1064 42000 …near 'inet6 PATH '$')) AS j' at line 1
> exec ^
: string
| SELECT * FROM JSON_TABLE('[]', '$[*]' COLUMNS (a inet6(3) PATH '$')) AS j
< error
| 1064 42000 …near 'inet6(3) PATH '$')) AS j' at line 1
> exec ^
: string
| SELECT * FROM JSON_TABLE('[]', '$[*]' COLUMNS (a uuid PATH '$')) AS j
< error
| 1064 42000 …near 'uuid PATH '$')) AS j' at line 1
> exec ^
: string
| ALTER TABLE t ADD PRIMARY IF NOT EXISTS (a)
< error
| 1064 42000 …near 'IF NOT EXISTS (a)' at line 1
> exec ^
: string
| ALTER TABLE t ADD PRIMARY (a)
< error
| 1064 42000 …near '(a)' at line 1
> exec ^
: string
| CREATE TABLE p1 (a INT, PRIMARY (a))
< error
| 1064 42000 …near '(a))' at line 1
> exec ^
: string
| ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS `c9` CHECK (a > 0)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS c10 CHECK (a > 0) NOT ENFORCED
< error
| 1064 42000 …near 'NOT ENFORCED' at line 1
> exec ^
: string
| ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS c14 (a > 0)
< error
| 1064 42000 …near '(a > 0)' at line 1
> exec ^
: string
| ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS c11
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS c12 PRIMARY KEY (a)
< error
| 1064 42000 …near 'PRIMARY KEY (a)' at line 1
> exec ^
: string
| ALTER TABLE t ADD CONSTRAINT IF c13 CHECK (a > 0)
< error
| 1064 42000 …near 'c13 CHECK (a > 0)' at line 1
> exec ^
: string
| ALTER TABLE t ADD INDEX IF NOT EXISTS select (a)
< error
| 1064 42000 …near 'select (a)' at line 1
> exec ^
: string
| ALTER TABLE t ADD INDEX IF NOT EXISTS `select` (a)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD INDEX IF NOT EXISTS (a) COMMENT 'c'
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD UNIQUE IF NOT EXISTS USING BTREE (a)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD UNIQUE IF NOT EXISTS n40 USING BTREE (a)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD PRIMARY KEY IF NOT EXISTS USING BTREE (a)
< ok 0 0
> exec ^
: string
| CREATE INDEX n41 ON t (a) WAIT 0X10
< error
| 1064 42000 …near '0X10' at line 1
> exec ^
: string
| CREATE INDEX n42 ON t (a) WAIT X'10'
< error
| 1064 42000 …near 'X'10'' at line 1
> exec ^
: string
| CREATE INDEX n43 ON t (a) WAIT b'1'
< error
| 1064 42000 …near 'b'1'' at line 1
> exec ^
: string
| CREATE INDEX n44 ON t (a) WAIT 0b10
< error
| 1064 42000 …near '0b10' at line 1
> exec ^
: string
| CREATE INDEX n45 ON t (a) WAIT (5)
< error
| 1064 42000 …near '(5)' at line 1
> exec ^
: string
| CREATE INDEX n46 ON t (a) WAIT @v
< error
| 1064 42000 …near '@v' at line 1
> exec ^
: string
| DROP INDEX ki ON t WAIT 3 ALGORITHM=DEFAULT
< error
| 1064 42000 …near 'ALGORITHM=DEFAULT' at line 1
> exec ^
: string
| DROP INDEX ki ON t WAIT 3 NOWAIT
< error
| 1064 42000 …near 'NOWAIT' at line 1
> exec ^
: string
| DROP INDEX ki ON t COMMENT 'c'
< error
| 1064 42000 …near 'COMMENT 'c'' at line 1
> exec ^
: string
| DROP INDEX ki ON t WAIT 3
< ok 0 0
> exec ^
: string
| DROP INDEX ki ON t NOWAIT
< ok 0 0
> exec ^
: string
| DROP INDEX IF EXISTS ki ON t WAIT 1.5
< ok 0 0
> exec ^
: string
| DROP INDEX ki ON t WAIT
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| DROP INDEX ki WAIT 3 ON t
< error
| 1064 42000 …near 'WAIT 3 ON t' at line 1
> exec ^
: string
| ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS CHECK (a > 0)
< error
| 1064 42000 …near 'CHECK (a > 0)' at line 1
> exec ^
: string
| ALTER TABLE t ADD (CONSTRAINT IF NOT EXISTS c8 CHECK (a > 0))
< error
| 1064 42000 …near 'IF NOT EXISTS c8 CHECK (a > 0))' at line 1
> exec ^
: string
| ALTER TABLE t ADD (CONSTRAINT IF NOT EXISTS CHECK (a > 0))
< error
| 1064 42000 …near 'IF NOT EXISTS CHECK (a > 0))' at line 1
> exec ^
: string
| ALTER TABLE t ADD (FOREIGN KEY IF NOT EXISTS n30 (a) REFERENCES t (id))
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD (SPATIAL IF NOT EXISTS n31 (g))
< ok 0 0
> exec ^
: string
| CREATE TABLE u8 (a INT, CONSTRAINT IF NOT EXISTS c CHECK (a > 0))
< error
| 1064 42000 …near 'IF NOT EXISTS c CHECK (a > 0))' at line 1
> exec ^
: string
| CREATE TABLE u9 (a INT, CONSTRAINT IF NOT EXISTS CHECK (a > 0))
< error
| 1064 42000 …near 'IF NOT EXISTS CHECK (a > 0))' at line 1
> exec ^
: string
| CREATE FULLTEXT INDEX r1 USING BTREE ON t (c)
< error
| 1064 42000 …near 'USING BTREE ON t (c)' at line 1
> exec ^
: string
| CREATE FULLTEXT INDEX r2 TYPE BTREE ON t (c)
< error
| 1064 42000 …near 'TYPE BTREE ON t (c)' at line 1
> exec ^
: string
| CREATE SPATIAL INDEX r3 USING RTREE ON t (g)
< error
| 1064 42000 …near 'USING RTREE ON t (g)' at line 1
> exec ^
: string
| CREATE SPATIAL INDEX r4 USING BTREE ON t (g)
< error
| 1064 42000 …near 'USING BTREE ON t (g)' at line 1
> exec ^
: string
| CREATE FULLTEXT INDEX r5 ON t (c) TYPE BTREE
< error
| 1064 42000 …near 'TYPE BTREE' at line 1
> exec ^
: string
| CREATE SPATIAL INDEX r6 ON t (g) TYPE RTREE
< error
| 1064 42000 …near 'TYPE RTREE' at line 1
> exec ^
: string
| CREATE TABLE r7 (c TEXT, FULLTEXT k USING BTREE (c))
< error
| 1064 42000 …near 'USING BTREE (c))' at line 1
> exec ^
: string
| CREATE TABLE r8 (g GEOMETRY NOT NULL, SPATIAL k USING RTREE (g))
< error
| 1064 42000 …near 'USING RTREE (g))' at line 1
> exec ^
: string
| ALTER TABLE t ADD FULLTEXT r9 (c) USING BTREE
< error
| 1064 42000 …near 'USING BTREE' at line 1
> exec ^
: string
| ALTER TABLE t ADD SPATIAL r10 (g) USING RTREE
< error
| 1064 42000 …near 'USING RTREE' at line 1
> exec ^
: string
| ALTER TABLE t ADD FULLTEXT r11 USING BTREE (c)
< error
| 1064 42000 …near 'USING BTREE (c)' at line 1
> exec ^
: string
| ALTER TABLE t ADD SPATIAL r12 USING RTREE (g)
< error
| 1064 42000 …near 'USING RTREE (g)' at line 1
> exec ^
: string
| CREATE TABLE z30 (a uuid(3,2))
< ok 0 0
> exec ^
: string
| CREATE TABLE z31 (a inet6(1.5))
< ok 0 0
> exec ^
: string
| CREATE TABLE z32 (a inet4(3))
< ok 0 0
> exec ^
: string
| CREATE TABLE z33 (a inet6())
< error
| 1064 42000 …near '))' at line 1
> exec ^
: string
| CREATE TABLE z34 (a inet6(a))
< error
| 1064 42000 …near 'a))' at line 1
> exec ^
: string
| CREATE TABLE z35 (a uuid(1e5))
< error
| 1064 42000 …near '1e5))' at line 1
