CREATE TABLE t (b INT PRIMARY KEY, utf8mb4 INT, d INT, KEY idx_d (d));
CREATE VIEW utf8mb4 AS SELECT b, utf8mb4, d FROM t
