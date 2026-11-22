package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseZonemapInteger reads an attribute's integer into *v.
func (p *Parser) parseZonemapInteger(v **int) error {
	if p.cur.Type != tokICONST {
		return p.syntaxErrorAtCur()
	}
	n, err := p.parseIntValue()
	if err != nil {
		return err
	}
	*v = &n
	return nil
}

// parseCreateMaterializedZonemapStmt parses a CREATE MATERIALIZED ZONEMAP
// statement after ZONEMAP.
//
// BNF: oracle/parser/bnf/CREATE-MATERIALIZED-ZONEMAP.bnf
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-MATERIALIZED-ZONEMAP.html
//
//	CREATE MATERIALIZED ZONEMAP [ IF NOT EXISTS ]
//	    [ schema. ] zonemap_name
//	    [ zonemap_attributes ]
//	    [ zonemap_refresh_clause ]
//	    [ { ENABLE | DISABLE } PRUNING ]
//	    { create_zonemap_on_table | create_zonemap_as_subquery }
//
//	create_zonemap_on_table:
//	    ON [ schema. ] table ( column [, column ]... )
//	create_zonemap_as_subquery:
//	    [ ( column_alias [, column_alias ]... ) ] AS query_block
//	zonemap_attributes:
//	    [ TABLESPACE tablespace_name ] [ SCALE integer ]
//	    [ PCTFREE integer ] [ PCTUSED integer ] [ { CACHE | NOCACHE } ]
//	zonemap_refresh_clause:
//	    REFRESH [ { FAST | COMPLETE | FORCE } ]
//	    [ { ON DEMAND | ON COMMIT | ON LOAD | ON DATA MOVEMENT | ON LOAD DATA MOVEMENT } ]
func (p *Parser) parseCreateMaterializedZonemapStmt(start int) (*nodes.CreateMaterializedZonemapStmt, error) {
	stmt := &nodes.CreateMaterializedZonemapStmt{Loc: nodes.Loc{Start: start}}

	if p.cur.Type == kwIF {
		if p.peekNext().Type == kwEXISTS {
			return nil, p.syntaxErrorAtCur()
		}
		if p.peekNext().Type == kwNOT {
			ifTok := p.cur
			p.advance() // IF
			p.advance() // NOT
			if p.cur.Type != kwEXISTS {
				return nil, p.syntaxErrorAtTok(ifTok)
			}
			p.advance() // EXISTS
			stmt.IfNotExists = true
		}
	}
	if !p.isCreateTableNameStart() {
		return nil, p.syntaxErrorAtCur()
	}
	var err error
	if stmt.Name, err = p.parseObjectName(); err != nil {
		return nil, err
	}

	for {
		switch {
		case p.cur.Type == kwTABLESPACE:
			p.advance()
			stmt.Tablespace = p.cur.Str
			if err := p.parseCreateTableName(); err != nil {
				return nil, err
			}
			continue
		case p.isKeywordStr("SCALE"):
			p.advance()
			if err := p.parseZonemapInteger(&stmt.Scale); err != nil {
				return nil, err
			}
			continue
		case p.cur.Type == kwPCTFREE:
			p.advance()
			if err := p.parseZonemapInteger(&stmt.PctFree); err != nil {
				return nil, err
			}
			continue
		case p.isKeywordStr("PCTUSED"):
			p.advance()
			if err := p.parseZonemapInteger(&stmt.PctUsed); err != nil {
				return nil, err
			}
			continue
		case p.cur.Type == kwCACHE:
			p.advance()
			stmt.Cache = true
			continue
		case p.cur.Type == kwNOCACHE:
			p.advance()
			stmt.NoCache = true
			continue
		case p.isKeywordStr("REFRESH"):
			p.advance()
			if stmt.RefreshMethod, stmt.RefreshOn, err = p.parseZonemapRefresh(false); err != nil {
				return nil, err
			}
			continue
		case p.cur.Type == kwENABLE, p.cur.Type == kwDISABLE:
			enable := p.cur.Type == kwENABLE
			p.advance()
			if err := p.expectCreateTableWord("PRUNING"); err != nil {
				return nil, err
			}
			stmt.EnablePruning, stmt.DisablePruning = enable, !enable
			continue
		}
		break
	}

	switch {
	case p.cur.Type == kwON:
		p.advance() // ON
		if !p.isCreateTableNameStart() {
			return nil, p.syntaxErrorAtCur()
		}
		if stmt.OnTable, err = p.parseObjectName(); err != nil {
			return nil, err
		}
		if stmt.OnColumns, err = p.parseGraphNameListUnique(true); err != nil {
			return nil, err
		}
	case p.cur.Type == '(' || p.cur.Type == kwAS:
		if p.cur.Type == '(' {
			if err := p.expectCreateTableToken('('); err != nil {
				return nil, err
			}
			firstTok := map[string]Token{}
			for {
				tok := p.cur
				if err := p.parseCreateTableName(); err != nil {
					return nil, err
				}
				if prev, dup := firstTok[tok.Str]; dup {
					return nil, p.syntaxErrorAtTok(prev)
				}
				firstTok[tok.Str] = tok
				stmt.ColumnAliases = append(stmt.ColumnAliases, tok.Str)
				if p.cur.Type != ',' {
					break
				}
				p.advance() // ','
			}
			if err := p.expectCreateTableToken(')'); err != nil {
				return nil, err
			}
		}
		if err := p.expectCreateTableToken(kwAS); err != nil {
			return nil, err
		}
		if stmt.AsQuery, err = p.parseSelectStmt(); err != nil {
			return nil, err
		}
		if stmt.AsQuery == nil || !selectStmtHasTargets(stmt.AsQuery) {
			return nil, p.syntaxErrorAtCur()
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseZonemapRefresh reads a zonemap_refresh_clause after REFRESH, which
// names a method, a time, or both: Oracle 23ai refuses REFRESH alone.
// ON STATEMENT is ALTER's alone.
//
//	REFRESH [ { FAST | COMPLETE | FORCE } ]
//	[ { ON DEMAND | ON COMMIT | ON LOAD | ON DATA MOVEMENT | ON LOAD DATA MOVEMENT
//	  | ON STATEMENT } ]
func (p *Parser) parseZonemapRefresh(alter bool) (string, string, error) {
	method := ""
	if p.isKeywordStr("FAST") || p.isKeywordStr("COMPLETE") || p.isKeywordStr("FORCE") {
		method = p.cur.Str
		p.advance()
	}
	// After a method, an ON that starts no refresh time is CREATE's
	// create_zonemap_on_table (REFRESH FAST ON t (c)).
	if p.cur.Type != kwON || method != "" && !alter && !p.isZonemapRefreshTimingNext() {
		if method == "" {
			return "", "", p.syntaxErrorAtCur()
		}
		return method, "", nil
	}
	p.advance() // ON
	refreshOn := ""
	switch {
	case p.isKeywordStr("DEMAND"), p.cur.Type == kwCOMMIT:
		refreshOn = "ON " + p.cur.Str
		p.advance()
	case alter && p.isKeywordStr("STATEMENT"):
		refreshOn = "ON STATEMENT"
		p.advance()
	case p.isKeywordStr("LOAD"):
		p.advance()
		refreshOn = "ON LOAD"
		if !alter && p.isKeywordStr("DATA") {
			p.advance()
			if err := p.expectCreateTableWord("MOVEMENT"); err != nil {
				return "", "", err
			}
			refreshOn = "ON LOAD DATA MOVEMENT"
		}
	case p.isKeywordStr("DATA"):
		p.advance()
		if err := p.expectCreateTableWord("MOVEMENT"); err != nil {
			return "", "", err
		}
		refreshOn = "ON DATA MOVEMENT"
	default:
		return "", "", p.syntaxErrorAtCur()
	}
	return method, refreshOn, nil
}

// parseAlterMaterializedZonemapStmt parses an ALTER MATERIALIZED ZONEMAP
// statement after ZONEMAP.
//
// BNF: oracle/parser/bnf/ALTER-MATERIALIZED-ZONEMAP.bnf
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-MATERIALIZED-ZONEMAP.html
//
//	ALTER MATERIALIZED ZONEMAP [ IF EXISTS ] [ schema. ] zonemap_name
//	    { alter_zonemap_attributes
//	    | zonemap_refresh_clause
//	    | { ENABLE | DISABLE } PRUNING
//	    | COMPILE
//	    | REBUILD
//	    | UNUSABLE
//	    } ;
//
//	alter_zonemap_attributes:
//	    [ PCTFREE integer ] [ PCTUSED integer ] [ { CACHE | NOCACHE } ]
func (p *Parser) parseAlterMaterializedZonemapStmt(start int) (*nodes.AlterMaterializedZonemapStmt, error) {
	stmt := &nodes.AlterMaterializedZonemapStmt{Loc: nodes.Loc{Start: start}}

	if p.cur.Type == kwIF && p.peekNext().Type == kwEXISTS {
		p.advance() // IF
		p.advance() // EXISTS
		stmt.IfExists = true
	}
	if !p.isCreateTableNameStart() {
		return nil, p.syntaxErrorAtCur()
	}
	var err error
	if stmt.Name, err = p.parseObjectName(); err != nil {
		return nil, err
	}

	switch {
	case p.cur.Type == kwPCTFREE, p.isKeywordStr("PCTUSED"), p.cur.Type == kwCACHE, p.cur.Type == kwNOCACHE:
		stmt.Action = "ATTRS"
		for {
			switch {
			case stmt.PctFree == nil && p.cur.Type == kwPCTFREE:
				p.advance()
				if err := p.parseZonemapInteger(&stmt.PctFree); err != nil {
					return nil, err
				}
				continue
			case stmt.PctUsed == nil && p.isKeywordStr("PCTUSED"):
				p.advance()
				if err := p.parseZonemapInteger(&stmt.PctUsed); err != nil {
					return nil, err
				}
				continue
			case !stmt.Cache && !stmt.NoCache && (p.cur.Type == kwCACHE || p.cur.Type == kwNOCACHE):
				stmt.Cache, stmt.NoCache = p.cur.Type == kwCACHE, p.cur.Type == kwNOCACHE
				p.advance()
				continue
			}
			break
		}
	case p.isKeywordStr("REFRESH"):
		p.advance()
		stmt.Action = "REFRESH"
		if stmt.RefreshMethod, stmt.RefreshOn, err = p.parseZonemapRefresh(true); err != nil {
			return nil, err
		}
	case p.cur.Type == kwENABLE, p.cur.Type == kwDISABLE:
		stmt.Action = "ENABLE_PRUNING"
		if p.cur.Type == kwDISABLE {
			stmt.Action = "DISABLE_PRUNING"
		}
		p.advance()
		if err := p.expectCreateTableWord("PRUNING"); err != nil {
			return nil, err
		}
	case p.isKeywordStr("COMPILE"), p.isKeywordStr("REBUILD"), p.isKeywordStr("UNUSABLE"):
		stmt.Action = p.cur.Str
		p.advance()
	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// isZonemapRefreshTimingNext reports whether the ON at cur starts a
// refresh time of CREATE MATERIALIZED ZONEMAP rather than its
// create_zonemap_on_table: ON DEMAND, ON COMMIT, ON LOAD or ON DATA, not
// followed by the column list or the dot of a table named so.
func (p *Parser) isZonemapRefreshTimingNext() bool {
	toks := p.peekAhead(2)
	next := toks[0]
	if next.Type != kwCOMMIT && !isCreateTableWordTok(next, "DEMAND") &&
		!isCreateTableWordTok(next, "LOAD") && !isCreateTableWordTok(next, "DATA") {
		return false
	}
	return toks[1].Type != '(' && toks[1].Type != '.'
}
