package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Objects around sharing and governance, read as ast.ObjectStmt: EXTERNAL
// CONSUMER (with its programmatic access tokens), CONTACT, MULTI PARTY
// APPROVAL POLICY and RESTRICTED SESSION SCOPE (RSS)
// ---------------------------------------------------------------------------

func init() {
	registerObject(&objectForm{
		words:  []string{"EXTERNAL", "CONSUMER"},
		create: parseCreateExternalConsumer,
		alter:  parseAlterExternalConsumer,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-external-consumer
		drop:      true,
		orReplace: true,
		// No props: ADD and ROTATE PAT carry their own runs, which
		// parseAlterExternalConsumer checks against the token's properties.
	})
	registerObject(&objectForm{
		words:  []string{"CONTACT"},
		create: parseCreateContact,
		alter:  parseAlterContact,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-contact
		drop:           true,
		dropNoIfExists: true,
		orReplace:      true,
		props:          contactProps,
		unset:          []string{},
	})
	registerObject(&objectForm{
		words:  []string{"MULTI", "PARTY", "APPROVAL", "POLICY"},
		create: parseObjectAsBody,
		alter:  parseAlterMultiPartyApprovalPolicy,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-multi-party-approval-policy
		drop:       true,
		orReplace:  true,
		noIfExists: true,
		props:      objectProps{},
		unset:      []string{},
	})
	for _, words := range [][]string{{"RESTRICTED", "SESSION", "SCOPE"}, {"RSS"}} {
		registerObject(&objectForm{
			words:  words,
			create: parseCreateRestrictedSessionScope,
			alter:  parseAlterRestrictedSessionScope,
			// https://docs.snowflake.com/en/sql-reference/sql/drop-restricted-session-scope
			drop:    true,
			orAlter: true,
			props:   objectProps{"COMMENT": propString},
		})
	}
}

// externalConsumerProps are CREATE EXTERNAL CONSUMER's properties.
var externalConsumerProps = objectProps{
	"REGION":  propString,
	"COMMENT": propString,
	"EMAIL":   propString,
}

// externalConsumerAlterProps are ALTER EXTERNAL CONSUMER ... SET's and
// UNSET's.
var externalConsumerAlterProps = objectProps{
	"COMMENT": propString,
	"EMAIL":   propString,
}

// patAddProps are ADD PAT's properties, patModifyProps MODIFY PAT ... SET's,
// and patRotateProps ROTATE PAT's.
var (
	patAddProps    = objectProps{"DAYS_TO_EXPIRY": propNumber, "COMMENT": propString}
	patModifyProps = patAddProps
	patRotateProps = objectProps{"EXPIRE_ROTATED_TOKEN_AFTER_HOURS": propNumber}
)

// checkObjectClauseProps holds clause c's properties (a run, a SET, or an
// UNSET) to props, and an UNSET's keys to unset.
func checkObjectClauseProps(kind string, c *ast.ObjectClause, props objectProps, unset []string) error {
	if c == nil {
		return nil
	}
	return (&objectForm{props: props, unset: unset}).checkClauses(&ast.ObjectStmt{Kind: kind, Clauses: []*ast.ObjectClause{c}})
}

// parseCreateExternalConsumer reads what follows the name of
//
//	CREATE [ OR REPLACE ] EXTERNAL CONSUMER [ IF NOT EXISTS ] <name>
//	  [ REGION = '<region>' ] [ COMMENT = '<string>' ] [ EMAIL = '<email>' ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-external-consumer
func parseCreateExternalConsumer(p *Parser, s *ast.ObjectStmt) error {
	c, err := p.parseObjectProperties(s, "", false)
	if err != nil {
		return err
	}
	if err := rejectObjectTags(s); err != nil {
		return err
	}
	return checkObjectClauseProps(s.Kind, c, externalConsumerProps, nil)
}

