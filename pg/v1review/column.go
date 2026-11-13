package v1review

import "github.com/bytebase/omni/pg/ast"

// statementNode returns a top-level statement without its RawStmt wrapper.
func statementNode(stmt ast.Node) ast.Node {
	if raw, ok := stmt.(*ast.RawStmt); ok {
		return raw.Stmt
	}
	return stmt
}

// tableColumns returns the column definitions a CREATE TABLE writes, in
// order.
func tableColumns(n *ast.CreateStmt) []*ast.ColumnDef {
	if n.TableElts == nil {
		return nil
	}
	var columns []*ast.ColumnDef
	for _, item := range n.TableElts.Items {
		if col, ok := item.(*ast.ColumnDef); ok {
			columns = append(columns, col)
		}
	}
	return columns
}

// alterCommands returns an ALTER statement's actions, in order.
func alterCommands(n *ast.AlterTableStmt) []*ast.AlterTableCmd {
	if n.Cmds == nil {
		return nil
	}
	var cmds []*ast.AlterTableCmd
	for _, item := range n.Cmds.Items {
		if cmd, ok := item.(*ast.AlterTableCmd); ok {
			cmds = append(cmds, cmd)
		}
	}
	return cmds
}

// lastTypeName returns the last part of a type's name as the parser stores
// it, without schema or modifiers: "int4" for INT, "bpchar" for CHAR(10).
func lastTypeName(tn *ast.TypeName) string {
	if tn == nil || tn.Names == nil {
		return ""
	}
	name := ""
	for _, item := range tn.Names.Items {
		if s, ok := item.(*ast.String); ok {
			name = s.Str
		}
	}
	return name
}
