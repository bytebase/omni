-- pin: default dyrnq/doris:3.1.4 sha256:88f99819c0877567b77d78cc8367f0e2dd4b6c1602fe17c6f4b1e332d8e315dd 3.1.4-rc02

== TestSyncInstance default#1
> query
| SELECT VERSION()
< result
< column "VERSION()" "TEXT"
< row
: bytes
| 5.7.99
> query
| SHOW VARIABLES LIKE 'lower_case_table_names'
< result
< column "Variable_name" "CHAR"
< column "Value" "CHAR"
< column "Default_Value" "CHAR"
< column "Changed" "CHAR"
< row
: bytes
| lower_case_table_names
: bytes
| 0
: bytes
| 0
: bytes
| 0
> query
|
| 		SELECT
| 			SCHEMA_NAME,
| 			DEFAULT_CHARACTER_SET_NAME,
| 			DEFAULT_COLLATION_NAME
| 		FROM information_schema.SCHEMATA
| 		WHERE LOWER(SCHEMA_NAME) NOT IN ('information_schema', '_statistics_')
| 		ORDER BY SCHEMA_NAME
< result
< column "SCHEMA_NAME" "CHAR"
< column "DEFAULT_CHARACTER_SET_NAME" "CHAR"
< column "DEFAULT_COLLATION_NAME" "CHAR"
< row
: bytes
| __internal_schema
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_bin
< row
: bytes
| mysql
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_bin
< row
: bytes
| sync_inst_db
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_bin
