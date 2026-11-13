-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8

== TestInsertFormsMatchEngine default#1
> exec
| SET @omni_sweep = ?
: string
| insert into t values ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values ()
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t value ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t value ()
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| replace into t values ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| replace into t values ()
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| replace into t value ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| replace into t value ()
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert t values ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert t values ()
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| replace t value ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| replace t value ()
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values (), ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values (), ()
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t value (), ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t value (), ()
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values (), (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values (), (1)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1), ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1), ()
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t (a) values ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t (a) values ()
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t (a) value ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t (a) value ()
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t (a, b) values (), ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t (a, b) values (), ()
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t partition (p0) values ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t partition (p0) values ()
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values () on duplicate key update a = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values () on duplicate key update a = 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values( )
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values( )
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values (/* c */)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values (/* c */)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| explain insert into t values ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| explain insert into t values ()
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t () values ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t () values ()
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t () value ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t () value ()
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| replace into t () values ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| replace into t () values ()
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| replace into t () value ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| replace into t () value ()
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| replace t () values ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| replace t () values ()
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t ( ) values ( )
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t ( ) values ( )
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t (/* c */) values (/* c */)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t (/* c */) values (/* c */)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t () values (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t () values (1)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t () value (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t () value (1)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t () values (), (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t () values (), (1)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t partition (p0) () values ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t partition (p0) () values ()
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into db.t () values ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'db.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into db.t () values ()
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert ignore into t () values ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert ignore into t () values ()
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert low_priority into t () values ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert low_priority into t () values ()
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| replace low_priority into t () values ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| replace low_priority into t () values ()
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t () values () on duplicate key update a = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t () values () on duplicate key update a = 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| explain insert into t () values ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| explain insert into t () values ()
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t () select 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t () select 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t () (select 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t () (select 1)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t () with c as (select 1) select * from c
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t () with c as (select 1) select * from c
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t () table u
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'table u' at line 1
> exec
| SET @omni_sweep = ?
: string
| replace into t () table u
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'table u' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t () values row()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row()' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t () values row(1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t () set a = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'set a = 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (a) set a = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'set a = 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| replace into t (a) set a = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'set a = 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| replace into t () set a = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'set a = 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (a)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t () on duplicate key update a = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'on duplicate key update a = 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (()) values ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ')) values ()' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (,) values ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ') values ()' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (a,) values (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ') values (1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (,a) values (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'a) values (1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t value
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '()' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values () ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '()' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values default
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'default' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (),
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values ,()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '()' at line 1
> exec
| SET @omni_sweep = ?
: string
| replace into t values () on duplicate key update a = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'on duplicate key update a = 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (@a) values (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '@a) values (1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (a, @b) values (1, 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '@b) values (1, 2)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (@a) select 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '@a) select 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| replace into t (@a) values (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '@a) values (1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (t.*) values (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t (t.*) values (1)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t (`@a`) values (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t (`@a`) values (1)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t (value) values (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t (value) values (1)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t (values) values (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ') values (1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (table) values (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'table) values (1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (select) values (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ') values (1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (with) values (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ') values (1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values row(1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values row(1), row(2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1), row(2)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values row()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row()' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values row(), row()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(), row()' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (a) values row(1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| replace into t values row(1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| replace into t values row()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row()' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t value row(1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t value row()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row()' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t value row(1), row(2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1), row(2)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (a) value row(1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t () value row()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row()' at line 1
> exec
| SET @omni_sweep = ?
: string
| replace into t value row(1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values row(1), (2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1), (2)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1), row(2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(2)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t value (1), row(2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(2)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t value row(1), (2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1), (2)' at line 1
> exec
| SET @omni_sweep = ?
: string
| replace into t values (1), row(2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(2)' at line 1
> exec
| SET @omni_sweep = ?
: string
| replace into t values row(1), (2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1), (2)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values row 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values row
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values row(default)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(default)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values row(default, 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(default, 1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values row((1))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row((1))' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values row(1) on duplicate key update a = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1) on duplicate key update a = 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values row(1) as new
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1) as new' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values row(1) order by 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1) order by 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values row(1) limit 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1) limit 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values row(1), row(2) order by 1 limit 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1), row(2) order by 1 limit 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values row(1) union values row(2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1) union values row(2)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values row(1) union all select 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1) union all select 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values row(1) union values row(2) on duplicate key update a = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1) union values row(2) on duplicate key update a = 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t select 1 union values row(1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (select 1) union values row(1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (values row(1))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1))' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t ((values row(1)))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1)))' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (a) (values row(1))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1))' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (values row(1)) on duplicate key update a = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1)) on duplicate key update a = 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (values row(1)) union (values row(2))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1)) union (values row(2))' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (values row(1)) order by 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1)) order by 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (a) values row(1) order by 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1) order by 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t () values row(1) order by 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1) order by 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| replace into t values row(1) order by 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1) order by 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| replace into t values row(1) union values row(2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1) union values row(2)' at line 1
> exec
| SET @omni_sweep = ?
: string
| replace into t (values row(1))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1))' at line 1
> exec
| SET @omni_sweep = ?
: string
| replace into t () (values row(1))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1))' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) union values row(2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(2)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values row(1) for update
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1) for update' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values row(1) order by 1 for update
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1) order by 1 for update' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) limit 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) limit 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) order by 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) order by 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values () order by 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values () order by 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t () values () limit 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t () values () limit 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t () values (1) union values (2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t () values (1) union values (2)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t () value () union values ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t () value () union values ()
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t value () union value ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'value ()' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (values ())
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t (values ())
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t () (values ())
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t () (values ())
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t (select 1) on duplicate key update a = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t (select 1) on duplicate key update a = 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t (a) (select 1) on duplicate key update a = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t (a) (select 1) on duplicate key update a = 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t ((select 1)) on duplicate key update a = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t ((select 1)) on duplicate key update a = 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t (select 1) as new
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'as new' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (table u)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'table u)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (table u) on duplicate key update a = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'table u) on duplicate key update a = 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t table u
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'table u' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t table u order by a
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'table u order by a' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (a) table u
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'table u' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t table u as new
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'table u as new' at line 1
> exec
| SET @omni_sweep = ?
: string
| replace into t table u
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'table u' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (values row(1) into @x)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1) into @x)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (table u into @x)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'table u into @x)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t ((values row(1) into @x))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1) into @x))' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t ((table u into @x))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'table u into @x))' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (a) (values row(1) into @x)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1) into @x)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (a) (table u into @x)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'table u into @x)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t () (values row(1) into @x)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1) into @x)' at line 1
> exec
| SET @omni_sweep = ?
: string
| replace into t (values row(1) into @x)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1) into @x)' at line 1
> exec
| SET @omni_sweep = ?
: string
| replace into t (table u into @x)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'table u into @x)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (select 1 into @x)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'into @x)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t ((select 1 into @x))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'into @x))' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (a) (select 1 into @x)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'into @x)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (table u) into @x
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'table u) into @x' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (select 1) union (values row(1) into @x)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1) into @x)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (values row(1) into @x) union (select 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1) into @x) union (select 1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (select 1) union (table u into @x)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'table u into @x)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (table u into @x) union (select 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'table u into @x) union (select 1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (select 1) union (select 2 into @x)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'into @x)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (select 1 into @x) union (select 2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'into @x) union (select 2)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t select 1 union (values row(1) into @x)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1) into @x)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t select 1 union (table u into @x)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'table u into @x)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values row(1) union (values row(2) into @x)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1) union (values row(2) into @x)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t select 1 union (select 2 into @x)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'into @x)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t ((values row(1)) into @x)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1)) into @x)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t ((table u) into @x)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'table u) into @x)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (values row(1)) union values row(2) into @x
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1)) union values row(2) into @x' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (select 1) union values row(1) into @x
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1) into @x' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (select 1) union table u into @x
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'table u into @x' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (select 1) order by 1 into @x
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'into @x' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (values row(1)) order by 1 into @x
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1)) order by 1 into @x' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (table u) order by 1 into @x
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'table u) order by 1 into @x' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (select 1) limit 1 into @x
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'into @x' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values row(1) into @x
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1) into @x' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t table u into @x
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'table u into @x' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t table u limit 1 into @x
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'table u limit 1 into @x' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t select 1 into @x
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'into @x' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t select 1 union select 2 into @x
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'into @x' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t select 1 union table u into @x
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'table u into @x' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t select 1 union values row(1) into @x
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1) into @x' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values row(1) union values row(2) into @x
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1) union values row(2) into @x' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t table u union table v into @x
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'table u union table v into @x' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values row(1) union select 1 into @x
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1) union select 1 into @x' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (select 1) into @x
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'into @x' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (values row(1)) into @x
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1)) into @x' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (default)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values (default)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values (default, 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values (default, 1)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values ((1))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values ((1))
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values (row(1,2))
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values (row(1,2))
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1), (2,3)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1), (2,3)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values (*)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '*)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (t.*)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '*)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1, *)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '*)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1, t.*)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '*)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (t.*, 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '*, 1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (db.t.*)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '*)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t value (*)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '*)' at line 1
> exec
| SET @omni_sweep = ?
: string
| replace into t values (t.*)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '*)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values row(*)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(*)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values row(t.*)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(t.*)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select * from (values row(*)) v
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(*)) v' at line 1
> exec
| SET @omni_sweep = ?
: string
| select * from (values row(t.*)) v
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(t.*)) v' at line 1
> exec
| SET @omni_sweep = ?
: string
| values row(*)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(*)' at line 1
> exec
| SET @omni_sweep = ?
: string
| values row(t.*)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(t.*)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (,)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ')' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (,1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t value (,1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| replace into t values (,1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t () values (,1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (a) values (,1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1), (,2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '2)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (,1), (,2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '1), (,2)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (,,1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ',1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (,1,)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '1,)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (,1,2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '1,2)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (,default)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'default)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (,t.*)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 't.*)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1,)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ')' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1,,2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '2)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1),(,)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ')' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values row(,)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(,)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values row(,1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(,1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values row(,1), row(,2)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(,1), row(,2)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values row(1,)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1,)' at line 1
> exec
| SET @omni_sweep = ?
: string
| values row(,1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(,1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| select * from (values row(,1)) v
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(,1)) v' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert high_priority t values ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert high_priority t values ()
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert delayed into t values ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert delayed into t values ()
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert delayed t set a = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert delayed t set a = 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert delayed into t select 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert delayed into t select 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| replace delayed into t values ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| replace delayed into t values ()
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| replace delayed t values (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| replace delayed t values (1)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| replace high_priority into t values ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'high_priority into t values ()' at line 1
> exec
| SET @omni_sweep = ?
: string
| replace high_priority into t values (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'high_priority into t values (1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| replace ignore into t values (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'ignore into t values (1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert low_priority ignore into t values (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert low_priority ignore into t values (1)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert ignore low_priority into t values (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'low_priority into t values (1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert high_priority delayed into t values (1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'delayed into t values (1)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values () as new
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'as new' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t () values () as new
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'as new' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t () values () as new (a, b)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'as new (a, b)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values () as new on duplicate key update a = new.a
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'as new on duplicate key update a = new.a' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (), () as new
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'as new' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) as new
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'as new' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) as new (a)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'as new (a)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) as new ()
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'as new ()' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) as new (a,)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'as new (a,)' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) new
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'new' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t set a = 1 as new
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'as new' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t set a = 1 as new on duplicate key update a = new.a
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'as new on duplicate key update a = new.a' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t set a = 1 as new (a)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'as new (a)' at line 1
> exec
| SET @omni_sweep = ?
: string
| replace into t values (1) as new
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'as new' at line 1
> exec
| SET @omni_sweep = ?
: string
| replace into t set a = 1 as new
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'as new' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t select 1 as new
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t select 1 as new
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t (a) select 1 as new from u as new
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t (a) select 1 as new from u as new
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t set a := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t set a := 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t set a:=1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t set a:=1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t set a := default
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t set a := default
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t set a := 1, b = 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t set a := 1, b = 2
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t set a = 1, b := 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t set a = 1, b := 2
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t set t.a := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t set t.a := 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t set `@a` := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t set `@a` := 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| replace into t set a := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| replace into t set a := 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| replace t set a := default
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| replace t set a := default
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t set a : = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ': = 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t set a :=
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t set a := 1 on duplicate key update a := 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t set a := 1 on duplicate key update a := 2
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t set a = 1 on duplicate key update a := 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t set a = 1 on duplicate key update a := 2
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t set a = @b := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t set a = @b := 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) on duplicate key update a := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) on duplicate key update a := 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values(1) on duplicate key update a:=1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values(1) on duplicate key update a:=1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) on duplicate key update a := values(a)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) on duplicate key update a := values(a)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) on duplicate key update a := default
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) on duplicate key update a := default
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) on duplicate key update a = 1, b := 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) on duplicate key update a = 1, b := 2
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) on duplicate key update t.a := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) on duplicate key update t.a := 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) on duplicate key update `@a` := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) on duplicate key update `@a` := 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) on duplicate key update a : = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ': = 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t (a) select 1 on duplicate key update a := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t (a) select 1 on duplicate key update a := 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values row(1) on duplicate key update a := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'row(1) on duplicate key update a := 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t () values () on duplicate key update a := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t () values () on duplicate key update a := 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1), (2) as new (a) on duplicate key update a := new.a
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'as new (a) on duplicate key update a := new.a' at line 1
> exec
| SET @omni_sweep = ?
: string
| INSERT INTO t (a,b,c) VALUES (1,2,3),(4,5,6) ON DUPLICATE KEY UPDATE c:=VALUES(a)+VALUES(b)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| INSERT INTO t (a,b,c) VALUES (1,2,3),(4,5,6) ON DUPLICATE KEY UPDATE c:=VALUES(a)+VALUES(b)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| INSERT IGNORE INTO t (a,b,c) VALUES (1,2,3),(4,5,6) ON DUPLICATE KEY UPDATE c:=VALUES(a)+VALUES(b)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| INSERT IGNORE INTO t (a,b,c) VALUES (1,2,3),(4,5,6) ON DUPLICATE KEY UPDATE c:=VALUES(a)+VALUES(b)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| update t set a := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| update t set a := 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| update t set a:=1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| update t set a:=1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| update t set a := default
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| update t set a := default
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| update t set a := 1, b = 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| update t set a := 1, b = 2
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| update t set t.a := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| update t set t.a := 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| update t, u set t.a := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| update t, u set t.a := 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| update t set `@a` := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| update t set `@a` := 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| update t set a := (select 1)
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| update t set a := (select 1)
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| update t set a := 1 where b = 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| update t set a := 1 where b = 2
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| update t set a := 1 order by b limit 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| update t set a := 1 order by b limit 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| update t set a : = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ': = 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| update t set a :=
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '' at line 1
> exec
| SET @omni_sweep = ?
: string
| update t set a := b := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| update t set a = b := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| update t set a := @b := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| update t set a := @b := 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| update t set a = @b := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| update t set a = @b := 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| update t set a = 1, b := @c := 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| update t set a = 1, b := @c := 2
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| explain update t set a := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| explain update t set a := 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| load data infile 'f' into table t set a = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| load data infile 'f' into table t set a = 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| load data infile 'f' into table t set a = default
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| load data infile 'f' into table t set a = default
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| load data infile 'f' into table t set a := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| load data infile 'f' into table t set a := 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| load data infile 'f' into table t set a:=1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| load data infile 'f' into table t set a:=1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| load data infile 'f' into table t set a := default
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| load data infile 'f' into table t set a := default
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| load data infile 'f' into table t set a = 1, b := 2
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| load data infile 'f' into table t set a = 1, b := 2
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| load data infile 'f' into table t set t.a := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| load data infile 'f' into table t set t.a := 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| load data local infile 'f' into table t set a := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| load data local infile 'f' into table t set a := 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| load xml infile 'f' into table t set a := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| load xml infile 'f' into table t set a := 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| load data infile 'f' into table t set a : = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ': = 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| load data infile 'f' into table t (a, @b) set a := @b
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| load data infile 'f' into table t (a, @b) set a := @b
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| update t set @a = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '@a = 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| update t set @a := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '@a := 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| update t set @@autocommit = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '@@autocommit = 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| update t set @@autocommit := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '@@autocommit := 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| update t set t.* = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '* = 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| update t set t.* := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '* := 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| update t set `@a` = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| update t set `@a` = 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t set @a = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '@a = 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t set @a := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '@a := 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t set @@autocommit = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '@@autocommit = 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t set @@autocommit := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '@@autocommit := 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t set t.* = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '* = 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t set `@a` = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t set `@a` = 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) on duplicate key update @a = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '@a = 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) on duplicate key update @a := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '@a := 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) on duplicate key update t.* = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '* = 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| load data infile 'f' into table t set @a = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '@a = 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| load data infile 'f' into table t set @a := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '@a := 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| load data infile 'f' into table t set t.* = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near '* = 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| load data infile 'f' into table t set `@a` = 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| load data infile 'f' into table t set `@a` = 1
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| select a := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near ':= 1' at line 1
> exec
| SET @omni_sweep = ?
: string
| select @a := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| set @a := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| set autocommit := 1
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< ok 0 0
> exec
| DEALLOCATE PREPARE omni_sweep_stmt
< ok 0 0
> exec
| SET @omni_sweep = ?
: string
| insert into t values () returning a
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values () returning a
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t () values () returning a
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t () values () returning a
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t () select 1 returning a
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t () select 1 returning a
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values (), () returning *
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values (), () returning *
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| replace into t () values () returning a
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| replace into t () values () returning a
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t set a := 1 returning a
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t set a := 1 returning a
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) on duplicate key update a := 1 returning a
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec
| SET @omni_sweep = ?
: string
| insert into t values (1) on duplicate key update a := 1 returning a
| ; SELECT FROM
< ok 0 0
> exec
| PREPARE omni_sweep_stmt FROM @omni_sweep
< error
| 1064 42000 …near 'SELECT FROM' at line 2
