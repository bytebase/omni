CREATE DATABASE db_test;
CREATE DATABASE `db_test` PROPERTIES ("replication_allocation" = "tag.location.group_1:3");
CREATE DATABASE `db_test` PROPERTIES ("storage_vault_name" = "hdfs_demo_vault");
ALTER DATABASE example_db SET DATA QUOTA 10995116277760;
ALTER DATABASE example_db RENAME example_db2;
ALTER DATABASE example_db SET REPLICA QUOTA 102400;
ALTER DATABASE example_db SET PROPERTIES("replication_allocation" = "tag.location.default:2");
ALTER DATABASE example_db SET PROPERTIES("replication_allocation" = "");
ALTER DATABASE example_db SET PROPERTIES("storage_vault_name" = "hdfs_demo_vault");
ALTER DATABASE example_db SET PROPERTIES("storage_vault_name" = "");
DROP DATABASE db_test;
SHOW DATABASES;
SHOW DATABASES FROM hms_catalog;
SHOW DATABASES LIKE 'infor%'
