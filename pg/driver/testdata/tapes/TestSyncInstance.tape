-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12
-- pin: pg12 postgres:12.22-alpine sha256:7c8f4870583184ebadf7f17a6513620aac5f365a7938dc6a6911c1d5df2f481a REL_12_22
-- pin: pg13 postgres:13.23-alpine sha256:fb9065b6e3e213bdc07edd372a5b2a26245840b7fb65d1fd8b6700106d51805c REL_13_23
-- pin: pg14 postgres:14.25-alpine sha256:14078c4aa0db8e2e8492866138cafd796fc39afd5ecb5b1c6f26a912b07529f5 REL_14_25
-- pin: pg15 postgres:15.20-alpine sha256:5cb339462bde5e6b617acedf6c2a3008f3aadc44a3c0507dd22226ee1279aee8 REL_15_20
-- pin: pg16 postgres:16.16-alpine sha256:0adbf25a9e96d38823bb813f3ac0ed41fadb665c88f1df7fe6dbfb060cbdc106 REL_16_16
-- pin: pg18 postgres:18.7-alpine sha256:75dcac3dfb1cd01ca3536910bbc84128759362f8b118bf8c55b38a94f68f3b73 REL_18_7

== TestSyncInstance/default default#1
> exec
|
| CREATE ROLE sync_nologin NOLOGIN;
| CREATE ROLE sync_admin SUPERUSER NOINHERIT CREATEROLE CREATEDB REPLICATION BYPASSRLS LOGIN VALID UNTIL '2030-06-01 12:00:00+00';
| CREATE ROLE sync_plain LOGIN;
< ok 0 -
> query
| SHOW server_version_num
< result
< column "server_version_num" "TEXT"
< row
: string
| 170012
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
| postgres
: bool true
: bool true
: bool true
: bool true
: bool true
: bool true
: null
: bool true
< row
: string
| sync_admin
: bool true
: bool false
: bool true
: bool true
: bool true
: bool true
: time 2030-06-01T12:00:00Z
: bool true
< row
: string
| sync_nologin
: bool false
: bool true
: bool false
: bool false
: bool false
: bool false
: null
: bool false
< row
: string
| sync_plain
: bool false
: bool true
: bool false
: bool false
: bool true
: bool false
: null
: bool false
> query
| SELECT datname, pg_encoding_to_char(encoding), datcollate, pg_catalog.pg_get_userbyid(datdba) as db_owner FROM pg_database ORDER BY datname;
< result
< column "datname" "NAME"
< column "pg_encoding_to_char" "NAME"
< column "datcollate" "TEXT"
< column "db_owner" "NAME"
< row
: string
| omni_testsyncinstance
: string
| UTF8
: string
| en_US.utf8
: string
| postgres
< row
: string
| postgres
: string
| UTF8
: string
| en_US.utf8
: string
| postgres
< row
: string
| template0
: string
| UTF8
: string
| en_US.utf8
: string
| postgres
< row
: string
| template1
: string
| UTF8
: string
| en_US.utf8
: string
| postgres

== TestSyncInstance/pg12 pg12#1
> exec
|
| CREATE ROLE sync_nologin NOLOGIN;
| CREATE ROLE sync_admin SUPERUSER NOINHERIT CREATEROLE CREATEDB REPLICATION BYPASSRLS LOGIN VALID UNTIL '2030-06-01 12:00:00+00';
| CREATE ROLE sync_plain LOGIN;
< ok 0 -
> query
| SHOW server_version_num
< result
< column "server_version_num" "TEXT"
< row
: string
| 120022
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
| postgres
: bool true
: bool true
: bool true
: bool true
: bool true
: bool true
: null
: bool true
< row
: string
| sync_admin
: bool true
: bool false
: bool true
: bool true
: bool true
: bool true
: time 2030-06-01T12:00:00Z
: bool true
< row
: string
| sync_nologin
: bool false
: bool true
: bool false
: bool false
: bool false
: bool false
: null
: bool false
< row
: string
| sync_plain
: bool false
: bool true
: bool false
: bool false
: bool true
: bool false
: null
: bool false
> query
| SELECT datname, pg_encoding_to_char(encoding), datcollate, pg_catalog.pg_get_userbyid(datdba) as db_owner FROM pg_database ORDER BY datname;
< result
< column "datname" "NAME"
< column "pg_encoding_to_char" "NAME"
< column "datcollate" "NAME"
< column "db_owner" "NAME"
< row
: string
| omni_testsyncinstance
: string
| UTF8
: string
| en_US.utf8
: string
| postgres
< row
: string
| postgres
: string
| UTF8
: string
| en_US.utf8
: string
| postgres
< row
: string
| template0
: string
| UTF8
: string
| en_US.utf8
: string
| postgres
< row
: string
| template1
: string
| UTF8
: string
| en_US.utf8
: string
| postgres

