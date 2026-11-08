package parser

import "testing"

// The forms SELECT ... INTO OUTFILE's FIELDS clause takes, and the forms it
// does not. TiDB reads the same Fields rule here that LOAD DATA reads, so the
// field items are a non-empty list in any order, a repeat allowed with the
// last of a kind winning, and TiDB's own DEFINED NULL BY among them -- none of
// which a fixed one-of-each read accepted.
//
// Every case was measured against pingcap/tidb:v8.5.5 through PREPARE on
// 2026-11-07, and TestOutfileFieldsFormsMatchTiDB keeps the verdicts a
// recording rather than a guess.
//
// Ref: SelectStmtIntoOption, Fields, FieldItem and FieldTerminator in
// pingcap/tidb pkg/parser/parser.y.

var outfileFieldsFormCases = []loadDataFormCase{
	// --- the list is one or more items, in any order, and may repeat ---
	{"select * into outfile '/x' fields terminated by 'a'", true},
	{"select * into outfile '/x' fields terminated by 'a' terminated by 'b'", true},
	{"select * into outfile '/x' fields optionally enclosed by 'a' optionally enclosed by 'b'", true},
	{"select * into outfile '/x' fields terminated by 'a' enclosed by 'b' escaped by 'c'", true},
	{"select * into outfile '/x' fields enclosed by 'a' optionally enclosed by 'b'", true},
	// FIELDS with nothing after it names no item, and LINES is not an item.
	{"select * into outfile '/x' fields", false},
	{"select * into outfile '/x' fields lines terminated by 'a'", false},
	// ENCLOSED BY is not optional after OPTIONALLY, and the value is one byte.
	{"select * into outfile '/x' fields optionally", false},
	{"select * into outfile '/x' fields optionally enclosed", false},
	{"select * into outfile '/x' fields optionally enclosed by 'ab'", false},

	// --- DEFINED NULL BY, which TiDB has here as in LOAD DATA ---
	{"select * into outfile '/x' fields defined null by 'n'", true},
	{"select * into outfile '/x' fields defined null by 'n' optionally enclosed", true},
	{"select * into outfile '/x' fields defined null by 'a' defined null by 'b'", true},
	{"select * into outfile '/x' fields terminated by 'a' defined null by 'b' enclosed by 'c'", true},
	// A BY after OPTIONALLY ENCLOSED makes it the next item, an ENCLOSED one.
	{"select * into outfile '/x' fields defined null by 'a' optionally enclosed by 'b'", true},
	// OPTIONALLY with no ENCLOSED after it is still incomplete.
	{"select * into outfile '/x' fields defined null by 'a' optionally", false},
	// Every word of DEFINED NULL BY is a word, not a name.
	{"select * into outfile '/x' fields `defined` null by 'a'", false},
}

// TestOutfileFieldsForms holds this parser to the verdicts above.
func TestOutfileFieldsForms(t *testing.T) {
	for _, c := range outfileFieldsFormCases {
		t.Run(c.sql, func(t *testing.T) {
			_, err := Parse(c.sql)
			if accepts := err == nil; accepts != c.accept {
				t.Errorf("omni %ss %q, TiDB v8.5.5 %ss it: %v",
					verdictName(accepts), c.sql, verdictName(c.accept), err)
			}
		})
	}
}
