package parser

import (
	"strings"

	"github.com/bytebase/omni/starrocks/ast"
)

// GRANT and REVOKE follow StarRocks.g4's grantRoleStatement,
// revokeRoleStatement (:1666, :1671), grantPrivilegeStatement and
// revokePrivilegeStatement (:1690, :1704).

// privTarget is what a privilege GRANT or REVOKE names after ON.
type privTarget struct {
	objectType     string
	objects        []*ast.ObjectName
	functionArgs   [][]*ast.TypeName
	inDatabase     string
	inAllDatabases bool
}

// grantHead is a GRANT or REVOKE read up to its grantee: a role list, or a
// privilege list and its target.
type grantHead struct {
	roles      []string
	privileges []string
	target     privTarget
}

// parseGrant parses grantRoleStatement and grantPrivilegeStatement. On
// entry, cur is the token after GRANT.
func (p *Parser) parseGrant(startLoc ast.Loc) (ast.Node, error) {
	head, err := p.parseGrantHead(kwTO)
	if err != nil {
		return nil, err
	}
	granteeType, grantee, err := p.parseGrantRevokeClause()
	if err != nil {
		return nil, err
	}
	stmt := &ast.GrantStmt{
		Roles:          head.roles,
		Privileges:     head.privileges,
		ObjectType:     head.target.objectType,
		Objects:        head.target.objects,
		FunctionArgs:   head.target.functionArgs,
		InDatabase:     head.target.inDatabase,
		InAllDatabases: head.target.inAllDatabases,
		ToType:         granteeType,
		Grantees:       []string{grantee},
	}
	if len(stmt.Objects) > 0 {
		stmt.Object = stmt.Objects[0]
	}
	if head.roles == nil && p.cur.Kind == kwWITH {
		// (WITH GRANT OPTION)?
		p.advance() // consume WITH
		if _, err := p.expect(kwGRANT); err != nil {
			return nil, err
		}
		if !p.atGrantWord("OPTION") {
			return nil, p.syntaxErrorAtCur()
		}
		p.advance()
		stmt.WithGrantOption = true
	}
	stmt.Loc = ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End}
	return stmt, nil
}

// parseRevoke parses revokeRoleStatement and revokePrivilegeStatement. On
// entry, cur is the token after REVOKE.
func (p *Parser) parseRevoke(startLoc ast.Loc) (ast.Node, error) {
	head, err := p.parseGrantHead(kwFROM)
	if err != nil {
		return nil, err
	}
	revokeeType, revokee, err := p.parseGrantRevokeClause()
	if err != nil {
		return nil, err
	}
	stmt := &ast.RevokeStmt{
		Roles:          head.roles,
		Privileges:     head.privileges,
		ObjectType:     head.target.objectType,
		Objects:        head.target.objects,
		FunctionArgs:   head.target.functionArgs,
		InDatabase:     head.target.inDatabase,
		InAllDatabases: head.target.inAllDatabases,
		FromType:       revokeeType,
		Revokees:       []string{revokee},
	}
	if len(stmt.Objects) > 0 {
		stmt.Object = stmt.Objects[0]
	}
	stmt.Loc = ast.Loc{Start: startLoc.Start, End: p.prev.Loc.End}
	return stmt, nil
}

