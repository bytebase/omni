// Package parser - create_table.go implements T-SQL CREATE TABLE statement parsing.
package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/mssql/ast"
)

// Column-level enum value sets for various T-SQL column annotations.
var (
	encryptionTypeValues   = newOptionSet().withIdents("DETERMINISTIC", "RANDOMIZED")
	ledgerColumnKindValues = newOptionSet().withIdents("TRANSACTION_ID", "SEQUENCE_NUMBER")
	ledgerStartEndValues   = newOptionSet().withIdents("START", "END")
	distributionTypeValues = newOptionSet().withIdents("HASH", "ROUND_ROBIN", "REPLICATE")
)

// createTableOptions defines the valid option names for CREATE TABLE ... WITH (...).
// Derived from SqlScriptDOM TSql170.g createTableOption rule.
var createTableOptions = newOptionSet(
	kwMEMORY_OPTIMIZED,
	kwDISTRIBUTION,
	kwPARTITION,
	kwCLUSTERED,
	kwHEAP,
).withIdents(
	"DATA_COMPRESSION",
	"XML_COMPRESSION",
	"FILETABLE_DIRECTORY",
	"FILETABLE_COLLATE_FILENAME",
	"FILETABLE_PRIMARY_KEY_CONSTRAINT_NAME",
	"FILETABLE_STREAMID_UNIQUE_CONSTRAINT_NAME",
	"FILETABLE_FULLPATH_UNIQUE_CONSTRAINT_NAME",
	"DURABILITY",
	"REMOTE_DATA_ARCHIVE",
	"SYSTEM_VERSIONING",
	"LEDGER",
	"LOCATION",
	"DATA_DELETION",
	"DATA_RETENTION",
	"LOCK_ESCALATION",
)

// parseCreateTableStmt parses a CREATE TABLE statement.
//
// BNF: mssql/parser/bnf/create-table-transact-sql.bnf
//
//	CREATE TABLE
//	    { database_name.schema_name.table_name | schema_name.table_name | table_name }
//	    [ AS FileTable ]
//	    ( { <column_definition>
//	        | <computed_column_definition>
//	        | <column_set_definition>
//	        | [ <table_constraint> ] [ ,...n ]
//	        | [ <table_index> ] }
//	          [ ,...n ]
//	          [ PERIOD FOR SYSTEM_TIME ( system_start_time_column_name
//	             , system_end_time_column_name ) ]
//	      )
//	    [ ON { partition_scheme_name ( partition_column_name )
//	           | filegroup
//	           | "default" } ]
//	    [ TEXTIMAGE_ON { filegroup | "default" } ]
//	    [ FILESTREAM_ON { partition_scheme_name | filegroup | "default" } ]
//	    [ WITH ( <table_option> [ ,...n ] ) ]
func (p *Parser) parseCreateTableStmt() (*nodes.CreateTableStmt, error) {
	loc := p.pos()

	// Completion: after CREATE TABLE → identifier (new table name)
	if p.collectMode() {
		p.addRuleCandidate("identifier")
		return nil, errCollecting
	}

	stmt := &nodes.CreateTableStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Table name
	var err error
	stmt.Name, err = p.parseTableRef()
	if err != nil {
		return nil, err
	}
	if stmt.Name == nil {
		return nil, p.unexpectedToken()
	}

	// AS EDGE | AS FILETABLE, with no list (createTableStatement,
	// TSql160.g:26543: As ( 'FILETABLE' | 'EDGE' )). AS NODE comes only after
	// a list (graphAsNodeOrEdge), and neither takes one after it: SQL Server
	// refuses CREATE TABLE a AS NODE(a INT) and CREATE TABLE a AS EDGE(a INT)
	// (Msg 102, measured). Their statements still take the table options
	// that follow a list, as
	// CREATE TABLE ft AS FILETABLE WITH (FILETABLE_DIRECTORY = 'docs') does.
	// Note: AS CLONE OF is handled at the dispatch level (parser.go)
	if p.cur.Type == kwAS {
		next := p.peekNext()
		if next.Type == kwEDGE || next.Type == kwFILETABLE {
			p.advance() // AS
			p.advance() // EDGE or FILETABLE
			stmt.IsEdge = next.Type == kwEDGE
			stmt.IsFileTable = next.Type == kwFILETABLE
			return p.parseCreateTableOptions(stmt)
		}
	}

	// Column and constraint definitions. SQL Server requires them, except on
	// AS EDGE and AS FILETABLE: CREATE TABLE a and CREATE TABLE a AS NODE are
	// Msg 102 (measured).
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	// Completion: inside CREATE TABLE column list → constraint keywords / columnref
	if p.collectMode() {
		p.addRuleCandidate("identifier")
		p.addTokenCandidate(kwCONSTRAINT)
		p.addTokenCandidate(kwPRIMARY)
		p.addTokenCandidate(kwUNIQUE)
		p.addTokenCandidate(kwFOREIGN)
		p.addTokenCandidate(kwCHECK)
		p.addTokenCandidate(kwINDEX)
		return nil, errCollecting
	}

	var cols []nodes.Node
	var constraints []nodes.Node

	// CREATE TABLE column/constraint list is the one site where SqlScriptDOM
	// tolerates a trailing comma (see scriptdom harness verification).
	items, listErr := p.parseCommaList(')', commaListAllowTrail, func() (nodes.Node, error) {
		// Completion: at start of each column/constraint slot.
		if p.collectMode() {
			p.addRuleCandidate("identifier")
			p.addTokenCandidate(kwCONSTRAINT)
			p.addTokenCandidate(kwPRIMARY)
			p.addTokenCandidate(kwUNIQUE)
			p.addTokenCandidate(kwFOREIGN)
			p.addTokenCandidate(kwCHECK)
			p.addTokenCandidate(kwINDEX)
			return nil, errCollecting
		}
		// PERIOD FOR SYSTEM_TIME ( start , end ), every part required
		// (tablePeriodDefinition, TSql160.g:28778) — stored on stmt
		// directly, not collected here.
		if p.cur.Type == kwPERIOD && p.peekNext().Type == kwFOR {
			p.advance() // PERIOD
			p.advance() // FOR
			if _, err := p.expect(kwSYSTEM_TIME); err != nil {
				return nil, err
			}
			if _, err := p.expect('('); err != nil {
				return nil, err
			}
			startCol, ok := p.parseIdentifier()
			if !ok {
				return nil, p.syntaxErrorAtCur()
			}
			if _, err := p.expect(','); err != nil {
				return nil, err
			}
			endCol, ok := p.parseIdentifier()
			if !ok {
				return nil, p.syntaxErrorAtCur()
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
			stmt.PeriodStartCol = startCol
			stmt.PeriodEndCol = endCol
			// Return a sentinel so parseCommaList accepts this as "one item".
			return &nodes.String{Str: "__period__"}, nil
		}
		if p.cur.Type == kwINDEX {
			idx, err := p.parseInlineTableIndex()
			if err != nil {
				return nil, err
			}
			if idx != nil {
				if stmt.Indexes == nil {
					stmt.Indexes = &nodes.List{}
				}
				stmt.Indexes.Items = append(stmt.Indexes.Items, idx)
			}
			return idx, nil
		}
		if p.cur.Type == kwCONSTRAINT || p.cur.Type == kwPRIMARY ||
			p.cur.Type == kwUNIQUE || p.cur.Type == kwCHECK ||
			p.cur.Type == kwFOREIGN {
			constraint, err := p.parseTableConstraint()
			if err != nil {
				return nil, err
			}
			if constraint != nil {
				constraints = append(constraints, constraint)
			}
			return constraint, nil
		}
		col, err := p.parseTypedColumnDef()
		if err != nil {
			return nil, err
		}
		if col != nil {
			cols = append(cols, col)
		}
		return col, nil
	})
	if listErr != nil {
		return nil, listErr
	}
	rparen, err := p.expect(')')
	if err != nil {
		return nil, err
	}

	if len(cols) > 0 {
		stmt.Columns = &nodes.List{Items: cols}
	}
	if len(constraints) > 0 {
		stmt.Constraints = &nodes.List{Items: constraints}
	}

	if _, err := p.parseCreateTableOptions(stmt); err != nil {
		return nil, err
	}
	if !tableDefHasColumn(items, stmt.IsEdge) {
		return nil, p.newParseError(rparen.Loc, "syntax error at or near \")\"")
	}
	return stmt, nil
}

// tableDefHasColumn reports whether a table definition's elements hold what
// SQL Server's parser requires of them, which tableDefinitionCreateTable
// (TSql160.g:26628) and declareTableBodyMain (TSql160.g:16716) do not say: a column with a type. A computed column, a
// PERIOD, a constraint or an index alone is refused at the closing
// parenthesis (Msg 102, measured): CREATE TABLE a(UNIQUE(a)) and
// CREATE TABLE a(a AS 1) are, and so are a table type's, a table variable's
// and a function's return table. An AS EDGE table, which has columns of its
// own, needs a typed column, a constraint or an index: SQL Server takes
// CREATE TABLE a(UNIQUE(a)) AS EDGE and refuses CREATE TABLE a(a AS 1) AS EDGE.
func tableDefHasColumn(items []nodes.Node, edge bool) bool {
	for _, item := range items {
		switch n := item.(type) {
		case *nodes.ColumnDef:
			if n.DataType != nil {
				return true
			}
		case *nodes.String: // the PERIOD sentinel
		case nil:
		default:
			if edge {
				return true
			}
		}
	}
	return false
}

// parseCreateTableOptions parses what follows a CREATE TABLE body: ON
// <filegroup>, TEXTIMAGE_ON, FILESTREAM_ON and WITH ( table_options ), and
// ends the statement.
func (p *Parser) parseCreateTableOptions(stmt *nodes.CreateTableStmt) (*nodes.CreateTableStmt, error) {
	// Completion: after closing paren → ON/WITH/;
	if p.collectMode() {
		p.addTokenCandidate(kwON)
		p.addTokenCandidate(kwWITH)
		p.addTokenCandidate(';')
		return nil, errCollecting
	}

	// ON { partition_scheme_name ( partition_column_name ) | filegroup | "default" }
	if p.cur.Type == kwON {
		p.advance()
		if p.isIdentLike() || p.cur.Type == kwPRIMARY {
			name := p.cur.Str
			p.advance()
			if p.cur.Type == '(' {
				// partition_scheme_name ( partition_column_name )
				p.advance()
				partCol, _ := p.parseIdentifier()
				p.expect(')')
				stmt.OnFilegroup = name + "(" + partCol + ")"
			} else {
				stmt.OnFilegroup = name
			}
		}
	}

	// TEXTIMAGE_ON filegroup
	if p.cur.Type == kwTEXTIMAGE_ON {
		p.advance()
		if p.isIdentLike() || p.cur.Type == kwPRIMARY {
			stmt.TextImageOn = p.cur.Str
			p.advance()
		}
	}

	// FILESTREAM_ON filegroup
	if p.cur.Type == kwFILESTREAM_ON {
		p.advance()
		if p.isIdentLike() || p.cur.Type == kwPRIMARY {
			stmt.FilestreamOn = p.cur.Str
			p.advance()
		}
	}

	// WITH ( table_options )
	if p.cur.Type == kwWITH {
		p.advance()
		tableOpts, err := p.parseTableOptions()
		if err != nil {
			return nil, err
		}
		stmt.TableOptions = tableOpts
	}

	// AS NODE | AS EDGE (graph tables) -- can also appear after closing paren
	if !stmt.IsNode && !stmt.IsEdge && !stmt.IsFileTable && p.cur.Type == kwAS {
		next := p.peekNext()
		if next.Type == kwNODE {
			p.advance() // AS
			p.advance() // NODE
			stmt.IsNode = true
		} else if next.Type == kwEDGE {
			p.advance() // AS
			p.advance() // EDGE
			stmt.IsEdge = true
		}
	}

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseTableOptions parses ( option = value [, ...] ) for CREATE TABLE WITH clause.
//
//	<table_option> ::=
//	    DATA_COMPRESSION = { NONE | ROW | PAGE }
//	  | XML_COMPRESSION = { ON | OFF }
//	  | MEMORY_OPTIMIZED = ON
//	  | DURABILITY = { SCHEMA_ONLY | SCHEMA_AND_DATA }
//	  | SYSTEM_VERSIONING = ON [ ( HISTORY_TABLE = schema.table
//	      [, DATA_CONSISTENCY_CHECK = { ON | OFF } ]
//	      [, HISTORY_RETENTION_PERIOD = { INFINITE | number { DAY | DAYS | WEEK | WEEKS
//	                                    | MONTH | MONTHS | YEAR | YEARS } } ] ) ]
//	  | REMOTE_DATA_ARCHIVE = { ON [ ( <table_stretch_options> ) ] | OFF ( MIGRATION_STATE = PAUSED ) }
//	  | DATA_DELETION = ON { ( FILTER_COLUMN = column_name,
//	      RETENTION_PERIOD = { INFINITE | number { DAY | DAYS | ... } } ) }
//	  | LEDGER = ON [ ( <ledger_option> [ ,...n ] ) ] | OFF
//	  | FILETABLE_DIRECTORY = <directory_name>
//	  | FILETABLE_COLLATE_FILENAME = { <collation_name> | database_default }
//	  | FILETABLE_PRIMARY_KEY_CONSTRAINT_NAME = <constraint_name>
//	  | FILETABLE_STREAMID_UNIQUE_CONSTRAINT_NAME = <constraint_name>
//	  | FILETABLE_FULLPATH_UNIQUE_CONSTRAINT_NAME = <constraint_name>
func (p *Parser) parseTableOptions() (*nodes.List, error) {
	if p.cur.Type != '(' {
		return nil, nil
	}
	p.advance()

	items, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
		opt, err := p.parseOneTableOption()
		if err != nil {
			return nil, err
		}
		if opt == nil {
			return nil, p.unexpectedToken()
		}
		return opt, nil
	})
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, nil
	}
	return &nodes.List{Items: items}, nil
}

