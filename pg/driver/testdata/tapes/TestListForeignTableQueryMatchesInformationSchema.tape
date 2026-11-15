-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestListForeignTableQueryMatchesInformationSchema default#1
> exec
|
| CREATE EXTENSION postgres_fdw;
| CREATE SERVER loopback_server FOREIGN DATA WRAPPER postgres_fdw OPTIONS (host 'localhost', dbname 'postgres');
| CREATE SCHEMA fdw_schema;
| CREATE FOREIGN TABLE fdw_schema.remote_orders (
|     id integer NOT NULL,
|     note varchar(30)
| ) SERVER loopback_server OPTIONS (schema_name 'public', table_name 'orders');
| CREATE FOREIGN TABLE public.remote_items (
|     item_id bigint
| ) SERVER loopback_server OPTIONS (schema_name 'public', table_name 'items');
|
< ok 0 -
> query
| SELECT
| 	foreign_table.foreign_table_schema,
| 	foreign_table.foreign_table_name,
| 	foreign_table.foreign_server_catalog,
| 	foreign_table.foreign_server_name
| FROM information_schema.foreign_tables AS foreign_table
| ORDER BY foreign_table.foreign_table_schema, foreign_table.foreign_table_name;
< result
< column "foreign_table_schema" "NAME"
< column "foreign_table_name" "NAME"
< column "foreign_server_catalog" "NAME"
< column "foreign_server_name" "NAME"
< row
: string
| fdw_schema
: string
| remote_orders
: string
| omni_testlistforeigntablequerymatchesinformationschema
: string
| loopback_server
< row
: string
| public
: string
| remote_items
: string
| omni_testlistforeigntablequerymatchesinformationschema
: string
| loopback_server
> query
| SELECT
| 		n.nspname AS foreign_table_schema,
| 		c.relname AS foreign_table_name,
| 		current_database() AS foreign_server_catalog,
| 		s.srvname AS foreign_server_name
| 	FROM pg_catalog.pg_foreign_table t
| 		JOIN pg_catalog.pg_class c ON c.oid = t.ftrelid
| 		JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
| 		JOIN pg_catalog.pg_foreign_server s ON s.oid = t.ftserver
| 	WHERE c.relkind = 'f'
| 	ORDER BY n.nspname, c.relname;
< result
< column "foreign_table_schema" "NAME"
< column "foreign_table_name" "NAME"
< column "foreign_server_catalog" "NAME"
< column "foreign_server_name" "NAME"
< row
: string
| fdw_schema
: string
| remote_orders
: string
| omni_testlistforeigntablequerymatchesinformationschema
: string
| loopback_server
< row
: string
| public
: string
| remote_items
: string
| omni_testlistforeigntablequerymatchesinformationschema
: string
| loopback_server
