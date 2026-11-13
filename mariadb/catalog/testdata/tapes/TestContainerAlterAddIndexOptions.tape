-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestContainerAlterAddIndexOptions default#1
> exec
| CREATE TABLE t (a INT, b INT, c TEXT, v VARCHAR(10), g GEOMETRY NOT NULL)
< ok 0 0
> exec
| ALTER TABLE t ADD INDEX i (a) COMMENT 'hello'
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE KEY u (b) COMMENT 'u note' KEY_BLOCK_SIZE=8
< ok 0 0
> exec
| ALTER TABLE t ADD CONSTRAINT cu UNIQUE (v) COMMENT 'cu note'
< ok 0 0
> exec
| ALTER TABLE t ADD KEY k (v) KEY_BLOCK_SIZE=4
< ok 0 0
> exec
| ALTER TABLE t ADD FULLTEXT INDEX f (c) COMMENT 'ft note'
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
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   `c` text DEFAULT NULL,
|   `v` varchar(10) DEFAULT NULL,
|   `g` geometry NOT NULL,
|   UNIQUE KEY `u` (`b`) KEY_BLOCK_SIZE=8 COMMENT 'u note',
|   UNIQUE KEY `cu` (`v`) COMMENT 'cu note',
|   KEY `i` (`a`) COMMENT 'hello',
|   KEY `k` (`v`) KEY_BLOCK_SIZE=4,
|   SPATIAL KEY `s` (`g`) COMMENT 'sp note',
|   FULLTEXT KEY `f` (`c`) COMMENT 'ft note'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