// parseOneTableOption parses a single NAME = VALUE table option.
func (p *Parser) parseOneTableOption() (*nodes.TableOption, error) {
	loc := p.pos()
	if !p.isValidOption(createTableOptions) {
		return nil, nil
	}
	name := strings.ToUpper(p.cur.Str)
	p.advance()

	opt := &nodes.TableOption{
		Name: name,
		Loc:  nodes.Loc{Start: loc, End: -1},
	}

	// Every table option takes '=' and a value, and SQL Server refuses both
	// `WITH (DATA_COMPRESSION)` and `WITH (FILETABLE_DIRECTORY = )` with
	// Msg 102 (measured).
	if p.cur.Type != '=' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Type == ')' || p.cur.Type == ',' || p.cur.Type == tokEOF {
		return nil, p.syntaxErrorAtCur()
	}

	// SYSTEM_VERSIONING = { ON | OFF } [ ( ... ) ]
	if name == "SYSTEM_VERSIONING" {
		if p.isValidOption(onOffOptionValues) {
			opt.Value = strings.ToUpper(p.cur.Str)
			p.advance()
		}
		// Optional sub-options
		if opt.Value == "ON" && p.cur.Type == '(' {
			p.advance()
			for p.cur.Type != ')' && p.cur.Type != tokEOF {
				if p.isIdentLike() {
					subName := strings.ToUpper(p.cur.Str)
					p.advance()
					if p.cur.Type == '=' {
						p.advance()
					}
					switch subName {
					case "HISTORY_TABLE":
						// schema.table - collect dotted name
						var parts []string
						if p.isIdentLike() {
							parts = append(parts, p.cur.Str)
							p.advance()
						}
						for p.cur.Type == '.' {
							p.advance()
							if p.isIdentLike() {
								parts = append(parts, p.cur.Str)
								p.advance()
							}
						}
						opt.HistoryTable = strings.Join(parts, ".")
					case "DATA_CONSISTENCY_CHECK":
						// DATA_CONSISTENCY_CHECK = { ON | OFF }
						if p.isValidOption(onOffOptionValues) {
							opt.DataConsistencyCheck = strings.ToUpper(p.cur.Str)
							p.advance()
						}
					case "HISTORY_RETENTION_PERIOD":
						// Collect value tokens until comma or paren
						var valParts []string
						for p.cur.Type != ',' && p.cur.Type != ')' && p.cur.Type != tokEOF {
							valParts = append(valParts, p.cur.Str)
							p.advance()
						}
						opt.HistoryRetentionPeriod = strings.Join(valParts, " ")
					default:
						// Consume unknown sub-option value structurally (ON/OFF/identifier/number)
						if p.cur.Type == kwON {
							p.advance()
						} else if p.cur.Type == kwOFF {
							p.advance()
						} else if p.isIdentLike() {
							p.advance()
						} else if p.cur.Type == tokICONST {
							p.advance()
						}
					}
				}
				if _, ok := p.match(','); !ok {
					break
				}
			}
			p.expect(')')
		}
	} else {
		// Simple NAME = VALUE. FILETABLE_DIRECTORY is the one that takes a
		// string (<directory_name>), which the option keeps as written,
		// quotes and all; the constraint-name options beside it take an
		// identifier, and the server refuses `FILETABLE_PRIMARY_KEY_
		// CONSTRAINT_NAME = 'pk'` (measured).
		if name == "FILETABLE_DIRECTORY" && (p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST) {
			opt.Value = p.tokenText(p.cur)
			p.advance()
		} else if p.isIdentLike() {
			opt.Value = strings.ToUpper(p.cur.Str)
			p.advance()
		} else if p.cur.Type == kwON {
			opt.Value = "ON"
			p.advance()
		} else if p.cur.Type == kwOFF {
			opt.Value = "OFF"
			p.advance()
		}
	}

	opt.Loc.End = p.prevEnd()
	return opt, nil
}

