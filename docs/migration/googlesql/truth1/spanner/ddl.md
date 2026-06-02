# Cloud Spanner GoogleSQL truth1 — DDL

All forms extracted from official Cloud Spanner GoogleSQL DDL reference.
Primary source: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language
version: spanner-current

---

## DDL-001: CREATE DATABASE

```
CREATE DATABASE database_id

database_id: {a–z}[{a–z|0–9|_|-}+]{a–z|0–9}
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#create_database

```sql
CREATE DATABASE my_database;
```

---

## DDL-002: ALTER DATABASE options

```
ALTER DATABASE database_id
  SET OPTIONS ( options_def [, ...] )

options_def:
    default_leader = { 'region' | null }
  | optimizer_version = { integer | null }
  | optimizer_statistics_package = { 'package_name' | null }
  | version_retention_period = { 'duration' | null }
  | default_sequence_kind = { 'bit_reversed_positive' | null }
  | default_time_zone = { 'time_zone_name' | null }
  | read_lease_regions = { 'region_name [, ...]' | null }
  | columnar_policy = { 'columnar_policy_value' | null }
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#alter_database

```sql
ALTER DATABASE my_database SET OPTIONS (default_leader = 'us-east1');
ALTER DATABASE my_database SET OPTIONS (version_retention_period = '7d', optimizer_version = 3);
ALTER DATABASE my_database SET OPTIONS (default_sequence_kind = 'bit_reversed_positive');
```

---

## DDL-003: CREATE TABLE

```
CREATE TABLE [IF NOT EXISTS] [[schema_name.]table_name] (
  column_def [, ...]
  [, table_constraint [, ...]]
)
PRIMARY KEY ( key_part [, ...] )
[, INTERLEAVE IN PARENT [schema_name.]table_name [ON DELETE { CASCADE | NO ACTION }]]
[, ROW DELETION POLICY ( OLDER_THAN ( timestamp_column, INTERVAL interval_expr ) )]
[OPTIONS ( table_option [, ...] )]

column_def:
  column_name { scalar_type | array_type }
  [NOT NULL]
  [DEFAULT ( expression )]
  [AS ( expression ) { STORED | VIRTUAL }]
  [REFERENCES [schema_name.]ref_table ( ref_column ) [ON DELETE { CASCADE | NO ACTION }]]
  [CONSTRAINT constraint_name] [CHECK ( expression )]
  [OPTIONS ( column_option [, ...] )]

scalar_type:
    BOOL
  | INT64
  | FLOAT32
  | FLOAT64
  | NUMERIC
  | STRING ( { length | MAX } )
  | BYTES ( { length | MAX } )
  | JSON
  | DATE
  | TIMESTAMP
  | TOKENLIST
  | { proto_type_name }

array_type:
  ARRAY < scalar_type >
  [ ( vector_length => integer_literal ) ]

key_part:
  column_name [{ ASC | DESC }]

table_constraint:
    [CONSTRAINT constraint_name]
    { PRIMARY KEY ( key_part [, ...] )
      | FOREIGN KEY ( column_name [, ...] )
          REFERENCES [schema_name.]ref_table ( ref_column [, ...] )
          [ON DELETE { CASCADE | NO ACTION }]
      | CHECK ( expression )
    }

table_option:
    row_deletion_policy = { ROW IS OLDER THAN ... }
  | (reserved for future OPTIONS keys)

column_option:
  allow_commit_timestamp = { true | false }
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#create_table

```sql
CREATE TABLE Singers (
  SingerId   INT64 NOT NULL,
  FirstName  STRING(1024),
  LastName   STRING(1024),
  BirthDate  DATE,
  SingerInfo BYTES(MAX),
) PRIMARY KEY (SingerId);

CREATE TABLE Albums (
  SingerId     INT64 NOT NULL,
  AlbumId      INT64 NOT NULL,
  Title        STRING(MAX),
  LastUpdated  TIMESTAMP OPTIONS (allow_commit_timestamp = true),
) PRIMARY KEY (SingerId, AlbumId),
  INTERLEAVE IN PARENT Singers ON DELETE CASCADE;

