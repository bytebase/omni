-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestOracle_DatabaseDefinitionRestoresTimestampsWithoutImplicitDefaults default#1
> exec
| DROP DATABASE IF EXISTS `omni_testoracle_databasedefinitionrestorestimestampswit_1ebe6132`
< ok 0 0
> exec
| CREATE DATABASE `omni_testoracle_databasedefinitionrestorestimestampswit_1ebe6132`
< ok 1 0
> exec
| USE `omni_testoracle_databasedefinitionrestorestimestampswit_1ebe6132`
< ok 0 0
> exec
| SET explicit_defaults_for_timestamp = OFF
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
| CREATE TABLE `t` (
|   `ts` timestamp NOT NULL
| ) ENGINE=InnoDB
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
| SHOW CREATE TABLE t
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t
: bytes
| CREATE TABLE `t` (
|   `ts` timestamp NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT @@explicit_defaults_for_timestamp
< result
< column "@@explicit_defaults_for_timestamp" "BIGINT"
< row
: int64 0
