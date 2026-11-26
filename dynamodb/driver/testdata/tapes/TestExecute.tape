-- pin: driver amazon/dynamodb-local:3.3.2 sha256:6ada631ec7aa5fc83446855eb17233e3dad184fec47477e8140c6c569cfe0c86 3.3.2

== TestExecute driver#1
> raw
| {"op":"executeTransaction","statement":"DELETE FROM \"writes\" WHERE id = 'w1';"}
< raw
| {}
> raw
| {"op":"executeTransaction","statement":"\nINSERT INTO \"writes\" VALUE {'id':'w1','name':'first'};"}
< raw
| {}
> raw
| {"op":"executeStatement","statement":"SELECT * FROM \"writes\" WHERE id = 'w1'"}
< raw
| {"items":[{"id":{"S":"w1"},"name":{"S":"first"}}]}
> raw
| {"op":"executeTransaction","statement":"DELETE FROM \"writes\" WHERE id = 'w1'"}
< raw
| {}
