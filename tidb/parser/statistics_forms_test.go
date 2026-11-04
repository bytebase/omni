package parser

import "testing"

// The forms TiDB's standalone statistics statements take, and -- the point of
// this file -- the forms they do not.
//
// The corpus in testdata/pingcap_parser_corpus.tsv pins what pingcap's own
// parser test writes, and for this family that is 44 accepting rows and
// almost no rejections. A reader that took any number and any word after
// WITH would turn every one of those rows green while accepting `WITH 4
// PENGUINS`, so the rejections are written here by hand: an unknown unit, a
// missing number, a non-integer or out-of-range number, a unit with no
// number, two WITH clauses, a trailing comma, the arms that do not combine,
// and each word written backquoted.
//
// Every case below was measured against pingcap/tidb:v8.5.5 through PREPARE
// on 2026-11-04, and the recorded verdict is what that server answered.
// TestStatisticsFormsMatchTiDB is what keeps it a recording: an image upgrade
// that changes an answer fails there rather than leaving a stale expectation
// behind.
//
// Ref: AnalyzeTableStmt, DropStatsStmt, LoadStatsStmt, CreateStatisticsStmt
// and DropStatisticsStmt in pingcap/tidb pkg/parser/parser.y.

// statisticsFormCase is a statement and whether TiDB's parser accepts it.
type statisticsFormCase struct {
	sql    string
	accept bool
}

