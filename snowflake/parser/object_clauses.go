package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Clauses several objects read by their documented shape share
// ---------------------------------------------------------------------------

// rejectObjectTags fails when a clause of s holds tags, for an object whose
// page documents no TAG.
func rejectObjectTags(s *ast.ObjectStmt) error {
	for _, c := range s.Clauses {
		if c.Tags != nil {
			return &ParseError{
				Position: c.Loc.Start, End: c.Loc.End,
				Message: "syntax error: " + strings.ToLower(s.Kind) + " takes no TAG",
			}
		}
	}
	return nil
}

// parseObjectStringList reads '<string>' [ , ... ], each string a nameless
// option, as a clause's list of strings is kept.
func (p *Parser) parseObjectStringList() ([]*ast.CopyOption, error) {
	var out []*ast.CopyOption
	for {
		if p.cur.Type != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		lit, err := p.parseCopyLiteral()
		if err != nil {
			return nil, err
		}
		out = append(out, &ast.CopyOption{Lit: lit, Loc: lit.Loc})
		if p.cur.Type != ',' {
			return out, nil
		}
		p.advance() // consume ','
	}
}

// parseObjectAsString reads AS '<string>' into a clause holding the string as
// a nameless option. cur is AS.
func (p *Parser) parseObjectAsString(s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	p.advance() // consume AS
	if p.cur.Type != tokString {
		return p.syntaxErrorAtCur()
	}
	lit, err := p.parseCopyLiteral()
	if err != nil {
		return err
	}
	p.addObjectKeywordsClause(s, "AS", start).Options = []*ast.CopyOption{{Lit: lit, Loc: lit.Loc}}
	return nil
}

// parseObjectFromStage reads FROM @<stage> into a clause holding the stage
// location, verbatim, as a nameless option's Words. cur is FROM.
func (p *Parser) parseObjectFromStage(s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	p.advance() // consume FROM
	if p.cur.Type != '@' {
		return p.syntaxErrorAtCur()
	}
	loc, err := p.parseStageRef()
	if err != nil {
		return err
	}
	p.addObjectKeywordsClause(s, "FROM", start).Options = []*ast.CopyOption{{Words: loc.Raw, Loc: loc.Loc}}
	return nil
}

// parseObjectAddVersion reads the [ [ IF NOT EXISTS ] <version_name> ] after
// keywords (ADD VERSION, ADD LIVE VERSION), which the caller has consumed
// from start, into a clause. The name is absent when FROM follows.
func (p *Parser) parseObjectAddVersion(s *ast.ObjectStmt, keywords string, start int) error {
	c := p.addObjectKeywordsClause(s, keywords, start)
	if p.cur.Type == kwFROM {
		return nil
	}
	ifNotExists := false
	if err := p.parseIfNotExists(&ifNotExists); err != nil {
		return err
	}
	if ifNotExists {
		c.Keywords += " IF NOT EXISTS"
	}
	name, err := p.parseObjectSimpleName()
	if err != nil {
		return err
	}
	c.Names = []*ast.ObjectName{name}
	c.Loc.End = p.prev.Loc.End
	return nil
}

// parseObjectRenameSimple reads RENAME TO <new_name> for an object whose name
// is a single <name>. cur is RENAME.
func (p *Parser) parseObjectRenameSimple(s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	p.advance() // consume RENAME
	if _, err := p.expect(kwTO); err != nil {
		return err
	}
	name, err := p.parseObjectSimpleName()
	if err != nil {
		return err
	}
	p.addObjectKeywordsClause(s, "RENAME TO", start).Names = []*ast.ObjectName{name}
	return nil
}
