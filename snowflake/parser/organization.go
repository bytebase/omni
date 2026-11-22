package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Organization objects, read as ast.ObjectStmt: ORGANIZATION USER,
// ORGANIZATION USER GROUP, ORGANIZATION ACCOUNT, ORGANIZATION PROFILE, and
// ALTER ORGANIZATION itself
// ---------------------------------------------------------------------------

func init() {
	registerObject(&objectForm{
		words: []string{"ORGANIZATION"},
		alter: parseAlterOrganization,
		// ALTER ORGANIZATION names no object and takes no IF EXISTS.
		alterNoName: true,
		noIfExists:  true,
		props:       organizationProps,
	})
	registerObject(&objectForm{
		words:  []string{"ORGANIZATION", "USER"},
		create: parseCreateOrganizationUser,
		alter:  parseAlterOrganizationUser,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-organization-user
		drop:       true,
		nameParts:  1,
		props:      organizationUserProps,
		alterProps: organizationUserAlterProps,
		unset:      []string{"EMAIL", "DISPLAY_NAME", "FIRST_NAME", "MIDDLE_NAME", "LAST_NAME", "TYPE", "COMMENT"},
	})
	registerObject(&objectForm{
		words:  []string{"ORGANIZATION", "USER", "GROUP"},
		create: parseCreateOrganizationUserGroup,
		alter:  parseAlterOrganizationUserGroup,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-organization-user-group
		drop:      true,
		nameParts: 1,
		props:     objectProps{"IS_GRANTABLE": propBool},
		unset:     []string{},
	})
	registerObject(&objectForm{
		words:  []string{"ORGANIZATION", "ACCOUNT"},
		create: parseCreateOrganizationAccount,
		alter:  parseAlterOrganizationAccount,
		// ALTER ORGANIZATION ACCOUNT acts on the current account: it names
		// none and takes no IF EXISTS. CREATE writes no IF NOT EXISTS. No
		// props: ALTER ... SET takes the account parameters ALTER ACCOUNT
		// does, so CREATE's properties are checked by the create function.
		alterNoName:   true,
		noIfExists:    true,
		noIfNotExists: true,
		nameParts:     1,
	})
	registerObject(&objectForm{
		words:  []string{"ORGANIZATION", "PROFILE"},
		create: parseCreateOrganizationProfile,
		alter:  parseAlterOrganizationProfile,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-organization-profile
		drop:           true,
		dropNoIfExists: true,
		nameParts:      1,
		props:          objectProps{"PUBLISH": propBool},
		unset:          []string{},
	})
}

// organizationProps are the properties ALTER ORGANIZATION ... SET documents.
var organizationProps = objectProps{
	"DEFAULT_TENANT_TYPE": propOneOf("INTERNAL", "EXTERNAL"),
	"DOMAIN_NAMES":        propStrings,
}

// organizationUserProps are CREATE ORGANIZATION USER's properties.
var organizationUserProps = objectProps{
	"EMAIL":        propString,
	"LOGIN_NAME":   propString,
	"DISPLAY_NAME": propString,
	"FIRST_NAME":   propString,
	"MIDDLE_NAME":  propString,
	"LAST_NAME":    propString,
	"TYPE":         propOneOf("PERSON", "SERVICE"),
	"COMMENT":      propString,
}

// organizationUserAlterProps are ALTER ORGANIZATION USER ... SET's and
// UNSET's: CREATE's but LOGIN_NAME.
var organizationUserAlterProps = func() objectProps {
	m := objectProps{}
	for k, v := range organizationUserProps {
		if k != "LOGIN_NAME" {
			m[k] = v
		}
	}
	return m
}()

