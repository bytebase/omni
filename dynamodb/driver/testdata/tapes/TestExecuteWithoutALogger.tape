-- pin: default amazon/dynamodb-local:3.3.2 sha256:6ada631ec7aa5fc83446855eb17233e3dad184fec47477e8140c6c569cfe0c86 3.3.2

== TestExecuteWithoutALogger default#1
> raw
| {"op":"executeTransaction","statement":"DELETE FROM \"writes\" WHERE id = 'w4'"}
< raw
| {}
