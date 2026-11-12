-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestViewExecCommentNames_OceanBase default#1
> exec
| DROP DATABASE IF EXISTS ob_view_exec_names
< ok 0 0
> exec
| CREATE DATABASE ob_view_exec_names
< ok 1 0
> exec
| USE ob_view_exec_names
< ok 0 0
> exec
| CREATE TABLE n1 (a INT, b INT)
< ok 0 0
> exec
| CREATE VIEW v0 AS SELECT a/*!50100 +2*/ FROM n1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'ob_view_exec_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| v0
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a/*!50100 +2
> exec
| CREATE VIEW v1 AS SELECT a/*!+3*/ + 1 FROM n1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'ob_view_exec_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| v1
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a/*!+3*/ + 1
> exec
| CREATE VIEW v2 AS SELECT /*!1234 */ + 1 FROM n1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'ob_view_exec_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| v2
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| 1234 */ + 1
> exec
| CREATE VIEW v3 AS SELECT 2 * /*!100005 */ FROM n1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'ob_view_exec_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| v3
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| 2 * /*!100005
> exec
| CREATE VIEW v4 AS SELECT 1 + /*!1234567 */ FROM n1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'ob_view_exec_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| v4
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| 1 + /*!1234567
> exec
| CREATE VIEW v5 AS SELECT 2 * /*!5 */ FROM n1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'ob_view_exec_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| v5
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| 2 * /*!5
> exec
| CREATE VIEW v6 AS SELECT 1 + /*!12345 2 */ FROM n1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'ob_view_exec_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| v6
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| 1 + /*!12345 2
> exec
| CREATE VIEW v7 AS SELECT 1 + /*!99999 2 */ FROM n1
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'ob_view_exec_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| v7
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| 1 + /*!99999 2
> exec
| CREATE VIEW v8 AS SELECT 1 + 1 /*! from */
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'ob_view_exec_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| v8
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| 1 + 1
> exec
| CREATE VIEW v9 AS SELECT 1 /*! + 1 */ + 1 /*! from */
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'ob_view_exec_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| v9
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| 1 /*! + 1 */ + 1
> exec
| CREATE VIEW v10 AS SELECT 1 /*! + 10 */ + 100 /*! , 2 */ /*! from */
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'ob_view_exec_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| v10
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| 1 /*! + 10 */ + 100
> exec
| CREATE VIEW v11 AS SELECT a + /*! 1 */ 2, b FROM n1 /*! where */
< ok 0 0
> query
| SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'ob_view_exec_names' AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION
: string
| v11
< result
< column "COLUMN_NAME" "VARCHAR"
< row
: bytes
| a + /*! 1 */ 2
< row
: bytes
| b
> exec
| DROP DATABASE IF EXISTS ob_view_exec_names
< ok 13 0