// organizationAccountProps are CREATE ORGANIZATION ACCOUNT's properties.
var organizationAccountProps = objectProps{
	"ADMIN_NAME":           propString,
	"ADMIN_PASSWORD":       propString,
	"ADMIN_RSA_PUBLIC_KEY": propString,
	"FIRST_NAME":           propString,
	"LAST_NAME":            propString,
	"EMAIL":                propString,
	"MUST_CHANGE_PASSWORD": propBool,
	"EDITION":              propOneOf("ENTERPRISE", "BUSINESS_CRITICAL"),
	"REGION_GROUP":         propName,
	"REGION":               propName,
	"COMMENT":              propString,
}

// parseCreateOrganizationUser reads what follows the name of
//
//	CREATE ORGANIZATION USER [ IF NOT EXISTS ] <name>
//	  [ EMAIL = '<string>' ] [ LOGIN_NAME = '<string>' ] [ DISPLAY_NAME = '<string>' ]
//	  [ FIRST_NAME = '<string>' ] [ MIDDLE_NAME = '<string>' ] [ LAST_NAME = '<string>' ]
//	  [ TYPE = { PERSON | SERVICE } ] [ COMMENT = '<string>' ]
//
// Duplicate property validation is deferred to the semantic/catalog layer,
// matching CREATE WAREHOUSE.
//
// https://docs.snowflake.com/en/sql-reference/sql/create-organization-user
func parseCreateOrganizationUser(p *Parser, s *ast.ObjectStmt) error {
	if _, err := p.parseObjectPropertiesWith(s, "", false, true); err != nil {
		return err
	}
	return rejectObjectTags(s)
}

// parseAlterOrganizationUser reads what follows the name of
//
//	ALTER ORGANIZATION USER [ IF EXISTS ] <name> SET <property> [ ... ]
//	ALTER ORGANIZATION USER <name> UNSET <property> [ , ... ]
//
// Duplicate property validation is deferred to the semantic/catalog layer,
// matching CREATE/ALTER WAREHOUSE.
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-organization-user
func parseAlterOrganizationUser(p *Parser, s *ast.ObjectStmt) error {
	if s.IfExists && p.cur.Type == kwUNSET {
		// The page writes IF EXISTS on SET alone.
		return p.syntaxErrorAtCur()
	}
	if p.cur.Type == kwSET {
		start := p.cur.Loc.Start
		p.advance() // consume SET
		c, err := p.parseObjectPropertiesWith(s, "SET", true, true)
		if err != nil {
			return err
		}
		c.Loc.Start = start
		return rejectObjectTags(s)
	}
	ok, err := p.parseObjectAlterCommon(s, actUnset)
	if err != nil {
		return err
	}
	if !ok {
		return p.syntaxErrorAtCur()
	}
	return rejectObjectTags(s)
}

// parseCreateOrganizationUserGroup reads what follows the name of
//
//	CREATE ORGANIZATION USER GROUP [ IF NOT EXISTS ] <name>
//	  [ IS_GRANTABLE = { TRUE | FALSE } ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-organization-user-group
func parseCreateOrganizationUserGroup(p *Parser, s *ast.ObjectStmt) error {
	if _, err := p.parseObjectProperties(s, "", false); err != nil {
		return err
	}
	return rejectObjectTags(s)
}

