-- pin: default gcr.io/cloud-spanner-emulator/emulator:1.5.54 sha256:caf1bd24c081e005837b5977bae5a250e25cb4da9f25ec1abc91936ad67e4de2 v1.5.54

== TestCorpusSpannerDifferential default#1
> query
| SELECT CAST('123' AS INT64)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT CAST(3.14 AS FLOAT64)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT SAFE_CAST('not_a_number' AS INT64)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| -- returns NULL
| SELECT CAST(TRUE AS INT64)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT NULL
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT CAST(NULL AS INT64)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| CREATE DATABASE my_database
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"Internal","message":"GOOGLESQL_RET_CHECK failure (backend/schema/updater/schema_updater.cc:1182) false Unsupported ddl statement: 1"}
> query
| CREATE TABLE Singers (
|   SingerId   INT64 NOT NULL,
|   FirstName  STRING(1024),
|   LastName   STRING(1024),
|   BirthDate  DATE,
|   SingerInfo BYTES(MAX),
| ) PRIMARY KEY (SingerId)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| CREATE TABLE Albums (
|   SingerId     INT64 NOT NULL,
|   AlbumId      INT64 NOT NULL,
|   Title        STRING(MAX),
|   LastUpdated  TIMESTAMP OPTIONS (allow_commit_timestamp = true),
| ) PRIMARY KEY (SingerId, AlbumId),
|   INTERLEAVE IN PARENT Singers ON DELETE CASCADE
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Table not found: Singers"}
> query
| CREATE TABLE PlayerStats (
|   PlayerId   INT64 NOT NULL,
|   CreatedAt  TIMESTAMP NOT NULL,
|   Score      INT64,
| ) PRIMARY KEY (PlayerId, CreatedAt),
|   ROW DELETION POLICY (OLDER_THAN(CreatedAt, INTERVAL 30 DAY))
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| CREATE TABLE Orders (
|   OrderId     INT64 NOT NULL,
|   Status      STRING(20) DEFAULT ('pending'),
|   CreatedAt   TIMESTAMP DEFAULT (CURRENT_TIMESTAMP()),
| ) PRIMARY KEY (OrderId)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| CREATE TABLE Events (
|   EventId    INT64 NOT NULL,
|   LastMod    TIMESTAMP OPTIONS (allow_commit_timestamp = true),
| ) PRIMARY KEY (EventId)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| CREATE TABLE Chapters (
|   SingerId  INT64 NOT NULL,
|   AlbumId   INT64 NOT NULL,
|   ChapterId INT64 NOT NULL,
|   Title     STRING(MAX),
| ) PRIMARY KEY (SingerId, AlbumId, ChapterId),
|   INTERLEAVE IN PARENT Albums ON DELETE CASCADE
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Table not found: Albums"}
> query
| -- NO ACTION (default)
| CREATE TABLE Details (
|   ParentId INT64 NOT NULL,
|   DetailId INT64 NOT NULL,
| ) PRIMARY KEY (ParentId, DetailId),
|   INTERLEAVE IN PARENT Parent ON DELETE NO ACTION
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Table not found: Parent"}
> query
| CREATE TABLE Orders (
|   OrderId    INT64 NOT NULL,
|   CustomerId INT64 NOT NULL,
|   CONSTRAINT FK_Customers
|     FOREIGN KEY (CustomerId) REFERENCES Customers (CustomerId) ON DELETE NO ACTION,
| ) PRIMARY KEY (OrderId)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Table not found: Customers"}
> query
| CREATE TABLE Products (
|   ProductId INT64 NOT NULL,
|   Price     FLOAT64,
|   CONSTRAINT chk_positive_price CHECK (Price > 0),
| ) PRIMARY KEY (ProductId)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| CREATE TABLE Sessions (
|   SessionId   INT64 NOT NULL,
|   CreatedAt   TIMESTAMP NOT NULL,
| ) PRIMARY KEY (SessionId),
|   ROW DELETION POLICY (OLDER_THAN(CreatedAt, INTERVAL 90 DAY))
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| ALTER TABLE Singers ADD COLUMN Nickname STRING(100)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Table not found: Singers"}
> query
| ALTER TABLE Singers ADD COLUMN IF NOT EXISTS Bio STRING(MAX)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Table not found: Singers"}
> query
| ALTER TABLE Singers DROP COLUMN Nickname
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Table not found: Singers"}
> query
| ALTER TABLE Orders ADD CONSTRAINT FK_Cust FOREIGN KEY (CustomerId) REFERENCES Customers (CustomerId)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Table not found: Orders"}
> query
| ALTER TABLE Orders DROP CONSTRAINT FK_Cust
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Table not found: Orders"}
> query
| ALTER TABLE Singers RENAME TO Artists
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Table not found: Singers"}
> query
| DROP TABLE Singers
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Table not found: Singers"}
> query
| DROP TABLE IF EXISTS temp_table
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| CREATE INDEX SingersByName ON Singers(LastName, FirstName)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Table not found: Singers"}
> query
| CREATE UNIQUE INDEX UniqueEmail ON Users(Email)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Table not found: Users"}
> query
| CREATE NULL_FILTERED INDEX AlbumsByTitle ON Albums(Title)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Table not found: Albums"}
> query
| CREATE INDEX AlbumsByArtist ON Albums(SingerId, Title) STORING (MarketingBudget)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Table not found: Albums"}
> query
| CREATE INDEX SongsBySinger ON Songs(SingerId, AlbumId, TrackId)
|   STORING (SongName), INTERLEAVE IN Albums
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Table not found: Songs"}
> query
| ALTER INDEX SingersByName ADD STORED COLUMN BirthDate
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Index not found: SingersByName"}
> query
| ALTER INDEX SingersByName DROP STORED COLUMN BirthDate
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Index not found: SingersByName"}
> query
| DROP INDEX SingersByName
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Index not found: SingersByName"}
> query
| DROP INDEX IF EXISTS SingersByName
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| DROP VIEW SingerView
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"View not found: SingerView"}
> query
| DROP VIEW IF EXISTS TopAlbums
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| CREATE CHANGE STREAM MyStream FOR ALL
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| CREATE CHANGE STREAM SingersStream
|   FOR Singers,
|       Albums (Title, MarketingBudget)
|   OPTIONS (retention_period = '7d', value_capture_type = 'NEW_ROW')
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"FailedPrecondition","message":"Change Stream SingersStream cannot track Singers because either Change Streams do not track database objects of this type, or this Table does not exist."}
> query
| ALTER CHANGE STREAM MyStream SET FOR Singers, Albums
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Change Stream not found: MyStream"}
> query
| ALTER CHANGE STREAM MyStream SET OPTIONS (retention_period = '14d')
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Change Stream not found: MyStream"}
> query
| ALTER CHANGE STREAM MyStream DROP FOR ALL
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Change Stream not found: MyStream"}
> query
| DROP CHANGE STREAM MyStream
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Change Stream not found: MyStream"}
> query
| CREATE SEQUENCE MySeq OPTIONS (
|   sequence_kind = 'bit_reversed_positive'
| )
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| CREATE SEQUENCE OrderSeq OPTIONS (
|   sequence_kind = 'bit_reversed_positive',
|   skip_range_min = 1,
|   skip_range_max = 1000,
|   start_with_counter = 5000
| )
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| ALTER SEQUENCE MySeq SET OPTIONS (start_with_counter = 9000)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Sequence not found: MySeq"}
> query
| DROP SEQUENCE MySeq
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Sequence not found: MySeq"}
> query
| DROP SEQUENCE IF EXISTS MySeq
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| CREATE SCHEMA myschema
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| CREATE SCHEMA IF NOT EXISTS analytics
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| DROP SCHEMA myschema
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Schema not found: myschema"}
> query
| DROP SCHEMA IF EXISTS analytics
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| CREATE ROLE analyst
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| CREATE ROLE read_only_user
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| DROP ROLE analyst
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| GRANT ROLE analyst TO ROLE senior_analyst
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| REVOKE SELECT ON TABLE Singers FROM ROLE analyst
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| REVOKE INSERT, UPDATE ON TABLE Albums FROM ROLE editor
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| REVOKE ROLE analyst FROM ROLE senior_analyst
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| CREATE LOCALITY GROUP hot_data OPTIONS (storage = 'ssd')
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| CREATE LOCALITY GROUP cold_data OPTIONS (storage = 'hdd', ssd_to_hdd_spill_timespan = '30d')
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| ALTER LOCALITY GROUP cold_data SET OPTIONS (ssd_to_hdd_spill_timespan = '14d')
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Locality group not found: cold_data"}
> query
| DROP LOCALITY GROUP cold_data
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Locality group not found: cold_data"}
> query
| CREATE PLACEMENT us_placement OPTIONS (instance_partition = 'us-partition')
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| DROP PLACEMENT us_placement
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Placement not found: us_placement"}
> query
| CREATE PROTO BUNDLE (
|   `my.package.MyMessage`,
|   `my.package.AnotherMessage`
| )
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"InvalidArgument","message":"Contents of proto_descriptor_file cannot be empty"}
> query
| ALTER PROTO BUNDLE INSERT (`my.package.NewMessage`)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"FailedPrecondition","message":"Proto bundle does not yet exist; cannot alter it"}
> query
| ALTER PROTO BUNDLE UPDATE (`my.package.MyMessage`)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"FailedPrecondition","message":"Proto bundle does not yet exist; cannot alter it"}
> query
| ALTER PROTO BUNDLE DELETE (`my.package.OldMessage`)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"FailedPrecondition","message":"Proto bundle does not yet exist; cannot alter it"}
> query
| -- Requires TOKENLIST column created via TOKENIZE_* functions
| CREATE SEARCH INDEX SingerNameIndex ON Singers (SingerNameTokens)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Table not found: Singers"}
> query
| CREATE SEARCH INDEX AlbumTitleIdx ON Albums (TitleTokens)
|   STORING (MarketingBudget)
|   PARTITION BY SingerId
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Table not found: Albums"}
> query
| DROP SEARCH INDEX SingerNameIndex
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Index not found: SingerNameIndex"}
> query
| CREATE TABLE Products (
|   ProductId    INT64 NOT NULL,
|   Name         STRING(MAX),
|   Description  STRING(MAX),
|   NameTokens   TOKENLIST AS (TOKENIZE_FULLTEXT(Name)) STORED HIDDEN,
| ) PRIMARY KEY (ProductId)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| CREATE TABLE Products (
|   ProductId  INT64 NOT NULL,
|   Embedding  ARRAY<FLOAT32>(vector_length=>128),
| ) PRIMARY KEY (ProductId)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| CREATE VECTOR INDEX ProductEmbeddingIdx
|   ON Products (Embedding)
|   OPTIONS (distance_type = 'COSINE', num_leaves = 100)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Table not found: Products"}
> query
| DROP VECTOR INDEX ProductEmbeddingIdx
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Index not found: ProductEmbeddingIdx"}
> query
| CREATE TABLE myschema.Employees (
|   EmpId INT64 NOT NULL,
|   Name  STRING(100),
| ) PRIMARY KEY (EmpId)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Schema not found: myschema"}
> query
| ALTER TABLE myschema.Employees ADD COLUMN Department STRING(50)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Table not found: myschema.Employees"}
> query
| CREATE TABLE Metrics (
|   MetricId   INT64 NOT NULL,
|   HotValue   FLOAT64 OPTIONS (locality_group = 'hot_data'),
|   ColdValue  FLOAT64 OPTIONS (locality_group = 'cold_data'),
| ) PRIMARY KEY (MetricId)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"semantic","code":"NotFound","message":"Locality group not found: hot_data"}
> query
| INSERT INTO Singers (SingerId, FirstName, LastName)
| VALUES (1, 'Marc', 'Richards')
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"dml","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:13] INSERT INTO Singers (SingerId, FirstName, LastName)             ^"}
> query
| INSERT INTO Singers (SingerId, FirstName, LastName)
| VALUES (1, 'Marc', 'Richards'),
|        (2, 'Catalina', 'Smith')
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"dml","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:13] INSERT INTO Singers (SingerId, FirstName, LastName)             ^"}
> query
| INSERT OR IGNORE INTO Singers (SingerId, FirstName, LastName)
| VALUES (1, 'Marc', 'Richards')
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"dml","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:23] INSERT OR IGNORE INTO Singers (SingerId, FirstName, LastName)                       ^"}
> query
| INSERT OR UPDATE INTO Singers (SingerId, FirstName, LastName)
| VALUES (1, 'Marc', 'Richards')
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"dml","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:23] INSERT OR UPDATE INTO Singers (SingerId, FirstName, LastName)                       ^"}
> query
| INSERT INTO ArchivedSingers (SingerId, FirstName, LastName)
| SELECT SingerId, FirstName, LastName FROM Singers WHERE Retired = TRUE
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"dml","reason":"semantic","code":"InvalidArgument","message":"Table not found: ArchivedSingers [at 1:13] INSERT INTO ArchivedSingers (SingerId, FirstName, LastName)             ^"}
> query
| INSERT INTO Singers (SingerId, FirstName, LastName)
| VALUES (3, 'Alice', 'Liu')
| THEN RETURN SingerId, FirstName
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"dml","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:13] INSERT INTO Singers (SingerId, FirstName, LastName)             ^"}
> query
| INSERT OR UPDATE INTO Singers (SingerId, FirstName, LastName)
| VALUES (1, 'Marc', 'Richards')
| THEN RETURN WITH ACTION AS op_type *
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"dml","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:23] INSERT OR UPDATE INTO Singers (SingerId, FirstName, LastName)                       ^"}
> query
| DELETE FROM Singers WHERE SingerId = 5
| THEN RETURN SingerId
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"dml","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:13] DELETE FROM Singers WHERE SingerId = 5             ^"}
> query
| UPDATE Singers
| SET FirstName = 'Alice'
| WHERE SingerId = 1
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"dml","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:8] UPDATE Singers        ^"}
> query
| UPDATE Albums
| SET MarketingBudget = MarketingBudget * 1.1
| WHERE SingerId = 1
|   AND AlbumId IN (SELECT AlbumId FROM TopAlbums)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"dml","reason":"semantic","code":"InvalidArgument","message":"Table not found: Albums [at 1:8] UPDATE Albums        ^"}
> query
| UPDATE Singers s
| SET s.LastName = u.NewLastName
| FROM Updates u
| WHERE s.SingerId = u.SingerId
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"dml","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:8] UPDATE Singers s        ^"}
> query
| DELETE FROM Singers WHERE SingerId = 1
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"dml","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:13] DELETE FROM Singers WHERE SingerId = 1             ^"}
> query
| DELETE FROM Albums
| WHERE SingerId NOT IN (SELECT SingerId FROM Singers)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"dml","reason":"semantic","code":"InvalidArgument","message":"Table not found: Albums [at 1:13] DELETE FROM Albums             ^"}
> query
| DELETE Singers AS s WHERE s.SingerId = 42
| THEN RETURN s.FirstName, s.LastName
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"dml","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:8] DELETE Singers AS s WHERE s.SingerId = 42        ^"}
> query
| @{PDML_MAX_PARALLELISM=4}
| DELETE FROM OldEvents WHERE CreatedAt < '2020-01-01'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"dml","reason":"semantic","code":"InvalidArgument","message":"Table not found: OldEvents [at 2:13] DELETE FROM OldEvents WHERE CreatedAt \u003c '2020-01-01'             ^"}
> query
| @{LOCK_SCANNED_RANGES=shared}
| UPDATE Accounts SET Balance = 0 WHERE Status = 'closed'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"dml","reason":"semantic","code":"InvalidArgument","message":"Table not found: Accounts [at 2:8] UPDATE Accounts SET Balance = 0 WHERE Status = 'closed'        ^"}
> query
| SELECT CASE Status
|   WHEN 'active' THEN 1
|   WHEN 'inactive' THEN 0
|   ELSE -1
| END AS StatusCode
| FROM Accounts
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Accounts [at 6:6] FROM Accounts      ^"}
> query
| SELECT CASE
|   WHEN Score > 90 THEN 'A'
|   WHEN Score > 80 THEN 'B'
|   WHEN Score > 70 THEN 'C'
|   ELSE 'F'
| END AS Grade
| FROM Students
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Students [at 7:6] FROM Students      ^"}
> query
| SELECT CAST('123' AS INT64)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT CAST(3.14 AS STRING)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT CAST('2024-01-01' AS DATE)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT CAST(TRUE AS INT64)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| -- 1
| SELECT CAST(CURRENT_TIMESTAMP() AS DATE)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT SAFE_CAST('not_a_number' AS INT64)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| -- NULL
| SELECT SAFE_CAST('2024-01-01' AS DATE)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT SAFE_CAST(NULL AS STRING)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT EXTRACT(YEAR FROM BirthDate) AS birth_year FROM Singers
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:56] SELECT EXTRACT(YEAR FROM BirthDate) AS birth_year FROM Singers                                                        ^"}
> query
| SELECT EXTRACT(HOUR FROM CURRENT_TIMESTAMP() AT TIME ZONE 'America/New_York')
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT EXTRACT(DATE FROM TIMESTAMP '2024-06-01 12:30:00 UTC')
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT EXTRACT(DAYOFWEEK FROM CURRENT_DATE())
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT CONCAT(FirstName, ' ', LastName) FROM Singers
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:46] SELECT CONCAT(FirstName, ' ', LastName) FROM Singers                                              ^"}
> query
| SELECT UPPER(LastName), LOWER(FirstName) FROM Singers
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:47] SELECT UPPER(LastName), LOWER(FirstName) FROM Singers                                               ^"}
> query
| SELECT CURRENT_TIMESTAMP(), CURRENT_DATE()
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT COUNT(DISTINCT SingerId) FROM Albums
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Albums [at 1:38] SELECT COUNT(DISTINCT SingerId) FROM Albums                                      ^"}
> query
| SELECT SUM(DISTINCT Score) FROM Scores
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Scores [at 1:33] SELECT SUM(DISTINCT Score) FROM Scores                                 ^"}
> query
| SELECT ARRAY_AGG(SingerId IGNORE NULLS) FROM Singers
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:46] SELECT ARRAY_AGG(SingerId IGNORE NULLS) FROM Singers                                              ^"}
> query
| SELECT STRING_AGG(Name RESPECT NULLS ORDER BY Name) FROM Tags
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Tags [at 1:58] SELECT STRING_AGG(Name RESPECT NULLS ORDER BY Name) FROM Tags                                                          ^"}
> query
| SELECT LAST_VALUE(Score IGNORE NULLS) OVER (PARTITION BY Player ORDER BY Ts) FROM Scores
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Scores [at 1:83] ...NULLS) OVER (PARTITION BY Player ORDER BY Ts) FROM Scores                                                       ^"}
> query
| SELECT ARRAY_AGG(SingerId ORDER BY LastName) FROM Singers
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:51] SELECT ARRAY_AGG(SingerId ORDER BY LastName) FROM Singers                                                   ^"}
> query
| SELECT STRING_AGG(Title, ', ' ORDER BY AlbumId DESC) FROM Albums
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Albums [at 1:59] SELECT STRING_AGG(Title, ', ' ORDER BY AlbumId DESC) FROM Albums                                                           ^"}
> query
| SELECT ARRAY_AGG(SingerId ORDER BY LastName LIMIT 5) FROM Singers
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:59] SELECT ARRAY_AGG(SingerId ORDER BY LastName LIMIT 5) FROM Singers                                                           ^"}
> query
| SELECT SingerId, AlbumId,
|   ROW_NUMBER() OVER (PARTITION BY SingerId ORDER BY AlbumId) AS rn,
|   SUM(MarketingBudget) OVER (PARTITION BY SingerId
|                              ORDER BY AlbumId
|                              ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS running_budget
| FROM Albums
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Albums [at 6:6] FROM Albums      ^"}
> query
| -- Named window
| SELECT AVG(Score) OVER w AS avg_score
| FROM Scores
| WINDOW w AS (PARTITION BY PlayerId ORDER BY Ts ROWS BETWEEN 5 PRECEDING AND CURRENT ROW)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Scores [at 3:6] FROM Scores      ^"}
> query
| SELECT [1, 2, 3][OFFSET(0)]
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| -- 1
| SELECT [1, 2, 3][OFFSET(2)]
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| -- 3
| SELECT arr[OFFSET(0)] FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 2:28] SELECT arr[OFFSET(0)] FROM t                            ^"}
> query
| SELECT [1, 2, 3][ORDINAL(1)]
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| -- 1
| SELECT [1, 2, 3][ORDINAL(3)]
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT [1, 2, 3][SAFE_OFFSET(10)]
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT [1, 2, 3][SAFE_ORDINAL(10)]
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT [1, 2, 3]
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT ARRAY<INT64>[1, 2, 3]
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT ARRAY['a', 'b', 'c']
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT ARRAY(SELECT SingerId FROM Singers ORDER BY SingerId)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:35] SELECT ARRAY(SELECT SingerId FROM Singers ORDER BY SingerId)                                   ^"}
> query
| SELECT s.name FROM StructTable AS s
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: StructTable [at 1:20] SELECT s.name FROM StructTable AS s                    ^"}
> query
| SELECT STRUCT(1 AS id, 'Alice' AS name).name
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| -- 'Alice'
| SELECT (SELECT AS STRUCT SingerId, FirstName FROM Singers LIMIT 1).FirstName
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 2:51] SELECT (SELECT AS STRUCT SingerId, FirstName FROM Singers LIMIT 1).FirstName                                                   ^"}
> query
| SELECT TIMESTAMP_ADD(CURRENT_TIMESTAMP(), INTERVAL 1 HOUR)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT DATE_ADD(CURRENT_DATE(), INTERVAL 30 DAY)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 7 DAY)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT DATE_DIFF(DATE '2024-12-31', DATE '2024-01-01', DAY)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT IF(Score > 50, 'pass', 'fail') FROM Students
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Students [at 1:44] SELECT IF(Score \u003e 50, 'pass', 'fail') FROM Students                                            ^"}
> query
| SELECT IFNULL(NickName, FirstName) FROM Singers
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:41] SELECT IFNULL(NickName, FirstName) FROM Singers                                         ^"}
> query
| SELECT NULLIF(Status, 'unknown') FROM Orders
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Orders [at 1:39] SELECT NULLIF(Status, 'unknown') FROM Orders                                       ^"}
> query
| SELECT COALESCE(NickName, FirstName, 'Anonymous') FROM Singers
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:56] SELECT COALESCE(NickName, FirstName, 'Anonymous') FROM Singers                                                        ^"}
> query
| SELECT * FROM Singers WHERE SingerId IN (1, 2, 3)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:15] SELECT * FROM Singers WHERE SingerId IN (1, 2, 3)               ^"}
> query
| SELECT * FROM Singers WHERE SingerId NOT IN (SELECT SingerId FROM Banned)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:15] SELECT * FROM Singers WHERE SingerId NOT IN (SELECT SingerId FROM Banned)               ^"}
> query
| SELECT * FROM Singers WHERE SingerId IN UNNEST(@singer_ids)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:15] SELECT * FROM Singers WHERE SingerId IN UNNEST(@singer_ids)               ^"}
> query
| SELECT * FROM Singers
| WHERE EXISTS (SELECT 1 FROM Albums WHERE Albums.SingerId = Singers.SingerId)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:15] SELECT * FROM Singers               ^"}
> query
| SELECT * FROM Singers WHERE LastName LIKE 'Sm%'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:15] SELECT * FROM Singers WHERE LastName LIKE 'Sm%'               ^"}
> query
| SELECT * FROM Singers WHERE FirstName NOT LIKE '_a%'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:15] SELECT * FROM Singers WHERE FirstName NOT LIKE '_a%'               ^"}
> query
| SELECT * FROM Singers WHERE BirthDate IS NULL
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:15] SELECT * FROM Singers WHERE BirthDate IS NULL               ^"}
> query
| SELECT * FROM Orders WHERE Processed IS NOT NULL
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Orders [at 1:15] SELECT * FROM Orders WHERE Processed IS NOT NULL               ^"}
> query
| SELECT * FROM Flags WHERE Active IS TRUE
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Flags [at 1:15] SELECT * FROM Flags WHERE Active IS TRUE               ^"}
> query
| SELECT * FROM Singers WHERE SingerId BETWEEN 1 AND 100
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:15] SELECT * FROM Singers WHERE SingerId BETWEEN 1 AND 100               ^"}
> query
| SELECT * FROM Albums WHERE MarketingBudget NOT BETWEEN 0 AND 10000
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Albums [at 1:15] SELECT * FROM Albums WHERE MarketingBudget NOT BETWEEN 0 AND 10000               ^"}
> query
| SELECT 1 + 2, 10 - 3, 4 * 5, 8 / 2
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT -Score, +Amount FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:29] SELECT -Score, +Amount FROM t                             ^"}
> query
| SELECT * FROM t WHERE a = b
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:15] SELECT * FROM t WHERE a = b               ^"}
> query
| SELECT * FROM t WHERE x != y
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:15] SELECT * FROM t WHERE x != y               ^"}
> query
| SELECT * FROM t WHERE Score >= 90
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:15] SELECT * FROM t WHERE Score \u003e= 90               ^"}
> query
| SELECT * FROM t WHERE a AND b
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:15] SELECT * FROM t WHERE a AND b               ^"}
> query
| SELECT * FROM t WHERE x OR y
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:15] SELECT * FROM t WHERE x OR y               ^"}
> query
| SELECT * FROM t WHERE NOT active
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:15] SELECT * FROM t WHERE NOT active               ^"}
> query
| SELECT 5 & 3, 5 | 3, 5 ^ 3, ~5
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT flags << 2, flags >> 1 FROM t
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:36] SELECT flags \u003c\u003c 2, flags \u003e\u003e 1 FROM t                                    ^"}
> query
| SELECT 'Hello' || ', ' || 'World'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT FirstName || ' ' || LastName AS FullName FROM Singers
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:54] SELECT FirstName || ' ' || LastName AS FullName FROM Singers                                                      ^"}
> query
| SELECT [1,2] || [3,4]
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT (SELECT MAX(MarketingBudget) FROM Albums WHERE SingerId = s.SingerId)
| FROM Singers s
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 2:6] FROM Singers s      ^"}
> query
| SELECT ARRAY(SELECT AlbumId FROM Albums WHERE SingerId = s.SingerId) AS albums
| FROM Singers s
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 2:6] FROM Singers s      ^"}
> query
| SELECT * FROM t WHERE x > ALL (SELECT score FROM scores)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:15] SELECT * FROM t WHERE x \u003e ALL (SELECT score FROM scores)               ^"}
> query
| SELECT * FROM t WHERE x = ANY (SELECT val FROM vals)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: t [at 1:15] SELECT * FROM t WHERE x = ANY (SELECT val FROM vals)               ^"}
> query
| INSERT INTO Events (EventId, LastMod) VALUES (1, PENDING_COMMIT_TIMESTAMP())
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"dml","reason":"semantic","code":"InvalidArgument","message":"Table not found: Events [at 1:13] INSERT INTO Events (EventId, LastMod) VALUES (1, PENDING_COMMIT_TIMESTAMP())             ^"}
> query
| UPDATE Events SET LastMod = PENDING_COMMIT_TIMESTAMP() WHERE EventId = 1
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"dml","reason":"semantic","code":"InvalidArgument","message":"Table not found: Events [at 1:8] UPDATE Events SET LastMod = PENDING_COMMIT_TIMESTAMP() WHERE EventId = 1        ^"}
> query
| SELECT * FROM Singers WHERE SingerId IN UNNEST(@singer_ids)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:15] SELECT * FROM Singers WHERE SingerId IN UNNEST(@singer_ids)               ^"}
> query
| SELECT * FROM Albums WHERE Title NOT IN UNNEST(@excluded_titles)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Albums [at 1:15] SELECT * FROM Albums WHERE Title NOT IN UNNEST(@excluded_titles)               ^"}
> query
| @{USE_ADDITIONAL_PARALLELISM=TRUE}
| SELECT * FROM Singers
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 2:15] SELECT * FROM Singers               ^"}
> query
| @{OPTIMIZER_VERSION=2, OPTIMIZER_STATISTICS_PACKAGE='auto_20240101'}
| SELECT * FROM Albums
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Albums [at 2:15] SELECT * FROM Albums               ^"}
> query
| SELECT * FROM Singers @{FORCE_INDEX=SingersByName}
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:15] SELECT * FROM Singers @{FORCE_INDEX=SingersByName}               ^"}
> query
| SELECT * FROM Albums @{FORCE_INDEX=AlbumsByTitle, GROUPBY_SCAN_OPTIMIZATION=TRUE} AS a
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Albums [at 1:15] SELECT * FROM Albums @{FORCE_INDEX=AlbumsByTitle, GROUPBY_SCAN_OPTIMIZATION=T...               ^"}
> query
| SELECT *
| FROM Singers s
| JOIN @{JOIN_METHOD=HASH_JOIN} Albums a ON s.SingerId = a.SingerId
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 2:6] FROM Singers s      ^"}
> query
| SELECT *
| FROM Singers s
| JOIN @{JOIN_METHOD=APPLY_JOIN} Albums a ON s.SingerId = a.SingerId
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 2:6] FROM Singers s      ^"}
> query
| SELECT * FROM Singers @{FORCE_INDEX=SingersByLastName}
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:15] SELECT * FROM Singers @{FORCE_INDEX=SingersByLastName}               ^"}
> query
| SELECT * FROM Singers @{FORCE_INDEX=_BASE_TABLE}
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:15] SELECT * FROM Singers @{FORCE_INDEX=_BASE_TABLE}               ^"}
> query
| SELECT SingerId, COUNT(*) FROM Albums @{GROUPBY_SCAN_OPTIMIZATION=TRUE} GROUP BY SingerId
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Albums [at 1:32] SELECT SingerId, COUNT(*) FROM Albums @{GROUPBY_SCAN_OPTIMIZATION=TRUE} GROUP...                                ^"}
> query
| @{OPTIMIZER_VERSION=4}
| SELECT * FROM Singers
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 2:15] SELECT * FROM Singers               ^"}
> query
| @{OPTIMIZER_VERSION=latest}
| SELECT * FROM Singers
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 2:15] SELECT * FROM Singers               ^"}
> query
| @{OPTIMIZER_STATISTICS_PACKAGE='auto_20240601_14_00_00'}
| SELECT * FROM Albums
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Albums [at 2:15] SELECT * FROM Albums               ^"}
> query
| @{USE_ADDITIONAL_PARALLELISM=TRUE}
| SELECT COUNT(*) FROM LargeTable
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: LargeTable [at 2:22] SELECT COUNT(*) FROM LargeTable                      ^"}
> query
| @{PDML_MAX_PARALLELISM=4}
| DELETE FROM OldLogs WHERE CreatedAt < '2023-01-01'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"dml","reason":"semantic","code":"InvalidArgument","message":"Table not found: OldLogs [at 2:13] DELETE FROM OldLogs WHERE CreatedAt \u003c '2023-01-01'             ^"}
> query
| SELECT *
| FROM Singers s
| INNER JOIN @{JOIN_METHOD=HASH_JOIN} Albums a ON s.SingerId = a.SingerId
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 2:6] FROM Singers s      ^"}
> query
| SELECT *
| FROM Singers s
| LEFT JOIN @{JOIN_METHOD=APPLY_JOIN} Albums a ON s.SingerId = a.SingerId
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 2:6] FROM Singers s      ^"}
> query
| @{LOCK_SCANNED_RANGES=exclusive}
| SELECT * FROM Accounts WHERE AccountId = 42 FOR UPDATE
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Accounts [at 2:15] SELECT * FROM Accounts WHERE AccountId = 42 FOR UPDATE               ^"}
> query
| @{LOCK_SCANNED_RANGES=shared}
| UPDATE Counters SET Value = Value + 1 WHERE CounterId = 1
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"dml","reason":"semantic","code":"InvalidArgument","message":"Table not found: Counters [at 2:8] UPDATE Counters SET Value = Value + 1 WHERE CounterId = 1        ^"}
> query
| SELECT * FROM LargeTable @{SCAN_METHOD=BATCH}
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: LargeTable [at 1:15] SELECT * FROM LargeTable @{SCAN_METHOD=BATCH}               ^"}
> query
| @{OPTIMIZER_VERSION=4, USE_ADDITIONAL_PARALLELISM=TRUE}
| SELECT * FROM Singers @{FORCE_INDEX=SingersByName}
| JOIN @{JOIN_METHOD=HASH_JOIN} Albums a ON Singers.SingerId = a.SingerId
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 2:15] SELECT * FROM Singers @{FORCE_INDEX=SingersByName}               ^"}
> query
| SELECT SingerId, first_name, MyColumn FROM MyTable
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: MyTable [at 1:44] SELECT SingerId, first_name, MyColumn FROM MyTable                                            ^"}
> query
| SELECT `SELECT` FROM `my table`
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: `my table` [at 1:22] SELECT `SELECT` FROM `my table`                      ^"}
> query
| SELECT `my-column` FROM `schema.name`
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: `schema.name` [at 1:25] SELECT `my-column` FROM `schema.name`                         ^"}
> query
| CREATE TABLE `order` (id INT64) PRIMARY KEY (id)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| SELECT "Hello, World!"
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT "She said \"hi\""
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT '''Line1
| Line2
| Line3'''
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT """Hello
| World"""
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT r'abc\ndef'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| -- literal backslash-n, NOT a newline
| SELECT r"a+b*c"
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT R'''no\escape'''
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT b'abc'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT B'\x00\xFF'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT b"""binary\x01data"""
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT rb'no\escape'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT br'\xFF'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT '\n\t\r'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT '\x41'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| -- 'A'
| SELECT 'A'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| -- 'A'
| SELECT '\U00000041'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT 42, -100, +7
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT 0x1A, 0xFF, 0XABC123
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT 3.14, .5, 1.
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT 1.23e10, 4.5E-3, 1e2
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT -0.001, +1.5
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT NUMERIC '123.456'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT NUMERIC '-9999999999.999999999'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT CAST('+inf' AS FLOAT64)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT CAST('-Inf' AS FLOAT32)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT CAST('NaN' AS FLOAT64)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT TRUE, FALSE, true, false
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT NULL
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT COALESCE(NULL, 'default')
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT DATE '2024-06-01'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT TIMESTAMP '2024-06-01 12:00:00 UTC'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT TIMESTAMP '2024-06-01T12:00:00.000000000Z'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT JSON '{"name": "Alice", "age": 30}'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT * FROM Singers WHERE SingerId = @singer_id
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:15] SELECT * FROM Singers WHERE SingerId = @singer_id               ^"}
> query
| SELECT * FROM Albums WHERE SingerId = @id AND Title = @title
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Albums [at 1:15] SELECT * FROM Albums WHERE SingerId = @id AND Title = @title               ^"}
> query
| INSERT INTO Singers (SingerId, FirstName) VALUES (@id, @name)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"dml","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:13] INSERT INTO Singers (SingerId, FirstName) VALUES (@id, @name)             ^"}
> query
| -- This is a single-line comment
| SELECT 1
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| -- inline comment
|
| /* This is a
|    multi-line comment */
| SELECT /* inline block comment */ 2
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| -- Using reserved keyword as identifier requires backticks
| CREATE TABLE `order` (id INT64 NOT NULL) PRIMARY KEY (id)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| SELECT `select` FROM `from`
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: `from` [at 1:22] SELECT `select` FROM `from`                      ^"}
> query
| -- All equivalent:
| SELECT singerid FROM singers
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: singers [at 2:22] SELECT singerid FROM singers                      ^"}
> query
| SELECT SingerId FROM Singers
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:22] SELECT SingerId FROM Singers                      ^"}
> query
| SELECT SINGERID FROM SINGERS
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: SINGERS [at 1:22] SELECT SINGERID FROM SINGERS                      ^"}
> query
| CREATE TABLE A (id INT64) PRIMARY KEY (id)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| CREATE TABLE B (id INT64) PRIMARY KEY (id)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"ddl","reason":"ok","code":"OK"}
> query
| SELECT SingerId, FirstName FROM Singers ORDER BY LastName LIMIT 10 OFFSET 5
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:33] SELECT SingerId, FirstName FROM Singers ORDER BY LastName LIMIT 10 OFFSET 5                                 ^"}
> query
| (SELECT 1 AS n UNION ALL SELECT 2 AS n) ORDER BY n DESC
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT *
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"SELECT * must have a FROM clause [at 1:8] SELECT *        ^"}
> query
| SELECT * EXCEPT (SingerInfo)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"SELECT * must have a FROM clause [at 1:8] SELECT * EXCEPT (SingerInfo)        ^"}
> query
| SELECT * REPLACE (UPPER(LastName) AS LastName)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"SELECT * must have a FROM clause [at 1:8] SELECT * REPLACE (UPPER(LastName) AS LastName)        ^"}
> query
| SELECT DISTINCT LastName FROM Singers
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:31] SELECT DISTINCT LastName FROM Singers                               ^"}
> query
| SELECT AS STRUCT SingerId, FirstName, LastName FROM Singers
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:53] SELECT AS STRUCT SingerId, FirstName, LastName FROM Singers                                                     ^"}
> query
| SELECT AS VALUE SingerId FROM Singers
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:31] SELECT AS VALUE SingerId FROM Singers                               ^"}
> query
| SELECT s.FirstName FROM Singers AS s
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:25] SELECT s.FirstName FROM Singers AS s                         ^"}
> query
| SELECT a.Title FROM myschema.Albums AS a
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: myschema.Albums [at 1:21] SELECT a.Title FROM myschema.Albums AS a                     ^"}
> query
| SELECT * FROM (SELECT 1 AS x) AS sub
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT * FROM Singers s CROSS JOIN Albums a
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:15] SELECT * FROM Singers s CROSS JOIN Albums a               ^"}
> query
| SELECT s.FirstName, a.Title
| FROM Singers s INNER JOIN Albums a ON s.SingerId = a.SingerId
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 2:6] FROM Singers s INNER JOIN Albums a ON s.SingerId = a.SingerId      ^"}
> query
| SELECT s.FirstName, a.Title
| FROM Singers s LEFT OUTER JOIN Albums a USING (SingerId)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 2:6] FROM Singers s LEFT OUTER JOIN Albums a USING (SingerId)      ^"}
> query
| SELECT s.FirstName, a.Title
| FROM Singers s FULL OUTER JOIN Albums a ON s.SingerId = a.SingerId
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 2:6] FROM Singers s FULL OUTER JOIN Albums a ON s.SingerId = a.SingerId      ^"}
> query
| -- HASH join hint
| SELECT * FROM Singers s HASH JOIN Albums a ON s.SingerId = a.SingerId
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 2:15] SELECT * FROM Singers s HASH JOIN Albums a ON s.SingerId = a.SingerId               ^"}
> query
| SELECT * FROM Singers WHERE LastName = 'Richards'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:15] SELECT * FROM Singers WHERE LastName = 'Richards'               ^"}
> query
| SELECT * FROM Albums WHERE MarketingBudget > 10000 AND SingerId = 1
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Albums [at 1:15] SELECT * FROM Albums WHERE MarketingBudget \u003e 10000 AND SingerId = 1               ^"}
> query
| SELECT LastName, COUNT(*) FROM Singers GROUP BY LastName
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:32] SELECT LastName, COUNT(*) FROM Singers GROUP BY LastName                                ^"}
> query
| SELECT SingerId, AlbumId, SUM(Revenue) FROM Sales GROUP BY ROLLUP (SingerId, AlbumId)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Sales [at 1:45] SELECT SingerId, AlbumId, SUM(Revenue) FROM Sales GROUP BY ROLLUP (SingerId, ...                                             ^"}
> query
| SELECT LastName, COUNT(*) AS cnt
| FROM Singers
| GROUP BY LastName
| HAVING cnt > 1
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 2:6] FROM Singers      ^"}
> query
| SELECT * FROM Singers ORDER BY LastName ASC, FirstName DESC
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:15] SELECT * FROM Singers ORDER BY LastName ASC, FirstName DESC               ^"}
> query
| SELECT * FROM Singers ORDER BY BirthDate NULLS LAST
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:15] SELECT * FROM Singers ORDER BY BirthDate NULLS LAST               ^"}
> query
| SELECT * FROM Singers LIMIT 5
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:15] SELECT * FROM Singers LIMIT 5               ^"}
> query
| SELECT * FROM Singers LIMIT 10 OFFSET 20
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:15] SELECT * FROM Singers LIMIT 10 OFFSET 20               ^"}
> query
| SELECT * FROM Singers WHERE SingerId = 1 FOR UPDATE
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:15] SELECT * FROM Singers WHERE SingerId = 1 FOR UPDATE               ^"}
> query
| SELECT SingerId FROM Singers
| UNION ALL
| SELECT SingerId FROM ArchivedSingers
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:22] SELECT SingerId FROM Singers                      ^"}
> query
| SELECT LastName FROM Singers
| UNION DISTINCT
| SELECT LastName FROM ArchivedSingers
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:22] SELECT LastName FROM Singers                      ^"}
> query
| SELECT SingerId FROM Singers
| INTERSECT DISTINCT
| SELECT SingerId FROM Albums
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:22] SELECT SingerId FROM Singers                      ^"}
> query
| SELECT SingerId FROM Singers
| EXCEPT DISTINCT
| SELECT SingerId FROM BlacklistSingers
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:22] SELECT SingerId FROM Singers                      ^"}
> query
| WITH TopSingers AS (
|   SELECT SingerId, FirstName, LastName
|   FROM Singers
|   ORDER BY LastName
|   LIMIT 10
| )
| SELECT * FROM TopSingers WHERE FirstName LIKE 'A%'
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 3:8]   FROM Singers        ^"}
> query
| -- Multiple CTEs
| WITH
|   Cte1 AS (SELECT 1 AS n),
|   Cte2 AS (SELECT n + 1 AS n FROM Cte1)
| SELECT * FROM Cte2
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| -- Scalar subquery
| SELECT (SELECT MAX(MarketingBudget) FROM Albums) AS MaxBudget
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Albums [at 2:42] SELECT (SELECT MAX(MarketingBudget) FROM Albums) AS MaxBudget                                          ^"}
> query
| -- EXISTS
| SELECT * FROM Singers WHERE EXISTS (SELECT 1 FROM Albums WHERE Albums.SingerId = Singers.SingerId)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 2:15] SELECT * FROM Singers WHERE EXISTS (SELECT 1 FROM Albums WHERE Albums.SingerI...               ^"}
> query
| -- IN subquery
| SELECT * FROM Singers WHERE SingerId IN (SELECT SingerId FROM Albums WHERE Title LIKE 'A%')
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 2:15] SELECT * FROM Singers WHERE SingerId IN (SELECT SingerId FROM Albums WHERE Ti...               ^"}
> query
| -- FROM subquery
| SELECT * FROM (SELECT SingerId, COUNT(*) AS albums FROM Albums GROUP BY SingerId) AS sub
| WHERE albums > 2
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Albums [at 2:57] ...FROM (SELECT SingerId, COUNT(*) AS albums FROM Albums GROUP BY SingerId) A...                                                   ^"}
> query
| SELECT item FROM UNNEST(['apple', 'banana', 'cherry']) AS item
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| SELECT item, pos
| FROM UNNEST(['apple', 'banana', 'cherry']) AS item WITH OFFSET AS pos
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
> query
| -- Unnesting an array column
| SELECT SingerId, album
| FROM Singers, UNNEST(AlbumList) AS album
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 3:6] FROM Singers, UNNEST(AlbumList) AS album      ^"}
> query
| SELECT s.Name.FirstName, s.Name.LastName FROM StructTable AS s
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: StructTable [at 1:47] SELECT s.Name.FirstName, s.Name.LastName FROM StructTable AS s                                               ^"}
> query
| SELECT arr_elem FROM Table t, t.ArrayColumn arr_elem
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Table [at 1:22] SELECT arr_elem FROM Table t, t.ArrayColumn arr_elem                      ^"}
> query
| SELECT * FROM Singers TABLESAMPLE BERNOULLI (10 PERCENT)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:15] SELECT * FROM Singers TABLESAMPLE BERNOULLI (10 PERCENT)               ^"}
> query
| SELECT * FROM Albums TABLESAMPLE RESERVOIR (100 ROWS)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Albums [at 1:15] SELECT * FROM Albums TABLESAMPLE RESERVOIR (100 ROWS)               ^"}
> query
| SELECT s.FirstName, a.Title
| FROM Singers s, Albums a
| WHERE s.SingerId = a.SingerId
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 2:6] FROM Singers s, Albums a      ^"}
> query
| SELECT AS STRUCT SingerId, FirstName, LastName FROM Singers LIMIT 1
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:53] SELECT AS STRUCT SingerId, FirstName, LastName FROM Singers LIMIT 1                                                     ^"}
> query
| SELECT AS VALUE SingerId FROM Singers
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Singers [at 1:31] SELECT AS VALUE SingerId FROM Singers                               ^"}
> query
| SELECT SingerId, AlbumId,
|   ROW_NUMBER() OVER (PARTITION BY SingerId ORDER BY MarketingBudget DESC) AS rn
| FROM Albums
| QUALIFY rn = 1
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"QUALIFY is not supported [at 4:1] QUALIFY rn = 1 ^"}
> query
| SELECT SingerId, AlbumId, MarketingBudget,
|   AVG(MarketingBudget) OVER w AS AvgBudget
| FROM Albums
| WINDOW w AS (PARTITION BY SingerId ORDER BY AlbumId
|              ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"semantic","code":"InvalidArgument","message":"Table not found: Albums [at 3:6] FROM Albums      ^"}
> query
| (SELECT 1 AS x) UNION ALL (SELECT 2 AS x)
< result
< column "verdict" "JSON"
< row
: string
| {"verdict":"accept","kind":"query","reason":"ok","code":"OK"}
