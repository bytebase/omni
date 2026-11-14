// Package parser - partition.go implements T-SQL CREATE/ALTER/DROP PARTITION FUNCTION/SCHEME parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// alterPartitionFunctionActions matches SqlScriptDOM's closed set for
// ALTER PARTITION FUNCTION: only SPLIT or MERGE after the empty `()`.
var alterPartitionFunctionActions = newOptionSet().withIdents("SPLIT", "MERGE")

// parseCreatePartitionFunctionStmt parses a CREATE PARTITION FUNCTION statement.
//
// BNF: mssql/parser/bnf/create-partition-function-transact-sql.bnf
//
//	CREATE PARTITION FUNCTION partition_function_name ( input_parameter_type )
//	AS RANGE [ LEFT | RIGHT ]
//	FOR VALUES ( [ boundary_value [ ,...n ] ] )
//
// The grammar is createPartitionFunction and partitionFunctionTypes
// (TSql160.g:10046, TSql160.g:10075): every part is required, the type may
// carry a collation, and a boundary value is an expression that names no
// column (Msg 128).
func (p *Parser) parseCreatePartitionFunctionStmt() (*nodes.CreatePartitionFunctionStmt, error) {
	loc := p.pos()
	// FUNCTION keyword already consumed by caller

	stmt := &nodes.CreatePartitionFunctionStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Name
	name, ok := p.parseIdentifier()
	if !ok {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Name = name

	// ( input_parameter_type [ COLLATE collation_name ] )
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	dt, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	if dt == nil {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.InputType = dt
	if p.cur.Type == kwCOLLATE {
		p.advance()
		coll, ok := p.parseIdentifier()
		if !ok {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Collation = coll
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	// AS RANGE [LEFT|RIGHT]
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwRANGE); err != nil {
		return nil, err
	}
	if p.cur.Type == kwLEFT || p.cur.Type == kwRIGHT {
		stmt.Range = strings.ToUpper(p.cur.Str)
		p.advance()
	}

	// FOR VALUES ( [ boundary_value [,...n] ] )
	if _, err := p.expect(kwFOR); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwVALUES); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var vals []nodes.Node
	if p.cur.Type != ')' {
		for {
			expr, err := p.partitionParseBoundary()
			if err != nil {
				return nil, err
			}
			vals = append(vals, expr)
			if _, ok := p.match(','); !ok {
				break
			}
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	stmt.Values = &nodes.List{Items: vals}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// partitionParseBoundary parses a boundary value of a partition function: an
// expression, which SQL Server refuses when it is a bare column name
// (Msg 128 "The name is not permitted in this context").
func (p *Parser) partitionParseBoundary() (nodes.ExprNode, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.syntaxErrorAtCur()
	}
	if _, ok := expr.(*nodes.ColumnRef); ok {
		return nil, p.newParseError(expr.(*nodes.ColumnRef).Loc.Start, "a column name is not permitted in a boundary value")
	}
	return expr, nil
}

// parseAlterPartitionFunctionStmt parses an ALTER PARTITION FUNCTION statement.
//
// BNF: mssql/parser/bnf/alter-partition-function-transact-sql.bnf
//
//	ALTER PARTITION FUNCTION partition_function_name ()
//	{ SPLIT | MERGE } RANGE ( boundary_value )
//
// The grammar is alterPartitionFunctionStatement (TSql160.g:12230): the
// empty parentheses are required, and RANGE ( boundary_value ) is optional.
func (p *Parser) parseAlterPartitionFunctionStmt() (*nodes.AlterPartitionFunctionStmt, error) {
	loc := p.pos()
	// FUNCTION keyword already consumed by caller

	stmt := &nodes.AlterPartitionFunctionStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Name
	name, ok := p.parseIdentifier()
	if !ok {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Name = name

	// ()
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	// SPLIT or MERGE — SqlScriptDOM's AlterPartitionFunctionStatement closed set.
	if !p.isValidOption(alterPartitionFunctionActions) {
		return nil, p.unexpectedToken()
	}
	stmt.Action = strings.ToUpper(p.cur.Str)
	p.advance()

	// [ RANGE ( boundary_value ) ]
	if p.cur.Type == kwRANGE {
		p.advance()
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		v, err := p.partitionParseBoundary()
		if err != nil {
			return nil, err
		}
		stmt.BoundaryValue = v
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
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
//
// The grammar is createPartitionScheme (TSql160.g:10153): a filegroup is a
// string or an identifier (stringOrIdentifier, TSql160.g:33229), so PRIMARY
// must be bracketed.
func (p *Parser) parseCreatePartitionSchemeStmt() (*nodes.CreatePartitionSchemeStmt, error) {
	loc := p.pos()
	// SCHEME keyword already consumed by caller

	stmt := &nodes.CreatePartitionSchemeStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Name
	name, ok := p.parseIdentifier()
	if !ok {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Name = name

	// AS PARTITION function_name
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwPARTITION); err != nil {
		return nil, err
	}
	fn, ok := p.parseIdentifier()
	if !ok {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.FunctionName = fn

	// [ ALL ] TO ( filegroup [ ,...n ] )
	allTo := false
	if p.cur.Type == kwALL {
		p.advance()
		allTo = true
	}
	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var fgs []nodes.Node
	for {
		fgname, err := p.partitionParseFileGroup()
		if err != nil {
			return nil, err
		}
		if allTo && len(fgs) == 0 {
			stmt.AllToFileGroup = strings.ToUpper(fgname)
		}
		fgs = append(fgs, &nodes.String{Str: strings.ToUpper(fgname)})
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	if !allTo {
		stmt.FileGroups = &nodes.List{Items: fgs}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// partitionParseFileGroup parses a filegroup name: a string or an
// identifier (stringOrIdentifier, TSql160.g:33229).
func (p *Parser) partitionParseFileGroup() (string, error) {
	if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
		s := p.cur.Str
		p.advance()
		return s, nil
	}
	name, ok := p.parseIdentifier()
	if !ok {
		return "", p.syntaxErrorAtCur()
	}
	return name, nil
}

// parseAlterPartitionSchemeStmt parses an ALTER PARTITION SCHEME statement.
//
// BNF: mssql/parser/bnf/alter-partition-scheme-transact-sql.bnf
//
//	ALTER PARTITION SCHEME partition_scheme_name
//	NEXT USED [ filegroup_name ]
//
// The grammar is alterPartitionSchemeStatement (TSql160.g:12263).
func (p *Parser) parseAlterPartitionSchemeStmt() (*nodes.AlterPartitionSchemeStmt, error) {
	loc := p.pos()
	// SCHEME keyword already consumed by caller

	stmt := &nodes.AlterPartitionSchemeStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Name
	name, ok := p.parseIdentifier()
	if !ok {
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Name = name

	// NEXT USED [filegroup_name]
	if _, err := p.expect(kwNEXT); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwUSED); err != nil {
		return nil, err
	}
	if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST || p.isIdentLike() {
		fg, err := p.partitionParseFileGroup()
		if err != nil {
			return nil, err
		}
		stmt.FileGroup = fg
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}
