-- pin: emulator mcr.microsoft.com/cosmosdb/linux/azure-cosmos-emulator:2.14.28 sha256:9a4f69121f20c1d6502db84625c71ebb441ab61b01880a7b2d3d8031c6dac429 2.14.28

== TestQueryExplainReportsTheGatewayError emulator#1
> raw
| {"op":"Plan","database":"cosmosfixture","container":"no-such-container","statement":"SELECT * FROM c"}
< raw
| {"error":{"response":{"method":"POST","url":"https://127.0.0.1:8081/dbs/cosmosfixture/colls/no-such-container/docs","status":"404 Not Found","statusCode":404,"errorCode":"404 Not Found","body":"{\"code\":\"NotFound\",\"message\":\"Message: {\\\"Errors\\\":[\\\"Resource Not Found. Learn more: https://aka.ms/cosmosdb-tsg-not-found\\\"]}\\r\\nActivityId: 00000000-0000-0000-0000-000000000000, Request URI: /apps/DocDbApp/services/DocDbMaster0/partitions/780e44f4-38c8-11e6-8106-8cdcd42c33be/replicas/1p/, RequestStats: , SDK: Microsoft.Azure.Documents.Common/2.14.0\"}"}}}
