-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

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
|   `id` int NOT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
|   `id` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
: null
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
: null
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
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
< error
| 1101 42000 BLOB, TEXT, GEOMETRY or JSON column 'b' can't have a default value
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
| CREATE TABLE c20_6b (t TEXT DEFAULT 'hello')
< error
| 1101 42000 BLOB, TEXT, GEOMETRY or JSON column 't' can't have a default value
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
| CREATE TABLE c20_6c (g GEOMETRY DEFAULT 'x')
< error
| 1101 42000 BLOB, TEXT, GEOMETRY or JSON column 'g' can't have a default value
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
| CREATE TABLE c20_6d (j JSON DEFAULT '[]')
< error
| 1101 42000 BLOB, TEXT, GEOMETRY or JSON column 'j' can't have a default value
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
| CREATE TABLE c20_7 (
| 			id INT PRIMARY KEY,
| 			tags  JSON      DEFAULT (JSON_ARRAY()),
| 			meta  JSON      DEFAULT (JSON_OBJECT('v', 1)),
| 			blob1 BLOB      DEFAULT (SUBSTRING('abcdef', 1, 3)),
| 			pt    POINT     DEFAULT (POINT(0, 0)),
| 			uuid  BINARY(16) DEFAULT (UUID_TO_BIN(UUID()))
| 		)
< ok 0 0
> query
| SELECT COLUMN_DEFAULT, EXTRA FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='c20_7' AND COLUMN_NAME='tags'
< result
< column "COLUMN_DEFAULT" "TEXT"
< column "EXTRA" "VARCHAR"
< row
: bytes
| json_array()
: bytes
| DEFAULT_GENERATED
> query
| SELECT COLUMN_DEFAULT, EXTRA FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='c20_7' AND COLUMN_NAME='meta'
< result
< column "COLUMN_DEFAULT" "TEXT"
< column "EXTRA" "VARCHAR"
< row
: bytes
| json_object(_utf8mb4\'v\',1)
: bytes
| DEFAULT_GENERATED
> query
| SELECT COLUMN_DEFAULT, EXTRA FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='c20_7' AND COLUMN_NAME='blob1'
< result
< column "COLUMN_DEFAULT" "TEXT"
< column "EXTRA" "VARCHAR"
< row
: bytes
| substr(_utf8mb4\'abcdef\',1,3)
: bytes
| DEFAULT_GENERATED
> query
| SELECT COLUMN_DEFAULT, EXTRA FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='c20_7' AND COLUMN_NAME='pt'
< result
< column "COLUMN_DEFAULT" "TEXT"
< column "EXTRA" "VARCHAR"
< row
: bytes
| point(0,0)
: bytes
| DEFAULT_GENERATED
> query
| SELECT COLUMN_DEFAULT, EXTRA FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='c20_7' AND COLUMN_NAME='uuid'
< result
< column "COLUMN_DEFAULT" "TEXT"
< column "EXTRA" "VARCHAR"
< row
: bytes
| uuid_to_bin(uuid())
: bytes
| DEFAULT_GENERATED
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
| CREATE TABLE c20_8a (a INT, b INT AS (a + 1) DEFAULT 0)
< error
| 1221 HY000 Incorrect usage of DEFAULT and generated column
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
| 1221 HY000 Incorrect usage of DEFAULT and generated column
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
| 1221 HY000 Incorrect usage of DEFAULT and generated column
