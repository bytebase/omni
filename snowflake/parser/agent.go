package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// AGENT and EXTERNAL AGENT (Cortex Agents), read as an ast.ObjectStmt
// ---------------------------------------------------------------------------

func init() {
	registerObject(&objectForm{
		words:     []string{"AGENT"},
		create:    parseCreateAgent,
		alter:     parseAlterAgent,
		orReplace: true,
		temp:      true,
		secure:    true,
		// https://docs.snowflake.com/en/sql-reference/sql/alter-agent writes no
		// IF EXISTS.
		noIfExists: true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-agent
		drop: true,
		// https://docs.snowflake.com/en/sql-reference/sql/undrop-agent
		undrop: true,
		props:  objectProps{"SECURE": propBool, "COMMENT": propString, "PROFILE": propString},
		unset:  []string{},
	})
	registerObject(&objectForm{
		words:     []string{"EXTERNAL", "AGENT"},
		create:    parseCreateExternalAgent,
		alter:     parseAlterExternalAgent,
		orReplace: true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-external-agent
		drop:  true,
		props: objectProps{"COMMENT": propString},
		unset: []string{},
	})
}

// parseCreateAgent reads what follows the name of
//
//	CREATE [ OR REPLACE ] [ { TEMP | TEMPORARY } ] [ SECURE ] AGENT [ IF NOT EXISTS ] <name>
//	  [ COPY GRANTS ]
//	  [ SECURE = { TRUE | FALSE } ]
//	  [ COMMENT = '<comment>' ]
//	  [ PROFILE = '<profile_object>' ]
//	  FROM SPECIFICATION '<specification_text>'
//
//	CREATE [ OR REPLACE ] [ SECURE ] AGENT [ IF NOT EXISTS ] <name>
//	  [ COPY GRANTS ]
//	  [ COMMENT = '<comment>' ]
//	  FROM @<stage_path>
//
// https://docs.snowflake.com/en/sql-reference/sql/create-agent
//
// The FROM SPECIFICATION clause holds the specification as its one nameless
// option; a FROM clause holds the stage path the same way, verbatim in Words.
func parseCreateAgent(p *Parser, s *ast.ObjectStmt) error {
	if start := p.cur.Loc.Start; p.acceptWords("COPY", "GRANTS") {
		p.addObjectKeywordsClause(s, "COPY GRANTS", start)
	}
	props, err := p.parseObjectProperties(s, "", false)
	if err != nil {
		return err
	}
	start := p.cur.Loc.Start
	if _, err := p.expect(kwFROM); err != nil {
		return err
	}
	if p.curIs('@') {
		// The stage form takes neither TEMP/TEMPORARY nor properties other than COMMENT.
		if s.Temporary {
			return &ParseError{Position: start, End: p.cur.Loc.End, Message: "syntax error: TEMPORARY is not documented with FROM @stage"}
		}
		if props != nil {
			if err := aiCheckOptions(props.Options, objectProps{"COMMENT": propString}); err != nil {
				return err
			}
		}
		return p.parseAiStageFrom(s, start)
	}
	if err := p.expectWords("SPECIFICATION"); err != nil {
		return err
	}
	// The text may be a $$ string.
	text, err := p.expect(tokString)
	if err != nil {
		return err
	}
	spec := &ast.CopyOption{Lit: &ast.Literal{Kind: ast.LitString, Value: text.Str, Loc: text.Loc}, Loc: text.Loc}
	p.addObjectKeywordsClause(s, "FROM SPECIFICATION", start).Options = []*ast.CopyOption{spec}
	return nil
}

// parseAlterAgent reads what follows the name of
//
//	ALTER AGENT <name> SET
//	  [ COMMENT = '<string>' ] [ PROFILE = '<string>' ] [ SECURE = { TRUE | FALSE } ]
//	ALTER AGENT <name> MODIFY LIVE VERSION SET SPECIFICATION = <specification>
//	ALTER AGENT <name> SET TAG <tag_name> = '<tag_value>' [ , ... ]
//	ALTER AGENT <name> UNSET TAG <tag_name> [ , ... ]
//
// The specification is a quoted (or dollar-quoted) YAML or JSON string. The
// page writes no IF EXISTS.
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-agent
func parseAlterAgent(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	if p.acceptWords("MODIFY", "LIVE", "VERSION", "SET") {
		if !p.curIsWord("SPECIFICATION") || !p.startsObjectProperty() {
			return p.syntaxErrorAtCur()
		}
		spec, err := p.parseObjectProperty()
		if err != nil {
			return err
		}
		if !propString.matches(spec) {
			return &ParseError{Position: spec.Loc.Start, End: spec.Loc.End, Message: "syntax error: a specification is a string"}
		}
		if err := aiCheckOptions([]*ast.CopyOption{spec}, nil); err != nil {
			return err
		}
		p.addObjectKeywordsClause(s, "MODIFY LIVE VERSION SET", start).Options = []*ast.CopyOption{spec}
		return nil
	}
	if ok, err := p.parseObjectAlterCommon(s, actSet|actTags); ok || err != nil {
		return err
	}
	return p.syntaxErrorAtCur()
}

// parseCreateExternalAgent reads what follows the name of
//
//	CREATE [ OR REPLACE ] EXTERNAL AGENT [ IF NOT EXISTS ] <name>
//	  [ WITH VERSION <version_name> ]
//	  [ COMMENT = '<string_literal>' ]
//
// https://docs.snowflake.com/en/sql-reference/sql/create-external-agent
func parseCreateExternalAgent(p *Parser, s *ast.ObjectStmt) error {
	if start := p.cur.Loc.Start; p.acceptWords("WITH", "VERSION") {
		if err := p.parseAiNameClause(s, "WITH VERSION", start); err != nil {
			return err
		}
	}
	_, err := p.parseObjectProperties(s, "", false)
	return err
}

// parseAlterExternalAgent reads what follows the name of
//
//	ALTER EXTERNAL AGENT [ IF EXISTS ] <name> SET COMMENT = '<string_literal>'
//	ALTER EXTERNAL AGENT [ IF EXISTS ] <name> ADD VERSION <version_name>
//
// https://docs.snowflake.com/en/sql-reference/sql/alter-external-agent
func parseAlterExternalAgent(p *Parser, s *ast.ObjectStmt) error {
	if start := p.cur.Loc.Start; p.acceptWords("ADD", "VERSION") {
		return p.parseAiNameClause(s, "ADD VERSION", start)
	}
	if ok, err := p.parseObjectAlterCommon(s, actSet); ok || err != nil {
		return err
	}
	return p.syntaxErrorAtCur()
}

// parseAiNameClause reads one name (a version, run or column) into a clause of keywords,
// begun at start and already consumed.
func (p *Parser) parseAiNameClause(s *ast.ObjectStmt, keywords string, start int) error {
	name, err := p.parseAiSingleName()
	if err != nil {
		return err
	}
	p.addObjectKeywordsClause(s, keywords, start).Names = []*ast.ObjectName{name}
	return nil
}
