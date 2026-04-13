ALTER TABLE example_db.my_table ADD PARTITION p1 VALUES LESS THAN ("2014-01-01");
ALTER TABLE example_db.my_table ADD PARTITION p1 VALUES LESS THAN ("2015-01-01") DISTRIBUTED BY HASH(k1) BUCKETS 20;
ALTER TABLE example_db.my_table ADD PARTITION p1 VALUES LESS THAN ("2015-01-01") ("replication_num" = "1");
ALTER TABLE example_db.my_table MODIFY PARTITION p1 SET("replication_num" = "1");
ALTER TABLE example_db.my_table MODIFY PARTITION (p1, p2, p4) SET("replication_num" = "1");
ALTER TABLE example_db.my_table MODIFY PARTITION (*) SET("storage_medium" = "HDD");
ALTER TABLE example_db.my_table DROP PARTITION p1;
ALTER TABLE example_db.my_table DROP PARTITION p1, DROP PARTITION p2, DROP PARTITION p3;
ALTER TABLE example_db.my_table ADD PARTITION p1 VALUES [("2014-01-01"), ("2014-02-01"))
