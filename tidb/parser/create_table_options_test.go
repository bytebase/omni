package parser

import "testing"

// The CREATE TABLE table options and column-type oddities of the TiDB grammar
// gap inventory (docs/plans/2026-11-03-tidb-grammar-gap-inventory.md, Stage 1,
// "CREATE TABLE clauses"), with the malformed variants of each one beside it.
//
// The corpus in testdata records whether a statement parses and nothing more,
// and it only holds statements pingcap's own test writes, so it cannot see an
// option that takes a value it should refuse. These cases can: each is checked
// against the parser with no container, and against a live TiDB v8.5.5 in
// TestCreateTableOptionsMatchTiDB, so the verdict written here is the engine's
// rather than a memory of it.
//
// The rules the rejections rest on, each read from v8.5.5:
//
//   - CHECKSUM and TABLE_CHECKSUM, its other name, take an integer: not a
//     string, not DEFAULT, not a negative, and — unlike MySQL, whose
//     `ulong_num` is wider — not a decimal.
//   - The STATS_ options each take the value of their own kind: a count is an
//     integer, STATS_SAMPLE_RATE a number, and STATS_COL_CHOICE and
//     STATS_COL_LIST strings. STATS_AUTO_RECALC, which TiDB shares with MySQL,
//     takes 0, 1 or DEFAULT alone.
//   - SERIAL DEFAULT VALUE is three words and all three are required.
//   - SQL_TSI_YEAR is another name for YEAR, with the same optional width and
//     the same numeric modifiers; `()` is not a width.
//   - An ENUM or SET member may be a hex or bit literal, but not a number.
//   - GLOBAL TEMPORARY requires an ON COMMIT clause, which nothing else may
//     carry, and which comes after the options and the partition clause.
//   - SECONDARY_ENGINE_ATTRIBUTE and ENGINE_ATTRIBUTE are MySQL 8.0 options
//     TiDB does not have, as a table option or as a column option.
var createTableOptionsAccepted = []string{
	// CHECKSUM and TABLE_CHECKSUM.
	`create table t (c int) table_checksum = 0`,
	`create table t (c int) table_checksum 1`,
	`create table t (c int) TABLE_CHECKSUM = 2`,
	`create table t (c int) checksum = 1`,
	`create table t (c int) checksum 0`,

	// The TiDB statistics options, each with and without the `=`.
	`CREATE TABLE t (a int) STATS_BUCKETS=1`,
	`CREATE TABLE t (a int) STATS_BUCKETS 1`,
	`CREATE TABLE t (a int) STATS_TOPN=1`,
	`CREATE TABLE t (a int) STATS_TOPN 0`,
	`CREATE TABLE t (a int) STATS_SAMPLE_RATE=0.1`,
	`CREATE TABLE t (a int) STATS_SAMPLE_RATE 1`,
	`CREATE TABLE t (a int) STATS_COL_CHOICE='all'`,
	`CREATE TABLE t (a int) STATS_COL_CHOICE 'list'`,
	// The server settles which strings are meaningful; the grammar takes any.
	`CREATE TABLE t (a int) STATS_COL_CHOICE='whatever'`,
	`CREATE TABLE t (a int, b int) STATS_COL_LIST='a,b'`,
	`CREATE TABLE t (a int, b int) STATS_COL_LIST 'a'`,
	// Several at once, with and without commas, and ahead of a partition clause.
	`CREATE TABLE t (a int) STATS_BUCKETS=1,STATS_TOPN=1`,
	`CREATE TABLE t (a int) STATS_BUCKETS=1 STATS_TOPN=1`,
	`CREATE TABLE t (a int) STATS_COL_CHOICE 'all' STATS_COL_LIST 'a'`,
	`CREATE TABLE t (a int) STATS_BUCKETS=1,STATS_TOPN=1 PARTITION BY RANGE (a) (PARTITION p1 VALUES LESS THAN (200))`,
	// STATS_AUTO_RECALC's whole value set.
	`CREATE TABLE t (a int) STATS_AUTO_RECALC=0`,
	`CREATE TABLE t (a int) STATS_AUTO_RECALC=1`,
	`CREATE TABLE t (a int) STATS_AUTO_RECALC=DEFAULT`,
	`CREATE TABLE t (a int) stats_auto_recalc default`,

	// SERIAL DEFAULT VALUE, on any type and in any position among the other
	// column options, and repeatable.
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

	// SQL_TSI_YEAR, and the modifiers YEAR itself now takes.
	`CREATE TABLE t (a SQL_TSI_YEAR)`,
	`CREATE TABLE t (a SQL_TSI_YEAR(4), b SQL_TSI_YEAR)`,
	`CREATE TABLE t (a SQL_TSI_YEAR(4) ZEROFILL)`,
	`CREATE TABLE t (a SQL_TSI_YEAR UNSIGNED)`,
	`CREATE TABLE t (a SQL_TSI_YEAR(4) unsigned zerofill)`,
	`CREATE TABLE t (a SQL_TSI_YEAR(4) not null)`,
	`CREATE TABLE t (a YEAR ZEROFILL)`,
	`CREATE TABLE t (a YEAR(4) ZEROFILL)`,

	// ENUM and SET members written as hex or bit literals, in both spellings.
	`create table t (c1 enum(0x61, 'b'), c2 set(0x61, 'b'))`,
	`create table t (c1 enum(0b01100001, 'b'), c2 set(0b01100001, 'b'))`,
	`create table t (c1 enum(x'61'))`,
	`create table t (c1 enum(b'01100001'))`,
	`create table t (c1 enum(0x6162))`,
	`create table t (c1 set(0b0110000101100010))`,
	`create table t (c1 enum(0x61) charset utf8mb4)`,
	`create table t (c1 enum(0x61, 'b') character set utf8mb4)`,
	// Digits that do not fill whole bytes are padded on the left, and the
	// contents of a quoted hex literal may themselves begin with `0b`, which is
	// two bytes and not a prefix.
	`create table t (c1 enum(0x6))`,
	`create table t (c1 enum(0x616))`,
	`create table t (c1 enum(0b1))`,
	`create table t (c1 enum(0b101))`,
	`create table t (c1 enum(x'0b61'))`,
	`create table t (c1 enum(x'0b'))`,
	`create table t (c1 enum(x'0B61'))`,
	`create table t (c1 set(x'0b61'))`,
	`create table t (c1 enum(b'1'))`,
	`create table t (c1 enum(b'101'))`,
	// A quoted literal with nothing between its quotes is a member with no
	// bytes in it.
	`create table t (c1 enum(x''))`,
	`create table t (c1 enum(b''))`,
	// The engine then refuses the duplicate member with 1291, which is not a
	// syntax error: 0x61 and 'a' are the same member.
	`create table t (c1 enum('a', 0x61))`,

	// A name may begin with a digit as long as it is not digits alone.
	"create table `123` (123a1 int)",
	`create table t (123a1 int)`,
	`create table t (1a int)`,
	`create table t (123abc int)`,
	`create table t (12e int)`,
	`create table t (1_2 int)`,
	`create table t (123a1 int, 456b2 int)`,
	`create table 1t (a int)`,
	`create table 9t.t1 (a int)`,

	// GLOBAL TEMPORARY and its ON COMMIT clause.
	`create global temporary table t (a int, b varchar(255)) on commit delete rows`,
	`create global temporary table t (a int) ON COMMIT PRESERVE ROWS`,
	`create global temporary table if not exists t (a int) on commit delete rows`,
	`create global temporary table t (a int, b varchar(255)) partition by hash(a) partitions 10 on commit delete rows`,
	`create global temporary table t (a int) comment 'x' on commit delete rows`,
	`create global temporary table t like t1 on commit delete rows`,
	`drop global temporary table t`,
	`drop global temporary table if exists t`,
	`drop global temporary table t1, t2`,
	// The ordinary temporary table is untouched by any of it.
	`create temporary table t (a int, b varchar(255))`,
	`drop temporary table t`,
}

