-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12
-- pin: pg12 postgres:12.22-alpine sha256:7c8f4870583184ebadf7f17a6513620aac5f365a7938dc6a6911c1d5df2f481a REL_12_22
-- pin: pg18 postgres:18.7-alpine sha256:75dcac3dfb1cd01ca3536910bbc84128759362f8b118bf8c55b38a94f68f3b73 REL_18_7

== TestCountAffectedRows/default default#1
> exec
|
| 				CREATE TABLE s (id int PRIMARY KEY, flag bool NOT NULL);
| 				INSERT INTO s SELECT g, g <= 10 FROM generate_series(1, 100) g;
| 				CREATE TABLE big (id int PRIMARY KEY, s_id int NOT NULL, v int);
| 				INSERT INTO big SELECT g, (g % 100) + 1, 0 FROM generate_series(1, 1000) g;
| 				CREATE TABLE big_archive (LIKE big);
| 				ANALYZE s;
| 				ANALYZE big;
| 				ANALYZE big_archive;
| 			
< ok 0 -
> query
| EXPLAIN (FORMAT JSON) DELETE FROM big USING s WHERE big.s_id = s.id AND s.flag;
< result
< column "QUERY PLAN" "JSON"
< row
: bytes
| [
|   {
|     "Plan": {
|       "Node Type": "ModifyTable",
|       "Operation": "Delete",
|       "Parallel Aware": false,
|       "Async Capable": false,
|       "Relation Name": "big",
|       "Alias": "big",
|       "Startup Cost": 2.12,
|       "Total Cost": 20.86,
|       "Plan Rows": 0,
|       "Plan Width": 0,
|       "Plans": [
|         {
|           "Node Type": "Hash Join",
|           "Parent Relationship": "Outer",
|           "Parallel Aware": false,
|           "Async Capable": false,
|           "Join Type": "Inner",
|           "Startup Cost": 2.12,
|           "Total Cost": 20.86,
|           "Plan Rows": 100,
|           "Plan Width": 12,
|           "Inner Unique": true,
|           "Hash Cond": "(big.s_id = s.id)",
|           "Plans": [
|             {
|               "Node Type": "Seq Scan",
|               "Parent Relationship": "Outer",
|               "Parallel Aware": false,
|               "Async Capable": false,
|               "Relation Name": "big",
|               "Alias": "big",
|               "Startup Cost": 0.00,
|               "Total Cost": 16.00,
|               "Plan Rows": 1000,
|               "Plan Width": 10
|             },
|             {
|               "Node Type": "Hash",
|               "Parent Relationship": "Inner",
|               "Parallel Aware": false,
|               "Async Capable": false,
|               "Startup Cost": 2.00,
|               "Total Cost": 2.00,
|               "Plan Rows": 10,
|               "Plan Width": 10,
|               "Plans": [
|                 {
|                   "Node Type": "Seq Scan",
|                   "Parent Relationship": "Outer",
|                   "Parallel Aware": false,
|                   "Async Capable": false,
|                   "Relation Name": "s",
|                   "Alias": "s",
|                   "Startup Cost": 0.00,
|                   "Total Cost": 2.00,
|                   "Plan Rows": 10,
|                   "Plan Width": 10,
|                   "Filter": "flag"
|                 }
|               ]
|             }
|           ]
|         }
|       ]
|     }
|   }
| ]
> query
| EXPLAIN (FORMAT JSON) UPDATE big SET v = 1 FROM s WHERE big.s_id = s.id AND s.flag;
< result
< column "QUERY PLAN" "JSON"
< row
: bytes
| [
|   {
|     "Plan": {
|       "Node Type": "ModifyTable",
|       "Operation": "Update",
|       "Parallel Aware": false,
|       "Async Capable": false,
|       "Relation Name": "big",
|       "Alias": "big",
|       "Startup Cost": 2.12,
|       "Total Cost": 20.86,
|       "Plan Rows": 0,
|       "Plan Width": 0,
|       "Plans": [
|         {
|           "Node Type": "Hash Join",
|           "Parent Relationship": "Outer",
|           "Parallel Aware": false,
|           "Async Capable": false,
|           "Join Type": "Inner",
|           "Startup Cost": 2.12,
|           "Total Cost": 20.86,
|           "Plan Rows": 100,
|           "Plan Width": 16,
|           "Inner Unique": true,
|           "Hash Cond": "(big.s_id = s.id)",
|           "Plans": [
|             {
|               "Node Type": "Seq Scan",
|               "Parent Relationship": "Outer",
|               "Parallel Aware": false,
|               "Async Capable": false,
|               "Relation Name": "big",
|               "Alias": "big",
|               "Startup Cost": 0.00,
|               "Total Cost": 16.00,
|               "Plan Rows": 1000,
|               "Plan Width": 10
|             },
|             {
|               "Node Type": "Hash",
|               "Parent Relationship": "Inner",
|               "Parallel Aware": false,
|               "Async Capable": false,
|               "Startup Cost": 2.00,
|               "Total Cost": 2.00,
|               "Plan Rows": 10,
|               "Plan Width": 10,
|               "Plans": [
|                 {
|                   "Node Type": "Seq Scan",
|                   "Parent Relationship": "Outer",
|                   "Parallel Aware": false,
|                   "Async Capable": false,
|                   "Relation Name": "s",
|                   "Alias": "s",
|                   "Startup Cost": 0.00,
|                   "Total Cost": 2.00,
|                   "Plan Rows": 10,
|                   "Plan Width": 10,
|                   "Filter": "flag"
|                 }
|               ]
|             }
|           ]
|         }
|       ]
|     }
|   }
| ]
> query
| EXPLAIN (FORMAT JSON) DELETE FROM big WHERE s_id IN (SELECT id FROM s WHERE flag);
< result
< column "QUERY PLAN" "JSON"
< row
: bytes
| [
|   {
|     "Plan": {
|       "Node Type": "ModifyTable",
|       "Operation": "Delete",
|       "Parallel Aware": false,
|       "Async Capable": false,
|       "Relation Name": "big",
|       "Alias": "big",
|       "Startup Cost": 2.12,
|       "Total Cost": 20.86,
|       "Plan Rows": 0,
|       "Plan Width": 0,
|       "Plans": [
|         {
|           "Node Type": "Hash Join",
|           "Parent Relationship": "Outer",
|           "Parallel Aware": false,
|           "Async Capable": false,
|           "Join Type": "Inner",
|           "Startup Cost": 2.12,
|           "Total Cost": 20.86,
|           "Plan Rows": 100,
|           "Plan Width": 12,
|           "Inner Unique": true,
|           "Hash Cond": "(big.s_id = s.id)",
|           "Plans": [
|             {
|               "Node Type": "Seq Scan",
|               "Parent Relationship": "Outer",
|               "Parallel Aware": false,
|               "Async Capable": false,
|               "Relation Name": "big",
|               "Alias": "big",
|               "Startup Cost": 0.00,
|               "Total Cost": 16.00,
|               "Plan Rows": 1000,
|               "Plan Width": 10
|             },
|             {
|               "Node Type": "Hash",
|               "Parent Relationship": "Inner",
|               "Parallel Aware": false,
|               "Async Capable": false,
|               "Startup Cost": 2.00,
|               "Total Cost": 2.00,
|               "Plan Rows": 10,
|               "Plan Width": 10,
|               "Plans": [
|                 {
|                   "Node Type": "Seq Scan",
|                   "Parent Relationship": "Outer",
|                   "Parallel Aware": false,
|                   "Async Capable": false,
|                   "Relation Name": "s",
|                   "Alias": "s",
|                   "Startup Cost": 0.00,
|                   "Total Cost": 2.00,
|                   "Plan Rows": 10,
|                   "Plan Width": 10,
|                   "Filter": "flag"
|                 }
|               ]
|             }
|           ]
|         }
|       ]
|     }
|   }
| ]
> query
| EXPLAIN (FORMAT JSON) INSERT INTO big_archive SELECT big.* FROM big JOIN s ON big.s_id = s.id WHERE s.flag;
< result
< column "QUERY PLAN" "JSON"
< row
: bytes
| [
|   {
|     "Plan": {
|       "Node Type": "ModifyTable",
|       "Operation": "Insert",
|       "Parallel Aware": false,
|       "Async Capable": false,
|       "Relation Name": "big_archive",
|       "Alias": "big_archive",
|       "Startup Cost": 2.12,
|       "Total Cost": 20.86,
|       "Plan Rows": 0,
|       "Plan Width": 0,
|       "Plans": [
|         {
|           "Node Type": "Hash Join",
|           "Parent Relationship": "Outer",
|           "Parallel Aware": false,
|           "Async Capable": false,
|           "Join Type": "Inner",
|           "Startup Cost": 2.12,
|           "Total Cost": 20.86,
|           "Plan Rows": 100,
|           "Plan Width": 12,
|           "Inner Unique": true,
|           "Hash Cond": "(big.s_id = s.id)",
|           "Plans": [
|             {
|               "Node Type": "Seq Scan",
|               "Parent Relationship": "Outer",
|               "Parallel Aware": false,
|               "Async Capable": false,
|               "Relation Name": "big",
|               "Alias": "big",
|               "Startup Cost": 0.00,
|               "Total Cost": 16.00,
|               "Plan Rows": 1000,
|               "Plan Width": 12
|             },
|             {
|               "Node Type": "Hash",
|               "Parent Relationship": "Inner",
|               "Parallel Aware": false,
|               "Async Capable": false,
|               "Startup Cost": 2.00,
|               "Total Cost": 2.00,
|               "Plan Rows": 10,
|               "Plan Width": 4,
|               "Plans": [
|                 {
|                   "Node Type": "Seq Scan",
|                   "Parent Relationship": "Outer",
|                   "Parallel Aware": false,
|                   "Async Capable": false,
|                   "Relation Name": "s",
|                   "Alias": "s",
|                   "Startup Cost": 0.00,
|                   "Total Cost": 2.00,
|                   "Plan Rows": 10,
|                   "Plan Width": 4,
|                   "Filter": "flag"
|                 }
|               ]
|             }
|           ]
|         }
|       ]
|     }
|   }
| ]
> query
| EXPLAIN (FORMAT JSON) MERGE INTO big USING s ON big.s_id = s.id AND s.flag WHEN MATCHED THEN UPDATE SET v = 1;
< result
< column "QUERY PLAN" "JSON"
< row
: bytes
| [
|   {
|     "Plan": {
|       "Node Type": "ModifyTable",
|       "Operation": "Merge",
|       "Parallel Aware": false,
|       "Async Capable": false,
|       "Relation Name": "big",
|       "Alias": "big",
|       "Startup Cost": 2.12,
|       "Total Cost": 20.86,
|       "Plan Rows": 0,
|       "Plan Width": 0,
|       "Plans": [
|         {
|           "Node Type": "Hash Join",
|           "Parent Relationship": "Outer",
|           "Parallel Aware": false,
|           "Async Capable": false,
|           "Join Type": "Inner",
|           "Startup Cost": 2.12,
|           "Total Cost": 20.86,
|           "Plan Rows": 100,
|           "Plan Width": 12,
|           "Inner Unique": true,
|           "Hash Cond": "(big.s_id = s.id)",
|           "Plans": [
|             {
|               "Node Type": "Seq Scan",
|               "Parent Relationship": "Outer",
|               "Parallel Aware": false,
|               "Async Capable": false,
|               "Relation Name": "big",
|               "Alias": "big",
|               "Startup Cost": 0.00,
|               "Total Cost": 16.00,
|               "Plan Rows": 1000,
|               "Plan Width": 10
|             },
|             {
|               "Node Type": "Hash",
|               "Parent Relationship": "Inner",
|               "Parallel Aware": false,
|               "Async Capable": false,
|               "Startup Cost": 2.00,
|               "Total Cost": 2.00,
|               "Plan Rows": 10,
|               "Plan Width": 10,
|               "Plans": [
|                 {
|                   "Node Type": "Seq Scan",
|                   "Parent Relationship": "Outer",
|                   "Parallel Aware": false,
|                   "Async Capable": false,
|                   "Relation Name": "s",
|                   "Alias": "s",
|                   "Startup Cost": 0.00,
|                   "Total Cost": 2.00,
|                   "Plan Rows": 10,
|                   "Plan Width": 10,
|                   "Filter": "flag"
|                 }
|               ]
|             }
|           ]
|         }
|       ]
|     }
|   }
| ]
> query
| EXPLAIN (FORMAT JSON) SELECT * FROM big;
< result
< column "QUERY PLAN" "JSON"
< row
: bytes
| [
|   {
|     "Plan": {
|       "Node Type": "Seq Scan",
|       "Parallel Aware": false,
|       "Async Capable": false,
|       "Relation Name": "big",
|       "Alias": "big",
|       "Startup Cost": 0.00,
|       "Total Cost": 16.00,
|       "Plan Rows": 1000,
|       "Plan Width": 12
|     }
|   }
| ]
> exec
|
| 				CREATE SCHEMA app;
| 				CREATE TABLE app.only_in_app (id int);
| 				INSERT INTO app.only_in_app SELECT generate_series(1, 50);
| 				ANALYZE app.only_in_app;
| 			
< ok 0 -
> begin
< ok
> exec
| SET LOCAL search_path TO "app"
< ok 0 -
> query
| EXPLAIN (FORMAT JSON) DELETE FROM only_in_app;
< result
< column "QUERY PLAN" "JSON"
< row
: bytes
| [
|   {
|     "Plan": {
|       "Node Type": "ModifyTable",
|       "Operation": "Delete",
|       "Parallel Aware": false,
|       "Async Capable": false,
|       "Relation Name": "only_in_app",
|       "Alias": "only_in_app",
|       "Startup Cost": 0.00,
|       "Total Cost": 1.50,
|       "Plan Rows": 0,
|       "Plan Width": 0,
|       "Plans": [
|         {
|           "Node Type": "Seq Scan",
|           "Parent Relationship": "Outer",
|           "Parallel Aware": false,
|           "Async Capable": false,
|           "Relation Name": "only_in_app",
|           "Alias": "only_in_app",
|           "Startup Cost": 0.00,
|           "Total Cost": 1.50,
|           "Plan Rows": 50,
|           "Plan Width": 6
|         }
|       ]
|     }
|   }
| ]
> rollback
< ok
> query
| EXPLAIN (FORMAT JSON) DELETE FROM only_in_app;
< error
| 42P01 relation "only_in_app" does not exist
| position: 35

