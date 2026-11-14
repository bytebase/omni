package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// This file reads CREATE PLUGGABLE DATABASE and ALTER PLUGGABLE DATABASE by
// the grammar of
// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/CREATE-PLUGGABLE-DATABASE.html
// and
// https://docs.oracle.com/en/database/oracle/oracle-database/26/sqlrf/ALTER-PLUGGABLE-DATABASE.html
// (bnf/CREATE-PLUGGABLE-DATABASE.bnf, bnf/ALTER-PLUGGABLE-DATABASE.bnf),
// corrected where the pinned Oracle reads them otherwise: the compat rows and
// the probe decide. As in ALTER DATABASE (alter_database_clauses.go), a
// clause that stops where Oracle reports the missing operand with a code of
// the statement rather than of the parser ends the statement.

// pdbKind is the form of CREATE PLUGGABLE DATABASE: from the seed, a clone
// (FROM or AS PROXY FROM), or a plug-in (USING 'filename').
type pdbKind int

const (
	pdbFromSeed pdbKind = iota
	pdbClone
	pdbFromXML
)

// pdbName reads a PDB name: an identifier Oracle does not reserve
// (ORA-65000 for SELECT).
func (p *Parser) pdbName() (*nodes.ObjectName, error) {
	start := p.pos()
	name, err := p.adbName()
	if err != nil {
		return nil, err
	}
	return &nodes.ObjectName{Name: name, Loc: nodes.Loc{Start: start, End: p.prev.End}}, nil
}

// pdbOption appends the option the caller read from start.
func (p *Parser) pdbOption(opts *nodes.List, key, value string, start int) *nodes.DDLOption {
	opt := &nodes.DDLOption{Key: key, Value: value, Loc: nodes.Loc{Start: start, End: p.prev.End}}
	opts.Items = append(opts.Items, opt)
	return opt
}

// pdbExpectEquals consumes '=' or fails at the current token.
func (p *Parser) pdbExpectEquals() error {
	if p.cur.Type != '=' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// parseCreatePluggableDatabaseStmt parses CREATE PLUGGABLE DATABASE, CREATE
// PLUGGABLE DATABASE consumed.
//
//	CREATE PLUGGABLE DATABASE pdb_name [ AS APPLICATION CONTAINER ]
//	  { create_pdb_from_seed | create_pdb_clone | create_pdb_from_xml }
//	create_pdb_clone ::= { FROM src_pdb_name [ @ dblink ] | AS PROXY FROM src_pdb_name [ @ dblink ] }
//	create_pdb_from_xml ::= USING 'filename'
//
// Each form's clauses follow in any order, each at most once (ORA-12990).
// Oracle requires the seed form's ADMIN USER, reporting it missing at the
// end of the statement (ORA-00922), and does not take AS SEED (ORA-02000),
// pdb_snapshot_clause, HOST, PORT, or a DECRYPT USING or EXTERNAL STORE
// keystore clause in a clone.
func (p *Parser) parseCreatePluggableDatabaseStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "CREATE",
		ObjectType: nodes.OBJECT_PLUGGABLE_DATABASE,
		Loc:        nodes.Loc{Start: start},
	}
	opts := &nodes.List{}
	if p.adbEndKept() {
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}
	name, err := p.pdbName()
	if err != nil {
		return nil, err
	}
	stmt.Name = name
	kind := pdbFromSeed
	dblink := false
	if p.cur.Type == kwAS && p.isIdentLikeStrAt(p.peekNext(), "APPLICATION") {
		optStart := p.pos()
		p.advance()
		p.advance()
		if err := p.adbExpect("CONTAINER"); err != nil {
			return nil, err
		}
		p.pdbOption(opts, "AS_APPLICATION_CONTAINER", "", optStart)
	}
	optStart := p.pos()
	key := "FROM"
	switch {
	case p.cur.Type == kwAS:
		key = "AS_PROXY_FROM"
		p.advance()
		if p.adbWord("SEED") {
			p.advance()
			return nil, p.syntaxErrorAtCur()
		}
		if err := p.adbExpect("PROXY"); err != nil {
			return nil, err
		}
		if p.cur.Type != kwFROM {
			return nil, p.syntaxErrorAtCur()
		}
		fallthrough
	case p.cur.Type == kwFROM:
		p.advance()
		src, link, err := p.pdbSource()
		if err != nil {
			return nil, err
		}
		dblink = link != ""
		if link != "" {
			src += "@" + link
		}
		p.pdbOption(opts, key, src, optStart)
		kind = pdbClone
	case p.adbWord("USING"):
		p.advance()
		kind = pdbFromXML
		if !p.adbEndKept() {
			f, err := p.adbString()
			if err != nil {
				return nil, err
			}
			p.pdbOption(opts, "USING", f, optStart)
		}
	}
	tolerated, err := p.parseCreatePDBClauses(opts, kind, dblink)
	if err != nil {
		return nil, err
	}
	if kind == pdbFromSeed && !tolerated && !p.pdbHasOption(opts, "ADMIN_USER") {
		return nil, p.syntaxErrorAtCur()
	}
	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// pdbHasOption reports whether opts holds an option with key.
func (p *Parser) pdbHasOption(opts *nodes.List, key string) bool {
	for _, item := range opts.Items {
		if o, ok := item.(*nodes.DDLOption); ok && o.Key == key {
			return true
		}
	}
	return false
}

// pdbSource reads src_pdb_name [ @ dblink ], the dblink a dotted name.
// Oracle ends the statement without a complaint after a bare @.
func (p *Parser) pdbSource() (string, string, error) {
	src, err := p.adbName()
	if err != nil {
		return "", "", err
	}
	if p.cur.Type != '@' {
		return src, "", nil
	}
	p.advance()
	if p.adbEndKept() {
		return src, "", nil
	}
	link, err := p.adbDottedName()
	if err != nil {
		return "", "", err
	}
	return src, link, nil
}

