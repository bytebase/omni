-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_TriggerDDLEchoesItsStatement default#1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE `order` (a INT)
< ok 0 0
> exec
| CREATE TABLE tg (a INT)
< ok 0 0
> exec
| CREATE TRIGGER `my trigger` BEFORE INSERT ON `order` FOR EACH ROW SET NEW.a = 1
< ok 0 0
> query
| SHOW CREATE TRIGGER `my trigger`
< result
< column "Trigger" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< column "Created" "DATETIME"
< row
: bytes
| my trigger
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER `my trigger` BEFORE INSERT ON `order` FOR EACH ROW SET NEW.a = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
: time 2026-11-13T00:00:00Z
> exec
| CREATE TRIGGER tr_bare BEFORE INSERT ON tg FOR EACH ROW SET NEW.a = 1
< ok 0 0
> query
| SHOW CREATE TRIGGER `tr_bare`
< result
< column "Trigger" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< column "Created" "DATETIME"
< row
: bytes
| tr_bare
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER tr_bare BEFORE INSERT ON tg FOR EACH ROW SET NEW.a = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
: time 2026-11-13T00:00:00Z
> exec
| create   trigger   tr_spaced   before  insert on   tg
|   for each row
|    set NEW.a = 1
< ok 0 0
> query
| SHOW CREATE TRIGGER `tr_spaced`
< result
< column "Trigger" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< column "Created" "DATETIME"
< row
: bytes
| tr_spaced
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` trigger   tr_spaced   before  insert on   tg
|   for each row
|    set NEW.a = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
: time 2026-11-13T00:00:00Z
> exec
| CREATE TRIGGER IF NOT EXISTS tr_ifne BEFORE INSERT ON tg FOR EACH ROW SET NEW.a = 1
< ok 0 0
> query
| SHOW CREATE TRIGGER `tr_ifne`
< result
< column "Trigger" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< column "Created" "DATETIME"
< row
: bytes
| tr_ifne
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER IF NOT EXISTS tr_ifne BEFORE INSERT ON tg FOR EACH ROW SET NEW.a = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
: time 2026-11-13T00:00:00Z
> exec
| CREATE OR REPLACE TRIGGER tr_replace BEFORE INSERT ON tg FOR EACH ROW SET NEW.a = 1
< ok 0 0
> query
| SHOW CREATE TRIGGER `tr_replace`
< result
< column "Trigger" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< column "Created" "DATETIME"
< row
: bytes
| tr_replace
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER tr_replace BEFORE INSERT ON tg FOR EACH ROW SET NEW.a = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
: time 2026-11-13T00:00:00Z
> exec
| CREATE TRIGGER tr_follows BEFORE INSERT ON tg FOR EACH ROW FOLLOWS tr_bare SET NEW.a = 1
< ok 0 0
> query
| SHOW CREATE TRIGGER `tr_follows`
< result
< column "Trigger" "VARCHAR"
< column "sql_mode" "VARCHAR"
< column "SQL Original Statement" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< column "Database Collation" "VARCHAR"
< column "Created" "DATETIME"
< row
: bytes
| tr_follows
: bytes
| STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
: bytes
| CREATE DEFINER=`root`@`%` TRIGGER tr_follows BEFORE INSERT ON tg FOR EACH ROW SET NEW.a = 1
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
: bytes
| utf8mb4_uca1400_ai_ci
: time 2026-11-13T00:00:00Z
