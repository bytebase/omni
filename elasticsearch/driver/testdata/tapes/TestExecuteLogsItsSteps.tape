-- pin: driver docker.elastic.co/elasticsearch/elasticsearch:8.19.25 sha256:5bfdf4313b9bc05c0d46fe219d6e67710eb0c7654ed6c275167bdf9beefc56f0 8.19.25

== TestExecuteLogsItsSteps driver#1
> raw
| {"method":"PUT","route":"/logged","body":"{\"settings\": {\"number_of_shards\": 1}}\n"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"acknowledged\":true,\"index\":\"logged\",\"shards_acknowledged\":true}"}
> raw
| {"method":"GET","route":"/missing-index/_search"}
< raw
| {"status":404,"contentType":"application/json","body":"{\"error\":{\"index\":\"missing-index\",\"index_uuid\":\"_na_\",\"reason\":\"no such index [missing-index]\",\"resource.id\":\"missing-index\",\"resource.type\":\"index_or_alias\",\"root_cause\":[{\"index\":\"missing-index\",\"index_uuid\":\"_na_\",\"reason\":\"no such index [missing-index]\",\"resource.id\":\"missing-index\",\"resource.type\":\"index_or_alias\",\"type\":\"index_not_found_exception\"}],\"type\":\"index_not_found_exception\"},\"status\":404}"}
