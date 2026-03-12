// Package parser - backup_restore.go implements T-SQL BACKUP and RESTORE statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseBackupStmt parses a BACKUP DATABASE or BACKUP LOG statement.
//
//	BACKUP { DATABASE | LOG } database_name
//	    TO { DISK | URL | TAPE } = 'path'
//	    [WITH option [, ...]]
func (p *Parser) parseBackupStmt() *nodes.BackupStmt {
	loc := p.pos()
	p.advance() // consume BACKUP

	stmt := &nodes.BackupStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// DATABASE or LOG (or identifier like LOG)
	if p.cur.Type == kwDATABASE {
		stmt.Type = "DATABASE"
		p.advance()
	} else if p.isIdentLike() && strings.EqualFold(p.cur.Str, "LOG") {
		stmt.Type = "LOG"
		p.advance()
	} else if p.isIdentLike() && strings.EqualFold(p.cur.Str, "CERTIFICATE") {
		stmt.Type = "CERTIFICATE"
		p.advance()
	} else {
		stmt.Type = "DATABASE"
	}

	// Database name (not for CERTIFICATE)
	if stmt.Type != "CERTIFICATE" {
		if p.isIdentLike() {
			stmt.Database = p.cur.Str
			p.advance()
		}
	}

	// TO { DISK | URL | TAPE | ... } = 'path'
	if p.cur.Type == kwTO {
		p.advance() // consume TO
		// consume DISK / URL / TAPE / identifier
		if p.isIdentLike() || p.cur.Type == kwFILE {
			p.advance()
		}
		// = 'path'
		if _, ok := p.match('='); ok {
			if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
				stmt.Target = p.cur.Str
				p.advance()
			}
		}
	}

	// WITH options — consume everything inside WITH (...) or comma-separated keywords
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		stmt.Options = p.consumeBackupWithOptions()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseRestoreStmt parses a RESTORE DATABASE / LOG / HEADERONLY / FILELISTONLY statement.
//
//	RESTORE { DATABASE | LOG | HEADERONLY | FILELISTONLY | VERIFYONLY | LABELONLY }
//	    [database_name]
//	    FROM { DISK | URL | TAPE } = 'path'
//	    [WITH option [, ...]]
func (p *Parser) parseRestoreStmt() *nodes.RestoreStmt {
	loc := p.pos()
	p.advance() // consume RESTORE

	stmt := &nodes.RestoreStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Determine restore type
	if p.cur.Type == kwDATABASE {
		stmt.Type = "DATABASE"
		p.advance()
	} else if p.isIdentLike() {
		upper := strings.ToUpper(p.cur.Str)
		switch upper {
		case "LOG":
			stmt.Type = "LOG"
			p.advance()
		case "HEADERONLY":
			stmt.Type = "HEADERONLY"
			p.advance()
		case "FILELISTONLY":
			stmt.Type = "FILELISTONLY"
			p.advance()
		case "VERIFYONLY":
			stmt.Type = "VERIFYONLY"
			p.advance()
		case "LABELONLY":
			stmt.Type = "LABELONLY"
			p.advance()
		default:
			stmt.Type = "DATABASE"
		}
	} else {
		stmt.Type = "DATABASE"
	}

	// Database name (optional for HEADERONLY/FILELISTONLY/VERIFYONLY/LABELONLY)
	switch stmt.Type {
	case "HEADERONLY", "FILELISTONLY", "VERIFYONLY", "LABELONLY":
		// no database name expected before FROM
	default:
		if p.isIdentLike() && p.cur.Type != kwFROM {
			stmt.Database = p.cur.Str
			p.advance()
		}
	}

	// FROM { DISK | URL | TAPE | ... } = 'path'
	if p.cur.Type == kwFROM {
		p.advance() // consume FROM
		// consume DISK / URL / TAPE / identifier
		if p.isIdentLike() || p.cur.Type == kwFILE {
			p.advance()
		}
		// = 'path'
		if _, ok := p.match('='); ok {
			if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
				stmt.Source = p.cur.Str
				p.advance()
			}
		}
	}

	// WITH options
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		stmt.Options = p.consumeBackupWithOptions()
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// consumeBackupWithOptions consumes backup/restore WITH option tokens.
// Options may be a parenthesized list or a comma-separated list of bare keywords/values.
// We capture the raw text as String nodes for each top-level option.
func (p *Parser) consumeBackupWithOptions() *nodes.List {
	var opts []nodes.Node

	if p.cur.Type == '(' {
		// Parenthesized options: consume balanced parens
		depth := 0
		var sb strings.Builder
		for p.cur.Type != tokEOF {
			if p.cur.Type == '(' {
				depth++
				if depth > 1 {
					sb.WriteString("(")
				}
				p.advance()
			} else if p.cur.Type == ')' {
				depth--
				if depth == 0 {
					p.advance()
					break
				}
				sb.WriteString(")")
				p.advance()
			} else if p.cur.Type == ',' && depth == 1 {
				s := strings.TrimSpace(sb.String())
				if s != "" {
					opts = append(opts, &nodes.String{Str: s})
				}
				sb.Reset()
				p.advance()
			} else {
				if sb.Len() > 0 {
					sb.WriteString(" ")
				}
				sb.WriteString(p.cur.Str)
				p.advance()
			}
		}
		s := strings.TrimSpace(sb.String())
		if s != "" {
			opts = append(opts, &nodes.String{Str: s})
		}
	} else {
		// Bare comma-separated options
		for {
			if p.cur.Type == tokEOF || p.cur.Type == ';' || p.isStatementStart() {
				break
			}
			var sb strings.Builder
			// Collect tokens for this option until we hit a comma or end
			for p.cur.Type != tokEOF && p.cur.Type != ';' && p.cur.Type != ',' &&
				!p.isStatementStart() {
				if sb.Len() > 0 {
					sb.WriteString(" ")
				}
				sb.WriteString(p.cur.Str)
				p.advance()
			}
			s := strings.TrimSpace(sb.String())
			if s != "" {
				opts = append(opts, &nodes.String{Str: s})
			}
			if p.cur.Type == ',' {
				p.advance()
			} else {
				break
			}
		}
	}

	if len(opts) == 0 {
		return nil
	}
	return &nodes.List{Items: opts}
}
