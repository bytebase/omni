-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestContainer_Section_1_17_ForeignKeys default#1
> exec
| DROP TABLE IF EXISTS t_fk1
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_parent1
< ok 0 0
> exec
| CREATE TABLE t_parent1 (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE t_fk1 (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_parent1(id))
< ok 0 0
> query
| SHOW CREATE TABLE t_fk1
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_fk1
: bytes
| CREATE TABLE `t_fk1` (
|   `id` int DEFAULT NULL,
|   `pid` int DEFAULT NULL,
|   KEY `pid` (`pid`),
|   CONSTRAINT `t_fk1_ibfk_1` FOREIGN KEY (`pid`) REFERENCES `t_parent1` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_fk2
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_parent2
< ok 0 0
> exec
| CREATE TABLE t_parent2 (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE t_fk2 (id INT, pid INT, CONSTRAINT `fk_name` FOREIGN KEY (pid) REFERENCES t_parent2(id))
< ok 0 0
> query
| SHOW CREATE TABLE t_fk2
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_fk2
: bytes
| CREATE TABLE `t_fk2` (
|   `id` int DEFAULT NULL,
|   `pid` int DEFAULT NULL,
|   KEY `fk_name` (`pid`),
|   CONSTRAINT `fk_name` FOREIGN KEY (`pid`) REFERENCES `t_parent2` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_fk3
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_parent3
< ok 0 0
> exec
| CREATE TABLE t_parent3 (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE t_fk3 (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_parent3(id) ON DELETE CASCADE)
< ok 0 0
> query
| SHOW CREATE TABLE t_fk3
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_fk3
: bytes
| CREATE TABLE `t_fk3` (
|   `id` int DEFAULT NULL,
|   `pid` int DEFAULT NULL,
|   KEY `pid` (`pid`),
|   CONSTRAINT `t_fk3_ibfk_1` FOREIGN KEY (`pid`) REFERENCES `t_parent3` (`id`) ON DELETE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_fk4
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_parent4
< ok 0 0
> exec
| CREATE TABLE t_parent4 (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE t_fk4 (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_parent4(id) ON DELETE SET NULL)
< ok 0 0
> query
| SHOW CREATE TABLE t_fk4
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_fk4
: bytes
| CREATE TABLE `t_fk4` (
|   `id` int DEFAULT NULL,
|   `pid` int DEFAULT NULL,
|   KEY `pid` (`pid`),
|   CONSTRAINT `t_fk4_ibfk_1` FOREIGN KEY (`pid`) REFERENCES `t_parent4` (`id`) ON DELETE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_fk5
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_parent5
< ok 0 0
> exec
| CREATE TABLE t_parent5 (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE t_fk5 (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_parent5(id) ON DELETE SET DEFAULT)
< ok 0 0
> query
| SHOW CREATE TABLE t_fk5
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_fk5
: bytes
| CREATE TABLE `t_fk5` (
|   `id` int DEFAULT NULL,
|   `pid` int DEFAULT NULL,
|   KEY `pid` (`pid`),
|   CONSTRAINT `t_fk5_ibfk_1` FOREIGN KEY (`pid`) REFERENCES `t_parent5` (`id`) ON DELETE SET DEFAULT
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_fk6
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_parent6
< ok 0 0
> exec
| CREATE TABLE t_parent6 (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE t_fk6 (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_parent6(id) ON DELETE RESTRICT)
< ok 0 0
> query
| SHOW CREATE TABLE t_fk6
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_fk6
: bytes
| CREATE TABLE `t_fk6` (
|   `id` int DEFAULT NULL,
|   `pid` int DEFAULT NULL,
|   KEY `pid` (`pid`),
|   CONSTRAINT `t_fk6_ibfk_1` FOREIGN KEY (`pid`) REFERENCES `t_parent6` (`id`) ON DELETE RESTRICT
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_fk7
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_parent7
< ok 0 0
> exec
| CREATE TABLE t_parent7 (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE t_fk7 (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_parent7(id) ON DELETE NO ACTION)
< ok 0 0
> query
| SHOW CREATE TABLE t_fk7
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_fk7
: bytes
| CREATE TABLE `t_fk7` (
|   `id` int DEFAULT NULL,
|   `pid` int DEFAULT NULL,
|   KEY `pid` (`pid`),
|   CONSTRAINT `t_fk7_ibfk_1` FOREIGN KEY (`pid`) REFERENCES `t_parent7` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_fk8
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_parent8
< ok 0 0
> exec
| CREATE TABLE t_parent8 (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE t_fk8 (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_parent8(id) ON UPDATE CASCADE)
< ok 0 0
> query
| SHOW CREATE TABLE t_fk8
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_fk8
: bytes
| CREATE TABLE `t_fk8` (
|   `id` int DEFAULT NULL,
|   `pid` int DEFAULT NULL,
|   KEY `pid` (`pid`),
|   CONSTRAINT `t_fk8_ibfk_1` FOREIGN KEY (`pid`) REFERENCES `t_parent8` (`id`) ON UPDATE CASCADE
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_fk9
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_parent9
< ok 0 0
> exec
| CREATE TABLE t_parent9 (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE t_fk9 (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_parent9(id) ON UPDATE SET NULL)
< ok 0 0
> query
| SHOW CREATE TABLE t_fk9
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_fk9
: bytes
| CREATE TABLE `t_fk9` (
|   `id` int DEFAULT NULL,
|   `pid` int DEFAULT NULL,
|   KEY `pid` (`pid`),
|   CONSTRAINT `t_fk9_ibfk_1` FOREIGN KEY (`pid`) REFERENCES `t_parent9` (`id`) ON UPDATE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_fk10
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_parent10
< ok 0 0
> exec
| CREATE TABLE t_parent10 (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE t_fk10 (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_parent10(id) ON DELETE CASCADE ON UPDATE SET NULL)
< ok 0 0
> query
| SHOW CREATE TABLE t_fk10
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_fk10
: bytes
| CREATE TABLE `t_fk10` (
|   `id` int DEFAULT NULL,
|   `pid` int DEFAULT NULL,
|   KEY `pid` (`pid`),
|   CONSTRAINT `t_fk10_ibfk_1` FOREIGN KEY (`pid`) REFERENCES `t_parent10` (`id`) ON DELETE CASCADE ON UPDATE SET NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_fk11
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_parent11
< ok 0 0
> exec
| CREATE TABLE t_parent11 (id INT NOT NULL, id2 INT NOT NULL, PRIMARY KEY (id), KEY (id2)) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE t_fk11 (id INT, pid INT, pid2 INT, FOREIGN KEY (pid) REFERENCES t_parent11(id), FOREIGN KEY (pid2) REFERENCES t_parent11(id2))
< ok 0 0
> query
| SHOW CREATE TABLE t_fk11
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_fk11
: bytes
| CREATE TABLE `t_fk11` (
|   `id` int DEFAULT NULL,
|   `pid` int DEFAULT NULL,
|   `pid2` int DEFAULT NULL,
|   KEY `pid` (`pid`),
|   KEY `pid2` (`pid2`),
|   CONSTRAINT `t_fk11_ibfk_1` FOREIGN KEY (`pid`) REFERENCES `t_parent11` (`id`),
|   CONSTRAINT `t_fk11_ibfk_2` FOREIGN KEY (`pid2`) REFERENCES `t_parent11` (`id2`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_fk12
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_parent12
< ok 0 0
> exec
| CREATE TABLE t_parent12 (id INT NOT NULL, PRIMARY KEY (id)) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE t_fk12 (id INT, pid INT, FOREIGN KEY (pid) REFERENCES t_parent12(id))
< ok 0 0
> query
| SHOW CREATE TABLE t_fk12
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_fk12
: bytes
| CREATE TABLE `t_fk12` (
|   `id` int DEFAULT NULL,
|   `pid` int DEFAULT NULL,
|   KEY `pid` (`pid`),
|   CONSTRAINT `t_fk12_ibfk_1` FOREIGN KEY (`pid`) REFERENCES `t_parent12` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_fk13
< ok 0 0
> exec
| CREATE TABLE t_fk13 (id INT NOT NULL, parent_id INT, PRIMARY KEY (id), FOREIGN KEY (parent_id) REFERENCES t_fk13(id))
< ok 0 0
> query
| SHOW CREATE TABLE t_fk13
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_fk13
: bytes
| CREATE TABLE `t_fk13` (
|   `id` int NOT NULL,
|   `parent_id` int DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `parent_id` (`parent_id`),
|   CONSTRAINT `t_fk13_ibfk_1` FOREIGN KEY (`parent_id`) REFERENCES `t_fk13` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> exec
| DROP TABLE IF EXISTS t_fk14
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_parent14
< ok 0 0
> exec
| CREATE TABLE t_parent14 (x INT NOT NULL, y INT NOT NULL, PRIMARY KEY (x, y)) ENGINE=InnoDB
< ok 0 0
> exec
| CREATE TABLE t_fk14 (id INT, a INT, b INT, FOREIGN KEY (a, b) REFERENCES t_parent14(x, y))
< ok 0 0
> query
| SHOW CREATE TABLE t_fk14
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| t_fk14
: bytes
| CREATE TABLE `t_fk14` (
|   `id` int DEFAULT NULL,
|   `a` int DEFAULT NULL,
|   `b` int DEFAULT NULL,
|   KEY `a` (`a`,`b`),
|   CONSTRAINT `t_fk14_ibfk_1` FOREIGN KEY (`a`, `b`) REFERENCES `t_parent14` (`x`, `y`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
