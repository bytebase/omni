-- pin: default clickhouse/clickhouse-server:26.8.38.9 sha256:5d30ac138948b43fd928842be7da1dcd2de712d32ba6c42368ff5009ba2072ff 26.8.38.9

== TestQueryTypes default#1
> exec
|
| CREATE TABLE typed (
|   id Int32,
|   dec Decimal(9, 2),
|   dec_n Nullable(Decimal(18, 4)),
|   dec_wide Decimal(76, 10),
|   u UUID,
|   u_n Nullable(UUID),
|   arr Array(Int32),
|   arr_n Array(Nullable(String)),
|   arr_nested Array(Array(UInt8)),
|   arr_dec Array(Decimal(9, 3)),
|   arr_uuid Array(UUID),
|   arr_dt Array(DateTime('Asia/Tokyo')),
|   arr_f Array(Float64),
|   arr_enum Array(Enum8('a' = 1, 'b' = 2)),
|   arr_low Array(LowCardinality(String)),
|   tup Tuple(Int32, String),
|   tup_named Tuple(x Float64, y Nullable(Int8)),
|   tup_nested Tuple(Array(UInt16), Map(String, UUID)),
|   m Map(String, UInt64),
|   m_float Map(Float64, Array(String)),
|   m_dec Map(String, Nullable(Decimal(9, 2))),
|   i128 Int128,
|   ip4 IPv4,
|   ip6 IPv6,
|   arr_ip Array(IPv4),
|   pt Point
| ) ENGINE = Memory;
< ok 0 -
> exec
|
|
| INSERT INTO typed VALUES
|   (1, 1.50, -12.3400, '-12345678901234567890.0123456789',
|    '0f9a2d7e-5b1c-4e3a-9d8f-7a6b5c4d3e2f', '00000000-0000-0000-0000-000000000001',
|    [1, -2, 3], ['a', NULL, 'b,c'], [[0, 255], []], [1.5, -0.001],
|    ['0f9a2d7e-5b1c-4e3a-9d8f-7a6b5c4d3e2f'], ['2026-03-04 05:06:07'], [1.5, nan, inf, -inf],
|    ['b', 'a'], ['x', 'x'],
|    (7, 'seven'), (2.5, NULL), ([1, 2], {'k': '00000000-0000-0000-0000-000000000002'}),
|    {'b': 2, 'a': 1}, {2.5: ['x'], -1: []}, {'p': 1.25, 'q': NULL},
|    -170141183460469231731687303715884105728, '127.0.0.1', '::ffff:10.0.0.1', ['192.168.0.1'],
|    (1.5, -2)),
|   (2, 0, NULL, 0,
|    '00000000-0000-0000-0000-000000000000', NULL,
|    [], [], [], [], [], [], [], [], [],
|    (0, ''), (0, 0), ([], {}),
|    {}, {}, {},
|    0, '0.0.0.0', '::', [],
|    (0, 0));
< ok 0 -
> query
| SELECT * FROM typed ORDER BY id
< result
< column "id" "Int32" "int32"
< column "dec" "Decimal(9, 2)" "decimal.Decimal" precision 9 scale 2
< column "dec_n" "Nullable(Decimal(18, 4))" "*decimal.Decimal" precision 18 scale 4
< column "dec_wide" "Decimal(76, 10)" "decimal.Decimal" precision 76 scale 10
< column "u" "UUID" "uuid.UUID"
< column "u_n" "Nullable(UUID)" "*uuid.UUID"
< column "arr" "Array(Int32)" "[]int32"
< column "arr_n" "Array(Nullable(String))" "[]*string"
< column "arr_nested" "Array(Array(UInt8))" "[][]uint8"
< column "arr_dec" "Array(Decimal(9, 3))" "[]decimal.Decimal" precision 9 scale 3
< column "arr_uuid" "Array(UUID)" "[]uuid.UUID"
< column "arr_dt" "Array(DateTime('Asia/Tokyo'))" "[]time.Time"
< column "arr_f" "Array(Float64)" "[]float64"
< column "arr_enum" "Array(Enum8('a' = 1, 'b' = 2))" "[]string"
< column "arr_low" "Array(LowCardinality(String))" "[]string"
< column "tup" "Tuple(Int32, String)" "[]interface {}"
< column "tup_named" "Tuple(x Float64, y Nullable(Int8))" "map[string]interface {}"
< column "tup_nested" "Tuple(Array(UInt16), Map(String, UUID))" "[]interface {}"
< column "m" "Map(String, UInt64)" "map[string]uint64"
< column "m_float" "Map(Float64, Array(String))" "map[float64][]string"
< column "m_dec" "Map(String, Nullable(Decimal(9, 2)))" "map[string]*decimal.Decimal"
< column "i128" "Int128" "*big.Int"
< column "ip4" "IPv4" "net.IP"
< column "ip6" "IPv6" "net.IP"
< column "arr_ip" "Array(IPv4)" "[]net.IP"
< column "pt" "Point" "orb.Point"
< row
: int32 1
: string
| 1.5
: string
| -12.34
: string
| -12345678901234567890.0123456789
: string
| 0f9a2d7e-5b1c-4e3a-9d8f-7a6b5c4d3e2f
: string
| 00000000-0000-0000-0000-000000000001
: go []int32
| [1, -2, 3]
: go []*string
| [&"a", nil, &"b,c"]
: go [][]uint8
| [[0, 255], []]
: go []decimal.Decimal
| ["1500e-3", "-1e-3"]
: go []uuid.UUID
| ["0f9a2d7e-5b1c-4e3a-9d8f-7a6b5c4d3e2f"]
: go []time.Time
| ["2026-03-04T05:06:07+09:00 Asia/Tokyo"]
: go []float64
| [1.5, NaN, +Inf, -Inf]
: go []string
| ["b", "a"]
: go []string
| ["x", "x"]
: go []interface {}
| [int32(7), string("seven")]
: go map[string]interface {}
| {"x": float64(2.5), "y": *int8(nil)}
: go []interface {}
| [[]uint16([1, 2]), map[string]uuid.UUID({"k": "00000000-0000-0000-0000-000000000002"})]
: go map[string]uint64
| {"a": 1, "b": 2}
: go map[float64][]string
| {-1: [], 2.5: ["x"]}
: go map[string]*decimal.Decimal
| {"p": &"125e-2", "q": nil}
: go big.Int
| "-170141183460469231731687303715884105728"
: go net.IP
| [127, 0, 0, 1]
: go net.IP
| [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 10, 0, 0, 1]
: go []net.IP
| [[192, 168, 0, 1]]
: go orb.Point
| [1.5, -2]
< row
: int32 2
: string
| 0
: null
: string
| 0
: string
| 00000000-0000-0000-0000-000000000000
: null
: go []int32
| []
: go []*string
| []
: go [][]uint8
| []
: go []decimal.Decimal
| []
: go []uuid.UUID
| []
: go []time.Time
| []
: go []float64
| []
: go []string
| []
: go []string
| []
: go []interface {}
| [int32(0), string("")]
: go map[string]interface {}
| {"x": float64(0), "y": *int8(&0)}
: go []interface {}
| [[]uint16([]), map[string]uuid.UUID({})]
: go map[string]uint64
| {}
: go map[float64][]string
| {}
: go map[string]*decimal.Decimal
| {}
: go big.Int
| "0"
: go net.IP
| [0, 0, 0, 0]
: go net.IP
| [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
: go []net.IP
| []
: go orb.Point
| [0, 0]
