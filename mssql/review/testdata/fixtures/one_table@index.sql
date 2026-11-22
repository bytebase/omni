CREATE TABLE dbo.a (b int NOT NULL CONSTRAINT pk_a PRIMARY KEY, a int, d int);
CREATE INDEX a ON dbo.a (a);
