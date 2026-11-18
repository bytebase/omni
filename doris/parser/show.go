package parser

import (
	"strings"

	"github.com/bytebase/omni/doris/ast"
)

// ---------------------------------------------------------------------------
// SHOW
// ---------------------------------------------------------------------------

// parseShow parses a SHOW statement.
//
// Strategy: consume SHOW, then dispatch on the next 1-3 keywords to identify
// the variant. For common variants (TABLES, DATABASES, COLUMNS, CREATE TABLE,
// VARIABLES, PARTITIONS, GRANTS) populate the ShowStmt fields properly.
// For all other variants, capture the remaining tokens as raw text in Args.
//
// On entry, SHOW has NOT yet been consumed (cur == kwSHOW).
func (p *Parser) parseShow() (ast.Node, error) {
	if p.strictTrailing {
		if err := p.checkShowGrammar(dorisShowForms); err != nil {
			return nil, err
		}
	}
	startLoc := p.cur.Loc
	p.advance() // consume SHOW

	stmt := &ast.ShowStmt{Loc: startLoc}

	// Optional FULL / EXTENDED / BRIEF / ALL modifiers before the variant keyword.
	if p.cur.Kind == kwFULL {
		stmt.Full = true
		p.advance()
	} else if p.cur.Kind == kwEXTENDED {
		stmt.Extended = true
		p.advance()
	} else if p.cur.Kind == kwBRIEF {
		// SHOW BRIEF CREATE TABLE ... — record and consume
		stmt.Args = "BRIEF"
		p.advance()
	} else if p.cur.Kind == kwALL {
		// SHOW ALL GRANTS — record ALL in args, let dispatch handle GRANTS
		stmt.Args = "ALL"
		p.advance()
	} else if p.cur.Kind == kwTEMPORARY {
		// SHOW TEMPORARY PARTITIONS — record TEMPORARY
		stmt.Args = "TEMPORARY"
		p.advance()
	}

	// SHOW JOB (T8.1) — special-case routing
	if p.cur.Kind == kwJOB {
		return p.parseShowJob(startLoc)
	}

	// SHOW ANALYZE / SHOW STATS / SHOW CONSTRAINTS (T8.3) — special-case routing
	if p.cur.Kind == kwANALYZE {
		all := stmt.Args == "ALL"
		return p.parseShowAnalyze(startLoc, all, false)
	}
	if p.cur.Kind == kwQUEUED && p.peekNext().Kind == kwANALYZE {
		p.advance() // consume QUEUED
		return p.parseShowAnalyze(startLoc, false, true)
	}
	if p.cur.Kind == kwCONSTRAINTS {
		return p.parseShowConstraints(startLoc)
	}
	if p.cur.Kind == kwSTATS {
		return p.parseShowStats(startLoc, "")
	}
	// SHOW COLUMN/INDEX/PARTITION/TABLE STATS
	if p.peekNext().Kind == kwSTATS {
		switch p.cur.Kind {
		case kwCOLUMN:
			p.advance()
			return p.parseShowStats(startLoc, "COLUMN")
		case kwINDEX:
			p.advance()
			return p.parseShowStats(startLoc, "INDEX")
		case kwPARTITION:
			p.advance()
			return p.parseShowStats(startLoc, "PARTITION")
		case kwTABLE:
			p.advance()
			return p.parseShowStats(startLoc, "TABLE")
		}
	}

	// SHOW [ALL] ROUTINE LOAD (T6.2) — special-case routing
	if p.cur.Kind == kwROUTINE {
		// SHOW ALL ROUTINE LOAD case: the ALL was already consumed and stored in stmt.Args.
		// parseShowRoutineLoad needs to know about it via a passthrough.
		all := stmt.Args == "ALL"
		node, err := p.parseShowRoutineLoad(startLoc)
		if err != nil {
			return nil, err
		}
		if all {
			if rl, ok := node.(*ast.ShowRoutineLoadStmt); ok {
				rl.All = true
			}
		}
		return node, nil
	}

	switch p.cur.Kind {
	case kwTABLES:
		return p.parseShowTables(stmt)
	case kwDATABASES, kwSCHEMAS:
		return p.parseShowDatabases(stmt)
	case kwCOLUMNS, kwFIELDS:
		return p.parseShowColumns(stmt)
	case kwCREATE:
		// SHOW CREATE ROUTINE LOAD FOR job_name (T6.2)
		if p.peekNext().Kind == kwROUTINE {
			p.advance() // consume CREATE
			return p.parseShowRoutineLoad(startLoc)
		}
		return p.parseShowCreate(stmt)
	case kwVARIABLES:
		return p.parseShowVariables(stmt)
	case kwPARTITIONS:
		return p.parseShowPartitions(stmt)
	case kwGRANTS:
		return p.parseShowGrants(stmt)
	case kwROLES:
		return p.parseShowGenericNamed(stmt, "ROLES")
	case kwCATALOGS:
		return p.parseShowCatalogs(stmt)
	case kwALTER:
		return p.parseShowAlterTable(stmt)
	case kwTABLE:
		// SHOW TABLE STATUS [FROM db] [LIKE 'pat']
		return p.parseShowTableStatus(stmt)
	default:
		// Generic fallthrough: capture the variant type word(s) plus remaining
		// tokens as raw text.
		return p.parseShowGeneric(stmt)
	}
}

