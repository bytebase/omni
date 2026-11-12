-- pin: default pingcap/tidb:v8.5.5 sha256:f2178ff6cd26f190c64a92cf867148ec6ee6fa31e214cc402bfbbb6bf5f70f26 v8.5.5

== TestAlterRangeMatchesTiDB default#1
> exec
| ALTER RANGE global PLACEMENT POLICY DEFAULT
< ok 0 0
> exec
| ALTER RANGE meta PLACEMENT POLICY DEFAULT
< ok 0 0
> exec
| DROP PLACEMENT POLICY IF EXISTS omni_range_p
< ok 0 0
> exec
| CREATE PLACEMENT POLICY omni_range_p PRIMARY_REGION="r" REGIONS="r"
< ok 0 0
> exec
| ALTER RANGE global PLACEMENT POLICY omni_range_nosuch
< error
| 8239 HY000 Unknown placement policy 'omni_range_nosuch'
> exec
| ALTER RANGE foo PLACEMENT POLICY DEFAULT
< error
| 1105 HY000 range name is not supported
> exec
| ALTER RANGE foo PLACEMENT POLICY omni_range_nosuch
< error
| 1105 HY000 range name is not supported
> exec
| ALTER RANGE global PLACEMENT POLICY omni_range_p
< ok 0 0
> exec
| DROP PLACEMENT POLICY omni_range_p
< error
| 8241 HY000 Placement policy 'omni_range_p' is still in use
> exec
| ALTER RANGE GLOBAL PLACEMENT POLICY = `default`
< ok 0 0
> exec
| ALTER RANGE META PLACEMENT POLICY = 'OMNI_RANGE_P'
< ok 0 0
> exec
| DROP PLACEMENT POLICY omni_range_p
< error
| 8241 HY000 Placement policy 'omni_range_p' is still in use
> exec
| ALTER RANGE `meta` PLACEMENT POLICY SET DEFAULT
< ok 0 0
> exec
| DROP PLACEMENT POLICY omni_range_p
< ok 0 0
> exec
| ALTER RANGE global PLACEMENT POLICY omni_range_p
< error
| 8239 HY000 Unknown placement policy 'omni_range_p'
> exec
| ALTER RANGE global PLACEMENT POLICY DEFAULT
< ok 0 0
> exec
| ALTER RANGE meta PLACEMENT POLICY DEFAULT
< ok 0 0
> exec
| DROP PLACEMENT POLICY IF EXISTS omni_range_p
< ok 0 0
