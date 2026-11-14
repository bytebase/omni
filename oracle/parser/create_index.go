package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCreateIndexStmt parses a CREATE INDEX statement.
// The CREATE keyword has already been consumed. The current token is
// UNIQUE, BITMAP, MULTIVALUE, or INDEX.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-INDEX.html
//
//	CREATE [ UNIQUE ] [ BITMAP ] [ MULTIVALUE ] INDEX [ IF NOT EXISTS ]
//	    [ schema. ] index_name
//	    [ index_ilm_clause ]
//	    { cluster_index_clause | table_index_clause | bitmap_join_index_clause }
//	    [ { DEFERRED | IMMEDIATE } INVALIDATION ]
//
// A cluster index is neither UNIQUE (ORA-01715) nor BITMAP (ORA-25111), and
// a join index is a bitmap index (ORA-25959 at FROM): Oracle reports each
// at the word that makes it one. The clauses after the column list are read
// by parseCreateIndexAttributes (index_clauses.go).
func (p *Parser) parseCreateIndexStmt(start int) (*nodes.CreateIndexStmt, error) {
	stmt := &nodes.CreateIndexStmt{
		Loc: nodes.Loc{Start: start},
	}

	if p.cur.Type == kwUNIQUE {
		stmt.Unique = true
		p.advance()
	}
	if p.cur.Type == kwBITMAP {
		stmt.Bitmap = true
		p.advance()
	}
	if p.isIdentLikeStr("MULTIVALUE") {
		stmt.Multivalue = true
		p.advance()
	}
	if p.cur.Type != kwINDEX {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	// IF NOT EXISTS: IF NOT without EXISTS is ORA-11543 at IF.
	if p.cur.Type == kwIF && p.peekNext().Type == kwNOT {
		ifTok := p.advance()
		p.advance() // NOT
		if p.cur.Type != kwEXISTS {
			return nil, p.syntaxErrorAtTok(ifTok)
		}
		p.advance()
		stmt.IfNotExists = true
	}

	var err error
	stmt.Name, err = p.parseReservedCheckedObjectName()
	if err != nil {
		return nil, err
	}
	if stmt.Name == nil || stmt.Name.Name == "" {
		return nil, p.syntaxErrorAtCur()
	}

	if p.cur.Type != kwON {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	// cluster_index_clause: ON CLUSTER [ schema. ] cluster [ index_attributes ]
	if p.cur.Type == kwCLUSTER {
		if stmt.Unique || stmt.Bitmap {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		stmt.Cluster, err = p.parseObjectName()
		if err != nil {
			return nil, err
		}
		if stmt.Cluster == nil || stmt.Cluster.Name == "" {
			return nil, p.syntaxErrorAtCur()
		}
		if err := p.parseCreateIndexAttributes(stmt); err != nil {
			return nil, err
		}
		if err := p.parseCreateIndexInvalidation(stmt); err != nil {
			return nil, err
		}
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	// table_index_clause or bitmap_join_index_clause: ON [ schema. ] table [ t_alias ]
	stmt.Table, err = p.parseObjectName()
	if err != nil {
		return nil, err
	}
	if stmt.Table == nil || stmt.Table.Name == "" {
		return nil, p.syntaxErrorAtCur()
	}
	if p.cur.Type != '(' && p.isIdentLike() {
		stmt.Alias, err = p.parseIdentifier()
		if err != nil {
			return nil, err
		}
	}

	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	stmt.Columns, err = p.parseIndexColumnList()
	if err != nil {
		return nil, err
	}
	if stmt.Columns == nil || stmt.Columns.Len() == 0 {
		return nil, p.syntaxErrorAtCur()
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	// bitmap_join_index_clause: FROM [ schema. ] table [, ...] WHERE join_condition
	if p.cur.Type == kwFROM {
		if !stmt.Bitmap {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		stmt.FromTables = &nodes.List{}
		for {
			tbl, err := p.parseRequiredObjectName()
			if err != nil {
				return nil, err
			}
			stmt.FromTables.Items = append(stmt.FromTables.Items, tbl)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if p.cur.Type != kwWHERE {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		stmt.Where, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
		if stmt.Where == nil {
			return nil, p.syntaxErrorAtCur()
		}
	}
	if err := p.parseCreateIndexAttributes(stmt); err != nil {
		return nil, err
	}
	if err := p.parseCreateIndexInvalidation(stmt); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseCreateIndexInvalidation parses optional { DEFERRED | IMMEDIATE } INVALIDATION.
func (p *Parser) parseCreateIndexInvalidation(stmt *nodes.CreateIndexStmt) error {
	if !p.isIxInvalidationStart() {
		return nil
	}
	inv, err := p.ixInvalidation()
	if err != nil {
		return err
	}
	stmt.Invalidation = inv
	return nil
}

func appendCreateIndexOption(stmt *nodes.CreateIndexStmt, key, value string, loc nodes.Loc) {
	if loc.End <= loc.Start {
		return
	}
	if stmt.Options == nil {
		stmt.Options = &nodes.List{}
	}
	stmt.Options.Items = append(stmt.Options.Items, &nodes.DDLOption{Key: key, Value: value, Loc: loc})
}

// collectParenthesizedBlock consumes and returns a balanced parenthesized block if present.
func (p *Parser) collectParenthesizedBlock() string {
	if p.cur.Type != '(' {
		return ""
	}
	var tokens []string
	depth := 1
	tokens = append(tokens, p.ddlOptionTokenText(p.cur))
	p.advance()
	for depth > 0 && p.cur.Type != tokEOF {
		tokens = append(tokens, p.ddlOptionTokenText(p.cur))
		if p.cur.Type == '(' {
			depth++
		} else if p.cur.Type == ')' {
			depth--
		}
		p.advance()
	}
	return strings.Join(tokens, " ")
}

// parseIndexColumnList parses a comma-separated list of index columns.
func (p *Parser) parseIndexColumnList() (*nodes.List, error) {
	list := &nodes.List{}
	for {
		col, parseErr424 := p.parseIndexColumn()
		if parseErr424 != nil {
			return nil, parseErr424
		}
		if col == nil {
			break
		}
		list.Items = append(list.Items, col)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
		if p.cur.Type == ')' || p.cur.Type == tokEOF {
			return nil, p.syntaxErrorAtCur()
		}
	}
	return list, nil
}

// parseIndexColumn parses a single index column expression with optional ASC/DESC.
//
//	index_expr [ ASC | DESC ]
func (p *Parser) parseIndexColumn() (*nodes.IndexColumn, error) {
	start := p.pos()
	expr, parseErr425 := p.parseExpr()
	if parseErr425 != nil {
		return nil, parseErr425
	}
	if expr == nil {
		return nil, nil
	}
	// A column name cut short at its dot (a.) is ORA-00931 at what follows.
	if p.prev.Type == '.' {
		return nil, p.syntaxErrorAtCur()
	}

	col := &nodes.IndexColumn{
		Expr: expr,
		Loc:  nodes.Loc{Start: start},
	}

	// ASC | DESC
	switch p.cur.Type {
	case kwASC:
		col.Dir = nodes.SORTBY_ASC
		p.advance()
	case kwDESC:
		col.Dir = nodes.SORTBY_DESC
		p.advance()
	}

	// NULLS FIRST | NULLS LAST
	if p.cur.Type == kwNULLS {
		p.advance()
		switch p.cur.Type {
		case kwFIRST:
			col.NullOrder = nodes.SORTBY_NULLS_FIRST
			p.advance()
		case kwLAST:
			col.NullOrder = nodes.SORTBY_NULLS_LAST
			p.advance()
		}
	}

	col.Loc.End = p.prev.End
	return col, nil
}

// ---------------------------------------------------------------------------
// Operator parameter types
// ---------------------------------------------------------------------------

// parseOperatorType reads a parameter or return type of an operator binding
// or an indextype's operator list. Oracle takes a type name there, never
// its size or precision: BINDING (NUMBER(1)) is ORA-00907 and RETURN
// NUMBER(1) ORA-29827 at the '(', TIMESTAMP WITH TIME ZONE ORA-00907 at
// WITH. LONG RAW and DOUBLE PRECISION are the two-word names it takes.
func (p *Parser) parseOperatorType() (string, error) {
	if !p.isIdentLike() {
		return "", p.syntaxErrorAtCur()
	}
	name := p.advance().Str
	switch {
	case name == "LONG" && p.isKeywordStr("RAW"):
		name += " " + p.advance().Str
	case name == "DOUBLE" && p.isKeywordStr("PRECISION"):
		name += " " + p.advance().Str
	case p.cur.Type == '.':
		p.advance()
		if !p.isIdentLike() {
			return "", p.syntaxErrorAtCur()
		}
		name += "." + p.advance().Str
	}
	if p.cur.Type == '(' {
		return "", p.syntaxErrorAtCur()
	}
	return name, nil
}

// parseOperatorTypeList reads ( type [, type ]... ). The list may not be
// empty: BINDING () is ORA-00902 at the ')'.
func (p *Parser) parseOperatorTypeList() ([]string, error) {
	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	var types []string
	for {
		t, err := p.parseOperatorType()
		if err != nil {
			return nil, err
		}
		types = append(types, t)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	return types, nil
}

// parseOperatorTypeName reads [ schema. ] name for the operators, types and
// functions of these statements.
func (p *Parser) parseOperatorTypeName() (*nodes.ObjectName, error) {
	return p.ixObjectName()
}

// parseSharingClause reads SHARING = { METADATA | NONE } after SHARING.
func (p *Parser) parseSharingClause() (string, error) {
	p.advance() // SHARING
	if p.cur.Type != '=' {
		return "", p.syntaxErrorAtCur()
	}
	p.advance()
	if !p.isKeywordStr("METADATA") && !p.isKeywordStr("NONE") {
		return "", p.syntaxErrorAtCur()
	}
	return p.advance().Str, nil
}

// ---------------------------------------------------------------------------
// CREATE INDEXTYPE
// ---------------------------------------------------------------------------

// parseCreateIndextypeStmt parses a CREATE INDEXTYPE statement.
// Called after CREATE [OR REPLACE] has been consumed and INDEXTYPE consumed.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-INDEXTYPE.html
//
//	CREATE [ OR REPLACE ] INDEXTYPE [ IF NOT EXISTS ] [ schema. ] indextype
//	    [ SHARING = { METADATA | NONE } ]
//	    FOR [ schema. ] operator ( parameter_type [, parameter_type ]... )
//	        [, [ schema. ] operator ( parameter_type [, parameter_type ]... ) ]...
//	    using_type_clause
//	    [ WITH LOCAL [ RANGE ] PARTITION ]
//	    [ storage_table_clause ]
//	    [ array_DML_clause ]
//
// The statement may end after the name or after USING (see ixEnd).
func (p *Parser) parseCreateIndextypeStmt(start int, orReplace bool) (*nodes.CreateIndextypeStmt, error) {
	stmt := &nodes.CreateIndextypeStmt{
		OrReplace: orReplace,
		Loc:       nodes.Loc{Start: start},
	}

	// IF NOT EXISTS: IF NOT without EXISTS is ORA-11543 at IF.
	if p.cur.Type == kwIF && p.peekNext().Type == kwNOT {
		ifTok := p.advance()
		p.advance() // NOT
		if p.cur.Type != kwEXISTS {
			return nil, p.syntaxErrorAtTok(ifTok)
		}
		p.advance()
		stmt.IfNotExists = true
	}

	var err error
	stmt.Name, err = p.parseOperatorTypeName()
	if err != nil {
		return nil, err
	}
	if p.ixEnd() {
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	if p.isKeywordStr("SHARING") {
		if stmt.Sharing, err = p.parseSharingClause(); err != nil {
			return nil, err
		}
	}

	if p.cur.Type != kwFOR {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	for {
		op, err := p.parseIndextypeOp()
		if err != nil {
			return nil, err
		}
		stmt.Operators = append(stmt.Operators, op)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}

	if p.cur.Type != kwUSING {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if p.ixEnd() {
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}
	if stmt.UsingType, err = p.parseOperatorTypeName(); err != nil {
		return nil, err
	}

	seen := map[string]bool{}
	for p.cur.Type == kwWITH {
		p.advance()
		switch {
		case p.cur.Type == kwLOCAL:
			if err := p.ixOnce(seen, "LOCAL"); err != nil {
				return nil, err
			}
			p.advance()
			if p.isKeywordStr("RANGE") {
				stmt.WithRange = true
				p.advance()
			}
			if p.cur.Type != kwPARTITION {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			stmt.WithLocal = true
		case p.isKeywordStr("SYSTEM") || p.cur.Type == kwUSER:
			if err := p.ixOnce(seen, "STORAGE"); err != nil {
				return nil, err
			}
			if stmt.StorageTable, err = p.parseStorageTableTail(); err != nil {
				return nil, err
			}
		case p.isKeywordStr("ARRAY"):
			if err := p.ixOnce(seen, "ARRAY"); err != nil {
				return nil, err
			}
			if err := p.parseArrayDMLTail(); err != nil {
				return nil, err
			}
			stmt.ArrayDML = true
		default:
			return nil, p.syntaxErrorAtCur()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseStorageTableTail reads storage_table_clause after WITH:
//
//	{ SYSTEM | USER } MANAGED STORAGE TABLES
func (p *Parser) parseStorageTableTail() (string, error) {
	kind := p.advance().Str // SYSTEM or USER
	for _, w := range []string{"MANAGED", "STORAGE", "TABLES"} {
		if err := p.ixWord(w); err != nil {
			return "", err
		}
	}
	return kind, nil
}

// parseArrayDMLTail reads array_DML_clause after WITH:
//
//	ARRAY DML [ ( [ schema. ] type , [ schema. ] varray_type ) ]
//
// The page allows a list of pairs; Oracle takes one (WITH ARRAY DML (a, a),
// (a, a) is refused at the second).
func (p *Parser) parseArrayDMLTail() error {
	p.advance() // ARRAY
	if err := p.ixWord("DML"); err != nil {
		return err
	}
	if p.cur.Type != '(' {
		return nil
	}
	p.advance()
	for i := 0; i < 2; i++ {
		if i > 0 {
			if p.cur.Type != ',' {
				return p.syntaxErrorAtCur()
			}
			p.advance()
		}
		if name, err := p.parseOperatorType(); err != nil {
			return err
		} else if name == "" {
			return p.syntaxErrorAtCur()
		}
	}
	if p.cur.Type != ')' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// parseIndextypeOp reads [ schema. ] operator ( parameter_type [, ...] ).
func (p *Parser) parseIndextypeOp() (*nodes.IndextypeOp, error) {
	start := p.pos()
	name, err := p.parseOperatorTypeName()
	if err != nil {
		return nil, err
	}
	op := &nodes.IndextypeOp{Name: name, Loc: nodes.Loc{Start: start}}
	if op.ParamTypes, err = p.parseOperatorTypeList(); err != nil {
		return nil, err
	}
	op.Loc.End = p.prev.End
	return op, nil
}

// ---------------------------------------------------------------------------
// ALTER INDEXTYPE
// ---------------------------------------------------------------------------

// parseAlterIndextypeStmt parses an ALTER INDEXTYPE statement.
// Called after ALTER has been consumed and INDEXTYPE consumed.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-INDEXTYPE.html
//
//	ALTER INDEXTYPE [ IF EXISTS ] [ schema. ] indextype
//	    { { ADD | DROP } [ schema. ] operator ( parameter_type [, parameter_type ]... )
//	        [ , { ADD | DROP } [ schema. ] operator ( parameter_type [, parameter_type ]... ) ]...
//	        [ using_type_clause ]
//	    | COMPILE
//	    }
//	    [ WITH LOCAL PARTITION ]
//	    [ storage_table_clause ]
//
// Oracle takes ADD a(NUMBER) ADD a(NUMBER) without the comma, and nothing
// after COMPILE (COMPILE WITH LOCAL PARTITION is ORA-00922 at WITH).
func (p *Parser) parseAlterIndextypeStmt(start int) (*nodes.AlterIndextypeStmt, error) {
	stmt := &nodes.AlterIndextypeStmt{
		Loc: nodes.Loc{Start: start},
	}
	if p.cur.Type == kwIF && p.peekNext().Type == kwEXISTS {
		stmt.IfExists = true
		p.advance()
		p.advance()
	}
	var err error
	stmt.Name, err = p.parseOperatorTypeName()
	if err != nil {
		return nil, err
	}

	switch {
	case p.isKeywordStr("COMPILE"):
		stmt.Action = "COMPILE"
		p.advance()
		stmt.Loc.End = p.prev.End
		return stmt, nil
	case p.cur.Type == kwADD || p.cur.Type == kwDROP:
		stmt.Action = "ADD_DROP"
	default:
		return nil, p.syntaxErrorAtCur()
	}

	for p.cur.Type == kwADD || p.cur.Type == kwDROP {
		modStart := p.pos()
		isAdd := p.advance().Type == kwADD
		op, err := p.parseIndextypeOp()
		if err != nil {
			return nil, err
		}
		stmt.Modifications = append(stmt.Modifications, &nodes.IndextypeModOp{
			Add:        isAdd,
			Name:       op.Name,
			ParamTypes: op.ParamTypes,
			Loc:        nodes.Loc{Start: modStart, End: p.prev.End},
		})
		if p.cur.Type == ',' {
			p.advance()
			if p.cur.Type != kwADD && p.cur.Type != kwDROP {
				return nil, p.syntaxErrorAtCur()
			}
		}
	}

	// using_type_clause: USING [ schema. ] implementation_type [ array_DML_clause ]
	if p.cur.Type == kwUSING {
		p.advance()
		if p.ixEnd() {
			stmt.Loc.End = p.prev.End
			return stmt, nil
		}
		if stmt.UsingType, err = p.parseOperatorTypeName(); err != nil {
			return nil, err
		}
		if p.cur.Type == kwWITH && p.peekNext().Str == "ARRAY" {
			p.advance()
			if err := p.parseArrayDMLTail(); err != nil {
				return nil, err
			}
			stmt.ArrayDML = true
		}
	}

	seen := map[string]bool{}
	for p.cur.Type == kwWITH {
		p.advance()
		switch {
		case p.cur.Type == kwLOCAL:
			if err := p.ixOnce(seen, "LOCAL"); err != nil {
				return nil, err
			}
			p.advance()
			if p.cur.Type != kwPARTITION {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			stmt.WithLocal = true
		case p.isKeywordStr("SYSTEM") || p.cur.Type == kwUSER:
			if err := p.ixOnce(seen, "STORAGE"); err != nil {
				return nil, err
			}
			if stmt.StorageTable, err = p.parseStorageTableTail(); err != nil {
				return nil, err
			}
		default:
			return nil, p.syntaxErrorAtCur()
		}
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// CREATE OPERATOR
// ---------------------------------------------------------------------------

// parseCreateOperatorStmt parses a CREATE OPERATOR statement.
// Called after CREATE [OR REPLACE] has been consumed and OPERATOR consumed.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-OPERATOR.html
//
//	CREATE [ OR REPLACE | IF NOT EXISTS ] OPERATOR
//	    [ schema. ] operator
//	    binding_clause
//	    [ SHARING = { METADATA | NONE } ]
//
// Where Oracle answers otherwise than the page: IF NOT EXISTS follows
// OPERATOR (CREATE IF NOT EXISTS OPERATOR is ORA-00901 at IF), SHARING
// follows the name (after the binding it is ORA-00922), and a further
// binding is , ( parameter_type [, ...] ) RETURN ... without BINDING.
func (p *Parser) parseCreateOperatorStmt(start int, orReplace bool, ifNotExists bool) (*nodes.CreateOperatorStmt, error) {
	if ifNotExists {
		return nil, p.syntaxErrorAtCur()
	}
	stmt := &nodes.CreateOperatorStmt{
		OrReplace: orReplace,
		Loc:       nodes.Loc{Start: start},
	}
	if p.cur.Type == kwIF && p.peekNext().Type == kwNOT {
		ifTok := p.advance()
		p.advance() // NOT
		if p.cur.Type != kwEXISTS {
			return nil, p.syntaxErrorAtTok(ifTok)
		}
		p.advance()
		stmt.IfNotExists = true
	}

	var err error
	stmt.Name, err = p.parseOperatorTypeName()
	if err != nil {
		return nil, err
	}

	if p.isKeywordStr("SHARING") {
		if stmt.Sharing, err = p.parseSharingClause(); err != nil {
			return nil, err
		}
	}

	if !p.isKeywordStr("BINDING") {
		return nil, p.syntaxErrorAtCur()
	}
	bindingStart := p.pos()
	p.advance()
	for {
		binding, err := p.parseOperatorBinding(false, bindingStart)
		if err != nil {
			return nil, err
		}
		stmt.Bindings = append(stmt.Bindings, binding)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
		bindingStart = p.pos()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseOperatorBinding reads a binding after BINDING (or after the comma
// that starts a further one):
//
//	( parameter_type [, parameter_type ]... )
//	    RETURN return_type
//	    [ implementation_clause ]
//	    using_function_clause
//
//	implementation_clause:
//	    { ANCILLARY TO primary_operator ( parameter_type [, ...] ) | context_clause }
//
//	context_clause:
//	    [ WITH INDEX CONTEXT , SCAN CONTEXT implementation_type ]
//	    [ COMPUTE ANCILLARY DATA ] [ WITH COLUMN CONTEXT ]
//
//	using_function_clause:
//	    USING [ schema. ] [ package_name. | type_name. ] function_name
//
// With alter set (ALTER OPERATOR's page) WITH INDEX CONTEXT must lead the
// context_clause; CREATE OPERATOR takes COMPUTE ANCILLARY DATA only after
// it too (ORA-29930 after DATA).
func (p *Parser) parseOperatorBinding(alter bool, start int) (*nodes.OperatorBinding, error) {
	binding := &nodes.OperatorBinding{Loc: nodes.Loc{Start: start}}
	var err error
	if binding.ParamTypes, err = p.parseOperatorTypeList(); err != nil {
		return nil, err
	}
	if p.cur.Type != kwRETURN {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if binding.ReturnType, err = p.parseOperatorType(); err != nil {
		return nil, err
	}

	if p.isKeywordStr("ANCILLARY") {
		p.advance()
		if p.cur.Type != kwTO {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if binding.AncillaryTo, err = p.parseOperatorTypeName(); err != nil {
			return nil, err
		}
		if binding.AncillaryParams, err = p.parseOperatorTypeList(); err != nil {
			return nil, err
		}
	} else {
		if p.cur.Type == kwWITH && p.peekNext().Type == kwINDEX {
			p.advance() // WITH
			p.advance() // INDEX
			if err := p.ixWord("CONTEXT"); err != nil {
				return nil, err
			}
			if p.cur.Type != ',' {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			if err := p.ixWord("SCAN"); err != nil {
				return nil, err
			}
			if err := p.ixWord("CONTEXT"); err != nil {
				return nil, err
			}
			scan, err := p.parseOperatorTypeName()
			if err != nil {
				return nil, err
			}
			binding.WithIndexCtx = true
			binding.ScanCtxType = scan.Name
			if scan.Schema != "" {
				binding.ScanCtxType = scan.Schema + "." + scan.Name
			}
		}
		if p.isKeywordStr("COMPUTE") {
			if alter && !binding.WithIndexCtx {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			if err := p.ixWord("ANCILLARY"); err != nil {
				return nil, err
			}
			if err := p.ixWord("DATA"); err != nil {
				return nil, err
			}
			if !binding.WithIndexCtx {
				return nil, p.syntaxErrorAtCur()
			}
			binding.ComputeAnc = true
		}
		if p.cur.Type == kwWITH && p.peekNext().Str == "COLUMN" {
			if alter && !binding.WithIndexCtx {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance() // WITH
			p.advance() // COLUMN
			if err := p.ixWord("CONTEXT"); err != nil {
				return nil, err
			}
			binding.WithColumnCtx = true
		}
	}

	if p.cur.Type != kwUSING {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if binding.UsingFunc, err = p.parseOperatorFunctionName(); err != nil {
		return nil, err
	}
	binding.Loc.End = p.prev.End
	return binding, nil
}

// parseOperatorFunctionName reads [ schema. ] [ package_name. | type_name. ]
// function_name: at most three parts (a.b.c.d is ORA-00922 at the third
// dot).
func (p *Parser) parseOperatorFunctionName() (*nodes.ObjectName, error) {
	obj := &nodes.ObjectName{Loc: nodes.Loc{Start: p.pos()}}
	var parts []string
	for {
		name, err := p.ixName()
		if err != nil {
			return nil, err
		}
		parts = append(parts, name)
		if p.cur.Type != '.' || len(parts) == 3 {
			break
		}
		p.advance()
	}
	obj.Name = parts[len(parts)-1]
	if len(parts) > 1 {
		obj.Schema = strings.Join(parts[:len(parts)-1], ".")
	}
	obj.Loc.End = p.prev.End
	return obj, nil
}

// ---------------------------------------------------------------------------
// ALTER OPERATOR
// ---------------------------------------------------------------------------

// parseAlterOperatorStmt parses an ALTER OPERATOR statement.
// Called after ALTER has been consumed and OPERATOR consumed.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-OPERATOR.html
//
//	ALTER OPERATOR [ IF EXISTS ] [ schema. ] operator
//	    { add_binding_clause
//	    | drop_binding_clause
//	    | COMPILE
//	    }
//
//	add_binding_clause:  ADD BINDING ( datatype [, datatype ]... ) RETURN ( datatype ) ...
//	drop_binding_clause: DROP BINDING ( datatype [, datatype ]... ) [ FORCE ]
//
// Oracle reads RETURN datatype without the page's parentheses: RETURN
// (NUMBER) is ORA-00902 at the '('.
func (p *Parser) parseAlterOperatorStmt(start int) (*nodes.AlterOperatorStmt, error) {
	stmt := &nodes.AlterOperatorStmt{
		Loc: nodes.Loc{Start: start},
	}
	if p.cur.Type == kwIF && p.peekNext().Type == kwEXISTS {
		stmt.IfExists = true
		p.advance()
		p.advance()
	}
	var err error
	stmt.Name, err = p.parseOperatorTypeName()
	if err != nil {
		return nil, err
	}

	switch {
	case p.isKeywordStr("COMPILE"):
		stmt.Action = "COMPILE"
		p.advance()
	case p.cur.Type == kwADD:
		stmt.Action = "ADD_BINDING"
		p.advance()
		bindingStart := p.pos()
		if err := p.ixWord("BINDING"); err != nil {
			return nil, err
		}
		if stmt.Binding, err = p.parseOperatorBinding(true, bindingStart); err != nil {
			return nil, err
		}
	case p.cur.Type == kwDROP:
		stmt.Action = "DROP_BINDING"
		p.advance()
		if err := p.ixWord("BINDING"); err != nil {
			return nil, err
		}
		if stmt.DropTypes, err = p.parseOperatorTypeList(); err != nil {
			return nil, err
		}
		if p.cur.Type == kwFORCE {
			stmt.DropForce = true
			p.advance()
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}
