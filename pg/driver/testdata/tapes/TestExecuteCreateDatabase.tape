-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12
-- pin: pg12 postgres:12.22-alpine sha256:7c8f4870583184ebadf7f17a6513620aac5f365a7938dc6a6911c1d5df2f481a REL_12_22

== TestExecuteCreateDatabase/default default#1
> exec
| CREATE ROLE created_db_owner NOLOGIN;
< ok 0 -
> query
| SELECT datname, pg_encoding_to_char(encoding), datcollate, pg_catalog.pg_get_userbyid(datdba) as db_owner FROM pg_database ORDER BY datname;
< result
< column "datname" "NAME"
< column "pg_encoding_to_char" "NAME"
< column "datcollate" "TEXT"
< column "db_owner" "NAME"
< row
: string
| omni_testexecutecreatedatabase
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
> exec
| CREATE DATABASE "omni_new_db"
| ENCODING "UTF8";
< ok 0 -
> exec
|
| ALTER DATABASE "omni_new_db" OWNER TO "created_db_owner";
< ok 0 -
> query
| SELECT pg_catalog.pg_get_userbyid(datdba) FROM pg_database WHERE datname = 'omni_new_db'
< result
< column "pg_get_userbyid" "NAME"
< row
: string
| created_db_owner
> query
| SELECT datname, pg_encoding_to_char(encoding), datcollate, pg_catalog.pg_get_userbyid(datdba) as db_owner FROM pg_database ORDER BY datname;
< result
< column "datname" "NAME"
< column "pg_encoding_to_char" "NAME"
< column "datcollate" "TEXT"
< column "db_owner" "NAME"
< row
: string
| omni_new_db
: string
| UTF8
: string
| en_US.utf8
: string
| created_db_owner
< row
: string
| omni_testexecutecreatedatabase
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

== TestExecuteCreateDatabase/pg12 pg12#1
> exec
| CREATE ROLE created_db_owner NOLOGIN;
< ok 0 -
> query
| SELECT datname, pg_encoding_to_char(encoding), datcollate, pg_catalog.pg_get_userbyid(datdba) as db_owner FROM pg_database ORDER BY datname;
< result
< column "datname" "NAME"
< column "pg_encoding_to_char" "NAME"
< column "datcollate" "NAME"
< column "db_owner" "NAME"
< row
: string
| omni_testexecutecreatedatabase
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
> exec
| CREATE DATABASE "omni_new_db"
| ENCODING "UTF8";
< ok 0 -
> exec
|
| ALTER DATABASE "omni_new_db" OWNER TO "created_db_owner";
< ok 0 -
> query
| SELECT pg_catalog.pg_get_userbyid(datdba) FROM pg_database WHERE datname = 'omni_new_db'
< result
< column "pg_get_userbyid" "NAME"
< row
: string
| created_db_owner
> query
| SELECT datname, pg_encoding_to_char(encoding), datcollate, pg_catalog.pg_get_userbyid(datdba) as db_owner FROM pg_database ORDER BY datname;
< result
< column "datname" "NAME"
< column "pg_encoding_to_char" "NAME"
< column "datcollate" "NAME"
< column "db_owner" "NAME"
< row
: string
| omni_new_db
: string
| UTF8
: string
| en_US.utf8
: string
| created_db_owner
< row
: string
| omni_testexecutecreatedatabase
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
