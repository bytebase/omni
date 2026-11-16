-- pin: emulator mcr.microsoft.com/cosmosdb/linux/azure-cosmos-emulator:2.14.28 sha256:9a4f69121f20c1d6502db84625c71ebb441ab61b01880a7b2d3d8031c6dac429 2.14.28

== TestQueryExplain emulator#1
> raw
| {"op":"Plan","database":"cosmosfixture","container":"cities","statement":"SELECT * FROM c WHERE c.country = \"US\""}
< raw
| {"plan":{"partitionedQueryExecutionInfoVersion":2,"queryInfo":{"distinctType":"None","top":null,"offset":null,"limit":null,"orderBy":[],"orderByExpressions":[],"groupByExpressions":[],"groupByAliases":[],"aggregates":[],"groupByAliasToAggregateType":{},"rewrittenQuery":"","hasSelectValue":false,"dCountInfo":null,"hasNonStreamingOrderBy":false},"queryRanges":[{"min":"05C17769B99008565400","max":"05C17769B99008565400","isMinInclusive":true,"isMaxInclusive":true}],"hybridSearchQueryInfo":null}}
> raw
| {"op":"Plan","database":"cosmosfixture","container":"cities","statement":"SELECT c.seq FROM c ORDER BY c.seq"}
< raw
| {"plan":{"partitionedQueryExecutionInfoVersion":2,"queryInfo":{"distinctType":"None","top":null,"offset":null,"limit":null,"orderBy":["Ascending"],"orderByExpressions":["c.seq"],"groupByExpressions":[],"groupByAliases":[],"aggregates":[],"groupByAliasToAggregateType":{},"rewrittenQuery":"SELECT c._rid, [{\"item\": c.seq}] AS orderByItems, {\"seq\": c.seq} AS payload\nFROM c\nWHERE ({documentdb-formattableorderbyquery-filter})\nORDER BY c.seq","hasSelectValue":false,"dCountInfo":null,"hasNonStreamingOrderBy":false},"queryRanges":[{"min":"","max":"FF","isMinInclusive":true,"isMaxInclusive":false}],"hybridSearchQueryInfo":null}}
> raw
| {"op":"Plan","database":"cosmosfixture","container":"cities","statement":"SELECT VALUE COUNT(1) FROM c"}
< raw
| {"plan":{"partitionedQueryExecutionInfoVersion":2,"queryInfo":{"distinctType":"None","top":null,"offset":null,"limit":null,"orderBy":[],"orderByExpressions":[],"groupByExpressions":[],"groupByAliases":[],"aggregates":["Count"],"groupByAliasToAggregateType":{},"rewrittenQuery":"SELECT VALUE [{\"item\": COUNT(1)}]\nFROM c","hasSelectValue":true,"dCountInfo":null,"hasNonStreamingOrderBy":false},"queryRanges":[{"min":"","max":"FF","isMinInclusive":true,"isMaxInclusive":false}],"hybridSearchQueryInfo":null}}
> raw
| {"op":"Plan","database":"cosmosfixture","container":"cities","statement":"SELECT * FROM c WHERE c.country = \"US\""}
< raw
| {"plan":{"partitionedQueryExecutionInfoVersion":2,"queryInfo":{"distinctType":"None","top":null,"offset":null,"limit":null,"orderBy":[],"orderByExpressions":[],"groupByExpressions":[],"groupByAliases":[],"aggregates":[],"groupByAliasToAggregateType":{},"rewrittenQuery":"","hasSelectValue":false,"dCountInfo":null,"hasNonStreamingOrderBy":false},"queryRanges":[{"min":"05C17769B99008565400","max":"05C17769B99008565400","isMinInclusive":true,"isMaxInclusive":true}],"hybridSearchQueryInfo":null}}
> raw
| {"op":"Plan","database":"cosmosfixture","container":"cities","statement":" SELECT c.seq FROM c ORDER BY c.seq"}
< raw
| {"plan":{"partitionedQueryExecutionInfoVersion":2,"queryInfo":{"distinctType":"None","top":null,"offset":null,"limit":null,"orderBy":["Ascending"],"orderByExpressions":["c.seq"],"groupByExpressions":[],"groupByAliases":[],"aggregates":[],"groupByAliasToAggregateType":{},"rewrittenQuery":"SELECT c._rid, [{\"item\": c.seq}] AS orderByItems, {\"seq\": c.seq} AS payload\nFROM c\nWHERE ({documentdb-formattableorderbyquery-filter})\nORDER BY c.seq","hasSelectValue":false,"dCountInfo":null,"hasNonStreamingOrderBy":false},"queryRanges":[{"min":"","max":"FF","isMinInclusive":true,"isMaxInclusive":false}],"hybridSearchQueryInfo":null}}
> raw
| {"op":"Plan","database":"cosmosfixture","container":"cities","statement":"SELECT c.seq FROM c ORDER BY c.seq"}
< raw
| {"plan":{"partitionedQueryExecutionInfoVersion":2,"queryInfo":{"distinctType":"None","top":null,"offset":null,"limit":null,"orderBy":["Ascending"],"orderByExpressions":["c.seq"],"groupByExpressions":[],"groupByAliases":[],"aggregates":[],"groupByAliasToAggregateType":{},"rewrittenQuery":"SELECT c._rid, [{\"item\": c.seq}] AS orderByItems, {\"seq\": c.seq} AS payload\nFROM c\nWHERE ({documentdb-formattableorderbyquery-filter})\nORDER BY c.seq","hasSelectValue":false,"dCountInfo":null,"hasNonStreamingOrderBy":false},"queryRanges":[{"min":"","max":"FF","isMinInclusive":true,"isMaxInclusive":false}],"hybridSearchQueryInfo":null}}
