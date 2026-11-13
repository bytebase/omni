-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseAlterTableClauses default#1
> exec
| DROP DATABASE IF EXISTS obalter_00
< ok 0 0
> exec
| CREATE DATABASE obalter_00
< ok 1 0
> exec
| USE obalter_00
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, a INT)
< ok 0 0
> exec
| CREATE TABLE u (id INT NOT NULL PRIMARY KEY, a INT)
< ok 0 0
> exec
| ALTER TABLE u ADD COLUMN seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT
< ok 0 0
> exec
| ALTER TABLE u ADD UNIQUE KEY uk_seq (seq)
< ok 0 0
> exec
| ALTER TABLE u DROP INDEX uk_seq
< ok 0 0
> exec
| DROP DATABASE IF EXISTS obalter_01
< ok 0 0
> exec
| CREATE DATABASE obalter_01
< ok 1 0
> exec
| USE obalter_01
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL AUTO_INCREMENT, b INT NOT NULL, KEY ka (a), KEY kb (b))
< ok 0 0
> exec
| ALTER TABLE t MODIFY COLUMN b INT NOT NULL AUTO_INCREMENT
< error
| 1075 42000 Incorrect table definition; there can be only one auto column
> exec
| ALTER TABLE t MODIFY COLUMN a INT NOT NULL, MODIFY COLUMN b INT NOT NULL AUTO_INCREMENT
< error
| 1075 42000 Incorrect table definition; there can be only one auto column
> exec
| ALTER TABLE t MODIFY COLUMN a INT NOT NULL
< ok 0 0
> exec
| ALTER TABLE t MODIFY COLUMN b INT NOT NULL AUTO_INCREMENT
< ok 0 0
> exec
| DROP DATABASE IF EXISTS obalter_02
< ok 0 0
> exec
| CREATE DATABASE obalter_02
< ok 1 0
> exec
| USE obalter_02
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, c INT NOT NULL, a INT, PRIMARY KEY (id), UNIQUE KEY uk_c (c), KEY ka (a))
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN seq BIGINT NOT NULL AUTO_INCREMENT, DROP PRIMARY KEY, ADD PRIMARY KEY (seq)
< error
| 1235 0A000 Multiple complex DDLs about primary in single stmt not supported
> exec
| ALTER TABLE t DROP INDEX uk_c, ADD UNIQUE KEY uk_c (c, a)
< error
| 4179 HY000 Dropping and adding indexes at the same time is a high-risk operation, which is not supported
> exec
| ALTER TABLE t DROP INDEX ka, DROP COLUMN a
< error
| 1235 0A000 table redefinition and drop index in single statment not supported
> exec
| ALTER TABLE t MODIFY COLUMN c INT NOT NULL AUTO_INCREMENT, ADD UNIQUE KEY uk_c2 (c)
< error
| 1235 0A000 modify auto_increment and add index in single statment not supported
> exec
| ALTER TABLE t DROP PRIMARY KEY, ADD PRIMARY KEY (c), DROP INDEX uk_c
< error
| 1235 0A000 Multiple complex DDLs about primary in single stmt not supported
> exec
| ALTER TABLE t DROP PRIMARY KEY, ADD PRIMARY KEY (c)
< ok 0 0
> exec
| CREATE TABLE u (a INT)
< ok 0 0
> exec
| ALTER TABLE u ADD COLUMN seq BIGINT NOT NULL AUTO_INCREMENT, ADD PRIMARY KEY (seq)
< error
| 1176 42000 Key 'seq' doesn't exist in table 'u'
> exec
| ALTER TABLE u ADD COLUMN seq BIGINT NOT NULL AUTO_INCREMENT
< ok 0 0
> exec
| ALTER TABLE u ADD PRIMARY KEY (seq)
< ok 0 0
> exec
| DROP DATABASE IF EXISTS obalter_03
< ok 0 0
> exec
| CREATE DATABASE obalter_03
< ok 1 0
> exec
| USE obalter_03
< ok 0 0
> exec
| CREATE TABLE t (a INT NOT NULL, c INT NOT NULL, b INT NOT NULL, PRIMARY KEY (c, a, b))
< ok 0 0
> exec
| ALTER TABLE t DROP COLUMN a
< error
| 1235 0A000 drop rowkey column is not supported
> exec
| ALTER TABLE t MODIFY COLUMN b INT DEFAULT NULL
< error
| 1171 42000 All parts of a PRIMARY KEY must be NOT NULL; if you need NULL in a key, use UNIQUE instead
> exec
| ALTER TABLE t DROP PRIMARY KEY, ADD PRIMARY KEY (c)
< ok 0 0
> exec
| ALTER TABLE t DROP COLUMN a
< ok 0 0
> exec
| ALTER TABLE t MODIFY COLUMN b INT DEFAULT NULL
< ok 0 0
> exec
| ALTER TABLE t DROP COLUMN c
< error
| 1235 0A000 drop rowkey column is not supported
> exec
| ALTER TABLE t DROP PRIMARY KEY
< ok 0 0
> exec
| ALTER TABLE t DROP COLUMN c
< ok 0 0
> exec
| DROP DATABASE IF EXISTS obalter_04
< ok 0 0
> exec
| CREATE DATABASE obalter_04
< ok 1 0
> exec
| USE obalter_04
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT, a INT, UNIQUE KEY uk_c (c), KEY k_ca (c, a))
< ok 0 0
> exec
| ALTER TABLE t DROP COLUMN c
< ok 0 0
> exec
| ALTER TABLE t DROP INDEX uk_c
< error
| 1091 42000 Can't DROP 'uk_c'; check that column/key exists
> exec
| ALTER TABLE t DROP INDEX k_ca
< ok 0 0
> exec
| DROP DATABASE IF EXISTS obalter_05
< ok 0 0
> exec
| CREATE DATABASE obalter_05
< ok 1 0
> exec
| USE obalter_05
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2010), PARTITION pmax VALUES LESS THAN MAXVALUE)
< ok 0 0
> exec
| ALTER TABLE t REMOVE PARTITIONING
< error
| 1064 42000 …near 'PARTITIONING' at line 1
> exec
| ALTER TABLE t DROP COLUMN dt
< error
| 3855 HY000 Column 'dt' has a partitioning function dependency and cannot be dropped or renamed.
> exec
| ALTER TABLE t DROP PRIMARY KEY, ADD PRIMARY KEY (id)
< error
| 1503 HY000 A PRIMARY KEY must include all columns in the table's partitioning function
> exec
| ALTER TABLE t PARTITION BY HASH (id) PARTITIONS 4
< ok 0 0
> exec
| ALTER TABLE t DROP PRIMARY KEY, ADD PRIMARY KEY (id)
< ok 0 0
> exec
| ALTER TABLE t DROP COLUMN dt
< ok 0 0
> exec
| ALTER TABLE t PARTITION BY HASH (id) PARTITIONS 1
< ok 0 0
> exec
| ALTER TABLE t DROP PARTITION p0
< error
| 1235 0A000 Drop hash partition not supported
> exec
| DROP DATABASE IF EXISTS obalter_06
< ok 0 0
> exec
| CREATE DATABASE obalter_06
< ok 1 0
> exec
| USE obalter_06
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY)
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN g GEOMETRY NOT NULL
< error
| 1138 22004 Invalid use of NULL value
> exec
| ALTER TABLE t ADD COLUMN p POINT NOT NULL /*!80003 SRID 0 */
< error
| 1138 22004 Invalid use of NULL value
> exec
| ALTER TABLE t ADD COLUMN j JSON NOT NULL
< ok 0 0
> exec
| ALTER TABLE t ADD COLUMN g GEOMETRY /*!80003 SRID 0 */
< ok 0 0
> exec
| ALTER TABLE t MODIFY COLUMN g GEOMETRY NOT NULL /*!80003 SRID 0 */
< ok 0 0
> exec
| ALTER TABLE t MODIFY COLUMN g GEOMETRY NOT NULL /*!80003 SRID 4326 */
< error
| 1235 0A000 alter geometry srid not supported
> exec
| ALTER TABLE t MODIFY COLUMN g GEOMETRY NOT NULL
< error
| 1235 0A000 alter geometry srid not supported
> exec
| ALTER TABLE t MODIFY COLUMN g GEOMETRY NOT NULL /*!80003 SRID 0 */
< error
| 1138 22004 Invalid use of NULL value
> exec
| ALTER TABLE t MODIFY COLUMN g GEOMETRY NOT NULL /*!80003 SRID 0 */ COMMENT 'c'
< error
| 1138 22004 Invalid use of NULL value
> exec
| ALTER TABLE t MODIFY COLUMN g POINT NOT NULL /*!80003 SRID 0 */
< error
| 1138 22004 Invalid use of NULL value
> exec
| ALTER TABLE t MODIFY COLUMN g POINT /*!80003 SRID 0 */ COMMENT 'c'
< ok 0 0
> exec
| ALTER TABLE t MODIFY COLUMN g POINT NOT NULL /*!80003 SRID 0 */ COMMENT 'c'
< ok 0 0
> exec
| ALTER TABLE t MODIFY COLUMN j JSON NOT NULL COMMENT 'c'
< ok 0 0
> exec
| DROP DATABASE IF EXISTS obalter_07
< ok 0 0
> exec
| CREATE DATABASE obalter_07
< ok 1 0
> exec
| USE obalter_07
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT, g INT GENERATED ALWAYS AS (a + 1) STORED, v INT GENERATED ALWAYS AS (a + 1) VIRTUAL, KEY kg (g))
< ok 0 0
> exec
| ALTER TABLE t MODIFY COLUMN g INT GENERATED ALWAYS AS (a + 2) STORED
< error
| 1235 0A000 Modify generated column definition not supported
> exec
| ALTER TABLE t MODIFY COLUMN v INT GENERATED ALWAYS AS (a + 2) VIRTUAL
< error
| 1235 0A000 Modify generated column definition not supported
> exec
| ALTER TABLE t MODIFY COLUMN g INT GENERATED ALWAYS AS (a + 1) STORED COMMENT 'c'
< ok 0 0
> exec
| ALTER TABLE t DROP COLUMN a
< error
| 3108 HY000 Column 'a' has a generated column dependency
> exec
| ALTER TABLE t DROP COLUMN g
< ok 0 0
> exec
| ALTER TABLE t DROP INDEX kg
< error
| 1091 42000 Can't DROP 'kg'; check that column/key exists
> exec
| ALTER TABLE t ADD COLUMN g INT GENERATED ALWAYS AS (a + 2) STORED
< ok 0 0
> exec
| ALTER TABLE t ADD KEY kg (g)
< ok 0 0
