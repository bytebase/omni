-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseRefusals default#1
> exec
| DROP DATABASE IF EXISTS obrefuse_00
< ok 0 0
> exec
| CREATE DATABASE obrefuse_00
< ok 1 0
> exec
| USE obrefuse_00
< ok 0 0
> exec
| CREATE TABLE t1 (a JSON DEFAULT (JSON_ARRAY()))
< error
| 1101 42000 BLOB, TEXT, GEOMETRY or JSON column 'a' can't have a default value
> exec
| CREATE TABLE t2 (a JSON DEFAULT '[]')
< error
| 1101 42000 BLOB, TEXT, GEOMETRY or JSON column 'a' can't have a default value
> exec
| CREATE TABLE t3 (a JSON DEFAULT NULL)
< ok 0 0
> exec
| CREATE TABLE t4 (a BLOB DEFAULT (SUBSTRING('abcdef', 1, 3)))
< error
| 1067 42000 Invalid default value for 'a'
> exec
| CREATE TABLE t5 (a TEXT DEFAULT ('x'))
< error
| 1067 42000 Invalid default value for 'a'
> exec
| CREATE TABLE t6 (a TEXT DEFAULT 'x')
< error
| 1067 42000 Invalid default value for 'a'
> exec
| CREATE TABLE t6 (a TEXT DEFAULT NULL, b BLOB DEFAULT NULL, c POINT DEFAULT NULL)
< ok 0 0
> exec
| CREATE TABLE t7 (a POINT DEFAULT (POINT(0, 0)))
< error
| 1101 42000 BLOB, TEXT, GEOMETRY or JSON column 'a' can't have a default value
> exec
| CREATE TABLE t8 (a BINARY(16) DEFAULT (UUID_TO_BIN(UUID())))
< ok 0 0
> exec
| CREATE TABLE t9 (id INT)
< ok 0 0
> exec
| ALTER TABLE t9 ADD COLUMN a JSON DEFAULT (JSON_ARRAY())
< error
| 1101 42000 BLOB, TEXT, GEOMETRY or JSON column 'a' can't have a default value
> exec
| ALTER TABLE t9 ADD COLUMN b JSON
< ok 0 0
> exec
| ALTER TABLE t9 ALTER COLUMN b SET DEFAULT (JSON_ARRAY())
< error
| 1101 42000 BLOB, TEXT, GEOMETRY or JSON column 'b' can't have a default value
> exec
| ALTER TABLE t9 MODIFY COLUMN b JSON DEFAULT (JSON_OBJECT())
< error
| 1101 42000 BLOB, TEXT, GEOMETRY or JSON column 'b' can't have a default value
> exec
| DROP DATABASE IF EXISTS obrefuse_01
< ok 0 0
> exec
| CREATE DATABASE obrefuse_01
< ok 1 0
> exec
| USE obrefuse_01
< ok 0 0
> exec
| CREATE TABLE t1 (a INT, b INT, c INT, PRIMARY KEY (a, b), UNIQUE KEY uc (c))
< ok 0 0
> exec
| ALTER TABLE t1 DROP COLUMN a
< error
| 1235 0A000 drop rowkey column is not supported
> exec
| ALTER TABLE t1 DROP COLUMN b
< error
| 1235 0A000 drop rowkey column is not supported
> exec
| ALTER TABLE t1 DROP COLUMN c
< ok 0 0
> exec
| CREATE TABLE t2 (a INT PRIMARY KEY, b INT)
< ok 0 0
> exec
| ALTER TABLE t2 DROP COLUMN a
< error
| 1235 0A000 drop rowkey column is not supported
> exec
| ALTER TABLE t2 DROP COLUMN a, DROP PRIMARY KEY
< error
| 1235 0A000 Multiple complex DDLs about primary in single stmt not supported
> exec
| DROP DATABASE IF EXISTS obrefuse_02
< ok 0 0
> exec
| CREATE DATABASE obrefuse_02
< ok 1 0
> exec
| USE obrefuse_02
< ok 0 0
> exec
| CREATE TABLE t1 (a INT, b INT)
< ok 0 0
> exec
| CREATE TABLE t2 (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW v1 AS SELECT t1.a, t2.a FROM t1 JOIN t2 ON t1.a = t2.a
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| CREATE VIEW v2 AS SELECT a, a FROM t1
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| CREATE VIEW v3 AS SELECT a AS x, b AS X FROM t1
< error
| 1060 42S21 Duplicate column name 'x'
> exec
| CREATE VIEW v4 (x, x) AS SELECT a, b FROM t1
< error
| 1060 42S21 Duplicate column name 'x'
> exec
| CREATE VIEW v5 (x, y) AS SELECT a, a FROM t1
< ok 0 0
> exec
| CREATE VIEW v6 AS SELECT * FROM t1 JOIN t2 ON t1.a = t2.a
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| CREATE VIEW v7 AS SELECT 1, 1
< ok 0 0
> exec
| CREATE VIEW v8 AS SELECT a + 1, a + 1 FROM t1
< ok 0 0
> exec
| CREATE VIEW v9 AS SELECT a, 1 AS a FROM t1
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| CREATE OR REPLACE VIEW v5 AS SELECT a, a FROM t1
< error
| 1060 42S21 Duplicate column name 'a'
> exec
| CREATE VIEW v10 AS SELECT a FROM t1 UNION SELECT b FROM t1
< ok 0 0
> exec
| DROP DATABASE IF EXISTS obrefuse_03
< ok 0 0
> exec
| CREATE DATABASE obrefuse_03
< ok 1 0
> exec
| USE obrefuse_03
< ok 0 0
> exec
| CREATE EVENT e1 ON SCHEDULE EVERY 1 HOUR DO SELECT 1
< ok 0 0
> exec
| ALTER EVENT e1 RENAME TO e1b
< error
| 1235 0A000 modify EVENT NAME that has already started not supported
> exec
| CREATE EVENT e2 ON SCHEDULE EVERY 1 HOUR STARTS CURRENT_TIMESTAMP + INTERVAL 1 DAY DO SELECT 1
< ok 0 0
> exec
| ALTER EVENT e2 RENAME TO e2b
< ok 0 0
> exec
| CREATE EVENT e3 ON SCHEDULE AT CURRENT_TIMESTAMP + INTERVAL 1 DAY DO SELECT 1
< ok 0 0
> exec
| ALTER EVENT e3 RENAME TO e3b
< ok 0 0
> exec
| CREATE EVENT e4 ON SCHEDULE EVERY 1 HOUR DISABLE DO SELECT 1
< ok 0 0
> exec
| ALTER EVENT e4 RENAME TO e4b
< error
| 1235 0A000 modify EVENT NAME that has already started not supported
> exec
| ALTER EVENT e1 COMMENT 'x'
< ok 0 0
> exec
| ALTER EVENT e2b ON SCHEDULE EVERY 1 HOUR
< ok 0 0
> exec
| ALTER EVENT e2b RENAME TO e2c
< error
| 1235 0A000 modify EVENT NAME that has already started not supported
> exec
| DROP DATABASE IF EXISTS obrefuse_04
< ok 0 0
> exec
| CREATE DATABASE obrefuse_04
< ok 1 0
> exec
| USE obrefuse_04
< ok 0 0
> exec
| CREATE TABLE t1 (a INT) COMPRESSION='ZLIB'
< error
| 1210 HY000 Incorrect arguments to ZLIB
> exec
| CREATE TABLE t2 (a INT) COMPRESSION='LZ4'
< error
| 1210 HY000 Incorrect arguments to LZ4
> exec
| CREATE TABLE t3 (a INT) COMPRESSION='NONE'
< ok 0 0
> exec
| CREATE TABLE t4 (a INT) COMPRESSION='none'
< ok 0 0
> exec
| CREATE TABLE t5 (a INT) COMPRESSION='zstd_1.3.8'
< ok 0 0
> exec
| CREATE TABLE t6 (a INT) COMPRESSION='zstd_1.0'
< ok 0 0
> exec
| CREATE TABLE t7 (a INT) COMPRESSION='lz4_1.0'
< ok 0 0
> exec
| CREATE TABLE t8 (a INT) COMPRESSION='snappy_1.0'
< ok 0 0
> exec
| CREATE TABLE t9 (a INT) COMPRESSION='zlib_1.0'
< error
| 1235 0A000 zlib compressor not supported
> exec
| CREATE TABLE t10 (a INT) COMPRESSION='lz4_1.9.1'
< ok 0 0
> exec
| CREATE TABLE t11 (a INT) COMPRESSION='ZSTD_1.3.8'
< ok 0 0
> exec
| CREATE TABLE t12 (a INT) COMPRESSION=''
< error
| 1210 HY000 Incorrect arguments to 
> exec
| ALTER TABLE t3 COMPRESSION='ZLIB'
< error
| 1210 HY000 Incorrect arguments to ZLIB
> exec
| ALTER TABLE t3 COMPRESSION='lz4_1.0'
< ok 0 0
> exec
| DROP DATABASE IF EXISTS obrefuse_05
< ok 0 0
> exec
| CREATE DATABASE obrefuse_05
< ok 1 0
> exec
| USE obrefuse_05
< ok 0 0
> exec
| CREATE TABLE t1 (id INT NOT NULL, val INT NOT NULL) PARTITION BY RANGE (val) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (30))
< ok 0 0
> exec
| ALTER TABLE t1 REORGANIZE PARTITION p1 INTO (PARTITION p1a VALUES LESS THAN (20), PARTITION p1b VALUES LESS THAN (30))
< error
| 1235 0A000 reorganizing table without primary key(s) is not supported
> exec
| CREATE TABLE t2 (id INT NOT NULL, val INT NOT NULL, PRIMARY KEY (id, val)) PARTITION BY RANGE (val) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (30))
< ok 0 0
> exec
| ALTER TABLE t2 REORGANIZE PARTITION p1 INTO (PARTITION p1a VALUES LESS THAN (20), PARTITION p1b VALUES LESS THAN (30))
< error
| 1235 0A000 mismatching between primary key prefix and partition key is not supported
> exec
| CREATE TABLE t3 (id INT NOT NULL, val INT NOT NULL, PRIMARY KEY (val, id)) PARTITION BY RANGE (val) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (30))
< ok 0 0
> exec
| ALTER TABLE t3 REORGANIZE PARTITION p1 INTO (PARTITION p1a VALUES LESS THAN (20), PARTITION p1b VALUES LESS THAN (30))
< ok 0 0
> exec
| CREATE TABLE t4 (id INT NOT NULL, val INT NOT NULL, PRIMARY KEY (val)) PARTITION BY RANGE COLUMNS (val) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (30))
< ok 0 0
> exec
| ALTER TABLE t4 REORGANIZE PARTITION p1 INTO (PARTITION p1a VALUES LESS THAN (20), PARTITION p1b VALUES LESS THAN (30))
< ok 0 0
> exec
| CREATE TABLE t5 (a INT NOT NULL, b INT NOT NULL, c INT NOT NULL, PRIMARY KEY (a, b, c)) PARTITION BY RANGE COLUMNS (b, a) (PARTITION p0 VALUES LESS THAN (10, 10), PARTITION p1 VALUES LESS THAN (30, 30))
< ok 0 0
> exec
| ALTER TABLE t5 REORGANIZE PARTITION p1 INTO (PARTITION p1a VALUES LESS THAN (20, 20), PARTITION p1b VALUES LESS THAN (30, 30))
< error
| 1235 0A000 mismatching between primary key prefix and partition key is not supported
> exec
| CREATE TABLE t6 (id INT NOT NULL, val INT NOT NULL, PRIMARY KEY (val, id)) PARTITION BY RANGE (val + 1) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (30))
< ok 0 0
> exec
| ALTER TABLE t6 REORGANIZE PARTITION p1 INTO (PARTITION p1a VALUES LESS THAN (20), PARTITION p1b VALUES LESS THAN (30))
< error
| 1235 0A000 mismatching between primary key prefix and partition key is not supported
> exec
| CREATE TABLE t7 (a INT NOT NULL, b INT NOT NULL, c INT NOT NULL, PRIMARY KEY (a, b, c)) PARTITION BY RANGE COLUMNS (a, b) (PARTITION p0 VALUES LESS THAN (10, 10), PARTITION p1 VALUES LESS THAN (30, 30))
< ok 0 0
> exec
| ALTER TABLE t7 REORGANIZE PARTITION p1 INTO (PARTITION p1a VALUES LESS THAN (20, 20), PARTITION p1b VALUES LESS THAN (30, 30))
< ok 0 0
> exec
| DROP DATABASE IF EXISTS obrefuse_06
< ok 0 0
> exec
| CREATE DATABASE obrefuse_06
< ok 1 0
> exec
| USE obrefuse_06
< ok 0 0
> exec
| SET SESSION sql_generate_invisible_primary_key = 1
< error
| 1193 HY000 Unknown system variable 'sql_generate_invisible_primary_key'
> exec
| SET SESSION show_gipk_in_create_table_and_information_schema = 1
< error
| 1193 HY000 Unknown system variable 'show_gipk_in_create_table_and_information_schema'
> exec
| SET sql_generate_invisible_primary_key = ON
< error
| 1193 HY000 Unknown system variable 'sql_generate_invisible_primary_key'
> exec
| SET @@session.sql_generate_invisible_primary_key = 0
< error
| 1193 HY000 Unknown system variable 'sql_generate_invisible_primary_key'
> exec
| CREATE TABLE t (my_row_id INT, a INT)
< ok 0 0
