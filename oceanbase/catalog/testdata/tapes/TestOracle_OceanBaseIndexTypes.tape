-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseIndexTypes default#1
> exec
| DROP DATABASE IF EXISTS obindextypes
< ok 0 0
> exec
| CREATE DATABASE obindextypes
< ok 1 0
> exec
| USE obindextypes
< ok 0 0
> exec
| CREATE TABLE t1 (a INT, KEY k (a))
< ok 0 0
> query
| SELECT DISTINCT INDEX_TYPE FROM information_schema.STATISTICS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND INDEX_NAME = ?
: string
| obindextypes
: string
| t1
: string
| k
< result
< column "INDEX_TYPE" "VARCHAR"
< row
: bytes
| BTREE
> exec
| CREATE TABLE t2 (a INT, KEY k (a) USING BTREE)
< ok 0 0
> query
| SELECT DISTINCT INDEX_TYPE FROM information_schema.STATISTICS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND INDEX_NAME = ?
: string
| obindextypes
: string
| t2
: string
| k
< result
< column "INDEX_TYPE" "VARCHAR"
< row
: bytes
| BTREE
> exec
| CREATE TABLE t3 (a INT, KEY k (a) USING HASH)
< ok 0 0
> query
| SELECT DISTINCT INDEX_TYPE FROM information_schema.STATISTICS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND INDEX_NAME = ?
: string
| obindextypes
: string
| t3
: string
| k
< result
< column "INDEX_TYPE" "VARCHAR"
< row
: bytes
| HASH
> exec
| CREATE TABLE t4 (a INT, KEY k USING HASH (a))
< ok 0 0
> query
| SELECT DISTINCT INDEX_TYPE FROM information_schema.STATISTICS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND INDEX_NAME = ?
: string
| obindextypes
: string
| t4
: string
| k
< result
< column "INDEX_TYPE" "VARCHAR"
< row
: bytes
| HASH
> exec
| CREATE TABLE t5 (a INT, UNIQUE KEY k (a) USING HASH) ENGINE=InnoDB
< ok 0 0
> query
| SELECT DISTINCT INDEX_TYPE FROM information_schema.STATISTICS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND INDEX_NAME = ?
: string
| obindextypes
: string
| t5
: string
| k
< result
< column "INDEX_TYPE" "VARCHAR"
< row
: bytes
| HASH
> exec
| CREATE TABLE t6 (a INT)
< ok 0 0
> exec
| CREATE INDEX k ON t6 (a) USING HASH
< ok 0 0
> query
| SELECT DISTINCT INDEX_TYPE FROM information_schema.STATISTICS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND INDEX_NAME = ?
: string
| obindextypes
: string
| t6
: string
| k
< result
< column "INDEX_TYPE" "VARCHAR"
< row
: bytes
| HASH
> exec
| CREATE TABLE t7 (a INT)
< ok 0 0
> exec
| ALTER TABLE t7 ADD KEY k (a) USING HASH
< ok 0 0
> query
| SELECT DISTINCT INDEX_TYPE FROM information_schema.STATISTICS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND INDEX_NAME = ?
: string
| obindextypes
: string
| t7
: string
| k
< result
< column "INDEX_TYPE" "VARCHAR"
< row
: bytes
| HASH
> exec
| CREATE TABLE t8 (g GEOMETRY NOT NULL SRID 0, SPATIAL KEY k (g))
< ok 0 0
> query
| SELECT DISTINCT INDEX_TYPE FROM information_schema.STATISTICS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND INDEX_NAME = ?
: string
| obindextypes
: string
| t8
: string
| k
< result
< column "INDEX_TYPE" "VARCHAR"
< row
: bytes
| BTREE
> exec
| CREATE TABLE t9 (g GEOMETRY NOT NULL SRID 0, SPATIAL KEY k (g) USING BTREE)
< ok 0 0
> query
| SELECT DISTINCT INDEX_TYPE FROM information_schema.STATISTICS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND INDEX_NAME = ?
: string
| obindextypes
: string
| t9
: string
| k
< result
< column "INDEX_TYPE" "VARCHAR"
< row
: bytes
| BTREE
> exec
| CREATE TABLE t11 (g GEOMETRY NOT NULL SRID 0)
< ok 0 0
> exec
| ALTER TABLE t11 ADD SPATIAL KEY k (g) USING BTREE
< ok 0 0
> query
| SELECT DISTINCT INDEX_TYPE FROM information_schema.STATISTICS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND INDEX_NAME = ?
: string
| obindextypes
: string
| t11
: string
| k
< result
< column "INDEX_TYPE" "VARCHAR"
< row
: bytes
| BTREE
