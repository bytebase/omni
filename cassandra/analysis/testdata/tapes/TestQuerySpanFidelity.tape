-- pin: driver cassandra:5.0.9 sha256:acd754f07c73f3a166bb477354b56ce9a05690c3e5e2aaacc1aab2afddba6436 5.0.9

== TestQuerySpanFidelity driver#1
> raw
| {"op":"spanFixture","statements":["CREATE KEYSPACE shop WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}","CREATE KEYSPACE crm WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}","CREATE KEYSPACE \"Audit\" WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}","CREATE TYPE shop.address (street text, city text, zip text)","CREATE TABLE shop.users (id int PRIMARY KEY, name text, email text, phone text, ssn text, score double, created timestamp, tags set\u003ctext\u003e, prefs map\u003ctext, text\u003e, history list\u003cint\u003e, home frozen\u003caddress\u003e)","CREATE TABLE shop.orders (customer_id int, order_id timeuuid, amount decimal, status text, note text, PRIMARY KEY ((customer_id), order_id)) WITH CLUSTERING ORDER BY (order_id DESC)","CREATE TABLE shop.events (tenant text, day date, seq int, kind text, payload blob, region text STATIC, PRIMARY KEY ((tenant, day), seq))","CREATE TABLE shop.\"Accounts\" (\"AccountId\" int PRIMARY KEY, \"DisplayName\" text, \"eMail\" text, plain text)","CREATE TABLE shop.counters (id int PRIMARY KEY, hits counter)","CREATE TABLE shop.pairs (id int PRIMARY KEY, pair tuple\u003cint, text\u003e)","CREATE MATERIALIZED VIEW shop.users_by_email AS SELECT email, id, name, phone FROM shop.users WHERE email IS NOT NULL AND id IS NOT NULL PRIMARY KEY (email, id)","CREATE MATERIALIZED VIEW shop.orders_by_status AS SELECT * FROM shop.orders WHERE status IS NOT NULL AND customer_id IS NOT NULL AND order_id IS NOT NULL PRIMARY KEY (status, customer_id, order_id)","CREATE TABLE crm.contacts (id uuid PRIMARY KEY, full_name text, email text, owner text, vec vector\u003cfloat, 3\u003e)","CREATE TABLE \"Audit\".\"Log\" (id int PRIMARY KEY, \"Who\" text, what text)"]}
< raw
| {"keyspaces":[{"name":"Audit","tables":[{"name":"Log","columns":[{"name":"Who","kind":"regular","position":-1,"type":"text"},{"name":"id","kind":"partition_key","position":0,"type":"int"},{"name":"what","kind":"regular","position":-1,"type":"text"}]}]},{"name":"crm","tables":[{"name":"contacts","columns":[{"name":"email","kind":"regular","position":-1,"type":"text"},{"name":"full_name","kind":"regular","position":-1,"type":"text"},{"name":"id","kind":"partition_key","position":0,"type":"uuid"},{"name":"owner","kind":"regular","position":-1,"type":"text"},{"name":"vec","kind":"regular","position":-1,"type":"vector\u003cfloat, 3\u003e"}]}]},{"name":"shop","tables":[{"name":"Accounts","columns":[{"name":"AccountId","kind":"partition_key","position":0,"type":"int"},{"name":"DisplayName","kind":"regular","position":-1,"type":"text"},{"name":"eMail","kind":"regular","position":-1,"type":"text"},{"name":"plain","kind":"regular","position":-1,"type":"text"}]},{"name":"counters","columns":[{"name":"hits","kind":"regular","position":-1,"type":"counter"},{"name":"id","kind":"partition_key","position":0,"type":"int"}]},{"name":"events","columns":[{"name":"day","kind":"partition_key","position":1,"type":"date"},{"name":"kind","kind":"regular","position":-1,"type":"text"},{"name":"payload","kind":"regular","position":-1,"type":"blob"},{"name":"region","kind":"static","position":-1,"type":"text"},{"name":"seq","kind":"clustering","position":0,"type":"int"},{"name":"tenant","kind":"partition_key","position":0,"type":"text"}]},{"name":"orders","columns":[{"name":"amount","kind":"regular","position":-1,"type":"decimal"},{"name":"customer_id","kind":"partition_key","position":0,"type":"int"},{"name":"note","kind":"regular","position":-1,"type":"text"},{"name":"order_id","kind":"clustering","position":0,"type":"timeuuid"},{"name":"status","kind":"regular","position":-1,"type":"text"}]},{"name":"pairs","columns":[{"name":"id","kind":"partition_key","position":0,"type":"int"},{"name":"pair","kind":"regular","position":-1,"type":"frozen\u003ctuple\u003cint, text\u003e\u003e"}]},{"name":"users","columns":[{"name":"created","kind":"regular","position":-1,"type":"timestamp"},{"name":"email","kind":"regular","position":-1,"type":"text"},{"name":"history","kind":"regular","position":-1,"type":"list\u003cint\u003e"},{"name":"home","kind":"regular","position":-1,"type":"frozen\u003caddress\u003e"},{"name":"id","kind":"partition_key","position":0,"type":"int"},{"name":"name","kind":"regular","position":-1,"type":"text"},{"name":"phone","kind":"regular","position":-1,"type":"text"},{"name":"prefs","kind":"regular","position":-1,"type":"map\u003ctext, text\u003e"},{"name":"score","kind":"regular","position":-1,"type":"double"},{"name":"ssn","kind":"regular","position":-1,"type":"text"},{"name":"tags","kind":"regular","position":-1,"type":"set\u003ctext\u003e"}]}],"views":[{"name":"orders_by_status","base":"orders","columns":[{"name":"amount","kind":"regular","position":-1,"type":"decimal"},{"name":"customer_id","kind":"clustering","position":0,"type":"int"},{"name":"note","kind":"regular","position":-1,"type":"text"},{"name":"order_id","kind":"clustering","position":1,"type":"timeuuid"},{"name":"status","kind":"partition_key","position":0,"type":"text"}]},{"name":"users_by_email","base":"users","columns":[{"name":"email","kind":"partition_key","position":0,"type":"text"},{"name":"id","kind":"clustering","position":0,"type":"int"},{"name":"name","kind":"regular","position":-1,"type":"text"},{"name":"phone","kind":"regular","position":-1,"type":"text"}]}]}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT id, name FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"id"},{"keyspace":"shop","table":"users","name":"name"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT name, id FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"name"},{"keyspace":"shop","table":"users","name":"id"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT id, name, name FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"id"},{"keyspace":"shop","table":"users","name":"name"},{"keyspace":"shop","table":"users","name":"name"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT NAME, Email FROM Users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"name"},{"keyspace":"shop","table":"users","name":"email"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT * FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"id"},{"keyspace":"shop","table":"users","name":"created"},{"keyspace":"shop","table":"users","name":"email"},{"keyspace":"shop","table":"users","name":"history"},{"keyspace":"shop","table":"users","name":"home"},{"keyspace":"shop","table":"users","name":"name"},{"keyspace":"shop","table":"users","name":"phone"},{"keyspace":"shop","table":"users","name":"prefs"},{"keyspace":"shop","table":"users","name":"score"},{"keyspace":"shop","table":"users","name":"ssn"},{"keyspace":"shop","table":"users","name":"tags"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT * FROM orders"}
< raw
| {"columns":[{"keyspace":"shop","table":"orders","name":"customer_id"},{"keyspace":"shop","table":"orders","name":"order_id"},{"keyspace":"shop","table":"orders","name":"amount"},{"keyspace":"shop","table":"orders","name":"note"},{"keyspace":"shop","table":"orders","name":"status"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT * FROM events"}
< raw
| {"columns":[{"keyspace":"shop","table":"events","name":"tenant"},{"keyspace":"shop","table":"events","name":"day"},{"keyspace":"shop","table":"events","name":"seq"},{"keyspace":"shop","table":"events","name":"region"},{"keyspace":"shop","table":"events","name":"kind"},{"keyspace":"shop","table":"events","name":"payload"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT * FROM counters"}
< raw
| {"columns":[{"keyspace":"shop","table":"counters","name":"id"},{"keyspace":"shop","table":"counters","name":"hits"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT * FROM pairs"}
< raw
| {"columns":[{"keyspace":"shop","table":"pairs","name":"id"},{"keyspace":"shop","table":"pairs","name":"pair"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT id AS user_id, name AS n FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"user_id"},{"keyspace":"shop","table":"users","name":"n"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT id AS \"ID\", name AS \"Name\" FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"ID"},{"keyspace":"shop","table":"users","name":"Name"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT name AS email FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"email"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT email FROM users WHERE id = 1"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"email"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT tags, prefs, history, home FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"tags"},{"keyspace":"shop","table":"users","name":"prefs"},{"keyspace":"shop","table":"users","name":"history"},{"keyspace":"shop","table":"users","name":"home"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT hits FROM counters"}
< raw
| {"columns":[{"keyspace":"shop","table":"counters","name":"hits"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT pair FROM pairs"}
< raw
| {"columns":[{"keyspace":"shop","table":"pairs","name":"pair"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT customer_id, order_id, amount FROM orders WHERE customer_id = 1 ORDER BY order_id ASC"}
< raw
| {"columns":[{"keyspace":"shop","table":"orders","name":"customer_id"},{"keyspace":"shop","table":"orders","name":"order_id"},{"keyspace":"shop","table":"orders","name":"amount"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT * FROM orders WHERE customer_id = 1 LIMIT 10"}
< raw
| {"columns":[{"keyspace":"shop","table":"orders","name":"customer_id"},{"keyspace":"shop","table":"orders","name":"order_id"},{"keyspace":"shop","table":"orders","name":"amount"},{"keyspace":"shop","table":"orders","name":"note"},{"keyspace":"shop","table":"orders","name":"status"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT * FROM orders PER PARTITION LIMIT 1"}
< raw
| {"columns":[{"keyspace":"shop","table":"orders","name":"customer_id"},{"keyspace":"shop","table":"orders","name":"order_id"},{"keyspace":"shop","table":"orders","name":"amount"},{"keyspace":"shop","table":"orders","name":"note"},{"keyspace":"shop","table":"orders","name":"status"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT status FROM orders WHERE status = 'open' ALLOW FILTERING"}
< raw
| {"columns":[{"keyspace":"shop","table":"orders","name":"status"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT region FROM events WHERE tenant = 't' AND day = '2024-01-01'"}
< raw
| {"columns":[{"keyspace":"shop","table":"events","name":"region"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT name FROM users WHERE token(id) \u003e 0"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"name"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT name FROM shop.users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"name"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT DISTINCT id FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"id"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT DISTINCT customer_id FROM orders"}
< raw
| {"columns":[{"keyspace":"shop","table":"orders","name":"customer_id"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT DISTINCT tenant, day FROM events"}
< raw
| {"columns":[{"keyspace":"shop","table":"events","name":"tenant"},{"keyspace":"shop","table":"events","name":"day"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT DISTINCT tenant, day, region FROM events"}
< raw
| {"columns":[{"keyspace":"shop","table":"events","name":"tenant"},{"keyspace":"shop","table":"events","name":"day"},{"keyspace":"shop","table":"events","name":"region"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT JSON id, name FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"[json]"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT JSON * FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"[json]"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT JSON * FROM orders"}
< raw
| {"columns":[{"keyspace":"shop","table":"orders","name":"[json]"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT JSON DISTINCT id FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"[json]"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT JSON name AS n FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"[json]"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT JSON home.city FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"[json]"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT JSON count(*) FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"[json]"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT JSON \"AccountId\", \"eMail\" FROM \"Accounts\""}
< raw
| {"columns":[{"keyspace":"shop","table":"Accounts","name":"[json]"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT WRITETIME(name) FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"writetime(name)"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT TTL(email) FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"ttl(email)"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT writetime(name) AS w, ttl(name) AS t FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"w"},{"keyspace":"shop","table":"users","name":"t"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT id, WRITETIME(email), email FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"id"},{"keyspace":"shop","table":"users","name":"writetime(email)"},{"keyspace":"shop","table":"users","name":"email"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT MAXWRITETIME(history) FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"maxwritetime(history)"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT writetime(region) FROM events"}
< raw
| {"columns":[{"keyspace":"shop","table":"events","name":"writetime(region)"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT count(*) FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"count"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT COUNT(1) FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"count"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT count(email) FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"system.count(email)"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT max(score), min(score), avg(score), sum(score) FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"system.max(score)"},{"keyspace":"shop","table":"users","name":"system.min(score)"},{"keyspace":"shop","table":"users","name":"system.avg(score)"},{"keyspace":"shop","table":"users","name":"system.sum(score)"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT sum(score) AS total FROM users WHERE id IN (1, 2)"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"total"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT customer_id, sum(amount) FROM orders GROUP BY customer_id"}
< raw
| {"columns":[{"keyspace":"shop","table":"orders","name":"customer_id"},{"keyspace":"shop","table":"orders","name":"system.sum(amount)"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT customer_id, count(*) AS n FROM orders GROUP BY customer_id"}
< raw
| {"columns":[{"keyspace":"shop","table":"orders","name":"customer_id"},{"keyspace":"shop","table":"orders","name":"n"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT tenant, day, max(seq) FROM events GROUP BY tenant, day"}
< raw
| {"columns":[{"keyspace":"shop","table":"events","name":"tenant"},{"keyspace":"shop","table":"events","name":"day"},{"keyspace":"shop","table":"events","name":"system.max(seq)"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT toJson(home) FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"system.tojson(home)"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT toJson(email) AS j FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"j"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT toJson(prefs), toJson(tags) FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"system.tojson(prefs)"},{"keyspace":"shop","table":"users","name":"system.tojson(tags)"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT token(id) FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"system.token(id)"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT token(id), id FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"system.token(id)"},{"keyspace":"shop","table":"users","name":"id"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT blobAsText(payload) FROM events"}
< raw
| {"columns":[{"keyspace":"shop","table":"events","name":"system.blobastext(payload)"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT textAsBlob(email) FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"system.textasblob(email)"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT CAST(score AS text) FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"cast(score as text)"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT CAST(id AS bigint) AS big FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"big"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT toDate(created), toUnixTimestamp(created) FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"system.todate(created)"},{"keyspace":"shop","table":"users","name":"system.tounixtimestamp(created)"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT toTimestamp(order_id), toUnixTimestamp(order_id) FROM orders"}
< raw
| {"columns":[{"keyspace":"shop","table":"orders","name":"system.totimestamp(order_id)"},{"keyspace":"shop","table":"orders","name":"system.tounixtimestamp(order_id)"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT now() FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"system.now()"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT uuid(), id FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"system.uuid()"},{"keyspace":"shop","table":"users","name":"id"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT id + 1 FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"id + 1"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT score * 2 AS doubled FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"doubled"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT -score FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"-score"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT id, (int) 1 FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"id"},{"keyspace":"shop","table":"users","name":"(int)1"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT mask_default(email) FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"system.mask_default(email)"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT mask_inner(phone, 1, 2) FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"system.mask_inner(phone, 1, 2)"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT mask_hash(ssn) AS h FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"h"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT collection_count(tags) FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"system.collection_count(tags)"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT map_keys(prefs) FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"system.map_keys(prefs)"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT home.city FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"home.city"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT home.city AS c, home.zip FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"c"},{"keyspace":"shop","table":"users","name":"home.zip"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT prefs['lang'] FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"prefs['lang']"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT prefs['a'..'m'] AS some FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"some"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT tags['red'] FROM users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"tags['red']"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT * FROM users_by_email"}
< raw
| {"columns":[{"keyspace":"shop","table":"users_by_email","name":"email"},{"keyspace":"shop","table":"users_by_email","name":"id"},{"keyspace":"shop","table":"users_by_email","name":"name"},{"keyspace":"shop","table":"users_by_email","name":"phone"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT email, name FROM users_by_email WHERE email = 'a@b.c'"}
< raw
| {"columns":[{"keyspace":"shop","table":"users_by_email","name":"email"},{"keyspace":"shop","table":"users_by_email","name":"name"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT * FROM orders_by_status"}
< raw
| {"columns":[{"keyspace":"shop","table":"orders_by_status","name":"status"},{"keyspace":"shop","table":"orders_by_status","name":"customer_id"},{"keyspace":"shop","table":"orders_by_status","name":"order_id"},{"keyspace":"shop","table":"orders_by_status","name":"amount"},{"keyspace":"shop","table":"orders_by_status","name":"note"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT status, amount FROM orders_by_status"}
< raw
| {"columns":[{"keyspace":"shop","table":"orders_by_status","name":"status"},{"keyspace":"shop","table":"orders_by_status","name":"amount"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT count(*) FROM users_by_email"}
< raw
| {"columns":[{"keyspace":"shop","table":"users_by_email","name":"count"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT JSON * FROM users_by_email"}
< raw
| {"columns":[{"keyspace":"shop","table":"users_by_email","name":"[json]"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT * FROM \"Accounts\""}
< raw
| {"columns":[{"keyspace":"shop","table":"Accounts","name":"AccountId"},{"keyspace":"shop","table":"Accounts","name":"DisplayName"},{"keyspace":"shop","table":"Accounts","name":"eMail"},{"keyspace":"shop","table":"Accounts","name":"plain"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT \"AccountId\", \"DisplayName\" FROM \"Accounts\""}
< raw
| {"columns":[{"keyspace":"shop","table":"Accounts","name":"AccountId"},{"keyspace":"shop","table":"Accounts","name":"DisplayName"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT \"eMail\" AS \"Mail\" FROM \"Accounts\""}
< raw
| {"columns":[{"keyspace":"shop","table":"Accounts","name":"Mail"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT plain, PLAIN FROM \"Accounts\""}
< raw
| {"columns":[{"keyspace":"shop","table":"Accounts","name":"plain"},{"keyspace":"shop","table":"Accounts","name":"plain"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT \"DisplayName\" AS display FROM shop.\"Accounts\""}
< raw
| {"columns":[{"keyspace":"shop","table":"Accounts","name":"display"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT * FROM \"Audit\".\"Log\""}
< raw
| {"columns":[{"keyspace":"Audit","table":"Log","name":"id"},{"keyspace":"Audit","table":"Log","name":"Who"},{"keyspace":"Audit","table":"Log","name":"what"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT \"Who\", what FROM \"Audit\".\"Log\""}
< raw
| {"columns":[{"keyspace":"Audit","table":"Log","name":"Who"},{"keyspace":"Audit","table":"Log","name":"what"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT * FROM crm.contacts"}
< raw
| {"columns":[{"keyspace":"crm","table":"contacts","name":"id"},{"keyspace":"crm","table":"contacts","name":"email"},{"keyspace":"crm","table":"contacts","name":"full_name"},{"keyspace":"crm","table":"contacts","name":"owner"},{"keyspace":"crm","table":"contacts","name":"vec"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT full_name, email FROM crm.contacts"}
< raw
| {"columns":[{"keyspace":"crm","table":"contacts","name":"full_name"},{"keyspace":"crm","table":"contacts","name":"email"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT id, owner FROM crm.contacts WHERE id = 5132b130-ae79-11e4-ab27-0800200c9a66"}
< raw
| {"columns":[{"keyspace":"crm","table":"contacts","name":"id"},{"keyspace":"crm","table":"contacts","name":"owner"}]}
> raw
| {"op":"span","keyspace":"shop","statement":"SELECT similarity_cosine(vec, [0.1, 0.2, 0.3]) FROM crm.contacts"}
< raw
| {"columns":[{"keyspace":"crm","table":"contacts","name":"system.similarity_cosine(vec, [0.1, 0.2, 0.3])"}]}
> raw
| {"op":"span","keyspace":"crm","statement":"SELECT * FROM contacts"}
< raw
| {"columns":[{"keyspace":"crm","table":"contacts","name":"id"},{"keyspace":"crm","table":"contacts","name":"email"},{"keyspace":"crm","table":"contacts","name":"full_name"},{"keyspace":"crm","table":"contacts","name":"owner"},{"keyspace":"crm","table":"contacts","name":"vec"}]}
> raw
| {"op":"span","keyspace":"crm","statement":"SELECT email FROM contacts"}
< raw
| {"columns":[{"keyspace":"crm","table":"contacts","name":"email"}]}
> raw
| {"op":"span","keyspace":"crm","statement":"SELECT name, email FROM shop.users"}
< raw
| {"columns":[{"keyspace":"shop","table":"users","name":"name"},{"keyspace":"shop","table":"users","name":"email"}]}
> raw
| {"op":"span","keyspace":"crm","statement":"SELECT * FROM shop.users_by_email"}
< raw
| {"columns":[{"keyspace":"shop","table":"users_by_email","name":"email"},{"keyspace":"shop","table":"users_by_email","name":"id"},{"keyspace":"shop","table":"users_by_email","name":"name"},{"keyspace":"shop","table":"users_by_email","name":"phone"}]}
