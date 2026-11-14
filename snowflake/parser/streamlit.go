package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// STREAMLIT, read as an ast.ObjectStmt
// ---------------------------------------------------------------------------

func init() {
	registerObject(&objectForm{
		words:     []string{"STREAMLIT"},
		create:    parseCreateStreamlit,
		alter:     parseAlterStreamlit,
		orReplace: true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-streamlit
		drop:  true,
		props: streamlitProps,
		unset: []string{"SECRETS", "EXTERNAL_ACCESS_INTEGRATIONS", "QUERY_WAREHOUSE", "TITLE", "COMMENT"},
	})
}

// streamlitProps are the properties CREATE STREAMLIT and ALTER STREAMLIT ...
// SET document, in both their current and their ROOT_LOCATION forms.
var streamlitProps = objectProps{
	"ROOT_LOCATION":                propString,
	"MAIN_FILE":                    propString,
	"QUERY_WAREHOUSE":              propName,
	"RUNTIME_NAME":                 propString,
	"COMPUTE_POOL":                 propName,
	"COMMENT":                      propString,
	"TITLE":                        propString,
	"IMPORTS":                      propStrings,
	"EXTERNAL_ACCESS_INTEGRATIONS": propNames,
	"SECRETS":                      propGroup,
}

// streamlitCheckPush holds PUSH's properties to the ones the page documents,
// each set once. The page writes their values as bare placeholders
// (<git_author_email>, ...), which an email address cannot be, so the ones a
// person types (USERNAME, NAME, EMAIL, COMMENT) take a string or a name;
// GIT_CREDENTIALS names a secret, and PASSWORD is a string.
func streamlitCheckPush(opts []*ast.CopyOption) error {
	for i, o := range opts {
		ok := true
		for _, prev := range opts[:i] {
			ok = ok && prev.Name != o.Name
		}
		switch o.Name {
		case "USERNAME", "NAME", "EMAIL", "COMMENT":
			ok = ok && (propName.matches(o) || propString.matches(o))
		case "GIT_CREDENTIALS":
			ok = ok && propName.matches(o)
		case "PASSWORD":
			ok = ok && propString.matches(o)
		default:
			ok = false
		}
		if !ok {
			return &ParseError{
				Position: o.Loc.Start, End: o.Loc.End,
				Message: "syntax error at or near " + o.Name,
			}
		}
	}
	return nil
}

// parseCreateStreamlit reads what follows the name of
//
//	CREATE [ OR REPLACE ] STREAMLIT [ IF NOT EXISTS ] <name>
//	  [ FROM <source_location> ]
//	  [ MAIN_FILE = '<path_to_main_file_in_root_directory>' ]
//	  [ QUERY_WAREHOUSE = <warehouse_name> ] [ RUNTIME_NAME = '<runtime_name>' ]
//	  [ COMPUTE_POOL = <compute_pool_name> ] [ COMMENT = '<string_literal>' ]
//	  [ TITLE = '<app_title>' ]
//	  [ IMPORTS = ( '<stage_path_and_file_name_to_read>' [ , ... ] ) ]
//	  [ EXTERNAL_ACCESS_INTEGRATIONS = ( <integration_name> [ , ... ] ) ]
//	  [ SECRETS = ( '<snowflake_secret_name>' = <snowflake_secret> [ , ... ] ) ]
//
//	CREATE [ OR REPLACE ] STREAMLIT [ IF NOT EXISTS ] <name>
//	  ROOT_LOCATION = '<stage_path_and_root_directory>'
//	  MAIN_FILE = '<path_to_main_file_in_root_directory>'
//	  [ QUERY_WAREHOUSE = <warehouse_name> ] [ COMMENT = '<string_literal>' ]
//	  [ TITLE = '<app_title>' ] [ IMPORTS = ( ... ) ]
//	  [ EXTERNAL_ACCESS_INTEGRATIONS = ( <integration_name> [ , ... ] ) ]
//
// The source location is written as a string ('@stage/path'), as the page's
// parameters do. The legacy ROOT_LOCATION form takes no FROM and needs
// MAIN_FILE.
//
// https://docs.snowflake.com/en/sql-reference/sql/create-streamlit
func parseCreateStreamlit(p *Parser, s *ast.ObjectStmt) error {
	from := p.cur.Type == kwFROM
	if from {
		if err := p.parseStreamlitFrom(s); err != nil {
			return err
		}
	}
	props, err := p.parseObjectProperties(s, "", false)
	if err != nil {
		return err
	}
	if props != nil && streamlitHasProperty(props, "ROOT_LOCATION") {
		if from {
			return &ParseError{
				Position: props.Loc.Start, End: props.Loc.End,
				Message: "syntax error: ROOT_LOCATION cannot follow FROM",
			}
		}
		return p.requireObjectProperties(props, "MAIN_FILE")
	}
	return nil
}

