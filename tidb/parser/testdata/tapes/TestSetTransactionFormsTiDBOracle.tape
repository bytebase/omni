-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestSetTransactionFormsTiDBOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| set transaction read only as of timestamp now()
< ok 0 0
> exec ^
: string
| set transaction read only as of timestamp '2021-04-21 00:42:12'
< ok 0 0
> exec ^
: string
| SET TRANSACTION READ ONLY AS OF TIMESTAMP NOW()
< ok 0 0
> exec ^
: string
| set session transaction read only as of timestamp now()
< ok 0 0
> exec ^
: string
| set global transaction read only as of timestamp now()
< ok 0 0
> exec ^
: string
| set local transaction read only as of timestamp now()
< error
| 1064 42000 …line 1 column 26 near "read only as of timestamp now()" 
> exec ^
: string
| set persist transaction read only as of timestamp now()
< error
| 1064 42000 …line 1 column 23 near "transaction read only as of timestamp now()" 
> exec ^
: string
| set persist_only transaction read only as of timestamp now()
< error
| 1064 42000 …line 1 column 28 near "transaction read only as of timestamp now()" 
> exec ^
: string
| set @@session transaction read only as of timestamp now()
< error
| 1064 42000 …line 1 column 25 near "transaction read only as of timestamp now()" 
> exec ^
: string
| set transaction read only as of timestamp now() - interval 3 second
< ok 0 0
> exec ^
: string
| set transaction read only as of timestamp tidb_bounded_staleness(now(), now())
< ok 0 0
> exec ^
: string
| set transaction read only as of timestamp (select now())
< ok 0 0
> exec ^
: string
| set transaction read only as of timestamp a
< ok 0 0
> exec ^
: string
| set transaction read only as of timestamp @a
< ok 0 0
> exec ^
: string
| set transaction read only as of timestamp 1 or 2
< ok 0 0
> exec ^
: string
| set transaction read only as /* c */ of timestamp now()
< ok 0 0
> exec ^
: string
| set transaction read only as
| of
| timestamp now()
< ok 0 0
> exec ^
: string
| set transaction read only as of
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| set transaction read only as of timestamp
< error
| 1064 42000 …line 1 column 41 near "" 
> exec ^
: string
| set transaction read only as of now()
< error
| 1064 42000 …line 1 column 35 near "now()" 
> exec ^
: string
| set transaction read only as of time now()
< error
| 1064 42000 …line 1 column 36 near "time now()" 
> exec ^
: string
| set transaction read only as now()
< error
| 1064 42000 …line 1 column 28 near "as now()" 
> exec ^
: string
| set transaction read only as timestamp now()
< error
| 1064 42000 …line 1 column 28 near "as timestamp now()" 
> exec ^
: string
| set transaction read only of timestamp now()
< error
| 1064 42000 …line 1 column 28 near "of timestamp now()" 
> exec ^
: string
| set transaction read only timestamp now()
< error
| 1064 42000 …line 1 column 35 near "timestamp now()" 
> exec ^
: string
| set transaction read only as `of` timestamp now()
< error
| 1064 42000 …line 1 column 28 near "as `of` timestamp now()" 
> exec ^
: string
| set transaction read only as of `timestamp` now()
< error
| 1064 42000 …line 1 column 43 near "`timestamp` now()" 
> exec ^
: string
| set transaction read only as of timestamp now() a
< error
| 1064 42000 …line 1 column 49 near "a" 
> exec ^
: string
| set transaction read only as of timestamp now() as of timestamp now()
< error
| 1064 42000 …line 1 column 53 near "of timestamp now()" 
> exec ^
: string
| set transaction read only as of timestamp now() read write
< error
| 1064 42000 …line 1 column 52 near "read write" 
> exec ^
: string
| set transaction read only as of timestamp now() for update
< error
| 1064 42000 …line 1 column 51 near "for update" 
> exec ^
: string
| set transaction read write as of timestamp now()
< error
| 1064 42000 …line 1 column 32 near "of timestamp now()" 
> exec ^
: string
| set transaction as of timestamp now()
< error
| 1064 42000 …line 1 column 21 near "of timestamp now()" 
> exec ^
: string
| set transaction isolation level serializable as of timestamp now()
< error
| 1064 42000 …line 1 column 50 near "of timestamp now()" 
> exec ^
: string
| set transaction read only, as of timestamp now()
< error
| 1064 42000 …line 1 column 32 near "of timestamp now()" 
> exec ^
: string
| set transaction read only as of timestamp now(), isolation level serializable
< ok 0 0
> exec ^
: string
| set transaction isolation level serializable, read only as of timestamp now()
< ok 0 0
> exec ^
: string
| set transaction read only as of timestamp now(), read write
< ok 0 0
> exec ^
: string
| set transaction read write, read only as of timestamp now()
< ok 0 0
> exec ^
: string
| set transaction read only as of timestamp 1, read only as of timestamp 2
< ok 0 0
> exec ^
: string
| set transaction read only as of timestamp now(),
< error
| 1064 42000 …line 1 column 48 near "" 
> exec ^
: string
| set transaction read only as of timestamp 1, 2
< error
| 1064 42000 …line 1 column 46 near "2" 
> exec ^
: string
| set session transaction isolation level read committed, read only as of timestamp now()
< ok 0 0
> exec ^
: string
| set global transaction read only as of timestamp now(), isolation level repeatable read
< ok 0 0
> exec ^
: string
| set transaction
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set session transaction
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| set global transaction
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| set transaction read
< error
| 1064 42000 …line 1 column 20 near "" 
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
| set transaction read committed
< error
| 1064 42000 …line 1 column 30 near "committed" 
> exec ^
: string
| set transaction only
< error
| 1064 42000 …line 1 column 20 near "only" 
> exec ^
: string
| set transaction write
< error
| 1064 42000 …line 1 column 21 near "write" 
> exec ^
: string
| set transaction isolation
< error
| 1064 42000 …line 1 column 25 near "" 
> exec ^
: string
| set transaction isolation level
< error
| 1064 42000 …line 1 column 31 near "" 
> exec ^
: string
| set transaction level serializable
< error
| 1064 42000 …line 1 column 21 near "level serializable" 
> exec ^
: string
| set transaction isolation serializable
< error
| 1064 42000 …line 1 column 38 near "serializable" 
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
| set transaction isolation level repeatable
< error
| 1064 42000 …line 1 column 42 near "" 
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
| set transaction isolation level read
< error
| 1064 42000 …line 1 column 36 near "" 
> exec ^
: string
| set transaction isolation level read only
< error
| 1064 42000 …line 1 column 41 near "only" 
> exec ^
: string
| set transaction isolation level committed
< error
| 1064 42000 …line 1 column 41 near "committed" 
> exec ^
: string
| set transaction isolation level snapshot
< error
| 1064 42000 …line 1 column 40 near "snapshot" 
> exec ^
: string
| set transaction isolation level `serializable`
< error
| 1064 42000 …line 1 column 46 near "`serializable`" 
> exec ^
: string
| set transaction isolation level serializable, read only
< ok 0 0
> exec ^
: string
| set transaction read only, isolation level serializable
< ok 0 0
> exec ^
: string
| set transaction read only, read write
< ok 0 0
> exec ^
: string
| set transaction isolation level serializable, isolation level read committed
< ok 0 0
> exec ^
: string
| set transaction read only,
< error
| 1064 42000 …line 1 column 26 near "" 
> exec ^
: string
| set transaction, read only
< error
| 1064 42000 …line 1 column 16 near ", read only" 
> exec ^
: string
| set transaction read only read write
< error
| 1064 42000 …line 1 column 30 near "read write" 
> exec ^
: string
| set transaction isolation level serializable read only
< error
| 1064 42000 …line 1 column 49 near "read only" 
> exec ^
: string
| set transaction read only foo
< error
| 1064 42000 …line 1 column 29 near "foo" 
> exec ^
: string
| set local transaction read only
< error
| 1064 42000 …line 1 column 26 near "read only" 
> exec ^
: string
| set persist transaction read only
< error
| 1064 42000 …line 1 column 23 near "transaction read only" 
> exec ^
: string
| set persist_only transaction read only
< error
| 1064 42000 …line 1 column 28 near "transaction read only" 
> exec ^
: string
| set session transaction isolation level serializable
< ok 0 0
> exec ^
: string
| set global transaction isolation level read committed
< ok 0 0
> exec ^
: string
| set local transaction isolation level serializable
< error
| 1064 42000 …line 1 column 31 near "isolation level serializable" 
> exec ^
: string
| set @@global transaction read only
< error
| 1064 42000 …line 1 column 24 near "transaction read only" 
> exec ^
: string
| set @@transaction = 1
< ok 0 0
> exec ^
: string
| set transaction_isolation = 'READ-COMMITTED'
< ok 0 0
> exec ^
: string
| set persist transaction = 1
< error
| 1064 42000 …line 1 column 23 near "transaction = 1" 
> exec ^
: string
| set a = 1, transaction read only
< error
| 1064 42000 …line 1 column 27 near "read only" 
> exec ^
: string
| set transaction read only, a = 1
< error
| 1064 42000 …line 1 column 28 near "a = 1" 
> exec ^
: string
| set transaction read only, @a = 1
< error
| 1064 42000 …line 1 column 29 near "@a = 1" 
> exec ^
: string
| set transaction read only, transaction read write
< error
| 1064 42000 …line 1 column 38 near "transaction read write" 
> exec ^
: string
| set transaction isolation level serializable = 1
< error
| 1064 42000 …line 1 column 46 near "= 1" 
> exec ^
: string
| set transaction read only as of timestamp a = 1
< error
| 1054 42S22 Unknown column 'a' in ''
> exec ^
: string
| set transaction read only as of timestamp a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| set transaction = 1
< ok 0 0
> exec ^
: string
| set transaction := 1
< ok 0 0
> exec ^
: string
| set session transaction = 1
< ok 0 0
> exec ^
: string
| set session transaction := 1
< ok 0 0
> exec ^
: string
| set global transaction = 1
< ok 0 0
> exec ^
: string
| set global transaction := 1
< ok 0 0
> exec ^
: string
| set local transaction = 1
< ok 0 0
> exec ^
: string
| set local transaction := 1
< ok 0 0
> exec ^
: string
| set transaction.a = 1
< ok 0 0
> exec ^
: string
| set transaction .a = 1
< ok 0 0
> exec ^
: string
| set transaction. a = 1
< ok 0 0
> exec ^
: string
| set transaction . a = 1
< ok 0 0
> exec ^
: string
| set session transaction.a = 1
< ok 0 0
> exec ^
: string
| set global transaction.a = 1
< ok 0 0
> exec ^
: string
| set local transaction.a = 1
< ok 0 0
> exec ^
: string
| set transaction.transaction = 1
< ok 0 0
> exec ^
: string
| set transaction.isolation = 1
< ok 0 0
> exec ^
: string
| set transaction.select = 1
< ok 0 0
> exec ^
: string
| set transaction.global = 1
< ok 0 0
> exec ^
: string
| set a.transaction = 1
< ok 0 0
> exec ^
: string
| set transaction = 1, a = 2
< ok 0 0
> exec ^
: string
| set a = 1, transaction = 2
< ok 0 0
> exec ^
: string
| set transaction = 1, transaction = 2
< ok 0 0
> exec ^
: string
| set transaction = 1, isolation = 2
< ok 0 0
> exec ^
: string
| set session a = 1, transaction = 2
< ok 0 0
> exec ^
: string
| set global a = 1, transaction = 2
< ok 0 0
> exec ^
: string
| set transaction = 1, a := 2
< ok 0 0
> exec ^
: string
| set session transaction = 1, a = 2
< ok 0 0
> exec ^
: string
| set global transaction = 1, a = 2
< ok 0 0
> exec ^
: string
| set session transaction.a = 1, b = 2
< ok 0 0
> exec ^
: string
| set transaction = (select 1)
< ok 0 0
> exec ^
: string
| set transaction = default
< ok 0 0
> exec ^
: string
| set transaction = a
< ok 0 0
> exec ^
: string
| set transaction = @a
< ok 0 0
> exec ^
: string
| set transaction = 1 + 1
< ok 0 0
> exec ^
: string
| set transaction = 1;
< ok 0 0
> exec ^
: string
| set transaction /* c */ = 1
< ok 0 0
> exec ^
: string
| set /* c */ transaction /* c */ = /* c */ 1
< ok 0 0
> exec ^
: string
| set TRANSACTION = 1
< ok 0 0
> exec ^
: string
| set Transaction.A = 1
< ok 0 0
> exec ^
: string
| set `transaction` = 1
< ok 0 0
> exec ^
: string
| set `transaction`.a = 1
< ok 0 0
> exec ^
: string
| set session `transaction` = 1
< ok 0 0
> exec ^
: string
| set isolation = 1
< ok 0 0
> exec ^
: string
| set isolation := 1
< ok 0 0
> exec ^
: string
| set session isolation = 1
< ok 0 0
> exec ^
: string
| set global isolation = 1
< ok 0 0
> exec ^
: string
| set local isolation = 1
< ok 0 0
> exec ^
: string
| set isolation.a = 1
< ok 0 0
> exec ^
: string
| set a.isolation = 1
< ok 0 0
> exec ^
: string
| set `isolation` = 1
< ok 0 0
> exec ^
: string
| set level = 1
< ok 0 0
> exec ^
: string
| set only = 1
< ok 0 0
> exec ^
: string
| set @transaction = 1
< ok 0 0
> exec ^
: string
| set @@session.transaction = 1
< ok 0 0
> exec ^
: string
| set @@global.transaction = 1
< ok 0 0
> exec ^
: string
| set @@local.transaction = 1
< ok 0 0
> exec ^
: string
| set @isolation = 1
< ok 0 0
> exec ^
: string
| set @@isolation = 1
< ok 0 0
> exec ^
: string
| trace set transaction = 1
< ok 0 0
> exec ^
: string
| trace set transaction read only
< ok 0 0
> exec ^
: string
| set transaction read only, transaction = 1
< error
| 1064 42000 …line 1 column 38 near "transaction = 1" 
> exec ^
: string
| set transaction isolation level serializable, a = 1
< error
| 1064 42000 …line 1 column 47 near "a = 1" 
> exec ^
: string
| set a = 1, transaction isolation level serializable
< error
| 1064 42000 …line 1 column 32 near "isolation level serializable" 
> exec ^
: string
| set a = 1, session transaction read only
< error
| 1064 42000 …line 1 column 35 near "read only" 
> exec ^
: string
| set a = 1, global transaction read only
< error
| 1064 42000 …line 1 column 34 near "read only" 
> exec ^
: string
| set transaction = 1, read only
< error
| 1064 42000 …line 1 column 25 near "read only" 
> exec ^
: string
| set transaction = 1 read only
< error
| 1064 42000 …line 1 column 24 near "read only" 
> exec ^
: string
| set transaction.a.b = 1
< error
| 1064 42000 …line 1 column 18 near ".b = 1" 
> exec ^
: string
| set transaction.a
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set transaction.
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set transaction.a read only
< error
| 1064 42000 …line 1 column 22 near "read only" 
> exec ^
: string
| set transaction =
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set transaction :=
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set transaction == 1
< error
| 1064 42000 …line 1 column 18 near "= 1" 
> exec ^
: string
| set transaction != 1
< error
| 1064 42000 …line 1 column 18 near "!= 1" 
> exec ^
: string
| set transaction < 1
< error
| 1064 42000 …line 1 column 17 near "< 1" 
> exec ^
: string
| set transaction 1
< error
| 1064 42000 …line 1 column 17 near "1" 
> exec ^
: string
| set transaction (1)
< error
| 1064 42000 …line 1 column 17 near "(1)" 
> exec ^
: string
| set transaction, a = 1
< error
| 1064 42000 …line 1 column 16 near ", a = 1" 
> exec ^
: string
| set transaction;
< error
| 1064 42000 …line 1 column 16 near ";" 
> exec ^
: string
| set `transaction` read only
< error
| 1064 42000 …line 1 column 22 near "read only" 
> exec ^
: string
| set isolation level = 1
< error
| 1064 42000 …line 1 column 19 near "level = 1" 
> exec ^
: string
| set isolation level serializable
< error
| 1064 42000 …line 1 column 19 near "level serializable" 
> exec ^
: string
| set read = 1
< error
| 1064 42000 …line 1 column 8 near "read = 1" 
> exec ^
: string
| set read only = 1
< error
| 1064 42000 …line 1 column 8 near "read only = 1" 
> exec ^
: string
| set write = 1
< error
| 1064 42000 …line 1 column 9 near "write = 1" 
> exec ^
: string
| set session session transaction = 1
< error
| 1064 42000 …line 1 column 31 near "transaction = 1" 
> exec ^
: string
| set global session transaction = 1
< error
| 1064 42000 …line 1 column 30 near "transaction = 1" 
> exec ^
: string
| set session global transaction read only
< error
| 1064 42000 …line 1 column 30 near "transaction read only" 
> exec ^
: string
| set persist transaction := 1
< error
| 1064 42000 …line 1 column 23 near "transaction := 1" 
> exec ^
: string
| set persist_only transaction = 1
< error
| 1064 42000 …line 1 column 28 near "transaction = 1" 
> exec ^
: string
| set persist transaction.a = 1
< error
| 1064 42000 …line 1 column 23 near "transaction.a = 1" 
> exec ^
: string
| set persist isolation = 1
< error
| 1064 42000 …line 1 column 21 near "isolation = 1" 
> exec ^
: string
| set persist a = 1
< error
| 1064 42000 …line 1 column 13 near "a = 1" 
> exec ^
: string
| set persist_only a = 1
< error
| 1064 42000 …line 1 column 18 near "a = 1" 
> exec ^
: string
| set a = 1, persist b = 2
< error
| 1064 42000 …line 1 column 20 near "b = 2" 
> exec ^
: string
| set persist = 1
< ok 0 0
> exec ^
: string
| set persist := 1
< ok 0 0
> exec ^
: string
| set persist_only = 1
< ok 0 0
> exec ^
: string
| set persist_only := 1
< ok 0 0
> exec ^
: string
| set persist.a = 1
< ok 0 0
> exec ^
: string
| set persist_only.a = 1
< ok 0 0
> exec ^
: string
| set session persist = 1
< ok 0 0
> exec ^
: string
| set global persist = 1
< ok 0 0
> exec ^
: string
| set a = 1, persist = 2
< ok 0 0
> exec ^
: string
| set a = 1, persist_only = 2
< ok 0 0
> exec ^
: string
| set @@persist.a = 1
< ok 0 0
> exec ^
: string
| set @@persist_only.a = 1
< ok 0 0
> exec ^
: string
| set global = 1
< ok 0 0
> exec ^
: string
| set session = 1
< ok 0 0
> exec ^
: string
| set local = 1
< ok 0 0
> exec ^
: string
| set global := 1
< ok 0 0
> exec ^
: string
| set session := 1
< ok 0 0
> exec ^
: string
| set global.a = 1
< ok 0 0
> exec ^
: string
| set session.a = 1
< ok 0 0
> exec ^
: string
| set local.a = 1
< ok 0 0
> exec ^
: string
| set global . a = 1
< ok 0 0
> exec ^
: string
| set session.transaction = 1
< ok 0 0
> exec ^
: string
| set global.transaction = 1
< ok 0 0
> exec ^
: string
| set session global = 1
< ok 0 0
> exec ^
: string
| set session session = 1
< ok 0 0
> exec ^
: string
| set session local = 1
< ok 0 0
> exec ^
: string
| set local session = 1
< ok 0 0
> exec ^
: string
| set global = 1, session = 2
< ok 0 0
> exec ^
: string
| set session = (select 1)
< ok 0 0
> exec ^
: string
| set a.global = 1
< ok 0 0
> exec ^
: string
| set a.session = 1
< ok 0 0
> exec ^
: string
| set a.local = 1
< ok 0 0
> exec ^
: string
| set a . global = 1
< ok 0 0
> exec ^
: string
| set global.a.b = 1
< error
| 1064 42000 …line 1 column 13 near ".b = 1" 
> exec ^
: string
| set names = 1
< ok 0 0
> exec ^
: string
| set names := 1
< ok 0 0
> exec ^
: string
| set names.a = 1
< ok 0 0
> exec ^
: string
| set names utf8
< ok 0 0
> exec ^
: string
| set names default
< ok 0 0
> exec ^
: string
| set charset = 1
< ok 0 0
> exec ^
: string
| set charset := 1
< ok 0 0
> exec ^
: string
| set charset.a = 1
< ok 0 0
> exec ^
: string
| set charset utf8
< ok 0 0
> exec ^
: string
| set charset default
< ok 0 0
> exec ^
: string
| set character = 1
< error
| 1064 42000 …line 1 column 15 near "= 1" 
> exec ^
: string
| set character set utf8
< ok 0 0
> exec ^
: string
| set role = 1
< ok 0 0
> exec ^
: string
| set role := 1
< ok 0 0
> exec ^
: string
| set role.a = 1
< ok 0 0
> exec ^
: string
| set role all
< ok 0 0
> exec ^
: string
| set role none
< ok 0 0
> exec ^
: string
| set resource = 1
< ok 0 0
> exec ^
: string
| set resource := 1
< ok 0 0
> exec ^
: string
| set resource.a = 1
< ok 0 0
> exec ^
: string
| set resource group rg
< ok 0 0
> exec ^
: string
| set binding = 1
< ok 0 0
> exec ^
: string
| set config = 1
< ok 0 0
> exec ^
: string
| set default = 1
< error
| 1064 42000 …line 1 column 13 near "= 1" 
> exec ^
: string
| set default.a = 1
< ok 0 0
> exec ^
: string
| set default role none to u
< ok 0 0
> exec ^
: string
| set password = 1
< error
| 1064 42000 …line 1 column 16 near "1" 
> exec ^
: string
| set password := 1
< ok 0 0
> exec ^
: string
| set password.a = 1
< ok 0 0
> exec ^
: string
| set password = 'x'
< ok 0 0
> exec ^
: string
| set password = password('x')
< ok 0 0
> exec ^
: string
| set password for u = 'x'
< ok 0 0
> exec ^
: string
| set session password = 1
< ok 0 0
> exec ^
: string
| set timestamp = 1
< ok 0 0
> exec ^
: string
| set time = 1
< ok 0 0
> exec ^
: string
| set date = 1
< ok 0 0
> exec ^
: string
| set user = 1
< ok 0 0
> exec ^
: string
| set status = 1
< ok 0 0
> exec ^
: string
| set session timestamp = 1
< ok 0 0
> exec ^
: string
| set a = 1, timestamp = 2
< ok 0 0
> exec ^
: string
| set timestamp.a = 1
< ok 0 0
> exec ^
: string
| set a.timestamp = 1
< ok 0 0
> exec ^
: string
| set a.select = 1
< ok 0 0
> exec ^
: string
| set select.a = 1
< ok 0 0
> exec ^
: string
| set a . transaction = 1
< ok 0 0
> exec ^
: string
| select transaction from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select transaction from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select isolation from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select isolation from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select transaction.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select transaction.a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a.transaction from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a.transaction from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a transaction from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a transaction from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a as transaction from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a as transaction from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from transaction
< error
| 1146 42S02 Table 'test.transaction' doesn't exist
> exec ^
: string
| select * from transaction
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t transaction
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t transaction
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select * from t as isolation
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select * from t as isolation
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table transaction (transaction int, isolation int)
< ok 0 0
> exec ^
: string
| update t set transaction = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t set transaction = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t (transaction) values (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t (transaction) values (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create role transaction
< error
| 1064 42000 …line 1 column 23 near "" 
> exec ^
: string
| drop role transaction
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| grant transaction to u
< error
| 1064 42000 …line 1 column 20 near "to u" 
> exec ^
: string
| grant isolation to u
< error
| 1064 42000 …line 1 column 18 near "to u" 
> exec ^
: string
| revoke transaction from u
< error
| 1064 42000 …line 1 column 23 near "from u" 
> exec ^
: string
| set role transaction
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| set role isolation
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set default role transaction to u
< error
| 1064 42000 …line 1 column 31 near "to u" 
> exec ^
: string
| grant select on t to u as a with role transaction
< error
| 1064 42000 …line 1 column 25 near "as a with role transaction" 
> exec ^
: string
| grant select on t to u as a with role isolation
< error
| 1064 42000 …line 1 column 25 near "as a with role isolation" 
> exec ^
: string
| reset persist transaction
< error
| 1064 42000 …line 1 column 5 near "reset persist transaction" 
> exec ^
: string
| reset persist isolation
< error
| 1064 42000 …line 1 column 5 near "reset persist isolation" 
> exec ^
: string
| fetch c into transaction
< error
| 1064 42000 …line 1 column 5 near "fetch c into transaction" 
> exec ^
: string
| set array = 1
< error
| 1064 42000 …line 1 column 9 near "array = 1" 
> exec ^
: string
| set until = 1
< error
| 1064 42000 …line 1 column 9 near "until = 1" 
> exec ^
: string
| set session array = 1
< error
| 1064 42000 …line 1 column 17 near "array = 1" 
> exec ^
: string
| set array.a = 1
< ok 0 0
> exec ^
: string
| set until.a = 1
< ok 0 0
> exec ^
: string
| set a.array = 1
< ok 0 0
> exec ^
: string
| set a.until = 1
< ok 0 0
> exec ^
: string
| select array from t
< error
| 1064 42000 …line 1 column 12 near "array from t" 
> exec ^
: string
| select until from t
< error
| 1064 42000 …line 1 column 12 near "until from t" 
> exec ^
: string
| select array.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select array.a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select until.a from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select until.a from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a.array from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a.array from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a.until from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select a.until from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table array (a int)
< error
| 1064 42000 …line 1 column 18 near "array (a int)" 
> exec ^
: string
| create table until (a int)
< error
| 1064 42000 …line 1 column 18 near "until (a int)" 
> exec ^
: string
| select cast(a as unsigned array) from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select cast(a as unsigned array) from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a json, index i ((cast(a as unsigned array))))
< ok 0 0
> exec ^
: string
| set a = 1, global b = 2
< ok 0 0
> exec ^
: string
| set a = 1, session b = 2
< ok 0 0
> exec ^
: string
| set a = 1, local b = 2
< ok 0 0
> exec ^
: string
| set global a = 1, session b = 2
< ok 0 0
> exec ^
: string
| set global a = 1, global b = 2
< ok 0 0
> exec ^
: string
| set session a = 1, local b = 2
< ok 0 0
> exec ^
: string
| set local a = 1, global b = 2
< ok 0 0
> exec ^
: string
| set global a = 1, b = 2
< ok 0 0
> exec ^
: string
| set a = 1, b = 2, global c = 3
< ok 0 0
> exec ^
: string
| set a = 1, session transaction = 2
< ok 0 0
> exec ^
: string
| set a = 1, global transaction = 2
< ok 0 0
> exec ^
: string
| set a = 1, local transaction = 2
< ok 0 0
> exec ^
: string
| set transaction = 1, session transaction = 2
< ok 0 0
> exec ^
: string
| set transaction = 1, global transaction = 2
< ok 0 0
> exec ^
: string
| set transaction = 1, local transaction = 2
< ok 0 0
> exec ^
: string
| set a = 1, global b := 2
< ok 0 0
> exec ^
: string
| set a = 1, global b.c = 2
< ok 0 0
> exec ^
: string
| set a = 1, global
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set a = 1, global =
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| set a = 1, global = 2
< ok 0 0
> exec ^
: string
| set a = 1, session = 2
< ok 0 0
> exec ^
: string
| set a = 1, local := 2
< ok 0 0
> exec ^
: string
| set a = 1, global.b = 2
< ok 0 0
> exec ^
: string
| set a = 1, global global = 2
< ok 0 0
> exec ^
: string
| set a = 1, global session = 2
< ok 0 0
> exec ^
: string
| set a = 1, global session b = 2
< error
| 1064 42000 …line 1 column 27 near "b = 2" 
> exec ^
: string
| set a = 1, persist_only b = 2
< error
| 1064 42000 …line 1 column 25 near "b = 2" 
> exec ^
: string
| set a = 1, global @b = 2
< error
| 1064 42000 …line 1 column 20 near "@b = 2" 
> exec ^
: string
| set a = 1, global @@b = 2
< error
| 1064 42000 …line 1 column 21 near "@@b = 2" 
> exec ^
: string
| set a = 1, session @@b = 2
< error
| 1064 42000 …line 1 column 22 near "@@b = 2" 
> exec ^
: string
| set global @a = 1
< error
| 1064 42000 …line 1 column 13 near "@a = 1" 
> exec ^
: string
| set session @a = 1
< error
| 1064 42000 …line 1 column 14 near "@a = 1" 
> exec ^
: string
| set local @a = 1
< error
| 1064 42000 …line 1 column 12 near "@a = 1" 
> exec ^
: string
| set global @@a = 1
< error
| 1064 42000 …line 1 column 14 near "@@a = 1" 
> exec ^
: string
| set session @@a = 1
< error
| 1064 42000 …line 1 column 15 near "@@a = 1" 
> exec ^
: string
| set global @@global.a = 1
< error
| 1064 42000 …line 1 column 21 near "@@global.a = 1" 
> exec ^
: string
| set global names utf8
< error
| 1064 42000 …line 1 column 21 near "utf8" 
> exec ^
: string
| set session names utf8
< error
| 1064 42000 …line 1 column 22 near "utf8" 
> exec ^
: string
| set local names utf8
< error
| 1064 42000 …line 1 column 20 near "utf8" 
> exec ^
: string
| set global charset utf8
< error
| 1064 42000 …line 1 column 23 near "utf8" 
> exec ^
: string
| set global character set utf8
< error
| 1064 42000 …line 1 column 20 near "character set utf8" 
> exec ^
: string
| set global names = 1
< ok 0 0
> exec ^
: string
| set global charset = 1
< ok 0 0
> exec ^
: string
| set a = 1, global names utf8
< error
| 1064 42000 …line 1 column 28 near "utf8" 
> exec ^
: string
| set a = on
< ok 0 0
> exec ^
: string
| set a = off
< ok 0 0
> exec ^
: string
| set a = on, b = 2
< ok 0 0
> exec ^
: string
| set a = on;
< ok 0 0
> exec ^
: string
| set a := on
< ok 0 0
> exec ^
: string
| set global a = on
< ok 0 0
> exec ^
: string
| set session a = on
< ok 0 0
> exec ^
: string
| set @@a = on
< ok 0 0
> exec ^
: string
| set @@global.a = on
< ok 0 0
> exec ^
: string
| set @a = on
< error
| 1064 42000 …line 1 column 11 near "on" 
> exec ^
: string
| set @a := on
< error
| 1064 42000 …line 1 column 12 near "on" 
> exec ^
: string
| set a = on + 1
< error
| 1064 42000 …line 1 column 12 near "+ 1" 
> exec ^
: string
| set a = (on)
< error
| 1064 42000 …line 1 column 11 near "on)" 
> exec ^
: string
| set a = on b
< error
| 1064 42000 …line 1 column 12 near "b" 
> exec ^
: string
| set a = on on
< error
| 1064 42000 …line 1 column 13 near "on" 
> exec ^
: string
| set a = 1, b = on
< ok 0 0
> exec ^
: string
| set transaction = on
< ok 0 0
> exec ^
: string
| set a.b = on
< ok 0 0
> exec ^
: string
| set a = binary
< ok 0 0
> exec ^
: string
| set a = binary, b = 2
< ok 0 0
> exec ^
: string
| set a = binary;
< ok 0 0
> exec ^
: string
| set a := binary
< ok 0 0
> exec ^
: string
| set global a = binary
< ok 0 0
> exec ^
: string
| set @@a = binary
< ok 0 0
> exec ^
: string
| set @@character_set_results = binary
< ok 0 0
> exec ^
: string
| set @@global.a = binary
< ok 0 0
> exec ^
: string
| set @a = binary
< error
| 1064 42000 …line 1 column 15 near "" 
> exec ^
: string
| set a = binary 'x'
< ok 0 0
> exec ^
: string
| set a = binary b
< error
| 1054 42S22 Unknown column 'b' in ''
> exec ^
: string
| set a = binary b
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| set a = binary 1 + 1
< ok 0 0
> exec ^
: string
| set a = binary + 1
< ok 0 0
> exec ^
: string
| set a = binary b, c = 2
< error
| 1054 42S22 Unknown column 'b' in ''
> exec ^
: string
| set a = binary b, c = 2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| set @a = binary 'x'
< ok 0 0
> exec ^
: string
| set a = (binary)
< error
| 1064 42000 …line 1 column 16 near ")" 
> exec ^
: string
| set a = binary binary
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| set transaction = binary
< ok 0 0
> exec ^
: string
| set a = default
< ok 0 0
> exec ^
: string
| set @@a = default
< ok 0 0
> exec ^
: string
| set @a = default
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set @a := default
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set a = default, b = 2
< ok 0 0
> exec ^
: string
| set a = default + 1
< error
| 1064 42000 …line 1 column 17 near "+ 1" 
> exec ^
: string
| set a = true
< ok 0 0
> exec ^
: string
| set a = null
< ok 0 0
> exec ^
: string
| set a = 'on'
< ok 0 0
> exec ^
: string
| set a = all
< error
| 1064 42000 …line 1 column 11 near "all" 
> exec ^
: string
| set @a = 1, global b = 2
< ok 0 0
> exec ^
: string
| set @@a = 1, session b = 2
< ok 0 0
> exec ^
: string
| set @@a = 1, b = on
< ok 0 0
> exec ^
: string
| set names 'utf8'
< ok 0 0
> exec ^
: string
| set names "utf8"
< ok 0 0
> exec ^
: string
| set names binary
< ok 0 0
> exec ^
: string
| set names `utf8`
< ok 0 0
> exec ^
: string
| set names 1
< error
| 1064 42000 …line 1 column 11 near "1" 
> exec ^
: string
| set names @a
< error
| 1064 42000 …line 1 column 12 near "@a" 
> exec ^
: string
| set names
< error
| 1064 42000 …line 1 column 9 near "" 
> exec ^
: string
| set names default collate utf8_bin
< error
| 1064 42000 …line 1 column 25 near "collate utf8_bin" 
> exec ^
: string
| set names utf8 collate utf8_bin
< ok 0 0
> exec ^
: string
| set names utf8 collate 'utf8_bin'
< ok 0 0
> exec ^
: string
| set names 'utf8' collate 'utf8_bin'
< ok 0 0
> exec ^
: string
| set names utf8 collate default
< ok 0 0
> exec ^
: string
| set names utf8 collate binary
< error
| 1064 42000 …line 1 column 29 near "binary" 
> exec ^
: string
| set names binary collate binary
< error
| 1064 42000 …line 1 column 31 near "binary" 
> exec ^
: string
| set names binary collate default
< ok 0 0
> exec ^
: string
| set names utf8 collate
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| set names utf8 collate 1
< error
| 1064 42000 …line 1 column 24 near "1" 
> exec ^
: string
| set names utf8 utf8
< error
| 1064 42000 …line 1 column 19 near "utf8" 
> exec ^
: string
| set names utf8, a = 1
< ok 0 0
> exec ^
: string
| set names 'utf8', a = 1
< ok 0 0
> exec ^
: string
| set names utf8 collate utf8_bin, a = 1
< ok 0 0
> exec ^
: string
| set names default, a = 1
< ok 0 0
> exec ^
: string
| set names binary, a = 1
< ok 0 0
> exec ^
: string
| set names utf8, @@session.sql_mode=1;
< ok 0 0
> exec ^
: string
| set @@session.sql_mode=1, names utf8, charset utf8;
< ok 0 0
> exec ^
: string
| set a = 1, names utf8
< ok 0 0
> exec ^
: string
| set a = 1, names 'utf8'
< ok 0 0
> exec ^
: string
| set a = 1, names default
< ok 0 0
> exec ^
: string
| set a = 1, names binary
< ok 0 0
> exec ^
: string
| set a = 1, names utf8 collate utf8_bin
< ok 0 0
> exec ^
: string
| set a = 1, names = 2
< ok 0 0
> exec ^
: string
| set a = 1, names := 2
< ok 0 0
> exec ^
: string
| set a = 1, names.b = 2
< ok 0 0
> exec ^
: string
| set a = 1, names
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| set names utf8, names latin1
< ok 0 0
> exec ^
: string
| set names utf8, charset utf8
< ok 0 0
> exec ^
: string
| set names = 1, names utf8
< ok 0 0
> exec ^
: string
| set charset 'utf8'
< ok 0 0
> exec ^
: string
| set charset binary
< ok 0 0
> exec ^
: string
| set charset
< error
| 1064 42000 …line 1 column 11 near "" 
> exec ^
: string
| set charset 1
< error
| 1064 42000 …line 1 column 13 near "1" 
> exec ^
: string
| set charset utf8 collate utf8_bin
< error
| 1064 42000 …line 1 column 24 near "collate utf8_bin" 
> exec ^
: string
| set charset utf8, a = 1
< ok 0 0
> exec ^
: string
| set a = 1, charset utf8
< ok 0 0
> exec ^
: string
| set a = 1, charset 'utf8'
< ok 0 0
> exec ^
: string
| set a = 1, charset default
< ok 0 0
> exec ^
: string
| set a = 1, charset binary
< ok 0 0
> exec ^
: string
| set a = 1, charset = 2
< ok 0 0
> exec ^
: string
| set a = 1, charset
< error
| 1064 42000 …line 1 column 18 near "" 
> exec ^
: string
| set character set 'utf8mb4';
< ok 0 0
> exec ^
: string
| SET CHARACTER SET 'utf8mb4';
< ok 0 0
> exec ^
: string
| set character set binary
< ok 0 0
> exec ^
: string
| set character set default
< ok 0 0
> exec ^
: string
| set character set
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| set character set utf8 collate utf8_bin
< error
| 1064 42000 …line 1 column 30 near "collate utf8_bin" 
> exec ^
: string
| set character set utf8, a = 1
< ok 0 0
> exec ^
: string
| set a = 1, character set utf8
< ok 0 0
> exec ^
: string
| set a = 1, character set default
< ok 0 0
> exec ^
: string
| set a = 1, character = 2
< error
| 1064 42000 …line 1 column 22 near "= 2" 
> exec ^
: string
| set a = 1, character
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| set char set default
< ok 0 0
> exec ^
: string
| set char set utf8
< ok 0 0
> exec ^
: string
| set char set 'utf8'
< ok 0 0
> exec ^
: string
| set char set binary
< ok 0 0
> exec ^
: string
| set char set
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| set char set utf8, a = 1
< ok 0 0
> exec ^
: string
| set a = 1, char set utf8
< ok 0 0
> exec ^
: string
| set a = 1, char set default
< ok 0 0
> exec ^
: string
| set char = 1
< error
| 1064 42000 …line 1 column 10 near "= 1" 
> exec ^
: string
| set char utf8
< error
| 1064 42000 …line 1 column 13 near "utf8" 
> exec ^
: string
| set a = 1, char = 2
< error
| 1064 42000 …line 1 column 17 near "= 2" 
> exec ^
: string
| set charset set utf8
< error
| 1064 42000 …line 1 column 15 near "set utf8" 
> exec ^
: string
| set character utf8
< error
| 1064 42000 …line 1 column 18 near "utf8" 
> exec ^
: string
| set a = 1, role = 2
< ok 0 0
> exec ^
: string
| set a = 1, role b
< error
| 1064 42000 …line 1 column 17 near "b" 
> exec ^
: string
| set a = 1, role all
< error
| 1064 42000 …line 1 column 19 near "all" 
> exec ^
: string
| set a = 1, default role none to u
< error
| 1064 42000 …line 1 column 18 near "default role none to u" 
> exec ^
: string
| set a = 1, password = 2
< ok 0 0
> exec ^
: string
| set a = 1, password := 2
< ok 0 0
> exec ^
: string
| set a = 1, password = 'x'
< ok 0 0
> exec ^
: string
| set a = 1, password for u = 'x'
< error
| 1064 42000 …line 1 column 23 near "for u = 'x'" 
> exec ^
: string
| set a = 1, resource = 2
< ok 0 0
> exec ^
: string
| set a = 1, resource group rg
< error
| 1064 42000 …line 1 column 25 near "group rg" 
> exec ^
: string
| set a = 1, binding = 2
< ok 0 0
> exec ^
: string
| set a = 1, binding enabled for select 1
< error
| 1064 42000 …line 1 column 26 near "enabled for select 1" 
> exec ^
: string
| set a = 1, config = 2
< ok 0 0
> exec ^
: string
| set role all, a = 1
< error
| 1064 42000 …line 1 column 13 near ", a = 1" 
> exec ^
: string
| set password = 'x', a = 1
< error
| 1064 42000 …line 1 column 19 near ", a = 1" 
> exec ^
: string
| set resource group rg, a = 1
< error
| 1064 42000 …line 1 column 22 near ", a = 1" 
> exec ^
: string
| set default role none to u, a = 1
< error
| 1064 42000 …line 1 column 31 near "= 1" 
> exec ^
: string
| set a = 1,
< error
| 1064 42000 …line 1 column 10 near "" 
> exec ^
: string
| set a = 1,,
< error
| 1064 42000 …line 1 column 11 near "," 
> exec ^
: string
| set , a = 1
< error
| 1064 42000 …line 1 column 5 near ", a = 1" 
> exec ^
: string
| set a .select = 1
< ok 0 0
> exec ^
: string
| set a . /* c */ select = 1
< error
| 1064 42000 …line 1 column 22 near "select = 1" 
> exec ^
: string
| set a . `select` = 1
< ok 0 0
> exec ^
: string
| set a . 1 = 1
< error
| 1064 42000 …line 1 column 9 near "1 = 1" 
> exec ^
: string
| set a . = 1
< error
| 1064 42000 …line 1 column 9 near "= 1" 
> exec ^
: string
| set a . select . b = 1
< error
| 1064 42000 …line 1 column 16 near ". b = 1" 
> exec ^
: string
| set select . a = 1
< error
| 1064 42000 …line 1 column 10 near "select . a = 1" 
> exec ^
: string
| set select .a = 1
< error
| 1064 42000 …line 1 column 10 near "select .a = 1" 
> exec ^
: string
| set select. a = 1
< ok 0 0
> exec ^
: string
| set @@a . select = 1
< error
| 1064 42000 …line 1 column 9 near ". select = 1" 
> exec ^
: string
| set foreign_key_checks = on
< ok 0 0
> exec ^
: string
| set foreign_key_checks = off
< ok 0 0
> exec ^
: string
| set foreign_key_checks = 'on'
< ok 0 0
> exec ^
: string
| set foreign_key_checks = binary
< ok 0 0
> exec ^
: string
| set @@session.foreign_key_checks = 1, global foreign_key_checks = 0
< ok 0 0
> exec ^
: string
| set global foreign_key_checks = 0, session foreign_key_checks = 1
< ok 0 0
> exec ^
: string
| trace set a = on
< ok 0 0
> exec ^
: string
| trace set a = 1, global b = 2
< ok 0 0
> exec ^
: string
| trace set names 'utf8'
< ok 0 0
> exec ^
: string
| set a.@b = 1
< error
| 1064 42000 …line 1 column 8 near "@b = 1" 
> exec ^
: string
| set a = default(b)
< ok 0 0
> exec ^
: string
| set @a = default(b)
< ok 0 0
> exec ^
: string
| set a = default(b) + 1
< error
| 1054 42S22 Unknown column 'b' in 'field list'
> exec ^
: string
| set a = default(b) + 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| set a = binary)
< error
| 1064 42000 …line 1 column 15 near ")" 
> exec ^
: string
| set a = on)
< error
| 1064 42000 …line 1 column 11 near ")" 
> exec ^
: string
| set @@a.b = 1
< ok 0 0
> exec ^
: string
| set global a = default
< ok 0 0
> exec ^
: string
| set a = 1, b = binary
< ok 0 0
> exec ^
: string
| set a = 1, b = default
< ok 0 0
> exec ^
: string
| set names utf8 collate default, a = 1
< ok 0 0
> exec ^
: string
| set names binary collate default, a = 1
< ok 0 0
> exec ^
: string
| set a = 1, names utf8 collate 'utf8_bin', b = 2
< ok 0 0
> exec ^
: string
| set a = 1, global b = 2, names utf8, session c = on, charset binary, @d = 1, @@e = binary
< ok 0 0
> exec ^
: string
| set NAMES 'utf8'
< ok 0 0
> exec ^
: string
| set Names Utf8 Collate Utf8_Bin
< ok 0 0
> exec ^
: string
| set CHAR SET utf8
< ok 0 0
> exec ^
: string
| set /* c */ names /* c */ 'utf8'
< ok 0 0
> exec ^
: string
| set char /* c */ set utf8
< ok 0 0
> exec ^
: string
| set global.names = 1
< ok 0 0
> exec ^
: string
| set session names = 1
< ok 0 0
> exec ^
: string
| set session char set utf8
< error
| 1064 42000 …line 1 column 16 near "char set utf8" 
> exec ^
: string
| set session binary = 1
< error
| 1064 42000 …line 1 column 18 near "binary = 1" 
> exec ^
: string
| set on = 1
< error
| 1064 42000 …line 1 column 6 near "on = 1" 
> exec ^
: string
| set binary = 1
< error
| 1064 42000 …line 1 column 10 near "binary = 1" 
> exec ^
: string
| set a = on, on = 1
< error
| 1064 42000 …line 1 column 14 near "on = 1" 
> exec ^
: string
| create procedure p() begin set a = binary; end
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| create procedure p() begin set a = binary; end
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create procedure p() begin set a = binary; select 1; end
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| create procedure p() begin set a = binary; select 1; end
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create procedure p() begin set a = on; end
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| create procedure p() begin set a = on; end
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create procedure p() begin set a = 1, global b = 2; end
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| create procedure p() begin set a = 1, global b = 2; end
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create procedure p() begin set names 'utf8'; end
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| create procedure p() begin set names 'utf8'; end
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create procedure p() begin set @a = default; end
< error
| 1064 42000 …line 1 column 44 near "; end" 
> exec ^
: string
| create procedure p() begin set a = binary end
< error
| 1064 42000 …line 1 column 45 near "" 
> exec ^
: string
| start transaction read only as of timestamp now()
< ok 0 0
> exec ^
: string
| start transaction read only as of
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| start transaction read only as of timestamp
< error
| 1064 42000 …line 1 column 43 near "" 
> exec ^
: string
| start transaction read only as now()
< error
| 1064 42000 …line 1 column 30 near "as now()" 
> exec ^
: string
| start transaction read only as timestamp now()
< error
| 1064 42000 …line 1 column 30 near "as timestamp now()" 
> exec ^
: string
| start transaction read only as `of` timestamp now()
< error
| 1064 42000 …line 1 column 30 near "as `of` timestamp now()" 
> exec ^
: string
| start transaction read only of timestamp now()
< error
| 1064 42000 …line 1 column 30 near "of timestamp now()" 
> exec ^
: string
| start transaction read only timestamp now()
< error
| 1064 42000 …line 1 column 37 near "timestamp now()" 
> exec ^
: string
| start transaction read only as of timestamp now(), read write
< error
| 1064 42000 …line 1 column 50 near ", read write" 
