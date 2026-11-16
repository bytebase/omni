-- pin: default mongo:8.0.38 sha256:d97fea17ccb8d004cd137be8dab54d70d0df11be69c9f422aefe4484c8414ab7 r8.0.38

== TestDumpRestores default#1
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.createCollection(\"users\");"}
< raw
| {"value":{"operation":"OpCreateCollection","values":[{"ok":{"$numberInt":"1"}}]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.users.insertMany([\n  {_id: 1, name: \"alice\", age: 30, tags: [\"admin\", \"staff\"], profile: {city: \"lisbon\", zip: \"1000-001\"}},\n  {_id: 2, name: \"bob\", age: 17, tags: [], profile: {city: \"porto\"}},\n  {_id: 3, name: \"carol\", age: 42, contact: {email: \"carol@example.com\"}}\n]);"}
< raw
| {"value":{"operation":"OpInsertMany","values":[{"acknowledged":true,"insertedIds":[{"$numberInt":"1"},{"$numberInt":"2"},{"$numberInt":"3"}]}]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.users.createIndex({name: 1}, {unique: true});"}
< raw
| {"value":{"operation":"OpCreateIndex","values":["name_1"]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.users.createIndex({age: -1, name: 1});"}
< raw
| {"value":{"operation":"OpCreateIndex","values":["age_-1_name_1"]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.createCollection(\"orders\");"}
< raw
| {"value":{"operation":"OpCreateCollection","values":[{"ok":{"$numberInt":"1"}}]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.orders.insertOne({_id: 1, customer_id: 1, total: 100.5, created_at: 1767323045000});"}
< raw
| {"value":{"operation":"OpInsertOne","values":[{"acknowledged":true,"insertedId":{"$numberInt":"1"}}]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.orders.createIndex({customer_id: 1, created_at: -1});"}
< raw
| {"value":{"operation":"OpCreateIndex","values":["customer_id_1_created_at_-1"]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.createCollection(\"events\");"}
< raw
| {"value":{"operation":"OpCreateCollection","values":[{"ok":{"$numberInt":"1"}}]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.events.insertMany([{_id: \"a\", kind: \"login\"}, {_id: \"b\", kind: \"logout\"}]);"}
< raw
| {"value":{"operation":"OpInsertMany","values":[{"acknowledged":true,"insertedIds":["a","b"]}]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.events.createIndex({kind: 1, _id: -1});"}
< raw
| {"value":{"operation":"OpCreateIndex","values":["kind_1__id_-1"]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.createCollection(\"empty\");"}
< raw
| {"value":{"operation":"OpCreateCollection","values":[{"ok":{"$numberInt":"1"}}]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.createCollection(\"nested\");"}
< raw
| {"value":{"operation":"OpCreateCollection","values":[{"ok":{"$numberInt":"1"}}]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.nested.insertOne({_id: 1, doc: {a: [1, 2], b: \"x\"}, list: [{k: 1}, {k: 2}]});"}
< raw
| {"value":{"operation":"OpInsertOne","values":[{"acknowledged":true,"insertedId":{"$numberInt":"1"}}]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.runCommand({create: \"v_adults\", viewOn: \"users\", pipeline: [{$match: {age: {$gte: 18}}}]});"}
< raw
| {"value":{"operation":"OpRunCommand","values":[{"ok":{"$numberDouble":"1.0"}}]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.createCollection(\"sessions\");"}
< raw
| {"value":{"operation":"OpCreateCollection","values":[{"ok":{"$numberInt":"1"}}]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.sessions.createIndex({created_at: 1}, {expireAfterSeconds: 3600});"}
< raw
| {"value":{"operation":"OpCreateIndex","values":["created_at_1"]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.runCommand({create: \"telemetry\", timeseries: {timeField: \"at\", metaField: \"sensor\", granularity: \"seconds\"}});"}
< raw
| {"value":{"operation":"OpRunCommand","values":[{"ok":{"$numberDouble":"1.0"}}]}}
> raw
| {"op":"ListDatabaseNames"}
< raw
| {"value":{"names":["admin","config","local","omni_driver"]}}
> raw
| {"op":"ListCollections","database":"omni_driver"}
< raw
| {"value":{"docs":[{"idIndex":{"v":{"$numberInt":"2"},"key":{"_id":{"$numberInt":"1"}},"name":"_id_"},"info":{"readOnly":false,"uuid":{"$binary":{"base64":"AAAAAAAAAAAAAAAAAAAAAQ==","subType":"04"}}},"name":"empty","options":{},"type":"collection"},{"idIndex":{"v":{"$numberInt":"2"},"key":{"_id":{"$numberInt":"1"}},"name":"_id_"},"info":{"readOnly":false,"uuid":{"$binary":{"base64":"AAAAAAAAAAAAAAAAAAAAAg==","subType":"04"}}},"name":"events","options":{},"type":"collection"},{"idIndex":{"v":{"$numberInt":"2"},"key":{"_id":{"$numberInt":"1"}},"name":"_id_"},"info":{"readOnly":false,"uuid":{"$binary":{"base64":"AAAAAAAAAAAAAAAAAAAAAw==","subType":"04"}}},"name":"nested","options":{},"type":"collection"},{"idIndex":{"v":{"$numberInt":"2"},"key":{"_id":{"$numberInt":"1"}},"name":"_id_"},"info":{"readOnly":false,"uuid":{"$binary":{"base64":"AAAAAAAAAAAAAAAAAAAABA==","subType":"04"}}},"name":"orders","options":{},"type":"collection"},{"idIndex":{"v":{"$numberInt":"2"},"key":{"_id":{"$numberInt":"1"}},"name":"_id_"},"info":{"readOnly":false,"uuid":{"$binary":{"base64":"AAAAAAAAAAAAAAAAAAAABQ==","subType":"04"}}},"name":"sessions","options":{},"type":"collection"},{"info":{"readOnly":false,"uuid":{"$binary":{"base64":"AAAAAAAAAAAAAAAAAAAABg==","subType":"04"}}},"name":"system.buckets.telemetry","options":{"validator":{"$jsonSchema":{"bsonType":"object","required":["_id","control","data"],"properties":{"_id":{"bsonType":"objectId"},"control":{"bsonType":"object","required":["version","min","max"],"properties":{"version":{"bsonType":"number"},"min":{"bsonType":"object","required":["at"],"properties":{"at":{"bsonType":"date"}}},"max":{"bsonType":"object","required":["at"],"properties":{"at":{"bsonType":"date"}}},"closed":{"bsonType":"bool"},"count":{"bsonType":"number","minimum":{"$numberInt":"1"}}},"additionalProperties":false},"data":{"bsonType":"object"},"meta":{}},"additionalProperties":false}},"clusteredIndex":true,"timeseries":{"timeField":"at","metaField":"sensor","granularity":"seconds","bucketMaxSpanSeconds":{"$numberInt":"3600"}}},"type":"collection"},{"idIndex":{"v":{"$numberInt":"2"},"key":{"_id":{"$numberInt":"1"}},"name":"_id_"},"info":{"readOnly":false,"uuid":{"$binary":{"base64":"AAAAAAAAAAAAAAAAAAAABw==","subType":"04"}}},"name":"system.views","options":{},"type":"collection"},{"info":{"readOnly":false},"name":"telemetry","options":{"timeseries":{"timeField":"at","metaField":"sensor","granularity":"seconds","bucketMaxSpanSeconds":{"$numberInt":"3600"}}},"type":"timeseries"},{"idIndex":{"v":{"$numberInt":"2"},"key":{"_id":{"$numberInt":"1"}},"name":"_id_"},"info":{"readOnly":false,"uuid":{"$binary":{"base64":"AAAAAAAAAAAAAAAAAAAACA==","subType":"04"}}},"name":"users","options":{},"type":"collection"},{"info":{"readOnly":true},"name":"v_adults","options":{"viewOn":"users","pipeline":[{"$match":{"age":{"$gte":{"$numberInt":"18"}}}}]},"type":"view"}]}}
> raw
| {"op":"EstimatedDocumentCount","database":"omni_driver","collection":"empty"}
< raw
| {"value":{"count":0}}
> raw
| {"op":"CollStats","database":"omni_driver","collection":"empty"}
< raw
| {"value":{"size":{"$numberInt":"0"},"totalIndexSize":{"$numberInt":"4096"}}}
> raw
| {"op":"ListIndexes","database":"omni_driver","collection":"empty"}
< raw
| {"value":{"docs":[{"key":{"_id":{"$numberInt":"1"}},"name":"_id_","v":{"$numberInt":"2"}}]}}
> raw
| {"op":"EstimatedDocumentCount","database":"omni_driver","collection":"events"}
< raw
| {"value":{"count":2}}
> raw
| {"op":"CollStats","database":"omni_driver","collection":"events"}
< raw
| {"value":{"size":{"$numberInt":"65"},"totalIndexSize":{"$numberInt":"24576"}}}
> raw
| {"op":"ListIndexes","database":"omni_driver","collection":"events"}
< raw
| {"value":{"docs":[{"key":{"_id":{"$numberInt":"1"}},"name":"_id_","v":{"$numberInt":"2"}},{"key":{"kind":{"$numberInt":"1"},"_id":{"$numberInt":"-1"}},"name":"kind_1__id_-1","v":{"$numberInt":"2"}}]}}
> raw
| {"op":"EstimatedDocumentCount","database":"omni_driver","collection":"nested"}
< raw
| {"value":{"count":1}}
> raw
| {"op":"CollStats","database":"omni_driver","collection":"nested"}
< raw
| {"value":{"size":{"$numberInt":"96"},"totalIndexSize":{"$numberInt":"4096"}}}
> raw
| {"op":"ListIndexes","database":"omni_driver","collection":"nested"}
< raw
| {"value":{"docs":[{"key":{"_id":{"$numberInt":"1"}},"name":"_id_","v":{"$numberInt":"2"}}]}}
> raw
| {"op":"EstimatedDocumentCount","database":"omni_driver","collection":"orders"}
< raw
| {"value":{"count":1}}
> raw
| {"op":"CollStats","database":"omni_driver","collection":"orders"}
< raw
| {"value":{"size":{"$numberInt":"66"},"totalIndexSize":{"$numberInt":"24576"}}}
> raw
| {"op":"ListIndexes","database":"omni_driver","collection":"orders"}
< raw
| {"value":{"docs":[{"key":{"_id":{"$numberInt":"1"}},"name":"_id_","v":{"$numberInt":"2"}},{"key":{"customer_id":{"$numberInt":"1"},"created_at":{"$numberInt":"-1"}},"name":"customer_id_1_created_at_-1","v":{"$numberInt":"2"}}]}}
> raw
| {"op":"EstimatedDocumentCount","database":"omni_driver","collection":"sessions"}
< raw
| {"value":{"count":0}}
> raw
| {"op":"CollStats","database":"omni_driver","collection":"sessions"}
< raw
| {"value":{"size":{"$numberInt":"0"},"totalIndexSize":{"$numberInt":"8192"}}}
> raw
| {"op":"ListIndexes","database":"omni_driver","collection":"sessions"}
< raw
| {"value":{"docs":[{"key":{"_id":{"$numberInt":"1"}},"name":"_id_","v":{"$numberInt":"2"}},{"expireAfterSeconds":{"$numberInt":"3600"},"key":{"created_at":{"$numberInt":"1"}},"name":"created_at_1","v":{"$numberInt":"2"}}]}}
> raw
| {"op":"EstimatedDocumentCount","database":"omni_driver","collection":"telemetry"}
< raw
| {"value":{"count":0}}
> raw
| {"op":"CollStats","database":"omni_driver","collection":"telemetry"}
< raw
| {"value":{"size":{"$numberInt":"0"},"totalIndexSize":{"$numberInt":"4096"}}}
> raw
| {"op":"ListIndexes","database":"omni_driver","collection":"telemetry"}
< raw
| {"value":{"docs":[{"key":{"sensor":{"$numberInt":"1"},"at":{"$numberInt":"1"}},"name":"sensor_1_at_1","v":{"$numberInt":"2"}}]}}
> raw
| {"op":"EstimatedDocumentCount","database":"omni_driver","collection":"users"}
< raw
| {"value":{"count":3}}
> raw
| {"op":"CollStats","database":"omni_driver","collection":"users"}
< raw
| {"value":{"size":{"$numberInt":"285"},"totalIndexSize":{"$numberInt":"45056"}}}
> raw
| {"op":"ListIndexes","database":"omni_driver","collection":"users"}
< raw
| {"value":{"docs":[{"key":{"_id":{"$numberInt":"1"}},"name":"_id_","v":{"$numberInt":"2"}},{"key":{"age":{"$numberInt":"-1"},"name":{"$numberInt":"1"}},"name":"age_-1_name_1","v":{"$numberInt":"2"}},{"key":{"name":{"$numberInt":"1"}},"name":"name_1","unique":true,"v":{"$numberInt":"2"}}]}}
> raw
| {"op":"Execute","database":"omni_restored","statement":"db.createCollection(\"empty\");"}
< raw
| {"value":{"operation":"OpCreateCollection","values":[{"ok":{"$numberInt":"1"}}]}}
> raw
| {"op":"Execute","database":"omni_restored","statement":"db.runCommand({\"createIndexes\": \"empty\", \"indexes\": [{\"key\": {\"_id\":1}, \"name\": \"_id_\", \"v\": 2}]});"}
< raw
| {"value":{"operation":"OpRunCommand","values":[{"numIndexesBefore":{"$numberInt":"1"},"numIndexesAfter":{"$numberInt":"1"},"note":"all indexes already exist","ok":{"$numberDouble":"1.0"}}]}}
> raw
| {"op":"Execute","database":"omni_restored","statement":"db.createCollection(\"events\");"}
< raw
| {"value":{"operation":"OpCreateCollection","values":[{"ok":{"$numberInt":"1"}}]}}
> raw
| {"op":"Execute","database":"omni_restored","statement":"db.runCommand({\"createIndexes\": \"events\", \"indexes\": [{\"key\": {\"_id\":1}, \"name\": \"_id_\", \"v\": 2}]});"}
< raw
| {"value":{"operation":"OpRunCommand","values":[{"numIndexesBefore":{"$numberInt":"1"},"numIndexesAfter":{"$numberInt":"1"},"note":"all indexes already exist","ok":{"$numberDouble":"1.0"}}]}}
> raw
| {"op":"Execute","database":"omni_restored","statement":"db.runCommand({\"createIndexes\": \"events\", \"indexes\": [{\"key\": {\"kind\":1,\"_id\":-1}, \"name\": \"kind_1__id_-1\", \"v\": 2}]});"}
< raw
| {"value":{"operation":"OpRunCommand","values":[{"numIndexesBefore":{"$numberInt":"1"},"numIndexesAfter":{"$numberInt":"2"},"createdCollectionAutomatically":false,"ok":{"$numberDouble":"1.0"}}]}}
> raw
| {"op":"Execute","database":"omni_restored","statement":"db.createCollection(\"nested\");"}
< raw
| {"value":{"operation":"OpCreateCollection","values":[{"ok":{"$numberInt":"1"}}]}}
> raw
| {"op":"Execute","database":"omni_restored","statement":"db.runCommand({\"createIndexes\": \"nested\", \"indexes\": [{\"key\": {\"_id\":1}, \"name\": \"_id_\", \"v\": 2}]});"}
< raw
| {"value":{"operation":"OpRunCommand","values":[{"numIndexesBefore":{"$numberInt":"1"},"numIndexesAfter":{"$numberInt":"1"},"note":"all indexes already exist","ok":{"$numberDouble":"1.0"}}]}}
> raw
| {"op":"Execute","database":"omni_restored","statement":"db.createCollection(\"orders\");"}
< raw
| {"value":{"operation":"OpCreateCollection","values":[{"ok":{"$numberInt":"1"}}]}}
> raw
| {"op":"Execute","database":"omni_restored","statement":"db.runCommand({\"createIndexes\": \"orders\", \"indexes\": [{\"key\": {\"_id\":1}, \"name\": \"_id_\", \"v\": 2}]});"}
< raw
| {"value":{"operation":"OpRunCommand","values":[{"numIndexesBefore":{"$numberInt":"1"},"numIndexesAfter":{"$numberInt":"1"},"note":"all indexes already exist","ok":{"$numberDouble":"1.0"}}]}}
> raw
| {"op":"Execute","database":"omni_restored","statement":"db.runCommand({\"createIndexes\": \"orders\", \"indexes\": [{\"key\": {\"customer_id\":1,\"created_at\":-1}, \"name\": \"customer_id_1_created_at_-1\", \"v\": 2}]});"}
< raw
| {"value":{"operation":"OpRunCommand","values":[{"numIndexesBefore":{"$numberInt":"1"},"numIndexesAfter":{"$numberInt":"2"},"createdCollectionAutomatically":false,"ok":{"$numberDouble":"1.0"}}]}}
> raw
| {"op":"Execute","database":"omni_restored","statement":"db.createCollection(\"sessions\");"}
< raw
| {"value":{"operation":"OpCreateCollection","values":[{"ok":{"$numberInt":"1"}}]}}
> raw
| {"op":"Execute","database":"omni_restored","statement":"db.runCommand({\"createIndexes\": \"sessions\", \"indexes\": [{\"key\": {\"_id\":1}, \"name\": \"_id_\", \"v\": 2}]});"}
< raw
| {"value":{"operation":"OpRunCommand","values":[{"numIndexesBefore":{"$numberInt":"1"},"numIndexesAfter":{"$numberInt":"1"},"note":"all indexes already exist","ok":{"$numberDouble":"1.0"}}]}}
> raw
| {"op":"Execute","database":"omni_restored","statement":"db.runCommand({\"createIndexes\": \"sessions\", \"indexes\": [{\"key\": {\"created_at\":1}, \"name\": \"created_at_1\", \"expireAfterSeconds\": 3600, \"v\": 2}]});"}
< raw
| {"value":{"operation":"OpRunCommand","values":[{"numIndexesBefore":{"$numberInt":"1"},"numIndexesAfter":{"$numberInt":"2"},"createdCollectionAutomatically":false,"ok":{"$numberDouble":"1.0"}}]}}
> raw
| {"op":"Execute","database":"omni_restored","statement":"db.runCommand({\"create\": \"telemetry\", \"timeseries\": {\"timeField\": \"at\", \"metaField\": \"sensor\", \"granularity\": \"seconds\", \"bucketMaxSpanSeconds\": 3600}});"}
< raw
| {"value":{"operation":"OpRunCommand","values":[{"ok":{"$numberDouble":"1.0"}}]}}
> raw
| {"op":"Execute","database":"omni_restored","statement":"db.runCommand({\"createIndexes\": \"telemetry\", \"indexes\": [{\"key\": {\"sensor\":1,\"at\":1}, \"name\": \"sensor_1_at_1\", \"v\": 2}]});"}
< raw
| {"value":{"operation":"OpRunCommand","values":[{"numIndexesBefore":{"$numberInt":"1"},"numIndexesAfter":{"$numberInt":"1"},"note":"all indexes already exist","ok":{"$numberDouble":"1.0"}}]}}
> raw
| {"op":"Execute","database":"omni_restored","statement":"db.createCollection(\"users\");"}
< raw
| {"value":{"operation":"OpCreateCollection","values":[{"ok":{"$numberInt":"1"}}]}}
> raw
| {"op":"Execute","database":"omni_restored","statement":"db.runCommand({\"createIndexes\": \"users\", \"indexes\": [{\"key\": {\"_id\":1}, \"name\": \"_id_\", \"v\": 2}]});"}
< raw
| {"value":{"operation":"OpRunCommand","values":[{"numIndexesBefore":{"$numberInt":"1"},"numIndexesAfter":{"$numberInt":"1"},"note":"all indexes already exist","ok":{"$numberDouble":"1.0"}}]}}
> raw
| {"op":"Execute","database":"omni_restored","statement":"db.runCommand({\"createIndexes\": \"users\", \"indexes\": [{\"key\": {\"age\":-1,\"name\":1}, \"name\": \"age_-1_name_1\", \"v\": 2}]});"}
< raw
| {"value":{"operation":"OpRunCommand","values":[{"numIndexesBefore":{"$numberInt":"1"},"numIndexesAfter":{"$numberInt":"2"},"createdCollectionAutomatically":false,"ok":{"$numberDouble":"1.0"}}]}}
> raw
| {"op":"Execute","database":"omni_restored","statement":"db.runCommand({\"createIndexes\": \"users\", \"indexes\": [{\"key\": {\"name\":1}, \"name\": \"name_1\", \"unique\": true, \"v\": 2}]});"}
< raw
| {"value":{"operation":"OpRunCommand","values":[{"numIndexesBefore":{"$numberInt":"2"},"numIndexesAfter":{"$numberInt":"3"},"createdCollectionAutomatically":false,"ok":{"$numberDouble":"1.0"}}]}}
> raw
| {"op":"Execute","database":"omni_restored","statement":"db.runCommand({\"create\": \"v_adults\", \"viewOn\": \"users\", \"pipeline\": [{\"$match\": {\"age\": {\"$gte\": 18}}}]});"}
< raw
| {"value":{"operation":"OpRunCommand","values":[{"ok":{"$numberDouble":"1.0"}}]}}
> raw
| {"op":"ListDatabaseNames"}
< raw
| {"value":{"names":["admin","config","local","omni_driver","omni_restored"]}}
> raw
| {"op":"ListCollections","database":"omni_restored"}
< raw
| {"value":{"docs":[{"idIndex":{"v":{"$numberInt":"2"},"key":{"_id":{"$numberInt":"1"}},"name":"_id_"},"info":{"readOnly":false,"uuid":{"$binary":{"base64":"AAAAAAAAAAAAAAAAAAAACQ==","subType":"04"}}},"name":"empty","options":{},"type":"collection"},{"idIndex":{"v":{"$numberInt":"2"},"key":{"_id":{"$numberInt":"1"}},"name":"_id_"},"info":{"readOnly":false,"uuid":{"$binary":{"base64":"AAAAAAAAAAAAAAAAAAAACg==","subType":"04"}}},"name":"events","options":{},"type":"collection"},{"idIndex":{"v":{"$numberInt":"2"},"key":{"_id":{"$numberInt":"1"}},"name":"_id_"},"info":{"readOnly":false,"uuid":{"$binary":{"base64":"AAAAAAAAAAAAAAAAAAAACw==","subType":"04"}}},"name":"nested","options":{},"type":"collection"},{"idIndex":{"v":{"$numberInt":"2"},"key":{"_id":{"$numberInt":"1"}},"name":"_id_"},"info":{"readOnly":false,"uuid":{"$binary":{"base64":"AAAAAAAAAAAAAAAAAAAADA==","subType":"04"}}},"name":"orders","options":{},"type":"collection"},{"idIndex":{"v":{"$numberInt":"2"},"key":{"_id":{"$numberInt":"1"}},"name":"_id_"},"info":{"readOnly":false,"uuid":{"$binary":{"base64":"AAAAAAAAAAAAAAAAAAAADQ==","subType":"04"}}},"name":"sessions","options":{},"type":"collection"},{"info":{"readOnly":false,"uuid":{"$binary":{"base64":"AAAAAAAAAAAAAAAAAAAADg==","subType":"04"}}},"name":"system.buckets.telemetry","options":{"validator":{"$jsonSchema":{"bsonType":"object","required":["_id","control","data"],"properties":{"_id":{"bsonType":"objectId"},"control":{"bsonType":"object","required":["version","min","max"],"properties":{"version":{"bsonType":"number"},"min":{"bsonType":"object","required":["at"],"properties":{"at":{"bsonType":"date"}}},"max":{"bsonType":"object","required":["at"],"properties":{"at":{"bsonType":"date"}}},"closed":{"bsonType":"bool"},"count":{"bsonType":"number","minimum":{"$numberInt":"1"}}},"additionalProperties":false},"data":{"bsonType":"object"},"meta":{}},"additionalProperties":false}},"clusteredIndex":true,"timeseries":{"timeField":"at","metaField":"sensor","granularity":"seconds","bucketMaxSpanSeconds":{"$numberInt":"3600"}}},"type":"collection"},{"idIndex":{"v":{"$numberInt":"2"},"key":{"_id":{"$numberInt":"1"}},"name":"_id_"},"info":{"readOnly":false,"uuid":{"$binary":{"base64":"AAAAAAAAAAAAAAAAAAAADw==","subType":"04"}}},"name":"system.views","options":{},"type":"collection"},{"info":{"readOnly":false},"name":"telemetry","options":{"timeseries":{"timeField":"at","metaField":"sensor","granularity":"seconds","bucketMaxSpanSeconds":{"$numberInt":"3600"}}},"type":"timeseries"},{"idIndex":{"v":{"$numberInt":"2"},"key":{"_id":{"$numberInt":"1"}},"name":"_id_"},"info":{"readOnly":false,"uuid":{"$binary":{"base64":"AAAAAAAAAAAAAAAAAAAAEA==","subType":"04"}}},"name":"users","options":{},"type":"collection"},{"info":{"readOnly":true},"name":"v_adults","options":{"viewOn":"users","pipeline":[{"$match":{"age":{"$gte":{"$numberInt":"18"}}}}]},"type":"view"}]}}
> raw
| {"op":"EstimatedDocumentCount","database":"omni_restored","collection":"empty"}
< raw
| {"value":{"count":0}}
> raw
| {"op":"CollStats","database":"omni_restored","collection":"empty"}
< raw
| {"value":{"size":{"$numberInt":"0"},"totalIndexSize":{"$numberInt":"4096"}}}
> raw
| {"op":"ListIndexes","database":"omni_restored","collection":"empty"}
< raw
| {"value":{"docs":[{"key":{"_id":{"$numberInt":"1"}},"name":"_id_","v":{"$numberInt":"2"}}]}}
> raw
| {"op":"EstimatedDocumentCount","database":"omni_restored","collection":"events"}
< raw
| {"value":{"count":0}}
> raw
| {"op":"CollStats","database":"omni_restored","collection":"events"}
< raw
| {"value":{"size":{"$numberInt":"0"},"totalIndexSize":{"$numberInt":"8192"}}}
> raw
| {"op":"ListIndexes","database":"omni_restored","collection":"events"}
< raw
| {"value":{"docs":[{"key":{"_id":{"$numberInt":"1"}},"name":"_id_","v":{"$numberInt":"2"}},{"key":{"kind":{"$numberInt":"1"},"_id":{"$numberInt":"-1"}},"name":"kind_1__id_-1","v":{"$numberInt":"2"}}]}}
> raw
| {"op":"EstimatedDocumentCount","database":"omni_restored","collection":"nested"}
< raw
| {"value":{"count":0}}
> raw
| {"op":"CollStats","database":"omni_restored","collection":"nested"}
< raw
| {"value":{"size":{"$numberInt":"0"},"totalIndexSize":{"$numberInt":"4096"}}}
> raw
| {"op":"ListIndexes","database":"omni_restored","collection":"nested"}
< raw
| {"value":{"docs":[{"key":{"_id":{"$numberInt":"1"}},"name":"_id_","v":{"$numberInt":"2"}}]}}
> raw
| {"op":"EstimatedDocumentCount","database":"omni_restored","collection":"orders"}
< raw
| {"value":{"count":0}}
> raw
| {"op":"CollStats","database":"omni_restored","collection":"orders"}
< raw
| {"value":{"size":{"$numberInt":"0"},"totalIndexSize":{"$numberInt":"8192"}}}
> raw
| {"op":"ListIndexes","database":"omni_restored","collection":"orders"}
< raw
| {"value":{"docs":[{"key":{"_id":{"$numberInt":"1"}},"name":"_id_","v":{"$numberInt":"2"}},{"key":{"customer_id":{"$numberInt":"1"},"created_at":{"$numberInt":"-1"}},"name":"customer_id_1_created_at_-1","v":{"$numberInt":"2"}}]}}
> raw
| {"op":"EstimatedDocumentCount","database":"omni_restored","collection":"sessions"}
< raw
| {"value":{"count":0}}
> raw
| {"op":"CollStats","database":"omni_restored","collection":"sessions"}
< raw
| {"value":{"size":{"$numberInt":"0"},"totalIndexSize":{"$numberInt":"8192"}}}
> raw
| {"op":"ListIndexes","database":"omni_restored","collection":"sessions"}
< raw
| {"value":{"docs":[{"key":{"_id":{"$numberInt":"1"}},"name":"_id_","v":{"$numberInt":"2"}},{"expireAfterSeconds":{"$numberInt":"3600"},"key":{"created_at":{"$numberInt":"1"}},"name":"created_at_1","v":{"$numberInt":"2"}}]}}
> raw
| {"op":"EstimatedDocumentCount","database":"omni_restored","collection":"telemetry"}
< raw
| {"value":{"count":0}}
> raw
| {"op":"CollStats","database":"omni_restored","collection":"telemetry"}
< raw
| {"value":{"size":{"$numberInt":"0"},"totalIndexSize":{"$numberInt":"4096"}}}
> raw
| {"op":"ListIndexes","database":"omni_restored","collection":"telemetry"}
< raw
| {"value":{"docs":[{"key":{"sensor":{"$numberInt":"1"},"at":{"$numberInt":"1"}},"name":"sensor_1_at_1","v":{"$numberInt":"2"}}]}}
> raw
| {"op":"EstimatedDocumentCount","database":"omni_restored","collection":"users"}
< raw
| {"value":{"count":0}}
> raw
| {"op":"CollStats","database":"omni_restored","collection":"users"}
< raw
| {"value":{"size":{"$numberInt":"0"},"totalIndexSize":{"$numberInt":"12288"}}}
> raw
| {"op":"ListIndexes","database":"omni_restored","collection":"users"}
< raw
| {"value":{"docs":[{"key":{"_id":{"$numberInt":"1"}},"name":"_id_","v":{"$numberInt":"2"}},{"key":{"age":{"$numberInt":"-1"},"name":{"$numberInt":"1"}},"name":"age_-1_name_1","v":{"$numberInt":"2"}},{"key":{"name":{"$numberInt":"1"}},"name":"name_1","unique":true,"v":{"$numberInt":"2"}}]}}
> raw
| {"op":"Disconnect"}
< raw
| {}
