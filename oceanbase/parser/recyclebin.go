package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// parsePurgeStmt parses OceanBase's purge_stmt, which drops an object out of
// the recycle bin, and MySQL's PURGE BINARY LOGS beside it.
//
// Ref: purge_stmt in sql_parser_mysql_mode.y (v4.3.5_CE_BP6, lines 20190-20215).
//
//	PURGE {TABLE | INDEX} relation_factor
//	PURGE {DATABASE | SCHEMA} database_factor
//	PURGE TENANT relation_name
//	PURGE RECYCLEBIN
//
// p.cur is PURGE.
func (p *Parser) parsePurgeStmt() (nodes.Node, error) {
	// The binary-log form is MySQL's, and OceanBase's grammar has no rule for
	// it: its purge_stmt names TABLE, INDEX, a database, a tenant and the
	// recycle bin, and nothing else.
	if p.peekNext().Type == kwBINARY || p.peekNext().Type == kwMASTER {
		return p.parsePurgeBinaryLogsStmt()
	}

	start := p.pos()
	p.advance() // consume PURGE

	stmt := &nodes.PurgeStmt{Loc: nodes.Loc{Start: start}}

	if p.isBareWord("RECYCLEBIN") {
		p.advance()
		stmt.Kind = "RECYCLEBIN"
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	switch {
	case p.cur.Type == kwTABLE:
		stmt.Kind = "TABLE"
		p.advance()
	case p.cur.Type == kwINDEX:
		stmt.Kind = "INDEX"
		p.advance()
	case p.cur.Type == kwDATABASE || p.isBareWord("SCHEMA"):
		stmt.Kind = "DATABASE"
		p.advance()
	case p.isBareWord("TENANT"):
		stmt.Kind = "TENANT"
		p.advance()
	default:
		return nil, p.syntaxErrorAtCur()
	}

	name, err := p.parseRelationFactor()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseFlashbackStmt parses OceanBase's flashback_stmt, which restores an
// object out of the recycle bin.
//
// Ref: flashback_stmt in sql_parser_mysql_mode.y (v4.3.5_CE_BP6, lines
// 20152-20188).
//
//	FLASHBACK TABLE relation_factor TO BEFORE DROP [RENAME TO relation_factor]
//	FLASHBACK {DATABASE | SCHEMA} database_factor TO BEFORE DROP [RENAME TO database_factor]
//	FLASHBACK TENANT relation_name TO BEFORE DROP [RENAME TO relation_name]
//
// FLASHBACK is no keyword of this lexer's, so the word this reads is an
// identifier; p.cur is it.
func (p *Parser) parseFlashbackStmt() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume FLASHBACK

	stmt := &nodes.FlashbackStmt{Loc: nodes.Loc{Start: start}}

	switch {
	case p.cur.Type == kwTABLE:
		stmt.Kind = "TABLE"
		p.advance()
	case p.cur.Type == kwDATABASE || p.isBareWord("SCHEMA"):
		stmt.Kind = "DATABASE"
		p.advance()
	case p.isBareWord("TENANT"):
		stmt.Kind = "TENANT"
		p.advance()
	default:
		return nil, p.syntaxErrorAtCur()
	}

	name, err := p.parseRelationFactor()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwBEFORE); err != nil {
		return nil, err
	}
	if !p.isBareWord("DROP") {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	if p.cur.Type == kwRENAME {
		p.advance()
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		renameTo, err := p.parseRelationFactor()
		if err != nil {
			return nil, err
		}
		stmt.RenameTo = renameTo
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseDumpMemoryStmt parses OceanBase's dump_memory_stmt: the observer's
// memory diagnostics.
//
// Ref: dump_memory_stmt in sql_parser_mysql_mode.y (v4.3.5_CE_BP6, lines
// 20236-20291).
//
//	DUMP ENTITY ALL
//	DUMP ENTITY P_ENTITY = 'string', SLOT_IDX = n
//	DUMP CHUNK ALL
//	DUMP CHUNK TENANT_ID = n, CTX_ID = name_or_string
//	DUMP CHUNK P_CHUNK = 'string'
//	DUMP MEMORY LEAK
//
// The SET OPTION arm of the same rule is read by parseSetStmt.
// DUMP is no keyword of this lexer's, so p.cur is an identifier.
func (p *Parser) parseDumpMemoryStmt() (nodes.Node, error) {
	start := p.pos()
	p.advance() // consume DUMP

	stmt := &nodes.DumpMemoryStmt{Loc: nodes.Loc{Start: start}}

	if !p.isBareWord("ENTITY") && !p.isBareWord("CHUNK") && !p.isBareWord("MEMORY") {
		return nil, p.syntaxErrorAtCur()
	}
	what := strings.ToUpper(p.cur.Str)
	p.advance()

	if what == "MEMORY" {
		if !p.isBareWord("LEAK") {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		stmt.Kind = "MEMORY_LEAK"
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	if p.cur.Type == kwALL {
		p.advance()
		stmt.Kind = what + "_ALL"
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	// The two named forms are key = value pairs, in a fixed order: P_ENTITY and
	// SLOT_IDX for DUMP ENTITY, TENANT_ID and CTX_ID for DUMP CHUNK.
	want := map[string][]string{
		"ENTITY": {"P_ENTITY", "SLOT_IDX"},
		"CHUNK":  {"TENANT_ID", "CTX_ID"},
	}[what]
	// DUMP CHUNK P_CHUNK = 'string' is the one single-valued form.
	if what == "CHUNK" && p.isBareWord("P_CHUNK") {
		want = []string{"P_CHUNK"}
	}

	for i, key := range want {
		if i > 0 {
			if _, err := p.expect(','); err != nil {
				return nil, err
			}
		}
		if !p.isBareWord(key) {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		if key == "SLOT_IDX" || key == "TENANT_ID" {
			if p.cur.Type != tokICONST {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.Args = append(stmt.Args, p.cur.Str)
			p.advance()
			continue
		}
		// CTX_ID takes a name or a string; the others take a string alone.
		if key == "CTX_ID" {
			if p.cur.Type != tokSCONST && p.cur.Type != tokIDENT {
				return nil, p.syntaxErrorAtCur()
			}
			stmt.Args = append(stmt.Args, p.cur.Str)
			p.advance()
			continue
		}
		if p.cur.Type != tokSCONST {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Args = append(stmt.Args, p.cur.Str)
		p.advance()
	}

	switch len(want) {
	case 1:
		stmt.Kind = "CHUNK_NAME"
	case 2:
		stmt.Kind = what
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseSetOptionLeakStmt parses dump_memory_stmt's SET OPTION arm, which is one
// statement of the rule and not a SET of a variable:
//
//	SET OPTION LEAK_MOD = 'string'
//	SET OPTION LEAK_RATE = n
//
// p.cur is OPTION, and atSetOptionLeak has seen the word after it.
func (p *Parser) parseSetOptionLeakStmt(start int) (nodes.Node, error) {
	p.advance() // consume OPTION
	stmt := &nodes.DumpMemoryStmt{Loc: nodes.Loc{Start: start}}
	key := strings.ToUpper(p.cur.Str)
	p.advance() // consume LEAK_MOD or LEAK_RATE
	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	if key == "LEAK_RATE" {
		if p.cur.Type != tokICONST {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Kind = "LEAK_RATE"
	} else {
		if p.cur.Type != tokSCONST {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Kind = "LEAK_MOD"
	}
	stmt.Args = append(stmt.Args, p.cur.Str)
	p.advance()
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// atSetOptionLeak reports whether the tokens are SET OPTION LEAK_MOD or SET
// OPTION LEAK_RATE, which is dump_memory_stmt and no assignment of a variable.
func (p *Parser) atSetOptionLeak() bool {
	if p.cur.Type != kwOPTION {
		return false
	}
	next := p.peekNext()
	return (next.Type == tokIDENT || next.Type >= 700) &&
		(strings.EqualFold(next.Str, "LEAK_MOD") || strings.EqualFold(next.Str, "LEAK_RATE"))
}
