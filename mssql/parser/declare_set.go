// Package parser - declare_set.go implements T-SQL DECLARE and SET statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
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
func (p *Parser) parseDeclareStmt() *nodes.DeclareStmt {
	loc := p.pos()
	p.advance() // consume DECLARE

	stmt := &nodes.DeclareStmt{
		Loc: nodes.Loc{Start: loc},
	}

	var vars []nodes.Node
	for {
		vd := p.parseVariableDecl()
		if vd == nil {
			break
		}
		vars = append(vars, vd)
		if _, ok := p.match(','); !ok {
			break
		}
	}
	stmt.Variables = &nodes.List{Items: vars}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseVariableDecl parses a single variable declaration.
//
// BNF: mssql/parser/bnf/declare-local-variable-transact-sql.bnf
//
//	@local_variable [AS] data_type [ = value ]
//	| @cursor_variable_name CURSOR
//	| @table_variable_name [AS] TABLE ( { <column_definition> | <table_constraint> | <table_index> } [ , ...n ] )
func (p *Parser) parseVariableDecl() *nodes.VariableDecl {
	if p.cur.Type != tokVARIABLE {
		return nil
	}

	loc := p.pos()
	vd := &nodes.VariableDecl{
		Name: p.cur.Str,
		Loc:  nodes.Loc{Start: loc},
	}
	p.advance() // consume @var

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
				col := p.parseColumnDef()
				if col != nil {
					cols = append(cols, col)
				}
				if _, ok := p.match(','); !ok {
					break
				}
			}
			_, _ = p.expect(')')
			vd.TableDef = &nodes.List{Items: cols}
		}
		vd.Loc.End = p.pos()
		return vd
	}

	// CURSOR type
	if p.cur.Type == kwCURSOR {
		p.advance()
		vd.IsCursor = true
		vd.Loc.End = p.pos()
		return vd
	}

	// Data type
	vd.DataType = p.parseDataType()

	// Optional default value
	if p.cur.Type == '=' {
		p.advance()
		vd.Default = p.parseExpr()
	}

	vd.Loc.End = p.pos()
	return vd
}

// parseSetStmt parses a SET statement.
//
// BNF: mssql/parser/bnf/set-transact-sql.bnf
//
//	SET @local_variable { = | += | -= | *= | /= | %= | &= | ^= | |= } expression
//	SET @local_variable.property_name = expression
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
//	SET IDENTITY_INSERT table_name { ON | OFF }
//	SET STATISTICS { IO | TIME | PROFILE | XML } { ON | OFF }
func (p *Parser) parseSetStmt() nodes.StmtNode {
	loc := p.pos()
	p.advance() // consume SET

	if p.cur.Type == tokVARIABLE {
		// SET @var { = | += | -= | *= | /= | %= | &= | ^= | |= } expr
		stmt := &nodes.SetStmt{
			Loc: nodes.Loc{Start: loc},
		}
		stmt.Variable = p.cur.Str
		p.advance()
		// Check for compound assignment operators or simple =
		if op := p.isCompoundAssign(); op != "" {
			stmt.Operator = op
			p.advance()
			stmt.Value = p.parseExpr()
		} else if _, err := p.expect('='); err == nil {
			stmt.Value = p.parseExpr()
		}
		stmt.Loc.End = p.pos()
		return stmt
	}

	// SET session option
	return p.parseSetOptionStmt(loc)
}

