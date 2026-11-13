-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestScenario_C23 default#1
> exec
| DROP DATABASE IF EXISTS testdb
< ok 0 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a VARCHAR(10), b VARCHAR(10))
< ok 0 0
> exec
| CREATE VIEW v_concat AS SELECT CONCAT(a, b) AS c1, CONCAT('x','','y') AS c2 FROM t
< ok 0 0
> exec
| INSERT INTO t VALUES ('foo', NULL), ('foo', 'bar')
< ok 2 0
> query
| SELECT CONCAT(a, b) FROM t ORDER BY b IS NULL DESC
< result
< column "CONCAT(a, b)" "VARCHAR"
< row
: null
< row
: bytes
| foobar
> query
| SELECT CONCAT('x', NULL, 'y')
< result
< column "CONCAT('x', NULL, 'y')" "VARCHAR"
< row
: null
> query
| SELECT CONCAT('x', '', 'y')
< result
< column "CONCAT('x', '', 'y')" "VARCHAR"
< row
: bytes
| xy
> query
| SELECT IS_NULLABLE FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v_concat' AND COLUMN_NAME='c1'
< result
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| YES
> query
| SELECT IS_NULLABLE FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v_concat' AND COLUMN_NAME='c2'
< result
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| YES
> exec
| DROP DATABASE IF EXISTS testdb
< ok 2 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (a VARCHAR(10), b VARCHAR(10), c VARCHAR(10))
< ok 0 0
> exec
| CREATE VIEW v_ws AS SELECT CONCAT_WS(',', a, b, c) AS d1 FROM t
< ok 0 0
> exec
| INSERT INTO t VALUES ('x', NULL, 'z'), (NULL, NULL, NULL)
< ok 2 0
> query
| SELECT CONCAT_WS(',', a, b, c) FROM t ORDER BY a IS NULL, a
< result
< column "CONCAT_WS(',', a, b, c)" "VARCHAR"
< row
: bytes
| x,z
< row
: bytes
|
> query
| SELECT CONCAT_WS(',', 'x', NULL, 'z')
< result
< column "CONCAT_WS(',', 'x', NULL, 'z')" "VARCHAR"
< row
: bytes
| x,z
> query
| SELECT CONCAT_WS(',', NULL, NULL, NULL)
< result
< column "CONCAT_WS(',', NULL, NULL, NULL)" "VARCHAR"
< row
: bytes
|
> query
| SELECT CONCAT_WS(NULL, 'x', 'y')
< result
< column "CONCAT_WS(NULL, 'x', 'y')" "VARCHAR"
< row
: null
> query
| SELECT IS_NULLABLE FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v_ws' AND COLUMN_NAME='d1'
< result
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| YES
> exec
| DROP DATABASE IF EXISTS testdb
< ok 2 0
> exec
| CREATE DATABASE testdb
< ok 1 0
> exec
| USE testdb
< ok 0 0
> exec
| CREATE TABLE t (first_name VARCHAR(20), middle_name VARCHAR(20), last_name VARCHAR(20))
< ok 0 0
> exec
| CREATE VIEW v_name AS SELECT
| 			   CONCAT(first_name, ' ', middle_name, ' ', last_name) AS bad,
| 			   CONCAT(first_name, ' ', IFNULL(middle_name, ''), ' ', last_name) AS rescue_ifnull,
| 			   CONCAT(first_name, ' ', COALESCE(middle_name, ''), ' ', last_name) AS rescue_coalesce,
| 			   CONCAT_WS(' ', first_name, middle_name, last_name) AS rescue_ws
| 			 FROM t
< ok 0 0
> exec
| INSERT INTO t VALUES ('Ada', NULL, 'Lovelace')
< ok 1 0
> query
| SELECT bad, rescue_ifnull, rescue_coalesce, rescue_ws FROM v_name
< result
< column "bad" "VARCHAR"
< column "rescue_ifnull" "VARCHAR"
< column "rescue_coalesce" "VARCHAR"
< column "rescue_ws" "VARCHAR"
< row
: null
: bytes
| Ada  Lovelace
: bytes
| Ada  Lovelace
: bytes
| Ada Lovelace
> query
| SELECT IS_NULLABLE FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v_name' AND COLUMN_NAME='bad'
< result
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| YES
> query
| SELECT IS_NULLABLE FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v_name' AND COLUMN_NAME='rescue_ifnull'
< result
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| YES
> query
| SELECT IS_NULLABLE FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v_name' AND COLUMN_NAME='rescue_coalesce'
< result
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| YES
> query
| SELECT IS_NULLABLE FROM information_schema.COLUMNS
| 			 WHERE TABLE_SCHEMA='testdb' AND TABLE_NAME='v_name' AND COLUMN_NAME='rescue_ws'
< result
< column "IS_NULLABLE" "VARCHAR"
< row
: bytes
| YES
