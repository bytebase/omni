package deparse

import (
	"fmt"
	"reflect"

	ast "github.com/bytebase/omni/oceanbase/ast"
)

// maxAutoNameBytes is the longest auto-generated column name OceanBase keeps in a
// view, in bytes: a longer one becomes Name_exp_N. MySQL counts the 64 in
// characters. catalog's TestOracle_OceanBaseViewColumnNames holds it to the engine.
const maxAutoNameBytes = 64

// sourceAutoName names an unaliased select item from the text it was parsed
// from, as MySQL does (PTI_expr_with_alias::do_itemize, sql/parse_tree_items.cc):
// an item that named itself when the parser built it keeps that name, and any
// other item is named after its text from its first token to its last, as
// written. A name too long for a column becomes Name_exp_N. ok is false when
// the item's text cannot be read from source.
//
// Verified on MySQL 8.0 with SHOW CREATE VIEW: (with c as (select 1 as x)
// select x from   c union select 2) keeps its case and spacing, 0004 is
// "0004", true is "true", - 2 is "- 2", (a) is "a", (1) is "1", and a
// 65-character expression is Name_exp_N.
//
// The text keeps comments, and here OceanBase parts from MySQL: MySQL drops
// the markers of an executable comment from a name (a/*!50100 +2*/ is
// "a +2"), and OceanBase keeps them, so the name is the text between the
// item's first and last token exactly. OceanBase 4.3.5.6 names a/*!50100 +2*/
// "a/*!50100 +2", 2 * /*!100005 */ "2 * /*!100005", and 1 /*! + 1 */ + 1 in
// a statement it reads with its executable comments as plain ones (see
// oceanbase/parser/exec_fallback.go) "1 /*! + 1 */ + 1", measured on
// 2026-11-05 (catalog's viewExecCommentNameCases).
func sourceAutoName(source string, expr ast.ExprNode, position int) (string, bool) {
	if name, ok := selfName(source, expr); ok {
		if name == "" {
			return fmt.Sprintf("Name_exp_%d", position), true
		}
		return name, true
	}
	name, ok := nodeText(source, expr)
	if !ok {
		return "", false
	}
	if len(name) > maxAutoNameBytes {
		return fmt.Sprintf("Name_exp_%d", position), true
	}
	return name, true
}

// selfName returns the name an item gives itself as the parser builds it: a
// column reference its column, a string literal its value (an adjacent-literal
// run its first segment), NULL "NULL", and a number its literal text.
// Parentheses pass the item inside through, so (a) is named "a". ok is false
// for every other item, which MySQL names after its text.
func selfName(source string, expr ast.ExprNode) (string, bool) {
	switch n := expr.(type) {
	case *ast.ParenExpr:
		return selfName(source, n.Expr)
	case *ast.ColumnRef:
		if n.Star {
			return "", false
		}
		return n.Column, true
	case *ast.StringLit:
		if n.Concatenated {
			return n.FirstSegment, true
		}
		return n.Value, true
	case *ast.NullLit:
		return "NULL", true
	case *ast.IntLit, *ast.FloatLit:
		return nodeText(source, n)
	}
	return "", false
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
