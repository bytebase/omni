-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_DiffMultiTableRoundTrip default#1
> exec
| DROP DATABASE IF EXISTS diff_multi
< ok 0 0
> exec
| CREATE DATABASE diff_multi DEFAULT CHARSET=utf8mb4
< ok 1 0
> exec
| USE diff_multi
< ok 0 0
> exec
| CREATE TABLE parent (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY, name VARCHAR(50) NOT NULL)
< ok 0 0
> exec
| CREATE TABLE child (id INT NOT NULL PRIMARY KEY, pid INT, amount DECIMAL(10,2) DEFAULT 0, KEY pid_idx (pid), CONSTRAINT fk_pid FOREIGN KEY (pid) REFERENCES parent (id))
< ok 0 0
> query
| SHOW CREATE TABLE diff_multi.parent
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| parent
: bytes
| CREATE TABLE `parent` (
|   `id` int NOT NULL AUTO_INCREMENT,
|   `name` varchar(50) NOT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
> query
| SHOW CREATE TABLE diff_multi.child
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| child
: bytes
| CREATE TABLE `child` (
|   `id` int NOT NULL,
|   `pid` int DEFAULT NULL,
|   `amount` decimal(10,2) DEFAULT '0.00',
|   PRIMARY KEY (`id`),
|   KEY `pid_idx` (`pid`),
|   CONSTRAINT `fk_pid` FOREIGN KEY (`pid`) REFERENCES `parent` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
