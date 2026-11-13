package parser

import "testing"

// Legacy column type words as TiDB v8.5.5 reads them, every verdict measured with
// PREPARE on pingcap/tidb:v8.5.5 on 2026-11-06.

type legacyTypeCase struct {
	sql    string
	accept bool
}

var legacyTypeCases = []legacyTypeCase{
	// OptCharsetWithOptBinary (parser.y:13370): ASCII, UNICODE or BYTE after a TEXT
	// type, ENUM, SET, LONG or LONG VARCHAR, in place of a charset clause.
	{"create table t (a text ascii)", true},
	{"create table t (a tinytext ascii)", true},
	{"create table t (a mediumtext ascii)", true},
	{"create table t (a longtext ascii)", true},
	{"create table t (a text(10) byte)", true},
	{"create table t (a text unicode)", true},
	{"create table t (a text ascii binary)", false},
	{"create table t (a text ascii charset latin1)", false},
	{"create table t (a text binary ascii)", false},
	{"create table t (a text ascii collate latin1_bin)", true},
	{"create table t (a varchar(10) ascii)", false},
	{"create table t (a char(10) byte)", false},
	{"create table t (a enum('x') ascii)", true},
	{"create table t (a set('x') byte)", true},
	{"create table t (a long ascii)", true},
	{"create table t (a long byte)", true},
	{"create table t (a long varchar ascii)", true},
	{"create table t (a long varbinary ascii)", false},
	{"create table t (a blob ascii)", false},
	{"create table t (ascii int, byte int)", true},
	{"create table t (a text ascii collate)", false},
	{"create table t (a text byte collate not null)", false},
	{"create table t (a text collate)", false},
	{"create table t (a text collate not null)", false},
	{"create table t (a varchar(3) collate)", false},
	{"create table t (a varchar(3) charset)", false},
	{"create table t (a varchar(3) character set)", false},
	{"create table t (a varchar(3) charset not null)", false},
	{"create table t (a enum('x') charset)", false},
	{"create table t (a long collate)", false},
	{"create table t (a varchar(3) charset binary collate binary)", true},

	// FieldOpts (parser.y:13472): any run of UNSIGNED, SIGNED and ZEROFILL after a
	// numeric type, BOOL and YEAR.
	{"create table t (y year(4) unsigned zerofill zerofill)", true},
	{"create table t (y year signed unsigned zerofill)", true},
	{"create table t (a int signed unsigned)", true},
	{"create table t (a int zerofill unsigned)", true},
	{"create table t (a int unsigned unsigned)", true},
	{"create table t (a decimal(5,2) zerofill unsigned)", true},
	{"create table t (a double unsigned signed)", true},
	{"create table t (a bool unsigned)", true},
	{"create table t (a boolean zerofill)", true},
	{"create table t (a bit(3) unsigned)", false},
	{"create table t (a varchar(3) unsigned)", false},

	// "STORAGE" StorageMedia (parser.y:3768): DEFAULT, DISK or MEMORY, parsed and
	// ignored.
	{"create table t (a int storage disk)", true},
	{"create table t (a int storage default)", true},
	{"create table t (a int storage memory)", true},
	{"create table t (a int storage)", false},
	{"create table t (a int storage foo)", false},
	{"create table t (a int storage disk storage memory)", true},
	{"alter table t add column a int storage disk", true},
	{"alter table t modify a int storage memory", true},
}

// TestLegacyTypeForms holds omni to the recorded verdicts without a server.
func TestLegacyTypeForms(t *testing.T) {
	for _, c := range legacyTypeCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := parsesAsOneStatement(c.sql)
			if got != c.accept {
				t.Errorf("Parse(%q) accepts=%v, want %v (TiDB %s): %v", c.sql, got, c.accept, tidbRelease, err)
			}
		})
	}
}
