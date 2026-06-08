package ast

// NodeTag identifies the type of a PL/pgSQL AST node.
type NodeTag int

// NodeTag constants for all PL/pgSQL AST node types.
const (
	T_Invalid NodeTag = 0

	T_PLBlock NodeTag = iota + 1000 // offset from pg/ast tags to avoid collision
	T_PLDeclare
	T_PLCursorDecl
	T_PLAliasDecl
	T_PLIf
	T_PLCase
	T_PLCaseWhen
	T_PLLoop
	T_PLWhile
	T_PLForI
	T_PLForS
	T_PLForC
	T_PLForDynS
	T_PLForEachA
	T_PLReturn
	T_PLReturnNext
	T_PLReturnQuery
	T_PLAssign
	T_PLExecSQL
	T_PLDynExecute
	T_PLPerform
	T_PLCall
	T_PLRaise
	T_PLAssert
	T_PLGetDiag
	T_PLGetDiagItem
	T_PLOpen
	T_PLFetch
	T_PLClose
	T_PLExit
	T_PLCommit
	T_PLRollback
	T_PLNull
	T_PLException
	T_PLExceptionWhen
	T_PLElsIf
	T_PLRaiseOption
)
