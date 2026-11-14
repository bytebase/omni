package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// Native Apps: APPLICATION, APPLICATION ROLE and the ALTER APPLICATION
// subcommand pages, read as ast.ObjectStmt
// ---------------------------------------------------------------------------
//
// A clause keyword followed by a bare literal rather than `= value`
// (USING '<path>', PATCH <num>) is recorded as an option named by the keyword,
// so a clause's values all sit in its Options.

func init() {
	registerObject(&objectForm{
		words:  []string{"APPLICATION"},
		create: parseCreateApplication,
		alter:  parseAlterApplication,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-application
		drop: true,
	})
	registerObject(&objectForm{
		words:  []string{"APPLICATION", "ROLE"},
		create: parseCreateApplicationRole,
		alter:  parseAlterApplicationRole,
		// https://docs.snowflake.com/en/sql-reference/sql/create-application-role
		orReplace: true,
		orAlter:   true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-application-role
		drop: true,
	})
	// The pages below write ALTER APPLICATION with no application name: they
	// run in an app's setup script against the app itself. Each registers its
	// leading keywords as the object's, so the name read after them is the
	// specification's or the configuration's.
	registerObject(&objectForm{
		words: []string{"APPLICATION", "SET", "SPECIFICATION"},
		alter: parseAlterApplicationSetAppSpec,
	})
	registerObject(&objectForm{
		words: []string{"APPLICATION", "DROP", "SPECIFICATION"},
		alter: parseAlterApplicationDropNamed,
	})
	registerObject(&objectForm{
		words: []string{"APPLICATION", "SET", "CONFIGURATION", "DEFINITION"},
		alter: parseAlterApplicationSetConfigurationDefinition,
	})
	registerObject(&objectForm{
		words: []string{"APPLICATION", "DROP", "CONFIGURATION", "DEFINITION"},
		alter: parseAlterApplicationDropNamed,
	})
}

// parseCreateApplication reads what follows the name of
//
//	CREATE APPLICATION <name> FROM APPLICATION PACKAGE <package_name>
//	  [ USING RELEASE CHANNEL { QA | ALPHA | DEFAULT }
//	  | USING '<path_to_version_directory>'
//	  | USING VERSION <version_identifier> [ PATCH <patch_num> ] ]
//	  [ DEBUG_MODE = { TRUE | FALSE } ] [ COMMENT = '<string_literal>' ]
//	  [ [ WITH ] TAG ( <tag_name> = '<tag_value>' [ , ... ] ) ]
//	  [ AUTHORIZE_TELEMETRY_EVENT_SHARING = { TRUE | FALSE } ]
//	  [ AUTHORIZE_RESTRICTED_PROVIDER_REMOTE_OPERATIONS_UNTIL = '<value>' ]
//	  [ WITH FEATURE POLICY = <policy_name> ]
//	CREATE APPLICATION <name> FROM LISTING <listing_name>
//	  [ USING RELEASE CHANNEL { QA | ALPHA | DEFAULT } ]
//	  [ COMMENT = ... ] [ [ WITH ] TAG ( ... ) ] [ BACKGROUND_INSTALL = { TRUE | FALSE } ]
//	  [ AUTHORIZE_TELEMETRY_EVENT_SHARING = ... ]
//	  [ AUTHORIZE_RESTRICTED_PROVIDER_REMOTE_OPERATIONS_UNTIL = ... ]
//	  [ WITH FEATURE POLICY = <policy_name> ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-application
func parseCreateApplication(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	keywords := "FROM APPLICATION PACKAGE"
	switch {
	case p.acceptWords("FROM", "APPLICATION", "PACKAGE"):
	case p.acceptWords("FROM", "LISTING"):
		keywords = "FROM LISTING"
	default:
		return p.syntaxErrorAtCur()
	}
	from, err := p.parseObjectName()
	if err != nil {
		return err
	}
	p.addObjectKeywordsClause(s, keywords, start).Names = []*ast.ObjectName{from}

	if start := p.cur.Loc.Start; p.curIsWord("USING") {
		switch {
		case p.acceptWords("USING", "RELEASE", "CHANNEL"):
			if err := p.parseApplicationReleaseChannel(s, start); err != nil {
				return err
			}
		case keywords == "FROM LISTING":
			return p.syntaxErrorAtCur()
		case p.acceptWords("USING", "VERSION"):
			if err := p.parseApplicationUsingVersion(s, "USING VERSION", start); err != nil {
				return err
			}
		default:
			c := &ast.ObjectClause{Keywords: "USING", Loc: ast.Loc{Start: start}}
			if err := p.parseApplicationLiteral(c, "USING", tokString); err != nil {
				return err
			}
			c.Loc.End = p.prev.Loc.End
			s.Clauses = append(s.Clauses, c)
		}
	}
	if _, err := p.parseObjectProperties(s, "", false); err != nil {
		return err
	}
	if start := p.cur.Loc.Start; p.acceptWords("WITH", "FEATURE", "POLICY") {
		if _, err := p.expect('='); err != nil {
			return err
		}
		policy, err := p.parseObjectName()
		if err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "WITH FEATURE POLICY", start).Names = []*ast.ObjectName{policy}
	}
	return nil
}

