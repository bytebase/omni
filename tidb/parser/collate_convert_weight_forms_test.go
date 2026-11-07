package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// COLLATE, CONVERT(... USING ...) and WEIGHT_STRING(... AS ...) as TiDB v8.5.5
// reads them. These close the collate/convert/weight_string corpus-gap
// cluster: the reserved word BINARY is a collation name (ColumnOption "COLLATE"
// CollationName and SimpleExpr "COLLATE" CollationName, where CollationName is
// StringName | binaryType) and a charset name (CONVERT "(" Expression "USING"
// CharsetName ")", where CharsetName is StringName | binaryType), and
// WEIGHT_STRING carries its own AS clause (FunctionCallNonKeyword weight_string
// arms, parser.y:8367,8374,8381).
//
// Every expectation here was measured against TiDB v8.5.5 through PREPARE;
// TestCollateConvertWeightFormsMatchTiDB holds the file to the server.
//
// Ref: TiDB v8.5.5 parser.y CollationName/CharsetName (11023-11071),
// ColumnOption (3760), SimpleExpr (7879,8056), FunctionCallNonKeyword (8367).

type collateConvertWeightCase struct {
	sql    string
	accept bool
}

var collateConvertWeightCases = []collateConvertWeightCase{
	// COLLATE's collation name: an identifier, a string, or the reserved word
	// BINARY, as a column option and as an operand postfix. TiDB takes BINARY
	// in both, where MySQL takes it only as the column option and MariaDB in
	// neither.
	{"ALTER TABLE t ADD c JSON COLLATE BINARY", true},
	{"ALTER TABLE t ADD c JSON COLLATE utf8mb4_bin", true},
	{"CREATE TABLE t (a INT COLLATE BINARY)", true},
	{"CREATE TABLE t (a INT COLLATE utf8mb4_bin)", true},
	{"CREATE TABLE t (a INT COLLATE 'utf8mb4_bin')", true},
	{"DO x COLLATE BINARY", true},
	{"DO x COLLATE utf8mb4_bin", true},
	{"DO x COLLATE 'utf8mb4_bin'", true},
	{"DO x COLLATE", false},

	// CONVERT(expr USING charset): the charset name carries the same BINARY
	// arm, beside the identifier and string forms.
	{"DO CONVERT(x USING BINARY)", true},
	{"DO CONVERT(x USING utf8mb4)", true},
	{"DO CONVERT(x USING 'utf8mb4')", true},
	{"DO CONVERT(x USING)", false},

	// WEIGHT_STRING(expr [AS CHAR(n) | AS CHARACTER(n) | AS BINARY(n)]). CHAR
	// and CHARACTER are one keyword, the length is required, and TiDB takes 0
	// where MySQL and MariaDB answer 1064. A general cast type is not an AS
	// clause.
	{"DO WEIGHT_STRING(x)", true},
	{"DO WEIGHT_STRING(x AS CHAR(1))", true},
	{"DO WEIGHT_STRING(x AS CHARACTER(1))", true},
	{"DO WEIGHT_STRING(x AS BINARY(1))", true},
	{"DO WEIGHT_STRING(x AS CHAR(0))", true},
	{"DO WEIGHT_STRING(x AS CHAR)", false},
	{"DO WEIGHT_STRING(x AS BINARY)", false},
	{"DO WEIGHT_STRING(x AS DECIMAL(1))", false},
	{"DO WEIGHT_STRING(x AS)", false},
}

func TestCollateConvertWeightFormsAccepted(t *testing.T) {
	for _, c := range collateConvertWeightCases {
		if !c.accept {
			continue
		}
		t.Run(c.sql, func(t *testing.T) {
			if _, err := Parse(c.sql); err != nil {
				t.Fatalf("Parse(%q) error: %v", c.sql, err)
			}
		})
	}
}

func TestCollateConvertWeightFormsRejected(t *testing.T) {
	for _, c := range collateConvertWeightCases {
		if c.accept {
			continue
		}
		t.Run(c.sql, func(t *testing.T) {
			if _, err := Parse(c.sql); err == nil {
				t.Fatalf("Parse(%q) was accepted, but TiDB %s rejects it", c.sql, tidbCorpusVersion)
			}
		})
	}
}

// TestWeightStringASClauseAST checks that the AS clause lands on the call node,
// so a consumer and the deparser can see CHAR(n) / BINARY(n).
func TestWeightStringASClauseAST(t *testing.T) {
	cases := []struct {
		sql     string
		wantAs  string
		wantLen uint64
	}{
		{"DO WEIGHT_STRING(x AS CHAR(5))", "CHAR", 5},
		{"DO WEIGHT_STRING(x AS CHARACTER(3))", "CHAR", 3},
		{"DO WEIGHT_STRING(x AS BINARY(7))", "BINARY", 7},
		{"DO WEIGHT_STRING(x)", "", 0},
	}
	for _, c := range cases {
		t.Run(c.sql, func(t *testing.T) {
			list, err := Parse(c.sql)
			if err != nil {
				t.Fatalf("Parse(%q): %v", c.sql, err)
			}
			do, ok := list.Items[0].(*nodes.DoStmt)
			if !ok {
				t.Fatalf("statement is %T, want *DoStmt", list.Items[0])
			}
			fc, ok := do.Exprs[0].(*nodes.FuncCallExpr)
			if !ok {
				t.Fatalf("expr is %T, want *FuncCallExpr", do.Exprs[0])
			}
			if fc.WeightAs != c.wantAs || fc.WeightAsLen != c.wantLen {
				t.Errorf("WeightAs=%q WeightAsLen=%d, want %q %d", fc.WeightAs, fc.WeightAsLen, c.wantAs, c.wantLen)
			}
		})
	}
}
