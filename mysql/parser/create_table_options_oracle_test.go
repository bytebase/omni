package parser

import "testing"

// The CREATE TABLE table options and column-type oddities the TiDB grammar gap
// inventory groups under "CREATE TABLE clauses"
// (docs/plans/2026-11-03-tidb-grammar-gap-inventory.md), checked here because a
// fix in TiDB's fork parent usually belongs in the parent too — and because the
// ones that do not belong here have to be read from MySQL 8.0 rather than
// assumed. TABLE_CHECKSUM, SERIAL DEFAULT VALUE, SQL_TSI_YEAR, the numeric
// modifiers on YEAR and an ENUM member written as a hex or bit literal are all
// MySQL's; the STATS_ counters and GLOBAL TEMPORARY are TiDB's own.
//
// Nothing is asserted about which is which: each statement goes to MySQL 8.0
// and the parser has to answer the same, so the verdicts are a measurement and
// not a memory of one. The malformed variants matter most — an option with no
// value, an ENUM with no members, a width written `()` — because they are what
// a corpus of statements someone meant to write cannot show.

func TestCreateTableOptionsMatchMySQL(t *testing.T) {
	cases := []string{
		// CHECKSUM and TABLE_CHECKSUM, its other name, and the values they refuse.
		`create table t (c int) table_checksum = 0`,
		`create table t (c int) table_checksum 1`,
		`create table t (c int) TABLE_CHECKSUM = 2`,
		`create table t (c int) table_checksum = 1.5`,
		`create table t (c int) table_checksum`,
		`create table t (c int) table_checksum =`,
		`create table t (c int) table_checksum = 'abc'`,
		`create table t (c int) table_checksum = -1`,
		`create table t (c int) TABLE_CHECKSUM = DEFAULT`,
		`create table t (c int) checksum = 1`,
		`create table t (c int) checksum = 1.5`,
		`create table t (c int) checksum = 'abc'`,
		`create table t (c int) checksum = DEFAULT`,

		// STATS_AUTO_RECALC takes 0, 1 and DEFAULT and no other value.
		`CREATE TABLE t (a int) STATS_AUTO_RECALC=0`,
		`CREATE TABLE t (a int) STATS_AUTO_RECALC=1`,
		`CREATE TABLE t (a int) STATS_AUTO_RECALC=DEFAULT`,
		`CREATE TABLE t (a int) stats_auto_recalc default`,
		`CREATE TABLE t (a int) stats_auto_recalc 2`,
		`CREATE TABLE t (a int) stats_auto_recalc = 10`,
		`CREATE TABLE t (a int) STATS_AUTO_RECALC='abc'`,
		`CREATE TABLE t (a int) STATS_AUTO_RECALC='1'`,
		`CREATE TABLE t (a int) STATS_AUTO_RECALC=-1`,
		`CREATE TABLE t (a int) STATS_AUTO_RECALC=`,

		// The statistics options TiDB adds and this engine answers for itself.
		`CREATE TABLE t (a int) STATS_BUCKETS=1`,
		`CREATE TABLE t (a int) STATS_TOPN=1`,
		`CREATE TABLE t (a int) STATS_SAMPLE_RATE=0.1`,
		`CREATE TABLE t (a int) STATS_COL_CHOICE='all'`,
		`CREATE TABLE t (a int, b int) STATS_COL_LIST='a,b'`,

		// SERIAL DEFAULT VALUE, in every position among the other column options,
		// and with one of its three words missing.
		`create table t (a int serial default value)`,
		`create table t (a bigint serial default value)`,
		`create table t (a varchar(3) serial default value)`,
		`create table t (a int serial default value null)`,
		`create table t (a int serial default value not null)`,
		`create table t (a int serial default value primary key)`,
		`create table t (a int unique serial default value)`,
		`create table t (a int comment 'x' serial default value)`,
		`create table t (a int serial default value serial default value)`,
		`create table t (a int serial default value, b int)`,
		`create table t (a int serial)`,
		`create table t (a int serial default)`,
		`create table t (a int serial value)`,

		// SQL_TSI_YEAR names the YEAR type, with the same optional width and the
		// same numeric modifiers; `()` is not a width and no charset follows.
		`CREATE TABLE t (a SQL_TSI_YEAR)`,
		`CREATE TABLE t (a SQL_TSI_YEAR(4), b SQL_TSI_YEAR)`,
		`CREATE TABLE t (a SQL_TSI_YEAR(4) ZEROFILL)`,
		`CREATE TABLE t (a SQL_TSI_YEAR UNSIGNED)`,
		`CREATE TABLE t (a SQL_TSI_YEAR(4) unsigned zerofill)`,
		`CREATE TABLE t (a SQL_TSI_YEAR(4) not null)`,
		`CREATE TABLE t (a YEAR ZEROFILL)`,
		`CREATE TABLE t (a YEAR(4) ZEROFILL)`,
		`CREATE TABLE t (a SQL_TSI_YEAR())`,
		`create table t (a year(4) binary)`,
		// SQL_TSI_DAY, the one other unit name that could have been a type, is in
		// tidb/parser/create_table_options_test.go and not here: MariaDB's grammar
		// takes any word as a type name and answers 4161, unknown data type, rather
		// than 1064, so an engine comparison cannot read a grammar verdict from it.

		// A width written with nothing in it is not a width, on any type that may
		// leave one out — the gap #623 left open when it tightened the types whose
		// width is required.
		`CREATE TABLE t (a YEAR())`,
		`CREATE TABLE t (a int())`,
		`CREATE TABLE t (a char())`,
		`CREATE TABLE t (a datetime())`,
		`CREATE TABLE t (a bit())`,
		`CREATE TABLE t (a blob())`,
		`CREATE TABLE t (a binary())`,
		`CREATE TABLE t (a decimal())`,
		`CREATE TABLE t (a float())`,
		`CREATE TABLE t (a decimal(10,))`,
		`CREATE TABLE t (a decimal(,2))`,
		`CREATE TABLE t (a float(10,))`,
		`CREATE TABLE t (a int(,))`,
		`CREATE TABLE t (a int(1))`,
		`CREATE TABLE t (a double(10,2))`,

		// An ENUM or SET member may be a hex or bit literal, in either spelling,
		// but not a number and not nothing.
		`create table t (c1 enum(0x61, 'b'), c2 set(0x61, 'b'))`,
		`create table t (c1 enum(0b01100001, 'b'), c2 set(0b01100001, 'b'))`,
		`create table t (c1 enum(x'61'))`,
		`create table t (c1 enum(b'01100001'))`,
		`create table t (c1 enum(0x6162))`,
		`create table t (c1 set(0b0110000101100010))`,
		`create table t (c1 enum(0x61) charset utf8mb4)`,
		`create table t (c1 enum(0x61, 'b') character set utf8mb4)`,
		`create table t (c1 enum())`,
		`create table t (c1 set())`,
		`create table t (c1 enum(1))`,
		`create table t (c1 enum(1.5))`,
		`create table t (c1 enum('a',))`,
		`create table t (c1 enum(0x))`,
		`create table t (c1 enum(0b))`,
		`create table t (c1 enum(0xZZ))`,

		// A name may begin with a digit as long as it is not digits alone, and the
		// words that are numbers are not names.
		"create table `123` (123a1 int)",
		`create table t (123a1 int)`,
		`create table t (1a int)`,
		`create table t (123abc int)`,
		`create table t (12e int)`,
		`create table t (1_2 int)`,
		`create table t (123a1 int, 456b2 int)`,
		`create table 1t (a int)`,
		`create table 9t.t1 (a int)`,
		`create table 123 (123a1 int)`,
		`create table t (1e1 int)`,
		`create table t (0x61 int)`,
		`select 1e1 from t`,
		`select 1E+10 from t`,
		`select 1.5e3 from t`,
		`select 1e+ from t`,

		// ENGINE_ATTRIBUTE and SECONDARY_ENGINE_ATTRIBUTE take a string and
		// require one, as a table option and as a column option.
		`CREATE TABLE t (id INT) SECONDARY_ENGINE_ATTRIBUTE='{"key":"value"}'`,
		`CREATE TABLE t (id INT) SECONDARY_ENGINE_ATTRIBUTE=`,
		`CREATE TABLE t (id INT) SECONDARY_ENGINE_ATTRIBUTE`,
		`CREATE TABLE t (id INT) TABLESPACE ts1 SECONDARY_ENGINE_ATTRIBUTE='{"key":"value"}'`,
		`CREATE TABLE t (id INT SECONDARY_ENGINE_ATTRIBUTE='{"key":"value"}')`,
		`CREATE TABLE t (id INT SECONDARY_ENGINE_ATTRIBUTE=)`,
		`CREATE TABLE t (id INT SECONDARY_ENGINE_ATTRIBUTE)`,
		`CREATE TABLE t (id INT) ENGINE_ATTRIBUTE='{"key":"value"}'`,
		`CREATE TABLE t (id INT ENGINE_ATTRIBUTE='{"key":"value"}')`,
		`CREATE TABLE t (id INT ENGINE_ATTRIBUTE)`,
		`CREATE TABLE t (id INT ENGINE_ATTRIBUTE=)`,

		// GLOBAL TEMPORARY and its ON COMMIT clause are TiDB's, and this engine
		// answers for whether it has them.
		`create global temporary table t (a int) on commit delete rows`,
		`create temporary table t (a int) on commit delete rows`,
		`create table t (a int) on commit delete rows`,
		`create global table t (a int)`,
		`drop global temporary table t`,
		`create temporary table t (a int)`,
		`drop temporary table t`,
	}
	o := startParserOracle(t)
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			assertOracleMatch(t, o, sql)
		})
	}
}
