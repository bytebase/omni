package parser

import "github.com/bytebase/omni/snowflake/ast"

// ---------------------------------------------------------------------------
// CREATE statement dispatch
// ---------------------------------------------------------------------------

// parseCreateStmt parses CREATE ... and dispatches to the appropriate
// sub-parser based on the object type keyword (TABLE, etc.).
func (p *Parser) parseCreateStmt() (ast.Node, error) {
	createTok := p.advance() // consume CREATE
	start := createTok.Loc

	// OR { REPLACE | ALTER }. These are mutually exclusive prefixes: REPLACE
	// drops and recreates the object, ALTER (a preview form documented at
	// CREATE OR ALTER) modifies it in place. Only one OR-modifier is consumed
	// here, so a second one (CREATE OR ALTER OR REPLACE ...) is left for the
	// object-type dispatch to reject.
	orReplace := false
	orAlter := false
	if p.cur.Type == kwOR {
		switch p.peekNext().Type {
		case kwREPLACE:
			p.advance() // consume OR
			p.advance() // consume REPLACE
			orReplace = true
		case kwALTER:
			p.advance() // consume OR
			p.advance() // consume ALTER
			orAlter = true
		}
	}

	// Optional SECURE modifier. SECURE applies to VIEW / MATERIALIZED VIEW
	// (T2.4) and to FUNCTION / EXTERNAL FUNCTION / PROCEDURE (T4.5). It is
	// checked before the temporary modifiers since it precedes the object type.
	secure := false
	if p.cur.Type == kwSECURE {
		p.advance()
		secure = true
	}

	// Optional RECURSIVE modifier (VIEW-only).
	recursive := false
	if p.cur.Type == kwRECURSIVE {
		p.advance()
		recursive = true
	}

	// If SECURE or RECURSIVE were consumed, dispatch early on the object type.
	// RECURSIVE pairs only with VIEW; SECURE additionally precedes
	// FUNCTION / EXTERNAL FUNCTION / PROCEDURE.
	if secure || recursive {
		switch p.cur.Type {
		case kwVIEW:
			return p.parseCreateViewStmt(start, orReplace, orAlter, secure, recursive)
		case kwMATERIALIZED:
			p.advance() // consume MATERIALIZED
			return p.parseCreateMaterializedViewStmt(start, orReplace, orAlter, secure)
		case kwFUNCTION:
			if recursive {
				return p.unsupported("CREATE")
			}
			// CREATE [OR REPLACE] SECURE FUNCTION ... (T4.5).
			return p.parseCreateFunctionStmt(start, orReplace, orAlter, true, false)
		case kwPROCEDURE:
			if recursive {
				return p.unsupported("CREATE")
			}
			// CREATE [OR REPLACE] SECURE PROCEDURE ... (T4.5).
			return p.parseCreateProcedureStmt(start, orReplace, orAlter, true)
		case kwEXTERNAL:
			if recursive {
				return p.unsupported("CREATE")
			}
			// CREATE [OR REPLACE] SECURE EXTERNAL FUNCTION ... (T4.5).
			p.advance() // consume EXTERNAL
			if p.cur.Type != kwFUNCTION {
				// EXTERNAL TABLE (and any other EXTERNAL object) is owned by another node.
				return p.unsupported("CREATE")
			}
			return p.parseCreateExternalFunctionStmt(start, orReplace, orAlter, true)
		default:
			return p.unsupported("CREATE")
		}
	}

	// Optional [LOCAL|GLOBAL] TRANSIENT|TEMPORARY|TEMP|VOLATILE
	transient := false
	temporary := false
	volatile := false

	// Consume optional LOCAL/GLOBAL prefix (they don't change semantics here)
	if p.cur.Type == kwLOCAL || p.cur.Type == kwGLOBAL {
		p.advance()
	}

	switch p.cur.Type {
	case kwTRANSIENT:
		p.advance()
		transient = true
	case kwTEMPORARY, kwTEMP:
		p.advance()
		temporary = true
	case kwVOLATILE:
		p.advance()
		volatile = true
	}

	// Optional HYBRID modifier (Unistore). HYBRID is not a reserved keyword (it
	// lexes as a plain identifier), so it is detected the non-reserved way via
	// curIsWord and only when immediately followed by TABLE — otherwise a bare
	// "HYBRID" identifier object (e.g. CREATE HYBRID VIEW, were that a thing) is
	// left for the dispatch switch to reject. CREATE HYBRID TABLE is the only
	// documented form, so HYBRID pairs only with TABLE.
	hybrid := false
	if p.curIsWord("HYBRID") && p.peekNext().Type == kwTABLE {
		p.advance() // consume HYBRID
		hybrid = true
	}

	switch p.cur.Type {
	case kwTABLE:
		return p.parseCreateTableStmt(start, orReplace, orAlter, transient, temporary, volatile, hybrid)
	case kwDATABASE:
		// CREATE [OR REPLACE] DATABASE ROLE ... (T4.6) vs CREATE DATABASE ... (T2.1).
		// DATABASE ROLE is disambiguated by the ROLE keyword following DATABASE.
		if p.peekNext().Type == kwROLE {
			p.advance() // consume DATABASE
			return p.parseCreateRoleStmt(start, orReplace, orAlter, true)
		}
		return p.parseCreateDatabaseStmt(start, orReplace, orAlter, transient)
	case kwSCHEMA:
		return p.parseCreateSchemaStmt(start, orReplace, orAlter, transient)
	case kwVIEW:
		return p.parseCreateViewStmt(start, orReplace, orAlter, false, false)
	case kwMATERIALIZED:
		p.advance() // consume MATERIALIZED
		return p.parseCreateMaterializedViewStmt(start, orReplace, orAlter, false)
	case kwSTAGE:
		// CREATE [OR REPLACE] [TEMP|TEMPORARY] STAGE ... (T4.1).
		return p.parseCreateStageStmt(start, orReplace, orAlter, temporary)
	case kwFILE_FORMAT, kwFILE:
		// CREATE [OR REPLACE] [TEMP|TEMPORARY|VOLATILE] FILE FORMAT ... (T4.2).
		// FILE FORMAT lexes as one FILE_FORMAT token or two FILE+FORMAT tokens;
		// both dispatch here. Per the docs VOLATILE is a synonym of TEMPORARY for
		// a file format, so either modifier sets the statement's Temporary flag.
		return p.parseCreateFileFormatStmt(start, orReplace, orAlter, temporary || volatile)
	case kwPIPE:
		// CREATE [OR REPLACE] PIPE ... (T4.3).
		return p.parseCreatePipeStmt(start, orReplace, orAlter)
	case kwSTREAM:
		// CREATE [OR REPLACE] STREAM ... (T4.3).
		return p.parseCreateStreamStmt(start, orReplace, orAlter)
	case kwTASK:
		// CREATE [OR REPLACE] TASK ... (T4.3).
		return p.parseCreateTaskStmt(start, orReplace, orAlter)
	case kwALERT:
		// CREATE [OR REPLACE] ALERT ... (T4.3).
		return p.parseCreateAlertStmt(start, orReplace, orAlter)
	case kwFUNCTION:
		// CREATE [OR REPLACE] [TEMP|TEMPORARY] FUNCTION ... (T4.5).
		return p.parseCreateFunctionStmt(start, orReplace, orAlter, false, temporary)
	case kwPROCEDURE:
		// CREATE [OR REPLACE] PROCEDURE ... (T4.5).
		return p.parseCreateProcedureStmt(start, orReplace, orAlter, false)
	case kwDYNAMIC:
		// CREATE [OR REPLACE] [TRANSIENT] DYNAMIC [ICEBERG] TABLE ... (T4.4).
		return p.parseCreateDynamicTableStmt(start, orReplace, orAlter, transient)
	case kwEVENT:
		// CREATE [OR REPLACE] EVENT TABLE ... (T4.4).
		return p.parseCreateEventTableStmt(start, orReplace, orAlter)
	case kwSEQUENCE:
		// CREATE [OR REPLACE] SEQUENCE ... (T4.4).
		return p.parseCreateSequenceStmt(start, orReplace, orAlter)
	case kwWAREHOUSE:
		// CREATE [ OR REPLACE | OR ALTER ] WAREHOUSE ... (gap-warehouse).
		return p.parseCreateWarehouseStmt(start, orReplace, orAlter)
	case kwEXTERNAL:
		// CREATE [OR REPLACE] EXTERNAL { FUNCTION (T4.5) | TABLE (T4.4) | VOLUME (T4.7) }.
		// The object is disambiguated by what follows EXTERNAL. EXTERNAL TABLE keeps
		// cur on the EXTERNAL keyword (its sub-parser consumes both EXTERNAL and
		// TABLE), matching the DYNAMIC/EVENT two-keyword handling. VOLUME is not a
		// reserved keyword, so EXTERNAL VOLUME lexes as kwEXTERNAL followed by a
		// "VOLUME" identifier; that branch is taken here.
		if p.peekNext().Type == kwTABLE {
			return p.parseCreateExternalTableStmt(start, orReplace, orAlter)
		}
		p.advance() // consume EXTERNAL
		if p.curIsWord("VOLUME") {
			p.advance() // consume VOLUME (a plain identifier; not a reserved keyword)
			return p.parseCreateExternalVolumeStmt(start, orReplace, orAlter)
		}
		if p.cur.Type != kwFUNCTION {
			return p.unsupported("CREATE")
		}
		return p.parseCreateExternalFunctionStmt(start, orReplace, orAlter, false)
	case kwSTORAGE, kwAPI, kwNOTIFICATION, kwSECURITY, kwRESOURCE, kwSECRET, kwCONNECTION, kwGIT:
		// CREATE [OR REPLACE] account-level integration objects (T4.7):
		// { STORAGE | API | NOTIFICATION | SECURITY } INTEGRATION, RESOURCE MONITOR,
		// SECRET, CONNECTION, GIT REPOSITORY. parseCreateIntegrationStmt consumes the
		// object-type keyword(s).
		return p.parseCreateIntegrationStmt(start, orReplace, orAlter)
	case kwTAG:
		// CREATE [OR REPLACE] TAG ... (T4.9).
		return p.parseCreateTagStmt(start, orReplace, orAlter)
	case kwSEMANTIC:
		// CREATE [OR REPLACE] SEMANTIC VIEW ... (T4.9). The sub-parser consumes
		// SEMANTIC + VIEW.
		return p.parseCreateSemanticViewStmt(start, orReplace, orAlter)
	case kwDATASET:
		// CREATE [OR REPLACE] DATASET ... (T4.9).
		return p.parseCreateDatasetStmt(start, orReplace, orAlter)
	case kwFAILOVER, kwREPLICATION:
		// CREATE [OR REPLACE] { FAILOVER | REPLICATION } GROUP ... (T4.8).
		// parseCreateReplicationGroupStmt consumes the kind keyword + GROUP.
		return p.parseCreateReplicationGroupStmt(start, orReplace, orAlter)
	case kwACCOUNT, kwMANAGED:
		// CREATE ACCOUNT / CREATE [OR REPLACE] MANAGED ACCOUNT ... (T4.8).
		// parseCreateAccountStmt consumes the (MANAGED) ACCOUNT keyword(s).
		return p.parseCreateAccountStmt(start, orReplace, orAlter)
	case kwSHARE:
		// CREATE [OR REPLACE] SHARE ... (T4.8).
		return p.parseCreateShareStmt(start, orReplace, orAlter)
	case kwROLE:
		// CREATE [OR REPLACE] ROLE ... (T4.6). DATABASE ROLE is dispatched by the
		// kwDATABASE case above.
		return p.parseCreateRoleStmt(start, orReplace, orAlter, false)
	case kwUSER:
		// CREATE [OR REPLACE] USER ... (T4.6).
		return p.parseCreateUserStmt(start, orReplace, orAlter)
	case kwMASKING, kwSESSION, kwPASSWORD, kwNETWORK, kwROW:
		// CREATE [OR REPLACE] NETWORK RULE ... (gap-network-rule) vs the policy
		// objects below. NETWORK is reserved but RULE is not, so NETWORK RULE lexes
		// as kwNETWORK followed by a "RULE" identifier; NETWORK POLICY lexes as
		// kwNETWORK followed by kwPOLICY. The RULE form is dispatched here before the
		// shared policy path.
		if p.cur.Type == kwNETWORK && p.peekIsWord("RULE") {
			return p.parseCreateNetworkRuleStmt(start, orReplace)
		}
		// CREATE [OR REPLACE] { MASKING | ROW ACCESS | SESSION | PASSWORD | NETWORK }
		// POLICY ... (T4.6). parseCreatePolicyDispatch consumes the kind keyword(s)
		// and the POLICY keyword, then delegates. A leading keyword not followed by
		// the expected POLICY (or ACCESS POLICY, for ROW) is not a policy statement.
		return p.parseCreatePolicyDispatch(start, orReplace, orAlter)
	default:
		// AUTHENTICATION is not a reserved keyword, so CREATE AUTHENTICATION POLICY
		// lexes as an "AUTHENTICATION" identifier followed by the POLICY keyword;
		// that form is dispatched here (T4.6).
		if p.curIsWord("AUTHENTICATION") && p.peekNext().Type == kwPOLICY {
			return p.parseCreatePolicyDispatch(start, orReplace, orAlter)
		}
		return p.unsupported("CREATE")
	}
}

