--
-- Table structure for `odd-name`
--
CREATE TABLE omni_testdump.`odd-name`
(
    `id` UInt64
)
ENGINE = MergeTree
ORDER BY id
SETTINGS index_granularity = 8192;

--
-- Table structure for src
--
CREATE TABLE omni_testdump.src
(
    `id` UInt64,
    `name` String,
    INDEX name_bloom name TYPE bloom_filter GRANULARITY 4
)
ENGINE = MergeTree
PRIMARY KEY id
ORDER BY (id, name)
SETTINGS index_granularity = 4096
COMMENT 'the source';

--
-- Table structure for target
--
CREATE TABLE omni_testdump.target
(
    `id` UInt64,
    `name` String
)
ENGINE = MergeTree
ORDER BY id
SETTINGS index_granularity = 8192;

--
-- Dictionary structure for dict
--
CREATE DICTIONARY omni_testdump.dict
(
    `id` UInt64,
    `name` String
)
PRIMARY KEY id
SOURCE(CLICKHOUSE(TABLE 'src'))
LIFETIME(MIN 0 MAX 300)
LAYOUT(FLAT());

--
-- Materialized view structure for mv_base
--
CREATE MATERIALIZED VIEW omni_testdump.mv_base TO omni_testdump.target
(
    `id` UInt64
)
AS SELECT id
FROM omni_testdump.src;

--
-- Materialized view structure for mv_own
--
CREATE MATERIALIZED VIEW omni_testdump.mv_own
(
    `id` UInt64
)
ENGINE = MergeTree
ORDER BY id
SETTINGS index_granularity = 8192
AS SELECT id
FROM omni_testdump.src;

--
-- Materialized view structure for mv_to
--
CREATE MATERIALIZED VIEW omni_testdump.mv_to TO omni_testdump.target
(
    `id` UInt64,
    `name` String
)
AS SELECT
    id,
    name
FROM omni_testdump.src;

--
-- View structure for v_names
--
CREATE VIEW omni_testdump.v_names
(
    `id` UInt64,
    `name` String
)
AS SELECT
    id,
    name
FROM omni_testdump.src;

--
-- View structure for a_first
--
CREATE VIEW omni_testdump.a_first
(
    `id` UInt64
)
AS SELECT id
FROM omni_testdump.v_names;

--
-- Materialized view structure for a_top
--
CREATE MATERIALIZED VIEW omni_testdump.a_top TO omni_testdump.target
(
    `id` UInt64
)
AS SELECT id
FROM omni_testdump.mv_base;

