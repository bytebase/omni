-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestZeroLengthKeyPartMatchesTiDB default#1
> exec
| DROP DATABASE IF EXISTS testzerolengthkeypartmatchestidb
< ok 0 0
> exec
| CREATE DATABASE testzerolengthkeypartmatchestidb
< ok 0 0
> exec
| USE testzerolengthkeypartmatchestidb
< ok 0 0
> exec
| CREATE TABLE zl_t (s VARCHAR(10))
< ok 0 0
> exec
| CREATE TABLE zl_a (s VARCHAR(10), KEY (s(0)))
< error
| 1391 HY000 Key part 's' length cannot be 0
> exec
| CREATE TABLE zl_b (s VARCHAR(10), KEY (nosuch(0)))
< error
| 1391 HY000 Key part 'nosuch' length cannot be 0
> exec
| CREATE TABLE zl_c (a INT, KEY (a(0)))
< error
| 1391 HY000 Key part 'a' length cannot be 0
> exec
| ALTER TABLE zl_t ADD KEY (s(0))
< error
| 1391 HY000 Key part 's' length cannot be 0
> exec
| CREATE INDEX zl_i ON zl_t (s(0))
< error
| 1391 HY000 Key part 's' length cannot be 0
> exec
| CREATE INDEX zl_i ON zl_nosuch (s(0))
< error
| 1391 HY000 Key part 's' length cannot be 0
> exec
| CREATE TABLE zl_d (s VARCHAR(10), KEY (s(5)))
< ok 0 0
> exec
| DROP DATABASE IF EXISTS testzerolengthkeypartmatchestidb
< ok 0 0
