package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/redshift/ast"
)

func appendNodeLists(left, right *nodes.List) *nodes.List {
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}
	left.Items = append(left.Items, right.Items...)
	return left
}

func (p *Parser) redshiftWordEqual(word string) bool {
	return strings.EqualFold(p.cur.Str, word)
}

func (p *Parser) consumeRedshiftWord(word string) bool {
	if !p.redshiftWordEqual(word) {
		return false
	}
	p.advance()
	return true
}

// parseRedshiftOptionWord parses the keyword a table attribute takes (EVEN,
// ALL, KEY, YES, NO). A quoted name is no keyword, so only an unquoted word
// parses here.
func (p *Parser) parseRedshiftOptionWord() (string, error) {
	if p.isColLabel() && !p.isDelimitedIdentToken(p.cur) {
		tok := p.advance()
		return tok.Str, nil
	}
	return "", p.syntaxErrorAtCur()
}

func (p *Parser) parseRedshiftCreateTableAttributes() (*nodes.List, error) {
	var items []nodes.Node
	for {
		switch {
		case p.acceptRedshiftKeyword("diststyle"):
			style, err := p.parseRedshiftOptionWord()
			if err != nil {
				return nil, err
			}
			items = append(items, makeDefElem("diststyle", &nodes.String{Str: strings.ToLower(style)}))
		case p.acceptRedshiftKeyword("distkey"):
			var cols nodes.Node
			if p.cur.Type == '(' {
				p.advance()
				list, err := p.parseColumnList()
				if err != nil {
					return nil, err
				}
				if list == nil {
					return nil, p.syntaxErrorAtCur()
				}
				cols = list
				if _, err := p.expect(')'); err != nil {
					return nil, err
				}
			}
			items = append(items, makeDefElem("distkey", cols))
		case p.redshiftWordEqual("compound") || p.redshiftWordEqual("interleaved") || p.redshiftWordEqual("sortkey"):
			sortStyle := ""
			if p.redshiftWordEqual("compound") || p.redshiftWordEqual("interleaved") {
				sortStyle = strings.ToLower(p.advance().Str)
			}
			if !p.acceptRedshiftKeyword("sortkey") {
				return nil, p.syntaxErrorAtCur()
			}
			if p.acceptRedshiftKeyword("auto") {
				if sortStyle != "" {
					items = append(items, makeDefElem("sortstyle", &nodes.String{Str: sortStyle}))
				}
				items = append(items, makeDefElem("sortkey", &nodes.String{Str: "auto"}))
				continue
			}
			if _, err := p.expect('('); err != nil {
				return nil, err
			}
			cols, err := p.parseColumnList()
			if err != nil {
				return nil, err
			}
			if cols == nil {
				return nil, p.syntaxErrorAtCur()
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			if sortStyle != "" {
				items = append(items, makeDefElem("sortstyle", &nodes.String{Str: sortStyle}))
			}
			items = append(items, makeDefElem("sortkey", cols))
		case p.acceptRedshiftKeyword("encode"):
			encoding := "auto"
			if p.isColLabel() {
				encoding = strings.ToLower(p.advance().Str)
			}
			items = append(items, makeDefElem("encode", &nodes.String{Str: encoding}))
		case p.acceptRedshiftKeyword("backup"):
			backup, err := p.parseRedshiftOptionWord()
			if err != nil {
				return nil, err
			}
			items = append(items, makeDefElem("backup", &nodes.String{Str: strings.ToLower(backup)}))
		default:
			if len(items) == 0 {
				return nil, nil
			}
			return &nodes.List{Items: items}, nil
		}
	}
}

func (p *Parser) parseRedshiftCreateTableAttributesInto(options *nodes.List) (*nodes.List, error) {
	attrs, err := p.parseRedshiftCreateTableAttributes()
	if err != nil {
		return nil, err
	}
	return appendNodeLists(options, attrs), nil
}

func (p *Parser) parseRedshiftColumnAttribute() (nodes.Node, error) {
	switch {
	case p.acceptRedshiftKeyword("encode"):
		encoding, err := p.parseRedshiftOptionWord()
		if err != nil {
			return nil, err
		}
		return makeDefElem("encode", &nodes.String{Str: strings.ToLower(encoding)}), nil
	case p.acceptRedshiftKeyword("distkey"):
		return makeDefElem("distkey", nil), nil
	case p.acceptRedshiftKeyword("sortkey"):
		return makeDefElem("sortkey", nil), nil
	case p.cur.Type == IDENTITY_P:
		return p.parseRedshiftIdentityColumnConstraint()
	// optional-probe: Redshift column attributes are optional in the column constraint loop.
	default:
		return nil, nil
	}
}

func (p *Parser) parseRedshiftIdentityColumnConstraint() (*nodes.Constraint, error) {
	p.advance() // IDENTITY

	var opts *nodes.List
	if p.cur.Type == '(' {
		var err error
		opts, err = p.parseRedshiftIdentitySeedStepOptions()
		if err != nil {
			return nil, err
		}
	}

	return &nodes.Constraint{
		Contype:       nodes.CONSTR_IDENTITY,
		GeneratedWhen: 'd',
		Options:       opts,
		Loc:           nodes.NoLoc(),
	}, nil
}
