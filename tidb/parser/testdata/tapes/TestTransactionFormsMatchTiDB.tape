-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestTransactionFormsMatchTiDB default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| begin
< ok 0 0
> exec ^
: string
| begin pessimistic
< ok 0 0
> exec ^
: string
| begin optimistic
< ok 0 0
> exec ^
: string
| BEGIN PESSIMISTIC
< ok 0 0
> exec ^
: string
| Begin Optimistic
< ok 0 0
> exec ^
: string
| begin /* c */ pessimistic
< ok 0 0
> exec ^
: string
| begin pessimistic;
< ok 0 0
> exec ^
: string
| begin work
< error
| 1064 42000 …line 1 column 10 near "work" 
> exec ^
: string
| begin work;
< error
| 1064 42000 …line 1 column 10 near "work;" 
> exec ^
: string
| begin pessimistic work
< error
| 1064 42000 …line 1 column 22 near "work" 
> exec ^
: string
| begin work pessimistic
< error
| 1064 42000 …line 1 column 10 near "work pessimistic" 
> exec ^
: string
| begin optimistic pessimistic
< error
| 1064 42000 …line 1 column 28 near "pessimistic" 
> exec ^
: string
| begin pessimistic, optimistic
< error
| 1064 42000 …line 1 column 18 near ", optimistic" 
> exec ^
: string
| begin transaction
< error
| 1064 42000 …line 1 column 17 near "transaction" 
> exec ^
: string
| begin read only
< error
| 1064 42000 …line 1 column 10 near "read only" 
> exec ^
: string
| begin with consistent snapshot
< error
| 1064 42000 …line 1 column 10 near "with consistent snapshot" 
> exec ^
: string
| begin `pessimistic`
< error
| 1064 42000 …line 1 column 19 near "`pessimistic`" 
> exec ^
: string
| begin 'pessimistic'
< error
| 1064 42000 …line 1 column 19 near "'pessimistic'" 
> exec ^
: string
| start transaction
< ok 0 0
> exec ^
: string
| start transaction read only
< ok 0 0
> exec ^
: string
| start transaction read write
< ok 0 0
> exec ^
: string
| start transaction with consistent snapshot
< ok 0 0
> exec ^
: string
| start transaction with causal consistency only
< ok 0 0
> exec ^
: string
| START TRANSACTION WITH CAUSAL CONSISTENCY ONLY
< ok 0 0
> exec ^
: string
| start transaction /*!40108 with consistent snapshot */
< ok 0 0
> exec ^
: string
| start
< error
| 1064 42000 …line 1 column 5 near "" 
> exec ^
: string
| start work
< error
| 1064 42000 …line 1 column 10 near "work" 
> exec ^
: string
| start transaction work
< error
| 1064 42000 …line 1 column 22 near "work" 
> exec ^
: string
| start transaction pessimistic
< error
| 1064 42000 …line 1 column 29 near "pessimistic" 
> exec ^
: string
| start transaction optimistic
< error
| 1064 42000 …line 1 column 28 near "optimistic" 
> exec ^
: string
| start transaction isolation level read committed
< error
| 1064 42000 …line 1 column 27 near "isolation level read committed" 
> exec ^
: string
| start transaction read
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| start transaction with
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| start transaction with consistent
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| start transaction with snapshot
< error
| 1064 42000 …line 1 column 31 near "snapshot" 
> exec ^
: string
| start transaction with consistent snapshot only
< error
| 1064 42000 …line 1 column 47 near "only" 
> exec ^
: string
| start transaction with causal
< error
| 1064 42000 …line 1 column 29 near "" 
> exec ^
: string
| start transaction with causal consistency
< error
| 1064 42000 …line 1 column 41 near "" 
> exec ^
: string
| start transaction with causal only
< error
| 1064 42000 …line 1 column 34 near "only" 
> exec ^
: string
| start transaction with `causal` consistency only
< error
| 1064 42000 …line 1 column 31 near "`causal` consistency only" 
> exec ^
: string
| start transaction with causal `consistency` only
< error
| 1064 42000 …line 1 column 43 near "`consistency` only" 
> exec ^
: string
| start transaction read only, with consistent snapshot
< error
| 1064 42000 …line 1 column 28 near ", with consistent snapshot" 
> exec ^
: string
| start transaction with consistent snapshot, read only
< error
| 1064 42000 …line 1 column 43 near ", read only" 
> exec ^
: string
| start transaction read only, read write
< error
| 1064 42000 …line 1 column 28 near ", read write" 
> exec ^
: string
| start transaction read write, read only
< error
| 1064 42000 …line 1 column 29 near ", read only" 
> exec ^
: string
| start transaction with consistent snapshot, read write
< error
| 1064 42000 …line 1 column 43 near ", read write" 
> exec ^
: string
| start transaction with consistent snapshot, with consistent snapshot
< error
| 1064 42000 …line 1 column 43 near ", with consistent snapshot" 
> exec ^
: string
| start transaction with consistent snapshot read only
< error
| 1064 42000 …line 1 column 47 near "read only" 
> exec ^
: string
| start transaction read only read write
< error
| 1064 42000 …line 1 column 32 near "read write" 
> exec ^
: string
| start transaction read only,
< error
| 1064 42000 …line 1 column 28 near "," 
> exec ^
: string
| start transaction, read only
< error
| 1064 42000 …line 1 column 18 near ", read only" 
> exec ^
: string
| start transaction read only as of timestamp now()
< ok 0 0
> exec ^
: string
| start transaction read only as of timestamp '2015-09-21 00:07:01'
< ok 0 0
> exec ^
: string
| start transaction read only as of timestamp now() - interval 3 second
< ok 0 0
> exec ^
: string
| start transaction read only as of timestamp now() + 1
< error
| 1105 HY000 cannot set read timestamp to a future time, readTS: …, currentTS: …
> exec ^
: string
| start transaction read only as of timestamp now() + 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| start transaction read only as of timestamp now() and 1
< error
| 1292 22007 Incorrect time value: '1'
> exec ^
: string
| start transaction read only as of timestamp now() and 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| start transaction read only as of timestamp (select now())
< ok 0 0
> exec ^
: string
| start transaction read only as of timestamp a
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec ^
: string
| start transaction read only as of timestamp a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| start transaction read only as of timestamp @a
< error
| 8135 HY000 invalid as of timestamp: as of timestamp cannot be NULL
> exec ^
: string
| start transaction read only as of timestamp @a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| start transaction read only as of timestamp tidb_bounded_staleness(now(), now())
< ok 0 0
> exec ^
: string
| start transaction read only as
< error
| 1064 42000 …line 1 column 30 near "as" 
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
| start transaction read only as of time now()
< error
| 1064 42000 …line 1 column 38 near "time now()" 
> exec ^
: string
| start transaction read only as of timestamp now() for update
< error
| 1064 42000 …line 1 column 53 near "for update" 
> exec ^
: string
| start transaction read only as of timestamp now() read write
< error
| 1064 42000 …line 1 column 54 near "read write" 
> exec ^
: string
| start transaction read only as of timestamp now(), read write
< error
| 1064 42000 …line 1 column 50 near ", read write" 
> exec ^
: string
| start transaction read write as of timestamp now()
< error
| 1064 42000 …line 1 column 34 near "of timestamp now()" 
> exec ^
: string
| start transaction as of timestamp now()
< error
| 1064 42000 …line 1 column 23 near "of timestamp now()" 
> exec ^
: string
| start transaction with consistent snapshot as of timestamp now()
< error
| 1064 42000 …line 1 column 48 near "of timestamp now()" 
> exec ^
: string
| commit
< ok 0 0
> exec ^
: string
| commit;
< ok 0 0
> exec ^
: string
| commit and chain
< ok 0 0
> exec ^
: string
| commit and no chain
< ok 0 0
> exec ^
: string
| commit release
< ok 0 0
> exec ^
: string
| commit no release
< ok 0 0
> exec ^
: string
| commit and chain no release
< ok 0 0
> exec ^
: string
| commit and no chain release
< ok 0 0
> exec ^
: string
| commit and no chain no release
< ok 0 0
> exec ^
: string
| commit work
< error
| 1064 42000 …line 1 column 11 near "work" 
> exec ^
: string
| commit work and chain
< error
| 1064 42000 …line 1 column 11 near "work and chain" 
> exec ^
: string
| commit work release
< error
| 1064 42000 …line 1 column 11 near "work release" 
> exec ^
: string
| commit transaction
< error
| 1064 42000 …line 1 column 18 near "transaction" 
> exec ^
: string
| commit to savepoint s
< error
| 1064 42000 …line 1 column 9 near "to savepoint s" 
> exec ^
: string
| commit and chain release
< error
| 1064 42000 …line 1 column 24 near "release" 
> exec ^
: string
| commit and
< error
| 1064 42000 …line 1 column 10 near "" 
> exec ^
: string
| commit and no
< error
| 1064 42000 …line 1 column 13 near "" 
> exec ^
: string
| commit no
< error
| 1064 42000 …line 1 column 9 near "" 
> exec ^
: string
| commit chain
< error
| 1064 42000 …line 1 column 12 near "chain" 
> exec ^
: string
| commit no chain
< error
| 1064 42000 …line 1 column 15 near "chain" 
> exec ^
: string
| commit and no chain no
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| commit and no chain and
< error
| 1064 42000 …line 1 column 23 near "and" 
> exec ^
: string
| commit and chain no release no
< error
| 1064 42000 …line 1 column 30 near "no" 
> exec ^
: string
| commit release no
< error
| 1064 42000 …line 1 column 17 near "no" 
> exec ^
: string
| commit release release
< error
| 1064 42000 …line 1 column 22 near "release" 
> exec ^
: string
| commit no release release
< error
| 1064 42000 …line 1 column 25 near "release" 
> exec ^
: string
| commit and chain and chain
< error
| 1064 42000 …line 1 column 20 near "and chain" 
> exec ^
: string
| commit no release and chain
< error
| 1064 42000 …line 1 column 21 near "and chain" 
> exec ^
: string
| rollback
< ok 0 0
> exec ^
: string
| rollback and chain
< ok 0 0
> exec ^
: string
| rollback and no chain
< ok 0 0
> exec ^
: string
| rollback release
< ok 0 0
> exec ^
: string
| rollback no release
< ok 0 0
> exec ^
: string
| rollback and chain no release
< ok 0 0
> exec ^
: string
| rollback and no chain release
< ok 0 0
> exec ^
: string
| rollback and no chain no release
< ok 0 0
> exec ^
: string
| rollback work
< error
| 1064 42000 …line 1 column 13 near "work" 
> exec ^
: string
| rollback transaction
< error
| 1064 42000 …line 1 column 20 near "transaction" 
> exec ^
: string
| rollback and chain release
< error
| 1064 42000 …line 1 column 26 near "release" 
> exec ^
: string
| rollback and chain no
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| rollback and
< error
| 1064 42000 …line 1 column 12 near "" 
> exec ^
: string
| rollback no
< error
| 1064 42000 …line 1 column 11 near "" 
> exec ^
: string
| rollback no release and no chain
< error
| 1064 42000 …line 1 column 23 near "and no chain" 
> exec ^
: string
| rollback to savepoint s
< ok 0 0
> exec ^
: string
| rollback to s
< ok 0 0
> exec ^
: string
| rollback to `s`
< ok 0 0
> exec ^
: string
| rollback to work
< ok 0 0
> exec ^
: string
| rollback to savepoint work
< ok 0 0
> exec ^
: string
| rollback to begin
< ok 0 0
> exec ^
: string
| rollback to chain
< ok 0 0
> exec ^
: string
| rollback work to savepoint s
< error
| 1064 42000 …line 1 column 13 near "work to savepoint s" 
> exec ^
: string
| rollback work to s
< error
| 1064 42000 …line 1 column 13 near "work to s" 
> exec ^
: string
| rollback to
< error
| 1064 42000 …line 1 column 11 near "" 
> exec ^
: string
| rollback savepoint s
< error
| 1064 42000 …line 1 column 18 near "savepoint s" 
> exec ^
: string
| rollback to 's'
< error
| 1064 42000 …line 1 column 15 near "'s'" 
> exec ^
: string
| rollback to savepoint 's'
< error
| 1064 42000 …line 1 column 25 near "'s'" 
> exec ^
: string
| rollback to release
< error
| 1064 42000 …line 1 column 19 near "release" 
> exec ^
: string
| rollback to and
< error
| 1064 42000 …line 1 column 15 near "and" 
> exec ^
: string
| rollback to savepoint s s
< error
| 1064 42000 …line 1 column 25 near "s" 
> exec ^
: string
| rollback to savepoint s and chain
< error
| 1064 42000 …line 1 column 27 near "and chain" 
> exec ^
: string
| rollback to s release
< error
| 1064 42000 …line 1 column 21 near "release" 
> exec ^
: string
| rollback to savepoint
< ok 0 0
> exec ^
: string
| rollback to savepoint;
< ok 0 0
> exec ^
: string
| rollback to savepoint savepoint
< ok 0 0
> exec ^
: string
| rollback to savepoint `savepoint`
< ok 0 0
> exec ^
: string
| rollback to `savepoint`
< ok 0 0
> exec ^
: string
| rollback to savepoint savepoint savepoint
< error
| 1064 42000 …line 1 column 41 near "savepoint" 
> exec ^
: string
| rollback to savepoint and chain
< error
| 1064 42000 …line 1 column 25 near "and chain" 
> exec ^
: string
| savepoint s
< ok 0 0
> exec ^
: string
| savepoint `s`
< ok 0 0
> exec ^
: string
| savepoint work
< ok 0 0
> exec ^
: string
| savepoint savepoint
< ok 0 0
> exec ^
: string
| savepoint begin
< ok 0 0
> exec ^
: string
| savepoint commit
< ok 0 0
> exec ^
: string
| savepoint chain
< ok 0 0
> exec ^
: string
| savepoint pessimistic
< ok 0 0
> exec ^
: string
| savepoint
< error
| 1064 42000 …line 1 column 9 near "" 
> exec ^
: string
| savepoint s s
< error
| 1064 42000 …line 1 column 13 near "s" 
> exec ^
: string
| savepoint 's'
< error
| 1064 42000 …line 1 column 13 near "'s'" 
> exec ^
: string
| savepoint release
< error
| 1064 42000 …line 1 column 17 near "release" 
> exec ^
: string
| release savepoint s
< ok 0 0
> exec ^
: string
| release savepoint `s`
< ok 0 0
> exec ^
: string
| release savepoint savepoint
< ok 0 0
> exec ^
: string
| release s
< error
| 1064 42000 …line 1 column 9 near "s" 
> exec ^
: string
| release work
< error
| 1064 42000 …line 1 column 12 near "work" 
> exec ^
: string
| release
< error
| 1064 42000 …line 1 column 7 near "" 
> exec ^
: string
| release savepoint
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| release savepoint 's'
< error
| 1064 42000 …line 1 column 21 near "'s'" 
> exec ^
: string
| release savepoint s s
< error
| 1064 42000 …line 1 column 21 near "s" 
> exec ^
: string
| release savepoint savepoint savepoint
< error
| 1064 42000 …line 1 column 37 near "savepoint" 
> exec ^
: string
| create procedure p() begin commit; end
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| create procedure p() begin commit; end
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create procedure p() begin rollback; end
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| create procedure p() begin rollback; end
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create procedure p() begin rollback to s; end
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| create procedure p() begin rollback to s; end
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create procedure p() begin rollback to savepoint; end
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| create procedure p() begin rollback to savepoint; end
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create procedure p() begin commit work; end
< error
| 1064 42000 …line 1 column 38 near "work; end" 
> exec ^
: string
| create procedure p() begin rollback work; end
< error
| 1064 42000 …line 1 column 40 near "work; end" 
> exec ^
: string
| create procedure p() begin commit and chain release; end
< error
| 1064 42000 …line 1 column 51 near "release; end" 
> exec ^
: string
| create procedure p() begin work: while 1 do select 1; end while work; end
< error
| 8108 HY000 Unsupported type *resolve.NodeW
> exec ^
: string
| create procedure p() begin work: while 1 do select 1; end while work; end
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
