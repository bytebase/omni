-- pin: default mariadb:11.8.8 sha256:24e76fcec8c003a0362d0dd53f4806e7e79458d7fdeaf47437760e19496f5a9c mariadb-11.8.8
-- pin: mariadb1011 mariadb:10.11.19 sha256:7d829180b5b2a9ba74f6242876ab401a0d7dd338a807e522ef92bc9698760f20 mariadb-10.11.19
-- pin: mariadb107 mariadb:10.7.8 sha256:9a48ac9f196f3d4fd6fea2cab59a49df9e7ca459bf14b2f7b85a0e38a5454571 mariadb-10.7.8
-- pin: mariadb114 mariadb:11.4.13 sha256:5e73baf8179b29f115101bb102763a47a017361f3e9792daccba6646c3ab1532 mariadb-11.4.13

== TestSyncInstance/default default#1
> query
| SELECT VERSION()
< result
< column "VERSION()" "VARCHAR"
< row
: bytes
| 11.8.8-MariaDB-ubu2404
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
| 	SELECT
| 		user,
| 		host
| 	FROM mysql.user
| 	WHERE user <> ''
| 		AND user NOT LIKE 'mysql.%'
| 	ORDER BY user, host
| 	
< result
< column "User" "CHAR"
< column "Host" "CHAR"
< row
: bytes
| mariadb.sys
: bytes
| localhost
< row
: bytes
| root
: bytes
| %
< row
: bytes
| root
: bytes
| localhost
< row
: bytes
| sync_reader
: bytes
| %
> query
| SHOW GRANTS FOR 'mariadb.sys'@'localhost'
< result
< column "Grants for mariadb.sys@localhost" "VARCHAR"
< row
: bytes
| GRANT USAGE ON *.* TO `mariadb.sys`@`localhost`
< row
: bytes
| GRANT SELECT, DELETE ON `mysql`.`global_priv` TO `mariadb.sys`@`localhost`
> query
| SHOW GRANTS FOR 'root'@'%'
< result
< column "Grants for root@%" "VARCHAR"
< row
: bytes
| GRANT ALL PRIVILEGES ON *.* TO `root`@`%` IDENTIFIED BY PASSWORD '*94BDCEBE19083CE2A1F959FD02F964C7AF4CFC29' WITH GRANT OPTION
< row
: bytes
| GRANT PROXY ON ''@'%' TO 'root'@'%' WITH GRANT OPTION
> query
| SHOW GRANTS FOR 'root'@'localhost'
< result
< column "Grants for root@localhost" "VARCHAR"
< row
: bytes
| GRANT ALL PRIVILEGES ON *.* TO `root`@`localhost` IDENTIFIED BY PASSWORD '*94BDCEBE19083CE2A1F959FD02F964C7AF4CFC29' WITH GRANT OPTION
< row
: bytes
| GRANT PROXY ON ''@'%' TO 'root'@'localhost' WITH GRANT OPTION
> query
| SHOW GRANTS FOR 'sync_reader'@'%'
< result
< column "Grants for sync_reader@%" "VARCHAR"
< row
: bytes
| GRANT USAGE ON *.* TO `sync_reader`@`%` IDENTIFIED BY PASSWORD '*0EB931DD0E0457476ABFBD89F6B34EA17B7C84C5'
< row
: bytes
| GRANT SELECT ON `sync_inst_db`.* TO `sync_reader`@`%`
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
< column "SCHEMA_NAME" "VARCHAR"
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| sync_inst_db
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
< row
: bytes
| test
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci

== TestSyncInstance/mariadb1011 mariadb1011#1
> query
| SELECT VERSION()
< result
< column "VERSION()" "VARCHAR"
< row
: bytes
| 10.11.19-MariaDB-ubu2204
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
| 	SELECT
| 		user,
| 		host
| 	FROM mysql.user
| 	WHERE user <> ''
| 		AND user NOT LIKE 'mysql.%'
| 	ORDER BY user, host
| 	
< result
< column "User" "CHAR"
< column "Host" "CHAR"
< row
: bytes
| mariadb.sys
: bytes
| localhost
< row
: bytes
| root
: bytes
| %
< row
: bytes
| root
: bytes
| localhost
< row
: bytes
| sync_reader
: bytes
| %
> query
| SHOW GRANTS FOR 'mariadb.sys'@'localhost'
< result
< column "Grants for mariadb.sys@localhost" "VARCHAR"
< row
: bytes
| GRANT USAGE ON *.* TO `mariadb.sys`@`localhost`
< row
: bytes
| GRANT SELECT, DELETE ON `mysql`.`global_priv` TO `mariadb.sys`@`localhost`
> query
| SHOW GRANTS FOR 'root'@'%'
< result
< column "Grants for root@%" "VARCHAR"
< row
: bytes
| GRANT ALL PRIVILEGES ON *.* TO `root`@`%` IDENTIFIED BY PASSWORD '*94BDCEBE19083CE2A1F959FD02F964C7AF4CFC29' WITH GRANT OPTION
< row
: bytes
| GRANT PROXY ON ``@`%` TO `root`@`%` WITH GRANT OPTION
> query
| SHOW GRANTS FOR 'root'@'localhost'
< result
< column "Grants for root@localhost" "VARCHAR"
< row
: bytes
| GRANT ALL PRIVILEGES ON *.* TO `root`@`localhost` IDENTIFIED BY PASSWORD '*94BDCEBE19083CE2A1F959FD02F964C7AF4CFC29' WITH GRANT OPTION
< row
: bytes
| GRANT PROXY ON ``@`%` TO `root`@`localhost` WITH GRANT OPTION
> query
| SHOW GRANTS FOR 'sync_reader'@'%'
< result
< column "Grants for sync_reader@%" "VARCHAR"
< row
: bytes
| GRANT USAGE ON *.* TO `sync_reader`@`%` IDENTIFIED BY PASSWORD '*0EB931DD0E0457476ABFBD89F6B34EA17B7C84C5'
< row
: bytes
| GRANT SELECT ON `sync_inst_db`.* TO `sync_reader`@`%`
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
< column "SCHEMA_NAME" "VARCHAR"
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| sync_inst_db
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
< row
: bytes
| test
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestSyncInstance/mariadb107 mariadb107#1
> query
| SELECT VERSION()
< result
< column "VERSION()" "VARCHAR"
< row
: bytes
| 10.7.8-MariaDB-1:10.7.8+maria~ubu2004
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
| 	SELECT
| 		user,
| 		host
| 	FROM mysql.user
| 	WHERE user <> ''
| 		AND user NOT LIKE 'mysql.%'
| 	ORDER BY user, host
| 	
< result
< column "User" "CHAR"
< column "Host" "CHAR"
< row
: bytes
| mariadb.sys
: bytes
| localhost
< row
: bytes
| root
: bytes
| %
< row
: bytes
| root
: bytes
| localhost
< row
: bytes
| sync_reader
: bytes
| %
> query
| SHOW GRANTS FOR 'mariadb.sys'@'localhost'
< result
< column "Grants for mariadb.sys@localhost" "VARCHAR"
< row
: bytes
| GRANT USAGE ON *.* TO `mariadb.sys`@`localhost`
< row
: bytes
| GRANT SELECT, DELETE ON `mysql`.`global_priv` TO `mariadb.sys`@`localhost`
> query
| SHOW GRANTS FOR 'root'@'%'
< result
< column "Grants for root@%" "VARCHAR"
< row
: bytes
| GRANT ALL PRIVILEGES ON *.* TO `root`@`%` IDENTIFIED BY PASSWORD '*94BDCEBE19083CE2A1F959FD02F964C7AF4CFC29' WITH GRANT OPTION
> query
| SHOW GRANTS FOR 'root'@'localhost'
< result
< column "Grants for root@localhost" "VARCHAR"
< row
: bytes
| GRANT ALL PRIVILEGES ON *.* TO `root`@`localhost` IDENTIFIED BY PASSWORD '*94BDCEBE19083CE2A1F959FD02F964C7AF4CFC29' WITH GRANT OPTION
< row
: bytes
| GRANT PROXY ON ''@'%' TO 'root'@'localhost' WITH GRANT OPTION
> query
| SHOW GRANTS FOR 'sync_reader'@'%'
< result
< column "Grants for sync_reader@%" "VARCHAR"
< row
: bytes
| GRANT USAGE ON *.* TO `sync_reader`@`%` IDENTIFIED BY PASSWORD '*0EB931DD0E0457476ABFBD89F6B34EA17B7C84C5'
< row
: bytes
| GRANT SELECT ON `sync_inst_db`.* TO `sync_reader`@`%`
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
< column "SCHEMA_NAME" "VARCHAR"
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| sync_inst_db
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
< row
: bytes
| test
: bytes
| utf8mb4
: bytes
| utf8mb4_general_ci

