package parser

import (
	"github.com/bytebase/omni/cassandra/ast"
)

// parseDescribe parses describeStatement (Parser.g:1410):
//
//	( DESCRIBE | DESC )
//	( CLUSTER | [FULL] SCHEMA | KEYSPACES | [ONLY] KEYSPACE [ks] | TABLES
//	| TABLE t | INDEX i | MATERIALIZED VIEW v | TYPES | TYPE t | FUNCTIONS
//	| FUNCTION f | AGGREGATES | AGGREGATE a | [ks.]name )
//	[ WITH INTERNALS ]
//
// CLUSTER, KEYSPACES, TABLES, TYPES, FUNCTIONS and AGGREGATES carry a
// syntactic predicate, so they are always the keyword form. TYPE, FUNCTION,
// AGGREGATE and ONLY are unreserved keywords: followed by no name they are
// the name described (DESCRIBE TYPE describes an object named type).
func (p *Parser) parseDescribe() (*ast.DescribeStmt, error) {
	start := p.curLoc()
	p.advance() // DESCRIBE | DESC

	stmt := &ast.DescribeStmt{}
	// takesName reports whether the token after a keyword starts the name
	// that keyword's form requires.
	takesName := func() bool {
		return isCQLName(p.peekNext())
	}
	// qualified parses the [ks.]name a keyword form takes.
	qualified := func(kind ast.DescribeKind) error {
		stmt.Kind = kind
		p.advance() // the keyword
		parse := p.parseCQLQualifiedName
		if kind == ast.DescribeType {
			parse = p.parseUserTypeQualifiedName
		}
		name, err := parse()
		if err != nil {
			return err
		}
		stmt.Name = name
		return nil
	}

	var err error
	switch {
	case p.curIs("cluster"):
		stmt.Kind = ast.DescribeCluster
		p.advance()
	case p.cur.Type == tokFULL:
		p.advance() // FULL
		if !p.curIs("schema") {
			return nil, p.errorf("expected SCHEMA after DESCRIBE FULL, got %s", p.tokenDesc())
		}
		p.advance()
		stmt.Kind = ast.DescribeSchema
		stmt.Full = true
	case p.curIs("schema"):
		stmt.Kind = ast.DescribeSchema
		p.advance()
	case p.cur.Type == tokKEYSPACES:
		stmt.Kind = ast.DescribeKeyspaces
		p.advance()
	case p.curIs("only") && p.peekNext().Type == tokKEYSPACE:
		p.advance() // ONLY
		stmt.Only = true
		err = p.parseDescribeKeyspace(stmt)
	case p.cur.Type == tokKEYSPACE:
		err = p.parseDescribeKeyspace(stmt)
	case p.curIs("tables"):
		stmt.Kind = ast.DescribeTables
		p.advance()
	case p.cur.Type == tokTABLE:
		err = qualified(ast.DescribeTable)
	case p.cur.Type == tokINDEX:
		err = qualified(ast.DescribeIndex)
	case p.cur.Type == tokMATERIALIZED:
		p.advance() // MATERIALIZED
		if p.cur.Type != tokVIEW {
			return nil, p.errorf("expected VIEW, got %s", p.tokenDesc())
		}
		err = qualified(ast.DescribeMaterializedView)
	case p.curIs("types"):
		stmt.Kind = ast.DescribeTypes
		p.advance()
	case p.cur.Type == tokTYPE && takesName():
		err = qualified(ast.DescribeType)
	case p.cur.Type == tokFUNCTIONS:
		stmt.Kind = ast.DescribeFunctions
		p.advance()
	case p.cur.Type == tokFUNCTION && takesName():
		err = qualified(ast.DescribeFunction)
	case p.curIs("aggregates"):
		stmt.Kind = ast.DescribeAggregates
		p.advance()
	case p.cur.Type == tokAGGREGATE && takesName():
		err = qualified(ast.DescribeAggregate)
	default:
		stmt.Kind = ast.DescribeName
		stmt.Name, err = p.parseCQLQualifiedName()
	}
	if err != nil {
		return nil, err
	}

	if p.cur.Type == tokWITH {
		p.advance() // WITH
		if !p.curIs("internals") {
			return nil, p.errorf("expected INTERNALS, got %s", p.tokenDesc())
		}
		p.advance()
		stmt.WithInternals = true
	}

	stmt.Loc = p.makeLoc(start)
	return stmt, nil
}

// parseDescribeKeyspace parses KEYSPACE [ks] of DESCRIBE [ONLY] KEYSPACE.
func (p *Parser) parseDescribeKeyspace(stmt *ast.DescribeStmt) error {
	stmt.Kind = ast.DescribeKeyspace
	p.advance() // KEYSPACE
	if !isCQLName(p.cur) {
		return nil
	}
	ks, err := p.parseNameIf(isCQLName)
	if err != nil {
		return err
	}
	stmt.Name = &ast.QualifiedName{Parts: []*ast.Identifier{ks}, Loc: ks.Loc}
	return nil
}
