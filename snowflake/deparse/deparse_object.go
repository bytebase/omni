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
		if err := w.writeObjectClause(n.Kind, c); err != nil {
			return err
		}
	}
	return nil
}

// writeObjectClause writes one clause of an object of kind: its keywords, then
// the names, keys, properties, tags and query it holds.
func (w *writer) writeObjectClause(kind string, c *ast.ObjectClause) error {
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
	if tt := c.TimeTravel; tt != nil {
		// A point in time: { AT | BEFORE } ( <anchor> => <expr> ), as a
		// POSTGRES INSTANCE's FORK writes it. Keywords already wrote AT or
		// BEFORE.
		w.buf.WriteString(" (")
		w.buf.WriteString(timeTravelAnchorName(tt.Anchor))
		w.buf.WriteString(" => ")
		if err := w.writeExpr(tt.Expr); err != nil {
			return err
		}
		w.buf.WriteByte(')')
		return nil
	}
	commas := objectOptionCommas(kind, c.Keywords)
	for i, o := range c.Options {
		// A clause's bare entries and nameless options are lists
		// (listsWith); some pages separate a clause's properties with
		// commas (objectOptionCommas).
		if i > 0 && (commas || listsWith(c.Options[i-1], o)) {
			w.buf.WriteByte(',')
		}
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

// listsWith reports whether option o continues a comma list after prev: a
// clause's bare entries (VECTOR INDEXES a, b) or its nameless options
// (REGION GROUPS 'a', 'b').
func listsWith(prev, o *ast.CopyOption) bool {
	return (o.Bare && prev.Bare) || (o.Name == "" && prev.Name == "")
}

// objectOptionCommas reports whether a clause of an object of kind separates
// its properties with commas, as the object's page writes them.
func objectOptionCommas(kind, keywords string) bool {
	switch keywords {
	case "SET CONTACT":
		return true
	case "SET":
		return kind == "TASK GRAPH"
	}
	return false
}

// writeObjectOption writes a property as `NAME = value`, or its value alone
// when it has no name (an element of a list of groups). A bare entry is its
// name alone, followed by its Group in parentheses when it has one
// (VECTOR INDEXES' <column> ( model = '<model>' )).
func (w *writer) writeObjectOption(o *ast.CopyOption) error {
	if o.Name == "EXECUTE AS USER" && !o.Bare {
		// TASK GRAPH's TASK_DEFAULTS ( EXECUTE AS USER <user_name> ): no '='.
		w.buf.WriteString("EXECUTE AS USER ")
		w.buf.WriteString(o.Words)
		return nil
	}
	if o.Name != "" {
		w.buf.WriteString(o.Name)
		if o.Bare {
			if o.Group != nil {
				w.buf.WriteString(" (")
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

// timeTravelAnchorName is the keyword a time-travel anchor is written with.
func timeTravelAnchorName(a ast.TimeTravelAnchor) string {
	switch a {
	case ast.TimeTravelOffset:
		return "OFFSET"
	case ast.TimeTravelStatement:
		return "STATEMENT"
	case ast.TimeTravelStream:
		return "STREAM"
	}
	return "TIMESTAMP"
}