== TestSyncInstance/pg13 pg13#1
> exec
|
| CREATE ROLE sync_nologin NOLOGIN;
| CREATE ROLE sync_admin SUPERUSER NOINHERIT CREATEROLE CREATEDB REPLICATION BYPASSRLS LOGIN VALID UNTIL '2030-06-01 12:00:00+00';
| CREATE ROLE sync_plain LOGIN;
< ok 0 -
> query
| SHOW server_version_num
< result
< column "server_version_num" "TEXT"
< row
: string
| 130023
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
| postgres
: bool true
: bool true
: bool true
: bool true
: bool true
: bool true
: null
: bool true
< row
: string
| sync_admin
: bool true
: bool false
: bool true
: bool true
: bool true
: bool true
: time 2030-06-01T12:00:00Z
: bool true
< row
: string
| sync_nologin
: bool false
: bool true
: bool false
: bool false
: bool false
: bool false
: null
: bool false
< row
: string
| sync_plain
: bool false
: bool true
: bool false
: bool false
: bool true
: bool false
: null
: bool false
> query
| SELECT datname, pg_encoding_to_char(encoding), datcollate, pg_catalog.pg_get_userbyid(datdba) as db_owner FROM pg_database ORDER BY datname;
< result
< column "datname" "NAME"
< column "pg_encoding_to_char" "NAME"
< column "datcollate" "NAME"
< column "db_owner" "NAME"
< row
: string
| omni_testsyncinstance
: string
| UTF8
: string
| en_US.utf8
: string
| postgres
< row
: string
| postgres
: string
| UTF8
: string
| en_US.utf8
: string
| postgres
< row
: string
| template0
: string
| UTF8
: string
| en_US.utf8
: string
| postgres
< row
: string
| template1
: string
| UTF8
: string
| en_US.utf8
: string
| postgres

== TestSyncInstance/pg14 pg14#1
> exec
|
| CREATE ROLE sync_nologin NOLOGIN;
| CREATE ROLE sync_admin SUPERUSER NOINHERIT CREATEROLE CREATEDB REPLICATION BYPASSRLS LOGIN VALID UNTIL '2030-06-01 12:00:00+00';
| CREATE ROLE sync_plain LOGIN;
< ok 0 -
> query
| SHOW server_version_num
< result
< column "server_version_num" "TEXT"
< row
: string
| 140025
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
| postgres
: bool true
: bool true
: bool true
: bool true
: bool true
: bool true
: null
: bool true
< row
: string
| sync_admin
: bool true
: bool false
: bool true
: bool true
: bool true
: bool true
: time 2030-06-01T12:00:00Z
: bool true
< row
: string
| sync_nologin
: bool false
: bool true
: bool false
: bool false
: bool false
: bool false
: null
: bool false
< row
: string
| sync_plain
: bool false
: bool true
: bool false
: bool false
: bool true
: bool false
: null
: bool false
> query
| SELECT datname, pg_encoding_to_char(encoding), datcollate, pg_catalog.pg_get_userbyid(datdba) as db_owner FROM pg_database ORDER BY datname;
< result
< column "datname" "NAME"
< column "pg_encoding_to_char" "NAME"
< column "datcollate" "NAME"
< column "db_owner" "NAME"
< row
: string
| omni_testsyncinstance
: string
| UTF8
: string
| en_US.utf8
: string
| postgres
< row
: string
| postgres
: string
| UTF8
: string
| en_US.utf8
: string
| postgres
< row
: string
| template0
: string
| UTF8
: string
| en_US.utf8
: string
| postgres
< row
: string
| template1
: string
| UTF8
: string
| en_US.utf8
: string
| postgres

