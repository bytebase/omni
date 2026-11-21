-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestOracle_ViewOperatorColumnTypes default#1
> exec
| CREATE TABLE t_ops (id INT NOT NULL PRIMARY KEY, name VARCHAR(20) NOT NULL, ch CHAR(5), txt TEXT, doc JSON)
< ok 0 0
> exec
| CREATE OR REPLACE VIEW v_ops AS SELECT
| 			id & 1 AS c_band,
| 			id | 1 AS c_bor,
| 			id ^ 1 AS c_bxor,
| 			id << 1 AS c_shl,
| 			id >> 1 AS c_shr,
| 			~id AS c_bnot,
| 			id XOR 1 AS c_lxor,
| 			BINARY name AS c_bin_vc,
| 			BINARY ch AS c_bin_c,
| 			BINARY id AS c_bin_int,
| 			BINARY txt AS c_bin_txt,
| 			BINARY 'abc' AS c_bin_lit,
| 			doc->'$.x' AS c_jextract,
| 			doc->>'$.x' AS c_junquote
| 		FROM t_ops
< ok 0 0
> query
| SELECT COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE FROM information_schema.COLUMNS
| 		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'v_ops' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR"
< column "COLUMN_TYPE" "TEXT" "string"
< column "IS_NULLABLE" "VARCHAR" "string"
< row
: bytes
| c_band
: bytes
| bigint unsigned
: bytes
| NO
< row
: bytes
| c_bor
: bytes
| bigint unsigned
: bytes
| NO
< row
: bytes
| c_bxor
: bytes
| bigint unsigned
: bytes
| NO
< row
: bytes
| c_shl
: bytes
| bigint unsigned
: bytes
| NO
< row
: bytes
| c_shr
: bytes
| bigint unsigned
: bytes
| NO
< row
: bytes
| c_bnot
: bytes
| bigint unsigned
: bytes
| NO
< row
: bytes
| c_lxor
: bytes
| int
: bytes
| NO
< row
: bytes
| c_bin_vc
: bytes
| varbinary(80)
: bytes
| YES
< row
: bytes
| c_bin_c
: bytes
| varbinary(20)
: bytes
| YES
< row
: bytes
| c_bin_int
: bytes
| varbinary(11)
: bytes
| YES
< row
: bytes
| c_bin_txt
: bytes
| mediumblob
: bytes
| YES
< row
: bytes
| c_bin_lit
: bytes
| varbinary(12)
: bytes
| YES
< row
: bytes
| c_jextract
: bytes
| json
: bytes
| YES
< row
: bytes
| c_junquote
: bytes
| longtext
: bytes
| YES
