package parser

import (
	"github.com/bytebase/omni/trino/ast"
)

// This file is part of the `parser-routines` DAG node (with function_def.go):
// it implements the SQL-routine control-flow body language — the legacy
// TrinoParser.g4 `controlStatement`, `caseStatementWhenClause`, `elseIfClause`,
// `elseClause`, `variableDeclaration`, and `sqlStatementList` rules. The
// CREATE/DROP/WITH FUNCTION statement shells live in function_def.go.
//
// A routine body is the `controlStatement` that follows the function signature
// (functionSpecification = FUNCTION declaration returnsClause characteristic*
// controlStatement). Its alternatives are:
//
//	RETURN valueExpression                                                   # returnStatement
//	SET identifier = expression                                             # assignmentStatement
//	CASE expression caseWhen+ elseClause? END CASE                          # simpleCaseStatement
//	CASE caseWhen+ elseClause? END CASE                                     # searchedCaseStatement
//	IF expression THEN sqlStatementList elseIfClause* elseClause? END IF    # ifStatement
//	ITERATE identifier                                                      # iterateStatement
//	LEAVE identifier                                                        # leaveStatement
//	BEGIN (variableDeclaration ;)* sqlStatementList? END                    # compoundStatement
//	(label:)? LOOP sqlStatementList END LOOP                                # loopStatement
//	(label:)? WHILE expression DO sqlStatementList END WHILE                # whileStatement
//	(label:)? REPEAT sqlStatementList UNTIL expression END REPEAT           # repeatStatement
//
//	caseStatementWhenClause : WHEN expression THEN sqlStatementList
//	elseIfClause            : ELSEIF expression THEN sqlStatementList
//	elseClause              : ELSE sqlStatementList
//	variableDeclaration     : DECLARE identifier (, identifier)* type (DEFAULT valueExpression)?
//	sqlStatementList        : (controlStatement ;)+
//
// As with expr.go / datatypes.go, the body-statement node types are
// PARSER-PACKAGE types. They satisfy a parser-local RoutineStatement interface
// (Span() marker) rather than ast.Node — exactly how Expr is a parser-local
// interface — so no ast.NodeTag is minted for the inner control statements;
// only the three top-level FUNCTION statement shells (function_def.go) get tags.
//
// Adjudicated against the live Trino 481 oracle. Oracle-confirmed facts baked
// in (see oracle_routines_test.go for the differential corpus):
//
//	R1 (RETURN takes a valueExpression, NOT a full expression). Trino 481
//	   REJECTS `RETURN x IN (1,2,3)`, `RETURN x > 0 AND x < 10`, and
//	   `RETURN x BETWEEN 1 AND 10` with SYNTAX_ERROR — a bare predicate / boolean
//	   is not a valueExpression. parseRoutineStatement therefore calls
//	   parseValueExpr (not parseExpr) for RETURN. Likewise DEFAULT in a DECLARE
//	   is a valueExpression (`DECLARE x int DEFAULT 1 IN (1,2)` is rejected).
//	R2 (SET / IF / WHILE / CASE / UNTIL take a full expression). The condition
//	   positions DO accept predicates/booleans: `SET y = x IN (1,2)`,
//	   `IF x IN (1,2,3) THEN …`, `WHILE x > 0 AND x < 10 DO …`,
//	   `CASE WHEN x IN (1,2) THEN …`, `REPEAT … UNTIL i >= n AND i > 0 …` are all
//	   accepted. Those positions call parseExpr.
//	R3 (assignment target is a single identifier). `SET x = 1` is accepted but
//	   `SET x.y = 1` is a SYNTAX_ERROR — the grammar is `SET identifier = …`, not
//	   a qualifiedName. parseAssignmentStatement reads exactly one identifier.
//	R4 (ITERATE / LEAVE require a label). `LEAVE;` and `ITERATE;` with no label
//	   are SYNTAX_ERRORs in 481; the label identifier is mandatory.
//	R5 (labels attach only to LOOP / WHILE / REPEAT). `BEGIN lbl: BEGIN … END; END`
//	   is rejected — a `label:` prefix is NOT allowed on a compound statement.
//	R6 (DECLARE precede statements; both optional). In a compound, every
//	   variableDeclaration must come before the first sqlStatementList statement
//	   (`BEGIN SET x=1; DECLARE y int; …` is rejected). Both the declaration block
//	   and the statement list are optional: `BEGIN END` and `BEGIN DECLARE x int; END`
//	   are both accepted.
//	R7 (each statement in a sqlStatementList is terminated by ';'). Inside
//	   BEGIN/IF/LOOP/WHILE/REPEAT bodies every controlStatement ends with ';'
//	   (`BEGIN RETURN 1 END` — missing ';' — is rejected). The trailing
//	   controlStatement of the whole functionSpecification (the RETURN/BEGIN form)
//	   is NOT semicolon-terminated; that ';' belongs to the statement splitter.

