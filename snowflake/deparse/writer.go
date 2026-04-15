package deparse

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// writer accumulates SQL text for departing AST nodes.
// It intentionally has no indentation — we produce single-line output
// optimised for round-trip correctness, not readability.
type writer struct {
	buf strings.Builder
}

// String returns the accumulated SQL string.
func (w *writer) String() string {
	return w.buf.String()
}

// writeString appends s verbatim.
func (w *writer) writeString(s string) {
	w.buf.WriteString(s)
}

// writeByte appends one byte.
func (w *writer) writeByte(b byte) {
	w.buf.WriteByte(b)
}

// writeKeyword appends a SQL keyword in UPPER CASE, preceded by a space if
// the buffer is non-empty and doesn't already end with a space or '('.
func (w *writer) writeKeyword(kw string) {
	w.ensureSpace()
	w.buf.WriteString(kw)
}

// writeIdent writes an identifier, re-quoting it if Quoted is true.
func (w *writer) writeIdent(id ast.Ident) {
	w.ensureSpace()
	w.buf.WriteString(id.String())
}

// writeIdentNoSpace writes an identifier without prepending a space.
// Used when preceding context (e.g. '(', '.', ',') already handles spacing.
func (w *writer) writeIdentNoSpace(id ast.Ident) {
	w.buf.WriteString(id.String())
}

// writeObjectName writes a possibly-qualified object name.
func (w *writer) writeObjectName(n *ast.ObjectName) {
	w.ensureSpace()
	w.writeObjectNameNoSpace(n)
}

// writeObjectNameNoSpace writes an object name without a preceding space.
func (w *writer) writeObjectNameNoSpace(n *ast.ObjectName) {
	w.buf.WriteString(n.String())
}

// ensureSpace inserts a single space if needed between adjacent tokens.
// It does NOT add a space if the buffer is empty, or the last char is one
// of ' ', '(', '.', or ':'.
func (w *writer) ensureSpace() {
	if w.buf.Len() == 0 {
		return
	}
	last := w.buf.String()[w.buf.Len()-1]
	switch last {
	case ' ', '(', '.', ':':
		// no space needed
	default:
		w.buf.WriteByte(' ')
	}
}

// writeCommaList writes a parenthesised, comma-separated list of items.
// fn is called once per item.
func (w *writer) writeCommaList(n int, fn func(i int)) {
	w.buf.WriteByte('(')
	for i := 0; i < n; i++ {
		if i > 0 {
			w.buf.WriteString(", ")
		}
		fn(i)
	}
	w.buf.WriteByte(')')
}

// writeSep writes items separated by sep (e.g. ", " or " ").
// fn is called once per item.
func (w *writer) writeSep(n int, sep string, fn func(i int)) {
	for i := 0; i < n; i++ {
		if i > 0 {
			w.buf.WriteString(sep)
		}
		fn(i)
	}
}
