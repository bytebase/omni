-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestGrantWordsStayNamesOnTiDB default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| SELECT * FROM san
< error
| 1146 42S02 Table 'test.san' doesn't exist
> exec ^
: string
| SELECT * FROM san
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM token_issuer
< error
| 1146 42S02 Table 'test.token_issuer' doesn't exist
> exec ^
: string
| SELECT * FROM token_issuer
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t san WHERE san.id = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t san WHERE san.id = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT * FROM t token_issuer WHERE token_issuer.id = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT * FROM t token_issuer WHERE token_issuer.id = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DELETE FROM t san WHERE san.id = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| DELETE FROM t san WHERE san.id = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| UPDATE t token_issuer SET token_issuer.a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| UPDATE t token_issuer SET token_issuer.a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT san, token_issuer FROM t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| SELECT san, token_issuer FROM t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| SELECT 1 AS san
< ok 0 0
> exec ^
: string
| SELECT 1 AS token_issuer
< ok 0 0
> exec ^
: string
| CREATE TABLE san (token_issuer int)
< ok 0 0
> exec ^
: string
| CREATE TABLE t (san varchar(64), token_issuer varchar(64), INDEX san (san))
< ok 0 0
> exec ^
: string
| ALTER TABLE t ADD COLUMN san varchar(64)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ADD COLUMN san varchar(64)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| DROP TABLE token_issuer
< ok 0 0
> exec ^
: string
| CREATE USER san@'%' REQUIRE SAN 'DNS:example-user'
< ok 0 0
> exec ^
: string
| CREATE USER token_issuer@'%' REQUIRE TOKEN_ISSUER 'issuer-abc'
< ok 0 0
> exec ^
: string
| GRANT SELECT ON san.* TO token_issuer
< ok 0 0
