// Package parser - create_table.go implements T-SQL CREATE TABLE statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/tsql/ast"
)

// parseCreateTableStmt parses a CREATE TABLE statement.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-table-transact-sql
//
//	CREATE TABLE name ( col_def | constraint [, ...] )
func (p *Parser) parseCreateTableStmt() *nodes.CreateTableStmt {
	loc := p.pos()

	stmt := &nodes.CreateTableStmt{
		Loc: nodes.Loc{Start: loc},
	}

	// Table name
	stmt.Name = p.parseTableRef()

	// Column and constraint definitions
	if _, err := p.expect('('); err != nil {
		stmt.Loc.End = p.pos()
		return stmt
	}

	var cols []nodes.Node
	var constraints []nodes.Node

	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		// Check if this is a table-level constraint
		if p.cur.Type == kwCONSTRAINT || p.cur.Type == kwPRIMARY ||
			p.cur.Type == kwUNIQUE || p.cur.Type == kwCHECK ||
			p.cur.Type == kwFOREIGN {
			constraint := p.parseTableConstraint()
			if constraint != nil {
				constraints = append(constraints, constraint)
			}
		} else {
			col := p.parseColumnDef()
			if col != nil {
				cols = append(cols, col)
			}
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}
	_, _ = p.expect(')')

	if len(cols) > 0 {
		stmt.Columns = &nodes.List{Items: cols}
	}
	if len(constraints) > 0 {
		stmt.Constraints = &nodes.List{Items: constraints}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseColumnDef parses a column definition.
//
//	col_def = name type [IDENTITY(seed,incr)] [NULL|NOT NULL] [DEFAULT expr]
//	          [CONSTRAINT name ...] [COLLATE name]
func (p *Parser) parseColumnDef() *nodes.ColumnDef {
	loc := p.pos()

	name, ok := p.parseIdentifier()
	if !ok {
		return nil
	}

	col := &nodes.ColumnDef{
		Name: name,
		Loc:  nodes.Loc{Start: loc},
	}

	// Check for computed column: name AS expr
	if p.cur.Type == kwAS {
		p.advance()
		compLoc := p.pos()
		expr := p.parseExpr()
		persisted := false
		if p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "persisted") {
			persisted = true
			p.advance()
		}
		col.Computed = &nodes.ComputedColumnDef{
			Expr:      expr,
			Persisted: persisted,
			Loc:       nodes.Loc{Start: compLoc},
		}
		col.Loc.End = p.pos()
		return col
	}

	// Data type
	col.DataType = p.parseDataType()

	// Column-level options in any order
	for {
		consumed := false

		// IDENTITY
		if p.cur.Type == kwIDENTITY {
			col.Identity = p.parseIdentitySpec()
			consumed = true
		}

		// NULL / NOT NULL
		if p.cur.Type == kwNULL {
			p.advance()
			col.Nullable = &nodes.NullableSpec{NotNull: false, Loc: nodes.Loc{Start: p.pos()}}
			consumed = true
		} else if p.cur.Type == kwNOT {
			next := p.peekNext()
			if next.Type == kwNULL {
				p.advance() // NOT
				p.advance() // NULL
				col.Nullable = &nodes.NullableSpec{NotNull: true, Loc: nodes.Loc{Start: p.pos()}}
				consumed = true
			}
		}

		// DEFAULT
		if p.cur.Type == kwDEFAULT {
			p.advance()
			col.DefaultExpr = p.parseExpr()
			consumed = true
		}

		// COLLATE
		if p.cur.Type == kwCOLLATE {
			p.advance()
			if p.isIdentLike() {
				col.Collation = p.cur.Str
				p.advance()
			}
			consumed = true
		}

		// CONSTRAINT (inline column constraint)
		if p.cur.Type == kwCONSTRAINT {
			constraint := p.parseColumnConstraint()
			if constraint != nil {
				if col.Constraints == nil {
					col.Constraints = &nodes.List{}
				}
				col.Constraints.Items = append(col.Constraints.Items, constraint)
			}
			consumed = true
		}

		// PRIMARY KEY / UNIQUE (without CONSTRAINT keyword)
		if p.cur.Type == kwPRIMARY || p.cur.Type == kwUNIQUE {
			constraint := p.parseInlineConstraint("")
			if constraint != nil {
				if col.Constraints == nil {
					col.Constraints = &nodes.List{}
				}
				col.Constraints.Items = append(col.Constraints.Items, constraint)
			}
			consumed = true
		}

		// CHECK (without CONSTRAINT keyword)
		if p.cur.Type == kwCHECK {
			constraint := p.parseInlineConstraint("")
			if constraint != nil {
				if col.Constraints == nil {
					col.Constraints = &nodes.List{}
				}
				col.Constraints.Items = append(col.Constraints.Items, constraint)
			}
			consumed = true
		}

		// REFERENCES (inline FK without CONSTRAINT keyword)
		if p.cur.Type == kwREFERENCES {
			constraint := p.parseInlineConstraint("")
			if constraint != nil {
				if col.Constraints == nil {
					col.Constraints = &nodes.List{}
				}
				col.Constraints.Items = append(col.Constraints.Items, constraint)
			}
			consumed = true
		}

		if !consumed {
			break
		}
	}

	col.Loc.End = p.pos()
	return col
}

// parseIdentitySpec parses IDENTITY(seed, increment).
func (p *Parser) parseIdentitySpec() *nodes.IdentitySpec {
	loc := p.pos()
	p.advance() // consume IDENTITY

	spec := &nodes.IdentitySpec{
		Seed:      1,
		Increment: 1,
		Loc:       nodes.Loc{Start: loc},
	}

	if p.cur.Type == '(' {
		p.advance()
		if p.cur.Type == tokICONST {
			spec.Seed = p.cur.Ival
			p.advance()
		}
		if _, ok := p.match(','); ok {
			if p.cur.Type == tokICONST {
				spec.Increment = p.cur.Ival
				p.advance()
			}
		}
		_, _ = p.expect(')')
	}

	spec.Loc.End = p.pos()
	return spec
}

// parseColumnConstraint parses CONSTRAINT name followed by constraint type.
func (p *Parser) parseColumnConstraint() *nodes.ConstraintDef {
	p.advance() // consume CONSTRAINT
	name, _ := p.parseIdentifier()
	return p.parseInlineConstraint(name)
}

// parseInlineConstraint parses a constraint type (PRIMARY KEY, UNIQUE, CHECK, DEFAULT, REFERENCES).
func (p *Parser) parseInlineConstraint(name string) *nodes.ConstraintDef {
	loc := p.pos()
	cd := &nodes.ConstraintDef{
		Name: name,
		Loc:  nodes.Loc{Start: loc},
	}

	switch p.cur.Type {
	case kwPRIMARY:
		p.advance() // PRIMARY
		p.match(kwKEY)
		cd.Type = nodes.ConstraintPrimaryKey
		p.parseClusteredOption(cd)
	case kwUNIQUE:
		p.advance()
		cd.Type = nodes.ConstraintUnique
		p.parseClusteredOption(cd)
	case kwCHECK:
		p.advance()
		cd.Type = nodes.ConstraintCheck
		if _, err := p.expect('('); err == nil {
			cd.Expr = p.parseExpr()
			_, _ = p.expect(')')
		}
	case kwDEFAULT:
		p.advance()
		cd.Type = nodes.ConstraintDefault
		cd.Expr = p.parseExpr()
	case kwREFERENCES:
		p.advance()
		cd.Type = nodes.ConstraintForeignKey
		cd.RefTable = p.parseTableRef()
		if p.cur.Type == '(' {
			cd.RefColumns = p.parseParenIdentList()
		}
		p.parseReferentialActions(cd)
	default:
		return nil
	}

	cd.Loc.End = p.pos()
	return cd
}

// parseTableConstraint parses a table-level constraint.
func (p *Parser) parseTableConstraint() *nodes.ConstraintDef {
	loc := p.pos()
	var name string

	if p.cur.Type == kwCONSTRAINT {
		p.advance()
		name, _ = p.parseIdentifier()
	}

	cd := &nodes.ConstraintDef{
		Name: name,
		Loc:  nodes.Loc{Start: loc},
	}

	switch p.cur.Type {
	case kwPRIMARY:
		p.advance() // PRIMARY
		p.match(kwKEY)
		cd.Type = nodes.ConstraintPrimaryKey
		p.parseClusteredOption(cd)
		if p.cur.Type == '(' {
			cd.Columns = p.parseParenIdentList()
		}
	case kwUNIQUE:
		p.advance()
		cd.Type = nodes.ConstraintUnique
		p.parseClusteredOption(cd)
		if p.cur.Type == '(' {
			cd.Columns = p.parseParenIdentList()
		}
	case kwCHECK:
		p.advance()
		cd.Type = nodes.ConstraintCheck
		if _, err := p.expect('('); err == nil {
			cd.Expr = p.parseExpr()
			_, _ = p.expect(')')
		}
	case kwFOREIGN:
		p.advance() // FOREIGN
		p.match(kwKEY)
		cd.Type = nodes.ConstraintForeignKey
		if p.cur.Type == '(' {
			cd.Columns = p.parseParenIdentList()
		}
		if _, ok := p.match(kwREFERENCES); ok {
			cd.RefTable = p.parseTableRef()
			if p.cur.Type == '(' {
				cd.RefColumns = p.parseParenIdentList()
			}
			p.parseReferentialActions(cd)
		}
	case kwDEFAULT:
		p.advance()
		cd.Type = nodes.ConstraintDefault
		cd.Expr = p.parseExpr()
		// FOR column
		if _, ok := p.match(kwFOR); ok {
			p.parseIdentifier() // column name (not stored separately but consumed)
		}
	default:
		return nil
	}

	cd.Loc.End = p.pos()
	return cd
}

// parseClusteredOption parses optional CLUSTERED/NONCLUSTERED.
func (p *Parser) parseClusteredOption(cd *nodes.ConstraintDef) {
	if p.cur.Type == kwCLUSTERED {
		p.advance()
		v := true
		cd.Clustered = &v
	} else if p.cur.Type == kwNONCLUSTERED {
		p.advance()
		v := false
		cd.Clustered = &v
	}
}

// parseReferentialActions parses ON DELETE/UPDATE actions.
func (p *Parser) parseReferentialActions(cd *nodes.ConstraintDef) {
	for p.cur.Type == kwON {
		p.advance()
		if p.cur.Type == kwDELETE {
			p.advance()
			cd.OnDelete = p.parseRefAction()
		} else if p.cur.Type == kwUPDATE {
			p.advance()
			cd.OnUpdate = p.parseRefAction()
		} else {
			break
		}
	}
}

// parseRefAction parses a referential action (CASCADE, SET NULL, SET DEFAULT, NO ACTION).
func (p *Parser) parseRefAction() nodes.ReferentialAction {
	if _, ok := p.match(kwCASCADE); ok {
		return nodes.RefActCascade
	}
	if p.cur.Type == kwSET {
		p.advance()
		if _, ok := p.match(kwNULL); ok {
			return nodes.RefActSetNull
		}
		if _, ok := p.match(kwDEFAULT); ok {
			return nodes.RefActSetDefault
		}
	}
	// NO ACTION
	if p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "no") {
		p.advance()
		if p.cur.Type == tokIDENT && strings.EqualFold(p.cur.Str, "action") {
			p.advance()
		}
		return nodes.RefActNoAction
	}
	return nodes.RefActNone
}

// parseParenIdentList parses (ident, ident, ...).
func (p *Parser) parseParenIdentList() *nodes.List {
	p.advance() // consume (
	var items []nodes.Node
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		name, ok := p.parseIdentifier()
		if !ok {
			break
		}
		items = append(items, &nodes.String{Str: name})
		if _, ok := p.match(','); !ok {
			break
		}
	}
	_, _ = p.expect(')')
	return &nodes.List{Items: items}
}
