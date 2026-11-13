-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestJSONQueryFormsOceanBaseOracle default#1
> exec
| PREPARE omni_alter_clause_stmt FROM ?
: string
| SELECT JSON_QUERY(a, '$')
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' RETURNING JSON)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_QUERY(a, '$' RETURNING CHAR(10) CHARSET ascii)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' RETURNING UNSIGNED)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' RETURNING NUMBER(10,2))
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' TRUNCATE)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' ALLOW SCALARS)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' DISALLOW SCALARS)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' PRETTY)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' ASCII)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' WITHOUT WRAPPER)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' WITHOUT ARRAY WRAPPER)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' WITH WRAPPER)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' WITH ARRAY WRAPPER)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' WITH UNCONDITIONAL WRAPPER)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' WITH CONDITIONAL WRAPPER)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' WITH UNCONDITIONAL ARRAY WRAPPER)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' WITH CONDITIONAL ARRAY WRAPPER)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' ASIS)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT JSON_QUERY(a, '$' ERROR ON ERROR)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' NULL ON ERROR)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' EMPTY ON ERROR)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' EMPTY ARRAY ON ERROR)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' EMPTY OBJECT ON ERROR)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' ERROR ON EMPTY)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' NULL ON EMPTY)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' EMPTY ON EMPTY)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' EMPTY ARRAY ON EMPTY)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' EMPTY OBJECT ON EMPTY)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' ERROR ON MISMATCH)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' NULL ON MISMATCH)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' DOT ON MISMATCH)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' MULTIVALUE)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' ERROR ON ERROR NULL ON EMPTY)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' NULL ON EMPTY ERROR ON ERROR)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' ERROR ON ERROR NULL ON MISMATCH)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' NULL ON EMPTY NULL ON MISMATCH)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' ERROR ON ERROR NULL ON EMPTY NULL ON MISMATCH)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' EMPTY ARRAY ON EMPTY EMPTY OBJECT ON ERROR DOT ON MISMATCH)
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' RETURNING CHAR(10) TRUNCATE ALLOW SCALARS PRETTY ASCII WITH CONDITIONAL ARRAY WRAPPER ASIS NULL ON EMPTY NULL ON ERROR NULL ON MISMATCH MULTIVALUE)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT json_query(`a`,'$.ids' RETURNING unsigned WITHOUT ARRAY WRAPPER asis error on error null on empty null on mismatch)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT json_query(`a`,'$.codes' RETURNING char(40) charset utf8mb4 WITHOUT ARRAY WRAPPER asis error on error null on empty null on mismatch)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT json_query(`a`,'$.amts' RETURNING number(10,2) WITHOUT ARRAY WRAPPER asis error on error null on empty null on mismatch)
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| CREATE TABLE t2 (id INT PRIMARY KEY, j JSON, INDEX mv ((json_query(`j`,'$.ids' RETURNING unsigned WITHOUT ARRAY WRAPPER asis error on error null on empty null on mismatch))))
< error
| 1046 3D000 No database selected
> exec ^
: string
| SELECT JSON_QUERY((a), '$')
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(-a, '$')
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a->'$.x', '$')
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a COLLATE utf8mb4_bin, '$')
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, N'$')
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, "$")
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, '$' '.x')
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, _utf8mb4'$')
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT JSON_QUERY(a, _latin1 X'24')
< error
| 5443 42000 invalid data type for return value
> exec ^
: string
| SELECT `json_query`(a, '$')
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec ^
: string
| SELECT `json_query`(1, 2)
< error
| 1582 42000 Incorrect parameter count
> exec ^
: string
| SELECT d.json_query(1, 2)
< error
| 1054 42S22 Unknown column '' in 'field list'
> exec ^
: string
| SELECT json_query FROM (SELECT 1 AS json_query) x
< ok 0 0
> exec ^
: string
| SELECT JSON_QUERY()
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT JSON_QUERY(*)
< error
| 1064 42000 …near '*)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$', '$')
< error
| 1064 42000 …near ', '$')' at line 1
> exec ^
: string
| SELECT JSON_QUERY(1, 2)
< error
| 1064 42000 …near '2)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, a)
< error
| 1064 42000 …near 'a)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, @p)
< error
| 1064 42000 …near '@p)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, CONCAT('$'))
< error
| 1064 42000 …near 'CONCAT('$'))' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a + 1, '$')
< error
| 1064 42000 …near '+ 1, '$')' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a IS NULL, '$')
< error
| 1064 42000 …near 'IS NULL, '$')' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' RETURNING)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' RETURNING VARCHAR(10))
< error
| 1064 42000 …near 'VARCHAR(10))' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' RETURNING TEXT)
< error
| 1064 42000 …near 'TEXT)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' RETURNING JSON RETURNING JSON)
< error
| 1064 42000 …near 'RETURNING JSON)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' TRUNCATE RETURNING JSON)
< error
| 1064 42000 …near 'RETURNING JSON)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' PRETTY ALLOW SCALARS)
< error
| 1064 42000 …near 'ALLOW SCALARS)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' ASCII PRETTY)
< error
| 1064 42000 …near 'PRETTY)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' WITH WRAPPER ASCII)
< error
| 1064 42000 …near 'ASCII)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' ASIS WITH WRAPPER)
< error
| 1064 42000 …near 'WITH WRAPPER)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' MULTIVALUE ASIS)
< error
| 1064 42000 …near 'ASIS)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' MULTIVALUE MULTIVALUE)
< error
| 1064 42000 …near 'MULTIVALUE)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' ALLOW SCALARS DISALLOW SCALARS)
< error
| 1064 42000 …near 'DISALLOW SCALARS)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' NULL ON ERROR NULL ON ERROR)
< error
| 1064 42000 …near 'ERROR)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' NULL ON EMPTY NULL ON EMPTY)
< error
| 1064 42000 …near 'EMPTY)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' NULL ON MISMATCH NULL ON MISMATCH)
< error
| 1064 42000 …near 'NULL ON MISMATCH)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' NULL ON MISMATCH NULL ON EMPTY)
< error
| 1064 42000 …near 'NULL ON EMPTY)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' NULL ON MISMATCH ERROR ON ERROR)
< error
| 1064 42000 …near 'ERROR ON ERROR)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' ALLOW)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' SCALARS)
< error
| 1064 42000 …near 'SCALARS)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' WITHOUT)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' WITH)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' WITHOUT ARRAY)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' WITH CONDITIONAL)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' WITHOUT CONDITIONAL WRAPPER)
< error
| 1064 42000 …near 'CONDITIONAL WRAPPER)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' WITH ARRAY CONDITIONAL WRAPPER)
< error
| 1064 42000 …near 'CONDITIONAL WRAPPER)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' ERROR ON)
< error
| 1064 42000 …near ')' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' NULL ERROR)
< error
| 1064 42000 …near 'ERROR)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' ON ERROR)
< error
| 1064 42000 …near 'ON ERROR)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' EMPTY ON MISMATCH)
< error
| 1064 42000 …near 'MISMATCH)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' EMPTY OBJECT ON MISMATCH)
< error
| 1064 42000 …near 'MISMATCH)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' DOT ON ERROR)
< error
| 1064 42000 …near 'ERROR)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' DOT ON EMPTY)
< error
| 1064 42000 …near 'EMPTY)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' DEFAULT 'x' ON ERROR)
< error
| 1064 42000 …near 'DEFAULT 'x' ON ERROR)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' DEFAULT 'x' ON EMPTY)
< error
| 1064 42000 …near 'DEFAULT 'x' ON EMPTY)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' RETURNING JSON,)
< error
| 1064 42000 …near ',)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' `asis`)
< error
| 1064 42000 …near '`asis`)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' `pretty`)
< error
| 1064 42000 …near '`pretty`)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' `multivalue`)
< error
| 1064 42000 …near '`multivalue`)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' ALLOW `scalars`)
< error
| 1064 42000 …near '`scalars`)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' ERROR ON `error`)
< error
| 1064 42000 …near '`error`)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' NULL ON `mismatch`)
< error
| 1064 42000 …near '`mismatch`)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' `dot` ON MISMATCH)
< error
| 1064 42000 …near '`dot` ON MISMATCH)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' `empty` ON ERROR)
< error
| 1064 42000 …near '`empty` ON ERROR)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' EMPTY `object` ON ERROR)
< error
| 1064 42000 …near '`object` ON ERROR)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' WITH `conditional` WRAPPER)
< error
| 1064 42000 …near '`conditional` WRAPPER)' at line 1
> exec ^
: string
| SELECT `json_query`(a, '$' RETURNING JSON)
< error
| 1064 42000 …near 'JSON)' at line 1
> exec ^
: string
| SELECT d.json_query(a, '$' RETURNING JSON)
< error
| 1064 42000 …near 'JSON)' at line 1
> exec ^
: string
| SELECT JSON_QUERY(a, '$' RETURNING JSON) OVER ()
< error
| 1064 42000 …near '()' at line 1
