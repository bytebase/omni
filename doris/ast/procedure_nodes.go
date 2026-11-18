package ast

// This file holds AST node types for stored procedure DDL/DML (T8.2):
// CREATE PROCEDURE, CALL, and DROP PROCEDURE.

// ProcedureParam represents one parameter in a CREATE PROCEDURE parameter list.
//
//	[IN | OUT | INOUT] param_name type
type ProcedureParam struct {
	Direction string // "IN", "OUT", "INOUT", or "" (defaults to IN)
	Name      string
	Type      *TypeName
	Loc       Loc
}

// Tag implements Node.
func (n *ProcedureParam) Tag() NodeTag { return T_ProcedureParam }

// CreateProcedureStmt represents the CREATE and REPLACE forms of the
// grammar's createProcedure alternative (DorisParser.g4:43):
//
//	CREATE [OR REPLACE] PROCEDURE proc_name([param_list])
//	    [COMMENT 'text']
//	    BEGIN
//	        ... procedure body ...
//	    END
//	REPLACE PROCEDURE proc_name([param_list]) ...
//
// The rule reads the parameter list and what follows it as free text
// (LEFT_PAREN .*? RIGHT_PAREN .*?), so a statement that is not a parameter
// list and a BEGIN...END body — CREATE PROCEDURE a(zzz), CREATE PROCEDURE
// a() zzz — parses too; its text is kept in RawText. The rule has no
// IF NOT EXISTS, so neither does this node.
type CreateProcedureStmt struct {
	Name       *ObjectName
	OrReplace  bool
	Replace    bool // the REPLACE verb, without CREATE
	Parameters []*ProcedureParam
	Body       string // raw BEGIN...END body text
	RawText    string // the free text when it is not a parameter list and body
	Comment    string
	Loc        Loc
}

// Tag implements Node.
func (n *CreateProcedureStmt) Tag() NodeTag { return T_CreateProcedureStmt }

// AlterProcedureStmt represents ALTER PROCEDURE, the third verb of the
// grammar's createProcedure alternative (DorisParser.g4:43), which reads
// everything after the name as free text:
//
//	ALTER PROCEDURE proc_name(...)
type AlterProcedureStmt struct {
	Name    *ObjectName
	RawText string // the text after the name, as written
	Loc     Loc
}

// Tag implements Node.
func (n *AlterProcedureStmt) Tag() NodeTag { return T_AlterProcedureStmt }

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
//	DROP (PROCEDURE | PROC) [IF EXISTS] proc_name
type DropProcedureStmt struct {
	Name     *ObjectName
	IfExists bool
	Loc      Loc
}

// Tag implements Node.
func (n *DropProcedureStmt) Tag() NodeTag { return T_DropProcedureStmt }
