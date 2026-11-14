package deparse

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// writeObjectStmt writes an object read by its documented shape: the verb and
// its modifiers, the object's keywords and name, then each clause in order.
func (w *writer) writeObjectStmt(n *ast.ObjectStmt) error {
	w.buf.WriteString(n.Verb)
	if n.OrReplace {
		w.buf.WriteString(" OR REPLACE")
	} else if n.OrAlter {
		w.buf.WriteString(" OR ALTER")
	}
	w.buf.WriteByte(' ')
	w.buf.WriteString(n.Kind)
	if n.IfNotExists {
		w.buf.WriteString(" IF NOT EXISTS")
	}
	if n.IfExists {
		w.buf.WriteString(" IF EXISTS")
	}
	if n.Name != nil {
		w.buf.WriteByte(' ')
		w.writeObjectNameNoSpace(n.Name)
	}
	for _, c := range n.Clauses {
		if ok, err := w.writeServiceClause(n, c); ok || err != nil {
			if err != nil {
				return err
			}
			continue
		}
		if err := w.writeObjectClause(c); err != nil {
			return err
		}
	}
	return nil
}

// writeObjectClause writes one clause: its keywords, then the names, keys,
// properties, tags and query it holds.
func (w *writer) writeObjectClause(c *ast.ObjectClause) error {
	if c.Keywords != "" {
		w.buf.WriteByte(' ')
		w.buf.WriteString(c.Keywords)
	}
	for i, name := range c.Names {
		if i == 0 {
			w.buf.WriteByte(' ')
		} else {
			w.buf.WriteString(", ")
		}
		w.writeObjectNameNoSpace(name)
	}
	if len(c.Keys) > 0 {
		w.buf.WriteByte(' ')
		w.buf.WriteString(strings.Join(c.Keys, ", "))
	}
	for _, o := range c.Options {
		w.buf.WriteByte(' ')
		if err := w.writeObjectOption(o); err != nil {
			return err
		}
	}
	if len(c.Tags) > 0 {
		switch c.Keywords {
		case "SET TAG":
			w.buf.WriteByte(' ')
			w.writeTagAssignments(c.Tags)
		case "SET":
			w.buf.WriteString(" TAG ")
			w.writeTagAssignments(c.Tags)
		default:
			w.buf.WriteString(" WITH TAG (")
			w.writeTagAssignments(c.Tags)
			w.buf.WriteByte(')')
		}
	}
	if c.Query != nil {
		w.buf.WriteByte(' ')
		return w.writeNode(c.Query)
	}
	return nil
}

// writeObjectOption writes a property as `NAME = value`, or its value alone
// when it has no name (an element of a list of groups).
func (w *writer) writeObjectOption(o *ast.CopyOption) error {
	if o.Name != "" {
		w.buf.WriteString(o.Name)
		if o.Bare {
			return nil
		}
		w.buf.WriteString(" = ")
	}
	switch {
	case o.Lit != nil:
		return w.writeCopyOptionLiteral(o.Lit)
	case o.Words != "":
		w.buf.WriteString(o.Words)
	case o.List != nil:
		w.buf.WriteByte('(')
		for i, lit := range o.List {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			if err := w.writeCopyOptionLiteral(lit); err != nil {
				return err
			}
		}
		w.buf.WriteByte(')')
	case o.Group != nil:
		w.buf.WriteByte('(')
		for i, e := range o.Group {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			if err := w.writeObjectOption(e); err != nil {
				return err
			}
		}
		w.buf.WriteByte(')')
	}
	return nil
}
