-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestBackupRestoreRoundTrip default#1
> exec
| CREATE TABLE limits (id INT PRIMARY KEY, a INT, b INT, g INT AS (a + b), UNIQUE KEY uk_b (b))
< ok 0 0
> exec
| INSERT INTO limits (id, a, b) VALUES (1, 1, 10), (2, 2, 20), (3, 3, 30), (4, 4, 40)
< ok 4 0
> query
| SELECT * FROM `limits` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< column "b" "INT"
< column "g" "INT"
< row
: int64 1
: int64 1
: int64 10
: int64 11
< row
: int64 2
: int64 2
: int64 20
: int64 22
< row
: int64 3
: int64 3
: int64 30
: int64 33
< row
: int64 4
: int64 4
: int64 40
: int64 44
> query
| SELECT COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_t_mysql_analysis_backup
: string
| limits
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< row
: bytes
| id
: bytes
|
< row
: bytes
| a
: bytes
|
< row
: bytes
| b
: bytes
|
< row
: bytes
| g
: bytes
| (`a` + `b`)
> query
| SELECT INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY INDEX_NAME, SEQ_IN_INDEX
: string
| omni_t_mysql_analysis_backup
: string
| limits
< result
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "TEXT"
< row
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| uk_b
: int64 0
: bytes
| b
> exec
| CREATE TABLE `omni_t_mysql_analysis_backup`.`bak_limits_omni_t_mysql_analysis_backup` LIKE `omni_t_mysql_analysis_backup`.`limits`
< ok 0 0
> exec
|
| INSERT INTO `omni_t_mysql_analysis_backup`.`bak_limits_omni_t_mysql_analysis_backup` (`id`,`a`,`b`)  (SELECT `limits`.`id`,`limits`.`a`,`limits`.`b` FROM limits WHERE id = 1)
|   UNION DISTINCT
|   (SELECT `limits`.`id`,`limits`.`a`,`limits`.`b` FROM limits WHERE id > 1 ORDER BY id DESC LIMIT 1)
< ok 2 0
> exec
| DELETE FROM limits WHERE id = 1
< ok 1 0
> exec
|  DELETE FROM limits WHERE id > 1 ORDER BY id DESC LIMIT 1
< ok 1 0
> query
| SELECT * FROM `limits` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< column "b" "INT"
< column "g" "INT"
< row
: int64 2
: int64 2
: int64 20
: int64 22
< row
: int64 3
: int64 3
: int64 30
: int64 33
> exec
| /*
| Original SQL:
| DELETE FROM limits WHERE id = 1;
|  DELETE FROM limits WHERE id > 1 ORDER BY id DESC LIMIT 1;
| */
| INSERT INTO `omni_t_mysql_analysis_backup`.`limits` (`id`, `a`, `b`) SELECT `id`, `a`, `b` FROM `omni_t_mysql_analysis_backup`.`bak_limits_omni_t_mysql_analysis_backup`
< ok 2 0
> query
| SELECT * FROM `limits` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< column "b" "INT"
< column "g" "INT"
< row
: int64 1
: int64 1
: int64 10
: int64 11
< row
: int64 2
: int64 2
: int64 20
: int64 22
< row
: int64 3
: int64 3
: int64 30
: int64 33
< row
: int64 4
: int64 4
: int64 40
: int64 44
> exec
| CREATE TABLE update_limits (id INT PRIMARY KEY, a INT, g INT AS (a * 2) STORED)
< ok 0 0
> exec
| INSERT INTO update_limits (id, a) VALUES (1, 1), (2, 2), (3, 3)
< ok 3 0
> query
| SELECT * FROM `update_limits` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< column "g" "INT"
< row
: int64 1
: int64 1
: int64 2
< row
: int64 2
: int64 2
: int64 4
< row
: int64 3
: int64 3
: int64 6
> query
| SELECT COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_t_mysql_analysis_backup
: string
| update_limits
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< row
: bytes
| id
: bytes
|
< row
: bytes
| a
: bytes
|
< row
: bytes
| g
: bytes
| (`a` * 2)
> query
| SELECT INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY INDEX_NAME, SEQ_IN_INDEX
: string
| omni_t_mysql_analysis_backup
: string
| update_limits
< result
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "TEXT"
< row
: bytes
| PRIMARY
: int64 0
: bytes
| id
> exec
| CREATE TABLE `omni_t_mysql_analysis_backup`.`bak_update_limits_omni_t_mysql_analysis_backup` LIKE `omni_t_mysql_analysis_backup`.`update_limits`
< ok 0 0
> exec
|
| INSERT INTO `omni_t_mysql_analysis_backup`.`bak_update_limits_omni_t_mysql_analysis_backup` (`id`,`a`)  (SELECT `update_limits`.`id`,`update_limits`.`a` FROM update_limits WHERE id = 1)
|   UNION DISTINCT
|   (SELECT `update_limits`.`id`,`update_limits`.`a` FROM update_limits WHERE id > 1 ORDER BY id LIMIT 1)
< ok 2 0
> exec
| UPDATE update_limits SET a = a + 100 WHERE id = 1
< ok 1 0
> exec
|  UPDATE update_limits SET a = a + 100 WHERE id > 1 ORDER BY id LIMIT 1
< ok 1 0
> query
| SELECT * FROM `update_limits` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< column "g" "INT"
< row
: int64 1
: int64 101
: int64 202
< row
: int64 2
: int64 102
: int64 204
< row
: int64 3
: int64 3
: int64 6
> exec
| /*
| Original SQL:
| UPDATE update_limits SET a = a + 100 WHERE id = 1;
|  UPDATE update_limits SET a = a + 100 WHERE id > 1 ORDER BY id LIMIT 1;
| */
| INSERT INTO `omni_t_mysql_analysis_backup`.`update_limits` (`id`, `a`) SELECT `id`, `a` FROM `omni_t_mysql_analysis_backup`.`bak_update_limits_omni_t_mysql_analysis_backup` ON DUPLICATE KEY UPDATE `a` = VALUES(`a`)
< ok 4 0
> query
| SELECT * FROM `update_limits` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< column "g" "INT"
< row
: int64 1
: int64 1
: int64 2
< row
: int64 2
: int64 2
: int64 4
< row
: int64 3
: int64 3
: int64 6
> exec
| CREATE TABLE self_join (id INT PRIMARY KEY, a INT, b INT)
< ok 0 0
> exec
| INSERT INTO self_join VALUES (1, 1, 1), (2, 2, 2), (3, 3, 3)
< ok 3 0
> query
| SELECT * FROM `self_join` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< column "b" "INT"
< row
: int64 1
: int64 1
: int64 1
< row
: int64 2
: int64 2
: int64 2
< row
: int64 3
: int64 3
: int64 3
> query
| SELECT COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_t_mysql_analysis_backup
: string
| self_join
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< row
: bytes
| id
: bytes
|
< row
: bytes
| a
: bytes
|
< row
: bytes
| b
: bytes
|
> query
| SELECT INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY INDEX_NAME, SEQ_IN_INDEX
: string
| omni_t_mysql_analysis_backup
: string
| self_join
< result
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "TEXT"
< row
: bytes
| PRIMARY
: int64 0
: bytes
| id
> exec
| CREATE TABLE `omni_t_mysql_analysis_backup`.`bak_self_join_omni_t_mysql_analysis_backup` LIKE `omni_t_mysql_analysis_backup`.`self_join`
< ok 0 0
> exec
|
| INSERT INTO `omni_t_mysql_analysis_backup`.`bak_self_join_omni_t_mysql_analysis_backup`  (SELECT `x`.* FROM self_join x JOIN self_join y ON x.id = y.id + 1)
|   UNION DISTINCT
|   (SELECT `y`.* FROM self_join x JOIN self_join y ON x.id = y.id + 1)
< ok 3 0
> exec
| UPDATE self_join x JOIN self_join y ON x.id = y.id + 1 SET x.a = 10, y.b = 20
< ok 4 0
> query
| SELECT * FROM `self_join` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< column "b" "INT"
< row
: int64 1
: int64 1
: int64 20
< row
: int64 2
: int64 10
: int64 20
< row
: int64 3
: int64 10
: int64 3
> exec
| /*
| Original SQL:
| UPDATE self_join x JOIN self_join y ON x.id = y.id + 1 SET x.a = 10, y.b = 20;
| */
| INSERT INTO `omni_t_mysql_analysis_backup`.`self_join` (`id`, `a`, `b`) SELECT `id`, `a`, `b` FROM `omni_t_mysql_analysis_backup`.`bak_self_join_omni_t_mysql_analysis_backup` ON DUPLICATE KEY UPDATE `a` = VALUES(`a`), `b` = VALUES(`b`)
< ok 6 0
> query
| SELECT * FROM `self_join` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< column "b" "INT"
< row
: int64 1
: int64 1
: int64 1
< row
: int64 2
: int64 2
: int64 2
< row
: int64 3
: int64 3
: int64 3
> exec
| CREATE TABLE comments (id INT PRIMARY KEY, note VARCHAR(10))
< ok 0 0
> exec
| INSERT INTO comments VALUES (1, '*/'), (2, 'x')
< ok 2 0
> query
| SELECT * FROM `comments` ORDER BY 1, 2
< result
< column "id" "INT"
< column "note" "VARCHAR"
< row
: int64 1
: bytes
| */
< row
: int64 2
: bytes
| x
> query
| SELECT COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_t_mysql_analysis_backup
: string
| comments
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< row
: bytes
| id
: bytes
|
< row
: bytes
| note
: bytes
|
> query
| SELECT INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY INDEX_NAME, SEQ_IN_INDEX
: string
| omni_t_mysql_analysis_backup
: string
| comments
< result
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "TEXT"
< row
: bytes
| PRIMARY
: int64 0
: bytes
| id
> exec
| CREATE TABLE `omni_t_mysql_analysis_backup`.`bak_comments_omni_t_mysql_analysis_backup` LIKE `omni_t_mysql_analysis_backup`.`comments`
< ok 0 0
> exec
|
| INSERT INTO `omni_t_mysql_analysis_backup`.`bak_comments_omni_t_mysql_analysis_backup`  SELECT `comments`.* FROM comments WHERE note = '*/'
< ok 1 0
> exec
| DELETE /* note */ FROM comments WHERE note = '*/'
< ok 1 0
> query
| SELECT * FROM `comments` ORDER BY 1, 2
< result
< column "id" "INT"
< column "note" "VARCHAR"
< row
: int64 2
: bytes
| x
> exec
| /*
| Original SQL:
| DELETE /* note * / FROM comments WHERE note = '* /';
| */
| INSERT INTO `omni_t_mysql_analysis_backup`.`comments` (`id`, `note`) SELECT `id`, `note` FROM `omni_t_mysql_analysis_backup`.`bak_comments_omni_t_mysql_analysis_backup`
< ok 1 0
> query
| SELECT * FROM `comments` ORDER BY 1, 2
< result
< column "id" "INT"
< column "note" "VARCHAR"
< row
: int64 1
: bytes
| */
< row
: int64 2
: bytes
| x
> exec
| CREATE TABLE cte_alias (id INT PRIMARY KEY, a INT)
< ok 0 0
> exec
| INSERT INTO cte_alias VALUES (1, 1), (2, 2)
< ok 2 0
> query
| SELECT * FROM `cte_alias` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< row
: int64 1
: int64 1
< row
: int64 2
: int64 2
> query
| SELECT COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_t_mysql_analysis_backup
: string
| cte_alias
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< row
: bytes
| id
: bytes
|
< row
: bytes
| a
: bytes
|
> query
| SELECT INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY INDEX_NAME, SEQ_IN_INDEX
: string
| omni_t_mysql_analysis_backup
: string
| cte_alias
< result
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "TEXT"
< row
: bytes
| PRIMARY
: int64 0
: bytes
| id
> exec
| CREATE TABLE `omni_t_mysql_analysis_backup`.`bak_cte_alias_omni_t_mysql_analysis_backup` LIKE `omni_t_mysql_analysis_backup`.`cte_alias`
< ok 0 0
> exec
|
| INSERT INTO `omni_t_mysql_analysis_backup`.`bak_cte_alias_omni_t_mysql_analysis_backup`  WITH c AS (SELECT 1 AS id) SELECT `c`.* FROM cte_alias AS c WHERE c.id = 1
< ok 1 0
> exec
| WITH c AS (SELECT 1 AS id) DELETE c FROM cte_alias AS c WHERE c.id = 1
< ok 1 0
> query
| SELECT * FROM `cte_alias` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< row
: int64 2
: int64 2
> exec
| /*
| Original SQL:
| WITH c AS (SELECT 1 AS id) DELETE c FROM cte_alias AS c WHERE c.id = 1;
| */
| INSERT INTO `omni_t_mysql_analysis_backup`.`cte_alias` (`id`, `a`) SELECT `id`, `a` FROM `omni_t_mysql_analysis_backup`.`bak_cte_alias_omni_t_mysql_analysis_backup`
< ok 1 0
> query
| SELECT * FROM `cte_alias` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< row
: int64 1
: int64 1
< row
: int64 2
: int64 2
> exec
| CREATE TABLE generated_key (id INT, a INT, g INT AS (a + 1), UNIQUE KEY uk_g (g))
< ok 0 0
> exec
| INSERT INTO generated_key (id, a) VALUES (1, 1), (2, 2)
< ok 2 0
> query
| SELECT * FROM `generated_key` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< column "g" "INT"
< row
: int64 1
: int64 1
: int64 2
< row
: int64 2
: int64 2
: int64 3
> query
| SELECT COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_t_mysql_analysis_backup
: string
| generated_key
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< row
: bytes
| id
: bytes
|
< row
: bytes
| a
: bytes
|
< row
: bytes
| g
: bytes
| (`a` + 1)
> query
| SELECT INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY INDEX_NAME, SEQ_IN_INDEX
: string
| omni_t_mysql_analysis_backup
: string
| generated_key
< result
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "TEXT"
< row
: bytes
| uk_g
: int64 0
: bytes
| g
> exec
| CREATE TABLE `omni_t_mysql_analysis_backup`.`bak_generated_key_omni_t_mysql_analysis_backup` LIKE `omni_t_mysql_analysis_backup`.`generated_key`
< ok 0 0
> exec
|
| INSERT INTO `omni_t_mysql_analysis_backup`.`bak_generated_key_omni_t_mysql_analysis_backup` (`id`,`a`)  SELECT `generated_key`.`id`,`generated_key`.`a` FROM generated_key WHERE id = 1
< ok 1 0
> exec
| UPDATE generated_key SET a = a + 10 WHERE id = 1
< ok 1 0
> query
| SELECT * FROM `generated_key` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< column "g" "INT"
< row
: int64 1
: int64 11
: int64 12
< row
: int64 2
: int64 2
: int64 3
> exec
| CREATE TABLE expression_key (id INT, a INT, UNIQUE KEY uk_a ((a + 1)))
< ok 0 0
> exec
| INSERT INTO expression_key VALUES (1, 1), (2, 2)
< ok 2 0
> query
| SELECT * FROM `expression_key` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< row
: int64 1
: int64 1
< row
: int64 2
: int64 2
> query
| SELECT COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_t_mysql_analysis_backup
: string
| expression_key
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< row
: bytes
| id
: bytes
|
< row
: bytes
| a
: bytes
|
> query
| SELECT INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY INDEX_NAME, SEQ_IN_INDEX
: string
| omni_t_mysql_analysis_backup
: string
| expression_key
< result
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "INT"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "TEXT"
< row
: bytes
| uk_a
: int64 0
: bytes
| (`a` + 1)
> exec
| CREATE TABLE `omni_t_mysql_analysis_backup`.`bak_expression_key_omni_t_mysql_analysis_backup` LIKE `omni_t_mysql_analysis_backup`.`expression_key`
< ok 0 0
> exec
|
| INSERT INTO `omni_t_mysql_analysis_backup`.`bak_expression_key_omni_t_mysql_analysis_backup`  SELECT `expression_key`.* FROM expression_key WHERE id = 1
< ok 1 0
> exec
| UPDATE expression_key SET a = a + 10 WHERE id = 1
< ok 1 0
> query
| SELECT * FROM `expression_key` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< row
: int64 1
: int64 11
< row
: int64 2
: int64 2
