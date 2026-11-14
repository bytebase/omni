-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestBackupRestoreRoundTrip default#1
> exec
| CREATE TABLE limits (id INT PRIMARY KEY, a INT, b INT, g INT AS (a + b), UNIQUE KEY uk_b (b))
< ok 0 0
> exec
| INSERT INTO limits (id, a, b) VALUES (1, 1, 10), (2, 2, 20), (3, 3, 30), (4, 4, 40), (5, 5, 50), (6, 6, 60), (7, 7, 70), (8, 8, 80)
< ok 8 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| limits
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| VIRTUAL GENERATED
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
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
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
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| limits
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| VIRTUAL GENERATED
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| limits
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| VIRTUAL GENERATED
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| overlapping_limits
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
> query
| SELECT `id` FROM `test`.`overlapping_limits` ORDER BY 1
< result
< column "id" "INT"
< row
: int64 1
< row
: int64 2
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
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
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE overlapping_rows (id INT PRIMARY KEY, a INT)
< ok 0 0
> exec
| INSERT INTO overlapping_rows VALUES (1, 1), (2, 2), (3, 3), (4, 4), (5, 5), (6, 6), (7, 7), (8, 8)
< ok 8 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| overlapping_rows
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
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
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| overlapping_rows
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| overlapping_rows
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| updates
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| STORED GENERATED
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
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
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
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| updates
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| STORED GENERATED
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| updates
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| STORED GENERATED
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| one_to_many
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| children
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
< row
: bytes
| parent
: bytes
|
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
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
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
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| one_to_many
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| one_to_many
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| children
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
< row
: bytes
| parent
: bytes
|
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| keyless_join
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
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
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE keyless_rows (id INT, a INT)
< ok 0 0
> exec
| INSERT INTO keyless_rows VALUES (1, 1), (1, 1), (2, 2)
< ok 3 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| keyless_rows
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
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
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| keyless_rows
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| keyless_rows
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| parens
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
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
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| parens
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| parens
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| back`tick
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
< row
: bytes
| c`ol
: bytes
|
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
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
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
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| back`tick
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
< row
: bytes
| c`ol
: bytes
|
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| back`tick
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
< row
: bytes
| c`ol
: bytes
|
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
| CREATE TABLE single_key (id INT PRIMARY KEY, u INT NOT NULL, UNIQUE KEY uk_u (u))
< ok 0 0
> exec
| INSERT INTO single_key VALUES (1, 10), (2, 20)
< ok 2 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| single_key
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
< row
: bytes
| u
: bytes
|
> query
| SELECT `id`, `u` FROM `test`.`single_key` ORDER BY 1, 2
< result
< column "id" "INT"
< column "u" "INT"
< row
: int64 1
: int64 10
< row
: int64 2
: int64 20
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
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
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE `test`.`bak_0_single_key` LIKE `test`.`single_key`
< ok 0 0
> exec
|
| INSERT INTO `test`.`bak_0_single_key` (`id`,`u`) SELECT `single_key`.`id`,`single_key`.`u` FROM single_key WHERE id = 1
< ok 1 0
> exec
| UPDATE single_key SET u = 30 WHERE id = 1
< ok 1 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| single_key
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
< row
: bytes
| u
: bytes
|
> query
| SELECT `id`, `u` FROM `test`.`single_key` ORDER BY 1, 2
< result
< column "id" "INT"
< column "u" "INT"
< row
: int64 1
: int64 30
< row
: int64 2
: int64 20
> exec
| /*
| Original SQL:
| UPDATE single_key SET u = 30 WHERE id = 1;
| */
| INSERT INTO `test`.`single_key` (`id`, `u`) SELECT `id`, `u` FROM `test`.`bak_0_single_key` ON DUPLICATE KEY UPDATE `u` = VALUES(`u`)
< ok 2 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| single_key
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
< row
: bytes
| u
: bytes
|
> query
| SELECT `id`, `u` FROM `test`.`single_key` ORDER BY 1, 2
< result
< column "id" "INT"
< column "u" "INT"
< row
: int64 1
: int64 10
< row
: int64 2
: int64 20
> exec
| CREATE TABLE comments (id INT PRIMARY KEY, note VARCHAR(10))
< ok 0 0
> exec
| INSERT INTO comments VALUES (1, '*/'), (2, 'x')
< ok 2 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| comments
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
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
| bak_0_single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| comments
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| comments
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| cte_alias
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
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
| bak_0_single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| cte_alias
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| cte_alias
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| generated_key
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| VIRTUAL GENERATED
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
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
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
| bak_0_single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| generated_key
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| VIRTUAL GENERATED
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| expression_key
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
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
| bak_0_single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| expression_key
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| nullable_key
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| u
: bytes
|
< row
: bytes
| v
: bytes
|
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
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
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
| bak_0_single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
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
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| nullable_key
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| u
: bytes
|
< row
: bytes
| v
: bytes
|
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
> exec
| CREATE TABLE stamped (id INT PRIMARY KEY, v INT, ts TIMESTAMP NOT NULL DEFAULT '2000-01-01 00:00:00' ON UPDATE CURRENT_TIMESTAMP)
< ok 0 0
> exec
| INSERT INTO stamped VALUES (1, 1, '2000-01-01 00:00:00'), (2, 2, '2000-01-01 00:00:00')
< ok 2 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| stamped
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
< row
: bytes
| v
: bytes
|
< row
: bytes
| ts
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
> query
| SELECT `id`, `v`, `ts` FROM `test`.`stamped` ORDER BY 1, 2, 3
< result
< column "id" "INT"
< column "v" "INT"
< column "ts" "TIMESTAMP"
< row
: int64 1
: int64 1
: bytes
| 2000-01-01 00:00:00
< row
: int64 2
: int64 2
: bytes
| 2000-01-01 00:00:00
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
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
| bak_0_nullable_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| bak_0_single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| stamped
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
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE `test`.`bak_0_stamped` LIKE `test`.`stamped`
< ok 0 0
> exec
|
| INSERT INTO `test`.`bak_0_stamped` (`id`,`v`,`ts`) SELECT `stamped`.`id`,`stamped`.`v`,`stamped`.`ts` FROM stamped WHERE id = 1
< ok 1 0
> exec
| UPDATE stamped SET v = 5 WHERE id = 1
< ok 1 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| stamped
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
< row
: bytes
| v
: bytes
|
< row
: bytes
| ts
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
> query
| SELECT `id`, `v`, '-' FROM `test`.`stamped` ORDER BY 1, 2, 3
< result
< column "id" "INT"
< column "v" "INT"
< column "-" "VARCHAR"
< row
: int64 1
: int64 5
: bytes
| -
< row
: int64 2
: int64 2
: bytes
| -
> exec
| /*
| Original SQL:
| UPDATE stamped SET v = 5 WHERE id = 1;
| */
| INSERT INTO `test`.`stamped` (`id`, `v`, `ts`) SELECT `id`, `v`, `ts` FROM `test`.`bak_0_stamped` ON DUPLICATE KEY UPDATE `v` = VALUES(`v`), `ts` = VALUES(`ts`)
< ok 2 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| stamped
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
< row
: bytes
| v
: bytes
|
< row
: bytes
| ts
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
> query
| SELECT `id`, `v`, `ts` FROM `test`.`stamped` ORDER BY 1, 2, 3
< result
< column "id" "INT"
< column "v" "INT"
< column "ts" "TIMESTAMP"
< row
: int64 1
: int64 1
: bytes
| 2000-01-01 00:00:00
< row
: int64 2
: int64 2
: bytes
| 2000-01-01 00:00:00
> exec
| CREATE TABLE flags (id INT PRIMARY KEY, enabled INT)
< ok 0 0
> exec
| CREATE TABLE flagged (id INT PRIMARY KEY)
< ok 0 0
> exec
| INSERT INTO flags VALUES (1, 0)
< ok 1 0
> exec
| INSERT INTO flagged VALUES (1)
< ok 1 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| flagged
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
> query
| SELECT `id` FROM `test`.`flagged` ORDER BY 1
< result
< column "id" "INT"
< row
: int64 1
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| flagged
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| enabled
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
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
| bak_0_nullable_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| bak_0_single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_0_stamped
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
| flagged
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flags
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
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| stamped
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
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE fk_parent (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE fk_child (id INT PRIMARY KEY, parent_id INT, FOREIGN KEY fk_child_parent (parent_id) REFERENCES fk_parent (id) ON DELETE CASCADE)
< ok 0 0
> exec
| INSERT INTO fk_parent VALUES (1)
< ok 1 0
> exec
| INSERT INTO fk_child VALUES (1, 1)
< ok 1 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| fk_parent
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
> query
| SELECT `id` FROM `test`.`fk_parent` ORDER BY 1
< result
< column "id" "INT"
< row
: int64 1
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_child
: bytes
| parent_id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_parent
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flagged
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| enabled
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
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
| bak_0_nullable_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| bak_0_single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_0_stamped
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
| fk_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| 1
: bytes
| parent_id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flags
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
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| stamped
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
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| parent_id
: bytes
| test
: bytes
| fk_parent
: bytes
| id
: bytes
| NO ACTION
: bytes
| CASCADE
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE outer_keep (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE outer_gone (id INT PRIMARY KEY, a INT)
< ok 0 0
> exec
| INSERT INTO outer_keep VALUES (1), (2)
< ok 2 0
> exec
| INSERT INTO outer_gone VALUES (1, 1)
< ok 1 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| outer_gone
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT `id`, `a` FROM `test`.`outer_gone` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< row
: int64 1
: int64 1
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_child
: bytes
| parent_id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_parent
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flagged
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| enabled
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_keep
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
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
| bak_0_nullable_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| bak_0_single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_0_stamped
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
| fk_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| 1
: bytes
| parent_id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flags
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
| outer_gone
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| outer_keep
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
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| stamped
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
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| parent_id
: bytes
| test
: bytes
| fk_parent
: bytes
| id
: bytes
| NO ACTION
: bytes
| CASCADE
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE `test`.`bak_0_outer_gone` LIKE `test`.`outer_gone`
< ok 0 0
> exec
|
| INSERT INTO `test`.`bak_0_outer_gone` (`id`,`a`) SELECT DISTINCT `t`.`id`,`t`.`a` FROM outer_keep u LEFT JOIN outer_gone t ON t.id = u.id WHERE `t`.`id` IS NOT NULL
< ok 1 0
> exec
| DELETE t FROM outer_keep u LEFT JOIN outer_gone t ON t.id = u.id
< ok 1 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| outer_gone
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT `id`, `a` FROM `test`.`outer_gone` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
> exec
| /*
| Original SQL:
| DELETE t FROM outer_keep u LEFT JOIN outer_gone t ON t.id = u.id;
| */
| INSERT INTO `test`.`outer_gone` (`id`, `a`) SELECT `id`, `a` FROM `test`.`bak_0_outer_gone`
< ok 1 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| outer_gone
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT `id`, `a` FROM `test`.`outer_gone` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< row
: int64 1
: int64 1
> exec
| CREATE TABLE nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two (id INT PRIMARY KEY)
< ok 0 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
> query
| SELECT `id` FROM `test`.`nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one` ORDER BY 1
< result
< column "id" "INT"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_child
: bytes
| parent_id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_parent
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flagged
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| enabled
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_keep
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
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
| bak_0_nullable_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| bak_0_outer_gone
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
| bak_0_single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_0_stamped
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
| fk_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| 1
: bytes
| parent_id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flags
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
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
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
| outer_gone
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| outer_keep
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
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| stamped
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
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| parent_id
: bytes
| test
: bytes
| fk_parent
: bytes
| id
: bytes
| NO ACTION
: bytes
| CASCADE
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE swapped_key (id INT PRIMARY KEY, u INT NOT NULL, UNIQUE KEY uk_u (u))
< ok 0 0
> exec
| INSERT INTO swapped_key VALUES (1, 10), (2, 20)
< ok 2 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| swapped_key
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
< row
: bytes
| u
: bytes
|
> query
| SELECT `id`, `u` FROM `test`.`swapped_key` ORDER BY 1, 2
< result
< column "id" "INT"
< column "u" "INT"
< row
: int64 1
: int64 10
< row
: int64 2
: int64 20
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_child
: bytes
| parent_id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_parent
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flagged
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| enabled
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_keep
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
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
| bak_0_nullable_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| bak_0_outer_gone
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
| bak_0_single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_0_stamped
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
| fk_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| 1
: bytes
| parent_id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flags
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
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
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
| outer_gone
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| outer_keep
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
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| parent_id
: bytes
| test
: bytes
| fk_parent
: bytes
| id
: bytes
| NO ACTION
: bytes
| CASCADE
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE `test`.`bak_0_swapped_key` LIKE `test`.`swapped_key`
< ok 0 0
> exec
|
| INSERT INTO `test`.`bak_0_swapped_key` (`id`,`u`) SELECT `swapped_key`.`id`,`swapped_key`.`u` FROM swapped_key ORDER BY u DESC
< ok 2 0
> exec
| UPDATE swapped_key SET u = u + 10 ORDER BY u DESC
< ok 2 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| swapped_key
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
< row
: bytes
| u
: bytes
|
> query
| SELECT `id`, `u` FROM `test`.`swapped_key` ORDER BY 1, 2
< result
< column "id" "INT"
< column "u" "INT"
< row
: int64 1
: int64 20
< row
: int64 2
: int64 30
> exec
| CREATE TABLE separate_overlap (id INT PRIMARY KEY, a INT)
< ok 0 0
> exec
| INSERT INTO separate_overlap VALUES (1, 1)
< ok 1 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| separate_overlap
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT `id`, `a` FROM `test`.`separate_overlap` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< row
: int64 1
: int64 1
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| bak_0_swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_child
: bytes
| parent_id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_parent
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flagged
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| enabled
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_keep
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
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
| bak_0_nullable_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| bak_0_outer_gone
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
| bak_0_single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_0_stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| fk_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| 1
: bytes
| parent_id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flags
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
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
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
| outer_gone
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| outer_keep
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
| separate_overlap
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| parent_id
: bytes
| test
: bytes
| fk_parent
: bytes
| id
: bytes
| NO ACTION
: bytes
| CASCADE
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE generated_default (id INT PRIMARY KEY, a INT, g INT AS (a + 1))
< ok 0 0
> exec
| INSERT INTO generated_default (id, a) VALUES (1, 1), (2, 2)
< ok 2 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| generated_default
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| VIRTUAL GENERATED
> query
| SELECT `id`, `a`, `g` FROM `test`.`generated_default` ORDER BY 1, 2, 3
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
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| bak_0_swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_child
: bytes
| parent_id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_parent
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flagged
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| enabled
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_keep
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
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
| bak_0_nullable_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| bak_0_outer_gone
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
| bak_0_single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_0_stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| fk_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| 1
: bytes
| parent_id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| generated_default
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
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
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
| outer_gone
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| outer_keep
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
| separate_overlap
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| parent_id
: bytes
| test
: bytes
| fk_parent
: bytes
| id
: bytes
| NO ACTION
: bytes
| CASCADE
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE `test`.`bak_0_generated_default` LIKE `test`.`generated_default`
< ok 0 0
> exec
|
| INSERT INTO `test`.`bak_0_generated_default` (`id`,`a`) SELECT `generated_default`.`id`,`generated_default`.`a` FROM generated_default WHERE id = 1
< ok 1 0
> exec
| UPDATE generated_default SET generated_default.a = 5, generated_default.g = DEFAULT WHERE id = 1
< ok 1 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| generated_default
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| VIRTUAL GENERATED
> query
| SELECT `id`, `a`, `g` FROM `test`.`generated_default` ORDER BY 1, 2, 3
< result
< column "id" "INT"
< column "a" "INT"
< column "g" "INT"
< row
: int64 1
: int64 5
: int64 6
< row
: int64 2
: int64 2
: int64 3
> exec
| /*
| Original SQL:
| UPDATE generated_default SET generated_default.a = 5, generated_default.g = DEFAULT WHERE id = 1;
| */
| INSERT INTO `test`.`generated_default` (`id`, `a`) SELECT `id`, `a` FROM `test`.`bak_0_generated_default` ON DUPLICATE KEY UPDATE `a` = VALUES(`a`)
< ok 2 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| generated_default
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| VIRTUAL GENERATED
> query
| SELECT `id`, `a`, `g` FROM `test`.`generated_default` ORDER BY 1, 2, 3
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
> exec
| CREATE TABLE user_variable (id INT PRIMARY KEY)
< ok 0 0
> exec
| INSERT INTO user_variable VALUES (1), (2)
< ok 2 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| user_variable
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
> query
| SELECT `id` FROM `test`.`user_variable` ORDER BY 1
< result
< column "id" "INT"
< row
: int64 1
< row
: int64 2
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| bak_0_swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_child
: bytes
| parent_id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_parent
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flagged
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| enabled
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_keep
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| user_variable
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
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
| bak_0_generated_default
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
| bak_0_nullable_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| bak_0_outer_gone
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
| bak_0_single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_0_stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| fk_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| 1
: bytes
| parent_id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| generated_default
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
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
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
| outer_gone
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| outer_keep
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
| separate_overlap
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| user_variable
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| parent_id
: bytes
| test
: bytes
| fk_parent
: bytes
| id
: bytes
| NO ACTION
: bytes
| CASCADE
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE dependent_updates (id INT PRIMARY KEY, a INT, b INT)
< ok 0 0
> exec
| INSERT INTO dependent_updates VALUES (1, 1, 0)
< ok 1 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| dependent_updates
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT `id`, `a`, `b` FROM `test`.`dependent_updates` ORDER BY 1, 2, 3
< result
< column "id" "INT"
< column "a" "INT"
< column "b" "INT"
< row
: int64 1
: int64 1
: int64 0
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| bak_0_swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_child
: bytes
| parent_id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_parent
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flagged
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| enabled
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_keep
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| user_variable
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
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
| bak_0_generated_default
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
| bak_0_nullable_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| bak_0_outer_gone
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
| bak_0_single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_0_stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| dependent_updates
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
| fk_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| 1
: bytes
| parent_id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| generated_default
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
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
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
| outer_gone
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| outer_keep
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
| separate_overlap
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| user_variable
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| parent_id
: bytes
| test
: bytes
| fk_parent
: bytes
| id
: bytes
| NO ACTION
: bytes
| CASCADE
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE covered_updates (id INT PRIMARY KEY, a INT, b INT)
< ok 0 0
> exec
| INSERT INTO covered_updates VALUES (1, 1, 0), (2, 2, 0)
< ok 2 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| covered_updates
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT `id`, `a`, `b` FROM `test`.`covered_updates` ORDER BY 1, 2, 3
< result
< column "id" "INT"
< column "a" "INT"
< column "b" "INT"
< row
: int64 1
: int64 1
: int64 0
< row
: int64 2
: int64 2
: int64 0
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| bak_0_swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_child
: bytes
| parent_id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_parent
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flagged
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| enabled
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_keep
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| user_variable
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
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
| bak_0_generated_default
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
| bak_0_nullable_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| bak_0_outer_gone
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
| bak_0_single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_0_stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| covered_updates
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
| dependent_updates
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
| fk_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| 1
: bytes
| parent_id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| generated_default
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
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
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
| outer_gone
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| outer_keep
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
| separate_overlap
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| user_variable
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| parent_id
: bytes
| test
: bytes
| fk_parent
: bytes
| id
: bytes
| NO ACTION
: bytes
| CASCADE
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE `test`.`bak_0_covered_updates` LIKE `test`.`covered_updates`
< ok 0 0
> exec
|
| INSERT INTO `test`.`bak_0_covered_updates` (`id`,`a`,`b`) SELECT `covered_updates`.`id`,`covered_updates`.`a`,`covered_updates`.`b` FROM covered_updates WHERE id = 1
< ok 1 0
> exec
| CREATE TABLE `test`.`bak_1_covered_updates` LIKE `test`.`covered_updates`
< ok 0 0
> exec
|
| INSERT INTO `test`.`bak_1_covered_updates` (`id`,`a`,`b`) SELECT `covered_updates`.`id`,`covered_updates`.`a`,`covered_updates`.`b` FROM covered_updates WHERE a = 2
< ok 1 0
> exec
| UPDATE covered_updates SET a = 2 WHERE id = 1
< ok 1 0
> exec
|  UPDATE covered_updates SET a = 3 WHERE a = 2
< ok 2 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| covered_updates
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT `id`, `a`, `b` FROM `test`.`covered_updates` ORDER BY 1, 2, 3
< result
< column "id" "INT"
< column "a" "INT"
< column "b" "INT"
< row
: int64 1
: int64 3
: int64 0
< row
: int64 2
: int64 3
: int64 0
> exec
| /*
| Original SQL:
|  UPDATE covered_updates SET a = 3 WHERE a = 2;
| */
| INSERT INTO `test`.`covered_updates` (`id`, `a`, `b`) SELECT `id`, `a`, `b` FROM `test`.`bak_1_covered_updates` ON DUPLICATE KEY UPDATE `a` = VALUES(`a`)
< ok 2 0
> exec
| /*
| Original SQL:
| UPDATE covered_updates SET a = 2 WHERE id = 1;
| */
| INSERT INTO `test`.`covered_updates` (`id`, `a`, `b`) SELECT `id`, `a`, `b` FROM `test`.`bak_0_covered_updates` ON DUPLICATE KEY UPDATE `a` = VALUES(`a`)
< ok 2 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| covered_updates
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT `id`, `a`, `b` FROM `test`.`covered_updates` ORDER BY 1, 2, 3
< result
< column "id" "INT"
< column "a" "INT"
< column "b" "INT"
< row
: int64 1
: int64 1
: int64 0
< row
: int64 2
: int64 2
: int64 0
> exec
| CREATE TABLE random_rows (id INT PRIMARY KEY)
< ok 0 0
> exec
| INSERT INTO random_rows VALUES (1), (2), (3), (4)
< ok 4 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| random_rows
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
> query
| SELECT `id` FROM `test`.`random_rows` ORDER BY 1
< result
< column "id" "INT"
< row
: int64 1
< row
: int64 2
< row
: int64 3
< row
: int64 4
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| bak_0_swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_child
: bytes
| parent_id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_parent
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flagged
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| enabled
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_keep
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| random_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| user_variable
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
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
| bak_0_covered_updates
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
| bak_0_generated_default
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
| bak_0_nullable_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| bak_0_outer_gone
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
| bak_0_single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_0_stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_1_covered_updates
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
| covered_updates
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
| dependent_updates
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
| fk_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| 1
: bytes
| parent_id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| generated_default
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
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
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
| outer_gone
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| outer_keep
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
| random_rows
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_overlap
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| user_variable
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| parent_id
: bytes
| test
: bytes
| fk_parent
: bytes
| id
: bytes
| NO ACTION
: bytes
| CASCADE
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE use_target (id INT PRIMARY KEY)
< ok 0 0
> exec
| INSERT INTO use_target VALUES (1), (2)
< ok 2 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| use_target
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
> query
| SELECT `id` FROM `test`.`use_target` ORDER BY 1
< result
< column "id" "INT"
< row
: int64 1
< row
: int64 2
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| bak_0_swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_child
: bytes
| parent_id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_parent
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flagged
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| enabled
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_keep
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| random_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| use_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| user_variable
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
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
| bak_0_covered_updates
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
| bak_0_generated_default
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
| bak_0_nullable_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| bak_0_outer_gone
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
| bak_0_single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_0_stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_1_covered_updates
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
| covered_updates
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
| dependent_updates
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
| fk_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| 1
: bytes
| parent_id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| generated_default
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
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
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
| outer_gone
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| outer_keep
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
| random_rows
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_overlap
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| use_target
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| user_variable
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| parent_id
: bytes
| test
: bytes
| fk_parent
: bytes
| id
: bytes
| NO ACTION
: bytes
| CASCADE
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE replaced (id INT PRIMARY KEY, v INT)
< ok 0 0
> exec
| INSERT INTO replaced VALUES (1, 1)
< ok 1 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| replaced
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
< row
: bytes
| v
: bytes
|
> query
| SELECT `id`, `v` FROM `test`.`replaced` ORDER BY 1, 2
< result
< column "id" "INT"
< column "v" "INT"
< row
: int64 1
: int64 1
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| bak_0_swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_child
: bytes
| parent_id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_parent
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flagged
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| enabled
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_keep
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| random_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| replaced
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| replaced
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| use_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| user_variable
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
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
| bak_0_covered_updates
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
| bak_0_generated_default
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
| bak_0_nullable_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| bak_0_outer_gone
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
| bak_0_single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_0_stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_1_covered_updates
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
| covered_updates
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
| dependent_updates
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
| fk_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| 1
: bytes
| parent_id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| generated_default
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
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
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
| outer_gone
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| outer_keep
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
| random_rows
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| replaced
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_overlap
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| use_target
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| user_variable
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| parent_id
: bytes
| test
: bytes
| fk_parent
: bytes
| id
: bytes
| NO ACTION
: bytes
| CASCADE
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE self_assignment (id INT PRIMARY KEY)
< ok 0 0
> exec
| INSERT INTO self_assignment VALUES (1)
< ok 1 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| self_assignment
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
> query
| SELECT `id` FROM `test`.`self_assignment` ORDER BY 1
< result
< column "id" "INT"
< row
: int64 1
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| bak_0_swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_child
: bytes
| parent_id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_parent
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flagged
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| enabled
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_keep
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| random_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| replaced
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| replaced
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_assignment
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| use_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| user_variable
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
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
| bak_0_covered_updates
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
| bak_0_generated_default
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
| bak_0_nullable_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| bak_0_outer_gone
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
| bak_0_single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_0_stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_1_covered_updates
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
| covered_updates
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
| dependent_updates
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
| fk_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| 1
: bytes
| parent_id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| generated_default
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
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
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
| outer_gone
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| outer_keep
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
| random_rows
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| replaced
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| self_assignment
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_overlap
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| use_target
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| user_variable
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| parent_id
: bytes
| test
: bytes
| fk_parent
: bytes
| id
: bytes
| NO ACTION
: bytes
| CASCADE
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE chained_key (id INT PRIMARY KEY, a INT, g1 INT AS (a + 1), g2 INT AS (g1 + 1), UNIQUE KEY uk_g2 (g2))
< ok 0 0
> exec
| INSERT INTO chained_key (id, a) VALUES (1, 10), (2, 20)
< ok 2 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| chained_key
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| g1
: bytes
| VIRTUAL GENERATED
< row
: bytes
| g2
: bytes
| VIRTUAL GENERATED
> query
| SELECT `id`, `a`, `g1`, `g2` FROM `test`.`chained_key` ORDER BY 1, 2, 3, 4
< result
< column "id" "INT"
< column "a" "INT"
< column "g1" "INT"
< column "g2" "INT"
< row
: int64 1
: int64 10
: int64 11
: int64 12
< row
: int64 2
: int64 20
: int64 21
: int64 22
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| bak_0_swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| chained_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| chained_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| chained_key
: bytes
| g1
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| chained_key
: bytes
| g2
: bytes
| `g1` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_child
: bytes
| parent_id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_parent
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flagged
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| enabled
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_keep
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| random_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| replaced
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| replaced
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_assignment
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| use_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| user_variable
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
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
| bak_0_covered_updates
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
| bak_0_generated_default
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
| bak_0_nullable_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| bak_0_outer_gone
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
| bak_0_single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_0_stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_1_covered_updates
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
| chained_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| chained_key
: bytes
| uk_g2
: bytes
| 0
: bytes
| g2
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
| covered_updates
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
| dependent_updates
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
| fk_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| 1
: bytes
| parent_id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| generated_default
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
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
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
| outer_gone
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| outer_keep
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
| random_rows
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| replaced
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| self_assignment
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_overlap
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| use_target
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| user_variable
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| parent_id
: bytes
| test
: bytes
| fk_parent
: bytes
| id
: bytes
| NO ACTION
: bytes
| CASCADE
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE `test`.`bak_0_chained_key` LIKE `test`.`chained_key`
< ok 0 0
> exec
|
| INSERT INTO `test`.`bak_0_chained_key` (`id`,`a`) SELECT `chained_key`.`id`,`chained_key`.`a` FROM chained_key WHERE id = 1
< ok 1 0
> exec
| UPDATE chained_key SET a = a + 5 WHERE id = 1
< ok 1 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| chained_key
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| g1
: bytes
| VIRTUAL GENERATED
< row
: bytes
| g2
: bytes
| VIRTUAL GENERATED
> query
| SELECT `id`, `a`, `g1`, `g2` FROM `test`.`chained_key` ORDER BY 1, 2, 3, 4
< result
< column "id" "INT"
< column "a" "INT"
< column "g1" "INT"
< column "g2" "INT"
< row
: int64 1
: int64 15
: int64 16
: int64 17
< row
: int64 2
: int64 20
: int64 21
: int64 22
> exec
| CREATE TABLE executable_comment (id INT PRIMARY KEY)
< ok 0 0
> exec
| INSERT INTO executable_comment VALUES (1), (2)
< ok 2 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| executable_comment
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
> query
| SELECT `id` FROM `test`.`executable_comment` ORDER BY 1
< result
< column "id" "INT"
< row
: int64 1
< row
: int64 2
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_chained_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_chained_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_chained_key
: bytes
| g1
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_chained_key
: bytes
| g2
: bytes
| `g1` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| bak_0_swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| chained_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| chained_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| chained_key
: bytes
| g1
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| chained_key
: bytes
| g2
: bytes
| `g1` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| executable_comment
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_child
: bytes
| parent_id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_parent
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flagged
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| enabled
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_keep
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| random_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| replaced
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| replaced
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_assignment
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| use_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| user_variable
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
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
| bak_0_chained_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_chained_key
: bytes
| uk_g2
: bytes
| 0
: bytes
| g2
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
| bak_0_covered_updates
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
| bak_0_generated_default
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
| bak_0_nullable_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| bak_0_outer_gone
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
| bak_0_single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_0_stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_1_covered_updates
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
| chained_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| chained_key
: bytes
| uk_g2
: bytes
| 0
: bytes
| g2
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
| covered_updates
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
| dependent_updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| executable_comment
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
| fk_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| 1
: bytes
| parent_id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| generated_default
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
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
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
| outer_gone
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| outer_keep
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
| random_rows
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| replaced
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| self_assignment
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_overlap
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| use_target
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| user_variable
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| parent_id
: bytes
| test
: bytes
| fk_parent
: bytes
| id
: bytes
| NO ACTION
: bytes
| CASCADE
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE `test`.`bak_0_executable_comment` LIKE `test`.`executable_comment`
< ok 0 0
> exec
|
| INSERT INTO `test`.`bak_0_executable_comment` (`id`) SELECT `executable_comment`.`id` FROM executable_comment /*!80000 WHERE id = 1
< ok 1 0
> exec
| DELETE FROM executable_comment /*!80000 WHERE id = 1 */
< ok 1 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| executable_comment
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
> query
| SELECT `id` FROM `test`.`executable_comment` ORDER BY 1
< result
< column "id" "INT"
< row
: int64 2
> exec
| /*
| Original SQL:
| DELETE FROM executable_comment /*!80000 WHERE id = 1 * /;
| */
| INSERT INTO `test`.`executable_comment` (`id`) SELECT `id` FROM `test`.`bak_0_executable_comment`
< ok 1 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| executable_comment
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
> query
| SELECT `id` FROM `test`.`executable_comment` ORDER BY 1
< result
< column "id" "INT"
< row
: int64 1
< row
: int64 2
> exec
| CREATE TABLE separate_swap (id INT PRIMARY KEY, u INT NOT NULL, UNIQUE KEY uk_u (u))
< ok 0 0
> exec
| INSERT INTO separate_swap VALUES (1, 10), (2, 20)
< ok 2 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| separate_swap
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
< row
: bytes
| u
: bytes
|
> query
| SELECT `id`, `u` FROM `test`.`separate_swap` ORDER BY 1, 2
< result
< column "id" "INT"
< column "u" "INT"
< row
: int64 1
: int64 10
< row
: int64 2
: int64 20
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_chained_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_chained_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_chained_key
: bytes
| g1
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_chained_key
: bytes
| g2
: bytes
| `g1` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_executable_comment
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| bak_0_swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| chained_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| chained_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| chained_key
: bytes
| g1
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| chained_key
: bytes
| g2
: bytes
| `g1` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| executable_comment
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_child
: bytes
| parent_id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_parent
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flagged
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| enabled
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_keep
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| random_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| replaced
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| replaced
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_assignment
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| separate_swap
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_swap
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| use_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| user_variable
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
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
| bak_0_chained_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_chained_key
: bytes
| uk_g2
: bytes
| 0
: bytes
| g2
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
| bak_0_covered_updates
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
| bak_0_executable_comment
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
| bak_0_generated_default
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
| bak_0_nullable_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| bak_0_outer_gone
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
| bak_0_single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_0_stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_1_covered_updates
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
| chained_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| chained_key
: bytes
| uk_g2
: bytes
| 0
: bytes
| g2
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
| covered_updates
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
| dependent_updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| executable_comment
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
| fk_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| 1
: bytes
| parent_id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| generated_default
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
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
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
| outer_gone
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| outer_keep
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
| random_rows
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| replaced
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| self_assignment
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_overlap
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_swap
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_swap
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| use_target
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| user_variable
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| parent_id
: bytes
| test
: bytes
| fk_parent
: bytes
| id
: bytes
| NO ACTION
: bytes
| CASCADE
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE flag_source (id INT PRIMARY KEY, flag INT)
< ok 0 0
> exec
| CREATE TABLE flag_target (id INT PRIMARY KEY, v INT)
< ok 0 0
> exec
| INSERT INTO flag_source VALUES (1, 0)
< ok 1 0
> exec
| INSERT INTO flag_target VALUES (1, 0)
< ok 1 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| flag_target
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
< row
: bytes
| v
: bytes
|
> query
| SELECT `id`, `v` FROM `test`.`flag_target` ORDER BY 1, 2
< result
< column "id" "INT"
< column "v" "INT"
< row
: int64 1
: int64 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| flag_source
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
< row
: bytes
| flag
: bytes
|
> query
| SELECT `id`, `flag` FROM `test`.`flag_source` ORDER BY 1, 2
< result
< column "id" "INT"
< column "flag" "INT"
< row
: int64 1
: int64 0
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_chained_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_chained_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_chained_key
: bytes
| g1
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_chained_key
: bytes
| g2
: bytes
| `g1` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_executable_comment
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| bak_0_swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| chained_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| chained_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| chained_key
: bytes
| g1
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| chained_key
: bytes
| g2
: bytes
| `g1` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| executable_comment
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_child
: bytes
| parent_id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_parent
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flag_source
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flag_source
: bytes
| flag
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| flag_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flag_target
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| flagged
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| enabled
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_keep
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| random_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| replaced
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| replaced
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_assignment
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| separate_swap
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_swap
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| use_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| user_variable
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
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
| bak_0_chained_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_chained_key
: bytes
| uk_g2
: bytes
| 0
: bytes
| g2
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
| bak_0_covered_updates
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
| bak_0_executable_comment
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
| bak_0_generated_default
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
| bak_0_nullable_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| bak_0_outer_gone
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
| bak_0_single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_0_stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_1_covered_updates
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
| chained_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| chained_key
: bytes
| uk_g2
: bytes
| 0
: bytes
| g2
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
| covered_updates
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
| dependent_updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| executable_comment
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
| fk_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| 1
: bytes
| parent_id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flag_source
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flag_target
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| generated_default
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
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
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
| outer_gone
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| outer_keep
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
| random_rows
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| replaced
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| self_assignment
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_overlap
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_swap
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_swap
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| use_target
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| user_variable
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| parent_id
: bytes
| test
: bytes
| fk_parent
: bytes
| id
: bytes
| NO ACTION
: bytes
| CASCADE
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE summed_key (id INT PRIMARY KEY, a INT, b INT, g INT AS (a + b) STORED, UNIQUE KEY uk_g (g))
< ok 0 0
> exec
| INSERT INTO summed_key (id, a, b) VALUES (1, 10, 1), (2, 20, 11)
< ok 2 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| summed_key
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| STORED GENERATED
> query
| SELECT `id`, `a`, `b`, `g` FROM `test`.`summed_key` ORDER BY 1, 2, 3, 4
< result
< column "id" "INT"
< column "a" "INT"
< column "b" "INT"
< column "g" "INT"
< row
: int64 1
: int64 10
: int64 1
: int64 11
< row
: int64 2
: int64 20
: int64 11
: int64 31
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_chained_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_chained_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_chained_key
: bytes
| g1
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_chained_key
: bytes
| g2
: bytes
| `g1` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_executable_comment
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| bak_0_swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| chained_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| chained_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| chained_key
: bytes
| g1
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| chained_key
: bytes
| g2
: bytes
| `g1` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| executable_comment
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_child
: bytes
| parent_id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_parent
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flag_source
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flag_source
: bytes
| flag
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| flag_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flag_target
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| flagged
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| enabled
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_keep
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| random_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| replaced
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| replaced
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_assignment
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| separate_swap
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_swap
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| summed_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| summed_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| summed_key
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| summed_key
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| use_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| user_variable
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
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
| bak_0_chained_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_chained_key
: bytes
| uk_g2
: bytes
| 0
: bytes
| g2
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
| bak_0_covered_updates
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
| bak_0_executable_comment
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
| bak_0_generated_default
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
| bak_0_nullable_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| bak_0_outer_gone
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
| bak_0_single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_0_stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_1_covered_updates
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
| chained_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| chained_key
: bytes
| uk_g2
: bytes
| 0
: bytes
| g2
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
| covered_updates
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
| dependent_updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| executable_comment
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
| fk_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| 1
: bytes
| parent_id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flag_source
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flag_target
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| generated_default
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
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
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
| outer_gone
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| outer_keep
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
| random_rows
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| replaced
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| self_assignment
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_overlap
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_swap
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_swap
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| summed_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| summed_key
: bytes
| uk_g
: bytes
| 0
: bytes
| g
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| use_target
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| user_variable
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| parent_id
: bytes
| test
: bytes
| fk_parent
: bytes
| id
: bytes
| NO ACTION
: bytes
| CASCADE
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE `test`.`bak_0_summed_key` LIKE `test`.`summed_key`
< ok 0 0
> exec
|
| INSERT INTO `test`.`bak_0_summed_key` (`id`,`a`,`b`) SELECT `summed_key`.`id`,`summed_key`.`a`,`summed_key`.`b` FROM summed_key ORDER BY id
< ok 2 0
> exec
| UPDATE summed_key SET a = 0 ORDER BY id
< ok 2 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| summed_key
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| STORED GENERATED
> query
| SELECT `id`, `a`, `b`, `g` FROM `test`.`summed_key` ORDER BY 1, 2, 3, 4
< result
< column "id" "INT"
< column "a" "INT"
< column "b" "INT"
< column "g" "INT"
< row
: int64 1
: int64 0
: int64 1
: int64 1
< row
: int64 2
: int64 0
: int64 11
: int64 11
> exec
| CREATE TABLE sysvar_rows (id INT PRIMARY KEY)
< ok 0 0
> exec
| INSERT INTO sysvar_rows VALUES (1), (2)
< ok 2 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| sysvar_rows
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
> query
| SELECT `id` FROM `test`.`sysvar_rows` ORDER BY 1
< result
< column "id" "INT"
< row
: int64 1
< row
: int64 2
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_chained_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_chained_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_chained_key
: bytes
| g1
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_chained_key
: bytes
| g2
: bytes
| `g1` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_executable_comment
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| bak_0_summed_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_summed_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_summed_key
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_summed_key
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| bak_0_swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| chained_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| chained_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| chained_key
: bytes
| g1
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| chained_key
: bytes
| g2
: bytes
| `g1` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| executable_comment
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_child
: bytes
| parent_id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_parent
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flag_source
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flag_source
: bytes
| flag
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| flag_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flag_target
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| flagged
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| enabled
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_keep
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| random_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| replaced
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| replaced
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_assignment
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| separate_swap
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_swap
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| summed_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| summed_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| summed_key
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| summed_key
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| sysvar_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| use_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| user_variable
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
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
| bak_0_chained_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_chained_key
: bytes
| uk_g2
: bytes
| 0
: bytes
| g2
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
| bak_0_covered_updates
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
| bak_0_executable_comment
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
| bak_0_generated_default
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
| bak_0_nullable_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| bak_0_outer_gone
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
| bak_0_single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_0_stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_summed_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_summed_key
: bytes
| uk_g
: bytes
| 0
: bytes
| g
< row
: bytes
| bak_0_swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_1_covered_updates
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
| chained_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| chained_key
: bytes
| uk_g2
: bytes
| 0
: bytes
| g2
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
| covered_updates
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
| dependent_updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| executable_comment
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
| fk_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| 1
: bytes
| parent_id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flag_source
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flag_target
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| generated_default
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
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
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
| outer_gone
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| outer_keep
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
| random_rows
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| replaced
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| self_assignment
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_overlap
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_swap
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_swap
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| summed_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| summed_key
: bytes
| uk_g
: bytes
| 0
: bytes
| g
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| sysvar_rows
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
< row
: bytes
| use_target
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| user_variable
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| parent_id
: bytes
| test
: bytes
| fk_parent
: bytes
| id
: bytes
| NO ACTION
: bytes
| CASCADE
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE fk_stamped_parent (id INT PRIMARY KEY, v INT, ts TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6), UNIQUE KEY uk_ts (ts))
< ok 0 0
> exec
| CREATE TABLE fk_stamped_child (id INT PRIMARY KEY, parent_ts TIMESTAMP(6) NULL, CONSTRAINT fk_ts FOREIGN KEY (parent_ts) REFERENCES fk_stamped_parent (ts) ON UPDATE CASCADE)
< ok 0 0
> exec
| INSERT INTO fk_stamped_parent (id, v, ts) VALUES (1, 1, '2000-01-01 00:00:00')
< ok 1 0
> exec
| INSERT INTO fk_stamped_child VALUES (1, '2000-01-01 00:00:00')
< ok 1 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| fk_stamped_parent
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
< row
: bytes
| v
: bytes
|
< row
: bytes
| ts
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP(6)
> query
| SELECT `id`, `v`, `ts` FROM `test`.`fk_stamped_parent` ORDER BY 1, 2, 3
< result
< column "id" "INT"
< column "v" "INT"
< column "ts" "TIMESTAMP"
< row
: int64 1
: int64 1
: bytes
| 2000-01-01 00:00:00.000000
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| fk_stamped_child
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
< row
: bytes
| parent_ts
: bytes
|
> query
| SELECT `id`, `parent_ts` FROM `test`.`fk_stamped_child` ORDER BY 1, 2
< result
< column "id" "INT"
< column "parent_ts" "TIMESTAMP"
< row
: int64 1
: bytes
| 2000-01-01 00:00:00.000000
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_chained_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_chained_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_chained_key
: bytes
| g1
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_chained_key
: bytes
| g2
: bytes
| `g1` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_executable_comment
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| bak_0_summed_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_summed_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_summed_key
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_summed_key
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| bak_0_swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| chained_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| chained_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| chained_key
: bytes
| g1
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| chained_key
: bytes
| g2
: bytes
| `g1` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| executable_comment
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_child
: bytes
| parent_id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_parent
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_stamped_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_stamped_child
: bytes
| parent_ts
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_stamped_parent
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_stamped_parent
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_stamped_parent
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP(6)
< row
: bytes
| flag_source
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flag_source
: bytes
| flag
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| flag_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flag_target
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| flagged
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| enabled
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_keep
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| random_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| replaced
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| replaced
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_assignment
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| separate_swap
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_swap
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| summed_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| summed_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| summed_key
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| summed_key
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| sysvar_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| use_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| user_variable
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
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
| bak_0_chained_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_chained_key
: bytes
| uk_g2
: bytes
| 0
: bytes
| g2
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
| bak_0_covered_updates
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
| bak_0_executable_comment
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
| bak_0_generated_default
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
| bak_0_nullable_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| bak_0_outer_gone
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
| bak_0_single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_0_stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_summed_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_summed_key
: bytes
| uk_g
: bytes
| 0
: bytes
| g
< row
: bytes
| bak_0_swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_1_covered_updates
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
| chained_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| chained_key
: bytes
| uk_g2
: bytes
| 0
: bytes
| g2
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
| covered_updates
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
| dependent_updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| executable_comment
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
| fk_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| 1
: bytes
| parent_id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_stamped_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_stamped_child
: bytes
| fk_ts
: bytes
| 1
: bytes
| parent_ts
< row
: bytes
| fk_stamped_parent
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_stamped_parent
: bytes
| uk_ts
: bytes
| 0
: bytes
| ts
< row
: bytes
| flag_source
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flag_target
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| generated_default
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
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
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
| outer_gone
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| outer_keep
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
| random_rows
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| replaced
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| self_assignment
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_overlap
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_swap
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_swap
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| summed_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| summed_key
: bytes
| uk_g
: bytes
| 0
: bytes
| g
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| sysvar_rows
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
< row
: bytes
| use_target
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| user_variable
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| parent_id
: bytes
| test
: bytes
| fk_parent
: bytes
| id
: bytes
| NO ACTION
: bytes
| CASCADE
< row
: bytes
| fk_stamped_child
: bytes
| fk_ts
: bytes
| parent_ts
: bytes
| test
: bytes
| fk_stamped_parent
: bytes
| ts
: bytes
| CASCADE
: bytes
| NO ACTION
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE outer_update_keep (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE outer_update_target (id INT PRIMARY KEY, a INT)
< ok 0 0
> exec
| INSERT INTO outer_update_keep VALUES (1), (2)
< ok 2 0
> exec
| INSERT INTO outer_update_target VALUES (1, 1)
< ok 1 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| outer_update_target
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT `id`, `a` FROM `test`.`outer_update_target` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< row
: int64 1
: int64 1
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_chained_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_chained_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_chained_key
: bytes
| g1
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_chained_key
: bytes
| g2
: bytes
| `g1` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_executable_comment
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| bak_0_summed_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_summed_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_summed_key
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_summed_key
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| bak_0_swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| chained_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| chained_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| chained_key
: bytes
| g1
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| chained_key
: bytes
| g2
: bytes
| `g1` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| executable_comment
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_child
: bytes
| parent_id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_parent
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_stamped_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_stamped_child
: bytes
| parent_ts
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_stamped_parent
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_stamped_parent
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_stamped_parent
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP(6)
< row
: bytes
| flag_source
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flag_source
: bytes
| flag
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| flag_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flag_target
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| flagged
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| enabled
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_keep
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_update_keep
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_update_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_update_target
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| random_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| replaced
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| replaced
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_assignment
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| separate_swap
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_swap
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| summed_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| summed_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| summed_key
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| summed_key
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| sysvar_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| use_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| user_variable
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
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
| bak_0_chained_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_chained_key
: bytes
| uk_g2
: bytes
| 0
: bytes
| g2
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
| bak_0_covered_updates
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
| bak_0_executable_comment
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
| bak_0_generated_default
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
| bak_0_nullable_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| bak_0_outer_gone
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
| bak_0_single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_0_stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_summed_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_summed_key
: bytes
| uk_g
: bytes
| 0
: bytes
| g
< row
: bytes
| bak_0_swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_1_covered_updates
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
| chained_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| chained_key
: bytes
| uk_g2
: bytes
| 0
: bytes
| g2
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
| covered_updates
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
| dependent_updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| executable_comment
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
| fk_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| 1
: bytes
| parent_id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_stamped_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_stamped_child
: bytes
| fk_ts
: bytes
| 1
: bytes
| parent_ts
< row
: bytes
| fk_stamped_parent
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_stamped_parent
: bytes
| uk_ts
: bytes
| 0
: bytes
| ts
< row
: bytes
| flag_source
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flag_target
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| generated_default
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
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
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
| outer_gone
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| outer_keep
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| outer_update_keep
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| outer_update_target
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
| random_rows
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| replaced
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| self_assignment
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_overlap
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_swap
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_swap
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| summed_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| summed_key
: bytes
| uk_g
: bytes
| 0
: bytes
| g
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| sysvar_rows
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
< row
: bytes
| use_target
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| user_variable
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| parent_id
: bytes
| test
: bytes
| fk_parent
: bytes
| id
: bytes
| NO ACTION
: bytes
| CASCADE
< row
: bytes
| fk_stamped_child
: bytes
| fk_ts
: bytes
| parent_ts
: bytes
| test
: bytes
| fk_stamped_parent
: bytes
| ts
: bytes
| CASCADE
: bytes
| NO ACTION
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE `test`.`bak_0_outer_update_target` LIKE `test`.`outer_update_target`
< ok 0 0
> exec
|
| INSERT INTO `test`.`bak_0_outer_update_target` (`id`,`a`) SELECT DISTINCT `t`.`id`,`t`.`a` FROM outer_update_keep u LEFT JOIN outer_update_target t ON t.id = u.id WHERE `t`.`id` IS NOT NULL
< ok 1 0
> exec
| UPDATE outer_update_keep u LEFT JOIN outer_update_target t ON t.id = u.id SET t.a = 5
< ok 1 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| outer_update_target
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT `id`, `a` FROM `test`.`outer_update_target` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< row
: int64 1
: int64 5
> exec
| /*
| Original SQL:
| UPDATE outer_update_keep u LEFT JOIN outer_update_target t ON t.id = u.id SET t.a = 5;
| */
| INSERT INTO `test`.`outer_update_target` (`id`, `a`) SELECT `id`, `a` FROM `test`.`bak_0_outer_update_target` ON DUPLICATE KEY UPDATE `a` = VALUES(`a`)
< ok 2 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| outer_update_target
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
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
| SELECT `id`, `a` FROM `test`.`outer_update_target` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< row
: int64 1
: int64 1
> exec
| CREATE TABLE ai_source (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE ai_target (id INT PRIMARY KEY)
< ok 0 0
> exec
| INSERT INTO ai_target VALUES (1)
< ok 1 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| ai_target
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
> query
| SELECT `id` FROM `test`.`ai_target` ORDER BY 1
< result
< column "id" "INT"
< row
: int64 1
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| ai_source
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| ai_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_chained_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_chained_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_chained_key
: bytes
| g1
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_chained_key
: bytes
| g2
: bytes
| `g1` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_executable_comment
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_outer_update_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_outer_update_target
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| bak_0_summed_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_summed_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_summed_key
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_summed_key
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| bak_0_swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| chained_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| chained_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| chained_key
: bytes
| g1
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| chained_key
: bytes
| g2
: bytes
| `g1` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| executable_comment
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_child
: bytes
| parent_id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_parent
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_stamped_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_stamped_child
: bytes
| parent_ts
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_stamped_parent
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_stamped_parent
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_stamped_parent
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP(6)
< row
: bytes
| flag_source
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flag_source
: bytes
| flag
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| flag_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flag_target
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| flagged
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| enabled
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_keep
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_update_keep
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_update_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_update_target
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| random_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| replaced
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| replaced
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_assignment
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| separate_swap
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_swap
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| summed_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| summed_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| summed_key
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| summed_key
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| sysvar_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| use_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| user_variable
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
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
| ai_source
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| ai_target
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
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
| bak_0_chained_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_chained_key
: bytes
| uk_g2
: bytes
| 0
: bytes
| g2
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
| bak_0_covered_updates
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
| bak_0_executable_comment
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
| bak_0_generated_default
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
| bak_0_nullable_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| bak_0_outer_gone
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_outer_update_target
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
| bak_0_single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_0_stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_summed_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_summed_key
: bytes
| uk_g
: bytes
| 0
: bytes
| g
< row
: bytes
| bak_0_swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_1_covered_updates
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
| chained_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| chained_key
: bytes
| uk_g2
: bytes
| 0
: bytes
| g2
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
| covered_updates
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
| dependent_updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| executable_comment
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
| fk_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| 1
: bytes
| parent_id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_stamped_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_stamped_child
: bytes
| fk_ts
: bytes
| 1
: bytes
| parent_ts
< row
: bytes
| fk_stamped_parent
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_stamped_parent
: bytes
| uk_ts
: bytes
| 0
: bytes
| ts
< row
: bytes
| flag_source
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flag_target
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| generated_default
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
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
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
| outer_gone
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| outer_keep
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| outer_update_keep
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| outer_update_target
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
| random_rows
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| replaced
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| self_assignment
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_overlap
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_swap
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_swap
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| summed_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| summed_key
: bytes
| uk_g
: bytes
| 0
: bytes
| g
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| sysvar_rows
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
< row
: bytes
| use_target
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| user_variable
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| parent_id
: bytes
| test
: bytes
| fk_parent
: bytes
| id
: bytes
| NO ACTION
: bytes
| CASCADE
< row
: bytes
| fk_stamped_child
: bytes
| fk_ts
: bytes
| parent_ts
: bytes
| test
: bytes
| fk_stamped_parent
: bytes
| ts
: bytes
| CASCADE
: bytes
| NO ACTION
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE DATABASE IF NOT EXISTS omni_tidb_analysis_backup_other
< ok 0 0
> exec
| CREATE TABLE xdb_parent (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE omni_tidb_analysis_backup_other.xdb_child (id INT PRIMARY KEY, parent_id INT, CONSTRAINT fk_xdb_child FOREIGN KEY (parent_id) REFERENCES test.xdb_parent (id) ON DELETE CASCADE)
< ok 0 0
> exec
| INSERT INTO xdb_parent VALUES (1), (2)
< ok 2 0
> exec
| INSERT INTO omni_tidb_analysis_backup_other.xdb_child VALUES (1, 1), (2, 2)
< ok 2 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| xdb_parent
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
> query
| SELECT `id` FROM `test`.`xdb_parent` ORDER BY 1
< result
< column "id" "INT"
< row
: int64 1
< row
: int64 2
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
: string
| xdb_child
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
< row
: bytes
| parent_id
: bytes
|
> query
| SELECT `id`, `parent_id` FROM `omni_tidb_analysis_backup_other`.`xdb_child` ORDER BY 1, 2
< result
< column "id" "INT"
< column "parent_id" "INT"
< row
: int64 1
: int64 1
< row
: int64 2
: int64 2
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| ai_source
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| ai_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_chained_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_chained_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_chained_key
: bytes
| g1
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_chained_key
: bytes
| g2
: bytes
| `g1` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_executable_comment
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_outer_update_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_outer_update_target
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| bak_0_summed_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_summed_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_summed_key
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_summed_key
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| bak_0_swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| chained_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| chained_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| chained_key
: bytes
| g1
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| chained_key
: bytes
| g2
: bytes
| `g1` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| executable_comment
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_child
: bytes
| parent_id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_parent
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_stamped_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_stamped_child
: bytes
| parent_ts
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_stamped_parent
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_stamped_parent
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_stamped_parent
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP(6)
< row
: bytes
| flag_source
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flag_source
: bytes
| flag
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| flag_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flag_target
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| flagged
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| enabled
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_keep
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_update_keep
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_update_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_update_target
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| random_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| replaced
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| replaced
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_assignment
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| separate_swap
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_swap
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| summed_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| summed_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| summed_key
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| summed_key
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| sysvar_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| use_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| user_variable
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| xdb_parent
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
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
| ai_source
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| ai_target
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
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
| bak_0_chained_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_chained_key
: bytes
| uk_g2
: bytes
| 0
: bytes
| g2
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
| bak_0_covered_updates
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
| bak_0_executable_comment
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
| bak_0_generated_default
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
| bak_0_nullable_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| bak_0_outer_gone
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_outer_update_target
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
| bak_0_single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_0_stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_summed_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_summed_key
: bytes
| uk_g
: bytes
| 0
: bytes
| g
< row
: bytes
| bak_0_swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_1_covered_updates
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
| chained_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| chained_key
: bytes
| uk_g2
: bytes
| 0
: bytes
| g2
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
| covered_updates
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
| dependent_updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| executable_comment
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
| fk_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| 1
: bytes
| parent_id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_stamped_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_stamped_child
: bytes
| fk_ts
: bytes
| 1
: bytes
| parent_ts
< row
: bytes
| fk_stamped_parent
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_stamped_parent
: bytes
| uk_ts
: bytes
| 0
: bytes
| ts
< row
: bytes
| flag_source
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flag_target
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| generated_default
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
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
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
| outer_gone
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| outer_keep
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| outer_update_keep
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| outer_update_target
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
| random_rows
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| replaced
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| self_assignment
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_overlap
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_swap
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_swap
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| summed_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| summed_key
: bytes
| uk_g
: bytes
| 0
: bytes
| g
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| sysvar_rows
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
< row
: bytes
| use_target
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| user_variable
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| xdb_parent
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| parent_id
: bytes
| test
: bytes
| fk_parent
: bytes
| id
: bytes
| NO ACTION
: bytes
| CASCADE
< row
: bytes
| fk_stamped_child
: bytes
| fk_ts
: bytes
| parent_ts
: bytes
| test
: bytes
| fk_stamped_parent
: bytes
| ts
: bytes
| CASCADE
: bytes
| NO ACTION
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| xdb_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| xdb_child
: bytes
| parent_id
: bytes
|
: bytes
| YES
: bytes
|
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
< row
: bytes
| xdb_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| xdb_child
: bytes
| fk_xdb_child
: bytes
| 1
: bytes
| parent_id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< row
: bytes
| xdb_child
: bytes
| fk_xdb_child
: bytes
| parent_id
: bytes
| test
: bytes
| xdb_parent
: bytes
| id
: bytes
| NO ACTION
: bytes
| CASCADE
> exec
| CREATE DATABASE IF NOT EXISTS omni_tidb_analysis_backup_other
< ok 0 0
> exec
| CREATE TABLE xdb_twin (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE omni_tidb_analysis_backup_other.xdb_twin (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE omni_tidb_analysis_backup_other.xdb_twin_child (id INT PRIMARY KEY, twin_id INT, CONSTRAINT fk_xdb_twin_child FOREIGN KEY (twin_id) REFERENCES omni_tidb_analysis_backup_other.xdb_twin (id) ON DELETE CASCADE)
< ok 0 0
> exec
| INSERT INTO xdb_twin VALUES (1), (2)
< ok 2 0
> exec
| INSERT INTO omni_tidb_analysis_backup_other.xdb_twin VALUES (1), (2)
< ok 2 0
> exec
| INSERT INTO omni_tidb_analysis_backup_other.xdb_twin_child VALUES (1, 1), (2, 2)
< ok 2 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| xdb_twin
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
> query
| SELECT `id` FROM `test`.`xdb_twin` ORDER BY 1
< result
< column "id" "INT"
< row
: int64 1
< row
: int64 2
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
: string
| xdb_twin
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
> query
| SELECT `id` FROM `omni_tidb_analysis_backup_other`.`xdb_twin` ORDER BY 1
< result
< column "id" "INT"
< row
: int64 1
< row
: int64 2
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
: string
| xdb_twin_child
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
< row
: bytes
| twin_id
: bytes
|
> query
| SELECT `id`, `twin_id` FROM `omni_tidb_analysis_backup_other`.`xdb_twin_child` ORDER BY 1, 2
< result
< column "id" "INT"
< column "twin_id" "INT"
< row
: int64 1
: int64 1
< row
: int64 2
: int64 2
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| ai_source
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| ai_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_back`tick
: bytes
| c`ol
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_chained_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_chained_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_chained_key
: bytes
| g1
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_chained_key
: bytes
| g2
: bytes
| `g1` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_executable_comment
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_0_keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_outer_update_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_outer_update_target
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| bak_0_summed_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_summed_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_summed_key
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_0_summed_key
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| bak_0_swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_0_swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_1_covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_test
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| bak_overlapping_rows_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_overlapping_rows_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_test
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| chained_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| chained_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| chained_key
: bytes
| g1
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| chained_key
: bytes
| g2
: bytes
| `g1` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| comments
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| comments
: bytes
| note
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| covered_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| covered_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| cte_alias
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| cte_alias
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| dependent_updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| executable_comment
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| expression_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| expression_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_child
: bytes
| parent_id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_parent
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_stamped_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_stamped_child
: bytes
| parent_ts
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_stamped_parent
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| fk_stamped_parent
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| fk_stamped_parent
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP(6)
< row
: bytes
| flag_source
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flag_source
: bytes
| flag
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| flag_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flag_target
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| flagged
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| flags
: bytes
| enabled
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| generated_default
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_default
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| generated_key
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| generated_key
: bytes
| g
: bytes
| `a` + 1
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| keyless_children
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| keyless_children
: bytes
| parent
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| keyless_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| limits
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| limits
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| VIRTUAL GENERATED
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| nullable_key
: bytes
| u
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| nullable_key
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| one_to_many
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| one_to_many
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_gone
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_gone
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| outer_keep
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_update_keep
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_update_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| outer_update_target
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| overlapping_limits
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| overlapping_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| paren_peers
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| parens
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| random_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| replaced
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| replaced
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_assignment
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_overlap
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| separate_swap
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| separate_swap
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| single_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| stamped
: bytes
| v
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| stamped
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| DEFAULT_GENERATED on update CURRENT_TIMESTAMP
< row
: bytes
| summed_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| summed_key
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| summed_key
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| summed_key
: bytes
| g
: bytes
| `a` + `b`
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| swapped_key
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| swapped_key
: bytes
| u
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| sysvar_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| updates
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| b
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| updates
: bytes
| g
: bytes
| `a` * 2
: bytes
| YES
: bytes
| STORED GENERATED
< row
: bytes
| use_target
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| user_variable
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| xdb_parent
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| xdb_twin
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
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
| ai_source
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| ai_target
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
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
| bak_0_chained_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_chained_key
: bytes
| uk_g2
: bytes
| 0
: bytes
| g2
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
| bak_0_covered_updates
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
| bak_0_executable_comment
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
| bak_0_generated_default
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
| bak_0_nullable_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
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
| bak_0_outer_gone
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_outer_update_target
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
| bak_0_single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_0_stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_summed_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_summed_key
: bytes
| uk_g
: bytes
| 0
: bytes
| g
< row
: bytes
| bak_0_swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| bak_0_swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| bak_1_covered_updates
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
| chained_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| chained_key
: bytes
| uk_g2
: bytes
| 0
: bytes
| g2
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
| covered_updates
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
| dependent_updates
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| executable_comment
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
| fk_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| 1
: bytes
| parent_id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_stamped_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_stamped_child
: bytes
| fk_ts
: bytes
| 1
: bytes
| parent_ts
< row
: bytes
| fk_stamped_parent
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| fk_stamped_parent
: bytes
| uk_ts
: bytes
| 0
: bytes
| ts
< row
: bytes
| flag_source
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flag_target
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| generated_default
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
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
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
| outer_gone
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| outer_keep
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| outer_update_keep
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| outer_update_target
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
| random_rows
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| replaced
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| self_assignment
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_overlap
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_swap
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| separate_swap
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| single_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| stamped
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| summed_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| summed_key
: bytes
| uk_g
: bytes
| 0
: bytes
| g
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: bytes
| 0
: bytes
| u
< row
: bytes
| sysvar_rows
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
< row
: bytes
| use_target
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| user_variable
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| xdb_parent
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| xdb_twin
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| test
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: bytes
| parent_id
: bytes
| test
: bytes
| fk_parent
: bytes
| id
: bytes
| NO ACTION
: bytes
| CASCADE
< row
: bytes
| fk_stamped_child
: bytes
| fk_ts
: bytes
| parent_ts
: bytes
| test
: bytes
| fk_stamped_parent
: bytes
| ts
: bytes
| CASCADE
: bytes
| NO ACTION
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "TEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| xdb_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| xdb_child
: bytes
| parent_id
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| xdb_twin
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| xdb_twin_child
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| xdb_twin_child
: bytes
| twin_id
: bytes
|
: bytes
| YES
: bytes
|
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COALESCE(COLUMN_NAME, EXPRESSION) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "VARCHAR"
< column "COALESCE(COLUMN_NAME, EXPRESSION)" "VARCHAR"
< row
: bytes
| xdb_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| xdb_child
: bytes
| fk_xdb_child
: bytes
| 1
: bytes
| parent_id
< row
: bytes
| xdb_twin
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| xdb_twin_child
: bytes
| PRIMARY
: bytes
| 0
: bytes
| id
< row
: bytes
| xdb_twin_child
: bytes
| fk_xdb_twin_child
: bytes
| 1
: bytes
| twin_id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_SCHEMA, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_SCHEMA" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
< row
: bytes
| xdb_child
: bytes
| fk_xdb_child
: bytes
| parent_id
: bytes
| test
: bytes
| xdb_parent
: bytes
| id
: bytes
| NO ACTION
: bytes
| CASCADE
< row
: bytes
| xdb_twin_child
: bytes
| fk_xdb_twin_child
: bytes
| twin_id
: bytes
| omni_tidb_analysis_backup_other
: bytes
| xdb_twin
: bytes
| id
: bytes
| NO ACTION
: bytes
| CASCADE
> exec
| CREATE TABLE `test`.`bak_0_xdb_twin` LIKE `test`.`xdb_twin`
< ok 0 0
> exec
|
| INSERT INTO `test`.`bak_0_xdb_twin` (`id`) SELECT `xdb_twin`.`id` FROM xdb_twin WHERE id = 1
< ok 1 0
> exec
| DELETE FROM xdb_twin WHERE id = 1
< ok 1 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| xdb_twin
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
> query
| SELECT `id` FROM `test`.`xdb_twin` ORDER BY 1
< result
< column "id" "INT"
< row
: int64 2
> exec
| /*
| Original SQL:
| DELETE FROM xdb_twin WHERE id = 1;
| */
| INSERT INTO `test`.`xdb_twin` (`id`) SELECT `id` FROM `test`.`bak_0_xdb_twin`
< ok 1 0
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| test
: string
| xdb_twin
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
> query
| SELECT `id` FROM `test`.`xdb_twin` ORDER BY 1
< result
< column "id" "INT"
< row
: int64 1
< row
: int64 2
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
: string
| xdb_twin
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
> query
| SELECT `id` FROM `omni_tidb_analysis_backup_other`.`xdb_twin` ORDER BY 1
< result
< column "id" "INT"
< row
: int64 1
< row
: int64 2
> query
| SELECT COLUMN_NAME, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_tidb_analysis_backup_other
: string
| xdb_twin_child
< result
< column "COLUMN_NAME" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| id
: bytes
|
< row
: bytes
| twin_id
: bytes
|
> query
| SELECT `id`, `twin_id` FROM `omni_tidb_analysis_backup_other`.`xdb_twin_child` ORDER BY 1, 2
< result
< column "id" "INT"
< column "twin_id" "INT"
< row
: int64 1
: int64 1
< row
: int64 2
: int64 2
