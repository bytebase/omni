-- pin: default trinodb/trino:482 sha256:90b35b7c603eaa1f889bf03981a62b75f998ee6c0f851d9f4e341b49a57022b6 482

== TestJSON_OracleDifferential default#1
> exec
| SELECT JSON_EXISTS(description, 'lax $.children[*]?(@ > 10)')
< error
| COLUMN_NOT_FOUND 47 1:20
| line 1:20: Column 'description' cannot be resolved
> exec
| SELECT JSON_EXISTS(doc, 'strict $.name' PASSING 'Alice' AS name_param)
< error
| COLUMN_NOT_FOUND 47 1:20
| line 1:20: Column 'doc' cannot be resolved
> exec
| SELECT JSON_EXISTS(col, '$.x' FALSE ON ERROR)
< error
| COLUMN_NOT_FOUND 47 1:20
| line 1:20: Column 'col' cannot be resolved
> exec
| SELECT JSON_EXISTS(col, '$.x' TRUE ON ERROR)
< error
| COLUMN_NOT_FOUND 47 1:20
| line 1:20: Column 'col' cannot be resolved
> exec
| SELECT JSON_EXISTS(col, '$.x' UNKNOWN ON ERROR)
< error
| COLUMN_NOT_FOUND 47 1:20
| line 1:20: Column 'col' cannot be resolved
> exec
| SELECT JSON_EXISTS(col, '$.x' ERROR ON ERROR)
< error
| COLUMN_NOT_FOUND 47 1:20
| line 1:20: Column 'col' cannot be resolved
> exec
| SELECT JSON_EXISTS(col FORMAT JSON, '$.x')
< error
| COLUMN_NOT_FOUND 47 1:20
| line 1:20: Column 'col' cannot be resolved
> exec
| SELECT JSON_EXISTS(col FORMAT JSON ENCODING UTF8, '$.x')
< error
| COLUMN_NOT_FOUND 47 1:20
| line 1:20: Column 'col' cannot be resolved
> exec
| SELECT JSON_EXISTS(col, '$.x' PASSING 1 AS a, 2 AS b)
< error
| COLUMN_NOT_FOUND 47 1:20
| line 1:20: Column 'col' cannot be resolved
> exec
| SELECT JSON_EXISTS(c || d, '$.x')
< error
| COLUMN_NOT_FOUND 47 1:20
| line 1:20: Column 'c' cannot be resolved
> exec
| SELECT JSON_EXISTS(CAST(c AS varchar) FORMAT JSON, '$.x')
< error
| COLUMN_NOT_FOUND 47 1:25
| line 1:25: Column 'c' cannot be resolved
> exec
| SELECT JSON_EXISTS(c, '$.x' PASSING d FORMAT JSON AS doc)
< error
| COLUMN_NOT_FOUND 47 1:20
| line 1:20: Column 'c' cannot be resolved
> exec
| SELECT JSON_EXISTS(c, '$.x' PASSING d FORMAT JSON ENCODING UTF16 AS doc, 5 AS num)
< error
| COLUMN_NOT_FOUND 47 1:20
| line 1:20: Column 'c' cannot be resolved
> exec
| SELECT JSON_VALUE(description, 'lax $.children[0]' RETURNING tinyint)
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'description' cannot be resolved
> exec
| SELECT JSON_VALUE(doc, '$.name' DEFAULT 'unknown' ON EMPTY)
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'doc' cannot be resolved
> exec
| SELECT JSON_VALUE(col, 'strict $.id' ERROR ON ERROR)
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'col' cannot be resolved
> exec
| SELECT JSON_VALUE(col, '$.x' NULL ON EMPTY NULL ON ERROR)
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'col' cannot be resolved
> exec
| SELECT JSON_VALUE(col, '$.x' DEFAULT 1 ON EMPTY DEFAULT 2 ON ERROR)
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'col' cannot be resolved
> exec
| SELECT JSON_VALUE(col, '$.x' RETURNING varchar DEFAULT 'a' ON EMPTY)
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'col' cannot be resolved
> exec
| SELECT JSON_VALUE(c, '$.x' DEFAULT 1 + 2 ON EMPTY)
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'c' cannot be resolved
> exec
| SELECT JSON_VALUE(c, '$.x' DEFAULT 'a' || 'b' ON ERROR)
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'c' cannot be resolved
> exec
| SELECT JSON_VALUE(c, '$.x' RETURNING DECIMAL(10,2))
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'c' cannot be resolved
> exec
| SELECT JSON_VALUE(c, '$.x' RETURNING int ERROR ON EMPTY NULL ON ERROR)
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'c' cannot be resolved
> exec
| SELECT JSON_VALUE(c, '$.x')
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'c' cannot be resolved
> exec
| SELECT JSON_VALUE(c FORMAT JSON, '$.x')
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'c' cannot be resolved
> exec
| SELECT JSON_VALUE(c, '$.x' PASSING a AS x, b FORMAT JSON AS y)
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'c' cannot be resolved
> exec
| SELECT JSON_QUERY(description, 'lax $.children[last]' WITH ARRAY WRAPPER)
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'description' cannot be resolved
> exec
| SELECT JSON_QUERY(doc, 'strict $.tags' RETURNING varchar)
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'doc' cannot be resolved
> exec
| SELECT JSON_QUERY(col, 'lax $.x' OMIT QUOTES ON SCALAR STRING NULL ON EMPTY)
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'col' cannot be resolved
> exec
| SELECT JSON_QUERY(col, '$.x' WITHOUT WRAPPER)
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'col' cannot be resolved
> exec
| SELECT JSON_QUERY(col, '$.x' WITHOUT ARRAY WRAPPER)
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'col' cannot be resolved
> exec
| SELECT JSON_QUERY(col, '$.x' WITH CONDITIONAL ARRAY WRAPPER)
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'col' cannot be resolved
> exec
| SELECT JSON_QUERY(col, '$.x' WITH UNCONDITIONAL WRAPPER)
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'col' cannot be resolved
> exec
| SELECT JSON_QUERY(col, '$.x' KEEP QUOTES)
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'col' cannot be resolved
> exec
| SELECT JSON_QUERY(col, '$.x' EMPTY ARRAY ON EMPTY EMPTY OBJECT ON ERROR)
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'col' cannot be resolved
> exec
| SELECT JSON_QUERY(col, '$.x' RETURNING varchar FORMAT JSON WITH ARRAY WRAPPER)
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'col' cannot be resolved
> exec
| SELECT JSON_QUERY(c, '$.x' KEEP QUOTES ON SCALAR STRING)
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'c' cannot be resolved
> exec
| SELECT JSON_QUERY(c, '$.x' RETURNING ARRAY(int))
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'c' cannot be resolved
> exec
| SELECT JSON_QUERY(c, '$.x' WITH ARRAY WRAPPER KEEP QUOTES ERROR ON EMPTY ERROR ON ERROR)
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'c' cannot be resolved
> exec
| SELECT JSON_QUERY(c, '$.x')
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'c' cannot be resolved
> exec
| SELECT JSON_QUERY(c, '$.x' WITH WRAPPER)
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'c' cannot be resolved
> exec
| SELECT JSON_QUERY(c FORMAT JSON, '$.x' NULL ON EMPTY)
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'c' cannot be resolved
> exec
| SELECT JSON_OBJECT('x' : true, 'y' : 1.2, 'z' : 'text')
< ok - -
> exec
| SELECT JSON_OBJECT(KEY 'name' VALUE 'Alice', KEY 'age' VALUE 30)
< ok - -
> exec
| SELECT JSON_OBJECT('x' : null, 'y' : 1 ABSENT ON NULL)
< ok - -
> exec
| SELECT JSON_OBJECT('a' : 1 WITH UNIQUE KEYS)
< ok - -
> exec
| SELECT JSON_OBJECT('a' VALUE 1)
< ok - -
> exec
| SELECT JSON_OBJECT()
< ok - -
> exec
| SELECT JSON_OBJECT('a' : 1 NULL ON NULL RETURNING varchar)
< ok - -
> exec
| SELECT JSON_OBJECT('a' : 1 WITHOUT UNIQUE)
< ok - -
> exec
| SELECT JSON_OBJECT('a' : 1 WITH UNIQUE)
< ok - -
> exec
| SELECT JSON_OBJECT('k' VALUE 1 FORMAT JSON)
< error
| TYPE_MISMATCH 58 1:30
| line 1:30: Cannot read input of type integer as JSON using formatting JSON
> exec
| SELECT JSON_OBJECT(KEY 'a' VALUE 1, 'b' : 2)
< ok - -
> exec
| SELECT JSON_OBJECT(KEY VALUE 1)
< error
| COLUMN_NOT_FOUND 47 1:20
| line 1:20: Column 'key' cannot be resolved
> exec
| SELECT JSON_OBJECT(k VALUE 1)
< error
| COLUMN_NOT_FOUND 47 1:20
| line 1:20: Column 'k' cannot be resolved
> exec
| SELECT JSON_OBJECT(KEY k VALUE 1)
< error
| COLUMN_NOT_FOUND 47 1:24
| line 1:24: Column 'k' cannot be resolved
> exec
| SELECT JSON_OBJECT(KEY 'x' VALUE 1)
< ok - -
> exec
| SELECT JSON_OBJECT(KEY KEY VALUE 1)
< error
| COLUMN_NOT_FOUND 47 1:24
| line 1:24: Column 'key' cannot be resolved
> exec
| SELECT JSON_OBJECT(value VALUE 1)
< error
| COLUMN_NOT_FOUND 47 1:20
| line 1:20: Column 'value' cannot be resolved
> exec
| SELECT JSON_OBJECT('a' || 'b' VALUE 1)
< ok - -
> exec
| SELECT JSON_OBJECT(concat('a','b') VALUE 1)
< ok - -
> exec
| SELECT JSON_OBJECT('a' : 1 + 2)
< ok - -
> exec
| SELECT JSON_OBJECT('a' VALUE 1 + 2)
< ok - -
> exec
| SELECT JSON_OBJECT('k' VALUE col FORMAT JSON)
< error
| COLUMN_NOT_FOUND 47 1:30
| line 1:30: Column 'col' cannot be resolved
> exec
| SELECT JSON_OBJECT('k' : col FORMAT JSON ENCODING UTF8)
< error
| COLUMN_NOT_FOUND 47 1:26
| line 1:26: Column 'col' cannot be resolved
> exec
| SELECT JSON_OBJECT(RETURNING varchar)
< ok - -
> exec
| SELECT JSON_OBJECT('a' : 1, KEY 'b' VALUE 2, 'c' VALUE 3)
< ok - -
> exec
| SELECT JSON_OBJECT('a':1 WITHOUT UNIQUE KEYS)
< ok - -
> exec
| SELECT JSON_OBJECT(KEY 'a' VALUE 1 NULL ON NULL)
< ok - -
> exec
| SELECT JSON_OBJECT('a':1 RETURNING varchar FORMAT JSON)
< ok - -
> exec
| SELECT JSON_ARRAY(true, 12e-1, 'text')
< ok - -
> exec
| SELECT JSON_ARRAY(true, null, 1)
< ok - -
> exec
| SELECT JSON_ARRAY(1, 2, 3 NULL ON NULL)
< ok - -
> exec
| SELECT JSON_ARRAY(x, y RETURNING varbinary FORMAT JSON ENCODING UTF8)
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'x' cannot be resolved
> exec
| SELECT JSON_ARRAY()
< ok - -
> exec
| SELECT JSON_ARRAY(1 ABSENT ON NULL RETURNING varchar)
< ok - -
> exec
| SELECT JSON_ARRAY(NULL)
< ok - -
> exec
| SELECT JSON_ARRAY(1, NULL)
< ok - -
> exec
| SELECT JSON_ARRAY(1 NULL ON NULL)
< ok - -
> exec
| SELECT JSON_ARRAY(RETURNING varchar)
< ok - -
> exec
| SELECT JSON_ARRAY(absent)
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'absent' cannot be resolved
> exec
| SELECT JSON_ARRAY(x ABSENT ON NULL)
< error
| COLUMN_NOT_FOUND 47 1:19
| line 1:19: Column 'x' cannot be resolved
> exec
| SELECT JSON_ARRAY(1 RETURNING varchar FORMAT JSON)
< ok - -
> exec
| SELECT JSON_VALUE(JSON_QUERY(c, '$.a'), '$.b')
< error
| COLUMN_NOT_FOUND 47 1:30
| line 1:30: Column 'c' cannot be resolved
> exec
| SELECT JSON_OBJECT('k' : JSON_QUERY(c, '$.a'))
< error
| COLUMN_NOT_FOUND 47 1:37
| line 1:37: Column 'c' cannot be resolved
> exec
| SELECT JSON_ARRAY(JSON_VALUE(c, '$.a'), JSON_VALUE(c, '$.b'))
< error
| COLUMN_NOT_FOUND 47 1:30
| line 1:30: Column 'c' cannot be resolved
> exec
| SELECT JSON_EXISTS(col)
< error
| SYNTAX_ERROR 1 1:23
| line 1:23: mismatched input ')'. Expecting: '%', '(', '*', '+', ',', '-', '->', '.', '/', '::', 'AND', 'AT', 'FORMAT', 'OR', 'OVER', '[', '||', <predicate>, <string>
> exec
| SELECT JSON_EXISTS(col, '$.x' ON ERROR)
< error
| SYNTAX_ERROR 1 1:31
| line 1:31: mismatched input 'ON'. Expecting: ')', 'AS', 'ERROR', 'FALSE', 'PASSING', 'TRUE', 'UNKNOWN'
> exec
| SELECT JSON_EXISTS(c, somecol)
< error
| SYNTAX_ERROR 1 1:23
| line 1:23: mismatched input 'somecol'. Expecting: <string>
> exec
| SELECT JSON_EXISTS(c, '$.x' PASSING)
< error
| SYNTAX_ERROR 1 1:36
| line 1:36: mismatched input ')'. Expecting: <expression>
> exec
| SELECT JSON_EXISTS(c FORMAT XML, '$.x')
< error
| SYNTAX_ERROR 1 1:29
| line 1:29: mismatched input 'XML'. Expecting: 'JSON'
> exec
| SELECT JSON_EXISTS(c FORMAT JSON ENCODING UTF99, '$.x')
< error
| SYNTAX_ERROR 1 1:43
| line 1:43: mismatched input 'UTF99'. Expecting: 'UTF16', 'UTF32', 'UTF8'
> exec
| SELECT JSON_VALUE(c, '$.x' NULL ON ERROR ERROR ON EMPTY)
< error
| SYNTAX_ERROR 1 1:42
| line 1:42: mismatched input 'ERROR'. Expecting: ')'
> exec
| SELECT JSON_VALUE(c, '$.x' NULL ON EMPTY NULL ON EMPTY)
< error
| SYNTAX_ERROR 1 1:50
| line 1:50: mismatched input 'EMPTY'. Expecting: 'ERROR'
> exec
| SELECT JSON_VALUE(c, '$.x' NULL ON ERROR NULL ON ERROR)
< error
| SYNTAX_ERROR 1 1:42
| line 1:42: mismatched input 'NULL'. Expecting: ')'
> exec
| SELECT JSON_VALUE(c, '$.x' NULL ON EMPTY NULL ON ERROR NULL ON ERROR)
< error
| SYNTAX_ERROR 1 1:56
| line 1:56: mismatched input 'NULL'. Expecting: ')'
> exec
| SELECT JSON_VALUE(c, '$.x' RETURNING varchar FORMAT JSON)
< error
| SYNTAX_ERROR 1 1:46
| line 1:46: mismatched input 'FORMAT'. Expecting: '(', ')', 'ARRAY', 'DEFAULT', 'ERROR', 'NULL'
> exec
| SELECT JSON_VALUE(c, '$.x' || 'y')
< error
| SYNTAX_ERROR 1 1:28
| line 1:28: mismatched input '||'. Expecting: ')', 'AS', 'DEFAULT', 'ERROR', 'NULL', 'PASSING', 'RETURNING'
> exec
| SELECT JSON_QUERY(c, '$.x' KEEP QUOTES WITH ARRAY WRAPPER)
< error
| SYNTAX_ERROR 1 1:40
| line 1:40: mismatched input 'WITH'. Expecting: ')', 'EMPTY', 'ERROR', 'NULL', 'ON'
> exec
| SELECT JSON_QUERY(col, '$.x' ON SCALAR STRING)
< error
| SYNTAX_ERROR 1 1:30
| line 1:30: mismatched input 'ON'. Expecting: ')', 'AS', 'EMPTY', 'ERROR', 'KEEP', 'NULL', 'OMIT', 'PASSING', 'RETURNING', 'WITH', 'WITHOUT'
> exec
| SELECT JSON_QUERY(c, '$.x' KEEP QUOTES ON SCALAR)
< error
| SYNTAX_ERROR 1 1:49
| line 1:49: mismatched input ')'. Expecting: 'STRING'
> exec
| SELECT JSON_QUERY(c, '$.x' KEEP)
< error
| SYNTAX_ERROR 1 1:32
| line 1:32: mismatched input ')'. Expecting: 'QUOTES'
> exec
| SELECT JSON_QUERY(c, '$.x' NULL ON EMPTY NULL ON EMPTY)
< error
| SYNTAX_ERROR 1 1:50
| line 1:50: mismatched input 'EMPTY'. Expecting: 'ERROR'
> exec
| SELECT JSON_OBJECT('k')
< error
| SYNTAX_ERROR 1 1:23
| line 1:23: mismatched input ')'. Expecting: '%', '*', '+', '-', '.', '/', ':', 'AND', 'AT', 'OR', 'VALUE', '[', '||', <predicate>
> exec
| SELECT JSON_OBJECT('k' :)
< error
| SYNTAX_ERROR 1 1:25
| line 1:25: mismatched input ')'. Expecting: <expression>
> exec
| SELECT JSON_OBJECT('k' VALUE)
< error
| SYNTAX_ERROR 1 1:29
| line 1:29: mismatched input ')'. Expecting: <expression>
> exec
| SELECT JSON_OBJECT(NULL ON NULL)
< error
| SYNTAX_ERROR 1 1:25
| line 1:25: mismatched input 'ON'. Expecting: '%', '*', '+', '-', '.', '/', ':', 'AND', 'AT', 'OR', 'VALUE', '[', '||', <predicate>
> exec
| SELECT JSON_OBJECT(ABSENT ON NULL)
< error
| SYNTAX_ERROR 1 1:27
| line 1:27: mismatched input 'ON'. Expecting: '%', '(', '*', '+', '-', '->', '.', '/', ':', '::', 'AND', 'AT', 'OR', 'OVER', 'VALUE', '[', '||', <predicate>, <string>
> exec
| SELECT JSON_OBJECT(KEY)
< error
| SYNTAX_ERROR 1 1:23
| line 1:23: mismatched input ')'. Expecting: '%', '(', '*', '+', '-', '->', '.', '/', ':', '::', 'AND', 'AT', 'OR', 'OVER', 'VALUE', '[', '||', <expression>, <predicate>, <string>
> exec
| SELECT JSON_OBJECT('a':1, NULL ON NULL ABSENT ON NULL)
< error
| SYNTAX_ERROR 1 1:32
| line 1:32: mismatched input 'ON'. Expecting: '%', '*', '+', '-', '.', '/', ':', 'AND', 'AT', 'OR', 'VALUE', '[', '||', <predicate>
> exec
| SELECT JSON_ARRAY(NULL ON NULL)
< error
| SYNTAX_ERROR 1 1:24
| line 1:24: mismatched input 'ON'. Expecting: '%', ')', '*', '+', ',', '-', '.', '/', 'ABSENT', 'AND', 'AT', 'FORMAT', 'NULL', 'OR', 'RETURNING', '[', '||', <predicate>
> exec
| SELECT JSON_ARRAY(NULL ON NULL ABSENT ON NULL)
< error
| SYNTAX_ERROR 1 1:24
| line 1:24: mismatched input 'ON'. Expecting: '%', ')', '*', '+', ',', '-', '.', '/', 'ABSENT', 'AND', 'AT', 'FORMAT', 'NULL', 'OR', 'RETURNING', '[', '||', <predicate>
> exec
| SELECT JSON_ARRAY(ABSENT ON NULL)
< error
| SYNTAX_ERROR 1 1:26
| line 1:26: mismatched input 'ON'. Expecting: '%', '(', ')', '*', '+', ',', '-', '->', '.', '/', '::', 'ABSENT', 'AND', 'AT', 'FORMAT', 'NULL', 'OR', 'OVER', 'RETURNING', '[', '||', <predicate>, <string>
> exec
| SELECT JSON_ARRAY(1, ABSENT ON NULL)
< error
| SYNTAX_ERROR 1 1:29
| line 1:29: mismatched input 'ON'. Expecting: '%', '(', ')', '*', '+', ',', '-', '->', '.', '/', '::', 'ABSENT', 'AND', 'AT', 'FORMAT', 'NULL', 'OR', 'OVER', 'RETURNING', '[', '||', <predicate>, <string>
> exec
| SELECT JSON_ARRAY(1 NULL ON ERROR)
< error
| SYNTAX_ERROR 1 1:29
| line 1:29: mismatched input 'ERROR'. Expecting: 'NULL'
> exec
| SELECT JSON_ARRAY(1 NULL ON)
< error
| SYNTAX_ERROR 1 1:28
| line 1:28: mismatched input ')'. Expecting: 'NULL'
> exec
| SELECT JSON_EXISTS
< error
| SYNTAX_ERROR 1 1:19
| line 1:19: mismatched input '<EOF>'. Expecting: '('
> exec
| SELECT JSON_VALUE
< error
| SYNTAX_ERROR 1 1:18
| line 1:18: mismatched input '<EOF>'. Expecting: '('
> exec
| SELECT JSON_QUERY
< error
| SYNTAX_ERROR 1 1:18
| line 1:18: mismatched input '<EOF>'. Expecting: '('
> exec
| SELECT JSON_OBJECT
< error
| SYNTAX_ERROR 1 1:19
| line 1:19: mismatched input '<EOF>'. Expecting: '('
> exec
| SELECT JSON_ARRAY
< error
| SYNTAX_ERROR 1 1:18
| line 1:18: mismatched input '<EOF>'. Expecting: '('
