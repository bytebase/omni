package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// APPLICATION PACKAGE (Native Apps), read as an ast.ObjectStmt
// ---------------------------------------------------------------------------

func init() {
	registerObject(&objectForm{
		words:  []string{"APPLICATION", "PACKAGE"},
		create: parseCreateApplicationPackage,
		alter:  applicationTagless(applicationIfExistsOn(parseAlterApplicationPackage, "SET", "UNSET", "RENAME TO"), "SET"),
		// https://docs.snowflake.com/en/sql-reference/sql/drop-application-package
		drop:      true,
		nameParts: 1, // an account-level object
		props:     applicationPackageProps,
		// SET's properties, and a release channel's SET ACCOUNTS
		// (parseAlterApplicationPackage keeps ACCOUNTS to the channel).
		alterProps: applicationPackageAlterProps,
		unset: []string{"DATA_RETENTION_TIME_IN_DAYS", "MAX_DATA_EXTENSION_TIME_IN_DAYS",
			"DEFAULT_DDL_COLLATION", "COMMENT", "DISTRIBUTION"},
	})
}

// applicationPackageProps are the properties CREATE APPLICATION PACKAGE and
// ALTER APPLICATION PACKAGE ... SET document.
var applicationPackageProps = objectProps{
	"DATA_RETENTION_TIME_IN_DAYS":       propNumber,
	"MAX_DATA_EXTENSION_TIME_IN_DAYS":   propNumber,
	"DEFAULT_DDL_COLLATION":             propString,
	"COMMENT":                           propString,
	"DISTRIBUTION":                      propOneOf("INTERNAL", "EXTERNAL"),
	"LISTING_AUTO_REFRESH":              propBool,
	"MULTIPLE_INSTANCES":                propOneOf("TRUE"),
	"ENABLE_RELEASE_CHANNELS":           propBool,
	"AUTOMATIC_APPLICATION_MAINTENANCE": propBool,
}

// applicationPackageAlterProps are ALTER APPLICATION PACKAGE ... SET's:
// CREATE's, with ENABLE_RELEASE_CHANNELS only ever set TRUE, and the ACCOUNTS
// a release channel's SET ACCOUNTS = ( ... ) takes.
var applicationPackageAlterProps = func() objectProps {
	m := objectProps{"ACCOUNTS": propNames}
	for k, v := range applicationPackageProps {
		m[k] = v
	}
	m["ENABLE_RELEASE_CHANNELS"] = propOneOf("TRUE")
	return m
}()

// applicationReleaseDirectiveProps are the properties a release directive
// takes after its accounts: the version and patch it serves and when
// consumers upgrade to them.
var applicationReleaseDirectiveProps = objectProps{
	"VERSION":                       propName,
	"PATCH":                         propNumber,
	"UPGRADE_AFTER":                 propString,
	"UPGRADE_IN_MAINTENANCE_WINDOW": propBool,
	"UPGRADE_DEADLINE":              propString,
}

// parseCreateApplicationPackage reads what follows the name of
//
//	CREATE APPLICATION PACKAGE [ IF NOT EXISTS ] <name>
//	  [ DATA_RETENTION_TIME_IN_DAYS = <integer> ]
//	  [ MAX_DATA_EXTENSION_TIME_IN_DAYS = <integer> ]
//	  [ DEFAULT_DDL_COLLATION = '<collation_specification>' ]
//	  [ COMMENT = '<string_literal>' ] [ DISTRIBUTION = { INTERNAL | EXTERNAL } ]
//	  [ LISTING_AUTO_REFRESH = { TRUE | FALSE } ] [ MULTIPLE_INSTANCES = TRUE ]
//	  [ ENABLE_RELEASE_CHANNELS = { TRUE | FALSE } ]
//	  [ AUTOMATIC_APPLICATION_MAINTENANCE = { TRUE | FALSE } ]
//	  [ [ WITH ] TAG ( <tag_name> = '<tag_value>' [ , ... ] ) ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-application-package
func parseCreateApplicationPackage(p *Parser, s *ast.ObjectStmt) error {
	_, err := p.parseObjectProperties(s, "", false)
	return err
}

