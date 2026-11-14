package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// DBT PROJECT, DCM PROJECT and CODE BUNDLE, read as an ast.ObjectStmt
// ---------------------------------------------------------------------------
//
// CREATE, ALTER and DROP go through the object registry (object.go). EXECUTE
// DBT PROJECT, EXECUTE DCM PROJECT and EXECUTE CODE BUNDLE run the object, so
// they are read as an ObjectStmt with Verb "EXECUTE" (parseExecuteStmt sends
// them here). Their properties are read value by value, in the shapes the
// pages give, since some (IMPORTS, ARGUMENTS) hold more than a property run's
// generic value does.

func init() {
	registerObject(&objectForm{
		words:     []string{"DBT", "PROJECT"},
		create:    parseCreateDbtProject,
		alter:     parseAlterDbtProject,
		orReplace: true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-dbt-project
		drop:  true,
		props: dbtProjectProps,
	})
	registerObject(&objectForm{
		words:     []string{"DCM", "PROJECT"},
		create:    parseCreateDcmProject,
		alter:     parseAlterDcmProject,
		orReplace: true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-dcm-project
		drop:       true,
		props:      dcmProjectProps,
		alterProps: dcmProjectAlterProps,
	})
	// CODE BUNDLE documents CREATE alone (and SHOW and EXECUTE).
	registerObject(&objectForm{
		words:     []string{"CODE", "BUNDLE"},
		create:    parseCreateCodeBundle,
		orReplace: true,
		props:     objectProps{"COMMENT": propString},
	})
}

// dbtProjectProps are the properties CREATE DBT PROJECT and ALTER DBT PROJECT
// ... SET document, and the ones UNSET takes. The CREATE page writes
// DBT_VERSION and DEFAULT_TARGET as bare placeholders, but its examples and
// the ALTER page give both as strings.
var dbtProjectProps = objectProps{
	"DBT_VERSION":                  propString,
	"DEFAULT_TARGET":               propString,
	"DEFAULT_ENVIRONMENT":          propString,
	"EXTERNAL_ACCESS_INTEGRATIONS": propNames,
	"AUTO_COMPILE":                 propBool,
	"DEFAULT_WRITEBACK":            propBool,
	"COMMENT":                      propString,
}

// dcmProjectProps are CREATE DCM PROJECT's properties, which ALTER ... UNSET
// takes too.
var dcmProjectProps = objectProps{
	"LOG_LEVEL": propOneOf("DEBUG", "INFO", "WARN", "ERROR"),
	"COMMENT":   propString,
}

// dcmProjectAlterProps are ALTER DCM PROJECT ... SET's, whose LOG_LEVEL takes
// more levels than CREATE's.
var dcmProjectAlterProps = objectProps{
	"LOG_LEVEL": propOneOf("TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL", "OFF"),
	"COMMENT":   propString,
}

// parseCreateDbtProject reads what follows the name of
//
//	CREATE [ OR REPLACE ] DBT PROJECT [ IF NOT EXISTS ] <name>
//	  FROM '<source_location>'
//	  [ DBT_VERSION = '<version>' ] [ DEFAULT_TARGET = '<target>' ]
//	  [ DEFAULT_ENVIRONMENT = '<environment_name>' ]
//	  [ EXTERNAL_ACCESS_INTEGRATIONS = ( <integration_name> [ , ... ] ) ]
//	  [ AUTO_COMPILE = { TRUE | FALSE } ] [ DEFAULT_WRITEBACK = { TRUE | FALSE } ]
//	  [ COMMENT = '<string_literal>' ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-dbt-project
func parseCreateDbtProject(p *Parser, s *ast.ObjectStmt) error {
	if err := p.parseObjectFromString(s, "FROM"); err != nil {
		return err
	}
	if _, err := p.parseObjectProperties(s, "", false); err != nil {
		return err
	}
	return checkProjectClauses(s)
}

