package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// This file reads the property clauses of ALTER TABLE that other clauses
// share: supplemental logging, column_properties and the clustering
// options (oracle/parser/bnf/ALTER-TABLE.bnf).

// parseATSupplementalLogging reads supplemental_table_logging after
// { ADD | DROP } SUPPLEMENTAL:
//
//	ADD SUPPLEMENTAL LOG { supplemental_log_grp_clause | supplemental_id_key_clause } [, ...]
//	DROP SUPPLEMENTAL LOG { supplemental_id_key_clause | GROUP log_group } [, ...]
//	supplemental_log_grp_clause:
//	    GROUP log_group ( column [ NO LOG ] [, column [ NO LOG ] ]... ) [ ALWAYS ]
//	supplemental_id_key_clause:
//	    DATA ( { ALL | PRIMARY KEY | UNIQUE | FOREIGN KEY } [, ...] ) COLUMNS
func (p *Parser) parseATSupplementalLogging(add bool, opts *nodes.List) error {
	if err := p.atExpectWord("LOG"); err != nil {
		return err
	}
	for {
		start := p.pos()
		switch {
		case p.atWord("DATA"):
			p.advance()
			if err := p.atExpect('('); err != nil {
				return err
			}
			// Each kind of key stands once (ORA-32586).
			seen := map[int]bool{}
			for {
				if seen[p.cur.Type] {
					return p.syntaxErrorAtCur()
				}
				seen[p.cur.Type] = true
				switch {
				case p.cur.Type == kwALL, p.cur.Type == kwUNIQUE:
					p.advance()
				case p.cur.Type == kwPRIMARY, p.cur.Type == kwFOREIGN:
					p.advance()
					if err := p.atExpect(kwKEY); err != nil {
						return err
					}
				default:
					return p.syntaxErrorAtCur()
				}
				if p.cur.Type != ',' {
					break
				}
				p.advance()
			}
			// ALTER TABLE's page puts COLUMNS after DATA's list; the
			// object_properties of a nested table's storage table puts it
			// inside the parentheses, as its own converted rule does.
			insideColumns := p.cur.Type == kwCOLUMNS
			if insideColumns {
				p.advance()
			}
			if err := p.atExpect(')'); err != nil {
				return err
			}
			p.atOption(opts, start)
			start = p.pos()
			if !insideColumns {
				if err := p.atExpect(kwCOLUMNS); err != nil {
					return err
				}
			}
		case p.cur.Type == kwGROUP:
			p.advance()
			if err := p.parseATName(); err != nil {
				return err
			}
			if add {
				if err := p.atExpect('('); err != nil {
					return err
				}
				for {
					if err := p.parseATName(); err != nil {
						return err
					}
					if p.atWord("NO") {
						p.advance()
						if err := p.atExpectWord("LOG"); err != nil {
							return err
						}
					}
					if p.cur.Type != ',' {
						break
					}
					p.advance()
				}
				if err := p.atExpect(')'); err != nil {
					return err
				}
				if p.cur.Type == kwALWAYS {
					p.advance()
				}
			}
		default:
			return p.syntaxErrorAtCur()
		}
		p.atOption(opts, start)
		if p.cur.Type != ',' {
			return nil
		}
		p.advance()
	}
}