// parseGrantHead reads a GRANT or REVOKE from its first item through ON and
// the target, or up to to (TO or FROM) for a role list, which it consumes.
//
// The two lists share their first tokens (USAGE is a privilege and may name
// a role), so the items are read as runs of tokens up to ',' and the word
// after the list decides: ON makes them privileges, to roles.
func (p *Parser) parseGrantHead(to int) (*grantHead, error) {
	var items [][]Token
	for {
		var item []Token
		for p.cur.Kind != int(',') && p.cur.Kind != kwON && p.cur.Kind != to &&
			p.cur.Kind != tokEOF && p.cur.Kind != int(';') {
			item = append(item, p.advance())
		}
		if len(item) == 0 {
			return nil, p.syntaxErrorAtCur()
		}
		items = append(items, item)
		if p.cur.Kind != int(',') {
			break
		}
		p.advance() // consume ','
	}

	head := &grantHead{}
	switch p.cur.Kind {
	case to:
		// identifierOrStringList
		for _, item := range items {
			if len(item) != 1 || !isGrantNameToken(item[0]) {
				return nil, grantSyntaxError(item[len(item)-1])
			}
			head.roles = append(head.roles, item[0].Str)
		}
		p.advance() // consume TO / FROM
		return head, nil
	case kwON:
		for _, item := range items {
			priv, ok := grantPrivilegeType(item)
			if !ok {
				return nil, grantSyntaxError(item[0])
			}
			head.privileges = append(head.privileges, priv)
		}
		p.advance() // consume ON
	default:
		return nil, p.syntaxErrorAtCur()
	}

	impersonate := len(head.privileges) == 1 && head.privileges[0] == "IMPERSONATE"
	if err := p.parsePrivTarget(&head.target, to, impersonate); err != nil {
		return nil, err
	}
	if _, err := p.expect(to); err != nil {
		return nil, err
	}
	return head, nil
}

// grantPrivilegeType matches item to privilegeType (StarRocks.g4:1743):
//
//	ALL PRIVILEGES? | ALTER | APPLY | BLACKLIST
//	| CREATE (DATABASE | TABLE | VIEW | FUNCTION | GLOBAL FUNCTION | MATERIALIZED VIEW
//	    | RESOURCE | RESOURCE GROUP | EXTERNAL CATALOG | STORAGE VOLUME | WAREHOUSE | PIPE)
//	| DELETE | DROP | EXPORT | FILE | IMPERSONATE | INSERT | GRANT | NODE | OPERATE
//	| PLUGIN | REPOSITORY | REFRESH | SELECT | UPDATE | USAGE
func grantPrivilegeType(item []Token) (string, bool) {
	words := make([]string, len(item))
	for i, tok := range item {
		if words[i] = grantWord(tok); words[i] == "" {
			return "", false
		}
	}
	text := strings.Join(words, " ")
	switch text {
	case "ALL", "ALL PRIVILEGES", "ALTER", "APPLY", "BLACKLIST",
		"DELETE", "DROP", "EXPORT", "FILE", "IMPERSONATE", "INSERT", "GRANT", "NODE",
		"OPERATE", "PLUGIN", "REPOSITORY", "REFRESH", "SELECT", "UPDATE", "USAGE",
		"CREATE DATABASE", "CREATE TABLE", "CREATE VIEW", "CREATE FUNCTION",
		"CREATE GLOBAL FUNCTION", "CREATE MATERIALIZED VIEW", "CREATE RESOURCE",
		"CREATE RESOURCE GROUP", "CREATE EXTERNAL CATALOG", "CREATE STORAGE VOLUME",
		"CREATE WAREHOUSE", "CREATE PIPE":
		return text, true
	}
	return "", false
}

// grantWord returns the upper-case spelling of an unquoted word, or "" for a
// quoted name, a string or punctuation.
func grantWord(tok Token) string {
	if tok.Kind == tokIdent || tok.Kind >= 700 {
		return strings.ToUpper(tok.Str)
	}
	return ""
}

// atGrantWord reports whether cur is the unquoted word w.
func (p *Parser) atGrantWord(w string) bool {
	return grantWord(p.cur) == w
}

// privObjectTypes are privObjectType's words (StarRocks.g4:1753), keyed by
// their first word; a two-word type lists its second.
var privObjectTypes = map[string]string{
	"CATALOG": "", "DATABASE": "", "MATERIALIZED": "VIEW", "RESOURCE": "",
	"STORAGE": "VOLUME", "SYSTEM": "", "TABLE": "", "VIEW": "", "WAREHOUSE": "", "PIPE": "",
}

