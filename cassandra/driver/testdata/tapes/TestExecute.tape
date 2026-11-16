-- pin: driver cassandra:5.0.9 sha256:acd754f07c73f3a166bb477354b56ce9a05690c3e5e2aaacc1aab2afddba6436 5.0.9

== TestExecute driver#1
> raw
| {"op":"exec","statement":"CREATE TABLE exec_items (id int PRIMARY KEY, note text)"}
< raw
| {}
> raw
| {"op":"exec","statement":"INSERT INTO exec_items (id, note) VALUES (1, 'a')"}
< raw
| {}
> raw
| {"op":"exec","statement":"INSERT INTO exec_items (id, note) VALUES (2, 'b')"}
< raw
| {}
> raw
| {"op":"exec","statement":"UPDATE exec_items SET note = 'z' WHERE id = 1"}
< raw
| {}
> raw
| {"op":"exec","statement":"DELETE FROM exec_items WHERE id = 2"}
< raw
| {}
> raw
| {"op":"query","statement":"SELECT id, note FROM exec_items"}
< raw
| {"columns":[{"name":"id","type":9},{"name":"note","type":13}],"rows":[[{"kind":"int","text":"1"},{"kind":"string","text":"z"}]]}