// RoutineStatement is one statement of a SQL-routine body — a `controlStatement`
// grammar alternative. Like Expr it is a parser-local interface (no ast.NodeTag)
// because routine statements are only ever embedded inside a FUNCTION shell.
type RoutineStatement interface {
	// Span returns the source byte range covered by the statement.
	Span() ast.Loc
	// routineStmtNode is a marker preventing unrelated types from satisfying
	// the interface.
	routineStmtNode()
}

// ---------------------------------------------------------------------------
// RETURN
// ---------------------------------------------------------------------------

// ReturnStatement is `RETURN valueExpression` (R1). Value is the returned
// value-expression (predicates/booleans are NOT permitted bare — they would be
// a SYNTAX_ERROR; see R1).
type ReturnStatement struct {
	Value Expr
	Loc   ast.Loc
}

func (s *ReturnStatement) Span() ast.Loc    { return s.Loc }
func (s *ReturnStatement) routineStmtNode() {}

var _ RoutineStatement = (*ReturnStatement)(nil)

// ---------------------------------------------------------------------------
// SET (assignment)
// ---------------------------------------------------------------------------

// AssignmentStatement is `SET identifier = expression` (R2, R3). Target is the
// single assigned variable name (a qualified `x.y` target is rejected, R3);
// Value is a full expression (predicates allowed, R2).
type AssignmentStatement struct {
	Target *ast.Identifier
	Value  Expr
	Loc    ast.Loc
}

func (s *AssignmentStatement) Span() ast.Loc    { return s.Loc }
func (s *AssignmentStatement) routineStmtNode() {}

var _ RoutineStatement = (*AssignmentStatement)(nil)

// ---------------------------------------------------------------------------
// CASE (simple + searched)
// ---------------------------------------------------------------------------

// CaseWhenClause is `WHEN expression THEN sqlStatementList` (caseStatementWhenClause).
type CaseWhenClause struct {
	// Condition is the WHEN expression (a value to match in a simple CASE, a
	// boolean in a searched CASE — both are full expressions).
	Condition Expr
	// Body is the THEN statement list.
	Body []RoutineStatement
	Loc  ast.Loc
}

// CaseStatement is `CASE [operand] when+ [else] END CASE`. Operand is non-nil
// for the simple form (`CASE expr WHEN v THEN …`) and nil for the searched form
// (`CASE WHEN cond THEN …`).
type CaseStatement struct {
	Operand Expr // nil for the searched form
	Whens   []CaseWhenClause
	Else    []RoutineStatement // nil when no ELSE clause
	Loc     ast.Loc
}

func (s *CaseStatement) Span() ast.Loc    { return s.Loc }
func (s *CaseStatement) routineStmtNode() {}

var _ RoutineStatement = (*CaseStatement)(nil)

// ---------------------------------------------------------------------------
// IF / ELSEIF / ELSE
// ---------------------------------------------------------------------------

// ElseIfClause is `ELSEIF expression THEN sqlStatementList`.
type ElseIfClause struct {
	Condition Expr
	Body      []RoutineStatement
	Loc       ast.Loc
}

// IfStatement is `IF expression THEN body elseIf* [else] END IF`.
type IfStatement struct {
	Condition Expr
	Then      []RoutineStatement
	ElseIfs   []ElseIfClause
	Else      []RoutineStatement // nil when no ELSE clause
	Loc       ast.Loc
}

func (s *IfStatement) Span() ast.Loc    { return s.Loc }
func (s *IfStatement) routineStmtNode() {}

var _ RoutineStatement = (*IfStatement)(nil)

// ---------------------------------------------------------------------------
// ITERATE / LEAVE
// ---------------------------------------------------------------------------

// IterateStatement is `ITERATE identifier` (R4 — the label is mandatory).
type IterateStatement struct {
	Label *ast.Identifier
	Loc   ast.Loc
}

func (s *IterateStatement) Span() ast.Loc    { return s.Loc }
func (s *IterateStatement) routineStmtNode() {}

var _ RoutineStatement = (*IterateStatement)(nil)

