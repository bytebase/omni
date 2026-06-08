// Package ast defines PL/pgSQL AST node types.
// These represent the parse tree for PL/pgSQL function bodies.
package ast

// Node is the interface implemented by all PL/pgSQL AST nodes.
type Node interface {
	// Tag returns the NodeTag for this node type.
	Tag() NodeTag
}

// Loc represents a source location range (byte offsets).
// -1 means "unknown" for either field.
type Loc struct {
	Start int // inclusive start byte offset (-1 if unknown)
	End   int // exclusive end byte offset (-1 if unknown)
}

// NoLoc returns a Loc with both Start and End set to -1 (unknown).
func NoLoc() Loc {
	return Loc{Start: -1, End: -1}
}

// ScrollOption represents the SCROLL option for cursors.
type ScrollOption int

const (
	ScrollNone ScrollOption = iota // no SCROLL option specified
	ScrollYes                      // SCROLL
	ScrollNo                       // NO SCROLL
)

// RaiseLevel represents the severity level for RAISE statements.
type RaiseLevel int

const (
	RaiseLevelNone    RaiseLevel = iota // not specified (defaults to EXCEPTION)
	RaiseLevelDebug                     // DEBUG
	RaiseLevelLog                       // LOG
	RaiseLevelInfo                      // INFO
	RaiseLevelNotice                    // NOTICE
	RaiseLevelWarning                   // WARNING
	RaiseLevelError                     // EXCEPTION
)

// FetchDirection represents the direction for FETCH/MOVE statements.
type FetchDirection int

const (
	FetchNext        FetchDirection = iota // NEXT (default)
	FetchPrior                             // PRIOR
	FetchFirst                             // FIRST
	FetchLast                              // LAST
	FetchAbsolute                          // ABSOLUTE n
	FetchRelative                          // RELATIVE n
	FetchForward                           // FORWARD
	FetchForwardN                          // FORWARD n
	FetchForwardAll                        // FORWARD ALL
	FetchBackward                          // BACKWARD
	FetchBackwardN                         // BACKWARD n
	FetchBackwardAll                       // BACKWARD ALL
)

// -------------------------------------------------------------------
// Statement nodes
// -------------------------------------------------------------------

// PLBlock represents a PL/pgSQL block (BEGIN...END).
// PLOption represents a compiler option directive (#option, #print_strict_params, #variable_conflict).
type PLOption struct {
	Name  string // directive name (e.g. "dump", "print_strict_params", "variable_conflict")
	Value string // directive value (e.g. "on", "off", "error", "use_variable", "use_column")
	Loc   Loc    // source location range
}

// PLBlock represents a PL/pgSQL block (BEGIN...END).
type PLBlock struct {
	Label        string     // optional label (empty if unlabeled)
	Declarations []Node     // variable/cursor declarations (PLDeclare, PLCursorDecl, PLAliasDecl)
	Body         []Node     // statement list
	Exceptions   []Node     // exception handler clauses (PLExceptionWhen)
	Options      []PLOption // compiler option directives before the block
	Loc          Loc        // source location range
}

func (n *PLBlock) Tag() NodeTag { return T_PLBlock }

// PLDeclare represents a variable declaration.
type PLDeclare struct {
	Name      string // variable name
	TypeName  string // type reference as text (e.g. "integer", "public.my_type", "my_table.col%TYPE")
	Constant  bool   // CONSTANT flag
	NotNull   bool   // NOT NULL flag
	Collation string // COLLATE clause (empty if none)
	Default   string // default expression as text (empty if none)
	Loc       Loc    // source location range
}

func (n *PLDeclare) Tag() NodeTag { return T_PLDeclare }

// PLCursorDecl represents a cursor declaration.
type PLCursorDecl struct {
	Name   string        // cursor name
	Scroll ScrollOption  // SCROLL / NO SCROLL / none
	Args   []PLCursorArg // cursor arguments (may be nil)
	Query  string        // query text after FOR/IS
	Loc    Loc           // source location range
}

func (n *PLCursorDecl) Tag() NodeTag { return T_PLCursorDecl }

