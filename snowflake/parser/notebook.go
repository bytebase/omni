package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// NOTEBOOK, NOTEBOOK PROJECT and EXECUTE NOTEBOOK [ PROJECT ], read as an
// ast.ObjectStmt
// ---------------------------------------------------------------------------
//
// NOTEBOOK PROJECT is not registered as an object of its own: PROJECT is not
// reserved, so `CREATE NOTEBOOK project ...` and `CREATE NOTEBOOK PROJECT.s.n`
// create a notebook named by it. parseCreateNotebook reads the name first and
// turns to NOTEBOOK PROJECT when the name is the bare word PROJECT and another
// name follows it.
//
// EXECUTE NOTEBOOK runs a notebook, so it is read as an ObjectStmt too, with
// Verb "EXECUTE": Kind "NOTEBOOK" with its arguments as a clause, or Kind
// "NOTEBOOK PROJECT" with its properties.

func init() {
	registerObject(&objectForm{
		words:     []string{"NOTEBOOK"},
		create:    parseCreateNotebook,
		alter:     parseAlterNotebook,
		orReplace: true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-notebook writes
		// no IF EXISTS.
		drop:           true,
		dropNoIfExists: true,
		props: objectProps{
			"MAIN_FILE":                       propString,
			"COMMENT":                         propString,
			"QUERY_WAREHOUSE":                 propName,
			"IDLE_AUTO_SHUTDOWN_TIME_SECONDS": propNumber,
			"RUNTIME_NAME":                    propString,
			"COMPUTE_POOL":                    propString,
			"WAREHOUSE":                       propName,
		},
		alterProps: objectProps{
			"COMMENT":                         propString,
			"QUERY_WAREHOUSE":                 propName,
			"IDLE_AUTO_SHUTDOWN_TIME_SECONDS": propNumber,
			"SECRETS":                         propGroup,
		},
		unset: []string{},
	})
}

// notebookProjectProps are CREATE NOTEBOOK PROJECT's properties, and
// notebookExecuteProjectProps EXECUTE NOTEBOOK PROJECT's.
var (
	notebookProjectProps        = objectProps{"COMMENT": propString}
	notebookExecuteProjectProps = objectProps{
		"MAIN_FILE":                    propString,
		"COMPUTE_POOL":                 propString,
		"QUERY_WAREHOUSE":              propString,
		"RUNTIME":                      propString,
		"ARGUMENTS":                    propString,
		"REQUIREMENTS_FILE":            propString,
		"ARTIFACT_REPOSITORIES":        propNames,
		"EXTERNAL_ACCESS_INTEGRATIONS": propNames,
		"SECRETS":                      propNames,
	}
)

// notebookCheckOptions holds the options of a clause that the registry does
// not check (NOTEBOOK PROJECT's, EXECUTE's) to props.
func notebookCheckOptions(c *ast.ObjectClause, props objectProps) error {
	if c == nil {
		return nil
	}
	if len(c.Tags) > 0 {
		return &ParseError{Position: c.Loc.Start, End: c.Loc.End, Message: "syntax error: TAG is not documented here"}
	}
	for i, o := range c.Options {
		for _, prev := range c.Options[:i] {
			if prev.Name == o.Name {
				return &ParseError{
					Position: o.Loc.Start, End: o.Loc.End,
					Message: "syntax error: " + o.Name + " is set twice",
				}
			}
		}
		if v, ok := props[o.Name]; !ok || !v.matches(o) {
			return &ParseError{
				Position: o.Loc.Start, End: o.Loc.End,
				Message: "syntax error at or near " + o.Name,
			}
		}
	}
	return nil
}

// parseCreateNotebook reads what follows the name of
//
//	CREATE [ OR REPLACE ] NOTEBOOK [ IF NOT EXISTS ] <name>
//	  [ FROM '<source_location>' ]
//	  [ MAIN_FILE = '<main_file_name>' ]
//	  [ COMMENT = '<string_literal>' ]
//	  [ QUERY_WAREHOUSE = <warehouse_to_run_nb_and_sql_queries_in> ]
//	  [ IDLE_AUTO_SHUTDOWN_TIME_SECONDS = <number_of_seconds> ]
//	  [ RUNTIME_NAME = '<runtime_name>' ]
//	  [ COMPUTE_POOL = '<compute_pool_name>' ]
//	  [ WAREHOUSE = <warehouse_to_run_notebook_python_runtime> ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-notebook
//
// or, when the name is the word PROJECT and a name follows, of
// CREATE NOTEBOOK PROJECT (parseCreateNotebookProject).
func parseCreateNotebook(p *Parser, s *ast.ObjectStmt) error {
	if p.startsNotebookProjectName(s) {
		return p.parseCreateNotebookProject(s)
	}
	if p.cur.Type == kwFROM {
		if err := p.parseStreamlitFrom(s); err != nil {
			return err
		}
	}
	_, err := p.parseObjectProperties(s, "", false)
	return err
}

// startsNotebookProjectName reports whether the CREATE NOTEBOOK statement s,
// whose name has been read, is CREATE NOTEBOOK PROJECT: its name is the bare
// word PROJECT and cur begins IF NOT EXISTS or another name, not FROM or a
// property. It consumes nothing.
func (p *Parser) startsNotebookProjectName(s *ast.ObjectStmt) bool {
	n := s.Name
	if s.OrReplace || s.IfNotExists || !n.Database.IsEmpty() || !n.Schema.IsEmpty() ||
		n.Name.Quoted || strings.ToUpper(n.Name.Name) != "PROJECT" {
		return false
	}
	switch {
	case p.cur.Type == kwIF:
		return p.peekNext().Type == kwNOT
	case p.cur.Type == kwFROM || p.startsObjectProperty():
		return false
	}
	return p.cur.Type == tokIdent || p.cur.Type == tokQuotedIdent || p.cur.Type >= 700
}

