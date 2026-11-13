-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestSetTransactionAsOfFormsOracle default#1
> exec
| PREPARE omni_parses_stmt FROM ?
: string
| set transaction read only as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| set transaction read only as of timestamp '2021-04-21 00:42:12'
< error
| 1064 42000 …near 'as of timestamp '2021-04-21 00:42:12'' at line 1
> exec ^
: string
| SET TRANSACTION READ ONLY AS OF TIMESTAMP NOW()
< error
| 1064 42000 …near 'AS OF TIMESTAMP NOW()' at line 1
> exec ^
: string
| set session transaction read only as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| set global transaction read only as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| set local transaction read only as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| set persist transaction read only as of timestamp now()
< error
| 1064 42000 …near 'transaction read only as of timestamp now()' at line 1
> exec ^
: string
| set persist_only transaction read only as of timestamp now()
< error
| 1064 42000 …near 'transaction read only as of timestamp now()' at line 1
> exec ^
: string
| set @@session transaction read only as of timestamp now()
< error
| 1064 42000 …near 'transaction read only as of timestamp now()' at line 1
> exec ^
: string
| set transaction read only as of timestamp now() - interval 3 second
< error
| 1064 42000 …near 'as of timestamp now() - interval 3 second' at line 1
> exec ^
: string
| set transaction read only as of timestamp tidb_bounded_staleness(now(), now())
< error
| 1064 42000 …near 'as of timestamp tidb_bounded_staleness(now(), now())' at line 1
> exec ^
: string
| set transaction read only as of timestamp (select now())
< error
| 1064 42000 …near 'as of timestamp (select now())' at line 1
> exec ^
: string
| set transaction read only as of timestamp a
< error
| 1064 42000 …near 'as of timestamp a' at line 1
> exec ^
: string
| set transaction read only as of timestamp @a
< error
| 1064 42000 …near 'as of timestamp @a' at line 1
> exec ^
: string
| set transaction read only as of timestamp 1 or 2
< error
| 1064 42000 …near 'as of timestamp 1 or 2' at line 1
> exec ^
: string
| set transaction read only as /* c */ of timestamp now()
< error
| 1064 42000 …near 'as /* c */ of timestamp now()' at line 1
> exec ^
: string
| set transaction read only as
| of
| timestamp now()
< error
| 1064 42000 …near 'as
| of
| timestamp now()' at line 1
> exec ^
: string
| set transaction read only as of
< error
| 1064 42000 …near 'as of' at line 1
> exec ^
: string
| set transaction read only as of timestamp
< error
| 1064 42000 …near 'as of timestamp' at line 1
> exec ^
: string
| set transaction read only as of now()
< error
| 1064 42000 …near 'as of now()' at line 1
> exec ^
: string
| set transaction read only as of time now()
< error
| 1064 42000 …near 'as of time now()' at line 1
> exec ^
: string
| set transaction read only as now()
< error
| 1064 42000 …near 'as now()' at line 1
> exec ^
: string
| set transaction read only as timestamp now()
< error
| 1064 42000 …near 'as timestamp now()' at line 1
> exec ^
: string
| set transaction read only of timestamp now()
< error
| 1064 42000 …near 'of timestamp now()' at line 1
> exec ^
: string
| set transaction read only timestamp now()
< error
| 1064 42000 …near 'timestamp now()' at line 1
> exec ^
: string
| set transaction read only as `of` timestamp now()
< error
| 1064 42000 …near 'as `of` timestamp now()' at line 1
> exec ^
: string
| set transaction read only as of `timestamp` now()
< error
| 1064 42000 …near 'as of `timestamp` now()' at line 1
> exec ^
: string
| set transaction read only as of timestamp now() a
< error
| 1064 42000 …near 'as of timestamp now() a' at line 1
> exec ^
: string
| set transaction read only as of timestamp now() as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now() as of timestamp now()' at line 1
> exec ^
: string
| set transaction read only as of timestamp now() read write
< error
| 1064 42000 …near 'as of timestamp now() read write' at line 1
> exec ^
: string
| set transaction read only as of timestamp now() for update
< error
| 1064 42000 …near 'as of timestamp now() for update' at line 1
> exec ^
: string
| set transaction read write as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| set transaction as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| set transaction isolation level serializable as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| set transaction read only, as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| set transaction read only as of timestamp now(), isolation level serializable
< error
| 1064 42000 …near 'as of timestamp now(), isolation level serializable' at line 1
> exec ^
: string
| set transaction isolation level serializable, read only as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| set transaction read only as of timestamp now(), read write
< error
| 1064 42000 …near 'as of timestamp now(), read write' at line 1
> exec ^
: string
| set transaction read write, read only as of timestamp now()
< error
| 1064 42000 …near 'read only as of timestamp now()' at line 1
> exec ^
: string
| set transaction read only as of timestamp 1, read only as of timestamp 2
< error
| 1064 42000 …near 'as of timestamp 1, read only as of timestamp 2' at line 1
> exec ^
: string
| set transaction read only as of timestamp now(),
< error
| 1064 42000 …near 'as of timestamp now(),' at line 1
> exec ^
: string
| set transaction read only as of timestamp 1, 2
< error
| 1064 42000 …near 'as of timestamp 1, 2' at line 1
> exec ^
: string
| set session transaction isolation level read committed, read only as of timestamp now()
< error
| 1064 42000 …near 'as of timestamp now()' at line 1
> exec ^
: string
| set global transaction read only as of timestamp now(), isolation level repeatable read
< error
| 1064 42000 …near 'as of timestamp now(), isolation level repeatable read' at line 1
