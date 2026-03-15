package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseCreateIndexStmt parses a CREATE INDEX statement.
// The CREATE keyword has already been consumed. The current token is
// UNIQUE, BITMAP, MULTIVALUE, or INDEX.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-INDEX.html
//
//	CREATE [ UNIQUE ] [ BITMAP ] [ MULTIVALUE ] INDEX [ IF NOT EXISTS ]
//	    [ schema. ] index_name
//	    [ index_ilm_clause ]
//	    { cluster_index_clause | table_index_clause | bitmap_join_index_clause }
//	    [ { DEFERRED | IMMEDIATE } INVALIDATION ]
func (p *Parser) parseCreateIndexStmt(start int) *nodes.CreateIndexStmt {
	stmt := &nodes.CreateIndexStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Optional UNIQUE
	if p.cur.Type == kwUNIQUE {
		stmt.Unique = true
		p.advance()
	}

	// Optional BITMAP
	if p.cur.Type == kwBITMAP {
		stmt.Bitmap = true
		p.advance()
	}

	// Optional MULTIVALUE
	if p.isIdentLikeStr("MULTIVALUE") {
		stmt.Multivalue = true
		p.advance()
	}

	// INDEX keyword
	if p.cur.Type == kwINDEX {
		p.advance()
	}

	// Optional IF NOT EXISTS
	if p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwNOT {
			stmt.IfNotExists = true
			p.advance() // consume IF
			p.advance() // consume NOT
			if p.cur.Type == kwEXISTS {
				p.advance() // consume EXISTS
			}
		}
	}

	// Index name
	stmt.Name = p.parseObjectName()

	// ON clause
	if p.cur.Type == kwON {
		p.advance() // consume ON

		// cluster_index_clause: ON CLUSTER [schema.]cluster
		if p.cur.Type == kwCLUSTER {
			p.advance() // consume CLUSTER
			stmt.Cluster = p.parseObjectName()
			// Parse index_attributes
			p.parseCreateIndexAttributes(stmt)
			// Parse optional INVALIDATION
			p.parseCreateIndexInvalidation(stmt)
			stmt.Loc.End = p.pos()
			return stmt
		}

		// table_index_clause or bitmap_join_index_clause
		// Both start with [schema.]table
		stmt.Table = p.parseObjectName()

		// Optional table alias (single identifier before '(')
		if p.isIdentLike() && p.cur.Type != '(' && !p.isKeywordLike() {
			stmt.Alias = p.parseIdentifier()
		}
	}

	// ( column_list )
	if p.cur.Type == '(' {
		p.advance()
		stmt.Columns = p.parseIndexColumnList()
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	// Check for bitmap_join_index_clause: FROM ... WHERE ...
	if p.cur.Type == kwFROM {
		p.advance() // consume FROM
		stmt.FromTables = &nodes.List{}
		// Parse table list
		tbl := p.parseObjectName()
		stmt.FromTables.Items = append(stmt.FromTables.Items, tbl)
		for p.cur.Type == ',' {
			p.advance() // consume ','
			tbl = p.parseObjectName()
			stmt.FromTables.Items = append(stmt.FromTables.Items, tbl)
		}
		// WHERE join_condition
		if p.cur.Type == kwWHERE {
			p.advance() // consume WHERE
			stmt.Where = p.parseExpr()
		}
	}

	// Parse trailing options (index_properties + index_attributes)
	p.parseCreateIndexAttributes(stmt)

	// Parse optional INVALIDATION
	p.parseCreateIndexInvalidation(stmt)

	stmt.Loc.End = p.pos()
	return stmt
}

// parseCreateIndexInvalidation parses optional { DEFERRED | IMMEDIATE } INVALIDATION.
func (p *Parser) parseCreateIndexInvalidation(stmt *nodes.CreateIndexStmt) {
	if p.cur.Type == kwDEFERRED {
		stmt.Invalidation = "DEFERRED"
		p.advance()
		if p.isIdentLikeStr("INVALIDATION") {
			p.advance()
		}
	} else if p.cur.Type == kwIMMEDIATE {
		stmt.Invalidation = "IMMEDIATE"
		p.advance()
		if p.isIdentLikeStr("INVALIDATION") {
			p.advance()
		}
	}
}

