-- pin: default clickhouse/clickhouse-server:26.8.38.9 sha256:5d30ac138948b43fd928842be7da1dcd2de712d32ba6c42368ff5009ba2072ff 26.8.38.9

== TestDump default#1
> exec
| -- The schema the dump tests read: a MergeTree table with an index, a primary
| -- key, a setting and a comment; a table the materialized views write to; a name
| -- that needs backquotes; a user table whose name starts with the dot ClickHouse
| -- uses for a materialized view's storage; a name that holds a line break; a
| -- dictionary, and a table that reads it; a view, and a view over it under a name
| -- that sorts first; two materialized views, the first holding the other's name
| -- in a literal; row policies, one of them on the whole database; and a
| -- materialized view with its own engine, whose inner table only the server
| -- names.
| --
| -- omni_testdump stands for the database the test created, which the dictionary's
| -- source, the materialized view's literal and a row policy name.
| CREATE TABLE src (
|   id UInt64,
|   name String,
|   INDEX name_bloom name TYPE bloom_filter GRANULARITY 4
| )
| ENGINE = MergeTree
| ORDER BY (id, name)
| PRIMARY KEY id
| SETTINGS index_granularity = 4096
| COMMENT 'the source';
< ok 0 -
> exec
|
|
| CREATE TABLE target (
|   id UInt64,
|   name String
| )
| ENGINE = MergeTree
| ORDER BY id;
< ok 0 -
> exec
|
|
| CREATE TABLE `odd-name` (
|   id UInt64
| )
| ENGINE = MergeTree
| ORDER BY id;
< ok 0 -
> exec
|
|
| CREATE TABLE `.events` (
|   id UInt64
| )
| ENGINE = MergeTree
| ORDER BY id;
< ok 0 -
> exec
|
|
| -- A name the internal-table pattern must not take for one: its underscores are
| -- what a LIKE pattern would read as wildcards.
| CREATE TABLE `.innerXid.events` (
|   id UInt64
| )
| ENGINE = MergeTree
| ORDER BY id;
< ok 0 -
> exec
|
|
| CREATE TABLE `line
| break` (
|   id UInt64
| )
| ENGINE = MergeTree
| ORDER BY id;
< ok 0 -
> exec
|
|
| CREATE DICTIONARY dict (
|   id UInt64,
|   name String
| )
| PRIMARY KEY id
| SOURCE(CLICKHOUSE(TABLE 'src'))
| LAYOUT(FLAT())
| LIFETIME(300);
< ok 0 -
> exec
|
|
| -- A dictionary whose source is a view, and a table that loads it: creating the
| -- table queries the view, so both follow the views.
| CREATE VIEW v_source AS SELECT id FROM src;
< ok 0 -
> exec
|
|
| CREATE DICTIONARY view_dict (
|   id UInt64
| )
| PRIMARY KEY id
| SOURCE(CLICKHOUSE(TABLE 'v_source'))
| LAYOUT(FLAT())
| LIFETIME(300);
< ok 0 -
> exec
|
|
| CREATE TABLE view_dict_table (
|   id UInt64
| )
| ENGINE = Dictionary('omni_testdump.view_dict');
< ok 0 -
> exec
|
|
| CREATE TABLE dict_table (
|   id UInt64
| )
| ENGINE = Dictionary('omni_testdump.dict');
< ok 0 -
> exec
|
|
| -- A dictionary whose source is a materialized view, and a table that loads it
| -- under a name that sorts first: creating the table queries the view, so it
| -- follows the view.
| -- Its table alias is a reference with the spaces around it, in backquotes: a name
| -- inside backquotes is part of that name, not a reference to a_mv_dict_table,
| -- which loads the dictionary this view sources.
| CREATE MATERIALIZED VIEW z_src_mv TO target AS SELECT ` omni_testdump.a_mv_dict_table `.id, 'm' AS name FROM src AS ` omni_testdump.a_mv_dict_table `;
< ok 0 -
> exec
|
|
| CREATE DICTIONARY mv_dict (
|   id UInt64
| )
| PRIMARY KEY id
| SOURCE(CLICKHOUSE(TABLE 'z_src_mv'))
| LAYOUT(FLAT())
| LIFETIME(300);
< ok 0 -
> exec
|
|
| CREATE TABLE a_mv_dict_table (
|   id UInt64
| )
| ENGINE = Dictionary('omni_testdump.mv_dict');
< ok 0 -
> exec
|
|
| -- A dictionary whose source is a query, and a table that loads it under a name
| -- that sorts first: the query reads z_src_mv, which the server parses like any
| -- other statement.
| CREATE DICTIONARY query_dict (
|   id UInt64
| )
| PRIMARY KEY id
| SOURCE(CLICKHOUSE(QUERY 'SELECT id FROM omni_testdump.z_src_mv'))
| LAYOUT(FLAT())
| LIFETIME(300);
< ok 0 -
> exec
|
|
| CREATE TABLE a_query_dict_table (
|   id UInt64
| )
| ENGINE = Dictionary('omni_testdump.query_dict');
< ok 0 -
> exec
|
|
| CREATE VIEW v_names AS SELECT id, name FROM src;
< ok 0 -
> exec
|
|
| -- A view states its columns, which ClickHouse creates without resolving its
| -- query, so this one replays before the view it reads.
| CREATE VIEW a_first AS SELECT id FROM v_names;
< ok 0 -
> exec
|
|
| CREATE MATERIALIZED VIEW mv_base TO target AS SELECT id, 'omni_testdump.a_top' AS name FROM src;
< ok 0 -
> exec
|
|
|
| -- An ordinary view reads mv_base, and a materialized view reads that view: a
| -- materialized view resolves its source when it is created, through the view, so
| -- the dump has to write mv_base first whatever the names say. The literal above
| -- names a_top and is not a read.
| CREATE VIEW v_bridge AS SELECT id, name FROM mv_base;
< ok 0 -
> exec
|
|
| CREATE MATERIALIZED VIEW a_top TO target AS SELECT id, name FROM v_bridge;
< ok 0 -
> exec
|
|
| -- A materialized view whose name holds a line break, and one that reads it under
| -- a name that sorts first: the reference the dump looks for is the rendering
| -- ClickHouse writes, with the line break escaped.
| CREATE MATERIALIZED VIEW `nl
| mv` TO target AS SELECT id, 'x' AS name FROM src;
< ok 0 -
> exec
|
|
| CREATE MATERIALIZED VIEW a_after TO target AS SELECT id, 'y' AS name FROM `nl
| mv`;
< ok 0 -
> exec
|
|
| -- A materialized view reading a table on the Dictionary engine, under a name that
| -- sorts first: creating it needs the dictionary table, which loads the
| -- dictionary.
| CREATE MATERIALIZED VIEW a_dict TO target AS SELECT id, 'z' AS name FROM dict_table;
< ok 0 -
> exec
|
|
| -- Row policies: one a role reads through, a restrictive one whose name needs
| -- backquotes and that applies to everyone but that role, and one on every table
| -- of the database. A row policy is an access entity that outlives the database,
| -- so the dump states each with OR REPLACE.
| CREATE ROLE omni_dump_reader;
< ok 0 -
> exec
|
|
| CREATE ROW POLICY tenant ON src FOR SELECT USING name != '' TO omni_dump_reader;
< ok 0 -
> exec
|
|
| CREATE ROW POLICY `odd policy` ON src AS RESTRICTIVE FOR SELECT USING id > 0 TO ALL EXCEPT omni_dump_reader;
< ok 0 -
> exec
|
|
| CREATE ROW POLICY everything ON omni_testdump.* USING 1 TO omni_dump_reader;
< ok 0 -
> exec
|
|
| CREATE MATERIALIZED VIEW mv_own ENGINE = MergeTree ORDER BY id AS SELECT id FROM src; -- not restorable: TestDumpGaps
< ok 0 -
> exec
| SET format_display_secrets_in_show_and_select = 1
< ok 0 -
> query
|
| 		SELECT
| 			name
| 		FROM system.tables
| 		WHERE database = ? AND NOT match(name, ?)
| 		ORDER BY name
: string
| omni_testdump
: string
| ^[.]inner_id[.]([a-z_]+[.])?[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$
< result
< column "name" "String" "string"
< row
: string
| .events
< row
: string
| .innerXid.events
< row
: string
| a_after
< row
: string
| a_dict
< row
: string
| a_first
< row
: string
| a_mv_dict_table
< row
: string
| a_query_dict_table
< row
: string
| a_top
< row
: string
| dict
< row
: string
| dict_table
< row
: string
| line
| break
< row
: string
| mv_base
< row
: string
| mv_dict
< row
: string
| mv_own
< row
: string
| nl
| mv
< row
: string
| odd-name
< row
: string
| query_dict
< row
: string
| src
< row
: string
| target
< row
: string
| v_bridge
< row
: string
| v_names
< row
: string
| v_source
< row
: string
| view_dict
< row
: string
| view_dict_table
< row
: string
| z_src_mv
> query
| SHOW CREATE TABLE `omni_testdump`.`.events`
< result
< column "statement" "String" "string"
< row
: string
| CREATE TABLE omni_testdump.`.events`
| (
|     `id` UInt64
| )
| ENGINE = MergeTree
| ORDER BY id
| SETTINGS index_granularity = 8192
> query
| SHOW CREATE TABLE `omni_testdump`.`.innerXid.events`
< result
< column "statement" "String" "string"
< row
: string
| CREATE TABLE omni_testdump.`.innerXid.events`
| (
|     `id` UInt64
| )
| ENGINE = MergeTree
| ORDER BY id
| SETTINGS index_granularity = 8192
> query
| SHOW CREATE TABLE `omni_testdump`.`a_after`
< result
< column "statement" "String" "string"
< row
: string
| CREATE MATERIALIZED VIEW omni_testdump.a_after TO omni_testdump.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     'y' AS name
| FROM omni_testdump.`nl\nmv`
> query
| SHOW CREATE TABLE `omni_testdump`.`a_dict`
< result
< column "statement" "String" "string"
< row
: string
| CREATE MATERIALIZED VIEW omni_testdump.a_dict TO omni_testdump.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     'z' AS name
| FROM omni_testdump.dict_table
> query
| SHOW CREATE TABLE `omni_testdump`.`a_first`
< result
< column "statement" "String" "string"
< row
: string
| CREATE VIEW omni_testdump.a_first
| (
|     `id` UInt64
| )
| AS SELECT id
| FROM omni_testdump.v_names
> query
| SHOW CREATE TABLE `omni_testdump`.`a_mv_dict_table`
< result
< column "statement" "String" "string"
< row
: string
| CREATE TABLE omni_testdump.a_mv_dict_table
| (
|     `id` UInt64
| )
| ENGINE = Dictionary('omni_testdump.mv_dict')
> query
| SHOW CREATE TABLE `omni_testdump`.`a_query_dict_table`
< result
< column "statement" "String" "string"
< row
: string
| CREATE TABLE omni_testdump.a_query_dict_table
| (
|     `id` UInt64
| )
| ENGINE = Dictionary('omni_testdump.query_dict')
> query
| SHOW CREATE TABLE `omni_testdump`.`a_top`
< result
< column "statement" "String" "string"
< row
: string
| CREATE MATERIALIZED VIEW omni_testdump.a_top TO omni_testdump.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     name
| FROM omni_testdump.v_bridge
> query
| SHOW CREATE TABLE `omni_testdump`.`dict`
< result
< column "statement" "String" "string"
< row
: string
| CREATE DICTIONARY omni_testdump.dict
| (
|     `id` UInt64,
|     `name` String
| )
| PRIMARY KEY id
| SOURCE(CLICKHOUSE(TABLE 'src'))
| LIFETIME(MIN 0 MAX 300)
| LAYOUT(FLAT())
> query
| SHOW CREATE TABLE `omni_testdump`.`dict_table`
< result
< column "statement" "String" "string"
< row
: string
| CREATE TABLE omni_testdump.dict_table
| (
|     `id` UInt64
| )
| ENGINE = Dictionary('omni_testdump.dict')
> query
| SHOW CREATE TABLE `omni_testdump`.`line\nbreak`
< result
< column "statement" "String" "string"
< row
: string
| CREATE TABLE omni_testdump.`line\nbreak`
| (
|     `id` UInt64
| )
| ENGINE = MergeTree
| ORDER BY id
| SETTINGS index_granularity = 8192
> query
| SHOW CREATE TABLE `omni_testdump`.`mv_base`
< result
< column "statement" "String" "string"
< row
: string
| CREATE MATERIALIZED VIEW omni_testdump.mv_base TO omni_testdump.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     'omni_testdump.a_top' AS name
| FROM omni_testdump.src
> query
| SHOW CREATE TABLE `omni_testdump`.`mv_dict`
< result
< column "statement" "String" "string"
< row
: string
| CREATE DICTIONARY omni_testdump.mv_dict
| (
|     `id` UInt64
| )
| PRIMARY KEY id
| SOURCE(CLICKHOUSE(TABLE 'z_src_mv'))
| LIFETIME(MIN 0 MAX 300)
| LAYOUT(FLAT())
> query
| SHOW CREATE TABLE `omni_testdump`.`mv_own`
< result
< column "statement" "String" "string"
< row
: string
| CREATE MATERIALIZED VIEW omni_testdump.mv_own
| (
|     `id` UInt64
| )
| ENGINE = MergeTree
| ORDER BY id
| SETTINGS index_granularity = 8192
| AS SELECT id
| FROM omni_testdump.src
> query
| SHOW CREATE TABLE `omni_testdump`.`nl\nmv`
< result
< column "statement" "String" "string"
< row
: string
| CREATE MATERIALIZED VIEW omni_testdump.`nl\nmv` TO omni_testdump.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     'x' AS name
| FROM omni_testdump.src
> query
| SHOW CREATE TABLE `omni_testdump`.`odd-name`
< result
< column "statement" "String" "string"
< row
: string
| CREATE TABLE omni_testdump.`odd-name`
| (
|     `id` UInt64
| )
| ENGINE = MergeTree
| ORDER BY id
| SETTINGS index_granularity = 8192
> query
| SHOW CREATE TABLE `omni_testdump`.`query_dict`
< result
< column "statement" "String" "string"
< row
: string
| CREATE DICTIONARY omni_testdump.query_dict
| (
|     `id` UInt64
| )
| PRIMARY KEY id
| SOURCE(CLICKHOUSE(QUERY 'SELECT id FROM omni_testdump.z_src_mv'))
| LIFETIME(MIN 0 MAX 300)
| LAYOUT(FLAT())
> query
| SHOW CREATE TABLE `omni_testdump`.`src`
< result
< column "statement" "String" "string"
< row
: string
| CREATE TABLE omni_testdump.src
| (
|     `id` UInt64,
|     `name` String,
|     INDEX name_bloom name TYPE bloom_filter GRANULARITY 4
| )
| ENGINE = MergeTree
| PRIMARY KEY id
| ORDER BY (id, name)
| SETTINGS index_granularity = 4096
| COMMENT 'the source'
> query
| SHOW CREATE TABLE `omni_testdump`.`target`
< result
< column "statement" "String" "string"
< row
: string
| CREATE TABLE omni_testdump.target
| (
|     `id` UInt64,
|     `name` String
| )
| ENGINE = MergeTree
| ORDER BY id
| SETTINGS index_granularity = 8192
> query
| SHOW CREATE TABLE `omni_testdump`.`v_bridge`
< result
< column "statement" "String" "string"
< row
: string
| CREATE VIEW omni_testdump.v_bridge
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     name
| FROM omni_testdump.mv_base
> query
| SHOW CREATE TABLE `omni_testdump`.`v_names`
< result
< column "statement" "String" "string"
< row
: string
| CREATE VIEW omni_testdump.v_names
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     name
| FROM omni_testdump.src
> query
| SHOW CREATE TABLE `omni_testdump`.`v_source`
< result
< column "statement" "String" "string"
< row
: string
| CREATE VIEW omni_testdump.v_source
| (
|     `id` UInt64
| )
| AS SELECT id
| FROM omni_testdump.src
> query
| SHOW CREATE TABLE `omni_testdump`.`view_dict`
< result
< column "statement" "String" "string"
< row
: string
| CREATE DICTIONARY omni_testdump.view_dict
| (
|     `id` UInt64
| )
| PRIMARY KEY id
| SOURCE(CLICKHOUSE(TABLE 'v_source'))
| LIFETIME(MIN 0 MAX 300)
| LAYOUT(FLAT())
> query
| SHOW CREATE TABLE `omni_testdump`.`view_dict_table`
< result
< column "statement" "String" "string"
< row
: string
| CREATE TABLE omni_testdump.view_dict_table
| (
|     `id` UInt64
| )
| ENGINE = Dictionary('omni_testdump.view_dict')
> query
| SHOW CREATE TABLE `omni_testdump`.`z_src_mv`
< result
< column "statement" "String" "string"
< row
: string
| CREATE MATERIALIZED VIEW omni_testdump.z_src_mv TO omni_testdump.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     ` omni_testdump.a_mv_dict_table `.id,
|     'm' AS name
| FROM omni_testdump.src AS ` omni_testdump.a_mv_dict_table `
> query
| EXPLAIN AST CREATE MATERIALIZED VIEW omni_testdump.a_after TO omni_testdump.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     'y' AS name
| FROM omni_testdump.`nl\nmv`
< result
< column "explain" "String" "string"
< row
: string
| CreateQuery omni_testdump a_after (children 5)
< row
: string
|  Identifier omni_testdump
< row
: string
|  Identifier a_after
< row
: string
|  Columns definition (children 1)
< row
: string
|   ExpressionList (children 2)
< row
: string
|    ColumnDeclaration id (children 1)
< row
: string
|     DataType UInt64
< row
: string
|    ColumnDeclaration name (children 1)
< row
: string
|     DataType String
< row
: string
|  SelectWithUnionQuery (children 1)
< row
: string
|   ExpressionList (children 1)
< row
: string
|    SelectQuery (children 2)
< row
: string
|     ExpressionList (children 2)
< row
: string
|      Identifier id
< row
: string
|      Literal 'y' (alias name)
< row
: string
|     TablesInSelectQuery (children 1)
< row
: string
|      TablesInSelectQueryElement (children 1)
< row
: string
|       TableExpression (children 1)
< row
: string
|        TableIdentifier omni_testdump.nl
< row
: string
| mv
< row
: string
|  ViewTargets
> query
| EXPLAIN AST CREATE MATERIALIZED VIEW omni_testdump.a_dict TO omni_testdump.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     'z' AS name
| FROM omni_testdump.dict_table
< result
< column "explain" "String" "string"
< row
: string
| CreateQuery omni_testdump a_dict (children 5)
< row
: string
|  Identifier omni_testdump
< row
: string
|  Identifier a_dict
< row
: string
|  Columns definition (children 1)
< row
: string
|   ExpressionList (children 2)
< row
: string
|    ColumnDeclaration id (children 1)
< row
: string
|     DataType UInt64
< row
: string
|    ColumnDeclaration name (children 1)
< row
: string
|     DataType String
< row
: string
|  SelectWithUnionQuery (children 1)
< row
: string
|   ExpressionList (children 1)
< row
: string
|    SelectQuery (children 2)
< row
: string
|     ExpressionList (children 2)
< row
: string
|      Identifier id
< row
: string
|      Literal 'z' (alias name)
< row
: string
|     TablesInSelectQuery (children 1)
< row
: string
|      TablesInSelectQueryElement (children 1)
< row
: string
|       TableExpression (children 1)
< row
: string
|        TableIdentifier omni_testdump.dict_table
< row
: string
|  ViewTargets
> query
| EXPLAIN AST CREATE VIEW omni_testdump.a_first
| (
|     `id` UInt64
| )
| AS SELECT id
| FROM omni_testdump.v_names
< result
< column "explain" "String" "string"
< row
: string
| CreateQuery omni_testdump a_first (children 4)
< row
: string
|  Identifier omni_testdump
< row
: string
|  Identifier a_first
< row
: string
|  Columns definition (children 1)
< row
: string
|   ExpressionList (children 1)
< row
: string
|    ColumnDeclaration id (children 1)
< row
: string
|     DataType UInt64
< row
: string
|  SelectWithUnionQuery (children 1)
< row
: string
|   ExpressionList (children 1)
< row
: string
|    SelectQuery (children 2)
< row
: string
|     ExpressionList (children 1)
< row
: string
|      Identifier id
< row
: string
|     TablesInSelectQuery (children 1)
< row
: string
|      TablesInSelectQueryElement (children 1)
< row
: string
|       TableExpression (children 1)
< row
: string
|        TableIdentifier omni_testdump.v_names
> query
| EXPLAIN AST SELECT id FROM omni_testdump.z_src_mv
< result
< column "explain" "String" "string"
< row
: string
| SelectWithUnionQuery (children 1)
< row
: string
|  ExpressionList (children 1)
< row
: string
|   SelectQuery (children 2)
< row
: string
|    ExpressionList (children 1)
< row
: string
|     Identifier id
< row
: string
|    TablesInSelectQuery (children 1)
< row
: string
|     TablesInSelectQueryElement (children 1)
< row
: string
|      TableExpression (children 1)
< row
: string
|       TableIdentifier omni_testdump.z_src_mv
> query
| EXPLAIN AST CREATE MATERIALIZED VIEW omni_testdump.a_top TO omni_testdump.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     name
| FROM omni_testdump.v_bridge
< result
< column "explain" "String" "string"
< row
: string
| CreateQuery omni_testdump a_top (children 5)
< row
: string
|  Identifier omni_testdump
< row
: string
|  Identifier a_top
< row
: string
|  Columns definition (children 1)
< row
: string
|   ExpressionList (children 2)
< row
: string
|    ColumnDeclaration id (children 1)
< row
: string
|     DataType UInt64
< row
: string
|    ColumnDeclaration name (children 1)
< row
: string
|     DataType String
< row
: string
|  SelectWithUnionQuery (children 1)
< row
: string
|   ExpressionList (children 1)
< row
: string
|    SelectQuery (children 2)
< row
: string
|     ExpressionList (children 2)
< row
: string
|      Identifier id
< row
: string
|      Identifier name
< row
: string
|     TablesInSelectQuery (children 1)
< row
: string
|      TablesInSelectQueryElement (children 1)
< row
: string
|       TableExpression (children 1)
< row
: string
|        TableIdentifier omni_testdump.v_bridge
< row
: string
|  ViewTargets
> query
| EXPLAIN AST CREATE MATERIALIZED VIEW omni_testdump.mv_base TO omni_testdump.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     'omni_testdump.a_top' AS name
| FROM omni_testdump.src
< result
< column "explain" "String" "string"
< row
: string
| CreateQuery omni_testdump mv_base (children 5)
< row
: string
|  Identifier omni_testdump
< row
: string
|  Identifier mv_base
< row
: string
|  Columns definition (children 1)
< row
: string
|   ExpressionList (children 2)
< row
: string
|    ColumnDeclaration id (children 1)
< row
: string
|     DataType UInt64
< row
: string
|    ColumnDeclaration name (children 1)
< row
: string
|     DataType String
< row
: string
|  SelectWithUnionQuery (children 1)
< row
: string
|   ExpressionList (children 1)
< row
: string
|    SelectQuery (children 2)
< row
: string
|     ExpressionList (children 2)
< row
: string
|      Identifier id
< row
: string
|      Literal 'omni_testdump.a_top' (alias name)
< row
: string
|     TablesInSelectQuery (children 1)
< row
: string
|      TablesInSelectQueryElement (children 1)
< row
: string
|       TableExpression (children 1)
< row
: string
|        TableIdentifier omni_testdump.src
< row
: string
|  ViewTargets
> query
| EXPLAIN AST CREATE MATERIALIZED VIEW omni_testdump.mv_own
| (
|     `id` UInt64
| )
| ENGINE = MergeTree
| ORDER BY id
| SETTINGS index_granularity = 8192
| AS SELECT id
| FROM omni_testdump.src
< result
< column "explain" "String" "string"
< row
: string
| CreateQuery omni_testdump mv_own (children 5)
< row
: string
|  Identifier omni_testdump
< row
: string
|  Identifier mv_own
< row
: string
|  Columns definition (children 1)
< row
: string
|   ExpressionList (children 1)
< row
: string
|    ColumnDeclaration id (children 1)
< row
: string
|     DataType UInt64
< row
: string
|  SelectWithUnionQuery (children 1)
< row
: string
|   ExpressionList (children 1)
< row
: string
|    SelectQuery (children 2)
< row
: string
|     ExpressionList (children 1)
< row
: string
|      Identifier id
< row
: string
|     TablesInSelectQuery (children 1)
< row
: string
|      TablesInSelectQueryElement (children 1)
< row
: string
|       TableExpression (children 1)
< row
: string
|        TableIdentifier omni_testdump.src
< row
: string
|  ViewTargets (children 1)
< row
: string
|   Storage definition (children 3)
< row
: string
|    Function MergeTree
< row
: string
|    Identifier id
< row
: string
|    Set
> query
| EXPLAIN AST CREATE MATERIALIZED VIEW omni_testdump.`nl\nmv` TO omni_testdump.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     'x' AS name
| FROM omni_testdump.src
< result
< column "explain" "String" "string"
< row
: string
| CreateQuery omni_testdump nl
< row
: string
| mv (children 5)
< row
: string
|  Identifier omni_testdump
< row
: string
|  Identifier nl
< row
: string
| mv
< row
: string
|  Columns definition (children 1)
< row
: string
|   ExpressionList (children 2)
< row
: string
|    ColumnDeclaration id (children 1)
< row
: string
|     DataType UInt64
< row
: string
|    ColumnDeclaration name (children 1)
< row
: string
|     DataType String
< row
: string
|  SelectWithUnionQuery (children 1)
< row
: string
|   ExpressionList (children 1)
< row
: string
|    SelectQuery (children 2)
< row
: string
|     ExpressionList (children 2)
< row
: string
|      Identifier id
< row
: string
|      Literal 'x' (alias name)
< row
: string
|     TablesInSelectQuery (children 1)
< row
: string
|      TablesInSelectQueryElement (children 1)
< row
: string
|       TableExpression (children 1)
< row
: string
|        TableIdentifier omni_testdump.src
< row
: string
|  ViewTargets
> query
| EXPLAIN AST CREATE VIEW omni_testdump.v_bridge
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     name
| FROM omni_testdump.mv_base
< result
< column "explain" "String" "string"
< row
: string
| CreateQuery omni_testdump v_bridge (children 4)
< row
: string
|  Identifier omni_testdump
< row
: string
|  Identifier v_bridge
< row
: string
|  Columns definition (children 1)
< row
: string
|   ExpressionList (children 2)
< row
: string
|    ColumnDeclaration id (children 1)
< row
: string
|     DataType UInt64
< row
: string
|    ColumnDeclaration name (children 1)
< row
: string
|     DataType String
< row
: string
|  SelectWithUnionQuery (children 1)
< row
: string
|   ExpressionList (children 1)
< row
: string
|    SelectQuery (children 2)
< row
: string
|     ExpressionList (children 2)
< row
: string
|      Identifier id
< row
: string
|      Identifier name
< row
: string
|     TablesInSelectQuery (children 1)
< row
: string
|      TablesInSelectQueryElement (children 1)
< row
: string
|       TableExpression (children 1)
< row
: string
|        TableIdentifier omni_testdump.mv_base
> query
| EXPLAIN AST CREATE VIEW omni_testdump.v_names
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     id,
|     name
| FROM omni_testdump.src
< result
< column "explain" "String" "string"
< row
: string
| CreateQuery omni_testdump v_names (children 4)
< row
: string
|  Identifier omni_testdump
< row
: string
|  Identifier v_names
< row
: string
|  Columns definition (children 1)
< row
: string
|   ExpressionList (children 2)
< row
: string
|    ColumnDeclaration id (children 1)
< row
: string
|     DataType UInt64
< row
: string
|    ColumnDeclaration name (children 1)
< row
: string
|     DataType String
< row
: string
|  SelectWithUnionQuery (children 1)
< row
: string
|   ExpressionList (children 1)
< row
: string
|    SelectQuery (children 2)
< row
: string
|     ExpressionList (children 2)
< row
: string
|      Identifier id
< row
: string
|      Identifier name
< row
: string
|     TablesInSelectQuery (children 1)
< row
: string
|      TablesInSelectQueryElement (children 1)
< row
: string
|       TableExpression (children 1)
< row
: string
|        TableIdentifier omni_testdump.src
> query
| EXPLAIN AST CREATE VIEW omni_testdump.v_source
| (
|     `id` UInt64
| )
| AS SELECT id
| FROM omni_testdump.src
< result
< column "explain" "String" "string"
< row
: string
| CreateQuery omni_testdump v_source (children 4)
< row
: string
|  Identifier omni_testdump
< row
: string
|  Identifier v_source
< row
: string
|  Columns definition (children 1)
< row
: string
|   ExpressionList (children 1)
< row
: string
|    ColumnDeclaration id (children 1)
< row
: string
|     DataType UInt64
< row
: string
|  SelectWithUnionQuery (children 1)
< row
: string
|   ExpressionList (children 1)
< row
: string
|    SelectQuery (children 2)
< row
: string
|     ExpressionList (children 1)
< row
: string
|      Identifier id
< row
: string
|     TablesInSelectQuery (children 1)
< row
: string
|      TablesInSelectQueryElement (children 1)
< row
: string
|       TableExpression (children 1)
< row
: string
|        TableIdentifier omni_testdump.src
> query
| EXPLAIN AST CREATE MATERIALIZED VIEW omni_testdump.z_src_mv TO omni_testdump.target
| (
|     `id` UInt64,
|     `name` String
| )
| AS SELECT
|     ` omni_testdump.a_mv_dict_table `.id,
|     'm' AS name
| FROM omni_testdump.src AS ` omni_testdump.a_mv_dict_table `
< result
< column "explain" "String" "string"
< row
: string
| CreateQuery omni_testdump z_src_mv (children 5)
< row
: string
|  Identifier omni_testdump
< row
: string
|  Identifier z_src_mv
< row
: string
|  Columns definition (children 1)
< row
: string
|   ExpressionList (children 2)
< row
: string
|    ColumnDeclaration id (children 1)
< row
: string
|     DataType UInt64
< row
: string
|    ColumnDeclaration name (children 1)
< row
: string
|     DataType String
< row
: string
|  SelectWithUnionQuery (children 1)
< row
: string
|   ExpressionList (children 1)
< row
: string
|    SelectQuery (children 2)
< row
: string
|     ExpressionList (children 2)
< row
: string
|      Identifier  omni_testdump.a_mv_dict_table .id
< row
: string
|      Literal 'm' (alias name)
< row
: string
|     TablesInSelectQuery (children 1)
< row
: string
|      TablesInSelectQueryElement (children 1)
< row
: string
|       TableExpression (children 1)
< row
: string
|        TableIdentifier omni_testdump.src (alias  omni_testdump.a_mv_dict_table )
< row
: string
|  ViewTargets
> query
|
| 		SELECT
| 			short_name,
| 			table
| 		FROM system.row_policies
| 		WHERE database = ?
| 		ORDER BY table, short_name
: string
| omni_testdump
< result
< column "short_name" "String" "string"
< column "table" "String" "string"
< row
: string
| everything
: string
|
< row
: string
| odd policy
: string
| src
< row
: string
| tenant
: string
| src
> query
| SHOW CREATE ROW POLICY `everything` ON `omni_testdump`.*
< result
< column "CREATE ROW POLICY everything ON omni_testdump.*" "String" "string"
< row
: string
| CREATE ROW POLICY everything ON omni_testdump.* FOR SELECT USING 1 TO omni_dump_reader
> query
| SHOW CREATE ROW POLICY `odd policy` ON `omni_testdump`.`src`
< result
< column "CREATE ROW POLICY `odd policy` ON omni_testdump.src" "String" "string"
< row
: string
| CREATE ROW POLICY `odd policy` ON omni_testdump.src AS RESTRICTIVE FOR SELECT USING id > 0 TO ALL EXCEPT omni_dump_reader
> query
| SHOW CREATE ROW POLICY `tenant` ON `omni_testdump`.`src`
< result
< column "CREATE ROW POLICY tenant ON omni_testdump.src" "String" "string"
< row
: string
| CREATE ROW POLICY tenant ON omni_testdump.src FOR SELECT USING name != '' TO omni_dump_reader
