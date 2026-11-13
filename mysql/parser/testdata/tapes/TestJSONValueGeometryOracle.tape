-- pin: default mysql:8.0.46 sha256:7dcddc01f13bab2f15cde676d44d01f61fc9f99fe7785e86196dfc07d358ae2b mysql-8.0.46

== TestJSONValueGeometryOracle default#1
> exec
| PREPARE omni_jv FROM 'SELECT JSON_VALUE(a,''$.a'' RETURNING POINT)'
< error
| 4032 22S01 Invalid cast from JSON to POINT.
> exec
| PREPARE omni_jv FROM 'SELECT JSON_VALUE(a,''$.a'' RETURNING POINT)
| ; SELECT FROM'
< error
| 4032 22S01 Invalid cast from JSON to POINT.
