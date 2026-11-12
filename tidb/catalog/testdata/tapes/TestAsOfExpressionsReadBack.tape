-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestAsOfExpressionsReadBack default#1
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP @seed'
< error
| 8135 HY000 invalid as of timestamp: as of timestamp cannot be NULL
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp @`seed`'
< error
| 8135 HY000 invalid as of timestamp: as of timestamp cannot be NULL
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP @`a b`'
< error
| 8135 HY000 invalid as of timestamp: as of timestamp cannot be NULL
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp @`a b`'
< error
| 8135 HY000 invalid as of timestamp: as of timestamp cannot be NULL
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP @@autocommit'
< error
| 1292 22007 Incorrect time value: '1'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp @@autocommit'
< error
| 1292 22007 Incorrect time value: '1'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP @@global.autocommit'
< error
| 1292 22007 Incorrect time value: '1'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp @@global.autocommit'
< error
| 1292 22007 Incorrect time value: '1'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP @@session.autocommit'
< error
| 1292 22007 Incorrect time value: '1'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp @@session.autocommit'
< error
| 1292 22007 Incorrect time value: '1'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP @@local.autocommit'
< error
| 1292 22007 Incorrect time value: '1'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp @@local.autocommit'
< error
| 1292 22007 Incorrect time value: '1'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP @a := 1'
< error
| 1292 22007 Incorrect time value: '1'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp (@`a` := 1)'
< error
| 1292 22007 Incorrect time value: '1'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP @a := @b := 1'
< error
| 1292 22007 Incorrect time value: '1'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp (@`a` := (@`b` := 1))'
< error
| 1292 22007 Incorrect time value: '1'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP match(a) against(''x'')'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp match (`a`) against (''x'')'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP match(a, b) against(''x'' in boolean mode)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp match (`a`,`b`) against (''x'' in boolean mode)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP match(a) against(''x'' in natural language mode)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp match (`a`) against (''x'' in natural language mode)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP match(a) against(''x'' with query expansion)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp match (`a`) against (''x'' with query expansion)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP match(a) against(''x'' in natural language mode with query expansion)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp match (`a`) against (''x'' in natural language mode with query expansion)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP 1 member of (''[1]'')'
< error
| 1292 22007 Incorrect time value: '1'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp (1 member of (''[1]''))'
< error
| 1292 22007 Incorrect time value: '1'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP a member of (''[1]'')'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp (`a` member of (''[1]''))'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP default(a)'
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp default(`a`)'
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP values(a)'
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp values(`a`)'
< error
| 1054 42S22 Unknown column 'a' in 'field list'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP a->''$.x'''
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp json_extract(`a`,''$.x'')'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP a->>''$.x'''
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp json_unquote(json_extract(`a`,''$.x''))'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP count(*)'
< error
| 1111 HY000 Invalid use of group function
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp count(0)'
< error
| 1111 HY000 Invalid use of group function
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP sum(a) over ()'
< error
| 3593 HY000 You cannot use the window function 'sum' in this context.'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp sum(`a`) OVER ( )'
< error
| 3593 HY000 You cannot use the window function 'sum' in this context.'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP sum(a) over (partition by b order by a rows between 1 preceding and current row)'
< error
| 3593 HY000 You cannot use the window function 'sum' in this context.'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp sum(`a`) OVER (PARTITION BY `b` ORDER BY `a` ROWS BETWEEN 1 PRECEDING AND CURRENT ROW)'
< error
| 3593 HY000 You cannot use the window function 'sum' in this context.'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP row_number() over w'
< error
| 3593 HY000 You cannot use the window function 'row_number' in this context.'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp row_number() OVER `w`'
< error
| 3593 HY000 You cannot use the window function 'row_number' in this context.'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP group_concat(a order by a separator '','')'
< error
| 1111 HY000 Invalid use of group function
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp group_concat(`a` order by `a` ASC separator '','')'
< error
| 1111 HY000 Invalid use of group function
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP _utf8mb4''x'''
< error
| 1292 22007 Incorrect datetime value: 'x'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp _utf8mb4''x'''
< error
| 1292 22007 Incorrect datetime value: 'x'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP x''41'''
< error
| 1105 HY000 cannot convert datum from bit/hex literal to type timestamp
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp 0x41'
< error
| 1105 HY000 cannot convert datum from bit/hex literal to type timestamp
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP b''1'''
< error
| 1105 HY000 cannot convert datum from bit/hex literal to type timestamp
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp 0x01'
< error
| 1105 HY000 cannot convert datum from bit/hex literal to type timestamp
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP date ''2020-01-01'''
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp DATE''2020-01-01'''
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP true'
< error
| 1292 22007 Incorrect time value: '1'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp true'
< error
| 1292 22007 Incorrect time value: '1'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP null'
< error
| 8135 HY000 invalid as of timestamp: as of timestamp cannot be NULL
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp NULL'
< error
| 8135 HY000 invalid as of timestamp: as of timestamp cannot be NULL
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP 1e3'
< error
| 1105 HY000 cannot convert datum from double to type timestamp
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp 1e3'
< error
| 1105 HY000 cannot convert datum from double to type timestamp
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP interval 1 day + a'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp (`a` + interval 1 day)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP now()'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp now()'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP current_timestamp'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp now()'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP current_user'
< error
| 1292 22007 Incorrect datetime value: 'root@%'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp current_user()'
< error
| 1292 22007 Incorrect datetime value: 'root@%'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP not a'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp (not(`a`))'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP !a'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp (not(`a`))'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP -a'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp -(`a`)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP ~a'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp ~(`a`)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP a || b'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp (`a` or `b`)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP a && b'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp (`a` and `b`)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP a xor b'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp (`a` xor `b`)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP a div b'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp (`a` DIV `b`)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP a mod b'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp (`a` % `b`)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP a <=> b'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp (`a` <=> `b`)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP a in (1,2)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp (`a` in (1,2))'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP a in (select 1)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp `a` in (select 1)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP a = any (select 1)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp (`a` = any (select 1))'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP a > all (select 1)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp (`a` > all (select 1))'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP exists (select 1)'
< error
| 1292 22007 Incorrect time value: '1'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp exists(select 1)'
< error
| 1292 22007 Incorrect time value: '1'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP row(1,2) = row(1,2)'
< error
| 1292 22007 Incorrect time value: '1'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp (row(1,2) = row(1,2))'
< error
| 1292 22007 Incorrect time value: '1'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP a like ''x'' escape ''y'''
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp (`a` like ''x'' escape ''y'')'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP a regexp ''x'''
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp regexp_like(`a`,''x'')'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP a between 1 and 2'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp (`a` between 1 and 2)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP a is true'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp (`a` is true)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP a is not null'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp (`a` is not null)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP case when a then 1 else 2 end'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp (case when `a` then 1 else 2 end)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP case a when 1 then 2 end'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp (case `a` when 1 then 2 end)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP cast(a as char)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp cast(`a` as char charset utf8mb4)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP convert(a, char)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp cast(`a` as char charset utf8mb4)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP convert(a using utf8mb4)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp convert(`a` using utf8mb4)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP extract(day from a)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp extract(day from `a`)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP position(''a'' in b)'
< error
| 1054 42S22 Unknown column 'b' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp position(''a'' in `b`)'
< error
| 1054 42S22 Unknown column 'b' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP substring(a from 1 for 2)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp substr(`a`,1,2)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP trim(both ''x'' from a)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp trim(both ''x'' from `a`)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP a collate utf8mb4_bin'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp (`a` collate utf8mb4_bin)'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP binary a'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp `a`'
< error
| 1054 42S22 Unknown column 'a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP t.a'
< error
| 1054 42S22 Unknown column 't.a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp `t`.`a`'
< error
| 1054 42S22 Unknown column 't.a' in 'expression'
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl AS OF TIMESTAMP d.t.a'
< error
| 1054 42S22 Unknown column 'd.t.a' in 'expression'
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` as of timestamp `d`.`t`.`a`'
< error
| 1054 42S22 Unknown column 'd.t.a' in 'expression'
