-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12
-- pin: pg12 postgres:12.22-alpine sha256:7c8f4870583184ebadf7f17a6513620aac5f365a7938dc6a6911c1d5df2f481a REL_12_22

== TestExplainAndRoles/default default#1
> exec
|
| 				CREATE SCHEMA app;
| 				CREATE TABLE app.items (id int PRIMARY KEY);
| 				INSERT INTO app.items SELECT generate_series(1, 25);
| 				ANALYZE app.items;
| 				CREATE ROLE explain_owner NOLOGIN;
| 				GRANT USAGE ON SCHEMA app TO explain_owner;
| 				GRANT SELECT, INSERT, DELETE ON app.items TO explain_owner;
| 			
< ok 0 -
> exec
| ALTER DATABASE "omni_testexplainandroles" OWNER TO explain_owner;
< ok 0 -
> query
|
| 		SELECT
| 			u.rolname
| 		FROM
| 			pg_roles AS u JOIN pg_database AS d ON (d.datdba = u.oid)
| 		WHERE
| 			d.datname = current_database();
| 		
< result
< column "rolname" "NAME"
< row
: string
| explain_owner
> query
| SELECT CURRENT_USER
< result
< column "current_user" "NAME"
< row
: string
| postgres
> query
| SHOW search_path
< result
< column "search_path" "TEXT"
< row
: string
| "$user", public
> begin
< ok
> query
|
| 		SELECT
| 			u.rolname
| 		FROM
| 			pg_roles AS u JOIN pg_database AS d ON (d.datdba = u.oid)
| 		WHERE
| 			d.datname = current_database();
| 		
< result
< column "rolname" "NAME"
< row
: string
| explain_owner
> exec
| SET SESSION ROLE 'explain_owner';
< ok 0 -
> exec
| SET LOCAL search_path TO "app"
< ok 0 -
> query
| EXPLAIN DELETE FROM items WHERE id = 1
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Delete on items  (cost=0.00..1.31 rows=0 width=0)
< row
: string
|   ->  Seq Scan on items  (cost=0.00..1.31 rows=1 width=6)
< row
: string
|         Filter: (id = 1)
> rollback
< ok
> begin
< ok
> query
|
| 		SELECT
| 			u.rolname
| 		FROM
| 			pg_roles AS u JOIN pg_database AS d ON (d.datdba = u.oid)
| 		WHERE
| 			d.datname = current_database();
| 		
< result
< column "rolname" "NAME"
< row
: string
| explain_owner
> exec
| SET SESSION ROLE 'explain_owner';
< ok 0 -
> query
| EXPLAIN DELETE FROM missing_table
< error
| 42P01 relation "missing_table" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> query
|
| 		SELECT
| 			u.rolname
| 		FROM
| 			pg_roles AS u JOIN pg_database AS d ON (d.datdba = u.oid)
| 		WHERE
| 			d.datname = current_database();
| 		
< result
< column "rolname" "NAME"
< row
: string
| explain_owner
> exec
| SET SESSION ROLE 'explain_owner';
< ok 0 -
> exec
| SET LOCAL search_path TO "app"
< ok 0 -
> query
| EXPLAIN (FORMAT JSON) INSERT INTO items SELECT id + 100 FROM items
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
|       "Relation Name": "items",
|       "Alias": "items",
|       "Startup Cost": 0.00,
|       "Total Cost": 1.31,
|       "Plan Rows": 0,
|       "Plan Width": 0,
|       "Plans": [
|         {
|           "Node Type": "Seq Scan",
|           "Parent Relationship": "Outer",
|           "Parallel Aware": false,
|           "Async Capable": false,
|           "Relation Name": "items",
|           "Alias": "items_1",
|           "Startup Cost": 0.00,
|           "Total Cost": 1.31,
|           "Plan Rows": 25,
|           "Plan Width": 4
|         }
|       ]
|     }
|   }
| ]
> rollback
< ok