// parseCreatePDBClauses reads the clauses of a CREATE PLUGGABLE DATABASE
// form, in any order, each at most once. It reports whether the statement
// ended where Oracle reports the missing operand with a code of the
// statement (ORA-02199, ORA-02216, ORA-02236, ORA-65001), which compat reads
// as accepted.
func (p *Parser) parseCreatePDBClauses(opts *nodes.List, kind pdbKind, dblink bool) (bool, error) {
	seen := map[string]bool{}
	for !p.adbEnd() {
		optStart := p.pos()
		key := p.cur.Str
		if p.cur.Type == tokQIDENT || !p.isIdentLike() {
			return false, p.syntaxErrorAtCur()
		}
		switch key {
		case "NOLOGGING":
			key = "LOGGING"
		case "SNAPSHOT", "NO":
			key = "COPY_DATA"
		}
		if seen[key] {
			return false, p.syntaxErrorAtCur()
		}
		seen[key] = true
		switch {
		case kind == pdbFromSeed && key == "ADMIN":
			p.advance()
			if err := p.adbExpect("USER"); err != nil {
				return false, err
			}
			user, err := p.adbName()
			if err != nil {
				return false, err
			}
			if p.cur.Type != kwIDENTIFIED {
				return false, p.syntaxErrorAtCur()
			}
			p.advance()
			if p.cur.Type != kwBY {
				return false, p.syntaxErrorAtCur()
			}
			p.advance()
			if err := p.pdbPassword(); err != nil {
				return false, err
			}
			opt := p.pdbOption(opts, "ADMIN_USER", user, optStart)
			if p.adbWord("ROLES") {
				rolesStart := p.pos()
				p.advance()
				if err := p.pdbExpectEquals(); err != nil {
					return false, err
				}
				roles, err := p.pdbParenRoles()
				if err != nil {
					return false, err
				}
				opt.Items = &nodes.List{Items: []nodes.Node{&nodes.DDLOption{
					Key: "ROLES", Value: strings.Join(roles, ", "),
					Loc: nodes.Loc{Start: rolesStart, End: p.prev.End},
				}}}
			}
		case key == "PARALLEL" || p.cur.Type == kwPARALLEL:
			p.advance()
			val := ""
			if p.cur.Type == tokICONST {
				val = p.advance().Str
			}
			p.pdbOption(opts, "PARALLEL", val, optStart)
		case p.cur.Type == kwDEFAULT:
			p.advance()
			if err := p.adbExpect("TABLESPACE"); err != nil {
				return false, err
			}
			if p.adbEndKept() {
				return true, nil
			}
			ts, err := p.adbName()
			if err != nil {
				return false, err
			}
			opt := p.pdbOption(opts, "DEFAULT_TABLESPACE", ts, optStart)
			if kind == pdbFromSeed {
				if p.adbEndKept() {
					return true, nil
				}
				if p.adbWord("DATAFILE") {
					p.advance()
					spec, err := p.parseAlterDatabaseDatafileSpec()
					if err != nil {
						return false, err
					}
					opt.Items = &nodes.List{Items: []nodes.Node{&nodes.DDLOption{Key: "DATAFILE", Value: spec}}}
					opt.Loc.End = p.prev.End
				}
			}
		case key == "FILE_NAME_CONVERT", key == "SERVICE_NAME_CONVERT",
			kind == pdbFromXML && key == "SOURCE_FILE_NAME_CONVERT":
			p.advance()
			v, err := p.pdbConvert()
			if err != nil {
				return false, err
			}
			p.pdbOption(opts, key, v, optStart)
		case key == "STORAGE":
			v, err := p.pdbStorage()
			if err != nil {
				return false, err
			}
			p.pdbOption(opts, "STORAGE", v, optStart)
		case key == "PATH_PREFIX", key == "CREATE_FILE_DEST",
			kind == pdbFromXML && key == "SOURCE_FILE_DIRECTORY":
			p.advance()
			if err := p.pdbExpectEquals(); err != nil {
				return false, err
			}
			v := "NONE"
			if !p.adbWord("NONE") {
				var err error
				if v, err = p.adbString(); err != nil {
					return false, err
				}
			} else {
				p.advance()
			}
			p.pdbOption(opts, key, v, optStart)
		case key == "TEMPFILE":
			p.advance()
			if err := p.adbExpect("REUSE"); err != nil {
				return false, err
			}
			p.pdbOption(opts, "TEMPFILE_REUSE", "", optStart)
		case key == "USER_TABLESPACES", key == "STANDBYS":
			p.advance()
			v, err := p.pdbAllExceptNone()
			if err != nil {
				return false, err
			}
			p.pdbOption(opts, key, v, optStart)
		case key == "LOGGING":
			p.pdbOption(opts, "LOGGING", p.advance().Str, optStart)
		case key == "CONTAINER_MAP":
			v, err := p.pdbContainerMap()
			if err != nil {
				return false, err
			}
			p.pdbOption(opts, "CONTAINER_MAP", v, optStart)
		case kind != pdbFromXML && key == "COPY_DATA":
			v, err := p.pdbCopyData()
			if err != nil {
				return false, err
			}
			p.pdbOption(opts, "DATA", v, optStart)
		case key == "KEYSTORE":
			p.advance()
			if p.cur.Type != kwIDENTIFIED {
				return false, p.syntaxErrorAtCur()
			}
			p.advance()
			if p.cur.Type != kwBY {
				return false, p.syntaxErrorAtCur()
			}
			p.advance()
			if err := p.pdbPassword(); err != nil {
				return false, err
			}
			p.pdbOption(opts, "KEYSTORE", "", optStart)
		case kind == pdbClone && key == "USING":
			p.advance()
			v, err := p.pdbUsingSnapshot()
			if err != nil {
				return false, err
			}
			p.pdbOption(opts, "USING_SNAPSHOT", v, optStart)
		case kind == pdbClone && key == "RELOCATE":
			p.advance()
			if !dblink {
				return false, p.syntaxErrorAtCur()
			}
			var parts []string
			if p.adbWord("KEEP") {
				p.advance()
				if err := p.adbExpect("SOURCE"); err != nil {
					return false, err
				}
				parts = append(parts, "KEEP SOURCE")
			}
			if p.adbWord("AVAILABILITY") {
				p.advance()
				if !p.adbWord("NORMAL") && !p.adbWord("MAX") {
					return false, p.syntaxErrorAtCur()
				}
				parts = append(parts, "AVAILABILITY "+p.advance().Str)
			}
			p.pdbOption(opts, "RELOCATE", strings.Join(parts, " "), optStart)
		case kind == pdbClone && key == "REFRESH":
			p.advance()
			if err := p.adbExpect("MODE"); err != nil {
				return false, err
			}
			if !dblink {
				return false, p.syntaxErrorAtCur()
			}
			v, err := p.pdbRefreshMode()
			if err != nil {
				return false, err
			}
			p.pdbOption(opts, "REFRESH_MODE", v, optStart)
		case kind == pdbFromXML && (key == "COPY" || key == "MOVE" || key == "NOCOPY"):
			seen["COPY"], seen["MOVE"], seen["NOCOPY"] = true, true, true
			p.pdbOption(opts, p.advance().Str, "", optStart)
		case kind == pdbFromXML && p.cur.Type == kwAS:
			p.advance()
			if err := p.adbExpect("CLONE"); err != nil {
				return false, err
			}
			p.pdbOption(opts, "AS_CLONE", "", optStart)
		case kind == pdbFromXML && key == "DECRYPT":
			p.advance()
			if p.cur.Type != kwUSING {
				return false, p.syntaxErrorAtCur()
			}
			p.advance()
			if err := p.pdbPassword(); err != nil {
				return false, err
			}
			p.pdbOption(opts, "DECRYPT_USING", "", optStart)
		default:
			return false, p.syntaxErrorAtCur()
		}
	}
	return false, nil
}

