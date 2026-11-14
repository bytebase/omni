package deparse

import "github.com/bytebase/omni/snowflake/ast"

// writeServiceClause writes the clauses of SERVICE and EXECUTE [ INFERENCE ]
// JOB SERVICE that writeObjectClause cannot, and reports whether c was one:
//
//	USING ( <key> => <value> [ , ... ] )   a specification template's arguments
//	FROM ( <query> )                       an inference job's input
func (w *writer) writeServiceClause(n *ast.ObjectStmt, c *ast.ObjectClause) (bool, error) {
	switch n.Kind {
	case "SERVICE", "JOB SERVICE", "INFERENCE JOB SERVICE":
	default:
		return false, nil
	}
	switch {
	case c.Keywords == "USING":
		w.buf.WriteString(" USING (")
		for i, o := range c.Options {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			w.buf.WriteString(o.Name)
			w.buf.WriteString(" => ")
			v := *o
			v.Name = ""
			if err := w.writeObjectOption(&v); err != nil {
				return true, err
			}
		}
		w.buf.WriteByte(')')
		return true, nil
	case c.Keywords == "FROM" && c.Query != nil:
		w.buf.WriteString(" FROM (")
		if err := w.writeNode(c.Query); err != nil {
			return true, err
		}
		w.buf.WriteByte(')')
		return true, nil
	}
	return false, nil
}
