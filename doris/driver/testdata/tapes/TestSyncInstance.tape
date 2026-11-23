-- pin: default dyrnq/doris:3.1.4 sha256:88f99819c0877567b77d78cc8367f0e2dd4b6c1602fe17c6f4b1e332d8e315dd 3.1.4-rc02

== TestSyncInstance default#1
> query
| SELECT @@version_comment, VERSION()
< result
< column "@@version_comment" "TEXT"
< column "VERSION()" "TEXT"
< row
: bytes
| doris version doris-3.1.4-rc02-7f5ba43de6
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
| 		WHERE LOWER(SCHEMA_NAME) NOT IN ('__internal_schema', '_statistics_', 'information_schema', 'mysql')
| 		ORDER BY SCHEMA_NAME
< result
< column "SCHEMA_NAME" "CHAR"
< column "DEFAULT_CHARACTER_SET_NAME" "CHAR"
< column "DEFAULT_COLLATION_NAME" "CHAR"
< row
: bytes
| sync_inst_db
: bytes
| utf8mb4
: bytes
| utf8mb4_0900_bin
> query
| SHOW ALL GRANTS
< result
< column "UserIdentity" "CHAR"
< column "Comment" "CHAR"
< column "Password" "CHAR"
< column "Roles" "CHAR"
< column "GlobalPrivs" "CHAR"
< column "CatalogPrivs" "CHAR"
< column "DatabasePrivs" "CHAR"
< column "TablePrivs" "CHAR"
< column "ColPrivs" "CHAR"
< column "ResourcePrivs" "CHAR"
< column "CloudClusterPrivs" "CHAR"
< column "CloudStagePrivs" "CHAR"
< column "StorageVaultPrivs" "CHAR"
< column "WorkloadGroupPrivs" "CHAR"
< column "ComputeGroupPrivs" "CHAR"
< row
: bytes
| 'sync_app'@'%'
: bytes
|
: bytes
| Yes
: bytes
| sync_reader
: null
: null
: bytes
| internal.information_schema: Select_priv; internal.mysql: Select_priv; internal.sync_inst_db: Select_priv
: bytes
| internal.sync_inst_db.t: Load_priv
: null
: null
: null
: null
: null
: bytes
| normal: Usage_priv
: null
< row
: bytes
| 'sync_ops'@'10.0.0.%'
: bytes
| operations
: bytes
| No
: bytes
|
: null
: bytes
| internal: Create_priv
: bytes
| internal.information_schema: Select_priv; internal.mysql: Select_priv
: null
: bytes
| internal.sync_inst_db.t: Select_priv[v]
: bytes
| %: Usage_priv
: null
: null
: null
: bytes
| normal: Usage_priv
: null
< row
: bytes
| 'root'@'%'
: bytes
| ROOT
: bytes
| No
: bytes
| operator
: bytes
| Node_priv,Admin_priv
: null
: bytes
| internal.information_schema: Select_priv; internal.mysql: Select_priv
: null
: null
: null
: null
: null
: null
: bytes
| normal: Usage_priv
: null
< row
: bytes
| 'admin'@'%'
: bytes
| ADMIN
: bytes
| No
: bytes
| admin
: bytes
| Admin_priv
: null
: bytes
| internal.information_schema: Select_priv; internal.mysql: Select_priv
: null
: null
: null
: null
: null
: null
: bytes
| normal: Usage_priv
: null
> query
| SHOW ROLES
< result
< column "Name" "CHAR"
< column "Comment" "CHAR"
< column "Users" "CHAR"
< column "GlobalPrivs" "CHAR"
< column "CatalogPrivs" "CHAR"
< column "DatabasePrivs" "CHAR"
< column "TablePrivs" "CHAR"
< column "ResourcePrivs" "CHAR"
< column "CloudClusterPrivs" "CHAR"
< column "CloudStagePrivs" "CHAR"
< column "StorageVaultPrivs" "CHAR"
< column "WorkloadGroupPrivs" "CHAR"
< column "ComputeGroupPrivs" "CHAR"
< row
: bytes
| sync_reader
: bytes
| reads sync_inst_db
: bytes
| 'sync_app'@'%'
: null
: null
: bytes
| internal.sync_inst_db.*: Select_priv
: null
: null
: null
: null
: null
: null
: null
< row
: bytes
| admin
: null
: bytes
| 'admin'@'%'
: bytes
| Admin_priv
: null
: null
: null
: null
: null
: null
: null
: null
: null
< row
: bytes
| operator
: null
: bytes
| 'root'@'%'
: bytes
| Node_priv,Admin_priv
: null
: null
: null
: null
: null
: null
: null
: null
: null
