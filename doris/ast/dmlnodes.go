package ast

// This file holds AST node types for DML statements (T4.x).

// ---------------------------------------------------------------------------
// INSERT statement
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
