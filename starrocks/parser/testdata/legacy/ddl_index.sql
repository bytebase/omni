CREATE INDEX index1 ON table1 (col1) USING INVERTED;
CREATE INDEX index2 ON table1 (col1) USING NGRAM_BF PROPERTIES("gram_size" = "3", "bf_size" = "1024");
DROP INDEX IF EXISTS index_name ON table1;
BUILD INDEX index1 ON table1;
BUILD INDEX index1 ON table1 PARTITION(p1, p2)