// parseATSubstitutable reads substitutable_column_clause if one follows,
// setting *found, where found is not nil, to whether one did:
//
//	{ [ ELEMENT ] IS OF [ TYPE ] ( [ ONLY ] type ) | [ NOT ] SUBSTITUTABLE AT ALL LEVELS }
func (p *Parser) parseATSubstitutable(found *bool) error {
	switch {
	case p.atWord("ELEMENT"), p.cur.Type == kwIS:
		if p.atWord("ELEMENT") {
			p.advance()
		}
		if err := p.atExpect(kwIS); err != nil {
			return err
		}
		if err := p.atExpect(kwOF); err != nil {
			return err
		}
		if p.cur.Type == kwTYPE {
			p.advance()
		}
		if err := p.atExpect('('); err != nil {
			return err
		}
		p.atOptOneOf("ONLY")
		if err := p.parseATQualifiedName(); err != nil {
			return err
		}
		if err := p.atExpect(')'); err != nil {
			return err
		}
	case p.cur.Type == kwNOT, p.atWord("SUBSTITUTABLE"):
		if p.cur.Type == kwNOT {
			p.advance()
		}
		if err := p.atExpectWord("SUBSTITUTABLE"); err != nil {
			return err
		}
		if err := p.atExpectWord("AT"); err != nil {
			return err
		}
		if err := p.atExpect(kwALL); err != nil {
			return err
		}
		if err := p.atExpectWord("LEVELS"); err != nil {
			return err
		}
	default:
		return nil
	}
	if found != nil {
		*found = true
	}
	return nil
}

// parseATColumnProperties reads [ column_properties ]:
//
//	{ object_type_col_properties | nested_table_col_properties
//	| { varray_col_properties | LOB_storage_clause }
//	| XMLType_column_properties | json_storage_clause }...
//
// Each property is recorded as an option keyed by its first word. Oracle
// refuses object_type_col_properties (COLUMN column substitutable_column_clause)
// in ALTER TABLE, at COLUMN.
func (p *Parser) parseATColumnProperties(opts *nodes.List) error {
	for {
		start := p.pos()
		var err error
		switch {
		case p.cur.Type == kwNESTED:
			err = p.parseATNestedTableStorage()
		case p.isATLobStorageStart():
			if err = p.parseATLobStorage(); err == nil && p.cur.Type == '(' && p.peekNext().Type == kwPARTITION {
				err = p.parseATLobPartitionStorage()
			}
		case p.atWord("XMLTYPE"):
			err = p.parseATXMLTypeColumnProperties()
		case p.cur.Type == kwJSON && p.peekNext().Type == '(':
			// json_storage_clause: JSON ( column [, column ]... ) STORE AS
			// [ lob_segname ] ( LOB_parameters )
			p.advance()
			if err = p.parseATNameList(); err == nil {
				err = p.parseATLobStoreAs()
			}
		default:
			return nil
		}
		if err != nil {
			return err
		}
		p.atOption(opts, start)
	}
}

