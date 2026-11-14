package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// Native Apps: APPLICATION, APPLICATION ROLE and the ALTER APPLICATION
// subcommand pages, read as ast.ObjectStmt
// ---------------------------------------------------------------------------
//
// A clause keyword followed by a bare literal rather than `= value`
// (USING '<path>', PATCH <num>) is recorded as a clause of that keyword with a
// nameless option holding the literal, which deparse writes back as written.

func init() {
	registerObject(&objectForm{
		words:  []string{"APPLICATION"},
		create: parseCreateApplication,
		alter:  parseAlterApplication,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-application
		drop:          true,
		dropCascade:   true,
		nameParts:     1, // an account-level object
		noIfNotExists: true,
		props:         applicationProps,
		unset: []string{"COMMENT", "SHARE_EVENTS_WITH_PROVIDER", "DEBUG_MODE",
			"AUTHORIZE_RESTRICTED_PROVIDER_REMOTE_OPERATIONS_UNTIL"},
	})
	registerObject(&objectForm{
		words:  []string{"APPLICATION", "ROLE"},
		create: parseCreateApplicationRole,
		alter:  parseAlterApplicationRole,
		// https://docs.snowflake.com/en/sql-reference/sql/create-application-role
		orReplace: true,
		orAlter:   true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-application-role
		drop:  true,
		props: objectProps{"COMMENT": propString},
		// nameParts stays 0: the pages write the name QNAME (<app_name>.<role_name>).
	})
	// The pages below write ALTER APPLICATION with no application name: they
	// run in an app's setup script against the app itself. Each registers its
	// leading keywords as the object's, so the name read after them is the
	// specification's or the configuration's.
	registerObject(&objectForm{
		words:      []string{"APPLICATION", "SET", "SPECIFICATION"},
		alter:      parseAlterApplicationSetAppSpec,
		nameParts:  1,
		noIfExists: true,
		props:      applicationSpecProps,
	})
	registerObject(&objectForm{
		words:      []string{"APPLICATION", "DROP", "SPECIFICATION"},
		alter:      parseAlterApplicationDropNamed,
		nameParts:  1,
		noIfExists: true,
		props:      objectProps{},
	})
	registerObject(&objectForm{
		words:      []string{"APPLICATION", "SET", "CONFIGURATION", "DEFINITION"},
		alter:      parseAlterApplicationSetConfigurationDefinition,
		nameParts:  1,
		noIfExists: true,
		props:      applicationConfigurationProps,
	})
	registerObject(&objectForm{
		words:      []string{"APPLICATION", "DROP", "CONFIGURATION", "DEFINITION"},
		alter:      parseAlterApplicationDropNamed,
		nameParts:  1,
		noIfExists: true,
		props:      objectProps{},
	})
}

// applicationProps are the properties CREATE APPLICATION and ALTER
// APPLICATION ... SET document. DEBUG_MODE is CREATE's from a package's
// version or path and ALTER's; BACKGROUND_INSTALL is CREATE's from a listing;
// SHARE_EVENTS_WITH_PROVIDER is ALTER's.
var applicationProps = objectProps{
	"COMMENT":                           propString,
	"DEBUG_MODE":                        propBool,
	"SHARE_EVENTS_WITH_PROVIDER":        propBool,
	"BACKGROUND_INSTALL":                propBool,
	"AUTHORIZE_TELEMETRY_EVENT_SHARING": propBool,
	// '<timestamp>', 'INDEFINITE' or 'NEVER'
	"AUTHORIZE_RESTRICTED_PROVIDER_REMOTE_OPERATIONS_UNTIL": propString,
}

