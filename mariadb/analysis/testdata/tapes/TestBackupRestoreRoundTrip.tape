-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestBackupRestoreRoundTrip default#1
> exec
| CREATE TABLE limits (id INT PRIMARY KEY, a INT, b INT, g INT AS (a + b), UNIQUE KEY uk_b (b))
< ok 0 0
> exec
| INSERT INTO limits (id, a, b) VALUES (1, 1, 10), (2, 2, 20), (3, 3, 30), (4, 4, 40)
< ok 4 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
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
| SELECT `id`, `a`, `b`, `g` FROM `omni_testbackuprestoreroundtrip`.`limits` ORDER BY 1, 2, 3, 4
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
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: int64 0
: bytes
| b
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
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
| DELETE FROM limits WHERE id > 1 ORDER BY id DESC LIMIT 1
< ok 1 0
> exec
|  DELETE FROM limits WHERE id = 1
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
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
| SELECT `id`, `a`, `b`, `g` FROM `omni_testbackuprestoreroundtrip`.`limits` ORDER BY 1, 2, 3, 4
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
| DELETE FROM limits WHERE id > 1 ORDER BY id DESC LIMIT 1;
|  DELETE FROM limits WHERE id = 1;
| */
| INSERT INTO `omni_testbackuprestoreroundtrip`.`limits` (`id`, `a`, `b`) SELECT `id`, `a`, `b` FROM `omni_testbackuprestoreroundtrip`.`bak_limits_omni_testbackuprestoreroundtrip`
< ok 2 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
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
| SELECT `id`, `a`, `b`, `g` FROM `omni_testbackuprestoreroundtrip`.`limits` ORDER BY 1, 2, 3, 4
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
| CREATE TABLE overlapping_limits (id INT PRIMARY KEY)
< ok 0 0
> exec
| INSERT INTO overlapping_limits VALUES (1), (2)
< ok 2 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
: string
| overlapping_limits
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
> query
| SELECT `id` FROM `omni_testbackuprestoreroundtrip`.`overlapping_limits` ORDER BY 1
< result
< column "id" "INT"
< row
: int64 1
< row
: int64 2
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_limits_omni_testbackuprestoreroundtrip
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE updates (id INT PRIMARY KEY, a INT, b INT, g INT AS (a * 2) STORED)
< ok 0 0
> exec
| INSERT INTO updates (id, a, b) VALUES (1, 1, 1), (2, 2, 2), (3, 3, 3)
< ok 3 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
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
| SELECT `id`, `a`, `b`, `g` FROM `omni_testbackuprestoreroundtrip`.`updates` ORDER BY 1, 2, 3, 4
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
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_limits_omni_testbackuprestoreroundtrip
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE `omni_testbackuprestoreroundtrip`.`bak_updates_omni_testbackuprestoreroundtrip` LIKE `omni_testbackuprestoreroundtrip`.`updates`
< ok 0 0
> exec
|
| INSERT INTO `omni_testbackuprestoreroundtrip`.`bak_updates_omni_testbackuprestoreroundtrip` (`id`,`a`,`b`)  (SELECT `updates`.`id`,`updates`.`a`,`updates`.`b` FROM updates WHERE id = 1)
|   UNION DISTINCT
|   (SELECT `updates`.`id`,`updates`.`a`,`updates`.`b` FROM updates WHERE id = 2)
< ok 2 0
> exec
| UPDATE updates SET a = 10 WHERE id = 1
< ok 1 0
> exec
|  UPDATE updates SET b = 20 WHERE id = 2
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
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
| SELECT `id`, `a`, `b`, `g` FROM `omni_testbackuprestoreroundtrip`.`updates` ORDER BY 1, 2, 3, 4
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
: int64 3
: int64 3
: int64 6
> exec
| /*
| Original SQL:
| UPDATE updates SET a = 10 WHERE id = 1;
|  UPDATE updates SET b = 20 WHERE id = 2;
| */
| INSERT INTO `omni_testbackuprestoreroundtrip`.`updates` (`id`, `a`, `b`) SELECT `id`, `a`, `b` FROM `omni_testbackuprestoreroundtrip`.`bak_updates_omni_testbackuprestoreroundtrip` ON DUPLICATE KEY UPDATE `a` = VALUES(`a`), `b` = VALUES(`b`)
< ok 4 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
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
| SELECT `id`, `a`, `b`, `g` FROM `omni_testbackuprestoreroundtrip`.`updates` ORDER BY 1, 2, 3, 4
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
> exec
| CREATE TABLE self_join (id INT PRIMARY KEY, a INT, b INT)
< ok 0 0
> exec
| INSERT INTO self_join VALUES (1, 1, 1), (2, 2, 2), (3, 3, 3)
< ok 3 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
: string
| self_join
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
> query
| SELECT `id`, `a`, `b` FROM `omni_testbackuprestoreroundtrip`.`self_join` ORDER BY 1, 2, 3
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
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
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
| self_join
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| self_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| self_join
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE `omni_testbackuprestoreroundtrip`.`bak_self_join_omni_testbackuprestoreroundtrip` LIKE `omni_testbackuprestoreroundtrip`.`self_join`
< ok 0 0
> exec
|
| INSERT INTO `omni_testbackuprestoreroundtrip`.`bak_self_join_omni_testbackuprestoreroundtrip` (`id`,`a`,`b`)  (SELECT DISTINCT `x`.`id`,`x`.`a`,`x`.`b` FROM self_join x JOIN self_join y ON x.id = y.id + 1)
|   UNION DISTINCT
|   (SELECT DISTINCT `y`.`id`,`y`.`a`,`y`.`b` FROM self_join x JOIN self_join y ON x.id = y.id + 1)
< ok 3 0
> exec
| UPDATE self_join x JOIN self_join y ON x.id = y.id + 1 SET x.a = 10, y.b = 20
< ok 4 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
: string
| self_join
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
> query
| SELECT `id`, `a`, `b` FROM `omni_testbackuprestoreroundtrip`.`self_join` ORDER BY 1, 2, 3
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
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
: string
| self_join
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
> query
| SELECT `id`, `a`, `b` FROM `omni_testbackuprestoreroundtrip`.`self_join` ORDER BY 1, 2, 3
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
| omni_testbackuprestoreroundtrip
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
| SELECT `id`, `a` FROM `omni_testbackuprestoreroundtrip`.`one_to_many` ORDER BY 1, 2
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
| omni_testbackuprestoreroundtrip
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
| SELECT `id`, `parent` FROM `omni_testbackuprestoreroundtrip`.`children` ORDER BY 1, 2
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
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
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
| self_join
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| self_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| self_join
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE `omni_testbackuprestoreroundtrip`.`bak_one_to_many_omni_testbackuprestoreroundtrip` LIKE `omni_testbackuprestoreroundtrip`.`one_to_many`
< ok 0 0
> exec
|
| INSERT INTO `omni_testbackuprestoreroundtrip`.`bak_one_to_many_omni_testbackuprestoreroundtrip` (`id`,`a`)  SELECT DISTINCT `p`.`id`,`p`.`a` FROM one_to_many p JOIN children c ON c.parent = p.id WHERE p.id = 1
< ok 1 0
> exec
| DELETE p FROM one_to_many p JOIN children c ON c.parent = p.id WHERE p.id = 1
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
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
| SELECT `id`, `a` FROM `omni_testbackuprestoreroundtrip`.`one_to_many` ORDER BY 1, 2
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
| INSERT INTO `omni_testbackuprestoreroundtrip`.`one_to_many` (`id`, `a`) SELECT `id`, `a` FROM `omni_testbackuprestoreroundtrip`.`bak_one_to_many_omni_testbackuprestoreroundtrip`
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
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
| SELECT `id`, `a` FROM `omni_testbackuprestoreroundtrip`.`one_to_many` ORDER BY 1, 2
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
| omni_testbackuprestoreroundtrip
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
| SELECT `id`, `parent` FROM `omni_testbackuprestoreroundtrip`.`children` ORDER BY 1, 2
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
| omni_testbackuprestoreroundtrip
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
| SELECT `id`, `a` FROM `omni_testbackuprestoreroundtrip`.`keyless_join` ORDER BY 1, 2
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
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
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
| self_join
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| self_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| keyless_children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| self_join
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
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
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
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
| SELECT `id`, `a` FROM `omni_testbackuprestoreroundtrip`.`keyless_rows` ORDER BY 1, 2
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
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
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
| self_join
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| self_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| keyless_children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| self_join
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE `omni_testbackuprestoreroundtrip`.`bak_keyless_rows_omni_testbackuprestoreroundtrip` LIKE `omni_testbackuprestoreroundtrip`.`keyless_rows`
< ok 0 0
> exec
|
| INSERT INTO `omni_testbackuprestoreroundtrip`.`bak_keyless_rows_omni_testbackuprestoreroundtrip` (`id`,`a`)  SELECT `keyless_rows`.`id`,`keyless_rows`.`a` FROM keyless_rows WHERE id = 1
< ok 2 0
> exec
| DELETE FROM keyless_rows WHERE id = 1
< ok 2 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
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
| SELECT `id`, `a` FROM `omni_testbackuprestoreroundtrip`.`keyless_rows` ORDER BY 1, 2
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
| INSERT INTO `omni_testbackuprestoreroundtrip`.`keyless_rows` (`id`, `a`) SELECT `id`, `a` FROM `omni_testbackuprestoreroundtrip`.`bak_keyless_rows_omni_testbackuprestoreroundtrip`
< ok 2 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
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
| SELECT `id`, `a` FROM `omni_testbackuprestoreroundtrip`.`keyless_rows` ORDER BY 1, 2
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
| CREATE TABLE invisible (id INT PRIMARY KEY, secret INT INVISIBLE DEFAULT 0)
< ok 0 0
> exec
| INSERT INTO invisible (id, secret) VALUES (1, 42), (2, 43)
< ok 2 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
: string
| invisible
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| secret
> query
| SELECT `id`, `secret` FROM `omni_testbackuprestoreroundtrip`.`invisible` ORDER BY 1, 2
< result
< column "id" "INT"
< column "secret" "INT"
< row
: int64 1
: int64 42
< row
: int64 2
: int64 43
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| bak_keyless_rows_omni_testbackuprestoreroundtrip
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
| bak_keyless_rows_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
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
| invisible
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| invisible
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
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
| self_join
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| self_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| invisible
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| keyless_children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| self_join
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE `omni_testbackuprestoreroundtrip`.`bak_invisible_omni_testbackuprestoreroundtrip` LIKE `omni_testbackuprestoreroundtrip`.`invisible`
< ok 0 0
> exec
|
| INSERT INTO `omni_testbackuprestoreroundtrip`.`bak_invisible_omni_testbackuprestoreroundtrip` (`id`,`secret`)  SELECT `invisible`.`id`,`invisible`.`secret` FROM invisible WHERE id = 1
< ok 1 0
> exec
| DELETE FROM invisible WHERE id = 1
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
: string
| invisible
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| secret
> query
| SELECT `id`, `secret` FROM `omni_testbackuprestoreroundtrip`.`invisible` ORDER BY 1, 2
< result
< column "id" "INT"
< column "secret" "INT"
< row
: int64 2
: int64 43
> exec
| /*
| Original SQL:
| DELETE FROM invisible WHERE id = 1;
| */
| INSERT INTO `omni_testbackuprestoreroundtrip`.`invisible` (`id`, `secret`) SELECT `id`, `secret` FROM `omni_testbackuprestoreroundtrip`.`bak_invisible_omni_testbackuprestoreroundtrip`
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
: string
| invisible
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| secret
> query
| SELECT `id`, `secret` FROM `omni_testbackuprestoreroundtrip`.`invisible` ORDER BY 1, 2
< result
< column "id" "INT"
< column "secret" "INT"
< row
: int64 1
: int64 42
< row
: int64 2
: int64 43
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
| omni_testbackuprestoreroundtrip
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
| SELECT `id`, `a` FROM `omni_testbackuprestoreroundtrip`.`parens` ORDER BY 1, 2
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
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
< row
: bytes
| bak_keyless_rows_omni_testbackuprestoreroundtrip
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
| bak_keyless_rows_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
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
| invisible
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| invisible
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
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
| self_join
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| self_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| invisible
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| keyless_children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| parens
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| paren_peers
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| self_join
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE `omni_testbackuprestoreroundtrip`.`bak_parens_omni_testbackuprestoreroundtrip` LIKE `omni_testbackuprestoreroundtrip`.`parens`
< ok 0 0
> exec
|
| INSERT INTO `omni_testbackuprestoreroundtrip`.`bak_parens_omni_testbackuprestoreroundtrip` (`id`,`a`)  SELECT DISTINCT `parens`.`id`,`parens`.`a` FROM (parens JOIN paren_peers ON parens.id = paren_peers.id) JOIN paren_peers AS p ON p.id = parens.id WHERE parens.id = 1
< ok 1 0
> exec
| UPDATE (parens JOIN paren_peers ON parens.id = paren_peers.id) JOIN paren_peers AS p ON p.id = parens.id SET parens.a = 9 WHERE parens.id = 1
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
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
| SELECT `id`, `a` FROM `omni_testbackuprestoreroundtrip`.`parens` ORDER BY 1, 2
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
| INSERT INTO `omni_testbackuprestoreroundtrip`.`parens` (`id`, `a`) SELECT `id`, `a` FROM `omni_testbackuprestoreroundtrip`.`bak_parens_omni_testbackuprestoreroundtrip` ON DUPLICATE KEY UPDATE `a` = VALUES(`a`)
< ok 2 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
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
| SELECT `id`, `a` FROM `omni_testbackuprestoreroundtrip`.`parens` ORDER BY 1, 2
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
| omni_testbackuprestoreroundtrip
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
| SELECT `id`, `c``ol` FROM `omni_testbackuprestoreroundtrip`.`back``tick` ORDER BY 1, 2
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
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
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
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
< row
: bytes
| bak_keyless_rows_omni_testbackuprestoreroundtrip
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
| bak_keyless_rows_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
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
| invisible
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| invisible
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
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
| self_join
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| self_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| invisible
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| keyless_children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| parens
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| paren_peers
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| self_join
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE `omni_testbackuprestoreroundtrip`.`bak_back``tick_omni_testbackuprestoreroundtrip` LIKE `omni_testbackuprestoreroundtrip`.`back``tick`
< ok 0 0
> exec
|
| INSERT INTO `omni_testbackuprestoreroundtrip`.`bak_back``tick_omni_testbackuprestoreroundtrip` (`id`,`c``ol`)  SELECT `back``tick`.`id`,`back``tick`.`c``ol` FROM `back``tick` WHERE `c``ol` = 1
< ok 1 0
> exec
| DELETE FROM `back``tick` WHERE `c``ol` = 1
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
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
| SELECT `id`, `c``ol` FROM `omni_testbackuprestoreroundtrip`.`back``tick` ORDER BY 1, 2
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
| INSERT INTO `omni_testbackuprestoreroundtrip`.`back``tick` (`id`, `c``ol`) SELECT `id`, `c``ol` FROM `omni_testbackuprestoreroundtrip`.`bak_back``tick_omni_testbackuprestoreroundtrip`
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
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
| SELECT `id`, `c``ol` FROM `omni_testbackuprestoreroundtrip`.`back``tick` ORDER BY 1, 2
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
| CREATE TABLE t (id INT PRIMARY KEY, a INT)
< ok 0 0
> exec
| CREATE TABLE u (id INT PRIMARY KEY, b INT)
< ok 0 0
> exec
| DROP DATABASE IF EXISTS omni_mariadb_analysis_backup_other
< ok 0 0
> exec
| CREATE DATABASE omni_mariadb_analysis_backup_other
< ok 1 0
> exec
| CREATE TABLE omni_mariadb_analysis_backup_other.t (id INT PRIMARY KEY, b INT)
< ok 0 0
> exec
| CREATE TABLE omni_mariadb_analysis_backup_other.u (id INT PRIMARY KEY, a INT)
< ok 0 0
> exec
| INSERT INTO omni_mariadb_analysis_backup_other.t VALUES (1, 1)
< ok 1 0
> exec
| INSERT INTO omni_mariadb_analysis_backup_other.u VALUES (1, 1)
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
: string
| u
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
> query
| SELECT `id`, `a` FROM `omni_mariadb_analysis_backup_other`.`u` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< row
: int64 1
: int64 1
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
: string
| t
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| b
> query
| SELECT `id`, `b` FROM `omni_mariadb_analysis_backup_other`.`t` ORDER BY 1, 2
< result
< column "id" "INT"
< column "b" "INT"
< row
: int64 1
: int64 1
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
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
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
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
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
< row
: bytes
| bak_keyless_rows_omni_testbackuprestoreroundtrip
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
| bak_keyless_rows_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
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
| invisible
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| invisible
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
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
| self_join
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| self_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
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
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| invisible
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| keyless_children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| parens
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| paren_peers
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| self_join
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
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
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE `omni_testbackuprestoreroundtrip`.`bak_u_omni_mariadb_analysis_backup_other` LIKE `omni_mariadb_analysis_backup_other`.`u`
< ok 0 0
> exec
|
| INSERT INTO `omni_testbackuprestoreroundtrip`.`bak_u_omni_mariadb_analysis_backup_other` (`id`,`a`)  SELECT DISTINCT `u`.`id`,`u`.`a` FROM omni_mariadb_analysis_backup_other.t JOIN omni_mariadb_analysis_backup_other.u ON omni_mariadb_analysis_backup_other.t.id = omni_mariadb_analysis_backup_other.u.id
< ok 1 0
> exec
| UPDATE omni_mariadb_analysis_backup_other.t JOIN omni_mariadb_analysis_backup_other.u ON omni_mariadb_analysis_backup_other.t.id = omni_mariadb_analysis_backup_other.u.id SET a = 7
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
: string
| u
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
> query
| SELECT `id`, `a` FROM `omni_mariadb_analysis_backup_other`.`u` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< row
: int64 1
: int64 7
> exec
| /*
| Original SQL:
| UPDATE omni_mariadb_analysis_backup_other.t JOIN omni_mariadb_analysis_backup_other.u ON omni_mariadb_analysis_backup_other.t.id = omni_mariadb_analysis_backup_other.u.id SET a = 7;
| */
| INSERT INTO `omni_mariadb_analysis_backup_other`.`u` (`id`, `a`) SELECT `id`, `a` FROM `omni_testbackuprestoreroundtrip`.`bak_u_omni_mariadb_analysis_backup_other` ON DUPLICATE KEY UPDATE `a` = VALUES(`a`)
< ok 2 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
: string
| u
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
> query
| SELECT `id`, `a` FROM `omni_mariadb_analysis_backup_other`.`u` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< row
: int64 1
: int64 1
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
: string
| t
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| b
> query
| SELECT `id`, `b` FROM `omni_mariadb_analysis_backup_other`.`t` ORDER BY 1, 2
< result
< column "id" "INT"
< column "b" "INT"
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
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
: string
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
> query
| SELECT `id` FROM `omni_testbackuprestoreroundtrip`.`nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one` ORDER BY 1
< result
< column "id" "INT"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
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
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
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
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
< row
: bytes
| bak_keyless_rows_omni_testbackuprestoreroundtrip
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
| bak_keyless_rows_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
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
| invisible
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| invisible
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
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
| self_join
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| self_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
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
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| invisible
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| keyless_children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| parens
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| paren_peers
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| self_join
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
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
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE stamped (id INT PRIMARY KEY, v INT, ts TIMESTAMP NOT NULL DEFAULT '2000-01-01 00:00:00' ON UPDATE CURRENT_TIMESTAMP)
< ok 0 0
> exec
| INSERT INTO stamped VALUES (1, 1, '2000-01-01 00:00:00'), (2, 2, '2000-01-01 00:00:00')
< ok 2 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
: string
| stamped
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| v
< row
: bytes
| ts
> query
| SELECT `id`, `v`, `ts` FROM `omni_testbackuprestoreroundtrip`.`stamped` ORDER BY 1, 2, 3
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
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
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
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
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
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
< row
: bytes
| bak_keyless_rows_omni_testbackuprestoreroundtrip
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
| bak_keyless_rows_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
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
| invisible
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| invisible
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
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
| self_join
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| self_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
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
| on update current_timestamp()
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| invisible
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| keyless_children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| parens
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| paren_peers
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| self_join
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| stamped
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
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
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE `omni_testbackuprestoreroundtrip`.`bak_stamped_omni_testbackuprestoreroundtrip` LIKE `omni_testbackuprestoreroundtrip`.`stamped`
< ok 0 0
> exec
|
| INSERT INTO `omni_testbackuprestoreroundtrip`.`bak_stamped_omni_testbackuprestoreroundtrip` (`id`,`v`,`ts`)  SELECT `stamped`.`id`,`stamped`.`v`,`stamped`.`ts` FROM stamped WHERE id = 1
< ok 1 0
> exec
| UPDATE stamped SET v = 5 WHERE id = 1
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
: string
| stamped
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| v
< row
: bytes
| ts
> query
| SELECT `id`, `v`, `ts` FROM `omni_testbackuprestoreroundtrip`.`stamped` ORDER BY 1, 2, 3
< result
< column "id" "INT"
< column "v" "INT"
< column "ts" "TIMESTAMP"
< row
: int64 1
: int64 5
: bytes
| 2026-11-13 00:00:00
< row
: int64 2
: int64 2
: bytes
| 2000-01-01 00:00:00
> exec
| /*
| Original SQL:
| UPDATE stamped SET v = 5 WHERE id = 1;
| */
| INSERT INTO `omni_testbackuprestoreroundtrip`.`stamped` (`id`, `v`, `ts`) SELECT `id`, `v`, `ts` FROM `omni_testbackuprestoreroundtrip`.`bak_stamped_omni_testbackuprestoreroundtrip` ON DUPLICATE KEY UPDATE `v` = VALUES(`v`), `ts` = VALUES(`ts`)
< ok 2 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
: string
| stamped
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| v
< row
: bytes
| ts
> query
| SELECT `id`, `v`, `ts` FROM `omni_testbackuprestoreroundtrip`.`stamped` ORDER BY 1, 2, 3
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
| CREATE TABLE swapped_key (id INT PRIMARY KEY, u INT NOT NULL, UNIQUE KEY uk_u (u))
< ok 0 0
> exec
| INSERT INTO swapped_key VALUES (1, 10), (2, 20)
< ok 2 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
: string
| swapped_key
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| u
> query
| SELECT `id`, `u` FROM `omni_testbackuprestoreroundtrip`.`swapped_key` ORDER BY 1, 2
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
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
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
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
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
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
< row
: bytes
| bak_keyless_rows_omni_testbackuprestoreroundtrip
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
| bak_keyless_rows_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
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
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_stamped_omni_testbackuprestoreroundtrip
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
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| on update current_timestamp()
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
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
| invisible
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| invisible
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
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
| self_join
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| self_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
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
| on update current_timestamp()
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
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| invisible
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| keyless_children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| parens
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| paren_peers
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| self_join
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| stamped
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
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
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE `omni_testbackuprestoreroundtrip`.`bak_swapped_key_omni_testbackuprestoreroundtrip` LIKE `omni_testbackuprestoreroundtrip`.`swapped_key`
< ok 0 0
> exec
|
| INSERT INTO `omni_testbackuprestoreroundtrip`.`bak_swapped_key_omni_testbackuprestoreroundtrip` (`id`,`u`)  (SELECT `swapped_key`.`id`,`swapped_key`.`u` FROM swapped_key WHERE id = 1)
|   UNION DISTINCT
|   (SELECT `swapped_key`.`id`,`swapped_key`.`u` FROM swapped_key WHERE id = 1)
|   UNION DISTINCT
|   (SELECT `swapped_key`.`id`,`swapped_key`.`u` FROM swapped_key WHERE id = 2)
< ok 2 0
> exec
| UPDATE swapped_key SET u = 30 WHERE id = 1
< ok 1 0
> exec
|  UPDATE swapped_key SET u = 10 WHERE id = 2
< ok 1 0
> exec
|  UPDATE swapped_key SET u = 20 WHERE id = 1
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
: string
| swapped_key
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| u
> query
| SELECT `id`, `u` FROM `omni_testbackuprestoreroundtrip`.`swapped_key` ORDER BY 1, 2
< result
< column "id" "INT"
< column "u" "INT"
< row
: int64 1
: int64 20
< row
: int64 2
: int64 10
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
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
: string
| flagged
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
> query
| SELECT `id` FROM `omni_testbackuprestoreroundtrip`.`flagged` ORDER BY 1
< result
< column "id" "INT"
< row
: int64 1
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
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
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
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
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
< row
: bytes
| bak_keyless_rows_omni_testbackuprestoreroundtrip
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
| bak_keyless_rows_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
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
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_stamped_omni_testbackuprestoreroundtrip
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
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| on update current_timestamp()
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
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
| invisible
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| invisible
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
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
| self_join
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| self_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
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
| on update current_timestamp()
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
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| invisible
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| keyless_children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| parens
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| paren_peers
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| self_join
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| stamped
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
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
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
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
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
: string
| fk_parent
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
> query
| SELECT `id` FROM `omni_testbackuprestoreroundtrip`.`fk_parent` ORDER BY 1
< result
< column "id" "INT"
< row
: int64 1
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
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
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
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
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
< row
: bytes
| bak_keyless_rows_omni_testbackuprestoreroundtrip
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
| bak_keyless_rows_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
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
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_stamped_omni_testbackuprestoreroundtrip
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
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| on update current_timestamp()
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
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
| invisible
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| invisible
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
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
| self_join
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| self_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
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
| on update current_timestamp()
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
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: int64 1
: bytes
| parent_id
< row
: bytes
| fk_child
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| invisible
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| keyless_children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| parens
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| paren_peers
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| self_join
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| stamped
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
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
| fk_parent
: bytes
| id
: bytes
| RESTRICT
: bytes
| CASCADE
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
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
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
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
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
: string
| outer_gone
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
> query
| SELECT `id`, `a` FROM `omni_testbackuprestoreroundtrip`.`outer_gone` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< row
: int64 1
: int64 1
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
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
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
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
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
< row
: bytes
| bak_keyless_rows_omni_testbackuprestoreroundtrip
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
| bak_keyless_rows_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
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
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_stamped_omni_testbackuprestoreroundtrip
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
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| on update current_timestamp()
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
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
| invisible
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| invisible
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
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
| self_join
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| self_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
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
| on update current_timestamp()
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
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: int64 1
: bytes
| parent_id
< row
: bytes
| fk_child
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| invisible
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| keyless_children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| outer_gone
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| outer_keep
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| parens
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| paren_peers
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| self_join
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| stamped
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
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
| fk_parent
: bytes
| id
: bytes
| RESTRICT
: bytes
| CASCADE
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
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
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE `omni_testbackuprestoreroundtrip`.`bak_outer_gone_omni_testbackuprestoreroundtrip` LIKE `omni_testbackuprestoreroundtrip`.`outer_gone`
< ok 0 0
> exec
|
| INSERT INTO `omni_testbackuprestoreroundtrip`.`bak_outer_gone_omni_testbackuprestoreroundtrip` (`id`,`a`)  SELECT DISTINCT `t`.`id`,`t`.`a` FROM outer_keep u LEFT JOIN outer_gone t ON t.id = u.id WHERE `t`.`id` IS NOT NULL
< ok 1 0
> exec
| DELETE t FROM outer_keep u LEFT JOIN outer_gone t ON t.id = u.id
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
: string
| outer_gone
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
> query
| SELECT `id`, `a` FROM `omni_testbackuprestoreroundtrip`.`outer_gone` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
> exec
| /*
| Original SQL:
| DELETE t FROM outer_keep u LEFT JOIN outer_gone t ON t.id = u.id;
| */
| INSERT INTO `omni_testbackuprestoreroundtrip`.`outer_gone` (`id`, `a`) SELECT `id`, `a` FROM `omni_testbackuprestoreroundtrip`.`bak_outer_gone_omni_testbackuprestoreroundtrip`
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
: string
| outer_gone
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
> query
| SELECT `id`, `a` FROM `omni_testbackuprestoreroundtrip`.`outer_gone` ORDER BY 1, 2
< result
< column "id" "INT"
< column "a" "INT"
< row
: int64 1
: int64 1
> exec
| CREATE TABLE single_key (id INT PRIMARY KEY, u INT NOT NULL, UNIQUE KEY uk_u (u))
< ok 0 0
> exec
| INSERT INTO single_key VALUES (1, 10), (2, 20)
< ok 2 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
: string
| single_key
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| u
> query
| SELECT `id`, `u` FROM `omni_testbackuprestoreroundtrip`.`single_key` ORDER BY 1, 2
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
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
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
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
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
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
< row
: bytes
| bak_keyless_rows_omni_testbackuprestoreroundtrip
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
| bak_keyless_rows_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_outer_gone_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_outer_gone_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
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
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_stamped_omni_testbackuprestoreroundtrip
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
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| on update current_timestamp()
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
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
| invisible
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| invisible
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
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
| self_join
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| self_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
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
| on update current_timestamp()
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
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_outer_gone_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: int64 1
: bytes
| parent_id
< row
: bytes
| fk_child
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| invisible
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| keyless_children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| outer_gone
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| outer_keep
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| parens
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| paren_peers
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| self_join
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| stamped
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
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
| fk_parent
: bytes
| id
: bytes
| RESTRICT
: bytes
| CASCADE
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
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
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE `omni_testbackuprestoreroundtrip`.`bak_single_key_omni_testbackuprestoreroundtrip` LIKE `omni_testbackuprestoreroundtrip`.`single_key`
< ok 0 0
> exec
|
| INSERT INTO `omni_testbackuprestoreroundtrip`.`bak_single_key_omni_testbackuprestoreroundtrip` (`id`,`u`)  SELECT `single_key`.`id`,`single_key`.`u` FROM single_key WHERE id = 1
< ok 1 0
> exec
| UPDATE single_key SET u = 30 WHERE id = 1
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
: string
| single_key
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| u
> query
| SELECT `id`, `u` FROM `omni_testbackuprestoreroundtrip`.`single_key` ORDER BY 1, 2
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
| INSERT INTO `omni_testbackuprestoreroundtrip`.`single_key` (`id`, `u`) SELECT `id`, `u` FROM `omni_testbackuprestoreroundtrip`.`bak_single_key_omni_testbackuprestoreroundtrip` ON DUPLICATE KEY UPDATE `u` = VALUES(`u`)
< ok 2 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
: string
| single_key
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| u
> query
| SELECT `id`, `u` FROM `omni_testbackuprestoreroundtrip`.`single_key` ORDER BY 1, 2
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
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
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
| SELECT `id`, `note` FROM `omni_testbackuprestoreroundtrip`.`comments` ORDER BY 1, 2
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
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
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
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
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
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
< row
: bytes
| bak_keyless_rows_omni_testbackuprestoreroundtrip
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
| bak_keyless_rows_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_outer_gone_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_outer_gone_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
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
| bak_single_key_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_single_key_omni_testbackuprestoreroundtrip
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
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_stamped_omni_testbackuprestoreroundtrip
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
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| on update current_timestamp()
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
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
| invisible
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| invisible
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
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
| self_join
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| self_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
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
| on update current_timestamp()
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
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_outer_gone_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_single_key_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_single_key_omni_testbackuprestoreroundtrip
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| comments
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: int64 1
: bytes
| parent_id
< row
: bytes
| fk_child
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| invisible
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| keyless_children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| outer_gone
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| outer_keep
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| parens
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| paren_peers
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| self_join
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| stamped
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
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
| fk_parent
: bytes
| id
: bytes
| RESTRICT
: bytes
| CASCADE
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
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
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE `omni_testbackuprestoreroundtrip`.`bak_comments_omni_testbackuprestoreroundtrip` LIKE `omni_testbackuprestoreroundtrip`.`comments`
< ok 0 0
> exec
|
| INSERT INTO `omni_testbackuprestoreroundtrip`.`bak_comments_omni_testbackuprestoreroundtrip` (`id`,`note`)  SELECT `comments`.`id`,`comments`.`note` FROM comments WHERE note = '*/'
< ok 1 0
> exec
| DELETE /* note */ FROM comments WHERE note = '*/'
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
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
| SELECT `id`, `note` FROM `omni_testbackuprestoreroundtrip`.`comments` ORDER BY 1, 2
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
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
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
| SELECT `id`, `note` FROM `omni_testbackuprestoreroundtrip`.`comments` ORDER BY 1, 2
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
| CREATE TABLE returning_rows (id INT PRIMARY KEY, a INT)
< ok 0 0
> exec
| INSERT INTO returning_rows VALUES (1, 1), (2, 2)
< ok 2 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
: string
| returning_rows
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
> query
| SELECT `id`, `a` FROM `omni_testbackuprestoreroundtrip`.`returning_rows` ORDER BY 1, 2
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
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
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
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
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
| bak_comments_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_comments_omni_testbackuprestoreroundtrip
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
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
< row
: bytes
| bak_keyless_rows_omni_testbackuprestoreroundtrip
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
| bak_keyless_rows_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_outer_gone_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_outer_gone_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
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
| bak_single_key_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_single_key_omni_testbackuprestoreroundtrip
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
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_stamped_omni_testbackuprestoreroundtrip
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
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| on update current_timestamp()
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
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
| invisible
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| invisible
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
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
| returning_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| returning_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| self_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
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
| on update current_timestamp()
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
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_comments_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_outer_gone_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_single_key_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_single_key_omni_testbackuprestoreroundtrip
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| comments
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: int64 1
: bytes
| parent_id
< row
: bytes
| fk_child
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| invisible
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| keyless_children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| outer_gone
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| outer_keep
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| parens
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| paren_peers
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| returning_rows
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| self_join
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| stamped
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
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
| fk_parent
: bytes
| id
: bytes
| RESTRICT
: bytes
| CASCADE
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
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
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE `omni_testbackuprestoreroundtrip`.`bak_returning_rows_omni_testbackuprestoreroundtrip` LIKE `omni_testbackuprestoreroundtrip`.`returning_rows`
< ok 0 0
> exec
|
| INSERT INTO `omni_testbackuprestoreroundtrip`.`bak_returning_rows_omni_testbackuprestoreroundtrip` (`id`,`a`)  SELECT `returning_rows`.`id`,`returning_rows`.`a` FROM returning_rows WHERE id = 1
< ok 1 0
> exec
| DELETE FROM returning_rows WHERE id = 1 RETURNING id
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
: string
| returning_rows
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
> query
| SELECT `id`, `a` FROM `omni_testbackuprestoreroundtrip`.`returning_rows` ORDER BY 1, 2
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
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
: string
| returning_rows
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| a
> query
| SELECT `id`, `a` FROM `omni_testbackuprestoreroundtrip`.`returning_rows` ORDER BY 1, 2
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
| omni_testbackuprestoreroundtrip
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
| SELECT `id`, `a`, `g` FROM `omni_testbackuprestoreroundtrip`.`generated_key` ORDER BY 1, 2, 3
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
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
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
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
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
| bak_comments_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_comments_omni_testbackuprestoreroundtrip
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
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
< row
: bytes
| bak_keyless_rows_omni_testbackuprestoreroundtrip
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
| bak_keyless_rows_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_outer_gone_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_outer_gone_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_returning_rows_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_returning_rows_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
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
| bak_single_key_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_single_key_omni_testbackuprestoreroundtrip
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
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_stamped_omni_testbackuprestoreroundtrip
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
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| on update current_timestamp()
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
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
| invisible
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| invisible
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
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
| returning_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| returning_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| self_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
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
| on update current_timestamp()
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
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_comments_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_outer_gone_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_returning_rows_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_single_key_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_single_key_omni_testbackuprestoreroundtrip
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| comments
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: int64 1
: bytes
| parent_id
< row
: bytes
| fk_child
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| generated_key
: bytes
| uk_g
: int64 0
: bytes
| g
< row
: bytes
| invisible
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| keyless_children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| outer_gone
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| outer_keep
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| parens
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| paren_peers
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| returning_rows
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| self_join
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| stamped
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
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
| fk_parent
: bytes
| id
: bytes
| RESTRICT
: bytes
| CASCADE
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
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
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
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
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
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
| SELECT `id`, `a`, `g` FROM `omni_testbackuprestoreroundtrip`.`generated_key` ORDER BY 1, 2, 3
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
| CREATE TABLE nullable_key (u INT NULL, v INT, UNIQUE KEY uk_u (u))
< ok 0 0
> exec
| INSERT INTO nullable_key VALUES (NULL, 1), (2, 2)
< ok 2 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
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
| SELECT `u`, `v` FROM `omni_testbackuprestoreroundtrip`.`nullable_key` ORDER BY 1, 2
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
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
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
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
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
| bak_comments_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_comments_omni_testbackuprestoreroundtrip
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
| bak_generated_key_omni_testbackuprestoreroundtrip
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
| bak_generated_key_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_generated_key_omni_testbackuprestoreroundtrip
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
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
< row
: bytes
| bak_keyless_rows_omni_testbackuprestoreroundtrip
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
| bak_keyless_rows_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_outer_gone_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_outer_gone_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_returning_rows_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_returning_rows_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
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
| bak_single_key_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_single_key_omni_testbackuprestoreroundtrip
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
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_stamped_omni_testbackuprestoreroundtrip
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
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| on update current_timestamp()
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
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
| invisible
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| invisible
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
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
| returning_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| returning_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| self_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
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
| on update current_timestamp()
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
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_comments_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_generated_key_omni_testbackuprestoreroundtrip
: bytes
| uk_g
: int64 0
: bytes
| g
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_outer_gone_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_returning_rows_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_single_key_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_single_key_omni_testbackuprestoreroundtrip
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| comments
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: int64 1
: bytes
| parent_id
< row
: bytes
| fk_child
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| generated_key
: bytes
| uk_g
: int64 0
: bytes
| g
< row
: bytes
| invisible
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| keyless_children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| nullable_key
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| outer_gone
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| outer_keep
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| parens
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| paren_peers
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| returning_rows
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| self_join
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| stamped
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
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
| fk_parent
: bytes
| id
: bytes
| RESTRICT
: bytes
| CASCADE
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
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
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE `omni_testbackuprestoreroundtrip`.`bak_nullable_key_omni_testbackuprestoreroundtrip` LIKE `omni_testbackuprestoreroundtrip`.`nullable_key`
< ok 0 0
> exec
|
| INSERT INTO `omni_testbackuprestoreroundtrip`.`bak_nullable_key_omni_testbackuprestoreroundtrip` (`u`,`v`)  SELECT `nullable_key`.`u`,`nullable_key`.`v` FROM nullable_key WHERE u IS NULL
< ok 1 0
> exec
| UPDATE nullable_key SET v = 3 WHERE u IS NULL
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
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
| SELECT `u`, `v` FROM `omni_testbackuprestoreroundtrip`.`nullable_key` ORDER BY 1, 2
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
| CREATE TABLE generated_default (id INT PRIMARY KEY, a INT, g INT AS (a + 1))
< ok 0 0
> exec
| INSERT INTO generated_default (id, a) VALUES (1, 1), (2, 2)
< ok 2 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
: string
| generated_default
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
| SELECT `id`, `a`, `g` FROM `omni_testbackuprestoreroundtrip`.`generated_default` ORDER BY 1, 2, 3
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
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
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
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
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
| bak_comments_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_comments_omni_testbackuprestoreroundtrip
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
| bak_generated_key_omni_testbackuprestoreroundtrip
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
| bak_generated_key_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_generated_key_omni_testbackuprestoreroundtrip
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
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
< row
: bytes
| bak_keyless_rows_omni_testbackuprestoreroundtrip
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
| bak_keyless_rows_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_nullable_key_omni_testbackuprestoreroundtrip
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
| bak_nullable_key_omni_testbackuprestoreroundtrip
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
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_outer_gone_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_outer_gone_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_returning_rows_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_returning_rows_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
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
| bak_single_key_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_single_key_omni_testbackuprestoreroundtrip
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
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_stamped_omni_testbackuprestoreroundtrip
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
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| on update current_timestamp()
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
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
| invisible
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| invisible
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
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
| returning_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| returning_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| self_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
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
| on update current_timestamp()
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
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_comments_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_generated_key_omni_testbackuprestoreroundtrip
: bytes
| uk_g
: int64 0
: bytes
| g
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| bak_nullable_key_omni_testbackuprestoreroundtrip
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_outer_gone_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_returning_rows_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_single_key_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_single_key_omni_testbackuprestoreroundtrip
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| comments
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: int64 1
: bytes
| parent_id
< row
: bytes
| fk_child
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| generated_default
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| generated_key
: bytes
| uk_g
: int64 0
: bytes
| g
< row
: bytes
| invisible
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| keyless_children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| nullable_key
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| outer_gone
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| outer_keep
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| parens
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| paren_peers
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| returning_rows
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| self_join
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| stamped
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
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
| fk_parent
: bytes
| id
: bytes
| RESTRICT
: bytes
| CASCADE
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
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
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE `omni_testbackuprestoreroundtrip`.`bak_generated_default_omni_testbackuprestoreroundtrip` LIKE `omni_testbackuprestoreroundtrip`.`generated_default`
< ok 0 0
> exec
|
| INSERT INTO `omni_testbackuprestoreroundtrip`.`bak_generated_default_omni_testbackuprestoreroundtrip` (`id`,`a`)  SELECT `generated_default`.`id`,`generated_default`.`a` FROM generated_default WHERE id = 1
< ok 1 0
> exec
| UPDATE generated_default SET generated_default.a = 5, generated_default.g = DEFAULT WHERE id = 1
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
: string
| generated_default
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
| SELECT `id`, `a`, `g` FROM `omni_testbackuprestoreroundtrip`.`generated_default` ORDER BY 1, 2, 3
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
| INSERT INTO `omni_testbackuprestoreroundtrip`.`generated_default` (`id`, `a`) SELECT `id`, `a` FROM `omni_testbackuprestoreroundtrip`.`bak_generated_default_omni_testbackuprestoreroundtrip` ON DUPLICATE KEY UPDATE `a` = VALUES(`a`)
< ok 2 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
: string
| generated_default
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
| SELECT `id`, `a`, `g` FROM `omni_testbackuprestoreroundtrip`.`generated_default` ORDER BY 1, 2, 3
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
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
: string
| user_variable
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
> query
| SELECT `id` FROM `omni_testbackuprestoreroundtrip`.`user_variable` ORDER BY 1
< result
< column "id" "INT"
< row
: int64 1
< row
: int64 2
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
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
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
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
| bak_comments_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_comments_omni_testbackuprestoreroundtrip
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
| bak_generated_default_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_generated_default_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_generated_default_omni_testbackuprestoreroundtrip
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
| bak_generated_key_omni_testbackuprestoreroundtrip
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
| bak_generated_key_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_generated_key_omni_testbackuprestoreroundtrip
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
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
< row
: bytes
| bak_keyless_rows_omni_testbackuprestoreroundtrip
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
| bak_keyless_rows_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_nullable_key_omni_testbackuprestoreroundtrip
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
| bak_nullable_key_omni_testbackuprestoreroundtrip
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
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_outer_gone_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_outer_gone_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_returning_rows_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_returning_rows_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
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
| bak_single_key_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_single_key_omni_testbackuprestoreroundtrip
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
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_stamped_omni_testbackuprestoreroundtrip
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
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| on update current_timestamp()
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
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
| invisible
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| invisible
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
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
| returning_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| returning_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| self_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
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
| on update current_timestamp()
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
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_comments_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_generated_default_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_generated_key_omni_testbackuprestoreroundtrip
: bytes
| uk_g
: int64 0
: bytes
| g
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| bak_nullable_key_omni_testbackuprestoreroundtrip
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_outer_gone_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_returning_rows_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_single_key_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_single_key_omni_testbackuprestoreroundtrip
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| comments
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: int64 1
: bytes
| parent_id
< row
: bytes
| fk_child
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| generated_default
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| generated_key
: bytes
| uk_g
: int64 0
: bytes
| g
< row
: bytes
| invisible
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| keyless_children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| nullable_key
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| outer_gone
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| outer_keep
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| parens
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| paren_peers
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| returning_rows
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| self_join
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| stamped
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| user_variable
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
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
| fk_parent
: bytes
| id
: bytes
| RESTRICT
: bytes
| CASCADE
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
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
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE portioned (id INT, s DATE, e DATE, v INT, PERIOD FOR app_time (s, e), PRIMARY KEY (id, s))
< ok 0 0
> exec
| INSERT INTO portioned VALUES (1, '2020-01-01', '2020-05-01', 1)
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
: string
| portioned
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| s
< row
: bytes
| e
< row
: bytes
| v
> query
| SELECT `id`, `s`, `e`, `v` FROM `omni_testbackuprestoreroundtrip`.`portioned` ORDER BY 1, 2, 3, 4
< result
< column "id" "INT"
< column "s" "DATE"
< column "e" "DATE"
< column "v" "INT"
< row
: int64 1
: bytes
| 2020-01-01
: bytes
| 2020-05-01
: int64 1
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
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
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
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
| bak_comments_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_comments_omni_testbackuprestoreroundtrip
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
| bak_generated_default_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_generated_default_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_generated_default_omni_testbackuprestoreroundtrip
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
| bak_generated_key_omni_testbackuprestoreroundtrip
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
| bak_generated_key_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_generated_key_omni_testbackuprestoreroundtrip
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
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
< row
: bytes
| bak_keyless_rows_omni_testbackuprestoreroundtrip
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
| bak_keyless_rows_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_nullable_key_omni_testbackuprestoreroundtrip
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
| bak_nullable_key_omni_testbackuprestoreroundtrip
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
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_outer_gone_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_outer_gone_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_returning_rows_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_returning_rows_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
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
| bak_single_key_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_single_key_omni_testbackuprestoreroundtrip
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
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_stamped_omni_testbackuprestoreroundtrip
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
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| on update current_timestamp()
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
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
| invisible
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| invisible
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
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
| portioned
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| portioned
: bytes
| s
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| portioned
: bytes
| e
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| portioned
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
| returning_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| returning_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| self_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
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
| on update current_timestamp()
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
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_comments_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_generated_default_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_generated_key_omni_testbackuprestoreroundtrip
: bytes
| uk_g
: int64 0
: bytes
| g
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| bak_nullable_key_omni_testbackuprestoreroundtrip
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_outer_gone_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_returning_rows_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_single_key_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_single_key_omni_testbackuprestoreroundtrip
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| comments
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: int64 1
: bytes
| parent_id
< row
: bytes
| fk_child
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| generated_default
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| generated_key
: bytes
| uk_g
: int64 0
: bytes
| g
< row
: bytes
| invisible
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| keyless_children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| nullable_key
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| outer_gone
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| outer_keep
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| parens
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| paren_peers
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| portioned
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| portioned
: bytes
| PRIMARY
: int64 0
: bytes
| s
< row
: bytes
| returning_rows
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| self_join
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| stamped
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| user_variable
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
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
| fk_parent
: bytes
| id
: bytes
| RESTRICT
: bytes
| CASCADE
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
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
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE portion_delete (id INT, s DATE, e DATE, PERIOD FOR app_time (s, e), PRIMARY KEY (id, s))
< ok 0 0
> exec
| INSERT INTO portion_delete VALUES (1, '2020-01-01', '2020-05-01')
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
: string
| portion_delete
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| s
< row
: bytes
| e
> query
| SELECT `id`, `s`, `e` FROM `omni_testbackuprestoreroundtrip`.`portion_delete` ORDER BY 1, 2, 3
< result
< column "id" "INT"
< column "s" "DATE"
< column "e" "DATE"
< row
: int64 1
: bytes
| 2020-01-01
: bytes
| 2020-05-01
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
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
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
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
| bak_comments_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_comments_omni_testbackuprestoreroundtrip
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
| bak_generated_default_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_generated_default_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_generated_default_omni_testbackuprestoreroundtrip
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
| bak_generated_key_omni_testbackuprestoreroundtrip
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
| bak_generated_key_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_generated_key_omni_testbackuprestoreroundtrip
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
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
< row
: bytes
| bak_keyless_rows_omni_testbackuprestoreroundtrip
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
| bak_keyless_rows_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_nullable_key_omni_testbackuprestoreroundtrip
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
| bak_nullable_key_omni_testbackuprestoreroundtrip
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
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_outer_gone_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_outer_gone_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_returning_rows_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_returning_rows_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
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
| bak_single_key_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_single_key_omni_testbackuprestoreroundtrip
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
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_stamped_omni_testbackuprestoreroundtrip
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
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| on update current_timestamp()
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
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
| invisible
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| invisible
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
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
| portioned
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| portioned
: bytes
| s
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| portioned
: bytes
| e
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| portioned
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
| portion_delete
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| portion_delete
: bytes
| s
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| portion_delete
: bytes
| e
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| returning_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| returning_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| self_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
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
| on update current_timestamp()
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
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_comments_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_generated_default_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_generated_key_omni_testbackuprestoreroundtrip
: bytes
| uk_g
: int64 0
: bytes
| g
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| bak_nullable_key_omni_testbackuprestoreroundtrip
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_outer_gone_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_returning_rows_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_single_key_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_single_key_omni_testbackuprestoreroundtrip
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| comments
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: int64 1
: bytes
| parent_id
< row
: bytes
| fk_child
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| generated_default
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| generated_key
: bytes
| uk_g
: int64 0
: bytes
| g
< row
: bytes
| invisible
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| keyless_children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| nullable_key
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| outer_gone
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| outer_keep
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| parens
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| paren_peers
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| portioned
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| portioned
: bytes
| PRIMARY
: int64 0
: bytes
| s
< row
: bytes
| portion_delete
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| portion_delete
: bytes
| PRIMARY
: int64 0
: bytes
| s
< row
: bytes
| returning_rows
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| self_join
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| stamped
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| user_variable
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
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
| fk_parent
: bytes
| id
: bytes
| RESTRICT
: bytes
| CASCADE
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
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
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
> exec
| CREATE TABLE history (id INT PRIMARY KEY, v INT) WITH SYSTEM VERSIONING
< ok 0 0
> exec
| INSERT INTO history VALUES (1, 1)
< ok 1 0
> exec
| UPDATE history SET v = 2
< ok 1 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
: string
| history
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| id
< row
: bytes
| v
> query
| SELECT `id`, `v` FROM `omni_testbackuprestoreroundtrip`.`history` ORDER BY 1, 2
< result
< column "id" "INT"
< column "v" "INT"
< row
: int64 1
: int64 2
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
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
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
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
| bak_comments_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_comments_omni_testbackuprestoreroundtrip
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
| bak_generated_default_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_generated_default_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_generated_default_omni_testbackuprestoreroundtrip
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
| bak_generated_key_omni_testbackuprestoreroundtrip
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
| bak_generated_key_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_generated_key_omni_testbackuprestoreroundtrip
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
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
< row
: bytes
| bak_keyless_rows_omni_testbackuprestoreroundtrip
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
| bak_keyless_rows_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_limits_omni_testbackuprestoreroundtrip
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
| bak_nullable_key_omni_testbackuprestoreroundtrip
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
| bak_nullable_key_omni_testbackuprestoreroundtrip
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
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_outer_gone_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_outer_gone_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_returning_rows_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_returning_rows_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
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
| bak_single_key_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_single_key_omni_testbackuprestoreroundtrip
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
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_stamped_omni_testbackuprestoreroundtrip
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
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| ts
: bytes
|
: bytes
| NO
: bytes
| on update current_timestamp()
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_updates_omni_testbackuprestoreroundtrip
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
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
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
| history
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| history
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
| invisible
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| invisible
: bytes
| secret
: bytes
|
: bytes
| YES
: bytes
| INVISIBLE
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
| portioned
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| portioned
: bytes
| s
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| portioned
: bytes
| e
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| portioned
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
| portion_delete
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| portion_delete
: bytes
| s
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| portion_delete
: bytes
| e
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| returning_rows
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| returning_rows
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| self_join
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| self_join
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
| on update current_timestamp()
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
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
< row
: bytes
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
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
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| back`tick
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_back`tick_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_comments_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_generated_default_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_generated_key_omni_testbackuprestoreroundtrip
: bytes
| uk_g
: int64 0
: bytes
| g
< row
: bytes
| bak_invisible_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_limits_omni_testbackuprestoreroundtrip
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| bak_nullable_key_omni_testbackuprestoreroundtrip
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| bak_one_to_many_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_outer_gone_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_parens_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_returning_rows_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_self_join_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_single_key_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_single_key_omni_testbackuprestoreroundtrip
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| bak_stamped_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_swapped_key_omni_testbackuprestoreroundtrip
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| bak_updates_omni_testbackuprestoreroundtrip
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| bak_u_omni_mariadb_analysis_backup_other
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| comments
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| fk_child
: bytes
| fk_child_parent
: int64 1
: bytes
| parent_id
< row
: bytes
| fk_child
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| fk_parent
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| flagged
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| flags
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| generated_default
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| generated_key
: bytes
| uk_g
: int64 0
: bytes
| g
< row
: bytes
| history
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| invisible
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| keyless_children
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| limits
: bytes
| uk_b
: int64 0
: bytes
| b
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_one
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn_two
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| nullable_key
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| one_to_many
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| outer_gone
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| outer_keep
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| overlapping_limits
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| parens
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| paren_peers
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| portioned
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| portioned
: bytes
| PRIMARY
: int64 0
: bytes
| s
< row
: bytes
| portion_delete
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| portion_delete
: bytes
| PRIMARY
: int64 0
: bytes
| s
< row
: bytes
| returning_rows
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| self_join
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| single_key
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| stamped
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| swapped_key
: bytes
| uk_u
: int64 0
: bytes
| u
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| updates
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| user_variable
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_testbackuprestoreroundtrip
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
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
| fk_parent
: bytes
| id
: bytes
| RESTRICT
: bytes
| CASCADE
> query
| SELECT TABLE_NAME, COLUMN_NAME, IFNULL(GENERATION_EXPRESSION, ''), IS_NULLABLE, EXTRA FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "IFNULL(GENERATION_EXPRESSION, '')" "LONGTEXT"
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| t
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| t
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
| u
: bytes
| id
: bytes
|
: bytes
| NO
: bytes
|
< row
: bytes
| u
: bytes
| a
: bytes
|
: bytes
| YES
: bytes
|
> query
| SELECT TABLE_NAME, INDEX_NAME, NON_UNIQUE, COLUMN_NAME FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, INDEX_NAME, SEQ_IN_INDEX
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "INDEX_NAME" "VARCHAR"
< column "NON_UNIQUE" "BIGINT"
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| t
: bytes
| PRIMARY
: int64 0
: bytes
| id
< row
: bytes
| u
: bytes
| PRIMARY
: int64 0
: bytes
| id
> query
| SELECT k.TABLE_NAME, k.CONSTRAINT_NAME, k.COLUMN_NAME, k.REFERENCED_TABLE_NAME, k.REFERENCED_COLUMN_NAME, r.UPDATE_RULE, r.DELETE_RULE FROM information_schema.KEY_COLUMN_USAGE k JOIN information_schema.REFERENTIAL_CONSTRAINTS r ON r.CONSTRAINT_SCHEMA = k.CONSTRAINT_SCHEMA AND r.CONSTRAINT_NAME = k.CONSTRAINT_NAME WHERE k.TABLE_SCHEMA = ? AND k.REFERENCED_TABLE_NAME IS NOT NULL ORDER BY k.TABLE_NAME, k.CONSTRAINT_NAME, k.ORDINAL_POSITION
: string
| omni_mariadb_analysis_backup_other
< result
< column "TABLE_NAME" "VARCHAR"
< column "CONSTRAINT_NAME" "VARCHAR"
< column "COLUMN_NAME" "VARCHAR"
< column "REFERENCED_TABLE_NAME" "VARCHAR"
< column "REFERENCED_COLUMN_NAME" "VARCHAR"
< column "UPDATE_RULE" "VARCHAR"
< column "DELETE_RULE" "VARCHAR"
