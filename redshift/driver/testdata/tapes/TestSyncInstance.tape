-- pin: driver aws:redshift:1.0.539507 sha256:afc4041fe6e35ce273970308e430eb8245ea9c4bbf597c497bb8eb8bb00810d2 2026-11-17

== TestSyncInstance driver#1
> query
| SELECT version()
< result
< column "version" "VARCHAR" "string"
< row
: string
| PostgreSQL 8.0.2 on i686-pc-linux-gnu, compiled by GCC gcc (GCC) 3.4.2 20041017 (Red Hat 3.4.2-6.fc3), Redshift 1.0.539507
> query
|
| 	SELECT
| 		u.usename AS rolename,
| 		u.usesuper AS rolsuper,
| 		true AS rolinherit,
| 		false AS rolcreaterole,
| 		u.usecreatedb AS rolcreatedb,
| 		true AS rolcanlogin,
| 		-1 AS rolconnlimit,
| 		u.valuntil as rolvaliduntil
| 	FROM pg_catalog.pg_user u
| 	ORDER BY 1;
| 	
< result
< column "rolename" "NAME" "string"
< column "rolsuper" "BOOL" "bool"
< column "rolinherit" "BOOL" "bool"
< column "rolcreaterole" "BOOL" "bool"
< column "rolcreatedb" "BOOL" "bool"
< column "rolcanlogin" "BOOL" "bool"
< column "rolconnlimit" "INT4" "int32"
< column "rolvaliduntil" "702" "string"
< row
: string
| omni_verify_admin
: bool true
: bool true
: bool false
: bool true
: bool true
: int64 -1
: null
< row
: string
| rdsdb
: bool true
: bool true
: bool false
: bool true
: bool true
: int64 -1
: null
> query
| SELECT consumer_database FROM SVV_DATASHARES WHERE share_type = 'INBOUND' ORDER BY consumer_database;
< result
< column "consumer_database" "VARCHAR" "string"
> query
|
| 		SELECT datname,
| 		pg_encoding_to_char(encoding)
| 		FROM pg_database
| 		ORDER BY datname;
| 	
< result
< column "datname" "NAME" "string"
< column "pg_encoding_to_char" "NAME" "string"
< row
: string
| awsdatacatalog
: string
| UNICODE
< row
: string
| dev
: string
| UNICODE
< row
: string
| padb_harvest
: string
| UNICODE
< row
: string
| sys:internal
: string
| UNICODE
< row
: string
| template0
: string
| UNICODE
< row
: string
| template1
: string
| UNICODE
