package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// An object read by its documented shape (object.go) declares UNDROP with its
// objectForm's undrop flag. The objects below have parsers of their own, so
// they declare it here; each is read into the same ast.UndropObject.

// undropObject is one such object.
type undropObject struct {
	words  []string
	simple bool // the name is one part (grammar ID)
}

// undropObjects lists them, each with its page.
var undropObjects = []undropObject{
	// https://docs.snowflake.com/en/sql-reference/sql/undrop-account
	{words: []string{"ACCOUNT"}, simple: true},
	// https://docs.snowflake.com/en/sql-reference/sql/undrop-external-volume
	{words: []string{"EXTERNAL", "VOLUME"}, simple: true},
	// https://docs.snowflake.com/en/sql-reference/sql/undrop-stage
	{words: []string{"STAGE"}},
}

// startsUndropObject reports the undropObjects entry whose keywords start at
// cur, after UNDROP, or nil. It consumes nothing.
func (p *Parser) startsUndropObject() *undropObject {
	for i := range undropObjects {
		if p.atWords(undropObjects[i].words...) {
			return &undropObjects[i]
		}
	}
	return nil
}

// parseUndropOther reads UNDROP <o> <name> once UNDROP is consumed.
func (p *Parser) parseUndropOther(o *undropObject, start ast.Loc) (ast.Node, error) {
	p.acceptWords(o.words...)
	stmt := &ast.UndropStmt{Kind: ast.UndropObject, Object: strings.Join(o.words, " "), Loc: ast.Loc{Start: start.Start}}
	var err error
	if o.simple {
		stmt.Name, err = p.parseObjectSimpleName()
	} else {
		stmt.Name, err = p.parseObjectName()
	}
	if err != nil {
		return nil, err
	}
	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseUndropRename reads the RENAME TO <new_name> an object's UNDROP may
// carry (objectForm.undropRename):
//
//	UNDROP SNAPSHOT { <name> | IDENTIFIER( <id> ) } [ RENAME TO <new_snapshot_name> ]
//
// https://docs.snowflake.com/en/sql-reference/sql/undrop-snapshot
func (p *Parser) parseUndropRename(stmt *ast.UndropStmt) error {
	if p.cur.Type != kwRENAME {
		return nil
	}
	p.advance() // consume RENAME
	if _, err := p.expect(kwTO); err != nil {
		return err
	}
	name, err := p.parseObjectSimpleName()
	if err != nil {
		return err
	}
	stmt.NewName = name
	stmt.Loc.End = p.prev.Loc.End
	return nil
}
