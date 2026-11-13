-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseNationalCollation default#1
> exec
| CREATE TABLE t_bin (plain CHAR(10), named CHAR(10) CHARACTER SET utf8mb4, n1 NCHAR(10), n2 NVARCHAR(10), n3 NATIONAL VARCHAR(10), own NCHAR(10) COLLATE utf8mb4_bin) DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_bin
< ok 0 0
> exec
| ALTER TABLE t_bin ADD COLUMN n4 NCHAR(10)
< ok 0 0
> query
| SELECT COLUMN_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| t_bin
< result
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| plain
: bytes
| utf8mb4_bin
< row
: bytes
| named
: bytes
| utf8mb4_general_ci
< row
: bytes
| n1
: bytes
| utf8mb4_general_ci
< row
: bytes
| n2
: bytes
| utf8mb4_general_ci
< row
: bytes
| n3
: bytes
| utf8mb4_general_ci
< row
: bytes
| own
: bytes
| utf8mb4_bin
< row
: bytes
| n4
: bytes
| utf8mb4_general_ci
> exec
| CREATE TABLE t_unicode (plain CHAR(10), named CHAR(10) CHARACTER SET utf8mb4, n1 NCHAR(10), n2 NVARCHAR(10), n3 NATIONAL VARCHAR(10), own NCHAR(10) COLLATE utf8mb4_bin) COLLATE utf8mb4_unicode_ci
< ok 0 0
> exec
| ALTER TABLE t_unicode ADD COLUMN n4 NCHAR(10)
< ok 0 0
> query
| SELECT COLUMN_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| t_unicode
< result
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| plain
: bytes
| utf8mb4_unicode_ci
< row
: bytes
| named
: bytes
| utf8mb4_general_ci
< row
: bytes
| n1
: bytes
| utf8mb4_general_ci
< row
: bytes
| n2
: bytes
| utf8mb4_general_ci
< row
: bytes
| n3
: bytes
| utf8mb4_general_ci
< row
: bytes
| own
: bytes
| utf8mb4_bin
< row
: bytes
| n4
: bytes
| utf8mb4_general_ci
> exec
| CREATE TABLE t_latin1 (plain CHAR(10), named CHAR(10) CHARACTER SET utf8mb4, n1 NCHAR(10), n2 NVARCHAR(10), n3 NATIONAL VARCHAR(10), own NCHAR(10) COLLATE utf8mb4_bin) DEFAULT CHARSET latin1
< ok 0 0
> exec
| ALTER TABLE t_latin1 ADD COLUMN n4 NCHAR(10)
< ok 0 0
> query
| SELECT COLUMN_NAME, COLLATION_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| t_latin1
< result
< column "COLUMN_NAME" "VARCHAR"
< column "COLLATION_NAME" "VARCHAR"
< row
: bytes
| plain
: bytes
| latin1_swedish_ci
< row
: bytes
| named
: bytes
| utf8mb4_general_ci
< row
: bytes
| n1
: bytes
| utf8mb4_general_ci
< row
: bytes
| n2
: bytes
| utf8mb4_general_ci
< row
: bytes
| n3
: bytes
| utf8mb4_general_ci
< row
: bytes
| own
: bytes
| utf8mb4_bin
< row
: bytes
| n4
: bytes
| utf8mb4_general_ci
