-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestContainer_InvisiblePrimaryKey default#1
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT, PRIMARY KEY (a) INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT, CONSTRAINT pk PRIMARY KEY (a) COMMENT 'x' INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT, PRIMARY KEY (a) INVISIBLE VISIBLE)
< ok 0 0
> query
| SELECT DISTINCT INDEX_NAME, IS_VISIBLE FROM information_schema.STATISTICS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 't'
< result
< column "INDEX_NAME" "VARCHAR"
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| PRIMARY
: bytes
| YES
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, CONSTRAINT u UNIQUE (a) INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) USING HASH INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY u (a, b) INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) VISIBLE INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) INVISIBLE VISIBLE)
< ok 0 0
> query
| SELECT DISTINCT INDEX_NAME, IS_VISIBLE FROM information_schema.STATISTICS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 't'
< result
< column "INDEX_NAME" "VARCHAR"
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| u
: bytes
| YES
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY u (a) INVISIBLE, UNIQUE KEY v (b))
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY v (b), UNIQUE KEY u (a) INVISIBLE)
< ok 0 0
> query
| SELECT DISTINCT INDEX_NAME, IS_VISIBLE FROM information_schema.STATISTICS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 't'
< result
< column "INDEX_NAME" "VARCHAR"
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| u
: bytes
| NO
< row
: bytes
| v
: bytes
| YES
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, KEY k (b), UNIQUE KEY u (a) INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, KEY k (a) INVISIBLE, UNIQUE KEY v (b))
< ok 0 0
> query
| SELECT DISTINCT INDEX_NAME, IS_VISIBLE FROM information_schema.STATISTICS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 't'
< result
< column "INDEX_NAME" "VARCHAR"
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| k
: bytes
| NO
< row
: bytes
| v
: bytes
| YES
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT, UNIQUE KEY u (a) INVISIBLE)
< ok 0 0
> query
| SELECT DISTINCT INDEX_NAME, IS_VISIBLE FROM information_schema.STATISTICS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 't'
< result
< column "INDEX_NAME" "VARCHAR"
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| u
: bytes
| NO
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT, UNIQUE KEY u (a, b) INVISIBLE)
< ok 0 0
> query
| SELECT DISTINCT INDEX_NAME, IS_VISIBLE FROM information_schema.STATISTICS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 't'
< result
< column "INDEX_NAME" "VARCHAR"
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| u
: bytes
| NO
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT, UNIQUE KEY v (b), UNIQUE KEY u (a) INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT, PRIMARY KEY (b), UNIQUE KEY u (a) INVISIBLE)
< ok 0 0
> query
| SELECT DISTINCT INDEX_NAME, IS_VISIBLE FROM information_schema.STATISTICS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 't'
< result
< column "INDEX_NAME" "VARCHAR"
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| PRIMARY
: bytes
| YES
< row
: bytes
| u
: bytes
| NO
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY u (a) INVISIBLE, PRIMARY KEY (b))
< ok 0 0
> query
| SELECT DISTINCT INDEX_NAME, IS_VISIBLE FROM information_schema.STATISTICS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 't'
< result
< column "INDEX_NAME" "VARCHAR"
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| PRIMARY
: bytes
| YES
< row
: bytes
| u
: bytes
| NO
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT PRIMARY KEY, b INT NOT NULL, UNIQUE KEY u (b) INVISIBLE)
< ok 0 0
> query
| SELECT DISTINCT INDEX_NAME, IS_VISIBLE FROM information_schema.STATISTICS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 't'
< result
< column "INDEX_NAME" "VARCHAR"
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| PRIMARY
: bytes
| YES
< row
: bytes
| u
: bytes
| NO
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL UNIQUE, b INT NOT NULL, UNIQUE KEY u (b) INVISIBLE)
< ok 0 0
> query
| SELECT DISTINCT INDEX_NAME, IS_VISIBLE FROM information_schema.STATISTICS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 't'
< result
< column "INDEX_NAME" "VARCHAR"
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| a
: bytes
| YES
< row
: bytes
| u
: bytes
| NO
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a SERIAL, b INT NOT NULL, UNIQUE KEY u (b) INVISIBLE)
< ok 0 0
> query
| SELECT DISTINCT INDEX_NAME, IS_VISIBLE FROM information_schema.STATISTICS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 't'
< result
< column "INDEX_NAME" "VARCHAR"
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| a
: bytes
| YES
< row
: bytes
| u
: bytes
| NO
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) INVISIBLE, b INT NOT NULL UNIQUE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, /*!80000 UNIQUE KEY u (a) INVISIBLE, */ c INT NOT NULL UNIQUE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a VARCHAR(10) NOT NULL, UNIQUE KEY u (a(3)) INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a VARCHAR(10) NOT NULL, b INT NOT NULL, UNIQUE KEY u (a(3)), UNIQUE KEY v (b) INVISIBLE)
< ok 0 0
> query
| SELECT DISTINCT INDEX_NAME, IS_VISIBLE FROM information_schema.STATISTICS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 't'
< result
< column "INDEX_NAME" "VARCHAR"
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| u
: bytes
| YES
< row
: bytes
| v
: bytes
| NO
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a VARCHAR(10) NOT NULL, b INT NOT NULL, UNIQUE KEY u (a(3)) INVISIBLE, UNIQUE KEY v (b))
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b BLOB NOT NULL, UNIQUE KEY v (b(4)) INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a VARCHAR(10) NOT NULL, b INT NOT NULL, UNIQUE KEY v (b) INVISIBLE, UNIQUE KEY u (a(3)))
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT AUTO_INCREMENT, UNIQUE KEY u (a) INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, g INT AS (a + 1) STORED NOT NULL, UNIQUE KEY u (g) INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, g INT AS (a + 1) VIRTUAL NOT NULL, UNIQUE KEY u (g) INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) INVISIBLE) SELECT 1 AS a
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY f ((a + 1)), UNIQUE KEY v (b) INVISIBLE)
< ok 0 0
> query
| SELECT DISTINCT INDEX_NAME, IS_VISIBLE FROM information_schema.STATISTICS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 't'
< result
< column "INDEX_NAME" "VARCHAR"
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| f
: bytes
| YES
< row
: bytes
| v
: bytes
| NO
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT, UNIQUE KEY u ((a + 1)) INVISIBLE)
< ok 0 0
> query
| SELECT DISTINCT INDEX_NAME, IS_VISIBLE FROM information_schema.STATISTICS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 't'
< result
< column "INDEX_NAME" "VARCHAR"
< column "IS_VISIBLE" "VARCHAR"
< row
: bytes
| u
: bytes
| NO
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u ((a + 1)) INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
> exec
| DROP TABLE IF EXISTS t
< ok 0 0
> exec
| CREATE TABLE t (a INT, b INT NOT NULL, UNIQUE KEY f ((a + 1)), UNIQUE KEY v (b) INVISIBLE)
< error
| 3522 HY000 A primary key index cannot be invisible