// parseColumnDef parses a column definition, computed column, or column set.
//
// BNF: mssql/parser/bnf/create-table-transact-sql.bnf
//
//	<column_definition> ::=
//	column_name <data_type>
//	    [ FILESTREAM ]
//	    [ COLLATE collation_name ]
//	    [ SPARSE ]
//	    [ MASKED WITH ( FUNCTION = 'mask_function' ) ]
//	    [ [ CONSTRAINT constraint_name ] DEFAULT constant_expression ]
//	    [ IDENTITY [ ( seed , increment ) ] ]
//	    [ NOT FOR REPLICATION ]
//	    [ GENERATED ALWAYS AS { ROW | TRANSACTION_ID | SEQUENCE_NUMBER }
//	      { START | END } [ HIDDEN ] ]
//	    [ [ CONSTRAINT constraint_name ] { NULL | NOT NULL } ]
//	    [ ROWGUIDCOL ]
//	    [ ENCRYPTED WITH
//	        ( COLUMN_ENCRYPTION_KEY = key_name ,
//	          ENCRYPTION_TYPE = { DETERMINISTIC | RANDOMIZED } ,
//	          ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256'
//	        ) ]
//	    [ <column_constraint> [ ,...n ] ]
//	    [ <column_index> ]
//
//	<computed_column_definition> ::=
//	column_name AS computed_column_expression
//	[ PERSISTED [ NOT NULL ] ]
//	[ <column_constraint> ]
//
//	<column_set_definition> ::=
//	column_set_name XML COLUMN_SET FOR ALL_SPARSE_COLUMNS
func (p *Parser) parseColumnDef() (*nodes.ColumnDef, error) {
	loc := p.pos()

	name, ok := p.parseIdentifier()
	if !ok {
		return nil, nil
	}

	col := &nodes.ColumnDef{
		Name: name,
		Loc:  nodes.Loc{Start: loc, End: -1},
	}

	// Completion: after column name → type_name
	if p.collectMode() {
		p.addRuleCandidate("type_name")
		return nil, errCollecting
	}

	// Check for computed column: name AS expr [ PERSISTED [ NOT NULL ] ]
	if p.cur.Type == kwAS {
		compLoc := p.pos()
		p.advance()
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		persisted := false
		notNull := false
		if p.cur.Type == kwPERSISTED {
			persisted = true
			p.advance()
			// [ NOT NULL ]
			if p.cur.Type == kwNOT && p.peekNext().Type == kwNULL {
				p.advance() // NOT
				p.advance() // NULL
				notNull = true
			}
		}
		col.Computed = &nodes.ComputedColumnDef{
			Expr:      expr,
			Persisted: persisted,
			NotNull:   notNull,
			Loc:       nodes.Loc{Start: compLoc, End: p.prevEnd()},
		}
		// Computed columns can also have column constraints (PK, UNIQUE,
		// CHECK), but no foreign key: SQL Server refuses a AS 1 REFERENCES b
		// and a AS 1 CONSTRAINT c FOREIGN KEY REFERENCES b, PERSISTED or not
		// (measured, compat.tsv).
		for p.cur.Type == kwCONSTRAINT || p.cur.Type == kwPRIMARY || p.cur.Type == kwUNIQUE ||
			p.cur.Type == kwCHECK {
			var constraint *nodes.ConstraintDef
			var cErr error
			if p.cur.Type == kwCONSTRAINT {
				if kind := p.peekSecond(); kind.Type == kwFOREIGN || kind.Type == kwREFERENCES {
					return nil, p.newParseError(kind.Loc, "a computed column takes no foreign key")
				}
				constraint, cErr = p.parseColumnConstraint()
			} else {
				constraint, cErr = p.parseInlineConstraint("")
			}
			if cErr != nil {
				return nil, cErr
			}
			if constraint != nil {
				if col.Constraints == nil {
					col.Constraints = &nodes.List{}
				}
				col.Constraints.Items = append(col.Constraints.Items, constraint)
			}
		}
		col.Loc.End = p.prevEnd()
		return col, nil
	}

	// Data type
	var dtErr error
	col.DataType, dtErr = p.parseDataType()
	if dtErr != nil {
		return nil, dtErr
	}

	// Completion: after column data type → column option keywords
	if p.collectMode() {
		p.addTokenCandidate(kwNULL)
		p.addTokenCandidate(kwNOT)
		p.addTokenCandidate(kwDEFAULT)
		p.addTokenCandidate(kwIDENTITY)
		p.addTokenCandidate(kwPRIMARY)
		p.addTokenCandidate(kwUNIQUE)
		p.addTokenCandidate(kwCHECK)
		p.addTokenCandidate(kwREFERENCES)
		p.addTokenCandidate(kwCONSTRAINT)
		p.addTokenCandidate(kwCOLLATE)
		p.addTokenCandidate(kwROWGUIDCOL)
		return nil, errCollecting
	}

	// column_set_definition: column_set_name XML COLUMN_SET FOR ALL_SPARSE_COLUMNS.
	// SQL Server checks neither the type name (it takes `a a COLUMN_SET FOR
	// ALL_SPARSE_COLUMNS`) nor the words before FOR, but it requires FOR ALL_SPARSE_COLUMNS:
	// `a a COLUMN_SET` is Msg 102, and SPARSE COLUMN_SET is too.
	if p.cur.Type == kwCOLUMN_SET {
		p.advance() // COLUMN_SET
		if _, err := p.expect(kwFOR); err != nil {
			return nil, err
		}
		if _, err := p.expect(kwALL_SPARSE_COLUMNS); err != nil {
			return nil, err
		}
		col.IsColumnSet = true
		col.Loc.End = p.prevEnd()
		return col, nil
	}

	// COLLATE directly follows the data type: SQL Server rejects it after
	// any other column option (NOT NULL COLLATE x, PRIMARY KEY COLLATE x)
	// and a second COLLATE, with Msg 156.
	if p.cur.Type == kwCOLLATE {
		p.advance()
		collation, err := p.parseCollationName()
		if err != nil {
			return nil, err
		}
		col.Collation = collation
	}

	// The options SQL Server takes once each, in this order, before any
	// constraint: GENERATED ALWAYS ... [ HIDDEN ], or else SPARSE and
	// FILESTREAM in either order, then MASKED WITH or ENCRYPTED WITH (not
	// after FILESTREAM). After a constraint or one another out of this order,
	// SQL Server rejects them, as NOT NULL SPARSE and MASKED WITH (...)
	// ENCRYPTED WITH (...) are.
	if p.cur.Type == kwGENERATED {
		genAlways, genErr := p.parseGeneratedAlways()
		if genErr != nil {
			return nil, genErr
		}
		col.GeneratedAlways = genAlways
		if p.cur.Type == kwHIDDEN {
			col.Hidden = true
			p.advance()
		}
	} else {
		if err := p.parseColumnStorage(col); err != nil {
			return nil, err
		}
		if p.cur.Type == kwMASKED {
			p.advance() // MASKED
			fn, err := p.parseMaskedWithFunction()
			if err != nil {
				return nil, err
			}
			col.MaskFunction = fn
		} else if p.cur.Type == kwENCRYPTED && !col.Filestream {
			encWith, encErr := p.parseEncryptedWith()
			if encErr != nil {
				return nil, encErr
			}
			col.EncryptedWith = encWith
		}
	}

	// The constraints, with ROWGUIDCOL and IDENTITY among them, in any order
	// and any number of times: SQL Server's parser takes NULL NOT NULL and
	// two IDENTITY clauses, which fail later. NOT FOR REPLICATION belongs to
	// the IDENTITY or REFERENCES before it, and SQL Server rejects it
	// anywhere else.
	for {
		switch p.cur.Type {
		case kwROWGUIDCOL:
			col.Rowguidcol = true
			p.advance()
			continue
		case kwIDENTITY:
			identSpec, identErr := p.parseIdentitySpec()
			if identErr != nil {
				return nil, identErr
			}
			col.Identity = identSpec
			if p.cur.Type == kwNOT && p.peekNext().Type == kwFOR {
				if err := p.parseNotForReplication(); err != nil {
					return nil, err
				}
				col.NotForReplication = true
			}
			continue
		case kwNOT:
			if p.peekNext().Type != kwNULL {
				break
			}
			nullLoc := p.pos()
			p.advance() // NOT
			p.advance() // NULL
			col.Nullable = &nodes.NullableSpec{NotNull: true, Loc: nodes.Loc{Start: nullLoc, End: p.prevEnd()}}
			continue
		case kwNULL:
			nullLoc := p.pos()
			p.advance()
			col.Nullable = &nodes.NullableSpec{NotNull: false, Loc: nodes.Loc{Start: nullLoc, End: p.prevEnd()}}
			continue
		case kwDEFAULT:
			p.advance()
			// Completion: after DEFAULT → expression
			if p.collectMode() {
				p.addRuleCandidate("expression")
				return nil, errCollecting
			}
			defExpr, defErr := p.parseExpr()
			if defErr != nil {
				return nil, defErr
			}
			if defExpr == nil {
				return nil, p.unexpectedToken()
			}
			col.DefaultExpr = defExpr
			continue
		case kwINDEX:
			idx, err := p.parseInlineColumnIndex()
			if err != nil {
				return nil, err
			}
			// A column may declare several: INDEX ix_a INDEX ix_b parses
			// (measured, compat.tsv).
			if col.Indexes == nil {
				col.Indexes = &nodes.List{}
			}
			col.Indexes.Items = append(col.Indexes.Items, idx)
			continue
		case kwCONSTRAINT, kwPRIMARY, kwUNIQUE, kwCHECK, kwFOREIGN, kwREFERENCES:
			// CONSTRAINT name, then the constraint, or a constraint without
			// a name (PRIMARY KEY, UNIQUE, CHECK, inline REFERENCES).
			var constraint *nodes.ConstraintDef
			var cErr error
			if p.cur.Type == kwCONSTRAINT {
				constraint, cErr = p.parseColumnConstraint()
			} else {
				constraint, cErr = p.parseInlineConstraint("")
			}
			if cErr != nil {
				return nil, cErr
			}
			if constraint != nil {
				if col.Constraints == nil {
					col.Constraints = &nodes.List{}
				}
				col.Constraints.Items = append(col.Constraints.Items, constraint)
			}
			continue
		}
		break
	}

	col.Loc.End = p.prevEnd()
	return col, nil
}

