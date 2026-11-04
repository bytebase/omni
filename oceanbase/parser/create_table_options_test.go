package parser

import "testing"

// The CREATE TABLE table options and column-type oddities the TiDB grammar gap
// inventory groups under "CREATE TABLE clauses"
// (docs/plans/2026-11-03-tidb-grammar-gap-inventory.md), read against OceanBase
// rather than ported from the MySQL tree this one forked. OceanBase turns out to
// have far less of this corner of the grammar than MySQL does, so most of these
// rows are rejections and the parser owes none of them:
//
//   - TABLE_CHECKSUM is not an option name, in any spelling. CHECKSUM is, and
//     takes a number: `CHECKSUM = 'abc'` and `= DEFAULT` are each a 1064.
//   - STATS_AUTO_RECALC takes any integer, not MySQL's 0, 1 and DEFAULT alone,
//     and still refuses a string.
//   - The STATS_ counters TiDB adds — STATS_BUCKETS, STATS_TOPN,
//     STATS_SAMPLE_RATE, STATS_COL_CHOICE, STATS_COL_LIST — are not options
//     here.
//   - SERIAL DEFAULT VALUE is not a column option, although SERIAL is a type.
//   - SQL_TSI_YEAR is not a type name, and YEAR takes neither UNSIGNED nor
//     ZEROFILL.
//   - SECONDARY_ENGINE, ENGINE_ATTRIBUTE and SECONDARY_ENGINE_ATTRIBUTE are not
//     options, as a table option or as a column option. The two attributes are
//     no longer read here; SECONDARY_ENGINE still is, which is an
//     over-acceptance of its own and left to the PR that owns that option.
//   - GLOBAL TEMPORARY and ON COMMIT belong to TiDB alone.
//
// An ENUM or SET member written as a hex or bit literal, and a name that begins
// with a digit, are the two things in this family OceanBase does have.
//
// The engine column is what 5.7.25-OceanBase_CE-v4.3.5.6 answered on 2026-11-04,
// read with scripts/oceanbase-oracle.sh. OceanBase needs about 8 GB and CI has
// no room for it, so TestCreateTableOptionsOceanBaseOracle reads that column
// back from a live server only where the DSN is set; TestCreateTableOptions
// holds the parser to the record unconditionally, which is the gate CI runs.
type createTableOptionCase struct {
	sql string
	// engine is whether OceanBase CE 4.3.5.6 parses the statement.
	engine bool
	// pending marks a statement OceanBase accepts that this parser does not
	// read yet. Implementing it means clearing the flag in the same change.
	pending bool
}

