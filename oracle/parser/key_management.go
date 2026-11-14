package parser

import (
	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseAdministerKeyManagementStmt parses an ADMINISTER KEY MANAGEMENT
// statement. The current token is ADMINISTER.
//
// BNF: oracle/parser/bnf/ADMINISTER-KEY-MANAGEMENT.bnf
// Ref: https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ADMINISTER-KEY-MANAGEMENT.html
//
// The page's syntax diagrams quote the values its BNF conversion lost
// ('keystore_location', 'key_id', 'tag', 'secret', 'client_identifier',
// 'filename', 'encrypt_algorithm', 'backup_identifier'), and Oracle Free
// 23.26.3 takes no other token there. Where the page and the server differ,
// the server decides: FORCE KEYSTORE leads ALTER KEYSTORE PASSWORD, WITH
// IDENTIFIER precedes WITH BACKUP in MOVE KEYS and is optional, the
// MIGRATE forms follow SET KEY with no master key, algorithm or leading
// FORCE KEYSTORE, and USE KEY takes them too.
//
//	ADMINISTER KEY MANAGEMENT
//	  { CREATE KEYSTORE 'keystore_location' IDENTIFIED BY keystore_password
//	  | CREATE [ LOCAL ] AUTO_LOGIN KEYSTORE FROM KEYSTORE 'keystore_location'
//	        IDENTIFIED BY keystore_password
//	  | SET KEYSTORE OPEN [ FORCE KEYSTORE ] IDENTIFIED BY auth [ container ]
//	  | SET KEYSTORE CLOSE [ IDENTIFIED BY auth ] [ container ]
//	  | BACKUP KEYSTORE [ USING backup_identifier ] [ FORCE KEYSTORE ]
//	        IDENTIFIED BY auth [ TO 'keystore_location' ]
//	  | ALTER KEYSTORE PASSWORD [ FORCE KEYSTORE ] IDENTIFIED BY old_password
//	        SET new_password [ with_backup ]
//	  | MERGE KEYSTORE 'keystore1_location' [ IDENTIFIED BY keystore1_password ]
//	        { AND KEYSTORE 'keystore2_location' [ IDENTIFIED BY keystore2_password ]
//	          INTO NEW KEYSTORE 'keystore3_location' IDENTIFIED BY keystore3_password
//	        | INTO EXISTING KEYSTORE 'keystore2_location'
//	          IDENTIFIED BY keystore2_password [ with_backup ] }
//	  | [ FORCE ] ISOLATE KEYSTORE IDENTIFIED BY isolated_keystore_password
//	        FROM ROOT KEYSTORE [ FORCE KEYSTORE ] IDENTIFIED BY auth [ with_backup ]
//	  | UNITE KEYSTORE IDENTIFIED BY isolated_keystore_password
//	        WITH ROOT KEYSTORE [ FORCE KEYSTORE ] IDENTIFIED BY auth [ with_backup ]
//	  | { SET | CREATE } [ ENCRYPTION ] KEY [ 'mkid:mk' ] [ USING TAG 'tag' ]
//	        [ USING ALGORITHM 'encrypt_algorithm' ] [ FORCE KEYSTORE ]
//	        IDENTIFIED BY auth [ with_backup ] [ container ]
//	  | SET [ ENCRYPTION ] KEY [ USING TAG 'tag' ] IDENTIFIED BY HSM_auth_string
//	        { [ FORCE KEYSTORE ] MIGRATE USING software_keystore_password [ with_backup ]
//	        | [ FORCE KEYSTORE ] REVERSE MIGRATE USING HSM_auth_string [ with_backup ] }
//	  | USE [ ENCRYPTION ] KEY 'key_id' [ USING TAG 'tag' ] [ FORCE KEYSTORE ]
//	        IDENTIFIED BY auth [ with_backup ]
//	  | USE [ ENCRYPTION ] KEY 'key_id' [ USING TAG 'tag' ] IDENTIFIED BY HSM_auth_string
//	        [ FORCE KEYSTORE ] MIGRATE USING software_keystore_password [ with_backup ]
//	  | SET TAG 'tag' FOR 'key_id' [ FORCE KEYSTORE ] IDENTIFIED BY auth [ with_backup ]
//	  | EXPORT [ ENCRYPTION ] KEYS WITH SECRET secret TO 'filename'
//	        [ FORCE KEYSTORE ] IDENTIFIED BY auth [ with_identifier ]
//	  | IMPORT [ ENCRYPTION ] KEYS WITH SECRET secret FROM 'filename'
//	        [ FORCE KEYSTORE ] IDENTIFIED BY auth [ with_backup ]
//	  | MOVE [ ENCRYPTION ] KEYS TO NEW KEYSTORE 'keystore_location1'
//	        IDENTIFIED BY keystore1_password FROM [ FORCE ] KEYSTORE
//	        IDENTIFIED BY auth [ with_identifier ] [ with_backup ]
//	  | { ADD | UPDATE } SECRET 'secret' FOR CLIENT 'client_identifier'
//	        [ USING TAG 'tag' ] [ FORCE KEYSTORE ] IDENTIFIED BY auth [ with_backup ]
//	  | DELETE SECRET FOR CLIENT 'client_identifier' [ FORCE KEYSTORE ]
//	        IDENTIFIED BY auth [ with_backup ]
//	  | SWITCHOVER TO LIBRARY 'library_path' FOR ALL CONTAINERS }
//
//	auth:            { EXTERNAL STORE | keystore_password }
//	with_backup:     WITH BACKUP [ USING backup_identifier ]
//	container:       CONTAINER = { CURRENT | ALL }
//	with_identifier: WITH IDENTIFIER IN { key_id [, key_id ]... | ( subquery ) }
func (p *Parser) parseAdministerKeyManagementStmt() (nodes.StmtNode, error) {
	start := p.pos()
	p.advance() // ADMINISTER

	stmt := &nodes.AdminDDLStmt{
		Action:     "ADMINISTER",
		ObjectType: nodes.OBJECT_KEY_MANAGEMENT,
		Loc:        nodes.Loc{Start: start},
	}
	k := &keyManagementClauses{p: p, opts: &nodes.List{}}
	scopeStart := p.pos()
	if err := k.words("KEY", "MANAGEMENT"); err != nil {
		return nil, err
	}
	k.add("SCOPE", "KEY MANAGEMENT", scopeStart)
	if err := k.parseCommand(); err != nil {
		return nil, err
	}
	stmt.Options = k.opts
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// keyManagementClauses reads the clauses of one ADMINISTER KEY MANAGEMENT
// statement into its options, one *DDLOption per clause.
type keyManagementClauses struct {
	p    *Parser
	opts *nodes.List
}

func (k *keyManagementClauses) add(key, value string, start int) {
	k.opts.Items = append(k.opts.Items, &nodes.DDLOption{
		Key:   key,
		Value: value,
		Loc:   nodes.Loc{Start: start, End: k.p.prev.End},
	})
}

// word reports whether the current token is the unquoted word s.
func (k *keyManagementClauses) word(s string) bool {
	return k.p.isKeywordStr(s)
}

// nextWord reports whether the token after the current one is the
// unquoted word s.
func (k *keyManagementClauses) nextWord(s string) bool {
	return k.p.isIdentLikeStrAt(k.p.peekNext(), s)
}

// words consumes the unquoted words ws in order.
func (k *keyManagementClauses) words(ws ...string) error {
	for _, w := range ws {
		if !k.word(w) {
			return k.p.syntaxErrorAtCur()
		}
		k.p.advance()
	}
	return nil
}

// parseCommand parses the clauses after ADMINISTER KEY MANAGEMENT.
func (k *keyManagementClauses) parseCommand() error {
	p := k.p
	start := p.pos()
	switch {
	case p.cur.Type == kwCREATE:
		return k.parseCreate(start)
	case p.cur.Type == kwSET:
		return k.parseSet(start)
	case k.word("USE"):
		return k.parseUseKey(start)
	case k.word("BACKUP"):
		return k.parseBackupKeystore(start)
	case p.cur.Type == kwALTER:
		return k.parseAlterKeystorePassword(start)
	case k.word("MERGE"):
		return k.parseMergeKeystore(start)
	case k.word("FORCE"), k.word("ISOLATE"):
		return k.parseIsolateKeystore(start)
	case k.word("UNITE"):
		return k.parseUniteKeystore(start)
	case k.word("EXPORT"), k.word("IMPORT"):
		return k.parseExportImportKeys(start)
	case k.word("MOVE"):
		return k.parseMoveKeys(start)
	case k.word("ADD"), k.word("UPDATE"):
		return k.parseAddUpdateSecret(start)
	case k.word("DELETE"):
		return k.parseDeleteSecret(start)
	case k.word("SWITCHOVER"):
		return k.parseSwitchover(start)
	}
	return p.syntaxErrorAtCur()
}

// command consumes the words of a command and records them.
func (k *keyManagementClauses) command(name string, start int, ws ...string) error {
	if err := k.words(ws...); err != nil {
		return err
	}
	k.add("COMMAND", name, start)
	return nil
}

// errorAtCommand reports the word after SET or CREATE that leads no
// clause at the command's first word, where Oracle reports ORA-46601,
// 46602 and 46603.
func (k *keyManagementClauses) errorAtCommand(start int) error {
	if k.p.lexer.Err != nil {
		return k.p.lexerError()
	}
	return &ParseError{Message: "syntax error: invalid ADMINISTER KEY MANAGEMENT command", Position: start}
}

// parseCreate parses CREATE KEYSTORE, CREATE [ LOCAL ] AUTO_LOGIN KEYSTORE
// and CREATE [ ENCRYPTION ] KEY.
func (k *keyManagementClauses) parseCreate(start int) error {
	p := k.p
	switch next := p.peekNext(); {
	case p.isIdentLikeStrAt(next, "KEYSTORE"):
		if err := k.command("CREATE KEYSTORE", start, "CREATE", "KEYSTORE"); err != nil {
			return err
		}
		if err := k.stringValue("KEYSTORE", p.pos()); err != nil {
			return err
		}
		_, err := k.identifiedBy(false)
		return err
	case p.isIdentLikeStrAt(next, "LOCAL"), p.isIdentLikeStrAt(next, "AUTO_LOGIN"):
		p.advance() // CREATE
		name := "CREATE AUTO_LOGIN KEYSTORE"
		if k.word("LOCAL") {
			p.advance()
			name = "CREATE LOCAL AUTO_LOGIN KEYSTORE"
		}
		if err := k.command(name, start, "AUTO_LOGIN", "KEYSTORE"); err != nil {
			return err
		}
		locStart := p.pos()
		if err := k.words("FROM", "KEYSTORE"); err != nil {
			return err
		}
		if err := k.stringValue("FROM KEYSTORE", locStart); err != nil {
			return err
		}
		_, err := k.identifiedBy(false)
		return err
	case p.isIdentLikeStrAt(next, "ENCRYPTION"), p.isIdentLikeStrAt(next, "KEY"):
		p.advance() // CREATE
		return k.parseKey(start, "CREATE KEY")
	}
	return k.errorAtCommand(start)
}

// parseSet parses SET KEYSTORE { OPEN | CLOSE }, SET [ ENCRYPTION ] KEY and
// SET TAG.
func (k *keyManagementClauses) parseSet(start int) error {
	p := k.p
	switch next := p.peekNext(); {
	case p.isIdentLikeStrAt(next, "KEYSTORE"):
		p.advance() // SET
		p.advance() // KEYSTORE
		// compat.tsv records SET KEYSTORE ending the statement as
		// accepted, so omni takes it.
		if p.cur.Type == tokEOF || p.cur.Type == ';' {
			k.add("COMMAND", "SET KEYSTORE", start)
			return nil
		}
		switch {
		case k.word("OPEN"):
			p.advance()
			k.add("COMMAND", "SET KEYSTORE OPEN", start)
			if _, err := k.forceKeystore(); err != nil {
				return err
			}
			if _, err := k.identifiedBy(true); err != nil {
				return err
			}
		case k.word("CLOSE"):
			p.advance()
			k.add("COMMAND", "SET KEYSTORE CLOSE", start)
			if p.cur.Type == kwIDENTIFIED {
				if _, err := k.identifiedBy(true); err != nil {
					return err
				}
			}
		default:
			return k.errorAtCommand(start)
		}
		return k.container()
	case p.isIdentLikeStrAt(next, "ENCRYPTION"), p.isIdentLikeStrAt(next, "KEY"):
		p.advance() // SET
		return k.parseKey(start, "SET KEY")
	case p.isIdentLikeStrAt(next, "TAG"):
		if err := k.command("SET TAG", start, "SET", "TAG"); err != nil {
			return err
		}
		if err := k.stringValue("TAG", p.prev.Loc); err != nil {
			return err
		}
		forStart := p.pos()
		if err := k.words("FOR"); err != nil {
			return err
		}
		if err := k.stringValue("FOR", forStart); err != nil {
			return err
		}
		return k.keystoreAuth(true)
	}
	return k.errorAtCommand(start)
}

// parseKey parses SET or CREATE [ ENCRYPTION ] KEY and what follows it.
// The current token is ENCRYPTION or KEY.
func (k *keyManagementClauses) parseKey(start int, name string) error {
	p := k.p
	if k.word("ENCRYPTION") {
		p.advance()
	}
	if err := k.command(name, start, "KEY"); err != nil {
		return err
	}
	// Oracle takes the MIGRATE forms after SET KEY [ USING TAG 'tag' ] alone.
	migrate := name == "SET KEY"
	if p.cur.Type == tokSCONST {
		if err := k.stringValue("KEY", p.pos()); err != nil {
			return err
		}
		migrate = false
	}
	if p.cur.Type == kwUSING && k.nextWord("TAG") {
		if err := k.usingTag(); err != nil {
			return err
		}
	}
	if p.cur.Type == kwUSING && k.nextWord("ALGORITHM") {
		algStart := p.pos()
		p.advance() // USING
		p.advance() // ALGORITHM
		if err := k.stringValue("USING ALGORITHM", algStart); err != nil {
			return err
		}
		migrate = false
	}
	if force, err := k.forceKeystore(); err != nil {
		return err
	} else if force {
		migrate = false
	}
	store, err := k.identifiedBy(true)
	if err != nil {
		return err
	}
	if migrate && !store && k.isMigrateStart(true) {
		return k.migrate(true)
	}
	if err := k.withBackup(); err != nil {
		return err
	}
	return k.container()
}

// parseUseKey parses USE [ ENCRYPTION ] KEY.
func (k *keyManagementClauses) parseUseKey(start int) error {
	p := k.p
	p.advance() // USE
	if k.word("ENCRYPTION") {
		p.advance()
	}
	if err := k.command("USE KEY", start, "KEY"); err != nil {
		return err
	}
	if err := k.stringValue("KEY", p.pos()); err != nil {
		return err
	}
	if p.cur.Type == kwUSING {
		if err := k.usingTag(); err != nil {
			return err
		}
	}
	force, err := k.forceKeystore()
	if err != nil {
		return err
	}
	store, err := k.identifiedBy(true)
	if err != nil {
		return err
	}
	if !force && !store && k.isMigrateStart(false) {
		return k.migrate(false)
	}
	return k.withBackup()
}

// isMigrateStart reports whether the current token starts the MIGRATE
// tail of SET KEY or USE KEY: [ FORCE KEYSTORE ] [ REVERSE ] MIGRATE.
// REVERSE MIGRATE follows SET KEY alone.
func (k *keyManagementClauses) isMigrateStart(reverse bool) bool {
	p := k.p
	return p.cur.Type == kwFORCE || k.word("MIGRATE") ||
		reverse && k.word("REVERSE") && k.nextWord("MIGRATE")
}

// migrate parses [ FORCE KEYSTORE ] [ REVERSE ] MIGRATE USING password
// [ WITH BACKUP [ USING backup_identifier ] ].
func (k *keyManagementClauses) migrate(reverse bool) error {
	p := k.p
	if _, err := k.forceKeystore(); err != nil {
		return err
	}
	optStart := p.pos()
	key := "MIGRATE USING"
	if reverse && k.word("REVERSE") {
		p.advance()
		key = "REVERSE MIGRATE USING"
	}
	if err := k.words("MIGRATE", "USING"); err != nil {
		return err
	}
	value, err := k.password(false)
	if err != nil {
		return err
	}
	k.add(key, value, optStart)
	return k.withBackup()
}

// parseBackupKeystore parses BACKUP KEYSTORE.
func (k *keyManagementClauses) parseBackupKeystore(start int) error {
	p := k.p
	if err := k.command("BACKUP KEYSTORE", start, "BACKUP", "KEYSTORE"); err != nil {
		return err
	}
	if p.cur.Type == kwUSING {
		if err := k.backupIdentifier(); err != nil {
			return err
		}
	}
	if err := k.keystoreAuthOnly(true); err != nil {
		return err
	}
	if p.cur.Type == kwTO {
		toStart := p.pos()
		p.advance()
		return k.stringValue("TO", toStart)
	}
	return nil
}

// parseAlterKeystorePassword parses ALTER KEYSTORE PASSWORD. Oracle takes
// FORCE KEYSTORE before IDENTIFIED BY, not after the new password where
// the BNF shows it.
func (k *keyManagementClauses) parseAlterKeystorePassword(start int) error {
	p := k.p
	if err := k.command("ALTER KEYSTORE PASSWORD", start, "ALTER", "KEYSTORE", "PASSWORD"); err != nil {
		return err
	}
	if err := k.keystoreAuthOnly(false); err != nil {
		return err
	}
	setStart := p.pos()
	if p.cur.Type != kwSET {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	value, err := k.password(false)
	if err != nil {
		return err
	}
	k.add("SET", value, setStart)
	return k.withBackup()
}

// parseMergeKeystore parses MERGE KEYSTORE into a new or an existing
// keystore.
func (k *keyManagementClauses) parseMergeKeystore(start int) error {
	p := k.p
	if err := k.command("MERGE KEYSTORE", start, "MERGE", "KEYSTORE"); err != nil {
		return err
	}
	if err := k.stringValue("KEYSTORE", p.prev.Loc); err != nil {
		return err
	}
	if err := k.optionalIdentifiedBy(); err != nil {
		return err
	}
	if p.cur.Type == kwAND {
		andStart := p.pos()
		if err := k.words("AND", "KEYSTORE"); err != nil {
			return err
		}
		if err := k.stringValue("AND KEYSTORE", andStart); err != nil {
			return err
		}
		if err := k.optionalIdentifiedBy(); err != nil {
			return err
		}
		intoStart := p.pos()
		if err := k.words("INTO", "NEW", "KEYSTORE"); err != nil {
			return err
		}
		if err := k.stringValue("INTO NEW KEYSTORE", intoStart); err != nil {
			return err
		}
		_, err := k.identifiedBy(false)
		return err
	}
	intoStart := p.pos()
	if err := k.words("INTO", "EXISTING", "KEYSTORE"); err != nil {
		return err
	}
	if err := k.stringValue("INTO EXISTING KEYSTORE", intoStart); err != nil {
		return err
	}
	if _, err := k.identifiedBy(false); err != nil {
		return err
	}
	return k.withBackup()
}

// optionalIdentifiedBy parses the [ IDENTIFIED BY password ] after a
// keystore that MERGE KEYSTORE reads.
func (k *keyManagementClauses) optionalIdentifiedBy() error {
	if k.p.cur.Type != kwIDENTIFIED {
		return nil
	}
	_, err := k.identifiedBy(false)
	return err
}

// parseIsolateKeystore parses [ FORCE ] ISOLATE KEYSTORE.
func (k *keyManagementClauses) parseIsolateKeystore(start int) error {
	p := k.p
	name := "ISOLATE KEYSTORE"
	if p.cur.Type == kwFORCE {
		p.advance()
		name = "FORCE ISOLATE KEYSTORE"
	}
	if err := k.command(name, start, "ISOLATE", "KEYSTORE"); err != nil {
		return err
	}
	if _, err := k.identifiedBy(false); err != nil {
		return err
	}
	rootStart := p.pos()
	if err := k.words("FROM", "ROOT", "KEYSTORE"); err != nil {
		return err
	}
	k.add("FROM ROOT KEYSTORE", "", rootStart)
	return k.keystoreAuth(true)
}

// parseUniteKeystore parses UNITE KEYSTORE.
func (k *keyManagementClauses) parseUniteKeystore(start int) error {
	p := k.p
	if err := k.command("UNITE KEYSTORE", start, "UNITE", "KEYSTORE"); err != nil {
		return err
	}
	if _, err := k.identifiedBy(false); err != nil {
		return err
	}
	rootStart := p.pos()
	if err := k.words("WITH", "ROOT", "KEYSTORE"); err != nil {
		return err
	}
	k.add("WITH ROOT KEYSTORE", "", rootStart)
	return k.keystoreAuth(true)
}

// parseExportImportKeys parses EXPORT and IMPORT [ ENCRYPTION ] KEYS.
func (k *keyManagementClauses) parseExportImportKeys(start int) error {
	p := k.p
	export := k.word("EXPORT")
	p.advance() // EXPORT or IMPORT
	if k.word("ENCRYPTION") {
		p.advance()
	}
	name, file := "IMPORT KEYS", "FROM"
	if export {
		name, file = "EXPORT KEYS", "TO"
	}
	if err := k.command(name, start, "KEYS"); err != nil {
		return err
	}
	secretStart := p.pos()
	if err := k.words("WITH", "SECRET"); err != nil {
		return err
	}
	// The secret of EXPORT and IMPORT is an identifier: Oracle refuses a
	// string literal there (ORA-46609).
	secret, err := k.password(false)
	if err != nil {
		return err
	}
	k.add("WITH SECRET", secret, secretStart)
	fileStart := p.pos()
	if err := k.words(file); err != nil {
		return err
	}
	if err := k.stringValue(file, fileStart); err != nil {
		return err
	}
	if !export {
		return k.keystoreAuth(true)
	}
	if err := k.keystoreAuthOnly(true); err != nil {
		return err
	}
	if p.cur.Type == kwWITH {
		_, err := k.withIdentifier(true)
		return err
	}
	return nil
}

// parseMoveKeys parses MOVE [ ENCRYPTION ] KEYS.
func (k *keyManagementClauses) parseMoveKeys(start int) error {
	p := k.p
	p.advance() // MOVE
	if k.word("ENCRYPTION") {
		p.advance()
	}
	if err := k.command("MOVE KEYS", start, "KEYS"); err != nil {
		return err
	}
	toStart := p.pos()
	if err := k.words("TO", "NEW", "KEYSTORE"); err != nil {
		return err
	}
	if err := k.stringValue("TO NEW KEYSTORE", toStart); err != nil {
		return err
	}
	if _, err := k.identifiedBy(false); err != nil {
		return err
	}
	fromStart := p.pos()
	if err := k.words("FROM"); err != nil {
		return err
	}
	key := "FROM KEYSTORE"
	if p.cur.Type == kwFORCE {
		p.advance()
		key = "FROM FORCE KEYSTORE"
	}
	if err := k.words("KEYSTORE"); err != nil {
		return err
	}
	k.add(key, "", fromStart)
	if _, err := k.identifiedBy(true); err != nil {
		return err
	}
	if p.cur.Type == kwWITH && k.nextWord("IDENTIFIER") {
		subquery, err := k.withIdentifier(false)
		if err != nil {
			return err
		}
		// Oracle 23.26.3 takes a subquery here only before WITH BACKUP: one
		// that ends the statement, or is followed by anything else, it
		// answers with ORA-00921 at the end of the text.
		if subquery && p.cur.Type != kwWITH {
			return p.syntaxErrorAtCur()
		}
	}
	return k.withBackup()
}

// parseAddUpdateSecret parses { ADD | UPDATE } SECRET.
func (k *keyManagementClauses) parseAddUpdateSecret(start int) error {
	p := k.p
	verb := p.cur.Str
	if err := k.command(verb+" SECRET", start, verb, "SECRET"); err != nil {
		return err
	}
	if err := k.stringValue("SECRET", p.prev.Loc); err != nil {
		return err
	}
	if err := k.forClient(); err != nil {
		return err
	}
	if p.cur.Type == kwUSING {
		if err := k.usingTag(); err != nil {
			return err
		}
	}
	return k.keystoreAuth(true)
}

// parseDeleteSecret parses DELETE SECRET.
func (k *keyManagementClauses) parseDeleteSecret(start int) error {
	if err := k.command("DELETE SECRET", start, "DELETE", "SECRET"); err != nil {
		return err
	}
	if err := k.forClient(); err != nil {
		return err
	}
	return k.keystoreAuth(true)
}

// forClient parses FOR CLIENT 'client_identifier'.
func (k *keyManagementClauses) forClient() error {
	forStart := k.p.pos()
	if err := k.words("FOR", "CLIENT"); err != nil {
		return err
	}
	return k.stringValue("FOR CLIENT", forStart)
}

// parseSwitchover parses SWITCHOVER TO LIBRARY 'library_path' FOR ALL
// CONTAINERS.
func (k *keyManagementClauses) parseSwitchover(start int) error {
	p := k.p
	if err := k.command("SWITCHOVER TO LIBRARY", start, "SWITCHOVER", "TO", "LIBRARY"); err != nil {
		return err
	}
	if err := k.stringValue("LIBRARY", p.prev.Loc); err != nil {
		return err
	}
	forStart := p.pos()
	if err := k.words("FOR", "ALL", "CONTAINERS"); err != nil {
		return err
	}
	k.add("FOR ALL CONTAINERS", "", forStart)
	return nil
}

// keystoreAuth parses [ FORCE KEYSTORE ] IDENTIFIED BY auth [ WITH BACKUP
// [ USING backup_identifier ] ], the tail most clauses share.
func (k *keyManagementClauses) keystoreAuth(external bool) error {
	if err := k.keystoreAuthOnly(external); err != nil {
		return err
	}
	return k.withBackup()
}

// keystoreAuthOnly parses [ FORCE KEYSTORE ] IDENTIFIED BY auth.
func (k *keyManagementClauses) keystoreAuthOnly(external bool) error {
	if _, err := k.forceKeystore(); err != nil {
		return err
	}
	_, err := k.identifiedBy(external)
	return err
}

// forceKeystore parses an optional FORCE KEYSTORE and reports whether it
// was there.
func (k *keyManagementClauses) forceKeystore() (bool, error) {
	p := k.p
	if p.cur.Type != kwFORCE {
		return false, nil
	}
	start := p.pos()
	p.advance()
	if err := k.words("KEYSTORE"); err != nil {
		return false, err
	}
	k.add("FORCE KEYSTORE", "", start)
	return true, nil
}

// identifiedBy parses IDENTIFIED BY and a password, or EXTERNAL STORE where
// external says the clause takes it, and reports whether it read EXTERNAL
// STORE: a password "EXTERNAL STORE" records the same text.
func (k *keyManagementClauses) identifiedBy(external bool) (bool, error) {
	p := k.p
	start := p.pos()
	if p.cur.Type != kwIDENTIFIED {
		return false, p.syntaxErrorAtCur()
	}
	p.advance()
	if p.cur.Type != kwBY {
		return false, p.syntaxErrorAtCur()
	}
	p.advance()
	store := external && k.word("EXTERNAL")
	value, err := k.password(external)
	if err != nil {
		return false, err
	}
	k.add("IDENTIFIED BY", value, start)
	return store, nil
}

// password parses a keystore password, an identifier Oracle reads as
// written: a string literal, a number or a reserved word is no password
// (ORA-00988). Where external is set, EXTERNAL starts EXTERNAL STORE;
// elsewhere it is a password like any other word.
func (k *keyManagementClauses) password(external bool) (string, error) {
	p := k.p
	if external && k.word("EXTERNAL") {
		p.advance()
		if err := k.words("STORE"); err != nil {
			return "", err
		}
		return "EXTERNAL STORE", nil
	}
	if !k.isName() {
		return "", p.syntaxErrorAtCur()
	}
	return p.advance().Str, nil
}

// isName reports whether the current token is a name: a quoted identifier,
// or a word that is not reserved.
func (k *keyManagementClauses) isName() bool {
	p := k.p
	return p.cur.Type == tokQIDENT || p.isIdentLike() && !isOracleSQLReservedKeyword(p.cur)
}

// stringValue parses the string literal of a clause that started at
// start: a keystore location, key identifier, tag, secret, client
// identifier, file name or algorithm.
func (k *keyManagementClauses) stringValue(key string, start int) error {
	p := k.p
	if p.cur.Type != tokSCONST {
		return p.syntaxErrorAtCur()
	}
	k.add(key, p.advance().Str, start)
	return nil
}

// usingTag parses USING TAG 'tag'.
func (k *keyManagementClauses) usingTag() error {
	start := k.p.pos()
	if err := k.words("USING", "TAG"); err != nil {
		return err
	}
	return k.stringValue("USING TAG", start)
}

// withBackup parses an optional WITH BACKUP [ USING backup_identifier ].
func (k *keyManagementClauses) withBackup() error {
	p := k.p
	if p.cur.Type != kwWITH {
		return nil
	}
	start := p.pos()
	if err := k.words("WITH", "BACKUP"); err != nil {
		return err
	}
	k.add("WITH BACKUP", "", start)
	if p.cur.Type == kwUSING {
		return k.backupIdentifier()
	}
	return nil
}

// backupIdentifier parses USING backup_identifier: a string literal, a
// national one, or a quoted identifier; Oracle refuses an unquoted word
// (ORA-00931).
func (k *keyManagementClauses) backupIdentifier() error {
	p := k.p
	start := p.pos()
	p.advance() // USING
	switch p.cur.Type {
	case tokSCONST, tokNCHARLIT, tokQIDENT:
		k.add("USING", p.advance().Str, start)
		return nil
	}
	return p.syntaxErrorAtCur()
}

// container parses an optional CONTAINER = { CURRENT | ALL }.
func (k *keyManagementClauses) container() error {
	p := k.p
	if !k.word("CONTAINER") {
		return nil
	}
	start := p.pos()
	p.advance()
	if p.cur.Type != '=' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	if !k.word("CURRENT") && !k.word("ALL") {
		return p.syntaxErrorAtCur()
	}
	k.add("CONTAINER", p.advance().Str, start)
	return nil
}

// withIdentifier parses WITH IDENTIFIER IN { key_id [, key_id ]... |
// ( subquery ) }, and reports whether it read a subquery. A key identifier
// is a string literal or a word.
//
// After EXPORT KEYS (deferred), Oracle 23.26.3 does not parse what follows
// the parenthesis: it takes IN(a zzz), IN() and IN( alike and reads the text
// only when it runs the export. omni keeps that text, from the parenthesis
// to the end of the statement, as the option's value.
func (k *keyManagementClauses) withIdentifier(deferred bool) (bool, error) {
	p := k.p
	start := p.pos()
	if err := k.words("WITH", "IDENTIFIER"); err != nil {
		return false, err
	}
	if p.cur.Type != kwIN {
		return false, p.syntaxErrorAtCur()
	}
	p.advance()
	items := &nodes.List{}
	subquery := p.cur.Type == '('
	if subquery && deferred {
		textStart := p.pos()
		for p.cur.Type != tokEOF && p.cur.Type != ';' {
			p.advance()
		}
		k.add("WITH IDENTIFIER IN", p.source[textStart:p.prev.End], start)
		return true, nil
	}
	if subquery {
		p.advance()
		if !p.atQuery() {
			return false, p.syntaxErrorAtCur()
		}
		sub, err := p.parseSelectStmt()
		if err != nil {
			return false, err
		}
		if p.cur.Type != ')' {
			return false, p.syntaxErrorAtCur()
		}
		p.advance()
		items.Items = append(items.Items, sub)
	} else {
		for {
			if p.cur.Type != tokSCONST && !k.isName() {
				return false, p.syntaxErrorAtCur()
			}
			items.Items = append(items.Items, &nodes.String{Str: p.advance().Str})
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}
	k.opts.Items = append(k.opts.Items, &nodes.DDLOption{
		Key:   "WITH IDENTIFIER IN",
		Items: items,
		Loc:   nodes.Loc{Start: start, End: p.prev.End},
	})
	return subquery, nil
}
