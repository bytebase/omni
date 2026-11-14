package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// MCP SERVER and GATEWAY, objects defined by a specification, read as
// ast.ObjectStmt
// ---------------------------------------------------------------------------

func init() {
	registerObject(&objectForm{
		words:     []string{"MCP", "SERVER"},
		create:    parseFromSpecification,
		orReplace: true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-mcp-server
		drop: true,
	})
	registerObject(&objectForm{
		words:     []string{"GATEWAY"},
		create:    parseFromSpecification,
		alter:     parseFromSpecification,
		orReplace: true,
		// https://docs.snowflake.com/en/sql-reference/sql/drop-gateway
		drop: true,
	})
}

// parseFromSpecification reads FROM SPECIFICATION '<text>' into a clause
// holding the text as its one nameless option, as CREATE AGENT's is held. The
// text may be a $$ string. It reads what follows the name of
//
//	CREATE [ OR REPLACE ] MCP SERVER [ IF NOT EXISTS ] <name>
//	  FROM SPECIFICATION $$<specification_yaml>$$
//	CREATE [ OR REPLACE ] GATEWAY [ IF NOT EXISTS ] <name>
//	  FROM SPECIFICATION <specification_text>
//	ALTER GATEWAY [ IF EXISTS ] <name> FROM SPECIFICATION <specification_text>
//
// https://docs.snowflake.com/en/sql-reference/sql/create-mcp-server
// https://docs.snowflake.com/en/sql-reference/sql/create-gateway
// https://docs.snowflake.com/en/sql-reference/sql/alter-gateway
func parseFromSpecification(p *Parser, s *ast.ObjectStmt) error {
	start := p.cur.Loc.Start
	if _, err := p.expect(kwFROM); err != nil {
		return err
	}
	if err := p.expectWords("SPECIFICATION"); err != nil {
		return err
	}
	text, err := p.expect(tokString)
	if err != nil {
		return err
	}
	spec := &ast.CopyOption{Lit: &ast.Literal{Kind: ast.LitString, Value: text.Str, Loc: text.Loc}, Loc: text.Loc}
	p.addObjectKeywordsClause(s, "FROM SPECIFICATION", start).Options = []*ast.CopyOption{spec}
	return nil
}
