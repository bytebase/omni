package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// VARIABLE and TYPE (schema-level variables and user-defined types), read as
// ast.ObjectStmt
// ---------------------------------------------------------------------------

func init() {
	registerObject(&objectForm{
		words:     []string{"VARIABLE"},
		create:    applicationTagless(parseCreateVariable, ""),
		alter:     applicationTagless(parseAlterVariable, "SET"),
		orReplace: true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-variable
		drop: true,
		// https://docs.snowflake.com/en/sql-reference/sql/undrop-variable
		undrop: true,
		props:  objectProps{"COMMENT": propString},
	})
	registerObject(&objectForm{
		words:     []string{"TYPE"},
		create:    applicationTagless(parseCreateType, ""),
		alter:     applicationTagless(parseAlterType, "SET"),
		orReplace: true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-type writes no
		// IF EXISTS.
		drop:           true,
		dropNoIfExists: true,
		// https://docs.snowflake.com/en/sql-reference/sql/undrop-type
		undrop: true,
		props:  objectProps{"COMMENT": propString},
	})
}

// parseCreateVariable reads what follows the name of
//
//	CREATE [ OR REPLACE ] VARIABLE [ IF NOT EXISTS ] <name> <data_type>
//	  [ DEFAULT <default_value> ]
//	  [ COMMENT = '<string_literal>' ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-variable
//
// The data type is a clause of its own holding Type; DEFAULT holds its value
// as one nameless option whose Expr is the expression.
func parseCreateVariable(p *Parser, s *ast.ObjectStmt) error {
	if err := p.parseObjectTypeClause(s, ""); err != nil {
		return err
	}
	if p.cur.Type == kwDEFAULT {
		if err := p.parseObjectDefault(s, "DEFAULT"); err != nil {
			return err
		}
	}
	return p.parseObjectComment(s)
}

// parseAlterVariable reads what follows the name of
//
//	ALTER VARIABLE [ IF EXISTS ] <name> RENAME TO <new_name>
//	ALTER VARIABLE [ IF EXISTS ] <name> SET DEFAULT <default_value>
//	ALTER VARIABLE [ IF EXISTS ] <name> DROP DEFAULT
//	ALTER VARIABLE [ IF EXISTS ] <name> SET COMMENT = '<string_literal>'
//	ALTER VARIABLE [ IF EXISTS ] <name> UNSET COMMENT
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-variable
func parseAlterVariable(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	switch {
	case p.cur.Type == kwSET && p.peekNext().Type == kwDEFAULT:
		p.advance() // consume SET
		if err := p.parseObjectDefault(s, "SET DEFAULT"); err != nil {
			return err
		}
		s.Clauses[len(s.Clauses)-1].Loc.Start = start
		return nil
	case p.cur.Type == kwDROP && p.peekNext().Type == kwDEFAULT:
		p.advance() // consume DROP
		p.advance() // consume DEFAULT
		p.addObjectKeywordsClause(s, "DROP DEFAULT", start)
		return nil
	}
	return p.parseObjectAlterCommentOrRename(s, actSet|actUnset|actRename)
}

// parseCreateType reads what follows the name of
//
//	CREATE [ OR REPLACE ] TYPE [ IF NOT EXISTS ] <name> AS <type>
//	  [ COMMENT = '<string_literal>' ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-type
func parseCreateType(p *Parser, s *ast.ObjectStmt) error {
	if p.cur.Type != kwAS {
		return p.syntaxErrorAtCur()
	}
	if err := p.parseObjectTypeClause(s, "AS"); err != nil {
		return err
	}
	return p.parseObjectComment(s)
}

// parseAlterType reads what follows the name of
//
//	ALTER TYPE [ IF EXISTS ] <name> SET COMMENT = '<string_literal>'
//	ALTER TYPE [ IF EXISTS ] <name> UNSET COMMENT
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-type
func parseAlterType(p *Parser, s *ast.ObjectStmt) error {
	return p.parseObjectAlterCommentOrRename(s, actSet|actUnset)
}

// parseObjectAlterCommentOrRename reads one of the shared ALTER actions allow
// names, failing when none starts at cur.
func (p *Parser) parseObjectAlterCommentOrRename(s *ast.ObjectStmt, allow objectActions) error {
	ok, err := p.parseObjectAlterCommon(s, allow)
	if err != nil || ok {
		return err
	}
	return p.syntaxErrorAtCur()
}

// parseObjectTypeClause reads a data type, after keywords when they are not
// empty (cur is then their first), into a clause holding Type.
func (p *Parser) parseObjectTypeClause(s *ast.ObjectStmt, keywords string) error {
	start := p.cur.Loc.Start
	if keywords != "" {
		p.advance() // consume AS
	}
	typ, err := p.parseDataType()
	if err != nil {
		return err
	}
	s.Clauses = append(s.Clauses, &ast.ObjectClause{Keywords: keywords, Type: typ, Loc: ast.Loc{Start: start, End: p.prev.Loc.End}})
	return nil
}

// parseObjectDefault reads DEFAULT <expr> (cur is DEFAULT) into a clause of
// keywords holding the expression as one nameless option.
func (p *Parser) parseObjectDefault(s *ast.ObjectStmt, keywords string) error {
	start := p.cur.Loc.Start
	p.advance() // consume DEFAULT
	exprStart := p.cur.Loc.Start
	expr, err := p.parseExpr()
	if err != nil {
		return err
	}
	opt := &ast.CopyOption{Expr: expr, Loc: ast.Loc{Start: exprStart, End: p.prev.Loc.End}}
	p.addObjectKeywordsClause(s, keywords, start).Options = []*ast.CopyOption{opt}
	return nil
}

// parseObjectComment reads an optional COMMENT = '<string_literal>' run; the
// registry holds its name and value to the object's props.
func (p *Parser) parseObjectComment(s *ast.ObjectStmt) error {
	if !p.startsObjectProperty() {
		return nil
	}
	_, err := p.parseObjectProperties(s, "", false)
	return err
}