// parseAlterOrganizationUserGroup reads what follows the name of
//
//	ALTER ORGANIZATION USER GROUP [ IF EXISTS ] <name>
//	  { ADD | REMOVE } ORGANIZATION USERS <user_name> [ , ... ]
//	ALTER ORGANIZATION USER GROUP [ IF EXISTS ] <name>
//	  SET VISIBILITY = { ALL | ACCOUNTS <account_name> [ , ... ]
//	                   | REGION GROUPS '<region_group>' [ , ... ] }
//	ALTER ORGANIZATION USER GROUP [ IF EXISTS ] <name>
//	  SET IS_GRANTABLE = { TRUE | FALSE }
//
// SET VISIBILITY's value is a clause of its own: SET VISIBILITY = ALL,
// SET VISIBILITY = ACCOUNTS with the accounts' names, or SET VISIBILITY =
// REGION GROUPS with each region group a nameless option.
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-organization-user-group
func parseAlterOrganizationUserGroup(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	for _, verb := range []string{"ADD", "REMOVE"} {
		if p.acceptWords(verb, "ORGANIZATION", "USERS") {
			names, err := p.parseObjectSimpleNameList()
			if err != nil {
				return err
			}
			p.addObjectKeywordsClause(s, verb+" ORGANIZATION USERS", start).Names = names
			return nil
		}
	}
	if p.cur.Type == kwSET && isWordToken(p.peekNext(), "VISIBILITY") {
		p.advance() // consume SET
		p.advance() // consume VISIBILITY
		if _, err := p.expect('='); err != nil {
			return err
		}
		switch {
		case p.acceptWords("ALL"):
			p.addObjectKeywordsClause(s, "SET VISIBILITY = ALL", start)
		case p.acceptWords("ACCOUNTS"):
			names, err := p.parseObjectSimpleNameList()
			if err != nil {
				return err
			}
			p.addObjectKeywordsClause(s, "SET VISIBILITY = ACCOUNTS", start).Names = names
		case p.acceptWords("REGION", "GROUPS"):
			groups, err := p.parseObjectStringList()
			if err != nil {
				return err
			}
			p.addObjectKeywordsClause(s, "SET VISIBILITY = REGION GROUPS", start).Options = groups
		default:
			return p.syntaxErrorAtCur()
		}
		return nil
	}
	if p.cur.Type == kwSET {
		ok, err := p.parseObjectAlterCommon(s, actSet)
		if err != nil || !ok {
			return err
		}
		return rejectObjectTags(s)
	}
	return p.syntaxErrorAtCur()
}

// parseCreateOrganizationAccount reads what follows the name of
//
//	CREATE ORGANIZATION ACCOUNT <name>
//	  ADMIN_NAME = '<string>'
//	  { ADMIN_PASSWORD = '<string>' | ADMIN_RSA_PUBLIC_KEY = '<string>' }
//	  [ FIRST_NAME = '<string>' ] [ LAST_NAME = '<string>' ]
//	  EMAIL = '<string>' [ MUST_CHANGE_PASSWORD = { TRUE | FALSE } ]
//	  EDITION = { ENTERPRISE | BUSINESS_CRITICAL }
//	  [ REGION_GROUP = <region_group_id> ] [ REGION = <snowflake_region_id> ]
//	  [ COMMENT = '<string>' ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-organization-account
func parseCreateOrganizationAccount(p *Parser, s *ast.ObjectStmt) error {
	props, err := p.parseObjectProperties(s, "", true)
	if err != nil {
		return err
	}
	if err := rejectObjectTags(s); err != nil {
		return err
	}
	for _, o := range props.Options {
		v, ok := organizationAccountProps[o.Name]
		// ADMIN_NAME is a '<string>' on the reference page, and an unquoted
		// name in the guide's example (ADMIN_NAME = admin).
		// https://docs.snowflake.com/en/user-guide/organization-accounts
		if o.Name == "ADMIN_NAME" && propName.matches(o) {
			continue
		}
		if !ok || !v.matches(o) {
			return &ParseError{Position: o.Loc.Start, End: o.Loc.End, Message: "syntax error at or near " + o.Name}
		}
	}
	if err := p.requireObjectProperties(props, "ADMIN_NAME", "EMAIL", "EDITION"); err != nil {
		return err
	}
	// Exactly one of ADMIN_PASSWORD and ADMIN_RSA_PUBLIC_KEY.
	n := 0
	for _, o := range props.Options {
		if o.Name == "ADMIN_PASSWORD" || o.Name == "ADMIN_RSA_PUBLIC_KEY" {
			n++
		}
	}
	if n != 1 {
		return &ParseError{
			Position: props.Loc.Start, End: props.Loc.End,
			Message: "syntax error: one of ADMIN_PASSWORD and ADMIN_RSA_PUBLIC_KEY is required",
		}
	}
	return nil
}

