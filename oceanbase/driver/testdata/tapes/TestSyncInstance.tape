-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestSyncInstance default#1
> query
| SELECT VERSION()
< result
< column "VERSION()" "VARCHAR"
< row
: bytes
| 5.7.25-OceanBase_CE-v4.3.5.6
> query
| SHOW VARIABLES LIKE 'lower_case_table_names'
< result
< column "Variable_name" "VARCHAR" "string"
< column "Value" "VARCHAR"
< row
: bytes
| lower_case_table_names
: bytes
| 1
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
< column "user" "VARCHAR" "string"
< column "host" "VARCHAR" "string"
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
< column "Grants for root@%" "VARCHAR" "string"
< row
: bytes
| GRANT ALL PRIVILEGES ON *.* TO 'root' WITH GRANT OPTION
< row
: bytes
| GRANT ALL PRIVILEGES ON `oceanbase`.* TO 'root'
< row
: bytes
| GRANT ALL PRIVILEGES ON `__recyclebin`.* TO 'root'
< row
: bytes
| GRANT ALL PRIVILEGES ON `mysql`.* TO 'root'
< row
: bytes
| GRANT ALL PRIVILEGES ON `test`.* TO 'root'
< row
: bytes
| GRANT ALL PRIVILEGES ON `information_schema`.* TO 'root'
< row
: bytes
| GRANT ALL PRIVILEGES ON `__public`.* TO 'root'
> query
| SHOW GRANTS FOR 'sync_reader'@'%'
< result
< column "Grants for sync_reader@%" "VARCHAR" "string"
< row
: bytes
| GRANT USAGE ON *.* TO 'sync_reader'
< row
: bytes
| GRANT SELECT ON `sync_inst_db`.* TO 'sync_reader'
> query
|
| 		SELECT
| 			SCHEMA_NAME,
| 			DEFAULT_CHARACTER_SET_NAME,
| 			DEFAULT_COLLATION_NAME
| 		FROM information_schema.SCHEMATA
| 		WHERE LOWER(SCHEMA_NAME) NOT IN ('LBACSYS', 'ORAAUDITOR', 'SYS', '__public', 'information_schema', 'mysql', 'oceanbase', 'performance_schema', 'sys')
| 		ORDER BY SCHEMA_NAME
< result
< column "SCHEMA_NAME" "VARCHAR" "string"
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR" "string"
< column "DEFAULT_COLLATION_NAME" "VARCHAR" "string"
< row
: bytes
| sync_inst_db
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
< row
: bytes
| test
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci
