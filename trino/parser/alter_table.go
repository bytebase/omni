package parser

import (
	"github.com/bytebase/omni/trino/ast"
)

// This file is part of the parser-ddl DAG node. It implements Trino's
// ALTER TABLE statement family.
//
// Legacy ANTLR grammar (TrinoParser.g4) the implementation tracks:
//
//	| ALTER_ TABLE_ (IF_ EXISTS_)? from RENAME_ TO_ to                              # renameTable
//	| ALTER_ TABLE_ (IF_ EXISTS_)? tableName ADD_ COLUMN_ (IF_ NOT_ EXISTS_)? columnDefinition # addColumn
//	| ALTER_ TABLE_ (IF_ EXISTS_)? tableName RENAME_ COLUMN_ (IF_ EXISTS_)? from TO_ to # renameColumn
//	| ALTER_ TABLE_ (IF_ EXISTS_)? tableName DROP_ COLUMN_ (IF_ EXISTS_)? column     # dropColumn
//	| ALTER_ TABLE_ (IF_ EXISTS_)? tableName ALTER_ COLUMN_ columnName SET_ DATA_ TYPE_ type # setColumnType
//	| ALTER_ TABLE_ tableName SET_ AUTHORIZATION_ principal                         # setTableAuthorization
//	| ALTER_ TABLE_ tableName SET_ PROPERTIES_ propertyAssignments                  # setTableProperties
//	| ALTER_ TABLE_ tableName EXECUTE_ procedureName (LPAREN_ (callArgument (COMMA_ callArgument)*)? RPAREN_)? (WHERE_ where)? # tableExecute
//
// Adjudicated against the live Trino 481 oracle (which is ahead of the legacy
// grammar). Oracle-confirmed facts baked in:
//
//	D-AT1 (names ≤ 3 parts). The table name and the RENAME-TO target are bounded
//	   to catalog.schema.table. A column name in DROP/RENAME/ALTER COLUMN may be
//	   dotted (e.g. a nested-field path `c.field`), so column-position names are
//	   parsed as an unbounded qualifiedName, matching the legacy grammar.
//	D-AT2 (ADD COLUMN position). Trino 481 accepts a trailing
//	   `FIRST | LAST | AFTER col` on ADD COLUMN (alter-table.html), absent from
//	   the legacy grammar; it is accepted here. The full column constraint set
//	   (DEFAULT / NOT NULL / COMMENT / WITH) of a plain column definition applies.
//	D-AT3 (ALTER COLUMN docs forms). Beyond the legacy `SET DATA TYPE type`,
//	   Trino 481 accepts `ALTER COLUMN col SET DEFAULT expr`,
//	   `ALTER COLUMN col DROP DEFAULT`, and `ALTER COLUMN col DROP NOT NULL`
//	   (alter-table.html). All four are accepted; the oracle pins them.
//	D-AT4 (EXECUTE named args). ALTER TABLE … EXECUTE proc(name => v, …) uses the
//	   same callArgument grammar as CALL; the arg list and the WHERE clause are
//	   both optional.
//	D-AT5 (RENAME TO target is UNBOUNDED). The renameTable `to` target is an
//	   UNBOUNDED qualifiedName: Trino 481 accepts a 4+-part target (the over-bound
//	   failure is the semantic TABLE_NOT_FOUND, not a SYNTAX_ERROR), in contrast
//	   to the source name which is capped at 3. Same for ALTER VIEW / ALTER
//	   MATERIALIZED VIEW RENAME TO. (Oracle-confirmed; divergence ledger DDL4.)
//	D-AT6 (IF EXISTS only on the column/rename sub-forms). The legacy grammar
//	   carries `(IF_ EXISTS_)?` ONLY on renameTable / addColumn / renameColumn /
//	   dropColumn / setColumnType — NOT on setTableAuthorization /
//	   setTableProperties / tableExecute. Trino 481 rejects
//	   `ALTER TABLE IF EXISTS t SET …` and `ALTER TABLE IF EXISTS t EXECUTE …` as
//	   SYNTAX_ERRORs. IF EXISTS is parsed up-front (shared with the column
//	   sub-forms), then the SET/EXECUTE finishers reject it. (Oracle-confirmed.)

// ---------------------------------------------------------------------------
// AST nodes
// ---------------------------------------------------------------------------

// AlterTableKind classifies an ALTER TABLE sub-form.
type AlterTableKind int

