-- pin: default trinodb/trino:482 sha256:90b35b7c603eaa1f889bf03981a62b75f998ee6c0f851d9f4e341b49a57022b6 482

== TestGeneratedSDLParsesOnTrino default#1
> exec
| CREATE TABLE IF NOT EXISTS "default"."testtable" (
|     "id" bigint NOT NULL,
|     "name" varchar
| )
< ok - -
> exec
| CREATE TABLE IF NOT EXISTS "default"."single_col" (
|     "c" integer NOT NULL
| )
< ok - -
> exec
| CREATE TABLE IF NOT EXISTS "default"."typed" (
|     "a" varchar(10) NOT NULL,
|     "b" decimal(10,2),
|     "c" array(integer),
|     "d" map(varchar, bigint),
|     "e" row(x integer, y varchar),
|     "f" timestamp(3) with time zone
| )
< ok - -
> exec
| CREATE TABLE IF NOT EXISTS "default"."MixedCase" (
|     "MixedCol" bigint NOT NULL,
|     "table" varchar
| )
< ok - -
