-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

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
< column "Create Table" "TEXT"
< row
: bytes
| c20_1
: bytes
| CREATE TABLE `c20_1` (
|   `id` int(11) NOT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| c20_2
: bytes
| CREATE TABLE `c20_2` (
|   `id` int(11) DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| c20_3
: bytes
| CREATE TABLE `c20_3` (
|   `name` varchar(64) NOT NULL,
|   `code` char(4) NOT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| c20_4
: bytes
| CREATE TABLE `c20_4` (
|   `status` enum('active','archived','deleted') NOT NULL,
|   `kind` enum('a','b','c') DEFAULT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
< column "Create Table" "TEXT"
< row
: bytes
| c20_5
: bytes
| CREATE TABLE `c20_5` (
|   `created_at` datetime NOT NULL,
|   `birthday` date NOT NULL
| ) ORGANIZATION INDEX DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMPRESSION = 'zstd_1.3.8' REPLICA_NUM = 1 BLOCK_SIZE = 16384 USE_BLOOM_FILTER = FALSE ENABLE_MACRO_BLOCK_BLOOM_FILTER = FALSE TABLET_SIZE = 134217728 PCTFREE = 0
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
| 1067 42000 Invalid default value for 'b'
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
| 1067 42000 Invalid default value for 't'
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
< error
| 1101 42000 BLOB, TEXT, GEOMETRY or JSON column 'tags' can't have a default value
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
