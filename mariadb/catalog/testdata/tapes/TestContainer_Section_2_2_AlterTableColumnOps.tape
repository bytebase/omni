-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestContainer_Section_2_2_AlterTableColumnOps default#1
> exec
| DROP TABLE IF EXISTS t_add_end
< ok 0 0
> exec
| CREATE TABLE t_add_end (id INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE t_add_end ADD COLUMN name VARCHAR(100)
< ok 0 0
> query
| SHOW CREATE TABLE t_add_end
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_add_end
: bytes
| CREATE TABLE `t_add_end` (
|   `id` int(11) NOT NULL,
|   `name` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_add_first
< ok 0 0
> exec
| CREATE TABLE t_add_first (id INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE t_add_first ADD COLUMN name VARCHAR(100) FIRST
< ok 0 0
> query
| SHOW CREATE TABLE t_add_first
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_add_first
: bytes
| CREATE TABLE `t_add_first` (
|   `name` varchar(100) DEFAULT NULL,
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_add_after
< ok 0 0
> exec
| CREATE TABLE t_add_after (id INT PRIMARY KEY, age INT)
< ok 0 0
> exec
| ALTER TABLE t_add_after ADD COLUMN name VARCHAR(100) AFTER id
< ok 0 0
> query
| SHOW CREATE TABLE t_add_after
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_add_after
: bytes
| CREATE TABLE `t_add_after` (
|   `id` int(11) NOT NULL,
|   `name` varchar(100) DEFAULT NULL,
|   `age` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_add_multi
< ok 0 0
> exec
| CREATE TABLE t_add_multi (id INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE t_add_multi ADD COLUMN name VARCHAR(100), ADD COLUMN age INT, ADD COLUMN email VARCHAR(255)
< ok 0 0
> query
| SHOW CREATE TABLE t_add_multi
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_add_multi
: bytes
| CREATE TABLE `t_add_multi` (
|   `id` int(11) NOT NULL,
|   `name` varchar(100) DEFAULT NULL,
|   `age` int(11) DEFAULT NULL,
|   `email` varchar(255) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_drop_col
< ok 0 0
> exec
| CREATE TABLE t_drop_col (id INT PRIMARY KEY, name VARCHAR(100), age INT)
< ok 0 0
> exec
| ALTER TABLE t_drop_col DROP COLUMN age
< ok 0 0
> query
| SHOW CREATE TABLE t_drop_col
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_drop_col
: bytes
| CREATE TABLE `t_drop_col` (
|   `id` int(11) NOT NULL,
|   `name` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_drop_idx_col
< ok 0 0
> exec
| CREATE TABLE t_drop_idx_col (id INT PRIMARY KEY, a INT, b INT, KEY idx_ab (a, b))
< ok 0 0
> exec
| ALTER TABLE t_drop_idx_col DROP COLUMN b
< ok 0 0
> query
| SHOW CREATE TABLE t_drop_idx_col
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_drop_idx_col
: bytes
| CREATE TABLE `t_drop_idx_col` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_ab` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_drop_only_idx
< ok 0 0
> exec
| CREATE TABLE t_drop_only_idx (id INT PRIMARY KEY, a INT, KEY idx_a (a))
< ok 0 0
> exec
| ALTER TABLE t_drop_only_idx DROP COLUMN a
< ok 0 0
> query
| SHOW CREATE TABLE t_drop_only_idx
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_drop_only_idx
: bytes
| CREATE TABLE `t_drop_only_idx` (
|   `id` int(11) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_mod_type
< ok 0 0
> exec
| CREATE TABLE t_mod_type (id INT PRIMARY KEY, val SMALLINT)
< ok 0 0
> exec
| ALTER TABLE t_mod_type MODIFY COLUMN val INT
< ok 0 0
> query
| SHOW CREATE TABLE t_mod_type
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_mod_type
: bytes
| CREATE TABLE `t_mod_type` (
|   `id` int(11) NOT NULL,
|   `val` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_mod_widen
< ok 0 0
> exec
| CREATE TABLE t_mod_widen (id INT PRIMARY KEY, name VARCHAR(50))
< ok 0 0
> exec
| ALTER TABLE t_mod_widen MODIFY COLUMN name VARCHAR(200)
< ok 0 0
> query
| SHOW CREATE TABLE t_mod_widen
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_mod_widen
: bytes
| CREATE TABLE `t_mod_widen` (
|   `id` int(11) NOT NULL,
|   `name` varchar(200) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_mod_narrow
< ok 0 0
> exec
| CREATE TABLE t_mod_narrow (id INT PRIMARY KEY, name VARCHAR(200))
< ok 0 0
> exec
| ALTER TABLE t_mod_narrow MODIFY COLUMN name VARCHAR(50)
< ok 0 0
> query
| SHOW CREATE TABLE t_mod_narrow
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_mod_narrow
: bytes
| CREATE TABLE `t_mod_narrow` (
|   `id` int(11) NOT NULL,
|   `name` varchar(50) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_mod_bigint
< ok 0 0
> exec
| CREATE TABLE t_mod_bigint (id INT PRIMARY KEY, val INT)
< ok 0 0
> exec
| ALTER TABLE t_mod_bigint MODIFY COLUMN val BIGINT
< ok 0 0
> query
| SHOW CREATE TABLE t_mod_bigint
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_mod_bigint
: bytes
| CREATE TABLE `t_mod_bigint` (
|   `id` int(11) NOT NULL,
|   `val` bigint(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_mod_nn
< ok 0 0
> exec
| CREATE TABLE t_mod_nn (id INT PRIMARY KEY, name VARCHAR(100))
< ok 0 0
> exec
| ALTER TABLE t_mod_nn MODIFY COLUMN name VARCHAR(100) NOT NULL
< ok 0 0
> query
| SHOW CREATE TABLE t_mod_nn
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_mod_nn
: bytes
| CREATE TABLE `t_mod_nn` (
|   `id` int(11) NOT NULL,
|   `name` varchar(100) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_mod_rnn
< ok 0 0
> exec
| CREATE TABLE t_mod_rnn (id INT PRIMARY KEY, name VARCHAR(100) NOT NULL)
< ok 0 0
> exec
| ALTER TABLE t_mod_rnn MODIFY COLUMN name VARCHAR(100)
< ok 0 0
> query
| SHOW CREATE TABLE t_mod_rnn
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_mod_rnn
: bytes
| CREATE TABLE `t_mod_rnn` (
|   `id` int(11) NOT NULL,
|   `name` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_mod_pos
< ok 0 0
> exec
| CREATE TABLE t_mod_pos (a INT, b INT, c INT)
< ok 0 0
> exec
| ALTER TABLE t_mod_pos MODIFY COLUMN c INT FIRST
< ok 0 0
> query
| SHOW CREATE TABLE t_mod_pos
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_mod_pos
: bytes
| CREATE TABLE `t_mod_pos` (
|   `c` int(11) DEFAULT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_chg_rt
< ok 0 0
> exec
| CREATE TABLE t_chg_rt (id INT PRIMARY KEY, old_name VARCHAR(50))
< ok 0 0
> exec
| ALTER TABLE t_chg_rt CHANGE COLUMN old_name new_name VARCHAR(100)
< ok 0 0
> query
| SHOW CREATE TABLE t_chg_rt
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_chg_rt
: bytes
| CREATE TABLE `t_chg_rt` (
|   `id` int(11) NOT NULL,
|   `new_name` varchar(100) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_chg_st
< ok 0 0
> exec
| CREATE TABLE t_chg_st (id INT PRIMARY KEY, val INT)
< ok 0 0
> exec
| ALTER TABLE t_chg_st CHANGE COLUMN val val BIGINT
< ok 0 0
> query
| SHOW CREATE TABLE t_chg_st
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_chg_st
: bytes
| CREATE TABLE `t_chg_st` (
|   `id` int(11) NOT NULL,
|   `val` bigint(20) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_chg_idx
< ok 0 0
> exec
| CREATE TABLE t_chg_idx (id INT PRIMARY KEY, a INT, KEY idx_a (a))
< ok 0 0
> exec
| ALTER TABLE t_chg_idx CHANGE COLUMN a b INT
< ok 0 0
> query
| SHOW CREATE TABLE t_chg_idx
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_chg_idx
: bytes
| CREATE TABLE `t_chg_idx` (
|   `id` int(11) NOT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_a` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_ren_col
< ok 0 0
> exec
| CREATE TABLE t_ren_col (id INT PRIMARY KEY, old_col INT)
< ok 0 0
> exec
| ALTER TABLE t_ren_col RENAME COLUMN old_col TO new_col
< ok 0 0
> query
| SHOW CREATE TABLE t_ren_col
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ren_col
: bytes
| CREATE TABLE `t_ren_col` (
|   `id` int(11) NOT NULL,
|   `new_col` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_ren_idx
< ok 0 0
> exec
| CREATE TABLE t_ren_idx (id INT PRIMARY KEY, a INT, KEY idx_a (a))
< ok 0 0
> exec
| ALTER TABLE t_ren_idx RENAME COLUMN a TO b
< ok 0 0
> query
| SHOW CREATE TABLE t_ren_idx
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_ren_idx
: bytes
| CREATE TABLE `t_ren_idx` (
|   `id` int(11) NOT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `idx_a` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_set_def
< ok 0 0
> exec
| CREATE TABLE t_set_def (id INT PRIMARY KEY, val INT)
< ok 0 0
> exec
| ALTER TABLE t_set_def ALTER COLUMN val SET DEFAULT 42
< ok 0 0
> query
| SHOW CREATE TABLE t_set_def
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_set_def
: bytes
| CREATE TABLE `t_set_def` (
|   `id` int(11) NOT NULL,
|   `val` int(11) DEFAULT 42,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_drop_def
< ok 0 0
> exec
| CREATE TABLE t_drop_def (id INT PRIMARY KEY, val INT DEFAULT 10)
< ok 0 0
> exec
| ALTER TABLE t_drop_def ALTER COLUMN val DROP DEFAULT
< ok 0 0
> query
| SHOW CREATE TABLE t_drop_def
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_drop_def
: bytes
| CREATE TABLE `t_drop_def` (
|   `id` int(11) NOT NULL,
|   `val` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_vis
< ok 0 0
> exec
| CREATE TABLE t_vis (id INT PRIMARY KEY, val INT /*!80023 INVISIBLE */)
< ok 0 0
> exec
| ALTER TABLE t_vis MODIFY val INT
< ok 0 0
> query
| SHOW CREATE TABLE t_vis
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_vis
: bytes
| CREATE TABLE `t_vis` (
|   `id` int(11) NOT NULL,
|   `val` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_invis
< ok 0 0
> exec
| CREATE TABLE t_invis (id INT PRIMARY KEY, val INT)
< ok 0 0
> exec
| ALTER TABLE t_invis MODIFY val INT INVISIBLE
< ok 0 0
> query
| SHOW CREATE TABLE t_invis
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_invis
: bytes
| CREATE TABLE `t_invis` (
|   `id` int(11) NOT NULL,
|   `val` int(11) DEFAULT NULL INVISIBLE,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS t_fk_child
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_fk_parent
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_fk_child
< ok 0 0
> exec
| CREATE TABLE t_fk_parent (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t_fk_child (id INT PRIMARY KEY, pid INT, FOREIGN KEY (pid) REFERENCES t_fk_parent(id))
< ok 0 0
> exec
| ALTER TABLE t_fk_child DROP COLUMN pid
< error
| 1553 HY000 Cannot drop index 'pid': needed in a foreign key constraint
