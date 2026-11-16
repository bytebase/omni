-- pin: default trinodb/trino:482 sha256:90b35b7c603eaa1f889bf03981a62b75f998ee6c0f851d9f4e341b49a57022b6 482

== TestDumpGaps/quoted_identifier default#1
> exec
| CREATE SCHEMA memory.omni_gap_quote
< ok 0 -
> exec
| CREATE TABLE memory.omni_gap_quote."a""b" (id integer, "c""d" varchar)
< ok 0 -
> query
| SELECT table_schem FROM system.jdbc.schemas WHERE table_catalog = ? ORDER BY table_schem
: string
| memory
< result
< column "table_schem" "VARCHAR"
< row
: string
| default
< row
: string
| information_schema
< row
: string
| omni_gap_quote
> query
| SELECT table_schem, table_name FROM system.jdbc.tables WHERE table_cat = ? AND table_type = 'TABLE' ORDER BY table_schem, table_name
: string
| memory
< result
< column "table_schem" "VARCHAR"
< column "table_name" "VARCHAR"
< row
: string
| information_schema
: string
| applicable_roles
< row
: string
| information_schema
: string
| columns
< row
: string
| information_schema
: string
| enabled_roles
< row
: string
| information_schema
: string
| roles
< row
: string
| information_schema
: string
| schemata
< row
: string
| information_schema
: string
| table_privileges
< row
: string
| information_schema
: string
| tables
< row
: string
| information_schema
: string
| views
< row
: string
| omni_gap_quote
: string
| a"b
> query
| SELECT table_schem, table_name, column_name, type_name, is_nullable FROM system.jdbc.columns WHERE table_cat = ? ORDER BY table_schem, table_name, ordinal_position
: string
| memory
< result
< column "table_schem" "VARCHAR"
< column "table_name" "VARCHAR"
< column "column_name" "VARCHAR"
< column "type_name" "VARCHAR"
< column "is_nullable" "VARCHAR"
< row
: string
| information_schema
: string
| applicable_roles
: string
| grantee
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| applicable_roles
: string
| grantee_type
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| applicable_roles
: string
| role_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| applicable_roles
: string
| is_grantable
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| table_catalog
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| table_schema
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| table_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| column_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| ordinal_position
: string
| bigint
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| column_default
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| is_nullable
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| data_type
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| enabled_roles
: string
| role_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| roles
: string
| role_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| schemata
: string
| catalog_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| schemata
: string
| schema_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| grantor
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| grantor_type
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| grantee
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| grantee_type
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| table_catalog
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| table_schema
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| table_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| privilege_type
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| is_grantable
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| with_hierarchy
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| tables
: string
| table_catalog
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| tables
: string
| table_schema
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| tables
: string
| table_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| tables
: string
| table_type
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| views
: string
| table_catalog
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| views
: string
| table_schema
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| views
: string
| table_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| views
: string
| view_definition
: string
| varchar
: string
| YES
< row
: string
| omni_gap_quote
: string
| a"b
: string
| id
: string
| integer
: string
| YES
< row
: string
| omni_gap_quote
: string
| a"b
: string
| c"d
: string
| varchar
: string
| YES
> exec
| DROP SCHEMA IF EXISTS memory.omni_gap_quote CASCADE
< ok 0 -
> exec
| CREATE SCHEMA memory.omni_gap_quote
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS "information_schema"."applicable_roles" (
|     "grantee" varchar,
|     "grantee_type" varchar,
|     "role_name" varchar,
|     "is_grantable" varchar
| )
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS "information_schema"."columns" (
|     "table_catalog" varchar,
|     "table_schema" varchar,
|     "table_name" varchar,
|     "column_name" varchar,
|     "ordinal_position" bigint,
|     "column_default" varchar,
|     "is_nullable" varchar,
|     "data_type" varchar
| )
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS "information_schema"."enabled_roles" (
|     "role_name" varchar
| )
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS "information_schema"."roles" (
|     "role_name" varchar
| )
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS "information_schema"."schemata" (
|     "catalog_name" varchar,
|     "schema_name" varchar
| )
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS "information_schema"."table_privileges" (
|     "grantor" varchar,
|     "grantor_type" varchar,
|     "grantee" varchar,
|     "grantee_type" varchar,
|     "table_catalog" varchar,
|     "table_schema" varchar,
|     "table_name" varchar,
|     "privilege_type" varchar,
|     "is_grantable" varchar,
|     "with_hierarchy" varchar
| )
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS "information_schema"."tables" (
|     "table_catalog" varchar,
|     "table_schema" varchar,
|     "table_name" varchar,
|     "table_type" varchar
| )
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS "information_schema"."views" (
|     "table_catalog" varchar,
|     "table_schema" varchar,
|     "table_name" varchar,
|     "view_definition" varchar
| )
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS "omni_gap_quote"."a"b" (
|     "id" integer,
|     "c"d" varchar
| )
< error
| {"Status":200,"Type":"USER_ERROR","Message":"line 1:48: mismatched input 'b'. Expecting: '(', '.', 'AS', 'COMMENT', 'WITH'","SqlState":"","Name":"SYNTAX_ERROR","Code":1,"Line":1,"Column":48}