// isKeywordLike checks if current token looks like a known clause-starting keyword
// used to distinguish table aliases from trailing clauses.
func (p *Parser) isKeywordLike() bool {
	switch p.cur.Type {
	case kwREVERSE, kwTABLESPACE, kwLOCAL, kwGLOBAL, kwONLINE,
		kwPARALLEL, kwNOPARALLEL, kwCOMPRESS, kwNOCOMPRESS,
		kwLOGGING, kwNOLOGGING, kwINVISIBLE, kwPCTFREE,
		kwFROM, kwWHERE, kwPARTITION, kwDEFERRED, kwIMMEDIATE:
		return true
	}
	if p.isIdentLikeStr("SORT") || p.isIdentLikeStr("NOSORT") ||
		p.isIdentLikeStr("VISIBLE") || p.isIdentLikeStr("INDEXTYPE") ||
		p.isIdentLikeStr("INDEXING") || p.isIdentLikeStr("PARAMETERS") {
		return true
	}
	return false
}

// parseCreateIndexAttributes parses index_attributes and index_properties for CREATE INDEX.
//
//	index_attributes:
//	    [ physical_attributes_clause ]
//	    [ TABLESPACE { tablespace | DEFAULT } ]
//	    [ index_compression ]
//	    [ partial_index_clause ]
//	    [ parallel_clause ]
//	    [ { SORT | NOSORT } ]
//	    [ REVERSE ]
//	    [ { VISIBLE | INVISIBLE } ]
//	    [ logging_clause ]
//	    [ ONLINE ]
//
//	index_properties:
//	    [ global_partitioned_index | local_partitioned_index ]
//	    [ index_attributes ]
//	    [ domain_index_clause | XMLIndex_clause ]
func (p *Parser) parseCreateIndexAttributes(stmt *nodes.CreateIndexStmt) {
	for {
		switch p.cur.Type {
		case kwREVERSE:
			stmt.Reverse = true
			p.advance()
		case kwTABLESPACE:
			p.advance()
			stmt.Tablespace = p.parseIdentifier()
		case kwLOCAL:
			stmt.Local = true
			p.advance()
			// Skip local partitioned details (parenthesized partition list)
			p.skipParenthesizedBlock()
		case kwGLOBAL:
			stmt.Global = true
			p.advance()
			// GLOBAL PARTITION BY ... skip details
			if p.cur.Type == kwPARTITION {
				p.advance() // consume PARTITION
				if p.cur.Type == kwBY {
					p.advance() // consume BY
				}
				// Skip RANGE/HASH and partitioning details
				p.skipToNextClause()
			}
		case kwONLINE:
			stmt.Online = true
			p.advance()
		case kwPARALLEL:
			p.advance()
			if p.cur.Type == tokICONST {
				stmt.Parallel = p.cur.Str
				p.advance()
			} else {
				stmt.Parallel = "PARALLEL"
			}
		case kwNOPARALLEL:
			stmt.Parallel = "NOPARALLEL"
			p.advance()
		case kwCOMPRESS:
			p.advance()
			// COMPRESS ADVANCED [LOW | HIGH] or COMPRESS [integer]
			if p.isIdentLikeStr("ADVANCED") {
				p.advance() // consume ADVANCED
				if p.isIdentLikeStr("LOW") {
					stmt.Compress = "ADVANCED LOW"
					p.advance()
				} else if p.isIdentLikeStr("HIGH") {
					stmt.Compress = "ADVANCED HIGH"
					p.advance()
				} else {
					stmt.Compress = "ADVANCED"
				}
			} else if p.cur.Type == tokICONST {
				stmt.Compress = p.cur.Str
				p.advance()
			} else {
				stmt.Compress = "COMPRESS"
			}
		case kwNOCOMPRESS:
			stmt.Compress = "NOCOMPRESS"
			p.advance()
		case kwLOGGING:
			stmt.Logging = true
			p.advance()
		case kwNOLOGGING:
			stmt.NoLogging = true
			p.advance()
		case kwINVISIBLE:
			stmt.Invisible = true
			p.advance()
		case kwPCTFREE:
			p.advance()
			if p.cur.Type == tokICONST {
				stmt.PctFree = p.cur.Str
				p.advance()
			}
		default:
			// Check identifier-based keywords
			if p.isIdentLikeStr("VISIBLE") {
				stmt.Visible = true
				p.advance()
			} else if p.isIdentLikeStr("SORT") {
				stmt.Sort = true
				p.advance()
			} else if p.isIdentLikeStr("NOSORT") {
				stmt.NoSort = true
				p.advance()
			} else if p.isIdentLikeStr("INITRANS") {
				p.advance()
				if p.cur.Type == tokICONST {
					stmt.InitTrans = p.cur.Str
					p.advance()
				}
			} else if p.isIdentLikeStr("MAXTRANS") {
				p.advance()
				if p.cur.Type == tokICONST {
					stmt.MaxTrans = p.cur.Str
					p.advance()
				}
			} else if p.isIdentLikeStr("INDEXTYPE") {
				// domain_index_clause: INDEXTYPE IS [schema.]indextype
				p.advance() // consume INDEXTYPE
				if p.cur.Type == kwIS {
					p.advance() // consume IS
				}
				stmt.IndexType = p.parseObjectName()
				// Optional LOCAL and PARAMETERS
				if p.cur.Type == kwLOCAL {
					stmt.Local = true
					p.advance()
					p.skipParenthesizedBlock()
				}
				// Optional parallel_clause
				if p.cur.Type == kwPARALLEL {
					p.advance()
					if p.cur.Type == tokICONST {
						stmt.Parallel = p.cur.Str
						p.advance()
					} else {
						stmt.Parallel = "PARALLEL"
					}
				} else if p.cur.Type == kwNOPARALLEL {
					stmt.Parallel = "NOPARALLEL"
					p.advance()
				}
				// Optional PARAMETERS
				if p.isIdentLikeStr("PARAMETERS") {
					p.advance() // consume PARAMETERS
					if p.cur.Type == '(' {
						p.advance()
						if p.cur.Type == tokSCONST {
							stmt.Parameters = p.cur.Str
							p.advance()
						}
						if p.cur.Type == ')' {
							p.advance()
						}
					}
				}
			} else if p.isIdentLikeStr("INDEXING") {
				p.advance() // consume INDEXING
				if p.isIdentLikeStr("FULL") {
					stmt.IndexingFull = true
					p.advance()
				} else if p.isIdentLikeStr("PARTIAL") {
					stmt.IndexingPartial = true
					p.advance()
				}
			} else if p.isIdentLikeStr("PARAMETERS") {
				p.advance() // consume PARAMETERS
				if p.cur.Type == '(' {
					p.advance()
					if p.cur.Type == tokSCONST {
						stmt.Parameters = p.cur.Str
						p.advance()
					}
					if p.cur.Type == ')' {
						p.advance()
					}
				}
			} else if p.isIdentLikeStr("STORAGE") {
				// storage_clause - skip parenthesized block
				p.advance()
				p.skipParenthesizedBlock()
			} else if p.isIdentLikeStr("ANNOTATIONS") {
				p.advance()
				p.skipParenthesizedBlock()
			} else {
				return
			}
		}
	}
}