// parseAlterOrganizationAccount reads what follows ORGANIZATION ACCOUNT in
//
//	ALTER ORGANIZATION ACCOUNT SET { [ accountParams ] | [ objectParams ] | [ sessionParams ] }
//	ALTER ORGANIZATION ACCOUNT UNSET <param_name> [ , ... ]
//	ALTER ORGANIZATION ACCOUNT SET RESOURCE_MONITOR = <monitor_name>
//	ALTER ORGANIZATION ACCOUNT SET { PASSWORD | SESSION } POLICY <policy_name>
//	ALTER ORGANIZATION ACCOUNT UNSET { PASSWORD | SESSION } POLICY
//	ALTER ORGANIZATION ACCOUNT SET TAG <tag_name> = '<tag_value>' [ , ... ]
//	ALTER ORGANIZATION ACCOUNT UNSET TAG <tag_name> [ , ... ]
//
// The parameters are ALTER ACCOUNT's, which the page refers to.
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-organization-account
func parseAlterOrganizationAccount(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	if (p.cur.Type == kwSET || p.cur.Type == kwUNSET) &&
		(p.peekNext().Type == kwPASSWORD || p.peekNext().Type == kwSESSION) {
		verb := strings.ToUpper(p.cur.Str)
		snap := p.snapshot()
		p.advance() // consume SET or UNSET
		kind := strings.ToUpper(p.advance().Str)
		if p.cur.Type == kwPOLICY {
			p.advance() // consume POLICY
			c := p.addObjectKeywordsClause(s, verb+" "+kind+" POLICY", start)
			if verb == "SET" {
				name, err := p.parseObjectName()
				if err != nil {
					return err
				}
				c.Names = []*ast.ObjectName{name}
				c.Loc.End = p.prev.Loc.End
			}
			return nil
		}
		// A parameter named PASSWORD or SESSION.
		p.restore(snap)
	}
	if p.cur.Type == kwSET && p.peekNext().Type != kwTAG {
		return p.parseOrganizationAccountParams(s)
	}
	ok, err := p.parseObjectAlterCommon(s, actUnset|actTags)
	if err != nil {
		return err
	}
	if !ok {
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseOrganizationAccountParams reads SET <param> = <value> [ ... ], whose
// parameters the page lets blank spaces, commas or new lines separate. cur
// is SET.
func (p *Parser) parseOrganizationAccountParams(s *ast.ObjectStmt) error {
	c := &ast.ObjectClause{Keywords: "SET", Loc: ast.Loc{Start: p.cur.Loc.Start}}
	p.advance() // consume SET
	seen := map[string]bool{}
	for {
		if !p.startsObjectProperty() {
			return p.syntaxErrorAtCur()
		}
		opt, err := p.parseObjectProperty()
		if err != nil {
			return err
		}
		if seen[opt.Name] {
			return &ParseError{Position: opt.Loc.Start, End: opt.Loc.End, Message: "syntax error: " + opt.Name + " is set twice"}
		}
		seen[opt.Name] = true
		c.Options = append(c.Options, opt)
		if p.cur.Type == ',' {
			p.advance() // consume ','; a parameter must follow
			continue
		}
		if !p.startsObjectProperty() {
			break
		}
	}
	c.Loc.End = p.prev.Loc.End
	s.Clauses = append(s.Clauses, c)
	return nil
}

// parseCreateOrganizationProfile reads what follows the name of
//
//	CREATE ORGANIZATION PROFILE [ IF NOT EXISTS ] <name>
//	CREATE ORGANIZATION PROFILE [ IF NOT EXISTS ] <name> AS '<yaml_manifest_string>'
//	  [ VERSION <version_alias_name> ] [ PUBLISH = { TRUE | FALSE } ]
//	CREATE ORGANIZATION PROFILE [ IF NOT EXISTS ] <name> FROM @<yaml_manifest_stage_location>
//	  [ VERSION <version_alias_name> ] [ PUBLISH = { TRUE | FALSE } ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-organization-profile
func parseCreateOrganizationProfile(p *Parser, s *ast.ObjectStmt) error {
	switch {
	case p.cur.Type == kwAS:
		if err := p.parseObjectAsString(s); err != nil {
			return err
		}
	case p.cur.Type == kwFROM:
		if err := p.parseObjectFromStage(s); err != nil {
			return err
		}
	default:
		return nil
	}
	if start := p.cur.Loc.Start; p.cur.Type == kwVERSION {
		p.advance() // consume VERSION
		name, err := p.parseObjectSimpleName()
		if err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "VERSION", start).Names = []*ast.ObjectName{name}
	}
	if _, err := p.parseObjectProperties(s, "", false); err != nil {
		return err
	}
	return rejectObjectTags(s)
}

// parseAlterOrganizationProfile reads what follows the name of
//
//	ALTER ORGANIZATION PROFILE [ IF EXISTS ] <name> AS '<yaml_manifest_string>'
//	ALTER ORGANIZATION PROFILE [ IF EXISTS ] <name> RENAME TO <new_name>
//	ALTER ORGANIZATION PROFILE [ IF EXISTS ] <name> PUBLISH
//	ALTER ORGANIZATION PROFILE <name> ADD VERSION [ [ IF NOT EXISTS ] <version_name> ]
//	  FROM @<yaml_manifest_stage_location>
//	ALTER ORGANIZATION PROFILE <name> ADD LIVE VERSION [ [ IF NOT EXISTS ] <version_name> ] FROM LAST
//	ALTER ORGANIZATION PROFILE <name> { COMMIT | ABORT }
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-organization-profile
func parseAlterOrganizationProfile(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	switch {
	case p.cur.Type == kwAS:
		return p.parseObjectAsString(s)
	case p.cur.Type == kwRENAME:
		return p.parseObjectRenameSimple(s)
	case p.acceptWords("PUBLISH"):
		p.addObjectKeywordsClause(s, "PUBLISH", start)
		return nil
	}
	if s.IfExists {
		// The page writes IF EXISTS on AS, RENAME TO and PUBLISH alone.
		return p.syntaxErrorAtCur()
	}
	switch {
	case p.acceptWords("COMMIT"):
		p.addObjectKeywordsClause(s, "COMMIT", start)
		return nil
	case p.acceptWords("ABORT"):
		p.addObjectKeywordsClause(s, "ABORT", start)
		return nil
	case p.acceptWords("ADD", "VERSION"):
		if err := p.parseObjectAddVersion(s, "ADD VERSION", start); err != nil {
			return err
		}
		return p.parseObjectFromStage(s)
	case p.acceptWords("ADD", "LIVE", "VERSION"):
		if err := p.parseObjectAddVersion(s, "ADD LIVE VERSION", start); err != nil {
			return err
		}
		fromStart := p.cur.Loc.Start
		if err := p.expectWords("FROM", "LAST"); err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "FROM LAST", fromStart)
		return nil
	}
	return p.syntaxErrorAtCur()
}

