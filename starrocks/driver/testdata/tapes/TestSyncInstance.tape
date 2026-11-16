-- pin: default starrocks/allin1-ubuntu:3.4.10 sha256:8bc803ea7a5bb490c77117a456b32b792ccae698b1f0e47dfcafdac41542566d 3.4.10

== TestSyncInstance default#1
> query
| SELECT VERSION()
< result
< column "version()" "VARCHAR"
< row
: bytes
| 8.0.33
> query
| SHOW VARIABLES LIKE 'lower_case_table_names'
< result
< column "Variable_name" "VARCHAR"
< column "Value" "VARCHAR"
< row
: bytes
| lower_case_table_names
: bytes
| 0
> query
|
| 		SELECT
| 			SCHEMA_NAME,
| 			DEFAULT_CHARACTER_SET_NAME,
| 			DEFAULT_COLLATION_NAME
| 		FROM information_schema.SCHEMATA
| 		WHERE LOWER(SCHEMA_NAME) NOT IN ('information_schema', '_statistics_', 'sys')
| 		ORDER BY SCHEMA_NAME
< result
< column "SCHEMA_NAME" "VARCHAR"
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| sync_inst_db
: bytes
| utf8
: bytes
| utf8_general_ci
