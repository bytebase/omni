-- pin: driver cassandra:5.0.9 sha256:acd754f07c73f3a166bb477354b56ce9a05690c3e5e2aaacc1aab2afddba6436 5.0.9

== TestSyncDatabase driver#1
> raw
| {"op":"exec","statement":"-- The tables and rows every cassandra/driver test records against. A recording\n-- starts from an empty keyspace, so the tape holds these statements' answers\n-- and nothing a previous test left behind.\n\nCREATE TABLE omni_types (\n    id int PRIMARY KEY,\n    txt text,\n    ascii_col ascii,\n    vc varchar,\n    flag boolean,\n    bi bigint,\n    ti tinyint,\n    si smallint,\n    dbl double,\n    flt float,\n    dec decimal,\n    vi varint,\n    ts timestamp,\n    dt date,\n    tm time,\n    uid uuid,\n    tuid timeuuid,\n    bl blob,\n    dur duration,\n    lst list\u003cint\u003e,\n    st set\u003ctext\u003e,\n    mp map\u003ctext,int\u003e,\n    addr inet,\n    tup tuple\u003cint, text\u003e\n)"}
< raw
| {}
> raw
| {"op":"exec","statement":"INSERT INTO omni_types (\n    id, txt, ascii_col, vc, flag, bi, ti, si, dbl, flt, dec, vi, ts, dt, tm, uid, tuid, bl, dur, lst, st, mp, addr, tup\n) VALUES (\n    1, 'text value', 'ascii', 'varchar', true, 9223372036854775807, 7, 300, 3.5, 1.25, 1.5,\n    123456789012345678901234567890, '2026-01-02 03:04:05.678', '2026-01-02', '08:12:54.123456789',\n    6afaf7c2-c0c2-0ddc-7f22-71b8f0a0b0c0, 50554d6e-29bb-11e5-b345-feff819cdc9f,\n    0xdeadbeef, 1mo2d3h, [1, 2, 3], {'a', 'b'}, {'a': 1, 'b': 2}, '10.0.0.1', (7, 'seven')\n)"}
< raw
| {}
> raw
| {"op":"exec","statement":"INSERT INTO omni_types (id, txt) VALUES (2, NULL)"}
< raw
| {}
> raw
| {"op":"exec","statement":"CREATE TABLE omni_rows (\n    id int PRIMARY KEY,\n    label text\n)"}
< raw
| {}
> raw
| {"op":"exec","statement":"INSERT INTO omni_rows (id, label) VALUES (1, 'row-1')"}
< raw
| {}
> raw
| {"op":"exec","statement":"INSERT INTO omni_rows (id, label) VALUES (2, 'row-2')"}
< raw
| {}
> raw
| {"op":"exec","statement":"INSERT INTO omni_rows (id, label) VALUES (3, 'row-3')"}
< raw
| {}
> raw
| {"op":"exec","statement":"INSERT INTO omni_rows (id, label) VALUES (4, 'row-4')"}
< raw
| {}
> raw
| {"op":"exec","statement":"INSERT INTO omni_rows (id, label) VALUES (5, 'row-5')"}
< raw
| {}
> raw
| {"op":"exec","statement":"INSERT INTO omni_rows (id, label) VALUES (6, 'row-6')"}
< raw
| {}
> raw
| {"op":"exec","statement":"INSERT INTO omni_rows (id, label) VALUES (7, 'row-7')"}
< raw
| {}
> raw
| {"op":"exec","statement":"INSERT INTO omni_rows (id, label) VALUES (8, 'row-8')"}
< raw
| {}
> raw
| {"op":"exec","statement":"INSERT INTO omni_rows (id, label) VALUES (9, 'row-9')"}
< raw
| {}
> raw
| {"op":"exec","statement":"INSERT INTO omni_rows (id, label) VALUES (10, 'row-10')"}
< raw
| {}
> raw
| {"op":"exec","statement":"CREATE TABLE omni_events (\n    group_id int,\n    bucket int,\n    event_id int,\n    payload text,\n    note text static,\n    PRIMARY KEY ((group_id), bucket, event_id)\n) WITH CLUSTERING ORDER BY (bucket DESC, event_id ASC)"}
< raw
| {}
> raw
| {"op":"exec","statement":"INSERT INTO omni_events (group_id, bucket, event_id, payload, note) VALUES (1, 2, 3, 'p1', 'first group')"}
< raw
| {}
> raw
| {"op":"exec","statement":"INSERT INTO omni_events (group_id, bucket, event_id, payload) VALUES (1, 1, 4, 'p2')"}
< raw
| {}
> raw
| {"op":"exec","statement":"INSERT INTO omni_events (group_id, bucket, event_id, payload, note) VALUES (2, 9, 1, 'p3', 'second group')"}
< raw
| {}
> raw
| {"op":"exec","statement":"CREATE INDEX omni_events_payload_idx ON omni_events (payload)"}
< raw
| {}
> raw
| {"op":"query","statement":"SELECT keyspace_name FROM system_schema.keyspaces WHERE keyspace_name = 'omni_driver'"}
< raw
| {"columns":[{"name":"keyspace_name","type":13}],"rowColumns":["keyspace_name"],"rows":[[{"kind":"string","text":"omni_driver"}]]}
> raw
| {"op":"query","statement":"\n\t\tSELECT\n\t\t\ttable_name,\n\t\t\tcolumn_name,\n\t\t\tkind,\n\t\t\tposition,\n\t\t\tclustering_order,\n\t\t\ttype\n\t\tFROM system_schema.columns\n\t\tWHERE keyspace_name = 'omni_driver'\n\t"}
< raw
| {"columns":[{"name":"table_name","type":13},{"name":"column_name","type":13},{"name":"kind","type":13},{"name":"position","type":9},{"name":"clustering_order","type":13},{"name":"type","type":13}],"rowColumns":["table_name","column_name","kind","position","clustering_order","type"],"rows":[[{"kind":"string","text":"omni_events"},{"kind":"string","text":"bucket"},{"kind":"string","text":"clustering"},{"kind":"int","text":"0"},{"kind":"string","text":"desc"},{"kind":"string","text":"int"}],[{"kind":"string","text":"omni_events"},{"kind":"string","text":"event_id"},{"kind":"string","text":"clustering"},{"kind":"int","text":"1"},{"kind":"string","text":"asc"},{"kind":"string","text":"int"}],[{"kind":"string","text":"omni_events"},{"kind":"string","text":"group_id"},{"kind":"string","text":"partition_key"},{"kind":"int","text":"0"},{"kind":"string","text":"none"},{"kind":"string","text":"int"}],[{"kind":"string","text":"omni_events"},{"kind":"string","text":"note"},{"kind":"string","text":"static"},{"kind":"int","text":"-1"},{"kind":"string","text":"none"},{"kind":"string","text":"text"}],[{"kind":"string","text":"omni_events"},{"kind":"string","text":"payload"},{"kind":"string","text":"regular"},{"kind":"int","text":"-1"},{"kind":"string","text":"none"},{"kind":"string","text":"text"}],[{"kind":"string","text":"omni_rows"},{"kind":"string","text":"id"},{"kind":"string","text":"partition_key"},{"kind":"int","text":"0"},{"kind":"string","text":"none"},{"kind":"string","text":"int"}],[{"kind":"string","text":"omni_rows"},{"kind":"string","text":"label"},{"kind":"string","text":"regular"},{"kind":"int","text":"-1"},{"kind":"string","text":"none"},{"kind":"string","text":"text"}],[{"kind":"string","text":"omni_types"},{"kind":"string","text":"addr"},{"kind":"string","text":"regular"},{"kind":"int","text":"-1"},{"kind":"string","text":"none"},{"kind":"string","text":"inet"}],[{"kind":"string","text":"omni_types"},{"kind":"string","text":"ascii_col"},{"kind":"string","text":"regular"},{"kind":"int","text":"-1"},{"kind":"string","text":"none"},{"kind":"string","text":"ascii"}],[{"kind":"string","text":"omni_types"},{"kind":"string","text":"bi"},{"kind":"string","text":"regular"},{"kind":"int","text":"-1"},{"kind":"string","text":"none"},{"kind":"string","text":"bigint"}],[{"kind":"string","text":"omni_types"},{"kind":"string","text":"bl"},{"kind":"string","text":"regular"},{"kind":"int","text":"-1"},{"kind":"string","text":"none"},{"kind":"string","text":"blob"}],[{"kind":"string","text":"omni_types"},{"kind":"string","text":"dbl"},{"kind":"string","text":"regular"},{"kind":"int","text":"-1"},{"kind":"string","text":"none"},{"kind":"string","text":"double"}],[{"kind":"string","text":"omni_types"},{"kind":"string","text":"dec"},{"kind":"string","text":"regular"},{"kind":"int","text":"-1"},{"kind":"string","text":"none"},{"kind":"string","text":"decimal"}],[{"kind":"string","text":"omni_types"},{"kind":"string","text":"dt"},{"kind":"string","text":"regular"},{"kind":"int","text":"-1"},{"kind":"string","text":"none"},{"kind":"string","text":"date"}],[{"kind":"string","text":"omni_types"},{"kind":"string","text":"dur"},{"kind":"string","text":"regular"},{"kind":"int","text":"-1"},{"kind":"string","text":"none"},{"kind":"string","text":"duration"}],[{"kind":"string","text":"omni_types"},{"kind":"string","text":"flag"},{"kind":"string","text":"regular"},{"kind":"int","text":"-1"},{"kind":"string","text":"none"},{"kind":"string","text":"boolean"}],[{"kind":"string","text":"omni_types"},{"kind":"string","text":"flt"},{"kind":"string","text":"regular"},{"kind":"int","text":"-1"},{"kind":"string","text":"none"},{"kind":"string","text":"float"}],[{"kind":"string","text":"omni_types"},{"kind":"string","text":"id"},{"kind":"string","text":"partition_key"},{"kind":"int","text":"0"},{"kind":"string","text":"none"},{"kind":"string","text":"int"}],[{"kind":"string","text":"omni_types"},{"kind":"string","text":"lst"},{"kind":"string","text":"regular"},{"kind":"int","text":"-1"},{"kind":"string","text":"none"},{"kind":"string","text":"list\u003cint\u003e"}],[{"kind":"string","text":"omni_types"},{"kind":"string","text":"mp"},{"kind":"string","text":"regular"},{"kind":"int","text":"-1"},{"kind":"string","text":"none"},{"kind":"string","text":"map\u003ctext, int\u003e"}],[{"kind":"string","text":"omni_types"},{"kind":"string","text":"si"},{"kind":"string","text":"regular"},{"kind":"int","text":"-1"},{"kind":"string","text":"none"},{"kind":"string","text":"smallint"}],[{"kind":"string","text":"omni_types"},{"kind":"string","text":"st"},{"kind":"string","text":"regular"},{"kind":"int","text":"-1"},{"kind":"string","text":"none"},{"kind":"string","text":"set\u003ctext\u003e"}],[{"kind":"string","text":"omni_types"},{"kind":"string","text":"ti"},{"kind":"string","text":"regular"},{"kind":"int","text":"-1"},{"kind":"string","text":"none"},{"kind":"string","text":"tinyint"}],[{"kind":"string","text":"omni_types"},{"kind":"string","text":"tm"},{"kind":"string","text":"regular"},{"kind":"int","text":"-1"},{"kind":"string","text":"none"},{"kind":"string","text":"time"}],[{"kind":"string","text":"omni_types"},{"kind":"string","text":"ts"},{"kind":"string","text":"regular"},{"kind":"int","text":"-1"},{"kind":"string","text":"none"},{"kind":"string","text":"timestamp"}],[{"kind":"string","text":"omni_types"},{"kind":"string","text":"tuid"},{"kind":"string","text":"regular"},{"kind":"int","text":"-1"},{"kind":"string","text":"none"},{"kind":"string","text":"timeuuid"}],[{"kind":"string","text":"omni_types"},{"kind":"string","text":"tup"},{"kind":"string","text":"regular"},{"kind":"int","text":"-1"},{"kind":"string","text":"none"},{"kind":"string","text":"frozen\u003ctuple\u003cint, text\u003e\u003e"}],[{"kind":"string","text":"omni_types"},{"kind":"string","text":"txt"},{"kind":"string","text":"regular"},{"kind":"int","text":"-1"},{"kind":"string","text":"none"},{"kind":"string","text":"text"}],[{"kind":"string","text":"omni_types"},{"kind":"string","text":"uid"},{"kind":"string","text":"regular"},{"kind":"int","text":"-1"},{"kind":"string","text":"none"},{"kind":"string","text":"uuid"}],[{"kind":"string","text":"omni_types"},{"kind":"string","text":"vc"},{"kind":"string","text":"regular"},{"kind":"int","text":"-1"},{"kind":"string","text":"none"},{"kind":"string","text":"text"}],[{"kind":"string","text":"omni_types"},{"kind":"string","text":"vi"},{"kind":"string","text":"regular"},{"kind":"int","text":"-1"},{"kind":"string","text":"none"},{"kind":"string","text":"varint"}]]}
> raw
| {"op":"query","statement":"\n\t\tSELECT\n\t\t\ttable_name,\n\t\t\tindex_name,\n\t\t\tkind,\n\t\t\ttoJson(options)\n\t\tFROM system_schema.indexes\n\t\tWHERE keyspace_name = 'omni_driver'\n\t\tORDER BY table_name, index_name\n\t"}
< raw
| {"columns":[{"name":"table_name","type":13},{"name":"index_name","type":13},{"name":"kind","type":13},{"name":"system.tojson(options)","type":13}],"rowColumns":["table_name","index_name","kind","system.tojson(options)"],"rows":[[{"kind":"string","text":"omni_events"},{"kind":"string","text":"omni_events_payload_idx"},{"kind":"string","text":"COMPOSITES"},{"kind":"string","text":"{\"target\": \"payload\"}"}]]}
> raw
| {"op":"query","statement":"\n\t\tSELECT\n\t\t\ttable_name,\n\t\t\tcomment\n\t\tFROM system_schema.tables\n\t\tWHERE keyspace_name = 'omni_driver'\n\t\tORDER BY table_name\n\t"}
< raw
| {"columns":[{"name":"table_name","type":13},{"name":"comment","type":13}],"rowColumns":["table_name","comment"],"rows":[[{"kind":"string","text":"omni_events"},{"kind":"string"}],[{"kind":"string","text":"omni_rows"},{"kind":"string"}],[{"kind":"string","text":"omni_types"},{"kind":"string"}]]}
