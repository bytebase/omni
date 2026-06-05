package parser

import (
	"strings"

	"github.com/bytebase/omni/googlesql/ast"
)

// ---------------------------------------------------------------------------
// Spanner DDL — CHANGE STREAM (parser-ddl-spanner node)
// ---------------------------------------------------------------------------
//
// Cloud Spanner change streams (truth1 DDL-024/025/026):
//
//	CREATE CHANGE STREAM name { FOR ALL | FOR table_and_column [, ...] }?
//	  [OPTIONS ( option [, ...] )]
//	ALTER  CHANGE STREAM name { SET FOR { ALL | table_and_column [, ...] }
//	                          | DROP FOR ALL
//	                          | SET OPTIONS ( option [, ...] ) }
//	DROP   CHANGE STREAM name
//
//	table_and_column: table_name | table_name ( ) | table_name ( column [, ...] )
//
// DIALECT NOTE — these have NO first-class rule in the legacy ANTLR
// GoogleSQLParser.g4 (it models CHANGE STREAM via the generic-entity hook, which
// cannot express `CHANGE STREAM` as a two-word object-type with a FOR clause, so
// the legacy parser effectively over-rejects them). The omni parser models them
// directly. The authoritative oracle is the live Cloud Spanner emulator
// (oracle.md), which ACCEPTS every form here; verified in
// spanner_ddl_oracle_test.go. The `CHANGE` and `STREAM` words both lex as bare
// identifiers (non-reserved), so dispatch matches them by spelling.

// parseCreateChangeStream parses a CREATE CHANGE STREAM body. The shared CREATE
// prefix has been consumed by parseCreateStmt; cur is at the `CHANGE` identifier
// (peeked + confirmed `CHANGE STREAM` by the dispatcher). create_change_stream has
// no opt_or_replace / opt_create_scope, so a leading OR REPLACE / scope rejects.
func (p *Parser) parseCreateChangeStream(create Token, orReplace bool, scope string) (ast.Node, error) {
	if orReplace || scope != "" {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // CHANGE
	p.advance() // STREAM

	stmt := &ast.CreateChangeStreamStmt{}
	stmt.Loc.Start = create.Loc.Start

	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Optional FOR clause: FOR ALL | FOR table_and_column [, ...].
	if p.cur.Type == kwFOR {
		p.advance() // FOR
		stmt.HasFor = true
		if p.cur.Type == kwALL {
			p.advance() // ALL
			stmt.ForAll = true
		} else {
			tables, err := p.parseChangeStreamForTables()
			if err != nil {
				return nil, err
			}
			stmt.ForTables = tables
		}
	}

	// opt_options_list?
	if p.cur.Type == kwOPTIONS {
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseAlterChangeStream parses an ALTER CHANGE STREAM body. The ALTER keyword has
// been consumed by parseAlterStmt; cur is at `CHANGE` (dispatcher confirmed
// `CHANGE STREAM`).
func (p *Parser) parseAlterChangeStream(alter Token) (ast.Node, error) {
	p.advance() // CHANGE
	p.advance() // STREAM

	stmt := &ast.AlterChangeStreamStmt{}
	stmt.Loc.Start = alter.Loc.Start

	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	switch p.cur.Type {
	case kwSET:
		p.advance() // SET
		switch p.cur.Type {
		case kwFOR:
			p.advance() // FOR
			stmt.Action = ast.ChangeStreamSetFor
			if p.cur.Type == kwALL {
				p.advance() // ALL
				stmt.ForAll = true
			} else {
				tables, err := p.parseChangeStreamForTables()
				if err != nil {
					return nil, err
				}
				stmt.ForTables = tables
			}
		case kwOPTIONS:
			stmt.Action = ast.ChangeStreamSetOptions
			opts, err := p.parseOptionsList()
			if err != nil {
				return nil, err
			}
			stmt.Options = opts
		default:
			return nil, p.syntaxErrorAtCur()
		}
	case kwDROP:
		p.advance() // DROP
		if _, err := p.expect(kwFOR); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwALL); err != nil {
			return nil, err
		}
		stmt.Action = ast.ChangeStreamDropForAll
	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseDropChangeStream parses a DROP CHANGE STREAM body. The DROP keyword has
// been consumed by parseDropStmt; cur is at `CHANGE` (dispatcher confirmed
// `CHANGE STREAM`).
func (p *Parser) parseDropChangeStream(drop Token) (ast.Node, error) {
	p.advance() // CHANGE
	p.advance() // STREAM

	stmt := &ast.DropChangeStreamStmt{}
	stmt.Loc.Start = drop.Loc.Start

	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseChangeStreamForTables parses a non-empty `table_and_column [, ...]` list.
// cur is at the first table name. Each entry is:
//
//	table_name                     -> whole table (no parens)
//	table_name ( )                 -> ExplicitColumns, no columns
//	table_name ( column [, ...] )  -> ExplicitColumns + column list
func (p *Parser) parseChangeStreamForTables() ([]*ast.ChangeStreamTrackedTable, error) {
	var out []*ast.ChangeStreamTrackedTable
	for {
		entry, err := p.parseChangeStreamTrackedTable()
		if err != nil {
			return nil, err
		}
		out = append(out, entry)
		if _, ok := p.match(int(',')); ok {
			continue
		}
		break
	}
	return out, nil
}

// parseChangeStreamTrackedTable parses one table_and_column entry. cur is at the
// table name.
func (p *Parser) parseChangeStreamTrackedTable() (*ast.ChangeStreamTrackedTable, error) {
	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	entry := &ast.ChangeStreamTrackedTable{Name: name, Loc: name.Loc}

	// Optional `( )` or `( column [, ...] )`.
	if p.cur.Type == int('(') {
		p.advance() // '('
		entry.ExplicitColumns = true
		if p.cur.Type != int(')') {
			cols, err := p.parseChangeStreamColumnList()
			if err != nil {
				return nil, err
			}
			entry.Columns = cols
		}
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, err
		}
		entry.Loc.End = closeTok.Loc.End
	} else {
		entry.Loc.End = p.prev.Loc.End
	}
	return entry, nil
}

// parseChangeStreamColumnList parses a non-empty comma-separated column-name list
// inside a tracked table's `( … )`. cur is at the first column name.
func (p *Parser) parseChangeStreamColumnList() ([]string, error) {
	var cols []string
	for {
		tok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		name, err := p.identifierText(tok)
		if err != nil {
			return nil, err
		}
		cols = append(cols, name)
		if _, ok := p.match(int(',')); ok {
			continue
		}
		break
	}
	return cols, nil
}

// tokIsWord reports whether tok is an identifier or keyword whose spelling equals
// word (case-insensitive). It is the peek-friendly sibling of curIsWord, used by
// the CREATE/ALTER/DROP dispatchers to recognize the bare-identifier object words
// (CHANGE / STREAM / LOCALITY / GROUP) that the lexer does not tokenize as
// dedicated keywords.
func (p *Parser) tokIsWord(tok Token, word string) bool {
	if tok.Type == tokIdentifier {
		return strings.EqualFold(tok.Str, word)
	}
	if tok.Type >= keywordBase {
		return strings.EqualFold(TokenName(tok.Type), word)
	}
	return false
}
