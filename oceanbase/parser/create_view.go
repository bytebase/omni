package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// parseOptTableID reads opt_table_id, `TABLE_ID '=' INTNUM`, and reports whether
// it was written.
func (p *Parser) parseOptTableID() (int64, bool, error) {
	if !p.isBareWord("TABLE_ID") {
		return 0, false, nil
	}
	p.advance()
	if _, err := p.expect('='); err != nil {
		return 0, false, err
	}
	if p.cur.Type != tokICONST {
		return 0, false, p.syntaxErrorAtCur()
	}
	id := p.cur.Ival
	p.advance()
	return id, true, nil
}

// parseCreateViewStmt parses a CREATE VIEW statement.
// The caller has already consumed CREATE [OR REPLACE].
//
// Ref: create_view_stmt in OceanBase's src/sql/parser/sql_parser_mysql_mode.y
// (v4.3.5_CE_BP6, line 9717)
//
//	CREATE
//	    [OR REPLACE]
//	    [ALGORITHM = {UNDEFINED | MERGE | TEMPTABLE} [DEFINER = user]]
//	    [SQL SECURITY { DEFINER | INVOKER }]
//	    VIEW view_name [(column_list)]
//	    AS select_statement
//	    [WITH [CASCADED | LOCAL] CHECK OPTION]
func (p *Parser) parseCreateViewStmt(start int, orReplace bool) (*nodes.CreateViewStmt, error) {
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
		name, err := p.parseWordOf(viewAlgorithmWords)
		if err != nil {
			return nil, err
		}
		stmt.Algorithm = name
	}

	// DEFINER = user, only after ALGORITHM. The grammar lists opt_definer on its
	// own, but CREATE DEFINER also heads a stored program, and 4.3.5.6 reads it
	// that way: without ALGORITHM before it, CREATE [OR REPLACE] DEFINER = u VIEW
	// answers 1064 (measured 2026-11-05).
	if p.cur.Type == kwDEFINER && stmt.Algorithm != "" {
		definer, err := p.parseDefinerValue(false)
		if err != nil {
			return nil, err
		}
		stmt.Definer = definer
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
	ref, err := p.parseRelationFactor()
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

	// opt_table_id: TABLE_ID [=] n, after the column list and before AS.
	tableID, hasTableID, err := p.parseOptTableID()
	if err != nil {
		return nil, err
	}
	stmt.TableID, stmt.HasTableID = tableID, hasTableID

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
	sel, err := p.parseQueryStatement()
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
		// CHECK and OPTION are both words of the clause: `WITH CHECK`,
		// `WITH CASCADED CHECK` and `WITH LOCAL CHECK` are 1064s.
		if _, err := p.expect(kwCHECK); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwOPTION); err != nil {
			return nil, err
		}
		stmt.CheckOption = checkOption
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}
