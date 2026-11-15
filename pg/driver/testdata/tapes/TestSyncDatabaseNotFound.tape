-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestSyncDatabaseNotFound default#1
> query
| SELECT datname, pg_encoding_to_char(encoding), datcollate, pg_catalog.pg_get_userbyid(datdba) as db_owner FROM pg_database ORDER BY datname;
< result
< column "datname" "NAME"
< column "pg_encoding_to_char" "NAME"
< column "datcollate" "TEXT"
< column "db_owner" "NAME"
< row
: string
| omni_testsyncdatabasenotfound
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
