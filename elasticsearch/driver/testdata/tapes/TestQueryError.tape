-- pin: driver docker.elastic.co/elasticsearch/elasticsearch:8.19.25 sha256:5bfdf4313b9bc05c0d46fe219d6e67710eb0c7654ed6c275167bdf9beefc56f0 8.19.25

== TestQueryError driver#1
> raw
| {"method":"GET","route":"/missing-index/_search","body":"{\"query\": {\"match_all\": {}}}\n"}
< raw
| {"status":404,"contentType":"application/json","body":"{\"error\":{\"index\":\"missing-index\",\"index_uuid\":\"_na_\",\"reason\":\"no such index [missing-index]\",\"resource.id\":\"missing-index\",\"resource.type\":\"index_or_alias\",\"root_cause\":[{\"index\":\"missing-index\",\"index_uuid\":\"_na_\",\"reason\":\"no such index [missing-index]\",\"resource.id\":\"missing-index\",\"resource.type\":\"index_or_alias\",\"type\":\"index_not_found_exception\"}],\"type\":\"index_not_found_exception\"},\"status\":404}"}
