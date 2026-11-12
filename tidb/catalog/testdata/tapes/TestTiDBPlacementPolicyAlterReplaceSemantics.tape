-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestTiDBPlacementPolicyAlterReplaceSemantics default#1
> exec
| CREATE PLACEMENT POLICY pp_rb FOLLOWERS = 2
< ok 0 0
> exec
| ALTER PLACEMENT POLICY pp_rb PRIMARY_REGION = 'us-west-1' REGIONS = 'us-west-1'
< ok 0 0
> query
| SHOW CREATE PLACEMENT POLICY pp_rb
< result
< column "Policy" "VARCHAR"
< column "Create Policy" "VARCHAR"
< row
: bytes
| pp_rb
: bytes
| CREATE PLACEMENT POLICY `pp_rb` PRIMARY_REGION="us-west-1" REGIONS="us-west-1"
> exec
| DROP PLACEMENT POLICY IF EXISTS pp_rb
< ok 0 0
