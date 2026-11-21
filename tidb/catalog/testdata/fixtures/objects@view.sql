CREATE TABLE t (b INT PRIMARY KEY, ident INT, d INT, KEY idx_d (d));
CREATE VIEW ident AS SELECT b, ident, d FROM t