// parseAlterDbtProject reads what follows the name of
//
//	ALTER DBT PROJECT [ IF EXISTS ] <name> RENAME TO <new_name>
//	ALTER DBT PROJECT <name> DEPLOY FROM '<source_location>'
//	ALTER DBT PROJECT [ IF EXISTS ] <name> SET <property> [ ... ]
//	ALTER DBT PROJECT [ IF EXISTS ] <name> UNSET <property> [ , ... ]
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-dbt-project
func parseAlterDbtProject(p *Parser, s *ast.ObjectStmt) error {
	if p.atWords("DEPLOY") {
		if err := p.noIfExistsBefore(s); err != nil {
			return err
		}
		p.advance() // consume DEPLOY
		return p.parseObjectFromString(s, "DEPLOY FROM")
	}
	ok, err := p.parseObjectAlterCommon(s, actSet|actUnset|actRename)
	if err != nil {
		return err
	}
	if !ok {
		return p.syntaxErrorAtCur()
	}
	return checkProjectClauses(s)
}

// parseCreateDcmProject reads what follows the name of
//
//	CREATE [ OR REPLACE ] DCM PROJECT [ IF NOT EXISTS ] <name>
//	  [ LOG_LEVEL = { DEBUG | INFO | WARN | ERROR } ] [ COMMENT = '<string_literal>' ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-dcm-project
func parseCreateDcmProject(p *Parser, s *ast.ObjectStmt) error {
	if _, err := p.parseObjectProperties(s, "", false); err != nil {
		return err
	}
	return checkProjectClauses(s)
}

// parseAlterDcmProject reads what follows the name of
//
//	ALTER DCM PROJECT [ IF EXISTS ] <name> SET [ LOG_LEVEL = <log_level> ] [ COMMENT = '<string_literal>' ]
//	ALTER DCM PROJECT [ IF EXISTS ] <name> UNSET { LOG_LEVEL | COMMENT } [ , ... ]
//	ALTER DCM PROJECT <name> DROP DEPLOYMENT <deployment_name>
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-dcm-project
func parseAlterDcmProject(p *Parser, s *ast.ObjectStmt) error {
	if p.atWords("DROP", "DEPLOYMENT") {
		if err := p.noIfExistsBefore(s); err != nil {
			return err
		}
		start := p.cur.Loc.Start
		p.advance() // consume DROP
		p.advance() // consume DEPLOYMENT
		name, err := p.parseObjectSimpleName()
		if err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "DROP DEPLOYMENT", start).Names = []*ast.ObjectName{name}
		return nil
	}
	ok, err := p.parseObjectAlterCommon(s, actSet|actUnset)
	if err != nil {
		return err
	}
	if !ok {
		return p.syntaxErrorAtCur()
	}
	return checkProjectClauses(s)
}

// parseCreateCodeBundle reads what follows the name of
//
//	CREATE [ OR REPLACE ] CODE BUNDLE [ IF NOT EXISTS ] <database_name>.<schema_name>.<bundle_name>
//	  FROM '<source>' [ COMMENT = '<string_literal>' ]
//
// The page names the bundle by a fully qualified identifier.
//
// https://docs.snowflake.com/en/sql-reference/sql/create-code-bundle
func parseCreateCodeBundle(p *Parser, s *ast.ObjectStmt) error {
	if s.Name.Database.IsEmpty() || s.Name.Schema.IsEmpty() {
		return &ParseError{
			Position: s.Name.Loc.Start, End: s.Name.Loc.End,
			Message: "syntax error: a code bundle's name is <database>.<schema>.<name>",
		}
	}
	if err := p.parseObjectFromString(s, "FROM"); err != nil {
		return err
	}
	if _, err := p.parseObjectProperties(s, "", false); err != nil {
		return err
	}
	return checkProjectClauses(s)
}

// parseObjectFromString reads FROM '<string>' into a clause named keywords.
// cur is FROM.
func (p *Parser) parseObjectFromString(s *ast.ObjectStmt, keywords string) error {
	start := p.cur.Loc.Start
	if _, err := p.expect(kwFROM); err != nil {
		return err
	}
	opt, err := p.parseObjectStringValue()
	if err != nil {
		return err
	}
	p.addObjectKeywordsClause(s, keywords, start).Options = []*ast.CopyOption{opt}
	return nil
}

// parseObjectStringValue reads a '<string>' as a nameless option.
func (p *Parser) parseObjectStringValue() (*ast.CopyOption, error) {
	if p.cur.Type != tokString {
		return nil, p.syntaxErrorAtCur()
	}
	lit, err := p.parseCopyLiteral()
	if err != nil {
		return nil, err
	}
	return &ast.CopyOption{Lit: lit, Loc: lit.Loc}, nil
}

