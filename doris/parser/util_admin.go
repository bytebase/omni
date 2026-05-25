package parser

import (
	"strings"

	"github.com/bytebase/omni/doris/ast"
)

// ---------------------------------------------------------------------------
// BACKUP SNAPSHOT
// ---------------------------------------------------------------------------

// parseBackup parses:
//
//	BACKUP SNAPSHOT [db.]label TO repo_name
//	    [ON (tbl [PARTITION(p1,...)], ...)]
//	    [PROPERTIES("key"="value", ...)]
//
// On entry, cur is SNAPSHOT (BACKUP was consumed by the dispatch).
func (p *Parser) parseBackup(startLoc ast.Loc) (ast.Node, error) {
	if _, err := p.expect(kwSNAPSHOT); err != nil {
		return nil, err
	}

	stmt := &ast.BackupStmt{}

	// [db.]label — parse as multipart identifier
	labelName, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Label = labelName.String()
	stmt.LabelParts = labelName.Parts
	endLoc := ast.NodeLoc(labelName)

	// TO repo_name
	if _, err := p.expect(kwTO); err != nil {
		return nil, err
	}
	repo, repoLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Repo = repo
	endLoc = repoLoc

	// Optional ON (table_list)
	if p.cur.Kind == kwON {
		p.advance() // consume ON
		tables, loc, err := p.parseBackupTableList()
		if err != nil {
			return nil, err
		}
		stmt.Tables = tables
		endLoc = loc
	}

	// Optional PROPERTIES(...)
	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
		if len(props) > 0 {
			endLoc = ast.NodeLoc(props[len(props)-1])
		}
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseBackupTableList parses the parenthesized table list in BACKUP/RESTORE ON clause.
// We consume '(' tbl [PARTITION(...)], tbl [PARTITION(...)] ..., ')'.
// Individual PARTITION sub-clauses are consumed but not stored (best-effort).
func (p *Parser) parseBackupTableList() ([]*ast.ObjectName, ast.Loc, error) {
	startTok, err := p.expect(int('('))
	if err != nil {
		return nil, ast.NoLoc(), err
	}
	endLoc := startTok.Loc

	var tables []*ast.ObjectName
	for p.cur.Kind != int(')') && p.cur.Kind != tokEOF {
		tbl, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, ast.NoLoc(), err
		}
		tables = append(tables, tbl)
		endLoc = ast.NodeLoc(tbl)

		// Optional PARTITION(p1, p2, ...) — consume and discard
		if p.cur.Kind == kwPARTITION {
			p.advance() // consume PARTITION
			if p.cur.Kind == int('(') {
				_, loc, err := p.consumeParenGroup()
				if err != nil {
					return nil, ast.NoLoc(), err
				}
				endLoc = loc
			}
		}

		// Comma separator
		if p.cur.Kind == int(',') {
			p.advance()
		}
	}

	closeTok, err := p.expect(int(')'))
	if err != nil {
		return nil, ast.NoLoc(), err
	}
	endLoc = closeTok.Loc

	return tables, endLoc, nil
}

// ---------------------------------------------------------------------------
// RESTORE SNAPSHOT
// ---------------------------------------------------------------------------

// parseRestore parses:
//
//	RESTORE SNAPSHOT [db.]label FROM repo_name
//	    [ON (tbl [PARTITION(p1,...)], ...)]
//	    [PROPERTIES("key"="value", ...)]
//
// On entry, cur is SNAPSHOT (RESTORE was consumed by the dispatch).
func (p *Parser) parseRestore(startLoc ast.Loc) (ast.Node, error) {
	if _, err := p.expect(kwSNAPSHOT); err != nil {
		return nil, err
	}

	stmt := &ast.RestoreStmt{}

	labelName, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Label = labelName.String()
	stmt.LabelParts = labelName.Parts
	endLoc := ast.NodeLoc(labelName)

	// FROM repo_name
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}
	repo, repoLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	stmt.Repo = repo
	endLoc = repoLoc

	// Optional ON (table_list)
	if p.cur.Kind == kwON {
		p.advance()
		tables, loc, err := p.parseBackupTableList()
		if err != nil {
			return nil, err
		}
		stmt.Tables = tables
		endLoc = loc
	}

	// Optional PROPERTIES(...)
	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
		if len(props) > 0 {
			endLoc = ast.NodeLoc(props[len(props)-1])
		}
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// ---------------------------------------------------------------------------
// KILL
// ---------------------------------------------------------------------------

