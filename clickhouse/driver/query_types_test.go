package driver

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/driver"
)

// typesSetup holds the values the client answers in its own Go types: a
// Decimal and a UUID, which it states as text but names as the column's scan
// type, and an Array, a Tuple and a Map, which it answers as Go slices and
// maps of the element types, its own among them. The second row holds each
// column's empty value or NULL.
const typesSetup = `
CREATE TABLE typed (
  id Int32,
  dec Decimal(9, 2),
  dec_n Nullable(Decimal(18, 4)),
  dec_wide Decimal(76, 10),
  u UUID,
  u_n Nullable(UUID),
  arr Array(Int32),
  arr_n Array(Nullable(String)),
  arr_nested Array(Array(UInt8)),
  arr_dec Array(Decimal(9, 3)),
  arr_uuid Array(UUID),
  arr_dt Array(DateTime('Asia/Tokyo')),
  arr_f Array(Float64),
  arr_enum Array(Enum8('a' = 1, 'b' = 2)),
  arr_low Array(LowCardinality(String)),
  tup Tuple(Int32, String),
  tup_named Tuple(x Float64, y Nullable(Int8)),
  tup_nested Tuple(Array(UInt16), Map(String, UUID)),
  m Map(String, UInt64),
  m_float Map(Float64, Array(String)),
  m_dec Map(String, Nullable(Decimal(9, 2))),
  i128 Int128,
  ip4 IPv4,
  ip6 IPv6,
  arr_ip Array(IPv4),
  pt Point
) ENGINE = Memory;

INSERT INTO typed VALUES
  (1, 1.50, -12.3400, '-12345678901234567890.0123456789',
   '0f9a2d7e-5b1c-4e3a-9d8f-7a6b5c4d3e2f', '00000000-0000-0000-0000-000000000001',
   [1, -2, 3], ['a', NULL, 'b,c'], [[0, 255], []], [1.5, -0.001],
   ['0f9a2d7e-5b1c-4e3a-9d8f-7a6b5c4d3e2f'], ['2026-03-04 05:06:07'], [1.5, nan, inf, -inf],
   ['b', 'a'], ['x', 'x'],
   (7, 'seven'), (2.5, NULL), ([1, 2], {'k': '00000000-0000-0000-0000-000000000002'}),
   {'b': 2, 'a': 1}, {2.5: ['x'], -1: []}, {'p': 1.25, 'q': NULL},
   -170141183460469231731687303715884105728, '127.0.0.1', '::ffff:10.0.0.1', ['192.168.0.1'],
   (1.5, -2)),
  (2, 0, NULL, 0,
   '00000000-0000-0000-0000-000000000000', NULL,
   [], [], [], [], [], [], [], [], [],
   (0, ''), (0, 0), ([], {}),
   {}, {}, {},
   0, '0.0.0.0', '::', [],
   (0, 0));
`

