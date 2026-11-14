-- pin: default gcr.io/cloud-spanner-emulator/emulator:1.5.54 sha256:caf1bd24c081e005837b5977bae5a250e25cb4da9f25ec1abc91936ad67e4de2 v1.5.54

== TestTypeDifferential default#1
> query
| SELECT CAST(NULL AS INT64)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT CAST(NULL AS STRING)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT CAST(NULL AS BYTES)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT CAST(NULL AS BOOL)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT CAST(NULL AS FLOAT64)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT CAST(NULL AS FLOAT32)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT CAST(NULL AS NUMERIC)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT CAST(NULL AS BIGNUMERIC)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Type not found: BIGNUMERIC [at 1:21] SELECT CAST(NULL AS BIGNUMERIC)                     ^"}
> query
| SELECT CAST(NULL AS JSON)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT CAST(NULL AS DATE)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT CAST(NULL AS TIMESTAMP)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT CAST(NULL AS GEOGRAPHY)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Type not found: GEOGRAPHY [at 1:21] SELECT CAST(NULL AS GEOGRAPHY)                     ^"}
> query
| SELECT CAST(NULL AS TOKENLIST)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"TOKENLIST is an internal-only type and cannot be returned to the user by SQL query.;  when converting column  of type 0x… at position 0 in row cursor"}
> query
| SELECT CAST(NULL AS INT)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Type not found: INT [at 1:21] SELECT CAST(NULL AS INT)                     ^"}
> query
| SELECT CAST(NULL AS INTERVAL)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT CAST(NULL AS foo.bar.Baz)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Type not found: foo.bar.Baz [at 1:21] SELECT CAST(NULL AS foo.bar.Baz)                     ^"}
> query
| SELECT CAST(NULL AS `my.proto.Type`)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Type not found: `my.proto.Type` [at 1:21] SELECT CAST(NULL AS `my.proto.Type`)                     ^"}
> query
| SELECT CAST(NULL AS NUMERIC.foo)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Type not found: NUMERIC.foo [at 1:21] SELECT CAST(NULL AS NUMERIC.foo)                     ^"}
> query
| SELECT CAST(NULL AS STRING(100))
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Parameterized types are not supported [at 1:27] SELECT CAST(NULL AS STRING(100))                           ^"}
> query
| SELECT CAST(NULL AS BYTES(256))
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Parameterized types are not supported [at 1:26] SELECT CAST(NULL AS BYTES(256))                          ^"}
> query
| SELECT CAST(NULL AS NUMERIC(10, 2))
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Parameterized types are not supported [at 1:28] SELECT CAST(NULL AS NUMERIC(10, 2))                            ^"}
> query
| SELECT CAST(NULL AS STRING(MAX))
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Parameterized types are not supported [at 1:27] SELECT CAST(NULL AS STRING(MAX))                           ^"}
> query
| SELECT CAST(NULL AS NUMERIC(10, 2, 3))
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Parameterized types are not supported [at 1:28] SELECT CAST(NULL AS NUMERIC(10, 2, 3))                            ^"}
> query
| SELECT CAST(NULL AS NUMERIC(true))
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Parameterized types are not supported [at 1:28] SELECT CAST(NULL AS NUMERIC(true))                            ^"}
> query
| SELECT CAST(NULL AS NUMERIC('foo'))
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Parameterized types are not supported [at 1:28] SELECT CAST(NULL AS NUMERIC('foo'))                            ^"}
> query
| SELECT CAST(NULL AS INT64 COLLATE 'und:ci')
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Type with collation name is not supported in cast [at 1:27] SELECT CAST(NULL AS INT64 COLLATE 'und:ci')                           ^"}
> query
| SELECT CAST(NULL AS DATE COLLATE 'und:ci')
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Type with collation name is not supported in cast [at 1:26] SELECT CAST(NULL AS DATE COLLATE 'und:ci')                          ^"}
> query
| SELECT CAST(NULL AS ARRAY<INT64>)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT CAST(NULL AS ARRAY<ARRAY<INT64>>)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Arrays of arrays are not supported [at 1:21] SELECT CAST(NULL AS ARRAY\u003cARRAY\u003cINT64\u003e\u003e)                     ^"}
> query
| SELECT 1 FROM UNNEST(CAST(NULL AS ARRAY<INT64 COLLATE 'und:ci'>))
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Type with collation name is not supported in cast [at 1:47] SELECT 1 FROM UNNEST(CAST(NULL AS ARRAY\u003cINT64 COLLATE 'und:ci'\u003e))                                               ^"}
> query
| SELECT CAST(NULL AS RANGE<DATE>)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Type not found: RANGE\u003cDATE\u003e"}
> query
| SELECT CAST(NULL AS RANGE<INT64>)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Unsupported type: RANGE\u003cINT64\u003e is not supported [at 1:27] SELECT CAST(NULL AS RANGE\u003cINT64\u003e)                           ^"}
> query
| SELECT CAST(NULL AS MAP<INT64, STRING>)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"MAP datatype is not supported [at 1:21] SELECT CAST(NULL AS MAP\u003cINT64, STRING\u003e)                     ^"}
> query
| SELECT CAST(NULL AS FUNCTION<INT64 -> STRING>)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"FUNCTION type not supported [at 1:21] SELECT CAST(NULL AS FUNCTION\u003cINT64 -\u003e STRING\u003e)                     ^"}
> query
| SELECT CAST(NULL AS FUNCTION<(INT64, STRING) -> BOOL>)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"FUNCTION type not supported [at 1:21] SELECT CAST(NULL AS FUNCTION\u003c(INT64, STRING) -\u003e BOOL\u003e)                     ^"}
> query
| SELECT CAST(NULL AS FUNCTION<() -> INT64>)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"FUNCTION type not supported [at 1:21] SELECT CAST(NULL AS FUNCTION\u003c() -\u003e INT64\u003e)                     ^"}
> query
| SELECT STRUCT<x INT64, y STRING>(1, 'a').x
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT STRUCT<INT64, STRING>(1, 'a').field1
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Field name field1 does not exist in STRUCT\u003cINT64, STRING\u003e [at 1:38] SELECT STRUCT\u003cINT64, STRING\u003e(1, 'a').field1                                      ^"}
> query
| SELECT STRUCT<a ARRAY<INT64>>([1,2]).a
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT STRUCT<ARRAY<INT64>>([1,2]).field1
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Field name field1 does not exist in STRUCT\u003cARRAY\u003cINT64\u003e\u003e [at 1:36] SELECT STRUCT\u003cARRAY\u003cINT64\u003e\u003e([1,2]).field1                                    ^"}
> query
| SELECT STRUCT<x ARRAY<STRUCT<y INT64>>>([STRUCT<y INT64>(1)]).x
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT STRUCT<x INT64 COLLATE 'und:ci'>('a').x
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Type with collation name is not supported in literal value construction [at 1:23] SELECT STRUCT\u003cx INT64 COLLATE 'und:ci'\u003e('a').x                       ^"}
> query
| SELECT CAST(NULL AS ARRAY<>)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Expected \"\u003c\" but got \"\u003c\u003e\" [at 1:26] SELECT CAST(NULL AS ARRAY\u003c\u003e)                          ^"}
> query
| SELECT CAST(NULL AS ARRAY< >)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Unexpected \"\u003e\" [at 1:28] SELECT CAST(NULL AS ARRAY\u003c \u003e)                            ^"}
> query
| SELECT CAST(NULL AS ARRAY<INT64, STRING>)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Expected \"\u003e\" but got \",\" [at 1:32] SELECT CAST(NULL AS ARRAY\u003cINT64, STRING\u003e)                                ^"}
> query
| SELECT CAST(NULL AS ARRAY<INT64>>)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Expected \")\" or keyword FORMAT but got \"\u003e\" [at 1:33] SELECT CAST(NULL AS ARRAY\u003cINT64\u003e\u003e)                                 ^"}
> query
| SELECT CAST(NULL AS ARRAY)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Expected \"\u003c\" but got \")\" [at 1:26] SELECT CAST(NULL AS ARRAY)                          ^"}
> query
| SELECT CAST(NULL AS ARRAY.foo)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Expected \"\u003c\" but got \".\" [at 1:26] SELECT CAST(NULL AS ARRAY.foo)                          ^"}
> query
| SELECT CAST(NULL AS STRUCT)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Expected \"\u003c\" but got \")\" [at 1:27] SELECT CAST(NULL AS STRUCT)                           ^"}
> query
| SELECT STRUCT<x INT64,>(1).x
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Unexpected \"\u003e\" [at 1:23] SELECT STRUCT\u003cx INT64,\u003e(1).x                       ^"}
> query
| SELECT CAST(NULL AS RANGE)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Expected \"\u003c\" but got \")\" [at 1:26] SELECT CAST(NULL AS RANGE)                          ^"}
> query
| SELECT CAST(NULL AS RANGE<>)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Expected \"\u003c\" but got \"\u003c\u003e\" [at 1:26] SELECT CAST(NULL AS RANGE\u003c\u003e)                          ^"}
> query
| SELECT CAST(NULL AS INTERVAL.foo)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Expected \")\" or keyword FORMAT but got \".\" [at 1:29] SELECT CAST(NULL AS INTERVAL.foo)                             ^"}
> query
| SELECT CAST(NULL AS INTERVAL<INT64>)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Expected \")\" or keyword FORMAT but got \"\u003c\" [at 1:29] SELECT CAST(NULL AS INTERVAL\u003cINT64\u003e)                             ^"}
> query
| SELECT CAST(NULL AS STRING())
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Unexpected \")\" [at 1:28] SELECT CAST(NULL AS STRING())                            ^"}
> query
| SELECT CAST(NULL AS STRING(,))
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Unexpected \",\" [at 1:28] SELECT CAST(NULL AS STRING(,))                            ^"}
> query
| SELECT CAST(NULL AS NUMERIC(-5))
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Unexpected \"-\" [at 1:29] SELECT CAST(NULL AS NUMERIC(-5))                             ^"}
> query
| SELECT CAST(NULL AS NUMERIC(10,))
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Trailing comma in type parameter list is not allowed. [at 1:31] SELECT CAST(NULL AS NUMERIC(10,))                               ^"}
> query
| SELECT CAST(NULL AS ARRAY<INT64>(vector_length=>128))
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Unexpected identifier \"vector_length\" [at 1:34] SELECT CAST(NULL AS ARRAY\u003cINT64\u003e(vector_length=\u003e128))                                  ^"}
> query
| SELECT CAST(NULL AS MAP<>)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Expected \")\" or keyword FORMAT but got \"\u003c\u003e\" [at 1:24] SELECT CAST(NULL AS MAP\u003c\u003e)                        ^"}
> query
| SELECT CAST(NULL AS FUNCTION<>)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Expected \")\" or keyword FORMAT but got \"\u003c\u003e\" [at 1:29] SELECT CAST(NULL AS FUNCTION\u003c\u003e)                             ^"}
> query
| SELECT CAST(NULL AS MAP<INT64>)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Expected \",\" but got \"\u003e\" [at 1:30] SELECT CAST(NULL AS MAP\u003cINT64\u003e)                              ^"}
> query
| SELECT CAST(NULL AS FUNCTION<INT64>)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Expected -\u003e but got \"\u003e\" [at 1:35] SELECT CAST(NULL AS FUNCTION\u003cINT64\u003e)                                   ^"}
> query
| SELECT CAST(NULL AS ``)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"reject","kind":"query","reason":"syntax","code":"InvalidArgument","message":"Syntax error: Invalid empty identifier [at 1:21] SELECT CAST(NULL AS ``)                     ^"}
> query
| SELECT STRUCT<STRUCT<x INT64>>(STRUCT<x INT64>(1)).field1
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Field name field1 does not exist in STRUCT\u003cSTRUCT\u003cx INT64\u003e\u003e [at 1:52] SELECT STRUCT\u003cSTRUCT\u003cx INT64\u003e\u003e(STRUCT\u003cx INT64\u003e(1)).field1                                                    ^"}
> query
| SELECT STRUCT<x ARRAY<ARRAY<INT64>>>([[1]]).x
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Arrays of arrays are not supported [at 1:17] SELECT STRUCT\u003cx ARRAY\u003cARRAY\u003cINT64\u003e\u003e\u003e([[1]]).x                 ^"}
> query
| SELECT CAST(NULL AS ARRAY<STRUCT<a INT64 COLLATE 'x'>>)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Type with collation name is not supported in cast [at 1:42] SELECT CAST(NULL AS ARRAY\u003cSTRUCT\u003ca INT64 COLLATE 'x'\u003e\u003e)                                          ^"}
> query
| SELECT CAST(NULL AS ARRAY<NUMERIC(10, 2)>)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Parameterized types are not supported [at 1:34] SELECT CAST(NULL AS ARRAY\u003cNUMERIC(10, 2)\u003e)                                  ^"}
> query
| SELECT CAST(NULL AS STRING(0))
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Parameterized types are not supported [at 1:27] SELECT CAST(NULL AS STRING(0))                           ^"}
> query
| SELECT CAST(NULL AS STRING(0x10))
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Parameterized types are not supported [at 1:27] SELECT CAST(NULL AS STRING(0x10))                           ^"}
> query
| SELECT STRUCT<key INT64>(1).key
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT STRUCT<row INT64>(1).row
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
