package deparse

import (
	"fmt"

	"github.com/bytebase/omni/snowflake/ast"
)

// writeInsertStmt writes a single-table INSERT statement.
func (w *writer) writeInsertStmt(n *ast.InsertStmt) error {
	w.buf.WriteString("INSERT")
	if n.Overwrite {
		w.buf.WriteString(" OVERWRITE")
	}
	w.buf.WriteString(" INTO ")
	w.writeObjectNameNoSpace(n.Target)
	if len(n.Columns) > 0 {
		w.buf.WriteString(" (")
		for i, col := range n.Columns {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			w.buf.WriteString(col.String())
		}
		w.buf.WriteByte(')')
	}
	if n.Select != nil {
		w.buf.WriteByte(' ')
		return w.writeQueryExpr(n.Select)
	}
	// VALUES form
	w.buf.WriteString(" VALUES ")
	for i, row := range n.Values {
		if i > 0 {
			w.buf.WriteString(", ")
		}
		w.buf.WriteByte('(')
		for j, val := range row {
			if j > 0 {
				w.buf.WriteString(", ")
			}
			if err := writeExprNoLeadSpace(w, val); err != nil {
				return err
			}
		}
		w.buf.WriteByte(')')
	}
	return nil
}

// writeInsertMultiStmt writes an INSERT ALL / INSERT FIRST statement.
func (w *writer) writeInsertMultiStmt(n *ast.InsertMultiStmt) error {
	w.buf.WriteString("INSERT")
	if n.Overwrite {
		w.buf.WriteString(" OVERWRITE")
	}
	if n.First {
		w.buf.WriteString(" FIRST")
	} else {
		w.buf.WriteString(" ALL")
	}

	// For INSERT FIRST with WHEN branches, branches without a When condition
	// that follow conditional branches are ELSE branches.
	hasConditional := false
	for _, b := range n.Branches {
		if b.When != nil {
			hasConditional = true
			break
		}
	}

	inElse := false
	var prevWhen ast.Node
	for i, branch := range n.Branches {
		if branch.When != nil {
			// New WHEN guard — only write "WHEN ... THEN" once per condition
			// (multiple INTO branches can share the same When object).
			if i == 0 || branch.When != prevWhen {
				w.buf.WriteString(" WHEN ")
				if err := writeExprNoLeadSpace(w, branch.When); err != nil {
					return err
				}
				w.buf.WriteString(" THEN")
			}
			prevWhen = branch.When
			inElse = false
		} else if hasConditional && n.First && !inElse {
			// First nil-When branch after conditional branches → ELSE
			w.buf.WriteString(" ELSE")
			inElse = true
		}
		w.buf.WriteString(" INTO ")
		w.writeObjectNameNoSpace(branch.Target)
		if len(branch.Columns) > 0 {
			w.buf.WriteString(" (")
			for i, col := range branch.Columns {
				if i > 0 {
					w.buf.WriteString(", ")
				}
				w.buf.WriteString(col.String())
			}
			w.buf.WriteByte(')')
		}
		if len(branch.Values) > 0 {
			w.buf.WriteString(" VALUES (")
			for i, val := range branch.Values {
				if i > 0 {
					w.buf.WriteString(", ")
				}
				if err := writeExprNoLeadSpace(w, val); err != nil {
					return err
				}
			}
			w.buf.WriteByte(')')
		}
	}
	w.buf.WriteByte(' ')
	return w.writeQueryExpr(n.Select)
}

// writeUpdateStmt writes an UPDATE statement.
func (w *writer) writeUpdateStmt(n *ast.UpdateStmt) error {
	w.buf.WriteString("UPDATE ")
	w.writeObjectNameNoSpace(n.Target)
	if !n.Alias.IsEmpty() {
		w.buf.WriteString(" AS ")
		w.buf.WriteString(n.Alias.String())
	}
	w.buf.WriteString(" SET ")
	for i, set := range n.Sets {
		if i > 0 {
			w.buf.WriteString(", ")
		}
		w.writeObjectNameNoSpace(set.Column)
		w.buf.WriteString(" = ")
		if err := writeExprNoLeadSpace(w, set.Value); err != nil {
			return err
		}
	}
	if len(n.From) > 0 {
		w.buf.WriteString(" FROM ")
		for i, src := range n.From {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			if err := w.writeFromSource(src); err != nil {
				return err
			}
		}
	}
	if n.Where != nil {
		w.buf.WriteString(" WHERE ")
		if err := writeExprNoLeadSpace(w, n.Where); err != nil {
			return err
		}
	}
	return nil
}

