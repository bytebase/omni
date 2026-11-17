-- The schema the dump tests read: a MergeTree table with an index, a primary
-- key, a setting and a comment; a table the materialized views write to; a view,
-- and a view over it under a name that sorts first; a materialized view, and one
-- over that under a name that sorts first; a dictionary; and a materialized view
-- with its own engine, whose inner table only the server names.
CREATE TABLE src (
  id UInt64,
  name String,
  INDEX name_bloom name TYPE bloom_filter GRANULARITY 4
)
ENGINE = MergeTree
ORDER BY (id, name)
PRIMARY KEY id
SETTINGS index_granularity = 4096
COMMENT 'the source';

CREATE TABLE target (
  id UInt64,
  name String
)
ENGINE = MergeTree
ORDER BY id;

-- A name ClickHouse writes in backquotes, so the dump quotes it back.
CREATE TABLE `odd-name` (
  id UInt64
)
ENGINE = MergeTree
ORDER BY id;

CREATE VIEW v_names AS SELECT id, name FROM src;

-- A view states its columns, which ClickHouse creates without resolving its
-- query, so this one replays before the view it reads.
CREATE VIEW a_first AS SELECT id FROM v_names;

CREATE MATERIALIZED VIEW mv_to TO target AS SELECT id, name FROM src;

-- A materialized view, and one over it under a name that sorts first: a
-- materialized view resolves its source when it is created, so the dump has to
-- write mv_base before a_top whatever their names say.
CREATE MATERIALIZED VIEW mv_base TO target AS SELECT id FROM src;

CREATE MATERIALIZED VIEW a_top TO target AS SELECT id FROM mv_base;

CREATE DICTIONARY dict (
  id UInt64,
  name String
)
PRIMARY KEY id
SOURCE(CLICKHOUSE(TABLE 'src'))
LAYOUT(FLAT())
LIFETIME(300);

CREATE MATERIALIZED VIEW mv_own ENGINE = MergeTree ORDER BY id AS SELECT id FROM src; -- not restorable: TestDumpGaps
