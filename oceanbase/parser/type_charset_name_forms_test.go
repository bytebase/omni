package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// The name a column type's CHARACTER SET and COLLATE carry, as OceanBase 4.3.5 reads
// it, every verdict measured with PREPARE on oceanbase-ce:4.3.5.6 on 2026-11-12.
//
// charset_name is `NAME_OB | STRING_VALUE | BINARY` and collation_name `NAME_OB |
// STRING_VALUE`: a name may be a string literal in either quote. The column type read
// identifiers only, so `CHAR CHARSET 'binary'` was a syntax error here.
//
// The charsets and collations are real ones, since the server looks a name up.
//
// Not in this table, because omni still disagrees, quoted or not: where BINARY stands
// (OceanBase takes it before the character set and not after), BINARY as a collation,
// a character set on a national type, and a clause with its name left off, which this
// reader lets go.
//
// Ref: charset_name and collation_name in OceanBase 4.3.5 sql_parser_mysql_mode.y.

var typeCharsetNameCases = []parserFormCase{
	{"alter table t add c char charset 'binary'", true},
	{"alter table t add c char character set 'binary'", true},
	{"alter table t add c char char set 'binary'", true},
	{"alter table t add c char charset 'utf8mb4'", true},
	{`alter table t add c char charset "utf8mb4"`, true},
	{"alter table t add c char charset 'UTF8MB4'", true},
	{"alter table t add c char charset utf8mb4", true},
	{"alter table t add c char charset `utf8mb4`", true},
	{"alter table t add c char charset binary", true},
	// The collation after it, and on its own.
	{"alter table t add c char collate 'utf8mb4_bin'", true},
	{"alter table t add c char collate 'binary'", true},
	{"alter table t add c char charset 'utf8mb4' collate 'utf8mb4_bin'", true},
	{"alter table t add c char charset utf8mb4 collate 'utf8mb4_bin'", true},
	{"alter table t add c char charset 'utf8mb4' collate utf8mb4_bin", true},
	{"alter table t add c char charset 'utf8mb4'collate'utf8mb4_bin'", true},
	// Every type that takes the clause.
	{"alter table t add c varchar(10) charset 'utf8mb4'", true},
	{"alter table t add c varchar(10) character set 'utf8mb4' collate 'utf8mb4_bin'", true},
	{"alter table t add c text charset 'utf8mb4'", true},
	{"alter table t add c text collate 'utf8mb4_bin'", true},
	{"alter table t add c enum('a') charset 'utf8mb4'", true},
	{"alter table t add c set('a') charset 'utf8mb4'", true},
	{"alter table t add c long charset 'utf8mb4'", true},
	{"alter table t add c long varchar charset 'utf8mb4'", true},
	// Every statement that reads a column type, and the options after one.
	{"alter table t modify c char(10) charset 'utf8mb4'", true},
	{"alter table t change c c char(10) charset 'utf8mb4'", true},
	{"create table t2 (c char charset 'utf8mb4')", true},
	{"create table t2 (c char(1) charset 'utf8mb4' collate 'utf8mb4_bin')", true},
	{"create table t2 (c char(1) charset 'utf8mb4' not null)", true},
	{"alter table t add c char charset 'utf8mb4' default 'x'", true},

	// A plain string: no introducer and no other literal.
	{"alter table t add c char charset N'utf8mb4'", false},
	{"alter table t add c char charset _utf8mb4'utf8mb4'", false},
	{"alter table t add c char charset 0x62", false},
	{"alter table t add c char charset x'62'", false},
	{"alter table t add c char charset 1", false},
	{"alter table t add c char charset ?", false},
	{"alter table t add c char charset default", false},
	{"alter table t add c char charset = 'utf8mb4'", false},
	{"alter table t add c char charset 'utf8mb4' charset 'utf8mb4'", false},
	{"alter table t add c char collate 'utf8mb4_bin' charset 'utf8mb4'", false},
	{"alter table t add c int charset 'utf8mb4'", false},
}

// TestTypeCharsetNameForms holds omni to the recorded verdicts without a server.
func TestTypeCharsetNameForms(t *testing.T) {
	for _, c := range typeCharsetNameCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := parsesAsOneStatement(c.sql)
			if got != c.accept {
				t.Errorf("omni accepts=%v, want %v for %q: %v", got, c.accept, c.sql, err)
			}
		})
	}
}

// TestTypeCharsetNameTree holds what a quoted name leaves in the tree: the name without
// its quotes, as an unquoted one leaves it.
func TestTypeCharsetNameTree(t *testing.T) {
	for _, c := range []struct {
		sql              string
		charset, collate string
	}{
		{"create table t (c char charset 'utf8mb4')", "utf8mb4", ""},
		{`create table t (c char character set "utf8mb4" collate 'utf8mb4_bin')`, "utf8mb4", "utf8mb4_bin"},
		{"create table t (c text collate 'utf8mb4_bin')", "", "utf8mb4_bin"},
		{"create table t (c char charset 'binary')", "binary", ""},
	} {
		t.Run(c.sql, func(t *testing.T) {
			list, err := Parse(c.sql)
			if err != nil {
				t.Fatalf("Parse(%q): %v", c.sql, err)
			}
			typ := list.Items[0].(*nodes.CreateTableStmt).Columns[0].TypeName
			if typ.Charset != c.charset || typ.Collate != c.collate {
				t.Errorf("type has charset %q and collate %q, want %q and %q",
					typ.Charset, typ.Collate, c.charset, c.collate)
			}
		})
	}
}