// parseApplicationReleaseChannel reads the { QA | ALPHA | DEFAULT } after
// USING RELEASE CHANNEL, which start began, into a clause of keywords alone.
func (p *Parser) parseApplicationReleaseChannel(s *ast.ObjectStmt, start int) error {
	for _, ch := range []string{"QA", "ALPHA", "DEFAULT"} {
		if p.acceptWords(ch) {
			p.addObjectKeywordsClause(s, "USING RELEASE CHANNEL "+ch, start)
			return nil
		}
	}
	return p.syntaxErrorAtCur()
}

// parseApplicationUsingVersion reads <version_identifier> [ PATCH <num> ] after
// USING VERSION, which start began, into a clause named keywords.
func (p *Parser) parseApplicationUsingVersion(s *ast.ObjectStmt, keywords string, start int) error {
	version, err := p.parseObjectName()
	if err != nil {
		return err
	}
	c := p.addObjectKeywordsClause(s, keywords, start)
	c.Names = []*ast.ObjectName{version}
	if p.curIsWord("PATCH") {
		if err := p.parseApplicationLiteral(c, "PATCH", tokInt); err != nil {
			return err
		}
	}
	c.Loc.End = p.prev.Loc.End
	return nil
}

// parseAlterApplication reads what follows the name of
//
//	ALTER APPLICATION [ IF EXISTS ] <name> SET [ COMMENT = '<string>' ]
//	  [ SHARE_EVENTS_WITH_PROVIDER = { TRUE | FALSE } ] [ DEBUG_MODE = { TRUE | FALSE } ]
//	  [ AUTHORIZE_RESTRICTED_PROVIDER_REMOTE_OPERATIONS_UNTIL = '<value>' ]
//	ALTER APPLICATION [ IF EXISTS ] <name> UNSET <property> [ , ... ]
//	ALTER APPLICATION [ IF EXISTS ] <name> RENAME TO <new_name>
//	ALTER APPLICATION <name> { SET | UNSET } { FEATURE | MAINTENANCE } POLICY ...
//	ALTER APPLICATION <name> SET TAG <tag_name> = '<tag_value>' [ , ... ]
//	ALTER APPLICATION <name> UNSET TAG <tag_name> [ , ... ]
//	ALTER APPLICATION <name> SET SHARED TELEMETRY EVENTS ( '<event_definition>' [ , ... ] )
//	ALTER APPLICATION <name> SET AUTHORIZE_TELEMETRY_EVENT_SHARING = { TRUE | FALSE }
//	ALTER APPLICATION <name> UNSET REFERENCES [ ( '<reference_name>' [ , '<reference_alias>' ] ) ]
//	ALTER APPLICATION <name> UPGRADE
//	  [ USING VERSION <version_identifier> [ PATCH <patch_num> ] | USING '<path>' ]
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-application
//
// and the subcommand pages that name the application:
//
//	ALTER APPLICATION <name> { APPROVE | DECLINE } SPECIFICATION <app_spec_name>
//	  SEQUENCE_NUMBER = <sequence_number>
//	ALTER APPLICATION <name> SET CONFIGURATION <config> VALUE = '<value>'
//	ALTER APPLICATION <name> UNSET CONFIGURATION <config>
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-application-sequence-number
// https://docs.snowflake.com/en/sql-reference/sql/alter-application-set-configuration-value
// https://docs.snowflake.com/en/sql-reference/sql/alter-application-unset-configuration
func parseAlterApplication(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	switch {
	case p.acceptWords("SET", "FEATURE", "POLICY"):
		return p.parseApplicationPolicy(s, "SET FEATURE POLICY", start)
	case p.acceptWords("SET", "MAINTENANCE", "POLICY"):
		return p.parseApplicationPolicy(s, "SET MAINTENANCE POLICY", start)
	case p.acceptWords("UNSET", "FEATURE", "POLICY"):
		p.addObjectKeywordsClause(s, "UNSET FEATURE POLICY", start)
		return nil
	case p.acceptWords("UNSET", "MAINTENANCE", "POLICY"):
		p.addObjectKeywordsClause(s, "UNSET MAINTENANCE POLICY", start)
		return nil
	case p.acceptWords("SET", "SHARED", "TELEMETRY", "EVENTS"):
		opt := &ast.CopyOption{Name: "EVENTS", Loc: ast.Loc{Start: p.cur.Loc.Start}}
		if err := p.parseApplicationStringList(opt, 0); err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "SET SHARED TELEMETRY EVENTS", start).Options = []*ast.CopyOption{opt}
		return nil
	case p.acceptWords("UNSET", "REFERENCES"):
		c := p.addObjectKeywordsClause(s, "UNSET REFERENCES", start)
		if p.cur.Type == '(' {
			opt := &ast.CopyOption{Name: "REFERENCES", Loc: ast.Loc{Start: p.cur.Loc.Start}}
			if err := p.parseApplicationStringList(opt, 2); err != nil {
				return err
			}
			c.Options = []*ast.CopyOption{opt}
			c.Loc.End = p.prev.Loc.End
		}
		return nil
	case p.acceptWords("SET", "CONFIGURATION"):
		c, err := p.parseApplicationNamedClause(s, "SET CONFIGURATION", start)
		if err != nil {
			return err
		}
		return p.parseApplicationRequiredOption(c, "VALUE")
	case p.acceptWords("UNSET", "CONFIGURATION"):
		_, err := p.parseApplicationNamedClause(s, "UNSET CONFIGURATION", start)
		return err
	case p.acceptWords("APPROVE", "SPECIFICATION"):
		c, err := p.parseApplicationNamedClause(s, "APPROVE SPECIFICATION", start)
		if err != nil {
			return err
		}
		return p.parseApplicationRequiredOption(c, "SEQUENCE_NUMBER")
	case p.acceptWords("DECLINE", "SPECIFICATION"):
		c, err := p.parseApplicationNamedClause(s, "DECLINE SPECIFICATION", start)
		if err != nil {
			return err
		}
		return p.parseApplicationRequiredOption(c, "SEQUENCE_NUMBER")
	case p.acceptWords("UPGRADE"):
		upgradeEnd := p.prev.Loc.End
		if p.acceptWords("USING", "VERSION") {
			return p.parseApplicationUsingVersion(s, "UPGRADE USING VERSION", start)
		}
		c := p.addObjectKeywordsClause(s, "UPGRADE", start)
		c.Loc.End = upgradeEnd
		if p.curIsWord("USING") {
			if err := p.parseApplicationLiteral(c, "USING", tokString); err != nil {
				return err
			}
			c.Loc.End = p.prev.Loc.End
		}
		return nil
	}
	ok, err := p.parseObjectAlterCommon(s, actSet|actUnset|actRename|actTags)
	if err != nil || ok {
		return err
	}
	return p.syntaxErrorAtCur()
}

