-- pin: default mongo:8.0.38 sha256:d97fea17ccb8d004cd137be8dab54d70d0df11be69c9f422aefe4484c8414ab7 r8.0.38

== TestQuerySpanFidelity default#1
> raw
| {"op":"SpanFixture","statement":"// The starting databases of span.tsv: mongosh statements, one per line, each\n// run in the database the last \"// database:\" line names. A view is made with\n// the create command, since gomongo, which runs these lines as Bytebase runs a\n// statement, has no db.createView.\n// database: hr\ndb.employees.insertOne({_id: 1, name: \"ann\", salary: 100, dept: 1, boss: null})\ndb.departments.insertOne({_id: 1, title: \"eng\"})\ndb.runCommand({create: \"staff\", viewOn: \"employees\", pipeline: [{$project: {salary: 0}}]})\ndb.runCommand({create: \"staff_depts\", viewOn: \"staff\", pipeline: [{$lookup: {from: \"departments\", localField: \"dept\", foreignField: \"_id\", as: \"dept\"}}]})\n// database: shop\ndb.users.insertOne({_id: 1, name: \"a\", email: \"a@x\", ssn: \"1\", manager: null})\ndb.orders.insertOne({_id: 1, uid: 1, total: 5, sku: \"s1\"})\ndb.products.insertOne({_id: \"s1\", name: \"p\"})\ndb.reviews.insertOne({_id: 1, sku: \"s1\", body: \"ok\"})\ndb.archive.insertOne({_id: 2, name: \"old\"})\ndb.createCollection(\"empty\")\ndb.getCollection(\"Mixed Case\").insertOne({_id: 1})\ndb.getCollection(\"audit.log\").insertOne({_id: 1, uid: 1})\ndb.runCommand({create: \"active_users\", viewOn: \"users\", pipeline: [{$match: {name: {$ne: null}}}]})\ndb.runCommand({create: \"user_orders\", viewOn: \"users\", pipeline: [{$lookup: {from: \"orders\", localField: \"_id\", foreignField: \"uid\", as: \"orders\"}}]})\ndb.runCommand({create: \"active_user_orders\", viewOn: \"active_users\", pipeline: [{$lookup: {from: \"orders\", localField: \"_id\", foreignField: \"uid\", as: \"orders\"}}]})\ndb.runCommand({create: \"all_people\", viewOn: \"users\", pipeline: [{$unionWith: \"archive\"}]})\ndb.runCommand({create: \"order_products\", viewOn: \"orders\", pipeline: [{$lookup: {from: \"products\", let: {s: \"$sku\"}, pipeline: [{$match: {$expr: {$eq: [\"$_id\", \"$$s\"]}}}, {$lookup: {from: \"reviews\", localField: \"_id\", foreignField: \"sku\", as: \"reviews\"}}], as: \"product\"}}]})\ndb.runCommand({create: \"user_graph\", viewOn: \"users\", pipeline: [{$graphLookup: {from: \"user_orders\", startWith: \"$_id\", connectFromField: \"_id\", connectToField: \"_id\", as: \"graph\"}}]})\ndb.runCommand({create: \"order_facets\", viewOn: \"orders\", pipeline: [{$facet: {buyers: [{$lookup: {from: \"active_users\", localField: \"uid\", foreignField: \"_id\", as: \"buyer\"}}]}}]})\ndb.runCommand({create: \"Mixed View\", viewOn: \"Mixed Case\", pipeline: []})"}
< raw
| {"views":[{"database":"hr","name":"staff","definition":"{\"viewOn\":\"employees\",\"pipeline\":[{\"$project\":{\"salary\":{\"$numberInt\":\"0\"}}}]}"},{"database":"hr","name":"staff_depts","definition":"{\"viewOn\":\"staff\",\"pipeline\":[{\"$lookup\":{\"from\":\"departments\",\"localField\":\"dept\",\"foreignField\":\"_id\",\"as\":\"dept\"}}]}"},{"database":"shop","name":"Mixed View","definition":"{\"viewOn\":\"Mixed Case\",\"pipeline\":[]}"},{"database":"shop","name":"active_user_orders","definition":"{\"viewOn\":\"active_users\",\"pipeline\":[{\"$lookup\":{\"from\":\"orders\",\"localField\":\"_id\",\"foreignField\":\"uid\",\"as\":\"orders\"}}]}"},{"database":"shop","name":"active_users","definition":"{\"viewOn\":\"users\",\"pipeline\":[{\"$match\":{\"name\":{\"$ne\":null}}}]}"},{"database":"shop","name":"all_people","definition":"{\"viewOn\":\"users\",\"pipeline\":[{\"$unionWith\":\"archive\"}]}"},{"database":"shop","name":"order_facets","definition":"{\"viewOn\":\"orders\",\"pipeline\":[{\"$facet\":{\"buyers\":[{\"$lookup\":{\"from\":\"active_users\",\"localField\":\"uid\",\"foreignField\":\"_id\",\"as\":\"buyer\"}}]}}]}"},{"database":"shop","name":"order_products","definition":"{\"viewOn\":\"orders\",\"pipeline\":[{\"$lookup\":{\"from\":\"products\",\"let\":{\"s\":\"$sku\"},\"pipeline\":[{\"$match\":{\"$expr\":{\"$eq\":[\"$_id\",\"$$s\"]}}},{\"$lookup\":{\"from\":\"reviews\",\"localField\":\"_id\",\"foreignField\":\"sku\",\"as\":\"reviews\"}}],\"as\":\"product\"}}]}"},{"database":"shop","name":"user_graph","definition":"{\"viewOn\":\"users\",\"pipeline\":[{\"$graphLookup\":{\"from\":\"user_orders\",\"startWith\":\"$_id\",\"connectFromField\":\"_id\",\"connectToField\":\"_id\",\"as\":\"graph\"}}]}"},{"database":"shop","name":"user_orders","definition":"{\"viewOn\":\"users\",\"pipeline\":[{\"$lookup\":{\"from\":\"orders\",\"localField\":\"_id\",\"foreignField\":\"uid\",\"as\":\"orders\"}}]}"}]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.find()"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.find({name: \"a\"})"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.find({}, {email: 1, _id: 0})"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.find({ssn: \"1\"}).sort({name: 1}).limit(5)"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.find().skip(1).limit(1)"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.find({$expr: {$eq: [\"$name\", \"a\"]}})"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.find({manager: null}, {name: 1}).sort({_id: -1})"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.find().explain()"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.find().explain(\"executionStats\")"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.findOne({_id: 1})"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.findOne()"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.findOne({}, {email: 1})"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.countDocuments({})"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.countDocuments({name: \"a\"})"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.estimatedDocumentCount()"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.count()"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.distinct(\"name\")"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.distinct(\"email\", {name: \"a\"})"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.getCollection(\"users\").find()"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db[\"users\"].find()"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.getCollection(\"Mixed Case\").find()"}
< raw
| {"tables":["shop.Mixed Case"]}
> raw
| {"op":"Span","database":"shop","statement":"db.getCollection(\"audit.log\").find()"}
< raw
| {"tables":["shop.audit.log"]}
> raw
| {"op":"Span","database":"shop","statement":"db.getCollection(\"audit.log\").aggregate([{$lookup: {from: \"users\", localField: \"uid\", foreignField: \"_id\", as: \"user\"}}])"}
< raw
| {"tables":["shop.audit.log","shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.empty.find()"}
< raw
| {"tables":["shop.empty"]}
> raw
| {"op":"Span","database":"shop","statement":"db.nosuch.find()"}
< raw
| {"tables":["shop.nosuch"]}
> raw
| {"op":"Span","database":"shop","statement":"db.nosuch.aggregate([{$lookup: {from: \"nosuch2\", localField: \"a\", foreignField: \"b\", as: \"c\"}}])"}
< raw
| {"tables":["shop.nosuch"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([])"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$match: {name: \"a\"}}])"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$project: {name: 1}}])"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$group: {_id: \"$name\", n: {$sum: 1}}}])"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$sort: {name: 1}}, {$skip: 1}, {$limit: 2}])"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$addFields: {x: 1}}, {$set: {y: 2}}, {$unset: \"ssn\"}])"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$sortByCount: \"$name\"}])"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$count: \"n\"}])"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$replaceRoot: {newRoot: {n: \"$name\"}}}])"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$setWindowFields: {sortBy: {_id: 1}, output: {r: {$rank: {}}}}}])"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$lookup: {from: \"orders\", localField: \"_id\", foreignField: \"uid\", as: \"orders\"}}])"}
< raw
| {"tables":["shop.orders","shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$lookup: {from: \"orders\", localField: \"_id\", foreignField: \"uid\", as: \"orders\"}}, {$unwind: \"$orders\"}, {$project: {name: 1, \"orders.total\": 1}}])"}
< raw
| {"tables":["shop.orders","shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$lookup: {from: \"orders\", let: {u: \"$_id\"}, pipeline: [{$match: {$expr: {$eq: [\"$uid\", \"$$u\"]}}}], as: \"orders\"}}])"}
< raw
| {"tables":["shop.orders","shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$lookup: {from: \"orders\", pipeline: [{$match: {total: {$gt: 1}}}], as: \"big\"}}])"}
< raw
| {"tables":["shop.orders","shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$lookup: {from: \"orders\", localField: \"_id\", foreignField: \"uid\", pipeline: [{$project: {total: 1}}], as: \"orders\"}}])"}
< raw
| {"tables":["shop.orders","shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$lookup: {from: \"orders\", localField: \"_id\", foreignField: \"uid\", as: \"orders\"}}, {$lookup: {from: \"archive\", localField: \"_id\", foreignField: \"_id\", as: \"old\"}}])"}
< raw
| {"tables":["shop.archive","shop.orders","shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$lookup: {from: \"orders\", let: {u: \"$_id\"}, pipeline: [{$match: {$expr: {$eq: [\"$uid\", \"$$u\"]}}}, {$lookup: {from: \"products\", localField: \"sku\", foreignField: \"_id\", as: \"product\"}}], as: \"orders\"}}])"}
< raw
| {"tables":["shop.orders","shop.products","shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$lookup: {from: \"orders\", pipeline: [{$lookup: {from: \"products\", pipeline: [{$lookup: {from: \"reviews\", localField: \"_id\", foreignField: \"sku\", as: \"r\"}}], as: \"p\"}}], as: \"o\"}}])"}
< raw
| {"tables":["shop.orders","shop.products","shop.reviews","shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$lookup: {from: \"users\", localField: \"manager\", foreignField: \"_id\", as: \"users\"}}])"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$lookup: {from: \"active_users\", localField: \"_id\", foreignField: \"_id\", as: \"me\"}}])"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$lookup: {from: \"user_orders\", localField: \"_id\", foreignField: \"_id\", as: \"mine\"}}])"}
< raw
| {"tables":["shop.orders","shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.orders.aggregate([{$lookup: {from: \"users\", localField: \"uid\", foreignField: \"_id\", as: \"orders\"}}])"}
< raw
| {"tables":["shop.orders","shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.orders.aggregate([{$lookup: {from: \"nosuch\", localField: \"uid\", foreignField: \"_id\", as: \"x\"}}])"}
< raw
| {"tables":["shop.nosuch","shop.orders"]}
> raw
| {"op":"Span","database":"shop","statement":"db.orders.aggregate([{$lookup: {from: {db: \"hr\", coll: \"employees\"}, localField: \"uid\", foreignField: \"_id\", as: \"x\"}}])"}
< raw
| {"error":"FailedToParse"}
> raw
| {"op":"Span","database":"shop","statement":"db.orders.aggregate([{$match: {total: 5}}, {$lookup: {from: \"products\", localField: \"sku\", foreignField: \"_id\", as: \"product\"}}, {$replaceRoot: {newRoot: {$first: \"$product\"}}}])"}
< raw
| {"tables":["shop.orders","shop.products"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$graphLookup: {from: \"users\", startWith: \"$manager\", connectFromField: \"manager\", connectToField: \"_id\", as: \"chain\"}}])"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$graphLookup: {from: \"orders\", startWith: \"$_id\", connectFromField: \"uid\", connectToField: \"uid\", as: \"g\"}}])"}
< raw
| {"tables":["shop.orders","shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$graphLookup: {from: \"active_users\", startWith: \"$_id\", connectFromField: \"_id\", connectToField: \"_id\", as: \"g\"}}])"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$unionWith: \"archive\"}])"}
< raw
| {"tables":["shop.archive","shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$unionWith: {coll: \"archive\"}}])"}
< raw
| {"tables":["shop.archive","shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$unionWith: {coll: \"orders\", pipeline: [{$lookup: {from: \"products\", localField: \"sku\", foreignField: \"_id\", as: \"p\"}}]}}])"}
< raw
| {"tables":["shop.orders","shop.products","shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$unionWith: {coll: \"archive\", pipeline: [{$unionWith: \"orders\"}]}}])"}
< raw
| {"tables":["shop.archive","shop.orders","shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$unionWith: \"all_people\"}])"}
< raw
| {"tables":["shop.archive","shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$unionWith: \"active_users\"}, {$lookup: {from: \"orders\", localField: \"_id\", foreignField: \"uid\", as: \"o\"}}])"}
< raw
| {"tables":["shop.orders","shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$facet: {o: [{$lookup: {from: \"orders\", localField: \"_id\", foreignField: \"uid\", as: \"o\"}}], n: [{$count: \"n\"}]}}])"}
< raw
| {"tables":["shop.orders","shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$facet: {a: [{$lookup: {from: \"orders\", localField: \"_id\", foreignField: \"uid\", as: \"o\"}}], b: [{$lookup: {from: \"archive\", localField: \"_id\", foreignField: \"_id\", as: \"x\"}}]}}])"}
< raw
| {"tables":["shop.archive","shop.orders","shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$unionWith: {coll: \"orders\", pipeline: [{$facet: {p: [{$lookup: {from: \"products\", localField: \"sku\", foreignField: \"_id\", as: \"p\"}}]}}]}}])"}
< raw
| {"tables":["shop.orders","shop.products","shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$facet: {u: [{$unionWith: \"archive\"}]}}])"}
< raw
| {"tables":["shop.archive","shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$out: \"report\"}])"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$out: {db: \"hr\", coll: \"report\"}}])"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$lookup: {from: \"orders\", localField: \"_id\", foreignField: \"uid\", as: \"o\"}}, {$out: \"report\"}])"}
< raw
| {"tables":["shop.orders","shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$merge: {into: \"report\"}}])"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$merge: {into: {db: \"hr\", coll: \"report\"}, on: \"_id\", whenMatched: \"replace\"}}])"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.users.aggregate([{$unionWith: \"archive\"}, {$merge: \"report\"}])"}
< raw
| {"tables":["shop.archive","shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.active_users.find()"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.active_users.find({name: \"a\"}, {email: 1})"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.active_users.findOne()"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.active_users.countDocuments({})"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.active_users.distinct(\"name\")"}
< raw
| {"tables":["shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.user_orders.find()"}
< raw
| {"tables":["shop.orders","shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.active_user_orders.find({})"}
< raw
| {"tables":["shop.orders","shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.all_people.find()"}
< raw
| {"tables":["shop.archive","shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.order_products.find()"}
< raw
| {"tables":["shop.orders","shop.products","shop.reviews"]}
> raw
| {"op":"Span","database":"shop","statement":"db.user_graph.find()"}
< raw
| {"tables":["shop.orders","shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.order_facets.aggregate([])"}
< raw
| {"tables":["shop.orders","shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.active_users.aggregate([{$lookup: {from: \"user_orders\", localField: \"_id\", foreignField: \"_id\", as: \"o\"}}])"}
< raw
| {"tables":["shop.orders","shop.users"]}
> raw
| {"op":"Span","database":"shop","statement":"db.getCollection(\"Mixed View\").find()"}
< raw
| {"tables":["shop.Mixed Case"]}
> raw
| {"op":"Span","database":"shop","statement":"db.getSiblingDB(\"hr\").employees.find()"}
< raw
| {"tables":["hr.employees"]}
> raw
| {"op":"Span","database":"shop","statement":"db.getSiblingDB(\"hr\").employees.find({name: \"ann\"}).limit(1)"}
< raw
| {"tables":["hr.employees"]}
> raw
| {"op":"Span","database":"shop","statement":"db.getSiblingDB(\"hr\").getCollection(\"employees\").find()"}
< raw
| {"tables":["hr.employees"]}
> raw
| {"op":"Span","database":"shop","statement":"db.getSiblingDB(\"hr\").employees.findOne()"}
< raw
| {"tables":["hr.employees"]}
> raw
| {"op":"Span","database":"shop","statement":"db.getSiblingDB(\"hr\").employees.countDocuments({})"}
< raw
| {"tables":["hr.employees"]}
> raw
| {"op":"Span","database":"shop","statement":"db.getSiblingDB(\"hr\").staff.find()"}
< raw
| {"tables":["hr.employees"]}
> raw
| {"op":"Span","database":"shop","statement":"db.getSiblingDB(\"hr\").staff_depts.find()"}
< raw
| {"tables":["hr.departments","hr.employees"]}
> raw
| {"op":"Span","database":"shop","statement":"db.getSiblingDB(\"hr\").employees.aggregate([{$lookup: {from: \"departments\", localField: \"dept\", foreignField: \"_id\", as: \"d\"}}])"}
< raw
| {"tables":["hr.departments","hr.employees"]}
> raw
| {"op":"Span","database":"shop","statement":"db.getSiblingDB(\"hr\").employees.aggregate([{$unionWith: \"staff\"}])"}
< raw
| {"tables":["hr.employees"]}
> raw
| {"op":"Span","database":"shop","statement":"db.getSiblingDB(\"hr\").users.find()"}
< raw
| {"tables":["hr.users"]}
> raw
| {"op":"Span","database":"hr","statement":"db.employees.find()"}
< raw
| {"tables":["hr.employees"]}
> raw
| {"op":"Span","database":"hr","statement":"db.staff_depts.aggregate([{$lookup: {from: \"employees\", localField: \"boss\", foreignField: \"_id\", as: \"boss\"}}])"}
< raw
| {"tables":["hr.departments","hr.employees"]}
> raw
| {"op":"Span","database":"hr","statement":"db.departments.aggregate([{$graphLookup: {from: \"employees\", startWith: \"$_id\", connectFromField: \"dept\", connectToField: \"dept\", as: \"people\"}}])"}
< raw
| {"tables":["hr.departments","hr.employees"]}
> raw
| {"op":"Span","database":"hr","statement":"db.users.find()"}
< raw
| {"tables":["hr.users"]}
