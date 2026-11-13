-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestTypeNameMatchesFormatType default#1
> query
| SELECT format_type($1::text::regtype, NULL)
: string
| int
< result
< column "format_type" "TEXT"
< row
: string
| integer
> query ^
: string
| integer
< result
< column "format_type" "TEXT"
< row
: string
| integer
> query ^
: string
| pg_catalog.int4
< result
< column "format_type" "TEXT"
< row
: string
| integer
> query ^
: string
| smallint
< result
< column "format_type" "TEXT"
< row
: string
| smallint
> query ^
: string
| bigint
< result
< column "format_type" "TEXT"
< row
: string
| bigint
> query ^
: string
| boolean
< result
< column "format_type" "TEXT"
< row
: string
| boolean
> query ^
: string
| pg_catalog.bool
< result
< column "format_type" "TEXT"
< row
: string
| boolean
> query ^
: string
| real
< result
< column "format_type" "TEXT"
< row
: string
| real
> query ^
: string
| float
< result
< column "format_type" "TEXT"
< row
: string
| double precision
> query ^
: string
| float(10)
< result
< column "format_type" "TEXT"
< row
: string
| real
> query ^
: string
| float(30)
< result
< column "format_type" "TEXT"
< row
: string
| double precision
> query ^
: string
| double precision
< result
< column "format_type" "TEXT"
< row
: string
| double precision
> query ^
: string
| pg_catalog.float8
< result
< column "format_type" "TEXT"
< row
: string
| double precision
> query ^
: string
| numeric
< result
< column "format_type" "TEXT"
< row
: string
| numeric
> query ^
: string
| numeric(10,2)
< result
< column "format_type" "TEXT"
< row
: string
| numeric
> query ^
: string
| decimal
< result
< column "format_type" "TEXT"
< row
: string
| numeric
> query ^
: string
| dec(5)
< result
< column "format_type" "TEXT"
< row
: string
| numeric
> query ^
: string
| varchar
< result
< column "format_type" "TEXT"
< row
: string
| character varying
> query ^
: string
| varchar(10)
< result
< column "format_type" "TEXT"
< row
: string
| character varying
> query ^
: string
| character varying(3)
< result
< column "format_type" "TEXT"
< row
: string
| character varying
> query ^
: string
| char varying
< result
< column "format_type" "TEXT"
< row
: string
| character varying
> query ^
: string
| pg_catalog.varchar
< result
< column "format_type" "TEXT"
< row
: string
| character varying
> query ^
: string
| character
< result
< column "format_type" "TEXT"
< row
: string
| character
> query ^
: string
| char
< result
< column "format_type" "TEXT"
< row
: string
| character
> query ^
: string
| char(3)
< result
< column "format_type" "TEXT"
< row
: string
| character
> query ^
: string
| nchar(2)
< result
< column "format_type" "TEXT"
< row
: string
| character
> query ^
: string
| national character varying(4)
< result
< column "format_type" "TEXT"
< row
: string
| character varying
> query ^
: string
| pg_catalog.bpchar
< result
< column "format_type" "TEXT"
< row
: string
| character
> query ^
: string
| bit
< result
< column "format_type" "TEXT"
< row
: string
| bit
> query ^
: string
| bit(3)
< result
< column "format_type" "TEXT"
< row
: string
| bit
> query ^
: string
| bit varying
< result
< column "format_type" "TEXT"
< row
: string
| bit varying
> query ^
: string
| pg_catalog.varbit
< result
< column "format_type" "TEXT"
< row
: string
| bit varying
> query ^
: string
| time
< result
< column "format_type" "TEXT"
< row
: string
| time without time zone
> query ^
: string
| time(3)
< result
< column "format_type" "TEXT"
< row
: string
| time without time zone
> query ^
: string
| time with time zone
< result
< column "format_type" "TEXT"
< row
: string
| time with time zone
> query ^
: string
| time(2) without time zone
< result
< column "format_type" "TEXT"
< row
: string
| time without time zone
> query ^
: string
| pg_catalog.timetz
< result
< column "format_type" "TEXT"
< row
: string
| time with time zone
> query ^
: string
| timestamp
< result
< column "format_type" "TEXT"
< row
: string
| timestamp without time zone
> query ^
: string
| timestamp(6) with time zone
< result
< column "format_type" "TEXT"
< row
: string
| timestamp with time zone
> query ^
: string
| timestamp without time zone
< result
< column "format_type" "TEXT"
< row
: string
| timestamp without time zone
> query ^
: string
| pg_catalog.timestamptz
< result
< column "format_type" "TEXT"
< row
: string
| timestamp with time zone
> query ^
: string
| interval
< result
< column "format_type" "TEXT"
< row
: string
| interval
> query ^
: string
| interval year to month
< result
< column "format_type" "TEXT"
< row
: string
| interval
> query ^
: string
| interval(3)
< result
< column "format_type" "TEXT"
< row
: string
| interval
> query ^
: string
| json
< result
< column "format_type" "TEXT"
< row
: string
| json
> query ^
: string
| jsonb
< result
< column "format_type" "TEXT"
< row
: string
| jsonb
> query ^
: string
| text
< result
< column "format_type" "TEXT"
< row
: string
| text
> query ^
: string
| pg_catalog.text
< result
< column "format_type" "TEXT"
< row
: string
| text
> query ^
: string
| "char"
< result
< column "format_type" "TEXT"
< row
: string
| "char"
> query ^
: string
| name
< result
< column "format_type" "TEXT"
< row
: string
| name
> query ^
: string
| uuid
< result
< column "format_type" "TEXT"
< row
: string
| uuid
> query ^
: string
| bytea
< result
< column "format_type" "TEXT"
< row
: string
| bytea
> query ^
: string
| xml
< result
< column "format_type" "TEXT"
< row
: string
| xml
> query ^
: string
| int[]
< result
< column "format_type" "TEXT"
< row
: string
| integer[]
> query ^
: string
| integer[][]
< result
< column "format_type" "TEXT"
< row
: string
| integer[]
> query ^
: string
| int ARRAY
< result
< column "format_type" "TEXT"
< row
: string
| integer[]
> query ^
: string
| int ARRAY[3]
< result
< column "format_type" "TEXT"
< row
: string
| integer[]
> query ^
: string
| text[]
< result
< column "format_type" "TEXT"
< row
: string
| text[]
> query ^
: string
| double precision[]
< result
< column "format_type" "TEXT"
< row
: string
| double precision[]
> query ^
: string
| pg_catalog.int8[]
< result
< column "format_type" "TEXT"
< row
: string
| bigint[]
> exec
| CREATE SCHEMA s
< ok 0 -
> exec
| CREATE TYPE s.int2 AS (a int)
< ok 0 -
> exec
| CREATE TYPE s.int4 AS (a int)
< ok 0 -
> exec
| CREATE TYPE s.int8 AS (a int)
< ok 0 -
> exec
| CREATE TYPE s.bool AS (a int)
< ok 0 -
> exec
| CREATE TYPE s.float4 AS (a int)
< ok 0 -
> exec
| CREATE TYPE s.float8 AS (a int)
< ok 0 -
> exec
| CREATE TYPE s.bpchar AS (a int)
< ok 0 -
> exec
| CREATE TYPE s.varbit AS (a int)
< ok 0 -
> exec
| CREATE TYPE s.timetz AS (a int)
< ok 0 -
> exec
| CREATE TYPE s.timestamptz AS (a int)
< ok 0 -
> exec
| CREATE TYPE s."varchar" AS (a int)
< ok 0 -
> exec
| CREATE TYPE s."numeric" AS (a int)
< ok 0 -
> exec
| CREATE TYPE s."bit" AS (a int)
< ok 0 -
> exec
| CREATE TYPE s."time" AS (a int)
< ok 0 -
> exec
| CREATE TYPE s."timestamp" AS (a int)
< ok 0 -
> exec
| CREATE TYPE s."interval" AS (a int)
< ok 0 -
> exec
| CREATE TYPE s."json" AS (a int)
< ok 0 -
> exec
| SET search_path = s, pg_catalog
< ok 0 -
> query
| SELECT format_type($1::text::regtype, NULL)
: string
| int2
< result
< column "format_type" "TEXT"
< row
: string
| int2
> query ^
: string
| int4
< result
< column "format_type" "TEXT"
< row
: string
| int4
> query ^
: string
| int8
< result
< column "format_type" "TEXT"
< row
: string
| int8
> query ^
: string
| bool
< result
< column "format_type" "TEXT"
< row
: string
| bool
> query ^
: string
| float4
< result
< column "format_type" "TEXT"
< row
: string
| float4
> query ^
: string
| float8
< result
< column "format_type" "TEXT"
< row
: string
| float8
> query ^
: string
| bpchar
< result
< column "format_type" "TEXT"
< row
: string
| bpchar
> query ^
: string
| varbit
< result
< column "format_type" "TEXT"
< row
: string
| varbit
> query ^
: string
| timetz
< result
< column "format_type" "TEXT"
< row
: string
| timetz
> query ^
: string
| timestamptz
< result
< column "format_type" "TEXT"
< row
: string
| timestamptz
> query ^
: string
| "varchar"
< result
< column "format_type" "TEXT"
< row
: string
| "varchar"
> query ^
: string
| "numeric"
< result
< column "format_type" "TEXT"
< row
: string
| "numeric"
> query ^
: string
| "bit"
< result
< column "format_type" "TEXT"
< row
: string
| "bit"
> query ^
: string
| "time"
< result
< column "format_type" "TEXT"
< row
: string
| "time"
> query ^
: string
| "timestamp"
< result
< column "format_type" "TEXT"
< row
: string
| "timestamp"
> query ^
: string
| "interval"
< result
< column "format_type" "TEXT"
< row
: string
| "interval"
> query ^
: string
| "json"
< result
< column "format_type" "TEXT"
< row
: string
| "json"
> query ^
: string
| int4[]
< result
< column "format_type" "TEXT"
< row
: string
| int4[]
> query ^
: string
| timestamptz[]
< result
< column "format_type" "TEXT"
< row
: string
| timestamptz[]
> query ^
: string
| integer
< result
< column "format_type" "TEXT"
< row
: string
| integer
> query ^
: string
| character varying(3)
< result
< column "format_type" "TEXT"
< row
: string
| character varying
> query ^
: string
| timestamp with time zone
< result
< column "format_type" "TEXT"
< row
: string
| timestamp with time zone
> query ^
: string
| bit varying
< result
< column "format_type" "TEXT"
< row
: string
| bit varying
> query ^
: string
| pg_catalog.int4
< result
< column "format_type" "TEXT"
< row
: string
| integer