== TestCountAffectedRows/pg12 pg12#1
> exec
|
| 				CREATE TABLE s (id int PRIMARY KEY, flag bool NOT NULL);
| 				INSERT INTO s SELECT g, g <= 10 FROM generate_series(1, 100) g;
| 				CREATE TABLE big (id int PRIMARY KEY, s_id int NOT NULL, v int);
| 				INSERT INTO big SELECT g, (g % 100) + 1, 0 FROM generate_series(1, 1000) g;
| 				CREATE TABLE big_archive (LIKE big);
| 				ANALYZE s;
| 				ANALYZE big;
| 				ANALYZE big_archive;
| 			
< ok 0 -
> query
| EXPLAIN (FORMAT JSON) DELETE FROM big USING s WHERE big.s_id = s.id AND s.flag;
< result
< column "QUERY PLAN" "JSON"
< row
: bytes
| [
|   {
|     "Plan": {
|       "Node Type": "ModifyTable",
|       "Operation": "Delete",
|       "Parallel Aware": false,
|       "Relation Name": "big",
|       "Alias": "big",
|       "Startup Cost": 2.12,
|       "Total Cost": 20.86,
|       "Plan Rows": 100,
|       "Plan Width": 12,
|       "Plans": [
|         {
|           "Node Type": "Hash Join",
|           "Parent Relationship": "Member",
|           "Parallel Aware": false,
|           "Join Type": "Inner",
|           "Startup Cost": 2.12,
|           "Total Cost": 20.86,
|           "Plan Rows": 100,
|           "Plan Width": 12,
|           "Inner Unique": true,
|           "Hash Cond": "(big.s_id = s.id)",
|           "Plans": [
|             {
|               "Node Type": "Seq Scan",
|               "Parent Relationship": "Outer",
|               "Parallel Aware": false,
|               "Relation Name": "big",
|               "Alias": "big",
|               "Startup Cost": 0.00,
|               "Total Cost": 16.00,
|               "Plan Rows": 1000,
|               "Plan Width": 10
|             },
|             {
|               "Node Type": "Hash",
|               "Parent Relationship": "Inner",
|               "Parallel Aware": false,
|               "Startup Cost": 2.00,
|               "Total Cost": 2.00,
|               "Plan Rows": 10,
|               "Plan Width": 10,
|               "Plans": [
|                 {
|                   "Node Type": "Seq Scan",
|                   "Parent Relationship": "Outer",
|                   "Parallel Aware": false,
|                   "Relation Name": "s",
|                   "Alias": "s",
|                   "Startup Cost": 0.00,
|                   "Total Cost": 2.00,
|                   "Plan Rows": 10,
|                   "Plan Width": 10,
|                   "Filter": "flag"
|                 }
|               ]
|             }
|           ]
|         }
|       ]
|     }
|   }
| ]
> query
| EXPLAIN (FORMAT JSON) UPDATE big SET v = 1 FROM s WHERE big.s_id = s.id AND s.flag;
< result
< column "QUERY PLAN" "JSON"
< row
: bytes
| [
|   {
|     "Plan": {
|       "Node Type": "ModifyTable",
|       "Operation": "Update",
|       "Parallel Aware": false,
|       "Relation Name": "big",
|       "Alias": "big",
|       "Startup Cost": 2.12,
|       "Total Cost": 20.86,
|       "Plan Rows": 100,
|       "Plan Width": 24,
|       "Plans": [
|         {
|           "Node Type": "Hash Join",
|           "Parent Relationship": "Member",
|           "Parallel Aware": false,
|           "Join Type": "Inner",
|           "Startup Cost": 2.12,
|           "Total Cost": 20.86,
|           "Plan Rows": 100,
|           "Plan Width": 24,
|           "Inner Unique": true,
|           "Hash Cond": "(big.s_id = s.id)",
|           "Plans": [
|             {
|               "Node Type": "Seq Scan",
|               "Parent Relationship": "Outer",
|               "Parallel Aware": false,
|               "Relation Name": "big",
|               "Alias": "big",
|               "Startup Cost": 0.00,
|               "Total Cost": 16.00,
|               "Plan Rows": 1000,
|               "Plan Width": 14
|             },
|             {
|               "Node Type": "Hash",
|               "Parent Relationship": "Inner",
|               "Parallel Aware": false,
|               "Startup Cost": 2.00,
|               "Total Cost": 2.00,
|               "Plan Rows": 10,
|               "Plan Width": 10,
|               "Plans": [
|                 {
|                   "Node Type": "Seq Scan",
|                   "Parent Relationship": "Outer",
|                   "Parallel Aware": false,
|                   "Relation Name": "s",
|                   "Alias": "s",
|                   "Startup Cost": 0.00,
|                   "Total Cost": 2.00,
|                   "Plan Rows": 10,
|                   "Plan Width": 10,
|                   "Filter": "flag"
|                 }
|               ]
|             }
|           ]
|         }
|       ]
|     }
|   }
| ]
> query
| EXPLAIN (FORMAT JSON) DELETE FROM big WHERE s_id IN (SELECT id FROM s WHERE flag);
< result
< column "QUERY PLAN" "JSON"
< row
: bytes
| [
|   {
|     "Plan": {
|       "Node Type": "ModifyTable",
|       "Operation": "Delete",
|       "Parallel Aware": false,
|       "Relation Name": "big",
|       "Alias": "big",
|       "Startup Cost": 2.12,
|       "Total Cost": 20.86,
|       "Plan Rows": 100,
|       "Plan Width": 12,
|       "Plans": [
|         {
|           "Node Type": "Hash Join",
|           "Parent Relationship": "Member",
|           "Parallel Aware": false,
|           "Join Type": "Inner",
|           "Startup Cost": 2.12,
|           "Total Cost": 20.86,
|           "Plan Rows": 100,
|           "Plan Width": 12,
|           "Inner Unique": true,
|           "Hash Cond": "(big.s_id = s.id)",
|           "Plans": [
|             {
|               "Node Type": "Seq Scan",
|               "Parent Relationship": "Outer",
|               "Parallel Aware": false,
|               "Relation Name": "big",
|               "Alias": "big",
|               "Startup Cost": 0.00,
|               "Total Cost": 16.00,
|               "Plan Rows": 1000,
|               "Plan Width": 10
|             },
|             {
|               "Node Type": "Hash",
|               "Parent Relationship": "Inner",
|               "Parallel Aware": false,
|               "Startup Cost": 2.00,
|               "Total Cost": 2.00,
|               "Plan Rows": 10,
|               "Plan Width": 10,
|               "Plans": [
|                 {
|                   "Node Type": "Seq Scan",
|                   "Parent Relationship": "Outer",
|                   "Parallel Aware": false,
|                   "Relation Name": "s",
|                   "Alias": "s",
|                   "Startup Cost": 0.00,
|                   "Total Cost": 2.00,
|                   "Plan Rows": 10,
|                   "Plan Width": 10,
|                   "Filter": "flag"
|                 }
|               ]
|             }
|           ]
|         }
|       ]
|     }
|   }
| ]
> query
| EXPLAIN (FORMAT JSON) INSERT INTO big_archive SELECT big.* FROM big JOIN s ON big.s_id = s.id WHERE s.flag;
< result
< column "QUERY PLAN" "JSON"
< row
: bytes
| [
|   {
|     "Plan": {
|       "Node Type": "ModifyTable",
|       "Operation": "Insert",
|       "Parallel Aware": false,
|       "Relation Name": "big_archive",
|       "Alias": "big_archive",
|       "Startup Cost": 2.12,
|       "Total Cost": 20.86,
|       "Plan Rows": 100,
|       "Plan Width": 12,
|       "Plans": [
|         {
|           "Node Type": "Hash Join",
|           "Parent Relationship": "Member",
|           "Parallel Aware": false,
|           "Join Type": "Inner",
|           "Startup Cost": 2.12,
|           "Total Cost": 20.86,
|           "Plan Rows": 100,
|           "Plan Width": 12,
|           "Inner Unique": true,
|           "Hash Cond": "(big.s_id = s.id)",
|           "Plans": [
|             {
|               "Node Type": "Seq Scan",
|               "Parent Relationship": "Outer",
|               "Parallel Aware": false,
|               "Relation Name": "big",
|               "Alias": "big",
|               "Startup Cost": 0.00,
|               "Total Cost": 16.00,
|               "Plan Rows": 1000,
|               "Plan Width": 12
|             },
|             {
|               "Node Type": "Hash",
|               "Parent Relationship": "Inner",
|               "Parallel Aware": false,
|               "Startup Cost": 2.00,
|               "Total Cost": 2.00,
|               "Plan Rows": 10,
|               "Plan Width": 4,
|               "Plans": [
|                 {
|                   "Node Type": "Seq Scan",
|                   "Parent Relationship": "Outer",
|                   "Parallel Aware": false,
|                   "Relation Name": "s",
|                   "Alias": "s",
|                   "Startup Cost": 0.00,
|                   "Total Cost": 2.00,
|                   "Plan Rows": 10,
|                   "Plan Width": 4,
|                   "Filter": "flag"
|                 }
|               ]
|             }
|           ]
|         }
|       ]
|     }
|   }
| ]
> query
| EXPLAIN (FORMAT JSON) SELECT * FROM big;
< result
< column "QUERY PLAN" "JSON"
< row
: bytes
| [
|   {
|     "Plan": {
|       "Node Type": "Seq Scan",
|       "Parallel Aware": false,
|       "Relation Name": "big",
|       "Alias": "big",
|       "Startup Cost": 0.00,
|       "Total Cost": 16.00,
|       "Plan Rows": 1000,
|       "Plan Width": 12
|     }
|   }
| ]
> exec
|
| 				CREATE SCHEMA app;
| 				CREATE TABLE app.only_in_app (id int);
| 				INSERT INTO app.only_in_app SELECT generate_series(1, 50);
| 				ANALYZE app.only_in_app;
| 			
< ok 0 -
> begin
< ok
> exec
| SET LOCAL search_path TO "app"
< ok 0 -
> query
| EXPLAIN (FORMAT JSON) DELETE FROM only_in_app;
< result
< column "QUERY PLAN" "JSON"
< row
: bytes
| [
|   {
|     "Plan": {
|       "Node Type": "ModifyTable",
|       "Operation": "Delete",
|       "Parallel Aware": false,
|       "Relation Name": "only_in_app",
|       "Alias": "only_in_app",
|       "Startup Cost": 0.00,
|       "Total Cost": 1.50,
|       "Plan Rows": 50,
|       "Plan Width": 6,
|       "Plans": [
|         {
|           "Node Type": "Seq Scan",
|           "Parent Relationship": "Member",
|           "Parallel Aware": false,
|           "Relation Name": "only_in_app",
|           "Alias": "only_in_app",
|           "Startup Cost": 0.00,
|           "Total Cost": 1.50,
|           "Plan Rows": 50,
|           "Plan Width": 6
|         }
|       ]
|     }
|   }
| ]
> rollback
< ok
> query
| EXPLAIN (FORMAT JSON) DELETE FROM only_in_app;
< error
| 42P01 relation "only_in_app" does not exist
| position: 35

