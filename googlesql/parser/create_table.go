package parser

import (
	"github.com/bytebase/omni/googlesql/ast"
)

// ---------------------------------------------------------------------------
// Core DDL — CREATE TABLE + shared DDL helpers (parser-ddl node)
// ---------------------------------------------------------------------------
//
// This file ports the legacy ANTLR create_table_statement and the shared
// table-element / column / constraint / OPTIONS grammar
// (GoogleSQLParser.g4 §2.4-§2.6 — table_element_list, table_column_definition,
// table_constraint_definition, primary_key_spec, foreign_key_reference,
// generated_column_info, options_list, …), a hand-port of Google's open-source
// ZetaSQL reference and the grammar bytebase consumes for both BigQuery and
// Spanner. One omni parser serves both dialects: the production accepts the
// UNION of BigQuery and Spanner forms.
//
// ORACLE NOTE — Spanner is a SUBSET of the GoogleSQL union (oracle.md). The
// Spanner emulator (the differential harness) rejects BigQuery-only CREATE TABLE
// clauses (OR REPLACE, TEMP/TEMPORARY, PARTITION BY, CLUSTER BY, OPTIONS, CTAS,
// inline PRIMARY KEY, NOT ENFORCED, LIKE, COPY, GENERATED ALWAYS AS,
// STORED VOLATILE, VIRTUAL) as syntax errors — those rejects are NON-AUTHORITATIVE
// and are triangulated against the legacy GoogleSQLParser.g4 (which has every one
// of these productions) plus the BigQuery truth1 corpus. The Spanner-authoritative
// forms (INTERLEAVE IN PARENT, generated AS (expr) STORED, ROW DELETION POLICY,
// the trailing PRIMARY KEY ( … ) element list, the PK-required rule) are verified
// against the live emulator. See create_table_oracle_test.go.

// parseCreateStmt is the CREATE second-keyword dispatcher. The CREATE keyword has
// NOT yet been consumed when this is called (parseStmt peeks it). It consumes the
// optional CREATE-statement prefixes (OR REPLACE, create-scope) that are shared
// across CREATE forms, then dispatches on the object keyword.
//
// This node owns TABLE / VIEW / INDEX / SCHEMA / DATABASE. Other object kinds
// (FUNCTION / PROCEDURE / MODEL / EXTERNAL / CONNECTION / CONSTANT / SNAPSHOT /
// PROPERTY GRAPH / ROW ACCESS POLICY / PRIVILEGE RESTRICTION / generic entity, and
// the MATERIALIZED|APPROX VIEW and SEARCH|VECTOR INDEX variants) belong to the
// parser-ddl-{bigquery,spanner} nodes; until those land they route to the
// unsupported stub, which emits a "CREATE …" diagnostic so Diagnose reports them
// as "not yet supported" rather than "unknown statement".
func (p *Parser) parseCreateStmt() (ast.Node, error) {
	create := p.advance() // consume CREATE

	// opt_or_replace? opt_create_scope?  — these precede the object keyword in
	// the grammar for most CREATE forms. CREATE DATABASE and the bare CREATE
	// SCHEMA path do not take them, but parsing them up front is harmless: the
	// object dispatch below rejects an OR REPLACE / scope before DATABASE (no
	// such production) by routing to the right per-object parser, which would
	// then reject the unexpected leading keyword. We capture them and pass them
	// to the per-object parser so it can honor or reject them per its grammar.
	orReplace := false
	if p.cur.Type == kwOR {
		// OR REPLACE: consume both (REPLACE is required after OR here).
		p.advance() // OR
		if _, err := p.expect(kwREPLACE); err != nil {
			return nil, err
		}
		orReplace = true
	}
	scope := p.matchCreateScope()

	switch p.cur.Type {
	case kwTABLE:
		// CREATE TABLE — but CREATE TABLE FUNCTION (TVF) is a BigQuery-only object
		// owned by parser-ddl-bigquery; route it to the stub.
		if p.peekNext().Type == kwFUNCTION {
			return p.unsupported("CREATE TABLE FUNCTION")
		}
		return p.parseCreateTable(create, orReplace, scope)
	case kwVIEW:
		return p.parseCreateView(create, orReplace, scope, false /*recursive*/)
	case kwRECURSIVE:
		// RECURSIVE VIEW (the plain-view alt with RECURSIVE). MATERIALIZED/APPROX
		// RECURSIVE views are dialect-node territory; a bare RECURSIVE followed by
		// VIEW is the plain recursive view this node owns.
		if p.peekNext().Type == kwVIEW {
			p.advance() // RECURSIVE
			return p.parseCreateView(create, orReplace, scope, true /*recursive*/)
		}
		return p.unsupported("CREATE")
	case kwUNIQUE, kwNULL_FILTERED, kwINDEX:
		// CREATE [UNIQUE] [NULL_FILTERED] INDEX. The SEARCH/VECTOR index_type
		// variants are dialect-node territory; parseCreateIndex rejects a
		// SEARCH/VECTOR token after the qualifiers by routing to the stub.
		return p.parseCreateIndex(create, orReplace, scope)
	case kwSEARCH, kwVECTOR:
		// CREATE SEARCH|VECTOR INDEX — dialect-specific (parser-ddl-{bigquery,
		// spanner}).
		return p.unsupported("CREATE " + p.cur.Str)
	case kwSCHEMA:
		return p.parseCreateSchema(create, orReplace)
	case kwDATABASE:
		return p.parseCreateDatabase(create, orReplace, scope)
	default:
		// Any other object kind (FUNCTION / PROCEDURE / MODEL / EXTERNAL /
		// MATERIALIZED|APPROX VIEW / SNAPSHOT / generic entity / a bare identifier
		// such as the dispatch-coverage `CREATE x` probe) is not owned here; emit
		// a "CREATE …" not-yet-supported diagnostic.
		return p.unsupported("CREATE")
	}
}

// matchCreateScope consumes an opt_create_scope (TEMP | TEMPORARY | PUBLIC |
// PRIVATE) if present and returns its normalized keyword spelling, or "" when
// absent.
func (p *Parser) matchCreateScope() string {
	switch p.cur.Type {
	case kwTEMP:
		p.advance()
		return "TEMP"
	case kwTEMPORARY:
		p.advance()
		return "TEMPORARY"
	case kwPUBLIC:
		p.advance()
		return "PUBLIC"
	case kwPRIVATE:
		p.advance()
		return "PRIVATE"
	}
	return ""
}