// LeaveStatement is `LEAVE identifier` (R4 — the label is mandatory).
type LeaveStatement struct {
	Label *ast.Identifier
	Loc   ast.Loc
}

func (s *LeaveStatement) Span() ast.Loc    { return s.Loc }
func (s *LeaveStatement) routineStmtNode() {}

var _ RoutineStatement = (*LeaveStatement)(nil)

// ---------------------------------------------------------------------------
// BEGIN … END (compound)
// ---------------------------------------------------------------------------

// VariableDeclaration is `DECLARE identifier (, identifier)* type [DEFAULT valueExpression]`
// (variableDeclaration). One DECLARE may introduce several names of the same
// type. Default, when present, is a valueExpression (R1).
type VariableDeclaration struct {
	Names   []*ast.Identifier
	Type    *DataType
	Default Expr // nil when no DEFAULT
	Loc     ast.Loc
}

// CompoundStatement is `BEGIN (variableDeclaration ;)* sqlStatementList? END`
// (R5 — no label is permitted; R6 — declarations precede statements, both
// optional).
type CompoundStatement struct {
	Declarations []VariableDeclaration // nil when none
	Body         []RoutineStatement    // nil when the body is empty
	Loc          ast.Loc
}

func (s *CompoundStatement) Span() ast.Loc    { return s.Loc }
func (s *CompoundStatement) routineStmtNode() {}

var _ RoutineStatement = (*CompoundStatement)(nil)

// ---------------------------------------------------------------------------
// LOOP / WHILE / REPEAT (labeled iteration)
// ---------------------------------------------------------------------------

// LoopStatement is `(label:)? LOOP sqlStatementList END LOOP`. Label is non-nil
// when a `label:` prefix was present.
type LoopStatement struct {
	Label *ast.Identifier // nil when unlabeled
	Body  []RoutineStatement
	Loc   ast.Loc
}

func (s *LoopStatement) Span() ast.Loc    { return s.Loc }
func (s *LoopStatement) routineStmtNode() {}

var _ RoutineStatement = (*LoopStatement)(nil)

// WhileStatement is `(label:)? WHILE expression DO sqlStatementList END WHILE`.
type WhileStatement struct {
	Label     *ast.Identifier // nil when unlabeled
	Condition Expr
	Body      []RoutineStatement
	Loc       ast.Loc
}

func (s *WhileStatement) Span() ast.Loc    { return s.Loc }
func (s *WhileStatement) routineStmtNode() {}

var _ RoutineStatement = (*WhileStatement)(nil)

// RepeatStatement is `(label:)? REPEAT sqlStatementList UNTIL expression END REPEAT`.
// Condition is the UNTIL expression (a full expression, R2).
type RepeatStatement struct {
	Label     *ast.Identifier // nil when unlabeled
	Body      []RoutineStatement
	Condition Expr
	Loc       ast.Loc
}

func (s *RepeatStatement) Span() ast.Loc    { return s.Loc }
func (s *RepeatStatement) routineStmtNode() {}

var _ RoutineStatement = (*RepeatStatement)(nil)

// ---------------------------------------------------------------------------
// controlStatement dispatch
// ---------------------------------------------------------------------------

