-- pin: driver docker.elastic.co/elasticsearch/elasticsearch:8.19.25 sha256:5bfdf4313b9bc05c0d46fe219d6e67710eb0c7654ed6c275167bdf9beefc56f0 8.19.25

== TestQueryExplain driver#1
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
| {"method":"GET","route":"/alpha/_search","body":"{\"profile\":true,\"query\":{\"match_all\":{}}}\n"}
< raw
| {"status":200,"contentType":"application/json","body":"{\"_shards\":{\"failed\":0,\"skipped\":0,\"successful\":1,\"total\":1},\"hits\":{\"hits\":[{\"_id\":\"1\",\"_index\":\"alpha\",\"_score\":1.0,\"_source\":{\"n\":1,\"name\":\"row 1\"}},{\"_id\":\"2\",\"_index\":\"alpha\",\"_score\":1.0,\"_source\":{\"n\":2,\"name\":\"row 2\"}},{\"_id\":\"3\",\"_index\":\"alpha\",\"_score\":1.0,\"_source\":{\"n\":3,\"name\":\"row 3\"}}],\"max_score\":1.0,\"total\":{\"relation\":\"eq\",\"value\":3}},\"profile\":{\"shards\":[{\"aggregations\":[],\"cluster\":\"(local)\",\"fetch\":{\"breakdown\":{\"load_source\":0,\"load_source_count\":0,\"load_stored_fields\":0,\"load_stored_fields_count\":0,\"next_reader\":0,\"next_reader_count\":0},\"children\":[{\"breakdown\":{\"next_reader\":0,\"next_reader_count\":0,\"process\":0,\"process_count\":0},\"description\":\"\",\"time_in_nanos\":0,\"type\":\"FetchFieldsPhase\"},{\"breakdown\":{\"next_reader\":0,\"next_reader_count\":0,\"process\":0,\"process_count\":0},\"debug\":{\"fast_path\":3},\"description\":\"\",\"time_in_nanos\":0,\"type\":\"FetchSourcePhase\"},{\"breakdown\":{\"next_reader\":0,\"next_reader_count\":0,\"process\":0,\"process_count\":0},\"description\":\"\",\"time_in_nanos\":0,\"type\":\"StoredFieldsPhase\"}],\"debug\":{\"stored_fields\":[\"_id\",\"_routing\",\"_source\"]},\"description\":\"\",\"time_in_nanos\":0,\"type\":\"fetch\"},\"id\":\"0\",\"index\":\"alpha\",\"node_id\":\"0\",\"searches\":[{\"collector\":[{\"children\":[{\"name\":\"SimpleTopScoreDocCollector\",\"reason\":\"search_top_hits\",\"time_in_nanos\":0}],\"name\":\"QueryPhaseCollector\",\"reason\":\"search_query_phase\",\"time_in_nanos\":0}],\"query\":[{\"breakdown\":{\"advance\":0,\"advance_count\":0,\"build_scorer\":0,\"build_scorer_count\":0,\"compute_max_score\":0,\"compute_max_score_count\":0,\"count_weight\":0,\"count_weight_count\":0,\"create_weight\":0,\"create_weight_count\":0,\"match\":0,\"match_count\":0,\"next_doc\":0,\"next_doc_count\":0,\"score\":0,\"score_count\":0,\"set_min_competitive_score\":0,\"set_min_competitive_score_count\":0,\"shallow_advance\":0,\"shallow_advance_count\":0},\"description\":\"*:*\",\"time_in_nanos\":0,\"type\":\"MatchAllDocsQuery\"}],\"rewrite_time\":0}],\"shard_id\":0}]},\"timed_out\":false,\"took\":0}"}
