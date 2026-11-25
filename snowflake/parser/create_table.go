package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

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
	if p.curIs(kwOR) {
		switch {
		case p.nextIs(kwREPLACE):
			p.advance() // consume OR
			p.advance() // consume REPLACE
			orReplace = true
		case p.nextIs(kwALTER):
			p.advance() // consume OR
			p.advance() // consume ALTER
			orAlter = true
		}
	}

	// An object read by its documented shape (object.go) is dispatched by the
	// registry before the switches below.
	if f := p.matchObjectForm(formCreates); f != nil {
		return p.parseObjectStmt("CREATE", f, start, orReplace, orAlter)
	}
	// CATALOG / EXTERNAL ACCESS INTEGRATION (integration.go).
	if p.startsCatalogOrAccessIntegration() {
		return p.parseCreateIntegrationStmt(start, orReplace, orAlter)
	}
	// Policies and data metric functions read by their docs pages (policy.go,
	// data_metric_function.go).
	if d := p.startsDocPolicy(); d != nil {
		return p.parseCreateDocPolicy(start, orReplace, orAlter, d)
	}
	if p.startsCreateDataMetricFunction() {
		return p.parseCreateDataMetricFunction(start, orReplace, orAlter)
	}

	// Optional SECURE modifier. SECURE applies to VIEW / MATERIALIZED VIEW
	// (T2.4) and to FUNCTION / EXTERNAL FUNCTION / PROCEDURE (T4.5). It is
	// checked before the temporary modifiers since it precedes the object type.
	secure := false
	if p.curIs(kwSECURE) {
		p.advance()
		secure = true
	}

	// Optional RECURSIVE modifier (VIEW-only).
	recursive := false
	if p.curIs(kwRECURSIVE) {
		p.advance()
		recursive = true
	}

	// Optional INTERACTIVE modifier (MATERIALIZED VIEW-only). INTERACTIVE is not a
	// reserved keyword (it lexes as a plain identifier), so it is detected the
	// non-reserved way via curIsWord and only when immediately followed by
	// MATERIALIZED — otherwise a bare "INTERACTIVE" identifier object is left for
	// the dispatch switch to reject. INTERACTIVE pairs only with MATERIALIZED VIEW.
	interactive := false
	if p.curIsWord("INTERACTIVE") && p.nextIs(kwMATERIALIZED) {
		p.advance() // consume INTERACTIVE
		interactive = true
	}

	// If SECURE / RECURSIVE / INTERACTIVE were consumed, dispatch early on the
	// object type. RECURSIVE pairs only with VIEW; INTERACTIVE only with
	// MATERIALIZED VIEW; SECURE additionally precedes registered objects that
	// document it (AGENT), AGGREGATE FUNCTION, and FUNCTION / EXTERNAL
	// FUNCTION / PROCEDURE.
	if secure || recursive || interactive {
		if !recursive && !interactive {
			if f := p.matchObjectForm(formSecure); f != nil {
				return p.parseCreateObjectStmt(f, start, orReplace, orAlter, false, true)
			}
			if p.curWord("AGGREGATE") && p.nextIs(kwFUNCTION) {
				p.advance() // consume AGGREGATE
				return p.parseCreateAggregateFunctionStmt(start, orReplace, orAlter, true, false)
			}
		}
		switch {
		case p.curIs(kwVIEW):
			return p.parseCreateViewStmt(start, orReplace, orAlter, secure, recursive)
		case p.curIs(kwMATERIALIZED):
			p.advance() // consume MATERIALIZED
			return p.parseCreateMaterializedViewStmt(start, orReplace, orAlter, secure, interactive)
		case p.curIs(kwFUNCTION):
			if recursive {
				return p.unsupported("CREATE")
			}
			// CREATE [OR REPLACE] SECURE FUNCTION ... (T4.5).
			return p.parseCreateFunctionStmt(start, orReplace, orAlter, true, false)
		case p.curIs(kwPROCEDURE):
			if recursive {
				return p.unsupported("CREATE")
			}
			// CREATE [OR REPLACE] SECURE PROCEDURE ... (T4.5).
			return p.parseCreateProcedureStmt(start, orReplace, orAlter, true, false)
		case p.curIs(kwEXTERNAL):
			if recursive {
				return p.unsupported("CREATE")
			}
			// CREATE [OR REPLACE] SECURE EXTERNAL FUNCTION ... (T4.5).
			p.advance() // consume EXTERNAL
			if !p.curIs(kwFUNCTION) {
				// EXTERNAL TABLE (and any other EXTERNAL object) is owned by another node.
				return p.unsupported("CREATE")
			}
			return p.parseCreateExternalFunctionStmt(start, orReplace, orAlter, true)
		default:
			return p.unsupported("CREATE")
		}
	}

	// CREATE [ OR REPLACE ] AGGREGATE FUNCTION ... (LANGUAGE PYTHON).
	if p.curWord("AGGREGATE") && p.nextIs(kwFUNCTION) {
		p.advance() // consume AGGREGATE
		return p.parseCreateAggregateFunctionStmt(start, orReplace, orAlter, false, false)
	}

	// Optional [LOCAL|GLOBAL] TRANSIENT|TEMPORARY|TEMP|VOLATILE
	transient := false
	temporary := false
	volatile := false
	localOrGlobal := false

	// Consume optional LOCAL/GLOBAL prefix (they don't change semantics here)
	if p.curIs(kwLOCAL) || p.curIs(kwGLOBAL) {
		p.advance()
		localOrGlobal = true
	}

	switch {
	case p.curIs(kwTRANSIENT):
		p.advance()
		transient = true
	case p.curIs(kwTEMPORARY), p.curIs(kwTEMP):
		p.advance()
		temporary = true
		// TEMP[ORARY] READ ONLY TABLE ... CLONE
		// (https://docs.snowflake.com/en/sql-reference/sql/create-table).
		if p.curIs(kwREAD) && p.nextIs(kwONLY) {
			p.advance() // consume READ
			p.advance() // consume ONLY
			if !p.curIs(kwTABLE) {
				return nil, p.syntaxErrorAtCur()
			}
			return p.parseCreateTableBody(&ast.CreateTableStmt{
				OrReplace: orReplace,
				OrAlter:   orAlter,
				Temporary: true,
				ReadOnly:  true,
				Loc:       ast.Loc{Start: start.Start},
			})
		}
		if !localOrGlobal {
			if f := p.matchObjectForm(formTemp); f != nil {
				return p.parseCreateObjectStmt(f, start, orReplace, orAlter, true, false)
			}
			if p.curIs(kwSECURE) {
				p.advance() // consume SECURE
				if f := p.matchObjectForm(formTempSecure); f != nil {
					return p.parseCreateObjectStmt(f, start, orReplace, orAlter, true, true)
				}
				if p.curWord("AGGREGATE") && p.nextIs(kwFUNCTION) {
					p.advance() // consume AGGREGATE
					return p.parseCreateAggregateFunctionStmt(start, orReplace, orAlter, true, true)
				}
				if p.curIs(kwFUNCTION) {
					return p.parseCreateFunctionStmt(start, orReplace, orAlter, true, true)
				}
				return p.unsupported("CREATE")
			}
			if p.curWord("AGGREGATE") && p.nextIs(kwFUNCTION) {
				p.advance() // consume AGGREGATE
				return p.parseCreateAggregateFunctionStmt(start, orReplace, orAlter, false, true)
			}
		}
	case p.curIs(kwVOLATILE):
		p.advance()
		volatile = true
	}

	// CREATE [ OR REPLACE ] INTERACTIVE WAREHOUSE (warehouse.go).
	if !transient && !temporary && !volatile && !orAlter &&
		p.curNameWord("INTERACTIVE") && p.nextIs(kwWAREHOUSE) {
		return p.parseCreateInteractiveWarehouseStmt(start, orReplace)
	}

	// CREATE [ OR REPLACE ] INTERACTIVE TABLE
	// (https://docs.snowflake.com/en/sql-reference/sql/create-interactive-table).
	// INTERACTIVE is not a keyword token; it takes no temporary modifier.
	if !transient && !temporary && !volatile && !orAlter &&
		p.curIsWord("INTERACTIVE") && p.nextIs(kwTABLE) {
		p.advance() // consume INTERACTIVE
		return p.parseCreateTableBody(&ast.CreateTableStmt{
			OrReplace:   orReplace,
			Interactive: true,
			Loc:         ast.Loc{Start: start.Start},
		})
	}

	// Optional HYBRID modifier (Unistore). HYBRID is not a reserved keyword (it
	// lexes as a plain identifier), so it is detected the non-reserved way via
	// curIsWord and only when immediately followed by TABLE — otherwise a bare
	// "HYBRID" identifier object (e.g. CREATE HYBRID VIEW, were that a thing) is
	// left for the dispatch switch to reject. CREATE HYBRID TABLE is the only
	// documented form, so HYBRID pairs only with TABLE.
	hybrid := false
	if p.curIsWord("HYBRID") && p.nextIs(kwTABLE) {
		p.advance() // consume HYBRID
		hybrid = true
	}

	// Optional VERSIONED modifier: CREATE OR ALTER VERSIONED SCHEMA is the only
	// documented form (https://docs.snowflake.com/en/sql-reference/sql/create-versioned-schema).
	// VERSIONED is not a reserved keyword; it is taken only before SCHEMA.
	if p.curNameWord("VERSIONED") && p.nextIs(kwSCHEMA) {
		if !orAlter || transient || temporary || volatile {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume VERSIONED
		return p.parseCreateSchemaStmt(start, orReplace, orAlter, transient, true)
	}

	// CREATE [OR REPLACE] [TRANSIENT] ICEBERG TABLE (iceberg_table.go). ICEBERG is
	// not a reserved keyword, so it is matched by text and only before TABLE.
	if p.curIsWord("ICEBERG") && p.nextIs(kwTABLE) && !orAlter && !temporary && !volatile {
		return p.parseCreateIcebergTableStmt(start, orReplace, transient)
	}

	switch {
	case p.curIs(kwTABLE):
		return p.parseCreateTableStmt(start, orReplace, orAlter, transient, temporary, volatile, hybrid)
	case p.curIs(kwDATABASE):
		// CREATE [OR REPLACE] DATABASE ROLE ... (T4.6) vs CREATE DATABASE ... (T2.1).
		// DATABASE ROLE is disambiguated by the ROLE keyword following DATABASE.
		if p.nextIs(kwROLE) {
			p.advance() // consume DATABASE
			return p.parseCreateRoleStmt(start, orReplace, orAlter, true)
		}
		return p.parseCreateDatabaseStmt(start, orReplace, orAlter, transient)
	case p.curIs(kwSCHEMA):
		return p.parseCreateSchemaStmt(start, orReplace, orAlter, transient, false)
	case p.curIs(kwVIEW):
		return p.parseCreateViewStmt(start, orReplace, orAlter, false, false)
	case p.curIs(kwMATERIALIZED):
		p.advance() // consume MATERIALIZED
		return p.parseCreateMaterializedViewStmt(start, orReplace, orAlter, false, false)
	case p.curIs(kwSTAGE):
		// CREATE [OR REPLACE] [TEMP|TEMPORARY] STAGE ... (T4.1).
		return p.parseCreateStageStmt(start, orReplace, orAlter, temporary)
	case p.curIs(kwFILE_FORMAT), p.curIs(kwFILE):
		// CREATE [OR REPLACE] [TEMP|TEMPORARY|VOLATILE] FILE FORMAT ... (T4.2).
		// FILE FORMAT lexes as one FILE_FORMAT token or two FILE+FORMAT tokens;
		// both dispatch here. Per the docs VOLATILE is a synonym of TEMPORARY for
		// a file format, so either modifier sets the statement's Temporary flag.
		return p.parseCreateFileFormatStmt(start, orReplace, orAlter, temporary || volatile)
	case p.curIs(kwPIPE):
		// CREATE [OR REPLACE] PIPE ... (T4.3).
		return p.parseCreatePipeStmt(start, orReplace, orAlter)
	case p.curIs(kwSTREAM):
		// CREATE [OR REPLACE] STREAM ... (T4.3).
		return p.parseCreateStreamStmt(start, orReplace, orAlter)
	case p.curIs(kwTASK):
		// CREATE [OR REPLACE] TASK ... (T4.3), or TASK GRAPH (task_graph.go).
		task := func() (ast.Node, error) { return p.parseCreateTaskStmt(start, orReplace, orAlter) }
		if p.startsTaskGraph() {
			return p.parseTaskOrGraph(func() (ast.Node, error) {
				return p.parseObjectStmt("CREATE", taskGraphForm, start, orReplace, orAlter)
			}, task)
		}
		return task()
	case p.curIs(kwALERT):
		// CREATE [OR REPLACE] ALERT ... (T4.3).
		return p.parseCreateAlertStmt(start, orReplace, orAlter)
	case p.curIs(kwFUNCTION):
		// CREATE [OR REPLACE] [TEMP|TEMPORARY] FUNCTION ... (T4.5).
		return p.parseCreateFunctionStmt(start, orReplace, orAlter, false, temporary)
	case p.curIs(kwPROCEDURE):
		// CREATE [OR REPLACE] PROCEDURE ... (T4.5).
		return p.parseCreateProcedureStmt(start, orReplace, orAlter, false, transient || temporary || volatile)
	case p.curIs(kwDYNAMIC):
		// CREATE [OR REPLACE] [TRANSIENT] DYNAMIC [ICEBERG] TABLE ... (T4.4).
		return p.parseCreateDynamicTableStmt(start, orReplace, orAlter, transient)
	case p.curIs(kwEVENT):
		// CREATE [OR REPLACE] EVENT TABLE ... (T4.4).
		return p.parseCreateEventTableStmt(start, orReplace, orAlter)
	case p.curIs(kwSEQUENCE):
		// CREATE [OR REPLACE] SEQUENCE ... (T4.4).
		return p.parseCreateSequenceStmt(start, orReplace, orAlter)
	case p.curIs(kwWAREHOUSE):
		// CREATE [ OR REPLACE | OR ALTER ] WAREHOUSE ... (gap-warehouse).
		return p.parseCreateWarehouseStmt(start, orReplace, orAlter)
	case p.curIs(kwEXTERNAL):
		// CREATE [OR REPLACE] EXTERNAL { FUNCTION (T4.5) | TABLE (T4.4) | VOLUME (T4.7) }.
		// The object is disambiguated by what follows EXTERNAL. EXTERNAL TABLE keeps
		// cur on the EXTERNAL keyword (its sub-parser consumes both EXTERNAL and
		// TABLE), matching the DYNAMIC/EVENT two-keyword handling. VOLUME is not a
		// reserved keyword, so EXTERNAL VOLUME lexes as kwEXTERNAL followed by a
		// "VOLUME" identifier; that branch is taken here.
		if p.nextIs(kwTABLE) {
			return p.parseCreateExternalTableStmt(start, orReplace, orAlter)
		}
		p.advance() // consume EXTERNAL
		if p.curIsWord("VOLUME") {
			p.advance() // consume VOLUME (a plain identifier; not a reserved keyword)
			return p.parseCreateExternalVolumeStmt(start, orReplace, orAlter)
		}
		if !p.curIs(kwFUNCTION) {
			return p.unsupported("CREATE")
		}
		return p.parseCreateExternalFunctionStmt(start, orReplace, orAlter, false)
	case p.curIs(kwSTORAGE), p.curIs(kwAPI), p.curIs(kwNOTIFICATION), p.curIs(kwSECURITY), p.curIs(kwRESOURCE), p.curIs(kwSECRET), p.curIs(kwCONNECTION), p.curIs(kwGIT):
		// CREATE [OR REPLACE] account-level integration objects (T4.7):
		// { STORAGE | API | NOTIFICATION | SECURITY } INTEGRATION, RESOURCE MONITOR,
		// SECRET, CONNECTION, GIT REPOSITORY. parseCreateIntegrationStmt consumes the
		// object-type keyword(s).
		return p.parseCreateIntegrationStmt(start, orReplace, orAlter)
	case p.curIs(kwTAG):
		// CREATE [OR REPLACE] TAG ... (T4.9).
		return p.parseCreateTagStmt(start, orReplace, orAlter)
	case p.curIs(kwSEMANTIC):
		// CREATE [OR REPLACE] SEMANTIC VIEW ... (T4.9). The sub-parser consumes
		// SEMANTIC + VIEW.
		return p.parseCreateSemanticViewStmt(start, orReplace, orAlter)
	case p.curIs(kwDATASET):
		// CREATE [OR REPLACE] DATASET ... (T4.9).
		return p.parseCreateDatasetStmt(start, orReplace, orAlter)
	case p.curIs(kwIF):
		// CREATE [ OR REPLACE ] [ IF NOT EXISTS ] DATASET <name>, as the
		// dataset page writes it.
		if p.startsIfNotExistsDataset() {
			return p.parseCreateDatasetStmt(start, orReplace, orAlter)
		}
		return p.unsupported("CREATE")
	case p.curIs(kwFAILOVER), p.curIs(kwREPLICATION):
		// CREATE [OR REPLACE] { FAILOVER | REPLICATION } GROUP ... (T4.8).
		// parseCreateReplicationGroupStmt consumes the kind keyword + GROUP.
		return p.parseCreateReplicationGroupStmt(start, orReplace, orAlter)
	case p.curIs(kwACCOUNT), p.curIs(kwMANAGED):
		// CREATE ACCOUNT / CREATE [OR REPLACE] MANAGED ACCOUNT ... (T4.8).
		// parseCreateAccountStmt consumes the (MANAGED) ACCOUNT keyword(s).
		return p.parseCreateAccountStmt(start, orReplace, orAlter)
	case p.curIs(kwSHARE):
		// CREATE [OR REPLACE] SHARE ... (T4.8).
		return p.parseCreateShareStmt(start, orReplace, orAlter)
	case p.curIs(kwROLE):
		// CREATE [OR REPLACE] ROLE ... (T4.6). DATABASE ROLE is dispatched by the
		// kwDATABASE case above.
		return p.parseCreateRoleStmt(start, orReplace, orAlter, false)
	case p.curIs(kwUSER):
		// CREATE [OR REPLACE] USER ... (T4.6).
		return p.parseCreateUserStmt(start, orReplace, orAlter)
	case p.curIs(kwMASKING), p.curIs(kwSESSION), p.curIs(kwPASSWORD), p.curIs(kwNETWORK), p.curIs(kwROW):
		// CREATE [OR REPLACE] NETWORK RULE ... (gap-network-rule) vs the policy
		// objects below. NETWORK is reserved but RULE is not, so NETWORK RULE lexes
		// as kwNETWORK followed by a "RULE" identifier; NETWORK POLICY lexes as
		// kwNETWORK followed by kwPOLICY. The RULE form is dispatched here before the
		// shared policy path.
		if p.curIs(kwNETWORK) && p.peekIsWord("RULE") {
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
		if p.curIsWord("AUTHENTICATION") && p.nextIs(kwPOLICY) {
			return p.parseCreatePolicyDispatch(start, orReplace, orAlter)
		}
		// CREATE [ OR REPLACE ] AI FUNCTION ... (ai_function.go).
		if p.atWords("AI", "FUNCTION") && !orAlter && !transient && !temporary && !volatile && !hybrid {
			return p.parseCreateAiFunctionStmt(start, orReplace)
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
	return p.parseCreateTableBody(&ast.CreateTableStmt{
		OrReplace: orReplace,
		OrAlter:   orAlter,
		Transient: transient,
		Temporary: temporary,
		Volatile:  volatile,
		Hybrid:    hybrid,
		Loc:       ast.Loc{Start: start.Start},
	})
}

// parseCreateTableBody parses CREATE TABLE from the TABLE keyword on, into a
// statement whose modifiers the caller has set.
//
// The forms are those of https://docs.snowflake.com/en/sql-reference/sql/create-table
// and, for an interactive table,
// https://docs.snowflake.com/en/sql-reference/sql/create-interactive-table.
// A READ ONLY table takes only the CLONE form; an interactive table takes
// CLUSTER BY and ends in AS <query>.
func (p *Parser) parseCreateTableBody(stmt *ast.CreateTableStmt) (ast.Node, error) {
	p.advance() // consume TABLE

	// IF NOT EXISTS
	if p.curIs(kwIF) {
		if p.nextIs(kwNOT) {
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

	// Postfix IF NOT EXISTS: Snowflake (and the legacy grammar) also accept
	// the clause AFTER the table name — CREATE TABLE t1 IF NOT EXISTS (...).
	if !stmt.IfNotExists && p.curIs(kwIF) && p.nextIs(kwNOT) {
		p.advance() // consume IF
		p.advance() // consume NOT
		if _, err := p.expect(kwEXISTS); err != nil {
			return nil, err
		}
		stmt.IfNotExists = true
	}

	if stmt.ReadOnly && !p.curIs(kwCLONE) {
		return nil, p.syntaxErrorAtCur()
	}
	if stmt.Interactive && !p.curIs('(') && !p.curIs(kwCLUSTER) {
		return nil, p.syntaxErrorAtCur()
	}

	// Branch on what follows the table name
	switch {
	case p.curIs(kwLIKE):
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

	case p.curIs(kwCLONE):
		// CREATE TABLE ... CLONE source [AT|BEFORE ...]
		clone, err := p.parseCloneSource()
		if err != nil {
			return nil, err
		}
		stmt.Clone = clone
		// [ COPY GRANTS ] [ COPY TAGS ] [ COPY = ( ... ) ]
		for p.startsCopyClause() {
			if err := p.parseTableCopyClause(stmt); err != nil {
				return nil, err
			}
		}

	case p.curIs(kwFROM):
		if err := p.parseCreateTableFrom(stmt); err != nil {
			return nil, err
		}

	case p.curIs('('):
		// CREATE TABLE ... ( column_defs [, constraints] ) [table_properties] [AS SELECT]
		if err := p.parseColumnDeclItems(stmt); err != nil {
			return nil, err
		}
		if stmt.Interactive {
			if err := p.checkInteractiveColumns(stmt); err != nil {
				return nil, err
			}
		}
		if err := p.parseTableProperties(stmt); err != nil {
			return nil, err
		}
		// Optionally followed by AS SELECT (CTAS with column list)
		if err := p.parseCreateTableAsQuery(stmt); err != nil {
			return nil, err
		}

	default:
		// No column list: [ COPY GRANTS ] USING TEMPLATE ( <query> ); CTAS,
		// whose CLUSTER BY / COPY clauses come before AS; or just table
		// properties (e.g. CREATE TABLE t COMMENT = '...').
		if matched, err := p.tryParseUsingTemplate(stmt); err != nil || matched {
			if err != nil {
				return nil, err
			}
			break
		}
		if err := p.parseTableProperties(stmt); err != nil {
			return nil, err
		}
		if err := p.parseCreateTableAsQuery(stmt); err != nil {
			return nil, err
		}
	}

	if stmt.Interactive && stmt.AsSelect == nil {
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// checkInteractiveColumns rejects what an interactive table's column list
// cannot hold: its page gives each column a name and a data type, with an
// optional masking policy, and no constraint or index.
func (p *Parser) checkInteractiveColumns(stmt *ast.CreateTableStmt) error {
	bad := func(loc ast.Loc) error {
		return &ParseError{Position: loc.Start, End: loc.End, Message: "an interactive table's column takes a name, a data type and an optional masking policy"}
	}
	for _, col := range stmt.Columns {
		if col.DataType == nil || col.Default != nil || col.Identity != nil || col.InlineConstraint != nil ||
			col.NotNull || col.Nullable || col.Collate != "" || col.Comment != nil || col.Tags != nil || col.VirtualExpr != nil || col.ProjectionPolicy != nil {
			return bad(col.Loc)
		}
	}
	if len(stmt.Constraints) > 0 {
		return bad(stmt.Constraints[0].Loc)
	}
	if len(stmt.Indexes) > 0 {
		return bad(stmt.Indexes[0].Loc)
	}
	return nil
}

// tryParseUsingTemplate parses
//
//	[ COPY GRANTS ] USING TEMPLATE ( <query> )
//
// when it follows the table name, and reports whether it did; otherwise it
// leaves the parser where it was. The page gives this form only OR REPLACE
// (https://docs.snowflake.com/en/sql-reference/sql/create-table).
func (p *Parser) tryParseUsingTemplate(stmt *ast.CreateTableStmt) (bool, error) {
	snap := p.snapshot()
	copyGrants := false
	if p.startsCopyGrants() {
		p.advance() // consume COPY
		p.advance() // consume GRANTS
		copyGrants = true
	}
	if !p.curIs(kwUSING) || !p.peekIsWord("TEMPLATE") {
		p.restore(snap)
		return false, nil
	}
	if stmt.OrAlter || stmt.Transient || stmt.Temporary || stmt.Volatile || stmt.Hybrid || stmt.IfNotExists {
		return false, p.syntaxErrorAtCur()
	}
	p.advance() // consume USING
	p.advance() // consume TEMPLATE
	if _, err := p.expect('('); err != nil {
		return false, err
	}
	query, err := p.parseQueryExpr()
	if err != nil {
		return false, err
	}
	if _, err := p.expect(')'); err != nil {
		return false, err
	}
	stmt.CopyGrants = copyGrants
	stmt.UsingTemplate = query
	return true, nil
}

// parseCreateTableAsQuery parses the optional AS <query> that ends a CTAS.
// An interactive table's every form puts CLUSTER BY before its AS
// (https://docs.snowflake.com/en/sql-reference/sql/create-interactive-table).
func (p *Parser) parseCreateTableAsQuery(stmt *ast.CreateTableStmt) error {
	if !p.curIs(kwAS) {
		return nil
	}
	if stmt.Interactive && len(stmt.ClusterBy) == 0 {
		return p.syntaxErrorAtCur()
	}
	p.advance() // consume AS
	query, err := p.parseQueryExpr()
	if err != nil {
		return err
	}
	stmt.AsSelect = query
	return nil
}

// parseCreateTableFrom parses the two forms of CREATE TABLE that restore a
// table (https://docs.snowflake.com/en/sql-reference/sql/create-table):
//
//	FROM ARCHIVE OF <table> [ [ AS ] <alias> ] WHERE <expr>
//	FROM BACKUP SET <backup_set> IDENTIFIER '<backup_id>'
//
// The page gives FROM BACKUP SET no modifier, and FROM ARCHIVE OF only
// TRANSIENT and IF NOT EXISTS. On entry cur is FROM.
func (p *Parser) parseCreateTableFrom(stmt *ast.CreateTableStmt) error {
	// Modifiers neither form takes.
	if stmt.OrReplace || stmt.OrAlter || stmt.Temporary || stmt.Volatile || stmt.Hybrid {
		return p.syntaxErrorAtCur()
	}
	p.advance() // consume FROM
	if p.curIsWord("ARCHIVE") {
		p.advance() // consume ARCHIVE
		if _, err := p.expect(kwOF); err != nil {
			return err
		}
		table, err := p.parseObjectName()
		if err != nil {
			return err
		}
		src := &ast.ArchiveSource{Table: table}
		if p.curIs(kwAS) {
			p.advance() // consume AS
			alias, err := p.parseIdent()
			if err != nil {
				return err
			}
			src.Alias = alias
		} else if !p.curIs(kwWHERE) {
			alias, err := p.parseIdent()
			if err != nil {
				return err
			}
			src.Alias = alias
		}
		if _, err := p.expect(kwWHERE); err != nil {
			return err
		}
		where, err := p.parseExpr()
		if err != nil {
			return err
		}
		src.Where = where
		stmt.FromArchive = src
		return nil
	}
	if stmt.Transient || stmt.IfNotExists {
		return p.syntaxErrorAtCur()
	}
	backup, err := p.parseBackupSource()
	if err != nil {
		return err
	}
	stmt.FromBackup = backup
	return nil
}

// parseBackupSource parses BACKUP SET <backup_set> IDENTIFIER '<backup_id>',
// which follows FROM in CREATE [ DYNAMIC [ ICEBERG ] ] TABLE. On entry cur is
// BACKUP, which is not a keyword token.
func (p *Parser) parseBackupSource() (*ast.BackupSource, error) {
	if !p.curWord("BACKUP") {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance() // consume BACKUP
	if _, err := p.expect(kwSET); err != nil {
		return nil, err
	}
	set, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwIDENTIFIER); err != nil {
		return nil, err
	}
	id, err := p.expect(tokString)
	if err != nil {
		return nil, err
	}
	return &ast.BackupSource{BackupSet: set, ID: id.Str}, nil
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

	for !p.curIs(')') && !p.curIs(tokEOF) {
		// An inline INDEX element (HYBRID TABLE) shares its leading keyword with a
		// column whose name is literally "index" (INDEX is non-reserved).
		// tryParseTableIndex speculatively matches the INDEX <name> ( ... ) shape
		// and reports matched=false (rolling back) when it is actually a column
		// named "index", which then falls through to parseColumnDef below.
		if p.curIs(kwINDEX) {
			idx, matched, err := p.tryParseTableIndex()
			if err != nil {
				return err
			}
			if matched {
				stmt.Indexes = append(stmt.Indexes, idx)
				if p.curIs(',') {
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
			col, err := p.parseColumnDefIn(stmt.Iceberg)
			if err != nil {
				return err
			}
			stmt.Columns = append(stmt.Columns, col)
		}

		if p.curIs(',') {
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
	switch {
	case p.curIs(kwCONSTRAINT), p.curIs(kwPRIMARY), p.curIs(kwUNIQUE), p.curIs(kwFOREIGN), p.curIs(kwCHECK):
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
	if nameErr != nil || !p.curIs('(') {
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

	idx = &ast.TableIndex{
		Name:    name,
		Columns: cols,
	}

	// Optional INCLUDE ( <col> [, <col>...] ): the covering columns
	// (https://docs.snowflake.com/en/sql-reference/sql/create-hybrid-table).
	// The element is committed by now, so an error here is a real one.
	if p.curIs(kwINCLUDE) {
		p.advance() // consume INCLUDE
		include, err := p.parseIdentListInParens()
		if err != nil {
			return nil, false, err
		}
		idx.Include = include
	}

	idx.Loc = ast.Loc{Start: start, End: p.prev.Loc.End}
	return idx, true, nil
}

// ---------------------------------------------------------------------------
// Column definition
// ---------------------------------------------------------------------------

// parseColumnDef parses a single column definition:
//
//	name [data_type | AS (expr)] [column_options...]
func (p *Parser) parseColumnDef() (*ast.ColumnDef, error) {
	return p.parseColumnDefIn(false)
}

// parseColumnDefIn parses a column definition; iceberg enables the bare
// `<col_name> <col_type> AS <expr>` virtual column form of Parquet Iceberg
// tables (https://docs.snowflake.com/en/sql-reference/sql/create-iceberg-table-parquet).
func (p *Parser) parseColumnDefIn(iceberg bool) (*ast.ColumnDef, error) {
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
	if !p.curIs(kwAS) && !p.curIs(',') && !p.curIs(')') &&
		!p.curIs(kwNOT) && !p.curIs(kwNULL) && !p.curIs(kwDEFAULT) &&
		!p.curIs(kwIDENTITY) && !p.curIs(kwAUTOINCREMENT) &&
		!p.curIs(kwCOLLATE) && !p.curIs(kwCONSTRAINT) &&
		!p.curIs(kwPRIMARY) && !p.curIs(kwUNIQUE) && !p.curIs(kwFOREIGN) &&
		!p.curIs(kwCOMMENT) && !p.curIs(kwWITH) && !p.curIs(kwMASKING) &&
		!p.curIs(kwTAG) && !p.curIs(tokEOF) {
		dt, err := p.parseDataType()
		if err != nil {
			return nil, err
		}
		col.DataType = dt
	}

	if err := p.parseColumnOptionsIn(col, iceberg); err != nil {
		return nil, err
	}

	col.Loc.End = p.prev.Loc.End
	return col, nil
}

// parseColumnOptions parses the optional constraints and properties that
// follow a column name+type. Returns nil when no more column options are
// found.
func (p *Parser) parseColumnOptions(col *ast.ColumnDef) error {
	return p.parseColumnOptionsIn(col, false)
}

func (p *Parser) parseColumnOptionsIn(col *ast.ColumnDef, iceberg bool) error {
	for {
		switch {
		case p.curIs(kwNOT):
			// NOT NULL
			if p.nextIs(kwNULL) {
				p.advance() // consume NOT
				p.advance() // consume NULL
				col.NotNull = true
			} else {
				return nil
			}

		case p.curIs(kwNULL):
			p.advance() // consume NULL
			col.Nullable = true

		case p.curIs(kwDEFAULT):
			p.advance() // consume DEFAULT
			expr, err := p.parseExpr()
			if err != nil {
				return err
			}
			col.Default = expr

		case p.curIs(kwIDENTITY), p.curIs(kwAUTOINCREMENT):
			spec, err := p.parseIdentitySpec()
			if err != nil {
				return err
			}
			col.Identity = spec

		case p.curIs(kwCOLLATE):
			p.advance() // consume COLLATE
			tok, err := p.expect(tokString)
			if err != nil {
				return err
			}
			col.Collate = tok.Str

		case p.curIs(kwCONSTRAINT), p.curIs(kwPRIMARY), p.curIs(kwUNIQUE), p.curIs(kwFOREIGN), p.curIs(kwREFERENCES), p.curIs(kwCHECK):
			ic, err := p.parseInlineConstraint(true)
			if err != nil {
				return err
			}
			col.InlineConstraint = ic

		case p.curIs(kwCOMMENT):
			p.advance() // consume COMMENT
			tok, err := p.expect(tokString)
			if err != nil {
				return err
			}
			s := tok.Str
			col.Comment = &s

		case p.curIs(kwWITH):
			// WITH MASKING POLICY name [USING (...)]
			// WITH TAG (...)
			// WITH ROW ACCESS POLICY ... — consumed and discarded
			switch {
			case p.nextIs(kwMASKING):
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
				// Optional USING ( <col> [ , ... ] )
				if p.curIs(kwUSING) {
					p.advance() // consume USING
					if _, err := p.parseIdentListInParens(); err != nil {
						return err
					}
				}
			case p.nextIs(kwTAG):
				p.advance() // consume WITH
				tags, err := p.parseTagAssignments()
				if err != nil {
					return err
				}
				col.Tags = append(col.Tags, tags...)
			case p.nextIs(tokIdent):
				// WITH PROJECTION POLICY: nothing else begins WITH PROJECTION.
				if !p.peekIsWord("PROJECTION") {
					return nil
				}
				p.advance() // consume WITH
				policyName, err := p.parseProjectionPolicy()
				if err != nil {
					return err
				}
				col.ProjectionPolicy = policyName
			case p.nextIs(kwROW):
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
				if p.curIs(kwON) {
					p.advance() // consume ON
					if err := p.skipParenthesized(); err != nil {
						return err
					}
				}
			default:
				return nil
			}

		case p.curIs(kwMASKING):
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
			if p.curIs(kwUSING) {
				p.advance() // consume USING
				if _, err := p.parseIdentListInParens(); err != nil {
					return err
				}
			}

		case p.curIs(kwTAG):
			// TAG (...) — without WITH prefix
			tags, err := p.parseTagAssignments()
			if err != nil {
				return err
			}
			col.Tags = append(col.Tags, tags...)

		case p.curIs(kwAS):
			if col.VirtualExpr != nil {
				return p.syntaxErrorAtCur()
			}
			if iceberg && col.DataType != nil {
				p.advance() // consume AS
				expr, err := p.parseExpr()
				if err != nil {
					return err
				}
				if pe, ok := expr.(*ast.ParenExpr); ok {
					expr = pe.Expr
				}
				col.VirtualExpr = expr
				continue
			}
			if err := p.parseVirtualColumnExpr(col); err != nil {
				return err
			}

		default:
			switch {
			case !iceberg && p.startsGeneratedAlways():
				if col.VirtualExpr != nil {
					return p.syntaxErrorAtCur()
				}
				p.advance() // consume GENERATED
				p.advance() // consume ALWAYS
				if !p.curIs(kwAS) {
					return p.syntaxErrorAtCur()
				}
				if err := p.parseVirtualColumnExpr(col); err != nil {
					return err
				}
			case p.startsProjectionPolicy():
				policyName, err := p.parseProjectionPolicy()
				if err != nil {
					return err
				}
				col.ProjectionPolicy = policyName
			default:
				return nil
			}
		}
	}
}

// parseVirtualColumnExpr parses a virtual column's expression, with the
// optional GENERATED ALWAYS already consumed
// (https://docs.snowflake.com/en/sql-reference/sql/create-table):
//
//	AS ( <expr> ) [ VIRTUAL ]
//
// GENERATED ALWAYS and VIRTUAL restate what AS ( <expr> ) already means, so
// the AST keeps only the expression.
func (p *Parser) parseVirtualColumnExpr(col *ast.ColumnDef) error {
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
	if p.curIsWord("VIRTUAL") {
		p.advance() // consume VIRTUAL
	}
	return nil
}

// startsGeneratedAlways reports whether cur begins GENERATED ALWAYS, the
// optional lead-in of a virtual column's AS ( <expr> ). Neither word is a
// keyword token.
func (p *Parser) startsGeneratedAlways() bool {
	return p.curIsWord("GENERATED") && p.peekIsWord("ALWAYS")
}

// startsProjectionPolicy reports whether cur begins PROJECTION POLICY.
// PROJECTION is not a keyword token.
func (p *Parser) startsProjectionPolicy() bool {
	return p.curIsWord("PROJECTION") && p.nextIs(kwPOLICY)
}

// parseProjectionPolicy parses PROJECTION POLICY <name>, a column's
// projection policy (https://docs.snowflake.com/en/sql-reference/sql/create-table).
// On entry cur is PROJECTION; a leading WITH is already consumed.
func (p *Parser) parseProjectionPolicy() (*ast.ObjectName, error) {
	p.advance() // consume PROJECTION
	if _, err := p.expect(kwPOLICY); err != nil {
		return nil, err
	}
	return p.parseObjectName()
}

// ---------------------------------------------------------------------------
// Inline constraint
// ---------------------------------------------------------------------------

// parseInlineConstraint parses a column-level (inline) constraint
// (https://docs.snowflake.com/en/sql-reference/sql/create-table-constraint):
//
//	[CONSTRAINT name] { PRIMARY KEY | UNIQUE | [FOREIGN KEY] REFERENCES ... | CHECK ( expr ) }
//	  [constraint_properties]
//
// allowCheck is false where the object's page lists no CHECK constraint.
func (p *Parser) parseInlineConstraint(allowCheck bool) (*ast.InlineConstraint, error) {
	ic := &ast.InlineConstraint{
		Loc: p.cur.Loc,
	}

	// Optional CONSTRAINT name
	if p.curIs(kwCONSTRAINT) {
		p.advance() // consume CONSTRAINT
		name, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		ic.Name = name
	}

	switch {
	case p.curIs(kwPRIMARY):
		p.advance() // consume PRIMARY
		if _, err := p.expect(kwKEY); err != nil {
			return nil, err
		}
		ic.Type = ast.ConstrPrimaryKey

	case p.curIs(kwUNIQUE):
		p.advance() // consume UNIQUE
		ic.Type = ast.ConstrUnique

	case p.curIs(kwFOREIGN), p.curIs(kwREFERENCES):
		// An inline foreign key's FOREIGN KEY is optional.
		if p.curIs(kwFOREIGN) {
			p.advance() // consume FOREIGN
			if _, err := p.expect(kwKEY); err != nil {
				return nil, err
			}
		}
		ref, err := p.parseForeignKeyRef()
		if err != nil {
			return nil, err
		}
		ic.Type = ast.ConstrForeignKey
		ic.References = ref

	case p.curIs(kwCHECK):
		if !allowCheck {
			return nil, p.syntaxErrorAtCur()
		}
		expr, err := p.parseCheckExpr()
		if err != nil {
			return nil, err
		}
		ic.Type = ast.ConstrCheck
		ic.Check = expr

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

// parseOutOfLineConstraint parses a table-level (out-of-line) constraint
// (https://docs.snowflake.com/en/sql-reference/sql/create-table-constraint):
//
//	[CONSTRAINT name] { PRIMARY KEY (cols) | UNIQUE (cols)
//	  | FOREIGN KEY (cols) REFERENCES ... | CHECK ( expr ) }
//	  [constraint_properties] [COMMENT 'text']
//
// FOREIGN KEY is required here: the constraint page writes the out-of-line
// form without brackets around it.
func (p *Parser) parseOutOfLineConstraint() (*ast.TableConstraint, error) {
	con := &ast.TableConstraint{
		Loc: p.cur.Loc,
	}

	// Optional CONSTRAINT name
	if p.curIs(kwCONSTRAINT) {
		p.advance() // consume CONSTRAINT
		name, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		con.Name = name
	}

	switch {
	case p.curIs(kwPRIMARY):
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

	case p.curIs(kwUNIQUE):
		p.advance() // consume UNIQUE
		con.Type = ast.ConstrUnique
		cols, err := p.parseIdentListInParens()
		if err != nil {
			return nil, err
		}
		con.Columns = cols

	case p.curIs(kwFOREIGN):
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

	case p.curIs(kwCHECK):
		expr, err := p.parseCheckExpr()
		if err != nil {
			return nil, err
		}
		con.Type = ast.ConstrCheck
		con.Check = expr

	default:
		return nil, p.syntaxErrorAtCur()
	}

	if err := p.parseConstraintProperties(); err != nil {
		return nil, err
	}

	// Optional COMMENT 'text'
	if p.curIs(kwCOMMENT) {
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

// parseCheckExpr parses CHECK ( <expr> ) and returns the expression. On entry
// cur is CHECK.
func (p *Parser) parseCheckExpr() (ast.Node, error) {
	p.advance() // consume CHECK
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return expr, nil
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

	for p.curIs(',') {
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
//
// The constraint page writes the actions under one ON, as
// ON [ UPDATE action ] [ DELETE action ]
// (https://docs.snowflake.com/en/sql-reference/sql/create-table-constraint),
// so a DELETE action may also follow an UPDATE action without its own ON.
func (p *Parser) parseForeignKeyRefAfterReferences() (*ast.ForeignKeyRef, error) {
	table, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}

	ref := &ast.ForeignKeyRef{
		Table: table,
	}

	// Optional column list
	if p.curIs('(') {
		cols, err := p.parseIdentListInParens()
		if err != nil {
			return nil, err
		}
		ref.Columns = cols
	}

	// Optional MATCH FULL | PARTIAL | SIMPLE
	if p.curIs(kwMATCH) {
		p.advance() // consume MATCH
		switch {
		case p.curIs(kwFULL):
			p.advance()
			ref.Match = "FULL"
		case p.curIs(kwPARTIAL):
			p.advance()
			ref.Match = "PARTIAL"
		case p.curIs(kwSIMPLE):
			p.advance()
			ref.Match = "SIMPLE"
		default:
			return nil, p.syntaxErrorAtCur()
		}
	}

	// Optional ON DELETE / ON UPDATE
	for p.curIs(kwON) {
		switch {
		case p.nextIs(kwDELETE):
			p.advance() // consume ON
			p.advance() // consume DELETE
			action, err := p.parseReferenceAction()
			if err != nil {
				return nil, err
			}
			ref.OnDelete = action
		case p.nextIs(kwUPDATE):
			p.advance() // consume ON
			p.advance() // consume UPDATE
			action, err := p.parseReferenceAction()
			if err != nil {
				return nil, err
			}
			ref.OnUpdate = action
			if p.curIs(kwDELETE) {
				p.advance() // consume DELETE
				action, err := p.parseReferenceAction()
				if err != nil {
					return nil, err
				}
				ref.OnDelete = action
			}
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
	switch {
	case p.curIs(kwCASCADE):
		p.advance()
		return ast.RefActCascade, nil
	case p.curIs(kwSET):
		p.advance() // consume SET
		switch {
		case p.curIs(kwNULL):
			p.advance()
			return ast.RefActSetNull, nil
		case p.curIs(kwDEFAULT):
			p.advance()
			return ast.RefActSetDefault, nil
		default:
			return ast.RefActNone, p.syntaxErrorAtCur()
		}
	case p.curIs(kwRESTRICT):
		p.advance()
		return ast.RefActRestrict, nil
	case p.curIs(kwNO):
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
		switch {
		case p.curIs(kwNOT):
			if p.nextIs(kwENFORCED) {
				p.advance() // consume NOT
				p.advance() // consume ENFORCED
			} else if p.nextIs(kwDEFERRABLE) {
				p.advance() // consume NOT
				p.advance() // consume DEFERRABLE
			} else {
				return nil
			}
		case p.curIs(kwENFORCED):
			p.advance()
		case p.curIs(kwDEFERRABLE):
			p.advance()
		case p.curIs(kwINITIALLY):
			p.advance() // consume INITIALLY
			// Expect DEFERRED or IMMEDIATE
			if p.curIs(kwDEFERRED) || p.curIs(kwIMMEDIATE) {
				p.advance()
			} else {
				return p.syntaxErrorAtCur()
			}
		case p.curIs(kwENABLE):
			p.advance() // consume ENABLE
			// Optional VALIDATE / NOVALIDATE
			if p.curIs(kwVALIDATE) || p.curIs(kwNOVALIDATE) {
				p.advance()
			}
		case p.curIs(kwDISABLE):
			p.advance() // consume DISABLE
			// Optional VALIDATE / NOVALIDATE
			if p.curIs(kwVALIDATE) || p.curIs(kwNOVALIDATE) {
				p.advance()
			}
		case p.curIs(kwVALIDATE):
			p.advance()
		case p.curIs(kwNOVALIDATE):
			p.advance()
		case p.curIs(kwRELY):
			p.advance()
		case p.curIs(kwNORELY):
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
	if p.curIs('(') {
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
	if p.curIs(kwSTART) {
		p.advance() // consume START
		// Optional WITH
		if p.curIs(kwWITH) {
			p.advance() // consume WITH
		}
		// Optional =
		if p.curIs('=') {
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
	if p.curIs(kwINCREMENT) {
		p.advance() // consume INCREMENT
		// Optional BY
		if p.curIs(kwBY) {
			p.advance() // consume BY
		}
		// Optional =
		if p.curIs('=') {
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
	if p.curIs(kwORDER) {
		p.advance()
		b := true
		spec.Order = &b
	} else if p.curIs(kwNOORDER) {
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

	for !p.curIs(')') && !p.curIs(tokEOF) {
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

		if p.curIs(',') {
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
	if p.curIs(kwAT) || p.curIs(kwBEFORE) {
		atBefore := "AT"
		if p.curIs(kwBEFORE) {
			atBefore = "BEFORE"
		}
		p.advance() // consume AT or BEFORE

		if _, err := p.expect('('); err != nil {
			return nil, err
		}

		var kind string
		switch {
		case p.curIs(kwTIMESTAMP):
			p.advance()
			kind = "TIMESTAMP"
		case p.curIs(kwOFFSET):
			p.advance()
			kind = "OFFSET"
		case p.curIs(kwSTATEMENT):
			p.advance()
			kind = "STATEMENT"
		default:
			return nil, p.syntaxErrorAtCur()
		}

		// Expect => (tokAssoc)
		if _, err := p.expect(tokAssoc); err != nil {
			return nil, err
		}

		// Value: any expression. The docs accept a string literal for STATEMENT,
		// but TIMESTAMP/OFFSET commonly take a function call (e.g.
		// TO_TIMESTAMP_TZ('…')) or an arithmetic expression, so we parse the full
		// expression grammar here rather than only a bare string token.
		value, err := p.parseExpr()
		if err != nil {
			return nil, err
		}

		if _, err := p.expect(')'); err != nil {
			return nil, err
		}

		clone.AtBefore = atBefore
		clone.Kind = kind
		clone.Value = value
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
	// The governance clauses each page lists: CREATE TABLE's
	// (https://docs.snowflake.com/en/sql-reference/sql/create-table) and an
	// interactive table's, which has no WITH CONTACT.
	clauses := clauseRowAccess | clauseAggregation | clauseJoin | clauseStorageLifecycle | clauseContact
	if stmt.Interactive {
		clauses &^= clauseContact
	}
	for {
		if c := p.startsObjectClause(clauses); c != 0 {
			res, err := p.parseObjectClause(c)
			if err != nil {
				return err
			}
			if res.RowPolicy != nil {
				stmt.RowPolicy = res.RowPolicy
			}
			if res.AggregationPolicy != nil {
				stmt.AggregationPolicy = res.AggregationPolicy
			}
			stmt.Contacts = append(stmt.Contacts, res.Contacts...)
			continue
		}
		// An interactive table's refresh settings
		// (https://docs.snowflake.com/en/sql-reference/sql/create-interactive-table).
		if stmt.Interactive && (p.curIs(kwTARGET_LAG) || p.curIs(kwWAREHOUSE) || p.curIsWord("INITIALIZATION_WAREHOUSE")) {
			opt, err := p.parseRefreshSetting(false)
			if err != nil {
				return err
			}
			stmt.Options = append(stmt.Options, opt)
			continue
		}
		// Of the rest, an interactive table takes only CLUSTER BY,
		// COMMENT = '<string>' and [ WITH ] TAG ( ... ).
		if stmt.Interactive {
			switch {
			case p.curIs(kwCLUSTER), p.curIs(kwTAG):
			case p.curIs(kwWITH) && p.nextIs(kwTAG):
			case p.curIs(kwCOMMENT) && p.nextIs('='):
			default:
				return nil
			}
		}
		switch {
		case p.curIs(kwCLUSTER):
			// CLUSTER BY [LINEAR] (exprs)
			p.advance() // consume CLUSTER
			if _, err := p.expect(kwBY); err != nil {
				return err
			}
			if p.curIs(kwLINEAR) {
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

		case p.curIs(kwCOPY):
			// COPY GRANTS | COPY TAGS | COPY = ( ... )
			if !p.startsCopyClause() {
				return nil
			}
			if err := p.parseTableCopyClause(stmt); err != nil {
				return err
			}

		case p.curIs(kwCOMMENT):
			// COMMENT = 'text'
			p.advance() // consume COMMENT
			if p.curIs('=') {
				p.advance() // consume =
			}
			tok, err := p.expect(tokString)
			if err != nil {
				return err
			}
			s := tok.Str
			stmt.Comment = &s

		case p.curIs(kwWITH):
			switch {
			case p.nextIs(kwTAG):
				p.advance() // consume WITH
				tags, err := p.parseTagAssignments()
				if err != nil {
					return err
				}
				stmt.Tags = append(stmt.Tags, tags...)
			default:
				return nil
			}

		case p.curIs(kwTAG):
			// TAG (...) — without WITH prefix
			tags, err := p.parseTagAssignments()
			if err != nil {
				return err
			}
			stmt.Tags = append(stmt.Tags, tags...)

		case p.curIs(kwDATA_RETENTION_TIME_IN_DAYS):
			// DATA_RETENTION_TIME_IN_DAYS = n — consume and discard
			p.advance() // consume DATA_RETENTION_TIME_IN_DAYS
			if p.curIs('=') {
				p.advance()
			}
			if _, err := p.expect(tokInt); err != nil {
				return err
			}

		case p.curIs(kwMAX_DATA_EXTENSION_TIME_IN_DAYS):
			// MAX_DATA_EXTENSION_TIME_IN_DAYS = n — consume and discard
			p.advance() // consume MAX_DATA_EXTENSION_TIME_IN_DAYS
			if _, err := p.expect('='); err != nil {
				return err
			}
			if _, err := p.expect(tokInt); err != nil {
				return err
			}

		case p.curIs(kwCHANGE_TRACKING):
			// CHANGE_TRACKING = TRUE|FALSE — consume and discard
			p.advance() // consume CHANGE_TRACKING
			if p.curIs('=') {
				p.advance()
			}
			if p.curIs(kwTRUE) || p.curIs(kwFALSE) {
				p.advance()
			} else {
				return p.syntaxErrorAtCur()
			}

		case p.curIs(kwDEFAULT_DDL_COLLATION):
			// DEFAULT_DDL_COLLATION = 'string' — consume and discard
			p.advance() // consume DEFAULT_DDL_COLLATION
			if p.curIs('=') {
				p.advance()
			}
			if _, err := p.expect(tokString); err != nil {
				return err
			}

		case p.curIs(tokIdent):
			// The parameters CREATE TABLE lists that are not keyword tokens —
			// consume and discard, like the ones above:
			//	ICEBERG_DEFAULT_DDL_COLLATION = '<string>'
			//	{ ENABLE_SCHEMA_EVOLUTION | ERROR_LOGGING | ROW_TIMESTAMP } = { TRUE | FALSE }
			switch {
			case p.curIsWord("ICEBERG_DEFAULT_DDL_COLLATION"):
				p.advance() // consume the parameter name
				if _, err := p.expect('='); err != nil {
					return err
				}
				if _, err := p.expect(tokString); err != nil {
					return err
				}
			case p.curIsWord("ENABLE_SCHEMA_EVOLUTION") || p.curIsWord("ERROR_LOGGING") || p.curIsWord("ROW_TIMESTAMP"):
				p.advance() // consume the parameter name
				if _, err := p.expect('='); err != nil {
					return err
				}
				if !p.curIs(kwTRUE) && !p.curIs(kwFALSE) {
					return p.syntaxErrorAtCur()
				}
				p.advance() // consume TRUE / FALSE
			default:
				return nil
			}

		case p.curIs(kwSTAGE_FILE_FORMAT):
			// STAGE_FILE_FORMAT = (...) — consume and discard
			p.advance() // consume STAGE_FILE_FORMAT
			if p.curIs('=') {
				p.advance()
			}
			if err := p.skipParenthesized(); err != nil {
				return err
			}

		case p.curIs(kwSTAGE_COPY_OPTIONS):
			// STAGE_COPY_OPTIONS = (...) — consume and discard
			p.advance() // consume STAGE_COPY_OPTIONS
			if p.curIs('=') {
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
	for depth > 0 && !p.curIs(tokEOF) {
		switch {
		case p.curIs('('):
			depth++
		case p.curIs(')'):
			depth--
		}
		p.advance()
	}
	if depth != 0 {
		return p.syntaxErrorAtCur()
	}
	return nil
}

// ---------------------------------------------------------------------------
// Object-level clauses shared by the table-like objects
// ---------------------------------------------------------------------------

// objectClause names one of the object-level governance clauses that CREATE
// TABLE, the table variants, and CREATE [MATERIALIZED] VIEW list. Each page
// lists its own subset, so a caller passes the clauses its page allows.
type objectClause uint8

const (
	// [ WITH ] ROW ACCESS POLICY <name> ON ( <col> [ , ... ] )
	clauseRowAccess objectClause = 1 << iota
	// [ WITH ] AGGREGATION POLICY <name> [ ENTITY KEY ( <col> [ , ... ] ) ]
	clauseAggregation
	// [ WITH ] JOIN POLICY <name> [ ALLOWED JOIN KEYS ( <col> [ , ... ] ) ]
	clauseJoin
	// [ WITH ] STORAGE LIFECYCLE POLICY <name> ON ( <col> [ , ... ] )
	clauseStorageLifecycle
	// WITH CONTACT ( <purpose> = <contact> [ , ... ] )
	clauseContact
)

// startsObjectClause reports which allowed object clause begins at cur, with
// or without a leading WITH, or 0 when none does. Past a WITH only one token is
// visible, and each of these clauses is the only thing its second word can
// begin there.
func (p *Parser) startsObjectClause(allowed objectClause) objectClause {
	var c objectClause
	if p.curIs(kwWITH) {
		switch {
		case p.nextIs(kwROW):
			c = clauseRowAccess
		case p.nextIs(kwJOIN):
			c = clauseJoin
		case p.nextIs(kwSTORAGE):
			c = clauseStorageLifecycle
		case p.peekIsWord("AGGREGATION"):
			c = clauseAggregation
		case p.peekIsWord("CONTACT"):
			c = clauseContact
		}
	} else {
		switch {
		case p.curIs(kwROW) && p.nextIs(kwACCESS):
			c = clauseRowAccess
		case p.curIs(kwJOIN) && p.nextIs(kwPOLICY):
			c = clauseJoin
		case p.curIs(kwSTORAGE) && p.peekIsWord("LIFECYCLE"):
			c = clauseStorageLifecycle
		case p.curIsWord("AGGREGATION") && p.nextIs(kwPOLICY):
			c = clauseAggregation
		}
	}
	return c & allowed
}

// objectClauseResult is what parseObjectClause keeps of a clause, for the
// objects whose AST holds it: the row access policy, the aggregation policy's
// name, and the contacts. The rest of each clause is read and not kept.
type objectClauseResult struct {
	RowPolicy         *ast.RowAccessPolicy
	AggregationPolicy *ast.ObjectName
	Contacts          []*ast.CopyOption
}

// parseObjectClause parses the clause startsObjectClause reported.
func (p *Parser) parseObjectClause(c objectClause) (objectClauseResult, error) {
	var res objectClauseResult
	if p.curIs(kwWITH) {
		p.advance() // consume WITH
	}
	switch c {
	case clauseRowAccess:
		rp, err := p.parseRowAccessPolicyClause()
		if err != nil {
			return res, err
		}
		res.RowPolicy = rp
	case clauseAggregation:
		p.advance() // consume AGGREGATION
		name, err := p.parsePolicyName()
		if err != nil {
			return res, err
		}
		res.AggregationPolicy = name
		if p.curIsWord("ENTITY") {
			p.advance() // consume ENTITY
			if _, err := p.expect(kwKEY); err != nil {
				return res, err
			}
			if _, err := p.parseIdentListInParens(); err != nil {
				return res, err
			}
		}
	case clauseJoin:
		p.advance() // consume JOIN
		if _, err := p.parsePolicyName(); err != nil {
			return res, err
		}
		if p.curIsWord("ALLOWED") {
			p.advance() // consume ALLOWED
			if _, err := p.expect(kwJOIN); err != nil {
				return res, err
			}
			if _, err := p.expect(kwKEYS); err != nil {
				return res, err
			}
			if _, err := p.parseIdentListInParens(); err != nil {
				return res, err
			}
		}
	case clauseStorageLifecycle:
		p.advance() // consume STORAGE
		if !p.curIsWord("LIFECYCLE") {
			return res, p.syntaxErrorAtCur()
		}
		p.advance() // consume LIFECYCLE
		if _, err := p.parsePolicyName(); err != nil {
			return res, err
		}
		if _, err := p.expect(kwON); err != nil {
			return res, err
		}
		if _, err := p.parseIdentListInParens(); err != nil {
			return res, err
		}
	case clauseContact:
		contacts, err := p.parseWithContact()
		if err != nil {
			return res, err
		}
		res.Contacts = contacts
	default:
		return res, p.syntaxErrorAtCur()
	}
	return res, nil
}

// parseRefreshSetting parses one of the refresh settings an interactive or a
// cloned dynamic table takes
// (https://docs.snowflake.com/en/sql-reference/sql/create-interactive-table,
// https://docs.snowflake.com/en/sql-reference/sql/create-dynamic-table):
//
//	TARGET_LAG = { '<num> { seconds | minutes | hours | days }' | DOWNSTREAM }
//	{ WAREHOUSE | INITIALIZATION_WAREHOUSE } = <warehouse_name>
//
// DOWNSTREAM is a dynamic table's only; downstream says whether it is allowed.
// A warehouse name is an identifier, so a quoted one keeps its quotes in
// Words, as the option's text.
func (p *Parser) parseRefreshSetting(downstream bool) (*ast.CopyOption, error) {
	nameTok := p.advance() // consume the setting name
	opt := &ast.CopyOption{
		Name: strings.ToUpper(nameTok.Str),
		Loc:  ast.Loc{Start: nameTok.Loc.Start},
	}
	if _, err := p.expect('='); err != nil {
		return nil, err
	}
	switch {
	case opt.Name != "TARGET_LAG":
		name, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		opt.Words = name.String()
	case p.curIs(tokString):
		tok := p.advance()
		opt.Lit = &ast.Literal{Kind: ast.LitString, Value: tok.Str, Loc: tok.Loc}
	case downstream && p.curIsWord("DOWNSTREAM"):
		p.advance() // consume DOWNSTREAM
		opt.Words = "DOWNSTREAM"
	default:
		return nil, p.syntaxErrorAtCur()
	}
	opt.Loc.End = p.prev.Loc.End
	return opt, nil
}

// parsePolicyName parses POLICY <name> and returns the name.
func (p *Parser) parsePolicyName() (*ast.ObjectName, error) {
	if _, err := p.expect(kwPOLICY); err != nil {
		return nil, err
	}
	return p.parseObjectName()
}

// parseTableCopyClause parses CREATE TABLE's COPY GRANTS, COPY TAGS or
// COPY = ( ... ) into the fields CREATE ICEBERG TABLE also sets. The caller
// has confirmed startsCopyClause.
func (p *Parser) parseTableCopyClause(stmt *ast.CreateTableStmt) error {
	switch {
	case p.nextIs(kwGRANTS):
		p.advance() // consume COPY
		p.advance() // consume GRANTS
		stmt.CopyGrants = true
	case p.nextIs(kwTAGS):
		p.advance() // consume COPY
		p.advance() // consume TAGS
		stmt.CopyTags = true
	default:
		opt, err := p.parseCopyWhatOption()
		if err != nil {
			return err
		}
		stmt.Options = append(stmt.Options, opt)
	}
	return nil
}

// startsCopyClause reports whether cur begins COPY GRANTS, COPY TAGS, or
// COPY = ( ... ).
func (p *Parser) startsCopyClause() bool {
	if !p.curIs(kwCOPY) {
		return false
	}
	switch {
	case p.nextIs(kwGRANTS), p.nextIs(kwTAGS), p.nextIs('='):
		return true
	}
	return false
}

// parseCopyClause parses one of
//
//	COPY GRANTS
//	COPY TAGS
//	COPY = ( { COMMENTS | TAGS | GRANTS } [ , ... ] )
//
// and reports whether it copies grants. grantsOnly is true for the objects
// whose page lists only COPY GRANTS and COPY = ( GRANTS ) (event table,
// external table, materialized view). What else is copied is not kept.
func (p *Parser) parseCopyClause(grantsOnly bool) (bool, error) {
	p.advance() // consume COPY
	switch {
	case p.curIs(kwGRANTS):
		p.advance() // consume GRANTS
		return true, nil
	case p.curIs(kwTAGS):
		if grantsOnly {
			return false, p.syntaxErrorAtCur()
		}
		p.advance() // consume TAGS
		return false, nil
	}
	if _, err := p.expect('='); err != nil {
		return false, err
	}
	if _, err := p.expect('('); err != nil {
		return false, err
	}
	grants := false
	for {
		switch {
		case p.curIs(kwGRANTS):
			grants = true
		case !grantsOnly && (p.curIs(kwTAGS) || p.curIsWord("COMMENTS")):
		default:
			return false, p.syntaxErrorAtCur()
		}
		p.advance() // consume the copied item
		if !p.curIs(',') {
			break
		}
		p.advance() // consume ','
	}
	if _, err := p.expect(')'); err != nil {
		return false, err
	}
	return grants, nil
}