var createTableOptionCases = []createTableOptionCase{
	// CHECKSUM takes a number; TABLE_CHECKSUM is not a name here at all.
	{sql: "CREATE TABLE obcto (c int) checksum = 1", engine: true},
	{sql: "CREATE TABLE obcto (c int) checksum = 2", engine: true},
	// CHECKSUM takes an integer and not MySQL's wider `ulong_num`.
	{sql: "CREATE TABLE obcto (c int) checksum = 1.5"},
	{sql: "CREATE TABLE obcto (a int) STATS_AUTO_RECALC=1.5"},
	{sql: "CREATE TABLE obcto (c int) checksum = 'abc'"},
	{sql: "CREATE TABLE obcto (c int) checksum = DEFAULT"},
	{sql: "CREATE TABLE obcto (c int) table_checksum = 0"},
	{sql: "CREATE TABLE obcto (c int) table_checksum 1"},
	{sql: "CREATE TABLE obcto (c int) table_checksum = 1"},
	{sql: "CREATE TABLE obcto (c int) table_checksum"},
	{sql: "CREATE TABLE obcto (c int) table_checksum = 'abc'"},

	// STATS_AUTO_RECALC takes any integer, and DEFAULT, and no string.
	{sql: "CREATE TABLE obcto (a int) STATS_AUTO_RECALC=0", engine: true},
	{sql: "CREATE TABLE obcto (a int) STATS_AUTO_RECALC=1", engine: true},
	{sql: "CREATE TABLE obcto (a int) stats_auto_recalc 2", engine: true},
	{sql: "CREATE TABLE obcto (a int) stats_auto_recalc = 10", engine: true},
	{sql: "CREATE TABLE obcto (a int) STATS_AUTO_RECALC=DEFAULT", engine: true},
	{sql: "CREATE TABLE obcto (a int) stats_auto_recalc default", engine: true},
	{sql: "CREATE TABLE obcto (a int) STATS_AUTO_RECALC='abc'"},
	{sql: "CREATE TABLE obcto (a int) STATS_AUTO_RECALC='1'"},
	{sql: "CREATE TABLE obcto (a int) STATS_AUTO_RECALC=-1"},
	{sql: "CREATE TABLE obcto (a int) STATS_AUTO_RECALC="},

	// The statistics counters TiDB adds are not options here.
	{sql: "CREATE TABLE obcto (a int) STATS_BUCKETS=1"},
	{sql: "CREATE TABLE obcto (a int) STATS_TOPN=1"},
	{sql: "CREATE TABLE obcto (a int) STATS_SAMPLE_RATE=0.1"},
	{sql: "CREATE TABLE obcto (a int) STATS_COL_CHOICE='all'"},
	{sql: "CREATE TABLE obcto (a int, b int) STATS_COL_LIST='a,b'"},

	// SERIAL is a type; SERIAL DEFAULT VALUE is not a column option.
	{sql: "CREATE TABLE obcto (a serial)", engine: true},
	{sql: "CREATE TABLE obcto (a int serial default value)"},
	{sql: "CREATE TABLE obcto (a bigint serial default value)"},
	{sql: "CREATE TABLE obcto (a int serial default value null)"},
	{sql: "CREATE TABLE obcto (a int serial default value primary key)"},
	{sql: "CREATE TABLE obcto (a int serial)"},
	{sql: "CREATE TABLE obcto (a int serial default)"},

	// YEAR takes a width and nothing else, and SQL_TSI_YEAR is not a type.
	{sql: "CREATE TABLE obcto (a year)", engine: true},
	{sql: "CREATE TABLE obcto (a year(4))", engine: true},
	{sql: "CREATE TABLE obcto (a year zerofill)"},
	{sql: "CREATE TABLE obcto (a year(4) zerofill)"},
	{sql: "CREATE TABLE obcto (a year(4) binary)"},
	{sql: "CREATE TABLE obcto (a sql_tsi_year)"},
	{sql: "CREATE TABLE obcto (a SQL_TSI_YEAR(4))"},
	{sql: "CREATE TABLE obcto (a SQL_TSI_YEAR(4) ZEROFILL)"},

	// A width written with nothing in it is not a width, on any type that may
	// leave one out — the gap #623 left open when it tightened the types whose
	// width is required.
	{sql: "CREATE TABLE obcto (a int(1))", engine: true},
	{sql: "CREATE TABLE obcto (a double(10,2))", engine: true},
	{sql: "CREATE TABLE obcto (a year())"},
	{sql: "CREATE TABLE obcto (a int())"},
	{sql: "CREATE TABLE obcto (a char())"},
	{sql: "CREATE TABLE obcto (a datetime())"},
	{sql: "CREATE TABLE obcto (a bit())"},
	{sql: "CREATE TABLE obcto (a blob())"},
	{sql: "CREATE TABLE obcto (a binary())"},
	{sql: "CREATE TABLE obcto (a decimal())"},
	{sql: "CREATE TABLE obcto (a float())"},
	{sql: "CREATE TABLE obcto (a decimal(10,))"},
	{sql: "CREATE TABLE obcto (a decimal(,2))"},
	{sql: "CREATE TABLE obcto (a float(10,))"},
	{sql: "CREATE TABLE obcto (a int(,))"},

	// An ENUM or SET member may be a hex or bit literal, in either spelling.
	{sql: "CREATE TABLE obcto (c1 enum(0x61, 'b'), c2 set(0x61, 'b'))", engine: true},
	{sql: "CREATE TABLE obcto (c1 enum(0b01100001, 'b'), c2 set(0b01100001, 'b'))", engine: true},
	{sql: "CREATE TABLE obcto (c1 enum(x'61'))", engine: true},
	{sql: "CREATE TABLE obcto (c1 enum(b'01100001'))", engine: true},
	{sql: "CREATE TABLE obcto (c1 enum(0x6162))", engine: true},
	{sql: "CREATE TABLE obcto (c1 set(0b0110000101100010))", engine: true},
	{sql: "CREATE TABLE obcto (c1 enum(0x61) charset utf8mb4)", engine: true},
	{sql: "CREATE TABLE obcto (c1 enum(0x61, 'b') character set utf8mb4)", engine: true},
	// 0x61 is the member 'a', so the engine then refuses the duplicate with
	// 1291, which is not a syntax error.
	{sql: "CREATE TABLE obcto (c1 enum('a', 0x61))", engine: true},
	{sql: "CREATE TABLE obcto (c1 enum())"},
	{sql: "CREATE TABLE obcto (c1 set())"},
	{sql: "CREATE TABLE obcto (c1 enum(1))"},
	{sql: "CREATE TABLE obcto (c1 enum(1.5))"},
	{sql: "CREATE TABLE obcto (c1 enum('a',))"},
	{sql: "CREATE TABLE obcto (c1 enum(0x))"},
	{sql: "CREATE TABLE obcto (c1 enum(0b))"},
	{sql: "CREATE TABLE obcto (c1 enum(0xZZ))"},
	// Digits that do not fill whole bytes are padded on the left, and a quoted
	// hex literal, alone among the spellings, needs an even number of them. The
	// contents of one may themselves begin with `0b`, which is two bytes and
	// not a prefix, and nothing between the quotes is a member with no bytes in
	// it.
	{sql: "CREATE TABLE obcto (c1 enum(0x6))", engine: true},
	{sql: "CREATE TABLE obcto (c1 enum(0x616))", engine: true},
	{sql: "CREATE TABLE obcto (c1 enum(0b1))", engine: true},
	{sql: "CREATE TABLE obcto (c1 enum(0b101))", engine: true},
	{sql: "CREATE TABLE obcto (c1 enum(x'0b61'))", engine: true},
	{sql: "CREATE TABLE obcto (c1 enum(x'0b'))", engine: true},
	{sql: "CREATE TABLE obcto (c1 enum(x'0B61'))", engine: true},
	{sql: "CREATE TABLE obcto (c1 set(x'0b61'))", engine: true},
	{sql: "CREATE TABLE obcto (c1 enum(b'1'))", engine: true},
	{sql: "CREATE TABLE obcto (c1 enum(b'101'))", engine: true},
	{sql: "CREATE TABLE obcto (c1 enum(x''))", engine: true},
	{sql: "CREATE TABLE obcto (c1 enum(b''))", engine: true},
	{sql: "CREATE TABLE obcto (c1 enum(x'6'))"},
	{sql: "CREATE TABLE obcto (c1 enum(x'616'))"},
	{sql: "CREATE TABLE obcto (c1 enum(x'0x61'))"},
	{sql: "CREATE TABLE obcto (c1 enum(b'0b01'))"},

	// A name may begin with a digit as long as it is not digits alone, and the
	// words that are numbers are not names.
	{sql: "CREATE TABLE `123` (123a1 int)", engine: true},
	{sql: "CREATE TABLE obcto (123a1 int)", engine: true},
	{sql: "CREATE TABLE obcto (1a int)", engine: true},
	{sql: "CREATE TABLE obcto (123abc int)", engine: true},
	{sql: "CREATE TABLE obcto (12e int)", engine: true},
	{sql: "CREATE TABLE obcto (1_2 int)", engine: true},
	{sql: "CREATE TABLE obcto (123a1 int, 456b2 int)", engine: true},
	{sql: "CREATE TABLE 1t (a int)", engine: true},
	{sql: "CREATE TABLE 123 (123a1 int)"},
	{sql: "CREATE TABLE obcto (1e1 int)"},
	{sql: "CREATE TABLE obcto (0x61 int)"},
	{sql: "SELECT 1e1 FROM obcto", engine: true},
	{sql: "SELECT 1E+10 FROM obcto", engine: true},
	{sql: "SELECT 1.5e3 FROM obcto", engine: true},
	{sql: "SELECT 1e+ FROM obcto"},

	// The engine attributes, which OceanBase has in no position.
	{sql: `CREATE TABLE obcto (id INT) SECONDARY_ENGINE_ATTRIBUTE='{"key":"value"}'`},
	{sql: "CREATE TABLE obcto (id INT) SECONDARY_ENGINE_ATTRIBUTE="},
	{sql: "CREATE TABLE obcto (id INT) SECONDARY_ENGINE_ATTRIBUTE"},
	{sql: `CREATE TABLE obcto (id INT SECONDARY_ENGINE_ATTRIBUTE='{"key":"value"}')`},
	{sql: "CREATE TABLE obcto (id INT SECONDARY_ENGINE_ATTRIBUTE=)"},
	{sql: "CREATE TABLE obcto (id INT SECONDARY_ENGINE_ATTRIBUTE)"},
	{sql: `CREATE TABLE obcto (id INT) ENGINE_ATTRIBUTE='{"key":"value"}'`},
	{sql: `CREATE TABLE obcto (id INT ENGINE_ATTRIBUTE='{"key":"value"}')`},
	{sql: "CREATE TABLE obcto (id INT ENGINE_ATTRIBUTE)"},
	{sql: `ALTER TABLE obcto SECONDARY_ENGINE_ATTRIBUTE='{"a":1}'`},
	{sql: `ALTER TABLE obcto add column c int SECONDARY_ENGINE_ATTRIBUTE='{"a":1}'`},

	// GLOBAL TEMPORARY and ON COMMIT are TiDB's; the ordinary temporary table
	// is untouched.
	{sql: "CREATE GLOBAL TEMPORARY TABLE obcto (a int) ON COMMIT DELETE ROWS"},
	{sql: "CREATE TEMPORARY TABLE obcto (a int) ON COMMIT DELETE ROWS"},
	{sql: "CREATE TABLE obcto (a int) ON COMMIT DELETE ROWS"},
	{sql: "CREATE GLOBAL TABLE obcto (a int)"},
	{sql: "DROP GLOBAL TEMPORARY TABLE obcto"},
	{sql: "CREATE GLOBAL TEMPORARY VIEW obctov AS SELECT 1"},
	{sql: "CREATE GLOBAL TEMPORARY DATABASE obctod"},
	{sql: "DROP GLOBAL TEMPORARY DATABASE obctod"},
	{sql: "CREATE TEMPORARY TABLE obcto (a int)", engine: true},
	{sql: "DROP TEMPORARY TABLE obcto", engine: true},
}

// TestCreateTableOptions holds the parser to the recorded OceanBase verdicts. It
// needs no server, so it is the gate CI runs: the parser must read exactly what
// CE 4.3.5.6 parses, minus the pending rows.
func TestCreateTableOptions(t *testing.T) {
	for _, c := range createTableOptionCases {
		t.Run(c.sql, func(t *testing.T) {
			want := c.engine && !c.pending
			_, err := Parse(c.sql)
			if got := err == nil; got != want {
				t.Errorf("parser accepts=%v, want %v for %q (OceanBase accepts=%v, pending=%v): %v",
					got, want, c.sql, c.engine, c.pending, err)
			}
		})
	}
}