CREATE TABLE PlayerStats (
  PlayerId   INT64 NOT NULL,
  CreatedAt  TIMESTAMP NOT NULL,
  Score      INT64,
) PRIMARY KEY (PlayerId, CreatedAt),
  ROW DELETION POLICY (OLDER_THAN(CreatedAt, INTERVAL 30 DAY));
```

---

## DDL-004: Column DEFAULT value

```
column_name type DEFAULT ( expression )
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#column_default

```sql
CREATE TABLE Orders (
  OrderId     INT64 NOT NULL,
  Status      STRING(20) DEFAULT ('pending'),
  CreatedAt   TIMESTAMP DEFAULT (CURRENT_TIMESTAMP()),
) PRIMARY KEY (OrderId);
```

---

## DDL-005: Generated column (AS ... STORED / VIRTUAL)

```
column_name type [NOT NULL] AS ( expression ) { STORED | VIRTUAL }
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#generated_column

```sql
CREATE TABLE Employees (
  EmpId      INT64 NOT NULL,
  FirstName  STRING(100),
  LastName   STRING(100),
  FullName   STRING(200) AS (CONCAT(FirstName, ' ', LastName)) STORED,
  NameLen    INT64 AS (CHAR_LENGTH(LastName)) VIRTUAL,
) PRIMARY KEY (EmpId);
```

---

## DDL-006: allow_commit_timestamp column option

```
column_name TIMESTAMP [NOT NULL] OPTIONS (allow_commit_timestamp = { true | false })
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#allow_commit_timestamp

```sql
CREATE TABLE Events (
  EventId    INT64 NOT NULL,
  LastMod    TIMESTAMP OPTIONS (allow_commit_timestamp = true),
) PRIMARY KEY (EventId);
```

---

## DDL-007: INTERLEAVE IN PARENT clause

```
INTERLEAVE IN PARENT [schema_name.]parent_table_name
  [ON DELETE { CASCADE | NO ACTION }]
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#interleave_in_parent

```sql
CREATE TABLE Chapters (
  SingerId  INT64 NOT NULL,
  AlbumId   INT64 NOT NULL,
  ChapterId INT64 NOT NULL,
  Title     STRING(MAX),
) PRIMARY KEY (SingerId, AlbumId, ChapterId),
  INTERLEAVE IN PARENT Albums ON DELETE CASCADE;

-- NO ACTION (default)
CREATE TABLE Details (
  ParentId INT64 NOT NULL,
  DetailId INT64 NOT NULL,
) PRIMARY KEY (ParentId, DetailId),
  INTERLEAVE IN PARENT Parent ON DELETE NO ACTION;
```

---

## DDL-008: FOREIGN KEY constraint

```
[CONSTRAINT constraint_name]
FOREIGN KEY ( column_name [, ...] )
  REFERENCES [schema_name.]ref_table ( ref_column [, ...] )
  [ON DELETE { CASCADE | NO ACTION }]
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#foreign_key

```sql
CREATE TABLE Orders (
  OrderId    INT64 NOT NULL,
  CustomerId INT64 NOT NULL,
  CONSTRAINT FK_Customers
    FOREIGN KEY (CustomerId) REFERENCES Customers (CustomerId) ON DELETE NO ACTION,
) PRIMARY KEY (OrderId);
```

---

## DDL-009: CHECK constraint

```
[CONSTRAINT constraint_name] CHECK ( bool_expression )
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#check_constraint

```sql
CREATE TABLE Products (
  ProductId INT64 NOT NULL,
  Price     FLOAT64,
  CONSTRAINT chk_positive_price CHECK (Price > 0),
) PRIMARY KEY (ProductId);
```

---

## DDL-010: ROW DELETION POLICY (TTL)

```
ROW DELETION POLICY ( OLDER_THAN ( timestamp_column, INTERVAL n { DAY | HOUR | ... } ) )
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#row_deletion_policy

```sql
CREATE TABLE Sessions (
  SessionId   INT64 NOT NULL,
  CreatedAt   TIMESTAMP NOT NULL,
) PRIMARY KEY (SessionId),
  ROW DELETION POLICY (OLDER_THAN(CreatedAt, INTERVAL 90 DAY));
```

---

## DDL-011: ALTER TABLE — add column