// parseRoutineStatement parses one `controlStatement` and returns it as a
// RoutineStatement. A `label:` prefix (any identifier — including a non-reserved
// keyword — followed by ':') selects one of the labeled iteration forms
// (LOOP / WHILE / REPEAT — R5); otherwise the leading keyword selects the
// alternative.
//
// The label check MUST precede the keyword switch: a label name may itself be a
// non-reserved control keyword (Trino 481 accepts `loop: LOOP …`,
// `while: LOOP …`, `set: LOOP …`, etc. — every dispatch keyword except the
// reserved CASE is a valid label). Dispatching on the keyword first would
// mis-read `loop: LOOP …` as an unlabeled LOOP starting at the label token.
func (p *Parser) parseRoutineStatement() (RoutineStatement, error) {
	// A `label:` prefix: an identifier (or non-reserved keyword) immediately
	// followed by ':' begins a labeled LOOP / WHILE / REPEAT.
	if isIdentifierStart(p.cur.Kind) && p.peekNext().Kind == int(':') {
		labelTok := p.advance() // consume the label identifier
		label := identFromToken(labelTok)
		p.advance() // consume ':'
		switch p.cur.Kind {
		case kwLOOP:
			return p.parseLoopStatement(label, label.Loc.Start)
		case kwWHILE:
			return p.parseWhileStatement(label, label.Loc.Start)
		case kwREPEAT:
			return p.parseRepeatStatement(label, label.Loc.Start)
		default:
			// Only LOOP / WHILE / REPEAT take a label (R5); a label on anything
			// else (e.g. `lbl: BEGIN`) is a syntax error.
			return nil, p.syntaxErrorAtCur()
		}
	}

	switch p.cur.Kind {
	case kwRETURN:
		return p.parseReturnStatement()
	case kwSET:
		return p.parseAssignmentStatement()
	case kwCASE:
		return p.parseCaseStatement()
	case kwIF:
		return p.parseIfStatement()
	case kwITERATE:
		return p.parseIterateStatement()
	case kwLEAVE:
		return p.parseLeaveStatement()
	case kwBEGIN:
		return p.parseCompoundStatement()
	case kwLOOP:
		return p.parseLoopStatement(nil, p.cur.Loc.Start)
	case kwWHILE:
		return p.parseWhileStatement(nil, p.cur.Loc.Start)
	case kwREPEAT:
		return p.parseRepeatStatement(nil, p.cur.Loc.Start)
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parseReturnStatement parses `RETURN valueExpression` (R1). RETURN is current.
func (p *Parser) parseReturnStatement() (RoutineStatement, error) {
	retTok := p.advance() // consume RETURN
	// R1: RETURN takes a valueExpression, not a full expression. A bare
	// predicate/boolean (`RETURN x IN (1,2)`) is a SYNTAX_ERROR in Trino 481.
	val, err := p.parseValueExpr()
	if err != nil {
		return nil, err
	}
	return &ReturnStatement{
		Value: val,
		Loc:   ast.Loc{Start: retTok.Loc.Start, End: val.Span().End},
	}, nil
}

// parseAssignmentStatement parses `SET identifier = expression` (R2, R3). SET is
// current.
func (p *Parser) parseAssignmentStatement() (RoutineStatement, error) {
	setTok := p.advance() // consume SET
	// R3: the assignment target is a single identifier; `SET x.y = …` is a
	// SYNTAX_ERROR (the grammar is SET identifier EQ expression, not a
	// qualifiedName).
	target, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(int('=')); err != nil {
		return nil, err
	}
	// R2: the assigned value is a full expression (predicates allowed).
	val, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return &AssignmentStatement{
		Target: target,
		Value:  val,
		Loc:    ast.Loc{Start: setTok.Loc.Start, End: val.Span().End},
	}, nil
}

// parseCaseStatement parses the simple (`CASE expr WHEN …`) and searched
// (`CASE WHEN …`) case statements; the two share the `when+ else? END CASE`
// tail. CASE is current.
func (p *Parser) parseCaseStatement() (RoutineStatement, error) {
	caseTok := p.advance() // consume CASE
	s := &CaseStatement{Loc: ast.Loc{Start: caseTok.Loc.Start, End: caseTok.Loc.End}}

	// Simple form: a subject expression precedes the first WHEN. Searched form:
	// WHEN comes immediately after CASE.
	if p.cur.Kind != kwWHEN {
		operand, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		s.Operand = operand
	}

	// At least one WHEN clause is required (when+).
	if p.cur.Kind != kwWHEN {
		return nil, p.syntaxErrorAtCur()
	}
	for p.cur.Kind == kwWHEN {
		w, err := p.parseCaseWhenClause()
		if err != nil {
			return nil, err
		}
		s.Whens = append(s.Whens, w)
	}

	if p.cur.Kind == kwELSE {
		body, err := p.parseElseClause()
		if err != nil {
			return nil, err
		}
		s.Else = body
	}

	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}
	endCaseTok, err := p.expect(kwCASE)
	if err != nil {
		return nil, err
	}
	s.Loc.End = endCaseTok.Loc.End
	return s, nil
}

// parseCaseWhenClause parses `WHEN expression THEN sqlStatementList`. WHEN is
// current.
func (p *Parser) parseCaseWhenClause() (CaseWhenClause, error) {
	whenTok := p.advance() // consume WHEN
	cond, err := p.parseExpr()
	if err != nil {
		return CaseWhenClause{}, err
	}
	if _, err := p.expect(kwTHEN); err != nil {
		return CaseWhenClause{}, err
	}
	body, err := p.parseSQLStatementList()
	if err != nil {
		return CaseWhenClause{}, err
	}
	end := whenTok.Loc.End
	if n := len(body); n > 0 {
		end = body[n-1].Span().End
	}
	return CaseWhenClause{
		Condition: cond,
		Body:      body,
		Loc:       ast.Loc{Start: whenTok.Loc.Start, End: end},
	}, nil
}

// parseElseClause parses `ELSE sqlStatementList` and returns the statement list.
// ELSE is current.
func (p *Parser) parseElseClause() ([]RoutineStatement, error) {
	p.advance() // consume ELSE
	return p.parseSQLStatementList()
}

// parseIfStatement parses
// `IF expression THEN sqlStatementList elseIfClause* elseClause? END IF`. IF is
// current.
func (p *Parser) parseIfStatement() (RoutineStatement, error) {
	ifTok := p.advance() // consume IF
	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwTHEN); err != nil {
		return nil, err
	}
	thenBody, err := p.parseSQLStatementList()
	if err != nil {
		return nil, err
	}
	s := &IfStatement{
		Condition: cond,
		Then:      thenBody,
		Loc:       ast.Loc{Start: ifTok.Loc.Start, End: ifTok.Loc.End},
	}

	for p.cur.Kind == kwELSEIF {
		elseIfTok := p.advance() // consume ELSEIF
		eiCond, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(kwTHEN); err != nil {
			return nil, err
		}
		eiBody, err := p.parseSQLStatementList()
		if err != nil {
			return nil, err
		}
		end := elseIfTok.Loc.End
		if n := len(eiBody); n > 0 {
			end = eiBody[n-1].Span().End
		}
		s.ElseIfs = append(s.ElseIfs, ElseIfClause{
			Condition: eiCond,
			Body:      eiBody,
			Loc:       ast.Loc{Start: elseIfTok.Loc.Start, End: end},
		})
	}

	if p.cur.Kind == kwELSE {
		elseBody, err := p.parseElseClause()
		if err != nil {
			return nil, err
		}
		s.Else = elseBody
	}

	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}
	endIfTok, err := p.expect(kwIF)
	if err != nil {
		return nil, err
	}
	s.Loc.End = endIfTok.Loc.End
	return s, nil
}

