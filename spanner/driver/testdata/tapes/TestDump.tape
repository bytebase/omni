-- pin: driver gcr.io/cloud-spanner-emulator/emulator:1.5.54 sha256:caf1bd24c081e005837b5977bae5a250e25cb4da9f25ec1abc91936ad67e4de2 v1.5.54

== TestDump driver#1
> raw
| {"op":"getDatabaseDdl","database":"projects/omni/instances/dump/databases/dump"}
< raw
| {"ddl":["CREATE TABLE Singers (\n  Id INT64 NOT NULL,\n  FirstName STRING(1024),\n  LastName STRING(1024) NOT NULL,\n  Active BOOL,\n  Rating FLOAT64,\n  Photo BYTES(MAX),\n  Born DATE,\n  UpdatedAt TIMESTAMP,\n  Meta JSON,\n  Scores ARRAY\u003cINT64\u003e,\n) PRIMARY KEY(Id)","CREATE UNIQUE INDEX SingersByFirstName ON Singers(FirstName)","CREATE INDEX SingersByLastName ON Singers(LastName DESC)","CREATE TABLE Albums (\n  Id INT64 NOT NULL,\n  SingerId INT64 NOT NULL,\n  Title STRING(1024) NOT NULL DEFAULT ('untitled'),\n  CONSTRAINT FK_Singer FOREIGN KEY(SingerId) REFERENCES Singers(Id),\n) PRIMARY KEY(Id)","CREATE VIEW SingerNames SQL SECURITY INVOKER AS SELECT s.Id AS Id, s.LastName AS LastName FROM Singers s"]}
> raw
| {"op":"close"}
< raw
| {}