```
ALTER TABLE [schema_name.]table_name
  ADD COLUMN [IF NOT EXISTS] column_def
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#alter_table

```sql
ALTER TABLE Singers ADD COLUMN Nickname STRING(100);
ALTER TABLE Singers ADD COLUMN IF NOT EXISTS Bio STRING(MAX);
```

---

## DDL-012: ALTER TABLE — drop column

```
ALTER TABLE [schema_name.]table_name
  DROP COLUMN [IF EXISTS] column_name
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#alter_table

```sql
ALTER TABLE Singers DROP COLUMN Nickname;
```

---

## DDL-013: ALTER TABLE — alter column

```
ALTER TABLE [schema_name.]table_name
  ALTER COLUMN column_name
    { type [NOT NULL]
      | SET DEFAULT ( expression )
      | DROP DEFAULT
      | SET OPTIONS ( column_option [, ...] )
    }
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#alter_table

```sql
ALTER TABLE Singers ALTER COLUMN FirstName STRING(2048);
ALTER TABLE Singers ALTER COLUMN FirstName SET NOT NULL;
ALTER TABLE Events ALTER COLUMN LastMod SET OPTIONS (allow_commit_timestamp = true);
ALTER TABLE Orders ALTER COLUMN Status SET DEFAULT ('new');
ALTER TABLE Orders ALTER COLUMN Status DROP DEFAULT;
```

---

## DDL-014: ALTER TABLE — add/drop constraint

```
ALTER TABLE [schema_name.]table_name
  { ADD [CONSTRAINT constraint_name] { FOREIGN KEY ... | CHECK ... | PRIMARY KEY ... }
  | DROP CONSTRAINT [IF EXISTS] constraint_name
  }
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#alter_table

```sql
ALTER TABLE Orders ADD CONSTRAINT FK_Cust FOREIGN KEY (CustomerId) REFERENCES Customers (CustomerId);
ALTER TABLE Orders DROP CONSTRAINT FK_Cust;
```

---

## DDL-015: ALTER TABLE — set options

```
ALTER TABLE [schema_name.]table_name
  SET OPTIONS ( table_option [, ...] )
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#alter_table

```sql
ALTER TABLE Sessions SET OPTIONS (row_deletion_policy = null);
```

---

## DDL-016: ALTER TABLE — row deletion policy

```
ALTER TABLE [schema_name.]table_name
  { ADD ROW DELETION POLICY ( OLDER_THAN ( col, INTERVAL n DAY ) )
  | ALTER ROW DELETION POLICY ( OLDER_THAN ( col, INTERVAL n DAY ) )
  | DROP ROW DELETION POLICY
  }
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#alter_table

```sql
ALTER TABLE Sessions ADD ROW DELETION POLICY (OLDER_THAN(CreatedAt, INTERVAL 30 DAY));
ALTER TABLE Sessions ALTER ROW DELETION POLICY (OLDER_THAN(CreatedAt, INTERVAL 60 DAY));
ALTER TABLE Sessions DROP ROW DELETION POLICY;
```

---

## DDL-017: ALTER TABLE — rename

```
ALTER TABLE [schema_name.]table_name
  RENAME TO new_table_name
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#alter_table

```sql
ALTER TABLE Singers RENAME TO Artists;
```

---

## DDL-018: DROP TABLE

```
DROP TABLE [IF EXISTS] [schema_name.]table_name
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#drop_table

```sql
DROP TABLE Singers;
DROP TABLE IF EXISTS temp_table;
```

---

## DDL-019: CREATE INDEX

```
CREATE [UNIQUE] [NULL_FILTERED] INDEX [IF NOT EXISTS] index_name
  ON [schema_name.]table_name ( key_part [, ...] )
  [STORING ( stored_column [, ...] )]
  [, INTERLEAVE IN [schema_name.]parent_table_name]
  [WHERE filter_expression]

key_part:
  column_name [{ ASC | DESC }]
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#create_index

```sql
CREATE INDEX SingersByName ON Singers(LastName, FirstName);
CREATE UNIQUE INDEX UniqueEmail ON Users(Email);
CREATE NULL_FILTERED INDEX AlbumsByTitle ON Albums(Title);
CREATE INDEX AlbumsByArtist ON Albums(SingerId, Title) STORING (MarketingBudget);
CREATE INDEX SongsBySinger ON Songs(SingerId, AlbumId, TrackId)
  STORING (SongName), INTERLEAVE IN Albums;
