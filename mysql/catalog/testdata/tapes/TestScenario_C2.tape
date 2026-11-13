-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestScenario_C2 default#1
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
| CREATE TABLE t (c REAL)
< ok 0 0
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "DATA_TYPE" "LONGTEXT"
< row
: bytes
| double
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| double
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
| CREATE TABLE t (c BOOL)
< ok 0 0
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| tinyint(1)
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "DATA_TYPE" "LONGTEXT"
< row
: bytes
| tinyint
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
| CREATE TABLE t (c INTEGER)
< ok 0 0
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "DATA_TYPE" "LONGTEXT"
< row
: bytes
| int
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| int
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
| CREATE TABLE t (c BOOLEAN)
< ok 0 0
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| tinyint(1)
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
| CREATE TABLE t (a INT1, b INT2, cc INT3, d INT4, e INT8)
< ok 0 0
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'
< result
< column "DATA_TYPE" "LONGTEXT"
< row
: bytes
| tinyint
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "DATA_TYPE" "LONGTEXT"
< row
: bytes
| smallint
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='cc'
< result
< column "DATA_TYPE" "LONGTEXT"
< row
: bytes
| mediumint
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='d'
< result
< column "DATA_TYPE" "LONGTEXT"
< row
: bytes
| int
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='e'
< result
< column "DATA_TYPE" "LONGTEXT"
< row
: bytes
| bigint
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
| CREATE TABLE t (c MIDDLEINT)
< ok 0 0
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "DATA_TYPE" "LONGTEXT"
< row
: bytes
| mediumint
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
| CREATE TABLE t (c INT(11))
< ok 0 0
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| int
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
| CREATE TABLE t (c INT(5) ZEROFILL)
< ok 0 0
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| int(5) unsigned zerofill
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
| CREATE TABLE t (c SERIAL)
< ok 0 0
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| bigint unsigned
> query
| SELECT IS_NULLABLE, EXTRA FROM information_schema.COLUMNS
|              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "IS_NULLABLE" "VARCHAR"
< column "EXTRA" "VARCHAR"
< row
: bytes
| NO
: bytes
| auto_increment
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
| CREATE TABLE t (c NUMERIC(10,2))
< ok 0 0
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| decimal(10,2)
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
| CREATE TABLE t (a DEC(6,2), b FIXED(6,2))
< ok 0 0
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| decimal(6,2)
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| decimal(6,2)
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
| CREATE TABLE t (c DOUBLE PRECISION)
< ok 0 0
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "DATA_TYPE" "LONGTEXT"
< row
: bytes
| double
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
| CREATE TABLE t (a FLOAT4, b FLOAT8)
< ok 0 0
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'
< result
< column "DATA_TYPE" "LONGTEXT"
< row
: bytes
| float
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "DATA_TYPE" "LONGTEXT"
< row
: bytes
| double
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
| CREATE TABLE t (a FLOAT(10), b FLOAT(25))
< ok 0 0
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'
< result
< column "DATA_TYPE" "LONGTEXT"
< row
: bytes
| float
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "DATA_TYPE" "LONGTEXT"
< row
: bytes
| double
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
| CREATE TABLE t (c FLOAT(7,4))
< ok 0 0
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| float(7,4)
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
| CREATE TABLE t (a CHARACTER(10), b CHARACTER VARYING(20))
< ok 0 0
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| char(10)
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| varchar(20)
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
| CREATE TABLE t (a NATIONAL CHAR(10), b NCHAR(10))
< ok 0 0
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| char(10)
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| char(10)
> query
| SELECT IFNULL(CHARACTER_SET_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'
< result
< column "IFNULL(CHARACTER_SET_NAME,'')" "VARCHAR"
< row
: bytes
| utf8mb3
> query
| SELECT IFNULL(CHARACTER_SET_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "IFNULL(CHARACTER_SET_NAME,'')" "VARCHAR"
< row
: bytes
| utf8mb3
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
| CREATE TABLE t (
|   a NVARCHAR(10),
|   b NATIONAL VARCHAR(10),
|   cc NCHAR VARCHAR(10),
|   d NATIONAL CHAR VARYING(10),
|   e NCHAR VARYING(10)
| )
< ok 0 0
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| varchar(10)
> query
| SELECT IFNULL(CHARACTER_SET_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'
< result
< column "IFNULL(CHARACTER_SET_NAME,'')" "VARCHAR"
< row
: bytes
| utf8mb3
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| varchar(10)
> query
| SELECT IFNULL(CHARACTER_SET_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "IFNULL(CHARACTER_SET_NAME,'')" "VARCHAR"
< row
: bytes
| utf8mb3
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='cc'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| varchar(10)
> query
| SELECT IFNULL(CHARACTER_SET_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='cc'
< result
< column "IFNULL(CHARACTER_SET_NAME,'')" "VARCHAR"
< row
: bytes
| utf8mb3
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='d'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| varchar(10)
> query
| SELECT IFNULL(CHARACTER_SET_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='d'
< result
< column "IFNULL(CHARACTER_SET_NAME,'')" "VARCHAR"
< row
: bytes
| utf8mb3
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='e'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| varchar(10)
> query
| SELECT IFNULL(CHARACTER_SET_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='e'
< result
< column "IFNULL(CHARACTER_SET_NAME,'')" "VARCHAR"
< row
: bytes
| utf8mb3
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
| CREATE TABLE t (a LONG, b LONG VARCHAR, cc LONG VARBINARY)
< ok 0 0
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'
< result
< column "DATA_TYPE" "LONGTEXT"
< row
: bytes
| mediumtext
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "DATA_TYPE" "LONGTEXT"
< row
: bytes
| mediumtext
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='cc'
< result
< column "DATA_TYPE" "LONGTEXT"
< row
: bytes
| mediumblob
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
| CREATE TABLE t (a CHAR, b BINARY)
< ok 0 0
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| char(1)
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| binary(1)
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
| CREATE TABLE t (c VARCHAR)
< error
| 1064 42000 …near ')' at line 1
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
| CREATE TABLE t (
|   a TIMESTAMP NULL,
|   b DATETIME,
|   cc TIME,
|   d TIMESTAMP(6) NULL,
|   e DATETIME(6),
|   f TIME(3)
| )
< ok 0 0
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| timestamp
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| datetime
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='cc'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| time
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='d'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| timestamp(6)
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='e'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| datetime(6)
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='f'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| time(3)
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
| CREATE TABLE t (c YEAR(4))
< ok 0 0
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| year
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
| CREATE TABLE t (c BIT)
< ok 0 0
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| bit(1)
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
| SET SESSION sql_mode=''
< ok 0 0
> exec
| CREATE TABLE t (c VARCHAR(65536))
< ok 0 0
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS
|              WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "DATA_TYPE" "LONGTEXT"
< row
: bytes
| mediumtext
> exec
| SET SESSION sql_mode='STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION'
< ok 0 0
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
| CREATE TABLE t (
|   a TEXT(100),
|   b TEXT(1000),
|   cc TEXT(70000),
|   d TEXT(20000000)
| )
< ok 0 0
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'
< result
< column "DATA_TYPE" "LONGTEXT"
< row
: bytes
| text
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "DATA_TYPE" "LONGTEXT"
< row
: bytes
| text
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='cc'
< result
< column "DATA_TYPE" "LONGTEXT"
< row
: bytes
| mediumtext
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='d'
< result
< column "DATA_TYPE" "LONGTEXT"
< row
: bytes
| longtext
