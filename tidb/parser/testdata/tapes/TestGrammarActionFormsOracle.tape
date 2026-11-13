-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestGrammarActionFormsOracle default#1
> exec
| PREPARE omni_stmt FROM ?
: string
| create resource group x burstable=18446744073709551615
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| create resource group x burstable=18446744073709551615
| ; SELECT FROM
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| create resource group x burstable=9223372036854775808
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| create resource group x burstable=9223372036854775808
| ; SELECT FROM
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| create resource group x burstable=9223372036854775807
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT = (RU 18446744073709551615)
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT = (RU 18446744073709551615)
| ; SELECT FROM
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT = (RU 18446744073709551615 action=kill)
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT = (RU 18446744073709551615 action=kill)
| ; SELECT FROM
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT = (PROCESSED_KEYS 9223372036854775808 action=kill)
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT = (PROCESSED_KEYS 9223372036854775808 action=kill)
| ; SELECT FROM
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| create resource group x ru_per_sec=18446744073709551615
< ok 0 0
> exec ^
: string
| query watch remove 9223372036854775808
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| query watch remove 9223372036854775808
| ; SELECT FROM
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| query watch remove 18446744073709551615
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| query watch remove 18446744073709551615
| ; SELECT FROM
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| query watch remove 9223372036854775807
< ok 0 0
> exec ^
: string
| recover table by job 9223372036854775808
< error
| 1064 42000 …line 1 column 40 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| recover table by job 9223372036854775807
< ok 0 0
> exec ^
: string
| recover table t 9223372036854775808
< error
| 1064 42000 …line 1 column 35 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| cancel distribution job 9223372036854775808
< error
| 1064 42000 …line 1 column 43 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| split table t between (1) and (2) regions 9223372036854775808
< error
| 1064 42000 …line 1 column 61 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| split table t between (1) and (2) regions 9223372036854775807
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| split table t between (1) and (2) regions 9223372036854775807
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| show br job 9223372036854775808
< error
| 1064 42000 …line 1 column 31 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| show br job query 9223372036854775808
< error
| 1064 42000 …line 1 column 37 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| cancel br job 9223372036854775808
< error
| 1064 42000 …line 1 column 33 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| cancel import job 9223372036854775808
< error
| 1064 42000 …line 1 column 37 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| admin show ddl jobs 9223372036854775808
< error
| 1064 42000 …line 1 column 39 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| admin alter ddl jobs 9223372036854775808 thread=1
< error
| 1064 42000 …line 1 column 47 near "thread=1"9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| admin cancel ddl jobs 9223372036854775808
< error
| 1064 42000 …line 1 column 41 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| admin pause ddl jobs 9223372036854775808
< error
| 1064 42000 …line 1 column 40 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| admin resume ddl jobs 9223372036854775808
< error
| 1064 42000 …line 1 column 41 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| admin check index t i (1,2)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| admin check index t i (1,2)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| admin check index t i (9223372036854775808,2)
< error
| 1064 42000 …line 1 column 43 near ",2)"9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| show import job 9223372036854775808
< error
| 1064 42000 …line 1 column 35 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| show distribution job 9223372036854775808
< error
| 1064 42000 …line 1 column 41 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| show profile for query 9223372036854775808
< error
| 1064 42000 …line 1 column 42 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| create user u with max_queries_per_hour 9223372036854775808
< error
| 1064 42000 …line 1 column 59 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| create user u with max_queries_per_hour 9223372036854775807
< ok 0 0
> exec ^
: string
| create user u with max_user_connections 9223372036854775808
< error
| 1064 42000 …line 1 column 59 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| create user u password expire interval 9223372036854775808 day
< error
| 1064 42000 …line 1 column 62 near "day"9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| create user u failed_login_attempts 9223372036854775808
< error
| 1064 42000 …line 1 column 55 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| create user u password_lock_time 9223372036854775808
< error
| 1064 42000 …line 1 column 52 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| show profile for query foo
< error
| 1064 42000 …line 1 column 26 near "foo" 
> exec ^
: string
| show profile for query 1+1
< error
| 1064 42000 …line 1 column 25 near "+1" 
> exec ^
: string
| show profile for query 1
< ok 0 0
> exec ^
: string
| show profile for query
< error
| 1064 42000 …line 1 column 22 near "" 
> exec ^
: string
| create user u with max_queries_per_hour
< error
| 1064 42000 …line 1 column 39 near "" 
> exec ^
: string
| create user u with max_user_connections 'x'
< error
| 1064 42000 …line 1 column 43 near "'x'" 
> exec ^
: string
| create user u password expire interval day
< error
| 1064 42000 …line 1 column 42 near "day" 
> exec ^
: string
| create user u password expire interval 1
< error
| 1064 42000 …line 1 column 40 near "" 
> exec ^
: string
| create user u failed_login_attempts
< error
| 1064 42000 …line 1 column 35 near "" 
> exec ^
: string
| create user u password_lock_time
< error
| 1064 42000 …line 1 column 32 near "" 
> exec ^
: string
| alter user u with max_queries_per_hour 9223372036854775808
< error
| 1064 42000 …line 1 column 58 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| alter user u failed_login_attempts 9223372036854775808
< error
| 1064 42000 …line 1 column 54 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| create resource group x ru_per_sec=1 ru_per_sec=2
< error
| 1064 42000 …line 1 column 49 near ""Dupliated options specified 
> exec ^
: string
| create resource group x ru_per_sec=1, ru_per_sec=2
< error
| 1064 42000 …line 1 column 50 near ""Dupliated options specified 
> exec ^
: string
| create resource group x ru_per_sec=1 priority=low priority=high
< error
| 1064 42000 …line 1 column 63 near ""Dupliated options specified 
> exec ^
: string
| create resource group x ru_per_sec=1 query_limit=(exec_elapsed='1s' action=kill action=dryrun)
< error
| 1064 42000 …line 1 column 94 near ")"Dupliated runaway options specified 
> exec ^
: string
| create resource group x ru_per_sec=1 query_limit=(exec_elapsed='1s', exec_elapsed='2s' action=kill)
< ok 0 0
> exec ^
: string
| create resource group x ru_per_sec=1 query_limit=(exec_elapsed='foo' action=kill)
< error
| 1064 42000 …line 1 column 75 near "action=kill)"The EXEC_ELAPSED option is not a valid duration: time: invalid duration "foo" 
> exec ^
: string
| create resource group x ru_per_sec=1 query_limit=(exec_elapsed='1s' action=kill watch similar duration 'foo')
< error
| 1064 42000 …line 1 column 109 near ")"The WATCH DURATION option is not a valid duration: time: invalid duration "foo" 
> exec ^
: string
| create resource group x ru_per_sec=1 background=(task_types='a' task_types='b')
< error
| 1064 42000 …line 1 column 79 near ")"Dupliated background options specified 
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT = (RU 9223372036854775808 action=kill)
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| create resource group x ru_per_sec=1000 QUERY_LIMIT = (RU 9223372036854775808 action=kill)
| ; SELECT FROM
< error
| 1105 HY000 interface conversion: interface {} is uint64, not int64
> exec ^
: string
| query watch add action kill action dryrun sql text exact to 'x'
< error
| 1064 42000 …line 1 column 45 near "sql text exact to 'x'"Dupliated options specified 
> exec ^
: string
| query watch add resource group a resource group b sql text exact to 'x'
< error
| 1064 42000 …line 1 column 53 near "sql text exact to 'x'"Dupliated options specified 
> exec ^
: string
| calibrate resource start_time '2020-01-01' start_time '2020-01-01'
< error
| 1064 42000 …line 1 column 66 near ""Dupliated options specified 
> exec ^
: string
| calibrate resource duration 'foo'
< error
| 1064 42000 …line 1 column 33 near ""The DURATION option is not a valid duration: fail to read an integer 
> exec ^
: string
| calibrate resource duration '1h' duration '1h'
< error
| 1064 42000 …line 1 column 46 near ""Dupliated options specified 
> exec ^
: string
| calibrate resource start_time '2020-01-01' end_time '2020-01-02' duration '1h'
< ok 0 0
> exec ^
: string
| calibrate resource start_time '2020-01-01' duration '1h' start_time '2020-01-01'
< error
| 1064 42000 …line 1 column 80 near ""Dupliated options specified 
> exec ^
: string
| create placement policy x followers=0
< error
| 1064 42000 …line 1 column 37 near ""FOLLOWERS must be positive 
> exec ^
: string
| create placement policy x followers=1
< ok 0 0
> exec ^
: string
| alter placement policy x followers=0
< error
| 1064 42000 …line 1 column 36 near ""FOLLOWERS must be positive 
> exec ^
: string
| create table t (a int) followers=0
< error
| 1064 42000 …line 1 column 32 near "followers=0" 
> exec ^
: string
| create database d followers=0
< error
| 1064 42000 …line 1 column 27 near "followers=0" 
> exec ^
: string
| alter table t algorithm=foo
< error
| 1800 HY000 Unknown ALGORITHM 'algorithm'
> exec ^
: string
| alter table t algorithm=foo
| ; SELECT FROM
< error
| 1800 HY000 Unknown ALGORITHM 'algorithm'
> exec ^
: string
| alter table t algorithm foo
< error
| 1800 HY000 Unknown ALGORITHM 'algorithm'
> exec ^
: string
| alter table t algorithm foo
| ; SELECT FROM
< error
| 1800 HY000 Unknown ALGORITHM 'algorithm'
> exec ^
: string
| alter table t algorithm=`inplace`
< error
| 1800 HY000 Unknown ALGORITHM 'algorithm'
> exec ^
: string
| alter table t algorithm=`inplace`
| ; SELECT FROM
< error
| 1800 HY000 Unknown ALGORITHM 'algorithm'
> exec ^
: string
| alter table t algorithm=`default`
< error
| 1800 HY000 Unknown ALGORITHM 'algorithm'
> exec ^
: string
| alter table t algorithm=`default`
| ; SELECT FROM
< error
| 1800 HY000 Unknown ALGORITHM 'algorithm'
> exec ^
: string
| alter table t algorithm=inplace
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t algorithm=inplace
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t algorithm default
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t algorithm default
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t algorithm=copy
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t algorithm=copy
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t algorithm=instant
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t algorithm=instant
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t lock=foo
< error
| 1801 HY000 Unknown LOCK type 'foo'
> exec ^
: string
| alter table t lock=foo
| ; SELECT FROM
< error
| 1801 HY000 Unknown LOCK type 'foo'
> exec ^
: string
| alter table t lock foo
< error
| 1801 HY000 Unknown LOCK type 'foo'
> exec ^
: string
| alter table t lock foo
| ; SELECT FROM
< error
| 1801 HY000 Unknown LOCK type 'foo'
> exec ^
: string
| alter table t lock=none
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t lock=none
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t lock=`none`
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t lock=`none`
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t lock=`default`
< error
| 1801 HY000 Unknown LOCK type 'default'
> exec ^
: string
| alter table t lock=`default`
| ; SELECT FROM
< error
| 1801 HY000 Unknown LOCK type 'default'
> exec ^
: string
| alter table t lock default
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t lock default
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t lock=shared
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t lock=shared
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t lock=exclusive
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t lock=exclusive
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t lock=NONE
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t lock=NONE
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create index i on t(a) algorithm=foo
< error
| 1800 HY000 Unknown ALGORITHM 'algorithm'
> exec ^
: string
| create index i on t(a) algorithm=foo
| ; SELECT FROM
< error
| 1800 HY000 Unknown ALGORITHM 'algorithm'
> exec ^
: string
| create index i on t(a) algorithm=`inplace`
< error
| 1800 HY000 Unknown ALGORITHM 'algorithm'
> exec ^
: string
| create index i on t(a) algorithm=`inplace`
| ; SELECT FROM
< error
| 1800 HY000 Unknown ALGORITHM 'algorithm'
> exec ^
: string
| create index i on t(a) lock=foo
< error
| 1801 HY000 Unknown LOCK type 'foo'
> exec ^
: string
| create index i on t(a) lock=foo
| ; SELECT FROM
< error
| 1801 HY000 Unknown LOCK type 'foo'
> exec ^
: string
| create index i on t(a) lock=`none`
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| create index i on t(a) lock=`none`
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| drop index i on t algorithm=foo
< error
| 1800 HY000 Unknown ALGORITHM 'algorithm'
> exec ^
: string
| drop index i on t algorithm=foo
| ; SELECT FROM
< error
| 1800 HY000 Unknown ALGORITHM 'algorithm'
> exec ^
: string
| drop index i on t lock=foo
< error
| 1801 HY000 Unknown LOCK type 'foo'
> exec ^
: string
| drop index i on t lock=foo
| ; SELECT FROM
< error
| 1801 HY000 Unknown LOCK type 'foo'
> exec ^
: string
| drop index i on t algorithm=inplace
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| drop index i on t algorithm=inplace
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| drop index i on t lock=none
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| drop index i on t lock=none
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) partition by key algorithm=1 (a)
< ok 0 0
> exec ^
: string
| alter table t algorithm='inplace'
< error
| 1064 42000 …line 1 column 33 near "'inplace'" 
> exec ^
: string
| alter table t algorithm=1
< error
| 1064 42000 …line 1 column 25 near "1" 
> exec ^
: string
| alter table t lock='none'
< error
| 1064 42000 …line 1 column 25 near "'none'" 
> exec ^
: string
| alter table t lock=1
< error
| 1064 42000 …line 1 column 20 near "1" 
> exec ^
: string
| alter table t algorithm=select
< error
| 1064 42000 …line 1 column 30 near "select" 
> exec ^
: string
| alter table t lock=select
< error
| 1064 42000 …line 1 column 25 near "select" 
> exec ^
: string
| flashback cluster to tso 0
< error
| 1064 42000 …line 1 column 26 near ""Invalid TSO value provided: 0 
> exec ^
: string
| flashback cluster to tso 1
< ok 0 0
> exec ^
: string
| flashback table t to tso 0
< error
| 1064 42000 …line 1 column 26 near ""Invalid TSO value provided: 0 
> exec ^
: string
| flashback table t to tso 1
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| flashback table t to tso 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| flashback database d to tso 0
< error
| 1064 42000 …line 1 column 29 near ""Invalid TSO value provided: 0 
> exec ^
: string
| flashback schema d to tso 0
< error
| 1064 42000 …line 1 column 27 near ""Invalid TSO value provided: 0 
> exec ^
: string
| flashback database d to tso 1
< ok 0 0
> exec ^
: string
| create table t (a int, b int as (a) default 1)
< error
| 1221 HY000 Incorrect usage of DEFAULT and generated column
> exec ^
: string
| create table t (a int, b int as (a) default 1)
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of DEFAULT and generated column
> exec ^
: string
| create table t (a int, b int as (a) auto_increment)
< error
| 1221 HY000 Incorrect usage of AUTO_INCREMENT and generated column
> exec ^
: string
| create table t (a int, b int as (a) auto_increment)
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of AUTO_INCREMENT and generated column
> exec ^
: string
| create table t (a int, b timestamp as (a) on update current_timestamp)
< error
| 1221 HY000 Incorrect usage of ON UPDATE and generated column
> exec ^
: string
| create table t (a int, b timestamp as (a) on update current_timestamp)
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of ON UPDATE and generated column
> exec ^
: string
| create table t (a int, b int as (a) not null)
< ok 0 0
> exec ^
: string
| create table t (a int, b int as (a) comment 'c')
< ok 0 0
> exec ^
: string
| create table t (a int, b int as (a) primary key)
< error
| 3106 HY000 'Defining a virtual generated column as primary key' is not supported for generated columns.
> exec ^
: string
| create table t (a int, b int as (a) primary key)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int, b int default 1 as (a))
< error
| 1221 HY000 Incorrect usage of DEFAULT and generated column
> exec ^
: string
| create table t (a int, b int default 1 as (a))
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of DEFAULT and generated column
> exec ^
: string
| create table t (a int, b int auto_increment as (a))
< error
| 1221 HY000 Incorrect usage of AUTO_INCREMENT and generated column
> exec ^
: string
| create table t (a int, b int auto_increment as (a))
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of AUTO_INCREMENT and generated column
> exec ^
: string
| alter table t add column b int as (a) default 1
< error
| 1221 HY000 Incorrect usage of DEFAULT and generated column
> exec ^
: string
| alter table t add column b int as (a) default 1
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of DEFAULT and generated column
> exec ^
: string
| alter table t add column (b int as (a) default 1)
< error
| 1221 HY000 Incorrect usage of DEFAULT and generated column
> exec ^
: string
| alter table t add column (b int as (a) default 1)
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of DEFAULT and generated column
> exec ^
: string
| alter table t modify b int as (a) default 1
< error
| 1221 HY000 Incorrect usage of DEFAULT and generated column
> exec ^
: string
| alter table t modify b int as (a) default 1
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of DEFAULT and generated column
> exec ^
: string
| alter table t change a b int as (a) default 1
< error
| 1221 HY000 Incorrect usage of DEFAULT and generated column
> exec ^
: string
| alter table t change a b int as (a) default 1
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of DEFAULT and generated column
> exec ^
: string
| create table t (a int, b int as (a) serial default value)
< error
| 1221 HY000 Incorrect usage of AUTO_INCREMENT and generated column
> exec ^
: string
| create table t (a int, b int as (a) serial default value)
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of AUTO_INCREMENT and generated column
> exec ^
: string
| create table t (a int, b int serial default value as (a))
< error
| 1221 HY000 Incorrect usage of AUTO_INCREMENT and generated column
> exec ^
: string
| create table t (a int, b int serial default value as (a))
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of AUTO_INCREMENT and generated column
> exec ^
: string
| create table t (a int, b int generated always as (a) default 1)
< error
| 1221 HY000 Incorrect usage of DEFAULT and generated column
> exec ^
: string
| create table t (a int, b int generated always as (a) default 1)
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of DEFAULT and generated column
> exec ^
: string
| create table t (a int, b int default null generated always as (a))
< error
| 1221 HY000 Incorrect usage of DEFAULT and generated column
> exec ^
: string
| create table t (a int, b int default null generated always as (a))
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of DEFAULT and generated column
> exec ^
: string
| create global temporary table t (a int)
< error
| 1064 42000 …line 1 column 39 near ""GLOBAL TEMPORARY and ON COMMIT DELETE ROWS must appear together 
> exec ^
: string
| create global temporary table t (a int) on commit delete rows
< ok 0 0
> exec ^
: string
| create global temporary table t (a int) on commit preserve rows
< ok 0 0
> exec ^
: string
| create temporary table t (a int) on commit delete rows
< error
| 1064 42000 …line 1 column 54 near ""GLOBAL TEMPORARY and ON COMMIT DELETE ROWS must appear together 
> exec ^
: string
| create temporary table t (a int) on commit preserve rows
< error
| 1064 42000 …line 1 column 56 near ""GLOBAL TEMPORARY and ON COMMIT DELETE ROWS must appear together 
> exec ^
: string
| create table t (a int) on commit delete rows
< error
| 1064 42000 …line 1 column 44 near ""GLOBAL TEMPORARY and ON COMMIT DELETE ROWS must appear together 
> exec ^
: string
| create table t (a int) on commit preserve rows
< error
| 1064 42000 …line 1 column 46 near ""GLOBAL TEMPORARY and ON COMMIT DELETE ROWS must appear together 
> exec ^
: string
| create global temporary table t like u
< error
| 1064 42000 …line 1 column 38 near ""GLOBAL TEMPORARY and ON COMMIT DELETE ROWS must appear together 
> exec ^
: string
| create global temporary table t like u on commit delete rows
< error
| 1146 42S02 Table 'test.u' doesn't exist
> exec ^
: string
| create global temporary table t like u on commit delete rows
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create temporary table t like u on commit delete rows
< error
| 1064 42000 …line 1 column 53 near ""GLOBAL TEMPORARY and ON COMMIT DELETE ROWS must appear together 
> exec ^
: string
| create table t like u on commit delete rows
< error
| 1064 42000 …line 1 column 43 near ""GLOBAL TEMPORARY and ON COMMIT DELETE ROWS must appear together 
> exec ^
: string
| create global temporary table t select 1
< error
| 1064 42000 …line 1 column 40 near ""GLOBAL TEMPORARY and ON COMMIT DELETE ROWS must appear together 
> exec ^
: string
| create global temporary table t select 1 on commit delete rows
< error
| 1105 HY000 'CREATE TABLE ... SELECT' is not implemented yet
> exec ^
: string
| create global temporary table t select 1 on commit delete rows
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 0
< error
| 1504 HY000 Number of partitions = 0 is not an allowed value
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 0
| ; SELECT FROM
< error
| 1504 HY000 Number of partitions = 0 is not an allowed value
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 1
< ok 0 0
> exec ^
: string
| create table t (a int) partition by key(a) partitions 0
< error
| 1504 HY000 Number of partitions = 0 is not an allowed value
> exec ^
: string
| create table t (a int) partition by key(a) partitions 0
| ; SELECT FROM
< error
| 1504 HY000 Number of partitions = 0 is not an allowed value
> exec ^
: string
| create table t (a int) partition by range(a) subpartition by hash(a) subpartitions 0 (partition p values less than (1))
< error
| 1504 HY000 Number of subpartitions = 0 is not an allowed value
> exec ^
: string
| create table t (a int) partition by range(a) subpartition by hash(a) subpartitions 0 (partition p values less than (1))
| ; SELECT FROM
< error
| 1504 HY000 Number of subpartitions = 0 is not an allowed value
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 2 (partition x)
< error
| 1484 HY000 Wrong number of partitions defined, mismatch with previous setting
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 2 (partition x)
| ; SELECT FROM
< error
| 1484 HY000 Wrong number of partitions defined, mismatch with previous setting
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 2 (partition x, partition y)
< ok 0 0
> exec ^
: string
| create table t (a int) partition by range(a) (partition x values less than (1) (subpartition a), partition y values less than (2))
< error
| 1485 HY000 Wrong number of subpartitions defined, mismatch with previous setting
> exec ^
: string
| create table t (a int) partition by range(a) (partition x values less than (1) (subpartition a), partition y values less than (2))
| ; SELECT FROM
< error
| 1485 HY000 Wrong number of subpartitions defined, mismatch with previous setting
> exec ^
: string
| create table t (a int) partition by range(a) (partition x values less than (1) (subpartition a), partition y values less than (2) (subpartition b))
< error
| 1500 HY000 It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
> exec ^
: string
| create table t (a int) partition by range(a) (partition x values less than (1) (subpartition a), partition y values less than (2) (subpartition b))
| ; SELECT FROM
< error
| 1500 HY000 It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
> exec ^
: string
| create table t (a int) partition by range(a) subpartition by hash(a) subpartitions 2 (partition x values less than (1) (subpartition a))
< error
| 1485 HY000 Wrong number of subpartitions defined, mismatch with previous setting
> exec ^
: string
| create table t (a int) partition by range(a) subpartition by hash(a) subpartitions 2 (partition x values less than (1) (subpartition a))
| ; SELECT FROM
< error
| 1485 HY000 Wrong number of subpartitions defined, mismatch with previous setting
> exec ^
: string
| create table t (a int) partition by range(a) subpartition by hash(a) subpartitions 1 (partition x values less than (1) (subpartition a))
< ok 0 0
> exec ^
: string
| create table t (a int) partition by hash(a) (partition x (subpartition a))
< error
| 1500 HY000 It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
> exec ^
: string
| create table t (a int) partition by hash(a) (partition x (subpartition a))
| ; SELECT FROM
< error
| 1500 HY000 It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
> exec ^
: string
| create table t (a int) partition by range(a)
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| create table t (a int) partition by range(a)
| ; SELECT FROM
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| create table t (a int) partition by list(a)
< error
| 1492 HY000 For LIST partitions each partition must be defined
> exec ^
: string
| create table t (a int) partition by list(a)
| ; SELECT FROM
< error
| 1492 HY000 For LIST partitions each partition must be defined
> exec ^
: string
| create table t (a int) partition by range columns(a)
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| create table t (a int) partition by range columns(a)
| ; SELECT FROM
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| create table t (a int) partition by list columns(a)
< error
| 1492 HY000 For LIST partitions each partition must be defined
> exec ^
: string
| create table t (a int) partition by list columns(a)
| ; SELECT FROM
< error
| 1492 HY000 For LIST partitions each partition must be defined
> exec ^
: string
| create table t (a int) partition by range(a) partitions 3
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| create table t (a int) partition by range(a) partitions 3
| ; SELECT FROM
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| create table t (a int) partition by hash(a)
< ok 0 0
> exec ^
: string
| create table t (a int) partition by range(a) (partition x)
< error
| 1479 HY000 Syntax : RANGE PARTITIONING requires definition of VALUES LESS THAN for each partition
> exec ^
: string
| create table t (a int) partition by range(a) (partition x)
| ; SELECT FROM
< error
| 1479 HY000 Syntax : RANGE PARTITIONING requires definition of VALUES LESS THAN for each partition
> exec ^
: string
| create table t (a int) partition by list(a) (partition x)
< error
| 1479 HY000 Syntax : LIST PARTITIONING requires definition of VALUES IN for each partition
> exec ^
: string
| create table t (a int) partition by list(a) (partition x)
| ; SELECT FROM
< error
| 1479 HY000 Syntax : LIST PARTITIONING requires definition of VALUES IN for each partition
> exec ^
: string
| create table t (a int) partition by hash(a) (partition x values less than (1))
< error
| 1480 HY000 Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
> exec ^
: string
| create table t (a int) partition by hash(a) (partition x values less than (1))
| ; SELECT FROM
< error
| 1480 HY000 Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
> exec ^
: string
| create table t (a int) partition by hash(a) (partition x values in (1))
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES IN in partition definition
> exec ^
: string
| create table t (a int) partition by hash(a) (partition x values in (1))
| ; SELECT FROM
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES IN in partition definition
> exec ^
: string
| create table t (a int) partition by key(a) (partition x values less than (1))
< error
| 1480 HY000 Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
> exec ^
: string
| create table t (a int) partition by key(a) (partition x values less than (1))
| ; SELECT FROM
< error
| 1480 HY000 Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
> exec ^
: string
| create table t (a int) partition by list(a) (partition x values less than (1))
< error
| 1480 HY000 Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
> exec ^
: string
| create table t (a int) partition by list(a) (partition x values less than (1))
| ; SELECT FROM
< error
| 1480 HY000 Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
> exec ^
: string
| create table t (a int) partition by range(a) (partition x values in (1))
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES IN in partition definition
> exec ^
: string
| create table t (a int) partition by range(a) (partition x values in (1))
| ; SELECT FROM
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES IN in partition definition
> exec ^
: string
| create table t (a int) partition by range(a) (partition x values less than (1, 2))
< error
| 1657 HY000 Cannot have more than one value for this type of RANGE partitioning
> exec ^
: string
| create table t (a int) partition by range(a) (partition x values less than (1, 2))
| ; SELECT FROM
< error
| 1657 HY000 Cannot have more than one value for this type of RANGE partitioning
> exec ^
: string
| create table t (a int) partition by range(a) (partition x values less than maxvalue)
< ok 0 0
> exec ^
: string
| create table t (a int) partition by range(a) (partition x values less than (maxvalue))
< ok 0 0
> exec ^
: string
| create table t (a int, b int) partition by range columns(a,b) (partition x values less than (1))
< error
| 1653 HY000 Inconsistency in usage of column lists for partitioning
> exec ^
: string
| create table t (a int, b int) partition by range columns(a,b) (partition x values less than (1))
| ; SELECT FROM
< error
| 1653 HY000 Inconsistency in usage of column lists for partitioning
> exec ^
: string
| create table t (a int, b int) partition by range columns(a,b) (partition x values less than (1, 2))
< ok 0 0
> exec ^
: string
| create table t (a int, b int) partition by range columns(a,b) (partition x values less than maxvalue)
< error
| 1653 HY000 Inconsistency in usage of column lists for partitioning
> exec ^
: string
| create table t (a int, b int) partition by range columns(a,b) (partition x values less than maxvalue)
| ; SELECT FROM
< error
| 1653 HY000 Inconsistency in usage of column lists for partitioning
> exec ^
: string
| create table t (a int, b int) partition by range columns(a,b) (partition x values less than (maxvalue, maxvalue))
< ok 0 0
> exec ^
: string
| create table t (a int) partition by list(a) (partition x values in ((1, 2)))
< error
| 1658 HY000 Row expressions in VALUES IN only allowed for multi-field column partitioning
> exec ^
: string
| create table t (a int) partition by list(a) (partition x values in ((1, 2)))
| ; SELECT FROM
< error
| 1658 HY000 Row expressions in VALUES IN only allowed for multi-field column partitioning
> exec ^
: string
| create table t (a int) partition by list(a) (partition x values in (1, 2))
< ok 0 0
> exec ^
: string
| create table t (a int, b int) partition by list columns(a,b) (partition x values in (1, 2))
< error
| 1653 HY000 Inconsistency in usage of column lists for partitioning
> exec ^
: string
| create table t (a int, b int) partition by list columns(a,b) (partition x values in (1, 2))
| ; SELECT FROM
< error
| 1653 HY000 Inconsistency in usage of column lists for partitioning
> exec ^
: string
| create table t (a int, b int) partition by list columns(a,b) (partition x values in ((1, 2), (3, 4)))
< ok 0 0
> exec ^
: string
| create table t (a int, b int) partition by list columns(a,b) (partition x values in ((1, 2), 3))
< error
| 1653 HY000 Inconsistency in usage of column lists for partitioning
> exec ^
: string
| create table t (a int, b int) partition by list columns(a,b) (partition x values in ((1, 2), 3))
| ; SELECT FROM
< error
| 1653 HY000 Inconsistency in usage of column lists for partitioning
> exec ^
: string
| create table t (a int, b int) partition by list columns(a,b) (partition x values in ((1, 2), (3, 4, 5)))
< error
| 1653 HY000 Inconsistency in usage of column lists for partitioning
> exec ^
: string
| create table t (a int, b int) partition by list columns(a,b) (partition x values in ((1, 2), (3, 4, 5)))
| ; SELECT FROM
< error
| 1653 HY000 Inconsistency in usage of column lists for partitioning
> exec ^
: string
| create table t (a int, b int) partition by list columns(a,b) (partition x values in ((1, 2), default))
< ok 0 0
> exec ^
: string
| create table t (a int, b int) partition by list columns(a,b) (partition x values in (default))
< ok 0 0
> exec ^
: string
| create table t (a int, b int) partition by list columns(a,b) (partition x default)
< ok 0 0
> exec ^
: string
| create table t (a int) partition by list(a) (partition x default)
< ok 0 0
> exec ^
: string
| create table t (a int) partition by list(a) (partition x values in (default))
< ok 0 0
> exec ^
: string
| create table t (a int) partition by list(a) (partition x values in (default, 1))
< ok 0 0
> exec ^
: string
| create table t (a int) partition by list(a) (partition x values in (1, default))
< ok 0 0
> exec ^
: string
| create table t (a int) partition by hash(a) (partition x default)
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES IN in partition definition
> exec ^
: string
| create table t (a int) partition by hash(a) (partition x default)
| ; SELECT FROM
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES IN in partition definition
> exec ^
: string
| create table t (a int) partition by range(a) (partition x default)
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES IN in partition definition
> exec ^
: string
| create table t (a int) partition by range(a) (partition x default)
| ; SELECT FROM
< error
| 1480 HY000 Only LIST PARTITIONING can use VALUES IN in partition definition
> exec ^
: string
| create table t (a int) partition by system_time
< error
| 4128 HY000 Wrong Partitions: must have at least one HISTORY and exactly one last CURRENT
> exec ^
: string
| create table t (a int) partition by system_time
| ; SELECT FROM
< error
| 4128 HY000 Wrong Partitions: must have at least one HISTORY and exactly one last CURRENT
> exec ^
: string
| create table t (a int) partition by system_time (partition x history)
< error
| 4128 HY000 Wrong Partitions: must have at least one HISTORY and exactly one last CURRENT
> exec ^
: string
| create table t (a int) partition by system_time (partition x history)
| ; SELECT FROM
< error
| 4128 HY000 Wrong Partitions: must have at least one HISTORY and exactly one last CURRENT
> exec ^
: string
| create table t (a int) partition by system_time (partition x values less than (1), partition y current)
< error
| 1480 HY000 Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
> exec ^
: string
| create table t (a int) partition by system_time (partition x values less than (1), partition y current)
| ; SELECT FROM
< error
| 1480 HY000 Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
> exec ^
: string
| create table t (a int) partition by system_time (partition x, partition y)
< error
| 4128 HY000 Wrong Partitions: must have at least one HISTORY and exactly one last CURRENT
> exec ^
: string
| create table t (a int) partition by system_time (partition x, partition y)
| ; SELECT FROM
< error
| 4128 HY000 Wrong Partitions: must have at least one HISTORY and exactly one last CURRENT
> exec ^
: string
| create table t (a int) partition by hash(a) (partition x history, partition y current)
< error
| 4113 HY000 Wrong partitioning type, expected type: `SYSTEM_TIME`
> exec ^
: string
| create table t (a int) partition by hash(a) (partition x history, partition y current)
| ; SELECT FROM
< error
| 4113 HY000 Wrong partitioning type, expected type: `SYSTEM_TIME`
> exec ^
: string
| create table t (a int) partition by range(a) (partition x history, partition y current)
< error
| 4113 HY000 Wrong partitioning type, expected type: `SYSTEM_TIME`
> exec ^
: string
| create table t (a int) partition by range(a) (partition x history, partition y current)
| ; SELECT FROM
< error
| 4113 HY000 Wrong partitioning type, expected type: `SYSTEM_TIME`
> exec ^
: string
| alter table t partition by range(a)
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| alter table t partition by range(a)
| ; SELECT FROM
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| alter table t partition by range(a) update indexes (a local)
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| alter table t partition by range(a) update indexes (a local)
| ; SELECT FROM
< error
| 1492 HY000 For RANGE partitions each partition must be defined
> exec ^
: string
| alter table t partition by hash(a)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t partition by hash(a)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t partition by hash(a) partitions 0
< error
| 1504 HY000 Number of partitions = 0 is not an allowed value
> exec ^
: string
| alter table t partition by hash(a) partitions 0
| ; SELECT FROM
< error
| 1504 HY000 Number of partitions = 0 is not an allowed value
> exec ^
: string
| alter table t partition by hash(a) partitions 2 (partition x)
< error
| 1484 HY000 Wrong number of partitions defined, mismatch with previous setting
> exec ^
: string
| alter table t partition by hash(a) partitions 2 (partition x)
| ; SELECT FROM
< error
| 1484 HY000 Wrong number of partitions defined, mismatch with previous setting
> exec ^
: string
| alter table t partition by range(a) (partition x)
< error
| 1479 HY000 Syntax : RANGE PARTITIONING requires definition of VALUES LESS THAN for each partition
> exec ^
: string
| alter table t partition by range(a) (partition x)
| ; SELECT FROM
< error
| 1479 HY000 Syntax : RANGE PARTITIONING requires definition of VALUES LESS THAN for each partition
> exec ^
: string
| alter table t partition by list(a) (partition x values less than (1))
< error
| 1480 HY000 Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
> exec ^
: string
| alter table t partition by list(a) (partition x values less than (1))
| ; SELECT FROM
< error
| 1480 HY000 Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
> exec ^
: string
| alter table t partition by linear key algorithm=3 (a)
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| alter table t partition by linear key algorithm=3 (a)
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| alter table t add partition (partition x)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add partition (partition x)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add partition partitions 0
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add partition partitions 0
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add partition (partition x values less than (1, 2))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add partition (partition x values less than (1, 2))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add partition (partition x values in (1) (subpartition a))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add partition (partition x values in (1) (subpartition a))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t reorganize partition a into (partition x)
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t reorganize partition a into (partition x)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t coalesce partition 0
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t coalesce partition 0
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a int) partition by key algorithm=2 (a)
< ok 0 0
> exec ^
: string
| create table t (a int) partition by key algorithm=0 (a)
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create table t (a int) partition by key algorithm=0 (a)
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create table t (a int) partition by key algorithm=3 (a)
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create table t (a int) partition by key algorithm=3 (a)
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create table t (a int) partition by linear key algorithm=0 (a)
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create table t (a int) partition by linear key algorithm=0 (a)
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create table t (a int) partition by linear key algorithm=3 (a) partitions 4
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create table t (a int) partition by linear key algorithm=3 (a) partitions 4
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create table t (a int) partition by range(a) subpartition by key algorithm=3 (a) (partition x values less than (1))
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create table t (a int) partition by range(a) subpartition by key algorithm=3 (a) (partition x values less than (1))
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create table t (a int) partition by key algorithm 2 (a)
< error
| 1064 42000 …line 1 column 51 near "2 (a)" 
> exec ^
: string
| create table t (a int) partition by key algorithm = (a)
< error
| 1064 42000 …line 1 column 53 near "(a)" 
> exec ^
: string
| create table t (a int) partition by key algorithm = 'x' (a)
< error
| 1064 42000 …line 1 column 55 near "'x' (a)" 
> exec ^
: string
| create table t (a int) partition by key
< error
| 1064 42000 …line 1 column 39 near "" 
> exec ^
: string
| create table t (a int) partition by key partitions 2
< error
| 1064 42000 …line 1 column 50 near "partitions 2" 
> exec ^
: string
| create table t (a int) partition by hash(a) partitions
< error
| 1064 42000 …line 1 column 54 near "" 
> exec ^
: string
| create table t (a int) partition by hash(a) partitions 'x'
< error
| 1064 42000 …line 1 column 58 near "'x'" 
> exec ^
: string
| create table t (a int) partition by range(a) subpartition by hash(a) subpartitions (partition p values less than (1))
< error
| 1064 42000 …line 1 column 84 near "(partition p values less than (1))" 
> exec ^
: string
| create table t (a int) partition by range(a) subpartition by (partition p values less than (1))
< error
| 1064 42000 …line 1 column 62 near "(partition p values less than (1))" 
> exec ^
: string
| create table t (a int) partition by range(a) subpartition by linear (partition p values less than (1))
< error
| 1064 42000 …line 1 column 69 near "(partition p values less than (1))" 
> exec ^
: string
| create table t (a int, b int) partition by range(a) subpartition by hash(b) (partition x values less than (1) (subpartition y), partition z values less than (2))
< error
| 1485 HY000 Wrong number of subpartitions defined, mismatch with previous setting
> exec ^
: string
| create table t (a int, b int) partition by range(a) subpartition by hash(b) (partition x values less than (1) (subpartition y), partition z values less than (2))
| ; SELECT FROM
< error
| 1485 HY000 Wrong number of subpartitions defined, mismatch with previous setting
> exec ^
: string
| create table t (a int, b int) partition by range(a) subpartition by hash(b) (partition x values less than (1) (subpartition y), partition z values less than (2) (subpartition w))
< ok 0 0
> exec ^
: string
| create table t (a int, b int) partition by list columns(a,b) (partition x values in (default, (1, 2)))
< ok 0 0
> exec ^
: string
| create table t (a int, b int) partition by list columns(a,b) (partition x values in (default, (1, 2), (3, 4)))
< ok 0 0
> exec ^
: string
| create table t (a int, b int) partition by list columns(a,b) (partition x values in (default, (1, 2), (3, 4, 5)))
< error
| 1653 HY000 Inconsistency in usage of column lists for partitioning
> exec ^
: string
| create table t (a int, b int) partition by list columns(a,b) (partition x values in (default, (1, 2), (3, 4, 5)))
| ; SELECT FROM
< error
| 1653 HY000 Inconsistency in usage of column lists for partitioning
> exec ^
: string
| create table t (a int, b int) partition by list columns(a,b) (partition x values in ((1, 2), default))
< ok 0 0
> exec ^
: string
| create table t (a int, b int) partition by list columns(a,b) (partition x values in (default, 1))
< error
| 1653 HY000 Inconsistency in usage of column lists for partitioning
> exec ^
: string
| create table t (a int, b int) partition by list columns(a,b) (partition x values in (default, 1))
| ; SELECT FROM
< error
| 1653 HY000 Inconsistency in usage of column lists for partitioning
> exec ^
: string
| backup database * to 's3://a' snapshot = 1 day ago
< ok 0 0
> exec ^
: string
| backup database * to 's3://a' snapshot = 1 month ago
< error
| 1064 42000 …MONTH is not a constant time interval and cannot be used here
> exec ^
: string
| backup database * to 's3://a' snapshot = 1 year ago
< error
| 1064 42000 …YEAR is not a constant time interval and cannot be used here
> exec ^
: string
| backup database * to 's3://a' snapshot = 1 quarter ago
< error
| 1064 42000 …QUARTER is not a constant time interval and cannot be used here
> exec ^
: string
| backup database * to 's3://a' snapshot = 1 week ago
< ok 0 0
> exec ^
: string
| backup database * to 's3://a' snapshot = 1 microsecond ago
< ok 0 0
> exec ^
: string
| select 1 like 'a' escape 'xy'
< error
| 1210 HY000 Incorrect arguments to ESCAPE
> exec ^
: string
| select 1 like 'a' escape 'xy'
| ; SELECT FROM
< error
| 1210 HY000 Incorrect arguments to ESCAPE
> exec ^
: string
| select 1 like 'a' escape '||'
< error
| 1210 HY000 Incorrect arguments to ESCAPE
> exec ^
: string
| select 1 like 'a' escape '||'
| ; SELECT FROM
< error
| 1210 HY000 Incorrect arguments to ESCAPE
> exec ^
: string
| select 1 like 'a' escape ''
< ok 0 0
> exec ^
: string
| select 1 like 'a' escape 'x'
< ok 0 0
> exec ^
: string
| select 1 not like 'a' escape 'xy'
< error
| 1210 HY000 Incorrect arguments to ESCAPE
> exec ^
: string
| select 1 not like 'a' escape 'xy'
| ; SELECT FROM
< error
| 1210 HY000 Incorrect arguments to ESCAPE
> exec ^
: string
| select 1 like 'a' escape '\\'
< ok 0 0
> exec ^
: string
| select 1 like 'a' escape '€'
< error
| 1210 HY000 Incorrect arguments to ESCAPE
> exec ^
: string
| select 1 like 'a' escape '€'
| ; SELECT FROM
< error
| 1210 HY000 Incorrect arguments to ESCAPE
> exec ^
: string
| select cast(1 as float(54))
< error
| 1426 42000 Too-big precision 54 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| select cast(1 as float(54))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision 54 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| select cast(1 as float(53))
< ok 0 0
> exec ^
: string
| select cast(1 as float(25))
< ok 0 0
> exec ^
: string
| select cast(1 as float(255))
< error
| 1426 42000 Too-big precision 255 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| select cast(1 as float(255))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision 255 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| select cast(1 as float)
< ok 0 0
> exec ^
: string
| select convert(1, float(54))
< error
| 1426 42000 Too-big precision 54 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| select convert(1, float(54))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision 54 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| select convert(1, float(53))
< ok 0 0
> exec ^
: string
| select cast(1 as float(54,2))
< error
| 1426 42000 Too-big precision 54 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| select cast(1 as float(54,2))
| ; SELECT FROM
< error
| 1426 42000 Too-big precision 54 specified for 'CAST'. Maximum is 53.
> exec ^
: string
| select cast(1 as float(10,2))
< ok 0 0
> exec ^
: string
| select cast(1 as double(10))
< error
| 1064 42000 …line 1 column 24 near "(10))" 
> exec ^
: string
| select cast(1 as year)
< ok 0 0
> exec ^
: string
| select cast(1 as char unicode)
< error
| 1064 42000 …line 1 column 29 near "unicode)" 
> exec ^
: string
| select cast(1 as char charset foo)
< error
| 1115 42000 Unknown character set: 'foo'
> exec ^
: string
| select * from ``.t
< error
| 1102 42000 Incorrect database name ''
> exec ^
: string
| select * from ``.t
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name ''
> exec ^
: string
| select * from `a `.t
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| select * from `a `.t
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| select `a `.t.b from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select `a `.t.b from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table `a `.t (x int)
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| create table `a `.t (x int)
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| drop table `a `.t
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| drop table `a `.t
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| create index i on `a `.t (x)
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| create index i on `a `.t (x)
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| create database `a `
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| create database `a `
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table `a `.t add column y int
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| alter table `a `.t add column y int
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| create view `a `.v as select 1
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| create view `a `.v as select 1
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| select * from ``
< error
| 1146 42S02 Table 'test.' doesn't exist
> exec ^
: string
| select * from ``
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select 1 from `a `
< error
| 1146 42S02 Table 'test.a ' doesn't exist
> exec ^
: string
| select 1 from `a `
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select a from ``.`b`
< error
| 1102 42000 Incorrect database name ''
> exec ^
: string
| select a from ``.`b`
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name ''
> exec ^
: string
| select * from `` t
< error
| 1146 42S02 Table 'test.' doesn't exist
> exec ^
: string
| select * from `` t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select `a `.b from t
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| select `a `.b from t
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| update `a `.t set a=1
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| update `a `.t set a=1
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| delete from `a `.t
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| delete from `a `.t
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| insert into `a `.t values (1)
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| insert into `a `.t values (1)
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| load data infile '/x' into table `a `.t
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| load data infile '/x' into table `a `.t
| ; SELECT FROM
< error
| 1102 42000 Incorrect database name 'a '
> exec ^
: string
| select all distinct 1
< error
| 1221 HY000 Incorrect usage of ALL and DISTINCT
> exec ^
: string
| select all distinct 1
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of ALL and DISTINCT
> exec ^
: string
| select distinct all 1
< error
| 1221 HY000 Incorrect usage of ALL and DISTINCT
> exec ^
: string
| select distinct all 1
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of ALL and DISTINCT
> exec ^
: string
| select distinctrow all 1
< error
| 1221 HY000 Incorrect usage of ALL and DISTINCT
> exec ^
: string
| select distinctrow all 1
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of ALL and DISTINCT
> exec ^
: string
| select all distinctrow 1
< error
| 1221 HY000 Incorrect usage of ALL and DISTINCT
> exec ^
: string
| select all distinctrow 1
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of ALL and DISTINCT
> exec ^
: string
| select all all 1
< ok 0 0
> exec ^
: string
| select distinct distinct 1
< ok 0 0
> exec ^
: string
| select distinct distinctrow 1
< ok 0 0
> exec ^
: string
| select all 1
< ok 0 0
> exec ^
: string
| select all high_priority distinct 1
< error
| 1221 HY000 Incorrect usage of ALL and DISTINCT
> exec ^
: string
| select all high_priority distinct 1
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of ALL and DISTINCT
> exec ^
: string
| select distinct high_priority all 1
< error
| 1221 HY000 Incorrect usage of ALL and DISTINCT
> exec ^
: string
| select distinct high_priority all 1
| ; SELECT FROM
< error
| 1221 HY000 Incorrect usage of ALL and DISTINCT
> exec ^
: string
| create table t (a int) stats_auto_recalc=2
< error
| 1064 42000 …line 1 column 42 near ""The value of STATS_AUTO_RECALC must be one of [0|1|DEFAULT]. 
> exec ^
: string
| create table t (a int) stats_auto_recalc=0
< ok 0 0
> exec ^
: string
| create table t (a int) stats_auto_recalc=1
< ok 0 0
> exec ^
: string
| create table t (a int) stats_auto_recalc=default
< ok 0 0
> exec ^
: string
| create table t (a int) ttl=a+interval 1 day ttl_enable='foo'
< error
| 1064 42000 …line 1 column 60 near ""The TTL_ENABLE option has to be set 'ON' or 'OFF' 
> exec ^
: string
| create table t (a int) ttl=a+interval 1 day ttl_enable='ON'
< ok 0 0
> exec ^
: string
| create table t (a int) ttl=a+interval 1 day ttl_enable='off'
< ok 0 0
> exec ^
: string
| create table t (a int) ttl=a+interval 1 day ttl_job_interval='foo'
< error
| 1064 42000 …line 1 column 66 near ""The TTL_JOB_INTERVAL option is not a valid duration: fail to read an integer 
> exec ^
: string
| create table t (a int) ttl=a+interval 1 day ttl_job_interval='1h'
< ok 0 0
> exec ^
: string
| create table t (a int) ttl=a+interval 1 day ttl_job_interval='1d'
< ok 0 0
> exec ^
: string
| create table t (a int) ttl=a+interval 1 day ttl_job_interval=1
< error
| 1064 42000 …line 1 column 62 near "1" 
> exec ^
: string
| create table t (a int) ttl=a+interval 1 day ttl_job_interval=foo
< error
| 1064 42000 …line 1 column 64 near "foo" 
> exec ^
: string
| create table t (a double(10))
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create table t (a double(10))
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create table t (a double(10,2))
< ok 0 0
> exec ^
: string
| create table t (a float(10))
< ok 0 0
> exec ^
: string
| create table t (a real(10))
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create table t (a real(10))
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create table t (a year(2))
< error
| 1818 HY000 Supports only YEAR or YEAR(4) column
> exec ^
: string
| create table t (a year(2))
| ; SELECT FROM
< error
| 1818 HY000 Supports only YEAR or YEAR(4) column
> exec ^
: string
| create table t (a year(4))
< ok 0 0
> exec ^
: string
| create table t (a year(0))
< error
| 1818 HY000 Supports only YEAR or YEAR(4) column
> exec ^
: string
| create table t (a year(0))
| ; SELECT FROM
< error
| 1818 HY000 Supports only YEAR or YEAR(4) column
> exec ^
: string
| create table t (a year(5))
< error
| 1818 HY000 Supports only YEAR or YEAR(4) column
> exec ^
: string
| create table t (a year(5))
| ; SELECT FROM
< error
| 1818 HY000 Supports only YEAR or YEAR(4) column
> exec ^
: string
| create table t (a year)
< ok 0 0
> exec ^
: string
| alter table t add column a year(2)
< error
| 1818 HY000 Supports only YEAR or YEAR(4) column
> exec ^
: string
| alter table t add column a year(2)
| ; SELECT FROM
< error
| 1818 HY000 Supports only YEAR or YEAR(4) column
> exec ^
: string
| create table t (a float8(10))
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create table t (a float8(10))
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create table t (a float8(10,2))
< ok 0 0
> exec ^
: string
| create table t (a float4(10))
< ok 0 0
> exec ^
: string
| create table t (a double precision(10))
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create table t (a double precision(10))
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create table t (a double precision(10,2))
< ok 0 0
> exec ^
: string
| create table t (a double(10,0))
< ok 0 0
> exec ^
: string
| create table t (a real(10,0))
< ok 0 0
> exec ^
: string
| create function f() returns double(10) return 1
< error
| 1064 42000 …line 1 column 15 near "function f() returns double(10) return 1" 
> exec ^
: string
| create function f(a year(2)) returns int return 1
< error
| 1064 42000 …line 1 column 15 near "function f(a year(2)) returns int return 1" 
> exec ^
: string
| grant r1, select on *.* to u
< ok 0 0
> exec ^
: string
| grant select, r1 on *.* to u
< ok 0 0
> exec ^
: string
| grant r1@h on *.* to u
< error
| 1064 42000 …can't convert to PrivElem, type *auth.RoleIdentity
> exec ^
: string
| grant select, r1@h on *.* to u
< error
| 1064 42000 …can't convert to PrivElem, type *auth.RoleIdentity
> exec ^
: string
| grant r1@h, select on *.* to u
< error
| 1064 42000 …can't convert to PrivElem, type *auth.RoleIdentity
> exec ^
: string
| grant select to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant r1, select to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant select, r1 to u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| grant 'r1'@'h' to u
< ok 0 0
> exec ^
: string
| grant r1 on *.* to u
< ok 0 0
> exec ^
: string
| revoke r1@h on *.* from u
< error
| 1064 42000 …can't convert to PrivElem, type *auth.RoleIdentity
> exec ^
: string
| revoke select, r1@h on *.* from u
< error
| 1064 42000 …can't convert to PrivElem, type *auth.RoleIdentity
> exec ^
: string
| revoke select from u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| revoke r1, select from u
< error
| 1064 42000 …can't convert to RoleIdentity, type *ast.PrivElem
> exec ^
: string
| import into t (@a) from select 1
< error
| 1064 42000 …line 1 column 32 near ""Cannot use user variable(a) in IMPORT INTO FROM SELECT statement. 
> exec ^
: string
| import into t (a) from select 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t (a) from select 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t (a) set a=1 from select 1
< error
| 1064 42000 …line 1 column 39 near ""Cannot use SET clause in IMPORT INTO FROM SELECT statement. 
> exec ^
: string
| import into t from select 1
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t from select 1
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| import into t from '/x' with fields_enclosed_by='aa'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| import into t from '/x' with fields_enclosed_by='aa'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a double(10 /* , */))
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create table t (a double(10 /* , */))
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create table t (a double(10 /* , */ , 2))
< ok 0 0
> exec ^
: string
| create sequence s start with -9223372036854775809
< error
| 1064 42000 …line 1 column 49 near ""the Signed Value should be at the range of [-9223372036854775808, 9223372036854775807]. 
> exec ^
: string
| create sequence s start with -9223372036854775808
< ok 0 0
> exec ^
: string
| create sequence s start with 9223372036854775808
< error
| 1064 42000 …line 1 column 48 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| create sequence s start with 9223372036854775807
< ok 0 0
> exec ^
: string
| create sequence s increment by -9223372036854775809
< error
| 1064 42000 …line 1 column 51 near ""the Signed Value should be at the range of [-9223372036854775808, 9223372036854775807]. 
> exec ^
: string
| create sequence s maxvalue 9223372036854775808
< error
| 1064 42000 …line 1 column 46 near ""9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| create sequence s minvalue -9223372036854775809
< error
| 1064 42000 …line 1 column 47 near ""the Signed Value should be at the range of [-9223372036854775808, 9223372036854775807]. 
> exec ^
: string
| select setval(s, -9223372036854775809)
< error
| 1064 42000 …line 1 column 38 near ")"the Signed Value should be at the range of [-9223372036854775808, 9223372036854775807]. 
> exec ^
: string
| select setval(s, 9223372036854775808)
< error
| 1064 42000 …line 1 column 37 near ")"9223372036854775808 is out of range [–9223372036854775808,9223372036854775807] 
> exec ^
: string
| select setval(s, 1+1)
< error
| 1064 42000 …line 1 column 19 near "+1)" 
> exec ^
: string
| select setval(s)
< error
| 1064 42000 …line 1 column 16 near ")" 
> exec ^
: string
| select setval(s, 1)
< error
| 1146 42S02 Table 'test.s' doesn't exist
> exec ^
: string
| select setval(s, 1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select setval(s, -1)
< error
| 1146 42S02 Table 'test.s' doesn't exist
> exec ^
: string
| select setval(s, -1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select setval(s, +1)
< error
| 1146 42S02 Table 'test.s' doesn't exist
> exec ^
: string
| select setval(s, +1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select setval(1, 1)
< error
| 1064 42000 …line 1 column 15 near "1, 1)" 
> exec ^
: string
| select setval(s.t, 1)
< error
| 1146 42S02 Table 's.t' doesn't exist
> exec ^
: string
| select setval(s.t, 1)
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create database d encryption='x'
< error
| 1525 HY000 Incorrect argument (should be Y or N) value: 'x'
> exec ^
: string
| create database d encryption='x'
| ; SELECT FROM
< error
| 1525 HY000 Incorrect argument (should be Y or N) value: 'x'
> exec ^
: string
| create database d encryption='Y'
< ok 0 0
> exec ^
: string
| create database d encryption='n'
< ok 0 0
> exec ^
: string
| alter database d encryption='x'
< error
| 1525 HY000 Incorrect argument (should be Y or N) value: 'x'
> exec ^
: string
| alter database d encryption='x'
| ; SELECT FROM
< error
| 1525 HY000 Incorrect argument (should be Y or N) value: 'x'
> exec ^
: string
| create schema d encryption='yes'
< error
| 1525 HY000 Incorrect argument (should be Y or N) value: 'yes'
> exec ^
: string
| create schema d encryption='yes'
| ; SELECT FROM
< error
| 1525 HY000 Incorrect argument (should be Y or N) value: 'yes'
> exec ^
: string
| create table t (a int) encryption='x'
< error
| 1525 HY000 Incorrect argument (should be Y or N) value: 'x'
> exec ^
: string
| create table t (a int) encryption='x'
| ; SELECT FROM
< error
| 1525 HY000 Incorrect argument (should be Y or N) value: 'x'
> exec ^
: string
| create table t (a int) encryption='Y'
< ok 0 0
> exec ^
: string
| create table t (a int) encryption=Y
< error
| 1064 42000 …line 1 column 35 near "Y" 
> exec ^
: string
| create table t (a int) encryption=1
< error
| 1064 42000 …line 1 column 35 near "1" 
> exec ^
: string
| create table t (a int) encryption=''
< error
| 1525 HY000 Incorrect argument (should be Y or N) value: ''
> exec ^
: string
| create table t (a int) encryption=''
| ; SELECT FROM
< error
| 1525 HY000 Incorrect argument (should be Y or N) value: ''
> exec ^
: string
| create database d encryption=''
< error
| 1525 HY000 Incorrect argument (should be Y or N) value: ''
> exec ^
: string
| create database d encryption=''
| ; SELECT FROM
< error
| 1525 HY000 Incorrect argument (should be Y or N) value: ''
> exec ^
: string
| select _foo'x'
< error
| 1054 42S22 Unknown column '_foo' in 'field list'
> exec ^
: string
| select _foo'x'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create table t (a char(1) unicode)
< error
| 1064 42000 …line 1 column 33 near "unicode)" 
> exec ^
: string
| create table t (a char(1) charset ucs2)
< error
| 1115 42000 Unknown character set: 'ucs2'
> exec ^
: string
| load data infile '/x' into table t fields enclosed by 'aa'
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| load data infile '/x' into table t fields enclosed by 'aa'
| ; SELECT FROM
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| load data infile '/x' into table t fields optionally enclosed by 'aa'
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| load data infile '/x' into table t fields optionally enclosed by 'aa'
| ; SELECT FROM
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| load data infile '/x' into table t fields escaped by 'aa'
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| load data infile '/x' into table t fields escaped by 'aa'
| ; SELECT FROM
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| load data infile '/x' into table t fields enclosed by X'0D0D'
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| load data infile '/x' into table t fields enclosed by X'0D0D'
| ; SELECT FROM
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| load data infile '/x' into table t fields enclosed by B'110101101101011'
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| load data infile '/x' into table t fields enclosed by B'110101101101011'
| ; SELECT FROM
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| load data infile '/x' into table t fields escaped by X'6B6B'
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| load data infile '/x' into table t fields escaped by X'6B6B'
| ; SELECT FROM
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| load data infile '/x' into table t fields escaped by B'110101101101011'
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| load data infile '/x' into table t fields escaped by B'110101101101011'
| ; SELECT FROM
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| load data infile '/x' into table t fields enclosed by '\\'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x' into table t fields enclosed by '\\'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x' into table t fields enclosed by ''
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x' into table t fields enclosed by ''
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x' into table t fields enclosed by '€'
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| load data infile '/x' into table t fields enclosed by '€'
| ; SELECT FROM
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| load data infile '/x' into table t fields escaped by '€'
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| load data infile '/x' into table t fields escaped by '€'
| ; SELECT FROM
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| load data infile '/x' into table t fields escaped by 'é'
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| load data infile '/x' into table t fields escaped by 'é'
| ; SELECT FROM
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| load data infile '/x' into table t columns enclosed by 'aa'
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| load data infile '/x' into table t columns enclosed by 'aa'
| ; SELECT FROM
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| load data infile '/x' into table t fields terminated by 'aa'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x' into table t fields terminated by 'aa'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| load data infile '/x' into table t lines terminated by 'aa'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| load data infile '/x' into table t lines terminated by 'aa'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select 1 into outfile '/x' fields enclosed by 'aa'
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| select 1 into outfile '/x' fields enclosed by 'aa'
| ; SELECT FROM
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| select 1 into outfile '/x' fields escaped by 'aa'
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| select 1 into outfile '/x' fields escaped by 'aa'
| ; SELECT FROM
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| select 1 into outfile '/x' fields optionally enclosed by 'aa'
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| select 1 into outfile '/x' fields optionally enclosed by 'aa'
| ; SELECT FROM
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| select 1 into outfile '/x' fields enclosed by 'a'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| select 1 into outfile '/x' fields enclosed by 'a'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select 1 into outfile '/x' fields enclosed by X'0D0D'
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| select 1 into outfile '/x' fields enclosed by X'0D0D'
| ; SELECT FROM
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| select 1 into outfile '/x' fields enclosed by X'0D'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| select 1 into outfile '/x' fields enclosed by X'0D'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select 1 into outfile '/x' fields terminated by X'6B6B'
< error
| 1295 HY000 This command is not supported in the prepared statement protocol yet
> exec ^
: string
| select 1 into outfile '/x' fields terminated by X'6B6B'
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| select 1 into outfile '/x' fields escaped by B'110101101101011'
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| select 1 into outfile '/x' fields escaped by B'110101101101011'
| ; SELECT FROM
< error
| 1083 42000 Field separator argument is not what is expected; check the manual
> exec ^
: string
| select 1 into outfile '/x' fields enclosed 'a'
< error
| 1064 42000 …line 1 column 46 near "'a'" 
> exec ^
: string
| select 1 into outfile '/x' fields enclosed by
< error
| 1064 42000 …line 1 column 45 near "" 
> exec ^
: string
| create index i using hnsw on t (a)
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create index i using hnsw on t (a)
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create index i using hnsw on t (a) using btree
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| create index i using hnsw on t (a) using btree
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| create index i on t (a) using hnsw
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create index i on t (a) using hnsw
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create index i on t (a) type hnsw
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create index i on t (a) type hnsw
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create unique index i on t (a) using hnsw
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create unique index i on t (a) using hnsw
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create fulltext index i on t (a) using hnsw
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create fulltext index i on t (a) using hnsw
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create spatial index i on t (a) using hnsw
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create spatial index i on t (a) using hnsw
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create index i using btree on t (a) using hnsw
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create index i using btree on t (a) using hnsw
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create table t (a int, key i (a) using hnsw)
< ok 0 0
> exec ^
: string
| alter table t add key i (a) using hnsw
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add key i (a) using hnsw
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add vector index (a)
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| alter table t add vector index (a)
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| alter table t add vector index ((a))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add vector index ((a))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add vector index ((a), (b))
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| alter table t add vector index ((a), (b))
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| alter table t add vector index i using hnsw ((a))
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add vector index i using hnsw ((a))
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add vector index i using btree ((a))
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| alter table t add vector index i using btree ((a))
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| alter table t add vector index ((a)) using hash
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| alter table t add vector index ((a)) using hash
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| alter table t add vector index ((a)) using hypo
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| alter table t add vector index ((a)) using hypo
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| alter table t add vector index i using hnsw ((a)) using btree
< error
| 1146 42S02 Table 'test.t' doesn't exist
> exec ^
: string
| alter table t add vector index i using hnsw ((a)) using btree
| ; SELECT FROM
< error
| 1064 42000 …line 2 column 14 near "FROM" 
> exec ^
: string
| alter table t add vector index i using btree ((a)) using hnsw
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| alter table t add vector index i using btree ((a)) using hnsw
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| alter table t add column c int, add (vector index i (a))
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| alter table t add column c int, add (vector index i (a))
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create table t (a int, b int, vector index (b))
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create table t (a int, b int, vector index (b))
| ; SELECT FROM
< error
| 1149 42000 You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
> exec ^
: string
| create table t (a int, b int, vector index ((b)))
< ok 0 0
