package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// APPLICATION SERVICE (Native Apps), read as an ast.ObjectStmt
// ---------------------------------------------------------------------------

func init() {
	registerObject(&objectForm{
		words:  []string{"APPLICATION", "SERVICE"},
		create: parseCreateApplicationService,
		alter:  parseAlterApplicationService,
		// https://docs.snowflake.com/en/sql-reference/sql/create-application-service
		orAlter: true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-application-service
		drop: true,
	})
}

// parseCreateApplicationService reads what follows the name of
//
//	CREATE APPLICATION SERVICE [ IF NOT EXISTS ] <name>
//	  FROM ARTIFACT REPOSITORY <repository_name> PACKAGE <package_name> [ VERSION <version> ]
//	  [ EXTERNAL_ACCESS_INTEGRATIONS = ( <integration> [ , ... ] ) ]
//	  [ QUERY_WAREHOUSE = <warehouse> ] [ EXECUTE_AS_ROLE = <role> ]
//	  [ AUTO_RESUME = { TRUE | FALSE } ] [ AUTO_SUSPEND_SECS = <num> ]
//	  [ MIN_INSTANCES = <num> ] [ MAX_INSTANCES = <num> ]
//	  [ COMMENT = '<string_literal>' ] [ SPECIFICATION = '<specification>' ]
//	CREATE OR ALTER APPLICATION SERVICE <name>
//	  FROM ARTIFACT REPOSITORY <repository_name> PACKAGE <package_name> [ VERSION <version> ]
//	  [ COMMENT = '<string_literal>' ] SPECIFICATION = '<specification>'
//
// https://docs.snowflake.com/en/sql-reference/sql/create-application-service
func parseCreateApplicationService(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	if err := p.expectWords("FROM", "ARTIFACT", "REPOSITORY"); err != nil {
		return err
	}
	if _, err := p.parseApplicationNamedClause(s, "FROM ARTIFACT REPOSITORY", start); err != nil {
		return err
	}
	start = p.cur.Loc.Start
	if err := p.expectWords("PACKAGE"); err != nil {
		return err
	}
	if _, err := p.parseApplicationNamedClause(s, "PACKAGE", start); err != nil {
		return err
	}
	if start := p.cur.Loc.Start; p.curIsWord("VERSION") && p.peekNext().Type != '=' {
		p.advance() // consume VERSION
		if _, err := p.parseApplicationNamedClause(s, "VERSION", start); err != nil {
			return err
		}
	}
	props, err := p.parseObjectProperties(s, "", false)
	if err != nil {
		return err
	}
	if s.OrAlter {
		return p.requireObjectProperties(props, "SPECIFICATION")
	}
	return nil
}

// parseAlterApplicationService reads what follows the name of
//
//	ALTER APPLICATION SERVICE [ IF EXISTS ] <name> { SUSPEND | RESUME }
//	ALTER APPLICATION SERVICE [ IF EXISTS ] <name> UPGRADE [ TO VERSION <version> ]
//	ALTER APPLICATION SERVICE [ IF EXISTS ] <name> SET <property> [ ... ]
//	ALTER APPLICATION SERVICE [ IF EXISTS ] <name> UNSET <property> [ , ... ]
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-application-service
func parseAlterApplicationService(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	switch {
	case p.acceptWords("SUSPEND"):
		p.addObjectKeywordsClause(s, "SUSPEND", start)
		return nil
	case p.acceptWords("RESUME"):
		p.addObjectKeywordsClause(s, "RESUME", start)
		return nil
	case p.acceptWords("UPGRADE", "TO", "VERSION"):
		_, err := p.parseApplicationNamedClause(s, "UPGRADE TO VERSION", start)
		return err
	case p.acceptWords("UPGRADE"):
		p.addObjectKeywordsClause(s, "UPGRADE", start)
		return nil
	}
	ok, err := p.parseObjectAlterCommon(s, actSet|actUnset)
	if err != nil || ok {
		return err
	}
	return p.syntaxErrorAtCur()
}
