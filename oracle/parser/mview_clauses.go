package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// This file reads CREATE MATERIALIZED VIEW and ALTER MATERIALIZED VIEW by
// the grammar of their pages, reusing CREATE TABLE's readers for the
// physical properties, column properties and partitioning a materialized
// view shares with a table.
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-MATERIALIZED-VIEW.html
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-MATERIALIZED-VIEW.html
// BNF: oracle/parser/bnf/CREATE-MATERIALIZED-VIEW.bnf, oracle/parser/bnf/ALTER-MATERIALIZED-VIEW.bnf

// parseCreateMaterializedViewStmt reads a CREATE MATERIALIZED VIEW after
// its VIEW.
//
//	CREATE MATERIALIZED VIEW [ IF NOT EXISTS ] [ schema. ] materialized_view
//	    [ OF [ schema. ] object_type
//	    | ( { scoped_table_ref_constraint | column_alias [ ENCRYPT [ encryption_spec ] ] } [, ...] )
//	    ]
//	    [ DEFAULT COLLATION collation_name ]
//	    { ON PREBUILT TABLE [ { WITH | WITHOUT } REDUCED PRECISION ]
//	    | [ physical_properties ] [ materialized_view_props ]
//	    }
//	    [ USING INDEX [ physical_attributes_clause | TABLESPACE tablespace ] ... | USING NO INDEX ]
//	    [ create_mv_refresh ]
//	    [ FOR UPDATE ]
//	    [ evaluation_edition_clause ]
//	    [ { ENABLE | DISABLE } ON QUERY COMPUTATION ]
//	    [ query_rewrite_clause ]
//	    [ unusable_editions_clause ]
//	    [ annotations_clause ]
//	AS subquery
//
// The order is the pinned server's, which differs from the page's in places:
// OF and the alias list exclude each other (ORA-00905 at OF t(a), ORA-02000
// at (a) OF t); the alias list holds scoped_table_ref_constraint, which
// Oracle refuses on its own (ORA-02000 at SCOPE); the query rewrite and
// query computation clauses come in either order; the annotations_clause
// stands last before AS (after AS the query takes ANNOTATIONS as an alias,
// ORA-03048); and USING INDEX, FOR UPDATE and the unusable editions clause,
// left out of the page, stand where the 19c page has them.
func (p *Parser) parseCreateMaterializedViewStmt(stmt *nodes.CreateViewStmt) (*nodes.CreateViewStmt, error) {
	if p.cur.Type == kwIF {
		if p.peekNext().Type == kwEXISTS {
			return nil, p.syntaxErrorAtCur()
		}
		if p.peekNext().Type == kwNOT {
			// Oracle answers IF NOT without EXISTS at the IF (ORA-11543).
			ifTok := p.cur
			p.advance() // consume IF
			p.advance() // consume NOT
			if p.cur.Type != kwEXISTS {
				return nil, p.syntaxErrorAtTok(ifTok)
			}
			p.advance() // consume EXISTS
			stmt.IfNotExists = true
		}
	}

	var err error
	stmt.Name, err = p.parseReservedCheckedObjectName()
	if err != nil {
		return nil, err
	}
	if stmt.Name == nil || stmt.Name.Name == "" {
		return nil, p.syntaxErrorAtCur()
	}

	switch {
	case p.cur.Type == kwOF:
		p.advance() // consume OF
		objectType, err := p.parseRequiredObjectName()
		if err != nil {
			return nil, err
		}
		_ = objectType // the AST keeps no object type
	case p.cur.Type == '(':
		if err := p.parseMviewColumnAliases(stmt); err != nil {
			return nil, err
		}
	}

	if p.cur.Type == kwDEFAULT {
		p.advance() // consume DEFAULT
		if err := p.expectCreateTableWord("COLLATION"); err != nil {
			return nil, err
		}
		if !p.isCreateTableNameStart() {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.DefaultCollation, err = p.parseIdentifier()
		if err != nil {
			return nil, err
		}
	}

	defer func(saved bool) { p.mviewInmemoryColumn = saved }(p.mviewInmemoryColumn)
	p.mviewInmemoryColumn = false
	if p.cur.Type == kwON && isCreateTableWordTok(p.peekNext(), "PREBUILT") {
		if err := p.parseMviewPrebuiltTable(stmt); err != nil {
			return nil, err
		}
	} else if err := p.parseMviewProperties(stmt); err != nil {
		return nil, err
	}

	if p.cur.Type == kwUSING {
		switch next := p.peekNext(); {
		case next.Type == kwINDEX:
			p.advance() // consume USING
			p.advance() // consume INDEX
			if err := p.parseMviewUsingIndex(); err != nil {
				return nil, err
			}
		case isCreateTableWordTok(next, "NO"):
			p.advance() // consume USING
			p.advance() // consume NO
			if err := p.expectCreateTableToken(kwINDEX); err != nil {
				return nil, err
			}
		default:
			p.advance() // consume USING
			return nil, p.syntaxErrorAtCur()
		}
	}
	if p.cur.Type == kwON {
		p.advance() // consume ON (expected ON PREBUILT TABLE)
		return nil, p.syntaxErrorAtCur()
	}

	switch {
	case p.cur.Type == kwREFRESH:
		if err := p.parseMviewCreateRefresh(stmt); err != nil {
			return nil, err
		}
	case p.isKeywordStr("NEVER"):
		p.advance() // consume NEVER
		if err := p.expectCreateTableToken(kwREFRESH); err != nil {
			return nil, err
		}
		stmt.NeverRefresh = true
	}

	if p.cur.Type == kwFOR {
		p.advance() // consume FOR
		if err := p.expectCreateTableToken(kwUPDATE); err != nil {
			return nil, err
		}
	}

	if p.isKeywordStr("EVALUATE") {
		if err := p.parseMviewEvaluationEdition(); err != nil {
			return nil, err
		}
	}

	rewrite, computation := false, false
	for p.cur.Type == kwENABLE || p.cur.Type == kwDISABLE {
		enable := p.cur.Type == kwENABLE
		switch next := p.peekNext(); {
		case !rewrite && isCreateTableWordTok(next, "QUERY"):
			p.advance() // consume ENABLE / DISABLE
			p.advance() // consume QUERY
			if err := p.expectCreateTableToken(kwREWRITE); err != nil {
				return nil, err
			}
			stmt.EnableQuery, stmt.DisableQuery = enable, !enable
			rewrite = true
		case !computation && next.Type == kwON:
			p.advance() // consume ENABLE / DISABLE
			p.advance() // consume ON
			if err := p.expectCreateTableWord("QUERY"); err != nil {
				return nil, err
			}
			if err := p.expectCreateTableWord("COMPUTATION"); err != nil {
				return nil, err
			}
			stmt.EnableOnQueryComputation, stmt.DisableOnQueryComputation = enable, !enable
			computation = true
		default:
			return nil, p.syntaxErrorAtTok(next)
		}
	}
	if rewrite || computation {
		if err := p.parseMviewUnusableEditions(); err != nil {
			return nil, err
		}
	}

	if p.isKeywordStr("ANNOTATIONS") {
		if err := p.parseMviewAnnotations(false); err != nil {
			return nil, err
		}
	}

	if err := p.expectCreateTableToken(kwAS); err != nil {
		return nil, err
	}
	stmt.Query, err = p.parseSelectStmt()
	if err != nil {
		return nil, err
	}
	if !selectStmtHasTargets(stmt.Query) {
		return nil, p.syntaxErrorAtCur()
	}
	if p.mviewInmemoryColumn {
		// Oracle 23ai refuses column-level INMEMORY clauses on CREATE
		// MATERIALIZED VIEW with ORA-00922 at the final token of the query.
		return nil, p.syntaxErrorAtTok(p.prev)
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseMviewColumnAliases reads the parenthesized alias list of a
// materialized view.
//
//	( { scoped_table_ref_constraint | column_alias [ ENCRYPT [ encryption_spec ] ] } [, ...] )
//
//	scoped_table_ref_constraint:
//	    SCOPE FOR ( { REF ( column ) | column } ) IS [ schema. ] scope_table_name
//
// An alias takes no constraint or VISIBLE, as a view's does (ORA-00907).
func (p *Parser) parseMviewColumnAliases(stmt *nodes.CreateViewStmt) error {
	p.advance() // consume '('
	stmt.Columns = &nodes.List{}
	seenCols := map[string]Token{}
	for {
		if p.isKeywordStr("SCOPE") && p.peekNext().Type == kwFOR {
			if err := p.parseMviewScopeFor(true); err != nil {
				return err
			}
		} else {
			if !p.isCreateTableNameStart() {
				return p.syntaxErrorAtCur()
			}
			colTok := p.cur
			name, err := p.parseIdentifier()
			if err != nil {
				return err
			}
			if firstTok, dup := seenCols[name]; dup {
				return p.syntaxErrorAtTok(firstTok)
			}
			seenCols[name] = colTok
			stmt.Columns.Items = append(stmt.Columns.Items, &nodes.String{Str: name})
			if p.isKeywordStr("ENCRYPT") {
				p.advance() // consume ENCRYPT
				if err := p.parseMviewEncryptionSpec(); err != nil {
					return err
				}
			}
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return p.expectCreateTableToken(')')
}

// parseMviewScopeFor reads a scoped_table_ref_constraint. ref tells whether
// the column may be written REF ( column ), as CREATE's may.
//
//	SCOPE FOR ( { REF ( column ) | column } ) IS [ schema. ] scope_table_name
func (p *Parser) parseMviewScopeFor(ref bool) error {
	p.advance() // consume SCOPE
	if err := p.expectCreateTableToken(kwFOR); err != nil {
		return err
	}
	if err := p.expectCreateTableToken('('); err != nil {
		return err
	}
	if ref && p.cur.Type == kwREF && p.peekNext().Type == '(' {
		p.advance() // consume REF
		p.advance() // consume '('
		if err := p.parseCreateTableName(); err != nil {
			return err
		}
		if err := p.expectCreateTableToken(')'); err != nil {
			return err
		}
	} else if err := p.parseCreateTableName(); err != nil {
		return err
	}
	if err := p.expectCreateTableToken(')'); err != nil {
		return err
	}
	if err := p.expectCreateTableToken(kwIS); err != nil {
		return err
	}
	scopeTable, err := p.parseRequiredObjectName()
	_ = scopeTable // the AST keeps no scope table here
	return err
}

// parseMviewEncryptionSpec reads an encryption_spec after ENCRYPT, all of
// whose parts are optional.
//
//	[ USING 'encrypt_algorithm' ] [ IDENTIFIED BY password ]
//	[ 'integrity_algorithm' ] [ [ NO ] SALT ]
func (p *Parser) parseMviewEncryptionSpec() error {
	if p.cur.Type == kwUSING {
		p.advance() // consume USING
		if p.cur.Type != tokSCONST || !isValidOracleEncryptionAlgorithm(p.cur.Str) {
			return p.syntaxErrorAtCur()
		}
		p.advance()
	}
	if p.cur.Type == kwIDENTIFIED {
		p.advance() // consume IDENTIFIED
		if err := p.expectCreateTableToken(kwBY); err != nil {
			return err
		}
		if p.cur.Type == ')' {
			p.advance()
		} else {
			if !p.isIdentLike() {
				return p.syntaxErrorAtCur()
			}
			p.advance() // consume password
		}
	}
	if p.cur.Type == tokSCONST {
		if !isValidOracleIntegrityAlgorithm(p.cur.Str) {
			return p.syntaxErrorAtCur()
		}
		p.advance()
	}
	if p.isKeywordStr("NO") {
		p.advance() // consume NO
		if err := p.expectCreateTableWord("SALT"); err != nil {
			return err
		}
	} else if p.isKeywordStr("SALT") {
		p.advance()
	}
	return nil
}

// parseMviewPrebuiltTable reads ON PREBUILT TABLE, which stands in place of
// the physical and materialized view properties (ORA-02000 at the property
// after it, ORA-00922 at ON after one).
//
//	ON PREBUILT TABLE [ { WITH | WITHOUT } REDUCED PRECISION ]
func (p *Parser) parseMviewPrebuiltTable(stmt *nodes.CreateViewStmt) error {
	p.advance() // consume ON
	p.advance() // consume PREBUILT
	if err := p.expectCreateTableToken(kwTABLE); err != nil {
		return err
	}
	stmt.OnPrebuilt = true
	if p.cur.Type == kwWITH || p.isKeywordStr("WITHOUT") {
		with := p.cur.Type == kwWITH
		p.advance() // consume WITH / WITHOUT
		if err := p.expectCreateTableWord("REDUCED"); err != nil {
			return err
		}
		if err := p.expectCreateTableWord("PRECISION"); err != nil {
			return err
		}
		if with {
			stmt.ReducedPrec = "WITH_REDUCED"
		} else {
			stmt.ReducedPrec = "WITHOUT_REDUCED"
		}
	}
	return nil
}

// parseMviewProperties reads a materialized view's physical_properties and
// materialized_view_props, the build_clause last.
//
//	physical_properties:
//	    [ deferred_segment_creation ] [ segment_attributes_clause ] [ table_compression ]
//	    [ inmemory_table_clause ] [ { heap_org_table_clause | index_org_table_clause } ]
//	    [ CLUSTER cluster ( column [, column ] ... ) ]
//	materialized_view_props:
//	    [ column_properties ] [ table_partitioning_clauses ] [ CACHE | NOCACHE ]
//	    [ parallel_clause ] [ build_clause ]
//
// The pinned server takes the clauses before BUILD in any order, as it
// does a table's, but SEGMENT CREATION only first (ORA-12043 after another)
// and ORGANIZATION only first or after SEGMENT CREATION (ORA-64303); the
// index_org_table_clause's own parts (PCTTHRESHOLD, MAPPING TABLE,
// INCLUDING ... OVERFLOW) only after ORGANIZATION INDEX (ORA-12043). Of
// column_properties it takes nested tables, varrays, LOB storage, XMLType
// and JSON storage, and refuses object_type_col_properties (ORA-03048 at
// COLUMN) and LOB_partition_storage after a LOB_storage_clause (ORA-02000).
// The table clauses a materialized view has not, as ROWDEPENDENCIES or an
// ilm_clause, end the properties, and Oracle refuses them where AS is due.
func (p *Parser) parseMviewProperties(stmt *nodes.CreateViewStmt) error {
	// CREATE TABLE's readers keep what they read on a table; the view keeps
	// what it always has.
	table := &nodes.CreateTableStmt{}
	defer func(saved bool) { p.createTableIndexOrganized = saved }(p.createTableIndexOrganized)
	inmemoryTable, segment := false, false
	for read := 0; ; read++ {
		switch {
		case read == 0 && !segment && p.isKeywordStr("SEGMENT"):
			p.advance() // consume SEGMENT
			if err := p.expectCreateTableWord("CREATION"); err != nil {
				return err
			}
			if p.cur.Type != kwIMMEDIATE && p.cur.Type != kwDEFERRED {
				return p.syntaxErrorAtCur()
			}
			p.advance()
			// ORGANIZATION may still follow.
			segment = true
			read = -1
		case read == 0 && p.isKeywordStr("ORGANIZATION"):
			p.advance() // consume ORGANIZATION
			switch {
			case p.isKeywordStr("HEAP"):
				table.Organization = "HEAP"
			case p.cur.Type == kwINDEX:
				table.Organization = "INDEX"
				p.createTableIndexOrganized = true
			default:
				return p.syntaxErrorAtCur()
			}
			p.advance()
		case p.isKeywordStr("PCTFREE"), p.isKeywordStr("PCTUSED"), p.isKeywordStr("INITRANS"),
			p.isKeywordStr("MAXTRANS"):
			p.advance()
			if err := p.parseCreateTableInteger(); err != nil {
				return err
			}
		case p.isKeywordStr("STORAGE"):
			if err := p.parseCreateTableStorageClause(); err != nil {
				return err
			}
		case p.cur.Type == kwTABLESPACE:
			p.advance() // consume TABLESPACE
			if err := p.parseCreateTableName(); err != nil {
				return err
			}
		case p.cur.Type == kwLOGGING, p.cur.Type == kwNOLOGGING:
			p.advance()
		case p.isCreateTablePartitionCompressionStart():
			if err := p.parseCreateTablePartitionCompression(); err != nil {
				return err
			}
		case p.isCreateTableInmemoryStart():
			if err := p.parseCreateTableInmemory(table, true, &inmemoryTable, true); err != nil {
				return err
			}
		case table.Organization == "INDEX" && p.isKeywordStr("PCTTHRESHOLD"):
			p.advance() // consume PCTTHRESHOLD
			if err := p.parseCreateTableInteger(); err != nil {
				return err
			}
		case table.Organization == "INDEX" && p.isKeywordStr("MAPPING"):
			p.advance() // consume MAPPING
			if err := p.expectCreateTableToken(kwTABLE); err != nil {
				return err
			}
		case table.Organization == "INDEX" && p.isKeywordStr("NOMAPPING"):
			p.advance()
		case table.Organization == "INDEX" && p.isKeywordStr("INCLUDING"):
			p.advance() // consume INCLUDING
			if err := p.parseCreateTableName(); err != nil {
				return err
			}
			if !p.isKeywordStr("OVERFLOW") {
				return p.syntaxErrorAtCur()
			}
		case table.Organization == "INDEX" && p.isKeywordStr("OVERFLOW"):
			// OVERFLOW [ segment_attributes_clause ]: its attributes are
			// read as the properties that follow.
			p.advance()
		case p.cur.Type == kwCLUSTER:
			p.advance() // consume CLUSTER
			if err := p.parseCreateTableName(); err != nil {
				return err
			}
			if err := p.parseCreateTableColumnNamesUnique(false); err != nil {
				return err
			}
		case p.cur.Type == kwCACHE, p.cur.Type == kwNOCACHE:
			if p.cur.Type == kwCACHE {
				stmt.CacheMode = "CACHE"
			} else {
				stmt.CacheMode = "NOCACHE"
			}
			p.advance()
		case p.cur.Type == kwPARALLEL:
			p.advance() // consume PARALLEL
			stmt.ParallelMode = "PARALLEL"
			if p.cur.Type == tokICONST {
				stmt.ParallelDegree = p.cur.Str
				p.advance()
			}
		case p.cur.Type == kwNOPARALLEL:
			p.advance()
			stmt.ParallelMode = "NOPARALLEL"
		case p.isKeywordStr("NESTED"):
			if err := p.parseCreateTableNestedTableColProperties(true); err != nil {
				return err
			}
		case p.isKeywordStr("VARRAY"):
			if err := p.parseCreateTableVarrayColProperties(true); err != nil {
				return err
			}
		case p.isKeywordStr("LOB"):
			if err := p.parseCreateTableLobStorageClause(); err != nil {
				return err
			}
		case p.isKeywordStr("XMLTYPE"):
			if err := p.parseMviewXmltypeColumnProperties(); err != nil {
				return err
			}
		case p.cur.Type == kwJSON:
			if err := p.parseMviewJSONStorage(); err != nil {
				return err
			}
		case p.cur.Type == kwPARTITION:
			// PARTITION is takeable here whatever follows it: the clause is
			// `PARTITION BY …`, so a token that is not BY is where the
			// grammar stops (`PARTITION RANGE(a)…` is refused at the RANGE,
			// as Oracle refuses it), and parsePartitionClause reports it.
			if p.isMviewReferencePartitioningStart() {
				// Oracle refuses reference partitioning on a materialized
				// view at its constraint name once the clause's ')' closes
				// (ORA-14655).
				p.advance() // consume PARTITION
				p.advance() // consume BY
				p.advance() // consume REFERENCE
				if err := p.expectCreateTableToken('('); err != nil {
					return err
				}
				nameTok := p.cur
				if err := p.parseCreateTableName(); err != nil {
					return err
				}
				if err := p.expectCreateTableToken(')'); err != nil {
					return err
				}
				return p.syntaxErrorAtTok(nameTok)
			}
			savedMviewPart := p.createTablePartitionsOfMview
			p.createTablePartitionsOfMview = true
			partition, err := p.parsePartitionClause(table, true)
			p.createTablePartitionsOfMview = savedMviewPart
			if err != nil {
				return err
			}
			_ = partition // the AST keeps no partitioning of a view
		case p.isKeywordStr("CLUSTERING"):
			p.advance() // consume CLUSTERING
			if err := p.parseCreateTableClustering(); err != nil {
				return err
			}
		case p.isKeywordStr("NOMAPPING"), p.cur.Type == kwIMMEDIATE:
			p.advance()
			return p.syntaxErrorAtCur()
		case p.isKeywordStr("BUILD"):
			p.advance() // consume BUILD
			switch p.cur.Type {
			case kwIMMEDIATE:
				stmt.BuildMode = "IMMEDIATE"
			case kwDEFERRED:
				stmt.BuildMode = "DEFERRED"
			default:
				return p.syntaxErrorAtCur()
			}
			p.advance()
			return nil
		default:
			return nil
		}
	}
}

// isMviewReferencePartitioningStart reports whether PARTITION BY REFERENCE
// starts here.
func (p *Parser) isMviewReferencePartitioningStart() bool {
	ahead := p.peekAhead(2)
	return p.cur.Type == kwPARTITION && len(ahead) == 2 && ahead[0].Type == kwBY && ahead[1].Type == kwREFERENCE
}

// parseMviewXmltypeColumnProperties reads XMLType_column_properties.
//
//	XMLTYPE [ COLUMN ] column [ XMLType_storage ] [ XMLSchema_spec ]
//
//	XMLType_storage:
//	    STORE AS { OBJECT RELATIONAL
//	             | [ SECUREFILE | BASICFILE ] { CLOB | BINARY XML }
//	               [ { LOB_segname [ ( LOB_parameters ) ] | ( LOB_parameters ) } ]
//	             }
//	XMLSchema_spec:
//	    [ XMLSCHEMA XMLSchema_URL ] ELEMENT { element | XMLSchema_URL # element }
//	    [ STORE ALL VARRAYS AS { LOBS | TABLES } ]
//	    [ { ALLOW | DISALLOW } NONSCHEMA ] [ { ALLOW | DISALLOW } ANYSCHEMA ]
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-TABLE.html
func (p *Parser) parseMviewXmltypeColumnProperties() error {
	p.advance() // consume XMLTYPE
	if p.cur.Type == kwCOLUMN {
		p.advance()
	}
	if err := p.syntaxErrorIfReservedIdentifier(); err != nil {
		return err
	}
	if err := p.parseCreateTableItem(); err != nil {
		return err
	}
	objectRelational := false
	binary := false
	if p.isKeywordStr("STORE") {
		p.advance() // consume STORE
		if err := p.expectCreateTableToken(kwAS); err != nil {
			return err
		}
		if p.isKeywordStr("OBJECT") {
			p.advance() // consume OBJECT
			if err := p.expectCreateTableWord("RELATIONAL"); err != nil {
				return err
			}
			objectRelational = true
		} else {
			if p.isKeywordStr("SECUREFILE") || p.isKeywordStr("BASICFILE") {
				p.advance()
			}
			switch {
			case p.isKeywordStr("CLOB"):
				p.advance()
			case p.isKeywordStr("BINARY"):
				p.advance() // consume BINARY
				if err := p.expectCreateTableWord("XML"); err != nil {
					return err
				}
				binary = true
			default:
				return p.syntaxErrorAtCur()
			}
			// The pinned server names no LOB segment after BINARY XML
			// (ORA-02000 at the name).
			if !binary && p.isCreateTableNameStart() && !p.isMviewXmlschemaSpecStart() {
				p.advance() // consume LOB_segname
			}
			if p.cur.Type == '(' {
				if err := p.parseCreateTableLobStorageParameters(); err != nil {
					return err
				}
			}
		}
	}
	if !p.isMviewXmlschemaSpecStart() {
		if objectRelational || binary {
			return p.syntaxErrorAtCur()
		}
		return nil
	}
	if p.isKeywordStr("XMLSCHEMA") {
		p.advance() // consume XMLSCHEMA
		// Oracle refuses a string literal here (ORA-02000 at XMLSCHEMA 's'),
		// matching alter_table_properties.go's parseATXMLTypeColumnProperties.
		if err := p.parseCreateTableName(); err != nil {
			return err
		}
	}
	if err := p.expectCreateTableWord("ELEMENT"); err != nil {
		return err
	}
	if p.cur.Type == tokSCONST {
		// XMLSchema_URL # element
		p.advance()
		if err := p.expectCreateTableToken('#'); err != nil {
			return err
		}
	}
	if err := p.parseCreateTableName(); err != nil {
		return err
	}
	if p.isKeywordStr("STORE") {
		storeTok := p.advance() // consume STORE
		if err := p.expectCreateTableToken(kwALL); err != nil {
			return p.syntaxErrorAtTok(storeTok)
		}
		if err := p.expectCreateTableWord("VARRAYS"); err != nil {
			return err
		}
		if err := p.expectCreateTableToken(kwAS); err != nil {
			return err
		}
		if !p.isKeywordStr("LOBS") && !p.isKeywordStr("TABLES") {
			return p.syntaxErrorAtCur()
		}
		p.advance()
	}
	return nil
}

// isMviewXmlschemaSpecStart reports whether an XMLSchema_spec starts here.
func (p *Parser) isMviewXmlschemaSpecStart() bool {
	return p.isKeywordStr("XMLSCHEMA") || p.isKeywordStr("ELEMENT")
}

// parseMviewJSONStorage reads a JSON_storage_clause.
//
//	JSON ( json_column [, json_column ] ... ) [ STORE ] AS
//	    [ LOB_segname [ ( JSON_parameters ) ] | ( JSON_parameters ) ]
//
// JSON_parameters are the LOB_storage_parameters of a SECUREFILE LOB.
// Oracle takes STORE and what follows AS as optional, the page's
// json_storage_clause with its parameters left out: JSON (b) STORE AS
// followed by the view's own AS SELECT parses on the pinned server
// (ORA-00904 at the column it resolves), as JSON (b) AS AS SELECT does.
// One AS alone is the clause's, and leaves the view without its own
// (ORA-00906 at what follows).
func (p *Parser) parseMviewJSONStorage() error {
	p.advance() // consume JSON
	if err := p.parseCreateTableColumnNames(); err != nil {
		return err
	}
	if p.isKeywordStr("STORE") {
		p.advance()
	}
	if err := p.expectCreateTableToken(kwAS); err != nil {
		return err
	}
	if p.cur.Type == '(' {
		return p.parseCreateTableLobStorageParameters()
	}
	// A name here stands before the clause's parameters, before a
	// materialized view clause that follows, or before the view's own AS.
	// A word followed by anything else opens a clause that Oracle refuses
	// there: JSON (b) AS BUILD IMMEDIATE AS SELECT ... answers ORA-02000 at
	// the AS before the query, AS REFRESH FAST ORA-02000 at 48, AS USING
	// INDEX ORA-02461 at 46, where AS seg CACHE AS SELECT, AS seg BUILD
	// IMMEDIATE AS SELECT and AS BUILD AS SELECT name the segment (ORA-00904
	// at the column they resolve).
	if p.isCreateTableNameStart() {
		if next := p.peekNext(); next.Type != '(' && next.Type != kwAS && !p.isMviewClauseStart(next) {
			return p.syntaxErrorAtCur()
		}
		if err := p.parseCreateTableName(); err != nil {
			return err
		}
		if p.cur.Type == '(' {
			return p.parseCreateTableLobStorageParameters()
		}
	}
	return nil
}

// isMviewClauseStart reports whether a clause of CREATE MATERIALIZED VIEW
// starts at tok: what parseMviewProperties reads and what follows it.
func (p *Parser) isMviewClauseStart(tok Token) bool {
	switch tok.Type {
	case kwTABLESPACE, kwLOGGING, kwNOLOGGING, kwCACHE, kwNOCACHE, kwPARALLEL, kwNOPARALLEL,
		kwCLUSTER, kwPARTITION, kwJSON, kwUSING, kwREFRESH, kwFOR, kwENABLE, kwDISABLE, kwAS:
		return true
	}
	for _, w := range []string{
		"NESTED", "VARRAY", "LOB", "XMLTYPE", "BUILD", "INMEMORY", "PCTFREE", "PCTUSED",
		"INITRANS", "MAXTRANS", "STORAGE", "COMPRESS", "NOCOMPRESS", "PCTTHRESHOLD", "MAPPING",
		"NOMAPPING", "INCLUDING", "OVERFLOW", "NEVER", "EVALUATE", "ANNOTATIONS",
	} {
		if isCreateTableWordTok(tok, w) {
			return true
		}
	}
	return false
}

// parseMviewUsingIndex reads the index attributes after USING INDEX.
//
//	USING INDEX [ physical_attributes_clause | TABLESPACE tablespace ] ...
//
// The pinned server refuses the other index properties here (ORA-02000 at
// LOCAL, NOLOGGING or COMPRESS).
func (p *Parser) parseMviewUsingIndex() error {
	for {
		switch {
		case p.isKeywordStr("PCTFREE"), p.isKeywordStr("PCTUSED"), p.isKeywordStr("INITRANS"),
			p.isKeywordStr("MAXTRANS"):
			p.advance()
			if err := p.parseCreateTableInteger(); err != nil {
				return err
			}
		case p.isKeywordStr("STORAGE"):
			if err := p.parseCreateTableStorageClause(); err != nil {
				return err
			}
		case p.cur.Type == kwTABLESPACE:
			p.advance() // consume TABLESPACE
			if err := p.parseCreateTableName(); err != nil {
				return err
			}
		default:
			return nil
		}
	}
}

// parseMviewCreateRefresh reads create_mv_refresh's REFRESH form.
//
//	REFRESH
//	    { { FAST | COMPLETE | FORCE }
//	    | ON { DEMAND | COMMIT | STATEMENT }
//	    | { START WITH | NEXT } date
//	    | WITH { PRIMARY KEY | ROWID }
//	    | USING { DEFAULT [ MASTER | LOCAL ] ROLLBACK SEGMENT
//	            | [ MASTER | LOCAL ] ROLLBACK SEGMENT rollback_segment }
//	      [ DEFAULT [ MASTER | LOCAL ] ROLLBACK SEGMENT
//	      | [ MASTER | LOCAL ] ROLLBACK SEGMENT rollback_segment ] ...
//	    | USING { ENFORCED | TRUSTED } CONSTRAINTS
//	    } ...
//
// That is the 19c page's form: the pinned server takes the parts in any
// order and at least one (ORA-00905 at AS after REFRESH), the refresh method
// more than once, and the rollback segment clauses any number of times, but
// each other part once (ORA-02000 at a second START WITH, NEXT or WITH,
// ORA-12990 at a second ON, ORA-00905 at a second USING ... CONSTRAINTS).
func (p *Parser) parseMviewCreateRefresh(stmt *nodes.CreateViewStmt) error {
	p.advance() // consume REFRESH
	var on, start, next, with, constraints bool
	for read := 0; ; read++ {
		switch {
		case p.isKeywordStr("FAST"), p.isKeywordStr("COMPLETE"), p.cur.Type == kwFORCE:
			stmt.RefreshMethod = p.cur.Str
			p.advance()
		case !on && p.cur.Type == kwON:
			p.advance() // consume ON
			switch {
			case p.cur.Type == kwCOMMIT:
				stmt.RefreshMode = "ON COMMIT"
			case p.isKeywordStr("DEMAND"):
				stmt.RefreshMode = "ON DEMAND"
			case p.isKeywordStr("STATEMENT"):
				stmt.RefreshMode = "ON STATEMENT"
				stmt.RefreshOnStmt = true
			default:
				return p.syntaxErrorAtCur()
			}
			p.advance()
			on = true
		case !start && p.cur.Type == kwSTART:
			p.advance() // consume START
			if err := p.expectCreateTableToken(kwWITH); err != nil {
				return err
			}
			expr, err := p.parseMviewDateExpr()
			if err != nil {
				return err
			}
			stmt.StartWith = expr
			start = true
		case !next && p.cur.Type == kwNEXT:
			p.advance() // consume NEXT
			expr, err := p.parseMviewDateExpr()
			if err != nil {
				return err
			}
			stmt.Next = expr
			next = true
		case !with && p.cur.Type == kwWITH:
			p.advance() // consume WITH
			switch {
			case p.cur.Type == kwPRIMARY:
				p.advance() // consume PRIMARY
				if err := p.expectCreateTableToken(kwKEY); err != nil {
					return err
				}
				stmt.WithPK = true
			case p.cur.Type == kwROWID:
				p.advance()
				stmt.WithRowID = true
			default:
				return p.syntaxErrorAtCur()
			}
			with = true
		case p.cur.Type == kwUSING:
			p.advance() // consume USING
			if p.isKeywordStr("TRUSTED") || p.isKeywordStr("ENFORCED") {
				if constraints {
					return p.syntaxErrorAtCur()
				}
				p.advance() // consume TRUSTED / ENFORCED
				if err := p.expectCreateTableToken(kwCONSTRAINTS); err != nil {
					return err
				}
				constraints = true
				break
			}
			if err := p.parseMviewRollbackSegment(); err != nil {
				return err
			}
			for p.cur.Type == kwDEFAULT || p.cur.Type == kwROLLBACK || p.isKeywordStr("MASTER") || p.isKeywordStr("LOCAL") {
				if err := p.parseMviewRollbackSegment(); err != nil {
					return err
				}
			}
		default:
			if read == 0 {
				return p.syntaxErrorAtCur()
			}
			return nil
		}
	}
}

// parseMviewRollbackSegment reads one rollback segment clause of
// create_mv_refresh: a DEFAULT one names no segment, any other names one
// (ORA-02000 at a name after DEFAULT ... SEGMENT, ORA-00905 at AS after
// ROLLBACK SEGMENT).
//
//	{ DEFAULT [ MASTER | LOCAL ] ROLLBACK SEGMENT
//	| [ MASTER | LOCAL ] ROLLBACK SEGMENT rollback_segment }
func (p *Parser) parseMviewRollbackSegment() error {
	def := p.cur.Type == kwDEFAULT
	if def {
		p.advance() // consume DEFAULT
	}
	if p.isKeywordStr("MASTER") || p.isKeywordStr("LOCAL") {
		p.advance()
	}
	if err := p.expectCreateTableToken(kwROLLBACK); err != nil {
		return err
	}
	if err := p.expectCreateTableWord("SEGMENT"); err != nil {
		return err
	}
	if def {
		return nil
	}
	return p.parseCreateTableName()
}

// parseMviewDateExpr reads the date expression of START WITH or NEXT.
func (p *Parser) parseMviewDateExpr() (nodes.ExprNode, error) {
	if p.cur.Type == kwAS {
		return nil, p.syntaxErrorAtCur()
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.syntaxErrorAtCur()
	}
	return expr, nil
}

// parseMviewEvaluationEdition reads CREATE's evaluation_edition_clause.
//
//	EVALUATE USING { CURRENT EDITION | EDITION edition | NULL EDITION }
func (p *Parser) parseMviewEvaluationEdition() error {
	p.advance() // consume EVALUATE
	if err := p.expectCreateTableToken(kwUSING); err != nil {
		return err
	}
	return p.parseMviewEditionRef(true)
}

// parseMviewEditionRef reads { CURRENT EDITION | EDITION edition }, and
// NULL EDITION when null is set.
func (p *Parser) parseMviewEditionRef(null bool) error {
	switch {
	case p.cur.Type == kwCURRENT, null && p.cur.Type == kwNULL:
		p.advance() // consume CURRENT / NULL
		return p.expectCreateTableWord("EDITION")
	case p.isKeywordStr("EDITION"):
		p.advance() // consume EDITION
		return p.parseCreateTableName()
	}
	return p.syntaxErrorAtCur()
}

// parseMviewUnusableEditions reads CREATE's unusable_editions_clause, which
// the pinned server takes only after a query rewrite or query computation
// clause (ORA-00922 at UNUSABLE otherwise), in the 19c page's form: the
// 26 page's UNUSABLE EDITIONS is refused at EDITIONS (ORA-02000).
//
//	[ UNUSABLE BEFORE { CURRENT EDITION | EDITION edition } ]
//	[ UNUSABLE BEGINNING WITH { CURRENT EDITION | EDITION edition | NULL EDITION } ]
func (p *Parser) parseMviewUnusableEditions() error {
	if p.isKeywordStr("UNUSABLE") && isCreateTableWordTok(p.peekNext(), "BEFORE") {
		p.advance() // consume UNUSABLE
		p.advance() // consume BEFORE
		if err := p.parseMviewEditionRef(false); err != nil {
			return err
		}
	}
	if p.isKeywordStr("UNUSABLE") {
		p.advance() // consume UNUSABLE
		if err := p.expectCreateTableWord("BEGINNING"); err != nil {
			return err
		}
		if err := p.expectCreateTableToken(kwWITH); err != nil {
			return err
		}
		if err := p.parseMviewEditionRef(true); err != nil {
			return err
		}
	}
	return nil
}

// parseMviewAnnotations reads an annotations_clause. ALTER takes a single
// item; CREATE a list of them.
//
//	ANNOTATIONS ( annotation_item [, annotation_item ] ... )
//
//	annotation_item:
//	    [ ADD [ IF NOT EXISTS | OR REPLACE ] | DROP [ IF EXISTS ] | REPLACE ]
//	    annotation_name [ 'annotation_value' ]
//
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/annotations_clause.html
func (p *Parser) parseMviewAnnotations(single bool) error {
	p.advance() // consume ANNOTATIONS
	if err := p.expectCreateTableToken('('); err != nil {
		return err
	}
	for {
		switch {
		case p.cur.Type == kwADD:
			p.advance() // consume ADD
			switch {
			case p.cur.Type == kwIF:
				p.advance() // consume IF
				if err := p.expectCreateTableToken(kwNOT); err != nil {
					return err
				}
				if err := p.expectCreateTableToken(kwEXISTS); err != nil {
					return err
				}
			case p.cur.Type == kwOR:
				p.advance() // consume OR
				if err := p.expectCreateTableToken(kwREPLACE); err != nil {
					return err
				}
			}
		case p.cur.Type == kwDROP:
			p.advance() // consume DROP
			if p.cur.Type == kwIF {
				p.advance() // consume IF
				if err := p.expectCreateTableToken(kwEXISTS); err != nil {
					return err
				}
			}
		case p.cur.Type == kwREPLACE:
			p.advance()
		}
		if p.cur.Type != tokQIDENT && !p.isIdentLike() {
			return p.syntaxErrorAtCur()
		}
		p.advance() // consume annotation_name
		if p.cur.Type == tokSCONST {
			p.advance()
		}
		if single || p.cur.Type != ',' {
			break
		}
		p.advance() // consume ','
	}
	return p.expectCreateTableToken(')')
}
