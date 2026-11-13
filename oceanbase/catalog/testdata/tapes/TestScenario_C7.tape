-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

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
< ok 2 0
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
< ok 2 0
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
| 1064 42000 …near 'INVISIBLE)' at line 1
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
| CREATE TABLE t_noprefix (a TEXT, KEY (a))
< error
| 1167 42000 The used storage engine can't index column 'a'
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
< row
: bytes
| FULLTEXT
> exec
| DROP DATABASE IF EXISTS testdb
< ok 11 0
> exec
| CREATE DATABASE testdb
< ok 1 0
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
< column "Create Table" "TEXT"
< row
: bytes
| t1
: bytes
| CREATE TABLE `t1` (
|   `a` text DEFAULT NULL,
|   FULLTEXT KEY `a` (`a`) WITH PARSER space PARSER_PROPERTIES=(min_token_size=3,max_token_size=84) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> query
| SHOW CREATE TABLE t2
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t2
: bytes
| CREATE TABLE `t2` (
|   `a` text DEFAULT NULL,
|   FULLTEXT KEY `a` (`a`) WITH PARSER ngram PARSER_PROPERTIES=(ngram_token_size=2) BLOCK_SIZE 16384
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
> exec
| DROP DATABASE IF EXISTS testdb
< ok 14 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t_ok (g GEOMETRY NOT NULL SRID 4326, SPATIAL KEY (g))
< ok 0 0
> exec
| CREATE TABLE t_null (g GEOMETRY, SPATIAL KEY (g))
< error
| 1252 42000 All parts of a SPATIAL index must be NOT NULL.
> exec
| CREATE TABLE t_spbtree (g GEOMETRY NOT NULL, SPATIAL KEY (g) USING BTREE)
< ok 0 0
> exec
| DROP DATABASE IF EXISTS testdb
< ok 8 0
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
