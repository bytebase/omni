-- pin: default cockroachdb/cockroach:v25.2.3 sha256:4a412731d1cfc5d6dea7431b7ce383b502966c9d589163d63db9e1ff7b51495a v25.2.3

== TestQuery default#1
> exec
|
| CREATE SCHEMA app;
| CREATE TABLE app.people (
|     id INT8 PRIMARY KEY,
|     name STRING NOT NULL,
|     active BOOL,
|     amount DECIMAL(10, 2),
|     ratio FLOAT8,
|     born DATE,
|     at TIMESTAMP,
|     at_tz TIMESTAMPTZ
| );
| INSERT INTO app.people VALUES
|     (1, 'one', true, 10.50, 1.5, '2020-01-02', '2024-01-01 00:00:00', '2024-01-01 00:00:00+00'),
|     (2, 'two', false, 20.25, 2.5, '2021-02-03', '2024-01-02 00:00:00', '2024-01-02 00:00:00+00');
| CREATE SCHEMA other;
| CREATE TABLE other.thing (v INT8);
| INSERT INTO other.thing VALUES (7);
< ok 1 -
> query
| SHOW timezone
< result
< column "timezone" "TEXT"
< row
: string
| UTC
> query
| SELECT id, name, active, amount, ratio, born, at, at_tz FROM app.people ORDER BY id;
< result
< column "id" "INT8"
< column "name" "TEXT"
< column "active" "BOOL"
< column "amount" "NUMERIC"
< column "ratio" "FLOAT8"
< column "born" "DATE"
< column "at" "TIMESTAMP"
< column "at_tz" "TIMESTAMPTZ"
< row
: int64 1
: string
| one
: bool true
: string
| 10.50
: float64 1.5
: time 2020-01-02T00:00:00Z
: time 2024-01-01T00:00:00Z
: time 2024-01-01T00:00:00Z
< row
: int64 2
: string
| two
: bool false
: string
| 20.25
: float64 2.5
: time 2021-02-03T00:00:00Z
: time 2024-01-02T00:00:00Z
: time 2024-01-02T00:00:00Z
> query
| SHOW timezone
< result
< column "timezone" "TEXT"
< row
: string
| UTC
> query
| SELECT id, name FROM app.people ORDER BY id LIMIT 1;
< result
< column "id" "INT8"
< column "name" "TEXT"
< row
: int64 1
: string
| one
> query
| SHOW timezone
< result
< column "timezone" "TEXT"
< row
: string
| UTC
> query
| EXPLAIN (DISTSQL) SELECT id FROM app.people;
< result
< column "info" "TEXT"
< row
: string
| distribution: local
< row
: string
| vectorized: true
< row
: string
|
< row
: string
| • scan
< row
: string
|   missing stats
< row
: string
|   table: people@people_pkey
< row
: string
|   spans: FULL SCAN
< row
: string
|
< row
: string
| Diagram: https://cockroachdb.github.io/distsqlplan/decode.html#eJyMkF3r0zAYxe_9FOVcKWTa3uZK-W9CYf9trr0QpIzYPJthaRPzpMxR-t0l7cQXFEwgL-c5yfklI_irhcTm42H7rtxlL9dlVVcftq-yarPdPNWZ0dn74_45U96_9uS8JQj0TtNOdcSQn1CgEfDBtcTsQpLG2VDqb5C5gOn9EJPcCLQuEOSIaKIlSNTqs6UjKU3hTQ4BTVEZm8xYwt4u08lf6Q6BJ2eHrmeZGQ2Byqu0XqGZBNwQf-ZwVBeCLH4BK9eQ-ST-n-1I7F3P9BvWv5LyP5JWxdQIkL7MfzSC3RBaOgTXzt5lu5-RZ0ETx6VapDiOZf-jxDGQ6hb8RuBs3e1kNCTyR1v9ZXj0AumAunB6WPXF3eZr67tPWGdlmQSe1ZXWFCl0pjccTQt5VpZpml58HwDV-K2t
> exec
| SET search_path TO "other";
< ok 0 -
> query
| SHOW timezone
< result
< column "timezone" "TEXT"
< row
: string
| UTC
> query
| SELECT v FROM thing;
< result
< column "v" "INT8"
< row
: int64 7
> query
| SHOW timezone
< result
< column "timezone" "TEXT"
< row
: string
| UTC
> query
| SELECT 'abc'::BYTES AS b, ''::BYTES AS empty;
< result
< column "b" "BYTEA"
< column "empty" "BYTEA"
< row
: bytes
| abc
: bytes
|
> exec
| UPDATE app.people SET name = name WHERE id = 1;
< ok 1 -
> query
| SHOW timezone
< result
< column "timezone" "TEXT"
< row
: string
| UTC
> query
| INSERT INTO app.people (id, name) VALUES (3, 'three') RETURNING id, name;
< result
< column "id" "INT8"
< column "name" "TEXT"
< row
: int64 3
: string
| three
> query
| SHOW timezone
< result
< column "timezone" "TEXT"
< row
: string
| UTC
> query
| WITH d AS (DELETE FROM app.people WHERE id = 3 RETURNING id) SELECT count(*) FROM d;
< result
< column "count" "INT8"
< row
: int64 1
> query
| SHOW timezone
< result
< column "timezone" "TEXT"
< row
: string
| UTC
> query
| SELECT 1;
< result
< column "?column?" "INT8"
< row
: int64 1
> query
| SHOW timezone
< result
< column "timezone" "TEXT"
< row
: string
| UTC
> query
| SELECT * FROM app.missing;
< error
| 42P01 relation "app.missing" does not exist
> exec
| SET TIME ZONE 'America/New_York';
< ok 0 -
> query
| SHOW timezone
< result
< column "timezone" "TEXT"
< row
: string
| America/New_York
> query
| SELECT at_tz FROM app.people WHERE id = 1;
< result
< column "at_tz" "TIMESTAMPTZ"
< row
: time 2024-01-01T00:00:00Z