// parseAlterExternalConsumer reads what follows the name of
//
//	ALTER EXTERNAL CONSUMER [ IF EXISTS ] <name> SET { COMMENT = '<string>' | EMAIL = '<email>' } [ ... ]
//	ALTER EXTERNAL CONSUMER [ IF EXISTS ] <name> UNSET { COMMENT | EMAIL } [ , ... ]
//	ALTER EXTERNAL CONSUMER [ IF EXISTS ] <name> RENAME TO <new_name>
//	ALTER EXTERNAL CONSUMER <name> ADD <pat> [ IF NOT EXISTS ] <token_name>
//	  [ DAYS_TO_EXPIRY = <integer> ] [ COMMENT = '<string>' ]
//	ALTER EXTERNAL CONSUMER <name> MODIFY <pat> [ IF EXISTS ] <token_name>
//	  [ SET { DAYS_TO_EXPIRY = <integer> | COMMENT = '<string>' } [ ... ] ] [ UNSET COMMENT ]
//	ALTER EXTERNAL CONSUMER [ IF EXISTS ] <name> REMOVE <pat> [ IF EXISTS ] <token_name>
//	ALTER EXTERNAL CONSUMER <name> ROTATE <pat> [ IF EXISTS ] <token_name>
//	  [ EXPIRE_ROTATED_TOKEN_AFTER_HOURS = <integer> ]
//
// where <pat> is PROGRAMMATIC ACCESS TOKEN or PAT. A token action is a clause
// of its keywords as written (ADD PAT, MODIFY PROGRAMMATIC ACCESS TOKEN IF
// EXISTS, ...) holding the token's name; its properties follow as a run (ADD,
// ROTATE) or as SET and UNSET clauses (MODIFY).
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-external-consumer
// https://docs.snowflake.com/en/sql-reference/sql/alter-external-consumer-add-programmatic-access-token
// https://docs.snowflake.com/en/sql-reference/sql/alter-external-consumer-modify-programmatic-access-token
// https://docs.snowflake.com/en/sql-reference/sql/alter-external-consumer-remove-programmatic-access-token
// https://docs.snowflake.com/en/sql-reference/sql/alter-external-consumer-rotate-programmatic-access-token
func parseAlterExternalConsumer(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	for _, verb := range []string{"ADD", "MODIFY", "REMOVE", "ROTATE"} {
		if !isWordToken(p.cur, verb) {
			continue
		}
		next := p.peekNext()
		if !isWordToken(next, "PAT") && !isWordToken(next, "PROGRAMMATIC") {
			break
		}
		if s.IfExists && verb != "REMOVE" {
			// The pages write the consumer's IF EXISTS on REMOVE alone.
			return p.syntaxErrorAtCur()
		}
		return p.parseExternalConsumerToken(s, verb, start)
	}
	if p.cur.Type == kwRENAME {
		name, err := p.parseRenameTo()
		if err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "RENAME TO", start).Names = []*ast.ObjectName{name}
		return nil
	}
	ok, err := p.parseObjectAlterCommon(s, actSet|actUnset)
	if err != nil {
		return err
	}
	if !ok {
		return p.syntaxErrorAtCur()
	}
	if err := rejectObjectTags(s); err != nil {
		return err
	}
	return checkObjectClauseProps(s.Kind, s.Clauses[len(s.Clauses)-1], externalConsumerAlterProps, []string{"COMMENT", "EMAIL"})
}

// parseExternalConsumerToken reads a token action of ALTER EXTERNAL
// CONSUMER. cur is verb (ADD, MODIFY, REMOVE or ROTATE), which start begins.
func (p *Parser) parseExternalConsumerToken(s *ast.ObjectStmt, verb string, start int) error {
	p.advance() // consume verb
	keywords := verb + " PAT"
	if !p.acceptWords("PAT") {
		if err := p.expectWords("PROGRAMMATIC", "ACCESS", "TOKEN"); err != nil {
			return err
		}
		keywords = verb + " PROGRAMMATIC ACCESS TOKEN"
	}
	if verb == "ADD" {
		ifNotExists := false
		if err := p.parseIfNotExists(&ifNotExists); err != nil {
			return err
		}
		if ifNotExists {
			keywords += " IF NOT EXISTS"
		}
	} else {
		ifExists := false
		if err := p.parseIfExists(&ifExists); err != nil {
			return err
		}
		if ifExists {
			keywords += " IF EXISTS"
		}
	}
	name, err := p.parseObjectSimpleName()
	if err != nil {
		return err
	}
	p.addObjectKeywordsClause(s, keywords, start).Names = []*ast.ObjectName{name}
	switch verb {
	case "ADD", "ROTATE":
		props := patAddProps
		if verb == "ROTATE" {
			props = patRotateProps
		}
		c, err := p.parseObjectProperties(s, "", false)
		if err != nil {
			return err
		}
		if err := rejectObjectTags(s); err != nil {
			return err
		}
		return checkObjectClauseProps(s.Kind, c, props, nil)
	case "MODIFY":
		if p.cur.Type == kwSET {
			if _, err := p.parseObjectAlterCommon(s, actSet); err != nil {
				return err
			}
			if err := rejectObjectTags(s); err != nil {
				return err
			}
			if err := checkObjectClauseProps(s.Kind, s.Clauses[len(s.Clauses)-1], patModifyProps, nil); err != nil {
				return err
			}
		}
		if unsetStart := p.cur.Loc.Start; p.cur.Type == kwUNSET {
			p.advance() // consume UNSET
			if err := p.expectWords("COMMENT"); err != nil {
				return err
			}
			p.addObjectKeywordsClause(s, "UNSET", unsetStart).Keys = []string{"COMMENT"}
		}
	}
	return nil
}

