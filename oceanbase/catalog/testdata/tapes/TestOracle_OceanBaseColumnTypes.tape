-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseColumnTypes default#1
> exec
| CREATE TABLE t00 (c TINYINT)
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t00
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| tinyint(4)
: null
: null
> exec
| CREATE TABLE t01 (c SMALLINT)
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t01
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| smallint(6)
: null
: null
> exec
| CREATE TABLE t02 (c MEDIUMINT)
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t02
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| mediumint(9)
: null
: null
> exec
| CREATE TABLE t03 (c INT)
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t03
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| int(11)
: null
: null
> exec
| CREATE TABLE t04 (c INTEGER)
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t04
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| int(11)
: null
: null
> exec
| CREATE TABLE t05 (c BIGINT)
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t05
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| bigint(20)
: null
: null
> exec
| CREATE TABLE t06 (c TINYINT UNSIGNED)
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t06
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| tinyint(3) unsigned
: null
: null
> exec
| CREATE TABLE t07 (c SMALLINT UNSIGNED)
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t07
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| smallint(5) unsigned
: null
: null
> exec
| CREATE TABLE t08 (c MEDIUMINT UNSIGNED)
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t08
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| mediumint(8) unsigned
: null
: null
> exec
| CREATE TABLE t09 (c INT UNSIGNED)
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t09
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| int(10) unsigned
: null
: null
> exec
| CREATE TABLE t10 (c BIGINT UNSIGNED)
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t10
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| bigint(20) unsigned
: null
: null
> exec
| CREATE TABLE t11 (c INT(5))
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t11
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| int(5)
: null
: null
> exec
| CREATE TABLE t12 (c INT(5) UNSIGNED)
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t12
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| int(5) unsigned
: null
: null
> exec
| CREATE TABLE t13 (c INT(5) ZEROFILL)
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t13
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| int(5) unsigned zerofill
: null
: null
> exec
| CREATE TABLE t14 (c INT ZEROFILL)
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t14
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| int(10) unsigned zerofill
: null
: null
> exec
| CREATE TABLE t15 (c BIGINT(3))
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t15
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| bigint(3)
: null
: null
> exec
| CREATE TABLE t16 (c TINYINT(1))
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t16
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| tinyint(1)
: null
: null
> exec
| CREATE TABLE t17 (c BOOLEAN)
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t17
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| tinyint(1)
: null
: null
> exec
| CREATE TABLE t18 (c BOOL)
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t18
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| tinyint(1)
: null
: null
> exec
| CREATE TABLE t19 (c SERIAL)
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t19
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| bigint(20) unsigned
: null
: null
> exec
| CREATE TABLE t20 (c YEAR)
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t20
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| year(4)
: null
: null
> exec
| CREATE TABLE t21 (c YEAR(4))
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t21
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| year(4)
: null
: null
> exec
| CREATE TABLE t22 (c YEAR(2))
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t22
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| year(4)
: null
: null
> exec
| CREATE TABLE t23 (c YEAR(3))
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t23
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| year(4)
: null
: null
> exec
| CREATE TABLE t24 (c YEAR(5))
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t24
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| year(4)
: null
: null
> exec
| CREATE TABLE t25 (c BIT)
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t25
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| bit(1)
: null
: null
> exec
| CREATE TABLE t26 (c DECIMAL)
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t26
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| decimal(10,0)
: null
: null
> exec
| CREATE TABLE t27 (c FLOAT)
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t27
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| float
: null
: null
> exec
| CREATE TABLE t28 (c DOUBLE)
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t28
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| double
: null
: null
> exec
| CREATE TABLE t29 (c VARCHAR(10) CHARACTER SET utf8)
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t29
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| varchar(10)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE TABLE t30 (c VARCHAR(10) CHARACTER SET utf8mb3)
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t30
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| varchar(10)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE TABLE t31 (c VARCHAR(10) CHARACTER SET utf8mb4)
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t31
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| varchar(10)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE TABLE t32 (c VARCHAR(10) COLLATE utf8_general_ci)
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t32
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| varchar(10)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE TABLE t33 (c VARCHAR(10) COLLATE utf8mb3_bin)
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t33
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| varchar(10)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> exec
| CREATE TABLE t34 (c VARCHAR(10) CHARACTER SET utf8 COLLATE utf8_bin)
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t34
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| varchar(10)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> exec
| CREATE TABLE t35 (c VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8_bin)
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t35
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| varchar(10)
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
> exec
| CREATE TABLE t36 (c VARCHAR(10) CHARACTER SET latin1)
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t36
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| varchar(10)
: bytes
| latin1
: bytes
| latin1_swedish_ci
> exec
| CREATE TABLE t37 (c NATIONAL CHAR(10))
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t37
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| char(10)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE TABLE t38 (c NCHAR(10))
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t38
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| char(10)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE TABLE t39 (c NVARCHAR(10))
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t39
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| varchar(10)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE TABLE t40 (c NATIONAL VARCHAR(10))
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t40
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| varchar(10)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE TABLE t41 (c NCHAR VARCHAR(10))
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t41
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| varchar(10)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE TABLE t42 (c NATIONAL CHAR VARYING(10))
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t42
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| varchar(10)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE TABLE t43 (c VARCHAR(65535))
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t43
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| varchar(65535)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE TABLE t44 (c VARCHAR(65536))
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t44
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| varchar(65536)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE TABLE t45 (c VARCHAR(262144))
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t45
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| varchar(262144)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE TABLE t46 (c VARCHAR(262145))
< error
| 1074 42000 Column length too big for column 'c' (max = 262144)
> exec
| CREATE TABLE t47 (c CHAR(255))
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t47
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| char(255)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE TABLE t48 (c CHAR(256))
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t48
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| char(256)
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE TABLE t49 (c CHAR(257))
< error
| 1074 42000 Column length too big for column 'c' (max = 256)
> exec
| CREATE TABLE t50 (c TEXT(10))
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t50
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| tinytext
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE TABLE t51 (c TEXT(1000))
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t51
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| text
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> exec
| CREATE TABLE t52 (c BLOB(10))
< ok 0 0
> query
| SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'c'
: string
| t52
< result
< column "COLUMN_TYPE" "TEXT"
< column "CHARACTER_SET_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| tinyblob
: null
: null
