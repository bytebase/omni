package deparse

import (
	"fmt"
	"reflect"
	"strings"
	"unicode/utf8"

	ast "github.com/bytebase/omni/mariadb/ast"
)

// maxAutoNameChars is the longest auto-generated column name MariaDB keeps in
// a view (NAME_CHAR_LEN, counted in characters): a longer one becomes
// Name_exp_N.
const maxAutoNameChars = 64

// sourceAutoName names an unaliased select item from the text it was parsed
// from, as MariaDB does (the select_item rule of sql/sql_yacc.yy):
// an item that named itself when the parser built it keeps that name, and any
// other item is named after its text from its first token to its last, as the
// lexer copies it (see nameText). A name too long for a column becomes
// Name_exp_N. ok is false when the item's text cannot be read from source.
//
// Verified on MariaDB 11.8 with SHOW CREATE VIEW: (with c as (select 1 as x)
// select x from   c union select 2) keeps its case and spacing, 0004 is
// "0004", true is "TRUE", - 2 is "- 2", (a) is "a", (1) is "1", and a
// 65-character expression is Name_exp_N.
func sourceAutoName(source string, expr ast.ExprNode, position int) (string, bool) {
	if name, ok := selfName(source, expr); ok {
		if name == "" {
			return fmt.Sprintf("Name_exp_%d", position), true
		}
		return name, true
	}
	text, ok := itemText(source, expr)
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
// run its whole value), NULL "NULL", TRUE and FALSE "TRUE" and "FALSE", and a
// number its literal text.
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
		return n.Value, true
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

// itemText returns the text MariaDB names an item after: the text its Loc
// covers, extended back over the comments between it and the comma before
// it. Verified on MariaDB 11.8: in SELECT a+1, /* d */ b+1 the second item is
// named "/* d */ b+1", while SELECT /* c */ a+1 names its first item "a+1".
func itemText(source string, expr ast.ExprNode) (string, bool) {
	loc, ok := nodeLoc(expr)
	if !ok || loc.Start < 0 || loc.End > len(source) || loc.Start >= loc.End {
		return "", false
	}
	start := loc.Start
	i := start
	for {
		for i > 0 && isSpace(source[i-1]) {
			i--
		}
		if i >= 2 && source[i-2] == '*' && source[i-1] == '/' {
			open := strings.LastIndex(source[:i-2], "/*")
			if open < 0 {
				break
			}
			i = open
			continue
		}
		break
	}
	if i > 0 && source[i-1] == ',' {
		for i < start && isSpace(source[i]) {
			i++
		}
		start = i
	}
	return source[start:loc.End], true
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

// nameText returns an item's text as MariaDB's lexer copies it into the
// preprocessed query an item name is taken from: an executable comment
// /*!NNNNN ... */ loses its markers, and everything else, other comments
// included, stays as written.
//
// Verified on MariaDB 11.8: a /* c */ + 1 is named "a /* c */ + 1", a -- note
// <newline> + 1 keeps its comment, and a /*!+3*/ is named "a +3".
func nameText(text string) string {
	if !strings.Contains(text, "/*!") && !strings.Contains(text, "/*M!") && !strings.Contains(text, "*/") {
		return text
	}
	out := make([]byte, 0, len(text))
	open := false // an executable comment opened in text and not yet closed
	for i := 0; i < len(text); {
		c := text[i]
		switch {
		case c == '\'' || c == '"' || c == '`':
			j := quotedEnd(text, i)
			out = append(out, text[i:j]...)
			i = j
		case skippedComment(text, i):
			// A comment MariaDB skips unread drops out of the name with its body:
			// it names "1 /*!80023 +3*/ + 1" as "1  + 1".
			if j := skippedCommentEnd(text, i); j >= 0 {
				i = j
			} else {
				i = len(text)
			}
		case execOpenerAt(text, i) > 0:
			// Executable comments do not nest: an opener inside one is dropped,
			// and the first "*/" closes both.
			i = execOpenerAt(text, i)
			open = true
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
			(open || i+2 == len(text) || text[i+2] != '*'):
			// Outside a quote and a comment, "*/" closes an executable comment,
			// opened in the text or before it; "*/*" with none open is a
			// multiplication followed by a comment.
			open = false
			i += 2
		default:
			out = append(out, c)
			i++
		}
	}
	return string(out)
}

// skippedComment, skippedCommentEnd and execOpenerAt read comments as
// mariadb/parser's skippedComment, skippedCommentEnd and execOpener do. A
// comment MariaDB skips unread is "/*!" versioned from 50700 to 99999 (five
// digits, or six when a sixth follows), or "/*!" or "/*M!" versioned 999999,
// and ends where its nesting, counted over raw bytes, closes. Any other "/*!"
// or "/*M!" is executable: MariaDB 11.8.8 names "1 /*M!+3*/" as "1 +3".
func skippedComment(s string, i int) bool {
	if i+2 >= len(s) || s[i] != '/' || s[i+1] != '*' {
		return false
	}
	j := i + 3
	mariadb := false
	switch {
	case s[i+2] == '!':
	case s[i+2] == 'M' && i+3 < len(s) && s[i+3] == '!':
		j = i + 4
		mariadb = true
	default:
		return false
	}
	v, n := 0, 0
	for j < len(s) && n < 6 && s[j] >= '0' && s[j] <= '9' {
		v = v*10 + int(s[j]-'0')
		j++
		n++
	}
	if n == 6 && v == 999999 {
		return true
	}
	return !mariadb && n >= 5 && v >= 50700 && v <= 99999
}

func skippedCommentEnd(s string, i int) int {
	depth := 0
	for i+1 < len(s) {
		switch {
		case s[i] == '/' && s[i+1] == '*':
			depth++
			i += 2
		case s[i] == '*' && s[i+1] == '/':
			depth--
			i += 2
			if depth == 0 {
				return i
			}
		default:
			i++
		}
	}
	return -1
}

// execOpenerAt returns the offset after the opener of the executable comment
// starting at text[i], or 0 when none starts there.
func execOpenerAt(text string, i int) int {
	if i+2 >= len(text) || text[i] != '/' || text[i+1] != '*' {
		return 0
	}
	j := i + 3
	switch {
	case text[i+2] == '!':
	case text[i+2] == 'M' && i+3 < len(text) && text[i+3] == '!':
		j = i + 4
	default:
		return 0
	}
	if skippedComment(text, i) {
		return 0
	}
	for j < len(text) && text[j] >= '0' && text[j] <= '9' {
		j++
	}
	return j
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