// pdbPassword reads a password or secret: an identifier, quoted or not,
// reserved or not. Oracle does not take a string here (ORA-00988, ORA-65002).
func (p *Parser) pdbPassword() error {
	if p.cur.Type != tokIDENT && p.cur.Type != tokQIDENT && p.cur.Type < 2000 {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}

// pdbParenRoles reads ( role [, role ]... ), each role a word, quoted or
// not: Oracle takes the reserved CONNECT and RESOURCE here.
func (p *Parser) pdbParenRoles() ([]string, error) {
	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	var names []string
	for {
		if p.cur.Type != tokIDENT && p.cur.Type != tokQIDENT && p.cur.Type < 2000 {
			return nil, p.syntaxErrorAtCur()
		}
		names = append(names, p.advance().Str)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	return names, nil
}

// pdbParenStrings reads ( 'string' [, 'string' ]... ).
func (p *Parser) pdbParenStrings() ([]string, error) {
	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	items, err := p.adbStringList()
	if err != nil {
		return nil, err
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	return items, nil
}

// pdbConvert reads, the keyword consumed,
//
//	= { ( 'pattern' , 'replacement' [, 'pattern' , 'replacement' ]... ) | NONE }
func (p *Parser) pdbConvert() (string, error) {
	if err := p.pdbExpectEquals(); err != nil {
		return "", err
	}
	if p.adbWord("NONE") {
		return p.advance().Str, nil
	}
	if p.cur.Type != '(' {
		return "", p.syntaxErrorAtCur()
	}
	p.advance()
	var items []string
	for {
		from, err := p.adbString()
		if err != nil {
			return "", err
		}
		if p.cur.Type != ',' {
			return "", p.syntaxErrorAtCur()
		}
		p.advance()
		to, err := p.adbString()
		if err != nil {
			return "", err
		}
		items = append(items, from, to)
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if p.cur.Type != ')' {
		return "", p.syntaxErrorAtCur()
	}
	p.advance()
	return "(" + strings.Join(items, ", ") + ")", nil
}

// pdbAllExceptNone reads, USER_TABLESPACES or STANDBYS consumed,
//
//	= { ( 'name' [, 'name' ]... ) | ALL [ EXCEPT ( 'name' [, 'name' ]... ) ] | NONE }
//
// Oracle takes the names as strings (ORA-01780).
func (p *Parser) pdbAllExceptNone() (string, error) {
	if err := p.pdbExpectEquals(); err != nil {
		return "", err
	}
	switch {
	case p.adbWord("NONE"):
		return p.advance().Str, nil
	case p.cur.Type == kwALL:
		p.advance()
		if !p.adbWord("EXCEPT") && p.cur.Type != kwEXCEPT {
			return "ALL", nil
		}
		p.advance()
		names, err := p.pdbParenStrings()
		if err != nil {
			return "", err
		}
		return "ALL EXCEPT (" + strings.Join(names, ", ") + ")", nil
	}
	names, err := p.pdbParenStrings()
	if err != nil {
		return "", err
	}
	return "(" + strings.Join(names, ", ") + ")", nil
}

// pdbStorage reads
//
//	pdb_storage_clause ::= STORAGE { UNLIMITED | ( { MAXSIZE | MAX_AUDIT_SIZE | MAX_DIAG_SIZE }
//	  { size_clause | UNLIMITED } ... ) }
//
// Oracle takes the limits in any order, each at most once (ORA-12990),
// without commas (ORA-65044).
func (p *Parser) pdbStorage() (string, error) {
	p.advance()
	if p.adbWord("UNLIMITED") {
		return p.advance().Str, nil
	}
	if p.cur.Type != '(' {
		return "", p.syntaxErrorAtCur()
	}
	p.advance()
	seen := map[string]bool{}
	var parts []string
	for p.cur.Type != ')' || len(parts) == 0 {
		if !p.adbWord("MAXSIZE") && !p.adbWord("MAX_AUDIT_SIZE") && !p.adbWord("MAX_DIAG_SIZE") || seen[p.cur.Str] {
			return "", p.syntaxErrorAtCur()
		}
		seen[p.cur.Str] = true
		limit := p.advance().Str
		if p.adbWord("UNLIMITED") {
			parts = append(parts, limit+" "+p.advance().Str)
			continue
		}
		n, err := p.parseAlterDatabaseSize()
		if err != nil {
			return "", err
		}
		parts = append(parts, limit+" "+n)
	}
	p.advance()
	return "(" + strings.Join(parts, " ") + ")", nil
}

// pdbPartitionValues reads ( value [, value ]... ), each a literal or
// MAXVALUE (ORA-00907 for an expression).
func (p *Parser) pdbPartitionValues() ([]string, error) {
	if p.cur.Type != '(' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	var vals []string
	for {
		switch {
		case p.cur.Type == tokICONST, p.cur.Type == tokFCONST, p.cur.Type == tokSCONST,
			p.cur.Type == kwMAXVALUE, p.adbWord("MAXVALUE"):
			vals = append(vals, p.advance().Str)
		default:
			return nil, p.syntaxErrorAtCur()
		}
		if p.cur.Type != ',' {
			break
		}
		p.advance()
	}
	if p.cur.Type != ')' {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	return vals, nil
}

// pdbContainerMap reads
//
//	container_map_clause ::= CONTAINER_MAP UPDATE ( { add_partition_clause | split_partition_clause } )
//	add_partition_clause ::= ADD PARTITION [ partition_name ] VALUES { LESS THAN ( value ) | ( value ) }
//	split_partition_clause ::= SPLIT PARTITION partition_name AT ( value )
//	  INTO ( PARTITION partition_name , PARTITION partition_name )
//
// Oracle takes the name of an added partition left out, and a list
// partition's VALUES ( value ).
func (p *Parser) pdbContainerMap() (string, error) {
	p.advance()
	if err := p.adbExpect("UPDATE"); err != nil {
		return "", err
	}
	if p.cur.Type != '(' {
		return "", p.syntaxErrorAtCur()
	}
	p.advance()
	var v string
	switch {
	case p.cur.Type == kwADD:
		p.advance()
		if p.cur.Type != kwPARTITION {
			return "", p.syntaxErrorAtCur()
		}
		p.advance()
		v = "ADD PARTITION"
		if p.cur.Type != kwVALUES || p.isIdentLikeStrAt(p.peekNext(), "VALUES") {
			name, err := p.adbName()
			if err != nil {
				return "", err
			}
			v += " " + name
		}
		if p.cur.Type != kwVALUES {
			return "", p.syntaxErrorAtCur()
		}
		p.advance()
		if p.adbWord("LESS") {
			p.advance()
			if err := p.adbExpect("THAN"); err != nil {
				return "", err
			}
			v += " VALUES LESS THAN"
		} else {
			v += " VALUES"
		}
		vals, err := p.pdbPartitionValues()
		if err != nil {
			return "", err
		}
		v += " (" + strings.Join(vals, ", ") + ")"
	case p.adbWord("SPLIT"):
		p.advance()
		if p.cur.Type != kwPARTITION {
			return "", p.syntaxErrorAtCur()
		}
		p.advance()
		name, err := p.adbName()
		if err != nil {
			return "", err
		}
		if !p.adbWord("AT") {
			return "", p.syntaxErrorAtCur()
		}
		p.advance()
		vals, err := p.pdbPartitionValues()
		if err != nil {
			return "", err
		}
		if p.cur.Type != kwINTO {
			return "", p.syntaxErrorAtCur()
		}
		p.advance()
		if p.cur.Type != '(' {
			return "", p.syntaxErrorAtCur()
		}
		p.advance()
		var into []string
		for i := 0; i < 2; i++ {
			if i == 1 {
				if p.cur.Type != ',' {
					return "", p.syntaxErrorAtCur()
				}
				p.advance()
			}
			if p.cur.Type != kwPARTITION {
				return "", p.syntaxErrorAtCur()
			}
			p.advance()
			part, err := p.adbName()
			if err != nil {
				return "", err
			}
			into = append(into, part)
		}
		if p.cur.Type != ')' {
			return "", p.syntaxErrorAtCur()
		}
		p.advance()
		v = "SPLIT PARTITION " + name + " AT (" + strings.Join(vals, ", ") + ") INTO (" + strings.Join(into, ", ") + ")"
	default:
		return "", p.syntaxErrorAtCur()
	}
	if p.cur.Type != ')' {
		return "", p.syntaxErrorAtCur()
	}
	p.advance()
	return v, nil
}

// pdbCopyData reads SNAPSHOT COPY [ NO DATA ] | NO DATA.
func (p *Parser) pdbCopyData() (string, error) {
	var parts []string
	if p.adbWord("SNAPSHOT") {
		p.advance()
		if err := p.adbExpect("COPY"); err != nil {
			return "", err
		}
		parts = append(parts, "SNAPSHOT COPY")
		if !p.adbWord("NO") {
			return parts[0], nil
		}
	}
	p.advance()
	if err := p.adbExpect("DATA"); err != nil {
		return "", err
	}
	return strings.Join(append(parts, "NO DATA"), " "), nil
}

// pdbUsingSnapshot reads, USING consumed,
//
//	SNAPSHOT { snapshot_name | AT SCN integer | AT TIME 'timestamp' }
func (p *Parser) pdbUsingSnapshot() (string, error) {
	if err := p.adbExpect("SNAPSHOT"); err != nil {
		return "", err
	}
	if !p.adbWord("AT") {
		if !isPDBSnapshotName(p.cur) {
			return "", p.syntaxErrorAtCur()
		}
		return p.advance().Str, nil
	}
	p.advance()
	switch {
	case p.adbWord("SCN"):
		p.advance()
		n, err := p.adbInteger()
		return "AT SCN " + n, err
	case p.adbWord("TIME"):
		p.advance()
		s, err := p.adbString()
		return "AT TIME " + s, err
	}
	return "", p.syntaxErrorAtCur()
}

// isPDBSnapshotName reports whether tok can name a snapshot in USING
// SNAPSHOT: an identifier that is not one of Oracle's keywords, as SCN, TIME
// and A (of IS A SET) are (ORA-00922).
func isPDBSnapshotName(tok Token) bool {
	return tok.Type == tokQIDENT || tok.Type == tokIDENT && tok.Str != "A" && tok.Str != "TIME"
}

// pdbRefreshMode reads, REFRESH MODE consumed,
// { MANUAL | EVERY integer { MINUTES | HOURS } | NONE }.
func (p *Parser) pdbRefreshMode() (string, error) {
	switch {
	case p.adbWord("MANUAL"), p.adbWord("NONE"):
		return p.advance().Str, nil
	case p.adbWord("EVERY"):
		p.advance()
		n, err := p.adbInteger()
		if err != nil {
			return "", err
		}
		if !p.adbWord("MINUTES") && !p.adbWord("HOURS") {
			return "", p.syntaxErrorAtCur()
		}
		return "EVERY " + n + " " + p.advance().Str, nil
	}
	return "", p.syntaxErrorAtCur()
}

// isAlterPDBClauseStart reports whether the current token opens an ALTER
// PLUGGABLE DATABASE clause rather than naming the PDB.
func (p *Parser) isAlterPDBClauseStart() bool {
	if p.cur.Type == tokQIDENT {
		return false
	}
	switch p.cur.Str {
	case "UNPLUG", "DEFAULT", "SET", "RENAME", "DATAFILE", "TEMPFILE", "CREATE",
		"ADD", "DROP", "STORAGE", "LOGGING", "NOLOGGING", "ENABLE", "DISABLE",
		"FORCE", "REFRESH", "CONTAINERS", "PRIORITY", "RECOVER", "BACKUP",
		"BEGIN", "END", "OPEN", "CLOSE", "SAVE", "DISCARD", "APPLICATION",
		"SNAPSHOT", "MATERIALIZE", "PREPARE", "ALL":
		return p.isIdentLike()
	}
	return false
}

// isAlterPDBStateStart reports whether the current token opens
// pdb_change_state: OPEN, CLOSE, SAVE or DISCARD.
func (p *Parser) isAlterPDBStateStart() bool {
	return p.adbWord("OPEN") || p.adbWord("CLOSE") || p.adbWord("SAVE") || p.adbWord("DISCARD")
}

// parseAlterPluggableDatabaseStmt parses ALTER PLUGGABLE DATABASE, ALTER
// PLUGGABLE DATABASE consumed.
//
//	ALTER PLUGGABLE DATABASE
//	  { [ pdb_name ] { pdb_unplug_clause | pdb_settings_clauses | pdb_datafile_clause
//	                 | pdb_recovery_clauses | pdb_change_state | snapshot_clauses }
//	  | { pdb_name [, pdb_name ]... | ALL [ EXCEPT pdb_name [, pdb_name ]... ] } pdb_change_state
//	  | application_clauses }
//
// Oracle takes ALL EXCEPT's names without the page's parentheses
// (ORA-65000), and UNPLUG only after a PDB name (ORA-65000).
func (p *Parser) parseAlterPluggableDatabaseStmt(start int) (nodes.StmtNode, error) {
	stmt := &nodes.AdminDDLStmt{
		Action:     "ALTER",
		ObjectType: nodes.OBJECT_PLUGGABLE_DATABASE,
		Loc:        nodes.Loc{Start: start},
	}
	opts := &nodes.List{}
	if p.adbEndKept() {
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}
	stateOnly := false
	switch {
	case p.cur.Type == kwALL:
		targetStart := p.pos()
		p.advance()
		target := p.pdbOption(opts, "TARGET", "ALL", targetStart)
		if p.cur.Type == kwEXCEPT || p.adbWord("EXCEPT") {
			p.advance()
			names, err := p.pdbNameItems()
			if err != nil {
				return nil, err
			}
			target.Value = "ALL_EXCEPT"
			target.Items = names
			target.Loc.End = p.prev.End
		}
		stateOnly = true
	case p.adbWord("APPLICATION"):
	case !p.isAlterPDBClauseStart():
		if p.cur.Type != tokQIDENT && isOracleSQLReservedKeyword(p.cur) {
			// Oracle reads a reserved word here as the PDB name and refuses
			// it (ORA-65000) when a clause follows; otherwise it fails at
			// the token after it.
			tok := p.advance()
			if p.adbEnd() || p.isAlterPDBClauseStart() {
				return nil, p.syntaxErrorAtTok(tok)
			}
			return nil, p.syntaxErrorAtCur()
		}
		name, err := p.pdbName()
		if err != nil {
			return nil, err
		}
		stmt.Name = name
		if p.cur.Type == ',' {
			p.advance()
			names, err := p.pdbNameItems()
			if err != nil {
				return nil, err
			}
			first := &nodes.DDLOption{Key: "PDB", Value: name.Name, Loc: name.Loc}
			names.Items = append([]nodes.Node{first}, names.Items...)
			opts.Items = append(opts.Items, &nodes.DDLOption{
				Key: "TARGET", Value: "LIST", Items: names,
				Loc: nodes.Loc{Start: name.Loc.Start, End: p.prev.End},
			})
			stateOnly = true
		}
		if p.adbWord("APPLICATION") {
			return nil, p.syntaxErrorAtCur()
		}
	}
	{
		if stateOnly && !p.isAlterPDBStateStart() {
			return nil, p.syntaxErrorAtCur()
		}
		opt, err := p.parseAlterPDBClause(stmt.Name != nil)
		if err != nil {
			return nil, err
		}
		if opt != nil {
			opts.Items = append(opts.Items, opt)
		}
	}
	if len(opts.Items) > 0 {
		stmt.Options = opts
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// pdbNameItems reads pdb_name [, pdb_name ]... as options, each spanning its
// name.
func (p *Parser) pdbNameItems() (*nodes.List, error) {
	items := &nodes.List{}
	for {
		start := p.pos()
		name, err := p.adbName()
		if err != nil {
			return nil, err
		}
		items.Items = append(items.Items, p.adbItem("PDB", name, start))
		if p.cur.Type != ',' {
			return items, nil
		}
		p.advance()
	}
}

// parseAlterPDBClause reads the one clause of ALTER PLUGGABLE DATABASE.
// named says whether a PDB name came before it. It returns nil where the
// statement ended where Oracle reports the missing clause with a code of
// the statement.
func (p *Parser) parseAlterPDBClause(named bool) (*nodes.DDLOption, error) {
	optStart := p.pos()
	opt := &nodes.DDLOption{}
	if p.adbEnd() || p.cur.Type == tokQIDENT || !p.isIdentLike() {
		return nil, p.syntaxErrorAtCur()
	}
	word := p.cur.Str
	var err error
	switch word {
	case "UNPLUG":
		if !named {
			return nil, p.syntaxErrorAtCur()
		}
		err = p.parseAlterPDBUnplug(opt)
	case "DEFAULT":
		err = p.parseAlterPDBDefault(opt)
	case "SET":
		err = p.parseAlterPDBSet(opt)
	case "RENAME":
		err = p.parseAlterDatabaseRename(opt)
	case "DATAFILE":
		if p.peekNext().Type == kwALL {
			p.advance()
			p.advance()
			opt.Key, opt.Value = "DATAFILE", "ALL"
			if p.adbEndKept() {
				break
			}
			if p.cur.Type != kwONLINE && p.cur.Type != kwOFFLINE {
				return nil, p.syntaxErrorAtCur()
			}
			actionStart := p.pos()
			opt.Items = &nodes.List{Items: []nodes.Node{p.adbItem(p.advance().Str, "", actionStart)}}
			break
		}
		err = p.parseAlterDatabaseFile(opt, word)
	case "TEMPFILE":
		err = p.parseAlterDatabaseFile(opt, word)
	case "CREATE":
		p.advance()
		opt.Key = "CREATE_DATAFILE"
		if !p.adbWord("DATAFILE") {
			return nil, p.syntaxErrorAtCur()
		}
		err = p.parseAlterDatabaseCreateDatafile(opt)
	case "ADD", "DROP":
		err = p.parseAlterPDBAddDrop(opt, word)
	case "STORAGE":
		opt.Key = "STORAGE"
		opt.Value, err = p.pdbStorage()
	case "LOGGING", "NOLOGGING", "FORCE":
		// Oracle takes FORCE without ENABLE or DISABLE only when no PDB is
		// named (ORA-00922).
		if named && word == "FORCE" {
			return nil, p.syntaxErrorAtCur()
		}
		err = p.parseAlterPDBLogging(opt, "")
	case "ENABLE", "DISABLE":
		p.advance()
		switch {
		case p.adbWord("BACKUP"), p.adbWord("RECOVERY"):
			opt.Key = word + "_" + p.advance().Str
		case p.adbWord("LOST"):
			p.advance()
			opt.Key = word + "_LOST_WRITE_PROTECTION"
			if err = p.adbExpect("WRITE"); err == nil {
				err = p.adbExpect("PROTECTION")
			}
		default:
			err = p.parseAlterPDBLogging(opt, word)
		}
	case "REFRESH":
		p.advance()
		opt.Key = "REFRESH"
		if p.adbWord("MODE") {
			p.advance()
			opt.Key = "REFRESH_MODE"
			opt.Value, err = p.pdbRefreshMode()
		}
	case "CONTAINERS":
		err = p.parseAlterPDBContainers(opt)
	case "PRIORITY":
		p.advance()
		opt.Key = "PRIORITY"
		if p.adbWord("NONE") {
			opt.Value = p.advance().Str
		} else {
			opt.Value, err = p.adbInteger()
		}
	case "RECOVER":
		err = p.parseAlterDatabaseRecover(opt)
	case "BACKUP":
		p.advance()
		if !p.adbWord("BEGIN") && !p.adbWord("END") {
			return nil, p.syntaxErrorAtCur()
		}
		opt.Key = "BACKUP_" + p.advance().Str
	case "BEGIN", "END":
		p.advance()
		opt.Key = "BACKUP_" + word
		if !p.adbEndKept() {
			err = p.adbExpect("BACKUP")
		}
	case "OPEN":
		err = p.parseAlterPDBOpen(opt)
	case "CLOSE":
		err = p.parseAlterPDBClose(opt)
	case "SAVE", "DISCARD":
		p.advance()
		opt.Key = word + "_STATE"
		if err = p.adbExpect("STATE"); err == nil && p.adbWord("INSTANCES") {
			opt.Items, err = p.pdbInstancesOrServices("INSTANCES")
		}
	case "APPLICATION":
		if named {
			return nil, p.syntaxErrorAtCur()
		}
		err = p.parseAlterPDBApplication(opt)
	case "SNAPSHOT":
		p.advance()
		opt.Key = "SNAPSHOT"
		if !p.adbEndKept() {
			opt.Value, err = p.pdbRefreshMode()
		}
	case "MATERIALIZE":
		p.advance()
		opt.Key = "MATERIALIZE"
	case "PREPARE":
		err = p.parseAlterPDBPrepareMirror(opt)
	default:
		return nil, p.syntaxErrorAtCur()
	}
	if err != nil {
		return nil, err
	}
	if !p.adbEnd() {
		return nil, p.syntaxErrorAtCur()
	}
	opt.Loc = nodes.Loc{Start: optStart, End: p.prev.End}
	return opt, nil
}

// parseAlterPDBUnplug reads UNPLUG INTO 'filename' [ ENCRYPT USING transport_secret ].
func (p *Parser) parseAlterPDBUnplug(opt *nodes.DDLOption) error {
	p.advance()
	opt.Key = "UNPLUG_INTO"
	if p.cur.Type != kwINTO {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	f, err := p.adbString()
	if err != nil {
		return err
	}
	opt.Value = f
	if p.adbWord("ENCRYPT") {
		p.advance()
		if p.cur.Type != kwUSING {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		return p.pdbPassword()
	}
	return nil
}

// parseAlterPDBDefault reads
//
//	DEFAULT EDITION = edition_name
//	DEFAULT TABLESPACE tablespace_name
//	DEFAULT [ LOCAL ] TEMPORARY TABLESPACE tablespace_name
//
// Oracle takes EDITION with '=', as in ALTER DATABASE, not as the page shows
// (ORA-00927).
func (p *Parser) parseAlterPDBDefault(opt *nodes.DDLOption) error {
	p.advance()
	switch {
	case p.adbWord("EDITION"):
		p.advance()
		opt.Key = "DEFAULT_EDITION"
		if err := p.pdbExpectEquals(); err != nil {
			return err
		}
		name, err := p.adbName()
		opt.Value = name
		return err
	case p.adbWord("TABLESPACE"):
		p.advance()
		opt.Key = "DEFAULT_TABLESPACE"
	case p.adbWord("LOCAL"), p.adbWord("TEMPORARY"):
		opt.Key = "DEFAULT_TEMPORARY_TABLESPACE"
		if p.adbWord("LOCAL") {
			p.advance()
			opt.Key = "DEFAULT_LOCAL_TEMPORARY_TABLESPACE"
		}
		if err := p.adbExpect("TEMPORARY"); err != nil {
			return err
		}
		if err := p.adbExpect("TABLESPACE"); err != nil {
			return err
		}
	default:
		return p.syntaxErrorAtCur()
	}
	if p.adbEndKept() {
		return nil
	}
	name, err := p.adbName()
	opt.Value = name
	return err
}

// parseAlterPDBSet reads the clauses that begin with SET:
//
//	SET DEFAULT { BIGFILE | SMALLFILE } TABLESPACE
//	SET TIME_ZONE = 'time_zone'
//	SET STANDBY NOLOGGING FOR { LOAD PERFORMANCE | DATA AVAILABILITY }
//	SET MAX_PDB_SNAPSHOTS = integer
func (p *Parser) parseAlterPDBSet(opt *nodes.DDLOption) error {
	p.advance()
	switch {
	case p.cur.Type == kwDEFAULT:
		p.advance()
		opt.Key = "SET_DEFAULT_TABLESPACE_TYPE"
		if !p.adbWord("BIGFILE") && !p.adbWord("SMALLFILE") {
			return p.syntaxErrorAtCur()
		}
		opt.Value = p.advance().Str
		if p.adbEndKept() {
			return nil
		}
		return p.adbExpect("TABLESPACE")
	case p.adbWord("TIME_ZONE"):
		p.advance()
		opt.Key = "SET_TIME_ZONE"
		if err := p.pdbExpectEquals(); err != nil {
			return err
		}
		tz, err := p.adbString()
		opt.Value = tz
		return err
	case p.adbWord("STANDBY"):
		p.advance()
		opt.Key = "SET_STANDBY_NOLOGGING"
		if err := p.adbExpect("NOLOGGING"); err != nil {
			return err
		}
		if p.cur.Type != kwFOR {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		switch {
		case p.adbWord("LOAD"):
			p.advance()
			opt.Value = "LOAD PERFORMANCE"
			return p.adbExpect("PERFORMANCE")
		case p.adbWord("DATA"):
			p.advance()
			opt.Value = "DATA AVAILABILITY"
			return p.adbExpect("AVAILABILITY")
		}
		return p.syntaxErrorAtCur()
	case p.adbWord("MAX_PDB_SNAPSHOTS"):
		p.advance()
		opt.Key = "SET_MAX_PDB_SNAPSHOTS"
		if err := p.pdbExpectEquals(); err != nil {
			return err
		}
		n, err := p.adbInteger()
		opt.Value = n
		return err
	}
	return p.syntaxErrorAtCur()
}

// parseAlterPDBAddDrop reads
//
//	{ ADD | DROP } SUPPLEMENTAL LOG DATA [ ... ]
//	DROP SNAPSHOT snapshot_name
//	DROP MIRROR COPY mirror_name
func (p *Parser) parseAlterPDBAddDrop(opt *nodes.DDLOption, verb string) error {
	p.advance()
	switch {
	case p.adbWord("SUPPLEMENTAL"):
		p.advance()
		opt.Key = verb + "_SUPPLEMENTAL_LOG"
		if err := p.adbExpect("LOG"); err != nil {
			return err
		}
		if err := p.adbExpect("DATA"); err != nil {
			return err
		}
		v, err := p.parseAlterDatabaseSupplementalData()
		opt.Value = v
		return err
	case verb == "DROP" && p.adbWord("SNAPSHOT"):
		p.advance()
		opt.Key = "DROP_SNAPSHOT"
		name, err := p.adbName()
		opt.Value = name
		return err
	case verb == "DROP" && p.adbWord("MIRROR"):
		p.advance()
		opt.Key = "DROP_MIRROR_COPY"
		if err := p.adbExpect("COPY"); err != nil {
			return err
		}
		name, err := p.adbName()
		opt.Value = name
		return err
	}
	return p.syntaxErrorAtCur()
}

// parseAlterPDBLogging reads, ENABLE or DISABLE (verb) consumed where given,
//
//	[ FORCE ] { LOGGING | NOLOGGING }
//
// Oracle takes FORCE and the verb each without the other, not only the
// page's { ENABLE | DISABLE } FORCE { LOGGING | NOLOGGING } and logging_clause.
func (p *Parser) parseAlterPDBLogging(opt *nodes.DDLOption, verb string) error {
	var parts []string
	if verb != "" {
		parts = append(parts, verb)
	}
	if p.cur.Type == kwFORCE || p.adbWord("FORCE") {
		parts = append(parts, p.advance().Str)
	}
	if !p.adbWord("LOGGING") && !p.adbWord("NOLOGGING") {
		return p.syntaxErrorAtCur()
	}
	parts = append(parts, p.advance().Str)
	opt.Key = strings.Join(parts, "_")
	return nil
}

// parseAlterPDBContainers reads
//
//	CONTAINERS { DEFAULT TARGET = { ( container_name ) | NONE }
//	           | HOST = 'hostname' | PORT = port_number }
//
// Oracle takes the default target in parentheses (ORA-65324).
func (p *Parser) parseAlterPDBContainers(opt *nodes.DDLOption) error {
	p.advance()
	switch {
	case p.cur.Type == kwDEFAULT:
		p.advance()
		opt.Key = "CONTAINERS_DEFAULT_TARGET"
		if err := p.adbExpect("TARGET"); err != nil {
			return err
		}
		if err := p.pdbExpectEquals(); err != nil {
			return err
		}
		if p.adbWord("NONE") {
			opt.Value = p.advance().Str
			return nil
		}
		if p.cur.Type != '(' {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		name, err := p.adbName()
		if err != nil {
			return err
		}
		opt.Value = name
		if p.cur.Type != ')' {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		return nil
	case p.adbWord("HOST"):
		p.advance()
		opt.Key = "CONTAINERS_HOST"
		if err := p.pdbExpectEquals(); err != nil {
			return err
		}
		host, err := p.adbString()
		opt.Value = host
		return err
	case p.adbWord("PORT"):
		p.advance()
		opt.Key = "CONTAINERS_PORT"
		if err := p.pdbExpectEquals(); err != nil {
			return err
		}
		port, err := p.adbInteger()
		opt.Value = port
		return err
	}
	return p.syntaxErrorAtCur()
}

// parseAlterPDBOpen reads
//
//	OPEN [ READ WRITE [ UPGRADE ] | READ ONLY | HYBRID [ READ ONLY ] | UPGRADE | RESETLOGS ]
//	  [ RESTRICTED ] [ FORCE ] [ instances_clause ] [ services_clause ]
//
// Oracle takes RESETLOGS only as the mode, and READ WRITE UPGRADE.
func (p *Parser) parseAlterPDBOpen(opt *nodes.DDLOption) error {
	p.advance()
	opt.Key = "OPEN"
	var parts []string
	switch {
	case p.cur.Type == kwREAD:
		p.advance()
		switch {
		case p.cur.Type == kwWRITE:
			p.advance()
			parts = append(parts, "READ WRITE")
			if p.adbWord("UPGRADE") {
				parts = append(parts, p.advance().Str)
			}
		case p.cur.Type == kwONLY:
			p.advance()
			parts = append(parts, "READ ONLY")
		default:
			return p.syntaxErrorAtCur()
		}
	case p.adbWord("HYBRID"):
		p.advance()
		parts = append(parts, "HYBRID")
		if p.cur.Type == kwREAD {
			p.advance()
			if p.cur.Type != kwONLY {
				return p.syntaxErrorAtCur()
			}
			p.advance()
			parts = append(parts, "READ ONLY")
		}
	case p.adbWord("UPGRADE"), p.adbWord("RESETLOGS"):
		parts = append(parts, p.advance().Str)
	}
	if p.adbWord("RESTRICTED") {
		parts = append(parts, p.advance().Str)
	}
	if p.cur.Type == kwFORCE || p.adbWord("FORCE") {
		parts = append(parts, p.advance().Str)
	}
	opt.Value = strings.Join(parts, " ")
	items := &nodes.List{}
	if p.adbWord("INSTANCES") {
		inst, err := p.pdbInstancesOrServices("INSTANCES")
		if err != nil {
			return err
		}
		items.Items = append(items.Items, inst.Items...)
	}
	if p.adbWord("SERVICES") {
		svc, err := p.pdbInstancesOrServices("SERVICES")
		if err != nil {
			return err
		}
		items.Items = append(items.Items, svc.Items...)
	}
	if len(items.Items) > 0 {
		opt.Items = items
	}
	return nil
}

// parseAlterPDBClose reads
//
//	CLOSE [ IMMEDIATE | ABORT ] [ instances_clause ]
//	CLOSE [ IMMEDIATE ] { RELOCATE [ TO 'instance_name' ] | NORELOCATE }
//
// Oracle does not take the relocate_clause after ABORT or an
// instances_clause (ORA-03049).
func (p *Parser) parseAlterPDBClose(opt *nodes.DDLOption) error {
	p.advance()
	opt.Key = "CLOSE"
	abort := false
	if p.cur.Type == kwIMMEDIATE || p.adbWord("ABORT") {
		abort = p.adbWord("ABORT")
		opt.Value = p.advance().Str
	}
	switch {
	case p.adbWord("INSTANCES"):
		items, err := p.pdbInstancesOrServices("INSTANCES")
		opt.Items = items
		return err
	case !abort && p.adbWord("RELOCATE"):
		start := p.pos()
		p.advance()
		value := ""
		if p.cur.Type == kwTO {
			p.advance()
			if !p.adbEndKept() {
				inst, err := p.adbString()
				if err != nil {
					return err
				}
				value = inst
			}
		}
		opt.Items = &nodes.List{Items: []nodes.Node{p.adbItem("RELOCATE", value, start)}}
	case !abort && p.adbWord("NORELOCATE"):
		start := p.pos()
		p.advance()
		opt.Items = &nodes.List{Items: []nodes.Node{p.adbItem("NORELOCATE", "", start)}}
	}
	return nil
}

// pdbInstancesOrServices reads, at INSTANCES or SERVICES (key),
//
//	key = { ( 'name' [, 'name' ]... ) | ALL [ EXCEPT ( 'name' [, 'name' ]... ) ] }
//
// SERVICES also takes NONE. Oracle takes the names as strings (ORA-01780).
func (p *Parser) pdbInstancesOrServices(key string) (*nodes.List, error) {
	start := p.pos()
	p.advance()
	if err := p.pdbExpectEquals(); err != nil {
		return nil, err
	}
	var v string
	switch {
	case key == "SERVICES" && p.adbWord("NONE"):
		v = p.advance().Str
	case p.cur.Type == kwALL:
		p.advance()
		v = "ALL"
		if p.cur.Type == kwEXCEPT || p.adbWord("EXCEPT") {
			p.advance()
			names, err := p.pdbParenStrings()
			if err != nil {
				return nil, err
			}
			v = "ALL EXCEPT (" + strings.Join(names, ", ") + ")"
		}
	default:
		names, err := p.pdbParenStrings()
		if err != nil {
			return nil, err
		}
		v = "(" + strings.Join(names, ", ") + ")"
	}
	return &nodes.List{Items: []nodes.Node{p.adbItem(key, v, start)}}, nil
}

// parseAlterPDBPrepareMirror reads
//
//	PREPARE MIRROR COPY mirror_name
//	  [ WITH { UNPROTECTED | MIRROR | HIGH } REDUNDANCY ] [ FOR DATABASE database_name ]
//
// Oracle takes the redundancies of ALTER DATABASE's prepare_clause, not the
// page's EXTERNAL and NORMAL (ORA-00922).
func (p *Parser) parseAlterPDBPrepareMirror(opt *nodes.DDLOption) error {
	p.advance()
	opt.Key = "PREPARE_MIRROR_COPY"
	if err := p.adbExpect("MIRROR"); err != nil {
		return err
	}
	if err := p.adbExpect("COPY"); err != nil {
		return err
	}
	name, err := p.adbName()
	if err != nil {
		return err
	}
	opt.Value = name
	if p.cur.Type == kwWITH {
		p.advance()
		if !p.adbWord("UNPROTECTED") && !p.adbWord("MIRROR") && !p.adbWord("HIGH") {
			return p.syntaxErrorAtCur()
		}
		opt.Value += " WITH " + p.advance().Str + " REDUNDANCY"
		if err := p.adbExpect("REDUNDANCY"); err != nil {
			return err
		}
	}
	if p.cur.Type == kwFOR {
		p.advance()
		if p.cur.Type != kwDATABASE {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		name, err := p.adbName()
		if err != nil {
			return err
		}
		opt.Value += " FOR DATABASE " + name
	}
	return nil
}

// parseAlterPDBApplication reads
//
//	APPLICATION { app_name app_clause
//	            | app_name [, app_name ]... SYNC
//	            | ALL [ EXCEPT app_name [, app_name ]... ] SYNC }
//	app_clause ::=
//	  { BEGIN { INSTALL 'app_version' [ COMMENT 'comment' ]
//	          | PATCH number [ MINIMUM VERSION 'app_version' | COMMENT 'comment' ]...
//	          | UPGRADE 'start_app_version' TO 'end_app_version' [ COMMENT 'comment' ]
//	          | UNINSTALL }
//	  | END { INSTALL [ 'app_version' ] | PATCH [ number ] | UPGRADE [ TO 'end_app_version' ] | UNINSTALL }
//	  | SET PATCH number | SET VERSION 'app_version'
//	  | SET COMPATIBILITY VERSION { 'app_version' | CURRENT }
//	  | SYNC [ TO { 'app_version' | PATCH number } ] }
//
// Oracle takes INSTALL and UPGRADE without the page's VERSION and FROM
// (ORA-00922, ORA-00946), and ALL EXCEPT's names without parentheses.
func (p *Parser) parseAlterPDBApplication(opt *nodes.DDLOption) error {
	p.advance()
	opt.Key = "APPLICATION"
	if p.cur.Type == kwALL {
		p.advance()
		opt.Value = "ALL"
		if p.cur.Type == kwEXCEPT || p.adbWord("EXCEPT") {
			p.advance()
			names, err := p.pdbNameItems()
			if err != nil {
				return err
			}
			opt.Value = "ALL EXCEPT"
			opt.Items = names
		}
		return p.adbExpect("SYNC")
	}
	names, err := p.pdbNameItems()
	if err != nil {
		return err
	}
	opt.Items = names
	if len(names.Items) > 1 {
		return p.adbExpect("SYNC")
	}
	opt.Value = names.Items[0].(*nodes.DDLOption).Value
	clauseStart := p.pos()
	var parts []string
	switch {
	case p.cur.Type == kwBEGIN || p.adbWord("BEGIN"):
		p.advance()
		parts = append(parts, "BEGIN")
		switch {
		case p.adbWord("INSTALL"):
			p.advance()
			v, err := p.adbString()
			if err != nil {
				return err
			}
			parts = append(parts, "INSTALL "+v)
			if err := p.pdbAppComment(&parts); err != nil {
				return err
			}
		case p.adbWord("PATCH"):
			p.advance()
			n, err := p.adbInteger()
			if err != nil {
				return err
			}
			parts = append(parts, "PATCH "+n)
			for {
				if p.adbWord("MINIMUM") {
					p.advance()
					if err := p.adbExpect("VERSION"); err != nil {
						return err
					}
					v, err := p.adbString()
					if err != nil {
						return err
					}
					parts = append(parts, "MINIMUM VERSION "+v)
					continue
				}
				if !p.adbWord("COMMENT") {
					break
				}
				if err := p.pdbAppComment(&parts); err != nil {
					return err
				}
			}
		case p.adbWord("UPGRADE"):
			p.advance()
			from, err := p.adbString()
			if err != nil {
				return err
			}
			if p.cur.Type != kwTO {
				return p.syntaxErrorAtCur()
			}
			p.advance()
			to, err := p.adbString()
			if err != nil {
				return err
			}
			parts = append(parts, "UPGRADE "+from+" TO "+to)
			if err := p.pdbAppComment(&parts); err != nil {
				return err
			}
		case p.adbWord("UNINSTALL"):
			parts = append(parts, p.advance().Str)
		default:
			return p.syntaxErrorAtCur()
		}
	case p.cur.Type == kwEND || p.adbWord("END"):
		p.advance()
		parts = append(parts, "END")
		switch {
		case p.adbWord("INSTALL"):
			parts = append(parts, p.advance().Str)
			if p.cur.Type == tokSCONST {
				parts = append(parts, p.advance().Str)
			}
		case p.adbWord("PATCH"):
			parts = append(parts, p.advance().Str)
			if p.cur.Type == tokICONST {
				parts = append(parts, p.advance().Str)
			}
		case p.adbWord("UPGRADE"):
			parts = append(parts, p.advance().Str)
			if p.cur.Type == kwTO {
				p.advance()
				to, err := p.adbString()
				if err != nil {
					return err
				}
				parts = append(parts, "TO "+to)
			}
		case p.adbWord("UNINSTALL"):
			parts = append(parts, p.advance().Str)
		default:
			return p.syntaxErrorAtCur()
		}
	case p.cur.Type == kwSET:
		p.advance()
		parts = append(parts, "SET")
		switch {
		case p.adbWord("PATCH"):
			p.advance()
			n, err := p.adbInteger()
			if err != nil {
				return err
			}
			parts = append(parts, "PATCH "+n)
		case p.adbWord("VERSION"):
			p.advance()
			v, err := p.adbString()
			if err != nil {
				return err
			}
			parts = append(parts, "VERSION "+v)
		case p.adbWord("COMPATIBILITY"):
			p.advance()
			if err := p.adbExpect("VERSION"); err != nil {
				return err
			}
			if p.adbWord("CURRENT") {
				parts = append(parts, "COMPATIBILITY VERSION "+p.advance().Str)
				break
			}
			v, err := p.adbString()
			if err != nil {
				return err
			}
			parts = append(parts, "COMPATIBILITY VERSION "+v)
		default:
			return p.syntaxErrorAtCur()
		}
	case p.adbWord("SYNC"):
		p.advance()
		parts = append(parts, "SYNC")
		if p.cur.Type == kwTO {
			p.advance()
			if p.adbWord("PATCH") {
				p.advance()
				n, err := p.adbInteger()
				if err != nil {
					return err
				}
				parts = append(parts, "TO PATCH "+n)
			} else {
				v, err := p.adbString()
				if err != nil {
					return err
				}
				parts = append(parts, "TO "+v)
			}
		}
	default:
		return p.syntaxErrorAtCur()
	}
	opt.Items.Items = append(opt.Items.Items, p.adbItem("CLAUSE", strings.Join(parts, " "), clauseStart))
	return nil
}

// pdbAppComment reads [ COMMENT 'comment' ] into parts.
func (p *Parser) pdbAppComment(parts *[]string) error {
	if !p.adbWord("COMMENT") {
		return nil
	}
	p.advance()
	c, err := p.adbString()
	if err != nil {
		return err
	}
	*parts = append(*parts, "COMMENT "+c)
	return nil
}
