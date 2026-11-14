package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// POSTGRES INSTANCE (Snowflake Postgres), read as an ast.ObjectStmt
// ---------------------------------------------------------------------------

func init() {
	registerObject(&objectForm{
		words:  []string{"POSTGRES", "INSTANCE"},
		create: parseCreatePostgresInstance,
		alter:  parseAlterPostgresInstance,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-postgres-instance
		drop: true,
	})
}

// parseCreatePostgresInstance reads what follows the name of
//
//	CREATE POSTGRES INSTANCE [ IF NOT EXISTS ] <name>
//	  COMPUTE_FAMILY = '<compute_family>'
//	  STORAGE_SIZE_GB = <storage_gb>
//	  AUTHENTICATION_AUTHORITY = { POSTGRES | POSTGRES_OR_SNOWFLAKE | EXTERNAL_OAUTH }
//	  [ AUTHENTICATION_PROPERTIES = ( ISSUER = '<issuer>', AUDIENCE = '<audience>',
//	      ACCEPTED_TOKENS = ( ( MAPPING_CLAIM = '<claim>', SCOPES_CLAIM = '<claim>',
//	        SCOPES = { ( '<scope>' [ , ... ] ) | '<scope>' } ) [ , ... ] ) ) ]
//	  [ POSTGRES_VERSION = { 16 | 17 | 18 } ] [ NETWORK_POLICY = '<network_policy>' ]
//	  [ HIGH_AVAILABILITY = { TRUE | FALSE } ] [ STORAGE_INTEGRATION = '<integration>' ]
//	  [ POSTGRES_SETTINGS = '<json_string>' ] [ COMMENT = '<string>' ]
//	  [ [ WITH ] TAG ( <tag_name> = '<tag_value>' [ , ... ] ) ]
//
//	CREATE POSTGRES INSTANCE [ IF NOT EXISTS ] <name>
//	  FORK <original_instance_name>
//	  [ { AT | BEFORE } ( { TIMESTAMP => '<timestamp>' | OFFSET => <time_difference> } ) ]
//	  [ COMPUTE_FAMILY = '<compute_family>' ] [ STORAGE_SIZE_GB = <storage_gb> ]
//	  [ HIGH_AVAILABILITY = { TRUE | FALSE } ] [ POSTGRES_SETTINGS = '<json_string>' ]
//	  [ COMMENT = '<string>' ] [ [ WITH ] TAG ( <tag_name> = '<tag_value>' [ , ... ] ) ]
//
// The FORK clause records the original instance in Names; the AT or BEFORE
// clause records its TIMESTAMP or OFFSET as its one option.
//
// https://docs.snowflake.com/en/sql-reference/sql/create-postgres-instance
func parseCreatePostgresInstance(p *Parser, s *ast.ObjectStmt) error {
	if start := p.cur.Loc.Start; p.acceptWords("FORK") {
		src, err := p.parseIdent()
		if err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "FORK", start).Names = []*ast.ObjectName{{Name: src, Loc: src.Loc}}
		if p.cur.Type == kwAT || p.cur.Type == kwBEFORE {
			if err := parsePostgresInstanceForkPoint(p, s); err != nil {
				return err
			}
		}
		_, err = p.parseObjectProperties(s, "", false)
		return err
	}
	props, err := p.parseObjectProperties(s, "", true)
	if err != nil {
		return err
	}
	return p.requireObjectProperties(props, "COMPUTE_FAMILY", "STORAGE_SIZE_GB", "AUTHENTICATION_AUTHORITY")
}

// parsePostgresInstanceForkPoint reads a fork's point in time,
//
//	{ AT | BEFORE } ( { TIMESTAMP => '<timestamp>' | OFFSET => <time_difference> } )
//
// into a clause keyed AT or BEFORE whose one option is TIMESTAMP or OFFSET.
// cur is AT or BEFORE.
func parsePostgresInstanceForkPoint(p *Parser, s *ast.ObjectStmt) error {
	kw := p.advance() // consume AT or BEFORE
	c := &ast.ObjectClause{Keywords: "AT", Loc: ast.Loc{Start: kw.Loc.Start}}
	if kw.Type == kwBEFORE {
		c.Keywords = "BEFORE"
	}
	if _, err := p.expect('('); err != nil {
		return err
	}
	if p.cur.Type != kwTIMESTAMP && p.cur.Type != kwOFFSET {
		return p.syntaxErrorAtCur()
	}
	anchor := p.advance()
	opt := &ast.CopyOption{Name: "OFFSET", Loc: ast.Loc{Start: anchor.Loc.Start}}
	if anchor.Type == kwTIMESTAMP {
		opt.Name = "TIMESTAMP"
	}
	if _, err := p.expect(tokAssoc); err != nil {
		return err
	}
	if err := p.parseObjectValue(opt); err != nil {
		return err
	}
	if opt.Lit == nil {
		// A timestamp is a string and an offset a number of seconds.
		return &ParseError{Position: opt.Loc.Start, End: p.prev.Loc.End, Message: "syntax error: expected a literal"}
	}
	opt.Loc.End = p.prev.Loc.End
	if _, err := p.expect(')'); err != nil {
		return err
	}
	c.Options = []*ast.CopyOption{opt}
	c.Loc.End = p.prev.Loc.End
	s.Clauses = append(s.Clauses, c)
	return nil
}