== TestSyncInstance/pg15 pg15#1
> exec
|
| CREATE ROLE sync_nologin NOLOGIN;
| CREATE ROLE sync_admin SUPERUSER NOINHERIT CREATEROLE CREATEDB REPLICATION BYPASSRLS LOGIN VALID UNTIL '2030-06-01 12:00:00+00';
| CREATE ROLE sync_plain LOGIN;
< ok 0 -
> query
| SHOW server_version_num
< result
< column "server_version_num" "TEXT"
< row
: string
| 150020
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
| postgres
: bool true
: bool true
: bool true
: bool true
: bool true
: bool true
: null
: bool true
< row
: string
| sync_admin
: bool true
: bool false
: bool true
: bool true
: bool true
: bool true
: time 2030-06-01T12:00:00Z
: bool true
< row
: string
| sync_nologin
: bool false
: bool true
: bool false
: bool false
: bool false
: bool false
: null
: bool false
< row
: string
| sync_plain
: bool false
: bool true
: bool false
: bool false
: bool true
: bool false
: null
: bool false
> query
| SELECT datname, pg_encoding_to_char(encoding), datcollate, pg_catalog.pg_get_userbyid(datdba) as db_owner FROM pg_database ORDER BY datname;
< result
< column "datname" "NAME"
< column "pg_encoding_to_char" "NAME"
< column "datcollate" "TEXT"
< column "db_owner" "NAME"
< row
: string
| omni_testsyncinstance
: string
| UTF8
: string
| en_US.utf8
: string
| postgres
< row
: string
| postgres
: string
| UTF8
: string
| en_US.utf8
: string
| postgres
< row
: string
| template0
: string
| UTF8
: string
| en_US.utf8
: string
| postgres
< row
: string
| template1
: string
| UTF8
: string
| en_US.utf8
: string
| postgres

== TestSyncInstance/pg16 pg16#1
> exec
|
| CREATE ROLE sync_nologin NOLOGIN;
| CREATE ROLE sync_admin SUPERUSER NOINHERIT CREATEROLE CREATEDB REPLICATION BYPASSRLS LOGIN VALID UNTIL '2030-06-01 12:00:00+00';
| CREATE ROLE sync_plain LOGIN;
< ok 0 -
> query
| SHOW server_version_num
< result
< column "server_version_num" "TEXT"
< row
: string
| 160016
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
| postgres
: bool true
: bool true
: bool true
: bool true
: bool true
: bool true
: null
: bool true
< row
: string
| sync_admin
: bool true
: bool false
: bool true
: bool true
: bool true
: bool true
: time 2030-06-01T12:00:00Z
: bool true
< row
: string
| sync_nologin
: bool false
: bool true
: bool false
: bool false
: bool false
: bool false
: null
: bool false
< row
: string
| sync_plain
: bool false
: bool true
: bool false
: bool false
: bool true
: bool false
: null
: bool false
> query
| SELECT datname, pg_encoding_to_char(encoding), datcollate, pg_catalog.pg_get_userbyid(datdba) as db_owner FROM pg_database ORDER BY datname;
< result
< column "datname" "NAME"
< column "pg_encoding_to_char" "NAME"
< column "datcollate" "TEXT"
< column "db_owner" "NAME"
< row
: string
| omni_testsyncinstance
: string
| UTF8
: string
| en_US.utf8
: string
| postgres
< row
: string
| postgres
: string
| UTF8
: string
| en_US.utf8
: string
| postgres
< row
: string
| template0
: string
| UTF8
: string
| en_US.utf8
: string
| postgres
< row
: string
| template1
: string
| UTF8
: string
| en_US.utf8
: string
| postgres

== TestSyncInstance/pg18 pg18#1
> exec
|
| CREATE ROLE sync_nologin NOLOGIN;
| CREATE ROLE sync_admin SUPERUSER NOINHERIT CREATEROLE CREATEDB REPLICATION BYPASSRLS LOGIN VALID UNTIL '2030-06-01 12:00:00+00';
| CREATE ROLE sync_plain LOGIN;
< ok 0 -
> query
| SHOW server_version_num
< result
< column "server_version_num" "TEXT"
< row
: string
| 180007
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
| postgres
: bool true
: bool true
: bool true
: bool true
: bool true
: bool true
: null
: bool true
< row
: string
| sync_admin
: bool true
: bool false
: bool true
: bool true
: bool true
: bool true
: time 2030-06-01T12:00:00Z
: bool true
< row
: string
| sync_nologin
: bool false
: bool true
: bool false
: bool false
: bool false
: bool false
: null
: bool false
< row
: string
| sync_plain
: bool false
: bool true
: bool false
: bool false
: bool true
: bool false
: null
: bool false
> query
| SELECT datname, pg_encoding_to_char(encoding), datcollate, pg_catalog.pg_get_userbyid(datdba) as db_owner FROM pg_database ORDER BY datname;
< result
< column "datname" "NAME"
< column "pg_encoding_to_char" "NAME"
< column "datcollate" "TEXT"
< column "db_owner" "NAME"
< row
: string
| omni_testsyncinstance
: string
| UTF8
: string
| en_US.utf8
: string
| postgres
< row
: string
| postgres
: string
| UTF8
: string
| en_US.utf8
: string
| postgres
< row
: string
| template0
: string
| UTF8
: string
| en_US.utf8
: string
| postgres
< row
: string
| template1
: string
| UTF8
: string
| en_US.utf8
: string
| postgres
