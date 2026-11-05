package parser

import "testing"

// The forms TiDB's LOAD DATA and IMPORT INTO take, and -- the point of this
// file -- the forms they do not.
//
// The corpus in testdata/pingcap_parser_corpus.tsv pins what pingcap's own
// parser test writes for this family: 25 accepting rows and no rejection at
// all. A reader that took any word after WITH, any value after DEFINED NULL
// BY, and FORMAT anywhere would turn all 25 of them green while accepting
// `WITH 4 PENGUINS`, so the rejections are written here by hand: an empty
// FIELDS list, a LINES clause out of order, IGNORE counted in ROWS, FORMAT
// in the wrong place or spelled with a byte literal, a WITH list that is
// empty, has a trailing comma, repeats the clause or carries an expression,
// a bare OPTIONALLY with no ENCLOSED BY, IMPORT INTO with no FROM, and the
// pieces IMPORT INTO's query source refuses.
//
// Every case below was measured against pingcap/tidb:v8.5.5 through PREPARE
// on 2026-11-04, and the recorded verdict is what that server answered.
// TestLoadDataFormsMatchTiDB is what keeps it a recording: an image upgrade
// that changes an answer fails there rather than leaving a stale expectation
// behind.
//
// Ref: LoadDataStmt, ImportIntoStmt, Fields, FieldItem, FieldTerminator,
// Lines and LoadDataOption in pingcap/tidb pkg/parser/parser.y.

// loadDataFormCase is a statement and whether TiDB's parser accepts it.
type loadDataFormCase struct {
	sql    string
	accept bool
}

