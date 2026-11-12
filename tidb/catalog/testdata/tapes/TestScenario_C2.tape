-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestScenario_C2 default#1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (c REAL)
< ok 0 0
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "DATA_TYPE" "VARCHAR"
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
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
< column "DATA_TYPE" "VARCHAR"
< row
: bytes
| tinyint
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (c INTEGER)
< ok 0 0
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "DATA_TYPE" "VARCHAR"
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
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a INT1, b INT2, cc INT3, d INT4, e INT8)
< ok 0 0
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'
< result
< column "DATA_TYPE" "VARCHAR"
< row
: bytes
| tinyint
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "DATA_TYPE" "VARCHAR"
< row
: bytes
| smallint
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='cc'
< result
< column "DATA_TYPE" "VARCHAR"
< row
: bytes
| mediumint
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='d'
< result
< column "DATA_TYPE" "VARCHAR"
< row
: bytes
| int
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='e'
< result
< column "DATA_TYPE" "VARCHAR"
< row
: bytes
| bigint
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (c MIDDLEINT)
< ok 0 0
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "DATA_TYPE" "VARCHAR"
< row
: bytes
| mediumint
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (c INT(5) ZEROFILL)
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
|   `c` int(5) unsigned zerofill DEFAULT NULL
| ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin
> query
| SELECT COLUMN_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| int(5) unsigned
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (c DOUBLE PRECISION)
< ok 0 0
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='c'
< result
< column "DATA_TYPE" "VARCHAR"
< row
: bytes
| double
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a FLOAT4, b FLOAT8)
< ok 0 0
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'
< result
< column "DATA_TYPE" "VARCHAR"
< row
: bytes
| float
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "DATA_TYPE" "VARCHAR"
< row
: bytes
| double
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a FLOAT(10), b FLOAT(25))
< ok 0 0
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'
< result
< column "DATA_TYPE" "VARCHAR"
< row
: bytes
| float
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "DATA_TYPE" "VARCHAR"
< row
: bytes
| double
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
| utf8mb4
> query
| SELECT IFNULL(CHARACTER_SET_NAME,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "IFNULL(CHARACTER_SET_NAME,'')" "VARCHAR"
< row
: bytes
| utf8mb4
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
| utf8mb4
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
| utf8mb4
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
| utf8mb4
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
| utf8mb4
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
| utf8mb4
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a LONG, b LONG VARCHAR, cc LONG VARBINARY)
< ok 0 0
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='a'
< result
< column "DATA_TYPE" "VARCHAR"
< row
: bytes
| mediumtext
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "DATA_TYPE" "VARCHAR"
< row
: bytes
| mediumtext
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='cc'
< result
< column "DATA_TYPE" "VARCHAR"
< row
: bytes
| mediumblob
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (c VARCHAR)
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 26 near ")" 
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
| year(4)
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
< column "DATA_TYPE" "VARCHAR"
< row
: bytes
| mediumtext
> exec
| SET SESSION sql_mode='STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION'
< ok 0 0
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 0 0
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
< column "DATA_TYPE" "VARCHAR"
< row
: bytes
| text
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='b'
< result
< column "DATA_TYPE" "VARCHAR"
< row
: bytes
| text
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='cc'
< result
< column "DATA_TYPE" "VARCHAR"
< row
: bytes
| mediumtext
> query
| SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='t' AND COLUMN_NAME='d'
< result
< column "DATA_TYPE" "VARCHAR"
< row
: bytes
| longtext
