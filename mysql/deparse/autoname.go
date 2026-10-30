package deparse

import (
	"fmt"
	"reflect"
	"strings"
	"unicode/utf8"

	ast "github.com/bytebase/omni/mysql/ast"
)

// maxAutoNameChars is the longest auto-generated column name MySQL keeps in a
// view (NAME_CHAR_LEN, counted in characters): a longer one becomes
// Name_exp_N.
const maxAutoNameChars = 64

// sourceAutoName names an unaliased select item from the text it was parsed
// from, as MySQL does (PTI_expr_with_alias::do_itemize, sql/parse_tree_items.cc):
// an item that named itself when the parser built it keeps that name, and any
// other item is named after its text from its first token to its last, as the
// lexer copies it (see nameText). A name too long for a column becomes
// Name_exp_N. ok is false when the item's text cannot be read from source.
//
// Verified on MySQL 8.0 with SHOW CREATE VIEW: (with c as (select 1 as x)
// select x from   c union select 2) keeps its case and spacing, 0004 is
// "0004", true is "true", - 2 is "- 2", (a) is "a", (1) is "1", and a
// 65-character expression is Name_exp_N.
func sourceAutoName(source string, expr ast.ExprNode, position int) (string, bool) {
	if name, ok := selfName(source, expr); ok {
		if name == "" {
			return fmt.Sprintf("Name_exp_%d", position), true
		}
		return name, true
	}
	text, ok := nodeText(source, expr)
	if !ok {
		return "", false
	}
	name := nameText(text)
	if utf8.RuneCountInString(name) > maxAutoNameChars {
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

// nameText returns an item's text as MySQL's lexer copies it into the
// preprocessed query an item name is taken from: an executable comment
// /*!NNNNN ... */ loses its markers, and everything else, other comments
// included, stays as written.
//
// Verified on MySQL 8.0: a /* c */ + 1 is named "a /* c */ + 1", a -- note
// <newline> + 1 keeps its comment, and a/*!80000 +2*/ is named "a +2".
func nameText(text string) string {
	if !strings.Contains(text, "/*!") && !strings.Contains(text, "*/") {
		return text
	}
	out := make([]byte, 0, len(text))
	open := 0 // executable comments opened in text and not yet closed
	for i := 0; i < len(text); {
		c := text[i]
		switch {
		case c == '\'' || c == '"' || c == '`':
			j := quotedEnd(text, i)
			out = append(out, text[i:j]...)
			i = j
		case c == '/' && i+2 < len(text) && text[i+1] == '*' && text[i+2] == '!':
			i += 3
			for i < len(text) && text[i] >= '0' && text[i] <= '9' {
				i++
			}
			open++
		case c == '/' && i+1 < len(text) && text[i+1] == '*':
			j := commentEnd(text, i+2)
			out = append(out, text[i:j]...)
			i = j
		case c == '#' || isDashComment(text, i):
			j := i
			for j < len(text) && text[j] != '\n' {
				j++
			}
			out = append(out, text[i:j]...)
			i = j
		case c == '*' && i+1 < len(text) && text[i+1] == '/' &&
			(open > 0 || i+2 == len(text) || text[i+2] != '*'):
			// Outside a quote and a comment, "*/" closes an executable comment,
			// opened in the text or before it; "*/*" with none open is a
			// multiplication followed by a comment.
			if open > 0 {
				open--
			}
			i += 2
		default:
			out = append(out, c)
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

// quotedEnd returns the offset just past the quoted string or identifier
// starting at i. A doubled quote stands for itself; inside a string, a
// backslash escapes the next byte.
func quotedEnd(text string, i int) int {
	q := text[i]
	for j := i + 1; j < len(text); j++ {
		switch {
		case text[j] == '\\' && q != '`':
			j++
		case text[j] == q:
			if j+1 < len(text) && text[j+1] == q {
				j++
				continue
			}
			return j + 1
		}
	}
	return len(text)
}
