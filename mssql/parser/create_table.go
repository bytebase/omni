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

	// AS NODE | AS EDGE (graph tables) | AS FILETABLE
	// Note: AS CLONE OF is handled at the dispatch level (parser.go)
	if p.cur.Type == kwAS {
		next := p.peekNext()
		if next.Type == kwNODE {
			p.advance() // AS
			p.advance() // NODE
			stmt.IsNode = true
		} else if next.Type == kwEDGE {
			p.advance() // AS
			p.advance() // EDGE
			stmt.IsEdge = true
		} else if next.Type == kwFILETABLE {
			p.advance() // AS
			p.advance() // FILETABLE
			stmt.IsFileTable = true
		}
	}

	// Column and constraint definitions
	if _, err := p.expect('('); err != nil {
		stmt.Loc.End = p.prevEnd()
		return stmt, nil
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
	_, listErr := p.parseCommaList(')', commaListAllowTrail, func() (nodes.Node, error) {
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
		// PERIOD FOR SYSTEM_TIME — stored on stmt directly, not collected here.
		if p.cur.Type == kwPERIOD && p.peekNext().Type == kwFOR {
			p.advance() // PERIOD
			p.advance() // FOR
			if p.cur.Type == kwSYSTEM_TIME {
				p.advance()
			}
			if p.cur.Type == '(' {
				p.advance()
				startCol, _ := p.parseIdentifier()
				stmt.PeriodStartCol = startCol
				p.match(',')
				endCol, _ := p.parseIdentifier()
				stmt.PeriodEndCol = endCol
				p.expect(')')
			}
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
		col, err := p.parseColumnDef()
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
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	if len(cols) > 0 {
		stmt.Columns = &nodes.List{Items: cols}
	}
	if len(constraints) > 0 {
		stmt.Constraints = &nodes.List{Items: constraints}
	}

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
		if (p.isIdentLike() || p.cur.Type == kwPRIMARY) {
			stmt.TextImageOn = p.cur.Str
			p.advance()
		}
	}

	// FILESTREAM_ON filegroup
	if p.cur.Type == kwFILESTREAM_ON {
		p.advance()
		if (p.isIdentLike() || p.cur.Type == kwPRIMARY) {
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

	// Expect '='
	if p.cur.Type != '=' {
		opt.Loc.End = p.prevEnd()
		return opt, nil
	}
	p.advance()

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
		// Simple NAME = VALUE
		if p.isIdentLike() {
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
		p.advance()
		compLoc := p.pos()
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
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
			Loc:       nodes.Loc{Start: compLoc, End: -1},
		}
		// Computed columns can also have column constraints (PK, UNIQUE, etc.)
		for p.cur.Type == kwCONSTRAINT || p.cur.Type == kwPRIMARY || p.cur.Type == kwUNIQUE ||
			p.cur.Type == kwCHECK || p.cur.Type == kwFOREIGN || p.cur.Type == kwREFERENCES {
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

	// column_set_definition: column_set_name XML COLUMN_SET FOR ALL_SPARSE_COLUMNS
	if col.DataType != nil && strings.EqualFold(col.DataType.Name, "xml") &&
		p.cur.Type == kwCOLUMN_SET {
		p.advance() // COLUMN_SET
		if p.cur.Type == kwFOR {
			p.advance() // FOR
		}
		if p.cur.Type == kwALL_SPARSE_COLUMNS {
			p.advance()
		}
		col.IsColumnSet = true
		col.Loc.End = p.prevEnd()
		return col, nil
	}

	// Column-level options in any order
	for {
		consumed := false

		// FILESTREAM
		if p.cur.Type == kwFILESTREAM {
			col.Filestream = true
			p.advance()
			consumed = true
		}

		// SPARSE
		if p.cur.Type == kwSPARSE {
			col.Sparse = true
			p.advance()
			consumed = true
		}

		// ROWGUIDCOL
		if p.cur.Type == kwROWGUIDCOL {
			col.Rowguidcol = true
			p.advance()
			consumed = true
		}

		// HIDDEN (standalone, not after GENERATED ALWAYS)
		if p.cur.Type == kwHIDDEN {
			col.Hidden = true
			p.advance()
			consumed = true
		}

		// MASKED WITH (FUNCTION = 'mask_function')
		if p.cur.Type == kwMASKED {
			p.advance() // MASKED
			if p.cur.Type == kwWITH {
				p.advance() // WITH
				if p.cur.Type == '(' {
					p.advance()
					// FUNCTION = 'value'
					if p.cur.Type == kwFUNCTION {
						p.advance()
						if p.cur.Type == '=' {
							p.advance()
						}
						if p.cur.Type == tokSCONST {
							col.MaskFunction = p.cur.Str
							p.advance()
						}
					}
					p.expect(')')
				}
			}
			consumed = true
		}

		// ENCRYPTED WITH (...)
		if p.cur.Type == kwENCRYPTED {
			encWith, encErr := p.parseEncryptedWith()
			if encErr != nil {
				return nil, encErr
			}
			col.EncryptedWith = encWith
			consumed = true
		}

		// GENERATED ALWAYS AS { ROW | TRANSACTION_ID | SEQUENCE_NUMBER } { START | END } [ HIDDEN ]
		if p.cur.Type == kwGENERATED {
			genAlways, genErr := p.parseGeneratedAlways()
			if genErr != nil {
				return nil, genErr
			}
			col.GeneratedAlways = genAlways
			// GENERATED ALWAYS ... HIDDEN
			if p.cur.Type == kwHIDDEN {
				col.Hidden = true
				p.advance()
			}
			consumed = true
		}

		// IDENTITY
		if p.cur.Type == kwIDENTITY {
			identSpec, identErr := p.parseIdentitySpec()
			if identErr != nil {
				return nil, identErr
			}
			col.Identity = identSpec
			consumed = true
		}

		// NOT FOR REPLICATION / NOT NULL
		if p.cur.Type == kwNOT {
			next := p.peekNext()
			if next.Type == kwNULL {
				p.advance() // NOT
				p.advance() // NULL
				col.Nullable = &nodes.NullableSpec{NotNull: true, Loc: nodes.Loc{Start: p.pos(), End: -1}}
				consumed = true
			} else if next.Type == kwFOR {
				p.advance() // NOT
				p.advance() // FOR
				if p.cur.Type == kwREPLICATION {
					p.advance()
				}
				col.NotForReplication = true
				consumed = true
			}
		}

		// NULL
		if p.cur.Type == kwNULL {
			p.advance()
			col.Nullable = &nodes.NullableSpec{NotNull: false, Loc: nodes.Loc{Start: p.pos(), End: -1}}
			consumed = true
		}

		// DEFAULT
		if p.cur.Type == kwDEFAULT {
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
			consumed = true
		}

		// COLLATE
		if p.cur.Type == kwCOLLATE {
			p.advance()
			if p.isIdentLike() {
				col.Collation = p.cur.Str
				p.advance()
			}
			consumed = true
		}

		// CONSTRAINT (inline column constraint)
		if p.cur.Type == kwCONSTRAINT {
			constraint, cErr := p.parseColumnConstraint()
			if cErr != nil {
				return nil, cErr
			}
			if constraint != nil {
				if col.Constraints == nil {
					col.Constraints = &nodes.List{}
				}
				col.Constraints.Items = append(col.Constraints.Items, constraint)
			}
			consumed = true
		}

		// PRIMARY KEY / UNIQUE (without CONSTRAINT keyword)
		if p.cur.Type == kwPRIMARY || p.cur.Type == kwUNIQUE {
			constraint, cErr := p.parseInlineConstraint("")
			if cErr != nil {
				return nil, cErr
			}
			if constraint != nil {
				if col.Constraints == nil {
					col.Constraints = &nodes.List{}
				}
				col.Constraints.Items = append(col.Constraints.Items, constraint)
			}
			consumed = true
		}

		// CHECK (without CONSTRAINT keyword)
		if p.cur.Type == kwCHECK {
			constraint, cErr := p.parseInlineConstraint("")
			if cErr != nil {
				return nil, cErr
			}
			if constraint != nil {
				if col.Constraints == nil {
					col.Constraints = &nodes.List{}
				}
				col.Constraints.Items = append(col.Constraints.Items, constraint)
			}
			consumed = true
		}

		// REFERENCES (inline FK without CONSTRAINT keyword)
		if p.cur.Type == kwREFERENCES {
			constraint, cErr := p.parseInlineConstraint("")
			if cErr != nil {
				return nil, cErr
			}
			if constraint != nil {
				if col.Constraints == nil {
					col.Constraints = &nodes.List{}
				}
				col.Constraints.Items = append(col.Constraints.Items, constraint)
			}
			consumed = true
		}

		if !consumed {
			break
		}
	}

	col.Loc.End = p.prevEnd()
	return col, nil
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

	if p.cur.Type != kwWITH {
		spec.Loc.End = p.prevEnd()
		return spec, nil
	}
	p.advance() // WITH

	if p.cur.Type != '(' {
		spec.Loc.End = p.prevEnd()
		return spec, nil
	}
	p.advance()

	for p.cur.Type != ')' && p.cur.Type != tokEOF {
		if p.isIdentLike() {
			optName := strings.ToUpper(p.cur.Str)
			p.advance()
			if p.cur.Type == '=' {
				p.advance()
			}
			switch optName {
			case "COLUMN_ENCRYPTION_KEY":
				if p.isIdentLike() {
					spec.ColumnEncryptionKey = p.cur.Str
					p.advance()
				}
			case "ENCRYPTION_TYPE":
				if p.isValidOption(encryptionTypeValues) {
					spec.EncryptionType = strings.ToUpper(p.cur.Str)
					p.advance()
				}
			case "ALGORITHM":
				if p.cur.Type == tokSCONST {
					spec.Algorithm = p.cur.Str
					p.advance()
				}
			default:
				if p.isIdentLike() || p.cur.Type == tokSCONST {
					p.advance()
				}
			}
		}
		if _, ok := p.match(','); !ok {
			break
		}
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

	// ALWAYS
	if p.cur.Type == kwALWAYS {
		p.advance()
	}

	// AS
	if p.cur.Type == kwAS {
		p.advance()
	}

	// ROW | TRANSACTION_ID | SEQUENCE_NUMBER — ROW is kwROW (CoreKeyword),
	// TRANSACTION_ID / SEQUENCE_NUMBER are idents.
	if p.isValidOption(ledgerColumnKindValues) || p.cur.Type == kwROW {
		kind := strings.ToUpper(p.cur.Str)
		p.advance()
		spec.Kind = kind
	}

	// START | END — END is kwEND (CoreKeyword); START is an ident.
	if p.isValidOption(ledgerStartEndValues) || p.cur.Type == kwEND {
		startEnd := strings.ToUpper(p.cur.Str)
		if startEnd == "START" || startEnd == "END" {
			spec.StartEnd = startEnd
			p.advance()
		}
	}
	// Handle END keyword (it's a reserved keyword kwEND)
	if p.cur.Type == kwEND {
		spec.StartEnd = "END"
		p.advance()
	}

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
		_, _ = p.expect(')')
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
		p.parseConstraintWithOptions(cd)
		p.parseConstraintOnFilegroup(cd)
	case kwUNIQUE:
		p.advance()
		cd.Type = nodes.ConstraintUnique
		p.parseClusteredOption(cd)
		p.parseConstraintWithOptions(cd)
		p.parseConstraintOnFilegroup(cd)
	case kwCHECK:
		p.advance()
		cd.Type = nodes.ConstraintCheck
		// [ NOT FOR REPLICATION ]
		if p.cur.Type == kwNOT && p.peekNext().Type == kwFOR {
			p.advance() // NOT
			p.advance() // FOR
			if p.cur.Type == kwREPLICATION {
				p.advance()
			}
			cd.NotForReplication = true
		}
		if _, err := p.expect('('); err == nil {
			var exprErr error
			cd.Expr, exprErr = p.parseExpr()
			if exprErr != nil {
				return nil, exprErr
			}
			if cd.Expr == nil {
				return nil, p.unexpectedToken()
			}
			_, _ = p.expect(')')
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
	case kwREFERENCES:
		p.advance()
		// Completion: after REFERENCES → table_ref
		if p.collectMode() {
			p.addRuleCandidate("table_ref")
			return nil, errCollecting
		}
		cd.Type = nodes.ConstraintForeignKey
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
		p.parseReferentialActions(cd)
		// [ NOT FOR REPLICATION ]
		if p.cur.Type == kwNOT && p.peekNext().Type == kwFOR {
			p.advance() // NOT
			p.advance() // FOR
			if p.cur.Type == kwREPLICATION {
				p.advance()
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
func (p *Parser) parseConstraintOnFilegroup(cd *nodes.ConstraintDef) {
	if p.cur.Type != kwON {
		return
	}
	p.advance()
	if p.isIdentLike() || p.cur.Type == kwPRIMARY {
		name := p.cur.Str
		p.advance()
		if p.cur.Type == '(' {
			p.advance()
			col, _ := p.parseIdentifier()
			p.expect(')')
			cd.OnFilegroup = name + "(" + col + ")"
		} else {
			cd.OnFilegroup = name
		}
	}
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
	loc := p.pos()
	var name string

	if p.cur.Type == kwCONSTRAINT {
		p.advance()
		name, _ = p.parseIdentifier()
	}

	cd := &nodes.ConstraintDef{
		Name: name,
		Loc:  nodes.Loc{Start: loc, End: -1},
	}

	switch p.cur.Type {
	case kwPRIMARY:
		p.advance() // PRIMARY
		p.match(kwKEY)
		cd.Type = nodes.ConstraintPrimaryKey
		p.parseClusteredOption(cd)
		if p.cur.Type == '(' {
			var err error
			cd.Columns, err = p.parseIndexColumnList()
			if err != nil {
				return nil, err
			}
		}
		p.parseConstraintWithOptions(cd)
		p.parseConstraintOnFilegroup(cd)
	case kwUNIQUE:
		p.advance()
		cd.Type = nodes.ConstraintUnique
		p.parseClusteredOption(cd)
		if p.cur.Type == '(' {
			var err error
			cd.Columns, err = p.parseIndexColumnList()
			if err != nil {
				return nil, err
			}
		}
		p.parseConstraintWithOptions(cd)
		p.parseConstraintOnFilegroup(cd)
	case kwCHECK:
		p.advance()
		cd.Type = nodes.ConstraintCheck
		// [ NOT FOR REPLICATION ]
		if p.cur.Type == kwNOT && p.peekNext().Type == kwFOR {
			p.advance() // NOT
			p.advance() // FOR
			if p.cur.Type == kwREPLICATION {
				p.advance()
			}
			cd.NotForReplication = true
		}
		if _, err := p.expect('('); err == nil {
			var exprErr error
			cd.Expr, exprErr = p.parseExpr()
			if exprErr != nil {
				return nil, exprErr
			}
			_, _ = p.expect(')')
		}
	case kwFOREIGN:
		p.advance() // FOREIGN
		p.match(kwKEY)
		cd.Type = nodes.ConstraintForeignKey
		if p.cur.Type == '(' {
			var err error
			cd.Columns, err = p.parseParenIdentList()
			if err != nil {
				return nil, err
			}
		}
		if _, ok := p.match(kwREFERENCES); ok {
			var refErr error
			cd.RefTable, refErr = p.parseTableRef()
			if refErr != nil {
				return nil, refErr
			}
			if p.cur.Type == '(' {
				var pilErr error
				cd.RefColumns, pilErr = p.parseParenIdentList()
				if pilErr != nil {
					return nil, pilErr
				}
			}
			p.parseReferentialActions(cd)
		}
		// [ NOT FOR REPLICATION ]
		if p.cur.Type == kwNOT && p.peekNext().Type == kwFOR {
			p.advance() // NOT
			p.advance() // FOR
			if p.cur.Type == kwREPLICATION {
				p.advance()
			}
			cd.NotForReplication = true
		}
	case kwDEFAULT:
		p.advance()
		cd.Type = nodes.ConstraintDefault
		var exprErr error
		cd.Expr, exprErr = p.parseExpr()
		if exprErr != nil {
			return nil, exprErr
		}
		// FOR column
		if _, ok := p.match(kwFOR); ok {
			p.parseIdentifier() // column name (not stored separately but consumed)
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
			p.parseReferentialActions(cd)
		} else {
			return nil, nil
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
		p.match(kwTO)
		toTable, err := p.parseTableRef()
		if err != nil {
			return nil, err
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
func (p *Parser) parseReferentialActions(cd *nodes.ConstraintDef) {
	for p.cur.Type == kwON {
		p.advance()
		if p.cur.Type == kwDELETE {
			p.advance()
			cd.OnDelete = p.parseRefAction()
		} else if p.cur.Type == kwUPDATE {
			p.advance()
			cd.OnUpdate = p.parseRefAction()
		} else {
			break
		}
	}
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

// parseInlineTableIndex parses an inline INDEX definition inside a CREATE TABLE body.
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
func (p *Parser) parseInlineTableIndex() (*nodes.InlineIndexDef, error) {
	loc := p.pos()
	p.advance() // consume INDEX

	idx := &nodes.InlineIndexDef{
		Loc: nodes.Loc{Start: loc, End: -1},
	}

	// index_name
	name, _ := p.parseIdentifier()
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
		if p.cur.Type == kwCOLUMNSTORE {
			p.advance()
			idx.Columnstore = true
			if p.cur.Type == kwORDER {
				p.advance()
				if p.cur.Type == '(' {
					var err error
					idx.Columns, err = p.parseIndexColumnList()
					if err != nil {
						return nil, err
					}
				}
			}
			goto trailingOptions
		}
	} else if p.cur.Type == kwNONCLUSTERED {
		v := false
		idx.Clustered = &v
		p.advance()
	}

	// [ NONCLUSTERED ] COLUMNSTORE ( column_name [,...n] )
	if p.cur.Type == kwCOLUMNSTORE {
		p.advance()
		idx.Columnstore = true
		if p.cur.Type == '(' {
			var err error
			idx.Columns, err = p.parseIndexColumnList()
			if err != nil {
				return nil, err
			}
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
		p.advance()
		if p.cur.Type == '(' {
			var err error
			idx.IncludeCols, err = p.parseParenIdentList()
			if err != nil {
				return nil, err
			}
		}
	}

	// [ WHERE <filter_predicate> ]
	if _, ok := p.match(kwWHERE); ok {
		var err error
		idx.WhereClause, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
	}

	// [ WITH ( <index_option> [ ,...n ] ) ]
	if p.cur.Type == kwWITH {
		p.advance()
		if p.cur.Type == '(' {
			var err error
			idx.Options, err = p.parseOptionList()
			if err != nil {
				return nil, err
			}
		}
	}

	// [ ON { partition_scheme_name ( column_name ) | filegroup_name | default } ]
	if p.cur.Type == kwON {
		p.advance()
		if p.isIdentLike() || p.cur.Type == kwPRIMARY {
			fg := p.cur.Str
			p.advance()
			if p.cur.Type == '(' {
				p.advance()
				col, _ := p.parseIdentifier()
				p.expect(')')
				idx.OnFilegroup = fg + "(" + col + ")"
			} else {
				idx.OnFilegroup = fg
			}
		}
	}

	// [ FILESTREAM_ON { ... } ]
	if p.cur.Type == kwFILESTREAM_ON {
		p.advance()
		if (p.isIdentLike() || p.cur.Type == kwPRIMARY) {
			idx.FilestreamOn = p.cur.Str
			p.advance()
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
	if p.cur.Type == '(' {
		// Peek ahead to check if this is a simple column name list (CTAS)
		// or a column definition list (regular CREATE TABLE).
		// In CTAS, the paren list contains only identifiers separated by commas.
		stmt.Columns, err = p.parseParenIdentList()
		if err != nil {
			return nil, err
		}
	}

	// Optional WITH ( distribution_option [, table_option ...] )
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
	stmt.Query, err = p.parseSelectStmt()
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
	if _, err := p.expect(')'); err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, nil
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
	items, err := p.parseCommaList(')', commaListStrict, func() (nodes.Node, error) {
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