```

---

## DDL-020: ALTER INDEX

```
ALTER INDEX [schema_name.]index_name
  { ADD STORED COLUMN column_name
  | DROP STORED COLUMN column_name
  }
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#alter_index

```sql
ALTER INDEX SingersByName ADD STORED COLUMN BirthDate;
ALTER INDEX SingersByName DROP STORED COLUMN BirthDate;
```

---

## DDL-021: DROP INDEX

```
DROP INDEX [IF EXISTS] [schema_name.]index_name
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#drop_index

```sql
DROP INDEX SingersByName;
DROP INDEX IF EXISTS SingersByName;
```

---

## DDL-022: CREATE VIEW

```
CREATE [OR REPLACE] VIEW [schema_name.]view_name
  SQL SECURITY { INVOKER | DEFINER }
  AS query_statement
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#create_view

```sql
CREATE VIEW SingerView
  SQL SECURITY INVOKER
  AS SELECT SingerId, FirstName, LastName FROM Singers;

CREATE OR REPLACE VIEW TopAlbums
  SQL SECURITY DEFINER
  AS SELECT AlbumId, Title, SingerId FROM Albums WHERE MarketingBudget > 10000;
```

---

## DDL-023: DROP VIEW

```
DROP VIEW [IF EXISTS] [schema_name.]view_name
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#drop_view

```sql
DROP VIEW SingerView;
DROP VIEW IF EXISTS TopAlbums;
```

---

## DDL-024: CREATE CHANGE STREAM

```
CREATE CHANGE STREAM change_stream_name
  { FOR ALL
  | FOR table_and_column [, ...]
  }
  [OPTIONS ( change_stream_option [, ...] )]

table_and_column:
    table_name
  | table_name ( )
  | table_name ( column_name [, ...] )

change_stream_option:
    retention_period = { 'duration' | null }
  | value_capture_type = { 'OLD_AND_NEW_VALUES' | 'NEW_ROW' | 'NEW_VALUES' | null }
  | exclude_ttl_deletes = { true | false | null }
  | exclude_insert = { true | false | null }
  | exclude_update = { true | false | null }
  | exclude_delete = { true | false | null }
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#create_change_stream

```sql
CREATE CHANGE STREAM MyStream FOR ALL;

CREATE CHANGE STREAM SingersStream
  FOR Singers,
      Albums (Title, MarketingBudget)
  OPTIONS (retention_period = '7d', value_capture_type = 'NEW_ROW');
```

---

## DDL-025: ALTER CHANGE STREAM

```
ALTER CHANGE STREAM change_stream_name
  { SET FOR { ALL | table_and_column [, ...] }
  | DROP FOR ALL
  | SET OPTIONS ( change_stream_option [, ...] )
  }
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#alter_change_stream

```sql
ALTER CHANGE STREAM MyStream SET FOR Singers, Albums;
ALTER CHANGE STREAM MyStream SET OPTIONS (retention_period = '14d');
ALTER CHANGE STREAM MyStream DROP FOR ALL;
```

---

## DDL-026: DROP CHANGE STREAM

```
DROP CHANGE STREAM change_stream_name
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#drop_change_stream

```sql
DROP CHANGE STREAM MyStream;
```

---

## DDL-027: CREATE SEQUENCE

```
CREATE SEQUENCE [IF NOT EXISTS] [schema_name.]sequence_name
  [OPTIONS ( sequence_option [, ...] )]

sequence_option:
    sequence_kind = 'bit_reversed_positive'
  | skip_range_min = integer
  | skip_range_max = integer
  | start_with_counter = integer
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#create_sequence

```sql
CREATE SEQUENCE MySeq OPTIONS (
  sequence_kind = 'bit_reversed_positive'
);

CREATE SEQUENCE OrderSeq OPTIONS (
  sequence_kind = 'bit_reversed_positive',
  skip_range_min = 1,
  skip_range_max = 1000,
  start_with_counter = 5000
);
```

---

## DDL-028: ALTER SEQUENCE

```
ALTER SEQUENCE [IF EXISTS] [schema_name.]sequence_name
  SET OPTIONS ( sequence_option [, ...] )
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#alter_sequence

