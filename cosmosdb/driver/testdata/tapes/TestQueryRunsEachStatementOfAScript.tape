-- pin: emulator mcr.microsoft.com/cosmosdb/linux/azure-cosmos-emulator:2.14.28 sha256:9a4f69121f20c1d6502db84625c71ebb441ab61b01880a7b2d3d8031c6dac429 2.14.28

== TestQueryRunsEachStatementOfAScript emulator#1
> raw
| {"op":"Page","database":"cosmosfixture","container":"cities","statement":"SELECT c.name FROM c WHERE c.country = \"US\""}
< raw
| {"items":[{"name":"city-US-0"},{"name":"city-US-1"},{"name":"city-US-2"},{"name":"city-US-3"},{"name":"city-US-4"}],"pages":1}
> raw
| {"op":"Page","database":"cosmosfixture","container":"cities","statement":" SELECT c.name FROM c WHERE c.country = \"AE\""}
< raw
| {"items":[{"name":"city-AE-0"},{"name":"city-AE-1"},{"name":"city-AE-2"},{"name":"city-AE-3"},{"name":"city-AE-4"}],"pages":1}
