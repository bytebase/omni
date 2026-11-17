-- pin: default cockroachdb/cockroach:v25.2.3 sha256:4a412731d1cfc5d6dea7431b7ce383b502966c9d589163d63db9e1ff7b51495a v25.2.3

== TestSyncInstance default#1
> exec
|
| CREATE DATABASE sync_extra;
| CREATE ROLE sync_reader;
| CREATE USER sync_login;
| CREATE ROLE sync_admin WITH CREATEROLE CREATEDB LOGIN;
< ok 0 -
> query
| SHOW CLUSTER SETTING version;
< result
< column "version" "TEXT"
< row
: string
| 25.2
> query
|
| 		SELECT r.rolname, r.rolsuper, r.rolinherit, r.rolcreaterole, r.rolcreatedb, r.rolcanlogin, r.rolreplication, r.rolvaliduntil, r.rolbypassrls
| 		FROM pg_catalog.pg_roles r
| 		WHERE r.rolname !~ '^pg_'
| 		ORDER BY r.rolname;
| 	
< result
< column "rolname" "NAME"
< column "rolsuper" "BOOL"
< column "rolinherit" "BOOL"
< column "rolcreaterole" "BOOL"
< column "rolcreatedb" "BOOL"
< column "rolcanlogin" "BOOL"
< column "rolreplication" "BOOL"
< column "rolvaliduntil" "TIMESTAMPTZ"
< column "rolbypassrls" "BOOL"
< row
: string
| admin
: bool true
: bool true
: bool true
: bool true
: bool true
: bool false
: null
: bool false
< row
: string
| node
: bool true
: bool true
: bool true
: bool true
: bool false
: bool false
: null
: bool false
< row
: string
| root
: bool true
: bool true
: bool true
: bool true
: bool true
: bool false
: null
: bool false
< row
: string
| sync_admin
: bool false
: bool true
: bool true
: bool true
: bool true
: bool false
: null
: bool false
< row
: string
| sync_login
: bool false
: bool true
: bool false
: bool false
: bool true
: bool false
: null
: bool false
< row
: string
| sync_reader
: bool false
: bool true
: bool false
: bool false
: bool false
: bool false
: null
: bool false
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
| omni_testsyncinstance
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
| sync_extra
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
