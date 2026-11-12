package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// The name a column type's CHARACTER SET and COLLATE carry, as TiDB v8.5.5 reads it,
// every verdict measured with PREPARE on pingcap/tidb:v8.5.5 on 2026-11-12.
//
// CharsetName and CollationName are each `StringName | binaryType` (parser.y), and
// StringName is `stringLit | Identifier`: an identifier, a string literal in either
// quote, or the reserved word BINARY. The column type read only the first and the
// last, so `CHAR CHARSET 'binary'` was a syntax error here.
//
// The charsets and collations are real ones: TiDB looks a name up while it parses, and
// an unknown one answers 1115 or 1273 before a later syntax error can.

var charsetNameCases = []tableSampleCase{
	// OptBinary: CharsetKw CharsetName OptBinMod, and BINARY OptCharset.
	{"alter table t add c char charset 'binary'", true},
	{"alter table t add c char charset 'binary' binary", true},
	{"alter table t add c char binary charset 'binary'", true},
	{"alter table t add c char character set 'utf8mb4'", true},
	{"alter table t add c char char set 'utf8mb4'", true},
	{`alter table t add c char charset "utf8mb4"`, true},
	{"alter table t add c char charset 'UTF8MB4'", true},
	{"alter table t add c char charset utf8mb4", true},
	{"alter table t add c char charset `utf8mb4`", true},
	{"alter table t add c char charset binary", true},
	// The collation after it, and on its own.
	{"alter table t add c char collate 'utf8mb4_bin'", true},
	{"alter table t add c char collate 'binary'", true},
	{"alter table t add c char collate binary", true},
	{"alter table t add c char charset 'utf8mb4' collate 'utf8mb4_bin'", true},
	{"alter table t add c char charset utf8mb4 collate 'utf8mb4_bin'", true},
	{"alter table t add c char charset 'utf8mb4' collate utf8mb4_bin", true},
	{"alter table t add c char binary collate 'utf8mb4_bin'", true},
	{"alter table t add c char charset 'utf8mb4' binary collate 'utf8mb4_bin'", true},
	{"alter table t add c char charset 'utf8mb4'collate'utf8mb4_bin'", true},
	// Every type that takes the clause.
	{"alter table t add c varchar(10) charset 'utf8mb4'", true},
	{"alter table t add c text charset 'utf8mb4'", true},
	{"alter table t add c text collate 'utf8mb4_bin'", true},
	{"alter table t add c text ascii collate 'latin1_bin'", true},
	{"alter table t add c text byte collate 'binary'", true},
	{"alter table t add c enum('a') charset 'utf8mb4'", true},
	{"alter table t add c set('a') charset 'utf8mb4'", true},
	{"alter table t add c long charset 'utf8mb4'", true},
	{"alter table t add c long varchar charset 'utf8mb4'", true},
	{"alter table t add c nchar(1) charset 'utf8mb4'", true},
	{"alter table t add c national char(1) charset 'utf8mb4'", true},
	// Every statement that reads a column type, and the options after one.
	{"alter table t modify c char(10) charset 'utf8mb4'", true},
	{"alter table t change c c char(10) charset 'utf8mb4'", true},
	{"create table t2 (c char charset 'utf8mb4')", true},
	{"create table t2 (c char(1) charset 'utf8mb4' collate 'utf8mb4_bin')", true},
	{"create table t2 (c char(1) charset 'utf8mb4' not null)", true},
	{"alter table t add c char charset 'utf8mb4' default 'x'", true},
	{"alter table t add c char(1) charset 'utf8mb4' collate 'utf8mb4_bin' collate 'utf8mb4_bin'", true},

	// One string, and a plain one: no second string, no introducer, no other literal.
	{"alter table t add c char charset 'utf8' 'mb4'", false},
	{"alter table t add c char charset 'utf8mb4' 'x'", false},
	{"alter table t add c char charset N'utf8mb4'", false},
	{"alter table t add c char charset _utf8mb4'utf8mb4'", false},
	{"alter table t add c char charset 0x62", false},
	{"alter table t add c char charset x'62'", false},
	{"alter table t add c char charset 1", false},
	{"alter table t add c char charset null", false},
	{"alter table t add c char charset ?", false},
	// An Identifier, so no user variable and no introducer either.
	{"alter table t add c char charset @a", false},
	{"alter table t add c char charset _utf8mb4", false},
	{"alter table t add c char collate @a", false},
	{"alter table t add c char collate _utf8mb4", false},
	// The name is still required, and the clause's order unchanged.
	{"alter table t add c char charset", false},
	{"alter table t add c char charset default", false},
	{"alter table t add c char charset = 'utf8mb4'", false},
	{"alter table t add c char charset 'utf8mb4' charset 'utf8mb4'", false},
	{"alter table t add c char collate 'utf8mb4_bin' charset 'utf8mb4'", false},
	{"alter table t add c char collate 'utf8mb4_bin' binary", false},
	{"alter table t add c char charset 'utf8mb4' binary binary", false},
	{"alter table t add c char binary charset 'utf8mb4' binary", false},
	{"alter table t add c int charset 'utf8mb4'", false},
}

// TestCharsetNameForms holds omni to the recorded verdicts without a server.
func TestCharsetNameForms(t *testing.T) {
	for _, c := range charsetNameCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := parsesAsOneStatement(c.sql)
			if got != c.accept {
				t.Errorf("omni accepts=%v, want %v for %q: %v", got, c.accept, c.sql, err)
			}
		})
	}
}

// TestCharsetNameTree holds what a quoted name leaves in the tree: the name without
// its quotes, as an unquoted one leaves it.
func TestCharsetNameTree(t *testing.T) {
	for _, c := range []struct {
		sql              string
		charset, collate string
		binary           bool
	}{
		{"create table t (c char charset 'utf8mb4')", "utf8mb4", "", false},
		{`create table t (c char character set "utf8mb4" collate 'utf8mb4_bin')`, "utf8mb4", "utf8mb4_bin", false},
		{"create table t (c char binary charset 'binary')", "binary", "", true},
		{"create table t (c char charset 'binary' binary)", "binary", "", true},
		{"create table t (c text byte collate 'binary')", "binary", "binary", false},
	} {
		t.Run(c.sql, func(t *testing.T) {
			list, err := Parse(c.sql)
			if err != nil {
				t.Fatalf("Parse(%q): %v", c.sql, err)
			}
			typ := list.Items[0].(*nodes.CreateTableStmt).Columns[0].TypeName
			if typ.Charset != c.charset || typ.Collate != c.collate || typ.Binary != c.binary {
				t.Errorf("type has charset %q, collate %q, binary %v; want %q, %q, %v",
					typ.Charset, typ.Collate, typ.Binary, c.charset, c.collate, c.binary)
			}
		})
	}
}
