-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestDropPeriodIfExistsOracle default#1
> exec
| DROP TABLE IF EXISTS ta
< ok 0 0
> exec
| CREATE TABLE ta (id INT, s DATE, e DATE, PERIOD FOR app(s, e))
< ok 0 0
> exec
| ALTER TABLE ta DROP PERIOD IF EXISTS FOR app
< ok 0 0
> query
| SHOW CREATE TABLE ta
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ta
: bytes
| CREATE TABLE `ta` (
|   `id` int(11) DEFAULT NULL,
|   `s` date NOT NULL,
|   `e` date NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS ta
< ok 0 0
> exec
| CREATE TABLE ta (id INT, s DATE, e DATE, PERIOD FOR app(s, e))
< ok 0 0
> exec
| ALTER TABLE ta DROP PERIOD IF EXISTS FOR APP
< ok 0 0
> query
| SHOW CREATE TABLE ta
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ta
: bytes
| CREATE TABLE `ta` (
|   `id` int(11) DEFAULT NULL,
|   `s` date NOT NULL,
|   `e` date NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS ta
< ok 0 0
> exec
| CREATE TABLE ta (id INT, s DATE, e DATE, PERIOD FOR app(s, e))
< ok 0 0
> exec
| ALTER TABLE ta DROP PERIOD IF EXISTS FOR `app`
< ok 0 0
> query
| SHOW CREATE TABLE ta
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ta
: bytes
| CREATE TABLE `ta` (
|   `id` int(11) DEFAULT NULL,
|   `s` date NOT NULL,
|   `e` date NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS ta
< ok 0 0
> exec
| CREATE TABLE ta (id INT, s DATE, e DATE, PERIOD FOR app(s, e))
< ok 0 0
> exec
| ALTER TABLE ta DROP PERIOD IF EXISTS FOR nope
< ok 0 0
> query
| SHOW CREATE TABLE ta
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ta
: bytes
| CREATE TABLE `ta` (
|   `id` int(11) DEFAULT NULL,
|   `s` date NOT NULL,
|   `e` date NOT NULL,
|   PERIOD FOR `app` (`s`, `e`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS ta
< ok 0 0
> exec
| CREATE TABLE ta (id INT, s DATE, e DATE, PERIOD FOR app(s, e))
< ok 0 0
> exec
| ALTER TABLE ta DROP PERIOD FOR nope
< error
| 1091 42000 Can't DROP PERIOD `nope`; check that it exists
> query
| SHOW CREATE TABLE ta
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ta
: bytes
| CREATE TABLE `ta` (
|   `id` int(11) DEFAULT NULL,
|   `s` date NOT NULL,
|   `e` date NOT NULL,
|   PERIOD FOR `app` (`s`, `e`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS ta
< ok 0 0
> exec
| CREATE TABLE ta (id INT, s DATE, e DATE)
< ok 0 0
> exec
| ALTER TABLE ta DROP PERIOD IF EXISTS FOR app
< ok 0 0
> query
| SHOW CREATE TABLE ta
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ta
: bytes
| CREATE TABLE `ta` (
|   `id` int(11) DEFAULT NULL,
|   `s` date DEFAULT NULL,
|   `e` date DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS ta
< ok 0 0
> exec
| CREATE TABLE ta (id INT, s DATE, e DATE)
< ok 0 0
> exec
| ALTER TABLE ta DROP PERIOD FOR app
< error
| 1091 42000 Can't DROP PERIOD `app`; check that it exists
> query
| SHOW CREATE TABLE ta
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ta
: bytes
| CREATE TABLE `ta` (
|   `id` int(11) DEFAULT NULL,
|   `s` date DEFAULT NULL,
|   `e` date DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS ta
< ok 0 0
> exec
| CREATE TABLE ta (id INT, s DATE, e DATE, PERIOD FOR app(s, e))
< ok 0 0
> exec
| ALTER TABLE ta DROP PERIOD IF EXISTS FOR `SYSTEM_TIME`
< ok 0 0
> query
| SHOW CREATE TABLE ta
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ta
: bytes
| CREATE TABLE `ta` (
|   `id` int(11) DEFAULT NULL,
|   `s` date NOT NULL,
|   `e` date NOT NULL,
|   PERIOD FOR `app` (`s`, `e`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS ta
< ok 0 0
> exec
| CREATE TABLE ta (id INT) WITH SYSTEM VERSIONING
< ok 0 0
> exec
| ALTER TABLE ta DROP PERIOD IF EXISTS FOR `SYSTEM_TIME`
< ok 0 0
> query
| SHOW CREATE TABLE ta
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ta
: bytes
| CREATE TABLE `ta` (
|   `id` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci WITH SYSTEM VERSIONING
> exec
| DROP TABLE IF EXISTS ta
< ok 0 0
> exec
| CREATE TABLE ta (id INT, rs TIMESTAMP(6) GENERATED ALWAYS AS ROW START, re TIMESTAMP(6) GENERATED ALWAYS AS ROW END, PERIOD FOR SYSTEM_TIME(rs, re)) WITH SYSTEM VERSIONING
< ok 0 0
> exec
| ALTER TABLE ta DROP PERIOD IF EXISTS FOR `SYSTEM_TIME`
< ok 0 0
> query
| SHOW CREATE TABLE ta
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ta
: bytes
| CREATE TABLE `ta` (
|   `id` int(11) DEFAULT NULL,
|   `rs` timestamp(6) GENERATED ALWAYS AS ROW START,
|   `re` timestamp(6) GENERATED ALWAYS AS ROW END,
|   PERIOD FOR SYSTEM_TIME (`rs`, `re`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci WITH SYSTEM VERSIONING
> exec
| DROP TABLE IF EXISTS ta
< ok 0 0
> exec
| CREATE TABLE ta (id INT, s DATE, e DATE)
< ok 0 0
> exec
| ALTER TABLE ta DROP PERIOD IF EXISTS FOR `SYSTEM_TIME`
< ok 0 0
> query
| SHOW CREATE TABLE ta
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ta
: bytes
| CREATE TABLE `ta` (
|   `id` int(11) DEFAULT NULL,
|   `s` date DEFAULT NULL,
|   `e` date DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS ta
< ok 0 0
> exec
| CREATE TABLE ta (id INT, s DATE, e DATE, PERIOD FOR app(s, e)) WITH SYSTEM VERSIONING
< ok 0 0
> exec
| ALTER TABLE ta DROP PERIOD IF EXISTS FOR app
< ok 0 0
> query
| SHOW CREATE TABLE ta
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ta
: bytes
| CREATE TABLE `ta` (
|   `id` int(11) DEFAULT NULL,
|   `s` date NOT NULL,
|   `e` date NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci WITH SYSTEM VERSIONING
> exec
| DROP TABLE IF EXISTS ta
< ok 0 0
> exec
| CREATE TABLE ta (id INT, s DATE, e DATE, PERIOD FOR app(s, e), UNIQUE KEY k (id, app WITHOUT OVERLAPS))
< ok 0 0
> exec
| ALTER TABLE ta DROP PERIOD IF EXISTS FOR app
< error
| 4156 HY000 Period `app` is not found in table
> query
| SHOW CREATE TABLE ta
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ta
: bytes
| CREATE TABLE `ta` (
|   `id` int(11) DEFAULT NULL,
|   `s` date NOT NULL,
|   `e` date NOT NULL,
|   PERIOD FOR `app` (`s`, `e`),
|   UNIQUE KEY `k` (`id`,`app` WITHOUT OVERLAPS)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS ta
< ok 0 0
> exec
| CREATE TABLE ta (id INT, s DATE, e DATE, PERIOD FOR app(s, e), UNIQUE KEY k (id, app WITHOUT OVERLAPS))
< ok 0 0
> exec
| ALTER TABLE ta DROP PERIOD IF EXISTS FOR app, DROP KEY k
< ok 0 0
> query
| SHOW CREATE TABLE ta
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ta
: bytes
| CREATE TABLE `ta` (
|   `id` int(11) DEFAULT NULL,
|   `s` date NOT NULL,
|   `e` date NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS ta
< ok 0 0
> exec
| CREATE TABLE ta (id INT, s DATE, e DATE, PERIOD FOR app(s, e), UNIQUE KEY k (id, app WITHOUT OVERLAPS))
< ok 0 0
> exec
| ALTER TABLE ta DROP PERIOD IF EXISTS FOR nope
< ok 0 0
> query
| SHOW CREATE TABLE ta
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ta
: bytes
| CREATE TABLE `ta` (
|   `id` int(11) DEFAULT NULL,
|   `s` date NOT NULL,
|   `e` date NOT NULL,
|   PERIOD FOR `app` (`s`, `e`),
|   UNIQUE KEY `k` (`id`,`app` WITHOUT OVERLAPS)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS ta
< ok 0 0
> exec
| CREATE TABLE ta (id INT, s DATE, e DATE, PERIOD FOR app(s, e))
< ok 0 0
> exec
| ALTER TABLE ta DROP PERIOD IF EXISTS FOR app, DROP PERIOD FOR app
< error
| 1091 42000 Can't DROP PERIOD `app`; check that it exists
> query
| SHOW CREATE TABLE ta
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ta
: bytes
| CREATE TABLE `ta` (
|   `id` int(11) DEFAULT NULL,
|   `s` date NOT NULL,
|   `e` date NOT NULL,
|   PERIOD FOR `app` (`s`, `e`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS ta
< ok 0 0
> exec
| CREATE TABLE ta (id INT, s DATE, e DATE, PERIOD FOR app(s, e))
< ok 0 0
> exec
| ALTER TABLE ta DROP PERIOD FOR app, DROP PERIOD IF EXISTS FOR app
< ok 0 0
> query
| SHOW CREATE TABLE ta
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ta
: bytes
| CREATE TABLE `ta` (
|   `id` int(11) DEFAULT NULL,
|   `s` date NOT NULL,
|   `e` date NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS ta
< ok 0 0
> exec
| CREATE TABLE ta (id INT, s DATE, e DATE, PERIOD FOR app(s, e))
< ok 0 0
> exec
| ALTER TABLE ta DROP PERIOD IF EXISTS FOR app, DROP PERIOD IF EXISTS FOR app
< ok 0 0
> query
| SHOW CREATE TABLE ta
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ta
: bytes
| CREATE TABLE `ta` (
|   `id` int(11) DEFAULT NULL,
|   `s` date NOT NULL,
|   `e` date NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS ta
< ok 0 0
> exec
| CREATE TABLE ta (id INT, s DATE, e DATE, PERIOD FOR app(s, e))
< ok 0 0
> exec
| ALTER TABLE ta DROP PERIOD IF EXISTS FOR app, ADD PERIOD FOR p2(s, e)
< ok 0 0
> query
| SHOW CREATE TABLE ta
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| ta
: bytes
| CREATE TABLE `ta` (
|   `id` int(11) DEFAULT NULL,
|   `s` date NOT NULL,
|   `e` date NOT NULL,
|   PERIOD FOR `p2` (`s`, `e`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
