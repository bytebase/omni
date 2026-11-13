-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestContainer_AlterAddIndexOptions default#1
> exec
| CREATE TABLE t (a INT, b INT, c TEXT, v VARCHAR(10), g GEOMETRY NOT NULL)
< ok 0 0
> exec
| ALTER TABLE t ADD INDEX i (a) COMMENT 'hello' INVISIBLE
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE KEY u (b) COMMENT 'u note' KEY_BLOCK_SIZE=8
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE INDEX u2 (a, b) INVISIBLE
< ok 0 0
> exec
| ALTER TABLE t ADD CONSTRAINT cu UNIQUE (v) COMMENT 'cu note'
< ok 0 0
> exec
| ALTER TABLE t ADD KEY k (v) KEY_BLOCK_SIZE=4
< ok 0 0
> exec
| ALTER TABLE t ADD FULLTEXT INDEX f (c) COMMENT 'ft note' INVISIBLE
< ok 0 0
> exec
| ALTER TABLE t ADD SPATIAL INDEX s (g) COMMENT 'sp note'
< ok 0 0
> query
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   `c` text,
|   `v` varchar(10) DEFAULT NULL,
|   `g` geometry NOT NULL,
|   UNIQUE KEY `u` (`b`) COMMENT 'u note',
|   UNIQUE KEY `u2` (`a`,`b`) /*!80000 INVISIBLE */,
|   UNIQUE KEY `cu` (`v`) COMMENT 'cu note',
|   KEY `i` (`a`) COMMENT 'hello' /*!80000 INVISIBLE */,
|   KEY `k` (`v`),
|   SPATIAL KEY `s` (`g`) COMMENT 'sp note',
|   FULLTEXT KEY `f` (`c`) COMMENT 'ft note' /*!80000 INVISIBLE */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| CREATE TABLE tn (id INT)
< ok 0 0
> exec
| ALTER TABLE tn ADD PRIMARY KEY (id) INVISIBLE
< error
| 3522 HY000 A primary key index cannot be invisible
> query
| SHOW CREATE TABLE tn
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| tn
: bytes
| CREATE TABLE `tn` (
|   `id` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