// parseTypedColumnDef parses a column definition that SQL Server requires
// to name its type unless it is computed (columnDefinition, TSql160.g:28920):
// it refuses `a NULL` or `a PRIMARY KEY`, where a reserved word follows the
// name, with Msg 173. It returns (nil, nil) when no column name starts here.
func (p *Parser) parseTypedColumnDef() (*nodes.ColumnDef, error) {
	col, err := p.parseColumnDef()
	if err != nil || col == nil {
		return col, err
	}
	if col.DataType == nil && col.Computed == nil {
		// A timestamp column may leave out its name: SQL Server names the
		// column timestamp and types it so (rowversion docs).
		if strings.ToUpper(col.Name) != "TIMESTAMP" {
			return nil, p.syntaxErrorAtCur()
		}
		col.DataType = &nodes.DataType{Name: "timestamp", Loc: col.Loc}
	}
	return col, nil
}

// parseColumnStorage parses SPARSE and FILESTREAM, in either order and once
// each, of a column definition. SQL Server's parser takes FILESTREAM only on
// a varbinary(max) column (Msg 102 otherwise).
func (p *Parser) parseColumnStorage(col *nodes.ColumnDef) error {
	for {
		switch {
		case p.cur.Type == kwSPARSE && !col.Sparse:
			col.Sparse = true
			p.advance()
		case p.cur.Type == kwFILESTREAM && !col.Filestream:
			if col.DataType == nil || strings.ToUpper(col.DataType.Name) != "VARBINARY" || !col.DataType.MaxLength {
				return p.syntaxErrorAtCur()
			}
			col.Filestream = true
			p.advance()
		default:
			return nil
		}
	}
}

// parseMaskedWithFunction parses WITH ( FUNCTION = 'mask_function' ) after
// MASKED (maskedClause, TSql160.g) and returns the function's text.
func (p *Parser) parseMaskedWithFunction() (string, error) {
	if _, err := p.expect(kwWITH); err != nil {
		return "", err
	}
	if _, err := p.expect('('); err != nil {
		return "", err
	}
	if _, err := p.expect(kwFUNCTION); err != nil {
		return "", err
	}
	if _, err := p.expect('='); err != nil {
		return "", err
	}
	if p.cur.Type != tokSCONST && p.cur.Type != tokNSCONST {
		return "", p.syntaxErrorAtCur()
	}
	fn := p.cur.Str
	p.advance()
	if _, err := p.expect(')'); err != nil {
		return "", err
	}
	return fn, nil
}

