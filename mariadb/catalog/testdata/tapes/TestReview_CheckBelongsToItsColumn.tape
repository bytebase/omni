-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestReview_CheckBelongsToItsColumn default#1
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
| CREATE TABLE t (k INT, a INT, CONSTRAINT k CHECK (a > 0))
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
|   `k` int(11) DEFAULT NULL,
|   `a` int(11) DEFAULT NULL,
|   CONSTRAINT `k` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT, CONSTRAINT a CHECK (a > 0))
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
|   `a` int(11) DEFAULT NULL,
|   CONSTRAINT `a` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (k INT, a INT, CONSTRAINT k CHECK (a > 0))
< ok 0 0
> exec
| ALTER TABLE t RENAME COLUMN k TO z
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
|   `z` int(11) DEFAULT NULL,
|   `a` int(11) DEFAULT NULL,
|   CONSTRAINT `k` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT CHECK (b > 0), CONSTRAINT b CHECK (a > 0))
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
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL CHECK (`b` > 0),
|   CONSTRAINT `b` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT, CONSTRAINT b CHECK (a > 0), b INT CHECK (b > 0))
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
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL CHECK (`b` > 0),
|   CONSTRAINT `b` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (b INT CHECK (b > 0), c INT, CONSTRAINT b CHECK (c > 0))
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
|   `b` int(11) DEFAULT NULL CHECK (`b` > 0),
|   `c` int(11) DEFAULT NULL,
|   CONSTRAINT `b` CHECK (`c` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (j JSON, a INT, CONSTRAINT k CHECK (a > 0))
< ok 0 0
> exec
| ALTER TABLE t CHANGE j k JSON
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
|   `k` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (json_valid(`k`)),
|   `a` int(11) DEFAULT NULL,
|   CONSTRAINT `k` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (j JSON, a INT, CONSTRAINT k CHECK (a > 0))
< ok 0 0
> exec
| ALTER TABLE t RENAME COLUMN j TO k
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
|   `k` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (json_valid(`k`)),
|   `a` int(11) DEFAULT NULL,
|   CONSTRAINT `k` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (j JSON, a INT, CONSTRAINT j CHECK (a > 0))
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
|   `j` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (json_valid(`j`)),
|   `a` int(11) DEFAULT NULL,
|   CONSTRAINT `j` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (j JSON, a INT)
< ok 0 0
> exec
| ALTER TABLE t ADD CONSTRAINT j CHECK (a > 0)
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
|   `j` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (json_valid(`j`)),
|   `a` int(11) DEFAULT NULL,
|   CONSTRAINT `j` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (j LONGTEXT, a INT, CONSTRAINT j CHECK (a > 0))
< ok 0 0
> exec
| ALTER TABLE t MODIFY COLUMN j JSON
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
|   `j` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (json_valid(`j`)),
|   `a` int(11) DEFAULT NULL,
|   CONSTRAINT `j` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT, CONSTRAINT k CHECK (a > 0))
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN k JSON
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
|   `a` int(11) DEFAULT NULL,
|   `k` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (json_valid(`k`)),
|   CONSTRAINT `k` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (j JSON CHECK (j IS NOT NULL))
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
|   `j` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (`j` is not null)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (j JSON CHECK (j IS NOT NULL), a INT, CONSTRAINT ck CHECK (a > 0))
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
|   `j` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (`j` is not null),
|   `a` int(11) DEFAULT NULL,
|   CONSTRAINT `ck` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT CHECK (a > 0))
< ok 0 0
> exec
| ALTER TABLE t CHANGE a z INT
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
|   `z` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP DATABASE IF EXISTS testdb
< ok 1 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (j JSON)
< ok 0 0
> exec
| ALTER TABLE t CHANGE j z JSON
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
|   `z` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL CHECK (json_valid(`z`))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
