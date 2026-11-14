package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// noauditStatementOptions are the statement options of traditional
// auditing (the SQL statement shortcuts and the additional ones of the
// AUDIT (Traditional Auditing) page), which Oracle 23ai still reads in
// NOAUDIT, with NETWORK and DIRECT_PATH LOAD.
var noauditStatementOptions = newAuditNames(
	"ALTER SYSTEM", "CLUSTER", "CONTEXT", "DATABASE LINK", "DIMENSION",
	"DIRECTORY", "INDEX", "MATERIALIZED VIEW", "NOT EXISTS", "OUTLINE",
	"PLUGGABLE DATABASE", "PROCEDURE", "PROFILE", "PUBLIC DATABASE LINK",
	"PUBLIC SYNONYM", "ROLE", "ROLLBACK SEGMENT", "SEQUENCE", "SESSION",
	"SYNONYM", "SYSTEM AUDIT", "SYSTEM GRANT", "TABLE", "TABLESPACE",
	"TRIGGER", "TYPE", "USER", "VIEW",
	"ALTER SEQUENCE", "ALTER TABLE", "COMMENT TABLE", "DELETE TABLE",
	"EXECUTE DIRECTORY", "EXECUTE PROCEDURE", "GRANT DIRECTORY",
	"GRANT PROCEDURE", "GRANT SEQUENCE", "GRANT TABLE", "GRANT TYPE",
	"INSERT TABLE", "LOCK TABLE", "READ DIRECTORY", "SELECT SEQUENCE",
	"SELECT TABLE", "UPDATE TABLE", "WRITE DIRECTORY",
	"NETWORK", "DIRECT_PATH LOAD",
)

// noauditObjectOperations are the operations NOAUDIT ... ON an object
// names (ORA-00956 at TRUNCATE or MERGE).
var noauditObjectOperations = newAuditNames(
	"ALTER", "AUDIT", "COMMENT", "CREATE", "DEBUG", "DELETE", "EXECUTE",
	"FLASHBACK", "GRANT", "INDEX", "INSERT", "LOCK", "READ", "RENAME",
	"SELECT", "UPDATE", "USE", "WRITE",
)

// isNoauditOptionEnd reports whether the current token ends a NOAUDIT
// option.
func (p *Parser) isNoauditOptionEnd() bool {
	return isNoauditOptionEndTok(p.cur)
}

// isNoauditOptionEndTok reports whether tok ends a NOAUDIT option.
func isNoauditOptionEndTok(tok Token) bool {
	switch tok.Type {
	case ',', kwON, kwBY, kwWHENEVER, ';', tokEOF:
		return true
	}
	return tok.Type != tokQIDENT && tok.Str == "CONTAINER"
}

// parseNoauditTraditional parses the traditional forms of NOAUDIT, after
// NOAUDIT.
//
//	NOAUDIT
//	    { audit_operation_clause [ auditing_by_clause ]
//	    | audit_schema_object_clause
//	    }
//	    [ WHENEVER [ NOT ] SUCCESSFUL ]
//	    [ CONTAINER = { CURRENT | ALL } ] ;
//
// Oracle 23ai reads each option as the longest run of words that names a
// statement option, a system privilege, or an operation on objects; another
// fails at its start (ORA-00956). Operations need ON and nothing else does
// (ORA-00969 where ON is missing, ORA-03048 at ON after a statement
// option). An object option takes no BY (ORA-01708) and no CONTAINER
// (ORA-46504), and a BY names users, not ACCESS or SESSION (ORA-01718).
func (p *Parser) parseNoauditTraditional(stmt *nodes.NoauditStmt) error {
	// onOK: every option may audit objects ON; bareOK: every one may stand
	// without ON.
	onOK, bareOK := true, true
	if p.cur.Type == kwALL {
		p.advance()
		action := "ALL"
		if p.isKeywordStr("STATEMENTS") || p.cur.Type == kwPRIVILEGES {
			action += " " + p.cur.Str
			p.advance()
			onOK = false
		}
		stmt.Actions = append(stmt.Actions, action)
	} else {
		for {
			itemStart := p.cur
			name, n := p.matchAuditName(noauditStatementOptions, auditSystemPrivileges, noauditObjectOperations)
			if n == 0 {
				return p.syntaxErrorAtCur()
			}
			for range n {
				p.advance()
			}
			// An option ends where the clause after it starts, or it
			// names none (ORA-00956 at TYPE in TYPE BODY).
			if !p.isNoauditOptionEnd() {
				return p.syntaxErrorAtTok(itemStart)
			}
			onOK = onOK && noauditObjectOperations.names[name]
			bareOK = bareOK && (noauditStatementOptions.names[name] || auditSystemPrivileges.names[name])
			stmt.Actions = append(stmt.Actions, name)
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}

	if p.cur.Type == kwON {
		if !onOK {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		if err := p.parseTraditionalNoauditOnClause(stmt); err != nil {
			return err
		}
		// No BY follows an object option (ORA-01708 at the user).
		if p.cur.Type == kwBY {
			p.advance()
			return p.syntaxErrorAtCur()
		}
	} else {
		// An operation names the objects it audits ON.
		if !bareOK {
			return p.syntaxErrorAtCur()
		}
		if p.cur.Type == kwBY {
			p.advance()
			for {
				if !p.isIdentLike() {
					return p.syntaxErrorAtCur()
				}
				if err := p.syntaxErrorIfReservedIdentifier(); err != nil {
					return err
				}
				stmt.ByUsers2 = append(stmt.ByUsers2, p.cur.Str)
				p.advance()
				if p.cur.Type != ',' {
					break
				}
				p.advance()
			}
		}
	}

	if p.cur.Type == kwWHENEVER {
		p.advance()
		stmt.When = "WHENEVER"
		if p.cur.Type == kwNOT {
			stmt.When += " NOT"
			p.advance()
		}
		if p.cur.Type != kwSUCCESSFUL {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		stmt.When += " SUCCESSFUL"
	}

	if p.isKeywordStr("CONTAINER") {
		if stmt.Object != nil || stmt.OnDefault || stmt.OnDirectory != "" {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		if p.cur.Type != '=' {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		if p.cur.Type != kwALL && !p.isKeywordStr("CURRENT") {
			return p.syntaxErrorAtCur()
		}
		v := p.cur.Type == kwALL
		stmt.ContainerAll = &v
		p.advance()
	}
	return nil
}