// parseApplicationPolicy reads <policy_name> [ FORCE ] after SET FEATURE
// POLICY or SET MAINTENANCE POLICY, which start began.
func (p *Parser) parseApplicationPolicy(s *ast.ObjectStmt, keywords string, start int) error {
	if _, err := p.parseApplicationNamedClause(s, keywords, start); err != nil {
		return err
	}
	if start := p.cur.Loc.Start; p.acceptWords("FORCE") {
		p.addObjectKeywordsClause(s, "FORCE", start)
	}
	return nil
}

// parseApplicationNamedClause reads the one name a clause takes after its
// keywords, which start began, and appends the clause.
func (p *Parser) parseApplicationNamedClause(s *ast.ObjectStmt, keywords string, start int) (*ast.ObjectClause, error) {
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	c := p.addObjectKeywordsClause(s, keywords, start)
	c.Names = []*ast.ObjectName{name}
	return c, nil
}

// parseApplicationRequiredOption reads the `<key> = <value>` property a clause
// requires next (VALUE = '<value>', SEQUENCE_NUMBER = <num>) into c.
func (p *Parser) parseApplicationRequiredOption(c *ast.ObjectClause, key string) error {
	if !p.curIsWord(key) || !p.startsObjectProperty() {
		return p.syntaxErrorAtCur()
	}
	opt, err := p.parseObjectProperty()
	if err != nil {
		return err
	}
	c.Options = append(c.Options, opt)
	c.Loc.End = p.prev.Loc.End
	return nil
}

