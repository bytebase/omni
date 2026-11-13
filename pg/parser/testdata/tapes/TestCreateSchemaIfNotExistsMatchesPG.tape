-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestCreateSchemaIfNotExistsMatchesPG default#1
> exec
| SELECT 1/0;
| CREATE SCHEMA IF NOT EXISTS s CREATE TABLE n (id int)
< error
| 0A000 CREATE SCHEMA IF NOT EXISTS cannot include schema elements
| position: 43
> exec
| SELECT 1/0;
| CREATE SCHEMA IF NOT EXISTS AUTHORIZATION r CREATE TABLE n (id int)
< error
| 0A000 CREATE SCHEMA IF NOT EXISTS cannot include schema elements
| position: 57
> exec
| SELECT 1/0;
| CREATE SCHEMA IF NOT EXISTS s AUTHORIZATION r CREATE VIEW v AS SELECT 1 CREATE SEQUENCE q
< error
| 0A000 CREATE SCHEMA IF NOT EXISTS cannot include schema elements
| position: 59
> exec
| SELECT 1/0;
| CREATE SCHEMA IF NOT EXISTS s
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE SCHEMA IF NOT EXISTS s AUTHORIZATION r
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE SCHEMA IF NOT EXISTS AUTHORIZATION r
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE SCHEMA s CREATE TABLE n (id int)
< error
| 22012 division by zero
> exec
| SELECT 1/0;
| CREATE SCHEMA AUTHORIZATION r CREATE TABLE n (id int)
< error
| 22012 division by zero
