package ast

// This file holds AST node types for stored procedure DDL/DML (T8.2):
// CREATE PROCEDURE, CALL, and DROP PROCEDURE.

// ProcedureParam represents one parameter in a CREATE PROCEDURE parameter list.
//
//	[IN | OUT | INOUT] param_name type
type ProcedureParam struct {
	Direction string    // "IN", "OUT", "INOUT", or "" (defaults to IN)
	Name      string
	Type      *TypeName
	Loc       Loc
}

// Tag implements Node.
func (n *ProcedureParam) Tag() NodeTag { return T_ProcedureParam }

// CreateProcedureStmt represents:
//
//	CREATE [OR REPLACE] PROCEDURE [IF NOT EXISTS] proc_name([param_list])
//	    [COMMENT 'text']
//	    BEGIN
//	        ... procedure body ...
//	    END
type CreateProcedureStmt struct {
	Name        *ObjectName
	IfNotExists bool
	OrReplace   bool
	Parameters  []*ProcedureParam
	Body        string // raw BEGIN...END body text
	Comment     string
	Loc         Loc
}

// Tag implements Node.
func (n *CreateProcedureStmt) Tag() NodeTag { return T_CreateProcedureStmt }

// CallProcedureStmt represents:
//
//	CALL proc_name([args])
type CallProcedureStmt struct {
	Name *ObjectName
	Args []Node // expression arguments
	Loc  Loc
}

// Tag implements Node.
func (n *CallProcedureStmt) Tag() NodeTag { return T_CallProcedureStmt }

// DropProcedureStmt represents:
//
//	DROP PROCEDURE [IF EXISTS] proc_name
type DropProcedureStmt struct {
	Name     *ObjectName
	IfExists bool
	Loc      Loc
}

// Tag implements Node.
func (n *DropProcedureStmt) Tag() NodeTag { return T_DropProcedureStmt }
