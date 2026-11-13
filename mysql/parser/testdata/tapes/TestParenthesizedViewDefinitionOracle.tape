-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestParenthesizedViewDefinitionOracle default#1
> exec
| CREATE TABLE a (id INT, k VARCHAR(20))
< ok 0 0
> exec
| CREATE TABLE b (id INT, k VARCHAR(20))
< ok 0 0
> exec
| CREATE VIEW vb AS SELECT id, k FROM b
< ok 0 0
> exec
| CREATE VIEW r2 AS SELECT t1.id, t2.k FROM a t1 JOIN vb t2 ON t1.id = t2.id
< ok 0 0
> exec
| CREATE VIEW r4 AS SELECT t1.id, t2.k FROM a t1 JOIN vb t2 ON t1.id = t2.id AND t1.k = t2.k AND t2.k <> 'x' AND t2.k <> 'y'
< ok 0 0
> exec
| CREATE VIEW r5 AS SELECT id FROM r2 UNION ALL SELECT id FROM r4
< ok 0 0
> exec
| CREATE VIEW r8 AS SELECT t1.id FROM a t1 JOIN b t2 ON t1.id = t2.id AND t1.k = t2.k
< ok 0 0
> exec
| CREATE VIEW r9 AS SELECT t1.id FROM a t1 JOIN b t2 ON t1.id = t2.id AND t2.k <> 'x'
< ok 0 0
> exec
| CREATE VIEW r10 AS SELECT t1.id FROM a t1 JOIN b t2 ON t1.id = t2.id WHERE t2.k <> 'x'
< ok 0 0
> exec
| SELECT count(*) FROM a WHERE ((a.id = 1) AND (a.k = 'x'))
< ok 0 0
> exec
| SELECT count(*) FROM a WHERE (((a.id = 1) AND (a.k = 'x')))
< ok 0 0
> exec
| SELECT count(*) FROM a t1 JOIN b t2 ON (((t1.id = t2.id) AND (t1.k = t2.k)))
< ok 0 0
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
: string
| vb
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `omni_t_mysql_parser`.`b`.`id` AS `id`,`omni_t_mysql_parser`.`b`.`k` AS `k` from `omni_t_mysql_parser`.`b`
> query
| SHOW CREATE VIEW vb
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| vb
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`omni_test`@`%` SQL SECURITY DEFINER VIEW `vb` AS select `b`.`id` AS `id`,`b`.`k` AS `k` from `b`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
: string
| r2
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `t1`.`id` AS `id`,`omni_t_mysql_parser`.`t2`.`k` AS `k` from (`omni_t_mysql_parser`.`a` `t1` join `omni_t_mysql_parser`.`vb` `t2` on((`t1`.`id` = `omni_t_mysql_parser`.`t2`.`id`)))
> query
| SHOW CREATE VIEW r2
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| r2
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`omni_test`@`%` SQL SECURITY DEFINER VIEW `r2` AS select `t1`.`id` AS `id`,`t2`.`k` AS `k` from (`a` `t1` join `vb` `t2` on((`t1`.`id` = `t2`.`id`)))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
: string
| r4
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `t1`.`id` AS `id`,`omni_t_mysql_parser`.`t2`.`k` AS `k` from (`omni_t_mysql_parser`.`a` `t1` join `omni_t_mysql_parser`.`vb` `t2` on(((`t1`.`id` = `omni_t_mysql_parser`.`t2`.`id`) and (`t1`.`k` = `omni_t_mysql_parser`.`t2`.`k`) and (`omni_t_mysql_parser`.`t2`.`k` <> 'x') and (`omni_t_mysql_parser`.`t2`.`k` <> 'y'))))
> query
| SHOW CREATE VIEW r4
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| r4
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`omni_test`@`%` SQL SECURITY DEFINER VIEW `r4` AS select `t1`.`id` AS `id`,`t2`.`k` AS `k` from (`a` `t1` join `vb` `t2` on(((`t1`.`id` = `t2`.`id`) and (`t1`.`k` = `t2`.`k`) and (`t2`.`k` <> 'x') and (`t2`.`k` <> 'y'))))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
: string
| r5
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `omni_t_mysql_parser`.`r2`.`id` AS `id` from `omni_t_mysql_parser`.`r2` union all select `omni_t_mysql_parser`.`r4`.`id` AS `id` from `omni_t_mysql_parser`.`r4`
> query
| SHOW CREATE VIEW r5
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| r5
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`omni_test`@`%` SQL SECURITY DEFINER VIEW `r5` AS select `r2`.`id` AS `id` from `r2` union all select `r4`.`id` AS `id` from `r4`
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
: string
| r8
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `t1`.`id` AS `id` from (`omni_t_mysql_parser`.`a` `t1` join `omni_t_mysql_parser`.`b` `t2` on(((`t1`.`id` = `t2`.`id`) and (`t1`.`k` = `t2`.`k`))))
> query
| SHOW CREATE VIEW r8
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| r8
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`omni_test`@`%` SQL SECURITY DEFINER VIEW `r8` AS select `t1`.`id` AS `id` from (`a` `t1` join `b` `t2` on(((`t1`.`id` = `t2`.`id`) and (`t1`.`k` = `t2`.`k`))))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
: string
| r9
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `t1`.`id` AS `id` from (`omni_t_mysql_parser`.`a` `t1` join `omni_t_mysql_parser`.`b` `t2` on(((`t1`.`id` = `t2`.`id`) and (`t2`.`k` <> 'x'))))
> query
| SHOW CREATE VIEW r9
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| r9
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`omni_test`@`%` SQL SECURITY DEFINER VIEW `r9` AS select `t1`.`id` AS `id` from (`a` `t1` join `b` `t2` on(((`t1`.`id` = `t2`.`id`) and (`t2`.`k` <> 'x'))))
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
> query
| SELECT VIEW_DEFINITION FROM information_schema.VIEWS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
: string
| r10
< result
< column "VIEW_DEFINITION" "TEXT"
< row
: bytes
| select `t1`.`id` AS `id` from (`omni_t_mysql_parser`.`a` `t1` join `omni_t_mysql_parser`.`b` `t2` on((`t1`.`id` = `t2`.`id`))) where (`t2`.`k` <> 'x')
> query
| SHOW CREATE VIEW r10
< result
< column "View" "VARCHAR"
< column "Create View" "VARCHAR"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| r10
: bytes
| CREATE ALGORITHM=UNDEFINED DEFINER=`omni_test`@`%` SQL SECURITY DEFINER VIEW `r10` AS select `t1`.`id` AS `id` from (`a` `t1` join `b` `t2` on((`t1`.`id` = `t2`.`id`))) where (`t2`.`k` <> 'x')
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
