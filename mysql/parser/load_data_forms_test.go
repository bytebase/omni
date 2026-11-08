package parser

import "testing"

// The forms a LOAD DATA FIELDS or LINES clause takes, and -- the point of
// this file -- the forms it does not.
//
// Each option's value is a text_string, which is a string literal, a hex
// literal or a bit literal, so `TERMINATED BY X'6B6B'` is the two bytes
// "kk". This parser used to take only the string literal. The clause's items
// also come in any order, may repeat, and are never none: `FIELDS` and
// `LINES` with nothing after them are each a 1064, and so is a bare
// OPTIONALLY with no ENCLOSED BY, which this parser used to read.
//
// The clauses TiDB has grown on this statement -- FORMAT, FIELDS DEFINED
// NULL BY, the WITH option list -- and its IMPORT INTO are not here: each is
// a 1064 on mysql:8.0.46, mariadb:11.8.8 and OceanBase 4.3.5.
//
// Every verdict below was measured on 2026-11-04 against mysql:8.0.46,
// mariadb:11.8.8 and OceanBase 4.3.5 (oceanbase/oceanbase-ce:4.3.5-lts), all
// three of which answer the same for every case. TestLoadDataFieldFormsOracle holds them to mysql:8.0.
//
// Ref: field_term, line_term and text_string in sql_yacc.yy.

// loadDataFieldCase is a statement and whether the engine accepts it.
type loadDataFieldCase struct {
	sql    string
	accept bool
}

var loadDataFieldCases = []loadDataFieldCase{
	// --- the option value: a string, a hex or a bit literal ---
	{"load data infile '/x.csv' into table t fields terminated by X'6B6B'", true},
	{"load data infile '/x.csv' into table t fields terminated by x'6b6b'", true},
	{"load data infile '/x.csv' into table t fields terminated by 0x6B6B", true},
	{"load data infile '/x.csv' into table t fields terminated by B'110101101101011'", true},
	{"load data infile '/x.csv' into table t fields terminated by b'1101'", true},
	{"load data infile '/x.csv' into table t fields terminated by X''", true},
	{"load data infile '/x.csv' into table t fields enclosed by X'0D'", true},
	{"load data infile '/x.csv' into table t fields escaped by X'5C'", true},
	{"load data infile '/x.csv' into table t fields optionally enclosed by X'22'", true},
	{"load data infile '/x.csv' into table t columns terminated by X'09'", true},
	{"load data infile '/x.csv' into table t lines starting by X'6B6B' terminated by X'6B6B'", true},
	{"load data infile '/x.csv' into table t lines starting by B'1101' terminated by B'1101'", true},
	{"load data infile '/x.csv' into table t character set utf8mb4 fields terminated by X'09'", true},
	{"load data infile '/x.csv' replace into table t fields terminated by X'09'", true},
	// A hex literal with an odd number of digits is no literal at all.
	{"load data infile '/x.csv' into table t fields terminated by X'616'", false},

	// --- CHARACTER SET / CHARSET charset_name, a charset_name (ident_or_text
	// or the reserved word BINARY). This parser used to read a plain
	// identifier here, turning BINARY and a quoted name away.
	{"load data infile '/x.csv' into table t character set binary", true},
	{"load data infile '/x.csv' into table t character set 'utf8mb4'", true},
	{"load data infile '/x.csv' into table t charset binary", true},
	{"load data infile '/x.csv' into table t charset 'utf8mb4'", true},
	{"load data infile '/x.csv' into table t character set 1", false},

	// --- the item list: one or more, any order, repeats allowed ---
	{"load data infile '/x.csv' into table t fields terminated by 'a' terminated by 'b'", true},
	{"load data infile '/x.csv' into table t fields escaped by 'a' terminated by 'b'", true},
	{"load data infile '/x.csv' into table t fields enclosed by 'a' escaped by 'b' terminated by 'c'", true},
	{"load data infile '/x.csv' into table t lines terminated by 'b' starting by 'c'", true},
	{"load data infile '/x.csv' into table t lines starting by 'a' starting by 'b'", true},
	{"load data infile '/x.csv' into table t fields optionally enclosed by 'a' enclosed by 'b'", true},
	{"load data infile '/x.csv' into table t fields enclosed by 'a' optionally enclosed by 'b'", true},
	{`load data infile '/x.csv' into table t fields optionally enclosed by '"' terminated by ','`, true},
	// Neither clause names no item.
	{"load data infile '/x.csv' into table t fields", false},
	{"load data infile '/x.csv' into table t lines", false},
	{"load data infile '/x.csv' into table t fields lines terminated by 'a'", false},
	// ENCLOSED BY is not optional after OPTIONALLY.
	{"load data infile '/x.csv' into table t fields optionally enclosed ignore 1 lines", false},
	{"load data infile '/x.csv' into table t fields optionally `enclosed` by 'a'", false},

	// --- the column-and-user-variable list ---
	{"load data infile '/x.csv' into table t (a, @v) set b = @v", true},
	// A backquoted name is a column, even one spelled with an '@'.
	{"load data infile '/x.csv' into table t (`@v`)", true},
	{"load data infile '/x.csv' into table t (`@v`, a)", true},

	// --- what this engine does not have ---
	{"load data infile '/x.csv' into table t fields defined null by 'nil'", false},
	{"load data infile '/x.csv' into table t fields defined null by 'NULL' optionally enclosed", false},
	{"load data infile '/x.csv' format 'csv' into table t", false},
	{"load data infile '/x.csv' into table t with detached", false},
	{"load data infile '/x.csv' into table t with threads=10", false},
	{"import into t from '/file.csv'", false},
}

