-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestLegacyTypeFormsTiDBOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| create table t (a text ascii)
< ok 0 0
> exec ^
: string
| create table t (a tinytext ascii)
< ok 0 0
> exec ^
: string
| create table t (a mediumtext ascii)
< ok 0 0
> exec ^
: string
| create table t (a longtext ascii)
< ok 0 0
> exec ^
: string
| create table t (a text(10) byte)
< ok 0 0
> exec ^
: string
| create table t (a text unicode)
< error
| 1115 42000 Unknown character set: 'ucs2'
> exec ^
: string
| create table t (a text ascii binary)
< error
| 1064 42000 …line 1 column 35 near "binary)" 
> exec ^
: string
| create table t (a text ascii charset latin1)
< error
| 1064 42000 …line 1 column 36 near "charset latin1)" 
> exec ^
: string
| create table t (a text binary ascii)
< error
| 1064 42000 …line 1 column 35 near "ascii)" 
> exec ^
: string
| create table t (a text ascii collate latin1_bin)
< ok 0 0
> exec ^
: string
| create table t (a varchar(10) ascii)
< error
| 1064 42000 …line 1 column 35 near "ascii)" 
> exec ^
: string
| create table t (a char(10) byte)
< error
| 1064 42000 …line 1 column 31 near "byte)" 
> exec ^
: string
| create table t (a enum('x') ascii)
< ok 0 0
> exec ^
: string
| create table t (a set('x') byte)
< ok 0 0
> exec ^
: string
| create table t (a long ascii)
< ok 0 0
> exec ^
: string
| create table t (a long byte)
< ok 0 0
> exec ^
: string
| create table t (a long varchar ascii)
< ok 0 0
> exec ^
: string
| create table t (a long varbinary ascii)
< error
| 1064 42000 …line 1 column 38 near "ascii)" 
> exec ^
: string
| create table t (a blob ascii)
< error
| 1064 42000 …line 1 column 28 near "ascii)" 
> exec ^
: string
| create table t (ascii int, byte int)
< ok 0 0
> exec ^
: string
| create table t (a text ascii collate)
< error
| 1064 42000 …line 1 column 37 near ")" 
> exec ^
: string
| create table t (a text byte collate not null)
< error
| 1064 42000 …line 1 column 39 near "not null)" 
> exec ^
: string
| create table t (a text collate)
< error
| 1064 42000 …line 1 column 31 near ")" 
> exec ^
: string
| create table t (a text collate not null)
< error
| 1064 42000 …line 1 column 34 near "not null)" 
> exec ^
: string
| create table t (a varchar(3) collate)
< error
| 1064 42000 …line 1 column 37 near ")" 
> exec ^
: string
| create table t (a varchar(3) charset)
< error
| 1064 42000 …line 1 column 37 near ")" 
> exec ^
: string
| create table t (a varchar(3) character set)
< error
| 1064 42000 …line 1 column 43 near ")" 
> exec ^
: string
| create table t (a varchar(3) charset not null)
< error
| 1064 42000 …line 1 column 40 near "not null)" 
> exec ^
: string
| create table t (a enum('x') charset)
< error
| 1064 42000 …line 1 column 36 near ")" 
> exec ^
: string
| create table t (a long collate)
< error
| 1064 42000 …line 1 column 31 near ")" 
> exec ^
: string
| create table t (a varchar(3) charset binary collate binary)
< ok 0 0
> exec ^
: string
| create table t (y year(4) unsigned zerofill zerofill)
< ok 0 0
> exec ^
: string
| create table t (y year signed unsigned zerofill)
< ok 0 0
> exec ^
: string
| create table t (a int signed unsigned)
< ok 0 0
> exec ^
: string
| create table t (a int zerofill unsigned)
< ok 0 0
> exec ^
: string
| create table t (a int unsigned unsigned)
< ok 0 0
> exec ^
: string
| create table t (a decimal(5,2) zerofill unsigned)
< ok 0 0
> exec ^
: string
| create table t (a double unsigned signed)
< ok 0 0
> exec ^
: string
| create table t (a bool unsigned)
< ok 0 0
> exec ^
: string
| create table t (a boolean zerofill)
< ok 0 0
> exec ^
: string
| create table t (a bit(3) unsigned)
< error
| 1064 42000 …line 1 column 33 near "unsigned)" 
> exec ^
: string
| create table t (a varchar(3) unsigned)
< error
| 1064 42000 …line 1 column 37 near "unsigned)" 
> exec ^
: string
| create table t (a int storage disk)
< ok 0 0
> exec ^
: string
| create table t (a int storage default)
< ok 0 0
> exec ^
: string
| create table t (a int storage memory)
< ok 0 0
> exec ^
: string
| create table t (a int storage)
< error
| 1064 42000 …line 1 column 30 near ")" 
> exec ^
: string
| create table t (a int storage foo)
< error
| 1064 42000 …line 1 column 33 near "foo)" 
> exec ^
: string
| create table t (a int storage disk storage memory)
< ok 0 0
> exec ^
: string
| alter table t add column a int storage disk
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add column a int storage disk
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t modify a int storage memory
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t modify a int storage memory
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
