-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_RoutineRecreatePreservesSQLMode default#1
> exec
| DROP DATABASE IF EXISTS byt9832_sc_apply_rtn
< ok 0 0
> exec
| CREATE DATABASE byt9832_sc_apply_rtn DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE byt9832_sc_apply_rtn
< ok 0 0
> exec
| SET SESSION sql_mode = 'PIPES_AS_CONCAT,NO_BACKSLASH_ESCAPES'
< ok 0 0
> exec
| CREATE FUNCTION f(a INT) RETURNS INT DETERMINISTIC RETURN a + 1
< ok 0 0
> exec
| SET SESSION sql_mode = 'STRICT_TRANS_TABLES,NO_ENGINE_SUBSTITUTION'
< ok 0 0
> exec
| DROP FUNCTION `byt9832_sc_apply_rtn`.`f`
< ok 0 0
> exec
| SET @saved_sql_mode = @@sql_mode;
| SET sql_mode = 'PIPES_AS_CONCAT,NO_BACKSLASH_ESCAPES';
| SET @saved_cs_client = @@character_set_client;
| SET character_set_client = utf8mb4;
| SET @saved_coll_conn = @@collation_connection;
| SET collation_connection = utf8mb4_0900_ai_ci;
| CREATE FUNCTION `byt9832_sc_apply_rtn`.`f`(`a` INT) RETURNS int
|     DETERMINISTIC
| RETURN a + 2;
| SET collation_connection = @saved_coll_conn;
| SET character_set_client = @saved_cs_client;
| SET sql_mode = @saved_sql_mode
< ok 0 0
> query
| SELECT SQL_MODE FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA=? AND ROUTINE_NAME=?
: string
| byt9832_sc_apply_rtn
: string
| f
< result
< column "SQL_MODE" "TEXT"
< row
: bytes
| PIPES_AS_CONCAT,NO_BACKSLASH_ESCAPES
