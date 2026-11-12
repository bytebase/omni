-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestTiDBPlacementPolicyContainer default#1
> exec
| CREATE DATABASE IF NOT EXISTS omni_pp_test
< ok 0 0
> exec
| USE omni_pp_test
< ok 0 0
> exec
| CREATE PLACEMENT POLICY pp_cpr PRIMARY_REGION = 'us-east-1' REGIONS = 'us-east-1'
< ok 0 0
> exec
| DROP PLACEMENT POLICY IF EXISTS pp_cpr
< ok 0 0
> exec
| CREATE PLACEMENT POLICY pp_cf FOLLOWERS = 2
< ok 0 0
> exec
| DROP PLACEMENT POLICY IF EXISTS pp_cf
< ok 0 0
> exec
| CREATE PLACEMENT POLICY pp_ccl CONSTRAINTS='[+region=us-east-1]'
< ok 0 0
> exec
| DROP PLACEMENT POLICY IF EXISTS pp_ccl
< ok 0 0
> exec
| CREATE PLACEMENT POLICY pp_or FOLLOWERS = 2
< ok 0 0
> exec
| CREATE OR REPLACE PLACEMENT POLICY pp_or FOLLOWERS = 3
< ok 0 0
> exec
| DROP PLACEMENT POLICY IF EXISTS pp_or
< ok 0 0
> exec
| CREATE PLACEMENT POLICY pp_alter FOLLOWERS = 2
< ok 0 0
> exec
| ALTER PLACEMENT POLICY pp_alter PRIMARY_REGION = 'us-east-1' REGIONS = 'us-east-1'
< ok 0 0
> exec
| DROP PLACEMENT POLICY IF EXISTS pp_alter
< ok 0 0
> exec
| DROP PLACEMENT POLICY IF EXISTS never_defined_pp
< ok 0 0
> exec
| CREATE PLACEMENT POLICY pp_tbl PRIMARY_REGION='us-east-1' REGIONS='us-east-1'
< ok 0 0
> exec
| CREATE TABLE t_pp (id INT) PLACEMENT POLICY = pp_tbl
< ok 0 0
> exec
| DROP TABLE IF EXISTS t_pp; DROP PLACEMENT POLICY IF EXISTS pp_tbl
< ok 0 0
> exec
| DROP DATABASE IF EXISTS omni_pp_test
< ok 0 0