// privObjectTypePlurals are privObjectTypePlural's words (StarRocks.g4:1758),
// keyed by their first word; a two-word type lists its second.
var privObjectTypePlurals = map[string]string{
	"CATALOGS": "", "DATABASES": "", "FUNCTIONS": "", "GLOBAL": "FUNCTIONS",
	"MATERIALIZED": "VIEWS", "POLICIES": "", "RESOURCES": "", "RESOURCE": "GROUPS",
	"STORAGE": "VOLUMES", "TABLES": "", "USERS": "", "VIEWS": "", "WAREHOUSES": "", "PIPES": "",
}

// parsePrivTarget reads what follows ON in a privilege GRANT or REVOKE:
//
//	USER user (',' user)*                                    (IMPERSONATE only)
//	| GLOBAL? FUNCTION privFunctionObjectNameList
//	| SYSTEM
//	| privObjectType privObjectNameList
//	| ALL privObjectTypePlural (IN ALL DATABASES | IN DATABASE identifierOrString)?
//	| privObjectNameList
//
// to is TO or FROM, the word after the target.
func (p *Parser) parsePrivTarget(t *privTarget, to int, impersonate bool) error {
	ends := func(tok Token) bool {
		return tok.Kind == to || tok.Kind == int(',') || tok.Kind == int('.') ||
			tok.Kind == tokEOF || tok.Kind == int(';')
	}
	word := grantWord(p.cur)
	switch {
	case impersonate && p.cur.Kind == kwUSER && !ends(p.peekNext()):
		p.advance() // consume USER
		t.objectType = "USER"
		for {
			user, loc, err := p.parseGrantUser()
			if err != nil {
				return err
			}
			t.objects = append(t.objects, &ast.ObjectName{Parts: []string{user}, Loc: loc})
			if p.cur.Kind != int(',') {
				return nil
			}
			p.advance() // consume ','
		}

	case p.cur.Kind == kwFUNCTION || (p.cur.Kind == kwGLOBAL && p.peekNext().Kind == kwFUNCTION):
		t.objectType = "FUNCTION"
		if p.cur.Kind == kwGLOBAL {
			t.objectType = "GLOBAL FUNCTION"
			p.advance()
		}
		p.advance() // consume FUNCTION
		// privFunctionObjectNameList (:1735): qualifiedName '(' typeList ')' (',' ...)*
		for {
			name, err := p.parseMultipartIdentifier()
			if err != nil {
				return err
			}
			args, _, err := p.parseFunctionTypeList()
			if err != nil {
				return err
			}
			t.objects = append(t.objects, name)
			t.functionArgs = append(t.functionArgs, args)
			if p.cur.Kind != int(',') {
				return nil
			}
			p.advance() // consume ','
		}

	case p.cur.Kind == kwSYSTEM && p.peekNext().Kind == to:
		p.advance() // consume SYSTEM
		t.objectType = "SYSTEM"
		return nil

	case p.cur.Kind == kwALL:
		p.advance() // consume ALL
		first := grantWord(p.cur)
		second, ok := privObjectTypePlurals[first]
		if !ok {
			return p.syntaxErrorAtCur()
		}
		p.advance()
		t.objectType = "ALL " + first
		if second != "" {
			if !p.atGrantWord(second) {
				return p.syntaxErrorAtCur()
			}
			p.advance()
			t.objectType += " " + second
		}
		if p.cur.Kind == kwIN {
			p.advance() // consume IN
			if p.cur.Kind == kwALL {
				p.advance() // consume ALL
				if _, err := p.expect(kwDATABASES); err != nil {
					return err
				}
				t.inAllDatabases = true
				return nil
			}
			if _, err := p.expect(kwDATABASE); err != nil {
				return err
			}
			db, _, err := p.parseGrantName()
			if err != nil {
				return err
			}
			t.inDatabase = db
		}
		return nil
	}

	// privObjectType privObjectNameList, when the words name a type and
	// do not end the object name list themselves.
	if second, ok := privObjectTypes[word]; ok {
		if word == "RESOURCE" && grantWord(p.peekNext()) == "GROUP" {
			second = "GROUP" // RESOURCE GROUP, beside RESOURCE alone
		}
		switch {
		case second != "" && grantWord(p.peekNext()) == second:
			p.advance()
			p.advance()
			t.objectType = word + " " + second
		case second == "" && !ends(p.peekNext()):
			p.advance()
			t.objectType = word
		}
	}
	return p.parsePrivObjectNameList(t)
}

