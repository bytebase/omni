-- pin: default clickhouse/clickhouse-server:26.8.38.9 sha256:5d30ac138948b43fd928842be7da1dcd2de712d32ba6c42368ff5009ba2072ff 26.8.38.9

== TestQueryMaxResultBytes default#1
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
| SELECT * FROM q ORDER BY id
< result
< column "id" "Int32" "int32"
< column "i8" "Int8" "int8"
< column "i16" "Int16" "int16"
< column "i32" "Int32" "int32"
< column "i64" "Int64" "int64"
< column "u8" "UInt8" "uint8"
< column "u16" "UInt16" "uint16"
< column "u32" "UInt32" "uint32"
< column "u64" "UInt64" "uint64"
< column "f32" "Float32" "float32"
< column "f64" "Float64" "float64"
< column "s" "String" "string"
< column "b" "Bool" "bool"
< column "d" "Date" "time.Time"
< column "dt" "DateTime('UTC')" "time.Time"
< column "dt64" "DateTime64(3, 'UTC')" "time.Time"
< column "n" "Nullable(Int32)" "*int32"
< column "fs" "FixedString(4)" "string"
< row
: int32 1
: int8 -8
: int16 -16
: int32 -32
: int64 -64
: uint8 8
: uint16 16
: uint32 32
: uint64 64
: float32 1.5
: float64 2.5
: string
| a
: bool true
: time 2026-03-04T00:00:00Z
: time 2026-03-04T05:06:07Z
: time 2026-03-04T05:06:07.891Z
: int32 42
: string
| abcd
< row
: int32 2
: int8 0
: int16 0
: int32 0
: int64 0
: uint8 0
: uint16 0
: uint32 0
: uint64 0
: float32 0
: float64 0
: string
|
: bool false
: time 1970-01-01T00:00:00Z
: time 1970-01-01T00:00:00Z
: time 1970-01-01T00:00:00Z
: null
: string
| wxyz