var statisticsFormCases = []statisticsFormCase{
	// --- ANALYZE TABLE: the modifiers before TABLE ---
	{"analyze table t", true},
	{"analyze local table t", true},
	{"analyze no_write_to_binlog table t", true},
	{"analyze local incremental table t index a", true},
	{"analyze no_write_to_binlog incremental table t index a", true},
	// INCREMENTAL comes after the binlog modifier, not before it.
	{"analyze incremental local table t index a", false},
	// TABLE is not optional. omni used to match it leniently, which left
	// `ANALYZE t` parsing as an ANALYZE of t on every one of these servers.
	{"analyze t", false},
	{"analyze table", false},
	// INCREMENTAL is a word, not a name: a backquoted one is a table.
	{"analyze `incremental` table t index a", false},

	// --- ANALYZE TABLE: the column choice ---
	{"analyze table t all columns", true},
	{"analyze table t predicate columns", true},
	{"analyze table t columns c1,c2", true},
	{"analyze table t columns `c1`, `c2`", true},
	// ALL COLUMNS and PREDICATE COLUMNS take a list of tables; a COLUMNS
	// list takes one.
	{"analyze table t1,t2 all columns", true},
	{"analyze table t1,t2 predicate columns", true},
	{"analyze table t1,t2 columns c1", false},
	// The list is never empty and its names carry no qualifier.
	{"analyze table t columns", false},
	{"analyze table t columns c1,", false},
	{"analyze table t columns c1.c2", false},
	{"analyze table t columns t.c1,t.c2", false},
	// ALL and COLUMNS are keywords here, so neither is a backquoted name.
	{"analyze table t `all` columns", false},
	{"analyze table t all `columns`", false},
	// One choice only, and each choice word is two words or none.
	{"analyze table t columns c1,c2 all columns", false},
	{"analyze table t all", false},
	{"analyze table t predicate", false},
	{"analyze table t partition p0 partition p1", false},

	// --- ANALYZE TABLE: the INDEX arm ---
	{"analyze table t index", true},
	{"analyze table t index a", true},
	{"analyze table t index a,b", true},
	{"analyze table t index a, b", true},
	{"analyze table t index `a`", true},
	// TiDB writes IndexNameList with an empty base case and a left
	// recursion on the comma, so a leading comma parses and a trailing or
	// doubled one does not.
	{"analyze table t index ,a", true},
	{"analyze table t index ,a,b", true},
	{"analyze table t index a,", false},
	{"analyze table t index ,,a", false},
	{"analyze table t index a,,b", false},
	{"analyze table t index 4", false},
	// The INDEX arm takes one table and combines with no column choice.
	{"analyze table t1,t2 index a", false},
	{"analyze table t index a all columns", false},
	{"analyze table t index a predicate columns", false},
	{"analyze table t index a columns c", false},
	{"analyze table t all columns index a", false},
	{"analyze table t columns c1 index a", false},

	// --- ANALYZE TABLE: the PARTITION arm ---
	{"analyze table t partition a", true},
	{"analyze table t partition p0,p1", true},
	{"analyze table t partition `p0`", true},
	{"analyze table t partition a all columns", true},
	{"analyze table t partition a predicate columns", true},
	{"analyze table t partition a columns c1,c2", true},
	{"analyze table t partition a index b", true},
	{"analyze table t partition p0 index", true},
	{"analyze table t partition p0 index a,b with 4 buckets", true},
	{"analyze table db.t partition p0 index b", true},
	// One table, bare names, and no parentheses -- which LOCK STATS takes
	// and this does not.
	{"analyze table t1,t2 partition p0", false},
	{"analyze table t partition", false},
	{"analyze table t partition (p0)", false},
	{"analyze table t partition p0 all columns index a", false},
	{"analyze table t partition p0 index a all columns", false},

	// --- ANALYZE TABLE: INCREMENTAL reaches only the INDEX arms ---
	{"analyze incremental table t index", true},
	{"analyze incremental table t index idx", true},
	{"analyze incremental table t partition p0 index b", true},
	{"analyze incremental table t", false},
	{"analyze incremental table t all columns", false},
	{"analyze incremental table t columns c1", false},
	{"analyze incremental table t1,t2 index a", false},
	{"analyze incremental table t update histogram on b", false},
	{"analyze incremental table t partition p0", false},
	{"analyze incremental table t with 4 buckets", false},
	{"analyze incremental table t drop histogram on b", false},

	// --- ANALYZE TABLE: the WITH list ---
	{"analyze table t with 4 buckets", true},
	{"analyze table t with 4 topn", true},
	{"analyze table t with 4 cmsketch width", true},
	{"analyze table t with 4 cmsketch depth", true},
	{"analyze table t with 4 samples", true},
	{"analyze table t with 10 samplerate", true},
	{"analyze table t with 0.1 samplerate", true},
	{"analyze table t with 4 BUCKETS", true},
	{"analyze table t with 4 TopN", true},
	{"analyze table t with 4 CMSketch Width", true},
	{"analyze table t with 4 buckets, 4 topn, 4 cmsketch width, 4 cmsketch depth, 4 samples", true},
	{"analyze table t with 4 buckets,4 topn", true},
	{"analyze table t with 0.1 samplerate, 4 buckets", true},
	{"analyze table t with 4 buckets, 0.5 samplerate", true},
	// The list is not a set: a repeat parses and the parser keeps it.
	{"analyze table t with 4 buckets, 4 buckets", true},
	{"analyze table t with 4 buckets, 8 buckets", true},
	{"analyze table t with 4 samplerate, 4 samplerate", true},
	// Zero is a value; the range is uint64 and nothing wider.
	{"analyze table t with 0 buckets", true},
	{"analyze table t with 0 topn", true},
	{"analyze table t with 0 samples", true},
	{"analyze table t with 18446744073709551615 buckets", true},
	{"analyze table t with 18446744073709551616 buckets", false},
	{"analyze table t with 9999999999999999999 buckets", true},
	// The unit list is closed.
	{"analyze table t with 4 penguins", false},
	{"analyze table t with 4 cmsketch", false},
	{"analyze table t with 4 width", false},
	{"analyze table t with 4 depth", false},
	// A number with no unit, a unit with no number, and neither.
	{"analyze table t with 4", false},
	{"analyze table t with buckets", false},
	{"analyze table t with", false},
	// Only SAMPLERATE takes a number that is not an unsigned integer.
	{"analyze table t with 1.5 buckets", false},
	{"analyze table t with 0.5 topn", false},
	{"analyze table t with -1 buckets", false},
	{"analyze table t with -0.1 samplerate", false},
	{"analyze table t with 1.0 samplerate", true},
	{"analyze table t with .5 samplerate", true},
	{"analyze table t with 1e3 samplerate", true},
	{"analyze table t with 1024 samplerate", true},
	// A hex, bit or string literal is no number here, for either kind of
	// unit.
	{"analyze table t with 0x10 buckets", false},
	{"analyze table t with 0x10 samplerate", false},
	{"analyze table t with 4e0 buckets", false},
	{"analyze table t with '4' buckets", false},
	// One WITH carries the whole list, and the separator is the comma.
	{"analyze table t with 4 buckets 4 topn", false},
	{"analyze table t index a with 4 buckets with 4 topn", false},
	{"analyze table t with 4 buckets,", false},
	// The unit words are words, not names.
	{"analyze table t with 4 `buckets`", false},
	{"analyze table t index a with 4 `topn`", false},
	// Every arm takes the list.
	{"analyze table t index a with 4 buckets", true},
	{"analyze table t index a with 4 buckets, 4 topn", true},
	{"analyze table t index a with 0.1 samplerate", true},
	{"analyze table t partition a with 4 buckets", true},
	{"analyze table t partition a index b with 4 buckets", true},
	{"analyze table t partition p0 index a with 4 buckets", true},
	{"analyze table t partition p0 columns c1,c2", true},
	{"analyze table t all columns with 1024 buckets", true},
	{"analyze no_write_to_binlog incremental table t partition p0 index a", true},
	{"analyze table t partition a all columns with 1024 buckets", true},
	{"analyze table t partition a columns c1,c2 with 1024 buckets", true},
	{"analyze table t partition a predicate columns with 1024 buckets", true},
	{"analyze table t all columns with 4 buckets", true},
	{"analyze table t predicate columns with 4 buckets", true},
	{"analyze table t columns c1,c2 with 4 topn", true},
	{"analyze table t columns c1 with 4 buckets, 4 topn", true},
	{"analyze table t1,t2 all columns with 4 topn", true},
	{"analyze table t1,t2 predicate columns with 4 topn", true},
	{"analyze table t1, t2 all columns with 4 buckets, 0.5 samplerate", true},
	{"analyze local table t with 4 buckets", true},

	// --- ANALYZE TABLE: the histogram arms ---
	{"analyze table t update histogram on b", true},
	{"analyze table t update histogram on c1, c2", true},
	{"analyze table t drop histogram on b", true},
	{"analyze table t drop histogram on c1, c2", true},
	{"analyze no_write_to_binlog table t update histogram on b", true},
	// UPDATE HISTOGRAM takes the same WITH list every other arm does;
	// DROP HISTOGRAM takes none.
	{"analyze table t update histogram on b with 1024 buckets", true},
	{"analyze table t update histogram on b with 4 topn", true},
	{"analyze table t update histogram on b with 4 buckets, 4 topn", true},
	{"analyze table t update histogram on b with 0.5 samplerate", true},
	{"analyze table t drop histogram on b with 4 buckets", false},
	// One table, bare column names, no partition list, and no INCREMENTAL.
	{"analyze table t1,t2 update histogram on b", false},
	{"analyze table t1,t2 drop histogram on b", false},
	{"analyze table t update histogram on t.c1", false},
	{"analyze table t update histogram on t.c1, t.c2", false},
	{"analyze table t drop histogram on t.c1, t.c2", false},
	{"analyze table t partition p0 update histogram on b", false},
	{"analyze table t index a update histogram on b", false},
	{"analyze table t index a drop histogram on b", false},
	{"analyze table t all columns update histogram on b", false},
	// A tail the grammar has ended is a rejection, not the start of the
	// next statement.
	{"analyze table t drop table u", false},
	{"analyze table t partition p0 update histogram on b with 4 buckets", false},
	{"analyze table t update histogram", false},
	{"analyze table t drop histogram", false},
	// MySQL 8.0's USING DATA tail, which TiDB has not. omni used to take
	// it here.
	{"analyze table t update histogram on b using data '{}'", false},
	{"analyze table t update histogram on b, c using data '{}'", false},

	// --- ANALYZE TABLE: the table list ---
	{"analyze table t1,t2", true},
	{"analyze table t1, t2, t3", true},
	{"analyze table db1.t1", true},
	{"analyze table t1.*", false},

	// --- DROP STATS ---
	{"drop stats t", true},
	{"drop stats db.t", true},
	{"drop stats `t`", true},
	{"drop stats t1, t2, t3", true},
	{"drop stats t global", true},
	{"drop stats t partition p0", true},
	{"drop stats t partition p0, p1, p2", true},
	// GLOBAL and the partition list each take one table, and they do not
	// combine.
	{"drop stats t1,t2 global", false},
	{"drop stats t1,t2 partition p0", false},
	{"drop stats t global partition p0", false},
	{"drop stats t global, t2", false},
	{"drop stats", false},
	{"drop stats t partition", false},
	{"drop stats t partition p0,", false},
	// No parentheses, unlike LOCK STATS.
	{"drop stats t partition (p0)", false},
	{"drop stats t partition (p0,p1)", false},
	// STATS is a word: a backquoted one is the table DROP TABLE would
	// drop, and `DROP `STATS` t` is no statement at all.
	{"drop `stats` t", false},

	// --- LOAD STATS ---
	{"load stats '/tmp/stats.json'", true},
	{"load stats \"x\"", true},
	{"load stats", false},
	{"load stats 4", false},
	{"load stats 'a' 'b'", false},
	{"load `stats` 'x'", false},

	// --- CREATE STATISTICS ---
	{"create statistics stats1 (cardinality) on t(a,b,c)", true},
	{"create statistics stats1(cardinality) on t(a,b,c)", true},
	{"create statistics stats2 (dependency) on t(a,b)", true},
	{"create statistics stats3 (correlation) on t(a,b)", true},
	{"create statistics s1 (CARDINALITY) on t(a)", true},
	{"create statistics `s1` (cardinality) on t(a)", true},
	{"create statistics if not exists stats1 (cardinality) on t(a,b,c)", true},
	{"create statistics s1 (cardinality) on db.t(a)", true},
	// Unlike ANALYZE's COLUMNS list these names may be qualified, but
	// `t.*` is not a name.
	{"create statistics s1 (cardinality) on t(t.a)", true},
	{"create statistics s1 (cardinality) on t(d.t.a)", true},
	{"create statistics s1 (cardinality) on t(a.*)", false},
	// Exactly one type, from a closed list.
	{"create statistics stats3 on t(a,b)", false},
	{"create statistics if not exists stats3 on t(a,b)", false},
	{"create statistics s1 (frequency) on t(a,b)", false},
	{"create statistics s1 (cardinality,dependency) on t(a)", false},
	{"create statistics s1 (`cardinality`) on t(a)", false},
	// The column list is required and not empty.
	{"create statistics s1 (cardinality) on t", false},
	{"create statistics s1 (cardinality) on t()", false},
	{"create statistics s1 (cardinality) on t(a,)", false},
	{"create `statistics` s1 (cardinality) on t(a)", false},

	// --- DROP STATISTICS ---
	{"drop statistics stats1", true},
	{"drop statistics `s1`", true},
	{"drop statistics", false},
	{"drop statistics s1, s2", false},
	{"drop statistics if exists s1", false},
}

