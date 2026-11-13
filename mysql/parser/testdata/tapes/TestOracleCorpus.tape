-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracleCorpus default#1
> exec
| CREATE TABLE IF NOT EXISTS t (id INT, val INT, score INT, status INT, name VARCHAR(100), dept VARCHAR(50), salary DECIMAL(10,2), d DATE, created_at DATETIME)
< ok 0 0
> exec
| CREATE TABLE t1 (a INT)
< ok 0 0
> exec
| CREATE TABLE t2 (a INT UNSIGNED)
< ok 0 0
> exec
| CREATE TABLE t3 (a INT SIGNED)
< ok 0 0
> exec
| CREATE TABLE t4 (a INT(11))
< ok 0 0
> exec
| CREATE TABLE t5 (a INT UNSIGNED ZEROFILL)
< ok 0 0
> exec
| CREATE TABLE t6 (a TINYINT SIGNED)
< ok 0 0
> exec
| CREATE TABLE t7 (a TINYINT UNSIGNED ZEROFILL)
< ok 0 0
> exec
| CREATE TABLE t8 (a SMALLINT(5) UNSIGNED)
< ok 0 0
> exec
| CREATE TABLE t9 (a MEDIUMINT SIGNED)
< ok 0 0
> exec
| CREATE TABLE t10 (a BIGINT(20) UNSIGNED)
< ok 0 0
> exec
| CREATE TABLE t11 (a INT1 UNSIGNED)
< ok 0 0
> exec
| CREATE TABLE t12 (a INT8 SIGNED)
< ok 0 0
> exec
| DROP TABLE IF EXISTS t1
< ok 0 0
> exec
| DROP TABLE IF EXISTS t2
< ok 0 0
> exec
| DROP TABLE IF EXISTS t3
< ok 0 0
> exec
| DROP TABLE IF EXISTS t4
< ok 0 0
> exec
| DROP TABLE IF EXISTS t5
< ok 0 0
> exec
| DROP TABLE IF EXISTS t6
< ok 0 0
> exec
| DROP TABLE IF EXISTS t7
< ok 0 0
> exec
| DROP TABLE IF EXISTS t8
< ok 0 0
> exec
| DROP TABLE IF EXISTS t9
< ok 0 0
> exec
| DROP TABLE IF EXISTS t10
< ok 0 0
> exec
| DROP TABLE IF EXISTS t11
< ok 0 0
> exec
| DROP TABLE IF EXISTS t12
< ok 0 0
> exec
| CREATE TABLE t13 (a DECIMAL)
< ok 0 0
> exec
| CREATE TABLE t14 (a DECIMAL(10))
< ok 0 0
> exec
| CREATE TABLE t15 (a DECIMAL(10,2))
< ok 0 0
> exec
| CREATE TABLE t16 (a DECIMAL(10,2) UNSIGNED)
< ok 0 0
> exec
| CREATE TABLE t17 (a NUMERIC(10,2))
< ok 0 0
> exec
| CREATE TABLE t18 (a DEC(10,2))
< ok 0 0
> exec
| CREATE TABLE t19 (a FIXED(10,2))
< ok 0 0
> exec
| CREATE TABLE t20 (a FLOAT)
< ok 0 0
> exec
| CREATE TABLE t21 (a FLOAT(10,2))
< ok 0 0
> exec
| CREATE TABLE t22 (a FLOAT(24))
< ok 0 0
> exec
| CREATE TABLE t23 (a FLOAT(25))
< ok 0 0
> exec
| CREATE TABLE t24 (a DOUBLE)
< ok 0 0
> exec
| CREATE TABLE t25 (a DOUBLE PRECISION)
< ok 0 0
> exec
| CREATE TABLE t26 (a REAL)
< ok 0 0
> exec
| CREATE TABLE t27 (a FLOAT4)
< ok 0 0
> exec
| CREATE TABLE t28 (a FLOAT8)
< ok 0 0
> exec
| DROP TABLE IF EXISTS t13
< ok 0 0
> exec
| DROP TABLE IF EXISTS t14
< ok 0 0
> exec
| DROP TABLE IF EXISTS t15
< ok 0 0
> exec
| DROP TABLE IF EXISTS t16
< ok 0 0
> exec
| DROP TABLE IF EXISTS t17
< ok 0 0
> exec
| DROP TABLE IF EXISTS t18
< ok 0 0
> exec
| DROP TABLE IF EXISTS t19
< ok 0 0
> exec
| DROP TABLE IF EXISTS t20
< ok 0 0
> exec
| DROP TABLE IF EXISTS t21
< ok 0 0
> exec
| DROP TABLE IF EXISTS t22
< ok 0 0
> exec
| DROP TABLE IF EXISTS t23
< ok 0 0
> exec
| DROP TABLE IF EXISTS t24
< ok 0 0
> exec
| DROP TABLE IF EXISTS t25
< ok 0 0
> exec
| DROP TABLE IF EXISTS t26
< ok 0 0
> exec
| DROP TABLE IF EXISTS t27
< ok 0 0
> exec
| DROP TABLE IF EXISTS t28
< ok 0 0
> exec
| CREATE TABLE t29 (a CHAR(10))
< ok 0 0
> exec
| CREATE TABLE t30 (a CHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci)
< ok 0 0
> exec
| CREATE TABLE t31 (a VARCHAR(255))
< ok 0 0
> exec
| CREATE TABLE t32 (a TEXT)
< ok 0 0
> exec
| CREATE TABLE t33 (a TEXT(1000))
< ok 0 0
> exec
| CREATE TABLE t34 (a TINYTEXT)
< ok 0 0
> exec
| CREATE TABLE t35 (a MEDIUMTEXT CHARACTER SET latin1)
< ok 0 0
> exec
| CREATE TABLE t36 (a LONGTEXT)
< ok 0 0
> exec
| CREATE TABLE t37 (a LONG)
< ok 0 0
> exec
| CREATE TABLE t38 (a LONG VARCHAR)
< ok 0 0
> exec
| CREATE TABLE t39 (a BINARY(16))
< ok 0 0
> exec
| CREATE TABLE t40 (a VARBINARY(255))
< ok 0 0
> exec
| CREATE TABLE t41 (a BLOB)
< ok 0 0
> exec
| CREATE TABLE t42 (a BLOB(1000))
< ok 0 0
> exec
| CREATE TABLE t43 (a TINYBLOB)
< ok 0 0
> exec
| CREATE TABLE t44 (a MEDIUMBLOB)
< ok 0 0
> exec
| CREATE TABLE t45 (a LONGBLOB)
< ok 0 0
> exec
| CREATE TABLE t46 (a LONG VARBINARY)
< ok 0 0
> exec
| CREATE TABLE t47 (a NATIONAL CHAR(10))
< ok 0 0
> exec
| CREATE TABLE t48 (a NCHAR(10))
< ok 0 0
> exec
| CREATE TABLE t49 (a NVARCHAR(100))
< ok 0 0
> exec
| DROP TABLE IF EXISTS t29
< ok 0 0
> exec
| DROP TABLE IF EXISTS t30
< ok 0 0
> exec
| DROP TABLE IF EXISTS t31
< ok 0 0
> exec
| DROP TABLE IF EXISTS t32
< ok 0 0
> exec
| DROP TABLE IF EXISTS t33
< ok 0 0
> exec
| DROP TABLE IF EXISTS t34
< ok 0 0
> exec
| DROP TABLE IF EXISTS t35
< ok 0 0
> exec
| DROP TABLE IF EXISTS t36
< ok 0 0
> exec
| DROP TABLE IF EXISTS t37
< ok 0 0
> exec
| DROP TABLE IF EXISTS t38
< ok 0 0
> exec
| DROP TABLE IF EXISTS t39
< ok 0 0
> exec
| DROP TABLE IF EXISTS t40
< ok 0 0
> exec
| DROP TABLE IF EXISTS t41
< ok 0 0
> exec
| DROP TABLE IF EXISTS t42
< ok 0 0
> exec
| DROP TABLE IF EXISTS t43
< ok 0 0
> exec
| DROP TABLE IF EXISTS t44
< ok 0 0
> exec
| DROP TABLE IF EXISTS t45
< ok 0 0
> exec
| DROP TABLE IF EXISTS t46
< ok 0 0
> exec
| DROP TABLE IF EXISTS t47
< ok 0 0
> exec
| DROP TABLE IF EXISTS t48
< ok 0 0
> exec
| DROP TABLE IF EXISTS t49
< ok 0 0
> exec
| CREATE TABLE t50 (a DATE)
< ok 0 0
> exec
| CREATE TABLE t51 (a TIME)
< ok 0 0
> exec
| CREATE TABLE t52 (a TIME(3))
< ok 0 0
> exec
| CREATE TABLE t53 (a DATETIME)
< ok 0 0
> exec
| CREATE TABLE t54 (a DATETIME(6))
< ok 0 0
> exec
| CREATE TABLE t55 (a TIMESTAMP)
< ok 0 0
> exec
| CREATE TABLE t56 (a TIMESTAMP(3))
< ok 0 0
> exec
| CREATE TABLE t57 (a YEAR)
< ok 0 0
> exec
| CREATE TABLE t58 (a BIT(8))
< ok 0 0
> exec
| CREATE TABLE t59 (a BOOL)
< ok 0 0
> exec
| CREATE TABLE t60 (a JSON)
< ok 0 0
> exec
| CREATE TABLE t61 (a SERIAL)
< ok 0 0
> exec
| CREATE TABLE t62 (a ENUM('a','b','c'))
< ok 0 0
> exec
| CREATE TABLE t63 (a SET('x','y','z'))
< ok 0 0
> exec
| CREATE TABLE t64 (a GEOMETRY)
< ok 0 0
> exec
| CREATE TABLE t65 (a POINT)
< ok 0 0
> exec
| CREATE TABLE t66 (a LINESTRING)
< ok 0 0
> exec
| CREATE TABLE t67 (a POLYGON)
< ok 0 0
> exec
| CREATE TABLE t68 (a GEOMETRYCOLLECTION)
< ok 0 0
> exec
| DROP TABLE IF EXISTS t50
< ok 0 0
> exec
| DROP TABLE IF EXISTS t51
< ok 0 0
> exec
| DROP TABLE IF EXISTS t52
< ok 0 0
> exec
| DROP TABLE IF EXISTS t53
< ok 0 0
> exec
| DROP TABLE IF EXISTS t54
< ok 0 0
> exec
| DROP TABLE IF EXISTS t55
< ok 0 0
> exec
| DROP TABLE IF EXISTS t56
< ok 0 0
> exec
| DROP TABLE IF EXISTS t57
< ok 0 0
> exec
| DROP TABLE IF EXISTS t58
< ok 0 0
> exec
| DROP TABLE IF EXISTS t59
< ok 0 0
> exec
| DROP TABLE IF EXISTS t60
< ok 0 0
> exec
| DROP TABLE IF EXISTS t61
< ok 0 0
> exec
| DROP TABLE IF EXISTS t62
< ok 0 0
> exec
| DROP TABLE IF EXISTS t63
< ok 0 0
> exec
| DROP TABLE IF EXISTS t64
< ok 0 0
> exec
| DROP TABLE IF EXISTS t65
< ok 0 0
> exec
| DROP TABLE IF EXISTS t66
< ok 0 0
> exec
| DROP TABLE IF EXISTS t67
< ok 0 0
> exec
| DROP TABLE IF EXISTS t68
< ok 0 0
> exec
| SELECT RANK() OVER w FROM t WINDOW w AS (ORDER BY id)
< ok 0 0
> exec
| SELECT ROW_NUMBER() OVER (PARTITION BY dept ORDER BY salary DESC) FROM t
< ok 0 0
> exec
| SELECT LAG(val) OVER (ORDER BY id) FROM t
< ok 0 0
> exec
| SELECT LAG(val, 2, 0) OVER (ORDER BY id) FROM t
< ok 0 0
> exec
| SELECT NTH_VALUE(val, 3) OVER (ORDER BY id) FROM t
< ok 0 0
> exec
| SELECT SUM(val) OVER (ORDER BY id ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING) FROM t
< ok 0 0
> exec
| SELECT DENSE_RANK() OVER (ORDER BY score), PERCENT_RANK() OVER (ORDER BY score) FROM t
< ok 0 0
> exec
| SELECT NTILE(4) OVER (ORDER BY id) FROM t
< ok 0 0
> exec
| SELECT FIRST_VALUE(val) OVER (ORDER BY id) FROM t
< ok 0 0
> exec
| SELECT COUNT(*) OVER (ORDER BY id) FROM t
< ok 0 0
> exec
| SELECT COUNT(val) OVER (ORDER BY id) FROM t
< ok 0 0
> exec
| SELECT NTILE() OVER (ORDER BY id) FROM t
< error
| 1064 42000 …near ') OVER (ORDER BY id) FROM t' at line 1
> exec
| SELECT NTILE(2, 3) OVER (ORDER BY id) FROM t
< error
| 1064 42000 …near ', 3) OVER (ORDER BY id) FROM t' at line 1
> exec
| SELECT LEAD() OVER (ORDER BY id) FROM t
< error
| 1064 42000 …near ') OVER (ORDER BY id) FROM t' at line 1
> exec
| SELECT FIRST_VALUE() OVER (ORDER BY id) FROM t
< error
| 1064 42000 …near ') OVER (ORDER BY id) FROM t' at line 1
> exec
| SELECT NTH_VALUE(val) OVER (ORDER BY id) FROM t
< error
| 1064 42000 …near ') OVER (ORDER BY id) FROM t' at line 1
> exec
| SELECT ROW_NUMBER(id) OVER (ORDER BY id) FROM t
< error
| 1064 42000 …near 'id) OVER (ORDER BY id) FROM t' at line 1
> exec
| SELECT ABS(val) OVER (ORDER BY id) FROM t
< error
| 1064 42000 …near 'OVER (ORDER BY id) FROM t' at line 1
> exec
| SELECT INTERVAL(val, 1, 2) OVER () FROM t
< error
| 1064 42000 …near 'OVER () FROM t' at line 1
> exec
| SELECT `sum`(val) OVER (ORDER BY id) FROM t
< error
| 1064 42000 …near 'OVER (ORDER BY id) FROM t' at line 1
> exec
| SELECT NOW() + INTERVAL 1 DAY
< ok 0 0
> exec
| SELECT NOW() - INTERVAL 2 HOUR
< ok 0 0
> exec
| SELECT DATE_ADD('2024-01-01', INTERVAL 1 MONTH)
< ok 0 0
> exec
| SELECT DATE_SUB(NOW(), INTERVAL 30 MINUTE)
< ok 0 0
> exec
| SELECT TIMESTAMPADD(MINUTE, 30, NOW())
< ok 0 0
> exec
| SELECT EXTRACT(HOUR FROM NOW())
< ok 0 0
> exec
| SELECT * FROM t WHERE created_at > NOW() - INTERVAL 7 DAY
< ok 0 0
> exec
| SELECT * FROM t WHERE created_at BETWEEN NOW() - INTERVAL 1 MONTH AND NOW()
< ok 0 0
> exec
| CREATE TABLE t_rej1 (a VARCHAR UNSIGNED)
< error
| 1064 42000 …near 'UNSIGNED)' at line 1
> exec
| CREATE TABLE t_rej2 (a TEXT ZEROFILL)
< error
| 1064 42000 …near 'ZEROFILL)' at line 1
> exec
| CREATE TABLE t_rej3 (a JSON UNSIGNED)
< error
| 1064 42000 …near 'UNSIGNED)' at line 1
> exec
| CREATE TABLE select (a INT)
< error
| 1064 42000 …near 'select (a INT)' at line 1
> exec
| CREATE TABLE t_rej5 (select INT)
< error
| 1064 42000 …near 'INT)' at line 1
> exec
| ALTER TABLE t PARTITION BY
< error
| 1064 42000 …near '' at line 1
> exec
| ALTER TABLE t PARTITION BY RANGE(id)
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| ALTER TABLE t PARTITION BY RANGE(id)
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| ALTER TABLE t PARTITION BY RANGE(id)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT DATE_ADD(d, INTERVAL 1 INVALID_UNIT) FROM t
< error
| 1064 42000 …near 'INVALID_UNIT) FROM t' at line 1
> exec
| DROP TABLE IF EXISTS t_rej1
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_rej2
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_rej3
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_rej4
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_rej5
< ok 0 0