var createTableOptionsRejected = []string{
	// A checksum with no value, a string, DEFAULT, a negative or a decimal.
	`create table t (c int) table_checksum`,
	`create table t (c int) table_checksum =`,
	`create table t (c int) table_checksum = 'abc'`,
	`create table t (c int) table_checksum = -1`,
	`create table t (c int) TABLE_CHECKSUM = DEFAULT`,
	`create table t (c int) table_checksum = 1.5`,
	`create table t (c int) checksum = 'abc'`,
	`create table t (c int) checksum = DEFAULT`,
	`create table t (c int) checksum = 1.5`,

	// Each statistics option with a value of the wrong kind, or none.
	`CREATE TABLE t (a int) STATS_BUCKETS`,
	`CREATE TABLE t (a int) STATS_BUCKETS=`,
	`CREATE TABLE t (a int) STATS_BUCKETS='abc'`,
	`CREATE TABLE t (a int) STATS_BUCKETS="1"`,
	`CREATE TABLE t (a int) STATS_BUCKETS=1.5`,
	`CREATE TABLE t (a int) STATS_BUCKETS=-1`,
	`CREATE TABLE t (a int) STATS_BUCKETS=DEFAULT`,
	`CREATE TABLE t (a int) STATS_TOPN`,
	`CREATE TABLE t (a int) STATS_TOPN='abc'`,
	`CREATE TABLE t (a int) STATS_SAMPLE_RATE='abc'`,
	`CREATE TABLE t (a int) STATS_SAMPLE_RATE=-0.1`,
	`CREATE TABLE t (a int) STATS_SAMPLE_RATE=DEFAULT`,
	`CREATE TABLE t (a int) STATS_COL_CHOICE`,
	`CREATE TABLE t (a int) STATS_COL_CHOICE=`,
	`CREATE TABLE t (a int) STATS_COL_CHOICE=1`,
	`CREATE TABLE t (a int) STATS_COL_CHOICE=all`,
	`CREATE TABLE t (a int, b int) STATS_COL_LIST=1`,
	// STATS_OPTIONS is an ALTER TABLE option here and not a CREATE TABLE one.
	`create table t (a int) STATS_OPTIONS='str'`,
	// STATS_AUTO_RECALC takes 0, 1 and DEFAULT and no other value.
	`CREATE TABLE t (a int) stats_auto_recalc 2`,
	`CREATE TABLE t (a int) stats_auto_recalc = 10`,
	`CREATE TABLE t (a int) STATS_AUTO_RECALC='abc'`,
	`CREATE TABLE t (a int) STATS_AUTO_RECALC='1'`,
	`CREATE TABLE t (a int) STATS_AUTO_RECALC=-1`,
	`CREATE TABLE t (a int) STATS_AUTO_RECALC=`,

	// SERIAL DEFAULT VALUE with a word missing.
	`create table t (a int serial)`,
	`create table t (a int serial default)`,
	`create table t (a int serial value)`,

	// A width written with nothing in it is not a width, on any type that may
	// leave one out. No corpus row carries one, which is how `varcharacter()`
	// shipped parsing in #623.
	`CREATE TABLE t (a SQL_TSI_YEAR())`,
	`CREATE TABLE t (a YEAR())`,
	`CREATE TABLE t (a int())`,
	`CREATE TABLE t (a char())`,
	`CREATE TABLE t (a datetime())`,
	`CREATE TABLE t (a bit())`,
	`CREATE TABLE t (a blob())`,
	`CREATE TABLE t (a binary())`,
	`CREATE TABLE t (a decimal())`,
	`CREATE TABLE t (a float())`,
	// A precision opens the same way, and a scale after a comma is a number too.
	`CREATE TABLE t (a decimal(10,))`,
	`CREATE TABLE t (a decimal(,2))`,
	`CREATE TABLE t (a float(10,))`,
	`CREATE TABLE t (a int(,))`,
	// SQL_TSI_YEAR is the one unit name that is also a type.
	`CREATE TABLE t (a SQL_TSI_DAY)`,
	// YEAR takes the numeric modifiers and no charset.
	`create table t (a year(4) binary)`,

	// An ENUM or SET member list that is empty, or holds something that is not
	// a string, a hex literal or a bit literal.
	`create table t (c1 enum())`,
	`create table t (c1 set())`,
	`create table t (c1 enum(1))`,
	`create table t (c1 enum(1.5))`,
	`create table t (c1 enum('a',))`,
	// `0x` and `0b` with no digit of their base are names, not literals.
	`create table t (c1 enum(0x))`,
	`create table t (c1 enum(0b))`,
	`create table t (c1 enum(0xZZ))`,
	// A quoted hex literal, and only that spelling, needs an even number of
	// digits, where the unquoted one pads.
	`create table t (c1 enum(x'6'))`,
	`create table t (c1 enum(x'616'))`,
	// Neither spelling takes the other's prefix inside the quotes.
	`create table t (c1 enum(x'0x61'))`,
	`create table t (c1 enum(b'0b01'))`,

	// A name of digits alone, and the literals that are not names.
	`create table 123 (123a1 int)`,
	`create table t (1e1 int)`,
	`create table t (0x61 int)`,

	// ON COMMIT belongs to GLOBAL TEMPORARY, is required there, and comes last.
	`create global temporary table t (a int, b varchar(255))`,
	`create temporary table t (a int, b varchar(255)) on commit delete rows`,
	`create table t (a int, b varchar(255)) on commit delete rows`,
	`create global temporary table t (a int) on commit`,
	`create global temporary table t (a int) on commit rows`,
	`create global temporary table t (a int) on commit delete`,
	`create global temporary table t (a int) on commit preserve`,
	`create global temporary table t (a int) on commit delete rows comment 'x'`,
	`create global temporary table t (a int) on commit preserve rows partition by hash(a)`,
	`create global temporary table t like t1`,
	`create global table t (a int)`,
	`drop global table t`,
	// GLOBAL TEMPORARY reaches TABLE and nothing else, on either statement.
	`create global temporary view v as select 1`,
	`create global temporary index i on t (a)`,
	`create global temporary database d`,
	`create global temporary sequence s`,
	`drop global temporary database d`,
	`drop global temporary view v`,
	`drop global temporary index i on t`,

	// The MySQL 8.0 engine attributes, which TiDB has in no position.
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
	`ALTER TABLE t SECONDARY_ENGINE_ATTRIBUTE='{"a":1}'`,
	`alter table t ENGINE_ATTRIBUTE='{"a":1}'`,
	`alter table t add column c int SECONDARY_ENGINE_ATTRIBUTE='{"a":1}'`,
}

