package parser

import (
	"strings"

	"github.com/bytebase/omni/doris/ast"
)

// parseCreateIndex parses a CREATE INDEX statement.
//
//	CREATE INDEX [IF NOT EXISTS] index_name ON table_name (col1, col2, ...)
//	  [USING index_type]
//	  [COMMENT 'comment']
//	  [PROPERTIES("key"="value", ...)]
//
// The CREATE keyword has already been consumed when this method is called.
// cur is kwINDEX on entry. startLoc is the Loc of the CREATE token.
func (p *Parser) parseCreateIndex(startLoc ast.Loc) (ast.Node, error) {
	// Consume INDEX.
	p.advance()

	// Optional IF NOT EXISTS.
	ifNotExists := false
	if p.cur.Kind == kwIF {
		p.advance()
		if _, err := p.expect(kwNOT); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		ifNotExists = true
	}

	// Index name.
	indexName, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}

	// ON table_name.
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	tableName, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}

	// (col1, col2, ...).
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	columns, err := p.parseIdentifierList()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}

	// Optional clauses: USING, COMMENT, PROPERTIES (in any order, each at most once).
	var indexType string
	var comment string
	var properties []*ast.Property

	for {
		switch p.cur.Kind {
		case kwUSING:
			p.advance()
			typeTok := p.cur
			// Index type is one of: BITMAP, NGRAM_BF, INVERTED, ANN (all non-reserved).
			// Accept any non-reserved keyword or plain identifier as an index type name.
			if !isIdentifierToken(typeTok.Kind) {
				return nil, p.syntaxErrorAtCur()
			}
			// Normalize to lowercase for case-insensitive comparison.
			indexType = strings.ToLower(typeTok.Str)
			p.advance()
			continue
		case kwCOMMENT:
			p.advance()
			if p.cur.Kind != tokString {
				return nil, p.syntaxErrorAtCur()
			}
			comment = p.cur.Str
			p.advance()
			continue
		case kwPROPERTIES:
			var err error
			properties, err = p.parseProperties()
			if err != nil {
				return nil, err
			}
			continue
		}
		break
	}

	endLoc := p.prev.Loc
	stmt := &ast.CreateIndexStmt{
		Name:        indexName,
		Table:       tableName,
		Columns:     columns,
		IfNotExists: ifNotExists,
		IndexType:   indexType,
		Properties:  properties,
		Comment:     comment,
		Loc:         ast.Loc{Start: startLoc.Start, End: endLoc.End},
	}
	return stmt, nil
}

// parseDropIndex parses a DROP INDEX statement.
//
//	DROP INDEX [IF EXISTS] index_name ON table_name
//
// The DROP keyword has already been consumed when this method is called.
// cur is kwINDEX on entry. startLoc is the Loc of the DROP token.
func (p *Parser) parseDropIndex(startLoc ast.Loc) (ast.Node, error) {
	// Consume INDEX.
	p.advance()

	// Optional IF EXISTS.
	ifExists := false
	if p.cur.Kind == kwIF {
		p.advance()
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		ifExists = true
	}

	// Index name.
	indexName, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}

	// ON table_name.
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	tableName, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}

	endLoc := p.prev.Loc
	stmt := &ast.DropIndexStmt{
		Name:     indexName,
		Table:    tableName,
		IfExists: ifExists,
		Loc:      ast.Loc{Start: startLoc.Start, End: endLoc.End},
	}
	return stmt, nil
}

// parseBuildIndex parses a BUILD INDEX statement.
//
//	BUILD INDEX index_name ON table_name [PARTITIONS(p1, p2, ...)]
//
// cur is kwINDEX on entry (BUILD has already been consumed).
// startLoc is the Loc of the BUILD token.
func (p *Parser) parseBuildIndex(startLoc ast.Loc) (ast.Node, error) {
	// Consume INDEX.
	p.advance()

	// Index name.
	indexName, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}

	// ON table_name.
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	tableName, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}

	// Optional PARTITIONS(p1, p2, ...).
	var partitions []string
	if p.cur.Kind == kwPARTITIONS {
		p.advance()
		if _, err := p.expect(int('(')); err != nil {
			return nil, err
		}
		partitions, err = p.parseIdentifierList()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(int(')')); err != nil {
			return nil, err
		}
	}

	endLoc := p.prev.Loc
	stmt := &ast.BuildIndexStmt{
		Name:       indexName,
		Table:      tableName,
		Partitions: partitions,
		Loc:        ast.Loc{Start: startLoc.Start, End: endLoc.End},
	}
	return stmt, nil
}

// parseIdentifierList parses a comma-separated list of bare identifiers.
// Used for column lists and partition name lists.
// The surrounding parentheses are NOT consumed here.
func (p *Parser) parseIdentifierList() ([]string, error) {
	first, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	names := []string{first}

	for p.cur.Kind == int(',') {
		p.advance() // consume ','
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, nil
}

// parseProperties parses a PROPERTIES("key"="value", ...) clause.
// cur must be kwPROPERTIES on entry; it is consumed here.
func (p *Parser) parseProperties() ([]*ast.Property, error) {
	p.advance() // consume PROPERTIES

	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}

	var props []*ast.Property
	for p.cur.Kind != int(')') && p.cur.Kind != tokEOF {
		propLoc := p.cur.Loc

		// Key: string literal or identifier.
		key, _, err := p.parseIdentifierOrString()
		if err != nil {
			return nil, err
		}

		if _, err := p.expect(int('=')); err != nil {
			return nil, err
		}

		// Value: string literal or identifier.
		value, valueLoc, err := p.parseIdentifierOrString()
		if err != nil {
			return nil, err
		}

		props = append(props, &ast.Property{
			Key:   key,
			Value: value,
			Loc:   ast.Loc{Start: propLoc.Start, End: valueLoc.End},
		})

		if p.cur.Kind == int(',') {
			p.advance()
		} else {
			break
		}
	}

	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	return props, nil
}
