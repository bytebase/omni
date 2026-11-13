-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestSetNamesScopesFormsOracle default#1
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| set global sql_mode = ''
< ok 0 0
> exec ^
: string
| set global a = 1
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| set global a.b = 1
< error
| 1193 HY000 Unknown system variable 'a.b'
> exec ^
: string
| set global = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set global := 1
< error
| 1064 42000 …near ':= 1' at line 1
> exec ^
: string
| set global to 1
< error
| 1064 42000 …near 'to 1' at line 1
> exec ^
: string
| set global.a = 1
< error
| 1064 42000 …near 'global.a = 1' at line 1
> exec ^
: string
| set global
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set global sql_mode
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set global sql_mode =
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set global global sql_mode = ''
< error
| 1064 42000 …near 'global sql_mode = ''' at line 1
> exec ^
: string
| set global sql_mode := ''
< ok 0 0
> exec ^
: string
| set global sql_mode to ''
< error
| 1064 42000 …near 'to ''' at line 1
> exec ^
: string
| set sql_mode = '', global sql_mode = ''
< ok 0 0
> exec ^
: string
| set @a = 1, global sql_mode = ''
< ok 0 0
> exec ^
: string
| set @a = 1, global a = 1
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| set @a = 1, global = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set global sql_mode = '', @a = 1
< ok 0 0
> exec ^
: string
| set global sql_mode = '', sql_mode = ''
< ok 0 0
> exec ^
: string
| set global transaction = 1
< error
| 1193 HY000 Unknown system variable 'transaction'
> exec ^
: string
| set global `transaction` = 1
< error
| 1193 HY000 Unknown system variable 'transaction'
> exec ^
: string
| set global read = 1
< error
| 1064 42000 …near 'read = 1' at line 1
> exec ^
: string
| set `global` sql_mode = ''
< error
| 1064 42000 …near 'sql_mode = ''' at line 1
> exec ^
: string
| set `global` = 1
< error
| 1193 HY000 Unknown system variable 'global'
> exec ^
: string
| set global sql_mode = '' 2
< error
| 1064 42000 …near '2' at line 1
> exec ^
: string
| set session sql_mode = ''
< ok 0 0
> exec ^
: string
| set session a = 1
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| set session a.b = 1
< error
| 1193 HY000 Unknown system variable 'a.b'
> exec ^
: string
| set session = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set session := 1
< error
| 1064 42000 …near ':= 1' at line 1
> exec ^
: string
| set session to 1
< error
| 1064 42000 …near 'to 1' at line 1
> exec ^
: string
| set session.a = 1
< error
| 1064 42000 …near 'session.a = 1' at line 1
> exec ^
: string
| set session
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set session sql_mode
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set session sql_mode =
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set session session sql_mode = ''
< error
| 1064 42000 …near 'session sql_mode = ''' at line 1
> exec ^
: string
| set session sql_mode := ''
< ok 0 0
> exec ^
: string
| set session sql_mode to ''
< error
| 1064 42000 …near 'to ''' at line 1
> exec ^
: string
| set sql_mode = '', session sql_mode = ''
< ok 0 0
> exec ^
: string
| set @a = 1, session sql_mode = ''
< ok 0 0
> exec ^
: string
| set @a = 1, session a = 1
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| set @a = 1, session = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set session sql_mode = '', @a = 1
< ok 0 0
> exec ^
: string
| set session sql_mode = '', sql_mode = ''
< ok 0 0
> exec ^
: string
| set session transaction = 1
< error
| 1193 HY000 Unknown system variable 'transaction'
> exec ^
: string
| set session `transaction` = 1
< error
| 1193 HY000 Unknown system variable 'transaction'
> exec ^
: string
| set session read = 1
< error
| 1064 42000 …near 'read = 1' at line 1
> exec ^
: string
| set `session` sql_mode = ''
< error
| 1064 42000 …near 'sql_mode = ''' at line 1
> exec ^
: string
| set `session` = 1
< error
| 1193 HY000 Unknown system variable 'session'
> exec ^
: string
| set session sql_mode = '' 2
< error
| 1064 42000 …near '2' at line 1
> exec ^
: string
| set local sql_mode = ''
< ok 0 0
> exec ^
: string
| set local a = 1
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| set local a.b = 1
< error
| 1193 HY000 Unknown system variable 'a.b'
> exec ^
: string
| set local = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set local := 1
< error
| 1064 42000 …near ':= 1' at line 1
> exec ^
: string
| set local to 1
< error
| 1064 42000 …near 'to 1' at line 1
> exec ^
: string
| set local.a = 1
< error
| 1064 42000 …near 'local.a = 1' at line 1
> exec ^
: string
| set local
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set local sql_mode
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set local sql_mode =
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set local local sql_mode = ''
< error
| 1064 42000 …near 'local sql_mode = ''' at line 1
> exec ^
: string
| set local sql_mode := ''
< ok 0 0
> exec ^
: string
| set local sql_mode to ''
< error
| 1064 42000 …near 'to ''' at line 1
> exec ^
: string
| set sql_mode = '', local sql_mode = ''
< ok 0 0
> exec ^
: string
| set @a = 1, local sql_mode = ''
< ok 0 0
> exec ^
: string
| set @a = 1, local a = 1
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| set @a = 1, local = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set local sql_mode = '', @a = 1
< ok 0 0
> exec ^
: string
| set local sql_mode = '', sql_mode = ''
< ok 0 0
> exec ^
: string
| set local transaction = 1
< error
| 1193 HY000 Unknown system variable 'transaction'
> exec ^
: string
| set local `transaction` = 1
< error
| 1193 HY000 Unknown system variable 'transaction'
> exec ^
: string
| set local read = 1
< error
| 1064 42000 …near 'read = 1' at line 1
> exec ^
: string
| set `local` sql_mode = ''
< error
| 1064 42000 …near 'sql_mode = ''' at line 1
> exec ^
: string
| set `local` = 1
< error
| 1193 HY000 Unknown system variable 'local'
> exec ^
: string
| set local sql_mode = '' 2
< error
| 1064 42000 …near '2' at line 1
> exec ^
: string
| set persist sql_mode = ''
< ok 0 0
> exec ^
: string
| set persist a = 1
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| set persist a.b = 1
< error
| 1193 HY000 Unknown system variable 'a.b'
> exec ^
: string
| set persist = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set persist := 1
< error
| 1064 42000 …near ':= 1' at line 1
> exec ^
: string
| set persist to 1
< error
| 1064 42000 …near 'to 1' at line 1
> exec ^
: string
| set persist.a = 1
< error
| 1193 HY000 Unknown system variable 'persist.a'
> exec ^
: string
| set persist
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set persist sql_mode
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set persist sql_mode =
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set persist persist sql_mode = ''
< error
| 1064 42000 …near 'persist sql_mode = ''' at line 1
> exec ^
: string
| set persist sql_mode := ''
< ok 0 0
> exec ^
: string
| set persist sql_mode to ''
< error
| 1064 42000 …near 'to ''' at line 1
> exec ^
: string
| set sql_mode = '', persist sql_mode = ''
< ok 0 0
> exec ^
: string
| set @a = 1, persist sql_mode = ''
< ok 0 0
> exec ^
: string
| set @a = 1, persist a = 1
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| set @a = 1, persist = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set persist sql_mode = '', @a = 1
< ok 0 0
> exec ^
: string
| set persist sql_mode = '', sql_mode = ''
< ok 0 0
> exec ^
: string
| set persist transaction = 1
< error
| 1193 HY000 Unknown system variable 'transaction'
> exec ^
: string
| set persist `transaction` = 1
< error
| 1193 HY000 Unknown system variable 'transaction'
> exec ^
: string
| set persist read = 1
< error
| 1064 42000 …near 'read = 1' at line 1
> exec ^
: string
| set `persist` sql_mode = ''
< error
| 1064 42000 …near 'sql_mode = ''' at line 1
> exec ^
: string
| set `persist` = 1
< error
| 1193 HY000 Unknown system variable 'persist'
> exec ^
: string
| set persist sql_mode = '' 2
< error
| 1064 42000 …near '2' at line 1
> exec ^
: string
| set persist_only sql_mode = ''
< ok 0 0
> exec ^
: string
| set persist_only a = 1
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| set persist_only a.b = 1
< error
| 1193 HY000 Unknown system variable 'a.b'
> exec ^
: string
| set persist_only = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set persist_only := 1
< error
| 1064 42000 …near ':= 1' at line 1
> exec ^
: string
| set persist_only to 1
< error
| 1064 42000 …near 'to 1' at line 1
> exec ^
: string
| set persist_only.a = 1
< error
| 1193 HY000 Unknown system variable 'persist_only.a'
> exec ^
: string
| set persist_only
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set persist_only sql_mode
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set persist_only sql_mode =
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set persist_only persist_only sql_mode = ''
< error
| 1064 42000 …near 'persist_only sql_mode = ''' at line 1
> exec ^
: string
| set persist_only sql_mode := ''
< ok 0 0
> exec ^
: string
| set persist_only sql_mode to ''
< error
| 1064 42000 …near 'to ''' at line 1
> exec ^
: string
| set sql_mode = '', persist_only sql_mode = ''
< ok 0 0
> exec ^
: string
| set @a = 1, persist_only sql_mode = ''
< ok 0 0
> exec ^
: string
| set @a = 1, persist_only a = 1
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| set @a = 1, persist_only = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set persist_only sql_mode = '', @a = 1
< ok 0 0
> exec ^
: string
| set persist_only sql_mode = '', sql_mode = ''
< ok 0 0
> exec ^
: string
| set persist_only transaction = 1
< error
| 1193 HY000 Unknown system variable 'transaction'
> exec ^
: string
| set persist_only `transaction` = 1
< error
| 1193 HY000 Unknown system variable 'transaction'
> exec ^
: string
| set persist_only read = 1
< error
| 1064 42000 …near 'read = 1' at line 1
> exec ^
: string
| set `persist_only` sql_mode = ''
< error
| 1064 42000 …near 'sql_mode = ''' at line 1
> exec ^
: string
| set `persist_only` = 1
< error
| 1193 HY000 Unknown system variable 'persist_only'
> exec ^
: string
| set persist_only sql_mode = '' 2
< error
| 1064 42000 …near '2' at line 1
> exec ^
: string
| set @a = 1, session transaction = 2
< error
| 1193 HY000 Unknown system variable 'transaction'
> exec ^
: string
| set global . a = 1
< error
| 1064 42000 …near '. a = 1' at line 1
> exec ^
: string
| set persist . a = 1
< error
| 1064 42000 …near '. a = 1' at line 1
> exec ^
: string
| set global global = 1
< error
| 1064 42000 …near 'global = 1' at line 1
> exec ^
: string
| set global session = 1
< error
| 1064 42000 …near 'session = 1' at line 1
> exec ^
: string
| set session transaction.a = 1
< error
| 1193 HY000 Unknown system variable 'transaction.a'
> exec ^
: string
| set session add.a = 1
< error
| 1193 HY000 Unknown system variable 'add.a'
> exec ^
: string
| set global persist.a = 1
< error
| 1193 HY000 Unknown system variable 'persist.a'
> exec ^
: string
| set global default.a = 1
< error
| 1193 HY000 Unknown system variable 'default.a'
> exec ^
: string
| set @a = 1, global transaction = 2
< error
| 1193 HY000 Unknown system variable 'transaction'
> exec ^
: string
| set @a = 1, local transaction = 2
< error
| 1193 HY000 Unknown system variable 'transaction'
> exec ^
: string
| set @a = 1, local a = 2
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| set @a = 1, global.a = 2
< error
| 1064 42000 …near 'global.a = 2' at line 1
> exec ^
: string
| set @a = 1, global
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set @a = 1, global sql_mode
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set @a = 1, global sql_mode =
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set @a = 1, global global sql_mode = ''
< error
| 1064 42000 …near 'global sql_mode = ''' at line 1
> exec ^
: string
| set session sql_mode = '', global sql_mode = '', local sql_mode = ''
< ok 0 0
> exec ^
: string
| set global sql_mode = '', a = 1
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| set global names = 1
< error
| 1193 HY000 Unknown system variable 'names'
> exec ^
: string
| set global password = 1
< error
| 1193 HY000 Unknown system variable 'password'
> exec ^
: string
| set session names := 1
< error
| 1193 HY000 Unknown system variable 'names'
> exec ^
: string
| set global role = 1
< error
| 1193 HY000 Unknown system variable 'role'
> exec ^
: string
| set global charset = 1
< error
| 1193 HY000 Unknown system variable 'charset'
> exec ^
: string
| set @a = 1, global names = 2
< error
| 1193 HY000 Unknown system variable 'names'
> exec ^
: string
| set @a = 1, global role = 2
< error
| 1193 HY000 Unknown system variable 'role'
> exec ^
: string
| set @a = 1, global to 2
< error
| 1064 42000 …near 'to 2' at line 1
> exec ^
: string
| set global a to 1
< error
| 1064 42000 …near 'to 1' at line 1
> exec ^
: string
| set session a to 1
< error
| 1064 42000 …near 'to 1' at line 1
> exec ^
: string
| set `persist` a = 1
< error
| 1064 42000 …near 'a = 1' at line 1
> exec ^
: string
| set global @a = 1
< error
| 1064 42000 …near '@a = 1' at line 1
> exec ^
: string
| set global @@sql_mode = ''
< error
| 1064 42000 …near '@@sql_mode = ''' at line 1
> exec ^
: string
| set global @@global.sql_mode = ''
< error
| 1064 42000 …near '@@global.sql_mode = ''' at line 1
> exec ^
: string
| set @@global.sql_mode = ''
< ok 0 0
> exec ^
: string
| set @@global = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set @@global.a = 1
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| set @@global . sql_mode = ''
< ok 0 0
> exec ^
: string
| set @a = 1, @@global.sql_mode = ''
< ok 0 0
> exec ^
: string
| set session @a = 1
< error
| 1064 42000 …near '@a = 1' at line 1
> exec ^
: string
| set session @@sql_mode = ''
< error
| 1064 42000 …near '@@sql_mode = ''' at line 1
> exec ^
: string
| set session @@global.sql_mode = ''
< error
| 1064 42000 …near '@@global.sql_mode = ''' at line 1
> exec ^
: string
| set @@session.sql_mode = ''
< ok 0 0
> exec ^
: string
| set @@session = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set @@session.a = 1
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| set @@session . sql_mode = ''
< ok 0 0
> exec ^
: string
| set @a = 1, @@session.sql_mode = ''
< ok 0 0
> exec ^
: string
| set local @a = 1
< error
| 1064 42000 …near '@a = 1' at line 1
> exec ^
: string
| set local @@sql_mode = ''
< error
| 1064 42000 …near '@@sql_mode = ''' at line 1
> exec ^
: string
| set local @@global.sql_mode = ''
< error
| 1064 42000 …near '@@global.sql_mode = ''' at line 1
> exec ^
: string
| set @@local.sql_mode = ''
< ok 0 0
> exec ^
: string
| set @@local = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set @@local.a = 1
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| set @@local . sql_mode = ''
< ok 0 0
> exec ^
: string
| set @a = 1, @@local.sql_mode = ''
< ok 0 0
> exec ^
: string
| set persist @a = 1
< error
| 1064 42000 …near '@a = 1' at line 1
> exec ^
: string
| set persist @@sql_mode = ''
< error
| 1064 42000 …near '@@sql_mode = ''' at line 1
> exec ^
: string
| set persist @@global.sql_mode = ''
< error
| 1064 42000 …near '@@global.sql_mode = ''' at line 1
> exec ^
: string
| set @@persist.sql_mode = ''
< ok 0 0
> exec ^
: string
| set @@persist = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set @@persist.a = 1
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| set @@persist . sql_mode = ''
< ok 0 0
> exec ^
: string
| set @a = 1, @@persist.sql_mode = ''
< ok 0 0
> exec ^
: string
| set persist_only @a = 1
< error
| 1064 42000 …near '@a = 1' at line 1
> exec ^
: string
| set persist_only @@sql_mode = ''
< error
| 1064 42000 …near '@@sql_mode = ''' at line 1
> exec ^
: string
| set persist_only @@global.sql_mode = ''
< error
| 1064 42000 …near '@@global.sql_mode = ''' at line 1
> exec ^
: string
| set @@persist_only.sql_mode = ''
< ok 0 0
> exec ^
: string
| set @@persist_only = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set @@persist_only.a = 1
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| set @@persist_only . sql_mode = ''
< ok 0 0
> exec ^
: string
| set @a = 1, @@persist_only.sql_mode = ''
< ok 0 0
> exec ^
: string
| set @@sql_mode = ''
< ok 0 0
> exec ^
: string
| set @@sql_mode := ''
< ok 0 0
> exec ^
: string
| set @@sql_mode to ''
< error
| 1064 42000 …near 'to ''' at line 1
> exec ^
: string
| set @@a.b = 1
< error
| 1193 HY000 Unknown system variable 'a.b'
> exec ^
: string
| set @@global.a.b = 1
< error
| 1193 HY000 Unknown system variable 'a.b'
> exec ^
: string
| set @@ = 1
< error
| 1064 42000 …near ' = 1' at line 1
> exec ^
: string
| set @@. = 1
< error
| 1064 42000 …near '. = 1' at line 1
> exec ^
: string
| set @@a = 1
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| set @a = 1, global @b = 2
< error
| 1064 42000 …near '@b = 2' at line 1
> exec ^
: string
| set @a = 1, global @@sql_mode = ''
< error
| 1064 42000 …near '@@sql_mode = ''' at line 1
> exec ^
: string
| set @a = 1, @@sql_mode to ''
< error
| 1064 42000 …near 'to ''' at line 1
> exec ^
: string
| set global @@global.sql_mode to ''
< error
| 1064 42000 …near '@@global.sql_mode to ''' at line 1
> exec ^
: string
| set sql_mode = '', @@session.sql_mode = '', global sql_mode = ''
< ok 0 0
> exec ^
: string
| set transaction = 1
< error
| 1193 HY000 Unknown system variable 'transaction'
> exec ^
: string
| set transaction := 1
< error
| 1193 HY000 Unknown system variable 'transaction'
> exec ^
: string
| set transaction to 1
< error
| 1064 42000 …near 'to 1' at line 1
> exec ^
: string
| set transaction.a = 1
< error
| 1193 HY000 Unknown system variable 'transaction.a'
> exec ^
: string
| set `transaction` = 1
< error
| 1193 HY000 Unknown system variable 'transaction'
> exec ^
: string
| set `transaction`.a = 1
< error
| 1193 HY000 Unknown system variable 'transaction.a'
> exec ^
: string
| set a.transaction = 1
< error
| 1193 HY000 Unknown system variable 'a.transaction'
> exec ^
: string
| set @a = 1, transaction = 2
< error
| 1193 HY000 Unknown system variable 'transaction'
> exec ^
: string
| set @a = 1, `transaction` = 2
< error
| 1193 HY000 Unknown system variable 'transaction'
> exec ^
: string
| set transaction = 1, @a = 2
< error
| 1193 HY000 Unknown system variable 'transaction'
> exec ^
: string
| set isolation = 1
< error
| 1193 HY000 Unknown system variable 'isolation'
> exec ^
: string
| set level = 1
< error
| 1193 HY000 Unknown system variable 'level'
> exec ^
: string
| set read = 1
< error
| 1064 42000 …near 'read = 1' at line 1
> exec ^
: string
| set only = 1
< error
| 1193 HY000 Unknown system variable 'only'
> exec ^
: string
| set write = 1
< error
| 1064 42000 …near 'write = 1' at line 1
> exec ^
: string
| set names = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set names := 1
< error
| 1064 42000 …near ':= 1' at line 1
> exec ^
: string
| set names to 1
< error
| 1064 42000 …near 'to 1' at line 1
> exec ^
: string
| set names.a = 1
< error
| 1193 HY000 Unknown system variable 'names.a'
> exec ^
: string
| set charset = 1
< error
| 1193 HY000 Unknown system variable 'charset'
> exec ^
: string
| set charset := 1
< error
| 1193 HY000 Unknown system variable 'charset'
> exec ^
: string
| set charset.a = 1
< error
| 1193 HY000 Unknown system variable 'charset.a'
> exec ^
: string
| set character = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set character.a = 1
< error
| 1193 HY000 Unknown system variable 'character.a'
> exec ^
: string
| set char = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set role = 1
< error
| 1193 HY000 Unknown system variable 'role'
> exec ^
: string
| set resource = 1
< error
| 1193 HY000 Unknown system variable 'resource'
> exec ^
: string
| set resource := 1
< error
| 1193 HY000 Unknown system variable 'resource'
> exec ^
: string
| set resource.a = 1
< error
| 1193 HY000 Unknown system variable 'resource.a'
> exec ^
: string
| set old = 1
< error
| 1229 HY000 Variable 'old' is a GLOBAL variable and should be set with SET GLOBAL
> exec ^
: string
| set old = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| set new = 1
< ok 0 0
> exec ^
: string
| set role := 1
< error
| 1193 HY000 Unknown system variable 'role'
> exec ^
: string
| set role.a = 1
< error
| 1193 HY000 Unknown system variable 'role.a'
> exec ^
: string
| set password = 1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| set password := 1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| set password.a = 1
< error
| 1193 HY000 Unknown system variable 'password.a'
> exec ^
: string
| set role to 1
< error
| 1064 42000 …near 'to 1' at line 1
> exec ^
: string
| set charset to 1
< error
| 1064 42000 …near 'to 1' at line 1
> exec ^
: string
| set password to 1
< error
| 1064 42000 …near '1' at line 1
> exec ^
: string
| set transaction . a = 1
< error
| 1193 HY000 Unknown system variable 'transaction.a'
> exec ^
: string
| set role . a = 1
< error
| 1193 HY000 Unknown system variable 'role.a'
> exec ^
: string
| set @a = 1, role = 2
< error
| 1193 HY000 Unknown system variable 'role'
> exec ^
: string
| set @a = 1, charset = 2
< error
| 1193 HY000 Unknown system variable 'charset'
> exec ^
: string
| set @a = 1, names = 2
< error
| 1064 42000 …near '= 2' at line 1
> exec ^
: string
| set @a = 1, password = 2
< error
| 1193 HY000 Unknown system variable 'password'
> exec ^
: string
| set @a = 1, password := 2
< error
| 1193 HY000 Unknown system variable 'password'
> exec ^
: string
| set @a = 1, names := 2
< error
| 1064 42000 …near ':= 2' at line 1
> exec ^
: string
| set @a = 1, names.a = 2
< error
| 1193 HY000 Unknown system variable 'names.a'
> exec ^
: string
| set @a = 1, password.a = 2
< error
| 1193 HY000 Unknown system variable 'password.a'
> exec ^
: string
| set @a = 1, transaction.a = 2
< error
| 1193 HY000 Unknown system variable 'transaction.a'
> exec ^
: string
| set @a = 1, transaction := 2
< error
| 1193 HY000 Unknown system variable 'transaction'
> exec ^
: string
| set @a = 1, transaction to 2
< error
| 1064 42000 …near 'to 2' at line 1
> exec ^
: string
| set a.b = 1
< error
| 1193 HY000 Unknown system variable 'a.b'
> exec ^
: string
| set a . b = 1
< error
| 1193 HY000 Unknown system variable 'a.b'
> exec ^
: string
| set a.b.c = 1
< error
| 1064 42000 …near '.c = 1' at line 1
> exec ^
: string
| set a.b := 1
< error
| 1193 HY000 Unknown system variable 'a.b'
> exec ^
: string
| set `a`.b = 1
< error
| 1193 HY000 Unknown system variable 'a.b'
> exec ^
: string
| set a.`b` = 1
< error
| 1193 HY000 Unknown system variable 'a.b'
> exec ^
: string
| set `a`.`b` = 1
< error
| 1193 HY000 Unknown system variable 'a.b'
> exec ^
: string
| set default.a = 1
< error
| 1193 HY000 Unknown system variable 'default.a'
> exec ^
: string
| set new.a = 1
< error
| 1193 HY000 Unknown system variable 'new.a'
> exec ^
: string
| set old.a = 1
< error
| 1193 HY000 Unknown system variable 'old.a'
> exec ^
: string
| set a.new = 1
< error
| 1193 HY000 Unknown system variable 'a.new'
> exec ^
: string
| set a.old = 1
< error
| 1193 HY000 Unknown system variable 'a.old'
> exec ^
: string
| set a.default = 1
< error
| 1193 HY000 Unknown system variable 'a.default'
> exec ^
: string
| set a.global = 1
< error
| 1193 HY000 Unknown system variable 'a.global'
> exec ^
: string
| set a.select = 1
< error
| 1193 HY000 Unknown system variable 'a.select'
> exec ^
: string
| set .a = 1
< error
| 1064 42000 …near '.a = 1' at line 1
> exec ^
: string
| set a. = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set @a.b = 1
< ok 0 0
> exec ^
: string
| set a.@b = 1
< error
| 1064 42000 …near '@b = 1' at line 1
> exec ^
: string
| set a . select = 1
< error
| 1064 42000 …near 'select = 1' at line 1
> exec ^
: string
| set a . transaction = 1
< error
| 1193 HY000 Unknown system variable 'a.transaction'
> exec ^
: string
| set select .a = 1
< error
| 1064 42000 …near 'select .a = 1' at line 1
> exec ^
: string
| set add .a = 1
< error
| 1064 42000 …near 'add .a = 1' at line 1
> exec ^
: string
| set default . a = 1
< error
| 1193 HY000 Unknown system variable 'default.a'
> exec ^
: string
| set default .a = 1
< error
| 1193 HY000 Unknown system variable 'default.a'
> exec ^
: string
| set a to 1
< error
| 1064 42000 …near 'to 1' at line 1
> exec ^
: string
| set a.b to 1
< error
| 1064 42000 …near 'to 1' at line 1
> exec ^
: string
| set sql_mode to ''
< error
| 1064 42000 …near 'to ''' at line 1
> exec ^
: string
| set @a to 1
< error
| 1064 42000 …near 'to 1' at line 1
> exec ^
: string
| set @a = 1, a to 2
< error
| 1064 42000 …near 'to 2' at line 1
> exec ^
: string
| set @a = 1, @b to 2
< error
| 1064 42000 …near 'to 2' at line 1
> exec ^
: string
| set @`global` = 1
< ok 0 0
> exec ^
: string
| set @`session` = 1
< ok 0 0
> exec ^
: string
| set @`local` = 1
< ok 0 0
> exec ^
: string
| set @`persist` = 1
< ok 0 0
> exec ^
: string
| set @`persist_only` = 1
< ok 0 0
> exec ^
: string
| set a = 1
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| set a := 1
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| set `a` = 1
< error
| 1193 HY000 Unknown system variable 'a'
> exec ^
: string
| set sql_mode = ''
< ok 0 0
> exec ^
: string
| set `sql_mode` = ''
< ok 0 0
> exec ^
: string
| set sql_mode := ''
< ok 0 0
> exec ^
: string
| set 1 = 1
< error
| 1064 42000 …near '1 = 1' at line 1
> exec ^
: string
| set 'a' = 1
< error
| 1064 42000 …near ''a' = 1' at line 1
> exec ^
: string
| set a
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set a =
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set @a
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set @a =
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set @a = 1
< ok 0 0
> exec ^
: string
| set @a := 1
< ok 0 0
> exec ^
: string
| set @a = 1, @b = 2
< ok 0 0
> exec ^
: string
| set @`a` = 1
< ok 0 0
> exec ^
: string
| set a = 1,
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set , a = 1
< error
| 1064 42000 …near ', a = 1' at line 1
> exec ^
: string
| set @a = 1,
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set @a = 1 b = 2
< error
| 1064 42000 …near 'b = 2' at line 1
> exec ^
: string
| set @a = 1 2
< error
| 1064 42000 …near '2' at line 1
> exec ^
: string
| set a, b = 1
< error
| 1064 42000 …near ', b = 1' at line 1
> exec ^
: string
| set sql_mode = binary + 1
< ok 0 0
> exec ^
: string
| set @a = 1, add.a = 2
< error
| 1193 HY000 Unknown system variable 'add.a'
> exec ^
: string
| set GLOBAL SQL_MODE = ''
< ok 0 0
> exec ^
: string
| set Global sql_mode = '', Session sql_mode = ''
< ok 0 0
> exec ^
: string
| create procedure p() begin declare global int; set global = 1; end
< error
| 1064 42000 …near '= 1; end' at line 1
> exec ^
: string
| create procedure p() begin set global sql_mode = ''; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin set global sql_mode = ''; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin set global a; end
< error
| 1064 42000 …near '; end' at line 1
> exec ^
: string
| create procedure p() begin declare session int; set session = 1; end
< error
| 1064 42000 …near '= 1; end' at line 1
> exec ^
: string
| create procedure p() begin set session sql_mode = ''; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin set session sql_mode = ''; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin set session a; end
< error
| 1064 42000 …near '; end' at line 1
> exec ^
: string
| create procedure p() begin declare local int; set local = 1; end
< error
| 1064 42000 …near '= 1; end' at line 1
> exec ^
: string
| create procedure p() begin set local sql_mode = ''; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin set local sql_mode = ''; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin set local a; end
< error
| 1064 42000 …near '; end' at line 1
> exec ^
: string
| create procedure p() begin declare persist int; set persist = 1; end
< error
| 1064 42000 …near '= 1; end' at line 1
> exec ^
: string
| create procedure p() begin set persist sql_mode = ''; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin set persist sql_mode = ''; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin set persist a; end
< error
| 1064 42000 …near '; end' at line 1
> exec ^
: string
| create procedure p() begin declare persist_only int; set persist_only = 1; end
< error
| 1064 42000 …near '= 1; end' at line 1
> exec ^
: string
| create procedure p() begin set persist_only sql_mode = ''; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin set persist_only sql_mode = ''; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin set persist_only a; end
< error
| 1064 42000 …near '; end' at line 1
> exec ^
: string
| create procedure p() begin declare transaction int; set transaction = 1; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin declare transaction int; set transaction = 1; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin declare transaction int; set transaction = 1 2; end
< error
| 1064 42000 …near '2; end' at line 1
> exec ^
: string
| create procedure p() begin declare transaction int; set transaction := 1; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin declare transaction int; set transaction := 1; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin declare transaction int; set transaction to 1; end
< error
| 1064 42000 …near 'to 1; end' at line 1
> exec ^
: string
| create procedure p() begin declare transaction int; set @a = 1, transaction = 2; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin declare transaction int; set @a = 1, transaction = 2; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin declare a int; set a = 1; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin declare a int; set a = 1; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin declare a int; set a = 1 2; end
< error
| 1064 42000 …near '2; end' at line 1
> exec ^
: string
| create procedure p() begin declare a int; set a := 1; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin declare a int; set a := 1; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin declare a int; set a to 1; end
< error
| 1064 42000 …near 'to 1; end' at line 1
> exec ^
: string
| create procedure p() begin declare isolation int; set isolation = 1; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin declare isolation int; set isolation = 1; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin declare charset int; set charset = 1; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin declare charset int; set charset = 1; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin declare role int; set role = 1; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin declare role int; set role = 1; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin declare password int; set password = 1; end
< error
| 1064 42000 …near '1; end' at line 1
> exec ^
: string
| create procedure p() begin declare statement int; set statement = 1; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin declare statement int; set statement = 1; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin declare a int; set @b = 1, a = 2; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin declare a int; set @b = 1, a = 2; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin declare a int; set a = 1, global sql_mode = ''; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin declare a int; set a = 1, global sql_mode = ''; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin declare a int; set a = 1 2, global sql_mode = ''; end
< error
| 1064 42000 …near '2, global sql_mode = ''; end' at line 1
> exec ^
: string
| create procedure p() begin declare a int; set a = 1, @b = 2 3; end
< error
| 1064 42000 …near '3; end' at line 1
> exec ^
: string
| create procedure p() begin declare new int; set new = 1; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin declare new int; set new = 1; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create trigger g before insert on t for each row set new.a = 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create trigger g before insert on t for each row set new.a = 1
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create trigger g before insert on t for each row set new.a = 1, @b = 2
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create trigger g before insert on t for each row set new.a = 1, @b = 2
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create trigger g before insert on t for each row set @b = 1, new.a = 2
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create trigger g before insert on t for each row set @b = 1, new.a = 2
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
