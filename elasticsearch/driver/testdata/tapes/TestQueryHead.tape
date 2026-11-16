-- pin: driver docker.elastic.co/elasticsearch/elasticsearch:8.19.25 sha256:5bfdf4313b9bc05c0d46fe219d6e67710eb0c7654ed6c275167bdf9beefc56f0 8.19.25

== TestQueryHead driver#1
> raw
| {"method":"PUT","route":"/alpha","body":"{  \"settings\": {    \"number_of_shards\": 1,    \"number_of_replicas\": 0  }}\n"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"acknowledged\":true,\"index\":\"alpha\",\"shards_acknowledged\":true}"}
> raw
| {"method":"POST","route":"/alpha/_doc/1","body":"{\"name\": \"row 1\", \"n\": 1}\n"}
< raw
| {"status":201,"contentType":"application/json","body":"{\"_id\":\"1\",\"_index\":\"alpha\",\"_primary_term\":1,\"_seq_no\":0,\"_shards\":{\"failed\":0,\"successful\":1,\"total\":1},\"_version\":1,\"result\":\"created\"}"}
> raw
| {"method":"POST","route":"/alpha/_doc/2","body":"{\"name\": \"row 2\", \"n\": 2}\n"}
< raw
| {"status":201,"contentType":"application/json","body":"{\"_id\":\"2\",\"_index\":\"alpha\",\"_primary_term\":1,\"_seq_no\":1,\"_shards\":{\"failed\":0,\"successful\":1,\"total\":1},\"_version\":1,\"result\":\"created\"}"}
> raw
| {"method":"POST","route":"/alpha/_doc/3","body":"{\"name\": \"row 3\", \"n\": 3}\n"}
< raw
| {"status":201,"contentType":"application/json","body":"{\"_id\":\"3\",\"_index\":\"alpha\",\"_primary_term\":1,\"_seq_no\":2,\"_shards\":{\"failed\":0,\"successful\":1,\"total\":1},\"_version\":1,\"result\":\"created\"}"}
> raw
| {"method":"POST","route":"/alpha/_doc/big","body":"{\"big\": 9007199254740993, \"name\": \"big\"}\n"}
< raw
| {"status":201,"contentType":"application/json","body":"{\"_id\":\"big\",\"_index\":\"alpha\",\"_primary_term\":1,\"_seq_no\":3,\"_shards\":{\"failed\":0,\"successful\":1,\"total\":1},\"_version\":1,\"result\":\"created\"}"}
> raw
| {"method":"POST","route":"/alpha/_refresh"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"_shards\":{\"failed\":0,\"successful\":1,\"total\":1}}"}
> raw
| {"method":"PUT","route":"/beta","body":"{  \"settings\": {    \"number_of_shards\": 1,    \"number_of_replicas\": 0  }}\n"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"acknowledged\":true,\"index\":\"beta\",\"shards_acknowledged\":true}"}
> raw
| {"method":"POST","route":"/beta/_doc/1","body":"{\"name\": \"row 1\", \"n\": 1}\n"}
< raw
| {"status":201,"contentType":"application/json","body":"{\"_id\":\"1\",\"_index\":\"beta\",\"_primary_term\":1,\"_seq_no\":0,\"_shards\":{\"failed\":0,\"successful\":1,\"total\":1},\"_version\":1,\"result\":\"created\"}"}
> raw
| {"method":"POST","route":"/beta/_doc/2","body":"{\"name\": \"row 2\", \"n\": 2}\n"}
< raw
| {"status":201,"contentType":"application/json","body":"{\"_id\":\"2\",\"_index\":\"beta\",\"_primary_term\":1,\"_seq_no\":1,\"_shards\":{\"failed\":0,\"successful\":1,\"total\":1},\"_version\":1,\"result\":\"created\"}"}
> raw
| {"method":"POST","route":"/beta/_doc/3","body":"{\"name\": \"row 3\", \"n\": 3}\n"}
< raw
| {"status":201,"contentType":"application/json","body":"{\"_id\":\"3\",\"_index\":\"beta\",\"_primary_term\":1,\"_seq_no\":2,\"_shards\":{\"failed\":0,\"successful\":1,\"total\":1},\"_version\":1,\"result\":\"created\"}"}
> raw
| {"method":"POST","route":"/beta/_doc/4","body":"{\"name\": \"row 4\", \"n\": 4}\n"}
< raw
| {"status":201,"contentType":"application/json","body":"{\"_id\":\"4\",\"_index\":\"beta\",\"_primary_term\":1,\"_seq_no\":3,\"_shards\":{\"failed\":0,\"successful\":1,\"total\":1},\"_version\":1,\"result\":\"created\"}"}
> raw
| {"method":"POST","route":"/beta/_doc/5","body":"{\"name\": \"row 5\", \"n\": 5}\n"}
< raw
| {"status":201,"contentType":"application/json","body":"{\"_id\":\"5\",\"_index\":\"beta\",\"_primary_term\":1,\"_seq_no\":4,\"_shards\":{\"failed\":0,\"successful\":1,\"total\":1},\"_version\":1,\"result\":\"created\"}"}
> raw
| {"method":"POST","route":"/beta/_refresh"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"_shards\":{\"failed\":0,\"successful\":1,\"total\":1}}"}
> raw
| {"method":"PUT","route":"/hidden-one","body":"{  \"settings\": {    \"number_of_shards\": 1,    \"number_of_replicas\": 0,    \"index.hidden\": true  }}\n"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"acknowledged\":true,\"index\":\"hidden-one\",\"shards_acknowledged\":true}"}
> raw
| {"method":"HEAD","route":"/alpha"}
< raw
| {"status":200,"contentType":"application/json"}
> raw
| {"method":"HEAD","route":"/missing"}
< raw
| {"status":404,"contentType":"application/json"}
