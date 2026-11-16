CREATE TABLE Singers (
  Id INT64 NOT NULL,
  FirstName STRING(1024),
  LastName STRING(1024) NOT NULL,
  Active BOOL,
  Rating FLOAT64,
  Photo BYTES(MAX),
  Born DATE,
  UpdatedAt TIMESTAMP,
  Meta JSON,
  Scores ARRAY<INT64>,
) PRIMARY KEY(Id);
CREATE UNIQUE INDEX SingersByFirstName ON Singers(FirstName);
CREATE INDEX SingersByLastName ON Singers(LastName DESC);
CREATE TABLE Albums (
  Id INT64 NOT NULL,
  SingerId INT64 NOT NULL,
  Title STRING(1024) NOT NULL DEFAULT ('untitled'),
  CONSTRAINT FK_Singer FOREIGN KEY(SingerId) REFERENCES Singers(Id),
) PRIMARY KEY(Id);
CREATE VIEW SingerNames SQL SECURITY INVOKER AS SELECT s.Id AS Id, s.LastName AS LastName FROM Singers s;
