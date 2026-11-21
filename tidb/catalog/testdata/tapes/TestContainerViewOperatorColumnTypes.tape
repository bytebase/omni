-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestContainerViewOperatorColumnTypes default#1
> exec
| DROP DATABASE IF EXISTS vops_oracle
< ok 0 0
> exec
| CREATE DATABASE vops_oracle; USE vops_oracle
< ok 0 0
> exec
| CREATE TABLE t_ops (id INT NOT NULL PRIMARY KEY, name VARCHAR(20) NOT NULL, ch CHAR(5), doc JSON);CREATE VIEW v_ops AS SELECT
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
| 			BINARY 'abc' AS c_bin_lit,
| 			doc->'$.x' AS c_jextract,
| 			doc->>'$.x' AS c_junquote
| 		FROM t_ops;
< ok 0 0
> query
| SELECT COLUMN_NAME, COLUMN_TYPE FROM information_schema.COLUMNS
| 		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 'v_ops' ORDER BY ORDINAL_POSITION
: string
| vops_oracle
< result
< column "COLUMN_NAME" "VARCHAR"
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| c_band
: bytes
| bigint unsigned
< row
: bytes
| c_bor
: bytes
| bigint unsigned
< row
: bytes
| c_bxor
: bytes
| bigint unsigned
< row
: bytes
| c_shl
: bytes
| bigint unsigned
< row
: bytes
| c_shr
: bytes
| bigint unsigned
< row
: bytes
| c_bnot
: bytes
| bigint unsigned
< row
: bytes
| c_lxor
: bytes
| bigint
< row
: bytes
| c_bin_vc
: bytes
| binary(1)
< row
: bytes
| c_bin_c
: bytes
| binary(1)
< row
: bytes
| c_bin_int
: bytes
| binary(1)
< row
: bytes
| c_bin_lit
: bytes
| binary(1)
< row
: bytes
| c_jextract
: bytes
| json
< row
: bytes
| c_junquote
: bytes
| longtext
> exec
| DROP DATABASE IF EXISTS vops_oracle
< ok 0 0
