-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestScenario_C20 default#1
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
| CREATE TABLE c20_1 (id INT NOT NULL)
< ok 0 0
> query
| SELECT COLUMN_DEFAULT FROM information_schema.COLUMNS
| 		 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?
: string
| c20_1
: string
| id
< result
< column "COLUMN_DEFAULT" "TEXT"
< row
: null
> query
| SHOW CREATE TABLE c20_1
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c20_1
: bytes
| CREATE TABLE `c20_1` (
|   `id` int(11) NOT NULL
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
| CREATE TABLE c20_2 (id INT)
< ok 0 0
> query
| SHOW CREATE TABLE c20_2
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c20_2
: bytes
| CREATE TABLE `c20_2` (
|   `id` int(11) DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> query
| SELECT COLUMN_DEFAULT FROM information_schema.COLUMNS
| 		 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?
: string
| c20_2
: string
| id
< result
< column "COLUMN_DEFAULT" "TEXT"
< row
: bytes
| NULL
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
| CREATE TABLE c20_3 (name VARCHAR(64) NOT NULL, code CHAR(4) NOT NULL)
< ok 0 0
> query
| SELECT COLUMN_DEFAULT FROM information_schema.COLUMNS
| 		 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?
: string
| c20_3
: string
| name
< result
< column "COLUMN_DEFAULT" "TEXT"
< row
: null
> query ^
: string
| c20_3
: string
| code
< result
< column "COLUMN_DEFAULT" "TEXT"
< row
: null
> query
| SHOW CREATE TABLE c20_3
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c20_3
: bytes
| CREATE TABLE `c20_3` (
|   `name` varchar(64) NOT NULL,
|   `code` char(4) NOT NULL
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
| CREATE TABLE c20_4 (
| 			status ENUM('active','archived','deleted') NOT NULL,
| 			kind   ENUM('a','b','c')
| 		)
< ok 0 0
> query
| SELECT COLUMN_DEFAULT FROM information_schema.COLUMNS
| 		 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?
: string
| c20_4
: string
| status
< result
< column "COLUMN_DEFAULT" "TEXT"
< row
: null
> query ^
: string
| c20_4
: string
| kind
< result
< column "COLUMN_DEFAULT" "TEXT"
< row
: bytes
| NULL
> query
| SHOW CREATE TABLE c20_4
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c20_4
: bytes
| CREATE TABLE `c20_4` (
|   `status` enum('active','archived','deleted') NOT NULL,
|   `kind` enum('a','b','c') DEFAULT NULL
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
| CREATE TABLE c20_5 (
| 			created_at DATETIME NOT NULL,
| 			birthday   DATE     NOT NULL
| 		)
< ok 0 0
> query
| SELECT COLUMN_DEFAULT FROM information_schema.COLUMNS
| 		 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME=? AND COLUMN_NAME=?
: string
| c20_5
: string
| created_at
< result
< column "COLUMN_DEFAULT" "TEXT"
< row
: null
> query ^
: string
| c20_5
: string
| birthday
< result
< column "COLUMN_DEFAULT" "TEXT"
< row
: null
> query
| SHOW CREATE TABLE c20_5
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c20_5
: bytes
| CREATE TABLE `c20_5` (
|   `created_at` datetime NOT NULL,
|   `birthday` date NOT NULL
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
| CREATE TABLE c20_6a (b BLOB DEFAULT 'abc')
< ok 0 0
> query
| SHOW CREATE TABLE c20_6a
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c20_6a
: bytes
| CREATE TABLE `c20_6a` (
|   `b` blob DEFAULT 'abc'
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
| CREATE TABLE c20_6b (t TEXT DEFAULT 'hello')
< ok 0 0
> query
| SHOW CREATE TABLE c20_6b
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c20_6b
: bytes
| CREATE TABLE `c20_6b` (
|   `t` text DEFAULT 'hello'
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
| CREATE TABLE c20_6d (j JSON DEFAULT '[]')
< ok 0 0
> query
| SHOW CREATE TABLE c20_6d
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| c20_6d
: bytes
| CREATE TABLE `c20_6d` (
|   `j` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT '[]' CHECK (json_valid(`j`))
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
| CREATE TABLE c20_6c (g GEOMETRY DEFAULT 'x')
< error
| 1416 22003 Cannot get geometry object from data you send to the GEOMETRY field
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
| CREATE TABLE c20_8a (a INT, b INT AS (a + 1) DEFAULT 0)
< error
| 1064 42000 …near 'DEFAULT 0)' at line 1
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
| CREATE TABLE c20_8b (a INT, b INT GENERATED ALWAYS AS (a + 1) VIRTUAL DEFAULT 0)
< error
| 1064 42000 …near 'DEFAULT 0)' at line 1
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
| CREATE TABLE c20_8c (a INT, b INT GENERATED ALWAYS AS (a + 1) STORED DEFAULT (a * 2))
< error
| 1064 42000 …near 'DEFAULT (a * 2))' at line 1