== TestCountAffectedRows/pg18 pg18#1
> exec
|
| 				CREATE TABLE s (id int PRIMARY KEY, flag bool NOT NULL);
| 				INSERT INTO s SELECT g, g <= 10 FROM generate_series(1, 100) g;
| 				CREATE TABLE big (id int PRIMARY KEY, s_id int NOT NULL, v int);
| 				INSERT INTO big SELECT g, (g % 100) + 1, 0 FROM generate_series(1, 1000) g;
| 				CREATE TABLE big_archive (LIKE big);
| 				ANALYZE s;
| 				ANALYZE big;
| 				ANALYZE big_archive;
| 			
< ok 0 -
> query
| EXPLAIN (FORMAT JSON) DELETE FROM big USING s WHERE big.s_id = s.id AND s.flag;
< result
< column "QUERY PLAN" "JSON"
< row
: bytes
| [
|   {
|     "Plan": {
|       "Node Type": "ModifyTable",
|       "Operation": "Delete",
|       "Parallel Aware": false,
|       "Async Capable": false,
|       "Relation Name": "big",
|       "Alias": "big",
|       "Startup Cost": 2.12,
|       "Total Cost": 20.86,
|       "Plan Rows": 0,
|       "Plan Width": 0,
|       "Disabled": false,
|       "Plans": [
|         {
|           "Node Type": "Hash Join",
|           "Parent Relationship": "Outer",
|           "Parallel Aware": false,
|           "Async Capable": false,
|           "Join Type": "Inner",
|           "Startup Cost": 2.12,
|           "Total Cost": 20.86,
|           "Plan Rows": 100,
|           "Plan Width": 12,
|           "Disabled": false,
|           "Inner Unique": true,
|           "Hash Cond": "(big.s_id = s.id)",
|           "Plans": [
|             {
|               "Node Type": "Seq Scan",
|               "Parent Relationship": "Outer",
|               "Parallel Aware": false,
|               "Async Capable": false,
|               "Relation Name": "big",
|               "Alias": "big",
|               "Startup Cost": 0.00,
|               "Total Cost": 16.00,
|               "Plan Rows": 1000,
|               "Plan Width": 10,
|               "Disabled": false
|             },
|             {
|               "Node Type": "Hash",
|               "Parent Relationship": "Inner",
|               "Parallel Aware": false,
|               "Async Capable": false,
|               "Startup Cost": 2.00,
|               "Total Cost": 2.00,
|               "Plan Rows": 10,
|               "Plan Width": 10,
|               "Disabled": false,
|               "Plans": [
|                 {
|                   "Node Type": "Seq Scan",
|                   "Parent Relationship": "Outer",
|                   "Parallel Aware": false,
|                   "Async Capable": false,
|                   "Relation Name": "s",
|                   "Alias": "s",
|                   "Startup Cost": 0.00,
|                   "Total Cost": 2.00,
|                   "Plan Rows": 10,
|                   "Plan Width": 10,
|                   "Disabled": false,
|                   "Filter": "flag"
|                 }
|               ]
|             }
|           ]
|         }
|       ]
|     }
|   }
| ]
> query
| EXPLAIN (FORMAT JSON) UPDATE big SET v = 1 FROM s WHERE big.s_id = s.id AND s.flag;
< result
< column "QUERY PLAN" "JSON"
< row
: bytes
| [
|   {
|     "Plan": {
|       "Node Type": "ModifyTable",
|       "Operation": "Update",
|       "Parallel Aware": false,
|       "Async Capable": false,
|       "Relation Name": "big",
|       "Alias": "big",
|       "Startup Cost": 2.12,
|       "Total Cost": 20.86,
|       "Plan Rows": 0,
|       "Plan Width": 0,
|       "Disabled": false,
|       "Plans": [
|         {
|           "Node Type": "Hash Join",
|           "Parent Relationship": "Outer",
|           "Parallel Aware": false,
|           "Async Capable": false,
|           "Join Type": "Inner",
|           "Startup Cost": 2.12,
|           "Total Cost": 20.86,
|           "Plan Rows": 100,
|           "Plan Width": 16,
|           "Disabled": false,
|           "Inner Unique": true,
|           "Hash Cond": "(big.s_id = s.id)",
|           "Plans": [
|             {
|               "Node Type": "Seq Scan",
|               "Parent Relationship": "Outer",
|               "Parallel Aware": false,
|               "Async Capable": false,
|               "Relation Name": "big",
|               "Alias": "big",
|               "Startup Cost": 0.00,
|               "Total Cost": 16.00,
|               "Plan Rows": 1000,
|               "Plan Width": 10,
|               "Disabled": false
|             },
|             {
|               "Node Type": "Hash",
|               "Parent Relationship": "Inner",
|               "Parallel Aware": false,
|               "Async Capable": false,
|               "Startup Cost": 2.00,
|               "Total Cost": 2.00,
|               "Plan Rows": 10,
|               "Plan Width": 10,
|               "Disabled": false,
|               "Plans": [
|                 {
|                   "Node Type": "Seq Scan",
|                   "Parent Relationship": "Outer",
|                   "Parallel Aware": false,
|                   "Async Capable": false,
|                   "Relation Name": "s",
|                   "Alias": "s",
|                   "Startup Cost": 0.00,
|                   "Total Cost": 2.00,
|                   "Plan Rows": 10,
|                   "Plan Width": 10,
|                   "Disabled": false,
|                   "Filter": "flag"
|                 }
|               ]
|             }
|           ]
|         }
|       ]
|     }
|   }
| ]
> query
| EXPLAIN (FORMAT JSON) DELETE FROM big WHERE s_id IN (SELECT id FROM s WHERE flag);
< result
< column "QUERY PLAN" "JSON"
< row
: bytes
| [
|   {
|     "Plan": {
|       "Node Type": "ModifyTable",
|       "Operation": "Delete",
|       "Parallel Aware": false,
|       "Async Capable": false,
|       "Relation Name": "big",
|       "Alias": "big",
|       "Startup Cost": 2.12,
|       "Total Cost": 20.86,
|       "Plan Rows": 0,
|       "Plan Width": 0,
|       "Disabled": false,
|       "Plans": [
|         {
|           "Node Type": "Hash Join",
|           "Parent Relationship": "Outer",
|           "Parallel Aware": false,
|           "Async Capable": false,
|           "Join Type": "Inner",
|           "Startup Cost": 2.12,
|           "Total Cost": 20.86,
|           "Plan Rows": 100,
|           "Plan Width": 12,
|           "Disabled": false,
|           "Inner Unique": true,
|           "Hash Cond": "(big.s_id = s.id)",
|           "Plans": [
|             {
|               "Node Type": "Seq Scan",
|               "Parent Relationship": "Outer",
|               "Parallel Aware": false,
|               "Async Capable": false,
|               "Relation Name": "big",
|               "Alias": "big",
|               "Startup Cost": 0.00,
|               "Total Cost": 16.00,
|               "Plan Rows": 1000,
|               "Plan Width": 10,
|               "Disabled": false
|             },
|             {
|               "Node Type": "Hash",
|               "Parent Relationship": "Inner",
|               "Parallel Aware": false,
|               "Async Capable": false,
|               "Startup Cost": 2.00,
|               "Total Cost": 2.00,
|               "Plan Rows": 10,
|               "Plan Width": 10,
|               "Disabled": false,
|               "Plans": [
|                 {
|                   "Node Type": "Seq Scan",
|                   "Parent Relationship": "Outer",
|                   "Parallel Aware": false,
|                   "Async Capable": false,
|                   "Relation Name": "s",
|                   "Alias": "s",
|                   "Startup Cost": 0.00,
|                   "Total Cost": 2.00,
|                   "Plan Rows": 10,
|                   "Plan Width": 10,
|                   "Disabled": false,
|                   "Filter": "flag"
|                 }
|               ]
|             }
|           ]
|         }
|       ]
|     }
|   }
| ]
> query
| EXPLAIN (FORMAT JSON) INSERT INTO big_archive SELECT big.* FROM big JOIN s ON big.s_id = s.id WHERE s.flag;
< result
< column "QUERY PLAN" "JSON"
< row
: bytes
| [
|   {
|     "Plan": {
|       "Node Type": "ModifyTable",
|       "Operation": "Insert",
|       "Parallel Aware": false,
|       "Async Capable": false,
|       "Relation Name": "big_archive",
|       "Alias": "big_archive",
|       "Startup Cost": 2.12,
|       "Total Cost": 20.86,
|       "Plan Rows": 0,
|       "Plan Width": 0,
|       "Disabled": false,
|       "Plans": [
|         {
|           "Node Type": "Hash Join",
|           "Parent Relationship": "Outer",
|           "Parallel Aware": false,
|           "Async Capable": false,
|           "Join Type": "Inner",
|           "Startup Cost": 2.12,
|           "Total Cost": 20.86,
|           "Plan Rows": 100,
|           "Plan Width": 12,
|           "Disabled": false,
|           "Inner Unique": true,
|           "Hash Cond": "(big.s_id = s.id)",
|           "Plans": [
|             {
|               "Node Type": "Seq Scan",
|               "Parent Relationship": "Outer",
|               "Parallel Aware": false,
|               "Async Capable": false,
|               "Relation Name": "big",
|               "Alias": "big",
|               "Startup Cost": 0.00,
|               "Total Cost": 16.00,
|               "Plan Rows": 1000,
|               "Plan Width": 12,
|               "Disabled": false
|             },
|             {
|               "Node Type": "Hash",
|               "Parent Relationship": "Inner",
|               "Parallel Aware": false,
|               "Async Capable": false,
|               "Startup Cost": 2.00,
|               "Total Cost": 2.00,
|               "Plan Rows": 10,
|               "Plan Width": 4,
|               "Disabled": false,
|               "Plans": [
|                 {
|                   "Node Type": "Seq Scan",
|                   "Parent Relationship": "Outer",
|                   "Parallel Aware": false,
|                   "Async Capable": false,
|                   "Relation Name": "s",
|                   "Alias": "s",
|                   "Startup Cost": 0.00,
|                   "Total Cost": 2.00,
|                   "Plan Rows": 10,
|                   "Plan Width": 4,
|                   "Disabled": false,
|                   "Filter": "flag"
|                 }
|               ]
|             }
|           ]
|         }
|       ]
|     }
|   }
| ]
> query
| EXPLAIN (FORMAT JSON) MERGE INTO big USING s ON big.s_id = s.id AND s.flag WHEN MATCHED THEN UPDATE SET v = 1;
< result
< column "QUERY PLAN" "JSON"
< row
: bytes
| [
|   {
|     "Plan": {
|       "Node Type": "ModifyTable",
|       "Operation": "Merge",
|       "Parallel Aware": false,
|       "Async Capable": false,
|       "Relation Name": "big",
|       "Alias": "big",
|       "Startup Cost": 2.12,
|       "Total Cost": 20.86,
|       "Plan Rows": 0,
|       "Plan Width": 0,
|       "Disabled": false,
|       "Plans": [
|         {
|           "Node Type": "Hash Join",
|           "Parent Relationship": "Outer",
|           "Parallel Aware": false,
|           "Async Capable": false,
|           "Join Type": "Inner",
|           "Startup Cost": 2.12,
|           "Total Cost": 20.86,
|           "Plan Rows": 100,
|           "Plan Width": 12,
|           "Disabled": false,
|           "Inner Unique": true,
|           "Hash Cond": "(big.s_id = s.id)",
|           "Plans": [
|             {
|               "Node Type": "Seq Scan",
|               "Parent Relationship": "Outer",
|               "Parallel Aware": false,
|               "Async Capable": false,
|               "Relation Name": "big",
|               "Alias": "big",
|               "Startup Cost": 0.00,
|               "Total Cost": 16.00,
|               "Plan Rows": 1000,
|               "Plan Width": 10,
|               "Disabled": false
|             },
|             {
|               "Node Type": "Hash",
|               "Parent Relationship": "Inner",
|               "Parallel Aware": false,
|               "Async Capable": false,
|               "Startup Cost": 2.00,
|               "Total Cost": 2.00,
|               "Plan Rows": 10,
|               "Plan Width": 10,
|               "Disabled": false,
|               "Plans": [
|                 {
|                   "Node Type": "Seq Scan",
|                   "Parent Relationship": "Outer",
|                   "Parallel Aware": false,
|                   "Async Capable": false,
|                   "Relation Name": "s",
|                   "Alias": "s",
|                   "Startup Cost": 0.00,
|                   "Total Cost": 2.00,
|                   "Plan Rows": 10,
|                   "Plan Width": 10,
|                   "Disabled": false,
|                   "Filter": "flag"
|                 }
|               ]
|             }
|           ]
|         }
|       ]
|     }
|   }
| ]
> query
| EXPLAIN (FORMAT JSON) SELECT * FROM big;
< result
< column "QUERY PLAN" "JSON"
< row
: bytes
| [
|   {
|     "Plan": {
|       "Node Type": "Seq Scan",
|       "Parallel Aware": false,
|       "Async Capable": false,
|       "Relation Name": "big",
|       "Alias": "big",
|       "Startup Cost": 0.00,
|       "Total Cost": 16.00,
|       "Plan Rows": 1000,
|       "Plan Width": 12,
|       "Disabled": false
|     }
|   }
| ]
> exec
|
| 				CREATE SCHEMA app;
| 				CREATE TABLE app.only_in_app (id int);
| 				INSERT INTO app.only_in_app SELECT generate_series(1, 50);
| 				ANALYZE app.only_in_app;
| 			
< ok 0 -
> begin
< ok
> exec
| SET LOCAL search_path TO "app"
< ok 0 -
> query
| EXPLAIN (FORMAT JSON) DELETE FROM only_in_app;
< result
< column "QUERY PLAN" "JSON"
< row
: bytes
| [
|   {
|     "Plan": {
|       "Node Type": "ModifyTable",
|       "Operation": "Delete",
|       "Parallel Aware": false,
|       "Async Capable": false,
|       "Relation Name": "only_in_app",
|       "Alias": "only_in_app",
|       "Startup Cost": 0.00,
|       "Total Cost": 1.50,
|       "Plan Rows": 0,
|       "Plan Width": 0,
|       "Disabled": false,
|       "Plans": [
|         {
|           "Node Type": "Seq Scan",
|           "Parent Relationship": "Outer",
|           "Parallel Aware": false,
|           "Async Capable": false,
|           "Relation Name": "only_in_app",
|           "Alias": "only_in_app",
|           "Startup Cost": 0.00,
|           "Total Cost": 1.50,
|           "Plan Rows": 50,
|           "Plan Width": 6,
|           "Disabled": false
|         }
|       ]
|     }
|   }
| ]
> rollback
< ok
> query
| EXPLAIN (FORMAT JSON) DELETE FROM only_in_app;
< error
| 42P01 relation "only_in_app" does not exist
| position: 35
