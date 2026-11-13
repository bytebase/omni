-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestAddIfNotExistsOracle default#1
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD INDEX IF NOT EXISTS k (a)
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
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `a` (`a`),
|   KEY `k` (`b`),
|   CONSTRAINT `ck` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD INDEX IF NOT EXISTS n (a)
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
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `a` (`a`),
|   KEY `k` (`b`),
|   KEY `n` (`a`),
|   CONSTRAINT `ck` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD INDEX IF NOT EXISTS (b)
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
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `a` (`a`),
|   KEY `k` (`b`),
|   KEY `b` (`b`),
|   CONSTRAINT `ck` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE IF NOT EXISTS (a)
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
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `a` (`a`),
|   KEY `k` (`b`),
|   CONSTRAINT `ck` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE IF NOT EXISTS k (a)
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
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `a` (`a`),
|   KEY `k` (`b`),
|   CONSTRAINT `ck` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD CONSTRAINT nc UNIQUE IF NOT EXISTS k (a)
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
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `a` (`a`),
|   KEY `k` (`b`),
|   CONSTRAINT `ck` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD CONSTRAINT nc UNIQUE IF NOT EXISTS (a)
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
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `nc` (`a`),
|   KEY `a` (`a`),
|   KEY `k` (`b`),
|   CONSTRAINT `ck` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD INDEX IF NOT EXISTS k (nosuch)
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
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `a` (`a`),
|   KEY `k` (`b`),
|   CONSTRAINT `ck` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD INDEX IF NOT EXISTS n (a), ADD INDEX IF NOT EXISTS n (b)
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
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `a` (`a`),
|   KEY `k` (`b`),
|   KEY `n` (`a`),
|   CONSTRAINT `ck` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD INDEX n (a), ADD INDEX n (b)
< error
| 1061 42000 Duplicate key name 'n'
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD (INDEX IF NOT EXISTS n (a), INDEX IF NOT EXISTS n (b))
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
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `a` (`a`),
|   KEY `k` (`b`),
|   KEY `n` (`a`),
|   CONSTRAINT `ck` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE np ADD PRIMARY KEY IF NOT EXISTS (id), ADD PRIMARY KEY IF NOT EXISTS (id)
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD PRIMARY KEY IF NOT EXISTS (a), ADD PRIMARY KEY IF NOT EXISTS (b)
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
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `a` (`a`),
|   KEY `k` (`b`),
|   CONSTRAINT `ck` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD INDEX IF NOT EXISTS n (a), ADD INDEX n (b)
< error
| 1061 42000 Duplicate key name 'n'
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD INDEX n (a), ADD INDEX IF NOT EXISTS n (b)
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
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `a` (`a`),
|   KEY `k` (`b`),
|   KEY `n` (`a`),
|   CONSTRAINT `ck` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS c1 CHECK (a > 0), ADD CONSTRAINT IF NOT EXISTS c1 CHECK (a > 1)
< error
| 1826 HY000 Duplicate CHECK constraint name 'c1'
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS ck CHECK (a > 1), ADD CONSTRAINT IF NOT EXISTS ck CHECK (a > 2)
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
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `a` (`a`),
|   KEY `k` (`b`),
|   CONSTRAINT `ck` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE f ADD FOREIGN KEY IF NOT EXISTS fq (a) REFERENCES r (id), ADD FOREIGN KEY IF NOT EXISTS fq (a) REFERENCES r (id)
< ok 0 0
> query
| SHOW CREATE TABLE f
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| f
: bytes
| CREATE TABLE `f` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fq` (`a`),
|   CONSTRAINT `fkx` FOREIGN KEY (`a`) REFERENCES `r` (`id`),
|   CONSTRAINT `fq` FOREIGN KEY (`a`) REFERENCES `r` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD PRIMARY KEY IF NOT EXISTS (a)
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
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `a` (`a`),
|   KEY `k` (`b`),
|   CONSTRAINT `ck` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE np ADD PRIMARY KEY IF NOT EXISTS (id)
< ok 0 0
> query
| SHOW CREATE TABLE np
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| np
: bytes
| CREATE TABLE `np` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD PRIMARY KEY (a)
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS ck CHECK (a > 1)
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
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `a` (`a`),
|   KEY `k` (`b`),
|   CONSTRAINT `ck` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS ck2 CHECK (a > 1)
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
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `a` (`a`),
|   KEY `k` (`b`),
|   CONSTRAINT `ck` CHECK (`a` > 0),
|   CONSTRAINT `ck2` CHECK (`a` > 1)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE f ADD FOREIGN KEY IF NOT EXISTS fkx (a) REFERENCES r (id)
< ok 0 0
> query
| SHOW CREATE TABLE f
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| f
: bytes
| CREATE TABLE `f` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fkx` (`a`),
|   CONSTRAINT `fkx` FOREIGN KEY (`a`) REFERENCES `r` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE f ADD CONSTRAINT fkx FOREIGN KEY IF NOT EXISTS (a) REFERENCES r (id)
< ok 0 0
> query
| SHOW CREATE TABLE f
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| f
: bytes
| CREATE TABLE `f` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fkx` (`a`),
|   CONSTRAINT `fkx` FOREIGN KEY (`a`) REFERENCES `r` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE f ADD FOREIGN KEY IF NOT EXISTS (a) REFERENCES r (id)
< ok 0 0
> query
| SHOW CREATE TABLE f
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| f
: bytes
| CREATE TABLE `f` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `a` (`a`),
|   CONSTRAINT `f_ibfk_1` FOREIGN KEY (`a`) REFERENCES `r` (`id`),
|   CONSTRAINT `fkx` FOREIGN KEY (`a`) REFERENCES `r` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE f ADD FOREIGN KEY IF NOT EXISTS fkq (a) REFERENCES r (id)
< ok 0 0
> query
| SHOW CREATE TABLE f
< result
< column "Table" "VARCHAR"
< column "Create Table" "VARCHAR"
< row
: bytes
| f
: bytes
| CREATE TABLE `f` (
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `fkq` (`a`),
|   CONSTRAINT `fkq` FOREIGN KEY (`a`) REFERENCES `r` (`id`),
|   CONSTRAINT `fkx` FOREIGN KEY (`a`) REFERENCES `r` (`id`)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE np ADD PRIMARY KEY (id), ADD PRIMARY KEY IF NOT EXISTS (a)
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE np ADD PRIMARY KEY IF NOT EXISTS (a), ADD PRIMARY KEY (id)
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD CONSTRAINT c1 CHECK (a > 0), ADD CONSTRAINT IF NOT EXISTS c1 CHECK (a > 1)
< error
| 1826 HY000 Duplicate CHECK constraint name 'c1'
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS c1 CHECK (a > 0), ADD CONSTRAINT c1 CHECK (a > 1)
< error
| 1826 HY000 Duplicate CHECK constraint name 'c1'
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE np ADD b INT PRIMARY KEY, ADD PRIMARY KEY IF NOT EXISTS (id)
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE np ADD PRIMARY KEY IF NOT EXISTS (id), ADD b INT PRIMARY KEY
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE np ADD COLUMN (b INT PRIMARY KEY), ADD PRIMARY KEY IF NOT EXISTS (id)
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE np MODIFY id INT NOT NULL PRIMARY KEY, ADD PRIMARY KEY IF NOT EXISTS (id)
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE np CHANGE id id INT NOT NULL PRIMARY KEY, ADD PRIMARY KEY IF NOT EXISTS (id)
< error
| 1068 42000 Multiple primary key defined
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE np ADD INDEX IF NOT EXISTS `PRIMARY` (a)
< error
| 1280 42000 Incorrect index name 'PRIMARY'
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE np ADD UNIQUE IF NOT EXISTS `primary` (a)
< error
| 1280 42000 Incorrect index name 'primary'
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD INDEX IF NOT EXISTS `PRIMARY` (a)
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
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `a` (`a`),
|   KEY `k` (`b`),
|   CONSTRAINT `ck` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS c1 CHECK (a > 0), ADD INDEX c1 (a)
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
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `a` (`a`),
|   KEY `k` (`b`),
|   KEY `c1` (`a`),
|   CONSTRAINT `ck` CHECK (`a` > 0),
|   CONSTRAINT `c1` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD INDEX c2 (a), ADD CONSTRAINT IF NOT EXISTS c2 CHECK (a > 0)
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
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `a` (`a`),
|   KEY `k` (`b`),
|   KEY `c2` (`a`),
|   CONSTRAINT `ck` CHECK (`a` > 0),
|   CONSTRAINT `c2` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD INDEX n (a), ADD UNIQUE IF NOT EXISTS n (a)
< error
| 1061 42000 Duplicate key name 'n'
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE n (a), ADD INDEX IF NOT EXISTS n (a)
< error
| 1061 42000 Duplicate key name 'n'
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD INDEX IF NOT EXISTS n (a), ADD UNIQUE IF NOT EXISTS n (a)
< error
| 1061 42000 Duplicate key name 'n'
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE IF NOT EXISTS n (a), ADD UNIQUE IF NOT EXISTS n (b)
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
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   UNIQUE KEY `n` (`a`),
|   KEY `a` (`a`),
|   KEY `k` (`b`),
|   CONSTRAINT `ck` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD INDEX n (a), ADD KEY IF NOT EXISTS n (b)
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
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `a` (`a`),
|   KEY `k` (`b`),
|   KEY `n` (`a`),
|   CONSTRAINT `ck` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE IF NOT EXISTS k (a), ADD FULLTEXT IF NOT EXISTS k (a)
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
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `a` (`a`),
|   KEY `k` (`b`),
|   CONSTRAINT `ck` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD FOREIGN KEY IF NOT EXISTS k (b) REFERENCES r (id)
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
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `a` (`a`),
|   KEY `k` (`b`),
|   CONSTRAINT `k` FOREIGN KEY (`b`) REFERENCES `r` (`id`),
|   CONSTRAINT `ck` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
> exec
| CREATE TABLE r (id INT PRIMARY KEY)
< ok 0 0
> exec
| CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0))
< ok 0 0
> exec
| CREATE TABLE np (id INT NOT NULL, a INT)
< ok 0 0
> exec
| CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id))
< ok 0 0
> exec
| ALTER TABLE t ADD CONSTRAINT k FOREIGN KEY IF NOT EXISTS (b) REFERENCES r (id)
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
|   `id` int(11) NOT NULL,
|   `a` int(11) DEFAULT NULL,
|   `b` int(11) DEFAULT NULL,
|   PRIMARY KEY (`id`),
|   KEY `a` (`a`),
|   KEY `k` (`b`),
|   CONSTRAINT `k` FOREIGN KEY (`b`) REFERENCES `r` (`id`),
|   CONSTRAINT `ck` CHECK (`a` > 0)
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci
> exec
| DROP TABLE IF EXISTS f
< ok 0 0
> exec
| DROP TABLE IF EXISTS np
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| DROP TABLE IF EXISTS r
< ok 0 0
