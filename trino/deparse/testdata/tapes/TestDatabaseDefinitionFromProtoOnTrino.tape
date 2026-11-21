-- pin: default trinodb/trino:482 sha256:90b35b7c603eaa1f889bf03981a62b75f998ee6c0f851d9f4e341b49a57022b6 482

== TestDatabaseDefinitionFromProtoOnTrino default#1
> exec
| CREATE SCHEMA IF NOT EXISTS "default"
< ok - -
> exec
| CREATE FUNCTION "memory"."default".add_one(x bigint) RETURNS bigint COMMENT 'increment by one' RETURN x + 1
< ok - -
> exec
| CREATE TABLE IF NOT EXISTS "default"."items" (
|     "id" bigint DEFAULT 1 -- seed
|  NOT NULL,
|     "label" varchar DEFAULT 'ok' COMMENT 'display label'
| ) COMMENT 'catalog items'
< ok - -
> exec
| CREATE OR REPLACE VIEW "default"."v_items" COMMENT 'item view' AS
| SELECT id, "memory"."default".add_one(id) AS next_id FROM "memory"."default".items
< ok - -
> exec
| COMMENT ON COLUMN "default"."v_items"."id" IS 'item id'
< ok - -
> exec
| CREATE MATERIALIZED VIEW IF NOT EXISTS "default"."mv_items" AS
| SELECT id FROM "memory"."default".items
< error
| NOT_SUPPORTED 13 0:0
| This connector does not support creating materialized views