== TestDumpGaps/schema_not_created default#1
> exec
| CREATE SCHEMA memory.omni_gap_schema
< ok 0 -
> exec
| CREATE TABLE memory.omni_gap_schema.t (id integer)
< ok 0 -
> query
| SELECT table_schem FROM system.jdbc.schemas WHERE table_catalog = ? ORDER BY table_schem
: string
| memory
< result
< column "table_schem" "VARCHAR"
< row
: string
| default
< row
: string
| information_schema
< row
: string
| omni_gap_quote
< row
: string
| omni_gap_schema
> query
| SELECT table_schem, table_name FROM system.jdbc.tables WHERE table_cat = ? AND table_type = 'TABLE' ORDER BY table_schem, table_name
: string
| memory
< result
< column "table_schem" "VARCHAR"
< column "table_name" "VARCHAR"
< row
: string
| information_schema
: string
| applicable_roles
< row
: string
| information_schema
: string
| columns
< row
: string
| information_schema
: string
| enabled_roles
< row
: string
| information_schema
: string
| roles
< row
: string
| information_schema
: string
| schemata
< row
: string
| information_schema
: string
| table_privileges
< row
: string
| information_schema
: string
| tables
< row
: string
| information_schema
: string
| views
< row
: string
| omni_gap_schema
: string
| t
> query
| SELECT table_schem, table_name, column_name, type_name, is_nullable FROM system.jdbc.columns WHERE table_cat = ? ORDER BY table_schem, table_name, ordinal_position
: string
| memory
< result
< column "table_schem" "VARCHAR"
< column "table_name" "VARCHAR"
< column "column_name" "VARCHAR"
< column "type_name" "VARCHAR"
< column "is_nullable" "VARCHAR"
< row
: string
| information_schema
: string
| applicable_roles
: string
| grantee
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| applicable_roles
: string
| grantee_type
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| applicable_roles
: string
| role_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| applicable_roles
: string
| is_grantable
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| table_catalog
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| table_schema
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| table_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| column_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| ordinal_position
: string
| bigint
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| column_default
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| is_nullable
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| data_type
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| enabled_roles
: string
| role_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| roles
: string
| role_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| schemata
: string
| catalog_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| schemata
: string
| schema_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| grantor
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| grantor_type
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| grantee
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| grantee_type
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| table_catalog
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| table_schema
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| table_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| privilege_type
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| is_grantable
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| with_hierarchy
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| tables
: string
| table_catalog
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| tables
: string
| table_schema
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| tables
: string
| table_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| tables
: string
| table_type
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| views
: string
| table_catalog
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| views
: string
| table_schema
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| views
: string
| table_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| views
: string
| view_definition
: string
| varchar
: string
| YES
< row
: string
| omni_gap_schema
: string
| t
: string
| id
: string
| integer
: string
| YES
> exec
| DROP SCHEMA IF EXISTS memory.omni_gap_schema CASCADE
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS "information_schema"."applicable_roles" (
|     "grantee" varchar,
|     "grantee_type" varchar,
|     "role_name" varchar,
|     "is_grantable" varchar
| )
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS "information_schema"."columns" (
|     "table_catalog" varchar,
|     "table_schema" varchar,
|     "table_name" varchar,
|     "column_name" varchar,
|     "ordinal_position" bigint,
|     "column_default" varchar,
|     "is_nullable" varchar,
|     "data_type" varchar
| )
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS "information_schema"."enabled_roles" (
|     "role_name" varchar
| )
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS "information_schema"."roles" (
|     "role_name" varchar
| )
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS "information_schema"."schemata" (
|     "catalog_name" varchar,
|     "schema_name" varchar
| )
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS "information_schema"."table_privileges" (
|     "grantor" varchar,
|     "grantor_type" varchar,
|     "grantee" varchar,
|     "grantee_type" varchar,
|     "table_catalog" varchar,
|     "table_schema" varchar,
|     "table_name" varchar,
|     "privilege_type" varchar,
|     "is_grantable" varchar,
|     "with_hierarchy" varchar
| )
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS "information_schema"."tables" (
|     "table_catalog" varchar,
|     "table_schema" varchar,
|     "table_name" varchar,
|     "table_type" varchar
| )
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS "information_schema"."views" (
|     "table_catalog" varchar,
|     "table_schema" varchar,
|     "table_name" varchar,
|     "view_definition" varchar
| )
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS "omni_gap_schema"."t" (
|     "id" integer
| )
< error
| {"Status":200,"Type":"USER_ERROR","Message":"Schema omni_gap_schema not found","SqlState":"","Name":"NOT_FOUND","Code":5,"Line":0,"Column":0}

