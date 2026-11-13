package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// Where BINARY stands in a character type, as OceanBase 4.3.5 reads it, every verdict
// measured with PREPARE on oceanbase-ce:4.3.5.6 on 2026-11-13.
//
// The grammar is `VARCHAR string_length_i opt_binary opt_charset`, and so for CHARACTER,
// CHARACTER VARYING, the text types, ENUM and SET: BINARY comes before the character set
// and nowhere after it, where MySQL takes it on either side. A national type has
// opt_binary alone. COLLATE is a column attribute, which may follow either.
//
// Not in this table, because omni still disagrees: a character set on a national type and
// BINARY or a character set on a BLOB, all of which OceanBase refuses and this reader takes.
//
// Ref: data_type in OceanBase's sql_parser_mysql_mode.y (v4.3.5_CE_BP6, lines 6786 to 7017).
var typeBinaryCases = []parserFormCase{
	// BINARY alone, and before a character set in each of its spellings.
	{"alter table t add c varchar(10) binary", true},
	{"alter table t add c varchar(10) binary character set ascii", true},
	{"alter table t add c varchar(10) binary charset ascii", true},
	{"alter table t add c varchar(10) binary char set ascii", true},
	{"alter table t add c varchar(10) binary character set 'ascii'", true},
	{"alter table t add c varchar(10) binary charset binary", true},
	// A collation and the other attributes follow it.
	{"alter table t add c varchar(10) binary collate utf8mb4_bin", true},
	{"alter table t add c varchar(10) binary character set ascii collate ascii_general_ci", true},
	{"alter table t add c varchar(10) binary not null", true},
	{"alter table t add c varchar(10) binary character set ascii not null default 'x'", true},
	// Every type that takes it.
	{"alter table t add c char binary", true},
	{"alter table t add c char(5) binary character set ascii", true},
	{"alter table t add c character varying(5) binary charset ascii", true},
	{"alter table t add c tinytext binary", true},
	{"alter table t add c text binary character set ascii", true},
	{"alter table t add c text(100) binary charset ascii", true},
	{"alter table t add c mediumtext binary charset ascii", true},
	{"alter table t add c longtext binary charset ascii", true},
	{"alter table t add c long binary", true},
	{"alter table t add c long binary charset ascii", true},
	{"alter table t add c long varchar binary charset ascii", true},
	{"alter table t add c enum('x','y') binary", true},
	{"alter table t add c enum('x','y') binary character set ascii", true},
	{"alter table t add c set('x','y') binary charset ascii", true},
	{"alter table t add c nchar(5) binary", true},
	{"alter table t add c nvarchar(5) binary", true},
	{"alter table t add c national char(5) binary", true},
	{"alter table t add c national varchar(5) binary", true},
	{"alter table t add c nchar varchar(5) binary", true},
	{"alter table t add c national char varying(5) binary", true},
	// Every statement that reads a column type.
	{"create table t2 (c varchar(10) binary character set ascii)", true},
	{"alter table t modify c char(10) binary charset ascii", true},
	{"alter table t change c c char(10) binary charset ascii", true},

	// After the character set, or a collation, or another attribute, or itself.
	{"alter table t add c varchar(10) character set ascii binary", false},
	{"alter table t add c varchar(10) charset ascii binary", false},
	{"alter table t add c varchar(10) char set ascii binary", false},
	{"alter table t add c varchar(10) collate utf8mb4_bin binary", false},
	{"alter table t add c varchar(10) character set ascii collate ascii_bin binary", false},
	{"alter table t add c varchar(10) not null binary", false},
	{"alter table t add c varchar(10) binary binary", false},
	{"alter table t add c varchar(10) binary character set ascii binary", false},
	{"alter table t add c char(5) character set ascii binary", false},
	{"alter table t add c character varying(5) charset ascii binary", false},
	{"alter table t add c text character set ascii binary", false},
	{"alter table t add c longtext charset ascii binary", false},
	{"alter table t add c long charset ascii binary", false},
	{"alter table t add c enum('x','y') character set ascii binary", false},
	{"alter table t add c set('x','y') charset ascii binary", false},
	{"create table t2 (c varchar(10) character set ascii binary)", false},
	{"alter table t modify c char(10) charset ascii binary", false},
}

// TestTypeBinaryForms holds omni to the recorded verdicts without a server.
func TestTypeBinaryForms(t *testing.T) {
	for _, c := range typeBinaryCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := parsesAsOneStatement(c.sql)
			if got != c.accept {
				t.Errorf("omni accepts=%v, want %v for %q: %v", got, c.accept, c.sql, err)
			}
		})
	}
}

// TestTypeBinaryTree holds what BINARY leaves in the tree, with a character set after it
// and without one.
func TestTypeBinaryTree(t *testing.T) {
	for _, c := range []struct {
		sql              string
		charset, collate string
	}{
		{"create table t (c varchar(10) binary)", "", ""},
		{"create table t (c varchar(10) binary character set ascii)", "ascii", ""},
		{"create table t (c char(5) binary charset 'latin1' collate latin1_bin)", "latin1", "latin1_bin"},
	} {
		t.Run(c.sql, func(t *testing.T) {
			list, err := Parse(c.sql)
			if err != nil {
				t.Fatalf("Parse(%q): %v", c.sql, err)
			}
			typ := list.Items[0].(*nodes.CreateTableStmt).Columns[0].TypeName
			if !typ.Binary || typ.Charset != c.charset || typ.Collate != c.collate {
				t.Errorf("type has binary %v, charset %q and collate %q, want true, %q and %q",
					typ.Binary, typ.Charset, typ.Collate, c.charset, c.collate)
			}
			if got, want := c.sql[typ.Loc.Start:typ.Loc.End], c.sql[len("create table t (c "):len(c.sql)-1]; got != want {
				t.Errorf("type spans %q, want %q", got, want)
			}
		})
	}
}
