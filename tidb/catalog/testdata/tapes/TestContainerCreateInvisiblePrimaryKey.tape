-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestContainerCreateInvisiblePrimaryKey default#1
> exec
| DROP DATABASE IF EXISTS omni_invisible_pk
< ok 0 0
> exec
| CREATE DATABASE omni_invisible_pk
< ok 0 0
> exec
| USE omni_invisible_pk
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT, PRIMARY KEY (a) INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT, CONSTRAINT pk PRIMARY KEY (a) INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT, PRIMARY KEY (a) COMMENT 'x' INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT, PRIMARY KEY USING BTREE (a) INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT, PRIMARY KEY (a) INVISIBLE CLUSTERED)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (v VARCHAR(10), PRIMARY KEY (v) INVISIBLE COMMENT 'x' NONCLUSTERED)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT, PRIMARY KEY (a, b) INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT AUTO_INCREMENT, PRIMARY KEY (a) INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT, PRIMARY KEY (a) INVISIBLE, UNIQUE KEY u (a))
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT, PRIMARY KEY (a) VISIBLE INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT, PRIMARY KEY (a) INVISIBLE VISIBLE)
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
|   `a` int NOT NULL,
|   PRIMARY KEY (`a`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT, PRIMARY KEY (a) VISIBLE)
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
|   `a` int NOT NULL,
|   PRIMARY KEY (`a`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, CONSTRAINT u UNIQUE (a) INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) USING HASH INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY u (a, b) INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) INVISIBLE VISIBLE)
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
|   `a` int NOT NULL,
|   UNIQUE KEY `u` (`a`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY u (a) INVISIBLE, UNIQUE KEY v (b))
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY v (b), UNIQUE KEY u (a) INVISIBLE)
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
|   `a` int NOT NULL,
|   `b` int NOT NULL,
|   UNIQUE KEY `v` (`b`),
|   UNIQUE KEY `u` (`a`) /*!80000 INVISIBLE */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, KEY k (b), UNIQUE KEY u (a) INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, KEY k (a) INVISIBLE, UNIQUE KEY v (b))
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
|   `a` int NOT NULL,
|   `b` int NOT NULL,
|   KEY `k` (`a`) /*!80000 INVISIBLE */,
|   UNIQUE KEY `v` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT, UNIQUE KEY u (a) INVISIBLE)
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
|   `a` int DEFAULT NULL,
|   UNIQUE KEY `u` (`a`) /*!80000 INVISIBLE */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT, UNIQUE KEY u (a, b) INVISIBLE)
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
|   `a` int NOT NULL,
|   `b` int DEFAULT NULL,
|   UNIQUE KEY `u` (`a`,`b`) /*!80000 INVISIBLE */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT, UNIQUE KEY v (b), UNIQUE KEY u (a) INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT, PRIMARY KEY (b), UNIQUE KEY u (a) INVISIBLE)
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
|   `a` int NOT NULL,
|   `b` int NOT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `u` (`a`) /*!80000 INVISIBLE */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY u (a) INVISIBLE, PRIMARY KEY (b))
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
|   `a` int NOT NULL,
|   `b` int NOT NULL,
|   PRIMARY KEY (`b`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `u` (`a`) /*!80000 INVISIBLE */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT PRIMARY KEY, b INT NOT NULL, UNIQUE KEY u (b) INVISIBLE)
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
|   `a` int NOT NULL,
|   `b` int NOT NULL,
|   PRIMARY KEY (`a`) /*T![clustered_index] CLUSTERED */,
|   UNIQUE KEY `u` (`b`) /*!80000 INVISIBLE */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY u (a) INVISIBLE, UNIQUE KEY v (b), PRIMARY KEY (a, b))
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
|   `a` int NOT NULL,
|   `b` int NOT NULL,
|   UNIQUE KEY `u` (`a`) /*!80000 INVISIBLE */,
|   UNIQUE KEY `v` (`b`),
|   PRIMARY KEY (`a`,`b`) /*T![clustered_index] CLUSTERED */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL UNIQUE, b INT NOT NULL, UNIQUE KEY u (b) INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) INVISIBLE, b INT NOT NULL UNIQUE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a SERIAL, b INT NOT NULL, UNIQUE KEY u (b) INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT NOT NULL UNIQUE KEY, UNIQUE KEY u (a) INVISIBLE)
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
|   `a` int DEFAULT NULL,
|   `b` int NOT NULL,
|   UNIQUE KEY `u` (`a`) /*!80000 INVISIBLE */,
|   UNIQUE KEY `b` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a), b INT NOT NULL UNIQUE, UNIQUE KEY v (b) INVISIBLE)
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
|   `a` int NOT NULL,
|   `b` int NOT NULL,
|   UNIQUE KEY `u` (`a`),
|   UNIQUE KEY `v` (`b`) /*!80000 INVISIBLE */,
|   UNIQUE KEY `b` (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a VARCHAR(10) NOT NULL, UNIQUE KEY u (a(3)) INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a VARCHAR(10) NOT NULL, b INT NOT NULL, UNIQUE KEY u (a(3)), UNIQUE KEY v (b) INVISIBLE)
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
|   `a` varchar(10) NOT NULL,
|   `b` int NOT NULL,
|   UNIQUE KEY `u` (`a`(3)),
|   UNIQUE KEY `v` (`b`) /*!80000 INVISIBLE */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a VARCHAR(10) NOT NULL, b INT NOT NULL, UNIQUE KEY u (a(3)) INVISIBLE, UNIQUE KEY v (b))
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a TEXT NOT NULL, b INT NOT NULL, UNIQUE KEY u (a(10)), UNIQUE KEY v (b) INVISIBLE)
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
|   `a` text NOT NULL,
|   `b` int NOT NULL,
|   UNIQUE KEY `u` (`a`(10)),
|   UNIQUE KEY `v` (`b`) /*!80000 INVISIBLE */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT AUTO_INCREMENT, UNIQUE KEY u (a) INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, g INT AS (a + 1) STORED NOT NULL, UNIQUE KEY u (g) INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, g INT AS (a + 1) VIRTUAL NOT NULL, UNIQUE KEY u (g) INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u ((a + 1)) INVISIBLE)
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
|   `a` int NOT NULL,
|   UNIQUE KEY `u` ((`a` + 1)) /*!80000 INVISIBLE */
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY f ((a + 1)), UNIQUE KEY v (b) INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP DATABASE IF EXISTS omni_invisible_pk
< ok 0 0
