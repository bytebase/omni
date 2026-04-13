ALTER TABLE example_db.my_table SET ("bloom_filter_columns"="k1,k2,k3");
ALTER TABLE example_db.my_table SET ("colocate_with" = "t1");
ALTER TABLE example_db.my_table SET ("distribution_type" = "random");
ALTER TABLE example_db.my_table SET ("dynamic_partition.enable" = "false");
ALTER TABLE example_db.my_table SET ("dynamic_partition.enable" = "true", "dynamic_partition.time_unit" = "DAY", "dynamic_partition.end" = "3", "dynamic_partition.prefix" = "p", "dynamic_partition.buckets" = "32");
ALTER TABLE example_db.my_table SET ("in_memory" = "false");
ALTER TABLE example_db.my_table ENABLE FEATURE "BATCH_DELETE";
ALTER TABLE example_db.my_table ENABLE FEATURE "SEQUENCE_LOAD" WITH PROPERTIES ("function_column.sequence_type" = "Date");
ALTER TABLE example_db.my_table MODIFY DISTRIBUTION DISTRIBUTED BY HASH(k1) BUCKETS 50;
ALTER TABLE example_db.my_table MODIFY COMMENT "new comment";
ALTER TABLE example_db.my_table MODIFY COLUMN k1 COMMENT "k1", MODIFY COLUMN k2 COMMENT "k2";
ALTER TABLE example_db.mysql_table SET ("replication_num" = "2");
ALTER TABLE example_db.mysql_table SET ("default.replication_num" = "2");
ALTER TABLE example_db.mysql_table SET ("replication_allocation" = "tag.location.default: 1");
ALTER TABLE example_db.mysql_table SET ("default.replication_allocation" = "tag.location.default: 1");
ALTER TABLE example_db.mysql_table SET ("light_schema_change" = "true");
ALTER TABLE create_table_not_have_policy SET ("storage_policy" = "created_create_table_alter_policy");
ALTER TABLE create_table_partition MODIFY PARTITION (*) SET("storage_policy"="created_create_table_partition_alter_policy")
