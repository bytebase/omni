-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseViewOperatorTypesAndMVParenIndex default#1
> exec
| CREATE TABLE t_ops (id INT NOT NULL PRIMARY KEY, name VARCHAR(20) NOT NULL, ch CHAR(5), txt TEXT, doc JSON, `a(b` INT)
< ok 0 0
> exec
| CREATE VIEW v_ops AS SELECT
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
> exec
| CREATE MATERIALIZED VIEW mv_paren (PRIMARY KEY (id)) AS SELECT id, `a(b` FROM t_ops
< ok 0 0
> exec
| CREATE INDEX idx_paren ON mv_paren (`a(b`)
< ok 0 0
> query
| SELECT QUERY FROM oceanbase.DBA_MVIEWS WHERE OWNER = DATABASE() AND MVIEW_NAME = 'mv_paren'
< result
< column "QUERY" "TEXT" "string"
< row
: bytes
| select `omni_t_oceanbase_catalog_view_ops_and_mv`.`t_ops`.`id` AS `id`,`omni_t_oceanbase_catalog_view_ops_and_mv`.`t_ops`.`a(b` AS `a(b` from `omni_t_oceanbase_catalog_view_ops_and_mv`.`t_ops`
> query
| SELECT s.COLUMN_NAME FROM oceanbase.DBA_MVIEWS m
| 		JOIN information_schema.STATISTICS s ON s.TABLE_SCHEMA = m.OWNER AND s.TABLE_NAME = m.CONTAINER_NAME
| 		WHERE m.OWNER = DATABASE() AND m.MVIEW_NAME = 'mv_paren' AND s.INDEX_NAME = 'idx_paren'
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a(b
> query
| SELECT COLUMN_NAME, COLUMN_TYPE FROM information_schema.COLUMNS
| 		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'v_ops' ORDER BY ORDINAL_POSITION
< result
< column "COLUMN_NAME" "VARCHAR" "string"
< column "COLUMN_TYPE" "TEXT"
< row
: bytes
| c_band
: bytes
| bigint(20) unsigned
< row
: bytes
| c_bor
: bytes
| bigint(20) unsigned
< row
: bytes
| c_bxor
: bytes
| bigint(20) unsigned
< row
: bytes
| c_shl
: bytes
| bigint(20) unsigned
< row
: bytes
| c_shr
: bytes
| bigint(20) unsigned
< row
: bytes
| c_bnot
: bytes
| bigint(20) unsigned
< row
: bytes
| c_lxor
: bytes
| int(1)
< row
: bytes
| c_bin_vc
: bytes
| varbinary(20)
< row
: bytes
| c_bin_c
: bytes
| varbinary(5)
< row
: bytes
| c_bin_int
: bytes
| varbinary(11)
< row
: bytes
| c_bin_txt
: bytes
| varbinary(65535)
< row
: bytes
| c_bin_lit
: bytes
| varbinary(3)
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
