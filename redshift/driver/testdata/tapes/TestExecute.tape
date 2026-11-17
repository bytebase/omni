-- pin: driver aws:redshift:1.0.539507 sha256:afc4041fe6e35ce273970308e430eb8245ea9c4bbf597c497bb8eb8bb00810d2 2026-11-17

== TestExecute driver#1
> begin
< ok
> exec
| CREATE TEMP TABLE omni_execute_probe (a int)
< ok 0 -
> commit
< ok
> begin
< ok
> exec
| INSERT INTO omni_execute_probe VALUES (1), (2)
< ok 2 -
> commit
< ok
> begin
< ok
> exec
| INSERT INTO omni_execute_probe VALUES (3);
< ok 1 -
> exec
|  INSERT INTO omni_execute_probe VALUES ('not an int')
< error
| 22P02 invalid input syntax for integer: "not an int"
> rollback
< ok
> query
| SELECT COUNT(*) FROM omni_execute_probe
< result
< column "count" "INT8" "int64"
< row
: int64 2
