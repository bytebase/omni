-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOracle_OceanBaseViewAttributes default#1
> exec
| DROP DATABASE IF EXISTS obviews
< ok 0 0
> exec
| CREATE DATABASE obviews
< ok 1 0
> exec
| USE obviews
< ok 0 0
> exec
| CREATE TABLE t (id INT NOT NULL, a INT NOT NULL, n INT)
< ok 0 0
> exec
| CREATE TABLE u (id INT NOT NULL, b INT NOT NULL)
< ok 0 0
> exec
| CREATE TABLE s (l VARCHAR(10) CHARACTER SET latin1 COLLATE latin1_swedish_ci, m VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin)
< ok 0 0
> exec
| CREATE VIEW v_plain AS SELECT a, n FROM t
< ok 0 0
> query
| SELECT IS_UPDATABLE, SECURITY_TYPE FROM information_schema.VIEWS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
: string
| obviews
: string
| v_plain
< result
< column "IS_UPDATABLE" "VARCHAR"
< column "SECURITY_TYPE" "VARCHAR"
< row
: bytes
| YES
: bytes
| NONE
> query
| SELECT COLUMN_NAME, IS_NULLABLE, CHARACTER_SET_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| obviews
: string
| v_plain
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IS_NULLABLE" "VARCHAR"
< column "CHARACTER_SET_NAME" "VARCHAR"
< row
: bytes
| a
: bytes
| NO
: null
< row
: bytes
| n
: bytes
| YES
: null
> exec
| CREATE VIEW v_distinct AS SELECT DISTINCT a FROM t
< ok 0 0
> query
| SELECT IS_UPDATABLE, SECURITY_TYPE FROM information_schema.VIEWS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
: string
| obviews
: string
| v_distinct
< result
< column "IS_UPDATABLE" "VARCHAR"
< column "SECURITY_TYPE" "VARCHAR"
< row
: bytes
| YES
: bytes
| NONE
> query
| SELECT COLUMN_NAME, IS_NULLABLE, CHARACTER_SET_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| obviews
: string
| v_distinct
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IS_NULLABLE" "VARCHAR"
< column "CHARACTER_SET_NAME" "VARCHAR"
< row
: bytes
| a
: bytes
| NO
: null
> exec
| CREATE ALGORITHM=TEMPTABLE VIEW v_temptable AS SELECT a FROM t
< ok 0 0
> query
| SELECT IS_UPDATABLE, SECURITY_TYPE FROM information_schema.VIEWS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
: string
| obviews
: string
| v_temptable
< result
< column "IS_UPDATABLE" "VARCHAR"
< column "SECURITY_TYPE" "VARCHAR"
< row
: bytes
| YES
: bytes
| NONE
> query
| SELECT COLUMN_NAME, IS_NULLABLE, CHARACTER_SET_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| obviews
: string
| v_temptable
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IS_NULLABLE" "VARCHAR"
< column "CHARACTER_SET_NAME" "VARCHAR"
< row
: bytes
| a
: bytes
| NO
: null
> exec
| CREATE ALGORITHM=MERGE VIEW v_merge AS SELECT a FROM t
< ok 0 0
> query
| SELECT IS_UPDATABLE, SECURITY_TYPE FROM information_schema.VIEWS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
: string
| obviews
: string
| v_merge
< result
< column "IS_UPDATABLE" "VARCHAR"
< column "SECURITY_TYPE" "VARCHAR"
< row
: bytes
| YES
: bytes
| NONE
> query
| SELECT COLUMN_NAME, IS_NULLABLE, CHARACTER_SET_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| obviews
: string
| v_merge
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IS_NULLABLE" "VARCHAR"
< column "CHARACTER_SET_NAME" "VARCHAR"
< row
: bytes
| a
: bytes
| NO
: null
> exec
| CREATE VIEW v_group AS SELECT a, COUNT(*) AS c FROM t GROUP BY a
< ok 0 0
> query
| SELECT IS_UPDATABLE, SECURITY_TYPE FROM information_schema.VIEWS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
: string
| obviews
: string
| v_group
< result
< column "IS_UPDATABLE" "VARCHAR"
< column "SECURITY_TYPE" "VARCHAR"
< row
: bytes
| YES
: bytes
| NONE
> query
| SELECT COLUMN_NAME, IS_NULLABLE, CHARACTER_SET_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| obviews
: string
| v_group
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IS_NULLABLE" "VARCHAR"
< column "CHARACTER_SET_NAME" "VARCHAR"
< row
: bytes
| a
: bytes
| NO
: null
< row
: bytes
| c
: bytes
| YES
: null
> exec
| CREATE VIEW v_expr AS SELECT a + 1 AS e, 1 AS one, 'x' AS str, ABS(a) AS f, CAST(a AS SIGNED) AS g, id FROM t
< ok 0 0
> query
| SELECT IS_UPDATABLE, SECURITY_TYPE FROM information_schema.VIEWS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
: string
| obviews
: string
| v_expr
< result
< column "IS_UPDATABLE" "VARCHAR"
< column "SECURITY_TYPE" "VARCHAR"
< row
: bytes
| YES
: bytes
| NONE
> query
| SELECT COLUMN_NAME, IS_NULLABLE, CHARACTER_SET_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| obviews
: string
| v_expr
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IS_NULLABLE" "VARCHAR"
< column "CHARACTER_SET_NAME" "VARCHAR"
< row
: bytes
| e
: bytes
| NO
: null
< row
: bytes
| one
: bytes
| NO
: null
< row
: bytes
| str
: bytes
| NO
: bytes
| utf8mb4
< row
: bytes
| f
: bytes
| NO
: null
< row
: bytes
| g
: bytes
| NO
: null
< row
: bytes
| id
: bytes
| NO
: null
> exec
| CREATE VIEW v_inner AS SELECT t.a, u.b, t.n FROM t JOIN u ON t.id = u.id
< ok 0 0
> query
| SELECT IS_UPDATABLE, SECURITY_TYPE FROM information_schema.VIEWS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
: string
| obviews
: string
| v_inner
< result
< column "IS_UPDATABLE" "VARCHAR"
< column "SECURITY_TYPE" "VARCHAR"
< row
: bytes
| YES
: bytes
| NONE
> query
| SELECT COLUMN_NAME, IS_NULLABLE, CHARACTER_SET_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| obviews
: string
| v_inner
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IS_NULLABLE" "VARCHAR"
< column "CHARACTER_SET_NAME" "VARCHAR"
< row
: bytes
| a
: bytes
| NO
: null
< row
: bytes
| b
: bytes
| NO
: null
< row
: bytes
| n
: bytes
| YES
: null
> exec
| CREATE VIEW v_union AS SELECT a FROM t UNION SELECT b FROM u
< ok 0 0
> query
| SELECT IS_UPDATABLE, SECURITY_TYPE FROM information_schema.VIEWS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
: string
| obviews
: string
| v_union
< result
< column "IS_UPDATABLE" "VARCHAR"
< column "SECURITY_TYPE" "VARCHAR"
< row
: bytes
| YES
: bytes
| NONE
> query
| SELECT COLUMN_NAME, IS_NULLABLE, CHARACTER_SET_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| obviews
: string
| v_union
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IS_NULLABLE" "VARCHAR"
< column "CHARACTER_SET_NAME" "VARCHAR"
< row
: bytes
| a
: bytes
| NO
: null
> exec
| CREATE SQL SECURITY INVOKER VIEW v_invoker AS SELECT a FROM t
< ok 0 0
> query
| SELECT IS_UPDATABLE, SECURITY_TYPE FROM information_schema.VIEWS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
: string
| obviews
: string
| v_invoker
< result
< column "IS_UPDATABLE" "VARCHAR"
< column "SECURITY_TYPE" "VARCHAR"
< row
: bytes
| YES
: bytes
| NONE
> query
| SELECT COLUMN_NAME, IS_NULLABLE, CHARACTER_SET_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| obviews
: string
| v_invoker
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IS_NULLABLE" "VARCHAR"
< column "CHARACTER_SET_NAME" "VARCHAR"
< row
: bytes
| a
: bytes
| NO
: null
> exec
| CREATE SQL SECURITY DEFINER VIEW v_definer AS SELECT a FROM t
< ok 0 0
> query
| SELECT IS_UPDATABLE, SECURITY_TYPE FROM information_schema.VIEWS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
: string
| obviews
: string
| v_definer
< result
< column "IS_UPDATABLE" "VARCHAR"
< column "SECURITY_TYPE" "VARCHAR"
< row
: bytes
| YES
: bytes
| NONE
> query
| SELECT COLUMN_NAME, IS_NULLABLE, CHARACTER_SET_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| obviews
: string
| v_definer
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IS_NULLABLE" "VARCHAR"
< column "CHARACTER_SET_NAME" "VARCHAR"
< row
: bytes
| a
: bytes
| NO
: null
> exec
| CREATE VIEW v_left AS SELECT t.a, u.b FROM t LEFT JOIN u ON t.id = u.id
< ok 0 0
> query
| SELECT IS_UPDATABLE, SECURITY_TYPE FROM information_schema.VIEWS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
: string
| obviews
: string
| v_left
< result
< column "IS_UPDATABLE" "VARCHAR"
< column "SECURITY_TYPE" "VARCHAR"
< row
: bytes
| YES
: bytes
| NONE
> query
| SELECT COLUMN_NAME, IS_NULLABLE, CHARACTER_SET_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| obviews
: string
| v_left
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IS_NULLABLE" "VARCHAR"
< column "CHARACTER_SET_NAME" "VARCHAR"
< row
: bytes
| a
: bytes
| NO
: null
< row
: bytes
| b
: bytes
| NO
: null
> exec
| CREATE VIEW v_right AS SELECT t.a, u.b FROM t RIGHT JOIN u ON t.id = u.id
< ok 0 0
> query
| SELECT IS_UPDATABLE, SECURITY_TYPE FROM information_schema.VIEWS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
: string
| obviews
: string
| v_right
< result
< column "IS_UPDATABLE" "VARCHAR"
< column "SECURITY_TYPE" "VARCHAR"
< row
: bytes
| YES
: bytes
| NONE
> query
| SELECT COLUMN_NAME, IS_NULLABLE, CHARACTER_SET_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| obviews
: string
| v_right
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IS_NULLABLE" "VARCHAR"
< column "CHARACTER_SET_NAME" "VARCHAR"
< row
: bytes
| a
: bytes
| NO
: null
< row
: bytes
| b
: bytes
| NO
: null
> exec
| CREATE VIEW v_repeat AS SELECT REPEAT(l, 3) AS c FROM s
< ok 0 0
> query
| SELECT IS_UPDATABLE, SECURITY_TYPE FROM information_schema.VIEWS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
: string
| obviews
: string
| v_repeat
< result
< column "IS_UPDATABLE" "VARCHAR"
< column "SECURITY_TYPE" "VARCHAR"
< row
: bytes
| YES
: bytes
| NONE
> query
| SELECT COLUMN_NAME, IS_NULLABLE, CHARACTER_SET_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| obviews
: string
| v_repeat
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IS_NULLABLE" "VARCHAR"
< column "CHARACTER_SET_NAME" "VARCHAR"
< row
: bytes
| c
: bytes
| YES
: bytes
| latin1
> exec
| CREATE VIEW v_lpad_lm AS SELECT LPAD(l, 20, m) AS c FROM s
< ok 0 0
> query
| SELECT IS_UPDATABLE, SECURITY_TYPE FROM information_schema.VIEWS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
: string
| obviews
: string
| v_lpad_lm
< result
< column "IS_UPDATABLE" "VARCHAR"
< column "SECURITY_TYPE" "VARCHAR"
< row
: bytes
| YES
: bytes
| NONE
> query
| SELECT COLUMN_NAME, IS_NULLABLE, CHARACTER_SET_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| obviews
: string
| v_lpad_lm
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IS_NULLABLE" "VARCHAR"
< column "CHARACTER_SET_NAME" "VARCHAR"
< row
: bytes
| c
: bytes
| YES
: bytes
| utf8mb4
> exec
| CREATE VIEW v_lpad_ml AS SELECT LPAD(m, 20, l) AS c FROM s
< ok 0 0
> query
| SELECT IS_UPDATABLE, SECURITY_TYPE FROM information_schema.VIEWS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
: string
| obviews
: string
| v_lpad_ml
< result
< column "IS_UPDATABLE" "VARCHAR"
< column "SECURITY_TYPE" "VARCHAR"
< row
: bytes
| YES
: bytes
| NONE
> query
| SELECT COLUMN_NAME, IS_NULLABLE, CHARACTER_SET_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| obviews
: string
| v_lpad_ml
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IS_NULLABLE" "VARCHAR"
< column "CHARACTER_SET_NAME" "VARCHAR"
< row
: bytes
| c
: bytes
| YES
: bytes
| utf8mb4
> exec
| CREATE VIEW v_rpad_lm AS SELECT RPAD(l, 20, m) AS c FROM s
< ok 0 0
> query
| SELECT IS_UPDATABLE, SECURITY_TYPE FROM information_schema.VIEWS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
: string
| obviews
: string
| v_rpad_lm
< result
< column "IS_UPDATABLE" "VARCHAR"
< column "SECURITY_TYPE" "VARCHAR"
< row
: bytes
| YES
: bytes
| NONE
> query
| SELECT COLUMN_NAME, IS_NULLABLE, CHARACTER_SET_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| obviews
: string
| v_rpad_lm
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IS_NULLABLE" "VARCHAR"
< column "CHARACTER_SET_NAME" "VARCHAR"
< row
: bytes
| c
: bytes
| YES
: bytes
| utf8mb4
> exec
| CREATE VIEW v_rpad_ml AS SELECT RPAD(m, 20, l) AS c FROM s
< ok 0 0
> query
| SELECT IS_UPDATABLE, SECURITY_TYPE FROM information_schema.VIEWS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
: string
| obviews
: string
| v_rpad_ml
< result
< column "IS_UPDATABLE" "VARCHAR"
< column "SECURITY_TYPE" "VARCHAR"
< row
: bytes
| YES
: bytes
| NONE
> query
| SELECT COLUMN_NAME, IS_NULLABLE, CHARACTER_SET_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| obviews
: string
| v_rpad_ml
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IS_NULLABLE" "VARCHAR"
< column "CHARACTER_SET_NAME" "VARCHAR"
< row
: bytes
| c
: bytes
| YES
: bytes
| utf8mb4
> exec
| CREATE VIEW v_lpad_ll AS SELECT LPAD(l, 20, l) AS c FROM s
< ok 0 0
> query
| SELECT IS_UPDATABLE, SECURITY_TYPE FROM information_schema.VIEWS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
: string
| obviews
: string
| v_lpad_ll
< result
< column "IS_UPDATABLE" "VARCHAR"
< column "SECURITY_TYPE" "VARCHAR"
< row
: bytes
| YES
: bytes
| NONE
> query
| SELECT COLUMN_NAME, IS_NULLABLE, CHARACTER_SET_NAME FROM information_schema.COLUMNS
| 				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| obviews
: string
| v_lpad_ll
< result
< column "COLUMN_NAME" "VARCHAR"
< column "IS_NULLABLE" "VARCHAR"
< column "CHARACTER_SET_NAME" "VARCHAR"
< row
: bytes
| c
: bytes
| YES
: bytes
| latin1