func TestCreateTableOptionsAccepted(t *testing.T) {
	for _, sql := range createTableOptionsAccepted {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Fatalf("Parse(%q) error: %v", sql, err)
			}
		})
	}
}

func TestCreateTableOptionsRejected(t *testing.T) {
	for _, sql := range createTableOptionsRejected {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err == nil {
				t.Fatalf("Parse(%q) was accepted, but TiDB v8.5.5 rejects it", sql)
			}
		})
	}
}

// TestCreateTableOptionsMatchTiDB holds the two tables above to the engine
// itself, so the rejections above stay a measurement of v8.5.5 and not a
// reading of its grammar file.
func TestCreateTableOptionsMatchTiDB(t *testing.T) {
	tc := startTiDB(t)
	for want, cases := range map[bool][]string{
		true:  createTableOptionsAccepted,
		false: createTableOptionsRejected,
	} {
		for _, sql := range cases {
			t.Run(sql, func(t *testing.T) {
				accepts, err := tc.prepareSyntax(sql)
				if err != nil {
					t.Fatalf("the engine gave no verdict for %q: %v", sql, err)
				}
				if accepts != want {
					t.Errorf("TiDB %ss %q, this file records %s",
						verdictName(accepts), sql, verdictName(want))
				}
			})
		}
	}
}
