package deparse

import (
	"github.com/bytebase/omni/snowflake/ast"
)

// writer accumulates SQL text for departing AST nodes.
// It intentionally has no indentation — we produce single-line output
// optimised for round-trip correctness, not readability.
type writer struct {
	buf buffer
}

// buffer is the text a writer has produced. It appends like a
// strings.Builder, and it can also take back a space it already holds:
// writeExprNoLeadSpace writes an expression, whose first token may have put a
// space in front of itself (ensureSpace), and then removes that space.
//
// Cutting that byte out of a strings.Builder copied everything written so far,
// once per expression written that way, so deparse was quadratic in its
// output: 90,000 subscripts took 14 seconds. A taken-back byte is marked
// dead instead and left in place, and String drops the dead bytes in one
// pass. Every read (Len, lastByte, String) sees the text without them, the
// same text cutting produced, so the output is unchanged.
type buffer struct {
	b     []byte
	dead  []uint64 // bit i%64 of dead[i/64] is set when b[i] was taken back
	ndead int
}

// WriteString appends s.
func (b *buffer) WriteString(s string) {
	b.b = append(b.b, s...)
}

// WriteByte appends c. It never fails; the error is the io.ByteWriter shape.
func (b *buffer) WriteByte(c byte) error {
	b.b = append(b.b, c)
	return nil
}

// Len returns the length of the text, dead bytes excluded.
func (b *buffer) Len() int {
	return len(b.b) - b.ndead
}

// String returns the text, dead bytes excluded.
func (b *buffer) String() string {
	if b.ndead == 0 {
		return string(b.b)
	}
	out := make([]byte, 0, b.Len())
	for i, c := range b.b {
		if !b.isDead(i) {
			out = append(out, c)
		}
	}
	return string(out)
}

// lastByte returns the last byte of the text, or false when it is empty.
func (b *buffer) lastByte() (byte, bool) {
	for i := len(b.b) - 1; i >= 0; i-- {
		if !b.isDead(i) {
			return b.b[i], true
		}
	}
	return 0, false
}

// mark returns the position at which the next byte will be written, for a
// later dropLeadingSpace. Positions count dead bytes, so a mark stays valid
// while bytes before or after it are taken back.
func (b *buffer) mark() int {
	return len(b.b)
}

// dropLeadingSpace takes back the first byte of the text written since mark
// when it is a space: the byte that sits, once the dead are removed, where
// the text stood when mark was taken.
func (b *buffer) dropLeadingSpace(mark int) {
	i := mark
	for i < len(b.b) && b.isDead(i) {
		i++
	}
	if i == len(b.b) || b.b[i] != ' ' {
		return
	}
	for len(b.dead) <= i/64 {
		b.dead = append(b.dead, 0)
	}
	b.dead[i/64] |= 1 << (i % 64)
	b.ndead++
}

func (b *buffer) isDead(i int) bool {
	return i/64 < len(b.dead) && b.dead[i/64]&(1<<(i%64)) != 0
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
	last, ok := w.buf.lastByte()
	if !ok {
		return
	}
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
