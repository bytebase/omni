-- pin: default clickhouse/clickhouse-server:26.8.38.9 sha256:5d30ac138948b43fd928842be7da1dcd2de712d32ba6c42368ff5009ba2072ff 26.8.38.9

== TestExecuteStopsAtTheFailingStatement default#1
> exec
|
| CREATE TABLE e (id Int32, note String) ENGINE = Memory;
< ok 0 -
> exec
| INSERT INTO e VALUES (1, 'a');
< ok 0 -
> exec
|
| INSERT INTO no_such_table VALUES (2, 'b');
< error
| {"code":60,"name":"DB::Exception","message":"Table omni_testexecutestopsatthefailingstatement.no_such_table does not exist"}
> query
| SELECT count() FROM e
< result
< column "count()" "UInt64" "uint64"
< row
: uint64 1