// PLCursorArg represents a single cursor argument in a cursor declaration.
type PLCursorArg struct {
	Name     string // argument name
	TypeName string // argument type as text
}

// PLAliasDecl represents an ALIAS declaration.
type PLAliasDecl struct {
	Name    string // alias name
	RefName string // referenced variable or parameter (e.g. "$1", "param_name")
	Loc     Loc    // source location range
}

func (n *PLAliasDecl) Tag() NodeTag { return T_PLAliasDecl }

// PLIf represents an IF / ELSIF / ELSE statement.
type PLIf struct {
	Condition string // condition expression text (between IF and THEN)
	ThenBody  []Node // THEN statement list
	ElsIfs    []Node // ELSIF clauses (PLElsIf nodes)
	ElseBody  []Node // ELSE statement list (nil if no ELSE)
	Loc       Loc    // source location range
}

func (n *PLIf) Tag() NodeTag { return T_PLIf }

// PLElsIf represents an ELSIF clause within an IF statement.
type PLElsIf struct {
	Condition string // condition expression text
	Body      []Node // statement list
	Loc       Loc    // source location range
}

func (n *PLElsIf) Tag() NodeTag { return T_PLElsIf }

// PLCase represents a CASE statement (searched or simple).
type PLCase struct {
	TestExpr string // test expression text (empty for searched CASE)
	HasTest  bool   // true if this is a simple CASE (has test expression)
	Whens    []Node // WHEN clauses (PLCaseWhen nodes)
	ElseBody []Node // ELSE statement list (nil if no ELSE)
	Loc      Loc    // source location range
}

func (n *PLCase) Tag() NodeTag { return T_PLCase }

// PLCaseWhen represents a WHEN clause in a CASE statement.
type PLCaseWhen struct {
	Expr string // WHEN expression text
	Body []Node // statement list
	Loc  Loc    // source location range
}

func (n *PLCaseWhen) Tag() NodeTag { return T_PLCaseWhen }

// PLLoop represents an unconditional LOOP statement.
type PLLoop struct {
	Label string // optional label (empty if unlabeled)
	Body  []Node // statement list
	Loc   Loc    // source location range
}

func (n *PLLoop) Tag() NodeTag { return T_PLLoop }

// PLWhile represents a WHILE loop.
type PLWhile struct {
	Label     string // optional label
	Condition string // condition expression text (between WHILE and LOOP)
	Body      []Node // statement list
	Loc       Loc    // source location range
}

func (n *PLWhile) Tag() NodeTag { return T_PLWhile }

// PLForI represents an integer FOR loop (FOR i IN lower..upper).
type PLForI struct {
	Label   string // optional label
	Var     string // loop variable name
	Lower   string // lower bound expression text
	Upper   string // upper bound expression text
	Step    string // step expression text (empty if no BY clause)
	Reverse bool   // REVERSE flag
	Body    []Node // statement list
	Loc     Loc    // source location range
}

func (n *PLForI) Tag() NodeTag { return T_PLForI }

// PLForS represents a query FOR loop (FOR rec IN SELECT ...).
type PLForS struct {
	Label string // optional label
	Var   string // loop variable name
	Query string // query text
	Body  []Node // statement list
	Loc   Loc    // source location range
}

func (n *PLForS) Tag() NodeTag { return T_PLForS }

// PLForC represents a cursor FOR loop (FOR rec IN cursor_var).
type PLForC struct {
	Label     string // optional label
	Var       string // loop variable name
	CursorVar string // cursor variable name
	ArgQuery  string // argument expression text (empty if no args)
	Body      []Node // statement list
	Loc       Loc    // source location range
}

func (n *PLForC) Tag() NodeTag { return T_PLForC }

// PLForDynS represents a dynamic SQL FOR loop (FOR rec IN EXECUTE ...).
type PLForDynS struct {
	Label  string   // optional label
	Var    string   // loop variable name
	Query  string   // dynamic query expression text
	Params []string // USING parameter expressions (nil if no USING)
	Body   []Node   // statement list
	Loc    Loc      // source location range
}

