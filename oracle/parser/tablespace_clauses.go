package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// This file reads the tablespace statements: CREATE TABLESPACE (permanent,
// temporary and undo), CREATE TABLESPACE SET, ALTER TABLESPACE [SET] and
// DROP TABLESPACE [SET], with the clauses they share.
//
// BNF: oracle/parser/bnf/CREATE-TABLESPACE.bnf, CREATE-TABLESPACE-SET.bnf,
// ALTER-TABLESPACE.bnf, ALTER-TABLESPACE-SET.bnf, DROP-TABLESPACE.bnf,
// DROP-TABLESPACE-SET.bnf
// (https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-TABLESPACE.html
// and the pages beside it).
//
// Where the pages and the pinned Oracle 23.26.3 disagree, the parser follows
// Oracle, as its compat.tsv rows and probes on the server show:
//
//   - CREATE reads its options in any order, each at most once (a second
//     LOGGING is ORA-25119, a second ONLINE ORA-02198, a second DATAFILE
//     ORA-02197), and each tablespace kind takes its own set of them
//     (ORA-02180, ORA-25139, ORA-30024). ALTER takes exactly one clause.
//   - IF NOT EXISTS comes before the name: CREATE TABLESPACE IF NOT EXISTS a.
//   - LOCAL TEMPORARY TABLESPACE takes FOR { ALL | LEAF } before the name.
//   - logging_clause is LOGGING or NOLOGGING; FILESYSTEM_LIKE_LOGGING is
//     refused here.
//   - CREATE takes LOST WRITE PROTECTION with no ENABLE; ALTER takes ENABLE,
//     SUSPEND or REMOVE (DISABLE is ORA-02142).
//   - ALTER's tablespace_group_clause is TABLESPACE GROUP, as CREATE's.
//   - DEFAULT TABLE takes only COMPRESS or NOCOMPRESS; ROW STORE and COLUMN
//     STORE compression stand without TABLE.
//   - A permanent tablespace also takes the undocumented AUTOALLOCATE and
//     UNIFORM [ SIZE size_clause ], and ALTER the undocumented OFFLINE FOR
//     RECOVER.
//
// Oracle answers an error here with its SQL parser's code (ORA-00905,
// ORA-02000, ORA-03049, ...) or with one of the tablespace statements' own
// (ORA-02142 invalid ALTER TABLESPACE option, ORA-02180 invalid option for
// CREATE TABLESPACE, ORA-02236 invalid file name, ...); either is a syntax
// error, at the end of the text too.

// tsAtEnd reports whether the statement's text has ended.
func (p *Parser) tsAtEnd() bool {
	return p.cur.Type == ';' || p.cur.Type == tokEOF
}

// tsEndsBeforeOperand reports whether the text ends where Oracle answers a
// missing operand with a code that is not a syntax error, so the compat
// judge takes the statement: ORA-25121 for MINIMUM EXTENT's size, ORA-25147
// for UNIFORM SIZE's, ORA-02236 at byte 0 for DROP DATAFILE's file,
// ORA-02580 for a tablespace set's ADD DATAFILE, ORA-02557 for DROP
// TABLESPACE SET's name. Parse takes the statement there too; anywhere
// else a missing operand is a syntax error.
func (p *Parser) tsEndsBeforeOperand() bool {
	return p.tsAtEnd()
}

