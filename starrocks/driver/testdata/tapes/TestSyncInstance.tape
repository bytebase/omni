-- pin: default starrocks/allin1-ubuntu:3.4.10 sha256:8bc803ea7a5bb490c77117a456b32b792ccae698b1f0e47dfcafdac41542566d 3.4.10

== TestSyncInstance default#1
> query
| SELECT @@version_comment, VERSION()
< result
< column "@@version_comment" "VARCHAR"
< column "version()" "VARCHAR"
< row
: bytes
| 3.4.10-965be75
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
> query
| SHOW USERS
< result
< column "User" "VARCHAR"
< row
: bytes
| 'root'@'%'
< row
: bytes
| 'sync_ops'@'10.0.0.%'
< row
: bytes
| 'sync_app'@'%'
> query
| SHOW ROLES
< result
< column "Name" "VARCHAR"
< column "Builtin" "VARCHAR"
< column "Comment" "VARCHAR"
< row
: bytes
| root
: bytes
| true
: bytes
| built-in root role which has all privileges on all objects
< row
: bytes
| db_admin
: bytes
| true
: bytes
| built-in database administration role
< row
: bytes
| cluster_admin
: bytes
| true
: bytes
| built-in cluster administration role
< row
: bytes
| user_admin
: bytes
| true
: bytes
| built-in user administration role
< row
: bytes
| public
: bytes
| true
: bytes
| built-in public role which is owned by any user
< row
: bytes
| sync_reader
: bytes
| false
: bytes
| reads sync_inst_db
> query
| SHOW GRANTS FOR 'root'@'%'
< result
< column "UserIdentity" "VARCHAR"
< column "Catalog" "VARCHAR"
< column "Grants" "VARCHAR"
< row
: bytes
| 'root'@'%'
: null
: bytes
| GRANT 'root' TO 'root'@'%'
> query
| SHOW GRANTS FOR 'sync_app'@'%'
< result
< column "UserIdentity" "VARCHAR"
< column "Catalog" "VARCHAR"
< column "Grants" "VARCHAR"
< row
: bytes
| 'sync_app'@'%'
: null
: bytes
| GRANT 'sync_reader' TO 'sync_app'@'%'
< row
: bytes
| 'sync_app'@'%'
: bytes
| default_catalog
: bytes
| GRANT INSERT ON TABLE sync_inst_db.t TO USER 'sync_app'@'%'
> query
| SHOW GRANTS FOR 'sync_ops'@'10.0.0.%'
< result
< column "UserIdentity" "VARCHAR"
< column "Catalog" "VARCHAR"
< column "Grants" "VARCHAR"
< row
: bytes
| 'sync_ops'@'10.0.0.%'
: null
: bytes
| GRANT 'sync_reader' TO 'sync_ops'@'10.0.0.%'
< row
: bytes
| 'sync_ops'@'10.0.0.%'
: bytes
| default_catalog
: bytes
| GRANT CREATE TABLE ON DATABASE sync_inst_db TO USER 'sync_ops'@'10.0.0.%'
> query
| SHOW GRANTS FOR ROLE `cluster_admin`
< result
< column "UserIdentity" "VARCHAR"
< column "Catalog" "VARCHAR"
< column "Grants" "VARCHAR"
< row
: bytes
| cluster_admin
: bytes
| default_catalog
: bytes
| GRANT USAGE, ALTER, DROP ON ALL WAREHOUSES TO ROLE 'cluster_admin'
< row
: bytes
| cluster_admin
: bytes
| default_catalog
: bytes
| GRANT NODE, CREATE WAREHOUSE ON SYSTEM TO ROLE 'cluster_admin'
> query
| SHOW GRANTS FOR ROLE `db_admin`
< result
< column "UserIdentity" "VARCHAR"
< column "Catalog" "VARCHAR"
< column "Grants" "VARCHAR"
< row
: bytes
| db_admin
: null
: bytes
| GRANT DELETE, DROP, INSERT, SELECT, ALTER, EXPORT, UPDATE ON ALL TABLES IN ALL DATABASES TO ROLE 'db_admin'
< row
: bytes
| db_admin
: null
: bytes
| GRANT CREATE TABLE, DROP, ALTER, CREATE VIEW, CREATE FUNCTION, CREATE MATERIALIZED VIEW, CREATE PIPE ON ALL DATABASES TO ROLE 'db_admin'
< row
: bytes
| db_admin
: bytes
| default_catalog
: bytes
| GRANT CREATE RESOURCE, PLUGIN, FILE, BLACKLIST, OPERATE, CREATE EXTERNAL CATALOG, REPOSITORY, CREATE RESOURCE GROUP, CREATE GLOBAL FUNCTION, CREATE STORAGE VOLUME ON SYSTEM TO ROLE 'db_admin'
< row
: bytes
| db_admin
: bytes
| default_catalog
: bytes
| GRANT USAGE, ALTER, DROP ON ALL RESOURCES TO ROLE 'db_admin'
< row
: bytes
| db_admin
: bytes
| default_catalog
: bytes
| GRANT SELECT, ALTER, DROP ON ALL VIEWS IN ALL DATABASES TO ROLE 'db_admin'
< row
: bytes
| db_admin
: null
: bytes
| GRANT USAGE, CREATE DATABASE, DROP, ALTER ON ALL CATALOGS TO ROLE 'db_admin'
< row
: bytes
| db_admin
: bytes
| default_catalog
: bytes
| GRANT ALTER, REFRESH, DROP, SELECT ON ALL MATERIALIZED VIEWS IN ALL DATABASES TO ROLE 'db_admin'
< row
: bytes
| db_admin
: bytes
| default_catalog
: bytes
| GRANT USAGE, DROP ON ALL FUNCTIONS IN ALL DATABASES TO ROLE 'db_admin'
< row
: bytes
| db_admin
: bytes
| default_catalog
: bytes
| GRANT ALTER, DROP ON ALL RESOURCE GROUPS TO ROLE 'db_admin'
< row
: bytes
| db_admin
: bytes
| default_catalog
: bytes
| GRANT USAGE, DROP ON ALL GLOBAL FUNCTIONS TO ROLE 'db_admin'
< row
: bytes
| db_admin
: bytes
| default_catalog
: bytes
| GRANT DROP, ALTER, USAGE ON ALL STORAGE VOLUMES TO ROLE 'db_admin'
< row
: bytes
| db_admin
: bytes
| default_catalog
: bytes
| GRANT ALTER, DROP, USAGE ON ALL DATABASES TO ROLE 'db_admin'
> query
| SHOW GRANTS FOR ROLE `public`
< result
< column "UserIdentity" "VARCHAR"
< column "Catalog" "VARCHAR"
< column "Grants" "VARCHAR"
< row
: bytes
| public
: bytes
| default_catalog
: bytes
| GRANT SELECT ON ALL TABLES IN DATABASE information_schema TO ROLE 'public'
> query
| SHOW GRANTS FOR ROLE `root`
< result
< column "UserIdentity" "VARCHAR"
< column "Catalog" "VARCHAR"
< column "Grants" "VARCHAR"
< row
: bytes
| root
: null
: bytes
| GRANT DELETE, DROP, INSERT, SELECT, ALTER, EXPORT, UPDATE ON ALL TABLES IN ALL DATABASES TO ROLE 'root'
< row
: bytes
| root
: null
: bytes
| GRANT CREATE TABLE, DROP, ALTER, CREATE VIEW, CREATE FUNCTION, CREATE MATERIALIZED VIEW, CREATE PIPE ON ALL DATABASES TO ROLE 'root'
< row
: bytes
| root
: bytes
| default_catalog
: bytes
| GRANT USAGE, ALTER, DROP ON ALL WAREHOUSES TO ROLE 'root'
< row
: bytes
| root
: bytes
| default_catalog
: bytes
| GRANT GRANT, NODE, CREATE RESOURCE, PLUGIN, FILE, BLACKLIST, OPERATE, CREATE EXTERNAL CATALOG, REPOSITORY, CREATE RESOURCE GROUP, CREATE GLOBAL FUNCTION, CREATE STORAGE VOLUME, CREATE WAREHOUSE ON SYSTEM TO ROLE 'root'
< row
: bytes
| root
: bytes
| default_catalog
: bytes
| GRANT IMPERSONATE ON ALL USERS TO ROLE 'root'
< row
: bytes
| root
: bytes
| default_catalog
: bytes
| GRANT USAGE, ALTER, DROP ON ALL RESOURCES TO ROLE 'root'
< row
: bytes
| root
: bytes
| default_catalog
: bytes
| GRANT SELECT, ALTER, DROP ON ALL VIEWS IN ALL DATABASES TO ROLE 'root'
< row
: bytes
| root
: null
: bytes
| GRANT USAGE, CREATE DATABASE, DROP, ALTER ON ALL CATALOGS TO ROLE 'root'
< row
: bytes
| root
: bytes
| default_catalog
: bytes
| GRANT ALTER, REFRESH, DROP, SELECT ON ALL MATERIALIZED VIEWS IN ALL DATABASES TO ROLE 'root'
< row
: bytes
| root
: bytes
| default_catalog
: bytes
| GRANT USAGE, DROP ON ALL FUNCTIONS IN ALL DATABASES TO ROLE 'root'
< row
: bytes
| root
: bytes
| default_catalog
: bytes
| GRANT ALTER, DROP ON ALL RESOURCE GROUPS TO ROLE 'root'
< row
: bytes
| root
: bytes
| default_catalog
: bytes
| GRANT USAGE, DROP ON ALL GLOBAL FUNCTIONS TO ROLE 'root'
< row
: bytes
| root
: bytes
| default_catalog
: bytes
| GRANT DROP, ALTER, USAGE ON ALL STORAGE VOLUMES TO ROLE 'root'
< row
: bytes
| root
: bytes
| default_catalog
: bytes
| GRANT ALTER, DROP, USAGE ON ALL DATABASES TO ROLE 'root'
> query
| SHOW GRANTS FOR ROLE `sync_reader`
< result
< column "UserIdentity" "VARCHAR"
< column "Catalog" "VARCHAR"
< column "Grants" "VARCHAR"
< row
: bytes
| sync_reader
: bytes
| default_catalog
: bytes
| GRANT SELECT ON ALL TABLES IN DATABASE sync_inst_db TO ROLE 'sync_reader'
> query
| SHOW GRANTS FOR ROLE `user_admin`
< result
< column "UserIdentity" "VARCHAR"
< column "Catalog" "VARCHAR"
< column "Grants" "VARCHAR"
< row
: bytes
| user_admin
: bytes
| default_catalog
: bytes
| GRANT GRANT ON SYSTEM TO ROLE 'user_admin'
< row
: bytes
| user_admin
: bytes
| default_catalog
: bytes
| GRANT IMPERSONATE ON ALL USERS TO ROLE 'user_admin'
