package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// DROP of objects whose CREATE and ALTER omni reads on another path
// ---------------------------------------------------------------------------
//
// SECRET, PASSWORD POLICY, AUTHENTICATION POLICY, GIT REPOSITORY, EXTERNAL
// VOLUME, SEMANTIC VIEW and ACCOUNT are created and altered by parsers of
// their own (create_table.go, database_schema.go, role_user_policy.go,
// semantic_view.go, ...), so they are not registered objects (object.go),
// whose CREATE and ALTER the registry would then take over. Their DROP pages
// all write DROP <kind> [ IF EXISTS ] <name>; each is read here into a
// DropObject.

// docDrop is one DROP page.
type docDrop struct {
	words []string // the object's keywords, uppercased
	// nameParts is 1 where the page names the object by a single <name>
	// (an account-level object), 0 where db.schema.name may qualify it.
	nameParts int
	// gracePeriod says the page writes GRACE_PERIOD_IN_DAYS = <integer>
	// after the name (DROP ACCOUNT).
	gracePeriod bool
}

var docDrops = []docDrop{
	// https://docs.snowflake.com/en/sql-reference/sql/drop-secret
	{words: []string{"SECRET"}},
	// https://docs.snowflake.com/en/sql-reference/sql/drop-password-policy
	{words: []string{"PASSWORD", "POLICY"}},
	// https://docs.snowflake.com/en/sql-reference/sql/drop-authentication-policy
	{words: []string{"AUTHENTICATION", "POLICY"}},
	// https://docs.snowflake.com/en/sql-reference/sql/drop-git-repository
	{words: []string{"GIT", "REPOSITORY"}},
	// https://docs.snowflake.com/en/sql-reference/sql/drop-external-volume
	{words: []string{"EXTERNAL", "VOLUME"}, nameParts: 1},
	// https://docs.snowflake.com/en/sql-reference/sql/drop-semantic-view
	{words: []string{"SEMANTIC", "VIEW"}},
	// https://docs.snowflake.com/en/sql-reference/sql/drop-account
	//
	//	DROP ACCOUNT [ IF EXISTS ] <name> GRACE_PERIOD_IN_DAYS = <integer>
	{words: []string{"ACCOUNT"}, nameParts: 1, gracePeriod: true},
}

// startsDocDrop reports the docDrop whose keywords start at cur, the word
// after DROP. It consumes nothing.
func (p *Parser) startsDocDrop() *docDrop {
	if !p.curIs(tokIdent) && p.cur.Type < 700 {
		return nil
	}
	for i := range docDrops {
		if p.atWords(docDrops[i].words...) {
			return &docDrops[i]
		}
	}
	return nil
}

// parseDocDrop reads DROP <kind> [ IF EXISTS ] <name> of d. cur is d's first
// keyword; start is the DROP token's Loc.
func (p *Parser) parseDocDrop(start ast.Loc, d *docDrop) (ast.Node, error) {
	for range d.words {
		p.advance()
	}
	node, err := p.parseDropObject(ast.DropObject, start, true, false)
	if err != nil {
		return nil, err
	}
	ds := node.(*ast.DropStmt)
	ds.Object = strings.Join(d.words, " ")
	if d.nameParts > 0 {
		if _, err := p.checkFormNameParts(&objectForm{kind: ds.Object, nameParts: d.nameParts}, ds.Name); err != nil {
			return nil, err
		}
	}
	if d.gracePeriod {
		if err := p.expectWords("GRACE_PERIOD_IN_DAYS"); err != nil {
			return nil, err
		}
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		tok, err := p.expect(tokInt)
		if err != nil {
			return nil, err
		}
		ds.GracePeriodInDays = &ast.Literal{Kind: ast.LitInt, Value: tok.Str, Ival: tok.Ival, Loc: tok.Loc}
	}
	ds.Loc.End = p.prev.Loc.End
	return ds, nil
}
