CREATE TYPE dbo.a FROM int;
CREATE TABLE dbo.t (b int NOT NULL CONSTRAINT pk_t PRIMARY KEY, c int, d int);
CREATE INDEX t_d_idx ON dbo.t (d);
