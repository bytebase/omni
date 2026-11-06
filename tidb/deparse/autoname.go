package deparse

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	ast "github.com/bytebase/omni/tidb/ast"
)

// maxAutoNameBytes is the longest auto-generated column name TiDB keeps in a
// view's column list, in bytes: a longer one becomes name_exp_N there.
const maxAutoNameBytes = 64

// sourceAutoName names an unaliased select item from the text it was parsed
// from, as TiDB does (the field text its parser records for a select field):
// an item that named itself when the parser built it keeps that name, and any
// other item is named after its text from its first token to its last, as the
// lexer copies it (see nameText). A name too long for a column becomes
// Name_exp_N. ok is false when the item's text cannot be read from source.
//
// Verified on TiDB v8.5 with SHOW CREATE VIEW: (select   1) keeps its case
// and spacing, 0004 is "0004", true is "TRUE", - 2 is "- 2", (a) is "(a)",
// and (1) is "1". A name longer than 64 bytes stays Name_exp_N, as before:
// TiDB lists it as name_exp_N in the view's column list.
//
// With capLength false a long name is kept whole, as TiDB keeps it for the
// field's alias in a view body; only the view's column list shortens it.
func sourceAutoName(source string, expr ast.ExprNode, position int, capLength bool) (string, bool) {
	if name, ok := selfName(source, expr); ok {
		if name == "" && capLength {
			return fmt.Sprintf("Name_exp_%d", position), true
		}
		// TiDB names an empty string literal "" and writes it with no
		// alias: SELECT '' is stored as SELECT _UTF8MB4'' (verified on
		// v8.5.5).
		return name, true
	}
	text, ok := itemText(source, expr)
	if !ok {
		return "", false
	}
	name := nameText(text)
	if capLength && len(name) > maxAutoNameBytes {
		return fmt.Sprintf("Name_exp_%d", position), true
	}
	return name, true
}

// selfName returns the name an item gives itself as the parser builds it: a
// column reference its column, a string literal its value, NULL "NULL", TRUE
// and FALSE "TRUE" and "FALSE", and a number its literal text. Parentheses
// pass a literal inside through, so (1) is named "1", but not a column
// reference: (a) is named "(a)". ok is false for every other item, which TiDB
// names after its text.
func selfName(source string, expr ast.ExprNode) (string, bool) {
	switch n := expr.(type) {
	case *ast.ParenExpr, *ast.UnaryExpr:
		// Parentheses and unary plus pass a literal inside through
		// (getInnerFromParenthesesAndUnaryPlus): (1) and +1 are named "1".
		inner := unwrapParenPlus(expr)
		if inner == expr {
			return "", false
		}
		if _, ok := inner.(*ast.ColumnRef); ok {
			return "", false
		}
		return selfName(source, inner)
	case *ast.ColumnRef:
		if n.Star {
			return "", false
		}
		return n.Column, true
	case *ast.StringLit:
		// The characters before the first graphic one are dropped
		// (planner/core buildProjectionFieldNameFromExpressions):
		// SELECT '  x' names its column x. Adjacent strings are named after
		// the first: SELECT '  x' 'y' names its column x (v8.5.5).
		return strings.TrimLeftFunc(stringLitName(n), func(r rune) bool {
			return !unicode.In(r, graphRanges...)
		}), true
	case *ast.NullLit:
		return "NULL", true
	case *ast.BoolLit:
		if n.Value {
			return "TRUE", true
		}
		return "FALSE", true
	case *ast.IntLit, *ast.FloatLit:
		return nodeText(source, n)
	}
	return "", false
}

// unwrapParenPlus strips the parentheses and unary pluses around an
// expression.
func unwrapParenPlus(expr ast.ExprNode) ast.ExprNode {
	for {
		switch n := expr.(type) {
		case *ast.ParenExpr:
			expr = n.Expr
			continue
		case *ast.UnaryExpr:
			if n.Op == ast.UnaryPlus {
				expr = n.Operand
				continue
			}
		}
		return expr
	}
}

