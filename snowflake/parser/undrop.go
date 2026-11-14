package parser

import (
	"github.com/bytebase/omni/snowflake/ast"
)

// undropObject is an object UNDROP restores beyond the ones drop.go reads.
type undropObject struct {
	words  []string
	kind   ast.UndropObjectKind
	simple bool // the name is one part (grammar ID)
}

// undropObjects lists them, each with its page.
var undropObjects = []undropObject{
	// https://docs.snowflake.com/en/sql-reference/sql/undrop-account
	{words: []string{"ACCOUNT"}, kind: ast.UndropAccount, simple: true},
	// https://docs.snowflake.com/en/sql-reference/sql/undrop-agent
	{words: []string{"AGENT"}, kind: ast.UndropAgent},
	// https://docs.snowflake.com/en/sql-reference/sql/undrop-external-volume
	{words: []string{"EXTERNAL", "VOLUME"}, kind: ast.UndropExternalVolume, simple: true},
	// https://docs.snowflake.com/en/sql-reference/sql/undrop-notebook
	{words: []string{"NOTEBOOK"}, kind: ast.UndropNotebook},
	// https://docs.snowflake.com/en/sql-reference/sql/undrop-snapshot
	{words: []string{"SNAPSHOT"}, kind: ast.UndropSnapshot},
	// https://docs.snowflake.com/en/sql-reference/sql/undrop-stage
	{words: []string{"STAGE"}, kind: ast.UndropStage},
	// https://docs.snowflake.com/en/sql-reference/sql/undrop-streamlit
	{words: []string{"STREAMLIT"}, kind: ast.UndropStreamlit},
	// https://docs.snowflake.com/en/sql-reference/sql/undrop-type
	{words: []string{"TYPE"}, kind: ast.UndropType},
	// https://docs.snowflake.com/en/sql-reference/sql/undrop-variable
	{words: []string{"VARIABLE"}, kind: ast.UndropVariable},
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

// parseUndropMore reads UNDROP <o> once UNDROP is consumed: o's keywords, its
// name, and UNDROP SNAPSHOT's RENAME TO:
//
//	UNDROP SNAPSHOT { <name> | IDENTIFIER( <id> ) } [ RENAME TO <new_snapshot_name> ]
func (p *Parser) parseUndropMore(o *undropObject, start ast.Loc) (ast.Node, error) {
	p.acceptWords(o.words...)
	stmt := &ast.UndropStmt{Kind: o.kind, Loc: ast.Loc{Start: start.Start}}
	var err error
	if o.simple {
		stmt.Name, err = p.parseObjectSimpleName()
	} else {
		stmt.Name, err = p.parseObjectName()
	}
	if err != nil {
		return nil, err
	}
	if o.kind == ast.UndropSnapshot && p.cur.Type == kwRENAME {
		p.advance() // consume RENAME
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		if stmt.NewName, err = p.parseObjectSimpleName(); err != nil {
			return nil, err
		}
	}
	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}