// parseIterateStatement parses `ITERATE identifier` (R4). ITERATE is current.
func (p *Parser) parseIterateStatement() (RoutineStatement, error) {
	iterTok := p.advance() // consume ITERATE
	// R4: the label is mandatory; `ITERATE` with no label is a SYNTAX_ERROR.
	label, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	return &IterateStatement{
		Label: label,
		Loc:   ast.Loc{Start: iterTok.Loc.Start, End: label.Loc.End},
	}, nil
}

// parseLeaveStatement parses `LEAVE identifier` (R4). LEAVE is current.
func (p *Parser) parseLeaveStatement() (RoutineStatement, error) {
	leaveTok := p.advance() // consume LEAVE
	// R4: the label is mandatory; `LEAVE` with no label is a SYNTAX_ERROR.
	label, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	return &LeaveStatement{
		Label: label,
		Loc:   ast.Loc{Start: leaveTok.Loc.Start, End: label.Loc.End},
	}, nil
}

// parseCompoundStatement parses
// `BEGIN (variableDeclaration ;)* sqlStatementList? END` (R5, R6). BEGIN is
// current. No label is permitted (R5 — the caller never reaches here via the
// `label:` path).
func (p *Parser) parseCompoundStatement() (RoutineStatement, error) {
	beginTok := p.advance() // consume BEGIN
	s := &CompoundStatement{Loc: ast.Loc{Start: beginTok.Loc.Start, End: beginTok.Loc.End}}

	// R6: zero or more `DECLARE … ;` declarations, all before any statement.
	for p.cur.Kind == kwDECLARE {
		decl, err := p.parseVariableDeclaration()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(int(';')); err != nil {
			return nil, err
		}
		s.Declarations = append(s.Declarations, decl)
	}

	// R6: an optional sqlStatementList. It is absent when END follows the
	// declarations directly (`BEGIN END`, `BEGIN DECLARE x int; END`).
	if p.cur.Kind != kwEND {
		body, err := p.parseSQLStatementList()
		if err != nil {
			return nil, err
		}
		s.Body = body
	}

	endTok, err := p.expect(kwEND)
	if err != nil {
		return nil, err
	}
	s.Loc.End = endTok.Loc.End
	return s, nil
}