// parseApplicationOptions reads a run of `KEY = value` properties into c,
// where a clause's keywords and names come before its properties (MODIFY
// RELEASE DIRECTIVE <name> VERSION = ... PATCH = ...). need says whether an
// empty run is a syntax error.
func (p *Parser) parseApplicationOptions(c *ast.ObjectClause, need bool) error {
	if need && !p.startsObjectProperty() {
		return p.syntaxErrorAtCur()
	}
	for p.startsObjectProperty() {
		opt, err := p.parseObjectProperty()
		if err != nil {
			return err
		}
		c.Options = append(c.Options, opt)
	}
	c.Loc.End = p.prev.Loc.End
	return nil
}

// parseApplicationLiteral reads <word> <literal>, a clause keyword followed by
// a bare literal of type typ (USING '<path>', PATCH <num>), into an option of
// c named by the word. cur is the word.
func (p *Parser) parseApplicationLiteral(c *ast.ObjectClause, word string, typ int) error {
	if !p.curIsWord(word) {
		return p.syntaxErrorAtCur()
	}
	key := p.advance()
	if p.cur.Type != typ {
		return p.syntaxErrorAtCur()
	}
	lit, err := p.parseCopyLiteral()
	if err != nil {
		return err
	}
	c.Options = append(c.Options, &ast.CopyOption{Name: word, Lit: lit, Loc: ast.Loc{Start: key.Loc.Start, End: lit.Loc.End}})
	return nil
}

