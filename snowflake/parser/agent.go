package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// AGENT and EXTERNAL AGENT (Cortex Agents), read as an ast.ObjectStmt
// ---------------------------------------------------------------------------
//
// CREATE AGENT's modifiers before AGENT ({ TEMP | TEMPORARY }, SECURE) and its
// FROM @<stage_path> form are not read yet: ObjectStmt records neither.

func init() {
	registerObject(&objectForm{
		words:     []string{"AGENT"},
		create:    parseCreateAgent,
		alter:     parseAlterAgent,
		orReplace: true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-agent
		drop: true,
	})
	registerObject(&objectForm{
		words:     []string{"EXTERNAL", "AGENT"},
		create:    parseCreateExternalAgent,
		alter:     parseAlterExternalAgent,
		orReplace: true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-external-agent
		drop: true,
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
// https://docs.snowflake.com/en/sql-reference/sql/create-agent
//
// The FROM clause holds the specification as an option named SPECIFICATION.
func parseCreateAgent(p *Parser, s *ast.ObjectStmt) error {
	if start := p.cur.Loc.Start; p.acceptWords("COPY", "GRANTS") {
		p.addObjectKeywordsClause(s, "COPY GRANTS", start)
	}
	if _, err := p.parseObjectProperties(s, "", false); err != nil {
		return err
	}
	start := p.cur.Loc.Start
	if _, err := p.expect(kwFROM); err != nil {
		return err
	}
	spec, err := p.parseAgentSpecification()
	if err != nil {
		return err
	}
	p.addObjectKeywordsClause(s, "FROM", start).Options = []*ast.CopyOption{spec}
	return nil
}

// parseAgentSpecification reads SPECIFICATION '<specification_text>' (the
// text may be a $$ string) into an option named SPECIFICATION.
func (p *Parser) parseAgentSpecification() (*ast.CopyOption, error) {
	key := p.cur
	if err := p.expectWords("SPECIFICATION"); err != nil {
		return nil, err
	}
	text, err := p.expect(tokString)
	if err != nil {
		return nil, err
	}
	return &ast.CopyOption{
		Name: "SPECIFICATION",
		Lit:  &ast.Literal{Kind: ast.LitString, Value: text.Str, Loc: text.Loc},
		Loc:  ast.Loc{Start: key.Loc.Start, End: text.Loc.End},
	}, nil
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
	if s.IfExists {
		return &ParseError{Position: s.Name.Loc.Start, End: s.Name.Loc.End, Message: "syntax error: ALTER AGENT takes no IF EXISTS"}
	}
	start := p.cur.Loc.Start
	if p.acceptWords("MODIFY", "LIVE", "VERSION") {
		p.addObjectKeywordsClause(s, "MODIFY LIVE VERSION", start)
		setStart := p.cur.Loc.Start
		if _, err := p.expect(kwSET); err != nil {
			return err
		}
		if !p.curIsWord("SPECIFICATION") || !p.startsObjectProperty() {
			return p.syntaxErrorAtCur()
		}
		spec, err := p.parseObjectProperty()
		if err != nil {
			return err
		}
		if spec.Lit == nil || spec.Lit.Kind != ast.LitString {
			return &ParseError{Position: spec.Loc.Start, End: spec.Loc.End, Message: "syntax error: a specification is a string"}
		}
		p.addObjectKeywordsClause(s, "SET", setStart).Options = []*ast.CopyOption{spec}
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