// contactProps are CREATE CONTACT's and ALTER CONTACT ... SET's properties.
// USERS, EMAIL_DISTRIBUTION_LIST and URL exclude one another.
var contactProps = objectProps{
	"USERS":                   propStrings,
	"EMAIL_DISTRIBUTION_LIST": propString,
	"URL":                     propString,
	"COMMENT":                 propString,
}

// checkContactMethod fails when clause c sets more than one of USERS,
// EMAIL_DISTRIBUTION_LIST and URL.
func checkContactMethod(c *ast.ObjectClause) error {
	if c == nil {
		return nil
	}
	n := 0
	for _, o := range c.Options {
		switch o.Name {
		case "USERS", "EMAIL_DISTRIBUTION_LIST", "URL":
			n++
			if n > 1 {
				return &ParseError{
					Position: o.Loc.Start, End: o.Loc.End,
					Message: "syntax error: a contact takes one of USERS, EMAIL_DISTRIBUTION_LIST and URL",
				}
			}
		}
	}
	return nil
}

// parseCreateContact reads what follows the name of
//
//	CREATE [ OR REPLACE ] CONTACT [ IF NOT EXISTS ] <name>
//	  [ { USERS = ( '<user_name>' [ , ... ] ) | EMAIL_DISTRIBUTION_LIST = '<email>'
//	    | URL = '<url>' } ]
//	  [ COMMENT = '<string>' ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-contact
func parseCreateContact(p *Parser, s *ast.ObjectStmt) error {
	c, err := p.parseObjectProperties(s, "", false)
	if err != nil {
		return err
	}
	if err := rejectObjectTags(s); err != nil {
		return err
	}
	return checkContactMethod(c)
}

// parseAlterContact reads what follows the name of
//
//	ALTER CONTACT [ IF EXISTS ] <name> RENAME TO <new_name>
//	ALTER CONTACT [ IF EXISTS ] <name> SET
//	  { { USERS = ( '<user_name>' [ , ... ] ) | EMAIL_DISTRIBUTION_LIST = '<email>'
//	    | URL = '<url>' } [ COMMENT = '<string>' ] | COMMENT = '<string>' }
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-contact
func parseAlterContact(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	if p.cur.Type == kwRENAME {
		name, err := p.parseRenameTo()
		if err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "RENAME TO", start).Names = []*ast.ObjectName{name}
		return nil
	}
	ok, err := p.parseObjectAlterCommon(s, actSet)
	if err != nil {
		return err
	}
	if !ok {
		return p.syntaxErrorAtCur()
	}
	if err := rejectObjectTags(s); err != nil {
		return err
	}
	return checkContactMethod(s.Clauses[len(s.Clauses)-1])
}

// parseObjectAsBody reads the AS '<string>' a statement ends with, as
//
//	CREATE [ OR REPLACE ] MULTI PARTY APPROVAL POLICY [ IF NOT EXISTS ] <name> AS '<policy_body>'
//
// https://docs.snowflake.com/en/sql-reference/sql/create-multi-party-approval-policy
func parseObjectAsBody(p *Parser, s *ast.ObjectStmt) error {
	if p.cur.Type != kwAS {
		return p.syntaxErrorAtCur()
	}
	return p.parseObjectAsString(s)
}

