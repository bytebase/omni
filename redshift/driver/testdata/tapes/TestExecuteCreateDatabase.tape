-- pin: driver aws:redshift:1.0.539507 sha256:afc4041fe6e35ce273970308e430eb8245ea9c4bbf597c497bb8eb8bb00810d2 2026-11-17

== TestExecuteCreateDatabase driver#1
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
> exec
| CREATE DATABASE omni_create_probe
< ok 0 -
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
| omni_create_probe
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
| omni_create_probe
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
| omni_create_probe
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
| omni_create_probe
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
> exec
| DROP DATABASE omni_create_probe
< ok 0 -
