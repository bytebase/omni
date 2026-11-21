CREATE TABLE t (b int PRIMARY KEY, a int, d int);
CREATE VIEW a AS SELECT b, a, d FROM t;