// parseCreateNotebookProject reads, from the project's name on,
//
//	CREATE NOTEBOOK PROJECT <database_name>.<schema_name>.<project_name>
//	  FROM 'snow://workspace/<workspace_path>'
//	  [ COMMENT = '<string_literal>' ]
//
//	CREATE NOTEBOOK PROJECT [ IF NOT EXISTS ] <database_name>.<schema_name>.<project_name>
//	  FROM '@<database_name>.<schema_name>.<stage_name>'
//	  [ COMMENT = '<string_literal>' ]
//
// s's name, the word PROJECT, becomes part of its Kind.
//
// https://docs.snowflake.com/en/sql-reference/sql/create-notebook-project
func (p *Parser) parseCreateNotebookProject(s *ast.ObjectStmt) error {
	s.Kind = "NOTEBOOK PROJECT"
	if err := p.parseIfNotExists(&s.IfNotExists); err != nil {
		return err
	}
	name, err := p.parseObjectName()
	if err != nil {
		return err
	}
	s.Name = name
	if p.cur.Type != kwFROM {
		return p.syntaxErrorAtCur()
	}
	if err := p.parseStreamlitFrom(s); err != nil {
		return err
	}
	props, err := p.parseObjectProperties(s, "", false)
	if err != nil {
		return err
	}
	return notebookCheckOptions(props, notebookProjectProps)
}

// parseAlterNotebook reads what follows the name of
//
//	ALTER NOTEBOOK [ IF EXISTS ] <name> RENAME TO <new_name>
//	ALTER NOTEBOOK [ IF EXISTS ] <name> SET
//	  [ COMMENT = '<string_literal>' ]
//	  [ QUERY_WAREHOUSE = <warehouse_to_run_nb_and_sql_queries_in> ]
//	  [ IDLE_AUTO_SHUTDOWN_TIME_SECONDS = <number_of_seconds> ]
//	  [ SECRETS = ( '<snowflake_secret_variable_name>' = <snowflake_secret_name> [ , ... ] ) ]
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-notebook
func parseAlterNotebook(p *Parser, s *ast.ObjectStmt) error {
	if ok, err := p.parseObjectAlterCommon(s, actSet|actRename); ok || err != nil {
		return err
	}
	return p.syntaxErrorAtCur()
}

// parseExecuteNotebook reads
//
//	EXECUTE NOTEBOOK <name>( [ '<parameter_string>' [ , ... ] ] )
//
// https://docs.snowflake.com/en/sql-reference/sql/execute-notebook
//
// with its arguments as a clause of no keywords whose one option, unnamed,
// lists them, and
//
//	EXECUTE NOTEBOOK PROJECT <database_name>.<schema_name>.<project_name>
//	  MAIN_FILE = '<notebook_file_name>.ipynb'
//	  COMPUTE_POOL = '<compute_pool_name>'
//	  QUERY_WAREHOUSE = '<warehouse_name>'
//	  RUNTIME = '<runtime_version>'
//	  [ ARGUMENTS = '<string_of_args>' ]
//	  [ REQUIREMENTS_FILE = '<path/to/requirements.txt>' ]
//	  [ ARTIFACT_REPOSITORIES = ( <repository_name> [ , ... ] ) ]
//	  [ EXTERNAL_ACCESS_INTEGRATIONS = ( <integration_name> [ , ... ] ) ]
//	  [ SECRETS = ( <database_name>.<schema_name>.<secret_name> [ , ... ] ) ]
//
// https://docs.snowflake.com/en/sql-reference/sql/execute-notebook-project
//
// As in CREATE NOTEBOOK, a notebook may be named PROJECT. cur is NOTEBOOK;
// start is the EXECUTE token's Loc.
func (p *Parser) parseExecuteNotebook(start ast.Loc) (ast.Node, error) {
	p.advance() // consume NOTEBOOK
	s := &ast.ObjectStmt{Verb: "EXECUTE", Kind: "NOTEBOOK", Loc: ast.Loc{Start: start.Start}}
	project := p.curIsWord("PROJECT") && p.cur.Type != tokQuotedIdent && p.peekNext().Type != '(' &&
		p.peekNext().Type != '.'
	if project {
		p.advance() // consume PROJECT
		s.Kind = "NOTEBOOK PROJECT"
	}
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	s.Name = name
	if project {
		props, err := p.parseObjectProperties(s, "", true)
		if err != nil {
			return nil, err
		}
		if err := notebookCheckOptions(props, notebookExecuteProjectProps); err != nil {
			return nil, err
		}
		if err := p.requireObjectProperties(props, "MAIN_FILE", "COMPUTE_POOL", "QUERY_WAREHOUSE", "RUNTIME"); err != nil {
			return nil, err
		}
		s.Loc.End = p.prev.Loc.End
		return s, nil
	}
	argsStart := p.cur.Loc.Start
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	args := &ast.CopyOption{Loc: ast.Loc{Start: argsStart}, List: []*ast.Literal{}}
	for p.cur.Type != ')' {
		tok, err := p.expect(tokString)
		if err != nil {
			return nil, err
		}
		args.List = append(args.List, &ast.Literal{Kind: ast.LitString, Value: tok.Str, Loc: tok.Loc})
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
		if p.cur.Type == ')' {
			return nil, p.syntaxErrorAtCur()
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	args.Loc.End = p.prev.Loc.End
	p.addObjectKeywordsClause(s, "", argsStart).Options = []*ast.CopyOption{args}
	s.Loc.End = p.prev.Loc.End
	return s, nil
}
