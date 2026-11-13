-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestQuotedBuiltinCallOracle default#1
> exec
| CREATE TABLE tq (a INT, b VARCHAR(20), d DATE)
< ok 0 0
> exec
| CREATE FULLTEXT INDEX ft ON tq (b)
< ok 0 0
> exec
| SELECT `count`(a) FROM tq
< error
| 1630 42000 FUNCTION test.count does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT `count`(a) FROM tq
< error
| 1630 42000 FUNCTION test.count does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT `count`(a) FROM tq
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT `count`()
< error
| 1630 42000 FUNCTION test.count does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT `count`()
< error
| 1630 42000 FUNCTION test.count does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT `count`()
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT `sum`(a) FROM tq
< error
| 1630 42000 FUNCTION test.sum does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT `sum`(a) FROM tq
< error
| 1630 42000 FUNCTION test.sum does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT `sum`(a) FROM tq
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT `trim`(b) FROM tq
< error
| 1630 42000 FUNCTION test.trim does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT `trim`(b) FROM tq
< error
| 1630 42000 FUNCTION test.trim does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT `trim`(b) FROM tq
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT `trim`(b, 'x') FROM tq
< error
| 1630 42000 FUNCTION test.trim does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT `trim`(b, 'x') FROM tq
< error
| 1630 42000 FUNCTION test.trim does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT `trim`(b, 'x') FROM tq
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT `trim`() FROM tq
< error
| 1630 42000 FUNCTION test.trim does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT `trim`() FROM tq
< error
| 1630 42000 FUNCTION test.trim does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT `trim`() FROM tq
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT `substring`(b, 1, 2) FROM tq
< error
| 1630 42000 FUNCTION test.substring does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT `substring`(b, 1, 2) FROM tq
< error
| 1630 42000 FUNCTION test.substring does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT `substring`(b, 1, 2) FROM tq
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT `substring`(b) FROM tq
< error
| 1630 42000 FUNCTION test.substring does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT `substring`(b) FROM tq
< error
| 1630 42000 FUNCTION test.substring does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT `substring`(b) FROM tq
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT `substring`(b, 1, 2, 3) FROM tq
< error
| 1630 42000 FUNCTION test.substring does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT `substring`(b, 1, 2, 3) FROM tq
< error
| 1630 42000 FUNCTION test.substring does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT `substring`(b, 1, 2, 3) FROM tq
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT `substr`(b, 1) FROM tq
< error
| 1630 42000 FUNCTION test.substr does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT `substr`(b, 1) FROM tq
< error
| 1630 42000 FUNCTION test.substr does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT `substr`(b, 1) FROM tq
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT `group_concat`(b) FROM tq
< error
| 1630 42000 FUNCTION test.group_concat does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT `group_concat`(b) FROM tq
< error
| 1630 42000 FUNCTION test.group_concat does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT `group_concat`(b) FROM tq
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT `group_concat`(a, b) FROM tq
< error
| 1630 42000 FUNCTION test.group_concat does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT `group_concat`(a, b) FROM tq
< error
| 1630 42000 FUNCTION test.group_concat does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT `group_concat`(a, b) FROM tq
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT `rank`() FROM tq
< error
| 1630 42000 FUNCTION test.rank does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT `rank`() FROM tq
< error
| 1630 42000 FUNCTION test.rank does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT `rank`() FROM tq
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT `row_number`() FROM tq
< error
| 1630 42000 FUNCTION test.row_number does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT `row_number`() FROM tq
< error
| 1630 42000 FUNCTION test.row_number does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT `row_number`() FROM tq
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT `cast`(a) FROM tq
< error
| 1630 42000 FUNCTION test.cast does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT `cast`(a) FROM tq
< error
| 1630 42000 FUNCTION test.cast does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT `cast`(a) FROM tq
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT `extract`(d) FROM tq
< error
| 1630 42000 FUNCTION test.extract does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT `extract`(d) FROM tq
< error
| 1630 42000 FUNCTION test.extract does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT `extract`(d) FROM tq
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT `position`(b) FROM tq
< error
| 1630 42000 FUNCTION test.position does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT `position`(b) FROM tq
< error
| 1630 42000 FUNCTION test.position does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT `position`(b) FROM tq
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT `char`(65)
< error
| 1630 42000 FUNCTION test.char does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT `char`(65)
< error
| 1630 42000 FUNCTION test.char does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT `char`(65)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT `date_add`(d, 1) FROM tq
< error
| 1630 42000 FUNCTION test.date_add does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT `date_add`(d, 1) FROM tq
< error
| 1630 42000 FUNCTION test.date_add does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT `date_add`(d, 1) FROM tq
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT `f`(d + INTERVAL 1 DAY) FROM tq
< error
| 1305 42000 FUNCTION test.f does not exist
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT `f`(d + INTERVAL 1 DAY) FROM tq
< error
| 1305 42000 FUNCTION test.f does not exist
> exec ^
: string
| SELECT `f`(d + INTERVAL 1 DAY) FROM tq
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT test.count(a) FROM tq
< error
| 1630 42000 FUNCTION test.count does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT test.count(a) FROM tq
< error
| 1630 42000 FUNCTION test.count does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT test.count(a) FROM tq
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT test.replace(b, 'x', 'y') FROM tq
< error
| 1630 42000 FUNCTION test.replace does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT test.replace(b, 'x', 'y') FROM tq
< error
| 1630 42000 FUNCTION test.replace does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT test.replace(b, 'x', 'y') FROM tq
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
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
| SELECT `substring`(DISTINCT b) FROM tq
< error
| 1064 42000 …near 'DISTINCT b) FROM tq' at line 1
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
| 1064 42000 …near 'OVER () FROM tq' at line 1
> exec
| SELECT `rank`() OVER () FROM tq
< error
| 1064 42000 …near 'OVER () FROM tq' at line 1
> exec
| SELECT `row_number`() OVER () FROM tq
< error
| 1064 42000 …near 'OVER () FROM tq' at line 1
> exec
| SELECT `lag`(a) OVER () FROM tq
< error
| 1064 42000 …near 'OVER () FROM tq' at line 1
> exec
| SELECT `json_arrayagg`(a) OVER () FROM tq
< error
| 1064 42000 …near 'OVER () FROM tq' at line 1
> exec
| SELECT `json_objectagg`(a, b) OVER () FROM tq
< error
| 1064 42000 …near 'OVER () FROM tq' at line 1
> exec
| SELECT `count`(a) OVER w FROM tq WINDOW w AS ()
< error
| 1064 42000 …near 'OVER w FROM tq WINDOW w AS ()' at line 1
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
| 1064 42000 …near 'OVER () FROM tq' at line 1
> exec
| SELECT test.count(a) OVER () FROM tq
< error
| 1064 42000 …near 'OVER () FROM tq' at line 1
> exec
| SELECT test.trim(LEADING 'x' FROM b) FROM tq
< error
| 1128 HY000 Function 'test.trim' is not defined
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT test.trim(LEADING 'x' FROM b) FROM tq
< error
| 1128 HY000 Function 'test.trim' is not defined
> exec
| SELECT test.TRIM(TRAILING 'x' FROM b) FROM tq
< error
| 1128 HY000 Function 'test.TRIM' is not defined
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT test.TRIM(TRAILING 'x' FROM b) FROM tq
< error
| 1128 HY000 Function 'test.TRIM' is not defined
> exec
| SELECT test.trim('x' FROM b) FROM tq
< error
| 1128 HY000 Function 'test.trim' is not defined
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT test.trim('x' FROM b) FROM tq
< error
| 1128 HY000 Function 'test.trim' is not defined
> exec
| SELECT `test`.trim(LEADING 'x' FROM b) FROM tq
< error
| 1128 HY000 Function 'test.trim' is not defined
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT `test`.trim(LEADING 'x' FROM b) FROM tq
< error
| 1128 HY000 Function 'test.trim' is not defined
> exec
| SELECT test.substring(b FROM 1) FROM tq
< error
| 1128 HY000 Function 'test.substring' is not defined
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT test.substring(b FROM 1) FROM tq
< error
| 1128 HY000 Function 'test.substring' is not defined
> exec
| SELECT test.Substring(b FROM 1 FOR 2) FROM tq
< error
| 1128 HY000 Function 'test.Substring' is not defined
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT test.Substring(b FROM 1 FOR 2) FROM tq
< error
| 1128 HY000 Function 'test.Substring' is not defined
> exec
| SELECT test.substr(b FROM 1) FROM tq
< error
| 1128 HY000 Function 'test.substr' is not defined
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT test.substr(b FROM 1) FROM tq
< error
| 1128 HY000 Function 'test.substr' is not defined
> exec
| SELECT test.trim() FROM tq
< error
| 1630 42000 FUNCTION test.trim does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT test.trim() FROM tq
< error
| 1630 42000 FUNCTION test.trim does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT test.trim() FROM tq
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT test.trim(b) FROM tq
< error
| 1630 42000 FUNCTION test.trim does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT test.trim(b) FROM tq
< error
| 1630 42000 FUNCTION test.trim does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT test.trim(b) FROM tq
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT test.trim(b, 'x') FROM tq
< error
| 1630 42000 FUNCTION test.trim does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT test.trim(b, 'x') FROM tq
< error
| 1630 42000 FUNCTION test.trim does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT test.trim(b, 'x') FROM tq
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT test.trim(b, 'x', 1) FROM tq
< error
| 1630 42000 FUNCTION test.trim does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT test.trim(b, 'x', 1) FROM tq
< error
| 1630 42000 FUNCTION test.trim does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT test.trim(b, 'x', 1) FROM tq
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT test.substring() FROM tq
< error
| 1630 42000 FUNCTION test.substring does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT test.substring() FROM tq
< error
| 1630 42000 FUNCTION test.substring does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT test.substring() FROM tq
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT test.substring(b) FROM tq
< error
| 1630 42000 FUNCTION test.substring does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT test.substring(b) FROM tq
< error
| 1630 42000 FUNCTION test.substring does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT test.substring(b) FROM tq
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT test.substring(b, 1, 2) FROM tq
< error
| 1630 42000 FUNCTION test.substring does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT test.substring(b, 1, 2) FROM tq
< error
| 1630 42000 FUNCTION test.substring does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT test.substring(b, 1, 2) FROM tq
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT test.substring(b, 1, 2, 3) FROM tq
< error
| 1630 42000 FUNCTION test.substring does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT test.substring(b, 1, 2, 3) FROM tq
< error
| 1630 42000 FUNCTION test.substring does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT test.substring(b, 1, 2, 3) FROM tq
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT test.substr() FROM tq
< error
| 1630 42000 FUNCTION test.substr does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT test.substr() FROM tq
< error
| 1630 42000 FUNCTION test.substr does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT test.substr() FROM tq
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT test.substr(b, 1, 2, 3) FROM tq
< error
| 1630 42000 FUNCTION test.substr does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| SELECT test.substr(b, 1, 2, 3) FROM tq
< error
| 1630 42000 FUNCTION test.substr does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT test.substr(b, 1, 2, 3) FROM tq
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SELECT test.trim(DISTINCT b) FROM tq
< error
| 1064 42000 …near 'DISTINCT b) FROM tq' at line 1
> exec
| SELECT test.trim(*) FROM tq
< error
| 1064 42000 …near '*) FROM tq' at line 1
> exec
| SELECT test.trim('x' FROM b, 1) FROM tq
< error
| 1064 42000 …near ' 1) FROM tq' at line 1
> exec
| SELECT test.trim(d, INTERVAL 1 DAY) FROM tq
< error
| 1064 42000 …near ') FROM tq' at line 1
> exec
| SELECT test.substring(DISTINCT b, 1) FROM tq
< error
| 1064 42000 …near 'DISTINCT b, 1) FROM tq' at line 1
> exec
| SELECT test.substring(b FROM 1, 2) FROM tq
< error
| 1064 42000 …near ' 2) FROM tq' at line 1
> exec
| SELECT test.substring(b, 1 FOR 2) FROM tq
< error
| 1064 42000 …near 'FOR 2) FROM tq' at line 1
> exec
| SELECT test.substring(d, INTERVAL 1 DAY) FROM tq
< error
| 1064 42000 …near ') FROM tq' at line 1
> exec
| SELECT test.substring(b FROM INTERVAL 1 DAY) FROM tq
< error
| 1064 42000 …near ') FROM tq' at line 1
> exec
| SELECT test.trim(b) OVER () FROM tq
< error
| 1064 42000 …near 'OVER () FROM tq' at line 1
> exec
| SELECT test.substring(b, 1) OVER () FROM tq
< error
| 1064 42000 …near 'OVER () FROM tq' at line 1
> exec
| SELECT test.`trim`(LEADING 'x' FROM b) FROM tq
< error
| 1064 42000 …near 'LEADING 'x' FROM b) FROM tq' at line 1
> exec
| SELECT `test`.`trim`(LEADING 'x' FROM b) FROM tq
< error
| 1064 42000 …near 'LEADING 'x' FROM b) FROM tq' at line 1
> exec
| SELECT test.`substring`(b FROM 1) FROM tq
< error
| 1064 42000 …near 'FROM 1) FROM tq' at line 1
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
| SELECT JSON_ARRAYAGG(DISTINCT a) FROM tq
< ok 0 0
