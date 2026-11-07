package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// A charset introducer before a hex or bit literal, as TiDB v8.5.5 reads it, every
// verdict measured with PREPARE on pingcap/tidb:v8.5.5 on 2026-11-07.
//
// "UNDERSCORE_CHARSET" hexLit and bitLit (parser.y:7625, :7642) are a hex or bit
// literal carrying a charset. The introducer must name a charset TiDB's lexer
// knows: _notacharset is a column, so `_notacharset 0x61` is a 1064. The word
// introduces a string too (StringLiteral), which omni already read.

type introducerLiteralCase struct {
	sql    string
	accept bool
}

var introducerLiteralCases = []introducerLiteralCase{
	{"select _utf8 0xD0B1", true},
	{"select _utf8 X'D0B1'", true},
	{"select _utf8 0b1101000010110001", true},
	{"select _utf8 B'1101000010110001'", true},
	{"select _latin1 0x61", true},
	{"select _utf8mb4 x'61'", true},
	{"select _utf8 0x61 + 1", true},
	{"select _utf8 0xD0B1 from t", true},
	{"select _notacharset 0x61", false},
	{"select `_utf8` 0x61", false},
	{"select _utf8 'str'", true},
}

// TestIntroducerLiteralForms holds omni to the recorded verdicts without a server.
func TestIntroducerLiteralForms(t *testing.T) {
	for _, c := range introducerLiteralCases {
		t.Run(c.sql, func(t *testing.T) {
			got, err := parsesAsOneStatement(c.sql)
			if got != c.accept {
				t.Errorf("Parse(%q) accepts=%v, want %v (TiDB %s): %v", c.sql, got, c.accept, tidbCorpusVersion, err)
			}
		})
	}
}

// TestIntroducerLiteralTrees checks the introducer is kept on the literal and the
// Loc spans both tokens.
func TestIntroducerLiteralTrees(t *testing.T) {
	item := func(t *testing.T, sql string) nodes.ExprNode {
		t.Helper()
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", sql, err)
		}
		return list.Items[0].(*nodes.SelectStmt).TargetList[0]
	}
	t.Run("hex", func(t *testing.T) {
		h, ok := item(t, "select _utf8 0xD0B1").(*nodes.HexLit)
		if !ok || h.Value != "0xD0B1" || h.Charset != "_utf8" {
			t.Fatalf("got %#v, want HexLit 0xD0B1 _utf8", item(t, "select _utf8 0xD0B1"))
		}
		if got := "select _utf8 0xD0B1"[h.Loc.Start:h.Loc.End]; got != "_utf8 0xD0B1" {
			t.Errorf("Loc covers %q, want %q", got, "_utf8 0xD0B1")
		}
	})
	t.Run("bit", func(t *testing.T) {
		b, ok := item(t, "select _latin1 B'101'").(*nodes.BitLit)
		if !ok || b.Value != "101" || b.Charset != "_latin1" {
			t.Fatalf("got %#v, want BitLit 101 _latin1", item(t, "select _latin1 B'101'"))
		}
	})
	t.Run("plain hex keeps no charset", func(t *testing.T) {
		h, ok := item(t, "select 0x61").(*nodes.HexLit)
		if !ok || h.Charset != "" {
			t.Errorf("got %#v, want a HexLit with no charset", item(t, "select 0x61"))
		}
	})
}
