-- pin: driver cassandra:5.0.9 sha256:acd754f07c73f3a166bb477354b56ce9a05690c3e5e2aaacc1aab2afddba6436 5.0.9

== TestExecuteLogsItsSteps driver#1
> raw
| {"op":"exec","statement":"CREATE TABLE logged (id int PRIMARY KEY)"}
< raw
| {}
> raw
| {"op":"exec","statement":"INSERT INTO logged (id) VALUES (1)"}
< raw
| {}
> raw
| {"op":"exec","statement":"INSERT INTO logged (id) VALUES (2)"}
< raw
| {}
> raw
| {"op":"exec","statement":"INSERT INTO missing (id) VALUES (3)"}
< raw
| {"error":{"code":8704,"message":"table missing does not exist"}}
