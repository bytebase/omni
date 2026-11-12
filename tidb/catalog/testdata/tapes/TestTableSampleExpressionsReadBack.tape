-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestTableSampleExpressionsReadBack default#1
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(@seed) REPEATABLE (@seed)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (@`seed`) repeatable (@`seed`)'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(@`a b`) REPEATABLE (@`a b`)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (@`a b`) repeatable (@`a b`)'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(@@autocommit) REPEATABLE (@@autocommit)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (@@autocommit) repeatable (@@autocommit)'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(@@global.autocommit) REPEATABLE (@@global.autocommit)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (@@global.autocommit) repeatable (@@global.autocommit)'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(@@session.autocommit) REPEATABLE (@@session.autocommit)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (@@session.autocommit) repeatable (@@session.autocommit)'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(@@local.autocommit) REPEATABLE (@@local.autocommit)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (@@local.autocommit) repeatable (@@local.autocommit)'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(@a := 1) REPEATABLE (@a := 1)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions ((@`a` := 1)) repeatable ((@`a` := 1))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(@a := @b := 1) REPEATABLE (@a := @b := 1)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions ((@`a` := (@`b` := 1))) repeatable ((@`a` := (@`b` := 1)))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(match(a) against(''x'')) REPEATABLE (match(a) against(''x''))'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (match (`a`) against (''x'')) repeatable (match (`a`) against (''x''))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(match(a, b) against(''x'' in boolean mode)) REPEATABLE (match(a, b) against(''x'' in boolean mode))'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (match (`a`,`b`) against (''x'' in boolean mode)) repeatable (match (`a`,`b`) against (''x'' in boolean mode))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(match(a) against(''x'' in natural language mode)) REPEATABLE (match(a) against(''x'' in natural language mode))'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (match (`a`) against (''x'' in natural language mode)) repeatable (match (`a`) against (''x'' in natural language mode))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(match(a) against(''x'' with query expansion)) REPEATABLE (match(a) against(''x'' with query expansion))'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (match (`a`) against (''x'' with query expansion)) repeatable (match (`a`) against (''x'' with query expansion))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(match(a) against(''x'' in natural language mode with query expansion)) REPEATABLE (match(a) against(''x'' in natural language mode with query expansion))'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (match (`a`) against (''x'' in natural language mode with query expansion)) repeatable (match (`a`) against (''x'' in natural language mode with query expansion))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(1 member of (''[1]'')) REPEATABLE (1 member of (''[1]''))'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions ((1 member of (''[1]''))) repeatable ((1 member of (''[1]'')))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(a member of (''[1]'')) REPEATABLE (a member of (''[1]''))'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions ((`a` member of (''[1]''))) repeatable ((`a` member of (''[1]'')))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(default(a)) REPEATABLE (default(a))'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (default(`a`)) repeatable (default(`a`))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(values(a)) REPEATABLE (values(a))'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (values(`a`)) repeatable (values(`a`))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(a->''$.x'') REPEATABLE (a->''$.x'')'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (json_extract(`a`,''$.x'')) repeatable (json_extract(`a`,''$.x''))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(a->>''$.x'') REPEATABLE (a->>''$.x'')'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (json_unquote(json_extract(`a`,''$.x''))) repeatable (json_unquote(json_extract(`a`,''$.x'')))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(count(*)) REPEATABLE (count(*))'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (count(0)) repeatable (count(0))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(sum(a) over ()) REPEATABLE (sum(a) over ())'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (sum(`a`) OVER ( )) repeatable (sum(`a`) OVER ( ))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(sum(a) over (partition by b order by a rows between 1 preceding and current row)) REPEATABLE (sum(a) over (partition by b order by a rows between 1 preceding and current row))'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (sum(`a`) OVER (PARTITION BY `b` ORDER BY `a` ROWS BETWEEN 1 PRECEDING AND CURRENT ROW)) repeatable (sum(`a`) OVER (PARTITION BY `b` ORDER BY `a` ROWS BETWEEN 1 PRECEDING AND CURRENT ROW))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(row_number() over w) REPEATABLE (row_number() over w)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (row_number() OVER `w`) repeatable (row_number() OVER `w`)'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(group_concat(a order by a separator '','')) REPEATABLE (group_concat(a order by a separator '',''))'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (group_concat(`a` order by `a` ASC separator '','')) repeatable (group_concat(`a` order by `a` ASC separator '',''))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(_utf8mb4''x'') REPEATABLE (_utf8mb4''x'')'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (_utf8mb4''x'') repeatable (_utf8mb4''x'')'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(x''41'') REPEATABLE (x''41'')'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (0x41) repeatable (0x41)'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(b''1'') REPEATABLE (b''1'')'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (0x01) repeatable (0x01)'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(date ''2020-01-01'') REPEATABLE (date ''2020-01-01'')'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (DATE''2020-01-01'') repeatable (DATE''2020-01-01'')'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(true) REPEATABLE (true)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (true) repeatable (true)'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(null) REPEATABLE (null)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (NULL) repeatable (NULL)'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(1e3) REPEATABLE (1e3)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (1e3) repeatable (1e3)'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(interval 1 day + a) REPEATABLE (interval 1 day + a)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions ((`a` + interval 1 day)) repeatable ((`a` + interval 1 day))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(now()) REPEATABLE (now())'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (now()) repeatable (now())'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(current_timestamp) REPEATABLE (current_timestamp)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (now()) repeatable (now())'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(current_user) REPEATABLE (current_user)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (current_user()) repeatable (current_user())'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(not a) REPEATABLE (not a)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions ((not(`a`))) repeatable ((not(`a`)))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(!a) REPEATABLE (!a)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions ((not(`a`))) repeatable ((not(`a`)))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(-a) REPEATABLE (-a)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (-(`a`)) repeatable (-(`a`))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(~a) REPEATABLE (~a)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (~(`a`)) repeatable (~(`a`))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(a || b) REPEATABLE (a || b)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions ((`a` or `b`)) repeatable ((`a` or `b`))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(a && b) REPEATABLE (a && b)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions ((`a` and `b`)) repeatable ((`a` and `b`))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(a xor b) REPEATABLE (a xor b)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions ((`a` xor `b`)) repeatable ((`a` xor `b`))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(a div b) REPEATABLE (a div b)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions ((`a` DIV `b`)) repeatable ((`a` DIV `b`))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(a mod b) REPEATABLE (a mod b)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions ((`a` % `b`)) repeatable ((`a` % `b`))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(a <=> b) REPEATABLE (a <=> b)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions ((`a` <=> `b`)) repeatable ((`a` <=> `b`))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(a in (1,2)) REPEATABLE (a in (1,2))'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions ((`a` in (1,2))) repeatable ((`a` in (1,2)))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(a in (select 1)) REPEATABLE (a in (select 1))'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (`a` in (select 1)) repeatable (`a` in (select 1))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(a = any (select 1)) REPEATABLE (a = any (select 1))'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions ((`a` = any (select 1))) repeatable ((`a` = any (select 1)))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(a > all (select 1)) REPEATABLE (a > all (select 1))'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions ((`a` > all (select 1))) repeatable ((`a` > all (select 1)))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(exists (select 1)) REPEATABLE (exists (select 1))'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (exists(select 1)) repeatable (exists(select 1))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(row(1,2) = row(1,2)) REPEATABLE (row(1,2) = row(1,2))'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions ((row(1,2) = row(1,2))) repeatable ((row(1,2) = row(1,2)))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(a like ''x'' escape ''y'') REPEATABLE (a like ''x'' escape ''y'')'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions ((`a` like ''x'' escape ''y'')) repeatable ((`a` like ''x'' escape ''y''))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(a regexp ''x'') REPEATABLE (a regexp ''x'')'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (regexp_like(`a`,''x'')) repeatable (regexp_like(`a`,''x''))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(a between 1 and 2) REPEATABLE (a between 1 and 2)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions ((`a` between 1 and 2)) repeatable ((`a` between 1 and 2))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(a is true) REPEATABLE (a is true)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions ((`a` is true)) repeatable ((`a` is true))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(a is not null) REPEATABLE (a is not null)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions ((`a` is not null)) repeatable ((`a` is not null))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(case when a then 1 else 2 end) REPEATABLE (case when a then 1 else 2 end)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions ((case when `a` then 1 else 2 end)) repeatable ((case when `a` then 1 else 2 end))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(case a when 1 then 2 end) REPEATABLE (case a when 1 then 2 end)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions ((case `a` when 1 then 2 end)) repeatable ((case `a` when 1 then 2 end))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(cast(a as char)) REPEATABLE (cast(a as char))'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (cast(`a` as char charset utf8mb4)) repeatable (cast(`a` as char charset utf8mb4))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(convert(a, char)) REPEATABLE (convert(a, char))'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (cast(`a` as char charset utf8mb4)) repeatable (cast(`a` as char charset utf8mb4))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(convert(a using utf8mb4)) REPEATABLE (convert(a using utf8mb4))'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (convert(`a` using utf8mb4)) repeatable (convert(`a` using utf8mb4))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(extract(day from a)) REPEATABLE (extract(day from a))'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (extract(day from `a`)) repeatable (extract(day from `a`))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(position(''a'' in b)) REPEATABLE (position(''a'' in b))'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (position(''a'' in `b`)) repeatable (position(''a'' in `b`))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(substring(a from 1 for 2)) REPEATABLE (substring(a from 1 for 2))'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (substr(`a`,1,2)) repeatable (substr(`a`,1,2))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(trim(both ''x'' from a)) REPEATABLE (trim(both ''x'' from a))'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (trim(both ''x'' from `a`)) repeatable (trim(both ''x'' from `a`))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(a collate utf8mb4_bin) REPEATABLE (a collate utf8mb4_bin)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions ((`a` collate utf8mb4_bin)) repeatable ((`a` collate utf8mb4_bin))'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(binary a) REPEATABLE (binary a)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (`a`) repeatable (`a`)'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(t.a) REPEATABLE (t.a)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (`t`.`a`) repeatable (`t`.`a`)'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'SELECT a FROM tbl TABLESAMPLE REGIONS(d.t.a) REPEATABLE (d.t.a)'
< error
| 1146 42S02 Table 'test.tbl' doesn't exist
> exec
| PREPARE omni_parses FROM 'select `tbl`.`a` from `d`.`tbl` tablesample regions (`d`.`t`.`a`) repeatable (`d`.`t`.`a`)'
< error
| 1146 42S02 Table 'd.tbl' doesn't exist
