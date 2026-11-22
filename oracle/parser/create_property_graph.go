package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseGraphNameList reads ( name [, name ] ... ). When unique is true, a
// repeated name is rejected at the duplicate token (ORA-42416, ORA-42426).
func (p *Parser) parseGraphNameList() ([]string, error) {
	return p.parseGraphNameListUnique(false)
}

func (p *Parser) parseGraphNameListUnique(unique bool) ([]string, error) {
	if err := p.expectCreateTableToken('('); err != nil {
		return nil, err
	}
	var names []string
	var seen map[string]bool
	if unique {
		seen = map[string]bool{}
	}
	for {
		tok := p.cur
		names = append(names, tok.Str)
		if err := p.parseCreateTableName(); err != nil {
			return nil, err
		}
		if unique {
			if seen[tok.Str] {
				return nil, p.syntaxErrorAtTok(tok)
			}
			seen[tok.Str] = true
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance() // ','
	}
	if err := p.expectCreateTableToken(')'); err != nil {
		return nil, err
	}
	return names, nil
}

// parseCreatePropertyGraphStmt parses a CREATE PROPERTY GRAPH statement
// after GRAPH.
//
// BNF: oracle/parser/bnf/create-property-graph.bnf
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/create-property-graph.html
//
//	CREATE [ OR REPLACE ] PROPERTY GRAPH [ IF NOT EXISTS ]
//	    [ schema. ] graph_name
//	    vertex_tables_clause
//	    [ edge_tables_clause ]
//	    [ graph_options ] ;
//
//	vertex_tables_clause:
//	    VERTEX TABLES ( vertex_table_definition [, vertex_table_definition ]... )
//	edge_tables_clause:
//	    EDGE TABLES ( edge_table_definition [, edge_table_definition ]... )
//	graph_options:
//	    OPTIONS ( { ENFORCED MODE | TRUSTED MODE | { ALLOW | DISALLOW } MIXED PROPERTY TYPES } [, ...] )
//
// Oracle 23ai takes either option alone in OPTIONS and refuses the page's
// leading comma (OPTIONS(, ALLOW MIXED PROPERTY TYPES)).
func (p *Parser) parseCreatePropertyGraphStmt(start int, orReplace, ifNotExists bool) (*nodes.CreatePropertyGraphStmt, error) {
	stmt := &nodes.CreatePropertyGraphStmt{
		OrReplace:   orReplace,
		IfNotExists: ifNotExists,
		Loc:         nodes.Loc{Start: start},
	}

	// [ IF NOT EXISTS ]: Oracle answers an IF not followed by NOT EXISTS
	// at the IF.
	if !ifNotExists && p.cur.Type == kwIF {
		ifTok := p.cur
		p.advance() // IF
		if p.cur.Type != kwNOT || p.peekNext().Type != kwEXISTS {
			return nil, p.syntaxErrorAtTok(ifTok)
		}
		p.advance() // NOT
		p.advance() // EXISTS
		stmt.IfNotExists = true
	}

	if !p.isCreateTableNameStart() {
		return nil, p.syntaxErrorAtCur()
	}
	var err error
	if stmt.Name, err = p.parseObjectName(); err != nil {
		return nil, err
	}

	if err := p.expectCreateTableWord("VERTEX"); err != nil {
		return nil, err
	}
	if err := p.expectCreateTableWord("TABLES"); err != nil {
		return nil, err
	}
	if err := p.expectCreateTableToken('('); err != nil {
		return nil, err
	}
	seenTables := map[string]bool{}
	for {
		vtd, err := p.parseGraphTableDef(seenTables)
		if err != nil {
			return nil, err
		}
		stmt.VertexTables = append(stmt.VertexTables, vtd)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // ','
	}
	if err := p.expectCreateTableToken(')'); err != nil {
		return nil, err
	}

	if p.isKeywordStr("EDGE") {
		p.advance() // EDGE
		if err := p.expectCreateTableWord("TABLES"); err != nil {
			return nil, err
		}
		if err := p.expectCreateTableToken('('); err != nil {
			return nil, err
		}
		for {
			etd, err := p.parseGraphEdgeDef(seenTables)
			if err != nil {
				return nil, err
			}
			stmt.EdgeTables = append(stmt.EdgeTables, etd)
			if p.cur.Type != ',' {
				break
			}
			p.advance() // ','
		}
		if err := p.expectCreateTableToken(')'); err != nil {
			return nil, err
		}
	}

	if p.isKeywordStr("OPTIONS") {
		optStart := p.pos()
		p.advance() // OPTIONS
		if err := p.expectCreateTableToken('('); err != nil {
			return nil, err
		}
		opts := &nodes.GraphOptions{Loc: nodes.Loc{Start: optStart}}
		for {
			switch {
			case opts.Mode == "" && (p.isKeywordStr("ENFORCED") || p.isKeywordStr("TRUSTED")):
				opts.Mode = p.cur.Str
				p.advance()
				if err := p.expectCreateTableWord("MODE"); err != nil {
					return nil, err
				}
			case opts.MixedPropTypes == "" && (p.isKeywordStr("ALLOW") || p.isKeywordStr("DISALLOW")):
				opts.MixedPropTypes = p.cur.Str
				p.advance()
				for _, w := range []string{"MIXED", "PROPERTY", "TYPES"} {
					if err := p.expectCreateTableWord(w); err != nil {
						return nil, err
					}
				}
			default:
				return nil, p.syntaxErrorAtCur()
			}
			if p.cur.Type != ',' {
				break
			}
			p.advance() // ','
		}
		if err := p.expectCreateTableToken(')'); err != nil {
			return nil, err
		}
		opts.Loc.End = p.prev.End
		stmt.Options = opts
	}

	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseGraphElementNameAndKey reads graph_element_name_and_key.
//
//	[ schema. ] graph_element_object_name [ AS graph_element_name ]
//	[ KEY ( column_name [, column_name ]... ) ]
//
// Each element table in a property graph has a unique effective name (its AS
// alias if given, otherwise its unqualified object name; ORA-42404).
func (p *Parser) parseGraphElementNameAndKey(seenTables map[string]bool) (*nodes.ObjectName, string, []string, error) {
	if !p.isCreateTableNameStart() {
		return nil, "", nil, p.syntaxErrorAtCur()
	}
	nameStart := p.cur
	name, err := p.parseObjectName()
	if err != nil {
		return nil, "", nil, err
	}
	tableTok := p.prev
	if name.Schema == "" {
		tableTok = nameStart
	}
	var alias string
	if p.cur.Type == kwAS {
		p.advance()
		aliasTok := p.cur
		alias = aliasTok.Str
		if err := p.parseCreateTableName(); err != nil {
			return nil, "", nil, err
		}
		if seenTables[alias] {
			return nil, "", nil, p.syntaxErrorAtTok(aliasTok)
		}
		seenTables[alias] = true
	} else {
		if seenTables[name.Name] {
			return nil, "", nil, p.syntaxErrorAtTok(tableTok)
		}
		seenTables[name.Name] = true
	}
	var key []string
	if p.cur.Type == kwKEY {
		p.advance()
		if key, err = p.parseGraphNameListUnique(true); err != nil {
			return nil, "", nil, err
		}
	}
	return name, alias, key, nil
}

// parseGraphTableDef parses a vertex_table_definition.
//
//	graph_element_name_and_key [ graph_table_label_and_properties ]
func (p *Parser) parseGraphTableDef(seenTables map[string]bool) (*nodes.GraphTableDef, error) {
	def := &nodes.GraphTableDef{Loc: nodes.Loc{Start: p.pos()}}
	var err error
	if def.Name, def.Alias, def.KeyColumns, err = p.parseGraphElementNameAndKey(seenTables); err != nil {
		return nil, err
	}
	if err := p.parseGraphLabelAndProperties(def); err != nil {
		return nil, err
	}
	def.Loc.End = p.prev.End
	return def, nil
}

// parseGraphVertexReference reads the end of an edge's SOURCE or
// DESTINATION.
//
//	[ KEY ( column_name [, column_name ]... ) REFERENCES ] vertex_table_reference
//	vertex_table_reference: graph_element_name [ ( column_name [, column_name ]... ) ]
//
// Oracle 23ai takes a vertex table without its columns (SOURCE t), which
// its foreign key then supplies.
func (p *Parser) parseGraphVertexReference() ([]string, string, []string, error) {
	var key []string
	var err error
	if p.cur.Type == kwKEY {
		p.advance()
		if key, err = p.parseGraphNameListUnique(true); err != nil {
			return nil, "", nil, err
		}
		if err := p.expectCreateTableToken(kwREFERENCES); err != nil {
			return nil, "", nil, err
		}
	}
	ref := p.cur.Str
	if err := p.parseCreateTableName(); err != nil {
		return nil, "", nil, err
	}
	var cols []string
	if p.cur.Type == '(' {
		if cols, err = p.parseGraphNameListUnique(true); err != nil {
			return nil, "", nil, err
		}
	}
	return key, ref, cols, nil
}

// parseGraphEdgeDef parses an edge_table_definition.
//
//	graph_element_name_and_key
//	SOURCE [ KEY ( column_name [, column_name ]... ) REFERENCES ] vertex_table_reference
//	DESTINATION [ KEY ( column_name [, column_name ]... ) REFERENCES ] vertex_table_reference
//	[ graph_table_label_and_properties ]
func (p *Parser) parseGraphEdgeDef(seenTables map[string]bool) (*nodes.GraphEdgeDef, error) {
	def := &nodes.GraphEdgeDef{Loc: nodes.Loc{Start: p.pos()}}
	var err error
	if def.Name, def.Alias, def.KeyColumns, err = p.parseGraphElementNameAndKey(seenTables); err != nil {
		return nil, err
	}
	if err := p.expectCreateTableWord("SOURCE"); err != nil {
		return nil, err
	}
	if def.SourceKeyColumns, def.SourceRef, def.SourceRefColumns, err = p.parseGraphVertexReference(); err != nil {
		return nil, err
	}
	if err := p.expectCreateTableWord("DESTINATION"); err != nil {
		return nil, err
	}
	if def.DestKeyColumns, def.DestRef, def.DestRefColumns, err = p.parseGraphVertexReference(); err != nil {
		return nil, err
	}
	vtd := &nodes.GraphTableDef{}
	if err := p.parseGraphLabelAndProperties(vtd); err != nil {
		return nil, err
	}
	def.Labels = vtd.Labels
	def.Properties = vtd.Properties
	def.PropColumns = vtd.PropColumns
	def.PropAliases = vtd.PropAliases
	def.DefaultLabel = vtd.DefaultLabel
	def.Loc.End = p.prev.End
	return def, nil
}

// parseGraphLabelAndProperties parses graph_table_label_and_properties.
//
//	[ LABEL label_name [ LABEL label_name ]... ]
//	[ PROPERTIES [ ARE ] ALL COLUMNS [ EXCEPT ( column_name [, column_name ]... ) ]
//	| PROPERTIES ( { column_name [ AS property_name ] | value_expression AS property_name } [, ...] )
//	| NO PROPERTIES
//	]
//	[ DEFAULT LABEL ]
func (p *Parser) parseGraphLabelAndProperties(def *nodes.GraphTableDef) error {
	for p.isKeywordStr("LABEL") {
		p.advance()
		def.Labels = append(def.Labels, p.cur.Str)
		if err := p.parseCreateTableName(); err != nil {
			return err
		}
	}

	switch {
	case p.isKeywordStr("PROPERTIES"):
		p.advance()
		if p.cur.Type == '(' {
			def.Properties = "LIST"
			p.advance() // '('
			seenProps := map[string]bool{}
			for {
				start := p.cur
				if p.cur.Type == kwAS {
					return p.syntaxErrorAtCur()
				}
				expr, err := p.parseExpr()
				if err != nil {
					return err
				}
				if expr == nil {
					return p.syntaxErrorAtCur()
				}
				prop := p.source[start.Loc:p.prev.End]
				column := p.prev.Loc == start.Loc && start.Str != "" && start.Type != tokSCONST
				if column {
					prop = start.Str
				}
				def.PropColumns = append(def.PropColumns, prop)
				// A value expression names its property; only a column may
				// leave it out.
				if !column && p.cur.Type != kwAS {
					return p.syntaxErrorAtCur()
				}
				propTok := start
				propName := prop
				if p.cur.Type == kwAS {
					p.advance()
					propTok = p.cur
					propName = propTok.Str
					def.PropAliases = append(def.PropAliases, propName)
					if err := p.parseCreateTableName(); err != nil {
						return err
					}
				}
				if seenProps[propName] {
					return p.syntaxErrorAtTok(propTok)
				}
				seenProps[propName] = true
				if p.cur.Type != ',' {
					break
				}
				p.advance() // ','
			}
			if err := p.expectCreateTableToken(')'); err != nil {
				return err
			}
			break
		}
		if p.isKeywordStr("ARE") {
			p.advance()
		}
		if err := p.expectCreateTableToken(kwALL); err != nil {
			return err
		}
		if err := p.expectCreateTableToken(kwCOLUMNS); err != nil {
			return err
		}
		def.Properties = "ALL"
		if p.cur.Type == kwEXCEPT {
			p.advance()
			def.Properties = "ALL_EXCEPT"
			var err error
			if def.PropColumns, err = p.parseGraphNameListUnique(true); err != nil {
				return err
			}
		}
	case p.isKeywordStr("NO"):
		p.advance()
		if err := p.expectCreateTableWord("PROPERTIES"); err != nil {
			return err
		}
		def.Properties = "NO"
	}

	if p.cur.Type == kwDEFAULT {
		p.advance()
		if err := p.expectCreateTableWord("LABEL"); err != nil {
			return err
		}
		def.DefaultLabel = true
	}
	return nil
}
