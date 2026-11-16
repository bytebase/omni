-- pin: default amazon/dynamodb-local:3.3.2 sha256:6ada631ec7aa5fc83446855eb17233e3dad184fec47477e8140c6c569cfe0c86 3.3.2

== TestQueryErrorTextIsTheClients default#1
> raw
| {"op":"executeStatement","statement":"SELECT * FROM \"nope\""}
< raw
| {"error":{"text":"operation error DynamoDB: ExecuteStatement, https response error StatusCode: 400, RequestID: 00000000-0000-0000-0000-000000000001, ResourceNotFoundException: Cannot do operations on a non-existent table"}}