const (
	// AlterTableRename is ALTER TABLE [IF EXISTS] name RENAME TO new.
	AlterTableRename AlterTableKind = iota
	// AlterTableAddColumn is ALTER TABLE … ADD COLUMN [IF NOT EXISTS] coldef [pos].
	AlterTableAddColumn
	// AlterTableDropColumn is ALTER TABLE … DROP COLUMN [IF EXISTS] col.
	AlterTableDropColumn
	// AlterTableRenameColumn is ALTER TABLE … RENAME COLUMN [IF EXISTS] from TO to.
	AlterTableRenameColumn
	// AlterTableAlterColumn is ALTER TABLE … ALTER COLUMN col <action>.
	AlterTableAlterColumn
	// AlterTableSetAuthorization is ALTER TABLE name SET AUTHORIZATION principal.
	AlterTableSetAuthorization
	// AlterTableSetProperties is ALTER TABLE name SET PROPERTIES …
	AlterTableSetProperties
	// AlterTableExecute is ALTER TABLE name EXECUTE proc[(args)] [WHERE expr].
	AlterTableExecute
)

// ColumnPosition is the optional ADD COLUMN placement (FIRST | LAST | AFTER c).
type ColumnPosition int

const (
	// ColumnPositionNone means no placement clause.
	ColumnPositionNone ColumnPosition = iota
	// ColumnPositionFirst is FIRST.
	ColumnPositionFirst
	// ColumnPositionLast is LAST.
	ColumnPositionLast
	// ColumnPositionAfter is AFTER column (After holds the reference column).
	ColumnPositionAfter
)

// AlterColumnAction classifies an ALTER COLUMN sub-action.
type AlterColumnAction int

const (
	// AlterColumnSetDataType is SET DATA TYPE new_type.
	AlterColumnSetDataType AlterColumnAction = iota
	// AlterColumnSetDefault is SET DEFAULT expr (D-AT3).
	AlterColumnSetDefault
	// AlterColumnDropDefault is DROP DEFAULT (D-AT3).
	AlterColumnDropDefault
	// AlterColumnDropNotNull is DROP NOT NULL (D-AT3).
	AlterColumnDropNotNull
)

// AlterTableStmt is one ALTER TABLE statement. Only the fields relevant to Kind
// are populated.
type AlterTableStmt struct {
	Kind     AlterTableKind
	IfExists bool // the table-level IF EXISTS (renameTable / addColumn / dropColumn / renameColumn / setColumnType)
	Name     *ast.QualifiedName

	// RENAME TO
	NewName *ast.QualifiedName

	// ADD COLUMN
	ColumnIfNotExists bool
	NewColumn         *ColumnDefinition
	Position          ColumnPosition
	PositionAfter     *ast.Identifier // AFTER reference column (ColumnPositionAfter)

	// DROP / RENAME COLUMN
	ColumnIfExists bool
	Column         *ast.QualifiedName // DROP COLUMN target / RENAME COLUMN from
	RenamedColumn  *ast.Identifier    // RENAME COLUMN … TO new

	// ALTER COLUMN
	ColumnAction  AlterColumnAction
	AlterColumn   *ast.QualifiedName // the column being altered
	NewColumnType *DataType          // SET DATA TYPE
	NewDefault    Expr               // SET DEFAULT expr

	// SET AUTHORIZATION
	Authorization *Principal

	// SET PROPERTIES
	Properties []*Property

	// EXECUTE
	Procedure    *ast.Identifier
	HasArgList   bool // an (arg, …) list was present (even if empty)
	ExecuteArgs  []CallArgument
	ExecuteWhere Expr // optional WHERE expression

	Loc ast.Loc
}

func (n *AlterTableStmt) Tag() ast.NodeTag { return ast.T_Invalid }
func (n *AlterTableStmt) Span() ast.Loc    { return n.Loc }

var _ ast.Node = (*AlterTableStmt)(nil)

// ---------------------------------------------------------------------------
// Parser
// ---------------------------------------------------------------------------

// parseAlterTableStmt parses ALTER TABLE. On entry cur is ALTER (the dispatch
// has peeked TABLE as the second keyword). startOffset is ALTER's byte offset.
func (p *Parser) parseAlterTableStmt(startOffset int) (ast.Node, error) {
	p.advance() // consume ALTER
	p.advance() // consume TABLE

	ifExists, err := p.parseOptionalIfExists()
	if err != nil {
		return nil, err
	}
	// D-AT1: table name ≤ catalog.schema.table.
	name, err := p.parseBoundedQualifiedName(3, "table name")
	if err != nil {
		return nil, err
	}
	stmt := &AlterTableStmt{IfExists: ifExists, Name: name, Loc: ast.Loc{Start: startOffset}}

	switch p.cur.Kind {
	case kwRENAME:
		return p.finishAlterTableRename(stmt)
	case kwADD:
		return p.finishAlterTableAddColumn(stmt)
	case kwDROP:
		return p.finishAlterTableDropColumn(stmt)
	case kwALTER:
		return p.finishAlterTableAlterColumn(stmt)
	case kwSET:
		return p.finishAlterTableSet(stmt)
	case kwEXECUTE:
		return p.finishAlterTableExecute(stmt)
	default:
		return nil, p.syntaxErrorAtCur()
	}
}

