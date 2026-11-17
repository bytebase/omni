-- pin: default cockroachdb/cockroach:v25.2.3 sha256:4a412731d1cfc5d6dea7431b7ce383b502966c9d589163d63db9e1ff7b51495a v25.2.3

== TestCountAffectedRowsLive default#1
> exec
|
| CREATE SCHEMA app;
| CREATE TABLE app.big (id INT8 PRIMARY KEY);
| INSERT INTO app.big SELECT generate_series(1, 100);
| ANALYZE app.big;
< ok 0 -
> query
| EXPLAIN UPDATE app.big SET id = id WHERE id <= 50
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
| • update
< row
: string
| │ table: big
< row
: string
| │ set: id
< row
: string
| │ auto commit
< row
: string
| │
< row
: string
| └── • scan
< row
: string
|       estimated row count: 50 (50% of the table; stats collected 0 seconds ago)
< row
: string
|       table: big@big_pkey
< row
: string
|       spans: [ - /50]
< row
: string
|       locking strength: for update
> query
| EXPLAIN DELETE FROM app.big WHERE id <= 50
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
| • delete range
< row
: string
|   from: big
< row
: string
|   spans: [ - /50]
> begin
< ok
> exec
| SET LOCAL optimizer_use_delete_range_fast_path = off
< ok 0 -
> query
| EXPLAIN DELETE FROM app.big WHERE id <= 50
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
| • delete
< row
: string
| │ from: big
< row
: string
| │
< row
: string
| └── • scan
< row
: string
|       estimated row count: 50 (50% of the table; stats collected 0 seconds ago)
< row
: string
|       table: big@big_pkey
< row
: string
|       spans: [ - /50]
< row
: string
|       locking strength: for update
> rollback
< ok
