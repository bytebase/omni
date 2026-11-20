package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCreateViewStmt parses a CREATE [OR REPLACE] VIEW statement.
// The CREATE keyword has already been consumed. orReplace is set if OR REPLACE was parsed.
// CREATE MATERIALIZED VIEW is parseCreateMaterializedViewStmt's.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/CREATE-VIEW.html
//
//	CREATE [ OR REPLACE ]
//	    [ { FORCE | NOFORCE } ]
//	    [ EDITIONING ]
//	    [ { EDITIONABLE | NONEDITIONABLE } ]
//	    VIEW [ schema. ] view_name
//	    [ IF NOT EXISTS ]
//	    [ SHARING = { METADATA | DATA | EXTENDED DATA | NONE } ]
//	    [ ( column_alias [,...] ) ]
//	    [ DEFAULT COLLATION collation_name ]
//	    [ BEQUEATH { CURRENT_USER | DEFINER } ]
//	    AS select_statement
//	    [ WITH { READ ONLY | CHECK OPTION } [ CONSTRAINT constraint_name ] ]
//	    [ CONTAINER_MAP | CONTAINERS_DEFAULT ]
func (p *Parser) parseCreateViewStmt(start int, orReplace bool) (*nodes.CreateViewStmt, error) {
	stmt := &nodes.CreateViewStmt{
		OrReplace: orReplace,
		Loc:       nodes.Loc{Start: start},
	}

	// [FORCE | NOFORCE] [EDITIONING] [EDITIONABLE | NONEDITIONABLE] VIEW
	if p.cur.Type == kwFORCE {
		stmt.Force = true
		p.advance()
	} else if p.isIdentLike() && p.cur.Str == "NOFORCE" {
		stmt.NoForce = true
		p.advance()
	} else if p.isIdentLike() && p.cur.Str == "NO" {
		// NO FORCE
		stmt.NoForce = true
		p.advance()
		if p.cur.Type == kwFORCE {
			p.advance()
		}
	}

	// EDITIONING
	if p.isIdentLike() && p.cur.Str == "EDITIONING" {
		stmt.Editioning = true
		p.advance()
	}

	// EDITIONABLE | NONEDITIONABLE
	if p.isIdentLike() && p.cur.Str == "EDITIONABLE" {
		stmt.Editionable = "EDITIONABLE"
		p.advance()
	} else if p.isIdentLike() && p.cur.Str == "NONEDITIONABLE" {
		stmt.Editionable = "NONEDITIONABLE"
		p.advance()
	}

	if p.cur.Type != kwVIEW {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	return p.finishCreatePlainViewStmt(stmt)
}

// isViewConstraintStart reports whether the current position starts an
// out-of-line view constraint. CONSTRAINT is unambiguous in a view alias
// list; PRIMARY, FOREIGN, and UNIQUE are legal bare aliases (Oracle 23ai
// accepts CREATE VIEW v (primary) AS ...), so they count only when followed
// by their mandatory next token. CHECK constraints do not exist on views.
func (p *Parser) isViewConstraintStart() bool {
	switch p.cur.Type {
	case kwCONSTRAINT, kwCONSTRAINTS:
		return true
	case kwPRIMARY, kwFOREIGN:
		return p.peekNext().Type == kwKEY
	case kwUNIQUE:
		return p.peekNext().Type == '('
	}
	return false
}

// isIdentLikeStrAt checks if the given token is an identifier-like token with the given string.
func (p *Parser) isIdentLikeStrAt(tok Token, s string) bool {
	if tok.Type == tokIDENT || tok.Type >= kwAS {
		return tok.Str == s
	}
	return false
}

// parseJsonDualityViewKind parses [ RELATIONAL ] DUALITY VIEW after JSON.
// Oracle 23ai takes CREATE JSON [ RELATIONAL ] [ DUALITY ] VIEW, and
// requires DUALITY after ALTER JSON (ORA-02000 at VIEW).
func (p *Parser) parseJsonDualityViewKind(dualityRequired bool) error {
	if p.isKeywordStr("RELATIONAL") {
		p.advance()
	}
	if p.isKeywordStr("DUALITY") {
		p.advance()
	} else if dualityRequired {
		return p.syntaxErrorAtCur()
	}
	if p.cur.Type != kwVIEW {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// parseDualityViewReplication parses LOGICAL REPLICATION after ENABLE or
// DISABLE.
func (p *Parser) parseDualityViewReplication() error {
	return p.expectWords("LOGICAL", "REPLICATION")
}

// parseCreateJsonDualityViewStmt parses a CREATE JSON RELATIONAL DUALITY VIEW statement.
// Called after CREATE [OR REPLACE | IF NOT EXISTS] JSON RELATIONAL DUALITY VIEW has been consumed.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/create-json-relational-duality-view.html
//
//	CREATE [ OR REPLACE | IF NOT EXISTS ] JSON RELATIONAL DUALITY VIEW
//	    [ schema. ] view_name
//	    [ duality_view_replication_clause ]
//	    AS duality_view_subquery
func (p *Parser) parseCreateJsonDualityViewStmt(start int, orReplace, ifNotExists bool) (*nodes.CreateJsonDualityViewStmt, error) {
	stmt := &nodes.CreateJsonDualityViewStmt{
		OrReplace:   orReplace,
		IfNotExists: ifNotExists,
		Loc:         nodes.Loc{Start: start},
	}
	if !ifNotExists {
		var err error
		if stmt.IfNotExists, err = p.parseIfNotExistsAfterKind(orReplace); err != nil {
			return nil, err
		}
	}
	var parseErr611 error

	stmt.Name, parseErr611 = p.parseObjectName()
	if parseErr611 != nil {
		return nil, parseErr611
	}

	// duality_view_replication_clause: { ENABLE | DISABLE } LOGICAL REPLICATION
	if p.cur.Type == kwENABLE {
		p.advance()
		if err := p.parseDualityViewReplication(); err != nil {
			return nil, err
		}
		stmt.EnableLogicalReplication = true
	} else if p.cur.Type == kwDISABLE {
		p.advance()
		if err := p.parseDualityViewReplication(); err != nil {
			return nil, err
		}
		stmt.DisableLogicalReplication = true
	}

	// AS duality_view_subquery
	if p.cur.Type != kwAS {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	var parseErr612 error
	stmt.Query, parseErr612 = p.parseSelectStmt()
	if parseErr612 != nil {
		return nil, parseErr612
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseAlterJsonDualityViewStmt parses an ALTER JSON RELATIONAL DUALITY VIEW statement.
// Called after ALTER JSON RELATIONAL DUALITY VIEW has been consumed.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/23/sqlrf/alter-json-relational-duality-view.html
//
//	ALTER JSON RELATIONAL DUALITY VIEW [ schema. ] view_name
//	    { ENABLE | DISABLE } LOGICAL REPLICATION
func (p *Parser) parseAlterJsonDualityViewStmt(start int) (*nodes.AlterJsonDualityViewStmt, error) {
	stmt := &nodes.AlterJsonDualityViewStmt{
		Loc: nodes.Loc{Start: start},
	}
	var parseErr613 error

	stmt.Name, parseErr613 = p.parseObjectName()
	if parseErr613 != nil {
		return nil, parseErr613
	}

	if p.cur.Type == kwENABLE {
		p.advance()
		if err := p.parseDualityViewReplication(); err != nil {
			return nil, err
		}
		stmt.Action = "ENABLE_LOGICAL_REPLICATION"
	} else if p.cur.Type == kwDISABLE {
		p.advance()
		if err := p.parseDualityViewReplication(); err != nil {
			return nil, err
		}
		stmt.Action = "DISABLE_LOGICAL_REPLICATION"
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}
