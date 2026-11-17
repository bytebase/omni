--
-- Table structure for `.events`
--
CREATE TABLE omni_testdump.`.events`
(
    `id` UInt64
)
ENGINE = MergeTree
ORDER BY id
SETTINGS index_granularity = 8192;

--
-- Table structure for `.innerXid.events`
--
CREATE TABLE omni_testdump.`.innerXid.events`
(
    `id` UInt64
)
ENGINE = MergeTree
ORDER BY id
SETTINGS index_granularity = 8192;

--
-- Table structure for `line\nbreak`
--
CREATE TABLE omni_testdump.`line\nbreak`
(
    `id` UInt64
)
ENGINE = MergeTree
ORDER BY id
SETTINGS index_granularity = 8192;

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
-- Table structure for `src`
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
-- Table structure for `target`
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
-- Dictionary structure for `dict`
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
-- Dictionary structure for `mv_dict`
--
CREATE DICTIONARY omni_testdump.mv_dict
(
    `id` UInt64
)
PRIMARY KEY id
SOURCE(CLICKHOUSE(TABLE 'z_src_mv'))
LIFETIME(MIN 0 MAX 300)
LAYOUT(FLAT());

--
-- Dictionary structure for `view_dict`
--
CREATE DICTIONARY omni_testdump.view_dict
(
    `id` UInt64
)
PRIMARY KEY id
SOURCE(CLICKHOUSE(TABLE 'v_source'))
LIFETIME(MIN 0 MAX 300)
LAYOUT(FLAT());

--
-- View structure for `a_first`
--
CREATE VIEW omni_testdump.a_first
(
    `id` UInt64
)
AS SELECT id
FROM omni_testdump.v_names;

--
-- View structure for `v_bridge`
--
CREATE VIEW omni_testdump.v_bridge
(
    `id` UInt64,
    `name` String
)
AS SELECT
    id,
    name
FROM omni_testdump.mv_base;

--
-- View structure for `v_names`
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
-- View structure for `v_source`
--
CREATE VIEW omni_testdump.v_source
(
    `id` UInt64
)
AS SELECT id
FROM omni_testdump.src;

--
-- Table structure for `dict_table`
--
CREATE TABLE omni_testdump.dict_table
(
    `id` UInt64
)
ENGINE = Dictionary('omni_testdump.dict');

--
-- Materialized view structure for `mv_base`
--
CREATE MATERIALIZED VIEW omni_testdump.mv_base TO omni_testdump.target
(
    `id` UInt64,
    `name` String
)
AS SELECT
    id,
    'omni_testdump.a_top' AS name
FROM omni_testdump.src;

--
-- Materialized view structure for `mv_own`
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
-- Materialized view structure for `nl\nmv`
--
CREATE MATERIALIZED VIEW omni_testdump.`nl\nmv` TO omni_testdump.target
(
    `id` UInt64,
    `name` String
)
AS SELECT
    id,
    'x' AS name
FROM omni_testdump.src;

--
-- Table structure for `view_dict_table`
--
CREATE TABLE omni_testdump.view_dict_table
(
    `id` UInt64
)
ENGINE = Dictionary('omni_testdump.view_dict');

--
-- Materialized view structure for `z_src_mv`
--
CREATE MATERIALIZED VIEW omni_testdump.z_src_mv TO omni_testdump.target
(
    `id` UInt64,
    `name` String
)
AS SELECT
    ` omni_testdump.a_mv_dict_table `.id,
    'm' AS name
FROM omni_testdump.src AS ` omni_testdump.a_mv_dict_table `;

--
-- Materialized view structure for `a_after`
--
CREATE MATERIALIZED VIEW omni_testdump.a_after TO omni_testdump.target
(
    `id` UInt64,
    `name` String
)
AS SELECT
    id,
    'y' AS name
FROM omni_testdump.`nl\nmv`;

--
-- Materialized view structure for `a_dict`
--
CREATE MATERIALIZED VIEW omni_testdump.a_dict TO omni_testdump.target
(
    `id` UInt64,
    `name` String
)
AS SELECT
    id,
    'z' AS name
FROM omni_testdump.dict_table;

--
-- Table structure for `a_mv_dict_table`
--
CREATE TABLE omni_testdump.a_mv_dict_table
(
    `id` UInt64
)
ENGINE = Dictionary('omni_testdump.mv_dict');

--
-- Materialized view structure for `a_top`
--
CREATE MATERIALIZED VIEW omni_testdump.a_top TO omni_testdump.target
(
    `id` UInt64,
    `name` String
)
AS SELECT
    id,
    name
FROM omni_testdump.v_bridge;

