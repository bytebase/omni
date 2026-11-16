-- pin: default mongo:8.0.38 sha256:d97fea17ccb8d004cd137be8dab54d70d0df11be69c9f422aefe4484c8414ab7 r8.0.38

== TestQuery default#1
> raw
| {"op":"Execute","database":"omni_driver","statement":"db['test.collection'].insertMany([\n\t\t{\"name\": \"Alice\", \"age\": 30},\n\t\t{\"name\": \"Bob\", \"age\": 25},\n\t\t{\"name\": \"Charlie\", \"age\": 35}\n\t]);"}
< raw
| {"value":{"operation":"OpInsertMany","values":[{"acknowledged":true,"insertedIds":[{"$oid":"000000000000000000000001"},{"$oid":"000000000000000000000002"},{"$oid":"000000000000000000000003"}]}]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db[\"test.collection\"].find().limit(50)"}
< raw
| {"value":{"operation":"OpFind","values":[{"_id":{"$oid":"000000000000000000000001"},"name":"Alice","age":{"$numberInt":"30"}},{"_id":{"$oid":"000000000000000000000002"},"name":"Bob","age":{"$numberInt":"25"}},{"_id":{"$oid":"000000000000000000000003"},"name":"Charlie","age":{"$numberInt":"35"}}]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db['test.collection'].find().limit(50)"}
< raw
| {"value":{"operation":"OpFind","values":[{"_id":{"$oid":"000000000000000000000001"},"name":"Alice","age":{"$numberInt":"30"}},{"_id":{"$oid":"000000000000000000000002"},"name":"Bob","age":{"$numberInt":"25"}},{"_id":{"$oid":"000000000000000000000003"},"name":"Charlie","age":{"$numberInt":"35"}}]}}
> raw
| {"op":"Disconnect"}
< raw
| {}
