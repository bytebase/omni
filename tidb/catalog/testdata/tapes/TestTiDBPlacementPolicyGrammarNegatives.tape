-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestTiDBPlacementPolicyGrammarNegatives default#1
> exec
| CREATE PLACEMENT POLICY pp_neg
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 30 near "" 
> exec
| CREATE PLACEMENT POLICY pp_neg PRIMARY_REGION = 'us',
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 53 near "" 
> exec
| ALTER PLACEMENT POLICY pp_neg
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 29 near "" 
> exec
| ALTER PLACEMENT POLICY pp_neg PRIMARY_REGION = 'us',
< error
| 1064 42000 You have an error in your SQL syntax; check the manual that corresponds to your TiDB version for the right syntax to use line 1 column 52 near "" 
> exec
| DROP PLACEMENT POLICY IF EXISTS pp_neg
< ok 0 0