// parseSetOptionStmt parses SET session option statements.
//
// BNF: mssql/parser/bnf/set-transact-sql.bnf
// BNF: mssql/parser/bnf/set-transaction-isolation-level-transact-sql.bnf
//
//	SET TRANSACTION ISOLATION LEVEL
//	    { READ UNCOMMITTED | READ COMMITTED | REPEATABLE READ | SNAPSHOT | SERIALIZABLE }
//	SET IDENTITY_INSERT table_name { ON | OFF }
//	SET { option_name } { ON | OFF | value }
func (p *Parser) parseSetOptionStmt(loc int) *nodes.SetOptionStmt {
	stmt := &nodes.SetOptionStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// TRANSACTION ISOLATION LEVEL
	if p.cur.Type == kwTRANSACTION || (p.isIdentLike() && matchesKeywordCI(p.cur.Str, "TRANSACTION")) {
		p.advance() // consume TRANSACTION
		// ISOLATION
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "ISOLATION") {
			p.advance()
		}
		// LEVEL
		if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "LEVEL") {
			p.advance()
		}
		// isolation level name: READ UNCOMMITTED / READ COMMITTED / REPEATABLE READ / SNAPSHOT / SERIALIZABLE
		stmt.Option = "TRANSACTION ISOLATION LEVEL"
		valLoc := p.pos()
		var level string
		if p.cur.Type == kwREAD || (p.isIdentLike() && matchesKeywordCI(p.cur.Str, "READ")) {
			p.advance() // consume READ
			if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "UNCOMMITTED") {
				p.advance()
				level = "READ UNCOMMITTED"
			} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "COMMITTED") {
				p.advance()
				level = "READ COMMITTED"
			} else {
				level = "READ"
			}
		} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "REPEATABLE") {
			p.advance() // consume REPEATABLE
			if p.cur.Type == kwREAD || (p.isIdentLike() && matchesKeywordCI(p.cur.Str, "READ")) {
				p.advance()
			}
			level = "REPEATABLE READ"
		} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SERIALIZABLE") {
			p.advance()
			level = "SERIALIZABLE"
		} else if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "SNAPSHOT") {
			p.advance()
			level = "SNAPSHOT"
		}
		stmt.Value = &nodes.ColumnRef{Column: level, Loc: nodes.Loc{Start: valLoc, End: p.pos()}}
		stmt.Loc.End = p.pos()
		return stmt
	}

	// IDENTITY_INSERT table ON|OFF
	if p.cur.Type == kwIDENTITY_INSERT {
		p.advance() // consume IDENTITY_INSERT
		stmt.Option = "IDENTITY_INSERT"
		// table name
		tableRef := p.parseTableRef()
		// ON or OFF
		var onoff string
		if p.cur.Type == kwON {
			onoff = "ON"
			p.advance()
		} else if p.cur.Type == kwOFF {
			onoff = "OFF"
			p.advance()
		}
		valLoc := p.pos()
		// Store "table ON/OFF" as option value
		var tableName string
		if tableRef != nil {
			if tableRef.Schema != "" {
				tableName = tableRef.Schema + "." + tableRef.Object
			} else {
				tableName = tableRef.Object
			}
		}
		stmt.Value = &nodes.ColumnRef{Column: tableName + " " + onoff, Loc: nodes.Loc{Start: valLoc}}
		stmt.Loc.End = p.pos()
		return stmt
	}

	// SET OFFSETS keyword_list { ON | OFF }
	// BNF: mssql/parser/bnf/set-offsets-transact-sql.bnf
	//
	//	SET OFFSETS keyword_list { ON | OFF }
	//	keyword_list: SELECT, FROM, ORDER, COMPUTE, TABLE, PROCEDURE, STATEMENT, PARAM, EXECUTE
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "OFFSETS") {
		p.advance() // consume OFFSETS
		stmt.Option = "OFFSETS"
		// Collect comma-separated keyword list
		var keywords []string
		for {
			if p.cur.Type == kwON || p.cur.Type == kwOFF {
				break
			}
			if p.cur.Type == ';' || p.cur.Type == tokEOF || p.cur.Type == kwGO {
				break
			}
			if p.cur.Type == ',' {
				p.advance()
				continue
			}
			keywords = append(keywords, strings.ToUpper(p.cur.Str))
			p.advance()
		}
		if len(keywords) > 0 {
			stmt.Option = "OFFSETS " + strings.Join(keywords, ", ")
		}
		if p.cur.Type == kwON {
			onLoc := p.pos()
			p.advance()
			stmt.Value = &nodes.ColumnRef{Column: "ON", Loc: nodes.Loc{Start: onLoc}}
		} else if p.cur.Type == kwOFF {
			offLoc := p.pos()
			p.advance()
			stmt.Value = &nodes.ColumnRef{Column: "OFF", Loc: nodes.Loc{Start: offLoc}}
		}
		stmt.Loc.End = p.pos()
		return stmt
	}

	// Generic option name
	if p.isIdentLike() || p.cur.Type == kwNOCOUNT || p.cur.Type == kwXACT_ABORT ||
		p.cur.Type == kwROWCOUNT || p.cur.Type == kwTEXTSIZE ||
		p.cur.Type == kwSTATISTICS {
		stmt.Option = strings.ToUpper(p.cur.Str)
		p.advance()

		// STATISTICS IO|TIME|PROFILE|XML (multi-word option)
		if strings.EqualFold(stmt.Option, "STATISTICS") && p.isIdentLike() {
			stmt.Option = stmt.Option + " " + strings.ToUpper(p.cur.Str)
			p.advance()
		}

		if p.cur.Type == kwON {
			onLoc := p.pos()
			p.advance()
			stmt.Value = &nodes.ColumnRef{Column: "ON", Loc: nodes.Loc{Start: onLoc}}
		} else if p.cur.Type == kwOFF {
			offLoc := p.pos()
			p.advance()
			stmt.Value = &nodes.ColumnRef{Column: "OFF", Loc: nodes.Loc{Start: offLoc}}
		} else if p.cur.Type != ';' && p.cur.Type != tokEOF && p.cur.Type != kwGO {
			// Could be SET ROWCOUNT n, SET LANGUAGE ..., SET DATEFORMAT ..., etc.
			stmt.Value = p.parseExpr()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}
