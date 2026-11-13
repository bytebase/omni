-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestContainer_Section_1_16_IndexOptions default#1
> exec
| DROP TABLE IF EXISTS t_idx_comment
< ok 0 0
> exec
| CREATE TABLE t_idx_comment (a INT, INDEX idx_a (a) COMMENT 'description')
< ok 0 0
> query
| SHOW CREATE TABLE t_idx_comment
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_idx_comment
: bytes
| CREATE TABLE `t_idx_comment` (
|   `a` int(11) DEFAULT NULL,
|   KEY `idx_a` (`a`) COMMENT 'description'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_idx_vis
< ok 0 0
> exec
| CREATE TABLE t_idx_vis (a INT, INDEX idx_a (a) VISIBLE)
< ok 0 0
> query
| SHOW CREATE TABLE t_idx_vis
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_idx_vis
: bytes
| CREATE TABLE `t_idx_vis` (
|   `a` int(11) DEFAULT NULL,
|   KEY `idx_a` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_idx_kbs
< ok 0 0
> exec
| CREATE TABLE t_idx_kbs (a INT, INDEX idx_a (a) KEY_BLOCK_SIZE=4)
< ok 0 0
> query
| SHOW CREATE TABLE t_idx_kbs
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_idx_kbs
: bytes
| CREATE TABLE `t_idx_kbs` (
|   `a` int(11) DEFAULT NULL,
|   KEY `idx_a` (`a`) KEY_BLOCK_SIZE=4
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
