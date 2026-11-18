package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// This file reads OceanBase's external-table forms: the tbl_name alternatives
// that name a location instead of a table (url_expr), the SOURCE(...) and
// FILES(...) table functions with their property lists, and the url_spec a
// LOAD DATA reads from. The properties, formats and patterns are OceanBase's
// own: MySQL has no external table.
//
// Ref: external_properties_list, external_properties, compression_key,
// external_file_format_list, external_file_format, opt_pattern, pattern_expr,
// format_expr, url_expr, url_table_function_expr and location_expr in
// sql_parser_mysql_mode.y (oceanbase v4.3.5_CE_BP6, lines 9319-9560 and
// 14499-14549). The rule's rows in testdata/compat.tsv hold each alternative,
// and its near-misses, to the pinned server.

// atURLExpr reports whether the current token opens tbl_name's url_expr: a
// quoted path, SOURCE(...) or FILES(...).
func (p *Parser) atURLExpr() bool {
	if (p.cur.Type == kwSOURCE || p.cur.Type == kwFILES) && p.peekNext().Type == '(' {
		return true
	}
	return p.cur.Type == tokSCONST && p.peekNext().Type == '('
}

// atURLTableFunction reports whether the current token opens
// url_table_function_expr: SOURCE(...) or FILES(...), and not a quoted path.
func (p *Parser) atURLTableFunction() bool {
	return (p.cur.Type == kwSOURCE || p.cur.Type == kwFILES) && p.peekNext().Type == '('
}

// parseExternalTableRef parses url_expr opt_as_alias, the tbl_name alternative
// that reads a table from a location:
//
//	url_expr: STRING_VALUE '(' format_expr opt_pattern ')' | url_table_function_expr
//	url_table_function_expr:
//	    SOURCE '(' external_properties_list ')'
//	  | FILES '(' location_expr opt_comma format_expr [opt_comma pattern_expr] ')'
//	opt_as_alias: /* EMPTY */ | AS relation_name
//
// A path names its location itself; SOURCE(...) carries properties instead; and
// FILES(...) names its location with LOCATION.
func (p *Parser) parseExternalTableRef() (*nodes.ExternalTableRef, error) {
	ref, err := p.parseExternalTableRefBody()
	if err != nil {
		return nil, err
	}
	return p.externalTableRefAlias(ref)
}

