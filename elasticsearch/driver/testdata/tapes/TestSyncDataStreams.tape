-- pin: driver docker.elastic.co/elasticsearch/elasticsearch:8.19.25 sha256:5bfdf4313b9bc05c0d46fe219d6e67710eb0c7654ed6c275167bdf9beefc56f0 8.19.25

== TestSyncDataStreams driver#1
> raw
| {"method":"PUT","route":"/alpha","body":"{\"settings\": {\"number_of_shards\": 1, \"number_of_replicas\": 0}}\n"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"acknowledged\":true,\"index\":\"alpha\",\"shards_acknowledged\":true}"}
> raw
| {"method":"PUT","route":"/_index_template/omni-logs","body":"{\"index_patterns\": [\"logs-omni-*\"], \"data_stream\": {}, \"priority\": 500, \"template\": {\"settings\": {\"number_of_shards\": 1, \"number_of_replicas\": 0}, \"mappings\": {\"properties\": {\"@timestamp\": {\"type\": \"date\"}, \"user\": {\"type\": \"keyword\"}}}}}\n"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"acknowledged\":true}"}
> raw
| {"method":"PUT","route":"/_index_template/omni-quiet","body":"{\"index_patterns\": [\"quiet-omni-*\"], \"data_stream\": {\"hidden\": true}, \"priority\": 500, \"template\": {\"settings\": {\"number_of_shards\": 1, \"number_of_replicas\": 0, \"index.hidden\": true}}}\n"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"acknowledged\":true}"}
> raw
| {"method":"PUT","route":"/_data_stream/logs-omni-app"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"acknowledged\":true}"}
> raw
| {"method":"PUT","route":"/_data_stream/quiet-omni-app"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"acknowledged\":true}"}
> raw
| {"method":"POST","route":"/logs-omni-app/_rollover"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"acknowledged\":true,\"conditions\":{},\"dry_run\":false,\"lazy\":false,\"new_index\":\".ds-logs-omni-app-0000.00.00-000002\",\"old_index\":\".ds-logs-omni-app-0000.00.00-000001\",\"rolled_over\":true,\"shards_acknowledged\":true}"}
> raw
| {"method":"POST","route":"/logs-omni-app/_bulk?refresh=true","body":"{\"create\": {}}\n{\"@timestamp\": \"2026-11-25T00:00:00Z\", \"user\": \"alice\"}\n{\"create\": {}}\n{\"@timestamp\": \"2026-11-25T00:00:01Z\", \"user\": \"bob\"}\n"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"errors\":false,\"items\":[{\"create\":{\"status\":201}},{\"create\":{\"status\":201}}]}"}
> raw
| {"method":"POST","route":"/_aliases","body":"{\"actions\": [{\"add\": {\"index\": \"logs-omni-app\", \"alias\": \"app-logs\"}}]}\n"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"acknowledged\":true,\"errors\":false}"}
> raw
| {"method":"GET","route":"/*/_settings?filter_path=**.hidden\u0026expand_wildcards=all"}
< raw
| {"status":200,"contentType":"application/json","body":"{\".ds-logs-omni-app-0000.00.00-000001\":{\"settings\":{\"index\":{\"hidden\":\"true\"}}},\".ds-logs-omni-app-0000.00.00-000002\":{\"settings\":{\"index\":{\"hidden\":\"true\"}}},\".ds-quiet-omni-app-0000.00.00-000001\":{\"settings\":{\"index\":{\"hidden\":\"true\"}}},\".security-7\":{\"settings\":{\"index\":{\"hidden\":\"true\"}}}}"}
> raw
| {"method":"GET","route":"/_data_stream"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"data_streams\":[{\"name\":\"logs-omni-app\",\"hidden\":false,\"indices\":[{\"index_name\":\".ds-logs-omni-app-0000.00.00-000001\"},{\"index_name\":\".ds-logs-omni-app-0000.00.00-000002\"}]}]}"}
> raw
| {"method":"GET","route":"/_cat/indices?format=json\u0026pretty"}
< raw
| {"status":200,"contentType":"application/json","body":"[{\"index\":\".ds-logs-omni-app-0000.00.00-000001\",\"docs.count\":\"0\",\"store.size\":\"227b\"},{\"index\":\".ds-logs-omni-app-0000.00.00-000002\",\"docs.count\":\"2\",\"store.size\":\"3.8kb\"},{\"index\":\"alpha\",\"docs.count\":\"0\",\"store.size\":\"227b\"}]"}
> raw
| {"method":"GET","route":"/_alias"}
< raw
| {"status":200,"contentType":"application/json","body":"{\".security-7\":{\"aliases\":{\".security\":{}}},\"alpha\":{\"aliases\":{}},\"logs-omni-app\":{\"aliases\":{\"app-logs\":{}}}}"}
