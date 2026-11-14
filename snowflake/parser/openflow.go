package parser

import (
	"slices"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Openflow DEPLOYMENT, RUNTIME and DATA PLANE INTEGRATION, read as an
// ast.ObjectStmt
// ---------------------------------------------------------------------------

func init() {
	registerObject(&objectForm{
		words:  []string{"OPENFLOW", "DEPLOYMENT"},
		create: parseCreateOpenflowDeployment,
		alter:  parseAlterOpenflowDeployment,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-openflow-deployment
		drop:  true,
		props: openflowDeploymentProps,
		unset: openflowDeploymentAlterProps,
	})
	registerObject(&objectForm{
		words:  []string{"OPENFLOW", "RUNTIME"},
		create: parseCreateOpenflowRuntime,
		alter:  parseAlterOpenflowRuntime,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-openflow-runtime
		drop:        true,
		dropCascade: true,
		props:       openflowRuntimeProps,
		unset:       []string{"DISPLAY_NAME", "EXTERNAL_ACCESS_INTEGRATIONS", "COMMENT"},
	})
	registerObject(&objectForm{
		words: []string{"OPENFLOW", "DATA", "PLANE", "INTEGRATION"},
		alter: parseAlterOpenflowDataPlaneIntegration,
		// The page writes no IF EXISTS.
		noIfExists: true,
		props:      objectProps{"EVENT_TABLE": propString},
		unset:      []string{},
	})
}

// openflowDeploymentProps are the properties CREATE OPENFLOW DEPLOYMENT
// documents; ALTER ... SET takes openflowDeploymentAlterProps.
var openflowDeploymentProps = objectProps{
	"DEPLOYMENT_TYPE":                 propOneOf("SNOWFLAKE", "BYOC"),
	"VPC_TYPE":                        propString,
	"USE_PRIVATE_LINK":                propBool,
	"USE_USER_AUTH_OVER_PRIVATE_LINK": propBool,
	"CUSTOM_INGRESS_HOSTNAME":         propString,
	"DISPLAY_NAME":                    propString,
	"COMMENT":                         propString,
	"EVENT_TABLE":                     propStringOr("NONE"),
}

// openflowDeploymentAlterProps are the properties ALTER OPENFLOW DEPLOYMENT
// sets and unsets.
var openflowDeploymentAlterProps = []string{"DISPLAY_NAME", "COMMENT", "EVENT_TABLE"}

// openflowRuntimeProps are the properties CREATE OPENFLOW RUNTIME documents;
// ALTER ... SET takes openflowRuntimeAlterProps.
var openflowRuntimeProps = objectProps{
	"NODE_TYPE":                    propOneOf("SMALL", "MEDIUM", "LARGE"),
	"NODE_TYPE_TIER":               propString,
	"MIN_NODES":                    propNumber,
	"MAX_NODES":                    propNumber,
	"EXECUTE_AS_ROLE":              propName,
	"EXTERNAL_ACCESS_INTEGRATIONS": propNames,
	"DISPLAY_NAME":                 propString,
	"COMMENT":                      propString,
}

// openflowRuntimeAlterProps are the properties ALTER OPENFLOW RUNTIME ... SET
// documents: all of CREATE's but NODE_TYPE.
var openflowRuntimeAlterProps = []string{
	"DISPLAY_NAME", "MIN_NODES", "MAX_NODES", "NODE_TYPE_TIER", "EXTERNAL_ACCESS_INTEGRATIONS", "EXECUTE_AS_ROLE", "COMMENT",
}

// pgOpenflowOnlyProps fails when clause c sets a property outside names. A
// form whose CREATE and SET document different properties holds each to its
// own list, beside the shapes objectForm.props checks.
func pgOpenflowOnlyProps(c *ast.ObjectClause, names ...string) error {
	if c == nil {
		return nil
	}
	for _, o := range c.Options {
		if !slices.Contains(names, o.Name) {
			return &ParseError{
				Position: o.Loc.Start, End: o.Loc.End,
				Message: "syntax error at or near " + o.Name,
			}
		}
	}
	return nil
}

// pgOpenflowAlterCommon reads one of the shared ALTER actions allow names, as
// parseObjectAlterCommon does, and holds a SET's properties to setProps.
func pgOpenflowAlterCommon(p *Parser, s *ast.ObjectStmt, allow objectActions, setProps []string) (bool, error) {
	n := len(s.Clauses)
	ok, err := p.parseObjectAlterCommon(s, allow)
	if err != nil || !ok {
		return ok, err
	}
	if c := s.Clauses[n]; c.Keywords == "SET" {
		return true, pgOpenflowOnlyProps(c, setProps...)
	}
	return true, nil
}

// parseCreateOpenflowDeployment reads what follows the name of
//
//	CREATE OPENFLOW DEPLOYMENT [ IF NOT EXISTS ] <name>
//	  [ DEPLOYMENT_TYPE = { SNOWFLAKE | BYOC } ]
//	  [ VPC_TYPE = { 'MANAGED' | 'PROVIDED' } ]
//	  [ USE_PRIVATE_LINK = { TRUE | FALSE } ]
//	  [ USE_USER_AUTH_OVER_PRIVATE_LINK = { TRUE | FALSE } ]
//	  [ CUSTOM_INGRESS_HOSTNAME = '<hostname>' ]
//	  [ DISPLAY_NAME = '<display_name>' ] [ COMMENT = '<string>' ]
//	  [ EVENT_TABLE = { '<event_table_name>' | NONE } ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-openflow-deployment
func parseCreateOpenflowDeployment(p *Parser, s *ast.ObjectStmt) error {
	_, err := p.parseObjectProperties(s, "", false)
	return err
}

// parseAlterOpenflowDeployment reads what follows the name of
//
//	ALTER OPENFLOW DEPLOYMENT [ IF EXISTS ] <name> SET
//	  [ DISPLAY_NAME = '<display_name>' ] [ COMMENT = '<string>' ]
//	  [ EVENT_TABLE = { '<event_table_name>' | NONE } ]
//	ALTER OPENFLOW DEPLOYMENT [ IF EXISTS ] <name> UNSET { DISPLAY_NAME | COMMENT | EVENT_TABLE } [ , ... ]
//	ALTER OPENFLOW DEPLOYMENT [ IF EXISTS ] <name> RENAME TO <new_name>
//	ALTER OPENFLOW DEPLOYMENT [ IF EXISTS ] <name> { UPGRADE | TERMINATE }
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-openflow-deployment
func parseAlterOpenflowDeployment(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	for _, kw := range []string{"UPGRADE", "TERMINATE"} {
		if p.acceptWords(kw) {
			p.addObjectKeywordsClause(s, kw, start)
			return nil
		}
	}
	ok, err := pgOpenflowAlterCommon(p, s, actSet|actUnset|actRename, openflowDeploymentAlterProps)
	if err != nil || ok {
		return err
	}
	return p.syntaxErrorAtCur()
}

// parseCreateOpenflowRuntime reads what follows the name of
//
//	CREATE OPENFLOW RUNTIME [ IF NOT EXISTS ] [ <db>.<schema>. ]<name>
//	  IN DEPLOYMENT <deployment_name>
//	  NODE_TYPE = { SMALL | MEDIUM | LARGE }
//	  [ NODE_TYPE_TIER = { 'S1' | 'S2' | 'S3' | 'M4' | 'M6' | 'L8' } ]
//	  MIN_NODES = <num> MAX_NODES = <num>
//	  EXECUTE_AS_ROLE = <role_name>
//	  [ EXTERNAL_ACCESS_INTEGRATIONS = ( <integration_name> [ , ... ] ) ]
//	  [ DISPLAY_NAME = '<display_name>' ] [ COMMENT = '<string>' ]
//
// The name is one part or three, never two. The IN DEPLOYMENT clause records
// the deployment in Names.
//
// https://docs.snowflake.com/en/sql-reference/sql/create-openflow-runtime
func parseCreateOpenflowRuntime(p *Parser, s *ast.ObjectStmt) error {
	if s.Name.Database.IsEmpty() && !s.Name.Schema.IsEmpty() {
		return &ParseError{
			Position: s.Name.Loc.Start, End: s.Name.Loc.End,
			Message: "syntax error: a runtime's name is <name> or <db>.<schema>.<name>",
		}
	}
	start := p.cur.Loc.Start
	if err := p.expectWords("IN", "DEPLOYMENT"); err != nil {
		return err
	}
	dep, err := p.parseObjectSimpleName()
	if err != nil {
		return err
	}
	p.addObjectKeywordsClause(s, "IN DEPLOYMENT", start).Names = []*ast.ObjectName{dep}
	props, err := p.parseObjectProperties(s, "", true)
	if err != nil {
		return err
	}
	return p.requireObjectProperties(props, "NODE_TYPE", "MIN_NODES", "MAX_NODES", "EXECUTE_AS_ROLE")
}

// parseAlterOpenflowRuntime reads what follows the name of
//
//	ALTER OPENFLOW RUNTIME [ IF EXISTS ] <name> SET
//	  [ DISPLAY_NAME = '<display_name>' ] [ MIN_NODES = <num> ] [ MAX_NODES = <num> ]
//	  [ NODE_TYPE_TIER = { 'S1' | 'S2' | 'S3' | 'M4' | 'M6' | 'L8' } ]
//	  [ EXTERNAL_ACCESS_INTEGRATIONS = ( <integration_name> [ , ... ] ) ]
//	  [ EXECUTE_AS_ROLE = <role_name> ] [ COMMENT = '<string>' ]
//	ALTER OPENFLOW RUNTIME [ IF EXISTS ] <name> UNSET
//	  { DISPLAY_NAME | EXTERNAL_ACCESS_INTEGRATIONS | COMMENT } [ , ... ]
//	ALTER OPENFLOW RUNTIME [ IF EXISTS ] <name> { ADD | REMOVE }
//	  EXTERNAL_ACCESS_INTEGRATIONS = ( <integration_name> [ , ... ] )
//	ALTER OPENFLOW RUNTIME [ IF EXISTS ] <name> RENAME TO <new_name>
//	ALTER OPENFLOW RUNTIME [ IF EXISTS ] <name> UPGRADE [ RECOVERY | FORCE ]
//	ALTER OPENFLOW RUNTIME [ IF EXISTS ] <name> SUSPEND
//	ALTER OPENFLOW RUNTIME [ IF EXISTS ] <name> { RESUME | RESTART } [ RECOVERY ]
//	ALTER OPENFLOW RUNTIME [ IF EXISTS ] <name> TERMINATE [ CASCADE ]
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-openflow-runtime
func parseAlterOpenflowRuntime(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	// Each action and the one word that may follow it.
	for _, a := range [][]string{
		{"UPGRADE", "RECOVERY", "FORCE"},
		{"SUSPEND"},
		{"RESUME", "RECOVERY"},
		{"RESTART", "RECOVERY"},
		{"TERMINATE", "CASCADE"},
	} {
		if !p.acceptWords(a[0]) {
			continue
		}
		kw := a[0]
		for _, opt := range a[1:] {
			if p.acceptWords(opt) {
				kw += " " + opt
				break
			}
		}
		p.addObjectKeywordsClause(s, kw, start)
		return nil
	}
	for _, kw := range []string{"ADD", "REMOVE"} {
		if !p.acceptWords(kw) {
			continue
		}
		if !p.atWords("EXTERNAL_ACCESS_INTEGRATIONS") || !p.startsObjectProperty() {
			return p.syntaxErrorAtCur()
		}
		opt, err := p.parseObjectProperty()
		if err != nil {
			return err
		}
		if !propNames.matches(opt) {
			return &ParseError{Position: opt.Loc.Start, End: opt.Loc.End, Message: "syntax error at or near " + opt.Name}
		}
		p.addObjectKeywordsClause(s, kw, start).Options = []*ast.CopyOption{opt}
		return nil
	}
	ok, err := pgOpenflowAlterCommon(p, s, actSet|actUnset|actRename, openflowRuntimeAlterProps)
	if err != nil || ok {
		return err
	}
	return p.syntaxErrorAtCur()
}

// parseAlterOpenflowDataPlaneIntegration reads what follows the name of
//
//	ALTER OPENFLOW DATA PLANE INTEGRATION <name> SET EVENT_TABLE = '<event_table_name>'
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-oflow-data-plane
func parseAlterOpenflowDataPlaneIntegration(p *Parser, s *ast.ObjectStmt) error {
	ok, err := p.parseObjectAlterCommon(s, actSet)
	if err != nil || ok {
		return err
	}
	return p.syntaxErrorAtCur()
}
