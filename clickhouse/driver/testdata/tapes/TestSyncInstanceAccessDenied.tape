-- pin: default clickhouse/clickhouse-server:26.8.38.9 sha256:5d30ac138948b43fd928842be7da1dcd2de712d32ba6c42368ff5009ba2072ff 26.8.38.9

== TestSyncInstanceAccessDenied default#1
> exec
| CREATE ROLE omni_hidden;
< ok 0 -
> exec
|
| CREATE USER omni_users_only IDENTIFIED WITH no_password;
< ok 0 -
> exec
|
| GRANT SELECT ON system.users TO omni_users_only;
< ok 0 -
> exec
| EXECUTE AS omni_users_only
< ok 0 -
> query
| SELECT currentUser()
< result
< column "currentUser()" "String" "string"
< row
: string
| omni_users_only
> query
|
| 		SELECT
| 			name,
| 			if(empty(valid_until) OR arrayExists(t -> toUnixTimestamp(t) = 0, valid_until), '', toString(arrayMax(valid_until)))
| 		FROM system.users
| 		ORDER BY name
< result
< column "name" "String" "string"
< column "if(or(empty(valid_until), arrayExists(lambda(tuple(t), equals(toUnixTimestamp(t), 0)), valid_until)), '', toString(arrayMax(valid_until)))" "String" "string"
< row
: string
| default
: string
|
< row
: string
| omni_users_only
: string
|
> query
| SELECT name FROM system.roles ORDER BY name
< error
| {"code":497,"name":"DB::Exception","message":"omni_users_only: Not enough privileges. To execute this query, it's necessary to have the grant SELECT ON system.roles"}
