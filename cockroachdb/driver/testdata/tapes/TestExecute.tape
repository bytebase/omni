-- pin: default cockroachdb/cockroach:v25.2.3 sha256:4a412731d1cfc5d6dea7431b7ce383b502966c9d589163d63db9e1ff7b51495a v25.2.3

== TestExecute default#1
> exec
|
| CREATE SCHEMA app;
| CREATE TABLE app.account (id INT8 PRIMARY KEY, v INT8 NOT NULL);
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
| root
> begin
< ok
> exec
| SET LOCAL ROLE 'root'
< ok 0 -
> exec
| INSERT INTO app.account VALUES (1, 10);
< ok 1 -
> exec
| INSERT INTO app.account VALUES (2, 20);
< ok 1 -
> exec
| UPDATE app.account SET v = v + 1 WHERE id = 1;
< ok 1 -
> commit
< ok
> exec
| SET SESSION ROLE 'root'
< ok 0 -
> query
| SELECT sum(v) FROM app.account;
< result
< column "sum" "NUMERIC"
< row
: string
| 31
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
| root
> exec
| SET SESSION ROLE 'root'
< ok 0 -
> exec
| INSERT INTO app.account VALUES (3, 30);
< ok 1 -
> exec
| INSERT INTO app.account VALUES (4, 40);
< ok 1 -
> query
| SELECT count(*) FROM app.account WHERE id IN (3, 4);
< result
< column "count" "INT8"
< row
: int64 2
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
| root
> begin
< ok
> exec
| SET LOCAL ROLE 'root'
< ok 0 -
> exec
| INSERT INTO app.account VALUES (5, 50);
< ok 1 -
> exec
| INSERT INTO app.account VALUES (1, 99);
< error
| 23505 duplicate key value violates unique constraint "account_pkey"
| detail: Key (id)=(1) already exists.
| constraint: account_pkey
> rollback
< ok
> query
| SELECT count(*) FROM app.account WHERE id = 5;
< result
< column "count" "INT8"
< row
: int64 0
> query
| SELECT datname, pg_encoding_to_char(encoding), datcollate FROM pg_database ORDER BY datname;
< result
< column "datname" "NAME"
< column "pg_encoding_to_char" "TEXT"
< column "datcollate" "TEXT"
< row
: string
| defaultdb
: string
| UTF8
: string
| en_US.utf8
< row
: string
| omni_testexecute
: string
| UTF8
: string
| en_US.utf8
< row
: string
| postgres
: string
| UTF8
: string
| en_US.utf8
< row
: string
| system
: string
| UTF8
: string
| en_US.utf8
> exec
| CREATE DATABASE omni_execute_created;
< ok 0 -
> query
| SELECT count(*) FROM pg_database WHERE datname = 'omni_execute_created';
< result
< column "count" "INT8"
< row
: int64 1
> query
| SELECT datname, pg_encoding_to_char(encoding), datcollate FROM pg_database ORDER BY datname;
< result
< column "datname" "NAME"
< column "pg_encoding_to_char" "TEXT"
< column "datcollate" "TEXT"
< row
: string
| defaultdb
: string
| UTF8
: string
| en_US.utf8
< row
: string
| omni_execute_created
: string
| UTF8
: string
| en_US.utf8
< row
: string
| omni_testexecute
: string
| UTF8
: string
| en_US.utf8
< row
: string
| postgres
: string
| UTF8
: string
| en_US.utf8
< row
: string
| system
: string
| UTF8
: string
| en_US.utf8