// tsExpect consumes the word w or fails with a syntax error.
func (p *Parser) tsExpect(w string) error {
	if !p.isKeywordStr(w) {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// tsFinish ends a tablespace statement read into loc, or passes its error
// on.
func (p *Parser) tsFinish(loc *nodes.Loc, err error) error {
	if err != nil {
		return err
	}
	loc.End = p.prev.End
	return nil
}

// tsText is the statement's text from start to end, one token at a time,
// as the AST's summary fields keep it: words and unquoted strings joined by
// single spaces.
func (p *Parser) tsText(start, end int) string {
	lexer := NewLexerRange(p.source, start, end)
	var parts []string
	for {
		tok := lexer.NextToken()
		if tok.Type == tokEOF || lexer.Err != nil {
			break
		}
		s := tok.Str
		if s == "" && tok.Type > 0 && tok.Type < 256 {
			s = string(rune(tok.Type))
		}
		parts = append(parts, s)
	}
	return strings.Join(parts, " ")
}

// isTablespaceNameToken reports whether the current token can name a
// tablespace, a tablespace set or a tablespace group: an identifier, quoted
// or not, that is not a reserved word (ALTER TABLESPACE EXISTS a is
// ORA-02140).
func (p *Parser) isTablespaceNameToken() bool {
	switch {
	case p.cur.Type == tokIDENT || p.cur.Type == tokQIDENT:
		return true
	case p.cur.Type >= 2000:
		return !isOracleSQLReservedKeyword(p.cur)
	}
	return false
}

// parseTablespaceName reads a tablespace or tablespace set name: a single
// identifier (CREATE TABLESPACE a.b is ORA-02180 at the dot). A missing
// name is ORA-02216 or ORA-02140.
func (p *Parser) parseTablespaceName() (*nodes.ObjectName, error) {
	if !p.isTablespaceNameToken() {
		return nil, p.syntaxErrorAtCur()
	}
	tok := p.advance()
	return &nodes.ObjectName{Name: tok.Str, Loc: nodes.Loc{Start: tok.Loc, End: tok.End}}, nil
}

// parseTablespaceIfNotExists reads IF NOT EXISTS before a name. IF not
// followed by NOT is the name IF; IF NOT without EXISTS is ORA-11543 at IF.
func (p *Parser) parseTablespaceIfNotExists() (bool, error) {
	if p.cur.Type != kwIF || p.peekNext().Type != kwNOT {
		return false, nil
	}
	ifTok := p.cur
	p.advance() // IF
	p.advance() // NOT
	if p.cur.Type != kwEXISTS {
		return false, p.syntaxErrorAtTok(ifTok)
	}
	p.advance()
	return true, nil
}

// tsIfExists reads IF EXISTS before a name, if there. IF not followed by
// EXISTS is the name IF: ALTER TABLESPACE IF a DEFAULT is ORA-02142 at a.
func (p *Parser) tsIfExists() bool {
	if p.cur.Type != kwIF || p.peekNext().Type != kwEXISTS {
		return false
	}
	p.advance() // IF
	p.advance() // EXISTS
	return true
}

// parseTablespaceSize reads size_clause: integer [ K | M | G | T | P | E ].
// A size that is missing or not an integer is ORA-02237, ORA-02017 and the
// like.
func (p *Parser) parseTablespaceSize() (string, error) {
	if p.cur.Type != tokICONST {
		return "", p.syntaxErrorAtCur()
	}
	val := p.advance().Str
	if p.cur.Type != tokQIDENT {
		switch p.cur.Str {
		case "K", "M", "G", "T", "P", "E":
			if p.isIdentLike() {
				val += p.advance().Str
			}
		}
	}
	return val, nil
}

// tsSkipSize reads a size_clause the AST does not keep.
func (p *Parser) tsSkipSize() error {
	size, err := p.parseTablespaceSize()
	if err != nil || size == "" {
		return err
	}
	return nil
}

// parseTablespaceAutoextend reads autoextend_clause:
//
//	AUTOEXTEND { OFF | ON [ NEXT size_clause ] [ MAXSIZE { UNLIMITED | size_clause } ] }
//
// ON or OFF missing is ORA-02491, a bad MAXSIZE ORA-02494.
func (p *Parser) parseTablespaceAutoextend() (*nodes.AutoextendClause, error) {
	ac := &nodes.AutoextendClause{Loc: nodes.Loc{Start: p.pos()}}
	p.advance() // AUTOEXTEND
	switch {
	case p.isKeywordStr("OFF"):
		p.advance()
	case p.cur.Type == kwON:
		p.advance()
		ac.On = true
		if p.cur.Type == kwNEXT {
			p.advance()
			next, err := p.parseTablespaceSize()
			if err != nil {
				return nil, err
			}
			ac.Next = next
		}
		if p.isKeywordStr("MAXSIZE") {
			p.advance()
			if p.isKeywordStr("UNLIMITED") {
				p.advance()
				ac.MaxSize = "UNLIMITED"
			} else {
				size, err := p.parseTablespaceSize()
				if err != nil {
					return nil, err
				}
				ac.MaxSize = size
			}
		}
	default:
		return nil, p.syntaxErrorAtCur()
	}
	ac.Loc.End = p.prev.End
	return ac, nil
}

// parseTablespaceFileSpec reads one file_specification of a tablespace:
//
//	'filename' [ SIZE size_clause ] [ REUSE ] [ autoextend_clause ]
//
// The page shows the file name optional, but Oracle wants one after
// DATAFILE or TEMPFILE (CREATE TABLESPACE a DATAFILE SIZE 1M is ORA-02236
// at SIZE, and ALTER TABLESPACE a ADD DATAFILE ORA-02236 at its end).
func (p *Parser) parseTablespaceFileSpec() (*nodes.DatafileClause, error) {
	if p.cur.Type != tokSCONST {
		return nil, p.syntaxErrorAtCur()
	}
	df := &nodes.DatafileClause{Loc: nodes.Loc{Start: p.pos()}, Filename: p.advance().Str}
	if p.cur.Type == kwSIZE {
		p.advance()
		size, err := p.parseTablespaceSize()
		if err != nil {
			return nil, err
		}
		df.Size = size
	}
	if p.isKeywordStr("REUSE") {
		p.advance()
		df.Reuse = true
	}
	if p.isKeywordStr("AUTOEXTEND") {
		ac, err := p.parseTablespaceAutoextend()
		if err != nil {
			return nil, err
		}
		df.Autoextend = ac
	}
	df.Loc.End = p.prev.End
	return df, nil
}

// parseTablespaceFileSpecs reads file_specification [, file_specification]...
func (p *Parser) parseTablespaceFileSpecs() ([]*nodes.DatafileClause, error) {
	var files []*nodes.DatafileClause
	for {
		df, err := p.parseTablespaceFileSpec()
		if err != nil {
			return files, err
		}
		files = append(files, df)
		if p.cur.Type != ',' {
			return files, nil
		}
		p.advance()
	}
}

// ---------------------------------------------------------------------------
// CREATE TABLESPACE
// ---------------------------------------------------------------------------

// tsCreateKind is the CREATE TABLESPACE form an option list belongs to,
// which decides the options Oracle takes there.
type tsCreateKind int

const (
	tsPermanent tsCreateKind = iota // permanent_tablespace_clause
	tsTemporary                     // temporary_tablespace_clause
	tsUndo                          // undo_tablespace_clause
	tsSet                           // CREATE TABLESPACE SET, outside USING TEMPLATE
	tsTemplate                      // CREATE TABLESPACE SET ... USING TEMPLATE ( ... )
)

// parseCreateTablespaceFamily reads a CREATE TABLESPACE or CREATE
// TABLESPACE SET statement after CREATE, at BIGFILE, SMALLFILE, LOCAL,
// TEMPORARY, UNDO or TABLESPACE.
//
//	CREATE [ BIGFILE | SMALLFILE ]
//	    { TABLESPACE [ IF NOT EXISTS ] tablespace permanent options
//	    | [ LOCAL ] TEMPORARY TABLESPACE ... temporary options
//	    | UNDO TABLESPACE [ IF NOT EXISTS ] tablespace undo options
//	    }
//	CREATE TABLESPACE SET [ IF NOT EXISTS ] tablespace_set set options
func (p *Parser) parseCreateTablespaceFamily(start int) (nodes.StmtNode, error) {
	stmt := &nodes.CreateTablespaceStmt{Loc: nodes.Loc{Start: start}}
	switch {
	case p.isKeywordStr("BIGFILE"):
		p.advance()
		stmt.Bigfile = true
	case p.isKeywordStr("SMALLFILE"):
		p.advance()
		stmt.Smallfile = true
	}
	kind := tsPermanent
	if p.isKeywordStr("LOCAL") {
		p.advance()
		stmt.Local = true
		if p.cur.Type != kwTEMPORARY {
			return nil, p.syntaxErrorAtCur()
		}
	}
	temporaryTok := p.cur
	switch {
	case p.cur.Type == kwTEMPORARY:
		p.advance()
		stmt.Temporary = true
		kind = tsTemporary
	case p.isKeywordStr("UNDO"):
		p.advance()
		stmt.Undo = true
		kind = tsUndo
	}
	if p.cur.Type != kwTABLESPACE {
		if kind == tsUndo {
			// CREATE UNDO a is ORA-00901 at UNDO.
			return nil, p.syntaxErrorAtTok(temporaryTok)
		}
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	if kind == tsPermanent && !stmt.Bigfile && !stmt.Smallfile && p.cur.Type == kwSET {
		p.advance()
		return p.parseCreateTablespaceSetStmt(start)
	}
	if err := p.tsFinish(&stmt.Loc, p.parseCreateTablespaceBody(stmt, kind, temporaryTok)); err != nil {
		return nil, err
	}
	return stmt, nil
}

// parseCreateTablespaceBody reads a CREATE TABLESPACE statement after
// TABLESPACE.
//
//	LOCAL TEMPORARY TABLESPACE FOR { ALL | LEAF } [ IF NOT EXISTS ] tablespace ...
//
// Oracle refuses LOCAL TEMPORARY TABLESPACE without FOR (ORA-00901 at
// TEMPORARY), and FOR anywhere else.
func (p *Parser) parseCreateTablespaceBody(stmt *nodes.CreateTablespaceStmt, kind tsCreateKind, temporaryTok Token) error {
	if stmt.Local {
		if p.cur.Type != kwFOR {
			return p.syntaxErrorAtTok(temporaryTok)
		}
		p.advance()
		switch {
		case p.cur.Type == kwALL:
			stmt.ForLeaf = "ALL"
		case p.isKeywordStr("LEAF"):
			stmt.ForLeaf = "LEAF"
		default:
			return p.syntaxErrorAtCur()
		}
		p.advance()
	}
	ifNotExists, err := p.parseTablespaceIfNotExists()
	if err != nil {
		return err
	}
	stmt.IfNotExists = ifNotExists
	if stmt.Name, err = p.parseTablespaceName(); err != nil {
		return err
	}
	return p.parseCreateTablespaceOptions(stmt, nil, kind, map[string]bool{})
}

// tsCreateOptionKey names the CREATE TABLESPACE option the current token
// starts, or "" for none. Options that exclude each other share a key:
// Oracle takes each key once.
func (p *Parser) tsCreateOptionKey() string {
	switch {
	case p.cur.Type == kwLOGGING, p.cur.Type == kwNOLOGGING:
		return "LOGGING"
	case p.cur.Type == kwONLINE, p.cur.Type == kwOFFLINE:
		return "STATE"
	case p.cur.Type == kwDEFAULT:
		return "DEFAULT"
	case p.cur.Type == kwTABLESPACE:
		return "GROUP"
	case p.cur.Type == kwIN:
		return "SHARDSPACE"
	case p.cur.Type == kwUSING:
		return "TEMPLATE"
	case p.cur.Type == tokQIDENT || !p.isIdentLike():
		return ""
	}
	switch p.cur.Str {
	case "AUTOALLOCATE", "UNIFORM":
		return "ALLOCATION"
	case "DATAFILE", "TEMPFILE", "MINIMUM", "BLOCKSIZE", "FORCE", "ENCRYPTION",
		"EXTENT", "SEGMENT", "FLASHBACK", "LOST", "RETENTION":
		return p.cur.Str
	}
	return ""
}

// tsCreateOptionAllowed reports whether a CREATE TABLESPACE form takes the
// option of key: a permanent tablespace (and a tablespace set and its
// template) its permanent_tablespace_attrs, a temporary one its TEMPFILE,
// group, extent management and encryption, an undo one its DATAFILE,
// extent management, retention and encryption.
func tsCreateOptionAllowed(kind tsCreateKind, key string) bool {
	switch key {
	case "ENCRYPTION", "EXTENT":
		return true
	case "TEMPFILE", "GROUP":
		return kind == tsTemporary
	case "RETENTION":
		return kind == tsUndo
	case "DATAFILE":
		return kind != tsTemporary
	case "SHARDSPACE", "TEMPLATE":
		return kind == tsSet
	}
	return kind == tsPermanent || kind == tsSet || kind == tsTemplate
}

// parseCreateTablespaceOptions reads a CREATE TABLESPACE statement's
// options, in any order, each at most once (seen), into stmt; a tablespace
// set's IN SHARDSPACE goes into set. In a template it stops at ')'.
func (p *Parser) parseCreateTablespaceOptions(stmt *nodes.CreateTablespaceStmt, set *nodes.CreateTablespaceSetStmt, kind tsCreateKind, seen map[string]bool) error {
	for !p.tsAtEnd() {
		if kind == tsTemplate && p.cur.Type == ')' {
			return nil
		}
		key := p.tsCreateOptionKey()
		if key == "" || seen[key] || !tsCreateOptionAllowed(kind, key) {
			return p.syntaxErrorAtCur()
		}
		seen[key] = true
		if err := p.parseCreateTablespaceOption(stmt, set, key, seen); err != nil {
			return err
		}
	}
	return nil
}

// parseCreateTablespaceOption reads one CREATE TABLESPACE option, the one
// key names, at its first token.
func (p *Parser) parseCreateTablespaceOption(stmt *nodes.CreateTablespaceStmt, set *nodes.CreateTablespaceSetStmt, key string, seen map[string]bool) error {
	var err error
	switch key {
	case "DATAFILE", "TEMPFILE":
		// { DATAFILE | TEMPFILE } file_specification [, file_specification ]...
		p.advance()
		var files []*nodes.DatafileClause
		files, err = p.parseTablespaceFileSpecs()
		stmt.Datafiles = append(stmt.Datafiles, files...)
	case "MINIMUM":
		// MINIMUM EXTENT size_clause (ORA-25121)
		p.advance()
		if err = p.tsExpect("EXTENT"); err == nil && !p.tsEndsBeforeOperand() {
			stmt.MinimumExtent, err = p.parseTablespaceSize()
		}
	case "BLOCKSIZE":
		// BLOCKSIZE integer [ K ] (ORA-27035)
		p.advance()
		stmt.Blocksize, err = p.parseTablespaceSize()
	case "LOGGING":
		stmt.Logging = tsJoin(stmt.Logging, p.advance().Str)
	case "FORCE":
		// FORCE LOGGING
		p.advance()
		if err = p.tsExpect("LOGGING"); err == nil {
			stmt.Logging = tsJoin(stmt.Logging, "FORCE LOGGING")
		}
	case "ENCRYPTION":
		stmt.Encryption, stmt.EncryptionAlgorithm, err = p.parseCreateTablespaceEncryption()
	case "DEFAULT":
		p.advance()
		stmt.DefaultParams, err = p.parseTablespaceDefaultParams()
	case "STATE":
		if p.advance().Type == kwONLINE {
			stmt.Online = true
		} else {
			stmt.Offline = true
		}
	case "EXTENT":
		stmt.Extent, err = p.parseTablespaceExtentManagement()
	case "ALLOCATION":
		// AUTOALLOCATE | UNIFORM [ SIZE size_clause ]: undocumented, taken
		// by Oracle in a permanent tablespace (ORA-25145 when repeated).
		extent := p.advance().Str
		if extent == "UNIFORM" && p.cur.Type == kwSIZE {
			p.advance()
			extent += " SIZE"
			if !p.tsEndsBeforeOperand() {
				var size string
				if size, err = p.parseTablespaceSize(); err == nil {
					extent += " " + size
				}
			}
		}
		if stmt.Extent == "" {
			stmt.Extent = extent
		}
	case "SEGMENT":
		// SEGMENT SPACE MANAGEMENT { AUTO | MANUAL } (ORA-30571)
		p.advance()
		if err = p.tsExpect("SPACE"); err == nil {
			err = p.tsExpect("MANAGEMENT")
		}
		if err == nil {
			if p.isKeywordStr("AUTO") || p.isKeywordStr("MANUAL") {
				stmt.Segment = p.advance().Str
			} else {
				err = p.syntaxErrorAtCur()
			}
		}
	case "FLASHBACK":
		stmt.Flashback, err = p.parseTablespaceFlashback()
	case "LOST":
		// LOST WRITE PROTECTION (ORA-65482)
		p.advance()
		if err = p.tsExpect("WRITE"); err == nil {
			err = p.tsExpect("PROTECTION")
		}
		if err == nil {
			stmt.LostWriteProtection = "LOST WRITE PROTECTION"
		}
	case "GROUP":
		stmt.TablespaceGroup, err = p.parseTablespaceGroup()
	case "RETENTION":
		stmt.Retention, err = p.parseTablespaceRetention()
	case "SHARDSPACE":
		// IN SHARDSPACE shardspace_name: SHARDSPACE missing is ORA-02507,
		// the name ORA-00931.
		p.advance()
		if err = p.tsExpect("SHARDSPACE"); err == nil {
			if !p.isTablespaceNameToken() {
				return p.syntaxErrorAtCur()
			}
			set.Shardspace = p.advance().Str
		}
	case "TEMPLATE":
		// USING TEMPLATE ( options ), the options of a permanent tablespace.
		p.advance()
		if err = p.tsExpect("TEMPLATE"); err != nil {
			return err
		}
		if p.cur.Type != '(' {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		if err = p.parseCreateTablespaceOptions(stmt, set, tsTemplate, seen); err != nil {
			return err
		}
		if p.cur.Type != ')' {
			return p.syntaxErrorAtCur()
		}
		p.advance()
	}
	return err
}

// tsJoin joins two words of a summary field with a space.
func tsJoin(a, b string) string {
	if a == "" {
		return b
	}
	return a + " " + b
}

// parseCreateTablespaceEncryption reads CREATE's tablespace_encryption_clause:
//
//	ENCRYPTION [ USING 'encrypt_algorithm' ] [ ENCRYPT | DECRYPT ]
//
// Oracle takes it without ENCRYPT or DECRYPT. It returns the clause's
// ENCRYPT, DECRYPT or ENCRYPTION, and the algorithm.
func (p *Parser) parseCreateTablespaceEncryption() (string, string, error) {
	p.advance() // ENCRYPTION
	algorithm := ""
	if p.cur.Type == kwUSING {
		p.advance()
		if p.cur.Type != tokSCONST || !isValidOracleEncryptionAlgorithm(p.cur.Str) {
			return "", "", p.syntaxErrorAtCur() // ORA-28339
		}
		algorithm = p.advance().Str
	}
	if p.isKeywordStr("ENCRYPT") || p.isKeywordStr("DECRYPT") {
		return p.advance().Str, algorithm, nil
	}
	return "ENCRYPTION", algorithm, nil
}

// parseTablespaceExtentManagement reads extent_management_clause:
//
//	EXTENT MANAGEMENT { LOCAL [ AUTOALLOCATE | UNIFORM [ SIZE size_clause ] ] | DICTIONARY }
//
// Its errors are ORA-25141.
func (p *Parser) parseTablespaceExtentManagement() (string, error) {
	p.advance() // EXTENT
	if err := p.tsExpect("MANAGEMENT"); err != nil {
		return "", err
	}
	switch {
	case p.isKeywordStr("DICTIONARY"):
		p.advance()
		return "DICTIONARY", nil
	case p.isKeywordStr("LOCAL"):
		p.advance()
	default:
		return "", p.syntaxErrorAtCur()
	}
	switch {
	case p.isKeywordStr("AUTOALLOCATE"):
		p.advance()
		return "AUTOALLOCATE", nil
	case p.isKeywordStr("UNIFORM"):
		p.advance()
		if p.cur.Type != kwSIZE {
			return "UNIFORM", nil
		}
		p.advance()
		if p.tsEndsBeforeOperand() {
			return "UNIFORM SIZE", nil
		}
		size, err := p.parseTablespaceSize()
		if err != nil {
			return "", err
		}
		return "UNIFORM SIZE " + size, nil
	}
	return "LOCAL", nil
}

// parseTablespaceFlashback reads flashback_mode_clause: FLASHBACK { ON | OFF }
// (ORA-38722).
func (p *Parser) parseTablespaceFlashback() (string, error) {
	p.advance() // FLASHBACK
	if p.cur.Type == kwON || p.isKeywordStr("OFF") {
		return p.advance().Str, nil
	}
	return "", p.syntaxErrorAtCur()
}

// parseTablespaceRetention reads tablespace_retention_clause:
// RETENTION { GUARANTEE | NOGUARANTEE }.
func (p *Parser) parseTablespaceRetention() (string, error) {
	p.advance() // RETENTION
	if p.isKeywordStr("GUARANTEE") || p.isKeywordStr("NOGUARANTEE") {
		return p.advance().Str, nil
	}
	return "", p.syntaxErrorAtCur()
}

// parseTablespaceGroup reads tablespace_group_clause:
//
//	TABLESPACE GROUP { tablespace_group_name | '' }
//
// A string that is not empty is ORA-10914, as is a missing GROUP or name.
func (p *Parser) parseTablespaceGroup() (string, error) {
	p.advance() // TABLESPACE
	if err := p.tsExpect("GROUP"); err != nil {
		return "", err
	}
	switch {
	case p.cur.Type == tokSCONST:
		if p.cur.Str != "" {
			return "", p.syntaxErrorAtCur()
		}
		p.advance()
		return "''", nil
	case p.isTablespaceNameToken():
		return p.advance().Str, nil
	}
	return "", p.syntaxErrorAtCur()
}

// ---------------------------------------------------------------------------
// default_tablespace_params
// ---------------------------------------------------------------------------

// tsDefaultParamKey names the default_tablespace_params entry the current
// token starts, or "" for none. Each is taken once: a second table
// compression is ORA-14460, a second inmemory_clause ORA-03218.
// Compression written without TABLE (COMPRESS, NOCOMPRESS, ROW STORE,
// COLUMN STORE) is "BARE": Oracle takes no index compression beside it
// (DEFAULT COMPRESS INDEX NOCOMPRESS is ORA-14460 at INDEX).
func (p *Parser) tsDefaultParamKey() string {
	if p.cur.Type == tokQIDENT || !p.isIdentLike() {
		return ""
	}
	switch p.cur.Str {
	case "TABLE":
		return "COMPRESSION"
	case "COMPRESS", "NOCOMPRESS", "ROW", "COLUMN":
		return "BARE"
	case "INDEX":
		return "INDEX"
	case "INMEMORY", "NO":
		return "INMEMORY"
	case "ILM":
		return "ILM"
	case "STORAGE":
		return "STORAGE"
	}
	return ""
}

// parseTablespaceDefaultParams reads default_tablespace_params after
// DEFAULT, one or more, and returns their text. DEFAULT with none is
// ORA-00905. storage_clause ends the list: Oracle reads no parameter after
// it (DEFAULT STORAGE (INITIAL 1M) COMPRESS is ORA-02180 at COMPRESS).
func (p *Parser) parseTablespaceDefaultParams() (string, error) {
	start := p.pos()
	seen := map[string]bool{}
	for !seen["STORAGE"] {
		key := p.tsDefaultParamKey()
		if key == "" {
			break
		}
		if seen[key] || key == "BARE" && (seen["COMPRESSION"] || seen["INDEX"]) ||
			(key == "COMPRESSION" || key == "INDEX") && seen["BARE"] {
			return "", p.syntaxErrorAtCur()
		}
		seen[key] = true
		var err error
		switch key {
		case "COMPRESSION", "BARE":
			err = p.parseTablespaceTableCompression(true)
		case "INDEX":
			err = p.parseTablespaceIndexCompression()
		case "INMEMORY":
			err = p.parseTablespaceInmemoryClause()
		case "ILM":
			err = p.parseTablespaceILMClause()
		case "STORAGE":
			err = p.parseTablespaceStorageClause()
		}
		if err != nil {
			return "", err
		}
	}
	if len(seen) == 0 {
		return "", p.syntaxErrorAtCur()
	}
	return p.tsText(start, p.prev.End), nil
}

// parseTablespaceTableCompression reads default_table_compression, or with
// inDefault false the compression of an ILM policy:
//
//	[ TABLE ] { COMPRESS [ BASIC | FOR { OLTP | QUERY [ LOW | HIGH ] | ARCHIVE [ LOW | HIGH ] } ]
//	          | NOCOMPRESS }
//	| ROW STORE [ COMPRESS [ BASIC | ADVANCED ] ]
//	| COLUMN STORE [ COMPRESS [ FOR { QUERY | ARCHIVE } [ LOW | HIGH ] ] ]
//
// TABLE and NOCOMPRESS stand only in DEFAULT. TABLE takes only COMPRESS or
// NOCOMPRESS (DEFAULT TABLE ROW STORE COMPRESS is ORA-00922). A missing
// compression type after FOR is ORA-14464.
func (p *Parser) parseTablespaceTableCompression(inDefault bool) error {
	switch {
	case inDefault && p.cur.Type == kwTABLE:
		p.advance()
		switch {
		case p.cur.Type == kwCOMPRESS:
			p.advance()
			return p.parseTablespaceCompressType()
		case p.cur.Type == kwNOCOMPRESS:
			p.advance()
			return nil
		}
		return p.syntaxErrorAtCur()
	case p.cur.Type == kwCOMPRESS:
		p.advance()
		return p.parseTablespaceCompressType()
	case inDefault && p.cur.Type == kwNOCOMPRESS:
		p.advance()
		return nil
	case p.cur.Type == kwROW:
		p.advance()
		if err := p.tsExpect("STORE"); err != nil {
			return err
		}
		if p.cur.Type == kwCOMPRESS {
			p.advance()
			if p.isKeywordStr("BASIC") || p.isKeywordStr("ADVANCED") {
				p.advance()
			}
		}
		return nil
	case p.cur.Type == kwCOLUMN:
		p.advance()
		if err := p.tsExpect("STORE"); err != nil {
			return err
		}
		if p.cur.Type != kwCOMPRESS {
			return nil
		}
		p.advance()
		if p.cur.Type != kwFOR {
			return nil
		}
		p.advance()
		if !p.isKeywordStr("QUERY") && !p.isKeywordStr("ARCHIVE") {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		if p.isKeywordStr("LOW") || p.isKeywordStr("HIGH") {
			p.advance()
		}
		return nil
	}
	return p.syntaxErrorAtCur()
}

// parseTablespaceCompressType reads what follows COMPRESS:
//
//	[ BASIC | FOR { OLTP | QUERY [ LOW | HIGH ] | ARCHIVE [ LOW | HIGH ] } ]
func (p *Parser) parseTablespaceCompressType() error {
	switch {
	case p.isKeywordStr("BASIC"):
		p.advance()
	case p.cur.Type == kwFOR:
		p.advance()
		switch {
		case p.isKeywordStr("OLTP"):
			p.advance()
		case p.isKeywordStr("QUERY"), p.isKeywordStr("ARCHIVE"):
			p.advance()
			if p.isKeywordStr("LOW") || p.isKeywordStr("HIGH") {
				p.advance()
			}
		default:
			return p.syntaxErrorAtCur()
		}
	}
	return nil
}

// parseTablespaceIndexCompression reads default_index_compression:
//
//	INDEX { COMPRESS ADVANCED { LOW | HIGH } | NOCOMPRESS }
//
// The page shows ADVANCED and its level optional; Oracle wants both.
func (p *Parser) parseTablespaceIndexCompression() error {
	p.advance() // INDEX
	switch p.cur.Type {
	case kwNOCOMPRESS:
		p.advance()
		return nil
	case kwCOMPRESS:
		p.advance()
	default:
		return p.syntaxErrorAtCur()
	}
	if err := p.tsExpect("ADVANCED"); err != nil {
		return err
	}
	if !p.isKeywordStr("LOW") && !p.isKeywordStr("HIGH") {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// parseTablespaceInmemoryClause reads inmemory_clause:
//
//	INMEMORY [ inmemory_attributes ] | NO INMEMORY
//
// NO not followed by INMEMORY is ORA-03218 at NO.
func (p *Parser) parseTablespaceInmemoryClause() error {
	if p.isKeywordStr("NO") {
		noTok := p.advance()
		if !p.isKeywordStr("INMEMORY") {
			return p.syntaxErrorAtTok(noTok)
		}
		p.advance()
		return nil
	}
	p.advance() // INMEMORY
	return p.parseTablespaceInmemoryAttributes()
}

// tsNextIs reports whether the token after the current one is the word w.
func (p *Parser) tsNextIs(w string) bool {
	next := p.peekNext()
	return next.Type != tokQIDENT && (next.Type == tokIDENT || next.Type >= 2000) && next.Str == w
}

// parseTablespaceInmemoryAttributes reads inmemory_attributes, in any
// order:
//
//	inmemory_memcompress | PRIORITY { NONE | LOW | MEDIUM | HIGH | CRITICAL }
//	| DISTRIBUTE [ AUTO | BY { ROWID RANGE | PARTITION | SUBPARTITION } ]
//	| DUPLICATE [ ALL ] | NO DUPLICATE
//
// It stops before NO INMEMORY, which the caller reads as a second
// inmemory_clause.
func (p *Parser) parseTablespaceInmemoryAttributes() error {
	for {
		switch {
		case p.isKeywordStr("MEMCOMPRESS"):
			if err := p.parseTablespaceMemcompress(); err != nil {
				return err
			}
		case p.isKeywordStr("NO"):
			switch {
			case p.tsNextIs("MEMCOMPRESS"), p.tsNextIs("DUPLICATE"):
				p.advance()
				p.advance()
			case p.tsNextIs("INMEMORY"):
				return nil
			default:
				p.advance()
				return p.syntaxErrorAtCur()
			}
		case p.isKeywordStr("PRIORITY"):
			p.advance()
			switch p.cur.Str {
			case "NONE", "LOW", "MEDIUM", "HIGH", "CRITICAL":
				if p.isKeywordStr(p.cur.Str) {
					p.advance()
					continue
				}
			}
			return p.syntaxErrorAtCur()
		case p.isKeywordStr("DISTRIBUTE"):
			p.advance()
			switch {
			case p.isKeywordStr("AUTO"):
				p.advance()
			case p.cur.Type == kwBY:
				// Oracle names BY when what follows it is wrong (ORA-00922).
				byTok := p.advance()
				switch {
				case p.cur.Type == kwROWID:
					p.advance()
					if !p.isKeywordStr("RANGE") {
						return p.syntaxErrorAtTok(byTok)
					}
					p.advance()
				case p.cur.Type == kwPARTITION, p.isKeywordStr("SUBPARTITION"):
					// Oracle refuses distribution by partition in a
					// tablespace's defaults, at what follows it (ORA-03218).
					p.advance()
					return p.syntaxErrorAtCur()
				default:
					return p.syntaxErrorAtTok(byTok)
				}
			}
		case p.isKeywordStr("DUPLICATE"):
			p.advance()
			if p.cur.Type == kwALL {
				p.advance()
			}
		default:
			return nil
		}
	}
}

// parseTablespaceMemcompress reads MEMCOMPRESS FOR { DML | QUERY [ LOW | HIGH ]
// | CAPACITY [ LOW | HIGH ] }.
func (p *Parser) parseTablespaceMemcompress() error {
	p.advance() // MEMCOMPRESS
	if p.cur.Type != kwFOR {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	switch {
	case p.isKeywordStr("DML"):
		p.advance()
	case p.isKeywordStr("QUERY"), p.isKeywordStr("CAPACITY"):
		p.advance()
		if p.isKeywordStr("LOW") || p.isKeywordStr("HIGH") {
			p.advance()
		}
	default:
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseTablespaceILMClause reads ilm_clause:
//
//	ILM { ADD POLICY ilm_policy_clause
//	    | { DELETE | ENABLE | DISABLE } POLICY policy_name
//	    | DELETE_ALL | ENABLE_ALL | DISABLE_ALL }
//
// Nothing may follow DELETE_ALL, ENABLE_ALL or DISABLE_ALL (ORA-38346).
func (p *Parser) parseTablespaceILMClause() error {
	p.advance() // ILM
	switch {
	case p.cur.Type == kwADD:
		p.advance()
		if err := p.tsExpect("POLICY"); err != nil {
			return err
		}
		return p.parseTablespaceILMPolicy()
	case p.cur.Type == kwDELETE, p.isKeywordStr("ENABLE"), p.isKeywordStr("DISABLE"):
		p.advance()
		if err := p.tsExpect("POLICY"); err != nil {
			return err
		}
		if !p.isTablespaceNameToken() {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		return nil
	case p.isKeywordStr("DELETE_ALL"), p.isKeywordStr("ENABLE_ALL"), p.isKeywordStr("DISABLE_ALL"):
		p.advance()
		if !p.tsAtEnd() {
			return p.syntaxErrorAtCur()
		}
		return nil
	}
	return p.syntaxErrorAtCur()
}

// parseTablespaceILMPolicy reads ilm_policy_clause after ADD POLICY: a
// compression, tiering or inmemory policy. Any other action is ORA-38339.
//
//	ilm_compression_policy:
//	    table_compression { SEGMENT | GROUP | ROW } AFTER ilm_time_period OF ilm_condition
//	ilm_tiering_policy:
//	    TIER TO tablespace [ SEGMENT | GROUP ]
//	  | TIER TO tablespace READ ONLY [ SEGMENT | GROUP ] AFTER ilm_time_period OF ilm_condition
//	ilm_inmemory_policy:
//	    { SET INMEMORY [ inmemory_attributes ] | MODIFY INMEMORY [ inmemory_memcompress ]
//	    | NO INMEMORY } [ SEGMENT ] AFTER ilm_time_period OF ilm_condition
//
// A compression policy without its scope is ORA-38322. The ON function_name
// forms are refused in a tablespace (ORA-38340).
func (p *Parser) parseTablespaceILMPolicy() error {
	switch {
	case p.cur.Type == kwCOMPRESS, p.cur.Type == kwROW, p.cur.Type == kwCOLUMN:
		if err := p.parseTablespaceTableCompression(false); err != nil {
			return err
		}
		switch {
		case p.isKeywordStr("SEGMENT"), p.cur.Type == kwGROUP, p.cur.Type == kwROW:
			p.advance()
		default:
			return p.syntaxErrorAtCur()
		}
		return p.parseTablespaceILMAfter()
	case p.isKeywordStr("TIER"):
		p.advance()
		if err := p.tsExpect("TO"); err != nil {
			return err
		}
		if !p.isTablespaceNameToken() {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		readOnly := false
		if p.cur.Type == kwREAD {
			p.advance()
			if err := p.tsExpect("ONLY"); err != nil {
				return err
			}
			readOnly = true
		}
		switch {
		case p.isKeywordStr("SEGMENT"), p.cur.Type == kwGROUP:
			p.advance()
		case p.cur.Type == kwROW:
			// Oracle reads ROW here and refuses what follows it
			// (ORA-00922): a tiering policy has no row scope.
			p.advance()
			return p.syntaxErrorAtCur()
		}
		if readOnly {
			return p.parseTablespaceILMAfter()
		}
		return nil
	case p.cur.Type == kwSET:
		p.advance()
		if err := p.tsExpect("INMEMORY"); err != nil {
			return err
		}
		if err := p.parseTablespaceInmemoryAttributes(); err != nil {
			return err
		}
	case p.cur.Type == kwMODIFY:
		p.advance()
		if err := p.tsExpect("INMEMORY"); err != nil {
			return err
		}
		switch {
		case p.isKeywordStr("MEMCOMPRESS"):
			if err := p.parseTablespaceMemcompress(); err != nil {
				return err
			}
		case p.isKeywordStr("NO"):
			p.advance()
			if err := p.tsExpect("MEMCOMPRESS"); err != nil {
				return err
			}
		}
	case p.isKeywordStr("NO"):
		p.advance()
		if err := p.tsExpect("INMEMORY"); err != nil {
			return err
		}
	default:
		return p.syntaxErrorAtCur()
	}
	if p.isKeywordStr("SEGMENT") {
		p.advance()
	}
	return p.parseTablespaceILMAfter()
}

// parseTablespaceILMAfter reads AFTER ilm_time_period OF ilm_condition:
//
//	AFTER integer [ DAY | DAYS | MONTH | MONTHS | YEAR | YEARS ]
//	OF { NO ACCESS | NO MODIFICATION | LOW ACCESS | CREATION }
//
// Oracle takes the period without its unit, and LOW ACCESS.
func (p *Parser) parseTablespaceILMAfter() error {
	if err := p.tsExpect("AFTER"); err != nil {
		return err
	}
	if p.cur.Type != tokICONST {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	switch p.cur.Str {
	case "DAY", "DAYS", "MONTH", "MONTHS", "YEAR", "YEARS":
		if p.isKeywordStr(p.cur.Str) {
			p.advance()
		}
	}
	if err := p.tsExpect("OF"); err != nil {
		return err
	}
	switch {
	case p.isKeywordStr("NO"):
		p.advance()
		if p.isKeywordStr("ACCESS") || p.isKeywordStr("MODIFICATION") {
			p.advance()
			return nil
		}
		return p.syntaxErrorAtCur()
	case p.isKeywordStr("LOW"):
		p.advance()
		return p.tsExpect("ACCESS")
	case p.isKeywordStr("CREATION"):
		p.advance()
		return nil
	}
	return p.syntaxErrorAtCur()
}

// parseTablespaceStorageClause reads the storage_clause of
// default_tablespace_params
// (https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/storage_clause.html):
//
//	STORAGE ( { INITIAL size_clause | NEXT size_clause | MINEXTENTS integer
//	          | MAXEXTENTS { integer | UNLIMITED } | MAXSIZE { UNLIMITED | size_clause }
//	          | PCTINCREASE integer | ENCRYPT } ... )
//
// Oracle refuses the page's other options in a tablespace's default
// storage: OPTIMAL (ORA-02143), FREELISTS (ORA-02169), FREELIST GROUPS
// (ORA-02170), BUFFER_POOL, FLASH_CACHE and CELL_FLASH_CACHE (ORA-25125).
// A missing '(' is ORA-00906; the errors inside are the storage clause's
// own (ORA-02143 invalid STORAGE option, ORA-02145 missing STORAGE option).
func (p *Parser) parseTablespaceStorageClause() error {
	p.advance() // STORAGE
	if p.cur.Type != '(' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	n := 0
	for p.cur.Type != ')' {
		if err := p.parseTablespaceStorageOption(); err != nil {
			return err
		}
		n++
	}
	if n == 0 {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// parseTablespaceStorageOption reads one option of storage_clause.
func (p *Parser) parseTablespaceStorageOption() error {
	if p.cur.Type == tokQIDENT || !p.isIdentLike() {
		return p.syntaxErrorAtCur()
	}
	integer := func() error {
		if p.cur.Type != tokICONST {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		return nil
	}
	switch p.cur.Str {
	case "INITIAL", "NEXT":
		p.advance()
		return p.tsSkipSize()
	case "MINEXTENTS", "PCTINCREASE":
		p.advance()
		return integer()
	case "MAXEXTENTS":
		p.advance()
		if p.isKeywordStr("UNLIMITED") {
			p.advance()
			return nil
		}
		return integer()
	case "MAXSIZE":
		p.advance()
		if p.isKeywordStr("UNLIMITED") {
			p.advance()
			return nil
		}
		return p.tsSkipSize()
	case "ENCRYPT":
		p.advance()
		return nil
	}
	return p.syntaxErrorAtCur()
}

// ---------------------------------------------------------------------------
// CREATE TABLESPACE SET
// ---------------------------------------------------------------------------

// parseCreateTablespaceSetStmt reads CREATE TABLESPACE SET after SET.
//
// BNF: oracle/parser/bnf/CREATE-TABLESPACE-SET.bnf
//
//	CREATE TABLESPACE SET tablespace_set
//	    [ IN SHARDSPACE shardspace_name ]
//	    [ USING TEMPLATE ( [ DATAFILE file_specification ] [ permanent_tablespace_attrs ] ) ]
//
// Oracle reads it as CREATE TABLESPACE's option list, with IN SHARDSPACE
// and USING TEMPLATE among the options: it takes the permanent options
// outside the template too, and IF NOT EXISTS before the name. An option
// is taken once, in or out of the template (USING TEMPLATE (LOGGING)
// LOGGING is ORA-25119).
func (p *Parser) parseCreateTablespaceSetStmt(start int) (*nodes.CreateTablespaceSetStmt, error) {
	stmt := &nodes.CreateTablespaceSetStmt{Loc: nodes.Loc{Start: start}}
	opts := &nodes.CreateTablespaceStmt{}
	err := p.tsFinish(&stmt.Loc, func() error {
		ifNotExists, err := p.parseTablespaceIfNotExists()
		if err != nil {
			return err
		}
		stmt.IfNotExists = ifNotExists
		if stmt.Name, err = p.parseTablespaceName(); err != nil {
			return err
		}
		return p.parseCreateTablespaceOptions(opts, stmt, tsSet, map[string]bool{})
	}())
	if err != nil {
		return nil, err
	}
	stmt.Datafiles = opts.Datafiles
	stmt.Logging = opts.Logging
	stmt.Encryption = opts.Encryption
	stmt.DefaultParams = opts.DefaultParams
	stmt.Extent = opts.Extent
	stmt.Segment = opts.Segment
	stmt.Flashback = opts.Flashback
	return stmt, nil
}

// ---------------------------------------------------------------------------
// ALTER TABLESPACE
// ---------------------------------------------------------------------------

// parseAlterTablespaceStmt reads ALTER TABLESPACE [ SET ] after TABLESPACE
// or SET.
//
// BNF: oracle/parser/bnf/ALTER-TABLESPACE.bnf, ALTER-TABLESPACE-SET.bnf
//
//	ALTER TABLESPACE [ SET ] [ IF EXISTS ] tablespace alter_tablespace_attrs
//
// Oracle takes exactly one clause (ALTER TABLESPACE a LOGGING ONLINE is
// ORA-03048 at ONLINE), the same ones for a tablespace set; an unknown
// clause is ORA-02142.
func (p *Parser) parseAlterTablespaceStmt(start int, isSet bool) (*nodes.AlterTablespaceStmt, error) {
	stmt := &nodes.AlterTablespaceStmt{IsSet: isSet, Loc: nodes.Loc{Start: start}}
	if err := p.tsFinish(&stmt.Loc, p.parseAlterTablespaceBody(stmt)); err != nil {
		return nil, err
	}
	return stmt, nil
}

func (p *Parser) parseAlterTablespaceBody(stmt *nodes.AlterTablespaceStmt) error {
	stmt.IfExists = p.tsIfExists()
	var err error
	if stmt.Name, err = p.parseTablespaceName(); err != nil {
		return err
	}
	if err := p.parseAlterTablespaceClause(stmt); err != nil {
		return err
	}
	if !p.tsAtEnd() {
		return p.syntaxErrorAtCur()
	}
	return nil
}

// parseAlterTablespaceClause reads one alter_tablespace_attrs clause.
//
//	{ DEFAULT default_tablespace_params | MINIMUM EXTENT size_clause
//	| RESIZE size_clause | COALESCE | SHRINK SPACE [ KEEP size_clause ]
//	| RENAME TO new_tablespace_name | { BEGIN | END } BACKUP
//	| datafile_tempfile_clauses | tablespace_logging_clauses
//	| tablespace_group_clause | tablespace_state_clauses
//	| autoextend_clause | flashback_mode_clause
//	| tablespace_retention_clause | alter_tablespace_encryption
//	| { ENABLE | SUSPEND | REMOVE } LOST WRITE PROTECTION }
func (p *Parser) parseAlterTablespaceClause(stmt *nodes.AlterTablespaceStmt) error {
	var err error
	switch {
	case p.cur.Type == kwDEFAULT:
		p.advance()
		stmt.DefaultParams, err = p.parseTablespaceDefaultParams()
	case p.isKeywordStr("MINIMUM"):
		p.advance()
		if err = p.tsExpect("EXTENT"); err == nil && !p.tsEndsBeforeOperand() {
			stmt.MinimumExtent, err = p.parseTablespaceSize()
		}
	case p.isKeywordStr("RESIZE"):
		p.advance()
		stmt.Resize, err = p.parseTablespaceSize()
	case p.isKeywordStr("COALESCE"):
		p.advance()
		stmt.Coalesce = true
	case p.isKeywordStr("SHRINK"):
		err = p.parseAlterTablespaceShrink(stmt)
	case p.cur.Type == kwRENAME:
		err = p.parseAlterTablespaceRename(stmt)
	case p.isKeywordStr("BEGIN"), p.isKeywordStr("END"):
		// { BEGIN | END } BACKUP (ORA-01906)
		begin := p.advance().Str == "BEGIN"
		if err = p.tsExpect("BACKUP"); err == nil {
			stmt.BeginBackup, stmt.EndBackup = begin, !begin
		}
	case p.cur.Type == kwADD:
		err = p.parseAlterTablespaceAddFile(stmt)
	case p.cur.Type == kwDROP:
		// DROP { DATAFILE | TEMPFILE } { 'filename' | file_number }
		p.advance()
		switch {
		case p.isKeywordStr("DATAFILE"):
			stmt.DropDatafile = true
		case p.isKeywordStr("TEMPFILE"):
			stmt.DropTempfile = true
		default:
			return p.syntaxErrorAtCur()
		}
		p.advance()
		if p.tsEndsBeforeOperand() {
			return nil
		}
		// ALTER TABLESPACE SET uses OMF and takes no datafile name or number
		// (ORA-02580 at the operand; ALTER-TABLESPACE-SET.bnf).
		if stmt.IsSet && stmt.DropDatafile {
			return p.syntaxErrorAtCur()
		}
		if p.cur.Type != tokSCONST && p.cur.Type != tokICONST {
			return p.syntaxErrorAtCur() // ORA-02236
		}
		stmt.DropFileRef = p.advance().Str
	case p.isKeywordStr("DATAFILE"), p.isKeywordStr("TEMPFILE"):
		// { DATAFILE | TEMPFILE } { ONLINE | OFFLINE }
		datafile := p.advance().Str == "DATAFILE"
		switch p.cur.Type {
		case kwONLINE:
			stmt.DatafileOnline, stmt.TempfileOnline = datafile, !datafile
		case kwOFFLINE:
			stmt.DatafileOffline, stmt.TempfileOffline = datafile, !datafile
		default:
			return p.syntaxErrorAtCur()
		}
		p.advance()
	case p.cur.Type == kwLOGGING, p.cur.Type == kwNOLOGGING:
		stmt.Logging = p.advance().Str
	case p.isKeywordStr("FORCE"):
		// FORCE LOGGING
		p.advance()
		if err = p.tsExpect("LOGGING"); err == nil {
			stmt.ForceLogging = "FORCE LOGGING"
		}
	case p.isKeywordStr("NO"):
		// NO FORCE LOGGING
		p.advance()
		if err = p.tsExpect("FORCE"); err == nil {
			err = p.tsExpect("LOGGING")
		}
		if err == nil {
			stmt.ForceLogging = "NO FORCE LOGGING"
		}
	case p.cur.Type == kwTABLESPACE:
		stmt.TablespaceGroup, err = p.parseTablespaceGroup()
	case p.cur.Type == kwONLINE:
		p.advance()
		stmt.Online = true
	case p.cur.Type == kwOFFLINE:
		// OFFLINE [ NORMAL | TEMPORARY | IMMEDIATE | FOR RECOVER ]
		p.advance()
		stmt.Offline = true
		switch {
		case p.isKeywordStr("NORMAL"), p.cur.Type == kwTEMPORARY, p.cur.Type == kwIMMEDIATE:
			stmt.OfflineMode = p.advance().Str
		case p.cur.Type == kwFOR:
			p.advance()
			if err = p.tsExpect("RECOVER"); err == nil {
				stmt.OfflineMode = "FOR RECOVER"
			}
		}
	case p.cur.Type == kwREAD:
		// READ { ONLY | WRITE }
		p.advance()
		switch {
		case p.isKeywordStr("ONLY"):
			stmt.ReadOnly = true
		case p.isKeywordStr("WRITE"):
			stmt.ReadWrite = true
		default:
			return p.syntaxErrorAtCur()
		}
		p.advance()
	case p.isKeywordStr("PERMANENT"):
		p.advance()
		stmt.Permanent = true
	case p.cur.Type == kwTEMPORARY:
		p.advance()
		stmt.TempState = true
	case p.isKeywordStr("AUTOEXTEND"):
		stmt.Autoextend, err = p.parseTablespaceAutoextend()
	case p.isKeywordStr("FLASHBACK"):
		stmt.Flashback, err = p.parseTablespaceFlashback()
	case p.isKeywordStr("RETENTION"):
		stmt.Retention, err = p.parseTablespaceRetention()
	case p.isKeywordStr("ENCRYPTION"):
		start := p.pos()
		err = p.parseAlterTablespaceEncryption()
		stmt.Encryption = p.tsText(start, p.prev.End)
	case p.isKeywordStr("ENABLE"), p.isKeywordStr("SUSPEND"), p.isKeywordStr("REMOVE"):
		// { ENABLE | SUSPEND | REMOVE } LOST WRITE PROTECTION
		action := p.advance().Str
		for _, w := range []string{"LOST", "WRITE", "PROTECTION"} {
			if err = p.tsExpect(w); err != nil {
				return err
			}
		}
		stmt.LostWriteProtection = action
	default:
		return p.syntaxErrorAtCur()
	}
	return err
}

// parseAlterTablespaceShrink reads SHRINK SPACE [ KEEP size_clause ] or
// SHRINK TEMPFILE { 'filename' | file_number } [ KEEP size_clause ].
func (p *Parser) parseAlterTablespaceShrink(stmt *nodes.AlterTablespaceStmt) error {
	p.advance() // SHRINK
	keep := &stmt.ShrinkKeep
	switch {
	case p.isKeywordStr("SPACE"):
		p.advance()
		stmt.ShrinkSpace = true
	case p.isKeywordStr("TEMPFILE"):
		p.advance()
		if p.cur.Type != tokSCONST && p.cur.Type != tokICONST {
			return p.syntaxErrorAtCur()
		}
		stmt.ShrinkTempfile = p.advance().Str
		keep = &stmt.ShrinkTempfileKeep
	default:
		return p.syntaxErrorAtCur()
	}
	if !p.isKeywordStr("KEEP") {
		return nil
	}
	p.advance()
	size, err := p.parseTablespaceSize()
	*keep = size
	return err
}

// parseAlterTablespaceRename reads RENAME TO new_tablespace_name or
//
//	RENAME DATAFILE 'filename' [, 'filename' ]... TO 'filename' [, 'filename' ]...
//
// Anything else after RENAME is ORA-02152, a missing file name ORA-02236
// and a missing TO ORA-00946.
func (p *Parser) parseAlterTablespaceRename(stmt *nodes.AlterTablespaceStmt) error {
	p.advance() // RENAME
	switch {
	case p.cur.Type == kwTO:
		p.advance()
		name, err := p.parseTablespaceName()
		if err != nil {
			return err
		}
		stmt.RenameTo = name.Name
		return nil
	case p.isKeywordStr("DATAFILE"):
		p.advance()
	default:
		return p.syntaxErrorAtCur()
	}
	stmt.RenameDatafile = true
	files := func(into *[]string) error {
		for {
			if p.cur.Type != tokSCONST {
				return p.syntaxErrorAtCur()
			}
			*into = append(*into, p.advance().Str)
			if p.cur.Type != ',' {
				return nil
			}
			p.advance()
		}
	}
	if err := files(&stmt.RenameFrom); err != nil {
		return err
	}
	if p.cur.Type != kwTO {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return files(&stmt.RenameTo2)
}

// parseAlterTablespaceAddFile reads
//
//	ADD { DATAFILE | TEMPFILE } file_specification [, file_specification ]...
//
// The page shows the file specifications optional; Oracle wants one
// (ALTER TABLESPACE a ADD DATAFILE is ORA-02236 at its end).
func (p *Parser) parseAlterTablespaceAddFile(stmt *nodes.AlterTablespaceStmt) error {
	p.advance() // ADD
	switch {
	case p.isKeywordStr("DATAFILE"):
		stmt.AddDatafile = true
	case p.isKeywordStr("TEMPFILE"):
		stmt.AddTempfile = true
	default:
		return p.syntaxErrorAtCur()
	}
	p.advance()
	if stmt.IsSet && stmt.AddDatafile {
		if p.tsEndsBeforeOperand() {
			return nil
		}
		// ALTER TABLESPACE SET uses OMF and takes no file specification
		// (ORA-02580 at the operand; ALTER-TABLESPACE-SET.bnf).
		return p.syntaxErrorAtCur()
	}
	files, err := p.parseTablespaceFileSpecs()
	stmt.Datafiles = append(stmt.Datafiles, files...)
	return err
}

// parseAlterTablespaceEncryption reads alter_tablespace_encryption:
//
//	ENCRYPTION { ONLINE [ tablespace_encryption_spec ] { ENCRYPT | DECRYPT } [ ts_file_name_convert ]
//	           | OFFLINE [ tablespace_encryption_spec ] { ENCRYPT | DECRYPT } [ tablespace_encryption_spec ]
//	           | FINISH { ENCRYPT | DECRYPT } [ ts_file_name_convert ]
//	           | [ tablespace_encryption_spec ] { ENCRYPT | DECRYPT } [ ts_file_name_convert ] }
//
// Oracle wants ENCRYPT or DECRYPT after each mode (ENCRYPTION ONLINE
// FILE_NAME_CONVERT ... is ORA-02142), takes the algorithm after OFFLINE's
// ENCRYPT or DECRYPT too, and takes the clause with no mode.
func (p *Parser) parseAlterTablespaceEncryption() error {
	p.advance() // ENCRYPTION
	mode := ""
	switch {
	case p.cur.Type == kwONLINE, p.cur.Type == kwOFFLINE, p.isKeywordStr("FINISH"):
		mode = p.advance().Str
	}
	using := false
	if mode != "FINISH" && p.cur.Type == kwUSING {
		if err := p.parseTablespaceEncryptionSpec(); err != nil {
			return err
		}
		using = true
	}
	if !p.isKeywordStr("ENCRYPT") && !p.isKeywordStr("DECRYPT") {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	if mode == "OFFLINE" {
		if !using && p.cur.Type == kwUSING {
			return p.parseTablespaceEncryptionSpec()
		}
		if p.isKeywordStr("FILE_NAME_CONVERT") {
			// Oracle reads FILE_NAME_CONVERT here and refuses what
			// follows it (ORA-02142): OFFLINE converts no file names.
			p.advance()
			return p.syntaxErrorAtCur()
		}
		return nil
	}
	if p.isKeywordStr("FILE_NAME_CONVERT") {
		return p.parseTablespaceFileNameConvert()
	}
	return nil
}

// parseTablespaceEncryptionSpec reads tablespace_encryption_spec:
// USING 'encrypt_algorithm' (ORA-28339).
func (p *Parser) parseTablespaceEncryptionSpec() error {
	p.advance() // USING
	if p.cur.Type != tokSCONST || !isValidOracleEncryptionAlgorithm(p.cur.Str) {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// parseTablespaceFileNameConvert reads ts_file_name_convert:
//
//	FILE_NAME_CONVERT = ( 'pattern', 'replacement' [, 'pattern', 'replacement' ]... ) [ KEEP ]
func (p *Parser) parseTablespaceFileNameConvert() error {
	p.advance() // FILE_NAME_CONVERT
	if p.cur.Type != '=' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Type != '(' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Type != tokSCONST {
		return p.syntaxErrorAtCur()
	}
	for {
		if p.cur.Type != tokSCONST {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		if p.cur.Type != ',' {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		if p.cur.Type != tokSCONST {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if p.cur.Type != ')' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	if p.isKeywordStr("KEEP") {
		p.advance()
	}
	return nil
}

// ---------------------------------------------------------------------------
// DROP TABLESPACE
// ---------------------------------------------------------------------------

// parseDropTablespaceStmt reads DROP TABLESPACE [ SET ] after TABLESPACE or
// SET.
//
// BNF: oracle/parser/bnf/DROP-TABLESPACE.bnf, DROP-TABLESPACE-SET.bnf
//
//	DROP TABLESPACE [ IF EXISTS ] tablespace
//	    [ { DROP | KEEP } QUOTA ]
//	    [ INCLUDING CONTENTS [ { AND DATAFILES | KEEP DATAFILES } ] [ CASCADE CONSTRAINTS ] ]
//	DROP TABLESPACE SET tablespace_set
//	    [ INCLUDING CONTENTS [ { AND DATAFILES | KEEP DATAFILES } ] [ CASCADE CONSTRAINTS ] ]
//
// Oracle takes the QUOTA clause after INCLUDING CONTENTS too; anything else,
// or a clause twice, is ORA-02173.
func (p *Parser) parseDropTablespaceStmt(start int, isSet bool) (*nodes.DropTablespaceStmt, error) {
	stmt := &nodes.DropTablespaceStmt{IsSet: isSet, Loc: nodes.Loc{Start: start}}
	if err := p.tsFinish(&stmt.Loc, p.parseDropTablespaceBody(stmt)); err != nil {
		return nil, err
	}
	return stmt, nil
}

func (p *Parser) parseDropTablespaceBody(stmt *nodes.DropTablespaceStmt) error {
	stmt.IfExists = p.tsIfExists()
	if stmt.IsSet && !stmt.IfExists && p.tsEndsBeforeOperand() {
		return nil
	}
	var err error
	if stmt.Name, err = p.parseTablespaceName(); err != nil {
		return err
	}
	for !p.tsAtEnd() {
		switch {
		case p.isKeywordStr("INCLUDING") && !stmt.IncludingContents:
			p.advance()
			if err := p.tsExpect("CONTENTS"); err != nil { // ORA-01911
				return err
			}
			stmt.IncludingContents = true
			switch {
			case p.cur.Type == kwAND:
				p.advance()
				if err := p.tsExpect("DATAFILES"); err != nil {
					return err
				}
				stmt.AndDatafiles = true
			case p.isKeywordStr("KEEP") && p.tsNextIs("DATAFILES"):
				p.advance()
				p.advance()
				stmt.KeepDatafiles = true
			}
			if p.cur.Type == kwCASCADE {
				p.advance()
				if err := p.tsExpect("CONSTRAINTS"); err != nil {
					return err
				}
				stmt.CascadeConstraints = true
			}
		case !stmt.IsSet && !stmt.DropQuota && !stmt.KeepQuota && (p.cur.Type == kwDROP || p.isKeywordStr("KEEP")):
			drop := p.advance().Type == kwDROP
			if err := p.tsExpect("QUOTA"); err != nil {
				return err
			}
			stmt.DropQuota, stmt.KeepQuota = drop, !drop
		default:
			return p.syntaxErrorAtCur()
		}
	}
	return nil
}