var loadDataFormCases = []loadDataFormCase{
	// --- FIELDS and LINES option values: a string, a hex or a bit literal ---
	{"load data infile '/x.csv' into table t fields terminated by X'6B6B'", true},
	{"load data infile '/x.csv' into table t fields terminated by x'6b6b'", true},
	{"load data infile '/x.csv' into table t fields terminated by 0x6B6B", true},
	{"load data infile '/x.csv' into table t fields terminated by B'110101101101011'", true},
	{"load data infile '/x.csv' into table t fields terminated by b'1101'", true},
	// A byte literal with nothing between its quotes is the empty string.
	{"load data infile '/x.csv' into table t fields terminated by X''", true},
	{"load data infile '/x.csv' into table t fields enclosed by X'0D'", true},
	// A multi-byte ENCLOSED BY parses and fails later (1083), so it is no
	// rejection of this grammar.
	{"load data infile '/x.csv' into table t fields enclosed by X'0D0D'", false},
	{"load data infile '/x.csv' into table t fields escaped by X'5C'", true},
	{"load data infile '/x.csv' into table t fields optionally enclosed by X'22'", true},
	{"load data infile '/x.csv' into table t columns terminated by X'09'", true},
	{"load data infile '/x.csv' into table t lines starting by X'6B6B' terminated by X'6B6B'", true},
	{"load data infile '/x.csv' into table t lines starting by B'1101' terminated by B'1101'", true},
	// A hex literal with an odd number of digits is no literal at all.
	{"load data infile '/x.csv' into table t fields terminated by X'616'", false},

	// --- FIELDS: a list of one or more items, in any order ---
	{"load data infile '/x.csv' into table t fields terminated by 'a' terminated by 'b'", true},
	{"load data infile '/x.csv' into table t fields escaped by 'a' terminated by 'b'", true},
	{"load data infile '/x.csv' into table t fields enclosed by 'a' escaped by 'b' terminated by 'c'", true},
	// FIELDS with nothing after it names no item.
	{"load data infile '/x.csv' into table t fields", false},
	{"load data infile '/x.csv' into table t fields lines terminated by 'a'", false},
	// ENCLOSED BY is not optional after OPTIONALLY.
	{"load data infile '/x.csv' into table t fields optionally enclosed ignore 1 lines", false},
	{"load data infile '/x.csv' into table t fields optionally enclosed by '\"'", true},
	// A backquoted word is a name, and no name is an item here.
	{"load data infile '/x.csv' into table t fields optionally `enclosed` by 'a'", false},

	// --- FIELDS DEFINED NULL BY, which TiDB has and MySQL has not ---
	{"load data infile '/x.csv' into table t fields defined null by 'nil'", true},
	{"load data infile '/x.csv' into table t fields defined null by X'00'", true},
	{"load data infile '/x.csv' into table t fields defined null by B'1101'", true},
	{"load data infile '/x.csv' into table t fields defined null by 'a' defined null by 'b'", true},
	{"load data infile '/x.csv' into table t fields defined null by 'a' optionally enclosed defined null by 'b'", true},
	{"load data infile '/x.csv' into table t fields optionally enclosed by 'a' enclosed by 'b'", true},
	{"load data infile '/x.csv' into table t fields enclosed by 'a' optionally enclosed by 'b'", true},
	{"load data infile '/x.csv' into table t fields terminated by 'a' defined null by 'b' enclosed by 'c'", true},
	{"load data infile '/x.csv' into table t fields defined null by 'a' terminated by 'b'", true},
	// Its OPTIONALLY ENCLOSED tail carries no BY and no value.
	{"load data infile '/x.csv' into table t fields defined null by 'NULL' optionally enclosed", true},
	{"load data infile '/x.csv' into table t fields defined null by 'NULL' optionally enclosed ignore 1 lines", true},
	// A BY after it makes it the next item instead, an ENCLOSED one.
	{"load data infile '/x.csv' into table t fields defined null by 'a' optionally enclosed by 'b'", true},
	{"load data infile '/x.csv' into table t fields defined null by 'a' optionally enclosed by X'62'", true},
	// OPTIONALLY with no ENCLOSED after it at all.
	{"load data infile '/x.csv' into table t fields defined null by 'a' optionally", false},
	// Every word of DEFINED NULL BY is a word, not a name.
	{"load data infile '/x.csv' into table t fields `defined` null by 'a'", false},
	{"load data infile '/x.csv' into table t fields defined `null` by 'a'", false},
	{"load data infile '/x.csv' into table t fields defined null `by` 'a'", false},

	// --- LINES: STARTING before TERMINATED, each once ---
	{"load data infile '/x.csv' into table t lines", true},
	{"load data infile '/x.csv' into table t lines starting by 'a'", true},
	{"load data infile '/x.csv' into table t lines terminated by 'b'", true},
	{"load data infile '/x.csv' into table t lines starting by 'a' terminated by 'b'", true},
	{"load data infile '/x.csv' into table t lines terminated by 'b' starting by 'c'", false},
	{"load data infile '/x.csv' into table t lines starting by 'a' starting by 'b'", false},
	{"load data infile '/x.csv' into table t lines starting by 'a' terminated by 'b' starting by 'c'", false},

	// --- IGNORE n LINES, counted in LINES and not in ROWS ---
	{"load data infile '/x.csv' into table t ignore 2 lines", true},
	{"load data infile '/x.csv' into table t ignore 2 rows", false},
	// IGNORE comes after FIELDS and LINES, not before them.
	{"load data infile '/x.csv' into table t ignore 2 lines fields terminated by ','", false},

	// --- FORMAT, between the file and REPLACE/IGNORE ---
	{"load data infile '/x.csv' format 'delimited data' into table t", true},
	{"load data infile '/x.sql' format 'sql file' replace into table t", true},
	{"load data infile '/x.csv' format 'csv' ignore into table t", true},
	{"load data local infile '/x.parquet' format 'parquet' into table t", true},
	{"load data infile '/x.csv' into table t format 'csv'", false},
	{"load data infile '/x.csv' format 'csv' format 'csv' into table t", false},
	{"load data infile '/x.csv' format X'6161' into table t", false},
	{"load data infile '/x.csv' `format` 'csv' into table t", false},

	// --- WITH, the option list ---
	{"load data infile '/x.csv' into table t with detached", true},
	{"load data infile '/x.csv' into table t with threads=10", true},
	{"load data infile '/x.csv' into table t with threads=10, detached", true},
	{"load data infile '/x.csv' into table t with threads=-1", true},
	{"load data infile '/x.csv' into table t with threads=+2", true},
	{"load data infile '/x.csv' into table t with threads=1.5", true},
	{"load data infile '/x.csv' into table t with threads='abc'", true},
	{"load data infile '/x.csv' into table t with threads=X'3130'", true},
	{"load data infile '/x.csv' into table t with threads=true", true},
	{"load data infile '/x.csv' into table t with threads=null", true},
	{"load data infile '/x.csv' into table t with character_set='utf8'", true},
	{"load data infile '/x.csv' into table t with Threads=10", true},
	{"load data infile '/x.csv' into table t with detached, detached", true},
	// A backquoted option name is still an identifier, so it is a name here.
	{"load data infile '/x.csv' into table t with `detached`", true},
	{"load data infile '/x.csv' into table t with `threads`=10", true},
	{"load data infile '/x.csv' into table t fields terminated by 'a' with detached", true},
	{"load data infile '/x.csv' into table t (a) set b=1 with detached", true},
	// The list holds at least one option, with no trailing comma, no
	// reserved word for a name, no expression for a value, and one WITH to
	// a statement.
	{"load data infile '/x.csv' into table t with", false},
	{"load data infile '/x.csv' into table t with detached,", false},
	{"load data infile '/x.csv' into table t with threads=", false},
	{"load data infile '/x.csv' into table t with threads=@a", false},
	{"load data infile '/x.csv' into table t with select=1", false},
	{"load data infile '/x.csv' into table t with detached with detached", false},

	// --- the column-and-user-variable list and SET ---
	{"load data infile '/x.csv' into table t (a, @v) set b = @v", true},
	// A backquoted name is a column, even one spelled with an '@'.
	{"load data infile '/x.csv' into table t (`@v`)", true},
	{"load data infile '/x.csv' into table t (`@v`, a)", true},
	{"load data infile '/x.csv' into table t (a,@1) with detached", true},
	{"load data infile '/x.csv' into table t set b = default", true},

	// --- IMPORT INTO: the file source ---
	{"import into t from '/file.csv'", true},
	{"import into `t` from '/file.csv'", true},
	{"import into d.t from '/file.csv'", true},
	{"import into t (a,b) from '/file.csv'", true},
	{"import into t (a,@1) from '/file.csv'", true},
	{"import into t (@1) from '/file.csv'", true},
	// An empty pair of parentheses is a written list with no entries.
	{"import into t () from '/file.csv'", true},
	{"import into t (a,@1) set b=@1+100 from '/file.csv'", true},
	{"import into t set a=1 from '/file.csv'", true},
	{"import into t (a) set b=default from '/file.csv'", true},
	{"import into t from '/file.csv' format 'sql file'", true},
	{"import into t from '/file.csv' with detached", true},
	{"import into t from '/file.csv' with thread=1", true},
	{"import into t from '/file.csv' with detached, thread=1", true},
	{"import into t from '/file.csv' with detached=1", true},
	{"import into t from '/file.csv' with thread", true},
	{"import into t from '/file.csv' with __starting_ts=1", true},
	{"import into t from '/file.csv' format 'csv' with detached", true},
	{"import into t from '/file.csv' format 'csv' with cloud_storage_uri='s3://x'", true},
	{"import into t (a,b) set c=1 from '/file.csv' format 'csv' with thread=2, detached", true},
	// FROM is required, FORMAT follows the path and precedes WITH, the
	// format name is a string, and PARTITION is no part of this statement.
	{"import into t", false},
	{"import into t from", false},
	{"import into t from '/file.csv' with detached format 'csv'", false},
	{"import into t from 'f' format 'csv' format 'csv'", false},
	{"import into t from '/file.csv' format X'6363'", false},
	{"import into t from '/file.csv' `format` 'csv'", false},
	{"import into t from '/file.csv' with thread=1 ,", false},
	{"import into t partition (p0) from '/file.csv'", false},

	// --- IMPORT INTO: the query source ---
	{"import into t from select * from s", true},
	{"import into t from (select * from s)", true},
	{"import into t from with cte as (select 1) select * from cte", true},
	{"import into t from table s", true},
	{"import into t (a,b) from select * from s", true},
	{"import into t (a) from select a from s with thread=2", true},
	{"import into t (`@v`) from select 1", true},
	{"import into t from select a from s order by a with thread=2", true},
	{"import into t from select a from s limit 1 with thread=2", true},
	// After GROUP BY, a WITH is the query's own WITH ROLLUP, so an option
	// list needs its own WITH after it.
	{"import into t from select a from s group by a with rollup", true},
	{"import into t from select a from s group by a with rollup with thread=2", true},
	{"import into t from select a from s group by a with thread=2", false},
	// A query source carries neither a user variable nor a SET list, and a
	// VALUES statement is no query here.
	{"import into t (a,@1) from select a from s", false},
	{"import into t (a,@1) from table s", false},
	{"import into t set a=1 from select * from s", false},
	{"import into t (a) set b=1 from select * from s", false},
	{"import into t from values (1)", false},
	// The TABLE and VALUES primaries, their set operations and their tails.
	// The query-source tests in load_data_tree_test.go assert the same
	// verdicts, but as constants; these cases are re-measured against the
	// pinned container by TestLoadDataFormsMatchTiDB on every run, which is
	// what would catch one of those verdicts going stale.
	{"import into t from values row(1)", true},
	{"import into t from values row(1,1)", true},
	{"import into t from values row(1), row(2)", true},
	{"import into t from values row(1) union select 1", true},
	{"import into t from table t1 union table t2", true},
	{"import into t from table t1 union all table t2", true},
	{"import into t from table t1 intersect table t2", true},
	{"import into t from table t1 except table t2", true},
	{"import into t from table t1 order by a limit 1", true},
	// An incomplete source is still a 1064.
	{"import into t from table t1 union", false},
	{"import into t from table t1 union table", false},
	{"import into t from table", false},
	{"import into t from values", false},

	// --- IMPORT TABLE FROM is MySQL's, not TiDB's ---
	{"import table from '/tmp/a.sdi'", false},
}

// TestLoadDataForms holds this parser to the verdicts above.
func TestLoadDataForms(t *testing.T) {
	for _, c := range loadDataFormCases {
		t.Run(c.sql, func(t *testing.T) {
			_, err := Parse(c.sql)
			if accepts := err == nil; accepts != c.accept {
				t.Errorf("omni %ss %q, TiDB v8.5.5 %ss it: %v",
					verdictName(accepts), c.sql, verdictName(c.accept), err)
			}
		})
	}
}