// parseAlterMultiPartyApprovalPolicy reads what follows the name of
//
//	ALTER MULTI PARTY APPROVAL POLICY <name> RENAME TO <new_name>
//	ALTER MULTI PARTY APPROVAL POLICY <name> AS '<policy_body>'
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-multi-party-approval-policy
func parseAlterMultiPartyApprovalPolicy(p *Parser, s *ast.ObjectStmt) error {
	if p.cur.Type == kwRENAME {
		start := p.cur.Loc.Start
		name, err := p.parseRenameTo()
		if err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "RENAME TO", start).Names = []*ast.ObjectName{name}
		return nil
	}
	return parseObjectAsBody(p, s)
}

// parseCreateRestrictedSessionScope reads what follows the name of
//
//	CREATE { RESTRICTED SESSION SCOPE | RSS } [ IF NOT EXISTS ] <name>
//	  [ COMMENT = '<string>' ] AS '<definition>'
//	CREATE OR ALTER { RESTRICTED SESSION SCOPE | RSS } <name>
//	  [ COMMENT = '<string>' ] AS '<definition>'
//
// https://docs.snowflake.com/en/sql-reference/sql/create-restricted-session-scope
func parseCreateRestrictedSessionScope(p *Parser, s *ast.ObjectStmt) error {
	if s.OrAlter && s.IfNotExists {
		// CREATE OR ALTER writes no IF NOT EXISTS.
		return &ParseError{Position: s.Name.Loc.Start, End: s.Name.Loc.End, Message: "syntax error: CREATE OR ALTER takes no IF NOT EXISTS"}
	}
	if _, err := p.parseObjectProperties(s, "", false); err != nil {
		return err
	}
	if err := rejectObjectTags(s); err != nil {
		return err
	}
	return parseObjectAsBody(p, s)
}

// parseAlterRestrictedSessionScope reads what follows the name of
//
//	ALTER { RESTRICTED SESSION SCOPE | RSS } [ IF EXISTS ] <name> [ ADD | REMOVE ] AS '<definition>'
//	ALTER { RESTRICTED SESSION SCOPE | RSS } [ IF EXISTS ] <name> SET COMMENT = '<string>'
//	ALTER { RESTRICTED SESSION SCOPE | RSS } [ IF EXISTS ] <name> UNSET COMMENT
//
// ADD AS and REMOVE AS are clauses of their own keywords.
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-restricted-session-scope
func parseAlterRestrictedSessionScope(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	switch {
	case p.cur.Type == kwAS:
		return p.parseObjectAsString(s)
	case (isWordToken(p.cur, "ADD") || isWordToken(p.cur, "REMOVE")) && p.peekNext().Type == kwAS:
		verb := strings.ToUpper(p.advance().Str)
		if err := p.parseObjectAsString(s); err != nil {
			return err
		}
		c := s.Clauses[len(s.Clauses)-1]
		c.Keywords = verb + " AS"
		c.Loc.Start = start
		return nil
	}
	ok, err := p.parseObjectAlterCommon(s, actSet|actUnset)
	if err != nil {
		return err
	}
	if !ok {
		return p.syntaxErrorAtCur()
	}
	return rejectObjectTags(s)
}

// parseShowExternalConsumerScope reads the scope SHOW USER { PROGRAMMATIC
// ACCESS TOKENS | PATS } takes, once parseShowScope has consumed FOR and cur
// is EXTERNAL:
//
//	SHOW USER { PROGRAMMATIC ACCESS TOKENS | PATS } FOR EXTERNAL CONSUMER <name>
//
// https://docs.snowflake.com/en/sql-reference/sql/show-external-consumer-programmatic-access-tokens
func (p *Parser) parseShowExternalConsumerScope(stmt *ast.ShowStmt) error {
	if p.prev.Type != kwFOR || (stmt.ObjectClass != "USER PATS" && stmt.ObjectClass != "USER PROGRAMMATIC ACCESS TOKENS") {
		return p.syntaxErrorAtCur()
	}
	if err := p.expectWords("EXTERNAL", "CONSUMER"); err != nil {
		return err
	}
	name, err := p.parseObjectName()
	if err != nil {
		return err
	}
	stmt.Scope = ast.ShowScopeOther
	stmt.ScopeText = "EXTERNAL CONSUMER"
	stmt.ScopeName = name
	return nil
}
