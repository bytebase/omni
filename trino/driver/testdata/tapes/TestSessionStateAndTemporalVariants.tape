-- pin: default trinodb/trino:482 sha256:90b35b7c603eaa1f889bf03981a62b75f998ee6c0f851d9f4e341b49a57022b6 482

== TestSessionStateAndTemporalVariants/session_persists_SET_SESSION_across_Query_calls default#1
> exec
| SET SESSION query_max_run_time = '1m'
< ok 0 -
> query
| SHOW SESSION LIKE 'query_max_run_time'
< result
< column "Name" "VARCHAR"
< column "Value" "VARCHAR"
< column "Default" "VARCHAR"
< column "Type" "VARCHAR"
< column "Description" "VARCHAR"
< row
: string
| query_max_run_time
: string
| 1m
: string
| 100.00d
: string
| varchar
: string
| Maximum run time of a query (includes the queueing time)

== TestSessionStateAndTemporalVariants/temporal_precision_and_null_variants default#1
> query
| SELECT
| 			TIME '12:34:56' AS c_time_sec,
| 			TIME '12:34:56.789 +02:00' AS c_time_tz_ms,
| 			TIMESTAMP '2026-01-02 12:34:56.789' AS c_timestamp_ms,
| 			TIMESTAMP '2026-01-02 12:34:56 +02:00' AS c_timestamp_tz_sec,
| 			CAST(NULL AS DATE) AS c_null_date
< result
< column "c_time_sec" "TIME" precision 0 scale 0
< column "c_time_tz_ms" "TIME WITH TIME ZONE" precision 3 scale 0
< column "c_timestamp_ms" "TIMESTAMP" precision 3 scale 0
< column "c_timestamp_tz_sec" "TIMESTAMP WITH TIME ZONE" precision 0 scale 0
< column "c_null_date" "DATE"
< row
: time 0000-01-01T12:34:56Z
: time 0000-01-01T12:34:56.789+02:00
: time 2026-01-02T12:34:56.789Z
: time 2026-01-02T12:34:56+02:00
: null