// parseAlterApplicationPackage reads what follows the name of
//
//	ALTER APPLICATION PACKAGE [ IF EXISTS ] <name> SET <property> [ ... ]
//	ALTER APPLICATION PACKAGE [ IF EXISTS ] <name> UNSET <property> [ , ... ]
//	ALTER APPLICATION PACKAGE [ IF EXISTS ] <name> RENAME TO <new_name>
//	ALTER APPLICATION PACKAGE <name> SET TAG <tag_name> = '<tag_value>' [ , ... ]
//	ALTER APPLICATION PACKAGE <name> UNSET TAG <tag_name> [ , ... ]
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-application-package
//
// and of its version, release channel and release directive pages:
//
//	ALTER APPLICATION PACKAGE <name> ADD VERSION [ <version> ] USING '<path>' [ LABEL = '<label>' ]
//	ALTER APPLICATION PACKAGE <name> DROP VERSION <version> [ ASYNC ]
//	ALTER APPLICATION PACKAGE <name> ADD PATCH [ <patch_num> ] FOR VERSION [ <version> ]
//	  USING '<path>' [ LABEL = '<label>' ]
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-application-package-version
//
//	ALTER APPLICATION PACKAGE <name> REGISTER VERSION [ <version> ] USING '<path>' [ LABEL = '<label>' ]
//	ALTER APPLICATION PACKAGE <name> DEREGISTER VERSION <version>
//	ALTER APPLICATION PACKAGE <name> MODIFY RELEASE CHANNEL <channel> <channel_action>
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-application-package-release-channel
//
//	ALTER APPLICATION PACKAGE <name> MODIFY RELEASE DIRECTIVE <directive> <directive_action>
//	ALTER APPLICATION PACKAGE <name> SET DEFAULT RELEASE DIRECTIVE
//	  VERSION = <version> PATCH = <patch_num> [ <upgrade_property> ... ]
//	ALTER APPLICATION PACKAGE <name> SET RELEASE DIRECTIVE <directive>
//	  ACCOUNTS = ( <org>.<account> [ , ... ] ) VERSION = <version> PATCH = <patch_num>
//	  [ <upgrade_property> ... ] [ FORCE ]
//	ALTER APPLICATION PACKAGE <name> UNSET RELEASE DIRECTIVE <directive>
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-application-package-release-directive
func parseAlterApplicationPackage(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	switch {
	case p.acceptWords("ADD", "VERSION"):
		return p.parseApplicationPackageVersionUsing(s, "ADD VERSION", start)
	case p.acceptWords("REGISTER", "VERSION"):
		return p.parseApplicationPackageVersionUsing(s, "REGISTER VERSION", start)
	case p.acceptWords("ADD", "PATCH"):
		if p.cur.Type == tokInt {
			if _, err := p.parseApplicationLiteralClause(s, "ADD PATCH", start, tokInt); err != nil {
				return err
			}
		} else {
			p.addObjectKeywordsClause(s, "ADD PATCH", start)
		}
		start := p.cur.Loc.Start
		if err := p.expectWords("FOR", "VERSION"); err != nil {
			return err
		}
		return p.parseApplicationPackageVersionUsing(s, "FOR VERSION", start)
	case p.acceptWords("DROP", "VERSION"):
		if _, err := p.parseApplicationNamedClause(s, "DROP VERSION", start); err != nil {
			return err
		}
		if start := p.cur.Loc.Start; p.acceptWords("ASYNC") {
			p.addObjectKeywordsClause(s, "ASYNC", start)
		}
		return nil
	case p.acceptWords("DEREGISTER", "VERSION"):
		_, err := p.parseApplicationNamedClause(s, "DEREGISTER VERSION", start)
		return err
	case p.acceptWords("MODIFY", "RELEASE", "CHANNEL"):
		if _, err := p.parseApplicationNamedClause(s, "MODIFY RELEASE CHANNEL", start); err != nil {
			return err
		}
		return p.parseApplicationReleaseChannelAction(s)
	case p.atWords("MODIFY", "RELEASE", "DIRECTIVE"):
		return p.parseApplicationModifyReleaseDirective(s)
	}
	if ok, err := p.parseApplicationSetReleaseDirective(s, true); err != nil || ok {
		return err
	}
	ok, err := p.parseObjectAlterCommon(s, actSet|actUnset|actRename|actTags)
	if err != nil || ok {
		if err == nil {
			err = rejectApplicationPackageAccounts(s)
		}
		return err
	}
	return p.syntaxErrorAtCur()
}

// rejectApplicationPackageAccounts fails a package's own SET of ACCOUNTS,
// which alterProps admits for a release channel's SET ACCOUNTS alone.
func rejectApplicationPackageAccounts(s *ast.ObjectStmt) error {
	c := s.Clauses[len(s.Clauses)-1]
	if c.Keywords != "SET" {
		return nil
	}
	for _, o := range c.Options {
		if o.Name == "ACCOUNTS" {
			return applicationOptionError(o)
		}
	}
	return nil
}

