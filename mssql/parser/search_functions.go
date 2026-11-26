// Package parser - search_functions.go parses T-SQL full-text search
// predicates and semantic / full-text rowset functions:
//
//	CONTAINS / FREETEXT                       (search_condition predicate)
//	CONTAINSTABLE / FREETEXTTABLE             (FROM table source)
//	SEMANTICKEYPHRASETABLE /
//	  SEMANTICSIMILARITYTABLE /
//	  SEMANTICSIMILARITYDETAILSTABLE          (FROM table source)
//
// These keywords are all CoreKeyword tokens in T-SQL: they can only appear
// at the specific grammar positions handled below. They are NOT generic
// function-call scalar expressions.
//
// AST aligns with SqlScriptDOM:
//
//	FullTextPredicate     -> ast.FullTextPredicate
//	FullTextTableReference-> ast.FullTextTableRef
//	SemanticTableReference-> ast.SemanticTableRef
//
// Grammar references: TSql170.g: fulltextPredicate, fulltextTableReference,
// semanticKeyPhraseTableReference, semanticSimilarityTableReference,
// semanticSimilarityDetailsTableReference.
package parser

import (
	nodes "github.com/bytebase/omni/mssql/ast"
)

// parseFullTextPredicate parses a CONTAINS(...) or FREETEXT(...) predicate.
// Caller must ensure the current token is kwCONTAINS or kwFREETEXT.
//
//	fulltext_predicate =
//	    ( CONTAINS | FREETEXT ) '('
//	        ( column_ref
//	        | '*'
//	        | '(' ( '*' | column_ref ( ',' column_ref )* ) ')'
//	        | PROPERTY '(' column_ref ',' string_literal ')'
//	        ) ','
//	        ( string_literal | variable )
//	        [ ',' LANGUAGE language_term ]
//	    ')'
func (p *Parser) parseFullTextPredicate() (nodes.ExprNode, error) {
	loc := p.pos()
	result := &nodes.FullTextPredicate{Loc: nodes.Loc{Start: loc, End: -1}}
	switch p.cur.Type {
	case kwCONTAINS:
		result.Func = nodes.FullTextContains
	case kwFREETEXT:
		result.Func = nodes.FullTextFreeText
	default:
		return nil, p.unexpectedToken()
	}
	p.advance() // consume CONTAINS / FREETEXT

	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	cols, propName, err := p.parseFullTextColumnsArg()
	if err != nil {
		return nil, err
	}
	result.Columns = cols
	result.PropertyName = propName

	if _, err := p.expect(','); err != nil {
		return nil, err
	}

	val, err := p.parseFullTextValue()
	if err != nil {
		return nil, err
	}
	result.Value = val

	// Optional LANGUAGE <expr>
	if _, ok := p.match(','); ok {
		if _, err := p.expect(kwLANGUAGE); err != nil {
			return nil, err
		}
		lang, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if lang == nil {
			return nil, p.unexpectedToken()
		}
		result.LanguageTerm = lang
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	result.Loc.End = p.prevEnd()
	return result, nil
}

// parseFullTextColumnsArg parses the column portion of a full-text predicate
// or full-text table reference. Handles:
//   - single column reference
//   - '*'
//   - '(' col_list ')' or '(' '*' ')'
//   - PROPERTY '(' col_ref ',' string_literal ')'
//
// Returns the columns list and, when PROPERTY was used, the property-name
// literal (else nil).
func (p *Parser) parseFullTextColumnsArg() (*nodes.List, *nodes.Literal, error) {
	// PROPERTY ( col_ref , string_literal )
	if p.cur.Type == kwPROPERTY && p.peekNext().Type == '(' {
		p.advance() // PROPERTY
		p.advance() // (
		col, err := p.parseFullTextColumn()
		if err != nil {
			return nil, nil, err
		}
		if _, err := p.expect(','); err != nil {
			return nil, nil, err
		}
		lit, err := p.parseStringLiteral()
		if err != nil {
			return nil, nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, nil, err
		}
		return &nodes.List{Items: []nodes.Node{col}}, lit, nil
	}

	// '(' ... ')'
	if p.cur.Type == '(' {
		p.advance()
		var items []nodes.Node
		for {
			col, err := p.parseFullTextColumn()
			if err != nil {
				return nil, nil, err
			}
			items = append(items, col)
			if _, ok := p.match(','); !ok {
				break
			}
		}
		if _, err := p.expect(')'); err != nil {
			return nil, nil, err
		}
		return &nodes.List{Items: items}, nil, nil
	}

	// Single column or star
	col, err := p.parseFullTextColumn()
	if err != nil {
		return nil, nil, err
	}
	return &nodes.List{Items: []nodes.Node{col}}, nil, nil
}

// parseFullTextColumn parses one of:
//   - '*'
//   - [schema.]column  (qualified column reference)
func (p *Parser) parseFullTextColumn() (nodes.ExprNode, error) {
	if p.cur.Type == '*' {
		loc := p.pos()
		tok := p.advance()
		_ = tok
		return &nodes.StarExpr{
			Loc: nodes.Loc{Start: loc, End: p.prevEnd()},
		}, nil
	}
	loc := p.pos()
	name, ok := p.parseIdentifier()
	if !ok {
		return nil, p.unexpectedToken()
	}
	ref := &nodes.ColumnRef{Column: name, Loc: nodes.Loc{Start: loc, End: -1}}
	for p.cur.Type == '.' {
		p.advance()
		if p.cur.Type == '*' {
			p.advance()
			star := &nodes.StarExpr{
				Qualifier: ref.Column,
				Loc:       nodes.Loc{Start: loc, End: p.prevEnd()},
			}
			return star, nil
		}
		part, ok := p.parseIdentifier()
		if !ok {
			return nil, p.unexpectedToken()
		}
		// Shift parts left: old Column becomes Table, etc.
		ref.Server = ref.Database
		ref.Database = ref.Schema
		ref.Schema = ref.Table
		ref.Table = ref.Column
		ref.Column = part
	}
	ref.Loc.End = p.prevEnd()
	return ref, nil
}

// parseStringLiteral consumes a string literal token and returns a *Literal.
func (p *Parser) parseStringLiteral() (*nodes.Literal, error) {
	loc := p.pos()
	switch p.cur.Type {
	case tokSCONST:
		tok := p.advance()
		return &nodes.Literal{
			Type: nodes.LitString,
			Str:  tok.Str,
			Loc:  nodes.Loc{Start: loc, End: p.prevEnd()},
		}, nil
	case tokNSCONST:
		tok := p.advance()
		return &nodes.Literal{
			Type:    nodes.LitString,
			Str:     tok.Str,
			IsNChar: true,
			Loc:     nodes.Loc{Start: loc, End: p.prevEnd()},
		}, nil
	}
	return nil, p.unexpectedToken()
}

// parseFullTextValue parses the search value: a string literal or a variable.
func (p *Parser) parseFullTextValue() (nodes.ExprNode, error) {
	switch p.cur.Type {
	case tokSCONST, tokNSCONST:
		return p.parseStringLiteral()
	case tokVARIABLE:
		loc := p.pos()
		tok := p.advance()
		return &nodes.VariableRef{
			Name: tok.Str,
			Loc:  nodes.Loc{Start: loc, End: p.prevEnd()},
		}, nil
	}
	return nil, p.unexpectedToken()
}

// parseFullTextTableRef parses CONTAINSTABLE(...) / FREETEXTTABLE(...) used as
// a table source. Caller must ensure the current token is kwCONTAINSTABLE or
// kwFREETEXTTABLE.
//
//	fulltext_table_ref =
//	    ( CONTAINSTABLE | FREETEXTTABLE ) '('
//	        table_ref ','
//	        columns_arg ','                    -- same shapes as fulltext_predicate
//	        ( string_literal | variable )
//	        [ ',' LANGUAGE language_term ]     -- order of TopN and LANGUAGE may swap
//	        [ ',' ( integer_literal | variable ) ]
//	    ')' [ [AS] alias ]
func (p *Parser) parseFullTextTableRef() (nodes.TableExpr, error) {
	loc := p.pos()
	result := &nodes.FullTextTableRef{Loc: nodes.Loc{Start: loc, End: -1}}
	switch p.cur.Type {
	case kwCONTAINSTABLE:
		result.Func = nodes.FullTextContains
	case kwFREETEXTTABLE:
		result.Func = nodes.FullTextFreeText
	default:
		return nil, p.unexpectedToken()
	}
	p.advance() // consume the function keyword

	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	tbl, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	if tbl == nil {
		return nil, p.unexpectedToken()
	}
	result.Table = tbl

	if _, err := p.expect(','); err != nil {
		return nil, err
	}

	cols, propName, err := p.parseFullTextColumnsArg()
	if err != nil {
		return nil, err
	}
	result.Columns = cols
	result.PropertyName = propName

	if _, err := p.expect(','); err != nil {
		return nil, err
	}
	sc, err := p.parseFullTextValue()
	if err != nil {
		return nil, err
	}
	result.SearchCondition = sc

	// Options: [ , LANGUAGE language_term ] [ , top_n_by_rank ], in that
	// order: SQL Server refuses top_n_by_rank before LANGUAGE
	// (CONTAINSTABLE(a, a, 'a', 1, LANGUAGE 0x00) is Msg 102, measured).
	if _, ok := p.match(','); ok {
		more := true
		if p.cur.Type == kwLANGUAGE {
			p.advance()
			lang, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if lang == nil {
				return nil, p.unexpectedToken()
			}
			result.Language = lang
			_, more = p.match(',')
		}
		if more {
			topN, err := p.parseFullTextTopN()
			if err != nil {
				return nil, err
			}
			result.TopN = topN
		}
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	result.Alias = p.parseOptionalAlias()
	result.Loc.End = p.prevEnd()
	return result, nil
}

// parseFullTextTopN parses the optional top-N integer-or-variable argument
// to CONTAINSTABLE/FREETEXTTABLE.
func (p *Parser) parseFullTextTopN() (nodes.ExprNode, error) {
	loc := p.pos()
	switch p.cur.Type {
	case tokICONST:
		// integerOrVariable is an integer, and the lexer makes a binary
		// constant (0x10) a tokICONST too. CONTAINSTABLE(a, a, 'a', 0x00)
		// is an error.
		if procIsBinary(p.cur) {
			return nil, p.unexpectedToken()
		}
		tok := p.advance()
		return &nodes.Literal{
			Type: nodes.LitInteger,
			Ival: tok.Ival,
			Str:  tok.Str,
			Loc:  nodes.Loc{Start: loc, End: p.prevEnd()},
		}, nil
	case tokVARIABLE:
		tok := p.advance()
		return &nodes.VariableRef{
			Name: tok.Str,
			Loc:  nodes.Loc{Start: loc, End: p.prevEnd()},
		}, nil
	case '-':
		// Negative integer — allow for symmetry; SqlScriptDOM's
		// possibleNegativeConstant accepts '-' prefix for semantic table args
		// but not here; reject to stay strict.
		return nil, p.unexpectedToken()
	}
	return nil, p.unexpectedToken()
}

// parseSemanticTableRef parses one of SEMANTICKEYPHRASETABLE /
// SEMANTICSIMILARITYTABLE / SEMANTICSIMILARITYDETAILSTABLE. Caller must
// ensure the current token is one of those keywords.
//
//	semantic_key_phrase_table =
//	    SEMANTICKEYPHRASETABLE '(' table_ref ','
//	        ( column | '*' | '(' col_list ')' )
//	        [ ',' integer_or_variable ]
//	    ')' [ [AS] alias ]
//
//	semantic_similarity_table =
//	    SEMANTICSIMILARITYTABLE '(' table_ref ','
//	        ( column | '*' | '(' col_list ')' ) ','
//	        integer_or_variable
//	    ')' [ [AS] alias ]
//
//	semantic_similarity_details_table =
//	    SEMANTICSIMILARITYDETAILSTABLE '(' table_ref ','
//	        column ',' integer_or_variable ','
//	        column ',' integer_or_variable
//	    ')' [ [AS] alias ]
func (p *Parser) parseSemanticTableRef() (nodes.TableExpr, error) {
	loc := p.pos()
	result := &nodes.SemanticTableRef{Loc: nodes.Loc{Start: loc, End: -1}}
	var kind nodes.SemanticFunc
	switch p.cur.Type {
	case kwSEMANTICKEYPHRASETABLE:
		kind = nodes.SemanticKeyPhraseTable
	case kwSEMANTICSIMILARITYTABLE:
		kind = nodes.SemanticSimilarityTable
	case kwSEMANTICSIMILARITYDETAILSTABLE:
		kind = nodes.SemanticSimilarityDetailsTable
	default:
		return nil, p.unexpectedToken()
	}
	result.Func = kind
	p.advance()

	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	tbl, err := p.parseTableRef()
	if err != nil {
		return nil, err
	}
	if tbl == nil {
		return nil, p.unexpectedToken()
	}
	result.Table = tbl

	if _, err := p.expect(','); err != nil {
		return nil, err
	}

	switch kind {
	case nodes.SemanticKeyPhraseTable:
		cols, _, err := p.parseFullTextColumnsArg()
		if err != nil {
			return nil, err
		}
		result.Columns = cols
		// Optional , source_key
		if _, ok := p.match(','); ok {
			sk, err := p.parseSemanticKey()
			if err != nil {
				return nil, err
			}
			result.SourceKey = sk
		}
	case nodes.SemanticSimilarityTable:
		cols, _, err := p.parseFullTextColumnsArg()
		if err != nil {
			return nil, err
		}
		result.Columns = cols
		if _, err := p.expect(','); err != nil {
			return nil, err
		}
		sk, err := p.parseSemanticKey()
		if err != nil {
			return nil, err
		}
		result.SourceKey = sk
	case nodes.SemanticSimilarityDetailsTable:
		// source_column , source_key , matched_column , matched_key
		srcCol, err := p.parseFullTextColumn()
		if err != nil {
			return nil, err
		}
		srcRef, ok := srcCol.(*nodes.ColumnRef)
		if !ok {
			return nil, p.unexpectedToken()
		}
		result.Columns = &nodes.List{Items: []nodes.Node{srcCol}}
		if _, err := p.expect(','); err != nil {
			return nil, err
		}
		sk, err := p.parseSemanticKey()
		if err != nil {
			return nil, err
		}
		result.SourceKey = sk
		if _, err := p.expect(','); err != nil {
			return nil, err
		}
		matchedCol, err := p.parseFullTextColumn()
		if err != nil {
			return nil, err
		}
		matchedRef, ok := matchedCol.(*nodes.ColumnRef)
		if !ok {
			return nil, p.unexpectedToken()
		}
		result.MatchedColumn = matchedRef
		if _, err := p.expect(','); err != nil {
			return nil, err
		}
		mk, err := p.parseSemanticKey()
		if err != nil {
			return nil, err
		}
		result.MatchedKey = mk
		_ = srcRef
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	result.Alias = p.parseOptionalAlias()
	result.Loc.End = p.prevEnd()
	return result, nil
}

// parseSemanticKey parses the source_key and matched_key arguments of the
// semantic table functions, a possibleNegativeConstant
// (semanticKeyPhraseTableReference, TSql160.g:20414): any literal or a
// negative number, not a name or an expression (Msg 102, measured).
func (p *Parser) parseSemanticKey() (nodes.ExprNode, error) {
	return p.parsePossibleNegativeConstant(true)
}
