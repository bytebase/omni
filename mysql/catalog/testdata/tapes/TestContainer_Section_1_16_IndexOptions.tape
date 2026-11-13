-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

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
|   `a` int DEFAULT NULL,
|   KEY `idx_a` (`a`) COMMENT 'description'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_idx_invis
< ok 0 0
> exec
| CREATE TABLE t_idx_invis (a INT, INDEX idx_a (a) INVISIBLE)
< ok 0 0
> query
| SHOW CREATE TABLE t_idx_invis
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_idx_invis
: bytes
| CREATE TABLE `t_idx_invis` (
|   `a` int DEFAULT NULL,
|   KEY `idx_a` (`a`) /*!80000 INVISIBLE */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
|   `a` int DEFAULT NULL,
|   KEY `idx_a` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
|   `a` int DEFAULT NULL,
|   KEY `idx_a` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
