PUT /users
{"mappings":{"properties":{"name":{"type":"text","fields":{"keyword":{"type":"keyword"}}},"email":{"type":"keyword"},"ssn":{"type":"keyword"},"age":{"type":"integer"},"address":{"properties":{"city":{"type":"keyword"},"zip":{"type":"keyword"}}},"tags":{"type":"keyword"}}}}

PUT /users_archive
{"mappings":{"properties":{"name":{"type":"keyword"},"email":{"type":"keyword"},"archived_at":{"type":"date"}}}}

PUT /orders
{"mappings":{"properties":{"order_id":{"type":"keyword"},"user_email":{"type":"keyword"},"total":{"type":"double"},"items":{"type":"nested","properties":{"sku":{"type":"keyword"},"qty":{"type":"integer"}}}}}}

PUT /hidden_audit
{"settings":{"index":{"hidden":true}},"mappings":{"properties":{"actor":{"type":"keyword"},"action":{"type":"keyword"}}}}

PUT /.customers
{"mappings":{"properties":{"vip":{"type":"keyword"}}}}

PUT /_index_template/omni-events
{"index_patterns":["events-*"],"data_stream":{},"priority":500,"template":{"mappings":{"properties":{"@timestamp":{"type":"date"},"message":{"type":"text"},"user":{"type":"keyword"}}}}}

PUT /_data_stream/events-app

POST /_aliases
{"actions":[{"add":{"index":"users","alias":"people"}},{"add":{"index":"users","alias":"everyone"}},{"add":{"index":"users_archive","alias":"everyone"}},{"add":{"index":"users","alias":"adults","filter":{"range":{"age":{"gte":18}}}}}]}

POST /_bulk?refresh=true
{"index":{"_index":"users","_id":"1"}}
{"name":"Alice","email":"alice@example.com","ssn":"111-11-1111","age":34,"address":{"city":"Paris","zip":"75001"},"tags":["vip","eu"]}
{"index":{"_index":"users","_id":"2"}}
{"name":"Bob","email":"bob@example.com","ssn":"222-22-2222","age":12,"address":{"city":"Lyon","zip":"69001"},"tags":["new"]}
{"index":{"_index":"users_archive","_id":"1"}}
{"name":"Carol","email":"carol@example.com","archived_at":"2020-01-01"}
{"index":{"_index":"orders","_id":"o1"}}
{"order_id":"o1","user_email":"alice@example.com","total":12.5,"items":[{"sku":"a","qty":1},{"sku":"b","qty":2}]}
{"index":{"_index":"orders","_id":"o2"}}
{"order_id":"o2","user_email":"bob@example.com","total":3,"items":[{"sku":"c","qty":3}]}
{"index":{"_index":"hidden_audit","_id":"h1"}}
{"actor":"alice","action":"login"}
{"index":{"_index":".customers","_id":"c1"}}
{"vip":"alice"}
{"create":{"_index":"events-app","_id":"e1"}}
{"@timestamp":"2026-11-25T00:00:00Z","message":"login","user":"alice"}
