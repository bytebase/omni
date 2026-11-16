-- pin: default mongo:8.0.38 sha256:d97fea17ccb8d004cd137be8dab54d70d0df11be69c9f422aefe4484c8414ab7 r8.0.38

== TestQueryDistinctRendersEveryValueAsAJSONObject default#1
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.readings.insertMany([\n\t\t{_id: 1, value: 30},\n\t\t{_id: 2, value: 30},\n\t\t{_id: 3, value: 1.5},\n\t\t{_id: 4, value: null}\n\t]);"}
< raw
| {"value":{"operation":"OpInsertMany","values":[{"acknowledged":true,"insertedIds":[{"$numberInt":"1"},{"$numberInt":"2"},{"$numberInt":"3"},{"$numberInt":"4"}]}]}}
> raw
| {"op":"Execute","database":"omni_driver","statement":"db.readings.distinct(\"value\")"}
< raw
| {"value":{"operation":"OpDistinct","values":[null,{"$numberDouble":"1.5"},{"$numberInt":"30"}]}}
> raw
| {"op":"Disconnect"}
< raw
| {}
