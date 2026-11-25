package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// APPLICATION SERVICE (Native Apps), read as an ast.ObjectStmt
// ---------------------------------------------------------------------------

func init() {
	registerObject(&objectForm{
		words:  []string{"APPLICATION", "SERVICE"},
		create: applicationTagless(parseCreateApplicationService, ""),
		alter:  applicationTagless(parseAlterApplicationService, "SET"),
		// https://docs.snowflake.com/en/sql-reference/sql/create-application-service
		orAlter: true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-application-service
		drop:       true,
		props:      applicationServiceProps,
		alterProps: applicationServiceAlterProps,
		unset: []string{"AUTO_RESUME", "AUTO_SUSPEND_SECS", "EXTERNAL_ACCESS_INTEGRATIONS",
			"QUERY_WAREHOUSE", "MIN_INSTANCES", "MAX_INSTANCES", "COMMENT"},
	})
}

// applicationServiceProps are the properties CREATE APPLICATION SERVICE and
// ALTER APPLICATION SERVICE ... SET document. EXECUTE_AS_ROLE and
// SPECIFICATION are CREATE's alone.
var applicationServiceProps = objectProps{
	"EXTERNAL_ACCESS_INTEGRATIONS": propNames,
	"QUERY_WAREHOUSE":              propName,
	"EXECUTE_AS_ROLE":              propName,
	"AUTO_RESUME":                  propBool,
	"AUTO_SUSPEND_SECS":            propNumber,
	"MIN_INSTANCES":                propNumber,
	"MAX_INSTANCES":                propNumber,
	"COMMENT":                      propString,
	"SPECIFICATION":                propString,
}

// applicationServiceAlterProps are ALTER APPLICATION SERVICE ... SET's:
// CREATE's but EXECUTE_AS_ROLE and SPECIFICATION.
var applicationServiceAlterProps = func() objectProps {
	m := objectProps{}
	for k, v := range applicationServiceProps {
		if k != "EXECUTE_AS_ROLE" && k != "SPECIFICATION" {
			m[k] = v
		}
	}
	return m
}()

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
	if err := rejectApplicationOrAlterIfNotExists(s); err != nil {
		return err
	}
	start := p.cur.Loc.Start
	if err := p.expectWords("FROM", "ARTIFACT", "REPOSITORY"); err != nil {
		return err
	}
	repo, err := p.parseObjectName()
	if err != nil {
		return err
	}
	p.addObjectKeywordsClause(s, "FROM ARTIFACT REPOSITORY", start).Names = []*ast.ObjectName{repo}
	start = p.cur.Loc.Start
	if err := p.expectWords("PACKAGE"); err != nil {
		return err
	}
	if _, err := p.parseApplicationNamedClause(s, "PACKAGE", start); err != nil {
		return err
	}
	if start := p.cur.Loc.Start; p.curIsWord("VERSION") && !p.nextIs('=') {
		p.advance() // consume VERSION
		if _, err := p.parseApplicationNamedClause(s, "VERSION", start); err != nil {
			return err
		}
	}
	props, err := p.parseObjectProperties(s, "", false)
	if err != nil {
		return err
	}
	if !s.OrAlter {
		return nil
	}
	// CREATE OR ALTER takes COMMENT and the required SPECIFICATION alone.
	if err := p.requireObjectProperties(props, "SPECIFICATION"); err != nil {
		return err
	}
	for _, o := range props.Options {
		if o.Name != "COMMENT" && o.Name != "SPECIFICATION" {
			return applicationOptionError(o)
		}
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