// parseCreatePolicyDispatch identifies which policy object follows and consumes
// the policy-kind keyword(s) and the POLICY keyword, then delegates to
// parseCreatePolicyStmt. On entry cur is the first policy-kind word:
//
//	MASKING POLICY            -> PolicyMasking
//	ROW ACCESS POLICY         -> PolicyRowAccess
//	SESSION POLICY            -> PolicySession
//	PASSWORD POLICY           -> PolicyPassword
//	NETWORK POLICY            -> PolicyNetwork
//	AUTHENTICATION POLICY     -> PolicyAuthentication  (AUTHENTICATION is an identifier)
func (p *Parser) parseCreatePolicyDispatch(start ast.Loc, orReplace, orAlter bool) (ast.Node, error) {
	kind, err := p.consumePolicyKeywords()
	if err != nil {
		return nil, err
	}
	return p.parseCreatePolicyStmt(start, orReplace, orAlter, kind)
}

// ---------------------------------------------------------------------------
// CREATE TABLE statement parser
// ---------------------------------------------------------------------------

// parseCreateTableStmt parses the body of a CREATE [OR REPLACE] [...] TABLE
// statement. The CREATE keyword and optional modifiers have already been
// consumed; start is the Loc of the CREATE token.
func (p *Parser) parseCreateTableStmt(start ast.Loc, orReplace, orAlter, transient, temporary, volatile, hybrid bool) (ast.Node, error) {
	p.advance() // consume TABLE

	stmt := &ast.CreateTableStmt{
		OrReplace: orReplace,
		OrAlter:   orAlter,
		Transient: transient,
		Temporary: temporary,
		Volatile:  volatile,
		Hybrid:    hybrid,
		Loc:       ast.Loc{Start: start.Start},
	}

	// IF NOT EXISTS
	if p.cur.Type == kwIF {
		next := p.peekNext()
		if next.Type == kwNOT {
			p.advance() // consume IF
			p.advance() // consume NOT
			if _, err := p.expect(kwEXISTS); err != nil {
				return nil, err
			}
			stmt.IfNotExists = true
		}
	}

	// Table name
	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Branch on what follows the table name
	switch p.cur.Type {
	case kwLIKE:
		// CREATE TABLE ... LIKE source_table [table_properties]
		p.advance() // consume LIKE
		src, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Like = src
		if err := p.parseTableProperties(stmt); err != nil {
			return nil, err
		}

	case kwCLONE:
		// CREATE TABLE ... CLONE source [AT|BEFORE ...]
		clone, err := p.parseCloneSource()
		if err != nil {
			return nil, err
		}
		stmt.Clone = clone

	case kwAS:
		// CREATE TABLE ... AS SELECT (CTAS, no column list)
		p.advance() // consume AS
		query, err := p.parseQueryExpr()
		if err != nil {
			return nil, err
		}
		stmt.AsSelect = query

	case '(':
		// CREATE TABLE ... ( column_defs [, constraints] ) [table_properties] [AS SELECT]
		if err := p.parseColumnDeclItems(stmt); err != nil {
			return nil, err
		}
		if err := p.parseTableProperties(stmt); err != nil {
			return nil, err
		}
		// Optionally followed by AS SELECT (CTAS with column list)
		if p.cur.Type == kwAS {
			p.advance() // consume AS
			query, err := p.parseQueryExpr()
			if err != nil {
				return nil, err
			}
			stmt.AsSelect = query
		}

	default:
		// No body — possibly just table properties (e.g. CREATE TABLE t COMMENT = '...')
		if err := p.parseTableProperties(stmt); err != nil {
			return nil, err
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// Column declaration items
// ---------------------------------------------------------------------------

// parseColumnDeclItems parses the parenthesized list of column definitions
// and out-of-line constraints: ( item, item, ... ).
func (p *Parser) parseColumnDeclItems(stmt *ast.CreateTableStmt) error {
	if _, err := p.expect('('); err != nil {
		return err
	}

	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		// An inline INDEX element (HYBRID TABLE) shares its leading keyword with a
		// column whose name is literally "index" (INDEX is non-reserved).
		// tryParseTableIndex speculatively matches the INDEX <name> ( ... ) shape
		// and reports matched=false (rolling back) when it is actually a column
		// named "index", which then falls through to parseColumnDef below.
		if p.cur.Type == kwINDEX {
			idx, matched, err := p.tryParseTableIndex()
			if err != nil {
				return err
			}
			if matched {
				stmt.Indexes = append(stmt.Indexes, idx)
				if p.cur.Type == ',' {
					p.advance() // consume ','
					continue
				}
				break
			}
		}

		if p.isOutOfLineConstraintStart() {
			con, err := p.parseOutOfLineConstraint()
			if err != nil {
				return err
			}
			stmt.Constraints = append(stmt.Constraints, con)
		} else {
			col, err := p.parseColumnDef()
			if err != nil {
				return err
			}
			stmt.Columns = append(stmt.Columns, col)
		}

		if p.cur.Type == ',' {
			p.advance() // consume ','
		} else {
			break
		}
	}

	if _, err := p.expect(')'); err != nil {
		return err
	}
	return nil
}

// isOutOfLineConstraintStart returns true if the current token starts an
// out-of-line (table-level) constraint definition.
func (p *Parser) isOutOfLineConstraintStart() bool {
	switch p.cur.Type {
	case kwCONSTRAINT, kwPRIMARY, kwUNIQUE, kwFOREIGN:
		return true
	}
	return false
}

// tryParseTableIndex speculatively parses an inline INDEX element of a HYBRID
// TABLE column list:
//
//	INDEX <index_name> ( <col_name> [, <col_name>...] )
//
// On entry cur is the INDEX keyword. Because INDEX is a non-reserved keyword,
// the same leading token can begin a column definition whose name is literally
// "index" — e.g. `index INT`, `index VARCHAR(10)`, or even `index NUMBER(38,0)`
// (where the type's own parenthesized params would superficially resemble an
// index column list). To disambiguate unambiguously, the whole element shape is
// parsed under a snapshot: the index name, '(', a non-empty identifier list, and
// ')'. If every part matches, the index is committed (matched=true). If any part
// fails — wrong name, no '(', a numeric type param instead of a column name,
// missing ')' — the parser is rolled back to the INDEX keyword and matched=false
// is returned, so the caller parses a column definition instead. This never
// surfaces a hard error from the speculative path; a genuinely malformed INDEX
// element instead produces a column-definition error downstream.
func (p *Parser) tryParseTableIndex() (idx *ast.TableIndex, matched bool, err error) {
	snap := p.snapshot()

	start := p.cur.Loc.Start
	p.advance() // consume INDEX

	name, nameErr := p.parseIdent()
	if nameErr != nil || p.cur.Type != '(' {
		p.restore(snap)
		return nil, false, nil
	}

	cols, colsErr := p.parseIdentListInParens()
	if colsErr != nil {
		// Not an index element after all (e.g. a column named "index" with a
		// parenthesized type): roll back and let parseColumnDef handle it.
		p.restore(snap)
		return nil, false, nil
	}

	return &ast.TableIndex{
		Name:    name,
		Columns: cols,
		Loc:     ast.Loc{Start: start, End: p.prev.Loc.End},
	}, true, nil
}

// ---------------------------------------------------------------------------
// Column definition
// ---------------------------------------------------------------------------

// parseColumnDef parses a single column definition:
//
//	name [data_type | AS (expr)] [column_options...]
func (p *Parser) parseColumnDef() (*ast.ColumnDef, error) {
	name, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	col := &ast.ColumnDef{
		Name: name,
		Loc:  ast.Loc{Start: name.Loc.Start},
	}

	// The data type is optional when followed by AS (virtual column
	// defined with an expression). Peek ahead to decide.
	if p.cur.Type != kwAS && p.cur.Type != ',' && p.cur.Type != ')' &&
		p.cur.Type != kwNOT && p.cur.Type != kwNULL && p.cur.Type != kwDEFAULT &&
		p.cur.Type != kwIDENTITY && p.cur.Type != kwAUTOINCREMENT &&
		p.cur.Type != kwCOLLATE && p.cur.Type != kwCONSTRAINT &&
		p.cur.Type != kwPRIMARY && p.cur.Type != kwUNIQUE && p.cur.Type != kwFOREIGN &&
		p.cur.Type != kwCOMMENT && p.cur.Type != kwWITH && p.cur.Type != kwMASKING &&
		p.cur.Type != kwTAG && p.cur.Type != tokEOF {
		dt, err := p.parseDataType()
		if err != nil {
			return nil, err
		}
		col.DataType = dt
	}

	if err := p.parseColumnOptions(col); err != nil {
		return nil, err
	}

	col.Loc.End = p.prev.Loc.End
	return col, nil
}

// parseColumnOptions parses the optional constraints and properties that
// follow a column name+type. Returns nil when no more column options are
// found.
func (p *Parser) parseColumnOptions(col *ast.ColumnDef) error {
	for {
		switch p.cur.Type {
		case kwNOT:
			// NOT NULL
			next := p.peekNext()
			if next.Type == kwNULL {
				p.advance() // consume NOT
				p.advance() // consume NULL
				col.NotNull = true
			} else {
				return nil
			}

		case kwNULL:
			p.advance() // consume NULL
			col.Nullable = true

		case kwDEFAULT:
			p.advance() // consume DEFAULT
			expr, err := p.parseExpr()
			if err != nil {
				return err
			}
			col.Default = expr

		case kwIDENTITY, kwAUTOINCREMENT:
			spec, err := p.parseIdentitySpec()
			if err != nil {
				return err
			}
			col.Identity = spec

		case kwCOLLATE:
			p.advance() // consume COLLATE
			tok, err := p.expect(tokString)
			if err != nil {
				return err
			}
			col.Collate = tok.Str

		case kwCONSTRAINT, kwPRIMARY, kwUNIQUE, kwFOREIGN:
			ic, err := p.parseInlineConstraint()
			if err != nil {
				return err
			}
			col.InlineConstraint = ic

		case kwCOMMENT:
			p.advance() // consume COMMENT
			tok, err := p.expect(tokString)
			if err != nil {
				return err
			}
			s := tok.Str
			col.Comment = &s

		case kwWITH:
			// WITH MASKING POLICY name [USING (...)]
			// WITH TAG (...)
			// WITH ROW ACCESS POLICY ... — consumed and discarded
			next := p.peekNext()
			switch next.Type {
			case kwMASKING:
				p.advance() // consume WITH
				p.advance() // consume MASKING
				if _, err := p.expect(kwPOLICY); err != nil {
					return err
				}
				policyName, err := p.parseObjectName()
				if err != nil {
					return err
				}
				col.MaskingPolicy = policyName
				// Optional USING (col_list)
				if p.cur.Type == kwUSING {
					p.advance() // consume USING
					if err := p.skipParenthesized(); err != nil {
						return err
					}
				}
			case kwTAG:
				p.advance() // consume WITH
				tags, err := p.parseTagAssignments()
				if err != nil {
					return err
				}
				col.Tags = append(col.Tags, tags...)
			case kwROW:
				// WITH ROW ACCESS POLICY ... — consume and discard
				p.advance() // consume WITH
				p.advance() // consume ROW
				if _, err := p.expect(kwACCESS); err != nil {
					return err
				}
				if _, err := p.expect(kwPOLICY); err != nil {
					return err
				}
				// policy name
				if _, err := p.parseObjectName(); err != nil {
					return err
				}
				// Optional ON (cols)
				if p.cur.Type == kwON {
					p.advance() // consume ON
					if err := p.skipParenthesized(); err != nil {
						return err
					}
				}
			default:
				return nil
			}

		case kwMASKING:
			// MASKING POLICY name [USING (...)] — without WITH prefix
			p.advance() // consume MASKING
			if _, err := p.expect(kwPOLICY); err != nil {
				return err
			}
			policyName, err := p.parseObjectName()
			if err != nil {
				return err
			}
			col.MaskingPolicy = policyName
			if p.cur.Type == kwUSING {
				p.advance() // consume USING
				if err := p.skipParenthesized(); err != nil {
					return err
				}
			}

		case kwTAG:
			// TAG (...) — without WITH prefix
			tags, err := p.parseTagAssignments()
			if err != nil {
				return err
			}
			col.Tags = append(col.Tags, tags...)

		case kwAS:
			// AS (expr) — virtual column expression
			p.advance() // consume AS
			if _, err := p.expect('('); err != nil {
				return err
			}
			expr, err := p.parseExpr()
			if err != nil {
				return err
			}
			if _, err := p.expect(')'); err != nil {
				return err
			}
			col.VirtualExpr = expr

		default:
			return nil
		}
	}
}

// ---------------------------------------------------------------------------
// Inline constraint
// ---------------------------------------------------------------------------

// parseInlineConstraint parses a column-level (inline) constraint:
//
//	[CONSTRAINT name] PRIMARY KEY | UNIQUE | FOREIGN KEY REFERENCES ...
func (p *Parser) parseInlineConstraint() (*ast.InlineConstraint, error) {
	ic := &ast.InlineConstraint{
		Loc: p.cur.Loc,
	}

	// Optional CONSTRAINT name
	if p.cur.Type == kwCONSTRAINT {
		p.advance() // consume CONSTRAINT
		name, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		ic.Name = name
	}

	switch p.cur.Type {
	case kwPRIMARY:
		p.advance() // consume PRIMARY
		if _, err := p.expect(kwKEY); err != nil {
			return nil, err
		}
		ic.Type = ast.ConstrPrimaryKey

	case kwUNIQUE:
		p.advance() // consume UNIQUE
		ic.Type = ast.ConstrUnique

	case kwFOREIGN:
		p.advance() // consume FOREIGN
		if _, err := p.expect(kwKEY); err != nil {
			return nil, err
		}
		ref, err := p.parseForeignKeyRef()
		if err != nil {
			return nil, err
		}
		ic.Type = ast.ConstrForeignKey
		ic.References = ref

	default:
		return nil, p.syntaxErrorAtCur()
	}

	if err := p.parseConstraintProperties(); err != nil {
		return nil, err
	}

	ic.Loc.End = p.prev.Loc.End
	return ic, nil
}

// ---------------------------------------------------------------------------
// Out-of-line (table-level) constraint
// ---------------------------------------------------------------------------

// parseOutOfLineConstraint parses a table-level (out-of-line) constraint:
//
//	[CONSTRAINT name] PRIMARY KEY (cols) | UNIQUE (cols) | FOREIGN KEY (cols) REFERENCES ...
func (p *Parser) parseOutOfLineConstraint() (*ast.TableConstraint, error) {
	con := &ast.TableConstraint{
		Loc: p.cur.Loc,
	}

	// Optional CONSTRAINT name
	if p.cur.Type == kwCONSTRAINT {
		p.advance() // consume CONSTRAINT
		name, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		con.Name = name
	}

	switch p.cur.Type {
	case kwPRIMARY:
		p.advance() // consume PRIMARY
		if _, err := p.expect(kwKEY); err != nil {
			return nil, err
		}
		con.Type = ast.ConstrPrimaryKey
		cols, err := p.parseIdentListInParens()
		if err != nil {
			return nil, err
		}
		con.Columns = cols

	case kwUNIQUE:
		p.advance() // consume UNIQUE
		con.Type = ast.ConstrUnique
		cols, err := p.parseIdentListInParens()
		if err != nil {
			return nil, err
		}
		con.Columns = cols

	case kwFOREIGN:
		p.advance() // consume FOREIGN
		if _, err := p.expect(kwKEY); err != nil {
			return nil, err
		}
		con.Type = ast.ConstrForeignKey
		cols, err := p.parseIdentListInParens()
		if err != nil {
			return nil, err
		}
		con.Columns = cols
		ref, err := p.parseForeignKeyRef()
		if err != nil {
			return nil, err
		}
		con.References = ref

	default:
		return nil, p.syntaxErrorAtCur()
	}

	if err := p.parseConstraintProperties(); err != nil {
		return nil, err
	}

	// Optional COMMENT 'text'
	if p.cur.Type == kwCOMMENT {
		p.advance() // consume COMMENT
		tok, err := p.expect(tokString)
		if err != nil {
			return nil, err
		}
		s := tok.Str
		con.Comment = &s
	}

	con.Loc.End = p.prev.Loc.End
	return con, nil
}

// ---------------------------------------------------------------------------
// Identifier list helpers
// ---------------------------------------------------------------------------

// parseIdentListInParens parses ( ident, ident, ... ).
func (p *Parser) parseIdentListInParens() ([]ast.Ident, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	var idents []ast.Ident
	id, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	idents = append(idents, id)

	for p.cur.Type == ',' {
		p.advance() // consume ','
		id, err = p.parseIdent()
		if err != nil {
			return nil, err
		}
		idents = append(idents, id)
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return idents, nil
}

// ---------------------------------------------------------------------------
// Foreign key reference
// ---------------------------------------------------------------------------

// parseForeignKeyRef expects and parses the REFERENCES clause.
func (p *Parser) parseForeignKeyRef() (*ast.ForeignKeyRef, error) {
	if _, err := p.expect(kwREFERENCES); err != nil {
		return nil, err
	}
	return p.parseForeignKeyRefAfterReferences()
}

// parseForeignKeyRefAfterReferences parses the body of a REFERENCES clause
// after the REFERENCES keyword has already been consumed:
//
//	table_name [(cols)] [MATCH FULL|PARTIAL|SIMPLE] [ON DELETE action] [ON UPDATE action]
func (p *Parser) parseForeignKeyRefAfterReferences() (*ast.ForeignKeyRef, error) {
	table, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}

	ref := &ast.ForeignKeyRef{
		Table: table,
	}

	// Optional column list
	if p.cur.Type == '(' {
		cols, err := p.parseIdentListInParens()
		if err != nil {
			return nil, err
		}
		ref.Columns = cols
	}

	// Optional MATCH FULL | PARTIAL | SIMPLE
	if p.cur.Type == kwMATCH {
		p.advance() // consume MATCH
		switch p.cur.Type {
		case kwFULL:
			p.advance()
			ref.Match = "FULL"
		case kwPARTIAL:
			p.advance()
			ref.Match = "PARTIAL"
		case kwSIMPLE:
			p.advance()
			ref.Match = "SIMPLE"
		default:
			return nil, p.syntaxErrorAtCur()
		}
	}

	// Optional ON DELETE / ON UPDATE
	for p.cur.Type == kwON {
		next := p.peekNext()
		switch next.Type {
		case kwDELETE:
			p.advance() // consume ON
			p.advance() // consume DELETE
			action, err := p.parseReferenceAction()
			if err != nil {
				return nil, err
			}
			ref.OnDelete = action
		case kwUPDATE:
			p.advance() // consume ON
			p.advance() // consume UPDATE
			action, err := p.parseReferenceAction()
			if err != nil {
				return nil, err
			}
			ref.OnUpdate = action
		default:
			// Not an ON DELETE/UPDATE — stop
			return ref, nil
		}
	}

	return ref, nil
}

// parseReferenceAction parses a referential action keyword:
//
//	CASCADE | SET NULL | SET DEFAULT | RESTRICT | NO ACTION
func (p *Parser) parseReferenceAction() (ast.ReferenceAction, error) {
	switch p.cur.Type {
	case kwCASCADE:
		p.advance()
		return ast.RefActCascade, nil
	case kwSET:
		p.advance() // consume SET
		switch p.cur.Type {
		case kwNULL:
			p.advance()
			return ast.RefActSetNull, nil
		case kwDEFAULT:
			p.advance()
			return ast.RefActSetDefault, nil
		default:
			return ast.RefActNone, p.syntaxErrorAtCur()
		}
	case kwRESTRICT:
		p.advance()
		return ast.RefActRestrict, nil
	case kwNO:
		p.advance() // consume NO
		if _, err := p.expect(kwACTION); err != nil {
			return ast.RefActNone, err
		}
		return ast.RefActNoAction, nil
	default:
		return ast.RefActNone, p.syntaxErrorAtCur()
	}
}

// ---------------------------------------------------------------------------
// Constraint properties
// ---------------------------------------------------------------------------

// parseConstraintProperties consumes and discards Snowflake constraint
// enforcement/deferral properties. Returns nil when none are present or
// when all have been consumed.
func (p *Parser) parseConstraintProperties() error {
	for {
		switch p.cur.Type {
		case kwNOT:
			next := p.peekNext()
			if next.Type == kwENFORCED {
				p.advance() // consume NOT
				p.advance() // consume ENFORCED
			} else if next.Type == kwDEFERRABLE {
				p.advance() // consume NOT
				p.advance() // consume DEFERRABLE
			} else {
				return nil
			}
		case kwENFORCED:
			p.advance()
		case kwDEFERRABLE:
			p.advance()
		case kwINITIALLY:
			p.advance() // consume INITIALLY
			// Expect DEFERRED or IMMEDIATE
			if p.cur.Type == kwDEFERRED || p.cur.Type == kwIMMEDIATE {
				p.advance()
			} else {
				return p.syntaxErrorAtCur()
			}
		case kwENABLE:
			p.advance() // consume ENABLE
			// Optional VALIDATE / NOVALIDATE
			if p.cur.Type == kwVALIDATE || p.cur.Type == kwNOVALIDATE {
				p.advance()
			}
		case kwDISABLE:
			p.advance() // consume DISABLE
			// Optional VALIDATE / NOVALIDATE
			if p.cur.Type == kwVALIDATE || p.cur.Type == kwNOVALIDATE {
				p.advance()
			}
		case kwVALIDATE:
			p.advance()
		case kwNOVALIDATE:
			p.advance()
		case kwRELY:
			p.advance()
		case kwNORELY:
			p.advance()
		default:
			return nil
		}
	}
}

// ---------------------------------------------------------------------------
// IDENTITY / AUTOINCREMENT
// ---------------------------------------------------------------------------

// parseIdentitySpec parses an IDENTITY or AUTOINCREMENT specification:
//
//	IDENTITY | AUTOINCREMENT [(start, increment)]
//	  [START WITH [=] n] [INCREMENT BY [=] n] [ORDER | NOORDER]
func (p *Parser) parseIdentitySpec() (*ast.IdentitySpec, error) {
	p.advance() // consume IDENTITY or AUTOINCREMENT

	spec := &ast.IdentitySpec{}

	// Optional (start, increment)
	if p.cur.Type == '(' {
		p.advance() // consume '('
		startTok, err := p.expect(tokInt)
		if err != nil {
			return nil, err
		}
		startVal := startTok.Ival
		spec.Start = &startVal

		if _, err := p.expect(','); err != nil {
			return nil, err
		}

		incrTok, err := p.expect(tokInt)
		if err != nil {
			return nil, err
		}
		incrVal := incrTok.Ival
		spec.Increment = &incrVal

		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	// Optional START [WITH] [=] n
	if p.cur.Type == kwSTART {
		p.advance() // consume START
		// Optional WITH
		if p.cur.Type == kwWITH {
			p.advance() // consume WITH
		}
		// Optional =
		if p.cur.Type == '=' {
			p.advance()
		}
		tok, err := p.expect(tokInt)
		if err != nil {
			return nil, err
		}
		v := tok.Ival
		spec.Start = &v
	}

	// Optional INCREMENT [BY] [=] n
	if p.cur.Type == kwINCREMENT {
		p.advance() // consume INCREMENT
		// Optional BY
		if p.cur.Type == kwBY {
			p.advance() // consume BY
		}
		// Optional =
		if p.cur.Type == '=' {
			p.advance()
		}
		tok, err := p.expect(tokInt)
		if err != nil {
			return nil, err
		}
		v := tok.Ival
		spec.Increment = &v
	}

	// Optional ORDER / NOORDER
	if p.cur.Type == kwORDER {
		p.advance()
		b := true
		spec.Order = &b
	} else if p.cur.Type == kwNOORDER {
		p.advance()
		b := false
		spec.Order = &b
	}

	return spec, nil
}

// ---------------------------------------------------------------------------
// Tag assignments
// ---------------------------------------------------------------------------

// parseTagAssignments parses TAG ( name = 'value', ... ).
// The TAG keyword is consumed here.
func (p *Parser) parseTagAssignments() ([]*ast.TagAssignment, error) {
	if _, err := p.expect(kwTAG); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	var tags []*ast.TagAssignment

	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		tagName, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		valueTok, err := p.expect(tokString)
		if err != nil {
			return nil, err
		}
		tags = append(tags, &ast.TagAssignment{
			Name:  tagName,
			Value: valueTok.Str,
		})

		if p.cur.Type == ',' {
			p.advance() // consume ','
		} else {
			break
		}
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return tags, nil
}

// ---------------------------------------------------------------------------
// CLONE source
// ---------------------------------------------------------------------------

// parseCloneSource parses CLONE source_name [AT|BEFORE (kind => value)].
// The CLONE keyword is consumed here.
func (p *Parser) parseCloneSource() (*ast.CloneSource, error) {
	if _, err := p.expect(kwCLONE); err != nil {
		return nil, err
	}

	src, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}

	clone := &ast.CloneSource{
		Source: src,
	}

	// Optional AT|BEFORE (TIMESTAMP|OFFSET|STATEMENT => value)
	if p.cur.Type == kwAT || p.cur.Type == kwBEFORE {
		atBefore := "AT"
		if p.cur.Type == kwBEFORE {
			atBefore = "BEFORE"
		}
		p.advance() // consume AT or BEFORE

		if _, err := p.expect('('); err != nil {
			return nil, err
		}

		var kind string
		switch p.cur.Type {
		case kwTIMESTAMP:
			p.advance()
			kind = "TIMESTAMP"
		case kwOFFSET:
			p.advance()
			kind = "OFFSET"
		case kwSTATEMENT:
			p.advance()
			kind = "STATEMENT"
		default:
			return nil, p.syntaxErrorAtCur()
		}

		// Expect => (tokAssoc)
		if _, err := p.expect(tokAssoc); err != nil {
			return nil, err
		}

		// Value: string or expression — consume as a string token or expression
		valueTok, err := p.expect(tokString)
		if err != nil {
			return nil, err
		}

		if _, err := p.expect(')'); err != nil {
			return nil, err
		}

		clone.AtBefore = atBefore
		clone.Kind = kind
		clone.Value = valueTok.Str
	}

	return clone, nil
}

// ---------------------------------------------------------------------------
// Table properties
// ---------------------------------------------------------------------------

// parseTableProperties parses the optional table-level properties that follow
// the column definition list (or the table name for LIKE/CLONE forms).
// Stops when it encounters a token it doesn't recognize as a table property.
func (p *Parser) parseTableProperties(stmt *ast.CreateTableStmt) error {
	for {
		switch p.cur.Type {
		case kwCLUSTER:
			// CLUSTER BY [LINEAR] (exprs)
			p.advance() // consume CLUSTER
			if _, err := p.expect(kwBY); err != nil {
				return err
			}
			if p.cur.Type == kwLINEAR {
				p.advance() // consume LINEAR
				stmt.Linear = true
			}
			if _, err := p.expect('('); err != nil {
				return err
			}
			exprs, err := p.parseExprList()
			if err != nil {
				return err
			}
			if _, err := p.expect(')'); err != nil {
				return err
			}
			stmt.ClusterBy = exprs

		case kwCOPY:
			// COPY GRANTS
			next := p.peekNext()
			if next.Type == kwGRANTS {
				p.advance() // consume COPY
				p.advance() // consume GRANTS
				stmt.CopyGrants = true
			} else {
				return nil
			}

		case kwCOMMENT:
			// COMMENT = 'text'
			p.advance() // consume COMMENT
			if p.cur.Type == '=' {
				p.advance() // consume =
			}
			tok, err := p.expect(tokString)
			if err != nil {
				return err
			}
			s := tok.Str
			stmt.Comment = &s

		case kwWITH:
			next := p.peekNext()
			switch next.Type {
			case kwTAG:
				p.advance() // consume WITH
				tags, err := p.parseTagAssignments()
				if err != nil {
					return err
				}
				stmt.Tags = append(stmt.Tags, tags...)
			case kwROW:
				// WITH ROW ACCESS POLICY name ON (cols) — consume and discard
				p.advance() // consume WITH
				p.advance() // consume ROW
				if _, err := p.expect(kwACCESS); err != nil {
					return err
				}
				if _, err := p.expect(kwPOLICY); err != nil {
					return err
				}
				// policy name
				if _, err := p.parseObjectName(); err != nil {
					return err
				}
				// Optional ON (cols)
				if p.cur.Type == kwON {
					p.advance() // consume ON
					if err := p.skipParenthesized(); err != nil {
						return err
					}
				}
			default:
				return nil
			}

		case kwTAG:
			// TAG (...) — without WITH prefix
			tags, err := p.parseTagAssignments()
			if err != nil {
				return err
			}
			stmt.Tags = append(stmt.Tags, tags...)

		case kwDATA_RETENTION_TIME_IN_DAYS:
			// DATA_RETENTION_TIME_IN_DAYS = n — consume and discard
			p.advance() // consume DATA_RETENTION_TIME_IN_DAYS
			if p.cur.Type == '=' {
				p.advance()
			}
			if _, err := p.expect(tokInt); err != nil {
				return err
			}

		case kwCHANGE_TRACKING:
			// CHANGE_TRACKING = TRUE|FALSE — consume and discard
			p.advance() // consume CHANGE_TRACKING
			if p.cur.Type == '=' {
				p.advance()
			}
			if p.cur.Type == kwTRUE || p.cur.Type == kwFALSE {
				p.advance()
			} else {
				return p.syntaxErrorAtCur()
			}

		case kwDEFAULT_DDL_COLLATION:
			// DEFAULT_DDL_COLLATION = 'string' — consume and discard
			p.advance() // consume DEFAULT_DDL_COLLATION
			if p.cur.Type == '=' {
				p.advance()
			}
			if _, err := p.expect(tokString); err != nil {
				return err
			}

		case kwSTAGE_FILE_FORMAT:
			// STAGE_FILE_FORMAT = (...) — consume and discard
			p.advance() // consume STAGE_FILE_FORMAT
			if p.cur.Type == '=' {
				p.advance()
			}
			if err := p.skipParenthesized(); err != nil {
				return err
			}

		case kwSTAGE_COPY_OPTIONS:
			// STAGE_COPY_OPTIONS = (...) — consume and discard
			p.advance() // consume STAGE_COPY_OPTIONS
			if p.cur.Type == '=' {
				p.advance()
			}
			if err := p.skipParenthesized(); err != nil {
				return err
			}

		default:
			return nil
		}
	}
}

// ---------------------------------------------------------------------------
// Utility helpers
// ---------------------------------------------------------------------------

// skipParenthesized consumes tokens from '(' to the matching ')' inclusive,
// tracking nesting depth to handle nested parentheses.
func (p *Parser) skipParenthesized() error {
	if _, err := p.expect('('); err != nil {
		return err
	}
	depth := 1
	for depth > 0 && p.cur.Type != tokEOF {
		switch p.cur.Type {
		case '(':
			depth++
		case ')':
			depth--
		}
		p.advance()
	}
	if depth != 0 {
		return p.syntaxErrorAtCur()
	}
	return nil
}
