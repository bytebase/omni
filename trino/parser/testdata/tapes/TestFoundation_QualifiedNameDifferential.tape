-- pin: default trinodb/trino:482 sha256:90b35b7c603eaa1f889bf03981a62b75f998ee6c0f851d9f4e341b49a57022b6 482

== TestFoundation_QualifiedNameDifferential default#1
> exec
| SELECT * FROM t
< error
| TABLE_NOT_FOUND 46 1:15
| line 1:15: Table 'memory.default.t' does not exist
> exec
| SELECT * FROM orders
< error
| TABLE_NOT_FOUND 46 1:15
| line 1:15: Table 'memory.default.orders' does not exist
> exec
| SELECT * FROM public.orders
< error
| SCHEMA_NOT_FOUND 45 1:15
| line 1:15: Schema 'public' does not exist
> exec
| SELECT * FROM tpch.sf1.nation
< ok - -
> exec
| SELECT * FROM "My Table"
< error
| TABLE_NOT_FOUND 46 1:15
| line 1:15: Table 'memory.default."my table"' does not exist
> exec
| SELECT * FROM "public"."order"
< error
| SCHEMA_NOT_FOUND 45 1:15
| line 1:15: Schema 'public' does not exist
> exec
| SELECT * FROM cat."My Schema".tbl
< error
| CATALOG_NOT_FOUND 44 1:15
| line 1:15: Catalog 'cat' not found
> exec
| SELECT * FROM zone
< error
| TABLE_NOT_FOUND 46 1:15
| line 1:15: Table 'memory.default.zone' does not exist
> exec
| SELECT * FROM t.zone
< error
| SCHEMA_NOT_FOUND 45 1:15
| line 1:15: Schema 't' does not exist
> exec
| SELECT * FROM including
< error
| TABLE_NOT_FOUND 46 1:15
| line 1:15: Table 'memory.default.including' does not exist
> exec
| SELECT * FROM from
< error
| SYNTAX_ERROR 1 1:15
| line 1:15: mismatched input 'from'. Expecting: '(', 'JSON_TABLE', 'LATERAL', 'NEAREST', 'TABLE', 'UNNEST', <identifier>
> exec
| SELECT * FROM select
< error
| SYNTAX_ERROR 1 1:15
| line 1:15: mismatched input 'select'. Expecting: '(', 'JSON_TABLE', 'LATERAL', 'NEAREST', 'TABLE', 'UNNEST', <identifier>
> exec
| SELECT * FROM s.from
< error
| SYNTAX_ERROR 1 1:17
| line 1:17: mismatched input 'from'. Expecting: <identifier>
> exec
| SELECT * FROM public.order
< error
| SYNTAX_ERROR 1 1:22
| line 1:22: mismatched input 'order'. Expecting: <identifier>
> exec
| SELECT * FROM 1abc
< error
| SYNTAX_ERROR 1 1:15
| line 1:15: identifiers must not start with a digit; surround the identifier with double quotes
> exec
| SELECT * FROM t.1abc
< error
| SYNTAX_ERROR 1 1:16
| line 1:16: mismatched input '.1'. Expecting: ',', '.', 'AS', 'CROSS', 'EXCEPT', 'FETCH', 'FOR', 'FULL', 'GROUP', 'HAVING', 'INNER', 'INTERSECT', 'JOIN', 'LEFT', 'LIMIT', 'MATCH_RECOGNIZE', 'NATURAL', 'OFFSET', 'ORDER', 'RIGHT', 'TABLESAMPLE', 'UNION', 'WHERE', 'WINDOW', <EOF>, <identifier>
> exec
| SELECT * FROM `bt table`
< error
| SYNTAX_ERROR 1 1:15
| line 1:15: backquoted identifiers are not supported; use double quotes to quote identifiers
> exec
| SELECT * FROM t.
< error
| SYNTAX_ERROR 1 1:17
| line 1:17: mismatched input '<EOF>'. Expecting: <identifier>
