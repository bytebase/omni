package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parsePLSQLBlock parses a PL/SQL block:
//
//	[<<label>>] [DECLARE declarations] BEGIN statements [EXCEPTION handlers] END [label] ;
func (p *Parser) parsePLSQLBlock() (*nodes.PLSQLBlock, error) {
	start := p.pos()
	block := &nodes.PLSQLBlock{
		Loc: nodes.Loc{Start: start},
	}

	// Optional label: <<label>>
	if p.cur.Type == tokLABELOPEN {
		p.advance()
		var // consume <<
		parseErr829 error
		block.Label, parseErr829 = p.parseIdentifier()
		if parseErr829 != nil {
			return nil, parseErr829
		}
		if p.cur.Type == tokLABELCLOSE {
			p.advance() // consume >>
		}
	}

	// DECLARE section (optional)
	if p.cur.Type == kwDECLARE {
		p.advance()
		var // consume DECLARE
		parseErr830 error
		block.Declarations, parseErr830 = p.parsePLSQLDeclarations()
		if parseErr830 !=

			// BEGIN
			nil {
			return nil, parseErr830
		}
	} else if p.cur.Type != kwBEGIN && p.cur.Type != tokEOF {
		var parseErr830 error
		block.Declarations, parseErr830 = p.parsePLSQLDeclarations()
		if parseErr830 != nil {
			return nil, parseErr830
		}
	}

	if p.cur.Type != kwBEGIN {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume BEGIN
	var parseErr831 error

	// Statements
	block.Statements, parseErr831 = p.parsePLSQLStatements()
	if parseErr831 !=

		// EXCEPTION section (optional)
		nil {
		return nil, parseErr831
	}

	if p.cur.Type == kwEXCEPTION {
		p.advance()
		var // consume EXCEPTION
		parseErr832 error
		block.Exceptions, parseErr832 = p.parsePLSQLExceptionHandlers()
		if parseErr832 !=

			// END [label] ;
			nil {
			return nil, parseErr832
		}
	}

	if p.cur.Type != kwEND {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume END
	// Optional label after END
	if p.isIdentLike() && p.cur.Type != ';' && p.cur.Type != tokEOF {
		p.advance() // consume label
	}
	if p.cur.Type == ';' {
		p.advance() // consume ;
	}

	block.Loc.End = p.prev.End
	return block, nil
}

// parsePLSQLDeclarations parses the DECLARE section of a PL/SQL block.
// Stops when BEGIN is encountered.
func (p *Parser) parsePLSQLDeclarations() (*nodes.List, error) {
	decls := &nodes.List{}

	for p.cur.Type != kwBEGIN && p.cur.Type != tokEOF {
		decl, parseErr833 := p.parsePLSQLDeclaration()
		if parseErr833 != nil {
			return nil, parseErr833
		}
		if decl == nil {
			break
		}
		decls.Items = append(decls.Items, decl)
	}

	return decls, nil
}

// parsePLSQLDeclaration parses a single declaration.
func (p *Parser) parsePLSQLDeclaration() (nodes.Node, error) {
	// CURSOR declaration
	if p.cur.Type == kwCURSOR {
		return p.parsePLSQLCursorDecl()
	}

	// PRAGMA directive
	if p.cur.Type == kwPRAGMA {
		return p.parsePLSQLPragma()
	}

	// TYPE declaration
	if p.cur.Type == kwTYPE {
		next := p.peekNext()
		if next.Type != kwBODY { // not CREATE TYPE BODY
			return p.parsePLSQLTypeDecl()
		}
	}

	// Variable declaration: name [CONSTANT] type [NOT NULL] [:= | DEFAULT expr] ;
	if p.isIdentLike() {
		return p.parsePLSQLVarDecl()
	}

	return nil, nil
}

// parsePLSQLVarDecl parses a variable declaration.
func (p *Parser) parsePLSQLVarDecl() (*nodes.PLSQLVarDecl, error) {
	start := p.pos()
	decl := &nodes.PLSQLVarDecl{
		Loc: nodes.Loc{Start: start},
	}
	var parseErr834 error

	decl.Name, parseErr834 = p.parseIdentifier()
	if parseErr834 !=

		// Optional CONSTANT
		nil {
		return nil, parseErr834
	}

	if p.cur.Type == kwCONSTANT {
		decl.Constant = true
		p.advance()
	}
	var parseErr835 error

	// Type name
	decl.TypeName, parseErr835 = p.parseTypeName()
	if parseErr835 !=

		// Optional NOT NULL
		nil {
		return nil, parseErr835
	}

	if p.cur.Type == kwNOT {
		next := p.peekNext()
		if next.Type == kwNULL {
			decl.NotNull = true
			p.advance() // consume NOT
			p.advance() // consume NULL
		}
	}

	// Optional default value: := expr or DEFAULT expr
	if p.cur.Type == tokASSIGN {
		p.advance()
		var // consume :=
		parseErr836 error
		decl.Default, parseErr836 = p.parseExpr()
		if parseErr836 != nil {
			return nil, parseErr836
		}
	} else if p.cur.Type == kwDEFAULT {
		p.advance()
		var // consume DEFAULT
		parseErr837 error
		decl.Default, parseErr837 = p.parseExpr()
		if parseErr837 !=

			// Semicolon
			nil {
			return nil, parseErr837
		}
	}

	if p.cur.Type == ';' {
		p.advance()
	}

	decl.Loc.End = p.prev.End
	return decl, nil
}

// parsePLSQLCursorDecl parses a cursor declaration.
//
//	CURSOR name [(params)] IS select_stmt ;
func (p *Parser) parsePLSQLCursorDecl() (*nodes.PLSQLCursorDecl, error) {
	start := p.pos()
	p.advance() // consume CURSOR

	decl := &nodes.PLSQLCursorDecl{
		Loc: nodes.Loc{Start: start},
	}
	var parseErr838 error

	decl.Name, parseErr838 = p.parseIdentifier()
	if parseErr838 !=

		// Optional parameter list
		nil {
		return nil, parseErr838
	}

	if p.cur.Type == '(' {
		var parseErr839 error
		decl.Parameters, parseErr839 = p.parsePLSQLCursorParams()
		if parseErr839 !=

			// IS
			nil {
			return nil, parseErr839
		}
	}

	if p.cur.Type == kwIS {
		p.advance()
	}

	// SELECT statement
	if p.cur.Type == kwSELECT {
		var parseErr840 error
		decl.Query, parseErr840 = p.parseSelectStmt()
		if parseErr840 !=

			// Semicolon
			nil {
			return nil, parseErr840
		}
	}

	if p.cur.Type == ';' {
		p.advance()
	}

	decl.Loc.End = p.prev.End
	return decl, nil
}

// parsePLSQLCursorParams parses cursor parameter list.
func (p *Parser) parsePLSQLCursorParams() (*nodes.List, error) {
	params := &nodes.List{}
	p.advance() // consume (

	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		paramStart := p.pos()
		paramDecl := &nodes.PLSQLVarDecl{
			Loc: nodes.Loc{Start: paramStart},
		}
		var parseErr841 error
		paramDecl.Name, parseErr841 = p.parseIdentifier()
		if parseErr841 !=
			// Optional IN keyword (cursor params are always IN)
			nil {
			return nil, parseErr841
		}

		if p.cur.Type == kwIN {
			p.advance()
		}
		var parseErr842 error
		paramDecl.TypeName, parseErr842 = p.parseTypeName()
		if parseErr842 !=

			// Optional default
			nil {
			return nil, parseErr842
		}

		if p.cur.Type == tokASSIGN {
			p.advance()
			var parseErr843 error
			paramDecl.Default, parseErr843 = p.parseExpr()
			if parseErr843 != nil {
				return nil, parseErr843
			}
		} else if p.cur.Type == kwDEFAULT {
			p.advance()
			var parseErr844 error
			paramDecl.Default, parseErr844 = p.parseExpr()
			if parseErr844 != nil {
				return nil, parseErr844
			}
		}

		paramDecl.Loc.End = p.prev.End
		params.Items = append(params.Items, paramDecl)

		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ,
	}

	if p.cur.Type == ')' {
		p.advance()
	}
	return params, nil
}

// parsePLSQLStatements parses PL/SQL statements until END, EXCEPTION, ELSIF, ELSE, or EOF.
func (p *Parser) parsePLSQLStatements() (*nodes.List, error) {
	stmts := &nodes.List{}

	for {
		// Check for terminators
		switch p.cur.Type {
		case kwEND, kwEXCEPTION, kwELSIF, kwELSE, kwWHEN, tokEOF:
			return stmts, nil
		}

		stmt, parseErr845 := p.parsePLSQLStatement()
		if parseErr845 != nil {
			return nil, parseErr845
		}
		if stmt == nil {
			return stmts, nil
		}
		stmts.Items = append(stmts.Items, stmt)
	}
	return nil,

		// parsePLSQLStatement parses a single PL/SQL statement.
		nil
}

func (p *Parser) parsePLSQLStatement() (nodes.StmtNode, error) {
	// Handle optional label before statement
	if p.cur.Type == tokLABELOPEN {
		// This could be a labeled block or labeled loop
		return p.parsePLSQLBlock()
	}

	switch p.cur.Type {
	case kwBEGIN:
		return p.parsePLSQLBlock()

	case kwDECLARE:
		return p.parsePLSQLBlock()

	case kwIF:
		return p.parsePLSQLIf()

	case kwLOOP:
		return p.parsePLSQLBasicLoop()

	case kwWHILE:
		return p.parsePLSQLWhileLoop()

	case kwFOR:
		return p.parsePLSQLForLoop()

	case kwRETURN:
		return p.parsePLSQLReturn()

	case kwGOTO:
		return p.parsePLSQLGoto()

	case kwRAISE:
		return p.parsePLSQLRaise()

	case kwNULL:
		return p.parsePLSQLNull()

	case kwEXECUTE:
		return p.parsePLSQLExecImmediate()

	case kwOPEN:
		return p.parsePLSQLOpen()

	case kwFETCH:
		return p.parsePLSQLFetch()

	case kwCLOSE:
		return p.parsePLSQLClose()

	case kwEXIT:
		return p.parsePLSQLExit()

	case kwCONTINUE:
		return p.parsePLSQLContinue()

	case kwFORALL:
		return p.parsePLSQLForall()

	case kwPIPE:
		return p.parsePLSQLPipeRow()

	case kwCASE:
		return p.parsePLSQLCaseStmt()

	case kwMERGE:
		stmt, parseErr846 := p.parseMergeStmt()
		if parseErr846 != nil {
			return nil, parseErr846
		}
		if p.cur.Type == ';' {
			p.advance()
		}
		return stmt, nil

	// DML statements
	case kwSELECT, kwWITH:
		stmt, parseErr847 := p.parseSelectStmt()
		if parseErr847 != nil {
			return nil, parseErr847
		}
		if p.cur.Type == ';' {
			p.advance()
		}
		return stmt, nil

	case kwINSERT:
		stmt, parseErr848 := p.parseInsertStmt()
		if parseErr848 != nil {
			return nil, parseErr848
		}
		if p.cur.Type == ';' {
			p.advance()
		}
		return stmt, nil

	case kwUPDATE:
		stmt, parseErr849 := p.parseUpdateStmt()
		if parseErr849 != nil {
			return nil, parseErr849
		}
		if p.cur.Type == ';' {
			p.advance()
		}
		return stmt, nil

	case kwDELETE:
		stmt, parseErr850 := p.parseDeleteStmt()
		if parseErr850 != nil {
			return nil, parseErr850
		}
		if p.cur.Type == ';' {
			p.advance()
		}
		return stmt, nil

	default:
		// Try assignment: target := expr ;
		// tokBIND handles :NEW.col := expr in trigger bodies
		if p.isIdentLike() || p.cur.Type == tokQIDENT || p.cur.Type == tokBIND {
			return p.parsePLSQLAssignOrCall()
		}
		return nil, nil
	}
}

// parsePLSQLAssignOrCall parses an assignment statement (target := expr ;)
// or a procedure call (name(args) ;).
func (p *Parser) parsePLSQLAssignOrCall() (nodes.StmtNode, error) {
	start := p.pos()

	// Parse the target expression (could be column ref or function call)
	target, parseErr851 := p.parseExpr()
	if parseErr851 != nil {
		return nil, parseErr851

		// Check for := (assignment)
	}
	if target == nil {
		return nil, nil
	}

	if p.cur.Type == tokASSIGN {
		p.advance() // consume :=
		value, parseErr852 := p.parseExpr()
		if parseErr852 != nil {
			return nil, parseErr852
		}
		if p.cur.Type == ';' {
			p.advance()
		}
		return &nodes.PLSQLAssign{
			Target: target,
			Value:  value,
			Loc:    nodes.Loc{Start: start, End: p.prev.End},
		}, nil
	}

	// Otherwise it's a procedure call, consume the semicolon
	if p.cur.Type == ';' {
		p.advance()
	}

	return &nodes.PLSQLCall{
		Name: target,
		Loc:  nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parsePLSQLIf parses an IF/ELSIF/ELSE/END IF statement.
//
//	IF condition THEN statements [ELSIF condition THEN statements ...] [ELSE statements] END IF ;
func (p *Parser) parsePLSQLIf() (*nodes.PLSQLIf, error) {
	start := p.pos()
	p.advance() // consume IF

	ifStmt := &nodes.PLSQLIf{
		ElsIfs: &nodes.List{},
		Loc:    nodes.Loc{Start: start},
	}
	var parseErr853 error

	// Condition
	ifStmt.Condition, parseErr853 = p.parseExpr()
	if parseErr853 !=

		// THEN
		nil {
		return nil, parseErr853
	}

	if p.cur.Type == kwTHEN {
		p.advance()
	}
	var parseErr854 error

	// Then body
	ifStmt.Then, parseErr854 = p.parsePLSQLStatements()
	if parseErr854 !=

		// ELSIF clauses
		nil {
		return nil, parseErr854
	}

	for p.cur.Type == kwELSIF {
		elsifStart := p.pos()
		p.advance() // consume ELSIF

		elsif := &nodes.PLSQLElsIf{
			Loc: nodes.Loc{Start: elsifStart},
		}
		var parseErr855 error
		elsif.Condition, parseErr855 = p.parseExpr()
		if parseErr855 != nil {
			return nil, parseErr855
		}

		if p.cur.Type == kwTHEN {
			p.advance()
		}
		var parseErr856 error

		elsif.Then, parseErr856 = p.parsePLSQLStatements()
		if parseErr856 != nil {
			return nil, parseErr856
		}
		elsif.Loc.End = p.prev.End
		ifStmt.ElsIfs.Items = append(ifStmt.ElsIfs.Items, elsif)
	}

	// ELSE clause
	if p.cur.Type == kwELSE {
		p.advance()
		var parseErr857 error
		ifStmt.Else, parseErr857 = p.parsePLSQLStatements()
		if parseErr857 !=

			// END IF ;
			nil {
			return nil, parseErr857
		}
	}

	if p.cur.Type == kwEND {
		p.advance()
	}
	if p.cur.Type == kwIF {
		p.advance()
	}
	if p.cur.Type == ';' {
		p.advance()
	}

	ifStmt.Loc.End = p.prev.End
	return ifStmt, nil
}

// parsePLSQLBasicLoop parses a basic LOOP/END LOOP statement.
//
//	LOOP statements END LOOP [label] ;
func (p *Parser) parsePLSQLBasicLoop() (*nodes.PLSQLLoop, error) {
	start := p.pos()
	p.advance() // consume LOOP

	loop := &nodes.PLSQLLoop{
		Type: nodes.LOOP_BASIC,
		Loc:  nodes.Loc{Start: start},
	}
	var parseErr858 error

	loop.Statements, parseErr858 = p.parsePLSQLStatements()
	if parseErr858 !=

		// END LOOP [label] ;
		nil {
		return nil, parseErr858
	}

	if parseErr := p.consumeEndLoop(); parseErr != nil {
		return nil, parseErr
	}
	if p.cur.Type == ';' {
		p.advance()
	}

	loop.Loc.End = p.prev.End
	return loop, nil
}

// parsePLSQLWhileLoop parses a WHILE LOOP statement.
//
//	WHILE condition LOOP statements END LOOP [label] ;
func (p *Parser) parsePLSQLWhileLoop() (*nodes.PLSQLLoop, error) {
	start := p.pos()
	p.advance() // consume WHILE

	loop := &nodes.PLSQLLoop{
		Type: nodes.LOOP_WHILE,
		Loc:  nodes.Loc{Start: start},
	}
	var parseErr859 error

	loop.Condition, parseErr859 = p.parseExpr()
	if parseErr859 != nil {
		return nil, parseErr859
	}

	if p.cur.Type == kwLOOP {
		p.advance()
	}
	var parseErr860 error

	loop.Statements, parseErr860 = p.parsePLSQLStatements()
	if parseErr860 != nil {
		return nil, parseErr860
	}

	if parseErr := p.consumeEndLoop(); parseErr != nil {
		return nil, parseErr
	}
	if p.cur.Type == ';' {
		p.advance()
	}

	loop.Loc.End = p.prev.End
	return loop, nil
}

// parsePLSQLForLoop parses a FOR LOOP statement.
//
//	FOR var IN [REVERSE] lower..upper LOOP statements END LOOP [label] ;
//	FOR rec IN cursor [(args)] LOOP statements END LOOP [label] ;
func (p *Parser) parsePLSQLForLoop() (*nodes.PLSQLLoop, error) {
	start := p.pos()
	p.advance() // consume FOR

	loop := &nodes.PLSQLLoop{
		Loc: nodes.Loc{Start: start},
	}
	var parseErr861 error

	// Iterator variable
	loop.Iterator, parseErr861 = p.parseIdentifier()
	if parseErr861 !=

		// IN
		nil {
		return nil, parseErr861
	}

	if p.cur.Type == kwIN {
		p.advance()
	}

	// REVERSE (optional)
	if p.cur.Type == kwREVERSE {
		loop.Reverse = true
		p.advance()
	}

	// Determine: numeric range (lower..upper) or cursor FOR loop
	// Parse first expression
	expr1, parseErr862 := p.parseExpr()
	if parseErr862 != nil {
		return nil, parseErr862

		// Numeric FOR loop: lower..upper
	}

	if p.cur.Type == tokDOTDOT {

		loop.Type = nodes.LOOP_FOR
		loop.LowerBound = expr1
		p.advance()
		var // consume ..
		parseErr863 error
		loop.UpperBound, parseErr863 = p.parseExpr()
		if parseErr863 !=

			// Cursor FOR loop: cursor_name [(args)]
			nil {
			return nil, parseErr863
		}
	} else {

		loop.Type = nodes.LOOP_CURSOR_FOR
		// Extract cursor name from expr1
		if colRef, ok := expr1.(*nodes.ColumnRef); ok {
			loop.CursorName = colRef.Column
		}
		// Optional cursor arguments
		if p.cur.Type == '(' {
			var parseErr864 error
			loop.CursorArgs, parseErr864 = p.parsePLSQLArgList()
			if parseErr864 !=

				// LOOP
				nil {
				return nil, parseErr864
			}
		}
	}

	if p.cur.Type == kwLOOP {
		p.advance()
	}
	var parseErr865 error

	loop.Statements, parseErr865 = p.parsePLSQLStatements()
	if parseErr865 != nil {
		return nil, parseErr865
	}

	if parseErr := p.consumeEndLoop(); parseErr != nil {
		return nil, parseErr
	}
	if p.cur.Type == ';' {
		p.advance()
	}

	loop.Loc.End = p.prev.End
	return loop, nil
}

// consumeEndLoop consumes END LOOP [label].
func (p *Parser) consumeEndLoop() error {
	if p.cur.Type != kwEND {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Type != kwLOOP {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	// Optional label after END LOOP
	if p.isIdentLike() && p.cur.Type != ';' && p.cur.Type != tokEOF {
		p.advance()
	}
	return nil
}

// parsePLSQLReturn parses a RETURN [expr] ; statement.
func (p *Parser) parsePLSQLReturn() (*nodes.PLSQLReturn, error) {
	start := p.pos()
	p.advance() // consume RETURN

	ret := &nodes.PLSQLReturn{
		Loc: nodes.Loc{Start: start},
	}

	// Optional expression (not followed by ;)
	if p.cur.Type != ';' && p.cur.Type != tokEOF {
		var parseErr866 error
		ret.Expr, parseErr866 = p.parseExpr()
		if parseErr866 != nil {
			return nil, parseErr866
		}
	}

	if p.cur.Type == ';' {
		p.advance()
	}

	ret.Loc.End = p.prev.End
	return ret, nil
}

// parsePLSQLGoto parses a GOTO label ; statement.
func (p *Parser) parsePLSQLGoto() (*nodes.PLSQLGoto, error) {
	start := p.pos()
	p.advance()
	parseValue81, // consume GOTO
		parseErr82 := p.parseIdentifier()
	if parseErr82 != nil {
		return nil, parseErr82
	}
	g := &nodes.PLSQLGoto{
		Label: parseValue81,
		Loc:   nodes.Loc{Start: start},
	}

	if p.cur.Type == ';' {
		p.advance()
	}

	g.Loc.End = p.prev.End
	return g, nil
}

// parsePLSQLRaise parses a RAISE [exception_name] ; statement.
func (p *Parser) parsePLSQLRaise() (*nodes.PLSQLRaise, error) {
	start := p.pos()
	p.advance() // consume RAISE

	r := &nodes.PLSQLRaise{
		Loc: nodes.Loc{Start: start},
	}

	// Optional exception name
	if p.cur.Type != ';' && p.cur.Type != tokEOF && p.isIdentLike() {
		var parseErr867 error
		r.Exception, parseErr867 = p.parseIdentifier()
		if parseErr867 != nil {
			return nil, parseErr867
		}
	}

	if p.cur.Type == ';' {
		p.advance()
	}

	r.Loc.End = p.prev.End
	return r, nil
}

// parsePLSQLNull parses a NULL ; statement.
func (p *Parser) parsePLSQLNull() (*nodes.PLSQLNull, error) {
	start := p.pos()
	p.advance() // consume NULL

	if p.cur.Type == ';' {
		p.advance()
	}

	return &nodes.PLSQLNull{
		Loc: nodes.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parsePLSQLExecImmediate parses EXECUTE IMMEDIATE expr [INTO vars] [USING vars] ;
func (p *Parser) parsePLSQLExecImmediate() (*nodes.PLSQLExecImmediate, error) {
	start := p.pos()
	p.advance() // consume EXECUTE

	// IMMEDIATE
	if p.cur.Type == kwIMMEDIATE {
		p.advance()
	}

	stmt := &nodes.PLSQLExecImmediate{
		Loc: nodes.Loc{Start: start},
	}
	var parseErr868 error

	stmt.SQL, parseErr868 = p.parseExpr()
	if parseErr868 !=

		// INTO
		nil {
		return nil, parseErr868
	}

	if p.cur.Type == kwINTO {
		p.advance()
		var parseErr869 error
		stmt.Into, parseErr869 = p.parsePLSQLVarList()
		if parseErr869 !=

			// USING
			nil {
			return nil, parseErr869
		}
	}

	if p.cur.Type == kwUSING {
		p.advance()
		var parseErr870 error
		stmt.Using, parseErr870 = p.parsePLSQLVarList()
		if parseErr870 != nil {
			return nil, parseErr870
		}
	}

	if p.cur.Type == ';' {
		p.advance()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parsePLSQLOpen parses OPEN cursor [(args)] [FOR query] ;
func (p *Parser) parsePLSQLOpen() (*nodes.PLSQLOpen, error) {
	start := p.pos()
	p.advance() // consume OPEN

	o := &nodes.PLSQLOpen{
		Loc: nodes.Loc{Start: start},
	}
	var parseErr871 error

	o.Cursor, parseErr871 = p.parseIdentifier()
	if parseErr871 !=

		// Optional arguments
		nil {
		return nil, parseErr871
	}

	if p.cur.Type == '(' {
		var parseErr872 error
		o.Args, parseErr872 = p.parsePLSQLArgList()
		if parseErr872 !=

			// FOR query
			nil {
			return nil, parseErr872
		}
	}

	if p.cur.Type == kwFOR {
		p.advance()
		if p.cur.Type == kwSELECT || p.cur.Type == kwWITH {
			var parseErr873 error
			o.ForQuery, parseErr873 = p.parseSelectStmt()
			if parseErr873 != nil {
				return nil, parseErr873
			}
		}
	}

	if p.cur.Type == ';' {
		p.advance()
	}

	o.Loc.End = p.prev.End
	return o, nil
}

// parsePLSQLFetch parses FETCH cursor INTO vars [BULK COLLECT INTO vars] [LIMIT expr] ;
func (p *Parser) parsePLSQLFetch() (*nodes.PLSQLFetch, error) {
	start := p.pos()
	p.advance() // consume FETCH

	f := &nodes.PLSQLFetch{
		Loc: nodes.Loc{Start: start},
	}
	var parseErr874 error

	f.Cursor, parseErr874 = p.parseIdentifier()
	if parseErr874 !=

		// BULK COLLECT INTO
		nil {
		return nil, parseErr874
	}

	if p.cur.Type == kwBULK {
		f.Bulk = true
		p.advance() // consume BULK
		if p.cur.Type == kwCOLLECT {
			p.advance() // consume COLLECT
		}
	}

	// INTO
	if p.cur.Type == kwINTO {
		p.advance()
		var parseErr875 error
		f.Into, parseErr875 = p.parsePLSQLVarList()
		if parseErr875 !=

			// LIMIT
			nil {
			return nil, parseErr875
		}
	}

	if p.cur.Type == kwLIMIT {
		p.advance()
		var parseErr876 error
		f.Limit, parseErr876 = p.parseExpr()
		if parseErr876 != nil {
			return nil, parseErr876
		}
	}

	if p.cur.Type == ';' {
		p.advance()
	}

	f.Loc.End = p.prev.End
	return f, nil
}

// parsePLSQLClose parses CLOSE cursor ;
func (p *Parser) parsePLSQLClose() (*nodes.PLSQLClose, error) {
	start := p.pos()
	p.advance()
	parseValue83, // consume CLOSE
		parseErr84 := p.parseIdentifier()
	if parseErr84 != nil {
		return nil, parseErr84
	}
	c := &nodes.PLSQLClose{
		Cursor: parseValue83,
		Loc:    nodes.Loc{Start: start},
	}

	if p.cur.Type == ';' {
		p.advance()
	}

	c.Loc.End = p.prev.End
	return c, nil
}

// parsePLSQLExceptionHandlers parses exception handlers.
//
//	WHEN name [OR name ...] THEN statements
func (p *Parser) parsePLSQLExceptionHandlers() (*nodes.List, error) {
	handlers := &nodes.List{}

	for p.cur.Type == kwWHEN {
		start := p.pos()
		p.advance() // consume WHEN

		handler := &nodes.ExceptionHandler{
			Exceptions: &nodes.List{},
			Loc:        nodes.Loc{Start: start},
		}

		// Exception name(s)
		name, parseErr877 := p.parseIdentifier()
		if parseErr877 != nil {
			return nil, parseErr877
		}
		handler.Exceptions.Items = append(handler.Exceptions.Items, &nodes.String{Str: name})

		for p.cur.Type == kwOR {
			p.advance()
			var // consume OR
			parseErr878 error
			name, parseErr878 = p.parseIdentifier()
			if parseErr878 != nil {
				return nil, parseErr878
			}
			handler.Exceptions.Items = append(handler.Exceptions.Items, &nodes.String{Str: name})
		}

		// THEN
		if p.cur.Type == kwTHEN {
			p.advance()
		}
		var parseErr879 error

		handler.Statements, parseErr879 = p.parsePLSQLStatements()
		if parseErr879 != nil {
			return nil, parseErr879
		}
		handler.Loc.End = p.prev.End

		handlers.Items = append(handlers.Items, handler)
	}

	return handlers, nil
}

// parsePLSQLVarList parses a comma-separated list of variable references.
func (p *Parser) parsePLSQLVarList() (*nodes.List, error) {
	list := &nodes.List{}

	for {
		expr, parseErr880 := p.parseExpr()
		if parseErr880 != nil {
			return nil, parseErr880
		}
		if expr != nil {
			list.Items = append(list.Items, expr)
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ,
	}

	return list, nil
}

// parsePLSQLArgList parses a parenthesized argument list.
func (p *Parser) parsePLSQLArgList() (*nodes.List, error) {
	args := &nodes.List{}
	p.advance() // consume (

	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		expr, parseErr881 := p.parseExpr()
		if parseErr881 != nil {
			return nil, parseErr881
		}
		if expr != nil {
			args.Items = append(args.Items, expr)
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ,
	}

	if p.cur.Type == ')' {
		p.advance()
	}

	return args, nil
}

// parsePLSQLExit parses an EXIT [label] [WHEN condition] ; statement.
func (p *Parser) parsePLSQLExit() (nodes.StmtNode, error) {
	start := p.pos()
	p.advance() // consume EXIT

	stmt := &nodes.PLSQLExit{
		Loc: nodes.Loc{Start: start},
	}

	// Optional label (before WHEN or ;)
	if p.isIdentLike() && p.cur.Type != kwWHEN {
		stmt.Label = p.cur.Str
		p.advance()
	}

	// WHEN condition
	if p.cur.Type == kwWHEN {
		p.advance()
		var parseErr882 error
		stmt.Condition, parseErr882 = p.parseExpr()
		if parseErr882 != nil {
			return nil, parseErr882
		}
	}

	if p.cur.Type == ';' {
		p.advance()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parsePLSQLContinue parses a CONTINUE [label] [WHEN condition] ; statement.
func (p *Parser) parsePLSQLContinue() (nodes.StmtNode, error) {
	start := p.pos()
	p.advance() // consume CONTINUE

	stmt := &nodes.PLSQLContinue{
		Loc: nodes.Loc{Start: start},
	}

	// Optional label
	if p.isIdentLike() && p.cur.Type != kwWHEN {
		stmt.Label = p.cur.Str
		p.advance()
	}

	// WHEN condition
	if p.cur.Type == kwWHEN {
		p.advance()
		var parseErr883 error
		stmt.Condition, parseErr883 = p.parseExpr()
		if parseErr883 != nil {
			return nil, parseErr883
		}
	}

	if p.cur.Type == ';' {
		p.advance()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parsePLSQLForall parses a FORALL statement.
//
//	FORALL index IN lower..upper [SAVE EXCEPTIONS] dml_statement ;
func (p *Parser) parsePLSQLForall() (nodes.StmtNode, error) {
	start := p.pos()
	p.advance() // consume FORALL

	stmt := &nodes.PLSQLForall{
		Loc: nodes.Loc{Start: start},
	}

	// Index variable
	if p.isIdentLike() {
		stmt.Index = p.cur.Str
		p.advance()
	}

	// IN
	if p.cur.Type == kwIN {
		p.advance()
	}

	// lower..upper or VALUES OF or INDICES OF
	if p.isIdentLike() && (p.cur.Str == "VALUES" || p.cur.Str == "INDICES") {
		// VALUES OF / INDICES OF - skip to DML
		p.advance()
		if p.cur.Type == kwOF {
			p.advance()
		}
		var parseErr884 error
		// collection name
		stmt.Lower, parseErr884 = p.parseExpr()
		if parseErr884 != nil {
			return nil, parseErr884
		}
	} else {
		var parseErr885 error
		stmt.Lower, parseErr885 = p.parseExpr()
		if parseErr885 != nil {
			return nil, parseErr885
		}
		if p.cur.Type == tokDOTDOT {
			p.advance()
			var parseErr886 error
			stmt.Upper, parseErr886 = p.parseExpr()
			if parseErr886 !=

				// Optional SAVE EXCEPTIONS
				nil {
				return nil, parseErr886
			}
		}
	}

	if p.isIdentLike() && p.cur.Str == "SAVE" {
		p.advance()
		if p.cur.Type == kwEXCEPTION || p.isIdentLike() {
			p.advance()
		}
	}
	var parseErr887 error

	// DML statement body
	stmt.Body, parseErr887 = p.parsePLSQLStatement()
	if parseErr887 != nil {
		return nil, parseErr887
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parsePLSQLPipeRow parses a PIPE ROW statement.
//
//	PIPE ROW (expression) ;
func (p *Parser) parsePLSQLPipeRow() (nodes.StmtNode, error) {
	start := p.pos()
	p.advance() // consume PIPE

	// ROW
	if p.cur.Type == kwROW || (p.isIdentLike() && p.cur.Str == "ROW") {
		p.advance()
	}

	stmt := &nodes.PLSQLPipeRow{
		Loc: nodes.Loc{Start: start},
	}

	// (expression)
	if p.cur.Type == '(' {
		p.advance()
		var parseErr888 error
		stmt.Row, parseErr888 = p.parseExpr()
		if parseErr888 != nil {
			return nil, parseErr888
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	if p.cur.Type == ';' {
		p.advance()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parsePLSQLPragma parses a PRAGMA directive.
//
//	PRAGMA AUTONOMOUS_TRANSACTION ;
//	PRAGMA EXCEPTION_INIT ( exception, error_code ) ;
//	PRAGMA RESTRICT_REFERENCES ( ... ) ;
func (p *Parser) parsePLSQLPragma() (*nodes.PLSQLPragma, error) {
	start := p.pos()
	p.advance() // consume PRAGMA

	pragma := &nodes.PLSQLPragma{
		Loc: nodes.Loc{Start: start},
	}

	// Pragma name
	if p.isIdentLike() {
		pragma.Name = p.cur.Str
		p.advance()
	}

	// Optional arguments in parentheses
	if p.cur.Type == '(' {
		p.advance()
		pragma.Args = &nodes.List{}
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			arg, parseErr889 := p.parseExpr()
			if parseErr889 != nil {
				return nil, parseErr889
			}
			if arg != nil {
				pragma.Args.Items = append(pragma.Args.Items, arg)
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	if p.cur.Type == ';' {
		p.advance()
	}

	pragma.Loc.End = p.prev.End
	return pragma, nil
}

// parsePLSQLCaseStmt parses a PL/SQL CASE statement (distinct from CASE expression).
//
//	CASE [expr]
//	  WHEN expr THEN statements ...
//	  [ELSE statements]
//	END CASE ;
func (p *Parser) parsePLSQLCaseStmt() (nodes.StmtNode, error) {
	start := p.pos()
	p.advance() // consume CASE

	stmt := &nodes.PLSQLCase{
		Loc: nodes.Loc{Start: start},
	}

	// Optional search expression (simple CASE vs searched CASE)
	// If next is WHEN, it's a searched CASE
	if p.cur.Type != kwWHEN {
		var parseErr890 error
		stmt.Expr, parseErr890 = p.parseExpr()
		if parseErr890 !=

			// WHEN clauses
			nil {
			return nil, parseErr890
		}
	}

	for p.cur.Type == kwWHEN {
		whenStart := p.pos()
		p.advance() // consume WHEN
		when := &nodes.PLSQLWhen{
			Loc: nodes.Loc{Start: whenStart},
		}
		var parseErr891 error
		when.Expr, parseErr891 = p.parseExpr()
		if parseErr891 != nil {
			return nil, parseErr891
		}
		if p.cur.Type == kwTHEN {
			p.advance()
		}
		// Statements until next WHEN, ELSE, or END
		for p.cur.Type != kwWHEN && p.cur.Type != kwELSE && p.cur.Type != kwEND && p.cur.Type != tokEOF {
			s, parseErr892 := p.parsePLSQLStatement()
			if parseErr892 != nil {
				return nil, parseErr892
			}
			if s == nil {
				break
			}
			when.Stmts = append(when.Stmts, s)
		}
		when.Loc.End = p.prev.End
		stmt.Whens = append(stmt.Whens, when)
	}

	// ELSE
	if p.cur.Type == kwELSE {
		p.advance()
		for p.cur.Type != kwEND && p.cur.Type != tokEOF {
			s, parseErr893 := p.parsePLSQLStatement()
			if parseErr893 != nil {
				return nil, parseErr893
			}
			if s == nil {
				break
			}
			stmt.Else = append(stmt.Else, s)
		}
	}

	// END CASE ;
	if p.cur.Type == kwEND {
		p.advance()
	}
	if p.cur.Type == kwCASE {
		p.advance()
	}
	if p.cur.Type == ';' {
		p.advance()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parsePLSQLTypeDecl parses a PL/SQL TYPE declaration.
//
//	TYPE name IS TABLE OF type [INDEX BY type] ;
//	TYPE name IS VARRAY(n) OF type ;
//	TYPE name IS RECORD (field type [,...]) ;
//	TYPE name IS REF CURSOR [RETURN type] ;
func (p *Parser) parsePLSQLTypeDecl() (*nodes.PLSQLTypeDecl, error) {
	start := p.pos()
	p.advance() // consume TYPE

	decl := &nodes.PLSQLTypeDecl{
		Loc: nodes.Loc{Start: start},
	}

	// Type name
	if p.isIdentLike() {
		decl.Name = p.cur.Str
		p.advance()
	}

	// IS
	if p.cur.Type == kwIS {
		p.advance()
	}

	// TABLE OF / VARRAY / RECORD / REF CURSOR
	switch {
	case p.cur.Type == kwTABLE:
		decl.Kind = "TABLE"
		p.advance()
		if p.cur.Type == kwOF {
			p.advance()
		}
		var parseErr894 error
		decl.ElementType, parseErr894 = p.parseTypeName()
		if parseErr894 !=
			// INDEX BY
			nil {
			return nil, parseErr894
		}

		if p.cur.Type == kwINDEX {
			p.advance()
			if p.cur.Type == kwBY {
				p.advance()
			}
			var parseErr895 error
			decl.IndexBy, parseErr895 = p.parseTypeName()
			if parseErr895 != nil {
				return nil, parseErr895
			}
		}

	case p.isIdentLike() && p.cur.Str == "VARRAY" || p.isIdentLike() && p.cur.Str == "VARYING":
		decl.Kind = "VARRAY"
		p.advance()
		// Optional ARRAY keyword
		if p.isIdentLike() && p.cur.Str == "ARRAY" {
			p.advance()
		}
		// (limit)
		if p.cur.Type == '(' {
			p.advance()
			var parseErr896 error
			decl.Limit, parseErr896 = p.parseExpr()
			if parseErr896 != nil {
				return nil, parseErr896
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}
		if p.cur.Type == kwOF {
			p.advance()
		}
		var parseErr897 error
		decl.ElementType, parseErr897 = p.parseTypeName()
		if parseErr897 != nil {
			return nil, parseErr897
		}

	case p.isIdentLike() && p.cur.Str == "RECORD":
		decl.Kind = "RECORD"
		p.advance()
		// (field_name type [,...])
		if p.cur.Type == '(' {
			p.advance()
			decl.Fields = &nodes.List{}
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				field, parseErr898 := p.parsePLSQLVarDecl()
				if parseErr898 != nil {
					return nil, parseErr898
				}
				if field != nil {
					decl.Fields.Items = append(decl.Fields.Items, field)
				}
				if p.cur.Type == ',' {
					p.advance()
				} else {
					break
				}
			}
			if p.cur.Type == ')' {
				p.advance()
			}
		}

	case p.cur.Type == kwREF:
		decl.Kind = "REF_CURSOR"
		p.advance()
		if p.cur.Type == kwCURSOR {
			p.advance()
		}
		// RETURN type
		if p.cur.Type == kwRETURN {
			p.advance()
			var parseErr899 error
			decl.ReturnType, parseErr899 = p.parseTypeName()
			if parseErr899 != nil {
				return nil, parseErr899
			}
		}
	}

	if p.cur.Type == ';' {
		p.advance()
	}

	decl.Loc.End = p.prev.End
	return decl, nil
}
