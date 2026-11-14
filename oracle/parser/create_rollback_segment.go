package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCreateRollbackSegmentStmt parses a CREATE ROLLBACK SEGMENT statement.
// It is called after CREATE [ PUBLIC ] ROLLBACK SEGMENT; public is the
// PUBLIC, or nil.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-ROLLBACK-SEGMENT.html
//
//	CREATE [ PUBLIC ] ROLLBACK SEGMENT rollback_segment
//	    [ TABLESPACE tablespace ]
//	    [ storage_clause ]
//
// Oracle 23ai takes the two clauses in either order and the storage clause
// more than once, but TABLESPACE once (ORA-02215). The segment takes no
// schema (ORA-02176) and no reserved word (ORA-02175).
func (p *Parser) parseCreateRollbackSegmentStmt(start int, public *Token) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "CREATE",
		ObjectType: nodes.OBJECT_ROLLBACK_SEGMENT,
		Loc:        nodes.Loc{Start: start},
	}
	opts := &nodes.List{}
	if public != nil {
		opts.Items = append(opts.Items, &nodes.DDLOption{Key: "PUBLIC", Loc: nodes.Loc{Start: public.Loc, End: public.End}})
	}
	name, err := p.parseUnqualifiedName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	tablespace := false
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		optStart := p.pos()
		switch {
		case p.cur.Type == kwTABLESPACE && !tablespace:
			tablespace = true
			p.advance() // consume TABLESPACE
			ts, err := p.parseUnqualifiedName()
			if err != nil {
				return nil, err
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "TABLESPACE", Value: ts.Name, Loc: nodes.Loc{Start: optStart, End: p.prev.End}})
		case p.isKeywordStr("STORAGE"):
			if err := p.parseCreateTableStorageClause(); err != nil {
				return nil, err
			}
			opts.Items = append(opts.Items, &nodes.DDLOption{Key: "STORAGE", Value: p.source[optStart:p.prev.End], Loc: nodes.Loc{Start: optStart, End: p.prev.End}})
		default:
			return nil, p.syntaxErrorAtCur()
		}
	}

	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterRollbackSegmentStmt parses an ALTER ROLLBACK SEGMENT statement.
// It is called after ALTER ROLLBACK SEGMENT.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-ROLLBACK-SEGMENT.html
//
//	ALTER ROLLBACK SEGMENT rollback_segment
//	    { ONLINE | OFFLINE | storage_clause | SHRINK [ TO size_clause ] }
//
// Oracle 23ai takes one action and nothing after it (ORA-03048); the segment
// takes no schema (ORA-02244) and no reserved word (ORA-02245).
func (p *Parser) parseAlterRollbackSegmentStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "ALTER",
		ObjectType: nodes.OBJECT_ROLLBACK_SEGMENT,
		Loc:        nodes.Loc{Start: start},
	}
	name, err := p.parseUnqualifiedName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	optStart := p.pos()
	var key, value string
	switch {
	case p.isKeywordStr("ONLINE"), p.isKeywordStr("OFFLINE"):
		key = p.cur.Str
		p.advance()
	case p.isKeywordStr("STORAGE"):
		if err := p.parseCreateTableStorageClause(); err != nil {
			return nil, err
		}
		key, value = "STORAGE", p.source[optStart:p.prev.End]
	case p.isKeywordStr("SHRINK"):
		key = "SHRINK"
		p.advance() // consume SHRINK
		if p.cur.Type == kwTO {
			p.advance() // consume TO
			sizeStart := p.pos()
			if err := p.parseCreateTableSize(); err != nil {
				return nil, err
			}
			value = p.source[sizeStart:p.prev.End]
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Options = &nodes.List{Items: []nodes.Node{
		&nodes.DDLOption{Key: key, Value: value, Loc: nodes.Loc{Start: optStart, End: p.prev.End}},
	}}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}
