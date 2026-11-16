-- pin: driver cassandra:5.0.9 sha256:acd754f07c73f3a166bb477354b56ce9a05690c3e5e2aaacc1aab2afddba6436 5.0.9

== TestSyncDatabaseEmpty driver#1
> raw
| {"op":"query","statement":"\n\t\tSELECT\n\t\t\ttable_name,\n\t\t\tcolumn_name,\n\t\t\tkind,\n\t\t\tposition,\n\t\t\tclustering_order,\n\t\t\ttype\n\t\tFROM system_schema.columns\n\t\tWHERE keyspace_name = 'omni_driver'\n\t"}
< raw
| {"columns":[{"name":"table_name","type":13},{"name":"column_name","type":13},{"name":"kind","type":13},{"name":"position","type":9},{"name":"clustering_order","type":13},{"name":"type","type":13}]}
> raw
| {"op":"query","statement":"\n\t\tSELECT\n\t\t\ttable_name,\n\t\t\tindex_name,\n\t\t\tkind,\n\t\t\ttoJson(options)\n\t\tFROM system_schema.indexes\n\t\tWHERE keyspace_name = 'omni_driver'\n\t\tORDER BY table_name, index_name\n\t"}
< raw
| {"columns":[{"name":"table_name","type":13},{"name":"index_name","type":13},{"name":"kind","type":13},{"name":"system.tojson(options)","type":13}]}
> raw
| {"op":"query","statement":"\n\t\tSELECT\n\t\t\ttable_name,\n\t\t\tcomment\n\t\tFROM system_schema.tables\n\t\tWHERE keyspace_name = 'omni_driver'\n\t\tORDER BY table_name\n\t"}
< raw
| {"columns":[{"name":"table_name","type":13},{"name":"comment","type":13}]}