== TestDumpGaps/view default#1
> exec
| CREATE SCHEMA memory.omni_gap_view
< ok 0 -
> exec
| CREATE TABLE memory.omni_gap_view.t (id integer)
< ok 0 -
> exec
| CREATE VIEW memory.omni_gap_view.v AS SELECT id FROM memory.omni_gap_view.t
< ok 0 -
> query
| SELECT table_schem FROM system.jdbc.schemas WHERE table_catalog = ? ORDER BY table_schem
: string
| memory
< result
< column "table_schem" "VARCHAR"
< row
: string
| default
< row
: string
| information_schema
< row
: string
| omni_gap_quote
< row
: string
| omni_gap_view
> query
| SELECT table_schem, table_name FROM system.jdbc.tables WHERE table_cat = ? AND table_type = 'TABLE' ORDER BY table_schem, table_name
: string
| memory
< result
< column "table_schem" "VARCHAR"
< column "table_name" "VARCHAR"
< row
: string
| information_schema
: string
| applicable_roles
< row
: string
| information_schema
: string
| columns
< row
: string
| information_schema
: string
| enabled_roles
< row
: string
| information_schema
: string
| roles
< row
: string
| information_schema
: string
| schemata
< row
: string
| information_schema
: string
| table_privileges
< row
: string
| information_schema
: string
| tables
< row
: string
| information_schema
: string
| views
< row
: string
| omni_gap_view
: string
| t
> query
| SELECT table_schem, table_name, column_name, type_name, is_nullable FROM system.jdbc.columns WHERE table_cat = ? ORDER BY table_schem, table_name, ordinal_position
: string
| memory
< result
< column "table_schem" "VARCHAR"
< column "table_name" "VARCHAR"
< column "column_name" "VARCHAR"
< column "type_name" "VARCHAR"
< column "is_nullable" "VARCHAR"
< row
: string
| information_schema
: string
| applicable_roles
: string
| grantee
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| applicable_roles
: string
| grantee_type
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| applicable_roles
: string
| role_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| applicable_roles
: string
| is_grantable
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| table_catalog
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| table_schema
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| table_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| column_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| ordinal_position
: string
| bigint
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| column_default
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| is_nullable
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| data_type
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| enabled_roles
: string
| role_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| roles
: string
| role_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| schemata
: string
| catalog_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| schemata
: string
| schema_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| grantor
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| grantor_type
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| grantee
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| grantee_type
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| table_catalog
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| table_schema
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| table_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| privilege_type
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| is_grantable
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| with_hierarchy
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| tables
: string
| table_catalog
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| tables
: string
| table_schema
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| tables
: string
| table_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| tables
: string
| table_type
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| views
: string
| table_catalog
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| views
: string
| table_schema
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| views
: string
| table_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| views
: string
| view_definition
: string
| varchar
: string
| YES
< row
: string
| omni_gap_view
: string
| t
: string
| id
: string
| integer
: string
| YES
< row
: string
| omni_gap_view
: string
| v
: string
| id
: string
| integer
: string
| YES
> exec
| DROP SCHEMA IF EXISTS memory.omni_gap_view CASCADE
< ok 0 -
> exec
| CREATE SCHEMA memory.omni_gap_view
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS "information_schema"."applicable_roles" (
|     "grantee" varchar,
|     "grantee_type" varchar,
|     "role_name" varchar,
|     "is_grantable" varchar
| )
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS "information_schema"."columns" (
|     "table_catalog" varchar,
|     "table_schema" varchar,
|     "table_name" varchar,
|     "column_name" varchar,
|     "ordinal_position" bigint,
|     "column_default" varchar,
|     "is_nullable" varchar,
|     "data_type" varchar
| )
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS "information_schema"."enabled_roles" (
|     "role_name" varchar
| )
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS "information_schema"."roles" (
|     "role_name" varchar
| )
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS "information_schema"."schemata" (
|     "catalog_name" varchar,
|     "schema_name" varchar
| )
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS "information_schema"."table_privileges" (
|     "grantor" varchar,
|     "grantor_type" varchar,
|     "grantee" varchar,
|     "grantee_type" varchar,
|     "table_catalog" varchar,
|     "table_schema" varchar,
|     "table_name" varchar,
|     "privilege_type" varchar,
|     "is_grantable" varchar,
|     "with_hierarchy" varchar
| )
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS "information_schema"."tables" (
|     "table_catalog" varchar,
|     "table_schema" varchar,
|     "table_name" varchar,
|     "table_type" varchar
| )
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS "information_schema"."views" (
|     "table_catalog" varchar,
|     "table_schema" varchar,
|     "table_name" varchar,
|     "view_definition" varchar
| )
< ok 0 -
> exec
| CREATE TABLE IF NOT EXISTS "omni_gap_view"."t" (
|     "id" integer
| )
< ok 0 -
> query
| SELECT table_schem FROM system.jdbc.schemas WHERE table_catalog = ? ORDER BY table_schem
: string
| memory
< result
< column "table_schem" "VARCHAR"
< row
: string
| default
< row
: string
| information_schema
< row
: string
| omni_gap_quote
< row
: string
| omni_gap_view
> query
| SELECT table_schem, table_name FROM system.jdbc.tables WHERE table_cat = ? AND table_type = 'TABLE' ORDER BY table_schem, table_name
: string
| memory
< result
< column "table_schem" "VARCHAR"
< column "table_name" "VARCHAR"
< row
: string
| information_schema
: string
| applicable_roles
< row
: string
| information_schema
: string
| columns
< row
: string
| information_schema
: string
| enabled_roles
< row
: string
| information_schema
: string
| roles
< row
: string
| information_schema
: string
| schemata
< row
: string
| information_schema
: string
| table_privileges
< row
: string
| information_schema
: string
| tables
< row
: string
| information_schema
: string
| views
< row
: string
| omni_gap_view
: string
| t
> query
| SELECT table_schem, table_name, column_name, type_name, is_nullable FROM system.jdbc.columns WHERE table_cat = ? ORDER BY table_schem, table_name, ordinal_position
: string
| memory
< result
< column "table_schem" "VARCHAR"
< column "table_name" "VARCHAR"
< column "column_name" "VARCHAR"
< column "type_name" "VARCHAR"
< column "is_nullable" "VARCHAR"
< row
: string
| information_schema
: string
| applicable_roles
: string
| grantee
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| applicable_roles
: string
| grantee_type
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| applicable_roles
: string
| role_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| applicable_roles
: string
| is_grantable
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| table_catalog
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| table_schema
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| table_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| column_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| ordinal_position
: string
| bigint
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| column_default
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| is_nullable
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| columns
: string
| data_type
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| enabled_roles
: string
| role_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| roles
: string
| role_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| schemata
: string
| catalog_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| schemata
: string
| schema_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| grantor
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| grantor_type
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| grantee
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| grantee_type
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| table_catalog
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| table_schema
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| table_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| privilege_type
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| is_grantable
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| table_privileges
: string
| with_hierarchy
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| tables
: string
| table_catalog
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| tables
: string
| table_schema
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| tables
: string
| table_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| tables
: string
| table_type
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| views
: string
| table_catalog
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| views
: string
| table_schema
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| views
: string
| table_name
: string
| varchar
: string
| YES
< row
: string
| information_schema
: string
| views
: string
| view_definition
: string
| varchar
: string
| YES
< row
: string
| omni_gap_view
: string
| t
: string
| id
: string
| integer
: string
| YES
> query
| SELECT table_name FROM system.jdbc.tables WHERE table_cat = ? AND table_schem = ? AND table_type = 'VIEW'
: string
| memory
: string
| omni_gap_view
< result
< column "table_name" "VARCHAR"
