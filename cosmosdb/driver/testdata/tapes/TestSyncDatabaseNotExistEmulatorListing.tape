-- pin: emulator mcr.microsoft.com/cosmosdb/linux/azure-cosmos-emulator:2.14.28 sha256:9a4f69121f20c1d6502db84625c71ebb441ab61b01880a7b2d3d8031c6dac429 2.14.28

== TestSyncDatabaseNotExistEmulatorListing emulator#1
> raw
| {"op":"ListContainers","database":"no-such-database"}
< raw
| {"error":{"emulator":{"statusCode":404,"body":"{\"code\":\"NotFound\",\"message\":\"Message: {\\\"Errors\\\":[\\\"Owner resource does not exist\\\"]}\\r\\nActivityId: 00000000-0000-0000-0000-000000000000, Request URI: /apps/DocDbApp/services/DocDbMaster0/partitions/780e44f4-38c8-11e6-8106-8cdcd42c33be/replicas/1p/, RequestStats: , SDK: Microsoft.Azure.Documents.Common/2.14.0\"}"}}}