// parseAlterPostgresInstance reads what follows the name of
//
//	ALTER POSTGRES INSTANCE [ IF EXISTS ] <name> RENAME TO <new_name>
//	ALTER POSTGRES INSTANCE [ IF EXISTS ] <name> SET
//	  [ NETWORK_POLICY = '<network_policy>' ]
//	  [ AUTHENTICATION_AUTHORITY = { POSTGRES | POSTGRES_OR_SNOWFLAKE | EXTERNAL_OAUTH } ]
//	  [ AUTHENTICATION_PROPERTIES = ( ... ) ] [ COMMENT = '<string>' ]
//	  [ HIGH_AVAILABILITY = { TRUE | FALSE } ] [ COMPUTE_FAMILY = '<compute_family>' ]
//	  [ STORAGE_SIZE_GB = <storage_gb> ] [ STORAGE_INTEGRATION = '<integration>' ]
//	  [ POSTGRES_VERSION = { 16 | 17 | 18 } ] [ MAINTENANCE_WINDOW_START = <hour> ]
//	  [ POSTGRES_SETTINGS = '<json_string>' ]
//	  [ APPLY { IMMEDIATELY | ON '<timestamp>' } ]
//	ALTER POSTGRES INSTANCE [ IF EXISTS ] <name> UNSET
//	  { COMMENT | POSTGRES_SETTINGS | NETWORK_POLICY | MAINTENANCE_WINDOW_START | STORAGE_INTEGRATION } [ , ... ]
//	ALTER POSTGRES INSTANCE [ IF EXISTS ] <name> { SUSPEND | RESUME }
//	ALTER POSTGRES INSTANCE [ IF EXISTS ] <name> RESET ACCESS FOR { 'snowflake_admin' | 'application' }
//	ALTER POSTGRES INSTANCE [ IF EXISTS ] <name> SET TAG <tag_name> = '<tag_value>' [ , ... ]
//	ALTER POSTGRES INSTANCE [ IF EXISTS ] <name> UNSET TAG <tag_name> [ , ... ]
//
// RESET ACCESS FOR and APPLY ON record their string as the clause's one
// option, with no name.
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-postgres-instance
func parseAlterPostgresInstance(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	switch {
	case p.acceptWords("SUSPEND"):
		p.addObjectKeywordsClause(s, "SUSPEND", start)
		return nil
	case p.acceptWords("RESUME"):
		p.addObjectKeywordsClause(s, "RESUME", start)
		return nil
	case p.acceptWords("RESET", "ACCESS", "FOR"):
		lit, err := parsePostgresInstanceString(p)
		if err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "RESET ACCESS FOR", start).Options = []*ast.CopyOption{lit}
		return nil
	case p.cur.Type == kwSET && p.peekNext().Type != kwTAG:
		p.advance() // consume SET
		c, err := p.parseObjectProperties(s, "SET", false)
		if err != nil {
			return err
		}
		if c != nil {
			c.Loc.Start = start
		} else {
			if !p.atWords("APPLY") {
				return p.syntaxErrorAtCur()
			}
			// SET APPLY ...: a SET with no property before its APPLY.
			p.addObjectKeywordsClause(s, "SET", start)
		}
		applyStart := p.cur.Loc.Start
		switch {
		case p.acceptWords("APPLY", "IMMEDIATELY"):
			p.addObjectKeywordsClause(s, "APPLY IMMEDIATELY", applyStart)
		case p.acceptWords("APPLY", "ON"):
			lit, err := parsePostgresInstanceString(p)
			if err != nil {
				return err
			}
			p.addObjectKeywordsClause(s, "APPLY ON", applyStart).Options = []*ast.CopyOption{lit}
		case p.atWords("APPLY"):
			p.advance() // consume APPLY
			return p.syntaxErrorAtCur()
		}
		return nil
	}
	ok, err := p.parseObjectAlterCommon(s, actUnset|actRename|actTags)
	if err != nil || ok {
		return err
	}
	return p.syntaxErrorAtCur()
}

// parsePostgresInstanceString reads a '<string>' into an option with no name.
func parsePostgresInstanceString(p *Parser) (*ast.CopyOption, error) {
	if p.cur.Type != tokString {
		return nil, p.syntaxErrorAtCur()
	}
	lit, err := p.parseCopyLiteral()
	if err != nil {
		return nil, err
	}
	return &ast.CopyOption{Lit: lit, Loc: lit.Loc}, nil
}