// parseShowTables parses: SHOW [FULL] TABLES [FROM db] [LIKE 'pat'] [WHERE ...]
// On entry, cur == kwTABLES.
func (p *Parser) parseShowTables(stmt *ast.ShowStmt) (ast.Node, error) {
	p.advance() // consume TABLES
	stmt.Type = "TABLES"

	if err := p.parseShowFromLikeWhere(stmt); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseShowDatabases parses: SHOW DATABASES [FROM catalog] [LIKE 'pat'] [WHERE ...]
// On entry, cur == kwDATABASES or kwSCHEMAS.
func (p *Parser) parseShowDatabases(stmt *ast.ShowStmt) (ast.Node, error) {
	p.advance() // consume DATABASES/SCHEMAS
	stmt.Type = "DATABASES"

	// Optional FROM/IN <catalog>: SHOW DATABASES FROM hms_catalog.
	if p.cur.Kind == kwFROM || p.cur.Kind == kwIN {
		p.advance() // consume FROM/IN
		catalog, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.From = catalog
	}

	if err := p.parseShowLikeWhere(stmt); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseShowColumns parses: SHOW [FULL] COLUMNS FROM table [FROM db] [LIKE 'pat']
// On entry, cur == kwCOLUMNS or kwFIELDS.
func (p *Parser) parseShowColumns(stmt *ast.ShowStmt) (ast.Node, error) {
	p.advance() // consume COLUMNS/FIELDS
	stmt.Type = "COLUMNS"

	// FROM table
	if p.cur.Kind == kwFROM || p.cur.Kind == kwIN {
		p.advance() // consume FROM/IN
		target, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Target = target
	}

	// Optional second FROM db
	if p.cur.Kind == kwFROM || p.cur.Kind == kwIN {
		p.advance() // consume FROM/IN
		dbName, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.From = dbName
	}

	if err := p.parseShowLikeWhere(stmt); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseShowCreate parses: SHOW [BRIEF] CREATE {TABLE|VIEW|DATABASE|CATALOG|...} name
// On entry, cur == kwCREATE.
func (p *Parser) parseShowCreate(stmt *ast.ShowStmt) (ast.Node, error) {
	p.advance() // consume CREATE

	// Build type: "CREATE TABLE", "CREATE VIEW", "CREATE DATABASE", etc.
	typeWord := "CREATE"
	materializedView := false
	switch p.cur.Kind {
	case kwTABLE:
		typeWord = "CREATE TABLE"
		p.advance()
	case kwVIEW:
		typeWord = "CREATE VIEW"
		p.advance()
	case kwDATABASE, kwSCHEMA:
		typeWord = "CREATE DATABASE"
		p.advance()
	case kwCATALOG:
		typeWord = "CREATE CATALOG"
		p.advance()
	case kwLOAD:
		// SHOW CREATE LOAD FOR label (DorisParser.g4:391, the CUP parser's
		// show_create_load_stmt at DorisLegacyParser.g4:878).
		p.advance() // consume LOAD
		if _, err := p.expect(kwFOR); err != nil {
			return nil, err
		}
		typeWord = "CREATE LOAD"
	case kwMATERIALIZED:
		p.advance() // consume MATERIALIZED
		if p.cur.Kind == kwVIEW {
			p.advance()
		}
		typeWord = "CREATE MATERIALIZED VIEW"
		materializedView = true
	case kwFUNCTION, kwGLOBAL, kwSESSION, kwLOCAL:
		// SHOW CREATE [GLOBAL | SESSION | LOCAL] FUNCTION name(args) [(FROM | IN) db]
		return p.parseShowCreateFunction(stmt)
	default:
		// Unrecognized: let it fall through, name may be the next token
	}

	// Prepend BRIEF if it was the first modifier
	if stmt.Args == "BRIEF" {
		stmt.Type = "BRIEF " + typeWord
		stmt.Args = ""
	} else {
		stmt.Type = typeWord
	}

	// Optional target name
	if isIdentifierToken(p.cur.Kind) || p.cur.Kind == tokString {
		target, err := p.parseMultipartIdentifier()
		if err == nil {
			stmt.Target = target
		}
	}

	// SHOW CREATE MATERIALIZED VIEW name ON table (DorisParser.g4:355,
	// show_param at DorisLegacyParser.g4:1386).
	if materializedView && p.cur.Kind == kwON {
		p.advance() // consume ON
		on, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.On = on
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseShowCreateFunction reads the rest of
//
//	SHOW CREATE (GLOBAL | SESSION | LOCAL)? FUNCTION functionIdentifier
//	    LEFT_PAREN functionArguments? RIGHT_PAREN ((FROM | IN) multipartIdentifier)?
//
// (DorisParser.g4:291-293). cur is on the word after CREATE.
func (p *Parser) parseShowCreateFunction(stmt *ast.ShowStmt) (ast.Node, error) {
	stmt.Type = "CREATE FUNCTION"
	if p.cur.Kind == kwGLOBAL {
		stmt.Type = "CREATE GLOBAL FUNCTION"
	}
	if p.cur.Kind != kwFUNCTION {
		p.advance() // consume GLOBAL, SESSION or LOCAL
	}
	if _, err := p.expect(kwFUNCTION); err != nil {
		return nil, err
	}
	target, err := p.parseFunctionIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Target = target
	if _, _, err := p.parseFunctionArguments(); err != nil {
		return nil, err
	}
	if p.cur.Kind == kwFROM || p.cur.Kind == kwIN {
		p.advance() // consume FROM/IN
		db, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.From = strings.Join(db.Parts, ".")
	}
	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseShowVariables parses: SHOW [GLOBAL|SESSION] VARIABLES [LIKE 'pat'] [WHERE ...]
// On entry, cur == kwVARIABLES.
func (p *Parser) parseShowVariables(stmt *ast.ShowStmt) (ast.Node, error) {
	p.advance() // consume VARIABLES
	stmt.Type = "VARIABLES"

	if err := p.parseShowLikeWhere(stmt); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseShowPartitions parses:
//
//	SHOW [TEMPORARY] PARTITIONS FROM table [PARTITION(p1,...)]
//	    [WHERE ...] [ORDER BY ...] [LIMIT n]
//
// On entry, cur == kwPARTITIONS. stmt.Args may already contain "TEMPORARY".
func (p *Parser) parseShowPartitions(stmt *ast.ShowStmt) (ast.Node, error) {
	p.advance() // consume PARTITIONS
	stmt.Type = "PARTITIONS"
	// Note: stmt.Args may be "TEMPORARY" (set by parseShow before dispatch)

	// FROM table
	if p.cur.Kind == kwFROM || p.cur.Kind == kwIN {
		p.advance() // consume FROM
		target, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Target = target
	}

	// Remaining tokens captured as raw args (WHERE, ORDER BY, LIMIT, etc.)
	// Preserve pre-set Args (e.g. "TEMPORARY").
	rest := p.collectRemainingRaw()
	if rest != "" {
		if stmt.Args != "" {
			stmt.Args = stmt.Args + " " + rest
		} else {
			stmt.Args = rest
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseShowGrants parses: SHOW [ALL] GRANTS [FOR user]
// On entry, cur == kwGRANTS. stmt.Args may contain "ALL".
func (p *Parser) parseShowGrants(stmt *ast.ShowStmt) (ast.Node, error) {
	p.advance() // consume GRANTS
	stmt.Type = "GRANTS"

	// Optional FOR user_ident; preserve pre-set Args (e.g. "ALL").
	rest := p.collectRemainingRaw()
	if rest != "" {
		if stmt.Args != "" {
			stmt.Args = stmt.Args + " " + rest
		} else {
			stmt.Args = rest
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseShowTableStatus parses:
//
//	SHOW TABLE STATUS [FROM db[.table]] [LIKE 'pat' | WHERE ...]
//	SHOW TABLE CREATION [FROM db] [LIKE 'pat' | WHERE ...]
//	SHOW TABLE id
//
// (DorisParser.g4:275 #showTableStatus, :278 #showTableId and :351
// #showTableCreation; show_param at DorisLegacyParser.g4:1386 takes all
// three.) On entry, cur == kwTABLE.
func (p *Parser) parseShowTableStatus(stmt *ast.ShowStmt) (ast.Node, error) {
	p.advance() // consume TABLE
	stmt.Type = "TABLE"

	// SHOW TABLE CREATION: a variant of its own, not a table named CREATION.
	if p.cur.Kind == kwCREATION {
		p.advance() // consume CREATION
		stmt.Type = "TABLE CREATION"
		if err := p.parseShowFromLikeWhere(stmt); err != nil {
			return nil, err
		}
		stmt.Loc.End = p.prev.Loc.End
		return stmt, nil
	}

	// SHOW TABLE id names the table by id, which the AST keeps as the
	// variant's argument.
	if p.cur.Kind == tokInt {
		stmt.Args = p.sourceText(p.advance().Loc)
		stmt.Loc.End = p.prev.Loc.End
		return stmt, nil
	}

	// Optional STATUS
	if p.cur.Kind == kwSTATUS {
		p.advance()
		stmt.Type = "TABLE STATUS"
	}

	if err := p.parseShowFromLikeWhere(stmt); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseShowCatalogs parses: SHOW CATALOGS [LIKE 'pat']
// On entry, cur == kwCATALOGS.
func (p *Parser) parseShowCatalogs(stmt *ast.ShowStmt) (ast.Node, error) {
	p.advance() // consume CATALOGS
	stmt.Type = "CATALOGS"

	if err := p.parseShowLikeWhere(stmt); err != nil {
		return nil, err
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseShowAlterTable parses: SHOW ALTER TABLE {COLUMN|ROLLUP} [FROM db] [WHERE ...] ...
// On entry, cur == kwALTER.
func (p *Parser) parseShowAlterTable(stmt *ast.ShowStmt) (ast.Node, error) {
	p.advance() // consume ALTER

	if p.cur.Kind != kwTABLE {
		// Fallback: type is just "ALTER"
		stmt.Type = "ALTER"
		stmt.Args = p.collectRemainingRaw()
		stmt.Loc.End = p.prev.Loc.End
		return stmt, nil
	}
	p.advance() // consume TABLE

	subType := ""
	switch p.cur.Kind {
	case kwCOLUMN:
		subType = "COLUMN"
		p.advance()
	case kwROLLUP:
		subType = "ROLLUP"
		p.advance()
	default:
		subType = ""
	}

	if subType == "" {
		stmt.Type = "ALTER TABLE"
	} else {
		stmt.Type = "ALTER TABLE " + subType
	}

	// Optional FROM db
	if p.cur.Kind == kwFROM {
		p.advance()
		dbName, _, err := p.parseIdentifier()
		if err == nil {
			stmt.From = dbName
		} else if p.strictTrailing {
			// The raw capture below would swallow the token that is not a
			// database name, or there is none at all, so SHOW ALTER TABLE
			// COLUMN FROM and ... FROM 'x' parsed clean (container-verified
			// rejects). Best-effort keeps the raw recovery.
			return nil, err
		}
	}

	// Remaining tokens (WHERE, ORDER BY, LIMIT)
	stmt.Args = p.collectRemainingRaw()
	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseShowGenericNamed handles SHOW variants where the type is already known
// (e.g. ROLES, CATALOGS). It consumes the type keyword and captures the rest
// as Args.
func (p *Parser) parseShowGenericNamed(stmt *ast.ShowStmt, typeName string) (ast.Node, error) {
	p.advance() // consume the type keyword
	stmt.Type = typeName
	stmt.Args = p.collectRemainingRaw()
	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseShowGeneric handles all SHOW variants not specifically recognized above.
// It reads the next 1-2 tokens as the type word(s), then captures the rest
// as raw args.
func (p *Parser) parseShowGeneric(stmt *ast.ShowStmt) (ast.Node, error) {
	// Build the type from the next keyword(s).
	var typeParts []string

	// The BRIEF modifier may have been prepended in Args already.
	prefix := ""
	if stmt.Args == "BRIEF" {
		prefix = "BRIEF "
		stmt.Args = ""
	}

	// First type word.
	if p.cur.Kind != tokEOF {
		typeParts = append(typeParts, strings.ToUpper(p.cur.Str))
		p.advance()
	}

	// Second type word (for two-keyword variants like ALTER TABLE, TABLE STATUS, etc.)
	// Heuristic: consume a second keyword only when both are keyword tokens and
	// the combination commonly appears as a SHOW variant.
	if p.cur.Kind >= 700 && p.cur.Kind != tokEOF {
		second := strings.ToUpper(p.cur.Str)
		// Known two-word type prefixes.
		switch strings.Join(typeParts, "") + second {
		case "ALTERTABLE", "ALTERTABLECOLUMN",
			"TABLETSTATUS", "TABLESTATUS",
			"ALTERCOLUMN", "ALTERROLLUP":
			typeParts = append(typeParts, second)
			p.advance()
		default:
			// Keep only the first word.
		}
	}

	stmt.Type = prefix + strings.Join(typeParts, " ")
	stmt.Args = p.collectRemainingRaw()

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// Shared SHOW helpers
// ---------------------------------------------------------------------------

// parseShowFromLikeWhere optionally consumes [FROM db] [LIKE 'pat'] [WHERE expr]
// for SHOW TABLES and similar forms.
func (p *Parser) parseShowFromLikeWhere(stmt *ast.ShowStmt) error {
	// Optional FROM db: the CUP parser's opt_db names one database, and its
	// TABLE and TABLES forms also take `from_or_in ident '.' ident`
	// (DorisLegacyParser.g4:1386), which the ANTLR parser's
	// multipartIdentifier takes as well (DorisParser.g4:275-276).
	if p.cur.Kind == kwFROM || p.cur.Kind == kwIN {
		p.advance()
		if isIdentifierToken(p.cur.Kind) {
			name, err := p.parseMultipartIdentifier()
			if err == nil {
				stmt.From = strings.Join(name.Parts, ".")
			}
		} else if p.strictTrailing {
			// SHOW TABLES FROM with no database name ends at EOF, invisible
			// to the leftover-token check; the operand is required.
			return p.syntaxErrorAtCur()
		}
	}
	return p.parseShowLikeWhere(stmt)
}

// parseShowLikeWhere optionally consumes [LIKE 'pat'] [WHERE expr].
func (p *Parser) parseShowLikeWhere(stmt *ast.ShowStmt) error {
	if p.cur.Kind == kwLIKE {
		p.advance() // consume LIKE
		if p.cur.Kind == tokString {
			stmt.Like = p.cur.Str
			p.advance()
		} else if p.strictTrailing {
			// SHOW TABLES LIKE with no pattern ends at EOF, which the
			// leftover-token check cannot flag; the string is required. This
			// was once recorded as an engine accept, but the engine reports
			// "missing STRING_LITERAL at '<EOF>'", a parse failure the
			// container oracle used to miss (see dorisMissingTokenMarker).
			return p.syntaxErrorAtCur()
		}
		return nil
	}
	if p.cur.Kind == kwWHERE {
		p.advance() // consume WHERE
		expr, err := p.parseExpr()
		if err == nil {
			stmt.Where = expr
		} else if p.strictTrailing {
			// The clause parser may have consumed to EOF before failing
			// (SHOW TABLES WHERE '('), so the leftover-token check alone
			// cannot catch this — the error has to propagate.
			return err
		}
	}
	return nil
}

// collectRemainingRaw consumes all remaining tokens (up to EOF) and returns
// them joined by a single space. Used for variant-specific args we don't parse
// in detail.
func (p *Parser) collectRemainingRaw() string {
	var parts []string
	for p.cur.Kind != tokEOF {
		if p.cur.Str != "" {
			parts = append(parts, p.cur.Str)
		} else {
			// Single-char punctuation: use the actual character.
			if p.cur.Kind > 0 && p.cur.Kind < 128 {
				parts = append(parts, string(rune(p.cur.Kind)))
			}
		}
		p.advance()
	}
	return strings.Join(parts, " ")
}

// ---------------------------------------------------------------------------
// DESCRIBE / DESC
// ---------------------------------------------------------------------------

// parseDescribe parses DESCRIBE/DESC statements.
//
//	DESCRIBE [FULL] table_name [ALL VERBOSE]
//	DESC [FULL] table_name
//	DESCRIBE [FULL] table_name opt_partition_names
//	DESCRIBE FUNCTION name '(' [properties] ')' [table_alias]
//	DESCRIBE [explain options] query
//	DESCRIBE [explain options] INSERT ...
//
// The first four are DorisLegacyParser.g4:1629 describe_stmt, which the ANTLR
// parser's unsupportedDescribeStatement (DorisParser.g4:947-952) also derives;
// the last two are an explain command, so they parse as an ExplainStmt
// (DorisParser.g4:52, :1067).
//
// On entry, DESCRIBE or DESC has NOT yet been consumed.
func (p *Parser) parseDescribe() (ast.Node, error) {
	startLoc := p.cur.Loc
	p.advance() // consume DESCRIBE/DESC

	if p.describeExplains() {
		return p.parseExplainBody(&ast.ExplainStmt{Loc: startLoc})
	}

	stmt := &ast.DescribeStmt{Loc: startLoc}

	// DESCRIBE FUNCTION a() — table_valued_function_ref
	// (DorisLegacyParser.g4:2240), which DorisParser.g4:948 derives too. It is
	// a describe_stmt alternative of its own: FULL scopes to the table-name
	// alternative, so it is read after this.
	if p.cur.Kind == kwFUNCTION {
		if err := p.parseDescribeFunction(stmt); err != nil {
			return nil, err
		}
		stmt.Loc.End = p.prev.Loc.End
		return stmt, nil
	}

	// Optional FULL modifier.
	if p.cur.Kind == kwFULL {
		stmt.Full = true
		p.advance()
	}

	// Target table/function name.
	if isIdentifierToken(p.cur.Kind) || p.cur.Kind == tokString {
		target, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Target = target
	}

	// Optional partition list: describe_stmt's opt_partition_names
	// (DorisParser.g4:951 specifiedPartition, DorisLegacyParser.g4:2294).
	// Both require the table name first, so a partition clause without one is
	// left to the trailing-token check.
	if stmt.Target != nil {
		switch p.cur.Kind {
		case kwTEMPORARY, kwPARTITION, kwPARTITIONS:
			parts, temporary, err := p.parseDescribePartition()
			if err != nil {
				return nil, err
			}
			stmt.Partition, stmt.TempPartition = parts, temporary
		}
	}

	// Optional ALL VERBOSE suffix, which is an alternative to the partition
	// list rather than a tail on it (describe_stmt, DorisLegacyParser.g4:1629;
	// DorisParser.g4:950-951 describes the table or ALL, not both).
	if stmt.Partition == nil && !stmt.TempPartition && p.cur.Kind == kwALL {
		p.advance() // consume ALL
		if p.cur.Kind == kwVERBOSE {
			p.advance() // consume VERBOSE
			stmt.AllVerbose = true
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseDescribeFunction reads the rest of
//
//	FUNCTION name '(' propertyItemList? ')' tableAlias
//
// (DorisParser.g4:948; the CUP parser's table_valued_function_ref,
// DorisLegacyParser.g4:2240, takes the same shape, and its
// opt_key_value_map_in_paren is the same key=value list spelled with string
// literals). cur is FUNCTION.
func (p *Parser) parseDescribeFunction(stmt *ast.DescribeStmt) error {
	stmt.Function = true
	p.advance() // consume FUNCTION
	name, loc, err := p.parseIdentifier()
	if err != nil {
		return err
	}
	stmt.Target = &ast.ObjectName{Parts: []string{name}, Loc: loc}
	if _, err := p.expect(int('(')); err != nil {
		return err
	}
	if p.cur.Kind != int(')') { // the argument list is optional
		if stmt.Args, err = p.parsePropertyItems(p.parsePropertyItem); err != nil {
			return err
		}
	} else {
		p.advance() // consume ')'
	}
	// tableAlias : (AS? strictIdentifier identifierList?)? (:1393). The name
	// is required after AS, so a missing one is an error rather than an empty
	// alias.
	switch {
	case p.cur.Kind == kwAS:
		p.advance() // consume AS
		name, _, err := p.parseAliasIdentifier()
		if err != nil {
			return err
		}
		stmt.TableAlias = name
	case p.isAliasIdentToken():
		name, _, err := p.parseAliasIdentifier()
		if err != nil {
			return err
		}
		stmt.TableAlias = name
	}
	if stmt.TableAlias != "" && p.cur.Kind == int('(') {
		p.advance() // consume '('
		columns, err := p.parseIdentifierList()
		if err != nil {
			return err
		}
		stmt.TableColumns = columns
		if _, err := p.expect(int(')')); err != nil {
			return err
		}
	}
	return nil
}

// parseDescribePartition reads describe_stmt's opt_partition_names:
//
//	TEMPORARY? PARTITION (identifier | identifierList)
//	| TEMPORARY? PARTITIONS identifierList
//
// (DorisParser.g4:1627 specifiedPartition, which unsupportedDescribeStatement
// uses; the CUP parser's partition_names, DorisLegacyParser.g4:2311, also
// lists PARTITION and PARTITIONS with '*', which specifiedPartition does not
// have, and the engine refuses DESCRIBE t PARTITION(*).)
func (p *Parser) parseDescribePartition() ([]string, bool, error) {
	temporary := false
	if p.cur.Kind == kwTEMPORARY {
		temporary = true
		p.advance()
	}
	switch p.cur.Kind {
	case kwPARTITION:
		p.advance()
		if p.cur.Kind != int('(') { // PARTITION identifier
			name, _, err := p.parseIdentifier()
			return []string{name}, temporary, err
		}
	case kwPARTITIONS:
		p.advance()
	default:
		return nil, temporary, p.syntaxErrorAtCur()
	}
	if _, err := p.expect(int('(')); err != nil {
		return nil, temporary, err
	}
	var parts []string
	for {
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, temporary, err
		}
		parts = append(parts, name)
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}
	if _, err := p.expect(int(')')); err != nil {
		return nil, temporary, err
	}
	return parts, temporary, nil
}

// describeExplains reports whether what follows DESCRIBE/DESC is an explain
// command — an optional plan type, level and PROCESS, then the planned query
// or INSERT (DorisParser.g4:1061 explain, :52 statementBase) — rather than the
// table or function name describe_stmt names. It leaves the parser where it
// found it, so a bare word keeps describing a table.
func (p *Parser) describeExplains() bool {
	saved := p.save()
	defer p.restore(saved)

	// planType?
	switch p.cur.Kind {
	case kwPARSED, kwANALYZED, kwREWRITTEN, kwOPTIMIZED, kwPHYSICAL, kwSHAPE, kwMEMO:
		p.advance()
	}
	// level?
	switch p.cur.Kind {
	case kwVERBOSE, kwTREE, kwGRAPH, kwPLAN, kwDUMP:
		p.advance()
	}
	// PROCESS?
	if p.cur.Kind == kwPROCESS {
		p.advance()
	}

	switch p.cur.Kind {
	case kwSELECT, kwWITH, kwVALUES, kwINSERT, int('('):
		return true
	}
	return false
}

// ---------------------------------------------------------------------------
// EXPLAIN
// ---------------------------------------------------------------------------

// parseExplain parses EXPLAIN statements.
//
//	EXPLAIN [VERBOSE|GRAPH|PARSED|ANALYZED|REWRITTEN|PLAN|PHYSICAL|
//	         MEMO|SHAPE|DUMP|OPTIMIZED|PLAN PROCESS] query
//
// On entry, EXPLAIN has NOT yet been consumed.
func (p *Parser) parseExplain() (ast.Node, error) {
	startLoc := p.cur.Loc
	p.advance() // consume EXPLAIN

	stmt := &ast.ExplainStmt{Loc: startLoc}
	return p.parseExplainBody(stmt)
}

// parseExplainBody reads the modifiers and the explained statement of an
// explain command — EXPLAIN, DESCRIBE or DESC alike (DorisParser.g4:1061
// explain, :1067 explainCommand).
func (p *Parser) parseExplainBody(stmt *ast.ExplainStmt) (ast.Node, error) {
	// explain : explainCommand planType? level=(VERBOSE|TREE|GRAPH|PLAN|DUMP)?
	// PROCESS? (DorisParser.g4:1061), so the options are read in that order
	// and Type keeps the ones stated, space-separated — the same sequence
	// describeExplains looks ahead through.
	var options []string
	switch p.cur.Kind {
	case kwPARSED:
		options = append(options, "PARSED")
		p.advance()
	case kwANALYZED:
		options = append(options, "ANALYZED")
		p.advance()
	case kwREWRITTEN:
		options = append(options, "REWRITTEN")
		p.advance()
	case kwOPTIMIZED:
		options = append(options, "OPTIMIZED")
		p.advance()
	case kwPHYSICAL:
		options = append(options, "PHYSICAL")
		p.advance()
	case kwSHAPE:
		options = append(options, "SHAPE")
		p.advance()
	case kwMEMO:
		options = append(options, "MEMO")
		p.advance()
	}
	switch p.cur.Kind {
	case kwVERBOSE, kwTREE, kwGRAPH, kwDUMP:
		options = append(options, strings.ToUpper(p.advance().Str))
	case kwPLAN:
		options = append(options, "PLAN")
		p.advance()
	}
	// PROCESS is read once, after the level: the grammar gives it one slot,
	// so EXPLAIN PLAN PROCESS PROCESS is not a statement.
	if p.cur.Kind == kwPROCESS {
		if len(options) > 0 && options[len(options)-1] == "PLAN" {
			options[len(options)-1] = "PLAN PROCESS"
		} else {
			options = append(options, "PROCESS")
		}
		p.advance()
	}
	stmt.Type = strings.Join(options, " ")

	// Doris states the INSERT head: INSERT (INTO | OVERWRITE TABLE)
	// (DorisParser.g4:139, DorisLegacyParser.g4:1777 insert_stmt). parseInsert
	// also takes a bare INSERT table, a gap kept for the INSERT statement
	// itself; an explain command asks for the head the grammar names.
	if p.cur.Kind == kwINSERT {
		switch next := p.peekNext(); next.Kind {
		case kwINTO, kwOVERWRITE:
		default:
			return nil, syntaxErrorAt(next)
		}
	}

	// Parse the explained query. The RawQuery fallback is best-effort
	// recovery only: on the strict path it would swallow the nested failure
	// (EXPLAIN SELECT * FROM would parse clean with the garbage packed into
	// raw text), so strict mode propagates the error instead.
	if p.cur.Kind != tokEOF {
		query, err := p.parseStmt()
		if err == nil && query != nil {
			stmt.Query = query
		} else if p.strictTrailing {
			if err == nil {
				err = p.syntaxErrorAtCur()
			}
			return nil, err
		} else {
			// Fallback: wrap remaining tokens as a RawQuery.
			stmt.Query = &ast.RawQuery{RawText: p.collectRemainingRaw(), Loc: p.cur.Loc}
		}
	}

	// A bare EXPLAIN (or EXPLAIN followed only by a modifier or trivia) has
	// no explained statement; the engine requires one, so strict mode does
	// not accept the empty form.
	if p.strictTrailing && stmt.Query == nil {
		return nil, p.syntaxErrorAtCur()
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// USE
// ---------------------------------------------------------------------------

// parseUse parses a USE statement.
//
//	USE db_name
//	USE catalog_name.db_name
//	USE db_name@cluster_name
//	USE @cluster_name
//
// The last form is Doris' useCloudCluster with no database
// (DorisParser.g4:930), which switches the compute cluster alone.
//
// On entry, USE has NOT yet been consumed.
func (p *Parser) parseUse() (ast.Node, error) {
	startLoc := p.cur.Loc
	p.advance() // consume USE

	stmt := &ast.UseStmt{Loc: startLoc}

	// USE @cluster_name names the cluster and no database.
	if p.cur.Kind == int('@') {
		p.advance() // consume '@'
		cluster, clusterLoc, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Cluster = cluster
		stmt.Loc.End = clusterLoc.End
		return stmt, nil
	}

	// Parse the name: may be plain, catalog.db, or db@cluster.
	first, _, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	// Check for catalog.db form.
	if p.cur.Kind == int('.') {
		p.advance() // consume '.'
		db, _, err := p.parseIdentifierQualified()
		if err != nil {
			return nil, err
		}
		stmt.Catalog = first
		stmt.Database = db
	} else {
		stmt.Database = first
	}

	// Check for db@cluster form.
	if p.cur.Kind == int('@') {
		p.advance() // consume '@'
		cluster, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Cluster = cluster
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// SET (generic variable assignment)
// ---------------------------------------------------------------------------

// parseGenericSet parses #setOptions and #setTransaction
// (DorisParser.g4:882-890), as the CUP parser's set_stmt also takes them:
//
//	SET option (COMMA option)*
//	SET (GLOBAL | LOCAL | SESSION)? TRANSACTION characteristics
//
// SET has been consumed. A statement of one NAMES or CHARSET option keeps
// its Type; any other is a VARIABLE statement whose items may mix options.
func (p *Parser) parseGenericSet(startLoc ast.Loc) (ast.Node, error) {
	if p.cur.Kind == kwTRANSACTION {
		return p.parseSetTransaction(startLoc, "")
	}
	if (p.cur.Kind == kwGLOBAL || p.cur.Kind == kwSESSION || p.cur.Kind == kwLOCAL) &&
		p.peekNext().Kind == kwTRANSACTION {
		// LOCAL is a synonym for SESSION, matching the assignment scopes.
		scope := "SESSION"
		if p.cur.Kind == kwGLOBAL {
			scope = "GLOBAL"
		}
		p.advance() // consume the scope keyword
		return p.parseSetTransaction(startLoc, scope)
	}

	stmt := &ast.SetStmt{Loc: startLoc, Type: "VARIABLE"}
	var kind string
	for {
		item, k, err := p.parseSetOption()
		if err != nil {
			return nil, err
		}
		kind = k
		stmt.Items = append(stmt.Items, item)
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}
	// One NAMES or CHARSET option keeps its form's Type; a variable named
	// so (SET `names` = 1, SET @names = 1) stays a VARIABLE statement.
	if len(stmt.Items) == 1 {
		stmt.Type = kind
	}
	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseSetOption reads one optionWithType or optionWithoutType
// (DorisParser.g4:893-913):
//
//	(GLOBAL | LOCAL | SESSION) identifier EQ setValue
//	| (CHAR SET | CHARSET) (identifierOrText | DEFAULT)
//	| NAMES (identifierOrText | DEFAULT) (COLLATE identifierOrText)?
//	| LDAP_ADMIN_PASSWORD EQ (STRING_LITERAL | PASSWORD LEFT_PAREN STRING_LITERAL RIGHT_PAREN)
//	| (DOUBLEATSIGN ((GLOBAL | LOCAL | SESSION) DOT)?)? identifier EQ setValue
//	| ATSIGN identifier EQ expression
//	setValue : expression | DEFAULT | ON | ALL
//
// A scope word followed by '=' names the variable itself (SET GLOBAL = 1).
func (p *Parser) parseSetOption() (*ast.SetItem, string, error) {
	item := &ast.SetItem{Loc: p.cur.Loc}
	kind := "VARIABLE"
	switch {
	case p.cur.Kind == kwNAMES:
		p.advance() // consume NAMES
		item.Name = "names"
		kind = "NAMES"
		// Doris refuses NAMES = expression, which both grammars list, and a
		// DEFAULT after the charset (testdata/compat.tsv: SET NAMES = 1,
		// SET NAMES a DEFAULT).
		if p.cur.Kind == kwDEFAULT {
			p.advance()
			item.Raw = "DEFAULT"
		} else {
			val, _, err := p.parseIdentifierOrString()
			if err != nil {
				return nil, "", err
			}
			item.Raw = val
		}
		if p.cur.Kind == kwCOLLATE {
			p.advance()
			collation, _, err := p.parseIdentifierOrString()
			if err != nil {
				return nil, "", err
			}
			item.Raw += " COLLATE " + collation
		}

	case p.cur.Kind == kwCHARSET || (p.cur.Kind == kwCHAR && p.peekNext().Kind == kwSET):
		if p.advance().Kind == kwCHAR {
			p.advance() // consume SET
		}
		item.Name = "charset"
		kind = "CHARSET"
		if p.cur.Kind == kwDEFAULT {
			p.advance()
			item.Raw = "DEFAULT"
		} else {
			val, _, err := p.parseIdentifierOrString()
			if err != nil {
				return nil, "", err
			}
			item.Raw = val
		}

	case p.cur.Kind == kwLDAP_ADMIN_PASSWORD:
		p.advance()
		item.Name = "ldap_admin_password"
		if _, err := p.expect(int('=')); err != nil {
			return nil, "", err
		}
		if p.cur.Kind == kwPASSWORD {
			p.advance()
			if _, err := p.expect(int('(')); err != nil {
				return nil, "", err
			}
			if p.cur.Kind != tokString {
				return nil, "", p.syntaxErrorAtCur()
			}
			item.Raw = p.advance().Str
			if _, err := p.expect(int(')')); err != nil {
				return nil, "", err
			}
		} else if p.cur.Kind == tokString {
			item.Raw = p.advance().Str
		} else {
			return nil, "", p.syntaxErrorAtCur()
		}

	case p.cur.Kind == int('@'):
		p.advance() // consume '@'
		item.UserVariable = true
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, "", err
		}
		item.Name = name
		if _, err := p.expect(int('=')); err != nil {
			return nil, "", err
		}
		expr, err := p.parseExpr()
		if err != nil {
			return nil, "", err
		}
		item.Value = expr

	default:
		if p.cur.Kind == tokDoubleAt {
			p.advance() // consume @@
			if (p.cur.Kind == kwGLOBAL || p.cur.Kind == kwSESSION || p.cur.Kind == kwLOCAL) &&
				p.peekNext().Kind == int('.') {
				item.Scope = strings.ToUpper(p.cur.Str)
				p.advance() // consume the scope
				p.advance() // consume '.'
			}
		} else if (p.cur.Kind == kwGLOBAL || p.cur.Kind == kwSESSION || p.cur.Kind == kwLOCAL) &&
			p.peekNext().Kind != int('=') {
			item.Scope = strings.ToUpper(p.advance().Str)
		}
		var name string
		var err error
		if p.cur.Kind == kwGLOBAL || p.cur.Kind == kwSESSION || p.cur.Kind == kwLOCAL {
			// A scope word followed by '=' is the variable's name.
			name = p.advance().Str
		} else if name, _, err = p.parseIdentifier(); err != nil {
			return nil, "", err
		}
		item.Name = name
		// EQ: Doris refuses the CUP parser's := (SET a := DEFAULT).
		if _, err := p.expect(int('=')); err != nil {
			return nil, "", err
		}
		switch p.cur.Kind {
		case kwDEFAULT, kwON, kwALL:
			item.Raw = strings.ToUpper(p.advance().Str)
		default:
			expr, err := p.parseExpr()
			switch {
			case err == nil:
				item.Value = expr
			case p.strictTrailing:
				return nil, "", err
			default:
				// Best effort keeps the malformed value as raw text.
				item.Raw = p.collectUntilCommaOrEOF()
			}
		}
	}
	item.Loc.End = p.prev.Loc.End
	return item, kind, nil
}

// parseSetProperty parses #setUserProperties (DorisParser.g4:885) as the CUP
// parser's set_stmt also takes it:
//
//	SET PROPERTY (FOR identifierOrText)? STRING_LITERAL '=' (STRING_LITERAL | NULL) (',' ...)*
//
// SET has been consumed; cur is PROPERTY.
func (p *Parser) parseSetProperty(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume PROPERTY
	stmt := &ast.SetStmt{Loc: startLoc, Type: "PROPERTY"}
	if p.cur.Kind == kwFOR {
		p.advance() // consume FOR
		user, _, err := p.parseIdentifierOrString()
		if err != nil {
			return nil, err
		}
		stmt.PropertyUser = user
	}
	for {
		start := p.cur.Loc
		if p.cur.Kind != tokString {
			return nil, p.syntaxErrorAtCur()
		}
		key := p.advance().Str
		if _, err := p.expect(int('=')); err != nil {
			return nil, err
		}
		var val string
		switch p.cur.Kind {
		case tokString:
			val = p.advance().Str
		case kwNULL:
			p.advance()
			val = "NULL"
		default:
			return nil, p.syntaxErrorAtCur()
		}
		stmt.Properties = append(stmt.Properties, &ast.Property{Key: key, Value: val, Loc: ast.Loc{Start: start.Start, End: p.prev.Loc.End}})
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}
	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseSetNames parses: SET NAMES 'charset' [COLLATE 'collation']
// On entry cur == kwNAMES; startLoc is the SET Loc.
func (p *Parser) parseSetNames(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume NAMES
	stmt := &ast.SetStmt{Loc: startLoc, Type: "NAMES"}

	item := &ast.SetItem{Name: "names"}
	if p.cur.Kind == kwDEFAULT {
		// SET NAMES DEFAULT — DEFAULT lexes as a reserved keyword, so the
		// identifier/string gate below would misread the engine-valid form
		// as a missing charset (container-verified accept on both engines).
		tok := p.advance()
		item.Raw = "DEFAULT"
		item.Loc = tok.Loc
	} else if p.cur.Kind == tokString || isIdentifierToken(p.cur.Kind) {
		val, loc, err := p.parseIdentifierOrString()
		if err == nil {
			item.Raw = val
			item.Loc = loc
		}
	} else if p.strictTrailing {
		// SET NAMES with no charset ends at EOF, invisible to the
		// leftover-token check; the operand is required (engine-verified on
		// both engines).
		return nil, p.syntaxErrorAtCur()
	}

	// Optional COLLATE
	if p.cur.Kind == kwCOLLATE {
		p.advance()
		collation, _, _ := p.parseIdentifierOrString()
		if collation != "" {
			item.Raw += " COLLATE " + collation
		} else if p.strictTrailing {
			// A dangling COLLATE needs its collation name.
			return nil, p.syntaxErrorAtCur()
		}
	}

	stmt.Items = []*ast.SetItem{item}
	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseSetCharset parses: SET CHARSET 'charset'
// On entry cur == kwCHARSET; startLoc is the SET Loc.
func (p *Parser) parseSetCharset(startLoc ast.Loc) (ast.Node, error) {
	p.advance() // consume CHARSET
	stmt := &ast.SetStmt{Loc: startLoc, Type: "CHARSET"}

	item := &ast.SetItem{Name: "charset"}
	if p.cur.Kind == tokString || isIdentifierToken(p.cur.Kind) {
		val, loc, err := p.parseIdentifierOrString()
		if err == nil {
			item.Raw = val
			item.Loc = loc
		}
	} else if p.strictTrailing {
		// SET CHARSET with no charset name — the operand is required.
		return nil, p.syntaxErrorAtCur()
	}
	stmt.Items = []*ast.SetItem{item}
	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseSetTransaction parses: SET TRANSACTION { READ ONLY | READ WRITE | ISOLATION LEVEL ... }
// On entry cur == kwTRANSACTION; startLoc is the SET Loc.
// scope carries the optional GLOBAL/SESSION qualifier consumed by the
// caller, so SET SESSION TRANSACTION ... stays distinguishable from the
// unqualified form.
func (p *Parser) parseSetTransaction(startLoc ast.Loc, scope string) (ast.Node, error) {
	p.advance() // consume TRANSACTION
	stmt := &ast.SetStmt{Loc: startLoc, Type: "TRANSACTION"}

	// transaction_characteristic (, transaction_characteristic)*:
	//   ISOLATION LEVEL {READ UNCOMMITTED | READ COMMITTED | REPEATABLE READ
	//   | SERIALIZABLE} | READ ONLY | READ WRITE
	// Parsed for real rather than raw-captured to EOF — the raw capture ate
	// everything, so the strict trailing-token check could never see junk
	// like SET TRANSACTION BOGUS.
	parts, partial, err := p.parseTransactionCharacteristics()
	if err != nil {
		if p.strictTrailing {
			return nil, err
		}
		// Best-effort keeps its recovery contract: completion-style partial
		// input (SET TRANSACTION ISOLATION LEVEL READ) still yields the
		// statement. The interrupted characteristic and the unparsed
		// remainder stay one fragment — joining them as separate parts would
		// invent a comma the user never typed.
		if frag := strings.TrimSpace(partial + " " + p.collectRemainingRaw()); frag != "" {
			parts = append(parts, frag)
		}
	}

	item := &ast.SetItem{Name: "transaction", Scope: scope, Raw: strings.Join(parts, ", ")}
	stmt.Items = []*ast.SetItem{item}
	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// parseTransactionCharacteristics parses the characteristic list of
// SET TRANSACTION, returning whatever parsed cleanly plus the error that
// stopped it, so the caller can choose strict propagation or best-effort
// recovery.
func (p *Parser) parseTransactionCharacteristics() ([]string, string, error) {
	var parts []string
	for {
		iterStart := p.cur.Loc.Start
		switch p.cur.Kind {
		case kwISOLATION:
			p.advance() // consume ISOLATION
			if _, err := p.expect(kwLEVEL); err != nil {
				return parts, p.consumedCharacteristicText(iterStart), err
			}
			switch p.cur.Kind {
			case kwREAD:
				p.advance() // consume READ
				switch p.cur.Kind {
				case kwUNCOMMITTED, kwCOMMITTED:
					parts = append(parts, "ISOLATION LEVEL READ "+strings.ToUpper(p.advance().Str))
				default:
					return parts, p.consumedCharacteristicText(iterStart), p.syntaxErrorAtCur()
				}
			case kwREPEATABLE:
				p.advance() // consume REPEATABLE
				if _, err := p.expect(kwREAD); err != nil {
					return parts, p.consumedCharacteristicText(iterStart), err
				}
				parts = append(parts, "ISOLATION LEVEL REPEATABLE READ")
			case kwSERIALIZABLE:
				p.advance()
				parts = append(parts, "ISOLATION LEVEL SERIALIZABLE")
			default:
				return parts, p.consumedCharacteristicText(iterStart), p.syntaxErrorAtCur()
			}
		case kwREAD:
			p.advance() // consume READ
			switch p.cur.Kind {
			case kwONLY:
				p.advance()
				parts = append(parts, "READ ONLY")
			case kwWRITE:
				p.advance()
				parts = append(parts, "READ WRITE")
			default:
				return parts, p.consumedCharacteristicText(iterStart), p.syntaxErrorAtCur()
			}
		default:
			return parts, p.consumedCharacteristicText(iterStart), p.syntaxErrorAtCur()
		}
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}
	return parts, "", nil
}

// consumedCharacteristicText returns the text of the characteristic a parse
// error interrupted mid-way, so best-effort recovery does not lose it:
// ParseBestEffort("SET TRANSACTION ISOLATION LEVEL READ") keeps
// "ISOLATION LEVEL READ" in the item's Raw even though the tokens were
// already consumed when the error surfaced.
func (p *Parser) consumedCharacteristicText(start int) string {
	return strings.TrimSpace(p.input[start-p.baseOffset : p.cur.Loc.Start-p.baseOffset])
}

// parseSetItem parses a single SET assignment:
//
//	[GLOBAL|SESSION|LOCAL] var_name = expr
//	@@[GLOBAL.|SESSION.]var_name = expr
func (p *Parser) parseSetItem() (*ast.SetItem, error) {
	item := &ast.SetItem{Loc: p.cur.Loc}

	// Optional @@ prefix.
	if p.cur.Kind == tokDoubleAt {
		p.advance() // consume @@
		// Optional scope prefix: GLOBAL. or SESSION.
		if (p.cur.Kind == kwGLOBAL || p.cur.Kind == kwSESSION) && p.peekNext().Kind == int('.') {
			item.Scope = strings.ToUpper(p.cur.Str)
			p.advance() // consume GLOBAL/SESSION
			p.advance() // consume '.'
		}
		name, _, err := p.parseIdentifierQualified()
		if err != nil {
			return nil, err
		}
		item.Name = name
	} else {
		// Optional scope keyword.
		switch p.cur.Kind {
		case kwGLOBAL:
			item.Scope = "GLOBAL"
			p.advance()
		case kwSESSION, kwLOCAL:
			item.Scope = "SESSION"
			p.advance()
		}
		// Variable name.
		name, _, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		item.Name = name
	}

	// Consume = or :=
	hasAssign := false
	if p.cur.Kind == int('=') || p.cur.Kind == tokAssign {
		hasAssign = true
		p.advance()
	}

	// Parse the value expression.
	if p.cur.Kind != tokEOF && p.cur.Kind != int(',') {
		if p.strictTrailing && !hasAssign {
			// SET x 1 — the engine requires the assignment separator.
			return nil, p.syntaxErrorAtCur()
		}
		expr, err := p.parseExpr()
		if err == nil {
			item.Value = expr
		} else if p.strictTrailing {
			// The raw fallback would consume to the comma or EOF and hide the
			// malformed expression from the strict trailing-token check
			// (SET x = ( parsed clean). Strict mode propagates instead.
			return nil, err
		} else {
			// Fallback to raw.
			item.Raw = p.collectUntilCommaOrEOF()
		}
	} else if p.strictTrailing {
		// SET x / SET x = / SET x, y — the assignment value is required, and
		// the bare-name forms used to reach EOF or the comma before the
		// strict check could see anything wrong.
		return nil, p.syntaxErrorAtCur()
	}

	item.Loc.End = p.prev.Loc.End
	return item, nil
}

// collectUntilCommaOrEOF reads raw tokens until ',' or EOF.
func (p *Parser) collectUntilCommaOrEOF() string {
	var parts []string
	for p.cur.Kind != tokEOF && p.cur.Kind != int(',') {
		if p.cur.Str != "" {
			parts = append(parts, p.cur.Str)
		} else if p.cur.Kind > 0 && p.cur.Kind < 128 {
			parts = append(parts, string(rune(p.cur.Kind)))
		}
		p.advance()
	}
	return strings.Join(parts, " ")
}

// ---------------------------------------------------------------------------
// UNSET
// ---------------------------------------------------------------------------

// parseGenericUnset parses generic UNSET statements that are not handled by
// the dedicated UNSET DEFAULT STORAGE VAULT path:
//
//	UNSET [GLOBAL|SESSION] VARIABLE name [, name ...]
//	UNSET [GLOBAL|SESSION] VARIABLE ALL
//	UNSET [GLOBAL|SESSION] VARIABLE *
//
// On entry, UNSET has already been consumed; startLoc is the UNSET token Loc.
func (p *Parser) parseGenericUnset(startLoc ast.Loc) (ast.Node, error) {
	stmt := &ast.UnsetStmt{Loc: startLoc, Type: "VARIABLE"}

	// Optional scope.
	switch p.cur.Kind {
	case kwGLOBAL:
		stmt.Scope = "GLOBAL"
		p.advance()
	case kwSESSION, kwLOCAL:
		stmt.Scope = "SESSION"
		p.advance()
	}

	// Expect VARIABLE keyword.
	if p.cur.Kind == kwVARIABLE {
		p.advance()
	}

	// ALL or * — unset all variables.
	if p.cur.Kind == kwALL || p.cur.Kind == int('*') {
		stmt.All = true
		p.advance()
		stmt.Loc.End = p.prev.Loc.End
		return stmt, nil
	}

	// One or more variable names.
	for {
		name, _, err := p.parseIdentifier()
		if err != nil {
			break
		}
		stmt.Names = append(stmt.Names, name)
		if _, ok := p.match(int(',')); !ok {
			break
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}

// ---------------------------------------------------------------------------
// HELP
// ---------------------------------------------------------------------------

// parseHelp parses: HELP 'mask'
// On entry, HELP has NOT yet been consumed.
func (p *Parser) parseHelp() (ast.Node, error) {
	startLoc := p.cur.Loc
	p.advance() // consume HELP

	stmt := &ast.HelpStmt{Loc: startLoc}

	if p.cur.Kind == tokString || isIdentifierToken(p.cur.Kind) {
		mask, _, err := p.parseIdentifierOrString()
		if err == nil {
			stmt.Mask = mask
		}
	}

	stmt.Loc.End = p.prev.Loc.End
	return stmt, nil
}
