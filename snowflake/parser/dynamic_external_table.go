package parser

import (
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Table-variant DDL — CREATE / ALTER DYNAMIC / EXTERNAL / EVENT TABLE (T4.4)
// ---------------------------------------------------------------------------
//
// DYNAMIC / EXTERNAL / EVENT tables each carry a large, version-growing
// configuration vocabulary that the legacy ANTLR grammar enumerates as an
// already-stale subset (create_dynamic_table lacks REFRESH_MODE = ADAPTIVE |
// CUSTOM_INCREMENTAL, SCHEDULER, INITIALIZATION_WAREHOUSE, the ICEBERG variant's
// EXTERNAL_VOLUME / CATALOG / BASE_LOCATION / ICEBERG_VERSION, the REQUIRE USER
// and IMMUTABLE/FROZEN WHERE clauses, the REFRESH USING body, ...;
// create_external_table lacks the USING TEMPLATE form). Matching the merged STAGE
// (T4.1) / FILE FORMAT (T4.2) / COPY (T5.2) / pipeline (T4.3) approach, every such
// config parameter that is a `KEY = <value>` pair is captured open-ended
// (ast.CopyOption); only the structurally distinct anchors are modeled
// explicitly: the optional ICEBERG modifier, the column list, CLUSTER BY, CLONE,
// REQUIRE USER, the IMMUTABLE/FROZEN WHERE predicate, EXTERNAL TABLE's
// LOCATION = @stage / PARTITION BY / USING TEMPLATE clauses, and the terminal
// AS <query> / REFRESH USING (<dml>) body. The catalog/semantic layer, not the
// parser, validates that an option is real and legal.
//
// (DROP / UNDROP for DYNAMIC / EXTERNAL TABLE are handled by drop.go; ALTER EVENT
// TABLE is handled by ALTER TABLE per the legacy grammar, so there is no dedicated
// ALTER EVENT TABLE path here.)

// ---------------------------------------------------------------------------
// CREATE DYNAMIC TABLE
// ---------------------------------------------------------------------------

// parseCreateDynamicTableStmt parses
//
//	CREATE [ OR REPLACE ] [ TRANSIENT ] DYNAMIC [ ICEBERG ] TABLE [ IF NOT EXISTS ] <name>
//	  [ ( <col> [ <type> ] [ , ... ] ) ]
//	  [ <config options> ] [ CLUSTER BY ( ... ) ] [ REQUIRE USER ]
//	  [ { IMMUTABLE | FROZEN } WHERE ( <expr> ) ]
//	  { AS <query> | REFRESH USING ( <dml> ) }
//
//	CREATE [ OR REPLACE ] [ TRANSIENT ] DYNAMIC [ ICEBERG ] TABLE <name>
//	  CLONE <source> [ { AT | BEFORE } ( ... ) ]
//
// The CREATE keyword and the optional OR REPLACE / TRANSIENT modifiers have
// already been consumed by parseCreateStmt; start is the Loc of the CREATE token,
// transient records a preceding TRANSIENT, and cur is the DYNAMIC keyword. The
// two-keyword `DYNAMIC TABLE` (with an optional ICEBERG between them) is consumed
// here.
func (p *Parser) parseCreateDynamicTableStmt(start ast.Loc, orReplace, transient bool) (ast.Node, error) {
	p.advance() // consume DYNAMIC

	stmt := &ast.CreateDynamicTableStmt{
		OrReplace: orReplace,
		Transient: transient,
		Loc:       ast.Loc{Start: start.Start},
	}

	// Optional ICEBERG modifier (DYNAMIC ICEBERG TABLE). ICEBERG is not a reserved
	// keyword token, so it is matched by identifier text (and only when TABLE
	// follows, so a table literally named "iceberg" is not misread as the modifier).
	if p.curIsWord("ICEBERG") && p.peekNext().Type == kwTABLE {
		p.advance() // consume ICEBERG
		stmt.Iceberg = true
	}

	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}

	if err := p.parseIfNotExistsInto(&stmt.IfNotExists); err != nil {
		return nil, err
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// CLONE form: CREATE DYNAMIC TABLE <name> CLONE <source> [AT|BEFORE (...)].
	if p.cur.Type == kwCLONE {
		clone, err := p.parseCloneSource()
		if err != nil {
			return nil, err
		}
		stmt.Clone = clone
		stmt.Loc.End = p.prev.Loc.End
		return stmt, nil
	}

	// Optional materialized column list: ( <col> [ <type> ] [, ...] ).
	if p.cur.Type == '(' {
		cols, err := p.parseMaterializedColumnList()
		if err != nil {
			return nil, err
		}
		stmt.Columns = cols
	}

	// Config options, CLUSTER BY, REQUIRE USER, and the IMMUTABLE/FROZEN WHERE
	// clause may appear before the terminal AS / REFRESH body, in any order. AS /
	// REFRESH are the structural anchors that terminate the clause run.
	//
	// REQUIRE / FROZEN are not reserved keyword tokens (IMMUTABLE is), so the
	// contextual clauses they introduce are detected by identifier text plus a
	// lookahead (REQUIRE USER, { IMMUTABLE | FROZEN } WHERE). They are matched
	// before the open-ended option branch so they are not swallowed as KEY=value
	// options.
	for {
		switch {
		case p.cur.Type == kwCLUSTER:
			if err := p.parseClusterByInto(&stmt.ClusterBy, &stmt.Linear); err != nil {
				return nil, err
			}
			continue
		case p.curIsWord("REQUIRE") && p.peekNext().Type == kwUSER:
			// REQUIRE USER — a bare flag clause (no value).
			p.advance() // consume REQUIRE
			p.advance() // consume USER
			stmt.RequireUser = true
			continue
		case p.startsImmutableWhere():
			kind, expr, err := p.parseImmutableWhere()
			if err != nil {
				return nil, err
			}
			stmt.ImmutableKind = kind
			stmt.ImmutableWhere = expr
			continue
		case p.cur.Type == kwAS || p.cur.Type == kwREFRESH:
			// fall out of the loop to the body parse below.
		default:
			if p.startsDynamicTableOption() {
				opt, err := p.parseCopyOption()
				if err != nil {
					return nil, err
				}
				stmt.Options = append(stmt.Options, opt)
				continue
			}
		}
		break
	}

	// Terminal body: AS <query> | REFRESH USING ( <dml> ).
	switch p.cur.Type {
	case kwAS:
		p.advance() // consume AS
		query, err := p.parseQueryExpr()
		if err != nil {
			return nil, err
		}
		stmt.AsQuery = query
	case kwREFRESH:
		// REFRESH USING ( <dml_statement> ).
		p.advance() // consume REFRESH
		if _, err := p.expect(kwUSING); err != nil {
			return nil, err
		}
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		body, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		stmt.RefreshUsing = body
	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseMaterializedColumnList parses the optional ( <col> [ <type> ] [, ...] )
// list of a dynamic table. Per the legacy grammar (materialized_col_decl) each
// entry is a column name with an OPTIONAL data type — a dynamic-table column may
// be name-only when its type is inferred from the AS query (the corpus shows
// `order_day` with no type). MASKING POLICY / TAG / COMMENT column properties are
// captured by reusing the shared column-options reader, mirroring CREATE TABLE.
// On entry cur is '('.
func (p *Parser) parseMaterializedColumnList() ([]*ast.ColumnDef, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	var cols []*ast.ColumnDef
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		col, err := p.parseMaterializedColumnDef()
		if err != nil {
			return nil, err
		}
		cols = append(cols, col)

		if p.cur.Type == ',' {
			p.advance() // consume ','
			continue
		}
		break
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return cols, nil
}

// parseMaterializedColumnDef parses one dynamic-table column: a name, an optional
// data type, and the optional column properties (WITH MASKING POLICY / [WITH] TAG
// / COMMENT) handled by the shared parseColumnOptions reader. The data type is
// taken unless the next token already begins a column property, ',' or ')'.
func (p *Parser) parseMaterializedColumnDef() (*ast.ColumnDef, error) {
	name, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	col := &ast.ColumnDef{
		Name: name,
		Loc:  ast.Loc{Start: name.Loc.Start},
	}

	// The data type is optional. It is absent when the next token closes the
	// column (',' / ')') or begins a column property (WITH / MASKING / TAG /
	// COMMENT). Otherwise read a data type.
	if p.cur.Type != ',' && p.cur.Type != ')' && p.cur.Type != tokEOF &&
		p.cur.Type != kwWITH && p.cur.Type != kwMASKING && p.cur.Type != kwTAG &&
		p.cur.Type != kwCOMMENT {
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

// startsDynamicTableOption reports whether cur begins an open-ended dynamic-table
// option. It is the COPY-option predicate minus the keywords/words that anchor
// dedicated clauses (CLUSTER / AS / REFRESH, and the IMMUTABLE keyword and the
// contextual REQUIRE / FROZEN words), so those are never swallowed as option
// names.
func (p *Parser) startsDynamicTableOption() bool {
	switch p.cur.Type {
	case kwCLUSTER, kwAS, kwREFRESH, kwIMMUTABLE:
		return false
	}
	if p.curIsWord("REQUIRE") || p.curIsWord("FROZEN") {
		return false
	}
	return p.startsCopyOption()
}

// startsImmutableWhere reports whether cur begins a { IMMUTABLE | FROZEN } WHERE
// clause: the IMMUTABLE keyword (or the contextual FROZEN word) immediately
// followed by WHERE.
func (p *Parser) startsImmutableWhere() bool {
	if p.cur.Type == kwIMMUTABLE && p.peekNext().Type == kwWHERE {
		return true
	}
	return p.curIsWord("FROZEN") && p.peekNext().Type == kwWHERE
}

// parseImmutableWhere parses a { IMMUTABLE | FROZEN } WHERE ( <expr> ) clause,
// returning the uppercased selector ("IMMUTABLE" or "FROZEN") and the predicate
// expression. The caller must have confirmed startsImmutableWhere.
func (p *Parser) parseImmutableWhere() (string, ast.Node, error) {
	kind := strings.ToUpper(p.advance().Str) // consume IMMUTABLE / FROZEN
	if _, err := p.expect(kwWHERE); err != nil {
		return "", nil, err
	}
	if _, err := p.expect('('); err != nil {
		return "", nil, err
	}
	expr, err := p.parseExpr()
	if err != nil {
		return "", nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return "", nil, err
	}
	return kind, expr, nil
}

// parseClusterByInto parses a CLUSTER BY [ LINEAR ] ( <expr> [, ...] ) clause into
// *exprs and *linear. On entry cur is the CLUSTER keyword. Mirrors the CLUSTER BY
// handling in parseTableProperties (CREATE TABLE).
func (p *Parser) parseClusterByInto(exprs *[]ast.Node, linear *bool) error {
	p.advance() // consume CLUSTER
	if _, err := p.expect(kwBY); err != nil {
		return err
	}
	if p.cur.Type == kwLINEAR {
		p.advance() // consume LINEAR
		*linear = true
	}
	if _, err := p.expect('('); err != nil {
		return err
	}
	list, err := p.parseExprList()
	if err != nil {
		return err
	}
	if _, err := p.expect(')'); err != nil {
		return err
	}
	*exprs = list
	return nil
}

// ---------------------------------------------------------------------------
// ALTER DYNAMIC TABLE
// ---------------------------------------------------------------------------

// parseAlterDynamicTableStmt parses ALTER DYNAMIC TABLE [ IF EXISTS ] <name> <action>.
// The ALTER keyword has already been consumed; cur is the DYNAMIC keyword (the
// TABLE keyword that follows is consumed here).
//
//	{ SUSPEND | RESUME }
//	REFRESH [ COPY SESSION ]
//	RENAME TO <new_name>
//	SWAP WITH <other_name>
//	SET <settable params>
//	UNSET <property> [ , ... ]
//	SET TAG <tag> = '<value>' [ , ... ]
//	UNSET TAG <tag> [ , ... ]
//	CLUSTER BY ( <expr> [ , ... ] )
//	DROP CLUSTERING KEY
func (p *Parser) parseAlterDynamicTableStmt() (ast.Node, error) {
	altTok := p.advance() // consume DYNAMIC
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}
	stmt := &ast.AlterDynamicTableStmt{Loc: ast.Loc{Start: altTok.Loc.Start}}

	if err := p.parseIfExistsInto(&stmt.IfExists); err != nil {
		return nil, err
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	switch p.cur.Type {
	case kwSUSPEND:
		p.advance()
		stmt.Action = ast.AlterDynamicTableSuspend

	case kwRESUME:
		p.advance()
		stmt.Action = ast.AlterDynamicTableResume

	case kwREFRESH:
		// REFRESH [ COPY SESSION ].
		p.advance() // consume REFRESH
		stmt.Action = ast.AlterDynamicTableRefresh
		if p.cur.Type == kwCOPY && p.peekNext().Type == kwSESSION {
			p.advance() // consume COPY
			p.advance() // consume SESSION
			stmt.RefreshCopySession = true
		}

	case kwRENAME:
		// RENAME TO <new_name>.
		p.advance() // consume RENAME
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		newName, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterDynamicTableRename
		stmt.NewName = newName

	case kwSWAP:
		// SWAP WITH <other_name>.
		p.advance() // consume SWAP
		if _, err := p.expect(kwWITH); err != nil {
			return nil, err
		}
		other, err := p.parseObjectName()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterDynamicTableSwap
		stmt.NewName = other

	case kwSET:
		p.advance() // consume SET
		if p.cur.Type == kwTAG {
			tags, err := p.parseTagSetList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterDynamicTableSetTag
			stmt.Tags = tags
		} else {
			opts, err := p.parseRequiredOptions()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterDynamicTableSet
			stmt.Options = opts
		}

	case kwUNSET:
		p.advance() // consume UNSET
		if p.cur.Type == kwTAG {
			names, err := p.parseUnsetTagNameList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterDynamicTableUnsetTag
			stmt.UnsetTags = names
		} else {
			props, err := p.parseUnsetPropertyList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterDynamicTableUnset
			stmt.UnsetProps = props
		}

	case kwCLUSTER:
		// CLUSTER BY ( exprs ).
		if err := p.parseClusterByInto(&stmt.ClusterBy, &stmt.Linear); err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterDynamicTableClusterBy

	case kwDROP:
		// DROP CLUSTERING KEY.
		p.advance() // consume DROP
		if _, err := p.expect(kwCLUSTERING); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwKEY); err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterDynamicTableDropClusteringKey

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// CREATE EXTERNAL TABLE
// ---------------------------------------------------------------------------

// parseCreateExternalTableStmt parses
//
//	CREATE [ OR REPLACE ] EXTERNAL TABLE [ IF NOT EXISTS ] <name>
//	  ( <col> <type> AS ( <expr> ) [ inlineConstraint ] [ , ... ] )
//	  [ <config options> ] [ PARTITION BY ( ... ) ] [ WITH ] LOCATION = @stage
//	  [ <more options> ]
//
//	CREATE [ OR REPLACE ] EXTERNAL TABLE [ IF NOT EXISTS ] <name>
//	  USING TEMPLATE ( <query> )
//	  [ WITH ] LOCATION = @stage [ <options> ]
//
// The CREATE keyword and OR REPLACE have been consumed; cur is the EXTERNAL
// keyword (the TABLE keyword that follows is consumed here). LOCATION = @stage is
// mandatory per the docs; every other clause (REFRESH_ON_CREATE / AUTO_REFRESH /
// PATTERN / PARTITION_TYPE / FILE_FORMAT / TABLE_FORMAT / AWS_SNS_TOPIC /
// INTEGRATION / COMMENT / ...) is captured open-ended, with PARTITION BY, the
// column list / USING TEMPLATE, COPY GRANTS, and [WITH] TAG as the structural
// anchors.
func (p *Parser) parseCreateExternalTableStmt(start ast.Loc, orReplace bool) (ast.Node, error) {
	p.advance() // consume EXTERNAL
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}

	stmt := &ast.CreateExternalTableStmt{
		OrReplace: orReplace,
		Loc:       ast.Loc{Start: start.Start},
	}

	if err := p.parseIfNotExistsInto(&stmt.IfNotExists); err != nil {
		return nil, err
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// Either an explicit ( <external_col_decl_list> ) or the USING TEMPLATE form.
	switch {
	case p.cur.Type == '(':
		cols, err := p.parseExternalColumnList()
		if err != nil {
			return nil, err
		}
		stmt.Columns = cols
	case p.cur.Type == kwUSING:
		// USING TEMPLATE ( <query> ). TEMPLATE is not a reserved keyword token, so
		// it is matched by identifier text.
		p.advance() // consume USING
		if !p.curIsWord("TEMPLATE") {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance() // consume TEMPLATE
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		tmpl, err := p.parseQueryExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
		stmt.UsingTemplate = tmpl
	}

	// Trailing clauses, in any order: PARTITION BY, [WITH] LOCATION = @stage,
	// COPY GRANTS, [WITH] TAG, and the open-ended options. LOCATION is captured
	// specially because its `@stage` value is not a plain option word. A bare WITH
	// (the `WITH? LOCATION` of the legacy grammar) is consumed as a no-op anchor;
	// WITH TAG is routed to the tag clause.
	for {
		switch p.cur.Type {
		case kwPARTITION:
			// PARTITION BY ( exprs ). PARTITION_TYPE = USER_SPECIFIED is a *different*
			// token (kwPARTITION_TYPE), captured as an option, so only a PARTITION
			// followed by BY is the partition-by clause.
			if p.peekNext().Type == kwBY {
				p.advance() // consume PARTITION
				p.advance() // consume BY
				if _, err := p.expect('('); err != nil {
					return nil, err
				}
				exprs, err := p.parseExprList()
				if err != nil {
					return nil, err
				}
				if _, err := p.expect(')'); err != nil {
					return nil, err
				}
				stmt.PartitionBy = exprs
				continue
			}
		case kwLOCATION:
			// [WITH] LOCATION = @stage[/path].
			p.advance() // consume LOCATION
			if _, err := p.expect('='); err != nil {
				return nil, err
			}
			loc, err := p.parseStageRef()
			if err != nil {
				return nil, err
			}
			stmt.Location = loc
			continue
		case kwWITH:
			// A bare WITH preceding LOCATION (legacy `partition_by? WITH? LOCATION`)
			// is a no-op anchor; WITH TAG introduces the tag clause.
			if p.peekNext().Type == kwTAG {
				tags, err := p.parseStreamWithTags()
				if err != nil {
					return nil, err
				}
				stmt.Tags = append(stmt.Tags, tags...)
				continue
			}
			if p.peekNext().Type == kwROW {
				// WITH ROW ACCESS POLICY ... — consume and discard (mirrors CREATE TABLE).
				if err := p.skipRowAccessPolicy(); err != nil {
					return nil, err
				}
				continue
			}
			p.advance() // consume bare WITH (anchors the following LOCATION)
			continue
		case kwTAG:
			tags, err := p.parseTagAssignments()
			if err != nil {
				return nil, err
			}
			stmt.Tags = append(stmt.Tags, tags...)
			continue
		case kwCOPY:
			if p.startsCopyGrants() {
				p.advance() // consume COPY
				p.advance() // consume GRANTS
				stmt.CopyGrants = true
				continue
			}
		}
		if p.startsExternalTableOption() {
			opt, err := p.parseCopyOption()
			if err != nil {
				return nil, err
			}
			stmt.Options = append(stmt.Options, opt)
			continue
		}
		break
	}

	// LOCATION = @stage is mandatory in every documented form (column-list and
	// USING TEMPLATE) and in all three legacy-grammar alternatives, so a statement
	// without it is rejected rather than accepted with a nil Location.
	if stmt.Location == nil {
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// startsExternalTableOption reports whether cur begins an open-ended external-table
// option. It is the COPY-option predicate minus the keywords that anchor dedicated
// clauses (WITH / TAG / PARTITION / LOCATION / COPY), so those are never swallowed
// as option names.
func (p *Parser) startsExternalTableOption() bool {
	switch p.cur.Type {
	case kwWITH, kwTAG, kwPARTITION, kwLOCATION, kwCOPY:
		return false
	}
	return p.startsCopyOption()
}

// parseExternalColumnList parses ( <external_col_decl> [ , ... ] ). On entry cur
// is '('.
func (p *Parser) parseExternalColumnList() ([]*ast.ExternalColumnDef, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	var cols []*ast.ExternalColumnDef
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		col, err := p.parseExternalColumnDef()
		if err != nil {
			return nil, err
		}
		cols = append(cols, col)

		if p.cur.Type == ',' {
			p.advance() // consume ','
			continue
		}
		break
	}

	// An empty column list "()" is not a valid external-table declaration (the
	// legacy grammar's external_table_column_decl_list requires at least one
	// column). The USING TEMPLATE form, which omits the explicit list entirely,
	// takes a different parse path and never reaches here.
	if len(cols) == 0 {
		return nil, p.syntaxErrorAtCur()
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return cols, nil
}

// parseExternalColumnDef parses one external-table column:
//
//	<col_name> <col_type> AS ( <expr> | <id> ) [ NOT NULL | inlineConstraint ]
//
// Per the legacy grammar (external_table_column_decl) the data type is required
// and the AS clause is mandatory; the AS value is an expression that may be
// parenthesized (`AS (value:col::int)`) or bare (`AS TO_DATE(...)`).
func (p *Parser) parseExternalColumnDef() (*ast.ExternalColumnDef, error) {
	name, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	col := &ast.ExternalColumnDef{
		Name: name,
		Loc:  ast.Loc{Start: name.Loc.Start},
	}

	dt, err := p.parseDataType()
	if err != nil {
		return nil, err
	}
	col.DataType = dt

	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}

	// AS ( <expr> ) — parenthesized — or AS <expr> bare. Both reduce to a single
	// expression; a parenthesized form is parsed as a ParenExpr by the shared
	// expression parser, so parseExpr covers both without special-casing.
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	col.Expr = expr

	// Optional inline NOT NULL / PRIMARY KEY / UNIQUE / FOREIGN KEY constraint.
	if p.cur.Type == kwNOT && p.peekNext().Type == kwNULL {
		p.advance() // consume NOT
		p.advance() // consume NULL
		col.NotNull = true
	}
	if p.cur.Type == kwCONSTRAINT || p.cur.Type == kwPRIMARY ||
		p.cur.Type == kwUNIQUE || p.cur.Type == kwFOREIGN {
		ic, err := p.parseInlineConstraint()
		if err != nil {
			return nil, err
		}
		col.Constraint = ic
	}

	col.Loc.End = p.prev.Loc.End
	return col, nil
}

// skipRowAccessPolicy consumes and discards a WITH ROW ACCESS POLICY <name> [ ON
// ( <cols> ) ] clause, mirroring the CREATE TABLE handling. On entry cur is the
// WITH keyword and peekNext is ROW.
func (p *Parser) skipRowAccessPolicy() error {
	p.advance() // consume WITH
	p.advance() // consume ROW
	if _, err := p.expect(kwACCESS); err != nil {
		return err
	}
	if _, err := p.expect(kwPOLICY); err != nil {
		return err
	}
	if _, err := p.parseObjectName(); err != nil {
		return err
	}
	if p.cur.Type == kwON {
		p.advance() // consume ON
		if err := p.skipParenthesized(); err != nil {
			return err
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// ALTER EXTERNAL TABLE
// ---------------------------------------------------------------------------

// parseAlterExternalTableStmt parses ALTER EXTERNAL TABLE [ IF EXISTS ] <name> <action>.
// The ALTER keyword has already been consumed; cur is the EXTERNAL keyword (the
// TABLE keyword that follows is consumed here).
//
//	REFRESH [ '<relative-path>' ]
//	ADD FILES ( '<path>/<file>' [ , ... ] )
//	REMOVE FILES ( '<path>/<file>' [ , ... ] )
//	SET [ AUTO_REFRESH = { TRUE | FALSE } ] [ TAG <tag> = '<value>' [ , ... ] ]
//	UNSET TAG <tag> [ , ... ]
//	ADD PARTITION ( <col> = '<val>' [ , ... ] ) LOCATION '<path>'
//	DROP PARTITION LOCATION '<path>'
//
// Per the legacy grammar the optional IF EXISTS appears in two positions
// (before and after the name) depending on the alternative; both spellings are
// accepted by consuming IF EXISTS in either slot.
func (p *Parser) parseAlterExternalTableStmt() (ast.Node, error) {
	altTok := p.advance() // consume EXTERNAL
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}
	stmt := &ast.AlterExternalTableStmt{Loc: ast.Loc{Start: altTok.Loc.Start}}

	if err := p.parseIfExistsInto(&stmt.IfExists); err != nil {
		return nil, err
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	// IF EXISTS may also follow the name (the ADD/DROP PARTITION alternatives put
	// it there in the legacy grammar). Accept it in this slot too.
	if !stmt.IfExists {
		if err := p.parseIfExistsInto(&stmt.IfExists); err != nil {
			return nil, err
		}
	}

	switch p.cur.Type {
	case kwREFRESH:
		// REFRESH [ '<relative-path>' ].
		p.advance() // consume REFRESH
		stmt.Action = ast.AlterExternalTableRefresh
		if p.cur.Type == tokString {
			tok := p.advance()
			s := tok.Str
			stmt.RefreshPath = &s
		}

	case kwADD:
		p.advance() // consume ADD
		switch p.cur.Type {
		case kwFILES:
			files, err := p.parseFilesList()
			if err != nil {
				return nil, err
			}
			stmt.Action = ast.AlterExternalTableAddFiles
			stmt.Files = files
		case kwPARTITION:
			// ADD PARTITION ( <col> = '<val>' [, ...] ) LOCATION '<path>'.
			p.advance() // consume PARTITION
			parts, err := p.parsePartitionAssignments()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(kwLOCATION); err != nil {
				return nil, err
			}
			locTok, err := p.expect(tokString)
			if err != nil {
				return nil, err
			}
			s := locTok.Str
			stmt.Action = ast.AlterExternalTableAddPartition
			stmt.Partitions = parts
			stmt.Location = &s
		default:
			return nil, p.syntaxErrorAtCur()
		}

	case kwREMOVE:
		// REMOVE FILES ( ... ).
		p.advance() // consume REMOVE
		if p.cur.Type != kwFILES {
			return nil, p.syntaxErrorAtCur()
		}
		files, err := p.parseFilesList()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterExternalTableRemoveFiles
		stmt.Files = files

	case kwDROP:
		// DROP PARTITION LOCATION '<path>'.
		p.advance() // consume DROP
		if _, err := p.expect(kwPARTITION); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwLOCATION); err != nil {
			return nil, err
		}
		locTok, err := p.expect(tokString)
		if err != nil {
			return nil, err
		}
		s := locTok.Str
		stmt.Action = ast.AlterExternalTableDropPartition
		stmt.DropLocation = &s

	case kwSET:
		// SET [ AUTO_REFRESH = ... ] [ TAG <tag> = '<value>' [, ...] ]. Both parts
		// are optional in the legacy grammar but a SET with nothing settable is a
		// syntax error.
		p.advance() // consume SET
		sawSomething := false
		// Open-ended options (AUTO_REFRESH = TRUE/FALSE, ...) up to a TAG clause.
		for p.cur.Type != kwTAG && p.startsCopyOption() {
			opt, err := p.parseCopyOption()
			if err != nil {
				return nil, err
			}
			stmt.Options = append(stmt.Options, opt)
			sawSomething = true
		}
		if p.cur.Type == kwTAG {
			tags, err := p.parseTagSetList()
			if err != nil {
				return nil, err
			}
			stmt.Tags = tags
			sawSomething = true
		}
		if !sawSomething {
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Action = ast.AlterExternalTableSet

	case kwUNSET:
		// UNSET TAG <tag> [, ...]. parseUnsetTagNameList consumes the TAG keyword.
		p.advance() // consume UNSET
		if p.cur.Type != kwTAG {
			return nil, p.syntaxErrorAtCur()
		}
		names, err := p.parseUnsetTagNameList()
		if err != nil {
			return nil, err
		}
		stmt.Action = ast.AlterExternalTableUnsetTag
		stmt.UnsetTags = names

	default:
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseFilesList parses a FILES ( '<f>' [ , '<f>' ... ] ) clause and returns the
// string literals. The FILES keyword is consumed here.
func (p *Parser) parseFilesList() ([]*ast.Literal, error) {
	if _, err := p.expect(kwFILES); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var files []*ast.Literal
	for {
		tok, err := p.expect(tokString)
		if err != nil {
			return nil, err
		}
		files = append(files, &ast.Literal{Kind: ast.LitString, Value: tok.Str, Loc: tok.Loc})
		if p.cur.Type == ',' {
			p.advance() // consume ','
			continue
		}
		break
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return files, nil
}

// parsePartitionAssignments parses ( <col> = '<val>' [ , ... ] ) for an ADD
// PARTITION clause. On entry cur is '('.
func (p *Parser) parsePartitionAssignments() ([]*ast.ExternalTablePartition, error) {
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	var parts []*ast.ExternalTablePartition
	for {
		col, err := p.parseIdent()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect('='); err != nil {
			return nil, err
		}
		valTok, err := p.expect(tokString)
		if err != nil {
			return nil, err
		}
		parts = append(parts, &ast.ExternalTablePartition{Column: col, Value: valTok.Str})
		if p.cur.Type == ',' {
			p.advance() // consume ','
			continue
		}
		break
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return parts, nil
}

// ---------------------------------------------------------------------------
// CREATE EVENT TABLE
// ---------------------------------------------------------------------------

// parseCreateEventTableStmt parses
//
//	CREATE [ OR REPLACE ] EVENT TABLE [ IF NOT EXISTS ] <name>
//	  [ CLUSTER BY ( ... ) ] [ <config options> ] [ COPY GRANTS ]
//	  [ <with_row_access_policy> ] [ [ WITH ] TAG ( ... ) ]
//	  [ [ WITH ] COMMENT = '...' ]
//
// The CREATE keyword and OR REPLACE have been consumed; cur is the EVENT keyword
// (the TABLE keyword that follows is consumed here). An EVENT TABLE has a fixed
// system schema (no user column list); CLUSTER BY is the one structural shape and
// COPY GRANTS / [WITH] TAG are anchors; every other parameter
// (DATA_RETENTION_TIME_IN_DAYS / MAX_DATA_EXTENSION_TIME_IN_DAYS / CHANGE_TRACKING
// / DEFAULT_DDL_COLLATION / COMMENT / ...) is captured open-ended.
func (p *Parser) parseCreateEventTableStmt(start ast.Loc, orReplace bool) (ast.Node, error) {
	p.advance() // consume EVENT
	if _, err := p.expect(kwTABLE); err != nil {
		return nil, err
	}

	stmt := &ast.CreateEventTableStmt{
		OrReplace: orReplace,
		Loc:       ast.Loc{Start: start.Start},
	}

	if err := p.parseIfNotExistsInto(&stmt.IfNotExists); err != nil {
		return nil, err
	}

	name, err := p.parseObjectName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name

	for {
		switch p.cur.Type {
		case kwCLUSTER:
			if err := p.parseClusterByInto(&stmt.ClusterBy, &stmt.Linear); err != nil {
				return nil, err
			}
			continue
		case kwCOPY:
			if p.startsCopyGrants() {
				p.advance() // consume COPY
				p.advance() // consume GRANTS
				stmt.CopyGrants = true
				continue
			}
		case kwTAG:
			tags, err := p.parseTagAssignments()
			if err != nil {
				return nil, err
			}
			stmt.Tags = append(stmt.Tags, tags...)
			continue
		case kwWITH:
			switch p.peekNext().Type {
			case kwTAG:
				tags, err := p.parseStreamWithTags()
				if err != nil {
					return nil, err
				}
				stmt.Tags = append(stmt.Tags, tags...)
				continue
			case kwROW:
				if err := p.skipRowAccessPolicy(); err != nil {
					return nil, err
				}
				continue
			case kwCOMMENT:
				// WITH COMMENT = '...' — the WITH is a no-op anchor; route COMMENT to
				// the open-ended option reader.
				p.advance() // consume WITH
				opt, err := p.parseCopyOption()
				if err != nil {
					return nil, err
				}
				stmt.Options = append(stmt.Options, opt)
				continue
			}
		}
		if p.startsEventTableOption() {
			opt, err := p.parseCopyOption()
			if err != nil {
				return nil, err
			}
			stmt.Options = append(stmt.Options, opt)
			continue
		}
		break
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// startsEventTableOption reports whether cur begins an open-ended event-table
// option. It is the COPY-option predicate minus the keywords that anchor dedicated
// clauses (WITH / TAG / CLUSTER / COPY).
func (p *Parser) startsEventTableOption() bool {
	switch p.cur.Type {
	case kwWITH, kwTAG, kwCLUSTER, kwCOPY:
		return false
	}
	return p.startsCopyOption()
}