== TestSyncInstance/mariadb114 mariadb114#1
> query
| SELECT VERSION()
< result
< column "VERSION()" "VARCHAR"
< row
: bytes
| 11.4.13-MariaDB-ubu2404
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
| 	SELECT
| 		user,
| 		host
| 	FROM mysql.user
| 	WHERE user <> ''
| 		AND user NOT LIKE 'mysql.%'
| 	ORDER BY user, host
| 	
< result
< column "User" "CHAR"
< column "Host" "CHAR"
< row
: bytes
| mariadb.sys
: bytes
| localhost
< row
: bytes
| root
: bytes
| %
< row
: bytes
| root
: bytes
| localhost
< row
: bytes
| sync_reader
: bytes
| %
> query
| SHOW GRANTS FOR 'mariadb.sys'@'localhost'
< result
< column "Grants for mariadb.sys@localhost" "VARCHAR"
< row
: bytes
| GRANT USAGE ON *.* TO `mariadb.sys`@`localhost`
< row
: bytes
| GRANT SELECT, DELETE ON `mysql`.`global_priv` TO `mariadb.sys`@`localhost`
> query
| SHOW GRANTS FOR 'root'@'%'
< result
< column "Grants for root@%" "VARCHAR"
< row
: bytes
| GRANT ALL PRIVILEGES ON *.* TO `root`@`%` IDENTIFIED BY PASSWORD '*94BDCEBE19083CE2A1F959FD02F964C7AF4CFC29' WITH GRANT OPTION
< row
: bytes
| GRANT PROXY ON ``@`%` TO `root`@`%` WITH GRANT OPTION
> query
| SHOW GRANTS FOR 'root'@'localhost'
< result
< column "Grants for root@localhost" "VARCHAR"
< row
: bytes
| GRANT ALL PRIVILEGES ON *.* TO `root`@`localhost` IDENTIFIED BY PASSWORD '*94BDCEBE19083CE2A1F959FD02F964C7AF4CFC29' WITH GRANT OPTION
< row
: bytes
| GRANT PROXY ON ``@`%` TO `root`@`localhost` WITH GRANT OPTION
> query
| SHOW GRANTS FOR 'sync_reader'@'%'
< result
< column "Grants for sync_reader@%" "VARCHAR"
< row
: bytes
| GRANT USAGE ON *.* TO `sync_reader`@`%` IDENTIFIED BY PASSWORD '*0EB931DD0E0457476ABFBD89F6B34EA17B7C84C5'
< row
: bytes
| GRANT SELECT ON `sync_inst_db`.* TO `sync_reader`@`%`
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
< column "SCHEMA_NAME" "VARCHAR"
< column "DEFAULT_CHARACTER_SET_NAME" "VARCHAR"
< column "DEFAULT_COLLATION_NAME" "VARCHAR"
< row
: bytes
| sync_inst_db
: bytes
| utf8mb4
: bytes
| utf8mb4_unicode_ci
< row
: bytes
| test
: bytes
| utf8mb4
: bytes
| utf8mb4_uca1400_ai_ci
