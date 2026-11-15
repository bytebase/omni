-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12
-- pin: pg12 postgres:12.22-alpine sha256:7c8f4870583184ebadf7f17a6513620aac5f365a7938dc6a6911c1d5df2f481a REL_12_22
-- pin: pg18 postgres:18.7-alpine sha256:75dcac3dfb1cd01ca3536910bbc84128759362f8b118bf8c55b38a94f68f3b73 REL_18_7

== TestExecute/default default#1
> begin
< ok
> exec
| CREATE TABLE exec_t (id int PRIMARY KEY, v text);
< ok 0 -
> exec
|
| INSERT INTO exec_t VALUES (1, 'a'), (2, 'b'), (3, 'c');
< ok 3 -
> exec
|
| UPDATE exec_t SET v = 'x' WHERE id <= 2;
< ok 2 -
> commit
< ok
> exec
|
| CREATE INDEX CONCURRENTLY idx_exec_v ON exec_t (v);
< ok 0 -
> exec
|
| VACUUM exec_t;
< ok 0 -
> exec
|
| DROP INDEX CONCURRENTLY idx_exec_v;
< ok 0 -
> begin
< ok
> exec
| INSERT INTO exec_t VALUES (10, 'rolled_back');
< ok 1 -
> exec
|
| SELECT 1 / 0;
< error
| 22012 division by zero
> rollback
< ok
> query
| SELECT count(*) FROM exec_t WHERE id = 10
< result
< column "count" "INT8"
< row
: int64 0
> exec
| CREATE TABLE ac_order (id int);
< ok 0 -
> exec
|
| CREATE INDEX CONCURRENTLY idx_ac_order ON ac_order (id);
< ok 0 -
> exec
|
| INSERT INTO exec_t VALUES (20, 'SET ROLE reader; VACUUM');
< ok 1 -
> exec
|
| SELECT 1 / 0;
< error
| 22012 division by zero
> query
| SELECT count(*) FROM exec_t WHERE id = 20
< result
< column "count" "INT8"
< row
: int64 1
> exec
| DO $$ BEGIN INSERT INTO exec_t VALUES (30, 'do'); COMMIT; INSERT INTO exec_t VALUES (31, 'do2'); END $$;
< ok 0 -
> exec
|
| 				CREATE ROLE exec_tenant_owner NOLOGIN;
| 				CREATE ROLE exec_reader NOLOGIN;
| 				GRANT CREATE ON SCHEMA public TO exec_tenant_owner;
| 			
< ok 0 -
> exec
| ALTER DATABASE "omni_testexecute" OWNER TO exec_tenant_owner;
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
| exec_tenant_owner
> begin
< ok
> exec
| SET LOCAL ROLE 'exec_tenant_owner'
< ok 0 -
> exec
| CREATE TABLE tenant_owned (id int);
< ok 0 -
> exec
| SET LOCAL ROLE NONE;
| GRANT SELECT ON tenant_owned TO exec_reader;SET LOCAL ROLE 'exec_tenant_owner';
< ok 0 -
> commit
< ok
> exec
| SET SESSION ROLE 'exec_tenant_owner'
< ok 0 -
> query
| SELECT tableowner FROM pg_tables WHERE schemaname = 'public' AND tablename = 'tenant_owned'
< result
< column "tableowner" "NAME"
< row
: string
| exec_tenant_owner

== TestExecute/pg12 pg12#1
> begin
< ok
> exec
| CREATE TABLE exec_t (id int PRIMARY KEY, v text);
< ok 0 -
> exec
|
| INSERT INTO exec_t VALUES (1, 'a'), (2, 'b'), (3, 'c');
< ok 3 -
> exec
|
| UPDATE exec_t SET v = 'x' WHERE id <= 2;
< ok 2 -
> commit
< ok
> exec
|
| CREATE INDEX CONCURRENTLY idx_exec_v ON exec_t (v);
< ok 0 -
> exec
|
| VACUUM exec_t;
< ok 0 -
> exec
|
| DROP INDEX CONCURRENTLY idx_exec_v;
< ok 0 -
> begin
< ok
> exec
| INSERT INTO exec_t VALUES (10, 'rolled_back');
< ok 1 -
> exec
|
| SELECT 1 / 0;
< error
| 22012 division by zero
> rollback
< ok
> query
| SELECT count(*) FROM exec_t WHERE id = 10
< result
< column "count" "INT8"
< row
: int64 0
> exec
| CREATE TABLE ac_order (id int);
< ok 0 -
> exec
|
| CREATE INDEX CONCURRENTLY idx_ac_order ON ac_order (id);
< ok 0 -
> exec
|
| INSERT INTO exec_t VALUES (20, 'SET ROLE reader; VACUUM');
< ok 1 -
> exec
|
| SELECT 1 / 0;
< error
| 22012 division by zero
> query
| SELECT count(*) FROM exec_t WHERE id = 20
< result
< column "count" "INT8"
< row
: int64 1
> exec
| DO $$ BEGIN INSERT INTO exec_t VALUES (30, 'do'); COMMIT; INSERT INTO exec_t VALUES (31, 'do2'); END $$;
< ok 0 -
> exec
|
| 				CREATE ROLE exec_tenant_owner NOLOGIN;
| 				CREATE ROLE exec_reader NOLOGIN;
| 				GRANT CREATE ON SCHEMA public TO exec_tenant_owner;
| 			
< ok 0 -
> exec
| ALTER DATABASE "omni_testexecute" OWNER TO exec_tenant_owner;
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
| exec_tenant_owner
> begin
< ok
> exec
| SET LOCAL ROLE 'exec_tenant_owner'
< ok 0 -
> exec
| CREATE TABLE tenant_owned (id int);
< ok 0 -
> exec
| SET LOCAL ROLE NONE;
| GRANT SELECT ON tenant_owned TO exec_reader;SET LOCAL ROLE 'exec_tenant_owner';
< ok 0 -
> commit
< ok
> exec
| SET SESSION ROLE 'exec_tenant_owner'
< ok 0 -
> query
| SELECT tableowner FROM pg_tables WHERE schemaname = 'public' AND tablename = 'tenant_owned'
< result
< column "tableowner" "NAME"
< row
: string
| exec_tenant_owner

== TestExecute/pg18 pg18#1
> begin
< ok
> exec
| CREATE TABLE exec_t (id int PRIMARY KEY, v text);
< ok 0 -
> exec
|
| INSERT INTO exec_t VALUES (1, 'a'), (2, 'b'), (3, 'c');
< ok 3 -
> exec
|
| UPDATE exec_t SET v = 'x' WHERE id <= 2;
< ok 2 -
> commit
< ok
> exec
|
| CREATE INDEX CONCURRENTLY idx_exec_v ON exec_t (v);
< ok 0 -
> exec
|
| VACUUM exec_t;
< ok 0 -
> exec
|
| DROP INDEX CONCURRENTLY idx_exec_v;
< ok 0 -
> begin
< ok
> exec
| INSERT INTO exec_t VALUES (10, 'rolled_back');
< ok 1 -
> exec
|
| SELECT 1 / 0;
< error
| 22012 division by zero
> rollback
< ok
> query
| SELECT count(*) FROM exec_t WHERE id = 10
< result
< column "count" "INT8"
< row
: int64 0
> exec
| CREATE TABLE ac_order (id int);
< ok 0 -
> exec
|
| CREATE INDEX CONCURRENTLY idx_ac_order ON ac_order (id);
< ok 0 -
> exec
|
| INSERT INTO exec_t VALUES (20, 'SET ROLE reader; VACUUM');
< ok 1 -
> exec
|
| SELECT 1 / 0;
< error
| 22012 division by zero
> query
| SELECT count(*) FROM exec_t WHERE id = 20
< result
< column "count" "INT8"
< row
: int64 1
> exec
| DO $$ BEGIN INSERT INTO exec_t VALUES (30, 'do'); COMMIT; INSERT INTO exec_t VALUES (31, 'do2'); END $$;
< ok 0 -
> exec
|
| 				CREATE ROLE exec_tenant_owner NOLOGIN;
| 				CREATE ROLE exec_reader NOLOGIN;
| 				GRANT CREATE ON SCHEMA public TO exec_tenant_owner;
| 			
< ok 0 -
> exec
| ALTER DATABASE "omni_testexecute" OWNER TO exec_tenant_owner;
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
| exec_tenant_owner
> begin
< ok
> exec
| SET LOCAL ROLE 'exec_tenant_owner'
< ok 0 -
> exec
| CREATE TABLE tenant_owned (id int);
< ok 0 -
> exec
| SET LOCAL ROLE NONE;
| GRANT SELECT ON tenant_owned TO exec_reader;SET LOCAL ROLE 'exec_tenant_owner';
< ok 0 -
> commit
< ok
> exec
| SET SESSION ROLE 'exec_tenant_owner'
< ok 0 -
> query
| SELECT tableowner FROM pg_tables WHERE schemaname = 'public' AND tablename = 'tenant_owned'
< result
< column "tableowner" "NAME"
< row
: string
| exec_tenant_owner
