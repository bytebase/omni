-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestRequiredFormsOracle default#1
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| SELECT _latin1 X'61'
< ok 0 0
> exec ^
: string
| SELECT _latin1 0x61
< ok 0 0
> exec ^
: string
| SELECT _latin1 b'01100001'
< ok 0 0
> exec ^
: string
| SELECT _utf8mb4 x''
< ok 0 0
> exec ^
: string
| SELECT JSON_VALUE('{}', '$.a' DEFAULT _latin1 X'61' ON EMPTY)
< ok 0 0
> exec ^
: string
| SELECT JSON_VALUE('{}', '$.a' DEFAULT b'' ON EMPTY)
< ok 0 0
> exec ^
: string
| SELECT JSON_VALUE('{}', '$.a' DEFAULT +~1 ON EMPTY)
< error
| 1064 42000 …near '~1 ON EMPTY)' at line 1
> exec ^
: string
| SELECT JSON_VALUE('{}', '$.a' DEFAULT +!1 ON EMPTY)
< error
| 1064 42000 …near '!1 ON EMPTY)' at line 1
> exec ^
: string
| SELECT JSON_VALUE('{}', '$.a' DEFAULT -~1 ON EMPTY)
< error
| 1064 42000 …near '~1 ON EMPTY)' at line 1
> exec ^
: string
| SELECT JSON_VALUE('{}', '$.a' DEFAULT ~1 ON EMPTY)
< error
| 1064 42000 …near '~1 ON EMPTY)' at line 1
> exec ^
: string
| SELECT JSON_VALUE('{}', '$.a' DEFAULT - 1 ON EMPTY)
< ok 0 0
> exec ^
: string
| SELECT JSON_VALUE('{}', '$.a' DEFAULT + 1 ON EMPTY)
< ok 0 0
> exec ^
: string
| SELECT JSON_VALUE('{}', '$.a' DEFAULT --1 ON EMPTY)
< error
| 1064 42000 …near '-1 ON EMPTY)' at line 1
> exec ^
: string
| SELECT JSON_VALUE('{}', '$.a' DEFAULT -+1 ON EMPTY)
< error
| 1064 42000 …near '+1 ON EMPTY)' at line 1
> exec ^
: string
| SELECT JSON_VALUE('{}', '$.a' DEFAULT +-1 ON EMPTY)
< error
| 1064 42000 …near '-1 ON EMPTY)' at line 1
> exec ^
: string
| SELECT JSON_VALUE('{}', '$.a' DEFAULT -0x61 ON EMPTY)
< error
| 1064 42000 …near '0x61 ON EMPTY)' at line 1
> exec ^
: string
| SELECT JSON_VALUE('{}', '$.a' DEFAULT -b'1' ON EMPTY)
< error
| 1064 42000 …near 'b'1' ON EMPTY)' at line 1
> exec ^
: string
| SELECT JSON_VALUE('{}', '$.a' DEFAULT -1.5 ON EMPTY)
< ok 0 0
> exec ^
: string
| SELECT JSON_VALUE('{}', '$.a' DEFAULT -1e2 ON EMPTY)
< ok 0 0
> exec ^
: string
| SELECT JSON_VALUE('{}', '$.a' DEFAULT +1.5 ON EMPTY)
< ok 0 0
> exec ^
: string
| SELECT JSON_VALUE('{}', '$.a' DEFAULT -.5 ON EMPTY)
< ok 0 0
> exec ^
: string
| SELECT JSON_VALUE('{}', '$.a' DEFAULT - -1 ON EMPTY)
< error
| 1064 42000 …near '-1 ON EMPTY)' at line 1
> exec ^
: string
| SELECT JSON_VALUE('{}', '$.a' DEFAULT -(1) ON EMPTY)
< error
| 1064 42000 …near '(1) ON EMPTY)' at line 1
> exec ^
: string
| SELECT JSON_VALUE('{}', '$.a' DEFAULT -'1' ON EMPTY)
< error
| 1064 42000 …near ''1' ON EMPTY)' at line 1
> exec ^
: string
| SELECT JSON_VALUE('{}', '$.a' DEFAULT -TRUE ON EMPTY)
< error
| 1064 42000 …near 'TRUE ON EMPTY)' at line 1
> exec ^
: string
| SELECT JSON_VALUE('{}', _latin1 X'24')
< error
| 1064 42000 …near 'X'24')' at line 1
> exec ^
: string
| SELECT JSON_VALUE('{}', _latin1 b'00100100')
< error
| 1064 42000 …near 'b'00100100')' at line 1
> exec ^
: string
| SELECT JSON_VALUE('{}', _latin1 0x24)
< error
| 1064 42000 …near '0x24)' at line 1
> exec ^
: string
| SELECT JSON_VALUE('{}', _latin1'$')
< ok 0 0
> exec ^
: string
| SELECT JSON_VALUE('{}', '$' DEFAULT OFF ON EMPTY)
< error
| 1064 42000 …near 'OFF ON EMPTY)' at line 1
> exec ^
: string
| SELECT JSON_VALUE('{}', '$' DEFAULT FALSE ON EMPTY)
< ok 0 0
> exec ^
: string
| SELECT JSON_VALUE((1 IN (SELECT 1)), '$' NULL ON EMPTY)
< error
| 3146 22032 Invalid data type for JSON data in argument 1 to function json_value; a JSON string or JSON type is required.
> exec ^
: string
| SELECT JSON_VALUE((1 IN (SELECT 1)), '$' NULL ON EMPTY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(1 IN (SELECT 1), '$' NULL ON EMPTY)
< error
| 1064 42000 …near 'IN (SELECT 1), '$' NULL ON EMPTY)' at line 1
> exec ^
: string
| SELECT JSON_VALUE((1 = 1), '$')
< error
| 3146 22032 Invalid data type for JSON data in argument 1 to function json_value; a JSON string or JSON type is required.
> exec ^
: string
| SELECT JSON_VALUE((1 = 1), '$')
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(1 = 1, '$')
< error
| 1064 42000 …near '= 1, '$')' at line 1
> exec ^
: string
| SELECT JSON_VALUE('{}', '$' DEFAULT b'2' ON EMPTY)
< error
| 1064 42000 …near 'b'2' ON EMPTY)' at line 1
> exec ^
: string
| SELECT JSON_VALUE('{}', '$' DEFAULT X'6' ON EMPTY)
< error
| 1064 42000 …near 'X'6' ON EMPTY)' at line 1
> exec ^
: string
| SELECT JSON_VALUE('{}', '$' DEFAULT X'GG' ON EMPTY)
< error
| 1064 42000 …near 'X'GG' ON EMPTY)' at line 1
> exec ^
: string
| SELECT JSON_VALUE('{}', '$' DEFAULT 0x6 ON EMPTY)
< ok 0 0
> exec ^
: string
| SELECT JSON_VALUE('{}', '$' DEFAULT b'1' ON EMPTY)
< ok 0 0
> exec ^
: string
| SELECT `_latin1` X'61'
< error
| 1064 42000 …near 'X'61'' at line 1
> exec ^
: string
| SELECT JSON_VALUE('{}', '$' DEFAULT `_latin1` X'61' ON EMPTY)
< error
| 1064 42000 …near '`_latin1` X'61' ON EMPTY)' at line 1
> exec ^
: string
| SELECT `_latin1`'x'
< error
| 1054 42S22 Unknown column '_latin1' in 'field list'
> exec ^
: string
| SELECT `_latin1`'x'
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT _latin1 1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| SELECT _latin1 X'GG'
< error
| 1064 42000 …near 'X'GG'' at line 1
> exec ^
: string
| SELECT _latin1 b'2'
< error
| 1064 42000 …near 'b'2'' at line 1
> exec ^
: string
| SELECT _latin1 X'6'
< error
| 1064 42000 …near 'X'6'' at line 1
> exec ^
: string
| SELECT _latin1 0x6
< ok 0 0
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT N'x' ON ERROR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT N'x' ON ERROR)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT n'x' ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT n'x' ON EMPTY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT N'x' 'y' ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT N'x' 'y' ON EMPTY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT N 'x' ON ERROR)
< error
| 1064 42000 …near 'N 'x' ON ERROR)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT -N'x' ON ERROR)
< error
| 1064 42000 …near 'N'x' ON ERROR)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT N ON ERROR)
< error
| 1064 42000 …near 'N ON ERROR)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT -'x' ON ERROR)
< error
| 1064 42000 …near ''x' ON ERROR)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT -TRUE ON ERROR)
< error
| 1064 42000 …near 'TRUE ON ERROR)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT +'x' ON ERROR)
< error
| 1064 42000 …near ''x' ON ERROR)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT -X'41' ON ERROR)
< error
| 1064 42000 …near 'X'41' ON ERROR)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT -1.5 ON ERROR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT -1.5 ON ERROR)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT - 1 ON ERROR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT - 1 ON ERROR)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT --1 ON ERROR)
< error
| 1064 42000 …near '-1 ON ERROR)' at line 1
> exec ^
: string
| SELECT JSON_VALUE('{}','$.a')
< ok 0 0
> exec ^
: string
| SELECT json_value(a,'$.a')
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT json_value(a,'$.a')
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING CHAR(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING CHAR(3))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING SIGNED)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING SIGNED)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING DECIMAL(5,2))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING DECIMAL(5,2))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING JSON)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING JSON)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING CHAR(3) ARRAY)
< error
| 1064 42000 …near 'ARRAY)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' NULL ON ERROR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' NULL ON ERROR)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' ERROR ON ERROR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' ERROR ON ERROR)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT 'x' ON ERROR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT 'x' ON ERROR)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT 1 ON ERROR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT 1 ON ERROR)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT b ON ERROR)
< error
| 1064 42000 …near 'b ON ERROR)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' NULL ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' NULL ON EMPTY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' ERROR ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' ERROR ON EMPTY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT 'x' ON EMPTY)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT 'x' ON EMPTY)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' NULL ON EMPTY NULL ON ERROR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' NULL ON EMPTY NULL ON ERROR)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' NULL ON ERROR NULL ON EMPTY)
< error
| 1064 42000 …near 'NULL ON EMPTY)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' NULL ON EMPTY NULL ON EMPTY)
< error
| 1064 42000 …near 'EMPTY)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' NULL ON ERROR NULL ON ERROR)
< error
| 1064 42000 …near 'NULL ON ERROR)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING CHAR(3) NULL ON EMPTY DEFAULT 'x' ON ERROR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING CHAR(3) NULL ON EMPTY DEFAULT 'x' ON ERROR)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' NULL ON EMPTY RETURNING CHAR(3))
< error
| 1064 42000 …near 'RETURNING CHAR(3))' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING CHAR(3) TRUNCATE)
< error
| 1064 42000 …near 'TRUNCATE)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' TRUNCATE)
< error
| 1064 42000 …near 'TRUNCATE)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' ASCII)
< error
| 1064 42000 …near 'ASCII)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a',1)
< error
| 1064 42000 …near ',1)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT JSON_VALUE()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT JSON_VALUE(1,2)
< error
| 1064 42000 …near '2)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,b)
< error
| 1064 42000 …near 'b)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,concat('$','a'))
< error
| 1064 42000 …near 'concat('$','a'))' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,_utf8mb4'$.a')
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,_utf8mb4'$.a')
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,N'$.a')
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,N'$.a')
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,X'24')
< error
| 1064 42000 …near 'X'24')' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,"$.a")
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,"$.a")
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$' '.a')
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$' '.a')
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(1+1,'$.a')
< error
| 1064 42000 …near '+1,'$.a')' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a+1,'$.a')
< error
| 1064 42000 …near '+1,'$.a')' at line 1
> exec ^
: string
| SELECT JSON_VALUE((1+1),'$.a')
< error
| 3146 22032 Invalid data type for JSON data in argument 1 to function json_value; a JSON string or JSON type is required.
> exec ^
: string
| SELECT JSON_VALUE((1+1),'$.a')
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a->'$.b','$.a')
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a->'$.b','$.a')
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(NOT a,'$.a')
< error
| 1064 42000 …near 'NOT a,'$.a')' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a IS NULL,'$.a')
< error
| 1064 42000 …near 'IS NULL,'$.a')' at line 1
> exec ^
: string
| SELECT JSON_VALUE(DISTINCT a,'$.a')
< error
| 1064 42000 …near 'DISTINCT a,'$.a')' at line 1
> exec ^
: string
| SELECT JSON_VALUE(*,'$.a')
< error
| 1064 42000 …near '*,'$.a')' at line 1
> exec ^
: string
| SELECT `json_value`(a,'$.a')
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.json_value does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT `json_value`(a,'$.a')
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT d.json_value(a,'$.a')
< error
| 1370 42000 execute command denied to user 'omni_test'@'%' for routine 'd.json_value'
> exec ^
: string
| SELECT d.json_value(a,'$.a')
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE (a,'$.a')
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE (a,'$.a')
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT json_value FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT json_value FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CREATE TABLE json_value (json_value INT)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a JSON, b INT AS (JSON_VALUE(a,'$.b' RETURNING SIGNED)))
< ok 0 0
> exec ^
: string
| CREATE TABLE t (a JSON, KEY k ((JSON_VALUE(a,'$.b' RETURNING CHAR(3)))))
< ok 0 0
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT -1 ON ERROR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT -1 ON ERROR)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT +1 ON ERROR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT +1 ON ERROR)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT NULL ON ERROR)
< error
| 1064 42000 …near 'NULL ON ERROR)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT 1.5 ON ERROR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT 1.5 ON ERROR)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT TRUE ON ERROR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT TRUE ON ERROR)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT DATE'2000-01-01' ON ERROR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT DATE'2000-01-01' ON ERROR)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT X'41' ON ERROR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT X'41' ON ERROR)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT _utf8mb4'x' ON ERROR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT _utf8mb4'x' ON ERROR)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT 1+1 ON ERROR)
< error
| 1064 42000 …near '+1 ON ERROR)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT (1) ON ERROR)
< error
| 1064 42000 …near '(1) ON ERROR)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT @x ON ERROR)
< error
| 1064 42000 …near '@x ON ERROR)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' DEFAULT ? ON ERROR)
< error
| 1064 42000 …near '? ON ERROR)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING SIGNED INT)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING SIGNED INT)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING UNSIGNED)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING UNSIGNED)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING DATE)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING DATE)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING DATETIME(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING DATETIME(3))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING FLOAT)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING FLOAT)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING DOUBLE)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING DOUBLE)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING BINARY(3))
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING BINARY(3))
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING YEAR)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING YEAR)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING SIGNED RETURNING SIGNED)
< error
| 1064 42000 …near 'RETURNING SIGNED)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING SIGNED TRUNCATE ASCII)
< error
| 1064 42000 …near 'TRUNCATE ASCII)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING CHAR(3) ASCII TRUNCATE)
< error
| 1064 42000 …near 'TRUNCATE)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING CHAR(3) TRUNCATE ASCII)
< error
| 1064 42000 …near 'TRUNCATE ASCII)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' ASCII TRUNCATE)
< error
| 1064 42000 …near 'ASCII TRUNCATE)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' TRUNCATE NULL ON EMPTY)
< error
| 1064 42000 …near 'TRUNCATE NULL ON EMPTY)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' NULL ON EMPTY TRUNCATE)
< error
| 1064 42000 …near 'TRUNCATE)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,'$.a' RETURNING SIGNED ASCII)
< error
| 1064 42000 …near 'ASCII)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(-a,'$.a')
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(-a,'$.a')
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(!a,'$.a')
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(!a,'$.a')
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a COLLATE utf8mb4_bin,'$.a')
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a COLLATE utf8mb4_bin,'$.a')
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(f(a),'$.a')
< error
| 1305 42000 FUNCTION omni_t_mysql_parser.f does not exist
> exec ^
: string
| SELECT JSON_VALUE(f(a),'$.a')
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(@x,'$.a')
< error
| 3146 22032 Invalid data type for JSON data in argument 1 to function json_value; a JSON string or JSON type is required.
> exec ^
: string
| SELECT JSON_VALUE(@x,'$.a')
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(t.a,'$.a')
< error
| 1109 42S02 Unknown table 't' in field list
> exec ^
: string
| SELECT JSON_VALUE(t.a,'$.a')
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(CASE WHEN 1 THEN a END,'$.a')
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(CASE WHEN 1 THEN a END,'$.a')
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a->>'$.b','$.a')
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a->>'$.b','$.a')
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a||b,'$.a')
< error
| 1064 42000 …near '||b,'$.a')' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,'')
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_VALUE(a,'')
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT JSON_VALUE(a,NULL)
< error
| 1064 42000 …near 'NULL)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,@p)
< error
| 1064 42000 …near '@p)' at line 1
> exec ^
: string
| SELECT JSON_VALUE(a,?)
< error
| 1064 42000 …near '?)' at line 1
> exec ^
: string
| CREATE TABLE tl (a key)
< error
| 1064 42000 …near 'key)' at line 1
> exec ^
: string
| CREATE TABLE tl (a primary key)
< error
| 1064 42000 …near 'primary key)' at line 1
> exec ^
: string
| CREATE TABLE tl (a unique)
< error
| 1064 42000 …near 'unique)' at line 1
> exec ^
: string
| CREATE TABLE tl (a unique key)
< error
| 1064 42000 …near 'unique key)' at line 1
> exec ^
: string
| CREATE TABLE tl (a comment 'x')
< error
| 1064 42000 …near 'comment 'x')' at line 1
> exec ^
: string
| CREATE TABLE tl (a not null)
< error
| 1064 42000 …near 'not null)' at line 1
> exec ^
: string
| CREATE TABLE tl (a null)
< error
| 1064 42000 …near 'null)' at line 1
> exec ^
: string
| CREATE TABLE tl (a default 1)
< error
| 1064 42000 …near 'default 1)' at line 1
> exec ^
: string
| CREATE TABLE tl (a auto_increment)
< error
| 1064 42000 …near 'auto_increment)' at line 1
> exec ^
: string
| CREATE TABLE tl (a collate utf8mb4_bin)
< error
| 1064 42000 …near 'collate utf8mb4_bin)' at line 1
> exec ^
: string
| CREATE TABLE tl (a references r (id))
< error
| 1064 42000 …near 'references r (id))' at line 1
> exec ^
: string
| CREATE TABLE tl (a check (a > 0))
< error
| 1064 42000 …near 'check (a > 0))' at line 1
> exec ^
: string
| CREATE TABLE tl (a generated always as (1))
< error
| 1064 42000 …near 'generated always as (1))' at line 1
> exec ^
: string
| CREATE TABLE tl (a generated always as (1) stored)
< error
| 1064 42000 …near 'generated always as (1) stored)' at line 1
> exec ^
: string
| CREATE TABLE tl (a generated always as (1) virtual)
< error
| 1064 42000 …near 'generated always as (1) virtual)' at line 1
> exec ^
: string
| CREATE TABLE tl (a as (1))
< error
| 1064 42000 …near 'as (1))' at line 1
> exec ^
: string
| CREATE TABLE tl (a as (1) stored)
< error
| 1064 42000 …near 'as (1) stored)' at line 1
> exec ^
: string
| CREATE TABLE tl (a on update now())
< error
| 1064 42000 …near 'on update now())' at line 1
> exec ^
: string
| CREATE TABLE tl (a column_format fixed)
< error
| 1064 42000 …near 'column_format fixed)' at line 1
> exec ^
: string
| CREATE TABLE tl (a storage disk)
< error
| 1064 42000 …near 'storage disk)' at line 1
> exec ^
: string
| CREATE TABLE tl (a srid 0)
< error
| 1064 42000 …near 'srid 0)' at line 1
> exec ^
: string
| CREATE TABLE tl (a invisible)
< error
| 1064 42000 …near 'invisible)' at line 1
> exec ^
: string
| CREATE TABLE tl (a visible)
< error
| 1064 42000 …near 'visible)' at line 1
> exec ^
: string
| CREATE TABLE tl (a auto_random)
< error
| 1064 42000 …near 'auto_random)' at line 1
> exec ^
: string
| CREATE TABLE tl (a)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| CREATE TABLE tl (a, b int)
< error
| 1064 42000 …near ', b int)' at line 1
> exec ^
: string
| CREATE TABLE tl (a) SELECT 1 AS a
< error
| 1064 42000 …near ') SELECT 1 AS a' at line 1
> exec ^
: string
| CREATE TABLE tl (a) AS SELECT 1 AS a
< error
| 1064 42000 …near ') AS SELECT 1 AS a' at line 1
> exec ^
: string
| CREATE TABLE tl (a, b) AS SELECT 1 AS a, 2 AS b
< error
| 1064 42000 …near ', b) AS SELECT 1 AS a, 2 AS b' at line 1
> exec ^
: string
| CREATE TABLE tl (a key) SELECT 1 AS a
< error
| 1064 42000 …near 'key) SELECT 1 AS a' at line 1
> exec ^
: string
| CREATE TABLE tl (a default 1) SELECT 1 AS a
< error
| 1064 42000 …near 'default 1) SELECT 1 AS a' at line 1
> exec ^
: string
| CREATE TABLE tl (a not null) AS SELECT 1 AS a
< error
| 1064 42000 …near 'not null) AS SELECT 1 AS a' at line 1
> exec ^
: string
| CREATE TABLE tl (a comment 'x') AS SELECT 1 AS a
< error
| 1064 42000 …near 'comment 'x') AS SELECT 1 AS a' at line 1
> exec ^
: string
| CREATE TEMPORARY TABLE tl (a key)
< error
| 1064 42000 …near 'key)' at line 1
> exec ^
: string
| ALTER TABLE t ADD a key
< error
| 1064 42000 …near 'key' at line 1
> exec ^
: string
| ALTER TABLE t ADD COLUMN a comment 'x'
< error
| 1064 42000 …near 'comment 'x'' at line 1
> exec ^
: string
| ALTER TABLE t ADD COLUMN a generated always as (1)
< error
| 1064 42000 …near 'generated always as (1)' at line 1
> exec ^
: string
| ALTER TABLE t ADD (a key)
< error
| 1064 42000 …near 'key)' at line 1
> exec ^
: string
| ALTER TABLE t ADD COLUMN (a not null, b int)
< error
| 1064 42000 …near 'not null, b int)' at line 1
> exec ^
: string
| ALTER TABLE t MODIFY a comment 'x'
< error
| 1064 42000 …near 'comment 'x'' at line 1
> exec ^
: string
| ALTER TABLE t MODIFY COLUMN a not null
< error
| 1064 42000 …near 'not null' at line 1
> exec ^
: string
| ALTER TABLE t MODIFY a generated always as (1)
< error
| 1064 42000 …near 'generated always as (1)' at line 1
> exec ^
: string
| ALTER TABLE t CHANGE a b key
< error
| 1064 42000 …near 'key' at line 1
> exec ^
: string
| ALTER TABLE t CHANGE COLUMN a b comment 'x'
< error
| 1064 42000 …near 'comment 'x'' at line 1
> exec ^
: string
| ALTER TABLE t CHANGE a b generated always as (1)
< error
| 1064 42000 …near 'generated always as (1)' at line 1
> exec ^
: string
| CREATE TABLE tl (a int key)
< ok 0 0
> exec ^
: string
| CREATE TABLE tl (a int generated always as (1))
< ok 0 0
> exec ^
: string
| CREATE TABLE tl (a int as (1) stored)
< ok 0 0
> exec ^
: string
| CREATE TABLE tl (a serial)
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD COLUMN a int comment 'x'
< ok 0 0
> exec ^
: string
| ALTER TABLE t MODIFY a int not null
< ok 0 0
> exec ^
: string
| ALTER TABLE t CHANGE a b int key
< ok 0 0
> exec ^
: string
| SELECT d.comment(1)
< error
| 1370 42000 execute command denied to user 'omni_test'@'%' for routine 'd.comment'
> exec ^
: string
| SELECT d.comment(1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT `comment`(1)
< error
| 1630 42000 FUNCTION omni_t_mysql_parser.comment does not exist. Check the 'Function Name Parsing and Resolution' section in the Reference Manual
> exec ^
: string
| SELECT `comment`(1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT d.vector(1)
< error
| 1370 42000 execute command denied to user 'omni_test'@'%' for routine 'd.vector'
> exec ^
: string
| SELECT d.vector(1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT comment FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT comment FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 FROM comment
< error
| 1146 42S02 Table 'omni_t_mysql_parser.comment' doesn't exist
> exec ^
: string
| SELECT 1 FROM comment
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| INSERT INTO t VALUES (1) ON DUPLICATE KEY UPDATE a = VALUES(a)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| INSERT INTO t VALUES (1) ON DUPLICATE KEY UPDATE a = VALUES(a)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT 1 = ANY (SELECT 1)
< ok 0 0
> exec ^
: string
| SELECT 1 FROM t WHERE a = SOME (SELECT 1)
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT 1 FROM t WHERE a = SOME (SELECT 1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT rank() OVER ()
< ok 0 0
> exec ^
: string
| SELECT lag(a) OVER () FROM t
< error
| 1146 42S02 Table 'omni_t_mysql_parser.t' doesn't exist
> exec ^
: string
| SELECT lag(a) OVER () FROM t
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| SELECT ntile(1) OVER ()
< ok 0 0
> exec ^
: string
| SET @a = comment(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| DO comment(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT a FROM t ORDER BY comment(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| SELECT a FROM t WHERE comment(1) = 1
< error
| 1064 42000 …near '(1) = 1' at line 1
> exec ^
: string
| CREATE TABLE tf (a int DEFAULT (comment(1)))
< error
| 1064 42000 …near '(1)))' at line 1
> exec ^
: string
| CREATE TABLE tf (a int, b int AS (comment(1)))
< error
| 1064 42000 …near '(1)))' at line 1
> exec ^
: string
| CREATE VIEW vf AS SELECT comment(1)
< error
| 1064 42000 …near '(1)' at line 1
> exec ^
: string
| CREATE FUNCTION comment() RETURNS INT RETURN 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| CREATE FUNCTION comment() RETURNS INT RETURN 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| CALL comment(1)
< error
| 1305 42000 PROCEDURE omni_t_mysql_parser.comment does not exist
> exec ^
: string
| CALL comment(1)
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| DROP FUNCTION comment
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| DROP FUNCTION comment
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