// noIfExistsBefore fails when s was written ALTER ... IF EXISTS, for an
// action whose page writes no IF EXISTS. cur is the action's first word.
func (p *Parser) noIfExistsBefore(s *ast.ObjectStmt) error {
	if !s.IfExists {
		return nil
	}
	return &ParseError{
		Position: p.cur.Loc.Start, End: p.cur.Loc.End,
		Message: "syntax error: ALTER " + s.Kind + " ... " + strings.ToUpper(p.cur.Str) + " takes no IF EXISTS",
	}
}

// checkProjectClauses holds a project's or bundle's clauses to their pages,
// where the generic property check does not: none documents a tag, and an
// external access integration is named by one part.
func checkProjectClauses(s *ast.ObjectStmt) error {
	for _, c := range s.Clauses {
		if c.Tags != nil {
			return &ParseError{
				Position: c.Loc.Start, End: c.Loc.End,
				Message: "syntax error: " + s.Kind + " takes no TAG",
			}
		}
		for _, o := range c.Options {
			if o.Name != "EXTERNAL_ACCESS_INTEGRATIONS" {
				continue
			}
			for _, e := range o.Group {
				if nameParts(e.Name) != 1 {
					return &ParseError{
						Position: e.Loc.Start, End: e.Loc.End,
						Message: "syntax error: " + e.Name + " takes no qualifier here",
					}
				}
			}
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// EXECUTE DBT PROJECT | DCM PROJECT | CODE BUNDLE
// ---------------------------------------------------------------------------

// startsExecuteProject reports whether cur, after EXECUTE, starts EXECUTE DBT
// PROJECT, EXECUTE DCM PROJECT or EXECUTE CODE BUNDLE. It consumes nothing.
func (p *Parser) startsExecuteProject() bool {
	return p.atWords("DBT", "PROJECT") || p.atWords("DCM", "PROJECT") || p.atWords("CODE", "BUNDLE")
}

// parseExecuteProject reads an EXECUTE of a dbt project, a DCM project or a
// code bundle. cur is the object's first keyword; start is EXECUTE's Loc.
func (p *Parser) parseExecuteProject(start ast.Loc) (ast.Node, error) {
	s := &ast.ObjectStmt{Verb: "EXECUTE", Loc: ast.Loc{Start: start.Start}}
	var err error
	switch {
	case p.acceptWords("DBT", "PROJECT"):
		s.Kind = "DBT PROJECT"
		err = p.parseExecuteDbtProject(s)
	case p.acceptWords("DCM", "PROJECT"):
		s.Kind = "DCM PROJECT"
		err = p.parseExecuteDcmProject(s)
	case p.acceptWords("CODE", "BUNDLE"):
		s.Kind = "CODE BUNDLE"
		err = p.parseExecuteCodeBundle(s)
	default:
		return nil, p.syntaxErrorAtCur()
	}
	if err != nil {
		return nil, err
	}
	s.Loc.End = p.prev.Loc.End
	return s, nil
}

// execValue reads a property's value into opt, in the one shape the page
// gives it. cur is the value's first token.
type execValue func(p *Parser, opt *ast.CopyOption) error

// dbtExecuteProps are the properties EXECUTE DBT PROJECT <name> takes.
var dbtExecuteProps = map[string]execValue{
	"ARGS":                         execString,
	"DBT_VERSION":                  execString,
	"EXTERNAL_ACCESS_INTEGRATIONS": execNames,
	"ENVIRONMENT":                  execString,
	"ENV_VARS":                     execStringPairs,
	"IMPORTS":                      execImports,
	"WRITEBACK":                    execBool,
}

// dbtExecuteWorkspaceProps are the properties EXECUTE DBT PROJECT [ FROM
// WORKSPACE <name> ] takes: PROJECT_ROOT in place of IMPORTS and WRITEBACK.
var dbtExecuteWorkspaceProps = map[string]execValue{
	"ARGS":                         execString,
	"DBT_VERSION":                  execString,
	"EXTERNAL_ACCESS_INTEGRATIONS": execNames,
	"ENVIRONMENT":                  execString,
	"ENV_VARS":                     execStringPairs,
	"PROJECT_ROOT":                 execString,
}

// parseExecuteDbtProject reads what follows EXECUTE DBT PROJECT in
//
//	EXECUTE DBT PROJECT [ IF EXISTS ] <name>
//	  [ ARGS = '<string>' ] [ DBT_VERSION = '<version>' ]
//	  [ EXTERNAL_ACCESS_INTEGRATIONS = ( <integration_name> [ , ... ] ) ]
//	  [ ENVIRONMENT = '<environment_name>' ]
//	  [ ENV_VARS = ( '<key>' = '<value>' [ , ... ] ) ]
//	  [ IMPORTS = ( { '<source_location>' | <system_function> } [ AS '<alias>' ] [ , ... ] ) ]
//	  [ WRITEBACK = { TRUE | FALSE } ]
//	EXECUTE DBT PROJECT [ IF EXISTS ] [ FROM WORKSPACE <name> ]
//	  [ ARGS = ... ] [ DBT_VERSION = ... ] [ EXTERNAL_ACCESS_INTEGRATIONS = ... ]
//	  [ ENVIRONMENT = ... ] [ ENV_VARS = ... ] [ PROJECT_ROOT = '<subdirectory_path>' ]
//
// The second form brackets FROM WORKSPACE, so it may name no project at all.
//
// https://docs.snowflake.com/en/sql-reference/sql/execute-dbt-project
func (p *Parser) parseExecuteDbtProject(s *ast.ObjectStmt) error {
	if err := p.parseIfExists(&s.IfExists); err != nil {
		return err
	}
	props := dbtExecuteWorkspaceProps
	switch {
	case p.atWords("FROM", "WORKSPACE"):
		start := p.cur.Loc.Start
		p.advance() // consume FROM
		p.advance() // consume WORKSPACE
		name, err := p.parseObjectName()
		if err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "FROM WORKSPACE", start).Names = []*ast.ObjectName{name}
	case p.cur.Type == tokEOF || p.cur.Type == ';' || p.startsObjectProperty():
		// The second form with no workspace.
	default:
		name, err := p.parseObjectName()
		if err != nil {
			return err
		}
		s.Name = name
		props = dbtExecuteProps
	}
	return p.parseExecuteProperties(s, props)
}

// parseExecuteProperties reads a run of `KEY = value` properties, each one
// props names, at most once, into a clause with no keywords.
func (p *Parser) parseExecuteProperties(s *ast.ObjectStmt, props map[string]execValue) error {
	c := &ast.ObjectClause{Loc: ast.Loc{Start: p.cur.Loc.Start}}
	seen := map[string]bool{}
	for p.startsObjectProperty() {
		key := p.cur
		name := strings.ToUpper(key.Str)
		read, ok := props[name]
		if !ok {
			return p.syntaxErrorAtCur()
		}
		if seen[name] {
			return &ParseError{
				Position: key.Loc.Start, End: key.Loc.End,
				Message: "syntax error: " + name + " is set twice",
			}
		}
		seen[name] = true
		p.advance() // consume the key
		p.advance() // consume '='
		opt := &ast.CopyOption{Name: name, Loc: ast.Loc{Start: key.Loc.Start}}
		if err := read(p, opt); err != nil {
			return err
		}
		opt.Loc.End = p.prev.Loc.End
		c.Options = append(c.Options, opt)
	}
	if len(c.Options) > 0 {
		c.Loc.End = p.prev.Loc.End
		s.Clauses = append(s.Clauses, c)
	}
	return nil
}

// execString reads '<string>'.
func execString(p *Parser, opt *ast.CopyOption) error {
	v, err := p.parseObjectStringValue()
	if err != nil {
		return err
	}
	opt.Lit = v.Lit
	return nil
}

// execBool reads TRUE or FALSE.
func execBool(p *Parser, opt *ast.CopyOption) error {
	if !isWordToken(p.cur, "TRUE") && !isWordToken(p.cur, "FALSE") {
		return p.syntaxErrorAtCur()
	}
	opt.Words = strings.ToUpper(p.advance().Str)
	return nil
}

// execNames reads ( <name> [ , ... ] ) of one-part names.
func execNames(p *Parser, opt *ast.CopyOption) error {
	return p.parseExecList(opt, func(e *ast.CopyOption) error {
		if p.cur.Type != tokIdent && p.cur.Type != tokQuotedIdent && p.cur.Type < 700 {
			return p.syntaxErrorAtCur()
		}
		w, err := p.parseObjectWord()
		if err != nil {
			return err
		}
		if nameParts(w) != 1 {
			return &ParseError{
				Position: e.Loc.Start, End: p.prev.Loc.End,
				Message: "syntax error: " + w + " takes no qualifier here",
			}
		}
		e.Name = w
		e.Bare = true
		return nil
	})
}

// execStringPairs reads ( '<key>' = '<value>' [ , ... ] ).
func execStringPairs(p *Parser, opt *ast.CopyOption) error {
	return p.parseExecList(opt, func(e *ast.CopyOption) error {
		if p.cur.Type != tokString {
			return p.syntaxErrorAtCur()
		}
		key := p.advance()
		e.Name = p.srcSlice(key.Loc.Start, key.Loc.End) // a string key keeps its quotes
		if _, err := p.expect('='); err != nil {
			return err
		}
		v, err := p.parseObjectStringValue()
		if err != nil {
			return err
		}
		e.Lit = v.Lit
		return nil
	})
}

// dbtImportFunctions are the system functions IMPORTS takes; the page says
// other system functions are not supported there.
var dbtImportFunctions = map[string]bool{
	"SYSTEM$DBT_GET_LAST_SUCCESSFUL_RUN_TARGET": true,
	"SYSTEM$DBT_GET_LAST_FAILED_RUN_TARGET":     true,
	"SYSTEM$DBT_GET_LAST_RUN_TARGET":            true,
	"SYSTEM$LOCATE_DBT_ARTIFACTS":               true,
	"SYSTEM$LOCATE_DBT_ARCHIVE":                 true,
}

// execImports reads
//
//	( { '<source_location>' | <system_function>( ... ) } [ AS '<alias>' ] [ , ... ] )
func execImports(p *Parser, opt *ast.CopyOption) error {
	return p.parseExecList(opt, func(e *ast.CopyOption) error {
		switch {
		case p.cur.Type == tokString:
			v, err := p.parseObjectStringValue()
			if err != nil {
				return err
			}
			e.Lit = v.Lit
		case p.cur.Type == tokIdent && dbtImportFunctions[strings.ToUpper(p.cur.Str)] && p.peekNext().Type == '(':
			x, err := p.parseExpr()
			if err != nil {
				return err
			}
			if _, ok := x.(*ast.FuncCallExpr); !ok {
				return &ParseError{
					Position: e.Loc.Start, End: p.prev.Loc.End,
					Message: "syntax error: IMPORTS takes a source location or a system function call",
				}
			}
			e.Expr = x
		default:
			return p.syntaxErrorAtCur()
		}
		if p.cur.Type == kwAS {
			p.advance() // consume AS
			v, err := p.parseObjectStringValue()
			if err != nil {
				return err
			}
			e.Alias = v.Lit
		}
		return nil
	})
}

// execBundleArgs reads ( { '<arg>' | :<variable> } [ , ... ] ).
func execBundleArgs(p *Parser, opt *ast.CopyOption) error {
	return p.parseExecList(opt, func(e *ast.CopyOption) error {
		if p.cur.Type == tokString {
			v, err := p.parseObjectStringValue()
			if err != nil {
				return err
			}
			e.Lit = v.Lit
			return nil
		}
		ref, err := p.parseColonVariable()
		if err != nil {
			return err
		}
		e.Name = ref
		e.Bare = true
		return nil
	})
}

// parseColonVariable reads a SQL variable reference :<name>, the colon
// written against the name, and returns it verbatim.
func (p *Parser) parseColonVariable() (string, error) {
	if p.cur.Type != ':' {
		return "", p.syntaxErrorAtCur()
	}
	colon := p.advance()
	if p.cur.Loc.Start != colon.Loc.End || (p.cur.Type != tokIdent && p.cur.Type != tokQuotedIdent && p.cur.Type < 700) {
		return "", p.syntaxErrorAtCur()
	}
	name := p.advance()
	return p.srcSlice(colon.Loc.Start, name.Loc.End), nil
}

// parseExecList reads ( <element> [ , ... ] ) into opt.Group, reading each
// element with read. cur is '('.
func (p *Parser) parseExecList(opt *ast.CopyOption, read func(e *ast.CopyOption) error) error {
	if _, err := p.expect('('); err != nil {
		return err
	}
	for {
		e := &ast.CopyOption{Loc: ast.Loc{Start: p.cur.Loc.Start}}
		if err := read(e); err != nil {
			return err
		}
		e.Loc.End = p.prev.Loc.End
		opt.Group = append(opt.Group, e)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	_, err := p.expect(')')
	return err
}

// parseExecuteCodeBundle reads what follows EXECUTE CODE BUNDLE in
//
//	EXECUTE CODE BUNDLE <name>
//	  ENTRYPOINT = '<path>'
//	  [ ARGUMENTS = ( { '<arg>' | :<variable> } [ , ... ] ) ]
//	  [ WITH SPECIFICATION $$ <yaml_spec> $$ ]
//
// https://docs.snowflake.com/en/sql-reference/sql/execute-code-bundle
func (p *Parser) parseExecuteCodeBundle(s *ast.ObjectStmt) error {
	name, err := p.parseObjectName()
	if err != nil {
		return err
	}
	s.Name = name
	c := &ast.ObjectClause{Loc: ast.Loc{Start: p.cur.Loc.Start}}
	for _, prop := range []struct {
		name string
		read execValue
		need bool
	}{{"ENTRYPOINT", execString, true}, {"ARGUMENTS", execBundleArgs, false}} {
		if !p.atWords(prop.name) || p.peekNext().Type != '=' {
			if prop.need {
				return p.syntaxErrorAtCur()
			}
			continue
		}
		key := p.advance()
		p.advance() // consume '='
		opt := &ast.CopyOption{Name: prop.name, Loc: ast.Loc{Start: key.Loc.Start}}
		if err := prop.read(p, opt); err != nil {
			return err
		}
		opt.Loc.End = p.prev.Loc.End
		c.Options = append(c.Options, opt)
	}
	c.Loc.End = p.prev.Loc.End
	s.Clauses = append(s.Clauses, c)
	if start := p.cur.Loc.Start; p.cur.Type == kwWITH && isWordToken(p.peekNext(), "SPECIFICATION") {
		p.advance() // consume WITH
		p.advance() // consume SPECIFICATION
		spec, err := p.parseObjectStringValue()
		if err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "WITH SPECIFICATION", start).Options = []*ast.CopyOption{spec}
	}
	return nil
}

// parseExecuteDcmProject reads what follows EXECUTE DCM PROJECT in
//
//	EXECUTE DCM PROJECT <name> PLAN [ DELTA ] <source> [ OUTPUT_PATH '<output_path>' ]
//	EXECUTE DCM PROJECT <name> DEPLOY [ AS <alias> ] <source>
//	EXECUTE DCM PROJECT <name> TEST ALL
//	EXECUTE DCM PROJECT <name> PREVIEW <object_name> USING CONFIGURATION <config_name>
//	  [ ENVIRONMENT ( <bind_variable> ) ] FROM '<source_path>' [ LIMIT <rows> ]
//	EXECUTE DCM PROJECT <name> PURGE [ AS <alias> ]
//
// https://docs.snowflake.com/en/sql-reference/sql/execute-dcm-project
func (p *Parser) parseExecuteDcmProject(s *ast.ObjectStmt) error {
	name, err := p.parseObjectName()
	if err != nil {
		return err
	}
	s.Name = name
	start := p.cur.Loc.Start
	switch {
	case p.acceptWords("PLAN"):
		c := p.addObjectKeywordsClause(s, "PLAN", start)
		if p.acceptWords("DELTA") {
			c.Keywords = "PLAN DELTA"
			c.Loc.End = p.prev.Loc.End
		}
		if err := p.parseDcmSource(s); err != nil {
			return err
		}
		if start := p.cur.Loc.Start; p.acceptWords("OUTPUT_PATH") {
			path, err := p.parseObjectStringValue()
			if err != nil {
				return err
			}
			p.addObjectKeywordsClause(s, "OUTPUT_PATH", start).Options = []*ast.CopyOption{path}
		}
		return nil
	case p.acceptWords("DEPLOY"):
		if err := p.parseDcmAlias(s, "DEPLOY", start); err != nil {
			return err
		}
		return p.parseDcmSource(s)
	case p.acceptWords("TEST", "ALL"):
		p.addObjectKeywordsClause(s, "TEST ALL", start)
		return nil
	case p.acceptWords("PREVIEW"):
		target, err := p.parseObjectName()
		if err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "PREVIEW", start).Names = []*ast.ObjectName{target}
		using := p.cur.Loc.Start
		if err := p.expectWords("USING", "CONFIGURATION"); err != nil {
			return err
		}
		config, err := p.parseObjectSimpleName()
		if err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "USING CONFIGURATION", using).Names = []*ast.ObjectName{config}
		if err := p.parseDcmEnvironment(s); err != nil {
			return err
		}
		if err := p.parseObjectFromString(s, "FROM"); err != nil {
			return err
		}
		if start := p.cur.Loc.Start; p.cur.Type == kwLIMIT {
			p.advance() // consume LIMIT
			if p.cur.Type != tokInt {
				return p.syntaxErrorAtCur()
			}
			lit, err := p.parseCopyLiteral()
			if err != nil {
				return err
			}
			p.addObjectKeywordsClause(s, "LIMIT", start).Options = []*ast.CopyOption{{Lit: lit, Loc: lit.Loc}}
		}
		return nil
	case p.acceptWords("PURGE"):
		return p.parseDcmAlias(s, "PURGE", start)
	}
	return p.syntaxErrorAtCur()
}