// parseApplicationPackageVersionUsing reads
// [ <version> ] USING '<path>' [ LABEL = '<label>' ] after ADD VERSION,
// REGISTER VERSION or ADD PATCH [ <num> ] FOR VERSION, whose keywords start
// began: a clause of those keywords naming the version, and a USING clause
// holding the path and label.
func (p *Parser) parseApplicationPackageVersionUsing(s *ast.ObjectStmt, keywords string, start int) error {
	c := p.addObjectKeywordsClause(s, keywords, start)
	if !p.curIsWord("USING") {
		version, err := p.parseObjectSimpleName()
		if err != nil {
			return err
		}
		c.Names = []*ast.ObjectName{version}
		c.Loc.End = p.prev.Loc.End
	}
	start = p.cur.Loc.Start
	if err := p.expectWords("USING"); err != nil {
		return err
	}
	using, err := p.parseApplicationLiteralClause(s, "USING", start, tokString)
	if err != nil {
		return err
	}
	if p.curIsWord("LABEL") {
		return p.parseApplicationRequiredOption(using, "LABEL", propString)
	}
	return nil
}

// parseApplicationReleaseChannelAction reads what follows
// MODIFY RELEASE CHANNEL <channel>:
//
//	{ ADD | DROP } VERSION <version>
//	{ ADD | REMOVE | SET } ACCOUNTS = ( <org>.<account> [ , ... ] )
//	SET DEFAULT RELEASE DIRECTIVE VERSION = <version> PATCH = <patch_num> [ ... ]
//	SET RELEASE DIRECTIVE <directive> ACCOUNTS = ( ... ) VERSION = <version> PATCH = <patch_num> [ ... ]
//	MODIFY RELEASE DIRECTIVE <directive> <directive_action>
//	UNSET RELEASE DIRECTIVE <directive>
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-application-package-release-channel
func (p *Parser) parseApplicationReleaseChannelAction(s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	switch {
	case p.acceptWords("ADD", "VERSION"):
		_, err := p.parseApplicationNamedClause(s, "ADD VERSION", start)
		return err
	case p.acceptWords("DROP", "VERSION"):
		_, err := p.parseApplicationNamedClause(s, "DROP VERSION", start)
		return err
	case p.atWords("MODIFY", "RELEASE", "DIRECTIVE"):
		return p.parseApplicationModifyReleaseDirective(s)
	}
	if ok, err := p.parseApplicationSetReleaseDirective(s, false); err != nil || ok {
		return err
	}
	for _, verb := range []string{"ADD", "REMOVE", "SET"} {
		if p.curIsWord(verb) && p.peekIsWord("ACCOUNTS") {
			p.advance() // consume the verb
			return p.parseApplicationAccounts(p.addObjectKeywordsClause(s, verb, start))
		}
	}
	return p.syntaxErrorAtCur()
}

// parseApplicationSetReleaseDirective reads one of
//
//	SET DEFAULT RELEASE DIRECTIVE VERSION = <version> PATCH = <patch_num>
//	  [ UPGRADE_AFTER = '<timestamp>' ] [ UPGRADE_IN_MAINTENANCE_WINDOW = { TRUE | FALSE } ]
//	  [ UPGRADE_DEADLINE = '<timestamp>' ]
//	SET RELEASE DIRECTIVE <directive> ACCOUNTS = ( ... ) VERSION = <version>
//	  PATCH = <patch_num> [ ... ] [ FORCE ]
//	UNSET RELEASE DIRECTIVE <directive>
//
// when one starts at cur, and reports whether it did. force says whether
// SET RELEASE DIRECTIVE may end in FORCE, as the release directive page
// writes it and the release channel page does not.
func (p *Parser) parseApplicationSetReleaseDirective(s *ast.ObjectStmt, force bool) (bool, error) {
	start := p.cur.Loc.Start
	switch {
	case p.acceptWords("SET", "DEFAULT", "RELEASE", "DIRECTIVE"):
		c := p.addObjectKeywordsClause(s, "SET DEFAULT RELEASE DIRECTIVE", start)
		if err := p.parseApplicationOptions(c, applicationReleaseDirectiveProps); err != nil {
			return true, err
		}
		return true, p.requireObjectProperties(c, "VERSION", "PATCH")
	case p.acceptWords("SET", "RELEASE", "DIRECTIVE"):
		c, err := p.parseApplicationNamedClause(s, "SET RELEASE DIRECTIVE", start)
		if err != nil {
			return true, err
		}
		if err := p.parseApplicationAccounts(c); err != nil {
			return true, err
		}
		if err := p.parseApplicationOptions(c, applicationReleaseDirectiveProps); err != nil {
			return true, err
		}
		if err := p.requireObjectProperties(c, "VERSION", "PATCH"); err != nil {
			return true, err
		}
		if start := p.cur.Loc.Start; force && p.acceptWords("FORCE") {
			p.addObjectKeywordsClause(s, "FORCE", start)
		}
		return true, nil
	case p.acceptWords("UNSET", "RELEASE", "DIRECTIVE"):
		_, err := p.parseApplicationNamedClause(s, "UNSET RELEASE DIRECTIVE", start)
		return true, err
	}
	return false, nil
}

