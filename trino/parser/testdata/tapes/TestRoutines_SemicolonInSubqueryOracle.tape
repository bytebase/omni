-- pin: default trinodb/trino:482 sha256:90b35b7c603eaa1f889bf03981a62b75f998ee6c0f851d9f4e341b49a57022b6 482

== TestRoutines_SemicolonInSubqueryOracle default#1
> exec
| CREATE FUNCTION f() RETURNS int BEGIN RETURN (SELECT 1; ); END
< error
| SYNTAX_ERROR 1 1:55
| line 1:55: mismatched input ';'. Expecting: '%', ')', '*', '+', ',', '-', '.', '/', 'AND', 'AS', 'AT', 'EXCEPT', 'FETCH', 'FROM', 'GROUP', 'HAVING', 'INTERSECT', 'LIMIT', 'OFFSET', 'OR', 'ORDER', 'UNION', 'WHERE', 'WINDOW', '[', '||', <identifier>, <predicate>
> exec
| CREATE FUNCTION f() RETURNS int BEGIN RETURN (SELECT 1;); END
< error
| SYNTAX_ERROR 1 1:55
| line 1:55: mismatched input ';'. Expecting: '%', ')', '*', '+', ',', '-', '.', '/', 'AND', 'AS', 'AT', 'EXCEPT', 'FETCH', 'FROM', 'GROUP', 'HAVING', 'INTERSECT', 'LIMIT', 'OFFSET', 'OR', 'ORDER', 'UNION', 'WHERE', 'WINDOW', '[', '||', <identifier>, <predicate>
> exec
| CREATE FUNCTION f() RETURNS int BEGIN RETURN (SELECT 1; SELECT 2); END
< error
| SYNTAX_ERROR 1 1:55
| line 1:55: mismatched input ';'. Expecting: '%', ')', '*', '+', ',', '-', '.', '/', 'AND', 'AS', 'AT', 'EXCEPT', 'FETCH', 'FROM', 'GROUP', 'HAVING', 'INTERSECT', 'LIMIT', 'OFFSET', 'OR', 'ORDER', 'UNION', 'WHERE', 'WINDOW', '[', '||', <identifier>, <predicate>
> exec
| CREATE FUNCTION f() RETURNS int BEGIN RETURN (SELECT 1; garbage); END
< error
| SYNTAX_ERROR 1 1:55
| line 1:55: mismatched input ';'. Expecting: '%', ')', '*', '+', ',', '-', '.', '/', 'AND', 'AS', 'AT', 'EXCEPT', 'FETCH', 'FROM', 'GROUP', 'HAVING', 'INTERSECT', 'LIMIT', 'OFFSET', 'OR', 'ORDER', 'UNION', 'WHERE', 'WINDOW', '[', '||', <identifier>, <predicate>
> exec
| CREATE FUNCTION f() RETURNS int BEGIN RETURN (SELECT 1 /* c */ ; /* d */ ); END
< error
| SYNTAX_ERROR 1 1:64
| line 1:64: mismatched input ';'. Expecting: '%', ')', '*', '+', ',', '-', '.', '/', 'AND', 'AS', 'AT', 'EXCEPT', 'FETCH', 'FROM', 'GROUP', 'HAVING', 'INTERSECT', 'LIMIT', 'OFFSET', 'OR', 'ORDER', 'UNION', 'WHERE', 'WINDOW', '[', '||', <identifier>, <predicate>
> exec
| CREATE FUNCTION f() RETURNS int BEGIN RETURN 1 + (SELECT 1;;); END
< error
| SYNTAX_ERROR 1 1:59
| line 1:59: mismatched input ';'. Expecting: '%', ')', '*', '+', ',', '-', '.', '/', 'AND', 'AS', 'AT', 'EXCEPT', 'FETCH', 'FROM', 'GROUP', 'HAVING', 'INTERSECT', 'LIMIT', 'OFFSET', 'OR', 'ORDER', 'UNION', 'WHERE', 'WINDOW', '[', '||', <identifier>, <predicate>
> exec
| CREATE FUNCTION f() RETURNS int BEGIN RETURN (SELECT (SELECT 1;)); END
< error
| SYNTAX_ERROR 1 1:63
| line 1:63: mismatched input ';'. Expecting: '%', ')', '*', '+', ',', '-', '.', '/', 'AND', 'AS', 'AT', 'EXCEPT', 'FETCH', 'FROM', 'GROUP', 'HAVING', 'INTERSECT', 'LIMIT', 'OFFSET', 'OR', 'ORDER', 'UNION', 'WHERE', 'WINDOW', '[', '||', <identifier>, <predicate>
> exec
| CREATE FUNCTION f() RETURNS int BEGIN RETURN (SELECT (SELECT 1;) + !1); END
< error
| SYNTAX_ERROR 1 1:63
| line 1:63: mismatched input ';'. Expecting: '%', ')', '*', '+', ',', '-', '.', '/', 'AND', 'AS', 'AT', 'EXCEPT', 'FETCH', 'FROM', 'GROUP', 'HAVING', 'INTERSECT', 'LIMIT', 'OFFSET', 'OR', 'ORDER', 'UNION', 'WHERE', 'WINDOW', '[', '||', <identifier>, <predicate>
> exec
| CREATE FUNCTION f() RETURNS int BEGIN RETURN (SELECT (SELECT 1;) ; ); END
< error
| SYNTAX_ERROR 1 1:63
| line 1:63: mismatched input ';'. Expecting: '%', ')', '*', '+', ',', '-', '.', '/', 'AND', 'AS', 'AT', 'EXCEPT', 'FETCH', 'FROM', 'GROUP', 'HAVING', 'INTERSECT', 'LIMIT', 'OFFSET', 'OR', 'ORDER', 'UNION', 'WHERE', 'WINDOW', '[', '||', <identifier>, <predicate>
> exec
| CREATE FUNCTION f() RETURNS int BEGIN RETURN (SELECT (SELECT 1;) FROM t x y); END
< error
| SYNTAX_ERROR 1 1:63
| line 1:63: mismatched input ';'. Expecting: '%', ')', '*', '+', ',', '-', '.', '/', 'AND', 'AS', 'AT', 'EXCEPT', 'FETCH', 'FROM', 'GROUP', 'HAVING', 'INTERSECT', 'LIMIT', 'OFFSET', 'OR', 'ORDER', 'UNION', 'WHERE', 'WINDOW', '[', '||', <identifier>, <predicate>
> exec
| CREATE FUNCTION f() RETURNS int BEGIN RETURN (SELECT * FROM (SELECT 1;) t); END
< error
| SYNTAX_ERROR 1 1:70
| line 1:70: mismatched input ';'. Expecting: '%', ')', '*', '+', ',', '-', '.', '/', 'AND', 'AS', 'AT', 'EXCEPT', 'FETCH', 'FROM', 'GROUP', 'HAVING', 'INTERSECT', 'LIMIT', 'OFFSET', 'OR', 'ORDER', 'UNION', 'WHERE', 'WINDOW', '[', '||', <identifier>, <predicate>
> exec
| CREATE FUNCTION f() RETURNS int BEGIN RETURN (VALUES 1;); END
< error
| SYNTAX_ERROR 1 1:55
| line 1:55: mismatched input ';'. Expecting: '%', ')', '*', '+', ',', '-', '.', '/', 'AND', 'AT', 'EXCEPT', 'FETCH', 'INTERSECT', 'LIMIT', 'OFFSET', 'OR', 'ORDER', 'UNION', '[', '||', <predicate>
> exec
| CREATE FUNCTION f() RETURNS int BEGIN RETURN (TABLE t;); END
< error
| SYNTAX_ERROR 1 1:54
| line 1:54: mismatched input ';'. Expecting: ')', '.', 'EXCEPT', 'FETCH', 'INTERSECT', 'LIMIT', 'OFFSET', 'ORDER', 'UNION'
> exec
| CREATE FUNCTION f() RETURNS int BEGIN RETURN (WITH x AS (SELECT 1) SELECT * FROM x;); END
< error
| SYNTAX_ERROR 1 1:83
| line 1:83: mismatched input ';'. Expecting: ')', ',', '.', 'AS', 'CROSS', 'EXCEPT', 'FETCH', 'FOR', 'FULL', 'GROUP', 'HAVING', 'INNER', 'INTERSECT', 'JOIN', 'LEFT', 'LIMIT', 'MATCH_RECOGNIZE', 'NATURAL', 'OFFSET', 'ORDER', 'RIGHT', 'TABLESAMPLE', 'UNION', 'WHERE', 'WINDOW', <identifier>
> exec
| CREATE FUNCTION f() RETURNS boolean BEGIN RETURN EXISTS (SELECT 1;); END
< error
| SYNTAX_ERROR 1 1:66
| line 1:66: mismatched input ';'. Expecting: '%', ')', '*', '+', ',', '-', '.', '/', 'AND', 'AS', 'AT', 'EXCEPT', 'FETCH', 'FROM', 'GROUP', 'HAVING', 'INTERSECT', 'LIMIT', 'OFFSET', 'OR', 'ORDER', 'UNION', 'WHERE', 'WINDOW', '[', '||', <identifier>, <predicate>
> exec
| CREATE FUNCTION f() RETURNS int RETURN CASE WHEN true THEN (SELECT 1;) END
< error
| SYNTAX_ERROR 1 1:69
| line 1:69: mismatched input ';'. Expecting: '%', ')', '*', '+', ',', '-', '.', '/', 'AND', 'AS', 'AT', 'EXCEPT', 'FETCH', 'FROM', 'GROUP', 'HAVING', 'INTERSECT', 'LIMIT', 'OFFSET', 'OR', 'ORDER', 'UNION', 'WHERE', 'WINDOW', '[', '||', <identifier>, <predicate>
> exec
| CREATE FUNCTION f() RETURNS int RETURN if(true, (SELECT 1;), 2)
< error
| SYNTAX_ERROR 1 1:58
| line 1:58: mismatched input ';'. Expecting: '%', ')', '*', '+', ',', '-', '.', '/', 'AND', 'AS', 'AT', 'EXCEPT', 'FETCH', 'FROM', 'GROUP', 'HAVING', 'INTERSECT', 'LIMIT', 'OFFSET', 'OR', 'ORDER', 'UNION', 'WHERE', 'WINDOW', '[', '||', <identifier>, <predicate>
> exec
| WITH FUNCTION f() RETURNS int BEGIN RETURN (SELECT 1;); END SELECT f()
< error
| SYNTAX_ERROR 1 1:53
| line 1:53: mismatched input ';'. Expecting: '%', ')', '*', '+', ',', '-', '.', '/', 'AND', 'AS', 'AT', 'EXCEPT', 'FETCH', 'FROM', 'GROUP', 'HAVING', 'INTERSECT', 'LIMIT', 'OFFSET', 'OR', 'ORDER', 'UNION', 'WHERE', 'WINDOW', '[', '||', <identifier>, <predicate>
> exec
| WITH FUNCTION f() RETURNS int RETURN 1 SELECT CASE WHEN true THEN (SELECT 1;) END
< error
| SYNTAX_ERROR 1 1:76
| line 1:76: mismatched input ';'. Expecting: '%', ')', '*', '+', ',', '-', '.', '/', 'AND', 'AS', 'AT', 'EXCEPT', 'FETCH', 'FROM', 'GROUP', 'HAVING', 'INTERSECT', 'LIMIT', 'OFFSET', 'OR', 'ORDER', 'UNION', 'WHERE', 'WINDOW', '[', '||', <identifier>, <predicate>
> exec
| CREATE FUNCTION f() RETURNS int BEGIN DECLARE x int DEFAULT (SELECT 1;); RETURN x; END
< error
| SYNTAX_ERROR 1 1:70
| line 1:70: mismatched input ';'. Expecting: '%', ')', '*', '+', ',', '-', '.', '/', 'AND', 'AS', 'AT', 'EXCEPT', 'FETCH', 'FROM', 'GROUP', 'HAVING', 'INTERSECT', 'LIMIT', 'OFFSET', 'OR', 'ORDER', 'UNION', 'WHERE', 'WINDOW', '[', '||', <identifier>, <predicate>
> exec
| CREATE FUNCTION f() RETURNS int BEGIN WHILE (SELECT 1;) > 0 DO RETURN 1; END WHILE; RETURN 2; END
< error
| SYNTAX_ERROR 1 1:54
| line 1:54: mismatched input ';'. Expecting: '%', ')', '*', '+', ',', '-', '.', '/', 'AND', 'AS', 'AT', 'EXCEPT', 'FETCH', 'FROM', 'GROUP', 'HAVING', 'INTERSECT', 'LIMIT', 'OFFSET', 'OR', 'ORDER', 'UNION', 'WHERE', 'WINDOW', '[', '||', <identifier>, <predicate>
> exec
| CREATE FUNCTION f() RETURNS int BEGIN RETURN (SELECT 1); END
< error
| NOT_SUPPORTED 13 1:1
| line 1:1: Catalog and schema must be specified when function schema is not configured
> exec
| CREATE FUNCTION f() RETURNS int BEGIN RETURN (SELECT 1) + (SELECT (SELECT 2)); END
< error
| NOT_SUPPORTED 13 1:1
| line 1:1: Catalog and schema must be specified when function schema is not configured
> exec
| WITH FUNCTION f() RETURNS int BEGIN RETURN (SELECT 1); END SELECT (SELECT f())
< error
| NOT_SUPPORTED 13 1:44
| line 1:44: Queries are not allowed in functions