func (n *PLForDynS) Tag() NodeTag { return T_PLForDynS }

// PLForEachA represents a FOREACH ... IN ARRAY loop.
type PLForEachA struct {
	Label     string // optional label
	Var       string // loop variable name
	SliceDim  int    // SLICE dimension (0 if no SLICE)
	ArrayExpr string // array expression text
	Body      []Node // statement list
	Loc       Loc    // source location range
}

func (n *PLForEachA) Tag() NodeTag { return T_PLForEachA }

// PLReturn represents a RETURN statement.
type PLReturn struct {
	Expr string // expression text (empty for bare RETURN)
	Loc  Loc    // source location range
}

func (n *PLReturn) Tag() NodeTag { return T_PLReturn }

// PLReturnNext represents a RETURN NEXT statement.
type PLReturnNext struct {
	Expr string // expression text (empty for bare RETURN NEXT with OUT params)
	Loc  Loc    // source location range
}

func (n *PLReturnNext) Tag() NodeTag { return T_PLReturnNext }

// PLReturnQuery represents a RETURN QUERY statement.
type PLReturnQuery struct {
	Query    string   // static query text (empty if dynamic)
	DynQuery string   // dynamic query expression text (empty if static)
	Params   []string // USING parameter expressions for dynamic query (nil if static)
	Loc      Loc      // source location range
}

func (n *PLReturnQuery) Tag() NodeTag { return T_PLReturnQuery }

// PLAssign represents a variable assignment statement.
type PLAssign struct {
	Target string // variable target (may include field/subscript: "rec.field", "arr[1]")
	Expr   string // expression text (RHS of :=)
	Loc    Loc    // source location range
}

func (n *PLAssign) Tag() NodeTag { return T_PLAssign }

// PLExecSQL represents a SQL statement executed inline (INSERT, UPDATE, DELETE, SELECT, etc.).
type PLExecSQL struct {
	SQLText string   // the SQL statement text
	Into    []string // INTO target variable names (nil if no INTO)
	Strict  bool     // INTO STRICT flag
	Loc     Loc      // source location range
}

func (n *PLExecSQL) Tag() NodeTag { return T_PLExecSQL }

// PLDynExecute represents an EXECUTE dynamic SQL statement.
type PLDynExecute struct {
	Query  string   // dynamic query expression text
	Into   []string // INTO target variable names (nil if no INTO)
	Strict bool     // INTO STRICT flag
	Params []string // USING parameter expressions (nil if no USING)
	Loc    Loc      // source location range
}

func (n *PLDynExecute) Tag() NodeTag { return T_PLDynExecute }

// PLPerform represents a PERFORM statement.
type PLPerform struct {
	Expr string // expression text (everything after PERFORM)
	Loc  Loc    // source location range
}

func (n *PLPerform) Tag() NodeTag { return T_PLPerform }

// PLCall represents a CALL or DO statement.
type PLCall struct {
	SQLText string // the full CALL/DO statement text
	Loc     Loc    // source location range
}

func (n *PLCall) Tag() NodeTag { return T_PLCall }

// PLRaise represents a RAISE statement.
type PLRaise struct {
	Level    RaiseLevel // severity level (DEBUG, LOG, INFO, NOTICE, WARNING, EXCEPTION)
	CondName string     // condition name (e.g. "division_by_zero") or empty
	SQLState string     // SQLSTATE code (e.g. "22012") or empty
	Message  string     // message format text or empty
	Params   []string   // message format parameters (% substitutions)
	Options  []Node     // USING options (PLRaiseOption nodes)
	Loc      Loc        // source location range
}

func (n *PLRaise) Tag() NodeTag { return T_PLRaise }

// PLRaiseOption represents a single USING option in a RAISE statement.
type PLRaiseOption struct {
	OptType string // option name (MESSAGE, DETAIL, HINT, ERRCODE, COLUMN, CONSTRAINT, DATATYPE, TABLE, SCHEMA)
	Expr    string // option value expression text
	Loc     Loc    // source location range
}

func (n *PLRaiseOption) Tag() NodeTag { return T_PLRaiseOption }

