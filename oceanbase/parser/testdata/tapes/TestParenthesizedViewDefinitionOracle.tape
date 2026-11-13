-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

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
| select `omni_t_oceanbase_parser`.`b`.`id` AS `id`,`omni_t_oceanbase_parser`.`b`.`k` AS `k` from `omni_t_oceanbase_parser`.`b`
> query
| SHOW CREATE VIEW vb
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| vb
: bytes
| CREATE VIEW `vb` AS select `omni_t_oceanbase_parser`.`b`.`id` AS `id`,`omni_t_oceanbase_parser`.`b`.`k` AS `k` from `omni_t_oceanbase_parser`.`b`
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
| select `t1`.`id` AS `id`,`omni_t_oceanbase_parser`.`t2`.`k` AS `k` from (`omni_t_oceanbase_parser`.`a` `t1` join `omni_t_oceanbase_parser`.`vb` `t2` on (`t1`.`id` = `omni_t_oceanbase_parser`.`t2`.`id`))
> query
| SHOW CREATE VIEW r2
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| r2
: bytes
| CREATE VIEW `r2` AS select `t1`.`id` AS `id`,`omni_t_oceanbase_parser`.`t2`.`k` AS `k` from (`omni_t_oceanbase_parser`.`a` `t1` join `omni_t_oceanbase_parser`.`vb` `t2` on (`t1`.`id` = `omni_t_oceanbase_parser`.`t2`.`id`))
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
| select `t1`.`id` AS `id`,`omni_t_oceanbase_parser`.`t2`.`k` AS `k` from (`omni_t_oceanbase_parser`.`a` `t1` join `omni_t_oceanbase_parser`.`vb` `t2` on (`t1`.`id` = `omni_t_oceanbase_parser`.`t2`.`id`) and (`t1`.`k` = `omni_t_oceanbase_parser`.`t2`.`k`) and (`omni_t_oceanbase_parser`.`t2`.`k` <> 'x') and (`omni_t_oceanbase_parser`.`t2`.`k` <> 'y'))
> query
| SHOW CREATE VIEW r4
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| r4
: bytes
| CREATE VIEW `r4` AS select `t1`.`id` AS `id`,`omni_t_oceanbase_parser`.`t2`.`k` AS `k` from (`omni_t_oceanbase_parser`.`a` `t1` join `omni_t_oceanbase_parser`.`vb` `t2` on (`t1`.`id` = `omni_t_oceanbase_parser`.`t2`.`id`) and (`t1`.`k` = `omni_t_oceanbase_parser`.`t2`.`k`) and (`omni_t_oceanbase_parser`.`t2`.`k` <> 'x') and (`omni_t_oceanbase_parser`.`t2`.`k` <> 'y'))
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
| (select `omni_t_oceanbase_parser`.`r2`.`id` AS `id` from `omni_t_oceanbase_parser`.`r2`) union  all (select `omni_t_oceanbase_parser`.`r4`.`id` AS `id` from `omni_t_oceanbase_parser`.`r4`)
> query
| SHOW CREATE VIEW r5
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| r5
: bytes
| CREATE VIEW `r5` AS (select `omni_t_oceanbase_parser`.`r2`.`id` AS `id` from `omni_t_oceanbase_parser`.`r2`) union  all (select `omni_t_oceanbase_parser`.`r4`.`id` AS `id` from `omni_t_oceanbase_parser`.`r4`)
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
| select `t1`.`id` AS `id` from (`omni_t_oceanbase_parser`.`a` `t1` join `omni_t_oceanbase_parser`.`b` `t2` on (`t1`.`id` = `t2`.`id`) and (`t1`.`k` = `t2`.`k`))
> query
| SHOW CREATE VIEW r8
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| r8
: bytes
| CREATE VIEW `r8` AS select `t1`.`id` AS `id` from (`omni_t_oceanbase_parser`.`a` `t1` join `omni_t_oceanbase_parser`.`b` `t2` on (`t1`.`id` = `t2`.`id`) and (`t1`.`k` = `t2`.`k`))
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
| select `t1`.`id` AS `id` from (`omni_t_oceanbase_parser`.`a` `t1` join `omni_t_oceanbase_parser`.`b` `t2` on (`t1`.`id` = `t2`.`id`) and (`t2`.`k` <> 'x'))
> query
| SHOW CREATE VIEW r9
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| r9
: bytes
| CREATE VIEW `r9` AS select `t1`.`id` AS `id` from (`omni_t_oceanbase_parser`.`a` `t1` join `omni_t_oceanbase_parser`.`b` `t2` on (`t1`.`id` = `t2`.`id`) and (`t2`.`k` <> 'x'))
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
| select `t1`.`id` AS `id` from (`omni_t_oceanbase_parser`.`a` `t1` join `omni_t_oceanbase_parser`.`b` `t2` on (`t1`.`id` = `t2`.`id`)) where (`t2`.`k` <> 'x')
> query
| SHOW CREATE VIEW r10
< result
< column "View" "VARCHAR"
< column "Create View" "TEXT"
< column "character_set_client" "VARCHAR"
< column "collation_connection" "VARCHAR"
< row
: bytes
| r10
: bytes
| CREATE VIEW `r10` AS select `t1`.`id` AS `id` from (`omni_t_oceanbase_parser`.`a` `t1` join `omni_t_oceanbase_parser`.`b` `t2` on (`t1`.`id` = `t2`.`id`)) where (`t2`.`k` <> 'x')
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
