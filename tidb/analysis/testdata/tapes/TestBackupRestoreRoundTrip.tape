-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestBackupRestoreRoundTrip default#1
> exec
| CREATE TABLE limits (id INT PRIMARY KEY, a INT, b INT, g INT AS (a + b), UNIQUE KEY uk_b (b))
< ok 0 0
> exec
| INSERT INTO limits (id, a, b) VALUES (1, 1, 10), (2, 2, 20), (3, 3, 30), (4, 4, 40), (5, 5, 50), (6, 6, 60), (7, 7, 70), (8, 8, 80)
< ok 8 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| limits
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
< row
: bytes
| b
< row
: bytes
| g
> query
| SELECT `id`, `a`, `b`, `g` FROM `test`.`limits` ORDER BY 1, 2, 3, 4
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
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
< row
: bytes
| limits
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| limits
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
| INSERT INTO `test`.`bak_limits_test` (`id`,`a`,`b`)  (SELECT `limits`.`id`,`limits`.`a`,`limits`.`b` FROM limits WHERE id = 1)
|   UNION DISTINCT
|   (SELECT `limits`.`id`,`limits`.`a`,`limits`.`b` FROM limits WHERE id = 2)
|   UNION DISTINCT
|   (SELECT `limits`.`id`,`limits`.`a`,`limits`.`b` FROM limits WHERE id = 3)
|   UNION DISTINCT
|   (SELECT `limits`.`id`,`limits`.`a`,`limits`.`b` FROM limits WHERE id = 4)
|   UNION DISTINCT
|   (SELECT `limits`.`id`,`limits`.`a`,`limits`.`b` FROM limits WHERE id = 5)
|   UNION DISTINCT
|   (SELECT `limits`.`id`,`limits`.`a`,`limits`.`b` FROM limits WHERE id > 5 ORDER BY id DESC LIMIT 1)
< ok 6 0
> exec
| DELETE FROM limits WHERE id > 5 ORDER BY id DESC LIMIT 1
< ok 1 0
> exec
|
| DELETE FROM limits WHERE id = 1
< ok 1 0
> exec
|
| DELETE FROM limits WHERE id = 2
< ok 1 0
> exec
|
| DELETE FROM limits WHERE id = 3
< ok 1 0
> exec
|
| DELETE FROM limits WHERE id = 4
< ok 1 0
> exec
|
| DELETE FROM limits WHERE id = 5
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| limits
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
< row
: bytes
| b
< row
: bytes
| g
> query
| SELECT `id`, `a`, `b`, `g` FROM `test`.`limits` ORDER BY 1, 2, 3, 4
< result
< column "id" "INT"
< column "a" "INT"
< column "b" "INT"
< column "g" "INT"
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
> exec
| /*
| Original SQL:
| DELETE FROM limits WHERE id > 5 ORDER BY id DESC LIMIT 1;
|
| DELETE FROM limits WHERE id = 1;
|
| DELETE FROM limits WHERE id = 2;
|
| DELETE FROM limits WHERE id = 3;
|
| DELETE FROM limits WHERE id = 4;
|
| DELETE FROM limits WHERE id = 5;
| */
| INSERT INTO `test`.`limits` (`id`, `a`, `b`) SELECT `id`, `a`, `b` FROM `test`.`bak_limits_test`
< ok 6 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| limits
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
< row
: bytes
| b
< row
: bytes
| g
> query
| SELECT `id`, `a`, `b`, `g` FROM `test`.`limits` ORDER BY 1, 2, 3, 4
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
> exec
| CREATE TABLE overlapping_limits (id INT PRIMARY KEY)
< ok 0 0
> exec
| INSERT INTO overlapping_limits VALUES (1), (2)
< ok 2 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| overlapping_limits
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
> query
| SELECT `id` FROM `test`.`overlapping_limits` ORDER BY 1
< result
< column "id" "INT"
< row
: int64 1
< row
: int64 2
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
< row
: bytes
| bak_limits_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_limits_test
: bytes
| uk_b
: bytes
| 0
: bytes
| b
< row
: bytes
| limits
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: bytes
| 0
: bytes
| b
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> exec
| CREATE TABLE overlapping_rows (id INT PRIMARY KEY, a INT)
< ok 0 0
> exec
| INSERT INTO overlapping_rows VALUES (1, 1), (2, 2), (3, 3), (4, 4), (5, 5), (6, 6), (7, 7), (8, 8)
< ok 8 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| overlapping_rows
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
> query
| SELECT `id`, `a` FROM `test`.`overlapping_rows` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< row
: int64 1
: int64 1
< row
: int64 2
: int64 2
< row
: int64 3
: int64 3
< row
: int64 4
: int64 4
< row
: int64 5
: int64 5
< row
: int64 6
: int64 6
< row
: int64 7
: int64 7
< row
: int64 8
: int64 8
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
< row
: bytes
| bak_limits_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_limits_test
: bytes
| uk_b
: bytes
| 0
: bytes
| b
< row
: bytes
| limits
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: bytes
| 0
: bytes
| b
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| overlapping_rows
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> exec
| CREATE TABLE `test`.`bak_overlapping_rows_test` LIKE `test`.`overlapping_rows`
< ok 0 0
> exec
|
| INSERT INTO `test`.`bak_overlapping_rows_test` (`id`,`a`)  (SELECT `overlapping_rows`.`id`,`overlapping_rows`.`a` FROM overlapping_rows WHERE id IN (1, 2))
|   UNION DISTINCT
|   (SELECT `overlapping_rows`.`id`,`overlapping_rows`.`a` FROM overlapping_rows WHERE id IN (2, 3))
|   UNION DISTINCT
|   (SELECT `overlapping_rows`.`id`,`overlapping_rows`.`a` FROM overlapping_rows WHERE id IN (3, 4))
|   UNION DISTINCT
|   (SELECT `overlapping_rows`.`id`,`overlapping_rows`.`a` FROM overlapping_rows WHERE id IN (4, 5))
|   UNION DISTINCT
|   (SELECT `overlapping_rows`.`id`,`overlapping_rows`.`a` FROM overlapping_rows WHERE id IN (5, 6))
|   UNION DISTINCT
|   (SELECT `overlapping_rows`.`id`,`overlapping_rows`.`a` FROM overlapping_rows WHERE id IN (6, 7))
< ok 7 0
> exec
| DELETE FROM overlapping_rows WHERE id IN (1, 2)
< ok 2 0
> exec
|
| DELETE FROM overlapping_rows WHERE id IN (2, 3)
< ok 1 0
> exec
|
| DELETE FROM overlapping_rows WHERE id IN (3, 4)
< ok 1 0
> exec
|
| DELETE FROM overlapping_rows WHERE id IN (4, 5)
< ok 1 0
> exec
|
| DELETE FROM overlapping_rows WHERE id IN (5, 6)
< ok 1 0
> exec
|
| DELETE FROM overlapping_rows WHERE id IN (6, 7)
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| overlapping_rows
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
> query
| SELECT `id`, `a` FROM `test`.`overlapping_rows` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< row
: int64 8
: int64 8
> exec
| /*
| Original SQL:
| DELETE FROM overlapping_rows WHERE id IN (1, 2);
|
| DELETE FROM overlapping_rows WHERE id IN (2, 3);
|
| DELETE FROM overlapping_rows WHERE id IN (3, 4);
|
| DELETE FROM overlapping_rows WHERE id IN (4, 5);
|
| DELETE FROM overlapping_rows WHERE id IN (5, 6);
|
| DELETE FROM overlapping_rows WHERE id IN (6, 7);
| */
| INSERT INTO `test`.`overlapping_rows` (`id`, `a`) SELECT `id`, `a` FROM `test`.`bak_overlapping_rows_test`
< ok 7 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| overlapping_rows
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
> query
| SELECT `id`, `a` FROM `test`.`overlapping_rows` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< row
: int64 1
: int64 1
< row
: int64 2
: int64 2
< row
: int64 3
: int64 3
< row
: int64 4
: int64 4
< row
: int64 5
: int64 5
< row
: int64 6
: int64 6
< row
: int64 7
: int64 7
< row
: int64 8
: int64 8
> exec
| CREATE TABLE updates (id INT PRIMARY KEY, a INT, b INT, g INT AS (a * 2) STORED)
< ok 0 0
> exec
| INSERT INTO updates (id, a, b) VALUES (1, 1, 1), (2, 2, 2), (3, 3, 3), (4, 4, 4), (5, 5, 5), (6, 6, 6)
< ok 6 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| updates
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
< row
: bytes
| b
< row
: bytes
| g
> query
| SELECT `id`, `a`, `b`, `g` FROM `test`.`updates` ORDER BY 1, 2, 3, 4
< result
< column "id" "INT"
< column "a" "INT"
< column "b" "INT"
< column "g" "INT"
< row
: int64 1
: int64 1
: int64 1
: int64 2
< row
: int64 2
: int64 2
: int64 2
: int64 4
< row
: int64 3
: int64 3
: int64 3
: int64 6
< row
: int64 4
: int64 4
: int64 4
: int64 8
< row
: int64 5
: int64 5
: int64 5
: int64 10
< row
: int64 6
: int64 6
: int64 6
: int64 12
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
< row
: bytes
| bak_limits_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_limits_test
: bytes
| uk_b
: bytes
| 0
: bytes
| b
< row
: bytes
| bak_overlapping_rows_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: bytes
| 0
: bytes
| b
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| overlapping_rows
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> exec
| CREATE TABLE `test`.`bak_updates_test` LIKE `test`.`updates`
< ok 0 0
> exec
|
| INSERT INTO `test`.`bak_updates_test` (`id`,`a`,`b`)  (SELECT `updates`.`id`,`updates`.`a`,`updates`.`b` FROM updates WHERE id = 1)
|   UNION DISTINCT
|   (SELECT `updates`.`id`,`updates`.`a`,`updates`.`b` FROM updates WHERE id = 2)
|   UNION DISTINCT
|   (SELECT `updates`.`id`,`updates`.`a`,`updates`.`b` FROM updates WHERE id = 3)
|   UNION DISTINCT
|   (SELECT `updates`.`id`,`updates`.`a`,`updates`.`b` FROM updates WHERE id = 4)
|   UNION DISTINCT
|   (SELECT `updates`.`id`,`updates`.`a`,`updates`.`b` FROM updates WHERE id = 5)
|   UNION DISTINCT
|   (SELECT `updates`.`id`,`updates`.`a`,`updates`.`b` FROM updates WHERE id = 6)
< ok 6 0
> exec
| UPDATE updates SET a = 10 WHERE id = 1
< ok 1 0
> exec
|
| UPDATE updates SET b = 20 WHERE id = 2
< ok 1 0
> exec
|
| UPDATE updates SET a = a + 100 WHERE id = 3
< ok 1 0
> exec
|
| UPDATE updates SET a = a + 100 WHERE id = 4
< ok 1 0
> exec
|
| UPDATE updates SET a = a + 100 WHERE id = 5
< ok 1 0
> exec
|
| UPDATE updates SET a = a + 100 WHERE id = 6
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| updates
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
< row
: bytes
| b
< row
: bytes
| g
> query
| SELECT `id`, `a`, `b`, `g` FROM `test`.`updates` ORDER BY 1, 2, 3, 4
< result
< column "id" "INT"
< column "a" "INT"
< column "b" "INT"
< column "g" "INT"
< row
: int64 1
: int64 10
: int64 1
: int64 20
< row
: int64 2
: int64 2
: int64 20
: int64 4
< row
: int64 3
: int64 103
: int64 3
: int64 206
< row
: int64 4
: int64 104
: int64 4
: int64 208
< row
: int64 5
: int64 105
: int64 5
: int64 210
< row
: int64 6
: int64 106
: int64 6
: int64 212
> exec
| /*
| Original SQL:
| UPDATE updates SET a = 10 WHERE id = 1;
|
| UPDATE updates SET b = 20 WHERE id = 2;
|
| UPDATE updates SET a = a + 100 WHERE id = 3;
|
| UPDATE updates SET a = a + 100 WHERE id = 4;
|
| UPDATE updates SET a = a + 100 WHERE id = 5;
|
| UPDATE updates SET a = a + 100 WHERE id = 6;
| */
| INSERT INTO `test`.`updates` (`id`, `a`, `b`) SELECT `id`, `a`, `b` FROM `test`.`bak_updates_test` ON DUPLICATE KEY UPDATE `a` = VALUES(`a`), `b` = VALUES(`b`)
< ok 12 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| updates
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
< row
: bytes
| b
< row
: bytes
| g
> query
| SELECT `id`, `a`, `b`, `g` FROM `test`.`updates` ORDER BY 1, 2, 3, 4
< result
< column "id" "INT"
< column "a" "INT"
< column "b" "INT"
< column "g" "INT"
< row
: int64 1
: int64 1
: int64 1
: int64 2
< row
: int64 2
: int64 2
: int64 2
: int64 4
< row
: int64 3
: int64 3
: int64 3
: int64 6
< row
: int64 4
: int64 4
: int64 4
: int64 8
< row
: int64 5
: int64 5
: int64 5
: int64 10
< row
: int64 6
: int64 6
: int64 6
: int64 12
> exec
| CREATE TABLE one_to_many (id INT PRIMARY KEY, a INT)
< ok 0 0
> exec
| CREATE TABLE children (id INT PRIMARY KEY, parent INT)
< ok 0 0
> exec
| INSERT INTO one_to_many VALUES (1, 1), (2, 2)
< ok 2 0
> exec
| INSERT INTO children VALUES (1, 1), (2, 1), (3, 2)
< ok 3 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| one_to_many
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
> query
| SELECT `id`, `a` FROM `test`.`one_to_many` ORDER BY 1, 2
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
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| children
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| parent
> query
| SELECT `id`, `parent` FROM `test`.`children` ORDER BY 1, 2
< result
< column "id" "INT"
< column "parent" "INT"
< row
: int64 1
: int64 1
< row
: int64 2
: int64 1
< row
: int64 3
: int64 2
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
< row
: bytes
| bak_limits_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_limits_test
: bytes
| uk_b
: bytes
| 0
: bytes
| b
< row
: bytes
| bak_overlapping_rows_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_updates_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: bytes
| 0
: bytes
| b
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| overlapping_rows
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> exec
| CREATE TABLE `test`.`bak_0_one_to_many` LIKE `test`.`one_to_many`
< ok 0 0
> exec
|
| INSERT INTO `test`.`bak_0_one_to_many` (`id`,`a`) SELECT DISTINCT `p`.`id`,`p`.`a` FROM one_to_many p JOIN children c ON c.parent = p.id WHERE p.id = 1
< ok 1 0
> exec
| DELETE p FROM one_to_many p JOIN children c ON c.parent = p.id WHERE p.id = 1
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| one_to_many
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
> query
| SELECT `id`, `a` FROM `test`.`one_to_many` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< row
: int64 2
: int64 2
> exec
| /*
| Original SQL:
| DELETE p FROM one_to_many p JOIN children c ON c.parent = p.id WHERE p.id = 1;
| */
| INSERT INTO `test`.`one_to_many` (`id`, `a`) SELECT `id`, `a` FROM `test`.`bak_0_one_to_many`
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| one_to_many
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
> query
| SELECT `id`, `a` FROM `test`.`one_to_many` ORDER BY 1, 2
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
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| children
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| parent
> query
| SELECT `id`, `parent` FROM `test`.`children` ORDER BY 1, 2
< result
< column "id" "INT"
< column "parent" "INT"
< row
: int64 1
: int64 1
< row
: int64 2
: int64 1
< row
: int64 3
: int64 2
> exec
| CREATE TABLE keyless_join (id INT, a INT)
< ok 0 0
> exec
| CREATE TABLE keyless_children (id INT PRIMARY KEY, parent INT)
< ok 0 0
> exec
| INSERT INTO keyless_join VALUES (1, 1), (1, 1)
< ok 2 0
> exec
| INSERT INTO keyless_children VALUES (1, 1), (2, 1)
< ok 2 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| keyless_join
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
> query
| SELECT `id`, `a` FROM `test`.`keyless_join` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< row
: int64 1
: int64 1
< row
: int64 1
: int64 1
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
< row
: bytes
| bak_0_one_to_many
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_limits_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_limits_test
: bytes
| uk_b
: bytes
| 0
: bytes
| b
< row
: bytes
| bak_overlapping_rows_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_updates_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| keyless_children
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: bytes
| 0
: bytes
| b
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| overlapping_rows
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> exec
| CREATE TABLE keyless_rows (id INT, a INT)
< ok 0 0
> exec
| INSERT INTO keyless_rows VALUES (1, 1), (1, 1), (2, 2)
< ok 3 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| keyless_rows
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
> query
| SELECT `id`, `a` FROM `test`.`keyless_rows` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< row
: int64 1
: int64 1
< row
: int64 1
: int64 1
< row
: int64 2
: int64 2
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
< row
: bytes
| bak_0_one_to_many
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_limits_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_limits_test
: bytes
| uk_b
: bytes
| 0
: bytes
| b
< row
: bytes
| bak_overlapping_rows_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_updates_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| keyless_children
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: bytes
| 0
: bytes
| b
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| overlapping_rows
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> exec
| CREATE TABLE `test`.`bak_0_keyless_rows` LIKE `test`.`keyless_rows`
< ok 0 0
> exec
|
| INSERT INTO `test`.`bak_0_keyless_rows` (`id`,`a`) SELECT `keyless_rows`.`id`,`keyless_rows`.`a` FROM keyless_rows WHERE id = 1
< ok 2 0
> exec
| DELETE FROM keyless_rows WHERE id = 1
< ok 2 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| keyless_rows
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
> query
| SELECT `id`, `a` FROM `test`.`keyless_rows` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< row
: int64 2
: int64 2
> exec
| /*
| Original SQL:
| DELETE FROM keyless_rows WHERE id = 1;
| */
| INSERT INTO `test`.`keyless_rows` (`id`, `a`) SELECT `id`, `a` FROM `test`.`bak_0_keyless_rows`
< ok 2 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| keyless_rows
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
> query
| SELECT `id`, `a` FROM `test`.`keyless_rows` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< row
: int64 1
: int64 1
< row
: int64 1
: int64 1
< row
: int64 2
: int64 2
> exec
| CREATE TABLE parens (id INT PRIMARY KEY, a INT)
< ok 0 0
> exec
| CREATE TABLE paren_peers (id INT PRIMARY KEY)
< ok 0 0
> exec
| INSERT INTO parens VALUES (1, 1), (2, 2)
< ok 2 0
> exec
| INSERT INTO paren_peers VALUES (1), (2)
< ok 2 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| parens
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
> query
| SELECT `id`, `a` FROM `test`.`parens` ORDER BY 1, 2
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
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
< row
: bytes
| bak_0_one_to_many
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_limits_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_limits_test
: bytes
| uk_b
: bytes
| 0
: bytes
| b
< row
: bytes
| bak_overlapping_rows_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_updates_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| keyless_children
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: bytes
| 0
: bytes
| b
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| overlapping_rows
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| paren_peers
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| parens
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> exec
| CREATE TABLE `test`.`bak_0_parens` LIKE `test`.`parens`
< ok 0 0
> exec
|
| INSERT INTO `test`.`bak_0_parens` (`id`,`a`) SELECT DISTINCT `parens`.`id`,`parens`.`a` FROM (parens JOIN paren_peers ON parens.id = paren_peers.id) JOIN paren_peers AS p ON p.id = parens.id WHERE parens.id = 1
< ok 1 0
> exec
| UPDATE (parens JOIN paren_peers ON parens.id = paren_peers.id) JOIN paren_peers AS p ON p.id = parens.id SET parens.a = 9 WHERE parens.id = 1
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| parens
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
> query
| SELECT `id`, `a` FROM `test`.`parens` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< row
: int64 1
: int64 9
< row
: int64 2
: int64 2
> exec
| /*
| Original SQL:
| UPDATE (parens JOIN paren_peers ON parens.id = paren_peers.id) JOIN paren_peers AS p ON p.id = parens.id SET parens.a = 9 WHERE parens.id = 1;
| */
| INSERT INTO `test`.`parens` (`id`, `a`) SELECT `id`, `a` FROM `test`.`bak_0_parens` ON DUPLICATE KEY UPDATE `a` = VALUES(`a`)
< ok 2 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| parens
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
> query
| SELECT `id`, `a` FROM `test`.`parens` ORDER BY 1, 2
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
| CREATE TABLE `back``tick` (id INT PRIMARY KEY, `c``ol` INT)
< ok 0 0
> exec
| INSERT INTO `back``tick` VALUES (1, 1), (2, 2)
< ok 2 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| back`tick
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| c`ol
> query
| SELECT `id`, `c``ol` FROM `test`.`back``tick` ORDER BY 1, 2
< result
< column "id" "INT"
< column "c`ol" "INT"
< row
: int64 1
: int64 1
< row
: int64 2
: int64 2
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_one_to_many
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_parens
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_limits_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_limits_test
: bytes
| uk_b
: bytes
| 0
: bytes
| b
< row
: bytes
| bak_overlapping_rows_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_updates_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| keyless_children
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: bytes
| 0
: bytes
| b
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| overlapping_rows
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| paren_peers
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| parens
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> exec
| CREATE TABLE `test`.`bak_0_back``tick` LIKE `test`.`back``tick`
< ok 0 0
> exec
|
| INSERT INTO `test`.`bak_0_back``tick` (`id`,`c``ol`) SELECT `back``tick`.`id`,`back``tick`.`c``ol` FROM `back``tick` WHERE `c``ol` = 1
< ok 1 0
> exec
| DELETE FROM `back``tick` WHERE `c``ol` = 1
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| back`tick
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| c`ol
> query
| SELECT `id`, `c``ol` FROM `test`.`back``tick` ORDER BY 1, 2
< result
< column "id" "INT"
< column "c`ol" "INT"
< row
: int64 2
: int64 2
> exec
| /*
| Original SQL:
| DELETE FROM `back``tick` WHERE `c``ol` = 1;
| */
| INSERT INTO `test`.`back``tick` (`id`, `c``ol`) SELECT `id`, `c``ol` FROM `test`.`bak_0_back``tick`
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| back`tick
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| c`ol
> query
| SELECT `id`, `c``ol` FROM `test`.`back``tick` ORDER BY 1, 2
< result
< column "id" "INT"
< column "c`ol" "INT"
< row
: int64 1
: int64 1
< row
: int64 2
: int64 2
> exec
| CREATE TABLE comments (id INT PRIMARY KEY, note VARCHAR(10))
< ok 0 0
> exec
| INSERT INTO comments VALUES (1, '*/'), (2, 'x')
< ok 2 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| comments
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| note
> query
| SELECT `id`, `note` FROM `test`.`comments` ORDER BY 1, 2
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
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_back`tick
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_one_to_many
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_parens
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_limits_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_limits_test
: bytes
| uk_b
: bytes
| 0
: bytes
| b
< row
: bytes
| bak_overlapping_rows_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_updates_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| comments
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| keyless_children
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: bytes
| 0
: bytes
| b
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| overlapping_rows
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| paren_peers
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| parens
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| updates
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
| INSERT INTO `test`.`bak_0_comments` (`id`,`note`) SELECT `comments`.`id`,`comments`.`note` FROM comments WHERE note = '*/'
< ok 1 0
> exec
| DELETE /* note */ FROM comments WHERE note = '*/'
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| comments
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| note
> query
| SELECT `id`, `note` FROM `test`.`comments` ORDER BY 1, 2
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
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| comments
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| note
> query
| SELECT `id`, `note` FROM `test`.`comments` ORDER BY 1, 2
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
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| cte_alias
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
> query
| SELECT `id`, `a` FROM `test`.`cte_alias` ORDER BY 1, 2
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
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_back`tick
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_comments
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_one_to_many
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_parens
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_limits_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_limits_test
: bytes
| uk_b
: bytes
| 0
: bytes
| b
< row
: bytes
| bak_overlapping_rows_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_updates_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| comments
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| cte_alias
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| keyless_children
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: bytes
| 0
: bytes
| b
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| overlapping_rows
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| paren_peers
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| parens
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| updates
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
| INSERT INTO `test`.`bak_0_cte_alias` (`id`,`a`) WITH c AS (SELECT 1 AS id) SELECT `c`.`id`,`c`.`a` FROM cte_alias AS c WHERE c.id = 1
< ok 1 0
> exec
| WITH c AS (SELECT 1 AS id) DELETE c FROM cte_alias AS c WHERE c.id = 1
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| cte_alias
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
> query
| SELECT `id`, `a` FROM `test`.`cte_alias` ORDER BY 1, 2
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
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| cte_alias
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
> query
| SELECT `id`, `a` FROM `test`.`cte_alias` ORDER BY 1, 2
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
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| generated_key
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
< row
: bytes
| g
> query
| SELECT `id`, `a`, `g` FROM `test`.`generated_key` ORDER BY 1, 2, 3
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
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_back`tick
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_comments
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_cte_alias
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_one_to_many
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_parens
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_limits_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_limits_test
: bytes
| uk_b
: bytes
| 0
: bytes
| b
< row
: bytes
| bak_overlapping_rows_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_updates_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| comments
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| cte_alias
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| generated_key
: bytes
| uk_g
: bytes
| 0
: bytes
| g
< row
: bytes
| keyless_children
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: bytes
| 0
: bytes
| b
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| overlapping_rows
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| paren_peers
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| parens
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
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
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| generated_key
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
< row
: bytes
| g
> query
| SELECT `id`, `a`, `g` FROM `test`.`generated_key` ORDER BY 1, 2, 3
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
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| expression_key
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
> query
| SELECT `id`, `a` FROM `test`.`expression_key` ORDER BY 1, 2
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
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_generated_key
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_generated_key
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| expression_key
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| expression_key
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_back`tick
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_comments
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_cte_alias
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_generated_key
: bytes
| uk_g
: bytes
| 0
: bytes
| g
< row
: bytes
| bak_0_one_to_many
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_parens
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_limits_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_limits_test
: bytes
| uk_b
: bytes
| 0
: bytes
| b
< row
: bytes
| bak_overlapping_rows_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_updates_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| comments
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| cte_alias
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| expression_key
: bytes
| uk_a
: bytes
| 0
: bytes
| NULL
< row
: bytes
| generated_key
: bytes
| uk_g
: bytes
| 0
: bytes
| g
< row
: bytes
| keyless_children
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: bytes
| 0
: bytes
| b
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| overlapping_rows
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| paren_peers
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| parens
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> exec
| CREATE TABLE `test`.`bak_0_expression_key` LIKE `test`.`expression_key`
< ok 0 0
> exec
|
| INSERT INTO `test`.`bak_0_expression_key` (`id`,`a`) SELECT `expression_key`.`id`,`expression_key`.`a` FROM expression_key WHERE id = 1
< ok 1 0
> exec
| UPDATE expression_key SET a = a + 10 WHERE id = 1
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| expression_key
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
> query
| SELECT `id`, `a` FROM `test`.`expression_key` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< row
: int64 1
: int64 11
< row
: int64 2
: int64 2
> exec
| CREATE TABLE nullable_key (u INT NULL, v INT, UNIQUE KEY uk_u (u))
< ok 0 0
> exec
| INSERT INTO nullable_key VALUES (NULL, 1), (2, 2)
< ok 2 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| nullable_key
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| u
< row
: bytes
| v
> query
| SELECT `u`, `v` FROM `test`.`nullable_key` ORDER BY 1, 2
< result
< column "u" "INT"
< column "v" "INT"
< row
: null
: int64 1
< row
: int64 2
: int64 2
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_expression_key
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_expression_key
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_generated_key
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_generated_key
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| expression_key
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| expression_key
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
< row
: bytes
| nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
< row
: bytes
| nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_back`tick
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_comments
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_cte_alias
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_expression_key
: bytes
| uk_a
: bytes
| 0
: bytes
| NULL
< row
: bytes
| bak_0_generated_key
: bytes
| uk_g
: bytes
| 0
: bytes
| g
< row
: bytes
| bak_0_one_to_many
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_parens
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_limits_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_limits_test
: bytes
| uk_b
: bytes
| 0
: bytes
| b
< row
: bytes
| bak_overlapping_rows_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_updates_test
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| comments
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| cte_alias
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| expression_key
: bytes
| uk_a
: bytes
| 0
: bytes
| NULL
< row
: bytes
| generated_key
: bytes
| uk_g
: bytes
| 0
: bytes
| g
< row
: bytes
| keyless_children
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: bytes
| 0
: bytes
| b
< row
: bytes
| nullable_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| overlapping_rows
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| paren_peers
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| parens
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> exec
| CREATE TABLE `test`.`bak_0_nullable_key` LIKE `test`.`nullable_key`
< ok 0 0
> exec
|
| INSERT INTO `test`.`bak_0_nullable_key` (`u`,`v`) SELECT `nullable_key`.`u`,`nullable_key`.`v` FROM nullable_key WHERE u IS NULL
< ok 1 0
> exec
| UPDATE nullable_key SET v = 3 WHERE u IS NULL
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| nullable_key
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| u
< row
: bytes
| v
> query
| SELECT `u`, `v` FROM `test`.`nullable_key` ORDER BY 1, 2
< result
< column "u" "INT"
< column "v" "INT"
< row
: null
: int64 3
< row
: int64 2
: int64 2
