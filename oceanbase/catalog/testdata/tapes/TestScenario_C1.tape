-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestScenario_C1 default#1
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
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE child (
|     a INT, CONSTRAINT child_ibfk_5 FOREIGN KEY (a) REFERENCES p(id),
|     b INT, FOREIGN KEY (b) REFERENCES p(id)
| )
< ok 0 0
> query
| SELECT CONSTRAINT_NAME FROM information_schema.TABLE_CONSTRAINTS
|         WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='child'
|               AND CONSTRAINT_TYPE='FOREIGN KEY'
|         ORDER BY CONSTRAINT_NAME
< result
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| child_ibfk_5
< row
: bytes
| child_OBFK_0000000000000001
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
| CREATE TABLE child (
|     a INT, b INT,
|     CONSTRAINT child_ibfk_20 FOREIGN KEY (a) REFERENCES p(id)
| )
< ok 0 0
> exec
| ALTER TABLE child ADD FOREIGN KEY (b) REFERENCES p(id)
< ok 0 0
> query
| SELECT CONSTRAINT_NAME FROM information_schema.TABLE_CONSTRAINTS
|         WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='child'
|               AND CONSTRAINT_TYPE='FOREIGN KEY'
|         ORDER BY CONSTRAINT_NAME
< result
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| child_ibfk_20
< row
: bytes
| child_OBFK_0000000000000002
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
| CREATE TABLE t (id INT) PARTITION BY HASH(id) PARTITIONS 4
< ok 0 0
> query
| SELECT PARTITION_NAME FROM information_schema.PARTITIONS
|         WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
|         ORDER BY PARTITION_ORDINAL_POSITION
< result
< column "PARTITION_NAME" "VARCHAR"
< row
: bytes
| p0
< row
: bytes
| p1
< row
: bytes
| p2
< row
: bytes
| p3
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
| CREATE TABLE t (a INT, CHECK (a > 0))
< ok 0 0
> query
| SELECT CONSTRAINT_NAME FROM information_schema.CHECK_CONSTRAINTS
|             WHERE CONSTRAINT_SCHEMA='testdb'
< result
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| t_OBCHECK_0000000000000003
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
| CREATE TABLE t (a INT, UNIQUE KEY (a))
< ok 0 0
> query
| SELECT INDEX_NAME FROM information_schema.STATISTICS
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND NON_UNIQUE=0
< result
< column "INDEX_NAME" "VARCHAR"
< row
: bytes
| a
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
| CREATE TABLE t (a INT, UNIQUE KEY a (a), UNIQUE KEY (a))
< ok 0 0
> query
| SELECT INDEX_NAME FROM information_schema.STATISTICS
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND NON_UNIQUE=0
|             ORDER BY INDEX_NAME
< result
< column "INDEX_NAME" "VARCHAR"
< row
: bytes
| a
< row
: bytes
| a_2
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
| CREATE TABLE t (
|     a INT,
|     CONSTRAINT my_pk PRIMARY KEY (a)
| )
< ok 0 0
> query
| SELECT CONSTRAINT_NAME FROM information_schema.TABLE_CONSTRAINTS
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND CONSTRAINT_TYPE='PRIMARY KEY'
< result
< column "CONSTRAINT_NAME" "VARCHAR"
< row
: bytes
| PRIMARY
> query
| SELECT INDEX_NAME FROM information_schema.STATISTICS
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND INDEX_NAME='PRIMARY'
< result
< column "INDEX_NAME" "VARCHAR"
< row
: bytes
| PRIMARY
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
| CREATE TABLE t (a INT, UNIQUE KEY `PRIMARY` (a))
< error
| 1280 42000 Incorrect index name 'PRIMARY'
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
| CREATE TABLE t (
|     a INT,
|     b INT,
|     c INT,
|     KEY (b, c)
| )
< ok 0 0
> query
| SELECT INDEX_NAME, COLUMN_NAME FROM information_schema.STATISTICS
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
|             ORDER BY INDEX_NAME, SEQ_IN_INDEX
< result
< column "INDEX_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| b
: bytes
| b
< row
: bytes
| b
: bytes
| c
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
| CREATE TABLE t (`PRIMARY` INT)
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE KEY (`PRIMARY`)
< ok 0 0
> query
| SELECT INDEX_NAME FROM information_schema.STATISTICS
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND NON_UNIQUE=0
< result
< column "INDEX_NAME" "VARCHAR"
< row
: bytes
| PRIMARY
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
| CREATE TABLE t (
|     a INT,
|     INDEX ((a + 1)),
|     INDEX ((a * 2))
| )
< ok 0 0
> query
| SELECT DISTINCT INDEX_NAME FROM information_schema.STATISTICS
|             WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t'
|             ORDER BY INDEX_NAME
< result
< column "INDEX_NAME" "VARCHAR"
< row
: bytes
| functional_index
< row
: bytes
| functional_index_2
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
| CREATE TABLE t (a INT, INDEX fx ((a + 1), (a * 2)))
< ok 0 0
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "TEXT"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int(11) DEFAULT NULL,
|   KEY `fx` ((`a` + 1), (`a` * 2)) BLOCK_SIZE 16384 LOCAL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
| CREATE TABLE t1 (a INT, CONSTRAINT mychk CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE t2 (a INT, CONSTRAINT mychk CHECK (a < 100))
< error
| 3822 HY000 Duplicate check constraint name 'mychk'.
