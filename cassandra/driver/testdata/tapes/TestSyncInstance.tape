-- pin: driver cassandra:5.0.9 sha256:acd754f07c73f3a166bb477354b56ce9a05690c3e5e2aaacc1aab2afddba6436 5.0.9

== TestSyncInstance driver#1
> raw
| {"op":"query","statement":"SELECT release_version FROM system.local"}
< raw
| {"columns":[{"name":"release_version","type":13}],"rows":[[{"kind":"string","text":"5.0.9"}]]}
> raw
| {"op":"query","statement":"SELECT keyspace_name FROM system_schema.keyspaces"}
< raw
| {"columns":[{"name":"keyspace_name","type":13}],"rows":[[{"kind":"string","text":"system_auth"}],[{"kind":"string","text":"system_schema"}],[{"kind":"string","text":"omni_instance_b"}],[{"kind":"string","text":"system_distributed"}],[{"kind":"string","text":"system"}],[{"kind":"string","text":"omni_instance_a"}],[{"kind":"string","text":"system_traces"}]]}