// parseKill parses:
//
//	KILL [CONNECTION | QUERY] id
//
// On entry, cur is CONNECTION | QUERY | the id token (KILL was consumed by dispatch).
func (p *Parser) parseKill(startLoc ast.Loc) (ast.Node, error) {
	stmt := &ast.KillStmt{}
	endLoc := startLoc

	switch p.cur.Kind {
	case kwCONNECTION:
		stmt.Kind = "CONNECTION"
		endLoc = p.cur.Loc
		p.advance()
	case kwQUERY:
		stmt.Kind = "QUERY"
		endLoc = p.cur.Loc
		p.advance()
	}

	// Target: integer id or string query id
	switch p.cur.Kind {
	case tokInt:
		stmt.Target = p.cur.Str
		endLoc = p.cur.Loc
		p.advance()
	case tokString:
		stmt.Target = p.cur.Str
		endLoc = p.cur.Loc
		p.advance()
	default:
		// Best-effort: accept any token as a target id
		if p.cur.Kind != tokEOF {
			stmt.Target = p.cur.Str
			endLoc = p.cur.Loc
			p.advance()
		}
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// ---------------------------------------------------------------------------
// LOCK TABLES
// ---------------------------------------------------------------------------

// parseLockTables parses:
//
//	LOCK TABLES tbl [AS alias] {READ [LOCAL] | [LOW_PRIORITY] WRITE} [, ...]
//
// On entry, cur is TABLES (LOCK was consumed by dispatch).
func (p *Parser) parseLockTables(startLoc ast.Loc) (ast.Node, error) {
	if _, err := p.expect(kwTABLES); err != nil {
		return nil, err
	}

	stmt := &ast.LockTablesStmt{}
	endLoc := startLoc

	for p.cur.Kind != tokEOF {
		item, err := p.parseLockItem()
		if err != nil {
			return nil, err
		}
		stmt.Items = append(stmt.Items, item)
		endLoc = ast.NodeLoc(item)

		if p.cur.Kind == int(',') {
			p.advance()
		} else {
			break
		}
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseLockItem parses one lock entry: tbl [AS alias] {READ [LOCAL] | [LOW_PRIORITY] WRITE}.
func (p *Parser) parseLockItem() (*ast.LockItem, error) {
	tbl, err := p.parseMultipartIdentifier()
	if err != nil {
		return nil, err
	}

	item := &ast.LockItem{Table: tbl}
	endLoc := ast.NodeLoc(tbl)

	// Optional AS alias
	if p.cur.Kind == kwAS {
		p.advance()
		alias, aliasLoc, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		item.Alias = alias
		endLoc = aliasLoc
	}

	// Lock mode
	switch p.cur.Kind {
	case kwREAD:
		p.advance()
		endLoc = p.prev.Loc
		if p.cur.Kind == kwLOCAL {
			p.advance()
			endLoc = p.prev.Loc
			item.Mode = "READ LOCAL"
		} else {
			item.Mode = "READ"
		}
	case kwLOW_PRIORITY:
		p.advance()
		if _, err := p.expect(kwWRITE); err != nil {
			return nil, err
		}
		endLoc = p.prev.Loc
		item.Mode = "LOW_PRIORITY WRITE"
	case kwWRITE:
		p.advance()
		endLoc = p.prev.Loc
		item.Mode = "WRITE"
	default:
		// Best-effort: no lock mode token found; leave Mode empty
	}

	item.Loc = ast.NodeLoc(tbl).Merge(endLoc)
	return item, nil
}

// parseUnlockTables parses:
//
//	UNLOCK TABLES
//
// On entry, cur is TABLES (UNLOCK was consumed by dispatch).
func (p *Parser) parseUnlockTables(startLoc ast.Loc) (ast.Node, error) {
	endLoc := startLoc
	if p.cur.Kind == kwTABLES {
		endLoc = p.cur.Loc
		p.advance()
	}
	return &ast.UnlockTablesStmt{Loc: startLoc.Merge(endLoc)}, nil
}

// ---------------------------------------------------------------------------
// INSTALL / UNINSTALL PLUGIN
// ---------------------------------------------------------------------------

// parseInstallPlugin parses:
//
//	INSTALL PLUGIN FROM 'source'
//	INSTALL PLUGIN FROM SONAME 'library'
//	(with optional PROPERTIES(...))
//
// On entry, cur is PLUGIN (INSTALL was consumed by dispatch).
func (p *Parser) parseInstallPlugin(startLoc ast.Loc) (ast.Node, error) {
	if _, err := p.expect(kwPLUGIN); err != nil {
		return nil, err
	}
	if _, err := p.expect(kwFROM); err != nil {
		return nil, err
	}

	stmt := &ast.InstallPluginStmt{}
	endLoc := p.prev.Loc

	// Optional SONAME keyword
	if p.cur.Kind == kwSONAME {
		stmt.IsSoname = true
		p.advance()
	}

	// Source path / URL / library name — string literal
	src, srcLoc, err := p.parseIdentifierOrString()
	if err != nil {
		return nil, err
	}
	stmt.Source = src
	endLoc = srcLoc

	// Optional PROPERTIES(...)
	if p.cur.Kind == kwPROPERTIES {
		props, err := p.parseProperties()
		if err != nil {
			return nil, err
		}
		stmt.Properties = props
		if len(props) > 0 {
			endLoc = ast.NodeLoc(props[len(props)-1])
		}
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// parseUninstallPlugin parses:
//
//	UNINSTALL PLUGIN name
//
// On entry, cur is PLUGIN (UNINSTALL was consumed by dispatch).
func (p *Parser) parseUninstallPlugin(startLoc ast.Loc) (ast.Node, error) {
	if _, err := p.expect(kwPLUGIN); err != nil {
		return nil, err
	}

	name, nameLoc, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	return &ast.UninstallPluginStmt{
		Name: name,
		Loc:  startLoc.Merge(nameLoc),
	}, nil
}

// ---------------------------------------------------------------------------
// WARM UP
// ---------------------------------------------------------------------------

// parseWarmUp parses:
//
//	WARM UP CLUSTER cluster_name FROM ...
//	WARM UP COMPUTE GROUP cg WITH ...
//
// On entry, cur is UP (WARM was consumed by dispatch).
func (p *Parser) parseWarmUp(startLoc ast.Loc) (ast.Node, error) {
	if _, err := p.expect(kwUP); err != nil {
		return nil, err
	}

	stmt := &ast.WarmUpStmt{}
	endLoc := p.prev.Loc

	switch p.cur.Kind {
	case kwCLUSTER:
		stmt.Verb = "CLUSTER"
		p.advance()
	case kwCOMPUTE:
		p.advance() // consume COMPUTE
		stmt.Verb = "COMPUTE GROUP"
		if p.cur.Kind == kwGROUP {
			p.advance()
		}
	default:
		// Tolerate missing verb
		if p.cur.Kind != tokEOF {
			stmt.Verb = strings.ToUpper(p.cur.Str)
			p.advance()
		}
	}

	// Collect remainder as raw text
	raw, rawEnd := p.collectRawUntilEOF()
	stmt.Target = raw
	if rawEnd.IsValid() {
		endLoc = rawEnd
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// ---------------------------------------------------------------------------
// CLEAN
// ---------------------------------------------------------------------------

// parseClean parses:
//
//	CLEAN ALL PROFILE
//	CLEAN LABEL [label] [FROM db]
//	CLEAN QUERY STATS [FROM db] [ALL]
//	... and other variants
//
// On entry, cur is ALL | LABEL | QUERY | etc. (CLEAN was consumed by dispatch).
func (p *Parser) parseClean(startLoc ast.Loc) (ast.Node, error) {
	stmt := &ast.CleanStmt{}
	endLoc := startLoc

	switch p.cur.Kind {
	case kwALL:
		stmt.Verb = "ALL"
		p.advance()
		if p.cur.Kind == kwPROFILE {
			stmt.Verb = "ALL PROFILE"
			endLoc = p.cur.Loc
			p.advance()
		}
	case kwLABEL:
		stmt.Verb = "LABEL"
		endLoc = p.cur.Loc
		p.advance()
	case kwQUERY:
		stmt.Verb = "QUERY"
		p.advance()
		if p.cur.Kind == kwSTATS {
			stmt.Verb = "QUERY STATS"
			endLoc = p.cur.Loc
			p.advance()
		}
	default:
		if p.cur.Kind != tokEOF {
			stmt.Verb = strings.ToUpper(p.cur.Str)
			endLoc = p.cur.Loc
			p.advance()
		}
	}

	// Collect remaining raw text (FROM db, ALL, etc.)
	raw, rawEnd := p.collectRawUntilEOF()
	stmt.Target = raw
	if rawEnd.IsValid() {
		endLoc = rawEnd
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// ---------------------------------------------------------------------------
// CANCEL (generic)
// ---------------------------------------------------------------------------

// parseCancelGeneric parses non-MTMV CANCEL commands:
//
//	CANCEL LOAD [FROM db] WHERE ...
//	CANCEL EXPORT [FROM db] WHERE ...
//	CANCEL ALTER TABLE COLUMN|ROLLUP FROM table
//	CANCEL BACKUP [FROM db]
//	CANCEL RESTORE [FROM db]
//	CANCEL BUILD INDEX ON table
//
// On entry, cur is the keyword immediately after CANCEL (e.g. LOAD, EXPORT,
// ALTER, BACKUP, RESTORE, BUILD).  cancelTok.Loc is the CANCEL token's Loc.
func (p *Parser) parseCancelGeneric(cancelLoc ast.Loc) (ast.Node, error) {
	stmt := &ast.CancelStmt{}
	endLoc := cancelLoc

	switch p.cur.Kind {
	case kwLOAD:
		stmt.Verb = "LOAD"
		endLoc = p.cur.Loc
		p.advance()
	case kwEXPORT:
		stmt.Verb = "EXPORT"
		endLoc = p.cur.Loc
		p.advance()
	case kwBACKUP:
		stmt.Verb = "BACKUP"
		endLoc = p.cur.Loc
		p.advance()
	case kwRESTORE:
		stmt.Verb = "RESTORE"
		endLoc = p.cur.Loc
		p.advance()
	case kwALTER:
		p.advance() // consume ALTER
		stmt.Verb = "ALTER"
		if p.cur.Kind == kwTABLE {
			stmt.Verb = "ALTER TABLE"
			endLoc = p.cur.Loc
			p.advance()
		}
	case kwBUILD:
		p.advance() // consume BUILD
		stmt.Verb = "BUILD"
		if p.cur.Kind == kwINDEX {
			stmt.Verb = "BUILD INDEX"
			endLoc = p.cur.Loc
			p.advance()
		}
	default:
		if p.cur.Kind != tokEOF {
			stmt.Verb = strings.ToUpper(p.cur.Str)
			endLoc = p.cur.Loc
			p.advance()
		}
	}

	raw, rawEnd := p.collectRawUntilEOF()
	stmt.Target = raw
	if rawEnd.IsValid() {
		endLoc = rawEnd
	}

	stmt.Loc = cancelLoc.Merge(endLoc)
	return stmt, nil
}

// ---------------------------------------------------------------------------
// RECOVER
// ---------------------------------------------------------------------------

// parseRecover parses:
//
//	RECOVER DATABASE [name | dbid id] [AS new_name]
//	RECOVER TABLE [db.]name [tblid id] [AS new_name]
//	RECOVER PARTITION [name | pid id] FROM [db.]table [AS new_name]
//
// On entry, cur is DATABASE | TABLE | PARTITION (RECOVER was consumed by dispatch).
func (p *Parser) parseRecover(startLoc ast.Loc) (ast.Node, error) {
	stmt := &ast.RecoverStmt{}
	endLoc := startLoc

	switch p.cur.Kind {
	case kwDATABASE, kwSCHEMA:
		stmt.Verb = "DATABASE"
		p.advance()
	case kwTABLE:
		stmt.Verb = "TABLE"
		p.advance()
	case kwPARTITION:
		stmt.Verb = "PARTITION"
		p.advance()
	default:
		return nil, p.syntaxErrorAtCur()
	}

	// Parse object name (optional for DATABASE when only dbid is given,
	// but in practice the legacy corpus always has a name or id).
	// We accept identifier-or-nothing to be tolerant.
	if isIdentifierToken(p.cur.Kind) {
		name, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.Name = name
		endLoc = ast.NodeLoc(name)
	}

	// Optional numeric id — the legacy corpus uses bare integers: 12345
	// or uses keywords dbid/tblid/pid before the integer.
	switch p.cur.Kind {
	case kwDORIS_INTERNAL_TABLE_ID:
		// legacy: just consume as id keyword
		p.advance()
		if p.cur.Kind == tokInt {
			stmt.ID = p.cur.Str
			endLoc = p.cur.Loc
			p.advance()
		}
	case tokInt:
		stmt.ID = p.cur.Str
		endLoc = p.cur.Loc
		p.advance()
	default:
		// Check for identifier keywords used as id labels (dbid, tblid, pid)
		if p.cur.Kind == tokIdent {
			lower := strings.ToLower(p.cur.Str)
			if lower == "dbid" || lower == "tblid" || lower == "pid" {
				p.advance() // consume the label
				if p.cur.Kind == tokInt {
					stmt.ID = p.cur.Str
					endLoc = p.cur.Loc
					p.advance()
				}
			}
		}
	}

	// Optional AS new_name (may appear before or after FROM for PARTITION)
	if p.cur.Kind == kwAS {
		p.advance()
		newName, newNameLoc, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.NewName = newName
		endLoc = newNameLoc
	}

	// For PARTITION: optional FROM [db.]table
	if stmt.Verb == "PARTITION" && p.cur.Kind == kwFROM {
		p.advance() // consume FROM
		fromTable, err := p.parseMultipartIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.FromTable = fromTable
		endLoc = ast.NodeLoc(fromTable)
	}

	// Second AS new_name position (when AS follows FROM for TABLE/DATABASE)
	if stmt.NewName == "" && p.cur.Kind == kwAS {
		p.advance()
		newName, newNameLoc, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		stmt.NewName = newName
		endLoc = newNameLoc
	}

	stmt.Loc = startLoc.Merge(endLoc)
	return stmt, nil
}

// ---------------------------------------------------------------------------
// Helper: collectRawUntilEOF
// ---------------------------------------------------------------------------

// collectRawUntilEOF consumes all remaining tokens up to EOF (or ';') and
// returns them as a single space-joined string plus the Loc of the last token.
// Used by best-effort parsers (WARM UP, CLEAN, CANCEL) to capture unparsed
// trailing text.
func (p *Parser) collectRawUntilEOF() (string, ast.Loc) {
	var parts []string
	endLoc := ast.NoLoc()

	for p.cur.Kind != tokEOF && p.cur.Kind != int(';') {
		str := p.cur.Str
		if str == "" {
			str = TokenName(p.cur.Kind)
		}
		parts = append(parts, str)
		endLoc = p.cur.Loc
		p.advance()
	}

	return strings.Join(parts, " "), endLoc
}
