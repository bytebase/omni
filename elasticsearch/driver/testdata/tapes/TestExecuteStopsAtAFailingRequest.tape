-- pin: driver docker.elastic.co/elasticsearch/elasticsearch:8.19.25 sha256:5bfdf4313b9bc05c0d46fe219d6e67710eb0c7654ed6c275167bdf9beefc56f0 8.19.25

== TestExecuteStopsAtAFailingRequest driver#1
> raw
| {"method":"PUT","route":"/kept","body":"{\"settings\": {\"number_of_shards\": 1, \"number_of_replicas\": 0}}\n"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"acknowledged\":true,\"index\":\"kept\",\"shards_acknowledged\":true}"}
> raw
| {"method":"GET","route":"/missing-index/_search","body":"{\"query\": {\"match_all\": {}}}\n"}
< raw
| {"status":404,"contentType":"application/json","body":"{\"error\":{\"index\":\"missing-index\",\"index_uuid\":\"_na_\",\"reason\":\"no such index [missing-index]\",\"resource.id\":\"missing-index\",\"resource.type\":\"index_or_alias\",\"root_cause\":[{\"index\":\"missing-index\",\"index_uuid\":\"_na_\",\"reason\":\"no such index [missing-index]\",\"resource.id\":\"missing-index\",\"resource.type\":\"index_or_alias\",\"type\":\"index_not_found_exception\"}],\"type\":\"index_not_found_exception\"},\"status\":404}"}
> raw
| {"method":"GET","route":"/*/_settings?filter_path=**.hidden\u0026expand_wildcards=all"}
< raw
| {"status":200,"contentType":"application/json","body":"{\".security-7\":{\"settings\":{\"index\":{\"hidden\":\"true\"}}}}"}
> raw
| {"method":"GET","route":"/_data_stream"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"data_streams\":[]}"}
> raw
| {"method":"GET","route":"/_cat/indices?format=json\u0026pretty"}
< raw
| {"status":200,"contentType":"application/json","body":"[{\"index\":\"kept\",\"docs.count\":\"0\",\"store.size\":\"227b\"}]"}
> raw
| {"method":"GET","route":"/_alias"}
< raw
| {"status":200,"contentType":"application/json","body":"{\".security-7\":{\"aliases\":{\".security\":{}}},\"kept\":{\"aliases\":{}}}"}
