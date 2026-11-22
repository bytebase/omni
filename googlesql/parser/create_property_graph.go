package parser

import (
	"strings"

	"github.com/bytebase/omni/googlesql/ast"
)

// ---------------------------------------------------------------------------
// CREATE PROPERTY GRAPH (parser-gql node)
// ---------------------------------------------------------------------------
//
// This file ports create_property_graph_statement and its element-table grammar
// (GoogleSQL.g4; googlesql.tm:2662 — element_table_list /
// element_table_definition / label_and_properties / properties_clause /
// derived_property_list / source+dest node-table clauses) — a hand-port of
// ZetaSQL's GoogleSQL graph DDL (BigQuery / Spanner Graph):
//
//	CREATE [OR REPLACE] PROPERTY GRAPH [IF NOT EXISTS] <path>
//	  NODE TABLES ( <element_table>, … )
//	  [EDGE TABLES ( <element_table>, … )]
//	  [OPTIONS(…)]
//
// The OPTIONS list follows the graph content: googlesql.tm:2666 ends the rule
// with edge_table_clause? opt_options_list, and Spanner's graph DDL reference
// writes property_graph_content [OPTIONS (key=value, …)]. The legacy ANTLR
// GoogleSQLParser.g4 put the list before NODE TABLES; the pinned grammar and the
// emulator do not, and neither derives the leading position.
//
// One omni parser serves both dialects; it accepts the BigQuery+Spanner union.

// parseCreatePropertyGraph parses the body of a CREATE PROPERTY GRAPH statement
// after the shared CREATE prefix (OR REPLACE) has been consumed and cur is at
// the PROPERTY keyword. create_property_graph_statement has NO opt_create_scope,
// so the caller rejects a leading TEMP/TEMPORARY/PUBLIC/PRIVATE before reaching
// here.
//
//	PROPERTY GRAPH opt_if_not_exists path_expression
//	  NODE TABLES element_table_list opt_edge_table_clause? opt_options_list?
func (p *Parser) parseCreatePropertyGraph(create Token, orReplace bool) (ast.Node, error) {
	p.advance() // PROPERTY
	if _, err := p.expect(kwGRAPH); err != nil {
		return nil, err
	}

	stmt := &ast.CreatePropertyGraphStmt{OrReplace: orReplace}
	stmt.Loc.Start = create.Loc.Start

	ifTok := p.cur // anchor for the OR-REPLACE/IF-NOT-EXISTS conflict diagnostic
	ifNotExists, err := p.parseIfNotExists()
	if err != nil {
		return nil, err
	}
	stmt.IfNotExists = ifNotExists

	// OR REPLACE and IF NOT EXISTS are mutually exclusive (a GoogleSQL-wide
	// invariant — you cannot both replace an object and skip-if-it-exists).
	// DIVERGENCE from the legacy .g4, which (a) requires IF NOT EXISTS on every
	// CREATE PROPERTY GRAPH — `opt_or_replace? PROPERTY GRAPH opt_if_not_exists
	// path_expression`, where opt_if_not_exists carries NO `?` (a hand-port slip:
	// every other CREATE makes it optional) — and (b) permits OR REPLACE
	// alongside it. The live Spanner emulator is authoritative for this Spanner-
	// Graph form: it makes IF NOT EXISTS OPTIONAL and REJECTS the OR REPLACE + IF
	// NOT EXISTS combination with a true DDL parse error ("CREATE PROPERTY GRAPH
	// IF NOT EXISTS cannot be used with other existence modifiers such as
	// `OR REPLACE`"). We follow the oracle: IF NOT EXISTS optional, the pair
	// rejected. See divergence ledger + the suite:graph-query rows of testdata/compat.tsv.
	if orReplace && ifNotExists && p.isSpanner() {
		return nil, &ParseError{
			Position: ifTok.Loc.Start, End: ifTok.Loc.End,
			Message: "syntax error: CREATE PROPERTY GRAPH IF NOT EXISTS cannot be used with OR REPLACE",
		}
	}

	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	stmt.Loc.End = name.Loc.End

	// NODE TABLES <element_table_list> — required in Spanner; optional in
	// BigQuery/union mode (BigQuery accepts `CREATE PROPERTY GRAPH ds.g` and
	// `CREATE PROPERTY GRAPH ds.g EDGE TABLES (...)`).
	if p.cur.Type == kwNODE || p.isSpanner() {
		if _, err := p.expect(kwNODE); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwTABLES); err != nil {
			return nil, err
		}
		nodeTables, end, err := p.parseElementTableList()
		if err != nil {
			return nil, err
		}
		stmt.NodeTables = nodeTables
		stmt.Loc.End = end
	}

	// opt_edge_table_clause?  — EDGE TABLES <element_table_list>.
	if p.cur.Type == kwEDGE {
		p.advance() // EDGE
		if _, err := p.expect(kwTABLES); err != nil {
			return nil, err
		}
		edgeTables, eend, err := p.parseElementTableList()
		if err != nil {
			return nil, err
		}
		stmt.EdgeTables = edgeTables
		stmt.Loc.End = eend
	}

	// opt_options_list?  — OPTIONS ( … ), after the graph content.
	if p.cur.Type == kwOPTIONS {
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
		stmt.Loc.End = p.prev.Loc.End
	}

	return stmt, nil
}