// applicationSpecProps are the properties of ALTER APPLICATION SET
// SPECIFICATION, over every TYPE the page documents.
var applicationSpecProps = objectProps{
	"TYPE":        propOneOf("EXTERNAL_ACCESS", "SECURITY_INTEGRATION", "LISTING", "CONNECTION", "SETTING"),
	"LABEL":       propString,
	"DESCRIPTION": propString,
	// TYPE = EXTERNAL_ACCESS
	"HOST_PORTS":         propStrings,
	"PRIVATE_HOST_PORTS": propStrings,
	// TYPE = SECURITY_INTEGRATION
	"OAUTH_TYPE":                   propString,
	"OAUTH_TOKEN_ENDPOINT":         propString,
	"OAUTH_AUTHORIZATION_ENDPOINT": propString,
	"OAUTH_ALLOWED_SCOPES":         propStrings,
	// TYPE = LISTING
	"TARGET_ACCOUNTS":                   propString,
	"LISTING":                           propName,
	"AUTO_FULFILLMENT_REFRESH_SCHEDULE": propString,
	// TYPE = CONNECTION
	"SERVER_APPLICATION":       propName,
	"SERVER_APPLICATION_ROLES": propNames,
	// TYPE = SETTING
	"SETTING": propName,
	"VALUE":   propString,
}

