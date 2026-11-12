-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestOceanBaseQuotedBuiltinCall default#1
> exec
| CREATE DATABASE IF NOT EXISTS obprobe_qb
< ok 1 0
> exec
| USE obprobe_qb
< ok 0 0
> exec
| DROP TABLE IF EXISTS tu
< ok 0 0
> exec
| CREATE TABLE tu (a DATETIME, b DATETIME)
< ok 0 0
> exec
| DROP TABLE IF EXISTS tq
< ok 0 0
> exec
| CREATE TABLE tq (a INT, b VARCHAR(20), d DATE)
< ok 0 0
> exec
| SELECT `timestampdiff`(1, 2, 3)
< ok 0 0
> exec
| SELECT `TIMESTAMPDIFF`(SECOND, a, b) FROM tu
< error
| 1054 42S22 Unknown column 'SECOND' in 'field list'
> exec
| SELECT `timestampadd`(a, 2, 3) FROM tu
< ok 0 0
> exec
| SELECT `timestampadd`(DAY, 1, a) FROM tu
< error
| 1054 42S22 Unknown column 'DAY' in 'field list'
> exec
| SELECT `get_format`(1, 2)
< ok 0 0
> exec
| SELECT `get_format`(DATE, 'USA')
< error
| 1054 42S22 Unknown column 'DATE' in 'field list'
> exec
| SELECT `timestampdiff`()
< error
| 1582 42000 Incorrect parameter count in the call to native function 'timestampdiff'
> exec
| SELECT `weight_string`('a')
< error
| 4016 HY000 Internal error
> exec
| SELECT `timestampdiff`(DAY_HOUR, a, b) FROM tu
< error
| 1064 42000 …near ', a, b) FROM tu' at line 1
> exec
| SELECT `weight_string`('a' AS CHAR(2))
< error
| 1064 42000 …near 'CHAR(2))' at line 1
> exec
| SELECT TIMESTAMPDIFF(SECOND, a, b) FROM tu
< ok 0 0
> exec
| SELECT timestampdiff(1, 2, 3)
< error
| 1064 42000 …near '1, 2, 3)' at line 1
> exec
| SELECT GET_FORMAT(1, 2)
< error
| 1064 42000 …near '1, 2)' at line 1
> exec
| SELECT WEIGHT_STRING('a' AS CHAR(2))
< ok 0 0
> exec
| SELECT `count`(a) FROM tq
< error
| 1305 42000 FUNCTION count does not exist
> exec
| SELECT `count`()
< error
| 1305 42000 FUNCTION count does not exist
> exec
| SELECT `sum`(a) FROM tq
< error
| 1305 42000 FUNCTION sum does not exist
> exec
| SELECT `trim`(b) FROM tq
< ok 0 0
> exec
| SELECT `trim`(b, 'x') FROM tq
< ok 0 0
> exec
| SELECT `trim`() FROM tq
< error
| 1582 42000 Incorrect parameter count in the call to native function 'trim'
> exec
| SELECT `substring`(b, 1, 2) FROM tq
< error
| 1305 42000 FUNCTION substring does not exist
> exec
| SELECT `substring`(b) FROM tq
< error
| 1305 42000 FUNCTION substring does not exist
> exec
| SELECT `substring`(b, 1, 2, 3) FROM tq
< error
| 1305 42000 FUNCTION substring does not exist
> exec
| SELECT `substr`(b, 1) FROM tq
< ok 0 0
> exec
| SELECT `group_concat`(b) FROM tq
< error
| 1305 42000 FUNCTION group_concat does not exist
> exec
| SELECT `group_concat`(a, b) FROM tq
< error
| 1305 42000 FUNCTION group_concat does not exist
> exec
| SELECT `rank`() FROM tq
< error
| 1305 42000 FUNCTION rank does not exist
> exec
| SELECT `row_number`() FROM tq
< error
| 1305 42000 FUNCTION row_number does not exist
> exec
| SELECT `extract`(d) FROM tq
< error
| 1582 42000 Incorrect parameter count in the call to native function 'extract'
> exec
| SELECT `position`(b) FROM tq
< error
| 1582 42000 Incorrect parameter count in the call to native function 'position'
> exec
| SELECT `char`(65)
< error
| 1582 42000 Incorrect parameter count in the call to native function 'char'
> exec
| SELECT `date_add`(d, 1) FROM tq
< error
| 1582 42000 Incorrect parameter count in the call to native function 'date_add'
> exec
| SELECT `f`(d + INTERVAL 1 DAY) FROM tq
< error
| 1305 42000 FUNCTION f does not exist
> exec
| SELECT test.count(a) FROM tq
< error
| 1305 42000 FUNCTION test.count does not exist
> exec
| SELECT test.trim(b) FROM tq
< error
| 1305 42000 FUNCTION test.trim does not exist
> exec
| SELECT test.trim(b, 'x') FROM tq
< error
| 1305 42000 FUNCTION test.trim does not exist
> exec
| SELECT test.substring(b, 1, 2) FROM tq
< error
| 1305 42000 FUNCTION test.substring does not exist
> exec
| SELECT test.substr(b, 1) FROM tq
< error
| 1305 42000 FUNCTION test.substr does not exist
> exec
| SELECT `count`(*) FROM tq
< error
| 1064 42000 …near '*) FROM tq' at line 1
> exec
| SELECT `COUNT`(*) FROM tq
< error
| 1064 42000 …near '*) FROM tq' at line 1
> exec
| SELECT `count`(*, a) FROM tq
< error
| 1064 42000 …near '*, a) FROM tq' at line 1
> exec
| SELECT `count`(DISTINCT a) FROM tq
< error
| 1064 42000 …near 'DISTINCT a) FROM tq' at line 1
> exec
| SELECT `count`(DISTINCT a, b) FROM tq
< error
| 1064 42000 …near 'DISTINCT a, b) FROM tq' at line 1
> exec
| SELECT `sum`(DISTINCT a) FROM tq
< error
| 1064 42000 …near 'DISTINCT a) FROM tq' at line 1
> exec
| SELECT `avg`(DISTINCT a) FROM tq
< error
| 1064 42000 …near 'DISTINCT a) FROM tq' at line 1
> exec
| SELECT `min`(DISTINCT a) FROM tq
< error
| 1064 42000 …near 'DISTINCT a) FROM tq' at line 1
> exec
| SELECT `max`(DISTINCT a) FROM tq
< error
| 1064 42000 …near 'DISTINCT a) FROM tq' at line 1
> exec
| SELECT `json_arrayagg`(DISTINCT a) FROM tq
< error
| 1064 42000 …near 'DISTINCT a) FROM tq' at line 1
> exec
| SELECT `group_concat`(DISTINCT b) FROM tq
< error
| 1064 42000 …near 'DISTINCT b) FROM tq' at line 1
> exec
| SELECT `group_concat`(b SEPARATOR ',') FROM tq
< error
| 1064 42000 …near 'SEPARATOR ',') FROM tq' at line 1
> exec
| SELECT `group_concat`(b ORDER BY b) FROM tq
< error
| 1064 42000 …near 'ORDER BY b) FROM tq' at line 1
> exec
| SELECT `group_concat`(b ORDER BY b SEPARATOR ',') FROM tq
< error
| 1064 42000 …near 'ORDER BY b SEPARATOR ',') FROM tq' at line 1
> exec
| SELECT `trim`(LEADING 'x' FROM b) FROM tq
< error
| 1064 42000 …near 'LEADING 'x' FROM b) FROM tq' at line 1
> exec
| SELECT `trim`(BOTH FROM b) FROM tq
< error
| 1064 42000 …near 'BOTH FROM b) FROM tq' at line 1
> exec
| SELECT `trim`('x' FROM b) FROM tq
< error
| 1064 42000 …near 'FROM b) FROM tq' at line 1
> exec
| SELECT `substring`(b FROM 1) FROM tq
< error
| 1064 42000 …near 'FROM 1) FROM tq' at line 1
> exec
| SELECT `substring`(b FROM 1 FOR 2) FROM tq
< error
| 1064 42000 …near 'FROM 1 FOR 2) FROM tq' at line 1
> exec
| SELECT `substr`(b FROM 1) FROM tq
< error
| 1064 42000 …near 'FROM 1) FROM tq' at line 1
> exec
| SELECT `date_add`(d, INTERVAL 1 DAY) FROM tq
< error
| 1064 42000 …near ') FROM tq' at line 1
> exec
| SELECT `adddate`(d, INTERVAL 1 DAY) FROM tq
< error
| 1064 42000 …near ') FROM tq' at line 1
> exec
| SELECT `sum`(a) OVER () FROM tq
< error
| 1064 42000 …near '() FROM tq' at line 1
> exec
| SELECT `rank`() OVER () FROM tq
< error
| 1064 42000 …near '() FROM tq' at line 1
> exec
| SELECT `row_number`() OVER () FROM tq
< error
| 1064 42000 …near '() FROM tq' at line 1
> exec
| SELECT `lag`(a) OVER () FROM tq
< error
| 1064 42000 …near '() FROM tq' at line 1
> exec
| SELECT `json_arrayagg`(a) OVER () FROM tq
< error
| 1064 42000 …near '() FROM tq' at line 1
> exec
| SELECT `json_objectagg`(a, b) OVER () FROM tq
< error
| 1064 42000 …near '() FROM tq' at line 1
> exec
| SELECT `count`(a) OVER w FROM tq WINDOW w AS ()
< error
| 1064 42000 …near 'w FROM tq WINDOW w AS ()' at line 1
> exec
| SELECT `count`(ALL a) FROM tq
< error
| 1064 42000 …near 'ALL a) FROM tq' at line 1
> exec
| SELECT `cast`(a AS CHAR) FROM tq
< error
| 1064 42000 …near 'CHAR) FROM tq' at line 1
> exec
| SELECT `convert`(b USING utf8mb4) FROM tq
< error
| 1064 42000 …near 'USING utf8mb4) FROM tq' at line 1
> exec
| SELECT `convert`(b, CHAR) FROM tq
< error
| 1064 42000 …near ') FROM tq' at line 1
> exec
| SELECT `extract`(YEAR FROM d) FROM tq
< error
| 1064 42000 …near 'FROM d) FROM tq' at line 1
> exec
| SELECT `position`('a' IN b) FROM tq
< error
| 1064 42000 …near 'b) FROM tq' at line 1
> exec
| SELECT `char`(65 USING utf8mb4)
< error
| 1064 42000 …near 'USING utf8mb4)' at line 1
> exec
| SELECT * FROM tq WHERE `match`(b) AGAINST ('x')
< error
| 1064 42000 …near 'AGAINST ('x')' at line 1
> exec
| SELECT * FROM `json_table`('[1]', '$[*]' COLUMNS (x INT PATH '$')) AS jt
< error
| 1064 42000 …near '('[1]', '$[*]' COLUMNS (x INT PATH '$')) AS jt' at line 1
> exec
| SELECT test.count(*) FROM tq
< error
| 1064 42000 …near '*) FROM tq' at line 1
> exec
| SELECT test.count(DISTINCT a) FROM tq
< error
| 1064 42000 …near 'DISTINCT a) FROM tq' at line 1
> exec
| SELECT test.sum(DISTINCT a) FROM tq
< error
| 1064 42000 …near 'DISTINCT a) FROM tq' at line 1
> exec
| SELECT test.trim(LEADING 'x' FROM b) FROM tq
< error
| 1064 42000 …near 'LEADING 'x' FROM b) FROM tq' at line 1
> exec
| SELECT test.trim(TRAILING 'x' FROM b) FROM tq
< error
| 1064 42000 …near 'TRAILING 'x' FROM b) FROM tq' at line 1
> exec
| SELECT test.trim('x' FROM b) FROM tq
< error
| 1064 42000 …near 'FROM b) FROM tq' at line 1
> exec
| SELECT `test`.trim(LEADING 'x' FROM b) FROM tq
< error
| 1064 42000 …near '(LEADING 'x' FROM b) FROM tq' at line 1
> exec
| SELECT test.`trim`(LEADING 'x' FROM b) FROM tq
< error
| 1064 42000 …near 'LEADING 'x' FROM b) FROM tq' at line 1
> exec
| SELECT test.substring(b FROM 1) FROM tq
< error
| 1064 42000 …near 'FROM 1) FROM tq' at line 1
> exec
| SELECT test.substring(b FROM 1 FOR 2) FROM tq
< error
| 1064 42000 …near 'FROM 1 FOR 2) FROM tq' at line 1
> exec
| SELECT test.substr(b FROM 1) FROM tq
< error
| 1064 42000 …near 'FROM 1) FROM tq' at line 1
> exec
| SELECT test.group_concat(b SEPARATOR ',') FROM tq
< error
| 1064 42000 …near 'SEPARATOR ',') FROM tq' at line 1
> exec
| SELECT test.group_concat(b ORDER BY b) FROM tq
< error
| 1064 42000 …near 'ORDER BY b) FROM tq' at line 1
> exec
| SELECT test.group_concat(DISTINCT b) FROM tq
< error
| 1064 42000 …near 'DISTINCT b) FROM tq' at line 1
> exec
| SELECT test.date_add(d, INTERVAL 1 DAY) FROM tq
< error
| 1064 42000 …near ') FROM tq' at line 1
> exec
| SELECT test.sum(a) OVER () FROM tq
< error
| 1064 42000 …near '() FROM tq' at line 1
> exec
| SELECT test.count(a) OVER () FROM tq
< error
| 1064 42000 …near '() FROM tq' at line 1
> exec
| SELECT COUNT(*) FROM tq
< ok 0 0
> exec
| SELECT COUNT(DISTINCT a) FROM tq
< ok 0 0
> exec
| SELECT TRIM(LEADING 'x' FROM b) FROM tq
< ok 0 0
> exec
| SELECT SUBSTRING(b FROM 1 FOR 2) FROM tq
< ok 0 0
> exec
| SELECT GROUP_CONCAT(b ORDER BY b SEPARATOR ',') FROM tq
< ok 0 0
> exec
| SELECT SUM(a) OVER () FROM tq
< ok 0 0
> exec
| SELECT DATE_ADD(d, INTERVAL 1 DAY) FROM tq
< ok 0 0
> exec
| SELECT JSON_ARRAYAGG(a) OVER () FROM tq
< ok 0 0
> exec
| DROP DATABASE IF EXISTS obprobe_qb
< ok 2 0