// TestStatisticsForms holds this parser to statisticsFormCases. It needs no
// container, so it guards the grammar wherever it runs.
func TestStatisticsForms(t *testing.T) {
	for _, c := range statisticsFormCases {
		t.Run(c.sql, func(t *testing.T) {
			_, err := Parse(c.sql)
			if accepts := err == nil; accepts != c.accept {
				t.Errorf("omni %ss %q, TiDB %s %ss it (err=%v)",
					verdictName(accepts), c.sql, tidbCorpusVersion, verdictName(c.accept), err)
			}
		})
	}
}

// statisticsWordsAsNames are statements in which each word this family reads
// bare is an ordinary name instead: a table, a column, an alias, an index.
// None of these words is registered in the lexer, and these statements are
// what says so -- registering one would break every site that reads a name,
// which is why the family reads them with isBareWord.
//
// Every statement here parses on v8.5.5 too, which
// TestStatisticsWordsStayNamesOnTiDB measures.
var statisticsWordsAsNames = []string{
	// As a table name, bare and qualified.
	"SELECT * FROM stats",
	"SELECT * FROM statistics",
	"SELECT * FROM topn",
	"SELECT * FROM cmsketch",
	"SELECT * FROM samples",
	"SELECT * FROM samplerate",
	"SELECT * FROM cardinality",
	"SELECT * FROM dependency",
	"SELECT * FROM correlation",
	"SELECT * FROM incremental",
	"SELECT * FROM predicate",
	"SELECT * FROM width",
	"SELECT * FROM depth",
	"SELECT * FROM buckets",
	// As an implicit table alias, the site that broke when the SHOW
	// keywords were registered.
	"DELETE FROM t stats WHERE stats.id = 1",
	"DELETE FROM t topn WHERE topn.id = 1",
	"DELETE FROM t cmsketch WHERE cmsketch.id = 1",
	"DELETE FROM t samplerate WHERE samplerate.id = 1",
	"DELETE FROM t incremental WHERE incremental.id = 1",
	"DELETE FROM t predicate WHERE predicate.id = 1",
	// As a select-list alias and a column name.
	"SELECT 1 AS statistics",
	"SELECT 1 AS samples",
	"SELECT 1 AS correlation",
	"SELECT depth, width, topn FROM t",
	"SELECT t.cardinality FROM t",
	// As names in DDL, which is where a registered keyword hurts most.
	"CREATE TABLE stats (samples int, width int, depth int, topn int)",
	"CREATE TABLE statistics (cardinality int, dependency int, correlation int)",
	"CREATE TABLE t (samplerate double, incremental int, predicate int, INDEX cmsketch (incremental))",
	"CREATE INDEX topn ON t (a)",
	"DROP TABLE stats",
	"DROP TABLE statistics",
	"ALTER TABLE t ADD COLUMN samplerate double",
	// And beside the statements of this family, which is what makes the
	// bare-word reading load-bearing rather than incidental.
	"ANALYZE TABLE stats",
	"ANALYZE TABLE statistics COLUMNS samples, width",
	"ANALYZE TABLE t INDEX topn WITH 4 BUCKETS",
	"DROP STATS stats",
	"DROP STATISTICS statistics",
	"CREATE STATISTICS statistics (CARDINALITY) ON stats(samples)",
}

// TestStatisticsWordsStayNames holds the other half of the keyword decision:
// every word this family reads bare is still an ordinary name everywhere
// else. It needs no container.
func TestStatisticsWordsStayNames(t *testing.T) {
	for _, sql := range statisticsWordsAsNames {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err != nil {
				t.Errorf("Parse(%q): %v", sql, err)
			}
		})
	}
}
