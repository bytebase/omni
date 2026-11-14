package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// IMAGE REPOSITORY and ARTIFACT REPOSITORY, read as an ast.ObjectStmt
// ---------------------------------------------------------------------------

func init() {
	registerObject(&objectForm{
		words:     []string{"IMAGE", "REPOSITORY"},
		create:    parseCreateImageRepository,
		alter:     parseAlterImageRepository,
		orReplace: true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-image-repository
		drop:       true,
		props:      imageRepositoryProps,
		alterProps: objectProps{"COMMENT": propString},
		unset:      []string{"COMMENT"},
	})
	registerObject(&objectForm{
		words:     []string{"ARTIFACT", "REPOSITORY"},
		create:    parseCreateArtifactRepository,
		alter:     parseAlterArtifactRepository,
		orReplace: true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-artifact-repository
		drop:       true,
		props:      artifactRepositoryProps,
		alterProps: artifactRepositoryAlterProps,
		unset:      []string{"COMMENT", "AUTHENTICATION_SECRET", "PACKAGE_RETENTION"},
	})
}

// imageRepositoryProps are the properties CREATE IMAGE REPOSITORY documents;
// ALTER ... SET takes COMMENT alone.
var imageRepositoryProps = objectProps{
	"ENCRYPTION": propGroup,
	"COMMENT":    propString,
}

// artifactRepositoryProps are the properties CREATE ARTIFACT REPOSITORY
// documents.
// API_INTEGRATION is a string as the Syntax block writes it (one example on
// the page writes a bare name).
var artifactRepositoryProps = objectProps{
	"TYPE":                  propOneOf("APPLICATION", "PYPI"),
	"API_INTEGRATION":       propString,
	"INDEX_URL":             propString,
	"AUTHENTICATION_SECRET": propQName,
	"PACKAGE_RETENTION":     propOneOf("ON_OBJECT_CREATION"),
	"COMMENT":               propString,
}

// artifactRepositoryAlterProps are ALTER ARTIFACT REPOSITORY ... SET's:
// CREATE's but TYPE and API_INTEGRATION.
var artifactRepositoryAlterProps = objectProps{
	"INDEX_URL":             propString,
	"AUTHENTICATION_SECRET": propQName,
	"PACKAGE_RETENTION":     propOneOf("ON_OBJECT_CREATION"),
	"COMMENT":               propString,
}

// parseCreateImageRepository reads what follows the name of
//
//	CREATE [ OR REPLACE ] IMAGE REPOSITORY [ IF NOT EXISTS ] <name>
//	  [ ENCRYPTION = ( TYPE = 'SNOWFLAKE_FULL' | TYPE = 'SNOWFLAKE_SSE' ) ]
//	  [ COMMENT = '<string_literal>' ]
//	  [ [ WITH ] TAG ( <tag_name> = '<tag_value>' [ , ... ] ) ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-image-repository
func parseCreateImageRepository(p *Parser, s *ast.ObjectStmt) error {
	props, err := p.parseObjectProperties(s, "", false)
	if err != nil || props == nil {
		return err
	}
	// ENCRYPTION = ( TYPE = 'SNOWFLAKE_FULL' | TYPE = 'SNOWFLAKE_SSE' )
	for _, o := range props.Options {
		if o.Name != "ENCRYPTION" {
			continue
		}
		if len(o.Group) != 1 || o.Group[0].Name != "TYPE" || o.Group[0].Lit == nil ||
			o.Group[0].Lit.Kind != ast.LitString ||
			(o.Group[0].Lit.Value != "SNOWFLAKE_FULL" && o.Group[0].Lit.Value != "SNOWFLAKE_SSE") {
			return &ParseError{Position: o.Loc.Start, End: o.Loc.End, Message: "syntax error at or near ENCRYPTION"}
		}
	}
	return nil
}

// parseAlterImageRepository reads what follows the name of
//
//	ALTER IMAGE REPOSITORY [ IF EXISTS ] <name> SET COMMENT = '<string_literal>'
//	ALTER IMAGE REPOSITORY [ IF EXISTS ] <name> SET TAG <tag_name> = '<tag_value>' [ , ... ]
//	ALTER IMAGE REPOSITORY [ IF EXISTS ] <name> UNSET TAG <tag_name> [ , ... ]
//	ALTER IMAGE REPOSITORY [ IF EXISTS ] <name> UNSET COMMENT
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-image-repository
func parseAlterImageRepository(p *Parser, s *ast.ObjectStmt) error {
	ok, err := p.parseObjectAlterCommon(s, actSet|actTags|actUnset)
	if err != nil {
		return err
	}
	if !ok {
		return p.syntaxErrorAtCur()
	}
	return checkServiceSet(s)
}

// parseCreateArtifactRepository reads what follows the name of
//
//	CREATE [ OR REPLACE ] ARTIFACT REPOSITORY [ IF NOT EXISTS ] <name>
//	  TYPE = { APPLICATION | PYPI }
//	  [ API_INTEGRATION = '<api_integration_name>' ]
//	  [ INDEX_URL = '<url>' ] [ AUTHENTICATION_SECRET = <secret_name> ]
//	  [ PACKAGE_RETENTION = ON_OBJECT_CREATION ]
//	  [ [ WITH ] TAG ( <tag_name> = '<tag_value>' [ , ... ] ) ]
//	  [ COMMENT = '<string_literal>' ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-artifact-repository
func parseCreateArtifactRepository(p *Parser, s *ast.ObjectStmt) error {
	props, err := p.parseObjectProperties(s, "", true)
	if err != nil {
		return err
	}
	return p.requireObjectProperties(props, "TYPE")
}

// parseAlterArtifactRepository reads what follows the name of
//
//	ALTER ARTIFACT REPOSITORY [ IF EXISTS ] <name> SET
//	  [ COMMENT = '<string_literal>' ] [ INDEX_URL = '<url>' ]
//	  [ AUTHENTICATION_SECRET = <secret_name> ]
//	  [ PACKAGE_RETENTION = ON_OBJECT_CREATION ]
//	ALTER ARTIFACT REPOSITORY [ IF EXISTS ] <name> UNSET
//	  { COMMENT | AUTHENTICATION_SECRET | PACKAGE_RETENTION }
//	ALTER ARTIFACT REPOSITORY [ IF EXISTS ] <name> SET TAG <tag_name> = '<tag_value>' [ , ... ]
//	ALTER ARTIFACT REPOSITORY [ IF EXISTS ] <name> UNSET TAG <tag_name> [ , ... ]
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-artifact-repository
func parseAlterArtifactRepository(p *Parser, s *ast.ObjectStmt) error {
	ok, err := p.parseObjectAlterCommon(s, actSet|actTags|actUnset)
	if err != nil {
		return err
	}
	if !ok {
		return p.syntaxErrorAtCur()
	}
	// UNSET names one property.
	if c := s.Clauses[len(s.Clauses)-1]; c.Keywords == "UNSET" && len(c.Keys) != 1 {
		return &ParseError{Position: c.Loc.Start, End: c.Loc.End, Message: "syntax error: UNSET takes one property"}
	}
	return checkServiceSet(s)
}
