package parser

import (
	"strings"

	nodes "github.com/bytebase/omni/oracle/ast"
)

// parseGrantStmt parses a GRANT statement.
//
// BNF: oracle/parser/bnf/GRANT.bnf
//
//	grant::=
//	    grant_system_privileges
//	  | grant_schema_privileges
//	  | grant_object_privileges
//	  | grant_roles_to_programs
//
//	grant_system_privileges::=
//	    GRANT { system_privilege [, system_privilege ]...
//	          | ALL PRIVILEGES
//	          | role [, role ]...
//	          }
//	        TO { grantee_clause [, grantee_clause ]... }
//	        [ IDENTIFIED BY password [, password ]... ]
//	        [ WITH { ADMIN | DELEGATE } OPTION ]
//	        [ CONTAINER = { CURRENT | ALL } ]
//
//	grant_schema_privileges::=
//	    GRANT { schema_privilege [, schema_privilege ]...
//	          | ALL [ PRIVILEGES ]
//	          }
//	        ON SCHEMA schema
//	        TO { grantee_clause [, grantee_clause ]... }
//	        [ WITH ADMIN OPTION ]
//	        [ CONTAINER = { CURRENT | ALL } ]
//
//	grant_object_privileges::=
//	    GRANT { object_privilege [ ( column [, column ]... ) ]
//	              [, object_privilege [ ( column [, column ]... ) ] ]...
//	          | ALL [ PRIVILEGES ]
//	          }
//	        on_object_clause
//	        TO { grantee_clause [, grantee_clause ]... }
//	        [ WITH { GRANT | HIERARCHY } OPTION ]
//	        [ CONTAINER = { CURRENT | ALL } ]
//
//	grant_roles_to_programs::=
//	    GRANT role [, role ]...
//	        TO program_unit [, program_unit ]...
//	        [ WITH { ADMIN | DELEGATE } OPTION ]
//	        [ CONTAINER = { CURRENT | ALL } ]
//
// How Oracle 23ai reads the privileges is in parseGrantPrivileges and the
// ON clause in parseGrantOnClause. Its grantees are names, one word each
// (ORA-03048 at a '.' or '@'). A grant to program units takes nothing after
// them (ORA-03048 at WITH). A grant of system privileges or roles takes the
// word of its WITH ... OPTION and the value of CONTAINER = as words, and
// refuses a wrong one after the parse (ORA-01939, ORA-65013); an object
// grant takes WITH GRANT OPTION or WITH HIERARCHY OPTION alone (ORA-00993).
func (p *Parser) parseGrantStmt() (nodes.StmtNode, error) {
	start := p.pos()
	p.advance() // consume GRANT

	stmt := &nodes.GrantStmt{Loc: nodes.Loc{Start: start}}
	privs, err := p.parseGrantPrivileges()
	if err != nil {
		return nil, err
	}
	stmt.AllPriv, stmt.Privileges = privs.all, privs.names()

	on := grantOnNone
	if p.cur.Type == kwON {
		if on, stmt.OnType, stmt.OnObject, err = p.parseGrantOnClause(privs); err != nil {
			return nil, err
		}
	} else if err := privs.checkWithoutOn(p); err != nil {
		return nil, err
	}

	if p.cur.Type != kwTO {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	programs := on == grantOnNone && p.isGrantProgramUnitStart()
	if stmt.Grantees, err = p.parseGrantees(programs, false); err != nil {
		return nil, err
	}
	if programs {
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	if on == grantOnNone && p.cur.Type == kwIDENTIFIED {
		// IDENTIFIED BY password [, password ]...
		p.advance()
		if p.cur.Type != kwBY {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		for {
			if !p.isIdentLike() {
				return nil, p.syntaxErrorAtCur()
			}
			p.advance()
			if p.cur.Type != ',' {
				break
			}
			p.advance()
		}
	}

	if p.cur.Type == kwWITH {
		p.advance()
		switch {
		case on == grantOnObject && p.cur.Type == kwGRANT:
			stmt.WithGrant = true
			p.advance()
		case on == grantOnObject && p.isKeywordStr("HIERARCHY"):
			p.advance()
		case on == grantOnObject:
			return nil, p.syntaxErrorAtCur()
		case p.cur.Type != kwOPTION && p.isIdentLike():
			stmt.WithAdmin = p.isKeywordStr("ADMIN")
			p.advance()
		}
		if p.cur.Type != kwOPTION {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
	}

	if err := p.parseGrantContainer(); err != nil {
		return nil, err
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// parseRevokeStmt parses a REVOKE statement.
//
// BNF: oracle/parser/bnf/REVOKE.bnf
//
//	REVOKE
//	  { revoke_system_privileges
//	  | revoke_schema_privileges
//	  | revoke_object_privileges
//	  | revoke_roles_from_programs
//	  }
//	  [ CONTAINER = { CURRENT | ALL } ] ;
//
//	revoke_system_privileges:
//	    REVOKE { system_privilege [, system_privilege ]...
//	           | role [, role ]...
//	           | ALL PRIVILEGES }
//	    FROM revokee_clause
//
//	revoke_schema_privileges:
//	    REVOKE { schema_privilege [, schema_privilege ]...
//	           | ALL [ PRIVILEGES ] }
//	    ON SCHEMA schema_name
//	    FROM revokee_clause
//
//	revoke_object_privileges:
//	    REVOKE { object_privilege [, object_privilege ]...
//	           | ALL [ PRIVILEGES ] }
//	    on_object_clause
//	    FROM revokee_clause
//	    [ CASCADE CONSTRAINTS ]
//	    [ FORCE ]
//
//	revokee_clause:
//	    { user | role | PUBLIC } [, { user | role | PUBLIC } ]...
//
//	revoke_roles_from_programs:
//	    REVOKE { role [, role ]... | ALL }
//	    FROM { PROCEDURE | FUNCTION | PACKAGE } [ schema. ] program_unit
//
// The privileges and the ON clause read as GRANT's do. Oracle 23ai takes
// CASCADE CONSTRAINTS or FORCE, not both (ORA-03049 at the second), after
// any revokees, refusing them where they do not apply after the parse, and
// a list of program units.
func (p *Parser) parseRevokeStmt() (nodes.StmtNode, error) {
	start := p.pos()
	p.advance() // consume REVOKE
	// The page states each revoke_roles_from_programs branch with the
	// statement's own REVOKE, so Oracle 23ai reads a repeated REVOKE before
	// the roles there (REVOKE REVOKE r FROM PROCEDURE p). It does not before
	// a revokee (REVOKE REVOKE r FROM u) and refuses that statement.
	repeated := false
	if p.cur.Type == kwREVOKE && p.revokeRepeatsBeforeRoles() {
		p.advance()
		repeated = true
	}

	stmt := &nodes.RevokeStmt{Loc: nodes.Loc{Start: start}}
	privs, err := p.parseGrantPrivileges()
	if err != nil {
		return nil, err
	}
	stmt.AllPriv, stmt.Privileges = privs.all, privs.names()

	on := grantOnNone
	if p.cur.Type == kwON {
		if on, stmt.OnType, stmt.OnObject, err = p.parseGrantOnClause(privs); err != nil {
			return nil, err
		}
	} else if err := privs.checkWithoutOn(p); err != nil {
		return nil, err
	}

	if p.cur.Type != kwFROM {
		return nil, p.syntaxErrorAtCur()
	}
	p.advance()
	programs := on == grantOnNone && p.isGrantProgramUnitStart()
	if repeated && !programs {
		return nil, p.syntaxErrorAtCur()
	}
	if stmt.Grantees, err = p.parseGrantees(programs, repeated); err != nil {
		return nil, err
	}
	if programs {
		stmt.Loc.End = p.prev.End
		return stmt, nil
	}

	switch {
	case p.cur.Type == kwCASCADE:
		p.advance()
		if !p.isKeywordStr("CONSTRAINTS") {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
	case p.isKeywordStr("FORCE"):
		p.advance()
	}

	if err := p.parseGrantContainer(); err != nil {
		return nil, err
	}
	stmt.Loc.End = p.prev.End
	return stmt, nil
}

// grantOnKind is what the ON clause of a GRANT or REVOKE names.
type grantOnKind int

const (
	grantOnNone grantOnKind = iota
	grantOnObject
	grantOnSchema
	grantOnUser
)

// grantPrivilege is one item of a GRANT or REVOKE privilege list: the words
// that name it, where it starts, and where its column list opens, if it has
// one.
type grantPrivilege struct {
	name       string
	start      Token
	columns    Token
	hasCols    bool
	columnList *nodes.List
}

// grantPrivileges is the privilege list of a GRANT or REVOKE.
type grantPrivileges struct {
	all        bool // ALL [ PRIVILEGES ]
	allStart   Token
	privileges bool // PRIVILEGES after ALL
	items      []grantPrivilege
}

func (g *grantPrivileges) names() *nodes.List {
	list := &nodes.List{}
	for _, item := range g.items {
		list.Items = append(list.Items, &nodes.String{Str: item.name})
	}
	return list
}

// grantObjectPrivileges are the privileges Oracle 23ai grants on an object
// (ORA-00990 at another one, as at a system privilege).
var grantObjectPrivileges = newAuditNames(
	"ALTER", "COLLECT DEBUG INFO", "CREATE", "DEBUG", "DELETE", "EXECUTE",
	"FLASHBACK", "INDEX", "INSERT", "KEEP SEQUENCE", "MERGE VIEW",
	"ON COMMIT REFRESH", "QUERY REWRITE", "READ", "REFERENCES", "SELECT",
	"SIGN", "UNDER", "UPDATE", "USE", "WRITE",
)

// grantUserPrivileges are the privileges Oracle 23ai grants ON USER.
var grantUserPrivileges = newAuditNames(
	"INHERIT PRIVILEGES", "INHERIT REMOTE PRIVILEGES",
)

// grantColumnPrivileges are the object privileges that take a column list
// (ORA-00969 at the '(' of another).
var grantColumnPrivileges = map[string]bool{"INSERT": true, "UPDATE": true, "REFERENCES": true}

// parseGrantPrivileges parses the privilege list of a GRANT or REVOKE.
// Oracle 23ai reads each item as the longest run of words that names a
// privilege, or else as one word, reserved or not, that names a role
// (GRANT CONNECT TO u, GRANT FROM TO u parse); an item it reads that way
// and the ',', ON, TO or FROM that does not follow it fail at the item
// (ORA-00990: GRANT a a TO u, GRANT u, a b TO u). ALL [ PRIVILEGES ]
// stands alone.
func (p *Parser) parseGrantPrivileges() (*grantPrivileges, error) {
	g := &grantPrivileges{}
	if p.cur.Type == kwALL {
		g.all = true
		g.allStart = p.cur
		p.advance()
		if p.cur.Type == kwPRIVILEGES {
			g.privileges = true
			p.advance()
		}
		return g, nil
	}
	for {
		item := grantPrivilege{start: p.cur}
		name, n := p.matchAuditName(auditSystemPrivileges, grantObjectPrivileges, grantUserPrivileges)
		switch {
		case n > 0 && !(p.cur.Type == kwON && n == 1):
			item.name = name
			for range n {
				p.advance()
			}
		case p.isIdentLike():
			var err error
			if item.name, err = p.parseIdentifier(); err != nil {
				return nil, err
			}
		default:
			return nil, p.syntaxErrorAtCur()
		}
		if p.cur.Type == '(' {
			if !grantColumnPrivileges[item.name] {
				return nil, p.syntaxErrorAtCur()
			}
			item.columns, item.hasCols = p.cur, true
			var err error
			if item.columnList, err = p.parseIdentifierListInParens(); err != nil {
				return nil, err
			}
		}
		switch p.cur.Type {
		case ',', kwON, kwTO, kwFROM:
		default:
			return nil, p.syntaxErrorAtTok(item.start)
		}
		g.items = append(g.items, item)
		if p.cur.Type != ',' {
			return g, nil
		}
		p.advance()
	}
}

// checkWithoutOn checks the privileges of a grant or revoke with no ON
// clause: system privileges, roles, or ALL [ PRIVILEGES ]. A run of words
// that names no system privilege fails at its start, and a column list at
// its '(' (ORA-00969).
func (g *grantPrivileges) checkWithoutOn(p *Parser) error {
	for _, item := range g.items {
		if item.hasCols {
			return p.syntaxErrorAtTok(item.columns)
		}
		// A quoted role name may hold blanks; words that run on name a
		// system privilege.
		if item.start.Type != tokQIDENT && strings.Contains(item.name, " ") && !auditSystemPrivileges.names[item.name] {
			return p.syntaxErrorAtTok(item.start)
		}
	}
	return nil
}

// parseGrantOnClause parses the ON clause of a GRANT or REVOKE, at ON,
// after checking the privileges against what it names.
//
//	on_object_clause::=
//	    ON [ schema. ] object
//	  | ON USER user [, user ]...
//	  | ON DIRECTORY directory_name
//	  | ON EDITION edition_name
//	  | ON MINING MODEL [ schema. ] mining_model_name
//	  | ON JAVA { SOURCE | RESOURCE } [ schema. ] object
//	  | ON SQL TRANSLATION PROFILE [ schema. ] profile
//
// Oracle 23ai grants object privileges on an object, INHERIT [ REMOTE ]
// PRIVILEGES on one user, and system privileges that name ANY, or ALL
// PRIVILEGES, on a schema; another privilege fails at its start (ORA-00990),
// and a system privilege on a schema at the ON. An object takes no object
// type before it (ORA-00903 at TABLE) and no database link (ORA-02021).
func (p *Parser) parseGrantOnClause(g *grantPrivileges) (grantOnKind, nodes.ObjectType, *nodes.ObjectName, error) {
	onTok := p.cur
	p.advance() // consume ON

	kind, objType := grantOnObject, nodes.OBJECT_TABLE
	next := p.peekNext()
	switch {
	case p.isKeywordStr("SCHEMA") && next.Type != '.':
		kind = grantOnSchema
	case p.cur.Type == kwUSER && next.Type != '.':
		kind, objType = grantOnUser, nodes.OBJECT_USER
	case p.isKeywordStr("DIRECTORY") && next.Type != '.':
		objType = nodes.OBJECT_DIRECTORY
	case p.isKeywordStr("EDITION") && next.Type != '.':
		objType = nodes.OBJECT_EDITION
	case p.isKeywordStr("JAVA") && next.Type != '.':
		objType = nodes.OBJECT_JAVA
	}

	// The privileges first: Oracle 23ai fails at a wrong one before it
	// reads the object.
	for _, item := range g.items {
		var ok bool
		switch kind {
		case grantOnSchema:
			if auditSystemPrivileges.names[item.name] && !strings.Contains(" "+item.name+" ", " ANY ") {
				return 0, 0, nil, p.syntaxErrorAtTok(onTok)
			}
			ok = auditSystemPrivileges.names[item.name]
		case grantOnUser:
			ok = grantUserPrivileges.names[item.name]
		default:
			ok = grantObjectPrivileges.names[item.name]
		}
		if !ok {
			return 0, 0, nil, p.syntaxErrorAtTok(item.start)
		}
		if item.hasCols && kind != grantOnObject {
			return 0, 0, nil, p.syntaxErrorAtTok(item.columns)
		}
	}
	if kind == grantOnSchema && g.all && !g.privileges {
		return 0, 0, nil, p.syntaxErrorAtTok(g.allStart)
	}

	var err error
	var obj *nodes.ObjectName
	switch {
	case kind == grantOnSchema, kind == grantOnUser:
		p.advance() // consume SCHEMA or USER
		obj, err = p.parseGrantObjectName(false)
	case objType == nodes.OBJECT_DIRECTORY, objType == nodes.OBJECT_EDITION:
		p.advance()
		obj, err = p.parseGrantObjectName(false)
	case objType == nodes.OBJECT_JAVA:
		p.advance()
		if !p.isKeywordStr("SOURCE") && !p.isKeywordStr("RESOURCE") {
			return 0, 0, nil, p.syntaxErrorAtCur()
		}
		p.advance()
		obj, err = p.parseGrantObjectName(true)
	case p.isKeywordStr("MINING") && next.Type == kwMODEL:
		p.advance()
		p.advance()
		obj, err = p.parseGrantObjectName(true)
	case p.isKeywordStr("SQL") && p.isGrantTranslationProfileAhead():
		p.advance() // SQL
		p.advance() // TRANSLATION
		p.advance() // PROFILE
		obj, err = p.parseGrantObjectName(true)
	default:
		obj, err = p.parseGrantObjectName(true)
	}
	if err != nil {
		return 0, 0, nil, err
	}
	return kind, objType, obj, nil
}

// isGrantTranslationProfileAhead reports whether SQL TRANSLATION PROFILE
// starts at the current token.
func (p *Parser) isGrantTranslationProfileAhead() bool {
	ahead := p.peekAhead(2)
	return ahead[0].Type != tokQIDENT && ahead[0].Str == "TRANSLATION" &&
		ahead[1].Type != tokQIDENT && ahead[1].Str == "PROFILE"
}

// parseGrantObjectName parses the name an ON clause grants on: one name, or
// with qualified a [ schema. ] name, never a reserved word unquoted
// (ORA-00903) nor a database link.
func (p *Parser) parseGrantObjectName(qualified bool) (*nodes.ObjectName, error) {
	start := p.pos()
	part := func() (string, error) {
		if !p.isIdentLike() {
			return "", p.syntaxErrorAtCur()
		}
		if err := p.syntaxErrorIfReservedIdentifier(); err != nil {
			return "", err
		}
		return p.parseIdentifier()
	}
	name, err := part()
	if err != nil {
		return nil, err
	}
	obj := &nodes.ObjectName{Name: name}
	if qualified && p.cur.Type == '.' {
		p.advance()
		obj.Schema = name
		if obj.Name, err = part(); err != nil {
			return nil, err
		}
	}
	// A database link is no object granted on (ORA-02021 at its name).
	if p.cur.Type == '@' {
		p.advance()
		return nil, p.syntaxErrorAtCur()
	}
	obj.Loc = nodes.Loc{Start: start, End: p.prev.End}
	return obj, nil
}

// revokeRepeatsBeforeRoles reports whether the REVOKE the parser stands on
// repeats the statement's keyword before a program-unit revoke's roles. A
// FROM there may itself be a role (REVOKE REVOKE FROM, r FROM PROCEDURE p);
// only a FROM followed by the target's kind is the delimiter.
func (p *Parser) revokeRepeatsBeforeRoles() bool {
	if p.peekNext().Type == kwFROM {
		ahead := p.peekAhead(2)
		if isGrantProgramUnitKind(ahead[1]) {
			return false
		}
	}
	return p.revokeProgramTargetAhead()
}

// revokeProgramTargetAhead reports whether the statement's FROM names a
// program unit.
func (p *Parser) revokeProgramTargetAhead() bool {
	lexer := *p.lexer
	next := lexer.NextToken
	var tok Token
	if p.hasNext {
		tok = p.nextBuf
	} else {
		tok = next()
	}
	for tok.Type != tokEOF && tok.Type != ';' {
		if tok.Type == kwFROM {
			tok = next()
			if isGrantProgramUnitKind(tok) {
				return true
			}
			continue
		}
		tok = next()
	}
	return false
}

// isGrantProgramUnitStart reports whether the grantees of a GRANT or
// REVOKE are program units: [ schema. ] { FUNCTION | PROCEDURE | PACKAGE }
// and a name. A kind word alone names a grantee: GRANT r TO FUNCTION
// parses on Oracle 23ai.
func (p *Parser) isGrantProgramUnitStart() bool {
	ahead := p.peekAhead(3)
	if isGrantProgramUnitKind(p.cur) {
		return isGrantWordToken(ahead[0])
	}
	return p.isIdentLike() && ahead[0].Type == '.' && isGrantProgramUnitKind(ahead[1])
}

// isGrantWordToken reports whether tok is a word: an identifier, quoted or
// not, or a keyword.
func isGrantWordToken(tok Token) bool {
	return tok.Type == tokIDENT || tok.Type == tokQIDENT || tok.Type >= 2000
}

func isGrantProgramUnitKind(tok Token) bool {
	switch tok.Type {
	case kwFUNCTION, kwPROCEDURE, kwPACKAGE:
		return true
	}
	return false
}

// parseGrantees parses the grantees of a GRANT or revokees of a REVOKE:
// names, or with programs program units, each [ schema. ] { FUNCTION |
// PROCEDURE | PACKAGE } name. A program unit is kept as its kind and name.
func (p *Parser) parseGrantees(programs bool, repeatedKind bool) (*nodes.List, error) {
	list := &nodes.List{}
	for {
		if programs {
			item, err := p.parseGrantProgramUnit(repeatedKind)
			if err != nil {
				return nil, err
			}
			list.Items = append(list.Items, &nodes.String{Str: item})
		} else {
			if !p.isIdentLike() {
				return nil, p.syntaxErrorAtCur()
			}
			name, err := p.parseIdentifier()
			if err != nil {
				return nil, err
			}
			list.Items = append(list.Items, &nodes.String{Str: name})
		}
		if p.cur.Type != ',' {
			return list, nil
		}
		p.advance()
	}
}

// parseGrantProgramUnit reads one program-unit grantee or revokee:
// [ schema. ] { FUNCTION | PROCEDURE | PACKAGE } name. With repeatedKind,
// the page's revoke_program_unit carries its own kind after the one the
// statement gives (FROM PROCEDURE FUNCTION p), optionally schema-qualified.
func (p *Parser) parseGrantProgramUnit(repeatedKind bool) (string, error) {
	prefix := ""
	if !isGrantProgramUnitKind(p.cur) {
		if !p.isIdentLike() {
			return "", p.syntaxErrorAtCur()
		}
		schema, err := p.parseIdentifier()
		if err != nil {
			return "", err
		}
		if p.cur.Type != '.' {
			return "", p.syntaxErrorAtCur()
		}
		p.advance() // consume '.'
		prefix = schema + "."
	}
	if !isGrantProgramUnitKind(p.cur) {
		return "", p.syntaxErrorAtCur()
	}
	kind := p.cur.Str
	p.advance()
	if repeatedKind && isGrantProgramUnitKind(p.cur) && isGrantWordToken(p.peekNext()) {
		unitKind := p.cur.Str
		p.advance()
		obj, err := p.parseGrantObjectName(false)
		if err != nil {
			return "", err
		}
		name := obj.Name
		if obj.Schema != "" {
			name = obj.Schema + "." + name
		}
		return prefix + kind + " " + unitKind + " " + name, nil
	}
	if repeatedKind && p.isIdentLike() && p.peekNext().Type == '.' {
		ahead := p.peekAhead(3)
		if isGrantProgramUnitKind(ahead[1]) && isGrantWordToken(ahead[2]) {
			schema, err := p.parseIdentifier()
			if err != nil {
				return "", err
			}
			p.advance() // consume '.'
			unitKind := p.cur.Str
			p.advance()
			obj, err := p.parseGrantObjectName(false)
			if err != nil {
				return "", err
			}
			return prefix + kind + " " + schema + "." + unitKind + " " + obj.Name, nil
		}
	}
	obj, err := p.parseGrantObjectName(true)
	if err != nil {
		return "", err
	}
	name := obj.Name
	if obj.Schema != "" {
		name = obj.Schema + "." + name
	}
	return prefix + kind + " " + name, nil
}

// parseGrantContainer parses [ CONTAINER = value ], whose value Oracle 23ai
// reads as a word and checks after the parse (ORA-65013 for ALLOCATE).
func (p *Parser) parseGrantContainer() error {
	if !p.isKeywordStr("CONTAINER") {
		return nil
	}
	p.advance()
	if p.cur.Type != '=' {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	if !p.isIdentLike() {
		return p.syntaxErrorAtCur()
	}
	p.advance()
	return nil
}
