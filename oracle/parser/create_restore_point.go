package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCreateRestorePointStmt parses a CREATE RESTORE POINT statement. It is
// called after CREATE [ CLEAN ] RESTORE POINT; clean is the CLEAN, or nil.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-RESTORE-POINT.html
//
//	CREATE [ CLEAN ] RESTORE POINT restore_point
//	    [ FOR PLUGGABLE DATABASE pdb_name ]
//	    [ AS OF { TIMESTAMP | SCN } expr ]
//	    [ PRESERVE | GUARANTEE FLASHBACK DATABASE ]
//
// The page lets PRESERVE and GUARANTEE FLASHBACK DATABASE stand together;
// Oracle 23ai takes one of them, in this order after the other clauses, and
// nothing after it (ORA-03048, ORA-03049). The restore point and the PDB take
// no schema (ORA-03048) and no reserved word (ORA-03050).
func (p *Parser) parseCreateRestorePointStmt(start int, clean *Token) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "CREATE",
		ObjectType: nodes.OBJECT_RESTORE_POINT,
		Loc:        nodes.Loc{Start: start},
	}
	opts := &nodes.List{}
	if clean != nil {
		opts.Items = append(opts.Items, &nodes.DDLOption{Key: "CLEAN", Loc: nodes.Loc{Start: clean.Loc, End: clean.End}})
	}
	name, err := p.parseUnqualifiedName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	add := func(key, value string, optStart int, items *nodes.List) {
		opts.Items = append(opts.Items, &nodes.DDLOption{Key: key, Value: value, Items: items, Loc: nodes.Loc{Start: optStart, End: p.prev.End}})
	}

	// [ FOR PLUGGABLE DATABASE pdb_name ]
	if p.cur.Type == kwFOR {
		optStart := p.pos()
		p.advance() // consume FOR
		if err := p.expectCreateTableWord("PLUGGABLE"); err != nil {
			return nil, err
		}
		if err := p.expectCreateTableToken(kwDATABASE); err != nil {
			return nil, err
		}
		pdb, err := p.parseUnqualifiedName()
		if err != nil {
			return nil, err
		}
		add("FOR PLUGGABLE DATABASE", pdb.Name, optStart, nil)
	}

	// [ AS OF { TIMESTAMP | SCN } expr ]
	if p.cur.Type == kwAS {
		optStart := p.pos()
		p.advance() // consume AS
		if err := p.expectCreateTableWord("OF"); err != nil {
			return nil, err
		}
		if !p.isKeywordStr("TIMESTAMP") && !p.isKeywordStr("SCN") {
			return nil, p.syntaxErrorAtCur()
		}
		kind := p.cur.Str
		p.advance()
		expr, err := p.parseRequiredExpr()
		if err != nil {
			return nil, err
		}
		if _, ok := expr.(*nodes.ColumnRef); ok && !p.isStatementEnd() &&
			!p.isKeywordStr("PRESERVE") && !p.isKeywordStr("GUARANTEE") {
			return nil, p.syntaxErrorAtNode(expr)
		}
		add("AS OF", kind, optStart, &nodes.List{Items: []nodes.Node{expr}})
	}

	// [ PRESERVE | GUARANTEE FLASHBACK DATABASE ]
	optStart := p.pos()
	switch {
	case p.isKeywordStr("PRESERVE"):
		p.advance()
		add("PRESERVE", "", optStart, nil)
	case p.isKeywordStr("GUARANTEE"):
		p.advance() // consume GUARANTEE
		if err := p.expectCreateTableToken(kwFLASHBACK); err != nil {
			return nil, err
		}
		if err := p.expectCreateTableToken(kwDATABASE); err != nil {
			return nil, err
		}
		add("GUARANTEE FLASHBACK DATABASE", "", optStart, nil)
	}

	stmt.Options = opts
	stmt.Loc.End = p.prev.End
	return stmt, nil
}
