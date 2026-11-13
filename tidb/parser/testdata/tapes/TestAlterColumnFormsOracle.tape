-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestAlterColumnFormsOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| ALTER TABLE t ALTER COLUMN b SET NOT NULL
< error
| 1064 42000 …line 1 column 36 near "NOT NULL" 
> exec ^
: string
| ALTER TABLE t ALTER b SET NOT NULL
< error
| 1064 42000 …line 1 column 29 near "NOT NULL" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b DROP NOT NULL
< error
| 1064 42000 …line 1 column 37 near "NOT NULL" 
> exec ^
: string
| ALTER TABLE t ALTER b DROP NOT NULL
< error
| 1064 42000 …line 1 column 30 near "NOT NULL" 
> exec ^
: string
| alter table t alter column b set not null
< error
| 1064 42000 …line 1 column 36 near "not null" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b SET NULL
< error
| 1064 42000 …line 1 column 37 near "NULL" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b DROP NULL
< error
| 1064 42000 …line 1 column 38 near "NULL" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b SET NOT
< error
| 1064 42000 …line 1 column 36 near "NOT" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b DROP NOT
< error
| 1064 42000 …line 1 column 37 near "NOT" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN t.b SET NOT NULL
< error
| 1064 42000 …line 1 column 38 near "NOT NULL" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN t.b DROP NOT NULL
< error
| 1064 42000 …line 1 column 39 near "NOT NULL" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b SET DEFAULT 1, ALTER COLUMN b SET NOT NULL
< error
| 1064 42000 …line 1 column 66 near "NOT NULL" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b DROP NOT NULL, ALTER COLUMN b DROP DEFAULT
< error
| 1064 42000 …line 1 column 37 near "NOT NULL, ALTER COLUMN b DROP DEFAULT" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b SET
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b DROP
< error
| 1064 42000 …line 1 column 33 near "" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b
< error
| 1064 42000 …line 1 column 28 near "" 
> exec ^
: string
| ALTER TABLE t ALTER b
< error
| 1064 42000 …line 1 column 21 near "" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b SET DEFAULT 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b SET DEFAULT 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b SET DEFAULT NULL
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER COLUMN b SET DEFAULT NULL
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER b DROP DEFAULT
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER b DROP DEFAULT
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t ALTER b DROP DEFAULT, ALTER b SET DEFAULT 2
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t ALTER b DROP DEFAULT, ALTER b SET DEFAULT 2
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t MODIFY b INT NOT NULL
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t MODIFY b INT NOT NULL
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t MODIFY COLUMN b INT NULL
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t MODIFY COLUMN b INT NULL
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| ALTER TABLE t CHANGE b b INT NOT NULL
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| ALTER TABLE t CHANGE b b INT NOT NULL
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
