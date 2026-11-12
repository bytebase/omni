-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestContainerInvalidDefaultForms default#1
> exec
| DROP DATABASE IF EXISTS testcontainerinvaliddefaultforms
< ok 0 0
> exec
| CREATE DATABASE testcontainerinvaliddefaultforms
< ok 0 0
> exec
| USE testcontainerinvaliddefaultforms
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a date DEFAULT '')
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a datetime DEFAULT '')
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a datetime(3) DEFAULT '')
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a timestamp DEFAULT '')
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a year DEFAULT '')
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a int DEFAULT '')
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a tinyint DEFAULT '')
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a bigint DEFAULT '')
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a decimal(10,2) DEFAULT '')
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a double DEFAULT '')
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a bool DEFAULT '')
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a enum('x','y') DEFAULT '')
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a enum('','x') DEFAULT '')
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a enum('a','') DEFAULT '')
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a enum('','x') NOT NULL DEFAULT '')
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a enum(' ','x') DEFAULT ' ')
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a enum('x') DEFAULT ' ')
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a date DEFAULT ' ')
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a date NOT NULL DEFAULT '')
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a varchar(10) DEFAULT '')
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a char(5) DEFAULT '')
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a set('x','y') DEFAULT '')
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a binary(4) DEFAULT '')
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a bit(4) DEFAULT '')
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a date DEFAULT '2020-01-01')
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a int DEFAULT '5')
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a decimal DEFAULT '1.5')
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a enum('x') DEFAULT 'x')
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL)
< ok 0 0
> exec
| ALTER TABLE t ALTER COLUMN a SET DEFAULT NULL
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE t ALTER COLUMN a SET DEFAULT NULL
< error
| 1171 42000 All parts of a PRIMARY KEY must be NOT NULL; if you need NULL in a key, use UNIQUE instead
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (a,b))
< ok 0 0
> exec
| ALTER TABLE t ALTER COLUMN b SET DEFAULT NULL
< error
| 1171 42000 All parts of a PRIMARY KEY must be NOT NULL; if you need NULL in a key, use UNIQUE instead
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL UNIQUE)
< ok 0 0
> exec
| ALTER TABLE t ALTER COLUMN a SET DEFAULT NULL
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (x date)
< ok 0 0
> exec
| ALTER TABLE t ALTER COLUMN x SET DEFAULT ''
< error
| 1067 42000 Invalid default value for 'x'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (x int)
< ok 0 0
> exec
| ALTER TABLE t ALTER COLUMN x SET DEFAULT ''
< error
| 1067 42000 Invalid default value for 'x'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT)
< ok 0 0
> exec
| ALTER TABLE t ALTER COLUMN a SET DEFAULT NULL
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL)
< ok 0 0
> exec
| ALTER TABLE t ALTER COLUMN a SET DEFAULT 5
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (x varchar(5))
< ok 0 0
> exec
| ALTER TABLE t ALTER COLUMN x SET DEFAULT ''
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a enum('','x'))
< ok 0 0
> exec
| ALTER TABLE t ALTER COLUMN a SET DEFAULT ''
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a enum('x','y'))
< ok 0 0
> exec
| ALTER TABLE t ALTER COLUMN a SET DEFAULT ''
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL DEFAULT NULL)
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT PRIMARY KEY DEFAULT NULL)
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT DEFAULT NULL NOT NULL)
< error
| 1067 42000 Invalid default value for 'a'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT DEFAULT NULL, PRIMARY KEY(a))
< error
| 1171 42000 All parts of a PRIMARY KEY must be NOT NULL; if you need NULL in a key, use UNIQUE instead
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL DEFAULT NULL, PRIMARY KEY(a))
< error
| 1171 42000 All parts of a PRIMARY KEY must be NOT NULL; if you need NULL in a key, use UNIQUE instead
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT auto_increment DEFAULT NULL, PRIMARY KEY(a))
< error
| 1171 42000 All parts of a PRIMARY KEY must be NOT NULL; if you need NULL in a key, use UNIQUE instead
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT DEFAULT NULL)
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT DEFAULT NULL UNIQUE)
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL NULL DEFAULT NULL)
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (x int)
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN b INT NOT NULL DEFAULT NULL
< error
| 1067 42000 Invalid default value for 'b'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (x int)
< ok 0 0
> exec
| ALTER TABLE t MODIFY COLUMN x INT NOT NULL DEFAULT NULL
< error
| 1067 42000 Invalid default value for 'x'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (x int)
< ok 0 0
> exec
| ALTER TABLE t CHANGE COLUMN x y INT NOT NULL DEFAULT NULL
< error
| 1067 42000 Invalid default value for 'y'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (x int)
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN c date DEFAULT ''
< error
| 1067 42000 Invalid default value for 'c'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (x int)
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN b INT DEFAULT NULL
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (x int)
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN a enum('','x') DEFAULT ''
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (x int)
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN c enum('x','y') DEFAULT ''
< error
| 1067 42000 Invalid default value for 'c'
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a int auto_increment, b int auto_increment)
< error
| 1075 42000 Incorrect table definition; there can be only one auto column and it must be defined as a key
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a int auto_increment primary key, b int auto_increment)
< error
| 1075 42000 Incorrect table definition; there can be only one auto column and it must be defined as a key
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a int auto_increment, b int auto_increment, primary key(a))
< error
| 1075 42000 Incorrect table definition; there can be only one auto column and it must be defined as a key
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a int auto_increment unique, b int auto_increment unique)
< error
| 1075 42000 Incorrect table definition; there can be only one auto column and it must be defined as a key
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a SERIAL, b SERIAL)
< error
| 1075 42000 Incorrect table definition; there can be only one auto column and it must be defined as a key
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a SERIAL, b int auto_increment)
< error
| 1075 42000 Incorrect table definition; there can be only one auto column and it must be defined as a key
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a int auto_increment primary key)
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a int auto_increment)
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a SERIAL)
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a bigint auto_random primary key, b int auto_increment)
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (id int)
< ok 0 0
> exec
| CREATE INDEX idx_c ON t (c)
< error
| 1072 42000 column does not exist: c
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (id int)
< ok 0 0
> exec
| CREATE UNIQUE INDEX uidx ON t (c)
< error
| 1072 42000 column does not exist: c
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (id int, INDEX idx_c (c))
< error
| 1072 42000 column does not exist: c
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (id int, UNIQUE KEY uk_c (c))
< error
| 1072 42000 column does not exist: c
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (id int, PRIMARY KEY (c))
< error
| 1072 42000 Key column 'c' doesn't exist in table
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (id int, UNIQUE INDEX ui (id, c))
< error
| 1072 42000 column does not exist: c
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (id varchar(10), INDEX idx_p (c(3)))
< error
| 1072 42000 column does not exist: c
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (id int)
< ok 0 0
> exec
| ALTER TABLE t ADD INDEX idx_c (c)
< error
| 1072 42000 column does not exist: c
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (id int)
< ok 0 0
> exec
| ALTER TABLE t ADD UNIQUE KEY uk (c)
< error
| 1072 42000 column does not exist: c
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (id int)
< ok 0 0
> exec
| ALTER TABLE t ADD PRIMARY KEY (c)
< error
| 1072 42000 column does not exist: c
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (id int, KEY idx_x (id))
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a int, INDEX idx (A))
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a int, b int as (a+1), INDEX idx (b))
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (id int, INDEX idx_e ((id+1)))
< ok 0 0
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (id int, FULLTEXT KEY ft_c (c))
< ok 0 0
> exec
| DROP DATABASE IF EXISTS testcontainerinvaliddefaultforms
< ok 0 0
