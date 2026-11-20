-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestSimulationCausesOracle default#1
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER VIEW a AS SELECT b FROM a
< error
| 1347 HY000 'test.a' is not VIEW
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| ALTER VIEW a AS SELECT b FROM a
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE VIEW a AS SELECT 1 AS b
< ok 0 0
> exec
| ALTER VIEW a AS SELECT 1 AS b
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| DROP VIEW a
< error
| 1347 HY000 'test.a' is not VIEW
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| DROP VIEW a
< error
| 1051 42S02 Unknown table 'test.a'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(b)
< error
| 1283 HY000 Column 'b' cannot be part of FULLTEXT index
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT idx_d(c)
< error
| 1061 42000 Duplicate key name 'idx_d'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(x), ADD x INT
< error
| 1283 HY000 Column 'x' cannot be part of FULLTEXT index
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(x), ADD x TEXT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX ft ON a(b)
< error
| 1283 HY000 Column 'b' cannot be part of FULLTEXT index
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (b INT, FULLTEXT KEY ft (b))
< error
| 1283 HY000 Column 'b' cannot be part of FULLTEXT index
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE(a)
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(a) REFERENCES a
< error
| 1072 42000 Key column 'a' doesn't exist in table
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP CHECK nope
< error
| 3821 HY000 Check constraint 'nope' is not found in the table.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP CONSTRAINT nope
< error
| 3940 HY000 Constraint 'nope' does not exist.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, CONSTRAINT fk1 FOREIGN KEY (c) REFERENCES p(id))
< ok 0 0
> exec
| ALTER TABLE a DROP CHECK fk1
< error
| 3821 HY000 Check constraint 'fk1' is not found in the table.
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `fk1` (`c`),
|   CONSTRAINT `fk1` FOREIGN KEY (`c`) REFERENCES `p` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, CONSTRAINT ck1 CHECK (c > 0))
< ok 0 0
> exec
| ALTER TABLE a DROP FOREIGN KEY ck1
< error
| 1091 42000 Can't DROP 'ck1'; check that column/key exists
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   CONSTRAINT `ck1` CHECK ((`c` > 0))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, CONSTRAINT ck1 CHECK (c > 0))
< ok 0 0
> exec
| ALTER TABLE a DROP CONSTRAINT ck1
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER CHECK nope ENFORCED
< error
| 3821 HY000 Check constraint 'nope' is not found in the table.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER CONSTRAINT nope NOT ENFORCED
< error
| 3940 HY000 Constraint 'nope' does not exist.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY idx_d(missing)
< error
| 1061 42000 Duplicate key name 'idx_d'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PRIMARY KEY(missing)
< error
| 1068 42000 Multiple primary key defined
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(a) WITH PARSER a
< error
| 1128 HY000 Function 'a' is not defined
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b TEXT)
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT(b) WITH PARSER ngram
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX i ON a(missing)
< error
| 1072 42000 Key column 'missing' doesn't exist in table
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX idx_d ON a(missing)
< error
| 1061 42000 Duplicate key name 'idx_d'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX i ON a(missing) WITH PARSER a
< error
| 1128 HY000 Function 'a' is not defined
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (b INT, KEY(missing))
< error
| 1072 42000 Key column 'missing' doesn't exist in table
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(c2), ADD c2 INT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(c), DROP COLUMN c
< error
| 1072 42000 Key column 'c' doesn't exist in table
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(c), DROP COLUMN c, ADD COLUMN c INT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(y), CHANGE c y INT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME COLUMN c TO y, ADD KEY(c)
< error
| 1072 42000 Key column 'c' doesn't exist in table
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(x), ADD COLUMN x INT, RENAME COLUMN c TO y
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN x INT, CHANGE x y INT
< error
| 1054 42S22 Unknown column 'x' in 'a'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(x(11)), ADD x VARCHAR(10)
< error
| 1089 HY000 Incorrect prefix key; the used key part isn't a string, the used length is longer than the key part, or the storage engine doesn't support unique prefix keys
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(x), ADD x TEXT
< error
| 1170 42000 BLOB/TEXT column 'x' used in key specification without a key length
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (id INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE a ADD FOREIGN KEY (x2) REFERENCES s(id), ADD x2 VARCHAR(10)
< error
| 3780 HY000 Referencing column 'x2' and referenced column 'id' in foreign key constraint 'a_ibfk_1' are incompatible.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (id INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE a ADD FOREIGN KEY (x2) REFERENCES s(id) ON DELETE SET NULL, ADD x2 INT NOT NULL
< error
| 1830 HY000 Column 'x2' cannot be NOT NULL: needed in a foreign key constraint 'a_ibfk_1' SET NULL
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (id INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE a ADD FOREIGN KEY (x2) REFERENCES s(id), ADD x2 INT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL KEY(x), ADD x GEOMETRY
< error
| 1252 42000 All parts of a SPATIAL index must be NOT NULL
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PRIMARY KEY(x2), ADD x2 INT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(x), ADD COLUMN x INT, RENAME COLUMN x TO y
< error
| 1054 42S22 Unknown column 'x' in 'a'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME COLUMN x TO y
< error
| 1054 42S22 Unknown column 'x' in 'a'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN x INT, ADD INDEX idx_x (x), DROP COLUMN x
< error
| 1091 42000 Can't DROP 'x'; check that column/key exists
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (b INT, FULLTEXT(b) WITH PARSER a)
< error
| 1128 HY000 Function 'a' is not defined
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE VIEW a AS SELECT ABS(a)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE VIEW a AS SELECT 1 UNION SELECT a
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE VIEW a AS SELECT a HAVING 1
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE VIEW a AS SELECT 1 WHERE a
< error
| 1054 42S22 Unknown column 'a' in 'where clause'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE VIEW a AS SELECT (SELECT 1 AS x HAVING x)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT b FROM a
< error
| 1050 42S01 Table 'a' already exists
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SERIAL
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `a` bigint unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `a` (`a`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR BYTE
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `a` binary(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a CHAR COLLATE BINARY
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `a` binary(1) DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD a SERIAL NULL
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `a` bigint unsigned AUTO_INCREMENT,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `a` (`a`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT)
< ok 0 0
> exec
| ALTER TABLE a MODIFY c SERIAL
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   `c` bigint unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `c` (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT)
< ok 0 0
> exec
| ALTER TABLE a CHANGE c y SERIAL
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   `y` bigint unsigned NOT NULL AUTO_INCREMENT,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `y` (`y`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REMOVE PARTITIONING
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION (PARTITION p0)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION PARTITIONS 2
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DROP PARTITION p0
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TRUNCATE PARTITION ALL
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a COALESCE PARTITION 1
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REORGANIZE PARTITION p0 INTO (PARTITION p1)
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REORGANIZE PARTITION p0 INTO (PARTITION p1 VALUES LESS THAN (b))
< error
| 1054 42S22 Unknown column 'b' in 'partition function'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REORGANIZE PARTITION p0 INTO (PARTITION p1 ENGINE nosuch)
< error
| 1286 42000 Unknown storage engine 'nosuch'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REORGANIZE PARTITION p0 INTO (PARTITION p1 VALUES LESS THAN (10))
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT) PARTITION BY RANGE(b) SUBPARTITION BY HASH(b) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> exec
| ALTER TABLE a REORGANIZE PARTITION p0 INTO (PARTITION p1 VALUES LESS THAN (10) ENGINE InnoDB (SUBPARTITION s0 ENGINE nosuch))
< error
| 1286 42000 Unknown storage engine 'nosuch'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION (PARTITION p0 VALUES LESS THAN (b))
< error
| 1054 42S22 Unknown column 'b' in 'partition function'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| ALTER TABLE a REORGANIZE PARTITION p0 INTO (PARTITION p1 ENGINE nosuch)
< error
| 1286 42000 Unknown storage engine 'nosuch'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a REBUILD PARTITION ALL
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DISCARD PARTITION ALL TABLESPACE
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a IMPORT PARTITION ALL TABLESPACE
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (b INT)
< ok 0 0
> exec
| ALTER TABLE a EXCHANGE PARTITION p0 WITH TABLE s
< error
| 1505 HY000 Partition management on a not partitioned table is not possible
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a EXCHANGE PARTITION p0 WITH TABLE a
< error
| 1066 42000 Not unique table/alias: 'a'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ANALYZE PARTITION ALL
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CHECK PARTITION ALL
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT) PARTITION BY HASH(b) PARTITIONS 2
< ok 0 0
> exec
| ALTER TABLE a REMOVE PARTITIONING
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD e BIT(65)
< error
| 1439 42000 Display width out of range for column 'e' (max = 64)
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD e INT(256)
< error
| 1439 42000 Display width out of range for column 'e' (max = 255)
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD e INT(4294967296)
< error
| 1439 42000 Display width out of range for column 'e' (max = 4294967295)
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD e CHAR(4294967296)
< error
| 1439 42000 Display width out of range for column 'e' (max = 4294967295)
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD e DECIMAL(4294967296)
< error
| 1439 42000 Display width out of range for column 'e' (max = 4294967295)
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD b INT(256)
< error
| 1439 42000 Display width out of range for column 'b' (max = 255)
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| ALTER TABLE a ADD e INT(256)
< error
| 1439 42000 Display width out of range for column 'e' (max = 255)
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY b INT(256)
< error
| 1439 42000 Display width out of range for column 'b' (max = 255)
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE s (b INT(256))
< error
| 1439 42000 Display width out of range for column 'b' (max = 255)
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION (PARTITION p1 VALUES LESS THAN (b))
< error
| 1054 42S22 Unknown column 'b' in 'partition function'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION (PARTITION p1 ENGINE nosuch)
< error
| 1286 42000 Unknown storage engine 'nosuch'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PARTITION (PARTITION p1 VALUES LESS THAN (b) ENGINE nosuch)
< error
| 1054 42S22 Unknown column 'b' in 'partition function'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENGINE = nosuch
< error
| 1286 42000 Unknown storage engine 'nosuch'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| ALTER TABLE a ENGINE = nosuch
< error
| 1286 42000 Unknown storage engine 'nosuch'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (b INT) ENGINE = nosuch
< error
| 1286 42000 Unknown storage engine 'nosuch'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD e INT(256), ENGINE = nosuch
< error
| 1439 42000 Display width out of range for column 'e' (max = 255)
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENGINE = nosuch, ADD e INT(256)
< error
| 1286 42000 Unknown storage engine 'nosuch'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (b INT(256)) ENGINE = nosuch
< error
| 1439 42000 Display width out of range for column 'b' (max = 255)
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (b INT) ENGINE = PERFORMANCE_SCHEMA
< error
| 1683 HY000 Invalid performance_schema usage.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENGINE = PERFORMANCE_SCHEMA
< error
| 1031 HY000 Table storage engine for 'a' doesn't have this option
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| ALTER TABLE a ENGINE = PERFORMANCE_SCHEMA
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| CREATE TABLE s (b INT) ENGINE = nosuch
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| ALTER TABLE a ENGINE = nosuch
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| SET SESSION sql_mode = 'TRADITIONAL'
< ok 0 0
> exec
| CREATE TABLE s (b INT) ENGINE = nosuch
< error
| 1286 42000 Unknown storage engine 'nosuch'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| SET SESSION sql_mode = 'ANSI'
< ok 0 0
> exec
| CREATE TABLE s (b INT) ENGINE = nosuch
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| SET SESSION sql_mode = 'NO_ENGINE_SUBSTITUTION'
< ok 0 0
> exec
| ALTER TABLE a ENGINE = nosuch
< error
| 1286 42000 Unknown storage engine 'nosuch'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (b INT) ENGINE = MERGE
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENGINE = MERGE
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (b INT) ENGINE = HEAP
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (b INT) PARTITION BY RANGE(b) (PARTITION p0 VALUES LESS THAN (b))
< error
| 1054 42S22 Unknown column 'b' in 'partition function'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (b INT) PARTITION BY RANGE(b) SUBPARTITION BY HASH(b) (PARTITION p0 VALUES LESS THAN (10) ENGINE InnoDB (SUBPARTITION s0 ENGINE nosuch))
< error
| 1286 42000 Unknown storage engine 'nosuch'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (b INT) ENGINE = 'DEFAULT'
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENGINE = 'DEFAULT'
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD e FLOAT(256,2)
< error
| 1439 42000 Display width out of range for column 'e' (max = 255)
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD e DOUBLE(256,2)
< error
| 1439 42000 Display width out of range for column 'e' (max = 255)
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD e FLOAT(255,2)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD e INT(255)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD e BIT(64)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY(zzz)
< error
| 1488 HY000 Field in list of fields for partition function not found in table
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY HASH(zzz)
< error
| 1054 42S22 Unknown column 'zzz' in 'partition function'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY(c, c)
< error
| 1652 HY000 Duplicate partition field name 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY() SUBPARTITION BY HASH(zzz)
< error
| 1500 HY000 It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (b INT PRIMARY KEY, c INT) PARTITION BY HASH(b) PARTITIONS 2 SUBPARTITION BY HASH(c)
< error
| 1500 HY000 It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE(c)
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE(c)
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY HASH(c)
< error
| 1146 42S02 Table 'test.a' doesn't exist
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (b INT PRIMARY KEY, c INT) PARTITION BY KEY(zzz)
< error
| 1488 HY000 Field in list of fields for partition function not found in table
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY KEY()
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY LINEAR KEY() PARTITIONS 1
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY RANGE(b) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(c)
< error
| 1687 42000 A SPATIAL index may only contain a geometrical type column
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(b)
< error
| 1687 42000 A SPATIAL index may only contain a geometrical type column
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD g POINT
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(g)
< error
| 1252 42000 All parts of a SPATIAL index must be NOT NULL
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD g POINT NOT NULL
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL(g)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY USING RTREE(c)
< error
| 1687 42000 A SPATIAL index may only contain a geometrical type column
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD g POINT NOT NULL
< ok 0 0
> exec
| ALTER TABLE a ADD KEY USING RTREE(g)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `g` point NOT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   SPATIAL KEY `g` (`g`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE s (b INT NOT NULL, SPATIAL KEY (b))
< error
| 1687 42000 A SPATIAL index may only contain a geometrical type column
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL REFERENCE SYSTEM 5000
< error
| 3708 SR006 Missing mandatory attribute NAME.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL REFERENCE SYSTEM 5000 DEFINITION 's'
< error
| 3708 SR006 Missing mandatory attribute NAME.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL REFERENCE SYSTEM 5000 NAME 's'
< error
| 3708 SR006 Missing mandatory attribute DEFINITION.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE OR REPLACE SPATIAL REFERENCE SYSTEM 5000 ORGANIZATION 's' IDENTIFIED BY 1
< error
| 3708 SR006 Missing mandatory attribute NAME.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL REFERENCE SYSTEM 5000 NAME ''
< error
| 3710 SR006 The spatial reference system name can't be an empty string or start or end with whitespace.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL REFERENCE SYSTEM 5000 NAME ' s'
< error
| 3710 SR006 The spatial reference system name can't be an empty string or start or end with whitespace.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL REFERENCE SYSTEM 5000 NAME 's' DEFINITION ''
< error
| 3517 SR002 Can't parse the spatial reference system definition of SRID 5000.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE SPATIAL REFERENCE SYSTEM 0 NAME 's' DEFINITION 's'
< error
| 3714 SR000 SRID 0 is not modifiable.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW a AS SELECT *
< error
| 1096 HY000 No tables used
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT *
< error
| 1096 HY000 No tables used
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT t.*
< error
| 1051 42S02 Unknown table 't'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT t.* FROM a t
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT count(*)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a;
< error
| 1050 42S01 Table 'a' already exists
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z;
< error
| 4028 HY000 A table must have at least one visible column.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (KEY (b));
< error
| 1072 42000 Key column 'b' doesn't exist in table
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z PARTITION BY KEY();
< error
| 4028 HY000 A table must have at least one visible column.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT * FROM (SELECT *) d
< error
| 1096 HY000 No tables used
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS WITH c AS (SELECT *) SELECT * FROM c
< error
| 1096 HY000 No tables used
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS WITH c AS (SELECT *) SELECT 1
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT * FROM DUAL
< error
| 1096 HY000 No tables used
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE `dual` (x INT)
< ok 0 0
> exec
| CREATE VIEW v AS SELECT * FROM `dual`
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS WITH c1 AS (SELECT *), c2 AS (SELECT * FROM c1) SELECT 1
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS WITH c AS (SELECT *), d AS (WITH c AS (SELECT * FROM c) SELECT * FROM c) SELECT * FROM d
< error
| 1096 HY000 No tables used
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS WITH RECURSIVE c AS (SELECT 1 AS x) SELECT * FROM c
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE c2 (x INT)
< ok 0 0
> exec
| CREATE VIEW v AS WITH c1 AS (SELECT * FROM c2), c2 AS (SELECT *) SELECT * FROM c1
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS WITH c AS (SELECT *) SELECT * FROM (WITH c AS (SELECT 1 AS x) SELECT x FROM c) d
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT missing.* FROM a
< error
| 1051 42S02 Unknown table 'missing'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT other.a.* FROM a
< error
| 1051 42S02 Unknown table 'other.a'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT T.* FROM a AS t
< error
| 1051 42S02 Unknown table 'T'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT jt.* FROM JSON_TABLE('[1]', '$[*]' COLUMNS (x INT PATH '$')) AS jt
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT t.* FROM a AS t
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS WITH `dual` AS (SELECT *) SELECT 1 FROM DUAL
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE VIEW v AS SELECT 1 FROM DUAL
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE DATABASE cteother
< ok 1 0
> exec
| CREATE TABLE cteother.c (x INT)
< ok 0 0
> exec
| CREATE VIEW v AS WITH c AS (SELECT *) SELECT * FROM cteother.c
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 1 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE z (LIKE a);
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| ALTER TABLE nope ADD e INT DEFAULT NOW()
< error
| 1067 42000 Invalid default value for 'e'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD b INT DEFAULT CURRENT_TIMESTAMP
< error
| 1067 42000 Invalid default value for 'b'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a MODIFY b INT DEFAULT NOW()
< error
| 1067 42000 Invalid default value for 'b'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD e INT DEFAULT (NOW())
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD e DATETIME DEFAULT NOW()
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT DEFAULT NOW());
< error
| 1067 42000 Invalid default value for 'a'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a DATETIME DEFAULT NOW());
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT INVISIBLE);
< error
| 4028 HY000 A table must have at least one visible column.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = ON
< ok 0 0
> exec
| CREATE TABLE z (a INT INVISIBLE);
< error
| 4028 HY000 A table must have at least one visible column.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT INVISIBLE, b INT);
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT, b INT INVISIBLE)
< ok 0 0
> exec
| ALTER TABLE z ALTER COLUMN a SET INVISIBLE
< error
| 4028 HY000 A table must have at least one visible column.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a VARCHAR(10), INDEX ((LOWER(a))))
< ok 0 0
> exec
| ALTER TABLE z ALTER COLUMN a SET INVISIBLE
< error
| 4028 HY000 A table must have at least one visible column.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT, b INT INVISIBLE)
< ok 0 0
> exec
| ALTER TABLE z ALTER COLUMN a SET INVISIBLE, ADD COLUMN c INT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT, b INT)
< ok 0 0
> exec
| ALTER TABLE z ALTER COLUMN a SET INVISIBLE, ALTER COLUMN b SET INVISIBLE
< error
| 4028 HY000 A table must have at least one visible column.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT, b INT INVISIBLE)
< ok 0 0
> exec
| ALTER TABLE z MODIFY a INT INVISIBLE
< error
| 4028 HY000 A table must have at least one visible column.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT INVISIBLE) AS SELECT 1 AS b;
< ok 1 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a COLLATE utf8mb4
< error
| 1273 HY000 Unknown collation: 'utf8mb4'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD c2 INT COLLATE utf8mb4
< error
| 1273 HY000 Unknown collation: 'utf8mb4'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a COLLATE utf8mb4 zzz
< error
| 1273 HY000 Unknown collation: 'utf8mb4'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a COLLATE utf8mb4_bin zzz
< error
| 1064 42000 …near 'zzz' at line 1
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z SELECT nosuchcol
< error
| 1054 42S22 Unknown column 'nosuchcol' in 'field list'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z AS SELECT *
< error
| 1096 HY000 No tables used
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a SELECT nosuchcol
< error
| 1050 42S01 Table 'a' already exists
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z AS SELECT 1 AS b
< ok 1 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT) AS SELECT nosuchcol
< error
| 1054 42S22 Unknown column 'nosuchcol' in 'field list'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT) AS SELECT 1
< ok 1 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE AT nosuchcol DO COMMIT
< error
| 1054 42S22 Unknown column 'nosuchcol' in 'field list'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE DEFINER = root EVENT e ON SCHEDULE AT CURRENT_TIMESTAMP + INTERVAL 1 DAY DO COMMIT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT e ON SCHEDULE AT nosuchcol DO COMMIT
< error
| 1054 42S22 Unknown column 'nosuchcol' in 'field list'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER EVENT nosuchevent ON SCHEDULE AT nosuchcol DO COMMIT
< error
| 1054 42S22 Unknown column 'nosuchcol' in 'field list'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE EVENT e ON SCHEDULE AT CURRENT_TIMESTAMP + INTERVAL 1 DAY DO COMMIT
< ok 0 0
> exec
| ALTER EVENT e ON SCHEDULE AT CURRENT_TIMESTAMP + INTERVAL 2 DAY DO COMMIT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT FIXED
< error
| 1478 HY000 Table storage engine 'InnoDB' does not support the create option 'ROW_TYPE'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a COMPRESSION 's'
< error
| 1478 HY000 Table storage engine 'InnoDB' does not support the create option 'COMPRESSION'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a COMPRESSION=''
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a COMPRESSION='ZLIB'
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a COMPRESSION='none'
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT FIXED, ADD KEY(zzz)
< error
| 1072 42000 Key column 'zzz' doesn't exist in table
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENGINE=MyISAM, ROW_FORMAT FIXED
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=FIXED
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CHECK(zzz)
< error
| 3812 HY000 An expression of non-boolean type specified to a check constraint 'a_chk_1'.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD zzz INT CHECK(zzz)
< error
| 3812 HY000 An expression of non-boolean type specified to a check constraint 'a_chk_1'.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD CHECK(b > 0)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD zzz INT AS(zzz)
< error
| 3107 HY000 Generated column can refer only to generated columns defined prior to it.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD zzz INT AS(b + 1)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER INDEX zzz VISIBLE
< error
| 1176 42000 Key 'zzz' doesn't exist in table 'a'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a RENAME KEY zzz TO yyy
< error
| 1176 42000 Key 'zzz' doesn't exist in table 'a'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD zzz INT SRID 1
< error
| 1221 HY000 Incorrect usage of SRID and non-geometry column
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a WITH VALIDATION
< error
| 1221 HY000 Incorrect usage of ALTER and WITH VALIDATION
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a WITHOUT VALIDATION
< error
| 1221 HY000 Incorrect usage of ALTER and WITH VALIDATION
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a LIKE a
< error
| 1066 42000 Not unique table/alias: 'a'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a LIKE a
< error
| 1066 42000 Not unique table/alias: 'a'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(b, b)
< error
| 1060 42S21 Duplicate column name 'b'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(c, c) REFERENCES a(b)
< error
| 1060 42S21 Duplicate column name 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| ALTER TABLE a OPTIMIZE PARTITION ALL
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a START TRANSACTION
< error
| 3979 HY000 START TRANSACTION clause cannot be used with ALTER TABLE command.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENCRYPTION 's'
< error
| 3184 HY000 Invalid encryption option.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD g GEOMETRYCOLLECTION
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `g` geomcollection DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD g GEOMCOLLECTION
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `g` geomcollection DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY((1))
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `functional_index` ((1))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE((b + 1))
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   UNIQUE KEY `functional_index` (((`b` + 1))),
|   KEY `idx_d` (`d`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY((b))
< error
| 3762 HY000 Functional index on a column is not supported. Consider using a regular index instead.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(((b)))
< error
| 3762 HY000 Functional index on a column is not supported. Consider using a regular index instead.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY idx_d ((b))
< error
| 3762 HY000 Functional index on a column is not supported. Consider using a regular index instead.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY((nosuch))
< error
| 3762 HY000 Functional index on a column is not supported. Consider using a regular index instead.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX idx_d ON a ((b))
< error
| 3762 HY000 Functional index on a column is not supported. Consider using a regular index instead.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (b INT, KEY ((nosuch)))
< error
| 3762 HY000 Functional index on a column is not supported. Consider using a regular index instead.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY((nosuch + 1))
< error
| 1054 42S22 Unknown column 'nosuch' in 'functional index'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY((x2 + 1)), ADD x2 INT
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   `x2` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`),
|   KEY `functional_index` (((`x2` + 1)))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PRIMARY KEY((b + 1))
< error
| 3756 HY000 The primary key cannot be a functional index
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (b INT, PRIMARY KEY((b + 1)))
< error
| 3756 HY000 The primary key cannot be a functional index
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, s VARCHAR(10) CHARACTER SET latin1) DEFAULT CHARSET=latin1 COLLATE=latin1_bin
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARSET DEFAULT
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   `s` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, s VARCHAR(10) CHARACTER SET latin1) DEFAULT CHARSET=latin1
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARSET DEFAULT COLLATE utf8mb4_bin
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   `s` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARSET DEFAULT COLLATE latin1_bin
< error
| 1253 42000 COLLATION 'latin1_bin' is not valid for CHARACTER SET 'utf8mb4'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (e ENUM('a', 'A') CHARACTER SET utf8mb4 COLLATE utf8mb4_bin) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARSET DEFAULT
< error
| 1291 HY000 Column 'e' has duplicated value 'a' in ENUM
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a SECONDARY_LOAD
< error
| 3889 HY000 Secondary engine operation failed. No secondary engine defined.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALGORITHM DEFAULT, SECONDARY_LOAD
< error
| 3889 HY000 Secondary engine operation failed. No secondary engine defined.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a SECONDARY_UNLOAD
< error
| 3889 HY000 Secondary engine operation failed. No secondary engine defined.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY) SECONDARY_ENGINE = rapid
< ok 0 0
> exec
| ALTER TABLE a SECONDARY_LOAD
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   PRIMARY KEY (`b`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci SECONDARY_ENGINE=rapid
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY) SECONDARY_ENGINE = rapid
< ok 0 0
> exec
| ALTER TABLE a SECONDARY_UNLOAD
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY) SECONDARY_ENGINE = rapid
< ok 0 0
> exec
| ALTER TABLE a SECONDARY_ENGINE = NULL
< ok 0 0
> exec
| ALTER TABLE a SECONDARY_LOAD
< error
| 3889 HY000 Secondary engine operation failed. No secondary engine defined.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TABLESPACE nosuch_ts
< error
| 3510 HY000 Tablespace nosuch_ts doesn't exist.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(nosuch_col), TABLESPACE nosuch_ts
< error
| 3510 HY000 Tablespace nosuch_ts doesn't exist.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE a (b INT) TABLESPACE nosuch_ts
< error
| 3510 HY000 Tablespace nosuch_ts doesn't exist.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TABLESPACE mysql
< error
| 3723 HY000 The table 'a' may not be created in the reserved tablespace 'mysql'.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TABLESPACE innodb_custom
< error
| 3119 42000 InnoDB: Tablespace names starting with `innodb_` are reserved.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TABLESPACE innodb_temporary
< error
| 1478 HY000 InnoDB: Tablespace `innodb_temporary` can only contain TEMPORARY tables.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TABLESPACE innodb_file_per_table
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) /*!50100 TABLESPACE `innodb_file_per_table` */ ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENGINE=MyISAM, TABLESPACE nosuch_ts
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_oracle ADD DATAFILE 'ts_oracle.ibd' ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TABLESPACE ts_oracle
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) /*!50100 TABLESPACE `ts_oracle` */ ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLESPACE ts_oracle
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_oracle ADD DATAFILE 'ts_oracle.ibd' ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a TABLESPACE TS_ORACLE
< error
| 3510 HY000 Tablespace TS_ORACLE doesn't exist.
> exec
| DROP TABLESPACE ts_oracle
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AUTOEXTEND_SIZE 1
< error
| 4025 HY000 AUTOEXTEND_SIZE value should be between 4M and 4096M
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AUTOEXTEND_SIZE 8192M
< error
| 4025 HY000 AUTOEXTEND_SIZE value should be between 4M and 4096M
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AUTOEXTEND_SIZE 5M
< error
| 4023 HY000 AUTOEXTEND_SIZE should be a multiple of 4M
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AUTOEXTEND_SIZE 0
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a AUTOEXTEND_SIZE 4M
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(nosuch_col), AUTOEXTEND_SIZE 1
< error
| 1072 42000 Key column 'nosuch_col' doesn't exist in table
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENGINE=MyISAM, AUTOEXTEND_SIZE 1
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (b INT) AUTOEXTEND_SIZE 1
< error
| 4025 HY000 AUTOEXTEND_SIZE value should be between 4M and 4096M
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a IMPORT TABLESPACE
< error
| 1813 HY000 Tablespace 'test/a' exists.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DISCARD TABLESPACE
< ok 0 0
> exec
| ALTER TABLE a IMPORT TABLESPACE
< error
| 1812 HY000 Tablespace is missing for table `test`.`a`.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DISCARD TABLESPACE
< ok 0 0
> exec
| TRUNCATE TABLE a
< error
| 1814 HY000 Tablespace has been discarded for table 'test/a'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT) ENGINE=MyISAM
< ok 0 0
> exec
| ALTER TABLE a IMPORT TABLESPACE
< error
| 1031 HY000 Table storage engine for 'a' doesn't have this option
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_oracle ADD DATAFILE 'ts_oracle.ibd' ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE a (b INT) TABLESPACE ts_oracle
< ok 0 0
> exec
| ALTER TABLE a DISCARD TABLESPACE
< error
| 1148 42000 InnoDB: Cannot discard table `test/a` because it is in a general tablespace. It must be file-per-table.
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLESPACE ts_oracle
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(c) REFERENCES a
< error
| 1239 42000 Incorrect foreign key definition for 'foreign key without name': Key reference and table reference don't match
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(b, c) REFERENCES a(b)
< error
| 1239 42000 Incorrect foreign key definition for 'foreign key without name': Key reference and table reference don't match
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE TABLE z (b INT, FOREIGN KEY(b) REFERENCES a)
< error
| 1239 42000 Incorrect foreign key definition for 'foreign key without name': Key reference and table reference don't match
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| SET SESSION foreign_key_checks = 0
< ok 0 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(c) REFERENCES nosuch_tbl
< error
| 1239 42000 Incorrect foreign key definition for 'foreign key without name': Key reference and table reference don't match
> exec
| SET SESSION foreign_key_checks = ON
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE p (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, CONSTRAINT fk1 FOREIGN KEY (c) REFERENCES p(id))
< ok 0 0
> exec
| ALTER TABLE a ADD CONSTRAINT fk1 FOREIGN KEY(c) REFERENCES p
< error
| 1239 42000 Incorrect foreign key definition for 'fk1': Key reference and table reference don't match
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(c(1))
< error
| 1089 HY000 Incorrect prefix key; the used key part isn't a string, the used length is longer than the key part, or the storage engine doesn't support unique prefix keys
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD PRIMARY KEY(c(1))
< error
| 1089 HY000 Incorrect prefix key; the used key part isn't a string, the used length is longer than the key part, or the storage engine doesn't support unique prefix keys
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD FOREIGN KEY(c(1)) REFERENCES a(b)
< error
| 1089 HY000 Incorrect prefix key; the used key part isn't a string, the used length is longer than the key part, or the storage engine doesn't support unique prefix keys
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX i ON a(c(1))
< error
| 1089 HY000 Incorrect prefix key; the used key part isn't a string, the used length is longer than the key part, or the storage engine doesn't support unique prefix keys
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (c DATETIME, KEY(c(1)))
< error
| 1089 HY000 Incorrect prefix key; the used key part isn't a string, the used length is longer than the key part, or the storage engine doesn't support unique prefix keys
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (c ENUM('a','b'), KEY(c(1)))
< error
| 1089 HY000 Incorrect prefix key; the used key part isn't a string, the used length is longer than the key part, or the storage engine doesn't support unique prefix keys
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (g POINT NOT NULL, SPATIAL KEY(g(1)))
< error
| 1089 HY000 Incorrect prefix key; the used key part isn't a string, the used length is longer than the key part, or the storage engine doesn't support unique prefix keys
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (b BIT(8), KEY(b(4)))
< error
| 1089 HY000 Incorrect prefix key; the used key part isn't a string, the used length is longer than the key part, or the storage engine doesn't support unique prefix keys
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (b YEAR, KEY(b(2)))
< error
| 1089 HY000 Incorrect prefix key; the used key part isn't a string, the used length is longer than the key part, or the storage engine doesn't support unique prefix keys
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (s VARCHAR(8), KEY(s(4)))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (j JSON, KEY(j))
< error
| 3152 42000 JSON column 'j' supports indexing only via generated columns on a specified JSON path.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| ALTER TABLE nope ADD x INT ON UPDATE CURRENT_TIMESTAMP
< error
| 1294 HY000 Invalid ON UPDATE clause for 'x' column
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD b INT ON UPDATE CURRENT_TIMESTAMP
< error
| 1294 HY000 Invalid ON UPDATE clause for 'b' column
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| ALTER TABLE nope ADD x DATETIME(3) ON UPDATE CURRENT_TIMESTAMP
< error
| 1294 HY000 Invalid ON UPDATE clause for 'x' column
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| ALTER TABLE nope ADD x DATETIME(3) DEFAULT CURRENT_TIMESTAMP
< error
| 1067 42000 Invalid default value for 'x'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x INT AUTO_INCREMENT
< error
| 1075 42000 Incorrect table definition; there can be only one auto column and it must be defined as a key
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x INT AUTO_INCREMENT, ADD KEY(x)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x INT AUTO_INCREMENT, ADD KEY(c, x)
< error
| 1075 42000 Incorrect table definition; there can be only one auto column and it must be defined as a key
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (c INT NOT NULL, x INT AUTO_INCREMENT, PRIMARY KEY (c, x)) ENGINE=MyISAM
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT AUTO_INCREMENT, c INT, KEY k (b))
< ok 0 0
> exec
| DROP INDEX k ON a
< error
| 1075 42000 Incorrect table definition; there can be only one auto column and it must be defined as a key
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT AUTO_INCREMENT, c INT, KEY k (b))
< ok 0 0
> exec
| ALTER TABLE a DROP INDEX k
< error
| 1075 42000 Incorrect table definition; there can be only one auto column and it must be defined as a key
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x ENUM('s', 's')
< error
| 1291 HY000 Column 'x' has duplicated value 's' in ENUM
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x ENUM('s', 'S')
< error
| 1291 HY000 Column 'x' has duplicated value 's' in ENUM
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x ENUM('s ', 's') COLLATE utf8mb4_bin
< error
| 1291 HY000 Column 'x' has duplicated value 's' in ENUM
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x ENUM('s', 'S') COLLATE utf8mb4_bin
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x ENUM('s', 'S') NOT NULL COLLATE utf8mb4_bin
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x SET('a', 'A')
< error
| 1291 HY000 Column 'x' has duplicated value 'a' in SET
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| ALTER TABLE a ADD x ENUM('s', 's')
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x INT DEFAULT (x)
< error
| 3767 HY000 Default value expression of column 'x' cannot refer to a column defined after it if that column is a generated column or has an expression as default value.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER COLUMN c SET DEFAULT (c)
< error
| 3767 HY000 Default value expression of column 'c' cannot refer to a column defined after it if that column is a generated column or has an expression as default value.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (x INT DEFAULT (y), y INT DEFAULT (1))
< error
| 3767 HY000 Default value expression of column 'x' cannot refer to a column defined after it if that column is a generated column or has an expression as default value.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (x INT DEFAULT (y), y INT AS (1))
< error
| 3767 HY000 Default value expression of column 'x' cannot refer to a column defined after it if that column is a generated column or has an expression as default value.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (x INT DEFAULT (y), y INT)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (y INT DEFAULT (1), x INT DEFAULT (y))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (x INT DEFAULT (dt), dt DATETIME DEFAULT NOW())
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (x INT DEFAULT (dt), dt DATETIME DEFAULT (NOW()))
< error
| 3767 HY000 Default value expression of column 'x' cannot refer to a column defined after it if that column is a generated column or has an expression as default value.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (id INT AUTO_INCREMENT PRIMARY KEY, x INT DEFAULT (id))
< error
| 3768 HY000 Default value expression of column 'x' cannot refer to an auto-increment column.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x INT DEFAULT (nosuch + 1)
< error
| 1054 42S22 Unknown column 'nosuch' in 'default value expression'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x INT DEFAULT ((SELECT 1))
< error
| 3769 HY000 Default value expression of column 'x' contains a disallowed function.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x INT DEFAULT (@v)
< error
| 3772 HY000 Default value expression of column 'x' cannot refer user or system variables.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x INT DEFAULT 's'
< error
| 1067 42000 Invalid default value for 'x'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ALTER COLUMN c SET DEFAULT 's'
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x INT DEFAULT '123'
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x INT DEFAULT ('s')
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x TINYINT DEFAULT 128
< error
| 1067 42000 Invalid default value for 'x'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x TINYINT UNSIGNED DEFAULT -1
< error
| 1067 42000 Invalid default value for 'x'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x DECIMAL(4,2) DEFAULT 100.00
< error
| 1067 42000 Invalid default value for 'x'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x VARCHAR(3) DEFAULT 'abcd'
< error
| 1067 42000 Invalid default value for 'x'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x BIT(1) DEFAULT 's'
< error
| 1067 42000 Invalid default value for 'x'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x BIT(8) DEFAULT 's'
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x ENUM('a', 'b') DEFAULT 'c'
< error
| 1067 42000 Invalid default value for 'x'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x YEAR DEFAULT 's'
< error
| 1067 42000 Invalid default value for 'x'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x DATE DEFAULT 's'
< error
| 1067 42000 Invalid default value for 'x'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x DATE DEFAULT '0000-00-00'
< error
| 1067 42000 Invalid default value for 'x'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| ALTER TABLE a ADD x DATE DEFAULT '0000-00-00'
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_oracle ADD DATAFILE 'ts_oracle.ibd' ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE z (b INT) TABLESPACE ts_oracle PARTITION BY HASH (b) PARTITIONS 2
< error
| 1478 HY000 InnoDB : A partitioned table is not allowed in a shared tablespace.
> exec
| DROP TABLESPACE ts_oracle
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_oracle ADD DATAFILE 'ts_oracle.ibd' ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE z (b INT) PARTITION BY RANGE (b) (PARTITION p0 VALUES LESS THAN (10) TABLESPACE ts_oracle)
< error
| 1478 HY000 InnoDB : A partitioned table is not allowed in a shared tablespace.
> exec
| DROP TABLESPACE ts_oracle
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (b INT) PARTITION BY RANGE (b) (PARTITION p0 VALUES LESS THAN (10) TABLESPACE nosuch_ts)
< error
| 3510 HY000 Tablespace nosuch_ts doesn't exist.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (b INT) PARTITION BY RANGE (b) (PARTITION p0 VALUES LESS THAN (10) TABLESPACE innodb_file_per_table)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_oracle ADD DATAFILE 'ts_oracle.ibd' ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE a (b INT) TABLESPACE ts_oracle
< ok 0 0
> exec
| ALTER TABLE a PARTITION BY HASH (b) PARTITIONS 2
< error
| 1478 HY000 InnoDB : A partitioned table is not allowed in a shared tablespace.
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLESPACE ts_oracle
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_oracle ADD DATAFILE 'ts_oracle.ibd' ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE a (b INT) PARTITION BY RANGE (b) (PARTITION p0 VALUES LESS THAN (10))
< ok 0 0
> exec
| ALTER TABLE a TABLESPACE ts_oracle
< error
| 1478 HY000 InnoDB : A partitioned table is not allowed in a shared tablespace.
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLESPACE ts_oracle
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a DISCARD TABLESPACE
< ok 0 0
> exec ^
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x ENUM('a', 'á')
< error
| 1291 HY000 Column 'x' has duplicated value 'a' in ENUM
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x ENUM('a', 'á') COLLATE utf8mb4_0900_as_ci
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x ENUM('a', 'A') COLLATE utf8mb4_0900_as_ci
< error
| 1291 HY000 Column 'x' has duplicated value 'a' in ENUM
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x ENUM('ß', 'ss') COLLATE utf8mb4_0900_ai_ci
< error
| 1291 HY000 Column 'x' has duplicated value 'ß' in ENUM
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x ENUM('ß', 'ss') COLLATE utf8mb4_general_ci
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x ENUM('ß', 's') COLLATE utf8mb4_general_ci
< error
| 1291 HY000 Column 'x' has duplicated value 'ß' in ENUM
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| SET SESSION sql_mode = 'STRICT_TRANS_TABLES,ALLOW_INVALID_DATES'
< ok 0 0
> exec
| ALTER TABLE a ADD x DATE DEFAULT '2024-02-31'
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| SET SESSION sql_mode = 'STRICT_TRANS_TABLES,ALLOW_INVALID_DATES'
< ok 0 0
> exec
| ALTER TABLE a ADD x DATETIME DEFAULT '2024-02-31 12:00:00'
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| SET SESSION sql_mode = 'STRICT_TRANS_TABLES,ALLOW_INVALID_DATES'
< ok 0 0
> exec
| ALTER TABLE a ADD x TIMESTAMP DEFAULT '2024-02-31 12:00:00'
< error
| 1067 42000 Invalid default value for 'x'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| SET SESSION sql_mode = 'STRICT_TRANS_TABLES,ALLOW_INVALID_DATES'
< ok 0 0
> exec
| ALTER TABLE a ADD x DATE DEFAULT '2024-02-32'
< error
| 1067 42000 Invalid default value for 'x'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (id INT PRIMARY KEY, c INT, x INT DEFAULT (c + 1))
< ok 0 0
> exec
| ALTER TABLE a DROP COLUMN c
< error
| 3766 HY000 Column 'c' of table 'a' has a default value expression dependency and cannot be dropped or renamed.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (id INT PRIMARY KEY, c INT, x INT DEFAULT (c + 1))
< ok 0 0
> exec
| ALTER TABLE a RENAME COLUMN c TO c2
< error
| 3766 HY000 Column 'c' of table 'a' has a default value expression dependency and cannot be dropped or renamed.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (id INT PRIMARY KEY, c INT, x INT DEFAULT (c + 1))
< ok 0 0
> exec
| ALTER TABLE a CHANGE COLUMN c c2 INT
< error
| 3766 HY000 Column 'c' of table 'a' has a default value expression dependency and cannot be dropped or renamed.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (id INT PRIMARY KEY, c INT, x INT DEFAULT (c + 1))
< ok 0 0
> exec
| ALTER TABLE a DROP COLUMN c, ADD COLUMN c INT
< error
| 3766 HY000 Column 'c' of table 'a' has a default value expression dependency and cannot be dropped or renamed.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (id INT PRIMARY KEY, c INT, x INT DEFAULT (c + 1))
< ok 0 0
> exec
| ALTER TABLE a RENAME COLUMN c TO c2, ALTER COLUMN x SET DEFAULT (c2 + 1)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (id INT PRIMARY KEY, c INT, x INT DEFAULT (c + 1))
< ok 0 0
> exec
| ALTER TABLE a DROP COLUMN c, DROP COLUMN x
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE s (c SET('18446744073709551615','m1','m2','m3','m4','m5','m6','m7','m8','m9','m10','m11','m12','m13','m14','m15','m16','m17','m18','m19','m20','m21','m22','m23','m24','m25','m26','m27','m28','m29','m30','m31','m32','m33','m34','m35','m36','m37','m38','m39','m40','m41','m42','m43','m44','m45','m46','m47','m48','m49','m50','m51','m52','m53','m54','m55','m56','m57','m58','m59','m60','m61','m62','m63') DEFAULT 18446744073709551615)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE s (c SET('18446744073709551616','m1','m2','m3','m4','m5','m6','m7','m8','m9','m10','m11','m12','m13','m14','m15','m16','m17','m18','m19','m20','m21','m22','m23','m24','m25','m26','m27','m28','m29','m30','m31','m32','m33','m34','m35','m36','m37','m38','m39','m40','m41','m42','m43','m44','m45','m46','m47','m48','m49','m50','m51','m52','m53','m54','m55','m56','m57','m58','m59','m60','m61','m62','m63') DEFAULT 18446744073709551616)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE s (c SET('m0','m1','m2','m3','m4','m5','m6','m7','m8','m9','m10','m11','m12','m13','m14','m15','m16','m17','m18','m19','m20','m21','m22','m23','m24','m25','m26','m27','m28','m29','m30','m31','m32','m33','m34','m35','m36','m37','m38','m39','m40','m41','m42','m43','m44','m45','m46','m47','m48','m49','m50','m51','m52','m53','m54','m55','m56','m57','m58','m59','m60','m61','m62','m63') DEFAULT 18446744073709551615)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c SET('a','3') DEFAULT 3)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c SET('a','4') DEFAULT 4)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_myisam ADD DATAFILE 'ts_myisam.ibd' ENGINE=MyISAM
< error
| 1478 HY000 Table storage engine 'MyISAM' does not support the create option 'CREATE TABLESPACE'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_mem ADD DATAFILE 'ts_mem.ibd' ENGINE=MEMORY
< error
| 1478 HY000 Table storage engine 'MEMORY' does not support the create option 'CREATE TABLESPACE'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_bogus ADD DATAFILE 'ts_bogus.ibd' ENGINE=Bogus
< error
| 1286 42000 Unknown storage engine 'Bogus'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| CREATE TABLESPACE ts_sub ADD DATAFILE 'ts_sub.ibd' ENGINE=Bogus
< ok 0 0
> exec
| DROP TABLESPACE ts_sub
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_eng ADD DATAFILE 'ts_eng.ibd' ENGINE=InnoDB
< ok 0 0
> exec
| ALTER TABLESPACE ts_eng ENGINE=MyISAM
< error
| 3642 HY000 Engine 'MyISAM' does not match stored engine 'InnoDB' for tablespace 'ts_eng'
> exec
| DROP TABLESPACE ts_eng
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_eng ADD DATAFILE 'ts_eng.ibd' ENGINE=InnoDB
< ok 0 0
> exec
| DROP TABLESPACE ts_eng ENGINE=MyISAM
< error
| 3642 HY000 Engine 'MyISAM' does not match stored engine 'InnoDB' for tablespace 'ts_eng'
> exec
| DROP TABLESPACE ts_eng
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| DROP TABLESPACE ts_missing ENGINE=MyISAM
< error
| 3510 HY000 Tablespace ts_missing doesn't exist.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) DEFAULT 01)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` char(1) DEFAULT '1'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) DEFAULT -0)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` char(1) DEFAULT '0'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) DEFAULT 1e0)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` char(1) DEFAULT '1'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(3) DEFAULT 1.50)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(2) DEFAULT .5)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(3) DEFAULT .5)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` char(3) DEFAULT '0.5'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) DEFAULT 1e-2)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(2) DEFAULT 1e-2)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` char(2) DEFAULT '0'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(4) DEFAULT 1.23e10)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(5) DEFAULT 1.23e10)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` char(5) DEFAULT '1e10'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c BINARY(1) DEFAULT 01)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` binary(1) DEFAULT '1'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c BINARY(3) DEFAULT 1.50)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_m ADD DATAFILE 'ts_m.ibd'
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ENGINE=MyISAM, TABLESPACE ts_m
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `d` int DEFAULT NULL,
|   PRIMARY KEY (`b`),
|   KEY `idx_d` (`d`)
| ) /*!50100 TABLESPACE `ts_m` */ ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLESPACE ts_m
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_m ADD DATAFILE 'ts_m.ibd'
< ok 0 0
> exec
| CREATE TABLE a (b INT) ENGINE=InnoDB TABLESPACE ts_m
< ok 0 0
> exec
| ALTER TABLE a ENGINE=MyISAM
< ok 0 0
> exec
| DROP TABLESPACE ts_m
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_m ADD DATAFILE 'ts_m.ibd'
< ok 0 0
> exec
| CREATE TABLE a (b INT) ENGINE=MyISAM TABLESPACE ts_m
< ok 0 0
> exec
| DROP TABLESPACE ts_m
< ok 0 0
> exec
| ALTER TABLE a ENGINE=InnoDB
< error
| 3510 HY000 Tablespace ts_m doesn't exist.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c DATETIME DEFAULT '2024-01-01 00:00:00+14:00')
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c DATETIME DEFAULT '2024-01-01 00:00:00+14:01')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c DATETIME DEFAULT '2024-01-01 00:00:00+99:00')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c DATETIME DEFAULT '2024-01-01 12:00:00-13:59')
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c DATETIME DEFAULT '2024-01-01 12:00:00-14:00')
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c DATETIME DEFAULT '2024-01-01 12:00:00-14:01')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c DATETIME DEFAULT '2024-01-01 00:00:00-00:00')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c DATETIME DEFAULT '2024-01-01 00:00:00Z')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c TIMESTAMP DEFAULT '1970-01-01 02:00:00+05:00')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c TIMESTAMP DEFAULT '1970-01-01 00:00:01-05:00')
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(2) DEFAULT 'ab   ')
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` char(2) DEFAULT 'ab'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c VARCHAR(2) DEFAULT 'a   ')
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` varchar(2) DEFAULT 'a '
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c BINARY(2) DEFAULT 'ab   ')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) DEFAULT 0x41)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) DEFAULT 0x4142)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) DEFAULT 0xE4B8AD)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET utf8mb4 DEFAULT 0xFF)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET latin1 DEFAULT 0xFF)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) DEFAULT b'0100000101000010')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| CREATE TABLE a (c TINYINT DEFAULT 128)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| CREATE TABLE a (c DECIMAL(3,1) DEFAULT 'abc')
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x ENUM('ff', 'ﬀ') COLLATE utf8mb4_0900_ai_ci
< error
| 1291 HY000 Column 'x' has duplicated value 'ff' in ENUM
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x ENUM('ff', 'ﬀ') COLLATE utf8mb4_0900_as_ci
< error
| 1291 HY000 Column 'x' has duplicated value 'ff' in ENUM
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x ENUM('ff', 'ﬀ') COLLATE utf8mb4_0900_as_cs
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x ENUM('ff', 'ﬀ') COLLATE utf8mb4_general_ci
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x ENUM('dz', 'ʣ') COLLATE utf8mb4_0900_ai_ci
< error
| 1291 HY000 Column 'x' has duplicated value 'dz' in ENUM
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x ENUM('l·', 'ŀ') COLLATE utf8mb4_0900_ai_ci
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x ENUM('a', '𝐀') COLLATE utf8mb4_0900_ai_ci
< error
| 1291 HY000 Column 'x' has duplicated value 'a' in ENUM
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE UNDO TABLESPACE u_dup ADD DATAFILE 'u_dup.ibu'
< ok 0 0
> exec
| CREATE UNDO TABLESPACE u_dup ADD DATAFILE 'u_dup2.ibu'
< error
| 1813 HY000 Tablespace 'u_dup' exists.
> exec
| ALTER UNDO TABLESPACE u_dup SET INACTIVE
< ok 0 0
> exec
| SELECT SLEEP(1)
< ok 0 0
> exec
| DROP UNDO TABLESPACE u_dup
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE UNDO TABLESPACE u_ext ADD DATAFILE 'u_ext.ibd'
< error
| 3121 HY000 The ADD DATAFILE filepath must end with '.ibu'.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| ALTER UNDO TABLESPACE u_nosuch SET INACTIVE
< error
| 3510 HY000 Tablespace u_nosuch doesn't exist.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| DROP UNDO TABLESPACE u_nosuch
< error
| 3510 HY000 Tablespace u_nosuch doesn't exist.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE UNDO TABLESPACE u_act ADD DATAFILE 'u_act.ibu'
< ok 0 0
> exec
| DROP UNDO TABLESPACE u_act
< error
| 1529 HY000 Failed to drop UNDO TABLESPACE u_act
> exec
| ALTER UNDO TABLESPACE u_act SET INACTIVE
< ok 0 0
> exec
| SELECT SLEEP(1)
< ok 0 0
> exec
| DROP UNDO TABLESPACE u_act
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE UNDO TABLESPACE u_inact ADD DATAFILE 'u_inact.ibu'
< ok 0 0
> exec
| ALTER UNDO TABLESPACE u_inact SET INACTIVE
< ok 0 0
> exec
| SELECT SLEEP(1)
< ok 0 0
> exec
| DROP UNDO TABLESPACE u_inact
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE UNDO TABLESPACE u_kind ADD DATAFILE 'u_kind.ibu'
< ok 0 0
> exec
| DROP TABLESPACE u_kind
< error
| 3119 42000 Cannot DROP TABLESPACE `u_kind` because it is an undo tablespace.  Please use DROP UNDO TABLESPACE.
> exec
| ALTER UNDO TABLESPACE u_kind SET INACTIVE
< ok 0 0
> exec
| SELECT SLEEP(1)
< ok 0 0
> exec
| DROP UNDO TABLESPACE u_kind
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c FLOAT(4,2) DEFAULT 100)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c FLOAT(4,2) DEFAULT 99.994)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c FLOAT(4,2) DEFAULT 99.995)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c DOUBLE(4,2) DEFAULT 100)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE s (pt POINT NOT NULL)
< ok 0 0
> exec
| CREATE UNIQUE INDEX u USING RTREE ON s(pt)
< error
| 3728 HY000 Spatial indexes can't be primary or unique indexes.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE UNIQUE INDEX u USING RTREE ON a(b)
< error
| 1687 42000 A SPATIAL index may only contain a geometrical type column
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE s (pt POINT NOT NULL, UNIQUE KEY u (pt) USING RTREE)
< error
| 3728 HY000 Spatial indexes can't be primary or unique indexes.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE s (pt POINT NOT NULL)
< ok 0 0
> exec
| ALTER TABLE s ADD UNIQUE KEY u (pt) USING RTREE
< error
| 3728 HY000 Spatial indexes can't be primary or unique indexes.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET SESSION time_zone = '-13:59'
< ok 0 0
> exec
| CREATE TABLE a (c TIMESTAMP DEFAULT '2038-01-19 03:14:07')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET SESSION time_zone = '+14:00'
< ok 0 0
> exec
| CREATE TABLE a (c TIMESTAMP DEFAULT '1970-01-01 00:00:01')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET SESSION time_zone = '+14:00'
< ok 0 0
> exec
| CREATE TABLE a (c TIMESTAMP DEFAULT '1970-01-01 14:00:01')
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET SESSION time_zone = '-14:00'
< error
| 1298 HY000 Unknown or incorrect time zone: '-14:00'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x ENUM('ae', 'ä') COLLATE utf8mb4_de_pb_0900_ai_ci
< error
| 1291 HY000 Column 'x' has duplicated value 'ae' in ENUM
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x ENUM('ae', 'ä') COLLATE utf8mb4_de_pb_0900_ai_ci
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x ENUM('ae', 'ä') COLLATE utf8mb4_german2_ci
< error
| 1291 HY000 Column 'x' has duplicated value 'ae' in ENUM
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x ENUM('aa', 'å') COLLATE utf8mb4_da_0900_ai_ci
< error
| 1291 HY000 Column 'x' has duplicated value 'aa' in ENUM
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c INT DEFAULT 1.5)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` int DEFAULT '2'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c DECIMAL(5,2) DEFAULT 1.235)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` decimal(5,2) DEFAULT '1.24'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c FLOAT(4,2) DEFAULT 1.5)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` float(4,2) DEFAULT '1.50'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c VARBINARY(2) DEFAULT 007)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` varbinary(2) DEFAULT '7'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c VARBINARY(4) DEFAULT 1e-2)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` varbinary(4) DEFAULT '0.01'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c BINARY(3) DEFAULT -0.0)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` binary(3) DEFAULT '0.0'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT) PARTITION BY RANGE (b) (PARTITION p0 VALUES LESS THAN (10) TABLESPACE = innodb_file_per_table)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (`b`)
| (PARTITION p0 VALUES LESS THAN (10) TABLESPACE = `innodb_file_per_table` ENGINE = InnoDB) */
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT) PARTITION BY RANGE (b) SUBPARTITION BY HASH (b) (PARTITION p0 VALUES LESS THAN (10) TABLESPACE = innodb_file_per_table (SUBPARTITION sp0, SUBPARTITION sp1))
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
| /*!50100 PARTITION BY RANGE (`b`)
| SUBPARTITION BY HASH (`b`)
| (PARTITION p0 VALUES LESS THAN (10)
|  (SUBPARTITION sp0 TABLESPACE = `innodb_file_per_table` ENGINE = InnoDB,
|   SUBPARTITION sp1 TABLESPACE = `innodb_file_per_table` ENGINE = InnoDB)) */
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET @tz = '+14:00'
< ok 0 0
> exec
| SET time_zone = @tz
< ok 0 0
> exec
| CREATE TABLE a (c TIMESTAMP DEFAULT '1970-01-01 00:00:01')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET SESSION time_zone = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET @'tz_q' = '-13:59'
< ok 0 0
> exec
| SET time_zone = @tz_q
< ok 0 0
> exec
| CREATE TABLE a (c TIMESTAMP DEFAULT '2038-01-19 03:14:07')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET SESSION time_zone = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET @tz_cat = CONCAT('+', '14:00')
< ok 0 0
> exec
| SET time_zone = @tz_cat
< ok 0 0
> exec
| CREATE TABLE a (c TIMESTAMP DEFAULT '1970-01-01 14:00:01')
< ok 0 0
> exec
| SET SESSION time_zone = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET time_zone = @omni_unset_tz
< error
| 1231 42000 Variable 'time_zone' can't be set to the value of 'NULL'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET @tz_null = NULL
< ok 0 0
> exec
| SET time_zone = @tz_null
< error
| 1231 42000 Variable 'time_zone' can't be set to the value of 'NULL'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET @tz_num = 1
< ok 0 0
> exec
| SET time_zone = @tz_num
< error
| 1232 42000 Incorrect argument type to variable 'time_zone'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET @tz_bad = 'not_a_tz'
< ok 0 0
> exec
| SET time_zone = @tz_bad
< error
| 1298 HY000 Unknown or incorrect time zone: 'not_a_tz'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c DATE DEFAULT 10101)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` date DEFAULT '2001-01-01'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c DATE DEFAULT 101)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` date DEFAULT '2000-01-01'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c DATE DEFAULT 100)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c DATE DEFAULT '10101')
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c DATE DEFAULT 13101)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c DATE DEFAULT '13101')
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c DATE DEFAULT '101')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c DATE DEFAULT 1010101)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` date DEFAULT '0101-01-01'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c TIMESTAMP DEFAULT 1010101)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c TIMESTAMP DEFAULT 00010101)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` timestamp NULL DEFAULT '2001-01-01 00:00:00'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c DATETIME(3) DEFAULT 101010101.5)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` datetime(3) DEFAULT '2000-01-01 01:01:01.500'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c DATETIME DEFAULT 20200101.5)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c DATETIME DEFAULT '9999-12-31 23:59:59.6')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c TIMESTAMP DEFAULT '2038-01-19 03:14:07.6')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c TIMESTAMP(1) DEFAULT '2038-01-19 03:14:07.6')
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` timestamp(1) NULL DEFAULT '2038-01-19 03:14:07.6'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET utf16 DEFAULT 0xFFFD)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET utf16le DEFAULT 0xFDFF)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET utf16 DEFAULT 0xD800)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET utf16 DEFAULT 0xDC00)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET utf16 DEFAULT 0xD83DDE00)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET utf16 DEFAULT 0x41)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET utf16 DEFAULT 0x004142)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(2) CHARACTER SET utf16 DEFAULT 0x004142)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET utf32 DEFAULT 0x41)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET utf8mb3 DEFAULT '😀')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET ascii DEFAULT 'é')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET SESSION sql_mode = 'STRICT_TRANS_TABLES,TIME_TRUNCATE_FRACTIONAL'
< ok 0 0
> exec
| CREATE TABLE a (c DATETIME DEFAULT '9999-12-31 23:59:59.9')
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` datetime DEFAULT '9999-12-31 23:59:59'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET SESSION sql_mode = 'STRICT_TRANS_TABLES,TIME_TRUNCATE_FRACTIONAL'
< ok 0 0
> exec
| CREATE TABLE a (c TIMESTAMP DEFAULT '2038-01-19 03:14:07.9')
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` timestamp NULL DEFAULT '2038-01-19 03:14:07'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET SESSION sql_mode = 'STRICT_TRANS_TABLES,TIME_TRUNCATE_FRACTIONAL'
< ok 0 0
> exec
| CREATE TABLE a (c TIMESTAMP DEFAULT '1970-01-01 00:00:00.9')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c TIMESTAMP DEFAULT '1970-01-01 00:00:00.9')
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` timestamp NULL DEFAULT '1970-01-01 00:00:01'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET SESSION sql_mode = 'TIME_TRUNCATE_FRACTIONAL'
< ok 0 0
> exec
| CREATE TABLE a (c TIME DEFAULT '838:59:59.9')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET SESSION sql_mode = 'STRICT_TRANS_TABLES,TIME_TRUNCATE_FRACTIONAL'
< ok 0 0
> exec
| CREATE TABLE a (c TIME DEFAULT '838:59:58.9')
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` time DEFAULT '838:59:58'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c TIME DEFAULT '838:59:58.9')
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` time DEFAULT '838:59:59'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET SESSION sql_mode = 'STRICT_TRANS_TABLES,TIME_TRUNCATE_FRACTIONAL'
< ok 0 0
> exec
| CREATE TABLE a (c TIME(2) DEFAULT 10203.459)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` time(2) DEFAULT '01:02:03.45'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c TIME DEFAULT b'0011000100110010001110100011000000110000001110100011000000110000')
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` time DEFAULT '12:00:00'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c TIME DEFAULT 0x31323A30303A3030)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` time DEFAULT '12:00:00'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c DATE DEFAULT b'00110010001100000011001000110000001011010011000000110001001011010011000000110001')
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` date DEFAULT '2020-01-01'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c YEAR DEFAULT 0x07E4)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` year DEFAULT '2020'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c YEAR DEFAULT 0x32303230)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c ENUM('a','b') DEFAULT b'01100001')
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` enum('a','b') DEFAULT 'a'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c ENUM('a','b') DEFAULT 0x7A)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c SET('a','b') DEFAULT 0x612C62)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` set('a','b') DEFAULT 'a,b'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c SET('a','b') DEFAULT b'1')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET SESSION time_zone = '+08:00'
< ok 0 0
> exec
| CREATE TABLE a (dt DATETIME DEFAULT '2024-01-01 00:00:00+14:00', ts TIMESTAMP DEFAULT '2024-01-01 00:00:00+14:00', d DATE DEFAULT '2024-01-01 00:00:00+14:00')
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `dt` datetime DEFAULT '2023-12-31 18:00:00',
|   `ts` timestamp NULL DEFAULT '2023-12-31 18:00:00',
|   `d` date DEFAULT '2023-12-31'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET SESSION time_zone = '+00:00'
< ok 0 0
> exec
| CREATE TABLE a (dt3 DATETIME(3) DEFAULT '2024-01-01 00:00:00.5678+05:30', dth DATETIME DEFAULT 0x323032342D30312D30312030303A30303A30302B31343A3030)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `dt3` datetime(3) DEFAULT '2023-12-31 18:30:00.568',
|   `dth` datetime DEFAULT '2023-12-31 10:00:00'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET SESSION time_zone = '+05:00'
< ok 0 0
> exec
| CREATE TABLE a (dt DATETIME DEFAULT '9999-12-31 23:00:00+00:00')
< error
| 1067 42000 Invalid default value for 'dt'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET SESSION time_zone = '-05:00'
< ok 0 0
> exec
| CREATE TABLE a (dt DATETIME DEFAULT '0000-01-01 01:00:00+00:00')
< error
| 1067 42000 Invalid default value for 'dt'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (dt DATETIME, ts TIMESTAMP NULL)
< ok 0 0
> exec
| SET SESSION time_zone = '-05:00'
< ok 0 0
> exec
| ALTER TABLE a ALTER COLUMN dt SET DEFAULT '2024-01-01 02:00:00+00:00'
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `dt` datetime DEFAULT '2023-12-31 21:00:00',
|   `ts` timestamp NULL DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_bad_ext ADD DATAFILE 'ts_bad_ext.txt' ENGINE=InnoDB
< error
| 3121 HY000 The ADD DATAFILE filepath must end with '.ibd'.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_emp_df ADD DATAFILE '' ENGINE=InnoDB
< error
| 3121 HY000 The ADD DATAFILE filepath does not have a proper filename.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_dot_df ADD DATAFILE '.ibd' ENGINE=InnoDB
< error
| 3121 HY000 The ADD DATAFILE filepath does not have a proper filename.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_up_df ADD DATAFILE 'ts_up_df.IBD' ENGINE=InnoDB
< error
| 3121 HY000 The ADD DATAFILE filepath must end with '.ibd'.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_dir_df ADD DATAFILE '../ts_dir_df.ibd' ENGINE=InnoDB
< error
| 3121 HY000 The DATAFILE location must be in a known directory.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_df1 ADD DATAFILE 'shared_df.ibd' ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLESPACE ts_df2 ADD DATAFILE 'shared_df.ibd' ENGINE=InnoDB
< error
| 3606 HY000 Duplicate file name for tablespace 'ts_df2'
> exec
| DROP TABLESPACE ts_df1
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_fbs_bad1 ADD DATAFILE 'ts_fbs_bad1.ibd' FILE_BLOCK_SIZE=1 ENGINE=InnoDB
< error
| 1478 HY000 InnoDB does not support FILE_BLOCK_SIZE=1
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_fbs_32k ADD DATAFILE 'ts_fbs_32k.ibd' FILE_BLOCK_SIZE=32K ENGINE=InnoDB
< error
| 1478 HY000 InnoDB: Cannot create a tablespace with FILE_BLOCK_SIZE=32768 because INNODB_PAGE_SIZE=16384.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_fbs_8k ADD DATAFILE 'ts_fbs_8k.ibd' FILE_BLOCK_SIZE=8K ENGINE=InnoDB
< ok 0 0
> exec
| DROP TABLESPACE ts_fbs_8k
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_fbs_prec ADD DATAFILE 'ts_fbs_prec.txt' FILE_BLOCK_SIZE=1 AUTOEXTEND_SIZE=1 ENGINE=InnoDB
< error
| 1478 HY000 InnoDB does not support FILE_BLOCK_SIZE=1
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c INT) AUTOEXTEND_SIZE=4M
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` int DEFAULT NULL
| ) /*!80023 AUTOEXTEND_SIZE=4194304 */ ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c INT)
< ok 0 0
> exec
| ALTER TABLE a AUTOEXTEND_SIZE=8M
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` int DEFAULT NULL
| ) /*!80023 AUTOEXTEND_SIZE=8388608 */ ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c INT) AUTOEXTEND_SIZE=8M
< ok 0 0
> exec
| ALTER TABLE a AUTOEXTEND_SIZE=0
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` int DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE src_ae (c INT) AUTOEXTEND_SIZE=12M
< ok 0 0
> exec
| CREATE TABLE a LIKE src_ae
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` int DEFAULT NULL
| ) /*!80023 AUTOEXTEND_SIZE=12582912 */ ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_ae_gen ADD DATAFILE 'ts_ae_gen.ibd'
< ok 0 0
> exec
| CREATE TABLE a (c INT) TABLESPACE ts_ae_gen AUTOEXTEND_SIZE=8M
< error
| 4024 HY000 InnoDB: "AUTOEXTEND_SIZE" not allowed with general tablespaces
> exec
| DROP TABLESPACE ts_ae_gen
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TEMPORARY TABLE a_tmp_ae (c INT) AUTOEXTEND_SIZE=8M
< error
| 4026 HY000 AUTOEXTEND_SIZE clause is not valid for temporary tablespace.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c INT) TABLESPACE innodb_system AUTOEXTEND_SIZE=8M
< error
| 1478 HY000 InnoDB : AUTOEXTEND_SIZE is not accepted for system tablespace.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c DOUBLE DEFAULT 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c FLOAT DEFAULT 0x000000000000000001)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c DOUBLE DEFAULT 0x8000000000000000)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c DOUBLE UNSIGNED DEFAULT 0xFFFFFFFFFFFFFFFF)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` double unsigned DEFAULT '1.8446744073709552e19'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c FLOAT DEFAULT 0x7FFFFFFFFFFFFFFF)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` float DEFAULT '9.22337e18'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c DECIMAL(30,2) DEFAULT 0x8000000000000000)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c FLOAT DEFAULT 0x41424344)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` float DEFAULT '1094860000'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c BIGINT DEFAULT 0x000000000000000001)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| CREATE TABLE a (d DOUBLE DEFAULT 0xFFFFFFFFFFFFFFFF, b BIGINT DEFAULT 0x000000000000000001)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `d` double DEFAULT '9.223372036854776e18',
|   `b` bigint DEFAULT '9223372036854775807'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET sjis DEFAULT 0x82A0)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET sjis DEFAULT 0x82)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET ujis DEFAULT 0x8FA2AF)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET ujis DEFAULT 0x8FA2)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET gb18030 DEFAULT 0x81308130)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET gb18030 DEFAULT 0x813081)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET big5 DEFAULT 0xA440)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET euckr DEFAULT 0x815B)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET @tz = b'001010110011000100110100001110100011000000110000'
< ok 0 0
> exec
| SET time_zone = @tz
< ok 0 0
> exec
| CREATE TABLE a (c TIMESTAMP DEFAULT '1970-01-01 00:00:01')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET SESSION time_zone = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET @tz_b = b'001010110011000100110100001110100011000000110000'
< ok 0 0
> exec
| SET time_zone = @tz_b
< ok 0 0
> exec
| CREATE TABLE a (c TIMESTAMP DEFAULT '1970-01-01 14:00:01')
< ok 0 0
> exec
| SET SESSION time_zone = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (y0 YEAR DEFAULT 0, ys0 YEAR DEFAULT '0', ys000 YEAR DEFAULT '000', ys0000 YEAR DEFAULT '0000', y70 YEAR DEFAULT 70, ys01 YEAR DEFAULT '01')
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `y0` year DEFAULT '0000',
|   `ys0` year DEFAULT '2000',
|   `ys000` year DEFAULT '2000',
|   `ys0000` year DEFAULT '0000',
|   `y70` year DEFAULT '1970',
|   `ys01` year DEFAULT '2001'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (e1 ENUM('a','b') DEFAULT 'A', e2 ENUM('2','y') DEFAULT 2, e3 ENUM('ae','b') COLLATE utf8mb4_de_pb_0900_ai_ci DEFAULT 'ä', s1 SET('a','b') DEFAULT 'B,A,b')
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `e1` enum('a','b') DEFAULT 'a',
|   `e2` enum('2','y') DEFAULT 'y',
|   `e3` enum('ae','b') CHARACTER SET utf8mb4 COLLATE utf8mb4_de_pb_0900_ai_ci DEFAULT 'ae',
|   `s1` set('a','b') DEFAULT 'a,b'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET time_zone = '+5:00'
< ok 0 0
> exec
| CREATE TABLE a (c TIMESTAMP DEFAULT '1970-01-01 05:00:00')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET SESSION time_zone = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET time_zone = '-5:00'
< ok 0 0
> exec
| CREATE TABLE a (c DATETIME DEFAULT '2024-01-01 00:00:00+00:00')
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` datetime DEFAULT '2023-12-31 19:00:00'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET SESSION time_zone = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_datadir_abs ADD DATAFILE '/var/lib/mysql/ts_datadir_abs.ibd' ENGINE=InnoDB
< ok 0 0
> exec
| DROP TABLESPACE ts_datadir_abs
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_circ ADD DATAFILE '/var/lib/mysql/no_such_dir/../ts_circ.ibd' ENGINE=InnoDB
< error
| 3121 HY000 The ADD DATAFILE filepath cannot contain circular directory references.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET @mode_null = NULL
< ok 0 0
> exec
| SET sql_mode = @mode_null
< error
| 1231 42000 Variable 'sql_mode' can't be set to the value of 'NULL'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET sql_mode = @omni_unset_mode
< error
| 1231 42000 Variable 'sql_mode' can't be set to the value of 'NULL'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET @mode_bad = 'NOT_A_REAL_MODE'
< ok 0 0
> exec
| SET sql_mode = @mode_bad
< error
| 1231 42000 Variable 'sql_mode' can't be set to the value of 'NOT_A_REAL_MODE'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET @mode_flt = 1.5
< ok 0 0
> exec
| SET sql_mode = @mode_flt
< error
| 1232 42000 Incorrect argument type to variable 'sql_mode'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET ucs2 DEFAULT 0xD800)
< error
| 3507 HY000 Failed to update columns dictionary object.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET ucs2 DEFAULT 0xD83DDE00)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(2) CHARACTER SET ucs2 DEFAULT 0xD83DDE00)
< error
| 3507 HY000 Failed to update columns dictionary object.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c ENUM('a','b') DEFAULT 1)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c SET('a','b') DEFAULT 3)
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c ENUM('あ') CHARACTER SET sjis DEFAULT 0x82A0)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` enum('あ') CHARACTER SET sjis COLLATE sjis_japanese_ci DEFAULT 'あ'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c SET('あ','い') CHARACTER SET sjis DEFAULT 0x82A22C82A0)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` set('あ','い') CHARACTER SET sjis COLLATE sjis_japanese_ci DEFAULT 'あ,い'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (e1 ENUM('a','b') CHARACTER SET utf16 DEFAULT 0x41, e2 ENUM('é','b') CHARACTER SET latin1 DEFAULT 0xE9)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `e1` enum('a','b') CHARACTER SET utf16 COLLATE utf16_general_ci DEFAULT 'a',
|   `e2` enum('é','b') CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT 'é'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET @saved_sql_mode = @@sql_mode
< ok 0 0
> exec
| SET sql_mode = ''
< ok 0 0
> exec
| SET sql_mode = @saved_sql_mode
< ok 0 0
> exec
| SET sql_mode = @@sql_mode
< ok 0 0
> exec
| CREATE TABLE a (c DATE DEFAULT '0000-00-00')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET GLOBAL innodb_directories = '/tmp'
< error
| 1238 HY000 Variable 'innodb_directories' is a read only variable
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c VARCHAR(10), FULLTEXT KEY fk1 (c(100)))
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` varchar(10) DEFAULT NULL,
|   FULLTEXT KEY `fk1` (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c VARCHAR(10))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT INDEX fk2 (c(100))
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` varchar(10) DEFAULT NULL,
|   FULLTEXT KEY `fk2` (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c VARCHAR(10))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX fk3 ON a (c(100))
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` varchar(10) DEFAULT NULL,
|   FULLTEXT KEY `fk3` (`c`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET sjis DEFAULT '😀')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET sjis DEFAULT '😀')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET cp1251 DEFAULT 'é')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET cp1251 DEFAULT 'Д')
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` char(1) CHARACTER SET cp1251 COLLATE cp1251_general_ci DEFAULT 'Д'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET swe7 DEFAULT '[')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET swe7 DEFAULT 'Å')
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` char(1) CHARACTER SET swe7 COLLATE swe7_swedish_ci DEFAULT 'Å'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET latin1 DEFAULT '€')
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` char(1) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT '€'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(4) CHARACTER SET gb18030 DEFAULT '😀')
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` char(4) CHARACTER SET gb18030 COLLATE gb18030_chinese_ci DEFAULT 0x9439FC36
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE FUNCTION f() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| CREATE TABLE a (c INT DEFAULT (f()))
< error
| 3770 HY000 Default value expression of column 'c' contains a disallowed function: `f`.
> exec
| DROP FUNCTION IF EXISTS f
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c INT DEFAULT (unknown_fn()))
< error
| 3770 HY000 Default value expression of column 'c' contains a disallowed function: `unknown_fn`.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c VARCHAR(64) DEFAULT (VERSION()))
< error
| 3770 HY000 Default value expression of column 'c' contains a disallowed function: version().
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c INT)
< ok 0 0
> exec
| ALTER TABLE a ALTER COLUMN c SET DEFAULT (unknown_fn())
< error
| 3770 HY000 Default value expression of column 'c' contains a disallowed function: `unknown_fn`.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (x INT DEFAULT (ABS(-1)), y VARCHAR(36) DEFAULT (UUID()), z INT DEFAULT (IF(1,2,3)), g GEOMETRY DEFAULT (POINT(0,0)))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET utf8mb4 DEFAULT N'😀')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET utf8mb4 DEFAULT _utf8mb3'😀')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET utf8mb4 DEFAULT _ascii'é')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET utf8mb4 DEFAULT _sjis'😀')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET utf8mb4 DEFAULT _latin1'€')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(10) CHARACTER SET utf8mb4 DEFAULT _latin1'€')
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` char(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'â‚¬'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (id INT PRIMARY KEY) ENGINE=InnoDB
< ok 0 0
> exec
| ALTER TABLE a DISCARD TABLESPACE
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN c INT
< error
| 1814 HY000 Tablespace has been discarded for table 'test/a'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (id INT PRIMARY KEY) ENGINE=InnoDB
< ok 0 0
> exec
| ALTER TABLE a DISCARD TABLESPACE
< ok 0 0
> exec
| ALTER TABLE a COMMENT = 'ok'
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='ok'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE p_fk (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c_fk (id INT PRIMARY KEY, pid INT, CONSTRAINT fk_p FOREIGN KEY (pid) REFERENCES p_fk(id))
< ok 0 0
> exec
| ALTER TABLE p_fk DISCARD TABLESPACE
< error
| 1451 23000 Cannot delete or update a parent row: a foreign key constraint fails ()
> exec
| DROP TABLE IF EXISTS c_fk
< ok 0 0
> exec
| DROP TABLE IF EXISTS p_fk
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE p_fk (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c_fk (id INT PRIMARY KEY, pid INT, CONSTRAINT fk_p FOREIGN KEY (pid) REFERENCES p_fk(id))
< ok 0 0
> exec
| ALTER TABLE c_fk DISCARD TABLESPACE
< ok 0 0
> exec
| DROP TABLE IF EXISTS c_fk
< ok 0 0
> exec
| DROP TABLE IF EXISTS p_fk
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE p_fk (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE c_fk (id INT PRIMARY KEY, pid INT, CONSTRAINT fk_p FOREIGN KEY (pid) REFERENCES p_fk(id))
< ok 0 0
> exec
| SET foreign_key_checks = 0
< ok 0 0
> exec
| ALTER TABLE p_fk DISCARD TABLESPACE
< ok 0 0
> exec
| SET foreign_key_checks = 1
< ok 0 0
> exec
| DROP TABLE IF EXISTS c_fk
< ok 0 0
> exec
| DROP TABLE IF EXISTS p_fk
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE `a/b` ENGINE=InnoDB
< error
| 3119 42000 InnoDB: A general tablespace name cannot contain '/'.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_ren_slash ENGINE=InnoDB
< ok 0 0
> exec
| ALTER TABLESPACE ts_ren_slash RENAME TO `a/b`
< error
| 3119 42000 InnoDB: A general tablespace name cannot contain '/'.
> exec
| DROP TABLESPACE ts_ren_slash
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET time_zone = '+005:00'
< ok 0 0
> exec
| CREATE TABLE a (c TIMESTAMP DEFAULT '1970-01-01 05:00:00')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET SESSION time_zone = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET time_zone = '+5:0'
< ok 0 0
> exec
| CREATE TABLE a (c TIMESTAMP DEFAULT '1970-01-01 05:00:00')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET SESSION time_zone = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET @mode = 2 + 2
< ok 0 0
> exec
| SET sql_mode = @mode
< ok 0 0
> exec
| CREATE TABLE a (c VARCHAR(10) DEFAULT 'abc')
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE "a" (
|   "c" varchar(10) DEFAULT 'abc'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_emp_enc ADD DATAFILE 'ts_emp_enc.ibd' ENCRYPTION=''
< error
| 3184 HY000 Invalid encryption option.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET @m = 3 / 2
< ok 0 0
> exec
| SET sql_mode = @m
< error
| 1232 42000 Incorrect argument type to variable 'sql_mode'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET @m = 4 / 2
< ok 0 0
> exec
| SET sql_mode = @m
< error
| 1232 42000 Incorrect argument type to variable 'sql_mode'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_alt_emp ENGINE=InnoDB
< ok 0 0
> exec
| ALTER TABLESPACE ts_alt_emp ADD DATAFILE ''
< error
| 1533 HY000 Failed to alter: TABLESPACE ts_alt_emp
> exec
| DROP TABLESPACE ts_alt_emp
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_alt_drp ENGINE=InnoDB
< ok 0 0
> exec
| ALTER TABLESPACE ts_alt_drp DROP DATAFILE ''
< error
| 3629 HY000 Tablespace 'ts_alt_drp' does not have a file named ''
> exec
| DROP TABLESPACE ts_alt_drp
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_fbs8 ADD DATAFILE 'ts_fbs8.ibd' FILE_BLOCK_SIZE=8192 ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE a (id INT) TABLESPACE ts_fbs8 ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=4
< error
| 1478 HY000 InnoDB: Tablespace `ts_fbs8` uses block size 8192 and cannot contain a table with physical page size 4096
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLESPACE ts_fbs8
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_fbs8_ok ADD DATAFILE 'ts_fbs8_ok.ibd' FILE_BLOCK_SIZE=8192 ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE a (id INT) TABLESPACE ts_fbs8_ok ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLESPACE ts_fbs8_ok
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_fbs8_uncomp ADD DATAFILE 'ts_fbs8_uncomp.ibd' FILE_BLOCK_SIZE=8192 ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE a (id INT) TABLESPACE ts_fbs8_uncomp
< error
| 1478 HY000 InnoDB: Tablespace `ts_fbs8_uncomp` uses block size 8192 and cannot contain a table with physical page size 16384
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLESPACE ts_fbs8_uncomp
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_fbs16 ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE a (id INT) TABLESPACE ts_fbs16 ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8
< error
| 1478 HY000 InnoDB: Tablespace `ts_fbs16` cannot contain a COMPRESSED table
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLESPACE ts_fbs16
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_fbs8_alt ADD DATAFILE 'ts_fbs8_alt.ibd' FILE_BLOCK_SIZE=8192 ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE a (id INT) ENGINE=InnoDB
< ok 0 0
> exec
| ALTER TABLE a TABLESPACE ts_fbs8_alt ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=4
< error
| 1478 HY000 InnoDB: Tablespace `ts_fbs8_alt` uses block size 8192 and cannot contain a table with physical page size 4096
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLESPACE ts_fbs8_alt
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET NAMES ascii
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(10) CHARACTER SET utf8mb4 DEFAULT 'é')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec ^
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET NAMES ascii
< ok 0 0
> exec
| CREATE TABLE a (s ENUM('a') DEFAULT 'é')
< error
| 1067 42000 Invalid default value for 's'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec ^
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET NAMES sjis
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(10) CHARACTER SET utf8mb4 DEFAULT '😀')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec ^
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET NAMES latin1
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET utf8mb4 DEFAULT 'é')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec ^
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (id INT)
< ok 0 0
> exec
| SET NAMES latin1
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN d CHAR(1) CHARACTER SET utf8mb4 DEFAULT 'é'
< error
| 1067 42000 Invalid default value for 'd'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec ^
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET utf8mb4)
< ok 0 0
> exec
| SET NAMES latin1
< ok 0 0
> exec
| ALTER TABLE a ALTER COLUMN c SET DEFAULT 'é'
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec ^
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET NAMES binary
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(10) CHARACTER SET utf8mb4 DEFAULT 'é')
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec ^
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET NAMES binary
< ok 0 0
> exec
| CREATE TABLE a (s ENUM('a') DEFAULT 'é')
< error
| 1067 42000 Invalid default value for 's'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec ^
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET character_set_client = ascii
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(10) CHARACTER SET utf8mb4 DEFAULT 'é')
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` char(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '??'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec ^
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET NAMES ascii
< ok 0 0
> exec
| SET character_set_client = utf8mb4
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(10) CHARACTER SET utf8mb4 DEFAULT 'é')
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` char(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '?'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec ^
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET CHARACTER SET latin1
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET utf8mb4 DEFAULT 'é')
< error
| 1067 42000 Invalid default value for 'c'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec ^
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET NAMES nonexistent
< error
| 1115 42000 Unknown character set: 'nonexistent'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET character_set_client = nonexistent
< error
| 1115 42000 Unknown character set: 'nonexistent'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (id INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE a DISCARD TABLESPACE
< ok 0 0
> exec
| ALTER TABLE a DROP PRIMARY KEY
< error
| 1814 HY000 Tablespace has been discarded for table 'a'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET sql_mode=''
< ok 0 0
> exec
| CREATE TABLE a (a TINYINT DEFAULT 999)
< error
| 1067 42000 Invalid default value for 'a'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b BINARY(2) DEFAULT _latin1'é')
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` binary(2) DEFAULT 'é'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b VARBINARY(2) DEFAULT _latin1'é')
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` varbinary(2) DEFAULT 'é'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET NAMES latin1
< ok 0 0
> exec
| CREATE TABLE a (b BINARY(2) DEFAULT 'é')
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec ^
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b BINARY(2) DEFAULT _latin1'aé')
< error
| 1067 42000 Invalid default value for 'b'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c TINYBLOB, KEY(c(256)))
< error
| 1071 42000 Specified key was too long; max key length is 255 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c TINYTEXT CHARACTER SET utf8mb4, KEY(c(64)))
< error
| 1071 42000 Specified key was too long; max key length is 255 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c TINYTEXT CHARACTER SET utf8mb4, KEY(c(63)))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c BLOB, KEY(c(3073)))
< error
| 1071 42000 Specified key was too long; max key length is 3072 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c BLOB, KEY(c(1001))) ENGINE=MyISAM
< error
| 1071 42000 Specified key was too long; max key length is 1000 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET sql_mode=''
< ok 0 0
> exec
| CREATE TABLE a (c TINYBLOB, KEY(c(256)))
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` tinyblob,
|   KEY `c` (`c`(255))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_tmp_tbl ADD DATAFILE 'ts_tmp_tbl.ibd' ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TEMPORARY TABLE a (id INT) TABLESPACE ts_tmp_tbl
< error
| 1478 HY000 InnoDB: Tablespace `ts_tmp_tbl` cannot contain TEMPORARY tables.
> exec
| DROP TABLESPACE ts_tmp_tbl
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TEMPORARY TABLE a (id INT) TABLESPACE innodb_system
< error
| 1478 HY000 InnoDB: Tablespace `innodb_system` cannot contain TEMPORARY tables.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TEMPORARY TABLE a (id INT) TABLESPACE innodb_file_per_table
< error
| 1478 HY000 InnoDB: TABLESPACE=innodb_file_per_table option is disallowed for temporary tables with INNODB_STRICT_MODE=ON. This option is deprecated and will be removed in a future release
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TEMPORARY TABLE a (id INT) TABLESPACE innodb_temporary
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET @m = 4294967296 * 4294967296
< error
| 1690 22003 BIGINT value is out of range in '(4294967296 * 4294967296)'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET @m = 9223372036854775807 + 1
< error
| 1690 22003 BIGINT value is out of range in '(9223372036854775807 + 1)'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET @m = -9223372036854775808 - 1
< error
| 1690 22003 BIGINT value is out of range in '(-(9223372036854775808) - 1)'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET collation_connection = latin1_swedish_ci
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(10) CHARACTER SET utf8mb4 DEFAULT '😀')
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` char(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '?'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec ^
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET collation_connection = nonexistent_coll
< error
| 1273 HY000 Unknown collation: 'nonexistent_coll'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_ea_bad ADD DATAFILE 'ts_ea_bad.ibd' ENGINE_ATTRIBUTE='not-json'
< error
| 3980 HY000 Invalid json attribute, error: "Invalid value." at pos 1: 'ot-json'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_ea_ok ADD DATAFILE 'ts_ea_ok.ibd' ENGINE_ATTRIBUTE='{"a":1}'
< ok 0 0
> exec
| ALTER TABLESPACE ts_ea_ok ENGINE_ATTRIBUTE='not-json'
< error
| 3980 HY000 Invalid json attribute, error: "Invalid value." at pos 1: 'ot-json'
> exec
| DROP TABLESPACE ts_ea_ok
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE UNDO TABLESPACE u_rel_ok ADD DATAFILE 'u_rel_ok.ibu'
< ok 0 0
> exec
| ALTER UNDO TABLESPACE u_rel_ok SET INACTIVE
< ok 0 0
> exec
| SELECT SLEEP(1)
< ok 0 0
> exec
| DROP UNDO TABLESPACE u_rel_ok
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET collation_connection = utf8mb4_not_a_real_collation_ci
< error
| 1273 HY000 Unknown collation: 'utf8mb4_not_a_real_collation_ci'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET NAMES utf8mb4 COLLATE utf8mb4_not_a_real_collation_ci
< error
| 1273 HY000 Unknown collation: 'utf8mb4_not_a_real_collation_ci'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET character_set_client = utf8mb4
< ok 0 0
> exec
| SET collation_connection = latin1_swedish_ci
< ok 0 0
> exec
| CREATE TABLE a (b BINARY(1) DEFAULT 'é')
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `b` binary(1) DEFAULT 0xE9
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec ^
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET character_set_client = utf8mb4
< ok 0 0
> exec
| SET collation_connection = sjis_japanese_ci
< ok 0 0
> exec
| CREATE TABLE a (b BINARY(2) DEFAULT 'あ')
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec ^
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET character_set_client = utf8mb4
< ok 0 0
> exec
| SET collation_connection = sjis_japanese_ci
< ok 0 0
> exec
| CREATE TABLE a (b BINARY(1) DEFAULT 'あ')
< error
| 1067 42000 Invalid default value for 'b'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec ^
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET @p = @@SESSION.innodb_page_size
< error
| 1238 HY000 Variable 'innodb_page_size' is a GLOBAL variable
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET @p = @@LOCAL.innodb_page_size
< error
| 1238 HY000 Variable 'innodb_page_size' is a GLOBAL variable
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET @p = @@PERSIST.innodb_page_size
< error
| 1193 HY000 Unknown system variable 'PERSIST.innodb_page_size'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET @p = @@GLOBAL.innodb_page_size
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a VARCHAR(1) DEFAULT (SUBSTRING ('x',1,1)))
< error
| 3770 HY000 Default value expression of column 'a' contains a disallowed function: `SUBSTRING`.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET sql_mode = 'IGNORE_SPACE'
< ok 0 0
> exec
| CREATE TABLE a (a VARCHAR(1) DEFAULT (SUBSTRING ('x',1,1)))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET sql_mode = 'ANSI'
< ok 0 0
> exec
| CREATE TABLE a (a VARCHAR(1) DEFAULT (SUBSTRING ('x',1,1)))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET sql_mode = 'IGNORE_SPACE'
< ok 0 0
> exec
| CREATE TABLE a (a VARCHAR(1) DEFAULT (my_udf ('x')))
< error
| 3770 HY000 Default value expression of column 'a' contains a disallowed function: `my_udf`.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET character_set_client = ucs2
< error
| 1231 42000 Variable 'character_set_client' can't be set to the value of 'ucs2'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET character_set_client = utf16
< error
| 1231 42000 Variable 'character_set_client' can't be set to the value of 'utf16'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET character_set_client = utf16le
< error
| 1231 42000 Variable 'character_set_client' can't be set to the value of 'utf16le'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET character_set_client = utf32
< error
| 1231 42000 Variable 'character_set_client' can't be set to the value of 'utf32'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET NAMES utf16
< error
| 1231 42000 Variable 'character_set_client' can't be set to the value of 'utf16'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET CHARACTER SET utf32
< error
| 1231 42000 Variable 'character_set_client' can't be set to the value of 'utf32'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET NAMES utf8mb4 COLLATE latin1_swedish_ci
< error
| 1253 42000 COLLATION 'latin1_swedish_ci' is not valid for CHARACTER SET 'utf8mb4'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET character_set_client = NULL
< error
| 1231 42000 Variable 'character_set_client' can't be set to the value of 'NULL'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET GLOBAL character_set_client = latin1
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET utf8mb4 DEFAULT 'é')
< ok 0 0
> exec
| SET GLOBAL character_set_client = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec ^
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET PERSIST character_set_client = latin1
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(1) CHARACTER SET utf8mb4 DEFAULT 'é')
< ok 0 0
> exec
| SET GLOBAL character_set_client = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec ^
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET GLOBAL collation_connection = latin1_swedish_ci
< ok 0 0
> exec
| CREATE TABLE a (b BINARY(2) DEFAULT 'é')
< ok 0 0
> exec
| SET GLOBAL collation_connection = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec ^
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET GLOBAL character_set_client = ascii
< ok 0 0
> exec
| SET character_set_client = DEFAULT
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(10) CHARACTER SET utf8mb4 DEFAULT 'é')
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` char(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '??'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET GLOBAL character_set_client = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec ^
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET system_time_zone = 'UTC'
< error
| 1238 HY000 Variable 'system_time_zone' is a read only variable
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET PERSIST_ONLY system_time_zone = 'UTC'
< error
| 1238 HY000 Variable 'system_time_zone' is a non persistent read only variable
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET foreign_key_checks = 'bad'
< error
| 1231 42000 Variable 'foreign_key_checks' can't be set to the value of 'bad'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET GLOBAL sql_mode='', foreign_key_checks=0
< ok 0 0
> exec
| CREATE TABLE a (id INT, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES missing_parent(id))
< error
| 1824 HY000 Failed to open the referenced table 'missing_parent'
> exec
| SET GLOBAL sql_mode = DEFAULT
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET GLOBAL sql_mode='', @u=1, foreign_key_checks=0
< ok 0 0
> exec
| CREATE TABLE a (id INT, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES missing_parent(id))
< error
| 1824 HY000 Failed to open the referenced table 'missing_parent'
> exec
| SET GLOBAL sql_mode = DEFAULT
< ok 0 0
> exec
| SET GLOBAL foreign_key_checks = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET @@GLOBAL.sql_mode='', foreign_key_checks=0
< ok 0 0
> exec
| CREATE TABLE a (id INT, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES missing_parent(id))
< ok 0 0
> exec
| SET GLOBAL sql_mode = DEFAULT
< ok 0 0
> exec
| SET SESSION foreign_key_checks = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_enc_def ADD DATAFILE 'ts_enc_def.ibd'
< ok 0 0
> exec
| CREATE TABLE a (id INT) TABLESPACE ts_enc_def ENCRYPTION='Y'
< error
| 3825 HY000 Request to create 'encrypted' table while using an 'unencrypted' tablespace.
> exec
| DROP TABLESPACE ts_enc_def
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_enc_n ADD DATAFILE 'ts_enc_n.ibd' ENCRYPTION='N'
< ok 0 0
> exec
| CREATE TABLE a (id INT)
< ok 0 0
> exec
| ALTER TABLE a ENCRYPTION='Y', TABLESPACE ts_enc_n
< error
| 3825 HY000 Request to create 'encrypted' table while using an 'unencrypted' tablespace.
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLESPACE ts_enc_n
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET time_zone = 'Invalid/Zone'
< error
| 1298 HY000 Unknown or incorrect time zone: 'Invalid/Zone'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_tmp_like ADD DATAFILE 'ts_tmp_like.ibd'
< ok 0 0
> exec
| CREATE TABLE s (a INT) TABLESPACE ts_tmp_like
< ok 0 0
> exec
| CREATE TEMPORARY TABLE a LIKE s
< error
| 1478 HY000 InnoDB: Tablespace `ts_tmp_like` cannot contain TEMPORARY tables.
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLESPACE ts_tmp_like
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE s (a INT) TABLESPACE innodb_system
< ok 0 0
> exec
| CREATE TEMPORARY TABLE a LIKE s
< error
| 1478 HY000 InnoDB: Tablespace `innodb_system` cannot contain TEMPORARY tables.
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE s (a INT) TABLESPACE innodb_file_per_table
< ok 0 0
> exec
| CREATE TEMPORARY TABLE a LIKE s
< error
| 1478 HY000 InnoDB: TABLESPACE=innodb_file_per_table option is disallowed for temporary tables with INNODB_STRICT_MODE=ON. This option is deprecated and will be removed in a future release
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_myisam_like ADD DATAFILE 'ts_myisam_like.ibd'
< ok 0 0
> exec
| CREATE TABLE s (a INT) ENGINE=MyISAM TABLESPACE ts_myisam_like
< ok 0 0
> exec
| CREATE TEMPORARY TABLE a LIKE s
< ok 0 0
> exec
| DROP TEMPORARY TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLESPACE ts_myisam_like
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (id INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE a DISCARD TABLESPACE
< ok 0 0
> exec ^
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (id INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE a DISCARD TABLESPACE
< ok 0 0
> exec
| SET sql_mode = 'TRADITIONAL'
< ok 0 0
> exec
| ALTER TABLE a DISCARD TABLESPACE
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET @mode = ''
< ok 0 0
> exec
| SET @mode = 1.5 + 1
< ok 0 0
> exec
| SET sql_mode = @mode
< error
| 1232 42000 Incorrect argument type to variable 'sql_mode'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a TEXT CHARACTER SET ascii, b TEXT CHARACTER SET ascii, KEY(a(2000), b(2000)))
< error
| 1071 42000 Specified key was too long; max key length is 3072 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET sql_mode = ''
< ok 0 0
> exec
| CREATE TABLE a (a TEXT CHARACTER SET ascii, b TEXT CHARACTER SET ascii, KEY(a(2000), b(2000)))
< error
| 1071 42000 Specified key was too long; max key length is 3072 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a VARCHAR(500) CHARACTER SET utf8mb4, b VARCHAR(500) CHARACTER SET utf8mb4, KEY(a, b))
< error
| 1071 42000 Specified key was too long; max key length is 3072 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a TEXT CHARACTER SET ascii, b TEXT CHARACTER SET ascii, KEY(a(2000), b(1072)))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a TEXT CHARACTER SET ascii, b TEXT CHARACTER SET ascii)
< ok 0 0
> exec
| ALTER TABLE a ADD KEY(a(2000), b(2000))
< error
| 1071 42000 Specified key was too long; max key length is 3072 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a TEXT CHARACTER SET ascii, b TEXT CHARACTER SET ascii, KEY(a(600), b(500))) ENGINE=MyISAM
< error
| 1071 42000 Specified key was too long; max key length is 1000 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET SESSION sql_mode = REPLACE(@@sql_mode, 'ONLY_FULL_GROUP_BY', '')
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET PERSIST_ONLY innodb_page_size = 32768
< error
| 1238 HY000 Variable 'innodb_page_size' is a non persistent read only variable
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET PERSIST_ONLY innodb_page_size = 12345
< error
| 1238 HY000 Variable 'innodb_page_size' is a non persistent read only variable
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET PERSIST_ONLY innodb_directories = '/var/lib/mysql'
< error
| 1238 HY000 Variable 'innodb_directories' is a non persistent read only variable
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a VARCHAR(300) CHARACTER SET utf8mb4, KEY(a(200))) ROW_FORMAT=COMPACT
< error
| 1071 42000 Specified key was too long; max key length is 767 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a VARCHAR(300) CHARACTER SET utf8mb4, KEY(a(200))) ROW_FORMAT=REDUNDANT
< error
| 1071 42000 Specified key was too long; max key length is 767 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a VARCHAR(300) CHARACTER SET utf8mb4, KEY(a(191))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a VARCHAR(300) CHARACTER SET utf8mb4, KEY(a(200))) ROW_FORMAT=DYNAMIC
< ok 0 0
> exec
| ALTER TABLE a ROW_FORMAT=COMPACT
< error
| 1071 42000 Specified key was too long; max key length is 767 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a VARCHAR(300) CHARACTER SET utf8mb4, KEY(a(191))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| ALTER TABLE a ADD KEY k2 (a(200))
< error
| 1071 42000 Specified key was too long; max key length is 767 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (id INT) TABLESPACE innodb_system ROW_FORMAT=COMPRESSED
< error
| 1478 HY000 InnoDB: Tablespace `innodb_system` cannot contain a COMPRESSED table
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET sql_mode = missing_fn_db.CONCAT('ANSI')
< error
| 1305 42000 FUNCTION missing_fn_db.CONCAT does not exist
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET @u_qfn = missing_fn_db.CONCAT('ANSI')
< error
| 1305 42000 FUNCTION missing_fn_db.CONCAT does not exist
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET sql_mode = ''
< ok 0 0
> exec
| CREATE TABLE a (a TEXT CHARACTER SET utf8mb4, UNIQUE KEY(a(1000)))
< error
| 1071 42000 Specified key was too long; max key length is 3072 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET sql_mode = ''
< ok 0 0
> exec
| CREATE TABLE a (a TEXT CHARACTER SET utf8mb4, KEY(a(1000)))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a TEXT CHARACTER SET utf8mb4)
< ok 0 0
> exec
| SET sql_mode = ''
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE KEY (a(1000))
< error
| 1071 42000 Specified key was too long; max key length is 3072 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a TEXT CHARACTER SET utf8mb4)
< ok 0 0
> exec
| SET sql_mode = ''
< ok 0 0
> exec
| CREATE UNIQUE INDEX idx_uq ON a (a(1000))
< error
| 1071 42000 Specified key was too long; max key length is 3072 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a TEXT CHARACTER SET ascii, n BIGINT, KEY(a(3072), n))
< error
| 1071 42000 Specified key was too long; max key length is 3072 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a TEXT CHARACTER SET ascii, n BIGINT, KEY(a(3064), n))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a TEXT CHARACTER SET ascii, d DATETIME(6), KEY(a(3065), d))
< error
| 1071 42000 Specified key was too long; max key length is 3072 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a TEXT CHARACTER SET ascii, d DECIMAL(10,2), KEY(a(3068), d))
< error
| 1071 42000 Specified key was too long; max key length is 3072 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a TEXT CHARACTER SET ascii, b BIT(16), KEY(a(3071), b))
< error
| 1071 42000 Specified key was too long; max key length is 3072 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a TEXT CHARACTER SET ascii, e ENUM('x','y'), KEY(a(3072), e))
< error
| 1071 42000 Specified key was too long; max key length is 3072 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET sql_mode = IFNULL('', 'ANSI', 'TRADITIONAL')
< error
| 1582 42000 Incorrect parameter count in the call to native function 'IFNULL'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET sql_mode = CONCAT()
< error
| 1582 42000 Incorrect parameter count in the call to native function 'CONCAT'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a VARCHAR(10), KEY idx ((LOWER(a))))
< ok 0 0
> exec
| ALTER TABLE a RENAME COLUMN a TO b
< error
| 3837 HY000 Column 'a' has a functional index dependency and cannot be dropped or renamed.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a VARCHAR(10), KEY idx ((LOWER(a))))
< ok 0 0
> exec
| ALTER TABLE a DROP COLUMN a
< error
| 3837 HY000 Column 'a' has a functional index dependency and cannot be dropped or renamed.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a VARCHAR(10), KEY idx ((LOWER(a))))
< ok 0 0
> exec
| ALTER TABLE a CHANGE COLUMN a b INT
< error
| 3837 HY000 Column 'a' has a functional index dependency and cannot be dropped or renamed.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a VARCHAR(10), KEY idx ((LOWER(a))))
< ok 0 0
> exec
| ALTER TABLE a RENAME COLUMN a TO b, DROP INDEX idx
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a VARCHAR(10), KEY idx ((LOWER(a))))
< ok 0 0
> exec
| ALTER TABLE a MODIFY COLUMN a VARCHAR(20)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET time_zone = '+014:00'
< ok 0 0
> exec
| SET time_zone = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET time_zone = '+1:000'
< ok 0 0
> exec
| SET time_zone = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET time_zone = '++01:00'
< error
| 1298 HY000 Unknown or incorrect time zone: '++01:00'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a INT DEFAULT (IFNULL(1)))
< error
| 1582 42000 Incorrect parameter count in the call to native function 'IFNULL'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a INT DEFAULT (IFNULL(1, 2, 3)))
< error
| 1582 42000 Incorrect parameter count in the call to native function 'IFNULL'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a VARCHAR(10) DEFAULT (CONCAT()))
< error
| 1582 42000 Incorrect parameter count in the call to native function 'CONCAT'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a INT DEFAULT (ABS()))
< error
| 1582 42000 Incorrect parameter count in the call to native function 'ABS'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a INT DEFAULT (ABS(1, 2)))
< error
| 1582 42000 Incorrect parameter count in the call to native function 'ABS'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (id INT)
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN b INT DEFAULT (IFNULL(1))
< error
| 1582 42000 Incorrect parameter count in the call to native function 'IFNULL'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (id INT)
< ok 0 0
> exec
| ALTER TABLE a ALTER COLUMN id SET DEFAULT (IFNULL(1))
< error
| 1582 42000 Incorrect parameter count in the call to native function 'IFNULL'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a INT DEFAULT (IFNULL(1, 2)))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (j JSON DEFAULT (JSON_OBJECT('k')))
< error
| 1582 42000 Incorrect parameter count in the call to native function 'JSON_OBJECT'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (j JSON DEFAULT (JSON_OBJECT('k', 1, 'k2')))
< error
| 1582 42000 Incorrect parameter count in the call to native function 'JSON_OBJECT'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (j JSON DEFAULT (JSON_OBJECT()))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (j JSON DEFAULT (JSON_OBJECT('k', 1)))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (j JSON DEFAULT (JSON_SET('{}', '$.a')))
< error
| 1582 42000 Incorrect parameter count in the call to native function 'JSON_SET'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (j JSON DEFAULT (JSON_SET('{}', '$.a', 1)))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (j JSON DEFAULT (JSON_INSERT('{}', '$.a')))
< error
| 1582 42000 Incorrect parameter count in the call to native function 'JSON_INSERT'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (j JSON DEFAULT (JSON_REPLACE('{}', '$.a')))
< error
| 1582 42000 Incorrect parameter count in the call to native function 'JSON_REPLACE'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (j JSON DEFAULT (JSON_ARRAY_APPEND('[]', '$')))
< error
| 1582 42000 Incorrect parameter count in the call to native function 'JSON_ARRAY_APPEND'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (j JSON DEFAULT (JSON_ARRAY_INSERT('[]', '$[0]')))
< error
| 1582 42000 Incorrect parameter count in the call to native function 'JSON_ARRAY_INSERT'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a VARCHAR(400), b VARCHAR(400), KEY ((CONCAT(a, '')), (CONCAT(b, ''))))
< error
| 1071 42000 Specified key was too long; max key length is 3072 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a VARCHAR(400), b VARCHAR(400), UNIQUE KEY ((CONCAT(a, '')), (CONCAT(b, ''))))
< error
| 1071 42000 Specified key was too long; max key length is 3072 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a VARCHAR(400), b VARCHAR(400))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY ((CONCAT(a, '')), (CONCAT(b, '')))
< error
| 1071 42000 Specified key was too long; max key length is 3072 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a VARCHAR(400), b VARCHAR(400))
< ok 0 0
> exec
| ALTER TABLE a ADD UNIQUE KEY ((CONCAT(a, '')), (CONCAT(b, '')))
< error
| 1071 42000 Specified key was too long; max key length is 3072 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (id INT)
< ok 0 0
> exec
| ALTER TABLE a ADD KEY ((CONCAT(a, '')), (CONCAT(b, ''))), ADD a VARCHAR(400), ADD b VARCHAR(400)
< error
| 1071 42000 Specified key was too long; max key length is 3072 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a VARCHAR(400), b VARCHAR(400))
< ok 0 0
> exec
| CREATE INDEX idx_c ON a ((CONCAT(a, '')), (CONCAT(b, '')))
< error
| 1071 42000 Specified key was too long; max key length is 3072 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (x DOUBLE DEFAULT (ST_X()))
< error
| 1582 42000 Incorrect parameter count in the call to native function 'ST_X'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (x DOUBLE DEFAULT (ST_X(POINT(1, 2))))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (x DOUBLE DEFAULT (ST_X(POINT(1, 2), 3)))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (x INT DEFAULT (MBRWITHIN(POINT(1, 2))))
< error
| 1582 42000 Incorrect parameter count in the call to native function 'MBRWITHIN'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (x VARCHAR(100) DEFAULT (ROLES_GRAPHML('x')))
< error
| 1582 42000 Incorrect parameter count in the call to native function 'ROLES_GRAPHML'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (x INT DEFAULT (ST_CONTAINS(POINT(1, 2))))
< error
| 1582 42000 Incorrect parameter count in the call to native function 'ST_CONTAINS'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a TEXT, b TEXT, KEY ((SUBSTRING(a, 1, 400)), (SUBSTRING(b, 1, 400))))
< error
| 1071 42000 Specified key was too long; max key length is 3072 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (s VARCHAR(10) DEFAULT (TRIM(LEADING 'x' FROM 'xxx')))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (s VARCHAR(10) DEFAULT (TRIM(TRAILING 'x' FROM 'xxx')))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (s VARCHAR(10) DEFAULT (TRIM(BOTH 'x' FROM 'xxx')))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (s VARCHAR(10) DEFAULT (TRIM('x' FROM 'xxx')))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (id INT PRIMARY KEY, v INT) ENGINE=InnoDB
< ok 0 0
> exec
| ALTER TABLE a DISCARD TABLESPACE
< ok 0 0
> exec
| ALTER TABLE a ADD KEY k (v)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (id INT PRIMARY KEY, v INT) ENGINE=InnoDB
< ok 0 0
> exec
| ALTER TABLE a DISCARD TABLESPACE
< ok 0 0
> exec
| ALTER TABLE a ADD KEY k (v), ALGORITHM=COPY
< error
| 1814 HY000 Tablespace has been discarded for table 'a'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (id INT, v INT) ENGINE=InnoDB
< ok 0 0
> exec
| ALTER TABLE a DISCARD TABLESPACE
< ok 0 0
> exec
| ALTER TABLE a ADD PRIMARY KEY (id)
< error
| 1814 HY000 Tablespace has been discarded for table 'test/a'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (id INT PRIMARY KEY) ENGINE=InnoDB
< ok 0 0
> exec
| ALTER TABLE a DISCARD TABLESPACE
< ok 0 0
> exec
| ALTER TABLE a FORCE
< error
| 1814 HY000 Tablespace has been discarded for table 'test/a'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (id INT PRIMARY KEY, v INT) ENGINE=InnoDB
< ok 0 0
> exec
| ALTER TABLE a DISCARD TABLESPACE
< ok 0 0
> exec
| CREATE INDEX idx_v ON a (v)
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET sql_mode = IF(1, 'ANSI', '')
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (id INT, s VARCHAR(16) DEFAULT (REPLACE('abc', 'a', 'd')))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (d DOUBLE DEFAULT (ATAN2(1)))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (s VARBINARY(64) DEFAULT (AES_ENCRYPT('a', 'b', 'c', 'd', 'e', 'f')))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (s VARCHAR(64) DEFAULT (AES_DECRYPT('a', 'b', 'c', 'd', 'e', 'f')))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (v VARCHAR(20) CHARACTER SET utf8mb4, KEY k(v))
< ok 0 0
> exec
| ALTER TABLE a MODIFY COLUMN v VARCHAR(1000) CHARACTER SET utf8mb4
< error
| 1071 42000 Specified key was too long; max key length is 3072 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (v VARCHAR(800)) DEFAULT CHARSET=latin1, ENGINE=InnoDB
< ok 0 0
> exec
| ALTER TABLE a ADD KEY k(v)
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARACTER SET utf8mb4
< error
| 1071 42000 Specified key was too long; max key length is 3072 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (x INT, y INT DEFAULT (a.x + 1), KEY ((a.x + 1)))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (c1 TINYTEXT, c2 TEXT, c3 MEDIUMTEXT) DEFAULT CHARSET=latin1
< ok 0 0
> exec
| ALTER TABLE a CONVERT TO CHARACTER SET utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c1` text,
|   `c2` mediumtext,
|   `c3` longtext
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (x VARBINARY(600), y VARBINARY(600), KEY k(x, y)) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (v VARCHAR(200) CHARACTER SET utf8mb4) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| ALTER TABLE a ADD KEY k(v), ROW_FORMAT=DYNAMIC
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (j JSON UNIQUE)
< error
| 3152 42000 JSON column 'j' supports indexing only via generated columns on a specified JSON path.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (v VARCHAR(900) CHARACTER SET utf8mb4 PRIMARY KEY)
< error
| 1071 42000 Specified key was too long; max key length is 3072 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET sql_mode = ''
< ok 0 0
> exec
| CREATE TABLE a (v VARCHAR(900) CHARACTER SET utf8mb4, KEY i(v))
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `v` varchar(900) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
|   KEY `i` (`v`(768))
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (a VARCHAR(1) CHARACTER SET latin1, b VARCHAR(510) CHARACTER SET utf8mb4, c VARBINARY(1050), KEY ((CONCAT(a, b)), c))
< error
| 1071 42000 Specified key was too long; max key length is 3072 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (i INT DEFAULT 2.5E0, j INT DEFAULT 2.5, k TINYINT DEFAULT -128.5E0)
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `i` int DEFAULT '2',
|   `j` int DEFAULT '3',
|   `k` tinyint DEFAULT '-128'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET sql_mode = 'STRICT_ALL_TABLES'
< ok 0 0
> exec
| CREATE TABLE a (c ENUM('e', 'é') COLLATE latin1_general_ci)
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET sql_mode = 'STRICT_ALL_TABLES'
< ok 0 0
> exec
| CREATE TABLE a (c ENUM('😀', '😁') COLLATE utf8mb4_unicode_520_ci)
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20))
< ok 0 0
> exec
| ALTER TABLE a DISCARD PARTITION p0 TABLESPACE
< ok 0 0
> exec
| ALTER TABLE a IMPORT PARTITION p1 TABLESPACE
< error
| 1813 HY000 Tablespace 'test/a#p#p1' exists.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20))
< ok 0 0
> exec
| ALTER TABLE a DISCARD PARTITION p0 TABLESPACE
< ok 0 0
> exec
| ALTER TABLE a TRUNCATE PARTITION p0
< error
| 1814 HY000 Tablespace has been discarded for table './test/a'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (id INT) PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20))
< ok 0 0
> exec
| ALTER TABLE a DISCARD PARTITION p0 TABLESPACE
< ok 0 0
> exec
| ALTER TABLE a TRUNCATE PARTITION p1
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET sql_mode = CASE 'ansi' WHEN 'ANSI' THEN 'TRADITIONAL' ELSE '' END
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET sql_mode = IF(0 = 'abc', 'ANSI', '')
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET sql_mode = IF(TRIM(X'09') = X'', 'ANSI', '')
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET sql_mode = IF(CHAR_LENGTH(X'C3A9') = 2, 'ANSI', '')
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET sql_mode = IF(1, 'ANSI', sim_oracle_missing.no_such())
< error
| 1305 42000 FUNCTION sim_oracle_missing.no_such does not exist
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET sql_mode = CAST('ANSI' AS CHAR)
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET sql_mode = IFNULL('ANSI', sim_oracle_missing.no_such())
< error
| 1305 42000 FUNCTION sim_oracle_missing.no_such does not exist
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET sql_mode = IF(IF(1, 1, 'x') = '01', 'ANSI', '')
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (id INT PRIMARY KEY, pid INT, CONSTRAINT fk_self FOREIGN KEY (pid) REFERENCES a(id))
< ok 0 0
> exec
| ALTER TABLE a DISCARD TABLESPACE
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (x DOUBLE DEFAULT (ATAN2(1)))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (b VARBINARY(1), KEY ((CONCAT(b, X'00'))))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (b VARBINARY(1), KEY ((CONCAT(b, b'0000000100000001'))))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (v VARCHAR(10), KEY ((CONCAT(v, X'4142'))))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (b VARBINARY(1), KEY ((CONCAT(b, REPEAT(X'00', 600)))))
< error
| 3757 HY000 Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| ALTER UNDO TABLESPACE MYSQL SET INACTIVE
< error
| 3510 HY000 Tablespace MYSQL doesn't exist.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| DROP UNDO TABLESPACE MYSQL
< error
| 3510 HY000 Tablespace MYSQL doesn't exist.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| DROP UNDO TABLESPACE INNODB_SYSTEM
< error
| 3510 HY000 Tablespace INNODB_SYSTEM doesn't exist.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| ALTER UNDO TABLESPACE mysql SET INACTIVE
< error
| 3119 42000 Cannot ALTER UNDO TABLESPACE `mysql` because it is a general tablespace.  Please use ALTER TABLESPACE.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| DROP UNDO TABLESPACE mysql
< error
| 3119 42000 Cannot DROP UNDO TABLESPACE `mysql` because it is a general tablespace.  Please use DROP TABLESPACE.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| DROP UNDO TABLESPACE innodb_system
< error
| 3119 42000 InnoDB: `innodb_system` is a reserved tablespace name.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| DROP UNDO TABLESPACE innodb_temporary
< error
| 3119 42000 InnoDB: `innodb_temporary` is a reserved tablespace name.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| DROP UNDO TABLESPACE innodb_x
< error
| 3510 HY000 Tablespace innodb_x doesn't exist.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| ALTER UNDO TABLESPACE INNODB_X SET INACTIVE
< error
| 3510 HY000 Tablespace INNODB_X doesn't exist.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TEMPORARY TABLE z (a INT AUTO_INCREMENT)
< error
| 1075 42000 Incorrect table definition; there can be only one auto column and it must be defined as a key
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (b VARBINARY(1), KEY ((CONCAT(b, REPEAT(X'00', 128)))))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (b VARBINARY(1), KEY ((CONCAT(b, REPEAT(X'00', 512)))))
< error
| 3757 HY000 Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (b VARBINARY(10), KEY ((REPEAT(b, 51))))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (b VARBINARY(10), KEY ((REPEAT(b, 52))))
< error
| 3757 HY000 Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (v VARCHAR(2), KEY ((REPEAT(v, 9223372036854775807))))
< error
| 3757 HY000 Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, KEY i ((c + 1)))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY j ((`!hidden!i!0!0` + 1))
< error
| 1054 42S22 Unknown column '!hidden!i!0!0' in 'functional index'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, KEY i ((c + 1)))
< ok 0 0
> exec
| CREATE INDEX j ON a ((`!hidden!i!0!0` + 1))
< error
| 1054 42S22 Unknown column '!hidden!i!0!0' in 'functional index'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, s VARCHAR(300)) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| CREATE INDEX i ON a ((LOWER(s)))
< error
| 1071 42000 Specified key was too long; max key length is 767 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, s VARCHAR(300)) ROW_FORMAT=REDUNDANT
< ok 0 0
> exec
| ALTER TABLE a ADD KEY ((LOWER(s)))
< error
| 1071 42000 Specified key was too long; max key length is 767 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, s VARCHAR(300)) ROW_FORMAT=DYNAMIC
< ok 0 0
> exec
| CREATE INDEX i ON a ((LOWER(s)))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (s VARCHAR(191), KEY ((LOWER(s)))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (s VARCHAR(192), KEY ((LOWER(s)))) ROW_FORMAT=COMPACT
< error
| 1071 42000 Specified key was too long; max key length is 767 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| CREATE TABLE z (s VARCHAR(300), KEY ((LOWER(s)))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT, s VARCHAR(10), b VARBINARY(1), KEY ((RPAD(b, 512, X'00')))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT, s VARCHAR(10), b VARBINARY(1), KEY ((LPAD(b, 512, X'00')))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT, s VARCHAR(10), b VARBINARY(1), KEY ((RPAD(b, 513, X'00')))) ROW_FORMAT=COMPACT
< error
| 3757 HY000 Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT, s VARCHAR(10), b VARBINARY(1), KEY ((a = 1))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT, s VARCHAR(10), b VARBINARY(1), KEY ((a < 1 AND a > 0))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT, s VARCHAR(10), b VARBINARY(1), KEY ((NOT a))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT, s VARCHAR(10), b VARBINARY(1), KEY ((a IS NULL))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT, s VARCHAR(10), b VARBINARY(1), KEY ((a IN (1, 2)))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT, s VARCHAR(10), b VARBINARY(1), KEY ((a BETWEEN 1 AND 2))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT, s VARCHAR(10), b VARBINARY(1), KEY ((s LIKE 'x%'))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT, s VARCHAR(10), b VARBINARY(1), KEY ((a <=> 1))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT, KEY ((-a))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (s VARCHAR(10), KEY ((-s))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a TINYINT, KEY ((-a))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT, KEY ((~a))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT, KEY ((+a))) ROW_FORMAT=COMPACT
< error
| 3762 HY000 Functional index on a column is not supported. Consider using a regular index instead.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (d DECIMAL(10,2), KEY ((-d))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT, KEY ((1))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT, KEY ((1.5))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT, KEY ((1e0))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT, KEY (('abc'))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT, KEY ((TRUE))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT, KEY ((DATE '2024-01-01'))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT, KEY ((REPEAT('x', 200)))) ROW_FORMAT=COMPACT
< error
| 1071 42000 Specified key was too long; max key length is 767 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (s VARCHAR(100), KEY ((REPLACE(s, 'x', 'y')))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (s VARCHAR(100), KEY ((REPLACE(s, 'x', 'yy')))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (s VARCHAR(100), KEY ((REPLACE(s, 'xx', 'yyy')))) ROW_FORMAT=COMPACT
< error
| 1071 42000 Specified key was too long; max key length is 767 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (s VARCHAR(191), KEY ((REPLACE(s, 'x', '')))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (s VARCHAR(192), KEY ((REPLACE(s, 'x', '')))) ROW_FORMAT=COMPACT
< error
| 1071 42000 Specified key was too long; max key length is 767 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (s VARCHAR(100), KEY ((REPLACE(s, '', 'yyyyy')))) ROW_FORMAT=COMPACT
< error
| 1071 42000 Specified key was too long; max key length is 767 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (s VARCHAR(10), KEY ((SHA1(s)))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (s VARCHAR(10), KEY ((MD5(s)))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (s VARCHAR(10), KEY ((HEX(s)))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (s VARCHAR(10), KEY (s(0)))
< error
| 1391 HY000 Key part 's' length cannot be 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (s VARCHAR(10), KEY (nosuch(0)))
< error
| 1391 HY000 Key part 'nosuch' length cannot be 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (s VARCHAR(10), PRIMARY KEY (s(0)))
< error
| 1391 HY000 Key part 's' length cannot be 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD KEY (c(0))
< error
| 1391 HY000 Key part 'c' length cannot be 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| CREATE INDEX i ON a (c(0))
< error
| 1391 HY000 Key part 'c' length cannot be 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE INDEX i ON nosuch (s(0))
< error
| 1391 HY000 Key part 's' length cannot be 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (s VARCHAR(10), KEY (s(11)))
< error
| 1089 HY000 Incorrect prefix key; the used key part isn't a string, the used length is longer than the key part, or the storage engine doesn't support unique prefix keys
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT, KEY ((SPACE(512))))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT, KEY ((SPACE(513))))
< error
| 3757 HY000 Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT, KEY ((SPACE(9223372036854775807))))
< error
| 3757 HY000 Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT, KEY ((SPACE(a))))
< error
| 3757 HY000 Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT, KEY ((SPACE(191)))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (a INT, KEY ((SPACE(192)))) ROW_FORMAT=COMPACT
< error
| 1071 42000 Specified key was too long; max key length is 767 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (c1 VARCHAR(10), c2 VARCHAR(10), c3 VARCHAR(10), c4 VARCHAR(10), c5 VARCHAR(10), c6 VARCHAR(10), c7 VARCHAR(10), c8 VARCHAR(10), c9 VARCHAR(10), c10 VARCHAR(10), c11 VARCHAR(10), c12 VARCHAR(10), c13 VARCHAR(10), KEY ((ASCII(c1)), (ASCII(c2)), (ASCII(c3)), (ASCII(c4)), (ASCII(c5)), (ASCII(c6)), (ASCII(c7)), (ASCII(c8)), (ASCII(c9)), (ASCII(c10)), (ASCII(c11)), (ASCII(c12)), (ASCII(c13))))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (s VARCHAR(10), KEY ((ASCII(s)))) ROW_FORMAT=COMPACT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET sql_mode = IF(1, 'ANSI', no_such())
< error
| 1305 42000 FUNCTION test.no_such does not exist
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET @x = IF(1, 1, no_such())
< error
| 1305 42000 FUNCTION test.no_such does not exist
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET sql_mode = IF(1, 'ANSI', COUNT(1))
< error
| 1111 HY000 Invalid use of group function
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET sql_mode = IF(1, 'ANSI', ROW_NUMBER() OVER ())
< error
| 3593 HY000 You cannot use the window function 'row_number' in this context.'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET sql_mode = IF(1, 'ANSI', VERSION())
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET sql_mode = IF(1, 'ANSI', REPEAT('a', 2))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE FUNCTION f_set_branch() RETURNS INT DETERMINISTIC RETURN 1
< ok 0 0
> exec
| SET @x = IF(1, 1, f_set_branch())
< ok 0 0
> exec
| DROP FUNCTION IF EXISTS f_set_branch
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (s VARCHAR(10) DEFAULT (REPEAT('a', 3)))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (s VARCHAR(10) DEFAULT (REPEAT('a')))
< error
| 1064 42000 …near ')))' at line 1
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (s VARCHAR(64) DEFAULT (CREATE_DIGEST('SHA256', 'a')))
< error
| 3770 HY000 Default value expression of column 's' contains a disallowed function: `CREATE_DIGEST`.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLESPACE ts_impl1 ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLESPACE ts_impl2 ADD DATAFILE 'ts_impl1.ibd' ENGINE=InnoDB
< ok 0 0
> exec
| DROP TABLESPACE ts_impl2
< ok 0 0
> exec
| DROP TABLESPACE ts_impl1
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((REVERSE(b))))
< error
| 3757 HY000 Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((LEFT(b, 10))))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((SUBSTRING(b, 1, 10))))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((CONCAT(b, 'x'))))
< error
| 3757 HY000 Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((LOWER(b))))
< error
| 3757 HY000 Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((LEFT(b, 512))))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((LEFT(b, 513))))
< error
| 3757 HY000 Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((REVERSE(tb))))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((UPPER(tb))))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((CONCAT(tb, tb, 'x'))))
< error
| 3757 HY000 Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((REVERSE(vb))))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((CONCAT(vb, 'x'))))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((REVERSE(bn))))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((REVERSE(vbw))))
< error
| 3757 HY000 Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((CONCAT(vbw, vbw))))
< error
| 3757 HY000 Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((LEFT(vbw, 512))))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((LEFT(vbw, 513))))
< error
| 3757 HY000 Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((REVERSE(tt))))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((REVERSE(ttl))))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((LEFT(tx, 10))))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((REVERSE(v))))
< error
| 3757 HY000 Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((LEFT(v, 512))))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((LEFT(v, 513))))
< error
| 3757 HY000 Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((REVERSE(l))))
< error
| 3757 HY000 Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((CAST(i AS CHAR(512)))))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((CAST(i AS CHAR(513)))))
< error
| 3757 HY000 Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((CAST(i AS BINARY(513)))))
< error
| 3757 HY000 Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((CONCAT(v5, v5))))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((CONCAT(v5, v5, v5))))
< error
| 3757 HY000 Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((LOWER(v))))
< error
| 3757 HY000 Cannot create a functional index on an expression that returns a BLOB or TEXT. Please consider using CAST.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (i INT, b BLOB, tb TINYBLOB, vb VARBINARY(10), vbw VARBINARY(4000), bn BINARY(4), tt TINYTEXT, ttl TINYTEXT CHARACTER SET latin1, tx TEXT, v VARCHAR(600), l VARCHAR(600) CHARACTER SET latin1, v5 VARCHAR(256), w1 VARCHAR(500), w2 VARCHAR(500), KEY ((LOWER(w1)), (LOWER(w2))))
< error
| 1071 42000 Specified key was too long; max key length is 3072 bytes
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE INDEX i ON a ((CONNECTION_ID()))
< error
| 3758 HY000 Expression of functional index 'i' contains a disallowed function.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE INDEX i ON a ((UTC_TIMESTAMP()))
< error
| 3758 HY000 Expression of functional index 'i' contains a disallowed function.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE INDEX i ON a ((UTC_TIMESTAMP(3)))
< error
| 3758 HY000 Expression of functional index 'i' contains a disallowed function.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE INDEX i ON a ((CURRENT_DATE))
< error
| 3758 HY000 Expression of functional index 'i' contains a disallowed function.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE INDEX i ON a ((CURRENT_USER()))
< error
| 3758 HY000 Expression of functional index 'i' contains a disallowed function.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE INDEX i ON a ((DATABASE()))
< error
| 3758 HY000 Expression of functional index 'i' contains a disallowed function.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE INDEX i ON a ((UNIX_TIMESTAMP(c)))
< error
| 3758 HY000 Expression of functional index 'i' contains a disallowed function.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE INDEX i ON a ((UUID_SHORT()))
< error
| 3758 HY000 Expression of functional index 'i' contains a disallowed function.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE INDEX i ON a ((STATEMENT_DIGEST_TEXT('select 1')))
< error
| 3758 HY000 Expression of functional index 'i' contains a disallowed function.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE INDEX i ON a ((LAST_INSERT_ID()))
< error
| 3758 HY000 Expression of functional index 'i' contains a disallowed function.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE INDEX i ON a ((VERSION()))
< error
| 3758 HY000 Expression of functional index 'i' contains a disallowed function.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE INDEX i ON a ((GET_LOCK('a', 1)))
< error
| 3758 HY000 Expression of functional index 'i' contains a disallowed function.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE INDEX i ON a ((VALUES(c)))
< error
| 3758 HY000 Expression of functional index 'i' contains a disallowed function.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE INDEX i ON a ((NAME_CONST('a', 1)))
< error
| 3758 HY000 Expression of functional index 'i' contains a disallowed function.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE INDEX i ON a ((nosuchfn(c)))
< error
| 3758 HY000 Expression of functional index 'i' contains a disallowed function.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE INDEX i ON a ((d.nosuchfn(c)))
< error
| 3758 HY000 Expression of functional index 'i' contains a disallowed function.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE INDEX i ON a ((MATCH(c) AGAINST ('x')))
< error
| 3758 HY000 Expression of functional index 'i' contains a disallowed function.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE INDEX i ON a ((@x))
< error
| 3772 HY000 Default value expression of column 'i' cannot refer user or system variables.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE INDEX i ON a ((@x + NOW()))
< error
| 3772 HY000 Default value expression of column 'i' cannot refer user or system variables.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE INDEX i ON a ((NOW() + @x))
< error
| 3758 HY000 Expression of functional index 'i' contains a disallowed function.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE INDEX i ON a ((SLEEP(@x)))
< error
| 3772 HY000 Default value expression of column 'i' cannot refer user or system variables.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE INDEX i ON a ((COUNT(c)))
< error
| 1111 HY000 Invalid use of group function
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE INDEX i ON a ((GROUP_CONCAT(c)))
< error
| 1111 HY000 Invalid use of group function
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE INDEX i ON a ((COUNT(c) + NOW()))
< error
| 3758 HY000 Expression of functional index 'i' contains a disallowed function.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE INDEX i ON a ((ROW_NUMBER() OVER ()))
< error
| 3593 HY000 You cannot use the window function 'row_number' in this context.'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE INDEX i ON a ((COUNT(c) OVER ()))
< error
| 3593 HY000 You cannot use the window function 'count' in this context.'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE INDEX i ON a ((ANY_VALUE(c)))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE INDEX i ON a ((CHAR(65)))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE INDEX i ON a ((ELT(1, c)))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX i ON a ((LOWER(c)))
< error
| 3759 HY000 Fulltext functional index is not supported.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX i ON a (c, (LOWER(c)))
< error
| 3759 HY000 Fulltext functional index is not supported.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX i ON a ((CAST(c AS JSON)))
< error
| 3759 HY000 Fulltext functional index is not supported.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX i ON a ((c))
< error
| 3762 HY000 Functional index on a column is not supported. Consider using a regular index instead.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX i ON a ((NOW()))
< error
| 3758 HY000 Expression of functional index 'i' contains a disallowed function.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX i ON a ((COUNT(c)))
< error
| 1111 HY000 Invalid use of group function
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| CREATE FULLTEXT INDEX i ON a ((LOWER(nosuch)))
< error
| 1054 42S22 Unknown column 'nosuch' in 'functional index'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE s (g GEOMETRY NOT NULL SRID 0)
< ok 0 0
> exec
| CREATE SPATIAL INDEX i ON s ((ST_SRID(g)))
< error
| 3760 HY000 Spatial functional index is not supported.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE s (g GEOMETRY NOT NULL SRID 0)
< ok 0 0
> exec
| CREATE SPATIAL INDEX i ON s ((g))
< error
| 3762 HY000 Functional index on a column is not supported. Consider using a regular index instead.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT ((LOWER(c)))
< error
| 3759 HY000 Fulltext functional index is not supported.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT ((c))
< error
| 3762 HY000 Functional index on a column is not supported. Consider using a regular index instead.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| ALTER TABLE a ADD FULLTEXT nosuch_idx (nosuch, (LOWER(c)))
< error
| 1072 42000 Key column 'nosuch' doesn't exist in table
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c VARCHAR(10))
< ok 0 0
> exec
| ALTER TABLE a ADD SPATIAL ((LOWER(c)))
< error
| 3760 HY000 Spatial functional index is not supported.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (c VARCHAR(10), FULLTEXT ((LOWER(c))))
< error
| 3759 HY000 Fulltext functional index is not supported.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (c VARCHAR(10), FULLTEXT ((c)))
< error
| 3762 HY000 Functional index on a column is not supported. Consider using a regular index instead.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (c VARCHAR(10), FULLTEXT ((NOW())))
< error
| 3758 HY000 Expression of functional index 'functional_index' contains a disallowed function.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (c VARCHAR(10), FULLTEXT ((LOWER(nosuch))))
< error
| 1054 42S22 Unknown column 'nosuch' in 'functional index'
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE TABLE z (c VARCHAR(10), SPATIAL ((LOWER(c))))
< error
| 3760 HY000 Spatial functional index is not supported.
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET NAMES latin1
< ok 0 0
> exec
| CREATE TABLE a (c CHAR(10) CHARACTER SET utf8mb4 DEFAULT 'é')
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `c` char(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'Ã©'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET SESSION foreign_key_checks = 1
< ok 0 0
> exec
| SET time_zone = 'America/New_York'
< ok 0 0
> exec
| CREATE TABLE a (dt DATETIME DEFAULT '2026-06-01 12:00:00+00:00')
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `dt` datetime DEFAULT '2026-06-01 08:00:00'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET SESSION foreign_key_checks = 1
< ok 0 0
> exec
| CREATE TABLESPACE ts_sdl_plain ADD DATAFILE 'ts_sdl_plain.ibd'
< ok 0 0
> exec
| CREATE TABLE a (id INT PRIMARY KEY) TABLESPACE ts_sdl_plain
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `id` int NOT NULL,
|   PRIMARY KEY (`id`)
| ) /*!50100 TABLESPACE `ts_sdl_plain` */ ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET SESSION foreign_key_checks = 1
< ok 0 0
> exec
| INSTALL PLUGIN keyring_file SONAME 'keyring_file.so'
< ok 0 0
> exec
| CREATE TABLESPACE ts_enc_chain ADD DATAFILE 'ts_enc_chain.ibd'
< ok 0 0
> exec
| CREATE TABLE a (id INT PRIMARY KEY) TABLESPACE ts_enc_chain
< ok 0 0
> exec
| ALTER TABLESPACE ts_enc_chain ENCRYPTION='Y'
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN c INT
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `id` int NOT NULL,
|   `c` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) /*!50100 TABLESPACE `ts_enc_chain` */ ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci /*!80016 ENCRYPTION='Y' */
> exec
| ALTER TABLESPACE ts_enc_chain ENCRYPTION='N'
< ok 0 0
> exec
| ALTER TABLE a ADD COLUMN c2 INT
< ok 0 0
> query
| SHOW CREATE TABLE a
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| a
: bytes
| CREATE TABLE `a` (
|   `id` int NOT NULL,
|   `c` int DEFAULT NULL,
|   `c2` int DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) /*!50100 TABLESPACE `ts_enc_chain` */ ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| SET SESSION default_table_encryption = ON
< ok 0 0
> exec
| CREATE TABLESPACE ts_def_enc ADD DATAFILE 'ts_def_enc.ibd'
< ok 0 0
> exec
| CREATE TABLE a_def_enc (id INT) TABLESPACE ts_def_enc
< error
| 3825 HY000 Request to create 'unencrypted' table while using an 'encrypted' tablespace.
> exec
| CREATE TABLESPACE ts_omit_enc ADD DATAFILE 'ts_omit_enc.ibd' ENCRYPTION='Y'
< ok 0 0
> exec
| CREATE TABLE a_def_omit_y (id INT) TABLESPACE ts_omit_enc
< error
| 3825 HY000 Request to create 'unencrypted' table while using an 'encrypted' tablespace.
> exec
| CREATE TABLE a_def_enc (id INT) TABLESPACE ts_def_enc ENCRYPTION='Y'
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> query
| SELECT DATABASE()
< result
< column "DATABASE()" "VARCHAR"
< row
: bytes
| test
> exec
| CREATE DATABASE db_cs_bin CHARACTER SET latin1 COLLATE latin1_bin
< ok 1 0
> exec
| USE db_cs_bin
< ok 0 0
> exec
| SET CHARACTER SET utf8mb4
< ok 0 0
> query
| SELECT @@character_set_client, @@collation_connection
< result
< column "@@character_set_client" "VARCHAR"
< column "@@collation_connection" "VARCHAR"
< row
: bytes
| utf8mb4
: bytes
| latin1_bin
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET SESSION foreign_key_checks = 1
< ok 0 0
> exec
| SET foreign_key_checks = 0, sql_mode = 'NO_SUCH_MODE'
< error
| 1231 42000 Variable 'sql_mode' can't be set to the value of 'NO_SUCH_MODE'
> query
| SELECT @@foreign_key_checks
< result
< column "@@foreign_key_checks" "BIGINT"
< row
: int64 1
> exec
| SET @x_seq = 1, @y_seq = @x_seq + 1
< ok 0 0
> query
| SELECT @y_seq
< result
< column "@y_seq" "DOUBLE" precision 9223372036854775807 scale 9223372036854775807
< row
: null
> exec
| SET SESSION sql_mode = REPLACE(@@sql_mode, 'ONLY_FULL_GROUP_BY', '')
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> query
| SELECT DATABASE()
< result
< column "DATABASE()" "VARCHAR"
< row
: bytes
| db_cs_bin
> exec
| CREATE DATABASE d_qfn; USE d_qfn;
< ok 0 0
> exec
| SET sql_mode = d_qfn.CONCAT('ANSI')
< error
| 1305 42000 FUNCTION d_qfn.CONCAT does not exist
> exec
| CREATE FUNCTION d_qfn.CONCAT(s VARCHAR(20)) RETURNS VARCHAR(20) DETERMINISTIC RETURN ''
< ok 0 0
> exec
| SET sql_mode = d_qfn.CONCAT('ANSI')
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| CREATE DATABASE d_qual_own; USE d_qual_own;
< ok 0 0
> exec
| CREATE TABLE q_own (a INT, b INT DEFAULT (d_qual_own.q_own.a + 1), KEY ((d_qual_own.q_own.a + 1)))
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| SET sql_mode = IF(1, 'ANSI', @@no_such_variable)
< error
| 1193 HY000 Unknown system variable 'no_such_variable'
> exec
| SET sql_mode = IF(1e2 = 100, 'ANSI', '')
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| REAL_AS_FLOAT,PIPES_AS_CONCAT,ANSI_QUOTES,IGNORE_SPACE,ONLY_FULL_GROUP_BY,ANSI
> exec
| SET sql_mode = IF('123euro' = 123, 'STRICT_ALL_TABLES', '')
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| STRICT_ALL_TABLES
> exec
| SET sql_mode = IF(1e2 + 1 = 101, 'ANSI', '')
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| REAL_AS_FLOAT,PIPES_AS_CONCAT,ANSI_QUOTES,IGNORE_SPACE,ONLY_FULL_GROUP_BY,ANSI
> exec
| SET sql_mode = 1e2 + 1
< error
| 1232 42000 Incorrect argument type to variable 'sql_mode'
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| SET sql_mode = ELT(1, 'ANSI', '')
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| REAL_AS_FLOAT,PIPES_AS_CONCAT,ANSI_QUOTES,IGNORE_SPACE,ONLY_FULL_GROUP_BY,ANSI
> exec
| SET sql_mode = ELT('2', '', 'STRICT_ALL_TABLES')
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| STRICT_ALL_TABLES
> exec
| SET sql_mode = ELT('1.6', 'ANSI', 'STRICT_ALL_TABLES')
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| REAL_AS_FLOAT,PIPES_AS_CONCAT,ANSI_QUOTES,IGNORE_SPACE,ONLY_FULL_GROUP_BY,ANSI
> exec
| SET sql_mode = ELT(1.6, 'ANSI', 'STRICT_ALL_TABLES')
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| STRICT_ALL_TABLES
> exec
| SET sql_mode = ELT(2.5, 'ANSI', 'TRADITIONAL', 'STRICT_ALL_TABLES')
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| STRICT_ALL_TABLES
> exec
| SET sql_mode = ELT(2.5e0, 'ANSI', 'TRADITIONAL', 'STRICT_ALL_TABLES')
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| STRICT_TRANS_TABLES,STRICT_ALL_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,TRADITIONAL,NO_ENGINE_SUBSTITUTION
> exec
| SET sql_mode = ELT(' 2', 'ANSI', '')
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
|
> exec
| SET sql_mode = ELT(1, 'ANSI', @@sql_mode)
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| REAL_AS_FLOAT,PIPES_AS_CONCAT,ANSI_QUOTES,IGNORE_SPACE,ONLY_FULL_GROUP_BY,ANSI
> exec
| SET sql_mode = ELT(4294967297, 'ANSI', 'TRADITIONAL')
< error
| 1231 42000 Variable 'sql_mode' can't be set to the value of 'NULL'
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| REAL_AS_FLOAT,PIPES_AS_CONCAT,ANSI_QUOTES,IGNORE_SPACE,ONLY_FULL_GROUP_BY,ANSI
> exec
| SET sql_mode = ELT(0, 'ANSI')
< error
| 1231 42000 Variable 'sql_mode' can't be set to the value of 'NULL'
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| REAL_AS_FLOAT,PIPES_AS_CONCAT,ANSI_QUOTES,IGNORE_SPACE,ONLY_FULL_GROUP_BY,ANSI
> exec
| SET sql_mode = ELT(-1, 'ANSI')
< error
| 1231 42000 Variable 'sql_mode' can't be set to the value of 'NULL'
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| REAL_AS_FLOAT,PIPES_AS_CONCAT,ANSI_QUOTES,IGNORE_SPACE,ONLY_FULL_GROUP_BY,ANSI
> exec
| SET sql_mode = ELT(NULL, 'ANSI')
< error
| 1231 42000 Variable 'sql_mode' can't be set to the value of 'NULL'
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| REAL_AS_FLOAT,PIPES_AS_CONCAT,ANSI_QUOTES,IGNORE_SPACE,ONLY_FULL_GROUP_BY,ANSI
> exec
| SET sql_mode = ELT('x', 'ANSI')
< error
| 1231 42000 Variable 'sql_mode' can't be set to the value of 'NULL'
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| REAL_AS_FLOAT,PIPES_AS_CONCAT,ANSI_QUOTES,IGNORE_SPACE,ONLY_FULL_GROUP_BY,ANSI
> exec
| SET sql_mode = ELT(1, NULL, 'ANSI')
< error
| 1231 42000 Variable 'sql_mode' can't be set to the value of 'NULL'
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| REAL_AS_FLOAT,PIPES_AS_CONCAT,ANSI_QUOTES,IGNORE_SPACE,ONLY_FULL_GROUP_BY,ANSI
> exec
| SET sql_mode = ELT(1)
< error
| 1582 42000 Incorrect parameter count in the call to native function 'ELT'
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| REAL_AS_FLOAT,PIPES_AS_CONCAT,ANSI_QUOTES,IGNORE_SPACE,ONLY_FULL_GROUP_BY,ANSI
> exec
| SET sql_mode = ELT(1, 'ANSI', @@no_such_variable)
< error
| 1193 HY000 Unknown system variable 'no_such_variable'
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| REAL_AS_FLOAT,PIPES_AS_CONCAT,ANSI_QUOTES,IGNORE_SPACE,ONLY_FULL_GROUP_BY,ANSI
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| ALTER TABLESPACE mysql AUTOEXTEND_SIZE=4M
< ok 0 0
> exec
| ALTER TABLESPACE mysql ENGINE=InnoDB
< ok 0 0
> exec
| ALTER TABLESPACE mysql AUTOEXTEND_SIZE=3M
< error
| 4025 HY000 AUTOEXTEND_SIZE value should be between 4M and 4096M
> exec
| ALTER TABLESPACE mysql ENGINE=MyISAM
< error
| 3642 HY000 Engine 'MyISAM' does not match stored engine 'InnoDB' for tablespace 'mysql'
> exec
| ALTER TABLESPACE mysql ENCRYPTION='Z' ENGINE_ATTRIBUTE='bad'
< error
| 3980 HY000 Invalid json attribute, error: "Invalid value." at pos 0: 'bad'
> exec
| ALTER TABLESPACE mysql ADD DATAFILE 'x.ibd'
< error
| 1533 HY000 Failed to alter: TABLESPACE mysql
> exec
| ALTER TABLESPACE mysql DROP DATAFILE 'x.ibd'
< error
| 3629 HY000 Tablespace 'mysql' does not have a file named 'x.ibd'
> exec
| ALTER TABLESPACE mysql DROP DATAFILE 'mysql.ibd'
< error
| 1533 HY000 Failed to alter: TABLESPACE mysql
> exec
| ALTER TABLESPACE mysql RENAME TO x
< error
| 3119 42000 InnoDB: `mysql` is a reserved tablespace name.
> exec
| ALTER TABLESPACE mysql RENAME TO MYSQL
< error
| 3119 42000 InnoDB: `mysql` is a reserved tablespace name.
> exec
| ALTER TABLESPACE mysql RENAME TO mysql
< error
| 1813 HY000 Tablespace 'mysql' exists.
> exec
| ALTER TABLESPACE mysql RENAME TO innodb_temporary
< error
| 1813 HY000 Tablespace 'innodb_temporary' exists.
> exec
| ALTER TABLESPACE MYSQL ENCRYPTION='N'
< error
| 3510 HY000 Tablespace MYSQL doesn't exist.
> exec
| ALTER TABLESPACE mysql ENCRYPTION='Y'
< ok 0 0
> query
| SELECT ENCRYPTION FROM information_schema.INNODB_TABLESPACES WHERE NAME = 'mysql'
< result
< column "ENCRYPTION" "VARCHAR"
< row
: bytes
| Y
> exec
| ALTER TABLESPACE mysql ENCRYPTION='N'
< ok 0 0
> query
| SELECT ENCRYPTION FROM information_schema.INNODB_TABLESPACES WHERE NAME = 'mysql'
< result
< column "ENCRYPTION" "VARCHAR"
< row
: bytes
| N
> exec
| SET @x = CAST(3.6 AS SIGNED)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 4
> exec
| SET @x = CAST(-3.6 AS SIGNED)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| -4
> exec
| SET @x = CAST(2.5 AS SIGNED)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 3
> exec
| SET @x = CAST(-2.5 AS SIGNED)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| -3
> exec
| SET @x = CAST(2.5e0 AS SIGNED)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 2
> exec
| SET @x = CAST(3.5e0 AS SIGNED)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 4
> exec
| SET @x = CAST(-1.5e0 AS SIGNED)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| -2
> exec
| SET @x = CAST('3.6' AS SIGNED)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 3
> exec
| SET @x = CAST(' 3.6x' AS SIGNED)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 3
> exec
| SET @x = CAST('1e2' AS SIGNED)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 1
> exec
| SET @x = CAST(3.6 AS UNSIGNED)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 4
> exec
| SET @x = CAST(-1.6 AS UNSIGNED)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 18446744073709551614
> exec
| SET @x = CAST(-1 AS UNSIGNED)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 18446744073709551615
> exec
| SET @x = CAST('-1' AS UNSIGNED)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 18446744073709551615
> exec
| SET @x = CAST(99999999999999999999 AS SIGNED)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 9223372036854775807
> exec
| SET @x = CAST(-99999999999999999999 AS SIGNED)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| -9223372036854775808
> exec
| SET @x = CAST(99999999999999999999 AS UNSIGNED)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 18446744073709551615
> exec
| SET @x = CAST(-99999999999999999999 AS UNSIGNED)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 9223372036854775808
> exec
| SET @x = CAST('99999999999999999999' AS SIGNED)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| -1
> exec
| SET @x = CAST('10000000000000000000' AS SIGNED)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| -8446744073709551616
> exec
| SET @x = CAST('-99999999999999999999' AS SIGNED)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| -9223372036854775808
> exec
| SET @x = CAST('-99999999999999999999' AS UNSIGNED)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 9223372036854775808
> exec
| SET @x = CAST(18446744073709551615 AS SIGNED)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| -1
> exec
| SET @x = CAST(9223372036854775808 AS SIGNED)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| -9223372036854775808
> exec
| SET @x = CAST(9223372036854775807.5 AS SIGNED)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 9223372036854775807
> exec
| SET @x = CAST(1e30 AS SIGNED)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 9223372036854775807
> exec
| SET @x = CAST(1e30 AS UNSIGNED)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 9223372036854775807
> exec
| SET @x = CAST(1.5e19 AS UNSIGNED)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 9223372036854775807
> exec
| SET @x = CAST(-2.5e0 AS UNSIGNED)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 18446744073709551614
> exec
| SET @x = CAST(0.5 AS SIGNED)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 1
> exec
| SET @x = CAST(-0.5 AS UNSIGNED)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 18446744073709551615
> exec
| SET @x = SUBSTRING('ANSI', 2, 9223372036854775807)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "LONGTEXT"
< row
: bytes
| NSI
> exec
| SET @x = SUBSTRING('ANSI', -2, 9223372036854775807)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "LONGTEXT"
< row
: bytes
| SI
> exec
| SET @x = SUBSTRING(_binary'ANSI', 2, 9223372036854775807)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "LONGTEXT"
< row
: bytes
| NSI
> exec
| SET @x = SUBSTRING('ANSI', 9223372036854775807, 2)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "LONGTEXT"
< row
: bytes
|
> exec
| SET @x = SUBSTRING('ANSI', -9223372036854775807, 2)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "LONGTEXT"
< row
: bytes
|
> exec
| SET @x = SUBSTRING('ANSI', 2, 18446744073709551615)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "LONGTEXT"
< row
: bytes
| NSI
> exec
| SET @x = SUBSTRING('ANSI', 18446744073709551615)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "LONGTEXT"
< row
: bytes
|
> exec
| SET @x = SUBSTRING('ANSI', 2, -1)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "LONGTEXT"
< row
: bytes
|
> exec
| SET @x = LEFT('ANSI', 9223372036854775807)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "LONGTEXT"
< row
: bytes
| ANSI
> exec
| SET @x = LEFT('ANSI', 18446744073709551615)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "LONGTEXT"
< row
: bytes
| ANSI
> exec
| SET @x = RIGHT('ANSI', 9223372036854775807)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "LONGTEXT"
< row
: bytes
| ANSI
> exec
| SET @x = ABS(-4)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 4
> exec
| SET @x = ABS('-4')
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 4
> exec
| SET @x = ABS(-4.5)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 4.5
> exec
| SET @x = ABS('abc')
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 0
> exec
| SET @x = ABS('-1.5e0')
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 1.5
> exec
| SET @x = ABS(-1.5e0)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 1.5
> exec
| SET @x = ABS(' -3x')
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 3
> exec
| SET @x = ABS(4)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 4
> exec
| SET @x = ABS(-9223372036854775807)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 9223372036854775807
> exec
| SET @x = ABS(18446744073709551615)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 18446744073709551615
> exec
| SET @x = NOT 0
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 1
> exec
| SET @x = NOT 1
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 0
> exec
| SET @x = NOT 'abc'
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 1
> exec
| SET @x = NOT '1x'
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 0
> exec
| SET @x = !0
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 1
> exec
| SET @x = NOT 0.0
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 1
> exec
| SET @x = NOT 0.5
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 0
> exec
| SET @x = LEFT('ANSI', 1.6)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "LONGTEXT"
< row
: bytes
| AN
> exec
| SET @x = LEFT('ANSI', 1.5)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "LONGTEXT"
< row
: bytes
| AN
> exec
| SET @x = LEFT('ANSI', '1.6')
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "LONGTEXT"
< row
: bytes
| A
> exec
| SET @x = RIGHT('ANSI', 2.5e0)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "LONGTEXT"
< row
: bytes
| SI
> exec
| SET @x = LEFT('ANSI', 1.4)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "LONGTEXT"
< row
: bytes
| A
> exec
| SET @x = SUBSTRING('ANSI', 1.6, 1.6)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "LONGTEXT"
< row
: bytes
| NS
> exec
| SET @x = SUBSTRING('ANSI', '2.6')
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "LONGTEXT"
< row
: bytes
| NSI
> exec
| SET @x = SUBSTRING('ANSI', 2.5, 2.5e0)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "LONGTEXT"
< row
: bytes
| SI
> exec
| SET @x = LEFT('ANSI', 'x')
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "LONGTEXT"
< row
: bytes
|
> exec
| SET @x = LEFT('ANSI', -0.4)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "LONGTEXT"
< row
: bytes
|
> exec
| SET @x = -99999999999999999999
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| -99999999999999999999
> exec
| SET @x = 99999999999999999999 + 1
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 100000000000000000000
> exec
| SET @x = -(99999999999999999999)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| -99999999999999999999
> exec
| SET @x = ROUND(3.6)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 4
> exec
| SET @x = ROUND(-3.5)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| -4
> exec
| SET @x = ROUND(2.5e0)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 2
> exec
| SET @x = ROUND(3.5e0)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 4
> exec
| SET @x = ROUND(1.2345, 2)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 1.23
> exec
| SET @x = ROUND(1234.5, -2)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 1200
> exec
| SET @x = ROUND('3.6')
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 4
> exec
| SET @x = ROUND(1.2, 3)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 1.2
> exec
| SET @x = ROUND(3.14159, 2.6)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 3.142
> exec
| SET @x = ROUND(-1.25, 1)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| -1.3
> exec
| SET @x = ROUND(1.25e0, 1)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 1.2
> exec
| SET @x = ROUND(15, -1)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 20
> exec
| SET @x = ROUND(-15, -1)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| -20
> exec
| SET @x = ROUND(14, -5)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 0
> exec
| SET @x = ROUND(1.5, 40)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 1.5
> exec
| SET @x = TRUNCATE(3.69, 1)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 3.6
> exec
| SET @x = TRUNCATE(-3.69, 1)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| -3.6
> exec
| SET @x = TRUNCATE(1234, -2)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 1200
> exec
| SET @x = TRUNCATE(1.2, 3)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 1.2
> exec
| SET @x = TRUNCATE('3.69', 1)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 3.6
> exec
| SET @x = TRUNCATE(3.69e0, 1)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 3.6
> exec
| SET @x = FLOOR(-3.5)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| -4
> exec
| SET @x = FLOOR(3.5)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 3
> exec
| SET @x = CEIL(3.2)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 4
> exec
| SET @x = CEILING(-3.2)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| -3
> exec
| SET @x = FLOOR(2.5e0)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 2
> exec
| SET @x = FLOOR('3.7')
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 3
> exec
| SET @x = CEIL('x')
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 0
> exec
| SET @x = CEIL(-0.5)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 0
> exec
| SET @x = FLOOR(99999999999999999999.5)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 99999999999999999999
> exec
| SET @x = SIGN(-2.5)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| -1
> exec
| SET @x = SIGN(0)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 0
> exec
| SET @x = SIGN('x')
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 0
> exec
| SET @x = SIGN('-3')
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| -1
> exec
| SET @x = SIGN(2e0)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 1
> exec
| SET @x = MOD(5, 4)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 1
> exec
| SET @x = MOD(-5, 4)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| -1
> exec
| SET @x = MOD(5.5, 2)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 1.5
> exec
| SET @x = MOD('7', 4)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 3
> exec
| SET @x = MOD(5e0, 2)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 1
> exec
| SET @x = ~-2
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 1
> exec
| SET @x = ~0
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 18446744073709551615
> exec
| SET @x = ~'3'
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 18446744073709551612
> exec
| SET @x = ~1.6
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 18446744073709551613
> exec
| SET @x = ~1.5e0
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 18446744073709551613
> exec
| SET @x = ~'x'
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 18446744073709551615
> exec
| SET @x = BIT_COUNT(1)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 1
> exec
| SET @x = BIT_COUNT(-1)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 64
> exec
| SET @x = BIT_COUNT('3')
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 2
> exec
| SET @x = BIT_COUNT(1.6)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 1
> exec
| SET @x = BIT_COUNT(1.5e0)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 1
> exec
| SET @x = BIT_COUNT('x')
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 0
> exec
| SET @x = BIT_COUNT(18446744073709551615)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 64
> exec
| SET @x = BIT_COUNT(X'FF')
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 8
> exec
| SET @x = BIT_COUNT(b'101')
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 2
> exec
| SET @x = ASCII(X'01')
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 1
> exec
| SET @x = ASCII('A')
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 65
> exec
| SET @x = ASCII('')
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 0
> exec
| SET @x = ASCII(65)
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 54
> exec
| SET @x = 1 << 63
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 9223372036854775808
> exec
| SET @x = (1 << 63) > 0
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 1
> exec
| SET @x = -1 | 0
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 18446744073709551615
> exec
| SET @x = 1 << 64
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 0
> exec
| SET @x = 5 >> 1
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 2
> exec
| SET @x = 1.6 | 0
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 2
> exec
| SET @x = '3' ^ 1
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 2
> exec
| SET @x = -1 & 255
< ok 0 0
> query
| SELECT CAST(@x AS CHAR)
< result
< column "CAST(@x AS CHAR)" "VARCHAR"
< row
: bytes
| 255
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| SET SESSION sql_mode = IF(REPLACE('ansi', X'78', '') = 'ANSI', 'ANSI', '')
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| REAL_AS_FLOAT,PIPES_AS_CONCAT,ANSI_QUOTES,IGNORE_SPACE,ONLY_FULL_GROUP_BY,ANSI
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| SET SESSION sql_mode = IF(REPLACE('ansi', 'x', X'78') = 'ANSI', 'ANSI', '')
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| REAL_AS_FLOAT,PIPES_AS_CONCAT,ANSI_QUOTES,IGNORE_SPACE,ONLY_FULL_GROUP_BY,ANSI
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| SET SESSION sql_mode = IF(REPLACE(X'616E7369', 'x', '') = 'ANSI', 'ANSI', '')
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
|
> exec
| SET @x_max = 9223372036854775807
< ok 0 0
> exec
| SET SESSION sql_mode = 'ANSI'
< ok 0 0
> exec
| SET SESSION sql_mode = 0 AND (@x_max + 1)
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
|
> exec
| SET SESSION sql_mode = 'ANSI'
< ok 0 0
> exec
| SET SESSION sql_mode = 1 OR (@x_max + 1)
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| REAL_AS_FLOAT
> exec
| SET SESSION sql_mode = 'ANSI'
< ok 0 0
> exec
| SET SESSION sql_mode = 1 AND (@x_max + 1)
< error
| 1690 22003 BIGINT value is out of range in '((@"x_max") + 1)'
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| REAL_AS_FLOAT,PIPES_AS_CONCAT,ANSI_QUOTES,IGNORE_SPACE,ONLY_FULL_GROUP_BY,ANSI
> exec
| SET SESSION sql_mode = 'ANSI'
< ok 0 0
> exec
| SET SESSION sql_mode = NULL AND (@x_max + 1)
< error
| 1690 22003 BIGINT value is out of range in '((@"x_max") + 1)'
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| REAL_AS_FLOAT,PIPES_AS_CONCAT,ANSI_QUOTES,IGNORE_SPACE,ONLY_FULL_GROUP_BY,ANSI
> exec
| SET SESSION sql_mode = 'ANSI'
< ok 0 0
> exec
| SET SESSION sql_mode = 0 OR (@x_max + 1)
< error
| 1690 22003 BIGINT value is out of range in '((@"x_max") + 1)'
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| REAL_AS_FLOAT,PIPES_AS_CONCAT,ANSI_QUOTES,IGNORE_SPACE,ONLY_FULL_GROUP_BY,ANSI
> exec
| SET @b_cast = CAST('a' AS CHAR(4294967295))
< error
| 1690 22003 char(n) value is out of range in 'cast_as_char'
> exec
| SET @b_cast = CAST('a' AS BINARY(4294967295))
< ok 0 0
> query
| SELECT @b_cast IS NULL
< result
< column "@b_cast IS NULL" "BIGINT" "int64"
< row
: int64 1
> exec
| SET @b_cast = CAST('a' AS BINARY(67108865))
< ok 0 0
> query
| SELECT @b_cast IS NULL
< result
< column "@b_cast IS NULL" "BIGINT" "int64"
< row
: int64 1
> exec
| SET @b_cast = CAST('a' AS CHAR(67108865))
< ok 0 0
> query
| SELECT @b_cast IS NULL
< result
< column "@b_cast IS NULL" "BIGINT" "int64"
< row
: int64 1
> exec
| SET @b_cast = CAST('a' AS CHAR(67108864))
< ok 0 0
> query
| SELECT @b_cast IS NULL
< result
< column "@b_cast IS NULL" "BIGINT" "int64"
< row
: int64 0
> exec
| SET @b_cast = CAST('a' AS CHAR(4294967295) CHARACTER SET latin1)
< ok 0 0
> query
| SELECT @b_cast IS NULL
< result
< column "@b_cast IS NULL" "BIGINT" "int64"
< row
: int64 1
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| SET SESSION sql_mode = FLOOR(3.5)
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| REAL_AS_FLOAT,PIPES_AS_CONCAT
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| SET SESSION sql_mode = CEIL(3.2)
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| ANSI_QUOTES
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| SET SESSION sql_mode = ROUND(3.6)
< error
| 1232 42000 Incorrect argument type to variable 'sql_mode'
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
|
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| SET SESSION sql_mode = ROUND(4)
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| ANSI_QUOTES
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| SET SESSION sql_mode = ROUND(3.6e0)
< error
| 1232 42000 Incorrect argument type to variable 'sql_mode'
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
|
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| SET SESSION sql_mode = TRUNCATE(3.6, 0)
< error
| 1232 42000 Incorrect argument type to variable 'sql_mode'
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
|
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| SET SESSION sql_mode = TRUNCATE(4, 1)
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| ANSI_QUOTES
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| SET SESSION sql_mode = ROUND(14, -1)
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| PIPES_AS_CONCAT,IGNORE_SPACE
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| SET SESSION sql_mode = SIGN(-2.5)
< error
| 1231 42000 Variable 'sql_mode' can't be set to the value of '-1'
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
|
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| SET SESSION sql_mode = SIGN(2.5)
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| REAL_AS_FLOAT
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| SET SESSION sql_mode = FLOOR('3.7')
< error
| 1232 42000 Incorrect argument type to variable 'sql_mode'
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
|
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| SET SESSION sql_mode = ABS(-4.0)
< error
| 1232 42000 Incorrect argument type to variable 'sql_mode'
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
|
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| SET SESSION sql_mode = FLOOR(99999999999999999999.5)
< error
| 1232 42000 Incorrect argument type to variable 'sql_mode'
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
|
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| SET SESSION sql_mode = ROUND('4')
< error
| 1232 42000 Incorrect argument type to variable 'sql_mode'
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
|
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| SET SESSION sql_mode = MOD(5, 4)
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| REAL_AS_FLOAT
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| SET SESSION sql_mode = MOD(5.5, 2)
< error
| 1232 42000 Incorrect argument type to variable 'sql_mode'
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
|
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| SET SESSION sql_mode = ~-2
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| REAL_AS_FLOAT
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| SET SESSION sql_mode = ~0
< error
| 1231 42000 Variable 'sql_mode' can't be set to the value of '18446744073709551615'
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
|
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| SET SESSION sql_mode = BIT_COUNT(1)
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| REAL_AS_FLOAT
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| SET SESSION sql_mode = ASCII(X'01')
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| REAL_AS_FLOAT
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| SET SESSION sql_mode = ((1 << 63) > 0)
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| REAL_AS_FLOAT
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| SET SESSION sql_mode = IF(IF(1, REPLACE('ansi', X'78', ''), '') = 'ANSI', 'ANSI', '')
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| REAL_AS_FLOAT,PIPES_AS_CONCAT,ANSI_QUOTES,IGNORE_SPACE,ONLY_FULL_GROUP_BY,ANSI
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| SET @m_round = ROUND(3.6)
< ok 0 0
> exec
| SET SESSION sql_mode = @m_round
< error
| 1232 42000 Incorrect argument type to variable 'sql_mode'
> exec
| SET foreign_key_checks = 0, @x_abs = ABS(-9223372036854775808)
< error
| 1690 22003 BIGINT value is out of range in 'abs(-(9223372036854775808))'
> query
| SELECT CAST(@@foreign_key_checks AS CHAR)
< result
< column "CAST(@@foreign_key_checks AS CHAR)" "VARCHAR"
< row
: bytes
| 1
> exec
| SET foreign_key_checks = 0, @x_scope = @@SESSION.max_connections
< error
| 1238 HY000 Variable 'max_connections' is a GLOBAL variable
> query
| SELECT CAST(@@foreign_key_checks AS CHAR)
< result
< column "CAST(@@foreign_key_checks AS CHAR)" "VARCHAR"
< row
: bytes
| 1
> exec
| SET foreign_key_checks = 0, @x_scope = @@GLOBAL.warning_count
< error
| 1238 HY000 Variable 'warning_count' is a SESSION variable
> query
| SELECT CAST(@@foreign_key_checks AS CHAR)
< result
< column "CAST(@@foreign_key_checks AS CHAR)" "VARCHAR"
< row
: bytes
| 1
> exec
| SET @x_scope = @@GLOBAL.max_connections
< ok 0 0
> exec
| SET @x_scope = @@SESSION.warning_count
< ok 0 0
> exec
| SET @x_scope = @@max_connections
< ok 0 0
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| SET SESSION sql_mode = NOT 0
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| REAL_AS_FLOAT
> exec
| SET SESSION sql_mode = ''
< ok 0 0
> exec
| SET @m_abs = ABS(-4)
< ok 0 0
> exec
| SET SESSION sql_mode = @m_abs
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| ANSI_QUOTES
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
> exec
| SET SESSION time_zone = '+00:00'
< ok 0 0
> exec
| CREATE TABLE z (ts TIMESTAMP DEFAULT _utf8mb4'2026-01-01 12:00:00', tn TIMESTAMP DEFAULT N'2026-01-01 12:00:00')
< ok 0 0
> exec
| SET SESSION time_zone = '+03:00'
< ok 0 0
> query
| SHOW CREATE TABLE z
< result
< column "Table" "VARCHAR" "string"
< column "Create Table" "VARCHAR" "string"
< row
: bytes
| z
: bytes
| CREATE TABLE `z` (
|   `ts` timestamp NULL DEFAULT '2026-01-01 15:00:00',
|   `tn` timestamp NULL DEFAULT '2026-01-01 15:00:00'
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ENCRYPTION='Y'
> exec
| SET SESSION time_zone = DEFAULT
< ok 0 0
> exec
| SET SESSION sql_mode = CAST(3.6 AS SIGNED)
< ok 0 0
> query
| SELECT @@sql_mode
< result
< column "@@sql_mode" "VARCHAR"
< row
: bytes
| ANSI_QUOTES
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| SET SESSION time_zone = DEFAULT
< ok 0 0
> exec
| ALTER TABLESPACE mysql ENCRYPTION='N'
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| USE `db_cs_bin`
< ok 0 0
> exec
| DROP DATABASE IF EXISTS d_qual_own
< ok 2 0
> exec
| USE `db_cs_bin`
< ok 0 0
> exec
| DROP DATABASE IF EXISTS d_qfn
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| SET SESSION foreign_key_checks = 1
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| USE `test`
< ok 0 0
> exec
| DROP DATABASE IF EXISTS db_cs_bin
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| DROP TABLE IF EXISTS a_def_omit_y
< ok 0 0
> exec
| DROP TABLESPACE ts_omit_enc
< ok 0 0
> exec
| SET SESSION default_table_encryption = OFF
< ok 0 0
> exec
| DROP TABLE IF EXISTS a_def_enc
< ok 0 0
> exec
| DROP TABLESPACE ts_def_enc
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLESPACE ts_enc_chain
< ok 0 0
> exec
| UNINSTALL PLUGIN keyring_file
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLESPACE ts_sdl_plain
< ok 0 0
> exec
| SET time_zone = DEFAULT
< ok 0 0
> exec
| SET NAMES utf8mb4
< ok 0 0
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