func TestQueryTypes(t *testing.T) {
	d := openTestDriver(t)
	execCorpus(t, d.db, typesSetup)
	results := query(t, d, "SELECT * FROM typed ORDER BY id", driver.QueryOptions{})
	if len(results) != 1 || results[0].Err != nil {
		t.Fatalf("results = %+v", results)
	}
	result := results[0]
	wantColumns := []driver.Column{
		{Name: "id", TypeName: "INT32"},
		{Name: "dec", TypeName: "DECIMAL(9, 2)"},
		{Name: "dec_n", TypeName: "NULLABLE(DECIMAL(18, 4))"},
		{Name: "dec_wide", TypeName: "DECIMAL(76, 10)"},
		{Name: "u", TypeName: "UUID"},
		{Name: "u_n", TypeName: "NULLABLE(UUID)"},
		{Name: "arr", TypeName: "ARRAY(INT32)"},
		{Name: "arr_n", TypeName: "ARRAY(NULLABLE(STRING))"},
		{Name: "arr_nested", TypeName: "ARRAY(ARRAY(UINT8))"},
		{Name: "arr_dec", TypeName: "ARRAY(DECIMAL(9, 3))"},
		{Name: "arr_uuid", TypeName: "ARRAY(UUID)"},
		{Name: "arr_dt", TypeName: "ARRAY(DATETIME('ASIA/TOKYO'))"},
		{Name: "arr_f", TypeName: "ARRAY(FLOAT64)"},
		{Name: "arr_enum", TypeName: "ARRAY(ENUM8('A' = 1, 'B' = 2))"},
		{Name: "arr_low", TypeName: "ARRAY(LOWCARDINALITY(STRING))"},
		{Name: "tup", TypeName: "TUPLE(INT32, STRING)"},
		{Name: "tup_named", TypeName: "TUPLE(X FLOAT64, Y NULLABLE(INT8))"},
		{Name: "tup_nested", TypeName: "TUPLE(ARRAY(UINT16), MAP(STRING, UUID))"},
		{Name: "m", TypeName: "MAP(STRING, UINT64)"},
		{Name: "m_float", TypeName: "MAP(FLOAT64, ARRAY(STRING))"},
		{Name: "m_dec", TypeName: "MAP(STRING, NULLABLE(DECIMAL(9, 2)))"},
		{Name: "i128", TypeName: "INT128"},
		{Name: "ip4", TypeName: "IPV4"},
		{Name: "ip6", TypeName: "IPV6"},
		{Name: "arr_ip", TypeName: "ARRAY(IPV4)"},
		{Name: "pt", TypeName: "POINT"},
	}
	if diff := cmp.Diff(wantColumns, result.Columns); diff != "" {
		t.Errorf("columns (-want +got):\n%s", diff)
	}
	str := func(s string) driver.Value { return driver.Value{Kind: driver.String, String: s} }
	null := driver.Value{Kind: driver.Null}
	want := [][]driver.Value{
		{
			{Kind: driver.Int32, Int64: 1},
			str("1.5"),
			str("-12.34"),
			str("-12345678901234567890.0123456789"),
			str("0f9a2d7e-5b1c-4e3a-9d8f-7a6b5c4d3e2f"),
			str("00000000-0000-0000-0000-000000000001"),
			str(`[1,-2,3]`),
			str(`["a",null,"b,c"]`),
			str(`[[0,255],[]]`),
			str(`["1.5","-0.001"]`),
			str(`["0f9a2d7e-5b1c-4e3a-9d8f-7a6b5c4d3e2f"]`),
			str(`["2026-03-04T05:06:07+09:00"]`),
			str(`[1.5,"nan","inf","-inf"]`),
			str(`["b","a"]`),
			str(`["x","x"]`),
			str(`[7,"seven"]`),
			str(`{"x":2.5,"y":null}`),
			str(`[[1,2],{"k":"00000000-0000-0000-0000-000000000002"}]`),
			str(`{"a":1,"b":2}`),
			str(`{"-1":[],"2.5":["x"]}`),
			str(`{"p":"1.25","q":null}`),
			str("-170141183460469231731687303715884105728"),
			str("127.0.0.1"),
			str("10.0.0.1"),
			str(`["192.168.0.1"]`),
			str("POINT(1.5 -2)"),
		},
		{
			{Kind: driver.Int32, Int64: 2},
			str("0"),
			null,
			str("0"),
			str("00000000-0000-0000-0000-000000000000"),
			null,
			str(`[]`),
			str(`[]`),
			str(`[]`),
			str(`[]`),
			str(`[]`),
			str(`[]`),
			str(`[]`),
			str(`[]`),
			str(`[]`),
			str(`[0,""]`),
			str(`{"x":0,"y":0}`),
			str(`[[],{}]`),
			str(`{}`),
			str(`{}`),
			str(`{}`),
			str("0"),
			str("0.0.0.0"),
			str("::"),
			str(`[]`),
			str("POINT(0 0)"),
		},
	}
	if diff := cmp.Diff(want, result.Rows); diff != "" {
		t.Errorf("rows (-want +got):\n%s", diff)
	}
}
