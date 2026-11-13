-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestScenario_C7 default#1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY (a))
< ok 0 0
> query
| SELECT INDEX_TYPE FROM information_schema.STATISTICS
|              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND INDEX_NAME='a'
< result
< column "INDEX_TYPE" "VARCHAR"
< row
: bytes
| BTREE
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c (a INT, CONSTRAINT c_ibfk_1 FOREIGN KEY (a) REFERENCES p(id))
< ok 0 0
> query
| SELECT INDEX_NAME FROM information_schema.STATISTICS
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='c'
|             ORDER BY INDEX_NAME
< result
< column "INDEX_NAME" "VARCHAR"
< row
: bytes
| c_ibfk_1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 2 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY (a) USING HASH) ENGINE=InnoDB
< ok 0 0
> query
| SELECT INDEX_TYPE FROM information_schema.STATISTICS
|              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND INDEX_NAME='a'
< result
< column "INDEX_TYPE" "VARCHAR"
< row
: bytes
| HASH
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t1 (a INT, KEY (a))
< ok 0 0
> exec
| CREATE TABLE t2 (a INT, KEY (a) USING BTREE)
< ok 0 0
> query
| SHOW CREATE TABLE t1
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t1
: bytes
| CREATE TABLE `t1` (
|   `a` int(11) DEFAULT NULL,
|   KEY `a` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SHOW CREATE TABLE t2
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t2
: bytes
| CREATE TABLE `t2` (
|   `a` int(11) DEFAULT NULL,
|   KEY `a` (`a`) USING BTREE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS testdb
< ok 2 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT, UNIQUE KEY (a))
< ok 0 0
> exec
| INSERT INTO testdb.t VALUES (NULL)
< ok 1 0
> exec ^
< ok 1 0
> exec ^
< ok 1 0
> query
| SELECT COUNT(*) FROM testdb.t
< result
< column "COUNT(*)" "BIGINT"
< row
: int64 3
> exec
| INSERT INTO testdb.t VALUES (1)
< ok 1 0
> exec ^
< error
| 1062 23000 Duplicate entry '1' for key 'a'
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY ix (a))
< ok 0 0
> query
| SELECT IGNORED FROM information_schema.STATISTICS
|              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND INDEX_NAME='ix'
< result
< column "IGNORED" "VARCHAR"
< row
: bytes
| NO
> exec
| CREATE TABLE t_pkinv (a INT, PRIMARY KEY (a) IGNORED)
< error
| 4174 HY000 A primary key cannot be marked as IGNORE
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t_noprefix (a TEXT, KEY ix (a))
< ok 0 0
> query
| SHOW CREATE TABLE t_noprefix
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_noprefix
: bytes
| CREATE TABLE `t_noprefix` (
|   `a` text DEFAULT NULL,
|   KEY `ix` (`a`(768))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| CREATE TABLE t_prefix (a TEXT, KEY (a(100)))
< ok 0 0
> query
| SELECT SUB_PART FROM information_schema.STATISTICS
|              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t_prefix' AND INDEX_NAME='a'
< result
< column "SUB_PART" "BIGINT"
< row
: int64 100
> exec
| CREATE TABLE t_ft (a TEXT, FULLTEXT KEY (a))
< ok 0 0
> query
| SELECT INDEX_TYPE FROM information_schema.STATISTICS
|              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t_ft' AND INDEX_NAME='a'
< result
< column "INDEX_TYPE" "VARCHAR"
< row
: bytes
| FULLTEXT
> exec
| DROP DATABASE IF EXISTS testdb
< ok 3 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t1 (a TEXT, FULLTEXT KEY (a))
< ok 0 0
> query
| SHOW CREATE TABLE t1
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t1
: bytes
| CREATE TABLE `t1` (
|   `a` text DEFAULT NULL,
|   FULLTEXT KEY `a` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| CREATE TABLE t2 (a TEXT, FULLTEXT KEY (a) WITH PARSER ngram)
< error
| 1128 HY000 Function 'ngram' is not defined
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t_ok (g GEOMETRY NOT NULL, SPATIAL KEY (g))
< ok 0 0
> exec
| CREATE TABLE t_srid (g GEOMETRY NOT NULL SRID 4326)
< error
| 1064 42000 …near '4326)' at line 1
> exec
| CREATE TABLE t_null (g GEOMETRY, SPATIAL KEY (g))
< error
| 1252 42000 All parts of a SPATIAL index must be NOT NULL
> exec
| CREATE TABLE t_spbtree (g GEOMETRY NOT NULL, SPATIAL KEY (g) USING BTREE)
< error
| 1064 42000 …near 'USING BTREE)' at line 1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, PRIMARY KEY (a), UNIQUE KEY uk (a))
< ok 0 0
> query
| SELECT INDEX_NAME FROM information_schema.STATISTICS
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
|             ORDER BY INDEX_NAME
< result
< column "INDEX_NAME" "VARCHAR"
< row
: bytes
| PRIMARY
< row
: bytes
| uk