// parseAlterStreamlit reads what follows the name of
//
//	ALTER STREAMLIT [ IF EXISTS ] <name> SET <property> [ ... ]
//	ALTER STREAMLIT [ IF EXISTS ] <name> UNSET <property> [ , ... ]
//	ALTER STREAMLIT [ IF EXISTS ] <name> RENAME TO <new_name>
//	ALTER STREAMLIT <name> COMMIT
//	ALTER STREAMLIT <name> PUSH [ TO <git_branch_uri> ]
//	  [ { GIT_CREDENTIALS = <snowflake_secret>
//	    | USERNAME = <git_username> PASSWORD = <git_password> }
//	    NAME = <git_author_name> EMAIL = <git_author_email> ]
//	  [ COMMENT = <git_push_comment> ]
//	ALTER STREAMLIT <name> ABORT
//	ALTER STREAMLIT <name> PULL
//	ALTER STREAMLIT <name> ADD LIVE VERSION FROM LAST
//
// The page writes IF EXISTS only before SET, UNSET and RENAME TO.
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-streamlit
func parseAlterStreamlit(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	if ok, err := p.parseObjectAlterCommon(s, actSet|actUnset|actRename); ok || err != nil {
		return err
	}
	if s.IfExists {
		return p.syntaxErrorAtCur()
	}
	switch {
	case p.acceptWords("COMMIT"):
		p.addObjectKeywordsClause(s, "COMMIT", start)
	case p.acceptWords("ABORT"):
		p.addObjectKeywordsClause(s, "ABORT", start)
	case p.acceptWords("PULL"):
		p.addObjectKeywordsClause(s, "PULL", start)
	case p.acceptWords("ADD", "LIVE", "VERSION", "FROM", "LAST"):
		p.addObjectKeywordsClause(s, "ADD LIVE VERSION FROM LAST", start)
	case p.acceptWords("PUSH"):
		return p.parseStreamlitPush(s, start)
	default:
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseStreamlitPush reads PUSH [ TO <git_branch_uri> ] and its
// `KEY = value` properties into one clause: PUSH, or PUSH TO whose first
// option, nameless, holds the branch. PUSH has been consumed; start is where
// it began.
func (p *Parser) parseStreamlitPush(s *ast.ObjectStmt, start int) error {
	c := p.addObjectKeywordsClause(s, "PUSH", start)
	props := c.Options
	if p.cur.Type == kwTO {
		p.advance() // consume TO
		c.Keywords = "PUSH TO"
		branch := &ast.CopyOption{Loc: ast.Loc{Start: p.cur.Loc.Start}}
		if p.cur.Type == '(' {
			return p.syntaxErrorAtCur()
		}
		if err := p.parseObjectValue(branch); err != nil {
			return err
		}
		branch.Loc.End = p.prev.Loc.End
		if !propName.matches(branch) && !propString.matches(branch) {
			return &ParseError{Position: branch.Loc.Start, End: branch.Loc.End, Message: "syntax error: a branch is a string or a name"}
		}
		c.Options = append(c.Options, branch)
	}
	for p.startsObjectProperty() {
		opt, err := p.parseObjectProperty()
		if err != nil {
			return err
		}
		props = append(props, opt)
		c.Options = append(c.Options, opt)
	}
	c.Loc.End = p.prev.Loc.End
	return streamlitCheckPush(props)
}

// parseStreamlitFrom reads FROM '<source_location>' into a FROM clause whose
// one option, unnamed, holds the string, as CREATE STREAMLIT and CREATE
// NOTEBOOK write their source. cur is FROM.
func (p *Parser) parseStreamlitFrom(s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	p.advance() // consume FROM
	tok, err := p.expect(tokString)
	if err != nil {
		return err
	}
	src := &ast.CopyOption{
		Lit: &ast.Literal{Kind: ast.LitString, Value: tok.Str, Loc: tok.Loc},
		Loc: tok.Loc,
	}
	p.addObjectKeywordsClause(s, "FROM", start).Options = []*ast.CopyOption{src}
	return nil
}

// streamlitHasProperty reports whether clause c sets the property name.
func streamlitHasProperty(c *ast.ObjectClause, name string) bool {
	for _, o := range c.Options {
		if o.Name == name {
			return true
		}
	}
	return false
}