```sql
ALTER SEQUENCE MySeq SET OPTIONS (start_with_counter = 9000);
```

---

## DDL-029: DROP SEQUENCE

```
DROP SEQUENCE [IF EXISTS] [schema_name.]sequence_name
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#drop_sequence

```sql
DROP SEQUENCE MySeq;
DROP SEQUENCE IF EXISTS MySeq;
```

---

## DDL-030: CREATE SCHEMA

```
CREATE SCHEMA [IF NOT EXISTS] schema_name
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#create_schema

```sql
CREATE SCHEMA myschema;
CREATE SCHEMA IF NOT EXISTS analytics;
```

---

## DDL-031: DROP SCHEMA

```
DROP SCHEMA [IF EXISTS] schema_name
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#drop_schema

```sql
DROP SCHEMA myschema;
DROP SCHEMA IF EXISTS analytics;
```

---

## DDL-032: CREATE ROLE

```
CREATE ROLE role_name
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#create_role

```sql
CREATE ROLE analyst;
CREATE ROLE read_only_user;
```

---

## DDL-033: DROP ROLE

```
DROP ROLE role_name
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#drop_role

```sql
DROP ROLE analyst;
```

---

## DDL-034: GRANT (privileges to role)

```
GRANT { privilege [, ...] }
  ON { TABLE table_name [, ...]
       | VIEW view_name [, ...]
       | CHANGE STREAM change_stream_name [, ...]
       | TABLE FUNCTION tvf_name [, ...]
     }
  TO ROLE role_name [, ...]

privilege:
    SELECT
  | SELECT ( column_name [, ...] )
  | INSERT
  | INSERT ( column_name [, ...] )
  | UPDATE
  | UPDATE ( column_name [, ...] )
  | DELETE
  | EXECUTE
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#grant

```sql
GRANT SELECT ON TABLE Singers TO ROLE analyst;
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE Albums TO ROLE editor;
GRANT SELECT (SingerId, FirstName, LastName) ON TABLE Singers TO ROLE read_only_user;
GRANT EXECUTE ON TABLE FUNCTION get_singers TO ROLE analyst;
```

---

## DDL-035: GRANT (role to role)

```
GRANT ROLE role_name [, ...]
  TO ROLE target_role_name [, ...]
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#grant_role

```sql
GRANT ROLE analyst TO ROLE senior_analyst;
```

---

## DDL-036: REVOKE (privileges from role)

```
REVOKE { privilege [, ...] }
  ON { TABLE table_name [, ...]
       | VIEW view_name [, ...]
       | CHANGE STREAM change_stream_name [, ...]
       | TABLE FUNCTION tvf_name [, ...]
     }
  FROM ROLE role_name [, ...]
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#revoke

```sql
REVOKE SELECT ON TABLE Singers FROM ROLE analyst;
REVOKE INSERT, UPDATE ON TABLE Albums FROM ROLE editor;
```

---

## DDL-037: REVOKE (role from role)

```
REVOKE ROLE role_name [, ...]
  FROM ROLE target_role_name [, ...]
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#revoke_role

```sql
REVOKE ROLE analyst FROM ROLE senior_analyst;
```

---

## DDL-038: CREATE MODEL

```
CREATE [OR REPLACE] MODEL [IF NOT EXISTS] [schema_name.]model_name
  INPUT ( column_def [, ...] )
  OUTPUT ( column_def [, ...] )
  REMOTE
  OPTIONS ( model_option [, ...] )

model_option:
    endpoint = 'endpoint_url'
  | endpoints = ['endpoint_url', ...]
  | default_batch_size = integer
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#create_model

```sql
CREATE MODEL TextEmbedder
  INPUT (text STRING(MAX))
  OUTPUT (embeddings ARRAY<FLOAT64>)
  REMOTE
  OPTIONS (endpoints = ['//aiplatform.googleapis.com/projects/my_project/locations/us-central1/publishers/google/models/textembedding-gecko']);
```

---

## DDL-039: ALTER MODEL

```
ALTER MODEL [IF EXISTS] [schema_name.]model_name
  SET OPTIONS ( model_option [, ...] )
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#alter_model

```sql
ALTER MODEL TextEmbedder SET OPTIONS (default_batch_size = 5);
```

---

## DDL-040: DROP MODEL

```
DROP MODEL [IF EXISTS] [schema_name.]model_name
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#drop_model

