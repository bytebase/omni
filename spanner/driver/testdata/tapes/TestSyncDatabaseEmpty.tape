-- pin: default gcr.io/cloud-spanner-emulator/emulator:1.5.54 sha256:caf1bd24c081e005837b5977bae5a250e25cb4da9f25ec1abc91936ad67e4de2 v1.5.54

== TestSyncDatabaseEmpty default#1
> raw
| {"op":"getDatabase","database":"projects/omni/instances/empty/databases/empty"}
< raw
| {"dialect":"GOOGLE_STANDARD_SQL"}
> raw
| {"op":"query","statement":"\n    SELECT \n      TABLE_SCHEMA,\n      TABLE_NAME,\n      COLUMN_NAME,\n      ORDINAL_POSITION,\n      COLUMN_DEFAULT,\n      IS_NULLABLE = 'YES',\n      SPANNER_TYPE\n    FROM INFORMATION_SCHEMA.COLUMNS\n    WHERE TABLE_SCHEMA NOT IN ('INFORMATION_SCHEMA', 'SPANNER_SYS')\n    ORDER BY TABLE_SCHEMA, TABLE_NAME, ORDINAL_POSITION\n  "}
< raw
| {"columns":[{"name":"TABLE_SCHEMA","type":{"code":"STRING"}},{"name":"TABLE_NAME","type":{"code":"STRING"}},{"name":"COLUMN_NAME","type":{"code":"STRING"}},{"name":"ORDINAL_POSITION","type":{"code":"INT64"}},{"name":"COLUMN_DEFAULT","type":{"code":"STRING"}},{"name":"","type":{"code":"BOOL"}},{"name":"SPANNER_TYPE","type":{"code":"STRING"}}]}
> raw
| {"op":"query","statement":"\n    SELECT\n      TABLE_SCHEMA,\n      TABLE_NAME,\n      INDEX_NAME,\n      IS_UNIQUE,\n      INDEX_TYPE = 'PRIMARY_KEY' AS IS_PRIMARY,\n      ARRAY (\n        SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.INDEX_COLUMNS AS ic\n        WHERE ic.TABLE_SCHEMA = i.TABLE_SCHEMA AND ic.TABLE_NAME = i.TABLE_NAME AND ic.INDEX_NAME = i.INDEX_NAME\n        ORDER BY ic.ORDINAL_POSITION\n      )\n    FROM INFORMATION_SCHEMA.INDEXES as i\n    WHERE TABLE_SCHEMA NOT IN ('INFORMATION_SCHEMA', 'SPANNER_SYS') AND SPANNER_IS_MANAGED = FALSE\n    ORDER BY TABLE_SCHEMA, TABLE_NAME, INDEX_NAME\n  "}
< raw
| {"columns":[{"name":"TABLE_SCHEMA","type":{"code":"STRING"}},{"name":"TABLE_NAME","type":{"code":"STRING"}},{"name":"INDEX_NAME","type":{"code":"STRING"}},{"name":"IS_UNIQUE","type":{"code":"BOOL"}},{"name":"IS_PRIMARY","type":{"code":"BOOL"}},{"name":"","type":{"code":"ARRAY","arrayElementType":{"code":"STRING"}}}]}
> raw
| {"op":"query","statement":"\n    WITH t AS (\n      SELECT DISTINCT\n        tc.CONSTRAINT_SCHEMA,\n        tc.CONSTRAINT_NAME,\n        tc.TABLE_NAME,\n        rc.UNIQUE_CONSTRAINT_SCHEMA,\n        rc.UNIQUE_CONSTRAINT_NAME,\n        ccu.TABLE_SCHEMA AS REFERENCED_TABLE_SCHEMA,\n        ccu.TABLE_NAME AS REFERENCED_TABLE_NAME,\n        rc.DELETE_RULE,\n        rc.UPDATE_RULE,\n        rc.MATCH_OPTION\n      FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS AS tc\n      JOIN INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS AS rc\n      ON tc.CONSTRAINT_SCHEMA = rc.CONSTRAINT_SCHEMA AND tc.CONSTRAINT_NAME = rc.CONSTRAINT_NAME \n      JOIN INFORMATION_SCHEMA.CONSTRAINT_COLUMN_USAGE as ccu\n      ON tc.CONSTRAINT_SCHEMA = ccu.CONSTRAINT_SCHEMA AND tc.CONSTRAINT_NAME = ccu.CONSTRAINT_NAME \n      WHERE tc.CONSTRAINT_TYPE = 'FOREIGN KEY'\n    )\n    SELECT\n      t.CONSTRAINT_SCHEMA,\n      t.CONSTRAINT_NAME,\n      t.TABLE_NAME,\n      t.REFERENCED_TABLE_SCHEMA,\n      t.REFERENCED_TABLE_NAME,\n      ARRAY (\n        SELECT\n          COLUMN_NAME\n        FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE AS kcu\n        WHERE t.CONSTRAINT_SCHEMA = kcu.CONSTRAINT_SCHEMA AND t.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME\n        ORDER BY kcu.ORDINAL_POSITION\n      ) AS TABLE_COLUMNS,\n      ARRAY (\n        SELECT\n          COLUMN_NAME\n        FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE AS kcu\n        WHERE t.UNIQUE_CONSTRAINT_SCHEMA = kcu.CONSTRAINT_SCHEMA AND t.UNIQUE_CONSTRAINT_NAME = kcu.CONSTRAINT_NAME\n        ORDER BY kcu.ORDINAL_POSITION\n      ) AS REFERENCED_TABLE_COLUMNS,\n      t.DELETE_RULE,\n      t.UPDATE_RULE,\n      t.MATCH_OPTION\n    FROM t\n  "}
< raw
| {"columns":[{"name":"CONSTRAINT_SCHEMA","type":{"code":"STRING"}},{"name":"CONSTRAINT_NAME","type":{"code":"STRING"}},{"name":"TABLE_NAME","type":{"code":"STRING"}},{"name":"REFERENCED_TABLE_SCHEMA","type":{"code":"STRING"}},{"name":"REFERENCED_TABLE_NAME","type":{"code":"STRING"}},{"name":"TABLE_COLUMNS","type":{"code":"ARRAY","arrayElementType":{"code":"STRING"}}},{"name":"REFERENCED_TABLE_COLUMNS","type":{"code":"ARRAY","arrayElementType":{"code":"STRING"}}},{"name":"DELETE_RULE","type":{"code":"STRING"}},{"name":"UPDATE_RULE","type":{"code":"STRING"}},{"name":"MATCH_OPTION","type":{"code":"STRING"}}]}
> raw
| {"op":"query","statement":"\n    SELECT\n      TABLE_SCHEMA,\n      TABLE_NAME\n    FROM INFORMATION_SCHEMA.TABLES\n    WHERE TABLE_SCHEMA NOT IN ('INFORMATION_SCHEMA', 'SPANNER_SYS') AND TABLE_TYPE = 'BASE TABLE'\n    ORDER BY TABLE_SCHEMA, TABLE_NAME\n  "}
< raw
| {"columns":[{"name":"TABLE_SCHEMA","type":{"code":"STRING"}},{"name":"TABLE_NAME","type":{"code":"STRING"}}]}
> raw
| {"op":"query","statement":"\n  SELECT\n    TABLE_SCHEMA,\n    TABLE_NAME,\n    VIEW_DEFINITION\n  FROM INFORMATION_SCHEMA.VIEWS\n  WHERE TABLE_SCHEMA NOT IN ('INFORMATION_SCHEMA', 'SPANNER_SYS')\n  "}
< raw
| {"columns":[{"name":"TABLE_SCHEMA","type":{"code":"STRING"}},{"name":"TABLE_NAME","type":{"code":"STRING"}},{"name":"VIEW_DEFINITION","type":{"code":"STRING"}}]}
> raw
| {"op":"close"}
< raw
| {}
