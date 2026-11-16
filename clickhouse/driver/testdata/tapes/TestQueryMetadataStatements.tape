-- pin: default clickhouse/clickhouse-server:26.8.38.9 sha256:5d30ac138948b43fd928842be7da1dcd2de712d32ba6c42368ff5009ba2072ff 26.8.38.9

== TestQueryMetadataStatements default#1
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
| SHOW TABLES;
< result
< column "name" "String" "string"
< row
: string
| q
> query
|
| DESCRIBE TABLE q;
< result
< column "name" "String" "string"
< column "type" "String" "string"
< column "default_type" "String" "string"
< column "default_expression" "String" "string"
< column "comment" "String" "string"
< column "codec_expression" "String" "string"
< column "ttl_expression" "String" "string"
< row
: string
| id
: string
| Int32
: string
|
: string
|
: string
|
: string
|
: string
|
< row
: string
| i8
: string
| Int8
: string
|
: string
|
: string
|
: string
|
: string
|
< row
: string
| i16
: string
| Int16
: string
|
: string
|
: string
|
: string
|
: string
|
< row
: string
| i32
: string
| Int32
: string
|
: string
|
: string
|
: string
|
: string
|
< row
: string
| i64
: string
| Int64
: string
|
: string
|
: string
|
: string
|
: string
|
< row
: string
| u8
: string
| UInt8
: string
|
: string
|
: string
|
: string
|
: string
|
< row
: string
| u16
: string
| UInt16
: string
|
: string
|
: string
|
: string
|
: string
|
< row
: string
| u32
: string
| UInt32
: string
|
: string
|
: string
|
: string
|
: string
|
< row
: string
| u64
: string
| UInt64
: string
|
: string
|
: string
|
: string
|
: string
|
< row
: string
| f32
: string
| Float32
: string
|
: string
|
: string
|
: string
|
: string
|
< row
: string
| f64
: string
| Float64
: string
|
: string
|
: string
|
: string
|
: string
|
< row
: string
| s
: string
| String
: string
|
: string
|
: string
|
: string
|
: string
|
< row
: string
| b
: string
| Bool
: string
|
: string
|
: string
|
: string
|
: string
|
< row
: string
| d
: string
| Date
: string
|
: string
|
: string
|
: string
|
: string
|
< row
: string
| dt
: string
| DateTime('UTC')
: string
|
: string
|
: string
|
: string
|
: string
|
< row
: string
| dt64
: string
| DateTime64(3, 'UTC')
: string
|
: string
|
: string
|
: string
|
: string
|
< row
: string
| n
: string
| Nullable(Int32)
: string
|
: string
|
: string
|
: string
|
: string
|
< row
: string
| fs
: string
| FixedString(4)
: string
|
: string
|
: string
|
: string
|
: string
|
> query
|
| EXISTS TABLE q
< result
< column "result" "UInt8" "uint8"
< row
: uint8 1
