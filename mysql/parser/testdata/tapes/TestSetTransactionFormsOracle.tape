-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestSetTransactionFormsOracle default#1
> exec
| PREPARE omni_db_option_stmt FROM ?
: string
| set transaction
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set session transaction
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set global transaction
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set local transaction
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set persist transaction
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set persist_only transaction
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set transaction read
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set transaction read only,
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set transaction read write,
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set transaction, read only
< error
| 1064 42000 …near ', read only' at line 1
> exec ^
: string
| set transaction, isolation level serializable
< error
| 1064 42000 …near ', isolation level serializable' at line 1
> exec ^
: string
| set transaction only
< error
| 1064 42000 …near 'only' at line 1
> exec ^
: string
| set transaction write
< error
| 1064 42000 …near 'write' at line 1
> exec ^
: string
| set transaction read committed
< error
| 1064 42000 …near 'committed' at line 1
> exec ^
: string
| set transaction read uncommitted
< error
| 1064 42000 …near 'uncommitted' at line 1
> exec ^
: string
| set transaction repeatable read
< error
| 1064 42000 …near 'repeatable read' at line 1
> exec ^
: string
| set transaction serializable
< error
| 1064 42000 …near 'serializable' at line 1
> exec ^
: string
| set transaction isolation
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set transaction level serializable
< error
| 1064 42000 …near 'level serializable' at line 1
> exec ^
: string
| set transaction isolation serializable
< error
| 1064 42000 …near 'serializable' at line 1
> exec ^
: string
| set transaction isolation level
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set transaction isolation level repeatable
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set transaction isolation level repeatable write
< error
| 1064 42000 …near 'write' at line 1
> exec ^
: string
| set transaction isolation level read
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set transaction isolation level read only
< error
| 1064 42000 …near 'only' at line 1
> exec ^
: string
| set transaction isolation level read write
< error
| 1064 42000 …near 'write' at line 1
> exec ^
: string
| set transaction isolation level committed
< error
| 1064 42000 …near 'committed' at line 1
> exec ^
: string
| set transaction isolation level uncommitted
< error
| 1064 42000 …near 'uncommitted' at line 1
> exec ^
: string
| set transaction isolation level snapshot
< error
| 1064 42000 …near 'snapshot' at line 1
> exec ^
: string
| set transaction isolation level default
< error
| 1064 42000 …near 'default' at line 1
> exec ^
: string
| set transaction isolation level `serializable`
< error
| 1064 42000 …near '`serializable`' at line 1
> exec ^
: string
| set transaction isolation level 'serializable'
< error
| 1064 42000 …near ''serializable'' at line 1
> exec ^
: string
| set transaction `read` only
< error
| 1064 42000 …near '`read` only' at line 1
> exec ^
: string
| set transaction read `only`
< error
| 1064 42000 …near '`only`' at line 1
> exec ^
: string
| set transaction `isolation` level serializable
< error
| 1064 42000 …near '`isolation` level serializable' at line 1
> exec ^
: string
| set transaction isolation `level` serializable
< error
| 1064 42000 …near '`level` serializable' at line 1
> exec ^
: string
| set transaction read only
< ok 0 0
> exec ^
: string
| set transaction read write
< ok 0 0
> exec ^
: string
| set transaction read only;
< ok 0 0
> exec ^
: string
| SET TRANSACTION READ ONLY
< ok 0 0
> exec ^
: string
| SET TRANSACTION ISOLATION LEVEL READ COMMITTED
< ok 0 0
> exec ^
: string
| set transaction /* c */ read /* c */ only
< ok 0 0
> exec ^
: string
| set transaction isolation -- c
| level serializable
< ok 0 0
> exec ^
: string
| set transaction isolation level serializable
< ok 0 0
> exec ^
: string
| set transaction isolation level repeatable read
< ok 0 0
> exec ^
: string
| set transaction isolation level read committed
< ok 0 0
> exec ^
: string
| set transaction isolation level read uncommitted
< ok 0 0
> exec ^
: string
| set transaction read only, isolation level serializable
< ok 0 0
> exec ^
: string
| set transaction read write, isolation level read committed
< ok 0 0
> exec ^
: string
| set transaction isolation level serializable, read only
< ok 0 0
> exec ^
: string
| set transaction isolation level repeatable read, read write
< ok 0 0
> exec ^
: string
| set transaction isolation level read uncommitted, read write
< ok 0 0
> exec ^
: string
| set transaction isolation level read committed, read only
< ok 0 0
> exec ^
: string
| set transaction read only, read write
< error
| 1064 42000 …near 'read write' at line 1
> exec ^
: string
| set transaction read only, read only
< error
| 1064 42000 …near 'read only' at line 1
> exec ^
: string
| set transaction read write, read only
< error
| 1064 42000 …near 'read only' at line 1
> exec ^
: string
| set transaction read write, read write
< error
| 1064 42000 …near 'read write' at line 1
> exec ^
: string
| set transaction isolation level serializable, isolation level read committed
< error
| 1064 42000 …near 'isolation level read committed' at line 1
> exec ^
: string
| set transaction isolation level serializable, isolation level serializable
< error
| 1064 42000 …near 'isolation level serializable' at line 1
> exec ^
: string
| set transaction read only, isolation level serializable, read write
< error
| 1064 42000 …near ', read write' at line 1
> exec ^
: string
| set transaction read only, isolation level serializable, isolation level read committed
< error
| 1064 42000 …near ', isolation level read committed' at line 1
> exec ^
: string
| set transaction isolation level serializable, read only, isolation level read committed
< error
| 1064 42000 …near ', isolation level read committed' at line 1
> exec ^
: string
| set transaction isolation level serializable, read only, read write
< error
| 1064 42000 …near ', read write' at line 1
> exec ^
: string
| set transaction read only, isolation level serializable,
< error
| 1064 42000 …near ',' at line 1
> exec ^
: string
| set transaction isolation level serializable, read only,
< error
| 1064 42000 …near ',' at line 1
> exec ^
: string
| set transaction isolation level serializable,
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set transaction read only, isolation
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set transaction read only, isolation level
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set transaction read only, isolation level read
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set transaction read only, isolation level repeatable
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set transaction read only, read
< error
| 1064 42000 …near 'read' at line 1
> exec ^
: string
| set transaction read only, serializable
< error
| 1064 42000 …near 'serializable' at line 1
> exec ^
: string
| set transaction isolation level serializable, read
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set transaction isolation level serializable, only
< error
| 1064 42000 …near 'only' at line 1
> exec ^
: string
| set transaction isolation level serializable, isolation
< error
| 1064 42000 …near 'isolation' at line 1
> exec ^
: string
| set transaction read only read write
< error
| 1064 42000 …near 'read write' at line 1
> exec ^
: string
| set transaction read only isolation level serializable
< error
| 1064 42000 …near 'isolation level serializable' at line 1
> exec ^
: string
| set transaction isolation level serializable read only
< error
| 1064 42000 …near 'read only' at line 1
> exec ^
: string
| set transaction isolation level serializable isolation level read committed
< error
| 1064 42000 …near 'isolation level read committed' at line 1
> exec ^
: string
| set transaction read only foo
< error
| 1064 42000 …near 'foo' at line 1
> exec ^
: string
| set transaction read write foo
< error
| 1064 42000 …near 'foo' at line 1
> exec ^
: string
| set transaction isolation level serializable foo
< error
| 1064 42000 …near 'foo' at line 1
> exec ^
: string
| set transaction read only, foo
< error
| 1064 42000 …near 'foo' at line 1
> exec ^
: string
| set transaction read only with consistent snapshot
< error
| 1064 42000 …near 'with consistent snapshot' at line 1
> exec ^
: string
| set transaction with consistent snapshot
< error
| 1064 42000 …near 'with consistent snapshot' at line 1
> exec ^
: string
| set transaction read only for update
< error
| 1064 42000 …near 'for update' at line 1
> exec ^
: string
| set session transaction read only
< ok 0 0
> exec ^
: string
| set session transaction read write
< ok 0 0
> exec ^
: string
| set session transaction isolation level serializable
< ok 0 0
> exec ^
: string
| set session transaction isolation level read committed, read only
< ok 0 0
> exec ^
: string
| set session transaction read only, read write
< error
| 1064 42000 …near 'read write' at line 1
> exec ^
: string
| set session transaction read
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set global transaction read only
< ok 0 0
> exec ^
: string
| set global transaction read write
< ok 0 0
> exec ^
: string
| set global transaction isolation level serializable
< ok 0 0
> exec ^
: string
| set global transaction read only, isolation level serializable
< ok 0 0
> exec ^
: string
| set global transaction read only, read write
< error
| 1064 42000 …near 'read write' at line 1
> exec ^
: string
| set global transaction isolation
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set local transaction read only
< ok 0 0
> exec ^
: string
| set local transaction read write
< ok 0 0
> exec ^
: string
| set local transaction isolation level serializable
< ok 0 0
> exec ^
: string
| set local transaction isolation level read committed, read only
< ok 0 0
> exec ^
: string
| set local transaction read only, read write
< error
| 1064 42000 …near 'read write' at line 1
> exec ^
: string
| set local transaction read
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set persist transaction read only
< ok 0 0
> exec ^
: string
| set persist transaction read write
< ok 0 0
> exec ^
: string
| set persist transaction isolation level serializable
< ok 0 0
> exec ^
: string
| set persist transaction read only, isolation level serializable
< ok 0 0
> exec ^
: string
| set persist transaction read only, read write
< error
| 1064 42000 …near 'read write' at line 1
> exec ^
: string
| set persist transaction read
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set persist_only transaction read only
< ok 0 0
> exec ^
: string
| set persist_only transaction read write
< ok 0 0
> exec ^
: string
| set persist_only transaction isolation level serializable
< ok 0 0
> exec ^
: string
| set persist_only transaction read only, isolation level serializable
< ok 0 0
> exec ^
: string
| set persist_only transaction read only, read write
< error
| 1064 42000 …near 'read write' at line 1
> exec ^
: string
| set persist_only transaction read
< error
| 1064 42000 …near '' at line 1
> exec ^
: string
| set @@session transaction read only
< error
| 1064 42000 …near 'transaction read only' at line 1
> exec ^
: string
| set @@global transaction read only
< error
| 1064 42000 …near 'transaction read only' at line 1
> exec ^
: string
| set @@local transaction read only
< error
| 1064 42000 …near 'transaction read only' at line 1
> exec ^
: string
| set @@persist transaction read only
< error
| 1064 42000 …near 'transaction read only' at line 1
> exec ^
: string
| set @@persist_only transaction read only
< error
| 1064 42000 …near 'transaction read only' at line 1
> exec ^
: string
| set session session transaction read only
< error
| 1064 42000 …near 'session transaction read only' at line 1
> exec ^
: string
| set global session transaction read only
< error
| 1064 42000 …near 'session transaction read only' at line 1
> exec ^
: string
| set session global transaction read only
< error
| 1064 42000 …near 'global transaction read only' at line 1
> exec ^
: string
| set local local transaction read only
< error
| 1064 42000 …near 'local transaction read only' at line 1
> exec ^
: string
| set session.transaction read only
< error
| 1064 42000 …near 'session.transaction read only' at line 1
> exec ^
: string
| set `session` transaction read only
< error
| 1064 42000 …near 'transaction read only' at line 1
> exec ^
: string
| set `global` transaction read only
< error
| 1064 42000 …near 'transaction read only' at line 1
> exec ^
: string
| set `transaction` read only
< error
| 1064 42000 …near 'read only' at line 1
> exec ^
: string
| set @a = 1, transaction read only
< error
| 1064 42000 …near 'read only' at line 1
> exec ^
: string
| set sql_mode = '', transaction read only
< error
| 1064 42000 …near 'read only' at line 1
> exec ^
: string
| set session sql_mode = '', transaction read only
< error
| 1064 42000 …near 'read only' at line 1
> exec ^
: string
| set sql_mode = '', session transaction read only
< error
| 1064 42000 …near 'read only' at line 1
> exec ^
: string
| set transaction read only, @a = 1
< error
| 1064 42000 …near '@a = 1' at line 1
> exec ^
: string
| set transaction read only, sql_mode = ''
< error
| 1064 42000 …near 'sql_mode = ''' at line 1
> exec ^
: string
| set transaction read only, session sql_mode = ''
< error
| 1064 42000 …near 'session sql_mode = ''' at line 1
> exec ^
: string
| set transaction isolation level serializable, @a = 1
< error
| 1064 42000 …near '@a = 1' at line 1
> exec ^
: string
| set transaction read only, transaction read write
< error
| 1064 42000 …near 'transaction read write' at line 1
> exec ^
: string
| set transaction read only, session transaction read write
< error
| 1064 42000 …near 'session transaction read write' at line 1
> exec ^
: string
| set transaction isolation level serializable = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set transaction read only = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set transaction read = 1
< error
| 1064 42000 …near '= 1' at line 1
> exec ^
: string
| set @transaction = 1
< ok 0 0
> exec ^
: string
| set transaction_isolation = 'READ-COMMITTED'
< ok 0 0
> exec ^
: string
| set session transaction_isolation = 'READ-COMMITTED'
< ok 0 0
> exec ^
: string
| set global transaction_isolation = 'READ-COMMITTED'
< ok 0 0
> exec ^
: string
| set transaction_read_only = 1
< ok 0 0
> exec ^
: string
| set session transaction_read_only = 1
< ok 0 0
> exec ^
: string
| create procedure p() set transaction read only
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() set transaction read only
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin set transaction read only; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin set transaction read only; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin set transaction isolation level serializable, read only; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin set transaction isolation level serializable, read only; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin set transaction; end
< error
| 1064 42000 …near '; end' at line 1
> exec ^
: string
| create procedure p() begin set transaction read; end
< error
| 1064 42000 …near '; end' at line 1
> exec ^
: string
| create procedure p() begin set transaction read only, read write; end
< error
| 1064 42000 …near 'read write; end' at line 1
> exec ^
: string
| create procedure p() begin set local transaction read only; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin set local transaction read only; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec ^
: string
| create procedure p() begin set persist transaction read only; end
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| create procedure p() begin set persist transaction read only; end
| ; SELECT FROM
< error
| 1064 42000 …near 'SELECT FROM' at line 2
