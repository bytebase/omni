// Package parser - partition.go implements T-SQL CREATE/ALTER/DROP PARTITION FUNCTION/SCHEME parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseCreatePartitionFunctionStmt parses a CREATE PARTITION FUNCTION statement.
//
// BNF: mssql/parser/bnf/create-partition-function-transact-sql.bnf
//
//	CREATE PARTITION FUNCTION partition_function_name ( input_parameter_type )
//	AS RANGE [ LEFT | RIGHT ]
//	FOR VALUES ( [ boundary_value [ ,...n ] ] )
func (p *Parser) parseCreatePartitionFunctionStmt() (*nodes.CreatePartitionFunctionStmt, error) {
	loc := p.pos()
	// FUNCTION keyword already consumed by caller

	stmt := &nodes.CreatePartitionFunctionStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Name
	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// ( input_parameter_type )
	if p.cur.Type == '(' {
		p.advance()
		stmt.InputType , _ = p.parseDataType()
		p.match(')')
	}

	// AS RANGE [LEFT|RIGHT]
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "AS") {
		p.advance()
	}
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "RANGE") {
		p.advance()
	}
	if p.isIdentLike() || p.cur.Type == kwLEFT || p.cur.Type == kwRIGHT {
		if matchesKeywordCI(p.cur.Str, "LEFT") || matchesKeywordCI(p.cur.Str, "RIGHT") {
			stmt.Range = strings.ToUpper(p.cur.Str)
			p.advance()
		}
	}

	// FOR VALUES ( boundary_value [,...n] )
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "FOR") {
		p.advance()
	}
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "VALUES") {
		p.advance()
	}
	if p.cur.Type == '(' {
		p.advance()
		var vals []nodes.Node
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			expr, _ := p.parseExpr()
			if expr != nil {
				vals = append(vals, expr)
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		p.match(')')
		stmt.Values = &nodes.List{Items: vals}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterPartitionFunctionStmt parses an ALTER PARTITION FUNCTION statement.
//
// BNF: mssql/parser/bnf/alter-partition-function-transact-sql.bnf
//
//	ALTER PARTITION FUNCTION partition_function_name ()
//	{ SPLIT | MERGE } RANGE ( boundary_value )
func (p *Parser) parseAlterPartitionFunctionStmt() (*nodes.AlterPartitionFunctionStmt, error) {
	loc := p.pos()
	// FUNCTION keyword already consumed by caller

	stmt := &nodes.AlterPartitionFunctionStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Name
	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// ()
	if p.cur.Type == '(' {
		p.advance()
		p.match(')')
	}

	// SPLIT or MERGE
	if p.isIdentLike() {
		stmt.Action = strings.ToUpper(p.cur.Str)
		p.advance()
	}

	// RANGE ( boundary_value )
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "RANGE") {
		p.advance()
	}
	if p.cur.Type == '(' {
		p.advance()
		stmt.BoundaryValue, _ = p.parseExpr()
		p.match(')')
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCreatePartitionSchemeStmt parses a CREATE PARTITION SCHEME statement.
//
// BNF: mssql/parser/bnf/create-partition-scheme-transact-sql.bnf
//
//	CREATE PARTITION SCHEME partition_scheme_name
//	AS PARTITION partition_function_name
//	[ ALL ] TO ( { file_group_name | [ PRIMARY ] } [ ,...n ] )
func (p *Parser) parseCreatePartitionSchemeStmt() (*nodes.CreatePartitionSchemeStmt, error) {
	loc := p.pos()
	// SCHEME keyword already consumed by caller

	stmt := &nodes.CreatePartitionSchemeStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Name
	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// AS PARTITION function_name
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "AS") {
		p.advance()
	}
	if p.cur.Type == kwPARTITION {
		p.advance()
	}
	if p.isIdentLike() {
		stmt.FunctionName = p.cur.Str
		p.advance()
	}

	// ALL TO filegroup | TO ( filegrouplist )
	allTo := false
	if p.cur.Type == kwALL {
		p.advance()
		allTo = true
	}
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "TO") {
		p.advance()
	}

	if p.cur.Type == '(' {
		p.advance()
		var fgs []nodes.Node
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			if p.isIdentLike() || p.cur.Type == kwPRIMARY {
				fgname := p.cur.Str
				p.advance()
				if allTo && len(fgs) == 0 {
					stmt.AllToFileGroup = strings.ToUpper(fgname)
				}
				fgs = append(fgs, &nodes.String{Str: strings.ToUpper(fgname)})
			}
			if _, ok := p.match(','); !ok {
				break
			}
		}
		p.match(')')
		if !allTo {
			stmt.FileGroups = &nodes.List{Items: fgs}
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseAlterPartitionSchemeStmt parses an ALTER PARTITION SCHEME statement.
//
// BNF: mssql/parser/bnf/alter-partition-scheme-transact-sql.bnf
//
//	ALTER PARTITION SCHEME partition_scheme_name
//	NEXT USED [ filegroup_name ]
func (p *Parser) parseAlterPartitionSchemeStmt() (*nodes.AlterPartitionSchemeStmt, error) {
	loc := p.pos()
	// SCHEME keyword already consumed by caller

	stmt := &nodes.AlterPartitionSchemeStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Name
	if p.isIdentLike() {
		stmt.Name = p.cur.Str
		p.advance()
	}

	// NEXT USED [filegroup_name]
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "NEXT") {
		p.advance()
	}
	if p.isIdentLike() && matchesKeywordCI(p.cur.Str, "USED") {
		p.advance()
	}
	if p.isIdentLike() {
		stmt.FileGroup = p.cur.Str
		p.advance()
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}
