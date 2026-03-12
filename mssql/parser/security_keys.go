// Package parser - security_keys.go implements parsing for security key/cert/credential statements.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseSecurityKeyStmt parses CREATE/ALTER/DROP statements for security objects:
// MASTER KEY, SYMMETRIC KEY, ASYMMETRIC KEY, CERTIFICATE, CREDENTIAL.
// Also handles OPEN/CLOSE SYMMETRIC KEY and BACKUP CERTIFICATE.
//
// These statements are parsed generically: we capture the action, object type,
// name, and skip remaining tokens as options.
func (p *Parser) parseSecurityKeyStmt(action string) *nodes.SecurityKeyStmt {
	loc := p.pos()
	stmt := &nodes.SecurityKeyStmt{
		Action: action,
		Loc:    nodes.Loc{Start: loc},
	}

	// Determine object type
	switch {
	case p.matchIdentCI("MASTER"):
		// MASTER KEY
		p.match(kwKEY)
		stmt.ObjectType = "MASTER KEY"
	case p.matchIdentCI("SYMMETRIC"):
		// SYMMETRIC KEY key_name
		p.match(kwKEY)
		stmt.ObjectType = "SYMMETRIC KEY"
		name, _ := p.parseIdentifier()
		stmt.Name = name
	case p.matchIdentCI("ASYMMETRIC"):
		// ASYMMETRIC KEY key_name
		p.match(kwKEY)
		stmt.ObjectType = "ASYMMETRIC KEY"
		name, _ := p.parseIdentifier()
		stmt.Name = name
	case p.matchIdentCI("CERTIFICATE"):
		stmt.ObjectType = "CERTIFICATE"
		name, _ := p.parseIdentifier()
		stmt.Name = name
	case p.matchIdentCI("CREDENTIAL"):
		stmt.ObjectType = "CREDENTIAL"
		name, _ := p.parseIdentifier()
		stmt.Name = name
	default:
		stmt.Loc.End = p.pos()
		return stmt
	}

	// Consume remaining tokens as options until we hit a statement boundary
	p.skipSecurityKeyOptions(stmt)

	stmt.Loc.End = p.pos()
	return stmt
}

// parseOpenSymmetricKeyStmt parses OPEN SYMMETRIC KEY statements.
//
//	OPEN SYMMETRIC KEY key_name DECRYPTION BY ...
func (p *Parser) parseOpenSymmetricKeyStmt() *nodes.SecurityKeyStmt {
	loc := p.pos()
	p.advance() // consume OPEN
	stmt := &nodes.SecurityKeyStmt{
		Action: "OPEN",
		Loc:    nodes.Loc{Start: loc},
	}

	if p.matchIdentCI("SYMMETRIC") {
		p.match(kwKEY)
		stmt.ObjectType = "SYMMETRIC KEY"
	}

	name, _ := p.parseIdentifier()
	stmt.Name = name

	p.skipSecurityKeyOptions(stmt)

	stmt.Loc.End = p.pos()
	return stmt
}

// parseCloseSymmetricKeyStmt parses CLOSE SYMMETRIC KEY or CLOSE ALL SYMMETRIC KEYS.
//
//	CLOSE { SYMMETRIC KEY key_name | ALL SYMMETRIC KEYS }
func (p *Parser) parseCloseSymmetricKeyStmt() *nodes.SecurityKeyStmt {
	loc := p.pos()
	p.advance() // consume CLOSE
	stmt := &nodes.SecurityKeyStmt{
		Action: "CLOSE",
		Loc:    nodes.Loc{Start: loc},
	}

	if p.cur.Type == kwALL {
		p.advance() // consume ALL
		p.matchIdentCI("SYMMETRIC")
		p.matchIdentCI("KEYS")
		stmt.ObjectType = "ALL SYMMETRIC KEYS"
	} else if p.matchIdentCI("SYMMETRIC") {
		p.match(kwKEY)
		stmt.ObjectType = "SYMMETRIC KEY"
		name, _ := p.parseIdentifier()
		stmt.Name = name
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseBackupCertificateStmt parses BACKUP CERTIFICATE cert_name TO FILE = 'path' ...
func (p *Parser) parseBackupCertificateStmt() *nodes.SecurityKeyStmt {
	loc := p.pos()
	p.advance() // consume BACKUP
	stmt := &nodes.SecurityKeyStmt{
		Action: "BACKUP",
		Loc:    nodes.Loc{Start: loc},
	}

	if p.matchIdentCI("CERTIFICATE") {
		stmt.ObjectType = "CERTIFICATE"
		name, _ := p.parseIdentifier()
		stmt.Name = name
	} else if p.matchIdentCI("MASTER") {
		p.match(kwKEY)
		stmt.ObjectType = "MASTER KEY"
	}

	p.skipSecurityKeyOptions(stmt)

	stmt.Loc.End = p.pos()
	return stmt
}

// skipSecurityKeyOptions consumes remaining tokens until a statement boundary,
// collecting option strings.
func (p *Parser) skipSecurityKeyOptions(stmt *nodes.SecurityKeyStmt) {
	var opts []nodes.Node
	var buf strings.Builder
	depth := 0

	for p.cur.Type != tokEOF && p.cur.Type != ';' {
		// Break on keywords that start new statements (unless nested in parens)
		if depth == 0 && isStmtStart(p.cur.Type) {
			break
		}
		if p.cur.Type == '(' {
			depth++
		} else if p.cur.Type == ')' {
			if depth > 0 {
				depth--
			}
		}
		if buf.Len() > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(p.cur.Str)
		if buf.Len() == 0 {
			// For single-char tokens with empty Str
			buf.WriteRune(rune(p.cur.Type))
		}
		p.advance()
	}

	if buf.Len() > 0 {
		opts = append(opts, &nodes.String{Str: buf.String()})
	}
	if len(opts) > 0 {
		stmt.Options = &nodes.List{Items: opts}
	}
}

// isStmtStart returns true if the token type could start a new statement.
func isStmtStart(t int) bool {
	switch t {
	case kwSELECT, kwINSERT, kwUPDATE, kwDELETE, kwCREATE, kwALTER, kwDROP,
		kwGRANT, kwREVOKE, kwDENY, kwEXEC, kwEXECUTE, kwDECLARE, kwSET,
		kwIF, kwWHILE, kwBEGIN, kwRETURN, kwGOTO, kwWAITFOR,
		kwPRINT, kwRAISERROR, kwTHROW, kwUSE, kwGO,
		kwCOMMIT, kwROLLBACK, kwSAVE, kwTRUNCATE, kwMERGE,
		kwOPEN, kwFETCH, kwCLOSE, kwDEALLOCATE, kwBACKUP:
		return true
	}
	return false
}
