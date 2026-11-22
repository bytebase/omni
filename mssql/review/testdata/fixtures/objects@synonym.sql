CREATE TABLE dbo.t (b int NOT NULL CONSTRAINT pk_t PRIMARY KEY, a int, d int);
CREATE SYNONYM dbo.a FOR dbo.t;