```sql
DROP MODEL TextEmbedder;
DROP MODEL IF EXISTS MyModel;
```

---

## DDL-041: CREATE LOCALITY GROUP

```
CREATE LOCALITY GROUP locality_group_name
  [OPTIONS ( storage_option [, ...] )]

storage_option:
    storage = { 'ssd' | 'hdd' }
  | ssd_to_hdd_spill_timespan = 'duration'
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#create_locality_group

```sql
CREATE LOCALITY GROUP hot_data OPTIONS (storage = 'ssd');
CREATE LOCALITY GROUP cold_data OPTIONS (storage = 'hdd', ssd_to_hdd_spill_timespan = '30d');
```

---

## DDL-042: ALTER LOCALITY GROUP

```
ALTER LOCALITY GROUP locality_group_name
  SET OPTIONS ( storage_option [, ...] )
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#alter_locality_group

```sql
ALTER LOCALITY GROUP cold_data SET OPTIONS (ssd_to_hdd_spill_timespan = '14d');
```

---

## DDL-043: DROP LOCALITY GROUP

```
DROP LOCALITY GROUP locality_group_name
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#drop_locality_group

```sql
DROP LOCALITY GROUP cold_data;
```

---

## DDL-044: CREATE PLACEMENT

```
CREATE PLACEMENT placement_name
  [OPTIONS ( placement_option [, ...] )]

placement_option:
    instance_partition = 'partition_id'
  | default_leader = 'leader_region'
  | read_lease_regions = { 'region_name [, ...]' | null }
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#create_placement

```sql
CREATE PLACEMENT us_placement OPTIONS (instance_partition = 'us-partition');
```

---

## DDL-045: DROP PLACEMENT

```
DROP PLACEMENT placement_name
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#drop_placement

```sql
DROP PLACEMENT us_placement;
```

---

## DDL-046: CREATE PROTO BUNDLE

```
CREATE PROTO BUNDLE ( proto_type_name [, ...] )
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#create_proto_bundle

```sql
CREATE PROTO BUNDLE (
  `my.package.MyMessage`,
  `my.package.AnotherMessage`
);
```

---

## DDL-047: ALTER PROTO BUNDLE

```
ALTER PROTO BUNDLE
  { INSERT ( proto_type_name [, ...] )
  | UPDATE ( proto_type_name [, ...] )
  | DELETE ( proto_type_name [, ...] )
  } [, ...]
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#alter_proto_bundle

```sql
ALTER PROTO BUNDLE INSERT (`my.package.NewMessage`);
ALTER PROTO BUNDLE UPDATE (`my.package.MyMessage`);
ALTER PROTO BUNDLE DELETE (`my.package.OldMessage`);
```

---

## DDL-048: CREATE SEARCH INDEX

```
CREATE SEARCH INDEX [IF NOT EXISTS] index_name
  ON table_name ( tokenlist_column [, ...] )
  [STORING ( column [, ...] )]
  [PARTITION BY column [, ...]]
  [ORDER BY column [, ...]]
  [WHERE filter_expression]
  [INTERLEAVE IN parent_table]
  [OPTIONS ( search_index_option [, ...] )]

search_index_option:
    sort_order_sharding = { true | false }
  | disable_automatic_uid_column = { true | false }
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#create_search_index

```sql
-- Requires TOKENLIST column created via TOKENIZE_* functions
CREATE SEARCH INDEX SingerNameIndex ON Singers (SingerNameTokens);
CREATE SEARCH INDEX AlbumTitleIdx ON Albums (TitleTokens)
  STORING (MarketingBudget)
  PARTITION BY SingerId;
```

---

## DDL-049: ALTER SEARCH INDEX

```
ALTER SEARCH INDEX [IF EXISTS] index_name
  { ADD STORED COLUMN column_name
  | DROP STORED COLUMN column_name
  }
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#alter_search_index

```sql
ALTER SEARCH INDEX SingerNameIndex ADD STORED COLUMN FirstName;
ALTER SEARCH INDEX SingerNameIndex DROP STORED COLUMN FirstName;
```

---

## DDL-050: DROP SEARCH INDEX

```
DROP SEARCH INDEX [IF EXISTS] index_name
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#drop_search_index

