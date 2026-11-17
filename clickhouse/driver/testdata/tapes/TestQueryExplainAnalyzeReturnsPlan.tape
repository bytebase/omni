-- pin: default clickhouse/clickhouse-server:26.8.38.9 sha256:5d30ac138948b43fd928842be7da1dcd2de712d32ba6c42368ff5009ba2072ff 26.8.38.9

== TestQueryExplainAnalyzeReturnsPlan default#1
> exec
|
| CREATE TABLE q (
|   id Int32,
|   i8 Int8,
|   i16 Int16,
|   i32 Int32,
|   i64 Int64,
|   u8 UInt8,
|   u16 UInt16,
|   u32 UInt32,
|   u64 UInt64,
|   f32 Float32,
|   f64 Float64,
|   s String,
|   b Bool,
|   d Date,
|   dt DateTime('UTC'),
|   dt64 DateTime64(3, 'UTC'),
|   n Nullable(Int32),
|   fs FixedString(4)
| ) ENGINE = Memory;
< ok 0 -
> exec
|
|
| INSERT INTO q VALUES
|   (1, -8, -16, -32, -64, 8, 16, 32, 64, 1.5, 2.5, 'a', true,
|    '2026-03-04', '2026-03-04 05:06:07', '2026-03-04 05:06:07.891', 42, 'abcd'),
|   (2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, '', false,
|    '1970-01-01', '1970-01-01 00:00:00', '1970-01-01 00:00:00', NULL, 'wxyz');
< ok 0 -
> query
| EXPLAIN ANALYZE SELECT count() FROM q
< result
< column "explain" "String" "string"
< row
: string
| Query summary:
< row
: string
|   Time:        811.68 us (planning 515.96 us · execution 295.72 us)
< row
: string
|   Read:        2.00 rows, 2.00 B (6.76 thousand rows/s., 6.76 KB/s.)
< row
: string
|   Peak memory: 31.48 KiB
< row
: string
|
< row
: string
| Output: count()
< row
: string
|
< row
: string
| Expression ((Project names + Projection))
< row
: string
| │  I/O: rows 1 → 1 · 8 B → 8 B
< row
: string
| │    time 868.00 ns (0.3%) · parallelism 0.94/1
< row
: string
| └──Aggregating
< row
: string
|    │  Keys:
< row
: string
|    │  Aggregates: count()
< row
: string
|    │  Skip merging: 0
< row
: string
|    │  I/O: rows 2 → 1 (50.00%) · 0 B → 8 B
< row
: string
|    │    Stage (partial aggregation): time 47.38 us (16.0%) · parallelism 1.00/1
< row
: string
|    │    Stage (final aggregation): time 4.18 us (1.4%) · parallelism 0.98/1
< row
: string
|    └──Expression ((Before GROUP BY + Change column names to column identifiers))
< row
: string
|       │  I/O: rows 2 → 2 · 2 B → 0 B
< row
: string
|       │    time 2.70 us (0.9%) · parallelism 0.98/1
< row
: string
|       └──ReadFromMemoryStorage
< row
: string
|             Output: i8
< row
: string
|             I/O: rows 0 → 2 · 0 B → 2 B
< row
: string
|               time 22.92 us (7.8%) · parallelism 0.99/1
