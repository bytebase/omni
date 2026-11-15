-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestOracle_MigrationEscapesCommentsForSQLMode default#1
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_migrationescapescommentsforsqlmode`
< ok 0 0
> exec
| CREATE DATABASE `omni_testoracle_migrationescapescommentsforsqlmode`
< ok 1 0
> exec
| USE `omni_testoracle_migrationescapescommentsforsqlmode`
< ok 0 0
> exec
| SET @saved_sql_mode = @@sql_mode
< ok 0 0
> exec
| SET sql_mode = 'NO_BACKSLASH_ESCAPES'
< ok 0 0
> exec
| CREATE FUNCTION `omni_testoracle_migrationescapescommentsforsqlmode`.`f`() RETURNS int(11)
|     DETERMINISTIC
|     COMMENT 'a\b'
| RETURN 1
< ok 0 0
> exec
| SET sql_mode = @saved_sql_mode
< ok 0 0
> exec
| SET @saved_sql_mode = @@sql_mode
< ok 0 0
> exec
| SET sql_mode = 'NO_BACKSLASH_ESCAPES'
< ok 0 0
> exec
| CREATE EVENT `omni_testoracle_migrationescapescommentsforsqlmode`.`e` ON SCHEDULE EVERY 1 DAY STARTS '2030-01-01 00:00:00' ON COMPLETION NOT PRESERVE ENABLE COMMENT 'c\d' DO SELECT 1
< ok 0 0
> exec
| SET sql_mode = @saved_sql_mode
< ok 0 0
> query
| SELECT ROUTINE_COMMENT FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = DATABASE() AND ROUTINE_NAME = 'f'
< result
< column "ROUTINE_COMMENT" "TEXT"
< row
: bytes
| a\b
> query
| SELECT EVENT_COMMENT FROM information_schema.EVENTS WHERE EVENT_SCHEMA = DATABASE() AND EVENT_NAME = 'e'
< result
< column "EVENT_COMMENT" "VARCHAR"
< row
: bytes
| c\d
