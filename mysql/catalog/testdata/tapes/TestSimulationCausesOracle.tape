-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestSimulationCausesOracle default#1
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
