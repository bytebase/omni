package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/tidb/ast"
)

// parseCreateViewStmt parses a CREATE VIEW statement.
// The caller has already consumed CREATE [OR REPLACE].
//
// Ref: https://dev.mysql.com/doc/refman/8.0/en/create-view.html
//
//	CREATE
//	    [OR REPLACE]
//	    [ALGORITHM = {UNDEFINED | MERGE | TEMPTABLE}]
//	    [DEFINER = user]
//	    [SQL SECURITY { DEFINER | INVOKER }]
//	    VIEW view_name [(column_list)]
//	    AS select_statement
//	    [WITH [CASCADED | LOCAL] CHECK OPTION]
func (p *Parser) parseCreateViewStmt(orReplace bool) (*nodes.CreateViewStmt, error) {
	start := p.pos()
	stmt := &nodes.CreateViewStmt{
		Loc:       nodes.Loc{Start: start},
		OrReplace: orReplace,
	}

	// ALGORITHM = {UNDEFINED | MERGE | TEMPTABLE}
	if p.cur.Type == kwALGORITHM {
		p.advance()
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		name, _, err := p.parseKeywordOrIdent()
		if err != nil {
			return nil, err
		}
		stmt.Algorithm = name
	}

	// DEFINER = user
	if p.cur.Type == kwDEFINER {
		stmt.Definer = p.parseDefinerValue()
	}

	// SQL SECURITY { DEFINER | INVOKER }
	if p.cur.Type == kwSQL {
		p.advance()
		if p.cur.Type != kwSECURITY {
			return nil, &ParseError{
				Message:  "expected SECURITY after SQL",
				Position: p.cur.Loc,
			}
		}
		p.advance()
		name, _, err := p.parseKeywordOrIdent()
		if err != nil {
			return nil, err
		}
		stmt.SqlSecurity = name
	}

	// VIEW keyword
	if _, err := p.expect(kwVIEW); err != nil {
		return nil, err
	}

	// Completion: after CREATE VIEW, identifier context (no specific candidates).
	p.checkCursor()
	if p.collectMode() {
		// No specific candidates — user defines a new view name.
		return nil, &ParseError{Message: "collecting"}
	}

	// View name
	ref, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	stmt.Name = ref

	// Optional column list
	if p.cur.Type == '(' {
		cols, err := p.parseParenIdentList()
		if err != nil {
			return nil, err
		}
		stmt.Columns = cols
	}

	// AS
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}

	// Completion: after AS, offer SELECT keyword.
	p.checkCursor()
	if p.collectMode() {
		p.addTokenCandidate(kwSELECT)
		return nil, &ParseError{Message: "collecting"}
	}

	// SELECT statement — capture raw text
	selectStart := p.pos()
	sel, err := p.parseSelectStmt()
	if err != nil {
		return nil, err
	}
	stmt.Select = sel
	stmt.SelectText = strings.TrimSpace(p.inputText(selectStart, p.pos()))

	// WITH [CASCADED | LOCAL] CHECK OPTION
	if p.cur.Type == kwWITH {
		p.advance()
		checkOption := "CASCADED" // default
		if p.cur.Type == kwCASCADED {
			p.advance()
		} else if p.cur.Type == kwLOCAL {
			checkOption = "LOCAL"
			p.advance()
		}
		if p.cur.Type == kwCHECK {
			p.advance()
			if p.cur.Type == kwOPTION {
				p.advance()
			}
		}
		stmt.CheckOption = checkOption
	}

	stmt.Loc.End = p.pos()
	return stmt, nil
}