// parseCreateTable parses the body of a CREATE TABLE statement after the shared
// CREATE prefix (OR REPLACE / scope) has been consumed and cur is at the TABLE
// keyword. It follows create_table_statement:
//
//	TABLE opt_if_not_exists? maybe_dashed_path_expression table_element_list?
//	  opt_spanner_table_options? opt_like_path_expression? opt_clone_table?
//	  opt_copy_table? opt_default_collate_clause?
//	  partition_by_clause_prefix_no_hint? cluster_by_clause_prefix_no_hint?
//	  opt_ttl_clause? with_connection_clause? opt_options_list? as_query?
func (p *Parser) parseCreateTable(create Token, orReplace bool, scope string) (ast.Node, error) {
	p.advance() // consume TABLE

	stmt := &ast.CreateTableStmt{OrReplace: orReplace, Scope: scope}
	stmt.Loc.Start = create.Loc.Start

	ifNotExists, err := p.parseIfNotExists()
	if err != nil {
		return nil, err
	}
	stmt.IfNotExists = ifNotExists

	name, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// table_element_list? — `( table_element , … )`. Optional (CREATE TABLE LIKE /
	// COPY / CLONE / CTAS omit it).
	if p.cur.Type == int('(') {
		cols, cons, err := p.parseTableElementList()
		if err != nil {
			return nil, err
		}
		stmt.Columns = cols
		stmt.Constraints = cons
	}

	// opt_spanner_table_options? — trailing `PRIMARY KEY ( … )` plus optional
	// comma-led `, INTERLEAVE IN PARENT …` and/or `, ROW DELETION POLICY ( … )`.
	hasSpannerPK := false
	if p.cur.Type == kwPRIMARY {
		pk, interleave, rowDeletion, err := p.parseSpannerTableOptions()
		if err != nil {
			return nil, err
		}
		stmt.PrimaryKey = pk
		stmt.HasPrimaryKey = true
		stmt.Interleave = interleave
		stmt.RowDeletion = rowDeletion
		hasSpannerPK = true
	}

	// opt_like_path_expression?
	if p.cur.Type == kwLIKE {
		p.advance()
		path, err := p.parseTablePath()
		if err != nil {
			return nil, err
		}
		stmt.Like = path
	}

	// opt_clone_table? — CLONE <data source> (path [AT SYSTEM TIME] [WHERE]); this
	// node captures only the source path (the time-travel / WHERE filter of CLONE
	// DATA-style sources is dialect-node territory).
	if p.cur.Type == kwCLONE {
		p.advance()
		path, err := p.parseTablePath()
		if err != nil {
			return nil, err
		}
		stmt.Clone = path
	}

	// opt_copy_table?
	if p.cur.Type == kwCOPY {
		p.advance()
		path, err := p.parseTablePath()
		if err != nil {
			return nil, err
		}
		stmt.Copy = path
	}

	// opt_default_collate_clause?  — DEFAULT collate_clause.
	if p.cur.Type == kwDEFAULT {
		p.advance() // DEFAULT
		coll, _, err := p.parseCollateClause()
		if err != nil {
			return nil, err
		}
		stmt.DefaultCollate = coll
	}

	// partition_by_clause_prefix_no_hint?  — PARTITION BY expr (, expr)*.
	if p.cur.Type == kwPARTITION {
		exprs, err := p.parsePartitionByNoHint()
		if err != nil {
			return nil, err
		}
		stmt.PartitionBy = exprs
	}

	// cluster_by_clause_prefix_no_hint?  — CLUSTER BY expr (, expr)*.
	if p.cur.Type == kwCLUSTER {
		exprs, err := p.parseClusterByNoHint()
		if err != nil {
			return nil, err
		}
		stmt.ClusterBy = exprs
	}

	// opt_ttl_clause?  — a standalone (non-comma-led) ROW DELETION POLICY ( expr ).
	// After a Spanner trailing PRIMARY KEY, ROW DELETION POLICY is comma-led and
	// is consumed in parseSpannerTableOptions; a NON-comma ROW after the PK is a
	// syntax error (oracle: `PRIMARY KEY (x) ROW DELETION POLICY (…)` without the
	// comma rejects). So this standalone arm only applies when there was no
	// trailing Spanner PK (it leaves a bare ROW after a Spanner PK unconsumed,
	// which the outer EOF check then reports as the syntax error it is).
	if p.cur.Type == kwROW && !hasSpannerPK {
		expr, err := p.parseRowDeletionPolicyBody()
		if err != nil {
			return nil, err
		}
		stmt.RowDeletion = expr
	}

	// with_connection_clause?  — WITH connection_clause (WITH CONNECTION …). The
	// connection clause itself is dialect-node territory; recognized here so a
	// BigQuery CREATE TABLE … WITH CONNECTION does not mis-parse, but the body is
	// skipped to the connection path. Kept minimal: WITH CONNECTION <path>.
	if p.cur.Type == kwWITH && p.peekNext().Type == kwCONNECTION {
		if err := p.skipWithConnection(); err != nil {
			return nil, err
		}
	}

	// opt_options_list?
	if p.cur.Type == kwOPTIONS {
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		stmt.Options = opts
	}

	// as_query?  — AS query (CTAS).
	if p.cur.Type == kwAS {
		p.advance() // AS
		q, err := p.parseQuery()
		if err != nil {
			return nil, err
		}
		stmt.AsQuery = q
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseTableElementList parses `( table_element (, table_element)* ,? )` — the
// opening '(' is the current token. Each element is a column definition or a
// table-level constraint. A trailing comma before ')' is allowed (the grammar's
// COMMA_SYMBOL?). Returns the columns and constraints in source order (the AST
// keeps them in two slices; the original interleaving is not load-bearing for
// the query-span consumer).
func (p *Parser) parseTableElementList() ([]*ast.ColumnDef, []*ast.TableConstraint, error) {
	p.advance() // consume '('

	var cols []*ast.ColumnDef
	var cons []*ast.TableConstraint

	// Empty element list `()` is legal (oracle: `CREATE TABLE T () PRIMARY KEY ()`
	// accepts).
	if p.cur.Type == int(')') {
		p.advance()
		return cols, cons, nil
	}

	for {
		col, con, err := p.parseTableElement()
		if err != nil {
			return nil, nil, err
		}
		if col != nil {
			cols = append(cols, col)
		} else {
			cons = append(cons, con)
		}

		if _, ok := p.match(int(',')); ok {
			// Trailing comma: `, )` ends the list.
			if p.cur.Type == int(')') {
				break
			}
			continue
		}
		break
	}

	if _, err := p.expect(int(')')); err != nil {
		return nil, nil, err
	}
	return cols, cons, nil
}

// parseTableElement parses one table_element: either a table_constraint_definition
// (when it begins with PRIMARY / FOREIGN / CHECK / CONSTRAINT) or a
// table_column_definition (otherwise — a column name followed by its schema).
// Exactly one of the returned (column, constraint) is non-nil.
func (p *Parser) parseTableElement() (*ast.ColumnDef, *ast.TableConstraint, error) {
	switch p.cur.Type {
	case kwPRIMARY, kwFOREIGN, kwCHECK, kwCONSTRAINT:
		con, err := p.parseTableConstraint()
		if err != nil {
			return nil, nil, err
		}
		return nil, con, nil
	default:
		col, err := p.parseColumnDefinition()
		if err != nil {
			return nil, nil, err
		}
		return col, nil, nil
	}
}

// parseColumnDefinition parses a table_column_definition:
//
//	identifier table_column_schema column_attributes? opt_options_list?
//
// where table_column_schema is `column_schema_inner collate_clause?
// opt_column_info?` or `generated_column_info`, column_attributes is
// `column_attribute+ constraint_enforcement?` (PRIMARY KEY / FOREIGN KEY /
// HIDDEN / NOT NULL), and opt_column_info is a DEFAULT or generated clause.
func (p *Parser) parseColumnDefinition() (*ast.ColumnDef, error) {
	if !isIdentifierStart(p.cur.Type) {
		return nil, p.syntaxErrorAtCur()
	}
	nameTok := p.advance()
	name, err := p.identifierText(nameTok)
	if err != nil {
		return nil, err
	}
	col := &ast.ColumnDef{Name: name, Loc: nameTok.Loc}

	// table_column_schema: either a leading generated_column_info (a column with
	// no explicit type — the `AS (expr)` / `GENERATED … AS` form) or
	// `column_schema_inner collate_clause? opt_column_info?`.
	if p.atGeneratedMode() {
		if err := p.parseGeneratedColumnInfo(col); err != nil {
			return nil, err
		}
	} else {
		// column_schema_inner — the column type (parseType covers
		// simple/array/struct/range + type parameters).
		dt, err := p.parseType()
		if err != nil {
			return nil, err
		}
		col.Type = &ast.TypeRef{Text: dt.String(), Loc: dt.Loc}

		// collate_clause? after the type.
		if p.cur.Type == kwCOLLATE {
			coll, _, err := p.parseCollateClause()
			if err != nil {
				return nil, err
			}
			col.Collate = coll
		}

		// opt_column_info?: a DEFAULT or a generated clause (mutually exclusive —
		// the legacy grammar rejects both with a dedicated error).
		switch {
		case p.cur.Type == kwDEFAULT:
			p.advance() // DEFAULT
			expr, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			col.Default = expr
			// The grammar's invalid_generated_column guard: a generated clause must
			// not follow a DEFAULT. If one does, report the canonical error.
			if p.atGeneratedMode() {
				return nil, &ParseError{
					Loc: p.cur.Loc,
					Msg: `syntax error: "DEFAULT" and "GENERATED ALWAYS AS" clauses must not be both provided for the column`,
				}
			}
		case p.atGeneratedMode():
			if err := p.parseGeneratedColumnInfo(col); err != nil {
				return nil, err
			}
			// invalid_default_column guard: a DEFAULT must not follow a generated
			// clause.
			if p.cur.Type == kwDEFAULT {
				return nil, &ParseError{
					Loc: p.cur.Loc,
					Msg: `syntax error: "DEFAULT" and "GENERATED ALWAYS AS" clauses must not be both provided for the column`,
				}
			}
		}
	}

	// column_attributes?: PRIMARY KEY / FOREIGN KEY / HIDDEN / NOT NULL, in any
	// order, with an optional trailing constraint_enforcement.
	if err := p.parseColumnAttributes(col); err != nil {
		return nil, err
	}

	// opt_options_list?
	if p.cur.Type == kwOPTIONS {
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		col.Options = opts
	}

	col.Loc.End = p.prev.Loc.End
	return col, nil
}

// atGeneratedMode reports whether cur begins a generated_mode prefix: GENERATED
// (any of the GENERATED forms) or a bare AS (the Spanner `AS (expr) STORED`
// generated-column spelling). It does NOT match the AS that introduces a CTAS or
// view body — generated mode appears only in column-schema position, which is
// where this is called.
func (p *Parser) atGeneratedMode() bool {
	return p.cur.Type == kwGENERATED || p.cur.Type == kwAS
}

// parseGeneratedColumnInfo parses a generated_column_info into col:
//
//	generated_mode '(' expression ')' stored_mode?
//	generated_mode identity_column_info      (BigQuery IDENTITY — recognized
//	                                          structurally but rare; captured as
//	                                          a generated column with no expr)
//
// generated_mode is `AS | GENERATED AS | GENERATED ALWAYS AS | GENERATED BY
// DEFAULT AS`. stored_mode is `STORED [VOLATILE]`; the Spanner-only `VIRTUAL`
// spelling (which the legacy grammar does NOT carry and the emulator rejects —
// divergence) is also accepted here, captured as Stored="VIRTUAL", because the
// BigQuery+Spanner truth1 corpus documents it (oracle non-authoritative for it).
func (p *Parser) parseGeneratedColumnInfo(col *ast.ColumnDef) error {
	mode, err := p.parseGeneratedMode()
	if err != nil {
		return err
	}
	col.GenMode = mode

	// generated_mode identity_column_info — IDENTITY(...) form. IDENTITY is a
	// keyword (kwIDENTITY in the keyword table), so the lexer always emits the
	// keyword token, never tokIdentifier. We parse the body STRUCTURALLY (per the
	// .g4 identity_column_info) rather than skipping balanced parens, so a malformed
	// body such as `IDENTITY (FOO BAR)` or `IDENTITY (START WITH)` is rejected
	// instead of silently accepted. This enables the (BigQuery-valid)
	// `GENERATED ALWAYS AS IDENTITY(…)` form while still rejecting garbage. An
	// identity column carries no generated expression and no stored_mode (the
	// query-span path does not inspect it). The Spanner emulator rejects this
	// BigQuery-only form (non-authoritative; triangulated from the .g4 +
	// BigQuery docs).
	if p.cur.Type == kwIDENTITY {
		if err := p.parseIdentityColumnInfo(); err != nil {
			return err
		}
		col.Stored = "" // identity columns have no stored_mode
		return nil
	}

	// generated_mode '(' expression ')' stored_mode?
	if _, err := p.expect(int('(')); err != nil {
		return err
	}
	expr, err := p.parseExpr()
	if err != nil {
		return err
	}
	if _, err := p.expect(int(')')); err != nil {
		return err
	}
	col.Generated = expr

	// stored_mode?  — STORED [VOLATILE] | VIRTUAL.
	switch {
	case p.cur.Type == kwSTORED:
		p.advance() // STORED
		if p.cur.Type == kwVOLATILE {
			p.advance()
			col.Stored = "STORED VOLATILE"
		} else {
			col.Stored = "STORED"
		}
	case p.cur.Type == tokIdentifier && (p.cur.Str == "VIRTUAL" || p.cur.Str == "virtual"):
		// Spanner VIRTUAL generated column. Not in the legacy .g4 stored_mode and
		// rejected by the Spanner emulator (divergence: truth1 documents it, oracle
		// non-authoritative). Accepted so the documented Spanner form parses.
		p.advance()
		col.Stored = "VIRTUAL"
	}
	return nil
}

// parseGeneratedMode parses a generated_mode prefix and returns the enum.
//
//	AS                         -> GenModeAs
//	GENERATED AS               -> GenModeGeneratedAs
//	GENERATED ALWAYS AS        -> GenModeGeneratedAlways
//	GENERATED BY DEFAULT AS    -> GenModeGeneratedByDefault
func (p *Parser) parseGeneratedMode() (ast.GeneratedMode, error) {
	if p.cur.Type == kwAS {
		p.advance()
		return ast.GenModeAs, nil
	}
	// GENERATED ...
	if _, err := p.expect(kwGENERATED); err != nil {
		return 0, err
	}
	switch p.cur.Type {
	case kwAS:
		p.advance()
		return ast.GenModeGeneratedAs, nil
	case kwALWAYS:
		p.advance() // ALWAYS
		if _, err := p.expect(kwAS); err != nil {
			return 0, err
		}
		return ast.GenModeGeneratedAlways, nil
	case kwBY:
		p.advance() // BY
		if _, err := p.expect(kwDEFAULT); err != nil {
			return 0, err
		}
		if _, err := p.expect(kwAS); err != nil {
			return 0, err
		}
		return ast.GenModeGeneratedByDefault, nil
	default:
		return 0, p.syntaxErrorAtCur()
	}
}

// parseIdentityColumnInfo parses an identity_column_info (the BigQuery IDENTITY
// generated-column form). cur is at the IDENTITY keyword:
//
//	IDENTITY '(' opt_start_with? opt_increment_by? opt_maxvalue? opt_minvalue?
//	  opt_cycle? ')'
//	  opt_start_with   : START WITH signed_numeric_literal
//	  opt_increment_by : INCREMENT BY signed_numeric_literal
//	  opt_maxvalue     : MAXVALUE signed_numeric_literal
//	  opt_minvalue     : MINVALUE signed_numeric_literal
//	  opt_cycle        : CYCLE | NO CYCLE
//
// The clauses are each optional and appear AT MOST ONCE in this fixed order; we
// consume them sequentially, which enforces both order and at-most-once (an
// out-of-order or repeated clause falls through to the ')' expectation and is
// rejected). The body is validated rather than skipped so malformed identity
// bodies (e.g. `IDENTITY (FOO BAR)`, `IDENTITY (START WITH)`) reject. The parsed
// values are not retained (the query-span path does not inspect them); only the
// accept/reject decision is load-bearing.
func (p *Parser) parseIdentityColumnInfo() error {
	p.advance() // IDENTITY
	if _, err := p.expect(int('(')); err != nil {
		return err
	}

	if p.cur.Type == kwSTART {
		p.advance() // START
		if _, err := p.expect(kwWITH); err != nil {
			return err
		}
		if err := p.parseSignedNumericLiteral(); err != nil {
			return err
		}
	}
	if p.cur.Type == kwINCREMENT {
		p.advance() // INCREMENT
		if _, err := p.expect(kwBY); err != nil {
			return err
		}
		if err := p.parseSignedNumericLiteral(); err != nil {
			return err
		}
	}
	if p.cur.Type == kwMAXVALUE {
		p.advance() // MAXVALUE
		if err := p.parseSignedNumericLiteral(); err != nil {
			return err
		}
	}
	if p.cur.Type == kwMINVALUE {
		p.advance() // MINVALUE
		if err := p.parseSignedNumericLiteral(); err != nil {
			return err
		}
	}
	// opt_cycle: CYCLE | NO CYCLE.
	switch p.cur.Type {
	case kwCYCLE:
		p.advance()
	case kwNO:
		p.advance() // NO
		if _, err := p.expect(kwCYCLE); err != nil {
			return err
		}
	}

	if _, err := p.expect(int(')')); err != nil {
		return err
	}
	return nil
}

// parseSignedNumericLiteral consumes a signed_numeric_literal:
//
//	integer_literal | numeric_literal | bignumeric_literal | floating_point_literal
//	| MINUS integer_literal | MINUS floating_point_literal
//
// where numeric_literal / bignumeric_literal are the typed-string forms
// `(NUMERIC|DECIMAL) 'str'` / `(BIGNUMERIC|BIGDECIMAL) 'str'`. A leading sign is
// allowed only before a plain integer or float (per the grammar). It does not
// retain the value — it is used only to validate an identity_column_info clause.
func (p *Parser) parseSignedNumericLiteral() error {
	if p.cur.Type == int('-') {
		p.advance() // MINUS — only valid before an integer or float literal
		switch p.cur.Type {
		case tokInteger, tokFloat:
			p.advance()
			return nil
		default:
			return p.syntaxErrorAtCur()
		}
	}
	switch p.cur.Type {
	case tokInteger, tokFloat:
		p.advance()
		return nil
	case kwNUMERIC, kwDECIMAL, kwBIGNUMERIC, kwBIGDECIMAL:
		// numeric_literal / bignumeric_literal: a type-keyword prefix + string_literal.
		// string_literal is one-or-more ADJACENT STRING_LITERAL components (GoogleSQL
		// concatenates adjacent strings, e.g. NUMERIC '1' '2'), so consume the full
		// component run via the shared helper rather than a single token.
		p.advance() // the type-keyword prefix
		if p.cur.Type != tokString {
			return p.syntaxErrorAtCur()
		}
		if _, err := p.parseStringLiteral(); err != nil {
			return err
		}
		return nil
	default:
		return p.syntaxErrorAtCur()
	}
}

// parseColumnAttributes parses an optional column_attributes run into col:
// `column_attribute+ constraint_enforcement?`, where each column_attribute is
// PRIMARY KEY, a foreign_key_column_attribute (`[CONSTRAINT id] REFERENCES …`),
// HIDDEN, or NOT NULL. Attributes may appear in any order; the loop ends at the
// first token that is not an attribute introducer. A trailing
// constraint_enforcement (`[NOT] ENFORCED`) is captured into col.Enforced.
func (p *Parser) parseColumnAttributes(col *ast.ColumnDef) error {
	for {
		switch {
		case p.cur.Type == kwPRIMARY:
			p.advance() // PRIMARY
			if _, err := p.expect(kwKEY); err != nil {
				return err
			}
			col.PrimaryKey = true
		case p.cur.Type == kwHIDDEN:
			p.advance()
			col.Hidden = true
		case p.cur.Type == kwNOT && p.peekNext().Type == kwNULL:
			p.advance() // NOT
			p.advance() // NULL
			col.NotNull = true
		case p.cur.Type == kwCONSTRAINT || p.cur.Type == kwREFERENCES:
			// foreign_key_column_attribute: opt_constraint_identity? foreign_key_reference.
			if p.cur.Type == kwCONSTRAINT {
				p.advance() // CONSTRAINT
				nameTok, err := p.expectIdentifier()
				if err != nil {
					return err
				}
				col.FKName = nameTok.Str
			}
			fk, err := p.parseForeignKeyReference()
			if err != nil {
				return err
			}
			col.ForeignKey = fk
		default:
			// constraint_enforcement? at the end of the attribute run.
			if enf, ok := p.tryParseEnforcement(); ok {
				col.Enforced = enf
			}
			return nil
		}
	}
}

// parseTableConstraint parses a table_constraint_definition:
//
//	primary_key_spec
//	| table_constraint_spec
//	| identifier identifier table_constraint_spec   (CONSTRAINT name CHECK/FOREIGN)
//
// In practice the named form is `CONSTRAINT name {CHECK|FOREIGN KEY} …` — the
// grammar's `identifier identifier` covers `CONSTRAINT <name>`. cur is at
// PRIMARY / FOREIGN / CHECK / CONSTRAINT.
func (p *Parser) parseTableConstraint() (*ast.TableConstraint, error) {
	start := p.cur.Loc

	// CONSTRAINT name {PRIMARY KEY | CHECK | FOREIGN KEY} …
	name := ""
	if p.cur.Type == kwCONSTRAINT {
		p.advance() // CONSTRAINT
		nameTok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		name = nameTok.Str
	}

	switch p.cur.Type {
	case kwPRIMARY:
		con, err := p.parsePrimaryKeySpec(start)
		if err != nil {
			return nil, err
		}
		con.Name = name
		return con, nil
	case kwFOREIGN:
		con, err := p.parseForeignKeyConstraintSpec(start)
		if err != nil {
			return nil, err
		}
		con.Name = name
		return con, nil
	case kwCHECK:
		con, err := p.parseCheckConstraintSpec(start)
		if err != nil {
			return nil, err
		}
		con.Name = name
		return con, nil
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// parsePrimaryKeySpec parses a primary_key_spec:
//
//	PRIMARY KEY primary_key_element_list constraint_enforcement? opt_options_list?
func (p *Parser) parsePrimaryKeySpec(start ast.Loc) (*ast.TableConstraint, error) {
	p.advance() // PRIMARY
	if _, err := p.expect(kwKEY); err != nil {
		return nil, err
	}
	keys, err := p.parseKeyPartList(true /*pkElement*/)
	if err != nil {
		return nil, err
	}
	con := &ast.TableConstraint{Kind: ast.ConstraintPrimaryKey, KeyParts: keys, Loc: ast.Loc{Start: start.Start}}
	if enf, ok := p.tryParseEnforcement(); ok {
		con.Enforced = enf
	}
	if p.cur.Type == kwOPTIONS {
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		con.Options = opts
	}
	con.Loc.End = p.prev.Loc.End
	return con, nil
}

// parseForeignKeyConstraintSpec parses the FOREIGN KEY alternative of
// table_constraint_spec:
//
//	FOREIGN KEY column_list foreign_key_reference constraint_enforcement?
//	  opt_options_list?
func (p *Parser) parseForeignKeyConstraintSpec(start ast.Loc) (*ast.TableConstraint, error) {
	p.advance() // FOREIGN
	if _, err := p.expect(kwKEY); err != nil {
		return nil, err
	}
	cols, err := p.parseColumnList()
	if err != nil {
		return nil, err
	}
	fk, err := p.parseForeignKeyReference()
	if err != nil {
		return nil, err
	}
	con := &ast.TableConstraint{
		Kind:       ast.ConstraintForeignKey,
		Columns:    cols,
		ForeignKey: fk,
		Loc:        ast.Loc{Start: start.Start},
	}
	if enf, ok := p.tryParseEnforcement(); ok {
		con.Enforced = enf
	}
	if p.cur.Type == kwOPTIONS {
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		con.Options = opts
	}
	con.Loc.End = p.prev.Loc.End
	return con, nil
}

// parseCheckConstraintSpec parses the CHECK alternative of table_constraint_spec:
//
//	CHECK '(' expression ')' constraint_enforcement? opt_options_list?
func (p *Parser) parseCheckConstraintSpec(start ast.Loc) (*ast.TableConstraint, error) {
	p.advance() // CHECK
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	con := &ast.TableConstraint{Kind: ast.ConstraintCheck, Check: expr, Loc: ast.Loc{Start: start.Start}}
	if enf, ok := p.tryParseEnforcement(); ok {
		con.Enforced = enf
	}
	if p.cur.Type == kwOPTIONS {
		opts, err := p.parseOptionsList()
		if err != nil {
			return nil, err
		}
		con.Options = opts
	}
	con.Loc.End = p.prev.Loc.End
	return con, nil
}

// parseForeignKeyReference parses a foreign_key_reference:
//
//	REFERENCES path_expression column_list opt_foreign_key_match?
//	  opt_foreign_key_action?
func (p *Parser) parseForeignKeyReference() (*ast.ForeignKeyRef, error) {
	refTok, err := p.expect(kwREFERENCES)
	if err != nil {
		return nil, err
	}
	table, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	cols, err := p.parseColumnList()
	if err != nil {
		return nil, err
	}
	fk := &ast.ForeignKeyRef{
		Table:   table,
		Columns: cols,
		Loc:     ast.Loc{Start: refTok.Loc.Start},
	}

	// opt_foreign_key_match?  — MATCH {SIMPLE | FULL | NOT DISTINCT}.
	if p.cur.Type == kwMATCH {
		p.advance() // MATCH
		switch p.cur.Type {
		case kwSIMPLE:
			p.advance()
			fk.Match = "SIMPLE"
		case kwFULL:
			p.advance()
			fk.Match = "FULL"
		case kwNOT:
			p.advance() // NOT
			if _, err := p.expect(kwDISTINCT); err != nil {
				return nil, err
			}
			fk.Match = "NOT DISTINCT"
		default:
			return nil, p.syntaxErrorAtCur()
		}
	}

	// opt_foreign_key_action?  — ON UPDATE / ON DELETE in either order.
	if err := p.parseForeignKeyActions(fk); err != nil {
		return nil, err
	}

	fk.Loc.End = p.prev.Loc.End
	return fk, nil
}

// parseForeignKeyActions parses the opt_foreign_key_action pair (ON UPDATE and/or
// ON DELETE, in either order) into fk. The grammar allows update-then-delete or
// delete-then-update, each optional.
func (p *Parser) parseForeignKeyActions(fk *ast.ForeignKeyRef) error {
	for p.cur.Type == kwON {
		p.advance() // ON
		switch p.cur.Type {
		case kwUPDATE:
			p.advance()
			act, err := p.parseForeignKeyAction()
			if err != nil {
				return err
			}
			fk.OnUpdate = act
		case kwDELETE:
			p.advance()
			act, err := p.parseForeignKeyAction()
			if err != nil {
				return err
			}
			fk.OnDelete = act
		default:
			return p.syntaxErrorAtCur()
		}
	}
	return nil
}

// parseForeignKeyAction parses a foreign_key_action: NO ACTION | RESTRICT |
// CASCADE | SET NULL.
func (p *Parser) parseForeignKeyAction() (ast.ForeignKeyAction, error) {
	switch p.cur.Type {
	case kwNO:
		p.advance() // NO
		if _, err := p.expect(kwACTION); err != nil {
			return 0, err
		}
		return ast.FKActionNoAction, nil
	case kwRESTRICT:
		p.advance()
		return ast.FKActionRestrict, nil
	case kwCASCADE:
		p.advance()
		return ast.FKActionCascade, nil
	case kwSET:
		p.advance() // SET
		if _, err := p.expect(kwNULL); err != nil {
			return 0, err
		}
		return ast.FKActionSetNull, nil
	default:
		return 0, p.syntaxErrorAtCur()
	}
}

// parseSpannerTableOptions parses opt_spanner_table_options plus the comma-led
// Spanner trailing clauses that follow the primary key:
//
//	spanner_primary_key:  PRIMARY KEY primary_key_element_list
//	, INTERLEAVE IN PARENT maybe_dashed_path foreign_key_on_delete   (interleave)
//	, ROW DELETION POLICY ( expression )                             (TTL)
//
// cur is at the trailing PRIMARY keyword. The Spanner grammar threads INTERLEAVE
// and ROW DELETION POLICY as comma-prefixed trailing clauses after the PK (the
// comma is REQUIRED — oracle: `PRIMARY KEY (x) ROW DELETION POLICY (…)` without
// the comma rejects, with it accepts). Either may appear (in either order);
// returns the PK key parts, the optional interleave clause, and the optional
// row-deletion-policy expression.
func (p *Parser) parseSpannerTableOptions() ([]*ast.KeyPart, *ast.InterleaveClause, ast.Node, error) {
	p.advance() // PRIMARY
	if _, err := p.expect(kwKEY); err != nil {
		return nil, nil, nil, err
	}
	keys, err := p.parseKeyPartList(true /*pkElement*/)
	if err != nil {
		return nil, nil, nil, err
	}

	var interleave *ast.InterleaveClause
	var rowDeletion ast.Node
	// Comma-led trailing clauses. Each may appear AT MOST ONCE: the legacy grammar
	// threads opt_spanner_interleave_in_parent_clause? (zero-or-one) and a single
	// ROW DELETION POLICY, and the emulator rejects a repeat (oracle: two INTERLEAVE
	// clauses → syntax error). A duplicate is reported as a syntax error at the
	// repeated keyword.
	//
	// ORDER is fixed: INTERLEAVE must precede ROW DELETION POLICY. The live Spanner
	// grammar accepts `… , INTERLEAVE …, ROW DELETION POLICY …` but REJECTS the
	// reverse `… , ROW DELETION POLICY …, INTERLEAVE …` (reproduced both directions
	// against the emulator; the checked-in .g4 is stale for the comma-led RDP, so
	// the live oracle is authoritative here). We enforce it by rejecting a
	// `, INTERLEAVE` once a ROW DELETION POLICY has already been consumed (error
	// positioned at the out-of-order INTERLEAVE keyword), mirroring the
	// duplicate-clause guards.
	for p.cur.Type == int(',') {
		switch p.peekNext().Type {
		case kwINTERLEAVE:
			if interleave != nil {
				p.advance() // ',' — position the error at the duplicate INTERLEAVE
				return nil, nil, nil, p.syntaxErrorAtCur()
			}
			if rowDeletion != nil {
				// INTERLEAVE after ROW DELETION POLICY is out of order — reject at the
				// INTERLEAVE keyword.
				p.advance() // ','
				return nil, nil, nil, p.syntaxErrorAtCur()
			}
			il, err := p.parseInterleaveInParent()
			if err != nil {
				return nil, nil, nil, err
			}
			interleave = il
		case kwROW:
			if rowDeletion != nil {
				p.advance() // ',' — position the error at the duplicate ROW
				return nil, nil, nil, p.syntaxErrorAtCur()
			}
			p.advance() // ','
			expr, err := p.parseRowDeletionPolicyBody()
			if err != nil {
				return nil, nil, nil, err
			}
			rowDeletion = expr
		default:
			// A comma not introducing a known trailing clause ends the option run;
			// leave it for the caller's next clause / EOF check.
			return keys, interleave, rowDeletion, nil
		}
	}
	return keys, interleave, rowDeletion, nil
}

// parseInterleaveInParent parses a comma-led opt_spanner_interleave_in_parent_clause:
// `, INTERLEAVE IN PARENT maybe_dashed_path_expression foreign_key_on_delete`.
// cur is at the ','.
func (p *Parser) parseInterleaveInParent() (*ast.InterleaveClause, error) {
	start := p.cur.Loc
	p.advance() // ','
	p.advance() // INTERLEAVE
	if _, err := p.expect(kwIN); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwPARENT); err != nil {
		return nil, err
	}
	parent, err := p.parseTablePath()
	if err != nil {
		return nil, err
	}
	// foreign_key_on_delete: ON DELETE action (required by the grammar here).
	if _, err := p.expect(kwON); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwDELETE); err != nil {
		return nil, err
	}
	act, err := p.parseForeignKeyAction()
	if err != nil {
		return nil, err
	}
	return &ast.InterleaveClause{
		Parent:   parent,
		OnDelete: act,
		Loc:      ast.Loc{Start: start.Start, End: p.prev.Loc.End},
	}, nil
}

// parseKeyPartList parses a parenthesized key-part list. When pkElement is true
// it follows primary_key_element_list (`( primary_key_element (, …)* )?` — the
// element is `identifier asc_or_desc? null_order?`, and the list may be empty).
// When false it follows index_order_by_and_options (each element is a general
// `expression collate_clause? asc_or_desc? null_order? opt_options_list?`).
func (p *Parser) parseKeyPartList(pkElement bool) ([]*ast.KeyPart, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	var parts []*ast.KeyPart

	// Empty list `()` (legal for a Spanner PRIMARY KEY with no columns —
	// oracle-confirmed `PRIMARY KEY ()` accepts).
	if p.cur.Type == int(')') {
		p.advance()
		return parts, nil
	}

	for {
		var part *ast.KeyPart
		var err error
		if pkElement {
			part, err = p.parsePrimaryKeyElement()
		} else {
			part, err = p.parseIndexKeyElement()
		}
		if err != nil {
			return nil, err
		}
		parts = append(parts, part)
		if _, ok := p.match(int(',')); ok {
			continue
		}
		break
	}

	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	return parts, nil
}

// parsePrimaryKeyElement parses a primary_key_element: `identifier asc_or_desc?
// null_order?`.
func (p *Parser) parsePrimaryKeyElement() (*ast.KeyPart, error) {
	nameTok, err := p.expectIdentifier()
	if err != nil {
		return nil, err
	}
	name, err := p.identifierText(nameTok)
	if err != nil {
		return nil, err
	}
	part := &ast.KeyPart{Name: name, Loc: nameTok.Loc}
	part.Direction = p.matchAscDesc()
	part.NullOrder = p.matchNullOrder()
	part.Loc.End = p.prev.Loc.End
	return part, nil
}

// parseColumnList parses a column_list: `( identifier (, identifier)* )`.
// Requires at least one identifier.
func (p *Parser) parseColumnList() ([]string, error) {
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	var cols []string
	for {
		tok, err := p.expectIdentifier()
		if err != nil {
			return nil, err
		}
		name, err := p.identifierText(tok)
		if err != nil {
			return nil, err
		}
		cols = append(cols, name)
		if _, ok := p.match(int(',')); ok {
			continue
		}
		break
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	return cols, nil
}

// matchAscDesc consumes an asc_or_desc (ASC | DESC) if present and returns its
// keyword, or "" when absent.
func (p *Parser) matchAscDesc() string {
	switch p.cur.Type {
	case kwASC:
		p.advance()
		return "ASC"
	case kwDESC:
		p.advance()
		return "DESC"
	}
	return ""
}

// matchNullOrder consumes a null_order (NULLS FIRST | NULLS LAST) if present and
// returns "FIRST"/"LAST", or "" when absent. It consumes the NULLS keyword ONLY
// when it is immediately followed by FIRST or LAST (the grammar's only valid
// continuations); a bare `NULLS` with no FIRST/LAST is left entirely unconsumed
// so the caller's next expect surfaces the syntax error (otherwise `NULLS` would
// be silently swallowed and `PRIMARY KEY (x NULLS)` would wrongly accept).
func (p *Parser) matchNullOrder() string {
	if p.cur.Type != kwNULLS {
		return ""
	}
	switch p.peekNext().Type {
	case kwFIRST:
		p.advance() // NULLS
		p.advance() // FIRST
		return "FIRST"
	case kwLAST:
		p.advance() // NULLS
		p.advance() // LAST
		return "LAST"
	}
	// Bare NULLS (no FIRST/LAST): do NOT consume it; leave it for the caller.
	return ""
}

// tryParseEnforcement consumes a constraint_enforcement (`[NOT] ENFORCED`) if
// present and returns its rendered form ("ENFORCED" | "NOT ENFORCED", true), or
// ("", false) when absent.
func (p *Parser) tryParseEnforcement() (string, bool) {
	if p.cur.Type == kwNOT && p.peekNext().Type == kwENFORCED {
		p.advance() // NOT
		p.advance() // ENFORCED
		return "NOT ENFORCED", true
	}
	if p.cur.Type == kwENFORCED {
		p.advance()
		return "ENFORCED", true
	}
	return "", false
}

// parseIfNotExists consumes an opt_if_not_exists (IF NOT EXISTS) if present.
func (p *Parser) parseIfNotExists() (bool, error) {
	if p.cur.Type != kwIF {
		return false, nil
	}
	p.advance() // IF
	if _, err := p.expect(kwNOT); err != nil {
		return false, err
	}
	if _, err := p.expect(kwEXISTS); err != nil {
		return false, err
	}
	return true, nil
}

// parseIfExists consumes an opt_if_exists (IF EXISTS) if present.
func (p *Parser) parseIfExists() (bool, error) {
	if p.cur.Type != kwIF {
		return false, nil
	}
	p.advance() // IF
	if _, err := p.expect(kwEXISTS); err != nil {
		return false, err
	}
	return true, nil
}

// parseOptionsList parses an opt_options_list body `OPTIONS options_list` where
// cur is the OPTIONS keyword:
//
//	OPTIONS '(' (options_entry (, options_entry)*)? ')'
//
// Each options_entry is `identifier_in_hints options_assignment_operator
// expression_or_proto`. The empty `OPTIONS ()` form is legal.
func (p *Parser) parseOptionsList() ([]*ast.OptionsEntry, error) {
	p.advance() // OPTIONS
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	var entries []*ast.OptionsEntry
	if p.cur.Type == int(')') {
		p.advance()
		return entries, nil
	}
	for {
		entry, err := p.parseOptionsEntry()
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
		if _, ok := p.match(int(',')); ok {
			continue
		}
		break
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	return entries, nil
}

// parseOptionsEntry parses one options_entry: `identifier_in_hints
// options_assignment_operator expression_or_proto`. The key is any identifier or
// keyword-as-identifier (identifier_in_hints is permissive); the operator is one
// of `=`, `+=`, `-=`; the value is an expression or the bare PROTO keyword.
func (p *Parser) parseOptionsEntry() (*ast.OptionsEntry, error) {
	if !isAnyKeywordIdentifier(p.cur.Type) {
		return nil, p.syntaxErrorAtCur()
	}
	keyTok := p.advance()
	key := keyTok.Str
	if key == "" {
		key = p.tokenSource(keyTok)
	}
	entry := &ast.OptionsEntry{Name: key, Loc: keyTok.Loc}

	// options_assignment_operator: = | += | -=.
	op, err := p.parseOptionsAssignmentOperator()
	if err != nil {
		return nil, err
	}
	entry.Op = op

	// expression_or_proto: PROTO | expression.
	if p.cur.Type == kwPROTO {
		protoTok := p.advance()
		entry.Value = &ast.Literal{Kind: ast.LitString, Value: "PROTO", Loc: protoTok.Loc}
	} else {
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		entry.Value = expr
	}
	entry.Loc.End = p.prev.Loc.End
	return entry, nil
}

// parseOptionsAssignmentOperator consumes the `=`, `+=`, or `-=` operator in an
// options_entry and returns its spelling. The lexer emits `+=` / `-=` as the
// tokPlusEqual / tokSubEqual tokens (see tokens.go); `=` is the single-char '='.
func (p *Parser) parseOptionsAssignmentOperator() (string, error) {
	switch p.cur.Type {
	case int('='):
		p.advance()
		return "=", nil
	case tokPlusEqual:
		p.advance()
		return "+=", nil
	case tokMinusEqual:
		p.advance()
		return "-=", nil
	default:
		return "", p.syntaxErrorAtCur()
	}
}

// parsePartitionByNoHint parses a partition_by_clause_prefix_no_hint: `PARTITION
// BY expression (, expression)*`. cur is at PARTITION.
func (p *Parser) parsePartitionByNoHint() ([]ast.Node, error) {
	p.advance() // PARTITION
	if _, err := p.expect(kwBY); err != nil {
		return nil, err
	}
	return p.parseExprCommaList()
}

// parseClusterByNoHint parses a cluster_by_clause_prefix_no_hint: `CLUSTER BY
// expression (, expression)*`. cur is at CLUSTER.
func (p *Parser) parseClusterByNoHint() ([]ast.Node, error) {
	p.advance() // CLUSTER
	if _, err := p.expect(kwBY); err != nil {
		return nil, err
	}
	return p.parseExprCommaList()
}

// parseRowDeletionPolicyBody parses an opt_ttl_clause: `ROW DELETION POLICY '('
// expression ')'`. cur is at ROW. Returns the policy expression.
func (p *Parser) parseRowDeletionPolicyBody() (ast.Node, error) {
	p.advance() // ROW
	if _, err := p.expect(kwDELETION); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwPOLICY); err != nil {
		return nil, err
	}
	if _, err := p.expect(int('(')); err != nil {
		return nil, err
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, err
	}
	return expr, nil
}

// skipWithConnection consumes a `WITH CONNECTION <path>` clause (BigQuery
// with_connection_clause). The connection-clause AST is dialect-node territory;
// the core CREATE TABLE only needs to step over it so a following OPTIONS / AS
// query parses. cur is at WITH (with CONNECTION next).
func (p *Parser) skipWithConnection() error {
	p.advance() // WITH
	p.advance() // CONNECTION
	// connection_clause is `connection_name` (a path) or DEFAULT. Consume a path
	// or the DEFAULT keyword.
	if p.cur.Type == kwDEFAULT {
		p.advance()
		return nil
	}
	if _, err := p.parseTablePath(); err != nil {
		return err
	}
	return nil
}

// skipBalancedParens consumes a `( … )` run with balanced nesting starting at the
// current token (which must be '('). Used for sub-clauses whose body the core
// node does not need to model structurally (e.g. IDENTITY(...)). It always
// advances at least the opening paren.
func (p *Parser) skipBalancedParens() error {
	if p.cur.Type != int('(') {
		return p.syntaxErrorAtCur()
	}
	depth := 0
	for p.cur.Type != tokEOF {
		switch p.cur.Type {
		case int('('):
			depth++
		case int(')'):
			depth--
			if depth == 0 {
				p.advance()
				return nil
			}
		}
		p.advance()
	}
	return &ParseError{Loc: p.cur.Loc, Msg: "syntax error at end of input"}
}
