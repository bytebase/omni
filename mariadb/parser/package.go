package parser

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// MariaDB's packages, which its default sql_mode reads as well as ORACLE does: a package's
// specification (CREATE PACKAGE) and its body (CREATE PACKAGE BODY), dropped and shown as a
// routine is. PACKAGE and BODY are no tokens of this lexer, so they are read as bare words.
//
// omni reads no CREATE PACKAGE BODY. Its action looks the specification up while it parses
// and refuses a package that does not exist (1305), so the server, asked without one,
// refuses every body; omni holds no list of packages.
//
// Ref: create_routine's sp_handler_package_spec arm, sp_handler_package_spec,
// sp_handler_package_body, drop_routine and show_param in MariaDB 11.8 sql/sql_yacc.yy.

// parseCreatePackageStmt parses CREATE PACKAGE, p.cur being PACKAGE. The caller has read
// CREATE [OR REPLACE] [DEFINER = user].
//
//	CREATE [OR REPLACE] [DEFINER = user] PACKAGE [IF NOT EXISTS] sp_name
//	    [COMMENT 'string' | SQL SECURITY {DEFINER | INVOKER}] ...
//	    [FUNCTION name (params) RETURNS type [characteristic ...];
//	     | PROCEDURE name (params) [characteristic ...];] ...
//	END
//
// A package publishes a function and a procedure by one name, but not two of either (1304).
//
// Ref: create_routine, opt_create_package_chistics_init, create_package_chistic and
// opt_package_specification_element_list in MariaDB 11.8 sql/sql_yacc.yy.
func (p *Parser) parseCreatePackageStmt(start int, orReplace bool, definer string) (*nodes.CreatePackageStmt, error) {
	p.advance() // consume PACKAGE

	// After PACKAGE, BODY is always the keyword: `CREATE PACKAGE body END` is the body of a
	// package named END.
	if p.atWord("BODY") {
		return nil, &ParseError{Message: "CREATE PACKAGE BODY is not supported", Position: p.cur.Loc}
	}

	stmt := &nodes.CreatePackageStmt{
		Loc:       nodes.Loc{Start: start},
		OrReplace: orReplace,
		Definer:   definer,
	}
	ifNotExists, err := p.parseCreateIfNotExists(orReplace)
	if err != nil {
		return nil, err
	}
	stmt.IfNotExists = ifNotExists

	ref, err := p.parseNamedRoutineRef()
	if err != nil {
		return nil, err
	}
	stmt.Name = ref

	// create_package_chistic: COMMENT and SQL SECURITY, and nothing else of a routine's.
	for {
		chStart := p.pos()
		switch {
		case p.cur.Type == kwCOMMENT:
			p.advance()
			if p.cur.Type != tokSCONST {
				return nil, p.syntaxErrorAtCur()
			}
			val := p.cur.Str
			p.advance()
			stmt.Characteristics = append(stmt.Characteristics, &nodes.RoutineCharacteristic{
				Loc: nodes.Loc{Start: chStart, End: p.prev.End}, Name: "COMMENT", Value: val,
			})
			continue
		case p.cur.Type == kwSQL && p.peekNext().Type == kwSECURITY:
			p.advance()
			p.advance()
			var val string
			switch p.cur.Type {
			case kwDEFINER:
				val = "DEFINER"
			case kwINVOKER:
				val = "INVOKER"
			default:
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			stmt.Characteristics = append(stmt.Characteristics, &nodes.RoutineCharacteristic{
				Loc: nodes.Loc{Start: chStart, End: p.prev.End}, Name: "SQL SECURITY", Value: val,
			})
			continue
		}
		break
	}

	seen := map[string]bool{}
	for p.cur.Type == kwFUNCTION || p.cur.Type == kwPROCEDURE {
		routine, err := p.parsePackageRoutineSpec()
		if err != nil {
			return nil, err
		}
		kind := "FUNCTION"
		if routine.IsProcedure {
			kind = "PROCEDURE"
		}
		key := kind + " " + strings.ToLower(routine.Name)
		if seen[key] {
			return nil, &ParseError{
				Message:  fmt.Sprintf("%s %s.%s already exists", kind, ref.Name, routine.Name),
				Position: routine.Loc.Start,
			}
		}
		seen[key] = true
		stmt.Routines = append(stmt.Routines, routine)
	}

	if _, err := p.expect(kwEND); err != nil {
		return nil, err
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parsePackageRoutineSpec parses one function or procedure of a package's specification,
// p.cur being FUNCTION or PROCEDURE, through its ';'. Its name takes no database.
//
//	FUNCTION name ([param [, param] ...]) RETURNS type [characteristic ...];
//	PROCEDURE name ([param [, param] ...]) [characteristic ...];
//
// Ref: package_specification_element, package_specification_function and
// package_specification_procedure in MariaDB 11.8 sql/sql_yacc.yy.
func (p *Parser) parsePackageRoutineSpec() (*nodes.PackageRoutineSpec, error) {
	routine := &nodes.PackageRoutineSpec{
		Loc:         nodes.Loc{Start: p.pos()},
		IsProcedure: p.cur.Type == kwPROCEDURE,
	}
	p.advance() // consume FUNCTION or PROCEDURE

	name, namePos, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	// The name is checked as a routine's (1458), and may not hold a dot even quoted, which
	// would read as the package's own pkg.routine.
	//
	// Ref: LEX::package_routine_start and LEX::make_sp_name_package_routine in MariaDB 11.8
	// sql/sql_lex.cc.
	if err := refuseRoutineName(name, namePos); err != nil {
		return nil, err
	}
	if strings.Contains(name, ".") {
		return nil, &ParseError{Message: "Incorrect routine name '" + name + "'", Position: namePos}
	}
	routine.Name = name

	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	// Each entry is a routine of its own, as package_routine_start makes it, so its parameters
	// are read as a standalone routine's: named once each (1330), and ROW(...), TYPE OF and ROW
	// TYPE OF among their types.
	defer p.enterStoredProgram(routineScope{noResultSet: !routine.IsProcedure, function: !routine.IsProcedure})()
	if p.cur.Type != ')' {
		var names []string
		for {
			param, err := p.parseSpParam(routine.IsProcedure, names)
			if err != nil {
				return nil, err
			}
			names = append(names, param.Name)
			routine.Params = append(routine.Params, param)
			if _, ok := p.match(','); !ok {
				break
			}
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	if !routine.IsProcedure {
		if _, err := p.expect(kwRETURNS); err != nil {
			return nil, err
		}
		dt, err := p.parseDataType()
		if err != nil {
			return nil, err
		}
		routine.Returns = dt
	}

	for {
		ch, ok, err := p.parseRoutineCharacteristic()
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		routine.Characteristics = append(routine.Characteristics, ch)
	}

	routine.Loc.End = p.prev.End
	if _, err := p.expect(';'); err != nil {
		return nil, err
	}
	return routine, nil
}

// parseDropPackageStmt parses DROP PACKAGE [BODY], p.cur being PACKAGE. BODY after PACKAGE
// is always the keyword, so `DROP PACKAGE BODY;` is a 1064.
//
//	DROP PACKAGE [BODY] [IF EXISTS] sp_name
//
// Ref: drop_routine and sp_handler in MariaDB 11.8 sql/sql_yacc.yy.
func (p *Parser) parseDropPackageStmt(start int) (*nodes.DropPackageStmt, error) {
	p.advance() // consume PACKAGE
	stmt := &nodes.DropPackageStmt{Loc: nodes.Loc{Start: start}}
	if p.atWord("BODY") {
		p.advance()
		stmt.Body = true
	}
	if p.cur.Type == kwIF {
		p.advance()
		if _, err := p.expect(kwEXISTS_KW); err != nil {
			return nil, err
		}
		stmt.IfExists = true
	}
	ref, err := p.parseNamedRoutineRef()
	if err != nil {
		return nil, err
	}
	stmt.Name = ref
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseShowCreatePackage parses the rest of SHOW CREATE PACKAGE [BODY] sp_name, p.cur being
// PACKAGE. Here BODY names a package unless a name follows it: `SHOW CREATE PACKAGE BODY`
// shows the package named BODY, as `SHOW CREATE PACKAGE BODY.a` shows a in database BODY.
//
// Ref: show_param's CREATE PACKAGE_MARIADB_SYM arms in MariaDB 11.8 sql/sql_yacc.yy.
func (p *Parser) parseShowCreatePackage(stmt *nodes.ShowStmt) error {
	p.advance() // consume PACKAGE
	stmt.Type = "CREATE PACKAGE"
	if p.atWord("BODY") {
		switch p.peekNext().Type {
		case ';', '.', tokEOF:
		default:
			p.advance()
			stmt.Type = "CREATE PACKAGE BODY"
		}
	}
	ref, err := p.parseNamedRoutineRef()
	if err != nil {
		return err
	}
	stmt.From = ref
	return nil
}

// parseShowPackage parses the rest of SHOW PACKAGE [BODY] STATUS and SHOW PACKAGE BODY CODE,
// p.cur being PACKAGE.
//
//	SHOW PACKAGE [BODY] STATUS [LIKE 'pattern' | WHERE expr]
//	SHOW PACKAGE BODY CODE sp_name
//
// Ref: show_param's PACKAGE_MARIADB_SYM arms in MariaDB 11.8 sql/sql_yacc.yy.
func (p *Parser) parseShowPackage(stmt *nodes.ShowStmt) error {
	p.advance() // consume PACKAGE
	body := false
	if p.atWord("BODY") {
		p.advance()
		body = true
	}
	switch {
	case p.cur.Type == kwSTATUS:
		p.advance()
		stmt.Type = "PACKAGE STATUS"
		if body {
			stmt.Type = "PACKAGE BODY STATUS"
		}
		return p.parseShowLikeOrWhere(stmt)
	case body && p.cur.Type == kwCODE:
		p.advance()
		stmt.Type = "PACKAGE BODY CODE"
		ref, err := p.parseNamedRoutineRef()
		if err != nil {
			return err
		}
		stmt.From = ref
		return nil
	}
	return p.syntaxErrorAtCur()
}
