// Package parser - declare_set.go implements T-SQL DECLARE and SET statement parsing.
package parser

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// setOptions defines the valid option names for SET <option> ON/OFF/value statements.
// This combines:
//   - Predicate SET options (from SqlScriptDOM PredicateSetOptionsHelper — 25 options)
//   - General SET commands (from SqlScriptDOM GeneralSetCommandTypeHelper — 7 options)
//   - SET ROWCOUNT / SET TEXTSIZE / SET STATISTICS (keyword-based)
var setOptions = newOptionSet(
	// Registered keywords that are valid SET options:
	kwNOCOUNT,    // SET NOCOUNT ON/OFF
	kwXACT_ABORT, // SET XACT_ABORT ON/OFF
	kwROWCOUNT,   // SET ROWCOUNT n
	kwTEXTSIZE,   // SET TEXTSIZE n
	kwSTATISTICS, // SET STATISTICS IO/TIME/PROFILE/XML ON/OFF
	kwLANGUAGE,   // SET LANGUAGE ...
).withIdents(
	// Predicate SET options (scanned as identifiers):
	"QUOTED_IDENTIFIER",
	"CONCAT_NULL_YIELDS_NULL",
	"CURSOR_CLOSE_ON_COMMIT",
	"ARITHABORT",
	"ARITHIGNORE",
	"FMTONLY",
	"NOEXEC",
	"NUMERIC_ROUNDABORT",
	"PARSEONLY",
	"ANSI_DEFAULTS",
	"ANSI_NULL_DFLT_OFF",
	"ANSI_NULL_DFLT_ON",
	"ANSI_NULLS",
	"ANSI_PADDING",
	"ANSI_WARNINGS",
	"FORCEPLAN",
	"SHOWPLAN_ALL",
	"SHOWPLAN_TEXT",
	"SHOWPLAN_XML",
	"IMPLICIT_TRANSACTIONS",
	"REMOTE_PROC_TRANSACTIONS",
	"DISABLE_DEF_CNST_CHK",
	"NO_BROWSETABLE",
	// General SET commands (scanned as identifiers):
	"DATEFORMAT",
	"DATEFIRST",
	"DEADLOCK_PRIORITY",
	"LOCK_TIMEOUT",
	"QUERY_GOVERNOR_COST_LIMIT",
	"CONTEXT_INFO",
)

// setOnOffListOptions are the setOptions that may share one ON or OFF with
// others in a list: SET ANSI_NULLS, QUOTED_IDENTIFIER OFF. SQL Server 2022
// rejects SHOWPLAN_ALL, SHOWPLAN_TEXT, and SHOWPLAN_XML in a list ("The SET
// SHOWPLAN statements must be the only statements in the batch") and does not
// recognize DISABLE_DEF_CNST_CHK, so neither is listed here.
var setOnOffListOptions = newOptionSet(
	kwNOCOUNT,
	kwXACT_ABORT,
).withIdents(
	"QUOTED_IDENTIFIER",
	"CONCAT_NULL_YIELDS_NULL",
	"CURSOR_CLOSE_ON_COMMIT",
	"ARITHABORT",
	"ARITHIGNORE",
	"FMTONLY",
	"NOEXEC",
	"NUMERIC_ROUNDABORT",
	"PARSEONLY",
	"ANSI_DEFAULTS",
	"ANSI_NULL_DFLT_OFF",
	"ANSI_NULL_DFLT_ON",
	"ANSI_NULLS",
	"ANSI_PADDING",
	"ANSI_WARNINGS",
	"FORCEPLAN",
	"IMPLICIT_TRANSACTIONS",
	"REMOTE_PROC_TRANSACTIONS",
	"NO_BROWSETABLE",
)

// setStatisticsOptions are the options of SET STATISTICS, which may also be
// listed: SET STATISTICS IO, TIME ON.
var setStatisticsOptions = newOptionSet(
	kwTIME,
	kwXML,
).withIdents(
	"IO",
	"PROFILE",
)

