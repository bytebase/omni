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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
> exec
| CREATE TABLE a (b INT PRIMARY KEY, c INT, d INT, KEY idx_d (d))
< ok 0 0
> exec
| ALTER TABLE a ADD x ENUM('ff', 'ﬀ') COLLATE utf8mb4_0900_ai_ci
< error
| 1291 HY000 Column 'x' has duplicated value 'ff' in ENUM
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
> exec ^
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
> exec ^
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
> exec ^
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
> exec
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
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
| SET SESSION sql_mode = DEFAULT
< ok 0 0
> exec
| DROP VIEW IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS a
< ok 0 0
> exec
| DROP TABLE IF EXISTS p
< ok 0 0
> exec
| DROP TABLE IF EXISTS s
< ok 0 0
> exec
| DROP TABLE IF EXISTS z
< ok 0 0
> exec
| DROP VIEW IF EXISTS v
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = OFF
< ok 0 0
> exec
| DROP DATABASE IF EXISTS cteother
< ok 0 0
> exec
| DROP TABLE IF EXISTS `dual`
< ok 0 0
> exec
| DROP TABLE IF EXISTS c2
< ok 0 0
> exec
| DROP EVENT IF EXISTS e
< ok 0 0
