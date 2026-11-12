-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestScenario_C7 default#1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
| BTREE
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
|   `a` int DEFAULT NULL,
|   KEY `a` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
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
|   `a` int DEFAULT NULL,
|   KEY `a` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT, UNIQUE KEY (a))
< ok 0 0
> exec
| INSERT INTO testdb.t VALUES (NULL)
< ok 1 0
> exec
| INSERT INTO testdb.t VALUES (NULL)
< ok 1 0
> exec
| INSERT INTO testdb.t VALUES (NULL)
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
> exec
| INSERT INTO testdb.t VALUES (1)
< error
| 1062 23000 Duplicate entry '1' for key 't.a'
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT, KEY ix (a))
< ok 0 0
> query
| SELECT IS_VISIBLE FROM information_schema.STATISTICS
|              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND INDEX_NAME='ix'
< result
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| YES
> exec
| CREATE TABLE t_pkinv (a INT, PRIMARY KEY (a) INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t_noprefix (a TEXT, KEY (a))
< error
| 1170 42000 BLOB/TEXT column 'a' used in key specification without a key length
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
| SELECT INDEX_NAME FROM information_schema.STATISTICS
|              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t_ft'
< result
< column "INDEX_NAME" "VARCHAR"
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t1 (a TEXT, FULLTEXT KEY (a))
< ok 0 0
> exec
| CREATE TABLE t2 (a TEXT, FULLTEXT KEY (a) WITH PARSER ngram)
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
|   `a` text DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
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
|   `a` text DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t_ok (g GEOMETRY NOT NULL SRID 4326, SPATIAL KEY (g))
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 29 near "GEOMETRY NOT NULL SRID 4326, SPATIAL KEY (g))" 
> exec
| CREATE TABLE t_null (g GEOMETRY, SPATIAL KEY (g))
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 31 near "GEOMETRY, SPATIAL KEY (g))" 
> exec
| CREATE TABLE t_spbtree (g GEOMETRY NOT NULL, SPATIAL KEY (g) USING BTREE)
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 34 near "GEOMETRY NOT NULL, SPATIAL KEY (g) USING BTREE)" 
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