```sql
DROP SEARCH INDEX SingerNameIndex;
```

---

## DDL-051: TOKENLIST column (for full-text / vector search)

```
column_name TOKENLIST [NOT NULL]
  [AS ( tokenize_expression ) STORED]
  [OPTIONS ( allow_commit_timestamp = { true | false } )]

-- Common tokenize functions used in generated TOKENLIST columns:
TOKENIZE_FULLTEXT(string_col [, language_tag => '...'])
TOKENIZE_NGRAMS(string_col, ngram_size_min => n, ngram_size_max => n)
TOKENIZE_NUMBER(numeric_col [, ...])
TOKENIZE_BOOL(bool_col)
TOKEN(value [, ...])
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#tokenlist_column

```sql
CREATE TABLE Products (
  ProductId    INT64 NOT NULL,
  Name         STRING(MAX),
  Description  STRING(MAX),
  NameTokens   TOKENLIST AS (TOKENIZE_FULLTEXT(Name)) STORED HIDDEN,
) PRIMARY KEY (ProductId);
```

---

## DDL-052: Vector index (VECTOR INDEX) — approximate nearest neighbor

```
CREATE VECTOR INDEX [IF NOT EXISTS] index_name
  ON table_name ( embedding_column )
  OPTIONS ( distance_type = { 'COSINE' | 'DOT_PRODUCT' | 'EUCLIDEAN' },
            num_leaves = integer
            [, num_leaves_to_search = integer]
          )

-- vector_length on ARRAY column definition:
column_name ARRAY<FLOAT32>(vector_length=>integer)
column_name ARRAY<FLOAT64>(vector_length=>integer)
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#create_vector_index

```sql
CREATE TABLE Products (
  ProductId  INT64 NOT NULL,
  Embedding  ARRAY<FLOAT32>(vector_length=>128),
) PRIMARY KEY (ProductId);

CREATE VECTOR INDEX ProductEmbeddingIdx
  ON Products (Embedding)
  OPTIONS (distance_type = 'COSINE', num_leaves = 100);
```

---

## DDL-053: DROP VECTOR INDEX

```
DROP VECTOR INDEX [IF EXISTS] index_name
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#drop_vector_index

```sql
DROP VECTOR INDEX ProductEmbeddingIdx;
```

---

## DDL-054: NEXT VALUE FOR (sequence usage in DDL/DML)

```
NEXT VALUE FOR [schema_name.]sequence_name
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#sequences

```sql
-- Used as a DEFAULT expression
CREATE TABLE Orders (
  OrderId INT64 DEFAULT (NEXT VALUE FOR OrderSeq) NOT NULL,
) PRIMARY KEY (OrderId);
```

---

## DDL-055: GET_NEXT_SEQUENCE_VALUE function

```
GET_NEXT_SEQUENCE_VALUE(SEQUENCE sequence_name)
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#sequences

```sql
-- Used as DEFAULT expression
CREATE TABLE Orders (
  OrderId INT64 DEFAULT (GET_NEXT_SEQUENCE_VALUE(SEQUENCE OrderSeq)) NOT NULL,
) PRIMARY KEY (OrderId);
```

---

## DDL-056: Table with schema-qualified name

```
CREATE TABLE schema_name.table_name ( ... ) PRIMARY KEY ( ... )
ALTER TABLE schema_name.table_name ...
DROP TABLE schema_name.table_name
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#named_schema

```sql
CREATE TABLE myschema.Employees (
  EmpId INT64 NOT NULL,
  Name  STRING(100),
) PRIMARY KEY (EmpId);

ALTER TABLE myschema.Employees ADD COLUMN Department STRING(50);
```

---

## DDL-057: Columnar storage (LOCALITY GROUP assignment on column)

```
column_name type [NOT NULL]
  OPTIONS ( locality_group = 'locality_group_name' )
```

source_url: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#locality_group

```sql
CREATE TABLE Metrics (
  MetricId   INT64 NOT NULL,
  HotValue   FLOAT64 OPTIONS (locality_group = 'hot_data'),
  ColdValue  FLOAT64 OPTIONS (locality_group = 'cold_data'),
) PRIMARY KEY (MetricId);
```