// parseDcmAlias adds the clause of keywords, which the caller consumed from
// start, and reads its optional AS <alias> into it.
func (p *Parser) parseDcmAlias(s *ast.ObjectStmt, keywords string, start int) error {
	c := p.addObjectKeywordsClause(s, keywords, start)
	if p.cur.Type != kwAS {
		return nil
	}
	p.advance() // consume AS
	alias, err := p.parseObjectSimpleName()
	if err != nil {
		return err
	}
	c.Keywords = keywords + " AS"
	c.Names = []*ast.ObjectName{alias}
	c.Loc.End = p.prev.Loc.End
	return nil
}

// parseDcmSource reads where PLAN and DEPLOY take their definitions from:
//
//	[ USING [ CONFIGURATION <config_name> ] [ ( <name> => <value> [ , ... ] ) ] ]
//	  [ ENVIRONMENT ( <bind_variable> ) ] FROM '<source_path>'
//	| FROM DEPLOYMENT { <deployment_name> | LAST | LATEST }
//
// A USING clause holds a configuration, template values or both: the page
// writes both optional, but a bare USING says nothing.
func (p *Parser) parseDcmSource(s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	if p.cur.Type == kwFROM && isWordToken(p.peekNext(), "DEPLOYMENT") {
		p.advance() // consume FROM
		p.advance() // consume DEPLOYMENT
		ref, err := p.parseObjectSimpleName()
		if err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "FROM DEPLOYMENT", start).Names = []*ast.ObjectName{ref}
		return nil
	}
	if p.cur.Type == kwUSING {
		p.advance() // consume USING
		c := &ast.ObjectClause{Keywords: "USING", Loc: ast.Loc{Start: start}}
		if p.acceptWords("CONFIGURATION") {
			config, err := p.parseObjectSimpleName()
			if err != nil {
				return err
			}
			c.Keywords = "USING CONFIGURATION"
			c.Names = []*ast.ObjectName{config}
		}
		if p.cur.Type == '(' {
			args, err := p.parseNamedArgList()
			if err != nil {
				return err
			}
			c.Args = args
		}
		if c.Names == nil && c.Args == nil {
			return p.syntaxErrorAtCur()
		}
		c.Loc.End = p.prev.Loc.End
		s.Clauses = append(s.Clauses, c)
	}
	if err := p.parseDcmEnvironment(s); err != nil {
		return err
	}
	return p.parseObjectFromString(s, "FROM")
}