// parseApplicationStringList reads ( '<string>' [ , ... ] ) into opt.List,
// holding at most max strings when max is positive. cur is '('.
func (p *Parser) parseApplicationStringList(opt *ast.CopyOption, max int) error {
	if _, err := p.expect('('); err != nil {
		return err
	}
	for {
		if p.cur.Type != tokString || (max > 0 && len(opt.List) == max) {
			return p.syntaxErrorAtCur()
		}
		lit, err := p.parseCopyLiteral()
		if err != nil {
			return err
		}
		opt.List = append(opt.List, lit)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	if _, err := p.expect(')'); err != nil {
		return err
	}
	opt.Loc.End = p.prev.Loc.End
	return nil
}

// rejectApplicationIfExists fails an ALTER whose page does not document
// IF EXISTS, as the pages that name no application do not.
func rejectApplicationIfExists(s *ast.ObjectStmt) error {
	if !s.IfExists {
		return nil
	}
	return &ParseError{
		Position: s.Name.Loc.Start, End: s.Name.Loc.End,
		Message: "syntax error: IF EXISTS is not allowed in ALTER " + s.Kind,
	}
}

// parseAlterApplicationSetAppSpec reads what follows the specification's name
// in
//
//	ALTER APPLICATION SET SPECIFICATION <app_spec_name>
//	  TYPE = { EXTERNAL_ACCESS | SECURITY_INTEGRATION | LISTING | CONNECTION | SETTING }
//	  LABEL = '<label>' DESCRIPTION = '<description>' <type_properties>
//
// where the properties of each TYPE follow it (HOST_PORTS = ( ... ),
// OAUTH_TYPE = '...', TARGET_ACCOUNTS = '...', SERVER_APPLICATION = <name>,
// SETTING = <name>, ...).
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-application-set-app-spec
func parseAlterApplicationSetAppSpec(p *Parser, s *ast.ObjectStmt) error {
	if err := rejectApplicationIfExists(s); err != nil {
		return err
	}
	props, err := p.parseObjectProperties(s, "", true)
	if err != nil {
		return err
	}
	return p.requireObjectProperties(props, "TYPE", "LABEL", "DESCRIPTION")
}

// parseAlterApplicationSetConfigurationDefinition reads what follows the
// configuration's name in
//
//	ALTER APPLICATION SET CONFIGURATION DEFINITION <config>
//	  TYPE = { APPLICATION_NAME | STRING | SECRET_AUTHORIZATION } [ SECRET = <secret_name> ]
//	  LABEL = '<label>' DESCRIPTION = '<description>'
//	  APPLICATION_ROLES = ( <app_role> [ , ... ] ) [ SENSITIVE = { TRUE | FALSE } ]
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-application-set-configuration-definition
func parseAlterApplicationSetConfigurationDefinition(p *Parser, s *ast.ObjectStmt) error {
	if err := rejectApplicationIfExists(s); err != nil {
		return err
	}
	props, err := p.parseObjectProperties(s, "", true)
	if err != nil {
		return err
	}
	return p.requireObjectProperties(props, "TYPE", "LABEL", "DESCRIPTION", "APPLICATION_ROLES")
}

// parseAlterApplicationDropNamed reads the end of
//
//	ALTER APPLICATION DROP SPECIFICATION <app_spec_name>
//	ALTER APPLICATION DROP CONFIGURATION DEFINITION <config>
//
// which take nothing after the name.
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-application-drop-app-spec
// https://docs.snowflake.com/en/sql-reference/sql/alter-application-drop-configuration-definition
func parseAlterApplicationDropNamed(p *Parser, s *ast.ObjectStmt) error {
	return rejectApplicationIfExists(s)
}

// parseCreateApplicationRole reads what follows the name of
//
//	CREATE [ OR REPLACE ] APPLICATION ROLE [ IF NOT EXISTS ] <name>
//	  [ COMMENT = '<string_literal>' ]
//	CREATE OR ALTER APPLICATION ROLE <name> [ COMMENT = '<string_literal>' ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-application-role
func parseCreateApplicationRole(p *Parser, s *ast.ObjectStmt) error {
	_, err := p.parseObjectProperties(s, "", false)
	return err
}

// parseAlterApplicationRole reads what follows the name of
//
//	ALTER APPLICATION ROLE [ IF EXISTS ] <name> RENAME TO <new_name>
//	ALTER APPLICATION ROLE [ IF EXISTS ] <name> SET COMMENT = '<string_literal>'
//	ALTER APPLICATION ROLE [ IF EXISTS ] <name> UNSET COMMENT
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-application-role
func parseAlterApplicationRole(p *Parser, s *ast.ObjectStmt) error {
	ok, err := p.parseObjectAlterCommon(s, actSet|actUnset|actRename)
	if err != nil || ok {
		return err
	}
	return p.syntaxErrorAtCur()
}

// startsShowApplicationPackageScope reports whether cur, after a SHOW's
// IN APPLICATION, is PACKAGE followed by the package's name. It consumes
// nothing.
func (p *Parser) startsShowApplicationPackageScope() bool {
	return p.curIsWord("PACKAGE") && p.nextStartsShowScopeName()
}

// startsShowListingScope reports whether cur, after a SHOW's IN, is LISTING
// followed by the listing's name. It consumes nothing.
func (p *Parser) startsShowListingScope() bool {
	return p.curIsWord("LISTING") && p.nextStartsShowScopeName()
}

// nextStartsShowScopeName reports whether the token after cur can begin a
// SHOW scope's name. It consumes nothing.
func (p *Parser) nextStartsShowScopeName() bool {
	snap := p.snapshot()
	defer p.restore(snap)
	p.advance()
	return p.startsScopeName()
}
