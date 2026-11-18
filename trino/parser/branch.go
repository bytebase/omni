package parser

import "github.com/bytebase/omni/trino/ast"

// This file implements the table-branch statements of Trino 482's grammar:
// CREATE/DROP/ALTER BRANCH and SHOW BRANCHES (statement alternatives
// createBranch, dropBranch, fastForwardBranch and showBranches). They are
// statement-level nodes only; no catalog in this repository stores branches.

// BranchKind classifies a BranchStmt.
type BranchKind int

const (
	BranchCreate BranchKind = iota
	BranchDrop
	BranchFastForward
	BranchShow
)

// BranchStmt is one of the BRANCH statements. Name is the branch name (the
// target name for FastForward); Table is the table the branch belongs to.
type BranchStmt struct {
	Kind        BranchKind
	OrReplace   bool
	IfExists    bool
	IfNotExists bool
	Name        *ast.Identifier
	Target      *ast.Identifier
	Table       *ast.QualifiedName
	From        *ast.Identifier
	Properties  []*Property
	InKeyword   bool // SHOW BRANCHES IN TABLE (vs FROM)
	Loc         ast.Loc
}

func (n *BranchStmt) Tag() ast.NodeTag { return ast.T_Invalid }
func (n *BranchStmt) Span() ast.Loc    { return n.Loc }

var _ ast.Node = (*BranchStmt)(nil)

// parseCreateBranchStmt parses
// `CREATE [OR REPLACE] BRANCH [IF NOT EXISTS] name [WITH (props)] IN TABLE t
// [FROM source]`. CREATE is current.
func (p *Parser) parseCreateBranchStmt(start int) (ast.Node, error) {
	p.advance() // consume CREATE
	s := &BranchStmt{Kind: BranchCreate, Loc: ast.Loc{Start: start}}
	if p.cur.Kind == kwOR {
		p.advance()
		if _, err := p.expect(kwREPLACE); err != nil {
			return nil, err
		}
		s.OrReplace = true
	}
	if _, err := p.expect(kwBRANCH); err != nil {
		return nil, err
	}
	ifNotExists, err := p.parseOptionalIfNotExists()
	if err != nil {
		return nil, err
	}
	s.IfNotExists = ifNotExists
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	s.Name = name
	if p.cur.Kind == kwWITH {
		props, err := p.parseOptionalWithProperties()
		if err != nil {
			return nil, err
		}
		s.Properties = props
	}
	if _, err := p.expect(kwIN); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}
	table, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	s.Table = table
	s.Loc.End = table.Loc.End
	if p.cur.Kind == kwFROM {
		p.advance()
		from, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		s.From = from
		s.Loc.End = from.Loc.End
	}
	if s.OrReplace && s.IfNotExists {
		return nil, &ParseError{
			Position: s.Name.Loc.Start,
			End:      s.Name.Loc.End,
			Message:  "CREATE BRANCH cannot combine OR REPLACE with IF NOT EXISTS",
		}
	}
	return s, nil
}

// parseDropBranchStmt parses
// `DROP BRANCH [IF EXISTS] name IN TABLE t`. DROP is current.
func (p *Parser) parseDropBranchStmt(start int) (ast.Node, error) {
	p.advance() // consume DROP
	p.advance() // consume BRANCH
	ifExists, err := p.parseOptionalIfExists()
	if err != nil {
		return nil, err
	}
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwIN); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}
	table, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	return &BranchStmt{
		Kind:     BranchDrop,
		IfExists: ifExists,
		Name:     name,
		Table:    table,
		Loc:      ast.Loc{Start: start, End: table.Loc.End},
	}, nil
}

// parseAlterBranchStmt parses
// `ALTER BRANCH source IN TABLE t FAST FORWARD TO target`. ALTER is current.
func (p *Parser) parseAlterBranchStmt(start int) (ast.Node, error) {
	p.advance() // consume ALTER
	p.advance() // consume BRANCH
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwIN); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}
	table, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwFAST); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwFORWARD); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}
	target, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	return &BranchStmt{
		Kind:   BranchFastForward,
		Name:   name,
		Target: target,
		Table:  table,
		Loc:    ast.Loc{Start: start, End: target.Loc.End},
	}, nil
}

// parseShowBranches parses `SHOW BRANCHES (FROM|IN) TABLE t`. SHOW has already
// been consumed and BRANCHES is current.
func (p *Parser) parseShowBranches(start int) (ast.Node, error) {
	p.advance() // consume BRANCHES
	in := false
	switch p.cur.Kind {
	case kwIN:
		in = true
	case kwFROM:
	default:
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume FROM / IN
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}
	table, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	return &BranchStmt{
		Kind:      BranchShow,
		Table:     table,
		InKeyword: in,
		Loc:       ast.Loc{Start: start, End: table.Loc.End},
	}, nil
}

// parseAlterAuthorizationStmt parses
// `ALTER ownedEntityKind qualifiedName SET AUTHORIZATION principal`, where
// ownedEntityKind is TABLE, SCHEMA, VIEW, MATERIALIZED VIEW or any identifier.
// used by the generic ALTER dispatch for the kinds with no dedicated parser.
func (p *Parser) parseAlterAuthorizationStmt(start int) (ast.Node, error) {
	p.advance() // consume ALTER
	kind := ""
	if p.cur.Kind == kwMATERIALIZED && p.peekNext().Kind == kwVIEW {
		p.advance() // MATERIALIZED
		p.advance() // VIEW
		kind = "MATERIALIZED VIEW"
	} else {
		id, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		kind = id.String()
	}
	name, err := p.parseQualifiedName()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwSET); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwAUTHORIZATION); err != nil {
		return nil, err
	}
	princ, err := p.parseAuthorizationPrincipal()
	if err != nil {
		return nil, err
	}
	return &AlterTableStmt{
		Kind:          AlterTableSetAuthorization,
		Name:          name,
		Authorization: princ,
		EntityKind:    kind,
		Loc:           ast.Loc{Start: start, End: princ.Loc.End},
	}, nil
}
