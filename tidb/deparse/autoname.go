package deparse

import (
	"fmt"
	"reflect"
	"strings"

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
	if len(name) > maxAutoNameBytes {
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
	case *ast.ParenExpr:
		inner := n.Expr
		for {
			p, ok := inner.(*ast.ParenExpr)
			if !ok {
				break
			}
			inner = p.Expr
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

// nameText returns an item's text as TiDB names it: an executable comment
// /*!NNNNN ... */ loses its markers, any other block comment loses its
// closing "*/", and everything else stays as written.
//
// Verified on TiDB v8.5: a /* c */ + 1 is named "a /* c  + 1", a -- note
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
			body := text[i:j]
			if strings.HasSuffix(body, "*/") && len(body) >= 4 {
				body = body[:len(body)-2]
			}
			out = append(out, body...)
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
