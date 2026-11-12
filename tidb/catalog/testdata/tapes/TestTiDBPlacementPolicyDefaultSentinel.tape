-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestTiDBPlacementPolicyDefaultSentinel default#1
> exec
| CREATE DATABASE IF NOT EXISTS omni_pp_default
< ok 0 0
> exec
| USE omni_pp_default
< ok 0 0
> exec
| CREATE TABLE t_def_id (id INT) PLACEMENT POLICY = default
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_def_id
< ok 0 0
> exec
| CREATE TABLE t_def_str (id INT) PLACEMENT POLICY = 'default'
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_def_str
< ok 0 0
> exec
| CREATE TABLE t_def_up (id INT) PLACEMENT POLICY = DEFAULT
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_def_up
< ok 0 0
> exec
| DROP DATABASE IF EXISTS omni_pp_default
< ok 0 0