// parseATXMLTypeColumnProperties reads XMLType_column_properties:
//
//	XMLTYPE [ COLUMN ] column [ XMLType_storage ] [ XMLSchema_spec ]
//	XMLType_storage:
//	    STORE { AS { OBJECT RELATIONAL
//	               | [ SECUREFILE | BASICFILE ] { CLOB | BINARY XML }
//	                   [ { LOB_segname [ ( LOB_parameters ) ] | ( LOB_parameters ) } ] }
//	          | ALL VARRAYS AS { LOBS | TABLES } }
//	XMLSchema_spec:
//	    [ XMLSCHEMA XMLSchema_URL ] ELEMENT { element | XMLSchema_URL # element }
//	    [ STORE ALL VARRAYS AS { LOBS | TABLES } ]
//
// Oracle refuses the page's [ { ALLOW | DISALLOW } { NONSCHEMA | ANYSCHEMA } ]
// after the element here (ORA-02000).
func (p *Parser) parseATXMLTypeColumnProperties() error {
	p.advance() // XMLTYPE
	if p.cur.Type == kwCOLUMN {
		p.advance()
	}
	if err := p.parseATName(); err != nil {
		return err
	}
	storeAllVarrays := func() error {
		p.advance() // ALL
		if err := p.atExpectWord("VARRAYS"); err != nil {
			return err
		}
		if err := p.atExpect(kwAS); err != nil {
			return err
		}
		return p.atOneOf("LOBS", "TABLES")
	}
	if p.atWord("STORE") {
		p.advance()
		if p.cur.Type == kwALL {
			if err := storeAllVarrays(); err != nil {
				return err
			}
		} else {
			if err := p.atExpect(kwAS); err != nil {
				return err
			}
			if p.atWord("OBJECT") {
				p.advance()
				if err := p.atExpectWord("RELATIONAL"); err != nil {
					return err
				}
			} else {
				p.atOptOneOf("SECUREFILE", "BASICFILE")
				clob := p.cur.Type == kwCLOB
				switch {
				case clob:
					p.advance()
				case p.atWord("BINARY"):
					p.advance()
					if err := p.atExpectWord("XML"); err != nil {
						return err
					}
				default:
					return p.syntaxErrorAtCur()
				}
				// Oracle takes a LOB_segname after CLOB only (ORA-02000 after
				// BINARY XML).
				if clob && p.cur.Type != '(' && p.isATNameStart() && !p.atWord("XMLSCHEMA") && !p.atWord("ELEMENT") {
					p.advance()
				}
				if p.cur.Type == '(' {
					if err := p.parseATLobStorageParameters(); err != nil {
						return err
					}
				}
			}
		}
	}
	if p.atWord("XMLSCHEMA") || p.atWord("ELEMENT") {
		if p.atWord("XMLSCHEMA") {
			p.advance()
			// Oracle refuses a string here (ORA-02000 at XMLSCHEMA 's').
			if err := p.parseATName(); err != nil {
				return err
			}
		}
		if err := p.atExpectWord("ELEMENT"); err != nil {
			return err
		}
		if p.cur.Type == tokSCONST {
			p.advance()
		} else if err := p.parseATName(); err != nil {
			return err
		}
		if p.atWord("STORE") {
			p.advance()
			if p.cur.Type != kwALL {
				return p.syntaxErrorAtCur()
			}
			if err := storeAllVarrays(); err != nil {
				return err
			}
		}
	}
	return nil
}

// parseATClusteringWhen reads what follows MODIFY CLUSTERING:
//
//	[ clustering_when ] [ zonemap_clause ]
//	clustering_when:
//	    [ { YES | NO } ON LOAD ] [ { YES | NO } ON DATA MOVEMENT ]
//	zonemap_clause:
//	    { WITH MATERIALIZED ZONEMAP [ ( zonemap_name ) ] | WITHOUT MATERIALIZED ZONEMAP }
//
// The page gives MODIFY CLUSTERING attribute_clustering_clause, which
// Oracle refuses (ORA-03049 at the second CLUSTERING); the clause above is
// the one of ALTER TABLE's own clustering section.
func (p *Parser) parseATClusteringWhen(opts *nodes.List) error {
	yesNoOn := func() bool {
		return (p.atWord("YES") || p.atWord("NO")) && p.peekNext().Type == kwON
	}
	if yesNoOn() && p.isIdentLikeStrAt(p.peekAhead(2)[1], "LOAD") {
		start := p.pos()
		p.advance()
		p.advance()
		p.advance()
		p.atOption(opts, start)
	}
	if yesNoOn() {
		start := p.pos()
		p.advance()
		p.advance()
		if err := p.atExpectWord("DATA"); err != nil {
			return err
		}
		if err := p.atExpectWord("MOVEMENT"); err != nil {
			return err
		}
		p.atOption(opts, start)
	}
	start := p.pos()
	switch {
	case p.cur.Type == kwWITH:
		p.advance()
		if err := p.atExpectWord("MATERIALIZED"); err != nil {
			return err
		}
		if err := p.atExpectWord("ZONEMAP"); err != nil {
			return err
		}
		if p.cur.Type == '(' {
			p.advance()
			if err := p.parseATName(); err != nil {
				return err
			}
			if err := p.atExpect(')'); err != nil {
				return err
			}
		}
	case p.atWord("WITHOUT"):
		p.advance()
		if err := p.atExpectWord("MATERIALIZED"); err != nil {
			return err
		}
		if err := p.atExpectWord("ZONEMAP"); err != nil {
			return err
		}
	default:
		return nil
	}
	p.atOption(opts, start)
	return nil
}
