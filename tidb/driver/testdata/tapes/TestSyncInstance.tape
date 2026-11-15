-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestSyncInstance default#1
> query
| SELECT VERSION()
< result
< column "VERSION()" "VARCHAR"
< row
: bytes
| 8.0.11-TiDB-v8.5.5
> query
| SHOW VARIABLES LIKE 'lower_case_table_names'
< result
< column "Variable_name" "VARCHAR"
< column "Value" "VARCHAR"
< row
: bytes
| lower_case_table_names
: bytes
| 2
> query
|
| 	SELECT
| 		user,
| 		host
| 	FROM mysql.user
| 	WHERE user <> ''
| 		AND user NOT LIKE 'mysql.%'
| 	ORDER BY user, host
| 	
< result
< column "user" "CHAR"
< column "host" "CHAR"
< row
: bytes
| root
: bytes
| %
< row
: bytes
| sync_reader
: bytes
| %
> query
| SHOW GRANTS FOR 'root'@'%'
< result
< column "Grants for root@%" "VARCHAR"
< row
: bytes
| GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' WITH GRANT OPTION
> query
| SHOW GRANTS FOR 'sync_reader'@'%'
< result
< column "Grants for sync_reader@%" "VARCHAR"
< row
: bytes
| GRANT USAGE ON *.* TO 'sync_reader'@'%'
< row
: bytes
| GRANT SELECT ON `sync_inst_db`.* TO 'sync_reader'@'%'
> query
|
| 		SELECT
| 			SCHEMA_NAME,
| 			DEFAULT_CHARACTER_SET_NAME,
| 			DEFAULT_COLLATION_NAME
| 		FROM information_schema.SCHEMATA
| 		WHERE LOWER(SCHEMA_NAME) NOT IN ('information_schema', 'metrics_schema', 'mysql', 'performance_schema', 'sys')
| 		ORDER BY SCHEMA_NAME
< result
< column "SCHEMA_NAME" "VARCHAR"
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| sync_inst_db
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
< row
: bytes
| test
: bytes
| utf8mb4
: bytes
| utf8mb4_bin
