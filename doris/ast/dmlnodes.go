package ast

// This file holds AST node types for DML statements (T4.x).

// ---------------------------------------------------------------------------
// INSERT statement (T4.1)
// ---------------------------------------------------------------------------

// InsertStmt represents an INSERT INTO or INSERT OVERWRITE TABLE statement.
//
//	INSERT [INTO | OVERWRITE TABLE] [TEMPORARY PARTITION] table_name
//	    [PARTITION(p1, p2, ...) | PARTITION(*)]
//	    [WITH LABEL label_name]
//	    [(col1, col2, ...)]
//	    { VALUES (expr, ...) [, (...)] | SELECT ... | WITH ... SELECT ... }
type InsertStmt struct {
	Overwrite     bool        // true for INSERT OVERWRITE TABLE; false for INSERT INTO
	Target        *ObjectName // target table name
	Label         string      // optional WITH LABEL name; empty if absent
	Partition     []string    // optional PARTITION(p1, p2, ...) names; nil if absent
	TempPartition bool        // true if TEMPORARY PARTITION was specified
	PartitionStar bool        // true if PARTITION(*) was used
	Columns       []string    // optional column list; nil if absent
	// Source: exactly one of Values or Query is non-nil.
	Values [][]Node // VALUES rows: each inner slice is one row's expressions
	Query  Node     // SELECT or WITH…SELECT (*SelectStmt); nil if VALUES form
	Loc    Loc
}

// Tag implements Node.
func (n *InsertStmt) Tag() NodeTag { return T_InsertStmt }

// Compile-time assertion that *InsertStmt satisfies Node.
var _ Node = (*InsertStmt)(nil)

// ---------------------------------------------------------------------------
// Assignment: col = expr (used in UPDATE SET clause) (T4.2)
// ---------------------------------------------------------------------------

// Assignment represents a single SET assignment in an UPDATE statement:
//
//	col = expr
//	t.col = expr
type Assignment struct {
	Column *ObjectName // column name; may be qualified (t.col)
	Value  Node        // right-hand expression
	Loc    Loc
}

// Tag implements Node.
func (n *Assignment) Tag() NodeTag { return T_Assignment }

// Compile-time assertion that *Assignment satisfies Node.
var _ Node = (*Assignment)(nil)

// ---------------------------------------------------------------------------
// UpdateStmt (T4.2)
// ---------------------------------------------------------------------------

// UpdateStmt represents an UPDATE statement:
//
//	UPDATE table [AS alias]
//	    SET col1 = expr1 [, col2 = expr2 ...]
//	    [FROM table_refs]
//	    [WHERE condition]
type UpdateStmt struct {
	Target      *ObjectName   // table to update
	TargetAlias string        // optional alias (AS alias or bare alias)
	Assignments []*Assignment // SET clause assignments
	From        []Node        // optional FROM clause table references
	Where       Node          // optional WHERE expression
	Loc         Loc
}

// Tag implements Node.
func (n *UpdateStmt) Tag() NodeTag { return T_UpdateStmt }

// Compile-time assertion that *UpdateStmt satisfies Node.
var _ Node = (*UpdateStmt)(nil)

// ---------------------------------------------------------------------------
// DeleteStmt (T4.2)
// ---------------------------------------------------------------------------

// DeleteStmt represents a DELETE statement:
//
//	DELETE FROM table [AS alias]
//	    [PARTITION(p1 [, p2 ...])]
//	    [USING table_refs]
//	    [WHERE condition]
type DeleteStmt struct {
	Target      *ObjectName // table to delete from
	TargetAlias string      // optional alias
	Partition   []string    // optional PARTITION(p1, p2, ...) names
	Using       []Node      // optional USING clause table references
	Where       Node        // optional WHERE expression
	Loc         Loc
}

// Tag implements Node.
func (n *DeleteStmt) Tag() NodeTag { return T_DeleteStmt }

// Compile-time assertion that *DeleteStmt satisfies Node.
var _ Node = (*DeleteStmt)(nil)