// parseDcmEnvironment reads an optional ENVIRONMENT ( <bind_variable> ),
// where the variable is a session variable ($name) or a Snowflake Scripting
// variable (:name), kept verbatim as the one bare element of a list.
func (p *Parser) parseDcmEnvironment(s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	if !p.acceptWords("ENVIRONMENT") {
		return nil
	}
	opt := &ast.CopyOption{Loc: ast.Loc{Start: p.cur.Loc.Start}}
	err := p.parseExecList(opt, func(e *ast.CopyOption) error {
		if p.cur.Type == tokVariable {
			tok := p.advance()
			e.Name = p.srcSlice(tok.Loc.Start, tok.Loc.End)
		} else {
			ref, err := p.parseColonVariable()
			if err != nil {
				return err
			}
			e.Name = ref
		}
		e.Bare = true
		return nil
	})
	if err != nil {
		return err
	}
	if len(opt.Group) != 1 {
		return &ParseError{
			Position: opt.Group[1].Loc.Start, End: opt.Group[1].Loc.End,
			Message: "syntax error: ENVIRONMENT takes one bind variable",
		}
	}
	opt.Loc.End = p.prev.Loc.End
	p.addObjectKeywordsClause(s, "ENVIRONMENT", start).Options = []*ast.CopyOption{opt}
	return nil
}