// parseExternalTableRefBody parses url_expr alone, without its opt_as_alias,
// which url_spec never takes.
func (p *Parser) parseExternalTableRefBody() (*nodes.ExternalTableRef, error) {
	ref := &nodes.ExternalTableRef{Loc: nodes.Loc{Start: p.pos()}}
	switch {
	case p.cur.Type == tokSCONST:
		ref.Path = p.advance().Str
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		format, err := p.parseExternalFileFormatList()
		if err != nil {
			return nil, err
		}
		ref.Format = format
		if p.isBareWord("PATTERN") {
			pattern, err := p.parseExternalPattern()
			if err != nil {
				return nil, err
			}
			ref.Pattern, ref.PatternSet = pattern, true
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	case p.cur.Type == kwFILES:
		ref.Function = "FILES"
		p.advance()
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		// location_expr: LOCATION opt_equal_mark STRING_VALUE.
		if !p.isBareWord("LOCATION") {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		p.match('=')
		if p.cur.Type != tokSCONST {
			return nil, p.syntaxErrorAtCur()
		}
		ref.Path = p.advance().Str
		p.match(',')
		if !p.isBareWord("FORMAT") {
			return nil, p.syntaxErrorAtCur()
		}
		format, err := p.parseExternalFileFormatList()
		if err != nil {
			return nil, err
		}
		ref.Format = format
		// opt_comma: the pattern may follow with or without a comma, and a comma
		// with nothing after it is a 1064.
		comma := false
		if _, ok := p.match(','); ok {
			comma = true
		}
		if p.isBareWord("PATTERN") {
			pattern, err := p.parseExternalPattern()
			if err != nil {
				return nil, err
			}
			ref.Pattern, ref.PatternSet = pattern, true
		} else if comma {
			return nil, p.syntaxErrorAtCur()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	default:
		// SOURCE '(' external_properties_list ')'.
		ref.Function = "SOURCE"
		p.advance()
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		source, err := p.parseExternalProperties()
		if err != nil {
			return nil, err
		}
		ref.Source = source
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}
	ref.Loc.End = p.prev.End
	return ref, nil
}

// externalTableRefAlias reads the AS relation_name a url_expr may carry. The
// alias is written with AS: `DELETE FROM SOURCE(...) a` is a 1064 where
// `DELETE FROM SOURCE(...) AS a` parses (opt_as_alias).
func (p *Parser) externalTableRefAlias(ref *nodes.ExternalTableRef) (*nodes.ExternalTableRef, error) {
	if _, ok := p.match(kwAS); !ok {
		return ref, nil
	}
	alias, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	ref.Alias = alias
	ref.Loc.End = p.prev.End
	return ref, nil
}

// parseURLSpec parses LOAD DATA's url_spec, at SOURCE, FILES or '(':
//
//	url_spec: url_table_function_expr
//	        | '(' select_with_opt_hint select_expr_list FROM url_table_function_expr ')'
//
// The query form reads its columns from the source, so it is kept as the SELECT
// the grammar builds.
func (p *Parser) parseURLSpec() (*nodes.ExternalTableRef, *nodes.SelectStmt, error) {
	if !p.atURLTableFunction() {
		if _, err := p.expect('('); err != nil {
			return nil, nil, err
		}
		// The query's Loc starts at its SELECT, as every SelectStmt's does, not
		// at the url_spec's parenthesis.
		start := p.pos()
		if _, err := p.expect(kwSELECT); err != nil {
			return nil, nil, err
		}
		targets, err := p.parseSelectExprList()
		if err != nil {
			return nil, nil, err
		}
		if _, err := p.expect(kwFROM); err != nil {
			return nil, nil, err
		}
		ref, err := p.parseURLTableFunction()
		if err != nil {
			return nil, nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, nil, err
		}
		return nil, &nodes.SelectStmt{
			Loc:        nodes.Loc{Start: start, End: p.prev.End},
			TargetList: targets,
			From:       []nodes.TableExpr{ref},
		}, nil
	}
	ref, err := p.parseURLTableFunction()
	if err != nil {
		return nil, nil, err
	}
	return ref, nil, nil
}

// parseURLTableFunction parses url_table_function_expr, which is what url_spec
// ends in: SOURCE(...) or FILES(...), with no quoted path and no alias.
func (p *Parser) parseURLTableFunction() (*nodes.ExternalTableRef, error) {
	if !p.atURLTableFunction() {
		return nil, p.syntaxErrorAtCur()
	}
	return p.parseExternalTableRefBody()
}

// parseExternalProperties parses external_properties_list, after the '(' of
// SOURCE(...):
//
//	external_properties: TYPE | ACCESSTYPE | ACCESSID | ACCESSKEY | STSTOKEN |
//	    ENDPOINT | TUNNEL_ENDPOINT | PROJECT_NAME | SCHEMA_NAME | QUOTA_NAME |
//	    TABLE_NAME | REGION | API_MODE | SPLIT = STRING_VALUE
//	  | {COMPRESSION_CODE | COMPRESSION} = compression_name
//	  | COLLECT_STATISTICS_ON_CREATE = BOOL_VALUE
//
// A comma between two entries is optional.
func (p *Parser) parseExternalProperties() ([]*nodes.ExternalProperty, error) {
	var props []*nodes.ExternalProperty
	for {
		if p.cur.Type == ')' {
			if len(props) == 0 {
				return nil, p.syntaxErrorAtCur()
			}
			break
		}
		start := p.pos()
		word := strings.ToUpper(p.cur.Str)
		var value nodes.ExprNode
		switch word {
		case "TYPE", "ACCESSTYPE", "ACCESSID", "ACCESSKEY", "STSTOKEN", "ENDPOINT",
			"TUNNEL_ENDPOINT", "PROJECT_NAME", "SCHEMA_NAME", "QUOTA_NAME",
			"TABLE_NAME", "REGION", "API_MODE", "SPLIT":
			v, err := p.parseExternalStringValue()
			if err != nil {
				return nil, err
			}
			value = v
		case "COMPRESSION_CODE", "COMPRESSION":
			// compression_key COMP_EQ compression_name: the word COMPRESSION_CODE
			// or COMPRESSION itself, then a name.
			p.advance()
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			v, err := p.parseExternalCompressionName()
			if err != nil {
				return nil, err
			}
			value = v
		case "COLLECT_STATISTICS_ON_CREATE":
			p.advance()
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			v, err := p.parseExternalBoolValue()
			if err != nil {
				return nil, err
			}
			value = v
		default:
			return nil, p.syntaxErrorAtCur()
		}
		props = append(props, &nodes.ExternalProperty{
			Loc:   nodes.Loc{Start: start, End: p.prev.End},
			Name:  word,
			Value: value,
		})
		if p.cur.Type == ',' {
			p.advance()
			if p.cur.Type == ')' {
				return nil, p.syntaxErrorAtCur()
			}
		}
	}
	return props, nil
}

// parseExternalFileFormatList parses format_expr, at its FORMAT:
//
//	format_expr: FORMAT opt_equal_mark '(' external_file_format_list ')'
func (p *Parser) parseExternalFileFormatList() ([]*nodes.ExternalProperty, error) {
	if !p.isBareWord("FORMAT") {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	p.match('=')
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var props []*nodes.ExternalProperty
	for {
		if p.cur.Type == ')' {
			if len(props) == 0 {
				return nil, p.syntaxErrorAtCur()
			}
			break
		}
		prop, err := p.parseExternalFileFormat()
		if err != nil {
			return nil, err
		}
		props = append(props, prop)
		if p.cur.Type == ',' {
			p.advance()
			if p.cur.Type == ')' {
				return nil, p.syntaxErrorAtCur()
			}
		}
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return props, nil
}

// parseExternalFileFormat parses one external_file_format. The value's shape
// follows the word: a string, an expression, a number, a boolean, a
// parenthesized list, or a name.
//
// Ref: external_file_format in sql_parser_mysql_mode.y
// (oceanbase v4.3.5_CE_BP6, lines 9436-9524).
func (p *Parser) parseExternalFileFormat() (*nodes.ExternalProperty, error) {
	start := p.pos()
	word := strings.ToUpper(p.cur.Str)
	var value nodes.ExprNode
	switch word {
	case "TYPE", "ENCODING", "COLUMN_INDEX_TYPE", "FILE_EXTENSION":
		v, err := p.parseExternalStringValue()
		if err != nil {
			return nil, err
		}
		value = v
	case "FIELD_DELIMITER", "LINE_DELIMITER", "ESCAPE",
		"FIELD_OPTIONALLY_ENCLOSED_BY", "FIELD_ENCLOSED_BY":
		p.advance()
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		v, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if v == nil {
			return nil, p.syntaxErrorAtCur()
		}
		value = v
	case "SKIP_HEADER", "ROW_INDEX_STRIDE":
		p.advance()
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		if p.cur.Type != tokICONST {
			return nil, p.syntaxErrorAtCur()
		}
		value = &nodes.IntLit{Loc: nodes.Loc{Start: p.cur.Loc, End: p.cur.End}, Value: p.cur.Ival}
		p.advance()
	case "SKIP_BLANK_LINES", "TRIM_SPACE", "EMPTY_FIELD_AS_NULL", "PARSE_HEADER",
		"IGNORE_LAST_EMPTY_COLUMN":
		p.advance()
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		v, err := p.parseExternalBoolValue()
		if err != nil {
			return nil, err
		}
		value = v
	case "NULL_IF":
		p.advance()
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		v, err := p.parseExternalExprList()
		if err != nil {
			return nil, err
		}
		value = v
	case "ROW_GROUP_SIZE", "STRIPE_SIZE", "COMPRESSION_BLOCK_SIZE":
		// file_size_const: INTNUM | STRING_VALUE.
		p.advance()
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		switch p.cur.Type {
		case tokICONST:
			value = &nodes.IntLit{Loc: nodes.Loc{Start: p.cur.Loc, End: p.cur.End}, Value: p.cur.Ival}
			p.advance()
		case tokSCONST:
			value = externalStringLit(p.advance())
		default:
			return nil, p.syntaxErrorAtCur()
		}
	case "COMPRESSION", "BINARY_FORMAT":
		p.advance()
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		v, err := p.parseExternalCompressionName()
		if err != nil {
			return nil, err
		}
		value = v
	case "COLUMN_BLOOM_FILTER":
		p.advance()
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		listStart := p.pos()
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		var items []nodes.ExprNode
		for {
			if p.cur.Type != tokICONST {
				return nil, p.syntaxErrorAtCur()
			}
			items = append(items, &nodes.IntLit{Loc: nodes.Loc{Start: p.cur.Loc, End: p.cur.End}, Value: p.cur.Ival})
			p.advance()
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		value = &nodes.RowExpr{Loc: nodes.Loc{Start: listStart, End: p.prev.End}, Items: items}
	default:
		return nil, p.syntaxErrorAtCur()
	}
	return &nodes.ExternalProperty{
		Loc:   nodes.Loc{Start: start, End: p.prev.End},
		Name:  word,
		Value: value,
	}, nil
}

// parseExternalStringValue reads the COMP_EQ STRING_VALUE that follows a
// property's word.
func (p *Parser) parseExternalStringValue() (nodes.ExprNode, error) {
	p.advance() // the word
	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	if p.cur.Type != tokSCONST {
		return nil, p.syntaxErrorAtCur()
	}
	return externalStringLit(p.advance()), nil
}

// parseExternalBoolValue reads the COMP_EQ BOOL_VALUE that follows a property's
// word.
func (p *Parser) parseExternalBoolValue() (nodes.ExprNode, error) {
	if !p.isBareWord("TRUE") && !p.isBareWord("FALSE") {
		return nil, p.syntaxErrorAtCur()
	}
	tok := p.advance()
	return &nodes.BoolLit{
		Loc:   nodes.Loc{Start: tok.Loc, End: tok.End},
		Value: strings.EqualFold(tok.Str, "TRUE"),
	}, nil
}

// parseExternalCompressionName reads compression_name and binary_format, which
// are NAME_OB, STRING_VALUE or an unreserved keyword.
func (p *Parser) parseExternalCompressionName() (nodes.ExprNode, error) {
	if p.cur.Type == tokSCONST {
		return externalStringLit(p.advance()), nil
	}
	name, _, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	return &nodes.ColumnRef{Loc: nodes.Loc{Start: p.prev.Loc, End: p.prev.End}, Column: name}, nil
}

// parseExternalExprList reads the '(' expr_list ')' that NULL_IF takes, as a
// RowExpr of its expressions.
func (p *Parser) parseExternalExprList() (nodes.ExprNode, error) {
	start := p.pos()
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var items []nodes.ExprNode
	for {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		items = append(items, expr)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.RowExpr{Loc: nodes.Loc{Start: start, End: p.prev.End}, Items: items}, nil
}

// parseExternalPattern reads pattern_expr: PATTERN opt_equal_mark STRING_VALUE.
func (p *Parser) parseExternalPattern() (string, error) {
	p.advance() // PATTERN
	p.match('=')
	if p.cur.Type != tokSCONST {
		return "", p.syntaxErrorAtCur()
	}
	return p.advance().Str, nil
}

// externalStringLit makes a string literal of a quoted token.
func externalStringLit(tok Token) *nodes.StringLit {
	return &nodes.StringLit{Loc: nodes.Loc{Start: tok.Loc, End: tok.End}, Value: tok.Str}
}