// parseEncryptedWith parses ENCRYPTED WITH (...).
//
//	ENCRYPTED WITH
//	    ( COLUMN_ENCRYPTION_KEY = key_name ,
//	      ENCRYPTION_TYPE = { DETERMINISTIC | RANDOMIZED } ,
//	      ALGORITHM = 'AEAD_AES_256_CBC_HMAC_SHA_256'
//	    )
func (p *Parser) parseEncryptedWith() (*nodes.EncryptedWithSpec, error) {
	loc := p.pos()
	p.advance() // ENCRYPTED

	spec := &nodes.EncryptedWithSpec{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// SQL Server wants all three options, once each, in any order
	// (columnEncryptionClause, TSql160.g; Msg 102 near ENCRYPTED when one
	// is missing).
	if _, err := p.expect(kwWITH); err != nil {
		return nil, err
	}
	if _, err := p.expect('('); err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	for {
		name, err := p.tableOptionName(seen, "COLUMN_ENCRYPTION_KEY", "ENCRYPTION_TYPE", "ALGORITHM")
		if err != nil {
			return nil, err
		}
		switch name {
		case "COLUMN_ENCRYPTION_KEY":
			key, ok := p.parseIdentifier()
			if !ok {
				return nil, p.syntaxErrorAtCur()
			}
			spec.ColumnEncryptionKey = key
		case "ENCRYPTION_TYPE":
			if !p.isValidOption(encryptionTypeValues) {
				return nil, p.syntaxErrorAtCur()
			}
			spec.EncryptionType = strings.ToUpper(p.cur.Str)
			p.advance()
		case "ALGORITHM":
			if p.cur.Type != tokSCONST && p.cur.Type != tokNSCONST {
				return nil, p.syntaxErrorAtCur()
			}
			spec.Algorithm = p.cur.Str
			p.advance()
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}
	if len(seen) != 3 {
		return nil, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	spec.Loc.End = p.prevEnd()
	return spec, nil
}

// parseGeneratedAlways parses GENERATED ALWAYS AS { ROW | TRANSACTION_ID | SEQUENCE_NUMBER } { START | END }.
//
//	GENERATED ALWAYS AS { ROW | TRANSACTION_ID | SEQUENCE_NUMBER } { START | END } [ HIDDEN ]
func (p *Parser) parseGeneratedAlways() (*nodes.GeneratedAlwaysSpec, error) {
	loc := p.pos()
	p.advance() // GENERATED

	spec := &nodes.GeneratedAlwaysSpec{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// ALWAYS AS { ROW | TRANSACTION_ID | SEQUENCE_NUMBER } { START | END }
	// (generatedAlwaysClause, TSql160.g:29325), every word required.
	if _, err := p.expect(kwALWAYS); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwAS); err != nil {
		return nil, err
	}
	switch {
	case p.cur.Type == kwROW:
		spec.Kind = "ROW"
	case p.isValidOption(ledgerColumnKindValues):
		spec.Kind = strings.ToUpper(p.cur.Str)
	default:
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	switch {
	case p.cur.Type == kwEND:
		spec.StartEnd = "END"
	case p.isValidOption(ledgerStartEndValues):
		spec.StartEnd = strings.ToUpper(p.cur.Str)
	default:
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()

	spec.Loc.End = p.prevEnd()
	return spec, nil
}

// parseIdentitySpec parses IDENTITY(seed, increment).
func (p *Parser) parseIdentitySpec() (*nodes.IdentitySpec, error) {
	loc := p.pos()
	p.advance() // consume IDENTITY

	spec := &nodes.IdentitySpec{
		Seed:      1,
		Increment: 1,
		Loc:       nodes.Loc{Start: loc, End: -1},
	}

	if p.cur.Type == '(' {
		p.advance()
		// Completion: inside IDENTITY( → numeric
		if p.collectMode() {
			p.addRuleCandidate("numeric")
			return nil, errCollecting
		}
		if p.cur.Type == tokICONST {
			spec.Seed = p.cur.Ival
			p.advance()
		}
		if _, ok := p.match(','); ok {
			if p.cur.Type == tokICONST {
				spec.Increment = p.cur.Ival
				p.advance()
			}
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	}

	spec.Loc.End = p.prevEnd()
	return spec, nil
}

// parseColumnConstraint parses CONSTRAINT name followed by constraint type.
func (p *Parser) parseColumnConstraint() (*nodes.ConstraintDef, error) {
	p.advance() // consume CONSTRAINT
	name, ok := p.parseIdentifier()
	if !ok {
		return nil, p.unexpectedToken()
	}
	return p.parseInlineConstraint(name)
}

// parseInlineConstraint parses a constraint type (PRIMARY KEY, UNIQUE, CHECK, DEFAULT, REFERENCES).
//
// BNF: mssql/parser/bnf/create-table-transact-sql.bnf
//
//	<column_constraint> ::=
//	[ CONSTRAINT constraint_name ]
//	{
//	   { PRIMARY KEY | UNIQUE }
//	        [ CLUSTERED | NONCLUSTERED ]
//	        [ WITH FILLFACTOR = fillfactor | WITH ( <index_option> [ ,...n ] ) ]
//	        [ ON { partition_scheme_name ( partition_column_name )
//	            | filegroup | "default" } ]
//	  | [ FOREIGN KEY ]
//	        REFERENCES [ schema_name. ] referenced_table_name [ ( ref_column ) ]
//	        [ ON DELETE { NO ACTION | CASCADE | SET NULL | SET DEFAULT } ]
//	        [ ON UPDATE { NO ACTION | CASCADE | SET NULL | SET DEFAULT } ]
//	        [ NOT FOR REPLICATION ]
//	  | CHECK [ NOT FOR REPLICATION ] ( logical_expression )
//	}
func (p *Parser) parseInlineConstraint(name string) (*nodes.ConstraintDef, error) {
	loc := p.pos()
	cd := &nodes.ConstraintDef{
		Name: name,
		Loc:  nodes.Loc{Start: loc, End: -1},
	}

	switch p.cur.Type {
	case kwPRIMARY:
		p.advance() // PRIMARY
		if _, err := p.expect(kwKEY); err != nil {
			return nil, err
		}
		cd.Type = nodes.ConstraintPrimaryKey
		p.parseClusteredOption(cd)
		if err := p.parseColumnConstraintColumnList(cd); err != nil {
			return nil, err
		}
		p.parseConstraintWithOptions(cd)
		if err := p.parseConstraintOnFilegroup(cd); err != nil {
			return nil, err
		}
	case kwUNIQUE:
		p.advance()
		cd.Type = nodes.ConstraintUnique
		p.parseClusteredOption(cd)
		if err := p.parseColumnConstraintColumnList(cd); err != nil {
			return nil, err
		}
		p.parseConstraintWithOptions(cd)
		if err := p.parseConstraintOnFilegroup(cd); err != nil {
			return nil, err
		}
	case kwCHECK:
		p.advance()
		cd.Type = nodes.ConstraintCheck
		// [ NOT FOR REPLICATION ]
		if p.cur.Type == kwNOT && p.peekNext().Type == kwFOR {
			if err := p.parseNotForReplication(); err != nil {
				return nil, err
			}
			cd.NotForReplication = true
		}
		if _, err := p.expect('('); err == nil {
			// A CHECK constraint's expression is a search condition.
			p.enterSearchCondition()
			var exprErr error
			cd.Expr, exprErr = p.parseExpr()
			p.leaveSearchCondition()
			if exprErr != nil {
				return nil, exprErr
			}
			if cd.Expr == nil {
				return nil, p.unexpectedToken()
			}
			if _, err := p.expect(')'); err != nil {
				return nil, err
			}
		}
	case kwDEFAULT:
		p.advance()
		cd.Type = nodes.ConstraintDefault
		var exprErr error
		cd.Expr, exprErr = p.parseExpr()
		if exprErr != nil {
			return nil, exprErr
		}
		if cd.Expr == nil {
			return nil, p.unexpectedToken()
		}
	case kwFOREIGN, kwREFERENCES:
		// foreignKeyColumnConstraint (TSql160.g:29998): [ FOREIGN KEY
		// [ ( column [ ,...n ] ) ] ] REFERENCES ..., in CREATE and ALTER TABLE
		// only (ThrowSyntaxErrorIfNotCreateAlterTable): a table variable's,
		// a table type's and a function's table refuse it (measured,
		// compat.tsv).
		if p.tableTypeBody {
			return nil, p.syntaxErrorAtCur()
		}
		cd.Type = nodes.ConstraintForeignKey
		if p.cur.Type == kwFOREIGN {
			p.advance() // FOREIGN
			if _, err := p.expect(kwKEY); err != nil {
				return nil, err
			}
			if p.cur.Type == '(' {
				var err error
				cd.Columns, err = p.parseParenIdentList()
				if err != nil {
					return nil, err
				}
			}
			if p.cur.Type != kwREFERENCES {
				return nil, p.syntaxErrorAtCur()
			}
		}
		p.advance() // REFERENCES
		// Completion: after REFERENCES → table_ref
		if p.collectMode() {
			p.addRuleCandidate("table_ref")
			return nil, errCollecting
		}
		var refErr error
		cd.RefTable, refErr = p.parseTableRef()
		if refErr != nil {
			return nil, refErr
		}
		if cd.RefTable == nil {
			return nil, p.unexpectedToken()
		}
		if p.cur.Type == '(' {
			var pilErr error
			cd.RefColumns, pilErr = p.parseParenIdentList()
			if pilErr != nil {
				return nil, pilErr
			}
		}
		if err := p.parseReferentialActions(cd); err != nil {
			return nil, err
		}
		// [ NOT FOR REPLICATION ]
		if p.cur.Type == kwNOT && p.peekNext().Type == kwFOR {
			if err := p.parseNotForReplication(); err != nil {
				return nil, err
			}
			cd.NotForReplication = true
		}
	default:
		return nil, nil
	}

	cd.Loc.End = p.prevEnd()
	return cd, nil
}

// parseConstraintWithOptions parses [ WITH FILLFACTOR = N | WITH ( index_options ) ] on PK/UNIQUE constraints.
func (p *Parser) parseConstraintWithOptions(cd *nodes.ConstraintDef) {
	if p.cur.Type != kwWITH {
		return
	}
	p.advance()
	if p.cur.Type == '(' {
		cd.IndexOptions, _ = p.parseAlterIndexOptions()
	} else if p.cur.Type == kwFILLFACTOR {
		p.advance()
		if p.cur.Type == '=' {
			p.advance()
		}
		if p.cur.Type == tokICONST {
			cd.Fillfactor = int(p.cur.Ival)
			p.advance()
		}
	}
}

// parseConstraintOnFilegroup parses [ ON { partition_scheme(col) | filegroup | "default" } ].
// uniqueConstraintTailOpt (TSql160.g:29806) writes it On
// filegroupOrPartitionScheme, a stringOrIdentifier (TSql160.g:17224): SQL
// Server takes UNIQUE (a) ON 'a' and refuses a bare ON (measured).
func (p *Parser) parseConstraintOnFilegroup(cd *nodes.ConstraintDef) error {
	if p.cur.Type != kwON {
		return nil
	}
	p.advance()
	switch {
	case p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST:
		cd.OnFilegroup = p.cur.Str
		p.advance()
	case p.isIdentLike() || p.cur.Type == kwPRIMARY:
		name := p.cur.Str
		p.advance()
		if p.cur.Type == '(' {
			p.advance()
			col, ok := p.parseIdentifier()
			if !ok {
				return p.syntaxErrorAtCur()
			}
			if _, err := p.expect(')'); err != nil {
				return err
			}
			cd.OnFilegroup = name + "(" + col + ")"
		} else {
			cd.OnFilegroup = name
		}
	default:
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseColumnConstraintColumnList parses the column list SQL Server also
// takes on a PRIMARY KEY or UNIQUE constraint written beside a column, as
// `ALTER TABLE a ADD x a UNIQUE(x)` and `ADD CONTAINED a PRIMARY KEY(a)`
// (measured): the constraint is then a table constraint with no columns of
// its own, and the list names them.
func (p *Parser) parseColumnConstraintColumnList(cd *nodes.ConstraintDef) error {
	if p.cur.Type != '(' {
		return nil
	}
	cols, err := p.parseIndexColumnList()
	if err != nil {
		return err
	}
	cd.Columns = cols
	return nil
}

// parseTableConstraint parses a table-level constraint.
//
// BNF: mssql/parser/bnf/create-table-transact-sql.bnf
//
//	<table_constraint> ::=
//	[ CONSTRAINT constraint_name ]
//	{
//	    { PRIMARY KEY | UNIQUE }
//	        [ CLUSTERED | NONCLUSTERED ]
//	        ( column_name [ ASC | DESC ] [ ,...n ] )
//	        [ WITH FILLFACTOR = fillfactor | WITH ( <index_option> [ ,...n ] ) ]
//	        [ ON { partition_scheme_name (partition_column_name)
//	            | filegroup | "default" } ]
//	    | FOREIGN KEY
//	        ( column_name [ ,...n ] )
//	        REFERENCES referenced_table_name [ ( ref_column [ ,...n ] ) ]
//	        [ ON DELETE { NO ACTION | CASCADE | SET NULL | SET DEFAULT } ]
//	        [ ON UPDATE { NO ACTION | CASCADE | SET NULL | SET DEFAULT } ]
//	        [ NOT FOR REPLICATION ]
//	    | CHECK [ NOT FOR REPLICATION ] ( logical_expression )
//	}
func (p *Parser) parseTableConstraint() (*nodes.ConstraintDef, error) {
	return p.parseTableConstraintFor(false)
}

// parseAlterTableAddConstraint parses a table constraint of ALTER TABLE's ADD
// list, the one place a DEFAULT ... FOR column constraint is taken
// (defaultTableConstraint, TSql160.g:29670): CREATE TABLE refuses it (Msg
// 102, measured, compat.tsv).
func (p *Parser) parseAlterTableAddConstraint() (*nodes.ConstraintDef, error) {
	return p.parseTableConstraintFor(true)
}

func (p *Parser) parseTableConstraintFor(alterTableAdd bool) (*nodes.ConstraintDef, error) {
	loc := p.pos()
	var name string

	if p.cur.Type == kwCONSTRAINT {
		p.advance()
		var ok bool
		if name, ok = p.parseIdentifier(); !ok {
			return nil, p.syntaxErrorAtCur()
		}
	}

	cd := &nodes.ConstraintDef{
		Name: name,
		Loc:  nodes.Loc{Start: loc, End: -1},
	}

	switch p.cur.Type {
	case kwPRIMARY:
		p.advance() // PRIMARY
		if _, err := p.expect(kwKEY); err != nil {
			return nil, err
		}
		cd.Type = nodes.ConstraintPrimaryKey
		p.parseClusteredOption(cd)
		if p.cur.Type != '(' {
			return nil, p.syntaxErrorAtCur()
		}
		var err error
		cd.Columns, err = p.parseIndexColumnList()
		if err != nil {
			return nil, err
		}
		p.parseConstraintWithOptions(cd)
		if err := p.parseConstraintOnFilegroup(cd); err != nil {
			return nil, err
		}
	case kwUNIQUE:
		p.advance()
		cd.Type = nodes.ConstraintUnique
		p.parseClusteredOption(cd)
		if p.cur.Type != '(' {
			return nil, p.syntaxErrorAtCur()
		}
		var err error
		cd.Columns, err = p.parseIndexColumnList()
		if err != nil {
			return nil, err
		}
		p.parseConstraintWithOptions(cd)
		if err := p.parseConstraintOnFilegroup(cd); err != nil {
			return nil, err
		}
	case kwCHECK:
		p.advance()
		cd.Type = nodes.ConstraintCheck
		// [ NOT FOR REPLICATION ]
		if p.cur.Type == kwNOT && p.peekNext().Type == kwFOR {
			if err := p.parseNotForReplication(); err != nil {
				return nil, err
			}
			cd.NotForReplication = true
		}
		if _, err := p.expect('('); err != nil {
			return nil, err
		}
		// A CHECK constraint's expression is a search condition.
		p.enterSearchCondition()
		var exprErr error
		cd.Expr, exprErr = p.parseExpr()
		p.leaveSearchCondition()
		if exprErr != nil {
			return nil, exprErr
		}
		if cd.Expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		if _, err := p.expect(')'); err != nil {
			return nil, err
		}
	case kwFOREIGN:
		p.advance() // FOREIGN
		if _, err := p.expect(kwKEY); err != nil {
			return nil, err
		}
		cd.Type = nodes.ConstraintForeignKey
		// The column list may be left out (foreignConstraintColumnsOpt).
		if p.cur.Type == '(' {
			var err error
			cd.Columns, err = p.parseParenIdentList()
			if err != nil {
				return nil, err
			}
		}
		if _, ok := p.match(kwREFERENCES); !ok {
			return nil, p.syntaxErrorAtCur()
		}
		{
			if p.collectMode() {
				p.addRuleCandidate("table_ref")
				return nil, errCollecting
			}
			var refErr error
			cd.RefTable, refErr = p.parseTableRef()
			if refErr != nil {
				return nil, refErr
			}
			if cd.RefTable == nil {
				return nil, p.syntaxErrorAtCur()
			}
			if p.cur.Type == '(' {
				var pilErr error
				cd.RefColumns, pilErr = p.parseParenIdentList()
				if pilErr != nil {
					return nil, pilErr
				}
			}
			if err := p.parseReferentialActions(cd); err != nil {
				return nil, err
			}
		}
		// [ NOT FOR REPLICATION ]
		if p.cur.Type == kwNOT && p.peekNext().Type == kwFOR {
			if err := p.parseNotForReplication(); err != nil {
				return nil, err
			}
			cd.NotForReplication = true
		}
	case kwDEFAULT:
		if !alterTableAdd {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		cd.Type = nodes.ConstraintDefault
		var exprErr error
		cd.Expr, exprErr = p.parseExpr()
		if exprErr != nil {
			return nil, exprErr
		}
		if cd.Expr == nil {
			return nil, p.syntaxErrorAtCur()
		}
		// FOR column [ WITH VALUES ] (defaultConstraint, TSql160.g)
		if _, err := p.expect(kwFOR); err != nil {
			return nil, err
		}
		colName, ok := p.parseIdentifier()
		if !ok {
			return nil, p.syntaxErrorAtCur()
		}
		cd.Columns = &nodes.List{Items: []nodes.Node{&nodes.String{Str: colName}}}
		if p.cur.Type == kwWITH && p.peekNext().Type == kwVALUES {
			p.advance() // WITH
			p.advance() // VALUES
		}
	default:
		// Check for EDGE CONSTRAINT: CONSTRAINT name CONNECTION (...)
		if p.cur.Type == kwCONNECTION {
			cd.Type = nodes.ConstraintEdge
			var ecErr error
			cd.EdgeConnections, ecErr = p.parseEdgeConstraintConnections()
			if ecErr != nil {
				return nil, ecErr
			}
			// A graph connection takes ON DELETE { NO ACTION | CASCADE } and
			// nothing more: ON UPDATE, SET NULL, SET DEFAULT and a second ON
			// are Msg 156 (measured).
			if p.cur.Type == kwON {
				p.advance()
				if _, err := p.expect(kwDELETE); err != nil {
					return nil, err
				}
				switch {
				case p.cur.Type == kwNO:
					p.advance()
					if _, err := p.expect(kwACTION); err != nil {
						return nil, err
					}
					cd.OnDelete = nodes.RefActNoAction
				case p.cur.Type == kwCASCADE:
					p.advance()
					cd.OnDelete = nodes.RefActCascade
				default:
					return nil, p.syntaxErrorAtCur()
				}
			}
		} else {
			return nil, p.syntaxErrorAtCur()
		}
	}

	cd.Loc.End = p.prevEnd()
	return cd, nil
}

// parseEdgeConstraintConnections parses CONNECTION ( from_table TO to_table [, ...] ).
//
// Ref: https://learn.microsoft.com/en-us/sql/relational-databases/tables/graph-edge-constraints
//
//	CONNECTION ( from_table TO to_table [ , from_table TO to_table ... ] )
func (p *Parser) parseEdgeConstraintConnections() (*nodes.List, error) {
	p.advance() // consume CONNECTION

	var conns []nodes.Node
	if _, err := p.expect('('); err != nil {
		return nil, err
	}

	for {
		loc := p.pos()
		fromTable, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		if fromTable == nil {
			return nil, p.syntaxErrorAtCur()
		}
		if _, err := p.expect(kwTO); err != nil {
			return nil, err
		}
		toTable, err := p.parseTableRef()
		if err != nil {
			return nil, err
		}
		if toTable == nil {
			return nil, p.syntaxErrorAtCur()
		}
		conn := &nodes.EdgeConnectionDef{
			FromTable: fromTable,
			ToTable:   toTable,
			Loc:       nodes.Loc{Start: loc, End: p.prevEnd()},
		}
		conns = append(conns, conn)
		if _, ok := p.match(','); !ok {
			break
		}
	}

	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.List{Items: conns}, nil
}

// parseClusteredOption parses optional CLUSTERED/NONCLUSTERED.
func (p *Parser) parseClusteredOption(cd *nodes.ConstraintDef) {
	if p.cur.Type == kwCLUSTERED {
		p.advance()
		v := true
		cd.Clustered = &v
	} else if p.cur.Type == kwNONCLUSTERED {
		p.advance()
		v := false
		cd.Clustered = &v
	}
}

// parseReferentialActions parses ON DELETE/UPDATE actions.
func (p *Parser) parseReferentialActions(cd *nodes.ConstraintDef) error {
	for p.cur.Type == kwON {
		p.advance()
		if p.cur.Type == kwDELETE {
			p.advance()
			cd.OnDelete = p.parseRefAction()
			if cd.OnDelete == nodes.RefActNone {
				return p.syntaxErrorAtCur()
			}
		} else if p.cur.Type == kwUPDATE {
			p.advance()
			cd.OnUpdate = p.parseRefAction()
			if cd.OnUpdate == nodes.RefActNone {
				return p.syntaxErrorAtCur()
			}
		} else {
			// ON with neither DELETE nor UPDATE: SQL Server refuses it
			// (Msg 102), as it refuses an action word it does not know.
			return p.syntaxErrorAtCur()
		}
	}
	return nil
}

// parseRefAction parses a referential action (CASCADE, SET NULL, SET DEFAULT, NO ACTION).
func (p *Parser) parseRefAction() nodes.ReferentialAction {
	if _, ok := p.match(kwCASCADE); ok {
		return nodes.RefActCascade
	}
	if p.cur.Type == kwSET {
		p.advance()
		if _, ok := p.match(kwNULL); ok {
			return nodes.RefActSetNull
		}
		if _, ok := p.match(kwDEFAULT); ok {
			return nodes.RefActSetDefault
		}
	}
	// NO ACTION
	if p.cur.Type == kwNO {
		p.advance()
		if p.cur.Type == kwACTION {
			p.advance()
		}
		return nodes.RefActNoAction
	}
	return nodes.RefActNone
}

// parseInlineTableIndex parses an inline INDEX definition inside a CREATE
// TABLE body and in ALTER TABLE's ADD list (tableElement's
// inlineIndexTableDefinition, TSql160.g:28735).
//
// BNF: mssql/parser/bnf/create-table-transact-sql.bnf
//
//	<table_index> ::=
//	{
//	    INDEX index_name [ UNIQUE ] [ CLUSTERED | NONCLUSTERED ]
//	         ( column_name [ ASC | DESC ] [ ,...n ] )
//	    | INDEX index_name CLUSTERED COLUMNSTORE [ ORDER (column_name [ ,...n ] ) ]
//	    | INDEX index_name [ NONCLUSTERED ] COLUMNSTORE ( column_name [ ,...n ] )
//	    [ INCLUDE ( column_name [ ,...n ] ) ]
//	    [ WHERE <filter_predicate> ]
//	    [ WITH ( <index_option> [ ,...n ] ) ]
//	    [ ON { partition_scheme_name ( column_name )
//	         | filegroup_name | default } ]
//	    [ FILESTREAM_ON { filestream_filegroup_name | partition_scheme_name | "NULL" } ]
//	}
//
// What SQL Server refuses while it parses (its answers are in
// mssql/parser/testdata/compat.tsv): a name is required, a columnstore index
// that is not clustered takes a key list and a clustered one takes none
// (TSql160.g:28739's action), INCLUDE follows a key list, ON a filegroup, and
// a clustered index takes no filter (TSql160.g:17072's filterProhibited).
func (p *Parser) parseInlineTableIndex() (*nodes.InlineIndexDef, error) {
	return p.parseInlineIndex(false)
}

// parseInlineColumnIndex parses the INDEX a column definition may carry
// (inlineIndexColumnDefinition, TSql160.g:28800): the table index's form with
// the key list left out or given, and no COLUMNSTORE. DECLARE @t TABLE (a INT
// INDEX ix CLUSTERED) and CREATE TABLE a (a INT INDEX ix) parse (measured,
// compat.tsv).
func (p *Parser) parseInlineColumnIndex() (*nodes.InlineIndexDef, error) {
	return p.parseInlineIndex(true)
}

// parseInlineIndex parses a table's or, with column set, a column's INDEX.
// Either may name a HASH index, which is nonclustered (inlineIndexType,
// TSql160.g:28624): INDEX ix CLUSTERED HASH is refused.
func (p *Parser) parseInlineIndex(column bool) (*nodes.InlineIndexDef, error) {
	loc := p.pos()
	p.advance() // consume INDEX

	idx := &nodes.InlineIndexDef{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// index_name
	name, ok := p.parseIdentifier()
	if !ok {
		return nil, p.unexpectedToken()
	}
	idx.Name = name

	// [ UNIQUE ]
	if p.cur.Type == kwUNIQUE {
		idx.Unique = true
		p.advance()
	}

	// [ CLUSTERED | NONCLUSTERED ]
	if p.cur.Type == kwCLUSTERED {
		v := true
		idx.Clustered = &v
		p.advance()
		// CLUSTERED COLUMNSTORE [ ORDER ( col [,...n] ) ]
		if p.cur.Type == kwCOLUMNSTORE && !column {
			p.advance()
			idx.Columnstore = true
			if p.cur.Type == kwORDER {
				p.advance()
				if p.cur.Type != '(' {
					return nil, p.unexpectedToken()
				}
				var err error
				idx.Columns, err = p.parseIndexColumnList()
				if err != nil {
					return nil, err
				}
			}
			goto trailingOptions
		}
	} else if p.cur.Type == kwNONCLUSTERED {
		v := false
		idx.Clustered = &v
		p.advance()
	}

	// [ NONCLUSTERED ] HASH
	if p.cur.Type == kwHASH {
		if idx.Clustered != nil && *idx.Clustered {
			return nil, p.syntaxErrorAtCur()
		}
		idx.Hash = true
		p.advance()
	}

	// [ NONCLUSTERED ] COLUMNSTORE ( column_name [,...n] ), which a column's
	// index refuses.
	if p.cur.Type == kwCOLUMNSTORE && column {
		return nil, p.syntaxErrorAtCur()
	}
	if p.cur.Type == kwCOLUMNSTORE && !idx.Hash {
		p.advance()
		idx.Columnstore = true
		if p.cur.Type != '(' {
			return nil, p.newParseError(p.cur.Loc, "a nonclustered columnstore index takes a key list")
		}
		var err error
		idx.Columns, err = p.parseIndexColumnList()
		if err != nil {
			return nil, err
		}
		goto trailingOptions
	}

	// ( column_name [ ASC | DESC ] [ ,...n ] )
	if p.cur.Type == '(' {
		var err error
		idx.Columns, err = p.parseIndexColumnList()
		if err != nil {
			return nil, err
		}
	}

trailingOptions:
	// [ INCLUDE ( column_name [ ,...n ] ) ]
	if p.cur.Type == kwINCLUDE {
		if idx.Columns == nil {
			return nil, p.unexpectedToken()
		}
		p.advance()
		if p.cur.Type != '(' {
			return nil, p.unexpectedToken()
		}
		var err error
		idx.IncludeCols, err = p.parseParenIdentList()
		if err != nil {
			return nil, err
		}
	}

	// [ WHERE <filter_predicate> ]
	if p.cur.Type == kwWHERE {
		if idx.Clustered != nil && *idx.Clustered {
			return nil, p.newParseError(p.cur.Loc, "a clustered index takes no filter")
		}
		p.advance()
		var err error
		idx.WhereClause, err = p.indexFilterExpression()
		if err != nil {
			return nil, err
		}
	}

	// [ WITH ( <index_option> [ ,...n ] ) ]
	if p.cur.Type == kwWITH {
		p.advance()
		if p.cur.Type != '(' {
			return nil, p.unexpectedToken()
		}
		var err error
		idx.Options, err = p.indexOptionListRule(inlineIndexOptionRule(idx.Columnstore, idx.Hash))
		if err != nil {
			return nil, err
		}
	}
	// A HASH index takes a WITH list: INDEX ix HASH (a) is refused, and
	// INDEX ix HASH (a) WITH (FILLFACTOR = 80) parses (measured, compat.tsv).
	if idx.Hash && idx.Options == nil {
		return nil, p.syntaxErrorAtCur()
	}

	// [ ON { partition_scheme_name ( column_name ) | filegroup_name | default } ]
	// filegroupOrPartitionScheme is stringOrIdentifier (TSql160.g:17224).
	if p.cur.Type == kwON {
		p.advance()
		switch {
		case p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST:
			idx.OnFilegroup = p.cur.Str
			p.advance()
		case p.isIdentLike():
			fg := p.cur.Str
			p.advance()
			if p.cur.Type == '(' {
				p.advance()
				col, ok := p.parseIdentifier()
				if !ok {
					return nil, p.unexpectedToken()
				}
				if _, err := p.expect(')'); err != nil {
					return nil, err
				}
				idx.OnFilegroup = fg + "(" + col + ")"
			} else {
				idx.OnFilegroup = fg
			}
		default:
			return nil, p.unexpectedToken()
		}
	}

	// [ FILESTREAM_ON { filestream_filegroup_name | partition_scheme_name | "NULL" } ]
	// fileStreamOn is 'FILESTREAM_ON' stringOrIdentifier (TSql160.g:29200).
	if p.cur.Type == kwFILESTREAM_ON {
		p.advance()
		switch {
		case p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST:
			idx.FilestreamOn = p.cur.Str
			p.advance()
		case p.isIdentLike():
			idx.FilestreamOn = p.cur.Str
			p.advance()
		default:
			return nil, p.unexpectedToken()
		}
	}

	idx.Loc.End = p.prevEnd()
	return idx, nil
}

// isCTAS detects whether the current token stream (after table name) represents
// a CREATE TABLE AS SELECT statement. This is called during lookahead and the
// caller will restore lexer state afterwards. The tableName parameter is unused
// but kept for signature consistency with the dispatch.
//
// CTAS patterns after table name:
//
//	WITH ( DISTRIBUTION = ... ) AS SELECT ...
//	( col_list ) WITH ( ... ) AS SELECT ...
//	( col_list ) AS SELECT ...
//	AS SELECT ...
func (p *Parser) isCTAS(_ *nodes.TableRef) bool {
	// Skip optional column list: ( col1, col2, ... )
	if p.cur.Type == '(' {
		depth := 1
		p.advance()
		for depth > 0 && p.cur.Type != tokEOF {
			if p.cur.Type == '(' {
				depth++
			} else if p.cur.Type == ')' {
				depth--
			}
			p.advance()
		}
	}

	// Check for WITH ( DISTRIBUTION = ... ) pattern (Synapse table options)
	if p.cur.Type == kwWITH && p.peekNext().Type == '(' {
		p.advance() // consume WITH
		depth := 1
		p.advance() // consume (
		for depth > 0 && p.cur.Type != tokEOF {
			if p.cur.Type == '(' {
				depth++
			} else if p.cur.Type == ')' {
				depth--
			}
			p.advance()
		}
	}

	// After optional column list and optional WITH, check for AS SELECT or AS WITH (CTE)
	if p.cur.Type == kwAS {
		next := p.peekNext()
		if next.Type == kwSELECT || next.Type == kwWITH || next.Type == '(' {
			return true
		}
	}
	return false
}

// parseCreateTableAsSelectStmt parses a CREATE TABLE AS SELECT (CTAS) statement
// for Azure Synapse Analytics / Analytics Platform System.
//
// Ref: https://learn.microsoft.com/en-us/sql/t-sql/statements/create-table-as-select-azure-sql-data-warehouse
//
//	CREATE TABLE { database_name.schema_name.table_name | schema_name.table_name | table_name }
//	    [ ( column_name [ ,...n ] ) ]
//	    WITH (
//	      <distribution_option>
//	      [ , <table_option> [ ,...n ] ]
//	    )
//	    AS <select_statement>
//	    OPTION <query_hint>
//
//	<distribution_option> ::=
//	    {
//	        DISTRIBUTION = HASH ( distribution_column_name [, ...n] )
//	      | DISTRIBUTION = ROUND_ROBIN
//	      | DISTRIBUTION = REPLICATE
//	    }
//
//	<table_option> ::=
//	    {
//	        CLUSTERED COLUMNSTORE INDEX
//	      | CLUSTERED COLUMNSTORE INDEX ORDER ( column [,...n] )
//	      | HEAP
//	      | CLUSTERED INDEX ( { index_column_name [ ASC | DESC ] } [ ,...n ] )
//	    }
//	    | PARTITION ( partition_column_name RANGE [ LEFT | RIGHT ]
//	        FOR VALUES ( [ boundary_value [,...n] ] ) )
func (p *Parser) parseCreateTableAsSelectStmt() (*nodes.CreateTableAsSelectStmt, error) {
	loc := p.pos()

	stmt := &nodes.CreateTableAsSelectStmt{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// Table name
	var err error
	stmt.Name, err = p.parseTableRef()
	if err != nil {
		return nil, err
	}

	// Optional column list: ( col1, col2, ... )
	columns := false
	if p.cur.Type == '(' {
		// Peek ahead to check if this is a simple column name list (CTAS)
		// or a column definition list (regular CREATE TABLE).
		// In CTAS, the paren list contains only identifiers separated by commas.
		stmt.Columns, err = p.parseParenIdentList()
		if err != nil {
			return nil, err
		}
		columns = true
	}

	// WITH ( distribution_option [, table_option ...] ), required when the
	// statement names its columns: ctasCreateTableStatement is
	// ( columnNameList )? withTableOptions As selectStatement, and the only
	// alternative that reaches AS without WITH writes no column list
	// (TSql160.g:523-535, 26575-26584).
	if columns && p.cur.Type != kwWITH {
		return nil, p.syntaxErrorAtCur()
	}
	if p.cur.Type == kwWITH {
		p.advance() // consume WITH
		stmt.Options, err = p.parseCTASWithOptions()
		if err != nil {
			return nil, err
		}
	}

	// AS <select_statement>
	if p.cur.Type == kwAS {
		p.advance() // consume AS
	}

	// Parse the SELECT statement
	stmt.Query, err = p.parseSelectStmtWithCTE()
	if err != nil {
		return nil, err
	}

	// Optional OPTION ( query_hint )
	// This is already handled inside parseSelectStmt via parseOptionClause

	stmt.Loc.End = p.prevEnd()
	return stmt, nil
}

// parseCTASWithOptions parses the WITH ( ... ) clause for a CTAS statement.
// Options include DISTRIBUTION, CLUSTERED COLUMNSTORE INDEX, CLUSTERED INDEX,
// HEAP, and PARTITION.
func (p *Parser) parseCTASWithOptions() (*nodes.List, error) {
	if p.cur.Type != '(' {
		return nil, nil
	}
	p.advance() // consume (

	var items []nodes.Node
	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		opt, err := p.parseCTASOption()
		if err != nil {
			return nil, err
		}
		if opt != nil {
			items = append(items, opt)
		}
		if _, ok := p.match(','); !ok {
			break
		}
	}
	closing := p.cur.Loc
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	// withTableOptions is one option at least (TSql160.g:26710-26723), so WITH ()
	// is refused at the parenthesis that closes the empty list.
	if len(items) == 0 {
		return nil, p.newParseError(closing, `syntax error at or near ")"`)
	}
	return &nodes.List{Items: items}, nil
}

// parseCTASOption parses a single CTAS WITH option.
func (p *Parser) parseCTASOption() (*nodes.TableOption, error) {
	loc := p.pos()

	// CLUSTERED COLUMNSTORE INDEX [ORDER (col, ...)]
	if p.cur.Type == kwCLUSTERED {
		next := p.peekNext()
		if next.Type == kwCOLUMNSTORE {
			p.advance() // consume CLUSTERED
			p.advance() // consume COLUMNSTORE
			if p.cur.Type == kwINDEX {
				p.advance() // consume INDEX
			}
			opt := &nodes.TableOption{
				Name: "CLUSTERED COLUMNSTORE INDEX",
				Loc:  nodes.Loc{Start: loc, End: -1},
			}
			// Optional ORDER (col, ...)
			if p.cur.Type == kwORDER {
				p.advance() // consume ORDER
				if p.cur.Type == '(' {
					var cols []string
					p.advance()
					for p.cur.Type != ')' && p.cur.Type != tokEOF {
						name, ok := p.parseIdentifier()
						if !ok {
							break
						}
						cols = append(cols, name)
						if _, ok := p.match(','); !ok {
							break
						}
					}
					p.expect(')')
					opt.Value = "ORDER(" + strings.Join(cols, ", ") + ")"
				}
			}
			opt.Loc.End = p.prevEnd()
			return opt, nil
		}
		// CLUSTERED INDEX ( col [ASC|DESC], ... )
		if next.Type == kwINDEX {
			p.advance() // consume CLUSTERED
			p.advance() // consume INDEX
			opt := &nodes.TableOption{
				Name: "CLUSTERED INDEX",
				Loc:  nodes.Loc{Start: loc, End: -1},
			}
			if p.cur.Type == '(' {
				p.advance()
				var parts []string
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					name, ok := p.parseIdentifier()
					if !ok {
						break
					}
					entry := name
					if p.cur.Type == kwASC {
						entry += " ASC"
						p.advance()
					} else if p.cur.Type == kwDESC {
						entry += " DESC"
						p.advance()
					}
					parts = append(parts, entry)
					if _, ok := p.match(','); !ok {
						break
					}
				}
				p.expect(')')
				opt.Value = strings.Join(parts, ", ")
			}
			opt.Loc.End = p.prevEnd()
			return opt, nil
		}
	}

	// HEAP
	if p.cur.Type == kwHEAP {
		p.advance()
		opt := &nodes.TableOption{
			Name: "HEAP",
			Loc:  nodes.Loc{Start: loc, End: -1},
		}
		opt.Loc.End = p.prevEnd()
		return opt, nil
	}

	// PARTITION ( col RANGE [LEFT|RIGHT] FOR VALUES ( ... ) )
	if p.cur.Type == kwPARTITION {
		p.advance() // consume PARTITION
		opt := &nodes.TableOption{
			Name: "PARTITION",
			Loc:  nodes.Loc{Start: loc, End: -1},
		}
		if p.cur.Type == '(' {
			p.advance()
			// partition_column_name
			colName, _ := p.parseIdentifier()
			opt.Value = colName

			// RANGE [LEFT|RIGHT]
			rangeDir := ""
			if p.cur.Type == kwRANGE {
				p.advance()
				if p.cur.Type == kwLEFT {
					rangeDir = "LEFT"
					p.advance()
				} else if p.cur.Type == kwRIGHT {
					rangeDir = "RIGHT"
					p.advance()
				}
			}
			if rangeDir != "" {
				opt.Value += " RANGE " + rangeDir
			} else {
				opt.Value += " RANGE"
			}

			// FOR VALUES ( value, ... )
			if p.cur.Type == kwFOR {
				p.advance() // consume FOR
				if p.cur.Type == kwVALUES {
					p.advance() // consume VALUES
				}
				if p.cur.Type == '(' {
					p.advance()
					var vals []string
					for p.cur.Type != ')' && p.cur.Type != tokEOF {
						// Boundary values can be integers, strings, dates, etc.
						val := p.cur.Str
						if p.cur.Type == tokICONST || p.cur.Type == tokFCONST {
							val = p.cur.Str
						} else if p.cur.Type == tokSCONST || p.cur.Type == tokNSCONST {
							val = "'" + p.cur.Str + "'"
						} else if p.cur.Type == '-' {
							p.advance()
							val = "-" + p.cur.Str
						}
						vals = append(vals, val)
						p.advance()
						if _, ok := p.match(','); !ok {
							break
						}
					}
					p.expect(')')
					opt.Value += " FOR VALUES(" + strings.Join(vals, ", ") + ")"
				}
			}
			p.expect(')')
		}
		opt.Loc.End = p.prevEnd()
		return opt, nil
	}

	// DISTRIBUTION = HASH(col [,...]) | ROUND_ROBIN | REPLICATE
	if p.cur.Type == kwDISTRIBUTION {
		p.advance() // consume DISTRIBUTION
		opt := &nodes.TableOption{
			Name: "DISTRIBUTION",
			Loc:  nodes.Loc{Start: loc, End: -1},
		}
		if p.cur.Type == '=' {
			p.advance()
		}
		if p.isValidOption(distributionTypeValues) {
			distType := strings.ToUpper(p.cur.Str)
			p.advance()
			if distType == "HASH" && p.cur.Type == '(' {
				p.advance()
				var cols []string
				for p.cur.Type != ')' && p.cur.Type != tokEOF {
					name, ok := p.parseIdentifier()
					if !ok {
						break
					}
					cols = append(cols, name)
					if _, ok := p.match(','); !ok {
						break
					}
				}
				p.expect(')')
				opt.Value = "HASH(" + strings.Join(cols, ", ") + ")"
			} else {
				opt.Value = distType
			}
		}
		opt.Loc.End = p.prevEnd()
		return opt, nil
	}

	// Fallback: parse as NAME = VALUE
	return p.parseOneTableOption()
}

// parseParenIdentList parses (ident, ident, ...).
func (p *Parser) parseParenIdentList() (*nodes.List, error) {
	p.advance() // consume (
	if p.collectMode() {
		p.addRuleCandidate("columnref")
		return nil, errCollecting
	}
	items, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
		if p.collectMode() {
			p.addRuleCandidate("columnref")
			return nil, errCollecting
		}
		name, ok := p.parseIdentifier()
		if !ok {
			return nil, p.unexpectedToken()
		}
		return &nodes.String{Str: name}, nil
	})
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}
	return &nodes.List{Items: items}, nil
}

// parseNotForReplication consumes NOT FOR REPLICATION, the cursor on NOT.
// replicationClauseOpt (TSql160.g:29984) spells all three words, and SQL
// Server refuses CHECK NOT FOR (1 = 1) (measured, Msg 102).
func (p *Parser) parseNotForReplication() error {
	p.advance() // NOT
	p.advance() // FOR
	if p.cur.Type != kwREPLICATION {
		return p.syntaxErrorAtCur()
	}
	p.advance() // REPLICATION
	return nil
}
