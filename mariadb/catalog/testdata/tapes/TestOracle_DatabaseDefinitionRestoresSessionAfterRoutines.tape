-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestOracle_DatabaseDefinitionRestoresSessionAfterRoutines default#1
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_databasedefinitionrestoressessionafterroutines`
< ok 0 0
> exec
| CREATE DATABASE `omni_testoracle_databasedefinitionrestoressessionafterroutines`
< ok 1 0
> exec
| USE `omni_testoracle_databasedefinitionrestoressessionafterroutines`
< ok 0 0
> exec
| SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0
< ok 0 0
> exec
| SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0
< ok 0 0
> exec
| SET @OLD_EXPLICIT_DEFAULTS_FOR_TIMESTAMP=@@EXPLICIT_DEFAULTS_FOR_TIMESTAMP, EXPLICIT_DEFAULTS_FOR_TIMESTAMP=1
< ok 0 0
> exec
| SET @saved_cs_client = @@character_set_client, @saved_cs_results = @@character_set_results, @saved_col_connection = @@collation_connection, @saved_sql_mode = @@sql_mode
< ok 0 0
> exec
| SET character_set_client = latin1
< ok 0 0
> exec
| SET character_set_results = latin1
< ok 0 0
> exec
| SET collation_connection = latin1_swedish_ci
< ok 0 0
> exec
| SET sql_mode = 'NO_BACKSLASH_ESCAPES'
< ok 0 0
> exec
| CREATE FUNCTION `f`() RETURNS int(11)
| RETURN 1
< ok 0 0
> exec
| SET character_set_client = @saved_cs_client, character_set_results = @saved_cs_results, collation_connection = @saved_col_connection, sql_mode = @saved_sql_mode
< ok 0 0
> exec
| CREATE OR REPLACE VIEW `v` AS select 'a\\b' AS `x`
< ok 0 0
> exec
| SET EXPLICIT_DEFAULTS_FOR_TIMESTAMP=@OLD_EXPLICIT_DEFAULTS_FOR_TIMESTAMP
< ok 0 0
> exec
| SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS
< ok 0 0
> exec
| SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS
< ok 0 0
> query
| SELECT x FROM v
< result
< column "x" "VARCHAR"
< row
: bytes
| a\b
