-- pin: driver docker.elastic.co/elasticsearch/elasticsearch:8.19.25 sha256:5bfdf4313b9bc05c0d46fe219d6e67710eb0c7654ed6c275167bdf9beefc56f0 8.19.25

== TestSyncDatabaseAliases driver#1
> raw
| {"method":"PUT","route":"/alpha","body":"{\"settings\": {\"number_of_shards\": 1, \"number_of_replicas\": 0}}\n"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"acknowledged\":true,\"index\":\"alpha\",\"shards_acknowledged\":true}"}
> raw
| {"method":"PUT","route":"/beta","body":"{\"settings\": {\"number_of_shards\": 1, \"number_of_replicas\": 0}}\n"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"acknowledged\":true,\"index\":\"beta\",\"shards_acknowledged\":true}"}
> raw
| {"method":"PUT","route":"/hidden-one","body":"{\"settings\": {\"number_of_shards\": 1, \"number_of_replicas\": 0, \"index.hidden\": true}}\n"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"acknowledged\":true,\"index\":\"hidden-one\",\"shards_acknowledged\":true}"}
> raw
| {"method":"POST","route":"/_aliases","body":"{\"actions\": [{\"add\": {\"index\": \"alpha\", \"alias\": \"both\"}}, {\"add\": {\"index\": \"beta\", \"alias\": \"both\"}}, {\"add\": {\"index\": \"beta\", \"alias\": \"beta-recent\", \"filter\": {\"range\": {\"n\": {\"gte\": 3}}}}}, {\"add\": {\"index\": \"hidden-one\", \"alias\": \"hidden-alias\"}}]}\n"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"acknowledged\":true,\"errors\":false}"}
> raw
| {"method":"GET","route":"/*/_settings?filter_path=**.hidden\u0026expand_wildcards=all"}
< raw
| {"status":200,"contentType":"application/json","body":"{\".security-7\":{\"settings\":{\"index\":{\"hidden\":\"true\"}}},\"hidden-one\":{\"settings\":{\"index\":{\"hidden\":\"true\"}}}}"}
> raw
| {"method":"GET","route":"/_data_stream"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"data_streams\":[]}"}
> raw
| {"method":"GET","route":"/_cat/indices?format=json\u0026pretty"}
< raw
| {"status":200,"contentType":"application/json","body":"[{\"index\":\"alpha\",\"docs.count\":\"0\",\"store.size\":\"227b\"},{\"index\":\"beta\",\"docs.count\":\"0\",\"store.size\":\"227b\"},{\"index\":\"hidden-one\",\"docs.count\":\"0\",\"store.size\":\"227b\"}]"}
> raw
| {"method":"GET","route":"/_alias"}
< raw
| {"status":200,"contentType":"application/json","body":"{\".security-7\":{\"aliases\":{\".security\":{}}},\"alpha\":{\"aliases\":{\"both\":{}}},\"beta\":{\"aliases\":{\"beta-recent\":{},\"both\":{}}},\"hidden-one\":{\"aliases\":{\"hidden-alias\":{}}}}"}
