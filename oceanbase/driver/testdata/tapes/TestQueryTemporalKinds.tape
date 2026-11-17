-- pin: default oceanbase/oceanbase-ce:4.3.5.6-106000012026040916 sha256:31086a6900c21c479c2bcd942b6a28c53b17a51f4e9b9eb8eafcc596adfcd2e3 v4.3.5_CE_BP6

== TestQueryTemporalKinds/default-loc default#1
> exec
| DROP TABLE IF EXISTS qtimes;
< ok 0 0
> exec
|
| CREATE TABLE qtimes (
|   id BIGINT PRIMARY KEY,
|   ts TIMESTAMP(6) NULL,
|   dt DATETIME(6) NULL,
|   dt0 DATETIME NULL,
|   d DATE NULL
| );
< ok 0 0
> exec
|
| INSERT INTO qtimes VALUES
|   (1, '2026-03-04 05:06:07.123456', '2026-03-04 05:06:07.123456', '2026-03-04 05:06:07', '2026-03-04'),
|   (2, NULL, NULL, NULL, NULL);
< ok 2 0
> query
| SELECT * FROM qtimes ORDER BY id;
< result
< column "id" "BIGINT" "int64"
< column "ts" "TIMESTAMP" precision 6 scale 6
< column "dt" "DATETIME" precision 6 scale 6
< column "dt0" "DATETIME" precision 0 scale 0
< column "d" "DATE"
< row
: int64 1
: time 2026-03-04T05:06:07.123456Z
: time 2026-03-04T05:06:07.123456Z
: time 2026-03-04T05:06:07Z
: time 2026-03-04T00:00:00Z
< row
: int64 2
: null
: null
: null
: null

== TestQueryTemporalKinds/loc-plus8 default#1
> exec
| DROP TABLE IF EXISTS qtimes;
< ok 0 0
> exec
|
| CREATE TABLE qtimes (
|   id BIGINT PRIMARY KEY,
|   ts TIMESTAMP(6) NULL,
|   dt DATETIME(6) NULL,
|   dt0 DATETIME NULL,
|   d DATE NULL
| );
< ok 0 0
> exec
|
| INSERT INTO qtimes VALUES
|   (1, '2026-03-04 05:06:07.123456', '2026-03-04 05:06:07.123456', '2026-03-04 05:06:07', '2026-03-04'),
|   (2, NULL, NULL, NULL, NULL);
< ok 2 0
> query
| SELECT * FROM qtimes ORDER BY id;
< result
< column "id" "BIGINT" "int64"
< column "ts" "TIMESTAMP" precision 6 scale 6
< column "dt" "DATETIME" precision 6 scale 6
< column "dt0" "DATETIME" precision 0 scale 0
< column "d" "DATE"
< row
: int64 1
: time 2026-03-04T05:06:07.123456+08:00
: time 2026-03-04T05:06:07.123456+08:00
: time 2026-03-04T05:06:07+08:00
: time 2026-03-04T00:00:00+08:00
< row
: int64 2
: null
: null
: null
: null