// finishAlterTableRename handles `RENAME { TO new | COLUMN [IF EXISTS] from TO to }`.
func (p *Parser) finishAlterTableRename(stmt *AlterTableStmt) (ast.Node, error) {
	p.advance() // consume RENAME
	if _, ok := p.match(kwTO); ok {
		// renameTable. D-AT5: the RENAME TO target is UNBOUNDED — Trino 481
		// accepts a 4+-part target (the over-bound failure is the semantic
		// TABLE_NOT_FOUND, not a SYNTAX_ERROR), unlike the source name which is
		// capped at 3. So the target uses the plain qualifiedName, not the
		// bounded form.
		newName, err := p.parseQualifiedName()
		if err != nil {
			return nil, err
		}
		stmt.Kind = AlterTableRename
		stmt.NewName = newName
		stmt.Loc.End = newName.Loc.End
		return stmt, nil
	}
	// renameColumn: RENAME COLUMN [IF EXISTS] from TO to.
	if _, err := p.expect(kwCOLUMN); err != nil {
		return nil, err
	}
	colIfExists, err := p.parseOptionalIfExists()
	if err != nil {
		return nil, err
	}
	from, err := p.parseQualifiedName() // column path may be dotted (nested field)
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}
	to, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Kind = AlterTableRenameColumn
	stmt.ColumnIfExists = colIfExists
	stmt.Column = from
	stmt.RenamedColumn = to
	stmt.Loc.End = to.Loc.End
	return stmt, nil
}

// finishAlterTableAddColumn handles `ADD COLUMN [IF NOT EXISTS] columnDefinition
// [FIRST | LAST | AFTER col]` (D-AT2).
func (p *Parser) finishAlterTableAddColumn(stmt *AlterTableStmt) (ast.Node, error) {
	p.advance() // consume ADD
	if _, err := p.expect(kwCOLUMN); err != nil {
		return nil, err
	}
	colIfNotExists, err := p.parseOptionalIfNotExists()
	if err != nil {
		return nil, err
	}
	col, err := p.parseColumnDefinition()
	if err != nil {
		return nil, err
	}
	stmt.Kind = AlterTableAddColumn
	stmt.ColumnIfNotExists = colIfNotExists
	stmt.NewColumn = col
	stmt.Loc.End = col.Loc.End

	// D-AT2: optional placement clause.
	switch p.cur.Kind {
	case kwFIRST:
		tok := p.advance()
		stmt.Position = ColumnPositionFirst
		stmt.Loc.End = tok.Loc.End
	case kwLAST:
		tok := p.advance()
		stmt.Position = ColumnPositionLast
		stmt.Loc.End = tok.Loc.End
	case kwAFTER:
		p.advance() // consume AFTER
		after, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Position = ColumnPositionAfter
		stmt.PositionAfter = after
		stmt.Loc.End = after.Loc.End
	}
	return stmt, nil
}

// finishAlterTableDropColumn handles `DROP COLUMN [IF EXISTS] column`.
func (p *Parser) finishAlterTableDropColumn(stmt *AlterTableStmt) (ast.Node, error) {
	p.advance() // consume DROP
	if _, err := p.expect(kwCOLUMN); err != nil {
		return nil, err
	}
	colIfExists, err := p.parseOptionalIfExists()
	if err != nil {
		return nil, err
	}
	col, err := p.parseQualifiedName() // dotted nested-field path allowed
	if err != nil {
		return nil, err
	}
	stmt.Kind = AlterTableDropColumn
	stmt.ColumnIfExists = colIfExists
	stmt.Column = col
	stmt.Loc.End = col.Loc.End
	return stmt, nil
}

