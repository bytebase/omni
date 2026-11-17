-- pin: default trinodb/trino:482 sha256:90b35b7c603eaa1f889bf03981a62b75f998ee6c0f851d9f4e341b49a57022b6 482

== TestSyncDatabaseMissing default#1
> query
| SELECT table_schem FROM system.jdbc.schemas WHERE table_catalog = ? ORDER BY table_schem
: string
| omni_no_such_catalog
< result
< column "table_schem" "VARCHAR"
> query
| SELECT table_schem, table_name FROM system.jdbc.tables WHERE table_cat = ? AND table_type = 'TABLE' ORDER BY table_schem, table_name
: string
| omni_no_such_catalog
< result
< column "table_schem" "VARCHAR"
< column "table_name" "VARCHAR"
> query
| SELECT table_schema, table_name, view_definition FROM "omni_no_such_catalog".information_schema.views ORDER BY table_schema, table_name
< error
| {"Status":200,"Type":"USER_ERROR","Message":"line 1:55: Catalog 'omni_no_such_catalog' not found","SqlState":"","Name":"CATALOG_NOT_FOUND","Code":44,"Line":1,"Column":55}
