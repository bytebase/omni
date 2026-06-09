ALTER TABLE example_db.my_table
ADD ROLLUP example_rollup_index(k1, k3, v1, v2);
ALTER TABLE example_db.my_table
ADD ROLLUP example_rollup_index2 (k1, v1)
FROM example_rollup_index;
ALTER TABLE example_db.my_table
ADD ROLLUP example_rollup_index(k1, k3, v1)
PROPERTIES("timeout" = "3600");
ALTER TABLE example_db.my_table
DROP ROLLUP example_rollup_index2;
ALTER TABLE example_db.my_table
DROP ROLLUP example_rollup_index2,example_rollup_index3