== TestExplainAndRoles/pg12 pg12#1
> exec
|
| 				CREATE SCHEMA app;
| 				CREATE TABLE app.items (id int PRIMARY KEY);
| 				INSERT INTO app.items SELECT generate_series(1, 25);
| 				ANALYZE app.items;
| 				CREATE ROLE explain_owner NOLOGIN;
| 				GRANT USAGE ON SCHEMA app TO explain_owner;
| 				GRANT SELECT, INSERT, DELETE ON app.items TO explain_owner;
| 			
< ok 0 -
> exec
| ALTER DATABASE "omni_testexplainandroles" OWNER TO explain_owner;
< ok 0 -
> query
|
| 		SELECT
| 			u.rolname
| 		FROM
| 			pg_roles AS u JOIN pg_database AS d ON (d.datdba = u.oid)
| 		WHERE
| 			d.datname = current_database();
| 		
< result
< column "rolname" "NAME"
< row
: string
| explain_owner
> query
| SELECT CURRENT_USER
< result
< column "current_user" "NAME"
< row
: string
| postgres
> query
| SHOW search_path
< result
< column "search_path" "TEXT"
< row
: string
| "$user", public
> begin
< ok
> query
|
| 		SELECT
| 			u.rolname
| 		FROM
| 			pg_roles AS u JOIN pg_database AS d ON (d.datdba = u.oid)
| 		WHERE
| 			d.datname = current_database();
| 		
< result
< column "rolname" "NAME"
< row
: string
| explain_owner
> exec
| SET SESSION ROLE 'explain_owner';
< ok 0 -
> exec
| SET LOCAL search_path TO "app"
< ok 0 -
> query
| EXPLAIN DELETE FROM items WHERE id = 1
< result
< column "QUERY PLAN" "TEXT"
< row
: string
| Delete on items  (cost=0.00..1.31 rows=1 width=6)
< row
: string
|   ->  Seq Scan on items  (cost=0.00..1.31 rows=1 width=6)
< row
: string
|         Filter: (id = 1)
> rollback
< ok
> begin
< ok
> query
|
| 		SELECT
| 			u.rolname
| 		FROM
| 			pg_roles AS u JOIN pg_database AS d ON (d.datdba = u.oid)
| 		WHERE
| 			d.datname = current_database();
| 		
< result
< column "rolname" "NAME"
< row
: string
| explain_owner
> exec
| SET SESSION ROLE 'explain_owner';
< ok 0 -
> query
| EXPLAIN DELETE FROM missing_table
< error
| 42P01 relation "missing_table" does not exist
| position: 21
> rollback
< ok
> begin
< ok
> query
|
| 		SELECT
| 			u.rolname
| 		FROM
| 			pg_roles AS u JOIN pg_database AS d ON (d.datdba = u.oid)
| 		WHERE
| 			d.datname = current_database();
| 		
< result
< column "rolname" "NAME"
< row
: string
| explain_owner
> exec
| SET SESSION ROLE 'explain_owner';
< ok 0 -
> exec
| SET LOCAL search_path TO "app"
< ok 0 -
> query
| EXPLAIN (FORMAT JSON) INSERT INTO items SELECT id + 100 FROM items
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
|       "Relation Name": "items",
|       "Alias": "items",
|       "Startup Cost": 0.00,
|       "Total Cost": 1.31,
|       "Plan Rows": 25,
|       "Plan Width": 4,
|       "Plans": [
|         {
|           "Node Type": "Seq Scan",
|           "Parent Relationship": "Member",
|           "Parallel Aware": false,
|           "Relation Name": "items",
|           "Alias": "items_1",
|           "Startup Cost": 0.00,
|           "Total Cost": 1.31,
|           "Plan Rows": 25,
|           "Plan Width": 4
|         }
|       ]
|     }
|   }
| ]
> rollback
< ok