// parseVariableDeclaration parses
// `DECLARE identifier (, identifier)* type [DEFAULT valueExpression]`. DECLARE is
// current. The terminating ';' is consumed by the caller (compound statement).
func (p *Parser) parseVariableDeclaration() (VariableDeclaration, error) {
	declTok := p.advance() // consume DECLARE
	first, err := p.parseIdentifier()
	if err != nil {
		return VariableDeclaration{}, err
	}
	names := []*ast.Identifier{first}
	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		next, err := p.parseIdentifier()
		if err != nil {
			return VariableDeclaration{}, err
		}
		names = append(names, next)
	}

	typ, err := p.parseType()
	if err != nil {
		return VariableDeclaration{}, err
	}

	decl := VariableDeclaration{
		Names: names,
		Type:  typ,
		Loc:   ast.Loc{Start: declTok.Loc.Start, End: typ.Loc.End},
	}

	if p.cur.Kind == kwDEFAULT {
		p.advance() // consume DEFAULT
		// R1: DEFAULT is a valueExpression, not a full expression.
		def, err := p.parseValueExpr()
		if err != nil {
			return VariableDeclaration{}, err
		}
		decl.Default = def
		decl.Loc.End = def.Span().End
	}

	return decl, nil
}

// parseLoopStatement parses `LOOP sqlStatementList END LOOP`. LOOP is current;
// label (may be nil) and start are supplied by the dispatcher (so the node span
// covers a leading `label:`).
func (p *Parser) parseLoopStatement(label *ast.Identifier, start int) (RoutineStatement, error) {
	p.advance() // consume LOOP
	body, err := p.parseSQLStatementList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}
	endLoopTok, err := p.expect(kwLOOP)
	if err != nil {
		return nil, err
	}
	return &LoopStatement{
		Label: label,
		Body:  body,
		Loc:   ast.Loc{Start: start, End: endLoopTok.Loc.End},
	}, nil
}

// parseWhileStatement parses `WHILE expression DO sqlStatementList END WHILE`.
// WHILE is current; label (may be nil) and start come from the dispatcher.
func (p *Parser) parseWhileStatement(label *ast.Identifier, start int) (RoutineStatement, error) {
	p.advance() // consume WHILE
	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwDO); err != nil {
		return nil, err
	}
	body, err := p.parseSQLStatementList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}
	endWhileTok, err := p.expect(kwWHILE)
	if err != nil {
		return nil, err
	}
	return &WhileStatement{
		Label:     label,
		Condition: cond,
		Body:      body,
		Loc:       ast.Loc{Start: start, End: endWhileTok.Loc.End},
	}, nil
}

// parseRepeatStatement parses
// `REPEAT sqlStatementList UNTIL expression END REPEAT`. REPEAT is current;
// label (may be nil) and start come from the dispatcher.
func (p *Parser) parseRepeatStatement(label *ast.Identifier, start int) (RoutineStatement, error) {
	p.advance() // consume REPEAT
	body, err := p.parseSQLStatementList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwUNTIL); err != nil {
		return nil, err
	}
	// R2: the UNTIL condition is a full expression.
	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}
	endRepeatTok, err := p.expect(kwREPEAT)
	if err != nil {
		return nil, err
	}
	return &RepeatStatement{
		Label:     label,
		Body:      body,
		Condition: cond,
		Loc:       ast.Loc{Start: start, End: endRepeatTok.Loc.End},
	}, nil
}

// parseSQLStatementList parses `(controlStatement ;)+` (sqlStatementList). At
// least one statement is required, and EACH statement is terminated by ';' (R7).
// The list ends when the next token is not a controlStatement start (e.g. END,
// ELSE, ELSEIF, UNTIL).
func (p *Parser) parseSQLStatementList() ([]RoutineStatement, error) {
	var stmts []RoutineStatement
	for {
		stmt, err := p.parseRoutineStatement()
		if err != nil {
			return nil, err
		}
		// R7: every statement in a sqlStatementList is ';'-terminated.
		if _, err := p.expect(int(';')); err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
		if !p.routineStatementStarts() {
			break
		}
	}
	return stmts, nil
}

// routineStatementStarts reports whether the current token can begin another
// controlStatement, so parseSQLStatementList knows when to stop. It mirrors the
// dispatch in parseRoutineStatement: the fixed leading keywords plus a `label:`
// prefix (identifier followed by ':'). Terminators (END, ELSE, ELSEIF, UNTIL)
// return false.
func (p *Parser) routineStatementStarts() bool {
	switch p.cur.Kind {
	case kwRETURN, kwSET, kwCASE, kwIF, kwITERATE, kwLEAVE, kwBEGIN, kwLOOP, kwWHILE, kwREPEAT:
		return true
	default:
		return isIdentifierStart(p.cur.Kind) && p.peekNext().Kind == int(':')
	}
}