// finishAlterTableAlterColumn handles `ALTER COLUMN col { SET DATA TYPE type |
// SET DEFAULT expr | DROP DEFAULT | DROP NOT NULL }` (D-AT3).
func (p *Parser) finishAlterTableAlterColumn(stmt *AlterTableStmt) (ast.Node, error) {
	p.advance() // consume ALTER
	if _, err := p.expect(kwCOLUMN); err != nil {
		return nil, err
	}
	col, err := p.parseQualifiedName() // dotted nested-field path allowed
	if err != nil {
		return nil, err
	}
	stmt.Kind = AlterTableAlterColumn
	stmt.AlterColumn = col

	switch p.cur.Kind {
	case kwSET:
		p.advance() // consume SET
		switch p.cur.Kind {
		case kwDATA:
			p.advance() // consume DATA
			if _, err := p.expect(kwTYPE); err != nil {
				return nil, err
			}
			newType, err := p.parseType()
			if err != nil {
				return nil, err
			}
			stmt.ColumnAction = AlterColumnSetDataType
			stmt.NewColumnType = newType
			stmt.Loc.End = newType.Loc.End
		case kwDEFAULT:
			p.advance() // consume DEFAULT
			def, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			stmt.ColumnAction = AlterColumnSetDefault
			stmt.NewDefault = def
			stmt.Loc.End = def.Span().End
		default:
			return nil, p.syntaxErrorAtCur()
		}
	case kwDROP:
		p.advance() // consume DROP
		switch p.cur.Kind {
		case kwDEFAULT:
			tok := p.advance()
			stmt.ColumnAction = AlterColumnDropDefault
			stmt.Loc.End = tok.Loc.End
		case kwNOT:
			p.advance() // consume NOT
			nullTok, err := p.expect(kwNULL)
			if err != nil {
				return nil, err
			}
			stmt.ColumnAction = AlterColumnDropNotNull
			stmt.Loc.End = nullTok.Loc.End
		default:
			return nil, p.syntaxErrorAtCur()
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}
	return stmt, nil
}

// finishAlterTableSet handles `SET { AUTHORIZATION principal | PROPERTIES … }`.
func (p *Parser) finishAlterTableSet(stmt *AlterTableStmt) (ast.Node, error) {
	// D-AT6 (oracle-pinned): the setTableAuthorization / setTableProperties
	// alternatives do NOT carry `IF EXISTS` (only rename/add/drop/rename/alter
	// column do). Trino 481 rejects `ALTER TABLE IF EXISTS t SET …` as a
	// SYNTAX_ERROR, so reject it here even though IF EXISTS was consumed
	// up-front (it is shared with the column sub-forms).
	if stmt.IfExists {
		return nil, &ParseError{Loc: stmt.Name.Loc, Msg: "IF EXISTS is not allowed on ALTER TABLE … SET"}
	}
	p.advance() // consume SET
	switch p.cur.Kind {
	case kwAUTHORIZATION:
		p.advance() // consume AUTHORIZATION
		princ, err := p.parseAuthorizationPrincipal()
		if err != nil {
			return nil, err
		}
		stmt.Kind = AlterTableSetAuthorization
		stmt.Authorization = princ
		stmt.Loc.End = princ.Loc.End
	case kwPROPERTIES:
		p.advance() // consume PROPERTIES
		props, err := p.parsePropertyAssignments()
		if err != nil {
			return nil, err
		}
		stmt.Kind = AlterTableSetProperties
		stmt.Properties = props
		stmt.Loc.End = p.prev.Loc.End
	default:
		return nil, p.syntaxErrorAtCur()
	}
	return stmt, nil
}

// finishAlterTableExecute handles `EXECUTE proc [(callArgument, …)] [WHERE expr]`
// (D-AT4). The argument list reuses the CALL callArgument grammar (positional or
// `name => expr`). Both the arg list and the WHERE clause are optional.
func (p *Parser) finishAlterTableExecute(stmt *AlterTableStmt) (ast.Node, error) {
	// D-AT6 (oracle-pinned): the tableExecute alternative does NOT carry
	// `IF EXISTS`; Trino 481 rejects `ALTER TABLE IF EXISTS t EXECUTE …` as a
	// SYNTAX_ERROR.
	if stmt.IfExists {
		return nil, &ParseError{Loc: stmt.Name.Loc, Msg: "IF EXISTS is not allowed on ALTER TABLE … EXECUTE"}
	}
	p.advance() // consume EXECUTE
	proc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Kind = AlterTableExecute
	stmt.Procedure = proc
	stmt.Loc.End = proc.Loc.End

	if _, ok := p.match(int('(')); ok {
		stmt.HasArgList = true
		if p.cur.Kind != int(')') {
			first, err := p.parseCallArgument()
			if err != nil {
				return nil, err
			}
			stmt.ExecuteArgs = append(stmt.ExecuteArgs, first)
			for {
				if _, ok := p.match(int(',')); !ok {
					break
				}
				next, err := p.parseCallArgument()
				if err != nil {
					return nil, err
				}
				stmt.ExecuteArgs = append(stmt.ExecuteArgs, next)
			}
		}
		closeTok, err := p.expect(int(')'))
		if err != nil {
			return nil, err
		}
		stmt.Loc.End = closeTok.Loc.End
	}

	if _, ok := p.match(kwWHERE); ok {
		where, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		stmt.ExecuteWhere = where
		stmt.Loc.End = where.Span().End
	}
	return stmt, nil
}