// parseApplicationModifyReleaseDirective reads
//
//	MODIFY RELEASE DIRECTIVE <directive>
//	  { VERSION = <version> PATCH = <patch_num> [ <upgrade_property> ... ]
//	  | ADD ACCOUNTS = ( <org>.<account> [ , ... ] ) [ VERSION = <version> PATCH = <patch_num> ] [ FORCE ]
//	  | REMOVE ACCOUNTS = ( <org>.<account> [ , ... ] ) [ VERSION = <version> PATCH = <patch_num> ] }
//
// cur is MODIFY.
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-application-package-release-directive
func (p *Parser) parseApplicationModifyReleaseDirective(s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	if err := p.expectWords("MODIFY", "RELEASE", "DIRECTIVE"); err != nil {
		return err
	}
	c, err := p.parseApplicationNamedClause(s, "MODIFY RELEASE DIRECTIVE", start)
	if err != nil {
		return err
	}
	start = p.cur.Loc.Start
	for _, verb := range []string{"ADD", "REMOVE"} {
		if p.curIsWord(verb) && p.peekIsWord("ACCOUNTS") {
			p.advance() // consume the verb
			c := p.addObjectKeywordsClause(s, verb, start)
			if err := p.parseApplicationAccounts(c); err != nil {
				return err
			}
			if p.curIsWord("VERSION") {
				if err := p.parseApplicationOptions(c, applicationAccountsVersionProps); err != nil {
					return err
				}
				if err := p.requireObjectProperties(c, "PATCH"); err != nil {
					return err
				}
			}
			if start := p.cur.Loc.Start; verb == "ADD" && p.acceptWords("FORCE") {
				p.addObjectKeywordsClause(s, "FORCE", start)
			}
			return nil
		}
	}
	if !p.curIsWord("VERSION") {
		return p.syntaxErrorAtCur()
	}
	if err := p.parseApplicationOptions(c, applicationReleaseDirectiveProps); err != nil {
		return err
	}
	return p.requireObjectProperties(c, "VERSION", "PATCH")
}

// applicationAccountsVersionProps are what ADD or REMOVE ACCOUNTS on a
// release directive may take after the accounts: the version and patch alone,
// without the upgrade schedule a version change takes.
var applicationAccountsVersionProps = objectProps{
	"VERSION": propName,
	"PATCH":   propNumber,
}

// parseApplicationAccounts reads ACCOUNTS = ( <org_name>.<account_name>
// [ , ... ] ) into an ACCOUNTS option of c.
func (p *Parser) parseApplicationAccounts(c *ast.ObjectClause) error {
	opt := &ast.CopyOption{Name: "ACCOUNTS", Loc: ast.Loc{Start: p.cur.Loc.Start}}
	if err := p.expectWords("ACCOUNTS"); err != nil {
		return err
	}
	if _, err := p.expect('='); err != nil {
		return err
	}
	if p.cur.Type != '(' {
		return p.syntaxErrorAtCur()
	}
	if err := p.parseObjectValue(opt); err != nil {
		return err
	}
	opt.Loc.End = p.prev.Loc.End
	if !propNames.matches(opt) {
		return applicationOptionError(opt)
	}
	for _, e := range opt.Group {
		org, account, ok := strings.Cut(e.Name, ".")
		if !ok || org == "" || account == "" || strings.Contains(account, ".") {
			return &ParseError{Position: e.Loc.Start, End: e.Loc.End, Message: "syntax error at or near " + e.Name}
		}
	}
	c.Options = append(c.Options, opt)
	c.Loc.End = opt.Loc.End
	return nil
}