// writeDeleteStmt writes a DELETE statement.
func (w *writer) writeDeleteStmt(n *ast.DeleteStmt) error {
	w.buf.WriteString("DELETE FROM ")
	w.writeObjectNameNoSpace(n.Target)
	if !n.Alias.IsEmpty() {
		w.buf.WriteString(" AS ")
		w.buf.WriteString(n.Alias.String())
	}
	if len(n.Using) > 0 {
		w.buf.WriteString(" USING ")
		for i, src := range n.Using {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			if err := w.writeFromSource(src); err != nil {
				return err
			}
		}
	}
	if n.Where != nil {
		w.buf.WriteString(" WHERE ")
		if err := writeExprNoLeadSpace(w, n.Where); err != nil {
			return err
		}
	}
	return nil
}

// writeMergeStmt writes a MERGE statement.
func (w *writer) writeMergeStmt(n *ast.MergeStmt) error {
	w.buf.WriteString("MERGE INTO ")
	w.writeObjectNameNoSpace(n.Target)
	if !n.TargetAlias.IsEmpty() {
		w.buf.WriteString(" AS ")
		w.buf.WriteString(n.TargetAlias.String())
	}
	w.buf.WriteString(" USING ")
	switch src := n.Source.(type) {
	case *ast.TableRef:
		if err := w.writeTableRef(src); err != nil {
			return err
		}
	case *ast.SelectStmt:
		w.buf.WriteByte('(')
		if err := w.writeSelectStmt(src); err != nil {
			return err
		}
		w.buf.WriteByte(')')
	case *ast.SubqueryExpr:
		w.buf.WriteByte('(')
		if err := w.writeQueryExpr(src.Query); err != nil {
			return err
		}
		w.buf.WriteByte(')')
	default:
		return fmt.Errorf("deparse: unsupported MERGE source type %T", n.Source)
	}
	if !n.SourceAlias.IsEmpty() {
		w.buf.WriteString(" AS ")
		w.buf.WriteString(n.SourceAlias.String())
	}
	w.buf.WriteString(" ON ")
	if err := writeExprNoLeadSpace(w, n.On); err != nil {
		return err
	}
	for _, when := range n.Whens {
		if err := w.writeMergeWhen(when); err != nil {
			return err
		}
	}
	return nil
}

func (w *writer) writeMergeWhen(mw *ast.MergeWhen) error {
	if mw.Matched {
		w.buf.WriteString(" WHEN MATCHED")
	} else if mw.BySource {
		w.buf.WriteString(" WHEN NOT MATCHED BY SOURCE")
	} else {
		w.buf.WriteString(" WHEN NOT MATCHED")
	}
	if mw.AndCond != nil {
		w.buf.WriteString(" AND ")
		if err := writeExprNoLeadSpace(w, mw.AndCond); err != nil {
			return err
		}
	}
	w.buf.WriteString(" THEN")
	switch mw.Action {
	case ast.MergeActionUpdate:
		w.buf.WriteString(" UPDATE SET ")
		for i, set := range mw.Sets {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			w.writeObjectNameNoSpace(set.Column)
			w.buf.WriteString(" = ")
			if err := writeExprNoLeadSpace(w, set.Value); err != nil {
				return err
			}
		}
	case ast.MergeActionDelete:
		w.buf.WriteString(" DELETE")
	case ast.MergeActionInsert:
		w.buf.WriteString(" INSERT")
		if mw.InsertDefault {
			w.buf.WriteString(" VALUES DEFAULT")
		} else {
			if len(mw.InsertCols) > 0 {
				w.buf.WriteString(" (")
				for i, col := range mw.InsertCols {
					if i > 0 {
						w.buf.WriteString(", ")
					}
					w.buf.WriteString(col.String())
				}
				w.buf.WriteByte(')')
			}
			if len(mw.InsertVals) > 0 {
				w.buf.WriteString(" VALUES (")
				for i, val := range mw.InsertVals {
					if i > 0 {
						w.buf.WriteString(", ")
					}
					if err := writeExprNoLeadSpace(w, val); err != nil {
						return err
					}
				}
				w.buf.WriteByte(')')
			}
		}
	}
	return nil
}