// parseAlterOrganization reads what follows ALTER ORGANIZATION in
//
//	ALTER ORGANIZATION SET DEFAULT_TENANT_TYPE = { INTERNAL | EXTERNAL }
//	ALTER ORGANIZATION SET DOMAIN_NAMES = ( '<domain>' [ , ... ] )
//	ALTER ORGANIZATION UNSET DOMAIN_NAMES
//	ALTER ORGANIZATION SET EVENT ROUTING TABLE <name> FOR ALL APPLICATION LISTINGS
//	ALTER ORGANIZATION UNSET EVENT ROUTING TABLE FOR ALL APPLICATION LISTINGS
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-organization
// https://docs.snowflake.com/en/sql-reference/sql/alter-organization-set-event-routing-table
// https://docs.snowflake.com/en/sql-reference/sql/alter-organization-unset-event-routing-table
func parseAlterOrganization(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	switch {
	case p.acceptWords("SET", "EVENT", "ROUTING", "TABLE"):
		name, err := p.parseObjectName()
		if err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "SET EVENT ROUTING TABLE", start).Names = []*ast.ObjectName{name}
		forStart := p.cur.Loc.Start
		if err := p.expectWords("FOR", "ALL", "APPLICATION", "LISTINGS"); err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "FOR ALL APPLICATION LISTINGS", forStart)
		return nil
	case p.acceptWords("UNSET", "EVENT", "ROUTING", "TABLE", "FOR", "ALL", "APPLICATION", "LISTINGS"):
		p.addObjectKeywordsClause(s, "UNSET EVENT ROUTING TABLE FOR ALL APPLICATION LISTINGS", start)
		return nil
	case p.cur.Type == kwSET:
		// One property a statement: each is a form of its own.
		p.advance() // consume SET
		if !p.startsObjectProperty() {
			return p.syntaxErrorAtCur()
		}
		opt, err := p.parseObjectProperty()
		if err != nil {
			return err
		}
		s.Clauses = append(s.Clauses, &ast.ObjectClause{
			Keywords: "SET", Options: []*ast.CopyOption{opt},
			Loc: ast.Loc{Start: start, End: p.prev.Loc.End},
		})
		return nil
	case p.acceptWords("UNSET", "DOMAIN_NAMES"):
		c := p.addObjectKeywordsClause(s, "UNSET", start)
		c.Keys = []string{"DOMAIN_NAMES"}
		return nil
	}
	return p.syntaxErrorAtCur()
}

