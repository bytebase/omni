package parser

import (
	"strings"

	"github.com/bytebase/omni/mongo/ast"
)

// parseShowCommand parses "show dbs", "show databases", "show collections",
// "show log [name]", "show automationNotices", "show nonGenuineMongoDBCheck", etc.
func (p *Parser) parseShowCommand() (ast.Node, error) {
	showTok := p.advance() // consume 'show'
	start := showTok.Loc

	var target string
	var arg string
	switch p.cur.Type {
	case kwDbs:
		target = "databases"
		p.advance()
	case kwDatabases:
		target = "databases"
		p.advance()
	case kwCollections:
		target = "collections"
		p.advance()
	case kwTables:
		target = "tables"
		p.advance()
	case kwProfile:
		target = "profile"
		p.advance()
	case kwUsers:
		target = "users"
		p.advance()
	case kwRoles:
		target = "roles"
		p.advance()
	case kwLog:
		target = "log"
		p.advance()
		// Optional log name argument on the same line (e.g. "show log global").
		if !p.hasNewlineBeforeCur() && p.isIdentLike(p.cur.Type) {
			arg = p.advance().Str
		}
	case kwLogs:
		target = "logs"
		p.advance()
	case kwStartupWarnings:
		target = "startupWarnings"
		p.advance()
	case kwAutomationNotices:
		target = "automationNotices"
		p.advance()
	case kwNonGenuineMongoDBCheck:
		target = "nonGenuineMongoDBCheck"
		p.advance()
	default:
		return nil, p.errorAtCur("expected show target (dbs, databases, collections, etc.)")
	}

	return &ast.ShowCommand{
		Target: target,
		Arg:    arg,
		Loc:    ast.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseUseCommand parses the direct shell command "use <database>".
func (p *Parser) parseUseCommand() (ast.Node, error) {
	useTok := p.advance() // consume 'use'
	start := useTok.Loc

	if p.hasNewlineBeforeCur() || p.cur.Type == tokEOF || p.cur.Type == ';' {
		return nil, p.errorAtCur("expected database name after 'use'")
	}

	var dbName string
	var nameStart, nameEnd int
	switch {
	case p.cur.Type == tokString:
		tok := p.advance()
		dbName = tok.Str
		nameStart, nameEnd = tok.Loc, tok.End
	case p.isIdentLike(p.cur.Type):
		tok := p.advance()
		dbName = tok.Str
		nameStart, nameEnd = tok.Loc, tok.End
		// Allow hyphenated database names without whitespace (e.g. "use my-db", "use my-1").
		for p.cur.Loc == nameEnd {
			if p.cur.Type == '-' {
				p.advance() // consume '-'
				if p.cur.Loc != p.prev.End || (!p.isIdentLike(p.cur.Type) && p.cur.Type != tokNumber) {
					return nil, p.errorAtCur("expected identifier after '-' in database name")
				}
				part := p.advance()
				dbName += "-" + part.Str
				nameEnd = part.End
				continue
			}
			// Lexer scans "-1" as a single negative tokNumber when '-' precedes digits.
			if p.cur.Type == tokNumber && strings.HasPrefix(p.cur.Str, "-") {
				part := p.advance()
				dbName += part.Str
				nameEnd = part.End
				continue
			}
			break
		}
	default:
		return nil, p.errorAtCur("expected database name after 'use'")
	}

	return &ast.NativeFunctionCall{
		Name: "use",
		Args: []ast.Node{
			&ast.StringLiteral{
				Value: dbName,
				Loc:   ast.Loc{Start: nameStart, End: nameEnd},
			},
		},
		Loc: ast.Loc{Start: start, End: p.prev.End},
	}, nil
}

// parseDirectShellCommand parses parentheseless @directShellCommand commands:
// exit [code], quit [code], it, cls.
func (p *Parser) parseDirectShellCommand() (ast.Node, error) {
	tok := p.advance()
	start := tok.Loc

	var args []ast.Node
	if (tok.Type == kwExit || tok.Type == kwQuit) && !p.hasNewlineBeforeCur() && p.cur.Type == tokNumber {
		numTok := p.advance()
		args = append(args, &ast.NumberLiteral{
			Value:   numTok.Str,
			IsFloat: strings.ContainsAny(numTok.Str, ".eE"),
			Loc:     ast.Loc{Start: numTok.Loc, End: numTok.End},
		})
	}

	return &ast.NativeFunctionCall{
		Name: tok.Str,
		Args: args,
		Loc:  ast.Loc{Start: start, End: p.prev.End},
	}, nil
}
