-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

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
| omni_testbackuprestoreroundtrip
: string
| limits
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
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
| SELECT INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
: string
| limits
< result
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
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
| CREATE TABLE `omni_testbackuprestoreroundtrip`.`bak_limits_omni_testbackuprestoreroundtrip` LIKE `omni_testbackuprestoreroundtrip`.`limits`
< ok 0 0
> exec
|
| INSERT INTO `omni_testbackuprestoreroundtrip`.`bak_limits_omni_testbackuprestoreroundtrip` (`id`,`a`,`b`)  (SELECT `limits`.`id`,`limits`.`a`,`limits`.`b` FROM limits WHERE id = 1)
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
| INSERT INTO `omni_testbackuprestoreroundtrip`.`limits` (`id`, `a`, `b`) SELECT `id`, `a`, `b` FROM `omni_testbackuprestoreroundtrip`.`bak_limits_omni_testbackuprestoreroundtrip`
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
| omni_testbackuprestoreroundtrip
: string
| update_limits
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
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
| SELECT INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
: string
| update_limits
< result
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| PRIMARY
: int64 0
: bytes
| id
> exec
| CREATE TABLE `omni_testbackuprestoreroundtrip`.`bak_update_limits_omni_testbackuprestoreroundtrip` LIKE `omni_testbackuprestoreroundtrip`.`update_limits`
< ok 0 0
> exec
|
| INSERT INTO `omni_testbackuprestoreroundtrip`.`bak_update_limits_omni_testbackuprestoreroundtrip` (`id`,`a`)  (SELECT `update_limits`.`id`,`update_limits`.`a` FROM update_limits WHERE id = 1)
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
| INSERT INTO `omni_testbackuprestoreroundtrip`.`update_limits` (`id`, `a`) SELECT `id`, `a` FROM `omni_testbackuprestoreroundtrip`.`bak_update_limits_omni_testbackuprestoreroundtrip` ON DUPLICATE KEY UPDATE `a` = VALUES(`a`)
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
| omni_testbackuprestoreroundtrip
: string
| self_join
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
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
| SELECT INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
: string
| self_join
< result
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| PRIMARY
: int64 0
: bytes
| id
> exec
| CREATE TABLE `omni_testbackuprestoreroundtrip`.`bak_self_join_omni_testbackuprestoreroundtrip` LIKE `omni_testbackuprestoreroundtrip`.`self_join`
< ok 0 0
> exec
|
| INSERT INTO `omni_testbackuprestoreroundtrip`.`bak_self_join_omni_testbackuprestoreroundtrip`  (SELECT `x`.* FROM self_join x JOIN self_join y ON x.id = y.id + 1)
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
| INSERT INTO `omni_testbackuprestoreroundtrip`.`self_join` (`id`, `a`, `b`) SELECT `id`, `a`, `b` FROM `omni_testbackuprestoreroundtrip`.`bak_self_join_omni_testbackuprestoreroundtrip` ON DUPLICATE KEY UPDATE `a` = VALUES(`a`), `b` = VALUES(`b`)
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
| omni_testbackuprestoreroundtrip
: string
| comments
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
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
| SELECT INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
: string
| comments
< result
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| PRIMARY
: int64 0
: bytes
| id
> exec
| CREATE TABLE `omni_testbackuprestoreroundtrip`.`bak_comments_omni_testbackuprestoreroundtrip` LIKE `omni_testbackuprestoreroundtrip`.`comments`
< ok 0 0
> exec
|
| INSERT INTO `omni_testbackuprestoreroundtrip`.`bak_comments_omni_testbackuprestoreroundtrip`  SELECT `comments`.* FROM comments WHERE note = '*/'
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
| INSERT INTO `omni_testbackuprestoreroundtrip`.`comments` (`id`, `note`) SELECT `id`, `note` FROM `omni_testbackuprestoreroundtrip`.`bak_comments_omni_testbackuprestoreroundtrip`
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
| omni_testbackuprestoreroundtrip
: string
| generated_key
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
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
| SELECT INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
: string
| generated_key
< result
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| uk_g
: int64 0
: bytes
| g
> exec
| CREATE TABLE `omni_testbackuprestoreroundtrip`.`bak_generated_key_omni_testbackuprestoreroundtrip` LIKE `omni_testbackuprestoreroundtrip`.`generated_key`
< ok 0 0
> exec
|
| INSERT INTO `omni_testbackuprestoreroundtrip`.`bak_generated_key_omni_testbackuprestoreroundtrip` (`id`,`a`)  SELECT `generated_key`.`id`,`generated_key`.`a` FROM generated_key WHERE id = 1
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
| CREATE TABLE returning_rows (id INT PRIMARY KEY, a INT)
< ok 0 0
> exec
| INSERT INTO returning_rows VALUES (1, 1), (2, 2)
< ok 2 0
> query
| SELECT * FROM `returning_rows` ORDER BY 1, 2
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
| omni_testbackuprestoreroundtrip
: string
| returning_rows
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
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
| SELECT INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
: string
| returning_rows
< result
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| PRIMARY
: int64 0
: bytes
| id
> exec
| CREATE TABLE `omni_testbackuprestoreroundtrip`.`bak_returning_rows_omni_testbackuprestoreroundtrip` LIKE `omni_testbackuprestoreroundtrip`.`returning_rows`
< ok 0 0
> exec
|
| INSERT INTO `omni_testbackuprestoreroundtrip`.`bak_returning_rows_omni_testbackuprestoreroundtrip`  SELECT `returning_rows`.* FROM returning_rows WHERE id = 1
< ok 1 0
> exec
| DELETE FROM returning_rows WHERE id = 1 RETURNING id
< ok 0 0
> query
| SELECT * FROM `returning_rows` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< row
: int64 2
: int64 2
> exec
| /*
| Original SQL:
| DELETE FROM returning_rows WHERE id = 1 RETURNING id;
| */
| INSERT INTO `omni_testbackuprestoreroundtrip`.`returning_rows` (`id`, `a`) SELECT `id`, `a` FROM `omni_testbackuprestoreroundtrip`.`bak_returning_rows_omni_testbackuprestoreroundtrip`
< ok 1 0
> query
| SELECT * FROM `returning_rows` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< row
: int64 1
: int64 1
< row
: int64 2
: int64 2
