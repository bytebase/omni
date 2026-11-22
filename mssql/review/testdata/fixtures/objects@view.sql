CREATE TABLE dbo.t (b int NOT NULL CONSTRAINT pk_t PRIMARY KEY, a int, d int);
GO
CREATE VIEW dbo.a AS SELECT b, a, d FROM dbo.t;
