-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_EventAlterPreservesSQLMode default#1
> exec
| DROP DATABASE IF EXISTS byt9832_sc_apply_evt
< ok 0 0
> exec
| CREATE DATABASE byt9832_sc_apply_evt DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE byt9832_sc_apply_evt
< ok 0 0
> exec
| SET GLOBAL event_scheduler = OFF
< ok 0 0
> exec
| SET SESSION sql_mode = 'PIPES_AS_CONCAT'
< ok 0 0
> exec
| SET SESSION time_zone = '+08:00'
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE EVERY 1 HOUR DISABLE DO SET @x = 1
< ok 0 0
> exec
| SET SESSION sql_mode = 'STRICT_TRANS_TABLES,NO_ENGINE_SUBSTITUTION'
< ok 0 0
> exec
| SET @saved_sql_mode = @@sql_mode;
| SET sql_mode = 'PIPES_AS_CONCAT';
| SET @saved_cs_client = @@character_set_client;
| SET character_set_client = utf8mb4;
| SET @saved_coll_conn = @@collation_connection;
| SET collation_connection = utf8mb4_0900_ai_ci;
| SET @saved_time_zone = @@time_zone;
| SET time_zone = '+08:00';
| ALTER EVENT `e` ON SCHEDULE EVERY 1 HOUR ON COMPLETION NOT PRESERVE DISABLE COMMENT '' DO SET @x = 2;
| SET time_zone = @saved_time_zone;
| SET collation_connection = @saved_coll_conn;
| SET character_set_client = @saved_cs_client;
| SET sql_mode = @saved_sql_mode
< ok 0 0
> query
| SELECT SQL_MODE, TIME_ZONE FROM information_schema.EVENTS WHERE EVENT_SCHEMA=? AND EVENT_NAME=?
: string
| byt9832_sc_apply_evt
: string
| e
< result
< column "SQL_MODE" "TEXT"
< column "TIME_ZONE" "VARCHAR"
< row
: bytes
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| SYSTEM
