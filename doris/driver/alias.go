package driver

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bytebase/omni/doris/ast"
	dorisparser "github.com/bytebase/omni/doris/parser"
)

// aliasFunctionSource returns an alias function's definition as source that
// the server renders back to the statement it states, so a dump of a restored
// alias function states the same statement again. Any other function's
// definition is returned as it is.
//
// SHOW CREATE FUNCTION states an alias function's body as the server rendered
// it, not as it was written, and replaying that rendering renders it again:
//
//   - every arithmetic operation comes back in parentheses of its own, and a
//     parenthesized one keeps the parentheses written around it, so x - y
//     comes back as (`x` - `y`) and replays as ((`x` - `y`));
//   - unary minus comes back as a product, -x as (-1 * `x`), which replays as
//     ((-1 * (1 * `x`))), one more product each time;
//   - a string literal comes back with its quotes doubled and its backslashes
//     as they are, so 'a\\b' comes back as 'a\b' and replays as 'ab'.
//
// The other renderings, comparisons and logical operators in parentheses of
// their own among them, replay to themselves. TestDumpAliasFunctionRestores
// holds the inverse to dyrnq/doris:3.1.4.
func aliasFunctionSource(stmt string) (string, error) {
	if !isAliasFunction(stmt) {
		return stmt, nil
	}
	escaped, err := escapeRenderedStrings(stmt)
	if err != nil {
		return "", err
	}
	file, err := dorisparser.Parse(escaped)
	if err != nil {
		return "", err
	}
	if len(file.Stmts) != 1 {
		return "", fmt.Errorf("%d statements, want one", len(file.Stmts))
	}
	create, ok := file.Stmts[0].(*ast.CreateFunctionStmt)
	if !ok || !create.Alias || create.Body == nil {
		return "", fmt.Errorf("%T is not an alias function's definition", file.Stmts[0])
	}
	var edits []textEdit
	var walkErr error
	ast.Inspect(create.Body, func(n ast.Node) bool {
		paren, ok := n.(*ast.ParenExpr)
		if !ok || walkErr != nil {
			return walkErr == nil
		}
		binary, ok := paren.Expr.(*ast.BinaryExpr)
		if !ok || !renderedInParens(binary.Op) {
			return true
		}
		if paren.Loc.Start < 0 || paren.Loc.End > len(escaped) || paren.Loc.End-paren.Loc.Start < 2 {
			walkErr = fmt.Errorf("no position for %q's parentheses", escaped)
			return false
		}
		closing := textEdit{start: paren.Loc.End - 1, end: paren.Loc.End}
		if negated := renderedNegation(binary); negated != nil {
			// (-1 * e) is -e: the minus replaces everything before the operand.
			at := ast.NodeLoc(negated).Start
			if at <= paren.Loc.Start || at > paren.Loc.End {
				walkErr = fmt.Errorf("no position for the operand of %q", escaped[paren.Loc.Start:paren.Loc.End])
				return false
			}
			edits = append(edits, textEdit{start: paren.Loc.Start, end: at, text: "-"}, closing)
			return true
		}
		edits = append(edits, textEdit{start: paren.Loc.Start, end: paren.Loc.Start + 1}, closing)
		return true
	})
	if walkErr != nil {
		return "", walkErr
	}
	return applyTextEdits(escaped, edits), nil
}

// isAliasFunction reports whether stmt is the server's statement for an alias
// function: CREATE [GLOBAL] ALIAS FUNCTION.
func isAliasFunction(stmt string) bool {
	fields := strings.Fields(stmt)
	if len(fields) > 1 && strings.EqualFold(fields[1], "GLOBAL") {
		fields = append(fields[:1:1], fields[2:]...)
	}
	return len(fields) > 2 && strings.EqualFold(fields[0], "CREATE") && strings.EqualFold(fields[1], "ALIAS") && strings.EqualFold(fields[2], "FUNCTION")
}

// renderedInParens reports whether the server renders an operation of op in
// parentheses of its own that replaying it keeps: the arithmetic and bitwise
// operators. A comparison is rendered in parentheses too, but replaying one
// drops the parentheses written around it, so its rendering is its own source.
func renderedInParens(op ast.BinaryOp) bool {
	switch op {
	case ast.BinAdd, ast.BinSub, ast.BinMul, ast.BinDiv, ast.BinMod, ast.BinIntDiv,
		ast.BinBitOr, ast.BinBitAnd, ast.BinBitXor, ast.BinShiftLeft, ast.BinShiftRight:
		return true
	}
	return false
}

// renderedNegation returns e when binary is the server's rendering of -e,
// -1 * e; a product a source writes with -1 renders with the operand as a
// product of its own ((-1 * (1 * e))), so the form is only a negation's.
func renderedNegation(binary *ast.BinaryExpr) ast.Node {
	if binary.Op != ast.BinMul {
		return nil
	}
	minus, ok := binary.Left.(*ast.UnaryExpr)
	if !ok || minus.Op != ast.UnaryMinus {
		return nil
	}
	if one, ok := minus.Expr.(*ast.Literal); !ok || one.Value != "1" {
		return nil
	}
	return binary.Right
}

// escapeRenderedStrings returns stmt with each string literal's backslashes
// escaped. The server renders a string literal as its value in single quotes
// with each quote doubled and nothing else escaped, so a backslash in the
// value would read as an escape when the statement replays.
func escapeRenderedStrings(stmt string) (string, error) {
	var b strings.Builder
	for i := 0; i < len(stmt); {
		switch c := stmt[i]; c {
		case '`':
			end := closingQuote(stmt, i, '`')
			if end < 0 {
				return "", fmt.Errorf("unterminated identifier at byte %d of %q", i, stmt)
			}
			b.WriteString(stmt[i:end])
			i = end
		case '\'':
			end := closingQuote(stmt, i, '\'')
			if end < 0 {
				return "", fmt.Errorf("unterminated string at byte %d of %q", i, stmt)
			}
			b.WriteString(strings.ReplaceAll(stmt[i:end], `\`, `\\`))
			i = end
		default:
			b.WriteByte(c)
			i++
		}
	}
	return b.String(), nil
}

// closingQuote returns the offset just past the quoted text that opens at
// start, where a doubled quote stands for one, or -1 when it does not close.
func closingQuote(s string, start int, quote byte) int {
	for i := start + 1; i < len(s); i++ {
		if s[i] != quote {
			continue
		}
		if i+1 < len(s) && s[i+1] == quote {
			i++
			continue
		}
		return i + 1
	}
	return -1
}

// textEdit replaces s[start:end] with text.
type textEdit struct {
	start, end int
	text       string
}

// applyTextEdits applies non-overlapping edits to s. A minus an edit writes is
// kept apart from a minus that follows it, which would otherwise start a
// comment.
func applyTextEdits(s string, edits []textEdit) string {
	sort.SliceStable(edits, func(i, j int) bool { return edits[i].start < edits[j].start })
	var b strings.Builder
	pending := false
	write := func(text string) {
		if text == "" {
			return
		}
		if pending && text[0] == '-' {
			b.WriteByte(' ')
		}
		b.WriteString(text)
		pending = strings.HasSuffix(text, "-")
	}
	at := 0
	for _, e := range edits {
		write(s[at:e.start])
		write(e.text)
		at = e.end
	}
	write(s[at:])
	return b.String()
}
