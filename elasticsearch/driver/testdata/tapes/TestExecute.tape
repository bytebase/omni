-- pin: driver docker.elastic.co/elasticsearch/elasticsearch:8.19.25 sha256:5bfdf4313b9bc05c0d46fe219d6e67710eb0c7654ed6c275167bdf9beefc56f0 8.19.25

== TestExecute driver#1
> raw
| {"method":"PUT","route":"/exec_items","body":"{\"settings\": {\"number_of_shards\": 1, \"number_of_replicas\": 0}}\n"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"acknowledged\":true,\"index\":\"exec_items\",\"shards_acknowledged\":true}"}
> raw
| {"method":"POST","route":"/exec_items/_doc/1","body":"{\"note\": \"a\"}\n"}
< raw
| {"status":201,"contentType":"application/json","body":"{\"_id\":\"1\",\"_index\":\"exec_items\",\"_primary_term\":1,\"_seq_no\":0,\"_shards\":{\"failed\":0,\"successful\":1,\"total\":1},\"_version\":1,\"result\":\"created\"}"}
> raw
| {"method":"POST","route":"/exec_items/_refresh"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"_shards\":{\"failed\":0,\"successful\":1,\"total\":1}}"}
> raw
| {"method":"GET","route":"/*/_settings?filter_path=**.hidden\u0026expand_wildcards=open,hidden"}
< raw
| {"status":200,"contentType":"application/json","body":"{\".security-7\":{\"settings\":{\"index\":{\"hidden\":\"true\"}}}}"}
> raw
| {"method":"GET","route":"/_cat/indices?format=json\u0026pretty"}
< raw
| {"status":200,"contentType":"application/json","body":"[{\"index\":\"exec_items\",\"docs.count\":\"1\",\"store.size\":\"4.5kb\"}]"}
> raw
| {"method":"GET","route":"/_alias"}
< raw
| {"status":200,"contentType":"application/json","body":"{\".security-7\":{\"aliases\":{\".security\":{}}},\"exec_items\":{\"aliases\":{}}}"}
