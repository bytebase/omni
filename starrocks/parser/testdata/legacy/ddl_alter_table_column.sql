ALTER TABLE example_db.my_table ADD COLUMN new_col INT KEY DEFAULT "0" AFTER key_1;
ALTER TABLE example_db.my_table ADD COLUMN new_col INT DEFAULT "0" AFTER value_1;
ALTER TABLE example_db.my_table ADD COLUMN new_col INT SUM DEFAULT "0" AFTER value_1;
ALTER TABLE example_db.my_table ADD COLUMN new_col INT KEY DEFAULT "0" FIRST;
ALTER TABLE example_db.my_table ADD COLUMN (new_col1 INT SUM DEFAULT "0", new_col2 INT SUM DEFAULT "0");
ALTER TABLE example_db.my_table ADD COLUMN (new_col1 INT KEY DEFAULT "0", new_col2 INT DEFAULT "0");
ALTER TABLE example_db.my_table DROP COLUMN col1;
ALTER TABLE example_db.my_table MODIFY COLUMN col1 BIGINT KEY DEFAULT "1" AFTER col2;
ALTER TABLE example_db.my_table MODIFY COLUMN val1 VARCHAR(64) REPLACE DEFAULT "abc";
ALTER TABLE example_db.my_table MODIFY COLUMN k3 VARCHAR(50) KEY NULL COMMENT 'to 50';
ALTER TABLE example_db.my_table ORDER BY (k_2, k_1, v_3, v_2, v_1);
ALTER TABLE example_db.my_table ADD COLUMN col INT DEFAULT "0" AFTER v_1, ORDER BY (k_2, k_1, v_3, v_2, v_1, col)
