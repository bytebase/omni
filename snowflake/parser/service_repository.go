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
		drop: true,
	})
	registerObject(&objectForm{
		words:     []string{"ARTIFACT", "REPOSITORY"},
		create:    parseCreateArtifactRepository,
		alter:     parseAlterArtifactRepository,
		orReplace: true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-artifact-repository
		drop: true,
	})
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
	_, err := p.parseObjectProperties(s, "", false)
	return err
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
	if err != nil || ok {
		return err
	}
	return p.syntaxErrorAtCur()
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
	if err != nil || ok {
		return err
	}
	return p.syntaxErrorAtCur()
}
