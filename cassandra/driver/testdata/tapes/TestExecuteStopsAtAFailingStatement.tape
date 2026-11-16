-- pin: driver cassandra:5.0.9 sha256:acd754f07c73f3a166bb477354b56ce9a05690c3e5e2aaacc1aab2afddba6436 5.0.9

== TestExecuteStopsAtAFailingStatement driver#1
> raw
| {"op":"exec","statement":"CREATE TABLE kept (id int PRIMARY KEY)"}
< raw
| {}
> raw
| {"op":"exec","statement":"INSERT INTO kept (id) VALUES (1)"}
< raw
| {}
> raw
| {"op":"exec","statement":"INSERT INTO missing (id) VALUES (2)"}
< raw
| {"error":{"code":8704,"message":"table missing does not exist"}}
> raw
| {"op":"query","statement":"SELECT id FROM kept"}
< raw
| {"columns":[{"name":"id","type":9}],"rowColumns":["id"],"rows":[[{"kind":"int","text":"1"}]]}
