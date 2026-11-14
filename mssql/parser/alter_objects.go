// Package parser - alter_objects.go implements T-SQL ALTER INDEX parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// onOffOptionValues is the ON | OFF value set.
var onOffOptionValues = newOptionSet().withIdents("ON", "OFF")

// indexOptions defines the valid option names for index WITH clauses.
// Derived from SqlScriptDOM IndexOptionHelper + IndexStateOption + spatial/JSON/vector options.
var indexOptions = newOptionSet(
	kwFILLFACTOR,     // FILLFACTOR = N
	kwMAXDOP,         // MAXDOP = N
	kwBUCKET_COUNT,   // BUCKET_COUNT = N (hash indexes)
	kwLOB_COMPACTION, // LOB_COMPACTION = ON|OFF
	kwORDER,          // ORDER (columnstore)
	kwTYPE,           // TYPE (vector indexes)
).withIdents(
	// indexStateOption — on/off options parsed via IndexOptionHelper
	"PAD_INDEX",
	"SORT_IN_TEMPDB",
	"STATISTICS_NORECOMPUTE",
	"STATISTICS_INCREMENTAL",
	"DROP_EXISTING",
	"ALLOW_ROW_LOCKS",
	"ALLOW_PAGE_LOCKS",
	"IGNORE_DUP_KEY",
	"ONLINE",
	"RESUMABLE",
	"COMPRESS_ALL_ROW_GROUPS",
	"OPTIMIZE_FOR_SEQUENTIAL_KEY",
	// expression-valued options
	"MAX_DURATION",
	"DATA_COMPRESSION",
	"XML_COMPRESSION",
	"COMPRESSION_DELAY",
	"WAIT_AT_LOW_PRIORITY",
	// spatial index options
	"BOUNDING_BOX",
	"GRIDS",
	"CELLS_PER_OBJECT",
	// JSON index options
	"OPTIMIZE_FOR_ARRAY_SEARCH",
	// vector index options
	"METRIC",
)

// parseAlterIndexOptions parses a parenthesized list of index options.
//
// <rebuild_index_option> / <reorganize_option> / <set_index_option> ::=
//
//	PAD_INDEX = { ON | OFF }
//	| FILLFACTOR = fillfactor
//	| SORT_IN_TEMPDB = { ON | OFF }
//	| IGNORE_DUP_KEY = { ON | OFF }
//	| STATISTICS_NORECOMPUTE = { ON | OFF }
//	| STATISTICS_INCREMENTAL = { ON | OFF }
//	| ONLINE = { ON [ ( <low_priority_lock_wait> ) ] | OFF }
//	| RESUMABLE = { ON | OFF }
//	| MAX_DURATION = <time> [ MINUTES ]
//	| ALLOW_ROW_LOCKS = { ON | OFF }
//	| ALLOW_PAGE_LOCKS = { ON | OFF }
//	| MAXDOP = max_degree_of_parallelism
//	| DATA_COMPRESSION = { NONE | ROW | PAGE | COLUMNSTORE | COLUMNSTORE_ARCHIVE }
//	    [ ON PARTITIONS ( { <partition_number> [ TO <partition_number> ] } [ , ...n ] ) ]
//	| XML_COMPRESSION = { ON | OFF }
//	    [ ON PARTITIONS ( { <partition_number> [ TO <partition_number> ] } [ , ...n ] ) ]
//	| LOB_COMPACTION = { ON | OFF }
//	| COMPRESS_ALL_ROW_GROUPS = { ON | OFF }
func (p *Parser) parseAlterIndexOptions() (*nodes.List, error) {
	p.advance() // consume '('

	opts := &nodes.List{}
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		if p.cur.Type == ',' {
			p.advance()
			continue
		}

		// Parse option_name = value (validated against known index options)
		if p.isValidOption(indexOptions) {
			name := strings.ToUpper(p.cur.Str)
			p.advance()

			if p.cur.Type == '=' {
				p.advance() // consume '='
				// Value can be ON, OFF, an identifier, or a number
				val := ""
				if p.cur.Type == kwON {
					val = "ON"
					p.advance()
					// ONLINE = ON ( <low_priority_lock_wait> ) -- skip nested parens
					if p.cur.Type == '(' {
						depth := 1
						p.advance() // consume '('
						for depth > 0 && p.cur.Type != tokEOF {
							if p.cur.Type == '(' {
								depth++
							} else if p.cur.Type == ')' {
								depth--
								if depth == 0 {
									p.advance()
									break
								}
							}
							p.advance()
						}
					}
				} else if p.cur.Type == kwOFF {
					val = "OFF"
					p.advance()
				} else if p.cur.Type == tokICONST || p.cur.Type == tokFCONST {
					val = p.cur.Str
					p.advance()
				} else if p.isIdentLike() {
					val = strings.ToUpper(p.cur.Str)
					p.advance()
				} else if p.cur.Type == '(' {
					// Nested parenthesized value like BOUNDING_BOX = (0, 0, 100, 100)
					val = p.consumeBalancedValueText()
				}
				opts.Items = append(opts.Items, &nodes.String{Str: name + "=" + val})
			} else {
				opts.Items = append(opts.Items, &nodes.String{Str: name})
			}

			// DATA_COMPRESSION / XML_COMPRESSION may have ON PARTITIONS (...)
			if p.cur.Type == kwON {
				next := p.peekNext()
				if next.Type == kwPARTITIONS {
					p.advance() // consume ON
					p.advance() // consume PARTITIONS
					if p.cur.Type == '(' {
						depth := 1
						p.advance() // consume '('
						for depth > 0 && p.cur.Type != tokEOF {
							if p.cur.Type == '(' {
								depth++
							} else if p.cur.Type == ')' {
								depth--
								if depth == 0 {
									p.advance()
									break
								}
							}
							p.advance()
						}
					}
				}
			}

			// MAX_DURATION may have MINUTES suffix
			if strings.HasPrefix(name, "MAX_DURATION") && p.cur.Type == kwMINUTES {
				p.advance()
			}
		} else if p.cur.Type == tokSCONST || p.cur.Type == tokICONST || p.cur.Type == tokFCONST {
			// Collect string/numeric literals as standalone option values
			opts.Items = append(opts.Items, &nodes.String{Str: p.cur.Str})
			p.advance()
		} else {
			// Unknown token in option list — reject
			return nil, p.unexpectedToken()
		}
	}

	p.match(')') // consume ')'
	return opts, nil
}

func (p *Parser) consumeBalancedValueText() string {
	start := p.pos()
	depth := 0
	for p.cur.Type != tokEOF {
		if p.cur.Type == '(' {
			depth++
		} else if p.cur.Type == ')' {
			depth--
			p.advance()
			if depth == 0 {
				return strings.TrimSpace(p.source[start:p.prevEnd()])
			}
			continue
		}
		p.advance()
	}
	return strings.TrimSpace(p.source[start:p.prevEnd()])
}