// parseElementTableList parses an element_table_list:
//
//	( <element_table_definition> (, <element_table_definition>)* [,] )
//
// A trailing comma before the close paren is permitted by the grammar. Returns
// the definitions (>= 1) and the end offset of the close paren.
func (p *Parser) parseElementTableList() ([]*ast.ElementTableDef, int, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, 0, err
	}
	var defs []*ast.ElementTableDef
	for {
		def, err := p.parseElementTableDef()
		if err != nil {
			return nil, 0, err
		}
		defs = append(defs, def)
		if p.cur.Type != int(',') {
			break
		}
		p.advance() // ,
		// Trailing comma: `, )` closes the list.
		if p.cur.Type == int(')') {
			break
		}
	}
	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, 0, err
	}
	return defs, closeTok.Loc.End, nil
}

// parseElementTableDef parses an element_table_definition:
//
//	<path> [AS alias] [KEY ( cols )]
//	  [SOURCE KEY ( cols ) REFERENCES <node> [( cols )]]
//	  [DESTINATION KEY ( cols ) REFERENCES <node> [( cols )]]
//	  [<opt_label_and_properties_clause>]
func (p *Parser) parseElementTableDef() (*ast.ElementTableDef, error) {
	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	def := &ast.ElementTableDef{Name: name, Loc: name.Loc}
	def.Loc.End = name.Loc.End

	// opt_as_alias_with_required_as?: AS identifier.
	if p.cur.Type == kwAS {
		p.advance() // AS
		tok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		alias, err := p.identifierText(tok)
		if err != nil {
			return nil, err
		}
		def.Alias = alias
		def.Loc.End = tok.Loc.End
	}

	// opt_key_clause?: KEY ( cols ).
	if p.cur.Type == kwKEY {
		p.advance() // KEY
		cols, err := p.parseColumnList()
		if err != nil {
			return nil, err
		}
		def.Key = cols
		def.Loc.End = p.prev.Loc.End
	}

	// opt_source_node_table_clause?: SOURCE KEY ( cols ) REFERENCES node [( cols )].
	if p.cur.Type == kwSOURCE {
		ref, err := p.parseNodeTableRef(kwSOURCE)
		if err != nil {
			return nil, err
		}
		def.Source = ref
		def.Loc.End = ref.Loc.End
	}

	// opt_dest_node_table_clause?: DESTINATION KEY ( cols ) REFERENCES node [( cols )].
	if p.cur.Type == kwDESTINATION {
		ref, err := p.parseNodeTableRef(kwDESTINATION)
		if err != nil {
			return nil, err
		}
		def.Dest = ref
		def.Loc.End = ref.Loc.End
	}

	// opt_options_list[default_label_options].
	if p.cur.Type == kwOPTIONS {
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		def.DefaultLabelOptions = opts
		def.Loc.End = p.prev.Loc.End
	}

	// opt_label_and_properties_clause?: properties_clause | label_and_properties+.
	labels, err := p.parseLabelAndPropertiesClause()
	if err != nil {
		return nil, err
	}
	if labels != nil {
		def.Labels = labels
		def.Loc.End = labels[len(labels)-1].Loc.End
	}

	// dynamic_label_and_properties?:
	//   ( DYNAMIC LABEL ( expr ) | DYNAMIC PROPERTIES ( expr ) )+
	// Duplicate DYNAMIC LABEL or DYNAMIC PROPERTIES is a syntax error
	// (MergeDynamicLabelProperties in googlesql.tm).
	for p.isUnquotedDynamic() {
		p.advance() // consume DYNAMIC
		switch p.cur.Type {
		case kwLABEL:
			if def.DynamicLabel != nil {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance() // LABEL
			if _, err := p.expect(int('(')); err != nil {
				return nil, err
			}
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			endTok, err := p.expect(int(')'))
			if err != nil {
				return nil, err
			}
			def.DynamicLabel = expr
			def.Loc.End = endTok.Loc.End
		case kwPROPERTIES:
			if def.DynamicProperties != nil {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance() // PROPERTIES
			if _, err := p.expect(int('(')); err != nil {
				return nil, err
			}
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			endTok, err := p.expect(int(')'))
			if err != nil {
				return nil, err
			}
			def.DynamicProperties = expr
			def.Loc.End = endTok.Loc.End
		default:
			return nil, p.syntaxErrorAtCur()
		}
	}

	return def, nil
}

func (p *Parser) isUnquotedDynamic() bool {
	return p.cur.Type == tokIdentifier &&
		p.cur.Loc.End-p.cur.Loc.Start == len("DYNAMIC") &&
		strings.EqualFold(p.cur.Str, "DYNAMIC")
}

// parseNodeTableRef parses a SOURCE / DESTINATION node-table reference
// (opt_source_node_table_clause / opt_dest_node_table_clause):
//
//	(SOURCE|DESTINATION) KEY ( cols ) REFERENCES <identifier> [( cols )]
//
// kw is the leading keyword (kwSOURCE or kwDESTINATION), at the current token.
func (p *Parser) parseNodeTableRef(kw int) (*ast.ElementTableRef, error) {
	lead := p.advance() // SOURCE | DESTINATION
	ref := &ast.ElementTableRef{Loc: lead.Loc}
	if _, err := p.expect(kwKEY); err != nil {
		return nil, err
	}
	cols, err := p.parseColumnList()
	if err != nil {
		return nil, err
	}
	ref.Columns = cols
	if _, err := p.expect(kwREFERENCES); err != nil {
		return nil, err
	}
	tok, err := p.expectIdentifier()
	if err != nil {
		return nil, err
	}
	node, err := p.identifierText(tok)
	if err != nil {
		return nil, err
	}
	ref.Node = node
	ref.Loc.End = tok.Loc.End
	// Optional referenced-column list: `( cols )`.
	if p.cur.Type == int('(') {
		refCols, err := p.parseColumnList()
		if err != nil {
			return nil, err
		}
		ref.RefColumns = refCols
		ref.Loc.End = p.prev.Loc.End
	}
	return ref, nil
}

// parseLabelAndPropertiesClause parses an opt_label_and_properties_clause:
//
//	properties_clause              (the bare form)
//	| label_and_properties+        (one or more `[DEFAULT] LABEL id [props]`)
//
// It returns nil when neither form is present (the clause is optional). The bare
// properties_clause form is normalized into a single LabelAndProperties entry
// with an empty LabelName.
func (p *Parser) parseLabelAndPropertiesClause() ([]*ast.LabelAndProperties, error) {
	// Bare properties_clause: starts with NO PROPERTIES | PROPERTIES … .
	if p.cur.Type == kwNO || p.cur.Type == kwPROPERTIES {
		lp := &ast.LabelAndProperties{Loc: p.cur.Loc}
		if err := p.parsePropertiesClause(lp); err != nil {
			return nil, err
		}
		return []*ast.LabelAndProperties{lp}, nil
	}
	// label_and_properties+: [DEFAULT] LABEL id [properties_clause].
	if p.cur.Type != kwDEFAULT && p.cur.Type != kwLABEL {
		return nil, nil
	}
	var labels []*ast.LabelAndProperties
	for p.cur.Type == kwDEFAULT || p.cur.Type == kwLABEL {
		lp, err := p.parseLabelAndProperties()
		if err != nil {
			return nil, err
		}
		labels = append(labels, lp)
	}
	return labels, nil
}

// parseLabelAndProperties parses one label_and_properties:
//
//	( DEFAULT LABEL opt_options_list | LABEL <identifier> ) [properties_clause]
func (p *Parser) parseLabelAndProperties() (*ast.LabelAndProperties, error) {
	lp := &ast.LabelAndProperties{Kind: ast.LabelPropsDefault, Loc: p.cur.Loc}
	start := p.cur.Loc.Start
	if p.cur.Type == kwDEFAULT {
		p.advance() // DEFAULT
		lblTok, err := p.expect(kwLABEL)
		if err != nil {
			return nil, err
		}
		lp.Default = true
		lp.Loc = ast.Loc{Start: start, End: lblTok.Loc.End}
		if p.cur.Type == kwOPTIONS {
			opts, err := p.parseOptionsList()
			if err != nil {
				return nil, err
			}
			lp.Options = opts
			lp.Loc.End = p.prev.Loc.End
		}
	} else {
		if _, err := p.expect(kwLABEL); err != nil {
			return nil, err
		}
		tok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		name, err := p.identifierText(tok)
		if err != nil {
			return nil, err
		}
		lp.LabelName = name
		lp.Loc = ast.Loc{Start: start, End: tok.Loc.End}
	}
	// Optional properties_clause.
	if p.cur.Type == kwNO || p.cur.Type == kwPROPERTIES {
		if err := p.parsePropertiesClause(lp); err != nil {
			return nil, err
		}
	}
	return lp, nil
}

// parsePropertiesClause parses a properties_clause into lp:
//
//	NO PROPERTIES
//	| PROPERTIES ARE? ALL COLUMNS [EXCEPT ( cols )]
//	| PROPERTIES ( <derived_property>, … )
//
// The cursor is at NO or PROPERTIES.
func (p *Parser) parsePropertiesClause(lp *ast.LabelAndProperties) error {
	if p.cur.Type == kwNO {
		p.advance() // NO
		end, err := p.expect(kwPROPERTIES)
		if err != nil {
			return err
		}
		lp.Kind = ast.LabelPropsNone
		lp.Loc.End = end.Loc.End
		return nil
	}
	// PROPERTIES …
	p.advance() // PROPERTIES
	// properties_all_columns: PROPERTIES ARE? ALL COLUMNS  — but the leading
	// PROPERTIES is already consumed. Distinguish ALL-COLUMNS from the
	// parenthesized derived-property list by the next token.
	switch {
	case p.cur.Type == kwARE || p.cur.Type == kwALL:
		// PROPERTIES [ARE] ALL COLUMNS [EXCEPT ( cols )].
		if p.cur.Type == kwARE {
			p.advance() // ARE
		}
		if _, err := p.expect(kwALL); err != nil {
			return err
		}
		end, err := p.expect(kwCOLUMNS)
		if err != nil {
			return err
		}
		lp.Kind = ast.LabelPropsAllColumns
		lp.Loc.End = end.Loc.End
		// opt_except_column_list: EXCEPT ( cols ).
		if p.cur.Type == kwEXCEPT {
			p.advance() // EXCEPT
			cols, err := p.parseColumnList()
			if err != nil {
				return err
			}
			lp.ExceptColumns = cols
			lp.Loc.End = p.prev.Loc.End
		}
		return nil
	case p.cur.Type == int('('):
		// PROPERTIES ( <derived_property>, … ).
		p.advance() // (
		lp.Kind = ast.LabelPropsList
		for {
			dp, err := p.parseDerivedProperty()
			if err != nil {
				return err
			}
			lp.PropsList = append(lp.PropsList, dp)
			if p.cur.Type != int(',') {
				break
			}
			p.advance() // ,
		}
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return err
		}
		lp.Loc.End = closeTok.Loc.End
		return nil
	default:
		return p.syntaxErrorAtCur()
	}
}

// parseDerivedProperty parses one derived_property:
// expression [AS alias] [OPTIONS (...)].
func (p *Parser) parseDerivedProperty() (*ast.DerivedProperty, error) {
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	dp := &ast.DerivedProperty{Expr: expr, Loc: ast.Loc{Start: nodeLoc(expr).Start, End: nodeLoc(expr).End}}
	if p.cur.Type == kwAS {
		p.advance() // AS
		tok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		alias, err := p.identifierText(tok)
		if err != nil {
			return nil, err
		}
		dp.Alias = alias
		dp.Loc.End = tok.Loc.End
	}
	if p.cur.Type == kwOPTIONS {
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		dp.Options = opts
		dp.Loc.End = p.prev.Loc.End
	}
	return dp, nil
}
