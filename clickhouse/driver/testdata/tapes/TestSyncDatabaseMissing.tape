-- pin: default clickhouse/clickhouse-server:26.8.38.9 sha256:5d30ac138948b43fd928842be7da1dcd2de712d32ba6c42368ff5009ba2072ff 26.8.38.9

== TestSyncDatabaseMissing default#1
> connect
< error
| {"code":81,"name":"DB::Exception","message":"Database omni_no_such_database does not exist"}
> connect
< error
| {"code":81,"name":"DB::Exception","message":"Database omni_no_such_database does not exist"}