// parsePrivObjectNameList reads privObjectNameList (StarRocks.g4:1731):
//
//	privObjectName (',' privObjectName)*
//	privObjectName : identifierOrStringOrStar ('.' identifierOrStringOrStar)?
func (p *Parser) parsePrivObjectNameList(t *privTarget) error {
	for {
		start := p.cur.Loc
		var parts []string
		for {
			if p.cur.Kind == int('*') {
				p.advance()
				parts = append(parts, "*")
			} else {
				name, _, err := p.parseGrantName()
				if err != nil {
					return err
				}
				parts = append(parts, name)
			}
			if len(parts) == 2 || p.cur.Kind != int('.') {
				break
			}
			p.advance() // consume '.'
		}
		t.objects = append(t.objects, &ast.ObjectName{
			Parts: parts,
			Loc:   ast.Loc{Start: start.Start, End: p.prev.Loc.End},
		})
		if p.cur.Kind != int(',') {
			return nil
		}
		p.advance() // consume ','
	}
}

// parseGrantRevokeClause reads the grantee after TO or FROM:
//
//	grantRevokeClause : USER? user | ROLE identifierOrString     (:1686)
//
// and, after a role list, grantRoleStatement's TO USER? user | TO ROLE
// identifierOrString, which is the same. USER and ROLE are nonReserved, so
// either may itself be the user's name. It returns USER, ROLE or "" and the
// grantee as written.
func (p *Parser) parseGrantRevokeClause() (string, string, error) {
	kind := ""
	if (p.cur.Kind == kwUSER || p.cur.Kind == kwROLE) && isGrantNameToken(p.peekNext()) {
		kind = strings.ToUpper(p.advance().Str)
	}
	if kind == "ROLE" {
		name, _, err := p.parseGrantName()
		return kind, name, err
	}
	user, _, err := p.parseGrantUser()
	return kind, user, err
}

// parseGrantUser reads user (StarRocks.g4:2892):
//
//	identifierOrString ('@' (identifierOrString | '[' identifierOrString ']'))?
func (p *Parser) parseGrantUser() (string, ast.Loc, error) {
	start := p.cur.Loc
	name, _, err := p.parseGrantName()
	if err != nil {
		return "", ast.NoLoc(), err
	}
	if p.cur.Kind == int('@') {
		p.advance() // consume '@'
		bracket := p.cur.Kind == int('[')
		if bracket {
			p.advance()
		}
		host, _, err := p.parseGrantName()
		if err != nil {
			return "", ast.NoLoc(), err
		}
		if bracket {
			if _, err := p.expect(int(']')); err != nil {
				return "", ast.NoLoc(), err
			}
			host = "[" + host + "]"
		}
		name += "@" + host
	}
	return name, ast.Loc{Start: start.Start, End: p.prev.Loc.End}, nil
}

// isGrantNameToken reports whether tok is an identifierOrString: a string,
// an identifier, or a word of StarRocks.g4's nonReserved.
func isGrantNameToken(tok Token) bool {
	return tok.Kind == tokString || isIdentifierToken(tok.Kind) || starRocksNonReservedWords[tok.Kind]
}

// parseGrantName reads identifierOrString, taking the nonReserved words
// omni reserves elsewhere (starRocksNonReservedWords).
func (p *Parser) parseGrantName() (string, ast.Loc, error) {
	if starRocksNonReservedWords[p.cur.Kind] {
		tok := p.advance()
		return tok.Str, tok.Loc, nil
	}
	return p.parseIdentifierOrString()
}

// grantSyntaxError is a syntax error at tok, a token already consumed.
func grantSyntaxError(tok Token) *ParseError {
	text := tok.Str
	if text == "" {
		text = TokenName(tok.Kind)
	}
	return &ParseError{Position: tok.Loc.Start, End: tok.Loc.End, Message: "syntax error at or near " + text}
}
