-- The schema the dump tests read: a MergeTree table with an index, a primary
-- key, a setting and a comment; a table the materialized views write to; a name
-- that needs backquotes; a user table whose name starts with the dot ClickHouse
-- uses for a materialized view's storage; a name that holds a line break; a
-- dictionary, and a table that reads it; a view, and a view over it under a name
-- that sorts first; two materialized views, the first holding the other's name
-- in a literal; and a materialized view with its own engine, whose inner table
-- only the server names.
--
-- $database stands for the database the test created, which the dictionary's
-- source and the materialized view's literal name.
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

CREATE TABLE `odd-name` (
  id UInt64
)
ENGINE = MergeTree
ORDER BY id;

CREATE TABLE `.events` (
  id UInt64
)
ENGINE = MergeTree
ORDER BY id;

-- A name the internal-table pattern must not take for one: its underscores are
-- what a LIKE pattern would read as wildcards.
CREATE TABLE `.innerXid.events` (
  id UInt64
)
ENGINE = MergeTree
ORDER BY id;

CREATE TABLE `line
break` (
  id UInt64
)
ENGINE = MergeTree
ORDER BY id;

CREATE DICTIONARY dict (
  id UInt64,
  name String
)
PRIMARY KEY id
SOURCE(CLICKHOUSE(TABLE 'src'))
LAYOUT(FLAT())
LIFETIME(300);

-- A dictionary whose source is a view, and a table that loads it: creating the
-- table queries the view, so both follow the views.
CREATE VIEW v_source AS SELECT id FROM src;

CREATE DICTIONARY view_dict (
  id UInt64
)
PRIMARY KEY id
SOURCE(CLICKHOUSE(TABLE 'v_source'))
LAYOUT(FLAT())
LIFETIME(300);

CREATE TABLE view_dict_table (
  id UInt64
)
ENGINE = Dictionary('$database.view_dict');

CREATE TABLE dict_table (
  id UInt64
)
ENGINE = Dictionary('$database.dict');

-- A dictionary whose source is a materialized view, and a table that loads it
-- under a name that sorts first: creating the table queries the view, so it
-- follows the view.
-- Its table alias is a reference with the spaces around it, in backquotes: a name
-- inside backquotes is part of that name, not a reference to a_mv_dict_table,
-- which loads the dictionary this view sources.
CREATE MATERIALIZED VIEW z_src_mv TO target AS SELECT ` $database.a_mv_dict_table `.id, 'm' AS name FROM src AS ` $database.a_mv_dict_table `;

CREATE DICTIONARY mv_dict (
  id UInt64
)
PRIMARY KEY id
SOURCE(CLICKHOUSE(TABLE 'z_src_mv'))
LAYOUT(FLAT())
LIFETIME(300);

CREATE TABLE a_mv_dict_table (
  id UInt64
)
ENGINE = Dictionary('$database.mv_dict');

CREATE VIEW v_names AS SELECT id, name FROM src;

-- A view states its columns, which ClickHouse creates without resolving its
-- query, so this one replays before the view it reads.
CREATE VIEW a_first AS SELECT id FROM v_names;

CREATE MATERIALIZED VIEW mv_base TO target AS SELECT id, '$database.a_top' AS name FROM src;


-- An ordinary view reads mv_base, and a materialized view reads that view: a
-- materialized view resolves its source when it is created, through the view, so
-- the dump has to write mv_base first whatever the names say. The literal above
-- names a_top and is not a read.
CREATE VIEW v_bridge AS SELECT id, name FROM mv_base;

CREATE MATERIALIZED VIEW a_top TO target AS SELECT id, name FROM v_bridge;

-- A materialized view whose name holds a line break, and one that reads it under
-- a name that sorts first: the reference the dump looks for is the rendering
-- ClickHouse writes, with the line break escaped.
CREATE MATERIALIZED VIEW `nl
mv` TO target AS SELECT id, 'x' AS name FROM src;

CREATE MATERIALIZED VIEW a_after TO target AS SELECT id, 'y' AS name FROM `nl
mv`;

-- A materialized view reading a table on the Dictionary engine, under a name that
-- sorts first: creating it needs the dictionary table, which loads the
-- dictionary.
CREATE MATERIALIZED VIEW a_dict TO target AS SELECT id, 'z' AS name FROM dict_table;

CREATE MATERIALIZED VIEW mv_own ENGINE = MergeTree ORDER BY id AS SELECT id FROM src; -- not restorable: TestDumpGaps
