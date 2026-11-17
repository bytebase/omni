-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestPingMissingDatabase default#1
> connect
< error
| 3D000 database "omni_no_such_database" does not exist
| severity: FATAL