// parseDeclareStmt parses a DECLARE statement.
//
// BNF: mssql/parser/bnf/declare-local-variable-transact-sql.bnf
//
//	DECLARE
//	{
//	  { @local_variable [AS] data_type [ = value ] }
//	  | { @cursor_variable_name CURSOR }
//	  | { @table_variable_name [AS] TABLE ( { <column_definition> | <table_constraint> | <table_index> } [ , ...n ] ) }
//	} [ , ...n ]
//
// declareStatements (TSql160.g:16366): a table variable is declared alone
// (declareTableBody, TSql160.g:16697), the others in a list where each one
// is required after a comma, and SQL Server refuses a name declared twice
// (Msg 134; names compare without case).
func (p *Parser) parseDeclareStmt() (*nodes.DeclareStmt, error) {
	loc := p.pos()
	p.advance() // consume DECLARE

	// Completion: after DECLARE → @ variable
	if p.collectMode() {
		p.addRuleCandidate("@variable")
		return nil, errCollecting
	}

	stmt := &nodes.DeclareStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	var vars []nodes.Node
	seen := map[string]bool{}
	for {
		vd, err := p.parseVariableDecl()
		if err != nil {
			return nil, err
		}
		if vd == nil {
			if p.collectMode() {
				break
			}
			return nil, p.syntaxErrorAtCur()
		}
		key := strings.ToLower(vd.Name)
		if seen[key] {
			return nil, p.newParseError(vd.Loc.Start, fmt.Sprintf("the variable name '%s' has already been declared", vd.Name))
		}
		seen[key] = true
		vars = append(vars, vd)
		if vd.IsTable {
			break
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}
	stmt.Variables = &nodes.List{Items: vars}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseVariableDecl parses a single variable declaration.
//
// BNF: mssql/parser/bnf/declare-local-variable-transact-sql.bnf
//
//	@local_variable [AS] data_type [ = value ]
//	| @cursor_variable_name CURSOR
//	| @table_variable_name [AS] TABLE ( { <column_definition> | <table_constraint> | <table_index> } [ , ...n ] )
func (p *Parser) parseVariableDecl() (*nodes.VariableDecl, error) {
	if p.cur.Type != tokVARIABLE {
		return nil, nil
	}

	loc := p.pos()
	vd := &nodes.VariableDecl{
		Name: p.cur.Str,
		Loc:  nodes.Loc{Start: loc, End: -1},
	}
	p.advance() // consume @var

	// Completion: after DECLARE @v → type_name
	if p.collectMode() {
		p.addRuleCandidate("type_name")
		return nil, errCollecting
	}

	// Optional AS keyword
	p.match(kwAS)

	// TABLE type
	if p.cur.Type == kwTABLE {
		p.advance()
		vd.IsTable = true
		if p.cur.Type == '(' {
			p.advance()
			var cols []nodes.Node
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				// Near-CREATE TABLE element grammar: a table variable body
				// takes column definitions, UNNAMED table constraints
				// (PRIMARY KEY/UNIQUE/CHECK — the CONSTRAINT name prefix is
				// rejected by the engine in this context, verified on SQL
				// Server 2019: Msg 156), and indexes.
				if p.cur.Type == kwPRIMARY ||
					p.cur.Type == kwUNIQUE || p.cur.Type == kwCHECK {
					constraint, err := p.parseTableConstraint()
					if err != nil {
						return nil, err
					}
					if constraint != nil {
						cols = append(cols, constraint)
					}
				} else if p.cur.Type == kwINDEX {
					idx, err := p.parseInlineTableIndex()
					if err != nil {
						return nil, err
					}
					if idx != nil {
						cols = append(cols, idx)
					}
				} else {
					col, err := p.parseTypedColumnDef()
					if err != nil {
						return nil, err
					}
					if col != nil {
						cols = append(cols, col)
					}
				}
				if _, ok := p.match(','); !ok {
					break
				}
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			// The body holds an element (declareTableBodyMain): SQL Server
			// refuses TABLE () (Msg 102).
			if len(cols) == 0 {
				return nil, p.transSyntaxErrorAtPrev()
			}
			vd.TableDef = &nodes.List{Items: cols}
		} else {
			return nil, p.syntaxErrorAtCur()
		}
		vd.Loc.End = p.prevEnd()
		return vd, nil
	}

	// CURSOR type
	if p.cur.Type == kwCURSOR {
		p.advance()
		vd.IsCursor = true
		vd.Loc.End = p.prevEnd()
		return vd, nil
	}

	// Data type
	dt, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	if dt == nil {
		return nil, p.newParseError(p.cur.Loc, "expected data type in variable declaration")
	}
	vd.DataType = dt

	// Optional default value, which = requires (declareVariableElement,
	// TSql160.g:16335).
	if p.cur.Type == '=' {
		p.advance()
		var err error
		vd.Default, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
		if vd.Default == nil {
			return nil, p.syntaxErrorAtCur()
		}
	}

	vd.Loc.End = p.prevEnd()
	return vd, nil
}

// parseSetStmt parses a SET statement.
//
// BNF: mssql/parser/bnf/set-transact-sql.bnf
//
//	SET @local_variable { = | += | -= | *= | /= | %= | &= | ^= | |= } expression
//	SET @local_variable.property_name = expression
//	SET @local_variable.method_name ( [ argument [ ,...n ] ] )
//	SET @local_variable = CURSOR [ FORWARD_ONLY | SCROLL ] ... FOR select_statement
//
// BNF: mssql/parser/bnf/set-transaction-isolation-level-transact-sql.bnf
//
//	SET TRANSACTION ISOLATION LEVEL
//	    { READ UNCOMMITTED | READ COMMITTED | REPEATABLE READ | SNAPSHOT | SERIALIZABLE }
//
// SET session options:
//
//	SET { option_name } { ON | OFF | value }
//	SET option_name [ , ...n ] { ON | OFF }
//	SET IDENTITY_INSERT table_name { ON | OFF }
//	SET STATISTICS { IO | TIME | PROFILE | XML } [ , ...n ] { ON | OFF }
//
// setVariableStatement (TSql160.g:16406): the assignment, a property or
// method after a dot, or a cursor definition is required. SQL Server
// refuses the :: form there (Msg 102).
func (p *Parser) parseSetStmt() (nodes.StmtNode, error) {
	loc := p.pos()
	p.advance() // consume SET

	// Completion: after SET → @variable, session option keywords
	if p.collectMode() {
		p.addRuleCandidate("@variable")
		p.addTokenCandidate(kwNOCOUNT)
		p.addTokenCandidate(kwXACT_ABORT)
		return nil, errCollecting
	}

	if p.cur.Type != tokVARIABLE {
		// SET session option
		return p.parseSetOptionStmt(loc)
	}

	// SET @var ...
	stmt := &nodes.SetStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}
	stmt.Variable = p.cur.Str
	p.advance()

	// SET @var.property = expr | SET @var.method ( [ args ] )
	if p.cur.Type == '.' {
		p.advance() // consume .
		name, ok := p.parseIdentifier()
		if !ok {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Property = name
		if p.cur.Type == '(' {
			callLoc := p.prev.Loc
			args, err := p.parseSetMethodArgs()
			if err != nil {
				return nil, err
			}
			stmt.Value = &nodes.FuncCallExpr{
				Name: &nodes.TableRef{Object: name, Loc: nodes.Loc{Start: callLoc, End: callLoc + len(name)}},
				Args: args,
				Loc:  nodes.Loc{Start: callLoc, End: p.prevEnd()},
			}
			stmt.Loc.End = p.prevEnd()
			return stmt, nil
		}
	}

	if op := p.isCompoundAssign(); op != "" {
		stmt.Operator = op
		p.advance()
	} else if _, err := p.expect('='); err != nil {
		return nil, err
	}

	// SET @var = CURSOR ... (cursorDefinition, TSql160.g:16783)
	if p.cur.Type == kwCURSOR && stmt.Operator == "" && stmt.Property == "" {
		cursor := &nodes.DeclareCursorStmt{Loc: nodes.Loc{Start: p.pos(), End: -1}}
		p.advance() // consume CURSOR
		if err := p.parseCursorDefinitionTail(cursor); err != nil {
			return nil, err
		}
		cursor.Loc.End = p.prevEnd()
		stmt.Cursor = cursor
		stmt.Loc.End = p.prevEnd()
		return stmt, nil
	}

	// Completion: after SET @v = → expression context
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		p.addRuleCandidate("func_name")
		return nil, errCollecting
	}
	val, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, p.newParseError(p.cur.Loc, "expected expression after SET assignment")
	}
	stmt.Value = val
	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// setCommandNames are the options of setCommand (TSql160.g:16555) that take
// a constant or a name: SET LANGUAGE us_english, SET DATEFIRST 1, ...
// (GeneralSetCommandTypeHelper). They may be listed:
// SET DATEFIRST 1, LANGUAGE us_english.
var setCommandNames = map[string]bool{
	"LANGUAGE": true, "DATEFORMAT": true, "DATEFIRST": true, "DEADLOCK_PRIORITY": true,
	"LOCK_TIMEOUT": true, "QUERY_GOVERNOR_COST_LIMIT": true, "CONTEXT_INFO": true,
}

// setOffsetsItems are the words of SET OFFSETS (offsetItem, TSql160.g:16520).
var setOffsetsItems = map[string]bool{
	"SELECT": true, "FROM": true, "ORDER": true, "COMPUTE": true, "TABLE": true,
	"PROCEDURE": true, "PROC": true, "EXECUTE": true, "EXEC": true,
	"STATEMENT": true, "PARAM": true,
}

// parseSetOptionStmt parses SET session option statements.
//
// BNF: mssql/parser/bnf/set-transact-sql.bnf
// BNF: mssql/parser/bnf/set-transaction-isolation-level-transact-sql.bnf
// BNF: mssql/parser/bnf/set-ansi-defaults-transact-sql.bnf (predicateSetStatement)
// BNF: mssql/parser/bnf/set-statistics-transact-sql.bnf (setStatisticsStatement)
//
//	SET TRANSACTION ISOLATION LEVEL
//	    { READ UNCOMMITTED | READ COMMITTED | REPEATABLE READ | SNAPSHOT | SERIALIZABLE }
//	SET IDENTITY_INSERT table_name { ON | OFF }
//	SET { option_name } { ON | OFF | value }
//	SET option_name [ , ...n ] { ON | OFF }
//	SET STATISTICS { IO | TIME | PROFILE | XML } [ , ...n ] { ON | OFF }
//
// The list form is the one the indexed views reference uses:
// https://learn.microsoft.com/en-us/sql/relational-databases/views/create-indexed-views
//
// setStatements (TSql160.g:16386) reads each form by its words, and each
// part is required: the isolation level (TSql160.g:16589), the table and ON
// or OFF of IDENTITY_INSERT (TSql160.g:16650), the OFFSETS words and ON or
// OFF (TSql160.g:16504), the STATISTICS options and ON or OFF
// (TSql160.g:16671), ROWCOUNT's integer or variable (TSql160.g:16494),
// TEXTSIZE's and ERRLVL's signed integer (TSql160.g:16640, 16661),
// FIPS_FLAGGER's OFF or level string, a command's constant or name
// (setCommand, TSql160.g:16555), and a predicate option's ON or OFF
// (predicateSetStatement, TSql160.g:16470).
func (p *Parser) parseSetOptionStmt(loc int) (*nodes.SetOptionStmt, error) {
	stmt := &nodes.SetOptionStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}
	done := func(value nodes.ExprNode) (*nodes.SetOptionStmt, error) {
		stmt.Value = value
		if stmt.Options == nil {
			stmt.Options = []string{stmt.Option}
		}
		stmt.Loc.End = p.prevEnd()
		return stmt, nil
	}

	switch {
	case p.cur.Type == kwTRANSACTION || p.cur.Type == kwTRAN:
		// TRANSACTION ISOLATION LEVEL
		p.advance() // consume TRANSACTION or TRAN
		if _, err := p.expect(kwISOLATION); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwLEVEL); err != nil {
			return nil, err
		}
		stmt.Option = "TRANSACTION ISOLATION LEVEL"
		valLoc := p.pos()
		var level string
		switch p.cur.Type {
		case kwREAD:
			p.advance() // consume READ
			switch p.cur.Type {
			case kwUNCOMMITTED:
				level = "READ UNCOMMITTED"
			case kwCOMMITTED:
				level = "READ COMMITTED"
			default:
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
		case kwREPEATABLE:
			p.advance() // consume REPEATABLE
			if _, err := p.expect(kwREAD); err != nil {
				return nil, err
			}
			level = "REPEATABLE READ"
		case kwSERIALIZABLE:
			p.advance()
			level = "SERIALIZABLE"
		case kwSNAPSHOT:
			p.advance()
			level = "SNAPSHOT"
		default:
			return nil, p.syntaxErrorAtCur()
		}
		return done(&nodes.ColumnRef{Column: level, Loc: nodes.Loc{Start: valLoc, End: p.prevEnd()}})

	case p.cur.Type == kwIDENTITY_INSERT:
		// IDENTITY_INSERT table ON|OFF
		p.advance() // consume IDENTITY_INSERT
		stmt.Option = "IDENTITY_INSERT"
		tableRef, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		if tableRef == nil {
			return nil, p.syntaxErrorAtCur()
		}
		valLoc := p.pos()
		onOff := p.parseSetOnOff()
		if onOff == nil {
			return nil, p.syntaxErrorAtCur()
		}
		// Store "table ON/OFF" as option value
		tableName := tableRef.Object
		if tableRef.Schema != "" {
			tableName = tableRef.Schema + "." + tableRef.Object
		}
		return done(&nodes.ColumnRef{Column: tableName + " " + onOff.Column, Loc: nodes.Loc{Start: valLoc, End: -1}})

	case p.cur.Type == kwOFFSETS:
		// SET OFFSETS keyword_list { ON | OFF }
		p.advance() // consume OFFSETS
		var keywords []string
		for {
			if !p.isKeywordOrIdent() || p.isDelimited(p.cur) || !setOffsetsItems[strings.ToUpper(p.cur.Str)] {
				return nil, p.syntaxErrorAtCur()
			}
			keywords = append(keywords, strings.ToUpper(p.cur.Str))
			p.advance()
			if _, ok := p.match(','); !ok {
				break
			}
		}
		stmt.Option = "OFFSETS " + strings.Join(keywords, ", ")
		onOff := p.parseSetOnOff()
		if onOff == nil {
			return nil, p.syntaxErrorAtCur()
		}
		return done(onOff)

	case p.cur.Type == kwROWCOUNT, p.cur.Type == kwTEXTSIZE, p.cur.Type == kwERRLVL:
		stmt.Option = strings.ToUpper(p.cur.Str)
		p.advance()
		kinds := procValueKinds{signedInt: true}
		if stmt.Option == "ROWCOUNT" {
			kinds = procValueKinds{integer: true, variable: true}
		}
		value, err := p.parseProcValue(kinds)
		if err != nil {
			return nil, err
		}
		return done(value)

	case p.procIsWord("FIPS_FLAGGER"):
		stmt.Option = "FIPS_FLAGGER"
		p.advance()
		if p.cur.Type == kwOFF {
			return done(p.parseSetOnOff())
		}
		// fipsFlaggerLevel (TSql160.g:16573): the pinned server takes
		// 'ENTRY', 'INTERMEDIATE' and 'FULL', as a plain string (Msg 102).
		switch strings.ToUpper(p.cur.Str) {
		case "ENTRY", "INTERMEDIATE", "FULL":
			if p.cur.Type == tokSCONST {
				value, err := p.parsePrimary()
				if err != nil {
					return nil, err
				}
				return done(value)
			}
		}
		return nil, p.syntaxErrorAtCur()
	}

	// SET STATISTICS IO|TIME|PROFILE|XML [, ...] ON|OFF
	if p.cur.Type == kwSTATISTICS {
		statsTok := p.advance()
		if p.cur.Type == ',' {
			return nil, p.newParseError(statsTok.Loc, "SET STATISTICS cannot be listed with other options")
		}
		for {
			if !p.isSetListOption(p.cur, setStatisticsOptions) {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.Options = append(stmt.Options, "STATISTICS "+strings.ToUpper(p.cur.Str))
			p.advance()
			if _, ok := p.match(','); !ok {
				break
			}
		}
		stmt.Option = stmt.Options[0]
		onOff := p.parseSetOnOff()
		if onOff == nil {
			return nil, p.syntaxErrorAtCur()
		}
		return done(onOff)
	}

	// SET command value [, command value ...]
	if (p.isKeywordOrIdent() && !p.isDelimited(p.cur)) && setCommandNames[strings.ToUpper(p.cur.Str)] {
		var value nodes.ExprNode
		for {
			if !p.isKeywordOrIdent() || p.isDelimited(p.cur) || !setCommandNames[strings.ToUpper(p.cur.Str)] {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.Options = append(stmt.Options, strings.ToUpper(p.cur.Str))
			nameTok := p.advance()
			if p.cur.Type == ',' {
				return nil, p.newParseError(nameTok.Loc, fmt.Sprintf("SET %s needs a value before the next option", strings.ToUpper(nameTok.Str)))
			}
			v, err := p.parseProcConstantOrIdentifier()
			if err != nil {
				return nil, err
			}
			if value == nil {
				value = v
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		stmt.Option = stmt.Options[0]
		return done(value)
	}

	// Predicate option — validated against setOptions
	if !p.isValidOption(setOptions) {
		return nil, p.newParseError(p.cur.Loc, "expected SET option name")
	}
	name := p.cur
	stmt.Option = strings.ToUpper(p.cur.Str)
	p.advance()
	stmt.Options = []string{stmt.Option}

	// SET option_name, option_name [, ...] { ON | OFF }
	if p.cur.Type == ',' {
		if !p.isSetListOption(name, setOnOffListOptions) {
			return nil, p.newParseError(name.Loc, fmt.Sprintf("SET %s cannot be listed with other options", stmt.Option))
		}
		for p.cur.Type == ',' {
			p.advance() // consume ,
			if !p.isSetListOption(p.cur, setOnOffListOptions) {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.Options = append(stmt.Options, strings.ToUpper(p.cur.Str))
			p.advance()
		}
	}
	onOff := p.parseSetOnOff()
	if onOff == nil {
		return nil, p.syntaxErrorAtCur()
	}
	return done(onOff)
}

// isSetListOption reports whether tok names an option of opts that may appear
// in a list of SET options. SQL Server does not accept a delimited name there:
// it rejects SET ANSI_NULLS, [QUOTED_IDENTIFIER] ON.
func (p *Parser) isSetListOption(tok Token, opts optionSet) bool {
	if !p.isValidOptionToken(opts, tok) {
		return false
	}
	if tok.Type == tokIDENT && tok.Loc >= 0 && tok.Loc < len(p.source) {
		if c := p.source[tok.Loc]; c == '[' || c == '"' {
			return false
		}
	}
	return true
}

// parseSetOnOff consumes ON or OFF and returns it as a ColumnRef, or returns
// nil when the current token is neither.
func (p *Parser) parseSetOnOff() *nodes.ColumnRef {
	var value string
	switch p.cur.Type {
	case kwON:
		value = "ON"
	case kwOFF:
		value = "OFF"
	default:
		return nil
	}
	loc := p.pos()
	p.advance()
	return &nodes.ColumnRef{Column: value, Loc: nodes.Loc{Start: loc, End: p.prevEnd()}}
}

// parseSetMethodArgs parses parenthesizedOptExpressionWithDefaultList
// (TSql160.g:32501): ( [ { expression | DEFAULT } [ ,...n ] ] ).
func (p *Parser) parseSetMethodArgs() (*nodes.List, error) {
	p.advance() // consume (
	var args []nodes.Node
	for p.cur.Type != ')' {
		var arg nodes.ExprNode
		if p.cur.Type == kwDEFAULT {
			tok := p.advance()
			arg = &nodes.Literal{Type: nodes.LitDefault, Loc: nodes.Loc{Start: tok.Loc, End: p.prevEnd()}}
		} else {
			var err error
			arg, err = p.parseExpr()
			if err != nil {
				return nil, err
			}
			if arg == nil {
				return nil, p.syntaxErrorAtCur()
			}
		}
		args = append(args, arg)
		if _, ok := p.match(','); !ok {
			break
		}
		if p.cur.Type == ')' {
			return nil, p.syntaxErrorAtCur()
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.List{Items: args}, nil
}
