-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestBackupRestoreRoundTrip default#1
> exec
| CREATE TABLE limits (id INT PRIMARY KEY, a INT, b INT, g INT AS (a + b), UNIQUE KEY uk_b (b))
< ok 0 0
> exec
| INSERT INTO limits (id, a, b) VALUES (1, 1, 10), (2, 2, 20), (3, 3, 30), (4, 4, 40), (5, 5, 50), (6, 6, 60), (7, 7, 70), (8, 8, 80), (9, 9, 90), (10, 10, 100), (11, 11, 110), (12, 12, 120)
< ok 12 0
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
< row
: int64 5
: int64 5
: int64 50
: int64 55
< row
: int64 6
: int64 6
: int64 60
: int64 66
< row
: int64 7
: int64 7
: int64 70
: int64 77
< row
: int64 8
: int64 8
: int64 80
: int64 88
< row
: int64 9
: int64 9
: int64 90
: int64 99
< row
: int64 10
: int64 10
: int64 100
: int64 110
< row
: int64 11
: int64 11
: int64 110
: int64 121
< row
: int64 12
: int64 12
: int64 120
: int64 132
> query
| SELECT COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, '') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
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
| `a` + `b`
> query
| SELECT INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY INDEX_NAME, SEQ_IN_INDEX
: string
| test
: string
| limits
< result
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
< row
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| uk_b
: bytes
| 0
: bytes
| b
> exec
| CREATE TABLE `test`.`bak_limits_test` LIKE `test`.`limits`
< ok 0 0
> exec
|
| INSERT INTO `test`.`bak_limits_test` (`id`,`a`,`b`)  (SELECT `limits`.`id`,`limits`.`a`,`limits`.`b` FROM limits WHERE id IN (1, 2) ORDER BY id DESC LIMIT 1)
|   UNION ALL
|   (SELECT `limits`.`id`,`limits`.`a`,`limits`.`b` FROM limits WHERE id IN (11, 12) ORDER BY id DESC LIMIT 1)
|   UNION ALL
|   (SELECT `limits`.`id`,`limits`.`a`,`limits`.`b` FROM limits WHERE id IN (3, 4) ORDER BY id DESC LIMIT 1)
|   UNION ALL
|   (SELECT `limits`.`id`,`limits`.`a`,`limits`.`b` FROM limits WHERE id IN (5, 6) ORDER BY id DESC LIMIT 1)
|   UNION ALL
|   (SELECT `limits`.`id`,`limits`.`a`,`limits`.`b` FROM limits WHERE id IN (7, 8) ORDER BY id DESC LIMIT 1)
|   UNION ALL
|   (SELECT `limits`.`id`,`limits`.`a`,`limits`.`b` FROM limits WHERE id IN (9, 10) ORDER BY id DESC LIMIT 1)
< ok 6 0
> exec
| DELETE FROM limits WHERE id IN (1, 2) ORDER BY id DESC LIMIT 1
< ok 1 0
> exec
|
| DELETE FROM limits WHERE id IN (3, 4) ORDER BY id DESC LIMIT 1
< ok 1 0
> exec
|
| DELETE FROM limits WHERE id IN (5, 6) ORDER BY id DESC LIMIT 1
< ok 1 0
> exec
|
| DELETE FROM limits WHERE id IN (7, 8) ORDER BY id DESC LIMIT 1
< ok 1 0
> exec
|
| DELETE FROM limits WHERE id IN (9, 10) ORDER BY id DESC LIMIT 1
< ok 1 0
> exec
|
| DELETE FROM limits WHERE id IN (11, 12) ORDER BY id DESC LIMIT 1
< ok 1 0
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
: int64 3
: int64 3
: int64 30
: int64 33
< row
: int64 5
: int64 5
: int64 50
: int64 55
< row
: int64 7
: int64 7
: int64 70
: int64 77
< row
: int64 9
: int64 9
: int64 90
: int64 99
< row
: int64 11
: int64 11
: int64 110
: int64 121
> exec
| /*
| Original SQL:
| DELETE FROM limits WHERE id IN (1, 2) ORDER BY id DESC LIMIT 1;
|
| DELETE FROM limits WHERE id IN (3, 4) ORDER BY id DESC LIMIT 1;
|
| DELETE FROM limits WHERE id IN (5, 6) ORDER BY id DESC LIMIT 1;
|
| DELETE FROM limits WHERE id IN (7, 8) ORDER BY id DESC LIMIT 1;
|
| DELETE FROM limits WHERE id IN (9, 10) ORDER BY id DESC LIMIT 1;
|
| DELETE FROM limits WHERE id IN (11, 12) ORDER BY id DESC LIMIT 1;
| */
| INSERT INTO `test`.`limits` (`id`, `a`, `b`) SELECT `id`, `a`, `b` FROM `test`.`bak_limits_test`
< ok 6 0
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
< row
: int64 5
: int64 5
: int64 50
: int64 55
< row
: int64 6
: int64 6
: int64 60
: int64 66
< row
: int64 7
: int64 7
: int64 70
: int64 77
< row
: int64 8
: int64 8
: int64 80
: int64 88
< row
: int64 9
: int64 9
: int64 90
: int64 99
< row
: int64 10
: int64 10
: int64 100
: int64 110
< row
: int64 11
: int64 11
: int64 110
: int64 121
< row
: int64 12
: int64 12
: int64 120
: int64 132
> exec
| CREATE TABLE update_rows (id INT PRIMARY KEY, a INT, g INT AS (a * 2) STORED)
< ok 0 0
> exec
| INSERT INTO update_rows (id, a) VALUES (1, 1), (2, 2), (3, 3)
< ok 3 0
> query
| SELECT * FROM `update_rows` ORDER BY 1, 2
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
| test
: string
| update_rows
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
| `a` * 2
> query
| SELECT INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY INDEX_NAME, SEQ_IN_INDEX
: string
| test
: string
| update_rows
< result
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
< row
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> exec
| CREATE TABLE `test`.`bak_0_update_rows` LIKE `test`.`update_rows`
< ok 0 0
> exec
|
| INSERT INTO `test`.`bak_0_update_rows` (`id`,`a`) SELECT `update_rows`.`id`,`update_rows`.`a` FROM update_rows WHERE id > 1 ORDER BY id LIMIT 1
< ok 1 0
> exec
| UPDATE update_rows SET a = a + 100 WHERE id > 1 ORDER BY id LIMIT 1
< ok 1 0
> query
| SELECT * FROM `update_rows` ORDER BY 1, 2
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
: int64 102
: int64 204
< row
: int64 3
: int64 3
: int64 6
> exec
| /*
| Original SQL:
| UPDATE update_rows SET a = a + 100 WHERE id > 1 ORDER BY id LIMIT 1;
| */
| INSERT INTO `test`.`update_rows` (`id`, `a`) SELECT `id`, `a` FROM `test`.`bak_0_update_rows` ON DUPLICATE KEY UPDATE `a` = VALUES(`a`)
< ok 2 0
> query
| SELECT * FROM `update_rows` ORDER BY 1, 2
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
| test
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
| test
: string
| comments
< result
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
< row
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> exec
| CREATE TABLE `test`.`bak_0_comments` LIKE `test`.`comments`
< ok 0 0
> exec
|
| INSERT INTO `test`.`bak_0_comments` SELECT `comments`.* FROM comments WHERE note = '*/'
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
| INSERT INTO `test`.`comments` (`id`, `note`) SELECT `id`, `note` FROM `test`.`bak_0_comments`
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
| test
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
| test
: string
| cte_alias
< result
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
< row
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> exec
| CREATE TABLE `test`.`bak_0_cte_alias` LIKE `test`.`cte_alias`
< ok 0 0
> exec
|
| INSERT INTO `test`.`bak_0_cte_alias` WITH c AS (SELECT 1 AS id) SELECT `c`.* FROM cte_alias AS c WHERE c.id = 1
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
| INSERT INTO `test`.`cte_alias` (`id`, `a`) SELECT `id`, `a` FROM `test`.`bak_0_cte_alias`
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
| test
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
| `a` + 1
> query
| SELECT INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY INDEX_NAME, SEQ_IN_INDEX
: string
| test
: string
| generated_key
< result
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
< row
: bytes
| uk_g
: bytes
| 0
: bytes
| g
> exec
| CREATE TABLE `test`.`bak_0_generated_key` LIKE `test`.`generated_key`
< ok 0 0
> exec
|
| INSERT INTO `test`.`bak_0_generated_key` (`id`,`a`) SELECT `generated_key`.`id`,`generated_key`.`a` FROM generated_key WHERE id = 1
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
| test
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
| test
: string
| expression_key
< result
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
< row
: bytes
| uk_a
: bytes
| 0
: bytes
| NULL
> exec
| CREATE TABLE `test`.`bak_0_expression_key` LIKE `test`.`expression_key`
< ok 0 0
> exec
|
| INSERT INTO `test`.`bak_0_expression_key` SELECT `expression_key`.* FROM expression_key WHERE id = 1
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
