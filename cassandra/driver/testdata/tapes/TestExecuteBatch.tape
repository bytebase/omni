-- pin: driver cassandra:5.0.9 sha256:acd754f07c73f3a166bb477354b56ce9a05690c3e5e2aaacc1aab2afddba6436 5.0.9

== TestExecuteBatch driver#1
> raw
| {"op":"exec","statement":"CREATE TABLE batched (id int PRIMARY KEY, n int)"}
< raw
| {}
> raw
| {"op":"exec","statement":"BEGIN BATCH\n\t\t\tINSERT INTO batched (id, n) VALUES (1, 10);\n\t\t\tINSERT INTO batched (id, n) VALUES (2, 20);\n\t\tAPPLY BATCH"}
< raw
| {}
> raw
| {"op":"query","statement":"SELECT id, n FROM batched"}
< raw
| {"columns":[{"name":"id","type":9},{"name":"n","type":9}],"rowColumns":["id","n"],"rows":[[{"kind":"int","text":"1"},{"kind":"int","text":"10"}],[{"kind":"int","text":"2"},{"kind":"int","text":"20"}]]}
