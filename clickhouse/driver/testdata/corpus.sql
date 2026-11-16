-- The schema the sync tests read. Nothing is inserted: a table's row count and
-- data size are the server's approximations, and an empty table makes both
-- zero in every recording.
CREATE TABLE events (
  id UInt64,
  name String COMMENT 'the event name',
  at DateTime64(3, 'UTC'),
  amount Nullable(Decimal(12, 2)),
  tags Array(String),
  attrs Map(String, String),
  fixed FixedString(4),
  flag Bool,
  level Enum8('info' = 1, 'warn' = 2),
  seen Date DEFAULT '2026-01-02'
)
ENGINE = MergeTree
ORDER BY (id, name)
PRIMARY KEY id
COMMENT 'events the sync reads';

CREATE TABLE skipped (
  url String,
  body String,
  INDEX url_bloom url TYPE bloom_filter GRANULARITY 4
)
ENGINE = MergeTree
ORDER BY url;

CREATE VIEW event_names AS SELECT id, name FROM events;