// skipParenthesizedBlock skips a parenthesized block if present.
func (p *Parser) skipParenthesizedBlock() {
	if p.cur.Type != '(' {
		return
	}
	depth := 1
	p.advance() // consume '('
	for depth > 0 && p.cur.Type != tokEOF {
		if p.cur.Type == '(' {
			depth++
		} else if p.cur.Type == ')' {
			depth--
		}
		p.advance()
	}
}

// skipToNextClause skips tokens until we reach a known clause keyword,
// semicolon, or EOF. This is used for complex partitioning clauses.
func (p *Parser) skipToNextClause() {
	for p.cur.Type != ';' && p.cur.Type != tokEOF {
		switch p.cur.Type {
		case kwDEFERRED, kwIMMEDIATE:
			return
		case '(':
			p.skipParenthesizedBlock()
		default:
			// Check if we hit a clause-ending keyword
			if p.isIdentLikeStr("INVALIDATION") || p.isIdentLikeStr("STORE") {
				return
			}
			p.advance()
		}
	}
}

// parseIndexColumnList parses a comma-separated list of index columns.
func (p *Parser) parseIndexColumnList() *nodes.List {
	list := &nodes.List{}
	for {
		col := p.parseIndexColumn()
		if col == nil {
			break
		}
		list.Items = append(list.Items, col)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return list
}

// parseIndexColumn parses a single index column expression with optional ASC/DESC.
//
//	index_expr [ ASC | DESC ]
func (p *Parser) parseIndexColumn() *nodes.IndexColumn {
	start := p.pos()
	expr := p.parseExpr()
	if expr == nil {
		return nil
	}

	col := &nodes.IndexColumn{
		Expr: expr,
		Loc:  nodes.Loc{Start: start},
	}

	// ASC | DESC
	switch p.cur.Type {
	case kwASC:
		col.Dir = nodes.SORTBY_ASC
		p.advance()
	case kwDESC:
		col.Dir = nodes.SORTBY_DESC
		p.advance()
	}

	// NULLS FIRST | NULLS LAST
	if p.cur.Type == kwNULLS {
		p.advance()
		switch p.cur.Type {
		case kwFIRST:
			col.NullOrder = nodes.SORTBY_NULLS_FIRST
			p.advance()
		case kwLAST:
			col.NullOrder = nodes.SORTBY_NULLS_LAST
			p.advance()
		}
	}

	col.Loc.End = p.pos()
	return col
}

// parseIndexOptions is a legacy wrapper kept for backward compatibility.
// It delegates to parseCreateIndexAttributes.
func (p *Parser) parseIndexOptions(stmt *nodes.CreateIndexStmt) {
	p.parseCreateIndexAttributes(stmt)
}

// ---------------------------------------------------------------------------
// CREATE INDEXTYPE
// ---------------------------------------------------------------------------

// parseCreateIndextypeStmt parses a CREATE INDEXTYPE statement.
// Called after CREATE [OR REPLACE] has been consumed and INDEXTYPE consumed.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-INDEXTYPE.html
//
//	CREATE [ OR REPLACE ] INDEXTYPE [ IF NOT EXISTS ] [ schema. ] indextype
//	    [ SHARING = { METADATA | NONE } ]
//	    FOR [ schema. ] operator ( parameter_type [, parameter_type ]... )
//	        [, [ schema. ] operator ( parameter_type [, parameter_type ]... ) ]...
//	    using_type_clause
//	    [ WITH LOCAL [ RANGE ] PARTITION ]
//	    [ storage_table_clause ]
//	    [ array_DML_clause ]
func (p *Parser) parseCreateIndextypeStmt(start int, orReplace bool) *nodes.CreateIndextypeStmt {
	stmt := &nodes.CreateIndextypeStmt{
		OrReplace: orReplace,
		Loc:       nodes.Loc{Start: start},
	}

	// Optional IF NOT EXISTS
	if p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwNOT {
			stmt.IfNotExists = true
			p.advance() // consume IF
			p.advance() // consume NOT
			if p.cur.Type == kwEXISTS {
				p.advance() // consume EXISTS
			}
		}
	}

	// Indextype name
	stmt.Name = p.parseObjectName()

	// Optional SHARING = { METADATA | NONE }
	if p.isIdentLikeStr("SHARING") {
		p.advance() // consume SHARING
		if p.cur.Type == '=' {
			p.advance() // consume =
		}
		if p.isIdentLikeStr("METADATA") {
			stmt.Sharing = "METADATA"
			p.advance()
		} else if p.isIdentLikeStr("NONE") {
			stmt.Sharing = "NONE"
			p.advance()
		} else if p.isIdentLike() {
			stmt.Sharing = p.cur.Str
			p.advance()
		}
	}

	// FOR operator_list
	if p.cur.Type == kwFOR {
		p.advance() // consume FOR
		stmt.Operators = p.parseIndextypeOperatorList()
	}

	// USING implementation_type
	if p.cur.Type == kwUSING {
		p.advance() // consume USING
		stmt.UsingType = p.parseObjectName()
	}

	// Optional WITH LOCAL [RANGE] PARTITION
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		if p.cur.Type == kwLOCAL {
			stmt.WithLocal = true
			p.advance() // consume LOCAL
			if p.isIdentLikeStr("RANGE") {
				stmt.WithRange = true
				p.advance() // consume RANGE
			}
			if p.cur.Type == kwPARTITION {
				p.advance() // consume PARTITION
			}
		} else if p.isIdentLikeStr("SYSTEM") || p.cur.Type == kwUSER {
			// storage_table_clause
			stmt.StorageTable = p.parseIdentifier()
			// MANAGED STORAGE TABLES
			if p.isIdentLikeStr("MANAGED") {
				p.advance()
			}
			if p.isIdentLikeStr("STORAGE") {
				p.advance()
			}
			if p.isIdentLikeStr("TABLES") {
				p.advance()
			}
		} else if p.isIdentLikeStr("ARRAY") {
			// array_DML_clause
			stmt.ArrayDML = true
			p.advance() // consume ARRAY
			if p.isIdentLikeStr("DML") {
				p.advance() // consume DML
			}
			// Optional type list
			p.skipParenthesizedBlock()
		}
	}

	// More optional WITH clauses
	for p.cur.Type == kwWITH {
		p.advance() // consume WITH
		if p.cur.Type == kwLOCAL {
			stmt.WithLocal = true
			p.advance()
			if p.isIdentLikeStr("RANGE") {
				stmt.WithRange = true
				p.advance()
			}
			if p.cur.Type == kwPARTITION {
				p.advance()
			}
		} else if p.isIdentLikeStr("SYSTEM") || p.cur.Type == kwUSER {
			stmt.StorageTable = p.parseIdentifier()
			if p.isIdentLikeStr("MANAGED") {
				p.advance()
			}
			if p.isIdentLikeStr("STORAGE") {
				p.advance()
			}
			if p.isIdentLikeStr("TABLES") {
				p.advance()
			}
		} else if p.isIdentLikeStr("ARRAY") {
			stmt.ArrayDML = true
			p.advance()
			if p.isIdentLikeStr("DML") {
				p.advance()
			}
			p.skipParenthesizedBlock()
		} else {
			break
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseIndextypeOperatorList parses operator list for CREATE INDEXTYPE.
//
//	[ schema. ] operator ( parameter_type [, parameter_type ]... )
//	[, [ schema. ] operator ( parameter_type [, parameter_type ]... ) ]...
func (p *Parser) parseIndextypeOperatorList() []*nodes.IndextypeOp {
	var ops []*nodes.IndextypeOp
	for {
		op := p.parseIndextypeOp()
		if op == nil {
			break
		}
		ops = append(ops, op)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return ops
}

// parseIndextypeOp parses a single operator with parameter types.
func (p *Parser) parseIndextypeOp() *nodes.IndextypeOp {
	if !p.isIdentLike() {
		return nil
	}
	start := p.pos()
	name := p.parseObjectName()
	op := &nodes.IndextypeOp{
		Name: name,
		Loc:  nodes.Loc{Start: start},
	}
	// ( parameter_type [, parameter_type ]... )
	if p.cur.Type == '(' {
		p.advance()
		op.ParamTypes = p.parseTypeNameList()
		if p.cur.Type == ')' {
			p.advance()
		}
	}
	op.Loc.End = p.pos()
	return op
}

// parseTypeNameList parses a comma-separated list of type names.
func (p *Parser) parseTypeNameList() []string {
	var types []string
	for {
		tn := p.parseTypeNameStr()
		if tn == "" {
			break
		}
		types = append(types, tn)
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return types
}

// parseTypeName parses a single type name (e.g. NUMBER, VARCHAR2, schema.type_name).
func (p *Parser) parseTypeNameStr() string {
	if !p.isIdentLike() {
		return ""
	}
	name := p.parseIdentifier()
	// Handle schema.type_name
	if p.cur.Type == '.' {
		p.advance()
		name += "." + p.parseIdentifier()
	}
	// Handle parameterized types like VARCHAR2(100)
	if p.cur.Type == '(' {
		name += "("
		p.advance()
		for p.cur.Type != ')' && p.cur.Type != tokEOF {
			if p.cur.Type == tokICONST {
				name += p.cur.Str
			} else if p.isIdentLike() {
				name += p.cur.Str
			} else if p.cur.Type == ',' {
				name += ","
			}
			p.advance()
		}
		if p.cur.Type == ')' {
			name += ")"
			p.advance()
		}
	}
	return name
}

// ---------------------------------------------------------------------------
// ALTER INDEXTYPE
// ---------------------------------------------------------------------------

// parseAlterIndextypeStmt parses an ALTER INDEXTYPE statement.
// Called after ALTER has been consumed and INDEXTYPE consumed.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-INDEXTYPE.html
//
//	ALTER INDEXTYPE [ IF EXISTS ] [ schema. ] indextype
//	    { { ADD | DROP } [ schema. ] operator ( parameter_type [, parameter_type ]... )
//	        [ , { ADD | DROP } [ schema. ] operator ( parameter_type [, parameter_type ]... ) ]...
//	        [ using_type_clause ]
//	    | COMPILE
//	    }
//	    [ WITH LOCAL PARTITION ]
//	    [ storage_table_clause ]
func (p *Parser) parseAlterIndextypeStmt(start int) *nodes.AlterIndextypeStmt {
	stmt := &nodes.AlterIndextypeStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Optional IF EXISTS
	if p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwEXISTS {
			stmt.IfExists = true
			p.advance() // consume IF
			p.advance() // consume EXISTS
		}
	}

	// Indextype name
	stmt.Name = p.parseObjectName()

	// Action
	if p.isIdentLikeStr("COMPILE") {
		stmt.Action = "COMPILE"
		p.advance()
	} else if p.cur.Type == kwADD || p.cur.Type == kwDROP {
		stmt.Action = "ADD_DROP"
		// Parse modifications list
		for p.cur.Type == kwADD || p.cur.Type == kwDROP {
			modStart := p.pos()
			isAdd := p.cur.Type == kwADD
			p.advance() // consume ADD/DROP
			name := p.parseObjectName()
			mod := &nodes.IndextypeModOp{
				Add:  isAdd,
				Name: name,
				Loc:  nodes.Loc{Start: modStart},
			}
			if p.cur.Type == '(' {
				p.advance()
				mod.ParamTypes = p.parseTypeNameList()
				if p.cur.Type == ')' {
					p.advance()
				}
			}
			mod.Loc.End = p.pos()
			stmt.Modifications = append(stmt.Modifications, mod)
			if p.cur.Type == ',' {
				p.advance()
			}
		}
		// Optional using_type_clause
		if p.cur.Type == kwUSING {
			p.advance() // consume USING
			stmt.UsingType = p.parseObjectName()
			// Optional array_DML_clause
			if p.cur.Type == kwWITH {
				p.advance()
				if p.isIdentLikeStr("ARRAY") {
					stmt.ArrayDML = true
					p.advance() // consume ARRAY
					if p.isIdentLikeStr("DML") {
						p.advance() // consume DML
					}
					p.skipParenthesizedBlock()
				}
			}
		}
	}

	// Optional WITH LOCAL PARTITION / storage_table_clause
	for p.cur.Type == kwWITH {
		p.advance() // consume WITH
		if p.cur.Type == kwLOCAL {
			stmt.WithLocal = true
			p.advance()
			if p.cur.Type == kwPARTITION {
				p.advance()
			}
		} else if p.isIdentLikeStr("SYSTEM") || p.cur.Type == kwUSER {
			stmt.StorageTable = p.parseIdentifier()
			if p.isIdentLikeStr("MANAGED") {
				p.advance()
			}
			if p.isIdentLikeStr("STORAGE") {
				p.advance()
			}
			if p.isIdentLikeStr("TABLES") {
				p.advance()
			}
		} else {
			break
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// ---------------------------------------------------------------------------
// CREATE OPERATOR
// ---------------------------------------------------------------------------

// parseCreateOperatorStmt parses a CREATE OPERATOR statement.
// Called after CREATE [OR REPLACE | IF NOT EXISTS] has been consumed and OPERATOR consumed.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-OPERATOR.html
//
//	CREATE [ OR REPLACE | IF NOT EXISTS ] OPERATOR
//	    [ schema. ] operator
//	    binding_clause
//	    [ SHARING = { METADATA | NONE } ]
func (p *Parser) parseCreateOperatorStmt(start int, orReplace bool, ifNotExists bool) *nodes.CreateOperatorStmt {
	stmt := &nodes.CreateOperatorStmt{
		OrReplace:   orReplace,
		IfNotExists: ifNotExists,
		Loc:         nodes.Loc{Start: start},
	}

	// Operator name
	stmt.Name = p.parseObjectName()

	// binding_clause: BINDING (types) RETURN type USING func
	// There can be multiple bindings (comma-separated in the BNF for CREATE, but
	// typically one binding per CREATE OPERATOR)
	if p.isIdentLikeStr("BINDING") {
		binding := p.parseOperatorBinding()
		if binding != nil {
			stmt.Bindings = append(stmt.Bindings, binding)
		}
	}

	// Optional SHARING = { METADATA | NONE }
	if p.isIdentLikeStr("SHARING") {
		p.advance() // consume SHARING
		if p.cur.Type == '=' {
			p.advance() // consume =
		}
		if p.isIdentLikeStr("METADATA") {
			stmt.Sharing = "METADATA"
			p.advance()
		} else if p.isIdentLikeStr("NONE") {
			stmt.Sharing = "NONE"
			p.advance()
		} else if p.isIdentLike() {
			stmt.Sharing = p.cur.Str
			p.advance()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}

// parseOperatorBinding parses a BINDING clause.
//
//	BINDING ( [ parameter_type [, parameter_type ]... ] )
//	    RETURN return_type
//	    [ implementation_clause ]
//	    using_function_clause
func (p *Parser) parseOperatorBinding() *nodes.OperatorBinding {
	if !p.isIdentLikeStr("BINDING") {
		return nil
	}
	start := p.pos()
	p.advance() // consume BINDING

	binding := &nodes.OperatorBinding{
		Loc: nodes.Loc{Start: start},
	}

	// ( parameter_types )
	if p.cur.Type == '(' {
		p.advance()
		if p.cur.Type != ')' {
			binding.ParamTypes = p.parseTypeNameList()
		}
		if p.cur.Type == ')' {
			p.advance()
		}
	}

	// RETURN return_type
	if p.cur.Type == kwRETURN {
		p.advance() // consume RETURN
		// Return type may be in parens: RETURN ( datatype )
		if p.cur.Type == '(' {
			p.advance()
			binding.ReturnType = p.parseTypeNameStr()
			if p.cur.Type == ')' {
				p.advance()
			}
		} else {
			binding.ReturnType = p.parseTypeNameStr()
		}
	}

	// Optional implementation_clause
	// { ANCILLARY TO primary_operator (types) | context_clause }
	if p.isIdentLikeStr("ANCILLARY") {
		p.advance() // consume ANCILLARY
		if p.cur.Type == kwTO {
			p.advance() // consume TO
		}
		binding.AncillaryTo = p.parseObjectName()
		if p.cur.Type == '(' {
			p.advance()
			binding.AncillaryParams = p.parseTypeNameList()
			if p.cur.Type == ')' {
				p.advance()
			}
		}
	} else if p.cur.Type == kwWITH {
		// context_clause: WITH INDEX CONTEXT, SCAN CONTEXT implementation_type
		p.advance() // consume WITH
		if p.cur.Type == kwINDEX {
			p.advance() // consume INDEX
			if p.isIdentLikeStr("CONTEXT") {
				binding.WithIndexCtx = true
				p.advance() // consume CONTEXT
			}
			if p.cur.Type == ',' {
				p.advance() // consume ','
			}
			if p.isIdentLikeStr("SCAN") {
				p.advance() // consume SCAN
				if p.isIdentLikeStr("CONTEXT") {
					p.advance() // consume CONTEXT
				}
				binding.ScanCtxType = p.parseIdentifier()
			}
		} else if p.isIdentLikeStr("COLUMN") {
			binding.WithColumnCtx = true
			p.advance() // consume COLUMN
			if p.isIdentLikeStr("CONTEXT") {
				p.advance() // consume CONTEXT
			}
		}
	}

	// COMPUTE ANCILLARY DATA
	if p.isIdentLikeStr("COMPUTE") {
		binding.ComputeAnc = true
		p.advance() // consume COMPUTE
		if p.isIdentLikeStr("ANCILLARY") {
			p.advance() // consume ANCILLARY
		}
		if p.isIdentLikeStr("DATA") {
			p.advance() // consume DATA
		}
	}

	// WITH COLUMN CONTEXT (can appear after COMPUTE ANCILLARY DATA)
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		if p.isIdentLikeStr("COLUMN") {
			binding.WithColumnCtx = true
			p.advance() // consume COLUMN
			if p.isIdentLikeStr("CONTEXT") {
				p.advance() // consume CONTEXT
			}
		}
	}

	// USING function
	if p.cur.Type == kwUSING {
		p.advance() // consume USING
		binding.UsingFunc = p.parseObjectName()
	}

	binding.Loc.End = p.pos()
	return binding
}

// ---------------------------------------------------------------------------
// ALTER OPERATOR
// ---------------------------------------------------------------------------

// parseAlterOperatorStmt parses an ALTER OPERATOR statement.
// Called after ALTER has been consumed and OPERATOR consumed.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-OPERATOR.html
//
//	ALTER OPERATOR [ IF EXISTS ] [ schema. ] operator
//	    { add_binding_clause
//	    | drop_binding_clause
//	    | COMPILE
//	    }
func (p *Parser) parseAlterOperatorStmt(start int) *nodes.AlterOperatorStmt {
	stmt := &nodes.AlterOperatorStmt{
		Loc: nodes.Loc{Start: start},
	}

	// Optional IF EXISTS
	if p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwEXISTS {
			stmt.IfExists = true
			p.advance() // consume IF
			p.advance() // consume EXISTS
		}
	}

	// Operator name
	stmt.Name = p.parseObjectName()

	// Action
	switch {
	case p.isIdentLikeStr("COMPILE"):
		stmt.Action = "COMPILE"
		p.advance()
	case p.cur.Type == kwADD:
		stmt.Action = "ADD_BINDING"
		p.advance() // consume ADD
		binding := p.parseOperatorBinding()
		stmt.Binding = binding
	case p.cur.Type == kwDROP:
		stmt.Action = "DROP_BINDING"
		p.advance() // consume DROP
		if p.isIdentLikeStr("BINDING") {
			p.advance() // consume BINDING
		}
		// ( datatype [, datatype ]... )
		if p.cur.Type == '(' {
			p.advance()
			stmt.DropTypes = p.parseTypeNameList()
			if p.cur.Type == ')' {
				p.advance()
			}
		}
		// Optional FORCE
		if p.cur.Type == kwFORCE {
			stmt.DropForce = true
			p.advance()
		}
	}

	stmt.Loc.End = p.pos()
	return stmt
}