// applicationConfigurationProps are the properties of ALTER APPLICATION SET
// CONFIGURATION DEFINITION.
var applicationConfigurationProps = objectProps{
	"TYPE":              propOneOf("APPLICATION_NAME", "STRING", "SECRET_AUTHORIZATION"),
	"SECRET":            propQName, // <db>.<secret>
	"LABEL":             propString,
	"DESCRIPTION":       propString,
	"APPLICATION_ROLES": propNames,
	"SENSITIVE":         propBool,
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
	if _, err := p.parseApplicationNamedClause(s, keywords, start); err != nil {
		return err
	}

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
			p.advance() // consume USING
			if _, err := p.parseApplicationLiteralClause(s, "USING", start, tokString); err != nil {
				return err
			}
		}
	}
	if _, err := p.parseObjectProperties(s, "", false); err != nil {
		return err
	}
	// WITH FEATURE POLICY = <policy_name> is a WITH clause holding a
	// FEATURE POLICY option, so deparse writes its '=' back.
	if start := p.cur.Loc.Start; p.atWords("WITH", "FEATURE", "POLICY") {
		p.advance() // consume WITH
		opt := &ast.CopyOption{Name: "FEATURE POLICY", Loc: ast.Loc{Start: p.cur.Loc.Start}}
		p.advance() // consume FEATURE
		p.advance() // consume POLICY
		if _, err := p.expect('='); err != nil {
			return err
		}
		policy, err := p.parseObjectWord()
		if err != nil {
			return err
		}
		opt.Words = policy
		opt.Loc.End = p.prev.Loc.End
		p.addObjectKeywordsClause(s, "WITH", start).Options = []*ast.CopyOption{opt}
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
// USING VERSION, which start began, into a clause named keywords and a PATCH
// clause.
func (p *Parser) parseApplicationUsingVersion(s *ast.ObjectStmt, keywords string, start int) error {
	if _, err := p.parseApplicationNamedClause(s, keywords, start); err != nil {
		return err
	}
	if start := p.cur.Loc.Start; p.acceptWords("PATCH") {
		_, err := p.parseApplicationLiteralClause(s, "PATCH", start, tokInt)
		return err
	}
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
		opt := &ast.CopyOption{Loc: ast.Loc{Start: p.cur.Loc.Start}}
		if err := p.parseApplicationStringList(opt, 0); err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "SET SHARED TELEMETRY EVENTS", start).Options = []*ast.CopyOption{opt}
		return nil
	case p.acceptWords("UNSET", "REFERENCES"):
		c := p.addObjectKeywordsClause(s, "UNSET REFERENCES", start)
		if p.cur.Type == '(' {
			opt := &ast.CopyOption{Loc: ast.Loc{Start: p.cur.Loc.Start}}
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
		return p.parseApplicationRequiredOption(c, "VALUE", propString)
	case p.acceptWords("UNSET", "CONFIGURATION"):
		_, err := p.parseApplicationNamedClause(s, "UNSET CONFIGURATION", start)
		return err
	case p.acceptWords("APPROVE", "SPECIFICATION"):
		c, err := p.parseApplicationNamedClause(s, "APPROVE SPECIFICATION", start)
		if err != nil {
			return err
		}
		return p.parseApplicationRequiredOption(c, "SEQUENCE_NUMBER", propNumber)
	case p.acceptWords("DECLINE", "SPECIFICATION"):
		c, err := p.parseApplicationNamedClause(s, "DECLINE SPECIFICATION", start)
		if err != nil {
			return err
		}
		return p.parseApplicationRequiredOption(c, "SEQUENCE_NUMBER", propNumber)
	case p.acceptWords("UPGRADE"):
		if p.acceptWords("USING", "VERSION") {
			return p.parseApplicationUsingVersion(s, "UPGRADE USING VERSION", start)
		}
		p.addObjectKeywordsClause(s, "UPGRADE", start)
		if start := p.cur.Loc.Start; p.acceptWords("USING") {
			_, err := p.parseApplicationLiteralClause(s, "USING", start, tokString)
			return err
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
	policy, err := p.parseObjectName()
	if err != nil {
		return err
	}
	p.addObjectKeywordsClause(s, keywords, start).Names = []*ast.ObjectName{policy}
	if start := p.cur.Loc.Start; p.acceptWords("FORCE") {
		p.addObjectKeywordsClause(s, "FORCE", start)
	}
	return nil
}

// parseApplicationNamedClause reads the one name a clause takes after its
// keywords, which start began, and appends the clause. The pages write each
// such name (a version, channel, directive, specification, configuration,
// package or listing) as a single <name>.
func (p *Parser) parseApplicationNamedClause(s *ast.ObjectStmt, keywords string, start int) (*ast.ObjectClause, error) {
	name, err := p.parseObjectSimpleName()
	if err != nil {
		return nil, err
	}
	c := p.addObjectKeywordsClause(s, keywords, start)
	c.Names = []*ast.ObjectName{name}
	return c, nil
}

// parseApplicationRequiredOption reads the `<key> = <value>` property a clause
// requires next (VALUE = '<value>', SEQUENCE_NUMBER = <num>) into c, holding
// its value to shape v.
func (p *Parser) parseApplicationRequiredOption(c *ast.ObjectClause, key string, v objectValue) error {
	if !p.curIsWord(key) || !p.startsObjectProperty() {
		return p.syntaxErrorAtCur()
	}
	opt, err := p.parseObjectProperty()
	if err != nil {
		return err
	}
	if !v.matches(opt) {
		return applicationOptionError(opt)
	}
	c.Options = append(c.Options, opt)
	c.Loc.End = p.prev.Loc.End
	return nil
}

// parseApplicationOptions reads a run of `KEY = value` properties into c,
// where a clause's keywords and names come before its properties (MODIFY
// RELEASE DIRECTIVE <name> VERSION = ... PATCH = ...), holding each to props.
// An empty run is a syntax error.
func (p *Parser) parseApplicationOptions(c *ast.ObjectClause, props objectProps) error {
	if !p.startsObjectProperty() {
		return p.syntaxErrorAtCur()
	}
	for p.startsObjectProperty() {
		opt, err := p.parseObjectProperty()
		if err != nil {
			return err
		}
		if v, ok := props[opt.Name]; !ok || !v.matches(opt) {
			return applicationOptionError(opt)
		}
		c.Options = append(c.Options, opt)
	}
	c.Loc.End = p.prev.Loc.End
	return nil
}

// applicationOptionError reports a property a clause does not document, or
// whose value has another shape.
func applicationOptionError(o *ast.CopyOption) error {
	return &ParseError{Position: o.Loc.Start, End: o.Loc.End, Message: "syntax error at or near " + o.Name}
}

// parseApplicationLiteralClause reads the bare literal of type typ that
// follows a clause's keywords (USING '<path>', PATCH <num>), which the caller
// consumed from start, into a clause holding it as a nameless option.
func (p *Parser) parseApplicationLiteralClause(s *ast.ObjectStmt, keywords string, start, typ int) (*ast.ObjectClause, error) {
	if p.cur.Type != typ {
		return nil, p.syntaxErrorAtCur()
	}
	lit, err := p.parseCopyLiteral()
	if err != nil {
		return nil, err
	}
	c := p.addObjectKeywordsClause(s, keywords, start)
	c.Options = []*ast.CopyOption{{Lit: lit, Loc: lit.Loc}}
	return c, nil
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
	return nil
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
