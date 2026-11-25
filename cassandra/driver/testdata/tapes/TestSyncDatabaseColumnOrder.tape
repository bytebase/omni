-- pin: driver cassandra:5.0.9 sha256:acd754f07c73f3a166bb477354b56ce9a05690c3e5e2aaacc1aab2afddba6436 5.0.9

== TestSyncDatabaseColumnOrder driver#1
> raw
| {"op":"exec","statement":"CREATE TABLE readings (\n    sensor text,\n    day date,\n    seq int,\n    alpha text,\n    zone text STATIC,\n    beta text,\n    PRIMARY KEY ((sensor, day), seq)\n)"}
< raw
| {}
> raw
| {"op":"exec","statement":"CREATE MATERIALIZED VIEW readings_by_alpha AS\n    SELECT sensor, day, seq, alpha, beta FROM readings\n    WHERE alpha IS NOT NULL AND sensor IS NOT NULL AND day IS NOT NULL AND seq IS NOT NULL\n    PRIMARY KEY (alpha, sensor, day, seq)"}
< raw
| {}
> raw
| {"op":"query","statement":"SELECT keyspace_name FROM system_schema.keyspaces WHERE keyspace_name = 'omni_driver'"}
< raw
| {"columns":[{"name":"keyspace_name","type":13}],"rowColumns":["keyspace_name"],"rows":[[{"kind":"string","text":"omni_driver"}]]}
> raw
| {"op":"query","statement":"\n\t\tSELECT\n\t\t\ttable_name,\n\t\t\tcolumn_name,\n\t\t\tkind,\n\t\t\tposition,\n\t\t\tclustering_order,\n\t\t\ttype\n\t\tFROM system_schema.columns\n\t\tWHERE keyspace_name = 'omni_driver'\n\t"}
< raw
| {"columns":[{"name":"table_name","type":13},{"name":"column_name","type":13},{"name":"kind","type":13},{"name":"position","type":9},{"name":"clustering_order","type":13},{"name":"type","type":13}],"rowColumns":["table_name","column_name","kind","position","clustering_order","type"],"rows":[[{"kind":"string","text":"readings"},{"kind":"string","text":"alpha"},{"kind":"string","text":"regular"},{"kind":"int","text":"-1"},{"kind":"string","text":"none"},{"kind":"string","text":"text"}],[{"kind":"string","text":"readings"},{"kind":"string","text":"beta"},{"kind":"string","text":"regular"},{"kind":"int","text":"-1"},{"kind":"string","text":"none"},{"kind":"string","text":"text"}],[{"kind":"string","text":"readings"},{"kind":"string","text":"day"},{"kind":"string","text":"partition_key"},{"kind":"int","text":"1"},{"kind":"string","text":"none"},{"kind":"string","text":"date"}],[{"kind":"string","text":"readings"},{"kind":"string","text":"sensor"},{"kind":"string","text":"partition_key"},{"kind":"int","text":"0"},{"kind":"string","text":"none"},{"kind":"string","text":"text"}],[{"kind":"string","text":"readings"},{"kind":"string","text":"seq"},{"kind":"string","text":"clustering"},{"kind":"int","text":"0"},{"kind":"string","text":"asc"},{"kind":"string","text":"int"}],[{"kind":"string","text":"readings"},{"kind":"string","text":"zone"},{"kind":"string","text":"static"},{"kind":"int","text":"-1"},{"kind":"string","text":"none"},{"kind":"string","text":"text"}],[{"kind":"string","text":"readings_by_alpha"},{"kind":"string","text":"alpha"},{"kind":"string","text":"partition_key"},{"kind":"int","text":"0"},{"kind":"string","text":"none"},{"kind":"string","text":"text"}],[{"kind":"string","text":"readings_by_alpha"},{"kind":"string","text":"beta"},{"kind":"string","text":"regular"},{"kind":"int","text":"-1"},{"kind":"string","text":"none"},{"kind":"string","text":"text"}],[{"kind":"string","text":"readings_by_alpha"},{"kind":"string","text":"day"},{"kind":"string","text":"clustering"},{"kind":"int","text":"1"},{"kind":"string","text":"asc"},{"kind":"string","text":"date"}],[{"kind":"string","text":"readings_by_alpha"},{"kind":"string","text":"sensor"},{"kind":"string","text":"clustering"},{"kind":"int","text":"0"},{"kind":"string","text":"asc"},{"kind":"string","text":"text"}],[{"kind":"string","text":"readings_by_alpha"},{"kind":"string","text":"seq"},{"kind":"string","text":"clustering"},{"kind":"int","text":"2"},{"kind":"string","text":"asc"},{"kind":"string","text":"int"}]]}
> raw
| {"op":"query","statement":"\n\t\tSELECT\n\t\t\ttable_name,\n\t\t\tindex_name,\n\t\t\tkind,\n\t\t\ttoJson(options)\n\t\tFROM system_schema.indexes\n\t\tWHERE keyspace_name = 'omni_driver'\n\t\tORDER BY table_name, index_name\n\t"}
< raw
| {"columns":[{"name":"table_name","type":13},{"name":"index_name","type":13},{"name":"kind","type":13},{"name":"system.tojson(options)","type":13}]}
> raw
| {"op":"query","statement":"\n\t\tSELECT\n\t\t\ttable_name,\n\t\t\tcomment\n\t\tFROM system_schema.tables\n\t\tWHERE keyspace_name = 'omni_driver'\n\t\tORDER BY table_name\n\t"}
< raw
| {"columns":[{"name":"table_name","type":13},{"name":"comment","type":13}],"rowColumns":["table_name","comment"],"rows":[[{"kind":"string","text":"readings"},{"kind":"string"}]]}
> raw
| {"op":"query","statement":"\n\t\tSELECT\n\t\t\tview_name,\n\t\t\tbase_table_name\n\t\tFROM system_schema.views\n\t\tWHERE keyspace_name = 'omni_driver'\n\t\tORDER BY view_name\n\t"}
< raw
| {"columns":[{"name":"view_name","type":13},{"name":"base_table_name","type":13}],"rowColumns":["view_name","base_table_name"],"rows":[[{"kind":"string","text":"readings_by_alpha"},{"kind":"string","text":"readings"}]]}