// TestLoadDataFieldForms holds this parser to the verdicts above.
func TestLoadDataFieldForms(t *testing.T) {
	for _, c := range loadDataFieldCases {
		t.Run(c.sql, func(t *testing.T) {
			_, err := Parse(c.sql)
			if accepts := err == nil; accepts != c.accept {
				t.Errorf("omni accepts=%v, want %v for %q: %v", accepts, c.accept, c.sql, err)
			}
		})
	}
}

// TestLoadDataFieldTree pins the value each option carries, which a verdict
// alone cannot show: a hex or bit literal reaches the node as the bytes it
// stands for, and a user variable in the column list is a variable and not a
// column whose name begins with '@'.
func TestLoadDataFieldTree(t *testing.T) {
	// X'6B6B' and B'110101101101011' are both the two bytes "kk": a bit
	// literal is padded on the left to fill whole bytes.
	ParseAndCompare(t, "LOAD DATA INFILE '/x.csv' INTO TABLE t FIELDS TERMINATED BY X'6B6B'",
		`{LOAD_DATA :loc 0 :infile "/x.csv" :table {TABLEREF :loc 37 :name t} :fields_terminated "kk"}`)
	ParseAndCompare(t, "LOAD DATA INFILE '/x.csv' INTO TABLE t FIELDS TERMINATED BY B'110101101101011'",
		`{LOAD_DATA :loc 0 :infile "/x.csv" :table {TABLEREF :loc 37 :name t} :fields_terminated "kk"}`)
	// A backquoted name is a column, even one spelled with an '@'. The
	// lexer strips the quotes, so only the source text tells the two apart.
	ParseAndCompare(t, "LOAD DATA INFILE '/x.csv' INTO TABLE t (`@v`)",
		"{LOAD_DATA :loc 0 :infile \"/x.csv\" :table {TABLEREF :loc 37 :name t} :columns {COLREF :loc 40 :col @v}}")
	// A plain ENCLOSED item after an OPTIONALLY one leaves the flag set
	// here, where TiDB clears it: mysql:8.0.46 and mariadb:11.8.8 both
	// write `7` for `INTO OUTFILE ... FIELDS OPTIONALLY ENCLOSED BY '"'
	// ENCLOSED BY '"'` and v8.5.5 writes `"7"` (read on 2026-11-04).
	ParseAndCompare(t, "LOAD DATA INFILE '/x.csv' INTO TABLE t FIELDS OPTIONALLY ENCLOSED BY 'a' ENCLOSED BY 'b'",
		`{LOAD_DATA :loc 0 :infile "/x.csv" :table {TABLEREF :loc 37 :name t} :fields_enclosed "b" :fields_optionally_enclosed true}`)
	ParseAndCompare(t, "LOAD DATA INFILE '/x.csv' INTO TABLE t (column1, @var1) SET column2 = @var1/100",
		`{LOAD_DATA :loc 0 :infile "/x.csv" :table {TABLEREF :loc 37 :name t} :columns {COLREF :loc 40 :col column1} {VAR :name var1 :loc 49} :set {ASSIGN :loc 60 :col {COLREF :loc 60 :col column2} :val {BINEXPR :op / :loc 70 :left {VAR :name var1 :loc 70} :right {INT_LIT :val 100 :loc 76}}}}`)
}