// itemText returns the text TiDB names an item after: the text its Loc
// covers, extended over the comments that follow it up to the next token.
// Verified on TiDB v8.5: SELECT a*2 /* g */ , b names its first item
// "a*2 /* g " (nameText then drops the comment's "*/").
func itemText(source string, expr ast.ExprNode) (string, bool) {
	loc, ok := nodeLoc(expr)
	if !ok || loc.Start < 0 || loc.End > len(source) || loc.Start >= loc.End {
		return "", false
	}
	end := loc.End
	for i := end; i < len(source); {
		switch c := source[i]; {
		case isSpace(c):
			i++
		case c == '/' && i+1 < len(source) && source[i+1] == '*' && (i+2 >= len(source) || source[i+2] != '!'):
			i = commentEnd(source, i+2)
			end = i
		case c == '#' || isDashComment(source, i):
			for i < len(source) && source[i] != '\n' {
				i++
			}
			end = i
		default:
			return source[loc.Start:end], true
		}
	}
	return source[loc.Start:end], true
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' || c == '\v'
}

// nodeText returns the source text a node's Loc covers.
func nodeText(source string, node ast.Node) (string, bool) {
	loc, ok := nodeLoc(node)
	if !ok || loc.Start < 0 || loc.End > len(source) || loc.Start >= loc.End {
		return "", false
	}
	return source[loc.Start:loc.End], true
}

// nodeLoc reads a node's Loc field. The AST gives its nodes no common
// accessor; every node struct carries one.
func nodeLoc(node ast.Node) (ast.Loc, bool) {
	v := reflect.ValueOf(node)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return ast.Loc{}, false
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return ast.Loc{}, false
	}
	f := v.FieldByName("Loc")
	if !f.IsValid() {
		return ast.Loc{}, false
	}
	loc, ok := f.Interface().(ast.Loc)
	return loc, ok
}

// nameText returns an item's text as TiDB names it. TiDB reads the text once,
// left to right, and drops two things wherever they stand, in a string, a
// quoted identifier and a comment alike:
//
//   - "/*!", with the version after it when one follows: an optional M and
//     five or six digits;
//   - "*/".
//
// What a drop leaves is not read again, so "/**/!x" becomes "/*!x".
//
// Verified with views on TiDB v8.5.5 (TestNameText, and the oracle cases in
// catalog's view_names_test.go), which names
//
//	a /* c */ + 1     "a /* c  + 1"
//	a/*!80000 +2*/    "a +2"
//	a */*c*/3         "a *c3"
//	a + '*/'          "a + ''"
//
// and keeps a -- note <newline> + 1 as written.
func nameText(text string) string {
	if !strings.Contains(text, "/*!") && !strings.Contains(text, "*/") {
		return text
	}
	out := make([]byte, 0, len(text))
	for i := 0; i < len(text); {
		switch {
		case strings.HasPrefix(text[i:], "/*!"):
			i += 3
			j := i
			if j < len(text) && text[j] == 'M' {
				j++
			}
			k := j
			for k < len(text) && k-j < 6 && text[k] >= '0' && text[k] <= '9' {
				k++
			}
			if k-j >= 5 {
				i = k
			}
		case strings.HasPrefix(text[i:], "*/"):
			i += 2
		default:
			out = append(out, text[i])
			i++
		}
	}
	return string(out)
}

// isDashComment reports whether a "-- " comment starts at i: two dashes and a
// space or control character, as MySQL requires.
func isDashComment(text string, i int) bool {
	return i+2 < len(text) && text[i] == '-' && text[i+1] == '-' && text[i+2] <= ' '
}

// commentEnd returns the offset just past the "*/" that closes a block comment
// whose body starts at i, or len(text) when it is not closed.
func commentEnd(text string, i int) int {
	for ; i+1 < len(text); i++ {
		if text[i] == '*' && text[i+1] == '/' {
			return i + 2
		}
	}
	return len(text)
}

// graphRanges are the characters TiDB counts as graphic at the start of a
// string literal that names a column (parser/mysql RangeGraph).
var graphRanges = []*unicode.RangeTable{
	unicode.No, unicode.Mn, unicode.Me, unicode.Pc, unicode.Pd, unicode.Ps,
	unicode.Pe, unicode.Pi, unicode.Pf, unicode.Po, unicode.Sm, unicode.Sc,
	unicode.Sk, unicode.So, unicode.Lu, unicode.Lt, unicode.Nl, unicode.Ll,
	unicode.Lm, unicode.Lo, unicode.Mc, unicode.Nd,
}

// stringLitName is the part of a string literal TiDB names a column after: the
// first of adjacent strings ('a' 'b' is named a), or the whole value.
func stringLitName(n *ast.StringLit) string {
	if n.Concatenated {
		return n.FirstSegment
	}
	return n.Value
}
