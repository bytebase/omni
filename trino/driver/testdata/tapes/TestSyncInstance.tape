-- pin: default trinodb/trino:482 sha256:90b35b7c603eaa1f889bf03981a62b75f998ee6c0f851d9f4e341b49a57022b6 482

== TestSyncInstance default#1
> query
| SELECT VERSION()
< result
< column "_col0" "VARCHAR"
< row
: string
| 482
> query
| SELECT catalog_name FROM system.metadata.catalogs ORDER BY catalog_name
< result
< column "catalog_name" "VARCHAR"
< row
: string
| jmx
< row
: string
| memory
< row
: string
| system
< row
: string
| tpcds
< row
: string
| tpch
