-- pin: default gcr.io/cloud-spanner-emulator/emulator:1.5.54 sha256:caf1bd24c081e005837b5977bae5a250e25cb4da9f25ec1abc91936ad67e4de2 v1.5.54

== TestExecuteCreateDatabasePostgres default#1
> raw
| {"op":"createDatabase","instancePath":"projects/omni/instances/pgcreatedb","createStatement":"CREATE DATABASE \"pgcreated\";","extraStatements":["CREATE TABLE written (id bigint NOT NULL, PRIMARY KEY (id))"],"dialect":"POSTGRESQL"}
< raw
| {}
> raw
| {"op":"close"}
< raw
| {}

== TestExecuteCreateDatabasePostgres default#2
> raw
| {"op":"getDatabase","database":"projects/omni/instances/pgcreatedb/databases/pgcreated"}
< raw
| {"dialect":"POSTGRESQL"}
> raw
| {"op":"query","statement":"\n    SELECT \n      TABLE_SCHEMA,\n      TABLE_NAME,\n      COLUMN_NAME,\n      ORDINAL_POSITION,\n      COLUMN_DEFAULT,\n      IS_NULLABLE = 'YES',\n      SPANNER_TYPE\n    FROM INFORMATION_SCHEMA.COLUMNS\n    WHERE LOWER(TABLE_SCHEMA) NOT IN ('information_schema', 'spanner_sys')\n    ORDER BY TABLE_SCHEMA, TABLE_NAME, ORDINAL_POSITION\n  "}
< raw
| {"columns":[{"name":"table_schema","type":{"code":"STRING"}},{"name":"table_name","type":{"code":"STRING"}},{"name":"column_name","type":{"code":"STRING"}},{"name":"ordinal_position","type":{"code":"INT64"}},{"name":"column_default","type":{"code":"STRING"}},{"name":"?column?","type":{"code":"BOOL"}},{"name":"spanner_type","type":{"code":"STRING"}}],"rows":[[{"type":{"code":"STRING"},"value":"public"},{"type":{"code":"STRING"},"value":"written"},{"type":{"code":"STRING"},"value":"id"},{"type":{"code":"INT64"},"value":"1"},{"type":{"code":"STRING"},"value":null},{"type":{"code":"BOOL"},"value":false},{"type":{"code":"STRING"},"value":"bigint"}]]}
> raw
| {"op":"query","statement":"\n    SELECT\n      TABLE_SCHEMA,\n      TABLE_NAME,\n      INDEX_NAME,\n      IS_UNIQUE = 'YES' AS IS_UNIQUE,\n      INDEX_TYPE = 'PRIMARY_KEY' AS IS_PRIMARY,\n      ARRAY (\n        SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.INDEX_COLUMNS AS ic\n        WHERE ic.TABLE_SCHEMA = i.TABLE_SCHEMA AND ic.TABLE_NAME = i.TABLE_NAME AND ic.INDEX_NAME = i.INDEX_NAME\n        ORDER BY ic.ORDINAL_POSITION\n      )\n    FROM INFORMATION_SCHEMA.INDEXES as i\n    WHERE LOWER(TABLE_SCHEMA) NOT IN ('information_schema', 'spanner_sys') AND SPANNER_IS_MANAGED = 'NO'\n    ORDER BY TABLE_SCHEMA, TABLE_NAME, INDEX_NAME\n  "}
< raw
| {"columns":[{"name":"table_schema","type":{"code":"STRING"}},{"name":"table_name","type":{"code":"STRING"}},{"name":"index_name","type":{"code":"STRING"}},{"name":"is_unique","type":{"code":"BOOL"}},{"name":"is_primary","type":{"code":"BOOL"}},{"name":"array","type":{"code":"ARRAY","arrayElementType":{"code":"STRING"}}}],"rows":[[{"type":{"code":"STRING"},"value":"public"},{"type":{"code":"STRING"},"value":"written"},{"type":{"code":"STRING"},"value":"PRIMARY_KEY"},{"type":{"code":"BOOL"},"value":true},{"type":{"code":"BOOL"},"value":true},{"type":{"code":"ARRAY","arrayElementType":{"code":"STRING"}},"value":["id"]}]]}
> raw
| {"op":"query","statement":"\n    WITH t AS (\n      SELECT DISTINCT\n        tc.CONSTRAINT_SCHEMA,\n        tc.CONSTRAINT_NAME,\n        tc.TABLE_NAME,\n        rc.UNIQUE_CONSTRAINT_SCHEMA,\n        rc.UNIQUE_CONSTRAINT_NAME,\n        ccu.TABLE_SCHEMA AS REFERENCED_TABLE_SCHEMA,\n        ccu.TABLE_NAME AS REFERENCED_TABLE_NAME,\n        rc.DELETE_RULE,\n        rc.UPDATE_RULE,\n        rc.MATCH_OPTION\n      FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS AS tc\n      JOIN INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS AS rc\n      ON tc.CONSTRAINT_SCHEMA = rc.CONSTRAINT_SCHEMA AND tc.CONSTRAINT_NAME = rc.CONSTRAINT_NAME \n      JOIN INFORMATION_SCHEMA.CONSTRAINT_COLUMN_USAGE as ccu\n      ON tc.CONSTRAINT_SCHEMA = ccu.CONSTRAINT_SCHEMA AND tc.CONSTRAINT_NAME = ccu.CONSTRAINT_NAME \n      WHERE tc.CONSTRAINT_TYPE = 'FOREIGN KEY'\n    )\n    SELECT\n      t.CONSTRAINT_SCHEMA,\n      t.CONSTRAINT_NAME,\n      t.TABLE_NAME,\n      t.REFERENCED_TABLE_SCHEMA,\n      t.REFERENCED_TABLE_NAME,\n      ARRAY (\n        SELECT\n          COLUMN_NAME\n        FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE AS kcu\n        WHERE t.CONSTRAINT_SCHEMA = kcu.CONSTRAINT_SCHEMA AND t.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME\n        ORDER BY kcu.ORDINAL_POSITION\n      ) AS TABLE_COLUMNS,\n      ARRAY (\n        SELECT\n          COLUMN_NAME\n        FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE AS kcu\n        WHERE t.UNIQUE_CONSTRAINT_SCHEMA = kcu.CONSTRAINT_SCHEMA AND t.UNIQUE_CONSTRAINT_NAME = kcu.CONSTRAINT_NAME\n        ORDER BY kcu.ORDINAL_POSITION\n      ) AS REFERENCED_TABLE_COLUMNS,\n      t.DELETE_RULE,\n      t.UPDATE_RULE,\n      t.MATCH_OPTION\n    FROM t\n  "}
< raw
| {"columns":[{"name":"constraint_schema","type":{"code":"STRING"}},{"name":"constraint_name","type":{"code":"STRING"}},{"name":"table_name","type":{"code":"STRING"}},{"name":"referenced_table_schema","type":{"code":"STRING"}},{"name":"referenced_table_name","type":{"code":"STRING"}},{"name":"table_columns","type":{"code":"ARRAY","arrayElementType":{"code":"STRING"}}},{"name":"referenced_table_columns","type":{"code":"ARRAY","arrayElementType":{"code":"STRING"}}},{"name":"delete_rule","type":{"code":"STRING"}},{"name":"update_rule","type":{"code":"STRING"}},{"name":"match_option","type":{"code":"STRING"}}]}
> raw
| {"op":"query","statement":"\n    SELECT\n      TABLE_SCHEMA,\n      TABLE_NAME\n    FROM INFORMATION_SCHEMA.TABLES\n    WHERE LOWER(TABLE_SCHEMA) NOT IN ('information_schema', 'spanner_sys') AND TABLE_TYPE = 'BASE TABLE'\n    ORDER BY TABLE_SCHEMA, TABLE_NAME\n  "}
< raw
| {"columns":[{"name":"table_schema","type":{"code":"STRING"}},{"name":"table_name","type":{"code":"STRING"}}],"rows":[[{"type":{"code":"STRING"},"value":"public"},{"type":{"code":"STRING"},"value":"written"}]]}
> raw
| {"op":"query","statement":"\n  SELECT\n    TABLE_SCHEMA,\n    TABLE_NAME,\n    VIEW_DEFINITION\n  FROM INFORMATION_SCHEMA.VIEWS\n  WHERE LOWER(TABLE_SCHEMA) NOT IN ('information_schema', 'spanner_sys')\n  "}
< raw
| {"columns":[{"name":"table_schema","type":{"code":"STRING"}},{"name":"table_name","type":{"code":"STRING"}},{"name":"view_definition","type":{"code":"STRING"}}]}
> raw
| {"op":"close"}
< raw
| {}