// parseNamedArgList reads ( <name> => <value> [ , ... ] ). A name is any word,
// as the page's template variables are (desc => '...'). cur is '('.
func (p *Parser) parseNamedArgList() ([]*ast.CallArg, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var args []*ast.CallArg
	for {
		if (p.cur.Type != tokIdent && p.cur.Type != tokQuotedIdent && p.cur.Type < 700) || p.peekNext().Type != tokAssoc {
			return nil, p.syntaxErrorAtCur()
		}
		tok := p.advance()
		p.advance() // consume '=>'
		value, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		args = append(args, &ast.CallArg{
			Name:  ast.Ident{Name: tok.Str, Quoted: tok.Type == tokQuotedIdent, Loc: tok.Loc},
			Value: value,
			Loc:   ast.Loc{Start: tok.Loc.Start, End: p.prev.Loc.End},
		})
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return args, nil
}

// showProjectScopes are the SHOW classes a project's scope serves:
// IN DBT PROJECT and IN DCM PROJECT.
var showProjectScopes = map[string][]string{
	"DBT PROJECT": {"VERSIONS"},
	"DCM PROJECT": {"DEPLOYMENTS", "ENTITIES", "GRANTS"},
}

// startsShowProjectScope reports whether cur, after SHOW <class> IN, is DBT
// PROJECT or DCM PROJECT for a class it serves, and returns the keywords. It
// consumes nothing.
func (p *Parser) startsShowProjectScope(class string) []string {
	for _, words := range [][]string{{"DBT", "PROJECT"}, {"DCM", "PROJECT"}} {
		if !p.atWords(words...) {
			continue
		}
		for _, c := range showProjectScopes[strings.Join(words, " ")] {
			if c == class {
				return words
			}
		}
	}
	return nil
}

// checkShowProjectScope fails a SHOW in a project's scope that carries a
// filter its page does not write: only SHOW ENTITIES filters by LIKE,
// STARTS WITH and LIMIT ... FROM.
func checkShowProjectScope(stmt *ast.ShowStmt) error {
	if stmt.Scope != ast.ShowScopeOther || showProjectScopes[stmt.ScopeText] == nil {
		return nil
	}
	entities := stmt.ObjectClass == "ENTITIES"
	if stmt.Terse || stmt.History || stmt.Privileges != nil ||
		(!entities && (stmt.HasLike || stmt.HasStarts || stmt.LimitFrom != "")) {
		return &ParseError{
			Position: stmt.Loc.Start, End: stmt.Loc.End,
			Message: "syntax error: SHOW " + stmt.ObjectClass + " IN " + stmt.ScopeText + " takes no such filter",
		}
	}
	return nil
}

// parseShowProjectScope reads the { DBT | DCM } PROJECT <name> after IN of
//
//	SHOW VERSIONS IN DBT PROJECT <name> [ LIMIT <rows> ]
//	SHOW DEPLOYMENTS IN DCM PROJECT <name> [ LIMIT <rows> ]
//	SHOW ENTITIES [ LIKE '<pattern>' ] IN DCM PROJECT <name> [ STARTS WITH '<s>' ] [ LIMIT <rows> [ FROM '<s>' ] ]
//	SHOW [ FUTURE ] GRANTS IN DCM PROJECT <name> [ LIMIT <rows> ]
//
// as a ShowScopeOther. words is what startsShowProjectScope matched.
//
// https://docs.snowflake.com/en/sql-reference/sql/show-versions-in-dbt-project
// https://docs.snowflake.com/en/sql-reference/sql/show-deployments-in-dcm-project
// https://docs.snowflake.com/en/sql-reference/sql/show-entities-in-dcm-project
// https://docs.snowflake.com/en/sql-reference/sql/show-grants-in-dcm-project
func (p *Parser) parseShowProjectScope(stmt *ast.ShowStmt, words []string) error {
	p.acceptWords(words...)
	stmt.Scope = ast.ShowScopeOther
	stmt.ScopeText = strings.Join(words, " ")
	name, err := p.parseObjectName()
	if err != nil {
		return err
	}
	stmt.ScopeName = name
	return nil
}
