package parser

import (
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

// collateCases pin where SQL Server 2022 accepts a COLLATE clause, as SET
// PARSEONLY ON reports it. {a}, {b} and {c} stand for three collation names.
//
//   - An expression operand takes one COLLATE: a second one straight after
//     it is Msg 156, in every expression position. Parentheses make a new
//     operand that takes its own, and so does the time zone of AT TIME ZONE.
//     COLLATE belongs to the operand before it, so any operator may follow
//     it: 'a' COLLATE x + 'b'.
//   - A column definition takes COLLATE once, directly after the data type:
//     after any other column option it is Msg 156.
//   - The collation name is a regular identifier: bracketed, quoted and
//     string forms are Msg 102.
var collateCases = []struct {
	sql    string
	accept bool
}{
	{"SELECT x COLLATE {a} FROM t", true},
	{"SELECT x COLLATE database_default FROM t", true},
	{"SELECT x COLLATE {a}, y COLLATE {b} FROM t", true},
	{"SELECT (x COLLATE {a}) COLLATE {b} FROM t", true},
	{"SELECT ((x COLLATE {a}) COLLATE {b}) COLLATE {c} FROM t", true},
	{"SELECT -(x COLLATE {a}) COLLATE {b} FROM t", true},
	{"SELECT COALESCE(x COLLATE {a}, y) COLLATE {b} FROM t", true},
	{"SELECT CASE WHEN 1 = 1 THEN 'a' COLLATE {a} END COLLATE {b}", true},
	{"SELECT x COLLATE {a} AT TIME ZONE 'UTC' COLLATE {b} FROM t", true},
	{"SELECT d AT TIME ZONE 'UTC' COLLATE {a} AT TIME ZONE 'UTC' COLLATE {b} FROM t", true},
	{"SELECT 1 FROM t WHERE x COLLATE {a} = y COLLATE {b}", true},
	{"SELECT x FROM t ORDER BY (x COLLATE {a}) COLLATE {b}", true},

	// COLLATE binds to the operand before it, so an operator may follow it.
	{"SELECT 'a' COLLATE {a} + 'b'", true},
	{"SELECT 'a' COLLATE {a} + 'b' COLLATE {b}", true},
	{"SELECT x COLLATE {a} * 2 FROM t", true},
	{"SELECT x COLLATE {a} - y COLLATE {b} FROM t", true},
	{"SELECT x COLLATE {a} & 1 FROM t", true},
	{"SELECT -x COLLATE {a} FROM t", true},
	{"SELECT CAST(x AS VARCHAR(10)) COLLATE {a} + 'b' FROM t", true},
	{"SELECT (SELECT 'a') COLLATE {a} + 'b'", true},
	{"SELECT d AT TIME ZONE 'UTC' COLLATE {a} + 'b' FROM t", true},
	{"SELECT x COLLATE {a} AT TIME ZONE 'UTC' + 'x' COLLATE {b} FROM t", true},
	{"SELECT 1 FROM t WHERE x COLLATE {a} + 'a' = y", true},
	{"SELECT 1 FROM t WHERE x COLLATE {a} LIKE 'a' COLLATE {b} + '%'", true},
	{"UPDATE t SET x = y COLLATE {a} + 'z'", true},
	{"CREATE TABLE t3 (a VARCHAR(10) DEFAULT 'x' COLLATE {a} + 'y')", true},

	{"SELECT 'a' COLLATE {a} COLLATE {b} + 'b'", false},
	{"SELECT 1, COLLATE {a}", false},
	{"SELECT 1 FROM t WHERE COLLATE {a} = 1", false},

	{"SELECT x COLLATE {a} COLLATE {b} FROM t", false},
	{"SELECT x COLLATE {a} COLLATE {a} FROM t", false},
	{"SELECT x COLLATE {a} COLLATE {b} COLLATE {c} FROM t", false},
	{"SELECT x COLLATE {a} -- comment\n  COLLATE {b} FROM t", false},
	{"SELECT -x COLLATE {a} COLLATE {b} FROM t", false},
	{"SELECT 'a' + 'b' COLLATE {a} COLLATE {b}", false},
	{"SELECT (x + y) COLLATE {a} COLLATE {b} FROM t", false},
	{"SELECT CAST(x AS VARCHAR(10)) COLLATE {a} COLLATE {b} FROM t", false},
	{"SELECT CASE WHEN 1 = 1 THEN 'a' END COLLATE {a} COLLATE {b}", false},
	{"SELECT (SELECT 'a') COLLATE {a} COLLATE {b}", false},
	{"SELECT d AT TIME ZONE 'UTC' COLLATE {a} COLLATE {b} FROM t", false},
	{"SELECT 1 FROM t WHERE x COLLATE {a} COLLATE {b} = y", false},
	{"SELECT x FROM t ORDER BY x COLLATE {a} COLLATE {b}", false},
	{"UPDATE t SET x = y COLLATE {a} COLLATE {b}", false},
	{"DECLARE @v VARCHAR(10) = 'x' COLLATE {a} COLLATE {b}", false},

	{"CREATE TABLE t3 (a VARCHAR(10) COLLATE {a} SPARSE NULL)", true},
	{"CREATE TABLE t3 (a VARCHAR(10) COLLATE {a} PRIMARY KEY)", true},
	{"CREATE TABLE t3 (a VARCHAR(10) COLLATE {a} DEFAULT 'x' COLLATE {b})", true},
	{"CREATE TABLE t3 (a VARCHAR(10) DEFAULT ('x' COLLATE {a}) COLLATE {b})", true},
	{"CREATE TABLE t3 (x VARCHAR(10), a AS (x COLLATE {a}) COLLATE {b})", true},
	{"ALTER TABLE t ALTER COLUMN a VARCHAR(10) COLLATE {a} NOT NULL", true},
	{"DECLARE @t TABLE (a VARCHAR(10) COLLATE {a} NOT NULL)", true},
	{"CREATE DATABASE d2 COLLATE {a}", true},
	{"ALTER DATABASE testdb COLLATE {a}", true},

	{"CREATE TABLE t3 (a VARCHAR(10) COLLATE {a} COLLATE {b})", false},
	{"CREATE TABLE t3 (a VARCHAR(10) DEFAULT 'x' COLLATE {a} COLLATE {b})", false},
	{"CREATE TABLE t3 (a VARCHAR(10) NOT NULL COLLATE {a})", false},
	{"CREATE TABLE t3 (a VARCHAR(10) COLLATE {a} NOT NULL COLLATE {b})", false},
	{"CREATE TABLE t3 (a VARCHAR(10) PRIMARY KEY COLLATE {a})", false},
	{"CREATE TABLE t3 (a VARCHAR(10) CHECK (a <> '') COLLATE {a})", false},
	{"CREATE TABLE t3 (x VARCHAR(10), a AS x COLLATE {a} COLLATE {b})", false},
	{"ALTER TABLE t ADD c VARCHAR(10) NOT NULL COLLATE {a}", false},
	{"DECLARE @t TABLE (a VARCHAR(10) NOT NULL COLLATE {a})", false},
	{"CREATE TYPE tt AS TABLE (a VARCHAR(10) COLLATE {a} COLLATE {b})", false},

	{"SELECT x COLLATE [Latin1_General_CI_AS] FROM t", false},
	{`SELECT x COLLATE "Latin1_General_CI_AS" FROM t`, false},
	{"SELECT x COLLATE [database_default] FROM t", false},
	{"CREATE TABLE t3 (a VARCHAR(10) COLLATE [Latin1_General_CI_AS])", false},
	{"ALTER TABLE t ALTER COLUMN a VARCHAR(10) COLLATE [Latin1_General_CI_AS]", false},
	{"CREATE DATABASE d2 COLLATE [Latin1_General_CI_AS]", false},
	{"ALTER DATABASE testdb COLLATE [Latin1_General_CI_AS]", false},
}

var collationNames = strings.NewReplacer(
	"{a}", "Latin1_General_CI_AS", "{b}", "Latin1_General_BIN", "{c}", "SQL_Latin1_General_CP1_CI_AS")

// TestCollateClauseOracle: omni accepts exactly the COLLATE clauses SQL
// Server accepts, and reports a rejected one where SqlScriptDOM does.
// omni used to take any number of COLLATE clauses in a row, anywhere in a
// column definition, and a bracketed or quoted collation name.
func TestCollateClauseOracle(t *testing.T) {
	oracle := startParserOracle(t)
	h := getHarness(t)
	for _, tc := range collateCases {
		sql := collationNames.Replace(tc.sql)
		t.Run(tc.sql, func(t *testing.T) {
			serverErr, err := oracle.parseError(sql)
			if err != nil {
				t.Fatal(err)
			}
			if (serverErr == nil) != tc.accept {
				t.Fatalf("SQL Server: %v, want accept = %v", serverErr, tc.accept)
			}
			_, omniErr := Parse(sql)
			if (omniErr == nil) != tc.accept {
				t.Fatalf("omni: %v, want accept = %v", omniErr, tc.accept)
			}
			if tc.accept {
				return
			}
			errs, _ := h.Shape(t, sql)["errors"].([]any)
			if len(errs) == 0 {
				t.Fatal("SqlScriptDOM accepted it")
			}
			first, _ := errs[0].(map[string]any)
			line, _ := first["line"].(float64)
			column, _ := first["column"].(float64)
			var pe *ParseError
			if !errors.As(omniErr, &pe) {
				t.Fatalf("omni error %v is not a *ParseError", omniErr)
			}
			if want := lineColumnOffset(sql, int(line), int(column)); pe.Position != want {
				t.Errorf("omni reports offset %d, SqlScriptDOM line %d column %d (offset %d)",
					pe.Position, int(line), int(column), want)
			}
		})
	}
}

// collatePrecedenceCases pin the operator shape of expressions that mix
// COLLATE, AT TIME ZONE, and the arithmetic and bitwise operators, in the
// notation of exprShape, which the SqlScriptDOM harness prints too. COLLATE
// binds tighter than every operator, unary ones included; AT TIME ZONE binds
// like + and takes a multiplication-level time zone.
var collatePrecedenceCases = []struct {
	sql  string
	want string
}{
	{"SELECT 'a' COLLATE {a} + 'b'", "(+ (collate # {a}) #)"},
	{"SELECT 'a' + 'b' COLLATE {a}", "(+ # (collate # {a}))"},
	{"SELECT 'a' COLLATE {a} + 'b' COLLATE {b} + 'c'", "(+ (+ (collate # {a}) (collate # {b})) #)"},
	{"SELECT x * y COLLATE {a} FROM t", "(* # (collate # {a}))"},
	{"SELECT x COLLATE {a} % 2 FROM t", "(% (collate # {a}) #)"},
	{"SELECT -x COLLATE {a} FROM t", "(u- (collate # {a}))"},
	{"SELECT ~x COLLATE {a} & 1 FROM t", "(& (u~ (collate # {a})) #)"},
	{"SELECT (x + y) COLLATE {a} FROM t", "(collate (paren (+ # #)) {a})"},
	{"SELECT CAST(x AS VARCHAR(10)) COLLATE {a} + 'b' FROM t", "(+ (collate # {a}) #)"},
	{"SELECT x COLLATE {a} AT TIME ZONE 'UTC' COLLATE {b} FROM t", "(atz (collate # {a}) (collate # {b}))"},
	{"SELECT 1 + d AT TIME ZONE 'UTC' COLLATE {a} FROM t", "(atz (+ # #) (collate # {a}))"},
	{"SELECT d AT TIME ZONE 'U' + 'TC' FROM t", "(+ (atz # #) #)"},
	{"SELECT 1 + d AT TIME ZONE 'UTC' + 2 FROM t", "(+ (atz (+ # #) #) #)"},
	{"SELECT d * 1 AT TIME ZONE 'UTC' FROM t", "(atz (* # #) #)"},
	{"SELECT d AT TIME ZONE 'UTC' * 2 FROM t", "(atz # (* # #))"},
	{"SELECT -d AT TIME ZONE 'UTC' FROM t", "(atz (u- #) #)"},
}

// exprShape renders an expression in the notation the SqlScriptDOM harness
// prints for a select item: operators, AT TIME ZONE as atz, parentheses, and
// collations, with every other operand as #.
func exprShape(n ast.Node) string {
	switch v := n.(type) {
	case *ast.BinaryExpr:
		ops := map[ast.BinaryOp]string{ast.BinOpAdd: "+", ast.BinOpSub: "-", ast.BinOpMul: "*", ast.BinOpDiv: "/",
			ast.BinOpMod: "%", ast.BinOpBitAnd: "&", ast.BinOpBitOr: "|", ast.BinOpBitXor: "^"}
		if op, ok := ops[v.Op]; ok {
			return "(" + op + " " + exprShape(v.Left) + " " + exprShape(v.Right) + ")"
		}
	case *ast.UnaryExpr:
		ops := map[ast.UnaryOp]string{ast.UnaryPlus: "u+", ast.UnaryMinus: "u-", ast.UnaryBitNot: "u~"}
		return "(" + ops[v.Op] + " " + exprShape(v.Operand) + ")"
	case *ast.AtTimeZoneExpr:
		return "(atz " + exprShape(v.Expr) + " " + exprShape(v.TimeZone) + ")"
	case *ast.ParenExpr:
		return "(paren " + exprShape(v.Expr) + ")"
	case *ast.CollateExpr:
		return "(collate " + exprShape(v.Expr) + " " + v.Collation + ")"
	}
	return "#"
}

// TestCollatePrecedence: omni builds the shape SqlScriptDOM builds, and SQL
// Server accepts each statement. omni applied COLLATE to the whole +/- chain
// before it, so it rejected 'a' COLLATE x + 'b' and collated 'a' + 'b' as a
// whole; it also stopped the chain at AT TIME ZONE, rejecting
// d AT TIME ZONE 'U' + 'TC'.
func TestCollatePrecedence(t *testing.T) {
	oracle := startParserOracle(t)
	h := getHarness(t)
	for _, tc := range collatePrecedenceCases {
		sql, want := collationNames.Replace(tc.sql), collationNames.Replace(tc.want)
		t.Run(tc.sql, func(t *testing.T) {
			if serverErr, err := oracle.parseError(sql); err != nil || serverErr != nil {
				t.Fatalf("SQL Server: %v %v, want accepted", serverErr, err)
			}
			out := h.Shape(t, sql)
			exprs, _ := out["exprs"].([]any)
			if out["errors"] != nil || len(exprs) == 0 {
				t.Fatalf("SqlScriptDOM: %v", out)
			}
			if exprs[0] != want {
				t.Errorf("SqlScriptDOM shape %v, want %s", exprs[0], want)
			}
			list, err := Parse(sql)
			if err != nil {
				t.Fatalf("omni: %v", err)
			}
			if got := exprShape(list.Items[0].(*ast.SelectStmt).TargetList.Items[0].(*ast.ResTarget).Val); got != want {
				t.Errorf("omni shape %s, want %s", got, want)
			}
			assertChildLocsWithinParent(t, sql, list)
		})
	}
}

// lineColumnOffset converts SqlScriptDOM's 1-based line and column to a byte
// offset into the ASCII sql.
func lineColumnOffset(sql string, line, column int) int {
	offset := 0
	for ; line > 1; line-- {
		offset += strings.IndexByte(sql[offset:], '\n') + 1
	}
	return offset + column - 1
}
