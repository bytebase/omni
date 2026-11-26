-- pin: driver amazon/dynamodb-local:3.3.2 sha256:6ada631ec7aa5fc83446855eb17233e3dad184fec47477e8140c6c569cfe0c86 3.3.2

== TestExecuteEndsAtTheFirstFailure driver#1
> raw
| {"op":"executeTransaction","statement":"DELETE FROM \"writes\" WHERE id = 'w2';"}
< raw
| {}
> raw
| {"op":"executeTransaction","statement":"\nINSERT INTO \"writes\" VALUE {'id':'w2','name':'first'};"}
< raw
| {}
> raw
| {"op":"executeTransaction","statement":"\nINSERT INTO \"writes\" VALUE {'id':'w2','name':'second'};"}
< raw
| {"error":{"text":"operation error DynamoDB: ExecuteTransaction, https response error StatusCode: 400, RequestID: 00000000-0000-0000-0000-000000000001, TransactionCanceledException: Transaction cancelled, please refer cancellation reasons for specific reasons [ValidationError]"}}
> raw
| {"op":"executeTransaction","statement":"DELETE FROM \"writes\" WHERE id = 'w2'"}
< raw
| {}