// parseShowOrganizationScope reads the scopes organization objects add to
// SHOW, once parseShowScope has consumed IN or FOR, and reports whether one
// started at cur:
//
//	SHOW ORGANIZATION USERS IN ORGANIZATION USER GROUP <name>
//	SHOW VERSIONS IN ORGANIZATION PROFILE <name>
//	SHOW EVENT ROUTING TABLE ON ORGANIZATION FOR ALL APPLICATION LISTINGS
//
// https://docs.snowflake.com/en/sql-reference/sql/show-organization-users
// https://docs.snowflake.com/en/sql-reference/sql/show-versions-in-organization-profile
// https://docs.snowflake.com/en/sql-reference/sql/show-event-routing-table-on-organization
func (p *Parser) parseShowOrganizationScope(stmt *ast.ShowStmt) (bool, error) {
	in := p.prev.Type == kwIN
	switch {
	case in && stmt.ObjectClass == "ORGANIZATION USERS" && p.acceptWords("ORGANIZATION", "USER", "GROUP"):
		stmt.ScopeText = "ORGANIZATION USER GROUP"
		name, err := p.parseObjectSimpleName()
		if err != nil {
			return true, err
		}
		stmt.Scope = ast.ShowScopeOther
		stmt.ScopeName = name
		return true, nil
	case in && stmt.ObjectClass == "VERSIONS" && p.acceptWords("ORGANIZATION", "PROFILE"):
		stmt.ScopeText = "ORGANIZATION PROFILE"
		name, err := p.parseObjectName()
		if err != nil {
			return true, err
		}
		stmt.Scope = ast.ShowScopeOther
		stmt.ScopeName = name
		return true, nil
	case !in && stmt.ObjectClass == "EVENT ROUTING TABLE ON ORGANIZATION" && p.acceptWords("ALL", "APPLICATION", "LISTINGS"):
		stmt.Scope = ast.ShowScopeOther
		stmt.ScopeText = "ALL APPLICATION LISTINGS"
		return true, nil
	}
	return false, nil
}