// PLAssert represents an ASSERT statement.
type PLAssert struct {
	Condition string // condition expression text
	Message   string // message expression text (empty if none)
	Loc       Loc    // source location range
}

func (n *PLAssert) Tag() NodeTag { return T_PLAssert }

// PLGetDiag represents a GET [CURRENT|STACKED] DIAGNOSTICS statement.
type PLGetDiag struct {
	IsStacked bool   // true for GET STACKED DIAGNOSTICS
	Items     []Node // diagnostic items (PLGetDiagItem nodes)
	Loc       Loc    // source location range
}

func (n *PLGetDiag) Tag() NodeTag { return T_PLGetDiag }

// PLGetDiagItem represents a single item in a GET DIAGNOSTICS statement.
type PLGetDiagItem struct {
	Target string // target variable name
	Kind   string // diagnostic item kind (e.g. "ROW_COUNT", "MESSAGE_TEXT")
	Loc    Loc    // source location range
}

func (n *PLGetDiagItem) Tag() NodeTag { return T_PLGetDiagItem }

// PLOpen represents an OPEN cursor statement.
type PLOpen struct {
	CursorVar string       // cursor variable name
	Query     string       // static query text (empty if bound cursor or dynamic)
	DynQuery  string       // dynamic query expression text (empty if static)
	Params    []string     // USING parameter expressions for dynamic query
	Scroll    ScrollOption // SCROLL / NO SCROLL option
	Args      string       // argument expression text for bound cursor (empty if none)
	Loc       Loc          // source location range
}

func (n *PLOpen) Tag() NodeTag { return T_PLOpen }

// PLFetch represents a FETCH or MOVE cursor statement.
type PLFetch struct {
	CursorVar string         // cursor variable name
	Direction FetchDirection // fetch direction
	Count     string         // direction count expression (for ABSOLUTE n, RELATIVE n, FORWARD n, BACKWARD n)
	Into      []string       // INTO target variable names (nil for MOVE)
	IsMove    bool           // true if this is MOVE (no INTO)
	Loc       Loc            // source location range
}

func (n *PLFetch) Tag() NodeTag { return T_PLFetch }

// PLClose represents a CLOSE cursor statement.
type PLClose struct {
	CursorVar string // cursor variable name
	Loc       Loc    // source location range
}

func (n *PLClose) Tag() NodeTag { return T_PLClose }

// PLExit represents an EXIT or CONTINUE statement.
type PLExit struct {
	IsContinue bool   // true for CONTINUE, false for EXIT
	Label      string // target label (empty for innermost loop)
	Condition  string // WHEN condition expression text (empty if unconditional)
	Loc        Loc    // source location range
}

func (n *PLExit) Tag() NodeTag { return T_PLExit }

// PLCommit represents a COMMIT statement.
type PLCommit struct {
	Chain bool // AND CHAIN flag
	Loc   Loc  // source location range
}

func (n *PLCommit) Tag() NodeTag { return T_PLCommit }

// PLRollback represents a ROLLBACK statement.
type PLRollback struct {
	Chain bool // AND CHAIN flag
	Loc   Loc  // source location range
}

func (n *PLRollback) Tag() NodeTag { return T_PLRollback }

// PLNull represents the NULL (no-op) statement.
type PLNull struct {
	Loc Loc // source location range
}

func (n *PLNull) Tag() NodeTag { return T_PLNull }

// PLException represents the EXCEPTION section of a block.
// (This is not used directly as a top-level node; PLBlock.Exceptions
// holds PLExceptionWhen nodes. Kept for potential future use.)
type PLException struct {
	Whens []Node // list of PLExceptionWhen nodes
	Loc   Loc    // source location range
}

func (n *PLException) Tag() NodeTag { return T_PLException }

// PLExceptionWhen represents a single WHEN clause in an EXCEPTION block.
type PLExceptionWhen struct {
	Conditions []string // condition names or SQLSTATE codes (joined by OR in source)
	Body       []Node   // handler statement list
	Loc        Loc      // source location range
}

func (n *PLExceptionWhen) Tag() NodeTag { return T_PLExceptionWhen }
