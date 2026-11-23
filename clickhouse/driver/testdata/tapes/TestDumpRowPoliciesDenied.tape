-- pin: default clickhouse/clickhouse-server:26.8.38.9 sha256:5d30ac138948b43fd928842be7da1dcd2de712d32ba6c42368ff5009ba2072ff 26.8.38.9

== TestDumpRowPoliciesDenied default#1
> exec
| CREATE TABLE t (id UInt64) ENGINE = MergeTree ORDER BY id;
< ok 0 -
> exec
|
| CREATE ROW POLICY hidden ON t USING id > 0 TO ALL;
< ok 0 -
> exec
|
| CREATE USER omni_no_policies IDENTIFIED WITH no_password;
< ok 0 -
> exec
|
| GRANT SELECT, SHOW ON `omni_testdumprowpoliciesdenied`.* TO omni_no_policies;
< ok 0 -
> exec
|
| GRANT SELECT ON system.tables TO omni_no_policies;
< ok 0 -
> exec
| EXECUTE AS omni_no_policies
< ok 0 -
> query
| SELECT currentUser()
< result
< column "currentUser()" "String" "string"
< row
: string
| omni_no_policies
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
| omni_testdumprowpoliciesdenied
: string
| ^[.]inner_id[.]([a-z_]+[.])?[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$
< result
< column "name" "String" "string"
< row
: string
| t
> query
| SHOW CREATE TABLE `omni_testdumprowpoliciesdenied`.`t`
< result
< column "statement" "String" "string"
< row
: string
| CREATE TABLE omni_testdumprowpoliciesdenied.t
| (
|     `id` UInt64
| )
| ENGINE = MergeTree
| ORDER BY id
| SETTINGS index_granularity = 8192
> query
|
| 		SELECT
| 			short_name,
| 			table
| 		FROM system.row_policies
| 		WHERE database = ?
| 		ORDER BY table, short_name
: string
| omni_testdumprowpoliciesdenied
< error
| {"code":497,"name":"DB::Exception","message":"omni_no_policies: Not enough privileges. To execute this query, it's necessary to have the grant SELECT ON system.row_policies"}
