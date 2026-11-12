-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestContainerAlterAddIndexOptions default#1
> exec
| DROP DATABASE IF EXISTS omni_alter_index
< ok 0 0
> exec
| CREATE DATABASE omni_alter_index
< ok 0 0
> exec
| USE omni_alter_index
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, a INT, b INT, v VARCHAR(10))
< ok 0 0
> exec
| ALTER TABLE t ADD PRIMARY KEY (id) COMMENT 'pk note'
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
|   `id` int NOT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   `v` varchar(10) DEFAULT NULL,
|   PRIMARY KEY (`id`) COMMENT 'pk note' /*T![clustered_index] NONCLUSTERED */,
|   KEY `i` (`a`) /*!80000 INVISIBLE */ COMMENT 'hello',
|   UNIQUE KEY `u` (`b`) COMMENT 'u note',
|   UNIQUE KEY `u2` (`a`,`b`) /*!80000 INVISIBLE */,
|   UNIQUE KEY `cu` (`v`) COMMENT 'cu note',
|   KEY `k` (`v`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP DATABASE IF EXISTS omni_alter_index
< ok 0 0
