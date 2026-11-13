-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestInsertFormsMatchEngine default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| insert into t values ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t value ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t value ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| replace into t values ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| replace into t values ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| replace into t value ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| replace into t value ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert t values ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert t values ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| replace t value ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| replace t value ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values (), ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values (), ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t value (), ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t value (), ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values (), (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values (), (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values (1), ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values (1), ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t (a) values ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t (a) values ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t (a) value ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t (a) value ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t (a, b) values (), ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t (a, b) values (), ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t partition (p0) values ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t partition (p0) values ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values () on duplicate key update a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values () on duplicate key update a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values( )
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values( )
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values (/* c */)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values (/* c */)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| explain insert into t values ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| explain insert into t values ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t () values ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t () values ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t () value ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t () value ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| replace into t () values ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| replace into t () values ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| replace into t () value ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| replace into t () value ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| replace t () values ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| replace t () values ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t ( ) values ( )
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t ( ) values ( )
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t (/* c */) values (/* c */)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t (/* c */) values (/* c */)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t () values (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t () values (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t () value (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t () value (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t () values (), (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t () values (), (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t partition (p0) () values ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t partition (p0) () values ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into db.t () values ()
< error
| 1146 42S02 Table 'db.t' doesn't exist
> exec ^
: string
| insert into db.t () values ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert ignore into t () values ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert ignore into t () values ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert low_priority into t () values ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert low_priority into t () values ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| replace low_priority into t () values ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| replace low_priority into t () values ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t () values () on duplicate key update a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t () values () on duplicate key update a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| explain insert into t () values ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| explain insert into t () values ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t () select 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t () select 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t () (select 1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t () (select 1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t () with c as (select 1) select * from c
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t () with c as (select 1) select * from c
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t () table u
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t () table u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| replace into t () table u
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| replace into t () table u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t () values row()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t () values row()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t () values row(1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t () values row(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t () set a = 1
< error
| 1064 42000 …line 1 column 20 near "set a = 1" 
> exec ^
: string
| insert into t (a) set a = 1
< error
| 1064 42000 …line 1 column 21 near "set a = 1" 
> exec ^
: string
| replace into t (a) set a = 1
< error
| 1064 42000 …line 1 column 22 near "set a = 1" 
> exec ^
: string
| replace into t () set a = 1
< error
| 1064 42000 …line 1 column 21 near "set a = 1" 
> exec ^
: string
| insert into t ()
< error
| 1064 42000 …line 1 column 16 near "" 
> exec ^
: string
| insert into t (a)
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| insert into t () on duplicate key update a = 1
< error
| 1064 42000 …line 1 column 19 near "on duplicate key update a = 1" 
> exec ^
: string
| insert into t (()) values ()
< error
| 1064 42000 …line 1 column 17 near ")) values ()" 
> exec ^
: string
| insert into t (,) values ()
< error
| 1064 42000 …line 1 column 16 near ",) values ()" 
> exec ^
: string
| insert into t (a,) values (1)
< error
| 1064 42000 …line 1 column 18 near ") values (1)" 
> exec ^
: string
| insert into t (,a) values (1)
< error
| 1064 42000 …line 1 column 16 near ",a) values (1)" 
> exec ^
: string
| insert into t values
< error
| 1064 42000 …line 1 column 20 near "" 
> exec ^
: string
| insert into t value
< error
| 1064 42000 …line 1 column 19 near "" 
> exec ^
: string
| insert into t values (1) ()
< error
| 1064 42000 …line 1 column 26 near "()" 
> exec ^
: string
| insert into t values () ()
< error
| 1064 42000 …line 1 column 25 near "()" 
> exec ^
: string
| insert into t values default
< error
| 1064 42000 …line 1 column 28 near "default" 
> exec ^
: string
| insert into t values (),
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| insert into t values ,()
< error
| 1064 42000 …line 1 column 22 near ",()" 
> exec ^
: string
| replace into t values () on duplicate key update a = 1
< error
| 1064 42000 …line 1 column 27 near "on duplicate key update a = 1" 
> exec ^
: string
| insert into t (@a) values (1)
< error
| 1064 42000 …line 1 column 17 near "@a) values (1)" 
> exec ^
: string
| insert into t (a, @b) values (1, 2)
< error
| 1064 42000 …line 1 column 20 near "@b) values (1, 2)" 
> exec ^
: string
| insert into t (@a) select 1
< error
| 1064 42000 …line 1 column 17 near "@a) select 1" 
> exec ^
: string
| replace into t (@a) values (1)
< error
| 1064 42000 …line 1 column 18 near "@a) values (1)" 
> exec ^
: string
| insert into t (t.*) values (1)
< error
| 1064 42000 …line 1 column 18 near "*) values (1)" 
> exec ^
: string
| insert into t (`@a`) values (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t (`@a`) values (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t (value) values (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t (value) values (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t (values) values (1)
< error
| 1064 42000 …line 1 column 22 near ") values (1)" 
> exec ^
: string
| insert into t (table) values (1)
< error
| 1064 42000 …line 1 column 21 near ") values (1)" 
> exec ^
: string
| insert into t (select) values (1)
< error
| 1064 42000 …line 1 column 22 near ") values (1)" 
> exec ^
: string
| insert into t (with) values (1)
< error
| 1064 42000 …line 1 column 20 near ") values (1)" 
> exec ^
: string
| insert into t values row(1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values row(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values row(1), row(2)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values row(1), row(2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values row()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values row()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values row(), row()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values row(), row()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t (a) values row(1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t (a) values row(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| replace into t values row(1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| replace into t values row(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| replace into t values row()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| replace into t values row()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t value row(1)
< error
| 1064 42000 …line 1 column 23 near "row(1)" 
> exec ^
: string
| insert into t value row()
< error
| 1064 42000 …line 1 column 23 near "row()" 
> exec ^
: string
| insert into t value row(1), row(2)
< error
| 1064 42000 …line 1 column 23 near "row(1), row(2)" 
> exec ^
: string
| insert into t (a) value row(1)
< error
| 1064 42000 …line 1 column 27 near "row(1)" 
> exec ^
: string
| insert into t () value row()
< error
| 1064 42000 …line 1 column 26 near "row()" 
> exec ^
: string
| replace into t value row(1)
< error
| 1064 42000 …line 1 column 24 near "row(1)" 
> exec ^
: string
| insert into t values row(1), (2)
< error
| 1064 42000 …line 1 column 30 near "(2)" 
> exec ^
: string
| insert into t values (1), row(2)
< error
| 1064 42000 …line 1 column 29 near "row(2)" 
> exec ^
: string
| insert into t value (1), row(2)
< error
| 1064 42000 …line 1 column 28 near "row(2)" 
> exec ^
: string
| insert into t value row(1), (2)
< error
| 1064 42000 …line 1 column 23 near "row(1), (2)" 
> exec ^
: string
| replace into t values (1), row(2)
< error
| 1064 42000 …line 1 column 30 near "row(2)" 
> exec ^
: string
| replace into t values row(1), (2)
< error
| 1064 42000 …line 1 column 31 near "(2)" 
> exec ^
: string
| insert into t values row 1
< error
| 1064 42000 …line 1 column 26 near "1" 
> exec ^
: string
| insert into t values row
< error
| 1064 42000 …line 1 column 24 near "" 
> exec ^
: string
| insert into t values row(default)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values row(default)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values row(default, 1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values row(default, 1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values row((1))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values row((1))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values row(1) on duplicate key update a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values row(1) on duplicate key update a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values row(1) as new
< error
| 1064 42000 …line 1 column 30 near "as new" 
> exec ^
: string
| insert into t values row(1) order by 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values row(1) order by 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values row(1) limit 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values row(1) limit 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values row(1), row(2) order by 1 limit 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values row(1), row(2) order by 1 limit 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values row(1) union values row(2)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values row(1) union values row(2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values row(1) union all select 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values row(1) union all select 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values row(1) union values row(2) on duplicate key update a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values row(1) union values row(2) on duplicate key update a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t select 1 union values row(1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t select 1 union values row(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t (select 1) union values row(1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t (select 1) union values row(1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t (values row(1))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t (values row(1))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t ((values row(1)))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t ((values row(1)))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t (a) (values row(1))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t (a) (values row(1))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t (values row(1)) on duplicate key update a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t (values row(1)) on duplicate key update a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t (values row(1)) union (values row(2))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t (values row(1)) union (values row(2))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t (values row(1)) order by 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t (values row(1)) order by 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t (a) values row(1) order by 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t (a) values row(1) order by 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t () values row(1) order by 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t () values row(1) order by 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| replace into t values row(1) order by 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| replace into t values row(1) order by 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| replace into t values row(1) union values row(2)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| replace into t values row(1) union values row(2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| replace into t (values row(1))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| replace into t (values row(1))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| replace into t () (values row(1))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| replace into t () (values row(1))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values (1) union values row(2)
< error
| 1064 42000 …line 1 column 30 near "union values row(2)" 
> exec ^
: string
| insert into t values row(1) for update
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values row(1) for update
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values row(1) order by 1 for update
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values row(1) order by 1 for update
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values (1) limit 1
< error
| 1064 42000 …line 1 column 30 near "limit 1" 
> exec ^
: string
| insert into t values (1) order by 1
< error
| 1064 42000 …line 1 column 30 near "order by 1" 
> exec ^
: string
| insert into t values () order by 1
< error
| 1064 42000 …line 1 column 29 near "order by 1" 
> exec ^
: string
| insert into t () values () limit 1
< error
| 1064 42000 …line 1 column 32 near "limit 1" 
> exec ^
: string
| insert into t () values (1) union values (2)
< error
| 1064 42000 …line 1 column 33 near "union values (2)" 
> exec ^
: string
| insert into t () value () union values ()
< error
| 1064 42000 …line 1 column 31 near "union values ()" 
> exec ^
: string
| insert into t value () union value ()
< error
| 1064 42000 …line 1 column 28 near "union value ()" 
> exec ^
: string
| insert into t (values ())
< error
| 1064 42000 …line 1 column 23 near "())" 
> exec ^
: string
| insert into t () (values ())
< error
| 1064 42000 …line 1 column 26 near "())" 
> exec ^
: string
| insert into t (select 1) on duplicate key update a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t (select 1) on duplicate key update a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t (a) (select 1) on duplicate key update a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t (a) (select 1) on duplicate key update a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t ((select 1)) on duplicate key update a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t ((select 1)) on duplicate key update a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t (select 1) as new
< error
| 1064 42000 …line 1 column 27 near "as new" 
> exec ^
: string
| insert into t (table u)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t (table u)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t (table u) on duplicate key update a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t (table u) on duplicate key update a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t table u
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t table u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t table u order by a
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t table u order by a
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t (a) table u
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t (a) table u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t table u as new
< error
| 1064 42000 …line 1 column 24 near "as new" 
> exec ^
: string
| replace into t table u
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| replace into t table u
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t (values row(1) into @x)
< error
| 1064 42000 …line 1 column 36 near "@x)" 
> exec ^
: string
| insert into t (table u into @x)
< error
| 1064 42000 …line 1 column 30 near "@x)" 
> exec ^
: string
| insert into t ((values row(1) into @x))
< error
| 1064 42000 …line 1 column 37 near "@x))" 
> exec ^
: string
| insert into t ((table u into @x))
< error
| 1064 42000 …line 1 column 31 near "@x))" 
> exec ^
: string
| insert into t (a) (values row(1) into @x)
< error
| 1064 42000 …line 1 column 40 near "@x)" 
> exec ^
: string
| insert into t (a) (table u into @x)
< error
| 1064 42000 …line 1 column 34 near "@x)" 
> exec ^
: string
| insert into t () (values row(1) into @x)
< error
| 1064 42000 …line 1 column 39 near "@x)" 
> exec ^
: string
| replace into t (values row(1) into @x)
< error
| 1064 42000 …line 1 column 37 near "@x)" 
> exec ^
: string
| replace into t (table u into @x)
< error
| 1064 42000 …line 1 column 31 near "@x)" 
> exec ^
: string
| insert into t (select 1 into @x)
< error
| 1064 42000 …line 1 column 31 near "@x)" 
> exec ^
: string
| insert into t ((select 1 into @x))
< error
| 1064 42000 …line 1 column 32 near "@x))" 
> exec ^
: string
| insert into t (a) (select 1 into @x)
< error
| 1064 42000 …line 1 column 35 near "@x)" 
> exec ^
: string
| insert into t (table u) into @x
< error
| 1064 42000 …line 1 column 28 near "into @x" 
> exec ^
: string
| insert into t (select 1) union (values row(1) into @x)
< error
| 1064 42000 …line 1 column 53 near "@x)" 
> exec ^
: string
| insert into t (values row(1) into @x) union (select 1)
< error
| 1064 42000 …line 1 column 36 near "@x) union (select 1)" 
> exec ^
: string
| insert into t (select 1) union (table u into @x)
< error
| 1064 42000 …line 1 column 47 near "@x)" 
> exec ^
: string
| insert into t (table u into @x) union (select 1)
< error
| 1064 42000 …line 1 column 30 near "@x) union (select 1)" 
> exec ^
: string
| insert into t (select 1) union (select 2 into @x)
< error
| 1064 42000 …line 1 column 48 near "@x)" 
> exec ^
: string
| insert into t (select 1 into @x) union (select 2)
< error
| 1064 42000 …line 1 column 31 near "@x) union (select 2)" 
> exec ^
: string
| insert into t select 1 union (values row(1) into @x)
< error
| 1064 42000 …line 1 column 51 near "@x)" 
> exec ^
: string
| insert into t select 1 union (table u into @x)
< error
| 1064 42000 …line 1 column 45 near "@x)" 
> exec ^
: string
| insert into t values row(1) union (values row(2) into @x)
< error
| 1064 42000 …line 1 column 56 near "@x)" 
> exec ^
: string
| insert into t select 1 union (select 2 into @x)
< error
| 1064 42000 …line 1 column 46 near "@x)" 
> exec ^
: string
| insert into t ((values row(1)) into @x)
< error
| 1064 42000 …line 1 column 35 near "into @x)" 
> exec ^
: string
| insert into t ((table u) into @x)
< error
| 1064 42000 …line 1 column 29 near "into @x)" 
> exec ^
: string
| insert into t (values row(1)) union values row(2) into @x
< error
| 1064 42000 …line 1 column 57 near "@x" 
> exec ^
: string
| insert into t (select 1) union values row(1) into @x
< error
| 1064 42000 …line 1 column 52 near "@x" 
> exec ^
: string
| insert into t (select 1) union table u into @x
< error
| 1064 42000 …line 1 column 46 near "@x" 
> exec ^
: string
| insert into t (select 1) order by 1 into @x
< error
| 1064 42000 …line 1 column 40 near "into @x" 
> exec ^
: string
| insert into t (values row(1)) order by 1 into @x
< error
| 1064 42000 …line 1 column 45 near "into @x" 
> exec ^
: string
| insert into t (table u) order by 1 into @x
< error
| 1064 42000 …line 1 column 39 near "into @x" 
> exec ^
: string
| insert into t (select 1) limit 1 into @x
< error
| 1064 42000 …line 1 column 37 near "into @x" 
> exec ^
: string
| insert into t values row(1) into @x
< error
| 1064 42000 …line 1 column 35 near "@x" 
> exec ^
: string
| insert into t table u into @x
< error
| 1064 42000 …line 1 column 29 near "@x" 
> exec ^
: string
| insert into t table u limit 1 into @x
< error
| 1064 42000 …line 1 column 37 near "@x" 
> exec ^
: string
| insert into t select 1 into @x
< error
| 1064 42000 …line 1 column 30 near "@x" 
> exec ^
: string
| insert into t select 1 union select 2 into @x
< error
| 1064 42000 …line 1 column 45 near "@x" 
> exec ^
: string
| insert into t select 1 union table u into @x
< error
| 1064 42000 …line 1 column 44 near "@x" 
> exec ^
: string
| insert into t select 1 union values row(1) into @x
< error
| 1064 42000 …line 1 column 50 near "@x" 
> exec ^
: string
| insert into t values row(1) union values row(2) into @x
< error
| 1064 42000 …line 1 column 55 near "@x" 
> exec ^
: string
| insert into t table u union table v into @x
< error
| 1064 42000 …line 1 column 43 near "@x" 
> exec ^
: string
| insert into t values row(1) union select 1 into @x
< error
| 1064 42000 …line 1 column 50 near "@x" 
> exec ^
: string
| insert into t (select 1) into @x
< error
| 1064 42000 …line 1 column 29 near "into @x" 
> exec ^
: string
| insert into t (values row(1)) into @x
< error
| 1064 42000 …line 1 column 34 near "into @x" 
> exec ^
: string
| insert into t values (default)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values (default)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values (default, 1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values (default, 1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values ((1))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values ((1))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values (row(1,2))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values (row(1,2))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values (1), (2,3)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values (1), (2,3)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values (*)
< error
| 1064 42000 …line 1 column 23 near "*)" 
> exec ^
: string
| insert into t values (t.*)
< error
| 1064 42000 …line 1 column 25 near "*)" 
> exec ^
: string
| insert into t values (1, *)
< error
| 1064 42000 …line 1 column 26 near "*)" 
> exec ^
: string
| insert into t values (1, t.*)
< error
| 1064 42000 …line 1 column 28 near "*)" 
> exec ^
: string
| insert into t values (t.*, 1)
< error
| 1064 42000 …line 1 column 25 near "*, 1)" 
> exec ^
: string
| insert into t values (db.t.*)
< error
| 1064 42000 …line 1 column 28 near "*)" 
> exec ^
: string
| insert into t value (*)
< error
| 1064 42000 …line 1 column 22 near "*)" 
> exec ^
: string
| replace into t values (t.*)
< error
| 1064 42000 …line 1 column 26 near "*)" 
> exec ^
: string
| insert into t values row(*)
< error
| 1064 42000 …line 1 column 26 near "*)" 
> exec ^
: string
| insert into t values row(t.*)
< error
| 1064 42000 …line 1 column 28 near "*)" 
> exec ^
: string
| select * from (values row(*)) v
< error
| 1064 42000 …line 1 column 27 near "*)) v" 
> exec ^
: string
| select * from (values row(t.*)) v
< error
| 1064 42000 …line 1 column 29 near "*)) v" 
> exec ^
: string
| values row(*)
< error
| 1064 42000 …line 1 column 12 near "*)" 
> exec ^
: string
| values row(t.*)
< error
| 1064 42000 …line 1 column 14 near "*)" 
> exec ^
: string
| insert into t values (,)
< error
| 1064 42000 …line 1 column 23 near ",)" 
> exec ^
: string
| insert into t values (,1)
< error
| 1064 42000 …line 1 column 23 near ",1)" 
> exec ^
: string
| insert into t value (,1)
< error
| 1064 42000 …line 1 column 22 near ",1)" 
> exec ^
: string
| replace into t values (,1)
< error
| 1064 42000 …line 1 column 24 near ",1)" 
> exec ^
: string
| insert into t () values (,1)
< error
| 1064 42000 …line 1 column 26 near ",1)" 
> exec ^
: string
| insert into t (a) values (,1)
< error
| 1064 42000 …line 1 column 27 near ",1)" 
> exec ^
: string
| insert into t values (1), (,2)
< error
| 1064 42000 …line 1 column 28 near ",2)" 
> exec ^
: string
| insert into t values (,1), (,2)
< error
| 1064 42000 …line 1 column 23 near ",1), (,2)" 
> exec ^
: string
| insert into t values (,,1)
< error
| 1064 42000 …line 1 column 23 near ",,1)" 
> exec ^
: string
| insert into t values (,1,)
< error
| 1064 42000 …line 1 column 23 near ",1,)" 
> exec ^
: string
| insert into t values (,1,2)
< error
| 1064 42000 …line 1 column 23 near ",1,2)" 
> exec ^
: string
| insert into t values (,default)
< error
| 1064 42000 …line 1 column 23 near ",default)" 
> exec ^
: string
| insert into t values (,t.*)
< error
| 1064 42000 …line 1 column 23 near ",t.*)" 
> exec ^
: string
| insert into t values (1,)
< error
| 1064 42000 …line 1 column 25 near ")" 
> exec ^
: string
| insert into t values (1,,2)
< error
| 1064 42000 …line 1 column 25 near ",2)" 
> exec ^
: string
| insert into t values (1),(,)
< error
| 1064 42000 …line 1 column 27 near ",)" 
> exec ^
: string
| insert into t values row(,)
< error
| 1064 42000 …line 1 column 26 near ",)" 
> exec ^
: string
| insert into t values row(,1)
< error
| 1064 42000 …line 1 column 26 near ",1)" 
> exec ^
: string
| insert into t values row(,1), row(,2)
< error
| 1064 42000 …line 1 column 26 near ",1), row(,2)" 
> exec ^
: string
| insert into t values row(1,)
< error
| 1064 42000 …line 1 column 28 near ")" 
> exec ^
: string
| values row(,1)
< error
| 1064 42000 …line 1 column 12 near ",1)" 
> exec ^
: string
| select * from (values row(,1)) v
< error
| 1064 42000 …line 1 column 27 near ",1)) v" 
> exec ^
: string
| insert high_priority t values ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert high_priority t values ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert delayed into t values ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert delayed into t values ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert delayed t set a = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert delayed t set a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert delayed into t select 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert delayed into t select 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| replace delayed into t values ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| replace delayed into t values ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| replace delayed t values (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| replace delayed t values (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| replace high_priority into t values ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| replace high_priority into t values ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| replace high_priority into t values (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| replace high_priority into t values (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| replace ignore into t values (1)
< error
| 1064 42000 …line 1 column 14 near "ignore into t values (1)" 
> exec ^
: string
| insert low_priority ignore into t values (1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert low_priority ignore into t values (1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert ignore low_priority into t values (1)
< error
| 1064 42000 …line 1 column 26 near "low_priority into t values (1)" 
> exec ^
: string
| insert high_priority delayed into t values (1)
< error
| 1064 42000 …line 1 column 28 near "delayed into t values (1)" 
> exec ^
: string
| insert into t values () as new
< error
| 1064 42000 …line 1 column 26 near "as new" 
> exec ^
: string
| insert into t () values () as new
< error
| 1064 42000 …line 1 column 29 near "as new" 
> exec ^
: string
| insert into t () values () as new (a, b)
< error
| 1064 42000 …line 1 column 29 near "as new (a, b)" 
> exec ^
: string
| insert into t values () as new on duplicate key update a = new.a
< error
| 1064 42000 …line 1 column 26 near "as new on duplicate key update a = new.a" 
> exec ^
: string
| insert into t values (), () as new
< error
| 1064 42000 …line 1 column 30 near "as new" 
> exec ^
: string
| insert into t values (1) as new
< error
| 1064 42000 …line 1 column 27 near "as new" 
> exec ^
: string
| insert into t values (1) as new (a)
< error
| 1064 42000 …line 1 column 27 near "as new (a)" 
> exec ^
: string
| insert into t values (1) as new ()
< error
| 1064 42000 …line 1 column 27 near "as new ()" 
> exec ^
: string
| insert into t values (1) as new (a,)
< error
| 1064 42000 …line 1 column 27 near "as new (a,)" 
> exec ^
: string
| insert into t values (1) new
< error
| 1064 42000 …line 1 column 28 near "new" 
> exec ^
: string
| insert into t set a = 1 as new
< error
| 1064 42000 …line 1 column 26 near "as new" 
> exec ^
: string
| insert into t set a = 1 as new on duplicate key update a = new.a
< error
| 1064 42000 …line 1 column 26 near "as new on duplicate key update a = new.a" 
> exec ^
: string
| insert into t set a = 1 as new (a)
< error
| 1064 42000 …line 1 column 26 near "as new (a)" 
> exec ^
: string
| replace into t values (1) as new
< error
| 1064 42000 …line 1 column 28 near "as new" 
> exec ^
: string
| replace into t set a = 1 as new
< error
| 1064 42000 …line 1 column 27 near "as new" 
> exec ^
: string
| insert into t select 1 as new
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t select 1 as new
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t (a) select 1 as new from u as new
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t (a) select 1 as new from u as new
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t set a := 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t set a := 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t set a:=1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t set a:=1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t set a := default
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t set a := default
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t set a := 1, b = 2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t set a := 1, b = 2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t set a = 1, b := 2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t set a = 1, b := 2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t set t.a := 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t set t.a := 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t set `@a` := 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t set `@a` := 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| replace into t set a := 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| replace into t set a := 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| replace t set a := default
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| replace t set a := default
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t set a : = 1
< error
| 1064 42000 …line 1 column 21 near ": = 1" 
> exec ^
: string
| insert into t set a :=
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| insert into t set a := 1 on duplicate key update a := 2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t set a := 1 on duplicate key update a := 2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t set a = 1 on duplicate key update a := 2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t set a = 1 on duplicate key update a := 2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t set a = @b := 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t set a = @b := 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values (1) on duplicate key update a := 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values (1) on duplicate key update a := 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values(1) on duplicate key update a:=1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values(1) on duplicate key update a:=1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values (1) on duplicate key update a := values(a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values (1) on duplicate key update a := values(a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values (1) on duplicate key update a := default
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values (1) on duplicate key update a := default
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values (1) on duplicate key update a = 1, b := 2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values (1) on duplicate key update a = 1, b := 2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values (1) on duplicate key update t.a := 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values (1) on duplicate key update t.a := 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values (1) on duplicate key update `@a` := 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values (1) on duplicate key update `@a` := 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values (1) on duplicate key update a : = 1
< error
| 1064 42000 …line 1 column 52 near ": = 1" 
> exec ^
: string
| insert into t (a) select 1 on duplicate key update a := 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t (a) select 1 on duplicate key update a := 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values row(1) on duplicate key update a := 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t values row(1) on duplicate key update a := 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t () values () on duplicate key update a := 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t () values () on duplicate key update a := 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values (1), (2) as new (a) on duplicate key update a := new.a
< error
| 1064 42000 …line 1 column 32 near "as new (a) on duplicate key update a := new.a" 
> exec ^
: string
| INSERT INTO t (a,b,c) VALUES (1,2,3),(4,5,6) ON DUPLICATE KEY UPDATE c:=VALUES(a)+VALUES(b)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| INSERT INTO t (a,b,c) VALUES (1,2,3),(4,5,6) ON DUPLICATE KEY UPDATE c:=VALUES(a)+VALUES(b)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| INSERT IGNORE INTO t (a,b,c) VALUES (1,2,3),(4,5,6) ON DUPLICATE KEY UPDATE c:=VALUES(a)+VALUES(b)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| INSERT IGNORE INTO t (a,b,c) VALUES (1,2,3),(4,5,6) ON DUPLICATE KEY UPDATE c:=VALUES(a)+VALUES(b)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update t set a := 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t set a := 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update t set a:=1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t set a:=1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update t set a := default
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t set a := default
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update t set a := 1, b = 2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t set a := 1, b = 2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update t set t.a := 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t set t.a := 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update t, u set t.a := 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t, u set t.a := 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update t set `@a` := 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t set `@a` := 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update t set a := (select 1)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t set a := (select 1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update t set a := 1 where b = 2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t set a := 1 where b = 2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update t set a := 1 order by b limit 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t set a := 1 order by b limit 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update t set a : = 1
< error
| 1064 42000 …line 1 column 16 near ": = 1" 
> exec ^
: string
| update t set a :=
< error
| 1064 42000 …line 1 column 17 near "" 
> exec ^
: string
| update t set a := b := 1
< error
| 1064 42000 …line 1 column 22 near ":= 1" 
> exec ^
: string
| update t set a = b := 1
< error
| 1064 42000 …line 1 column 21 near ":= 1" 
> exec ^
: string
| update t set a := @b := 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t set a := @b := 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update t set a = @b := 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t set a = @b := 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update t set a = 1, b := @c := 2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t set a = 1, b := @c := 2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| explain update t set a := 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| explain update t set a := 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile 'f' into table t set a = 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile 'f' into table t set a = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile 'f' into table t set a = default
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile 'f' into table t set a = default
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile 'f' into table t set a := 1
< error
| 1064 42000 …line 1 column 42 near ":= 1" 
> exec ^
: string
| load data infile 'f' into table t set a:=1
< error
| 1064 42000 …line 1 column 41 near ":=1" 
> exec ^
: string
| load data infile 'f' into table t set a := default
< error
| 1064 42000 …line 1 column 42 near ":= default" 
> exec ^
: string
| load data infile 'f' into table t set a = 1, b := 2
< error
| 1064 42000 …line 1 column 49 near ":= 2" 
> exec ^
: string
| load data infile 'f' into table t set t.a := 1
< error
| 1064 42000 …line 1 column 44 near ":= 1" 
> exec ^
: string
| load data local infile 'f' into table t set a := 1
< error
| 1064 42000 …line 1 column 48 near ":= 1" 
> exec ^
: string
| load xml infile 'f' into table t set a := 1
< error
| 1064 42000 …line 1 column 8 near "xml infile 'f' into table t set a := 1" 
> exec ^
: string
| load data infile 'f' into table t set a : = 1
< error
| 1064 42000 …line 1 column 41 near ": = 1" 
> exec ^
: string
| load data infile 'f' into table t (a, @b) set a := @b
< error
| 1064 42000 …line 1 column 50 near ":= @b" 
> exec ^
: string
| update t set @a = 1
< error
| 1064 42000 …line 1 column 15 near "@a = 1" 
> exec ^
: string
| update t set @a := 1
< error
| 1064 42000 …line 1 column 15 near "@a := 1" 
> exec ^
: string
| update t set @@autocommit = 1
< error
| 1064 42000 …line 1 column 25 near "@@autocommit = 1" 
> exec ^
: string
| update t set @@autocommit := 1
< error
| 1064 42000 …line 1 column 25 near "@@autocommit := 1" 
> exec ^
: string
| update t set t.* = 1
< error
| 1064 42000 …line 1 column 16 near "* = 1" 
> exec ^
: string
| update t set t.* := 1
< error
| 1064 42000 …line 1 column 16 near "* := 1" 
> exec ^
: string
| update t set `@a` = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| update t set `@a` = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t set @a = 1
< error
| 1064 42000 …line 1 column 20 near "@a = 1" 
> exec ^
: string
| insert into t set @a := 1
< error
| 1064 42000 …line 1 column 20 near "@a := 1" 
> exec ^
: string
| insert into t set @@autocommit = 1
< error
| 1064 42000 …line 1 column 30 near "@@autocommit = 1" 
> exec ^
: string
| insert into t set @@autocommit := 1
< error
| 1064 42000 …line 1 column 30 near "@@autocommit := 1" 
> exec ^
: string
| insert into t set t.* = 1
< error
| 1064 42000 …line 1 column 21 near "* = 1" 
> exec ^
: string
| insert into t set `@a` = 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| insert into t set `@a` = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| insert into t values (1) on duplicate key update @a = 1
< error
| 1064 42000 …line 1 column 51 near "@a = 1" 
> exec ^
: string
| insert into t values (1) on duplicate key update @a := 1
< error
| 1064 42000 …line 1 column 51 near "@a := 1" 
> exec ^
: string
| insert into t values (1) on duplicate key update t.* = 1
< error
| 1064 42000 …line 1 column 52 near "* = 1" 
> exec ^
: string
| load data infile 'f' into table t set @a = 1
< error
| 1064 42000 …line 1 column 40 near "@a = 1" 
> exec ^
: string
| load data infile 'f' into table t set @a := 1
< error
| 1064 42000 …line 1 column 40 near "@a := 1" 
> exec ^
: string
| load data infile 'f' into table t set t.* = 1
< error
| 1064 42000 …line 1 column 41 near "* = 1" 
> exec ^
: string
| load data infile 'f' into table t set `@a` = 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile 'f' into table t set `@a` = 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a := 1
< error
| 1064 42000 …line 1 column 11 near ":= 1" 
> exec ^
: string
| select @a := 1
< ok 0 0
> exec ^
: string
| set @a := 1
< ok 0 0
> exec ^
: string
| set autocommit := 1
< ok 0 0
> exec ^
: string
| import into t set a = 1 from 's3://b/f'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t set a = 1 from 's3://b/f'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t set a := 1 from 's3://b/f'
< error
| 1064 42000 …line 1 column 22 near ":= 1 from 's3://b/f'" 
> exec ^
: string
| import into t set @a = 1 from 's3://b/f'
< error
| 1064 42000 …line 1 column 20 near "@a = 1 from 's3://b/f'" 
> exec ^
: string
| import into t set t.* = 1 from 's3://b/f'
< error
| 1064 42000 …line 1 column 21 near "* = 1 from 's3://b/f'" 
> exec ^
: string
| import into t set `@a` = 1 from 's3://b/f'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t set `@a` = 1 from 's3://b/f'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| batch on a limit 1 update t set a := 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| batch on a limit 1 update t set a := 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| batch on a limit 1 insert into t values ()
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| batch on a limit 1 insert into t values ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| batch on a limit 1 insert into t () select 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| batch on a limit 1 insert into t () select 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| batch on a limit 1 insert into t select 1 on duplicate key update a := 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| batch on a limit 1 insert into t select 1 on duplicate key update a := 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| trace insert into t values ()
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| trace insert into t values ()
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create binding for insert into t () select 1 using insert into t () select 1
< ok 0 0
